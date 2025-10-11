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

// ClusterResourcesKind is the name of the type used to represent objects
// of type 'cluster_resources'.
const ClusterResourcesKind = "ClusterResources"

// ClusterResourcesLinkKind is the name of the type used to represent links
// to objects of type 'cluster_resources'.
const ClusterResourcesLinkKind = "ClusterResourcesLink"

// ClusterResourcesNilKind is the name of the type used to nil references
// to objects of type 'cluster_resources'.
const ClusterResourcesNilKind = "ClusterResourcesNil"

// ClusterResources represents the values of the 'cluster_resources' type.
//
// Cluster Resource which belongs to a cluster, example Cluster Deployment.
type ClusterResources struct {
	bitmap_           uint32
	id                string
	href              string
	clusterID         string
	creationTimestamp time.Time
	resources         map[string]string
}

// Kind returns the name of the type of the object.
func (o *ClusterResources) Kind() string {
	if o == nil {
		return ClusterResourcesNilKind
	}
	if o.bitmap_&1 != 0 {
		return ClusterResourcesLinkKind
	}
	return ClusterResourcesKind
}

// Link returns true if this is a link.
func (o *ClusterResources) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ClusterResources) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ClusterResources) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ClusterResources) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ClusterResources) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterResources) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster ID for the fetched resources
func (o *ClusterResources) ClusterID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster ID for the fetched resources
func (o *ClusterResources) GetClusterID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clusterID
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the resources were fetched.
func (o *ClusterResources) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the resources were fetched.
func (o *ClusterResources) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// Resources returns the value of the 'resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Returned map of cluster resources fetched.
func (o *ClusterResources) Resources() map[string]string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.resources
	}
	return nil
}

// GetResources returns the value of the 'resources' attribute and
// a flag indicating if the attribute has a value.
//
// Returned map of cluster resources fetched.
func (o *ClusterResources) GetResources() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.resources
	}
	return
}

// ClusterResourcesListKind is the name of the type used to represent list of objects of
// type 'cluster_resources'.
const ClusterResourcesListKind = "ClusterResourcesList"

// ClusterResourcesListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_resources'.
const ClusterResourcesListLinkKind = "ClusterResourcesListLink"

// ClusterResourcesNilKind is the name of the type used to nil lists of objects of
// type 'cluster_resources'.
const ClusterResourcesListNilKind = "ClusterResourcesListNil"

// ClusterResourcesList is a list of values of the 'cluster_resources' type.
type ClusterResourcesList struct {
	href  string
	link  bool
	items []*ClusterResources
}

// Kind returns the name of the type of the object.
func (l *ClusterResourcesList) Kind() string {
	if l == nil {
		return ClusterResourcesListNilKind
	}
	if l.link {
		return ClusterResourcesListLinkKind
	}
	return ClusterResourcesListKind
}

// Link returns true iif this is a link.
func (l *ClusterResourcesList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ClusterResourcesList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ClusterResourcesList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ClusterResourcesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterResourcesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterResourcesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterResourcesList) SetItems(items []*ClusterResources) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterResourcesList) Items() []*ClusterResources {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterResourcesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterResourcesList) Get(i int) *ClusterResources {
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
func (l *ClusterResourcesList) Slice() []*ClusterResources {
	var slice []*ClusterResources
	if l == nil {
		slice = make([]*ClusterResources, 0)
	} else {
		slice = make([]*ClusterResources, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterResourcesList) Each(f func(item *ClusterResources) bool) {
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
func (l *ClusterResourcesList) Range(f func(index int, item *ClusterResources) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
