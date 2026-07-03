---
outline_id: 2551f215-1721-47a4-9605-f4d07eb7b254
outline_collection: Testing
outline_url: /doc/system-requirements-879HIvL6gV
outline_updated: 2026-07-03T04:39:44.398Z
outline_revision: 4
outline_parent_id: bc84adf7-0156-4c47-8002-e9545a218a57
---

# System Requirements

Ensure your system meets these requirements before installing Testing Collection.

## Minimum Requirements

### Operating System

* **macOS**: 10.15 (Catalina) or later
* **Linux**: Ubuntu 18.04+, Debian 10+, CentOS 7+, or equivalent
* **Windows**: Windows 10 or later, Windows Server 2016+

### Runtime

* **Node.js**: 16.x or higher (LTS recommended)
* **npm**: 7.x or higher
* **Memory**: 2GB RAM minimum, 4GB recommended
* **Disk Space**: 500MB free space

## Recommended Requirements

For optimal performance:

* **Node.js**: 18.x LTS or 20.x LTS
* **Memory**: 8GB RAM
* **CPU**: Multi-core processor (4+ cores)
* **Disk Space**: 2GB free space

## Dependencies

### Required

* Node.js runtime
* npm or yarn package manager

### Optional

* Git (for version control integration)
* Docker (for containerized testing)
* Chrome/Chromium (for browser testing)

## Compatibility

### Node.js Versions

* ✅ Node.js 16.x
* ✅ Node.js 18.x (LTS)
* ✅ Node.js 20.x (LTS)
* ⚠️ Node.js 14.x (deprecated)

### Package Managers

* ✅ npm 7.x+
* ✅ yarn 1.22+
* ✅ pnpm 7.x+

## Checking Your System

Run these commands to verify your system:

```bash
# Check Node.js version
node --version

# Check npm version
npm --version

# Check available memory
free -h  # Linux
vm_stat  # macOS

# Check disk space
df -h
```

## Upgrading

If your system doesn't meet requirements:

### Upgrade Node.js

```bash
# Using nvm
nvm install 18
nvm use 18

# Using n
npm install -g n
n lts
```

### Upgrade npm

```bash
npm install -g npm@latest
```

## Next Steps

* [Installation Methods](./installation-methods.md)
* [Troubleshooting](./troubleshooting.md)