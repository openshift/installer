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

// GroupBuilder contains the data and logic needed to build 'group' objects.
//
// Representation of a group of users.
type GroupBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	users   *UserListBuilder
}

// NewGroup creates a new builder of 'group' objects.
func NewGroup() *GroupBuilder {
	return &GroupBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *GroupBuilder) Link(value bool) *GroupBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *GroupBuilder) ID(value string) *GroupBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *GroupBuilder) HREF(value string) *GroupBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GroupBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Users sets the value of the 'users' attribute to the given values.
func (b *GroupBuilder) Users(value *UserListBuilder) *GroupBuilder {
	b.users = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GroupBuilder) Copy(object *Group) *GroupBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.users != nil {
		b.users = NewUserList().Copy(object.users)
	} else {
		b.users = nil
	}
	return b
}

// Build creates a 'group' object using the configuration stored in the builder.
func (b *GroupBuilder) Build() (object *Group, err error) {
	object = new(Group)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.users != nil {
		object.users, err = b.users.Build()
		if err != nil {
			return
		}
	}
	return
}
