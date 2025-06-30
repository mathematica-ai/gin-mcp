package ginmcp

import (
	"encoding/json"
	"fmt"
	"log"

	"gin-mcp/handlers"
	"gin-mcp/registry"
	"gin-mcp/watcher"

	"github.com/gin-gonic/gin"
)

// MCPConfig holds configuration for the MCP server
type MCPConfig struct {
	ResourcesDir string // Directory to watch for MCP resources
	ToolsDir     string // Directory to watch for MCP tools
	Prefix       string // URL prefix for MCP endpoints (default: "/mcp")
	Port         string // Port for the MCP server (if standalone)
}

// DefaultConfig returns default configuration
func DefaultConfig() *MCPConfig {
	return &MCPConfig{
		ResourcesDir: "./resources",
		ToolsDir:     "./tools",
		Prefix:       "/mcp",
		Port:         ":8080",
	}
}

// MCP represents the MCP server
type MCP struct {
	config   *MCPConfig
	registry *registry.Registry
	handler  *handlers.MCPHandler
	watcher  *watcher.Watcher
}

// New creates a new MCP server instance
func New(config *MCPConfig) (*MCP, error) {
	if config == nil {
		config = DefaultConfig()
	}

	return &MCP{
		config:   config,
		registry: registry.NewRegistry(),
		handler:  handlers.NewMCPHandler(),
	}, nil
}

// SetupRoutes adds MCP server routes to an existing Gin router
func (m *MCP) SetupRoutes(router *gin.Engine) error {
	// Initialize the file watcher
	if err := m.initializeWatcher(); err != nil {
		return fmt.Errorf("failed to initialize watcher: %w", err)
	}

	// Create MCP route group
	mcpGroup := router.Group(m.config.Prefix)

	// Health check endpoint
	mcpGroup.GET("/health", m.healthHandler)

	// MCP Resources endpoints
	mcpGroup.GET("/resources", m.listResourcesHandler)
	mcpGroup.GET("/resources/:name", m.getResourceInfoHandler)
	mcpGroup.POST("/resources/:name", m.accessResourceHandler)

	// MCP Tools endpoints
	mcpGroup.GET("/tools", m.listToolsHandler)
	mcpGroup.GET("/tools/:name", m.getToolInfoHandler)
	mcpGroup.POST("/tools/:name", m.executeToolHandler)

	// Registry export endpoint (for debugging)
	mcpGroup.GET("/registry", m.exportRegistryHandler)

	log.Printf("✅ MCP server routes registered at prefix: %s", m.config.Prefix)
	return nil
}

// StartStandalone starts MCP server as a standalone server
func (m *MCP) StartStandalone() error {
	router := gin.New()

	// Setup middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Setup MCP server routes
	if err := m.SetupRoutes(router); err != nil {
		return err
	}

	log.Printf("🚀 Starting MCP server on port %s", m.config.Port)
	return router.Run(m.config.Port)
}

// Stop gracefully shuts down the MCP server
func (m *MCP) Stop() error {
	if m.watcher != nil {
		return m.watcher.Stop()
	}
	return nil
}

// GetRegistry returns the registry for direct access
func (m *MCP) GetRegistry() *registry.Registry {
	return m.registry
}

// GetHandler returns the MCP handler for direct access
func (m *MCP) GetHandler() *handlers.MCPHandler {
	return m.handler
}

// initializeWatcher sets up the file watcher for both resources and tools
func (m *MCP) initializeWatcher() error {
	watcher, err := watcher.NewWatcher(m.config.ResourcesDir, m.config.ToolsDir, m.registry)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	m.watcher = watcher

	if err := m.watcher.Start(); err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}

	return nil
}

// healthHandler handles health check requests
func (m *MCP) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":      "healthy",
		"service":     "gin-mcp",
		"mcp_version": "2025.06.18",
		"resources":   m.registry.GetResourceCount(),
		"tools":       m.registry.GetToolCount(),
		"watcher":     m.watcher.IsRunning(),
		"prefix":      m.config.Prefix,
	})
}

