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

// CloudVPC represents the values of the 'cloud_VPC' type.
//
// Description of a cloud provider virtual private cloud.
type CloudVPC struct {
	fieldSet_         []bool
	awsSecurityGroups []*SecurityGroup
	awsSubnets        []*Subnetwork
	cidrBlock         string
	id                string
	name              string
	subnets           []string
	redHatManaged     bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudVPC) Empty() bool {
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

// AWSSecurityGroups returns the value of the 'AWS_security_groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of AWS security groups with details.
func (o *CloudVPC) AWSSecurityGroups() []*SecurityGroup {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.awsSecurityGroups
	}
	return nil
}

// GetAWSSecurityGroups returns the value of the 'AWS_security_groups' attribute and
// a flag indicating if the attribute has a value.
//
// List of AWS security groups with details.
func (o *CloudVPC) GetAWSSecurityGroups() (value []*SecurityGroup, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.awsSecurityGroups
	}
	return
}

// AWSSubnets returns the value of the 'AWS_subnets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of AWS subnetworks with details.
func (o *CloudVPC) AWSSubnets() []*Subnetwork {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.awsSubnets
	}
	return nil
}

// GetAWSSubnets returns the value of the 'AWS_subnets' attribute and
// a flag indicating if the attribute has a value.
//
// List of AWS subnetworks with details.
func (o *CloudVPC) GetAWSSubnets() (value []*Subnetwork, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.awsSubnets
	}
	return
}

// CIDRBlock returns the value of the 'CIDR_block' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CIDR block of the virtual private cloud.
func (o *CloudVPC) CIDRBlock() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.cidrBlock
	}
	return ""
}

// GetCIDRBlock returns the value of the 'CIDR_block' attribute and
// a flag indicating if the attribute has a value.
//
// CIDR block of the virtual private cloud.
func (o *CloudVPC) GetCIDRBlock() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.cidrBlock
	}
	return
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of virtual private cloud.
func (o *CloudVPC) ID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of virtual private cloud.
func (o *CloudVPC) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.id
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of virtual private cloud according to its `Name` tag on AWS.
func (o *CloudVPC) Name() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of virtual private cloud according to its `Name` tag on AWS.
func (o *CloudVPC) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.name
	}
	return
}

// RedHatManaged returns the value of the 'red_hat_managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// If the resource is RH managed.
func (o *CloudVPC) RedHatManaged() bool {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.redHatManaged
	}
	return false
}

// GetRedHatManaged returns the value of the 'red_hat_managed' attribute and
// a flag indicating if the attribute has a value.
//
// If the resource is RH managed.
func (o *CloudVPC) GetRedHatManaged() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.redHatManaged
	}
	return
}

// Subnets returns the value of the 'subnets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of subnets used by the virtual private cloud.
func (o *CloudVPC) Subnets() []string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.subnets
	}
	return nil
}

// GetSubnets returns the value of the 'subnets' attribute and
// a flag indicating if the attribute has a value.
//
// List of subnets used by the virtual private cloud.
func (o *CloudVPC) GetSubnets() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.subnets
	}
	return
}

// CloudVPCListKind is the name of the type used to represent list of objects of
// type 'cloud_VPC'.
const CloudVPCListKind = "CloudVPCList"

// CloudVPCListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_VPC'.
const CloudVPCListLinkKind = "CloudVPCListLink"

// CloudVPCNilKind is the name of the type used to nil lists of objects of
// type 'cloud_VPC'.
const CloudVPCListNilKind = "CloudVPCListNil"

// CloudVPCList is a list of values of the 'cloud_VPC' type.
type CloudVPCList struct {
	href  string
	link  bool
	items []*CloudVPC
}

// Len returns the length of the list.
func (l *CloudVPCList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CloudVPCList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CloudVPCList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CloudVPCList) SetItems(items []*CloudVPC) {
	l.items = items
}

// Items returns the items of the list.
func (l *CloudVPCList) Items() []*CloudVPC {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CloudVPCList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudVPCList) Get(i int) *CloudVPC {
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
func (l *CloudVPCList) Slice() []*CloudVPC {
	var slice []*CloudVPC
	if l == nil {
		slice = make([]*CloudVPC, 0)
	} else {
		slice = make([]*CloudVPC, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudVPCList) Each(f func(item *CloudVPC) bool) {
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
func (l *CloudVPCList) Range(f func(index int, item *CloudVPC) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
