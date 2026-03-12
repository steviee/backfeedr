---
id: 13
title: Optimize Docker image build
status: open
priority: medium
labels: [infra]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Production-ready Docker build under 30MB.

## Requirements
- Multi-stage build
- Alpine base image
- Static binary (CGO_ENABLED=0, `-extldflags '-static'`)
- Healthcheck endpoint
- Non-root user

## Current
Builder stage exists but runtime is basic.

## Acceptance Criteria
- [ ] Image size <30MB
- [ ] Runs as non-root
- [ ] Healthcheck passes
- [ ] Single binary, no shell needed
