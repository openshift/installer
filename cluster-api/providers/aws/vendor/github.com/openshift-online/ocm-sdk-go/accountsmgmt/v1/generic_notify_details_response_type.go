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

// GenericNotifyDetailsResponseKind is the name of the type used to represent objects
// of type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseKind = "GenericNotifyDetailsResponse"

// GenericNotifyDetailsResponseLinkKind is the name of the type used to represent links
// to objects of type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseLinkKind = "GenericNotifyDetailsResponseLink"

// GenericNotifyDetailsResponseNilKind is the name of the type used to nil references
// to objects of type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseNilKind = "GenericNotifyDetailsResponseNil"

// GenericNotifyDetailsResponse represents the values of the 'generic_notify_details_response' type.
//
// class that defines notify details response in general.
type GenericNotifyDetailsResponse struct {
	bitmap_    uint32
	id         string
	href       string
	associates []string
	items      []*NotificationDetailsResponse
	recipients []string
}

// Kind returns the name of the type of the object.
func (o *GenericNotifyDetailsResponse) Kind() string {
	if o == nil {
		return GenericNotifyDetailsResponseNilKind
	}
	if o.bitmap_&1 != 0 {
		return GenericNotifyDetailsResponseLinkKind
	}
	return GenericNotifyDetailsResponseKind
}

// Link returns true if this is a link.
func (o *GenericNotifyDetailsResponse) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *GenericNotifyDetailsResponse) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *GenericNotifyDetailsResponse) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *GenericNotifyDetailsResponse) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *GenericNotifyDetailsResponse) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GenericNotifyDetailsResponse) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Associates returns the value of the 'associates' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Retrieved List of associates email address.
func (o *GenericNotifyDetailsResponse) Associates() []string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.associates
	}
	return nil
}

// GetAssociates returns the value of the 'associates' attribute and
// a flag indicating if the attribute has a value.
//
// Retrieved List of associates email address.
func (o *GenericNotifyDetailsResponse) GetAssociates() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.associates
	}
	return
}

// Items returns the value of the 'items' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Retrieved list of additional notify details parameters key-value.
func (o *GenericNotifyDetailsResponse) Items() []*NotificationDetailsResponse {
	if o != nil && o.bitmap_&16 != 0 {
		return o.items
	}
	return nil
}

// GetItems returns the value of the 'items' attribute and
// a flag indicating if the attribute has a value.
//
// Retrieved list of additional notify details parameters key-value.
func (o *GenericNotifyDetailsResponse) GetItems() (value []*NotificationDetailsResponse, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.items
	}
	return
}

// Recipients returns the value of the 'recipients' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Retrieved List of recipients username.
func (o *GenericNotifyDetailsResponse) Recipients() []string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.recipients
	}
	return nil
}

// GetRecipients returns the value of the 'recipients' attribute and
// a flag indicating if the attribute has a value.
//
// Retrieved List of recipients username.
func (o *GenericNotifyDetailsResponse) GetRecipients() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.recipients
	}
	return
}

// GenericNotifyDetailsResponseListKind is the name of the type used to represent list of objects of
// type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseListKind = "GenericNotifyDetailsResponseList"

// GenericNotifyDetailsResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseListLinkKind = "GenericNotifyDetailsResponseListLink"

// GenericNotifyDetailsResponseNilKind is the name of the type used to nil lists of objects of
// type 'generic_notify_details_response'.
const GenericNotifyDetailsResponseListNilKind = "GenericNotifyDetailsResponseListNil"

// GenericNotifyDetailsResponseList is a list of values of the 'generic_notify_details_response' type.
type GenericNotifyDetailsResponseList struct {
	href  string
	link  bool
	items []*GenericNotifyDetailsResponse
}

// Kind returns the name of the type of the object.
func (l *GenericNotifyDetailsResponseList) Kind() string {
	if l == nil {
		return GenericNotifyDetailsResponseListNilKind
	}
	if l.link {
		return GenericNotifyDetailsResponseListLinkKind
	}
	return GenericNotifyDetailsResponseListKind
}

// Link returns true iif this is a link.
func (l *GenericNotifyDetailsResponseList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *GenericNotifyDetailsResponseList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *GenericNotifyDetailsResponseList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *GenericNotifyDetailsResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GenericNotifyDetailsResponseList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GenericNotifyDetailsResponseList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GenericNotifyDetailsResponseList) SetItems(items []*GenericNotifyDetailsResponse) {
	l.items = items
}

// Items returns the items of the list.
func (l *GenericNotifyDetailsResponseList) Items() []*GenericNotifyDetailsResponse {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GenericNotifyDetailsResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GenericNotifyDetailsResponseList) Get(i int) *GenericNotifyDetailsResponse {
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
func (l *GenericNotifyDetailsResponseList) Slice() []*GenericNotifyDetailsResponse {
	var slice []*GenericNotifyDetailsResponse
	if l == nil {
		slice = make([]*GenericNotifyDetailsResponse, 0)
	} else {
		slice = make([]*GenericNotifyDetailsResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GenericNotifyDetailsResponseList) Each(f func(item *GenericNotifyDetailsResponse) bool) {
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
func (l *GenericNotifyDetailsResponseList) Range(f func(index int, item *GenericNotifyDetailsResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
