---
outline_id: c79d235a-dc50-43ab-b5bf-7493a798acc9
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/testing-best-practices-NTf0qEG3jD
outline_updated: 2026-07-02T10:18:18.329Z
outline_revision: 2
---



# Testing Best Practices

Guidelines and patterns for writing maintainable, effective tests.

## Test Organization

### Structure by Feature

```
tests/
├── auth/
│   ├── login.test.js
│   ├── logout.test.js
│   └── registration.test.js
├── user/
│   ├── profile.test.js
│   └── settings.test.js
└── payment/
    └── checkout.test.js
```

### Naming Conventions

```javascript
// Good: Descriptive test names
test('should return user when valid ID is provided', () => {});
test('should throw error when user not found', () => {});

// Bad: Vague test names
test('works', () => {});
test('test1', () => {});
```

## AAA Pattern

Arrange, Act, Assert:

```javascript
test('should calculate total price', () => {
  // Arrange
  const cart = new Cart();
  cart.addItem({ price: 10, quantity: 2 });
  
  // Act
  const total = cart.getTotal();
  
  // Assert
  expect(total).toBe(20);
});
```

## Test Independence

### Avoid Test Dependencies

```javascript
// Bad: Tests depend on execution order
test('create user', () => {
  user = createUser();
});

test('update user', () => {
  updateUser(user); // Depends on previous test
});

// Good: Each test is independent
test('create user', () => {
  const user = createUser();
  expect(user).toBeDefined();
});

test('update user', () => {
  const user = createUser(); // Setup own data
  const updated = updateUser(user);
  expect(updated).toBeDefined();
});
```

## Mock Wisely

### Mock External Dependencies

```javascript
// Good: Mock external API
jest.mock('./api');
test('fetch user data', async () => {
  api.fetchUser.mockResolvedValue({ name: 'John' });
  const user = await getUserData(1);
  expect(user.name).toBe('John');
});

// Bad: Mock internal logic
jest.mock('./calculator'); // Don't mock what you're testing
```

## Test Data Management

### Use Factories

```javascript
// factory.js
function createUser(overrides = {}) {
  return {
    id: 1,
    name: 'John Doe',
    email: 'john@example.com',
    ...overrides
  };
}

// test.js
test('should validate email', () => {
  const user = createUser({ email: 'invalid' });
  expect(validateUser(user)).toBe(false);
});
```

### Avoid Magic Numbers

```javascript
// Bad
expect(result).toBe(42);

// Good
const EXPECTED_TOTAL = 42;
expect(result).toBe(EXPECTED_TOTAL);
```

## Async Testing

### Always Handle Promises

```javascript
// Good
test('async operation', async () => {
  await expect(fetchData()).resolves.toBe('data');
});

// Bad: Missing await
test('async operation', () => {
  expect(fetchData()).resolves.toBe('data'); // Won't work
});
```

## Error Testing

### Test Error Cases

```javascript
test('should throw error for invalid input', () => {
  expect(() => divide(10, 0)).toThrow('Division by zero');
});

test('should reject promise on error', async () => {
  await expect(fetchData('invalid')).rejects.toThrow();
});
```

## Performance

### Keep Tests Fast

```javascript
// Good: Mock slow operations
jest.mock('./database');
test('query users', async () => {
  database.query.mockResolvedValue([{ id: 1 }]);
  const users = await getUsers();
  expect(users).toHaveLength(1);
});

// Bad: Real database calls in unit tests
test('query users', async () => {
  await database.connect(); // Slow
  const users = await database.query('SELECT * FROM users');
  expect(users).toBeDefined();
});
```

## Coverage Guidelines

### Aim for Meaningful Coverage

```javascript
// Good: Test business logic
test('should apply discount for premium users', () => {
  const user = { isPremium: true };
  const price = calculatePrice(100, user);
  expect(price).toBe(80);
});

// Bad: Test trivial code
test('getter returns value', () => {
  const obj = { value: 5 };
  expect(obj.value).toBe(5); // Pointless
});
```

## Avoid Common Pitfalls

### Don't Test Implementation Details

```javascript
// Bad: Testing internal state
test('counter increments', () => {
  counter.increment();
  expect(counter._value).toBe(1); // Private implementation
});

// Good: Test public behavior
test('counter increments', () => {
  counter.increment();
  expect(counter.getValue()).toBe(1);
});
```

### Don't Over-Mock

```javascript
// Bad: Mocking everything
jest.mock('./utils');
jest.mock('./helpers');
jest.mock('./validators');

// Good: Mock only external dependencies
jest.mock('./api');
```

## Documentation

### Write Clear Test Descriptions

```javascript
describe('User Authentication', () => {
  describe('login', () => {
    test('should return token when credentials are valid', () => {});
    test('should throw error when password is incorrect', () => {});
    test('should lock account after 5 failed attempts', () => {});
  });
});
```

## Continuous Improvement

### Review Test Failures

- Investigate flaky tests immediately
- Update tests when requirements change
- Remove obsolete tests
- Refactor duplicated test code

### Monitor Test Metrics

- Execution time
- Coverage percentage
- Flakiness rate
- Failure patterns

## Next Steps

- [Running Tests](./running-tests.md)
- [Writing Tests](./writing-tests.md)
- [Configuration](../configuration/README.md)

