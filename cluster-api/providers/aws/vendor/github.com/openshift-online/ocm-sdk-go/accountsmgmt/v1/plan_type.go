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

// PlanKind is the name of the type used to represent objects
// of type 'plan'.
const PlanKind = "Plan"

// PlanLinkKind is the name of the type used to represent links
// to objects of type 'plan'.
const PlanLinkKind = "PlanLink"

// PlanNilKind is the name of the type used to nil references
// to objects of type 'plan'.
const PlanNilKind = "PlanNil"

// Plan represents the values of the 'plan' type.
type Plan struct {
	bitmap_  uint32
	id       string
	href     string
	category string
	name     string
	type_    string
}

// Kind returns the name of the type of the object.
func (o *Plan) Kind() string {
	if o == nil {
		return PlanNilKind
	}
	if o.bitmap_&1 != 0 {
		return PlanLinkKind
	}
	return PlanKind
}

// Link returns true if this is a link.
func (o *Plan) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Plan) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Plan) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Plan) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Plan) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Plan) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Category returns the value of the 'category' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Plan) Category() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.category
	}
	return ""
}

// GetCategory returns the value of the 'category' attribute and
// a flag indicating if the attribute has a value.
func (o *Plan) GetCategory() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.category
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Plan) Name() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Plan) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.name
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Plan) Type() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
func (o *Plan) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.type_
	}
	return
}

// PlanListKind is the name of the type used to represent list of objects of
// type 'plan'.
const PlanListKind = "PlanList"

// PlanListLinkKind is the name of the type used to represent links to list
// of objects of type 'plan'.
const PlanListLinkKind = "PlanListLink"

// PlanNilKind is the name of the type used to nil lists of objects of
// type 'plan'.
const PlanListNilKind = "PlanListNil"

// PlanList is a list of values of the 'plan' type.
type PlanList struct {
	href  string
	link  bool
	items []*Plan
}

// Kind returns the name of the type of the object.
func (l *PlanList) Kind() string {
	if l == nil {
		return PlanListNilKind
	}
	if l.link {
		return PlanListLinkKind
	}
	return PlanListKind
}

// Link returns true iif this is a link.
func (l *PlanList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *PlanList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *PlanList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *PlanList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *PlanList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *PlanList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *PlanList) SetItems(items []*Plan) {
	l.items = items
}

// Items returns the items of the list.
func (l *PlanList) Items() []*Plan {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *PlanList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PlanList) Get(i int) *Plan {
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
func (l *PlanList) Slice() []*Plan {
	var slice []*Plan
	if l == nil {
		slice = make([]*Plan, 0)
	} else {
		slice = make([]*Plan, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PlanList) Each(f func(item *Plan) bool) {
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
func (l *PlanList) Range(f func(index int, item *Plan) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
