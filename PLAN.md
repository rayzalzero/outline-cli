# Outline CLI - Go Implementation Plan

> Status: Planning
> Date: 2026-06-30

## Overview

Build a cross-platform CLI tool in Go to replace existing bash scripts for Outline wiki synchronization.

## Goals

1. **Cross-platform binary** - Single build for Linux, macOS, Windows
2. **Feature parity** - Replicate all bash script functionality
3. **Better UX** - Colored output, progress bars, better error messages
4. **Maintainable** - Clean Go code with proper error handling

## Architecture (Git-like)

```
outline-cli/
├── cmd/
│   └── outline/           # Main CLI entry point
│       └── main.go
├── pkg/
│   ├── api/              # Outline API client
│   │   ├── client.go
│   │   ├── documents.go
│   │   ├── collections.go
│   │   └── auth.go
│   ├── repository/       # Git-like repository abstraction
│   │   ├── repository.go # Repository struct & methods
│   │   ├── init.go       # outline init
│   │   ├── clone.go      # outline clone
│   │   └── remote.go     # Remote management
│   ├── sync/             # Sync engine
│   │   ├── pull.go       # outline pull
│   │   ├── push.go       # outline push
│   │   ├── status.go     # outline status
│   │   ├── add.go        # outline add (track new files)
│   │   └── diff.go       # Compare working tree vs manifest
│   ├── manifest/         # Manifest management
│   │   ├── manifest.go   # Load/save/update manifest
│   │   └── entry.go      # ManifestEntry struct
│   ├── config/           # Configuration (.outline/config)
│   │   └── config.go
│   └── markdown/         # Markdown processing
│       ├── frontmatter.go # Parse/serialize frontmatter
│       └── content.go     # Strip frontmatter from content
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Repository Structure (Per Collection)

Each cloned collection becomes a repository with `.outline/` metadata:

```
jns-collection/
├── .outline/
│   ├── config              # Repository config
│   │   ├── api_key         # Encrypted API key
│   │   ├── base_url        # Outline instance URL
│   │   └── collection_id   # Current collection ID
│   ├── manifest.json       # Sync state (like git index)
│   └── collection.json     # Collection metadata
├── document-1.md
├── document-2.md
└── subfolder/
    └── nested-doc.md
```

### `.outline/config` Format

```ini
[remote "origin"]
    url = https://outline-rbi.jatismobile.com
    collection = 2e317a13-b7fa-469f-aef8-27474cf336ed

[auth]
    token = ${OUTLINE_API_KEY}

[sync]
    auto_pull = false
    conflict_strategy = prompt  # prompt | force-local | force-remote
```

### `.outline/collection.json` Format

```json
{
  "id": "2e317a13-b7fa-469f-aef8-27474cf336ed",
  "name": "JNS Documentation",
  "description": "Technical documentation for JNS system",
  "color": "#4CAF50",
  "icon": "book",
  "url": "https://outline-rbi.jatismobile.com/collection/jns",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2026-06-30T05:20:00Z"
}
```

## API Client Design

### Client Structure

```go
type Client struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
    return &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}
```

### API Methods

Based on research findings:

**Collections:**
- `ListCollections() ([]Collection, error)`
- `GetCollectionDocuments(id string) (*DocumentTree, error)`

**Documents:**
- `GetDocument(id string) (*Document, error)`
- `UpdateDocument(id, content string) (*Document, error)`
- `ListDocuments(collectionID string, limit, offset int) ([]Document, error)`
- `SearchDocuments(query, collectionID string) ([]Document, error)`

**Auth:**
- `ValidateToken() error`

## Manifest Structure

```go
type ManifestEntry struct {
    ID         string    `json:"id"`
    Revision   int       `json:"revision"`
    Hash       string    `json:"hash"`
    Updated    time.Time `json:"updated"`
    Collection string    `json:"collection"`
}

