/*
Copyright 2025 The ORC Authors.

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
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type VolumeTypeClient interface {
	ListVolumeTypes(ctx context.Context, listOpts volumetypes.ListOptsBuilder) iter.Seq2[*volumetypes.VolumeType, error]
	CreateVolumeType(ctx context.Context, opts volumetypes.CreateOptsBuilder) (*volumetypes.VolumeType, error)
	DeleteVolumeType(ctx context.Context, resourceID string) error
	GetVolumeType(ctx context.Context, resourceID string) (*volumetypes.VolumeType, error)
	UpdateVolumeType(ctx context.Context, id string, opts volumetypes.UpdateOptsBuilder) (*volumetypes.VolumeType, error)
}

type volumetypeClient struct{ client *gophercloud.ServiceClient }

// NewVolumeTypeClient returns a new OpenStack client.
func NewVolumeTypeClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (VolumeTypeClient, error) {
	client, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create volumetype service client: %v", err)
	}

	return &volumetypeClient{client}, nil
}

func (c volumetypeClient) ListVolumeTypes(ctx context.Context, listOpts volumetypes.ListOptsBuilder) iter.Seq2[*volumetypes.VolumeType, error] {
	pager := volumetypes.List(c.client, listOpts)
	return func(yield func(*volumetypes.VolumeType, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(volumetypes.ExtractVolumeTypes, yield))
	}
}

func (c volumetypeClient) CreateVolumeType(ctx context.Context, opts volumetypes.CreateOptsBuilder) (*volumetypes.VolumeType, error) {
	return volumetypes.Create(ctx, c.client, opts).Extract()
}

func (c volumetypeClient) DeleteVolumeType(ctx context.Context, resourceID string) error {
	return volumetypes.Delete(ctx, c.client, resourceID).ExtractErr()
}

func (c volumetypeClient) GetVolumeType(ctx context.Context, resourceID string) (*volumetypes.VolumeType, error) {
	return volumetypes.Get(ctx, c.client, resourceID).Extract()
}

func (c volumetypeClient) UpdateVolumeType(ctx context.Context, id string, opts volumetypes.UpdateOptsBuilder) (*volumetypes.VolumeType, error) {
	return volumetypes.Update(ctx, c.client, id, opts).Extract()
}

type volumetypeErrorClient struct{ error }

// NewVolumeTypeErrorClient returns a VolumeTypeClient in which every method returns the given error.
func NewVolumeTypeErrorClient(e error) VolumeTypeClient {
	return volumetypeErrorClient{e}
}

func (e volumetypeErrorClient) ListVolumeTypes(_ context.Context, _ volumetypes.ListOptsBuilder) iter.Seq2[*volumetypes.VolumeType, error] {
	return func(yield func(*volumetypes.VolumeType, error) bool) {
		yield(nil, e.error)
	}
}

func (e volumetypeErrorClient) CreateVolumeType(_ context.Context, _ volumetypes.CreateOptsBuilder) (*volumetypes.VolumeType, error) {
	return nil, e.error
}

func (e volumetypeErrorClient) DeleteVolumeType(_ context.Context, _ string) error {
	return e.error
}

func (e volumetypeErrorClient) GetVolumeType(_ context.Context, _ string) (*volumetypes.VolumeType, error) {
	return nil, e.error
}

func (e volumetypeErrorClient) UpdateVolumeType(_ context.Context, _ string, _ volumetypes.UpdateOptsBuilder) (*volumetypes.VolumeType, error) {
	return nil, e.error
}
