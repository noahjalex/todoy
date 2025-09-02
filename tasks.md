# Todo List CRUD App - Project Plan

## Project Overview
Create a simple Golang CRUD application for managing todo tasks with:
- HTML templates for frontend
- MySQL database
- HTMX for dynamic interactions
- Tailwind CSS for styling
- Basic CRUD operations (Create, Read, Update, Delete)

## Project Structure
```
todoy/
├── main.go                 # Main application entry point
├── models/
│   └── todo.go            # Todo model and database operations
├── handlers/
│   └── todo_handlers.go   # HTTP handlers for CRUD operations
├── templates/
│   ├── layout.html        # Base layout template
│   ├── index.html         # Home page listing all todos
│   ├── create.html        # Create todo form
│   └── edit.html          # Edit todo form
├── static/
│   └── style.css          # Custom CSS (if needed)
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

## Database Schema
```sql
CREATE TABLE todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Required Dependencies
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/gorilla/mux` - HTTP router (optional, can use net/http)

## Routes to Implement
- `GET /` - Home page (list all todos)
- `GET /todos/new` - Show create todo form
- `POST /todos` - Create new todo
- `GET /todos/{id}/edit` - Show edit todo form
- `PUT /todos/{id}` - Update existing todo
- `DELETE /todos/{id}` - Delete todo
- `POST /todos/{id}/toggle` - Toggle todo completion status

## Tasks Checklist

### Setup & Configuration
- [ ] Initialize Go module
- [ ] Install required dependencies
- [ ] Create project directory structure
- [ ] Set up database connection configuration

### Database Layer
- [ ] Create Todo model struct
- [ ] Implement database connection
- [ ] Create todo table schema
- [ ] Implement CRUD database operations:
  - [ ] GetAllTodos()
  - [ ] GetTodoByID()
  - [ ] CreateTodo()
  - [ ] UpdateTodo()
  - [ ] DeleteTodo()
  - [ ] ToggleTodoStatus()

### HTTP Handlers
- [ ] Create todo handlers file
- [ ] Implement home page handler (list todos)
- [ ] Implement create todo handlers (GET form, POST create)
- [ ] Implement edit todo handlers (GET form, PUT update)
- [ ] Implement delete todo handler
- [ ] Implement toggle completion handler

### Templates
- [ ] Create base layout template with Tailwind CSS and HTMX CDN
- [ ] Create index template (todo list with HTMX interactions)
- [ ] Create create todo form template
- [ ] Create edit todo form template
- [ ] Add responsive design and proper styling

### Main Application
- [ ] Set up HTTP server and routes
- [ ] Configure middleware (if needed)
- [ ] Initialize database connection
- [ ] Wire up all handlers and routes

### Testing & Documentation
- [ ] Test all CRUD operations
- [ ] Create README with setup instructions
- [ ] Add database setup instructions

## Technical Notes
- Use CDN links for Tailwind CSS and HTMX (no local files needed)
- Keep templates simple but functional
- Use HTMX for dynamic updates without full page reloads
- MySQL connection should be configurable via environment variables
- Error handling for database operations
- Basic form validation
