package handlers

import (
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"todoy/models"
)

type TodoHandler struct {
	TodoModel       *models.TodoModel
	TemplateManager *TemplateManager
}

// NewTodoHandler creates a new TodoHandler instance
func NewTodoHandler(todoModel *models.TodoModel, templateFS fs.FS) *TodoHandler {
	return &TodoHandler{
		TodoModel:       todoModel,
		TemplateManager: NewTemplateManager(templateFS),
	}
}

// HomeHandler displays all todos
func (h *TodoHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := h.TodoModel.GetAllTodos()
	if err != nil {
		h.TemplateManager.RenderError(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Title:       "Todo List",
		Todos:       todos,
		CurrentPath: r.URL.Path,
	}

	err = h.TemplateManager.Render(w, "index", data)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// CreateTodoFormHandler shows the create todo form
func (h *TodoHandler) CreateTodoFormHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:       "Create Todo",
		CurrentPath: r.URL.Path,
	}

	err := h.TemplateManager.Render(w, "create", data)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// CreateTodoHandler creates a new todo
func (h *TodoHandler) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.TemplateManager.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.TemplateManager.RenderError(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))

	if title == "" {
		h.TemplateManager.RenderError(w, "Title is required", http.StatusBadRequest)
		return
	}

	_, err = h.TodoModel.CreateTodo(title, description)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error creating todo", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return updated todo list for HTMX
		todos, err := h.TodoModel.GetAllTodos()
		if err != nil {
			h.TemplateManager.RenderError(w, "Error fetching todos", http.StatusInternalServerError)
			return
		}

		err = h.TemplateManager.RenderTodoList(w, todos)
		if err != nil {
			h.TemplateManager.RenderError(w, "Error rendering todos", http.StatusInternalServerError)
		}
		return
	}

	// Regular form submission - redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditTodoFormHandler shows the edit todo form
func (h *TodoHandler) EditTodoFormHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	idStr := strings.TrimSuffix(path, "/edit")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.TemplateManager.RenderError(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.TodoModel.GetTodoByID(id)
	if err != nil {
		h.TemplateManager.RenderError(w, "Todo not found", http.StatusNotFound)
		return
	}

	data := TemplateData{
		Title:       "Edit Todo",
		Todo:        todo,
		CurrentPath: r.URL.Path,
	}

	err = h.TemplateManager.Render(w, "edit", data)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// UpdateTodoHandler updates an existing todo
func (h *TodoHandler) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		h.TemplateManager.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	idStr := path

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.TemplateManager.RenderError(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.TemplateManager.RenderError(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))

	if title == "" {
		h.TemplateManager.RenderError(w, "Title is required", http.StatusBadRequest)
		return
	}

	_, err = h.TodoModel.UpdateTodo(id, title, description)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error updating todo", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteTodoHandler deletes a todo
func (h *TodoHandler) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.TemplateManager.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	idStr := path

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.TemplateManager.RenderError(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = h.TodoModel.DeleteTodo(id)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error deleting todo", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return updated todo list for HTMX
		todos, err := h.TodoModel.GetAllTodos()
		if err != nil {
			h.TemplateManager.RenderError(w, "Error fetching todos", http.StatusInternalServerError)
			return
		}

		err = h.TemplateManager.RenderTodoList(w, todos)
		if err != nil {
			h.TemplateManager.RenderError(w, "Error rendering todos", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ToggleTodoHandler toggles the completion status of a todo
func (h *TodoHandler) ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.TemplateManager.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	idStr := strings.TrimSuffix(path, "/toggle")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.TemplateManager.RenderError(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	_, err = h.TodoModel.ToggleTodoStatus(id)
	if err != nil {
		h.TemplateManager.RenderError(w, "Error toggling todo status", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return updated todo list for HTMX
		todos, err := h.TodoModel.GetAllTodos()
		if err != nil {
			h.TemplateManager.RenderError(w, "Error fetching todos", http.StatusInternalServerError)
			return
		}

		err = h.TemplateManager.RenderTodoList(w, todos)
		if err != nil {
			h.TemplateManager.RenderError(w, "Error rendering todos", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
