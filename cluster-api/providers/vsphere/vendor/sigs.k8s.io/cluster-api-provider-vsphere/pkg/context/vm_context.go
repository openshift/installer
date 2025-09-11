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

// Package context defines context objects for controllers.
package context

import (
	"context"
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

// VMContext is a Go context used with a VSphereVM.
type VMContext struct {
	*ControllerManagerContext
	ClusterModuleInfo    *string
	VSphereVM            *infrav1.VSphereVM
	PatchHelper          *patch.Helper
	Session              *session.Session
	VSphereFailureDomain *infrav1.VSphereFailureDomain
}

// String returns VSphereVMGroupVersionKind VSphereVMNamespace/VSphereVMName.
func (c *VMContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereVM.GroupVersionKind(), c.VSphereVM.Namespace, c.VSphereVM.Name)
}

// Patch updates the object and its status on the API server.
func (c *VMContext) Patch(ctx context.Context) error {
	return c.PatchHelper.Patch(ctx, c.VSphereVM, patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.VSphereVMReadyV1Beta2Condition,
		infrav1.VSphereVMVCenterAvailableV1Beta2Condition,
		infrav1.VSphereVMVirtualMachineProvisionedV1Beta2Condition,
		infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
		infrav1.VSphereVMGuestSoftPowerOffSucceededV1Beta2Condition,
		infrav1.VSphereVMPCIDevicesDetachedV1Beta2Condition,
		clusterv1.PausedV1Beta2Condition,
	}})
}

// GetSession returns this context's session.
func (c *VMContext) GetSession() *session.Session {
	return c.Session
}
