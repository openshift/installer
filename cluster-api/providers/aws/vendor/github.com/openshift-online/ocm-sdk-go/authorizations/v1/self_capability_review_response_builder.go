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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// SelfCapabilityReviewResponseBuilder contains the data and logic needed to build 'self_capability_review_response' objects.
//
// Representation of a capability review response.
type SelfCapabilityReviewResponseBuilder struct {
	bitmap_ uint32
	result  string
}

// NewSelfCapabilityReviewResponse creates a new builder of 'self_capability_review_response' objects.
func NewSelfCapabilityReviewResponse() *SelfCapabilityReviewResponseBuilder {
	return &SelfCapabilityReviewResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SelfCapabilityReviewResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Result sets the value of the 'result' attribute to the given value.
func (b *SelfCapabilityReviewResponseBuilder) Result(value string) *SelfCapabilityReviewResponseBuilder {
	b.result = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SelfCapabilityReviewResponseBuilder) Copy(object *SelfCapabilityReviewResponse) *SelfCapabilityReviewResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.result = object.result
	return b
}

// Build creates a 'self_capability_review_response' object using the configuration stored in the builder.
func (b *SelfCapabilityReviewResponseBuilder) Build() (object *SelfCapabilityReviewResponse, err error) {
	object = new(SelfCapabilityReviewResponse)
	object.bitmap_ = b.bitmap_
	object.result = b.result
	return
}
