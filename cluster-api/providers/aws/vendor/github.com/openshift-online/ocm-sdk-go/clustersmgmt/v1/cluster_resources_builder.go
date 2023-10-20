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

import (
	time "time"
)

// ClusterResourcesBuilder contains the data and logic needed to build 'cluster_resources' objects.
//
// Cluster Resource which belongs to a cluster, example Cluster Deployment.
type ClusterResourcesBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	clusterID         string
	creationTimestamp time.Time
	resources         map[string]string
}

// NewClusterResources creates a new builder of 'cluster_resources' objects.
func NewClusterResources() *ClusterResourcesBuilder {
	return &ClusterResourcesBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterResourcesBuilder) Link(value bool) *ClusterResourcesBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ClusterResourcesBuilder) ID(value string) *ClusterResourcesBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ClusterResourcesBuilder) HREF(value string) *ClusterResourcesBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterResourcesBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterResourcesBuilder) ClusterID(value string) *ClusterResourcesBuilder {
	b.clusterID = value
	b.bitmap_ |= 8
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ClusterResourcesBuilder) CreationTimestamp(value time.Time) *ClusterResourcesBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 16
	return b
}

// Resources sets the value of the 'resources' attribute to the given value.
func (b *ClusterResourcesBuilder) Resources(value map[string]string) *ClusterResourcesBuilder {
	b.resources = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterResourcesBuilder) Copy(object *ClusterResources) *ClusterResourcesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.creationTimestamp = object.creationTimestamp
	if len(object.resources) > 0 {
		b.resources = map[string]string{}
		for k, v := range object.resources {
			b.resources[k] = v
		}
	} else {
		b.resources = nil
	}
	return b
}

// Build creates a 'cluster_resources' object using the configuration stored in the builder.
func (b *ClusterResourcesBuilder) Build() (object *ClusterResources, err error) {
	object = new(ClusterResources)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.clusterID = b.clusterID
	object.creationTimestamp = b.creationTimestamp
	if b.resources != nil {
		object.resources = make(map[string]string)
		for k, v := range b.resources {
			object.resources[k] = v
		}
	}
	return
}
