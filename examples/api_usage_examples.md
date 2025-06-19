# Using the Workflow API with curl

This document provides examples of how to use the workflow API with curl commands, including how to authenticate using tokens.

## Authentication

All API endpoints require authentication using a Bearer token. You need to include the token in the `Authorization` header of your requests:

```
Authorization: Bearer your-token-here
```

Replace `your-token-here` with a valid token from your database.

### Rate Limiting

The API now implements rate limiting based on token policies. Each token has rate limits that restrict:

1. The number of requests per minute
2. The maximum number of requests per day

If you exceed these limits, the API will return a `429 Too Many Requests` status code with an error message.

For more information about token policies and rate limiting, see [Token Policy Management API](api_policy_examples.md).

## API Endpoints

### List All Workflows

```bash
curl -X GET \
  http://localhost:3000/api/workflows \
  -H 'Authorization: Bearer your-token-here'
```

### Get a Specific Workflow

```bash
curl -X GET \
  http://localhost:3000/api/workflows/1 \
  -H 'Authorization: Bearer your-token-here'
```

Replace `1` with the ID of the workflow you want to retrieve.

### Create a New Workflow

```bash
curl -X POST \
  http://localhost:3000/api/workflows \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "My Workflow",
    "description": "A simple workflow example",
    "code": "package main\n\nfunc main() {\n  println(\"Hello, World!\")\n}"
  }'
```

### Update a Workflow

```bash
curl -X PUT \
  http://localhost:3000/api/workflows/1 \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Updated Workflow",
    "description": "An updated workflow example",
    "code": "package main\n\nfunc main() {\n  println(\"Hello, Updated World!\")\n}",
    "status": "active"
  }'
```

Replace `1` with the ID of the workflow you want to update.

### Schedule a Workflow (Run Immediately)

```bash
curl -X POST \
  http://localhost:3000/api/workflows/schedule \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "workflow_id": 1,
    "parameters": {
      "max_age_days": 30,
      "dry_run": false
    }
  }'
```

Replace `1` with the ID of the workflow you want to schedule. This will run the workflow immediately.

### Schedule a Workflow for Later Execution

```bash
curl -X POST \
  http://localhost:3000/api/workflows/schedule \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "workflow_id": 1,
    "parameters": {
      "max_age_days": 30,
      "dry_run": false
    },
    "scheduled_at": "2023-12-31T23:59:59Z"
  }'
```

Replace `1` with the ID of the workflow you want to schedule, and adjust the `scheduled_at` timestamp to your desired execution time (in RFC3339 format).

## Response Format

All API endpoints return JSON responses. For example, scheduling a workflow returns details about the execution:

```json
{
  "id": 123,
  "workflow_id": 1,
  "status": "pending",
  "parameters": {"max_age_days": 30, "dry_run": false},
  "scheduled_at": "2023-12-31T23:59:59Z",
  "created_at": "2023-12-30T10:15:30Z",
  "updated_at": "2023-12-30T10:15:30Z"
}
```

## Error Handling

If your request fails, you'll receive an HTTP error status code along with an error message. Common errors include:

- 401 Unauthorized: Invalid or missing token
- 400 Bad Request: Invalid request parameters
- 500 Internal Server Error: Server-side error

Example error response:

```
Invalid or expired token
```

## Getting a Valid Token

Tokens are stored in the `tokens` table in the database. You need to have a valid token in this table to authenticate your API requests. The token must not be expired (the `expires_at` column should be a future date).

You can create a token directly in the database or use your application's authentication system to generate one.
