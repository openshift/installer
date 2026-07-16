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
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type KeyPairClient interface {
	ListKeyPairs(ctx context.Context, listOpts keypairs.ListOptsBuilder) iter.Seq2[*keypairs.KeyPair, error]
	CreateKeyPair(ctx context.Context, opts keypairs.CreateOptsBuilder) (*keypairs.KeyPair, error)
	DeleteKeyPair(ctx context.Context, name string) error
	GetKeyPair(ctx context.Context, name string) (*keypairs.KeyPair, error)
}

type keypairClient struct{ client *gophercloud.ServiceClient }

// NewKeyPairClient returns a new OpenStack client.
func NewKeyPairClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (KeyPairClient, error) {
	client, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create keypair service client: %v", err)
	}
	client.Microversion = NovaMinimumMicroversion

	return &keypairClient{client}, nil
}

func (c keypairClient) ListKeyPairs(ctx context.Context, listOpts keypairs.ListOptsBuilder) iter.Seq2[*keypairs.KeyPair, error] {
	pager := keypairs.List(c.client, listOpts)
	return func(yield func(*keypairs.KeyPair, error) bool) {
		_ = pager.EachPage(ctx, yieldPage(keypairs.ExtractKeyPairs, yield))
	}
}

func (c keypairClient) CreateKeyPair(ctx context.Context, opts keypairs.CreateOptsBuilder) (*keypairs.KeyPair, error) {
	return keypairs.Create(ctx, c.client, opts).Extract()
}

func (c keypairClient) DeleteKeyPair(ctx context.Context, name string) error {
	return keypairs.Delete(ctx, c.client, name, nil).ExtractErr()
}

func (c keypairClient) GetKeyPair(ctx context.Context, name string) (*keypairs.KeyPair, error) {
	return keypairs.Get(ctx, c.client, name, nil).Extract()
}

type keypairErrorClient struct{ error }

// NewKeyPairErrorClient returns a KeyPairClient in which every method returns the given error.
func NewKeyPairErrorClient(e error) KeyPairClient {
	return keypairErrorClient{e}
}

func (e keypairErrorClient) ListKeyPairs(_ context.Context, _ keypairs.ListOptsBuilder) iter.Seq2[*keypairs.KeyPair, error] {
	return func(yield func(*keypairs.KeyPair, error) bool) {
		yield(nil, e.error)
	}
}

func (e keypairErrorClient) CreateKeyPair(_ context.Context, _ keypairs.CreateOptsBuilder) (*keypairs.KeyPair, error) {
	return nil, e.error
}

func (e keypairErrorClient) DeleteKeyPair(_ context.Context, _ string) error {
	return e.error
}

func (e keypairErrorClient) GetKeyPair(_ context.Context, _ string) (*keypairs.KeyPair, error) {
	return nil, e.error
}
