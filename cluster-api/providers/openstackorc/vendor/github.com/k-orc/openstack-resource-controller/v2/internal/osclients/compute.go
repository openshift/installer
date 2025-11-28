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
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/tags"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

/*
NovaMinimumMicroversion is the minimum Nova microversion supported by CAPO and ORC.
2.71 corresponds to OpenStack Stein

For the canonical description of Nova microversions, see
https://docs.openstack.org/nova/latest/reference/api-microversion-history.html

CAPO uses server tags, which were added in microversion 2.52.
CAPO supports multiattach volume types, which were added in microversion 2.60.
ORC supports server groups by specifying Policy/Rules instead of Policies, which were added in microversion 2.64.
ORC requires the use of microversion 2.71 to support ServerGroups field in Server response.
*/
const NovaMinimumMicroversion = "2.71"

type ComputeClient interface {
	CreateFlavor(ctx context.Context, opts flavors.CreateOptsBuilder) (*flavors.Flavor, error)
	GetFlavor(ctx context.Context, id string) (*flavors.Flavor, error)
	DeleteFlavor(ctx context.Context, id string) error
	ListFlavors(ctx context.Context, listOpts flavors.ListOptsBuilder) iter.Seq2[*flavors.Flavor, error]

	CreateServer(ctx context.Context, createOpts servers.CreateOptsBuilder, schedulerHints servers.SchedulerHintOptsBuilder) (*servers.Server, error)
	DeleteServer(ctx context.Context, serverID string) error
	GetServer(ctx context.Context, serverID string) (*servers.Server, error)
	ListServers(ctx context.Context, listOpts servers.ListOptsBuilder) iter.Seq2[*servers.Server, error]
	UpdateServer(ctx context.Context, id string, opts servers.UpdateOptsBuilder) (*servers.Server, error)

	CreateServerGroup(ctx context.Context, createOpts servergroups.CreateOptsBuilder) (*servergroups.ServerGroup, error)
	DeleteServerGroup(ctx context.Context, serverGroupID string) error
	GetServerGroup(ctx context.Context, serverGroupID string) (*servergroups.ServerGroup, error)
	ListServerGroups(ctx context.Context, listOpts servergroups.ListOptsBuilder) iter.Seq2[*servergroups.ServerGroup, error]

	CreateVolumeAttachment(ctx context.Context, serverID string, createOpts volumeattach.CreateOptsBuilder) (*volumeattach.VolumeAttachment, error)
	DeleteVolumeAttachment(ctx context.Context, serverID, volumeID string) error

	ListAttachedInterfaces(ctx context.Context, serverID string) ([]attachinterfaces.Interface, error)
	CreateAttachedInterface(ctx context.Context, serverID string, createOpts attachinterfaces.CreateOptsBuilder) (*attachinterfaces.Interface, error)
	DeleteAttachedInterface(ctx context.Context, serverID, portID string) error

	ReplaceAllServerAttributesTags(ctx context.Context, resourceID string, opts tags.ReplaceAllOptsBuilder) ([]string, error)
}

type computeClient struct{ client *gophercloud.ServiceClient }

// NewComputeClient returns a new compute client.
func NewComputeClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ComputeClient, error) {
	compute, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service client: %v", err)
	}
	compute.Microversion = NovaMinimumMicroversion

	return &computeClient{compute}, nil
}

func (c computeClient) ListAvailabilityZones(ctx context.Context) ([]availabilityzones.AvailabilityZone, error) {
	allPages, err := availabilityzones.List(c.client).AllPages(ctx)
	if err != nil {
		return nil, err
	}
	return availabilityzones.ExtractAvailabilityZones(allPages)
}

func (c computeClient) GetFlavor(ctx context.Context, id string) (*flavors.Flavor, error) {
	return flavors.Get(ctx, c.client, id).Extract()
}

