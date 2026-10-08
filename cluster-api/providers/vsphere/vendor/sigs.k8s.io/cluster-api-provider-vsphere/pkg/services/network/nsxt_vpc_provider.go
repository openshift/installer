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
	"reflect"

	pkgerrors "github.com/pkg/errors"
	nsxvpcv1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	infrautilv1 "sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
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

func (vp *nsxtVPCNetworkProvider) SupportsIPv6DualStack() bool {
	return feature.Gates.Enabled(feature.IPv6DualStack)
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
			deprecatedv1beta1conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyV1Beta1Condition, vmwarev1.ClusterNetworkProvisionFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%s", condition.Message)
			conditions.Set(vspherecluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterNetworkReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterNetworkNotReadyReason,
				Message: condition.Message,
			})
			return pkgerrors.Errorf("subnetset ready status is: '%s' in cluster %s. reason: %s, message: %s",
				condition.Status, types.NamespacedName{Namespace: namespace, Name: clusterName}, condition.Reason, condition.Message)
		}
	}

	if !hasReadyCondition {
		deprecatedv1beta1conditions.MarkFalse(vspherecluster, vmwarev1.ClusterNetworkReadyV1Beta1Condition, vmwarev1.ClusterNetworkProvisionFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "No Ready status for SubnetSet")
		conditions.Set(vspherecluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterNetworkReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterNetworkNotReadyReason,
			Message: "No Ready status for SubnetSet",
		})
		return pkgerrors.Errorf("subnetset ready status in cluster %s has not been set", types.NamespacedName{Namespace: namespace, Name: clusterName})
	}

	deprecatedv1beta1conditions.MarkTrue(vspherecluster, vmwarev1.ClusterNetworkReadyV1Beta1Condition)
	conditions.Set(vspherecluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterNetworkReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterNetworkReadyReason,
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
		deprecatedv1beta1conditions.MarkTrue(cluster, vmwarev1.ClusterNetworkReadyV1Beta1Condition)
		conditions.Set(cluster, metav1.Condition{
			Type:   vmwarev1.VSphereClusterNetworkReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: vmwarev1.VSphereClusterNetworkReadyReason,
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

	if vp.SupportsIPv6DualStack() {
		ipFamily, err := infrautilv1.DetermineClusterIPFamily(clusterCtx.Cluster)
		if err != nil {
			return pkgerrors.Wrap(err, "failed to determine cluster IP family")
		}

		switch ipFamily {
		case infrautilv1.IPv4SingleStack:
			subnetset.Spec.IPAddressType = nsxvpcv1.IPAddressTypeIPv4
		case infrautilv1.IPv6SingleStack:
			subnetset.Spec.IPAddressType = nsxvpcv1.IPAddressTypeIPv6
		case infrautilv1.DualStackIPv4Primary, infrautilv1.DualStackIPv6Primary:
			subnetset.Spec.IPAddressType = nsxvpcv1.IPAddressTypeIPv4IPv6
		default:
			subnetset.Spec.IPAddressType = nsxvpcv1.IPAddressTypeIPv4
		}
	}

	subnetSetExists := true
	if err := vp.client.Get(ctx, client.ObjectKeyFromObject(subnetset), subnetset); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		subnetSetExists = false
	}
	originalSubnetSet := subnetset.DeepCopy()

	if err := ctrlutil.SetOwnerReference(
		clusterCtx.VSphereCluster,
		subnetset,
		vp.client.Scheme(),
	); err != nil {
		return pkgerrors.Wrapf(err, "error setting %s as owner of %s", klog.KObj(clusterCtx.VSphereCluster), klog.KObj(subnetset))
	}

	var err error
	if !subnetSetExists {
		err = vp.client.Create(ctx, subnetset)
	} else if !reflect.DeepEqual(originalSubnetSet, subnetset) {
		patch := client.MergeFrom(originalSubnetSet)
		err = vp.client.Patch(ctx, subnetset, patch)
	}
	if err != nil {
		deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.ClusterNetworkReadyV1Beta1Condition, vmwarev1.ClusterNetworkProvisionFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterNetworkReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterNetworkNotReadyReason,
			Message: err.Error(),
		})
		return pkgerrors.Wrap(err, "Failed to provision network")
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
func (vp *nsxtVPCNetworkProvider) ConfigureVirtualMachine(_ context.Context, clusterCtx *vmware.ClusterContext, machine *vmwarev1.VSphereMachine, vm *vmoprvhub.VirtualMachine) error {
	vm.Spec.Network = &vmoprvhub.VirtualMachineNetworkSpec{}

	ipamModes, err := getIPAMModes(clusterCtx)
	if err != nil {
		return err
	}

	// Set the VM primary interface
	if createSubnetSet(clusterCtx) {
		if machine.Spec.Network.Interfaces.Primary.IsDefined() {
			return pkgerrors.New("primary interface can not be configured when createSubnetSet is true")
		}
		networkName := clusterCtx.VSphereCluster.Name
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmoprvhub.VirtualMachineNetworkInterfaceSpec{
			Name: PrimaryInterfaceName,
			Network: &vmoprvhub.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       NetworkGVKNSXTVPCSubnetSet.Kind,
					APIVersion: NetworkGVKNSXTVPCSubnetSet.GroupVersion().String(),
				},
				Name: networkName,
			},
			IPAMModes: ipamModes,
		})
	} else {
		if !machine.Spec.Network.Interfaces.Primary.IsDefined() {
			return pkgerrors.New("primary interface must be configured when createSubnetSet is false")
		}
		primary := machine.Spec.Network.Interfaces.Primary
		var mtu *int64
		if primary.MTU != 0 {
			mtu = ptr.To(int64(primary.MTU))
		}
		vmInterface := vmoprvhub.VirtualMachineNetworkInterfaceSpec{
			Name: PrimaryInterfaceName,
			Network: &vmoprvhub.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       primary.NetworkRef.Kind,
					APIVersion: primary.NetworkRef.APIVersion,
				},
				Name: primary.NetworkRef.Name,
			},
			MTU:       mtu,
			IPAMModes: ipamModes,
		}
		setRoutes(&vmInterface, primary.Routes)
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmInterface)
	}

	// Set the VM secondary interfaces
	setVMSecondaryInterfaces(machine, vm, ipamModes)

	// Set the VM VLAN sub-interfaces
	return setVLANs(machine, vm)
}

