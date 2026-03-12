package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/steviee/backfeedr/internal/models"
)

// CrashStore handles crash-related database operations
type CrashStore struct {
	db *DB
}

// NewCrashStore creates a new crash store
func NewCrashStore(db *DB) *CrashStore {
	return &CrashStore{db: db}
}

// Create inserts a new crash report
func (s *CrashStore) Create(ctx context.Context, crash *models.Crash) error {
	stackJSON, _ := json.Marshal(crash.StackTrace)
	
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO crashes (
			id, app_id, group_hash, exception_type, exception_reason,
			stack_trace, app_version, build_number, os_version, device_model,
			locale, free_memory_mb, free_disk_mb, battery_level, is_charging,
			occurred_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		crash.ID, crash.AppID, crash.GroupHash, crash.ExceptionType,
		crash.ExceptionReason, string(stackJSON), crash.AppVersion,
		crash.BuildNumber, crash.OSVersion, crash.DeviceModel, crash.Locale,
		crash.FreeMemoryMB, crash.FreeDiskMB, crash.BatteryLevel,
		crash.IsCharging, crash.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("insert crash: %w", err)
	}
	return nil
}

// GetByID retrieves a crash by ID
func (s *CrashStore) GetByID(ctx context.Context, id string) (*models.Crash, error) {
	// TODO: implement
	return nil, nil
}

// List retrieves crashes with optional filtering
func (s *CrashStore) List(ctx context.Context, appID string, limit int) ([]*models.Crash, error) {
	// TODO: implement
	return nil, nil
}

// GetGroups returns crash groups for an app
func (s *CrashStore) GetGroups(ctx context.Context, appID string) ([]*models.CrashGroup, error) {
	// TODO: implement
	return nil, nil
}
