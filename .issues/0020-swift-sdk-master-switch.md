---
id: 20
title: Swift SDK - Master Data Collection Switch
status: open
priority: high
labels: [sdk, swift, privacy]
created: 2026-03-12
updated: 2026-03-12
---

## Description
Add a kill switch to completely disable all data collection from the SDK side.

## User Story
As an app developer, I want to let users opt-out of ALL data collection with a simple setting in my app. The SDK must respect this immediately and stop all network calls.

## API

```swift
// Set via app settings
@AppStorage("enableCrashReporting") 
var enableCrashReporting: Bool = true

// In app init
func updateBackfeedr() {
    Backfeedr.shared.setEnabled(enableCrashReporting)
    // OR
    Backfeedr.shared.isEnabled = enableCrashReporting
}
```

## Behavior

### When disabled (`isEnabled = false`):
- [ ] No crash reporting
- [ ] No event tracking
- [ ] No network calls
- [ ] Queue is cleared
- [ ] Queue persistence disabled
- [ ] `Backfeedr.configure()` returns early/silent
- [ ] All SDK methods become no-ops

### When re-enabled (`isEnabled = true`):
- [ ] Normal operation resumes
- [ ] Queue persistence re-enabled
- [ ] Configure happens on next app start

## User-visible Setting Example

```swift
struct SettingsView: View {
    @AppStorage("crashReporting") var crashReporting = true
    
    var body: some View {
        Toggle("Help improve the app", isOn: $crashReporting)
            .onChange(of: crashReporting) { newValue in
                Backfeedr.shared.isEnabled = newValue
            }
        
        Text("Send anonymous crash reports and usage statistics")
            .font(.caption)
            .foregroundColor(.secondary)
    }
}
```

## Implementation Details

"""
Check at every entry point:

func capture(_ error: Error) {
    guard Backfeedr.shared.isEnabled else { return }
    // ... handle crash
}

func track(_ name: String) {
    guard Backfeedr.shared.isEnabled else { return }
    // ... track event
}

// Queue as well
func enqueue(_ item: QueueItem) {
    guard Backfeedr.shared.isEnabled else { return }
    // ... add to queue
}
"""

## Acceptance Criteria
- [ ] `isEnabled` property on Backfeedr
- [ ] All entry points check flag before processing
- [ ] Queue cleared when disabled
- [ ] No network calls when disabled
- [ ] Re-enables gracefully
- [ ] Documentation shows app settings integration
