# Users API Documentation

This document provides examples and documentation for the Users API endpoints.

> **Note:** The API now supports both RESTful and Action-based URL styles. Both styles are documented below and provide identical functionality.

## Authentication

All API requests require authentication using a bearer token in the Authorization header:

```
Authorization: Bearer your_token_here
```

## Endpoints

### List All Users

**Request:**
```http
GET /api/users
Authorization: Bearer your_token_here
```

**Response:**
```json
[
  {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "status": "active",
    "last_login": "2023-06-15T10:30:45Z",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "username": "user1",
    "email": "user1@example.com",
    "status": "active",
    "last_login": "2023-06-14T08:15:30Z",
    "created_at": "2023-02-15T00:00:00Z",
    "updated_at": "2023-02-15T00:00:00Z"
  }
]
```

**Status Codes:**
- `200 OK`: Successfully retrieved the list of users
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

### Create a New User

**Request:**
```http
POST /api/users
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "status": "active"
}
```

**Notes:**
- The `status` field is optional and defaults to "active" if not provided

**Response:**
```json
{
  "id": 3,
  "username": "newuser",
  "email": "newuser@example.com",
  "status": "active",
  "created_at": "2023-06-16T14:22:10Z",
  "updated_at": "2023-06-16T14:22:10Z"
}
```

**Status Codes:**
- `201 Created`: User successfully created
- `400 Bad Request`: Invalid request body or missing required fields
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

### Get a Specific User

**Request:**
```http
GET /api/users/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "id": 2,
  "username": "user1",
  "email": "user1@example.com",
  "status": "active",
  "last_login": "2023-06-14T08:15:30Z",
  "created_at": "2023-02-15T00:00:00Z",
  "updated_at": "2023-02-15T00:00:00Z"
}
```

**Status Codes:**
- `200 OK`: Successfully retrieved the user
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

### Update a User

**Request:**
```http
PUT /api/users/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "username": "updateduser",
  "email": "updated@example.com",
  "password": "newpassword456",
  "status": "active"
}
```

**Notes:**
- You can update just one field by only including that field in the request body
- Password is optional in update requests
- Status can be set to "active" or "suspended"

**Response:**
```json
{
  "id": 2,
  "username": "updateduser",
  "email": "updated@example.com",
  "status": "active",
  "last_login": "2023-06-14T08:15:30Z",
  "created_at": "2023-02-15T00:00:00Z",
  "updated_at": "2023-06-16T15:30:45Z"
}
```

**Status Codes:**
- `200 OK`: User successfully updated
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

### Suspend a User

**Request:**
```http
POST /api/users/2/suspend
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "success": true,
  "message": "User suspended successfully"
}
```

**Status Codes:**
- `200 OK`: User successfully suspended
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

### Unsuspend a User

**Request:**
```http
POST /api/users/2/unsuspend
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "success": true,
  "message": "User unsuspended successfully"
}
```

**Status Codes:**
- `200 OK`: User successfully unsuspended
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

### Delete a User

**Request:**
```http
DELETE /api/users/2
Authorization: Bearer your_token_here
```

**Response:**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

**Notes:**
- This is a soft delete that sets the `deleted_at` timestamp
- Deleted users will not appear in list or get operations

**Status Codes:**
- `200 OK`: User successfully deleted
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

## Error Responses

Error responses will have the following format:

```
Error message
```

For example:
```
User not found
```

## Action-Based API Endpoints

The following endpoints provide the same functionality as the RESTful endpoints above, but with a more action-oriented URL structure.

### List All Users

**Request:**
```http
GET /api/users/list
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

### Create a New User

**Request:**
```http
POST /api/users/create
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "status": "active"
}
```

**Response:** Same as the RESTful endpoint.

### Get a Specific User

**Request:**
```http
GET /api/users/get/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

### Update a User

**Request:**
```http
PUT /api/users/update/2
Authorization: Bearer your_token_here
Content-Type: application/json

{
  "username": "updateduser",
  "email": "updated@example.com",
  "password": "newpassword456",
  "status": "active"
}
```

**Response:** Same as the RESTful endpoint.

### Delete a User

**Request:**
```http
DELETE /api/users/delete/2
Authorization: Bearer your_token_here
```

**Response:** Same as the RESTful endpoint.

## Rate Limiting

API requests are subject to rate limiting based on the token's policy. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.
