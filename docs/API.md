# Quote Vault API Documentation

## Overview

Quote Vault is a REST API for storing and retrieving inspirational quotes. This document provides comprehensive information about all available endpoints.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, no authentication is required for API access.

## Endpoints

### Health Check

#### GET /health

Check if the API is running and healthy.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

### Quotes

#### GET /quotes

Retrieve all quotes with pagination support.

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `limit` (optional, default: 10) - Number of quotes per page
- `category` (optional) - Filter by category

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/quotes?page=1&limit=5&category=motivation"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "text": "The only way to do great work is to love what you do.",
      "author": "Steve Jobs",
      "category": "motivation",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 5,
    "total": 25,
    "total_pages": 5
  }
}
```

#### POST /quotes

Add a new quote to the collection.

**Request Body:**
```json
{
  "text": "Be yourself; everyone else is already taken.",
  "author": "Oscar Wilde",
  "category": "wisdom"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Be yourself; everyone else is already taken.",
    "author": "Oscar Wilde",
    "category": "wisdom"
  }'
```

**Response:**
```json
{
  "data": {
    "id": 26,
    "text": "Be yourself; everyone else is already taken.",
    "author": "Oscar Wilde",
    "category": "wisdom",
    "created_at": "2024-01-15T11:30:00Z",
    "updated_at": "2024-01-15T11:30:00Z"
  }
}
```

#### GET /quotes/random

Get a random quote from all quotes or a specific category.

**Query Parameters:**
- `category` (optional) - Get random quote from specific category

**Example Request (all categories):**
```bash
curl http://localhost:8080/api/v1/quotes/random
```

**Example Request (specific category):**
```bash
curl "http://localhost:8080/api/v1/quotes/random?category=motivation"
```

**Response:**
```json
{
  "data": {
    "id": 15,
    "text": "Success is not final, failure is not fatal: it is the courage to continue that counts.",
    "author": "Winston Churchill",
    "category": "motivation",
    "created_at": "2024-01-15T09:45:00Z",
    "updated_at": "2024-01-15T09:45:00Z"
  }
}
```

## Error Responses

The API uses standard HTTP status codes and returns errors in the following format:

```json
{
  "error": {
    "message": "Validation failed",
    "details": [
      "text is required",
      "author must be at least 2 characters long"
    ]
  }
}
```

### Common Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server error

## Data Validation

### Quote Object

- `text`: Required, minimum 10 characters, maximum 1000 characters
- `author`: Required, minimum 2 characters, maximum 100 characters
- `category`: Required, minimum 2 characters, maximum 50 characters, lowercase letters and hyphens only

### Pagination

- `page`: Must be positive integer, minimum 1
- `limit`: Must be positive integer, minimum 1, maximum 100

## Rate Limiting

Currently, no rate limiting is implemented. This may be added in future versions.

## Examples

### Adding Multiple Quotes

```bash
# Add a motivation quote
curl -X POST http://localhost:8080/api/v1/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Your limitation—it's only your imagination.",
    "author": "Unknown",
    "category": "motivation"
  }'

# Add a wisdom quote
curl -X POST http://localhost:8080/api/v1/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "It is during our darkest moments that we must focus to see the light.",
    "author": "Aristotle",
    "category": "wisdom"
  }'
```

### Retrieving Quotes with Different Filters

```bash
# Get first page of all quotes
curl "http://localhost:8080/api/v1/quotes"

# Get second page with 20 quotes per page
curl "http://localhost:8080/api/v1/quotes?page=2&limit=20"

# Get all motivation quotes
curl "http://localhost:8080/api/v1/quotes?category=motivation"

# Get random quote from wisdom category
curl "http://localhost:8080/api/v1/quotes/random?category=wisdom"
```

## Support

For questions or issues, please refer to the project repository or create an issue on GitHub.