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

// WifRoleBuilder contains the data and logic needed to build 'wif_role' objects.
type WifRoleBuilder struct {
	bitmap_     uint32
	permissions []string
	roleId      string
	predefined  bool
}

// NewWifRole creates a new builder of 'wif_role' objects.
func NewWifRole() *WifRoleBuilder {
	return &WifRoleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifRoleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Permissions sets the value of the 'permissions' attribute to the given values.
func (b *WifRoleBuilder) Permissions(values ...string) *WifRoleBuilder {
	b.permissions = make([]string, len(values))
	copy(b.permissions, values)
	b.bitmap_ |= 1
	return b
}

// Predefined sets the value of the 'predefined' attribute to the given value.
func (b *WifRoleBuilder) Predefined(value bool) *WifRoleBuilder {
	b.predefined = value
	b.bitmap_ |= 2
	return b
}

// RoleId sets the value of the 'role_id' attribute to the given value.
func (b *WifRoleBuilder) RoleId(value string) *WifRoleBuilder {
	b.roleId = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifRoleBuilder) Copy(object *WifRole) *WifRoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.permissions != nil {
		b.permissions = make([]string, len(object.permissions))
		copy(b.permissions, object.permissions)
	} else {
		b.permissions = nil
	}
	b.predefined = object.predefined
	b.roleId = object.roleId
	return b
}

// Build creates a 'wif_role' object using the configuration stored in the builder.
func (b *WifRoleBuilder) Build() (object *WifRole, err error) {
	object = new(WifRole)
	object.bitmap_ = b.bitmap_
	if b.permissions != nil {
		object.permissions = make([]string, len(b.permissions))
		copy(object.permissions, b.permissions)
	}
	object.predefined = b.predefined
	object.roleId = b.roleId
	return
}
