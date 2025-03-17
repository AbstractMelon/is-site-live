package api

import (
	"context"
	"net/http"
	"time"

	"github.com/abstractmelon/is-site-live/internal/auth"
	"github.com/abstractmelon/is-site-live/internal/config"
	"github.com/abstractmelon/is-site-live/internal/database"
	"github.com/abstractmelon/is-site-live/internal/monitoring"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	router            *gin.Engine
	config            *config.Config
	db                *database.DB
	monitoringService *monitoring.Service
	httpServer        *http.Server
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, db *database.DB, monitoringService *monitoring.Service) *Server {
	// Create router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create server
	server := &Server{
		router:            router,
		config:            cfg,
		db:                db,
		monitoringService: monitoringService,
		httpServer: &http.Server{
			Addr:    cfg.Server.Address,
			Handler: router,
		},
	}

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes sets up the API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// Auth routes
	s.router.POST("/auth/register", s.register)
	s.router.POST("/auth/login", s.login)

	// Protected routes
	protected := s.router.Group("/")
	protected.Use(auth.AuthMiddleware(s.config.JWT))
	{
		// User routes
		protected.GET("/user", s.getCurrentUser)
		protected.PUT("/user", s.updateUser)

		// Site routes
		protected.POST("/sites", s.createSite)
		protected.GET("/sites", s.getUserSites)
		protected.GET("/sites/:id", s.getSite)
		protected.PUT("/sites/:id", s.updateSite)
		protected.DELETE("/sites/:id", s.deleteSite)

		// Custom domain routes
		protected.POST("/domains", s.createCustomDomain)
		protected.GET("/domains", s.getUserDomains)
		protected.DELETE("/domains/:id", s.deleteCustomDomain)
		protected.GET("/domains/:id/verify", s.verifyCustomDomain)
	}

	// Public routes
	s.router.GET("/user/:username", s.getUserProfile)
	s.router.GET("/site/:id/stats", s.getSiteStats)
	s.router.GET("/domain/:domain", s.getDomainDashboard)
}

// Start starts the API server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the API server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}
