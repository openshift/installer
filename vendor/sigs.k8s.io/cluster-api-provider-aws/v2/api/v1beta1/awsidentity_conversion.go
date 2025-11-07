/*
Copyright 2022 The Kubernetes Authors.

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
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// ConvertTo converts the v1beta1 AWSClusterControllerIdentity receiver to a v1beta2 AWSClusterControllerIdentity.
func (src *AWSClusterControllerIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterControllerIdentity)
	return Convert_v1beta1_AWSClusterControllerIdentity_To_v1beta2_AWSClusterControllerIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterControllerIdentity to a v1beta1 AWSClusterControllerIdentity.
func (dst *AWSClusterControllerIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterControllerIdentity)

	return Convert_v1beta2_AWSClusterControllerIdentity_To_v1beta1_AWSClusterControllerIdentity(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSClusterControllerIdentityList receiver to a v1beta2 AWSClusterControllerIdentityList.
func (src *AWSClusterControllerIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterControllerIdentityList)
	return Convert_v1beta1_AWSClusterControllerIdentityList_To_v1beta2_AWSClusterControllerIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterControllerIdentityList to a v1beta1 AWSClusterControllerIdentityList.
func (dst *AWSClusterControllerIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterControllerIdentityList)

	return Convert_v1beta2_AWSClusterControllerIdentityList_To_v1beta1_AWSClusterControllerIdentityList(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSClusterRoleIdentity receiver to a v1beta2 AWSClusterRoleIdentity.
func (src *AWSClusterRoleIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterRoleIdentity)
	return Convert_v1beta1_AWSClusterRoleIdentity_To_v1beta2_AWSClusterRoleIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterRoleIdentity to a v1beta1 AWSClusterRoleIdentity.
func (dst *AWSClusterRoleIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterRoleIdentity)

	return Convert_v1beta2_AWSClusterRoleIdentity_To_v1beta1_AWSClusterRoleIdentity(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSClusterRoleIdentityList receiver to a v1beta2 AWSClusterRoleIdentityList.
func (src *AWSClusterRoleIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterRoleIdentityList)
	return Convert_v1beta1_AWSClusterRoleIdentityList_To_v1beta2_AWSClusterRoleIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterRoleIdentityList to a v1beta1 AWSClusterRoleIdentityList.
func (dst *AWSClusterRoleIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterRoleIdentityList)

	return Convert_v1beta2_AWSClusterRoleIdentityList_To_v1beta1_AWSClusterRoleIdentityList(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSClusterStaticIdentity receiver to a v1beta2 AWSClusterStaticIdentity.
func (src *AWSClusterStaticIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterStaticIdentity)
	return Convert_v1beta1_AWSClusterStaticIdentity_To_v1beta2_AWSClusterStaticIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterStaticIdentity to a v1beta1 AWSClusterStaticIdentity.
func (dst *AWSClusterStaticIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterStaticIdentity)

	return Convert_v1beta2_AWSClusterStaticIdentity_To_v1beta1_AWSClusterStaticIdentity(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSClusterStaticIdentityList receiver to a v1beta2 AWSClusterStaticIdentityList.
func (src *AWSClusterStaticIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterStaticIdentityList)
	return Convert_v1beta1_AWSClusterStaticIdentityList_To_v1beta2_AWSClusterStaticIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterStaticIdentityList to a v1beta1 AWSClusterStaticIdentityList.
func (dst *AWSClusterStaticIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterStaticIdentityList)

	return Convert_v1beta2_AWSClusterStaticIdentityList_To_v1beta1_AWSClusterStaticIdentityList(src, dst, nil)
}
