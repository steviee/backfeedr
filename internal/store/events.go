package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/steviee/backfeedr/internal/models"
)

// EventStore handles event-related database operations
type EventStore struct {
	db *DB
}

// NewEventStore creates a new event store
func NewEventStore(db *DB) *EventStore {
	return &EventStore{db: db}
}

// Create inserts a new event
func (s *EventStore) Create(ctx context.Context, event *models.Event) error {
	propsJSON, _ := json.Marshal(event.Properties)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO events (
			id, app_id, type, name, properties, app_version,
			os_version, device_model, session_id, occurred_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.AppID, event.Type, event.Name,
		string(propsJSON), event.AppVersion, event.OSVersion,
		event.DeviceModel, event.SessionID, event.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

// CreateBatch inserts multiple events in a transaction
func (s *EventStore) CreateBatch(ctx context.Context, events []*models.Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events (
			id, app_id, type, name, properties, app_version,
			os_version, device_model, session_id, occurred_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		propsJSON, _ := json.Marshal(event.Properties)
		if _, err := stmt.ExecContext(ctx,
			event.ID, event.AppID, event.Type, event.Name,
			string(propsJSON), event.AppVersion, event.OSVersion,
			event.DeviceModel, event.SessionID, event.OccurredAt,
		); err != nil {
			return fmt.Errorf("insert event: %w", err)
		}
	}

	return tx.Commit()
}

// List retrieves events with filtering
func (s *EventStore) List(ctx context.Context, appID string, eventType string, limit int) ([]*models.Event, error) {
	// TODO: implement
	return nil, nil
}
