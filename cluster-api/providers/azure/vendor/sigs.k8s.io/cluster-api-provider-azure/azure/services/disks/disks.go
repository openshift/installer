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
	client, err := newClient(scope, scope.DefaultedAzureCallTimeout())
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

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	return azure.DeleteAll(ctx, s.Reconciler, s.Scope, s.Scope.DiskSpecs(), serviceName, infrav1.DisksReadyCondition)
}

// IsManaged returns always returns true as CAPZ does not support BYO disk.
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
