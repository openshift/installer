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

// Representation of the status of an external authentication provider.
type ExternalAuthStatusBuilder struct {
	fieldSet_ []bool
	message   string
	state     *ExternalAuthStateBuilder
}

// NewExternalAuthStatus creates a new builder of 'external_auth_status' objects.
func NewExternalAuthStatus() *ExternalAuthStatusBuilder {
	return &ExternalAuthStatusBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthStatusBuilder) Empty() bool {
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

// Message sets the value of the 'message' attribute to the given value.
func (b *ExternalAuthStatusBuilder) Message(value string) *ExternalAuthStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.message = value
	b.fieldSet_[0] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of the state of an external authentication provider.
func (b *ExternalAuthStatusBuilder) State(value *ExternalAuthStateBuilder) *ExternalAuthStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.state = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthStatusBuilder) Copy(object *ExternalAuthStatus) *ExternalAuthStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.message = object.message
	if object.state != nil {
		b.state = NewExternalAuthState().Copy(object.state)
	} else {
		b.state = nil
	}
	return b
}

// Build creates a 'external_auth_status' object using the configuration stored in the builder.
func (b *ExternalAuthStatusBuilder) Build() (object *ExternalAuthStatus, err error) {
	object = new(ExternalAuthStatus)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.message = b.message
	if b.state != nil {
		object.state, err = b.state.Build()
		if err != nil {
			return
		}
	}
	return
}
