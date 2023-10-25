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

// CapabilityReviewResponseListBuilder contains the data and logic needed to build
// 'capability_review_response' objects.
type CapabilityReviewResponseListBuilder struct {
	items []*CapabilityReviewResponseBuilder
}

// NewCapabilityReviewResponseList creates a new builder of 'capability_review_response' objects.
func NewCapabilityReviewResponseList() *CapabilityReviewResponseListBuilder {
	return new(CapabilityReviewResponseListBuilder)
}

// Items sets the items of the list.
func (b *CapabilityReviewResponseListBuilder) Items(values ...*CapabilityReviewResponseBuilder) *CapabilityReviewResponseListBuilder {
	b.items = make([]*CapabilityReviewResponseBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *CapabilityReviewResponseListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *CapabilityReviewResponseListBuilder) Copy(list *CapabilityReviewResponseList) *CapabilityReviewResponseListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*CapabilityReviewResponseBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewCapabilityReviewResponse().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'capability_review_response' objects using the
// configuration stored in the builder.
func (b *CapabilityReviewResponseListBuilder) Build() (list *CapabilityReviewResponseList, err error) {
	items := make([]*CapabilityReviewResponse, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(CapabilityReviewResponseList)
	list.items = items
	return
}
