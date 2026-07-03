---
outline_id: a4b72714-a41e-4110-8544-860712c2a2b2
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/getting-started-L09CJNFjYx
outline_updated: 2026-07-03T04:39:42.489Z
outline_revision: 2
---

# Getting Started

This guide will help you get up and running with Testing Collection in just a few minutes.

## Prerequisites

Before you begin, ensure you have:
- Node.js 16.x or higher
- npm or yarn package manager
- Basic knowledge of JavaScript/TypeScript

## Quick Start

### 1. Installation

Install Testing Collection using npm:

```bash
npm install -g testing-collection
```

Or using yarn:

```bash
yarn global add testing-collection
```

### 2. Initialize Your Project

Create a new project:

```bash
testing-collection init my-project
cd my-project
```

### 3. Run Your First Test

Create a simple test file `test/example.test.js`:

```javascript
describe('Example Test', () => {
  it('should pass', () => {
    expect(true).toBe(true);
  });
});
```

Run the test:

```bash
testing-collection run
```

## Next Steps

Now that you have Testing Collection installed and running, explore these topics:

- [Installation Guide](./guides/installation/README.md) - Detailed installation options
- [Configuration](./guides/configuration/README.md) - Customize your setup
- [Usage Guide](./guides/usage/README.md) - Learn advanced features
- [API Documentation](./api/README.md) - Integrate programmatically

## Common Tasks

### Running Specific Tests

```bash
testing-collection run --filter "MyTest"
```

### Watch Mode

```bash
testing-collection run --watch
```

### Generate Reports

```bash
testing-collection run --report html
```

## Troubleshooting

If you encounter issues:
1. Check the [Installation Guide](./guides/installation/troubleshooting.md)
2. Review [Configuration](./guides/configuration/common-issues.md)
3. Consult the [CLI Reference](./reference/cli-commands/README.md)