// getIPAMModes maps the ClusterIPFamily to the VM Operator's expected IPAMModes format.
// For dual-stack configurations, it always returns [IPv4, IPv6] regardless of which is primary.
// Note: This function is feature gated by IPv6DualStack.
func getIPAMModes(clusterCtx *vmware.ClusterContext) ([]corev1.IPFamily, error) {
	if !feature.Gates.Enabled(feature.IPv6DualStack) {
		return nil, nil
	}
	ipFamily, err := infrautilv1.DetermineClusterIPFamily(clusterCtx.Cluster)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to determine cluster IP family")
	}
	switch ipFamily {
	case infrautilv1.IPv4SingleStack:
		return []corev1.IPFamily{corev1.IPv4Protocol}, nil
	case infrautilv1.IPv6SingleStack:
		return []corev1.IPFamily{corev1.IPv6Protocol}, nil
	case infrautilv1.DualStackIPv4Primary, infrautilv1.DualStackIPv6Primary:
		return []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol}, nil
	default:
		return []corev1.IPFamily{corev1.IPv4Protocol}, nil
	}
}

func setRoutes(vmInterface *vmoprvhub.VirtualMachineNetworkInterfaceSpec, routes []vmwarev1.RouteSpec) {
	for _, route := range routes {
		vmInterface.Routes = append(vmInterface.Routes, vmoprvhub.VirtualMachineNetworkRouteSpec{
			To:  route.To,
			Via: route.Via,
		})
	}
}

func setVMSecondaryInterfaces(machine *vmwarev1.VSphereMachine, vm *vmoprvhub.VirtualMachine, ipamModes []corev1.IPFamily) {
	if len(machine.Spec.Network.Interfaces.Secondary) == 0 {
		return
	}
	for _, secondaryInterface := range machine.Spec.Network.Interfaces.Secondary {
		var mtu *int64
		if secondaryInterface.MTU != 0 {
			mtu = ptr.To(int64(secondaryInterface.MTU))
		}
		vmInterface := vmoprvhub.VirtualMachineNetworkInterfaceSpec{
			Name: secondaryInterface.Name,
			Network: &vmoprvhub.PartialObjectRef{
				TypeMeta: metav1.TypeMeta{
					Kind:       secondaryInterface.NetworkRef.Kind,
					APIVersion: secondaryInterface.NetworkRef.APIVersion,
				},
				Name: secondaryInterface.NetworkRef.Name,
			},
			MTU:       mtu,
			Gateway4:  "None",
			Gateway6:  "None",
			IPAMModes: ipamModes,
		}
		setRoutes(&vmInterface, secondaryInterface.Routes)
		vm.Spec.Network.Interfaces = append(vm.Spec.Network.Interfaces, vmInterface)
	}
}

func setVLANs(machine *vmwarev1.VSphereMachine, vm *vmoprvhub.VirtualMachine) error {
	if len(machine.Spec.Network.VLANs) == 0 {
		return nil
	}
	if !feature.Gates.Enabled(feature.VLANSubinterface) {
		return pkgerrors.New("invalid configuration: VLANs cannot be used as feature gate VLANSubinterface is not enabled")
	}
	if vm.Spec.Network == nil {
		vm.Spec.Network = &vmoprvhub.VirtualMachineNetworkSpec{}
	}
	for _, vlan := range machine.Spec.Network.VLANs {
		vm.Spec.Network.VLANs = append(vm.Spec.Network.VLANs, vmoprvhub.VirtualMachineNetworkVLANSpec{
			Name: vlan.Name,
			ID:   int64(ptr.Deref(vlan.ID, 0)),
			Link: vlan.Link,
		})
	}
	return nil
}
