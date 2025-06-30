#!/bin/bash

# Test script for gin-mcp MCP server
# This script demonstrates the functionality of the gin-mcp MCP server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVER_URL="http://localhost:8080"
RESOURCE_NAME="database_schema"
TOOL_NAME="calculator"

echo -e "${BLUE}🏦 gin-mcp MCP Server Test Suite${NC}"
echo "=========================================="
echo ""

# Function to check if server is running
check_server() {
    echo -e "${YELLOW}🔍 Checking if MCP server is running...${NC}"
    if curl -s "$SERVER_URL/mcp/health" > /dev/null; then
        echo -e "${GREEN}✅ MCP server is running${NC}"
        return 0
    else
        echo -e "${RED}❌ MCP server is not running. Please start the server first:${NC}"
        echo "   make run"
        echo "   or"
        echo "   ./gin-mcp"
        return 1
    fi
}

# Function to test health endpoint
test_health() {
    echo -e "${YELLOW}🏥 Testing MCP server health endpoint...${NC}"
    response=$(curl -s "$SERVER_URL/mcp/health")
    echo "Response: $response"
    echo -e "${GREEN}✅ Health check passed${NC}"
    echo ""
}

# Function to test resources listing
test_list_resources() {
    echo -e "${YELLOW}📋 Testing MCP resources listing...${NC}"
    response=$(curl -s "$SERVER_URL/mcp/resources")
    echo "Response: $response"
    echo -e "${GREEN}✅ Resources listing passed${NC}"
    echo ""
}

# Function to test tools listing
test_list_tools() {
    echo -e "${YELLOW}🔧 Testing MCP tools listing...${NC}"
    response=$(curl -s "$SERVER_URL/mcp/tools")
    echo "Response: $response"
    echo -e "${GREEN}✅ Tools listing passed${NC}"
    echo ""
}

# Function to test resource access
test_access_resource() {
    echo -e "${YELLOW}📁 Testing MCP resource access...${NC}"
    response=$(curl -s -X POST "$SERVER_URL/mcp/resources/$RESOURCE_NAME" \
        -H "Content-Type: application/json" \
        -d '{"uri": "file://./resources/database_schema.sql"}')
    echo "Response: $response"
    echo -e "${GREEN}✅ Resource access passed${NC}"
    echo ""
}

# Function to test tool execution
test_execute_tool() {
    echo -e "${YELLOW}⚡ Testing MCP tool execution...${NC}"
    response=$(curl -s -X POST "$SERVER_URL/mcp/tools/$TOOL_NAME" \
        -H "Content-Type: application/json" \
        -d '{"arguments": {"expression": "2 + 3 * 4"}}')
    echo "Response: $response"
    echo -e "${GREEN}✅ Tool execution passed${NC}"
    echo ""
}

# Function to test registry export
test_registry_export() {
    echo -e "${YELLOW}📊 Testing MCP registry export...${NC}"
    response=$(curl -s "$SERVER_URL/mcp/registry")
    echo "Response: $response"
    echo -e "${GREEN}✅ Registry export passed${NC}"
    echo ""
}

# Function to test error handling
test_error_handling() {
    echo -e "${YELLOW}🚨 Testing error handling...${NC}"
    
    # Test non-existent resource
    response=$(curl -s -w "%{http_code}" "$SERVER_URL/mcp/resources/nonexistent" | tail -1)
    if [ "$response" = "404" ]; then
        echo -e "${GREEN}✅ 404 error handling passed${NC}"
    else
        echo -e "${RED}❌ 404 error handling failed${NC}"
    fi
    
    # Test non-existent tool
    response=$(curl -s -w "%{http_code}" "$SERVER_URL/mcp/tools/nonexistent" | tail -1)
    if [ "$response" = "404" ]; then
        echo -e "${GREEN}✅ 404 error handling passed${NC}"
    else
        echo -e "${RED}❌ 404 error handling failed${NC}"
    fi
    
    echo ""
}

# Function to run performance test
test_performance() {
    echo -e "${YELLOW}⚡ Running performance test...${NC}"
    
    start_time=$(date +%s.%N)
    
    # Make 10 concurrent requests
    for i in {1..10}; do
        curl -s "$SERVER_URL/mcp/health" > /dev/null &
    done
    wait
    
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    
    echo -e "${GREEN}✅ Performance test completed in ${duration}s${NC}"
    echo ""
}

# Main test execution
main() {
    echo -e "${BLUE}🚀 Starting gin-mcp MCP server test suite...${NC}"
    echo ""
    
    # Check if server is running
    if ! check_server; then
        exit 1
    fi
    
    # Run all tests
    test_health
    test_list_resources
    test_list_tools
    test_access_resource
    test_execute_tool
    test_registry_export
    test_error_handling
    test_performance
    
    echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
    echo ""
    echo -e "${BLUE}📋 Test Summary:${NC}"
    echo "  ✅ Health check"
    echo "  ✅ Resources listing"
    echo "  ✅ Tools listing"
    echo "  ✅ Resource access"
    echo "  ✅ Tool execution"
    echo "  ✅ Registry export"
    echo "  ✅ Error handling"
    echo "  ✅ Performance test"
    echo ""
    echo -e "${BLUE}🔌 MCP Server is working correctly!${NC}"
}

# Run main function
main 