// Package client provides a reference implementation for the backfeedr API.
//
// This package serves three purposes:
//   1. Living documentation - Shows how to use the API
//   2. Reference implementation - Basis for SDKs in other languages
//   3. Testing tool - Can be used in CI for integration tests
//
// Basic usage:
//
//	import "github.com/steviee/backfeedr/pkg/client"
//
//	c := client.New("https://crashes.example.com", "bf_live_...")
//
//	// Send a crash report
//	crash := &client.CrashReport{
//	    ExceptionType: "EXC_BAD_ACCESS",
//	    ...
//	}
//	resp, err := c.SendCrash(crash)
//
package client
