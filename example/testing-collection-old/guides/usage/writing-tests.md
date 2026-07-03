---
outline_id: e7c93f8f-f44a-4f77-aebc-99b7b2c34309
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/writing-tests-Ui5fXKTa8T
outline_updated: 2026-07-03T04:39:44.669Z
outline_revision: 2
---

# Writing Tests

Guide to writing effective tests with Testing Collection.

## Test Structure

### Basic Test

```javascript
test('should add two numbers', () => {
  const result = add(2, 3);
  expect(result).toBe(5);
});
```

### Test Suites

```javascript
describe('Calculator', () => {
  describe('add', () => {
    test('should add positive numbers', () => {
      expect(add(2, 3)).toBe(5);
    });

    test('should add negative numbers', () => {
      expect(add(-2, -3)).toBe(-5);
    });
  });

  describe('subtract', () => {
    test('should subtract numbers', () => {
      expect(subtract(5, 3)).toBe(2);
    });
  });
});
```

## Setup and Teardown

### Before/After Hooks

```javascript
describe('Database Tests', () => {
  beforeAll(async () => {
    // Run once before all tests
    await database.connect();
  });

  beforeEach(async () => {
    // Run before each test
    await database.clear();
  });

  afterEach(async () => {
    // Run after each test
    await database.cleanup();
  });

  afterAll(async () => {
    // Run once after all tests
    await database.disconnect();
  });

  test('should insert record', async () => {
    const result = await database.insert({ name: 'John' });
    expect(result).toBeDefined();
  });
});
```

## Assertions

### Basic Matchers

```javascript
// Equality
expect(value).toBe(5);
expect(value).toEqual({ name: 'John' });
expect(value).not.toBe(null);

// Truthiness
expect(value).toBeTruthy();
expect(value).toBeFalsy();
expect(value).toBeNull();
expect(value).toBeUndefined();
expect(value).toBeDefined();
```

### Number Matchers

```javascript
expect(value).toBeGreaterThan(3);
expect(value).toBeGreaterThanOrEqual(3);
expect(value).toBeLessThan(10);
expect(value).toBeLessThanOrEqual(10);
expect(value).toBeCloseTo(0.3, 5);
```

### String Matchers

```javascript
expect(string).toMatch(/pattern/);
expect(string).toContain('substring');
expect(string).toHaveLength(10);
```

### Array Matchers

```javascript
expect(array).toContain(item);
expect(array).toHaveLength(3);
expect(array).toEqual(expect.arrayContaining([1, 2]));
```

### Object Matchers

```javascript
expect(object).toHaveProperty('key');
expect(object).toHaveProperty('key', 'value');
expect(object).toMatchObject({ key: 'value' });
expect(object).toEqual(expect.objectContaining({ key: 'value' }));
```

## Async Testing

### Async/Await

```javascript
test('async operation', async () => {
  const data = await fetchData();
  expect(data).toBe('result');
});
```

### Promises

```javascript
test('promise operation', () => {
  return fetchData().then(data => {
    expect(data).toBe('result');
  });
});
```

### Callbacks

```javascript
test('callback operation', done => {
  fetchData((error, data) => {
    expect(data).toBe('result');
    done();
  });
});
```

## Mocking

### Mock Functions

```javascript
const mockFn = jest.fn();
mockFn.mockReturnValue(42);

test('mock function', () => {
  const result = mockFn();
  expect(result).toBe(42);
  expect(mockFn).toHaveBeenCalled();
});
```

### Mock Modules

```javascript
jest.mock('./api');
const api = require('./api');

test('mock module', async () => {
  api.fetchUser.mockResolvedValue({ name: 'John' });
  const user = await api.fetchUser(1);
  expect(user.name).toBe('John');
});
```

### Spy on Methods

```javascript
const spy = jest.spyOn(object, 'method');

test('spy on method', () => {
  object.method();
  expect(spy).toHaveBeenCalled();
  spy.mockRestore();
});
```

## Snapshot Testing

```javascript
test('component snapshot', () => {
  const component = render(<Button>Click me</Button>);
  expect(component).toMatchSnapshot();
});
```

Update snapshots:

```bash
testing-collection run --update-snapshots
```

## Parameterized Tests

```javascript
test.each([
  [1, 1, 2],
  [2, 2, 4],
  [3, 3, 6]
])('add(%i, %i) should equal %i', (a, b, expected) => {
  expect(add(a, b)).toBe(expected);
});
```

## Test Skipping

```javascript
// Skip single test
test.skip('not ready yet', () => {
  // test code
});

// Skip suite
describe.skip('Feature', () => {
  // tests
});

// Only run specific test
test.only('focus on this', () => {
  // test code
});
```

## Next Steps

- [Running Tests](./running-tests.md)
- [Best Practices](./best-practices.md)
- [API Documentation](../../api/README.md)

