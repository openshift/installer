/*
Copyright 2024 The Kubernetes Authors.

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
	nsxvpcv1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1common "github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

// nsxtVPCNetworkProvider provisions nsx-vpc type cluster network.
type nsxtVPCNetworkProvider struct {
	client client.Client
}

// NSXTVpcNetworkProvider returns an instance of nsx-vpc type network provider.
func NSXTVpcNetworkProvider(client client.Client) services.NetworkProvider {
	return &nsxtVPCNetworkProvider{
		client: client,
	}
}

func (vp *nsxtVPCNetworkProvider) HasLoadBalancer() bool {
	return true
}

func (vp *nsxtVPCNetworkProvider) SupportsVMReadinessProbe() bool {
	// Note: The control plane VM network is private with nsx-vpc and
	// readiness probe would fail. Therefore, nsxvpcNetworkProvider
	// doesn't support VM readiness probe.
	return false
}

// verifyNsxtVpcSubnetSetStatus checks the status conditions of a given SubnetSet within a cluster context.
// If the subnet isn't ready, it is marked as false, and the function returns an error.
// If the subnet is ready, the function updates the VSphereCluster with a "true" status and returns nil.
func (vp *nsxtVPCNetworkProvider) verifyNsxtVpcSubnetSetStatus(vspherecluster *vmwarev1.VSphereCluster, subnetset *nsxvpcv1.SubnetSet) error {
	clusterName := vspherecluster.Name
	namespace := vspherecluster.Namespace
	hasReadyCondition := false

	for _, condition := range subnetset.Status.Conditions {
		if condition.Type != nsxvpcv1.Ready {
			continue
		}
		hasReadyCondition = true
		if condition.Status != corev1.ConditionTrue {
			conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, condition.Message)
			return errors.Errorf("subnetset ready status is: '%s' in cluster %s. reason: %s, message: %s",
				condition.Status, types.NamespacedName{Namespace: namespace, Name: clusterName}, condition.Reason, condition.Message)
		}
	}

	if !hasReadyCondition {
		conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, "No Ready status for SubnetSet")
		return errors.Errorf("subnetset ready status in cluster %s has not been set", types.NamespacedName{Namespace: namespace, Name: clusterName})
	}

	conditions.MarkTrue(vspherecluster, vmwarev1.ClusterNetworkReadyCondition)
	return nil
}

// VerifyNetworkStatus checks if the given runtime object is of type SubnetSet.
// If it is, then it calls verifyNsxVpcSubnetSetStatus with the SubnetSet to verify its status.
// If it's not, it returns an error.
func (vp *nsxtVPCNetworkProvider) VerifyNetworkStatus(_ context.Context, clusterCtx *vmware.ClusterContext, obj runtime.Object) error {
	subnetset, ok := obj.(*nsxvpcv1.SubnetSet)
	if !ok {
		return fmt.Errorf("expected NSX VPC SubnetSet but got %T", obj)
	}

	return vp.verifyNsxtVpcSubnetSetStatus(clusterCtx.VSphereCluster, subnetset)
}

// ProvisionClusterNetwork provisions a new network in the context of a given cluster.
// It constructs a new SubnetSet and attempts to create or patch it on the cluster.
// If it fails to do so, it marks the status of the VSphereCluster as false and returns an error.
// If it succeeds, it calls verifyNsxVpcSubnetSetStatus to verify the status of the newly created/patched SubnetSet.
func (vp *nsxtVPCNetworkProvider) ProvisionClusterNetwork(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	cluster := clusterCtx.VSphereCluster
	networkNamespace := cluster.Namespace
	networkName := cluster.Name
	log = log.WithValues("SubnetSet", klog.KRef(networkNamespace, networkName))

	log.Info("Provisioning ")
	defer log.Info("Finished provisioning")

	subnetset := &nsxvpcv1.SubnetSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: networkNamespace,
			Name:      networkName,
		},
		Spec: nsxvpcv1.SubnetSetSpec{},
	}

	_, err := ctrlutil.CreateOrPatch(ctx, vp.client, subnetset, func() error {
		if err := ctrlutil.SetOwnerReference(
			clusterCtx.VSphereCluster,
			subnetset,
			vp.client.Scheme(),
		); err != nil {
			return errors.Wrapf(err, "error setting %s as owner of %s", klog.KObj(clusterCtx.VSphereCluster), klog.KObj(subnetset))
		}

		return nil
	})
	if err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return errors.Wrap(err, "Failed to provision network")
	}

	return vp.verifyNsxtVpcSubnetSetStatus(clusterCtx.VSphereCluster, subnetset)
}

// GetClusterNetworkName returns the name of a valid cluster network if one exists.
func (vp *nsxtVPCNetworkProvider) GetClusterNetworkName(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error) {
	subnetset := &nsxvpcv1.SubnetSet{}
	cluster := clusterCtx.VSphereCluster
	namespacedName := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      cluster.Name,
	}
	if err := vp.client.Get(ctx, namespacedName, subnetset); err != nil {
		return "", err
	}
	return namespacedName.Name, nil
}

// The GetVMServiceAnnotations method always returns an empty map representing annotations.
func (vp *nsxtVPCNetworkProvider) GetVMServiceAnnotations(_ context.Context, _ *vmware.ClusterContext) (map[string]string, error) {
	// The value of the annotation lb.iaas.vmware.com/enable-endpoint-health-check is expected to be an empty string.
	return map[string]string{AnnotationEnableEndpointHealthCheckKey: ""}, nil
}

// ConfigureVirtualMachine configures a VirtualMachine object based on the networking configuration.
func (vp *nsxtVPCNetworkProvider) ConfigureVirtualMachine(_ context.Context, clusterCtx *vmware.ClusterContext, vm *vmoprv1.VirtualMachine) error {
	networkName := clusterCtx.VSphereCluster.Name
	if vm.Spec.Network == nil {
		vm.Spec.Network = &vmoprv1.VirtualMachineNetworkSpec{}
	}
	for _, vnif := range vm.Spec.Network.Interfaces {
		if vnif.Network.TypeMeta.GroupVersionKind() == NetworkGVKNSXTVPC && vnif.Network.Name == networkName {
			// expected network interface is already found
			return nil
		}
	}

	vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmoprv1.VirtualMachineNetworkInterfaceSpec{
		Name: fmt.Sprintf("eth%d", len(vm.Spec.Network.Interfaces)),
		Network: vmoprv1common.PartialObjectRef{
			TypeMeta: metav1.TypeMeta{
				Kind:       NetworkGVKNSXTVPC.Kind,
				APIVersion: NetworkGVKNSXTVPC.GroupVersion().String(),
			},
			Name: networkName,
		},
	})
	return nil
}
