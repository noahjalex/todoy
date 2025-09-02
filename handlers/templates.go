package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"
	"todoy/models"
)

// TemplateData represents the common data structure for all templates
type TemplateData struct {
	Title       string
	Todo        *models.Todo
	Todos       []models.Todo
	Error       string
	Success     string
	CurrentPath string
}

// TemplateManager handles all template operations
type TemplateManager struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

// NewTemplateManager creates a new template manager with common functions
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	// Define common template functions
	tm.funcMap = template.FuncMap{
		"datesEqual": func(t1, t2 time.Time) bool {
			return t1.Equal(t2)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006 15:04")
		},
		"formatDateLong": func(t time.Time) string {
			return t.Format("January 2, 2006 at 3:04 PM")
		},
		"isCompleted": func(completed bool) string {
			if completed {
				return "Completed"
			}
			return "Pending"
		},
		"completedClass": func(completed bool) string {
			if completed {
				return "border-green-500"
			}
			return "border-blue-500"
		},
		"completedTextClass": func(completed bool) string {
			if completed {
				return "line-through text-gray-500"
			}
			return ""
		},
		"toggleButtonClass": func(completed bool) string {
			if completed {
				return "bg-yellow-500 hover:bg-yellow-600 text-white"
			}
			return "bg-green-500 hover:bg-green-600 text-white"
		},
		"toggleButtonText": func(completed bool) string {
			if completed {
				return "Undo"
			}
			return "Complete"
		},
		"toggleButtonLongText": func(completed bool) string {
			if completed {
				return "Mark as Pending"
			}
			return "Mark as Complete"
		},
		"statusClass": func(completed bool) string {
			if completed {
				return "text-green-600"
			}
			return "text-blue-600"
		},
	}

	tm.loadTemplates()
	return tm
}

// loadTemplates loads and parses all templates
func (tm *TemplateManager) loadTemplates() {
	var err error

	// Parse templates with the common function map
	tm.templates["index"], err = template.New("layout.html").Funcs(tm.funcMap).ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		panic("Error loading index template: " + err.Error())
	}

	tm.templates["create"], err = template.New("layout.html").Funcs(tm.funcMap).ParseFiles("templates/layout.html", "templates/create.html")
	if err != nil {
		panic("Error loading create template: " + err.Error())
	}

	tm.templates["edit"], err = template.New("layout.html").Funcs(tm.funcMap).ParseFiles("templates/layout.html", "templates/edit.html")
	if err != nil {
		panic("Error loading edit template: " + err.Error())
	}
}

// Render executes a template with the given data
func (tm *TemplateManager) Render(w http.ResponseWriter, templateName string, data TemplateData) error {
	tmpl, exists := tm.templates[templateName]
	if !exists {
		log.Printf("Template '%s' not found. Available templates: %v", templateName, tm.getTemplateNames())
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return nil
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", templateName, err)
		return err
	}

	return nil
}

// getTemplateNames returns a list of available template names for debugging
func (tm *TemplateManager) getTemplateNames() []string {
	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// RenderTodoList renders just the todo list portion for HTMX requests
func (tm *TemplateManager) RenderTodoList(w http.ResponseWriter, todos []models.Todo) error {
	tmpl := `{{range .}}
	<div class="bg-white p-4 rounded-lg shadow-md border-l-4 {{completedClass .Completed}}">
		<div class="flex items-center justify-between">
			<div class="flex-1">
				<h3 class="font-semibold {{completedTextClass .Completed}}">{{.Title}}</h3>
				{{if .Description}}<p class="text-gray-600 mt-1 {{completedTextClass .Completed}}">{{.Description}}</p>{{end}}
				<p class="text-xs text-gray-400 mt-2">Created: {{formatDate .CreatedAt}}</p>
			</div>
			<div class="flex space-x-2 ml-4">
				<button hx-post="/todos/{{.ID}}/toggle" hx-target="#todo-list" hx-indicator="#loading"
						class="px-3 py-1 text-sm rounded transition duration-200 {{toggleButtonClass .Completed}}">
					{{toggleButtonText .Completed}}
				</button>
				<a href="/todos/{{.ID}}/edit" 
				   class="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 transition duration-200">
					Edit
				</a>
				<button hx-delete="/todos/{{.ID}}" hx-target="#todo-list" hx-confirm="Are you sure you want to delete this todo?" hx-indicator="#loading"
						class="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600 transition duration-200">
					Delete
				</button>
			</div>
		</div>
	</div>
	{{end}}`

	t := template.Must(template.New("todos").Funcs(tm.funcMap).Parse(tmpl))
	return t.Execute(w, todos)
}

// RenderError renders an error response
func (tm *TemplateManager) RenderError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// RenderJSON renders a JSON response (for future API endpoints)
func (tm *TemplateManager) RenderJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	// For now, just return success - can be extended later
	w.Write([]byte(`{"success": true}`))
	return nil
}
