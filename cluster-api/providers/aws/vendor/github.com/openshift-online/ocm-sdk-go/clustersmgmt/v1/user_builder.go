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

// UserBuilder contains the data and logic needed to build 'user' objects.
//
// Representation of a user.
type UserBuilder struct {
	bitmap_ uint32
	id      string
	href    string
}

// NewUser creates a new builder of 'user' objects.
func NewUser() *UserBuilder {
	return &UserBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *UserBuilder) Link(value bool) *UserBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *UserBuilder) ID(value string) *UserBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *UserBuilder) HREF(value string) *UserBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *UserBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UserBuilder) Copy(object *User) *UserBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	return b
}

// Build creates a 'user' object using the configuration stored in the builder.
func (b *UserBuilder) Build() (object *User, err error) {
	object = new(User)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	return
}
