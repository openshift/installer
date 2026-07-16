/*
Copyright The ORC Authors.

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

package trunk

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type trunkStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.TrunkApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.TrunkStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.Trunk, *osResourceT, *objectApplyT, *statusApplyT] = trunkStatusWriter{}

func (trunkStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.Trunk(name, namespace)
}

func (trunkStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.Trunk, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}
	return metav1.ConditionTrue, nil
}

func (trunkStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.TrunkResourceStatus().
		WithPortID(osResource.PortID).
		WithProjectID(osResource.ProjectID).
		WithName(osResource.Name).
		WithAdminStateUp(osResource.AdminStateUp).
		WithRevisionNumber(int64(osResource.RevisionNumber)).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt)).
		WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt))

	if osResource.Status != "" {
		resourceStatus.WithStatus(osResource.Status)
	}

	if osResource.TenantID != "" {
		resourceStatus.WithTenantID(osResource.TenantID)
	}

	if len(osResource.Tags) > 0 {
		resourceStatus.WithTags(osResource.Tags...)
	}

	if len(osResource.Subports) > 0 {
		subports := make([]*orcapplyconfigv1alpha1.TrunkSubportStatusApplyConfiguration, 0, len(osResource.Subports))
		for i := range osResource.Subports {
			sp := osResource.Subports[i]
			subports = append(subports,
				orcapplyconfigv1alpha1.TrunkSubportStatus().
					WithPortID(sp.PortID).
					WithSegmentationID(int32(sp.SegmentationID)).
					WithSegmentationType(sp.SegmentationType),
			)
		}
		resourceStatus.WithSubports(subports...)
	}

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}

	statusApply.WithResource(resourceStatus)
}
