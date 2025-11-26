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

// Environment represents the values of the 'environment' type.
//
// Description of an environment
type Environment struct {
	fieldSet_                 []bool
	backplaneURL              string
	lastLimitedSupportCheck   time.Time
	lastUpgradeAvailableCheck time.Time
	name                      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Environment) Empty() bool {
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

// BackplaneURL returns the value of the 'backplane_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// the backplane url for the environment
func (o *Environment) BackplaneURL() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.backplaneURL
	}
	return ""
}

// GetBackplaneURL returns the value of the 'backplane_URL' attribute and
// a flag indicating if the attribute has a value.
//
// the backplane url for the environment
func (o *Environment) GetBackplaneURL() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.backplaneURL
	}
	return
}

// LastLimitedSupportCheck returns the value of the 'last_limited_support_check' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// last time that the worker checked for limited support clusters
func (o *Environment) LastLimitedSupportCheck() time.Time {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.lastLimitedSupportCheck
	}
	return time.Time{}
}

// GetLastLimitedSupportCheck returns the value of the 'last_limited_support_check' attribute and
// a flag indicating if the attribute has a value.
//
// last time that the worker checked for limited support clusters
func (o *Environment) GetLastLimitedSupportCheck() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.lastLimitedSupportCheck
	}
	return
}

// LastUpgradeAvailableCheck returns the value of the 'last_upgrade_available_check' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// last time that the worker checked for available upgrades
func (o *Environment) LastUpgradeAvailableCheck() time.Time {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.lastUpgradeAvailableCheck
	}
	return time.Time{}
}

// GetLastUpgradeAvailableCheck returns the value of the 'last_upgrade_available_check' attribute and
// a flag indicating if the attribute has a value.
//
// last time that the worker checked for available upgrades
func (o *Environment) GetLastUpgradeAvailableCheck() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.lastUpgradeAvailableCheck
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// environment name
func (o *Environment) Name() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// environment name
func (o *Environment) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.name
	}
	return
}

// EnvironmentListKind is the name of the type used to represent list of objects of
// type 'environment'.
const EnvironmentListKind = "EnvironmentList"

// EnvironmentListLinkKind is the name of the type used to represent links to list
// of objects of type 'environment'.
const EnvironmentListLinkKind = "EnvironmentListLink"

// EnvironmentNilKind is the name of the type used to nil lists of objects of
// type 'environment'.
const EnvironmentListNilKind = "EnvironmentListNil"

// EnvironmentList is a list of values of the 'environment' type.
type EnvironmentList struct {
	href  string
	link  bool
	items []*Environment
}

// Len returns the length of the list.
func (l *EnvironmentList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *EnvironmentList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *EnvironmentList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *EnvironmentList) SetItems(items []*Environment) {
	l.items = items
}

// Items returns the items of the list.
func (l *EnvironmentList) Items() []*Environment {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *EnvironmentList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *EnvironmentList) Get(i int) *Environment {
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
func (l *EnvironmentList) Slice() []*Environment {
	var slice []*Environment
	if l == nil {
		slice = make([]*Environment, 0)
	} else {
		slice = make([]*Environment, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *EnvironmentList) Each(f func(item *Environment) bool) {
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
func (l *EnvironmentList) Range(f func(index int, item *Environment) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
