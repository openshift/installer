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
	"slices"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

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
