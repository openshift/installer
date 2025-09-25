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

// DecisionKind is the name of the type used to represent objects
// of type 'decision'.
const DecisionKind = "Decision"

// DecisionLinkKind is the name of the type used to represent links
// to objects of type 'decision'.
const DecisionLinkKind = "DecisionLink"

// DecisionNilKind is the name of the type used to nil references
// to objects of type 'decision'.
const DecisionNilKind = "DecisionNil"

// Decision represents the values of the 'decision' type.
//
// Representation of an decision.
type Decision struct {
	bitmap_       uint32
	id            string
	href          string
	createdAt     time.Time
	decidedBy     string
	decision      DecisionDecision
	justification string
	updatedAt     time.Time
}

// Kind returns the name of the type of the object.
func (o *Decision) Kind() string {
	if o == nil {
		return DecisionNilKind
	}
	if o.bitmap_&1 != 0 {
		return DecisionLinkKind
	}
	return DecisionKind
}

// Link returns true if this is a link.
func (o *Decision) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Decision) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Decision) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Decision) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Decision) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Decision) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the decision was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Decision) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the decision was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Decision) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// DecidedBy returns the value of the 'decided_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// User that decided.
func (o *Decision) DecidedBy() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.decidedBy
	}
	return ""
}

// GetDecidedBy returns the value of the 'decided_by' attribute and
// a flag indicating if the attribute has a value.
//
// User that decided.
func (o *Decision) GetDecidedBy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.decidedBy
	}
	return
}

// Decision returns the value of the 'decision' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the decision.
func (o *Decision) Decision() DecisionDecision {
	if o != nil && o.bitmap_&32 != 0 {
		return o.decision
	}
	return DecisionDecision("")
}

// GetDecision returns the value of the 'decision' attribute and
// a flag indicating if the attribute has a value.
//
// State of the decision.
func (o *Decision) GetDecision() (value DecisionDecision, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.decision
	}
	return
}

// Justification returns the value of the 'justification' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Justification of the decision.
func (o *Decision) Justification() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.justification
	}
	return ""
}

// GetJustification returns the value of the 'justification' attribute and
// a flag indicating if the attribute has a value.
//
// Justification of the decision.
func (o *Decision) GetJustification() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.justification
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the decision was lastly updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Decision) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&128 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the decision was lastly updated, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Decision) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// DecisionListKind is the name of the type used to represent list of objects of
// type 'decision'.
const DecisionListKind = "DecisionList"

// DecisionListLinkKind is the name of the type used to represent links to list
// of objects of type 'decision'.
const DecisionListLinkKind = "DecisionListLink"

// DecisionNilKind is the name of the type used to nil lists of objects of
// type 'decision'.
const DecisionListNilKind = "DecisionListNil"

// DecisionList is a list of values of the 'decision' type.
type DecisionList struct {
	href  string
	link  bool
	items []*Decision
}

// Kind returns the name of the type of the object.
func (l *DecisionList) Kind() string {
	if l == nil {
		return DecisionListNilKind
	}
	if l.link {
		return DecisionListLinkKind
	}
	return DecisionListKind
}

// Link returns true iif this is a link.
func (l *DecisionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *DecisionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *DecisionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *DecisionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *DecisionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *DecisionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *DecisionList) SetItems(items []*Decision) {
	l.items = items
}

// Items returns the items of the list.
func (l *DecisionList) Items() []*Decision {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *DecisionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DecisionList) Get(i int) *Decision {
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
func (l *DecisionList) Slice() []*Decision {
	var slice []*Decision
	if l == nil {
		slice = make([]*Decision, 0)
	} else {
		slice = make([]*Decision, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DecisionList) Each(f func(item *Decision) bool) {
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
func (l *DecisionList) Range(f func(index int, item *Decision) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
