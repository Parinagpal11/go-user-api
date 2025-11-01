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

## 🗓️ Day-by-Day Implementation Plan

### **Day 1 – Setup & Database**
- Initialize Go module and dependencies  
- Setup Docker Compose for PostgreSQL  
- Create `.env` and migration files  
- Test DB connection  

✅ *Deliverable:* Project runs & connects to DB successfully

---

### **Day 2 – User Model & Security**
- Define `User` struct with JSON tags  
- Implement password hashing (`bcrypt`)  
- Implement JWT generation & validation  
- Write tests for both utilities  

✅ *Deliverable:* Security utilities functional and tested

---

### **Day 3 – Authentication Endpoints**
- Implement `Register` and `Login` handlers  
- Add routes for `/api/auth/register` and `/api/auth/login`  
- Test with Postman  

✅ *Deliverable:* Register + Login + JWT flow working

---

### **Day 4 – Protected User Endpoints**
- Create Auth middleware (JWT validation)  
- Implement user CRUD operations  
- Apply middleware to `/api/users/*`  

✅ *Deliverable:* CRUD operations with authentication

---

### **Day 5 – Middleware & Error Handling**
- Add logging and CORS middleware  
- Implement standardized JSON error responses  
- Validate inputs and sanitize data  
- Move all secrets to `.env`  

✅ *Deliverable:* Production-ready backend

---

### **Day 6 – Documentation & Deployment**
- Write this README 😎  
- Create Dockerfile (multi-stage build)  
- Deploy using **Railway**, **Render**, or **Fly.io**  
- Add GitHub badges and CI setup  

✅ *Deliverable:* Fully documented and deployable API

---

## 🔧 Key Code Snippets

### 🧠 Main Server Setup

```go
func main() {
    // Load environment variables
    godotenv.Load()

    // Connect to database
    db := database.Connect()
    defer db.Close()

    // Setup router
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/api/auth/register", handlers.Register(db)).Methods("POST")
    r.HandleFunc("/api/auth/login", handlers.Login(db)).Methods("POST")

    // Protected routes
    api := r.PathPrefix("/api/users").Subrouter()
    api.Use(middleware.Auth)
    api.HandleFunc("", handlers.GetUsers(db)).Methods("GET")
    api.HandleFunc("/{id}", handlers.GetUser(db)).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", r))
}
```

---

### 🔐 Auth Middleware Example

```go
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        userID, err := utils.ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

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

## 🎯 Stretch Goals (Optional)

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

## 🚀 Next Steps

After completing this project, continue with:

**🧩 Project 3: Task Management App**
- Build a full-stack system using Next.js frontend + this Go API backend  
- Implement advanced features like notifications, analytics, and admin dashboards  

---

## 🏁 Author

**Pari Nagpal**  
_M.S. Computer Engineering @ Iowa State University_  
💻 GitHub: [@Parinagpal11](https://github.com/Parinagpal11)

---

⭐ **If you like this project, consider giving it a star!**
