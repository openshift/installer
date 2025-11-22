/*
Copyright 2024 The ORC Authors.

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

package subnet

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type objectApplyPT = *orcapplyconfigv1alpha1.SubnetApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.SubnetStatusApplyConfiguration

type subnetStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = subnetStatusWriter{}

func (subnetStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.Subnet(name, namespace)
}

func (subnetStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	// Subnet is available as soon as it exists
	return metav1.ConditionTrue, nil
}

func (subnetStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	status := orcapplyconfigv1alpha1.SubnetResourceStatus().
		WithName(osResource.Name).
		WithIPVersion(int32(osResource.IPVersion)).
		WithCIDR(osResource.CIDR).
		WithGatewayIP(osResource.GatewayIP).
		WithDNSPublishFixedIP(osResource.DNSPublishFixedIP).
		WithEnableDHCP(osResource.EnableDHCP).
		WithNetworkID(osResource.NetworkID).
		WithProjectID(osResource.ProjectID).
		WithRevisionNumber(int64(osResource.RevisionNumber)).
		WithTags(osResource.Tags...).
		WithDNSNameservers(osResource.DNSNameservers...)

	if osResource.Description != "" {
		status.WithDescription(osResource.Description)
	}
	if osResource.IPv6AddressMode != "" {
		status.WithIPv6AddressMode(osResource.IPv6AddressMode)
	}
	if osResource.IPv6RAMode != "" {
		status.WithIPv6RAMode(osResource.IPv6RAMode)
	}

	for i := range osResource.AllocationPools {
		status.WithAllocationPools(orcapplyconfigv1alpha1.AllocationPoolStatus().
			WithStart(osResource.AllocationPools[i].Start).
			WithEnd(osResource.AllocationPools[i].End))
	}

	for i := range osResource.HostRoutes {
		status.WithHostRoutes(orcapplyconfigv1alpha1.HostRouteStatus().
			WithDestination(osResource.HostRoutes[i].DestinationCIDR).
			WithNextHop(osResource.HostRoutes[i].NextHop))
	}

	statusApply.WithResource(status)
}