// listResourcesHandler returns a list of all available MCP resources
func (m *MCP) listResourcesHandler(c *gin.Context) {
	resources := m.registry.ListResources()

	resourceList := make([]gin.H, len(resources))
	for i, resource := range resources {
		resourceList[i] = gin.H{
			"name":      resource.Name,
			"type":      resource.Type,
			"file_path": resource.FilePath,
			"mime_type": resource.MimeType,
		}
	}

	c.JSON(200, gin.H{
		"resources": resourceList,
		"count":     len(resourceList),
	})
}

// getResourceInfoHandler returns information about a specific MCP resource
func (m *MCP) getResourceInfoHandler(c *gin.Context) {
	resourceName := c.Param("name")

	resource, exists := m.registry.GetResource(resourceName)
	if !exists {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Resource '%s' not found", resourceName),
		})
		return
	}

	c.JSON(200, gin.H{
		"name":      resource.Name,
		"type":      resource.Type,
		"file_path": resource.FilePath,
		"mime_type": resource.MimeType,
	})
}

// accessResourceHandler accesses MCP resource content
func (m *MCP) accessResourceHandler(c *gin.Context) {
	resourceName := c.Param("name")

	// Get the resource from registry
	resource, exists := m.registry.GetResource(resourceName)
	if !exists {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Resource '%s' not found", resourceName),
		})
		return
	}

	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Failed to read request body: %v", err),
		})
		return
	}

	// Access the resource
	result, err := m.handler.AccessResource(resource, body)
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Resource access failed: %v", err),
		})
		return
	}

	// Parse the result as JSON and return it
	var jsonResult interface{}
	if err := json.Unmarshal(result, &jsonResult); err != nil {
		// If the result is not valid JSON, return it as a string
		c.JSON(200, gin.H{
			"contents": []gin.H{
				{
					"uri":       fmt.Sprintf("file://%s", resource.FilePath),
					"mime_type": resource.MimeType,
					"text":      string(result),
				},
			},
		})
		return
	}

	c.JSON(200, jsonResult)
}

// listToolsHandler returns a list of all available MCP tools
func (m *MCP) listToolsHandler(c *gin.Context) {
	tools := m.registry.ListTools()

	toolList := make([]gin.H, len(tools))
	for i, tool := range tools {
		toolList[i] = gin.H{
			"name":         tool.Name,
			"description":  tool.Description,
			"type":         tool.Type,
			"file_path":    tool.FilePath,
			"input_schema": tool.InputSchema,
		}
	}

	c.JSON(200, gin.H{
		"tools": toolList,
		"count": len(toolList),
	})
}

// getToolInfoHandler returns information about a specific MCP tool
func (m *MCP) getToolInfoHandler(c *gin.Context) {
	toolName := c.Param("name")

	tool, exists := m.registry.GetTool(toolName)
	if !exists {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Tool '%s' not found", toolName),
		})
		return
	}

	c.JSON(200, gin.H{
		"name":         tool.Name,
		"description":  tool.Description,
		"type":         tool.Type,
		"file_path":    tool.FilePath,
		"input_schema": tool.InputSchema,
	})
}

// executeToolHandler executes an MCP tool with the provided input
func (m *MCP) executeToolHandler(c *gin.Context) {
	toolName := c.Param("name")

	// Get the tool from registry
	tool, exists := m.registry.GetTool(toolName)
	if !exists {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Tool '%s' not found", toolName),
		})
		return
	}

	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Failed to read request body: %v", err),
		})
		return
	}

	// Execute the tool
	result, err := m.handler.ExecuteTool(tool, body)
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Tool execution failed: %v", err),
		})
		return
	}

	// Parse the result as JSON and return it
	var jsonResult interface{}
	if err := json.Unmarshal(result, &jsonResult); err != nil {
		// If the result is not valid JSON, return it as a string
		c.JSON(200, gin.H{
			"content": []gin.H{
				{
					"type": "text",
					"text": string(result),
				},
			},
		})
		return
	}

	c.JSON(200, jsonResult)
}

// exportRegistryHandler exports the registry for debugging
func (m *MCP) exportRegistryHandler(c *gin.Context) {
	registryData, err := m.registry.ExportRegistry()
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to export registry: %v", err),
		})
		return
	}

	var registryMap map[string]interface{}
	if err := json.Unmarshal(registryData, &registryMap); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to parse registry data",
		})
		return
	}

	c.JSON(200, registryMap)
}
