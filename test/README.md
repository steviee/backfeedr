# Testing

## Integration Test

Run the full integration test:

```bash
./test/integration.sh
```

This test:
1. Builds server and client
2. Starts the server
3. Creates a test app
4. Sends crash reports
5. Sends events (single + batch)
6. Verifies database contents
7. Checks dashboard accessibility

## Manual Testing

```bash
# Build
make build

# Start server
./backfeedr

# In another terminal, send test data
./backfeedr-client --endpoint http://localhost:8080 \
  --api-key bf_live_... \
  --command send-crash \
  --file examples/crash.json

# Check dashboard
open http://localhost:8080
```