func (c computeClient) CreateFlavor(ctx context.Context, opts flavors.CreateOptsBuilder) (*flavors.Flavor, error) {
	return flavors.Create(ctx, c.client, opts).Extract()
}

func (c computeClient) DeleteFlavor(ctx context.Context, id string) error {
	return flavors.Delete(ctx, c.client, id).ExtractErr()
}

func (c computeClient) ListFlavors(ctx context.Context, opts flavors.ListOptsBuilder) iter.Seq2[*flavors.Flavor, error] {
	pager := flavors.ListDetail(c.client, opts)
	return func(yield func(*flavors.Flavor, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(flavors.ExtractFlavors, yield))
	}
}

func (c computeClient) CreateServer(ctx context.Context, createOpts servers.CreateOptsBuilder, schedulerHints servers.SchedulerHintOptsBuilder) (*servers.Server, error) {
	return servers.Create(ctx, c.client, createOpts, schedulerHints).Extract()
}

func (c computeClient) DeleteServer(ctx context.Context, serverID string) error {
	return servers.Delete(ctx, c.client, serverID).ExtractErr()
}

func (c computeClient) GetServer(ctx context.Context, serverID string) (*servers.Server, error) {
	return servers.Get(ctx, c.client, serverID).Extract()
}

func (c computeClient) ListServers(ctx context.Context, opts servers.ListOptsBuilder) iter.Seq2[*servers.Server, error] {
	pager := servers.List(c.client, opts)
	return func(yield func(*servers.Server, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(servers.ExtractServers, yield))
	}
}

func (c computeClient) UpdateServer(ctx context.Context, id string, opts servers.UpdateOptsBuilder) (*servers.Server, error) {
	return servers.Update(ctx, c.client, id, opts).Extract()
}

func (c computeClient) ListAttachedInterfaces(ctx context.Context, serverID string) ([]attachinterfaces.Interface, error) {
	interfaces, err := attachinterfaces.List(c.client, serverID).AllPages(ctx)
	if err != nil {
		return nil, err
	}
	return attachinterfaces.ExtractInterfaces(interfaces)
}

func (c computeClient) CreateAttachedInterface(ctx context.Context, serverID string, createOpts attachinterfaces.CreateOptsBuilder) (*attachinterfaces.Interface, error) {
	return attachinterfaces.Create(ctx, c.client, serverID, createOpts).Extract()
}

func (c computeClient) DeleteAttachedInterface(ctx context.Context, serverID, portID string) error {
	return attachinterfaces.Delete(ctx, c.client, serverID, portID).ExtractErr()
}

func (c computeClient) CreateServerGroup(ctx context.Context, createOpts servergroups.CreateOptsBuilder) (*servergroups.ServerGroup, error) {
	return servergroups.Create(ctx, c.client, createOpts).Extract()
}

func (c computeClient) DeleteServerGroup(ctx context.Context, serverGroupID string) error {
	return servergroups.Delete(ctx, c.client, serverGroupID).ExtractErr()
}

func (c computeClient) GetServerGroup(ctx context.Context, serverGroupID string) (*servergroups.ServerGroup, error) {
	return servergroups.Get(ctx, c.client, serverGroupID).Extract()
}

func (c computeClient) ListServerGroups(ctx context.Context, opts servergroups.ListOptsBuilder) iter.Seq2[*servergroups.ServerGroup, error] {
	pager := servergroups.List(c.client, opts)
	return func(yield func(*servergroups.ServerGroup, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(servergroups.ExtractServerGroups, yield))
	}
}

func (c computeClient) CreateVolumeAttachment(ctx context.Context, serverID string, createOpts volumeattach.CreateOptsBuilder) (*volumeattach.VolumeAttachment, error) {
	return volumeattach.Create(ctx, c.client, serverID, createOpts).Extract()
}

func (c computeClient) DeleteVolumeAttachment(ctx context.Context, serverID, volumeID string) error {
	return volumeattach.Delete(ctx, c.client, serverID, volumeID).ExtractErr()
}

