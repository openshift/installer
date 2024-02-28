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
	"strings"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
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

func Convert_v1beta1_SecurityGroupFilter_To_string(in *infrav1.SecurityGroupFilter, out *string, _ apiconversion.Scope) error {
	if in.ID != "" {
		*out = in.ID
	}
	return nil
}

func Convert_v1alpha6_SecurityGroupParam_To_v1beta1_SecurityGroupFilter(in *SecurityGroupParam, out *infrav1.SecurityGroupFilter, s apiconversion.Scope) error {
	// SecurityGroupParam is replaced by its contained SecurityGroupFilter in v1beta1
	err := Convert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(&in.Filter, out, s)
	if err != nil {
		return err
	}

	if in.UUID != "" {
		out.ID = in.UUID
	}
	if in.Name != "" {
		out.Name = in.Name
	}
	return nil
}

func Convert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupParam(in *infrav1.SecurityGroupFilter, out *SecurityGroupParam, s apiconversion.Scope) error {
	// SecurityGroupParam is replaced by its contained SecurityGroupFilter in v1beta1
	err := Convert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupFilter(in, &out.Filter, s)
	if err != nil {
		return err
	}

	if in.ID != "" {
		out.UUID = in.ID
	}
	if in.Name != "" {
		out.Name = in.Name
	}
	return nil
}

func Convert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in *SecurityGroupFilter, out *infrav1.SecurityGroupFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)

	// TenantID has been removed in v1beta1. Write it to ProjectID if ProjectID is not already set.
	if out.ProjectID == "" {
		out.ProjectID = in.TenantID
	}
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
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func Convert_v1alpha6_NetworkFilter_To_v1beta1_NetworkFilter(in *NetworkFilter, out *infrav1.NetworkFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha6_NetworkFilter_To_v1beta1_NetworkFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsTo(in.Tags, in.TagsAny, in.NotTags, in.NotTagsAny, &out.FilterByNeutronTags)
	return nil
}

func Convert_v1beta1_NetworkFilter_To_v1alpha6_NetworkFilter(in *infrav1.NetworkFilter, out *NetworkFilter, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_NetworkFilter_To_v1alpha6_NetworkFilter(in, out, s); err != nil {
		return err
	}
	infrav1.ConvertAllTagsFrom(&in.FilterByNeutronTags, &out.Tags, &out.TagsAny, &out.NotTags, &out.NotTagsAny)
	return nil
}

/* SubnetParam, SubnetFilter */

func restorev1alpha6SubnetFilter(previous *SubnetFilter, dst *SubnetFilter) {
	// The edge cases with multiple commas are too tricky in this direction,
	// so we just restore the whole thing.
	dst.Tags = previous.Tags
	dst.TagsAny = previous.TagsAny
	dst.NotTags = previous.NotTags
	dst.NotTagsAny = previous.NotTagsAny
}

func Convert_v1alpha6_SubnetParam_To_v1beta1_SubnetFilter(in *SubnetParam, out *infrav1.SubnetFilter, s apiconversion.Scope) error {
	if err := Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetFilter(&in.Filter, out, s); err != nil {
		return err
	}
	if in.UUID != "" {
		out.ID = in.UUID
	}
	return nil
}

