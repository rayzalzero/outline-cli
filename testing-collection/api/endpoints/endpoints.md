---
outline_id: 5abcae8f-51c2-4353-ac92-620e7b0d46d4
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/api-endpoints-VT2lmfqvV9
outline_updated: 2026-07-02T10:18:17.178Z
outline_revision: 2
---



# API Endpoints

Complete reference for Testing Collection API endpoints.

## Table of Contents

- [Test Endpoints](./test-endpoints.md)
- [Configuration Endpoints](./config-endpoints.md)
- [Report Endpoints](./report-endpoints.md)

## Base URL

```
https://api.testing-collection.example.com/v1
```

## Endpoints Overview

### Tests

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tests` | List all tests |
| POST | `/tests/run` | Run tests |
| GET | `/tests/:id` | Get test details |
| DELETE | `/tests/:id` | Delete test |

### Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/config` | Get configuration |
| PUT | `/config` | Update configuration |
| POST | `/config/validate` | Validate configuration |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports` | List reports |
| GET | `/reports/:id` | Get report details |
| POST | `/reports/generate` | Generate report |

## Common Parameters

### Pagination

```
?page=1&limit=50
```

### Filtering

```
?filter=auth&status=passed
```

### Sorting

```
?sort=createdAt&order=desc
```

## Response Format

### Success Response

```json
{
  "success": true,
  "data": { },
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 100
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {}
  }
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 429 | Rate Limit Exceeded |
| 500 | Internal Server Error |

## Rate Limits

- 1000 requests/hour for standard keys
- Rate limit info in response headers:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1625097600
```

## Next Steps

- [Test Endpoints](./test-endpoints.md)
- [Configuration Endpoints](./config-endpoints.md)
- [Report Endpoints](./report-endpoints.md)

