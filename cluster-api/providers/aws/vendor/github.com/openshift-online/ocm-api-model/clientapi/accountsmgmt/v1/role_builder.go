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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

type RoleBuilder struct {
	fieldSet_   []bool
	id          string
	href        string
	name        string
	permissions []*PermissionBuilder
}

// NewRole creates a new builder of 'role' objects.
func NewRole() *RoleBuilder {
	return &RoleBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *RoleBuilder) Link(value bool) *RoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *RoleBuilder) ID(value string) *RoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *RoleBuilder) HREF(value string) *RoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RoleBuilder) Empty() bool {
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

// Name sets the value of the 'name' attribute to the given value.
func (b *RoleBuilder) Name(value string) *RoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.name = value
	b.fieldSet_[3] = true
	return b
}

// Permissions sets the value of the 'permissions' attribute to the given values.
func (b *RoleBuilder) Permissions(values ...*PermissionBuilder) *RoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.permissions = make([]*PermissionBuilder, len(values))
	copy(b.permissions, values)
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RoleBuilder) Copy(object *Role) *RoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.name = object.name
	if object.permissions != nil {
		b.permissions = make([]*PermissionBuilder, len(object.permissions))
		for i, v := range object.permissions {
			b.permissions[i] = NewPermission().Copy(v)
		}
	} else {
		b.permissions = nil
	}
	return b
}

// Build creates a 'role' object using the configuration stored in the builder.
func (b *RoleBuilder) Build() (object *Role, err error) {
	object = new(Role)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	if b.permissions != nil {
		object.permissions = make([]*Permission, len(b.permissions))
		for i, v := range b.permissions {
			object.permissions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
