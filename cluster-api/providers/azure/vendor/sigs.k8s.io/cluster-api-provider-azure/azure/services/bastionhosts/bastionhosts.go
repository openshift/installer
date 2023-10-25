/*
Copyright 2020 The Kubernetes Authors.

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

package bastionhosts

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "bastionhosts"

// BastionScope defines the scope interface for a bastion host service.
type BastionScope interface {
	azure.ClusterScoper
	azure.AsyncStatusUpdater
	AzureBastionSpec() azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope BastionScope
	async.Reconciler
}

// New creates a new service.
func New(scope BastionScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.BastionHostsClientCreateOrUpdateResponse,
			armnetwork.BastionHostsClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a bastion host.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "bastionhosts.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	var resultingErr error
	if bastionSpec := s.Scope.AzureBastionSpec(); bastionSpec != nil {
		_, resultingErr = s.CreateOrUpdateResource(ctx, bastionSpec, serviceName)
	} else {
		return nil
	}

	s.Scope.UpdatePutStatus(infrav1.BastionHostReadyCondition, serviceName, resultingErr)
	return resultingErr
}

// Delete deletes the bastion host with the provided scope.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "bastionhosts.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	var resultingErr error
	if bastionSpec := s.Scope.AzureBastionSpec(); bastionSpec != nil {
		resultingErr = s.DeleteResource(ctx, bastionSpec, serviceName)
	} else {
		return nil
	}

	s.Scope.UpdateDeleteStatus(infrav1.BastionHostReadyCondition, serviceName, resultingErr)
	return resultingErr
}

// IsManaged returns always returns true as CAPZ does not support BYO bastion.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
