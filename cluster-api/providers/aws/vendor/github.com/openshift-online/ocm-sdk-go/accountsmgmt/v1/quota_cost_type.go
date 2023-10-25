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

// QuotaCost represents the values of the 'quota_cost' type.
type QuotaCost struct {
	bitmap_          uint32
	allowed          int
	cloudAccounts    []*CloudAccount
	consumed         int
	organizationID   string
	quotaID          string
	relatedResources []*RelatedResource
	version          string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *QuotaCost) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Allowed returns the value of the 'allowed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) Allowed() int {
	if o != nil && o.bitmap_&1 != 0 {
		return o.allowed
	}
	return 0
}

// GetAllowed returns the value of the 'allowed' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetAllowed() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.allowed
	}
	return
}

// CloudAccounts returns the value of the 'cloud_accounts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) CloudAccounts() []*CloudAccount {
	if o != nil && o.bitmap_&2 != 0 {
		return o.cloudAccounts
	}
	return nil
}

// GetCloudAccounts returns the value of the 'cloud_accounts' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetCloudAccounts() (value []*CloudAccount, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.cloudAccounts
	}
	return
}

// Consumed returns the value of the 'consumed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) Consumed() int {
	if o != nil && o.bitmap_&4 != 0 {
		return o.consumed
	}
	return 0
}

// GetConsumed returns the value of the 'consumed' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetConsumed() (value int, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.consumed
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) OrganizationID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.organizationID
	}
	return
}

// QuotaID returns the value of the 'quota_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) QuotaID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.quotaID
	}
	return ""
}

// GetQuotaID returns the value of the 'quota_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetQuotaID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.quotaID
	}
	return
}

// RelatedResources returns the value of the 'related_resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) RelatedResources() []*RelatedResource {
	if o != nil && o.bitmap_&32 != 0 {
		return o.relatedResources
	}
	return nil
}

// GetRelatedResources returns the value of the 'related_resources' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetRelatedResources() (value []*RelatedResource, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.relatedResources
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaCost) Version() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaCost) GetVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.version
	}
	return
}

// QuotaCostListKind is the name of the type used to represent list of objects of
// type 'quota_cost'.
const QuotaCostListKind = "QuotaCostList"

// QuotaCostListLinkKind is the name of the type used to represent links to list
// of objects of type 'quota_cost'.
const QuotaCostListLinkKind = "QuotaCostListLink"

// QuotaCostNilKind is the name of the type used to nil lists of objects of
// type 'quota_cost'.
const QuotaCostListNilKind = "QuotaCostListNil"

// QuotaCostList is a list of values of the 'quota_cost' type.
type QuotaCostList struct {
	href  string
	link  bool
	items []*QuotaCost
}

// Len returns the length of the list.
func (l *QuotaCostList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *QuotaCostList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *QuotaCostList) Get(i int) *QuotaCost {
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
func (l *QuotaCostList) Slice() []*QuotaCost {
	var slice []*QuotaCost
	if l == nil {
		slice = make([]*QuotaCost, 0)
	} else {
		slice = make([]*QuotaCost, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *QuotaCostList) Each(f func(item *QuotaCost) bool) {
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
func (l *QuotaCostList) Range(f func(index int, item *QuotaCost) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
