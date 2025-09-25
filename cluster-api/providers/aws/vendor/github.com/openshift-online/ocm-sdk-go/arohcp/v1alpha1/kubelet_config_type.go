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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// KubeletConfigKind is the name of the type used to represent objects
// of type 'kubelet_config'.
const KubeletConfigKind = "KubeletConfig"

// KubeletConfigLinkKind is the name of the type used to represent links
// to objects of type 'kubelet_config'.
const KubeletConfigLinkKind = "KubeletConfigLink"

// KubeletConfigNilKind is the name of the type used to nil references
// to objects of type 'kubelet_config'.
const KubeletConfigNilKind = "KubeletConfigNil"

// KubeletConfig represents the values of the 'kubelet_config' type.
//
// OCM representation of KubeletConfig, exposing the fields of Kubernetes
// KubeletConfig that can be managed by users
type KubeletConfig struct {
	bitmap_      uint32
	id           string
	href         string
	name         string
	podPidsLimit int
}

// Kind returns the name of the type of the object.
func (o *KubeletConfig) Kind() string {
	if o == nil {
		return KubeletConfigNilKind
	}
	if o.bitmap_&1 != 0 {
		return KubeletConfigLinkKind
	}
	return KubeletConfigKind
}

// Link returns true if this is a link.
func (o *KubeletConfig) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *KubeletConfig) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *KubeletConfig) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *KubeletConfig) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *KubeletConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *KubeletConfig) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Allows the user to specify the name to be used to identify this KubeletConfig.
// Optional. A name will be generated if not provided.
func (o *KubeletConfig) Name() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Allows the user to specify the name to be used to identify this KubeletConfig.
// Optional. A name will be generated if not provided.
func (o *KubeletConfig) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.name
	}
	return
}

// PodPidsLimit returns the value of the 'pod_pids_limit' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Allows the user to specify the podPidsLimit to be applied via KubeletConfig.
// Useful if workloads have greater PIDs limit requirements than the OCP default.
func (o *KubeletConfig) PodPidsLimit() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.podPidsLimit
	}
	return 0
}

// GetPodPidsLimit returns the value of the 'pod_pids_limit' attribute and
// a flag indicating if the attribute has a value.
//
// Allows the user to specify the podPidsLimit to be applied via KubeletConfig.
// Useful if workloads have greater PIDs limit requirements than the OCP default.
func (o *KubeletConfig) GetPodPidsLimit() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.podPidsLimit
	}
	return
}

// KubeletConfigListKind is the name of the type used to represent list of objects of
// type 'kubelet_config'.
const KubeletConfigListKind = "KubeletConfigList"

// KubeletConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'kubelet_config'.
const KubeletConfigListLinkKind = "KubeletConfigListLink"

// KubeletConfigNilKind is the name of the type used to nil lists of objects of
// type 'kubelet_config'.
const KubeletConfigListNilKind = "KubeletConfigListNil"

// KubeletConfigList is a list of values of the 'kubelet_config' type.
type KubeletConfigList struct {
	href  string
	link  bool
	items []*KubeletConfig
}

// Kind returns the name of the type of the object.
func (l *KubeletConfigList) Kind() string {
	if l == nil {
		return KubeletConfigListNilKind
	}
	if l.link {
		return KubeletConfigListLinkKind
	}
	return KubeletConfigListKind
}

// Link returns true iif this is a link.
func (l *KubeletConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *KubeletConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *KubeletConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *KubeletConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *KubeletConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *KubeletConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *KubeletConfigList) SetItems(items []*KubeletConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *KubeletConfigList) Items() []*KubeletConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *KubeletConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *KubeletConfigList) Get(i int) *KubeletConfig {
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
func (l *KubeletConfigList) Slice() []*KubeletConfig {
	var slice []*KubeletConfig
	if l == nil {
		slice = make([]*KubeletConfig, 0)
	} else {
		slice = make([]*KubeletConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *KubeletConfigList) Each(f func(item *KubeletConfig) bool) {
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
func (l *KubeletConfigList) Range(f func(index int, item *KubeletConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
