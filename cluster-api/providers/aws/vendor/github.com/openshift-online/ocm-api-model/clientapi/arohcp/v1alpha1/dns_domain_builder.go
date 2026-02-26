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

// Contains the properties of a DNS domain.
type DNSDomainBuilder struct {
	fieldSet_           []bool
	id                  string
	href                string
	cluster             *ClusterLinkBuilder
	clusterArch         ClusterArchitecture
	organization        *OrganizationLinkBuilder
	reservedAtTimestamp time.Time
	userDefined         bool
}

// NewDNSDomain creates a new builder of 'DNS_domain' objects.
func NewDNSDomain() *DNSDomainBuilder {
	return &DNSDomainBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *DNSDomainBuilder) Link(value bool) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *DNSDomainBuilder) ID(value string) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *DNSDomainBuilder) HREF(value string) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DNSDomainBuilder) Empty() bool {
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

// Cluster sets the value of the 'cluster' attribute to the given value.
//
// Definition of a cluster link.
func (b *DNSDomainBuilder) Cluster(value *ClusterLinkBuilder) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.cluster = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// ClusterArch sets the value of the 'cluster_arch' attribute to the given value.
//
// Possible cluster architectures.
func (b *DNSDomainBuilder) ClusterArch(value ClusterArchitecture) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.clusterArch = value
	b.fieldSet_[4] = true
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
//
// Definition of an organization link.
func (b *DNSDomainBuilder) Organization(value *OrganizationLinkBuilder) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.organization = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// ReservedAtTimestamp sets the value of the 'reserved_at_timestamp' attribute to the given value.
func (b *DNSDomainBuilder) ReservedAtTimestamp(value time.Time) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.reservedAtTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// UserDefined sets the value of the 'user_defined' attribute to the given value.
func (b *DNSDomainBuilder) UserDefined(value bool) *DNSDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.userDefined = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DNSDomainBuilder) Copy(object *DNSDomain) *DNSDomainBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.cluster != nil {
		b.cluster = NewClusterLink().Copy(object.cluster)
	} else {
		b.cluster = nil
	}
	b.clusterArch = object.clusterArch
	if object.organization != nil {
		b.organization = NewOrganizationLink().Copy(object.organization)
	} else {
		b.organization = nil
	}
	b.reservedAtTimestamp = object.reservedAtTimestamp
	b.userDefined = object.userDefined
	return b
}

// Build creates a 'DNS_domain' object using the configuration stored in the builder.
func (b *DNSDomainBuilder) Build() (object *DNSDomain, err error) {
	object = new(DNSDomain)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.cluster != nil {
		object.cluster, err = b.cluster.Build()
		if err != nil {
			return
		}
	}
	object.clusterArch = b.clusterArch
	if b.organization != nil {
		object.organization, err = b.organization.Build()
		if err != nil {
			return
		}
	}
	object.reservedAtTimestamp = b.reservedAtTimestamp
	object.userDefined = b.userDefined
	return
}
