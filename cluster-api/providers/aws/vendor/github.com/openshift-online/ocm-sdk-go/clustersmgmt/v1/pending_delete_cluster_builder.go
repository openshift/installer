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

// PendingDeleteClusterBuilder contains the data and logic needed to build 'pending_delete_cluster' objects.
//
// Represents a pending delete entry for a specific cluster.
type PendingDeleteClusterBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	cluster           *ClusterBuilder
	creationTimestamp time.Time
	bestEffort        bool
}

// NewPendingDeleteCluster creates a new builder of 'pending_delete_cluster' objects.
func NewPendingDeleteCluster() *PendingDeleteClusterBuilder {
	return &PendingDeleteClusterBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *PendingDeleteClusterBuilder) Link(value bool) *PendingDeleteClusterBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *PendingDeleteClusterBuilder) ID(value string) *PendingDeleteClusterBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *PendingDeleteClusterBuilder) HREF(value string) *PendingDeleteClusterBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PendingDeleteClusterBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BestEffort sets the value of the 'best_effort' attribute to the given value.
func (b *PendingDeleteClusterBuilder) BestEffort(value bool) *PendingDeleteClusterBuilder {
	b.bestEffort = value
	b.bitmap_ |= 8
	return b
}

// Cluster sets the value of the 'cluster' attribute to the given value.
//
// Definition of an _OpenShift_ cluster.
//
// The `cloud_provider` attribute is a reference to the cloud provider. When a
// cluster is retrieved it will be a link to the cloud provider, containing only
// the kind, id and href attributes:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "kind": "CloudProviderLink",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/cloud_providers/123"
//	  }
//	}
//
// ```
//
// When a cluster is created this is optional, and if used it should contain the
// identifier of the cloud provider to use:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "id": "123",
//	  }
//	}
//
// ```
//
// If not included, then the cluster will be created using the default cloud
// provider, which is currently Amazon Web Services.
//
// The region attribute is mandatory when a cluster is created.
//
// The `aws.access_key_id`, `aws.secret_access_key` and `dns.base_domain`
// attributes are mandatory when creation a cluster with your own Amazon Web
// Services account.
func (b *PendingDeleteClusterBuilder) Cluster(value *ClusterBuilder) *PendingDeleteClusterBuilder {
	b.cluster = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *PendingDeleteClusterBuilder) CreationTimestamp(value time.Time) *PendingDeleteClusterBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PendingDeleteClusterBuilder) Copy(object *PendingDeleteCluster) *PendingDeleteClusterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.bestEffort = object.bestEffort
	if object.cluster != nil {
		b.cluster = NewCluster().Copy(object.cluster)
	} else {
		b.cluster = nil
	}
	b.creationTimestamp = object.creationTimestamp
	return b
}

// Build creates a 'pending_delete_cluster' object using the configuration stored in the builder.
func (b *PendingDeleteClusterBuilder) Build() (object *PendingDeleteCluster, err error) {
	object = new(PendingDeleteCluster)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.bestEffort = b.bestEffort
	if b.cluster != nil {
		object.cluster, err = b.cluster.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	return
}
