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

// NodesInfoBuilder contains the data and logic needed to build 'nodes_info' objects.
//
// Provides information about the nodes in the cluster.
type NodesInfoBuilder struct {
	bitmap_ uint32
	nodes   []*NodeInfoBuilder
}

// NewNodesInfo creates a new builder of 'nodes_info' objects.
func NewNodesInfo() *NodesInfoBuilder {
	return &NodesInfoBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodesInfoBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Nodes sets the value of the 'nodes' attribute to the given values.
func (b *NodesInfoBuilder) Nodes(values ...*NodeInfoBuilder) *NodesInfoBuilder {
	b.nodes = make([]*NodeInfoBuilder, len(values))
	copy(b.nodes, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodesInfoBuilder) Copy(object *NodesInfo) *NodesInfoBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
