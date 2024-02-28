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

package v1alpha6

import (
	"errors"
	"strings"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversioncommon"
	optional "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const trueString = "true"

/* ExternalRouterIPParam */
/* SecurityGroupParam, SecurityGroupFilter */

func restorev1alpha6SecurityGroupFilter(previous *SecurityGroupFilter, dst *SecurityGroupFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
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

func Convert_v1beta1_SecurityGroupParam_To_string(in *infrav1.SecurityGroupParam, out *string, _ apiconversion.Scope) error {
	if in.ID != nil {
		*out = *in.ID
	}
	return nil
}

func Convert_v1alpha6_SecurityGroupParam_To_v1beta1_SecurityGroupParam(in *SecurityGroupParam, out *infrav1.SecurityGroupParam, s apiconversion.Scope) error {
	if in.UUID != "" {
		out.ID = &in.UUID
		return nil
	}

	outFilter := &infrav1.SecurityGroupFilter{}

	if in.Name != "" {
		outFilter.Name = in.Name
	} else {
		err := Convert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(&in.Filter, outFilter, s)
		if err != nil {
			return err
		}
	}

	if !outFilter.IsZero() {
		out.Filter = outFilter
	}
	return nil
}

func Convert_v1beta1_SecurityGroupParam_To_v1alpha6_SecurityGroupParam(in *infrav1.SecurityGroupParam, out *SecurityGroupParam, s apiconversion.Scope) error {
	if in.ID != nil {
		out.UUID = *in.ID
		return nil
	}

	if in.Filter != nil {
		err := Convert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupFilter(in.Filter, &out.Filter, s)
		if err != nil {
			return err
		}
	}

	return nil
}

func Convert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in *SecurityGroupFilter, out *infrav1.SecurityGroupFilter, s apiconversion.Scope) error {
	err := autoConvert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in, out, s)
	if err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupFilter(in *infrav1.SecurityGroupFilter, out *SecurityGroupFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* NetworkParam, NetworkFilter */

func restorev1alpha6NetworkFilter(previous *NetworkFilter, dst *NetworkFilter) {
	if previous == nil || dst == nil {
		return
	}

	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny

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

func Convert_v1alpha6_NetworkFilter_To_v1beta1_NetworkParam(in *NetworkFilter, out *infrav1.NetworkParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}
	outFilter := &infrav1.NetworkFilter{}
	if err := autoConvert_v1alpha6_NetworkFilter_To_v1beta1_NetworkFilter(in, outFilter, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &outFilter.FilterByNeutronTags)
	if !outFilter.IsZero() {
		out.Filter = outFilter
	}
	return nil
}

func Convert_v1beta1_NetworkParam_To_v1alpha6_NetworkFilter(in *infrav1.NetworkParam, out *NetworkFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}
	if in.Filter != nil {
		if err := autoConvert_v1beta1_NetworkFilter_To_v1alpha6_NetworkFilter(in.Filter, out, s); err != nil {
			return err
		}
		infrav1.ConvertAllTagsFrom(&in.Filter.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	}
	return nil
}

/* SubnetParam, SubnetFilter */

func restorev1alpha6SubnetFilter(previous *SubnetFilter, dst *SubnetFilter) {
	if previous == nil || dst == nil {
		return
	}

	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny

	// We didn't convert other fields if ID was set
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

func restorev1alpha6SubnetParam(previous *SubnetParam, dst *SubnetParam) {
	if previous == nil || dst == nil {
		return
	}

	if previous.UUID != "" {
		dst.Filter = previous.Filter
	} else {
		restorev1alpha6SubnetFilter(&previous.Filter, &dst.Filter)
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

func Convert_v1alpha6_SubnetParam_To_v1beta1_SubnetParam(in *SubnetParam, out *infrav1.SubnetParam, s apiconversion.Scope) error {
	if in.UUID != "" {
		out.ID = &in.UUID
		return nil
	}

	outFilter := &infrav1.SubnetFilter{}
	if err := Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetFilter(&in.Filter, outFilter, s); err != nil {
		return err
	}
	if !outFilter.IsZero() {
		out.Filter = outFilter
	}
	return nil
}

func Convert_v1beta1_SubnetParam_To_v1alpha6_SubnetParam(in *infrav1.SubnetParam, out *SubnetParam, s apiconversion.Scope) error {
	if in.ID != nil {
		out.UUID = *in.ID
		return nil
	}

	if in.Filter != nil {
		if err := Convert_v1beta1_SubnetFilter_To_v1alpha6_SubnetFilter(in.Filter, &out.Filter, s); err != nil {
			return err
		}
	}

	return nil
}

func Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetFilter(in *SubnetFilter, out *infrav1.SubnetFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha6_SubnetFilter_To_v1beta1_SubnetFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_SubnetFilter_To_v1alpha6_SubnetFilter(in *infrav1.SubnetFilter, out *SubnetFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_SubnetFilter_To_v1alpha6_SubnetFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

func Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetParam(in *SubnetFilter, out *infrav1.SubnetParam, s apiconversion.Scope) error {
	if in.ID != "" {
		out.ID = &in.ID
		return nil
	}

	outFilter := &infrav1.SubnetFilter{}
	if err := Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetFilter(in, outFilter, s); err != nil {
		return err
	}
	if !outFilter.IsZero() {
		out.Filter = outFilter
	}
	return nil
}

func Convert_v1beta1_SubnetParam_To_v1alpha6_SubnetFilter(in *infrav1.SubnetParam, out *SubnetFilter, s apiconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
		return nil
	}

	if in.Filter != nil {
		if err := Convert_v1beta1_SubnetFilter_To_v1alpha6_SubnetFilter(in.Filter, out, s); err != nil {
			return err
		}
	}
	return nil
}

/* PortOpts, BindingProfile */

func restorev1alpha6Port(previous *PortOpts, dst *PortOpts) {
	if len(dst.SecurityGroupFilters) == len(previous.SecurityGroupFilters) {
		for i := range dst.SecurityGroupFilters {
			restorev1alpha6SecurityGroupFilter(&previous.SecurityGroupFilters[i].Filter, &dst.SecurityGroupFilters[i].Filter)
		}
	}

	restorev1alpha6NetworkFilter(previous.Network, dst.Network)

	if len(dst.FixedIPs) == len(previous.FixedIPs) {
		for i := range dst.FixedIPs {
			prevFixedIP := &previous.FixedIPs[i]
			dstFixedIP := &dst.FixedIPs[i]

			if dstFixedIP.Subnet != nil && prevFixedIP.Subnet != nil {
				restorev1alpha6SubnetFilter(prevFixedIP.Subnet, dstFixedIP.Subnet)
			}
		}
	}
}

func Convert_v1alpha6_PortOpts_To_v1beta1_PortOpts(in *PortOpts, out *infrav1.PortOpts, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha6_PortOpts_To_v1beta1_PortOpts(in, out, s); err != nil {
		return err
	}

	if len(in.SecurityGroups) > 0 || len(in.SecurityGroupFilters) > 0 {
		out.SecurityGroups = make([]infrav1.SecurityGroupParam, len(in.SecurityGroups)+len(in.SecurityGroupFilters))
		for i := range in.SecurityGroupFilters {
			if err := Convert_v1alpha6_SecurityGroupParam_To_v1beta1_SecurityGroupParam(&in.SecurityGroupFilters[i], &out.SecurityGroups[i], s); err != nil {
				return err
			}
		}
		for i := range in.SecurityGroups {
			out.SecurityGroups[i+len(in.SecurityGroupFilters)] = infrav1.SecurityGroupParam{ID: &in.SecurityGroups[i]}
		}
	}

	// Profile is now a struct in v1beta1.
	var ovsHWOffload, trustedVF bool
	if strings.Contains(in.Profile["capabilities"], "switchdev") {
		ovsHWOffload = true
	}
	if in.Profile["trusted"] == trueString {
		trustedVF = true
	}
	if ovsHWOffload || trustedVF {
		out.Profile = &infrav1.BindingProfile{}
		if ovsHWOffload {
			out.Profile.OVSHWOffload = &ovsHWOffload
		}
		if trustedVF {
			out.Profile.TrustedVF = &trustedVF
		}
	}

	return nil
}

func Convert_v1beta1_PortOpts_To_v1alpha6_PortOpts(in *infrav1.PortOpts, out *PortOpts, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_PortOpts_To_v1alpha6_PortOpts(in, out, s); err != nil {
		return err
	}

	// The auto-generated function converts v1beta1 SecurityGroup to
	// v1alpha6 SecurityGroup, but v1alpha6 SecurityGroupFilter is more
	// appropriate. Unset them and convert to SecurityGroupFilter instead.
	out.SecurityGroups = nil
	if len(in.SecurityGroups) > 0 {
		out.SecurityGroupFilters = make([]SecurityGroupParam, len(in.SecurityGroups))
		for i := range in.SecurityGroups {
			if err := Convert_v1beta1_SecurityGroupParam_To_v1alpha6_SecurityGroupParam(&in.SecurityGroups[i], &out.SecurityGroupFilters[i], s); err != nil {
				return err
			}
		}
	}

	if in.Profile != nil {
		out.Profile = make(map[string]string)
		if ptr.Deref(in.Profile.OVSHWOffload, false) {
			(out.Profile)["capabilities"] = "[\"switchdev\"]"
		}
		if ptr.Deref(in.Profile.TrustedVF, false) {
			(out.Profile)["trusted"] = trueString
		}
	}

	return nil
}

func Convert_Map_string_To_Interface_To_v1beta1_BindingProfile(in map[string]string, out *infrav1.BindingProfile, _ apiconversion.Scope) error {
	for k, v := range in {
		if k == "capabilities" {
			if strings.Contains(v, "switchdev") {
				out.OVSHWOffload = ptr.To(true)
			}
		}
		if k == "trusted" && v == trueString {
			out.TrustedVF = ptr.To(true)
		}
	}
	return nil
}

func Convert_v1beta1_BindingProfile_To_Map_string_To_Interface(in *infrav1.BindingProfile, out map[string]string, _ apiconversion.Scope) error {
	if ptr.Deref(in.OVSHWOffload, false) {
		(out)["capabilities"] = "[\"switchdev\"]"
	}
	if ptr.Deref(in.TrustedVF, false) {
		(out)["trusted"] = trueString
	}
	return nil
}

/* FixedIP */
/* AddressPair */
/* Instance */
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

func Convert_v1alpha6_RootVolume_To_v1beta1_RootVolume(in *RootVolume, out *infrav1.RootVolume, s apiconversion.Scope) error {
	out.SizeGiB = in.Size
	out.Type = in.VolumeType
	return conversioncommon.Convert_string_To_Pointer_v1beta1_VolumeAvailabilityZone(&in.AvailabilityZone, &out.AvailabilityZone, s)
}

func Convert_v1beta1_RootVolume_To_v1alpha6_RootVolume(in *infrav1.RootVolume, out *RootVolume, s apiconversion.Scope) error {
	out.Size = in.SizeGiB
	out.VolumeType = in.Type
	return conversioncommon.Convert_Pointer_v1beta1_VolumeAvailabilityZone_To_string(&in.AvailabilityZone, &out.AvailabilityZone, s)
}

/* Network */

func Convert_v1alpha6_Network_To_v1beta1_NetworkStatusWithSubnets(in *Network, out *infrav1.NetworkStatusWithSubnets, s apiconversion.Scope) error {
	// PortOpts has been removed in v1beta1
	err := Convert_v1alpha6_Network_To_v1beta1_NetworkStatus(in, &out.NetworkStatus, s)
	if err != nil {
		return err
	}

	if in.Subnet != nil {
		out.Subnets = []infrav1.Subnet{infrav1.Subnet(*in.Subnet)}
	}
	return nil
}

func Convert_v1beta1_NetworkStatusWithSubnets_To_v1alpha6_Network(in *infrav1.NetworkStatusWithSubnets, out *Network, s apiconversion.Scope) error {
	// PortOpts has been removed in v1beta1
	err := Convert_v1beta1_NetworkStatus_To_v1alpha6_Network(&in.NetworkStatus, out, s)
	if err != nil {
		return err
	}

	// Can only down-convert a single subnet
	if len(in.Subnets) > 0 {
		out.Subnet = (*Subnet)(&in.Subnets[0])
	}
	return nil
}

func Convert_v1alpha6_Network_To_v1beta1_NetworkStatus(in *Network, out *infrav1.NetworkStatus, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	out.Tags = in.Tags

	return nil
}

func Convert_v1beta1_NetworkStatus_To_v1alpha6_Network(in *infrav1.NetworkStatus, out *Network, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	out.Tags = in.Tags

	return nil
}

/* Subnet */
/* Router */
/* LoadBalancer */
/* SecurityGroup */

func restorev1alpha6SecurityGroup(previous *SecurityGroup, dst *SecurityGroup) {
	if previous == nil || dst == nil {
		return
	}

	dst.Rules = previous.Rules
}

func Convert_v1beta1_SecurityGroupStatus_To_v1alpha6_SecurityGroup(in *infrav1.SecurityGroupStatus, out *SecurityGroup, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	return nil
}

func Convert_v1alpha6_SecurityGroup_To_v1beta1_SecurityGroupStatus(in *SecurityGroup, out *infrav1.SecurityGroupStatus, _ apiconversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	return nil
}

/* SecurityGroupRule */
/* ValueSpec */
/* OpenStackIdentityReference */

func Convert_v1alpha6_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in *OpenStackIdentityReference, out *infrav1.OpenStackIdentityReference, s apiconversion.Scope) error {
	return autoConvert_v1alpha6_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in, out, s)
}

func Convert_v1beta1_OpenStackIdentityReference_To_v1alpha6_OpenStackIdentityReference(in *infrav1.OpenStackIdentityReference, out *OpenStackIdentityReference, _ apiconversion.Scope) error {
	out.Name = in.Name
	return nil
}

/* APIServerLoadBalancer */

func restorev1beta1APIServerLoadBalancer(previous *infrav1.APIServerLoadBalancer, dst *infrav1.APIServerLoadBalancer) {
	if dst == nil || previous == nil {
		return
	}

	// AZ doesn't exist in v1alpha6, so always restore.
	dst.AvailabilityZone = previous.AvailabilityZone
}

/* Placeholders */

// conversion-gen registers these functions so we must provider stubs, but
// nothing should ever call them

func Convert_v1alpha6_NetworkFilter_To_v1beta1_NetworkFilter(_ *NetworkFilter, _ *infrav1.NetworkFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha6_NetworkFilter_To_v1beta1_NetworkFilter should not be called")
}

func Convert_v1beta1_NetworkFilter_To_v1alpha6_NetworkFilter(_ *infrav1.NetworkFilter, _ *NetworkFilter, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_NetworkFilter_To_v1alpha6_NetworkFilter should not be called")
}

func Convert_v1alpha6_NetworkParam_To_v1beta1_NetworkParam(_ *NetworkParam, _ *infrav1.NetworkParam, _ apiconversion.Scope) error {
	return errors.New("Convert_v1alpha6_NetworkParam_To_v1beta1_NetworkParam should not be called")
}

func Convert_v1beta1_NetworkParam_To_v1alpha6_NetworkParam(_ *infrav1.NetworkParam, _ *NetworkParam, _ apiconversion.Scope) error {
	return errors.New("Convert_v1beta1_NetworkParam_To_v1alpha6_NetworkParam should not be called")
}
