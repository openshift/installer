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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentity receiver to a v1beta1 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterStaticIdentity)
	if err := Convert_v1alpha3_AWSClusterStaticIdentity_To_v1beta1_AWSClusterStaticIdentity(r, dst, nil); err != nil {
		return err
	}

	dst.Spec.SecretRef = r.Spec.SecretRef.Name
	return nil
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentity receiver to a v1alpha3 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterStaticIdentity)

	if err := Convert_v1beta1_AWSClusterStaticIdentity_To_v1alpha3_AWSClusterStaticIdentity(src, r, nil); err != nil {
		return err
	}

	r.Spec.SecretRef.Name = src.Spec.SecretRef
	return nil
}

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentityList receiver to a v1beta1 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterStaticIdentityList)

	return Convert_v1alpha3_AWSClusterStaticIdentityList_To_v1beta1_AWSClusterStaticIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentityList receiver to a v1alpha3 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterStaticIdentityList)

	return Convert_v1beta1_AWSClusterStaticIdentityList_To_v1alpha3_AWSClusterStaticIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentity receiver to a v1beta1 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterRoleIdentity)

	return Convert_v1alpha3_AWSClusterRoleIdentity_To_v1beta1_AWSClusterRoleIdentity(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentity receiver to a v1alpha3 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterRoleIdentity)

	return Convert_v1beta1_AWSClusterRoleIdentity_To_v1alpha3_AWSClusterRoleIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentityList receiver to a v1beta1 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterRoleIdentityList)

	return Convert_v1alpha3_AWSClusterRoleIdentityList_To_v1beta1_AWSClusterRoleIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentityList receiver to a v1alpha3 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterRoleIdentityList)

	return Convert_v1beta1_AWSClusterRoleIdentityList_To_v1alpha3_AWSClusterRoleIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentity receiver to a v1beta1 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterControllerIdentity)

	return Convert_v1alpha3_AWSClusterControllerIdentity_To_v1beta1_AWSClusterControllerIdentity(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentity receiver to a v1alpha3 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterControllerIdentity)

	return Convert_v1beta1_AWSClusterControllerIdentity_To_v1alpha3_AWSClusterControllerIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentityList receiver to a v1beta1 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterControllerIdentityList)

	return Convert_v1alpha3_AWSClusterControllerIdentityList_To_v1beta1_AWSClusterControllerIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentityList receiver to a v1alpha3 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterControllerIdentityList)

	return Convert_v1beta1_AWSClusterControllerIdentityList_To_v1alpha3_AWSClusterControllerIdentityList(src, r, nil)
}

// Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec .
func Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec(in *AWSClusterStaticIdentitySpec, out *infrav1.AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec(in, out, s)
}

// Convert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec .
func Convert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in *infrav1.AWSClusterStaticIdentitySpec, out *AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in, out, s)
}
