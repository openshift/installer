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

package routetables

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "routetables"

// RouteTableScope defines the scope interface for route table service.
type RouteTableScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	RouteTableSpecs() []azure.ResourceSpecGetter
	IsVnetManaged() bool
}

// Service provides operations on azure resources.
type Service struct {
	Scope RouteTableScope
	async.Reconciler
}

// New creates a new service.
func New(scope RouteTableScope) (*Service, error) {
	client, err := newClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.RouteTablesClientCreateOrUpdateResponse,
			armnetwork.RouteTablesClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a set of route tables.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "routetables.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	if managed, err := s.IsManaged(ctx); err == nil && !managed {
		log.V(4).Info("Skipping route tables reconcile in custom vnet mode")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to check if route tables are managed")
	}

	return azure.ReconcileAll(ctx, s.Reconciler, s.Scope, s.Scope.RouteTableSpecs(), serviceName, infrav1.RouteTablesReadyCondition)
}

// Delete deletes route tables.
func (s *Service) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "routetables.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	// Only delete the route tables if their lifecycle is managed by this controller.
	// Route tables are managed if and only if the vnet is managed.
	if managed, err := s.IsManaged(ctx); err == nil && !managed {
		log.V(4).Info("Skipping route table deletion in custom vnet mode")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to check if route tables are managed")
	}

	return azure.DeleteAll(ctx, s.Reconciler, s.Scope, s.Scope.RouteTableSpecs(), serviceName, infrav1.RouteTablesReadyCondition)
}

// IsManaged returns true if the route tables' lifecycles are managed.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "routetables.Service.IsManaged")
	defer done()

	return s.Scope.IsVnetManaged(), nil
}
