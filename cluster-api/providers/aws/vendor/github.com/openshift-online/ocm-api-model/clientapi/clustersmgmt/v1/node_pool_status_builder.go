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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of the status of a node pool.
type NodePoolStatusBuilder struct {
	fieldSet_       []bool
	id              string
	href            string
	currentReplicas int
	message         string
	state           *NodePoolStateBuilder
}

// NewNodePoolStatus creates a new builder of 'node_pool_status' objects.
func NewNodePoolStatus() *NodePoolStatusBuilder {
	return &NodePoolStatusBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolStatusBuilder) Link(value bool) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolStatusBuilder) ID(value string) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *NodePoolStatusBuilder) HREF(value string) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolStatusBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// CurrentReplicas sets the value of the 'current_replicas' attribute to the given value.
func (b *NodePoolStatusBuilder) CurrentReplicas(value int) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.currentReplicas = value
	b.fieldSet_[3] = true
	return b
}

// Message sets the value of the 'message' attribute to the given value.
func (b *NodePoolStatusBuilder) Message(value string) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.message = value
	b.fieldSet_[4] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of the status of a node pool.
func (b *NodePoolStatusBuilder) State(value *NodePoolStateBuilder) *NodePoolStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.state = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolStatusBuilder) Copy(object *NodePoolStatus) *NodePoolStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.currentReplicas = object.currentReplicas
	b.message = object.message
	if object.state != nil {
		b.state = NewNodePoolState().Copy(object.state)
	} else {
		b.state = nil
	}
	return b
}

// Build creates a 'node_pool_status' object using the configuration stored in the builder.
func (b *NodePoolStatusBuilder) Build() (object *NodePoolStatus, err error) {
	object = new(NodePoolStatus)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.currentReplicas = b.currentReplicas
	object.message = b.message
	if b.state != nil {
		object.state, err = b.state.Build()
		if err != nil {
			return
		}
	}
	return
}
