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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// GcpDnsDomain represents configuration for Google Cloud Platform DNS domain settings
// used in cluster DNS configuration for GCP-hosted clusters.
type GcpDnsDomainBuilder struct {
	fieldSet_    []bool
	domainPrefix string
	networkId    string
	projectId    string
}

// NewGcpDnsDomain creates a new builder of 'gcp_dns_domain' objects.
func NewGcpDnsDomain() *GcpDnsDomainBuilder {
	return &GcpDnsDomainBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GcpDnsDomainBuilder) Empty() bool {
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

// DomainPrefix sets the value of the 'domain_prefix' attribute to the given value.
func (b *GcpDnsDomainBuilder) DomainPrefix(value string) *GcpDnsDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.domainPrefix = value
	b.fieldSet_[0] = true
	return b
}

// NetworkId sets the value of the 'network_id' attribute to the given value.
func (b *GcpDnsDomainBuilder) NetworkId(value string) *GcpDnsDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.networkId = value
	b.fieldSet_[1] = true
	return b
}

// ProjectId sets the value of the 'project_id' attribute to the given value.
func (b *GcpDnsDomainBuilder) ProjectId(value string) *GcpDnsDomainBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.projectId = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GcpDnsDomainBuilder) Copy(object *GcpDnsDomain) *GcpDnsDomainBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.domainPrefix = object.domainPrefix
	b.networkId = object.networkId
	b.projectId = object.projectId
	return b
}

// Build creates a 'gcp_dns_domain' object using the configuration stored in the builder.
func (b *GcpDnsDomainBuilder) Build() (object *GcpDnsDomain, err error) {
	object = new(GcpDnsDomain)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.domainPrefix = b.domainPrefix
	object.networkId = b.networkId
	object.projectId = b.projectId
	return
}
