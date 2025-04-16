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
)

// Tags defines a map of tags.
type Tags map[string]string

// Equals returns true if the tags are equal.
func (t Tags) Equals(other Tags) bool {
	return reflect.DeepEqual(t, other)
}

// HasMatchingSpecVersionHash returns true if the resource has been tagged with a matching resource spec hash value.
func (t Tags) HasMatchingSpecVersionHash(hash string) bool {
	value, ok := t[SpecVersionHashTagKey()]
	return ok && value == hash
}

// HasOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of this management tooling.
func (t Tags) HasOwned(cluster string) bool {
	value, ok := t[ClusterTagKey(cluster)]
	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// HasAzureCloudProviderOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of the in-tree cloud provider.
func (t Tags) HasAzureCloudProviderOwned(cluster string) bool {
	value, ok := t[ClusterAzureCloudProviderTagKey(cluster)]
	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// GetRole returns the Cluster API role for the tagged resource.
func (t Tags) GetRole() string {
	return t[NameAzureClusterAPIRole]
}

// Difference returns the difference between this map of tags and the other map of tags.
// Items are considered equals if key and value are equals.
func (t Tags) Difference(other Tags) Tags {
	res := make(Tags, len(t))

	for key, value := range t {
		if otherValue, ok := other[key]; ok && value == otherValue {
			continue
		}
		res[key] = value
	}

	return res
}

// Merge merges in tags from other. If a tag already exists, it is replaced by the tag in other.
func (t Tags) Merge(other Tags) {
	for k, v := range other {
		t[k] = v
	}
}

// AddSpecVersionHashTag adds a spec version hash to the Azure resource tags to determine quickly if state has changed.
func (t Tags) AddSpecVersionHashTag(hash string) Tags {
	t[SpecVersionHashTagKey()] = hash
	return t
}

// ResourceLifecycle configures the lifecycle of a resource.
type ResourceLifecycle string

const (
	// ResourceLifecycleOwned is the value we use when tagging resources to indicate
	// that the resource is considered owned and managed by the cluster,
	// and in particular that the lifecycle is tied to the lifecycle of the cluster.
	ResourceLifecycleOwned = ResourceLifecycle("owned")

	// ResourceLifecycleShared is the value we use when tagging resources to indicate
	// that the resource is shared between multiple clusters, and should not be destroyed
	// if the cluster is destroyed.
	ResourceLifecycleShared = ResourceLifecycle("shared")

	// NameKubernetesAzureCloudProviderPrefix is the tag name used by the cloud provider to logically
	// separate independent cluster resources. We use it to identify which resources we expect
	// to be permissive about state changes.
	// logically independent clusters running in the same AZ.
	// The tag key = NameKubernetesAzureCloudProviderPrefix + clusterID.
	// The tag value is an ownership value.
	NameKubernetesAzureCloudProviderPrefix = "kubernetes.io_cluster_"

	// NameAzureProviderPrefix is the tag prefix we use to differentiate
	// cluster-api-provider-azure owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameAzureProviderPrefix = "sigs.k8s.io_cluster-api-provider-azure_"

	// NameAzureProviderOwned is the tag name we use to differentiate
	// cluster-api-provider-azure owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameAzureProviderOwned = NameAzureProviderPrefix + "cluster_"

	// NameAzureClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	NameAzureClusterAPIRole = NameAzureProviderPrefix + "role"

	// APIServerRole describes the value for the apiserver role.
	APIServerRole = "apiserver"

	// APIServerRoleInternal describes the value for the apiserver-internal role,
	// an identifier for an internal load balancer serving apiserver traffic for cluster nodes.
	APIServerRoleInternal = "apiserver-internal"

	// NodeOutboundRole describes the value for the node outbound LB role.
	NodeOutboundRole = "nodeOutbound"

	// ControlPlaneOutboundRole describes the value for the control plane outbound LB role.
	ControlPlaneOutboundRole = "controlPlaneOutbound"

	// BastionRole describes the value for the bastion role.
	BastionRole = Bastion

	// CommonRole describes the value for the common role.
	CommonRole = "common"

	// VMTagsLastAppliedAnnotation is the key for the machine object annotation
	// which tracks the AdditionalTags in the Machine Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	// Deprecated: use azure.VMTagsLastAppliedAnnotation instead. This constant will be removed in v1beta2.
	VMTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags-vm"

	// RGTagsLastAppliedAnnotation is the key for the Azure Cluster object annotation
	// which tracks the AdditionalTags for Resource Group which is part in the Azure Cluster.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	// Deprecated: use azure.RGTagsLastAppliedAnnotation instead. This constant will be removed in v1beta2.
	RGTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags-rg"
)

// SpecVersionHashTagKey is the key for the spec version hash used to enable quick spec difference comparison.
func SpecVersionHashTagKey() string {
	return fmt.Sprintf("%s%s", NameAzureProviderPrefix, "spec-version-hash")
}

// ClusterTagKey generates the key for resources associated with a cluster.
func ClusterTagKey(name string) string {
	return fmt.Sprintf("%s%s", NameAzureProviderOwned, name)
}

// ClusterAzureCloudProviderTagKey generates the key for resources associated a cluster's Azure cloud provider.
func ClusterAzureCloudProviderTagKey(name string) string {
	return fmt.Sprintf("%s%s", NameKubernetesAzureCloudProviderPrefix, name)
}

// BuildParams is used to build tags around an azure resource.
type BuildParams struct {
	// Lifecycle determines the resource lifecycle.
	Lifecycle ResourceLifecycle

	// ClusterName is the cluster associated with the resource.
	ClusterName string

	// ResourceID is the unique identifier of the resource to be tagged.
	ResourceID string

	// Name is the name of the resource, it's applied as the tag "Name" on Azure.
	// +optional
	Name *string

	// Role is the role associated to the resource.
	// +optional
	Role *string

	// Any additional tags to be added to the resource.
	// +optional
	Additional Tags
}

// Build builds tags including the cluster tag and returns them in map form.
func Build(params BuildParams) Tags {
	tags := make(Tags)
	for k, v := range params.Additional {
		tags[k] = v
	}

	tags[ClusterTagKey(params.ClusterName)] = string(params.Lifecycle)
	if params.Role != nil {
		tags[NameAzureClusterAPIRole] = *params.Role
	}

	if params.Name != nil {
		tags["Name"] = *params.Name
	}

	return tags
}
