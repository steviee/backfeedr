---
id: 9
title: Daily metrics aggregation job
status: closed
priority: medium
labels: [db]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Aggregate events into daily_metrics table each day.

## Requirements
- Count sessions per app/day
- Count unique devices (hash)
- Count crashes
- Calculate avg session duration
- Run via cron or goroutine

## Table: daily_metrics
- app_id, date (PK)
- sessions, unique_devices
- crashes, errors
- avg_session_sec

## Acceptance Criteria
- [ ] Aggregation query works
- [ ] Scheduled execution
- [ ] Backfill capability
- [ ] Idempotent (rerun safe)
