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

// CloudProviderKind is the name of the type used to represent objects
// of type 'cloud_provider'.
const CloudProviderKind = "CloudProvider"

// CloudProviderLinkKind is the name of the type used to represent links
// to objects of type 'cloud_provider'.
const CloudProviderLinkKind = "CloudProviderLink"

// CloudProviderNilKind is the name of the type used to nil references
// to objects of type 'cloud_provider'.
const CloudProviderNilKind = "CloudProviderNil"

// CloudProvider represents the values of the 'cloud_provider' type.
//
// Cloud provider.
type CloudProvider struct {
	fieldSet_   []bool
	id          string
	href        string
	displayName string
	name        string
	regions     []*CloudRegion
}

// Kind returns the name of the type of the object.
func (o *CloudProvider) Kind() string {
	if o == nil {
		return CloudProviderNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return CloudProviderLinkKind
	}
	return CloudProviderKind
}

// Link returns true if this is a link.
func (o *CloudProvider) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *CloudProvider) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *CloudProvider) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *CloudProvider) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *CloudProvider) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudProvider) Empty() bool {
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

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the cloud provider for display purposes. It can contain any characters,
// including spaces.
func (o *CloudProvider) DisplayName() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the cloud provider for display purposes. It can contain any characters,
// including spaces.
func (o *CloudProvider) GetDisplayName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.displayName
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Human friendly identifier of the cloud provider, for example `aws`.
func (o *CloudProvider) Name() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Human friendly identifier of the cloud provider, for example `aws`.
func (o *CloudProvider) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.name
	}
	return
}

// Regions returns the value of the 'regions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// (optional) Provider's regions - only included when listing providers with `fetchRegions=true`.
func (o *CloudProvider) Regions() []*CloudRegion {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.regions
	}
	return nil
}

// GetRegions returns the value of the 'regions' attribute and
// a flag indicating if the attribute has a value.
//
// (optional) Provider's regions - only included when listing providers with `fetchRegions=true`.
func (o *CloudProvider) GetRegions() (value []*CloudRegion, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.regions
	}
	return
}

// CloudProviderListKind is the name of the type used to represent list of objects of
// type 'cloud_provider'.
const CloudProviderListKind = "CloudProviderList"

// CloudProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_provider'.
const CloudProviderListLinkKind = "CloudProviderListLink"

// CloudProviderNilKind is the name of the type used to nil lists of objects of
// type 'cloud_provider'.
const CloudProviderListNilKind = "CloudProviderListNil"

// CloudProviderList is a list of values of the 'cloud_provider' type.
type CloudProviderList struct {
	href  string
	link  bool
	items []*CloudProvider
}

// Kind returns the name of the type of the object.
func (l *CloudProviderList) Kind() string {
	if l == nil {
		return CloudProviderListNilKind
	}
	if l.link {
		return CloudProviderListLinkKind
	}
	return CloudProviderListKind
}

// Link returns true iif this is a link.
func (l *CloudProviderList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *CloudProviderList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *CloudProviderList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *CloudProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CloudProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CloudProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CloudProviderList) SetItems(items []*CloudProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *CloudProviderList) Items() []*CloudProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CloudProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudProviderList) Get(i int) *CloudProvider {
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
func (l *CloudProviderList) Slice() []*CloudProvider {
	var slice []*CloudProvider
	if l == nil {
		slice = make([]*CloudProvider, 0)
	} else {
		slice = make([]*CloudProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudProviderList) Each(f func(item *CloudProvider) bool) {
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
func (l *CloudProviderList) Range(f func(index int, item *CloudProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
