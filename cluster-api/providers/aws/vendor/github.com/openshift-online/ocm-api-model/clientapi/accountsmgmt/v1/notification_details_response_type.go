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

// NotificationDetailsResponseKind is the name of the type used to represent objects
// of type 'notification_details_response'.
const NotificationDetailsResponseKind = "NotificationDetailsResponse"

// NotificationDetailsResponseLinkKind is the name of the type used to represent links
// to objects of type 'notification_details_response'.
const NotificationDetailsResponseLinkKind = "NotificationDetailsResponseLink"

// NotificationDetailsResponseNilKind is the name of the type used to nil references
// to objects of type 'notification_details_response'.
const NotificationDetailsResponseNilKind = "NotificationDetailsResponseNil"

// NotificationDetailsResponse represents the values of the 'notification_details_response' type.
//
// This class is a single response item for the notify details list.
type NotificationDetailsResponse struct {
	fieldSet_ []bool
	id        string
	href      string
	key       string
	value     string
}

// Kind returns the name of the type of the object.
func (o *NotificationDetailsResponse) Kind() string {
	if o == nil {
		return NotificationDetailsResponseNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return NotificationDetailsResponseLinkKind
	}
	return NotificationDetailsResponseKind
}

// Link returns true if this is a link.
func (o *NotificationDetailsResponse) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *NotificationDetailsResponse) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NotificationDetailsResponse) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NotificationDetailsResponse) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NotificationDetailsResponse) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NotificationDetailsResponse) Empty() bool {
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

// Key returns the value of the 'key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the key of the response parameter.
func (o *NotificationDetailsResponse) Key() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.key
	}
	return ""
}

// GetKey returns the value of the 'key' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the key of the response parameter.
func (o *NotificationDetailsResponse) GetKey() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.key
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the value of the response parameter.
func (o *NotificationDetailsResponse) Value() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the value of the response parameter.
func (o *NotificationDetailsResponse) GetValue() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.value
	}
	return
}

// NotificationDetailsResponseListKind is the name of the type used to represent list of objects of
// type 'notification_details_response'.
const NotificationDetailsResponseListKind = "NotificationDetailsResponseList"

// NotificationDetailsResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'notification_details_response'.
const NotificationDetailsResponseListLinkKind = "NotificationDetailsResponseListLink"

// NotificationDetailsResponseNilKind is the name of the type used to nil lists of objects of
// type 'notification_details_response'.
const NotificationDetailsResponseListNilKind = "NotificationDetailsResponseListNil"

// NotificationDetailsResponseList is a list of values of the 'notification_details_response' type.
type NotificationDetailsResponseList struct {
	href  string
	link  bool
	items []*NotificationDetailsResponse
}

// Kind returns the name of the type of the object.
func (l *NotificationDetailsResponseList) Kind() string {
	if l == nil {
		return NotificationDetailsResponseListNilKind
	}
	if l.link {
		return NotificationDetailsResponseListLinkKind
	}
	return NotificationDetailsResponseListKind
}

// Link returns true iif this is a link.
func (l *NotificationDetailsResponseList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NotificationDetailsResponseList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NotificationDetailsResponseList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NotificationDetailsResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NotificationDetailsResponseList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NotificationDetailsResponseList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NotificationDetailsResponseList) SetItems(items []*NotificationDetailsResponse) {
	l.items = items
}

// Items returns the items of the list.
func (l *NotificationDetailsResponseList) Items() []*NotificationDetailsResponse {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NotificationDetailsResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NotificationDetailsResponseList) Get(i int) *NotificationDetailsResponse {
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
func (l *NotificationDetailsResponseList) Slice() []*NotificationDetailsResponse {
	var slice []*NotificationDetailsResponse
	if l == nil {
		slice = make([]*NotificationDetailsResponse, 0)
	} else {
		slice = make([]*NotificationDetailsResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NotificationDetailsResponseList) Each(f func(item *NotificationDetailsResponse) bool) {
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
func (l *NotificationDetailsResponseList) Range(f func(index int, item *NotificationDetailsResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
