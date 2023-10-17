/*
Copyright 2023 The Kubernetes Authors.

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

package groups

import (
	"context"

	asoannotations "github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceName is the name of this service.
const ServiceName = "group"

// Service provides operations on Azure resources.
type Service struct {
	Scope GroupScope
	aso.Reconciler
}

// GroupScope defines the scope interface for a group service.
type GroupScope interface {
	azure.AsyncStatusUpdater
	GroupSpec() azure.ASOResourceSpecGetter
	GetClient() client.Client
	ClusterName() string
}

// New creates a new service.
func New(scope GroupScope) *Service {
	return &Service{
		Scope:      scope,
		Reconciler: aso.New(scope.GetClient(), scope.ClusterName()),
	}
}

// Name returns the service name.
func (s *Service) Name() string {
	return ServiceName
}

// Reconcile idempotently creates or updates a resource group.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "groups.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	groupSpec := s.Scope.GroupSpec()
	if groupSpec == nil {
		return nil
	}

	_, err := s.CreateOrUpdateResource(ctx, groupSpec, ServiceName)
	s.Scope.UpdatePutStatus(infrav1.ResourceGroupReadyCondition, ServiceName, err)
	return err
}

// Delete deletes the resource group if it is managed by capz.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "groups.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	groupSpec := s.Scope.GroupSpec()
	if groupSpec == nil {
		return nil
	}

	err := s.DeleteResource(ctx, groupSpec, ServiceName)
	s.Scope.UpdateDeleteStatus(infrav1.ResourceGroupReadyCondition, ServiceName, err)
	return err
}

// IsManaged returns true if the ASO ResourceGroup was created by CAPZ,
// meaning that the resource group's lifecycle is managed.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return aso.IsManaged(ctx, s.Scope.GetClient(), s.Scope.GroupSpec(), s.Scope.ClusterName())
}

var _ azure.Pauser = (*Service)(nil)

// Pause implements azure.Pauser.
func (s *Service) Pause(ctx context.Context) error {
	groupSpec := s.Scope.GroupSpec()
	if groupSpec == nil {
		return nil
	}
	return aso.PauseResource(ctx, s.Scope.GetClient(), groupSpec, s.Scope.ClusterName(), ServiceName)
}

// ShouldDeleteIndividualResources returns false if the resource group is
// managed and reconciled by ASO, meaning that we can rely on a single resource
// group delete operation as opposed to deleting every individual resource.
func (s *Service) ShouldDeleteIndividualResources(ctx context.Context) bool {
	// Since this is a best effort attempt to speed up delete, we don't fail the delete if we can't get the RG status.
	// Instead, take the long way and delete all resources one by one.
	managed, err := s.IsManaged(ctx)
	if err != nil || !managed {
		return true
	}

	// For ASO, "managed" only tells us that we're allowed to delete the ASO
	// resource. We also need to check that deleting the ASO resource will really
	// delete the underlying resource group by checking the ASO reconcile-policy.
	spec := s.Scope.GroupSpec()
	group := spec.ResourceRef()
	err = s.Scope.GetClient().Get(ctx, client.ObjectKeyFromObject(group), group)
	return err != nil || group.GetAnnotations()[asoannotations.ReconcilePolicy] != string(asoannotations.ReconcilePolicyManage)
}
