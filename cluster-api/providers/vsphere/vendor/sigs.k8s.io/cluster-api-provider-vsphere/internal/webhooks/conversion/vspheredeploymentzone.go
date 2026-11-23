/*
Copyright 2026 The Kubernetes Authors.

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

package conversion

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
)

// VSphereDeploymentZone is a HubSpokeConverter for the VSphereDeploymentZone API type.
var VSphereDeploymentZone = conversion.NewHubSpokeConverter(&infrav1.VSphereDeploymentZone{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereDeploymentZone{}, ConvertVSphereDeploymentZoneHubToV1Beta1, ConvertVSphereDeploymentZoneV1Beta1ToHub),
)

// ConvertVSphereDeploymentZoneV1Beta1ToHub converts a v1beta1 VSphereDeploymentZone to a hub VSphereDeploymentZone.
func ConvertVSphereDeploymentZoneV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereDeploymentZone, dst *infrav1.VSphereDeploymentZone) error {
	return infrav1beta1.Convert_v1beta1_VSphereDeploymentZone_To_v1beta2_VSphereDeploymentZone(src, dst, nil)
}

// ConvertVSphereDeploymentZoneHubToV1Beta1 converts a hub VSphereDeploymentZone to a v1beta1 VSphereDeploymentZone.
func ConvertVSphereDeploymentZoneHubToV1Beta1(_ context.Context, src *infrav1.VSphereDeploymentZone, dst *infrav1beta1.VSphereDeploymentZone) error {
	return infrav1beta1.Convert_v1beta2_VSphereDeploymentZone_To_v1beta1_VSphereDeploymentZone(src, dst, nil)
}
