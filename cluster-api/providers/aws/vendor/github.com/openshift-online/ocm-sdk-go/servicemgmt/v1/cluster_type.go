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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// Cluster represents the values of the 'cluster' type.
//
// This represents the parameters needed by Managed Service to create a cluster.
type Cluster struct {
	bitmap_     uint32
	api         *ClusterAPI
	aws         *AWS
	displayName string
	href        string
	id          string
	name        string
	network     *Network
	nodes       *ClusterNodes
	properties  map[string]string
	region      *CloudRegion
	state       string
	multiAZ     bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Cluster) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// API returns the value of the 'API' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) API() *ClusterAPI {
	if o != nil && o.bitmap_&1 != 0 {
		return o.api
	}
	return nil
}

// GetAPI returns the value of the 'API' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetAPI() (value *ClusterAPI, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.api
	}
	return
}

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) AWS() *AWS {
	if o != nil && o.bitmap_&2 != 0 {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetAWS() (value *AWS, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.aws
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DisplayName is the name of the cluster for display purposes.
// It can contain spaces.
func (o *Cluster) DisplayName() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// DisplayName is the name of the cluster for display purposes.
// It can contain spaces.
func (o *Cluster) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.displayName
	}
	return
}

// Href returns the value of the 'href' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Href() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.href
	}
	return ""
}

// GetHref returns the value of the 'href' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetHref() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.href
	}
	return
}

// Id returns the value of the 'id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Id() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.id
	}
	return ""
}

// GetId returns the value of the 'id' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.id
	}
	return
}

// MultiAZ returns the value of the 'multi_AZ' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Flag indicating if the cluster should be created with nodes in
// different availability zones or all the nodes in a single one
// randomly selected.
func (o *Cluster) MultiAZ() bool {
	if o != nil && o.bitmap_&32 != 0 {
		return o.multiAZ
	}
	return false
}

// GetMultiAZ returns the value of the 'multi_AZ' attribute and
// a flag indicating if the attribute has a value.
//
// Flag indicating if the cluster should be created with nodes in
// different availability zones or all the nodes in a single one
// randomly selected.
func (o *Cluster) GetMultiAZ() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.multiAZ
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Name() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.name
	}
	return
}

// Network returns the value of the 'network' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Network() *Network {
	if o != nil && o.bitmap_&128 != 0 {
		return o.network
	}
	return nil
}

// GetNetwork returns the value of the 'network' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetNetwork() (value *Network, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.network
	}
	return
}

// Nodes returns the value of the 'nodes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Nodes() *ClusterNodes {
	if o != nil && o.bitmap_&256 != 0 {
		return o.nodes
	}
	return nil
}

// GetNodes returns the value of the 'nodes' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetNodes() (value *ClusterNodes, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.nodes
	}
	return
}

// Properties returns the value of the 'properties' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Properties() map[string]string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.properties
	}
	return nil
}

// GetProperties returns the value of the 'properties' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetProperties() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.properties
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) Region() *CloudRegion {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.region
	}
	return nil
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetRegion() (value *CloudRegion, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.region
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Cluster) State() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.state
	}
	return ""
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
func (o *Cluster) GetState() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.state
	}
	return
}

// ClusterListKind is the name of the type used to represent list of objects of
// type 'cluster'.
const ClusterListKind = "ClusterList"

// ClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster'.
const ClusterListLinkKind = "ClusterListLink"

// ClusterNilKind is the name of the type used to nil lists of objects of
// type 'cluster'.
const ClusterListNilKind = "ClusterListNil"

// ClusterList is a list of values of the 'cluster' type.
type ClusterList struct {
	href  string
	link  bool
	items []*Cluster
}

// Len returns the length of the list.
func (l *ClusterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterList) Get(i int) *Cluster {
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
func (l *ClusterList) Slice() []*Cluster {
	var slice []*Cluster
	if l == nil {
		slice = make([]*Cluster, 0)
	} else {
		slice = make([]*Cluster, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterList) Each(f func(item *Cluster) bool) {
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
func (l *ClusterList) Range(f func(index int, item *Cluster) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
