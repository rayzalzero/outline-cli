---
outline_id: 16d339c1-4f93-47e5-a8fd-2c39065b0f51
outline_collection: Testing
outline_url: /doc/api-keys-I9ijsQGzm2
outline_updated: 2026-07-03T04:39:43.605Z
outline_revision: 4
outline_parent_id: b44c33b8-f6c5-45b8-a14b-44a1baa26004
---

# API Keys

Manage API keys for Testing Collection authentication.

## Creating API Keys

### CLI Command

```bash
testing-collection auth create-key --name "Production API Key"
```

### API Endpoint

```http
POST /api/v1/auth/keys
Authorization: Bearer YOUR_ADMIN_TOKEN
Content-Type: application/json

{
  "name": "Production API Key",
  "scopes": ["read:tests", "write:tests"],
  "expiresIn": 2592000
}
```

### Response

```json
{
  "key": "tc_live_1234567890abcdef",
  "name": "Production API Key",
  "scopes": ["read:tests", "write:tests"],
  "createdAt": "2026-07-02T08:43:48.419Z",
  "expiresAt": "2026-08-01T08:43:48.419Z"
}
```

## Listing API Keys

### CLI Command

```bash
testing-collection auth list-keys
```

### API Endpoint

```http
GET /api/v1/auth/keys
Authorization: Bearer YOUR_ADMIN_TOKEN
```

### Response

```json
{
  "keys": [
    {
      "id": "key_123",
      "name": "Production API Key",
      "scopes": ["read:tests", "write:tests"],
      "lastUsed": "2026-07-02T08:00:00.000Z",
      "createdAt": "2026-07-01T08:43:48.419Z"
    }
  ]
}
```

## Revoking API Keys

### CLI Command

```bash
testing-collection auth revoke-key key_123
```

### API Endpoint

```http
DELETE /api/v1/auth/keys/key_123
Authorization: Bearer YOUR_ADMIN_TOKEN
```

## Key Scopes

### Available Scopes

* `read:tests` - View test results and history
* `write:tests` - Run tests and modify test data
* `read:config` - View configuration
* `write:config` - Modify configuration
* `read:reports` - Access test reports
* `admin` - Full administrative access

### Scope Examples

```json
{
  "scopes": ["read:tests", "read:reports"]
}
```

## Key Rotation

### Best Practices


1. Rotate keys every 90 days
2. Create new key before revoking old one
3. Update applications with new key
4. Revoke old key after verification

### Rotation Process

```bash
# Create new key
testing-collection auth create-key --name "Production API Key v2"

# Update application with new key
export API_KEY="tc_live_new_key"

# Test new key
testing-collection run --api-key $API_KEY

# Revoke old key
testing-collection auth revoke-key key_old
```

## Key Security

### Storage

```bash
# Environment variable
export TESTING_COLLECTION_API_KEY="tc_live_1234567890abcdef"

# .env file
TESTING_COLLECTION_API_KEY=tc_live_1234567890abcdef
```

### Never Commit Keys

```gitignore
# .gitignore
.env
.env.local
*.key
```

## Key Prefixes

* `tc_live_` - Production keys
* `tc_test_` - Test/development keys
* `tc_admin_` - Administrative keys

## Rate Limits by Key Type

| Key Type | Requests/Hour | Concurrent Tests |
|----------|---------------|------------------|
| Free     | 100           | 1                |
| Pro      | 1,000         | 5                |
| Enterprise | 10,000        | 20               |

## Monitoring Key Usage

### View Usage Stats

```bash
testing-collection auth stats key_123
```

### API Endpoint

```http
GET /api/v1/auth/keys/key_123/stats
Authorization: Bearer YOUR_ADMIN_TOKEN
```

### Response

```json
{
  "keyId": "key_123",
  "totalRequests": 1543,
  "lastUsed": "2026-07-02T08:00:00.000Z",
  "usage": {
    "today": 45,
    "thisWeek": 312,
    "thisMonth": 1543
  }
}
```

## Troubleshooting

### Invalid API Key

```json
{
  "error": "invalid_api_key",
  "message": "The provided API key is invalid or expired"
}
```

**Solution**: Verify key format and expiration

### Rate Limit Exceeded

```json
{
  "error": "rate_limit_exceeded",
  "message": "API rate limit exceeded. Retry after 3600 seconds"
}
```

**Solution**: Wait or upgrade plan

## Next Steps

* [OAuth](./oauth.md)
* [Security](./security.md)
* [Endpoints](../endpoints/README.md)