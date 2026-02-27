package core

import (
	"fmt"
	"strings"

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
		val := v.Value
		
		// If main value is empty, try to find a default from registries
		if val == "" && len(v.Registries) > 0 {
			// 1. Try common npm registry specifically for npm provider
			if k == "npm" {
				for reg, regInfo := range v.Registries {
					if strings.Contains(reg, "npmjs.org") && regInfo.Value != "" {
						val = regInfo.Value
						break
					}
				}
			}
			// 2. Fallback: take the first non-empty value from any registry
			if val == "" {
				for _, regInfo := range v.Registries {
					if regInfo.Value != "" {
						val = regInfo.Value
						break
					}
				}
			}
		}

		if val != "" {
			tokens[k] = val
			// Also add a suffixed version just in case (e.g. npm_token)
			tokens[k+"_token"] = val
		}

		// Also add all registries as keys
		for reg, regInfo := range v.Registries {
			if regInfo.Value == "" {
				continue
			}
			// Full URL key
			tokens[reg] = regInfo.Value
			
			// Simplified key (e.g. registry.npmjs.org)
			simpleReg := reg
			simpleReg = strings.TrimPrefix(simpleReg, "https://")
			simpleReg = strings.TrimPrefix(simpleReg, "http://")
			simpleReg = strings.TrimSuffix(simpleReg, "/")
			tokens[simpleReg] = regInfo.Value
		}
	}
	return tokens
}
