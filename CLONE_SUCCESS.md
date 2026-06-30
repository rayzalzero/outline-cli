# Clone Command - Successfully Working

## Test Results (2026-07-01)

✅ **Clone command fully functional with JWT session token**

### Test Case
```bash
export OUTLINE_TOKEN="eyJhbGci..." # JWT session token
outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 catatan-dev
```

### Results
- ✅ Collection metadata fetched: "Catatan - dev"
- ✅ Document tree fetched: 5 documents
- ✅ All documents downloaded with correct titles
- ✅ Frontmatter correctly populated with Outline metadata
- ✅ Manifest.json created with correct structure
- ✅ Config file created with proper format
- ✅ Nested directory structure preserved

### Sample Output
```
Cloning collection 25299a17-a07d-48d2-b0df-4c5b7827a719...
Collection: Catatan - dev
Fetching document tree...
  [1] naufal.md
  [2] naufal/test/test.md
  [3] naufal/test/coba-ada-2/coba-ada-2.md
  [4] naufal/test/coba-ada-ini/coba-ada-ini.md
  [5] naufal/untitled/untitled.md

Cloned 5 documents from Catatan - dev
```

### File Structure
```
catatan-dev/
├── .outline/
│   ├── config
│   ├── collection.json
│   └── manifest.json
├── naufal.md
└── naufal/
    ├── test/
    │   ├── test.md
    │   ├── coba-ada-2/
    │   │   └── coba-ada-2.md
    │   └── coba-ada-ini/
    │       └── coba-ada-ini.md
    └── untitled/
        └── untitled.md
```

### Frontmatter Example (naufal.md)
```yaml
---
outline_id: af7c13a3-43d8-4767-955a-0b78c5a3cfe0
outline_collection: Catatan - dev
outline_url: /doc/naufal-fsNXr2Zxj4
outline_updated: 2026-06-29T07:07:01.154Z
outline_revision: 3
---
```

### Manifest Example
```json
{
    "naufal.md": {
        "id": "af7c13a3-43d8-4767-955a-0b78c5a3cfe0",
        "revision": 3,
        "hash": "ce626639d007affd916ded717f1aa8ab",
        "updated": "2026-06-29T07:07:01.154Z",
        "collection": "Catatan - dev"
    }
}
```

## Bug Fixes Applied

### 1. Collections API Response Parsing
**Problem:** Double-parsing wrapper response
```go
// BEFORE (wrong)
var docsResp GetCollectionDocumentsResponse
json.Unmarshal(resp.Data, &docsResp)
return docsResp.Data

// AFTER (correct)
var nodes []DocumentNode
json.Unmarshal(resp.Data, &nodes)
return nodes
```

### 2. Documents API Response Parsing
**Problem:** Same double-parsing issue
```go
// BEFORE (wrong)
var docResp GetDocumentResponse
json.Unmarshal(resp.Data, &docResp)
return &docResp.Data

// AFTER (correct)
var doc Document
json.Unmarshal(resp.Data, &doc)
return &doc
```

### Root Cause
The `client.post()` method already extracts `data` field from API response:
```
API response: {data: {...}, ok: true}
client.post() returns: Response{Data: json.RawMessage(of {...})}
```

So we should unmarshal `resp.Data` directly to the target struct, not to a wrapper.

## Authentication Status

✅ **Both authentication methods work:**
1. `OUTLINE_API_KEY` - API token (ol_api_...)
2. `OUTLINE_TOKEN` - JWT session token (eyJhbGci...)

Priority: OUTLINE_API_KEY → OUTLINE_TOKEN

## Next Steps

1. ✅ Phase 1: outline init - DONE
2. ✅ Phase 2: outline clone - DONE
3. 🔄 Phase 3: outline status - Stub implemented
4. ⏳ Phase 4: outline pull
5. ⏳ Phase 5: outline push
