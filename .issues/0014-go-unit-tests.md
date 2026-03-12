---
id: 14
title: Write Go server unit tests
status: open
priority: medium
labels: [api]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Test coverage for Go backend.

## Test Coverage
- Config parsing
- Database operations (in-memory SQLite)
- API handlers (mock store)
- Auth middleware
- HMAC verification
- Crash grouping logic

## Tools
- `testing` package
- `testify/assert` (optional)
- `httptest` for handlers

## Acceptance Criteria
- [ ] Store layer tested
- [ ] API handlers tested
- [ ] Auth middleware tested
- [ ] CI runs tests
