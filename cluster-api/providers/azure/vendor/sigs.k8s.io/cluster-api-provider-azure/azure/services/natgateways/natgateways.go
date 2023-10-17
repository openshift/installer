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

package natgateways

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "natgateways"

// NatGatewayScope defines the scope interface for NAT gateway service.
type NatGatewayScope interface {
	azure.ClusterScoper
	azure.AsyncStatusUpdater
	SetNatGatewayIDInSubnets(natGatewayName string, natGatewayID string)
	NatGatewaySpecs() []azure.ResourceSpecGetter
}

// Service provides operations on azure resources.
type Service struct {
	Scope NatGatewayScope
	async.Reconciler
}

// New creates a new service.
func New(scope NatGatewayScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.NatGatewaysClientCreateOrUpdateResponse,
			armnetwork.NatGatewaysClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a NAT gateway.
// Only when the NAT gateway 'Name' property is defined we create the NAT gateway: it's opt-in.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "natgateways.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	if managed, err := s.IsManaged(ctx); err == nil && !managed {
		log.V(4).Info("Skipping nat gateways reconcile in custom vnet mode")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to check if NAT gateways are managed")
	}

	// We go through the list of NatGatewaySpecs to reconcile each one, independently of the resultingErr of the previous one.
	specs := s.Scope.NatGatewaySpecs()
	if len(specs) == 0 {
		return nil
	}

	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (ie. error creating) -> operationNotDoneError (ie. creating in progress) -> no error (ie. created)
	var resultingErr error
	for _, natGatewaySpec := range specs {
		result, err := s.CreateOrUpdateResource(ctx, natGatewaySpec, serviceName)
		if err != nil {
			if !azure.IsOperationNotDoneError(err) || resultingErr == nil {
				resultingErr = err
			}
		}
		if err == nil {
			natGateway, ok := result.(armnetwork.NatGateway)
			if !ok {
				// Return out of loop since this would be an unexpected fatal error
				resultingErr = errors.Errorf("created resource %T is not an armnetwork.NatGateway", result)
				break
			}

			// TODO: ideally we wouldn't need to set the subnet spec based on the result of the create operation
			s.Scope.SetNatGatewayIDInSubnets(natGatewaySpec.ResourceName(), *natGateway.ID)
		}
	}

	s.Scope.UpdatePutStatus(infrav1.NATGatewaysReadyCondition, serviceName, resultingErr)
	return resultingErr
}

// Delete deletes the NAT gateway with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "natgateways.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	if managed, err := s.IsManaged(ctx); err == nil && !managed {
		log.V(4).Info("Skipping nat gateway deletion in custom vnet mode")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to check if NAT gateways are managed")
	}

	specs := s.Scope.NatGatewaySpecs()
	if len(specs) == 0 {
		return nil
	}

	// We go through the list of NatGatewaySpecs to delete each one, independently of the resultingErr of the previous one.
	// If multiple errors occur, we return the most pressing one.
	//  Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (ie. error creating) -> operationNotDoneError (ie. creating in progress) -> no error (ie. created)
	var resultingErr error
	for _, natGatewaySpec := range specs {
		if err := s.DeleteResource(ctx, natGatewaySpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || resultingErr == nil {
				resultingErr = err
			}
		}
	}
	s.Scope.UpdateDeleteStatus(infrav1.NATGatewaysReadyCondition, serviceName, resultingErr)
	return resultingErr
}

// IsManaged returns true if the NAT gateways' lifecycles are managed.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "natgateways.Service.IsManaged")
	defer done()

	return s.Scope.IsVnetManaged(), nil
}
