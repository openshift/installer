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

// ClusterManagementReferenceBuilder contains the data and logic needed to build 'cluster_management_reference' objects.
//
// Cluster Mgmt reference settings of the cluster.
type ClusterManagementReferenceBuilder struct {
	bitmap_   uint32
	clusterId string
	href      string
}

// NewClusterManagementReference creates a new builder of 'cluster_management_reference' objects.
func NewClusterManagementReference() *ClusterManagementReferenceBuilder {
	return &ClusterManagementReferenceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterManagementReferenceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *ClusterManagementReferenceBuilder) ClusterId(value string) *ClusterManagementReferenceBuilder {
	b.clusterId = value
	b.bitmap_ |= 1
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ClusterManagementReferenceBuilder) Href(value string) *ClusterManagementReferenceBuilder {
	b.href = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterManagementReferenceBuilder) Copy(object *ClusterManagementReference) *ClusterManagementReferenceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clusterId = object.clusterId
	b.href = object.href
	return b
}

// Build creates a 'cluster_management_reference' object using the configuration stored in the builder.
func (b *ClusterManagementReferenceBuilder) Build() (object *ClusterManagementReference, err error) {
	object = new(ClusterManagementReference)
	object.bitmap_ = b.bitmap_
	object.clusterId = b.clusterId
	object.href = b.href
	return
}
