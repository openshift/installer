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

// SelfCapabilityReviewRequestListBuilder contains the data and logic needed to build
// 'self_capability_review_request' objects.
type SelfCapabilityReviewRequestListBuilder struct {
	items []*SelfCapabilityReviewRequestBuilder
}

// NewSelfCapabilityReviewRequestList creates a new builder of 'self_capability_review_request' objects.
func NewSelfCapabilityReviewRequestList() *SelfCapabilityReviewRequestListBuilder {
	return new(SelfCapabilityReviewRequestListBuilder)
}

// Items sets the items of the list.
func (b *SelfCapabilityReviewRequestListBuilder) Items(values ...*SelfCapabilityReviewRequestBuilder) *SelfCapabilityReviewRequestListBuilder {
	b.items = make([]*SelfCapabilityReviewRequestBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *SelfCapabilityReviewRequestListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *SelfCapabilityReviewRequestListBuilder) Copy(list *SelfCapabilityReviewRequestList) *SelfCapabilityReviewRequestListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*SelfCapabilityReviewRequestBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewSelfCapabilityReviewRequest().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'self_capability_review_request' objects using the
// configuration stored in the builder.
func (b *SelfCapabilityReviewRequestListBuilder) Build() (list *SelfCapabilityReviewRequestList, err error) {
	items := make([]*SelfCapabilityReviewRequest, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(SelfCapabilityReviewRequestList)
	list.items = items
	return
}
