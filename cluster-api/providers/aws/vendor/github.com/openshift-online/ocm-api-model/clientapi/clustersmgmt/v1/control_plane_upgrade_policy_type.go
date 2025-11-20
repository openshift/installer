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

// ControlPlaneUpgradePolicyKind is the name of the type used to represent objects
// of type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyKind = "ControlPlaneUpgradePolicy"

// ControlPlaneUpgradePolicyLinkKind is the name of the type used to represent links
// to objects of type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyLinkKind = "ControlPlaneUpgradePolicyLink"

// ControlPlaneUpgradePolicyNilKind is the name of the type used to nil references
// to objects of type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyNilKind = "ControlPlaneUpgradePolicyNil"

// ControlPlaneUpgradePolicy represents the values of the 'control_plane_upgrade_policy' type.
//
// Representation of an upgrade policy that can be set for a cluster.
type ControlPlaneUpgradePolicy struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	clusterID                  string
	creationTimestamp          time.Time
	lastUpdateTimestamp        time.Time
	nextRun                    time.Time
	schedule                   string
	scheduleType               ScheduleType
	state                      *UpgradePolicyState
	upgradeType                UpgradeType
	version                    string
	enableMinorVersionUpgrades bool
}

// Kind returns the name of the type of the object.
func (o *ControlPlaneUpgradePolicy) Kind() string {
	if o == nil {
		return ControlPlaneUpgradePolicyNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ControlPlaneUpgradePolicyLinkKind
	}
	return ControlPlaneUpgradePolicyKind
}

// Link returns true if this is a link.
func (o *ControlPlaneUpgradePolicy) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ControlPlaneUpgradePolicy) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ControlPlaneUpgradePolicy) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ControlPlaneUpgradePolicy) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ControlPlaneUpgradePolicy) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ControlPlaneUpgradePolicy) Empty() bool {
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
// Cluster ID this upgrade policy for control plane is defined for.
func (o *ControlPlaneUpgradePolicy) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster ID this upgrade policy for control plane is defined for.
func (o *ControlPlaneUpgradePolicy) GetClusterID() (value string, ok bool) {
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
func (o *ControlPlaneUpgradePolicy) CreationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Timestamp for creation of resource.
func (o *ControlPlaneUpgradePolicy) GetCreationTimestamp() (value time.Time, ok bool) {
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
func (o *ControlPlaneUpgradePolicy) EnableMinorVersionUpgrades() bool {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.enableMinorVersionUpgrades
	}
	return false
}

// GetEnableMinorVersionUpgrades returns the value of the 'enable_minor_version_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if minor version upgrades are allowed for automatic upgrades (for manual it's always allowed).
func (o *ControlPlaneUpgradePolicy) GetEnableMinorVersionUpgrades() (value bool, ok bool) {
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
func (o *ControlPlaneUpgradePolicy) LastUpdateTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Timestamp for last update that happened to resource.
func (o *ControlPlaneUpgradePolicy) GetLastUpdateTimestamp() (value time.Time, ok bool) {
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
func (o *ControlPlaneUpgradePolicy) NextRun() time.Time {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.nextRun
	}
	return time.Time{}
}

// GetNextRun returns the value of the 'next_run' attribute and
// a flag indicating if the attribute has a value.
//
// Next time the upgrade should run.
func (o *ControlPlaneUpgradePolicy) GetNextRun() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.nextRun
	}
	return
}

// Schedule returns the value of the 'schedule' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *ControlPlaneUpgradePolicy) Schedule() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.schedule
	}
	return ""
}

// GetSchedule returns the value of the 'schedule' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *ControlPlaneUpgradePolicy) GetSchedule() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.schedule
	}
	return
}

// ScheduleType returns the value of the 'schedule_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule type of the control plane upgrade.
func (o *ControlPlaneUpgradePolicy) ScheduleType() ScheduleType {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.scheduleType
	}
	return ScheduleType("")
}

// GetScheduleType returns the value of the 'schedule_type' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule type of the control plane upgrade.
func (o *ControlPlaneUpgradePolicy) GetScheduleType() (value ScheduleType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.scheduleType
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the upgrade policy for the hosted control plane.
func (o *ControlPlaneUpgradePolicy) State() *UpgradePolicyState {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.state
	}
	return nil
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State of the upgrade policy for the hosted control plane.
func (o *ControlPlaneUpgradePolicy) GetState() (value *UpgradePolicyState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.state
	}
	return
}

// UpgradeType returns the value of the 'upgrade_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Upgrade type of the control plane.
func (o *ControlPlaneUpgradePolicy) UpgradeType() UpgradeType {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.upgradeType
	}
	return UpgradeType("")
}

// GetUpgradeType returns the value of the 'upgrade_type' attribute and
// a flag indicating if the attribute has a value.
//
// Upgrade type of the control plane.
func (o *ControlPlaneUpgradePolicy) GetUpgradeType() (value UpgradeType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.upgradeType
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version is the desired upgrade version.
func (o *ControlPlaneUpgradePolicy) Version() string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Version is the desired upgrade version.
func (o *ControlPlaneUpgradePolicy) GetVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.version
	}
	return
}

// ControlPlaneUpgradePolicyListKind is the name of the type used to represent list of objects of
// type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyListKind = "ControlPlaneUpgradePolicyList"

// ControlPlaneUpgradePolicyListLinkKind is the name of the type used to represent links to list
// of objects of type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyListLinkKind = "ControlPlaneUpgradePolicyListLink"

// ControlPlaneUpgradePolicyNilKind is the name of the type used to nil lists of objects of
// type 'control_plane_upgrade_policy'.
const ControlPlaneUpgradePolicyListNilKind = "ControlPlaneUpgradePolicyListNil"

// ControlPlaneUpgradePolicyList is a list of values of the 'control_plane_upgrade_policy' type.
type ControlPlaneUpgradePolicyList struct {
	href  string
	link  bool
	items []*ControlPlaneUpgradePolicy
}

// Kind returns the name of the type of the object.
func (l *ControlPlaneUpgradePolicyList) Kind() string {
	if l == nil {
		return ControlPlaneUpgradePolicyListNilKind
	}
	if l.link {
		return ControlPlaneUpgradePolicyListLinkKind
	}
	return ControlPlaneUpgradePolicyListKind
}

// Link returns true iif this is a link.
func (l *ControlPlaneUpgradePolicyList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ControlPlaneUpgradePolicyList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ControlPlaneUpgradePolicyList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ControlPlaneUpgradePolicyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ControlPlaneUpgradePolicyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ControlPlaneUpgradePolicyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ControlPlaneUpgradePolicyList) SetItems(items []*ControlPlaneUpgradePolicy) {
	l.items = items
}

// Items returns the items of the list.
func (l *ControlPlaneUpgradePolicyList) Items() []*ControlPlaneUpgradePolicy {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ControlPlaneUpgradePolicyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ControlPlaneUpgradePolicyList) Get(i int) *ControlPlaneUpgradePolicy {
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
func (l *ControlPlaneUpgradePolicyList) Slice() []*ControlPlaneUpgradePolicy {
	var slice []*ControlPlaneUpgradePolicy
	if l == nil {
		slice = make([]*ControlPlaneUpgradePolicy, 0)
	} else {
		slice = make([]*ControlPlaneUpgradePolicy, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ControlPlaneUpgradePolicyList) Each(f func(item *ControlPlaneUpgradePolicy) bool) {
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
func (l *ControlPlaneUpgradePolicyList) Range(f func(index int, item *ControlPlaneUpgradePolicy) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
