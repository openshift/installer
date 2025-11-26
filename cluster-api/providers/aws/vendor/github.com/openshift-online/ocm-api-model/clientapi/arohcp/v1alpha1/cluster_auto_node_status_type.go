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

// ClusterAutoNodeStatus represents the values of the 'cluster_auto_node_status' type.
//
// Additional status information on the AutoNode configuration on this Cluster
type ClusterAutoNodeStatus struct {
	fieldSet_ []bool
	message   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterAutoNodeStatus) Empty() bool {
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

// Message returns the value of the 'message' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Messages relating to the status of the AutoNode installation on this Cluster
func (o *ClusterAutoNodeStatus) Message() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.message
	}
	return ""
}

// GetMessage returns the value of the 'message' attribute and
// a flag indicating if the attribute has a value.
//
// Messages relating to the status of the AutoNode installation on this Cluster
func (o *ClusterAutoNodeStatus) GetMessage() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.message
	}
	return
}

// ClusterAutoNodeStatusListKind is the name of the type used to represent list of objects of
// type 'cluster_auto_node_status'.
const ClusterAutoNodeStatusListKind = "ClusterAutoNodeStatusList"

// ClusterAutoNodeStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_auto_node_status'.
const ClusterAutoNodeStatusListLinkKind = "ClusterAutoNodeStatusListLink"

// ClusterAutoNodeStatusNilKind is the name of the type used to nil lists of objects of
// type 'cluster_auto_node_status'.
const ClusterAutoNodeStatusListNilKind = "ClusterAutoNodeStatusListNil"

// ClusterAutoNodeStatusList is a list of values of the 'cluster_auto_node_status' type.
type ClusterAutoNodeStatusList struct {
	href  string
	link  bool
	items []*ClusterAutoNodeStatus
}

// Len returns the length of the list.
func (l *ClusterAutoNodeStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterAutoNodeStatusList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterAutoNodeStatusList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterAutoNodeStatusList) SetItems(items []*ClusterAutoNodeStatus) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterAutoNodeStatusList) Items() []*ClusterAutoNodeStatus {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterAutoNodeStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterAutoNodeStatusList) Get(i int) *ClusterAutoNodeStatus {
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
func (l *ClusterAutoNodeStatusList) Slice() []*ClusterAutoNodeStatus {
	var slice []*ClusterAutoNodeStatus
	if l == nil {
		slice = make([]*ClusterAutoNodeStatus, 0)
	} else {
		slice = make([]*ClusterAutoNodeStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterAutoNodeStatusList) Each(f func(item *ClusterAutoNodeStatus) bool) {
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
func (l *ClusterAutoNodeStatusList) Range(f func(index int, item *ClusterAutoNodeStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
