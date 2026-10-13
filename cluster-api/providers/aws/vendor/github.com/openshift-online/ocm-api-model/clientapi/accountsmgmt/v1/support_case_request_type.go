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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// SupportCaseRequestKind is the name of the type used to represent objects
// of type 'support_case_request'.
const SupportCaseRequestKind = "SupportCaseRequest"

// SupportCaseRequestLinkKind is the name of the type used to represent links
// to objects of type 'support_case_request'.
const SupportCaseRequestLinkKind = "SupportCaseRequestLink"

// SupportCaseRequestNilKind is the name of the type used to nil references
// to objects of type 'support_case_request'.
const SupportCaseRequestNilKind = "SupportCaseRequestNil"

// SupportCaseRequest represents the values of the 'support_case_request' type.
type SupportCaseRequest struct {
	fieldSet_      []bool
	id             string
	href           string
	clusterId      string
	clusterUuid    string
	description    string
	eventStreamId  string
	severity       string
	subscriptionId string
	summary        string
}

// Kind returns the name of the type of the object.
func (o *SupportCaseRequest) Kind() string {
	if o == nil {
		return SupportCaseRequestNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return SupportCaseRequestLinkKind
	}
	return SupportCaseRequestKind
}

// Link returns true if this is a link.
func (o *SupportCaseRequest) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *SupportCaseRequest) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *SupportCaseRequest) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *SupportCaseRequest) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *SupportCaseRequest) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SupportCaseRequest) Empty() bool {
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

// ClusterId returns the value of the 'cluster_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (optional) cluster id of the cluster on which we create the support case for.
func (o *SupportCaseRequest) ClusterId() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.clusterId
	}
	return ""
}

// GetClusterId returns the value of the 'cluster_id' attribute and
// a flag indicating if the attribute has a value.
//
// (optional) cluster id of the cluster on which we create the support case for.
func (o *SupportCaseRequest) GetClusterId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.clusterId
	}
	return
}

// ClusterUuid returns the value of the 'cluster_uuid' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (optional) cluster uuid of the cluster on which we create the support case for.
func (o *SupportCaseRequest) ClusterUuid() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.clusterUuid
	}
	return ""
}

// GetClusterUuid returns the value of the 'cluster_uuid' attribute and
// a flag indicating if the attribute has a value.
//
// (optional) cluster uuid of the cluster on which we create the support case for.
func (o *SupportCaseRequest) GetClusterUuid() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.clusterUuid
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Support case desciption.
func (o *SupportCaseRequest) Description() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Support case desciption.
func (o *SupportCaseRequest) GetDescription() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.description
	}
	return
}

// EventStreamId returns the value of the 'event_stream_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (optional) event stream id for the support case so we could track it.
func (o *SupportCaseRequest) EventStreamId() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.eventStreamId
	}
	return ""
}

// GetEventStreamId returns the value of the 'event_stream_id' attribute and
// a flag indicating if the attribute has a value.
//
// (optional) event stream id for the support case so we could track it.
func (o *SupportCaseRequest) GetEventStreamId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.eventStreamId
	}
	return
}

// Severity returns the value of the 'severity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Support case severity.
func (o *SupportCaseRequest) Severity() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.severity
	}
	return ""
}

// GetSeverity returns the value of the 'severity' attribute and
// a flag indicating if the attribute has a value.
//
// Support case severity.
func (o *SupportCaseRequest) GetSeverity() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.severity
	}
	return
}

// SubscriptionId returns the value of the 'subscription_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (optional) subscription id of the subscription on which we create the support case for.
func (o *SupportCaseRequest) SubscriptionId() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.subscriptionId
	}
	return ""
}

// GetSubscriptionId returns the value of the 'subscription_id' attribute and
// a flag indicating if the attribute has a value.
//
// (optional) subscription id of the subscription on which we create the support case for.
func (o *SupportCaseRequest) GetSubscriptionId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.subscriptionId
	}
	return
}

// Summary returns the value of the 'summary' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Support case title.
func (o *SupportCaseRequest) Summary() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.summary
	}
	return ""
}

// GetSummary returns the value of the 'summary' attribute and
// a flag indicating if the attribute has a value.
//
// Support case title.
func (o *SupportCaseRequest) GetSummary() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.summary
	}
	return
}

// SupportCaseRequestListKind is the name of the type used to represent list of objects of
// type 'support_case_request'.
const SupportCaseRequestListKind = "SupportCaseRequestList"

// SupportCaseRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'support_case_request'.
const SupportCaseRequestListLinkKind = "SupportCaseRequestListLink"

// SupportCaseRequestNilKind is the name of the type used to nil lists of objects of
// type 'support_case_request'.
const SupportCaseRequestListNilKind = "SupportCaseRequestListNil"

// SupportCaseRequestList is a list of values of the 'support_case_request' type.
type SupportCaseRequestList struct {
	href  string
	link  bool
	items []*SupportCaseRequest
}

// Kind returns the name of the type of the object.
func (l *SupportCaseRequestList) Kind() string {
	if l == nil {
		return SupportCaseRequestListNilKind
	}
	if l.link {
		return SupportCaseRequestListLinkKind
	}
	return SupportCaseRequestListKind
}

// Link returns true iif this is a link.
func (l *SupportCaseRequestList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *SupportCaseRequestList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *SupportCaseRequestList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *SupportCaseRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SupportCaseRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SupportCaseRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SupportCaseRequestList) SetItems(items []*SupportCaseRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *SupportCaseRequestList) Items() []*SupportCaseRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SupportCaseRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SupportCaseRequestList) Get(i int) *SupportCaseRequest {
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
func (l *SupportCaseRequestList) Slice() []*SupportCaseRequest {
	var slice []*SupportCaseRequest
	if l == nil {
		slice = make([]*SupportCaseRequest, 0)
	} else {
		slice = make([]*SupportCaseRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SupportCaseRequestList) Each(f func(item *SupportCaseRequest) bool) {
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
func (l *SupportCaseRequestList) Range(f func(index int, item *SupportCaseRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
