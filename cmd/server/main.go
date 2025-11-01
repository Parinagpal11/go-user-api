package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/yourusername/go-user-api/internal/database"
	"github.com/yourusername/go-user-api/internal/handlers"
	"github.com/yourusername/go-user-api/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Initialize router
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	// Public routes
	r.HandleFunc("/api/auth/register", handlers.Register(db)).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.Login(db)).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api/users").Subrouter()
	api.Use(middleware.Auth)
	api.HandleFunc("", handlers.GetUsers(db)).Methods("GET")
	api.HandleFunc("/me", handlers.GetCurrentUser(db)).Methods("GET")
	api.HandleFunc("/{id}", handlers.GetUser(db)).Methods("GET")
	api.HandleFunc("/{id}", handlers.UpdateUser(db)).Methods("PUT")
	api.HandleFunc("/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
