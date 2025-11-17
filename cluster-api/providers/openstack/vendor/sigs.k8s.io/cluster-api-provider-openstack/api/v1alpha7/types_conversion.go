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
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversioncommon"
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

	// If ID was set we lost all other filter params
	if dst.ID != "" {
		dst.Name = previous.Name
		dst.Description = previous.Description
		dst.ProjectID = previous.ProjectID
	}
}

func restorev1alpha7SecurityGroup(previous *SecurityGroup, dst *SecurityGroup) {
	if previous == nil || dst == nil {
		return
	}

	dst.Rules = previous.Rules
}

func restorev1beta1SecurityGroupParam(previous *infrav1.SecurityGroupParam, dst *infrav1.SecurityGroupParam) {
	if previous == nil || dst == nil {
		return
	}

	if dst.Filter != nil && previous.Filter != nil {
		dst.Filter.Tags = previous.Filter.Tags
		dst.Filter.TagsAny = previous.Filter.TagsAny
		dst.Filter.NotTags = previous.Filter.NotTags
		dst.Filter.NotTagsAny = previous.Filter.NotTagsAny
	}
}

func Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupParam(in *SecurityGroupFilter, out *infrav1.SecurityGroupParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}

	filter := &infrav1.SecurityGroupFilter{}
	if err := autoConvert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in, filter, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &filter.FilterByNeutronTags)
	if !filter.IsZero() {
		out.Filter = filter
	}
	return nil
}

func Convert_v1beta1_SecurityGroupParam_To_v1alpha7_SecurityGroupFilter(in *infrav1.SecurityGroupParam, out *SecurityGroupFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}

	if in.Filter != nil {
		if err := autoConvert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter(in.Filter, out, s); err != nil {
			return err
		}
		infrav1.ConvertAllTagsFrom(&in.Filter.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	}
	return nil
}

/* NetworkFilter */

func restorev1alpha7NetworkFilter(previous *NetworkFilter, dst *NetworkFilter) {
	if previous == nil || dst == nil {
		return
	}

	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny

	// If ID was set we lost all over filter params
	if dst.ID != "" {
		dst.Name = previous.Name
		dst.Description = previous.Description
		dst.ProjectID = previous.ProjectID
	}
}

func restorev1beta1NetworkParam(previous *infrav1.NetworkParam, dst *infrav1.NetworkParam) {
	if previous == nil || dst == nil {
		return
	}

	if dst.Filter != nil && previous.Filter != nil {
		dst.Filter.FilterByNeutronTags = previous.Filter.FilterByNeutronTags
	}
}

func Convert_v1alpha7_NetworkFilter_To_v1beta1_NetworkParam(in *NetworkFilter, out *infrav1.NetworkParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}
	outFilter := &infrav1.NetworkFilter{}
	if err := autoConvert_v1alpha7_NetworkFilter_To_v1beta1_NetworkFilter(in, outFilter, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &outFilter.FilterByNeutronTags)
	if !outFilter.IsZero() {
		out.Filter = outFilter
	}
	return nil
}

func Convert_v1beta1_NetworkParam_To_v1alpha7_NetworkFilter(in *infrav1.NetworkParam, out *NetworkFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}

	if in.Filter != nil {
		if err := autoConvert_v1beta1_NetworkFilter_To_v1alpha7_NetworkFilter(in.Filter, out, s); err != nil {
			return err
		}
		infrav1.ConvertAllTagsFrom(&in.Filter.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	}
	return nil
}

/* SubnetFilter */

func restorev1alpha7SubnetFilter(previous *SubnetFilter, dst *SubnetFilter) {
	if previous == nil || dst == nil {
		return
	}

	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny

	// If ID was set we will have lost all other fields in up-conversion
	if previous.ID != "" {
		dst.Name = previous.Name
		dst.Description = previous.Description
		dst.ProjectID = previous.ProjectID
		dst.IPVersion = previous.IPVersion
		dst.GatewayIP = previous.GatewayIP
		dst.CIDR = previous.CIDR
		dst.IPv6AddressMode = previous.IPv6AddressMode
		dst.IPv6RAMode = previous.IPv6RAMode
	}
}

