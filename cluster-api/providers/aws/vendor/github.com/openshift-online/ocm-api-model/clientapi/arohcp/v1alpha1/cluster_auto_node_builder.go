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

// The AutoNode configuration for the Cluster.
type ClusterAutoNodeBuilder struct {
	fieldSet_ []bool
	mode      string
	status    *ClusterAutoNodeStatusBuilder
}

// NewClusterAutoNode creates a new builder of 'cluster_auto_node' objects.
func NewClusterAutoNode() *ClusterAutoNodeBuilder {
	return &ClusterAutoNodeBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAutoNodeBuilder) Empty() bool {
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

// Mode sets the value of the 'mode' attribute to the given value.
func (b *ClusterAutoNodeBuilder) Mode(value string) *ClusterAutoNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.mode = value
	b.fieldSet_[0] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Additional status information on the AutoNode configuration on this Cluster
func (b *ClusterAutoNodeBuilder) Status(value *ClusterAutoNodeStatusBuilder) *ClusterAutoNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAutoNodeBuilder) Copy(object *ClusterAutoNode) *ClusterAutoNodeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.mode = object.mode
	if object.status != nil {
		b.status = NewClusterAutoNodeStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	return b
}

// Build creates a 'cluster_auto_node' object using the configuration stored in the builder.
func (b *ClusterAutoNodeBuilder) Build() (object *ClusterAutoNode, err error) {
	object = new(ClusterAutoNode)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.mode = b.mode
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	return
}
