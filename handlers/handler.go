package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"gin-mcp/registry"
)

// MCPHandler handles MCP resource access and tool execution
type MCPHandler struct{}

// NewMCPHandler creates a new MCP handler
func NewMCPHandler() *MCPHandler {
	return &MCPHandler{}
}

// AccessResource accesses an MCP resource and returns its content
func (h *MCPHandler) AccessResource(resourceInfo *registry.ResourceInfo, input []byte) ([]byte, error) {
	// Validate input
	if err := h.ValidateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Read the resource file
	content, err := ioutil.ReadFile(resourceInfo.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource file: %w", err)
	}

	// Create MCP resource response
	response := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":       fmt.Sprintf("file://%s", resourceInfo.FilePath),
				"mime_type": resourceInfo.MimeType,
				"text":      string(content),
			},
		},
	}

	return json.Marshal(response)
}

// ExecuteTool executes an MCP tool with the given input and returns the result
func (h *MCPHandler) ExecuteTool(toolInfo *registry.ToolInfo, input []byte) ([]byte, error) {
	// Validate input
	if err := h.ValidateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	switch toolInfo.Type {
	case registry.GoPluginTool:
		return h.executeGoPlugin(toolInfo, input)
	case registry.PythonTool:
		return h.executePythonScript(toolInfo, input)
	default:
		return nil, fmt.Errorf("unsupported tool type: %s", toolInfo.Type)
	}
}

// executeGoPlugin executes a Go plugin tool
func (h *MCPHandler) executeGoPlugin(toolInfo *registry.ToolInfo, input []byte) ([]byte, error) {
	// Get the Execute function from the plugin
	executeFunc := toolInfo.Handler
	if executeFunc == nil {
		return nil, fmt.Errorf("no handler found for tool %s", toolInfo.Name)
	}

	// Type assert to the expected function signature
	execute, ok := executeFunc.(func([]byte) ([]byte, error))
	if !ok {
		return nil, fmt.Errorf("invalid Execute function signature for tool %s", toolInfo.Name)
	}

	// Execute the tool with timeout
	resultChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := execute(input)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		log.Printf("✅ Go plugin tool %s executed successfully", toolInfo.Name)
		return h.validateAndFormatOutput(result)
	case err := <-errChan:
		return nil, fmt.Errorf("go plugin execution failed: %w", err)
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("go plugin execution timed out after 30 seconds")
	}
}

// executePythonScript executes a Python script tool
func (h *MCPHandler) executePythonScript(toolInfo *registry.ToolInfo, input []byte) ([]byte, error) {
	// Create command to execute Python script
	cmd := exec.Command("python3", toolInfo.FilePath)

	// Set up input/output pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Python script: %w", err)
	}

	// Send input to the script
	if _, err := stdin.Write(input); err != nil {
		return nil, fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	// Wait for completion with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("Python script execution failed: %w, stderr: %s", err, stderr.String())
		}
	case <-time.After(30 * time.Second):
		cmd.Process.Kill()
		return nil, fmt.Errorf("Python script execution timed out after 30 seconds")
	}

	output := stdout.Bytes()
	log.Printf("✅ Python script tool %s executed successfully", toolInfo.Name)

	return h.validateAndFormatOutput(output)
}

// ValidateInput validates the input JSON
func (h *MCPHandler) ValidateInput(input []byte) error {
	if len(input) == 0 {
		return fmt.Errorf("input cannot be empty")
	}

	var jsonData interface{}
	if err := json.Unmarshal(input, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON input: %w", err)
	}

	return nil
}

// ValidateOutput validates the output JSON
func (h *MCPHandler) ValidateOutput(output []byte) error {
	if len(output) == 0 {
		return fmt.Errorf("output cannot be empty")
	}

	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON output: %w", err)
	}

	return nil
}

// validateAndFormatOutput validates and formats the tool output
func (h *MCPHandler) validateAndFormatOutput(output []byte) ([]byte, error) {
	// Try to validate as JSON first
	if err := h.ValidateOutput(output); err == nil {
		// If it's valid JSON, return as is
		return output, nil
	}

	// If it's not valid JSON, format it as MCP content
	response := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": strings.TrimSpace(string(output)),
			},
		},
	}

	return json.Marshal(response)
}
