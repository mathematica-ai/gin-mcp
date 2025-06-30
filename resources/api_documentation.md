# API Documentation

This document provides comprehensive documentation for the REST API endpoints.

## Authentication

All API requests require authentication using Bearer tokens in the Authorization header:

```
Authorization: Bearer <your-token>
```

## Base URL

```
https://api.example.com/v1
```

## Endpoints

### Users

#### GET /users
Retrieve a list of users.

**Query Parameters:**
- `page` (integer): Page number (default: 1)
- `limit` (integer): Number of items per page (default: 20)
- `search` (string): Search term for username or email

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5
  }
}
```

#### POST /users
Create a new user.

**Request Body:**
```json
{
  "username": "new_user",
  "email": "new@example.com",
  "password": "secure_password"
}
```

**Response:**
```json
{
  "id": 2,
  "username": "new_user",
  "email": "new@example.com",
  "created_at": "2024-01-15T11:00:00Z"
}
```

### Products

#### GET /products
Retrieve a list of products.

**Query Parameters:**
- `category_id` (integer): Filter by category
- `min_price` (decimal): Minimum price filter
- `max_price` (decimal): Maximum price filter
- `in_stock` (boolean): Filter by stock availability

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Sample Product",
      "description": "A sample product description",
      "price": 29.99,
      "category_id": 1,
      "stock_quantity": 50
    }
  ]
}
```

#### POST /products
Create a new product.

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 19.99,
  "category_id": 1,
  "stock_quantity": 100
}
```

### Orders

#### GET /orders
Retrieve a list of orders.

**Query Parameters:**
- `user_id` (integer): Filter by user
- `status` (string): Filter by order status
- `date_from` (date): Filter orders from date
- `date_to` (date): Filter orders to date

#### POST /orders
Create a new order.

**Request Body:**
```json
{
  "user_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Validation failed",
  "details": {
    "field": "error message"
  }
}
```

### 401 Unauthorized
```json
{
  "error": "Authentication required"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Rate Limiting

API requests are limited to 1000 requests per hour per API key. Rate limit headers are included in all responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642248000
```

## Pagination

List endpoints support pagination with the following response format:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
``` 