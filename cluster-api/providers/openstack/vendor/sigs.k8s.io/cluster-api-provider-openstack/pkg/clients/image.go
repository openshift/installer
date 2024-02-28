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
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
)

type ImageClient interface {
	ListImages(listOpts images.ListOptsBuilder) ([]images.Image, error)
}

type imageClient struct{ client *gophercloud.ServiceClient }

// NewImageClient returns a new glance client.
func NewImageClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ImageClient, error) {
	images, err := openstack.NewImageServiceV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create image service client: %v", err)
	}

	return imageClient{images}, nil
}

func (c imageClient) ListImages(listOpts images.ListOptsBuilder) ([]images.Image, error) {
	mc := metrics.NewMetricPrometheusContext("image", "list")
	pages, err := images.List(c.client, listOpts).AllPages()
	if mc.ObserveRequest(err) != nil {
		return nil, err
	}
	return images.ExtractImages(pages)
}

type imageErrorClient struct{ error }

// NewImageErrorClient returns an ImageClient in which every method returns the given error.
func NewImageErrorClient(e error) ImageClient {
	return imageErrorClient{e}
}

func (e imageErrorClient) ListImages(_ images.ListOptsBuilder) ([]images.Image, error) {
	return nil, e.error
}
