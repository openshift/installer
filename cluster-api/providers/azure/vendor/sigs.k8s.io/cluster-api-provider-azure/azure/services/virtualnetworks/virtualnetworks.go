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

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	"github.com/Azure/azure-service-operator/v2/pkg/common/labels"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const serviceName = "virtualnetworks"

// VNetScope defines the scope interface for a virtual network service.
type VNetScope interface {
	aso.Scope
	Vnet() *infrav1.VnetSpec
	VNetSpec() azure.ASOResourceSpecGetter[*asonetworkv1.VirtualNetwork]
	UpdateSubnetCIDRs(string, []string)
}

// New creates a new service.
func New(scope VNetScope) *aso.Service[*asonetworkv1.VirtualNetwork, VNetScope] {
	svc := aso.NewService[*asonetworkv1.VirtualNetwork](serviceName, scope)
	svc.Specs = []azure.ASOResourceSpecGetter[*asonetworkv1.VirtualNetwork]{scope.VNetSpec()}
	svc.ConditionType = infrav1.VNetReadyCondition
	svc.PostCreateOrUpdateResourceHook = postCreateOrUpdateResourceHook
	return svc
}

func postCreateOrUpdateResourceHook(ctx context.Context, scope VNetScope, existingVnet *asonetworkv1.VirtualNetwork, err error) error {
	if err != nil {
		return err
	}

	vnet := scope.Vnet()
	vnet.ID = ptr.Deref(existingVnet.Status.Id, "")
	vnet.Tags = existingVnet.Status.Tags

	// Update the subnet CIDRs if they already exist.
	// This makes sure the subnet CIDRs are up to date and there are no validation errors when updating the VNet.
	// Subnets that are not part of this cluster spec are silently ignored.
	subnets := &asonetworkv1.VirtualNetworksSubnetList{}
	err = scope.GetClient().List(ctx, subnets,
		client.InNamespace(existingVnet.Namespace),
		client.MatchingLabels{labels.OwnerNameLabel: existingVnet.Name},
	)
	if err != nil {
		return errors.Wrap(err, "failed to list subnets")
	}
	for _, subnet := range subnets.Items {
		scope.UpdateSubnetCIDRs(subnet.AzureName(), converters.GetSubnetAddresses(subnet))
	}
	// Only update the vnet's CIDRBlocks when we also updated subnets' since the vnet is created before
	// subnets to prevent an updated vnet CIDR from invalidating subnet CIDRs that were defaulted and do not
	// exist yet.
	if len(subnets.Items) > 0 && existingVnet.Status.AddressSpace != nil {
		vnet.CIDRBlocks = existingVnet.Status.AddressSpace.AddressPrefixes
	}

	return nil
}
