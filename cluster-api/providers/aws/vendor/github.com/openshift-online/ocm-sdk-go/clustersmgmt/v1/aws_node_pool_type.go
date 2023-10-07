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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

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
	bitmap_         uint32
	id              string
	href            string
	instanceProfile string
	instanceType    string
	tags            map[string]string
}

// Kind returns the name of the type of the object.
func (o *AWSNodePool) Kind() string {
	if o == nil {
		return AWSNodePoolNilKind
	}
	if o.bitmap_&1 != 0 {
		return AWSNodePoolLinkKind
	}
	return AWSNodePoolKind
}

// Link returns true iif this is a link.
func (o *AWSNodePool) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AWSNodePool) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSNodePool) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSNodePool) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSNodePool) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSNodePool) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// InstanceProfile returns the value of the 'instance_profile' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// InstanceProfile is the AWS EC2 instance profile, which is a container for an IAM role that the EC2 instance uses.
func (o *AWSNodePool) InstanceProfile() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.instanceProfile
	}
	return ""
}

// GetInstanceProfile returns the value of the 'instance_profile' attribute and
// a flag indicating if the attribute has a value.
//
// InstanceProfile is the AWS EC2 instance profile, which is a container for an IAM role that the EC2 instance uses.
func (o *AWSNodePool) GetInstanceProfile() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
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
	if o != nil && o.bitmap_&16 != 0 {
		return o.instanceType
	}
	return ""
}

// GetInstanceType returns the value of the 'instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// InstanceType is an ec2 instance type for node instances (e.g. m5.large).
func (o *AWSNodePool) GetInstanceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.instanceType
	}
	return
}

// Tags returns the value of the 'tags' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWSNodePool) Tags() map[string]string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.tags
	}
	return nil
}

// GetTags returns the value of the 'tags' attribute and
// a flag indicating if the attribute has a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWSNodePool) GetTags() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
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