func restorev1beta1SubnetParam(previous *infrav1.SubnetParam, dst *infrav1.SubnetParam) {
	if previous == nil || dst == nil {
		return
	}

	optional.RestoreString(&previous.ID, &dst.ID)

	if previous.Filter != nil && dst.Filter != nil {
		dst.Filter.FilterByNeutronTags = previous.Filter.FilterByNeutronTags
	}
}

func Convert_v1alpha7_SubnetFilter_To_v1beta1_SubnetParam(in *SubnetFilter, out *infrav1.SubnetParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}

	filter := &infrav1.SubnetFilter{}
	if err := autoConvert_v1alpha7_SubnetFilter_To_v1beta1_SubnetFilter(in, filter, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &filter.FilterByNeutronTags)

	if !filter.IsZero() {
		out.Filter = filter
	}
	return nil
}

func Convert_v1beta1_SubnetParam_To_v1alpha7_SubnetFilter(in *infrav1.SubnetParam, out *SubnetFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}

	if in.Filter == nil {
		return nil
	}

	if err := autoConvert_v1beta1_SubnetFilter_To_v1alpha7_SubnetFilter(in.Filter, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.Filter.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
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

	// If ID was set we lost all other filter params
	if dst.ID != "" {
		dst.Name = previous.Name
		dst.Description = previous.Description
		dst.ProjectID = previous.ProjectID
	}
}

func restorev1beta1RouterParam(previous *infrav1.RouterParam, dst *infrav1.RouterParam) {
	if previous == nil || dst == nil {
		return
	}

	optional.RestoreString(&previous.ID, &dst.ID)
	if previous.Filter != nil && dst.Filter != nil {
		dst.Filter.FilterByNeutronTags = previous.Filter.FilterByNeutronTags
	}
}

func Convert_v1alpha7_RouterFilter_To_v1beta1_RouterParam(in *RouterFilter, out *infrav1.RouterParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}

	filter := &infrav1.RouterFilter{}
	if err := autoConvert_v1alpha7_RouterFilter_To_v1beta1_RouterFilter(in, filter, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &filter.FilterByNeutronTags)
	if !filter.IsZero() {
		out.Filter = filter
	}
	return nil
}

