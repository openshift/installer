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

// DNSDomainKind is the name of the type used to represent objects
// of type 'DNS_domain'.
const DNSDomainKind = "DNSDomain"

// DNSDomainLinkKind is the name of the type used to represent links
// to objects of type 'DNS_domain'.
const DNSDomainLinkKind = "DNSDomainLink"

// DNSDomainNilKind is the name of the type used to nil references
// to objects of type 'DNS_domain'.
const DNSDomainNilKind = "DNSDomainNil"

// DNSDomain represents the values of the 'DNS_domain' type.
//
// Contains the properties of a DNS domain.
type DNSDomain struct {
	bitmap_             uint32
	id                  string
	href                string
	cluster             *ClusterLink
	organization        *OrganizationLink
	reservedAtTimestamp time.Time
	userDefined         bool
}

// Kind returns the name of the type of the object.
func (o *DNSDomain) Kind() string {
	if o == nil {
		return DNSDomainNilKind
	}
	if o.bitmap_&1 != 0 {
		return DNSDomainLinkKind
	}
	return DNSDomainKind
}

// Link returns true iif this is a link.
func (o *DNSDomain) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *DNSDomain) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *DNSDomain) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *DNSDomain) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *DNSDomain) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *DNSDomain) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Cluster returns the value of the 'cluster' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cluster that is registered with the DNS domain (optional).
func (o *DNSDomain) Cluster() *ClusterLink {
	if o != nil && o.bitmap_&8 != 0 {
		return o.cluster
	}
	return nil
}

// GetCluster returns the value of the 'cluster' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cluster that is registered with the DNS domain (optional).
func (o *DNSDomain) GetCluster() (value *ClusterLink, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.cluster
	}
	return
}

// Organization returns the value of the 'organization' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the organization that reserved the DNS domain.
func (o *DNSDomain) Organization() *OrganizationLink {
	if o != nil && o.bitmap_&16 != 0 {
		return o.organization
	}
	return nil
}

// GetOrganization returns the value of the 'organization' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the organization that reserved the DNS domain.
func (o *DNSDomain) GetOrganization() (value *OrganizationLink, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.organization
	}
	return
}

// ReservedAtTimestamp returns the value of the 'reserved_at_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the DNS domain was reserved.
func (o *DNSDomain) ReservedAtTimestamp() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.reservedAtTimestamp
	}
	return time.Time{}
}

// GetReservedAtTimestamp returns the value of the 'reserved_at_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the DNS domain was reserved.
func (o *DNSDomain) GetReservedAtTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.reservedAtTimestamp
	}
	return
}

// UserDefined returns the value of the 'user_defined' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this dns domain is user defined.
func (o *DNSDomain) UserDefined() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.userDefined
	}
	return false
}

// GetUserDefined returns the value of the 'user_defined' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this dns domain is user defined.
func (o *DNSDomain) GetUserDefined() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.userDefined
	}
	return
}

// DNSDomainListKind is the name of the type used to represent list of objects of
// type 'DNS_domain'.
const DNSDomainListKind = "DNSDomainList"

// DNSDomainListLinkKind is the name of the type used to represent links to list
// of objects of type 'DNS_domain'.
const DNSDomainListLinkKind = "DNSDomainListLink"

// DNSDomainNilKind is the name of the type used to nil lists of objects of
// type 'DNS_domain'.
const DNSDomainListNilKind = "DNSDomainListNil"

// DNSDomainList is a list of values of the 'DNS_domain' type.
type DNSDomainList struct {
	href  string
	link  bool
	items []*DNSDomain
}

// Kind returns the name of the type of the object.
func (l *DNSDomainList) Kind() string {
	if l == nil {
		return DNSDomainListNilKind
	}
	if l.link {
		return DNSDomainListLinkKind
	}
	return DNSDomainListKind
}

// Link returns true iif this is a link.
func (l *DNSDomainList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *DNSDomainList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *DNSDomainList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *DNSDomainList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *DNSDomainList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DNSDomainList) Get(i int) *DNSDomain {
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
func (l *DNSDomainList) Slice() []*DNSDomain {
	var slice []*DNSDomain
	if l == nil {
		slice = make([]*DNSDomain, 0)
	} else {
		slice = make([]*DNSDomain, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DNSDomainList) Each(f func(item *DNSDomain) bool) {
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
func (l *DNSDomainList) Range(f func(index int, item *DNSDomain) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
