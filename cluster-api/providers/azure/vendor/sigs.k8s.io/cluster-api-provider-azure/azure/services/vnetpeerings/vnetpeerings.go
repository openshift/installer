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

package vnetpeerings

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// ServiceName is the name of this service.
const ServiceName = "vnetpeerings"

// VnetPeeringScope defines the scope interface for a subnet service.
type VnetPeeringScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	VnetPeeringSpecs() []azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope VnetPeeringScope
	async.Reconciler
}

// New creates a new service.
func New(scope VnetPeeringScope) (*Service, error) {
	Client, err := NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.VirtualNetworkPeeringsClientCreateOrUpdateResponse,
			armnetwork.VirtualNetworkPeeringsClientDeleteResponse](scope, Client, Client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return ServiceName
}

// Reconcile idempotently creates or updates a peering.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vnetpeerings.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	return azure.ReconcileAll(ctx, s.Reconciler, s.Scope, s.Scope.VnetPeeringSpecs(), ServiceName, infrav1.VnetPeeringReadyCondition)
}

// Delete deletes the peering with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "vnetpeerings.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	return azure.DeleteAll(ctx, s.Reconciler, s.Scope, s.Scope.VnetPeeringSpecs(), ServiceName, infrav1.VnetPeeringReadyCondition)
}

// IsManaged returns always returns true as CAPZ does not support BYO VNet peering.
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
