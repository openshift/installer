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

// CloudAccount represents the values of the 'cloud_account' type.
type CloudAccount struct {
	bitmap_         uint32
	cloudAccountID  string
	cloudProviderID string
	contracts       []*Contract
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CloudAccount) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CloudAccountID returns the value of the 'cloud_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudAccount) CloudAccountID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cloudAccountID
	}
	return ""
}

// GetCloudAccountID returns the value of the 'cloud_account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudAccount) GetCloudAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cloudAccountID
	}
	return
}

// CloudProviderID returns the value of the 'cloud_provider_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudAccount) CloudProviderID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.cloudProviderID
	}
	return ""
}

// GetCloudProviderID returns the value of the 'cloud_provider_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudAccount) GetCloudProviderID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.cloudProviderID
	}
	return
}

// Contracts returns the value of the 'contracts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CloudAccount) Contracts() []*Contract {
	if o != nil && o.bitmap_&4 != 0 {
		return o.contracts
	}
	return nil
}

// GetContracts returns the value of the 'contracts' attribute and
// a flag indicating if the attribute has a value.
func (o *CloudAccount) GetContracts() (value []*Contract, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.contracts
	}
	return
}

// CloudAccountListKind is the name of the type used to represent list of objects of
// type 'cloud_account'.
const CloudAccountListKind = "CloudAccountList"

// CloudAccountListLinkKind is the name of the type used to represent links to list
// of objects of type 'cloud_account'.
const CloudAccountListLinkKind = "CloudAccountListLink"

// CloudAccountNilKind is the name of the type used to nil lists of objects of
// type 'cloud_account'.
const CloudAccountListNilKind = "CloudAccountListNil"

// CloudAccountList is a list of values of the 'cloud_account' type.
type CloudAccountList struct {
	href  string
	link  bool
	items []*CloudAccount
}

// Len returns the length of the list.
func (l *CloudAccountList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *CloudAccountList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CloudAccountList) Get(i int) *CloudAccount {
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
func (l *CloudAccountList) Slice() []*CloudAccount {
	var slice []*CloudAccount
	if l == nil {
		slice = make([]*CloudAccount, 0)
	} else {
		slice = make([]*CloudAccount, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CloudAccountList) Each(f func(item *CloudAccount) bool) {
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
func (l *CloudAccountList) Range(f func(index int, item *CloudAccount) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
