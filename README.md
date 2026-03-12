# backfeedr

> Self-hosted crash reporting & app metrics for iOS indie developers

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

⚠️ **Early Development — Not Production Ready**

This project is actively being developed. APIs may change, features are incomplete, and it's not yet ready for production use. Follow our [roadmap](#-roadmap) for progress.

**One container. SQLite. Privacy-first. No vendor lock-in.**

Built for indie developers who want control over their data. No Google, no Firebase, no cloud dependencies. Just a single Docker container and your own VPS.

## ✨ Why backfeedr?

| | Firebase Crashlytics | Sentry Self-Hosted | **backfeedr** |
|---|---|---|---|
| Setup complexity | SDK + console | 10+ services, Kafka, ClickHouse | **1 container** |
| Data ownership | Google | You | **You** |
| Resource usage | N/A (cloud) | 4GB+ RAM | **256MB RAM** |
| iOS native feel | ❌ | ⚠️ | **✅ Swift-first** |
| Price | Free (lock-in) | OSS but complex | **Forever free** |

## 🚀 Quick Start

```bash
# Clone the repo
git clone https://github.com/steviee/backfeedr.git
cd backfeedr

# Start with Docker
mkdir data
docker-compose up -d

# Or build from source
make build
./backfeedr
```

Visit `http://localhost:8080` to see your dashboard.

## 📱 iOS Integration

```swift
import BackfeedrKit

@main
struct MyApp: App {
    init() {
        Backfeedr.configure(
            endpoint: "https://crashes.yourserver.com",
            apiKey: "bf_live_..."
        )
    }
}
```

See [`sdk/swift`](sdk/swift) for the full SDK.

## 🛠️ Configuration

Via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKFEEDR_PORT` | `8080` | HTTP port |
| `BACKFEEDR_DB_PATH` | `/data/backfeedr.db` | SQLite database path |
| `BACKFEEDR_AUTH_TOKEN` | auto-generated | Dashboard auth token |
| `BACKFEEDR_RETENTION_DAYS` | `90` | Data retention period |

## 🤝 Contributing

We'd love your help! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

- 🐛 Found a bug? Open an issue
- 💡 Have an idea? Start a discussion
- 🔧 Want to code? Check out our [open issues](.issues/)

All contributions are welcome, from documentation to code.

## 📅 Roadmap

| Feature | Status | Notes |
|---------|--------|-------|
| Crash ingestion API | ✅ | `POST /api/v1/crashes` |
| Event ingestion API | ✅ | `POST /api/v1/events` + batch |
| API key authentication | ✅ | `X-Backfeedr-Key` header |
| HMAC request signing | ✅ | Optional, SHA-256 |
| Web dashboard | ✅ | HTMX + Pico CSS |
| App management | ✅ | Create, rotate keys, delete |
| Daily metrics aggregation | ✅ | Background worker |
| Data retention | ✅ | Auto-cleanup after 90 days |
| iOS SDK (Swift) | ✅ | Under `sdk/swift/` |
| Go reference client | ✅ | `cmd/backfeedr-client/` |
| **End-to-end tests** | ✅ | `make test-integration` |
| **Privacy Dashboard (GDPR)** | 🔄 | Issue #18 - device-specific data view |
| Swift SDK tests | 🔄 | Unit tests |
| SwiftUI example app | 🔄 | Demo app |
| Email/Slack alerts | ⏳ | Webhook notifications |
| Symbolication | ⏳ | dSYM upload & processing |
| Multi-user accounts | ⏳ | Team access |
| Grafana export | ⏳ | Metrics integration |

**Legend:** ✅ Done | 🔄 In Progress | ⏳ Planned

Check our [open issues](.issues/) for details.

## 📖 Documentation

- [API Reference](docs/API.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [iOS SDK Guide](docs/IOS_SDK.md)

## 🏗️ Architecture

```
[iOS App] ──HTTPS/JSON──> [backfeedr Go Server] ──> [SQLite]
                                  │
                          [HTMX Dashboard]
```

- **Go 1.22+** — Fast, single binary
- **SQLite (WAL mode)** — Single file, no setup
- **HTMX + Alpine.js** — Modern UI without build steps
- **Pico CSS** — Clean, responsive design

## 🔒 Security

- TLS 1.3 only
- HMAC request signing
- No PII collection by design
- API key authentication
- Rate limiting

See [SECURITY.md](SECURITY.md) for details.

## 📜 License

MIT © Stephan E. — see [LICENSE](LICENSE)

---

<p align="center">
  Made with ❤️ for the iOS indie dev community
</p>
