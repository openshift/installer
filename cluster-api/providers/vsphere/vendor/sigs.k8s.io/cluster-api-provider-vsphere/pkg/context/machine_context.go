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

package context

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// BaseMachineContext contains information about a CAPI Machine for VSphereMachine reconciliation.
type BaseMachineContext struct {
	ControllerManagerContext *ControllerManagerContext
	Cluster                  *clusterv1.Cluster
	Machine                  *clusterv1.Machine
	PatchHelper              *patch.Helper
}

// GetCluster returns the cluster for the BaseMachineContext.
func (c *BaseMachineContext) GetCluster() *clusterv1.Cluster {
	return c.Cluster
}

// GetMachine returns the Machine for the BaseMachineContext.
func (c *BaseMachineContext) GetMachine() *clusterv1.Machine {
	return c.Machine
}

// VIMMachineContext is a Go context used with a VSphereMachine.
type VIMMachineContext struct {
	*BaseMachineContext
	VSphereCluster *infrav1.VSphereCluster
	VSphereMachine *infrav1.VSphereMachine
}

// String returns VSphereMachineGroupVersionKind VSphereMachineNamespace/VSphereMachineName.
func (c *VIMMachineContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereMachine.GroupVersionKind(), c.VSphereMachine.Namespace, c.VSphereMachine.Name)
}

// Patch updates the object and its status on the API server.
func (c *VIMMachineContext) Patch(ctx context.Context) error {
	return c.PatchHelper.Patch(ctx, c.VSphereMachine)
}

// GetVSphereMachine sets the VSphereMachine for the VIMMachineContext.
func (c *VIMMachineContext) GetVSphereMachine() VSphereMachine {
	return c.VSphereMachine
}

// GetObjectMeta returns the ObjectMeta for the VSphereMachine in the VIMMachineContext.
func (c *VIMMachineContext) GetObjectMeta() metav1.ObjectMeta {
	return c.VSphereMachine.ObjectMeta
}

// SetBaseMachineContext sets the BaseMachineContext for the VIMMachineContext.
func (c *VIMMachineContext) SetBaseMachineContext(base *BaseMachineContext) {
	c.BaseMachineContext = base
}
