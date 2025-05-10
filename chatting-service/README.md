# 🗨️ Chatting Service API

A real-time messaging platform built with **Go (Fiber)** for the backend and **React** for the frontend.

---

## 🚀 Features

### 🔐 Authentication
- User registration with email/password
- JWT-based authentication
- Refresh token support
- Password change functionality

### 💬 Messaging
- Direct 1:1 messaging
- Broadcast messaging to multiple users
- Message status tracking (sent/delivered/read)
- Message history with pagination
- Conversation threads
- Message deletion

### 📎 Media Handling
- File uploads (JPEG, PNG, PDF)
- 10MB max file size
- Local storage with public URL access

### 🔄 Real-Time Features
- WebSocket-based real-time updates
- Online/offline status tracking

---

## 📡 API Endpoints

### 🧾 Authentication
| Method | Endpoint              | Description                 |
|--------|-----------------------|-----------------------------|
| POST   | `/auth/login`         | User login                  |
| POST   | `/auth/register`      | User registration           |
| POST   | `/auth/change-password` | Change password (auth)    |

### 👤 Users
| Method | Endpoint              | Description                 |
|--------|-----------------------|-----------------------------|
| GET    | `/api/users/profile`  | Get user profile            |
| PUT    | `/api/users/profile`  | Update user profile         |
| GET    | `/api/users/all`      | Get all users               |

### ✉️ Messages
| Method | Endpoint                                     | Description                          |
|--------|----------------------------------------------|--------------------------------------|
| POST   | `/api/messages`                              | Send direct message                  |
| POST   | `/api/messages/broadcast`                    | Send broadcast message               |
| GET    | `/api/messages/conversation/{userID}`        | Get conversation with a user         |
| GET    | `/api/messages/conversations`                | Get all user conversations           |
| PUT    | `/api/messages/{id}/read`                    | Mark message as read                 |
| DELETE | `/api/messages/{id}`                         | Delete message                       |

### 🖼️ Media
| Method | Endpoint              | Description                 |
|--------|-----------------------|-----------------------------|
| POST   | `/api/media/upload`   | Upload media file           |

### 🔌 WebSocket
| Method | Endpoint              | Description                 |
|--------|-----------------------|-----------------------------|
| GET    | `/ws`                 | WebSocket connection        |

---

## 🏗️ Architecture

### Clean Architecture Layers

#### **Domain**
- Core business logic and entities
- Models: `User`, `Message`, `Media`
- Repository interfaces
- Domain service interfaces

#### **Application**
- Use cases: `AuthService`, `MessageService`, `MediaService`, `UserService`

#### **Infrastructure**
- PostgreSQL repositories
- JWT implementation
- Local file storage
- WebSocket notifier

#### **Delivery**
- REST API using Fiber
- WebSocket handlers

---

## 🗃️ Database Schema

- `Users` table
- `Messages` table
- `MessageRecipients` join table (for broadcast support)
- PostgreSQL `ENUMs` for message types and status

---

## 💻 Frontend Application

Located in `twerlo_chat_app/FE`.

### 🧭 Running the Frontend

```bash
cd twerlo_chat_app/FE
npm install
npm run dev
```

The frontend will run on: [http://127.0.0.1:8081](http://127.0.0.1:8081)

---

## ⚙️ Getting Started

### 📋 Prerequisites

- Docker `20.10+`
- Go `1.21+`
- Node.js `16+`
- PostgreSQL `14`

### 🛠️ Installation

```bash
git clone https://github.com/AmeerHeiba/chatting-service.git
cd chatting-service

# Setup environment
cp .env.example .env
# Edit the .env file with your own configuration

# Start database
docker-compose up -d postgres

# Run migrations
go run cmd/migrate/main.go

# Start backend server
go run cmd/api/main.go
```

In a separate terminal:

```bash
cd twerlo_chat_app/FE
npm install
npm run dev
```

---

## 📄 API Documentation

- Swagger is available at: [http://localhost:8080/swagger](http://localhost:8080/swagger)

---

## 🧪 Development Workflow

### 🧬 Running Tests
```bash
go test ./...
```

### 🛠️ Code Generation (Swagger)

Swagger docs are generated using [`swaggo/swag`](https://github.com/swaggo/swag):

```bash
swag init
```

---

## 🔐 Environment Variables

Example variables in `.env`:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chatting_service

JWT_SECRET=your-secret-key

MEDIA_STORAGE_PATH=./uploads
MEDIA_BASE_URL=http://localhost:8080/media
```

---

## 🚢 Deployment

### Docker

```bash
docker-compose up --build
```

### Kubernetes

Helm charts – **coming soon**

---

## 📦 Postman Collection

To test the API easily:

1. **Import Collection:** `Twerlo chat app.postman_collection.json`
2. **Import Environment:** `Twerlo-env.json`
3. **Select the environment** from the dropdown
4. **Login via `/api/auth/login`** to auto-fetch token
5. Other requests will use the token automatically via `Bearer {{access_token}}`

---

## 📈 Planned Enhancements

- Group messaging
- Message reactions
- End-to-end encryption
- Push notifications
- Media upload to cloud storage (e.g., S3)
- Advanced message search

---

## 📬 Contact

For questions or collaboration: **Ameer Heiba**  
[GitHub](https://github.com/AmeerHeiba)
