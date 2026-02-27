package core

import (
	"fmt"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/dropdx/dropdx/packages/provider"
	"github.com/fatih/color"
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

	return p.Apply(e.getTemplateTokens())
}

/**
 * ApplyAll processes all registered providers.
 */
func (e *Engine) ApplyAll() error {
	if len(e.providers) == 0 {
		fmt.Printf("%s No providers configured. Add some in your config.yaml.\n", color.HiBlackString("ℹ"))
		return nil
	}

	tokens := e.getTemplateTokens()
	for name := range e.providers {
		if err := e.providers[name].Apply(tokens); err != nil {
			return fmt.Errorf("failed to apply provider %s: %w", name, err)
		}
	}
	return nil
}

func (e *Engine) getTemplateTokens() map[string]string {
	tokens := make(map[string]string)
	for k, v := range e.Config.Tokens {
		if v.Value != "" {
			tokens[k] = v.Value
		}
		
		// If the main value is empty but we have registries, 
		// use the first one as the default for the provider name key
		if v.Value == "" && len(v.Registries) > 0 {
			// Try to find npmjs first for npm provider
			if k == "npm" {
				if reg, ok := v.Registries["https://registry.npmjs.org/"]; ok {
					tokens[k] = reg.Value
				}
			}
			// If still empty, take the first available
			if tokens[k] == "" {
				for _, regInfo := range v.Registries {
					tokens[k] = regInfo.Value
					break
				}
			}
		}

		// Also add registries to the map so they can be accessed via URL
		for reg, regInfo := range v.Registries {
			tokens[reg] = regInfo.Value
		}
	}
	return tokens
}
