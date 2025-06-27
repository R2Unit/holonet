# Admin API Documentation

This document provides examples and documentation for the Admin API endpoints.

> **Note:** The API supports both RESTful and Action-based URL styles. Both styles are documented below and provide identical functionality.

## Authentication

All API requests require authentication using a bearer token in the Authorization header:

```
Authorization: Bearer your_token_here
```

## Endpoints

### Permissions

#### List All Permissions

**Request:**
```http
GET /api/admin/permissions
Authorization: Bearer your_token_here
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Administrator permission",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "name": "user",
    "description": "Regular user permission",
    "created_at": "2023-02-15T00:00:00Z",
    "updated_at": "2023-02-15T00:00:00Z"
  }
]
```

**Status Codes:**
- `200 OK`: Successfully retrieved the list of permissions
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Create a New Permission

**Request:**
```http
POST /api/admin/permissions
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "name": "editor",
  "description": "Editor permission"
}
```

**Response:**
```json
{
  "id": 3,
  "name": "editor",
  "description": "Editor permission",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T14:22:10Z"
}
```

**Status Codes:**
- `201 Created`: Permission successfully created
- `400 Bad Request`: Invalid request body or missing required fields
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Get a Specific Permission

**Request:**
```http
GET /api/admin/permissions/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "id": 2,
  "name": "user",
  "description": "Regular user permission",
  "created_at": "2023-02-15T00:00:00Z",
  "updated_at": "2023-02-15T00:00:00Z"
}
```

**Status Codes:**
- `200 OK`: Successfully retrieved the permission
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Permission not found
- `500 Internal Server Error`: Server error

#### Update a Permission

**Request:**
```http
PUT /api/admin/permissions/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "name": "basic-user",
  "description": "Basic user permission with limited access"
}
```

**Response:**
```json
{
  "id": 2,
  "name": "basic-user",
  "description": "Basic user permission with limited access",
  "created_at": "2023-02-15T00:00:00Z",
  "updated_at": "2023-06-16T15:30:45Z"
}
```

**Status Codes:**
- `200 OK`: Permission successfully updated
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Permission not found
- `500 Internal Server Error`: Server error

#### Delete a Permission

**Request:**
```http
DELETE /api/admin/permissions/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "message": "Permission deleted successfully"
}
```

**Status Codes:**
- `200 OK`: Permission successfully deleted
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Permission not found
- `500 Internal Server Error`: Server error

### Tokens

#### List All Tokens

**Request:**
```http
GET /api/admin/tokens
Authorization: Bearer your_token_here
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "token": "token1",
    "expires_at": "2023-06-17T14:22:10Z",
    "created_at": "2023-06-16T14:22:10Z",
    "updated_at": "2023-06-16T14:22:10Z",
    "status": "active",
    "policy_id": 1
  },
  {
    "id": 2,
    "user_id": 2,
    "token": "token2",
    "expires_at": "2023-06-18T14:22:10Z",
    "created_at": "2023-06-16T14:22:10Z",
    "updated_at": "2023-06-16T14:22:10Z",
    "status": "active",
    "policy_id": 2
  }
]
```

**Status Codes:**
- `200 OK`: Successfully retrieved the list of tokens
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Create a New Token

**Request:**
```http
POST /api/admin/tokens
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 3,
  "expires_at": "2023-07-16T14:22:10Z",
  "policy_id": 1
}
```

**Notes:**
- The `expires_at` field is optional and defaults to 24 hours from creation if not provided
- The `policy_id` field is optional

**Response:**
```json
{
  "id": 3,
  "user_id": 3,
  "token": "generated_token",
  "expires_at": "2023-07-16T14:22:10Z",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T14:22:10Z",
  "status": "active",
  "policy_id": 1
}
```

**Status Codes:**
- `201 Created`: Token successfully created
- `400 Bad Request`: Invalid request body or missing required fields
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Get a Specific Token

**Request:**
```http
GET /api/admin/tokens/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "id": 2,
  "user_id": 2,
  "token": "token2",
  "expires_at": "2023-06-18T14:22:10Z",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T14:22:10Z",
  "status": "active",
  "policy_id": 2
}
```

**Status Codes:**
- `200 OK`: Successfully retrieved the token
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Token not found
- `500 Internal Server Error`: Server error

#### Update a Token

**Request:**
```http
PUT /api/admin/tokens/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 2,
  "expires_at": "2023-08-16T14:22:10Z",
  "policy_id": 3
}
```

