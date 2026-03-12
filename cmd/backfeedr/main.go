package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/steviee/backfeedr/internal/config"
	"github.com/steviee/backfeedr/internal/server"
	"github.com/steviee/backfeedr/internal/store"
	"github.com/steviee/backfeedr/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := store.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Start background workers
	metricsWorker := worker.NewMetricsWorker(db)
	metricsWorker.Start()
	defer metricsWorker.Stop()

	retentionWorker := worker.NewRetentionWorker(db, cfg)
	retentionWorker.Start()
	defer retentionWorker.Stop()

	srv := server.New(cfg, db)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
