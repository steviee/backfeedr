package server

import (
	_ "embed"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/steviee/backfeedr/internal/config"
	"github.com/steviee/backfeedr/internal/store"
)

// Server holds the HTTP server and dependencies
type Server struct {
	cfg    *config.Config
	db     *store.DB
	router chi.Router
	srv    *http.Server
}

// New creates a new server instance
func New(cfg *config.Config, db *store.DB) *Server {
	s := &Server{
		cfg: cfg,
		db:  db,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestSize(int64(cfg.MaxBodySize)))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.handleHealth)
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.apiKeyMiddleware)
			r.Post("/crashes", s.handleCrash)
			r.Post("/events", s.handleEvent)
			r.Post("/events/batch", s.handleBatch)
		})
	})

	// Dashboard routes
	r.Group(func(r chi.Router) {
		r.Use(s.authTokenMiddleware)
		r.Get("/", s.handleDashboard)
		r.Get("/crashes", s.handleCrashList)
		r.Get("/crashes/{id}", s.handleCrashDetail)
		r.Get("/apps", s.handleAppList)
		r.Get("/settings", s.handleSettings)
	})

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(webStatic))))

	s.router = r
	s.srv = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	return s
}

//go:embed all:web/static
var webStatic embed.FS

// Start begins listening for requests
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleCrash(w http.ResponseWriter, r *http.Request) {
	// TODO: implement crash ingestion
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: implement event ingestion
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleBatch(w http.ResponseWriter, r *http.Request) {
	// TODO: implement batch ingestion
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: implement dashboard
	fmt.Fprint(w, "backfeedr dashboard")
}

func (s *Server) handleCrashList(w http.ResponseWriter, r *http.Request) {
	// TODO: implement crash list
}

func (s *Server) handleCrashDetail(w http.ResponseWriter, r *http.Request) {
	// TODO: implement crash detail
}

func (s *Server) handleAppList(w http.ResponseWriter, r *http.Request) {
	// TODO: implement app list
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	// TODO: implement settings
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate API key + HMAC
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate auth token
		next.ServeHTTP(w, r)
	})
}
