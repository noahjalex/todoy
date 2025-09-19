# Todoy - Todo List CRUD App

A simple Golang CRUD application for managing todo tasks with HTML templates, MySQL database, HTMX for dynamic interactions, and Tailwind CSS for styling.

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/noahjalex/todoy/releases):

- **Linux (x64)**: `todoy-linux-amd64`
- **Linux (ARM64)**: `todoy-linux-arm64`
- **macOS (Intel)**: `todoy-darwin-amd64`
- **macOS (Apple Silicon)**: `todoy-darwin-arm64`
- **Windows (x64)**: `todoy-windows-amd64.exe`

Make the binary executable (Linux/macOS):
```bash
chmod +x todoy-*
```

### Build from Source

Alternatively, build from source if you have Go installed:
```bash
git clone https://github.com/noahjalex/todoy.git
cd todoy
go build -o todoy main.go
```

## Features

- ✅ Create, Read, Update, Delete todos
- ✅ Toggle todo completion status
- ✅ Responsive design with Tailwind CSS
- ✅ Dynamic updates with HTMX (no page reloads)
- ✅ MySQL database storage
- ✅ Clean HTML templates
- ✅ Quick add form on home page

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: MySQL
- **Frontend**: HTML templates, HTMX, Tailwind CSS (CDN)
- **HTTP Router**: Standard library net/http

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
├── static/                # Static files (currently empty)
├── go.mod                 # Go module file
├── tasks.md               # Project planning document
└── README.md              # This file
```

## Prerequisites

- Go 1.19 or higher
- MySQL 5.7 or higher
- MySQL database named `todoapp` (or configure via environment variables)

## Database Setup

1. Create a MySQL database:
```sql
CREATE DATABASE todoapp;
```

2. The application will automatically create the required table on first run:
```sql
CREATE TABLE IF NOT EXISTS todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Installation & Setup

1. Clone or download the project
2. Install dependencies:
```bash
go mod tidy
```

3. Set environment variables (optional):
```bash
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=todoapp
export PORT=8080
```

4. Run the application:
```bash
go run main.go
```

5. Visit `http://localhost:8080` in your browser

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_USER | root | MySQL username |
| DB_PASSWORD | (empty) | MySQL password |
| DB_HOST | localhost | MySQL host |
| DB_PORT | 3306 | MySQL port |
| DB_NAME | todoapp | Database name |
| PORT | 8080 | Server port |

## API Routes

| Method | Route | Description |
|--------|-------|-------------|
| GET | / | Home page (list all todos) |
| GET | /todos/new | Show create todo form |
| POST | /todos | Create new todo |
| GET | /todos/{id}/edit | Show edit todo form |
| PUT | /todos/{id} | Update existing todo |
| DELETE | /todos/{id} | Delete todo |
| POST | /todos/{id}/toggle | Toggle todo completion status |

## Features

### HTMX Integration
- Dynamic todo list updates without page reloads
- Quick add form on home page
- Instant toggle completion status
- Confirmation dialogs for delete operations

### Responsive Design
- Mobile-friendly interface
- Clean, modern design with Tailwind CSS
- Visual indicators for completed todos
- Loading indicators for HTMX requests

### Database Operations
- Full CRUD operations
- Automatic timestamp tracking
- Error handling and validation
- Connection pooling

## Usage

1. **View Todos**: Visit the home page to see all todos
2. **Add Todo**: Use the quick add form or click "Add New Todo"
3. **Edit Todo**: Click the "Edit" button on any todo
4. **Complete Todo**: Click the "Complete" button to mark as done
5. **Delete Todo**: Click the "Delete" button (with confirmation)

## Development

To modify the application:

1. **Models**: Edit `models/todo.go` for database operations
2. **Handlers**: Edit `handlers/todo_handlers.go` for HTTP logic
3. **Templates**: Edit files in `templates/` for UI changes
4. **Styling**: Modify Tailwind classes in templates or add custom CSS in `static/`

## Building for Production

```bash
# Build binary
go build -o todoapp main.go

# Run binary
./todoapp
```

## Releases

This project uses semantic versioning and automated releases via GitHub Actions.

### Creating a Release

1. Update the `CHANGELOG.md` file with your changes
2. Use the release script:
   ```bash
   ./scripts/release.sh v1.0.0
   ```
3. The GitHub Action will automatically build binaries for all platforms and create a release

### Version Management

- Version information is embedded in the binary during build
- The application displays the version on startup
- Releases follow semantic versioning (MAJOR.MINOR.PATCH)

## License

This is a sample project for educational purposes.
