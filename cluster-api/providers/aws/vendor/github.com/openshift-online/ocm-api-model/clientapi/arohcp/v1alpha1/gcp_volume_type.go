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

// GCPVolume represents the values of the 'GCP_volume' type.
//
// Holds settings for an GCP storage volume.
type GCPVolume struct {
	fieldSet_ []bool
	size      int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPVolume) Empty() bool {
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

// Size returns the value of the 'size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Volume size in Gib.
func (o *GCPVolume) Size() int {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.size
	}
	return 0
}

// GetSize returns the value of the 'size' attribute and
// a flag indicating if the attribute has a value.
//
// Volume size in Gib.
func (o *GCPVolume) GetSize() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.size
	}
	return
}

// GCPVolumeListKind is the name of the type used to represent list of objects of
// type 'GCP_volume'.
const GCPVolumeListKind = "GCPVolumeList"

// GCPVolumeListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_volume'.
const GCPVolumeListLinkKind = "GCPVolumeListLink"

// GCPVolumeNilKind is the name of the type used to nil lists of objects of
// type 'GCP_volume'.
const GCPVolumeListNilKind = "GCPVolumeListNil"

// GCPVolumeList is a list of values of the 'GCP_volume' type.
type GCPVolumeList struct {
	href  string
	link  bool
	items []*GCPVolume
}

// Len returns the length of the list.
func (l *GCPVolumeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPVolumeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPVolumeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPVolumeList) SetItems(items []*GCPVolume) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPVolumeList) Items() []*GCPVolume {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPVolumeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPVolumeList) Get(i int) *GCPVolume {
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
func (l *GCPVolumeList) Slice() []*GCPVolume {
	var slice []*GCPVolume
	if l == nil {
		slice = make([]*GCPVolume, 0)
	} else {
		slice = make([]*GCPVolume, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPVolumeList) Each(f func(item *GCPVolume) bool) {
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
func (l *GCPVolumeList) Range(f func(index int, item *GCPVolume) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
