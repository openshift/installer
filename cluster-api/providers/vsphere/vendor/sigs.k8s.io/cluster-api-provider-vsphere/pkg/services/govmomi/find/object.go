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

package find

import (
	"context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/cluster"
)

// ManagedRefFinder is the method to find the reference of the type specified in the Failure Domain.
type ManagedRefFinder func(context.Context) ([]object.Reference, error)

func ObjectFunc(failureDomainType infrav1.FailureDomainType, topology infrav1.Topology, finder *find.Finder) ManagedRefFinder {
	var managedRefFunc ManagedRefFinder

	switch failureDomainType {
	case infrav1.ComputeClusterFailureDomain:
		managedRefFunc = func(ctx context.Context) ([]object.Reference, error) {
			computeResource, err := finder.ClusterComputeResource(ctx, *topology.ComputeCluster)
			if err != nil {
				return nil, err
			}
			return []object.Reference{computeResource.Reference()}, nil
		}
	case infrav1.DatacenterFailureDomain:
		managedRefFunc = func(ctx context.Context) ([]object.Reference, error) {
			dataCenter, err := finder.Datacenter(ctx, topology.Datacenter)
			if err != nil {
				return nil, err
			}
			return []object.Reference{dataCenter.Reference()}, nil
		}
	case infrav1.HostGroupFailureDomain:
		managedRefFunc = func(ctx context.Context) ([]object.Reference, error) {
			ccr, err := finder.ClusterComputeResource(ctx, *topology.ComputeCluster)
			if err != nil {
				return nil, err
			}
			return cluster.ListHostsFromGroup(ctx, ccr, topology.Hosts.HostGroupName)
		}
	}
	return managedRefFunc
}
