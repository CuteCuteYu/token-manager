# Token Manager - LLM Token Management System

A comprehensive REST API backend service for managing employee information and LLM token allocation. This system allows administrators to manage employees, issue tokens with quota limits, track token usage, and generate statistics.

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Project Structure](#project-structure)
5. [Prerequisites](#prerequisites)
6. [Installation](#installation)
7. [Running the Application](#running-the-application)
8. [API Documentation](#api-documentation)
9. [Frontend Setup](#frontend-setup)
10. [Data Storage](#data-storage)
11. [API Examples](#api-examples)
12. [Building and Testing](#building-and-testing)
13. [Code Style Guidelines](#code-style-guidelines)
14. [License](#license)

---

## Overview

Token Manager is a Go-based REST API system designed to manage LLM (Large Language Model) tokens within an organization. It provides a complete solution for:

- Employee registration and management
- Token issuance with configurable quotas
- Token usage tracking
- Usage statistics and analytics
- Token-employee relationship mapping

### Technology Stack

- **Backend**: Go 1.25+
- **Web Framework**: Gin
- **Data Storage**: JSON file (embedded)
- **Frontend**: Vanilla HTML/CSS/JavaScript

---

## Features

### Employee Management
- Create, read, update, and delete employee records
- Track employee department and position
- Store employee contact information

### Token Management
- Issue tokens with configurable quota limits
- Set token expiration dates (in days)
- Activate and revoke tokens
- Track token usage history

### Statistics and Analytics
- System-wide usage statistics
- Per-employee usage reports
- Token utilization percentages
- Active vs revoked token counts

### Token-Employee Mappings
- View all token-employee relationships
- Paginated mapping views
- Remaining quota tracking

---

## Architecture

### Design Patterns

The application follows several design patterns:

1. **Repository Pattern**: The `Store` type abstracts all data access
2. **Factory Pattern**: `NewStore()` acts as a factory for store initialization
3. **Handler Pattern**: HTTP handlers follow a consistent pattern for request processing
4. **DTO Pattern**: Request/Response DTOs separate API contract from internal models

### Thread Safety

All data access is protected by `sync.RWMutex`:
- Read operations use `RLock` for concurrent access
- Write operations use `Lock` for exclusive access
- All lock acquisitions use `defer` for automatic unlocking

### Data Flow

```
Client Request
    |
    v
Gin Router --> Handler --> Store Method --> DataStore --> JSON File
                    |
                    v
              Response DTO <-- Internal Models
```

---

## Project Structure

```
token-manager/
|
|-- [main.go](./main.go)          # Application entry point
|-- [handlers.go](./handlers.go)      # HTTP request handlers
|-- [models.go](./models.go)        # Data structures and DTOs
|-- [store.go](./store.go)         # Data persistence and business logic
|-- [data.json](./data.json)        # Data storage file (auto-created)
|-- [go.mod](./go.mod)           # Go module definition
|-- [go.sum](./go.sum)           # Go module checksums
|
|-- [index.html](./index.html)       # Frontend entry point
|-- css/
|   |-- [styles.css](./css/styles.css)   # All CSS styles
|
|-- js/
|   |-- [app.js](./js/app.js)       # Frontend JavaScript
|
|-- [README.md](./README.md)        # This file
|-- [README.zh.md](./README.zh.md)        # Chinese documentation
```

---

## Prerequisites

- **Go**: Version 1.25 or higher
- **Web Browser**: Modern browser with JavaScript enabled (for frontend)
- **Network Access**: Port 8080 available

---

## Installation

### 1. Clone or Download

```bash
# Navigate to project directory
cd token-manager
```

### 2. Install Dependencies

```bash
# Download Go module dependencies
go mod tidy
```

This will automatically download the Gin web framework and any other required packages.

### 3. Verify Installation

```bash
# Build the application
go build -o token-manager.exe
```

---

## Running the Application

### Starting the Server

```bash
# Run the application
go run main.go
```

The server will start on `http://localhost:8080`

### Expected Output

```
Server starting on :8080
```

### Data File Initialization

On first run, the application will automatically create a `data.json` file with an empty data structure:

```json
{
  "employees": [],
  "tokens": [],
  "usage_records": []
}
```

---

## API Documentation

### Base URL

```
http://localhost:8080/api
```

### Employee Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/employees | Create a new employee |
| GET | /api/employees | List all employees |
| GET | /api/employees/:id | Get employee by ID |
| PUT | /api/employees/:id | Update employee |
| DELETE | /api/employees/:id | Delete employee |

### Token Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/tokens/issue | Issue a new token |
| GET | /api/tokens | List all tokens |
| GET | /api/employees/:id/tokens | Get employee's tokens |
| POST | /api/tokens/:id/revoke | Revoke a token |
| POST | /api/tokens/use | Record token usage |
| GET | /api/tokens/:id/usage | Get token usage history |

### Statistics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/stats | Get system statistics |
| GET | /api/employees/:id/stats | Get employee statistics |

### Mapping Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/mappings | Get token-employee mappings |

### Query Parameters for /api/mappings

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number |
| page_size | integer | 10 | Items per page |

---

## Frontend Setup

### Running the Frontend

1. Ensure the backend server is running (`go run main.go`)
2. Open `index.html` in a web browser
3. The frontend will connect to `http://localhost:8080/api`

### Frontend Features

- **Overview Tab**: Dashboard with system statistics
- **Employees Tab**: Create and manage employees
- **Tokens Tab**: Issue and manage tokens
- **Mappings Tab**: View token-employee relationships

### Frontend File Structure

```
[index.html](./index.html)       - Main HTML structure
[css/styles.css](./css/styles.css)   - All styling (dark theme, responsive)
[js/app.js](./js/app.js)        - API communication and UI logic
```

---

## Data Storage

### Persistence

Data is stored in `data.json` using JSON format. The file is automatically:
- Created on first run
- Loaded on subsequent runs
- Updated after every write operation

### Data Schema

```json
{
  "employees": [
    {
      "id": "string",
      "name": "string",
      "department": "string",
      "email": "string",
      "position": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ],
  "tokens": [
    {
      "id": "string",
      "employee_id": "string",
      "token_value": "string",
      "total_quota": "integer",
      "used_quota": "integer",
      "is_active": "boolean",
      "issued_at": "timestamp",
      "expired_at": "timestamp",
      "revoked_at": "timestamp | null"
    }
  ],
  "usage_records": [
    {
      "id": "string",
      "token_id": "string",
      "used_at": "timestamp",
      "amount": "integer",
      "model": "string",
      "description": "string"
    }
  ]
}
```

---

## API Examples

### Create an Employee

```bash
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@company.com",
    "department": "Engineering",
    "position": "Software Engineer"
  }'
```

### Issue a Token

```bash
curl -X POST http://localhost:8080/api/tokens/issue \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "1234567890",
    "total_quota": 1000000,
    "days_valid": 30
  }'
```

### Use a Token

```bash
curl -X POST http://localhost:8080/api/tokens/use \
  -H "Content-Type: application/json" \
  -d '{
    "token_id": "token_123",
    "amount": 1500,
    "model": "gpt-4",
    "description": "Code review"
  }'
```

### Get System Statistics

```bash
curl http://localhost:8080/api/stats
```

### Revoke a Token

```bash
curl -X POST http://localhost:8080/api/tokens/token_123/revoke
```

---

## Building and Testing

### Building

```bash
# Build the application
go build -o token-manager.exe

# Build for different platforms
GOOS=windows GOARCH=amd64 go build -o token-manager.exe
GOOS=linux GOARCH=amd64 go build -o token-manager
```

### Running

```bash
# Run the built executable
./token-manager.exe

# Or run directly with Go
go run main.go
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# Run all tests (if any)
go test ./...

# Run a specific test
go test -run TestFunctionName ./...

# Run tests with verbose output
go test -v ./...
```

### Dependencies

The project uses a single external dependency:

- **github.com/gin-gonic/gin** (v1.12.0) - Web framework

View all dependencies in [`go.mod`](./go.mod).

---

## Code Style Guidelines

### Go Code Style

This project follows standard Go conventions:

1. **Import Organization**: Standard library first, then external packages
2. **Naming**: PascalCase for exported identifiers, camelCase for private
3. **Error Handling**: Return errors rather than panicking
4. **Documentation**: Comprehensive comments for all exported functions
5. **Formatting**: Automatic formatting with `go fmt`

### Frontend Code Style

1. **CSS**: CSS custom properties for theming
2. **JavaScript**: Modern ES6+ features
3. **Separation**: CSS and JS in separate files
4. **Comments**: JSDoc-style comments for functions

### File Organization

Each Go file has a specific purpose:

- [`main.go`](./main.go) - Application entry point
- [`handlers.go`](./handlers.go) - HTTP request handlers
- [`models.go`](./models.go) - Data structures
- [`store.go`](./store.go) - Data persistence and business logic

---

## License

This project is provided as-is for educational and internal use purposes.

---

## Additional Documentation

- [README.zh.md](./README.zh.md) - Chinese version of this documentation
