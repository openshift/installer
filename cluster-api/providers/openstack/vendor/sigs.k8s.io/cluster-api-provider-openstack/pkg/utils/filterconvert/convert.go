/*
Copyright 2024 The Kubernetes Authors.

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

package filterconvert

import (
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	securitygroups "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func SecurityGroupFilterToListOpts(securityGroupFilter *infrav1.SecurityGroupFilter) securitygroups.ListOpts {
	if securityGroupFilter == nil {
		return securitygroups.ListOpts{}
	}
	return securitygroups.ListOpts{
		Name:        securityGroupFilter.Name,
		Description: securityGroupFilter.Description,
		ProjectID:   securityGroupFilter.ProjectID,
		Tags:        infrav1.JoinTags(securityGroupFilter.Tags),
		TagsAny:     infrav1.JoinTags(securityGroupFilter.TagsAny),
		NotTags:     infrav1.JoinTags(securityGroupFilter.NotTags),
		NotTagsAny:  infrav1.JoinTags(securityGroupFilter.NotTagsAny),
	}
}

func SubnetFilterToListOpts(subnetFilter *infrav1.SubnetFilter) subnets.ListOpts {
	if subnetFilter == nil {
		return subnets.ListOpts{}
	}
	return subnets.ListOpts{
		Name:            subnetFilter.Name,
		Description:     subnetFilter.Description,
		ProjectID:       subnetFilter.ProjectID,
		IPVersion:       subnetFilter.IPVersion,
		GatewayIP:       subnetFilter.GatewayIP,
		CIDR:            subnetFilter.CIDR,
		IPv6AddressMode: subnetFilter.IPv6AddressMode,
		IPv6RAMode:      subnetFilter.IPv6RAMode,
		Tags:            infrav1.JoinTags(subnetFilter.Tags),
		TagsAny:         infrav1.JoinTags(subnetFilter.TagsAny),
		NotTags:         infrav1.JoinTags(subnetFilter.NotTags),
		NotTagsAny:      infrav1.JoinTags(subnetFilter.NotTagsAny),
	}
}

func NetworkFilterToListOpts(networkFilter *infrav1.NetworkFilter) networks.ListOpts {
	if networkFilter == nil {
		return networks.ListOpts{}
	}
	return networks.ListOpts{
		Name:        networkFilter.Name,
		Description: networkFilter.Description,
		ProjectID:   networkFilter.ProjectID,
		Tags:        infrav1.JoinTags(networkFilter.Tags),
		TagsAny:     infrav1.JoinTags(networkFilter.TagsAny),
		NotTags:     infrav1.JoinTags(networkFilter.NotTags),
		NotTagsAny:  infrav1.JoinTags(networkFilter.NotTagsAny),
	}
}

func RouterFilterToListOpts(routerFilter *infrav1.RouterFilter) routers.ListOpts {
	if routerFilter == nil {
		return routers.ListOpts{}
	}
	return routers.ListOpts{
		Name:        routerFilter.Name,
		Description: routerFilter.Description,
		ProjectID:   routerFilter.ProjectID,
		Tags:        infrav1.JoinTags(routerFilter.Tags),
		TagsAny:     infrav1.JoinTags(routerFilter.TagsAny),
		NotTags:     infrav1.JoinTags(routerFilter.NotTags),
		NotTagsAny:  infrav1.JoinTags(routerFilter.NotTagsAny),
	}
}

func ImageFilterToListOpts(imageFilter *infrav1.ImageFilter) (listOpts images.ListOpts) {
	if imageFilter == nil {
		return
	}

	if imageFilter.Name != nil && *imageFilter.Name != "" {
		listOpts.Name = *imageFilter.Name
	}

	if len(imageFilter.Tags) > 0 {
		listOpts.Tags = imageFilter.Tags
	}
	return
}
