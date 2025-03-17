package database

import (
	"context"
	"fmt"

	"github.com/abstractmelon/is-site-live/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents a database connection
type DB struct {
	Pool *pgxpool.Pool
}

// Connect establishes a connection to the database
func Connect(cfg config.DatabaseConfig) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Migrate runs database migrations
func Migrate(db *DB) error {
	// Create users table
	_, err := db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Create sites table
	_, err = db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS sites (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			url VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id, name)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create sites table: %v", err)
	}

	// Create checks table
	_, err = db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS checks (
			id SERIAL PRIMARY KEY,
			site_id INTEGER REFERENCES sites(id) ON DELETE CASCADE,
			status_code INTEGER,
			response_time INTEGER, -- in milliseconds
			is_up BOOLEAN NOT NULL,
			error_message TEXT,
			checked_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create checks table: %v", err)
	}

	// Create custom_domains table
	_, err = db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS custom_domains (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			domain VARCHAR(255) UNIQUE NOT NULL,
			verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create custom_domains table: %v", err)
	}

	// Create TimescaleDB hypertable for checks if TimescaleDB extension is available
	_, err = db.Pool.Exec(context.Background(), `
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM pg_extension WHERE extname = 'timescaledb'
			) THEN
				-- Create hypertable if TimescaleDB is available
				PERFORM create_hypertable('checks', 'checked_at', if_not_exists => TRUE);
			END IF;
		END
		$$;
	`)
	if err != nil {
		return fmt.Errorf("failed to create hypertable: %v", err)
	}

	return nil
}
