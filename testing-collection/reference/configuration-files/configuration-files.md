---
outline_id: ca3d39c1-7515-439c-afaf-8fdfe3d3cc74
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/configuration-files-reference-cX0hmlDRrw
outline_updated: 2026-07-02T10:56:29.772Z
outline_revision: 2
---



# Configuration Files Reference

Complete reference for Testing Collection configuration files.

## Table of Contents

- [Config File Format](./config-format.md)
- [Environment Variables](./environment-variables.md)
- [File Locations](./file-locations.md)

## Configuration Files

Testing Collection uses several configuration files:

### Main Configuration

- `testing-collection.config.js` - JavaScript config
- `testing-collection.config.ts` - TypeScript config
- `testing-collection.config.json` - JSON config

### Environment Files

- `.env` - Environment variables
- `.env.test` - Test environment
- `.env.local` - Local overrides

## File Priority

Configuration is loaded in this order (later overrides earlier):

1. Default configuration
2. Global configuration (`~/.testing-collection/config`)
3. Project configuration (`testing-collection.config.js`)
4. Environment variables
5. Command-line arguments

## Quick Reference

### Basic Config

```javascript
module.exports = {
  testDir: './tests',
  timeout: 30000,
  retries: 2,
  workers: 4,
  reporter: 'html'
};
```

### Environment Variables

```bash
TESTING_COLLECTION_API_KEY=tc_live_xxx
NODE_ENV=test
TIMEOUT=60000
```

## Next Steps

- [Config File Format](./config-format.md)
- [Environment Variables](./environment-variables.md)
- [File Locations](./file-locations.md)

