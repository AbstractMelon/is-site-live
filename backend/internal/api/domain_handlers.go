package api

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/abstractmelon/is-site-live/internal/models"
	"github.com/gin-gonic/gin"
)

// createCustomDomain creates a new custom domain
func (s *Server) createCustomDomain(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body
	var domainCreation models.CustomDomainCreation
	if err := c.ShouldBindJSON(&domainCreation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create custom domain
	var domain models.CustomDomain
	err := s.db.Pool.QueryRow(context.Background(), `
		INSERT INTO custom_domains (user_id, domain, verified)
		VALUES ($1, $2, false)
		RETURNING id, user_id, domain, verified, created_at, updated_at
	`, userID, domainCreation.Domain).Scan(
		&domain.ID,
		&domain.UserID,
		&domain.Domain,
		&domain.Verified,
		&domain.CreatedAt,
		&domain.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom domain"})
		return
	}

	// Return created domain
	c.JSON(http.StatusCreated, domain)
}

// getUserDomains gets all custom domains for the current user
func (s *Server) getUserDomains(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get domains from database
	rows, err := s.db.Pool.Query(context.Background(), `
		SELECT id, user_id, domain, verified, created_at, updated_at
		FROM custom_domains
		WHERE user_id = $1
		ORDER BY domain
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get custom domains"})
		return
	}
	defer rows.Close()

	var domains []models.CustomDomain
	for rows.Next() {
		var domain models.CustomDomain
		err := rows.Scan(
			&domain.ID,
			&domain.UserID,
			&domain.Domain,
			&domain.Verified,
			&domain.CreatedAt,
			&domain.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan custom domain"})
			return
		}
		domains = append(domains, domain)
	}

	// Return domains
	c.JSON(http.StatusOK, domains)
}

// deleteCustomDomain deletes a custom domain
func (s *Server) deleteCustomDomain(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get domain ID from URL
	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	// Delete domain
	result, err := s.db.Pool.Exec(context.Background(), `
		DELETE FROM custom_domains
		WHERE id = $1 AND user_id = $2
	`, domainID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete custom domain"})
		return
	}

	// Check if domain was found
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Custom domain not found"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "Custom domain deleted successfully"})
}

// verifyCustomDomain verifies a custom domain
func (s *Server) verifyCustomDomain(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get domain ID from URL
	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	// Get domain from database
	var domain models.CustomDomain
	err = s.db.Pool.QueryRow(context.Background(), `
		SELECT id, user_id, domain, verified, created_at, updated_at
		FROM custom_domains
		WHERE id = $1 AND user_id = $2
	`, domainID, userID).Scan(
		&domain.ID,
		&domain.UserID,
		&domain.Domain,
		&domain.Verified,
		&domain.CreatedAt,
		&domain.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Custom domain not found"})
		return
	}

	// Verify domain by checking CNAME record
	verified, err := verifyDomainCNAME(domain.Domain, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify domain"})
		return
	}

	// Update domain verification status
	err = s.db.Pool.QueryRow(context.Background(), `
		UPDATE custom_domains
		SET verified = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, user_id, domain, verified, created_at, updated_at
	`, verified, domain.ID).Scan(
		&domain.ID,
		&domain.UserID,
		&domain.Domain,
		&domain.Verified,
		&domain.CreatedAt,
		&domain.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update domain verification status"})
		return
	}

	// Return updated domain
	c.JSON(http.StatusOK, domain)
}

// getDomainDashboard gets the dashboard for a custom domain
func (s *Server) getDomainDashboard(c *gin.Context) {
	// Get domain from URL
	domainName := c.Param("domain")

	// Get domain from database
	var domain models.CustomDomain
	var userID int
	var username string
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT d.id, d.user_id, d.domain, d.verified, d.created_at, d.updated_at, u.id, u.username
		FROM custom_domains d
		JOIN users u ON d.user_id = u.id
		WHERE d.domain = $1 AND d.verified = true
	`, domainName).Scan(
		&domain.ID,
		&domain.UserID,
		&domain.Domain,
		&domain.Verified,
		&domain.CreatedAt,
		&domain.UpdatedAt,
		&userID,
		&username,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Custom domain not found or not verified"})
		return
	}

	// Get user's sites
	rows, err := s.db.Pool.Query(context.Background(), `
		SELECT id, user_id, name, url, created_at, updated_at
		FROM sites
		WHERE user_id = $1
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user's sites"})
		return
	}
	defer rows.Close()

	var sites []models.Site
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan site"})
			return
		}
		sites = append(sites, site)
	}

	// Return dashboard data
	c.JSON(http.StatusOK, gin.H{
		"domain": domain,
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
		"sites": sites,
	})
}

// verifyDomainCNAME verifies that a domain has a CNAME record pointing to the expected host
func verifyDomainCNAME(domain, expectedHost string) (bool, error) {
	// Look up CNAME records
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		// If there's an error, it might be that the CNAME doesn't exist
		return false, nil
	}

	// Remove trailing dot from CNAME if present
	cname = strings.TrimSuffix(cname, ".")
	
	// Check if CNAME matches expected host
	// We're being a bit flexible here - we'll check if the expected host is contained in the CNAME
	// This allows for subdomains and different domain formats
	return strings.Contains(cname, expectedHost) || strings.Contains(expectedHost, cname), nil
}
