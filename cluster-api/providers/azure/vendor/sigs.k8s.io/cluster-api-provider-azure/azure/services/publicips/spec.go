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

package publicips

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// PublicIPSpec defines the specification for a Public IP.
type PublicIPSpec struct {
	Name             string
	ResourceGroup    string
	ClusterName      string
	DNSName          string
	IsIPv6           bool
	Location         string
	ExtendedLocation *infrav1.ExtendedLocationSpec
	FailureDomains   []*string
	AdditionalTags   infrav1.Tags
	IPTags           []infrav1.IPTag
}

// ResourceName returns the name of the public IP.
func (s *PublicIPSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *PublicIPSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for public IPs.
func (s *PublicIPSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the public IP.
func (s *PublicIPSpec) Parameters(_ context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armnetwork.PublicIPAddress); !ok {
			return nil, errors.Errorf("%T is not an armnetwork.PublicIPAddress", existing)
		}
		// public IP already exists
		return nil, nil
	}

	addressVersion := armnetwork.IPVersionIPv4
	if s.IsIPv6 {
		addressVersion = armnetwork.IPVersionIPv6
	}

	// only set DNS properties if there is a DNS name specified
	var dnsSettings *armnetwork.PublicIPAddressDNSSettings
	if s.DNSName != "" {
		dnsSettings = &armnetwork.PublicIPAddressDNSSettings{
			DomainNameLabel: ptr.To(strings.Split(s.DNSName, ".")[0]),
			Fqdn:            ptr.To(s.DNSName),
		}
	}

	return armnetwork.PublicIPAddress{
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Additional:  s.AdditionalTags,
		})),
		SKU:              &armnetwork.PublicIPAddressSKU{Name: ptr.To(armnetwork.PublicIPAddressSKUNameStandard)},
		Name:             ptr.To(s.Name),
		Location:         ptr.To(s.Location),
		ExtendedLocation: converters.ExtendedLocationToNetworkSDK(s.ExtendedLocation),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   &addressVersion,
			PublicIPAllocationMethod: ptr.To(armnetwork.IPAllocationMethodStatic),
			DNSSettings:              dnsSettings,
			IPTags:                   converters.IPTagsToSDK(s.IPTags),
		},
		Zones: s.FailureDomains,
	}, nil
}
