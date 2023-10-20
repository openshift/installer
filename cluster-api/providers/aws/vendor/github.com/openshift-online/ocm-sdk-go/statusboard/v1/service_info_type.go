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

// ServiceInfo represents the values of the 'service_info' type.
//
// Definition of a Status Board service info.
type ServiceInfo struct {
	bitmap_    uint32
	fullname   string
	statusType string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ServiceInfo) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Fullname returns the value of the 'fullname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Full name of the service
func (o *ServiceInfo) Fullname() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.fullname
	}
	return ""
}

// GetFullname returns the value of the 'fullname' attribute and
// a flag indicating if the attribute has a value.
//
// Full name of the service
func (o *ServiceInfo) GetFullname() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.fullname
	}
	return
}

// StatusType returns the value of the 'status_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of the service status
func (o *ServiceInfo) StatusType() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.statusType
	}
	return ""
}

// GetStatusType returns the value of the 'status_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of the service status
func (o *ServiceInfo) GetStatusType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.statusType
	}
	return
}

// ServiceInfoListKind is the name of the type used to represent list of objects of
// type 'service_info'.
const ServiceInfoListKind = "ServiceInfoList"

// ServiceInfoListLinkKind is the name of the type used to represent links to list
// of objects of type 'service_info'.
const ServiceInfoListLinkKind = "ServiceInfoListLink"

// ServiceInfoNilKind is the name of the type used to nil lists of objects of
// type 'service_info'.
const ServiceInfoListNilKind = "ServiceInfoListNil"

// ServiceInfoList is a list of values of the 'service_info' type.
type ServiceInfoList struct {
	href  string
	link  bool
	items []*ServiceInfo
}

// Len returns the length of the list.
func (l *ServiceInfoList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ServiceInfoList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ServiceInfoList) Get(i int) *ServiceInfo {
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
func (l *ServiceInfoList) Slice() []*ServiceInfo {
	var slice []*ServiceInfo
	if l == nil {
		slice = make([]*ServiceInfo, 0)
	} else {
		slice = make([]*ServiceInfo, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ServiceInfoList) Each(f func(item *ServiceInfo) bool) {
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
func (l *ServiceInfoList) Range(f func(index int, item *ServiceInfo) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
