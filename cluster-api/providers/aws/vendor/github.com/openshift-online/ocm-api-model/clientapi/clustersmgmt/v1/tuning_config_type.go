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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// TuningConfigKind is the name of the type used to represent objects
// of type 'tuning_config'.
const TuningConfigKind = "TuningConfig"

// TuningConfigLinkKind is the name of the type used to represent links
// to objects of type 'tuning_config'.
const TuningConfigLinkKind = "TuningConfigLink"

// TuningConfigNilKind is the name of the type used to nil references
// to objects of type 'tuning_config'.
const TuningConfigNilKind = "TuningConfigNil"

// TuningConfig represents the values of the 'tuning_config' type.
//
// Representation of a tuning config.
type TuningConfig struct {
	fieldSet_ []bool
	id        string
	href      string
	name      string
	spec      interface{}
}

// Kind returns the name of the type of the object.
func (o *TuningConfig) Kind() string {
	if o == nil {
		return TuningConfigNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return TuningConfigLinkKind
	}
	return TuningConfigKind
}

// Link returns true if this is a link.
func (o *TuningConfig) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *TuningConfig) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *TuningConfig) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *TuningConfig) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *TuningConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TuningConfig) Empty() bool {
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

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the tuning config.
func (o *TuningConfig) Name() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the tuning config.
func (o *TuningConfig) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.name
	}
	return
}

// Spec returns the value of the 'spec' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Spec of the tuning config.
func (o *TuningConfig) Spec() interface{} {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.spec
	}
	return nil
}

// GetSpec returns the value of the 'spec' attribute and
// a flag indicating if the attribute has a value.
//
// Spec of the tuning config.
func (o *TuningConfig) GetSpec() (value interface{}, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.spec
	}
	return
}

// TuningConfigListKind is the name of the type used to represent list of objects of
// type 'tuning_config'.
const TuningConfigListKind = "TuningConfigList"

// TuningConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'tuning_config'.
const TuningConfigListLinkKind = "TuningConfigListLink"

// TuningConfigNilKind is the name of the type used to nil lists of objects of
// type 'tuning_config'.
const TuningConfigListNilKind = "TuningConfigListNil"

// TuningConfigList is a list of values of the 'tuning_config' type.
type TuningConfigList struct {
	href  string
	link  bool
	items []*TuningConfig
}

// Kind returns the name of the type of the object.
func (l *TuningConfigList) Kind() string {
	if l == nil {
		return TuningConfigListNilKind
	}
	if l.link {
		return TuningConfigListLinkKind
	}
	return TuningConfigListKind
}

// Link returns true iif this is a link.
func (l *TuningConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *TuningConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *TuningConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *TuningConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TuningConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TuningConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TuningConfigList) SetItems(items []*TuningConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *TuningConfigList) Items() []*TuningConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TuningConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TuningConfigList) Get(i int) *TuningConfig {
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
func (l *TuningConfigList) Slice() []*TuningConfig {
	var slice []*TuningConfig
	if l == nil {
		slice = make([]*TuningConfig, 0)
	} else {
		slice = make([]*TuningConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TuningConfigList) Each(f func(item *TuningConfig) bool) {
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
func (l *TuningConfigList) Range(f func(index int, item *TuningConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
