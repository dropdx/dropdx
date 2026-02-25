package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dropdx/dropdx/internal/config"
)

/**
 * TestApplyProvider ensures the engine correctly processes a template and
 * writes the resulting file with injected tokens.
 */
func TestApplyProvider(t *testing.T) {
	// 1. Setup temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dropdx-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmplPath := filepath.Join(tmpDir, "test.tmpl")
	targetPath := filepath.Join(tmpDir, "output", "result.txt")
	
	tmplContent := "Hello {{.name}}, your token is {{.token}}!"
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	// 2. Setup configuration with metadata structure
	cfg := &config.Config{
		Tokens: map[string]config.TokenInfo{
			"name":  {Value: "User"},
			"token": {Value: "secret123"},
		},
		Providers: map[string]config.Provider{
			"test": {
				Template: tmplPath,
				Target:   targetPath,
			},
		},
	}

	// 3. Run Engine
	engine := NewEngine(cfg)
	if err := engine.ApplyProvider("test"); err != nil {
		t.Fatalf("ApplyProvider() failed: %v", err)
	}

	// 4. Verify results
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	want := "Hello User, your token is secret123!"
	if string(got) != want {
		t.Errorf("ApplyProvider() content got = %q, want %q", string(got), want)
	}
}
