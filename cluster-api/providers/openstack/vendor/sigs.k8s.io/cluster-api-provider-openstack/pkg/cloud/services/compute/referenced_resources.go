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

package compute

import (
	"fmt"
	"slices"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

// ResolveMachineSpec is responsible for populating a ResolvedMachineSpec from
// an OpenStackMachineSpec and any external dependencies. The result contains no
// external dependencies, and does not require any complex logic on creation.
// Note that we only set the fields in ResolvedMachineSpec that are not set yet. This is ok because:
// - OpenStackMachine is immutable, so we can't change the spec after the machine is created.
// - the bastion is mutable, but we delete the bastion when the spec changes, so the bastion status will be empty.
func ResolveMachineSpec(scope *scope.WithLogger, spec *infrav1.OpenStackMachineSpec, resolved *infrav1.ResolvedMachineSpec, clusterResourceName, baseName string, openStackCluster *infrav1.OpenStackCluster, managedSecurityGroup *string) (changed bool, err error) {
	changed = false

	computeService, err := NewService(scope)
	if err != nil {
		return changed, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return changed, err
	}

	// ServerGroup is optional, so we only need to resolve it if it's set in the spec
	if spec.ServerGroup != nil && resolved.ServerGroupID == "" {
		serverGroupID, err := computeService.GetServerGroupID(spec.ServerGroup)
		if err != nil {
			return changed, err
		}
		resolved.ServerGroupID = serverGroupID
		changed = true
	}

	// Image is required, so we need to resolve it if it's not set
	if resolved.ImageID == "" {
		imageID, err := computeService.GetImageID(spec.Image)
		if err != nil {
			return changed, err
		}
		resolved.ImageID = imageID
		changed = true
	}

	// ConstructPorts requires the cluster network to have been set. We only
	// call this from places where we know it should have been set, but the
	// cluster status is externally-provided data so we check it anyway.
	if openStackCluster.Status.Network == nil {
		return changed, fmt.Errorf("called ResolveMachineSpec with nil OpenStackCluster.Status.Network")
	}

	// Network resources are required in order to get ports options.
	if len(resolved.Ports) == 0 {
		defaultNetwork := openStackCluster.Status.Network
		portsOpts, err := networkingService.ConstructPorts(spec, clusterResourceName, baseName, defaultNetwork, managedSecurityGroup, InstanceTags(spec, openStackCluster))
		if err != nil {
			return changed, err
		}
		resolved.Ports = portsOpts
		changed = true
	}

	return changed, nil
}

// InstanceTags returns the tags that should be applied to an instance.
// The tags are a deduplicated combination of the tags specified in the
// OpenStackMachineSpec and the ones specified on the OpenStackCluster.
func InstanceTags(spec *infrav1.OpenStackMachineSpec, openStackCluster *infrav1.OpenStackCluster) []string {
	machineTags := slices.Concat(spec.Tags, openStackCluster.Spec.Tags)

	seen := make(map[string]struct{}, len(machineTags))
	unique := make([]string, 0, len(machineTags))
	for _, tag := range machineTags {
		if _, ok := seen[tag]; !ok {
			seen[tag] = struct{}{}
			unique = append(unique, tag)
		}
	}
	return slices.Clip(unique)
}
