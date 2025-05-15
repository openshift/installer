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

package privatedns

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// ZoneSpec defines the specification for private dns zone.
type ZoneSpec struct {
	Name           string
	ResourceGroup  string
	ClusterName    string
	AdditionalTags infrav1.Tags
}

// ResourceName returns the name of the private dns zone.
func (s ZoneSpec) ResourceName() string {
	return s.Name
}

// OwnerResourceName is a no-op for private dns zone.
func (s ZoneSpec) OwnerResourceName() string {
	return ""
}

// ResourceGroupName returns the name of the resource group of the private dns zone.
func (s ZoneSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// Parameters returns the parameters for the private dns zone.
func (s ZoneSpec) Parameters(_ context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		_, ok := existing.(armprivatedns.PrivateZone)
		if !ok {
			return nil, errors.Errorf("%T is not an armprivatedns.PrivateZone", existing)
		}
		return nil, nil
	}

	return armprivatedns.PrivateZone{
		Location: ptr.To(azure.Global),
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Additional:  s.AdditionalTags,
		})),
	}, nil
}