func (c computeClient) ReplaceAllServerAttributesTags(ctx context.Context, resourceID string, opts tags.ReplaceAllOptsBuilder) ([]string, error) {
	return tags.ReplaceAll(ctx, c.client, resourceID, opts).Extract()
}

type computeErrorClient struct{ error }

// NewComputeErrorClient returns a ComputeClient in which every method returns the given error.
func NewComputeErrorClient(e error) ComputeClient {
	return computeErrorClient{e}
}
func (e computeErrorClient) CreateFlavor(ctx context.Context, opts flavors.CreateOptsBuilder) (*flavors.Flavor, error) {
	return nil, e.error
}
func (e computeErrorClient) GetFlavor(ctx context.Context, id string) (*flavors.Flavor, error) {
	return nil, e.error
}
func (e computeErrorClient) DeleteFlavor(ctx context.Context, id string) error {
	return e.error
}
func (e computeErrorClient) ListFlavors(_ context.Context, _ flavors.ListOptsBuilder) iter.Seq2[*flavors.Flavor, error] {
	return func(yield func(*flavors.Flavor, error) bool) {
		yield(nil, e.error)
	}
}

func (e computeErrorClient) ListAvailabilityZones(_ context.Context) ([]availabilityzones.AvailabilityZone, error) {
	return nil, e.error
}

func (e computeErrorClient) CreateServer(_ context.Context, _ servers.CreateOptsBuilder, _ servers.SchedulerHintOptsBuilder) (*servers.Server, error) {
	return nil, e.error
}

func (e computeErrorClient) DeleteServer(_ context.Context, _ string) error {
	return e.error
}

func (e computeErrorClient) GetServer(_ context.Context, _ string) (*servers.Server, error) {
	return nil, e.error
}

func (e computeErrorClient) ListServers(ctx context.Context, listOpts servers.ListOptsBuilder) iter.Seq2[*servers.Server, error] {
	return func(yield func(*servers.Server, error) bool) {
		yield(nil, e.error)
	}
}

func (e computeErrorClient) UpdateServer(_ context.Context, _ string, _ servers.UpdateOptsBuilder) (*servers.Server, error) {
	return nil, e.error
}

func (e computeErrorClient) CreateServerGroup(_ context.Context, _ servergroups.CreateOptsBuilder) (*servergroups.ServerGroup, error) {
	return nil, e.error
}

func (e computeErrorClient) DeleteServerGroup(_ context.Context, _ string) error {
	return e.error
}

func (e computeErrorClient) GetServerGroup(_ context.Context, _ string) (*servergroups.ServerGroup, error) {
	return nil, e.error
}

func (e computeErrorClient) ListServerGroups(ctx context.Context, listOpts servergroups.ListOptsBuilder) iter.Seq2[*servergroups.ServerGroup, error] {
	return func(yield func(*servergroups.ServerGroup, error) bool) {
		yield(nil, e.error)
	}
}

func (e computeErrorClient) CreateVolumeAttachment(_ context.Context, _ string, _ volumeattach.CreateOptsBuilder) (*volumeattach.VolumeAttachment, error) {
	return nil, e.error
}

func (e computeErrorClient) DeleteVolumeAttachment(_ context.Context, _, _ string) error {
	return e.error
}

func (e computeErrorClient) ListAttachedInterfaces(_ context.Context, _ string) ([]attachinterfaces.Interface, error) {
	return nil, e.error
}

func (e computeErrorClient) CreateAttachedInterface(_ context.Context, _ string, _ attachinterfaces.CreateOptsBuilder) (*attachinterfaces.Interface, error) {
	return nil, e.error
}

func (e computeErrorClient) DeleteAttachedInterface(_ context.Context, _, _ string) error {
	return e.error
}

func (e computeErrorClient) ReplaceAllServerAttributesTags(_ context.Context, _ string, _ tags.ReplaceAllOptsBuilder) ([]string, error) {
	return nil, e.error
}
