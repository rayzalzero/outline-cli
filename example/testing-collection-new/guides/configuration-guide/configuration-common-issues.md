---
outline_id: 22b72067-0c5e-4794-ada5-67d0a2060f57
outline_collection: Testing
outline_url: /doc/configuration-common-issues-LQfwkce4UL
outline_updated: 2026-07-03T04:39:44.247Z
outline_revision: 4
outline_parent_id: 8811561b-5d7b-4b9f-814b-34b5e9c52d64
---

# Configuration Common Issues

Troubleshoot common configuration problems.

## Configuration Not Loading

### Issue: Config file ignored

**Symptoms**:

* Tests run with default settings
* Custom config not applied

**Solutions**:

```bash
# Verify config file name
ls testing-collection.config.js

# Specify config explicitly
testing-collection run --config=./my-config.js

# Check config syntax
node -c testing-collection.config.js
```

## Path Resolution Issues

### Issue: Module not found

**Symptoms**:

```
Error: Cannot find module '@/components/Button'
```

**Solution**: Fix path mapping

```javascript
module.exports = {
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1'
  }
};
```

## Timeout Problems

### Issue: Tests timing out

**Symptoms**:

```
Error: Timeout - Async callback was not invoked within 30000ms
```

**Solutions**:

```javascript
// Increase global timeout
module.exports = {
  timeout: 60000
};

// Or per-test
test('slow test', async () => {
  test.setTimeout(120000);
});
```

## Worker Configuration Issues

### Issue: Too many workers

**Symptoms**:

* System slowdown
* Out of memory errors

**Solution**: Limit workers

```javascript
module.exports = {
  workers: 2,  // Reduce workers
  maxConcurrency: 5
};
```

## Coverage Issues

### Issue: Coverage not collected

**Symptoms**:

* No coverage report generated
* Coverage shows 0%

**Solutions**:

```javascript
module.exports = {
  collectCoverage: true,
  collectCoverageFrom: [
    'src/**/*.js',
    '!src/**/*.test.js'
  ]
};
```

## Reporter Problems

### Issue: Reports not generated

**Symptoms**:

* Missing HTML/JSON reports
* Empty report files

**Solutions**:

```javascript
module.exports = {
  reporter: [
    ['html', { 
      outputFolder: 'test-results',
      open: 'never'
    }]
  ]
};
```

Check permissions:

```bash
mkdir -p test-results
chmod 755 test-results
```

## Environment Variable Issues

### Issue: Environment variables not loaded

**Symptoms**:

* `process.env.VAR` is undefined
* Config values missing

**Solution**: Load dotenv

```javascript
require('dotenv').config();

module.exports = {
  timeout: process.env.TIMEOUT || 30000
};
```

## TypeScript Configuration

### Issue: TypeScript files not recognized

**Symptoms**:

```
Error: Cannot use import statement outside a module
```

**Solution**: Configure transform

```javascript
module.exports = {
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx']
};
```

## Watch Mode Issues

### Issue: Watch not detecting changes

**Symptoms**:

* Tests don't re-run on file save
* Watch mode stuck

**Solutions**:

```javascript
module.exports = {
  watchPathIgnorePatterns: [
    '/node_modules/',
    '/dist/'
  ]
};
```

Or restart watch:

```bash
# Kill and restart
testing-collection run --watch
```

## Validation Errors

### Issue: Invalid configuration

**Symptoms**:

```
Error: Invalid configuration object
```

**Solution**: Validate config

```bash
# Check syntax
node -c testing-collection.config.js

# Use schema validation
testing-collection validate-config
```

## Performance Issues

### Issue: Slow test execution

**Solutions**:

```javascript
module.exports = {
  // Enable caching
  cache: true,
  
  // Optimize workers
  workers: '50%',
  
  // Reduce retries
  retries: 0,
  
  // Disable coverage in dev
  collectCoverage: process.env.CI === 'true'
};
```

## Next Steps

* [Basic Configuration](./basic-configuration.md)
* [Advanced Options](./advanced-options.md)
* [Installation Troubleshooting](../installation/troubleshooting.md)