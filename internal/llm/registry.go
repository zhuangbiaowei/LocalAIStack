package llm

import (
	"fmt"
	"sort"
)

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

func (r *Registry) Register(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider is nil")
	}
	name := provider.Name()
	if name == "" {
		return fmt.Errorf("provider has empty name")
	}
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}
	r.providers[name] = provider
	return nil
}

func (r *Registry) Provider(name string) (Provider, error) {
	provider, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not found", name)
	}
	return provider, nil
}

func (r *Registry) Providers() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
