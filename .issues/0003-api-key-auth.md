---
id: 3
title: Add API key authentication middleware
status: open
priority: critical
labels: [api, auth]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Secure all ingestion endpoints with API key validation.

## Requirements
- Extract `X-Backfeedr-Key` header
- Look up app in database
- Reject invalid keys with 401
- Set rate limit: 100 req/min per key
- Track last used timestamp

## Key Format
- `bf_live_...` - production
- `bf_test_...` - test

## Acceptance Criteria
- [ ] Middleware validates all /api/* routes
- [ ] Returns 401 for missing/invalid keys
- [ ] Rate limiting enforced
- [ ] Test vs live key distinction
