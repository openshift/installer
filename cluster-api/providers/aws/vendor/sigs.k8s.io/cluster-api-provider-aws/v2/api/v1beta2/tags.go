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

package v1beta2

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Tags defines a map of tags.
type Tags map[string]string

// HasOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of this management tooling.
func (t Tags) HasOwned(cluster string) bool {
	value, ok := t[ClusterTagKey(cluster)]
	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// HasAWSCloudProviderOwned returns true if the tags contains a tag that marks the resource as owned by the cluster from the perspective of the in-tree cloud provider.
func (t Tags) HasAWSCloudProviderOwned(cluster string) bool {
	value, ok := t[ClusterAWSCloudProviderTagKey(cluster)]
	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// GetRole returns the Cluster API role for the tagged resource.
func (t Tags) GetRole() string {
	return t[NameAWSClusterAPIRole]
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

// Validate checks if tags are valid for the AWS API/Resources.
// Keys must have at least 1 and max 128 characters.
// Values must be max 256 characters long.
// Keys and Values can only have alphabets, numbers, spaces and _ . : / = + - @ as characters.
// Tag's key cannot have prefix "aws:".
// Max count of User tags for a specific resource can be 50.
func (t Tags) Validate() []*field.Error {
	// Defines the maximum number of user tags which can be created for a specific resource
	const maxUserTagsAllowed = 50
	var errs field.ErrorList
	var userTagCount = len(t)
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\_\.\:\=\+\-\@\/]*$`)

	for k, v := range t {
		if len(k) < 1 {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), k, "key cannot be empty"),
			)
		}
		if len(k) > 128 {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), k, "key cannot be longer than 128 characters"),
			)
		}
		if len(v) > 256 {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), v, "value cannot be longer than 256 characters"),
			)
		}
		if wrongUserTagNomenclature(k) {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), k, "user created tag's key cannot have prefix aws:"),
			)
		}
		val := re.MatchString(k)
		if !val {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), k, "key cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ ."),
			)
		}
		val = re.MatchString(v)
		if !val {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "additionalTags"), v, "value cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ ."),
			)
		}
	}

	if userTagCount > maxUserTagsAllowed {
		errs = append(errs,
			field.Invalid(field.NewPath("spec", "additionalTags"), t, "user created tags cannot be more than 50"),
		)
	}

	return errs
}

// Checks whether the tag created is user tag or not.
func wrongUserTagNomenclature(k string) bool {
	return len(k) > 3 && k[0:4] == "aws:"
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

	// NameKubernetesAWSCloudProviderPrefix is the tag name used by the cloud provider to logically
	// separate independent cluster resources. We use it to identify which resources we expect
	// to be permissive about state changes.
	// logically independent clusters running in the same AZ.
	// The tag key = NameKubernetesAWSCloudProviderPrefix + clusterID
	// The tag value is an ownership value.
	NameKubernetesAWSCloudProviderPrefix = "kubernetes.io/cluster/"

	// NameAWSProviderPrefix is the tag prefix we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameAWSProviderPrefix = "sigs.k8s.io/cluster-api-provider-aws/"

	// NameAWSProviderOwned is the tag name we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses NameKubernetesClusterPrefix.
	NameAWSProviderOwned = NameAWSProviderPrefix + "cluster/"

	// NameAWSClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	NameAWSClusterAPIRole = NameAWSProviderPrefix + "role"

	// NameAWSSubnetAssociation is the tag name we use to mark association for resources
	// dedicated to this cluster api provider implementation.
	NameAWSSubnetAssociation = NameAWSProviderPrefix + "association"

	// SecondarySubnetTagValue is the secondary subnet tag constant value.
	SecondarySubnetTagValue = "secondary"

	// APIServerRoleTagValue describes the value for the apiserver role.
	APIServerRoleTagValue = "apiserver"

	// BastionRoleTagValue describes the value for the bastion role.
	BastionRoleTagValue = "bastion"

	// CommonRoleTagValue describes the value for the common role.
	CommonRoleTagValue = "common"

	// PublicRoleTagValue describes the value for the public role.
	PublicRoleTagValue = "public"

	// PrivateRoleTagValue describes the value for the private role.
	PrivateRoleTagValue = "private"

	// MachineNameTagKey is the key for machine name.
	MachineNameTagKey = "MachineName"

	// LaunchTemplateBootstrapDataSecret is the tag we use to store the `<namespace>/<name>`
	// of the bootstrap secret that was used to create the user data for the latest launch
	// template version.
	LaunchTemplateBootstrapDataSecret = NameAWSProviderPrefix + "bootstrap-data-secret"

	// LaunchTemplateBootstrapDataHash is the tag we use to store the hash of the raw bootstrap data.
	// If bootstrap data is stored in S3, this hash relates to that data, not to the EC2 instance
	// user data which only references the S3 object. We store this tag on launch template versions
	// so that S3 bootstrap data objects can be deleted when they get outdated.
	LaunchTemplateBootstrapDataHash = NameAWSProviderPrefix + "bootstrap-data-hash"
)

// ClusterTagKey generates the key for resources associated with a cluster.
func ClusterTagKey(name string) string {
	return fmt.Sprintf("%s%s", NameAWSProviderOwned, name)
}

// ClusterAWSCloudProviderTagKey generates the key for resources associated a cluster's AWS cloud provider.
func ClusterAWSCloudProviderTagKey(name string) string {
	return fmt.Sprintf("%s%s", NameKubernetesAWSCloudProviderPrefix, name)
}

// BuildParams is used to build tags around an aws resource.
type BuildParams struct {
	// Lifecycle determines the resource lifecycle.
	Lifecycle ResourceLifecycle

	// ClusterName is the cluster associated with the resource.
	ClusterName string

	// ResourceID is the unique identifier of the resource to be tagged.
	ResourceID string

	// Name is the name of the resource, it's applied as the tag "Name" on AWS.
	// +optional
	Name *string

	// Role is the role associated to the resource.
	// +optional
	Role *string

	// Any additional tags to be added to the resource.
	// +optional
	Additional Tags
}

// WithMachineName tags the namespaced machine name
// The machine name will be tagged with key "MachineName".
func (b BuildParams) WithMachineName(m *clusterv1.Machine) BuildParams {
	machineNamespacedName := types.NamespacedName{Namespace: m.Namespace, Name: m.Name}
	b.Additional[MachineNameTagKey] = machineNamespacedName.String()
	return b
}

// WithCloudProvider tags the cluster ownership for a resource.
func (b BuildParams) WithCloudProvider(name string) BuildParams {
	b.Additional[ClusterAWSCloudProviderTagKey(name)] = string(ResourceLifecycleOwned)
	return b
}

// Build builds tags including the cluster tag and returns them in map form.
func Build(params BuildParams) Tags {
	tags := make(Tags)

	// Add the name tag first so that it can be overwritten by a user-provided tag in the `Additional` tags.
	if params.Name != nil {
		tags["Name"] = *params.Name
	}

	for k, v := range params.Additional {
		tags[k] = v
	}

	if params.ClusterName != "" {
		tags[ClusterTagKey(params.ClusterName)] = string(params.Lifecycle)
	}
	if params.Role != nil {
		tags[NameAWSClusterAPIRole] = *params.Role
	}

	return tags
}
