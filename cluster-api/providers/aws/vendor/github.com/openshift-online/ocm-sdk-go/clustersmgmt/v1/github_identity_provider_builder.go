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

// GithubIdentityProviderBuilder contains the data and logic needed to build 'github_identity_provider' objects.
//
// Details for `github` identity providers.
type GithubIdentityProviderBuilder struct {
	bitmap_       uint32
	ca            string
	clientID      string
	clientSecret  string
	hostname      string
	organizations []string
	teams         []string
}

// NewGithubIdentityProvider creates a new builder of 'github_identity_provider' objects.
func NewGithubIdentityProvider() *GithubIdentityProviderBuilder {
	return &GithubIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GithubIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CA sets the value of the 'CA' attribute to the given value.
func (b *GithubIdentityProviderBuilder) CA(value string) *GithubIdentityProviderBuilder {
	b.ca = value
	b.bitmap_ |= 1
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GithubIdentityProviderBuilder) ClientID(value string) *GithubIdentityProviderBuilder {
	b.clientID = value
	b.bitmap_ |= 2
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GithubIdentityProviderBuilder) ClientSecret(value string) *GithubIdentityProviderBuilder {
	b.clientSecret = value
	b.bitmap_ |= 4
	return b
}

// Hostname sets the value of the 'hostname' attribute to the given value.
func (b *GithubIdentityProviderBuilder) Hostname(value string) *GithubIdentityProviderBuilder {
	b.hostname = value
	b.bitmap_ |= 8
	return b
}

// Organizations sets the value of the 'organizations' attribute to the given values.
func (b *GithubIdentityProviderBuilder) Organizations(values ...string) *GithubIdentityProviderBuilder {
	b.organizations = make([]string, len(values))
	copy(b.organizations, values)
	b.bitmap_ |= 16
	return b
}

// Teams sets the value of the 'teams' attribute to the given values.
func (b *GithubIdentityProviderBuilder) Teams(values ...string) *GithubIdentityProviderBuilder {
	b.teams = make([]string, len(values))
	copy(b.teams, values)
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GithubIdentityProviderBuilder) Copy(object *GithubIdentityProvider) *GithubIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
