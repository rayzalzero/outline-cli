---
outline_id: 15ca9746-e4bd-408a-9db8-cb686c51e491
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/api-documentation-sjttdtMFQN
outline_updated: 2026-07-02T10:18:16.362Z
outline_revision: 2
---



# API Documentation

Complete API reference for Testing Collection.

## Table of Contents

- [Authentication](./authentication/README.md)
- [Endpoints](./endpoints/README.md)

## Overview

Testing Collection provides a comprehensive REST API for programmatic access to all testing functionality.

## Base URL

```
https://api.testing-collection.example.com/v1
```

## Authentication

All API requests require authentication. See [Authentication Guide](./authentication/README.md) for details.

## Request Format

```http
POST /api/v1/tests/run
Content-Type: application/json
Authorization: Bearer YOUR_API_TOKEN

{
  "testDir": "./tests",
  "filter": "auth"
}
```

## Response Format

```json
{
  "success": true,
  "data": {
    "results": [],
    "summary": {
      "total": 10,
      "passed": 8,
      "failed": 2
    }
  }
}
```

## Error Handling

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Test directory not found"
  }
}
```

## Rate Limiting

- 1000 requests per hour per API key
- Rate limit headers included in responses

## Pagination

```http
GET /api/v1/tests?page=1&limit=50
```

## Next Steps

- [Authentication](./authentication/README.md)
- [Endpoints](./endpoints/README.md)


