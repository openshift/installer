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

import (
	time "time"
)

// NodePoolUpgradePolicyKind is the name of the type used to represent objects
// of type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyKind = "NodePoolUpgradePolicy"

// NodePoolUpgradePolicyLinkKind is the name of the type used to represent links
// to objects of type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyLinkKind = "NodePoolUpgradePolicyLink"

// NodePoolUpgradePolicyNilKind is the name of the type used to nil references
// to objects of type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyNilKind = "NodePoolUpgradePolicyNil"

// NodePoolUpgradePolicy represents the values of the 'node_pool_upgrade_policy' type.
//
// Representation of an upgrade policy that can be set for a node pool.
type NodePoolUpgradePolicy struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	clusterID                  string
	creationTimestamp          time.Time
	lastUpdateTimestamp        time.Time
	nextRun                    time.Time
	nodePoolID                 string
	schedule                   string
	scheduleType               ScheduleType
	state                      *UpgradePolicyState
	upgradeType                UpgradeType
	version                    string
	enableMinorVersionUpgrades bool
}

// Kind returns the name of the type of the object.
func (o *NodePoolUpgradePolicy) Kind() string {
	if o == nil {
		return NodePoolUpgradePolicyNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return NodePoolUpgradePolicyLinkKind
	}
	return NodePoolUpgradePolicyKind
}

// Link returns true if this is a link.
func (o *NodePoolUpgradePolicy) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *NodePoolUpgradePolicy) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NodePoolUpgradePolicy) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NodePoolUpgradePolicy) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NodePoolUpgradePolicy) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodePoolUpgradePolicy) Empty() bool {
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
// Cluster ID this upgrade policy for node pool is defined for.
func (o *NodePoolUpgradePolicy) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster ID this upgrade policy for node pool is defined for.
func (o *NodePoolUpgradePolicy) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.clusterID
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Timestamp for creation of resource.
func (o *NodePoolUpgradePolicy) CreationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Timestamp for creation of resource.
func (o *NodePoolUpgradePolicy) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.creationTimestamp
	}
	return
}

// EnableMinorVersionUpgrades returns the value of the 'enable_minor_version_upgrades' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if minor version upgrades are allowed for automatic upgrades (for manual it's always allowed).
func (o *NodePoolUpgradePolicy) EnableMinorVersionUpgrades() bool {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.enableMinorVersionUpgrades
	}
	return false
}

// GetEnableMinorVersionUpgrades returns the value of the 'enable_minor_version_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if minor version upgrades are allowed for automatic upgrades (for manual it's always allowed).
func (o *NodePoolUpgradePolicy) GetEnableMinorVersionUpgrades() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.enableMinorVersionUpgrades
	}
	return
}

// LastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Timestamp for last update that happened to resource.
func (o *NodePoolUpgradePolicy) LastUpdateTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Timestamp for last update that happened to resource.
func (o *NodePoolUpgradePolicy) GetLastUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.lastUpdateTimestamp
	}
	return
}

// NextRun returns the value of the 'next_run' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Next time the upgrade should run.
func (o *NodePoolUpgradePolicy) NextRun() time.Time {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.nextRun
	}
	return time.Time{}
}

// GetNextRun returns the value of the 'next_run' attribute and
// a flag indicating if the attribute has a value.
//
// Next time the upgrade should run.
func (o *NodePoolUpgradePolicy) GetNextRun() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.nextRun
	}
	return
}

// NodePoolID returns the value of the 'node_pool_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Node Pool ID this upgrade policy is defined for.
func (o *NodePoolUpgradePolicy) NodePoolID() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.nodePoolID
	}
	return ""
}

// GetNodePoolID returns the value of the 'node_pool_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Node Pool ID this upgrade policy is defined for.
func (o *NodePoolUpgradePolicy) GetNodePoolID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.nodePoolID
	}
	return
}

