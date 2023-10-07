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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

// ManagementClusterParentBuilder contains the data and logic needed to build 'management_cluster_parent' objects.
//
// ManagementClusterParent reference settings of the cluster.
type ManagementClusterParentBuilder struct {
	bitmap_   uint32
	clusterId string
	href      string
	kind      string
	name      string
}

// NewManagementClusterParent creates a new builder of 'management_cluster_parent' objects.
func NewManagementClusterParent() *ManagementClusterParentBuilder {
	return &ManagementClusterParentBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagementClusterParentBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *ManagementClusterParentBuilder) ClusterId(value string) *ManagementClusterParentBuilder {
	b.clusterId = value
	b.bitmap_ |= 1
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ManagementClusterParentBuilder) Href(value string) *ManagementClusterParentBuilder {
	b.href = value
	b.bitmap_ |= 2
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *ManagementClusterParentBuilder) Kind(value string) *ManagementClusterParentBuilder {
	b.kind = value
	b.bitmap_ |= 4
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ManagementClusterParentBuilder) Name(value string) *ManagementClusterParentBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagementClusterParentBuilder) Copy(object *ManagementClusterParent) *ManagementClusterParentBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clusterId = object.clusterId
	b.href = object.href
	b.kind = object.kind
	b.name = object.name
	return b
}

// Build creates a 'management_cluster_parent' object using the configuration stored in the builder.
func (b *ManagementClusterParentBuilder) Build() (object *ManagementClusterParent, err error) {
	object = new(ManagementClusterParent)
	object.bitmap_ = b.bitmap_
	object.clusterId = b.clusterId
	object.href = b.href
	object.kind = b.kind
	object.name = b.name
	return
}
