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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Represents an external authentication configuration
type ExternalAuthConfigBuilder struct {
	fieldSet_     []bool
	id            string
	href          string
	externalAuths *ExternalAuthListBuilder
	state         ExternalAuthConfigState
	enabled       bool
}

// NewExternalAuthConfig creates a new builder of 'external_auth_config' objects.
func NewExternalAuthConfig() *ExternalAuthConfigBuilder {
	return &ExternalAuthConfigBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ExternalAuthConfigBuilder) Link(value bool) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ExternalAuthConfigBuilder) ID(value string) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ExternalAuthConfigBuilder) HREF(value string) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthConfigBuilder) Empty() bool {
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

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *ExternalAuthConfigBuilder) Enabled(value bool) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.enabled = value
	b.fieldSet_[3] = true
	return b
}

// ExternalAuths sets the value of the 'external_auths' attribute to the given values.
func (b *ExternalAuthConfigBuilder) ExternalAuths(value *ExternalAuthListBuilder) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.externalAuths = value
	b.fieldSet_[4] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of the possible values for the state field of an external authentication configuration
func (b *ExternalAuthConfigBuilder) State(value ExternalAuthConfigState) *ExternalAuthConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.state = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthConfigBuilder) Copy(object *ExternalAuthConfig) *ExternalAuthConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.enabled = object.enabled
	if object.externalAuths != nil {
		b.externalAuths = NewExternalAuthList().Copy(object.externalAuths)
	} else {
		b.externalAuths = nil
	}
	b.state = object.state
	return b
}

// Build creates a 'external_auth_config' object using the configuration stored in the builder.
func (b *ExternalAuthConfigBuilder) Build() (object *ExternalAuthConfig, err error) {
	object = new(ExternalAuthConfig)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	if b.externalAuths != nil {
		object.externalAuths, err = b.externalAuths.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	return
}
