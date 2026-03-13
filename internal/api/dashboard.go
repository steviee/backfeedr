package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/steviee/backfeedr/internal/store"
)

// DashboardHandler provides dashboard data API
type DashboardHandler struct {
	crashStore *store.CrashStore
	appStore   *store.AppStore
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(crashStore *store.CrashStore, appStore *store.AppStore) *DashboardHandler {
	return &DashboardHandler{
		crashStore: crashStore,
		appStore:   appStore,
	}
}

// OverviewData represents dashboard overview statistics
type OverviewData struct {
	CrashFreeRate7d float64           `json:"crash_free_rate_7d"`
	DAU             int               `json:"dau"`
	Sessions7d      int               `json:"sessions_7d"`
	Crashes7d       int               `json:"crashes_7d"`
	TopCrashes      []CrashGroupStats `json:"top_crashes"`
}

// CrashGroupStats represents aggregated crash statistics
type CrashGroupStats struct {
	GroupHash       string `json:"group_hash"`
	ExceptionType   string `json:"exception_type"`
	ExceptionReason string `json:"exception_reason"`
	Count           int    `json:"count"`
}

// GetOverview returns dashboard overview data
func (h *DashboardHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	now := time.Now().UTC()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// Get crashes in last 7 days
	crashes, _ := h.crashStore.ListWithTimeRange(ctx, "", sevenDaysAgo, now, 1000)
	crashCount := len(crashes)

	// Get crash groups in last 7 days
	groups, err := h.crashStore.GetGroupsWithTimeRange(ctx, "", sevenDaysAgo, now)
	topCrashes := make([]CrashGroupStats, 0)
	if err == nil && len(groups) > 0 {
		for i, g := range groups {
			if i >= 5 { // Top 5 only
				break
			}
			topCrashes = append(topCrashes, CrashGroupStats{
				GroupHash:       g.GroupHash,
				ExceptionType:   g.ExceptionType,
				ExceptionReason: g.ExceptionReason,
				Count:           g.Count,
			})
		}
	}

	// Calculate crash-free rate (placeholder - needs events data)
	crashFreeRate := 100.0
	if crashCount > 0 {
		crashFreeRate = 98.0
	}

	data := OverviewData{
		CrashFreeRate7d: crashFreeRate,
		DAU:             0, // TODO: Get from daily_metrics
		Sessions7d:      0, // TODO: Get from daily_metrics
		Crashes7d:       crashCount,
		TopCrashes:      topCrashes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
