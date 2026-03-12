package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// CrashReport represents a crash report payload
type CrashReport struct {
	ExceptionType   string       `json:"exception_type"`
	ExceptionReason string       `json:"exception_reason,omitempty"`
	StackTrace      []StackFrame `json:"stack_trace,omitempty"`
	AppVersion      string       `json:"app_version"`
	BuildNumber     string       `json:"build_number,omitempty"`
	OSVersion       string       `json:"os_version"`
	DeviceModel     string       `json:"device_model"`
	DeviceIDHash    string       `json:"device_id_hash"`
	Locale          string       `json:"locale,omitempty"`
	FreeMemoryMB    int          `json:"free_memory_mb,omitempty"`
	FreeDiskMB      int          `json:"free_disk_mb,omitempty"`
	BatteryLevel    float64      `json:"battery_level,omitempty"`
	IsCharging      bool         `json:"is_charging,omitempty"`
	OccurredAt      time.Time    `json:"occurred_at"`
}

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	Frame  int     `json:"frame"`
	Symbol string  `json:"symbol"`
	File   *string `json:"file,omitempty"`
	Line   *int    `json:"line,omitempty"`
}

// CrashResponse represents the API response for crash submission
type CrashResponse struct {
	ID        string `json:"id"`
	GroupHash string `json:"group_hash"`
}

// SendCrash sends a crash report to the API
func (c *Client) SendCrash(crash *CrashReport) (*CrashResponse, error) {
	body, err := json.Marshal(crash)
	if err != nil {
		return nil, fmt.Errorf("marshal crash: %w", err)
	}

	resp, err := c.doRequest("POST", "/api/v1/crashes", body)
	if err != nil {
		return nil, err
	}

	var result CrashResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
