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

// Details for `github` identity providers.
type GithubIdentityProviderBuilder struct {
	fieldSet_     []bool
	ca            string
	clientID      string
	clientSecret  string
	hostname      string
	organizations []string
	teams         []string
}

// NewGithubIdentityProvider creates a new builder of 'github_identity_provider' objects.
func NewGithubIdentityProvider() *GithubIdentityProviderBuilder {
	return &GithubIdentityProviderBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GithubIdentityProviderBuilder) Empty() bool {
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
func (b *GithubIdentityProviderBuilder) CA(value string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.ca = value
	b.fieldSet_[0] = true
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GithubIdentityProviderBuilder) ClientID(value string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.clientID = value
	b.fieldSet_[1] = true
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GithubIdentityProviderBuilder) ClientSecret(value string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.clientSecret = value
	b.fieldSet_[2] = true
	return b
}

// Hostname sets the value of the 'hostname' attribute to the given value.
func (b *GithubIdentityProviderBuilder) Hostname(value string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.hostname = value
	b.fieldSet_[3] = true
	return b
}

// Organizations sets the value of the 'organizations' attribute to the given values.
func (b *GithubIdentityProviderBuilder) Organizations(values ...string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.organizations = make([]string, len(values))
	copy(b.organizations, values)
	b.fieldSet_[4] = true
	return b
}

// Teams sets the value of the 'teams' attribute to the given values.
func (b *GithubIdentityProviderBuilder) Teams(values ...string) *GithubIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.teams = make([]string, len(values))
	copy(b.teams, values)
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GithubIdentityProviderBuilder) Copy(object *GithubIdentityProvider) *GithubIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.ca = object.ca
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	b.hostname = object.hostname
	if object.organizations != nil {
		b.organizations = make([]string, len(object.organizations))
		copy(b.organizations, object.organizations)
	} else {
		b.organizations = nil
	}
	if object.teams != nil {
		b.teams = make([]string, len(object.teams))
		copy(b.teams, object.teams)
	} else {
		b.teams = nil
	}
	return b
}

// Build creates a 'github_identity_provider' object using the configuration stored in the builder.
func (b *GithubIdentityProviderBuilder) Build() (object *GithubIdentityProvider, err error) {
	object = new(GithubIdentityProvider)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.ca = b.ca
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	object.hostname = b.hostname
	if b.organizations != nil {
		object.organizations = make([]string, len(b.organizations))
		copy(object.organizations, b.organizations)
	}
	if b.teams != nil {
		object.teams = make([]string, len(b.teams))
		copy(object.teams, b.teams)
	}
	return
}
