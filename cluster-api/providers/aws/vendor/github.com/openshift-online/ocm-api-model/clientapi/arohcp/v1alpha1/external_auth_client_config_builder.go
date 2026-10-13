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

// ExternalAuthClientConfig contains configuration for the platform's clients that
// need to request tokens from the issuer.
type ExternalAuthClientConfigBuilder struct {
	fieldSet_   []bool
	id          string
	component   *ClientComponentBuilder
	extraScopes []string
	secret      string
	type_       ExternalAuthClientType
}

// NewExternalAuthClientConfig creates a new builder of 'external_auth_client_config' objects.
func NewExternalAuthClientConfig() *ExternalAuthClientConfigBuilder {
	return &ExternalAuthClientConfigBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthClientConfigBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *ExternalAuthClientConfigBuilder) ID(value string) *ExternalAuthClientConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Component sets the value of the 'component' attribute to the given value.
//
// The reference of a component that will consume the client configuration.
func (b *ExternalAuthClientConfigBuilder) Component(value *ClientComponentBuilder) *ExternalAuthClientConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.component = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// ExtraScopes sets the value of the 'extra_scopes' attribute to the given values.
func (b *ExternalAuthClientConfigBuilder) ExtraScopes(values ...string) *ExternalAuthClientConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.extraScopes = make([]string, len(values))
	copy(b.extraScopes, values)
	b.fieldSet_[2] = true
	return b
}

// Secret sets the value of the 'secret' attribute to the given value.
func (b *ExternalAuthClientConfigBuilder) Secret(value string) *ExternalAuthClientConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.secret = value
	b.fieldSet_[3] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Representation of the possible values of an external authentication client's type
func (b *ExternalAuthClientConfigBuilder) Type(value ExternalAuthClientType) *ExternalAuthClientConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.type_ = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthClientConfigBuilder) Copy(object *ExternalAuthClientConfig) *ExternalAuthClientConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	b.type_ = object.type_
	return b
}

// Build creates a 'external_auth_client_config' object using the configuration stored in the builder.
func (b *ExternalAuthClientConfigBuilder) Build() (object *ExternalAuthClientConfig, err error) {
	object = new(ExternalAuthClientConfig)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	object.type_ = b.type_
	return
}
