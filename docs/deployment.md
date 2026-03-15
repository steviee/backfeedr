# Deployment Guide

This guide covers deploying backfeedr to production environments.

## Requirements

- Linux server (Ubuntu 22.04 LTS recommended)
- Docker and Docker Compose
- Domain name (for HTTPS)
- SSL certificate (Let's Encrypt recommended)

## Quick Deploy with Docker

### 1. Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt install docker-compose-plugin

# Create directory
mkdir -p ~/backfeedr && cd ~/backfeedr
```

### 2. Create docker-compose.yml

```yaml
version: '3.8'

services:
  backfeedr:
    image: ghcr.io/steviee/backfeedr:latest
    # Or build from source:
    # build: .
    container_name: backfeedr
    restart: unless-stopped
    ports:
      - "127.0.0.1:8080:8080"  # Only localhost, nginx handles external
    volumes:
      - ./data:/data
    environment:
      - BACKFEEDR_DB_PATH=/data/backfeedr.db
      - BACKFEEDR_AUTH_TOKEN=${AUTH_TOKEN}
      - BACKFEEDR_BASE_URL=https://crashes.yourdomain.com
      - BACKFEEDR_RETENTION_DAYS=90
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  nginx:
    image: nginx:alpine
    container_name: backfeedr-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
      - ./certbot-data:/etc/letsencrypt:ro
    depends_on:
      - backfeedr

  # Optional: Let's Encrypt for automatic HTTPS
  certbot:
    image: certbot/certbot
    container_name: backfeedr-certbot
    volumes:
      - ./certbot-data:/etc/letsencrypt
      - ./certbot-www:/var/www/certbot
    entrypoint: "/bin/sh -c 'trap exit TERM; while :; do certbot renew; sleep 12h & wait $${!}; done;'"
```

### 3. Configure Nginx

Create `nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=100r/m;

    # HTTP to HTTPS redirect
    server {
        listen 80;
        server_name crashes.yourdomain.com;
        
        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }
        
        location / {
            return 301 https://$server_name$request_uri;
        }
    }

    # HTTPS server
    server {
        listen 443 ssl http2;
        server_name crashes.yourdomain.com;

        # SSL certificates
        ssl_certificate /etc/letsencrypt/live/crashes.yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/crashes.yourdomain.com/privkey.pem;
        
        # SSL settings
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers on;
        
        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;

        # Proxy to backfeedr
        location / {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://backfeedr:8080;
            proxy_http_version 1.1;
            
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }
    }
}
```

### 4. Set Up SSL with Let's Encrypt

```bash
# Get initial certificate
docker run -it --rm \
  -v "$(pwd)/certbot-data:/etc/letsencrypt" \
  -v "$(pwd)/certbot-www:/var/www/certbot" \
  -p 80:80 \
  certbot/certbot certonly \
  --standalone \
  -d crashes.yourdomain.com \
  --agree-tos \
  --email your-email@example.com
```

### 5. Start Services

```bash
# Create .env file
echo "AUTH_TOKEN=$(openssl rand -hex 32)" > .env

# Start
docker-compose up -d

# Check logs
docker-compose logs -f
```

## Manual Deployment (Without Docker)

### 1. Install Go

```bash
# Download from https://go.dev/dl/
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. Build Application

```bash
# Clone
git clone https://github.com/steviee/backfeedr.git
cd backfeedr

# Build
go build -o backfeedr ./cmd/backfeedr

# Create systemd service
sudo tee /etc/systemd/system/backfeedr.service > /dev/null <<EOF
[Unit]
Description=backfeedr crash reporting server
After=network.target

[Service]
Type=simple
User=backfeedr
Group=backfeedr
WorkingDirectory=/opt/backfeedr
ExecStart=/opt/backfeedr/backfeedr
Restart=always
RestartSec=5
Environment="BACKFEEDR_PORT=8080"
Environment="BACKFEEDR_DB_PATH=/opt/backfeedr/data/backfeedr.db"
Environment="BACKFEEDR_AUTH_TOKEN=your-token-here"

[Install]
WantedBy=multi-user.target
EOF

# Create user and directories
sudo useradd -r -s /bin/false backfeedr
sudo mkdir -p /opt/backfeedr/data
sudo cp backfeedr /opt/backfeedr/
sudo chown -R backfeedr:backfeedr /opt/backfeedr

# Start service
sudo systemctl daemon-reload
sudo systemctl enable backfeedr
sudo systemctl start backfeedr
```

### 3. Configure Nginx (Same as Docker setup)

## Security Checklist

- [ ] Use HTTPS only (redirect HTTP to HTTPS)
- [ ] Set strong `BACKFEEDR_AUTH_TOKEN`
- [ ] Enable rate limiting
- [ ] Configure firewall (only 80/443 open)
- [ ] Regular backups of `./data` directory
- [ ] Monitor logs for suspicious activity
- [ ] Keep system and Docker images updated

## Backup Strategy

### Automated Backups

```bash
# Create backup script
cat > backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backups/backfeedr"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
mkdir -p $BACKUP_DIR
tar czf $BACKUP_DIR/backfeedr_$DATE.tar.gz /opt/backfeedr/data

# Keep only last 7 days
find $BACKUP_DIR -name "backfeedr_*.tar.gz" -mtime +7 -delete
EOF

chmod +x backup.sh

# Add to crontab (daily at 2 AM)
echo "0 2 * * * /opt/backfeedr/backup.sh" | crontab -
```

### Restore from Backup

```bash
# Stop service
docker-compose down
# or: sudo systemctl stop backfeedr

# Restore data
tar xzf backfeedr_20240315_020000.tar.gz -C /

# Start service
docker-compose up -d
# or: sudo systemctl start backfeedr
```

## Monitoring

### Health Checks

```bash
# Simple health check
curl -f https://crashes.yourdomain.com/api/v1/health

# With authentication
curl -f https://crashes.yourdomain.com/api/v1/overview
```

### Prometheus Metrics (Coming Soon)

```
/metrics - Prometheus-compatible metrics endpoint
```

### Log Monitoring

```bash
# View logs
docker-compose logs -f backfeedr

# Or with systemd
sudo journalctl -u backfeedr -f
```

## Scaling

### Vertical Scaling

Increase resources:

```yaml
services:
  backfeedr:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
```

### Horizontal Scaling (Advanced)

For high traffic, use:
- Load balancer (nginx/haproxy)
- Shared SQLite (Litestream for replication)
- Multiple backfeedr instances

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs backfeedr

# Check permissions
ls -la data/

# Test config
docker-compose config
```

### Database locked

SQLite WAL mode handles most cases. If locked:

```bash
# Stop container
docker-compose stop backfeedr

# Check for lock files
ls -la data/*.db*

# Remove if stale (be careful!)
rm data/*.db-shm data/*.db-wal

# Restart
docker-compose start backfeedr
```

### SSL Certificate Issues

```bash
# Renew manually
docker-compose run --rm certbot renew

# Check certificate
openssl x509 -in certbot-data/live/crashes.yourdomain.com/cert.pem -text -noout
```

## Updates

### Update Docker Image

```bash
# Pull latest
docker-compose pull

# Restart with new image
docker-compose up -d

# Check status
docker-compose ps
```

### Update from Source

```bash
# Pull latest code
git pull origin main

# Rebuild
docker-compose build

# Restart
docker-compose up -d
```

## Support

- 🐛 [Open an Issue](https://github.com/steviee/backfeedr/issues)
- 💡 [Start a Discussion](https://github.com/steviee/backfeedr/discussions)
