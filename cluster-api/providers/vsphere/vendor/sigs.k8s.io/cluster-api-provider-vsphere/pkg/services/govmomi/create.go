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

package govmomi

import (
	"context"

	"github.com/pkg/errors"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"

	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/vcenter"
)

// createVM creates a new VM with the data in the VMContext passed. This method does not wait
// for the new VM to be created.
func createVM(ctx context.Context, vmCtx *capvcontext.VMContext, bootstrapData []byte, format bootstrapv1.Format) error {
	if !vmCtx.Session.IsVC() {
		return errors.Errorf("expected VCenter client got %v", vmCtx.Session.ServiceContent.About.ApiType)
	}
	return vcenter.Clone(ctx, vmCtx, bootstrapData, format)
}
