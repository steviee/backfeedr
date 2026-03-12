package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/steviee/backfeedr/internal/models"
)

// AppStore handles app-related database operations
type AppStore struct {
	db *DB
}

// NewAppStore creates a new app store
func NewAppStore(db *DB) *AppStore {
	return &AppStore{db: db}
}

// Create creates a new app with a generated API key
func (s *AppStore) Create(ctx context.Context, app *models.App) error {
	if app.ID == "" {
		app.ID = generateULID()
	}
	if app.APIKey == "" {
		app.APIKey = generateAPIKey(true) // live key by default
	}
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now().UTC()
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO apps (id, name, bundle_id, api_key, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		app.ID, app.Name, app.BundleID, app.APIKey, app.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert app: %w", err)
	}
	return nil
}

// GetByID retrieves an app by ID
func (s *AppStore) GetByID(ctx context.Context, id string) (*models.App, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, bundle_id, api_key, created_at
		FROM apps WHERE id = ?`, id)

	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.BundleID, &app.APIKey, &app.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan app: %w", err)
	}
	return &app, nil
}

// GetByAPIKey retrieves an app by API key
func (s *AppStore) GetByAPIKey(ctx context.Context, apiKey string) (*models.App, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, bundle_id, api_key, created_at
		FROM apps WHERE api_key = ?`, apiKey)

	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.BundleID, &app.APIKey, &app.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan app: %w", err)
	}
	return &app, nil
}

// List retrieves all apps
func (s *AppStore) List(ctx context.Context) ([]*models.App, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, bundle_id, api_key, created_at
		FROM apps ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query apps: %w", err)
	}
	defer rows.Close()

	var apps []*models.App
	for rows.Next() {
		var app models.App
		err := rows.Scan(&app.ID, &app.Name, &app.BundleID, &app.APIKey, &app.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan app: %w", err)
		}
		apps = append(apps, &app)
	}

	return apps, rows.Err()
}

// RotateAPIKey generates a new API key for an app
func (s *AppStore) RotateAPIKey(ctx context.Context, id string) (string, error) {
	newKey := generateAPIKey(true)
	res, err := s.db.ExecContext(ctx, `
		UPDATE apps SET api_key = ? WHERE id = ?`, newKey, id)
	if err != nil {
		return "", fmt.Errorf("update api key: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return "", fmt.Errorf("app not found")
	}
	return newKey, nil
}

func generateAPIKey(isLive bool) string {
	prefix := "bf_test_"
	if isLive {
		prefix = "bf_live_"
	}

	b := make([]byte, 16)
	rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

func generateULID() string {
	// Simple ULID-like generation using timestamp + random
	// For production, use a proper ULID library
	now := time.Now().UnixMilli()
	timestamp := fmt.Sprintf("%010x", now)

	b := make([]byte, 10)
	rand.Read(b)
	random := hex.EncodeToString(b)

	return timestamp + random
}
