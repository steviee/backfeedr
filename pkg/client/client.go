// Package client provides a reference implementation of the backfeedr API client.
// This package serves as living documentation and can be used as a basis for SDKs.
package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main API client
type Client struct {
	Endpoint   string
	APIKey     string
	HTTPClient *http.Client
	Debug      bool
}

// New creates a new API client
func New(endpoint, apiKey string) *Client {
	return &Client{
		Endpoint:   endpoint,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// HealthCheck checks if the API is healthy
func (c *Client) HealthCheck(endpoint string) (bool, error) {
	url := endpoint + "/api/v1/health"

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// doRequest performs an authenticated request with HMAC signing
func (c *Client) doRequest(method, path string, body []byte) (*http.Response, error) {
	url := c.Endpoint + path

	// Create timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)

	// Calculate HMAC signature
	signature := c.signRequest(timestamp, body)

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Backfeedr-Key", c.APIKey)
	req.Header.Set("X-Backfeedr-Timestamp", timestamp)
	req.Header.Set("X-Backfeedr-Signature", "sha256="+signature)

	if c.Debug {
		fmt.Printf("[DEBUG] Request: %s %s\n", method, url)
		fmt.Printf("[DEBUG] Headers:\n")
		fmt.Printf("  X-Backfeedr-Key: %s\n", c.APIKey[:20]+"...")
		fmt.Printf("  X-Backfeedr-Timestamp: %s\n", timestamp)
		fmt.Printf("  X-Backfeedr-Signature: sha256=%s...\n", signature[:20])
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if c.Debug {
		fmt.Printf("[DEBUG] Response: %d %s\n", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// signRequest creates an HMAC signature for the request
func (c *Client) signRequest(timestamp string, body []byte) string {
	// Calculate body hash
	bodyHash := sha256.Sum256(body)
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	// Build payload: timestamp.body_hash
	payload := timestamp + "." + bodyHashHex

	// Calculate HMAC
	mac := hmac.New(sha256.New, []byte(c.APIKey))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// handleResponse processes the API response
func (c *Client) handleResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if c.Debug {
		fmt.Printf("[DEBUG] Response body: %s\n", string(body))
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	if v != nil && len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
	}

	return nil
}
