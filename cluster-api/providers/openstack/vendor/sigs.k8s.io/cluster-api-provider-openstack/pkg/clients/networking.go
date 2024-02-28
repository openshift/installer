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

package clients

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
)

type NetworkClient interface {
	ListFloatingIP(opts floatingips.ListOptsBuilder) ([]floatingips.FloatingIP, error)
	CreateFloatingIP(opts floatingips.CreateOptsBuilder) (*floatingips.FloatingIP, error)
	DeleteFloatingIP(id string) error
	GetFloatingIP(id string) (*floatingips.FloatingIP, error)
	UpdateFloatingIP(id string, opts floatingips.UpdateOptsBuilder) (*floatingips.FloatingIP, error)

	ListPort(opts ports.ListOptsBuilder) ([]ports.Port, error)
	CreatePort(opts ports.CreateOptsBuilder) (*ports.Port, error)
	DeletePort(id string) error
	GetPort(id string) (*ports.Port, error)
	UpdatePort(id string, opts ports.UpdateOptsBuilder) (*ports.Port, error)

	ListTrunk(opts trunks.ListOptsBuilder) ([]trunks.Trunk, error)
	CreateTrunk(opts trunks.CreateOptsBuilder) (*trunks.Trunk, error)
	DeleteTrunk(id string) error

	ListRouter(opts routers.ListOpts) ([]routers.Router, error)
	CreateRouter(opts routers.CreateOptsBuilder) (*routers.Router, error)
	DeleteRouter(id string) error
	GetRouter(id string) (*routers.Router, error)
	UpdateRouter(id string, opts routers.UpdateOptsBuilder) (*routers.Router, error)
	AddRouterInterface(id string, opts routers.AddInterfaceOptsBuilder) (*routers.InterfaceInfo, error)
	RemoveRouterInterface(id string, opts routers.RemoveInterfaceOptsBuilder) (*routers.InterfaceInfo, error)

	ListSecGroup(opts groups.ListOpts) ([]groups.SecGroup, error)
	CreateSecGroup(opts groups.CreateOptsBuilder) (*groups.SecGroup, error)
	DeleteSecGroup(id string) error
	GetSecGroup(id string) (*groups.SecGroup, error)
	UpdateSecGroup(id string, opts groups.UpdateOptsBuilder) (*groups.SecGroup, error)

	ListSecGroupRule(opts rules.ListOpts) ([]rules.SecGroupRule, error)
	CreateSecGroupRule(opts rules.CreateOptsBuilder) (*rules.SecGroupRule, error)
	DeleteSecGroupRule(id string) error
	GetSecGroupRule(id string) (*rules.SecGroupRule, error)

	ListNetwork(opts networks.ListOptsBuilder) ([]networks.Network, error)
	CreateNetwork(opts networks.CreateOptsBuilder) (*networks.Network, error)
	DeleteNetwork(id string) error
	GetNetwork(id string) (*networks.Network, error)
	UpdateNetwork(id string, opts networks.UpdateOptsBuilder) (*networks.Network, error)

	ListSubnet(opts subnets.ListOptsBuilder) ([]subnets.Subnet, error)
	CreateSubnet(opts subnets.CreateOptsBuilder) (*subnets.Subnet, error)
	DeleteSubnet(id string) error
	GetSubnet(id string) (*subnets.Subnet, error)
	UpdateSubnet(id string, opts subnets.UpdateOptsBuilder) (*subnets.Subnet, error)

	ListExtensions() ([]extensions.Extension, error)

	ReplaceAllAttributesTags(resourceType string, resourceID string, opts attributestags.ReplaceAllOptsBuilder) ([]string, error)
}

type networkClient struct {
	serviceClient *gophercloud.ServiceClient
}

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

