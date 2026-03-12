---
id: 18
title: Privacy Dashboard for End-Users (DSGVO/GDPR)
status: open
priority: medium
labels: [feature, privacy, gdpr]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Allow end-users to view and manage their own data via a privacy-focused dashboard.

## User Story
As an app user, I want to see what data was collected from my device and be able to download or delete it, so that I have control over my personal information (GDPR/DSGVO compliance).

## Technical Approach

### 1. Device-Linked Dashboard URL
- URL format: `https://crashes.example.com/privacy/{device_id_hash}`
- Device ID is hashed (SHA-256) and included in URL
- No authentication required for URL access (URL itself is the "secret")

### 2. 2FA-Style Authorization
- When URL is accessed, backend sends push notification to app
- App shows: "Request from backfeedr.example.com - Allow access?"
- User must confirm in app within 60 seconds
- Only then is the dashboard data displayed

### 3. Dashboard Features
- View all crashes from this device
- View all events from this device  
- Download data as JSON/CSV
- Request deletion (immediate or within 30 days)
- See what data is retained and why

### 4. Security Considerations
- URL must be unguessable (long hash, min 16 chars)
- Rate limiting (max 10 tries per hour per device)
- Audit log of all access attempts
- Optional: Email verification before deletion

## Implementation Steps

### Phase 1: Basic Privacy Dashboard
- [ ] Privacy handler with device lookup
- [ ] Query crashes/events by device_id_hash
- [ ] HTML template for data display
- [ ] JSON/CSV export

### Phase 2: 2FA Authorization
- [ ] In-app notification/approval system
- [ ] Backend token generation for approved sessions
- [ ] Session management (30 min expiry)
- [ ] Deny/revoke mechanism

### Phase 3: Self-Service Deletion
- [ ] Data deletion request endpoint
- [ ] Anonymization option (keep aggregated stats, remove PII)
- [ ] Confirmation workflow
- [ ] Audit trail

## API Endpoints

```
GET /privacy/{device_hash}       → Request access (triggers app notification)
POST /privacy/{device_hash}/auth → Confirm access from app
GET /privacy/{device_hash}/data  → View data (requires valid session)
DELETE /privacy/{device_hash}    → Request deletion
```

## UI Mockup

```
┌─────────────────────────────────────────┐
│  🔒 Your Privacy Data                  │
│                                         │
│  Device: iPhone15,2 (masked)           │
│                                         │
│  Crashes: 2 incidents                 │
│  Events: 156 tracked                  │
│                                         │
│  [Download My Data]  [Delete All]     │
│                                         │
│  Last access: 2026-03-12 14:32        │
└─────────────────────────────────────────┘
```

## Open Questions
- How to handle devices with no active app (user uninstalled)?
- Should developer be notified of deletion requests?
- Retention period after deletion request (immediate vs 30 days)?
- Can user see data from previous device if they reinstall app?

## References
- GDPR Article 15 (Right of access)
- GDPR Article 17 (Right to erasure)
- Apple App Store Privacy Requirements
- Google Play Data Safety Section

## Acceptance Criteria
- [ ] User can access privacy dashboard via unique URL
- [ ] User must authorize access via app notification
- [ ] User can view all their data
- [ ] User can download their data
- [ ] User can request deletion
- [ ] Audit log tracks all access attempts
- [ ] Rate limiting prevents brute force

---

**Note:** This feature requires research into:
- Secure URL generation best practices
- In-app notification architecture
- Cryptographic token handling
- Legal compliance verification
