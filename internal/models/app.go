package models

import "time"

// App represents a registered iOS application
type App struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BundleID  string    `json:"bundle_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

// IsTestKey returns true if this is a test API key
func (a *App) IsTestKey() bool {
	return len(a.APIKey) > 8 && a.APIKey[:8] == "bf_test_"
}
