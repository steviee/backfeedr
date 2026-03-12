---
id: 6
title: Create dashboard overview page
status: open
priority: high
labels: [dashboard]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Main dashboard showing key metrics at a glance.

## Requirements
- Crash-free rate (7-day sparkline)
- DAU/MAU counters
- Session counts (7d/30d)
- Top 5 crash groups
- HTMX-driven updates
- Pico CSS styling

## UI Elements
- Stats cards with sparklines
- Crash list with occurrence counts
- App selector dropdown
- Dark mode support

## Acceptance Criteria
- [ ] Overview loads with real data
- [ ] Sparklines render
- [ ] Responsive layout
- [ ] Auto-refresh (30s)
