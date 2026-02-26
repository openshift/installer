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

import (
	time "time"
)

// Cluster Resource which belongs to a cluster, example Cluster Deployment.
type ClusterResourcesBuilder struct {
	fieldSet_         []bool
	id                string
	href              string
	clusterID         string
	creationTimestamp time.Time
	resources         map[string]string
}

// NewClusterResources creates a new builder of 'cluster_resources' objects.
func NewClusterResources() *ClusterResourcesBuilder {
	return &ClusterResourcesBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterResourcesBuilder) Link(value bool) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ClusterResourcesBuilder) ID(value string) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ClusterResourcesBuilder) HREF(value string) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterResourcesBuilder) Empty() bool {
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

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterResourcesBuilder) ClusterID(value string) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ClusterResourcesBuilder) CreationTimestamp(value time.Time) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.creationTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// Resources sets the value of the 'resources' attribute to the given value.
func (b *ClusterResourcesBuilder) Resources(value map[string]string) *ClusterResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.resources = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterResourcesBuilder) Copy(object *ClusterResources) *ClusterResourcesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
