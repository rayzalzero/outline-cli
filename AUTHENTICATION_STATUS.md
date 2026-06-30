# Authentication Status - Outline CLI

**Date**: 2026-06-30  
**Status**: ✅ Dual Token Support Implemented

---

## Implementation Summary

### ✅ What's Implemented

1. **Automatic Token Type Detection**
   - `ol_api_` prefix → Detected as API Key
   - JWT format (3 dots) → Detected as JWT token
   - Graceful fallback to JWT if unknown format

2. **Dual Environment Variable Support**
   ```bash
   OUTLINE_API_KEY='...'  # Priority 1
   OUTLINE_TOKEN='...'    # Priority 2 (fallback)
   ```

3. **Updated API Client**
   - Generic `token` field (was `apiKey`)
   - `tokenType` field for detection result
   - Bearer token authentication for both types

4. **Documentation**
   - `AUTHENTICATION.md` - Complete guide
   - `README.md` - Quick reference with guide link
   - Troubleshooting section

---

## Current Limitation

⚠️ **Outline API Server Rejects JWT Tokens**

```bash
# JWT token test result:
curl -X POST https://outline-rbi.jatismobile.com/api/collections.info \
  -H "Authorization: Bearer eyJhbG..." \
  -d '{"id":"test-0Zs6CX3gQx"}'

# Response:
{"ok":false,"error":"authorization_error","status":403,"message":"Authorization error"}
```

**Root Cause**: Outline API requires API key (`ol_api_...`), not session tokens.

**Impact**: 
- ✅ CLI can detect and accept JWT tokens
- ❌ Server rejects them with 403
- ✅ Users must use API key from settings

---

## User Instructions

### Required: Get API Key

1. Open: https://outline-rbi.jatismobile.com/settings/api
2. Click "Create new token"
3. Copy token (format: `ol_api_...`)
4. Export:
   ```bash
   export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
   ```

### Test
```bash
outline clone <collection-id> <directory>
```

---

## Technical Details

### Token Detection Logic
```go
func detectTokenType(token string) string {
    if strings.HasPrefix(token, "ol_api_") {
        return "api_key"
    }
    if strings.Count(token, ".") == 2 {
        return "jwt"
    }
    return "jwt"
}
```

### HTTP Authorization Header
```
Authorization: Bearer <token>
```

Both API key and JWT use same header format. Server-side validates token type.

---

## Files Changed

- `pkg/api/client.go` - Token detection + generic field
- `pkg/cmd/clone.go` - OUTLINE_TOKEN fallback
- `AUTHENTICATION.md` - Complete auth guide
- `README.md` - Guide references

---

## Build Status

✅ All platforms rebuilt:
- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)

Binary size: ~9.2 MB

---

## Git Status

```
Commit: 42af263
Message: Add dual authentication support: API key + JWT token detection
Files: 4 changed, 203 insertions(+), 11 deletions(-)
```

---

## Conclusion

**Implementation**: ✅ Complete  
**Server Support**: ❌ JWT tokens rejected by API  
**User Action**: Must use API key (`ol_api_...`)

CLI is ready for use with proper API key.
