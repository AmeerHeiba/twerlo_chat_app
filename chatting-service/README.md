# Chatting Service

A real-time messaging platform built with Go, Fiber, and PostgreSQL implementing Clean Architecture.

## Features

- ✅ **User Authentication**: JWT-based auth with refresh tokens (implemented)
- ✅ **Direct Messaging**: 1:1 conversations (implemented)
- 🚧 **Broadcast Messaging**: Functional but needs WebSocket integration (in progress)
- 🚧 **Media Support**: Model ready - storage service in development (planned)
- ✅ **Transactional Safety**: Atomic operation guarantees (implemented)
- 🚧 **Real-Time Updates**: Interfaces defined - WebSocket implementation planned

## Architecture
Presentation → Application → Domain ← Infrastructure

### Core Patterns
- **Clean Architecture**: Domain-centric design
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
│ │ └── /storage # (Planned) Filesystem/S3 storage
│ ├── /application # Use cases/services
│ ├── /delivery # Transport layers
│ │ ├── /http # REST handlers (Fiber)
│ │ └── /websocket # (Planned) Real-time handlers
│ └── /shared # Common utilities (logging, errors)
├── /migrations # Database schema changes
├── go.mod # Go dependencies
├── go.sum
└── Dockerfile # Multi-stage build
```

## Tech Stack

| Component       | Technology          | Status        |
|-----------------|---------------------|---------------|
| Language        | Go 1.21+            | ✅ Implemented |
| Web Framework   | Fiber v2            | ✅ Implemented |
| Database        | PostgreSQL 14       | ✅ Implemented |
| ORM             | GORM                | ✅ Implemented |
| Real-Time       | (Planning)          | 🚧 Interfaces |
| Error Handling  | Custom middleware   | ✅ Implemented |

## Key Components

### Domain Layer
```go
// Message repository interface
type MessageRepository interface {
    Create(ctx context.Context, senderID uint, content, mediaURL string, 
           messageType MessageType) (*Message, error)
    FindConversation(ctx context.Context, user1ID, user2ID uint, 
                    query MessageQuery) ([]Message, error)
    // ... actual implemented methods ...
}
```

### Transaction Management
```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Atomic operations
    if err := tx.Create(&message).Error; err != nil {
        return err
    }
    return tx.Create(&recipient).Error
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
- Air (for live reload during development)

### Installation
1. Clone the repository:
```bash
   git clone https://github.com/AmeerHeiba/chatting-service.git
   cd chatting-service
```
2. Setup environment:
```bash
    cp .env.example .env
```
3. Start services:
```bash
    docker-compose up --d postgres
```
4. Run migrations:
```bash
    go run migrate/main.go
```
5. Start development server:
```bash
    air
```

## API Documentation
Interactive Swagger docs available at http://localhost:8080/swagger when running locally.

## Development Workflow
1. Start dependencies:
```bash
    docker-compose up -d postgres
```
2. Run application (with live reload):
```bash
    air
```
3. Run tests:
```bash
    go test ./...
```

## Planned Enhancements
- WebSocket real-time messaging
- Media upload service (local/S3)
- Swagger API documentation
- Advanced message status tracking