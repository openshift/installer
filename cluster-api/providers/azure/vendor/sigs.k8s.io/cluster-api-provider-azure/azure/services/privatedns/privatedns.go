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

package privatedns

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/tags"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "privatedns"

// Scope defines the scope interface for a private dns service.
type Scope interface {
	azure.ClusterDescriber
	azure.Authorizer
	azure.AsyncStatusUpdater
	PrivateDNSSpec() (zoneSpec azure.ResourceSpecGetter, linksSpec, recordsSpec []azure.ResourceSpecGetter)
}

// Service provides operations on Azure resources.
type Service struct {
	Scope              Scope
	TagsGetter         async.TagsGetter
	zoneReconciler     async.Reconciler
	vnetLinkReconciler async.Reconciler
	recordReconciler   async.Reconciler
}

// New creates a new private dns service.
func New(scope Scope) (*Service, error) {
	zoneClient, err := newPrivateZonesClient(scope)
	if err != nil {
		return nil, err
	}
	vnetLinkClient, err := newVirtualNetworkLinksClient(scope)
	if err != nil {
		return nil, err
	}
	recordSetsClient, err := newRecordSetsClient(scope)
	if err != nil {
		return nil, err
	}
	tagsClient, err := tags.NewClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope:      scope,
		TagsGetter: tagsClient,
		zoneReconciler: async.New[armprivatedns.PrivateZonesClientCreateOrUpdateResponse,
			armprivatedns.PrivateZonesClientDeleteResponse](scope, zoneClient, zoneClient),
		vnetLinkReconciler: async.New[armprivatedns.VirtualNetworkLinksClientCreateOrUpdateResponse,
			armprivatedns.VirtualNetworkLinksClientDeleteResponse](scope, vnetLinkClient, vnetLinkClient),
		recordReconciler: async.New[armprivatedns.RecordSetsClientCreateOrUpdateResponse,
			armprivatedns.RecordSetsClientDeleteResponse](scope, recordSetsClient, recordSetsClient),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile creates or updates the private zone, links it to the vnet, and creates DNS records.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	zoneSpec, links, records := s.Scope.PrivateDNSSpec()
	if zoneSpec == nil {
		return nil
	}

	managed, err := s.reconcileZone(ctx, zoneSpec)
	if managed {
		s.Scope.UpdatePutStatus(infrav1.PrivateDNSZoneReadyCondition, serviceName, err)
	}
	if err != nil {
		return err
	}

	managed, err = s.reconcileLinks(ctx, links)
	if managed {
		s.Scope.UpdatePutStatus(infrav1.PrivateDNSLinkReadyCondition, serviceName, err)
	}
	if err != nil {
		return err
	}

	err = s.reconcileRecords(ctx, records)
	s.Scope.UpdatePutStatus(infrav1.PrivateDNSRecordReadyCondition, serviceName, err)
	return err
}

// Delete deletes the private zone and vnet links.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	zoneSpec, links, _ := s.Scope.PrivateDNSSpec()
	if zoneSpec == nil {
		return nil
	}

	managed, err := s.deleteLinks(ctx, links)
	if managed {
		s.Scope.UpdateDeleteStatus(infrav1.PrivateDNSLinkReadyCondition, serviceName, err)
	}
	if err != nil {
		return err
	}

	managed, err = s.deleteZone(ctx, zoneSpec)
	if managed {
		s.Scope.UpdateDeleteStatus(infrav1.PrivateDNSZoneReadyCondition, serviceName, err)
		s.Scope.UpdateDeleteStatus(infrav1.PrivateDNSRecordReadyCondition, serviceName, err)
	}

	return err
}

// isVnetLinkManaged returns true if the vnet link has an owned tag with the cluster name as value,
// meaning that the vnet link lifecycle is managed.
func (s *Service) isVnetLinkManaged(ctx context.Context, spec azure.ResourceSpecGetter) (bool, error) {
	scope := azure.VirtualNetworkLinkID(s.Scope.SubscriptionID(), spec.ResourceGroupName(), spec.OwnerResourceName(), spec.ResourceName())
	result, err := s.TagsGetter.GetAtScope(ctx, scope)
	if err != nil {
		return false, err
	}

	tagsMap := make(map[string]*string)
	if result.Properties != nil && result.Properties.Tags != nil {
		tagsMap = result.Properties.Tags
	}

	tags := converters.MapToTags(tagsMap)
	return tags.HasOwned(s.Scope.ClusterName()), nil
}

// IsManaged returns true if the private DNS has an owned tag with the cluster name as value,
// meaning that the DNS lifecycle is managed.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	zoneSpec, _, _ := s.Scope.PrivateDNSSpec()
	if zoneSpec == nil {
		return false, errors.Errorf("no private dns zone spec available")
	}

	scope := azure.PrivateDNSZoneID(s.Scope.SubscriptionID(), zoneSpec.ResourceGroupName(), zoneSpec.ResourceName())
	result, err := s.TagsGetter.GetAtScope(ctx, scope)
	if err != nil {
		return false, err
	}

	tagsMap := make(map[string]*string)
	if result.Properties != nil && result.Properties.Tags != nil {
		tagsMap = result.Properties.Tags
	}

	tags := converters.MapToTags(tagsMap)
	return tags.HasOwned(s.Scope.ClusterName()), nil
}
