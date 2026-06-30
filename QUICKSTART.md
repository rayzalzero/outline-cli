# Outline CLI - Quick Start Guide

## Installation

Download the appropriate binary for your platform from `bin/`:
- Linux (amd64): `outline-linux-amd64`
- Linux (arm64): `outline-linux-arm64`
- macOS (Intel): `outline-darwin-amd64`
- macOS (Apple Silicon): `outline-darwin-arm64`
- Windows: `outline-windows-amd64.exe`

```bash
# Copy to your PATH
sudo cp bin/outline-linux-amd64 /usr/local/bin/outline
sudo chmod +x /usr/local/bin/outline
```

## Authentication

Set one of these environment variables:

```bash
# Option 1: API Token (priority 1)
export OUTLINE_API_KEY="ol_api_..."

# Option 2: JWT Session Token (priority 2)
export OUTLINE_TOKEN="eyJhbGci..."
```

To get JWT token from web UI:
1. Login to Outline
2. Open browser DevTools → Application → Cookies
3. Copy `accessToken` value

## Quick Start

```bash
# Clone a collection
outline clone <collection-id> <directory-name>

# Example
outline clone 25299a17-a07d-48d2-b0df-4c5b7827a719 my-docs

# Check status (not implemented yet)
cd my-docs
outline status
```

## What Gets Created

```
my-docs/
├── .outline/
│   ├── config              # Git-like config
│   ├── collection.json     # Collection metadata
│   └── manifest.json       # File tracking
└── *.md                    # Your documents with frontmatter
```

## Frontmatter Format

Each markdown file includes:

```yaml
---
outline_id: af7c13a3-43d8-4767-955a-0b78c5a3cfe0
outline_collection: Collection Name
outline_url: /doc/slug-abc123
outline_updated: 2026-06-29T07:07:01.154Z
outline_revision: 3
---

Your document content here...
```

## Commands

| Command | Status | Description |
|---------|--------|-------------|
| `outline init` | ✅ Working | Initialize repository |
| `outline clone <id> <dir>` | ✅ Working | Clone collection |
| `outline status` | 🔄 Stub | Show file status |
| `outline pull` | ⏳ Planned | Fetch remote changes |
| `outline push` | ⏳ Planned | Upload local changes |

## Troubleshooting

### Authentication Failed
- Check token format: `ol_api_*` or `eyJhbGci*`
- Verify token is not expired (JWT)
- Ensure environment variable is exported

### Clone Failed
- Verify collection ID is correct
- Check network connection
- Ensure you have access to the collection

### Empty Files
- If you see this, it's a bug - please report
- Check CLONE_SUCCESS.md for known issues

## Next Steps

After cloning:
1. Explore the `.outline/` directory
2. Check `manifest.json` for file tracking
3. Edit documents locally
4. Wait for `outline push` implementation to sync changes back

## Development

```bash
# Build from source
make build

# Build all platforms
make build-all

# Run tests (when available)
make test
```

## Documentation

- `README.md` - Full user guide
- `PLAN.md` - Design decisions
- `AUTHENTICATION.md` - Auth system details
- `FINAL_STATUS.md` - Current project status

## Support

For issues or questions:
- Check documentation in project root
- See `FINAL_STATUS.md` for known limitations
- Review `CLONE_SUCCESS.md` for test results

---

**Version:** Phase 1-2 Complete (2026-07-01)  
**Status:** Clone command fully functional ✅  
**Next:** Status command implementation
