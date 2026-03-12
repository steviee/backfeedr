---
id: 4
title: Implement HMAC request signing
status: open
priority: high
labels: [api, security]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Add HMAC-SHA256 signature verification for crash/event ingestion.

## Requirements
- Parse `X-Backfeedr-Timestamp` header (ISO8601)
- Parse `X-Backfeedr-Signature` header (`sha256=...`)
- Verify signature using API key as secret
- Reject requests older than 5 minutes (prevent replay)
- Fail on signature mismatch with 401

## Algorithm
```
payload = "{timestamp}.{body_hash}"
signature = HMAC-SHA256(key, payload)
```

## Acceptance Criteria
- [ ] Signature verification works
- [ ] Timestamp window enforced (5 min)
- [ ] Prevents replay attacks
- [ ] Graceful error messages
