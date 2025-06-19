# Token Policy Management API

This document provides examples of how to use the token policy management API with curl commands, including how to create, update, and delete token policies, and how to assign policies to tokens.

## Authentication

All API endpoints require authentication using a Bearer token. You need to include the token in the `Authorization` header of your requests:

```
Authorization: Bearer your-token-here
```

Replace `your-token-here` with a valid token from your database.

## Rate Limiting and Request Policies

The API now supports rate limiting and request policies for tokens. Each token can have a policy assigned to it that defines:

- **Rate Limit Per Minute**: Maximum number of requests allowed per minute
- **Maximum Requests Per Day**: Maximum number of requests allowed per day

If a token exceeds its rate limit, the API will return a `429 Too Many Requests` status code with an error message.

## Policy Management API Endpoints

### List All Policies

```bash
curl -X GET \
  http://localhost:3000/api/policies \
  -H 'Authorization: Bearer your-token-here'
```

### Get a Specific Policy

```bash
curl -X GET \
  http://localhost:3000/api/policies/1 \
  -H 'Authorization: Bearer your-token-here'
```

Replace `1` with the ID of the policy you want to retrieve.

### Create a New Policy

```bash
curl -X POST \
  http://localhost:3000/api/policies \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Standard API Access",
    "description": "Standard API access with moderate rate limits",
    "rate_limit_per_min": 60,
    "max_requests_per_day": 1000,
    "active": true
  }'
```

### Update a Policy

```bash
curl -X PUT \
  http://localhost:3000/api/policies/1 \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Premium API Access",
    "description": "Premium API access with higher rate limits",
    "rate_limit_per_min": 120,
    "max_requests_per_day": 5000,
    "active": true
  }'
```

Replace `1` with the ID of the policy you want to update.

### Delete a Policy

```bash
curl -X DELETE \
  http://localhost:3000/api/policies/1 \
  -H 'Authorization: Bearer your-token-here'
```

Replace `1` with the ID of the policy you want to delete. Note that you cannot delete a policy that is currently assigned to one or more tokens.

### Assign a Policy to a Token

```bash
curl -X PUT \
  http://localhost:3000/api/tokens/policy \
  -H 'Authorization: Bearer your-token-here' \
  -H 'Content-Type: application/json' \
  -d '{
    "token_id": 1,
    "policy_id": 2
  }'
```

Replace `1` with the ID of the token and `2` with the ID of the policy you want to assign.

## Response Format

All API endpoints return JSON responses. For example, creating a policy returns details about the policy:

```json
{
  "id": 1,
  "name": "Standard API Access",
  "description": "Standard API access with moderate rate limits",
  "rate_limit_per_min": 60,
  "max_requests_per_day": 1000,
  "active": true
}
```

## Error Handling

If your request fails, you'll receive an HTTP error status code along with an error message. Common errors include:

- 401 Unauthorized: Invalid or missing token
- 400 Bad Request: Invalid request parameters
- 404 Not Found: Policy or token not found
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server-side error

Example error response for rate limit exceeded:

```
Rate limit exceeded: Maximum 60 requests per minute allowed.
```

## Default Policy

If a token does not have a policy assigned, it will use the default rate limits:

- 60 requests per minute
- 1000 requests per day

You can assign a policy to a token to customize these limits.