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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

// ManagementClusterParent reference settings of the cluster.
type ManagementClusterParentBuilder struct {
	fieldSet_ []bool
	clusterId string
	href      string
	kind      string
	name      string
}

// NewManagementClusterParent creates a new builder of 'management_cluster_parent' objects.
func NewManagementClusterParent() *ManagementClusterParentBuilder {
	return &ManagementClusterParentBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagementClusterParentBuilder) Empty() bool {
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

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *ManagementClusterParentBuilder) ClusterId(value string) *ManagementClusterParentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.clusterId = value
	b.fieldSet_[0] = true
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ManagementClusterParentBuilder) Href(value string) *ManagementClusterParentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.href = value
	b.fieldSet_[1] = true
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *ManagementClusterParentBuilder) Kind(value string) *ManagementClusterParentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.kind = value
	b.fieldSet_[2] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ManagementClusterParentBuilder) Name(value string) *ManagementClusterParentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.name = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagementClusterParentBuilder) Copy(object *ManagementClusterParent) *ManagementClusterParentBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.clusterId = object.clusterId
	b.href = object.href
	b.kind = object.kind
	b.name = object.name
	return b
}

// Build creates a 'management_cluster_parent' object using the configuration stored in the builder.
func (b *ManagementClusterParentBuilder) Build() (object *ManagementClusterParent, err error) {
	object = new(ManagementClusterParent)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterId = b.clusterId
	object.href = b.href
	object.kind = b.kind
	object.name = b.name
	return
}
