# Getting Started with backfeedr

This guide will help you set up backfeedr and start collecting crash reports from your iOS app.

## Prerequisites

- **Server**: Linux/macOS/Windows with Docker or Go 1.22+
- **iOS App**: Swift project with iOS 15+ target
- **Network**: Server must be accessible from iOS devices

## Step 1: Install backfeedr Server

### Option A: Docker (Recommended)

```bash
# Clone repository
git clone https://github.com/steviee/backfeedr.git
cd backfeedr

# Create data directory
mkdir -p data

# Start with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Option B: Build from Source

```bash
# Clone repository
git clone https://github.com/steviee/backfeedr.git
cd backfeedr

# Install Go dependencies
go mod download

# Build binary
make build

# Run server
./backfeedr
```

## Step 2: Configure Server

Create a `.env` file or set environment variables:

```bash
# Server settings
export BACKFEEDR_PORT=8080
export BACKFEEDR_DB_PATH=./data/backfeedr.db
export BACKFEEDR_AUTH_TOKEN=your-secret-dashboard-token

# Optional settings
export BACKFEEDR_BASE_URL=https://your-domain.com
export BACKFEEDR_RETENTION_DAYS=90
```

## Step 3: Access Dashboard

Open your browser:

```
http://localhost:8080
```

Or if running on a server:

```
http://your-server-ip:8080
```

## Step 4: Create Your First App

1. Go to **Apps** in the dashboard
2. Click **"+ New App"**
3. Enter:
   - **Name**: MyApp
   - **Bundle ID**: com.example.myapp
4. Copy the generated **API key**

⚠️ **Important**: Save this API key - it's shown only once!

## Step 5: Integrate iOS SDK

### Install SDK

Add to your `Package.swift`:

```swift
.package(url: "https://github.com/steviee/backfeedr.git", from: "1.0.0")
```

Or in Xcode: **File → Add Package Dependencies**

### Configure SDK

In your `App.swift` or `AppDelegate.swift`:

```swift
import BackfeedrKit

@main
struct MyApp: App {
    init() {
        Backfeedr.configure(
            endpoint: "https://your-server.com",
            apiKey: "bf_live_your_api_key_here"
        )
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
```

### Test Integration

Add a test button to trigger a crash:

```swift
Button("Test Crash") {
    // This will crash the app
    let array: [String] = []
    let _ = array[0] // Index out of bounds
}
```

Run your app, tap the button, then reopen. The crash should appear in the dashboard within seconds.

## Step 6: Verify Setup

### Check Server Health

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{"status":"ok"}
```

### Test with Client

```bash
# Build client
go build -o backfeedr-client ./cmd/backfeedr-client

# Check health
./backfeedr-client --endpoint http://localhost:8080 --command health

# Send test crash
./backfeedr-client --endpoint http://localhost:8080 \
  --api-key bf_live_your_key \
  --command send-crash
```

## Step 7: Production Deployment

### Using Docker Compose

```yaml
version: '3.8'

services:
  backfeedr:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - BACKFEEDR_DB_PATH=/data/backfeedr.db
      - BACKFEEDR_AUTH_TOKEN=${AUTH_TOKEN}
      - BACKFEEDR_BASE_URL=https://crashes.yourdomain.com
    restart: unless-stopped
    
  # Optional: Reverse proxy with HTTPS
  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - backfeedr
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `BACKFEEDR_PORT` | No | 8080 | HTTP port |
| `BACKFEEDR_DB_PATH` | No | ./data/backfeedr.db | SQLite database path |
| `BACKFEEDR_AUTH_TOKEN` | No | (auto) | Dashboard auth token |
| `BACKFEEDR_BASE_URL` | No | http://localhost:8080 | External URL |
| `BACKFEEDR_RETENTION_DAYS` | No | 90 | Data retention period |

### HTTPS Setup

For production, use a reverse proxy like nginx or Caddy:

```nginx
server {
    listen 443 ssl;
    server_name crashes.yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Next Steps

- 📊 Explore the [Dashboard](dashboard.md)
- 📖 Read the [API Reference](API.md)
- 🍎 Learn about [iOS SDK](ios-sdk.md)
- 🚀 Check [Deployment Guide](deployment.md)

## Troubleshooting

### Server won't start

```bash
# Check port is free
lsof -i :8080

# Check permissions on data directory
ls -la data/

# Check logs
cat /tmp/backfeedr.log
```

### iOS app not sending crashes

1. Verify API key is correct
2. Check endpoint URL is reachable from device
3. Check iOS app logs for SDK errors
4. Ensure app has internet permission

### Dashboard not accessible

1. Check firewall rules
2. Verify `BACKFEEDR_BIND_ADDR` (use `0.0.0.0` for external access)
3. Check server is running: `docker ps` or `ps aux | grep backfeedr`

## Getting Help

- 🐛 [Open an Issue](https://github.com/steviee/backfeedr/issues)
- 💡 [Start a Discussion](https://github.com/steviee/backfeedr/discussions)
- 📧 See repository for maintainer contact
