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

// CrashTypes represents crash type distribution
type CrashTypes struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// DeviceDistribution represents crash distribution by device
type DeviceDistribution struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// GetDailyCrashes returns crashes per day for the last 7 days
func (h *MetricsHandler) GetDailyCrashes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	now := time.Now().UTC()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// Get crashes in last 7 days
	crashes, _ := h.crashStore.ListWithTimeRange(ctx, "", sevenDaysAgo, now, 1000)

	// Group by day
	dayCount := make(map[string]int)
	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	for i := 0; i < 7; i++ {
		day := sevenDaysAgo.AddDate(0, 0, i)
		dayCount[day.Format("Mon")] = 0
	}

	for _, c := range crashes {
		day := c.OccurredAt.Format("Mon")
		dayCount[day]++
	}

	data := make([]int, 7)
	for i, label := range labels {
		data[i] = dayCount[label]
	}

	result := DailyCrashes{
		Labels: labels,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetCrashTypes returns distribution by exception type
func (h *MetricsHandler) GetCrashTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all crashes
	crashes, _ := h.crashStore.ListWithTimeRange(ctx, "",
		time.Now().AddDate(0, 0, -90), // Last 90 days
		time.Now().UTC(),
		1000)

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

// GetDeviceDistribution returns crashes by device model
func (h *MetricsHandler) GetDeviceDistribution(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all crashes
	crashes, _ := h.crashStore.ListWithTimeRange(ctx, "",
		time.Now().AddDate(0, 0, -90),
		time.Now().UTC(),
		1000)

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
