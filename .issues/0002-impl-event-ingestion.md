---
id: 2
title: Implement event ingestion API
status: closed
priority: critical
labels: [api]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Implement event tracking endpoints for session metrics.

## Requirements
- `POST /api/v1/events` - single event
- `POST /api/v1/events/batch` - up to 100 events
- Event types: session_start, session_end, error, custom
- Store in events table
- Properties as JSON

## Acceptance Criteria
- [ ] Single event submission works
- [ ] Batch submission works (max 100)
- [ ] Validates event type enum
- [ ] Returns 201 on success
