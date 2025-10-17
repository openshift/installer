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

// AddOnNamespaceKind is the name of the type used to represent objects
// of type 'add_on_namespace'.
const AddOnNamespaceKind = "AddOnNamespace"

// AddOnNamespaceLinkKind is the name of the type used to represent links
// to objects of type 'add_on_namespace'.
const AddOnNamespaceLinkKind = "AddOnNamespaceLink"

// AddOnNamespaceNilKind is the name of the type used to nil references
// to objects of type 'add_on_namespace'.
const AddOnNamespaceNilKind = "AddOnNamespaceNil"

// AddOnNamespace represents the values of the 'add_on_namespace' type.
type AddOnNamespace struct {
	bitmap_     uint32
	id          string
	href        string
	annotations map[string]string
	labels      map[string]string
	name        string
}

// Kind returns the name of the type of the object.
func (o *AddOnNamespace) Kind() string {
	if o == nil {
		return AddOnNamespaceNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddOnNamespaceLinkKind
	}
	return AddOnNamespaceKind
}

// Link returns true if this is a link.
func (o *AddOnNamespace) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddOnNamespace) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOnNamespace) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOnNamespace) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOnNamespace) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnNamespace) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Annotations returns the value of the 'annotations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Annotations to be applied to this namespace.
func (o *AddOnNamespace) Annotations() map[string]string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.annotations
	}
	return nil
}

// GetAnnotations returns the value of the 'annotations' attribute and
// a flag indicating if the attribute has a value.
//
// Annotations to be applied to this namespace.
func (o *AddOnNamespace) GetAnnotations() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.annotations
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Labels to be applied to this namespace.
func (o *AddOnNamespace) Labels() map[string]string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// Labels to be applied to this namespace.
func (o *AddOnNamespace) GetLabels() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.labels
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the namespace.
func (o *AddOnNamespace) Name() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the namespace.
func (o *AddOnNamespace) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.name
	}
	return
}

// AddOnNamespaceListKind is the name of the type used to represent list of objects of
// type 'add_on_namespace'.
const AddOnNamespaceListKind = "AddOnNamespaceList"

// AddOnNamespaceListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_namespace'.
const AddOnNamespaceListLinkKind = "AddOnNamespaceListLink"

// AddOnNamespaceNilKind is the name of the type used to nil lists of objects of
// type 'add_on_namespace'.
const AddOnNamespaceListNilKind = "AddOnNamespaceListNil"

// AddOnNamespaceList is a list of values of the 'add_on_namespace' type.
type AddOnNamespaceList struct {
	href  string
	link  bool
	items []*AddOnNamespace
}

// Kind returns the name of the type of the object.
func (l *AddOnNamespaceList) Kind() string {
	if l == nil {
		return AddOnNamespaceListNilKind
	}
	if l.link {
		return AddOnNamespaceListLinkKind
	}
	return AddOnNamespaceListKind
}

// Link returns true iif this is a link.
func (l *AddOnNamespaceList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnNamespaceList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnNamespaceList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnNamespaceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnNamespaceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnNamespaceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnNamespaceList) SetItems(items []*AddOnNamespace) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnNamespaceList) Items() []*AddOnNamespace {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnNamespaceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnNamespaceList) Get(i int) *AddOnNamespace {
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
func (l *AddOnNamespaceList) Slice() []*AddOnNamespace {
	var slice []*AddOnNamespace
	if l == nil {
		slice = make([]*AddOnNamespace, 0)
	} else {
		slice = make([]*AddOnNamespace, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnNamespaceList) Each(f func(item *AddOnNamespace) bool) {
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
func (l *AddOnNamespaceList) Range(f func(index int, item *AddOnNamespace) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
