# Gin Assistant v2

A powerful CLI tool for generating and managing Go Gin framework projects with hot reload functionality and automatic service scaffolding.

## Features

- 🚀 **Project Initialization**: Quick setup of Gin projects with predefined structure
- 🔧 **Service Generation**: Automatic creation of service layers with repository pattern
- 🔄 **Hot Reload**: Development server with automatic restart on file changes
- 📁 **Project Structure**: Organized folder structure following best practices
- 🔌 **Database Integration**: Pre-configured GORM setup for MySQL
- 🌐 **CORS Support**: Built-in CORS configuration

## Installation

```bash
go install github.com/gin-assistantv2@latest
```

## Commands

### Initialize Project
Creates a new Gin project in the current directory with the complete folder structure:

```bash
gin-assistant2 init
```

This command will:
- Initialize Go module
- Install required dependencies
- Create project structure
- Generate base files

### Generate Service
Creates a new service with repository, service, and router files:

```bash
gin-assistant2 service <service_name>
```

Example:
```bash
gin-assistant2 service user
gin-assistant2 service product
```

### Run Development Server
Starts the development server with hot reload on port 8080 (default) or specified port:

```bash
gin-assistant2 run [port]
```

Examples:
```bash
gin-assistant2 run          # Runs on port 8080
gin-assistant2 run 3000     # Runs on port 3000
gin-assistant2 run 8000     # Runs on port 8000
```

**Port Requirements:**
- Port must be between 7000-9999
- Uses environment variable `PORT` if set
- Automatically frees port if in use

### Help
Display help information:

```bash
gin-assistant2 help
```

## Project Structure

After initialization, your project will have the following structure:

```
your-project/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Dependencies checksum
└── app/
    ├── run.go              # Server configuration
    ├── db/
    │   └── connection.go   # Database connection
    ├── helper/
    │   └── helper.go       # Utility functions
    ├── model/              # Database models
    └── service/            # Service modules
        ├── validator.go    # Input validation
        ├── clonestruct.go  # Struct utilities
        └── base.go         # Base service functions
```

## Generated Service Structure

When you create a service using `gin-assistant2 service <name>`, it generates:

```
app/service/<service_name>/
├── repository.go           # Data access layer
├── service.go             # Business logic layer
└── router.go              # HTTP route handlers
```

### Repository Pattern
Each service follows the repository pattern:

- **Repository**: Handles database operations
- **Service**: Contains business logic
- **Router**: Manages HTTP routes and handlers

## Dependencies

The tool automatically installs these dependencies:

- `github.com/gin-gonic/gin` - Web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `gorm.io/gorm` - ORM library
- `gorm.io/driver/mysql` - MySQL driver
- `github.com/joho/godotenv` - Environment variables
- `gorm.io/gorm/logger` - GORM logging

## Development Features

### Hot Reload
The development server includes:
- Automatic file watching
- Instant server restart on changes
- Cross-platform compatibility (Windows/Unix)
- Debounced restarts to prevent rapid firing

### Port Management
- Automatic port conflict detection
- Smart port freeing mechanism
- Environment variable support
- Configurable port range (7000-9999)

## Environment Variables

You can configure the application using environment variables:

```bash
# .env file
PORT=8080
```

## Usage Examples

### 1. Create a new project
```bash
mkdir my-api
cd my-api
gin-assistant2 init
```

### 2. Generate user service
```bash
gin-assistant2 service user
```

### 3. Start development server
```bash
gin-assistant2 run 8080
```

### 4. Your main.go will look like:
```go
package main

import "your-project/app"

func main() {
    app.Run()
}
```

## Cross-Platform Support

The tool supports:
- **Windows**: Uses `cmd` and `taskkill` for process management
- **Linux/macOS**: Uses `sh` and `kill` for process management
- **Path handling**: Automatic path separator conversion

## File Watching

The hot reload feature watches:
- All directories in the `app` folder
- Automatically adds new directories to the watcher
- Handles file create, modify, and delete events
- Uses debounced restart (1-second delay) to prevent rapid restarts

## Error Handling

The tool includes comprehensive error handling for:
- Port conflicts and resolution
- File system operations
- Process management
- Module initialization

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you encounter any issues or have questions:
1. Check the help command: `gin-assistant2 help`
2. Ensure you're in a valid Go project directory
3. Verify port availability (7000-9999 range)
4. Check file permissions for directory creation

## Version

Current version: v2.0.0

---

**Happy coding with Gin Assistant v2! 🚀**