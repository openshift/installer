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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// Details for `htpasswd` identity providers.
type HTPasswdIdentityProviderBuilder struct {
	fieldSet_ []bool
	password  string
	username  string
	users     *v1.HTPasswdUserListBuilder
}

// NewHTPasswdIdentityProvider creates a new builder of 'HT_passwd_identity_provider' objects.
func NewHTPasswdIdentityProvider() *HTPasswdIdentityProviderBuilder {
	return &HTPasswdIdentityProviderBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HTPasswdIdentityProviderBuilder) Empty() bool {
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

// Password sets the value of the 'password' attribute to the given value.
func (b *HTPasswdIdentityProviderBuilder) Password(value string) *HTPasswdIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.password = value
	b.fieldSet_[0] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *HTPasswdIdentityProviderBuilder) Username(value string) *HTPasswdIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.username = value
	b.fieldSet_[1] = true
	return b
}

// Users sets the value of the 'users' attribute to the given values.
func (b *HTPasswdIdentityProviderBuilder) Users(value *v1.HTPasswdUserListBuilder) *HTPasswdIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.users = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HTPasswdIdentityProviderBuilder) Copy(object *HTPasswdIdentityProvider) *HTPasswdIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.password = object.password
	b.username = object.username
	if object.users != nil {
		b.users = v1.NewHTPasswdUserList().Copy(object.users)
	} else {
		b.users = nil
	}
	return b
}

// Build creates a 'HT_passwd_identity_provider' object using the configuration stored in the builder.
func (b *HTPasswdIdentityProviderBuilder) Build() (object *HTPasswdIdentityProvider, err error) {
	object = new(HTPasswdIdentityProvider)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.password = b.password
	object.username = b.username
	if b.users != nil {
		object.users, err = b.users.Build()
		if err != nil {
			return
		}
	}
	return
}
