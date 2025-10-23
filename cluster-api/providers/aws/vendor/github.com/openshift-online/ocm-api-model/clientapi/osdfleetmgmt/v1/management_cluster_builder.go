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

import (
	time "time"
)

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
type ManagementClusterBuilder struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	dns                        *DNSBuilder
	cloudProvider              string
	clusterManagementReference *ClusterManagementReferenceBuilder
	creationTimestamp          time.Time
	labels                     []*LabelBuilder
	name                       string
	parent                     *ManagementClusterParentBuilder
	region                     string
	sector                     string
	status                     string
	updateTimestamp            time.Time
}

// NewManagementCluster creates a new builder of 'management_cluster' objects.
func NewManagementCluster() *ManagementClusterBuilder {
	return &ManagementClusterBuilder{
		fieldSet_: make([]bool, 14),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ManagementClusterBuilder) Link(value bool) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ManagementClusterBuilder) ID(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ManagementClusterBuilder) HREF(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagementClusterBuilder) Empty() bool {
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

// DNS sets the value of the 'DNS' attribute to the given value.
//
// DNS settings of the cluster.
func (b *ManagementClusterBuilder) DNS(value *DNSBuilder) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.dns = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *ManagementClusterBuilder) CloudProvider(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.cloudProvider = value
	b.fieldSet_[4] = true
	return b
}

// ClusterManagementReference sets the value of the 'cluster_management_reference' attribute to the given value.
//
// Cluster Mgmt reference settings of the cluster.
func (b *ManagementClusterBuilder) ClusterManagementReference(value *ClusterManagementReferenceBuilder) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.clusterManagementReference = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ManagementClusterBuilder) CreationTimestamp(value time.Time) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.creationTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *ManagementClusterBuilder) Labels(values ...*LabelBuilder) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.labels = make([]*LabelBuilder, len(values))
	copy(b.labels, values)
	b.fieldSet_[7] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ManagementClusterBuilder) Name(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.name = value
	b.fieldSet_[8] = true
	return b
}

// Parent sets the value of the 'parent' attribute to the given value.
//
// ManagementClusterParent reference settings of the cluster.
func (b *ManagementClusterBuilder) Parent(value *ManagementClusterParentBuilder) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.parent = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
func (b *ManagementClusterBuilder) Region(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.region = value
	b.fieldSet_[10] = true
	return b
}

// Sector sets the value of the 'sector' attribute to the given value.
func (b *ManagementClusterBuilder) Sector(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.sector = value
	b.fieldSet_[11] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *ManagementClusterBuilder) Status(value string) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.status = value
	b.fieldSet_[12] = true
	return b
}

// UpdateTimestamp sets the value of the 'update_timestamp' attribute to the given value.
func (b *ManagementClusterBuilder) UpdateTimestamp(value time.Time) *ManagementClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.updateTimestamp = value
	b.fieldSet_[13] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagementClusterBuilder) Copy(object *ManagementCluster) *ManagementClusterBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.dns != nil {
		b.dns = NewDNS().Copy(object.dns)
	} else {
		b.dns = nil
	}
	b.cloudProvider = object.cloudProvider
	if object.clusterManagementReference != nil {
		b.clusterManagementReference = NewClusterManagementReference().Copy(object.clusterManagementReference)
	} else {
		b.clusterManagementReference = nil
	}
	b.creationTimestamp = object.creationTimestamp
	if object.labels != nil {
		b.labels = make([]*LabelBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabel().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.name = object.name
	if object.parent != nil {
		b.parent = NewManagementClusterParent().Copy(object.parent)
	} else {
		b.parent = nil
	}
	b.region = object.region
	b.sector = object.sector
	b.status = object.status
	b.updateTimestamp = object.updateTimestamp
	return b
}

// Build creates a 'management_cluster' object using the configuration stored in the builder.
func (b *ManagementClusterBuilder) Build() (object *ManagementCluster, err error) {
	object = new(ManagementCluster)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.dns != nil {
		object.dns, err = b.dns.Build()
		if err != nil {
			return
		}
	}
	object.cloudProvider = b.cloudProvider
	if b.clusterManagementReference != nil {
		object.clusterManagementReference, err = b.clusterManagementReference.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	if b.labels != nil {
		object.labels = make([]*Label, len(b.labels))
		for i, v := range b.labels {
			object.labels[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.name = b.name
	if b.parent != nil {
		object.parent, err = b.parent.Build()
		if err != nil {
			return
		}
	}
	object.region = b.region
	object.sector = b.sector
	object.status = b.status
	object.updateTimestamp = b.updateTimestamp
	return
}
