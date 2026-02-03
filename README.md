# Go Backend Template

Backend template chuẩn với authentication, được viết bằng Golang.

## Tech Stack

- **Framework**: Gin
- **Database**: MongoDB + SAP SQL Server
- **Authentication**: External API (App Manager)

## Cấu trúc thư mục

```
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── config/             # Configuration & database connections
│   ├── handlers/           # HTTP handlers (controllers)
│   ├── middleware/         # Middleware (auth, cors)
│   ├── models/             # Data models
│   ├── services/           # Business logic
│   └── utils/              # Utilities
├── env/                    # Environment files
├── go.mod                  # Go modules
└── .env.example            # Environment template
```

## Setup

1. Copy environment file:

```bash
cp .env.example env/.env.development
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run server:

```bash
go run ./cmd/server
```

## Build

```bash
go build -o server ./cmd/server
./server
```

## API Endpoints

| Method | Path                   | Description                          |
| ------ | ---------------------- | ------------------------------------ |
| GET    | `/`                    | Health check                         |
| GET    | `/v1/auth/permissions` | Get user permissions (requires auth) |

## Adding New Modules

1. Create handler in `internal/handlers/`
2. Create service in `internal/services/` (if needed)
3. Register routes in `cmd/server/main.go` → `setupRoutes()`
4. Add new permissions in `internal/services/permission.go`
