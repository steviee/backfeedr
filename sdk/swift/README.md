# BackfeedrKit

iOS SDK for [backfeedr](https://github.com/steviee/backfeedr) — self-hosted crash reporting.

**Location:** This SDK is part of the main backfeedr repository under `sdk/swift/`.

## Installation

### Swift Package Manager

Add to your `Package.swift`:

```swift
.package(url: "https://github.com/steviee/backfeedr.git", from: "1.0.0")
```

Or in Xcode: File → Add Package Dependencies → `https://github.com/steviee/backfeedr.git`

Then add `BackfeedrKit` as a dependency to your target:

```swift
.target(
    name: "YourApp",
    dependencies: ["BackfeedrKit"]
)
```

## Setup

```swift
import BackfeedrKit

@main
struct MyApp: App {
    init() {
        Backfeedr.shared.configure(
            endpoint: "https://crashes.yourserver.com",
            apiKey: "bf_live_..."
        )
    }
}
```

## Usage

### Crash Reporting

```swift
// Automatic crash handling is enabled by default
// Manual non-fatal error reporting:
do {
    try riskyOperation()
} catch {
    Backfeedr.shared.capture(error, context: ["screen": "checkout"])
}
```

### Event Tracking

```swift
Backfeedr.shared.track("purchase_completed", properties: [
    "plan": "pro",
    "amount": 29.99
])

// Breadcrumbs for crash context
Backfeedr.shared.breadcrumb("User tapped checkout button")
```

## Configuration

```swift
let config = Backfeedr.Configuration(
    endpoint: "https://...",
    apiKey: "...",
    scrubPII: true,           // Remove personal data (default: true)
    enableHMAC: true,         // Sign requests (default: true)
    sessionTimeout: 300       // Session timeout in seconds
)

Backfeedr.shared.configure(endpoint: "...", apiKey: "...", config: config)
```

## Privacy

- ✅ No personal data collected by design
- ✅ On-device PII scrubbing
- ✅ Device IDs hashed (SHA-256)
- ✅ No third-party tracking

## Requirements

- iOS 15.0+
- macOS 12.0+
- Swift 5.9+

## License

MIT
