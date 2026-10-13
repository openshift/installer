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

// Provides information about the nodes in the cluster.
type NodesInfoBuilder struct {
	fieldSet_ []bool
	nodes     []*NodeInfoBuilder
}

// NewNodesInfo creates a new builder of 'nodes_info' objects.
func NewNodesInfo() *NodesInfoBuilder {
	return &NodesInfoBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodesInfoBuilder) Empty() bool {
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

// Nodes sets the value of the 'nodes' attribute to the given values.
func (b *NodesInfoBuilder) Nodes(values ...*NodeInfoBuilder) *NodesInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.nodes = make([]*NodeInfoBuilder, len(values))
	copy(b.nodes, values)
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodesInfoBuilder) Copy(object *NodesInfo) *NodesInfoBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.nodes != nil {
		b.nodes = make([]*NodeInfoBuilder, len(object.nodes))
		for i, v := range object.nodes {
			b.nodes[i] = NewNodeInfo().Copy(v)
		}
	} else {
		b.nodes = nil
	}
	return b
}

// Build creates a 'nodes_info' object using the configuration stored in the builder.
func (b *NodesInfoBuilder) Build() (object *NodesInfo, err error) {
	object = new(NodesInfo)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.nodes != nil {
		object.nodes = make([]*NodeInfo, len(b.nodes))
		for i, v := range b.nodes {
			object.nodes[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
