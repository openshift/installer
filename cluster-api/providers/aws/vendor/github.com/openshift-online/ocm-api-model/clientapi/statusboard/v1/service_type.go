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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

import (
	time "time"
)

// ServiceKind is the name of the type used to represent objects
// of type 'service'.
const ServiceKind = "Service"

// ServiceLinkKind is the name of the type used to represent links
// to objects of type 'service'.
const ServiceLinkKind = "ServiceLink"

// ServiceNilKind is the name of the type used to nil references
// to objects of type 'service'.
const ServiceNilKind = "ServiceNil"

// Service represents the values of the 'service' type.
//
// Definition of a Status Board Service.
type Service struct {
	fieldSet_       []bool
	id              string
	href            string
	application     *Application
	createdAt       time.Time
	currentStatus   string
	fullname        string
	lastPingAt      time.Time
	metadata        interface{}
	name            string
	owners          []*Owner
	serviceEndpoint string
	statusType      string
	statusUpdatedAt time.Time
	token           string
	updatedAt       time.Time
	private         bool
}

// Kind returns the name of the type of the object.
func (o *Service) Kind() string {
	if o == nil {
		return ServiceNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ServiceLinkKind
	}
	return ServiceKind
}

// Link returns true if this is a link.
func (o *Service) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Service) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Service) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Service) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Service) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Service) Empty() bool {
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

// Application returns the value of the 'application' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Application associated with the Service
func (o *Service) Application() *Application {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.application
	}
	return nil
}

// GetApplication returns the value of the 'application' attribute and
// a flag indicating if the attribute has a value.
//
// The Application associated with the Service
func (o *Service) GetApplication() (value *Application, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.application
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Service) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Service) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.createdAt
	}
	return
}

// CurrentStatus returns the value of the 'current_status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) CurrentStatus() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.currentStatus
	}
	return ""
}

// GetCurrentStatus returns the value of the 'current_status' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetCurrentStatus() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.currentStatus
	}
	return
}

// Fullname returns the value of the 'fullname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) Fullname() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.fullname
	}
	return ""
}

// GetFullname returns the value of the 'fullname' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetFullname() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.fullname
	}
	return
}

// LastPingAt returns the value of the 'last_ping_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) LastPingAt() time.Time {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.lastPingAt
	}
	return time.Time{}
}

// GetLastPingAt returns the value of the 'last_ping_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetLastPingAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.lastPingAt
	}
	return
}

// Metadata returns the value of the 'metadata' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) Metadata() interface{} {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.metadata
	}
	return nil
}

// GetMetadata returns the value of the 'metadata' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetMetadata() (value interface{}, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.metadata
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the Service
func (o *Service) Name() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the Service
func (o *Service) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.name
	}
	return
}

// Owners returns the value of the 'owners' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) Owners() []*Owner {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.owners
	}
	return nil
}

// GetOwners returns the value of the 'owners' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetOwners() (value []*Owner, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.owners
	}
	return
}

// Private returns the value of the 'private' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) Private() bool {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.private
	}
	return false
}

// GetPrivate returns the value of the 'private' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetPrivate() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.private
	}
	return
}

// ServiceEndpoint returns the value of the 'service_endpoint' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) ServiceEndpoint() string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.serviceEndpoint
	}
	return ""
}

// GetServiceEndpoint returns the value of the 'service_endpoint' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetServiceEndpoint() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.serviceEndpoint
	}
	return
}

// StatusType returns the value of the 'status_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) StatusType() string {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.statusType
	}
	return ""
}

// GetStatusType returns the value of the 'status_type' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetStatusType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.statusType
	}
	return
}

// StatusUpdatedAt returns the value of the 'status_updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) StatusUpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14] {
		return o.statusUpdatedAt
	}
	return time.Time{}
}

// GetStatusUpdatedAt returns the value of the 'status_updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetStatusUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14]
	if ok {
		value = o.statusUpdatedAt
	}
	return
}

// Token returns the value of the 'token' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Service) Token() string {
	if o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15] {
		return o.token
	}
	return ""
}

// GetToken returns the value of the 'token' attribute and
// a flag indicating if the attribute has a value.
func (o *Service) GetToken() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15]
	if ok {
		value = o.token
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Service) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Service) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16]
	if ok {
		value = o.updatedAt
	}
	return
}

// ServiceListKind is the name of the type used to represent list of objects of
// type 'service'.
const ServiceListKind = "ServiceList"

// ServiceListLinkKind is the name of the type used to represent links to list
// of objects of type 'service'.
const ServiceListLinkKind = "ServiceListLink"

// ServiceNilKind is the name of the type used to nil lists of objects of
// type 'service'.
const ServiceListNilKind = "ServiceListNil"

// ServiceList is a list of values of the 'service' type.
type ServiceList struct {
	href  string
	link  bool
	items []*Service
}

// Kind returns the name of the type of the object.
func (l *ServiceList) Kind() string {
	if l == nil {
		return ServiceListNilKind
	}
	if l.link {
		return ServiceListLinkKind
	}
	return ServiceListKind
}

// Link returns true iif this is a link.
func (l *ServiceList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ServiceList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ServiceList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ServiceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ServiceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ServiceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ServiceList) SetItems(items []*Service) {
	l.items = items
}

// Items returns the items of the list.
func (l *ServiceList) Items() []*Service {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ServiceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ServiceList) Get(i int) *Service {
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
func (l *ServiceList) Slice() []*Service {
	var slice []*Service
	if l == nil {
		slice = make([]*Service, 0)
	} else {
		slice = make([]*Service, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ServiceList) Each(f func(item *Service) bool) {
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
func (l *ServiceList) Range(f func(index int, item *Service) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
