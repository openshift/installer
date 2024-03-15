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

// ClientComponent represents the values of the 'client_component' type.
//
// The reference of a component that will consume the client configuration.
type ClientComponent struct {
	bitmap_   uint32
	name      string
	namespace string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClientComponent) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the component.
func (o *ClientComponent) Name() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the component.
func (o *ClientComponent) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.name
	}
	return
}

// Namespace returns the value of the 'namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The namespace of the component.
func (o *ClientComponent) Namespace() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.namespace
	}
	return ""
}

// GetNamespace returns the value of the 'namespace' attribute and
// a flag indicating if the attribute has a value.
//
// The namespace of the component.
func (o *ClientComponent) GetNamespace() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.namespace
	}
	return
}

// ClientComponentListKind is the name of the type used to represent list of objects of
// type 'client_component'.
const ClientComponentListKind = "ClientComponentList"

// ClientComponentListLinkKind is the name of the type used to represent links to list
// of objects of type 'client_component'.
const ClientComponentListLinkKind = "ClientComponentListLink"

// ClientComponentNilKind is the name of the type used to nil lists of objects of
// type 'client_component'.
const ClientComponentListNilKind = "ClientComponentListNil"

// ClientComponentList is a list of values of the 'client_component' type.
type ClientComponentList struct {
	href  string
	link  bool
	items []*ClientComponent
}

// Len returns the length of the list.
func (l *ClientComponentList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClientComponentList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClientComponentList) Get(i int) *ClientComponent {
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
func (l *ClientComponentList) Slice() []*ClientComponent {
	var slice []*ClientComponent
	if l == nil {
		slice = make([]*ClientComponent, 0)
	} else {
		slice = make([]*ClientComponent, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClientComponentList) Each(f func(item *ClientComponent) bool) {
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
func (l *ClientComponentList) Range(f func(index int, item *ClientComponent) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
