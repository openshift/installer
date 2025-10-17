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

package routetables

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// RouteTableSpec defines the specification for a route table.
type RouteTableSpec struct {
	Name           string
	ResourceGroup  string
	Location       string
	ClusterName    string
	AdditionalTags infrav1.Tags
}

// ResourceName returns the name of the route table.
func (s *RouteTableSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *RouteTableSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for route tables.
func (s *RouteTableSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the route table.
func (s *RouteTableSpec) Parameters(_ context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armnetwork.RouteTable); !ok {
			return nil, errors.Errorf("%T is not an armnetwork.RouteTable", existing)
		}
		// route table already exists
		// currently don't support specifying your own routes via spec.
		return nil, nil
	}
	return armnetwork.RouteTable{
		Location:   ptr.To(s.Location),
		Properties: &armnetwork.RouteTablePropertiesFormat{},
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Additional:  s.AdditionalTags,
		})),
	}, nil
}
