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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// ClusterNodesBuilder contains the data and logic needed to build 'cluster_nodes' objects.
type ClusterNodesBuilder struct {
	bitmap_           uint32
	availabilityZones []string
}

// NewClusterNodes creates a new builder of 'cluster_nodes' objects.
func NewClusterNodes() *ClusterNodesBuilder {
	return &ClusterNodesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterNodesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AvailabilityZones sets the value of the 'availability_zones' attribute to the given values.
func (b *ClusterNodesBuilder) AvailabilityZones(values ...string) *ClusterNodesBuilder {
	b.availabilityZones = make([]string, len(values))
	copy(b.availabilityZones, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterNodesBuilder) Copy(object *ClusterNodes) *ClusterNodesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.availabilityZones != nil {
		b.availabilityZones = make([]string, len(object.availabilityZones))
		copy(b.availabilityZones, object.availabilityZones)
	} else {
		b.availabilityZones = nil
	}
	return b
}

// Build creates a 'cluster_nodes' object using the configuration stored in the builder.
func (b *ClusterNodesBuilder) Build() (object *ClusterNodes, err error) {
	object = new(ClusterNodes)
	object.bitmap_ = b.bitmap_
	if b.availabilityZones != nil {
		object.availabilityZones = make([]string, len(b.availabilityZones))
		copy(object.availabilityZones, b.availabilityZones)
	}
	return
}
