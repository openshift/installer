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

// PendingDeleteClusterKind is the name of the type used to represent objects
// of type 'pending_delete_cluster'.
const PendingDeleteClusterKind = "PendingDeleteCluster"

// PendingDeleteClusterLinkKind is the name of the type used to represent links
// to objects of type 'pending_delete_cluster'.
const PendingDeleteClusterLinkKind = "PendingDeleteClusterLink"

// PendingDeleteClusterNilKind is the name of the type used to nil references
// to objects of type 'pending_delete_cluster'.
const PendingDeleteClusterNilKind = "PendingDeleteClusterNil"

// PendingDeleteCluster represents the values of the 'pending_delete_cluster' type.
//
// Represents a pending delete entry for a specific cluster.
type PendingDeleteCluster struct {
	bitmap_           uint32
	id                string
	href              string
	cluster           *Cluster
	creationTimestamp time.Time
	bestEffort        bool
}

// Kind returns the name of the type of the object.
func (o *PendingDeleteCluster) Kind() string {
	if o == nil {
		return PendingDeleteClusterNilKind
	}
	if o.bitmap_&1 != 0 {
		return PendingDeleteClusterLinkKind
	}
	return PendingDeleteClusterKind
}

// Link returns true if this is a link.
func (o *PendingDeleteCluster) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *PendingDeleteCluster) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *PendingDeleteCluster) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *PendingDeleteCluster) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *PendingDeleteCluster) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *PendingDeleteCluster) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// BestEffort returns the value of the 'best_effort' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Flag indicating if the cluster deletion should be best-effort mode or not.
func (o *PendingDeleteCluster) BestEffort() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.bestEffort
	}
	return false
}

// GetBestEffort returns the value of the 'best_effort' attribute and
// a flag indicating if the attribute has a value.
//
// Flag indicating if the cluster deletion should be best-effort mode or not.
func (o *PendingDeleteCluster) GetBestEffort() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.bestEffort
	}
	return
}

// Cluster returns the value of the 'cluster' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster is the details of the cluster that is pending deletion.
func (o *PendingDeleteCluster) Cluster() *Cluster {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cluster
	}
	return nil
}

// GetCluster returns the value of the 'cluster' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster is the details of the cluster that is pending deletion.
func (o *PendingDeleteCluster) GetCluster() (value *Cluster, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cluster
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the cluster was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *PendingDeleteCluster) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the cluster was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *PendingDeleteCluster) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// PendingDeleteClusterListKind is the name of the type used to represent list of objects of
// type 'pending_delete_cluster'.
const PendingDeleteClusterListKind = "PendingDeleteClusterList"

// PendingDeleteClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'pending_delete_cluster'.
const PendingDeleteClusterListLinkKind = "PendingDeleteClusterListLink"

// PendingDeleteClusterNilKind is the name of the type used to nil lists of objects of
// type 'pending_delete_cluster'.
const PendingDeleteClusterListNilKind = "PendingDeleteClusterListNil"

// PendingDeleteClusterList is a list of values of the 'pending_delete_cluster' type.
type PendingDeleteClusterList struct {
	href  string
	link  bool
	items []*PendingDeleteCluster
}

// Kind returns the name of the type of the object.
func (l *PendingDeleteClusterList) Kind() string {
	if l == nil {
		return PendingDeleteClusterListNilKind
	}
	if l.link {
		return PendingDeleteClusterListLinkKind
	}
	return PendingDeleteClusterListKind
}

// Link returns true iif this is a link.
func (l *PendingDeleteClusterList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *PendingDeleteClusterList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *PendingDeleteClusterList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *PendingDeleteClusterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *PendingDeleteClusterList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *PendingDeleteClusterList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *PendingDeleteClusterList) SetItems(items []*PendingDeleteCluster) {
	l.items = items
}

// Items returns the items of the list.
func (l *PendingDeleteClusterList) Items() []*PendingDeleteCluster {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *PendingDeleteClusterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PendingDeleteClusterList) Get(i int) *PendingDeleteCluster {
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
func (l *PendingDeleteClusterList) Slice() []*PendingDeleteCluster {
	var slice []*PendingDeleteCluster
	if l == nil {
		slice = make([]*PendingDeleteCluster, 0)
	} else {
		slice = make([]*PendingDeleteCluster, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PendingDeleteClusterList) Each(f func(item *PendingDeleteCluster) bool) {
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
func (l *PendingDeleteClusterList) Range(f func(index int, item *PendingDeleteCluster) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
