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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// Network represents the values of the 'network' type.
//
// Network configuration of a cluster.
type Network struct {
	bitmap_     uint32
	hostPrefix  int
	machineCIDR string
	podCIDR     string
	serviceCIDR string
	type_       string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Network) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HostPrefix returns the value of the 'host_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Network host prefix which is defaulted to `23` if not specified.
func (o *Network) HostPrefix() int {
	if o != nil && o.bitmap_&1 != 0 {
		return o.hostPrefix
	}
	return 0
}

// GetHostPrefix returns the value of the 'host_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// Network host prefix which is defaulted to `23` if not specified.
func (o *Network) GetHostPrefix() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.hostPrefix
	}
	return
}

// MachineCIDR returns the value of the 'machine_CIDR' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// IP address block from which to assign machine IP addresses, for example `10.0.0.0/16`.
func (o *Network) MachineCIDR() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.machineCIDR
	}
	return ""
}

// GetMachineCIDR returns the value of the 'machine_CIDR' attribute and
// a flag indicating if the attribute has a value.
//
// IP address block from which to assign machine IP addresses, for example `10.0.0.0/16`.
func (o *Network) GetMachineCIDR() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.machineCIDR
	}
	return
}

// PodCIDR returns the value of the 'pod_CIDR' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// IP address block from which to assign pod IP addresses, for example `10.128.0.0/14`.
func (o *Network) PodCIDR() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.podCIDR
	}
	return ""
}

// GetPodCIDR returns the value of the 'pod_CIDR' attribute and
// a flag indicating if the attribute has a value.
//
// IP address block from which to assign pod IP addresses, for example `10.128.0.0/14`.
func (o *Network) GetPodCIDR() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.podCIDR
	}
	return
}

// ServiceCIDR returns the value of the 'service_CIDR' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// IP address block from which to assign service IP addresses, for example `172.30.0.0/16`.
func (o *Network) ServiceCIDR() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.serviceCIDR
	}
	return ""
}

// GetServiceCIDR returns the value of the 'service_CIDR' attribute and
// a flag indicating if the attribute has a value.
//
// IP address block from which to assign service IP addresses, for example `172.30.0.0/16`.
func (o *Network) GetServiceCIDR() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.serviceCIDR
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The main controller responsible for rendering the core networking components.
func (o *Network) Type() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// The main controller responsible for rendering the core networking components.
func (o *Network) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.type_
	}
	return
}

// NetworkListKind is the name of the type used to represent list of objects of
// type 'network'.
const NetworkListKind = "NetworkList"

// NetworkListLinkKind is the name of the type used to represent links to list
// of objects of type 'network'.
const NetworkListLinkKind = "NetworkListLink"

// NetworkNilKind is the name of the type used to nil lists of objects of
// type 'network'.
const NetworkListNilKind = "NetworkListNil"

// NetworkList is a list of values of the 'network' type.
type NetworkList struct {
	href  string
	link  bool
	items []*Network
}

// Len returns the length of the list.
func (l *NetworkList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NetworkList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NetworkList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NetworkList) SetItems(items []*Network) {
	l.items = items
}

// Items returns the items of the list.
func (l *NetworkList) Items() []*Network {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NetworkList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NetworkList) Get(i int) *Network {
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
func (l *NetworkList) Slice() []*Network {
	var slice []*Network
	if l == nil {
		slice = make([]*Network, 0)
	} else {
		slice = make([]*Network, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NetworkList) Each(f func(item *Network) bool) {
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
func (l *NetworkList) Range(f func(index int, item *Network) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
