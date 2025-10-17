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

// ApplicationDependencyKind is the name of the type used to represent objects
// of type 'application_dependency'.
const ApplicationDependencyKind = "ApplicationDependency"

// ApplicationDependencyLinkKind is the name of the type used to represent links
// to objects of type 'application_dependency'.
const ApplicationDependencyLinkKind = "ApplicationDependencyLink"

// ApplicationDependencyNilKind is the name of the type used to nil references
// to objects of type 'application_dependency'.
const ApplicationDependencyNilKind = "ApplicationDependencyNil"

// ApplicationDependency represents the values of the 'application_dependency' type.
//
// Definition of a Status Board application dependency.
type ApplicationDependency struct {
	bitmap_     uint32
	id          string
	href        string
	application *Application
	createdAt   time.Time
	metadata    interface{}
	name        string
	owners      []*Owner
	service     *Service
	type_       string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *ApplicationDependency) Kind() string {
	if o == nil {
		return ApplicationDependencyNilKind
	}
	if o.bitmap_&1 != 0 {
		return ApplicationDependencyLinkKind
	}
	return ApplicationDependencyKind
}

// Link returns true if this is a link.
func (o *ApplicationDependency) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ApplicationDependency) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ApplicationDependency) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ApplicationDependency) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ApplicationDependency) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ApplicationDependency) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Application returns the value of the 'application' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The parent Application of the dependency.
func (o *ApplicationDependency) Application() *Application {
	if o != nil && o.bitmap_&8 != 0 {
		return o.application
	}
	return nil
}

// GetApplication returns the value of the 'application' attribute and
// a flag indicating if the attribute has a value.
//
// The parent Application of the dependency.
func (o *ApplicationDependency) GetApplication() (value *Application, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.application
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *ApplicationDependency) CreatedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *ApplicationDependency) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Metadata returns the value of the 'metadata' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Miscellaneous metadata about the application.
func (o *ApplicationDependency) Metadata() interface{} {
	if o != nil && o.bitmap_&32 != 0 {
		return o.metadata
	}
	return nil
}

// GetMetadata returns the value of the 'metadata' attribute and
// a flag indicating if the attribute has a value.
//
// Miscellaneous metadata about the application.
func (o *ApplicationDependency) GetMetadata() (value interface{}, ok bool) {
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
func (o *ApplicationDependency) Name() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the application.
func (o *ApplicationDependency) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.name
	}
	return
}

// Owners returns the value of the 'owners' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The application dependency owners (name and email)
func (o *ApplicationDependency) Owners() []*Owner {
	if o != nil && o.bitmap_&128 != 0 {
		return o.owners
	}
	return nil
}

// GetOwners returns the value of the 'owners' attribute and
// a flag indicating if the attribute has a value.
//
// The application dependency owners (name and email)
func (o *ApplicationDependency) GetOwners() (value []*Owner, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.owners
	}
	return
}

// Service returns the value of the 'service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Service associated with the dependency.
func (o *ApplicationDependency) Service() *Service {
	if o != nil && o.bitmap_&256 != 0 {
		return o.service
	}
	return nil
}

// GetService returns the value of the 'service' attribute and
// a flag indicating if the attribute has a value.
//
// The Service associated with the dependency.
func (o *ApplicationDependency) GetService() (value *Service, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.service
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The type of application dependency, e.g. soft or hard.
func (o *ApplicationDependency) Type() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// The type of application dependency, e.g. soft or hard.
func (o *ApplicationDependency) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.type_
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *ApplicationDependency) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *ApplicationDependency) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// ApplicationDependencyListKind is the name of the type used to represent list of objects of
// type 'application_dependency'.
const ApplicationDependencyListKind = "ApplicationDependencyList"

// ApplicationDependencyListLinkKind is the name of the type used to represent links to list
// of objects of type 'application_dependency'.
const ApplicationDependencyListLinkKind = "ApplicationDependencyListLink"

// ApplicationDependencyNilKind is the name of the type used to nil lists of objects of
// type 'application_dependency'.
const ApplicationDependencyListNilKind = "ApplicationDependencyListNil"

// ApplicationDependencyList is a list of values of the 'application_dependency' type.
type ApplicationDependencyList struct {
	href  string
	link  bool
	items []*ApplicationDependency
}

// Kind returns the name of the type of the object.
func (l *ApplicationDependencyList) Kind() string {
	if l == nil {
		return ApplicationDependencyListNilKind
	}
	if l.link {
		return ApplicationDependencyListLinkKind
	}
	return ApplicationDependencyListKind
}

// Link returns true iif this is a link.
func (l *ApplicationDependencyList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ApplicationDependencyList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ApplicationDependencyList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ApplicationDependencyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ApplicationDependencyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ApplicationDependencyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ApplicationDependencyList) SetItems(items []*ApplicationDependency) {
	l.items = items
}

// Items returns the items of the list.
func (l *ApplicationDependencyList) Items() []*ApplicationDependency {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ApplicationDependencyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ApplicationDependencyList) Get(i int) *ApplicationDependency {
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
func (l *ApplicationDependencyList) Slice() []*ApplicationDependency {
	var slice []*ApplicationDependency
	if l == nil {
		slice = make([]*ApplicationDependency, 0)
	} else {
		slice = make([]*ApplicationDependency, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ApplicationDependencyList) Each(f func(item *ApplicationDependency) bool) {
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
func (l *ApplicationDependencyList) Range(f func(index int, item *ApplicationDependency) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
