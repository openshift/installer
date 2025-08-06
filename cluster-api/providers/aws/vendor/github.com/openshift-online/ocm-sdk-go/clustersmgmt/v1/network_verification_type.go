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

// NetworkVerification represents the values of the 'network_verification' type.
type NetworkVerification struct {
	bitmap_           uint32
	cloudProviderData *CloudProviderData
	clusterId         string
	items             []*SubnetNetworkVerification
	platform          Platform
	total             int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NetworkVerification) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CloudProviderData returns the value of the 'cloud_provider_data' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cloud provider data to execute the network verification.
func (o *NetworkVerification) CloudProviderData() *CloudProviderData {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cloudProviderData
	}
	return nil
}

// GetCloudProviderData returns the value of the 'cloud_provider_data' attribute and
// a flag indicating if the attribute has a value.
//
// Cloud provider data to execute the network verification.
func (o *NetworkVerification) GetCloudProviderData() (value *CloudProviderData, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cloudProviderData
	}
	return
}

// ClusterId returns the value of the 'cluster_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster ID needed to execute the network verification.
func (o *NetworkVerification) ClusterId() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.clusterId
	}
	return ""
}

// GetClusterId returns the value of the 'cluster_id' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster ID needed to execute the network verification.
func (o *NetworkVerification) GetClusterId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.clusterId
	}
	return
}

// Items returns the value of the 'items' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details about each subnet network verification.
func (o *NetworkVerification) Items() []*SubnetNetworkVerification {
	if o != nil && o.bitmap_&4 != 0 {
		return o.items
	}
	return nil
}

// GetItems returns the value of the 'items' attribute and
// a flag indicating if the attribute has a value.
//
// Details about each subnet network verification.
func (o *NetworkVerification) GetItems() (value []*SubnetNetworkVerification, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.items
	}
	return
}

// Platform returns the value of the 'platform' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Platform needed to execute the network verification.
func (o *NetworkVerification) Platform() Platform {
	if o != nil && o.bitmap_&8 != 0 {
		return o.platform
	}
	return Platform("")
}

// GetPlatform returns the value of the 'platform' attribute and
// a flag indicating if the attribute has a value.
//
// Platform needed to execute the network verification.
func (o *NetworkVerification) GetPlatform() (value Platform, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.platform
	}
	return
}

// Total returns the value of the 'total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Amount of network verifier executions started.
func (o *NetworkVerification) Total() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.total
	}
	return 0
}

// GetTotal returns the value of the 'total' attribute and
// a flag indicating if the attribute has a value.
//
// Amount of network verifier executions started.
func (o *NetworkVerification) GetTotal() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.total
	}
	return
}

// NetworkVerificationListKind is the name of the type used to represent list of objects of
// type 'network_verification'.
const NetworkVerificationListKind = "NetworkVerificationList"

// NetworkVerificationListLinkKind is the name of the type used to represent links to list
// of objects of type 'network_verification'.
const NetworkVerificationListLinkKind = "NetworkVerificationListLink"

// NetworkVerificationNilKind is the name of the type used to nil lists of objects of
// type 'network_verification'.
const NetworkVerificationListNilKind = "NetworkVerificationListNil"

// NetworkVerificationList is a list of values of the 'network_verification' type.
type NetworkVerificationList struct {
	href  string
	link  bool
	items []*NetworkVerification
}

// Len returns the length of the list.
func (l *NetworkVerificationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NetworkVerificationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NetworkVerificationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NetworkVerificationList) SetItems(items []*NetworkVerification) {
	l.items = items
}

// Items returns the items of the list.
func (l *NetworkVerificationList) Items() []*NetworkVerification {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NetworkVerificationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NetworkVerificationList) Get(i int) *NetworkVerification {
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
func (l *NetworkVerificationList) Slice() []*NetworkVerification {
	var slice []*NetworkVerification
	if l == nil {
		slice = make([]*NetworkVerification, 0)
	} else {
		slice = make([]*NetworkVerification, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NetworkVerificationList) Each(f func(item *NetworkVerification) bool) {
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
func (l *NetworkVerificationList) Range(f func(index int, item *NetworkVerification) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
