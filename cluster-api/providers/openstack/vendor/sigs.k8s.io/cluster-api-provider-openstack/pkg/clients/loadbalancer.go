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

package clients

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/apiversions"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/providers"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
)

type LbClient interface {
	CreateLoadBalancer(opts loadbalancers.CreateOptsBuilder) (*loadbalancers.LoadBalancer, error)
	ListLoadBalancers(opts loadbalancers.ListOptsBuilder) ([]loadbalancers.LoadBalancer, error)
	GetLoadBalancer(id string) (*loadbalancers.LoadBalancer, error)
	DeleteLoadBalancer(id string, opts loadbalancers.DeleteOptsBuilder) error
	CreateListener(opts listeners.CreateOptsBuilder) (*listeners.Listener, error)
	ListListeners(opts listeners.ListOptsBuilder) ([]listeners.Listener, error)
	UpdateListener(id string, opts listeners.UpdateOpts) (*listeners.Listener, error)
	GetListener(id string) (*listeners.Listener, error)
	DeleteListener(id string) error
	CreatePool(opts pools.CreateOptsBuilder) (*pools.Pool, error)
	ListPools(opts pools.ListOptsBuilder) ([]pools.Pool, error)
	GetPool(id string) (*pools.Pool, error)
	DeletePool(id string) error
	CreatePoolMember(poolID string, opts pools.CreateMemberOptsBuilder) (*pools.Member, error)
	ListPoolMember(poolID string, opts pools.ListMembersOptsBuilder) ([]pools.Member, error)
	DeletePoolMember(poolID string, lbMemberID string) error
	CreateMonitor(opts monitors.CreateOptsBuilder) (*monitors.Monitor, error)
	ListMonitors(opts monitors.ListOptsBuilder) ([]monitors.Monitor, error)
	UpdateMonitor(id string, opts monitors.UpdateOptsBuilder) (*monitors.Monitor, error)
	DeleteMonitor(id string) error
	ListLoadBalancerProviders() ([]providers.Provider, error)
	ListOctaviaVersions() ([]apiversions.APIVersion, error)
	ListLoadBalancerFlavors() ([]flavors.Flavor, error)
}

type lbClient struct {
	serviceClient *gophercloud.ServiceClient
}

// NewLbClient returns a new loadbalancer client.
func NewLbClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (LbClient, error) {
	loadbalancerClient, err := openstack.NewLoadBalancerV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer service client: %v", err)
	}

	return &lbClient{loadbalancerClient}, nil
}

