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

type HTPasswdUserBuilder struct {
	fieldSet_      []bool
	id             string
	hashedPassword string
	password       string
	username       string
}

// NewHTPasswdUser creates a new builder of 'HT_passwd_user' objects.
func NewHTPasswdUser() *HTPasswdUserBuilder {
	return &HTPasswdUserBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HTPasswdUserBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *HTPasswdUserBuilder) ID(value string) *HTPasswdUserBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// HashedPassword sets the value of the 'hashed_password' attribute to the given value.
func (b *HTPasswdUserBuilder) HashedPassword(value string) *HTPasswdUserBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.hashedPassword = value
	b.fieldSet_[1] = true
	return b
}

// Password sets the value of the 'password' attribute to the given value.
func (b *HTPasswdUserBuilder) Password(value string) *HTPasswdUserBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.password = value
	b.fieldSet_[2] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *HTPasswdUserBuilder) Username(value string) *HTPasswdUserBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.username = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HTPasswdUserBuilder) Copy(object *HTPasswdUser) *HTPasswdUserBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.hashedPassword = object.hashedPassword
	b.password = object.password
	b.username = object.username
	return b
}

// Build creates a 'HT_passwd_user' object using the configuration stored in the builder.
func (b *HTPasswdUserBuilder) Build() (object *HTPasswdUser, err error) {
	object = new(HTPasswdUser)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.hashedPassword = b.hashedPassword
	object.password = b.password
	object.username = b.username
	return
}
