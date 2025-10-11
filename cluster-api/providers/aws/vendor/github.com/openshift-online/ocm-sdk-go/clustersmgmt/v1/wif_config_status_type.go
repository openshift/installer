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

// WifConfigStatus represents the values of the 'wif_config_status' type.
//
// Configuration status of a WifConfig.
type WifConfigStatus struct {
	bitmap_     uint32
	description string
	configured  bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifConfigStatus) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Configured returns the value of the 'configured' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the current status of the WifConfig resource configuration.
// - `false`: The WifConfig resource has a user configuration error.
// - `true`: The resources associated with the WifConfig object are properly configured and operational at the time of the check.
func (o *WifConfigStatus) Configured() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.configured
	}
	return false
}

// GetConfigured returns the value of the 'configured' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the current status of the WifConfig resource configuration.
// - `false`: The WifConfig resource has a user configuration error.
// - `true`: The resources associated with the WifConfig object are properly configured and operational at the time of the check.
func (o *WifConfigStatus) GetConfigured() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.configured
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Provides additional information about the WifConfig resource status.
// - When `Configured` is `false`, this field contains details about the user configuration error.
// - When `Configured` is `true`, this field may be empty or contain optional notes about the configuration.
func (o *WifConfigStatus) Description() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Provides additional information about the WifConfig resource status.
// - When `Configured` is `false`, this field contains details about the user configuration error.
// - When `Configured` is `true`, this field may be empty or contain optional notes about the configuration.
func (o *WifConfigStatus) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.description
	}
	return
}

// WifConfigStatusListKind is the name of the type used to represent list of objects of
// type 'wif_config_status'.
const WifConfigStatusListKind = "WifConfigStatusList"

// WifConfigStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_config_status'.
const WifConfigStatusListLinkKind = "WifConfigStatusListLink"

// WifConfigStatusNilKind is the name of the type used to nil lists of objects of
// type 'wif_config_status'.
const WifConfigStatusListNilKind = "WifConfigStatusListNil"

// WifConfigStatusList is a list of values of the 'wif_config_status' type.
type WifConfigStatusList struct {
	href  string
	link  bool
	items []*WifConfigStatus
}

// Len returns the length of the list.
func (l *WifConfigStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *WifConfigStatusList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *WifConfigStatusList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *WifConfigStatusList) SetItems(items []*WifConfigStatus) {
	l.items = items
}

// Items returns the items of the list.
func (l *WifConfigStatusList) Items() []*WifConfigStatus {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *WifConfigStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifConfigStatusList) Get(i int) *WifConfigStatus {
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
func (l *WifConfigStatusList) Slice() []*WifConfigStatus {
	var slice []*WifConfigStatus
	if l == nil {
		slice = make([]*WifConfigStatus, 0)
	} else {
		slice = make([]*WifConfigStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifConfigStatusList) Each(f func(item *WifConfigStatus) bool) {
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
func (l *WifConfigStatusList) Range(f func(index int, item *WifConfigStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