type Manifest map[string]ManifestEntry
```

## Commands (Git-like Workflow)

### 1. Init Command (NEW)

```bash
outline init
```

**Purpose:** Initialize current directory as outline workspace
**Flow:**
1. Create `.outline/` directory
2. Create empty `config` file
3. Create empty `manifest.json`
4. Prompt for API key and base URL
5. Save to `.outline/config`

**Output:**
```
Initialized empty Outline repository in /path/to/current/.outline/
```

### 2. Clone Command (NEW)

```bash
outline clone <collection-id> [directory]
outline clone <collection-url> [directory]
outline clone --all [directory]  # Clone all collections
```

**Purpose:** Clone collection(s) to local directory (like `git clone`)

**Flow:**
1. Create target directory (or use current)
2. Create `.outline/` metadata directory
3. Save config (API key from env or prompt)
4. Fetch collection metadata
5. Download all documents
6. Create manifest with initial state
7. Write files with frontmatter

**Examples:**
```bash
# Clone single collection by ID
outline clone 2e317a13-b7fa-469f-aef8-27474cf336ed jns-docs

# Clone by URL
outline clone https://outline-rbi.jatismobile.com/collection/jns jns-docs

# Clone all accessible collections
outline clone --all ~/outline-workspace
```

**Directory Structure After Clone:**
```
jns-docs/
├── .outline/
│   ├── config          # API key, base URL
│   ├── manifest.json   # Sync state
│   └── collection.json # Collection metadata (name, id, color)
├── document-1.md
├── document-2.md
└── subfolder/
    └── nested-doc.md
```

### 3. Pull Command

```bash
outline pull [--force] [--rebase]
```

**Purpose:** Fetch changes from Outline and merge to working tree

**Flow:**
1. Check if `.outline/` exists (not initialized? show error)
2. Load manifest
3. Fetch remote collection document tree
4. For each document:
   - Check if local file exists
   - Compare local hash vs manifest hash (detect local changes)
   - Compare manifest revision vs remote revision (detect remote changes)
   - **Conflict detection:**
     - Local unchanged + remote changed → Pull (fast-forward)
     - Local changed + remote unchanged → Skip (already up-to-date)
     - Local changed + remote changed → Conflict (show warning, skip unless --force)
   - Download content if needed
   - Update file with new content + frontmatter
5. Update manifest
6. Show summary (like git pull)

**Output Example:**
```
From https://outline-rbi.jatismobile.com/collection/jns
 * [updated]     document-1.md
 * [updated]     subfolder/nested-doc.md
 ! [conflict]    document-2.md (both modified)

Fast-forward: 2 documents updated
Conflicts: 1 (use --force to overwrite local changes)
```

### 4. Push Command

```bash
outline push [--force] [--dry-run]
```

**Purpose:** Push local changes to Outline

**Flow:**
1. Check if `.outline/` exists
2. Load manifest
3. Scan local files
4. For each tracked file:
   - Compute current hash
   - Compare with manifest hash
   - If unchanged → skip
   - If changed:
     - Fetch remote document info
     - Check remote revision
     - If remote revision ≠ manifest revision → conflict (skip unless --force)
     - Strip frontmatter
     - Upload content via API
     - Update local frontmatter with new revision
     - Update manifest
5. Show summary

**Output Example:**
```
To https://outline-rbi.jatismobile.com/collection/jns
 + [pushed]      document-1.md (rev 5 → 6)
 + [pushed]      document-2.md (rev 3 → 4)
 ! [rejected]    document-3.md (remote changed, pull first)

2 documents pushed
1 rejection (use 'outline pull' first)
```

### 5. Status Command

```bash
outline status
```

**Purpose:** Show working tree status (like `git status`)

**Flow:**
1. Check if `.outline/` exists
2. Load manifest
3. For each tracked file:
   - Check if exists locally
   - Compute hash
   - Compare with manifest
4. Detect untracked files (exist locally but not in manifest)
5. Show categorized output

**Output Example:**
```
On collection: JNS (2e317a13-b7fa-469f-aef8-27474cf336ed)
Your workspace is up to date with 'origin/main'.

Changes not pushed:
  (use "outline push" to upload changes)

    modified:   document-1.md
    modified:   subfolder/nested-doc.md

Untracked files:
  (use "outline add <file>" to track new documents)

    new-document.md

