package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/steviee/backfeedr/internal/store"
)

// APIKeyMiddleware validates API keys for ingestion endpoints
func APIKeyMiddleware(appStore *store.AppStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-Backfeedr-Key")
			if apiKey == "" {
				http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
				return
			}

			// Validate key format
			if !strings.HasPrefix(apiKey, "bf_live_") && !strings.HasPrefix(apiKey, "bf_test_") {
				http.Error(w, `{"error":"invalid API key format"}`, http.StatusUnauthorized)
				return
			}

			// Look up app
			app, err := appStore.GetByAPIKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
				return
			}
			if app == nil {
				http.Error(w, `{"error":"invalid API key"}`, http.StatusUnauthorized)
				return
			}

			// Add app ID to context
			ctx := context.WithValue(r.Context(), "app_id", app.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
