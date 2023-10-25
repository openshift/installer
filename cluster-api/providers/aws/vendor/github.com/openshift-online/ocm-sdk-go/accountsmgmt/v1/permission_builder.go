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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// PermissionBuilder contains the data and logic needed to build 'permission' objects.
type PermissionBuilder struct {
	bitmap_  uint32
	id       string
	href     string
	action   Action
	resource string
}

// NewPermission creates a new builder of 'permission' objects.
func NewPermission() *PermissionBuilder {
	return &PermissionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *PermissionBuilder) Link(value bool) *PermissionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *PermissionBuilder) ID(value string) *PermissionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *PermissionBuilder) HREF(value string) *PermissionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PermissionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Action sets the value of the 'action' attribute to the given value.
//
// Possible actions for a permission.
func (b *PermissionBuilder) Action(value Action) *PermissionBuilder {
	b.action = value
	b.bitmap_ |= 8
	return b
}

// Resource sets the value of the 'resource' attribute to the given value.
func (b *PermissionBuilder) Resource(value string) *PermissionBuilder {
	b.resource = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PermissionBuilder) Copy(object *Permission) *PermissionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.action = object.action
	b.resource = object.resource
	return b
}

// Build creates a 'permission' object using the configuration stored in the builder.
func (b *PermissionBuilder) Build() (object *Permission, err error) {
	object = new(Permission)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.action = b.action
	object.resource = b.resource
	return
}