2 files modified, 1 untracked
```

### 6. Add Command (NEW)

```bash
outline add <file>
outline add .  # Add all untracked files
```

**Purpose:** Track new local markdown files (like `git add`)

**Flow:**
1. Check if file exists
2. Parse frontmatter (if exists)
3. If no frontmatter or no outline_id:
   - Create new document via API
   - Get document ID and revision
   - Add frontmatter to file
4. Add to manifest

### 7. Remote Command (NEW)

```bash
outline remote show
outline remote set-url <new-url>
```

**Purpose:** Manage remote Outline instance

### 8. Collections Command

```bash
outline collections list
outline collections show <collection-id>
```

**Purpose:** List and inspect collections

### 9. Search Command

```bash
outline search <query> [--collection <id>]
```

**Purpose:** Search documents across Outline

## Configuration (Git-like Hierarchy)

### Global Config

```
~/.outline/config  # Global settings for all repositories
```

Format:
```ini
[user]
    name = Rayzal Zero
    email = rayzal@example.com

[api]
    token = ${OUTLINE_API_KEY}
    default_url = https://outline-rbi.jatismobile.com
    
[sync]
    conflict_strategy = prompt
    api_delay = 300ms
```

### Repository Config

```
.outline/config  # Per-repository settings (overrides global)
```

Format:
```ini
[remote "origin"]
    url = https://outline-rbi.jatismobile.com
    collection = 2e317a13-b7fa-469f-aef8-27474cf336ed

[auth]
    token = ${OUTLINE_API_KEY}

[sync]
    auto_pull = false
```

### Environment Variables

- `OUTLINE_API_KEY` - API token (highest priority)
- `OUTLINE_BASE_URL` - Override base URL
- `OUTLINE_CONFIG_DIR` - Override global config directory (default: `~/.outline`)

### Priority

CLI flag > Environment variable > Repository config > Global config > Default

## Frontmatter Processing

### Format

```yaml
---
outline_id: "uuid"
outline_collection: "collection-name"
outline_url: "https://outline-rbi.jatismobile.com/doc/..."
outline_updated: "2026-04-01T07:11:54.989Z"
outline_revision: 5
---
```

### Processing

```go
type Frontmatter struct {
    OutlineID         string    `yaml:"outline_id"`
    OutlineCollection string    `yaml:"outline_collection"`
    OutlineURL        string    `yaml:"outline_url"`
    OutlineUpdated    time.Time `yaml:"outline_updated"`
    OutlineRevision   int       `yaml:"outline_revision"`
}

func ParseMarkdown(content []byte) (*Frontmatter, string, error)
func SerializeMarkdown(fm *Frontmatter, content string) []byte
```

## Conflict Detection

### Cases

| Local Hash | Remote Revision | Action |
|------------|----------------|--------|
| Changed | Same | Push (local edit only) |
| Same | Different | Pull (remote edit only) |
| Changed | Different | Conflict (skip unless --force) |
| Same | Same | Skip (up to date) |

## Dependencies

```go
// go.mod
module github.com/rayzalzero/outline-cli

go 1.22

require (
    github.com/spf13/cobra v1.8.0       // CLI framework
    github.com/spf13/viper v1.18.2      // Config management
    gopkg.in/yaml.v3 v3.0.1             // YAML parsing
    github.com/fatih/color v1.16.0      // Colored output
    github.com/schollz/progressbar/v3   // Progress bars
)
```

## Build & Distribution

### Makefile

```makefile
.PHONY: build build-all install clean

build:
	go build -o bin/outline cmd/outline/main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/outline-linux-amd64 cmd/outline/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/outline-darwin-amd64 cmd/outline/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/outline-darwin-arm64 cmd/outline/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/outline-windows-amd64.exe cmd/outline/main.go

install:
	go install cmd/outline/main.go

clean:
	rm -rf bin/
