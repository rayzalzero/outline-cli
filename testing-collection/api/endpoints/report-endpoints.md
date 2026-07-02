---
outline_id: 280b69f4-f9a2-4e02-843b-51d68326aa2e
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/report-endpoints-B0Sk5o21OK
outline_updated: 2026-07-02T10:18:17.313Z
outline_revision: 2
---



# Report Endpoints

API endpoints for generating and managing test reports.

## List Reports

Get all generated reports.

### Request

```http
GET /api/v1/reports
Authorization: Bearer YOUR_API_KEY
```

### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| page | integer | Page number (default: 1) |
| limit | integer | Items per page (default: 50) |
| format | string | Filter by format (html/json/junit) |
| startDate | string | Filter from date (ISO 8601) |
| endDate | string | Filter to date (ISO 8601) |

### Response

```json
{
  "success": true,
  "data": {
    "reports": [
      {
        "id": "report_123",
        "format": "html",
        "runId": "run_abc123",
        "createdAt": "2026-07-02T08:43:48.419Z",
        "url": "https://reports.testing-collection.example.com/report_123.html"
      }
    ]
  },
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 25
  }
}
```

## Generate Report

Create a new test report.

### Request

```http
POST /api/v1/reports/generate
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json

{
  "runId": "run_abc123",
  "format": "html",
  "options": {
    "includeScreenshots": true,
    "includeLogs": true
  }
}
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| runId | string | Yes | Test run ID |
| format | string | Yes | Report format (html/json/junit) |
| options | object | No | Report options |

### Response

```json
{
  "success": true,
  "data": {
    "reportId": "report_123",
    "format": "html",
    "status": "generating",
    "url": "https://reports.testing-collection.example.com/report_123.html"
  }
}
```

## Get Report Details

Retrieve specific report information.

### Request

```http
GET /api/v1/reports/:id
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "report_123",
    "format": "html",
    "runId": "run_abc123",
    "status": "completed",
    "createdAt": "2026-07-02T08:43:48.419Z",
    "url": "https://reports.testing-collection.example.com/report_123.html",
    "summary": {
      "total": 10,
      "passed": 8,
      "failed": 2,
      "duration": 23926
    }
  }
}
```

## Download Report

Download report file.

### Request

```http
GET /api/v1/reports/:id/download
Authorization: Bearer YOUR_API_KEY
```

### Response

Binary file download with appropriate Content-Type header.

## Delete Report

Remove a report.

### Request

```http
DELETE /api/v1/reports/:id
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "message": "Report deleted successfully"
  }
}
```

## Get Report Summary

Get aggregated report statistics.

### Request

```http
GET /api/v1/reports/summary
Authorization: Bearer YOUR_API_KEY
```

### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| startDate | string | From date (ISO 8601) |
| endDate | string | To date (ISO 8601) |

### Response

```json
{
  "success": true,
  "data": {
    "totalReports": 100,
    "totalTests": 1000,
    "passRate": 85.5,
    "averageDuration": 15234,
    "trends": {
      "daily": [
        {
          "date": "2026-07-01",
          "passed": 45,
          "failed": 5
        }
      ]
    }
  }
}
```

## Error Responses

### Report Not Found

```json
{
  "success": false,
  "error": {
    "code": "REPORT_NOT_FOUND",
    "message": "Report with ID report_123 not found"
  }
}
```

### Invalid Format

```json
{
  "success": false,
  "error": {
    "code": "INVALID_FORMAT",
    "message": "Report format must be one of: html, json, junit"
  }
}
```

## Next Steps

- [Test Endpoints](./test-endpoints.md)
- [Configuration Endpoints](./config-endpoints.md)
- [Authentication](../authentication/README.md)

