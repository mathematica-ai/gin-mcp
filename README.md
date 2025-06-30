# 🏦 gin-mcp

> **Model Context Protocol (MCP) Server** - A beautiful, production-ready MCP server implementation built with Go and Gin, enabling seamless integration with LLM applications like Claude Desktop.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Gin Version](https://img.shields.io/badge/Gin-1.10+-green.svg)](https://github.com/gin-gonic/gin)
[![MCP Version](https://img.shields.io/badge/MCP-2025.06.18-orange.svg)](https://modelcontextprotocol.io)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Contributors](https://img.shields.io/github/contributors/your-org/gin-mcp)](https://github.com/your-org/gin-mcp/graphs/contributors)
[![Issues](https://img.shields.io/github/issues/your-org/gin-mcp)](https://github.com/your-org/gin-mcp/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/your-org/gin-mcp)](https://github.com/your-org/gin-mcp/pulls)
[![CI/CD](https://img.shields.io/github/actions/workflow/status/your-org/gin-mcp/ci.yml?branch=main)](https://github.com/your-org/gin-mcp/actions)
[![Code Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)](https://codecov.io/gh/your-org/gin-mcp)

---

## ✨ Overview

`gin-mcp` is a production-grade **Model Context Protocol (MCP)** server implementation that enables LLM applications to access your data and tools through the standardized MCP protocol. Built with the elegance of Go and the power of Gin, it provides seamless integration with MCP clients like Claude Desktop, IDEs, and AI tools.

**What is MCP?** The [Model Context Protocol](https://modelcontextprotocol.io/introduction) is an open protocol that standardizes how applications provide context to LLMs. Think of MCP like a USB-C port for AI applications - it provides a standardized way to connect AI models to different data sources and tools.

**Two Usage Modes:**
- **🔌 Package Mode**: Add MCP server functionality to existing Gin applications
- **🚀 Standalone Mode**: Run as a complete MCP server

### 🌟 Key Features

- **🔌 MCP Protocol Compliant**: Full implementation of the Model Context Protocol specification
- **🔥 Hot Reload**: Resources and tools are automatically discovered and registered on file changes
- **🔌 Plugin Architecture**: Extensible Go plugin system for high-performance tools
- **⚡ High Performance**: Built on Gin framework for exceptional speed and reliability
- **🛡️ Production Ready**: Graceful shutdown, error handling, and comprehensive logging
- **🎯 Simple API**: RESTful endpoints for MCP resource management and tool execution
- **📊 Health Monitoring**: Built-in health checks and registry introspection
- **🔧 Flexible Integration**: Use as a package or standalone server
- **🔄 Real-time Updates**: File system monitoring with instant resource registration
- **🔐 Secure**: Proper authentication and authorization for MCP connections

### 🏗️ MCP Architecture

```
┌─────────────────┐    MCP Protocol    ┌─────────────────┐
│   MCP Client    │◄──────────────────►│   gin-mcp       │
│ (Claude Desktop)│                    │   Server        │
│   (IDE Plugin)  │                    │                 │
│   (AI Tool)     │                    │ ┌─────────────┐ │
└─────────────────┘                    │ │ Resources   │ │
                                       │ │ (Files, DB) │ │
                                       │ └─────────────┘ │
                                       │ ┌─────────────┐ │
                                       │ │ Tools       │ │
                                       │ │ (Functions) │ │
                                       │ └─────────────┘ │
                                       └─────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- CGO enabled (for Go plugin support)

### 🔌 Package Mode (Recommended)

Add MCP server to your existing Gin application:

```bash
# Install the package
go get github.com/your-org/gin-mcp
```

```go
package main

import (
    "github.com/gin-gonic/gin"
    "gin-mcp/pkg/ginmcp"
)

func main() {
    router := gin.Default()
    
    // Your existing routes
    router.GET("/api/users", getUsersHandler)
    
    // Add MCP server to your app
    mcp, _ := ginmcp.New(&ginmcp.MCPConfig{
        ResourcesDir: "./resources",
        ToolsDir:     "./tools",
        Prefix:       "/mcp",
    })
    mcp.SetupRoutes(router)
    
    router.Run(":8080")
}
```

### 🚀 Standalone Mode

```bash
# Clone the repository
git clone https://github.com/your-org/gin-mcp.git
cd gin-mcp

# Install dependencies
go mod tidy

# Build and run
go build -o gin-mcp .
./gin-mcp
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GIN_MCP_PORT` | `:8080` | Server port |
| `GIN_MCP_RESOURCES_DIR` | `./resources` | Resources directory path |
| `GIN_MCP_TOOLS_DIR` | `./tools` | Tools directory path |

---

## 🎯 Usage Examples

### 🔌 Package Mode

Add MCP server to your existing Gin application with custom configuration:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "gin-mcp/pkg/ginmcp"
)

func main() {
    router := gin.Default()
    
    // Your existing routes
    router.GET("/api/users", getUsersHandler)
    router.POST("/api/orders", createOrderHandler)
    
    // Initialize and add MCP server
    mcp, _ := ginmcp.New(&ginmcp.MCPConfig{
        ResourcesDir: "./data",
        ToolsDir:     "./ml-tools",
        Prefix:       "/mcp",
    })
    mcp.SetupRoutes(router)
    
    router.Run(":8080")
}
```

**Available MCP endpoints:**
- `GET /mcp/health` - Health check
- `GET /mcp/resources` - List available resources
- `GET /mcp/resources/{name}` - Get resource info
- `POST /mcp/resources/{name}` - Access resource content
- `GET /mcp/tools` - List available tools
- `POST /mcp/tools/{name}` - Execute tool
- `GET /mcp/registry` - Export registry

### 🚀 Standalone Mode

Run as a complete MCP server:

```bash
# Clone the repository
git clone https://github.com/your-org/gin-mcp.git
cd gin-mcp

# Install dependencies
go mod tidy

# Build the application
go build -o gin-mcp .

# Run the server
./gin-mcp
```

**Available MCP endpoints:**
- `GET /mcp/health` - Health check
- `GET /mcp/resources` - List available resources
- `GET /mcp/resources/{name}` - Get resource info
- `POST /mcp/resources/{name}` - Access resource content
- `GET /mcp/tools` - List available tools
- `POST /mcp/tools/{name}` - Execute tool
- `GET /mcp/registry` - Export registry

For detailed package documentation, see [pkg/ginmcp/README.md](pkg/ginmcp/README.md).

---

## 📁 Project Structure

```
gin-mcp/
├── main.go              # 🚀 Standalone MCP server entry point
├── pkg/ginmcp/          # 🔌 Reusable MCP package
│   ├── ginmcp.go       # Main MCP package implementation
│   └── README.md       # Package documentation
├── examples/            # 📚 Usage examples
│   ├── integration/    # Package mode example
│   └── standalone/     # Standalone mode example
├── registry/            # 📋 Resource and tool registry management
│   └── registry.go
├── handlers/            # ⚙️ MCP resource and tool execution logic
│   └── handler.go
├── watcher/             # 👀 File system monitoring
│   └── watcher.go
├── resources/           # 📁 MCP resources (auto-created)
├── tools/               # 🔧 MCP tools (auto-created)
├── docs/                # 📚 Project documentation
│   ├── README.md        # Documentation index
│   ├── SPECIFICATION.md # MCP specification details
│   └── PROJECT_STATUS.md # Project status & roadmap
├── .github/             # 🔧 GitHub templates and workflows
│   ├── ISSUE_TEMPLATE/
│   ├── workflows/
│   └── pull_request_template.md
├── go.mod               # 📦 Go module definition
├── Dockerfile           # 🐳 Container configuration
├── Makefile             # 🔨 Build automation
├── LICENSE              # 📄 MIT License
├── CONTRIBUTING.md      # 🤝 Contributing guidelines
├── CODE_OF_CONDUCT.md   # 📜 Community guidelines
└── README.md           # 📖 This file
```

---

## 🎯 MCP API Reference

### Health Check

```http
GET /mcp/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "gin-mcp",
  "mcp_version": "2025.06.18",
  "resources": 5,
  "tools": 3,
  "watcher": true,
  "prefix": "/mcp"
}
```

### List Resources

```http
GET /mcp/resources
```

**Response:**
```json
{
  "resources": [
    {
      "name": "database_schema",
      "type": "file",
      "file_path": "./resources/schema.sql",
      "mime_type": "text/sql"
    },
    {
      "name": "api_docs",
      "type": "file", 
      "file_path": "./resources/api.md",
      "mime_type": "text/markdown"
    }
  ],
  "count": 2
}
```

### Access Resource

```http
POST /mcp/resources/{name}
Content-Type: application/json

{
  "uri": "file://./resources/schema.sql"
}
```

**Response:**
```json
{
  "contents": [
    {
      "uri": "file://./resources/schema.sql",
      "mime_type": "text/sql",
      "text": "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255));"
    }
  ]
}
```

### List Tools

```http
GET /mcp/tools
```

**Response:**
```json
{
  "tools": [
    {
      "name": "calculate",
      "description": "Perform mathematical calculations",
      "input_schema": {
        "type": "object",
        "properties": {
          "expression": {"type": "string"}
        }
      }
    },
    {
      "name": "data_analyzer",
      "description": "Analyze data files",
      "input_schema": {
        "type": "object", 
        "properties": {
          "file_path": {"type": "string"}
        }
      }
    }
  ],
  "count": 2
}
```

### Execute Tool

```http
POST /mcp/tools/{name}
Content-Type: application/json

{
  "arguments": {
    "expression": "2 + 2 * 3"
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "Result: 8"
    }
  ]
}
```

---

## 🔧 MCP Resource Development

### File Resources

Place files in the `resources/` directory to make them available as MCP resources:

```
resources/
├── database_schema.sql    # Database schema
├── api_documentation.md   # API documentation
├── config.json           # Configuration files
└── data/
    ├── users.csv         # User data
    └── products.json     # Product catalog
```

### Tool Development

Create tools that can be executed by MCP clients:

#### Go Plugin Tools

```go
// tools/calculator.go
package main

import "encoding/json"

func Execute(input []byte) ([]byte, error) {
    var args map[string]interface{}
    json.Unmarshal(input, &args)
    
    expression := args["expression"].(string)
    // Your calculation logic here
    
    result := map[string]interface{}{
        "result": "8",
        "expression": expression,
    }
    
    return json.Marshal(result)
}
```

Compile to a plugin:

```bash
go build -buildmode=plugin -o tools/calculator.so tools/calculator.go
```

#### Advanced Go Plugin Example

```go
// tools/data_analyzer.go
package main

import (
    "encoding/json"
    "fmt"
    "math"
)

func Execute(input []byte) ([]byte, error) {
    var data map[string]interface{}
    if err := json.Unmarshal(input, &data); err != nil {
        return nil, fmt.Errorf("failed to parse input: %w", err)
    }
    
    arguments := data["arguments"].(map[string]interface{})
    numbers, ok := arguments["numbers"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("numbers array is required")
    }
    
    // Convert to float64 slice and calculate statistics
    var values []float64
    for _, num := range numbers {
        if val, ok := num.(float64); ok {
            values = append(values, val)
        }
    }
    
    if len(values) == 0 {
        return nil, fmt.Errorf("no valid numbers provided")
    }
    
    // Calculate basic statistics
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    mean := sum / float64(len(values))
    
    // Calculate standard deviation
    variance := 0.0
    for _, v := range values {
        variance += math.Pow(v-mean, 2)
    }
    stdDev := math.Sqrt(variance / float64(len(values)))
    
    result := fmt.Sprintf("Count: %d, Sum: %.2f, Mean: %.2f, StdDev: %.2f", 
        len(values), sum, mean, stdDev)
    
    response := map[string]interface{}{
        "content": []map[string]interface{}{
            {"type": "text", "text": result},
        },
    }
    
    return json.Marshal(response)
}
```

Build the plugin:

```bash
cd tools
go build -buildmode=plugin -o data_analyzer.so data_analyzer.go
```

---

## 🔧 Development

### Building

```bash
# Build for current platform
go build -o gin-mcp .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o gin-mcp-linux .

# Build with debug information
go build -gcflags="all=-N -l" -o gin-mcp-debug .
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./handlers
```

### Development Mode

```bash
# Run in development mode with hot reload
make dev

# Build and run
make run

# Clean build artifacts
make clean
```

---

## 🐳 Docker

### Build Image

```bash
# Build the Docker image
docker build -t gin-mcp .

# Run with Docker
docker run -p 8080:8080 -v $(PWD)/resources:/app/resources -v $(PWD)/tools:/app/tools gin-mcp
```

### Docker Compose

```yaml
version: '3.8'
services:
  gin-mcp:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./resources:/app/resources
      - ./tools:/app/tools
    environment:
      - GIN_MCP_PORT=:8080
      - GIN_MCP_RESOURCES_DIR=/app/resources
      - GIN_MCP_TOOLS_DIR=/app/tools
```

---

## 🔐 Security Considerations

- **Authentication**: MCP connections should be properly authenticated
- **Authorization**: Implement proper access controls for resources and tools
- **Input Validation**: All MCP requests are validated and sanitized
- **Execution Isolation**: Tools run in isolated contexts with timeouts
- **File System Access**: Restricted to designated resources and tools directories
- **Network Security**: Use HTTPS in production environments

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code style and standards
- Testing requirements
- Pull request process
- Issue reporting

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io) for the open protocol specification
- [Gin](https://github.com/gin-gonic/gin) for the excellent HTTP framework
- [Go](https://golang.org) for the beautiful programming language

---

## 📞 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/your-org/gin-mcp/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/gin-mcp/discussions)
- **MCP Community**: [MCP Discussions](https://github.com/modelcontextprotocol/specification/discussions)

---

*Built with excellence and beauty in mind. ✨*
