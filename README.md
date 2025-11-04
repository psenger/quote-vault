# quote-vault

A lightweight REST API service for storing and retrieving inspirational quotes with support for random quote selection and categorization.

## Installation

```bash
go mod tidy
```

## Usage

Start the server:

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### API Endpoints

- `POST /quotes` - Add a new quote
- `GET /quotes/random` - Get a random quote from all categories
- `GET /quotes/random?category=motivation` - Get a random quote from specific category
- `GET /quotes` - List all quotes with pagination (default: page=1, limit=10)
- `GET /quotes?page=2&limit=5` - List quotes with custom pagination

### Example Usage

Add a new quote:
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"text":"The only way to do great work is to love what you do.","author":"Steve Jobs","category":"motivation"}'
```

Get a random quote:
```bash
curl http://localhost:8080/quotes/random
```

List quotes with pagination:
```bash
curl "http://localhost:8080/quotes?page=1&limit=5"
```