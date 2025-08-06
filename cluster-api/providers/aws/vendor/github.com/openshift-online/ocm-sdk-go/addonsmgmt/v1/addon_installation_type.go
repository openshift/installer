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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	time "time"
)

// AddonInstallationKind is the name of the type used to represent objects
// of type 'addon_installation'.
const AddonInstallationKind = "AddonInstallation"

// AddonInstallationLinkKind is the name of the type used to represent links
// to objects of type 'addon_installation'.
const AddonInstallationLinkKind = "AddonInstallationLink"

// AddonInstallationNilKind is the name of the type used to nil references
// to objects of type 'addon_installation'.
const AddonInstallationNilKind = "AddonInstallationNil"

// AddonInstallation represents the values of the 'addon_installation' type.
//
// Representation of addon installation
type AddonInstallation struct {
	bitmap_           uint32
	id                string
	href              string
	addon             *Addon
	addonVersion      *AddonVersion
	billing           *AddonInstallationBilling
	creationTimestamp time.Time
	csvName           string
	deletedTimestamp  time.Time
	desiredVersion    string
	operatorVersion   string
	parameters        *AddonInstallationParameterList
	state             AddonInstallationState
	stateDescription  string
	subscription      *ObjectReference
	updatedTimestamp  time.Time
}

// Kind returns the name of the type of the object.
func (o *AddonInstallation) Kind() string {
	if o == nil {
		return AddonInstallationNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddonInstallationLinkKind
	}
	return AddonInstallationKind
}

// Link returns true if this is a link.
func (o *AddonInstallation) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddonInstallation) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddonInstallation) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddonInstallation) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddonInstallation) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonInstallation) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Addon returns the value of the 'addon' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Addon installed
func (o *AddonInstallation) Addon() *Addon {
	if o != nil && o.bitmap_&8 != 0 {
		return o.addon
	}
	return nil
}

// GetAddon returns the value of the 'addon' attribute and
// a flag indicating if the attribute has a value.
//
// Addon installed
func (o *AddonInstallation) GetAddon() (value *Addon, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.addon
	}
	return
}

// AddonVersion returns the value of the 'addon_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Addon version of the addon
func (o *AddonInstallation) AddonVersion() *AddonVersion {
	if o != nil && o.bitmap_&16 != 0 {
		return o.addonVersion
	}
	return nil
}

// GetAddonVersion returns the value of the 'addon_version' attribute and
// a flag indicating if the attribute has a value.
//
// Addon version of the addon
func (o *AddonInstallation) GetAddonVersion() (value *AddonVersion, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.addonVersion
	}
	return
}

// Billing returns the value of the 'billing' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Billing of addon installation.
func (o *AddonInstallation) Billing() *AddonInstallationBilling {
	if o != nil && o.bitmap_&32 != 0 {
		return o.billing
	}
	return nil
}

// GetBilling returns the value of the 'billing' attribute and
// a flag indicating if the attribute has a value.
//
// Billing of addon installation.
func (o *AddonInstallation) GetBilling() (value *AddonInstallationBilling, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.billing
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the add-on was initially installed in the cluster.
func (o *AddonInstallation) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the add-on was initially installed in the cluster.
func (o *AddonInstallation) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// CsvName returns the value of the 'csv_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Current CSV installed on cluster
func (o *AddonInstallation) CsvName() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.csvName
	}
	return ""
}

// GetCsvName returns the value of the 'csv_name' attribute and
// a flag indicating if the attribute has a value.
//
// Current CSV installed on cluster
func (o *AddonInstallation) GetCsvName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.csvName
	}
	return
}

// DeletedTimestamp returns the value of the 'deleted_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the add-on installation deleted at.
func (o *AddonInstallation) DeletedTimestamp() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.deletedTimestamp
	}
	return time.Time{}
}

// GetDeletedTimestamp returns the value of the 'deleted_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the add-on installation deleted at.
func (o *AddonInstallation) GetDeletedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.deletedTimestamp
	}
	return
}

