# Chatting Service

A real-time messaging service built with Go, Fiber, PostgreSQL, and WebSockets implementing Clean Architecture and CQRS pattern.

## Features

- User authentication (JWT)
- 1:1 messaging with CQRS separation
- Broadcast messaging via RabbitMQ
- Media attachments with local storage
- Real-time updates via WebSockets

## Tech Stack

- **Backend**: Go 1.21, Fiber
- **Database**: PostgreSQL 14
- **Real-time**: WebSockets (Gorilla)
- **Broker**: RabbitMQ (for message fan-out)
- **Storage**: Local filesystem (extensible to S3)
- **Architecture**: Clean Architecture + CQRS

## Architecture Overview
Presentation → Application → Domain ← Infrastructure

### Key Patterns:
1. **Clean Architecture**:
   - Domain-centric design
   - Framework-independent core
   - Testable components

2. **CQRS**:
   - Separate command (write) and query (read) paths
   - Optimized read models for message history
   - Transactional write models for message sending

3. **Domain-Driven Design**:
   - Explicit bounded contexts (Auth, Messaging, Media)
   - Repository pattern for persistence

## Folder Structure
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
