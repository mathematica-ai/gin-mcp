package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gin-mcp/pkg/ginmcp"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Add your existing routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to my application with MCP server integration!",
			"version": "1.0.0",
			"mcp":     "Model Context Protocol Server",
		})
	})

	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"users": []string{"user1", "user2", "user3"},
		})
	})

	// Initialize MCP server with custom configuration
	mcpConfig := &ginmcp.MCPConfig{
		ResourcesDir: "./resources", // Directory for MCP resources
		ToolsDir:     "./tools",     // Directory for MCP tools
		Prefix:       "/mcp",        // URL prefix for MCP endpoints
	}

	mcp, err := ginmcp.New(mcpConfig)
	if err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
	}

	// Add MCP server routes to your existing router
	if err := mcp.SetupRoutes(router); err != nil {
		log.Fatalf("Failed to setup MCP server routes: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("📡 Received shutdown signal")

		// Gracefully stop MCP server
		if err := mcp.Stop(); err != nil {
			log.Printf("❌ Error stopping MCP server: %v", err)
		}

		log.Println("✅ Application shutdown complete")
		os.Exit(0)
	}()

	// Start the server
	log.Println("🚀 Starting application with MCP server integration...")
	log.Println("📁 MCP resources available at: /mcp/resources/*")
	log.Println("🔧 MCP tools available at: /mcp/tools/*")
	log.Println("🔍 Health check: GET /mcp/health")
	log.Println("📋 List resources: GET /mcp/resources")
	log.Println("📋 List tools: GET /mcp/tools")
	log.Println("⚡ Access resource: POST /mcp/resources/{name}")
	log.Println("⚡ Execute tool: POST /mcp/tools/{name}")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
