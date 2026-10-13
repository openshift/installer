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

import (
	time "time"
)

// ClusterMigrationKind is the name of the type used to represent objects
// of type 'cluster_migration'.
const ClusterMigrationKind = "ClusterMigration"

// ClusterMigrationLinkKind is the name of the type used to represent links
// to objects of type 'cluster_migration'.
const ClusterMigrationLinkKind = "ClusterMigrationLink"

// ClusterMigrationNilKind is the name of the type used to nil references
// to objects of type 'cluster_migration'.
const ClusterMigrationNilKind = "ClusterMigrationNil"

// ClusterMigration represents the values of the 'cluster_migration' type.
//
// Representation of a cluster migration.
type ClusterMigration struct {
	fieldSet_         []bool
	id                string
	href              string
	clusterID         string
	creationTimestamp time.Time
	sdnToOvn          *SdnToOvnClusterMigration
	state             *ClusterMigrationState
	type_             ClusterMigrationType
	updatedTimestamp  time.Time
}

// Kind returns the name of the type of the object.
func (o *ClusterMigration) Kind() string {
	if o == nil {
		return ClusterMigrationNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ClusterMigrationLinkKind
	}
	return ClusterMigrationKind
}

// Link returns true if this is a link.
func (o *ClusterMigration) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ClusterMigration) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ClusterMigration) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ClusterMigration) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ClusterMigration) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterMigration) Empty() bool {
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

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Internal cluster ID.
func (o *ClusterMigration) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Internal cluster ID.
func (o *ClusterMigration) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.clusterID
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the cluster migration was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ClusterMigration) CreationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the cluster migration was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ClusterMigration) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.creationTimestamp
	}
	return
}

// SdnToOvn returns the value of the 'sdn_to_ovn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details for `SdnToOvn` cluster migrations.
func (o *ClusterMigration) SdnToOvn() *SdnToOvnClusterMigration {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.sdnToOvn
	}
	return nil
}

// GetSdnToOvn returns the value of the 'sdn_to_ovn' attribute and
// a flag indicating if the attribute has a value.
//
// Details for `SdnToOvn` cluster migrations.
func (o *ClusterMigration) GetSdnToOvn() (value *SdnToOvnClusterMigration, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.sdnToOvn
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The state of the cluster migration.
func (o *ClusterMigration) State() *ClusterMigrationState {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.state
	}
	return nil
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// The state of the cluster migration.
func (o *ClusterMigration) GetState() (value *ClusterMigrationState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.state
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of cluster migration. The rest of the attributes will be populated according to this
// value. For example, if the type is `sdnToOvn` then only the `SdnToOvn` attribute will be
// populated.
func (o *ClusterMigration) Type() ClusterMigrationType {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.type_
	}
	return ClusterMigrationType("")
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of cluster migration. The rest of the attributes will be populated according to this
// value. For example, if the type is `sdnToOvn` then only the `SdnToOvn` attribute will be
// populated.
func (o *ClusterMigration) GetType() (value ClusterMigrationType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.type_
	}
	return
}

// UpdatedTimestamp returns the value of the 'updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the cluster migration was last updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ClusterMigration) UpdatedTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.updatedTimestamp
	}
	return time.Time{}
}

// GetUpdatedTimestamp returns the value of the 'updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the cluster migration was last updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *ClusterMigration) GetUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.updatedTimestamp
	}
	return
}

// ClusterMigrationListKind is the name of the type used to represent list of objects of
// type 'cluster_migration'.
const ClusterMigrationListKind = "ClusterMigrationList"

// ClusterMigrationListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_migration'.
const ClusterMigrationListLinkKind = "ClusterMigrationListLink"

// ClusterMigrationNilKind is the name of the type used to nil lists of objects of
// type 'cluster_migration'.
const ClusterMigrationListNilKind = "ClusterMigrationListNil"

// ClusterMigrationList is a list of values of the 'cluster_migration' type.
type ClusterMigrationList struct {
	href  string
	link  bool
	items []*ClusterMigration
}

// Kind returns the name of the type of the object.
func (l *ClusterMigrationList) Kind() string {
	if l == nil {
		return ClusterMigrationListNilKind
	}
	if l.link {
		return ClusterMigrationListLinkKind
	}
	return ClusterMigrationListKind
}

// Link returns true iif this is a link.
func (l *ClusterMigrationList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ClusterMigrationList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ClusterMigrationList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ClusterMigrationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterMigrationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterMigrationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterMigrationList) SetItems(items []*ClusterMigration) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterMigrationList) Items() []*ClusterMigration {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterMigrationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterMigrationList) Get(i int) *ClusterMigration {
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
func (l *ClusterMigrationList) Slice() []*ClusterMigration {
	var slice []*ClusterMigration
	if l == nil {
		slice = make([]*ClusterMigration, 0)
	} else {
		slice = make([]*ClusterMigration, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterMigrationList) Each(f func(item *ClusterMigration) bool) {
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
func (l *ClusterMigrationList) Range(f func(index int, item *ClusterMigration) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
