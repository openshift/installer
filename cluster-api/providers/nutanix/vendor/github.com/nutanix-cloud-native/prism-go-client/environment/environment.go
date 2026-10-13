/*
 A Kubernetes cluster interacts with its underlying infrastructure or its
 environment in general. Client code should be shielded from how settings
 specific to an environment is injected into the k8s cluster as this mechanism
 might change.

 This package provides abstraction for enviroment and all its sources to evolve
 indpendently from clients like k8s extensions or operators. Environment here
 should be understood as least common denominator across any k8s extensions.

 Multiple Kubernetes clusters might share certain properties of an environment
 like the project. However idea is that each k8s cluster environemt is unique
 to allow infrastructure resources created on behalf of a k8s cluster to be
 associated with this particular k8s cluster.
*/

package environment

import (
	"github.com/nutanix-cloud-native/prism-go-client/environment/providers/configmap"
	"github.com/nutanix-cloud-native/prism-go-client/environment/providers/local"
	"github.com/nutanix-cloud-native/prism-go-client/environment/providers/secretdir"
	"github.com/nutanix-cloud-native/prism-go-client/environment/types"
)

// Env is current environment. Environment is a singleton as k8s cluster
// can't span multiple environment. Providers of environment hide behind
// this interface.
// Clients of this package are free to reconfigure environment.
var Env = NewEnvironment(
	local.NewProvider(),     // Give local provider highest priority.
	secretdir.NewProvider(), // Secret and config map env shouldn't intersect,
	configmap.NewProvider(), // hence order doesn't matter.
)

// GetManagementEndpoint is convenience function
func GetManagementEndpoint(topology types.Topology) (
	*types.ManagementEndpoint, error,
) {
	return Env.GetManagementEndpoint(topology)
}

// Get is convenience function
func Get(topology types.Topology, key string) (interface{}, error) {
	return Env.Get(topology, key)
}

// env construct environment from various sources
type env struct {
	providers []types.Provider
}

// GetManagementEndpoint implements Environment interface
func (e *env) GetManagementEndpoint(topology types.Topology) (
	*types.ManagementEndpoint, error,
) {
	for _, provider := range e.providers {
		if me, err := provider.GetManagementEndpoint(topology); err == nil ||
			err != types.ErrNotFound {
			return me, err
		}
	}
	return nil, types.ErrNotFound
}

// Get implements Environment interface
func (e *env) Get(topology types.Topology, key string) (interface{}, error) {
	for _, provider := range e.providers {
		if val, err := provider.Get(topology, key); err == nil ||
			err != types.ErrNotFound {
			return val, err
		}
	}
	return nil, types.ErrNotFound
}

// NewEnvironment constructs new environment from given list of
// providers. The order of providers is honored during lookup.
// This function is meant to construct Env as environment is a
// singleton.
func NewEnvironment(p ...types.Provider) types.Environment {
	return &env{
		providers: p,
	}
}
