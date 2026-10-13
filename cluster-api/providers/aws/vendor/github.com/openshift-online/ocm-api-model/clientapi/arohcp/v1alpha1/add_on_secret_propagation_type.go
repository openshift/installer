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

// AddOnSecretPropagation represents the values of the 'add_on_secret_propagation' type.
//
// Representation of an addon secret propagation
type AddOnSecretPropagation struct {
	fieldSet_         []bool
	id                string
	destinationSecret string
	sourceSecret      string
	enabled           bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnSecretPropagation) Empty() bool {
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

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the secret propagation
func (o *AddOnSecretPropagation) ID() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the secret propagation
func (o *AddOnSecretPropagation) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// DestinationSecret returns the value of the 'destination_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DestinationSecret is location of the secret to be added
func (o *AddOnSecretPropagation) DestinationSecret() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.destinationSecret
	}
	return ""
}

// GetDestinationSecret returns the value of the 'destination_secret' attribute and
// a flag indicating if the attribute has a value.
//
// DestinationSecret is location of the secret to be added
func (o *AddOnSecretPropagation) GetDestinationSecret() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.destinationSecret
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates is this secret propagation is enabled for the addon
func (o *AddOnSecretPropagation) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates is this secret propagation is enabled for the addon
func (o *AddOnSecretPropagation) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.enabled
	}
	return
}

// SourceSecret returns the value of the 'source_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// SourceSecret is location of the source secret
func (o *AddOnSecretPropagation) SourceSecret() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.sourceSecret
	}
	return ""
}

// GetSourceSecret returns the value of the 'source_secret' attribute and
// a flag indicating if the attribute has a value.
//
// SourceSecret is location of the source secret
func (o *AddOnSecretPropagation) GetSourceSecret() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.sourceSecret
	}
	return
}

// AddOnSecretPropagationListKind is the name of the type used to represent list of objects of
// type 'add_on_secret_propagation'.
const AddOnSecretPropagationListKind = "AddOnSecretPropagationList"

// AddOnSecretPropagationListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_secret_propagation'.
const AddOnSecretPropagationListLinkKind = "AddOnSecretPropagationListLink"

// AddOnSecretPropagationNilKind is the name of the type used to nil lists of objects of
// type 'add_on_secret_propagation'.
const AddOnSecretPropagationListNilKind = "AddOnSecretPropagationListNil"

// AddOnSecretPropagationList is a list of values of the 'add_on_secret_propagation' type.
type AddOnSecretPropagationList struct {
	href  string
	link  bool
	items []*AddOnSecretPropagation
}

// Len returns the length of the list.
func (l *AddOnSecretPropagationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnSecretPropagationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnSecretPropagationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnSecretPropagationList) SetItems(items []*AddOnSecretPropagation) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnSecretPropagationList) Items() []*AddOnSecretPropagation {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnSecretPropagationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnSecretPropagationList) Get(i int) *AddOnSecretPropagation {
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
func (l *AddOnSecretPropagationList) Slice() []*AddOnSecretPropagation {
	var slice []*AddOnSecretPropagation
	if l == nil {
		slice = make([]*AddOnSecretPropagation, 0)
	} else {
		slice = make([]*AddOnSecretPropagation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnSecretPropagationList) Each(f func(item *AddOnSecretPropagation) bool) {
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
func (l *AddOnSecretPropagationList) Range(f func(index int, item *AddOnSecretPropagation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
