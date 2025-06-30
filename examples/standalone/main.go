package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gin-mcp/pkg/ginmcp"
)

func main() {
	// Initialize MCP server with default configuration
	mcp, err := ginmcp.New(nil)
	if err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
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

		log.Println("✅ MCP server shutdown complete")
		os.Exit(0)
	}()

	// Start MCP server as a standalone server
	log.Println("🚀 Starting MCP server...")
	log.Println("📁 Resources directory: ./resources")
	log.Println("🔧 Tools directory: ./tools")
	log.Println("🔍 Health check: GET /mcp/health")
	log.Println("📋 List resources: GET /mcp/resources")
	log.Println("📋 List tools: GET /mcp/tools")
	log.Println("⚡ Access resource: POST /mcp/resources/{name}")
	log.Println("⚡ Execute tool: POST /mcp/tools/{name}")
	log.Println("🔌 MCP Protocol Version: 2025.06.18")

	if err := mcp.StartStandalone(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}
