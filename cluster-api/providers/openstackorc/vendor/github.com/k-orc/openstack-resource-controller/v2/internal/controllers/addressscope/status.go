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

package addressscope

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

type addressscopeStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.AddressScopeApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.AddressScopeStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.AddressScope, *osResourceT, *objectApplyT, *statusApplyT] = addressscopeStatusWriter{}

func (addressscopeStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.AddressScope(name, namespace)
}

func (addressscopeStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.AddressScope, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}
	return metav1.ConditionTrue, nil
}

func (addressscopeStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.AddressScopeResourceStatus().
		WithProjectID(osResource.ProjectID).
		WithName(osResource.Name).
		WithShared(osResource.Shared).
		WithIPVersion(int32(osResource.IPVersion))

	statusApply.WithResource(resourceStatus)
}
