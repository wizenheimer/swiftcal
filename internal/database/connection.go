// internal/database/connection.go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizenheimer/swiftcal/pkg/logger"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewConnection(databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.GetLogger().Info("Database connection established")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

// internal/database/migrations/001_create_users.up.sql
/*
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    access_token TEXT,
    refresh_token TEXT,
    expiry_date TIMESTAMP WITH TIME ZONE,
    token_scope TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_expiry_date ON users(expiry_date);
*/

// internal/database/migrations/001_create_users.down.sql
/*
DROP TABLE IF EXISTS users;
*/

// internal/database/migrations/002_create_email_addresses.up.sql
/*
CREATE TABLE email_addresses (
    email VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_email_addresses_user_id ON email_addresses(user_id);
CREATE INDEX idx_email_addresses_default ON email_addresses(is_default) WHERE is_default = TRUE;
*/

// internal/database/migrations/002_create_email_addresses.down.sql
/*
DROP TABLE IF EXISTS email_addresses;
*/

// internal/database/migrations/003_create_pending_emails.up.sql
/*
CREATE TABLE pending_email_addresses (
    email VARCHAR(255) PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner_email VARCHAR(255) NOT NULL,
    verification_code UUID NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '24 hours')
);

CREATE INDEX idx_pending_emails_verification_code ON pending_email_addresses(verification_code);
CREATE INDEX idx_pending_emails_expires_at ON pending_email_addresses(expires_at);
*/

// internal/database/migrations/003_create_pending_emails.down.sql
/*
DROP TABLE IF EXISTS pending_email_addresses;
*/
