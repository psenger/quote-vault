# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./handlers
go test ./repository

# Format code
gofmt -w .
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DB_PATH` - SQLite database path (default: ./quotes.db)
- `LOG_LEVEL` - Logging level (default: info)
- `PAGE_SIZE` - Default pagination size (default: 10)
- `CORS_ORIGIN` - CORS allowed origins (default: *)

## Architecture

This is a Go REST API for storing and retrieving quotes, following a layered architecture:

```
main.go              → Application entry point, server setup, graceful shutdown
├── config/          → Environment configuration loading
├── router/          → Route definitions using gorilla/mux
├── middleware/      → HTTP middleware (CORS, logging, request ID, error handling)
├── handlers/        → HTTP request handlers (parse requests, call services, write responses)
├── services/        → Business logic layer
├── repository/      → Database access layer with context support
├── database/        → SQLite connection and direct DB operations
├── models/          → Data structures (Quote, request/response types)
├── validators/      → Request validation logic
├── errors/          → Custom error types (AppError with Code, Message, Type)
└── utils/           → Response helpers, pagination utilities
```

**Request flow:** Router → Middleware chain → Handler → Service → Repository → Database

## Key Patterns

- **Error handling**: Uses typed errors in `errors/errors.go` with `AppError` struct containing HTTP status code, message, and error type
- **Middleware**: Applied in order via `r.Use()` - RequestID → SecurityHeaders → CORS → Logging → ErrorHandler
- **Request ID**: Generated in middleware and propagated via context (`middleware.RequestIDKey`)
- **Response format**: Consistent JSON responses via `utils.SuccessResponse()` and `utils.ErrorResponse()`
- **Database**: SQLite with `mattn/go-sqlite3` driver; auto-creates tables on startup

## API Routes

All API routes are prefixed with `/api/v1`:
- `GET /health` and `/health/ready` - Health check endpoints
- `POST /api/v1/quotes` - Create quote
- `GET /api/v1/quotes` - List quotes (pagination via `page`, `limit` query params)
- `GET /api/v1/quotes/random` - Get random quote (optional `category` filter)
- `GET /api/v1/quotes/{id}` - Get quote by ID
- `PUT /api/v1/quotes/{id}` - Update quote
- `DELETE /api/v1/quotes/{id}` - Delete quote
- `GET /api/v1/categories` - List all categories

## Commit Convention

Use conventional commits: `feat:`, `fix:`, `docs:`, `style:`, `refactor:`, `test:`, `chore:`