func Convert_v1beta1_RouterParam_To_v1alpha7_RouterFilter(in *infrav1.RouterParam, out *RouterFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}

	if in.Filter != nil {
		if err := autoConvert_v1beta1_RouterFilter_To_v1alpha7_RouterFilter(in.Filter, out, s); err != nil {
			return err
		}
		infrav1.ConvertAllTagsFrom(&in.Filter.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	}
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
	restorev1beta1NetworkParam(previous.Network, dst.Network)

	optional.RestoreString(&previous.NameSuffix, &dst.NameSuffix)
	optional.RestoreString(&previous.Description, &dst.Description)
	optional.RestoreString(&previous.MACAddress, &dst.MACAddress)

	if len(dst.FixedIPs) == len(previous.FixedIPs) {
		for j := range dst.FixedIPs {
			prevFixedIP := &previous.FixedIPs[j]
			dstFixedIP := &dst.FixedIPs[j]

			optional.RestoreString(&prevFixedIP.IPAddress, &dstFixedIP.IPAddress)
			restorev1beta1SubnetParam(prevFixedIP.Subnet, dstFixedIP.Subnet)
		}
	}

	if len(dst.AllowedAddressPairs) == len(previous.AllowedAddressPairs) {
		for j := range dst.AllowedAddressPairs {
			prevAAP := &previous.AllowedAddressPairs[j]
			dstAAP := &dst.AllowedAddressPairs[j]

			optional.RestoreString(&prevAAP.MACAddress, &dstAAP.MACAddress)
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

	if len(dst.SecurityGroups) == len(previous.SecurityGroups) {
		for j := range dst.SecurityGroups {
			restorev1beta1SecurityGroupParam(&previous.SecurityGroups[j], &dst.SecurityGroups[j])
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
		out.SecurityGroups = make([]infrav1.SecurityGroupParam, len(in.SecurityGroupFilters))
		for i := range in.SecurityGroupFilters {
			if err := Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupParam(&in.SecurityGroupFilters[i], &out.SecurityGroups[i], s); err != nil {
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
			if err := Convert_v1beta1_SecurityGroupParam_To_v1alpha7_SecurityGroupFilter(&in.SecurityGroups[i], &out.SecurityGroupFilters[i], s); err != nil {
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

/* RootVolume */

func restorev1beta1BlockDeviceVolume(previous *infrav1.BlockDeviceVolume, dst *infrav1.BlockDeviceVolume) {
	if previous == nil || dst == nil {
		return
	}

	dstAZ := dst.AvailabilityZone
	previousAZ := previous.AvailabilityZone

	// Empty From (the default) will be converted to the explicit "Name"
	if dstAZ != nil && previousAZ != nil && dstAZ.From == "Name" {
		dstAZ.From = previousAZ.From
	}
}

func Convert_v1alpha7_RootVolume_To_v1beta1_RootVolume(in *RootVolume, out *infrav1.RootVolume, s apiconversion.Scope) error {
	out.SizeGiB = in.Size
	out.Type = in.VolumeType
	return conversioncommon.Convert_string_To_Pointer_v1beta1_VolumeAvailabilityZone(&in.AvailabilityZone, &out.AvailabilityZone, s)
}

func Convert_v1beta1_RootVolume_To_v1alpha7_RootVolume(in *infrav1.RootVolume, out *RootVolume, s apiconversion.Scope) error {
	out.Size = in.SizeGiB
	out.VolumeType = in.Type
	return conversioncommon.Convert_Pointer_v1beta1_VolumeAvailabilityZone_To_string(&in.AvailabilityZone, &out.AvailabilityZone, s)
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
	// Kind will be overwritten during restore if it was previously set to some other value, but if not then some value is still required
	out.Kind = "Secret"
	return nil
}

/* APIServerLoadBalancer */

func restorev1beta1APIServerLoadBalancer(previous *infrav1.APIServerLoadBalancer, dst *infrav1.APIServerLoadBalancer) {
	if dst == nil || previous == nil {
		return
	}

	// AZ doesn't exist in v1alpha7, so always restore.
	dst.AvailabilityZone = previous.AvailabilityZone
}

/* Placeholders */

// conversion-gen registers these functions so we must provider stubs, but
// nothing should ever call them

func Convert_v1alpha7_SubnetFilter_To_v1beta1_SubnetFilter(_ *SubnetFilter, _ *infrav1.SubnetFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha7_SubnetFilter_To_v1beta1_SubnetFilter should not be called")
}

func Convert_v1beta1_SubnetFilter_To_v1alpha7_SubnetFilter(_ *infrav1.SubnetFilter, _ *SubnetFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_SubnetFilter_To_v1alpha7_SubnetFilter should not be called")
}

func Convert_v1alpha7_NetworkFilter_To_v1beta1_NetworkFilter(_ *NetworkFilter, _ *infrav1.NetworkFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha7_NetworkFilter_To_v1beta1_NetworkFilter should not be called")
}

func Convert_v1beta1_NetworkFilter_To_v1alpha7_NetworkFilter(_ *infrav1.NetworkFilter, _ *NetworkFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_NetworkFilter_To_v1alpha7_NetworkFilter should not be called")
}

func Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(_ *SecurityGroupFilter, _ *infrav1.SecurityGroupFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha7_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter should not be called")
}

func Convert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter(_ *infrav1.SecurityGroupFilter, _ *SecurityGroupFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_SecurityGroupFilter_To_v1alpha7_SecurityGroupFilter should not be called")
}

func Convert_v1alpha7_RouterFilter_To_v1beta1_RouterFilter(_ *RouterFilter, _ *infrav1.RouterFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha7_RouterFilter_To_v1beta1_RouterFilter should not be called")
}

func Convert_v1beta1_RouterFilter_To_v1alpha7_RouterFilter(_ *infrav1.RouterFilter, _ *RouterFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_RouterFilter_To_v1alpha7_RouterFilter should not be called")
}
