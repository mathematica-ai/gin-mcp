# Makefile for gin-mcp
# A beautiful, production-ready MCP server implementation

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-dev docker-compose-down docker-compose-logs dev build-tools

# Default target
help: ## Show this help message
	@echo "🏦 gin-mcp - Model Context Protocol Server"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the gin-mcp binary
	@echo "🔨 Building gin-mcp MCP server..."
	go build -ldflags="-s -w" -o gin-mcp .
	@echo "✅ Build complete: ./gin-mcp"

# Run the application
run: build ## Build and run the application
	@echo "🚀 Starting gin-mcp MCP server..."
	./gin-mcp

# Development mode
dev: ## Run in development mode with debug logging
	@echo "🔧 Starting gin-mcp MCP server in development mode..."
	GIN_MODE=debug go run main.go

# Test the application
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v ./...

# Test with coverage
test-coverage: ## Run tests with coverage report
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📈 Coverage report generated: coverage.html"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -f gin-mcp
	rm -f coverage.out coverage.html
	rm -f tools/*.so
	@echo "✅ Clean complete"

# Build MCP tools
build-tools: ## Build the sample MCP tools
	@echo "🔌 Building sample MCP tools..."
	@echo "🔨 Building calculator tool..."
	go build -buildmode=plugin -o tools/calculator.so tools/calculator.go
	@echo "🐍 Making data analyzer executable..."
	chmod +x tools/data_analyzer.py
	@echo "✅ MCP tools built successfully"

# Docker build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t gin-mcp .
	@echo "✅ Docker image built: gin-mcp"

# Docker run
docker-run: ## Run with Docker
	@echo "🐳 Running gin-mcp MCP server with Docker..."
	docker run -p 8080:8080 \
		-v $(PWD)/resources:/app/resources \
		-v $(PWD)/tools:/app/tools \
		-e GIN_MCP_RESOURCES_DIR=/app/resources \
		-e GIN_MCP_TOOLS_DIR=/app/tools \
		gin-mcp

# Docker run with custom port
docker-run-custom: ## Run with Docker on custom port
	@echo "🐳 Running gin-mcp MCP server with Docker on custom port..."
	docker run -p 3000:8080 \
		-v $(PWD)/resources:/app/resources \
		-v $(PWD)/tools:/app/tools \
		-e GIN_MCP_RESOURCES_DIR=/app/resources \
		-e GIN_MCP_TOOLS_DIR=/app/tools \
		gin-mcp

# Docker Compose commands
docker-compose-up: ## Run with Docker Compose
	@echo "🐳 Starting gin-mcp MCP server with Docker Compose..."
	docker-compose up --build

docker-compose-dev: ## Run development version with Docker Compose
	@echo "🔧 Starting gin-mcp MCP server in development mode with Docker Compose..."
	docker-compose --profile dev up --build gin-mcp-dev

docker-compose-down: ## Stop Docker Compose services
	@echo "🛑 Stopping Docker Compose services..."
	docker-compose down

docker-compose-logs: ## View Docker Compose logs
	@echo "📋 Viewing Docker Compose logs..."
	docker-compose logs -f

# Install dependencies
deps: ## Install Go dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

# Format code
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint: ## Lint Go code
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Code linted"

# Security scan
security: ## Run security scan
	@echo "🔒 Running security scan..."
	gosec ./...
	@echo "✅ Security scan complete"

# Generate documentation
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	godoc -http=:6060 &
	@echo "📖 Documentation available at http://localhost:6060"

# Create sample data
sample-data: ## Create sample data files
	@echo "📁 Creating sample data..."
	@mkdir -p resources tools
	@echo "✅ Sample directories created"

# Full setup
setup: deps sample-data build-tools ## Full setup including dependencies and sample data
	@echo "🚀 Setup complete! Ready to run gin-mcp MCP server"

# Development setup
dev-setup: deps sample-data build-tools ## Development setup
	@echo "🔧 Development setup complete!"
	@echo "Run 'make dev' to start the MCP server in development mode"

# Production build
prod-build: ## Production build with optimizations
	@echo "🏭 Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o gin-mcp .
	@echo "✅ Production build complete"

# Benchmark tests
bench: ## Run benchmark tests
	@echo "⚡ Running benchmarks..."
	go test -bench=. ./...
	@echo "✅ Benchmarks complete"

# Check for updates
check-updates: ## Check for dependency updates
	@echo "🔄 Checking for updates..."
	go list -u -m all
	@echo "✅ Update check complete"

# Update dependencies
update-deps: ## Update dependencies
	@echo "⬆️  Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated"

# Show version info
version: ## Show version information
	@echo "📋 Version Information:"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build time: $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')"

# Help for MCP development
mcp-help: ## Show MCP development help
	@echo "🔌 MCP Development Guide:"
	@echo ""
	@echo "📁 Resources:"
	@echo "  - Place files in ./resources/ to make them available as MCP resources"
	@echo "  - Supported formats: .sql, .md, .json, .txt, .csv, .yaml, .xml"
	@echo ""
	@echo "🔧 Tools:"
	@echo "  - Go plugins: Create .go files with Execute() function, build with 'make build-tools'"
	@echo "  - Python scripts: Create .py files, make executable with chmod +x"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  1. make setup          # Full setup"
	@echo "  2. make run            # Start MCP server"
	@echo "  3. curl localhost:8080/mcp/health  # Test health"
	@echo "  4. curl localhost:8080/mcp/resources  # List resources"
	@echo "  5. curl localhost:8080/mcp/tools  # List tools"

# Default target
.DEFAULT_GOAL := help 