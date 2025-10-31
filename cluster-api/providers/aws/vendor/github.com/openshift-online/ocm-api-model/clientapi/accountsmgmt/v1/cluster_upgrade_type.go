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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	time "time"
)

// ClusterUpgrade represents the values of the 'cluster_upgrade' type.
type ClusterUpgrade struct {
	fieldSet_        []bool
	state            string
	updatedTimestamp time.Time
	version          string
	available        bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterUpgrade) Empty() bool {
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

// Available returns the value of the 'available' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterUpgrade) Available() bool {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.available
	}
	return false
}

// GetAvailable returns the value of the 'available' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterUpgrade) GetAvailable() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.available
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterUpgrade) State() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.state
	}
	return ""
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterUpgrade) GetState() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.state
	}
	return
}

// UpdatedTimestamp returns the value of the 'updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterUpgrade) UpdatedTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.updatedTimestamp
	}
	return time.Time{}
}

// GetUpdatedTimestamp returns the value of the 'updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterUpgrade) GetUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.updatedTimestamp
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterUpgrade) Version() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterUpgrade) GetVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.version
	}
	return
}

// ClusterUpgradeListKind is the name of the type used to represent list of objects of
// type 'cluster_upgrade'.
const ClusterUpgradeListKind = "ClusterUpgradeList"

// ClusterUpgradeListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_upgrade'.
const ClusterUpgradeListLinkKind = "ClusterUpgradeListLink"

// ClusterUpgradeNilKind is the name of the type used to nil lists of objects of
// type 'cluster_upgrade'.
const ClusterUpgradeListNilKind = "ClusterUpgradeListNil"

// ClusterUpgradeList is a list of values of the 'cluster_upgrade' type.
type ClusterUpgradeList struct {
	href  string
	link  bool
	items []*ClusterUpgrade
}

// Len returns the length of the list.
func (l *ClusterUpgradeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterUpgradeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterUpgradeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterUpgradeList) SetItems(items []*ClusterUpgrade) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterUpgradeList) Items() []*ClusterUpgrade {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterUpgradeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterUpgradeList) Get(i int) *ClusterUpgrade {
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
func (l *ClusterUpgradeList) Slice() []*ClusterUpgrade {
	var slice []*ClusterUpgrade
	if l == nil {
		slice = make([]*ClusterUpgrade, 0)
	} else {
		slice = make([]*ClusterUpgrade, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterUpgradeList) Each(f func(item *ClusterUpgrade) bool) {
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
func (l *ClusterUpgradeList) Range(f func(index int, item *ClusterUpgrade) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
