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

// IdentityProviderBuilder contains the data and logic needed to build 'identity_provider' objects.
//
// Representation of an identity provider.
type IdentityProviderBuilder struct {
	bitmap_       uint32
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
	return &IdentityProviderBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *IdentityProviderBuilder) Link(value bool) *IdentityProviderBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *IdentityProviderBuilder) ID(value string) *IdentityProviderBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *IdentityProviderBuilder) HREF(value string) *IdentityProviderBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *IdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// LDAP sets the value of the 'LDAP' attribute to the given value.
//
// Details for `ldap` identity providers.
func (b *IdentityProviderBuilder) LDAP(value *LDAPIdentityProviderBuilder) *IdentityProviderBuilder {
	b.ldap = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Challenge sets the value of the 'challenge' attribute to the given value.
func (b *IdentityProviderBuilder) Challenge(value bool) *IdentityProviderBuilder {
	b.challenge = value
	b.bitmap_ |= 16
	return b
}

// Github sets the value of the 'github' attribute to the given value.
//
// Details for `github` identity providers.
func (b *IdentityProviderBuilder) Github(value *GithubIdentityProviderBuilder) *IdentityProviderBuilder {
	b.github = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Gitlab sets the value of the 'gitlab' attribute to the given value.
//
// Details for `gitlab` identity providers.
func (b *IdentityProviderBuilder) Gitlab(value *GitlabIdentityProviderBuilder) *IdentityProviderBuilder {
	b.gitlab = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Google sets the value of the 'google' attribute to the given value.
//
// Details for `google` identity providers.
func (b *IdentityProviderBuilder) Google(value *GoogleIdentityProviderBuilder) *IdentityProviderBuilder {
	b.google = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// Htpasswd sets the value of the 'htpasswd' attribute to the given value.
//
// Details for `htpasswd` identity providers.
func (b *IdentityProviderBuilder) Htpasswd(value *HTPasswdIdentityProviderBuilder) *IdentityProviderBuilder {
	b.htpasswd = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// Login sets the value of the 'login' attribute to the given value.
func (b *IdentityProviderBuilder) Login(value bool) *IdentityProviderBuilder {
	b.login = value
	b.bitmap_ |= 512
	return b
}

// MappingMethod sets the value of the 'mapping_method' attribute to the given value.
//
// Controls how mappings are established between provider identities and user objects.
func (b *IdentityProviderBuilder) MappingMethod(value IdentityProviderMappingMethod) *IdentityProviderBuilder {
	b.mappingMethod = value
	b.bitmap_ |= 1024
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *IdentityProviderBuilder) Name(value string) *IdentityProviderBuilder {
	b.name = value
	b.bitmap_ |= 2048
	return b
}

// OpenID sets the value of the 'open_ID' attribute to the given value.
//
// Details for `openid` identity providers.
func (b *IdentityProviderBuilder) OpenID(value *OpenIDIdentityProviderBuilder) *IdentityProviderBuilder {
	b.openID = value
	if value != nil {
		b.bitmap_ |= 4096
	} else {
		b.bitmap_ &^= 4096
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Type of identity provider.
func (b *IdentityProviderBuilder) Type(value IdentityProviderType) *IdentityProviderBuilder {
	b.type_ = value
	b.bitmap_ |= 8192
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *IdentityProviderBuilder) Copy(object *IdentityProvider) *IdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
