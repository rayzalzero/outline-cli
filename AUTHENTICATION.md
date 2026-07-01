# Authentication Guide

## Token Types

Outline CLI supports two types of authentication tokens:

### 1. API Key (Recommended)

**Format:** `ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

**Where to get:**
1. Login to Outline
2. Go to Settings → API
3. Create new API key
4. Copy the key (starts with `ol_api_`)

**Permissions:**
- ✅ Read documents
- ✅ Update documents
- ✅ Create documents (collection-level)
- ✅ Delete documents

**Usage:**
```bash
export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
outline push
```

### 2. Session Token (JWT)

**Format:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx.xxxxx`

**Where to get:**
1. Login to Outline in browser
2. Open DevTools → Application → Cookies
3. Copy `accessToken` cookie value

**Permissions:**
- ✅ Read documents
- ✅ Update documents
- ⚠️  Create documents (limited - see below)
- ❌ Delete documents

**Limitation:**
Session tokens typically cannot create documents at collection root level.
The CLI automatically works around this by creating documents as children of existing documents.

**Usage:**
```bash
export OUTLINE_TOKEN='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx.xxxxx'
outline push
```

## How Token Detection Works

The CLI automatically detects token type:

```go
if strings.HasPrefix(token, "ol_api_") {
    // API Key → Authorization: Bearer header
} else {
    // JWT Session Token → Cookie: accessToken header
}
```

## Create Document Permission Strategy

When creating new documents, the CLI uses this fallback strategy:

1. **Try collection-level create** (with `collectionId`)
   - Works with API keys
   - Usually fails with session tokens (403 authorization_error)

2. **Fallback to parent document** (with `parentDocumentId`)
   - Finds existing document in same directory
   - Or uses any existing document as parent
   - Works with session tokens ✅

This allows session tokens to create documents by attaching them to existing documents instead of creating at collection root.

## HTTP Headers

### API Key Request
```http
POST /api/documents.create
Authorization: Bearer ol_api_xxxxxxxx
Content-Type: application/json
Accept: application/json
```

### Session Token Request
```http
POST /api/documents.create
Cookie: accessToken=eyJhbGci...
Content-Type: application/json
Accept: application/json
x-api-version: 3
x-editor-version: 13.0.0
```

## Troubleshooting

### Error: "authorization_error" on push

**Cause:** Session token lacks permission to create documents at collection root.

**Solution:** CLI automatically falls back to `parentDocumentId`. No action needed.

### Error: "no permission to create documents in this collection"

**Cause:** Both `collectionId` and `parentDocumentId` strategies failed.

**Solutions:**
1. Use API key instead: `export OUTLINE_API_KEY='ol_api_...'`
2. Check user permissions in Outline settings
3. Ensure at least one document exists in the collection (for parent fallback)

### Session token expired

**Symptoms:** 403 errors on all requests

**Solution:** Get fresh session token from browser:
1. Logout and login again in browser
2. Copy new `accessToken` from DevTools → Cookies
3. Update `OUTLINE_TOKEN` environment variable

## Security Notes

- Never commit tokens to git
- Use environment variables or config files with `.gitignore`
- API keys have broader permissions - use session tokens for temporary access
- Session tokens expire (typically 30 days)
- API keys don't expire unless revoked

## Examples

### Using API Key (full permissions)
```bash
export OUTLINE_API_KEY='ol_api_xxxxx'
outline clone <collection-id> docs
cd docs
echo "# New Doc" > new.md
outline add new.md
outline push  # ✅ Creates at collection root
```

### Using Session Token (limited permissions)
```bash
export OUTLINE_TOKEN='eyJhbGci...'
outline clone <collection-id> docs
cd docs
echo "# New Doc" > new.md
outline add new.md
outline push  # ✅ Creates as child of existing document (auto-fallback)
```

Both work! Session token just uses a different strategy under the hood.
