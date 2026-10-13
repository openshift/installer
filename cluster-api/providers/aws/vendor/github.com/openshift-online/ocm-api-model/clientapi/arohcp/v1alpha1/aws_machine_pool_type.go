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

// AWSMachinePoolKind is the name of the type used to represent objects
// of type 'AWS_machine_pool'.
const AWSMachinePoolKind = "AWSMachinePool"

// AWSMachinePoolLinkKind is the name of the type used to represent links
// to objects of type 'AWS_machine_pool'.
const AWSMachinePoolLinkKind = "AWSMachinePoolLink"

// AWSMachinePoolNilKind is the name of the type used to nil references
// to objects of type 'AWS_machine_pool'.
const AWSMachinePoolNilKind = "AWSMachinePoolNil"

// AWSMachinePool represents the values of the 'AWS_machine_pool' type.
//
// Representation of aws machine pool specific parameters.
type AWSMachinePool struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	additionalSecurityGroupIds []string
	availabilityZoneTypes      map[string]string
	spotMarketOptions          *AWSSpotMarketOptions
	subnetOutposts             map[string]string
	tags                       map[string]string
}

// Kind returns the name of the type of the object.
func (o *AWSMachinePool) Kind() string {
	if o == nil {
		return AWSMachinePoolNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AWSMachinePoolLinkKind
	}
	return AWSMachinePoolKind
}

// Link returns true if this is a link.
func (o *AWSMachinePool) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AWSMachinePool) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSMachinePool) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSMachinePool) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSMachinePool) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSMachinePool) Empty() bool {
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

// AdditionalSecurityGroupIds returns the value of the 'additional_security_group_ids' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional AWS Security Groups to be added machine pool. Note that machine pools can only be worker node at the time.
func (o *AWSMachinePool) AdditionalSecurityGroupIds() []string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.additionalSecurityGroupIds
	}
	return nil
}

// GetAdditionalSecurityGroupIds returns the value of the 'additional_security_group_ids' attribute and
// a flag indicating if the attribute has a value.
//
// Additional AWS Security Groups to be added machine pool. Note that machine pools can only be worker node at the time.
func (o *AWSMachinePool) GetAdditionalSecurityGroupIds() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.additionalSecurityGroupIds
	}
	return
}

// AvailabilityZoneTypes returns the value of the 'availability_zone_types' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Associates nodepool availability zones with zone types (e.g. wavelength, local).
func (o *AWSMachinePool) AvailabilityZoneTypes() map[string]string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.availabilityZoneTypes
	}
	return nil
}

// GetAvailabilityZoneTypes returns the value of the 'availability_zone_types' attribute and
// a flag indicating if the attribute has a value.
//
// Associates nodepool availability zones with zone types (e.g. wavelength, local).
func (o *AWSMachinePool) GetAvailabilityZoneTypes() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.availabilityZoneTypes
	}
	return
}

// SpotMarketOptions returns the value of the 'spot_market_options' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Use spot instances on this machine pool to reduce cost.
func (o *AWSMachinePool) SpotMarketOptions() *AWSSpotMarketOptions {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.spotMarketOptions
	}
	return nil
}

// GetSpotMarketOptions returns the value of the 'spot_market_options' attribute and
// a flag indicating if the attribute has a value.
//
// Use spot instances on this machine pool to reduce cost.
func (o *AWSMachinePool) GetSpotMarketOptions() (value *AWSSpotMarketOptions, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.spotMarketOptions
	}
	return
}

// SubnetOutposts returns the value of the 'subnet_outposts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Associates nodepool subnets with AWS Outposts.
func (o *AWSMachinePool) SubnetOutposts() map[string]string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.subnetOutposts
	}
	return nil
}

// GetSubnetOutposts returns the value of the 'subnet_outposts' attribute and
// a flag indicating if the attribute has a value.
//
// Associates nodepool subnets with AWS Outposts.
func (o *AWSMachinePool) GetSubnetOutposts() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.subnetOutposts
	}
	return
}

// Tags returns the value of the 'tags' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional keys and values that the machine pool provisioner will add as AWS tags to all AWS resources it creates.
//
// AWS tags must conform to the following standards:
// - Each resource may have a maximum of 25 tags
// - Tags beginning with "aws:" are reserved for system use and may not be set
// - Tag keys may be between 1 and 128 characters in length
// - Tag values may be between 0 and 256 characters in length
// - Tags may only contain letters, numbers, spaces, and the following characters: [_ . : / = + - @]
func (o *AWSMachinePool) Tags() map[string]string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.tags
	}
	return nil
}

// GetTags returns the value of the 'tags' attribute and
// a flag indicating if the attribute has a value.
//
// Optional keys and values that the machine pool provisioner will add as AWS tags to all AWS resources it creates.
//
// AWS tags must conform to the following standards:
// - Each resource may have a maximum of 25 tags
// - Tags beginning with "aws:" are reserved for system use and may not be set
// - Tag keys may be between 1 and 128 characters in length
// - Tag values may be between 0 and 256 characters in length
// - Tags may only contain letters, numbers, spaces, and the following characters: [_ . : / = + - @]
func (o *AWSMachinePool) GetTags() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.tags
	}
	return
}

// AWSMachinePoolListKind is the name of the type used to represent list of objects of
// type 'AWS_machine_pool'.
const AWSMachinePoolListKind = "AWSMachinePoolList"

// AWSMachinePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_machine_pool'.
const AWSMachinePoolListLinkKind = "AWSMachinePoolListLink"

// AWSMachinePoolNilKind is the name of the type used to nil lists of objects of
// type 'AWS_machine_pool'.
const AWSMachinePoolListNilKind = "AWSMachinePoolListNil"

// AWSMachinePoolList is a list of values of the 'AWS_machine_pool' type.
type AWSMachinePoolList struct {
	href  string
	link  bool
	items []*AWSMachinePool
}

// Kind returns the name of the type of the object.
func (l *AWSMachinePoolList) Kind() string {
	if l == nil {
		return AWSMachinePoolListNilKind
	}
	if l.link {
		return AWSMachinePoolListLinkKind
	}
	return AWSMachinePoolListKind
}

// Link returns true iif this is a link.
func (l *AWSMachinePoolList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AWSMachinePoolList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AWSMachinePoolList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AWSMachinePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSMachinePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSMachinePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSMachinePoolList) SetItems(items []*AWSMachinePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSMachinePoolList) Items() []*AWSMachinePool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSMachinePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSMachinePoolList) Get(i int) *AWSMachinePool {
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
func (l *AWSMachinePoolList) Slice() []*AWSMachinePool {
	var slice []*AWSMachinePool
	if l == nil {
		slice = make([]*AWSMachinePool, 0)
	} else {
		slice = make([]*AWSMachinePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSMachinePoolList) Each(f func(item *AWSMachinePool) bool) {
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
func (l *AWSMachinePoolList) Range(f func(index int, item *AWSMachinePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
