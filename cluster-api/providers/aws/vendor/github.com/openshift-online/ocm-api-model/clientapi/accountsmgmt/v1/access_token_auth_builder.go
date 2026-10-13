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

type AccessTokenAuthBuilder struct {
	fieldSet_ []bool
	auth      string
	email     string
}

// NewAccessTokenAuth creates a new builder of 'access_token_auth' objects.
func NewAccessTokenAuth() *AccessTokenAuthBuilder {
	return &AccessTokenAuthBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessTokenAuthBuilder) Empty() bool {
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

// Auth sets the value of the 'auth' attribute to the given value.
func (b *AccessTokenAuthBuilder) Auth(value string) *AccessTokenAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.auth = value
	b.fieldSet_[0] = true
	return b
}

// Email sets the value of the 'email' attribute to the given value.
func (b *AccessTokenAuthBuilder) Email(value string) *AccessTokenAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.email = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessTokenAuthBuilder) Copy(object *AccessTokenAuth) *AccessTokenAuthBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.auth = object.auth
	b.email = object.email
	return b
}

// Build creates a 'access_token_auth' object using the configuration stored in the builder.
func (b *AccessTokenAuthBuilder) Build() (object *AccessTokenAuth, err error) {
	object = new(AccessTokenAuth)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.auth = b.auth
	object.email = b.email
	return
}
