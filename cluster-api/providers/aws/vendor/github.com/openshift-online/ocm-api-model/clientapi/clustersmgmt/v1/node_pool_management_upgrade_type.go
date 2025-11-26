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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// NodePoolManagementUpgradeKind is the name of the type used to represent objects
// of type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeKind = "NodePoolManagementUpgrade"

// NodePoolManagementUpgradeLinkKind is the name of the type used to represent links
// to objects of type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeLinkKind = "NodePoolManagementUpgradeLink"

// NodePoolManagementUpgradeNilKind is the name of the type used to nil references
// to objects of type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeNilKind = "NodePoolManagementUpgradeNil"

// NodePoolManagementUpgrade represents the values of the 'node_pool_management_upgrade' type.
//
// Representation of node pool management.
type NodePoolManagementUpgrade struct {
	fieldSet_      []bool
	id             string
	href           string
	maxSurge       string
	maxUnavailable string
	type_          string
}

// Kind returns the name of the type of the object.
func (o *NodePoolManagementUpgrade) Kind() string {
	if o == nil {
		return NodePoolManagementUpgradeNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return NodePoolManagementUpgradeLinkKind
	}
	return NodePoolManagementUpgradeKind
}

// Link returns true if this is a link.
func (o *NodePoolManagementUpgrade) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *NodePoolManagementUpgrade) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NodePoolManagementUpgrade) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NodePoolManagementUpgrade) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NodePoolManagementUpgrade) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodePoolManagementUpgrade) Empty() bool {
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

// MaxSurge returns the value of the 'max_surge' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Maximum number of nodes in the NodePool of a ROSA HCP cluster that can be scheduled above the desired number of nodes during the upgrade.
func (o *NodePoolManagementUpgrade) MaxSurge() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.maxSurge
	}
	return ""
}

// GetMaxSurge returns the value of the 'max_surge' attribute and
// a flag indicating if the attribute has a value.
//
// Maximum number of nodes in the NodePool of a ROSA HCP cluster that can be scheduled above the desired number of nodes during the upgrade.
func (o *NodePoolManagementUpgrade) GetMaxSurge() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.maxSurge
	}
	return
}

// MaxUnavailable returns the value of the 'max_unavailable' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Maximum number of nodes in the NodePool of a ROSA HCP cluster that can be unavailable during the upgrade.
func (o *NodePoolManagementUpgrade) MaxUnavailable() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.maxUnavailable
	}
	return ""
}

// GetMaxUnavailable returns the value of the 'max_unavailable' attribute and
// a flag indicating if the attribute has a value.
//
// Maximum number of nodes in the NodePool of a ROSA HCP cluster that can be unavailable during the upgrade.
func (o *NodePoolManagementUpgrade) GetMaxUnavailable() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.maxUnavailable
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of strategy for handling upgrades.
func (o *NodePoolManagementUpgrade) Type() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of strategy for handling upgrades.
func (o *NodePoolManagementUpgrade) GetType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.type_
	}
	return
}

// NodePoolManagementUpgradeListKind is the name of the type used to represent list of objects of
// type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeListKind = "NodePoolManagementUpgradeList"

// NodePoolManagementUpgradeListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeListLinkKind = "NodePoolManagementUpgradeListLink"

// NodePoolManagementUpgradeNilKind is the name of the type used to nil lists of objects of
// type 'node_pool_management_upgrade'.
const NodePoolManagementUpgradeListNilKind = "NodePoolManagementUpgradeListNil"

// NodePoolManagementUpgradeList is a list of values of the 'node_pool_management_upgrade' type.
type NodePoolManagementUpgradeList struct {
	href  string
	link  bool
	items []*NodePoolManagementUpgrade
}

// Kind returns the name of the type of the object.
func (l *NodePoolManagementUpgradeList) Kind() string {
	if l == nil {
		return NodePoolManagementUpgradeListNilKind
	}
	if l.link {
		return NodePoolManagementUpgradeListLinkKind
	}
	return NodePoolManagementUpgradeListKind
}

// Link returns true iif this is a link.
func (l *NodePoolManagementUpgradeList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NodePoolManagementUpgradeList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NodePoolManagementUpgradeList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NodePoolManagementUpgradeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NodePoolManagementUpgradeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NodePoolManagementUpgradeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NodePoolManagementUpgradeList) SetItems(items []*NodePoolManagementUpgrade) {
	l.items = items
}

// Items returns the items of the list.
func (l *NodePoolManagementUpgradeList) Items() []*NodePoolManagementUpgrade {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NodePoolManagementUpgradeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodePoolManagementUpgradeList) Get(i int) *NodePoolManagementUpgrade {
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
func (l *NodePoolManagementUpgradeList) Slice() []*NodePoolManagementUpgrade {
	var slice []*NodePoolManagementUpgrade
	if l == nil {
		slice = make([]*NodePoolManagementUpgrade, 0)
	} else {
		slice = make([]*NodePoolManagementUpgrade, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodePoolManagementUpgradeList) Each(f func(item *NodePoolManagementUpgrade) bool) {
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
func (l *NodePoolManagementUpgradeList) Range(f func(index int, item *NodePoolManagementUpgrade) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
