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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// SubscriptionKind is the name of the type used to represent objects
// of type 'subscription'.
const SubscriptionKind = "Subscription"

// SubscriptionLinkKind is the name of the type used to represent links
// to objects of type 'subscription'.
const SubscriptionLinkKind = "SubscriptionLink"

// SubscriptionNilKind is the name of the type used to nil references
// to objects of type 'subscription'.
const SubscriptionNilKind = "SubscriptionNil"

// Subscription represents the values of the 'subscription' type.
//
// Definition of a subscription.
type Subscription struct {
	fieldSet_ []bool
	id        string
	href      string
}

// Kind returns the name of the type of the object.
func (o *Subscription) Kind() string {
	if o == nil {
		return SubscriptionNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return SubscriptionLinkKind
	}
	return SubscriptionKind
}

// Link returns true if this is a link.
func (o *Subscription) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Subscription) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Subscription) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Subscription) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Subscription) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Subscription) Empty() bool {
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

// SubscriptionListKind is the name of the type used to represent list of objects of
// type 'subscription'.
const SubscriptionListKind = "SubscriptionList"

// SubscriptionListLinkKind is the name of the type used to represent links to list
// of objects of type 'subscription'.
const SubscriptionListLinkKind = "SubscriptionListLink"

// SubscriptionNilKind is the name of the type used to nil lists of objects of
// type 'subscription'.
const SubscriptionListNilKind = "SubscriptionListNil"

// SubscriptionList is a list of values of the 'subscription' type.
type SubscriptionList struct {
	href  string
	link  bool
	items []*Subscription
}

// Kind returns the name of the type of the object.
func (l *SubscriptionList) Kind() string {
	if l == nil {
		return SubscriptionListNilKind
	}
	if l.link {
		return SubscriptionListLinkKind
	}
	return SubscriptionListKind
}

// Link returns true iif this is a link.
func (l *SubscriptionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *SubscriptionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *SubscriptionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *SubscriptionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SubscriptionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SubscriptionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SubscriptionList) SetItems(items []*Subscription) {
	l.items = items
}

// Items returns the items of the list.
func (l *SubscriptionList) Items() []*Subscription {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SubscriptionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SubscriptionList) Get(i int) *Subscription {
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
func (l *SubscriptionList) Slice() []*Subscription {
	var slice []*Subscription
	if l == nil {
		slice = make([]*Subscription, 0)
	} else {
		slice = make([]*Subscription, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SubscriptionList) Each(f func(item *Subscription) bool) {
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
func (l *SubscriptionList) Range(f func(index int, item *Subscription) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
