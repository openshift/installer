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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	time "time"
)

// ApplicationKind is the name of the type used to represent objects
// of type 'application'.
const ApplicationKind = "Application"

// ApplicationLinkKind is the name of the type used to represent links
// to objects of type 'application'.
const ApplicationLinkKind = "ApplicationLink"

// ApplicationNilKind is the name of the type used to nil references
// to objects of type 'application'.
const ApplicationNilKind = "ApplicationNil"

// Application represents the values of the 'application' type.
//
// Definition of a Status Board application.
type Application struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	fullname  string
	metadata  interface{}
	name      string
	owners    []*Owner
	product   *Product
	updatedAt time.Time
}

// Kind returns the name of the type of the object.
func (o *Application) Kind() string {
	if o == nil {
		return ApplicationNilKind
	}
	if o.bitmap_&1 != 0 {
		return ApplicationLinkKind
	}
	return ApplicationKind
}

// Link returns true iif this is a link.
func (o *Application) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Application) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Application) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Application) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Application) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Application) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Application) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Application) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Fullname returns the value of the 'fullname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The full name of the application.
func (o *Application) Fullname() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.fullname
	}
	return ""
}

// GetFullname returns the value of the 'fullname' attribute and
// a flag indicating if the attribute has a value.
//
// The full name of the application.
func (o *Application) GetFullname() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.fullname
	}
	return
}

// Metadata returns the value of the 'metadata' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Miscellaneous metadata about the application.
func (o *Application) Metadata() interface{} {
	if o != nil && o.bitmap_&32 != 0 {
		return o.metadata
	}
	return nil
}

// GetMetadata returns the value of the 'metadata' attribute and
// a flag indicating if the attribute has a value.
//
// Miscellaneous metadata about the application.
func (o *Application) GetMetadata() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.metadata
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the application.
func (o *Application) Name() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the application.
func (o *Application) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.name
	}
	return
}

// Owners returns the value of the 'owners' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The application owners (name and email)
func (o *Application) Owners() []*Owner {
	if o != nil && o.bitmap_&128 != 0 {
		return o.owners
	}
	return nil
}

// GetOwners returns the value of the 'owners' attribute and
// a flag indicating if the attribute has a value.
//
// The application owners (name and email)
func (o *Application) GetOwners() (value []*Owner, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.owners
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The group ID that the application belongs to.
func (o *Application) Product() *Product {
	if o != nil && o.bitmap_&256 != 0 {
		return o.product
	}
	return nil
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
//
// The group ID that the application belongs to.
func (o *Application) GetProduct() (value *Product, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.product
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Application) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&512 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Application) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// ApplicationListKind is the name of the type used to represent list of objects of
// type 'application'.
const ApplicationListKind = "ApplicationList"

// ApplicationListLinkKind is the name of the type used to represent links to list
// of objects of type 'application'.
const ApplicationListLinkKind = "ApplicationListLink"

// ApplicationNilKind is the name of the type used to nil lists of objects of
// type 'application'.
const ApplicationListNilKind = "ApplicationListNil"

// ApplicationList is a list of values of the 'application' type.
type ApplicationList struct {
	href  string
	link  bool
	items []*Application
}

// Kind returns the name of the type of the object.
func (l *ApplicationList) Kind() string {
	if l == nil {
		return ApplicationListNilKind
	}
	if l.link {
		return ApplicationListLinkKind
	}
	return ApplicationListKind
}

// Link returns true iif this is a link.
func (l *ApplicationList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ApplicationList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ApplicationList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ApplicationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ApplicationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ApplicationList) Get(i int) *Application {
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
func (l *ApplicationList) Slice() []*Application {
	var slice []*Application
	if l == nil {
		slice = make([]*Application, 0)
	} else {
		slice = make([]*Application, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ApplicationList) Each(f func(item *Application) bool) {
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
func (l *ApplicationList) Range(f func(index int, item *Application) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
