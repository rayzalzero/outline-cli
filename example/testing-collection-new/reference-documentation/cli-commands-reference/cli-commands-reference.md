---
outline_id: e3a67fd1-e11e-49a7-8a30-584a7b4a903e
outline_collection: Testing
outline_url: /doc/cli-commands-reference-YrqyV9Vq6V
outline_updated: 2026-07-03T04:39:43.456Z
outline_revision: 4
outline_parent_id: 93006387-0fed-4138-9725-81d5aa1c10eb
---

# CLI Commands Reference

Complete reference for Testing Collection command-line interface.

## Table of Contents

* [Run Command](./run.md)
* [Init Command](./init.md)
* [Config Command](./config.md)

## Global Options

Available for all commands:

| Option | Alias | Description |
|--------|-------|-------------|
| --version | -v    | Show version |
| --help | -h    | Show help   |
| --verbose |       | Verbose output |
| --silent |       | Suppress output |
| --config | -c    | Config file path |

## Command Overview

### testing-collection run

Execute tests.

```bash
testing-collection run [options] [files...]
```

[Full documentation](./run.md)

### testing-collection init

Initialize new project.

```bash
testing-collection init [options] [directory]
```

[Full documentation](./init.md)

### testing-collection config

Manage configuration.

```bash
testing-collection config <action> [options]
```

[Full documentation](./config.md)

## Exit Codes

| Code | Description |
|------|-------------|
| 0    | Success     |
| 1    | Test failures |
| 2    | Configuration error |
| 3    | Runtime error |
| 4    | Invalid arguments |

## Environment Variables

| Variable | Description |
|----------|-------------|
| TESTING_COLLECTION_API_KEY | API authentication key |
| TESTING_COLLECTION_CONFIG | Config file path |
| NODE_ENV | Environment (test/development/production) |

## Examples

### Run all tests

```bash
testing-collection run
```

### Run with coverage

```bash
testing-collection run --coverage
```

### Initialize project

```bash
testing-collection init my-project
```

### View configuration

```bash
testing-collection config show
```

## Next Steps

* [Run Command](./run.md)
* [Init Command](./init.md)
* [Config Command](./config.md)