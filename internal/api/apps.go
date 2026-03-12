package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/steviee/backfeedr/internal/models"
	"github.com/steviee/backfeedr/internal/store"
)

// AppHandler handles app management API
type AppHandler struct {
	appStore *store.AppStore
}

// NewAppHandler creates a new app handler
func NewAppHandler(appStore *store.AppStore) *AppHandler {
	return &AppHandler{appStore: appStore}
}

// List returns all apps
func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appStore.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list apps"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

// Create creates a new app
func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		BundleID string `json:"bundle_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.BundleID == "" {
		http.Error(w, `{"error":"name and bundle_id are required"}`, http.StatusBadRequest)
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
		"id":       app.ID,
		"name":     app.Name,
		"bundle_id": app.BundleID,
		"api_key":  app.APIKey,
		"created_at": app.CreatedAt,
	})
}

// Get returns a single app
func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	app, err := h.appStore.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, `{"error":"app not found"}`, http.StatusNotFound)
		return
	}

	// Don't return the API key in normal GET
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         app.ID,
		"name":       app.Name,
		"bundle_id":  app.BundleID,
		"created_at": app.CreatedAt,
	})
}

// RotateKey rotates the API key
func (h *AppHandler) RotateKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newKey, err := h.appStore.RotateAPIKey(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"failed to rotate key"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"api_key": newKey,
	})
}

// Delete deletes an app
func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id") // TODO: Use this for deletion

	// TODO: Implement soft delete
	w.WriteHeader(http.StatusNotImplemented)
}
