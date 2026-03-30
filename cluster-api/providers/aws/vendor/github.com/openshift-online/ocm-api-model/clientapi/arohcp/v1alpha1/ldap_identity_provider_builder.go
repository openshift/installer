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

// Details for `ldap` identity providers.
type LDAPIdentityProviderBuilder struct {
	fieldSet_    []bool
	ca           string
	url          string
	attributes   *LDAPAttributesBuilder
	bindDN       string
	bindPassword string
	insecure     bool
}

// NewLDAPIdentityProvider creates a new builder of 'LDAP_identity_provider' objects.
func NewLDAPIdentityProvider() *LDAPIdentityProviderBuilder {
	return &LDAPIdentityProviderBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LDAPIdentityProviderBuilder) Empty() bool {
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

// CA sets the value of the 'CA' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) CA(value string) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.ca = value
	b.fieldSet_[0] = true
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) URL(value string) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.url = value
	b.fieldSet_[1] = true
	return b
}

// Attributes sets the value of the 'attributes' attribute to the given value.
//
// LDAP attributes used to configure the LDAP identity provider.
func (b *LDAPIdentityProviderBuilder) Attributes(value *LDAPAttributesBuilder) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.attributes = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// BindDN sets the value of the 'bind_DN' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) BindDN(value string) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.bindDN = value
	b.fieldSet_[3] = true
	return b
}

// BindPassword sets the value of the 'bind_password' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) BindPassword(value string) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.bindPassword = value
	b.fieldSet_[4] = true
	return b
}

// Insecure sets the value of the 'insecure' attribute to the given value.
func (b *LDAPIdentityProviderBuilder) Insecure(value bool) *LDAPIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.insecure = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LDAPIdentityProviderBuilder) Copy(object *LDAPIdentityProvider) *LDAPIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
