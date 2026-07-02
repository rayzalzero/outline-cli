---
outline_id: 1002e270-55cd-45a9-88fc-df4a201918d3
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/test-endpoints-VuC6lSWUVj
outline_updated: 2026-07-02T10:18:17.471Z
outline_revision: 2
---



# Test Endpoints

API endpoints for managing and running tests.

## List Tests

Get a list of all tests.

### Request

```http
GET /api/v1/tests
Authorization: Bearer YOUR_API_KEY
```

### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| page | integer | Page number (default: 1) |
| limit | integer | Items per page (default: 50) |
| filter | string | Filter by test name |
| status | string | Filter by status (passed/failed/skipped) |
| sort | string | Sort field (name/createdAt/duration) |
| order | string | Sort order (asc/desc) |

### Response

```json
{
  "success": true,
  "data": {
    "tests": [
      {
        "id": "test_123",
        "name": "User authentication",
        "file": "tests/auth.test.js",
        "status": "passed",
        "duration": 1234,
        "createdAt": "2026-07-02T08:43:48.419Z"
      }
    ]
  },
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 100
  }
}
```

## Run Tests

Execute tests programmatically.

### Request

```http
POST /api/v1/tests/run
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json

{
  "testDir": "./tests",
  "filter": "auth",
  "workers": 4,
  "timeout": 30000,
  "retries": 2
}
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| testDir | string | Yes | Test directory path |
| filter | string | No | Filter test files |
| workers | integer | No | Number of parallel workers |
| timeout | integer | No | Timeout in milliseconds |
| retries | integer | No | Number of retries |
| coverage | boolean | No | Collect coverage |

### Response

```json
{
  "success": true,
  "data": {
    "runId": "run_abc123",
    "status": "running",
    "startedAt": "2026-07-02T08:43:48.419Z",
    "summary": {
      "total": 10,
      "passed": 0,
      "failed": 0,
      "skipped": 0
    }
  }
}
```

## Get Test Details

Retrieve details for a specific test.

### Request

```http
GET /api/v1/tests/:id
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "test_123",
    "name": "User authentication",
    "file": "tests/auth.test.js",
    "status": "passed",
    "duration": 1234,
    "assertions": 5,
    "errors": [],
    "createdAt": "2026-07-02T08:43:48.419Z",
    "updatedAt": "2026-07-02T08:43:49.653Z"
  }
}
```

## Get Test Run Status

Check the status of a running test.

### Request

```http
GET /api/v1/tests/runs/:runId
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "runId": "run_abc123",
    "status": "completed",
    "startedAt": "2026-07-02T08:43:48.419Z",
    "completedAt": "2026-07-02T08:44:12.345Z",
    "duration": 23926,
    "summary": {
      "total": 10,
      "passed": 8,
      "failed": 2,
      "skipped": 0
    },
    "results": [
      {
        "testId": "test_123",
        "name": "should login successfully",
        "status": "passed",
        "duration": 1234
      }
    ]
  }
}
```

## Stop Test Run

Cancel a running test.

### Request

```http
POST /api/v1/tests/runs/:runId/stop
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "runId": "run_abc123",
    "status": "cancelled",
    "message": "Test run cancelled successfully"
  }
}
```

## Delete Test

Remove a test from the system.

### Request

```http
DELETE /api/v1/tests/:id
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "message": "Test deleted successfully"
  }
}
```

## Error Responses

### Test Not Found

```json
{
  "success": false,
  "error": {
    "code": "TEST_NOT_FOUND",
    "message": "Test with ID test_123 not found"
  }
}
```

### Invalid Test Directory

```json
{
  "success": false,
  "error": {
    "code": "INVALID_TEST_DIR",
    "message": "Test directory './tests' does not exist"
  }
}
```

## Next Steps

- [Configuration Endpoints](./config-endpoints.md)
- [Report Endpoints](./report-endpoints.md)
- [Authentication](../authentication/README.md)

