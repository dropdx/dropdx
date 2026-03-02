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
	
	// Create a temporary map to store results to handle aliases easily
	type tokenSet struct {
		val   string
		items []config.TokenInfo
		regs  map[string]config.TokenInfo
	}
	results := make(map[string]tokenSet)

	for k, v := range e.Config.Tokens {
		ts := tokenSet{
			val:   v.Value,
			items: v.Items,
			regs:  v.Registries,
		}
		
		// Pick first item as default if main value is empty
		if ts.val == "" && len(ts.items) > 0 {
			ts.val = ts.items[0].Value
		}
		
		results[k] = ts
	}

	// Handle 'gh' as alias for 'github' if missing
	if _, hasGh := results["gh"]; !hasGh {
		if gt, hasGithub := results["github"]; hasGithub {
			results["gh"] = gt
		}
	}

	for k, ts := range results {
		val := ts.val
		
		// Collect all possible registry values
		var registryValues []string
		for _, regInfo := range ts.regs {
			if regInfo.Value != "" {
				registryValues = append(registryValues, regInfo.Value)
			}
		}

		// If main value is still empty, use the first available registry value
		if val == "" && len(registryValues) > 0 {
			// Specific logic for npm to prefer npmjs.org
			if k == "npm" {
				for reg, regInfo := range ts.regs {
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

		// If it's a list, expose individual ones with names and indices
		for i, item := range ts.items {
			if item.Value == "" {
				continue
			}
			
			// By index: github_0, github_1
			tokens[fmt.Sprintf("%s_%d", k, i)] = item.Value
			
			// By name: github_classic, github_fine_grained
			if item.Name != "" {
				cleanName := strings.ToLower(strings.ReplaceAll(item.Name, " ", "_"))
				tokens[fmt.Sprintf("%s_%s", k, cleanName)] = item.Value
			}
		}

		// Add all registries as keys
		for reg, regInfo := range ts.regs {
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
