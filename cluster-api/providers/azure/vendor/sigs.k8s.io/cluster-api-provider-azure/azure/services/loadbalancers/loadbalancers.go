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

package loadbalancers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const (
	serviceName           = "loadbalancers"
	httpsProbe            = "HTTPSProbe"
	httpsProbeRequestPath = "/readyz"
	lbRuleHTTPS           = "LBRuleHTTPS"
	outboundNAT           = "OutboundNATAllProtocols"
)

// LBScope defines the scope interface for a load balancer service.
type LBScope interface {
	azure.ClusterScoper
	azure.AsyncStatusUpdater
	LBSpecs() []azure.ResourceSpecGetter
}

// Service provides operations on Azure resources.
type Service struct {
	Scope LBScope
	async.Reconciler
}

// New creates a new service.
func New(scope LBScope) (*Service, error) {
	client, err := newClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armnetwork.LoadBalancersClientCreateOrUpdateResponse,
			armnetwork.LoadBalancersClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a load balancer.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "loadbalancers.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	return azure.ReconcileAll(ctx, s.Reconciler, s.Scope, s.Scope.LBSpecs(), serviceName, infrav1.LoadBalancersReadyCondition)
}

// Delete deletes the public load balancer with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "loadbalancers.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	return azure.DeleteAll(ctx, s.Reconciler, s.Scope, s.Scope.LBSpecs(), serviceName, infrav1.LoadBalancersReadyCondition)
}

// IsManaged returns always returns true as CAPZ does not support BYO load balancers.
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
