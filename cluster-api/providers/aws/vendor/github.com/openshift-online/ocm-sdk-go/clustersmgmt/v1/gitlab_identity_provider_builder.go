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

// GitlabIdentityProviderBuilder contains the data and logic needed to build 'gitlab_identity_provider' objects.
//
// Details for `gitlab` identity providers.
type GitlabIdentityProviderBuilder struct {
	bitmap_      uint32
	ca           string
	url          string
	clientID     string
	clientSecret string
}

// NewGitlabIdentityProvider creates a new builder of 'gitlab_identity_provider' objects.
func NewGitlabIdentityProvider() *GitlabIdentityProviderBuilder {
	return &GitlabIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GitlabIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CA sets the value of the 'CA' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) CA(value string) *GitlabIdentityProviderBuilder {
	b.ca = value
	b.bitmap_ |= 1
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) URL(value string) *GitlabIdentityProviderBuilder {
	b.url = value
	b.bitmap_ |= 2
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) ClientID(value string) *GitlabIdentityProviderBuilder {
	b.clientID = value
	b.bitmap_ |= 4
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) ClientSecret(value string) *GitlabIdentityProviderBuilder {
	b.clientSecret = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GitlabIdentityProviderBuilder) Copy(object *GitlabIdentityProvider) *GitlabIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.ca = object.ca
	b.url = object.url
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	return b
}

// Build creates a 'gitlab_identity_provider' object using the configuration stored in the builder.
func (b *GitlabIdentityProviderBuilder) Build() (object *GitlabIdentityProvider, err error) {
	object = new(GitlabIdentityProvider)
	object.bitmap_ = b.bitmap_
	object.ca = b.ca
	object.url = b.url
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	return
}
