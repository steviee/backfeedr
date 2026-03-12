// backfeedr-client - Reference test client for backfeedr API
//
// Usage:
//   backfeedr-client --endpoint https://crashes.example.com --api-key bf_live_... --command send-crash --file crash.json
//   backfeedr-client --endpoint https://crashes.example.com --api-key bf_live_... --command send-event --type session_start
//
// This client serves as:
//   - Living API documentation
//   - Reference implementation for SDKs
//   - CI integration test tool
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/steviee/backfeedr/pkg/client"
)

type Config struct {
	Endpoint string
	APIKey   string
	Command  string
	File     string
	EventType string
	Debug    bool
}

func main() {
	cfg := parseFlags()

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Endpoint, "endpoint", "", "API endpoint URL (required)")
	flag.StringVar(&cfg.APIKey, "api-key", "", "API key (required)")
	flag.StringVar(&cfg.Command, "command", "", "Command: send-crash, send-event, batch-events, health")
	flag.StringVar(&cfg.File, "file", "", "JSON file to send (for send-crash, send-event, batch-events)")
	flag.StringVar(&cfg.EventType, "type", "custom", "Event type (for send-event)")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "backfeedr-client - Reference test client for backfeedr API\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint URL --api-key KEY --command CMD [options]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  send-crash     Send a crash report\n")
		fmt.Fprintf(os.Stderr, "  send-event     Send a single event\n")
		fmt.Fprintf(os.Stderr, "  batch-events   Send multiple events\n")
		fmt.Fprintf(os.Stderr, "  health         Check API health\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  # Send a crash report\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint https://crashes.example.com \\\n")
		fmt.Fprintf(os.Stderr, "    --api-key bf_live_... --command send-crash --file crash.json\n\n")
		fmt.Fprintf(os.Stderr, "  # Send a session start event\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint https://crashes.example.com \\\n")
		fmt.Fprintf(os.Stderr, "    --api-key bf_live_... --command send-event --type session_start\n\n")
		fmt.Fprintf(os.Stderr, "  # Check health\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint https://crashes.example.com --command health\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Validate required flags
	if cfg.Command == "" {
		fmt.Fprintf(os.Stderr, "Error: --command is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Health doesn't require auth
	if cfg.Command != "health" {
		if cfg.Endpoint == "" {
			fmt.Fprintf(os.Stderr, "Error: --endpoint is required\n\n")
			flag.Usage()
			os.Exit(1)
		}
		if cfg.APIKey == "" {
			fmt.Fprintf(os.Stderr, "Error: --api-key is required\n\n")
			flag.Usage()
			os.Exit(1)
		}
	}

	return cfg
}

func run(cfg *Config) error {
	c := client.New(cfg.Endpoint, cfg.APIKey)
	c.Debug = cfg.Debug

	switch cfg.Command {
	case "send-crash":
		return sendCrash(c, cfg)
	case "send-event":
		return sendEvent(c, cfg)
	case "batch-events":
		return batchEvents(c, cfg)
	case "health":
		return checkHealth(c, cfg)
	default:
		return fmt.Errorf("unknown command: %s", cfg.Command)
	}
}

func sendCrash(c *client.Client, cfg *Config) error {
	var crash client.CrashReport

	if cfg.File != "" {
		// Load from file
		data, err := os.ReadFile(cfg.File)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		if err := json.Unmarshal(data, &crash); err != nil {
			return fmt.Errorf("parse JSON: %w", err)
		}
	} else {
		// Use example crash
		crash = exampleCrash()
	}

	resp, err := c.SendCrash(&crash)
	if err != nil {
		return fmt.Errorf("send crash: %w", err)
	}

	fmt.Printf("Crash sent successfully!\n")
	fmt.Printf("  ID:        %s\n", resp.ID)
	fmt.Printf("  Group Hash: %s\n", resp.GroupHash)

	return nil
}

func sendEvent(c *client.Client, cfg *Config) error {
	event := client.EventRequest{
		Type:         cfg.EventType,
		Name:         "test_event",
		OccurredAt:   time.Now().UTC(),
		AppVersion:   "1.0.0",
		OSVersion:    "18.3.1",
		DeviceModel:  "iPhone16,1",
		DeviceIDHash: "abc123...",
		Locale:       "de_DE",
	}

	if cfg.File != "" {
		// Load from file
		data, err := os.ReadFile(cfg.File)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("parse JSON: %w", err)
		}
	}

	resp, err := c.SendEvent(&event)
	if err != nil {
		return fmt.Errorf("send event: %w", err)
	}

	fmt.Printf("Event sent successfully!\n")
	fmt.Printf("  ID: %s\n", resp.ID)

	return nil
}

func batchEvents(c *client.Client, cfg *Config) error {
	events := []client.EventRequest{
		{Type: "session_start", OccurredAt: time.Now().UTC(), AppVersion: "1.0.0"},
		{Type: "custom", Name: "test", OccurredAt: time.Now().UTC(), AppVersion: "1.0.0"},
		{Type: "session_end", OccurredAt: time.Now().UTC(), AppVersion: "1.0.0"},
	}

	if cfg.File != "" {
		// Load from file
		data, err := os.ReadFile(cfg.File)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		if err := json.Unmarshal(data, &events); err != nil {
			return fmt.Errorf("parse JSON: %w", err)
		}
	}

	resp, err := c.SendBatch(events)
	if err != nil {
		return fmt.Errorf("send batch: %w", err)
	}

	fmt.Printf("Batch sent successfully!\n")
	fmt.Printf("  Count: %d\n", resp.Count)
	fmt.Printf("  IDs:   %v\n", resp.IDs)

	return nil
}

func checkHealth(c *client.Client, cfg *Config) error {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}

	healthy, err := c.HealthCheck(endpoint)
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}

	if healthy {
		fmt.Println("✓ API is healthy")
	} else {
		fmt.Println("✗ API is not responding")
		os.Exit(1)
	}

	return nil
}

func exampleCrash() client.CrashReport {
	return client.CrashReport{
		ExceptionType:   "EXC_BAD_ACCESS",
		ExceptionReason: "Attempted to dereference null pointer",
		StackTrace: []client.StackFrame{
			{Frame: 0, Symbol: "ContentView.body.getter", File: strPtr("ContentView.swift"), Line: intPtr(42)},
			{Frame: 1, Symbol: "SwiftUI.View.update"},
		},
		AppVersion:   "1.2.0",
		BuildNumber:  "47",
		OSVersion:    "18.3.1",
		DeviceModel:  "iPhone16,1",
		DeviceIDHash: "a9f3c721...",
		Locale:       "de_DE",
		FreeMemoryMB: 312,
		BatteryLevel: 0.67,
		OccurredAt:   time.Now().UTC(),
	}
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
