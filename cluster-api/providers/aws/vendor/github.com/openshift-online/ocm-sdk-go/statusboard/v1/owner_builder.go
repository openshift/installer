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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

// OwnerBuilder contains the data and logic needed to build 'owner' objects.
//
// Definition of a Status Board owner.
type OwnerBuilder struct {
	bitmap_  uint32
	id       string
	href     string
	email    string
	username string
}

// NewOwner creates a new builder of 'owner' objects.
func NewOwner() *OwnerBuilder {
	return &OwnerBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *OwnerBuilder) Link(value bool) *OwnerBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *OwnerBuilder) ID(value string) *OwnerBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *OwnerBuilder) HREF(value string) *OwnerBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OwnerBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Email sets the value of the 'email' attribute to the given value.
func (b *OwnerBuilder) Email(value string) *OwnerBuilder {
	b.email = value
	b.bitmap_ |= 8
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *OwnerBuilder) Username(value string) *OwnerBuilder {
	b.username = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OwnerBuilder) Copy(object *Owner) *OwnerBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.email = object.email
	b.username = object.username
	return b
}

// Build creates a 'owner' object using the configuration stored in the builder.
func (b *OwnerBuilder) Build() (object *Owner, err error) {
	object = new(Owner)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.email = b.email
	object.username = b.username
	return
}
