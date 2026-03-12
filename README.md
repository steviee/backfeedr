# backfeedr

Self-Hosted Crash Reporting & App Metrics for iOS Indie Devs

One container. SQLite. Privacy-first. No vendor lock-in.

## Quick Start

```bash
# Clone and build
git clone https://github.com/steviee/backfeedr.git
cd backfeedr
make build

# Run with Docker Compose
mkdir data
docker-compose up -d
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKFEEDR_PORT` | `8080` | HTTP port |
| `BACKFEEDR_DB_PATH` | `/data/backfeedr.db` | SQLite database path |
| `BACKFEEDR_AUTH_TOKEN` | auto-generated | Dashboard auth token |
| `BACKFEEDR_BASE_URL` | `http://localhost:8080` | External URL |
| `BACKFEEDR_RETENTION_DAYS` | `90` | Data retention |
| `BACKFEEDR_MAX_BODY_SIZE` | `1MB` | Max request size |
| `BACKFEEDR_RATE_LIMIT` | `100` | Requests per minute |

## API

### Authentication
All ingestion requests require an API key:
```
X-Backfeedr-Key: bf_live_...
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/crashes` | Submit crash report |
| POST | `/api/v1/events` | Submit event |
| POST | `/api/v1/events/batch` | Submit events batch |
| GET | `/api/v1/health` | Health check |

## iOS SDK

See [backfeedr/backfeedr-swift](https://github.com/steviee/backfeedr-swift) for the Swift SDK.

## License

MIT
