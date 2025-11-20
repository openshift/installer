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

type WifRoleBuilder struct {
	fieldSet_   []bool
	permissions []string
	roleId      string
	predefined  bool
}

// NewWifRole creates a new builder of 'wif_role' objects.
func NewWifRole() *WifRoleBuilder {
	return &WifRoleBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifRoleBuilder) Empty() bool {
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

// Permissions sets the value of the 'permissions' attribute to the given values.
func (b *WifRoleBuilder) Permissions(values ...string) *WifRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.permissions = make([]string, len(values))
	copy(b.permissions, values)
	b.fieldSet_[0] = true
	return b
}

// Predefined sets the value of the 'predefined' attribute to the given value.
func (b *WifRoleBuilder) Predefined(value bool) *WifRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.predefined = value
	b.fieldSet_[1] = true
	return b
}

// RoleId sets the value of the 'role_id' attribute to the given value.
func (b *WifRoleBuilder) RoleId(value string) *WifRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.roleId = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifRoleBuilder) Copy(object *WifRole) *WifRoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.permissions != nil {
		object.permissions = make([]string, len(b.permissions))
		copy(object.permissions, b.permissions)
	}
	object.predefined = b.predefined
	object.roleId = b.roleId
	return
}
