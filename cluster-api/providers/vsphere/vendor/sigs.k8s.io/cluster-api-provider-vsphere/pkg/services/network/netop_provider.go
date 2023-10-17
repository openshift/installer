/*
Copyright 2022 The Kubernetes Authors.

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

package network

import (
	"fmt"

	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

type netopNetworkProvider struct {
	client client.Client
}

func NetOpNetworkProvider(client client.Client) services.NetworkProvider {
	return &netopNetworkProvider{
		client: client,
	}
}

func (np *netopNetworkProvider) HasLoadBalancer() bool {
	return true
}

func (np *netopNetworkProvider) ProvisionClusterNetwork(ctx *vmware.ClusterContext) error {
	conditions.MarkTrue(ctx.VSphereCluster, infrav1.ClusterNetworkReadyCondition)
	return nil
}

func (np *netopNetworkProvider) getDefaultClusterNetwork(ctx *vmware.ClusterContext) (*netopv1.Network, error) {
	networkWithLabel, err := np.getDefaultClusterNetworkWithLabel(ctx, CAPVDefaultNetworkLabel)
	if networkWithLabel != nil && err == nil {
		return networkWithLabel, nil
	}

	ctx.Logger.Info("falling back to legacy label to identify default network", "label", legacyDefaultNetworkLabel)
	return np.getDefaultClusterNetworkWithLabel(ctx, legacyDefaultNetworkLabel)
}

func (np *netopNetworkProvider) getDefaultClusterNetworkWithLabel(ctx *vmware.ClusterContext, label string) (*netopv1.Network, error) {
	labels := map[string]string{
		label: "true",
	}
	networkList := &netopv1.NetworkList{}
	err := np.client.List(ctx, networkList, client.InNamespace(ctx.Cluster.Namespace), client.MatchingLabels(labels))
	if err != nil {
		return nil, err
	}

	switch len(networkList.Items) {
	case 0:
		return nil, fmt.Errorf("no default network found with label %s", label)
	case 1:
		return &networkList.Items[0], nil
	default:
		return nil, fmt.Errorf("more than one network found with label %s: %d", label, len(networkList.Items))
	}
}

func (np *netopNetworkProvider) getClusterNetwork(ctx *vmware.ClusterContext) (*netopv1.Network, error) {
	// A "NetworkName" can later be added to the Spec, but currently we only have a preselected default.
	return np.getDefaultClusterNetwork(ctx)
}

func (np *netopNetworkProvider) GetClusterNetworkName(ctx *vmware.ClusterContext) (string, error) {
	network, err := np.getClusterNetwork(ctx)
	if err != nil {
		return "", err
	}

	return network.Name, nil
}

func (np *netopNetworkProvider) GetVMServiceAnnotations(ctx *vmware.ClusterContext) (map[string]string, error) {
	networkName, err := np.GetClusterNetworkName(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{NetOpNetworkNameAnnotation: networkName}, nil
}

func (np *netopNetworkProvider) ConfigureVirtualMachine(ctx *vmware.ClusterContext, vm *vmopv1.VirtualMachine) error {
	network, err := np.getClusterNetwork(ctx)
	if err != nil {
		return err
	}

	for _, vnif := range vm.Spec.NetworkInterfaces {
		if vnif.NetworkType == string(network.Spec.Type) && vnif.NetworkName == network.Name {
			// Expected network interface already exists.
			return nil
		}
	}

	vm.Spec.NetworkInterfaces = append(vm.Spec.NetworkInterfaces, vmopv1.VirtualMachineNetworkInterface{
		NetworkName: network.Name,
		NetworkType: string(network.Spec.Type),
	})

	return nil
}

func (np *netopNetworkProvider) VerifyNetworkStatus(_ *vmware.ClusterContext, obj runtime.Object) error {
	if _, ok := obj.(*netopv1.Network); !ok {
		return fmt.Errorf("expected Net Operator Network but got %T", obj)
	}

	// Network doesn't have a []Conditions but the specific network type pointed to by ProviderRef might.
	// The VSphereDistributedNetwork does but it is not currently populated by net-operator.

	return nil
}
