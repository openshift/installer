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

// WifPool represents the values of the 'wif_pool' type.
type WifPool struct {
	fieldSet_        []bool
	identityProvider *WifIdentityProvider
	poolId           string
	poolName         string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifPool) Empty() bool {
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

// IdentityProvider returns the value of the 'identity_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identity provider configuration data that will be created as part of the
// workload identity pool.
func (o *WifPool) IdentityProvider() *WifIdentityProvider {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.identityProvider
	}
	return nil
}

// GetIdentityProvider returns the value of the 'identity_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Identity provider configuration data that will be created as part of the
// workload identity pool.
func (o *WifPool) GetIdentityProvider() (value *WifIdentityProvider, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.identityProvider
	}
	return
}

// PoolId returns the value of the 'pool_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Id of the workload identity pool.
func (o *WifPool) PoolId() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.poolId
	}
	return ""
}

// GetPoolId returns the value of the 'pool_id' attribute and
// a flag indicating if the attribute has a value.
//
// The Id of the workload identity pool.
func (o *WifPool) GetPoolId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.poolId
	}
	return
}

// PoolName returns the value of the 'pool_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The display name of the workload identity pool.
func (o *WifPool) PoolName() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.poolName
	}
	return ""
}

// GetPoolName returns the value of the 'pool_name' attribute and
// a flag indicating if the attribute has a value.
//
// The display name of the workload identity pool.
func (o *WifPool) GetPoolName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.poolName
	}
	return
}

// WifPoolListKind is the name of the type used to represent list of objects of
// type 'wif_pool'.
const WifPoolListKind = "WifPoolList"

// WifPoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_pool'.
const WifPoolListLinkKind = "WifPoolListLink"

// WifPoolNilKind is the name of the type used to nil lists of objects of
// type 'wif_pool'.
const WifPoolListNilKind = "WifPoolListNil"

// WifPoolList is a list of values of the 'wif_pool' type.
type WifPoolList struct {
	href  string
	link  bool
	items []*WifPool
}

// Len returns the length of the list.
func (l *WifPoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *WifPoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *WifPoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *WifPoolList) SetItems(items []*WifPool) {
	l.items = items
}

// Items returns the items of the list.
func (l *WifPoolList) Items() []*WifPool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *WifPoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifPoolList) Get(i int) *WifPool {
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
func (l *WifPoolList) Slice() []*WifPool {
	var slice []*WifPool
	if l == nil {
		slice = make([]*WifPool, 0)
	} else {
		slice = make([]*WifPool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifPoolList) Each(f func(item *WifPool) bool) {
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
func (l *WifPoolList) Range(f func(index int, item *WifPool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