func Convert_v1beta1_SubnetFilter_To_v1alpha6_SubnetParam(in *infrav1.SubnetFilter, out *SubnetParam, s apiconversion.Scope) error {
	if err := Convert_v1beta1_SubnetFilter_To_v1alpha6_SubnetFilter(in, &out.Filter, s); err != nil {
		return err
	}
	out.UUID = in.ID

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

/* PortOpts, BindingProfile */

func restorev1alpha6Port(previous *PortOpts, dst *PortOpts) {
	if len(dst.SecurityGroupFilters) == len(previous.SecurityGroupFilters) {
		for i := range dst.SecurityGroupFilters {
			restorev1alpha6SecurityGroupFilter(&previous.SecurityGroupFilters[i].Filter, &dst.SecurityGroupFilters[i].Filter)
		}
	}

	if dst.Network != nil && previous.Network != nil {
		restorev1alpha6NetworkFilter(previous.Network, dst.Network)
	}

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
		out.SecurityGroups = make([]infrav1.SecurityGroupFilter, 0, len(in.SecurityGroups)+len(in.SecurityGroupFilters))
		for i := range in.SecurityGroupFilters {
			sgParam := &in.SecurityGroupFilters[i]
			switch {
			case sgParam.UUID != "":
				out.SecurityGroups = append(out.SecurityGroups, infrav1.SecurityGroupFilter{ID: sgParam.UUID})
			case sgParam.Name != "":
				out.SecurityGroups = append(out.SecurityGroups, infrav1.SecurityGroupFilter{Name: sgParam.Name})
			case sgParam.Filter != (SecurityGroupFilter{}):
				out.SecurityGroups = append(out.SecurityGroups, infrav1.SecurityGroupFilter{})
				outSG := &out.SecurityGroups[len(out.SecurityGroups)-1]
				if err := Convert_v1alpha6_SecurityGroupFilter_To_v1beta1_SecurityGroupFilter(&sgParam.Filter, outSG, s); err != nil {
					return err
				}
			}
		}
		for _, id := range in.SecurityGroups {
			out.SecurityGroups = append(out.SecurityGroups, infrav1.SecurityGroupFilter{ID: id})
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
			securityGroupParam := &out.SecurityGroupFilters[i]
			if in.SecurityGroups[i].ID != "" {
				securityGroupParam.UUID = in.SecurityGroups[i].ID
			} else {
				if err := Convert_v1beta1_SecurityGroupFilter_To_v1alpha6_SecurityGroupFilter(&in.SecurityGroups[i], &securityGroupParam.Filter, s); err != nil {
					return err
				}
			}
		}
	}

	if in.Profile != nil {
		out.Profile = make(map[string]string)
		if pointer.BoolDeref(in.Profile.OVSHWOffload, false) {
			(out.Profile)["capabilities"] = "[\"switchdev\"]"
		}
		if pointer.BoolDeref(in.Profile.TrustedVF, false) {
			(out.Profile)["trusted"] = trueString
		}
	}

	return nil
}

func Convert_Map_string_To_Interface_To_v1beta1_BindingProfile(in map[string]string, out *infrav1.BindingProfile, _ apiconversion.Scope) error {
	for k, v := range in {
		if k == "capabilities" {
			if strings.Contains(v, "switchdev") {
				out.OVSHWOffload = pointer.Bool(true)
			}
		}
		if k == "trusted" && v == trueString {
			out.TrustedVF = pointer.Bool(true)
		}
	}
	return nil
}

func Convert_v1beta1_BindingProfile_To_Map_string_To_Interface(in *infrav1.BindingProfile, out map[string]string, _ apiconversion.Scope) error {
	if pointer.BoolDeref(in.OVSHWOffload, false) {
		(out)["capabilities"] = "[\"switchdev\"]"
	}
	if pointer.BoolDeref(in.TrustedVF, false) {
		(out)["trusted"] = trueString
	}
	return nil
}

/* FixedIP */
/* AddressPair */
/* Instance */
/* RootVolume */
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
/* APIServerLoadBalancer */
/* ValueSpec */
/* OpenStackIdentityReference */

func Convert_v1alpha6_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in *OpenStackIdentityReference, out *infrav1.OpenStackIdentityReference, s apiconversion.Scope) error {
	return autoConvert_v1alpha6_OpenStackIdentityReference_To_v1beta1_OpenStackIdentityReference(in, out, s)
}

func Convert_v1beta1_OpenStackIdentityReference_To_v1alpha6_OpenStackIdentityReference(in *infrav1.OpenStackIdentityReference, out *OpenStackIdentityReference, _ apiconversion.Scope) error {
	out.Name = in.Name
	return nil
}
