package models

import "time"

// DailyMetrics holds aggregated metrics for a single day
type DailyMetrics struct {
	AppID         string    `json:"app_id"`
	Date          time.Time `json:"date"`
	Sessions      int       `json:"sessions"`
	UniqueDevices int       `json:"unique_devices"`
	Crashes       int       `json:"crashes"`
	Errors        int       `json:"errors"`
	AvgSessionSec float64   `json:"avg_session_sec"`
}

// RetentionMetrics holds retention data
type RetentionMetrics struct {
	Day1  float64 `json:"day_1"`
	Day7  float64 `json:"day_7"`
	Day30 float64 `json:"day_30"`
}

// OverviewMetrics holds dashboard overview data
type OverviewMetrics struct {
	CrashFreeRate7d float64       `json:"crash_free_rate_7d"`
	DAU             int           `json:"dau"`
	MAU             int           `json:"mau"`
	Sessions7d      int           `json:"sessions_7d"`
	Sessions30d     int           `json:"sessions_30d"`
	TopCrashes      []*CrashGroup `json:"top_crashes"`
}
