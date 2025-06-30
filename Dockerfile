# Multi-stage build for gin-mcp MCP server
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for plugin support
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gin-mcp .

# Build sample Go tools
RUN cd tools && go build -buildmode=plugin -o calculator.so calculator.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata python3 py3-pip

# Install Python packages for sample tools
RUN pip3 install pandas numpy

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gin-mcp .

# Copy built tools
COPY --from=builder /app/tools/*.so ./tools/
COPY --from=builder /app/tools/*.py ./tools/

# Create MCP directories and copy sample resources
RUN mkdir -p /app/resources /app/tools
COPY resources/ ./resources/
RUN chmod +x ./tools/*.py

# Set proper ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check using the MCP health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/mcp/health || exit 1

# Set environment variables
ENV GIN_MCP_RESOURCES_DIR=/app/resources
ENV GIN_MCP_TOOLS_DIR=/app/tools
ENV GIN_MCP_PORT=:8080

# Run the MCP server
CMD ["./gin-mcp"] 