/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// PrivateLinkClusterConfigurationBuilder contains the data and logic needed to build 'private_link_cluster_configuration' objects.
//
// Manages the configuration for the Private Links.
type PrivateLinkClusterConfigurationBuilder struct {
	bitmap_    uint32
	principals []*PrivateLinkPrincipalBuilder
}

// NewPrivateLinkClusterConfiguration creates a new builder of 'private_link_cluster_configuration' objects.
func NewPrivateLinkClusterConfiguration() *PrivateLinkClusterConfigurationBuilder {
	return &PrivateLinkClusterConfigurationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PrivateLinkClusterConfigurationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Principals sets the value of the 'principals' attribute to the given values.
func (b *PrivateLinkClusterConfigurationBuilder) Principals(values ...*PrivateLinkPrincipalBuilder) *PrivateLinkClusterConfigurationBuilder {
	b.principals = make([]*PrivateLinkPrincipalBuilder, len(values))
	copy(b.principals, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PrivateLinkClusterConfigurationBuilder) Copy(object *PrivateLinkClusterConfiguration) *PrivateLinkClusterConfigurationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.principals != nil {
		b.principals = make([]*PrivateLinkPrincipalBuilder, len(object.principals))
		for i, v := range object.principals {
			b.principals[i] = NewPrivateLinkPrincipal().Copy(v)
		}
	} else {
		b.principals = nil
	}
	return b
}

// Build creates a 'private_link_cluster_configuration' object using the configuration stored in the builder.
func (b *PrivateLinkClusterConfigurationBuilder) Build() (object *PrivateLinkClusterConfiguration, err error) {
	object = new(PrivateLinkClusterConfiguration)
	object.bitmap_ = b.bitmap_
	if b.principals != nil {
		object.principals = make([]*PrivateLinkPrincipal, len(b.principals))
		for i, v := range b.principals {
			object.principals[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
