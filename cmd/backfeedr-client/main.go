// backfeedr-client - Reference test client for backfeedr API
//
// Usage:
//
//	backfeedr-client --endpoint https://crashes.example.com --api-key bf_live_... --command send-crash --file crash.json
//	backfeedr-client --endpoint https://crashes.example.com --api-key bf_live_... --command send-event --type session_start
//
// This client serves as:
//   - Living API documentation
//   - Reference implementation for SDKs
//   - CI integration test tool
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/steviee/backfeedr/pkg/client"
)

type Config struct {
	Endpoint  string
	APIKey    string
	Command   string
	File      string
	EventType string
	Debug     bool
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

	flag.StringVar(&cfg.Endpoint, "endpoint", "", "API endpoint URL")
	flag.StringVar(&cfg.APIKey, "api-key", "", "API key")
	flag.StringVar(&cfg.Command, "command", "", "Command: send-crash, send-event, batch-events, health, interactive")
	flag.StringVar(&cfg.File, "file", "", "JSON file to send")
	flag.StringVar(&cfg.EventType, "type", "custom", "Event type")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "backfeedr-client - Reference test client for backfeedr API\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --command interactive\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint URL --api-key KEY --command CMD\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  interactive    Interactive menu mode\n")
		fmt.Fprintf(os.Stderr, "  send-crash     Send a crash report\n")
		fmt.Fprintf(os.Stderr, "  send-event     Send a single event\n")
		fmt.Fprintf(os.Stderr, "  batch-events   Send multiple events\n")
		fmt.Fprintf(os.Stderr, "  health         Check API health\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  # Interactive mode\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --command interactive\n\n")
		fmt.Fprintf(os.Stderr, "  # Send a crash report\n")
		fmt.Fprintf(os.Stderr, "  backfeedr-client --endpoint https://crashes.example.com \\\n")
		fmt.Fprintf(os.Stderr, "    --api-key bf_live_... --command send-crash --file crash.json\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Interactive mode: no required flags
	if cfg.Command == "interactive" {
		return cfg
	}

	// For other commands, validate required flags
	if cfg.Command == "" {
		fmt.Fprintf(os.Stderr, "Error: --command is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Health doesn't require auth
	if cfg.Command != "health" && cfg.Command != "interactive" {
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
	// Interactive mode
	if cfg.Command == "interactive" {
		return runInteractive()
	}

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

// runInteractive starts the interactive menu mode
func runInteractive() error {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("\n🚀 backfeedr-client Interactive Mode")
		fmt.Println("====================================")
		fmt.Println("1. Send test crash")
		fmt.Println("2. Send test event")
		fmt.Println("3. Send batch events")
		fmt.Println("4. Check health")
		fmt.Println("5. Exit")
		fmt.Print("\nChoice: ")
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			interactiveSendCrash(reader)
		case "2":
			interactiveSendEvent(reader)
		case "3":
			interactiveSendBatch(reader)
		case "4":
			interactiveHealthCheck(reader)
		case "5":
			fmt.Println("👋 Goodbye!")
			return nil
		default:
			fmt.Println("❌ Invalid choice")
		}
	}
}

func interactiveSendCrash(reader *bufio.Reader) {
	fmt.Print("Enter endpoint URL [http://localhost:8080]: ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	
	fmt.Print("Enter API key [bf_live_test]: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		apiKey = "bf_live_test"
	}
	
	c := client.New(endpoint, apiKey)
	crash := client.CrashReport{
		ExceptionType: "EXC_BAD_ACCESS",
		ExceptionReason: "Test from interactive mode",
		AppVersion: "1.0.0",
		OSVersion: "18.0.0",
		DeviceModel: "iPhone16,1",
		OccurredAt: time.Now().UTC(),
	}
	
	resp, err := c.SendCrash(&crash)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Crash sent! ID: %s, Group: %s\n", resp.ID, resp.GroupHash)
}

func interactiveSendEvent(reader *bufio.Reader) {
	fmt.Print("Enter endpoint URL [http://localhost:8080]: ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	
	fmt.Print("Enter API key [bf_live_test]: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		apiKey = "bf_live_test"
	}
	
	fmt.Print("Enter event name [test]: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "test"
	}
	
	c := client.New(endpoint, apiKey)
	event := client.EventRequest{
		Type: "custom",
		Name: name,
		AppVersion: "1.0.0",
		OSVersion: "18.0.0",
		DeviceModel: "iPhone16,1",
		OccurredAt: time.Now().UTC(),
	}
	
	resp, err := c.SendEvent(&event)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Event sent! ID: %s\n", resp.ID)
}

func interactiveSendBatch(reader *bufio.Reader) {
	fmt.Print("Enter endpoint URL [http://localhost:8080]: ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	
	fmt.Print("Enter API key [bf_live_test]: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		apiKey = "bf_live_test"
	}
	
	c := client.New(endpoint, apiKey)
	events := []client.EventRequest{
		{Type: "session_start", AppVersion: "1.0.0", OSVersion: "18.0.0", DeviceModel: "iPhone16,1", OccurredAt: time.Now().UTC()},
		{Type: "custom", Name: "action", AppVersion: "1.0.0", OSVersion: "18.0.0", DeviceModel: "iPhone16,1", OccurredAt: time.Now().UTC()},
		{Type: "session_end", AppVersion: "1.0.0", OSVersion: "18.0.0", DeviceModel: "iPhone16,1", OccurredAt: time.Now().UTC()},
	}
	
	resp, err := c.SendBatch(events)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Batch sent! Count: %d, IDs: %v\n", resp.Count, resp.IDs)
}

func interactiveHealthCheck(reader *bufio.Reader) {
	fmt.Print("Enter endpoint URL [http://localhost:8080]: ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	
	c := client.New(endpoint, "")
	healthy, err := c.HealthCheck(endpoint)
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		return
	}
	
	if healthy {
		fmt.Println("✅ Server is healthy")
	} else {
		fmt.Println("❌ Server not responding")
	}
}
