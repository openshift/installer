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

// GcpAuthentication represents the values of the 'gcp_authentication' type.
//
// Google cloud platform authentication method of a cluster.
type GcpAuthentication struct {
	fieldSet_ []bool
	href      string
	id        string
	kind      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GcpAuthentication) Empty() bool {
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

// Href returns the value of the 'href' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Self link
func (o *GcpAuthentication) Href() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.href
	}
	return ""
}

// GetHref returns the value of the 'href' attribute and
// a flag indicating if the attribute has a value.
//
// Self link
func (o *GcpAuthentication) GetHref() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.href
	}
	return
}

// Id returns the value of the 'id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Unique identifier of the object
func (o *GcpAuthentication) Id() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetId returns the value of the 'id' attribute and
// a flag indicating if the attribute has a value.
//
// Unique identifier of the object
func (o *GcpAuthentication) GetId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// Kind returns the value of the 'kind' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the type of this object
func (o *GcpAuthentication) Kind() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.kind
	}
	return ""
}

// GetKind returns the value of the 'kind' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the type of this object
func (o *GcpAuthentication) GetKind() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.kind
	}
	return
}

// GcpAuthenticationListKind is the name of the type used to represent list of objects of
// type 'gcp_authentication'.
const GcpAuthenticationListKind = "GcpAuthenticationList"

// GcpAuthenticationListLinkKind is the name of the type used to represent links to list
// of objects of type 'gcp_authentication'.
const GcpAuthenticationListLinkKind = "GcpAuthenticationListLink"

// GcpAuthenticationNilKind is the name of the type used to nil lists of objects of
// type 'gcp_authentication'.
const GcpAuthenticationListNilKind = "GcpAuthenticationListNil"

// GcpAuthenticationList is a list of values of the 'gcp_authentication' type.
type GcpAuthenticationList struct {
	href  string
	link  bool
	items []*GcpAuthentication
}

// Len returns the length of the list.
func (l *GcpAuthenticationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GcpAuthenticationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GcpAuthenticationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GcpAuthenticationList) SetItems(items []*GcpAuthentication) {
	l.items = items
}

// Items returns the items of the list.
func (l *GcpAuthenticationList) Items() []*GcpAuthentication {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GcpAuthenticationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GcpAuthenticationList) Get(i int) *GcpAuthentication {
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
func (l *GcpAuthenticationList) Slice() []*GcpAuthentication {
	var slice []*GcpAuthentication
	if l == nil {
		slice = make([]*GcpAuthentication, 0)
	} else {
		slice = make([]*GcpAuthentication, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GcpAuthenticationList) Each(f func(item *GcpAuthentication) bool) {
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
func (l *GcpAuthenticationList) Range(f func(index int, item *GcpAuthentication) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
