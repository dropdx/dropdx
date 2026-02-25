package core

import (
	"fmt"

	"github.com/dcdavidev/dropdx/internal/config"
	"github.com/dcdavidev/dropdx/internal/provider"
)

/**
 * Engine is the core component that manages providers and tokens.
 */
type Engine struct {
	Config    *config.Config
	providers map[string]provider.Provider
}

/**
 * NewEngine creates a new instance of the Engine and registers providers.
 */
func NewEngine(cfg *config.Config) *Engine {
	e := &Engine{
		Config:    cfg,
		providers: make(map[string]provider.Provider),
	}

	// Register template providers from config
	for name, p := range cfg.Providers {
		e.providers[name] = provider.NewTemplateProvider(name, p.Template, p.Target)
	}

	return e
}

/**
 * ApplyProvider processes a specific provider's logic.
 */
func (e *Engine) ApplyProvider(name string) error {
	p, ok := e.providers[name]
	if !ok {
		return fmt.Errorf("provider %s not found in configuration", name)
	}

	// Extract only values for the template engine
	tokenValues := make(map[string]string)
	for k, v := range e.Config.Tokens {
		tokenValues[k] = v.Value
	}

	return p.Apply(tokenValues)
}

/**
 * ApplyAll processes all registered providers.
 */
func (e *Engine) ApplyAll() error {
	if len(e.providers) == 0 {
		fmt.Println("ℹ No providers configured. Add some in your config.yaml.")
		return nil
	}

	for name := range e.providers {
		if err := e.ApplyProvider(name); err != nil {
			return fmt.Errorf("failed to apply provider %s: %w", name, err)
		}
	}
	return nil
}
