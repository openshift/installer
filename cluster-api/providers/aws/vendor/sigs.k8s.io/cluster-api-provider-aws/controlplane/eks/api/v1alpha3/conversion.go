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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	clusterapiapiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clusterapiapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 AWSManagedControlPlane receiver to a v1beta1 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedControlPlane)

	if err := Convert_v1alpha3_AWSManagedControlPlane_To_v1beta1_AWSManagedControlPlane(r, dst, nil); err != nil {
		return err
	}

	restored := &v1beta1.AWSManagedControlPlane{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	dst.Status.IdentityProviderStatus = restored.Status.IdentityProviderStatus
	dst.Status.Bastion = restored.Status.Bastion
	dst.Spec.OIDCIdentityProviderConfig = restored.Spec.OIDCIdentityProviderConfig
	dst.Spec.KubeProxy = restored.Spec.KubeProxy

	return nil
}

// ConvertFrom converts the v1beta1 AWSManagedControlPlane receiver to a v1alpha3 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSManagedControlPlane)

	if err := Convert_v1beta1_AWSManagedControlPlane_To_v1alpha3_AWSManagedControlPlane(src, r, nil); err != nil {
		return err
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha3 AWSManagedControlPlaneList receiver to a v1beta1 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedControlPlaneList)

	return Convert_v1alpha3_AWSManagedControlPlaneList_To_v1beta1_AWSManagedControlPlaneList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedControlPlaneList receiver to a v1alpha3 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSManagedControlPlaneList)

	return Convert_v1beta1_AWSManagedControlPlaneList_To_v1alpha3_AWSManagedControlPlaneList(src, r, nil)
}

// Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint is a conversion function.
func Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint(in *clusterapiapiv1alpha3.APIEndpoint, out *clusterapiapiv1beta1.APIEndpoint, s apiconversion.Scope) error {
	return clusterapiapiv1alpha3.Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint(in, out, s)
}

// Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint is a conversion function.
func Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint(in *clusterapiapiv1beta1.APIEndpoint, out *clusterapiapiv1alpha3.APIEndpoint, s apiconversion.Scope) error {
	return clusterapiapiv1alpha3.Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint(in, out, s)
}

// Convert_v1alpha3_Bastion_To_v1beta1_Bastion is a conversion function.
func Convert_v1alpha3_Bastion_To_v1beta1_Bastion(in *infrav1alpha3.Bastion, out *infrav1beta1.Bastion, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Bastion_To_v1beta1_Bastion(in, out, s)
}

// Convert_v1beta1_Bastion_To_v1alpha3_Bastion is a conversion function.
func Convert_v1beta1_Bastion_To_v1alpha3_Bastion(in *infrav1beta1.Bastion, out *infrav1alpha3.Bastion, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1beta1_Bastion_To_v1alpha3_Bastion(in, out, s)
}

// Convert_v1alpha3_NetworkSpec_To_v1beta1_NetworkSpec is a conversion function.
func Convert_v1alpha3_NetworkSpec_To_v1beta1_NetworkSpec(in *infrav1alpha3.NetworkSpec, out *infrav1beta1.NetworkSpec, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_NetworkSpec_To_v1beta1_NetworkSpec(in, out, s)
}

// Convert_v1beta1_NetworkSpec_To_v1alpha3_NetworkSpec is a conversion function.
func Convert_v1beta1_NetworkSpec_To_v1alpha3_NetworkSpec(in *infrav1beta1.NetworkSpec, out *infrav1alpha3.NetworkSpec, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1beta1_NetworkSpec_To_v1alpha3_NetworkSpec(in, out, s)
}

// Convert_v1beta1_Instance_To_v1alpha3_Instance is a conversion function.
func Convert_v1beta1_Instance_To_v1alpha3_Instance(in *infrav1beta1.Instance, out *infrav1alpha3.Instance, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1beta1_Instance_To_v1alpha3_Instance(in, out, s)
}

// Convert_v1alpha3_Instance_To_v1beta1_Instance is a conversion function.
func Convert_v1alpha3_Instance_To_v1beta1_Instance(in *infrav1alpha3.Instance, out *infrav1beta1.Instance, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Instance_To_v1beta1_Instance(in, out, s)
}

// Convert_v1alpha3_Network_To_v1beta1_NetworkStatus is a conversion function.
func Convert_v1alpha3_Network_To_v1beta1_NetworkStatus(in *infrav1alpha3.Network, out *infrav1beta1.NetworkStatus, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Network_To_v1beta1_NetworkStatus(in, out, s)
}

// Convert_v1beta1_NetworkStatus_To_v1alpha3_Network is a conversion function.
func Convert_v1beta1_NetworkStatus_To_v1alpha3_Network(in *infrav1beta1.NetworkStatus, out *infrav1alpha3.Network, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1beta1_NetworkStatus_To_v1alpha3_Network(in, out, s)
}

func Convert_v1beta1_AWSManagedControlPlaneSpec_To_v1alpha3_AWSManagedControlPlaneSpec(in *v1beta1.AWSManagedControlPlaneSpec, out *AWSManagedControlPlaneSpec, scope apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSManagedControlPlaneSpec_To_v1alpha3_AWSManagedControlPlaneSpec(in, out, scope)
}

func Convert_v1beta1_AWSManagedControlPlaneStatus_To_v1alpha3_AWSManagedControlPlaneStatus(in *v1beta1.AWSManagedControlPlaneStatus, out *AWSManagedControlPlaneStatus, scope apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSManagedControlPlaneStatus_To_v1alpha3_AWSManagedControlPlaneStatus(in, out, scope)
}
