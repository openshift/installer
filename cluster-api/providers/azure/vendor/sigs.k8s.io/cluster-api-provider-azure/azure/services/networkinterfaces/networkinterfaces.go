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

package networkinterfaces

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "interfaces"

// NICScope defines the scope interface for a network interfaces service.
type NICScope interface {
	azure.ClusterDescriber
	azure.AsyncStatusUpdater
	NICSpecs() []azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope NICScope
	async.Reconciler
	resourceSKUCache *resourceskus.Cache
}

// New creates a new service.
func New(scope NICScope, skuCache *resourceskus.Cache) (*Service, error) {
	client, err := NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.InterfacesClientCreateOrUpdateResponse,
			armnetwork.InterfacesClientDeleteResponse](scope, client, client),
		resourceSKUCache: skuCache,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a network interface.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "networkinterfaces.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	specs := s.Scope.NICSpecs()
	if len(specs) == 0 {
		return nil
	}

	// We go through the list of NICSpecs to reconcile each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	var result error
	for _, nicSpec := range specs {
		if _, err := s.CreateOrUpdateResource(ctx, nicSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || result == nil {
				result = err
			}
		}
	}

	s.Scope.UpdatePutStatus(infrav1.NetworkInterfaceReadyCondition, serviceName, result)
	return result
}

// Delete deletes the network interface with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "networkinterfaces.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	specs := s.Scope.NICSpecs()
	if len(specs) == 0 {
		return nil
	}

	// We go through the list of NICSpecs to delete each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error deleting) -> operationNotDoneError (i.e. deleting in progress) -> no error (i.e. deleted)
	var result error
	for _, nicSpec := range specs {
		if err := s.DeleteResource(ctx, nicSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || result == nil {
				result = err
			}
		}
	}

	s.Scope.UpdateDeleteStatus(infrav1.NetworkInterfaceReadyCondition, serviceName, result)
	return result
}

// IsManaged returns always returns true as CAPZ does not support BYO network interfaces.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
