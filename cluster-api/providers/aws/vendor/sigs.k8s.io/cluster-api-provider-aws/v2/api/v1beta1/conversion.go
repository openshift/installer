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

package v1beta1

import (
	"unsafe"

	"k8s.io/apimachinery/pkg/conversion"

	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func Convert_v1beta2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in *v1beta2.AWSClusterSpec, out *AWSClusterSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in, out, s)
}

func Convert_v1beta1_AWSResourceReference_To_v1beta2_AWSResourceReference(in *AWSResourceReference, out *v1beta2.AWSResourceReference, s conversion.Scope) error {
	return autoConvert_v1beta1_AWSResourceReference_To_v1beta2_AWSResourceReference(in, out, s)
}

func Convert_v1beta1_AWSMachineSpec_To_v1beta2_AWSMachineSpec(in *AWSMachineSpec, out *v1beta2.AWSMachineSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_AWSMachineSpec_To_v1beta2_AWSMachineSpec(in, out, s)
}

func Convert_v1beta2_AWSLoadBalancerSpec_To_v1beta1_AWSLoadBalancerSpec(in *v1beta2.AWSLoadBalancerSpec, out *AWSLoadBalancerSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_AWSLoadBalancerSpec_To_v1beta1_AWSLoadBalancerSpec(in, out, s)
}

func Convert_v1beta2_NetworkStatus_To_v1beta1_NetworkStatus(in *v1beta2.NetworkStatus, out *NetworkStatus, s conversion.Scope) error {
	return autoConvert_v1beta2_NetworkStatus_To_v1beta1_NetworkStatus(in, out, s)
}

func Convert_v1beta2_AWSMachineSpec_To_v1beta1_AWSMachineSpec(in *v1beta2.AWSMachineSpec, out *AWSMachineSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_AWSMachineSpec_To_v1beta1_AWSMachineSpec(in, out, s)
}

func Convert_v1beta2_Instance_To_v1beta1_Instance(in *v1beta2.Instance, out *Instance, s conversion.Scope) error {
	return autoConvert_v1beta2_Instance_To_v1beta1_Instance(in, out, s)
}

func Convert_v1beta1_ClassicELB_To_v1beta2_LoadBalancer(in *ClassicELB, out *v1beta2.LoadBalancer, s conversion.Scope) error {
	out.Name = in.Name
	out.DNSName = in.DNSName
	out.Scheme = v1beta2.ELBScheme(in.Scheme)
	out.HealthCheck = (*v1beta2.ClassicELBHealthCheck)(in.HealthCheck)
	out.AvailabilityZones = in.AvailabilityZones
	out.ClassicElbAttributes = (v1beta2.ClassicELBAttributes)(in.Attributes)
	out.ClassicELBListeners = *(*[]v1beta2.ClassicELBListener)(unsafe.Pointer(&in.Listeners))
	out.SecurityGroupIDs = in.SecurityGroupIDs
	out.Tags = in.Tags
	out.SubnetIDs = in.SubnetIDs
	return nil
}

func Convert_v1beta2_LoadBalancer_To_v1beta1_ClassicELB(in *v1beta2.LoadBalancer, out *ClassicELB, s conversion.Scope) error {
	out.Name = in.Name
	out.DNSName = in.DNSName
	out.Scheme = ClassicELBScheme(in.Scheme)
	out.HealthCheck = (*ClassicELBHealthCheck)(in.HealthCheck)
	out.AvailabilityZones = in.AvailabilityZones
	out.Attributes = (ClassicELBAttributes)(in.ClassicElbAttributes)
	out.Listeners = *(*[]ClassicELBListener)(unsafe.Pointer(&in.ClassicELBListeners))
	out.SecurityGroupIDs = in.SecurityGroupIDs
	out.Tags = in.Tags
	out.SubnetIDs = in.SubnetIDs
	return nil
}

func Convert_v1beta2_IngressRule_To_v1beta1_IngressRule(in *v1beta2.IngressRule, out *IngressRule, s conversion.Scope) error {
	return autoConvert_v1beta2_IngressRule_To_v1beta1_IngressRule(in, out, s)
}

func Convert_v1beta2_VPCSpec_To_v1beta1_VPCSpec(in *v1beta2.VPCSpec, out *VPCSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_VPCSpec_To_v1beta1_VPCSpec(in, out, s)
}

func Convert_v1beta2_IPv6_To_v1beta1_IPv6(in *v1beta2.IPv6, out *IPv6, s conversion.Scope) error {
	return autoConvert_v1beta2_IPv6_To_v1beta1_IPv6(in, out, s)
}

func Convert_v1beta2_NetworkSpec_To_v1beta1_NetworkSpec(in *v1beta2.NetworkSpec, out *NetworkSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_NetworkSpec_To_v1beta1_NetworkSpec(in, out, s)
}

func Convert_v1beta2_S3Bucket_To_v1beta1_S3Bucket(in *v1beta2.S3Bucket, out *S3Bucket, s conversion.Scope) error {
	return autoConvert_v1beta2_S3Bucket_To_v1beta1_S3Bucket(in, out, s)
}

func Convert_v1beta2_Ignition_To_v1beta1_Ignition(in *v1beta2.Ignition, out *Ignition, s conversion.Scope) error {
	return autoConvert_v1beta2_Ignition_To_v1beta1_Ignition(in, out, s)
}
