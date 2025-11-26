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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// ConvertTo converts the v1beta1 AWSManagedControlPlane receiver to a v1beta2 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*ekscontrolplanev1.AWSManagedControlPlane)
	if err := Convert_v1beta1_AWSManagedControlPlane_To_v1beta2_AWSManagedControlPlane(r, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &ekscontrolplanev1.AWSManagedControlPlane{}
	if _, err := utilconversion.UnmarshalData(r, restored); err != nil {
		return err
	}

	dst.Spec.IdentityRef = r.Spec.IdentityRef
	dst.Spec.NetworkSpec = r.Spec.NetworkSpec
	dst.Spec.Region = r.Spec.Region
	dst.Spec.SSHKeyName = r.Spec.SSHKeyName
	dst.Spec.Version = r.Spec.Version
	dst.Spec.RoleName = r.Spec.RoleName
	dst.Spec.RoleAdditionalPolicies = r.Spec.RoleAdditionalPolicies
	dst.Spec.AdditionalTags = r.Spec.AdditionalTags

	if r.Spec.Logging != nil {
		dst.Spec.Logging = &ekscontrolplanev1.ControlPlaneLoggingSpec{}
		if err := Convert_v1beta1_ControlPlaneLoggingSpec_To_v1beta2_ControlPlaneLoggingSpec(r.Spec.Logging, dst.Spec.Logging, nil); err != nil {
			return err
		}
	}

	dst.Spec.SecondaryCidrBlock = r.Spec.SecondaryCidrBlock

	if err := Convert_v1beta1_EndpointAccess_To_v1beta2_EndpointAccess(&r.Spec.EndpointAccess, &dst.Spec.EndpointAccess, nil); err != nil {
		return err
	}

	if r.Spec.EncryptionConfig != nil {
		dst.Spec.EncryptionConfig = &ekscontrolplanev1.EncryptionConfig{}
		if err := Convert_v1beta1_EncryptionConfig_To_v1beta2_EncryptionConfig(r.Spec.EncryptionConfig, dst.Spec.EncryptionConfig, nil); err != nil {
			return err
		}
	}

	dst.Spec.ImageLookupFormat = r.Spec.ImageLookupFormat
	dst.Spec.ImageLookupOrg = r.Spec.ImageLookupOrg
	dst.Spec.ImageLookupBaseOS = r.Spec.ImageLookupBaseOS

	if r.Spec.IAMAuthenticatorConfig != nil {
		dst.Spec.IAMAuthenticatorConfig = &ekscontrolplanev1.IAMAuthenticatorConfig{}
		if err := Convert_v1beta1_IAMAuthenticatorConfig_To_v1beta2_IAMAuthenticatorConfig(r.Spec.IAMAuthenticatorConfig, dst.Spec.IAMAuthenticatorConfig, nil); err != nil {
			return err
		}
	}

	if r.Spec.Addons != nil {
		addons := []ekscontrolplanev1.Addon{}
		for _, addon := range *r.Spec.Addons {
			var convertedAddon ekscontrolplanev1.Addon
			if err := Convert_v1beta1_Addon_To_v1beta2_Addon(&addon, &convertedAddon, nil); err != nil {
				return err
			}
			addons = append(addons, convertedAddon)
		}
		dst.Spec.Addons = &addons
	}

	dst.Spec.Bastion = r.Spec.Bastion

	if r.Spec.TokenMethod != nil {
		Convert_v1beta1_EKSTokenMethod_To_v1beta2_EKSTokenMethod(r.Spec.TokenMethod, &dst.Spec.TokenMethod)
	}

	if err := Convert_v1beta1_VpcCni_To_v1beta2_VpcCni(&r.Spec.VpcCni, &dst.Spec.VpcCni, nil); err != nil {
		return err
	}
	dst.Spec.VpcCni.Disable = r.Spec.DisableVPCCNI

	if err := Convert_v1beta1_KubeProxy_To_v1beta2_KubeProxy(&r.Spec.KubeProxy, &dst.Spec.KubeProxy, nil); err != nil {
		return err
	}

	dst.Spec.AssociateOIDCProvider = r.Spec.AssociateOIDCProvider

	if r.Spec.OIDCIdentityProviderConfig != nil {
		dst.Spec.OIDCIdentityProviderConfig = &ekscontrolplanev1.OIDCIdentityProviderConfig{}
		if err := Convert_v1beta1_OIDCIdentityProviderConfig_To_v1beta2_OIDCIdentityProviderConfig(r.Spec.OIDCIdentityProviderConfig, dst.Spec.OIDCIdentityProviderConfig, nil); err != nil {
			return err
		}
	}

	dst.Spec.Partition = restored.Spec.Partition
	dst.Spec.RestrictPrivateSubnets = restored.Spec.RestrictPrivateSubnets
	dst.Spec.AccessConfig = restored.Spec.AccessConfig
	dst.Spec.RolePath = restored.Spec.RolePath
	dst.Spec.RolePermissionsBoundary = restored.Spec.RolePermissionsBoundary
	dst.Status.Version = restored.Status.Version
	dst.Spec.BootstrapSelfManagedAddons = restored.Spec.BootstrapSelfManagedAddons
	dst.Spec.UpgradePolicy = restored.Spec.UpgradePolicy
	return nil
}

