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

// FlavourKind is the name of the type used to represent objects
// of type 'flavour'.
const FlavourKind = "Flavour"

// FlavourLinkKind is the name of the type used to represent links
// to objects of type 'flavour'.
const FlavourLinkKind = "FlavourLink"

// FlavourNilKind is the name of the type used to nil references
// to objects of type 'flavour'.
const FlavourNilKind = "FlavourNil"

// Flavour represents the values of the 'flavour' type.
//
// Set of predefined properties of a cluster. For example, a _huge_ flavour can be a cluster
// with 10 infra nodes and 1000 compute nodes.
type Flavour struct {
	fieldSet_ []bool
	id        string
	href      string
	aws       *AWSFlavour
	gcp       *GCPFlavour
	name      string
	network   *Network
	nodes     *FlavourNodes
}

// Kind returns the name of the type of the object.
func (o *Flavour) Kind() string {
	if o == nil {
		return FlavourNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return FlavourLinkKind
	}
	return FlavourKind
}

// Link returns true if this is a link.
func (o *Flavour) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Flavour) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Flavour) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Flavour) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Flavour) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Flavour) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Default _Amazon Web Services_ settings of the cluster.
func (o *Flavour) AWS() *AWSFlavour {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
//
// Default _Amazon Web Services_ settings of the cluster.
func (o *Flavour) GetAWS() (value *AWSFlavour, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.aws
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Default _Google Cloud Platform_ settings of the cluster.
func (o *Flavour) GCP() *GCPFlavour {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
//
// Default _Google Cloud Platform_ settings of the cluster.
func (o *Flavour) GetGCP() (value *GCPFlavour, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.gcp
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Human friendly identifier of the cluster, for example `4`.
//
// NOTE: Currently for all flavours the `id` and `name` attributes have exactly the
// same values.
func (o *Flavour) Name() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Human friendly identifier of the cluster, for example `4`.
//
// NOTE: Currently for all flavours the `id` and `name` attributes have exactly the
// same values.
func (o *Flavour) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.name
	}
	return
}

// Network returns the value of the 'network' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Default network settings of the cluster.
//
// These can be overridden specifying in the cluster itself a different set of settings.
func (o *Flavour) Network() *Network {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.network
	}
	return nil
}

// GetNetwork returns the value of the 'network' attribute and
// a flag indicating if the attribute has a value.
//
// Default network settings of the cluster.
//
// These can be overridden specifying in the cluster itself a different set of settings.
func (o *Flavour) GetNetwork() (value *Network, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.network
	}
	return
}

// Nodes returns the value of the 'nodes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of nodes that will be used by default when creating a cluster that uses
// this flavour.
//
// These can be overridden specifying in the cluster itself a different number of nodes.
func (o *Flavour) Nodes() *FlavourNodes {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.nodes
	}
	return nil
}

// GetNodes returns the value of the 'nodes' attribute and
// a flag indicating if the attribute has a value.
//
// Number of nodes that will be used by default when creating a cluster that uses
// this flavour.
//
// These can be overridden specifying in the cluster itself a different number of nodes.
func (o *Flavour) GetNodes() (value *FlavourNodes, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.nodes
	}
	return
}

// FlavourListKind is the name of the type used to represent list of objects of
// type 'flavour'.
const FlavourListKind = "FlavourList"

// FlavourListLinkKind is the name of the type used to represent links to list
// of objects of type 'flavour'.
const FlavourListLinkKind = "FlavourListLink"

// FlavourNilKind is the name of the type used to nil lists of objects of
// type 'flavour'.
const FlavourListNilKind = "FlavourListNil"

// FlavourList is a list of values of the 'flavour' type.
type FlavourList struct {
	href  string
	link  bool
	items []*Flavour
}

// Kind returns the name of the type of the object.
func (l *FlavourList) Kind() string {
	if l == nil {
		return FlavourListNilKind
	}
	if l.link {
		return FlavourListLinkKind
	}
	return FlavourListKind
}

// Link returns true iif this is a link.
func (l *FlavourList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *FlavourList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *FlavourList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *FlavourList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *FlavourList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *FlavourList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *FlavourList) SetItems(items []*Flavour) {
	l.items = items
}

// Items returns the items of the list.
func (l *FlavourList) Items() []*Flavour {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *FlavourList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FlavourList) Get(i int) *Flavour {
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
func (l *FlavourList) Slice() []*Flavour {
	var slice []*Flavour
	if l == nil {
		slice = make([]*Flavour, 0)
	} else {
		slice = make([]*Flavour, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FlavourList) Each(f func(item *Flavour) bool) {
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
func (l *FlavourList) Range(f func(index int, item *Flavour) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
