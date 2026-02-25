package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dcdavidev/dropdx/internal/config"
)

/**
 * Engine is the core component that processes templates and writes target files.
 */
type Engine struct {
	Config *config.Config
}

/**
 * NewEngine creates a new instance of the Engine.
 */
func NewEngine(cfg *config.Config) *Engine {
	return &Engine{Config: cfg}
}

/**
 * ApplyProvider processes a specific provider's template and writes to the target path.
 */
func (e *Engine) ApplyProvider(name string) error {
	p, ok := e.Config.Providers[name]
	if !ok {
		return fmt.Errorf("provider %s not found in configuration", name)
	}

	// 1. Resolve paths
	base := config.GetBaseDir()
	
	// If path is relative, resolve it against the config base dir
	tmplPath := p.Template
	if !filepath.IsAbs(tmplPath) && !strings.HasPrefix(tmplPath, "~") {
		tmplPath = filepath.Join(base, tmplPath)
	}
	
	resolvedTmpl, err := config.ResolvePath(tmplPath)
	if err != nil {
		return err
	}

	resolvedTarget, err := config.ResolvePath(p.Target)
	if err != nil {
		return err
	}

	// 2. Read template
	content, err := os.ReadFile(resolvedTmpl)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", resolvedTmpl, err)
	}

	// 3. Process template
	// We use the provider name as the template name.
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	// Execute template passing the Tokens map
	if err := tmpl.Execute(&buf, e.Config.Tokens); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 4. Write to target
	// Ensure the target directory exists (e.g., if writing to ~/.config/myapp/config)
	if err := os.MkdirAll(filepath.Dir(resolvedTarget), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	if err := os.WriteFile(resolvedTarget, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write target file %s: %w", resolvedTarget, err)
	}

	fmt.Printf("✔ Applied provider: %s -> %s\n", name, resolvedTarget)
	return nil
}

/**
 * ApplyAll processes all providers defined in the configuration.
 */
func (e *Engine) ApplyAll() error {
	if len(e.Config.Providers) == 0 {
		fmt.Println("ℹ No providers configured. Add some in your config.yaml.")
		return nil
	}

	for name := range e.Config.Providers {
		if err := e.ApplyProvider(name); err != nil {
			return fmt.Errorf("failed to apply provider %s: %w", name, err)
		}
	}
	return nil
}
