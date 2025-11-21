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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accesstransparency/v1

// AccessRequestPostRequest represents the values of the 'access_request_post_request' type.
//
// Representation of an access request post request.
type AccessRequestPostRequest struct {
	fieldSet_             []bool
	clusterId             string
	deadline              string
	duration              string
	internalSupportCaseId string
	justification         string
	subscriptionId        string
	supportCaseId         string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccessRequestPostRequest) Empty() bool {
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

// ClusterId returns the value of the 'cluster_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster from which the Access Request belongs to.
func (o *AccessRequestPostRequest) ClusterId() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.clusterId
	}
	return ""
}

// GetClusterId returns the value of the 'cluster_id' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster from which the Access Request belongs to.
func (o *AccessRequestPostRequest) GetClusterId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.clusterId
	}
	return
}

// Deadline returns the value of the 'deadline' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long the Access Request can be in pending state waiting for a customer decision.
func (o *AccessRequestPostRequest) Deadline() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.deadline
	}
	return ""
}

// GetDeadline returns the value of the 'deadline' attribute and
// a flag indicating if the attribute has a value.
//
// How long the Access Request can be in pending state waiting for a customer decision.
func (o *AccessRequestPostRequest) GetDeadline() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.deadline
	}
	return
}

// Duration returns the value of the 'duration' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long the access will last after it's been approved.
func (o *AccessRequestPostRequest) Duration() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.duration
	}
	return ""
}

// GetDuration returns the value of the 'duration' attribute and
// a flag indicating if the attribute has a value.
//
// How long the access will last after it's been approved.
func (o *AccessRequestPostRequest) GetDuration() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.duration
	}
	return
}

// InternalSupportCaseId returns the value of the 'internal_support_case_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Internal support case id linking to jira ticket.
func (o *AccessRequestPostRequest) InternalSupportCaseId() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.internalSupportCaseId
	}
	return ""
}

// GetInternalSupportCaseId returns the value of the 'internal_support_case_id' attribute and
// a flag indicating if the attribute has a value.
//
// Internal support case id linking to jira ticket.
func (o *AccessRequestPostRequest) GetInternalSupportCaseId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.internalSupportCaseId
	}
	return
}

// Justification returns the value of the 'justification' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Justification of the Access Request.
func (o *AccessRequestPostRequest) Justification() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.justification
	}
	return ""
}

// GetJustification returns the value of the 'justification' attribute and
// a flag indicating if the attribute has a value.
//
// Justification of the Access Request.
func (o *AccessRequestPostRequest) GetJustification() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.justification
	}
	return
}

// SubscriptionId returns the value of the 'subscription_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Subscription from which the Access Request belongs to.
func (o *AccessRequestPostRequest) SubscriptionId() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.subscriptionId
	}
	return ""
}

// GetSubscriptionId returns the value of the 'subscription_id' attribute and
// a flag indicating if the attribute has a value.
//
// Subscription from which the Access Request belongs to.
func (o *AccessRequestPostRequest) GetSubscriptionId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.subscriptionId
	}
	return
}

// SupportCaseId returns the value of the 'support_case_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Support case ID linking to JIRA ticket.
func (o *AccessRequestPostRequest) SupportCaseId() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.supportCaseId
	}
	return ""
}

// GetSupportCaseId returns the value of the 'support_case_id' attribute and
// a flag indicating if the attribute has a value.
//
// Support case ID linking to JIRA ticket.
func (o *AccessRequestPostRequest) GetSupportCaseId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.supportCaseId
	}
	return
}

// AccessRequestPostRequestListKind is the name of the type used to represent list of objects of
// type 'access_request_post_request'.
const AccessRequestPostRequestListKind = "AccessRequestPostRequestList"

// AccessRequestPostRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'access_request_post_request'.
const AccessRequestPostRequestListLinkKind = "AccessRequestPostRequestListLink"

// AccessRequestPostRequestNilKind is the name of the type used to nil lists of objects of
// type 'access_request_post_request'.
const AccessRequestPostRequestListNilKind = "AccessRequestPostRequestListNil"

// AccessRequestPostRequestList is a list of values of the 'access_request_post_request' type.
type AccessRequestPostRequestList struct {
	href  string
	link  bool
	items []*AccessRequestPostRequest
}

// Len returns the length of the list.
func (l *AccessRequestPostRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AccessRequestPostRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AccessRequestPostRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AccessRequestPostRequestList) SetItems(items []*AccessRequestPostRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *AccessRequestPostRequestList) Items() []*AccessRequestPostRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AccessRequestPostRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccessRequestPostRequestList) Get(i int) *AccessRequestPostRequest {
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
func (l *AccessRequestPostRequestList) Slice() []*AccessRequestPostRequest {
	var slice []*AccessRequestPostRequest
	if l == nil {
		slice = make([]*AccessRequestPostRequest, 0)
	} else {
		slice = make([]*AccessRequestPostRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccessRequestPostRequestList) Each(f func(item *AccessRequestPostRequest) bool) {
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
func (l *AccessRequestPostRequestList) Range(f func(index int, item *AccessRequestPostRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
