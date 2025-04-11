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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// WifSupportBuilder contains the data and logic needed to build 'wif_support' objects.
type WifSupportBuilder struct {
	bitmap_   uint32
	principal string
	roles     []*WifRoleBuilder
}

// NewWifSupport creates a new builder of 'wif_support' objects.
func NewWifSupport() *WifSupportBuilder {
	return &WifSupportBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifSupportBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Principal sets the value of the 'principal' attribute to the given value.
func (b *WifSupportBuilder) Principal(value string) *WifSupportBuilder {
	b.principal = value
	b.bitmap_ |= 1
	return b
}

// Roles sets the value of the 'roles' attribute to the given values.
func (b *WifSupportBuilder) Roles(values ...*WifRoleBuilder) *WifSupportBuilder {
	b.roles = make([]*WifRoleBuilder, len(values))
	copy(b.roles, values)
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifSupportBuilder) Copy(object *WifSupport) *WifSupportBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.principal = object.principal
	if object.roles != nil {
		b.roles = make([]*WifRoleBuilder, len(object.roles))
		for i, v := range object.roles {
			b.roles[i] = NewWifRole().Copy(v)
		}
	} else {
		b.roles = nil
	}
	return b
}

// Build creates a 'wif_support' object using the configuration stored in the builder.
func (b *WifSupportBuilder) Build() (object *WifSupport, err error) {
	object = new(WifSupport)
	object.bitmap_ = b.bitmap_
	object.principal = b.principal
	if b.roles != nil {
		object.roles = make([]*WifRole, len(b.roles))
		for i, v := range b.roles {
			object.roles[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
