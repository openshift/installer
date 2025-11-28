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

package floatingip

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

const (
	FloatingIPStatusActive = "ACTIVE"
	FloatingIPStatusDown   = "DOWN"
)

type objectApplyPT = *orcapplyconfigv1alpha1.FloatingIPApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.FloatingIPStatusApplyConfiguration

type floatingipStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = floatingipStatusWriter{}

func (floatingipStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.FloatingIP(name, namespace)
}

func (floatingipStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}
	// Both active and down ips are Available
	if osResource.Status == FloatingIPStatusActive || osResource.Status == FloatingIPStatusDown {
		return metav1.ConditionTrue, nil
	}

	return metav1.ConditionFalse, nil
}

func (floatingipStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	status := orcapplyconfigv1alpha1.FloatingIPResourceStatus().
		WithFloatingNetworkID(osResource.FloatingNetworkID).
		WithPortID(osResource.PortID).
		WithTenantID(osResource.TenantID).
		WithProjectID(osResource.ProjectID).
		WithStatus(osResource.Status).
		WithRouterID(osResource.RouterID).
		WithTags(osResource.Tags...).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt)).
		WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt)).
		WithFixedIP(osResource.FixedIP).
		WithFloatingIP(osResource.FloatingIP).
		WithRevisionNumber(int64(osResource.RevisionNumber))

	if osResource.Description != "" {
		status.WithDescription(osResource.Description)
	}

	statusApply.WithResource(status)
}
