---
id: 1
title: Implement crash ingestion API endpoint
status: open
priority: critical
labels: [api]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Implement `POST /api/v1/crashes` endpoint that receives crash reports from iOS SDK.

## Requirements
- Parse JSON payload
- Validate required fields (exception_type, occurred_at)
- Generate group_hash from top 3 app frames
- Validate API key from `X-Backfeedr-Key` header
- Store crash in SQLite crashes table
- Return 201 on success, 401/400 on error

## Payload Structure
```json
{
  "exception_type": "EXC_BAD_ACCESS",
  "exception_reason": "...",
  "stack_trace": [...],
  "app_version": "1.2.0",
  "build_number": "47",
  "os_version": "18.3.1",
  "device_model": "iPhone16,1",
  "device_id_hash": "...",
  "locale": "de_DE",
  "free_memory_mb": 312,
  "battery_level": 0.67,
  "occurred_at": "2026-03-12T14:22:31Z"
}
```

## Acceptance Criteria
- [ ] Endpoint accepts valid crash reports
- [ ] Group hash calculated correctly
- [ ] API key validation works
- [ ] Returns appropriate status codes
