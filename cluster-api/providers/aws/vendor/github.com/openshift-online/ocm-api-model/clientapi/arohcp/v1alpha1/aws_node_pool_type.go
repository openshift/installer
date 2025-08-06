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

// AWSNodePoolKind is the name of the type used to represent objects
// of type 'AWS_node_pool'.
const AWSNodePoolKind = "AWSNodePool"

// AWSNodePoolLinkKind is the name of the type used to represent links
// to objects of type 'AWS_node_pool'.
const AWSNodePoolLinkKind = "AWSNodePoolLink"

// AWSNodePoolNilKind is the name of the type used to nil references
// to objects of type 'AWS_node_pool'.
const AWSNodePoolNilKind = "AWSNodePoolNil"

// AWSNodePool represents the values of the 'AWS_node_pool' type.
//
// Representation of aws node pool specific parameters.
type AWSNodePool struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	additionalSecurityGroupIds []string
	availabilityZoneTypes      map[string]string
	capacityReservation        *AWSCapacityReservation
	ec2MetadataHttpTokens      Ec2MetadataHttpTokens
	instanceProfile            string
	instanceType               string
	rootVolume                 *AWSVolume
	subnetOutposts             map[string]string
	tags                       map[string]string
}

// Kind returns the name of the type of the object.
func (o *AWSNodePool) Kind() string {
	if o == nil {
		return AWSNodePoolNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AWSNodePoolLinkKind
	}
	return AWSNodePoolKind
}

// Link returns true if this is a link.
func (o *AWSNodePool) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AWSNodePool) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSNodePool) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSNodePool) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSNodePool) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSNodePool) Empty() bool {
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
// Additional AWS Security Groups to be added node pool.
func (o *AWSNodePool) AdditionalSecurityGroupIds() []string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.additionalSecurityGroupIds
	}
	return nil
}

// GetAdditionalSecurityGroupIds returns the value of the 'additional_security_group_ids' attribute and
// a flag indicating if the attribute has a value.
//
// Additional AWS Security Groups to be added node pool.
func (o *AWSNodePool) GetAdditionalSecurityGroupIds() (value []string, ok bool) {
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
func (o *AWSNodePool) AvailabilityZoneTypes() map[string]string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.availabilityZoneTypes
	}
	return nil
}

// GetAvailabilityZoneTypes returns the value of the 'availability_zone_types' attribute and
// a flag indicating if the attribute has a value.
//
// Associates nodepool availability zones with zone types (e.g. wavelength, local).
func (o *AWSNodePool) GetAvailabilityZoneTypes() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.availabilityZoneTypes
	}
	return
}

// CapacityReservation returns the value of the 'capacity_reservation' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// If present it defines the AWS Capacity Reservation used for this NodePool
func (o *AWSNodePool) CapacityReservation() *AWSCapacityReservation {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.capacityReservation
	}
	return nil
}

// GetCapacityReservation returns the value of the 'capacity_reservation' attribute and
// a flag indicating if the attribute has a value.
//
// If present it defines the AWS Capacity Reservation used for this NodePool
func (o *AWSNodePool) GetCapacityReservation() (value *AWSCapacityReservation, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.capacityReservation
	}
	return
}

// Ec2MetadataHttpTokens returns the value of the 'ec_2_metadata_http_tokens' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (o *AWSNodePool) Ec2MetadataHttpTokens() Ec2MetadataHttpTokens {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.ec2MetadataHttpTokens
	}
	return Ec2MetadataHttpTokens("")
}

// GetEc2MetadataHttpTokens returns the value of the 'ec_2_metadata_http_tokens' attribute and
// a flag indicating if the attribute has a value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (o *AWSNodePool) GetEc2MetadataHttpTokens() (value Ec2MetadataHttpTokens, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.ec2MetadataHttpTokens
	}
	return
}

// InstanceProfile returns the value of the 'instance_profile' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// InstanceProfile is the AWS EC2 instance profile, which is a container for an IAM role that the EC2 instance uses.
func (o *AWSNodePool) InstanceProfile() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.instanceProfile
	}
	return ""
}

// GetInstanceProfile returns the value of the 'instance_profile' attribute and
// a flag indicating if the attribute has a value.
//
// InstanceProfile is the AWS EC2 instance profile, which is a container for an IAM role that the EC2 instance uses.
func (o *AWSNodePool) GetInstanceProfile() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.instanceProfile
	}
	return
}

