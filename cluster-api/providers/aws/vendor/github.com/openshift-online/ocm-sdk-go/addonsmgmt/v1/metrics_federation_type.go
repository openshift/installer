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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// MetricsFederation represents the values of the 'metrics_federation' type.
//
// Representation of Metrics Federation
type MetricsFederation struct {
	bitmap_     uint32
	matchLabels map[string]string
	matchNames  []string
	namespace   string
	portName    string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MetricsFederation) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// MatchLabels returns the value of the 'match_labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of labels used to discover the prometheus server(s) to be federated.
func (o *MetricsFederation) MatchLabels() map[string]string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.matchLabels
	}
	return nil
}

// GetMatchLabels returns the value of the 'match_labels' attribute and
// a flag indicating if the attribute has a value.
//
// List of labels used to discover the prometheus server(s) to be federated.
func (o *MetricsFederation) GetMatchLabels() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.matchLabels
	}
	return
}

// MatchNames returns the value of the 'match_names' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of series names to federate from the prometheus server.
func (o *MetricsFederation) MatchNames() []string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.matchNames
	}
	return nil
}

// GetMatchNames returns the value of the 'match_names' attribute and
// a flag indicating if the attribute has a value.
//
// List of series names to federate from the prometheus server.
func (o *MetricsFederation) GetMatchNames() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.matchNames
	}
	return
}

// Namespace returns the value of the 'namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespace where the prometheus server is running.
func (o *MetricsFederation) Namespace() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.namespace
	}
	return ""
}

// GetNamespace returns the value of the 'namespace' attribute and
// a flag indicating if the attribute has a value.
//
// Namespace where the prometheus server is running.
func (o *MetricsFederation) GetNamespace() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.namespace
	}
	return
}

// PortName returns the value of the 'port_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the name of the service port fronting the prometheus server.
func (o *MetricsFederation) PortName() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.portName
	}
	return ""
}

// GetPortName returns the value of the 'port_name' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the name of the service port fronting the prometheus server.
func (o *MetricsFederation) GetPortName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.portName
	}
	return
}

// MetricsFederationListKind is the name of the type used to represent list of objects of
// type 'metrics_federation'.
const MetricsFederationListKind = "MetricsFederationList"

// MetricsFederationListLinkKind is the name of the type used to represent links to list
// of objects of type 'metrics_federation'.
const MetricsFederationListLinkKind = "MetricsFederationListLink"

// MetricsFederationNilKind is the name of the type used to nil lists of objects of
// type 'metrics_federation'.
const MetricsFederationListNilKind = "MetricsFederationListNil"

// MetricsFederationList is a list of values of the 'metrics_federation' type.
type MetricsFederationList struct {
	href  string
	link  bool
	items []*MetricsFederation
}

// Len returns the length of the list.
func (l *MetricsFederationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *MetricsFederationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *MetricsFederationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *MetricsFederationList) SetItems(items []*MetricsFederation) {
	l.items = items
}

// Items returns the items of the list.
func (l *MetricsFederationList) Items() []*MetricsFederation {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *MetricsFederationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MetricsFederationList) Get(i int) *MetricsFederation {
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
func (l *MetricsFederationList) Slice() []*MetricsFederation {
	var slice []*MetricsFederation
	if l == nil {
		slice = make([]*MetricsFederation, 0)
	} else {
		slice = make([]*MetricsFederation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MetricsFederationList) Each(f func(item *MetricsFederation) bool) {
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
func (l *MetricsFederationList) Range(f func(index int, item *MetricsFederation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
