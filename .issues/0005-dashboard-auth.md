---
id: 5
title: Implement dashboard authentication
status: open
priority: high
labels: [auth, dashboard]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Add token-based authentication for web dashboard.

## Requirements
- Read `BACKFEEDR_AUTH_TOKEN` env var
- Cookie-based session after login
- Login form at `/login`
- Protect all dashboard routes
- Token generated on first init if not set

## Acceptance Criteria
- [ ] Login form works
- [ ] Cookie session persists
- [ ] Protected routes redirect to login
- [ ] Token regenerated on request (optional)
