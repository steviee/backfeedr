package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/steviee/backfeedr/internal/models"
	"github.com/steviee/backfeedr/internal/store"
)

// EventHandler handles event ingestion
type EventHandler struct {
	eventStore *store.EventStore
	appStore   *store.AppStore
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventStore *store.EventStore, appStore *store.AppStore) *EventHandler {
	return &EventHandler{
		eventStore: eventStore,
		appStore:   appStore,
	}
}

// HandleEvent handles POST /api/v1/events
func (h *EventHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	appID, ok := r.Context().Value("app_id").(string)
	if !ok || appID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req EventRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validate event type
	if !isValidEventType(req.Type) {
		http.Error(w, fmt.Sprintf(`{"error":"invalid event type: %s"}`, req.Type), http.StatusBadRequest)
		return
	}

	if req.OccurredAt.IsZero() {
		http.Error(w, `{"error":"missing occurred_at"}`, http.StatusBadRequest)
		return
	}

	event := &models.Event{
		ID:           generateEventULID(),
		AppID:        appID,
		Type:         req.Type,
		Name:         req.Name,
		Properties:   req.Properties,
		AppVersion:   req.AppVersion,
		OSVersion:    req.OSVersion,
		DeviceModel:  req.DeviceModel,
		DeviceIDHash: req.DeviceIDHash,
		SessionID:    req.SessionID,
		OccurredAt:   req.OccurredAt,
	}

	if err := h.eventStore.Create(r.Context(), event); err != nil {
		http.Error(w, `{"error":"failed to store event"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(EventResponse{
		ID: event.ID,
	})
}

// HandleBatch handles POST /api/v1/events/batch
func (h *EventHandler) HandleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	appID, ok := r.Context().Value("app_id").(string)
	if !ok || appID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req BatchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Limit batch size
	if len(req.Events) > 100 {
		http.Error(w, `{"error":"batch too large, max 100 events"}`, http.StatusBadRequest)
		return
	}

	var events []*models.Event
	for i, evReq := range req.Events {
		if !isValidEventType(evReq.Type) {
			http.Error(w, fmt.Sprintf(`{"error":"invalid event type at index %d: %s"}`, i, evReq.Type), http.StatusBadRequest)
			return
		}

		if evReq.OccurredAt.IsZero() {
			http.Error(w, fmt.Sprintf(`{"error":"missing occurred_at at index %d"}`, i), http.StatusBadRequest)
			return
		}

		events = append(events, &models.Event{
			ID:           generateEventULID(),
			AppID:        appID,
			Type:         evReq.Type,
			Name:         evReq.Name,
			Properties:   evReq.Properties,
			AppVersion:   evReq.AppVersion,
			OSVersion:    evReq.OSVersion,
			DeviceModel:  evReq.DeviceModel,
			DeviceIDHash: evReq.DeviceIDHash,
			SessionID:    evReq.SessionID,
			OccurredAt:   evReq.OccurredAt,
		})
	}

	if err := h.eventStore.CreateBatch(r.Context(), events); err != nil {
		http.Error(w, `{"error":"failed to store events"}`, http.StatusInternalServerError)
		return
	}

	ids := make([]string, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(BatchResponse{
		Count: len(ids),
		IDs:   ids,
	})
}

func isValidEventType(t string) bool {
	switch t {
	case models.EventTypeSessionStart,
		models.EventTypeSessionEnd,
		models.EventTypeError,
		models.EventTypeCustom:
		return true
	}
	return false
}

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

// EventResponse represents the API response
type EventResponse struct {
	ID string `json:"id"`
}

// BatchRequest represents a batch of events
type BatchRequest struct {
	Events []EventRequest `json:"events"`
}

// BatchResponse represents the batch API response
type BatchResponse struct {
	Count int      `json:"count"`
	IDs   []string `json:"ids"`
}

func generateEventULID() string {
	return fmt.Sprintf("%d%010d", time.Now().UnixMilli(), time.Now().Nanosecond())
}
