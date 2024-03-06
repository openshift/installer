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

// ExternalAuthClientConfigBuilder contains the data and logic needed to build 'external_auth_client_config' objects.
//
// ExternalAuthClientConfig contains configuration for the platform's clients that
// need to request tokens from the issuer.
type ExternalAuthClientConfigBuilder struct {
	bitmap_     uint32
	id          string
	component   *ClientComponentBuilder
	extraScopes []string
	secret      string
}

// NewExternalAuthClientConfig creates a new builder of 'external_auth_client_config' objects.
func NewExternalAuthClientConfig() *ExternalAuthClientConfigBuilder {
	return &ExternalAuthClientConfigBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthClientConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *ExternalAuthClientConfigBuilder) ID(value string) *ExternalAuthClientConfigBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Component sets the value of the 'component' attribute to the given value.
//
// The reference of a component that will consume the client configuration.
func (b *ExternalAuthClientConfigBuilder) Component(value *ClientComponentBuilder) *ExternalAuthClientConfigBuilder {
	b.component = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// ExtraScopes sets the value of the 'extra_scopes' attribute to the given values.
func (b *ExternalAuthClientConfigBuilder) ExtraScopes(values ...string) *ExternalAuthClientConfigBuilder {
	b.extraScopes = make([]string, len(values))
	copy(b.extraScopes, values)
	b.bitmap_ |= 4
	return b
}

// Secret sets the value of the 'secret' attribute to the given value.
func (b *ExternalAuthClientConfigBuilder) Secret(value string) *ExternalAuthClientConfigBuilder {
	b.secret = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthClientConfigBuilder) Copy(object *ExternalAuthClientConfig) *ExternalAuthClientConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	if object.component != nil {
		b.component = NewClientComponent().Copy(object.component)
	} else {
		b.component = nil
	}
	if object.extraScopes != nil {
		b.extraScopes = make([]string, len(object.extraScopes))
		copy(b.extraScopes, object.extraScopes)
	} else {
		b.extraScopes = nil
	}
	b.secret = object.secret
	return b
}

// Build creates a 'external_auth_client_config' object using the configuration stored in the builder.
func (b *ExternalAuthClientConfigBuilder) Build() (object *ExternalAuthClientConfig, err error) {
	object = new(ExternalAuthClientConfig)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	if b.component != nil {
		object.component, err = b.component.Build()
		if err != nil {
			return
		}
	}
	if b.extraScopes != nil {
		object.extraScopes = make([]string, len(b.extraScopes))
		copy(object.extraScopes, b.extraScopes)
	}
	object.secret = b.secret
	return
}
