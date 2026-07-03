---
outline_id: 1b21eb6c-c99c-4d5f-afbf-98e84c4c33ef
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/init-command-ObTnFLnjUS
outline_updated: 2026-07-03T04:39:44.834Z
outline_revision: 2
---

# Init Command

Initialize a new Testing Collection project.

## Syntax

```bash
testing-collection init [options] [directory]
```

## Options

| Option | Alias | Type | Description |
|--------|-------|------|-------------|
| --template | -t | string | Project template |
| --config | -c | boolean | Generate config file |
| --force | -f | boolean | Overwrite existing files |
| --yes | -y | boolean | Skip prompts |

## Templates

Available project templates:

- `basic` - Basic test setup (default)
- `typescript` - TypeScript configuration
- `react` - React testing setup
- `node` - Node.js testing setup
- `api` - API testing setup

## Examples

### Basic Initialization

```bash
# Initialize in current directory
testing-collection init

# Initialize in new directory
testing-collection init my-project

# With specific template
testing-collection init --template typescript
```

### Skip Prompts

```bash
# Use defaults
testing-collection init --yes

# With template
testing-collection init --template react --yes
```

### Force Overwrite

```bash
# Overwrite existing files
testing-collection init --force
```

## Generated Files

### Basic Template

```
my-project/
├── tests/
│   └── example.test.js
├── testing-collection.config.js
├── .gitignore
└── package.json
```

### TypeScript Template

```
my-project/
├── tests/
│   └── example.test.ts
├── testing-collection.config.ts
├── tsconfig.json
├── .gitignore
└── package.json
```

## Interactive Prompts

When running without `--yes`:

```
? Project name: my-project
? Template: (Use arrow keys)
  ❯ basic
    typescript
    react
    node
    api
? Generate config file? (Y/n)
? Install dependencies? (Y/n)
```

## Next Steps

After initialization:

```bash
cd my-project
npm install
testing-collection run
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 2 | Directory already exists (without --force) |
| 3 | Invalid template |

## Next Steps

- [Run Command](./run.md)
- [Config Command](./config.md)
- [Getting Started Guide](../../getting-started.md)

