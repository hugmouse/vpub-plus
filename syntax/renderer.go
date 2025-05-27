package syntax

import (
	"fmt"
	"maps"
	"slices"
)

// Renderer is the minimal interface every rendering engine must satisfy
type Renderer interface {
	// Convert turns source text into HTML (or any other one)
	//
	// For example, markdown -> HTML
	//
	// Note: wrap bool is used by vanilla renderer,
	// to wrap in <p> tags the text nodes.
	// You can ignore it if you so desire.
	Convert(input string, wrap bool) string
	Name() string
}

type RenderEngineRegistry struct {
	registry map[string]Renderer
}

type RendererNotFoundError struct {
	Name      string
	Available []string
}

func (e *RendererNotFoundError) Error() string {
	return fmt.Sprintf("syntax: renderer '%s' not found (available: %q)", e.Name, e.Available)
}

type RendererAlreadyExistsError struct {
	Name string
}

func (e *RendererAlreadyExistsError) Error() string {
	return fmt.Sprintf("syntax: renderer %q already exists", e.Name)
}

const (
	// ImageProxyURLBase is used for renderers that wish to support
	// image proxy service.
	//
	// To learn more, you can see convert_blackfriday.go source code
	ImageProxyURLBase = "/image-proxy?url="
)

// NewRegistry creates a new render engine registry
// where you can Register, Get and List your engines
func NewRegistry() *RenderEngineRegistry {
	return &RenderEngineRegistry{
		registry: make(map[string]Renderer),
	}
}

// Register makes a renderer available under the given name
//
// To get a renderer, use Get
func (r *RenderEngineRegistry) Register(name string, renderer Renderer) error {
	if _, exists := r.registry[name]; exists {
		return &RendererAlreadyExistsError{name}
	}
	r.registry[name] = renderer
	return nil
}

// Get returns the renderer registered under name, or an error
//
// To register a renderer, use Register
func (r *RenderEngineRegistry) Get(name string) (Renderer, error) {
	renderer, exists := r.registry[name]
	if !exists {
		return nil, &RendererNotFoundError{name, slices.Collect(maps.Keys(r.registry))}
	}
	return renderer, nil
}

// List returns all registered renderers
func (r *RenderEngineRegistry) List() []Renderer {
	return slices.Collect(maps.Values(r.registry))
}
