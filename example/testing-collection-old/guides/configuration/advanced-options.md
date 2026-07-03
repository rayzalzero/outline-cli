---
outline_id: 031c3595-ce74-496b-a302-b69c348b5047
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/advanced-configuration-options-A3CANbRiYQ
outline_updated: 2026-07-03T04:39:44.037Z
outline_revision: 2
---

# Advanced Configuration Options

Fine-tune Testing Collection behavior with advanced settings.

## Custom Test Environment

Define custom test environment:

```javascript
module.exports = {
  testEnvironment: 'node',  // or 'jsdom' for browser
  testEnvironmentOptions: {
    url: 'http://localhost',
    userAgent: 'Custom Agent'
  }
};
```

## Module Resolution

### Path Mapping

```javascript
module.exports = {
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
    '^@components/(.*)$': '<rootDir>/src/components/$1'
  }
};
```

### Transform Rules

```javascript
module.exports = {
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
    '^.+\\.jsx?$': 'babel-jest'
  }
};
```

## Coverage Configuration

### Basic Coverage

```javascript
module.exports = {
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageReporters: ['html', 'text', 'lcov']
};
```

### Coverage Thresholds

```javascript
module.exports = {
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  }
};
```

### Coverage Paths

```javascript
module.exports = {
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.test.{js,jsx,ts,tsx}',
    '!src/**/*.d.ts'
  ]
};
```

## Watch Mode

Configure file watching:

```javascript
module.exports = {
  watchPathIgnorePatterns: [
    '/node_modules/',
    '/dist/',
    '/.git/'
  ],
  watchPlugins: [
    'jest-watch-typeahead/filename',
    'jest-watch-typeahead/testname'
  ]
};
```

## Snapshot Testing

```javascript
module.exports = {
  snapshotSerializers: ['enzyme-to-json/serializer'],
  snapshotResolver: './snapshot-resolver.js'
};
```

## Custom Matchers

Add custom assertion matchers:

```javascript
module.exports = {
  setupFilesAfterEnv: ['./jest.setup.js']
};
```

```javascript
// jest.setup.js
expect.extend({
  toBeWithinRange(received, floor, ceiling) {
    const pass = received >= floor && received <= ceiling;
    return {
      pass,
      message: () => `expected ${received} to be within ${floor}-${ceiling}`
    };
  }
});
```

## Performance Optimization

### Cache Configuration

```javascript
module.exports = {
  cache: true,
  cacheDirectory: '.cache/testing-collection'
};
```

### Max Workers

```javascript
module.exports = {
  maxWorkers: '50%',  // Use 50% of CPU cores
  maxConcurrency: 10   // Max concurrent tests per worker
};
```

## Debugging Options

### Verbose Output

```javascript
module.exports = {
  verbose: true,
  silent: false
};
```

### Bail Configuration

Stop after N failures:

```javascript
module.exports = {
  bail: 1  // Stop after first failure
};
```

## Custom Reporters

### Built-in Reporters

```javascript
module.exports = {
  reporters: [
    'default',
    ['html', { outputFolder: 'reports' }],
    ['json', { outputFile: 'results.json' }],
    ['junit', { outputFile: 'junit.xml' }]
  ]
};
```

### Custom Reporter

```javascript
module.exports = {
  reporters: [
    'default',
    ['./custom-reporter.js', { option: 'value' }]
  ]
};
```

## Next Steps

- [Basic Configuration](./basic-configuration.md)
- [Common Issues](./common-issues.md)
- [Usage Guide](../usage/README.md)

