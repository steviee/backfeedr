package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/steviee/backfeedr/internal/store"
)

// MetricsHandler provides metrics data for charts
type MetricsHandler struct {
	crashStore *store.CrashStore
	appStore   *store.AppStore
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(crashStore *store.CrashStore, appStore *store.AppStore) *MetricsHandler {
	return &MetricsHandler{
		crashStore: crashStore,
		appStore:   appStore,
	}
}

// DailyCrashes represents daily crash counts
type DailyCrashes struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// GetDailyCrashes returns crashes per day for the last 7 days
func (h *MetricsHandler) GetDailyCrashes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get last 7 days
	days := 7
	today := time.Now().UTC().Truncate(24 * time.Hour)

	labels := make([]string, days)
	data := make([]int, days)

	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, -days+i+1)
		labels[i] = date.Format("Mon")

		// Start/end of day
		startOfDay := date
		endOfDay := date.Add(24 * time.Hour)

		// Count crashes in range
		crashes, _ := h.crashStore.List(ctx, "", 1000)
		count := 0
		for _, c := range crashes {
			if c.OccurredAt.After(startOfDay) && c.OccurredAt.Before(endOfDay) {
				count++
			}
		}
		data[i] = count
	}

	result := DailyCrashes{
		Labels: labels,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CrashTypes represents crash type distribution
type CrashTypes struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// GetCrashTypes returns distribution by exception type
func (h *MetricsHandler) GetCrashTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	crashes, _ := h.crashStore.List(ctx, "", 1000)

	typeCount := make(map[string]int)
	for _, c := range crashes {
		typeCount[c.ExceptionType]++
	}

	labels := make([]string, 0, len(typeCount))
	data := make([]int, 0, len(typeCount))

	for t, count := range typeCount {
		labels = append(labels, t)
		data = append(data, count)
	}

	result := CrashTypes{
		Labels: labels,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeviceDistribution represents crash distribution by device
type DeviceDistribution struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// GetDeviceDistribution returns crashes by device model
func (h *MetricsHandler) GetDeviceDistribution(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	crashes, _ := h.crashStore.List(ctx, "", 1000)

	deviceCount := make(map[string]int)
	for _, c := range crashes {
		deviceCount[c.DeviceModel]++
	}

	labels := make([]string, 0, len(deviceCount))
	data := make([]int, 0, len(deviceCount))

	for device, count := range deviceCount {
		if len(labels) < 10 { // Top 10 only
			labels = append(labels, device)
			data = append(data, count)
		}
	}

	result := DeviceDistribution{
		Labels: labels,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
