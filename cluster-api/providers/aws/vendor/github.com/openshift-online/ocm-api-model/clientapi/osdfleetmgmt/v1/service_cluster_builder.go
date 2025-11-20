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
type ServiceClusterBuilder struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	dns                        *DNSBuilder
	cloudProvider              string
	clusterManagementReference *ClusterManagementReferenceBuilder
	labels                     []*LabelBuilder
	name                       string
	provisionShardReference    *ProvisionShardReferenceBuilder
	region                     string
	sector                     string
	status                     string
}

// NewServiceCluster creates a new builder of 'service_cluster' objects.
func NewServiceCluster() *ServiceClusterBuilder {
	return &ServiceClusterBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ServiceClusterBuilder) Link(value bool) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ServiceClusterBuilder) ID(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ServiceClusterBuilder) HREF(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceClusterBuilder) Empty() bool {
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
func (b *ServiceClusterBuilder) DNS(value *DNSBuilder) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
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
func (b *ServiceClusterBuilder) CloudProvider(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.cloudProvider = value
	b.fieldSet_[4] = true
	return b
}

// ClusterManagementReference sets the value of the 'cluster_management_reference' attribute to the given value.
//
// Cluster Mgmt reference settings of the cluster.
func (b *ServiceClusterBuilder) ClusterManagementReference(value *ClusterManagementReferenceBuilder) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.clusterManagementReference = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *ServiceClusterBuilder) Labels(values ...*LabelBuilder) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.labels = make([]*LabelBuilder, len(values))
	copy(b.labels, values)
	b.fieldSet_[6] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ServiceClusterBuilder) Name(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.name = value
	b.fieldSet_[7] = true
	return b
}

// ProvisionShardReference sets the value of the 'provision_shard_reference' attribute to the given value.
//
// Provision Shard Reference of the cluster.
func (b *ServiceClusterBuilder) ProvisionShardReference(value *ProvisionShardReferenceBuilder) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.provisionShardReference = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
func (b *ServiceClusterBuilder) Region(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.region = value
	b.fieldSet_[9] = true
	return b
}

// Sector sets the value of the 'sector' attribute to the given value.
func (b *ServiceClusterBuilder) Sector(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.sector = value
	b.fieldSet_[10] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *ServiceClusterBuilder) Status(value string) *ServiceClusterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.status = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceClusterBuilder) Copy(object *ServiceCluster) *ServiceClusterBuilder {
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
	if object.labels != nil {
		b.labels = make([]*LabelBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabel().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.name = object.name
	if object.provisionShardReference != nil {
		b.provisionShardReference = NewProvisionShardReference().Copy(object.provisionShardReference)
	} else {
		b.provisionShardReference = nil
	}
	b.region = object.region
	b.sector = object.sector
	b.status = object.status
	return b
}

// Build creates a 'service_cluster' object using the configuration stored in the builder.
func (b *ServiceClusterBuilder) Build() (object *ServiceCluster, err error) {
	object = new(ServiceCluster)
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
	if b.provisionShardReference != nil {
		object.provisionShardReference, err = b.provisionShardReference.Build()
		if err != nil {
			return
		}
	}
	object.region = b.region
	object.sector = b.sector
	object.status = b.status
	return
}
