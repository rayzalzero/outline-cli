---
outline_id: 806856a0-4fff-4c81-86bf-28ad8c56127d
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/configuration-file-locations-oUMuOEPdkx
outline_updated: 2026-07-03T04:39:45.166Z
outline_revision: 2
---

# Configuration File Locations

Where Testing Collection looks for configuration files.

## Search Order

Testing Collection searches for configuration in this order:

1. Command-line specified config (`--config`)
2. Project root config files
3. Global user config
4. Default configuration

## Project Configuration

### Root Directory

Testing Collection looks for these files in the project root:

```
testing-collection.config.js
testing-collection.config.ts
testing-collection.config.json
testing-collection.config.cjs
testing-collection.config.mjs
```

### Example Structure

```
my-project/
├── testing-collection.config.js  ← Project config
├── .env                           ← Environment variables
├── tests/
└── package.json
```

## Global Configuration

### User Home Directory

```
~/.testing-collection/
├── config.js                      ← Global config
└── credentials                    ← API keys
```

### Platform-Specific Paths

**macOS/Linux:**
```
~/.testing-collection/config.js
```

**Windows:**
```
C:\Users\<username>\.testing-collection\config.js
```

## Environment Files

### Location

Environment files are loaded from the project root:

```
.env                    ← Default
.env.local              ← Local overrides (not committed)
.env.test               ← Test environment
.env.production         ← Production environment
```

### Loading Priority

1. `.env.local` (highest priority)
2. `.env.{NODE_ENV}`
3. `.env`

## Custom Locations

### Specify Config Path

```bash
# Via command line
testing-collection run --config ./config/test.config.js

# Via environment variable
export TESTING_COLLECTION_CONFIG=./config/test.config.js
testing-collection run
```

### Relative vs Absolute Paths

```bash
# Relative to current directory
--config ./config/test.config.js

# Absolute path
--config /Users/username/project/config/test.config.js
```

## Monorepo Configuration

### Workspace Structure

```
monorepo/
├── packages/
│   ├── app-a/
│   │   └── testing-collection.config.js
│   └── app-b/
│       └── testing-collection.config.js
└── testing-collection.config.js  ← Root config
```

### Shared Configuration

```javascript
// root/testing-collection.config.js
module.exports = {
  timeout: 30000,
  workers: 4
};

// packages/app-a/testing-collection.config.js
const baseConfig = require('../../testing-collection.config.js');

module.exports = {
  ...baseConfig,
  testDir: './tests'
};
```

## Cache and Output Directories

### Default Locations

```
my-project/
├── .testing-collection/           ← Cache directory
├── coverage/                      ← Coverage reports
├── test-results/                  ← Test reports
└── node_modules/.cache/           ← Build cache
```

### Custom Locations

```javascript
module.exports = {
  cacheDirectory: '.cache/testing-collection',
  coverageDirectory: 'reports/coverage',
  reportDirectory: 'reports/tests'
};
```

## Troubleshooting

### Config Not Found

```bash
# Verify config file exists
ls testing-collection.config.js

# Check search path
testing-collection config show --verbose
```

### Wrong Config Loaded

```bash
# Show which config is being used
testing-collection run --verbose

# Specify config explicitly
testing-collection run --config ./testing-collection.config.js
```

## Next Steps

- [Config File Format](./config-format.md)
- [Environment Variables](./environment-variables.md)
- [Configuration Guide](../../guides/configuration/README.md)

