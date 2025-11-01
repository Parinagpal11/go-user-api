# 🚀 REST API with Go — Complete Implementation Guide

A production-ready **User Management REST API** built using **Go**, **PostgreSQL**, and **Docker**, featuring full **authentication**, **middleware**, and **CRUD** operations.

---

## 📋 Project Overview

**Duration:** 4–6 days  
**Tech Stack:** Go · PostgreSQL · Docker  
**Goal:** Build a scalable user management API with authentication and JWT-based security.

---

## 🏗️ Architecture Overview

```
go-user-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── auth.go
│   │   └── users.go
│   ├── models/                  # Data models
│   │   └── user.go
│   ├── database/                # DB connection
│   │   └── db.go
│   ├── middleware/              # Auth, logging, CORS
│   │   ├── auth.go
│   │   └── logger.go
│   └── utils/                   # Helpers (JWT, password)
│       ├── jwt.go
│       └── hash.go
├── migrations/                  # SQL migrations
│   └── 001_create_users.sql
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🧩 Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🌐 API Endpoints

| Method | Endpoint | Description |
|---------|-----------|-------------|
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/login` | Login and receive JWT token |
| `GET` | `/api/users` | Get all users *(Protected)* |
| `GET` | `/api/users/:id` | Get user by ID *(Protected)* |
| `PUT` | `/api/users/:id` | Update user *(Protected)* |
| `DELETE` | `/api/users/:id` | Delete user *(Protected)* |
| `GET` | `/api/users/me` | Get current user *(Protected)* |

---




## ✅ Testing Checklist

- [x] Register new user  
- [x] Login and receive JWT token  
- [x] Access protected routes with valid JWT  
- [x] CRUD operations for users  
- [x] Proper error handling and response codes  
- [x] Passwords stored as hashed values  
- [x] CORS and middleware verified  

---

## 🎯 Future Improvements

- 🔁 Email verification flow  
- 🔑 Password reset  
- 🔄 Refresh tokens  
- ⚖️ Role-based access control  
- 🔍 Search, filter, pagination  
- 📄 Swagger/OpenAPI documentation  
- 🧪 CI/CD with GitHub Actions  

---

## 📚 Resources

- [Go Docs](https://go.dev/doc/)  
- [Gorilla Mux](https://github.com/gorilla/mux)  
- [PostgreSQL Driver (`lib/pq`)](https://github.com/lib/pq)  
- [JWT RFC 8725](https://datatracker.ietf.org/doc/html/rfc8725)  
- [REST API Design Principles](https://restfulapi.net/)  

---


## 🏁 Author

**Pari Nagpal**  
_M.S. Computer Engineering @ Iowa State University_  
💻 GitHub: [@Parinagpal11](https://github.com/Parinagpal11)

---

⭐ **If you like this project, consider giving it a star!**
