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

package availabilitysets

import (
	"context"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
)

// AvailabilitySetSpec defines the specification for an availability set.
type AvailabilitySetSpec struct {
	Name             string
	ResourceGroup    string
	ClusterName      string
	Location         string
	CloudEnvironment string
	SKU              *resourceskus.SKU
	AdditionalTags   infrav1.Tags
}

// ResourceName returns the name of the availability set.
func (s *AvailabilitySetSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *AvailabilitySetSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for availability sets.
func (s *AvailabilitySetSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the availability set.
func (s *AvailabilitySetSpec) Parameters(_ context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armcompute.AvailabilitySet); !ok {
			return nil, errors.Errorf("%T is not an armcompute.AvailabilitySet", existing)
		}
		// availability set already exists
		return nil, nil
	}

	faultDomainCount, err := getFaultDomainCount(s.SKU, s.CloudEnvironment)
	if err != nil {
		return nil, err
	}

	asParams := armcompute.AvailabilitySet{
		SKU: &armcompute.SKU{
			Name: ptr.To(string(armcompute.AvailabilitySetSKUTypesAligned)),
		},
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount: faultDomainCount,
		},
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Role:        ptr.To(infrav1.CommonRole),
			Additional:  s.AdditionalTags,
		})),
		Location: ptr.To(s.Location),
	}

	return asParams, nil
}

func getFaultDomainCount(SKU *resourceskus.SKU, cloudEnvironment string) (*int32, error) {
	// Azure Stack environments may not implement the resource SKU API
	// for availability sets. Use a default value instead.
	if strings.EqualFold(cloudEnvironment, azure.StackCloudName) {
		return ptr.To(int32(2)), nil
	}

	if SKU == nil {
		return nil, errors.New("unable to get required availability set SKU from machine cache")
	}

	var faultDomainCount *int32
	faultDomainCountStr, ok := SKU.GetCapability(resourceskus.MaximumPlatformFaultDomainCount)
	if !ok {
		return nil, errors.Errorf("unable to get required availability set SKU capability %s", resourceskus.MaximumPlatformFaultDomainCount)
	}
	count, err := strconv.ParseInt(faultDomainCountStr, 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse availability set fault domain count")
	}
	faultDomainCount = ptr.To[int32](int32(count))
	return faultDomainCount, nil
}
