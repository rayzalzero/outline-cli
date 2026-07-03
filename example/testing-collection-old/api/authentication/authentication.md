---
outline_id: b44c33b8-f6c5-45b8-a14b-44a1baa26004
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/authentication-vFbLMx2ede
outline_updated: 2026-07-03T04:39:43.011Z
outline_revision: 2
---

# Authentication

API authentication methods and security guidelines.

## Table of Contents

- [API Keys](./api-keys.md)
- [OAuth](./oauth.md)
- [Security](./security.md)

## Authentication Methods

Testing Collection supports multiple authentication methods:

1. **API Keys** - Simple token-based authentication
2. **OAuth 2.0** - Delegated authorization
3. **JWT Tokens** - Stateless authentication

## API Key Authentication

### Generate API Key

```bash
testing-collection auth create-key --name "My API Key"
```

### Use API Key

```http
GET /api/v1/tests
Authorization: Bearer YOUR_API_KEY
```

```javascript
const response = await fetch('https://api.testing-collection.example.com/v1/tests', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
});
```

## OAuth 2.0 Flow

### Authorization Request

```http
GET /oauth/authorize?
  client_id=YOUR_CLIENT_ID&
  redirect_uri=https://yourapp.com/callback&
  response_type=code&
  scope=read:tests write:tests
```

### Token Exchange

```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=AUTHORIZATION_CODE&
client_id=YOUR_CLIENT_ID&
client_secret=YOUR_CLIENT_SECRET&
redirect_uri=https://yourapp.com/callback
```

### Response

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_token_here"
}
```

## JWT Authentication

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "your_password"
}
```

### Response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 3600
}
```

### Use JWT Token

```http
GET /api/v1/tests
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Scopes

Available permission scopes:

- `read:tests` - Read test results
- `write:tests` - Run and modify tests
- `read:config` - Read configuration
- `write:config` - Modify configuration
- `admin` - Full access

## Token Refresh

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```

## Revoke Token

```http
POST /api/v1/auth/revoke
Authorization: Bearer YOUR_TOKEN
```

## Security Best Practices

1. **Store tokens securely** - Use environment variables
2. **Rotate keys regularly** - Generate new keys periodically
3. **Use HTTPS only** - Never send tokens over HTTP
4. **Limit scope** - Request minimum required permissions
5. **Monitor usage** - Track API key activity

## Error Responses

### Invalid Token

```json
{
  "error": "invalid_token",
  "message": "The access token is invalid or expired"
}
```

### Insufficient Permissions

```json
{
  "error": "insufficient_scope",
  "message": "Token does not have required scope: write:tests"
}
```

## Next Steps

- [API Keys](./api-keys.md) - Detailed key management
- [OAuth](./oauth.md) - OAuth implementation guide
- [Security](./security.md) - Security guidelines
- [Endpoints](../endpoints/README.md) - API endpoints

