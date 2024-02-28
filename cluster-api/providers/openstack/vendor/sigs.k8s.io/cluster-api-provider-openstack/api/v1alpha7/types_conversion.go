/*
Copyright 2023 The Kubernetes Authors.

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

package v1alpha7

import (
	"errors"

	apiconversion "k8s.io/apimachinery/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

/* SecurityGroupFilter */

func restorev1alpha7SecurityGroupFilter(previous *SecurityGroupFilter, dst *SecurityGroupFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func restorev1alpha7SecurityGroup(previous *SecurityGroup, dst *SecurityGroup) {
	if previous == nil || dst == nil {
		return
	}

	dst.Rules = previous.Rules
}

func Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in *SecurityGroupFilter, out *infrav1.SecurityGroupFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter(in *infrav1.SecurityGroupFilter, out *SecurityGroupFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* NetworkFilter */

func restorev1alpha7NetworkFilter(previous *NetworkFilter, dst *NetworkFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func Convert_v1alpha7_NetworkFilter_To_v1beta1_NetworkFilter(in *NetworkFilter, out *infrav1.NetworkFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha7_NetworkFilter_To_v1beta1_NetworkFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_NetworkFilter_To_v1alpha7_NetworkFilter(in *infrav1.NetworkFilter, out *NetworkFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_NetworkFilter_To_v1alpha7_NetworkFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* SubnetFilter */

func restorev1alpha7SubnetFilter(previous *SubnetFilter, dst *SubnetFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func Convert_v1alpha7_SubnetFilter_To_v1beta1_SubnetFilter(in *SubnetFilter, out *infrav1.SubnetFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha7_SubnetFilter_To_v1beta1_SubnetFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_SubnetFilter_To_v1alpha7_SubnetFilter(in *infrav1.SubnetFilter, out *SubnetFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_SubnetFilter_To_v1alpha7_SubnetFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* RouterFilter */

func restorev1alpha7RouterFilter(previous *RouterFilter, dst *RouterFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func Convert_v1alpha7_RouterFilter_To_v1beta1_RouterFilter(in *RouterFilter, out *infrav1.RouterFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha7_RouterFilter_To_v1beta1_RouterFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_RouterFilter_To_v1alpha7_RouterFilter(in *infrav1.RouterFilter, out *RouterFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_RouterFilter_To_v1alpha7_RouterFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* PortOpts */

func restorev1alpha7Port(previous *PortOpts, dst *PortOpts) {
	if len(dst.SecurityGroupFilters) == len(previous.SecurityGroupFilters) {
		for i := range dst.SecurityGroupFilters {
			restorev1alpha7SecurityGroupFilter(&previous.SecurityGroupFilters[i], &dst.SecurityGroupFilters[i])
		}
	}

	if dst.Network != nil && previous.Network != nil {
		restorev1alpha7NetworkFilter(previous.Network, dst.Network)
	}

	if len(dst.FixedIPs) == len(previous.FixedIPs) {
		for i := range dst.FixedIPs {
			prevFixedIP := &previous.FixedIPs[i]
			dstFixedIP := &dst.FixedIPs[i]

			if dstFixedIP.Subnet != nil && prevFixedIP.Subnet != nil {
				restorev1alpha7SubnetFilter(prevFixedIP.Subnet, dstFixedIP.Subnet)
			}
		}
	}
}

func restorev1beta1Port(previous *infrav1.PortOpts, dst *infrav1.PortOpts) {
	optional.RestoreString(&previous.NameSuffix, &dst.NameSuffix)
	optional.RestoreString(&previous.Description, &dst.Description)
	optional.RestoreString(&previous.MACAddress, &dst.MACAddress)

	if len(dst.FixedIPs) == len(previous.FixedIPs) {
		for j := range dst.FixedIPs {
			prevFixedIP := &previous.FixedIPs[j]
			dstFixedIP := &dst.FixedIPs[j]

			if dstFixedIP.IPAddress == nil || *dstFixedIP.IPAddress == "" {
				dstFixedIP.IPAddress = prevFixedIP.IPAddress
			}
		}
	}

	if len(dst.AllowedAddressPairs) == len(previous.AllowedAddressPairs) {
		for j := range dst.AllowedAddressPairs {
			prevAAP := &previous.AllowedAddressPairs[j]
			dstAAP := &dst.AllowedAddressPairs[j]

			if dstAAP.MACAddress == nil || *dstAAP.MACAddress == "" {
				dstAAP.MACAddress = prevAAP.MACAddress
			}
		}
	}

	optional.RestoreString(&previous.HostID, &dst.HostID)
	optional.RestoreString(&previous.VNICType, &dst.VNICType)

	if dst.Profile == nil && previous.Profile != nil {
		dst.Profile = &infrav1.BindingProfile{}
	}

	if dst.Profile != nil && previous.Profile != nil {
		dstProfile := dst.Profile
		prevProfile := previous.Profile

		if dstProfile.OVSHWOffload == nil || !*dstProfile.OVSHWOffload {
			dstProfile.OVSHWOffload = prevProfile.OVSHWOffload
		}

		if dstProfile.TrustedVF == nil || !*dstProfile.TrustedVF {
			dstProfile.TrustedVF = prevProfile.TrustedVF
		}
	}
}

func Convert_v1alpha7_PortOpts_To_v1beta1_PortOpts(in *PortOpts, out *infrav1.PortOpts, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha7_PortOpts_To_v1beta1_PortOpts(in, out, s); err != nil {
		return err
	}

	// Copy members of ResolvedPortSpecFields
	var allowedAddressPairs []infrav1.AddressPair
	if len(in.AllowedAddressPairs) > 0 {
		allowedAddressPairs = make([]infrav1.AddressPair, len(in.AllowedAddressPairs))
		for i := range in.AllowedAddressPairs {
			aap := &in.AllowedAddressPairs[i]
			allowedAddressPairs[i] = infrav1.AddressPair{
				MACAddress: &aap.MACAddress,
				IPAddress:  aap.IPAddress,
			}
		}
	}
	var valueSpecs []infrav1.ValueSpec
	if len(in.ValueSpecs) > 0 {
		valueSpecs = make([]infrav1.ValueSpec, len(in.ValueSpecs))
		for i, vs := range in.ValueSpecs {
			valueSpecs[i] = infrav1.ValueSpec(vs)
		}
	}
	out.AdminStateUp = in.AdminStateUp
	out.AllowedAddressPairs = allowedAddressPairs
	out.DisablePortSecurity = in.DisablePortSecurity
	out.PropagateUplinkStatus = in.PropagateUplinkStatus
	out.ValueSpecs = valueSpecs
	if err := errors.Join(
		optional.Convert_string_To_optional_String(&in.MACAddress, &out.MACAddress, s),
		optional.Convert_string_To_optional_String(&in.HostID, &out.HostID, s),
		optional.Convert_string_To_optional_String(&in.VNICType, &out.VNICType, s),
	); err != nil {
		return err
	}

	if len(in.SecurityGroupFilters) > 0 {
		out.SecurityGroups = make([]infrav1.SecurityGroupFilter, len(in.SecurityGroupFilters))
		for i := range in.SecurityGroupFilters {
			if err := Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(&in.SecurityGroupFilters[i], &out.SecurityGroups[i], s); err != nil {
				return err
			}
		}
	}

	if in.Profile != (BindingProfile{}) {
		out.Profile = &infrav1.BindingProfile{}
		if err := Convert_v1alpha7_BindingProfile_To_v1beta1_BindingProfile(&in.Profile, out.Profile, s); err != nil {
			return err
		}
	}

	return nil
}

func Convert_v1beta1_PortOpts_To_v1alpha7_PortOpts(in *infrav1.PortOpts, out *PortOpts, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_PortOpts_To_v1alpha7_PortOpts(in, out, s); err != nil {
		return err
	}

	// Copy members of ResolvedPortSpecFields
	var allowedAddressPairs []AddressPair
	if len(in.AllowedAddressPairs) > 0 {
		allowedAddressPairs = make([]AddressPair, len(in.AllowedAddressPairs))
		for i := range in.AllowedAddressPairs {
			inAAP := &in.AllowedAddressPairs[i]
			outAAP := &allowedAddressPairs[i]
			if err := optional.Convert_optional_String_To_string(&inAAP.MACAddress, &outAAP.MACAddress, s); err != nil {
				return err
			}
			outAAP.IPAddress = inAAP.IPAddress
		}
	}
	var valueSpecs []ValueSpec
	if len(in.ValueSpecs) > 0 {
		valueSpecs = make([]ValueSpec, len(in.ValueSpecs))
		for i, vs := range in.ValueSpecs {
			valueSpecs[i] = ValueSpec(vs)
		}
	}
	out.AdminStateUp = in.AdminStateUp
	out.AllowedAddressPairs = allowedAddressPairs
	out.DisablePortSecurity = in.DisablePortSecurity
	out.PropagateUplinkStatus = in.PropagateUplinkStatus
	out.ValueSpecs = valueSpecs
	if err := errors.Join(
		optional.Convert_optional_String_To_string(&in.MACAddress, &out.MACAddress, s),
		optional.Convert_optional_String_To_string(&in.HostID, &out.HostID, s),
		optional.Convert_optional_String_To_string(&in.VNICType, &out.VNICType, s),
	); err != nil {
		return err
	}

	if len(in.SecurityGroups) > 0 {
		out.SecurityGroupFilters = make([]SecurityGroupFilter, len(in.SecurityGroups))
		for i := range in.SecurityGroups {
			if err := Convert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter(&in.SecurityGroups[i], &out.SecurityGroupFilters[i], s); err != nil {
				return err
			}
		}
	}

	if in.Profile != nil {
		if err := Convert_v1beta1_BindingProfile_To_v1alpha7_BindingProfile(in.Profile, &out.Profile, s); err != nil {
			return err
		}
	}

	return nil
}

/* SecurityGroup */

func Convert_v1alpha7_SecurityGroup_To_v1beta1_SecurityGroupStatus(in *SecurityGroup, out *infrav1.SecurityGroupStatus, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name

	return nil
}

func Convert_v1beta1_SecurityGroupStatus_To_v1alpha7_SecurityGroup(in *infrav1.SecurityGroupStatus, out *SecurityGroup, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	return nil
}

/* OpenStackIdentityReference */

func Convert_v1alpha7_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in *OpenStackIdentityReference, out *infrav1.OpenStackIdentityReference, s apiconversion.Scope) error {
	return autoConvert_v1alpha7_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in, out, s)
}

func Convert_v1beta1_OpenStackIdentityReference_To_v1alpha7_OpenStackIdentityReference(in *infrav1.OpenStackIdentityReference, out *OpenStackIdentityReference, _ apiconversion.Scope) error {
	out.Name = in.Name
	return nil
}
