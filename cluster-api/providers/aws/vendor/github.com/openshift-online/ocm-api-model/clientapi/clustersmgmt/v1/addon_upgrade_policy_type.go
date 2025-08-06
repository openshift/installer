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

// AddonUpgradePolicyKind is the name of the type used to represent objects
// of type 'addon_upgrade_policy'.
const AddonUpgradePolicyKind = "AddonUpgradePolicy"

// AddonUpgradePolicyLinkKind is the name of the type used to represent links
// to objects of type 'addon_upgrade_policy'.
const AddonUpgradePolicyLinkKind = "AddonUpgradePolicyLink"

// AddonUpgradePolicyNilKind is the name of the type used to nil references
// to objects of type 'addon_upgrade_policy'.
const AddonUpgradePolicyNilKind = "AddonUpgradePolicyNil"

// AddonUpgradePolicy represents the values of the 'addon_upgrade_policy' type.
//
// Representation of an upgrade policy that can be set for a cluster.
type AddonUpgradePolicy struct {
	fieldSet_    []bool
	id           string
	href         string
	addonID      string
	clusterID    string
	nextRun      time.Time
	schedule     string
	scheduleType string
	upgradeType  string
	version      string
}

// Kind returns the name of the type of the object.
func (o *AddonUpgradePolicy) Kind() string {
	if o == nil {
		return AddonUpgradePolicyNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AddonUpgradePolicyLinkKind
	}
	return AddonUpgradePolicyKind
}

// Link returns true if this is a link.
func (o *AddonUpgradePolicy) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AddonUpgradePolicy) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddonUpgradePolicy) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddonUpgradePolicy) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddonUpgradePolicy) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonUpgradePolicy) Empty() bool {
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

// AddonID returns the value of the 'addon_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Addon ID this upgrade policy is defined for
func (o *AddonUpgradePolicy) AddonID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.addonID
	}
	return ""
}

// GetAddonID returns the value of the 'addon_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Addon ID this upgrade policy is defined for
func (o *AddonUpgradePolicy) GetAddonID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.addonID
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster ID this upgrade policy is defined for.
func (o *AddonUpgradePolicy) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster ID this upgrade policy is defined for.
func (o *AddonUpgradePolicy) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.clusterID
	}
	return
}

// NextRun returns the value of the 'next_run' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Next time the upgrade should run.
func (o *AddonUpgradePolicy) NextRun() time.Time {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.nextRun
	}
	return time.Time{}
}

// GetNextRun returns the value of the 'next_run' attribute and
// a flag indicating if the attribute has a value.
//
// Next time the upgrade should run.
func (o *AddonUpgradePolicy) GetNextRun() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.nextRun
	}
	return
}

// Schedule returns the value of the 'schedule' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *AddonUpgradePolicy) Schedule() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.schedule
	}
	return ""
}

// GetSchedule returns the value of the 'schedule' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule cron expression that defines automatic upgrade scheduling.
func (o *AddonUpgradePolicy) GetSchedule() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.schedule
	}
	return
}

// ScheduleType returns the value of the 'schedule_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Schedule type can be either "manual" (single execution) or "automatic" (re-occurring).
func (o *AddonUpgradePolicy) ScheduleType() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.scheduleType
	}
	return ""
}

// GetScheduleType returns the value of the 'schedule_type' attribute and
// a flag indicating if the attribute has a value.
//
// Schedule type can be either "manual" (single execution) or "automatic" (re-occurring).
func (o *AddonUpgradePolicy) GetScheduleType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.scheduleType
	}
	return
}

// UpgradeType returns the value of the 'upgrade_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Upgrade type specify the type of the upgrade. Must be "ADDON".
func (o *AddonUpgradePolicy) UpgradeType() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.upgradeType
	}
	return ""
}

// GetUpgradeType returns the value of the 'upgrade_type' attribute and
// a flag indicating if the attribute has a value.
//
// Upgrade type specify the type of the upgrade. Must be "ADDON".
func (o *AddonUpgradePolicy) GetUpgradeType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.upgradeType
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version is the desired upgrade version.
func (o *AddonUpgradePolicy) Version() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Version is the desired upgrade version.
func (o *AddonUpgradePolicy) GetVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.version
	}
	return
}

// AddonUpgradePolicyListKind is the name of the type used to represent list of objects of
// type 'addon_upgrade_policy'.
const AddonUpgradePolicyListKind = "AddonUpgradePolicyList"

// AddonUpgradePolicyListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_upgrade_policy'.
const AddonUpgradePolicyListLinkKind = "AddonUpgradePolicyListLink"

// AddonUpgradePolicyNilKind is the name of the type used to nil lists of objects of
// type 'addon_upgrade_policy'.
const AddonUpgradePolicyListNilKind = "AddonUpgradePolicyListNil"

// AddonUpgradePolicyList is a list of values of the 'addon_upgrade_policy' type.
type AddonUpgradePolicyList struct {
	href  string
	link  bool
	items []*AddonUpgradePolicy
}

// Kind returns the name of the type of the object.
func (l *AddonUpgradePolicyList) Kind() string {
	if l == nil {
		return AddonUpgradePolicyListNilKind
	}
	if l.link {
		return AddonUpgradePolicyListLinkKind
	}
	return AddonUpgradePolicyListKind
}

// Link returns true iif this is a link.
func (l *AddonUpgradePolicyList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddonUpgradePolicyList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddonUpgradePolicyList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddonUpgradePolicyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonUpgradePolicyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonUpgradePolicyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonUpgradePolicyList) SetItems(items []*AddonUpgradePolicy) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonUpgradePolicyList) Items() []*AddonUpgradePolicy {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonUpgradePolicyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonUpgradePolicyList) Get(i int) *AddonUpgradePolicy {
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
func (l *AddonUpgradePolicyList) Slice() []*AddonUpgradePolicy {
	var slice []*AddonUpgradePolicy
	if l == nil {
		slice = make([]*AddonUpgradePolicy, 0)
	} else {
		slice = make([]*AddonUpgradePolicy, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonUpgradePolicyList) Each(f func(item *AddonUpgradePolicy) bool) {
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
func (l *AddonUpgradePolicyList) Range(f func(index int, item *AddonUpgradePolicy) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
