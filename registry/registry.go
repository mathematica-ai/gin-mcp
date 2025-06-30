package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
)

// ResourceType represents the type of MCP resource
type ResourceType string

const (
	FileResource     ResourceType = "file"
	DatabaseResource ResourceType = "database"
	APIRResource     ResourceType = "api"
	UnknownResource  ResourceType = "unknown"
)

// ToolType represents the type of MCP tool
type ToolType string

const (
	GoPluginTool ToolType = "go_plugin"
	PythonTool   ToolType = "python"
	UnknownTool  ToolType = "unknown"
)

// ResourceInfo contains metadata about a registered MCP resource
type ResourceInfo struct {
	Name     string       `json:"name"`
	FilePath string       `json:"file_path"`
	Type     ResourceType `json:"type"`
	MimeType string       `json:"mime_type"`
}

// ToolInfo contains metadata about a registered MCP tool
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	FilePath    string                 `json:"file_path"`
	Type        ToolType               `json:"type"`
	InputSchema map[string]interface{} `json:"input_schema"`
	Handler     interface{}            `json:"-"`
}

// Registry manages the collection of available MCP resources and tools
type Registry struct {
	resources map[string]*ResourceInfo
	tools     map[string]*ToolInfo
	mutex     sync.RWMutex
}

// NewRegistry creates a new MCP registry
func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]*ResourceInfo),
		tools:     make(map[string]*ToolInfo),
	}
}

// RegisterResource adds a resource to the registry
func (r *Registry) RegisterResource(name, filePath string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	resourceType := r.determineResourceType(filePath)
	mimeType := r.determineMimeType(filePath)

	resourceInfo := &ResourceInfo{
		Name:     name,
		FilePath: filePath,
		Type:     resourceType,
		MimeType: mimeType,
	}

	r.resources[name] = resourceInfo

	log.Printf("✅ Registered MCP resource: %s (%s) at %s", name, resourceType, filePath)
	return nil
}

// RegisterTool adds a tool to the registry
func (r *Registry) RegisterTool(name, filePath, description string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	toolType := r.determineToolType(filePath)

	toolInfo := &ToolInfo{
		Name:        name,
		Description: description,
		FilePath:    filePath,
		Type:        toolType,
		InputSchema: r.generateInputSchema(toolType),
	}

	// Load the appropriate handler based on tool type
	handler, err := r.loadToolHandler(toolInfo)
	if err != nil {
		return fmt.Errorf("failed to load handler for tool %s: %w", name, err)
	}

	toolInfo.Handler = handler
	r.tools[name] = toolInfo

	log.Printf("✅ Registered MCP tool: %s (%s) at %s", name, toolType, filePath)
	return nil
}

// UnregisterResource removes a resource from the registry
func (r *Registry) UnregisterResource(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.resources[name]; exists {
		delete(r.resources, name)
		log.Printf("🗑️  Unregistered MCP resource: %s", name)
	}
}

// UnregisterTool removes a tool from the registry
func (r *Registry) UnregisterTool(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tools[name]; exists {
		delete(r.tools, name)
		log.Printf("🗑️  Unregistered MCP tool: %s", name)
	}
}

// GetResource retrieves a resource from the registry
func (r *Registry) GetResource(name string) (*ResourceInfo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	resource, exists := r.resources[name]
	return resource, exists
}

// GetTool retrieves a tool from the registry
func (r *Registry) GetTool(name string) (*ToolInfo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}

// ListResources returns all registered resources
func (r *Registry) ListResources() []*ResourceInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	resources := make([]*ResourceInfo, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}
	return resources
}

// ListTools returns all registered tools
func (r *Registry) ListTools() []*ToolInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tools := make([]*ToolInfo, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// determineResourceType identifies the type of resource based on file extension
func (r *Registry) determineResourceType(filePath string) ResourceType {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".sql":
		return DatabaseResource
	case ".json", ".yaml", ".yml", ".xml":
		return FileResource
	case ".md", ".txt", ".csv":
		return FileResource
	default:
		return FileResource
	}
}

// determineToolType identifies the type of tool based on file extension
func (r *Registry) determineToolType(filePath string) ToolType {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".so":
		return GoPluginTool
	case ".py":
		return PythonTool
	default:
		return UnknownTool
	}
}

// determineMimeType determines the MIME type based on file extension
func (r *Registry) determineMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/x-yaml"
	case ".xml":
		return "application/xml"
	case ".md":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".sql":
		return "text/sql"
	default:
		return mime.TypeByExtension(ext)
	}
}

// generateInputSchema generates a basic input schema for tools
func (r *Registry) generateInputSchema(toolType ToolType) map[string]interface{} {
	switch toolType {
	case GoPluginTool:
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"arguments": map[string]interface{}{
					"type":        "object",
					"description": "Tool arguments",
				},
			},
		}
	case PythonTool:
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"arguments": map[string]interface{}{
					"type":        "object",
					"description": "Tool arguments",
				},
			},
		}
	default:
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"arguments": map[string]interface{}{
					"type":        "object",
					"description": "Tool arguments",
				},
			},
		}
	}
}

// loadToolHandler creates the appropriate handler for the tool type
func (r *Registry) loadToolHandler(toolInfo *ToolInfo) (interface{}, error) {
	switch toolInfo.Type {
	case GoPluginTool:
		return r.loadGoPlugin(toolInfo.FilePath)
	case PythonTool:
		return r.loadPythonScript(toolInfo.FilePath)
	default:
		return nil, fmt.Errorf("unsupported tool type: %s", toolInfo.Type)
	}
}

// loadGoPlugin loads a Go plugin and returns the Execute function
func (r *Registry) loadGoPlugin(filePath string) (interface{}, error) {
	plug, err := plugin.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %w", filePath, err)
	}

	executeSymbol, err := plug.Lookup("Execute")
	if err != nil {
		return nil, fmt.Errorf("failed to find Execute function in plugin %s: %w", filePath, err)
	}

	return executeSymbol, nil
}

// loadPythonScript creates a handler for Python scripts
func (r *Registry) loadPythonScript(filePath string) (interface{}, error) {
	// For Python scripts, we return the file path as the handler
	// The actual execution will be handled by the handler package
	return filePath, nil
}

// GetResourceNames returns a list of all registered resource names
func (r *Registry) GetResourceNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.resources))
	for name := range r.resources {
		names = append(names, name)
	}
	return names
}

// GetToolNames returns a list of all registered tool names
func (r *Registry) GetToolNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetResourceCount returns the number of registered resources
func (r *Registry) GetResourceCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.resources)
}

// GetToolCount returns the number of registered tools
func (r *Registry) GetToolCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.tools)
}

// ExportRegistry exports the registry data for debugging
func (r *Registry) ExportRegistry() ([]byte, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	export := map[string]interface{}{
		"resources": r.resources,
		"tools":     r.tools,
		"counts": map[string]int{
			"resources": len(r.resources),
			"tools":     len(r.tools),
		},
	}

	return json.MarshalIndent(export, "", "  ")
}
