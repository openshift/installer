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

package v1alpha6

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func (subnetFilter SubnetFilter) ToListOpt() subnets.ListOpts {
	return subnets.ListOpts{
		Name:            subnetFilter.Name,
		Description:     subnetFilter.Description,
		ProjectID:       subnetFilter.ProjectID,
		IPVersion:       subnetFilter.IPVersion,
		GatewayIP:       subnetFilter.GatewayIP,
		CIDR:            subnetFilter.CIDR,
		IPv6AddressMode: subnetFilter.IPv6AddressMode,
		IPv6RAMode:      subnetFilter.IPv6RAMode,
		ID:              subnetFilter.ID,
		Tags:            subnetFilter.Tags,
		TagsAny:         subnetFilter.TagsAny,
		NotTags:         subnetFilter.NotTags,
		NotTagsAny:      subnetFilter.NotTagsAny,
	}
}

func (networkFilter NetworkFilter) ToListOpt() networks.ListOpts {
	return networks.ListOpts{
		Name:        networkFilter.Name,
		Description: networkFilter.Description,
		ProjectID:   networkFilter.ProjectID,
		ID:          networkFilter.ID,
		Tags:        networkFilter.Tags,
		TagsAny:     networkFilter.TagsAny,
		NotTags:     networkFilter.NotTags,
		NotTagsAny:  networkFilter.NotTagsAny,
	}
}
