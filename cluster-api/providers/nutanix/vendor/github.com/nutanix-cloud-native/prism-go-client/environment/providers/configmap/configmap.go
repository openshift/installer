// TODO implement provider sourcing settings from a Kubernetes config map.
// Config maps are usually mounted as directory and should not require access
// to k8s APIs.
package configmap

import "github.com/nutanix-cloud-native/prism-go-client/environment/types"

type provider struct{}

func (prov *provider) GetManagementEndpoint(
	topology types.Topology,
) (*types.ManagementEndpoint, error) {
	return nil, types.ErrNotFound
}

func (prov *provider) Get(topology types.Topology, key string) (
	interface{}, error,
) {
	return nil, types.ErrNotFound
}

func NewProvider() types.Provider {
	return &provider{}
}
