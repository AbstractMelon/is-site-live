package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abstractmelon/is-site-live/internal/api"
	"github.com/abstractmelon/is-site-live/internal/config"
	"github.com/abstractmelon/is-site-live/internal/database"
	"github.com/abstractmelon/is-site-live/internal/monitoring"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Create monitoring service
	monitoringService := monitoring.NewService(db)

	// Start the monitoring worker pool
	monitoringService.StartWorkerPool(cfg.Monitoring.Workers, cfg.Monitoring.Interval)

	// Create and setup API server
	server := api.NewServer(cfg, db, monitoringService)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", cfg.Server.Address)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Stop the monitoring service
	monitoringService.Stop()

	log.Println("Server exited properly")
}
