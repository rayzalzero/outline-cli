---
outline_id: b9b83ee9-4c52-4d87-92a5-eef947f733c1
outline_collection: Testing
outline_url: /doc/config-command-inJkZ0lOZN
outline_updated: 2026-07-03T04:39:44.817Z
outline_revision: 4
outline_parent_id: e3a67fd1-e11e-49a7-8a30-584a7b4a903e
---

# Config Command

Manage Testing Collection configuration.

## Syntax

```bash
testing-collection config <action> [options]
```

## Actions

### show

Display current configuration.

```bash
testing-collection config show
```

### set

Set configuration value.

```bash
testing-collection config set <key> <value>
```

### get

Get configuration value.

```bash
testing-collection config get <key>
```

### reset

Reset configuration to defaults.

```bash
testing-collection config reset
```

### validate

Validate configuration file.

```bash
testing-collection config validate
```

## Options

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --global | -g    | boolean | Use global config |
| --config | -c    | string | Config file path |
| --json |       | boolean | Output as JSON |

## Examples

### Show Configuration

```bash
# Show current config
testing-collection config show

# Show as JSON
testing-collection config show --json

# Show global config
testing-collection config show --global
```

### Set Values

```bash
# Set timeout
testing-collection config set timeout 60000

# Set workers
testing-collection config set workers 4

# Set nested value
testing-collection config set coverage.threshold 85
```

### Get Values

```bash
# Get timeout
testing-collection config get timeout

# Get nested value
testing-collection config get coverage.threshold
```

### Reset Configuration

```bash
# Reset to defaults
testing-collection config reset

# Reset with confirmation
testing-collection config reset --yes
```

### Validate Configuration

```bash
# Validate current config
testing-collection config validate

# Validate specific file
testing-collection config validate --config ./my-config.js
```

## Configuration Keys

### Test Execution

* `testDir` - Test directory path
* `timeout` - Global timeout (ms)
* `retries` - Retry count
* `workers` - Parallel workers

### Reporting

* `reporter` - Reporter type
* `reportDir` - Report output directory

### Coverage

* `coverage.enabled` - Enable coverage
* `coverage.threshold` - Minimum coverage %
* `coverage.reporter` - Coverage reporter

## Output Examples

### Show Command

```bash
$ testing-collection config show

Configuration:
  testDir: ./tests
  timeout: 30000
  retries: 2
  workers: 4
  reporter: html
```

### JSON Output

```bash
$ testing-collection config show --json

{
  "testDir": "./tests",
  "timeout": 30000,
  "retries": 2,
  "workers": 4,
  "reporter": "html"
}
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0    | Success     |
| 2    | Invalid configuration |
| 3    | Configuration file not found |

## Next Steps

* [Run Command](./run.md)
* [Init Command](./init.md)
* [Configuration Guide](../../guides/configuration/README.md)

\-e \n\nTest update Thu Jul  2 18:20:06 WIB 2026