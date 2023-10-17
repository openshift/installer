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

package v1beta1

import (
	"fmt"
	"reflect"
	"strings"
)

// Labels defines a map of tags.
type Labels map[string]string

// Equals returns true if the tags are equal.
func (in Labels) Equals(other Labels) bool {
	return reflect.DeepEqual(in, other)
}

// HasOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of this management tooling.
func (in Labels) HasOwned(cluster string) bool {
	value, ok := in[ClusterTagKey(cluster)]

	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// // HasOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of the in-tree cloud provider.
// func (in Labels) HasGCPCloudProviderOwned(cluster string) bool {
// 	value, ok := t[ClusterGCPCloudProviderTagKey(cluster)]
// 	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
// }

// GetRole returns the Cluster API role for the tagged resource.
func (in Labels) GetRole() string {
	return in[NameGCPClusterAPIRole]
}

// ToComputeFilter returns the string representation of the labels as a filter
// to be used in google compute sdk calls.
func (in Labels) ToComputeFilter() string {
	var builder strings.Builder
	for k, v := range in {
		builder.WriteString(fmt.Sprintf("(labels.%s = %q) ", k, v))
	}

	return builder.String()
}

// Difference returns the difference between this map of tags and the other map of tags.
// Items are considered equals if key and value are equals.
func (in Labels) Difference(other Labels) Labels {
	res := make(Labels, len(in))

	for key, value := range in {
		if otherValue, ok := other[key]; ok && value == otherValue {
			continue
		}
		res[key] = value
	}

	return res
}

// AddLabels adds (and overwrites) the current labels with the ones passed in.
func (in Labels) AddLabels(other Labels) Labels {
	for key, value := range other {
		if in == nil {
			in = make(map[string]string, len(other))
		}
		in[key] = value
	}

	return in
}

// ResourceLifecycle configures the lifecycle of a resource.
type ResourceLifecycle string

const (
	// ResourceLifecycleOwned is the value we use when tagging resources to indicate
	// that the resource is considered owned and managed by the cluster,
	// and in particular that the lifecycle is tied to the lifecycle of the cluster.
	ResourceLifecycleOwned = ResourceLifecycle("owned")

	// NameGCPProviderPrefix is the tag prefix we use to differentiate
	// cluster-api-provider-gcp owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameGCPProviderPrefix = "capg-"

	// NameGCPProviderOwned is the tag name we use to differentiate
	// cluster-api-provider-gcp owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameGCPProviderOwned = NameGCPProviderPrefix + "cluster-"

	// NameGCPClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	NameGCPClusterAPIRole = NameGCPProviderPrefix + "role"

	// APIServerRoleTagValue describes the value for the apiserver role.
	APIServerRoleTagValue = "apiserver"
)

// ClusterTagKey generates the key for resources associated with a cluster.
func ClusterTagKey(name string) string {
	return fmt.Sprintf("%s%s", NameGCPProviderOwned, name)
}

// ClusterGCPCloudProviderTagKey generates the key for resources associated a cluster's GCP cloud provider.
// func ClusterGCPCloudProviderTagKey(name string) string {
// return fmt.Sprintf("%s%s", NameKubernetesGCPCloudProviderPrefix, name)
// }

// BuildParams is used to build tags around an gcp resource.
type BuildParams struct {
	// Lifecycle determines the resource lifecycle.
	Lifecycle ResourceLifecycle

	// ClusterName is the cluster associated with the resource.
	ClusterName string

	// ResourceID is the unique identifier of the resource to be tagged.
	ResourceID string

	// Role is the role associated to the resource.
	// +optional
	Role *string

	// Any additional tags to be added to the resource.
	// +optional
	Additional Labels
}

// Build builds tags including the cluster tag and returns them in map form.
func Build(params BuildParams) Labels {
	tags := make(Labels)
	for k, v := range params.Additional {
		tags[strings.ToLower(k)] = strings.ToLower(v)
	}

	tags[ClusterTagKey(params.ClusterName)] = string(params.Lifecycle)
	if params.Role != nil {
		tags[NameGCPClusterAPIRole] = strings.ToLower(*params.Role)
	}

	return tags
}
