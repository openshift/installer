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

// LDAPIdentityProviderBuilder contains the data and logic needed to build 'LDAP_identity_provider' objects.
//
// Details for `ldap` identity providers.
type LDAPIdentityProviderBuilder struct {
	bitmap_      uint32
	ca           string
	url          string
	attributes   *LDAPAttributesBuilder
	bindDN       string
	bindPassword string
	insecure     bool
}

// NewLDAPIdentityProvider creates a new builder of 'LDAP_identity_provider' objects.
func NewLDAPIdentityProvider() *LDAPIdentityProviderBuilder {
	return &LDAPIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LDAPIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CA sets the value of the 'CA' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) CA(value string) *LDAPIdentityProviderBuilder {
	b.ca = value
	b.bitmap_ |= 1
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) URL(value string) *LDAPIdentityProviderBuilder {
	b.url = value
	b.bitmap_ |= 2
	return b
}

// Attributes sets the value of the 'attributes' attribute to the given value.
//
// LDAP attributes used to configure the LDAP identity provider.
func (b *LDAPIdentityProviderBuilder) Attributes(value *LDAPAttributesBuilder) *LDAPIdentityProviderBuilder {
	b.attributes = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// BindDN sets the value of the 'bind_DN' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) BindDN(value string) *LDAPIdentityProviderBuilder {
	b.bindDN = value
	b.bitmap_ |= 8
	return b
}

// BindPassword sets the value of the 'bind_password' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) BindPassword(value string) *LDAPIdentityProviderBuilder {
	b.bindPassword = value
	b.bitmap_ |= 16
	return b
}

// Insecure sets the value of the 'insecure' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) Insecure(value bool) *LDAPIdentityProviderBuilder {
	b.insecure = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LDAPIdentityProviderBuilder) Copy(object *LDAPIdentityProvider) *LDAPIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.ca = object.ca
	b.url = object.url
	if object.attributes != nil {
		b.attributes = NewLDAPAttributes().Copy(object.attributes)
	} else {
		b.attributes = nil
	}
	b.bindDN = object.bindDN
	b.bindPassword = object.bindPassword
	b.insecure = object.insecure
	return b
}

// Build creates a 'LDAP_identity_provider' object using the configuration stored in the builder.
func (b *LDAPIdentityProviderBuilder) Build() (object *LDAPIdentityProvider, err error) {
	object = new(LDAPIdentityProvider)
	object.bitmap_ = b.bitmap_
	object.ca = b.ca
	object.url = b.url
	if b.attributes != nil {
		object.attributes, err = b.attributes.Build()
		if err != nil {
			return
		}
	}
	object.bindDN = b.bindDN
	object.bindPassword = b.bindPassword
	object.insecure = b.insecure
	return
}
