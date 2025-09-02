package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"todoy/handlers"
	"todoy/models"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database configuration
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "todoapp")

	// Create database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Connected to MySQL database successfully!")

	// Initialize models
	todoModel := models.NewTodoModel(db)

	// Create table if it doesn't exist
	err = todoModel.CreateTable()
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoModel)

	// Set up routes
	mux := http.NewServeMux()

	// Home route
	mux.HandleFunc("/", todoHandler.HomeHandler)

	// Todo routes
	mux.HandleFunc("/todos/new", todoHandler.CreateTodoFormHandler)
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			todoHandler.CreateTodoHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Dynamic routes for specific todos
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/todos/")

		if strings.HasSuffix(path, "/edit") {
			todoHandler.EditTodoFormHandler(w, r)
		} else if strings.HasSuffix(path, "/toggle") {
			todoHandler.ToggleTodoHandler(w, r)
		} else if r.Method == http.MethodDelete {
			todoHandler.DeleteTodoHandler(w, r)
		} else if r.Method == http.MethodPost || r.Method == http.MethodPut {
			todoHandler.UpdateTodoHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Static files (if needed)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Add middleware for method override (to support PUT and DELETE from forms)
	handler := methodOverrideMiddleware(mux)

	// Start server
	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Visit http://localhost:%s to view the application\n", port)

	// Add request logging middleware
	loggedHandler := loggingMiddleware(handler)

	log.Fatal(http.ListenAndServe(":"+port, loggedHandler))
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// loggingMiddleware logs HTTP requests and responses
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// methodOverrideMiddleware allows forms to use PUT and DELETE methods
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.FormValue("_method")
			if method == "PUT" || method == "DELETE" {
				r.Method = method
			}
		}
		next.ServeHTTP(w, r)
	})
}
