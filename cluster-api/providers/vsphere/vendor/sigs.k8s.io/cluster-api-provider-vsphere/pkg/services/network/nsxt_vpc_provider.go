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
	"k8s.io/utils/ptr"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
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
			v1beta1conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", condition.Message)
			v1beta2conditions.Set(vspherecluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterNetworkNotReadyV1Beta2Reason,
				Message: condition.Message,
			})
			return errors.Errorf("subnetset ready status is: '%s' in cluster %s. reason: %s, message: %s",
				condition.Status, types.NamespacedName{Namespace: namespace, Name: clusterName}, condition.Reason, condition.Message)
		}
	}

	if !hasReadyCondition {
		v1beta1conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1beta1.ConditionSeverityWarning, "No Ready status for SubnetSet")
		v1beta2conditions.Set(vspherecluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterNetworkNotReadyV1Beta2Reason,
			Message: "No Ready status for SubnetSet",
		})
		return errors.Errorf("subnetset ready status in cluster %s has not been set", types.NamespacedName{Namespace: namespace, Name: clusterName})
	}

	v1beta1conditions.MarkTrue(vspherecluster, vmwarev1.ClusterNetworkReadyCondition)
	v1beta2conditions.Set(vspherecluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterNetworkReadyV1Beta2Reason,
	})
	return nil
}

func createSubnetSet(clusterCtx *vmware.ClusterContext) bool {
	return ptr.Deref(clusterCtx.VSphereCluster.Spec.Network.NSXVPC.CreateSubnetSet, true)
}

// VerifyNetworkStatus checks if the given runtime object is of type SubnetSet.
// If it is, then it calls verifyNsxVpcSubnetSetStatus with the SubnetSet to verify its status.
// If it's not, it returns an error.
func (vp *nsxtVPCNetworkProvider) VerifyNetworkStatus(ctx context.Context, clusterCtx *vmware.ClusterContext, obj runtime.Object) error {
	log := ctrl.LoggerFrom(ctx)
	if !createSubnetSet(clusterCtx) {
		log.V(5).Info("Skipping SubnetSet status check as CreateSubnetSet is false")
		return nil
	}
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

	if !createSubnetSet(clusterCtx) {
		log.Info("Skipping SubnetSet creation as CreateSubnetSet is false")
		v1beta1conditions.MarkTrue(cluster, vmwarev1.ClusterNetworkReadyCondition)
		v1beta2conditions.Set(cluster, metav1.Condition{
			Type:   vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: vmwarev1.VSphereClusterNetworkReadyV1Beta2Reason,
		})
		return nil
	}

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
		v1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.ClusterNetworkReadyCondition, vmwarev1.ClusterNetworkProvisionFailedReason, clusterv1beta1.ConditionSeverityWarning, "%v", err)
		v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterNetworkNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
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
func (vp *nsxtVPCNetworkProvider) ConfigureVirtualMachine(_ context.Context, clusterCtx *vmware.ClusterContext, machine *vmwarev1.VSphereMachine, vm *vmoprv1.VirtualMachine) error {
	vm.Spec.Network = &vmoprv1.VirtualMachineNetworkSpec{}

	// Set the VM primary interface
	if createSubnetSet(clusterCtx) {
		if machine.Spec.Network.Interfaces.Primary.IsDefined() {
			return errors.New("primary interface can not be configured when createSubnetSet is true")
		}
		networkName := clusterCtx.VSphereCluster.Name
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmoprv1.VirtualMachineNetworkInterfaceSpec{
			Name: PrimaryInterfaceName,
			Network: vmoprv1common.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       NetworkGVKNSXTVPCSubnetSet.Kind,
					APIVersion: NetworkGVKNSXTVPCSubnetSet.GroupVersion().String(),
				},
				Name: networkName,
			},
		})
	} else {
		if !machine.Spec.Network.Interfaces.Primary.IsDefined() {
			return errors.New("primary interface must be configured when createSubnetSet is false")
		}
		primary := machine.Spec.Network.Interfaces.Primary
		var mtu *int64
		if primary.MTU != 0 {
			mtu = ptr.To(int64(primary.MTU))
		}
		vmInterface := vmoprv1.VirtualMachineNetworkInterfaceSpec{
			Name: PrimaryInterfaceName,
			Network: vmoprv1common.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       primary.Network.Kind,
					APIVersion: primary.Network.APIVersion,
				},
				Name: primary.Network.Name,
			},
			MTU: mtu,
		}
		setRoutes(&vmInterface, primary.Routes)
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmInterface)
	}

	// Set the VM secondary interfaces
	setVMSecondaryInterfaces(machine, vm)
	return nil
}

func setRoutes(vmInterface *vmoprv1.VirtualMachineNetworkInterfaceSpec, routes []vmwarev1.RouteSpec) {
	for _, route := range routes {
		vmInterface.Routes = append(vmInterface.Routes, vmoprv1.VirtualMachineNetworkRouteSpec{
			To:  route.To,
			Via: route.Via,
		})
	}
}

func setVMSecondaryInterfaces(machine *vmwarev1.VSphereMachine, vm *vmoprv1.VirtualMachine) {
	if len(machine.Spec.Network.Interfaces.Secondary) == 0 {
		return
	}
	for _, secondaryInterface := range machine.Spec.Network.Interfaces.Secondary {
		var mtu *int64
		if secondaryInterface.MTU != 0 {
			mtu = ptr.To(int64(secondaryInterface.MTU))
		}
		vmInterface := vmoprv1.VirtualMachineNetworkInterfaceSpec{
			Name: secondaryInterface.Name,
			Network: vmoprv1common.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       secondaryInterface.Network.Kind,
					APIVersion: secondaryInterface.Network.APIVersion,
				},
				Name: secondaryInterface.Network.Name,
			},
			MTU:      mtu,
			Gateway4: "None",
			Gateway6: "None",
		}
		setRoutes(&vmInterface, secondaryInterface.Routes)
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmInterface)
	}
}
