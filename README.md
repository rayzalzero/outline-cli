# Outline CLI

> Git-like workflow for Outline wiki synchronization

A command-line tool for managing Outline wiki documents with a familiar Git-like interface. Clone collections, pull updates, push changes, and track working tree status.

## Features

- 🔄 **Git-like workflow** - `clone`, `pull`, `push`, `status` commands
- 📦 **Single binary** - No dependencies, just download and run
- 🌍 **Cross-platform** - Linux, macOS, Windows
- 🔐 **Secure** - API key authentication
- 📝 **Markdown-first** - Edit locally with your favorite editor
- 🔀 **Conflict detection** - 3-way merge logic (local, remote, manifest)

## Installation

### Download Binary

Download the latest release for your platform:

```bash
# Linux (amd64)
curl -L https://github.com/rayzalzero/outline-cli/releases/latest/download/outline-linux-amd64 -o outline
chmod +x outline
sudo mv outline /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/rayzalzero/outline-cli/releases/latest/download/outline-darwin-arm64 -o outline
chmod +x outline
sudo mv outline /usr/local/bin/

# Windows (amd64)
# Download from: https://github.com/rayzalzero/outline-cli/releases/latest/download/outline-windows-amd64.exe
```

### Build from Source

```bash
git clone https://github.com/rayzalzero/outline-cli.git
cd outline-cli
make build

# Install to /usr/local/bin
sudo mv bin/outline /usr/local/bin/
```

## Quick Start

### 1. Set API Key

Create an API key at: `https://your-outline.com/settings/api`

```bash
export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

**Note**: Use API key (`ol_api_...`), not JWT session token. See [AUTHENTICATION.md](AUTHENTICATION.md) for details.

### 2. Clone a Collection

```bash
# Clone by collection ID
outline clone 2e317a13-b7fa-469f-aef8-27474cf336ed jns-docs

cd jns-docs
```

### 3. Work with Documents

```bash
# Check status
outline status

# Edit a file
vim document-1.md

# Push changes
outline push

# Pull updates
outline pull
```

## Usage

### Initialize Repository

```bash
# Create new repository in current directory
outline init
```

### Clone Collection

```bash
# Clone by collection ID
outline clone <collection-id> [directory]

# Examples
outline clone 2e317a13-b7fa-469f-aef8-27474cf336ed jns-docs
outline clone abc123 ~/outline/docs

# Clone all accessible collections (coming soon)
outline clone --all ~/outline-workspace
```

### Check Status

```bash
# Show working tree status
outline status

# Example output:
# On collection: JNS Documentation
# Your workspace is up to date with 'origin'.
# 
# Changes not pushed:
#   modified:   document-1.md
#   modified:   subfolder/nested-doc.md
```

### Pull Updates

```bash
# Pull changes from remote
outline pull

# Force overwrite local changes
outline pull --force
```

### Push Changes

```bash
# Push local changes to remote
outline push

# Force overwrite remote changes
outline push --force

# Dry run (preview what would be pushed)
outline push --dry-run
```

### Track New Files

```bash
# Add new document to Outline
outline add new-document.md

# Add all untracked files
outline add .
```

### List Collections

```bash
# List all accessible collections
outline collections list

# Show collection details
outline collections show <collection-id>
```

### Search

```bash
# Search across all documents
outline search "query"

# Search within specific collection
outline search "query" --collection <collection-id>
```

## Configuration

### Repository Config

Each cloned collection has a `.outline/config` file:

```ini
[remote "origin"]
    url = https://outline-rbi.jatismobile.com
    collection = 2e317a13-b7fa-469f-aef8-27474cf336ed

[auth]
    token = ${OUTLINE_API_KEY}

[sync]
    auto_pull = false
    conflict_strategy = prompt  # prompt | force-local | force-remote
    api_delay = 300ms
```

### Global Config (Optional)

Create `~/.outline/config` for global settings:

```ini
[user]
    name = Your Name
    email = you@example.com

[api]
    token = ${OUTLINE_API_KEY}
    default_url = https://outline.example.com
```

### Environment Variables

- `OUTLINE_API_KEY` - API token (format: `ol_api_...`, required)
- `OUTLINE_TOKEN` - Alternative to OUTLINE_API_KEY (fallback)
- `OUTLINE_BASE_URL` - Outline instance URL (default: `https://outline-rbi.jatismobile.com`)

**Important**: Use API key (`ol_api_...`) from Outline settings, not JWT session token. 
See [AUTHENTICATION.md](AUTHENTICATION.md) for complete authentication guide.

## How It Works

### Directory Structure

```
jns-docs/
├── .outline/
│   ├── config              # Repository configuration
│   ├── manifest.json       # Sync state tracker
│   └── collection.json     # Collection metadata
├── document-1.md           # Your documents
├── document-2.md
└── subfolder/
    └── nested-doc.md
```

### Manifest

The manifest tracks sync state for conflict detection:

```json
{
  "document-1.md": {
    "id": "uuid",
    "revision": 5,
    "hash": "md5hash",
    "updated": "2026-06-30T05:20:00Z",
    "collection": "JNS Documentation"
  }
}
```

### Frontmatter

Each markdown file includes YAML frontmatter with Outline metadata:

```yaml
---
outline_id: "2e317a13-b7fa-469f-aef8-27474cf336ed"
outline_collection: "JNS Documentation"
outline_url: "https://outline.example.com/doc/slug"
outline_updated: "2026-06-30T05:20:00Z"
outline_revision: 5
---

# Document Title

Your content here...
```

### Conflict Detection

3-way merge logic (local hash, remote revision, manifest state):

| Local Hash | Remote Revision | Action |
|------------|----------------|--------|
| Changed | Same | Push (local edit only) |
| Same | Different | Pull (remote edit only) |
| Changed | Different | Conflict (skip unless --force) |
| Same | Same | Skip (up to date) |

## Development

### Prerequisites

- Go 1.22+

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
outline-cli/
├── cmd/
│   └── outline/           # Main CLI entry point
├── pkg/
│   ├── api/              # Outline API client
│   ├── cmd/              # Cobra commands
│   ├── config/           # Configuration
│   ├── manifest/         # Manifest management
│   ├── markdown/         # Frontmatter parsing
│   ├── repository/       # Repository abstraction
│   └── sync/             # Sync engine
├── Makefile
├── go.mod
└── README.md
```

## Roadmap

- [x] Phase 1: Core infrastructure
- [x] Phase 2: Init & Clone
- [ ] Phase 3: Status & Diff
- [ ] Phase 4: Pull
- [ ] Phase 5: Push
- [ ] Phase 6: Add & Track
- [ ] Phase 7: Collections & Search
- [ ] Phase 8: Release & Distribution

See [PLAN.md](PLAN.md) for detailed implementation plan.

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.

## Support

- 📖 Documentation: [PLAN.md](PLAN.md)
- 🐛 Issues: https://github.com/rayzalzero/outline-cli/issues
- 💬 Discussions: https://github.com/rayzalzero/outline-cli/discussions
