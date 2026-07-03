---
outline_id: 4fc3d7d3-a338-4304-b466-41984b2aad3b
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/environment-variables-fOhGusXUVt
outline_updated: 2026-07-03T04:39:45.086Z
outline_revision: 2
---

# Environment Variables

Environment variables for Testing Collection configuration.

## Core Variables

### Authentication

| Variable | Description | Example |
|----------|-------------|---------|
| TESTING_COLLECTION_API_KEY | API authentication key | `tc_live_xxx` |
| TESTING_COLLECTION_CLIENT_ID | OAuth client ID | `oauth_client_xxx` |
| TESTING_COLLECTION_CLIENT_SECRET | OAuth client secret | `oauth_secret_xxx` |

### Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| TESTING_COLLECTION_CONFIG | Config file path | `./custom-config.js` |
| NODE_ENV | Environment mode | `test`, `development`, `production` |

### Test Execution

| Variable | Description | Example |
|----------|-------------|---------|
| TIMEOUT | Global timeout (ms) | `60000` |
| WORKERS | Number of workers | `4` |
| RETRIES | Retry count | `2` |

### Reporting

| Variable | Description | Example |
|----------|-------------|---------|
| REPORTER | Reporter type | `html`, `json`, `junit` |
| REPORT_DIR | Report output directory | `./test-results` |

### Coverage

| Variable | Description | Example |
|----------|-------------|---------|
| COVERAGE | Enable coverage | `true`, `false` |
| COVERAGE_THRESHOLD | Minimum coverage % | `80` |

## Usage

### Command Line

```bash
TIMEOUT=60000 testing-collection run
```

### .env File

```bash
# .env
TESTING_COLLECTION_API_KEY=tc_live_xxx
NODE_ENV=test
TIMEOUT=60000
WORKERS=4
COVERAGE=true
```

### Load in Config

```javascript
require('dotenv').config();

module.exports = {
  timeout: process.env.TIMEOUT || 30000,
  workers: process.env.WORKERS || 4,
  collectCoverage: process.env.COVERAGE === 'true'
};
```

## Environment-Specific Files

### .env.test

```bash
NODE_ENV=test
BASE_URL=http://localhost:3000
TIMEOUT=30000
```

### .env.production

```bash
NODE_ENV=production
BASE_URL=https://api.example.com
TIMEOUT=60000
RETRIES=3
```

### .env.local

```bash
# Local overrides (not committed)
TESTING_COLLECTION_API_KEY=tc_test_local_xxx
```

## CI/CD Variables

### GitHub Actions

```yaml
env:
  TESTING_COLLECTION_API_KEY: ${{ secrets.API_KEY }}
  NODE_ENV: test
  WORKERS: 2
```

### GitLab CI

```yaml
variables:
  NODE_ENV: test
  WORKERS: "2"
```

## Security Best Practices

### Never Commit Secrets

```gitignore
# .gitignore
.env
.env.local
.env.*.local
```

### Use Secrets Management

```bash
# AWS Secrets Manager
export TESTING_COLLECTION_API_KEY=$(aws secretsmanager get-secret-value --secret-id api-key --query SecretString --output text)

# HashiCorp Vault
export TESTING_COLLECTION_API_KEY=$(vault kv get -field=api_key secret/testing-collection)
```

## Next Steps

- [Config File Format](./config-format.md)
- [File Locations](./file-locations.md)
- [Security Guide](../../api/authentication/security.md)

