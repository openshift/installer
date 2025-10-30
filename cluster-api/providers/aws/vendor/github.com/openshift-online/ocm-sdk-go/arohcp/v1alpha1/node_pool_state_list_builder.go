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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// NodePoolStateListBuilder contains the data and logic needed to build
// 'node_pool_state' objects.
type NodePoolStateListBuilder struct {
	items []*NodePoolStateBuilder
}

// NewNodePoolStateList creates a new builder of 'node_pool_state' objects.
func NewNodePoolStateList() *NodePoolStateListBuilder {
	return new(NodePoolStateListBuilder)
}

// Items sets the items of the list.
func (b *NodePoolStateListBuilder) Items(values ...*NodePoolStateBuilder) *NodePoolStateListBuilder {
	b.items = make([]*NodePoolStateBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *NodePoolStateListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *NodePoolStateListBuilder) Copy(list *NodePoolStateList) *NodePoolStateListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*NodePoolStateBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewNodePoolState().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'node_pool_state' objects using the
// configuration stored in the builder.
func (b *NodePoolStateListBuilder) Build() (list *NodePoolStateList, err error) {
	items := make([]*NodePoolState, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(NodePoolStateList)
	list.items = items
	return
}
