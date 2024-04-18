/*
Copyright 2023 The Kubernetes Authors.

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

package converters

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// ExtendedLocationToNetworkSDK converts an infrav1.ExtendedLocationSpec to an armnetwork.ExtendedLocation.
func ExtendedLocationToNetworkSDK(src *infrav1.ExtendedLocationSpec) *armnetwork.ExtendedLocation {
	if src == nil {
		return nil
	}
	return &armnetwork.ExtendedLocation{
		Name: ptr.To(src.Name),
		Type: ptr.To(armnetwork.ExtendedLocationTypes(src.Type)),
	}
}

// ExtendedLocationToNetworkASO converts an infrav1.ExtendedLocationSpec to an asonetworkv1.ExtendedLocation.
func ExtendedLocationToNetworkASO(src *infrav1.ExtendedLocationSpec) *asonetworkv1.ExtendedLocation {
	if src == nil {
		return nil
	}
	return &asonetworkv1.ExtendedLocation{
		Name: ptr.To(src.Name),
		Type: ptr.To(asonetworkv1.ExtendedLocationType(src.Type)),
	}
}

// ExtendedLocationToComputeSDK converts an infrav1.ExtendedLocationSpec to an armcompute.ExtendedLocation.
func ExtendedLocationToComputeSDK(src *infrav1.ExtendedLocationSpec) *armcompute.ExtendedLocation {
	if src == nil {
		return nil
	}
	return &armcompute.ExtendedLocation{
		Name: ptr.To(src.Name),
		Type: ptr.To(armcompute.ExtendedLocationTypes(src.Type)),
	}
}
