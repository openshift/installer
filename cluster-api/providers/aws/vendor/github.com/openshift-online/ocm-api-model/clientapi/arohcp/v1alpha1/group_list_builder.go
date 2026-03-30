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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// GroupListBuilder contains the data and logic needed to build
// 'group' objects.
type GroupListBuilder struct {
	items []*GroupBuilder
}

// NewGroupList creates a new builder of 'group' objects.
func NewGroupList() *GroupListBuilder {
	return new(GroupListBuilder)
}

// Items sets the items of the list.
func (b *GroupListBuilder) Items(values ...*GroupBuilder) *GroupListBuilder {
	b.items = make([]*GroupBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *GroupListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *GroupListBuilder) Copy(list *GroupList) *GroupListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*GroupBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewGroup().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'group' objects using the
// configuration stored in the builder.
func (b *GroupListBuilder) Build() (list *GroupList, err error) {
	items := make([]*Group, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(GroupList)
	list.items = items
	return
}
