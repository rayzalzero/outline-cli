---
outline_id: 45069dbd-06c3-40d1-a36e-021d6a959b79
outline_collection: Testing
outline_url: /doc/oauth-20-integration-pQovQ9x2KQ
outline_updated: 2026-07-03T04:39:43.688Z
outline_revision: 4
outline_parent_id: b44c33b8-f6c5-45b8-a14b-44a1baa26004
---

# OAuth 2.0 Integration

Implement OAuth 2.0 authentication for Testing Collection.

## OAuth Flow Overview

Testing Collection supports the Authorization Code flow with PKCE for secure authentication.

## Register OAuth Application

### Via Dashboard


1. Navigate to Settings > OAuth Applications
2. Click "Create New Application"
3. Fill in application details
4. Save client ID and secret

### Via API

```http
POST /api/v1/oauth/applications
Authorization: Bearer YOUR_ADMIN_TOKEN
Content-Type: application/json

{
  "name": "My Application",
  "redirectUris": ["https://myapp.com/callback"],
  "scopes": ["read:tests", "write:tests"]
}
```

### Response

```json
{
  "clientId": "oauth_client_abc123",
  "clientSecret": "oauth_secret_xyz789",
  "name": "My Application",
  "redirectUris": ["https://myapp.com/callback"]
}
```

## Authorization Code Flow

### Step 1: Authorization Request

Redirect user to authorization endpoint:

```
https://api.testing-collection.example.com/oauth/authorize?
  response_type=code&
  client_id=oauth_client_abc123&
  redirect_uri=https://myapp.com/callback&
  scope=read:tests write:tests&
  state=random_state_string
```

### Step 2: User Authorization

User logs in and grants permissions.

### Step 3: Authorization Code

User redirected back with code:

```
https://myapp.com/callback?
  code=auth_code_xyz&
  state=random_state_string
```

### Step 4: Token Exchange

Exchange code for access token:

```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=auth_code_xyz&
client_id=oauth_client_abc123&
client_secret=oauth_secret_xyz789&
redirect_uri=https://myapp.com/callback
```

### Response

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_abc123",
  "scope": "read:tests write:tests"
}
```

## PKCE Flow (Recommended)

### Step 1: Generate Code Verifier

```javascript
const codeVerifier = generateRandomString(128);
const codeChallenge = base64UrlEncode(sha256(codeVerifier));
```

### Step 2: Authorization Request with PKCE

```
https://api.testing-collection.example.com/oauth/authorize?
  response_type=code&
  client_id=oauth_client_abc123&
  redirect_uri=https://myapp.com/callback&
  scope=read:tests write:tests&
  state=random_state_string&
  code_challenge=BASE64_URL_ENCODED_CHALLENGE&
  code_challenge_method=S256
```

### Step 3: Token Exchange with PKCE

```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=auth_code_xyz&
client_id=oauth_client_abc123&
redirect_uri=https://myapp.com/callback&
code_verifier=ORIGINAL_CODE_VERIFIER
```

## Refresh Token Flow

### Request New Access Token

```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&
refresh_token=refresh_abc123&
client_id=oauth_client_abc123&
client_secret=oauth_secret_xyz789
```

### Response

```json
{
  "access_token": "new_access_token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "new_refresh_token"
}
```

## Implementation Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

async function getAccessToken(code) {
  const response = await axios.post('https://api.testing-collection.example.com/oauth/token', {
    grant_type: 'authorization_code',
    code: code,
    client_id: process.env.CLIENT_ID,
    client_secret: process.env.CLIENT_SECRET,
    redirect_uri: 'https://myapp.com/callback'
  });
  
  return response.data.access_token;
}
```

### Python

```python
import requests

def get_access_token(code):
    response = requests.post(
        'https://api.testing-collection.example.com/oauth/token',
        data={
            'grant_type': 'authorization_code',
            'code': code,
            'client_id': os.environ['CLIENT_ID'],
            'client_secret': os.environ['CLIENT_SECRET'],
            'redirect_uri': 'https://myapp.com/callback'
        }
    )
    return response.json()['access_token']
```

## Revoke Token

```http
POST /oauth/revoke
Content-Type: application/x-www-form-urlencoded

token=access_token_or_refresh_token&
client_id=oauth_client_abc123&
client_secret=oauth_secret_xyz789
```

## Error Responses

### Invalid Client

```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

### Invalid Grant

```json
{
  "error": "invalid_grant",
  "error_description": "Authorization code is invalid or expired"
}
```

### Unauthorized Client

```json
{
  "error": "unauthorized_client",
  "error_description": "Client is not authorized for this grant type"
}
```

## Security Best Practices


1. **Use PKCE** for public clients
2. **Validate state parameter** to prevent CSRF
3. **Store secrets securely** - never expose client secret
4. **Use HTTPS only** for all OAuth endpoints
5. **Implement token refresh** before expiration
6. **Revoke tokens** when no longer needed

## Next Steps

* [API Keys](./api-keys.md)
* [Security](./security.md)
* [Endpoints](../endpoints/README.md)