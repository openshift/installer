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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// RegistryLocation represents the values of the 'registry_location' type.
//
// RegistryLocation contains a location of the registry specified by the registry domain
// name. The domain name might include wildcards, like '*' or '??'.
type RegistryLocation struct {
	bitmap_    uint32
	domainName string
	insecure   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RegistryLocation) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// DomainName returns the value of the 'domain_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// domainName specifies a domain name for the registry
// In case the registry use non-standard (80 or 443) port, the port should be included
// in the domain name as well.
func (o *RegistryLocation) DomainName() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.domainName
	}
	return ""
}

// GetDomainName returns the value of the 'domain_name' attribute and
// a flag indicating if the attribute has a value.
//
// domainName specifies a domain name for the registry
// In case the registry use non-standard (80 or 443) port, the port should be included
// in the domain name as well.
func (o *RegistryLocation) GetDomainName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.domainName
	}
	return
}

// Insecure returns the value of the 'insecure' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// insecure indicates whether the registry is secure (https) or insecure (http)
// By default (if not specified) the registry is assumed as secure.
func (o *RegistryLocation) Insecure() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.insecure
	}
	return false
}

// GetInsecure returns the value of the 'insecure' attribute and
// a flag indicating if the attribute has a value.
//
// insecure indicates whether the registry is secure (https) or insecure (http)
// By default (if not specified) the registry is assumed as secure.
func (o *RegistryLocation) GetInsecure() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.insecure
	}
	return
}

// RegistryLocationListKind is the name of the type used to represent list of objects of
// type 'registry_location'.
const RegistryLocationListKind = "RegistryLocationList"

// RegistryLocationListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry_location'.
const RegistryLocationListLinkKind = "RegistryLocationListLink"

// RegistryLocationNilKind is the name of the type used to nil lists of objects of
// type 'registry_location'.
const RegistryLocationListNilKind = "RegistryLocationListNil"

// RegistryLocationList is a list of values of the 'registry_location' type.
type RegistryLocationList struct {
	href  string
	link  bool
	items []*RegistryLocation
}

// Len returns the length of the list.
func (l *RegistryLocationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *RegistryLocationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *RegistryLocationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *RegistryLocationList) SetItems(items []*RegistryLocation) {
	l.items = items
}

// Items returns the items of the list.
func (l *RegistryLocationList) Items() []*RegistryLocation {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *RegistryLocationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RegistryLocationList) Get(i int) *RegistryLocation {
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
func (l *RegistryLocationList) Slice() []*RegistryLocation {
	var slice []*RegistryLocation
	if l == nil {
		slice = make([]*RegistryLocation, 0)
	} else {
		slice = make([]*RegistryLocation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RegistryLocationList) Each(f func(item *RegistryLocation) bool) {
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
func (l *RegistryLocationList) Range(f func(index int, item *RegistryLocation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
