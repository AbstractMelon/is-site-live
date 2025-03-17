package monitoring

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/abstractmelon/is-site-live/internal/database"
	"github.com/abstractmelon/is-site-live/internal/models"
)

// Service handles the monitoring of sites
type Service struct {
	db           *database.DB
	client       *http.Client
	stopChan     chan struct{}
	wg           sync.WaitGroup
	sitesCache   map[int]models.Site
	sitesCacheMu sync.RWMutex
}

// NewService creates a new monitoring service
func NewService(db *database.DB) *Service {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// Parse the address to get the host
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					// If SplitHostPort fails, assume the address is just a host
					host = addr
				}

				// Check if the host is a private IP
				ip := net.ParseIP(host)
				if ip != nil && (ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()) {
					return nil, fmt.Errorf("connection to private IP %s is not allowed", host)
				}

				// Use the default dialer
				dialer := &net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}

	return &Service{
		db:         db,
		client:     client,
		stopChan:   make(chan struct{}),
		sitesCache: make(map[int]models.Site),
	}
}

// StartWorkerPool starts the worker pool for monitoring sites
func (s *Service) StartWorkerPool(numWorkers int, checkInterval time.Duration) {
	// Start the scheduler
	s.wg.Add(1)
	go s.scheduler(checkInterval)

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// Stop stops the monitoring service
func (s *Service) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// scheduler schedules site checks at regular intervals
func (s *Service) scheduler(interval time.Duration) {
	defer s.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Channel for site check tasks
	checkChan := make(chan models.Site)

	// Start a goroutine to send sites to the check channel
	go func() {
		for {
			select {
			case <-s.stopChan:
				close(checkChan)
				return
			case <-ticker.C:
				// Get all sites from the database
				sites, err := s.getAllSites()
				if err != nil {
					fmt.Printf("Error getting sites: %v\n", err)
					continue
				}

				// Update the sites cache
				s.updateSitesCache(sites)

				// Send each site to the check channel
				for _, site := range sites {
					select {
					case <-s.stopChan:
						close(checkChan)
						return
					case checkChan <- site:
						// Site sent for checking
					}
				}
			}
		}
	}()

	// Process site check tasks
	for site := range checkChan {
		select {
		case <-s.stopChan:
			return
		default:
			// Check the site
			s.checkSite(site)
		}
	}
}

// worker processes site check tasks
func (s *Service) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		default:
			// Wait for the next task
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// getAllSites gets all sites from the database
func (s *Service) getAllSites() ([]models.Site, error) {
	var sites []models.Site

	rows, err := s.db.Pool.Query(context.Background(), `
		SELECT id, user_id, name, url, created_at, updated_at
		FROM sites
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var site models.Site
		err := rows.Scan(
			&site.ID,
			&site.UserID,
			&site.Name,
			&site.URL,
			&site.CreatedAt,
			&site.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	return sites, nil
}

// updateSitesCache updates the sites cache
func (s *Service) updateSitesCache(sites []models.Site) {
	s.sitesCacheMu.Lock()
	defer s.sitesCacheMu.Unlock()

	// Clear the cache
	s.sitesCache = make(map[int]models.Site)

	// Add sites to the cache
	for _, site := range sites {
		s.sitesCache[site.ID] = site
	}
}

// checkSite checks a site and records the result
func (s *Service) checkSite(site models.Site) {
	// Parse the URL
	parsedURL, err := url.Parse(site.URL)
	if err != nil {
		s.recordCheckResult(site.ID, 0, 0, false, fmt.Sprintf("Invalid URL: %v", err))
		return
	}

	// Ensure the URL has a scheme
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		s.recordCheckResult(site.ID, 0, 0, false, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	// Set a user agent
	req.Header.Set("User-Agent", "IsItLive Monitoring/1.0")

	// Start the timer
	startTime := time.Now()

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		s.recordCheckResult(site.ID, 0, 0, false, fmt.Sprintf("Request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	// Calculate response time
	responseTime := int(time.Since(startTime).Milliseconds())

	// Determine if the site is up (2xx or 3xx status codes)
	isUp := resp.StatusCode >= 200 && resp.StatusCode < 400

	// Record the result
	s.recordCheckResult(site.ID, resp.StatusCode, responseTime, isUp, "")
}

// recordCheckResult records a check result in the database
func (s *Service) recordCheckResult(siteID, statusCode, responseTime int, isUp bool, errorMessage string) {
	_, err := s.db.Pool.Exec(context.Background(), `
		INSERT INTO checks (site_id, status_code, response_time, is_up, error_message, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, siteID, statusCode, responseTime, isUp, errorMessage, time.Now())
	if err != nil {
		fmt.Printf("Error recording check result: %v\n", err)
	}
}

// GetSiteStats gets the uptime statistics for a site
func (s *Service) GetSiteStats(siteID int) (*models.SiteWithStats, error) {
	// Get the site from the cache
	s.sitesCacheMu.RLock()
	site, ok := s.sitesCache[siteID]
	s.sitesCacheMu.RUnlock()

	if !ok {
		// If not in cache, get from database
		err := s.db.Pool.QueryRow(context.Background(), `
			SELECT id, user_id, name, url, created_at, updated_at
			FROM sites
			WHERE id = $1
		`, siteID).Scan(
			&site.ID,
			&site.UserID,
			&site.Name,
			&site.URL,
			&site.CreatedAt,
			&site.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// Get the current status (latest check)
	var currentStatus models.Check
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT id, site_id, status_code, response_time, is_up, error_message, checked_at
		FROM checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`, siteID).Scan(
		&currentStatus.ID,
		&currentStatus.SiteID,
		&currentStatus.StatusCode,
		&currentStatus.ResponseTime,
		&currentStatus.IsUp,
		&currentStatus.ErrorMessage,
		&currentStatus.CheckedAt,
	)
	if err != nil {
		// If no checks yet, return site with empty stats
		return &models.SiteWithStats{
			Site:           site,
			CurrentStatus:  nil,
			LifetimeStats:  models.UptimeStats{},
			Last7DaysStats: models.UptimeStats{},
			Last30DaysStats: models.UptimeStats{},
			Last90DaysStats: models.UptimeStats{},
		}, nil
	}

	// Get lifetime stats
	lifetimeStats, err := s.getUptimeStats(siteID, 0)
	if err != nil {
		return nil, err
	}

	// Get 7-day stats
	last7DaysStats, err := s.getUptimeStats(siteID, 7)
	if err != nil {
		return nil, err
	}

	// Get 30-day stats
	last30DaysStats, err := s.getUptimeStats(siteID, 30)
	if err != nil {
		return nil, err
	}

	// Get 90-day stats
	last90DaysStats, err := s.getUptimeStats(siteID, 90)
	if err != nil {
		return nil, err
	}

	return &models.SiteWithStats{
		Site:           site,
		CurrentStatus:  &currentStatus,
		LifetimeStats:  lifetimeStats,
		Last7DaysStats: last7DaysStats,
		Last30DaysStats: last30DaysStats,
		Last90DaysStats: last90DaysStats,
	}, nil
}

// getUptimeStats gets the uptime statistics for a site over a period of days
// If days is 0, get lifetime stats
func (s *Service) getUptimeStats(siteID, days int) (models.UptimeStats, error) {
	var stats models.UptimeStats
	var query string
	var args []interface{}

	if days > 0 {
		query = `
			SELECT
				COUNT(*) AS total_checks,
				SUM(CASE WHEN is_up THEN 1 ELSE 0 END) AS successful_checks,
				CASE WHEN COUNT(*) > 0 THEN
					(SUM(CASE WHEN is_up THEN 1 ELSE 0 END)::float / COUNT(*)) * 100
				ELSE 0 END AS uptime_percentage,
				CASE WHEN SUM(CASE WHEN is_up THEN 1 ELSE 0 END) > 0 THEN
					SUM(CASE WHEN is_up THEN response_time ELSE 0 END) / SUM(CASE WHEN is_up THEN 1 ELSE 0 END)
				ELSE 0 END AS average_response_time
			FROM checks
			WHERE site_id = $1 AND checked_at >= NOW() - INTERVAL '1 day' * $2
		`
		args = []interface{}{siteID, days}
	} else {
		query = `
			SELECT
				COUNT(*) AS total_checks,
				SUM(CASE WHEN is_up THEN 1 ELSE 0 END) AS successful_checks,
				CASE WHEN COUNT(*) > 0 THEN
					(SUM(CASE WHEN is_up THEN 1 ELSE 0 END)::float / COUNT(*)) * 100
				ELSE 0 END AS uptime_percentage,
				CASE WHEN SUM(CASE WHEN is_up THEN 1 ELSE 0 END) > 0 THEN
					SUM(CASE WHEN is_up THEN response_time ELSE 0 END) / SUM(CASE WHEN is_up THEN 1 ELSE 0 END)
				ELSE 0 END AS average_response_time
			FROM checks
			WHERE site_id = $1
		`
		args = []interface{}{siteID}
	}

	err := s.db.Pool.QueryRow(context.Background(), query, args...).Scan(
		&stats.TotalChecks,
		&stats.SuccessfulChecks,
		&stats.UptimePercentage,
		&stats.AverageResponseTime,
	)
	if err != nil {
		return models.UptimeStats{}, err
	}

	return stats, nil
}
