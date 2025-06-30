package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gin-mcp/registry"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors file system changes for MCP resources and tools
type Watcher struct {
	watcher      *fsnotify.Watcher
	registry     *registry.Registry
	resourcesDir string
	toolsDir     string
	isRunning    bool
	stopChan     chan bool
}

// NewWatcher creates a new file watcher for MCP resources and tools
func NewWatcher(resourcesDir, toolsDir string, registry *registry.Registry) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	return &Watcher{
		watcher:      watcher,
		registry:     registry,
		resourcesDir: resourcesDir,
		toolsDir:     toolsDir,
		stopChan:     make(chan bool),
	}, nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	// Create directories if they don't exist
	if err := w.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Add directories to watcher
	if err := w.watcher.Add(w.resourcesDir); err != nil {
		return fmt.Errorf("failed to add resources directory to watcher: %w", err)
	}

	if err := w.watcher.Add(w.toolsDir); err != nil {
		return fmt.Errorf("failed to add tools directory to watcher: %w", err)
	}

	// Initial scan of existing files
	if err := w.scanExistingFiles(); err != nil {
		return fmt.Errorf("failed to scan existing files: %w", err)
	}

	w.isRunning = true

	// Start watching for changes
	go w.watchLoop()

	log.Printf("👀 File watcher started for resources: %s, tools: %s", w.resourcesDir, w.toolsDir)
	return nil
}

// Stop stops the file watcher
func (w *Watcher) Stop() error {
	if !w.isRunning {
		return nil
	}

	w.isRunning = false
	w.stopChan <- true

	if err := w.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close watcher: %w", err)
	}

	log.Printf("🛑 File watcher stopped")
	return nil
}

// IsRunning returns whether the watcher is currently running
func (w *Watcher) IsRunning() bool {
	return w.isRunning
}

// ensureDirectories creates the resources and tools directories if they don't exist
func (w *Watcher) ensureDirectories() error {
	dirs := []string{w.resourcesDir, w.toolsDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// scanExistingFiles scans for existing files and registers them
func (w *Watcher) scanExistingFiles() error {
	// Scan resources directory
	if err := w.scanDirectory(w.resourcesDir, "resource"); err != nil {
		return fmt.Errorf("failed to scan resources directory: %w", err)
	}

	// Scan tools directory
	if err := w.scanDirectory(w.toolsDir, "tool"); err != nil {
		return fmt.Errorf("failed to scan tools directory: %w", err)
	}

	return nil
}

// scanDirectory scans a directory for files and registers them
func (w *Watcher) scanDirectory(dir, itemType string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		if itemType == "resource" {
			if err := w.registry.RegisterResource(name, filePath); err != nil {
				log.Printf("⚠️  Failed to register resource %s: %v", name, err)
			}
		} else if itemType == "tool" {
			description := fmt.Sprintf("MCP tool: %s", name)
			if err := w.registry.RegisterTool(name, filePath, description); err != nil {
				log.Printf("⚠️  Failed to register tool %s: %v", name, err)
			}
		}
	}

	return nil
}

// watchLoop monitors for file system events
func (w *Watcher) watchLoop() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("❌ Watcher error: %v", err)

		case <-w.stopChan:
			return
		}
	}
}

// handleEvent processes a file system event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Ignore temporary files and hidden files
	if strings.HasPrefix(filepath.Base(event.Name), ".") {
		return
	}

	// Determine if this is a resource or tool based on the directory
	isResource := strings.HasPrefix(event.Name, w.resourcesDir)
	isTool := strings.HasPrefix(event.Name, w.toolsDir)

	if !isResource && !isTool {
		return
	}

	name := strings.TrimSuffix(filepath.Base(event.Name), filepath.Ext(event.Name))

	switch event.Op {
	case fsnotify.Create, fsnotify.Write:
		if isResource {
			if err := w.registry.RegisterResource(name, event.Name); err != nil {
				log.Printf("⚠️  Failed to register resource %s: %v", name, err)
			} else {
				log.Printf("✅ Resource %s registered/updated", name)
			}
		} else if isTool {
			description := fmt.Sprintf("MCP tool: %s", name)
			if err := w.registry.RegisterTool(name, event.Name, description); err != nil {
				log.Printf("⚠️  Failed to register tool %s: %v", name, err)
			} else {
				log.Printf("✅ Tool %s registered/updated", name)
			}
		}

	case fsnotify.Remove:
		if isResource {
			w.registry.UnregisterResource(name)
			log.Printf("🗑️  Resource %s unregistered", name)
		} else if isTool {
			w.registry.UnregisterTool(name)
			log.Printf("🗑️  Tool %s unregistered", name)
		}

	case fsnotify.Rename:
		// Handle rename as remove + create
		if isResource {
			w.registry.UnregisterResource(name)
			log.Printf("🔄 Resource %s renamed", name)
		} else if isTool {
			w.registry.UnregisterTool(name)
			log.Printf("🔄 Tool %s renamed", name)
		}
	}
}

// GetStats returns statistics about the watcher
func (w *Watcher) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"running":       w.isRunning,
		"resources_dir": w.resourcesDir,
		"tools_dir":     w.toolsDir,
		"resources":     w.registry.GetResourceCount(),
		"tools":         w.registry.GetToolCount(),
	}
}
