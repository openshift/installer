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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// AddonConfig represents the values of the 'addon_config' type.
//
// Representation of an addon config.
// The attributes under it are to be used by the addon once its installed in the cluster.
type AddonConfig struct {
	fieldSet_                 []bool
	addOnEnvironmentVariables []*AddonEnvironmentVariable
	addOnSecretPropagations   []*AddonSecretPropagation
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonConfig) Empty() bool {
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

// AddOnEnvironmentVariables returns the value of the 'add_on_environment_variables' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of environment variables for the addon
func (o *AddonConfig) AddOnEnvironmentVariables() []*AddonEnvironmentVariable {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.addOnEnvironmentVariables
	}
	return nil
}

// GetAddOnEnvironmentVariables returns the value of the 'add_on_environment_variables' attribute and
// a flag indicating if the attribute has a value.
//
// List of environment variables for the addon
func (o *AddonConfig) GetAddOnEnvironmentVariables() (value []*AddonEnvironmentVariable, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.addOnEnvironmentVariables
	}
	return
}

// AddOnSecretPropagations returns the value of the 'add_on_secret_propagations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of secret propagations for the addon
func (o *AddonConfig) AddOnSecretPropagations() []*AddonSecretPropagation {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.addOnSecretPropagations
	}
	return nil
}

// GetAddOnSecretPropagations returns the value of the 'add_on_secret_propagations' attribute and
// a flag indicating if the attribute has a value.
//
// List of secret propagations for the addon
func (o *AddonConfig) GetAddOnSecretPropagations() (value []*AddonSecretPropagation, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.addOnSecretPropagations
	}
	return
}

// AddonConfigListKind is the name of the type used to represent list of objects of
// type 'addon_config'.
const AddonConfigListKind = "AddonConfigList"

// AddonConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_config'.
const AddonConfigListLinkKind = "AddonConfigListLink"

// AddonConfigNilKind is the name of the type used to nil lists of objects of
// type 'addon_config'.
const AddonConfigListNilKind = "AddonConfigListNil"

// AddonConfigList is a list of values of the 'addon_config' type.
type AddonConfigList struct {
	href  string
	link  bool
	items []*AddonConfig
}

// Len returns the length of the list.
func (l *AddonConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonConfigList) SetItems(items []*AddonConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonConfigList) Items() []*AddonConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonConfigList) Get(i int) *AddonConfig {
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
func (l *AddonConfigList) Slice() []*AddonConfig {
	var slice []*AddonConfig
	if l == nil {
		slice = make([]*AddonConfig, 0)
	} else {
		slice = make([]*AddonConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonConfigList) Each(f func(item *AddonConfig) bool) {
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
func (l *AddonConfigList) Range(f func(index int, item *AddonConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
