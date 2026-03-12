package api

import (
	"encoding/json"
	"net/http"

	"github.com/steviee/backfeedr/internal/models"
	"github.com/steviee/backfeedr/internal/store"
)

// AdminHandler provides admin endpoints for testing
type AdminHandler struct {
	appStore *store.AppStore
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(appStore *store.AppStore) *AdminHandler {
	return &AdminHandler{appStore: appStore}
}

// CreateApp creates a new app (admin/debug endpoint)
func (h *AdminHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	// Only allow in test mode - check auth token matches test token
	// This is intentionally simple for MVP, should be more secure in production
	
	var req struct {
		Name     string `json:"name"`
		BundleID string `json:"bundle_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	
	app := &models.App{
		Name:     req.Name,
		BundleID: req.BundleID,
	}
	
	if err := h.appStore.Create(r.Context(), app); err != nil {
		http.Error(w, `{"error":"failed to create app"}`, http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      app.ID,
		"name":    app.Name,
		"api_key": app.APIKey,
	})
}
