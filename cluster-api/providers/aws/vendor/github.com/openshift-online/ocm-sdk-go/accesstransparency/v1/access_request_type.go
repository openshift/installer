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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	time "time"
)

// AccessRequestKind is the name of the type used to represent objects
// of type 'access_request'.
const AccessRequestKind = "AccessRequest"

// AccessRequestLinkKind is the name of the type used to represent links
// to objects of type 'access_request'.
const AccessRequestLinkKind = "AccessRequestLink"

// AccessRequestNilKind is the name of the type used to nil references
// to objects of type 'access_request'.
const AccessRequestNilKind = "AccessRequestNil"

// AccessRequest represents the values of the 'access_request' type.
//
// Representation of an access request.
type AccessRequest struct {
	bitmap_               uint32
	id                    string
	href                  string
	clusterId             string
	createdAt             time.Time
	deadline              string
	deadlineAt            time.Time
	decisions             []*Decision
	duration              string
	internalSupportCaseId string
	justification         string
	organizationId        string
	requestedBy           string
	status                *AccessRequestStatus
	subscriptionId        string
	supportCaseId         string
	updatedAt             time.Time
}

// Kind returns the name of the type of the object.
func (o *AccessRequest) Kind() string {
	if o == nil {
		return AccessRequestNilKind
	}
	if o.bitmap_&1 != 0 {
		return AccessRequestLinkKind
	}
	return AccessRequestKind
}

// Link returns true iif this is a link.
func (o *AccessRequest) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AccessRequest) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AccessRequest) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AccessRequest) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AccessRequest) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccessRequest) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ClusterId returns the value of the 'cluster_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster from which the Access Request belongs to.
func (o *AccessRequest) ClusterId() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clusterId
	}
	return ""
}

// GetClusterId returns the value of the 'cluster_id' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster from which the Access Request belongs to.
func (o *AccessRequest) GetClusterId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clusterId
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the access request was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) CreatedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the access request was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Deadline returns the value of the 'deadline' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long the Access Request can be in pending state waiting for a customer decision.
func (o *AccessRequest) Deadline() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.deadline
	}
	return ""
}

// GetDeadline returns the value of the 'deadline' attribute and
// a flag indicating if the attribute has a value.
//
// How long the Access Request can be in pending state waiting for a customer decision.
func (o *AccessRequest) GetDeadline() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.deadline
	}
	return
}

// DeadlineAt returns the value of the 'deadline_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time for the deadline that the Access Request needs to be decided, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) DeadlineAt() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.deadlineAt
	}
	return time.Time{}
}

// GetDeadlineAt returns the value of the 'deadline_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time for the deadline that the Access Request needs to be decided, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) GetDeadlineAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.deadlineAt
	}
	return
}

// Decisions returns the value of the 'decisions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Decisions attached to the Access Request.
func (o *AccessRequest) Decisions() []*Decision {
	if o != nil && o.bitmap_&128 != 0 {
		return o.decisions
	}
	return nil
}

// GetDecisions returns the value of the 'decisions' attribute and
// a flag indicating if the attribute has a value.
//
// Decisions attached to the Access Request.
func (o *AccessRequest) GetDecisions() (value []*Decision, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.decisions
	}
	return
}

// Duration returns the value of the 'duration' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long the access will last after it's been approved.
func (o *AccessRequest) Duration() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.duration
	}
	return ""
}

// GetDuration returns the value of the 'duration' attribute and
// a flag indicating if the attribute has a value.
//
// How long the access will last after it's been approved.
func (o *AccessRequest) GetDuration() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.duration
	}
	return
}

// InternalSupportCaseId returns the value of the 'internal_support_case_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Internal support case id linking to jira ticket.
func (o *AccessRequest) InternalSupportCaseId() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.internalSupportCaseId
	}
	return ""
}

// GetInternalSupportCaseId returns the value of the 'internal_support_case_id' attribute and
// a flag indicating if the attribute has a value.
//
// Internal support case id linking to jira ticket.
func (o *AccessRequest) GetInternalSupportCaseId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.internalSupportCaseId
	}
	return
}

// Justification returns the value of the 'justification' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Justification of the Access Request.
func (o *AccessRequest) Justification() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.justification
	}
	return ""
}

// GetJustification returns the value of the 'justification' attribute and
// a flag indicating if the attribute has a value.
//
// Justification of the Access Request.
func (o *AccessRequest) GetJustification() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.justification
	}
	return
}

// OrganizationId returns the value of the 'organization_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Organization from which the Access Request belongs to.
func (o *AccessRequest) OrganizationId() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.organizationId
	}
	return ""
}

// GetOrganizationId returns the value of the 'organization_id' attribute and
// a flag indicating if the attribute has a value.
//
// Organization from which the Access Request belongs to.
func (o *AccessRequest) GetOrganizationId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.organizationId
	}
	return
}

// RequestedBy returns the value of the 'requested_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// User that requested the Access.
func (o *AccessRequest) RequestedBy() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.requestedBy
	}
	return ""
}

// GetRequestedBy returns the value of the 'requested_by' attribute and
// a flag indicating if the attribute has a value.
//
// User that requested the Access.
func (o *AccessRequest) GetRequestedBy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.requestedBy
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Access Request status.
func (o *AccessRequest) Status() *AccessRequestStatus {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Access Request status.
func (o *AccessRequest) GetStatus() (value *AccessRequestStatus, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.status
	}
	return
}

// SubscriptionId returns the value of the 'subscription_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Subscription from which the Access Request belongs to.
func (o *AccessRequest) SubscriptionId() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.subscriptionId
	}
	return ""
}

// GetSubscriptionId returns the value of the 'subscription_id' attribute and
// a flag indicating if the attribute has a value.
//
// Subscription from which the Access Request belongs to.
func (o *AccessRequest) GetSubscriptionId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.subscriptionId
	}
	return
}

// SupportCaseId returns the value of the 'support_case_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Support case ID linking to JIRA ticket.
func (o *AccessRequest) SupportCaseId() string {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.supportCaseId
	}
	return ""
}

// GetSupportCaseId returns the value of the 'support_case_id' attribute and
// a flag indicating if the attribute has a value.
//
// Support case ID linking to JIRA ticket.
func (o *AccessRequest) GetSupportCaseId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.supportCaseId
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the access request was lastly updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the access request was lastly updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequest) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// AccessRequestListKind is the name of the type used to represent list of objects of
// type 'access_request'.
const AccessRequestListKind = "AccessRequestList"

// AccessRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'access_request'.
const AccessRequestListLinkKind = "AccessRequestListLink"

// AccessRequestNilKind is the name of the type used to nil lists of objects of
// type 'access_request'.
const AccessRequestListNilKind = "AccessRequestListNil"

// AccessRequestList is a list of values of the 'access_request' type.
type AccessRequestList struct {
	href  string
	link  bool
	items []*AccessRequest
}

// Kind returns the name of the type of the object.
func (l *AccessRequestList) Kind() string {
	if l == nil {
		return AccessRequestListNilKind
	}
	if l.link {
		return AccessRequestListLinkKind
	}
	return AccessRequestListKind
}

// Link returns true iif this is a link.
func (l *AccessRequestList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AccessRequestList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AccessRequestList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AccessRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AccessRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccessRequestList) Get(i int) *AccessRequest {
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
func (l *AccessRequestList) Slice() []*AccessRequest {
	var slice []*AccessRequest
	if l == nil {
		slice = make([]*AccessRequest, 0)
	} else {
		slice = make([]*AccessRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccessRequestList) Each(f func(item *AccessRequest) bool) {
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
func (l *AccessRequestList) Range(f func(index int, item *AccessRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
