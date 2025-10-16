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

// ClusterMigrationState represents the values of the 'cluster_migration_state' type.
//
// Representation of a cluster migration state.
type ClusterMigrationState struct {
	bitmap_     uint32
	description string
	value       ClusterMigrationStateValue
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterMigrationState) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A longer description of the current state of the cluster migration.
func (o *ClusterMigrationState) Description() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// A longer description of the current state of the cluster migration.
func (o *ClusterMigrationState) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.description
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The current state of the cluster migration.
func (o *ClusterMigrationState) Value() ClusterMigrationStateValue {
	if o != nil && o.bitmap_&2 != 0 {
		return o.value
	}
	return ClusterMigrationStateValue("")
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// The current state of the cluster migration.
func (o *ClusterMigrationState) GetValue() (value ClusterMigrationStateValue, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.value
	}
	return
}

// ClusterMigrationStateListKind is the name of the type used to represent list of objects of
// type 'cluster_migration_state'.
const ClusterMigrationStateListKind = "ClusterMigrationStateList"

// ClusterMigrationStateListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_migration_state'.
const ClusterMigrationStateListLinkKind = "ClusterMigrationStateListLink"

// ClusterMigrationStateNilKind is the name of the type used to nil lists of objects of
// type 'cluster_migration_state'.
const ClusterMigrationStateListNilKind = "ClusterMigrationStateListNil"

// ClusterMigrationStateList is a list of values of the 'cluster_migration_state' type.
type ClusterMigrationStateList struct {
	href  string
	link  bool
	items []*ClusterMigrationState
}

// Len returns the length of the list.
func (l *ClusterMigrationStateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterMigrationStateList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterMigrationStateList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterMigrationStateList) SetItems(items []*ClusterMigrationState) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterMigrationStateList) Items() []*ClusterMigrationState {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterMigrationStateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterMigrationStateList) Get(i int) *ClusterMigrationState {
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
func (l *ClusterMigrationStateList) Slice() []*ClusterMigrationState {
	var slice []*ClusterMigrationState
	if l == nil {
		slice = make([]*ClusterMigrationState, 0)
	} else {
		slice = make([]*ClusterMigrationState, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterMigrationStateList) Each(f func(item *ClusterMigrationState) bool) {
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
func (l *ClusterMigrationStateList) Range(f func(index int, item *ClusterMigrationState) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
