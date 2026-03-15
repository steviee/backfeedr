# backfeedr Documentation

> Self-hosted crash reporting & app metrics for iOS indie developers

## Quick Links

- [Getting Started](getting-started.md) - Installation & Setup
- [API Reference](API.md) - Complete API documentation
- [Dashboard Guide](dashboard.md) - Using the web dashboard
- [iOS SDK](ios-sdk.md) - Swift SDK documentation
- [Deployment](deployment.md) - Production deployment
- [Contributing](../CONTRIBUTING.md) - How to contribute

## What is backfeedr?

backfeedr is a **self-hosted crash reporting and analytics platform** designed specifically for iOS indie developers who want:

- ✅ **Full data ownership** - Your data stays on your server
- ✅ **Privacy-first** - No third-party tracking, no data sharing
- ✅ **Simple setup** - Single Docker container, SQLite database
- ✅ **iOS-native** - Built with Swift developers in mind
- ✅ **Forever free** - Open source, no usage limits

## Features

### Backend
- **Crash Reporting** - Automatic crash detection with stack traces
- **Event Tracking** - Custom events, sessions, user flows
- **Real-time Dashboard** - Live metrics and visualizations
- **Data Retention** - Automatic cleanup after 90 days
- **API Authentication** - API keys with optional HMAC signing
- **Rate Limiting** - 100 requests/minute per key

### Dashboard
- **Crash Overview** - See all crashes at a glance
- **Crash Detail View** - Full stack traces with device info
- **Time Filtering** - Filter by 24h, 7d, 30d, 90d
- **Interactive Charts** - Line, doughnut, and bar charts
- **Crash Grouping** - Automatic grouping by exception type
- **Device Analytics** - See which devices are affected

### iOS SDK (Swift)
- **Automatic Crash Detection** - Catches uncaught exceptions
- **Manual Error Reporting** - Report non-fatal errors
- **Event Tracking** - Track user actions and sessions
- **Offline Queue** - Stores crashes when offline
- **PII Scrubbing** - Removes personal data before sending
- **Lightweight** - Minimal impact on app performance

## Architecture

```
┌─────────────┐     HTTPS/JSON      ┌─────────────────┐
│   iOS App   │ ────────────────────> │  backfeedr      │
│  (Swift)    │                     │  (Go Server)    │
└─────────────┘                     │                 │
                                    │  • SQLite DB    │
┌─────────────┐     HTTP            │  • HTMX UI      │
│   Browser   │ ──────────────────> │  • REST API     │
│  (Dashboard)│                     └─────────────────┘
└─────────────┘
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/steviee/backfeedr.git
cd backfeedr

# Start with Docker
mkdir data
docker-compose up -d

# Or build from source
make build
./backfeedr

# Visit dashboard
open http://localhost:8080
```

## Project Status

⚠️ **Early Development** - APIs may change, features are incomplete. Not yet production-ready.

See our [Roadmap](../README.md#-roadmap) for planned features.

## Support

- 🐛 [Open an Issue](https://github.com/steviee/backfeedr/issues)
- 💡 [Start a Discussion](https://github.com/steviee/backfeedr/discussions)
- 📧 Contact: See repository for maintainer contact

## License

MIT © Stephan E. - See [LICENSE](../LICENSE) for details.
