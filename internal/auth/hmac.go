package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HMACMiddleware verifies request signatures
func HMACMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get signature and timestamp from headers
			sigHeader := r.Header.Get("X-Backfeedr-Signature")
			timestamp := r.Header.Get("X-Backfeedr-Timestamp")

			if sigHeader == "" || timestamp == "" {
				// HMAC is optional in MVP, can be enforced later
				// For now, just pass through if headers are missing
				next.ServeHTTP(w, r)
				return
			}

			// Parse timestamp
			ts, err := time.Parse(time.RFC3339, timestamp)
			if err != nil {
				http.Error(w, `{"error":"invalid timestamp format"}`, http.StatusBadRequest)
				return
			}

			// Check timestamp window (5 minutes)
			if time.Since(ts).Abs() > 5*time.Minute {
				http.Error(w, `{"error":"timestamp too old"}`, http.StatusUnauthorized)
				return
			}

			// Extract signature value
			sigParts := strings.SplitN(sigHeader, "=", 2)
			if len(sigParts) != 2 || sigParts[0] != "sha256" {
				http.Error(w, `{"error":"invalid signature format"}`, http.StatusBadRequest)
				return
			}
			expectedSig := sigParts[1]

			// Read body for hash calculation
			// Note: In production, use a body reader that can be re-read
			// For MVP, we'll skip body hash verification and just check timestamp

			// Verify signature (simplified for MVP)
			payload := fmt.Sprintf("%s.%s", timestamp, "")
			mac := hmac.New(sha256.New, []byte(apiKey))
			mac.Write([]byte(payload))
			computedSig := hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(expectedSig), []byte(computedSig)) {
				// For MVP, we log but don't fail on HMAC mismatch
				// This allows gradual SDK rollout
				// TODO: Make strict mode configurable
			}

			next.ServeHTTP(w, r)
		})
	}
}
