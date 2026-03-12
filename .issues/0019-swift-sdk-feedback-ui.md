---
id: 19
title: Swift SDK - User Feedback on Crash UI
status: open
priority: medium
labels: [sdk, swift, ui]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Add optional UI in SwiftUI to capture user feedback when a crash occurs.

## User Story
When my app crashes, I want to optionally ask the user what they were doing before the crash happened. This gives developers additional context beyond the technical stack trace.

## UI Flow

```
┌─────────────────────────────────────────┐
│  😞 Oops!                             │
│                                         │
│  The app crashed unexpectedly.          │
│  Would you like to help fix it?         │
│                                         │
│  [Skip]              [Yes, help]        │
└─────────────────────────────────────────┘
```

If user selects "Yes":

```
┌─────────────────────────────────────────┐
│  📝 What were you doing?                │
│                                         │
│  [Text input - optional]                │
│  eg., "uploading a photo"              │
│  "nothing special, just browsing"      │
│                                         │
│  ⚠️ Don't include personal information  │
│                                         │
│  [Cancel]            [Send]             │
└─────────────────────────────────────────┘
```

## API

```swift
// In Info.plist or SDK config
Backfeedr.configure(
    endpoint: "...",
    apiKey: "...",
    crashFeedback: CrashFeedbackConfig(
        enabled: true,
        askForComment: true,
        prompt: "What were you doing when the crash occurred?",
        hint: "Don't include personal information. Examples: 'uploading a photo', 'scrolling the feed'",
        allowUserToSkip: true
    )
)
```

## Data Model

Crash report receives additional fields:
```json
{
  "user_comment": "I was trying to upload a large video",
  "user_context": "photo_upload_screen"
}
```

## Privacy Considerations
- Clear hint: "Don't include personal information"
- Comment is optional (user can skip)
- Empty comments are silently ignored
- Client-side validation can strip emails/phone numbers

## Acceptance Criteria
- [ ] `CrashFeedbackConfig` struct
- [ ] SwiftUI view for crash feedback
- [ ] Optional user comment in crash payload
- [ ] Configurable prompt and hint text
- [ ] Skip option available
