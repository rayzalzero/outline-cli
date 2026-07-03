---
outline_id: 8811561b-5d7b-4b9f-814b-34b5e9c52d64
outline_collection: Testing
outline_url: /doc/configuration-guide-QlrMnvcD6M
outline_updated: 2026-07-03T04:39:43.222Z
outline_revision: 4
outline_parent_id: 0cbee575-9a5a-4cef-943f-dddb5d9acce3
---

# Configuration Guide

Learn how to configure Testing Collection for your project.

## Table of Contents

* [Basic Configuration](./basic-configuration.md)
* [Advanced Options](./advanced-options.md)
* [Common Issues](./common-issues.md)

## Configuration File

Testing Collection uses a configuration file in your project root:

```javascript
// testing-collection.config.js
module.exports = {
  testDir: './tests',
  timeout: 30000,
  retries: 2,
  reporter: 'html',
  use: {
    headless: true,
    viewport: { width: 1280, height: 720 }
  }
};
```

## Quick Configuration

### Initialize Config

```bash
testing-collection init --config
```

### Common Settings

```javascript
module.exports = {
  // Test directory
  testDir: './tests',
  
  // Timeout per test
  timeout: 30000,
  
  // Retry failed tests
  retries: 2,
  
  // Reporter type
  reporter: 'html',
  
  // Parallel execution
  workers: 4
};
```

## Environment Variables

Set environment-specific configuration:

```bash
# .env
TEST_ENV=staging
API_URL=https://api.staging.example.com
TIMEOUT=60000
```

Load in config:

```javascript
require('dotenv').config();

module.exports = {
  timeout: process.env.TIMEOUT || 30000,
  use: {
    baseURL: process.env.API_URL
  }
};
```

## Configuration Formats

### JavaScript

```javascript
// testing-collection.config.js
module.exports = { /* config */ };
```

### TypeScript

```typescript
// testing-collection.config.ts
import { Config } from 'testing-collection';

const config: Config = { /* config */ };
export default config;
```

### JSON

```json
{
  "testDir": "./tests",
  "timeout": 30000
}
```

## Next Steps

* [Basic Configuration](./basic-configuration.md) - Essential settings
* [Advanced Options](./advanced-options.md) - Fine-tune behavior
* [Common Issues](./common-issues.md) - Troubleshoot config problems