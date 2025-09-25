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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// HTPasswdIdentityProviderBuilder contains the data and logic needed to build 'HT_passwd_identity_provider' objects.
//
// Details for `htpasswd` identity providers.
type HTPasswdIdentityProviderBuilder struct {
	bitmap_  uint32
	password string
	username string
	users    *v1.HTPasswdUserListBuilder
}

// NewHTPasswdIdentityProvider creates a new builder of 'HT_passwd_identity_provider' objects.
func NewHTPasswdIdentityProvider() *HTPasswdIdentityProviderBuilder {
	return &HTPasswdIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HTPasswdIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Password sets the value of the 'password' attribute to the given value.
func (b *HTPasswdIdentityProviderBuilder) Password(value string) *HTPasswdIdentityProviderBuilder {
	b.password = value
	b.bitmap_ |= 1
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *HTPasswdIdentityProviderBuilder) Username(value string) *HTPasswdIdentityProviderBuilder {
	b.username = value
	b.bitmap_ |= 2
	return b
}

// Users sets the value of the 'users' attribute to the given values.
func (b *HTPasswdIdentityProviderBuilder) Users(value *v1.HTPasswdUserListBuilder) *HTPasswdIdentityProviderBuilder {
	b.users = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HTPasswdIdentityProviderBuilder) Copy(object *HTPasswdIdentityProvider) *HTPasswdIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
