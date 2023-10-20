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

package disks

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "disks"

// DiskScope defines the scope interface for a disk service.
type DiskScope interface {
	azure.ClusterDescriber
	azure.AsyncStatusUpdater
	DiskSpecs() []azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope DiskScope
	async.Reconciler
}

// New creates a disks service.
func New(scope DiskScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armcompute.DisksClientCreateOrUpdateResponse,
			armcompute.DisksClientDeleteResponse](scope, nil, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile on disk is currently no-op. OS disks should only be deleted and will create with the VM automatically.
func (s *Service) Reconcile(ctx context.Context) error {
	_, _, done := tele.StartSpanWithLogger(ctx, "disks.Service.Reconcile")
	defer done()

	// DisksReadyCondition is set in the VM service.
	return nil
}

// Delete deletes the disk associated with a VM.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "disks.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	specs := s.Scope.DiskSpecs()
	if len(specs) == 0 {
		return nil
	}

	// We go through the list of DiskSpecs to delete each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	var result error
	for _, diskSpec := range specs {
		if err := s.DeleteResource(ctx, diskSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || result == nil {
				result = err
			}
		}
	}
	s.Scope.UpdateDeleteStatus(infrav1.DisksReadyCondition, serviceName, result)
	return result
}

// IsManaged returns always returns true as CAPZ does not support BYO disk.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
