---
outline_id: de287ac3-95a6-471f-b84a-0aff39fa4321
outline_collection: Testing
outline_url: /doc/configuration-files-reference-VG89w9SaB8
outline_updated: 2026-07-03T04:39:43.533Z
outline_revision: 4
outline_parent_id: 93006387-0fed-4138-9725-81d5aa1c10eb
---

# Configuration Files Reference

Complete reference for Testing Collection configuration files.

## Table of Contents

* [Config File Format](./config-format.md)
* [Environment Variables](./environment-variables.md)
* [File Locations](./file-locations.md)

## Configuration Files

Testing Collection uses several configuration files:

### Main Configuration

* `testing-collection.config.js` - JavaScript config
* `testing-collection.config.ts` - TypeScript config
* `testing-collection.config.json` - JSON config

### Environment Files

* `.env` - Environment variables
* `.env.test` - Test environment
* `.env.local` - Local overrides

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

* [Config File Format](./config-format.md)
* [Environment Variables](./environment-variables.md)
* [File Locations](./file-locations.md)