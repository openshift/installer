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

import (
	time "time"
)

// CloudResourceKind is the name of the type used to represent objects
// of type 'cloud_resource'.
const CloudResourceKind = "CloudResource"

// CloudResourceLinkKind is the name of the type used to represent links
// to objects of type 'cloud_resource'.
const CloudResourceLinkKind = "CloudResourceLink"

// CloudResourceNilKind is the name of the type used to nil references
// to objects of type 'cloud_resource'.
const CloudResourceNilKind = "CloudResourceNil"

// CloudResource represents the values of the 'cloud_resource' type.
type CloudResource struct {
	bitmap_        uint32
	id             string
	href           string
	category       string
	categoryPretty string
	cloudProvider  string
	cpuCores       int
	createdAt      time.Time
	genericName    string
	memory         int
	memoryPretty   string
	namePretty     string
	resourceType   string
	sizePretty     string
	updatedAt      time.Time
	active         bool
}

// Kind returns the name of the type of the object.
func (o *CloudResource) Kind() string {
	if o == nil {
		return CloudResourceNilKind
	}
	if o.bitmap_&1 != 0 {
		return CloudResourceLinkKind
	}
	return CloudResourceKind
}

// Link returns true iif this is a link.
func (o *CloudResource) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *CloudResource) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *CloudResource) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *CloudResource) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *CloudResource) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudResource) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Active returns the value of the 'active' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) Active() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.active
	}
	return false
}

// GetActive returns the value of the 'active' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetActive() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.active
	}
	return
}

// Category returns the value of the 'category' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) Category() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.category
	}
	return ""
}

// GetCategory returns the value of the 'category' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetCategory() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.category
	}
	return
}

// CategoryPretty returns the value of the 'category_pretty' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) CategoryPretty() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.categoryPretty
	}
	return ""
}

// GetCategoryPretty returns the value of the 'category_pretty' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetCategoryPretty() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.categoryPretty
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) CloudProvider() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.cloudProvider
	}
	return ""
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetCloudProvider() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// CpuCores returns the value of the 'cpu_cores' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) CpuCores() int {
	if o != nil && o.bitmap_&128 != 0 {
		return o.cpuCores
	}
	return 0
}

// GetCpuCores returns the value of the 'cpu_cores' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetCpuCores() (value int, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.cpuCores
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) CreatedAt() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// GenericName returns the value of the 'generic_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) GenericName() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.genericName
	}
	return ""
}

// GetGenericName returns the value of the 'generic_name' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetGenericName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.genericName
	}
	return
}

// Memory returns the value of the 'memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) Memory() int {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.memory
	}
	return 0
}

// GetMemory returns the value of the 'memory' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetMemory() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.memory
	}
	return
}

// MemoryPretty returns the value of the 'memory_pretty' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) MemoryPretty() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.memoryPretty
	}
	return ""
}

// GetMemoryPretty returns the value of the 'memory_pretty' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetMemoryPretty() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.memoryPretty
	}
	return
}

// NamePretty returns the value of the 'name_pretty' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) NamePretty() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.namePretty
	}
	return ""
}

// GetNamePretty returns the value of the 'name_pretty' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetNamePretty() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.namePretty
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) ResourceType() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetResourceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.resourceType
	}
	return
}

// SizePretty returns the value of the 'size_pretty' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) SizePretty() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.sizePretty
	}
	return ""
}

// GetSizePretty returns the value of the 'size_pretty' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetSizePretty() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.sizePretty
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudResource) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudResource) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// CloudResourceListKind is the name of the type used to represent list of objects of
// type 'cloud_resource'.
const CloudResourceListKind = "CloudResourceList"

// CloudResourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_resource'.
const CloudResourceListLinkKind = "CloudResourceListLink"

// CloudResourceNilKind is the name of the type used to nil lists of objects of
// type 'cloud_resource'.
const CloudResourceListNilKind = "CloudResourceListNil"

// CloudResourceList is a list of values of the 'cloud_resource' type.
type CloudResourceList struct {
	href  string
	link  bool
	items []*CloudResource
}

// Kind returns the name of the type of the object.
func (l *CloudResourceList) Kind() string {
	if l == nil {
		return CloudResourceListNilKind
	}
	if l.link {
		return CloudResourceListLinkKind
	}
	return CloudResourceListKind
}

// Link returns true iif this is a link.
func (l *CloudResourceList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *CloudResourceList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *CloudResourceList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *CloudResourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *CloudResourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudResourceList) Get(i int) *CloudResource {
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
func (l *CloudResourceList) Slice() []*CloudResource {
	var slice []*CloudResource
	if l == nil {
		slice = make([]*CloudResource, 0)
	} else {
		slice = make([]*CloudResource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudResourceList) Each(f func(item *CloudResource) bool) {
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
func (l *CloudResourceList) Range(f func(index int, item *CloudResource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
