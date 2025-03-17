package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	SMTP       SMTPConfig
	Monitoring MonitoringConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Address string
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	Secret string
}

// SMTPConfig holds the SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// MonitoringConfig holds the monitoring configuration
type MonitoringConfig struct {
	Interval time.Duration
	Workers  int
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Server config
	serverAddress := getEnv("SERVER_ADDRESS", ":8080")

	// Database config
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "isitlive")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "isitlive")

	// JWT config
	jwtSecret := getEnv("JWT_SECRET", "your_jwt_secret_key")

	// SMTP config
	smtpHost := getEnv("SMTP_HOST", "")
	smtpPortStr := getEnv("SMTP_PORT", "587")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	smtpUser := getEnv("SMTP_USER", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	smtpFrom := getEnv("SMTP_FROM", "")

	// Monitoring config
	monitoringIntervalStr := getEnv("MONITORING_INTERVAL", "60")
	monitoringInterval, _ := strconv.Atoi(monitoringIntervalStr)
	monitoringWorkersStr := getEnv("MONITORING_WORKERS", "10")
	monitoringWorkers, _ := strconv.Atoi(monitoringWorkersStr)

	return &Config{
		Server: ServerConfig{
			Address: serverAddress,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
		},
		SMTP: SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			User:     smtpUser,
			Password: smtpPassword,
			From:     smtpFrom,
		},
		Monitoring: MonitoringConfig{
			Interval: time.Duration(monitoringInterval) * time.Second,
			Workers:  monitoringWorkers,
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
