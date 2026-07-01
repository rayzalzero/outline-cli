# User Issue Resolution - Collection ID vs Slug

**Date:** 2026-07-01  
**Issue:** Error 403 when trying to clone "test-0Zs6CX3gQx"  
**Status:** ✅ RESOLVED

## Problem

User tried to clone collection with:
```bash
./outline clone test-0Zs6CX3gQx test
```

Got error:
```
Error: HTTP 403: {"ok":false,"error":"authorization_error","status":403}
```

## Root Cause

1. **Token not exported** - User didn't have OUTLINE_TOKEN or OUTLINE_API_KEY set in their shell
2. **Wrong ID format** - "test-0Zs6CX3gQx" is a **slug** (URL-friendly name), not a **collection ID** (UUID)

### What's the Difference?

| Type | Format | Example | Where to find |
|------|--------|---------|---------------|
| **Slug** | `name-shortcode` | `test-0Zs6CX3gQx` | In URL path |
| **Collection ID** | UUID | `25299a17-a07d-48d2-b0df-4c5b7827a719` | Need API call |

**URL Structure:**
```
https://outline-rbi.jatismobile.com/collection/test-0Zs6CX3gQx
                                              └─── slug ───┘

Collection ID is NOT visible in the URL!
```

## Solution

### Immediate Fix: New Command

Implemented `outline list` command to display all collections with their IDs:

```bash
$ outline list

Found 7 collection(s):

[1] Catatan - dev
    ID:  25299a17-a07d-48d2-b0df-4c5b7827a719
    URL: /collection/catatan-dev-LtgQVZCHeI

[2] JNS
    ID:  30f7de84-0e11-4e57-a2b1-aa7292dc6608
    URL: /collection/jns-yY1zI9VRK3
...
```

### Correct Usage

```bash
# 1. Export authentication token
export OUTLINE_TOKEN="eyJhbGci..."

# 2. List collections to get correct ID
outline list

# 3. Clone using UUID (not slug!)
outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 my-docs
```

## Implementation Details

### New File: `pkg/cmd/list.go`

- Calls `api.ListCollections()` 
- Displays collection name, ID, description, URL
- Shows usage hint at the end
- Validates authentication before API call

### Changes Made

1. Created `pkg/cmd/list.go` - New list command implementation
2. Updated `QUICKSTART.md` - Added troubleshooting for slug vs UUID
3. Rebuilt all platform binaries with new command
4. Updated command table to include `outline list`

### Git Commits

```
801aa4c docs: update quickstart with 'outline list' command
4df8fd9 feat: add 'outline list' command to display all collections
```

## Testing

Verified working with real Outline instance:

```bash
$ export OUTLINE_TOKEN="eyJhbGci..."
$ outline list

Fetching collections...

Found 7 collection(s):

[1] Catatan - dev
    ID:  25299a17-a07d-48d2-b0df-4c5b7827a719
    URL: /collection/catatan-dev-LtgQVZCHeI
...

Usage: outline clone <collection-id> <directory>
```

Then successfully cloned:

```bash
$ outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 test
Cloning collection 25299a17-a07d-48d2-b0df-4c5b7827a719...
Collection: Catatan - dev
Fetching document tree...
  [1] naufal.md
  [2] naufal/test/test.md
...
Cloned 5 documents from Catatan - dev
```

## Prevention

Added troubleshooting section in QUICKSTART.md:

```markdown
### "test-0Zs6CX3gQx" is not a valid collection ID
- Collection ID must be UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
- "test-0Zs6CX3gQx" is a **slug**, not a collection ID
- Use `outline list` to get the correct UUID
```

## Future Improvements

Considered but not implemented (may add later):

1. **Auto-resolve slug to UUID** - API call to convert slug → ID
2. **Better error message** - Detect slug pattern and suggest `outline list`
3. **URL parsing** - Accept full Outline URL and extract collection ID
4. **Fuzzy search** - `outline search "test"` to find collection by name

For now, `outline list` is sufficient and keeps the tool simple.

## Lessons Learned

1. **API uses UUIDs, not slugs** - All Outline API endpoints require collection ID (UUID), not the URL slug
2. **User confusion** - Slug is visible in URL, but ID is not - users naturally try to use what they see
3. **Discovery problem** - No easy way to get collection ID without API call
4. **Documentation gap** - Need to explain difference between slug and ID

## Status

✅ Issue resolved  
✅ New command implemented  
✅ Documentation updated  
✅ Binaries rebuilt  
✅ User can now successfully clone collections

---

**For user:** Use `outline list` to get collection IDs, then clone with the UUID.
