#!/bin/bash
# Integration Test for backfeedr
# Tests: Server start → Events → Database → Dashboard

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Starting backfeedr integration test..."

# Build binaries first
echo "📦 Building binaries..."
cd "$(dirname "$0")/.."
go build -o backfeedr ./cmd/backfeedr
go build -o backfeedr-client ./cmd/backfeedr-client

# Setup test environment
export BACKFEEDR_DB_PATH="./data/test.db"
export BACKFEEDR_AUTH_TOKEN="test_token_123"
export BACKFEEDR_PORT="8080"
export BACKFEEDR_BASE_URL="http://localhost:8080"

# Clean up old test data
rm -f "$BACKFEEDR_DB_PATH"
rm -rf ./data
mkdir -p ./data

echo "🚀 Starting server..."
./backfeedr > /tmp/backfeedr.log 2>&1 &
SERVER_PID=$!
sleep 3

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up..."
    kill $SERVER_PID 2>/dev/null || true
    rm -f "$BACKFEEDR_DB_PATH"
}
trap cleanup EXIT

# Test 1: Health check
echo "🏥 Test 1: Health check..."
if ! curl -fs http://localhost:8080/api/v1/health > /tmp/health.json 2>&1; then
    echo -e "${RED}❌ Health check failed${NC}"
    echo "Server log:"
    cat /tmp/backfeedr.log
    exit 1
fi
echo -e "${GREEN}✅ Health check passed${NC}"
cat /tmp/health.json

# Create test app via admin endpoint (direct DB insert for now)
echo "📱 Test 2: Creating test app..."
sqlite3 "$BACKFEEDR_DB_PATH" << 'EOF'
INSERT INTO apps (id, name, bundle_id, api_key, created_at) 
VALUES ('test-app-001', 'TestApp', 'com.example.test', 'bf_live_test123abc', datetime('now'));
EOF
TEST_API_KEY="bf_live_test123abc"
echo -e "${GREEN}✅ Test app created${NC}"

# Test 3: Send crash
echo "💥 Test 3: Sending crash report..."
cat > /tmp/crash.json << 'EOF'
{
  "exception_type": "EXC_BAD_ACCESS",
  "exception_reason": "Null pointer dereference",
  "stack_trace": [
    {"frame": 0, "symbol": "main.crash", "file": "main.go", "line": 42}
  ],
  "app_version": "1.0.0",
  "os_version": "18.3.1",
  "device_model": "iPhone16,1",
  "device_id_hash": "abc123",
  "locale": "de_DE",
  "occurred_at": "2026-03-12T14:00:00Z"
}
EOF

if ! ./backfeedr-client \
    --endpoint http://localhost:8080 \
    --api-key "$TEST_API_KEY" \
    --command send-crash \
    --file /tmp/crash.json > /tmp/crash_result.txt 2>&1; then
    echo -e "${RED}❌ Crash report failed${NC}"
    cat /tmp/crash_result.txt
    exit 1
fi
echo -e "${GREEN}✅ Crash report sent${NC}"
cat /tmp/crash_result.txt

# Test 4: Send events (single + batch)
echo "📊 Test 4: Sending events..."

# Single event
if ! ./backfeedr-client \
    --endpoint http://localhost:8080 \
    --api-key "$TEST_API_KEY" \
    --command send-event \
    --type session_start > /tmp/event_result.txt 2>&1; then
    echo -e "${RED}❌ Event send failed${NC}"
    cat /tmp/event_result.txt
    exit 1
fi
echo -e "${GREEN}✅ Single event sent${NC}"

# Batch events
cat > /tmp/batch.json << 'EOF'
{
  "events": [
    {"type": "custom", "name": "button_click", "properties": {"button": "submit"}, "app_version": "1.0.0", "occurred_at": "2026-03-12T14:01:00Z"},
    {"type": "error", "name": "network_error", "app_version": "1.0.0", "occurred_at": "2026-03-12T14:02:00Z"},
    {"type": "session_end", "app_version": "1.0.0", "occurred_at": "2026-03-12T14:03:00Z"}
  ]
}
EOF

if ! ./backfeedr-client \
    --endpoint http://localhost:8080 \
    --api-key "$TEST_API_KEY" \
    --command batch-events \
    --file /tmp/batch.json > /tmp/batch_result.txt 2>&1; then
    echo -e "${RED}❌ Batch events failed${NC}"
    cat /tmp/batch_result.txt
    exit 1
fi
echo -e "${GREEN}✅ Batch events sent${NC}"
cat /tmp/batch_result.txt

# Test 5: Verify database
echo "🗄️  Test 5: Verifying database..."
CRASH_COUNT=$(sqlite3 "$BACKFEEDR_DB_PATH" "SELECT COUNT(*) FROM crashes;" || echo "0")
EVENT_COUNT=$(sqlite3 "$BACKFEEDR_DB_PATH" "SELECT COUNT(*) FROM events;" || echo "0")

echo "   Crashes in DB: $CRASH_COUNT"
echo "   Events in DB: $EVENT_COUNT"

if [ "$CRASH_COUNT" -lt 1 ] || [ "$EVENT_COUNT" -lt 4 ]; then
    echo -e "${RED}❌ Database verification failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Database verification passed${NC}"

# Test 6: Dashboard accessible
echo "🖥️  Test 6: Dashboard check..."
if ! curl -fs http://localhost:8080/ > /tmp/dashboard.html 2>&1; then
    echo -e "${RED}❌ Dashboard unreachable${NC}"
    exit 1
fi
if ! grep -q "backfeedr" /tmp/dashboard.html; then
    echo -e "${RED}❌ Dashboard content missing${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Dashboard accessible${NC}"

# Test 7: Verify specific data in dashboard (HTMX response)
echo "📈 Test 7: Dashboard content..."
# Note: This will be enhanced when dashboard shows real data
echo -e "${YELLOW}⚠️  Dashboard data display - manual verification needed${NC}"

echo ""
echo -e "${GREEN}🎉 All integration tests passed!${NC}"
echo ""
echo "Summary:"
echo "  ✅ Server health"
echo "  ✅ App creation"
echo "  ✅ Crash ingestion"
echo "  ✅ Event ingestion (single + batch)"
echo "  ✅ Database storage (crashes: $CRASH_COUNT, events: $EVENT_COUNT)"
echo "  ✅ Dashboard accessible"
echo ""
