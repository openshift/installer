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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	time "time"
)

// ManagedServiceKind is the name of the type used to represent objects
// of type 'managed_service'.
const ManagedServiceKind = "ManagedService"

// ManagedServiceLinkKind is the name of the type used to represent links
// to objects of type 'managed_service'.
const ManagedServiceLinkKind = "ManagedServiceLink"

// ManagedServiceNilKind is the name of the type used to nil references
// to objects of type 'managed_service'.
const ManagedServiceNilKind = "ManagedServiceNil"

// ManagedService represents the values of the 'managed_service' type.
//
// Represents data about a running Managed Service.
type ManagedService struct {
	bitmap_      uint32
	id           string
	href         string
	addon        *StatefulObject
	cluster      *Cluster
	createdAt    time.Time
	expiredAt    time.Time
	parameters   []*ServiceParameter
	resources    []*StatefulObject
	service      string
	serviceState string
	updatedAt    time.Time
}

// Kind returns the name of the type of the object.
func (o *ManagedService) Kind() string {
	if o == nil {
		return ManagedServiceNilKind
	}
	if o.bitmap_&1 != 0 {
		return ManagedServiceLinkKind
	}
	return ManagedServiceKind
}

// Link returns true if this is a link.
func (o *ManagedService) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ManagedService) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ManagedService) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ManagedService) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ManagedService) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ManagedService) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Addon returns the value of the 'addon' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) Addon() *StatefulObject {
	if o != nil && o.bitmap_&8 != 0 {
		return o.addon
	}
	return nil
}

// GetAddon returns the value of the 'addon' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetAddon() (value *StatefulObject, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.addon
	}
	return
}

// Cluster returns the value of the 'cluster' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) Cluster() *Cluster {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cluster
	}
	return nil
}

// GetCluster returns the value of the 'cluster' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetCluster() (value *Cluster, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cluster
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) CreatedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// ExpiredAt returns the value of the 'expired_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) ExpiredAt() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.expiredAt
	}
	return time.Time{}
}

// GetExpiredAt returns the value of the 'expired_at' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetExpiredAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.expiredAt
	}
	return
}

// Parameters returns the value of the 'parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) Parameters() []*ServiceParameter {
	if o != nil && o.bitmap_&128 != 0 {
		return o.parameters
	}
	return nil
}

// GetParameters returns the value of the 'parameters' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetParameters() (value []*ServiceParameter, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.parameters
	}
	return
}

// Resources returns the value of the 'resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) Resources() []*StatefulObject {
	if o != nil && o.bitmap_&256 != 0 {
		return o.resources
	}
	return nil
}

// GetResources returns the value of the 'resources' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetResources() (value []*StatefulObject, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.resources
	}
	return
}

// Service returns the value of the 'service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) Service() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.service
	}
	return ""
}

// GetService returns the value of the 'service' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetService() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.service
	}
	return
}

// ServiceState returns the value of the 'service_state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) ServiceState() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.serviceState
	}
	return ""
}

// GetServiceState returns the value of the 'service_state' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetServiceState() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.serviceState
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ManagedService) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *ManagedService) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// ManagedServiceListKind is the name of the type used to represent list of objects of
// type 'managed_service'.
const ManagedServiceListKind = "ManagedServiceList"

// ManagedServiceListLinkKind is the name of the type used to represent links to list
// of objects of type 'managed_service'.
const ManagedServiceListLinkKind = "ManagedServiceListLink"

// ManagedServiceNilKind is the name of the type used to nil lists of objects of
// type 'managed_service'.
const ManagedServiceListNilKind = "ManagedServiceListNil"

// ManagedServiceList is a list of values of the 'managed_service' type.
type ManagedServiceList struct {
	href  string
	link  bool
	items []*ManagedService
}

// Kind returns the name of the type of the object.
func (l *ManagedServiceList) Kind() string {
	if l == nil {
		return ManagedServiceListNilKind
	}
	if l.link {
		return ManagedServiceListLinkKind
	}
	return ManagedServiceListKind
}

// Link returns true iif this is a link.
func (l *ManagedServiceList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ManagedServiceList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ManagedServiceList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ManagedServiceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetItems(items []*ManagedService) {
	l.items = items
}

// Items returns the items of the list.
func (l *ManagedServiceList) Items() []*ManagedService {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ManagedServiceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ManagedServiceList) Get(i int) *ManagedService {
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
func (l *ManagedServiceList) Slice() []*ManagedService {
	var slice []*ManagedService
	if l == nil {
		slice = make([]*ManagedService, 0)
	} else {
		slice = make([]*ManagedService, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ManagedServiceList) Each(f func(item *ManagedService) bool) {
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
func (l *ManagedServiceList) Range(f func(index int, item *ManagedService) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
