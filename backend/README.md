# My Project Backend

Go backend API server

## Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional)

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Install dependencies:

```bash
   go mod download
```

## Running the Server

### Development

```bash
go run cmd/api/main.go
```

### Build and Run

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

## API Endpoints

- `GET /health` - Health check
- `GET /api/users` - Get all users

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware functions
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   └── database/            # Database operations
├── pkg/                     # Public packages
└── api/                     # API documentation
```