// Schedule returns the value of the 'schedule' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *NodePoolUpgradePolicy) Schedule() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.schedule
	}
	return ""
}

// GetSchedule returns the value of the 'schedule' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *NodePoolUpgradePolicy) GetSchedule() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.schedule
	}
	return
}

// ScheduleType returns the value of the 'schedule_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule type of the upgrade.
func (o *NodePoolUpgradePolicy) ScheduleType() ScheduleType {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.scheduleType
	}
	return ScheduleType("")
}

// GetScheduleType returns the value of the 'schedule_type' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule type of the upgrade.
func (o *NodePoolUpgradePolicy) GetScheduleType() (value ScheduleType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.scheduleType
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the upgrade policy for the node pool.
func (o *NodePoolUpgradePolicy) State() *UpgradePolicyState {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.state
	}
	return nil
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State of the upgrade policy for the node pool.
func (o *NodePoolUpgradePolicy) GetState() (value *UpgradePolicyState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.state
	}
	return
}

// UpgradeType returns the value of the 'upgrade_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Upgrade type of the node pool.
func (o *NodePoolUpgradePolicy) UpgradeType() UpgradeType {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.upgradeType
	}
	return UpgradeType("")
}

// GetUpgradeType returns the value of the 'upgrade_type' attribute and
// a flag indicating if the attribute has a value.
//
// Upgrade type of the node pool.
func (o *NodePoolUpgradePolicy) GetUpgradeType() (value UpgradeType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.upgradeType
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version is the desired upgrade version.
func (o *NodePoolUpgradePolicy) Version() string {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Version is the desired upgrade version.
func (o *NodePoolUpgradePolicy) GetVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.version
	}
	return
}

// NodePoolUpgradePolicyListKind is the name of the type used to represent list of objects of
// type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyListKind = "NodePoolUpgradePolicyList"

// NodePoolUpgradePolicyListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyListLinkKind = "NodePoolUpgradePolicyListLink"

// NodePoolUpgradePolicyNilKind is the name of the type used to nil lists of objects of
// type 'node_pool_upgrade_policy'.
const NodePoolUpgradePolicyListNilKind = "NodePoolUpgradePolicyListNil"

// NodePoolUpgradePolicyList is a list of values of the 'node_pool_upgrade_policy' type.
type NodePoolUpgradePolicyList struct {
	href  string
	link  bool
	items []*NodePoolUpgradePolicy
}

// Kind returns the name of the type of the object.
func (l *NodePoolUpgradePolicyList) Kind() string {
	if l == nil {
		return NodePoolUpgradePolicyListNilKind
	}
	if l.link {
		return NodePoolUpgradePolicyListLinkKind
	}
	return NodePoolUpgradePolicyListKind
}

// Link returns true iif this is a link.
func (l *NodePoolUpgradePolicyList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NodePoolUpgradePolicyList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NodePoolUpgradePolicyList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NodePoolUpgradePolicyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NodePoolUpgradePolicyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NodePoolUpgradePolicyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NodePoolUpgradePolicyList) SetItems(items []*NodePoolUpgradePolicy) {
	l.items = items
}

// Items returns the items of the list.
func (l *NodePoolUpgradePolicyList) Items() []*NodePoolUpgradePolicy {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NodePoolUpgradePolicyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodePoolUpgradePolicyList) Get(i int) *NodePoolUpgradePolicy {
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
func (l *NodePoolUpgradePolicyList) Slice() []*NodePoolUpgradePolicy {
	var slice []*NodePoolUpgradePolicy
	if l == nil {
		slice = make([]*NodePoolUpgradePolicy, 0)
	} else {
		slice = make([]*NodePoolUpgradePolicy, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodePoolUpgradePolicyList) Each(f func(item *NodePoolUpgradePolicy) bool) {
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
func (l *NodePoolUpgradePolicyList) Range(f func(index int, item *NodePoolUpgradePolicy) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
