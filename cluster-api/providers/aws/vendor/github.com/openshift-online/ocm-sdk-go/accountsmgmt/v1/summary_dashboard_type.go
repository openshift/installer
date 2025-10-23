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

// SummaryDashboardKind is the name of the type used to represent objects
// of type 'summary_dashboard'.
const SummaryDashboardKind = "SummaryDashboard"

// SummaryDashboardLinkKind is the name of the type used to represent links
// to objects of type 'summary_dashboard'.
const SummaryDashboardLinkKind = "SummaryDashboardLink"

// SummaryDashboardNilKind is the name of the type used to nil references
// to objects of type 'summary_dashboard'.
const SummaryDashboardNilKind = "SummaryDashboardNil"

// SummaryDashboard represents the values of the 'summary_dashboard' type.
type SummaryDashboard struct {
	bitmap_ uint32
	id      string
	href    string
	metrics []*SummaryMetrics
}

// Kind returns the name of the type of the object.
func (o *SummaryDashboard) Kind() string {
	if o == nil {
		return SummaryDashboardNilKind
	}
	if o.bitmap_&1 != 0 {
		return SummaryDashboardLinkKind
	}
	return SummaryDashboardKind
}

// Link returns true if this is a link.
func (o *SummaryDashboard) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *SummaryDashboard) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *SummaryDashboard) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *SummaryDashboard) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *SummaryDashboard) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SummaryDashboard) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Metrics returns the value of the 'metrics' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SummaryDashboard) Metrics() []*SummaryMetrics {
	if o != nil && o.bitmap_&8 != 0 {
		return o.metrics
	}
	return nil
}

// GetMetrics returns the value of the 'metrics' attribute and
// a flag indicating if the attribute has a value.
func (o *SummaryDashboard) GetMetrics() (value []*SummaryMetrics, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.metrics
	}
	return
}

// SummaryDashboardListKind is the name of the type used to represent list of objects of
// type 'summary_dashboard'.
const SummaryDashboardListKind = "SummaryDashboardList"

// SummaryDashboardListLinkKind is the name of the type used to represent links to list
// of objects of type 'summary_dashboard'.
const SummaryDashboardListLinkKind = "SummaryDashboardListLink"

// SummaryDashboardNilKind is the name of the type used to nil lists of objects of
// type 'summary_dashboard'.
const SummaryDashboardListNilKind = "SummaryDashboardListNil"

// SummaryDashboardList is a list of values of the 'summary_dashboard' type.
type SummaryDashboardList struct {
	href  string
	link  bool
	items []*SummaryDashboard
}

// Kind returns the name of the type of the object.
func (l *SummaryDashboardList) Kind() string {
	if l == nil {
		return SummaryDashboardListNilKind
	}
	if l.link {
		return SummaryDashboardListLinkKind
	}
	return SummaryDashboardListKind
}

// Link returns true iif this is a link.
func (l *SummaryDashboardList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *SummaryDashboardList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *SummaryDashboardList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *SummaryDashboardList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SummaryDashboardList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SummaryDashboardList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SummaryDashboardList) SetItems(items []*SummaryDashboard) {
	l.items = items
}

// Items returns the items of the list.
func (l *SummaryDashboardList) Items() []*SummaryDashboard {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SummaryDashboardList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SummaryDashboardList) Get(i int) *SummaryDashboard {
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
func (l *SummaryDashboardList) Slice() []*SummaryDashboard {
	var slice []*SummaryDashboard
	if l == nil {
		slice = make([]*SummaryDashboard, 0)
	} else {
		slice = make([]*SummaryDashboard, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SummaryDashboardList) Each(f func(item *SummaryDashboard) bool) {
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
func (l *SummaryDashboardList) Range(f func(index int, item *SummaryDashboard) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
