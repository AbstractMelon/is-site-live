package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abstractmelon/is-site-live/internal/models"
	"github.com/gin-gonic/gin"
)

// createSite creates a new site
func (s *Server) createSite(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body
	var siteCreation models.SiteCreation
	if err := c.ShouldBindJSON(&siteCreation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create site
	var site models.Site
	err := s.db.Pool.QueryRow(context.Background(), `
		INSERT INTO sites (user_id, name, url)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, url, created_at, updated_at
	`, userID, siteCreation.Name, siteCreation.URL).Scan(
		&site.ID,
		&site.UserID,
		&site.Name,
		&site.URL,
		&site.CreatedAt,
		&site.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site"})
		return
	}

	// Return created site
	c.JSON(http.StatusCreated, site)
}

// getUserSites gets all sites for the current user
func (s *Server) getUserSites(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get sites from database
	rows, err := s.db.Pool.Query(context.Background(), `
		SELECT id, user_id, name, url, created_at, updated_at
		FROM sites
		WHERE user_id = $1
		ORDER BY name
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sites"})
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

	// Return sites
	c.JSON(http.StatusOK, sites)
}

// getSite gets a site by ID
func (s *Server) getSite(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get site ID from URL
	siteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}

	// Get site from database
	var site models.Site
	err = s.db.Pool.QueryRow(context.Background(), `
		SELECT id, user_id, name, url, created_at, updated_at
		FROM sites
		WHERE id = $1 AND user_id = $2
	`, siteID, userID).Scan(
		&site.ID,
		&site.UserID,
		&site.Name,
		&site.URL,
		&site.CreatedAt,
		&site.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Return site
	c.JSON(http.StatusOK, site)
}

// updateSite updates a site
func (s *Server) updateSite(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get site ID from URL
	siteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}

	// Bind request body
	var siteUpdate models.SiteCreation
	if err := c.ShouldBindJSON(&siteUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update site
	var site models.Site
	err = s.db.Pool.QueryRow(context.Background(), `
		UPDATE sites
		SET name = $1, url = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, url, created_at, updated_at
	`, siteUpdate.Name, siteUpdate.URL, siteID, userID).Scan(
		&site.ID,
		&site.UserID,
		&site.Name,
		&site.URL,
		&site.CreatedAt,
		&site.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Return updated site
	c.JSON(http.StatusOK, site)
}

// deleteSite deletes a site
func (s *Server) deleteSite(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get site ID from URL
	siteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}

	// Delete site
	result, err := s.db.Pool.Exec(context.Background(), `
		DELETE FROM sites
		WHERE id = $1 AND user_id = $2
	`, siteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete site"})
		return
	}

	// Check if site was found
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}

// getSiteStats gets the stats for a site
func (s *Server) getSiteStats(c *gin.Context) {
	// Get site ID from URL
	siteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}

	// Get site stats
	siteWithStats, err := s.monitoringService.GetSiteStats(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Return site stats
	c.JSON(http.StatusOK, siteWithStats)
}
