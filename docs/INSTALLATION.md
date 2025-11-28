# Installation Guide

## Prerequisites

- Go 1.19 or later
- SQLite (included with Go's database/sql)

## Quick Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/quote-vault.git
cd quote-vault
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Environment

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` file with your preferred settings:

```env
PORT=8080
DB_PATH=./quotes.db
LOG_LEVEL=info
ENVIRONMENT=development
```

### 4. Run the Application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## Development Setup

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o quote-vault main.go
./quote-vault
```

### Docker Setup (Optional)

If you prefer using Docker:

```bash
# Build the image
docker build -t quote-vault .

# Run the container
docker run -p 8080:8080 -v $(pwd)/quotes.db:/app/quotes.db quote-vault
```

## Configuration Options

| Variable | Description | Default |
|----------|-------------|----------|
| `PORT` | HTTP server port | 8080 |
| `DB_PATH` | SQLite database file path | ./quotes.db |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |
| `ENVIRONMENT` | Environment mode (development, production) | development |

## Health Check

Once running, verify the installation:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the `PORT` in your `.env` file
2. **Database permission errors**: Ensure the application has write permissions to the database directory
3. **Module not found**: Run `go mod tidy` to clean up dependencies

For more detailed API usage, see [API.md](API.md).