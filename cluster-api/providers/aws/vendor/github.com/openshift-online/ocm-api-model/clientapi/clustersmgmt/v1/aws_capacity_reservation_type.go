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

// AWSCapacityReservation represents the values of the 'AWS_capacity_reservation' type.
//
// AWS Capacity Reservation specification.
type AWSCapacityReservation struct {
	fieldSet_  []bool
	id         string
	marketType MarketType
	preference CapacityReservationPreference
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSCapacityReservation) Empty() bool {
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

// Id returns the value of the 'id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specify the target Capacity Reservation in which the EC2 instances will be launched.
func (o *AWSCapacityReservation) Id() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetId returns the value of the 'id' attribute and
// a flag indicating if the attribute has a value.
//
// Specify the target Capacity Reservation in which the EC2 instances will be launched.
func (o *AWSCapacityReservation) GetId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// MarketType returns the value of the 'market_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// marketType specifies the market type of the CapacityReservation for the EC2
// instances. Valid values are OnDemand, CapacityBlocks.
// "OnDemand": EC2 instances run as standard On-Demand instances.
// "CapacityBlocks": scheduled pre-purchased compute capacity.
func (o *AWSCapacityReservation) MarketType() MarketType {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.marketType
	}
	return MarketType("")
}

// GetMarketType returns the value of the 'market_type' attribute and
// a flag indicating if the attribute has a value.
//
// marketType specifies the market type of the CapacityReservation for the EC2
// instances. Valid values are OnDemand, CapacityBlocks.
// "OnDemand": EC2 instances run as standard On-Demand instances.
// "CapacityBlocks": scheduled pre-purchased compute capacity.
func (o *AWSCapacityReservation) GetMarketType() (value MarketType, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.marketType
	}
	return
}

// Preference returns the value of the 'preference' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// preference specifies how capacity reservations should be used with this NodePool
// "none": EC2 instances in this NodePool should never make use of capacity reservations. Note that this value cannot
// be specified if a capacity reservation Id is also specified
// "capacity-reservations-only": EC2 instances in this NodePool should only run in a capacity reservation
// "open": EC2 instances in this NodePool should run in an Open capacity reservation if available, otherwise run on demand.
// Note that this value cannot be specified if a capacity reservation Id is also specified
func (o *AWSCapacityReservation) Preference() CapacityReservationPreference {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.preference
	}
	return CapacityReservationPreference("")
}

// GetPreference returns the value of the 'preference' attribute and
// a flag indicating if the attribute has a value.
//
// preference specifies how capacity reservations should be used with this NodePool
// "none": EC2 instances in this NodePool should never make use of capacity reservations. Note that this value cannot
// be specified if a capacity reservation Id is also specified
// "capacity-reservations-only": EC2 instances in this NodePool should only run in a capacity reservation
// "open": EC2 instances in this NodePool should run in an Open capacity reservation if available, otherwise run on demand.
// Note that this value cannot be specified if a capacity reservation Id is also specified
func (o *AWSCapacityReservation) GetPreference() (value CapacityReservationPreference, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.preference
	}
	return
}

// AWSCapacityReservationListKind is the name of the type used to represent list of objects of
// type 'AWS_capacity_reservation'.
const AWSCapacityReservationListKind = "AWSCapacityReservationList"

// AWSCapacityReservationListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_capacity_reservation'.
const AWSCapacityReservationListLinkKind = "AWSCapacityReservationListLink"

// AWSCapacityReservationNilKind is the name of the type used to nil lists of objects of
// type 'AWS_capacity_reservation'.
const AWSCapacityReservationListNilKind = "AWSCapacityReservationListNil"

// AWSCapacityReservationList is a list of values of the 'AWS_capacity_reservation' type.
type AWSCapacityReservationList struct {
	href  string
	link  bool
	items []*AWSCapacityReservation
}

// Len returns the length of the list.
func (l *AWSCapacityReservationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSCapacityReservationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSCapacityReservationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSCapacityReservationList) SetItems(items []*AWSCapacityReservation) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSCapacityReservationList) Items() []*AWSCapacityReservation {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSCapacityReservationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSCapacityReservationList) Get(i int) *AWSCapacityReservation {
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
func (l *AWSCapacityReservationList) Slice() []*AWSCapacityReservation {
	var slice []*AWSCapacityReservation
	if l == nil {
		slice = make([]*AWSCapacityReservation, 0)
	} else {
		slice = make([]*AWSCapacityReservation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSCapacityReservationList) Each(f func(item *AWSCapacityReservation) bool) {
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
func (l *AWSCapacityReservationList) Range(f func(index int, item *AWSCapacityReservation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
