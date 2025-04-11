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

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1common "github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
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

func (np *nsxtNetworkProvider) SupportsVMReadinessProbe() bool {
	return true
}

// GetNSXTVirtualNetworkName returns the name of the NSX-T vnet object.
func GetNSXTVirtualNetworkName(clusterName string) string {
	return fmt.Sprintf("%s-vnet", clusterName)
}

func (np *nsxtNetworkProvider) verifyNSXTVirtualNetworkStatus(vspherecluster *vmwarev1.VSphereCluster, vnet *ncpv1.VirtualNetwork) error {
	clusterName := vspherecluster.Name
	namespace := vspherecluster.Namespace
	hasReadyCondition := false

	for _, condition := range vnet.Status.Conditions {
		if condition.Type != "Ready" {
			continue
		}
		hasReadyCondition = true
		if condition.Status != "True" {
			conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, condition.Message)
			return errors.Errorf("virtual network ready status is: '%s' in cluster %s. reason: %s, message: %s",
				condition.Status, types.NamespacedName{Namespace: namespace, Name: clusterName}, condition.Reason, condition.Message)
		}
	}

	if !hasReadyCondition {
		conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, "No Ready status for virtual network")
		return errors.Errorf("virtual network ready status in cluster %s has not been set", types.NamespacedName{Namespace: namespace, Name: clusterName})
	}

	conditions.MarkTrue(vspherecluster, vmwarev1.ClusterNetworkReadyCondition)
	return nil
}

func (np *nsxtNetworkProvider) VerifyNetworkStatus(_ context.Context, clusterCtx *vmware.ClusterContext, obj runtime.Object) error {
	vnet, ok := obj.(*ncpv1.VirtualNetwork)
	if !ok {
		return fmt.Errorf("expected NCP VirtualNetwork but got %T", obj)
	}

	return np.verifyNSXTVirtualNetworkStatus(clusterCtx.VSphereCluster, vnet)
}

func (np *nsxtNetworkProvider) ProvisionClusterNetwork(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	cluster := clusterCtx.VSphereCluster

	log.Info("Provisioning ", "vnet", GetNSXTVirtualNetworkName(cluster.Name))
	defer log.Info("Finished provisioning", "vnet", GetNSXTVirtualNetworkName(cluster.Name))

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
				return errors.Wrap(err, "failed to check if NCP supports firewall rules enforcement on GC T1 router")
			}
			// specify whitelist_source_ranges if needed and if NCP supports it
			if supportFW {
				// Find system namespace snat ip
				systemNSSnatIP, err := util.GetNamespaceNetSnatIP(ctx, np.client, SystemNamespace)
				if err != nil {
					return errors.Wrap(err, "failed to get Snat IP for kube-system")
				}
				log.V(4).Info("Got system namespace snat ip", "ip", systemNSSnatIP)

				// WhitelistSourceRanges accept cidrs only
				vnet.Spec.WhitelistSourceRanges = systemNSSnatIP + "/32"
			}
		}

		if err := ctrlutil.SetOwnerReference(
			clusterCtx.VSphereCluster,
			vnet,
			np.client.Scheme(),
		); err != nil {
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				clusterCtx.VSphereCluster.Namespace,
				clusterCtx.VSphereCluster.Name,
				vnet.Namespace,
				vnet.Name,
			)
		}

		return nil
	})
	if err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return errors.Wrap(err, "Failed to provision network")
	}

	return np.verifyNSXTVirtualNetworkStatus(clusterCtx.VSphereCluster, vnet)
}

// GetClusterNetworkName returns the name of a valid cluster network if one exists.
func (np *nsxtNetworkProvider) GetClusterNetworkName(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error) {
	vnet := &ncpv1.VirtualNetwork{}
	cluster := clusterCtx.VSphereCluster
	namespacedName := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      GetNSXTVirtualNetworkName(cluster.Name),
	}
	if err := np.client.Get(ctx, namespacedName, vnet); err != nil {
		return "", err
	}
	return namespacedName.Name, nil
}

func (np *nsxtNetworkProvider) GetVMServiceAnnotations(ctx context.Context, clusterCtx *vmware.ClusterContext) (map[string]string, error) {
	vnetName, err := np.GetClusterNetworkName(ctx, clusterCtx)
	if err != nil {
		return nil, err
	}

	return map[string]string{NSXTVNetSelectorKey: vnetName}, nil
}

// ConfigureVirtualMachine configures a VirtualMachine object based on the networking configuration.
func (np *nsxtNetworkProvider) ConfigureVirtualMachine(_ context.Context, clusterCtx *vmware.ClusterContext, vm *vmoprv1.VirtualMachine) error {
	nsxtClusterNetworkName := GetNSXTVirtualNetworkName(clusterCtx.VSphereCluster.Name)
	if vm.Spec.Network == nil {
		vm.Spec.Network = &vmoprv1.VirtualMachineNetworkSpec{}
	}
	for _, vnif := range vm.Spec.Network.Interfaces {
		if vnif.Network.TypeMeta.GroupVersionKind() == NetworkGVKNSXT && vnif.Network.Name == nsxtClusterNetworkName {
			// expected network interface is already found
			return nil
		}
	}
	vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmoprv1.VirtualMachineNetworkInterfaceSpec{
		Name: fmt.Sprintf("eth%d", len(vm.Spec.Network.Interfaces)),
		Network: vmoprv1common.PartialObjectRef{
			TypeMeta: metav1.TypeMeta{
				Kind:       NetworkGVKNSXT.Kind,
				APIVersion: NetworkGVKNSXT.GroupVersion().String(),
			},
			Name: nsxtClusterNetworkName,
		},
	})
	return nil
}
