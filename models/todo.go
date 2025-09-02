package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoModel struct {
	DB *sql.DB
}

// NewTodoModel creates a new TodoModel instance
func NewTodoModel(db *sql.DB) *TodoModel {
	return &TodoModel{DB: db}
}

// CreateTable creates the todos table if it doesn't exist
func (m *TodoModel) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := m.DB.Exec(query)
	return err
}

// GetAllTodos retrieves all todos from the database
func (m *TodoModel) GetAllTodos() ([]Todo, error) {
	query := "SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodoByID retrieves a specific todo by ID
func (m *TodoModel) GetTodoByID(id int) (*Todo, error) {
	query := "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?"

	var todo Todo
	err := m.DB.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, err
	}

	return &todo, nil
}

// CreateTodo creates a new todo
func (m *TodoModel) CreateTodo(title, description string) (*Todo, error) {
	query := "INSERT INTO todos (title, description) VALUES (?, ?)"

	result, err := m.DB.Exec(query, title, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.GetTodoByID(int(id))
}

// UpdateTodo updates an existing todo
func (m *TodoModel) UpdateTodo(id int, title, description string) (*Todo, error) {
	query := "UPDATE todos SET title = ?, description = ? WHERE id = ?"

	_, err := m.DB.Exec(query, title, description, id)
	if err != nil {
		return nil, err
	}

	return m.GetTodoByID(id)
}

// DeleteTodo deletes a todo by ID
func (m *TodoModel) DeleteTodo(id int) error {
	query := "DELETE FROM todos WHERE id = ?"

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}

// ToggleTodoStatus toggles the completed status of a todo
func (m *TodoModel) ToggleTodoStatus(id int) (*Todo, error) {
	query := "UPDATE todos SET completed = NOT completed WHERE id = ?"

	_, err := m.DB.Exec(query, id)
	if err != nil {
		return nil, err
	}

	return m.GetTodoByID(id)
}
