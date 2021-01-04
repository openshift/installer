package externalprovider

import (
	"errors"

	"github.com/openshift/installer/pkg/externalprovider/provider"
)

// ErrNoSuchProvider is returned or thrown in panic when the specified provider is not registered.
var ErrNoSuchProvider = errors.New("no provider with the specified name registered")

// Registry registers the providers to their names.
type Registry interface {
	// Register registers a provider.
	Register(provider provider.ExternalProvider)
	// Get returns a provider registered to a name or ErrNoSuchProvider if the provider is not registered.
	Get(name string) (provider.ExternalProvider, error)
	// MustGet returns a provider registered to a name or throws a panic with ErrNoSuchProvider if the provider is not
	// registered.
	MustGet(name string) provider.ExternalProvider
}

// NewRegistry creates a new copy of a Registry.
func NewRegistry() Registry {
	return &registry{
		providers: map[string]provider.ExternalProvider{},
	}
}

type registry struct {
	providers map[string]provider.ExternalProvider
}

func (r *registry) Register(provider provider.ExternalProvider) {
	r.providers[provider.Name()] = provider
}

func (r *registry) Get(name string) (provider.ExternalProvider, error) {
	if p, ok := r.providers[name]; ok {
		return p, nil
	}
	return nil, ErrNoSuchProvider
}

func (r *registry) MustGet(name string) provider.ExternalProvider {
	p, err := r.Get(name)
	if err != nil {
		panic(err)
	}
	return p
}
