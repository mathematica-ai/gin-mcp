# 🏦 gin-mcp Package

The `gin-mcp` package provides a reusable Gin middleware that adds **Model Context Protocol (MCP)** server functionality to your existing Gin applications.

## 🚀 Quick Start

### Installation

```bash
go get github.com/your-org/gin-mcp
```

### Basic Usage

```go
package main

import (
    "github.com/gin-gonic/gin"
    "gin-mcp/pkg/ginmcp"
)

func main() {
    // Create your existing Gin router
    router := gin.Default()
    
    // Add your existing routes
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello World"})
    })
    
    // Initialize MCP server
    mcp, err := ginmcp.New(&ginmcp.MCPConfig{
        ResourcesDir: "./resources",
        ToolsDir:     "./tools",
        Prefix:       "/mcp",
    })
    if err != nil {
        panic(err)
    }
    
    // Add MCP server routes to your router
    if err := mcp.SetupRoutes(router); err != nil {
        panic(err)
    }
    
    // Start your server
    router.Run(":8080")
}
```

## 📋 Configuration

### MCPConfig

```go
type MCPConfig struct {
    ResourcesDir string // Directory for MCP resources (default: "./resources")
    ToolsDir     string // Directory for MCP tools (default: "./tools")
    Prefix       string // URL prefix for MCP endpoints (default: "/mcp")
    Port         string // Port for standalone server (default: ":8080")
}
```

### Default Configuration

```go
config := ginmcp.DefaultConfig()
// Returns:
// - ResourcesDir: "./resources"
// - ToolsDir: "./tools"
// - Prefix: "/mcp"
// - Port: ":8080"
```

## 🎯 Usage Modes

### 1. Integration Mode (Recommended)

Add MCP server to your existing Gin application:

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
- `GET /mcp/tools/{name}` - Get tool info
- `POST /mcp/tools/{name}` - Execute tool
- `GET /mcp/registry` - Export registry

### 2. Standalone Mode

Run MCP server as a standalone server:

```go
package main

import (
    "gin-mcp/pkg/ginmcp"
)

func main() {
    mcp, _ := ginmcp.New(nil) // Use default config
    mcp.StartStandalone()
}
```

## 🔧 Advanced Usage

### Custom Configuration

```go
config := &ginmcp.MCPConfig{
    ResourcesDir: "/path/to/your/resources",
    ToolsDir:     "/path/to/your/tools",
    Prefix:       "/ai",  // Custom prefix
}

mcp, err := ginmcp.New(config)
```

### Direct Access to Components

```go
mcp, _ := ginmcp.New(nil)

// Access the registry directly
registry := mcp.GetRegistry()
resources := registry.ListResources()
tools := registry.ListTools()

// Access the handler directly
handler := mcp.GetHandler()
```

### Graceful Shutdown

```go
package main

import (
    "os"
    "os/signal"
    "syscall"
    "gin-mcp/pkg/ginmcp"
)

func main() {
    mcp, _ := ginmcp.New(nil)
    
    // Setup graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        mcp.Stop() // Gracefully stop MCP server
        os.Exit(0)
    }()
    
    mcp.StartStandalone()
}
```

## 📁 Directory Structure

```
your-project/
├── main.go
├── resources/           # MCP resources directory
│   ├── schema.sql      # Database schema
│   ├── api_docs.md     # API documentation
│   └── config.json     # Configuration files
├── tools/              # MCP tools directory
│   ├── calculator.so   # Go plugin (compiled)
│   ├── analyzer.so     # Compiled Go plugin
│   └── processor.so    # Another Go plugin
└── ...
```

## 🔍 API Endpoints

When integrated, MCP server provides these endpoints under your configured prefix:

### Health Check
```http
GET /mcp/health
```
```json
{
  "status": "healthy",
  "service": "gin-mcp",
  "mcp_version": "2025.06.18",
  "resources": 3,
  "tools": 2,
  "watcher": true,
  "prefix": "/mcp"
}
```

### List Resources
```http
GET /mcp/resources
```
```json
{
  "resources": [
    {
      "name": "database_schema",
      "type": "file",
      "file_path": "./resources/schema.sql",
      "mime_type": "text/sql"
    }
  ],
  "count": 1
}
```

