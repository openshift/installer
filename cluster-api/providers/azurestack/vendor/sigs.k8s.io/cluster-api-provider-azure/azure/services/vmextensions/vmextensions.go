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

package vmextensions

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "vmextensions"

// VMExtensionScope defines the scope interface for a vm extension service.
type VMExtensionScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	VMExtensionSpecs() []azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope VMExtensionScope
	async.Reconciler
}

// New creates a new vm extension service.
func New(scope VMExtensionScope) (*Service, error) {
	client, err := newClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armcompute.VirtualMachineExtensionsClientCreateOrUpdateResponse,
			armcompute.VirtualMachineExtensionsClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a VM extension.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vmextensions.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	specs := s.Scope.VMExtensionSpecs()
	if len(specs) == 0 {
		return nil
	}

	// We go through the list of ExtensionSpecs to reconcile each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	var resultErr error
	for _, extensionSpec := range specs {
		_, err := s.CreateOrUpdateResource(ctx, extensionSpec, serviceName)
		if err != nil {
			if !azure.IsOperationNotDoneError(err) || resultErr == nil {
				resultErr = err
			}
		}
	}

	if azure.IsOperationNotDoneError(resultErr) {
		resultErr = errors.Wrapf(resultErr, "extension is still in provisioning state. This likely means that bootstrapping has not yet completed on the VM")
	} else if resultErr != nil {
		resultErr = errors.Wrapf(resultErr, "extension state failed. This likely means the Kubernetes node bootstrapping process failed or timed out. Check VM boot diagnostics logs to learn more")
	}

	s.Scope.UpdatePutStatus(infrav1.BootstrapSucceededCondition, serviceName, resultErr)
	return resultErr
}

// Delete is a no-op. VM Extensions will be deleted as part of VM deletion.
func (s *Service) Delete(_ context.Context) error {
	return nil
}

// IsManaged returns always returns true as CAPZ does not support BYO VM extension.
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