// DesiredVersion returns the value of the 'desired_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version of the next scheduled upgrade
func (o *AddonInstallation) DesiredVersion() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.desiredVersion
	}
	return ""
}

// GetDesiredVersion returns the value of the 'desired_version' attribute and
// a flag indicating if the attribute has a value.
//
// Version of the next scheduled upgrade
func (o *AddonInstallation) GetDesiredVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.desiredVersion
	}
	return
}

// OperatorVersion returns the value of the 'operator_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version of the operator installed by the add-on.
func (o *AddonInstallation) OperatorVersion() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.operatorVersion
	}
	return ""
}

// GetOperatorVersion returns the value of the 'operator_version' attribute and
// a flag indicating if the attribute has a value.
//
// Version of the operator installed by the add-on.
func (o *AddonInstallation) GetOperatorVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.operatorVersion
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Parameters in the installation
func (o *AddonInstallation) Parameters() *AddonInstallationParameterList {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
//
// Parameters in the installation
func (o *AddonInstallation) GetParameters() (value *AddonInstallationParameterList, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.parameters
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Addon Installation State
func (o *AddonInstallation) State() AddonInstallationState {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.state
	}
	return AddonInstallationState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// Addon Installation State
func (o *AddonInstallation) GetState() (value AddonInstallationState, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.state
	}
	return
}

// StateDescription returns the value of the 'state_description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Reason for the current State.
func (o *AddonInstallation) StateDescription() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.stateDescription
	}
	return ""
}

// GetStateDescription returns the value of the 'state_description' attribute and
// a flag indicating if the attribute has a value.
//
// Reason for the current State.
func (o *AddonInstallation) GetStateDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.stateDescription
	}
	return
}

// Subscription returns the value of the 'subscription' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Subscription for the addon installation
func (o *AddonInstallation) Subscription() *ObjectReference {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.subscription
	}
	return nil
}

// GetSubscription returns the value of the 'subscription' attribute and
// a flag indicating if the attribute has a value.
//
// Subscription for the addon installation
func (o *AddonInstallation) GetSubscription() (value *ObjectReference, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.subscription
	}
	return
}

// UpdatedTimestamp returns the value of the 'updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the add-on installation information was last updated.
func (o *AddonInstallation) UpdatedTimestamp() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.updatedTimestamp
	}
	return time.Time{}
}

// GetUpdatedTimestamp returns the value of the 'updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the add-on installation information was last updated.
func (o *AddonInstallation) GetUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.updatedTimestamp
	}
	return
}

// AddonInstallationListKind is the name of the type used to represent list of objects of
// type 'addon_installation'.
const AddonInstallationListKind = "AddonInstallationList"

// AddonInstallationListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_installation'.
const AddonInstallationListLinkKind = "AddonInstallationListLink"

// AddonInstallationNilKind is the name of the type used to nil lists of objects of
// type 'addon_installation'.
const AddonInstallationListNilKind = "AddonInstallationListNil"

// AddonInstallationList is a list of values of the 'addon_installation' type.
type AddonInstallationList struct {
	href  string
	link  bool
	items []*AddonInstallation
}

// Kind returns the name of the type of the object.
func (l *AddonInstallationList) Kind() string {
	if l == nil {
		return AddonInstallationListNilKind
	}
	if l.link {
		return AddonInstallationListLinkKind
	}
	return AddonInstallationListKind
}

// Link returns true iif this is a link.
func (l *AddonInstallationList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddonInstallationList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddonInstallationList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddonInstallationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonInstallationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonInstallationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonInstallationList) SetItems(items []*AddonInstallation) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonInstallationList) Items() []*AddonInstallation {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonInstallationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonInstallationList) Get(i int) *AddonInstallation {
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
func (l *AddonInstallationList) Slice() []*AddonInstallation {
	var slice []*AddonInstallation
	if l == nil {
		slice = make([]*AddonInstallation, 0)
	} else {
		slice = make([]*AddonInstallation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonInstallationList) Each(f func(item *AddonInstallation) bool) {
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
func (l *AddonInstallationList) Range(f func(index int, item *AddonInstallation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
