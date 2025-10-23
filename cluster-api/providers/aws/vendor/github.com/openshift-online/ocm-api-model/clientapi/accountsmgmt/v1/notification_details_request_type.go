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

// NotificationDetailsRequest represents the values of the 'notification_details_request' type.
//
// This struct is a request to get a templated email to a user related to this.
// subscription/cluster.
type NotificationDetailsRequest struct {
	fieldSet_               []bool
	bccAddress              string
	clusterID               string
	clusterUUID             string
	logType                 string
	subject                 string
	subscriptionID          string
	includeRedHatAssociates bool
	internalOnly            bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NotificationDetailsRequest) Empty() bool {
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

// BccAddress returns the value of the 'bcc_address' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The BCC address to be included on the email that is sent.
func (o *NotificationDetailsRequest) BccAddress() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.bccAddress
	}
	return ""
}

// GetBccAddress returns the value of the 'bcc_address' attribute and
// a flag indicating if the attribute has a value.
//
// The BCC address to be included on the email that is sent.
func (o *NotificationDetailsRequest) GetBccAddress() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.bccAddress
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Cluster (internal id) the resource type belongs to.
func (o *NotificationDetailsRequest) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Cluster (internal id) the resource type belongs to.
func (o *NotificationDetailsRequest) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.clusterID
	}
	return
}

// ClusterUUID returns the value of the 'cluster_UUID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Cluster (external id) the resource type belongs to.
func (o *NotificationDetailsRequest) ClusterUUID() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.clusterUUID
	}
	return ""
}

// GetClusterUUID returns the value of the 'cluster_UUID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Cluster (external id) the resource type belongs to.
func (o *NotificationDetailsRequest) GetClusterUUID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.clusterUUID
	}
	return
}

// IncludeRedHatAssociates returns the value of the 'include_red_hat_associates' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether to include red hat associates in the email notification.
func (o *NotificationDetailsRequest) IncludeRedHatAssociates() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.includeRedHatAssociates
	}
	return false
}

// GetIncludeRedHatAssociates returns the value of the 'include_red_hat_associates' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether to include red hat associates in the email notification.
func (o *NotificationDetailsRequest) GetIncludeRedHatAssociates() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.includeRedHatAssociates
	}
	return
}

// InternalOnly returns the value of the 'internal_only' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the service log is internal only to RH.
func (o *NotificationDetailsRequest) InternalOnly() bool {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.internalOnly
	}
	return false
}

// GetInternalOnly returns the value of the 'internal_only' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the service log is internal only to RH.
func (o *NotificationDetailsRequest) GetInternalOnly() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.internalOnly
	}
	return
}

// LogType returns the value of the 'log_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the type of the service log.
func (o *NotificationDetailsRequest) LogType() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.logType
	}
	return ""
}

// GetLogType returns the value of the 'log_type' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the type of the service log.
func (o *NotificationDetailsRequest) GetLogType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.logType
	}
	return
}

// Subject returns the value of the 'subject' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The email subject.
func (o *NotificationDetailsRequest) Subject() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.subject
	}
	return ""
}

// GetSubject returns the value of the 'subject' attribute and
// a flag indicating if the attribute has a value.
//
// The email subject.
func (o *NotificationDetailsRequest) GetSubject() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.subject
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Subscription the resource type belongs to.
func (o *NotificationDetailsRequest) SubscriptionID() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Subscription the resource type belongs to.
func (o *NotificationDetailsRequest) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.subscriptionID
	}
	return
}

// NotificationDetailsRequestListKind is the name of the type used to represent list of objects of
// type 'notification_details_request'.
const NotificationDetailsRequestListKind = "NotificationDetailsRequestList"

// NotificationDetailsRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'notification_details_request'.
const NotificationDetailsRequestListLinkKind = "NotificationDetailsRequestListLink"

// NotificationDetailsRequestNilKind is the name of the type used to nil lists of objects of
// type 'notification_details_request'.
const NotificationDetailsRequestListNilKind = "NotificationDetailsRequestListNil"

// NotificationDetailsRequestList is a list of values of the 'notification_details_request' type.
type NotificationDetailsRequestList struct {
	href  string
	link  bool
	items []*NotificationDetailsRequest
}

// Len returns the length of the list.
func (l *NotificationDetailsRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NotificationDetailsRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NotificationDetailsRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NotificationDetailsRequestList) SetItems(items []*NotificationDetailsRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *NotificationDetailsRequestList) Items() []*NotificationDetailsRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NotificationDetailsRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NotificationDetailsRequestList) Get(i int) *NotificationDetailsRequest {
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
func (l *NotificationDetailsRequestList) Slice() []*NotificationDetailsRequest {
	var slice []*NotificationDetailsRequest
	if l == nil {
		slice = make([]*NotificationDetailsRequest, 0)
	} else {
		slice = make([]*NotificationDetailsRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NotificationDetailsRequestList) Each(f func(item *NotificationDetailsRequest) bool) {
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
func (l *NotificationDetailsRequestList) Range(f func(index int, item *NotificationDetailsRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
