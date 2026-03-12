// Package worker provides background jobs for data processing
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/steviee/backfeedr/internal/store"
)

// MetricsWorker aggregates events into daily metrics
type MetricsWorker struct {
	store      *store.DB
	interval   time.Duration
	stopChan   chan struct{}
}

// NewMetricsWorker creates a new metrics aggregation worker
func NewMetricsWorker(db *store.DB) *MetricsWorker {
	return &MetricsWorker{
		store:    db,
		interval: 24 * time.Hour, // Run daily
		stopChan: make(chan struct{}),
	}
}

// Start begins the background aggregation
func (w *MetricsWorker) Start() {
	log.Println("[worker] Starting metrics aggregation worker")
	
	// Run immediately on start, then at interval
	w.runAggregation()
	
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.runAggregation()
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts the worker
func (w *MetricsWorker) Stop() {
	close(w.stopChan)
}

// runAggregation performs the daily aggregation
func (w *MetricsWorker) runAggregation() {
	ctx := context.Background()
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	
	log.Printf("[worker] Running aggregation for %s", yesterday.Format("2006-01-02"))
	
	if err := w.aggregateSessions(ctx, yesterday); err != nil {
		log.Printf("[worker] Error aggregating sessions: %v", err)
	}
	
	if err := w.aggregateCrashes(ctx, yesterday); err != nil {
		log.Printf("[worker] Error aggregating crashes: %v", err)
	}
	
	if err := w.calculateAvgSessionDuration(ctx, yesterday); err != nil {
		log.Printf("[worker] Error calculating session duration: %v", err)
	}
	
	log.Printf("[worker] Aggregation complete for %s", yesterday.Format("2006-01-02"))
}

// aggregateSessions counts sessions and unique devices per app/day
func (w *MetricsWorker) aggregateSessions(ctx context.Context, date time.Time) error {
	query := `
		INSERT INTO daily_metrics (app_id, date, sessions, unique_devices)
		SELECT 
			app_id,
			DATE(occurred_at) as date,
			COUNT(DISTINCT session_id) as sessions,
			COUNT(DISTINCT device_id_hash) as unique_devices
		FROM events
		WHERE type = 'session_start'
			AND DATE(occurred_at) = ?
		GROUP BY app_id, DATE(occurred_at)
		ON CONFLICT(app_id, date) DO UPDATE SET
			sessions = excluded.sessions,
			unique_devices = excluded.unique_devices,
			updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := w.store.ExecContext(ctx, query, date.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("aggregate sessions: %w", err)
	}
	return nil
}

// aggregateCrashes counts crashes per app/day
func (w *MetricsWorker) aggregateCrashes(ctx context.Context, date time.Time) error {
	query := `
		INSERT INTO daily_metrics (app_id, date, crashes)
		SELECT 
			app_id,
			DATE(occurred_at) as date,
			COUNT(*) as crashes
		FROM crashes
		WHERE DATE(occurred_at) = ?
		GROUP BY app_id, DATE(occurred_at)
		ON CONFLICT(app_id, date) DO UPDATE SET
			crashes = excluded.crashes,
			updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := w.store.ExecContext(ctx, query, date.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("aggregate crashes: %w", err)
	}
	return nil
}

// calculateAvgSessionDuration calculates average session duration per app/day
func (w *MetricsWorker) calculateAvgSessionDuration(ctx context.Context, date time.Time) error {
	query := `
		WITH session_durations AS (
			SELECT 
				e1.app_id,
				DATE(e1.occurred_at) as date,
				e2.occurred_at - e1.occurred_at as duration_sec
			FROM events e1
			JOIN events e2 ON e1.session_id = e2.session_id
			WHERE e1.type = 'session_start'
				AND e2.type = 'session_end'
				AND DATE(e1.occurred_at) = ?
		)
		INSERT INTO daily_metrics (app_id, date, avg_session_sec)
		SELECT 
			app_id,
			date,
			AVG(duration_sec) as avg_session_sec
		FROM session_durations
		GROUP BY app_id, date
		ON CONFLICT(app_id, date) DO UPDATE SET
			avg_session_sec = excluded.avg_session_sec,
			updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := w.store.ExecContext(ctx, query, date.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("calculate avg session duration: %w", err)
	}
	return nil
}

// AggregateDate manually triggers aggregation for a specific date
func (w *MetricsWorker) AggregateDate(ctx context.Context, date time.Time) error {
	log.Printf("[worker] Manual aggregation for %s", date.Format("2006-01-02"))
	
	if err := w.aggregateSessions(ctx, date); err != nil {
		return err
	}
	
	if err := w.aggregateCrashes(ctx, date); err != nil {
		return err
	}
	
	if err := w.calculateAvgSessionDuration(ctx, date); err != nil {
		return err
	}
	
	return nil
}

// AggregateDateRange aggregates metrics for a date range (backfill)
func (w *MetricsWorker) AggregateDateRange(ctx context.Context, start, end time.Time) error {
	log.Printf("[worker] Backfilling from %s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if err := w.AggregateDate(ctx, d); err != nil {
			return fmt.Errorf("aggregate %s: %w", d.Format("2006-01-02"), err)
		}
	}
	
	return nil
}
