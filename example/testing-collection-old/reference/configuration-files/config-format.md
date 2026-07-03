---
outline_id: c6868afa-dc3a-4d81-9a07-4d6d78a5f76f
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/configuration-file-format-5kX1T2fixx
outline_updated: 2026-07-03T04:39:45.006Z
outline_revision: 2
---

# Configuration File Format

Detailed specification for Testing Collection configuration files.

## File Types

### JavaScript (Recommended)

```javascript
// testing-collection.config.js
module.exports = {
  testDir: './tests',
  timeout: 30000,
  retries: 2,
  workers: 4,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000'
  }
};
```

### TypeScript

```typescript
// testing-collection.config.ts
import { Config } from 'testing-collection';

const config: Config = {
  testDir: './tests',
  timeout: 30000,
  retries: 2,
  workers: 4,
  reporter: 'html'
};

export default config;
```

### JSON

```json
{
  "testDir": "./tests",
  "timeout": 30000,
  "retries": 2,
  "workers": 4,
  "reporter": "html"
}
```

## Configuration Schema

### Root Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| testDir | string | './tests' | Test directory |
| timeout | number | 30000 | Global timeout (ms) |
| retries | number | 0 | Retry count |
| workers | number | 4 | Parallel workers |
| reporter | string | 'html' | Reporter type |

### Test Matching

```javascript
{
  testMatch: '**/*.test.js',
  testIgnore: ['**/node_modules/**']
}
```

### Coverage

```javascript
{
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageReporters: ['html', 'text', 'lcov'],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  }
}
```

### Reporters

```javascript
{
  reporter: [
    ['html', { outputFolder: 'test-results' }],
    ['json', { outputFile: 'results.json' }],
    ['junit', { outputFile: 'junit.xml' }]
  ]
}
```

## Complete Example

```javascript
module.exports = {
  // Test execution
  testDir: './tests',
  testMatch: '**/*.test.js',
  testIgnore: ['**/node_modules/**', '**/dist/**'],
  
  // Timeouts and retries
  timeout: 30000,
  retries: process.env.CI ? 2 : 0,
  
  // Parallel execution
  workers: process.env.CI ? 2 : 4,
  maxConcurrency: 10,
  
  // Reporting
  reporter: [
    ['html', { outputFolder: 'test-results' }],
    ['json', { outputFile: 'results.json' }]
  ],
  
  // Coverage
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageReporters: ['html', 'text', 'lcov'],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  },
  
  // Setup/teardown
  globalSetup: './global-setup.js',
  globalTeardown: './global-teardown.js',
  
  // Environment
  testEnvironment: 'node',
  
  // Custom options
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    headless: true
  }
};
```

## Next Steps

- [Environment Variables](./environment-variables.md)
- [File Locations](./file-locations.md)
- [Configuration Guide](../../guides/configuration/README.md)

