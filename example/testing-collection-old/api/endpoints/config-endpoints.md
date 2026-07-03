---
outline_id: f13ba462-4b48-4702-9c5e-8b22fa2441c7
outline_collection: 54ef41a8-f481-4eb7-8eb4-efd2a7407bf6
outline_url: /doc/configuration-endpoints-MxXUZ2unOL
outline_updated: 2026-07-03T04:39:43.815Z
outline_revision: 2
---

# Configuration Endpoints

API endpoints for managing Testing Collection configuration.

## Get Configuration

Retrieve current configuration.

### Request

```http
GET /api/v1/config
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "testDir": "./tests",
    "timeout": 30000,
    "retries": 2,
    "workers": 4,
    "reporter": "html",
    "coverage": {
      "enabled": true,
      "threshold": 80
    }
  }
}
```

## Update Configuration

Modify configuration settings.

### Request

```http
PUT /api/v1/config
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json

{
  "timeout": 60000,
  "workers": 8,
  "coverage": {
    "enabled": true,
    "threshold": 85
  }
}
```

### Request Body

| Field | Type | Description |
|-------|------|-------------|
| testDir | string | Test directory path |
| timeout | integer | Global timeout (ms) |
| retries | integer | Retry count |
| workers | integer | Parallel workers |
| reporter | string | Reporter type |
| coverage | object | Coverage settings |

### Response

```json
{
  "success": true,
  "data": {
    "message": "Configuration updated successfully",
    "config": {
      "timeout": 60000,
      "workers": 8,
      "coverage": {
        "enabled": true,
        "threshold": 85
      }
    }
  }
}
```

## Validate Configuration

Check if configuration is valid.

### Request

```http
POST /api/v1/config/validate
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json

{
  "testDir": "./tests",
  "timeout": 30000,
  "workers": 4
}
```

### Response (Valid)

```json
{
  "success": true,
  "data": {
    "valid": true,
    "message": "Configuration is valid"
  }
}
```

### Response (Invalid)

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CONFIG",
    "message": "Configuration validation failed",
    "details": {
      "timeout": "Must be a positive integer",
      "testDir": "Directory does not exist"
    }
  }
}
```

## Reset Configuration

Reset to default configuration.

### Request

```http
POST /api/v1/config/reset
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "message": "Configuration reset to defaults",
    "config": {
      "testDir": "./tests",
      "timeout": 30000,
      "retries": 2,
      "workers": 4
    }
  }
}
```

## Get Configuration Schema

Retrieve configuration schema for validation.

### Request

```http
GET /api/v1/config/schema
Authorization: Bearer YOUR_API_KEY
```

### Response

```json
{
  "success": true,
  "data": {
    "schema": {
      "type": "object",
      "properties": {
        "testDir": {
          "type": "string",
          "description": "Test directory path"
        },
        "timeout": {
          "type": "integer",
          "minimum": 1000,
          "description": "Global timeout in milliseconds"
        }
      }
    }
  }
}
```

## Next Steps

- [Test Endpoints](./test-endpoints.md)
- [Report Endpoints](./report-endpoints.md)
- [Authentication](../authentication/README.md)

