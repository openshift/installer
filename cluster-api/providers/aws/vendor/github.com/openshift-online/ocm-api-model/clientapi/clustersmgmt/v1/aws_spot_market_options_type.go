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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// AWSSpotMarketOptionsKind is the name of the type used to represent objects
// of type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsKind = "AWSSpotMarketOptions"

// AWSSpotMarketOptionsLinkKind is the name of the type used to represent links
// to objects of type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsLinkKind = "AWSSpotMarketOptionsLink"

// AWSSpotMarketOptionsNilKind is the name of the type used to nil references
// to objects of type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsNilKind = "AWSSpotMarketOptionsNil"

// AWSSpotMarketOptions represents the values of the 'AWS_spot_market_options' type.
//
// Spot market options for AWS machine pool.
type AWSSpotMarketOptions struct {
	fieldSet_ []bool
	id        string
	href      string
	maxPrice  float64
}

// Kind returns the name of the type of the object.
func (o *AWSSpotMarketOptions) Kind() string {
	if o == nil {
		return AWSSpotMarketOptionsNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AWSSpotMarketOptionsLinkKind
	}
	return AWSSpotMarketOptionsKind
}

// Link returns true if this is a link.
func (o *AWSSpotMarketOptions) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AWSSpotMarketOptions) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSSpotMarketOptions) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSSpotMarketOptions) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSSpotMarketOptions) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSSpotMarketOptions) Empty() bool {
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

// MaxPrice returns the value of the 'max_price' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The max price for spot instance. Optional.
// If not set, use the on-demand price.
func (o *AWSSpotMarketOptions) MaxPrice() float64 {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.maxPrice
	}
	return 0.0
}

// GetMaxPrice returns the value of the 'max_price' attribute and
// a flag indicating if the attribute has a value.
//
// The max price for spot instance. Optional.
// If not set, use the on-demand price.
func (o *AWSSpotMarketOptions) GetMaxPrice() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.maxPrice
	}
	return
}

// AWSSpotMarketOptionsListKind is the name of the type used to represent list of objects of
// type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsListKind = "AWSSpotMarketOptionsList"

// AWSSpotMarketOptionsListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsListLinkKind = "AWSSpotMarketOptionsListLink"

// AWSSpotMarketOptionsNilKind is the name of the type used to nil lists of objects of
// type 'AWS_spot_market_options'.
const AWSSpotMarketOptionsListNilKind = "AWSSpotMarketOptionsListNil"

// AWSSpotMarketOptionsList is a list of values of the 'AWS_spot_market_options' type.
type AWSSpotMarketOptionsList struct {
	href  string
	link  bool
	items []*AWSSpotMarketOptions
}

// Kind returns the name of the type of the object.
func (l *AWSSpotMarketOptionsList) Kind() string {
	if l == nil {
		return AWSSpotMarketOptionsListNilKind
	}
	if l.link {
		return AWSSpotMarketOptionsListLinkKind
	}
	return AWSSpotMarketOptionsListKind
}

// Link returns true iif this is a link.
func (l *AWSSpotMarketOptionsList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AWSSpotMarketOptionsList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AWSSpotMarketOptionsList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AWSSpotMarketOptionsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSSpotMarketOptionsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSSpotMarketOptionsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSSpotMarketOptionsList) SetItems(items []*AWSSpotMarketOptions) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSSpotMarketOptionsList) Items() []*AWSSpotMarketOptions {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSSpotMarketOptionsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSSpotMarketOptionsList) Get(i int) *AWSSpotMarketOptions {
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
func (l *AWSSpotMarketOptionsList) Slice() []*AWSSpotMarketOptions {
	var slice []*AWSSpotMarketOptions
	if l == nil {
		slice = make([]*AWSSpotMarketOptions, 0)
	} else {
		slice = make([]*AWSSpotMarketOptions, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSSpotMarketOptionsList) Each(f func(item *AWSSpotMarketOptions) bool) {
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
func (l *AWSSpotMarketOptionsList) Range(f func(index int, item *AWSSpotMarketOptions) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
