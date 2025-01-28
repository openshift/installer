package client

import (
	"github.com/nutanix-cloud-native/prism-go-client/environment/types"
	v3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	v4 "github.com/nutanix-cloud-native/prism-go-client/v4"

	"github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
)

// NutanixClientCache is the cache of prism clients to be shared across the different controllers
var NutanixClientCache = v3.NewClientCache(v3.WithSessionAuth(true))

// NutanixClientCacheV4 is the cache of prism clients to be shared across the different controllers
var NutanixClientCacheV4 = v4.NewClientCache(v4.WithSessionAuth(true))

// CacheParams is the struct that implements ClientCacheParams interface from prism-go-client
type CacheParams struct {
	NutanixCluster          *v1beta1.NutanixCluster
	PrismManagementEndpoint *types.ManagementEndpoint
}

// Key is the namespace/name of te NutanixCluster CR
func (c *CacheParams) Key() string {
	return c.NutanixCluster.GetNamespacedName()
}

// ManagementEndpoint returns the management endpoint of the NutanixCluster CR
func (c *CacheParams) ManagementEndpoint() types.ManagementEndpoint {
	return *c.PrismManagementEndpoint
}
