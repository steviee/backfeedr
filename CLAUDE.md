# backfeedr — Implementation Brief for Claude Code

## Project Goal

Build a self-hosted crash reporting and app metrics system for iOS indie developers. 
Single Docker container, SQLite database, Go + HTMX dashboard, privacy-first.

---

## Tech Stack

- **Backend:** Go 1.22+
- **Database:** SQLite (WAL mode) via modernc.org/sqlite
- **HTTP Router:** chi (lightweight, idiomatic)
- **Dashboard:** Go html/template + HTMX + Alpine.js + Pico CSS
- **iOS SDK:** Swift 6, SPM, separate repo (backfeedr/backfeedr-swift)
- **Container:** Docker Alpine (~20MB image)

---

## Repository Structure

```
backfeedr/
├── cmd/backfeedr/
│   └── main.go              # Entry point
├── internal/
│   ├── server/
│   │   ├── server.go        # HTTP server setup
│   │   ├── routes.go        # Route registration
│   │   └── middleware.go    # Auth, rate-limit, logging
│   ├── api/
│   │   ├── crashes.go       # POST /api/v1/crashes
│   │   ├── events.go        # POST /api/v1/events
│   │   ├── batch.go         # POST /api/v1/events/batch
│   │   └── health.go        # GET /api/v1/health
│   ├── dashboard/
│   │   ├── handler.go       # HTMX page handlers
│   │   ├── metrics.go       # Metrics aggregation handlers
│   │   ├── templates/       # Go templates (embed.FS)
│   │   └── static/          # CSS, JS (Pico, HTMX, Alpine)
│   ├── store/
│   │   ├── sqlite.go        # DB connection + migrations
│   │   ├── crashes.go       # Crash CRUD + grouping
│   │   ├── events.go        # Event CRUD
│   │   ├── apps.go          # App management
│   │   └── metrics.go       # Aggregated metrics queries
│   ├── config/
│   │   └── config.go        # Env/config parsing
│   └── auth/
│       ├── token.go         # API key validation
│       └── hmac.go          # HMAC signature verification
├── web/
│   ├── templates/           # HTML templates
│   └── static/              # Static assets
├── migrations/              # SQL migration files
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── CLAUDE.md
├── README.md
└── LICENSE
```

---

## Data Model (SQLite)

### apps
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT (ULID) | Primary key |
| name | TEXT | App name |
| bundle_id | TEXT | Bundle identifier |
| api_key | TEXT | Ingestion API key |
| created_at | DATETIME | Creation timestamp |

### crashes
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT (ULID) | Primary key |
| app_id | TEXT | FK → apps |
| group_hash | TEXT | Hash for crash grouping |
| exception_type | TEXT | e.g. EXC_BAD_ACCESS |
| exception_reason | TEXT | Crash description |
| stack_trace | TEXT (JSON) | Symbolized stack trace |
| app_version | TEXT | CFBundleShortVersionString |
| build_number | TEXT | CFBundleVersion |
| os_version | TEXT | iOS version |
| device_model | TEXT | e.g. iPhone15,2 |
| locale | TEXT | de_DE |
| free_memory_mb | INTEGER | Available RAM |
| free_disk_mb | INTEGER | Available storage |
| battery_level | REAL | 0.0–1.0 |
| is_charging | BOOLEAN | Charging state |
| occurred_at | DATETIME | Crash timestamp |
| received_at | DATETIME | Server received |

### events
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT (ULID) | Primary key |
| app_id | TEXT | FK → apps |
| type | TEXT | session_start, session_end, error, custom |
| name | TEXT | Event name |
| properties | TEXT (JSON) | Key-value pairs |
| app_version | TEXT | App version |
| os_version | TEXT | iOS version |
| device_model | TEXT | Device model |
| session_id | TEXT | Session identifier |
| occurred_at | DATETIME | Timestamp |

### daily_metrics (materialized)
| Column | Type | Description |
|--------|------|-------------|
| app_id | TEXT | FK → apps |
| date | DATE | Day |
| sessions | INTEGER | Session count |
| unique_devices | INTEGER | DAU |
| crashes | INTEGER | Crash count |
| errors | INTEGER | Non-fatal errors |
| avg_session_sec | REAL | Avg session duration |

---

## API Design

All ingestion endpoints require API key in `X-Backfeedr-Key` header.
Responses are minimal JSON.

### Authentication
```
X-Backfeedr-Key: bf_live_a1b2c3d4e5f6...
X-Backfeedr-Timestamp: 2026-03-12T14:22:31Z
X-Backfeedr-Signature: sha256=a9f3c721...
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/crashes | Submit crash report |
| POST | /api/v1/events | Submit event(s) |
| POST | /api/v1/events/batch | Bulk events (max 100) |
| GET | /api/v1/health | Health check (no auth) |

### Crash Payload
```json
{
  "exception_type": "EXC_BAD_ACCESS",
  "exception_reason": "Attempted to dereference null pointer",
  "stack_trace": [
    { "frame": 0, "symbol": "ContentView.body.getter",
      "file": "ContentView.swift", "line": 42 }
  ],
  "app_version": "1.2.0",
  "build_number": "47",
  "os_version": "18.3.1",
  "device_model": "iPhone16,1",
  "device_id_hash": "a9f3...c721",
  "locale": "de_DE",
  "free_memory_mb": 312,
  "battery_level": 0.67,
  "occurred_at": "2026-03-12T14:22:31Z"
}
```

---

## Dashboard Views

| View | Content | Priority |
|------|---------|----------|
| Overview | Crash-free rate, DAU, sessions (sparklines), top crashes | MVP |
| Crashes | Grouped crash list, occurrences, affected versions | MVP |
| Crash Detail | Full stack trace, device distribution, breadcrumbs | MVP |
| Events | Event stream with filters | MVP |
| Metrics | DAU/MAU, retention (D1/D7/D30), version adoption | v1.1 |
| Apps | App management, API key generation | MVP |
| Settings | Retention policy, auth token, export | MVP |

---

## Security Layers

1. **PII Scrubbing** — No personal data collected
2. **TLS 1.3 Only** — No HTTP fallback
3. **HMAC Signing** — Request signature verification
4. **Rate Limiting** — 100 req/min per API key
5. **Timestamp Window** — 5 min tolerance (prevents replay)

---

## Implementation Order

1. `internal/config/config.go` — Config parsing
2. `internal/store/sqlite.go` — DB connection + migrations
3. `internal/store/crashes.go` — Crash storage
4. `internal/store/events.go` — Event storage
5. `internal/store/apps.go` — App management
6. `internal/auth/token.go` — API key validation
7. `internal/auth/hmac.go` — HMAC verification
8. `internal/server/server.go` — HTTP server setup
9. `internal/api/crashes.go` — Crash ingestion endpoint
10. `internal/api/events.go` — Event ingestion endpoint
11. `internal/dashboard/handler.go` — Dashboard handlers
12. `web/templates/` — HTML templates
13. `web/static/` — Static assets
14. `Dockerfile` + `docker-compose.yml`
15. Tests
16. README

---

## Build & Run

```bash
# Build
go build -o backfeedr ./cmd/backfeedr

# Run
./backfeedr

# Docker
docker build -t backfeedr .
docker run -p 8080:8080 -v ./data:/data backfeedr
```

---

## Out of Scope (MVP)

- Public key pinning (v1.0)
- DNS-TXT validation (v1.0)
- Apple App Attest (v1.2)
- SQLCipher encryption (v1.1)
- Alerting/webhooks (v1.1)
- Grafana integration
- Multi-user support
