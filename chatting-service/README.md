# Chatting Service

A real-time messaging platform built with Go, Fiber, and PostgreSQL implementing Clean Architecture and CQRS patterns.

## Features

- **User Authentication**: JWT-based auth with refresh tokens
- **Messaging**: 1:1 and group conversations
- **Real-Time Updates**: WebSocket notifications
- **Media Support**: File uploads with storage abstraction
- **Transactional Safety**: Atomic operation guarantees

## Architecture
Presentation → Application → Domain ← Infrastructure

### Core Patterns
- **Clean Architecture**: Domain-centric design
- **CQRS**: Separate command and query paths
- **Repository Pattern**: Persistence abstraction
- **Transaction Management**: Cross-operation atomicity


## Folder Structure
```
/chatting-service
├── /cmd
│ └── /api
│ └── main.go # App entry (dependency injection)
├── /internal
│ ├── /config # Env/config loading
│ ├── /domain # Entities, value objects, repo interfaces
│ ├── /infrastructure # External implementations
│ │ ├── /database # PostgreSQL repositories
│ │ ├── /brokers # RabbitMQ/Kafka adapters
│ │ └── /storage # Filesystem/S3 storage
│ ├── /application # Use cases/services
│ ├── /delivery # Transport layers
│ │ ├── /http # REST handlers (Fiber)
│ │ └── /websocket # Real-time handlers
│ └── /shared # Common utilities (logging, errors)
├── /migrations # Database schema changes
├── /pkg # Reusable library code
├── /public # Static files/uploads
├── /web # Frontend assets (HTML/JS/CSS)
├── go.mod # Go dependencies
├── go.sum
└── Dockerfile # Multi-stage build
```

## Tech Stack

| Component       | Technology          |
|-----------------|---------------------|
| Language        | Go 1.21+            |
| Web Framework   | Fiber v2            |
| Database        | PostgreSQL 14       |
| ORM             | GORM                |
| Real-Time       | Gorilla WebSocket   |
| Error Handling  | Custom middleware   |

## Key Components

### Domain Layer
```go
// Example repository interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
}
```

### Transaction Management
```go
// Atomic operation example
err := txManager.WithTransaction(ctx, func(ctx context.Context, repos *Repositories) error {
    // Transactional operations
})
```
### Error Handling

Structured error responses with:
HTTP status codes
Machine-readable error codes
Contextual message

## Getting Started

### Prerequisites
- Docker 20.10+
- Go 1.21+

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/AmeerHeiba/chatting-service.git
   cd chatting-service
2. Setup environment:
    cp .env.example .env
3. Start services:
    docker-compose up -d --build
4. Run migrations:
    docker-compose exec app go run cmd/migrate/main.go
5. Access services:
    API: http://localhost:8080
    RabbitMQ: http://localhost:15672 (guest/guest)
    PGAdmin: http://localhost:5050 (configure server)

## API Documentation
Interactive Swagger docs available at http://localhost:8080/swagger when running locally.

## Development Workflow
1. Start dependencies:
    docker-compose up -d postgres rabbitmq
2. Run application locally:
    go run cmd/api/main.go
3. Run tests:
    go test ./...
