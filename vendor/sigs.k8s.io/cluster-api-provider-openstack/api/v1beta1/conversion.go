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

package v1beta1

import (
	"strings"

	"k8s.io/utils/ptr"
)

// Hub marks OpenStackCluster as a conversion hub.
func (*OpenStackCluster) Hub() {}

// Hub marks OpenStackClusterList as a conversion hub.
func (*OpenStackClusterList) Hub() {}

// Hub marks OpenStackClusterTemplate as a conversion hub.
func (*OpenStackClusterTemplate) Hub() {}

// Hub marks OpenStackClusterTemplateList as a conversion hub.
func (*OpenStackClusterTemplateList) Hub() {}

// Hub marks OpenStackMachine as a conversion hub.
func (*OpenStackMachine) Hub() {}

// Hub marks OpenStackMachineList as a conversion hub.
func (*OpenStackMachineList) Hub() {}

// Hub marks OpenStackMachineTemplate as a conversion hub.
func (*OpenStackMachineTemplate) Hub() {}

// Hub marks OpenStackMachineTemplateList as a conversion hub.
func (*OpenStackMachineTemplateList) Hub() {}

// LegacyCalicoSecurityGroupRules returns a list of security group rules for calico
// that need to be applied to the control plane and worker security groups when
// managed security groups are enabled and upgrading to v1beta1.
func LegacyCalicoSecurityGroupRules() []SecurityGroupRuleSpec {
	return []SecurityGroupRuleSpec{
		{
			Name:                "BGP (calico)",
			Description:         ptr.To("Created by cluster-api-provider-openstack API conversion - BGP (calico)"),
			Direction:           "ingress",
			EtherType:           ptr.To("IPv4"),
			PortRangeMin:        ptr.To(179),
			PortRangeMax:        ptr.To(179),
			Protocol:            ptr.To("tcp"),
			RemoteManagedGroups: []ManagedSecurityGroupName{"controlplane", "worker"},
		},
		{
			Name:                "IP-in-IP (calico)",
			Description:         ptr.To("Created by cluster-api-provider-openstack API conversion - IP-in-IP (calico)"),
			Direction:           "ingress",
			EtherType:           ptr.To("IPv4"),
			Protocol:            ptr.To("4"),
			RemoteManagedGroups: []ManagedSecurityGroupName{"controlplane", "worker"},
		},
	}
}

// splitTags splits a comma separated list of tags into a slice of tags.
// If the input is an empty string, it returns nil representing no list rather
// than an empty list.
func splitTags(tags string) []NeutronTag {
	if tags == "" {
		return nil
	}

	var ret []NeutronTag
	for _, tag := range strings.Split(tags, ",") {
		if tag != "" {
			ret = append(ret, NeutronTag(tag))
		}
	}

	return ret
}

// JoinTags joins a slice of tags into a comma separated list of tags.
func JoinTags(tags []NeutronTag) string {
	var b strings.Builder
	for i := range tags {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(string(tags[i]))
	}
	return b.String()
}

func ConvertAllTagsTo(tags, tagsAny, notTags, notTagsAny string, neutronTags *FilterByNeutronTags) {
	neutronTags.Tags = splitTags(tags)
	neutronTags.TagsAny = splitTags(tagsAny)
	neutronTags.NotTags = splitTags(notTags)
	neutronTags.NotTagsAny = splitTags(notTagsAny)
}

func ConvertAllTagsFrom(neutronTags *FilterByNeutronTags, tags, tagsAny, notTags, notTagsAny *string) {
	*tags = JoinTags(neutronTags.Tags)
	*tagsAny = JoinTags(neutronTags.TagsAny)
	*notTags = JoinTags(neutronTags.NotTags)
	*notTagsAny = JoinTags(neutronTags.NotTagsAny)
}