func (l lbClient) CreateLoadBalancer(opts loadbalancers.CreateOptsBuilder) (*loadbalancers.LoadBalancer, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer", "create")
	lb, err := loadbalancers.Create(context.TODO(), l.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return lb, nil
}

func (l lbClient) ListLoadBalancers(opts loadbalancers.ListOptsBuilder) ([]loadbalancers.LoadBalancer, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer", "list")
	allPages, err := loadbalancers.List(l.serviceClient, opts).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return loadbalancers.ExtractLoadBalancers(allPages)
}

func (l lbClient) GetLoadBalancer(id string) (*loadbalancers.LoadBalancer, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer", "get")
	lb, err := loadbalancers.Get(context.TODO(), l.serviceClient, id).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return lb, nil
}

func (l lbClient) DeleteLoadBalancer(id string, opts loadbalancers.DeleteOptsBuilder) error {
	mc := metrics.NewMetricPrometheusContext("loadbalancer", "delete")
	err := loadbalancers.Delete(context.TODO(), l.serviceClient, id, opts).ExtractErr()
	if mc.ObserveRequestIgnoreNotFound(err) != nil && !capoerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (l lbClient) CreateListener(opts listeners.CreateOptsBuilder) (*listeners.Listener, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_listener", "create")
	listener, err := listeners.Create(context.TODO(), l.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return listener, nil
}

func (l lbClient) UpdateListener(id string, opts listeners.UpdateOpts) (*listeners.Listener, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_listener", "update")
	listener, err := listeners.Update(context.TODO(), l.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return listener, nil
}

func (l lbClient) ListListeners(opts listeners.ListOptsBuilder) ([]listeners.Listener, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_listener", "list")
	allPages, err := listeners.List(l.serviceClient, opts).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return listeners.ExtractListeners(allPages)
}

func (l lbClient) GetListener(id string) (*listeners.Listener, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_listener", "get")
	listener, err := listeners.Get(context.TODO(), l.serviceClient, id).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return listener, nil
}

func (l lbClient) DeleteListener(id string) error {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_listener", "delete")
	err := listeners.Delete(context.TODO(), l.serviceClient, id).ExtractErr()
	if mc.ObserveRequestIgnoreNotFound(err) != nil && !capoerrors.IsNotFound(err) {
		return fmt.Errorf("error deleting lbaas listener %s: %v", id, err)
	}
	return nil
}

func (l lbClient) CreatePool(opts pools.CreateOptsBuilder) (*pools.Pool, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_pool", "create")
	pool, err := pools.Create(context.TODO(), l.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return pool, nil
}

func (l lbClient) ListPools(opts pools.ListOptsBuilder) ([]pools.Pool, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_pool", "list")
	allPages, err := pools.List(l.serviceClient, opts).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return pools.ExtractPools(allPages)
}

func (l lbClient) GetPool(id string) (*pools.Pool, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_pool", "get")
	pool, err := pools.Get(context.TODO(), l.serviceClient, id).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return pool, nil
}

func (l lbClient) DeletePool(id string) error {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_pool", "delete")
	err := pools.Delete(context.TODO(), l.serviceClient, id).ExtractErr()
	if mc.ObserveRequestIgnoreNotFound(err) != nil && !capoerrors.IsNotFound(err) {
		return fmt.Errorf("error deleting lbaas pool %s: %v", id, err)
	}
	return nil
}

func (l lbClient) CreatePoolMember(poolID string, lbMemberOpts pools.CreateMemberOptsBuilder) (*pools.Member, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_member", "create")
	member, err := pools.CreateMember(context.TODO(), l.serviceClient, poolID, lbMemberOpts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, fmt.Errorf("error create lbmember: %s", err)
	}
	return member, nil
}

func (l lbClient) ListPoolMember(poolID string, opts pools.ListMembersOptsBuilder) ([]pools.Member, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_pool", "list")
	allPages, err := pools.ListMembers(l.serviceClient, poolID, opts).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return pools.ExtractMembers(allPages)
}

func (l lbClient) DeletePoolMember(poolID string, lbMemberID string) error {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_member", "delete")
	err := pools.DeleteMember(context.TODO(), l.serviceClient, poolID, lbMemberID).ExtractErr()
	if mc.ObserveRequest(err) != nil {
		return fmt.Errorf("error deleting lbmember: %s", err)
	}
	return nil
}

func (l lbClient) CreateMonitor(opts monitors.CreateOptsBuilder) (*monitors.Monitor, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_healthmonitor", "create")
	monitor, err := monitors.Create(context.TODO(), l.serviceClient, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return monitor, nil
}

func (l lbClient) ListMonitors(opts monitors.ListOptsBuilder) ([]monitors.Monitor, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_healthmonitor", "list")
	allPages, err := monitors.List(l.serviceClient, opts).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return monitors.ExtractMonitors(allPages)
}

func (l lbClient) UpdateMonitor(id string, opts monitors.UpdateOptsBuilder) (*monitors.Monitor, error) {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_healthmonitor", "update")
	monitor, err := monitors.Update(context.TODO(), l.serviceClient, id, opts).Extract()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return monitor, nil
}

func (l lbClient) DeleteMonitor(id string) error {
	mc := metrics.NewMetricPrometheusContext("loadbalancer_healthmonitor", "delete")
	err := monitors.Delete(context.TODO(), l.serviceClient, id).ExtractErr()
	if mc.ObserveRequestIgnoreNotFound(err) != nil && !capoerrors.IsNotFound(err) {
		return fmt.Errorf("error deleting lbaas monitor %s: %v", id, err)
	}
	return nil
}

func (l lbClient) ListLoadBalancerProviders() ([]providers.Provider, error) {
	allPages, err := providers.List(l.serviceClient, providers.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("listing providers: %v", err)
	}
	providersList, err := providers.ExtractProviders(allPages)
	if err != nil {
		return nil, fmt.Errorf("extracting loadbalancer providers pages: %v", err)
	}
	return providersList, nil
}

func (l lbClient) ListOctaviaVersions() ([]apiversions.APIVersion, error) {
	mc := metrics.NewMetricPrometheusContext("version", "list")
	allPages, err := apiversions.List(l.serviceClient).AllPages(context.TODO())
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return apiversions.ExtractAPIVersions(allPages)
}

func (l lbClient) ListLoadBalancerFlavors() ([]flavors.Flavor, error) {
	allPages, err := flavors.List(l.serviceClient, flavors.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("listing flavors: %v", err)
	}
	flavorList, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return nil, fmt.Errorf("extracting loadbalancer flavors pages: %v", err)
	}
	return flavorList, nil
}
