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
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

// VMContext is a Go context used with a VSphereVM.
type VMContext struct {
	*ControllerContext
	ClusterModuleInfo    *string
	VSphereVM            *infrav1.VSphereVM
	PatchHelper          *patch.Helper
	Logger               logr.Logger
	Session              *session.Session
	VSphereFailureDomain *infrav1.VSphereFailureDomain
}

// String returns VSphereVMGroupVersionKind VSphereVMNamespace/VSphereVMName.
func (c *VMContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereVM.GroupVersionKind(), c.VSphereVM.Namespace, c.VSphereVM.Name)
}

// Patch updates the object and its status on the API server.
func (c *VMContext) Patch() error {
	return c.PatchHelper.Patch(c, c.VSphereVM)
}

// GetLogger returns this context's logger.
func (c *VMContext) GetLogger() logr.Logger {
	return c.Logger
}

// GetSession returns this context's session.
func (c *VMContext) GetSession() *session.Session {
	return c.Session
}
