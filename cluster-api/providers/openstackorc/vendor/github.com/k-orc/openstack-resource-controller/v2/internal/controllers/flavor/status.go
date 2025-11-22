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

package flavor

import (
	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type flavorStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.FlavorApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.FlavorStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.Flavor, *flavors.Flavor, *objectApplyT, *statusApplyT] = flavorStatusWriter{}

func (flavorStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.Flavor(name, namespace)
}

func (flavorStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.Flavor, osResource *flavors.Flavor) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	// Flavor is available as soon as it exists
	return metav1.ConditionTrue, nil
}

func (flavorStatusWriter) ApplyResourceStatus(_ logr.Logger, osResource *flavors.Flavor, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.FlavorResourceStatus().
		WithName(osResource.Name).
		WithRAM(int32(osResource.RAM)).
		WithDisk(int32(osResource.Disk)).
		WithVcpus(int32(osResource.VCPUs)).
		WithIsPublic(osResource.IsPublic)

	if osResource.Swap > 0 {
		resourceStatus.WithSwap(int32(osResource.Swap))
	}
	if osResource.Ephemeral > 0 {
		resourceStatus.WithEphemeral(int32(osResource.Ephemeral))
	}
	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}
	statusApply.WithResource(resourceStatus)
}
