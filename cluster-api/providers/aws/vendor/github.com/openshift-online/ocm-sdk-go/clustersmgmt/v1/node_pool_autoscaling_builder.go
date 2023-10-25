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

// NodePoolAutoscalingBuilder contains the data and logic needed to build 'node_pool_autoscaling' objects.
//
// Representation of a autoscaling in a node pool.
type NodePoolAutoscalingBuilder struct {
	bitmap_    uint32
	id         string
	href       string
	maxReplica int
	minReplica int
}

// NewNodePoolAutoscaling creates a new builder of 'node_pool_autoscaling' objects.
func NewNodePoolAutoscaling() *NodePoolAutoscalingBuilder {
	return &NodePoolAutoscalingBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolAutoscalingBuilder) Link(value bool) *NodePoolAutoscalingBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolAutoscalingBuilder) ID(value string) *NodePoolAutoscalingBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *NodePoolAutoscalingBuilder) HREF(value string) *NodePoolAutoscalingBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolAutoscalingBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// MaxReplica sets the value of the 'max_replica' attribute to the given value.
func (b *NodePoolAutoscalingBuilder) MaxReplica(value int) *NodePoolAutoscalingBuilder {
	b.maxReplica = value
	b.bitmap_ |= 8
	return b
}

// MinReplica sets the value of the 'min_replica' attribute to the given value.
func (b *NodePoolAutoscalingBuilder) MinReplica(value int) *NodePoolAutoscalingBuilder {
	b.minReplica = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolAutoscalingBuilder) Copy(object *NodePoolAutoscaling) *NodePoolAutoscalingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.maxReplica = object.maxReplica
	b.minReplica = object.minReplica
	return b
}

// Build creates a 'node_pool_autoscaling' object using the configuration stored in the builder.
func (b *NodePoolAutoscalingBuilder) Build() (object *NodePoolAutoscaling, err error) {
	object = new(NodePoolAutoscaling)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.maxReplica = b.maxReplica
	object.minReplica = b.minReplica
	return
}