// ConvertFrom converts the v1beta2 AWSManagedControlPlane receiver to a v1beta1 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*ekscontrolplanev1.AWSManagedControlPlane)
	if err := Convert_v1beta2_AWSManagedControlPlane_To_v1beta1_AWSManagedControlPlane(src, r, nil); err != nil {
		return err
	}

	r.Spec.IdentityRef = src.Spec.IdentityRef
	r.Spec.NetworkSpec = src.Spec.NetworkSpec
	r.Spec.Region = src.Spec.Region
	r.Spec.SSHKeyName = src.Spec.SSHKeyName
	r.Spec.Version = src.Spec.Version
	r.Spec.RoleName = src.Spec.RoleName
	r.Spec.RoleAdditionalPolicies = src.Spec.RoleAdditionalPolicies
	r.Spec.AdditionalTags = src.Spec.AdditionalTags

	if src.Spec.Logging != nil {
		r.Spec.Logging = &ControlPlaneLoggingSpec{}
		if err := Convert_v1beta2_ControlPlaneLoggingSpec_To_v1beta1_ControlPlaneLoggingSpec(src.Spec.Logging, r.Spec.Logging, nil); err != nil {
			return err
		}
	}

	r.Spec.SecondaryCidrBlock = src.Spec.SecondaryCidrBlock

	if err := Convert_v1beta2_EndpointAccess_To_v1beta1_EndpointAccess(&src.Spec.EndpointAccess, &r.Spec.EndpointAccess, nil); err != nil {
		return err
	}

	if src.Spec.EncryptionConfig != nil {
		r.Spec.EncryptionConfig = &EncryptionConfig{}
		if err := Convert_v1beta2_EncryptionConfig_To_v1beta1_EncryptionConfig(src.Spec.EncryptionConfig, r.Spec.EncryptionConfig, nil); err != nil {
			return err
		}
	}

	r.Spec.ImageLookupFormat = src.Spec.ImageLookupFormat
	r.Spec.ImageLookupOrg = src.Spec.ImageLookupOrg
	r.Spec.ImageLookupBaseOS = src.Spec.ImageLookupBaseOS

	if src.Spec.IAMAuthenticatorConfig != nil {
		r.Spec.IAMAuthenticatorConfig = &IAMAuthenticatorConfig{}
		if err := Convert_v1beta2_IAMAuthenticatorConfig_To_v1beta1_IAMAuthenticatorConfig(src.Spec.IAMAuthenticatorConfig, r.Spec.IAMAuthenticatorConfig, nil); err != nil {
			return err
		}
	}

	if src.Spec.Addons != nil {
		addons := []Addon{}
		for _, addon := range *src.Spec.Addons {
			var convertedAddon Addon
			if err := Convert_v1beta2_Addon_To_v1beta1_Addon(&addon, &convertedAddon, nil); err != nil {
				return err
			}
			addons = append(addons, convertedAddon)
		}
		r.Spec.Addons = &addons
	}

	r.Spec.Bastion = src.Spec.Bastion

	if src.Spec.TokenMethod != nil {
		Convert_v1beta2_EKSTokenMethod_To_v1beta1_EKSTokenMethod(src.Spec.TokenMethod, &r.Spec.TokenMethod)
	}

	if err := Convert_v1beta2_VpcCni_To_v1beta1_VpcCni(&src.Spec.VpcCni, &r.Spec.VpcCni, nil); err != nil {
		return err
	}
	r.Spec.DisableVPCCNI = src.Spec.VpcCni.Disable

	if err := Convert_v1beta2_KubeProxy_To_v1beta1_KubeProxy(&src.Spec.KubeProxy, &r.Spec.KubeProxy, nil); err != nil {
		return err
	}

	r.Spec.AssociateOIDCProvider = src.Spec.AssociateOIDCProvider

	if src.Spec.OIDCIdentityProviderConfig != nil {
		r.Spec.OIDCIdentityProviderConfig = &OIDCIdentityProviderConfig{}
		if err := Convert_v1beta2_OIDCIdentityProviderConfig_To_v1beta1_OIDCIdentityProviderConfig(src.Spec.OIDCIdentityProviderConfig, r.Spec.OIDCIdentityProviderConfig, nil); err != nil {
			return err
		}
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 AWSManagedControlPlaneList receiver to a v1beta2 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*ekscontrolplanev1.AWSManagedControlPlaneList)

	return Convert_v1beta1_AWSManagedControlPlaneList_To_v1beta2_AWSManagedControlPlaneList(r, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSManagedControlPlaneList receiver to a v1beta1 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*ekscontrolplanev1.AWSManagedControlPlaneList)

	return Convert_v1beta2_AWSManagedControlPlaneList_To_v1beta1_AWSManagedControlPlaneList(src, r, nil)
}

