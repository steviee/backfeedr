# Examples

Example JSON payloads for the backfeedr API.

## Usage with reference client

```bash
# Send a crash report
backfeedr-client --endpoint https://crashes.example.com \
  --api-key bf_live_... \
  --command send-crash \
  --file crash.json

# Send an event
backfeedr-client --endpoint https://crashes.example.com \
  --api-key bf_live_... \
  --command send-event \
  --type custom \
  --file event.json

# Send batch events
backfeedr-client --endpoint https://crashes.example.com \
  --api-key bf_live_... \
  --command batch-events \
  --file batch.json
```

## Files

| File | Description |
|------|-------------|
| `crash.json` | Example crash report with stack trace |
| `event.json` | Example custom event with properties |
| `batch.json` | Example batch of multiple events |
