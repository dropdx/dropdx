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
		
		// Collect all possible registry values
		var registryValues []string
		for _, regInfo := range v.Registries {
			if regInfo.Value != "" {
				registryValues = append(registryValues, regInfo.Value)
			}
		}

		// If main value is empty, use the first available registry value
		if val == "" && len(registryValues) > 0 {
			// Specific logic for npm to prefer npmjs.org
			if k == "npm" {
				for reg, regInfo := range v.Registries {
					if strings.Contains(reg, "npmjs.org") && regInfo.Value != "" {
						val = regInfo.Value
						break
					}
				}
			}
			// If still empty, just take the first one
			if val == "" {
				val = registryValues[0]
			}
		}

		if val != "" {
			tokens[k] = val
			tokens[k+"_token"] = val
		}

		// Add all registries as keys
		for reg, regInfo := range v.Registries {
			if regInfo.Value == "" {
				continue
			}
			tokens[reg] = regInfo.Value
			
			// Clean key (e.g. registry.npmjs.org)
			clean := reg
			clean = strings.TrimPrefix(clean, "https://")
			clean = strings.TrimPrefix(clean, "http://")
			clean = strings.TrimSuffix(clean, "/")
			tokens[clean] = regInfo.Value
			
			// Underscore version (e.g. registry_npmjs_org)
			underscore := strings.ReplaceAll(strings.ReplaceAll(clean, ".", "_"), ":", "_")
			tokens[underscore] = regInfo.Value
		}
	}
	return tokens
}
