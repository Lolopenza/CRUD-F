# CRUD API (Go + PostgreSQL)

Backend REST API service written in Go using `net/http`, `gorilla/mux`, and PostgreSQL.

Designed to follow real-world backend patterns and demonstrates:
- Clean and modular project structure
- Separation of handlers and repositories
- Context usage and request timeouts
- Graceful shutdown handling
- JSON response helpers

---

## 🚀 Features

- Create, read, update, and delete user records
- PostgreSQL integration
- Request context with timeout
- Server read/write timeouts
- Graceful shutdown on interrupt
- Helper utilities for JSON and error responses

---

## 🧱 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go            # Application entry point
│
├── internal/
│   ├── handlers/              # HTTP handlers (API layer)
│   ├── repository/            # Database logic (data access layer)
│   ├── models/                # Domain models
│   └── httphelper/            # JSON & error handling utilities
│
├── go.mod
└── README.md
```

---

## ▶️ Running the Project

Make sure PostgreSQL is running and accessible. Then run the app:

```bash
go run cmd/api/main.go
```

By default, the server starts on:

**http://localhost:3838**

---

## 📡 API Endpoints

### Health Check
**GET** `/healthcheck`  
Returns server status and uptime info.

### Users
| Method | Endpoint        | Description        |
|:-------|:----------------|:------------------|
| POST   | `/users`        | Create a new user |
| GET    | `/users`        | Get all users     |
| GET    | `/users/{id}`   | Get user by ID    |
| PUT    | `/users/{id}`   | Update user by ID |
| DELETE | `/users/{id}`   | Delete user by ID |

---
