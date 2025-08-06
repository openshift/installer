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

// GCPNetwork represents the values of the 'GCP_network' type.
//
// GCP Network configuration of a cluster.
type GCPNetwork struct {
	bitmap_            uint32
	vpcName            string
	vpcProjectID       string
	computeSubnet      string
	controlPlaneSubnet string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPNetwork) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// VPCName returns the value of the 'VPC_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// VPC mame used by the cluster.
func (o *GCPNetwork) VPCName() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.vpcName
	}
	return ""
}

// GetVPCName returns the value of the 'VPC_name' attribute and
// a flag indicating if the attribute has a value.
//
// VPC mame used by the cluster.
func (o *GCPNetwork) GetVPCName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.vpcName
	}
	return
}

// VPCProjectID returns the value of the 'VPC_project_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the host project where the shared VPC exists.
func (o *GCPNetwork) VPCProjectID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.vpcProjectID
	}
	return ""
}

// GetVPCProjectID returns the value of the 'VPC_project_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the host project where the shared VPC exists.
func (o *GCPNetwork) GetVPCProjectID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.vpcProjectID
	}
	return
}

// ComputeSubnet returns the value of the 'compute_subnet' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Compute subnet used by the cluster.
func (o *GCPNetwork) ComputeSubnet() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.computeSubnet
	}
	return ""
}

// GetComputeSubnet returns the value of the 'compute_subnet' attribute and
// a flag indicating if the attribute has a value.
//
// Compute subnet used by the cluster.
func (o *GCPNetwork) GetComputeSubnet() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.computeSubnet
	}
	return
}

// ControlPlaneSubnet returns the value of the 'control_plane_subnet' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Control plane subnet used by the cluster.
func (o *GCPNetwork) ControlPlaneSubnet() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.controlPlaneSubnet
	}
	return ""
}

// GetControlPlaneSubnet returns the value of the 'control_plane_subnet' attribute and
// a flag indicating if the attribute has a value.
//
// Control plane subnet used by the cluster.
func (o *GCPNetwork) GetControlPlaneSubnet() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.controlPlaneSubnet
	}
	return
}

// GCPNetworkListKind is the name of the type used to represent list of objects of
// type 'GCP_network'.
const GCPNetworkListKind = "GCPNetworkList"

// GCPNetworkListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_network'.
const GCPNetworkListLinkKind = "GCPNetworkListLink"

// GCPNetworkNilKind is the name of the type used to nil lists of objects of
// type 'GCP_network'.
const GCPNetworkListNilKind = "GCPNetworkListNil"

// GCPNetworkList is a list of values of the 'GCP_network' type.
type GCPNetworkList struct {
	href  string
	link  bool
	items []*GCPNetwork
}

// Len returns the length of the list.
func (l *GCPNetworkList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPNetworkList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPNetworkList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPNetworkList) SetItems(items []*GCPNetwork) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPNetworkList) Items() []*GCPNetwork {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPNetworkList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPNetworkList) Get(i int) *GCPNetwork {
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
func (l *GCPNetworkList) Slice() []*GCPNetwork {
	var slice []*GCPNetwork
	if l == nil {
		slice = make([]*GCPNetwork, 0)
	} else {
		slice = make([]*GCPNetwork, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPNetworkList) Each(f func(item *GCPNetwork) bool) {
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
func (l *GCPNetworkList) Range(f func(index int, item *GCPNetwork) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
