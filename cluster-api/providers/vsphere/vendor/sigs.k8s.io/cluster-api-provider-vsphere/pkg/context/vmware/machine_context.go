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

package vmware

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// VMModifier allows a function to be passed to VM creation to modify its spec
// The hook is loosely typed so as to allow for different VirtualMachine backends.
type VMModifier func(runtime.Object) (runtime.Object, error)

// SupervisorMachineContext is a Go capvcontext used with a VSphereMachine.
type SupervisorMachineContext struct {
	*capvcontext.BaseMachineContext
	VSphereCluster *vmwarev1.VSphereCluster
	VSphereMachine *vmwarev1.VSphereMachine
	VMModifiers    []VMModifier
}

// String returns VSphereMachineGroupVersionKind VSphereMachineNamespace/VSphereMachineName.
func (c *SupervisorMachineContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereMachine.GroupVersionKind(), c.VSphereMachine.Namespace, c.VSphereMachine.Name)
}

// Patch updates the object and its status on the API server.
func (c *SupervisorMachineContext) Patch(ctx context.Context) error {
	return c.PatchHelper.Patch(ctx, c.VSphereMachine, patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.VSphereMachineReadyV1Beta2Condition,
		infrav1.VSphereMachineVirtualMachineProvisionedV1Beta2Condition,
		clusterv1.PausedV1Beta2Condition,
	}})
}

// GetVSphereMachine returns the VSphereMachine from the SupervisorMachineContext.
func (c *SupervisorMachineContext) GetVSphereMachine() capvcontext.VSphereMachine {
	return c.VSphereMachine
}

// GetReady return when the VSphereMachine is ready.
func (c *SupervisorMachineContext) GetReady() bool {
	return c.VSphereMachine.Status.Ready
}

// GetObjectMeta returns the metadata for the VSphereMachine from the SupervisorMachineContext.
func (c *SupervisorMachineContext) GetObjectMeta() metav1.ObjectMeta {
	return c.VSphereMachine.ObjectMeta
}

// GetClusterContext returns the Cluster and VSphereCluster from the SupervisorMachineContext.
func (c *SupervisorMachineContext) GetClusterContext() *ClusterContext {
	return &ClusterContext{
		Cluster:        c.Cluster,
		VSphereCluster: c.VSphereCluster,
	}
}

// SetBaseMachineContext sets the BaseMachineContext for the SupervisorMachineContext.
func (c *SupervisorMachineContext) SetBaseMachineContext(base *capvcontext.BaseMachineContext) {
	c.BaseMachineContext = base
}
