package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/steviee/backfeedr/internal/api"
	"github.com/steviee/backfeedr/internal/auth"
	"github.com/steviee/backfeedr/internal/config"
	"github.com/steviee/backfeedr/internal/dashboard"
	"github.com/steviee/backfeedr/internal/store"
)

// Static files are served from ./web/static directory
// TODO: Use embed when static files are ready

// Server holds the HTTP server and dependencies
type Server struct {
	cfg        *config.Config
	db         *store.DB
	crashStore *store.CrashStore
	appStore   *store.AppStore
	eventStore *store.EventStore
	router     chi.Router
	srv        *http.Server
}

// New creates a new server instance
func New(cfg *config.Config, db *store.DB) *Server {
	// Initialize stores
	crashStore := store.NewCrashStore(db)
	appStore := store.NewAppStore(db)
	eventStore := store.NewEventStore(db)

	// Create handlers
	crashHandler := api.NewCrashHandler(crashStore, appStore)
	eventHandler := api.NewEventHandler(eventStore, appStore)
	dashAPIHandler := api.NewDashboardHandler(crashStore, appStore)
	metricsHandler := api.NewMetricsHandler(crashStore, appStore)
	dashHandler, err := dashboard.NewHandler(crashStore, appStore, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create dashboard handler: %v", err))
	}

	s := &Server{
		cfg:        cfg,
		db:         db,
		crashStore: crashStore,
		appStore:   appStore,
		eventStore: eventStore,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestSize(int64(cfg.MaxBodySize)))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.handleHealth)
		r.Get("/overview", dashAPIHandler.GetOverview)

		// Metrics endpoints (public)
		r.Get("/metrics/daily-crashes", metricsHandler.GetDailyCrashes)
		r.Get("/metrics/crash-types", metricsHandler.GetCrashTypes)
		r.Get("/metrics/devices", metricsHandler.GetDeviceDistribution)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(appStore))
			r.Post("/crashes", crashHandler.HandleCrash)
			r.Post("/events", eventHandler.HandleEvent)
			r.Post("/events/batch", eventHandler.HandleBatch)
		})
	})

	// Dashboard routes
	r.Get("/", dashHandler.Index)
	r.Get("/crashes", dashHandler.CrashList)
	r.Get("/crashes/{id}", dashHandler.CrashDetail)
	r.Get("/apps", dashHandler.AppList)
	r.Get("/settings", dashHandler.Settings)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	s.router = r
	s.srv = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s
}

// Start begins listening for requests
func (s *Server) Start() error {
	addr := s.cfg.BindAddr + ":" + s.cfg.Port
	s.srv.Addr = addr

	// Print startup message with clickable URLs
	fmt.Printf("\n🚀 backfeedr server starting...\n")
	fmt.Printf("   Local:    http://localhost:%s\n", s.cfg.Port)
	fmt.Printf("   Network:  http://%s\n\n", addr)

	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
