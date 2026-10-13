/*
Copyright 2022 The Kubernetes Authors.

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

package resourcehealth

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "resourcehealth"

// ResourceHealthScope defines the scope interface for a resourcehealth service.
type ResourceHealthScope interface {
	azure.Authorizer
	AvailabilityStatusResourceURI() string
	AvailabilityStatusResource() conditions.Setter
}

// AvailabilityStatusFilterer transforms the condition derived from the
// availability status to allow the condition to be overridden in specific
// circumstances.
type AvailabilityStatusFilterer interface {
	AvailabilityStatusFilter(cond *clusterv1.Condition) *clusterv1.Condition
}

// Service provides operations on Azure resources.
type Service struct {
	Scope ResourceHealthScope
	client
}

// New creates a new service.
func New(scope ResourceHealthScope) (*Service, error) {
	cli, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope:  scope,
		client: cli,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile ensures the resource's availability status is reflected in its own status.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "resourcehealth.Service.Reconcile")
	defer done()

	if !feature.Gates.Enabled(feature.AKSResourceHealth) {
		conditions.Delete(s.Scope.AvailabilityStatusResource(), infrav1.AzureResourceAvailableCondition)
		return nil
	}

	resource := s.Scope.AvailabilityStatusResourceURI()
	availStatus, err := s.GetByResource(ctx, resource)
	if err != nil {
		return errors.Wrapf(err, "failed to get availability status for resource %s", resource)
	}
	log.V(2).Info("got availability status for resource", "resource", resource, "status", availStatus)

	cond := converters.SDKAvailabilityStatusToCondition(availStatus)
	if filterer, ok := s.Scope.(AvailabilityStatusFilterer); ok {
		cond = filterer.AvailabilityStatusFilter(cond)
	}

	conditions.Set(s.Scope.AvailabilityStatusResource(), cond)

	if cond.Status == corev1.ConditionFalse {
		return errors.Errorf("resource is not available: %s", cond.Message)
	}

	return nil
}

// Delete is a no-op.
func (s *Service) Delete(ctx context.Context) error {
	_, _, done := tele.StartSpanWithLogger(ctx, "resourcehealth.Service.Delete")
	defer done()

	return nil
}

// IsManaged always returns true.
func (s *Service) IsManaged(_ context.Context) (bool, error) {
	return true, nil
}
