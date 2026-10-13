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

package router

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

const (
	RouterStatusActive = "ACTIVE"
)

type objectApplyPT = *orcapplyconfigv1alpha1.RouterApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.RouterStatusApplyConfiguration

type routerStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = routerStatusWriter{}

func (routerStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.Router(name, namespace)
}

func (routerStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	if osResource.Status == RouterStatusActive {
		return metav1.ConditionTrue, nil
	}
	return metav1.ConditionFalse, nil
}

func (routerStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	status := orcapplyconfigv1alpha1.RouterResourceStatus().
		WithName(osResource.Name).
		WithProjectID(osResource.ProjectID).
		WithStatus(osResource.Status).
		WithTags(osResource.Tags...).
		WithAdminStateUp(osResource.AdminStateUp).
		WithAvailabilityZoneHints(osResource.AvailabilityZoneHints...)
	if osResource.Description != "" {
		status.WithDescription(osResource.Description)
	}
	if osResource.GatewayInfo.NetworkID != "" {
		status.WithExternalGateways(orcapplyconfigv1alpha1.ExternalGatewayStatus().
			WithNetworkID(osResource.GatewayInfo.NetworkID))
	}

	statusApply.WithResource(status)
}
