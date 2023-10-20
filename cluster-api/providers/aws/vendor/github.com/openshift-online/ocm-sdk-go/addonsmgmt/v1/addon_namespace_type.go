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

// AddonNamespace represents the values of the 'addon_namespace' type.
//
// Representation of an addon namespace object.
type AddonNamespace struct {
	bitmap_     uint32
	annotations map[string]string
	labels      map[string]string
	name        string
	enabled     bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonNamespace) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Annotations returns the value of the 'annotations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Annotations to be included in the addon namespace
func (o *AddonNamespace) Annotations() map[string]string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.annotations
	}
	return nil
}

// GetAnnotations returns the value of the 'annotations' attribute and
// a flag indicating if the attribute has a value.
//
// Annotations to be included in the addon namespace
func (o *AddonNamespace) GetAnnotations() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.annotations
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Enabled shows if this namespace object is in use
func (o *AddonNamespace) Enabled() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Enabled shows if this namespace object is in use
func (o *AddonNamespace) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.enabled
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Labels to be included in the addon namespace
func (o *AddonNamespace) Labels() map[string]string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// Labels to be included in the addon namespace
func (o *AddonNamespace) GetLabels() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.labels
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the namespace
func (o *AddonNamespace) Name() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the namespace
func (o *AddonNamespace) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.name
	}
	return
}

// AddonNamespaceListKind is the name of the type used to represent list of objects of
// type 'addon_namespace'.
const AddonNamespaceListKind = "AddonNamespaceList"

// AddonNamespaceListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_namespace'.
const AddonNamespaceListLinkKind = "AddonNamespaceListLink"

// AddonNamespaceNilKind is the name of the type used to nil lists of objects of
// type 'addon_namespace'.
const AddonNamespaceListNilKind = "AddonNamespaceListNil"

// AddonNamespaceList is a list of values of the 'addon_namespace' type.
type AddonNamespaceList struct {
	href  string
	link  bool
	items []*AddonNamespace
}

// Len returns the length of the list.
func (l *AddonNamespaceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonNamespaceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonNamespaceList) Get(i int) *AddonNamespace {
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
func (l *AddonNamespaceList) Slice() []*AddonNamespace {
	var slice []*AddonNamespace
	if l == nil {
		slice = make([]*AddonNamespace, 0)
	} else {
		slice = make([]*AddonNamespace, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonNamespaceList) Each(f func(item *AddonNamespace) bool) {
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
func (l *AddonNamespaceList) Range(f func(index int, item *AddonNamespace) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
