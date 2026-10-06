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

// GcpDnsDomain represents the values of the 'gcp_dns_domain' type.
//
// GcpDnsDomain represents configuration for Google Cloud Platform DNS domain settings
// used in cluster DNS configuration for GCP-hosted clusters.
type GcpDnsDomain struct {
	fieldSet_    []bool
	domainPrefix string
	networkId    string
	projectId    string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GcpDnsDomain) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// DomainPrefix returns the value of the 'domain_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DomainPrefix specifies the DNS domain prefix that will be used for cluster DNS resolution.
// This prefix is combined with the base domain to form the full DNS domain for the cluster.
func (o *GcpDnsDomain) DomainPrefix() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.domainPrefix
	}
	return ""
}

// GetDomainPrefix returns the value of the 'domain_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// DomainPrefix specifies the DNS domain prefix that will be used for cluster DNS resolution.
// This prefix is combined with the base domain to form the full DNS domain for the cluster.
func (o *GcpDnsDomain) GetDomainPrefix() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.domainPrefix
	}
	return
}

// NetworkId returns the value of the 'network_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// NetworkId is the GCP VPC Network identifier where the DNS configuration will be applied.
func (o *GcpDnsDomain) NetworkId() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.networkId
	}
	return ""
}

// GetNetworkId returns the value of the 'network_id' attribute and
// a flag indicating if the attribute has a value.
//
// NetworkId is the GCP VPC Network identifier where the DNS configuration will be applied.
func (o *GcpDnsDomain) GetNetworkId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.networkId
	}
	return
}

// ProjectId returns the value of the 'project_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ProjectId is the Google Cloud Platform project identifier where the DNS zone is hosted.
func (o *GcpDnsDomain) ProjectId() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.projectId
	}
	return ""
}

// GetProjectId returns the value of the 'project_id' attribute and
// a flag indicating if the attribute has a value.
//
// ProjectId is the Google Cloud Platform project identifier where the DNS zone is hosted.
func (o *GcpDnsDomain) GetProjectId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.projectId
	}
	return
}

// GcpDnsDomainListKind is the name of the type used to represent list of objects of
// type 'gcp_dns_domain'.
const GcpDnsDomainListKind = "GcpDnsDomainList"

// GcpDnsDomainListLinkKind is the name of the type used to represent links to list
// of objects of type 'gcp_dns_domain'.
const GcpDnsDomainListLinkKind = "GcpDnsDomainListLink"

// GcpDnsDomainNilKind is the name of the type used to nil lists of objects of
// type 'gcp_dns_domain'.
const GcpDnsDomainListNilKind = "GcpDnsDomainListNil"

// GcpDnsDomainList is a list of values of the 'gcp_dns_domain' type.
type GcpDnsDomainList struct {
	href  string
	link  bool
	items []*GcpDnsDomain
}

// Len returns the length of the list.
func (l *GcpDnsDomainList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GcpDnsDomainList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GcpDnsDomainList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GcpDnsDomainList) SetItems(items []*GcpDnsDomain) {
	l.items = items
}

// Items returns the items of the list.
func (l *GcpDnsDomainList) Items() []*GcpDnsDomain {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GcpDnsDomainList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GcpDnsDomainList) Get(i int) *GcpDnsDomain {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *GcpDnsDomainList) Slice() []*GcpDnsDomain {
	var slice []*GcpDnsDomain
	if l == nil {
		slice = make([]*GcpDnsDomain, 0)
	} else {
		slice = make([]*GcpDnsDomain, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GcpDnsDomainList) Each(f func(item *GcpDnsDomain) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *GcpDnsDomainList) Range(f func(index int, item *GcpDnsDomain) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
