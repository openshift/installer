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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

// SelfAccessReviewResponseListBuilder contains the data and logic needed to build
// 'self_access_review_response' objects.
type SelfAccessReviewResponseListBuilder struct {
	items []*SelfAccessReviewResponseBuilder
}

// NewSelfAccessReviewResponseList creates a new builder of 'self_access_review_response' objects.
func NewSelfAccessReviewResponseList() *SelfAccessReviewResponseListBuilder {
	return new(SelfAccessReviewResponseListBuilder)
}

// Items sets the items of the list.
func (b *SelfAccessReviewResponseListBuilder) Items(values ...*SelfAccessReviewResponseBuilder) *SelfAccessReviewResponseListBuilder {
	b.items = make([]*SelfAccessReviewResponseBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *SelfAccessReviewResponseListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *SelfAccessReviewResponseListBuilder) Copy(list *SelfAccessReviewResponseList) *SelfAccessReviewResponseListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*SelfAccessReviewResponseBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewSelfAccessReviewResponse().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'self_access_review_response' objects using the
// configuration stored in the builder.
func (b *SelfAccessReviewResponseListBuilder) Build() (list *SelfAccessReviewResponseList, err error) {
	items := make([]*SelfAccessReviewResponse, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(SelfAccessReviewResponseList)
	list.items = items
	return
}
