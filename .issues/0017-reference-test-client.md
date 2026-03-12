---
id: 17
title: Build reference test client in Go
status: in-progress
priority: high
labels: [api, testing, sdk]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Create a simple Go-based reference client that demonstrates the API usage and serves as a test client for CI.

## Goals
- Validate API contract
- Serve as living documentation
- Base for SDK design decisions
- CI integration tests

## Requirements

### Client Features
- [ ] Configure endpoint + API key
- [ ] Send crash report with full payload
- [ ] Send single event
- [ ] Send batch events
- [ ] Handle HMAC signing
- [ ] Validate responses

### Structure
```
cmd/backfeedr-client/
├── main.go          # CLI tool
└── client/
    ├── client.go    # HTTP client
    ├── crash.go     # Crash reporting
    ├── event.go     # Event tracking
    └── auth.go      # HMAC signing
```

### CLI Interface
```bash
backfeedr-client \
  --endpoint https://crashes.example.com \
  --api-key bf_live_... \
  --command send-crash \
  --file crash.json

backfeedr-client \
  --endpoint https://crashes.example.com \
  --api-key bf_live_... \
  --command send-event \
  --type session_start
```

### Acceptance Criteria
- [ ] Client builds successfully
- [ ] Can send crashes and events
- [ ] HMAC signing implemented
- [ ] Used in CI pipeline
- [ ] Documented as reference implementation

## Notes
- Keep simple — no complex abstractions
- Code should be easy to port to Swift
- Include example JSON payloads
- Error handling should demonstrate API behavior
