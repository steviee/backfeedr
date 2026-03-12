# Security Policy

## Reporting Security Issues

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue
2. Email security concerns to the maintainer
3. Allow time for a fix before disclosure

## Security Features

backfeedr implements multiple security layers:

### Data Minimization
- No PII collected
- No IP addresses logged
- Device IDs hashed (SHA-256)
- Automatic data retention

### Transport Security
- TLS 1.3 required
- HSTS headers
- No HTTP fallback

### Request Authentication
- HMAC-SHA256 signatures
- 5-minute timestamp window
- API key validation

### Runtime Protection
- Rate limiting (100 req/min/key)
- Request size limits (1MB)
- Input validation

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest main | ✅ |
| older commits | ❌ |

## Security Updates

Security fixes are released as soon as possible. Update promptly.
