---
outline_id: 0a071fb2-4bd1-4310-bbaf-da063a0abf1e
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/running-tests-66opty3H0n
outline_updated: 2026-07-03T04:39:44.572Z
outline_revision: 2
---

# Running Tests

Comprehensive guide to running tests with Testing Collection.

## Basic Execution

### Run All Tests

```bash
testing-collection run
```

### Run Specific Directory

```bash
testing-collection run tests/unit
```

### Run Specific File

```bash
testing-collection run tests/auth.test.js
```

## Filtering Tests

### By Pattern

```bash
# Match file names
testing-collection run --filter "auth"

# Multiple patterns
testing-collection run --filter "auth|user"
```

### By Test Name

```bash
# Match test descriptions
testing-collection run --grep "login"

# Exclude tests
testing-collection run --grep "login" --invert
```

### By Tag

```bash
# Run tagged tests
testing-collection run --tag @smoke

# Multiple tags
testing-collection run --tag "@smoke,@critical"
```

## Execution Modes

### Watch Mode

Automatically re-run tests on file changes:

```bash
testing-collection run --watch
```

Watch specific files:

```bash
testing-collection run --watch --filter "auth"
```

### Parallel Execution

```bash
# Use 4 workers
testing-collection run --workers 4

# Use 50% of CPU cores
testing-collection run --workers 50%
```

### Sequential Execution

```bash
testing-collection run --workers 1
```

## Debugging

### Debug Mode

```bash
testing-collection run --debug
```

### Verbose Output

```bash
testing-collection run --verbose
```

### Inspect Mode

```bash
testing-collection run --inspect
```

Attach debugger:

```bash
testing-collection run --inspect-brk
node inspect localhost:9229
```

## Retry Configuration

### Retry Failed Tests

```bash
# Retry up to 3 times
testing-collection run --retries 3
```

### Retry Only on CI

```bash
testing-collection run --retries-on-ci 2
```

## Timeout Options

### Global Timeout

```bash
testing-collection run --timeout 60000
```

### Per-Test Timeout

```javascript
test('slow operation', async () => {
  test.setTimeout(120000);
  // test code
});
```

## Reporting

### HTML Report

```bash
testing-collection run --report html
```

### JSON Report

```bash
testing-collection run --report json --output results.json
```

### Multiple Reports

```bash
testing-collection run --report html,json,junit
```

### Open Report

```bash
testing-collection run --report html --open
```

## Coverage

### Collect Coverage

```bash
testing-collection run --coverage
```

### Coverage Reporters

```bash
testing-collection run --coverage --coverage-reporter html,text,lcov
```

### Coverage Threshold

```bash
testing-collection run --coverage --coverage-threshold 80
```

## Environment Variables

### Set Environment

```bash
NODE_ENV=test testing-collection run
```

### Load from .env

```bash
testing-collection run --env-file .env.test
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run Tests
  run: testing-collection run --ci --coverage
```

### GitLab CI

```yaml
test:
  script:
    - testing-collection run --ci --report junit
```

### Jenkins

```groovy
sh 'testing-collection run --ci --report html'
```

## Advanced Options

### Bail on Failure

```bash
# Stop after first failure
testing-collection run --bail

# Stop after N failures
testing-collection run --bail 3
```

### Update Snapshots

```bash
testing-collection run --update-snapshots
```

### Clear Cache

```bash
testing-collection run --clear-cache
```

## Performance Tips

```bash
# Optimize for speed
testing-collection run \
  --workers 50% \
  --cache \
  --no-coverage \
  --bail
```

## Next Steps

- [Writing Tests](./writing-tests.md)
- [Best Practices](./best-practices.md)
- [Configuration](../configuration/README.md)