**Response:**
```json
{
  "id": 2,
  "user_id": 2,
  "token": "token2",
  "expires_at": "2023-08-16T14:22:10Z",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T15:30:45Z",
  "status": "active",
  "policy_id": 3
}
```

**Status Codes:**
- `200 OK`: Token successfully updated
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Token not found
- `500 Internal Server Error`: Server error

#### Revoke a Token

**Request:**
```http
POST /api/admin/tokens/2/revoke
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "id": 2,
  "user_id": 2,
  "token": "token2",
  "expires_at": "2023-06-18T14:22:10Z",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T15:30:45Z",
  "status": "revoked",
  "policy_id": 2
}
```

**Status Codes:**
- `200 OK`: Token successfully revoked
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Token not found
- `500 Internal Server Error`: Server error

#### Delete a Token

**Request:**
```http
DELETE /api/admin/tokens/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "message": "Token deleted successfully"
}
```

**Status Codes:**
- `200 OK`: Token successfully deleted
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Token not found
- `500 Internal Server Error`: Server error

### User Permissions

#### List All User Permissions

**Request:**
```http
GET /api/admin/user-permissions
Authorization: Bearer your_token_here
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "permission_id": 1,
    "created_at": "2023-06-16T14:22:10Z",
    "username": "admin",
    "permission_name": "admin"
  },
  {
    "id": 2,
    "user_id": 2,
    "permission_id": 2,
    "created_at": "2023-06-16T14:22:10Z",
    "username": "user1",
    "permission_name": "user"
  }
]
```

**Status Codes:**
- `200 OK`: Successfully retrieved the list of user permissions
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Assign a Permission to a User

**Request:**
```http
POST /api/admin/user-permissions
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 3,
  "permission_id": 2
}
```

**Response:**
```json
{
  "id": 3,
  "user_id": 3,
  "permission_id": 2,
  "created_at": "2023-06-16T14:22:10Z"
}
```

**Status Codes:**
- `201 Created`: Permission successfully assigned
- `400 Bad Request`: Invalid request body or missing required fields
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

#### Remove a Permission from a User

**Request:**
```http
DELETE /api/admin/user-permissions/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "message": "Permission removed successfully"
}
```

**Status Codes:**
- `200 OK`: Permission successfully removed
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User permission not found
- `500 Internal Server Error`: Server error

## Action-Based API Endpoints

The following endpoints provide the same functionality as the RESTful endpoints above, but with a more action-oriented URL structure.

### Permissions

#### List All Permissions

**Request:**
```http
GET /api/admin/permissions/list
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Create a New Permission

**Request:**
```http
POST /api/admin/permissions/create
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "name": "editor",
  "description": "Editor permission"
}
```

**Response:** Same as the RESTful endpoint.

#### Get a Specific Permission

**Request:**
```http
GET /api/admin/permissions/get/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Update a Permission

**Request:**
```http
PUT /api/admin/permissions/update/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "name": "basic-user",
  "description": "Basic user permission with limited access"
}
```

**Response:** Same as the RESTful endpoint.

#### Delete a Permission

**Request:**
```http
DELETE /api/admin/permissions/delete/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

### Tokens

#### List All Tokens

**Request:**
```http
GET /api/admin/tokens/list
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Create a New Token

**Request:**
```http
POST /api/admin/tokens/create
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 3,
  "expires_at": "2023-07-16T14:22:10Z",
  "policy_id": 1
}
```

**Response:** Same as the RESTful endpoint.

#### Get a Specific Token

**Request:**
```http
GET /api/admin/tokens/get/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Update a Token

**Request:**
```http
PUT /api/admin/tokens/update/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 2,
  "expires_at": "2023-08-16T14:22:10Z",
  "policy_id": 3
}
```

**Response:** Same as the RESTful endpoint.

#### Revoke a Token

**Request:**
```http
POST /api/admin/tokens/revoke/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Delete a Token

**Request:**
```http
DELETE /api/admin/tokens/delete/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

### User Permissions

#### List All User Permissions

**Request:**
```http
GET /api/admin/user-permissions/list
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

#### Assign a Permission to a User

**Request:**
```http
POST /api/admin/user-permissions/assign
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "user_id": 3,
  "permission_id": 2
}
```

**Response:** Same as the RESTful endpoint.

#### Remove a Permission from a User

**Request:**
```http
DELETE /api/admin/user-permissions/remove/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

## Error Responses

Error responses will have the following format:

```
Error message
```

For example:
```
Permission not found
```

## Rate Limiting

API requests are subject to rate limiting based on the token's policy. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.