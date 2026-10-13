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

package network

import (
	"strconv"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

const (
	NetworkStatusActive = "ACTIVE"
)

type networkStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.NetworkApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.NetworkStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.Network, *osclients.NetworkExt, *objectApplyT, *statusApplyT] = networkStatusWriter{}

func (networkStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.Network(name, namespace)
}

func (networkStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.Network, osResource *osclients.NetworkExt) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	if osResource.Status == NetworkStatusActive {
		return metav1.ConditionTrue, nil
	}
	return metav1.ConditionFalse, nil
}

func (networkStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osclients.NetworkExt, statusApply *orcapplyconfigv1alpha1.NetworkStatusApplyConfiguration) {
	networkResourceStatus := orcapplyconfigv1alpha1.NetworkResourceStatus().
		WithName(osResource.Name).
		WithAdminStateUp(osResource.AdminStateUp).
		WithAvailabilityZoneHints(osResource.AvailabilityZoneHints...).
		WithStatus(osResource.Status).
		WithProjectID(osResource.ProjectID).
		WithTags(osResource.Tags...).
		WithRevisionNumber(int64(osResource.RevisionNumber)).
		WithExternal(osResource.External).
		WithSubnets(osResource.Subnets...).
		WithMTU(int32(osResource.MTU)).
		WithPortSecurityEnabled(osResource.PortSecurityEnabled).
		WithShared(osResource.Shared).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt)).
		WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt))

	if osResource.Description != "" {
		networkResourceStatus.WithDescription(osResource.Description)
	}
	if osResource.DNSDomain != "" {
		networkResourceStatus.WithDNSDomain(osResource.DNSDomain)
	}
	if osResource.NetworkType != "" {
		providerProperties := orcapplyconfigv1alpha1.ProviderPropertiesStatus().
			WithNetworkType(osResource.NetworkType).
			WithPhysicalNetwork(osResource.PhysicalNetwork)

		if osResource.SegmentationID != "" {
			segmentationID, err := strconv.ParseInt(osResource.SegmentationID, 10, 32)
			if err != nil {
				log.V(logging.Info).Error(err, "Invalid segmentation ID", "segmentationID", osResource.SegmentationID)
			} else {
				providerProperties.WithSegmentationID(int32(segmentationID))
			}
		}
		networkResourceStatus.WithProvider(providerProperties)
	}

	statusApply.WithResource(networkResourceStatus)
}
