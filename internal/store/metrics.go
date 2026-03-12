package store

import (
	"context"
	"time"

	"github.com/steviee/backfeedr/internal/models"
)

// MetricsStore handles metrics aggregation
type MetricsStore struct {
	db *DB
}

// NewMetricsStore creates a new metrics store
func NewMetricsStore(db *DB) *MetricsStore {
	return &MetricsStore{db: db}
}

// GetDailyMetrics retrieves metrics for a date range
func (s *MetricsStore) GetDailyMetrics(ctx context.Context, appID string, start, end time.Time) ([]*models.DailyMetrics, error) {
	// TODO: implement
	return nil, nil
}

// AggregateDay computes daily metrics for a specific day
func (s *MetricsStore) AggregateDay(ctx context.Context, appID string, date time.Time) error {
	// TODO: implement aggregation
	return nil
}

// GetCrashFreeRate calculates crash-free rate for a period
func (s *MetricsStore) GetCrashFreeRate(ctx context.Context, appID string, days int) (float64, error) {
	// TODO: implement
	return 100.0, nil
}
