package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	stackJSON, err := json.Marshal(crash.StackTrace)
	if err != nil {
		return fmt.Errorf("marshal stack trace: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO crashes (
			id, app_id, group_hash, exception_type, exception_reason,
			stack_trace, app_version, build_number, os_version, device_model,
			locale, free_memory_mb, free_disk_mb, battery_level, is_charging,
			occurred_at, received_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		crash.ID, crash.AppID, crash.GroupHash, crash.ExceptionType,
		crash.ExceptionReason, string(stackJSON), crash.AppVersion,
		crash.BuildNumber, crash.OSVersion, crash.DeviceModel, crash.Locale,
		crash.FreeMemoryMB, crash.FreeDiskMB, crash.BatteryLevel,
		crash.IsCharging, crash.OccurredAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert crash: %w", err)
	}
	return nil
}

// GetByID retrieves a crash by ID
func (s *CrashStore) GetByID(ctx context.Context, id string) (*models.Crash, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, app_id, group_hash, exception_type, exception_reason,
			stack_trace, app_version, build_number, os_version, device_model,
			locale, free_memory_mb, free_disk_mb, battery_level, is_charging,
			occurred_at, received_at
		FROM crashes WHERE id = ?`, id)

	return s.scanCrash(row)
}

// ListWithTimeRange retrieves crashes within a time range
func (s *CrashStore) ListWithTimeRange(ctx context.Context, appID string, start, end time.Time, limit int) ([]*models.Crash, error) {
	var query string
	var args []interface{}

	if appID == "" {
		// Get all crashes (no app filter)
		query = `
			SELECT id, app_id, group_hash, exception_type, exception_reason,
				stack_trace, app_version, build_number, os_version, device_model,
				locale, free_memory_mb, free_disk_mb, battery_level, is_charging,
				occurred_at, received_at
			FROM crashes 
			WHERE occurred_at >= ? AND occurred_at <= ?
			ORDER BY occurred_at DESC
			LIMIT ?`
		args = []interface{}{start, end, limit}
	} else {
		query = `
			SELECT id, app_id, group_hash, exception_type, exception_reason,
				stack_trace, app_version, build_number, os_version, device_model,
				locale, free_memory_mb, free_disk_mb, battery_level, is_charging,
				occurred_at, received_at
			FROM crashes 
			WHERE app_id = ? AND occurred_at >= ? AND occurred_at <= ?
			ORDER BY occurred_at DESC
			LIMIT ?`
		args = []interface{}{appID, start, end, limit}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query crashes with time range: %w", err)
	}
	defer rows.Close()

	var crashes []*models.Crash
	for rows.Next() {
		crash, err := s.scanCrash(rows)
		if err != nil {
			return nil, err
		}
		crashes = append(crashes, crash)
	}

	return crashes, rows.Err()
}

// GetGroupsWithTimeRange returns crash groups within a time range
func (s *CrashStore) GetGroupsWithTimeRange(ctx context.Context, appID string, start, end time.Time) ([]*models.CrashGroup, error) {
	var query string
	var args []interface{}

	if appID == "" {
		query = `
			SELECT group_hash, exception_type, exception_reason,
				COUNT(*) as count,
				MIN(occurred_at) as first_seen,
				MAX(occurred_at) as last_seen,
				GROUP_CONCAT(DISTINCT app_version) as versions
			FROM crashes
			WHERE occurred_at >= ? AND occurred_at <= ?
			GROUP BY group_hash
			ORDER BY count DESC`
		args = []interface{}{start, end}
	} else {
		query = `
			SELECT group_hash, exception_type, exception_reason,
				COUNT(*) as count,
				MIN(occurred_at) as first_seen,
				MAX(occurred_at) as last_seen,
				GROUP_CONCAT(DISTINCT app_version) as versions
			FROM crashes
			WHERE app_id = ? AND occurred_at >= ? AND occurred_at <= ?
			GROUP BY group_hash
			ORDER BY count DESC`
		args = []interface{}{appID, start, end}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query crash groups with time range: %w", err)
	}
	defer rows.Close()

	var groups []*models.CrashGroup
	for rows.Next() {
		var g models.CrashGroup
		var versions string
		err := rows.Scan(
			&g.GroupHash, &g.ExceptionType, &g.ExceptionReason,
			&g.Count, &g.FirstSeen, &g.LastSeen, &versions)
		if err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		if versions != "" {
			g.AffectedVersions = []string{versions}
		}
		groups = append(groups, &g)
	}

	return groups, rows.Err()
}

// GetGroups returns crash groups for an app
func (s *CrashStore) GetGroups(ctx context.Context, appID string) ([]*models.CrashGroup, error) {
	query := `
		SELECT group_hash, exception_type, exception_reason,
			COUNT(*) as count,
			MIN(occurred_at) as first_seen,
			MAX(occurred_at) as last_seen,
			GROUP_CONCAT(DISTINCT app_version) as versions
		FROM crashes
		WHERE app_id = ?
		GROUP BY group_hash
		ORDER BY count DESC`

	rows, err := s.db.QueryContext(ctx, query, appID)
	if err != nil {
		return nil, fmt.Errorf("query crash groups: %w", err)
	}
	defer rows.Close()

	var groups []*models.CrashGroup
	for rows.Next() {
		var g models.CrashGroup
		var versions string
		err := rows.Scan(&g.GroupHash, &g.ExceptionType, &g.ExceptionReason,
			&g.Count, &g.FirstSeen, &g.LastSeen, &versions)
		if err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		// Parse versions
		if versions != "" {
			g.AffectedVersions = splitVersions(versions)
		}
		groups = append(groups, &g)
	}

	return groups, rows.Err()
}

func (s *CrashStore) scanCrash(scanner interface {
	Scan(...interface{}) error
}) (*models.Crash, error) {
	var c models.Crash
	var stackJSON string
	var receivedAt sql.NullTime

	err := scanner.Scan(
		&c.ID, &c.AppID, &c.GroupHash, &c.ExceptionType, &c.ExceptionReason,
		&stackJSON, &c.AppVersion, &c.BuildNumber, &c.OSVersion, &c.DeviceModel,
		&c.Locale, &c.FreeMemoryMB, &c.FreeDiskMB, &c.BatteryLevel,
		&c.IsCharging, &c.OccurredAt, &receivedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan crash: %w", err)
	}

	if receivedAt.Valid {
		c.ReceivedAt = receivedAt.Time
	}

	if stackJSON != "" {
		if err := json.Unmarshal([]byte(stackJSON), &c.StackTrace); err != nil {
			return nil, fmt.Errorf("unmarshal stack trace: %w", err)
		}
	}

	return &c, nil
}

// GenerateGroupHash creates a hash for crash grouping
func GenerateGroupHash(exceptionType string, frames []models.StackFrame) string {
	// Take top 3 app frames (frames with file reference)
	var appFrames []models.StackFrame
	for _, f := range frames {
		if f.File != nil && *f.File != "" {
			appFrames = append(appFrames, f)
			if len(appFrames) >= 3 {
				break
			}
		}
	}

	// Build hash input: exception_type:symbol1:symbol2:symbol3
	hashInput := exceptionType
	for _, f := range appFrames {
		hashInput += ":" + f.Symbol
	}

	h := sha256.New()
	h.Write([]byte(hashInput))
	return hex.EncodeToString(h.Sum(nil))[:16] // First 16 chars for brevity
}

func splitVersions(versions string) []string {
	// Simple split - SQLite GROUP_CONCAT uses ',' by default
	// This is a placeholder - real implementation might need more robust parsing
	return []string{versions}
}
