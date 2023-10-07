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

// STSCredentialRequestBuilder contains the data and logic needed to build 'STS_credential_request' objects.
//
// Representation of an credRequest
type STSCredentialRequestBuilder struct {
	bitmap_  uint32
	name     string
	operator *STSOperatorBuilder
}

// NewSTSCredentialRequest creates a new builder of 'STS_credential_request' objects.
func NewSTSCredentialRequest() *STSCredentialRequestBuilder {
	return &STSCredentialRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *STSCredentialRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *STSCredentialRequestBuilder) Name(value string) *STSCredentialRequestBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Operator sets the value of the 'operator' attribute to the given value.
//
// Representation of an sts operator
func (b *STSCredentialRequestBuilder) Operator(value *STSOperatorBuilder) *STSCredentialRequestBuilder {
	b.operator = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSCredentialRequestBuilder) Copy(object *STSCredentialRequest) *STSCredentialRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.name = b.name
	if b.operator != nil {
		object.operator, err = b.operator.Build()
		if err != nil {
			return
		}
	}
	return
}
