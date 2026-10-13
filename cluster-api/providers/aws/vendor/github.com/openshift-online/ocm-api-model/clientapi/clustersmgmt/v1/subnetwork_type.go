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

// Subnetwork represents the values of the 'subnetwork' type.
//
// AWS subnetwork object to be used while installing a cluster
type Subnetwork struct {
	fieldSet_        []bool
	cidrBlock        string
	availabilityZone string
	name             string
	subnetID         string
	public           bool
	redHatManaged    bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Subnetwork) Empty() bool {
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

// CIDRBlock returns the value of the 'CIDR_block' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The CIDR Block of the subnet.
func (o *Subnetwork) CIDRBlock() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.cidrBlock
	}
	return ""
}

// GetCIDRBlock returns the value of the 'CIDR_block' attribute and
// a flag indicating if the attribute has a value.
//
// The CIDR Block of the subnet.
func (o *Subnetwork) GetCIDRBlock() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.cidrBlock
	}
	return
}

// AvailabilityZone returns the value of the 'availability_zone' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The availability zone to which the subnet is related.
func (o *Subnetwork) AvailabilityZone() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.availabilityZone
	}
	return ""
}

// GetAvailabilityZone returns the value of the 'availability_zone' attribute and
// a flag indicating if the attribute has a value.
//
// The availability zone to which the subnet is related.
func (o *Subnetwork) GetAvailabilityZone() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.availabilityZone
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the subnet according to its `Name` tag on AWS.
func (o *Subnetwork) Name() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the subnet according to its `Name` tag on AWS.
func (o *Subnetwork) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.name
	}
	return
}

// Public returns the value of the 'public' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Whether or not it is a public subnet.
func (o *Subnetwork) Public() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.public
	}
	return false
}

// GetPublic returns the value of the 'public' attribute and
// a flag indicating if the attribute has a value.
//
// Whether or not it is a public subnet.
func (o *Subnetwork) GetPublic() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.public
	}
	return
}

// RedHatManaged returns the value of the 'red_hat_managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// If the resource is RH managed.
func (o *Subnetwork) RedHatManaged() bool {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.redHatManaged
	}
	return false
}

// GetRedHatManaged returns the value of the 'red_hat_managed' attribute and
// a flag indicating if the attribute has a value.
//
// If the resource is RH managed.
func (o *Subnetwork) GetRedHatManaged() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.redHatManaged
	}
	return
}

// SubnetID returns the value of the 'subnet_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The subnet ID to be used while installing a cluster.
func (o *Subnetwork) SubnetID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.subnetID
	}
	return ""
}

// GetSubnetID returns the value of the 'subnet_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The subnet ID to be used while installing a cluster.
func (o *Subnetwork) GetSubnetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.subnetID
	}
	return
}

// SubnetworkListKind is the name of the type used to represent list of objects of
// type 'subnetwork'.
const SubnetworkListKind = "SubnetworkList"

// SubnetworkListLinkKind is the name of the type used to represent links to list
// of objects of type 'subnetwork'.
const SubnetworkListLinkKind = "SubnetworkListLink"

// SubnetworkNilKind is the name of the type used to nil lists of objects of
// type 'subnetwork'.
const SubnetworkListNilKind = "SubnetworkListNil"

// SubnetworkList is a list of values of the 'subnetwork' type.
type SubnetworkList struct {
	href  string
	link  bool
	items []*Subnetwork
}

// Len returns the length of the list.
func (l *SubnetworkList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SubnetworkList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SubnetworkList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SubnetworkList) SetItems(items []*Subnetwork) {
	l.items = items
}

// Items returns the items of the list.
func (l *SubnetworkList) Items() []*Subnetwork {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SubnetworkList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SubnetworkList) Get(i int) *Subnetwork {
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
func (l *SubnetworkList) Slice() []*Subnetwork {
	var slice []*Subnetwork
	if l == nil {
		slice = make([]*Subnetwork, 0)
	} else {
		slice = make([]*Subnetwork, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SubnetworkList) Each(f func(item *Subnetwork) bool) {
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
func (l *SubnetworkList) Range(f func(index int, item *Subnetwork) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
