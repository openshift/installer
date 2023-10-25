/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manager

import (
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/network"
)

const (
	NSXNetworkProvider     = "NSX"
	VDSNetworkProvider     = "vsphere-network"
	DummyLBNetworkProvider = "DummyLBNetworkProvider"
)

// GetNetworkProvider will return a network provider instance based on the environment
// the cfg is used to initialize a client that talks directly to api-server without using the cache.
func GetNetworkProvider(ctx *context.ControllerManagerContext) (services.NetworkProvider, error) {
	switch ctx.NetworkProvider {
	case NSXNetworkProvider:
		// TODO: disableFirewall not configurable
		ctx.Logger.Info("Pick NSX-T network provider")
		return network.NsxtNetworkProvider(ctx.Client, "false"), nil
	case VDSNetworkProvider:
		ctx.Logger.Info("Pick NetOp (VDS) network provider")
		return network.NetOpNetworkProvider(ctx.Client), nil
	case DummyLBNetworkProvider:
		ctx.Logger.Info("Pick Dummy network provider")
		return network.DummyLBNetworkProvider(), nil
	default:
		ctx.Logger.Info("NetworkProvider not set. Pick Dummy network provider")
		return network.DummyNetworkProvider(), nil
	}
}
