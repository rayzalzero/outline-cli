# URL Clone Feature ✨

## What's New

You can now **paste URLs directly** from your browser when cloning collections!

## Supported Formats

The `outline clone` command now accepts **4 different input formats**:

### 1. UUID (Original)
```bash
outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 my-docs
```

### 2. Full URL (NEW! 🎉)
```bash
# Just copy-paste from browser!
outline clone https://outline-rbi.jatismobile.com/collection/catatan-dev-LtgQVZCHeI my-docs
```

### 3. Collection Path (NEW!)
```bash
outline clone /collection/catatan-dev-LtgQVZCHeI my-docs
```

### 4. Collection Slug (NEW!)
```bash
outline clone catatan-dev-LtgQVZCHeI my-docs
```

## How It Works

1. **Extract slug** from URL pattern
2. **Resolve to UUID** via Outline API (`collections.info`)
3. **Clone** using the UUID

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
# 1. Open browser → copy URL
# 2. outline clone <paste-URL> dir
# Done! ✅
```

## Technical Details

- Uses regex to extract slug from various URL formats
- Calls `collections.info` API to resolve slug → UUID
- Falls back to direct UUID if already provided
- Zero breaking changes (UUID still works)

## Tested Formats

All formats successfully tested:

| Input | Extracted | Result |
|-------|-----------|--------|
| `https://outline.com/collection/slug-abc123` | `slug-abc123` | ✅ Resolved to UUID |
| `/collection/slug-abc123` | `slug-abc123` | ✅ Resolved to UUID |
| `slug-abc123` | `slug-abc123` | ✅ Resolved to UUID |
| `25299a17-...` | (unchanged) | ✅ Used directly |

## Benefits

✅ **Easier**: No need to look up UUIDs  
✅ **Faster**: One-step copy-paste from browser  
✅ **Flexible**: Multiple input formats accepted  
✅ **Backward compatible**: Old UUID commands still work  

---

Updated: 2026-07-01  
Commit: a09ede8
