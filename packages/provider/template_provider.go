package provider

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dropdx/dropdx/packages/config"
)

/**
 * TemplateProvider implements the Provider interface for template-based files.
 */
type TemplateProvider struct {
	name     string
	template string
	target   string
}

/**
 * NewTemplateProvider creates a new instance of TemplateProvider.
 */
func NewTemplateProvider(name, tmpl, target string) *TemplateProvider {
	return &TemplateProvider{
		name:     name,
		template: tmpl,
		target:   target,
	}
}

/**
 * Name returns the name of the provider.
 */
func (tp *TemplateProvider) Name() string {
	return tp.name
}

/**
 * Apply processes the template and writes the target file.
 */
func (tp *TemplateProvider) Apply(tokens map[string]string) error {
	// 1. Resolve paths
	base := config.GetBaseDir()

	tmplPath := tp.template
	if !filepath.IsAbs(tmplPath) && !strings.HasPrefix(tmplPath, "~") {
		tmplPath = filepath.Join(base, tmplPath)
	}

	resolvedTmpl, err := config.ResolvePath(tmplPath)
	if err != nil {
		return err
	}

	resolvedTarget, err := config.ResolvePath(tp.target)
	if err != nil {
		return err
	}

	// 2. Read template
	content, err := os.ReadFile(resolvedTmpl)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", resolvedTmpl, err)
	}

	// 3. Process template
	tmpl, err := template.New(tp.name).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tokens); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 4. Write to target
	if err := os.MkdirAll(filepath.Dir(resolvedTarget), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	if err := os.WriteFile(resolvedTarget, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write target file %s: %w", resolvedTarget, err)
	}

	fmt.Printf("✔ Applied provider: %s -> %s\n", tp.name, resolvedTarget)
	return nil
}
