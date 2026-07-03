---
outline_id: 74c8f12f-d210-4bee-bb34-60062bd279db
outline_collection: Testing
outline_url: /doc/api-documentation-mpgECAD4oy
outline_updated: 2026-07-03T04:43:45.18Z
outline_revision: 5
---

# API Documentation

Complete API reference for Testing Collection.

## Table of Contents

* [Authentication](./authentication/README.md)
* [Endpoints](./endpoints/README.md)

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

* 1000 requests per hour per API key
* Rate limit headers included in responses

## Pagination

```http
GET /api/v1/tests?page=1&limit=50
```

## Next Steps

* [Authentication](./authentication/README.md)
* [Endpoints](./endpoints/README.md)