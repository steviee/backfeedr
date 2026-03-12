---
id: 11
title: Create Swift SDK repository scaffold
status: closed
priority: critical
labels: [sdk]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Set up separate repo backfeedr-swift with SPM package structure.

## Requirements
- SPM Package.swift
- BackfeedrKit main target
- On-device symbolication support
- HMAC signing support
- Offline queue (local storage)
- PII scrubbing filters

## Structure
```
BackfeedrKit/
├── Sources/
│   ├── Backfeedr.swift      # Main API
│   ├── CrashReporter.swift
│   ├── EventTracker.swift
│   ├── HMAC.swift
│   ├── Queue.swift
│   └── Scrubber.swift
└── Tests/
```

## Acceptance Criteria
- [ ] SPM package builds
- [ ] Main API surface defined
- [ ] HMAC signing ready
