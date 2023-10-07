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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	time "time"
)

// UserBuilder contains the data and logic needed to build 'user' objects.
//
// Definition of a Web RCA user.
type UserBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	email     string
	name      string
	updatedAt time.Time
	username  string
	fromAuth  bool
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *UserBuilder) CreatedAt(value time.Time) *UserBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *UserBuilder) DeletedAt(value time.Time) *UserBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// Email sets the value of the 'email' attribute to the given value.
func (b *UserBuilder) Email(value string) *UserBuilder {
	b.email = value
	b.bitmap_ |= 32
	return b
}

// FromAuth sets the value of the 'from_auth' attribute to the given value.
func (b *UserBuilder) FromAuth(value bool) *UserBuilder {
	b.fromAuth = value
	b.bitmap_ |= 64
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *UserBuilder) Name(value string) *UserBuilder {
	b.name = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *UserBuilder) UpdatedAt(value time.Time) *UserBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *UserBuilder) Username(value string) *UserBuilder {
	b.username = value
	b.bitmap_ |= 512
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UserBuilder) Copy(object *User) *UserBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	b.email = object.email
	b.fromAuth = object.fromAuth
	b.name = object.name
	b.updatedAt = object.updatedAt
	b.username = object.username
	return b
}

// Build creates a 'user' object using the configuration stored in the builder.
func (b *UserBuilder) Build() (object *User, err error) {
	object = new(User)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.email = b.email
	object.fromAuth = b.fromAuth
	object.name = b.name
	object.updatedAt = b.updatedAt
	object.username = b.username
	return
}
