package models

import "time"

// Crash represents a single crash report
type Crash struct {
	ID              string          `json:"id"`
	AppID           string          `json:"app_id"`
	GroupHash       string          `json:"group_hash"`
	ExceptionType   string          `json:"exception_type"`
	ExceptionReason string          `json:"exception_reason"`
	StackTrace      []StackFrame    `json:"stack_trace"`
	AppVersion      string          `json:"app_version"`
	BuildNumber     string          `json:"build_number"`
	OSVersion       string          `json:"os_version"`
	DeviceModel     string          `json:"device_model"`
	DeviceIDHash    string          `json:"device_id_hash"`
	Locale          string          `json:"locale"`
	FreeMemoryMB    int             `json:"free_memory_mb"`
	FreeDiskMB      int             `json:"free_disk_mb"`
	BatteryLevel    float64         `json:"battery_level"`
	IsCharging      bool            `json:"is_charging"`
	OccurredAt      time.Time       `json:"occurred_at"`
	ReceivedAt      time.Time       `json:"received_at"`
}

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	Frame  int     `json:"frame"`
	Symbol string  `json:"symbol"`
	File   *string `json:"file,omitempty"`
	Line   *int    `json:"line,omitempty"`
}

// CrashGroup represents grouped crashes by hash
type CrashGroup struct {
	GroupHash       string    `json:"group_hash"`
	ExceptionType   string    `json:"exception_type"`
	ExceptionReason string    `json:"exception_reason"`
	Count           int       `json:"count"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
	AffectedVersions []string `json:"affected_versions"`
}
