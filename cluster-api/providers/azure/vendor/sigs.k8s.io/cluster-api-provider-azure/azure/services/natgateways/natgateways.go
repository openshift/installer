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

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
)

const serviceName = "natgateways"

// NatGatewayScope defines the scope interface for NAT gateway service.
type NatGatewayScope interface {
	azure.ClusterScoper
	aso.Scope
	SetNatGatewayIDInSubnets(natGatewayName string, natGatewayID string)
	NatGatewaySpecs() []azure.ASOResourceSpecGetter[*asonetworkv1.NatGateway]
}

// Service provides operations on azure resources.
type Service struct {
	Scope NatGatewayScope
	*aso.Service[*asonetworkv1.NatGateway, NatGatewayScope]
}

// New creates a new service.
func New(scope NatGatewayScope) *Service {
	svc := aso.NewService[*asonetworkv1.NatGateway, NatGatewayScope](serviceName, scope)
	svc.Specs = scope.NatGatewaySpecs()
	svc.ConditionType = infrav1.NATGatewaysReadyCondition
	svc.PostCreateOrUpdateResourceHook = postCreateOrUpdateResourceHook
	return &Service{
		Scope:   scope,
		Service: svc,
	}
}

func postCreateOrUpdateResourceHook(_ context.Context, scope NatGatewayScope, result *asonetworkv1.NatGateway, err error) error {
	if err != nil {
		return err
	}
	// TODO: ideally we wouldn't need to set the subnet spec based on the result of the create operation
	// result only gets populated when the resource is created or if it already exists
	if result != nil && result.Status.Id != nil {
		scope.SetNatGatewayIDInSubnets(result.Name, *result.Status.Id)
	}
	return nil
}
