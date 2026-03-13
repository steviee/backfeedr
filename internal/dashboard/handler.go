package dashboard

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/steviee/backfeedr/internal/models"
	"github.com/steviee/backfeedr/internal/store"
)

//go:embed templates/*.html
var templateFS embed.FS

// Handler handles dashboard requests
type Handler struct {
	crashStore   *store.CrashStore
	appStore     *store.AppStore
	metricsStore *store.MetricsStore
	templates    *template.Template
}

// NewHandler creates a new dashboard handler
func NewHandler(crashStore *store.CrashStore, appStore *store.AppStore, metricsStore *store.MetricsStore) (*Handler, error) {
	// Parse templates from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &Handler{
		crashStore:   crashStore,
		appStore:     appStore,
		metricsStore: metricsStore,
		templates:    tmpl,
	}, nil
}

// Index handles the dashboard home page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Dashboard",
		"AppName": "backfeedr",
	}

	// Load real data
	ctx := r.Context()

	// Get total crashes
	crashes, _ := h.crashStore.List(ctx, "", 100)
	crashCount := len(crashes)

	// Get apps
	apps, _ := h.appStore.List(ctx)
	appCount := len(apps)

	// Get crash groups
	groups, _ := h.crashStore.GetGroups(ctx, "")
	groupCount := len(groups)

	data["Stats"] = []StatCard{
		{
			Label: "Total Crashes",
			Value: fmt.Sprintf("%d", crashCount),
			Class: "",
		},
		{
			Label: "Connected Apps",
			Value: fmt.Sprintf("%d", appCount),
			Class: "",
		},
		{
			Label: "Crash Groups",
			Value: fmt.Sprintf("%d", groupCount),
			Class: "",
		},
	}
	data["Crashes"] = crashes
	data["Groups"] = groups

	if err := h.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StatCard represents a single stat card
type StatCard struct {
	Label       string
	Value       string
	Class       string
	Change      string
	ChangeClass string
	ChartData   []int
}

// CrashList handles the crash list page
func (h *Handler) CrashList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get crashes with groups
	crashes, _ := h.crashStore.List(ctx, "", 100)
	groups, _ := h.crashStore.GetGroups(ctx, "")

	data := map[string]interface{}{
		"Title":   "Crashes",
		"Crashes": crashes,
		"Groups":  groups,
	}

	if err := h.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CrashDetail handles individual crash pages
func (h *Handler) CrashDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get crash ID from URL
	crashID := r.URL.Path[len("/crashes/"):]
	if crashID == "" {
		http.Error(w, "Missing crash ID", http.StatusBadRequest)
		return
	}

	crash, err := h.crashStore.GetByID(ctx, crashID)
	if err != nil || crash == nil {
		http.Error(w, "Crash not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Crash Detail",
		"Crash": crash,
	}

	if err := h.templates.ExecuteTemplate(w, "crash_detail", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AppList handles the apps list page
func (h *Handler) AppList(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Apps",
	}

	ctx := r.Context()
	apps, err := h.appStore.List(ctx)
	if err != nil {
		// Render with empty list
		apps = []*models.App{}
	}
	data["Apps"] = apps

	if err := h.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Settings handles the settings page
func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Settings",
	}

	if err := h.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// APIOverview returns overview data for HTMX updates
func (h *Handler) APIOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Return updated stats HTML
	fmt.Fprintf(w, `
		<div class="stats-grid" hx-swap-oob="true">
			<article class="stat-card">
				<div class="stat-label">Crash-Free Rate (7d)</div>
				<div class="stat-value">%.1f%%</div>
				<div class="chart-container">
					<canvas id="crash-free-chart"></canvas>
				</div>
			</article>
			<article class="stat-card">
				<div class="stat-label">Active Users (DAU)</div>
				<div class="stat-value">%d</div>
			</article>
			<article class="stat-card">
				<div class="stat-label">Sessions (7d)</div>
				<div class="stat-value">%d</div>
			</article>
			<article class="stat-card">
				<div class="stat-label">Crashes (7d)</div>
				<div class="stat-value text-red-500">%d</div>
			</article>
		</div>
	`, 98.5, 1234, 12456, 23)
}
