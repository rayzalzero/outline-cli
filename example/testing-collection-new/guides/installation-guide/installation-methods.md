---
outline_id: 4acbb585-cf77-4458-9c04-51a024f50542
outline_collection: Testing
outline_url: /doc/installation-methods-zB2JhafGhc
outline_updated: 2026-07-03T04:39:44.318Z
outline_revision: 4
outline_parent_id: bc84adf7-0156-4c47-8002-e9545a218a57
---

# Installation Methods

Explore different ways to install Testing Collection based on your needs.

## Global Installation

Install Testing Collection globally to use it across all projects.

### npm

```bash
npm install -g testing-collection
```

### yarn

```bash
yarn global add testing-collection
```

### pnpm

```bash
pnpm add -g testing-collection
```

## Local Installation

Install Testing Collection as a project dependency.

### npm

```bash
npm install --save-dev testing-collection
```

### yarn

```bash
yarn add -D testing-collection
```

### pnpm

```bash
pnpm add -D testing-collection
```

## From Source

Build and install from source code:

```bash
# Clone repository
git clone https://github.com/example/testing-collection.git
cd testing-collection

# Install dependencies
npm install

# Build
npm run build

# Link globally
npm link
```

## Docker Installation

Run Testing Collection in a Docker container:

```bash
# Pull image
docker pull testing-collection:latest

# Run container
docker run -it testing-collection:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  testing-collection:
    image: testing-collection:latest
    volumes:
      - ./tests:/app/tests
    command: run
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Install Testing Collection
  run: npm install -g testing-collection
```

### GitLab CI

```yaml
install:
  script:
    - npm install -g testing-collection
```

### Jenkins

```groovy
sh 'npm install -g testing-collection'
```

## Version Management

### Install Specific Version

```bash
npm install -g testing-collection@1.2.3
```

### Install Latest Beta

```bash
npm install -g testing-collection@beta
```

### List Available Versions

```bash
npm view testing-collection versions
```

## Offline Installation

For environments without internet access:

```bash
# On connected machine
npm pack testing-collection

# Transfer .tgz file to offline machine
npm install -g testing-collection-1.2.3.tgz
```

## Next Steps

* [System Requirements](./system-requirements.md)
* [Troubleshooting](./troubleshooting.md)
* [Configuration](../configuration/README.md)