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

// ManifestKind is the name of the type used to represent objects
// of type 'manifest'.
const ManifestKind = "Manifest"

// ManifestLinkKind is the name of the type used to represent links
// to objects of type 'manifest'.
const ManifestLinkKind = "ManifestLink"

// ManifestNilKind is the name of the type used to nil references
// to objects of type 'manifest'.
const ManifestNilKind = "ManifestNil"

// Manifest represents the values of the 'manifest' type.
//
// Representation of a manifestwork.
type Manifest struct {
	bitmap_           uint32
	id                string
	href              string
	creationTimestamp time.Time
	liveResource      interface{}
	spec              interface{}
	updatedTimestamp  time.Time
	workloads         []interface{}
}

// Kind returns the name of the type of the object.
func (o *Manifest) Kind() string {
	if o == nil {
		return ManifestNilKind
	}
	if o.bitmap_&1 != 0 {
		return ManifestLinkKind
	}
	return ManifestKind
}

// Link returns true if this is a link.
func (o *Manifest) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Manifest) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Manifest) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Manifest) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Manifest) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Manifest) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the manifest got created in OCM database.
func (o *Manifest) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the manifest got created in OCM database.
func (o *Manifest) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// LiveResource returns the value of the 'live_resource' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Transient value to represent the underlying live resource.
func (o *Manifest) LiveResource() interface{} {
	if o != nil && o.bitmap_&16 != 0 {
		return o.liveResource
	}
	return nil
}

// GetLiveResource returns the value of the 'live_resource' attribute and
// a flag indicating if the attribute has a value.
//
// Transient value to represent the underlying live resource.
func (o *Manifest) GetLiveResource() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.liveResource
	}
	return
}

// Spec returns the value of the 'spec' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Spec of Manifest Work object from open cluster management
// For more info please check https://open-cluster-management.io/concepts/manifestwork.
func (o *Manifest) Spec() interface{} {
	if o != nil && o.bitmap_&32 != 0 {
		return o.spec
	}
	return nil
}

// GetSpec returns the value of the 'spec' attribute and
// a flag indicating if the attribute has a value.
//
// Spec of Manifest Work object from open cluster management
// For more info please check https://open-cluster-management.io/concepts/manifestwork.
func (o *Manifest) GetSpec() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.spec
	}
	return
}

// UpdatedTimestamp returns the value of the 'updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the manifest got updated in OCM database.
func (o *Manifest) UpdatedTimestamp() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.updatedTimestamp
	}
	return time.Time{}
}

// GetUpdatedTimestamp returns the value of the 'updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the manifest got updated in OCM database.
func (o *Manifest) GetUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.updatedTimestamp
	}
	return
}

// Workloads returns the value of the 'workloads' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of k8s objects to deploy on a hosted cluster.
func (o *Manifest) Workloads() []interface{} {
	if o != nil && o.bitmap_&128 != 0 {
		return o.workloads
	}
	return nil
}

// GetWorkloads returns the value of the 'workloads' attribute and
// a flag indicating if the attribute has a value.
//
// List of k8s objects to deploy on a hosted cluster.
func (o *Manifest) GetWorkloads() (value []interface{}, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.workloads
	}
	return
}

// ManifestListKind is the name of the type used to represent list of objects of
// type 'manifest'.
const ManifestListKind = "ManifestList"

// ManifestListLinkKind is the name of the type used to represent links to list
// of objects of type 'manifest'.
const ManifestListLinkKind = "ManifestListLink"

// ManifestNilKind is the name of the type used to nil lists of objects of
// type 'manifest'.
const ManifestListNilKind = "ManifestListNil"

// ManifestList is a list of values of the 'manifest' type.
type ManifestList struct {
	href  string
	link  bool
	items []*Manifest
}

// Kind returns the name of the type of the object.
func (l *ManifestList) Kind() string {
	if l == nil {
		return ManifestListNilKind
	}
	if l.link {
		return ManifestListLinkKind
	}
	return ManifestListKind
}

// Link returns true iif this is a link.
func (l *ManifestList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ManifestList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ManifestList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ManifestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ManifestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ManifestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ManifestList) SetItems(items []*Manifest) {
	l.items = items
}

// Items returns the items of the list.
func (l *ManifestList) Items() []*Manifest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ManifestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ManifestList) Get(i int) *Manifest {
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
func (l *ManifestList) Slice() []*Manifest {
	var slice []*Manifest
	if l == nil {
		slice = make([]*Manifest, 0)
	} else {
		slice = make([]*Manifest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ManifestList) Each(f func(item *Manifest) bool) {
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
func (l *ManifestList) Range(f func(index int, item *Manifest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
