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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	time "time"
)

// RegistryAllowlist represents a single registry allowlist.
type RegistryAllowlistBuilder struct {
	fieldSet_         []bool
	id                string
	href              string
	cloudProvider     *CloudProviderBuilder
	creationTimestamp time.Time
	registries        []string
}

// NewRegistryAllowlist creates a new builder of 'registry_allowlist' objects.
func NewRegistryAllowlist() *RegistryAllowlistBuilder {
	return &RegistryAllowlistBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *RegistryAllowlistBuilder) Link(value bool) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *RegistryAllowlistBuilder) ID(value string) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *RegistryAllowlistBuilder) HREF(value string) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryAllowlistBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *RegistryAllowlistBuilder) CloudProvider(value *CloudProviderBuilder) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.cloudProvider = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *RegistryAllowlistBuilder) CreationTimestamp(value time.Time) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.creationTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// Registries sets the value of the 'registries' attribute to the given values.
func (b *RegistryAllowlistBuilder) Registries(values ...string) *RegistryAllowlistBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.registries = make([]string, len(values))
	copy(b.registries, values)
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryAllowlistBuilder) Copy(object *RegistryAllowlist) *RegistryAllowlistBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.cloudProvider != nil {
		b.cloudProvider = NewCloudProvider().Copy(object.cloudProvider)
	} else {
		b.cloudProvider = nil
	}
	b.creationTimestamp = object.creationTimestamp
	if object.registries != nil {
		b.registries = make([]string, len(object.registries))
		copy(b.registries, object.registries)
	} else {
		b.registries = nil
	}
	return b
}

// Build creates a 'registry_allowlist' object using the configuration stored in the builder.
func (b *RegistryAllowlistBuilder) Build() (object *RegistryAllowlist, err error) {
	object = new(RegistryAllowlist)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.cloudProvider != nil {
		object.cloudProvider, err = b.cloudProvider.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	if b.registries != nil {
		object.registries = make([]string, len(b.registries))
		copy(object.registries, b.registries)
	}
	return
}
