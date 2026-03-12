// Package worker provides background jobs for data processing
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/steviee/backfeedr/internal/config"
	"github.com/steviee/backfeedr/internal/store"
)

// RetentionWorker handles data retention and cleanup
type RetentionWorker struct {
	store     *store.DB
	config    *config.Config
	stopChan  chan struct{}
}

// NewRetentionWorker creates a new retention worker
func NewRetentionWorker(db *store.DB, cfg *config.Config) *RetentionWorker {
	return &RetentionWorker{
		store:    db,
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

// Start begins the retention worker
func (w *RetentionWorker) Start() {
	log.Println("[retention] Starting data retention worker")
	
	// Schedule daily at 3 AM
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	
	delay := nextRun.Sub(now)
	log.Printf("[retention] First cleanup scheduled for %s (in %v)", nextRun.Format("2006-01-02 15:04"), delay)
	
	// Wait until 3 AM, then run daily
	go func() {
		time.Sleep(delay)
		w.runRetention()
		
		ticker := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-ticker.C:
				w.runRetention()
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts the worker
func (w *RetentionWorker) Stop() {
	close(w.stopChan)
}

// runRetention performs the cleanup
func (w *RetentionWorker) runRetention() {
	ctx := context.Background()
	cutoff := time.Now().UTC().AddDate(0, 0, -w.config.RetentionDays)
	
	log.Printf("[retention] Running cleanup for data older than %s (%d days)", 
		cutoff.Format("2006-01-02"), w.config.RetentionDays)
	
	// Delete old events
	eventsDeleted, err := w.deleteOldEvents(ctx, cutoff)
	if err != nil {
		log.Printf("[retention] Error deleting events: %v", err)
	} else {
		log.Printf("[retention] Deleted %d old events", eventsDeleted)
	}
	
	// Delete old crashes
	crashesDeleted, err := w.deleteOldCrashes(ctx, cutoff)
	if err != nil {
		log.Printf("[retention] Error deleting crashes: %v", err)
	} else {
		log.Printf("[retention] Deleted %d old crashes", crashesDeleted)
	}
	
	log.Printf("[retention] Cleanup complete. Total deleted: %d", eventsDeleted+crashesDeleted)
}

// deleteOldEvents removes events older than cutoff
func (w *RetentionWorker) deleteOldEvents(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := w.store.ExecContext(ctx, 
		"DELETE FROM events WHERE occurred_at < ?", cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete events: %w", err)
	}
	
	return result.RowsAffected()
}

// deleteOldCrashes removes crashes older than cutoff
func (w *RetentionWorker) deleteOldCrashes(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := w.store.ExecContext(ctx,
		"DELETE FROM crashes WHERE occurred_at < ?", cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete crashes: %w", err)
	}
	
	return result.RowsAffected()
}

// CleanupNow manually triggers cleanup
func (w *RetentionWorker) CleanupNow(ctx context.Context) error {
	cutoff := time.Now().UTC().AddDate(0, 0, -w.config.RetentionDays)
	
	eventsDeleted, err := w.deleteOldEvents(ctx, cutoff)
	if err != nil {
		return fmt.Errorf("delete events: %w", err)
	}
	
	crashesDeleted, err := w.deleteOldCrashes(ctx, cutoff)
	if err != nil {
		return fmt.Errorf("delete crashes: %w", err)
	}
	
	log.Printf("[retention] Manual cleanup complete: %d events, %d crashes deleted", 
		eventsDeleted, crashesDeleted)
	
	return nil
}
