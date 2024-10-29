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
	"context"
	"errors"
	"fmt"
	"slices"

	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

// ResolveServerSpec is responsible for populating a ResolvedServerSpec from
// an OpenStackMachineSpec and any external dependencies. The result contains no
// external dependencies, and does not require any complex logic on creation.
// Note that we only set the fields in ResolvedServerSpec that are not set yet. This is ok because
// OpenStackServer is immutable, so we can't change the spec after the machine is created.
func ResolveServerSpec(ctx context.Context, scope *scope.WithLogger, k8sClient client.Client, openStackServer *infrav1alpha1.OpenStackServer) (bool, bool, error) {
	resolved := openStackServer.Status.Resolved
	if resolved == nil {
		resolved = &infrav1alpha1.ResolvedServerSpec{}
		openStackServer.Status.Resolved = resolved
	}

	spec := &openStackServer.Spec
	if resolved == nil {
		resolved = &infrav1alpha1.ResolvedServerSpec{}
		openStackServer.Status.Resolved = resolved
	}

	// If the server is bound to a cluster, we use the cluster name to generate the port description.
	var clusterName string
	if openStackServer.ObjectMeta.Labels[clusterv1.ClusterNameLabel] != "" {
		clusterName = openStackServer.ObjectMeta.Labels[clusterv1.ClusterNameLabel]
	}

	computeService, err := NewService(scope)
	if err != nil {
		return false, false, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return false, false, err
	}

	// A setter returns: done, changed, error
	type setterFn func() (bool, bool, error)

	serverGroup := func() (bool, bool, error) {
		if spec.ServerGroup == nil || resolved.ServerGroupID != "" {
			return true, false, nil
		}
		serverGroupID, err := computeService.GetServerGroupID(spec.ServerGroup)
		if err != nil {
			return false, false, err
		}
		resolved.ServerGroupID = serverGroupID
		return true, true, nil
	}

	imageID := func() (bool, bool, error) {
		if resolved.ImageID != "" {
			return true, false, nil
		}

		imageID, err := computeService.GetImageID(ctx, k8sClient, openStackServer.Namespace, spec.Image)
		if err != nil {
			return false, false, err
		}

		// If we didn't get an imageID it means we're waiting on a dependency.
		// Wait to be called again.
		if imageID == nil {
			return false, false, nil
		}
		resolved.ImageID = *imageID
		return true, true, nil
	}

	flavorID := func() (bool, bool, error) {
		if resolved.FlavorID != "" {
			return true, false, nil
		}

		flavorID, err := computeService.GetFlavorID(spec.FlavorID, spec.Flavor)
		if err != nil {
			return false, false, err
		}

		resolved.FlavorID = flavorID
		return true, true, nil
	}

	ports := func() (bool, bool, error) {
		if len(resolved.Ports) > 0 {
			return true, false, nil
		}

		// Network resources are required in order to get ports options.
		// Notes:
		// - clusterResourceName is not used in this context, so we pass an empty string. In the future,
		// we may want to remove that (it's only used for the port description) or allow a user to pass
		// a custom description.
		// - managedSecurityGroup is not used in this context, so we pass nil. The security groups are
		//   passed in the spec.SecurityGroups and spec.Ports.
		// - We run a safety check to ensure that the resolved.Ports has the same length as the spec.Ports.
		//   This is to ensure that we don't accidentally add ports to the resolved.Ports that are not in the spec.
		specTrunk := ptr.Deref(spec.Trunk, false)
		portsOpts, err := networkingService.ConstructPorts(spec.Ports, spec.SecurityGroups, specTrunk, clusterName, openStackServer.Name, nil, nil, spec.Tags)
		if err != nil {
			return false, false, err
		}
		if portsOpts != nil && len(portsOpts) != len(spec.Ports) {
			return false, false, fmt.Errorf("resolved.Ports has a different length than spec.Ports")
		}
		resolved.Ports = portsOpts
		return true, true, nil
	}

	// Execute all setters and collate their return values
	var errs []error
	changed := false
	done := true
	for _, setter := range []setterFn{serverGroup, imageID, flavorID, ports} {
		thisDone, thisChanged, err := setter()
		changed = changed || thisChanged
		done = done && thisDone
		if err != nil {
			errs = append(errs, err)
		}
	}

	return changed, done, errors.Join(errs...)
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
