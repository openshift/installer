/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha3

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereDeploymentZone to the Hub version (v1beta1).
func (src *VSphereDeploymentZone) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereDeploymentZone)
	return Convert_v1alpha3_VSphereDeploymentZone_To_v1beta1_VSphereDeploymentZone(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereDeploymentZone.
func (dst *VSphereDeploymentZone) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.VSphereDeploymentZone)
	return Convert_v1beta1_VSphereDeploymentZone_To_v1alpha3_VSphereDeploymentZone(src, dst, nil)
}

// ConvertTo converts this VSphereDeploymentZoneList to the Hub version (v1beta1).
func (src *VSphereDeploymentZoneList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereDeploymentZoneList)
	return Convert_v1alpha3_VSphereDeploymentZoneList_To_v1beta1_VSphereDeploymentZoneList(src, dst, nil)
}

// ConvertFrom converts this VSphereDeploymentZoneList to the Hub version (v1beta1).
func (dst *VSphereDeploymentZoneList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.VSphereDeploymentZoneList)
	return Convert_v1beta1_VSphereDeploymentZoneList_To_v1alpha3_VSphereDeploymentZoneList(src, dst, nil)
}
