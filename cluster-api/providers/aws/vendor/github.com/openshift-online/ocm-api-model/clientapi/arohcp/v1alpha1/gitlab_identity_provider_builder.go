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

// Details for `gitlab` identity providers.
type GitlabIdentityProviderBuilder struct {
	fieldSet_    []bool
	ca           string
	url          string
	clientID     string
	clientSecret string
}

// NewGitlabIdentityProvider creates a new builder of 'gitlab_identity_provider' objects.
func NewGitlabIdentityProvider() *GitlabIdentityProviderBuilder {
	return &GitlabIdentityProviderBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GitlabIdentityProviderBuilder) Empty() bool {
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
func (b *GitlabIdentityProviderBuilder) CA(value string) *GitlabIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.ca = value
	b.fieldSet_[0] = true
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) URL(value string) *GitlabIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.url = value
	b.fieldSet_[1] = true
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) ClientID(value string) *GitlabIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.clientID = value
	b.fieldSet_[2] = true
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GitlabIdentityProviderBuilder) ClientSecret(value string) *GitlabIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.clientSecret = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GitlabIdentityProviderBuilder) Copy(object *GitlabIdentityProvider) *GitlabIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.ca = object.ca
	b.url = object.url
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	return b
}

// Build creates a 'gitlab_identity_provider' object using the configuration stored in the builder.
func (b *GitlabIdentityProviderBuilder) Build() (object *GitlabIdentityProvider, err error) {
	object = new(GitlabIdentityProvider)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.ca = b.ca
	object.url = b.url
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	return
}
