---
outline_id: f2b0408a-1a93-4a01-b991-27a15ca8534f
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/basic-configuration-vSN1DQRrtb
outline_updated: 2026-07-02T10:18:17.683Z
outline_revision: 2
---



# Basic Configuration

Essential configuration options for Testing Collection.

## Test Directory

Specify where your tests are located:

```javascript
module.exports = {
  testDir: './tests'
};
```

Multiple directories:

```javascript
module.exports = {
  testDir: ['./tests', './integration-tests']
};
```

## Timeout Settings

### Global Timeout

```javascript
module.exports = {
  timeout: 30000  // 30 seconds
};
```

### Per-Test Timeout

```javascript
test('slow operation', async () => {
  test.setTimeout(60000);  // 60 seconds
  // test code
});
```

## Retry Configuration

Automatically retry failed tests:

```javascript
module.exports = {
  retries: 2  // Retry up to 2 times
};
```

Conditional retries:

```javascript
module.exports = {
  retries: process.env.CI ? 2 : 0
};
```

## Reporter Options

### Single Reporter

```javascript
module.exports = {
  reporter: 'html'
};
```

### Multiple Reporters

```javascript
module.exports = {
  reporter: [
    ['html', { outputFolder: 'test-results' }],
    ['json', { outputFile: 'results.json' }],
    ['junit', { outputFile: 'junit.xml' }]
  ]
};
```

## Parallel Execution

### Worker Configuration

```javascript
module.exports = {
  workers: 4  // Run 4 tests in parallel
};
```

Auto-detect CPU cores:

```javascript
module.exports = {
  workers: process.env.CI ? 2 : undefined  // Auto on local
};
```

## Test Matching

### Include Patterns

```javascript
module.exports = {
  testMatch: '**/*.test.js'
};
```

Multiple patterns:

```javascript
module.exports = {
  testMatch: [
    '**/*.test.js',
    '**/*.spec.js'
  ]
};
```

### Ignore Patterns

```javascript
module.exports = {
  testIgnore: [
    '**/node_modules/**',
    '**/dist/**',
    '**/*.skip.js'
  ]
};
```

## Global Setup

Run code before all tests:

```javascript
module.exports = {
  globalSetup: './global-setup.js'
};
```

```javascript
// global-setup.js
module.exports = async () => {
  console.log('Starting test suite...');
  // Setup database, start servers, etc.
};
```

## Global Teardown

Run code after all tests:

```javascript
module.exports = {
  globalTeardown: './global-teardown.js'
};
```

```javascript
// global-teardown.js
module.exports = async () => {
  console.log('Cleaning up...');
  // Close connections, cleanup resources
};
```

## Example Complete Config

```javascript
module.exports = {
  testDir: './tests',
  timeout: 30000,
  retries: 2,
  workers: 4,
  reporter: 'html',
  testMatch: '**/*.test.js',
  testIgnore: ['**/node_modules/**'],
  globalSetup: './setup.js',
  globalTeardown: './teardown.js'
};
```

## Next Steps

- [Advanced Options](./advanced-options.md)
- [Common Issues](./common-issues.md)
- [Usage Guide](../usage/README.md)

