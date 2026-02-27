package provider

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/fatih/color"
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

	fileName := filepath.Base(resolvedTarget)
	isShellConfig := fileName == ".bashrc" || fileName == ".zshrc" || fileName == ".profile" || fileName == ".bash_profile"

	if isShellConfig {
		// Append or update the block in shell config
		existingContent, err := os.ReadFile(resolvedTarget)
		var newFullContent []byte
		markerStart := fmt.Sprintf("# >>> dropdx: %s start >>>", tp.name)
		markerEnd := fmt.Sprintf("# <<< dropdx: %s end <<<", tp.name)
		newBlock := fmt.Sprintf("\n%s\n%s\n%s\n", markerStart, strings.TrimSpace(buf.String()), markerEnd)

		action := "Applied"
		if err == nil {
			// File exists, try to replace existing block
			contentStr := string(existingContent)
			startIdx := strings.Index(contentStr, markerStart)
			endIdx := strings.Index(contentStr, markerEnd)

			if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
				// Replace existing block
				action = "Updated"
				newFullContent = []byte(contentStr[:startIdx] + strings.TrimPrefix(newBlock, "\n") + contentStr[endIdx+len(markerEnd):])
			} else {
				// Append to end
				action = "Appended to"
				newFullContent = append(existingContent, []byte(newBlock)...)
			}
		} else {
			// File doesn't exist, create new
			newFullContent = []byte(newBlock)
		}

		if err := os.WriteFile(resolvedTarget, newFullContent, 0644); err != nil {
			return fmt.Errorf("failed to write shell config %s: %w", resolvedTarget, err)
		}

		fmt.Printf("%s %s %s -> %s\n", 
			color.GreenString("✔"), 
			color.YellowString(action),
			color.MagentaString(tp.name), 
			color.CyanString(resolvedTarget))
	} else {
		// Overwrite for other files (default behavior)
		if err := os.WriteFile(resolvedTarget, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write target file %s: %w", resolvedTarget, err)
		}
		fmt.Printf("%s Applied %s -> %s\n", 
			color.GreenString("✔"), 
			color.MagentaString(tp.name), 
			color.CyanString(resolvedTarget))
	}

	return nil
}
