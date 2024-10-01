/*
Copyright 2019 The Kubernetes Authors.

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

package services

import (
	"context"

	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
)

// VSphereMachineService is used for vsphere VM lifecycle and syncing with VSphereMachine types.
type VSphereMachineService interface {
	GetMachinesInCluster(ctx context.Context, namespace, clusterName string) ([]client.Object, error)
	FetchVSphereMachine(ctx context.Context, name types.NamespacedName) (capvcontext.MachineContext, error)
	FetchVSphereCluster(ctx context.Context, cluster *clusterv1.Cluster, machineContext capvcontext.MachineContext) (capvcontext.MachineContext, error)
	ReconcileDelete(ctx context.Context, machineCtx capvcontext.MachineContext) error
	SyncFailureReason(ctx context.Context, machineCtx capvcontext.MachineContext) (bool, error)
	ReconcileNormal(ctx context.Context, machineCtx capvcontext.MachineContext) (bool, error)
	GetHostInfo(ctx context.Context, machineCtx capvcontext.MachineContext) (string, error)
}

// VirtualMachineService is a service for creating/updating/deleting virtual
// machines on vSphere.
type VirtualMachineService interface {
	// ReconcileVM reconciles a VM with the intended state.
	ReconcileVM(ctx context.Context, vmCtx *capvcontext.VMContext) (infrav1.VirtualMachine, error)

	// DestroyVM powers off and removes a VM from the inventory.
	DestroyVM(ctx context.Context, vmCtx *capvcontext.VMContext) (reconcile.Result, infrav1.VirtualMachine, error)
}

// ControlPlaneEndpointService is a service for reconciling load balanced control plane endpoints.
type ControlPlaneEndpointService interface {
	// ReconcileControlPlaneEndpointService manages the lifecycle of a
	// control plane endpoint managed by a vmoperator VirtualMachineService
	ReconcileControlPlaneEndpointService(ctx context.Context, clusterCtx *vmware.ClusterContext, netProvider NetworkProvider) (*clusterv1.APIEndpoint, error)
}

// ResourcePolicyService is a service for reconciling a VirtualMachineSetResourcePolicy for a cluster.
type ResourcePolicyService interface {
	// ReconcileResourcePolicy ensures that a VirtualMachineSetResourcePolicy exists for the cluster
	// Returns the name of a policy if it exists, otherwise returns an error
	ReconcileResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error)
}

// NetworkProvider provision network resources and configures VM based on network type.
type NetworkProvider interface {
	// HasLoadBalancer indicates whether this provider has a load balancer for Services.
	HasLoadBalancer() bool

	// SupportsVMReadinessProbe indicates whether this provider support vm readiness probe.
	SupportsVMReadinessProbe() bool

	// ProvisionClusterNetwork creates network resource for a given cluster
	// This operation should be idempotent
	ProvisionClusterNetwork(ctx context.Context, clusterCtx *vmware.ClusterContext) error

	// GetClusterNetworkName returns the name of a valid cluster network if one exists
	// Returns an empty string if the operation is not supported
	GetClusterNetworkName(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error)

	// GetVMServiceAnnotations returns the annotations, if any, to place on a VM Service.
	GetVMServiceAnnotations(ctx context.Context, clusterCtx *vmware.ClusterContext) (map[string]string, error)

	// ConfigureVirtualMachine configures a VM for the particular network
	ConfigureVirtualMachine(ctx context.Context, clusterCtx *vmware.ClusterContext, vm *vmoprv1.VirtualMachine) error

	// VerifyNetworkStatus verifies the status of the network after vnet creation
	VerifyNetworkStatus(ctx context.Context, clusterCtx *vmware.ClusterContext, obj runtime.Object) error
}
