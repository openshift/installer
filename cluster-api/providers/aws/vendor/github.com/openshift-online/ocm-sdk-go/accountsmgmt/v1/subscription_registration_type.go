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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// SubscriptionRegistration represents the values of the 'subscription_registration' type.
//
// Registration of a new subscription.
type SubscriptionRegistration struct {
	bitmap_     uint32
	clusterUUID string
	consoleURL  string
	displayName string
	planID      PlanID
	status      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SubscriptionRegistration) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ClusterUUID returns the value of the 'cluster_UUID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// External cluster ID.
func (o *SubscriptionRegistration) ClusterUUID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.clusterUUID
	}
	return ""
}

// GetClusterUUID returns the value of the 'cluster_UUID' attribute and
// a flag indicating if the attribute has a value.
//
// External cluster ID.
func (o *SubscriptionRegistration) GetClusterUUID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.clusterUUID
	}
	return
}

// ConsoleURL returns the value of the 'console_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Console URL of subscription (optional).
func (o *SubscriptionRegistration) ConsoleURL() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.consoleURL
	}
	return ""
}

// GetConsoleURL returns the value of the 'console_URL' attribute and
// a flag indicating if the attribute has a value.
//
// Console URL of subscription (optional).
func (o *SubscriptionRegistration) GetConsoleURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.consoleURL
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Display name of subscription (optional).
func (o *SubscriptionRegistration) DisplayName() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// Display name of subscription (optional).
func (o *SubscriptionRegistration) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.displayName
	}
	return
}

// PlanID returns the value of the 'plan_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Plan ID of subscription.
func (o *SubscriptionRegistration) PlanID() PlanID {
	if o != nil && o.bitmap_&8 != 0 {
		return o.planID
	}
	return PlanID("")
}

// GetPlanID returns the value of the 'plan_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Plan ID of subscription.
func (o *SubscriptionRegistration) GetPlanID() (value PlanID, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.planID
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of subscription.
func (o *SubscriptionRegistration) Status() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of subscription.
func (o *SubscriptionRegistration) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.status
	}
	return
}

// SubscriptionRegistrationListKind is the name of the type used to represent list of objects of
// type 'subscription_registration'.
const SubscriptionRegistrationListKind = "SubscriptionRegistrationList"

// SubscriptionRegistrationListLinkKind is the name of the type used to represent links to list
// of objects of type 'subscription_registration'.
const SubscriptionRegistrationListLinkKind = "SubscriptionRegistrationListLink"

// SubscriptionRegistrationNilKind is the name of the type used to nil lists of objects of
// type 'subscription_registration'.
const SubscriptionRegistrationListNilKind = "SubscriptionRegistrationListNil"

// SubscriptionRegistrationList is a list of values of the 'subscription_registration' type.
type SubscriptionRegistrationList struct {
	href  string
	link  bool
	items []*SubscriptionRegistration
}

// Len returns the length of the list.
func (l *SubscriptionRegistrationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *SubscriptionRegistrationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SubscriptionRegistrationList) Get(i int) *SubscriptionRegistration {
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
func (l *SubscriptionRegistrationList) Slice() []*SubscriptionRegistration {
	var slice []*SubscriptionRegistration
	if l == nil {
		slice = make([]*SubscriptionRegistration, 0)
	} else {
		slice = make([]*SubscriptionRegistration, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SubscriptionRegistrationList) Each(f func(item *SubscriptionRegistration) bool) {
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
func (l *SubscriptionRegistrationList) Range(f func(index int, item *SubscriptionRegistration) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
