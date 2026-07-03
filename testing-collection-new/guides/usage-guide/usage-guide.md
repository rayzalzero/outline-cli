---
outline_id: 370a099c-1d4f-444b-ae03-84fdd06298d9
outline_collection: Testing
outline_url: /doc/usage-guide-mnrVhAaf1D
outline_updated: 2026-07-03T04:39:43.388Z
outline_revision: 4
outline_parent_id: 0cbee575-9a5a-4cef-943f-dddb5d9acce3
---

# Usage Guide

Learn how to use Testing Collection effectively in your projects.

## Table of Contents

* [Running Tests](./running-tests.md)
* [Writing Tests](./writing-tests.md)
* [Best Practices](./best-practices.md)

## Quick Start

### Run All Tests

```bash
testing-collection run
```

### Run Specific Tests

```bash
# By file pattern
testing-collection run --filter "auth"

# By test name
testing-collection run --grep "login"

# Specific file
testing-collection run tests/auth.test.js
```

## Common Commands

### Watch Mode

```bash
testing-collection run --watch
```

### Debug Mode

```bash
testing-collection run --debug
```

### Generate Reports

```bash
testing-collection run --report html
```

### Coverage

```bash
testing-collection run --coverage
```

## Test Organization

### Directory Structure

```
tests/
├── unit/
│   ├── auth.test.js
│   └── utils.test.js
├── integration/
│   ├── api.test.js
│   └── database.test.js
└── e2e/
    └── user-flow.test.js
```

### Naming Conventions

* Unit tests: `*.test.js` or `*.spec.js`
* Integration tests: `*.integration.test.js`
* E2E tests: `*.e2e.test.js`

## Basic Test Structure

```javascript
describe('Feature Name', () => {
  beforeAll(() => {
    // Setup before all tests
  });

  beforeEach(() => {
    // Setup before each test
  });

  test('should do something', () => {
    // Test code
    expect(result).toBe(expected);
  });

  afterEach(() => {
    // Cleanup after each test
  });

  afterAll(() => {
    // Cleanup after all tests
  });
});
```

## Assertions

```javascript
// Equality
expect(value).toBe(expected);
expect(value).toEqual(expected);

// Truthiness
expect(value).toBeTruthy();
expect(value).toBeFalsy();

// Numbers
expect(value).toBeGreaterThan(3);
expect(value).toBeLessThan(10);

// Strings
expect(string).toMatch(/pattern/);
expect(string).toContain('substring');

// Arrays
expect(array).toContain(item);
expect(array).toHaveLength(3);

// Objects
expect(object).toHaveProperty('key');
expect(object).toMatchObject({ key: 'value' });
```

## Async Testing

```javascript
// Using async/await
test('async operation', async () => {
  const result = await fetchData();
  expect(result).toBe('data');
});

// Using promises
test('promise operation', () => {
  return fetchData().then(result => {
    expect(result).toBe('data');
  });
});
```

## Next Steps

* [Running Tests](./running-tests.md) - Detailed execution options
* [Writing Tests](./writing-tests.md) - Test authoring guide
* [Best Practices](./best-practices.md) - Tips and patterns