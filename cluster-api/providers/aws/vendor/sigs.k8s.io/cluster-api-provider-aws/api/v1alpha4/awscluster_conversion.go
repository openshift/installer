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

package v1alpha4

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha4 AWSCluster receiver to a v1beta1 AWSCluster.
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSCluster)

	if err := Convert_v1alpha4_AWSCluster_To_v1beta1_AWSCluster(src, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.ControlPlaneLoadBalancer != nil {
		if dst.Spec.ControlPlaneLoadBalancer == nil {
			dst.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{}
		}
		restoreControlPlaneLoadBalancer(restored.Spec.ControlPlaneLoadBalancer, dst.Spec.ControlPlaneLoadBalancer)
	}

	dst.Spec.S3Bucket = restored.Spec.S3Bucket

	return nil
}

// restoreControlPlaneLoadBalancer manually restores the control plane loadbalancer data.
// Assumes restored and dst are non-nil.
func restoreControlPlaneLoadBalancer(restored, dst *infrav1.AWSLoadBalancerSpec) {
	dst.Name = restored.Name
	dst.HealthCheckProtocol = restored.HealthCheckProtocol
}

// ConvertFrom converts the v1beta1 AWSCluster receiver to a v1alpha4 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSCluster)

	if err := Convert_v1beta1_AWSCluster_To_v1alpha4_AWSCluster(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint .
func Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint(in *clusterv1alpha4.APIEndpoint, out *clusterv1.APIEndpoint, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint(in, out, s)
}

// Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint .
func Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint(in *clusterv1.APIEndpoint, out *clusterv1alpha4.APIEndpoint, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint(in, out, s)
}

// ConvertTo converts the v1alpha4 AWSClusterList receiver to a v1beta1 AWSClusterList.
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterList)

	return Convert_v1alpha4_AWSClusterList_To_v1beta1_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterList receiver to a v1alpha4 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterList)

	return Convert_v1beta1_AWSClusterList_To_v1alpha4_AWSClusterList(src, r, nil)
}

func Convert_v1beta1_AWSLoadBalancerSpec_To_v1alpha4_AWSLoadBalancerSpec(in *infrav1.AWSLoadBalancerSpec, out *AWSLoadBalancerSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSLoadBalancerSpec_To_v1alpha4_AWSLoadBalancerSpec(in, out, s)
}
