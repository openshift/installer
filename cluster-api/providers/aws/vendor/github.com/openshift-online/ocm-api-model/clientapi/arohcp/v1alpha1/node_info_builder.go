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

// Provides information about a node from specific type in the cluster.
type NodeInfoBuilder struct {
	fieldSet_ []bool
	amount    int
	type_     NodeType
}

// NewNodeInfo creates a new builder of 'node_info' objects.
func NewNodeInfo() *NodeInfoBuilder {
	return &NodeInfoBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodeInfoBuilder) Empty() bool {
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

// Amount sets the value of the 'amount' attribute to the given value.
func (b *NodeInfoBuilder) Amount(value int) *NodeInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.amount = value
	b.fieldSet_[0] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Type of node received via telemetry.
func (b *NodeInfoBuilder) Type(value NodeType) *NodeInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.type_ = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodeInfoBuilder) Copy(object *NodeInfo) *NodeInfoBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.amount = object.amount
	b.type_ = object.type_
	return b
}

// Build creates a 'node_info' object using the configuration stored in the builder.
func (b *NodeInfoBuilder) Build() (object *NodeInfo, err error) {
	object = new(NodeInfo)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.amount = b.amount
	object.type_ = b.type_
	return
}
