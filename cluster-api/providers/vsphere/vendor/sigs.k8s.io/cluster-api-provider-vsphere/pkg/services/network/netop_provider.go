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
	"context"
	"fmt"

	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

type netopNetworkProvider struct {
	client client.Client
}

// NetOpNetworkProvider returns a NetOp (VDS) Network Provider.
func NetOpNetworkProvider(client client.Client) services.NetworkProvider {
	return &netopNetworkProvider{
		client: client,
	}
}

// HasLoadBalancer is always true for the NetOp Network Provider.
func (np *netopNetworkProvider) HasLoadBalancer() bool {
	return true
}

// ProvisionClusterNetwork marks the ClusterNetworkReadyCondition true.
func (np *netopNetworkProvider) ProvisionClusterNetwork(_ context.Context, clusterCtx *vmware.ClusterContext) error {
	conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.ClusterNetworkReadyCondition)
	return nil
}

func (np *netopNetworkProvider) getDefaultClusterNetwork(ctx context.Context, clusterCtx *vmware.ClusterContext) (*netopv1.Network, error) {
	log := ctrl.LoggerFrom(ctx)

	networkWithLabel, err := np.getDefaultClusterNetworkWithLabel(ctx, clusterCtx, CAPVDefaultNetworkLabel)
	if networkWithLabel != nil && err == nil {
		return networkWithLabel, nil
	}

	log.Info("Falling back to legacy label to identify default network", "label", legacyDefaultNetworkLabel)
	return np.getDefaultClusterNetworkWithLabel(ctx, clusterCtx, legacyDefaultNetworkLabel)
}

func (np *netopNetworkProvider) getDefaultClusterNetworkWithLabel(ctx context.Context, clusterCtx *vmware.ClusterContext, label string) (*netopv1.Network, error) {
	labels := map[string]string{
		label: "true",
	}
	networkList := &netopv1.NetworkList{}
	err := np.client.List(ctx, networkList, client.InNamespace(clusterCtx.Cluster.Namespace), client.MatchingLabels(labels))
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

func (np *netopNetworkProvider) getClusterNetwork(ctx context.Context, clusterCtx *vmware.ClusterContext) (*netopv1.Network, error) {
	// A "NetworkName" can later be added to the Spec, but currently we only have a preselected default.
	return np.getDefaultClusterNetwork(ctx, clusterCtx)
}

// GetClusterNetworkName returns the name of the network for the passed cluster.
func (np *netopNetworkProvider) GetClusterNetworkName(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error) {
	network, err := np.getClusterNetwork(ctx, clusterCtx)
	if err != nil {
		return "", err
	}

	return network.Name, nil
}

// GetVMServiceAnnotations returns the name of the network in a map[string]string to allow usage in annotations.
func (np *netopNetworkProvider) GetVMServiceAnnotations(ctx context.Context, clusterCtx *vmware.ClusterContext) (map[string]string, error) {
	networkName, err := np.GetClusterNetworkName(ctx, clusterCtx)
	if err != nil {
		return nil, err
	}

	return map[string]string{NetOpNetworkNameAnnotation: networkName}, nil
}

// ConfigureVirtualMachine configures the NetworkInterfaces on a VM Operator virtual machine.
func (np *netopNetworkProvider) ConfigureVirtualMachine(ctx context.Context, clusterCtx *vmware.ClusterContext, vm *vmoprv1.VirtualMachine) error {
	network, err := np.getClusterNetwork(ctx, clusterCtx)
	if err != nil {
		return err
	}

	for _, vnif := range vm.Spec.NetworkInterfaces {
		if vnif.NetworkType == string(network.Spec.Type) && vnif.NetworkName == network.Name {
			// Expected network interface already exists.
			return nil
		}
	}

	vm.Spec.NetworkInterfaces = append(vm.Spec.NetworkInterfaces, vmoprv1.VirtualMachineNetworkInterface{
		NetworkName: network.Name,
		NetworkType: string(network.Spec.Type),
	})

	return nil
}

func (np *netopNetworkProvider) VerifyNetworkStatus(_ context.Context, _ *vmware.ClusterContext, obj runtime.Object) error {
	if _, ok := obj.(*netopv1.Network); !ok {
		return fmt.Errorf("expected Net Operator Network but got %T", obj)
	}

	// Network doesn't have a []Conditions but the specific network type pointed to by ProviderRef might.
	// The VSphereDistributedNetwork does but it is not currently populated by net-operator.

	return nil
}
