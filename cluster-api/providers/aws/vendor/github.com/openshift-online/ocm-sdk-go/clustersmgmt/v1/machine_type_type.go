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

// MachineTypeKind is the name of the type used to represent objects
// of type 'machine_type'.
const MachineTypeKind = "MachineType"

// MachineTypeLinkKind is the name of the type used to represent links
// to objects of type 'machine_type'.
const MachineTypeLinkKind = "MachineTypeLink"

// MachineTypeNilKind is the name of the type used to nil references
// to objects of type 'machine_type'.
const MachineTypeNilKind = "MachineTypeNil"

// MachineType represents the values of the 'machine_type' type.
//
// Machine type.
type MachineType struct {
	bitmap_       uint32
	id            string
	href          string
	cpu           *Value
	category      MachineTypeCategory
	cloudProvider *CloudProvider
	genericName   string
	memory        *Value
	name          string
	size          MachineTypeSize
	ccsOnly       bool
}

// Kind returns the name of the type of the object.
func (o *MachineType) Kind() string {
	if o == nil {
		return MachineTypeNilKind
	}
	if o.bitmap_&1 != 0 {
		return MachineTypeLinkKind
	}
	return MachineTypeKind
}

// Link returns true iif this is a link.
func (o *MachineType) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *MachineType) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *MachineType) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *MachineType) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *MachineType) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MachineType) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CCSOnly returns the value of the 'CCS_only' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// 'true' if the instance type is supported only for CCS clusters, 'false' otherwise.
func (o *MachineType) CCSOnly() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.ccsOnly
	}
	return false
}

// GetCCSOnly returns the value of the 'CCS_only' attribute and
// a flag indicating if the attribute has a value.
//
// 'true' if the instance type is supported only for CCS clusters, 'false' otherwise.
func (o *MachineType) GetCCSOnly() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.ccsOnly
	}
	return
}

// CPU returns the value of the 'CPU' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The amount of cpu's of the machine type.
func (o *MachineType) CPU() *Value {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cpu
	}
	return nil
}

// GetCPU returns the value of the 'CPU' attribute and
// a flag indicating if the attribute has a value.
//
// The amount of cpu's of the machine type.
func (o *MachineType) GetCPU() (value *Value, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cpu
	}
	return
}

// Category returns the value of the 'category' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The category which the machine type is suitable for.
func (o *MachineType) Category() MachineTypeCategory {
	if o != nil && o.bitmap_&32 != 0 {
		return o.category
	}
	return MachineTypeCategory("")
}

// GetCategory returns the value of the 'category' attribute and
// a flag indicating if the attribute has a value.
//
// The category which the machine type is suitable for.
func (o *MachineType) GetCategory() (value MachineTypeCategory, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.category
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cloud provider that the machine type belongs to.
func (o *MachineType) CloudProvider() *CloudProvider {
	if o != nil && o.bitmap_&64 != 0 {
		return o.cloudProvider
	}
	return nil
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cloud provider that the machine type belongs to.
func (o *MachineType) GetCloudProvider() (value *CloudProvider, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// GenericName returns the value of the 'generic_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Generic name for quota purposes, for example `highmem-4`.
// Cloud provider agnostic - many values are shared between "similar"
// machine types on different providers.
// Corresponds to `resource_name` values in "compute.node"  quota cost data.
func (o *MachineType) GenericName() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.genericName
	}
	return ""
}

// GetGenericName returns the value of the 'generic_name' attribute and
// a flag indicating if the attribute has a value.
//
// Generic name for quota purposes, for example `highmem-4`.
// Cloud provider agnostic - many values are shared between "similar"
// machine types on different providers.
// Corresponds to `resource_name` values in "compute.node"  quota cost data.
func (o *MachineType) GetGenericName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.genericName
	}
	return
}

// Memory returns the value of the 'memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The amount of memory of the machine type.
func (o *MachineType) Memory() *Value {
	if o != nil && o.bitmap_&256 != 0 {
		return o.memory
	}
	return nil
}

// GetMemory returns the value of the 'memory' attribute and
// a flag indicating if the attribute has a value.
//
// The amount of memory of the machine type.
func (o *MachineType) GetMemory() (value *Value, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.memory
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Human friendly identifier of the machine type, for example `r5.xlarge - Memory Optimized`.
func (o *MachineType) Name() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Human friendly identifier of the machine type, for example `r5.xlarge - Memory Optimized`.
func (o *MachineType) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.name
	}
	return
}

// Size returns the value of the 'size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The size of the machine type.
func (o *MachineType) Size() MachineTypeSize {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.size
	}
	return MachineTypeSize("")
}

// GetSize returns the value of the 'size' attribute and
// a flag indicating if the attribute has a value.
//
// The size of the machine type.
func (o *MachineType) GetSize() (value MachineTypeSize, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.size
	}
	return
}

// MachineTypeListKind is the name of the type used to represent list of objects of
// type 'machine_type'.
const MachineTypeListKind = "MachineTypeList"

// MachineTypeListLinkKind is the name of the type used to represent links to list
// of objects of type 'machine_type'.
const MachineTypeListLinkKind = "MachineTypeListLink"

// MachineTypeNilKind is the name of the type used to nil lists of objects of
// type 'machine_type'.
const MachineTypeListNilKind = "MachineTypeListNil"

// MachineTypeList is a list of values of the 'machine_type' type.
type MachineTypeList struct {
	href  string
	link  bool
	items []*MachineType
}

// Kind returns the name of the type of the object.
func (l *MachineTypeList) Kind() string {
	if l == nil {
		return MachineTypeListNilKind
	}
	if l.link {
		return MachineTypeListLinkKind
	}
	return MachineTypeListKind
}

// Link returns true iif this is a link.
func (l *MachineTypeList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *MachineTypeList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *MachineTypeList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *MachineTypeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *MachineTypeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MachineTypeList) Get(i int) *MachineType {
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
func (l *MachineTypeList) Slice() []*MachineType {
	var slice []*MachineType
	if l == nil {
		slice = make([]*MachineType, 0)
	} else {
		slice = make([]*MachineType, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MachineTypeList) Each(f func(item *MachineType) bool) {
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
func (l *MachineTypeList) Range(f func(index int, item *MachineType) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
