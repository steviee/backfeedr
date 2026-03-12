package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventRequest represents a single event payload
type EventRequest struct {
	Type         string                 `json:"type"`
	Name         string                 `json:"name,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	AppVersion   string                 `json:"app_version"`
	OSVersion    string                 `json:"os_version"`
	DeviceModel  string                 `json:"device_model"`
	DeviceIDHash string                 `json:"device_id_hash"`
	SessionID    string                 `json:"session_id,omitempty"`
	OccurredAt   time.Time              `json:"occurred_at"`
}

// EventResponse represents the API response for event submission
type EventResponse struct {
	ID string `json:"id"`
}

// BatchResponse represents the API response for batch submission
type BatchResponse struct {
	Count int      `json:"count"`
	IDs   []string `json:"ids"`
}

// SendEvent sends a single event to the API
func (c *Client) SendEvent(event *EventRequest) (*EventResponse, error) {
	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshal event: %w", err)
	}

	resp, err := c.doRequest("POST", "/api/v1/events", body)
	if err != nil {
		return nil, err
	}

	var result EventResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SendBatch sends multiple events to the API
func (c *Client) SendBatch(events []EventRequest) (*BatchResponse, error) {
	payload := struct {
		Events []EventRequest `json:"events"`
	}{
		Events: events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal batch: %w", err)
	}

	resp, err := c.doRequest("POST", "/api/v1/events/batch", body)
	if err != nil {
		return nil, err
	}

	var result BatchResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
