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
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// VMModifier allows a function to be passed to VM creation to modify its spec
// The hook is loosely typed so as to allow for different VirtualMachine backends.
type VMModifier func(runtime.Object) (runtime.Object, error)

// SupervisorMachineContext is a Go context used with a VSphereMachine.
type SupervisorMachineContext struct {
	*context.BaseMachineContext
	VSphereCluster *vmwarev1.VSphereCluster
	VSphereMachine *vmwarev1.VSphereMachine
	VMModifiers    []VMModifier
}

// String returns VSphereMachineGroupVersionKind VSphereMachineNamespace/VSphereMachineName.
func (c *SupervisorMachineContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereMachine.GroupVersionKind(), c.VSphereMachine.Namespace, c.VSphereMachine.Name)
}

// Patch updates the object and its status on the API server.
func (c *SupervisorMachineContext) Patch() error {
	return c.PatchHelper.Patch(c, c.VSphereMachine)
}

func (c *SupervisorMachineContext) GetVSphereMachine() context.VSphereMachine {
	return c.VSphereMachine
}

func (c *SupervisorMachineContext) GetObjectMeta() v1.ObjectMeta {
	return c.VSphereMachine.ObjectMeta
}

func (c *SupervisorMachineContext) GetClusterContext() *ClusterContext {
	return &ClusterContext{
		ControllerContext: c.ControllerContext,
		Cluster:           c.Cluster,
		VSphereCluster:    c.VSphereCluster,
		Logger:            c.GetLogger(),
	}
}

func (c *SupervisorMachineContext) SetBaseMachineContext(base *context.BaseMachineContext) {
	c.BaseMachineContext = base
}
