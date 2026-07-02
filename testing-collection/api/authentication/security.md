---
outline_id: 10485b23-d45d-4a93-9fa5-2a239085d9a7
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/security-best-practices-Rxsa3ghWLe
outline_updated: 2026-07-02T10:18:17.052Z
outline_revision: 2
---



# Security Best Practices

Security guidelines for Testing Collection API integration.

## Token Security

### Secure Storage

**Never hardcode tokens:**

```javascript
// Bad
const API_KEY = 'tc_live_1234567890abcdef';

// Good
const API_KEY = process.env.TESTING_COLLECTION_API_KEY;
```

### Environment Variables

```bash
# .env
TESTING_COLLECTION_API_KEY=tc_live_1234567890abcdef
TESTING_COLLECTION_CLIENT_SECRET=oauth_secret_xyz789
```

### Secrets Management

Use dedicated secrets management:

- AWS Secrets Manager
- HashiCorp Vault
- Azure Key Vault
- Google Secret Manager

## HTTPS Only

**Always use HTTPS:**

```javascript
// Bad
const API_URL = 'http://api.testing-collection.example.com';

// Good
const API_URL = 'https://api.testing-collection.example.com';
```

## Token Rotation

### Regular Rotation Schedule

- API Keys: Every 90 days
- OAuth tokens: Refresh before expiration
- Client secrets: Every 180 days

### Rotation Process

```bash
# 1. Create new key
NEW_KEY=$(testing-collection auth create-key --name "Rotated Key")

# 2. Update application
export TESTING_COLLECTION_API_KEY=$NEW_KEY

# 3. Verify new key works
testing-collection run --api-key $NEW_KEY

# 4. Revoke old key
testing-collection auth revoke-key $OLD_KEY_ID
```

## Rate Limiting

### Implement Backoff

```javascript
async function apiCallWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      
      if (response.status === 429) {
        const retryAfter = response.headers.get('Retry-After');
        await sleep(retryAfter * 1000);
        continue;
      }
      
      return response;
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await sleep(Math.pow(2, i) * 1000);
    }
  }
}
```

## Input Validation

### Sanitize User Input

```javascript
function sanitizeTestFilter(filter) {
  // Remove dangerous characters
  return filter.replace(/[^a-zA-Z0-9_-]/g, '');
}

const filter = sanitizeTestFilter(userInput);
```

## CORS Configuration

### Restrict Origins

```javascript
// Server-side CORS config
app.use(cors({
  origin: ['https://myapp.com', 'https://staging.myapp.com'],
  credentials: true
}));
```

## Audit Logging

### Log API Access

```javascript
function logApiCall(endpoint, userId, status) {
  logger.info({
    timestamp: new Date().toISOString(),
    endpoint,
    userId,
    status,
    ip: request.ip
  });
}
```

## Webhook Security

### Verify Webhook Signatures

```javascript
const crypto = require('crypto');

function verifyWebhookSignature(payload, signature, secret) {
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex');
  
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(expectedSignature)
  );
}
```

## OAuth Security

### State Parameter

```javascript
// Generate random state
const state = crypto.randomBytes(32).toString('hex');

// Store in session
session.oauthState = state;

// Include in authorization URL
const authUrl = `${OAUTH_URL}?state=${state}&...`;
```

### Validate State

```javascript
// On callback
if (req.query.state !== session.oauthState) {
  throw new Error('Invalid state parameter');
}
```

## IP Whitelisting

### Restrict API Access

```javascript
const allowedIPs = ['203.0.113.0/24', '198.51.100.0/24'];

function isIPAllowed(ip) {
  return allowedIPs.some(range => ipInRange(ip, range));
}
```

## Monitoring and Alerts

### Set Up Alerts

- Failed authentication attempts
- Unusual API usage patterns
- Rate limit violations
- Token expiration warnings

### Example Alert Config

```yaml
alerts:
  - name: Failed Auth Attempts
    condition: failed_auth > 10 in 5m
    action: notify_security_team
  
  - name: Rate Limit Exceeded
    condition: rate_limit_hits > 100 in 1h
    action: notify_ops_team
```

## Incident Response

### Security Incident Checklist

1. **Detect**: Monitor for suspicious activity
2. **Contain**: Revoke compromised tokens immediately
3. **Investigate**: Review audit logs
4. **Remediate**: Rotate all affected credentials
5. **Document**: Record incident details
6. **Review**: Update security policies

### Revoke Compromised Token

```bash
# Immediate revocation
testing-collection auth revoke-key $COMPROMISED_KEY_ID --force

# Notify security team
testing-collection auth notify-security \
  --incident "Token compromise" \
  --key-id $COMPROMISED_KEY_ID
```

## Compliance

### Data Protection

- Encrypt tokens at rest
- Use TLS 1.2+ for transit
- Implement token expiration
- Regular security audits

### GDPR Compliance

- Allow users to export data
- Implement data deletion
- Maintain audit trails
- Document data processing

## Security Checklist

- [ ] Use HTTPS for all API calls
- [ ] Store tokens in environment variables
- [ ] Implement token rotation
- [ ] Enable rate limiting
- [ ] Validate all user input
- [ ] Configure CORS properly
- [ ] Enable audit logging
- [ ] Verify webhook signatures
- [ ] Use state parameter in OAuth
- [ ] Set up security monitoring
- [ ] Have incident response plan
- [ ] Regular security reviews

## Next Steps

- [API Keys](./api-keys.md)
- [OAuth](./oauth.md)
- [Endpoints](../endpoints/README.md)

