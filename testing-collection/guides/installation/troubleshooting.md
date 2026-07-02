---
outline_id: 683d8c10-ee7b-4005-8fd3-77dfc73b7a26
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/installation-troubleshooting-9PgfqJuiSM
outline_updated: 2026-07-02T10:18:18.260Z
outline_revision: 2
---



# Installation Troubleshooting

Common installation issues and their solutions.

## Permission Errors

### Error: EACCES

```bash
npm ERR! Error: EACCES: permission denied
```

**Solution 1**: Use sudo (not recommended)

```bash
sudo npm install -g testing-collection
```

**Solution 2**: Fix npm permissions (recommended)

```bash
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
npm install -g testing-collection
```

## Network Issues

### Error: ETIMEDOUT

```bash
npm ERR! network request to https://registry.npmjs.org/testing-collection failed
```

**Solutions**:

```bash
# Use different registry
npm config set registry https://registry.npm.taobao.org

# Increase timeout
npm config set fetch-timeout 60000

# Use proxy
npm config set proxy http://proxy.company.com:8080
```

## Version Conflicts

### Error: Incompatible Node Version

```bash
npm ERR! engine Unsupported engine
```

**Solution**: Upgrade Node.js

```bash
# Using nvm
nvm install 18
nvm use 18

# Verify
node --version
```

## Missing Dependencies

### Error: Cannot find module

```bash
Error: Cannot find module 'some-dependency'
```

**Solution**: Clean install

```bash
# Remove node_modules and lock file
rm -rf node_modules package-lock.json

# Clear npm cache
npm cache clean --force

# Reinstall
npm install -g testing-collection
```

## Platform-Specific Issues

### macOS: Command Not Found

**Solution**: Add to PATH

```bash
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Windows: PowerShell Execution Policy

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Linux: Missing Build Tools

```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"
```

## Verification Issues

### Command Not Found After Installation

**Check installation location**:

```bash
npm list -g testing-collection
npm bin -g
```

**Add to PATH**:

```bash
export PATH="$(npm bin -g):$PATH"
```

## Getting Help

If issues persist:

1. Check [GitHub Issues](https://github.com/example/testing-collection/issues)
2. Run diagnostics: `testing-collection doctor`
3. Enable verbose logging: `npm install -g testing-collection --verbose`
4. Contact support with error logs

## Next Steps

- [System Requirements](./system-requirements.md)
- [Installation Methods](./installation-methods.md)
- [Configuration](../configuration/README.md)

