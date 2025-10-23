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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Manages the configuration for the Private Links.
type PrivateLinkConfigurationBuilder struct {
	fieldSet_  []bool
	principals *PrivateLinkPrincipalsBuilder
}

// NewPrivateLinkConfiguration creates a new builder of 'private_link_configuration' objects.
func NewPrivateLinkConfiguration() *PrivateLinkConfigurationBuilder {
	return &PrivateLinkConfigurationBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PrivateLinkConfigurationBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Principals sets the value of the 'principals' attribute to the given value.
//
// Contains a list of principals for the Private Link.
func (b *PrivateLinkConfigurationBuilder) Principals(value *PrivateLinkPrincipalsBuilder) *PrivateLinkConfigurationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.principals = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PrivateLinkConfigurationBuilder) Copy(object *PrivateLinkConfiguration) *PrivateLinkConfigurationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.principals != nil {
		b.principals = NewPrivateLinkPrincipals().Copy(object.principals)
	} else {
		b.principals = nil
	}
	return b
}

// Build creates a 'private_link_configuration' object using the configuration stored in the builder.
func (b *PrivateLinkConfigurationBuilder) Build() (object *PrivateLinkConfiguration, err error) {
	object = new(PrivateLinkConfiguration)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.principals != nil {
		object.principals, err = b.principals.Build()
		if err != nil {
			return
		}
	}
	return
}
