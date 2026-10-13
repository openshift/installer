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

// AzureNodesOutboundConnectivity represents the values of the 'azure_nodes_outbound_connectivity' type.
//
// The configuration of the node outbound connectivity
type AzureNodesOutboundConnectivity struct {
	fieldSet_    []bool
	outboundType string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureNodesOutboundConnectivity) Empty() bool {
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

// OutboundType returns the value of the 'outbound_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// OutboundType is the type of network outbound configuration.
// The default and only accepted value is 'load_balancer'.
// This value is immutable.
func (o *AzureNodesOutboundConnectivity) OutboundType() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.outboundType
	}
	return ""
}

// GetOutboundType returns the value of the 'outbound_type' attribute and
// a flag indicating if the attribute has a value.
//
// OutboundType is the type of network outbound configuration.
// The default and only accepted value is 'load_balancer'.
// This value is immutable.
func (o *AzureNodesOutboundConnectivity) GetOutboundType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.outboundType
	}
	return
}

// AzureNodesOutboundConnectivityListKind is the name of the type used to represent list of objects of
// type 'azure_nodes_outbound_connectivity'.
const AzureNodesOutboundConnectivityListKind = "AzureNodesOutboundConnectivityList"

// AzureNodesOutboundConnectivityListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_nodes_outbound_connectivity'.
const AzureNodesOutboundConnectivityListLinkKind = "AzureNodesOutboundConnectivityListLink"

// AzureNodesOutboundConnectivityNilKind is the name of the type used to nil lists of objects of
// type 'azure_nodes_outbound_connectivity'.
const AzureNodesOutboundConnectivityListNilKind = "AzureNodesOutboundConnectivityListNil"

// AzureNodesOutboundConnectivityList is a list of values of the 'azure_nodes_outbound_connectivity' type.
type AzureNodesOutboundConnectivityList struct {
	href  string
	link  bool
	items []*AzureNodesOutboundConnectivity
}

// Len returns the length of the list.
func (l *AzureNodesOutboundConnectivityList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureNodesOutboundConnectivityList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureNodesOutboundConnectivityList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureNodesOutboundConnectivityList) SetItems(items []*AzureNodesOutboundConnectivity) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureNodesOutboundConnectivityList) Items() []*AzureNodesOutboundConnectivity {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureNodesOutboundConnectivityList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureNodesOutboundConnectivityList) Get(i int) *AzureNodesOutboundConnectivity {
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
func (l *AzureNodesOutboundConnectivityList) Slice() []*AzureNodesOutboundConnectivity {
	var slice []*AzureNodesOutboundConnectivity
	if l == nil {
		slice = make([]*AzureNodesOutboundConnectivity, 0)
	} else {
		slice = make([]*AzureNodesOutboundConnectivity, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureNodesOutboundConnectivityList) Each(f func(item *AzureNodesOutboundConnectivity) bool) {
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
func (l *AzureNodesOutboundConnectivityList) Range(f func(index int, item *AzureNodesOutboundConnectivity) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
