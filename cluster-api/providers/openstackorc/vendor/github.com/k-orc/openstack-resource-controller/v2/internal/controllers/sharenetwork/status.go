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

package sharenetwork

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type sharenetworkStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.ShareNetworkApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.ShareNetworkStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.ShareNetwork, *osResourceT, *objectApplyT, *statusApplyT] = sharenetworkStatusWriter{}

func (sharenetworkStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.ShareNetwork(name, namespace)
}

func (sharenetworkStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.ShareNetwork, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		}
		return metav1.ConditionUnknown, nil
	}

	// Share networks become available immediately after creation
	// No async operations to wait for
	return metav1.ConditionTrue, nil
}

func (sharenetworkStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.ShareNetworkResourceStatus()

	if osResource.Name != "" {
		resourceStatus.WithName(osResource.Name)
	}

	if osResource.NeutronNetID != "" {
		resourceStatus.WithNeutronNetID(osResource.NeutronNetID)
	}

	if osResource.NeutronSubnetID != "" {
		resourceStatus.WithNeutronSubnetID(osResource.NeutronSubnetID)
	}

	if osResource.NetworkType != "" {
		resourceStatus.WithNetworkType(osResource.NetworkType)
	}

	// Always set CIDR field, even if empty, so it's always present in status
	resourceStatus.WithCIDR(osResource.CIDR)

	if osResource.ProjectID != "" {
		resourceStatus.WithProjectID(osResource.ProjectID)
	}

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}

	if osResource.SegmentationID != 0 {
		resourceStatus.WithSegmentationID(int32(osResource.SegmentationID))
	}

	if osResource.IPVersion != 0 {
		resourceStatus.WithIPVersion(int32(osResource.IPVersion))
	}

	if !osResource.CreatedAt.IsZero() {
		resourceStatus.WithCreatedAt(metav1.Time{Time: osResource.CreatedAt})
	}

	if !osResource.UpdatedAt.IsZero() {
		resourceStatus.WithUpdatedAt(metav1.Time{Time: osResource.UpdatedAt})
	}

	statusApply.WithResource(resourceStatus)
}
