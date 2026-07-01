# URL Clone Feature ✨

## What's New

You can now **paste URLs directly** from your browser when cloning collections - including **document URLs**!

## Supported Formats

The `outline clone` command now accepts **7 different input formats**:

### Collection Formats

#### 1. UUID (Original)
```bash
outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 my-docs
```

#### 2. Full Collection URL (NEW! 🎉)
```bash
# Just copy-paste from browser!
outline clone https://outline-rbi.jatismobile.com/collection/catatan-dev-LtgQVZCHeI my-docs
```

#### 3. Collection Path (NEW!)
```bash
outline clone /collection/catatan-dev-LtgQVZCHeI my-docs
```

#### 4. Collection Slug (NEW!)
```bash
outline clone catatan-dev-LtgQVZCHeI my-docs
```

### Document Formats (NEW! 🎉🎉)

When you provide a document URL/slug, the tool **automatically detects and clones the parent collection**.

#### 5. Full Document URL (NEW!)
```bash
# Copy from browser address bar!
outline clone https://outline-rbi.jatismobile.com/doc/test-0Zs6CX3gQx my-docs
```

#### 6. Document Path (NEW!)
```bash
outline clone /doc/test-0Zs6CX3gQx my-docs
```

#### 7. Document Slug (NEW!)
```bash
outline clone test-0Zs6CX3gQx my-docs
```

## How It Works

### Collection URLs
1. **Extract slug** from URL pattern
2. **Resolve to UUID** via `collections.info` API
3. **Clone** using the UUID

### Document URLs
1. **Detect** if input is document URL (contains `/doc/` or matches document pattern)
2. **Resolve document** via `documents.info` API
3. **Extract parent collection ID** from document metadata
4. **Clone** the parent collection

## Example Workflow

**Before:**
```bash
# 1. Open browser → copy URL
# 2. Run outline list → find UUID
# 3. Copy UUID
# 4. outline clone <UUID> dir
```

**Now:**
```bash
# 1. Open browser → copy URL (collection OR document!)
# 2. outline clone <paste-URL> dir
# Done! ✅
```

## Technical Details

- Uses regex to extract slug from various URL formats
- Calls `collections.info` API to resolve collection slug → UUID
- Calls `documents.info` API to resolve document slug → collection ID
- Fallback logic: try collection first, then document
- Zero breaking changes (UUID still works)

## Tested Formats

All 7 formats successfully tested:

| Input | Type | Result |
|-------|------|--------|
| `25299a17-...` | UUID | ✅ Used directly |
| `https://outline.com/collection/slug-abc` | Collection URL | ✅ Resolved to UUID |
| `/collection/slug-abc` | Collection path | ✅ Resolved to UUID |
| `catatan-dev-LtgQVZCHeI` | Collection slug | ✅ Resolved to UUID |
| `https://outline.com/doc/test-xyz` | Document URL | ✅ Resolved to parent collection |
| `/doc/test-xyz` | Document path | ✅ Resolved to parent collection |
| `test-0Zs6CX3gQx` | Document slug | ✅ Resolved to parent collection |

## Benefits

✅ **Easier**: No need to look up UUIDs  
✅ **Faster**: One-step copy-paste from browser  
✅ **Flexible**: 7 input formats accepted  
✅ **Smart**: Auto-detects document → parent collection  
✅ **Backward compatible**: Old UUID commands still work  

---

Updated: 2026-07-01  
Commit: (pending)
