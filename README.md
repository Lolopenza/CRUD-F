# CRUD API (Go + PostgreSQL)

Simple RESTful CRUD API written in Go using `net/http`, `gorilla/mux`, and PostgreSQL.

This project was built for learning purposes and covers:
- Clean project structure
- Handlers / repository separation
- Context usage
- Timeouts
- Graceful shutdown

---

## 🚀 Features

- Create user
- Get all users
- Get user by ID
- Update user
- Delete user
- PostgreSQL database
- Request context with timeout
- Server timeouts
- Graceful shutdown
- JSON helpers for responses

---

## 🧱 Project Structure

.
├── cmd/
│   └── api/
│       └── main.go        # application entry point
│
├── internal/
│   ├── handlers/          # HTTP handlers (HTTP layer)
│   ├── repository/        # DB logic (data layer)
│   ├── models/            # domain models
│   └── httphelper/        # JSON & error helpers
│
├── go.mod
└── README.md
