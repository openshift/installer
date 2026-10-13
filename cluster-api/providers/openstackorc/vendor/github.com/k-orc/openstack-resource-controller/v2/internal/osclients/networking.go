/*
Copyright 2021 The ORC Authors.

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

package osclients

import (
	"context"
	"fmt"
	"iter"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/provider"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type NetworkExt struct {
	networks.Network
	dns.NetworkDNSExt
	external.NetworkExternalExt
	mtu.NetworkMTUExt
	portsecurity.PortSecurityExt
	provider.NetworkProviderExt
}

type PortExt struct {
	ports.Port
	portsecurity.PortSecurityExt
	portsbinding.PortsBindingExt
}

type NetworkClient interface {
	ListPort(ctx context.Context, opts ports.ListOptsBuilder) iter.Seq2[*PortExt, error]
	CreatePort(ctx context.Context, opts ports.CreateOptsBuilder) (*PortExt, error)
	DeletePort(ctx context.Context, id string) error
	GetPort(ctx context.Context, id string) (*PortExt, error)
	UpdatePort(ctx context.Context, id string, opts ports.UpdateOptsBuilder) (*PortExt, error)

	ListFloatingIP(ctx context.Context, opts floatingips.ListOptsBuilder) iter.Seq2[*floatingips.FloatingIP, error]
	CreateFloatingIP(ctx context.Context, opts floatingips.CreateOptsBuilder) (*floatingips.FloatingIP, error)
	DeleteFloatingIP(ctx context.Context, id string) error
	GetFloatingIP(ctx context.Context, id string) (*floatingips.FloatingIP, error)
	UpdateFloatingIP(ctx context.Context, id string, opts floatingips.UpdateOptsBuilder) (*floatingips.FloatingIP, error)

	ListRouter(ctx context.Context, opts routers.ListOpts) iter.Seq2[*routers.Router, error]
	CreateRouter(ctx context.Context, opts routers.CreateOptsBuilder) (*routers.Router, error)
	DeleteRouter(ctx context.Context, id string) error
	GetRouter(ctx context.Context, id string) (*routers.Router, error)
	UpdateRouter(ctx context.Context, id string, opts routers.UpdateOptsBuilder) (*routers.Router, error)
	AddRouterInterface(ctx context.Context, id string, opts routers.AddInterfaceOptsBuilder) (*routers.InterfaceInfo, error)
	RemoveRouterInterface(ctx context.Context, id string, opts routers.RemoveInterfaceOptsBuilder) (*routers.InterfaceInfo, error)

	ListSecGroup(ctx context.Context, opts groups.ListOpts) iter.Seq2[*groups.SecGroup, error]
	CreateSecGroup(ctx context.Context, opts groups.CreateOptsBuilder) (*groups.SecGroup, error)
	DeleteSecGroup(ctx context.Context, id string) error
	GetSecGroup(ctx context.Context, id string) (*groups.SecGroup, error)
	UpdateSecGroup(ctx context.Context, id string, opts groups.UpdateOptsBuilder) (*groups.SecGroup, error)

	ListSecGroupRule(ctx context.Context, opts rules.ListOpts) ([]rules.SecGroupRule, error)
	CreateSecGroupRules(ctx context.Context, opts []rules.CreateOpts) ([]rules.SecGroupRule, error)
	DeleteSecGroupRule(ctx context.Context, id string) error
	GetSecGroupRule(ctx context.Context, id string) (*rules.SecGroupRule, error)

	ListNetwork(ctx context.Context, opts networks.ListOptsBuilder) iter.Seq2[*NetworkExt, error]
	CreateNetwork(ctx context.Context, opts networks.CreateOptsBuilder) (*NetworkExt, error)
	DeleteNetwork(ctx context.Context, id string) error
	GetNetwork(ctx context.Context, id string) (*NetworkExt, error)
	UpdateNetwork(ctx context.Context, id string, opts networks.UpdateOptsBuilder) (*NetworkExt, error)

	ListSubnet(ctx context.Context, opts subnets.ListOptsBuilder) iter.Seq2[*subnets.Subnet, error]
	CreateSubnet(ctx context.Context, opts subnets.CreateOptsBuilder) (*subnets.Subnet, error)
	DeleteSubnet(ctx context.Context, id string) error
	GetSubnet(ctx context.Context, id string) (*subnets.Subnet, error)
	UpdateSubnet(ctx context.Context, id string, opts subnets.UpdateOptsBuilder) (*subnets.Subnet, error)

	ReplaceAllAttributesTags(ctx context.Context, resourceType string, resourceID string, opts attributestags.ReplaceAllOptsBuilder) ([]string, error)
}

type networkClient struct {
	serviceClient *gophercloud.ServiceClient
}

var _ NetworkClient = &networkClient{}

// NewNetworkClient returns an instance of the networking service.
func NewNetworkClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (NetworkClient, error) {
	serviceClient, err := openstack.NewNetworkV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create networking service providerClient: %v", err)
	}

	return networkClient{serviceClient}, nil
}

func (c networkClient) AddRouterInterface(ctx context.Context, id string, opts routers.AddInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	return routers.AddInterface(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) RemoveRouterInterface(ctx context.Context, id string, opts routers.RemoveInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	return routers.RemoveInterface(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) ReplaceAllAttributesTags(ctx context.Context, resourceType string, resourceID string, opts attributestags.ReplaceAllOptsBuilder) ([]string, error) {
	return attributestags.ReplaceAll(ctx, c.serviceClient, resourceType, resourceID, opts).Extract()
}

func (c networkClient) ListRouter(ctx context.Context, opts routers.ListOpts) iter.Seq2[*routers.Router, error] {
	pager := routers.List(c.serviceClient, opts)
	return func(yield func(*routers.Router, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(routers.ExtractRouters, yield))
	}
}

func (c networkClient) ListFloatingIP(ctx context.Context, opts floatingips.ListOptsBuilder) iter.Seq2[*floatingips.FloatingIP, error] {
	pager := floatingips.List(c.serviceClient, opts)
	return func(yield func(*floatingips.FloatingIP, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(floatingips.ExtractFloatingIPs, yield))
	}
}

func (c networkClient) CreateFloatingIP(ctx context.Context, opts floatingips.CreateOptsBuilder) (*floatingips.FloatingIP, error) {
	fip, err := floatingips.Create(ctx, c.serviceClient, opts).Extract()
	if err != nil {
		return nil, err
	}
	return fip, nil
}

func (c networkClient) DeleteFloatingIP(ctx context.Context, id string) error {
	return floatingips.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetFloatingIP(ctx context.Context, id string) (*floatingips.FloatingIP, error) {
	return floatingips.Get(ctx, c.serviceClient, id).Extract()
}

func (c networkClient) UpdateFloatingIP(ctx context.Context, id string, opts floatingips.UpdateOptsBuilder) (*floatingips.FloatingIP, error) {
	return floatingips.Update(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) ListPort(ctx context.Context, opts ports.ListOptsBuilder) iter.Seq2[*PortExt, error] {
	extractPortExt := func(p pagination.Page) ([]PortExt, error) {
		var resources []PortExt
		err := ports.ExtractPortsInto(p, &resources)
		if err != nil {
			return nil, err
		}
		return resources, nil
	}
	pager := ports.List(c.serviceClient, opts)
	return func(yield func(*PortExt, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(extractPortExt, yield))
	}
}

func (c networkClient) CreatePort(ctx context.Context, opts ports.CreateOptsBuilder) (*PortExt, error) {
	createResult := ports.Create(ctx, c.serviceClient, opts)
	portExt := PortExt{}
	if err := createResult.ExtractInto(&portExt); err != nil {
		return nil, err
	}
	return &portExt, nil
}

func (c networkClient) DeletePort(ctx context.Context, id string) error {
	return ports.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetPort(ctx context.Context, id string) (*PortExt, error) {
	portExt := PortExt{}
	if err := ports.Get(ctx, c.serviceClient, id).ExtractInto(&portExt); err != nil {
		return nil, err
	}
	return &portExt, nil
}

func (c networkClient) UpdatePort(ctx context.Context, id string, opts ports.UpdateOptsBuilder) (*PortExt, error) {
	portExt := PortExt{}
	if err := ports.Update(ctx, c.serviceClient, id, opts).ExtractInto(&portExt); err != nil {
		return nil, err
	}
	return &portExt, nil
}

func (c networkClient) CreateTrunk(ctx context.Context, opts trunks.CreateOptsBuilder) (*trunks.Trunk, error) {
	return trunks.Create(ctx, c.serviceClient, opts).Extract()
}

func (c networkClient) DeleteTrunk(ctx context.Context, id string) error {
	return trunks.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) ListTrunkSubports(ctx context.Context, trunkID string) ([]trunks.Subport, error) {
	return trunks.GetSubports(ctx, c.serviceClient, trunkID).Extract()
}

func (c networkClient) RemoveSubports(ctx context.Context, id string, opts trunks.RemoveSubportsOpts) error {
	_, err := trunks.RemoveSubports(ctx, c.serviceClient, id, opts).Extract()
	return err
}

func (c networkClient) ListTrunk(ctx context.Context, opts trunks.ListOptsBuilder) ([]trunks.Trunk, error) {
	allPages, err := trunks.List(c.serviceClient, opts).AllPages(ctx)
	if err != nil {
		return nil, err
	}
	return trunks.ExtractTrunks(allPages)
}

func (c networkClient) CreateRouter(ctx context.Context, opts routers.CreateOptsBuilder) (*routers.Router, error) {
	return routers.Create(ctx, c.serviceClient, opts).Extract()
}

func (c networkClient) DeleteRouter(ctx context.Context, id string) error {
	return routers.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetRouter(ctx context.Context, id string) (*routers.Router, error) {
	return routers.Get(ctx, c.serviceClient, id).Extract()
}

func (c networkClient) UpdateRouter(ctx context.Context, id string, opts routers.UpdateOptsBuilder) (*routers.Router, error) {
	return routers.Update(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) ListSecGroup(ctx context.Context, opts groups.ListOpts) iter.Seq2[*groups.SecGroup, error] {
	pager := groups.List(c.serviceClient, opts)
	return func(yield func(*groups.SecGroup, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(groups.ExtractGroups, yield))
	}
}

func (c networkClient) CreateSecGroup(ctx context.Context, opts groups.CreateOptsBuilder) (*groups.SecGroup, error) {
	return groups.Create(ctx, c.serviceClient, opts).Extract()
}

func (c networkClient) DeleteSecGroup(ctx context.Context, id string) error {
	return groups.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetSecGroup(ctx context.Context, id string) (*groups.SecGroup, error) {
	return groups.Get(ctx, c.serviceClient, id).Extract()
}

func (c networkClient) UpdateSecGroup(ctx context.Context, id string, opts groups.UpdateOptsBuilder) (*groups.SecGroup, error) {
	return groups.Update(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) ListSecGroupRule(ctx context.Context, opts rules.ListOpts) ([]rules.SecGroupRule, error) {
	allPages, err := rules.List(c.serviceClient, opts).AllPages(ctx)
	if err != nil {
		return nil, err
	}
	return rules.ExtractRules(allPages)
}

func (c networkClient) CreateSecGroupRules(ctx context.Context, opts []rules.CreateOpts) ([]rules.SecGroupRule, error) {
	return rules.CreateBulk(ctx, c.serviceClient, opts).Extract()
}

func (c networkClient) DeleteSecGroupRule(ctx context.Context, id string) error {
	return rules.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetSecGroupRule(ctx context.Context, id string) (*rules.SecGroupRule, error) {
	return rules.Get(ctx, c.serviceClient, id).Extract()
}

func (c networkClient) ListNetwork(ctx context.Context, opts networks.ListOptsBuilder) iter.Seq2[*NetworkExt, error] {
	extractNetworkExt := func(p pagination.Page) ([]NetworkExt, error) {
		var resources []NetworkExt
		err := networks.ExtractNetworksInto(p, &resources)
		if err != nil {
			return nil, err
		}
		return resources, nil
	}
	pager := networks.List(c.serviceClient, opts)
	return func(yield func(*NetworkExt, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(extractNetworkExt, yield))
	}
}

func (c networkClient) CreateNetwork(ctx context.Context, opts networks.CreateOptsBuilder) (*NetworkExt, error) {
	createResult := networks.Create(ctx, c.serviceClient, opts)
	networkExt := NetworkExt{}
	if err := createResult.ExtractInto(&networkExt); err != nil {
		return nil, err
	}
	return &networkExt, nil
}

func (c networkClient) DeleteNetwork(ctx context.Context, id string) error {
	return networks.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetNetwork(ctx context.Context, id string) (*NetworkExt, error) {
	networkExt := NetworkExt{}
	if err := networks.Get(ctx, c.serviceClient, id).ExtractInto(&networkExt); err != nil {
		return nil, err
	}
	return &networkExt, nil
}

func (c networkClient) UpdateNetwork(ctx context.Context, id string, opts networks.UpdateOptsBuilder) (*NetworkExt, error) {
	networkExt := NetworkExt{}
	if err := networks.Update(ctx, c.serviceClient, id, opts).ExtractInto(&networkExt); err != nil {
		return nil, err
	}
	return &networkExt, nil
}

func (c networkClient) ListSubnet(ctx context.Context, opts subnets.ListOptsBuilder) iter.Seq2[*subnets.Subnet, error] {
	pager := subnets.List(c.serviceClient, opts)
	return func(yield func(*subnets.Subnet, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(subnets.ExtractSubnets, yield))
	}
}

func (c networkClient) CreateSubnet(ctx context.Context, opts subnets.CreateOptsBuilder) (*subnets.Subnet, error) {
	return subnets.Create(ctx, c.serviceClient, opts).Extract()
}

func (c networkClient) DeleteSubnet(ctx context.Context, id string) error {
	return subnets.Delete(ctx, c.serviceClient, id).ExtractErr()
}

func (c networkClient) GetSubnet(ctx context.Context, id string) (*subnets.Subnet, error) {
	return subnets.Get(ctx, c.serviceClient, id).Extract()
}

func (c networkClient) UpdateSubnet(ctx context.Context, id string, opts subnets.UpdateOptsBuilder) (*subnets.Subnet, error) {
	return subnets.Update(ctx, c.serviceClient, id, opts).Extract()
}

func (c networkClient) ListExtensions(ctx context.Context) ([]extensions.Extension, error) {
	allPages, err := extensions.List(c.serviceClient).AllPages(ctx)
	if err != nil {
		return nil, err
	}
	return extensions.ExtractExtensions(allPages)
}
