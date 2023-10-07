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

// SyncsetKind is the name of the type used to represent objects
// of type 'syncset'.
const SyncsetKind = "Syncset"

// SyncsetLinkKind is the name of the type used to represent links
// to objects of type 'syncset'.
const SyncsetLinkKind = "SyncsetLink"

// SyncsetNilKind is the name of the type used to nil references
// to objects of type 'syncset'.
const SyncsetNilKind = "SyncsetNil"

// Syncset represents the values of the 'syncset' type.
//
// Representation of a syncset.
type Syncset struct {
	bitmap_   uint32
	id        string
	href      string
	resources []interface{}
}

// Kind returns the name of the type of the object.
func (o *Syncset) Kind() string {
	if o == nil {
		return SyncsetNilKind
	}
	if o.bitmap_&1 != 0 {
		return SyncsetLinkKind
	}
	return SyncsetKind
}

// Link returns true iif this is a link.
func (o *Syncset) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Syncset) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Syncset) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Syncset) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Syncset) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Syncset) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Resources returns the value of the 'resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of k8s objects to configure for the cluster.
func (o *Syncset) Resources() []interface{} {
	if o != nil && o.bitmap_&8 != 0 {
		return o.resources
	}
	return nil
}

// GetResources returns the value of the 'resources' attribute and
// a flag indicating if the attribute has a value.
//
// List of k8s objects to configure for the cluster.
func (o *Syncset) GetResources() (value []interface{}, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.resources
	}
	return
}

// SyncsetListKind is the name of the type used to represent list of objects of
// type 'syncset'.
const SyncsetListKind = "SyncsetList"

// SyncsetListLinkKind is the name of the type used to represent links to list
// of objects of type 'syncset'.
const SyncsetListLinkKind = "SyncsetListLink"

// SyncsetNilKind is the name of the type used to nil lists of objects of
// type 'syncset'.
const SyncsetListNilKind = "SyncsetListNil"

// SyncsetList is a list of values of the 'syncset' type.
type SyncsetList struct {
	href  string
	link  bool
	items []*Syncset
}

// Kind returns the name of the type of the object.
func (l *SyncsetList) Kind() string {
	if l == nil {
		return SyncsetListNilKind
	}
	if l.link {
		return SyncsetListLinkKind
	}
	return SyncsetListKind
}

// Link returns true iif this is a link.
func (l *SyncsetList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *SyncsetList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *SyncsetList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *SyncsetList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *SyncsetList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SyncsetList) Get(i int) *Syncset {
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
func (l *SyncsetList) Slice() []*Syncset {
	var slice []*Syncset
	if l == nil {
		slice = make([]*Syncset, 0)
	} else {
		slice = make([]*Syncset, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SyncsetList) Each(f func(item *Syncset) bool) {
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
func (l *SyncsetList) Range(f func(index int, item *Syncset) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
