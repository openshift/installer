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

// AddOnConfigKind is the name of the type used to represent objects
// of type 'add_on_config'.
const AddOnConfigKind = "AddOnConfig"

// AddOnConfigLinkKind is the name of the type used to represent links
// to objects of type 'add_on_config'.
const AddOnConfigLinkKind = "AddOnConfigLink"

// AddOnConfigNilKind is the name of the type used to nil references
// to objects of type 'add_on_config'.
const AddOnConfigNilKind = "AddOnConfigNil"

// AddOnConfig represents the values of the 'add_on_config' type.
//
// Representation of an add-on config.
// The attributes under it are to be used by the addon once its installed in the cluster.
type AddOnConfig struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	addOnEnvironmentVariables []*AddOnEnvironmentVariable
	secretPropagations        []*AddOnSecretPropagation
}

// Kind returns the name of the type of the object.
func (o *AddOnConfig) Kind() string {
	if o == nil {
		return AddOnConfigNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AddOnConfigLinkKind
	}
	return AddOnConfigKind
}

// Link returns true if this is a link.
func (o *AddOnConfig) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AddOnConfig) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOnConfig) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOnConfig) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOnConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnConfig) Empty() bool {
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

// AddOnEnvironmentVariables returns the value of the 'add_on_environment_variables' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of environment variables for the addon
func (o *AddOnConfig) AddOnEnvironmentVariables() []*AddOnEnvironmentVariable {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.addOnEnvironmentVariables
	}
	return nil
}

// GetAddOnEnvironmentVariables returns the value of the 'add_on_environment_variables' attribute and
// a flag indicating if the attribute has a value.
//
// List of environment variables for the addon
func (o *AddOnConfig) GetAddOnEnvironmentVariables() (value []*AddOnEnvironmentVariable, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.addOnEnvironmentVariables
	}
	return
}

// SecretPropagations returns the value of the 'secret_propagations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of secret propagations for the addon
func (o *AddOnConfig) SecretPropagations() []*AddOnSecretPropagation {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.secretPropagations
	}
	return nil
}

// GetSecretPropagations returns the value of the 'secret_propagations' attribute and
// a flag indicating if the attribute has a value.
//
// List of secret propagations for the addon
func (o *AddOnConfig) GetSecretPropagations() (value []*AddOnSecretPropagation, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.secretPropagations
	}
	return
}

// AddOnConfigListKind is the name of the type used to represent list of objects of
// type 'add_on_config'.
const AddOnConfigListKind = "AddOnConfigList"

// AddOnConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_config'.
const AddOnConfigListLinkKind = "AddOnConfigListLink"

// AddOnConfigNilKind is the name of the type used to nil lists of objects of
// type 'add_on_config'.
const AddOnConfigListNilKind = "AddOnConfigListNil"

// AddOnConfigList is a list of values of the 'add_on_config' type.
type AddOnConfigList struct {
	href  string
	link  bool
	items []*AddOnConfig
}

// Kind returns the name of the type of the object.
func (l *AddOnConfigList) Kind() string {
	if l == nil {
		return AddOnConfigListNilKind
	}
	if l.link {
		return AddOnConfigListLinkKind
	}
	return AddOnConfigListKind
}

// Link returns true iif this is a link.
func (l *AddOnConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnConfigList) SetItems(items []*AddOnConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnConfigList) Items() []*AddOnConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnConfigList) Get(i int) *AddOnConfig {
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
func (l *AddOnConfigList) Slice() []*AddOnConfig {
	var slice []*AddOnConfig
	if l == nil {
		slice = make([]*AddOnConfig, 0)
	} else {
		slice = make([]*AddOnConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnConfigList) Each(f func(item *AddOnConfig) bool) {
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
func (l *AddOnConfigList) Range(f func(index int, item *AddOnConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
