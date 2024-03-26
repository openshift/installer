/*
Copyright 2022 Nutanix

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

package v1alpha4

import (
	infrav1beta1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this NutanixCluster to the Hub version (v1beta1).
func (src *NutanixCluster) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.NutanixCluster)
	return Convert_v1alpha4_NutanixCluster_To_v1beta1_NutanixCluster(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this NutanixCluster.
func (dst *NutanixCluster) ConvertFrom(srcRaw conversion.Hub) error { //nolint
	src := srcRaw.(*infrav1beta1.NutanixCluster)
	return Convert_v1beta1_NutanixCluster_To_v1alpha4_NutanixCluster(src, dst, nil)
}

// ConvertTo converts this NutanixClusterList to the Hub version (v1beta1).
func (src *NutanixClusterList) ConvertTo(dstRaw conversion.Hub) error { //nolint
	dst := dstRaw.(*infrav1beta1.NutanixClusterList)
	return Convert_v1alpha4_NutanixClusterList_To_v1beta1_NutanixClusterList(src, dst, nil)
}

// ConvertFrom converts from the Hub version(v1beta1) to this NutanixClusterList.
func (dst *NutanixClusterList) ConvertFrom(srcRaw conversion.Hub) error { //nolint
	src := srcRaw.(*infrav1beta1.NutanixClusterList)
	return Convert_v1beta1_NutanixClusterList_To_v1alpha4_NutanixClusterList(src, dst, nil)
}

// Convert_v1alpha4_NutanixClusterSpec_To_v1beta1_NutanixClusterSpec converts NutanixClusterSpec in NutanixClusterResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_NutanixClusterSpec_To_v1beta1_NutanixClusterSpec(in *NutanixClusterSpec, out *infrav1beta1.NutanixClusterSpec, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_NutanixClusterSpec_To_v1beta1_NutanixClusterSpec(in, out, s)
}

// Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint converts APIEndpoint in NutanixClusterResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint(in *capiv1alpha4.APIEndpoint, out *capiv1beta1.APIEndpoint, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint(in, out, s)
}

// Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint converts APIEndpoint in NutanixClusterResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint(in *capiv1beta1.APIEndpoint, out *capiv1alpha4.APIEndpoint, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint(in, out, s)
}

// Convert_v1beta1_NutanixClusterSpec_To_v1alpha4_NutanixClusterSpec converts NutanixClusterSpec in NutanixClusterResource from v1beta1 to v1alpha4 version.
//
//nolint:all
func Convert_v1beta1_NutanixClusterSpec_To_v1alpha4_NutanixClusterSpec(in *infrav1beta1.NutanixClusterSpec, out *NutanixClusterSpec, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1beta1_NutanixClusterSpec_To_v1alpha4_NutanixClusterSpec(in, out, s)
}

// Convert_v1alpha4_NutanixClusterStatus_To_v1beta1_NutanixClusterStatus converts NutanixClusterStatus in NutanixClusterResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_NutanixClusterStatus_To_v1beta1_NutanixClusterStatus(in *NutanixClusterStatus, out *infrav1beta1.NutanixClusterStatus, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_NutanixClusterStatus_To_v1beta1_NutanixClusterStatus(in, out, s)
}
