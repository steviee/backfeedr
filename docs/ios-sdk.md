# iOS SDK Guide

The backfeedr iOS SDK (BackfeedrKit) provides crash reporting and analytics for your Swift apps.

## Installation

### Swift Package Manager (Recommended)

In Xcode:

1. **File → Add Package Dependencies**
2. Enter: `https://github.com/steviee/backfeedr.git`
3. Select **Up to Next Major Version** (1.0.0)
4. Add to your app target

Or in `Package.swift`:

```swift
// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "YourApp",
    dependencies: [
        .package(url: "https://github.com/steviee/backfeedr.git", from: "1.0.0")
    ],
    targets: [
        .target(
            name: "YourApp",
            dependencies: ["BackfeedrKit"]
        )
    ]
)
```

## Configuration

### Basic Setup

In your `App.swift`:

```swift
import BackfeedrKit

@main
struct MyApp: App {
    init() {
        Backfeedr.configure(
            endpoint: "https://crashes.yourserver.com",
            apiKey: "bf_live_your_api_key"
        )
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
```

### Advanced Configuration

```swift
let config = Backfeedr.Configuration(
    endpoint: "https://crashes.yourserver.com",
    apiKey: "bf_live_your_api_key",
    scrubPII: true,           // Remove personal data (default: true)
    enableHMAC: true,         // Sign requests (default: true)
    sessionTimeout: 300       // Session timeout in seconds (default: 300)
)

Backfeedr.configure(
    endpoint: "https://crashes.yourserver.com",
    apiKey: "bf_live_your_api_key",
    config: config
)
```

## Automatic Crash Reporting

The SDK automatically catches uncaught exceptions. No additional code needed!

### How It Works

1. App crashes with uncaught exception
2. SDK captures stack trace
3. On next app launch, SDK sends crash report
4. Report appears in dashboard

## Manual Error Reporting

Report non-fatal errors:

```swift
do {
    try riskyOperation()
} catch {
    Backfeedr.shared.capture(error, context: [
        "screen": "checkout",
        "user_action": "purchase_attempt"
    ])
}
```

### Error Context

Add context to help debug:

```swift
Backfeedr.shared.capture(error, context: [
    "screen": "settings",
    "feature_flag": "new_ui_enabled",
    "network_status": "wifi"
])
```

⚠️ **Privacy**: Don't include personal data (emails, names, IDs). The SDK scrubs common PII patterns, but be careful.

## Event Tracking

### Custom Events

Track user actions:

```swift
Backfeedr.shared.track("purchase_completed", properties: [
    "plan": "pro",
    "amount": 29.99,
    "currency": "EUR"
])
```

### Session Tracking

Sessions are tracked automatically:

- **session_start** - When app becomes active
- **session_end** - When app goes to background

### Breadcrumbs

Leave breadcrumbs for crash context:

```swift
Backfeedr.shared.breadcrumb("User opened settings")
Backfeedr.shared.breadcrumb("User tapped 'Upgrade' button")
// ... later crash happens
// Breadcrumbs show in crash report
```

## Privacy Features

### PII Scrubbing

By default, the SDK removes:
- Email addresses
- Phone numbers
- IP addresses
- URLs with query parameters

### Opt-Out

Let users disable tracking:

```swift
// In your settings screen
Toggle("Help improve the app", isOn: $crashReportingEnabled)
    .onChange(of: crashReportingEnabled) { enabled in
        Backfeedr.shared.isEnabled = enabled
    }
```

When disabled:
- No crash reports sent
- No events tracked
- Queue is cleared

## Offline Support

The SDK works offline:

1. Crashes/events stored locally
2. Sent when connection available
3. Automatic retry with backoff

### Force Upload

Manually trigger upload:

```swift
Backfeedr.shared.flushQueue()
```

## Testing

### Test Crash

Add a debug button:

```swift
#if DEBUG
Button("Trigger Test Crash") {
    fatalError("Test crash from debug button")
}
#endif
```

### Verify Integration

Check SDK is configured:

```swift
if Backfeedr.shared.isConfigured {
    print("SDK ready!")
} else {
    print("SDK not configured")
}
```

## Best Practices

### 1. Configure Early

Configure SDK in `App.init()` or `application(_:didFinishLaunchingWithOptions:)`

### 2. Use Context

Always add context to errors:

```swift
// Good
Backfeedr.shared.capture(error, context: ["screen": "checkout"])

// Bad
Backfeedr.shared.capture(error)
```

### 3. Don't Track Sensitive Data

```swift
// Bad - includes user email
Backfeedr.shared.track("login", properties: [
    "email": user.email  // ❌ Don't do this
])

// Good - use anonymized ID
Backfeedr.shared.track("login", properties: [
    "user_id_hash": hash(user.id)  // ✅ OK
])
```

### 4. Test in Debug Builds

```swift
#if DEBUG
Backfeedr.shared.configure(
    endpoint: "http://localhost:8080",  // Local server
    apiKey: "bf_live_test"
)
#else
Backfeedr.shared.configure(
    endpoint: "https://crashes.yourserver.com",
    apiKey: "bf_live_production_key"
)
#endif
```

## Troubleshooting

### Crashes not appearing

1. Check API key is correct
2. Verify endpoint is reachable
3. Check iOS app logs:
   ```swift
   // Enable debug logging
   Backfeedr.shared.debug = true
   ```

### Symbolication

For symbolicated stack traces:

1. Build app with dSYM generation enabled
2. Upload dSYM to server (feature coming soon)
3. Or use on-device symbolication (limited)

### Network Errors

SDK automatically retries on network failure. Check:
- Server is running
- Device has internet
- Firewall allows connections

## API Reference

### Backfeedr

```swift
class Backfeedr {
    static let shared: Backfeedr
    
    func configure(endpoint: String, apiKey: String, config: Configuration?)
    func capture(_ error: Error, context: [String: Any]?)
    func track(_ name: String, properties: [String: Any]?)
    func breadcrumb(_ message: String)
    func flushQueue()
    
    var isEnabled: Bool
    var isConfigured: Bool
}
```

### Configuration

```swift
struct Configuration {
    let endpoint: String
    let apiKey: String
    let scrubPII: Bool
    let enableHMAC: Bool
    let sessionTimeout: TimeInterval
}
```

## Sample App

See `examples/SwiftUIExample` in the repository for a complete example.

## Support

- 🐛 [Open an Issue](https://github.com/steviee/backfeedr/issues)
- 💡 [Start a Discussion](https://github.com/steviee/backfeedr/discussions)
