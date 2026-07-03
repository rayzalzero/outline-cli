---
outline_id: d3b6d62f-7519-4e96-89ca-e391a4f3ec4e
outline_collection: Testing
outline_url: /doc/run-command-uKTU2u6Bbp
outline_updated: 2026-07-03T04:39:44.988Z
outline_revision: 4
outline_parent_id: e3a67fd1-e11e-49a7-8a30-584a7b4a903e
---

# Run Command

Execute tests with Testing Collection.

## Syntax

```bash
testing-collection run [options] [files...]
```

## Options

### Test Selection

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --filter | -f    | string | Filter test files by pattern |
| --grep | -g    | string | Filter tests by name |
| --tag  | -t    | string | Run tests with specific tag |

### Execution

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --workers | -w    | integer | Number of parallel workers |
| --timeout |       | integer | Global timeout (ms) |
| --retries | -r    | integer | Retry failed tests |
| --bail | -b    | integer | Stop after N failures |

### Output

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --report |       | string | Reporter type (html/json/junit) |
| --output | -o    | string | Output file path |
| --verbose | -v    | boolean | Verbose output |
| --silent |       | boolean | Suppress output |

### Coverage

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --coverage |       | boolean | Collect coverage |
| --coverage-reporter |       | string | Coverage reporter |
| --coverage-threshold |       | integer | Minimum coverage % |

### Watch Mode

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --watch |       | boolean | Watch mode  |
| --watch-all |       | boolean | Watch all files |

### Debugging

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --debug |       | boolean | Debug mode  |
| --inspect |       | boolean | Node inspector |
| --inspect-brk |       | boolean | Break on start |

## Examples

### Basic Usage

```bash
# Run all tests
testing-collection run

# Run specific directory
testing-collection run tests/unit

# Run specific file
testing-collection run tests/auth.test.js
```

### Filtering

```bash
# Filter by file pattern
testing-collection run --filter "auth"

# Filter by test name
testing-collection run --grep "login"

# Multiple filters
testing-collection run --filter "auth|user"
```

### Parallel Execution

```bash
# Use 4 workers
testing-collection run --workers 4

# Use 50% of CPU cores
testing-collection run --workers 50%

# Sequential execution
testing-collection run --workers 1
```

### Coverage

```bash
# Collect coverage
testing-collection run --coverage

# With threshold
testing-collection run --coverage --coverage-threshold 80

# Specific reporter
testing-collection run --coverage --coverage-reporter html
```

### Watch Mode

```bash
# Watch mode
testing-collection run --watch

# Watch specific files
testing-collection run --watch --filter "auth"
```

### Debugging

```bash
# Debug mode
testing-collection run --debug

# With inspector
testing-collection run --inspect

# Break on start
testing-collection run --inspect-brk
```

### Reporting

```bash
# HTML report
testing-collection run --report html

# Multiple reporters
testing-collection run --report html,json,junit

# Custom output
testing-collection run --report json --output results.json
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0    | All tests passed |
| 1    | Some tests failed |
| 2    | Configuration error |
| 3    | Runtime error |

## Next Steps

* [Init Command](./init.md)
* [Config Command](./config.md)
* [Running Tests Guide](../../guides/usage/running-tests.md)