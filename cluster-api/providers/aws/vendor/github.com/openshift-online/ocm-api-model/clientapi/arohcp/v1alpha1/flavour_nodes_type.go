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

// FlavourNodes represents the values of the 'flavour_nodes' type.
//
// Counts of different classes of nodes inside a flavour.
type FlavourNodes struct {
	fieldSet_ []bool
	master    int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *FlavourNodes) Empty() bool {
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

// Master returns the value of the 'master' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of master nodes of the cluster.
func (o *FlavourNodes) Master() int {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.master
	}
	return 0
}

// GetMaster returns the value of the 'master' attribute and
// a flag indicating if the attribute has a value.
//
// Number of master nodes of the cluster.
func (o *FlavourNodes) GetMaster() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.master
	}
	return
}

// FlavourNodesListKind is the name of the type used to represent list of objects of
// type 'flavour_nodes'.
const FlavourNodesListKind = "FlavourNodesList"

// FlavourNodesListLinkKind is the name of the type used to represent links to list
// of objects of type 'flavour_nodes'.
const FlavourNodesListLinkKind = "FlavourNodesListLink"

// FlavourNodesNilKind is the name of the type used to nil lists of objects of
// type 'flavour_nodes'.
const FlavourNodesListNilKind = "FlavourNodesListNil"

// FlavourNodesList is a list of values of the 'flavour_nodes' type.
type FlavourNodesList struct {
	href  string
	link  bool
	items []*FlavourNodes
}

// Len returns the length of the list.
func (l *FlavourNodesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *FlavourNodesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *FlavourNodesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *FlavourNodesList) SetItems(items []*FlavourNodes) {
	l.items = items
}

// Items returns the items of the list.
func (l *FlavourNodesList) Items() []*FlavourNodes {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *FlavourNodesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FlavourNodesList) Get(i int) *FlavourNodes {
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
func (l *FlavourNodesList) Slice() []*FlavourNodes {
	var slice []*FlavourNodes
	if l == nil {
		slice = make([]*FlavourNodes, 0)
	} else {
		slice = make([]*FlavourNodes, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FlavourNodesList) Each(f func(item *FlavourNodes) bool) {
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
func (l *FlavourNodesList) Range(f func(index int, item *FlavourNodes) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