// InstanceType returns the value of the 'instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// InstanceType is an ec2 instance type for node instances (e.g. m5.large).
func (o *AWSNodePool) InstanceType() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.instanceType
	}
	return ""
}

// GetInstanceType returns the value of the 'instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// InstanceType is an ec2 instance type for node instances (e.g. m5.large).
func (o *AWSNodePool) GetInstanceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.instanceType
	}
	return
}

// RootVolume returns the value of the 'root_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS Volume specification to be used to set custom worker disk size
func (o *AWSNodePool) RootVolume() *AWSVolume {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.rootVolume
	}
	return nil
}

// GetRootVolume returns the value of the 'root_volume' attribute and
// a flag indicating if the attribute has a value.
//
// AWS Volume specification to be used to set custom worker disk size
func (o *AWSNodePool) GetRootVolume() (value *AWSVolume, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.rootVolume
	}
	return
}

// SubnetOutposts returns the value of the 'subnet_outposts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Associates nodepool subnets with AWS Outposts.
func (o *AWSNodePool) SubnetOutposts() map[string]string {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.subnetOutposts
	}
	return nil
}

// GetSubnetOutposts returns the value of the 'subnet_outposts' attribute and
// a flag indicating if the attribute has a value.
//
// Associates nodepool subnets with AWS Outposts.
func (o *AWSNodePool) GetSubnetOutposts() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.subnetOutposts
	}
	return
}

// Tags returns the value of the 'tags' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates.
//
// AWS tags must conform to the following standards:
// - Each resource may have a maximum of 25 tags
// - Tags beginning with "aws:" are reserved for system use and may not be set
// - Tag keys may be between 1 and 128 characters in length
// - Tag values may be between 0 and 256 characters in length
// - Tags may only contain letters, numbers, spaces, and the following characters: [_ . : / = + - @]
func (o *AWSNodePool) Tags() map[string]string {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.tags
	}
	return nil
}

// GetTags returns the value of the 'tags' attribute and
// a flag indicating if the attribute has a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates.
//
// AWS tags must conform to the following standards:
// - Each resource may have a maximum of 25 tags
// - Tags beginning with "aws:" are reserved for system use and may not be set
// - Tag keys may be between 1 and 128 characters in length
// - Tag values may be between 0 and 256 characters in length
// - Tags may only contain letters, numbers, spaces, and the following characters: [_ . : / = + - @]
func (o *AWSNodePool) GetTags() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.tags
	}
	return
}

// AWSNodePoolListKind is the name of the type used to represent list of objects of
// type 'AWS_node_pool'.
const AWSNodePoolListKind = "AWSNodePoolList"

// AWSNodePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_node_pool'.
const AWSNodePoolListLinkKind = "AWSNodePoolListLink"

// AWSNodePoolNilKind is the name of the type used to nil lists of objects of
// type 'AWS_node_pool'.
const AWSNodePoolListNilKind = "AWSNodePoolListNil"

// AWSNodePoolList is a list of values of the 'AWS_node_pool' type.
type AWSNodePoolList struct {
	href  string
	link  bool
	items []*AWSNodePool
}

// Kind returns the name of the type of the object.
func (l *AWSNodePoolList) Kind() string {
	if l == nil {
		return AWSNodePoolListNilKind
	}
	if l.link {
		return AWSNodePoolListLinkKind
	}
	return AWSNodePoolListKind
}

// Link returns true iif this is a link.
func (l *AWSNodePoolList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AWSNodePoolList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AWSNodePoolList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AWSNodePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSNodePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSNodePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSNodePoolList) SetItems(items []*AWSNodePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSNodePoolList) Items() []*AWSNodePool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSNodePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSNodePoolList) Get(i int) *AWSNodePool {
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
func (l *AWSNodePoolList) Slice() []*AWSNodePool {
	var slice []*AWSNodePool
	if l == nil {
		slice = make([]*AWSNodePool, 0)
	} else {
		slice = make([]*AWSNodePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSNodePoolList) Each(f func(item *AWSNodePool) bool) {
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
func (l *AWSNodePoolList) Range(f func(index int, item *AWSNodePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
