# backfeedr API Reference

## Base URL

```
https://your-server.com/api/v1
```

## Authentication

Most endpoints require an API key in the `X-Backfeedr-Key` header:

```
X-Backfeedr-Key: bf_live_abc123...
```

## Endpoints

### Health Check

```
GET /api/v1/health
```

**No authentication required.**

#### Response

```json
{
  "status": "ok"
}
```

---

### Submit Crash

```
POST /api/v1/crashes
```

Submit a crash report from your iOS app.

#### Headers
- `Content-Type: application/json`
- `X-Backfeedr-Key: bf_live_...` (required)

#### Request Body

```json
{
  "exception_type": "EXC_BAD_ACCESS",
  "exception_reason": "Attempted to dereference null pointer",
  "stack_trace": [
    {
      "frame": 0,
      "symbol": "ContentView.body.getter",
      "file": "ContentView.swift",
      "line": 42
    }
  ],
  "app_version": "1.2.0",
  "build_number": "47",
  "os_version": "18.3.1",
  "device_model": "iPhone16,1",
  "device_id_hash": "abc123...",
  "locale": "de_DE",
  "free_memory_mb": 312,
  "battery_level": 0.67,
  "is_charging": false,
  "occurred_at": "2026-03-12T14:22:31Z"
}
```

#### Response

```json
{
  "id": "abc123...",
  "group_hash": "9eaa6e44..."
}
```

---

### Submit Event

```
POST /api/v1/events
```

Submit a single event.

#### Headers
- `Content-Type: application/json`
- `X-Backfeedr-Key: bf_live_...` (required)

#### Request Body

```json
{
  "type": "custom",
  "name": "button_click",
  "properties": {
    "button": "submit",
    "screen": "checkout"
  },
  "app_version": "1.2.0",
  "os_version": "18.3.1",
  "device_model": "iPhone16,1",
  "device_id_hash": "abc123...",
  "session_id": "sess_xyz789",
  "locale": "de_DE",
  "occurred_at": "2026-03-12T14:22:31Z"
}
```

**Event Types:** `session_start`, `session_end`, `error`, `custom`

#### Response

```json
{
  "id": "evt_abc123..."
}
```

---

### Submit Batch Events

```
POST /api/v1/events/batch
```

Submit up to 100 events at once.

#### Headers
- `Content-Type: application/json`
- `X-Backfeedr-Key: bf_live_...` (required)

#### Request Body

```json
{
  "events": [
    {
      "type": "session_start",
      "session_id": "sess_123",
      "occurred_at": "2026-03-12T14:20:00Z"
    },
    {
      "type": "custom",
      "name": "page_view",
      "occurred_at": "2026-03-12T14:21:00Z"
    }
  ]
}
```

#### Response

```json
{
  "count": 2,
  "ids": ["evt_1...", "evt_2..."]
}
```

---

### Metrics

These endpoints are public (no authentication required).

#### Daily Crashes

```
GET /api/v1/metrics/daily-crashes
```

#### Response

```json
{
  "labels": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
  "data": [5, 2, 8, 3, 1, 0, 2]
}
```

#### Crash Types

```
GET /api/v1/metrics/crash-types
```

#### Response

```json
{
  "labels": ["EXC_BAD_ACCESS", "SIGABRT"],
  "data": [45, 12]
}
```

#### Device Distribution

```
GET /api/v1/metrics/devices
```

#### Response

```json
{
  "labels": ["iPhone16,1", "iPhone15,2"],
  "data": [30, 15]
}
```

---

### Dashboard Overview

```
GET /api/v1/overview
```

**No authentication required.**

#### Response

```json
{
  "crash_free_rate_7d": 98.5,
  "dau": 1234,
  "sessions_7d": 12456,
  "crashes_7d": 23,
  "top_crashes": [
    {
      "group_hash": "9eaa6e44...",
      "exception_type": "EXC_BAD_ACCESS",
      "exception_reason": "Null pointer...",
      "count": 15
    }
  ]
}
```

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "invalid JSON"
}
```

### 401 Unauthorized

```json
{
  "error": "missing API key"
}
```

### 429 Rate Limited

```json
{
  "error": "rate limit exceeded"
}
```

---

## cURL Examples

### Send Crash

```bash
curl -X POST http://localhost:8080/api/v1/crashes \
  -H "Content-Type: application/json" \
  -H "X-Backfeedr-Key: bf_live_your_key" \
  -d @examples/crash.json
```

### Send Event

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "X-Backfeedr-Key: bf_live_your_key" \
  -d '{
    "type": "custom",
    "name": "purchase",
    "occurred_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
  }'
```

### Check Health

```bash
curl http://localhost:8080/api/v1/health
```

---

## Rate Limits

- **100 requests per minute** per API key
- Applies to: crashes, events, batch

Exceeding returns `429 Too Many Requests`.

---

## HMAC Signing (Optional)

For additional security, sign requests with HMAC-SHA256:

1. Add `X-Backfeedr-Timestamp` (ISO8601)
2. Calculate `payload = timestamp.body_hash`
3. Sign with `HMAC-SHA256(api_key, payload)`
4. Add `X-Backfeedr-Signature: sha256=...`

See reference client for implementation.
