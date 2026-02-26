package provider

/**
 * Provider is the interface that defines how a service configuration is applied.
 */
type Provider interface {
	// Name returns the identifier for this provider.
	Name() string
	// Apply processes the provider's logic with the given tokens.
	Apply(tokens map[string]string) error
}
