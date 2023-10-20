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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AWSVolume represents the values of the 'AWS_volume' type.
//
// Holds settings for an AWS storage volume.
type AWSVolume struct {
	bitmap_ uint32
	iops    int
	size    int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSVolume) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// IOPS returns the value of the 'IOPS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Volume provisioned IOPS.
func (o *AWSVolume) IOPS() int {
	if o != nil && o.bitmap_&1 != 0 {
		return o.iops
	}
	return 0
}

// GetIOPS returns the value of the 'IOPS' attribute and
// a flag indicating if the attribute has a value.
//
// Volume provisioned IOPS.
func (o *AWSVolume) GetIOPS() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.iops
	}
	return
}

// Size returns the value of the 'size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Volume size in Gib.
func (o *AWSVolume) Size() int {
	if o != nil && o.bitmap_&2 != 0 {
		return o.size
	}
	return 0
}

// GetSize returns the value of the 'size' attribute and
// a flag indicating if the attribute has a value.
//
// Volume size in Gib.
func (o *AWSVolume) GetSize() (value int, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.size
	}
	return
}

// AWSVolumeListKind is the name of the type used to represent list of objects of
// type 'AWS_volume'.
const AWSVolumeListKind = "AWSVolumeList"

// AWSVolumeListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_volume'.
const AWSVolumeListLinkKind = "AWSVolumeListLink"

// AWSVolumeNilKind is the name of the type used to nil lists of objects of
// type 'AWS_volume'.
const AWSVolumeListNilKind = "AWSVolumeListNil"

// AWSVolumeList is a list of values of the 'AWS_volume' type.
type AWSVolumeList struct {
	href  string
	link  bool
	items []*AWSVolume
}

// Len returns the length of the list.
func (l *AWSVolumeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AWSVolumeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSVolumeList) Get(i int) *AWSVolume {
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
func (l *AWSVolumeList) Slice() []*AWSVolume {
	var slice []*AWSVolume
	if l == nil {
		slice = make([]*AWSVolume, 0)
	} else {
		slice = make([]*AWSVolume, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSVolumeList) Each(f func(item *AWSVolume) bool) {
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
func (l *AWSVolumeList) Range(f func(index int, item *AWSVolume) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
