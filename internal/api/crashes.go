package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/steviee/backfeedr/internal/models"
	"github.com/steviee/backfeedr/internal/store"
)

// CrashHandler handles crash ingestion
type CrashHandler struct {
	crashStore *store.CrashStore
	appStore   *store.AppStore
}

// NewCrashHandler creates a new crash handler
func NewCrashHandler(crashStore *store.CrashStore, appStore *store.AppStore) *CrashHandler {
	return &CrashHandler{
		crashStore: crashStore,
		appStore:   appStore,
	}
}

// HandleCrash handles POST /api/v1/crashes
func (h *CrashHandler) HandleCrash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get app ID from context (set by auth middleware)
	appID, ok := r.Context().Value("app_id").(string)
	if !ok || appID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Read and parse body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req CrashRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ExceptionType == "" {
		http.Error(w, `{"error":"missing exception_type"}`, http.StatusBadRequest)
		return
	}
	if req.OccurredAt.IsZero() {
		http.Error(w, `{"error":"missing occurred_at"}`, http.StatusBadRequest)
		return
	}

	// Generate group hash
	groupHash := store.GenerateGroupHash(req.ExceptionType, req.StackTrace)

	// Create crash model
	crash := &models.Crash{
		ID:              generateULID(),
		AppID:           appID,
		GroupHash:       groupHash,
		ExceptionType:   req.ExceptionType,
		ExceptionReason: req.ExceptionReason,
		StackTrace:      req.StackTrace,
		AppVersion:      req.AppVersion,
		BuildNumber:     req.BuildNumber,
		OSVersion:       req.OSVersion,
		DeviceModel:     req.DeviceModel,
		DeviceIDHash:    req.DeviceIDHash,
		Locale:          req.Locale,
		FreeMemoryMB:    req.FreeMemoryMB,
		FreeDiskMB:      req.FreeDiskMB,
		BatteryLevel:    req.BatteryLevel,
		IsCharging:      req.IsCharging,
		OccurredAt:      req.OccurredAt,
	}

	// Store crash
	if err := h.crashStore.Create(r.Context(), crash); err != nil {
		http.Error(w, `{"error":"failed to store crash"}`, http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CrashResponse{
		ID:        crash.ID,
		GroupHash: crash.GroupHash,
	})
}

// CrashRequest represents the incoming crash payload
type CrashRequest struct {
	ExceptionType   string              `json:"exception_type"`
	ExceptionReason string              `json:"exception_reason"`
	StackTrace      []models.StackFrame `json:"stack_trace"`
	AppVersion      string              `json:"app_version"`
	BuildNumber     string              `json:"build_number"`
	OSVersion       string              `json:"os_version"`
	DeviceModel     string              `json:"device_model"`
	DeviceIDHash    string              `json:"device_id_hash"`
	Locale          string              `json:"locale"`
	FreeMemoryMB    int                 `json:"free_memory_mb"`
	FreeDiskMB      int                 `json:"free_disk_mb"`
	BatteryLevel    float64             `json:"battery_level"`
	IsCharging      bool                `json:"is_charging"`
	OccurredAt      time.Time           `json:"occurred_at"`
}

// CrashResponse represents the API response
type CrashResponse struct {
	ID        string `json:"id"`
	GroupHash string `json:"group_hash"`
}

func generateULID() string {
	// Simple ULID-like generation
	// Use proper ULID library for production
	return fmt.Sprintf("%d%010d", time.Now().UnixMilli(), rand.Int63())
}
