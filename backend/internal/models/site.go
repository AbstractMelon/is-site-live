package models

import (
	"time"
)

// Site represents a monitored website
type Site struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SiteCreation represents the data needed to create a new site
type SiteCreation struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	URL  string `json:"url" binding:"required,url"`
}

// Check represents a single uptime check for a site
type Check struct {
	ID           int       `json:"id"`
	SiteID       int       `json:"site_id"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int       `json:"response_time"` // in milliseconds
	IsUp         bool      `json:"is_up"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
}

// UptimeStats represents uptime statistics for a site
type UptimeStats struct {
	TotalChecks     int     `json:"total_checks"`
	SuccessfulChecks int     `json:"successful_checks"`
	UptimePercentage float64 `json:"uptime_percentage"`
	AverageResponseTime int  `json:"average_response_time"` // in milliseconds
}

// SiteWithStats represents a site with its uptime statistics
type SiteWithStats struct {
	Site        Site        `json:"site"`
	CurrentStatus *Check     `json:"current_status"`
	LifetimeStats UptimeStats `json:"lifetime_stats"`
	Last7DaysStats UptimeStats `json:"last_7_days_stats"`
	Last30DaysStats UptimeStats `json:"last_30_days_stats"`
	Last90DaysStats UptimeStats `json:"last_90_days_stats"`
}

// CustomDomain represents a custom domain for a user's dashboard
type CustomDomain struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Domain    string    `json:"domain"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CustomDomainCreation represents the data needed to create a new custom domain
type CustomDomainCreation struct {
	Domain string `json:"domain" binding:"required,fqdn"`
}
