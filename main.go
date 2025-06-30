package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gin-mcp/pkg/ginmcp"
)

func main() {
	// Get configuration from environment variables
	resourcesDir := os.Getenv("GIN_MCP_RESOURCES_DIR")
	if resourcesDir == "" {
		resourcesDir = "./resources"
	}

	toolsDir := os.Getenv("GIN_MCP_TOOLS_DIR")
	if toolsDir == "" {
		toolsDir = "./tools"
	}

	port := os.Getenv("GIN_MCP_PORT")
	if port == "" {
		port = ":8080"
	}

	// Initialize MCP server with configuration
	config := &ginmcp.MCPConfig{
		ResourcesDir: resourcesDir,
		ToolsDir:     toolsDir,
		Prefix:       "/mcp",
		Port:         port,
	}

	mcp, err := ginmcp.New(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize MCP server: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("📡 Received shutdown signal")
		mcp.Stop()
		os.Exit(0)
	}()

	// Start MCP server as a standalone server
	log.Printf("🚀 Starting gin-mcp MCP server on port %s", port)
	log.Printf("📁 Watching resources directory: %s", resourcesDir)
	log.Printf("🔧 Watching tools directory: %s", toolsDir)
	log.Printf("🔌 MCP Protocol Version: 2025.06.18")

	if err := mcp.StartStandalone(); err != nil {
		log.Fatalf("❌ Failed to start MCP server: %v", err)
	}
}