func (c networkClient) AddRouterInterface(id string, opts routers.AddInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	mc := metrics.NewMetricPrometheusContext("server_os_interface", "create")
	interfaceInfo, err := routers.AddInterface(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return interfaceInfo, nil
}

func (c networkClient) RemoveRouterInterface(id string, opts routers.RemoveInterfaceOptsBuilder) (*routers.InterfaceInfo, error) {
	mc := metrics.NewMetricPrometheusContext("server_os_interface", "delete")
	interfaceInfo, err := routers.RemoveInterface(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return interfaceInfo, nil
}

func (c networkClient) ReplaceAllAttributesTags(resourceType string, resourceID string, opts attributestags.ReplaceAllOptsBuilder) ([]string, error) {
	mc := metrics.NewMetricPrometheusContext("attributes_tags", "replace_all")
	tags, err := attributestags.ReplaceAll(c.serviceClient, resourceType, resourceID, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return tags, nil
}

func (c networkClient) ListRouter(opts routers.ListOpts) ([]routers.Router, error) {
	mc := metrics.NewMetricPrometheusContext("router", "list")
	allPages, err := routers.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return routers.ExtractRouters(allPages)
}

func (c networkClient) ListFloatingIP(opts floatingips.ListOptsBuilder) ([]floatingips.FloatingIP, error) {
	mc := metrics.NewMetricPrometheusContext("floating_ip", "list")
	allPages, err := floatingips.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return floatingips.ExtractFloatingIPs(allPages)
}

func (c networkClient) CreateFloatingIP(opts floatingips.CreateOptsBuilder) (*floatingips.FloatingIP, error) {
	mc := metrics.NewMetricPrometheusContext("floating_ip", "create")
	fip, err := floatingips.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return fip, nil
}

func (c networkClient) DeleteFloatingIP(id string) error {
	mc := metrics.NewMetricPrometheusContext("floating_ip", "delete")
	return mc.ObserveRequestIgnoreNotFound(floatingips.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetFloatingIP(id string) (*floatingips.FloatingIP, error) {
	mc := metrics.NewMetricPrometheusContext("floating_ip", "list")
	fip, err := floatingips.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return fip, nil
}

func (c networkClient) UpdateFloatingIP(id string, opts floatingips.UpdateOptsBuilder) (*floatingips.FloatingIP, error) {
	mc := metrics.NewMetricPrometheusContext("floating_ip", "update")
	fip, err := floatingips.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return fip, nil
}

func (c networkClient) ListPort(opts ports.ListOptsBuilder) ([]ports.Port, error) {
	mc := metrics.NewMetricPrometheusContext("port", "list")
	allPages, err := ports.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return ports.ExtractPorts(allPages)
}

func (c networkClient) CreatePort(opts ports.CreateOptsBuilder) (*ports.Port, error) {
	mc := metrics.NewMetricPrometheusContext("port", "create")
	port, err := ports.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return port, nil
}

func (c networkClient) DeletePort(id string) error {
	mc := metrics.NewMetricPrometheusContext("port", "delete")
	return mc.ObserveRequestIgnoreNotFound(ports.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetPort(id string) (*ports.Port, error) {
	mc := metrics.NewMetricPrometheusContext("port", "get")
	port, err := ports.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return port, nil
}

func (c networkClient) UpdatePort(id string, opts ports.UpdateOptsBuilder) (*ports.Port, error) {
	mc := metrics.NewMetricPrometheusContext("port", "update")
	port, err := ports.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return port, nil
}

func (c networkClient) CreateTrunk(opts trunks.CreateOptsBuilder) (*trunks.Trunk, error) {
	mc := metrics.NewMetricPrometheusContext("trunk", "create")
	trunk, err := trunks.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return trunk, nil
}

func (c networkClient) DeleteTrunk(id string) error {
	mc := metrics.NewMetricPrometheusContext("trunk", "delete")
	return mc.ObserveRequestIgnoreNotFound(trunks.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) ListTrunk(opts trunks.ListOptsBuilder) ([]trunks.Trunk, error) {
	mc := metrics.NewMetricPrometheusContext("trunk", "list")
	allPages, err := trunks.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return trunks.ExtractTrunks(allPages)
}

func (c networkClient) CreateRouter(opts routers.CreateOptsBuilder) (*routers.Router, error) {
	mc := metrics.NewMetricPrometheusContext("router", "create")
	router, err := routers.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return router, nil
}

func (c networkClient) DeleteRouter(id string) error {
	mc := metrics.NewMetricPrometheusContext("router", "delete")
	return mc.ObserveRequestIgnoreNotFound(routers.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetRouter(id string) (*routers.Router, error) {
	mc := metrics.NewMetricPrometheusContext("router", "get")
	router, err := routers.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return router, nil
}

func (c networkClient) UpdateRouter(id string, opts routers.UpdateOptsBuilder) (*routers.Router, error) {
	mc := metrics.NewMetricPrometheusContext("router", "update")
	router, err := routers.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return router, nil
}

func (c networkClient) ListSecGroup(opts groups.ListOpts) ([]groups.SecGroup, error) {
	mc := metrics.NewMetricPrometheusContext("group", "list")
	allPages, err := groups.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return groups.ExtractGroups(allPages)
}

func (c networkClient) CreateSecGroup(opts groups.CreateOptsBuilder) (*groups.SecGroup, error) {
	mc := metrics.NewMetricPrometheusContext("security_group", "create")
	group, err := groups.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return group, nil
}

func (c networkClient) DeleteSecGroup(id string) error {
	mc := metrics.NewMetricPrometheusContext("security_group", "delete")
	return mc.ObserveRequestIgnoreNotFound(groups.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetSecGroup(id string) (*groups.SecGroup, error) {
	mc := metrics.NewMetricPrometheusContext("security_group", "get")
	group, err := groups.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return group, nil
}

func (c networkClient) UpdateSecGroup(id string, opts groups.UpdateOptsBuilder) (*groups.SecGroup, error) {
	mc := metrics.NewMetricPrometheusContext("security_group", "update")
	group, err := groups.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return group, nil
}

func (c networkClient) ListSecGroupRule(opts rules.ListOpts) ([]rules.SecGroupRule, error) {
	mc := metrics.NewMetricPrometheusContext("security_group_rule", "list")
	allPages, err := rules.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return rules.ExtractRules(allPages)
}

func (c networkClient) CreateSecGroupRule(opts rules.CreateOptsBuilder) (*rules.SecGroupRule, error) {
	mc := metrics.NewMetricPrometheusContext("security_group_rule", "create")
	rule, err := rules.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return rule, nil
}

func (c networkClient) DeleteSecGroupRule(id string) error {
	mc := metrics.NewMetricPrometheusContext("security_group_rule", "delete")
	return mc.ObserveRequestIgnoreNotFound(rules.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetSecGroupRule(id string) (*rules.SecGroupRule, error) {
	mc := metrics.NewMetricPrometheusContext("security_group_rule", "get")
	rule, err := rules.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return rule, nil
}

func (c networkClient) ListNetwork(opts networks.ListOptsBuilder) ([]networks.Network, error) {
	mc := metrics.NewMetricPrometheusContext("network", "list")
	allPages, err := networks.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return networks.ExtractNetworks(allPages)
}

func (c networkClient) CreateNetwork(opts networks.CreateOptsBuilder) (*networks.Network, error) {
	mc := metrics.NewMetricPrometheusContext("network", "create")
	net, err := networks.Create(c.serviceClient, opts).Extract()
	if (mc.ObserveRequest(err)) != nil {
		return nil, err
	}
	return net, nil
}

func (c networkClient) DeleteNetwork(id string) error {
	mc := metrics.NewMetricPrometheusContext("network", "delete")
	return mc.ObserveRequestIgnoreNotFound(networks.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetNetwork(id string) (*networks.Network, error) {
	mc := metrics.NewMetricPrometheusContext("network", "get")
	net, err := networks.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return net, nil
}

func (c networkClient) UpdateNetwork(id string, opts networks.UpdateOptsBuilder) (*networks.Network, error) {
	mc := metrics.NewMetricPrometheusContext("network", "update")
	net, err := networks.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return net, nil
}

func (c networkClient) ListSubnet(opts subnets.ListOptsBuilder) ([]subnets.Subnet, error) {
	mc := metrics.NewMetricPrometheusContext("subnet", "list")
	allPages, err := subnets.List(c.serviceClient, opts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return subnets.ExtractSubnets(allPages)
}

func (c networkClient) CreateSubnet(opts subnets.CreateOptsBuilder) (*subnets.Subnet, error) {
	mc := metrics.NewMetricPrometheusContext("subnet", "create")
	subnet, err := subnets.Create(c.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return subnet, nil
}

func (c networkClient) DeleteSubnet(id string) error {
	mc := metrics.NewMetricPrometheusContext("subnet", "delete")
	return mc.ObserveRequestIgnoreNotFound(subnets.Delete(c.serviceClient, id).ExtractErr())
}

func (c networkClient) GetSubnet(id string) (*subnets.Subnet, error) {
	mc := metrics.NewMetricPrometheusContext("subnet", "get")
	subnet, err := subnets.Get(c.serviceClient, id).Extract()
	if mc.ObserveRequestIgnoreNotFound(err) != nil {
		return nil, err
	}
	return subnet, nil
}

func (c networkClient) UpdateSubnet(id string, opts subnets.UpdateOptsBuilder) (*subnets.Subnet, error) {
	mc := metrics.NewMetricPrometheusContext("subnet", "update")
	subnet, err := subnets.Update(c.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return subnet, nil
}

func (c networkClient) ListExtensions() ([]extensions.Extension, error) {
	mc := metrics.NewMetricPrometheusContext("network_extension", "list")
	allPages, err := extensions.List(c.serviceClient).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return extensions.ExtractExtensions(allPages)
}
