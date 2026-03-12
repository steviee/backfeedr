package store

import (
	"database/sql"
	"embed"
	"fmt"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB wraps sql.DB with backfeedr-specific operations
type DB struct {
	*sql.DB
}

// New opens a SQLite database connection
func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return &DB{db}, nil
}

// Migrate runs database migrations
func Migrate(db *DB) error {
	// For MVP, we'll use simple schema creation
	// In production, use golang-migrate
	schema := `
CREATE TABLE IF NOT EXISTS apps (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    bundle_id TEXT NOT NULL UNIQUE,
    api_key TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS crashes (
    id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL REFERENCES apps(id),
    group_hash TEXT NOT NULL,
    exception_type TEXT,
    exception_reason TEXT,
    stack_trace TEXT, -- JSON
    app_version TEXT,
    build_number TEXT,
    os_version TEXT,
    device_model TEXT,
    locale TEXT,
    free_memory_mb INTEGER,
    free_disk_mb INTEGER,
    battery_level REAL,
    is_charging BOOLEAN,
    occurred_at DATETIME,
    received_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_crashes_app_id ON crashes(app_id);
CREATE INDEX IF NOT EXISTS idx_crashes_group_hash ON crashes(group_hash);
CREATE INDEX IF NOT EXISTS idx_crashes_occurred_at ON crashes(occurred_at);

CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL REFERENCES apps(id),
    type TEXT NOT NULL,
    name TEXT,
    properties TEXT, -- JSON
    app_version TEXT,
    os_version TEXT,
    device_model TEXT,
    device_id_hash TEXT,
    session_id TEXT,
    occurred_at DATETIME,
    received_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_events_app_id ON events(app_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at ON events(occurred_at);

CREATE TABLE IF NOT EXISTS daily_metrics (
    app_id TEXT NOT NULL REFERENCES apps(id),
    date DATE NOT NULL,
    sessions INTEGER DEFAULT 0,
    unique_devices INTEGER DEFAULT 0,
    crashes INTEGER DEFAULT 0,
    errors INTEGER DEFAULT 0,
    avg_session_sec REAL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (app_id, date)
);
`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}
