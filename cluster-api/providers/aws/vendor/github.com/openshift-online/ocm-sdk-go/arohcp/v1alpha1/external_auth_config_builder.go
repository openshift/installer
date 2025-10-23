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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// ExternalAuthConfigBuilder contains the data and logic needed to build 'external_auth_config' objects.
//
// ExternalAuthConfig configuration
type ExternalAuthConfigBuilder struct {
	bitmap_       uint32
	externalAuths *v1.ExternalAuthListBuilder
	enabled       bool
}

// NewExternalAuthConfig creates a new builder of 'external_auth_config' objects.
func NewExternalAuthConfig() *ExternalAuthConfigBuilder {
	return &ExternalAuthConfigBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *ExternalAuthConfigBuilder) Enabled(value bool) *ExternalAuthConfigBuilder {
	b.enabled = value
	b.bitmap_ |= 1
	return b
}

// ExternalAuths sets the value of the 'external_auths' attribute to the given values.
func (b *ExternalAuthConfigBuilder) ExternalAuths(value *v1.ExternalAuthListBuilder) *ExternalAuthConfigBuilder {
	b.externalAuths = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthConfigBuilder) Copy(object *ExternalAuthConfig) *ExternalAuthConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.enabled = object.enabled
	if object.externalAuths != nil {
		b.externalAuths = v1.NewExternalAuthList().Copy(object.externalAuths)
	} else {
		b.externalAuths = nil
	}
	return b
}

// Build creates a 'external_auth_config' object using the configuration stored in the builder.
func (b *ExternalAuthConfigBuilder) Build() (object *ExternalAuthConfig, err error) {
	object = new(ExternalAuthConfig)
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	if b.externalAuths != nil {
		object.externalAuths, err = b.externalAuths.Build()
		if err != nil {
			return
		}
	}
	return
}
