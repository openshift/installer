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

// Representation of an credRequest
type STSCredentialRequestBuilder struct {
	fieldSet_ []bool
	name      string
	operator  *STSOperatorBuilder
}

// NewSTSCredentialRequest creates a new builder of 'STS_credential_request' objects.
func NewSTSCredentialRequest() *STSCredentialRequestBuilder {
	return &STSCredentialRequestBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *STSCredentialRequestBuilder) Empty() bool {
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

// Name sets the value of the 'name' attribute to the given value.
func (b *STSCredentialRequestBuilder) Name(value string) *STSCredentialRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.name = value
	b.fieldSet_[0] = true
	return b
}

// Operator sets the value of the 'operator' attribute to the given value.
//
// Representation of an sts operator
func (b *STSCredentialRequestBuilder) Operator(value *STSOperatorBuilder) *STSCredentialRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.operator = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSCredentialRequestBuilder) Copy(object *STSCredentialRequest) *STSCredentialRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.name = object.name
	if object.operator != nil {
		b.operator = NewSTSOperator().Copy(object.operator)
	} else {
		b.operator = nil
	}
	return b
}

// Build creates a 'STS_credential_request' object using the configuration stored in the builder.
func (b *STSCredentialRequestBuilder) Build() (object *STSCredentialRequest, err error) {
	object = new(STSCredentialRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	if b.operator != nil {
		object.operator, err = b.operator.Build()
		if err != nil {
			return
		}
	}
	return
}
