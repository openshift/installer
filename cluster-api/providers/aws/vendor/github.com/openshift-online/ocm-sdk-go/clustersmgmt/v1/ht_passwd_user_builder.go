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

// HTPasswdUserBuilder contains the data and logic needed to build 'HT_passwd_user' objects.
type HTPasswdUserBuilder struct {
	bitmap_        uint32
	id             string
	hashedPassword string
	password       string
	username       string
}

// NewHTPasswdUser creates a new builder of 'HT_passwd_user' objects.
func NewHTPasswdUser() *HTPasswdUserBuilder {
	return &HTPasswdUserBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HTPasswdUserBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *HTPasswdUserBuilder) ID(value string) *HTPasswdUserBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// HashedPassword sets the value of the 'hashed_password' attribute to the given value.
func (b *HTPasswdUserBuilder) HashedPassword(value string) *HTPasswdUserBuilder {
	b.hashedPassword = value
	b.bitmap_ |= 2
	return b
}

// Password sets the value of the 'password' attribute to the given value.
func (b *HTPasswdUserBuilder) Password(value string) *HTPasswdUserBuilder {
	b.password = value
	b.bitmap_ |= 4
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *HTPasswdUserBuilder) Username(value string) *HTPasswdUserBuilder {
	b.username = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HTPasswdUserBuilder) Copy(object *HTPasswdUser) *HTPasswdUserBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.hashedPassword = object.hashedPassword
	b.password = object.password
	b.username = object.username
	return b
}

// Build creates a 'HT_passwd_user' object using the configuration stored in the builder.
func (b *HTPasswdUserBuilder) Build() (object *HTPasswdUser, err error) {
	object = new(HTPasswdUser)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.hashedPassword = b.hashedPassword
	object.password = b.password
	object.username = b.username
	return
}