### Access Resource
```http
POST /mcp/resources/database_schema
Content-Type: application/json

{
  "uri": "file://./resources/schema.sql"
}
```
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
```json
{
  "tools": [
    {
      "name": "calculator",
      "description": "MCP tool: calculator",
      "type": "go_plugin",
      "file_path": "./tools/calculator.so",
      "input_schema": {
        "type": "object",
        "properties": {
          "arguments": {
            "type": "object",
            "description": "Tool arguments"
          }
        }
      }
    }
  ],
  "count": 1
}
```

### Execute Tool
```http
POST /mcp/tools/calculator
Content-Type: application/json

{
  "arguments": {
    "expression": "2 + 3 * 4"
  }
}
```
```json
{
  "content": [
    {
      "type": "text",
      "text": "Result: 14"
    }
  ]
}
```

## 🔌 MCP Resources

MCP resources are files that can be accessed by MCP clients. Supported formats:

- **SQL files** (`.sql`) - Database schemas and queries
- **Markdown files** (`.md`) - Documentation
- **JSON files** (`.json`) - Configuration and data
- **Text files** (`.txt`) - Plain text content
- **CSV files** (`.csv`) - Tabular data
- **YAML files** (`.yaml`, `.yml`) - Configuration
- **XML files** (`.xml`) - Structured data

## 🔧 MCP Tools

MCP tools are executable functions that can be called by MCP clients:

### Go Plugin Tools

Create a Go file with an `Execute` function:

```go
// tools/calculator.go
package main

import "encoding/json"

func Execute(input []byte) ([]byte, error) {
    var args map[string]interface{}
    json.Unmarshal(input, &args)
    
    expression := args["arguments"].(map[string]interface{})["expression"].(string)
    // Your calculation logic here
    
    result := map[string]interface{}{
        "content": []map[string]interface{}{
            {
                "type": "text",
                "text": fmt.Sprintf("Result: %v", result),
            },
        },
    }
    
    return json.Marshal(result)
}
```

Compile to a plugin:

```bash
go build -buildmode=plugin -o tools/calculator.so tools/calculator.go
```

### Advanced Go Plugin Tools

Create more sophisticated Go plugins with complex functionality:

```go
// tools/file_analyzer.go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func Execute(input []byte) ([]byte, error) {
    var data map[string]interface{}
    if err := json.Unmarshal(input, &data); err != nil {
        return nil, fmt.Errorf("failed to parse input: %w", err)
    }
    
    arguments := data["arguments"].(map[string]interface{})
    filePath, ok := arguments["file_path"].(string)
    if !ok {
        return nil, fmt.Errorf("file_path argument is required")
    }
    
    // Get file info
    info, err := os.Stat(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to get file info: %w", err)
    }
    
    // Read file content for analysis
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    // Perform analysis
    analysis := map[string]interface{}{
        "file_name":  info.Name(),
        "file_size":  info.Size(),
        "extension":  strings.ToLower(filepath.Ext(filePath)),
        "line_count": strings.Count(string(content), "\n") + 1,
        "char_count": len(content),
        "word_count": len(strings.Fields(string(content))),
    }
    
    resultText := fmt.Sprintf("File Analysis:\n- Name: %s\n- Size: %d bytes\n- Lines: %d\n- Words: %d\n- Characters: %d",
        analysis["file_name"], analysis["file_size"], analysis["line_count"], analysis["word_count"], analysis["char_count"])
    
    response := map[string]interface{}{
        "content": []map[string]interface{}{
            {"type": "text", "text": resultText},
        },
    }
    
    return json.Marshal(response)
}
```

Build the plugin:

```bash
cd tools
go build -buildmode=plugin -o file_analyzer.so file_analyzer.go
```

## 🛡️ Security Considerations

- **Authentication**: MCP connections should be properly authenticated
- **Authorization**: Implement proper access controls for resources and tools
- **Input Validation**: All MCP requests are validated and sanitized
- **Execution Isolation**: Tools run in isolated contexts with timeouts
- **File System Access**: Restricted to designated resources and tools directories
- **Network Security**: Use HTTPS in production environments

## 🔧 Environment Variables

- `GIN_MCP_RESOURCES_DIR` - Resources directory path
- `GIN_MCP_TOOLS_DIR` - Tools directory path
- `GIN_MCP_PORT` - Server port (standalone mode)

## 📚 Examples

See the `examples/` directory for complete working examples:

- `examples/integration/` - Integration with existing Gin app
- `examples/standalone/` - Standalone MCP server

## 🤝 Contributing

Please see the main project [Contributing Guide](../../CONTRIBUTING.md) for details on how to contribute to this package.

---

*"Add Model Context Protocol server capabilities to your Gin applications with elegance and precision."* ✨ 