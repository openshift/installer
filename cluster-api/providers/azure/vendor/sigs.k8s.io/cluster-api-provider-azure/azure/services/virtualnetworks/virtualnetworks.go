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

package virtualnetworks

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/tags"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "virtualnetworks"

// VNetScope defines the scope interface for a virtual network service.
type VNetScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	Vnet() *infrav1.VnetSpec
	VNetSpec() azure.ResourceSpecGetter
	ClusterName() string
	IsVnetManaged() bool
	UpdateSubnetCIDRs(string, []string)
}

// Service provides operations on Azure resources.
type Service struct {
	Scope VNetScope
	async.Reconciler
	async.Getter
	async.TagsGetter
}

// New creates a new service.
func New(scope VNetScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	tagsClient, err := tags.NewClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope:      scope,
		Getter:     client,
		TagsGetter: tagsClient,
		Reconciler: async.New[armnetwork.VirtualNetworksClientCreateOrUpdateResponse,
			armnetwork.VirtualNetworksClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a virtual network.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	vnetSpec := s.Scope.VNetSpec()
	if vnetSpec == nil {
		return nil
	}

	result, err := s.CreateOrUpdateResource(ctx, vnetSpec, serviceName)
	if err == nil && result != nil {
		existingVnet, ok := result.(armnetwork.VirtualNetwork)
		if !ok {
			return errors.Errorf("%T is not an armnetwork.VirtualNetwork", result)
		}
		vnet := s.Scope.Vnet()
		vnet.ID = ptr.Deref(existingVnet.ID, "")
		vnet.Tags = converters.MapToTags(existingVnet.Tags)

		var prefixes []string
		if existingVnet.Properties != nil && existingVnet.Properties.AddressSpace != nil {
			for _, prefix := range existingVnet.Properties.AddressSpace.AddressPrefixes {
				if prefix != nil {
					prefixes = append(prefixes, *prefix)
				}
			}
		}
		vnet.CIDRBlocks = prefixes

		// Update the subnet CIDRs if they already exist.
		// This makes sure the subnet CIDRs are up to date and there are no validation errors when updating the VNet.
		// Subnets that are not part of this cluster spec are silently ignored.
		if existingVnet.Properties.Subnets != nil {
			for _, subnet := range existingVnet.Properties.Subnets {
				s.Scope.UpdateSubnetCIDRs(ptr.Deref(subnet.Name, ""), converters.GetSubnetAddresses(subnet))
			}
		}
	}

	if s.Scope.IsVnetManaged() {
		s.Scope.UpdatePutStatus(infrav1.VNetReadyCondition, serviceName, err)
	}

	return err
}

// Delete deletes the virtual network if it is managed by capz.
func (s *Service) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	vnetSpec := s.Scope.VNetSpec()
	if vnetSpec == nil {
		return nil
	}

	// Check that the vnet is not BYO.
	managed, err := s.IsManaged(ctx)
	if err != nil {
		if azure.ResourceNotFound(err) {
			// already deleted or doesn't exist, cleanup status and return.
			s.Scope.DeleteLongRunningOperationState(vnetSpec.ResourceName(), serviceName, infrav1.DeleteFuture)
			s.Scope.UpdateDeleteStatus(infrav1.VNetReadyCondition, serviceName, nil)
			return nil
		}
		return errors.Wrap(err, "could not get VNet management state")
	}
	if !managed {
		log.Info("Skipping VNet deletion in custom vnet mode")
		return nil
	}

	err = s.DeleteResource(ctx, vnetSpec, serviceName)
	s.Scope.UpdateDeleteStatus(infrav1.VNetReadyCondition, serviceName, err)
	return err
}

// IsManaged returns true if the virtual network has an owned tag with the cluster name as value,
// meaning that the vnet's lifecycle is managed.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualnetworks.Service.IsManaged")
	defer done()

	spec := s.Scope.VNetSpec()
	if spec == nil {
		return false, errors.New("cannot get vnet to check if it is managed: spec is nil")
	}

	scope := azure.VNetID(s.Scope.SubscriptionID(), spec.ResourceGroupName(), spec.ResourceName())
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
