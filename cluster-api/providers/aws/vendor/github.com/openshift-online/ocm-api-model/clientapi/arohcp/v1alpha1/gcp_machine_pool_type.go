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

// GCPMachinePool represents the values of the 'GCP_machine_pool' type.
//
// Representation of gcp machine pool specific parameters.
type GCPMachinePool struct {
	fieldSet_  []bool
	secureBoot bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPMachinePool) Empty() bool {
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

// SecureBoot returns the value of the 'secure_boot' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Determines whether the Shielded VM's SecureBoot feature should be
// enabled for the nodes of the machine pool. If SecureBoot is not
// specified, the value of this attribute will remain unspecified and the
// SecureBoot's value specified in the `.gcp.security.secure_boot`
// attribute of the parent Cluster will be the one applied to the nodes of
// the machine pool.
// Immutable.
func (o *GCPMachinePool) SecureBoot() bool {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.secureBoot
	}
	return false
}

// GetSecureBoot returns the value of the 'secure_boot' attribute and
// a flag indicating if the attribute has a value.
//
// Determines whether the Shielded VM's SecureBoot feature should be
// enabled for the nodes of the machine pool. If SecureBoot is not
// specified, the value of this attribute will remain unspecified and the
// SecureBoot's value specified in the `.gcp.security.secure_boot`
// attribute of the parent Cluster will be the one applied to the nodes of
// the machine pool.
// Immutable.
func (o *GCPMachinePool) GetSecureBoot() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.secureBoot
	}
	return
}

// GCPMachinePoolListKind is the name of the type used to represent list of objects of
// type 'GCP_machine_pool'.
const GCPMachinePoolListKind = "GCPMachinePoolList"

// GCPMachinePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_machine_pool'.
const GCPMachinePoolListLinkKind = "GCPMachinePoolListLink"

// GCPMachinePoolNilKind is the name of the type used to nil lists of objects of
// type 'GCP_machine_pool'.
const GCPMachinePoolListNilKind = "GCPMachinePoolListNil"

// GCPMachinePoolList is a list of values of the 'GCP_machine_pool' type.
type GCPMachinePoolList struct {
	href  string
	link  bool
	items []*GCPMachinePool
}

// Len returns the length of the list.
func (l *GCPMachinePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPMachinePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPMachinePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPMachinePoolList) SetItems(items []*GCPMachinePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPMachinePoolList) Items() []*GCPMachinePool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPMachinePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPMachinePoolList) Get(i int) *GCPMachinePool {
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
func (l *GCPMachinePoolList) Slice() []*GCPMachinePool {
	var slice []*GCPMachinePool
	if l == nil {
		slice = make([]*GCPMachinePool, 0)
	} else {
		slice = make([]*GCPMachinePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPMachinePoolList) Each(f func(item *GCPMachinePool) bool) {
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
func (l *GCPMachinePoolList) Range(f func(index int, item *GCPMachinePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
