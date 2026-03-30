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
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type VolumeClient interface {
	ListVolumes(ctx context.Context, listOpts volumes.ListOptsBuilder) iter.Seq2[*volumes.Volume, error]
	CreateVolume(ctx context.Context, opts volumes.CreateOptsBuilder) (*volumes.Volume, error)
	DeleteVolume(ctx context.Context, resourceID string, opts volumes.DeleteOptsBuilder) error
	GetVolume(ctx context.Context, resourceID string) (*volumes.Volume, error)
	UpdateVolume(ctx context.Context, id string, opts volumes.UpdateOptsBuilder) (*volumes.Volume, error)
}

type volumeClient struct{ client *gophercloud.ServiceClient }

// NewVolumeClient returns a new OpenStack client.
func NewVolumeClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (VolumeClient, error) {
	client, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create volume service client: %v", err)
	}

	return &volumeClient{client}, nil
}

func (c volumeClient) ListVolumes(ctx context.Context, listOpts volumes.ListOptsBuilder) iter.Seq2[*volumes.Volume, error] {
	pager := volumes.List(c.client, listOpts)
	return func(yield func(*volumes.Volume, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(volumes.ExtractVolumes, yield))
	}
}

func (c volumeClient) CreateVolume(ctx context.Context, opts volumes.CreateOptsBuilder) (*volumes.Volume, error) {
	return volumes.Create(ctx, c.client, opts, nil).Extract()
}

func (c volumeClient) DeleteVolume(ctx context.Context, resourceID string, opts volumes.DeleteOptsBuilder) error {
	return volumes.Delete(ctx, c.client, resourceID, opts).ExtractErr()
}

func (c volumeClient) GetVolume(ctx context.Context, resourceID string) (*volumes.Volume, error) {
	return volumes.Get(ctx, c.client, resourceID).Extract()
}

func (c volumeClient) UpdateVolume(ctx context.Context, id string, opts volumes.UpdateOptsBuilder) (*volumes.Volume, error) {
	return volumes.Update(ctx, c.client, id, opts).Extract()
}

type volumeErrorClient struct{ error }

// NewVolumeErrorClient returns a VolumeClient in which every method returns the given error.
func NewVolumeErrorClient(e error) VolumeClient {
	return volumeErrorClient{e}
}

func (e volumeErrorClient) ListVolumes(_ context.Context, _ volumes.ListOptsBuilder) iter.Seq2[*volumes.Volume, error] {
	return func(yield func(*volumes.Volume, error) bool) {
		yield(nil, e.error)
	}
}

func (e volumeErrorClient) CreateVolume(_ context.Context, _ volumes.CreateOptsBuilder) (*volumes.Volume, error) {
	return nil, e.error
}

func (e volumeErrorClient) DeleteVolume(_ context.Context, _ string, _ volumes.DeleteOptsBuilder) error {
	return e.error
}

func (e volumeErrorClient) GetVolume(_ context.Context, _ string) (*volumes.Volume, error) {
	return nil, e.error
}

func (e volumeErrorClient) UpdateVolume(_ context.Context, _ string, _ volumes.UpdateOptsBuilder) (*volumes.Volume, error) {
	return nil, e.error
}
