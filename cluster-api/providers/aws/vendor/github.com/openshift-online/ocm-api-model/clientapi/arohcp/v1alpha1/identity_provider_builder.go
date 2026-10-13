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

// Representation of an identity provider.
type IdentityProviderBuilder struct {
	fieldSet_     []bool
	id            string
	href          string
	ldap          *LDAPIdentityProviderBuilder
	github        *GithubIdentityProviderBuilder
	gitlab        *GitlabIdentityProviderBuilder
	google        *GoogleIdentityProviderBuilder
	htpasswd      *HTPasswdIdentityProviderBuilder
	mappingMethod IdentityProviderMappingMethod
	name          string
	openID        *OpenIDIdentityProviderBuilder
	type_         IdentityProviderType
	challenge     bool
	login         bool
}

// NewIdentityProvider creates a new builder of 'identity_provider' objects.
func NewIdentityProvider() *IdentityProviderBuilder {
	return &IdentityProviderBuilder{
		fieldSet_: make([]bool, 14),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *IdentityProviderBuilder) Link(value bool) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *IdentityProviderBuilder) ID(value string) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *IdentityProviderBuilder) HREF(value string) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *IdentityProviderBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// LDAP sets the value of the 'LDAP' attribute to the given value.
//
// Details for `ldap` identity providers.
func (b *IdentityProviderBuilder) LDAP(value *LDAPIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.ldap = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Challenge sets the value of the 'challenge' attribute to the given value.
func (b *IdentityProviderBuilder) Challenge(value bool) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.challenge = value
	b.fieldSet_[4] = true
	return b
}

// Github sets the value of the 'github' attribute to the given value.
//
// Details for `github` identity providers.
func (b *IdentityProviderBuilder) Github(value *GithubIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.github = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Gitlab sets the value of the 'gitlab' attribute to the given value.
//
// Details for `gitlab` identity providers.
func (b *IdentityProviderBuilder) Gitlab(value *GitlabIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.gitlab = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Google sets the value of the 'google' attribute to the given value.
//
// Details for `google` identity providers.
func (b *IdentityProviderBuilder) Google(value *GoogleIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.google = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Htpasswd sets the value of the 'htpasswd' attribute to the given value.
//
// Details for `htpasswd` identity providers.
func (b *IdentityProviderBuilder) Htpasswd(value *HTPasswdIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.htpasswd = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Login sets the value of the 'login' attribute to the given value.
func (b *IdentityProviderBuilder) Login(value bool) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.login = value
	b.fieldSet_[9] = true
	return b
}

// MappingMethod sets the value of the 'mapping_method' attribute to the given value.
//
// Controls how mappings are established between provider identities and user objects.
func (b *IdentityProviderBuilder) MappingMethod(value IdentityProviderMappingMethod) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.mappingMethod = value
	b.fieldSet_[10] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *IdentityProviderBuilder) Name(value string) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.name = value
	b.fieldSet_[11] = true
	return b
}

// OpenID sets the value of the 'open_ID' attribute to the given value.
//
// Details for `openid` identity providers.
func (b *IdentityProviderBuilder) OpenID(value *OpenIDIdentityProviderBuilder) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.openID = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Type of identity provider.
func (b *IdentityProviderBuilder) Type(value IdentityProviderType) *IdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.type_ = value
	b.fieldSet_[13] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *IdentityProviderBuilder) Copy(object *IdentityProvider) *IdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.ldap != nil {
		b.ldap = NewLDAPIdentityProvider().Copy(object.ldap)
	} else {
		b.ldap = nil
	}
	b.challenge = object.challenge
	if object.github != nil {
		b.github = NewGithubIdentityProvider().Copy(object.github)
	} else {
		b.github = nil
	}
	if object.gitlab != nil {
		b.gitlab = NewGitlabIdentityProvider().Copy(object.gitlab)
	} else {
		b.gitlab = nil
	}
	if object.google != nil {
		b.google = NewGoogleIdentityProvider().Copy(object.google)
	} else {
		b.google = nil
	}
	if object.htpasswd != nil {
		b.htpasswd = NewHTPasswdIdentityProvider().Copy(object.htpasswd)
	} else {
		b.htpasswd = nil
	}
	b.login = object.login
	b.mappingMethod = object.mappingMethod
	b.name = object.name
	if object.openID != nil {
		b.openID = NewOpenIDIdentityProvider().Copy(object.openID)
	} else {
		b.openID = nil
	}
	b.type_ = object.type_
	return b
}

// Build creates a 'identity_provider' object using the configuration stored in the builder.
func (b *IdentityProviderBuilder) Build() (object *IdentityProvider, err error) {
	object = new(IdentityProvider)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.ldap != nil {
		object.ldap, err = b.ldap.Build()
		if err != nil {
			return
		}
	}
	object.challenge = b.challenge
	if b.github != nil {
		object.github, err = b.github.Build()
		if err != nil {
			return
		}
	}
	if b.gitlab != nil {
		object.gitlab, err = b.gitlab.Build()
		if err != nil {
			return
		}
	}
	if b.google != nil {
		object.google, err = b.google.Build()
		if err != nil {
			return
		}
	}
	if b.htpasswd != nil {
		object.htpasswd, err = b.htpasswd.Build()
		if err != nil {
			return
		}
	}
	object.login = b.login
	object.mappingMethod = b.mappingMethod
	object.name = b.name
	if b.openID != nil {
		object.openID, err = b.openID.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	return
}
