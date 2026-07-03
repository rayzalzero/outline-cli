---
outline_id: 93006387-0fed-4138-9725-81d5aa1c10eb
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/reference-documentation-TYfOVeR3JP
outline_updated: 2026-07-03T04:39:42.714Z
outline_revision: 2
---

# Reference Documentation

Technical reference materials for Testing Collection.

## Table of Contents

- [CLI Commands](./cli-commands/README.md)
- [Configuration Files](./configuration-files/README.md)

## Overview

This section provides detailed technical reference documentation including:

- Complete CLI command reference
- Configuration file specifications
- Environment variables
- Exit codes
- File formats

## Quick Links

### CLI Commands

- [Run Command](./cli-commands/run.md)
- [Init Command](./cli-commands/init.md)
- [Config Command](./cli-commands/config.md)

### Configuration Files

- [Config File Format](./configuration-files/config-format.md)
- [Environment Variables](./configuration-files/environment-variables.md)
- [File Locations](./configuration-files/file-locations.md)

## Conventions

### Command Syntax

```
testing-collection <command> [options] [arguments]
```

### Option Formats

- Short: `-v`
- Long: `--verbose`
- With value: `--timeout 30000`

### Configuration Syntax

```javascript
module.exports = {
  key: 'value'
};
```

## Next Steps

- [CLI Commands](./cli-commands/README.md)
- [Configuration Files](./configuration-files/README.md)
- [Guides](../guides/README.md)


