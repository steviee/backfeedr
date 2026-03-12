---
id: 10
title: Implement data retention policy
status: closed
priority: medium
labels: [db, infra]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Auto-delete old data based on retention setting.

## Requirements
- Config: `BACKFEEDR_RETENTION_DAYS` (default 90)
- Delete from crashes, events older than N days
- Keep daily_metrics (or aggregate further?)
- Run daily at 3 AM

## Acceptance Criteria
- [ ] Old data deleted
- [ ] Configurable retention
- [ ] Logs deletion count
- [ ] Keeps recent data intact
