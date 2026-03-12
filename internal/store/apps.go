package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

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
	
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO apps (id, name, bundle_id, api_key)
		VALUES (?, ?, ?, ?)`,
		app.ID, app.Name, app.BundleID, app.APIKey,
	)
	if err != nil {
		return fmt.Errorf("insert app: %w", err)
	}
	return nil
}

// GetByID retrieves an app by ID
func (s *AppStore) GetByID(ctx context.Context, id string) (*models.App, error) {
	// TODO: implement
	return nil, nil
}

// GetByAPIKey retrieves an app by API key
func (s *AppStore) GetByAPIKey(ctx context.Context, apiKey string) (*models.App, error) {
	// TODO: implement
	return nil, nil
}

// List retrieves all apps
func (s *AppStore) List(ctx context.Context) ([]*models.App, error) {
	// TODO: implement
	return nil, nil
}

// RotateAPIKey generates a new API key for an app
func (s *AppStore) RotateAPIKey(ctx context.Context, id string) (string, error) {
	// TODO: implement
	return "", nil
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
	// TODO: use proper ULID library
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
