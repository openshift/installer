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

	"github.com/pkg/errors"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// nsxtNetworkProvider provision nsx-t type cluster network.
type nsxtNetworkProvider struct {
	client    client.Client
	disableFW string // "true" means disable firewalls on GC network
}

// NsxtNetworkProvider returns an instance of nsx-t type network provider.
func NsxtNetworkProvider(client client.Client, disableFW string) services.NetworkProvider {
	return &nsxtNetworkProvider{
		client:    client,
		disableFW: disableFW,
	}
}

func (np *nsxtNetworkProvider) HasLoadBalancer() bool {
	return true
}

// GetNSXTVirtualNetworkName returns the name of the NSX-T vnet object.
func GetNSXTVirtualNetworkName(clusterName string) string {
	return fmt.Sprintf("%s-vnet", clusterName)
}

func (np *nsxtNetworkProvider) verifyNSXTVirtualNetworkStatus(ctx *vmware.ClusterContext, vnet *ncpv1.VirtualNetwork) error {
	clusterName := ctx.VSphereCluster.Name
	namespace := ctx.VSphereCluster.Namespace
	for _, condition := range vnet.Status.Conditions {
		if condition.Type == "Ready" && condition.Status != "True" {
			conditions.MarkFalse(ctx.VSphereCluster, infrav1.ClusterNetworkReadyCondition, infrav1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, condition.Message)
			return errors.Errorf("virtual network ready status is: '%s' in cluster %s. reason: %s, message: %s",
				condition.Status, types.NamespacedName{Namespace: namespace, Name: clusterName}, condition.Reason, condition.Message)
		}
	}

	conditions.MarkTrue(ctx.VSphereCluster, infrav1.ClusterNetworkReadyCondition)
	return nil
}

func (np *nsxtNetworkProvider) VerifyNetworkStatus(ctx *vmware.ClusterContext, obj runtime.Object) error {
	vnet, ok := obj.(*ncpv1.VirtualNetwork)
	if !ok {
		return fmt.Errorf("expected NCP VirtualNetwork but got %T", obj)
	}

	return np.verifyNSXTVirtualNetworkStatus(ctx, vnet)
}

func (np *nsxtNetworkProvider) ProvisionClusterNetwork(ctx *vmware.ClusterContext) error {
	cluster := ctx.VSphereCluster
	clusterKey := types.NamespacedName{Name: cluster.Name, Namespace: cluster.Namespace}

	ctx.Logger.V(2).Info("Provisioning ", "vnet", GetNSXTVirtualNetworkName(cluster.Name), "namespace", cluster.Namespace)
	defer ctx.Logger.V(2).Info("Finished provisioning", "vnet", GetNSXTVirtualNetworkName(cluster.Name), "namespace", cluster.Namespace)

	vnet := &ncpv1.VirtualNetwork{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      GetNSXTVirtualNetworkName(cluster.Name),
		},
	}

	_, err := ctrlutil.CreateOrPatch(ctx, np.client, vnet, func() error {
		// add or update vnet spec only if FW is enabled and if WhitelistSourceRanges is empty
		if np.disableFW != "true" && vnet.Spec.WhitelistSourceRanges == "" {
			supportFW, err := util.NCPSupportFW(ctx, np.client)
			if err != nil {
				ctx.Logger.Error(err, "failed to check if NCP supports firewall rules enforcement on GC T1 router")
				return err
			}
			// specify whitelist_source_ranges if needed and if NCP supports it
			if supportFW {
				// Find system namespace snat ip
				systemNSSnatIP, err := util.GetNamespaceNetSnatIP(ctx, np.client, SystemNamespace)
				if err != nil {
					ctx.Logger.Error(err, "failed to get Snat IP for kube-system")
					return err
				}
				ctx.Logger.V(4).Info("got system namespace snat ip",
					"cluster", clusterKey, "ip", systemNSSnatIP)

				// WhitelistSourceRanges accept cidrs only
				vnet.Spec.WhitelistSourceRanges = systemNSSnatIP + "/32"
			}
		}

		if err := ctrlutil.SetOwnerReference(
			ctx.VSphereCluster,
			vnet,
			ctx.Scheme,
		); err != nil {
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				ctx.VSphereCluster.Namespace,
				ctx.VSphereCluster.Name,
				vnet.Namespace,
				vnet.Name,
			)
		}

		return nil
	})
	if err != nil {
		conditions.MarkFalse(ctx.VSphereCluster, infrav1.ClusterNetworkReadyCondition, infrav1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		ctx.Logger.V(2).Info("Failed to provision network", "cluster", clusterKey)
		return err
	}

	return np.verifyNSXTVirtualNetworkStatus(ctx, vnet)
}

// Returns the name of a valid cluster network if one exists.
func (np *nsxtNetworkProvider) GetClusterNetworkName(ctx *vmware.ClusterContext) (string, error) {
	vnet := &ncpv1.VirtualNetwork{}
	cluster := ctx.VSphereCluster
	namespacedName := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      GetNSXTVirtualNetworkName(cluster.Name),
	}
	if err := np.client.Get(ctx, namespacedName, vnet); err != nil {
		return "", err
	}
	return namespacedName.Name, nil
}

func (np *nsxtNetworkProvider) GetVMServiceAnnotations(ctx *vmware.ClusterContext) (map[string]string, error) {
	vnetName, err := np.GetClusterNetworkName(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{NSXTVNetSelectorKey: vnetName}, nil
}

// ConfigureVirtualMachine configures a VirtualMachine object based on the networking configuration.
func (np *nsxtNetworkProvider) ConfigureVirtualMachine(ctx *vmware.ClusterContext, vm *vmopv1.VirtualMachine) error {
	nsxtClusterNetworkName := GetNSXTVirtualNetworkName(ctx.VSphereCluster.Name)
	for _, vnif := range vm.Spec.NetworkInterfaces {
		if vnif.NetworkType == NSXTTypeNetwork && vnif.NetworkName == nsxtClusterNetworkName {
			// expected network interface is already found
			return nil
		}
	}
	vm.Spec.NetworkInterfaces = append(vm.Spec.NetworkInterfaces, vmopv1.VirtualMachineNetworkInterface{
		NetworkName: nsxtClusterNetworkName,
		NetworkType: NSXTTypeNetwork,
	})
	return nil
}
