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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// FeatureToggleQueryRequestListBuilder contains the data and logic needed to build
// 'feature_toggle_query_request' objects.
type FeatureToggleQueryRequestListBuilder struct {
	items []*FeatureToggleQueryRequestBuilder
}

// NewFeatureToggleQueryRequestList creates a new builder of 'feature_toggle_query_request' objects.
func NewFeatureToggleQueryRequestList() *FeatureToggleQueryRequestListBuilder {
	return new(FeatureToggleQueryRequestListBuilder)
}

// Items sets the items of the list.
func (b *FeatureToggleQueryRequestListBuilder) Items(values ...*FeatureToggleQueryRequestBuilder) *FeatureToggleQueryRequestListBuilder {
	b.items = make([]*FeatureToggleQueryRequestBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *FeatureToggleQueryRequestListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *FeatureToggleQueryRequestListBuilder) Copy(list *FeatureToggleQueryRequestList) *FeatureToggleQueryRequestListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*FeatureToggleQueryRequestBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewFeatureToggleQueryRequest().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'feature_toggle_query_request' objects using the
// configuration stored in the builder.
func (b *FeatureToggleQueryRequestListBuilder) Build() (list *FeatureToggleQueryRequestList, err error) {
	items := make([]*FeatureToggleQueryRequest, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(FeatureToggleQueryRequestList)
	list.items = items
	return
}