func Convert_v1beta1_AWSManagedControlPlaneSpec_To_v1beta2_AWSManagedControlPlaneSpec(in *AWSManagedControlPlaneSpec, out *ekscontrolplanev1.AWSManagedControlPlaneSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSManagedControlPlaneSpec_To_v1beta2_AWSManagedControlPlaneSpec(in, out, s)
}

func Convert_v1beta2_VpcCni_To_v1beta1_VpcCni(in *ekscontrolplanev1.VpcCni, out *VpcCni, s apiconversion.Scope) error {
	return autoConvert_v1beta2_VpcCni_To_v1beta1_VpcCni(in, out, s)
}

// Convert_v1beta2_AWSManagedControlPlaneSpec_To_v1beta1_AWSManagedControlPlaneSpec is a generated conversion function.
func Convert_v1beta2_AWSManagedControlPlaneSpec_To_v1beta1_AWSManagedControlPlaneSpec(in *ekscontrolplanev1.AWSManagedControlPlaneSpec, out *AWSManagedControlPlaneSpec, scope apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSManagedControlPlaneSpec_To_v1beta1_AWSManagedControlPlaneSpec(in, out, scope)
}

// Convert_v1beta2_AWSManagedControlPlaneStatus_To_v1beta1_AWSManagedControlPlaneStatus is an autogenerated conversion function.
func Convert_v1beta2_AWSManagedControlPlaneStatus_To_v1beta1_AWSManagedControlPlaneStatus(in *ekscontrolplanev1.AWSManagedControlPlaneStatus, out *AWSManagedControlPlaneStatus, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSManagedControlPlaneStatus_To_v1beta1_AWSManagedControlPlaneStatus(in, out, s)
}

func Convert_v1beta1_EKSTokenMethod_To_v1beta2_EKSTokenMethod(src *EKSTokenMethod, dst **ekscontrolplanev1.EKSTokenMethod) {
	if src == nil {
		*dst = nil
		return
	}
	tokenMethod := ekscontrolplanev1.EKSTokenMethod(*src)
	*dst = &tokenMethod
}

func Convert_v1beta2_EKSTokenMethod_To_v1beta1_EKSTokenMethod(src *ekscontrolplanev1.EKSTokenMethod, dst **EKSTokenMethod) {
	if src == nil {
		*dst = nil
		return
	}
	tokenMethod := EKSTokenMethod(*src)
	*dst = &tokenMethod
}
