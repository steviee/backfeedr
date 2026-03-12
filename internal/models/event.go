package models

import "time"

// Event represents an analytics event
type Event struct {
	ID           string                 `json:"id"`
	AppID        string                 `json:"app_id"`
	Type         string                 `json:"type"` // session_start, session_end, error, custom
	Name         string                 `json:"name,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	AppVersion   string                 `json:"app_version"`
	OSVersion    string                 `json:"os_version"`
	DeviceModel  string                 `json:"device_model"`
	DeviceIDHash string                 `json:"device_id_hash"`
	SessionID    string                 `json:"session_id,omitempty"`
	OccurredAt   time.Time              `json:"occurred_at"`
	ReceivedAt   time.Time              `json:"received_at"`
}

// Event types
const (
	EventTypeSessionStart = "session_start"
	EventTypeSessionEnd   = "session_end"
	EventTypeError        = "error"
	EventTypeCustom       = "custom"
)