```

## Testing Strategy

1. **Unit Tests** - Each package has _test.go files
2. **Integration Tests** - Mock Outline API responses
3. **E2E Tests** - Test against real Outline instance (optional)

## Migration from Bash Scripts

### Compatibility

- Manifest format: 100% compatible
- Frontmatter format: 100% compatible
- File structure: 100% compatible
- No breaking changes for existing synced repos

### Migration Steps

1. Build Go binary
2. Test with `outline status` (read-only)
3. Test with `outline pull --dry-run`
4. Gradually replace bash scripts

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1)
- [ ] Project setup (go.mod, directory structure)
- [ ] API client (basic HTTP wrapper with retry logic)
- [ ] Config loading (global + repository)
- [ ] Manifest struct and serialization
- [ ] Repository detection (check for `.outline/`)

### Phase 2: Init & Clone (Week 1-2)
- [ ] `outline init` command
- [ ] `outline clone` command (single collection)
- [ ] `.outline/` directory creation
- [ ] Collection metadata fetching
- [ ] Document tree traversal
- [ ] Frontmatter generation

### Phase 3: Status & Diff (Week 2)
- [ ] `outline status` command
- [ ] File hash computation
- [ ] Manifest comparison
- [ ] Untracked file detection
- [ ] Colored output (modified, deleted, untracked)

### Phase 4: Pull (Week 2-3)
- [ ] `outline pull` command
- [ ] Conflict detection (3-way merge logic)
- [ ] Fast-forward updates
- [ ] File writing with directory creation
- [ ] Manifest update after pull
- [ ] Summary output

### Phase 5: Push (Week 3)
- [ ] `outline push` command
- [ ] Frontmatter stripping
- [ ] API update calls
- [ ] Remote revision checking
- [ ] Conflict handling
- [ ] Manifest update after push

### Phase 6: Add & Track (Week 3-4)
- [ ] `outline add` command
- [ ] Create new document via API
- [ ] Generate frontmatter for new files
- [ ] Add to manifest

### Phase 7: Advanced Features (Week 4)
- [ ] `outline clone --all` (all collections)
- [ ] `outline remote` command
- [ ] `outline collections` command
- [ ] `outline search` command
- [ ] Progress bars for long operations

### Phase 8: Polish & Release (Week 4-5)
- [ ] Error handling improvements
- [ ] Better help messages
- [ ] Build scripts (Makefile)
- [ ] Cross-platform binaries
- [ ] GitHub Actions for releases
- [ ] Documentation (README, USAGE)
- [ ] Migration guide from bash scripts

## Success Criteria

- [ ] Can init new repository
- [ ] Can clone collection from Outline
- [ ] Can pull updates (fast-forward and merge)
- [ ] Can push local changes to Outline
- [ ] Detects conflicts correctly (3-way merge)
- [ ] Can track new files with `outline add`
- [ ] Status shows clear working tree state
- [ ] Preserves frontmatter on round-trip
- [ ] Works on Linux, macOS, Windows
- [ ] Single binary, no dependencies
- [ ] Faster than bash scripts (parallel downloads)
- [ ] Git-like UX (familiar commands and output)
- [ ] Better error messages than bash scripts

## Workflow Examples

### Scenario 1: Clone and Start Working

```bash
# Clone JNS documentation
outline clone https://outline-rbi.jatismobile.com/collection/jns jns-docs
cd jns-docs

# Check status
outline status
# On collection: JNS Documentation
# Your workspace is up to date with 'origin'.

# Edit a file
vim document-1.md

# Check what changed
outline status
# Changes not pushed:
#   modified:   document-1.md

# Push changes
outline push
# To https://outline-rbi.jatismobile.com/collection/jns
#  + [pushed]      document-1.md (rev 5 → 6)
```

### Scenario 2: Pull Remote Changes

```bash
# Someone else updated docs on Outline web UI
outline pull
# From https://outline-rbi.jatismobile.com/collection/jns
#  * [updated]     document-2.md
#  * [updated]     subfolder/nested-doc.md
# 
# Fast-forward: 2 documents updated
```

### Scenario 3: Handle Conflicts

```bash
# You edited document-1.md locally
vim document-1.md

# Someone else also edited it on web
outline pull
# From https://outline-rbi.jatismobile.com/collection/jns
#  * [updated]     document-2.md
#  ! [conflict]    document-1.md (both modified)
# 
# 1 document updated, 1 conflict
# Use 'outline pull --force' to overwrite local changes
# Or push your changes with 'outline push --force'

# Manually resolve: choose to keep local
outline push --force
```

### Scenario 4: Add New Document

```bash
# Create new document locally
echo "# New Feature" > new-feature.md

# Check status
outline status
# Untracked files:
#   new-feature.md

# Track it (creates document on Outline)
outline add new-feature.md
# Created new document: new-feature.md (id: abc123...)
# Added frontmatter to file

# Now it's tracked
outline status
# Changes not pushed:
#   new:        new-feature.md

# Push to sync
outline push
```

### Scenario 5: Work with Multiple Collections

```bash
# Clone multiple collections
outline clone https://outline-rbi.jatismobile.com/collection/jns jns-docs
outline clone https://outline-rbi.jatismobile.com/collection/releases releases

# Work in each independently
cd jns-docs
outline status

cd ../releases
outline status
```
