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

// ComponentRouteKind is the name of the type used to represent objects
// of type 'component_route'.
const ComponentRouteKind = "ComponentRoute"

// ComponentRouteLinkKind is the name of the type used to represent links
// to objects of type 'component_route'.
const ComponentRouteLinkKind = "ComponentRouteLink"

// ComponentRouteNilKind is the name of the type used to nil references
// to objects of type 'component_route'.
const ComponentRouteNilKind = "ComponentRouteNil"

// ComponentRoute represents the values of the 'component_route' type.
//
// Representation of a Component Route.
type ComponentRoute struct {
	fieldSet_    []bool
	id           string
	href         string
	hostname     string
	tlsSecretRef string
}

// Kind returns the name of the type of the object.
func (o *ComponentRoute) Kind() string {
	if o == nil {
		return ComponentRouteNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ComponentRouteLinkKind
	}
	return ComponentRouteKind
}

// Link returns true if this is a link.
func (o *ComponentRoute) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ComponentRoute) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ComponentRoute) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ComponentRoute) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ComponentRoute) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ComponentRoute) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// Hostname returns the value of the 'hostname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Hostname of the route.
func (o *ComponentRoute) Hostname() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.hostname
	}
	return ""
}

// GetHostname returns the value of the 'hostname' attribute and
// a flag indicating if the attribute has a value.
//
// Hostname of the route.
func (o *ComponentRoute) GetHostname() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.hostname
	}
	return
}

// TlsSecretRef returns the value of the 'tls_secret_ref' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// TLS Secret reference of the route.
func (o *ComponentRoute) TlsSecretRef() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.tlsSecretRef
	}
	return ""
}

// GetTlsSecretRef returns the value of the 'tls_secret_ref' attribute and
// a flag indicating if the attribute has a value.
//
// TLS Secret reference of the route.
func (o *ComponentRoute) GetTlsSecretRef() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.tlsSecretRef
	}
	return
}

// ComponentRouteListKind is the name of the type used to represent list of objects of
// type 'component_route'.
const ComponentRouteListKind = "ComponentRouteList"

// ComponentRouteListLinkKind is the name of the type used to represent links to list
// of objects of type 'component_route'.
const ComponentRouteListLinkKind = "ComponentRouteListLink"

// ComponentRouteNilKind is the name of the type used to nil lists of objects of
// type 'component_route'.
const ComponentRouteListNilKind = "ComponentRouteListNil"

// ComponentRouteList is a list of values of the 'component_route' type.
type ComponentRouteList struct {
	href  string
	link  bool
	items []*ComponentRoute
}

// Kind returns the name of the type of the object.
func (l *ComponentRouteList) Kind() string {
	if l == nil {
		return ComponentRouteListNilKind
	}
	if l.link {
		return ComponentRouteListLinkKind
	}
	return ComponentRouteListKind
}

// Link returns true iif this is a link.
func (l *ComponentRouteList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ComponentRouteList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ComponentRouteList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ComponentRouteList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ComponentRouteList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ComponentRouteList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ComponentRouteList) SetItems(items []*ComponentRoute) {
	l.items = items
}

// Items returns the items of the list.
func (l *ComponentRouteList) Items() []*ComponentRoute {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ComponentRouteList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ComponentRouteList) Get(i int) *ComponentRoute {
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
func (l *ComponentRouteList) Slice() []*ComponentRoute {
	var slice []*ComponentRoute
	if l == nil {
		slice = make([]*ComponentRoute, 0)
	} else {
		slice = make([]*ComponentRoute, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ComponentRouteList) Each(f func(item *ComponentRoute) bool) {
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
func (l *ComponentRouteList) Range(f func(index int, item *ComponentRoute) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
