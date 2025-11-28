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

package port

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

const (
	PortStatusActive = "ACTIVE"
	PortStatusDown   = "DOWN"
)

type objectApplyPT = *orcapplyconfigv1alpha1.PortApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.PortStatusApplyConfiguration

type portStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = portStatusWriter{}

func (portStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.Port(name, namespace)
}

func (portStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	// Both active and down ports are Available
	if osResource.Status == PortStatusActive || osResource.Status == PortStatusDown {
		return metav1.ConditionTrue, nil
	}
	return metav1.ConditionFalse, nil
}

func (portStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	resourceStatus := orcapplyconfigv1alpha1.PortResourceStatus().
		WithName(osResource.Name).
		WithAdminStateUp(osResource.AdminStateUp).
		WithMACAddress(osResource.MACAddress).
		WithStatus(osResource.Status).
		WithProjectID(osResource.ProjectID).
		WithNetworkID(osResource.NetworkID).
		WithTags(osResource.Tags...).
		WithSecurityGroups(osResource.SecurityGroups...).
		WithPropagateUplinkStatus(osResource.PropagateUplinkStatus).
		WithVNICType(osResource.VNICType).
		WithPortSecurityEnabled(osResource.PortSecurityEnabled).
		WithRevisionNumber(int64(osResource.RevisionNumber)).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt)).
		WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt))

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}
	if osResource.DeviceID != "" {
		resourceStatus.WithDeviceID(osResource.DeviceID)
	}
	if osResource.DeviceOwner != "" {
		resourceStatus.WithDeviceOwner(osResource.DeviceOwner)
	}
	if len(osResource.AllowedAddressPairs) > 0 {
		allowedAddressPairs := make([]*orcapplyconfigv1alpha1.AllowedAddressPairStatusApplyConfiguration, len(osResource.AllowedAddressPairs))
		for i := range osResource.AllowedAddressPairs {
			allowedAddressPairs[i] = orcapplyconfigv1alpha1.AllowedAddressPairStatus().
				WithIP(osResource.AllowedAddressPairs[i].IPAddress).
				WithMAC(osResource.AllowedAddressPairs[i].MACAddress)
		}
		resourceStatus.WithAllowedAddressPairs(allowedAddressPairs...)
	}

	if len(osResource.FixedIPs) > 0 {
		fixedIPs := make([]*orcapplyconfigv1alpha1.FixedIPStatusApplyConfiguration, len(osResource.FixedIPs))
		for i := range osResource.FixedIPs {
			fixedIPs[i] = orcapplyconfigv1alpha1.FixedIPStatus().
				WithIP(osResource.FixedIPs[i].IPAddress).
				WithSubnetID(osResource.FixedIPs[i].SubnetID)
		}
		resourceStatus.WithFixedIPs(fixedIPs...)
	}

	statusApply.WithResource(resourceStatus)
}
