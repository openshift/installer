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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// MachinePoolAutoscalingKind is the name of the type used to represent objects
// of type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingKind = "MachinePoolAutoscaling"

// MachinePoolAutoscalingLinkKind is the name of the type used to represent links
// to objects of type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingLinkKind = "MachinePoolAutoscalingLink"

// MachinePoolAutoscalingNilKind is the name of the type used to nil references
// to objects of type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingNilKind = "MachinePoolAutoscalingNil"

// MachinePoolAutoscaling represents the values of the 'machine_pool_autoscaling' type.
//
// Representation of a autoscaling in a machine pool.
type MachinePoolAutoscaling struct {
	bitmap_     uint32
	id          string
	href        string
	maxReplicas int
	minReplicas int
}

// Kind returns the name of the type of the object.
func (o *MachinePoolAutoscaling) Kind() string {
	if o == nil {
		return MachinePoolAutoscalingNilKind
	}
	if o.bitmap_&1 != 0 {
		return MachinePoolAutoscalingLinkKind
	}
	return MachinePoolAutoscalingKind
}

// Link returns true iif this is a link.
func (o *MachinePoolAutoscaling) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *MachinePoolAutoscaling) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *MachinePoolAutoscaling) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *MachinePoolAutoscaling) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *MachinePoolAutoscaling) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MachinePoolAutoscaling) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// MaxReplicas returns the value of the 'max_replicas' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The maximum number of replicas for the machine pool.
func (o *MachinePoolAutoscaling) MaxReplicas() int {
	if o != nil && o.bitmap_&8 != 0 {
		return o.maxReplicas
	}
	return 0
}

// GetMaxReplicas returns the value of the 'max_replicas' attribute and
// a flag indicating if the attribute has a value.
//
// The maximum number of replicas for the machine pool.
func (o *MachinePoolAutoscaling) GetMaxReplicas() (value int, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.maxReplicas
	}
	return
}

// MinReplicas returns the value of the 'min_replicas' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The minimum number of replicas for the machine pool.
func (o *MachinePoolAutoscaling) MinReplicas() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.minReplicas
	}
	return 0
}

// GetMinReplicas returns the value of the 'min_replicas' attribute and
// a flag indicating if the attribute has a value.
//
// The minimum number of replicas for the machine pool.
func (o *MachinePoolAutoscaling) GetMinReplicas() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.minReplicas
	}
	return
}

// MachinePoolAutoscalingListKind is the name of the type used to represent list of objects of
// type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingListKind = "MachinePoolAutoscalingList"

// MachinePoolAutoscalingListLinkKind is the name of the type used to represent links to list
// of objects of type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingListLinkKind = "MachinePoolAutoscalingListLink"

// MachinePoolAutoscalingNilKind is the name of the type used to nil lists of objects of
// type 'machine_pool_autoscaling'.
const MachinePoolAutoscalingListNilKind = "MachinePoolAutoscalingListNil"

// MachinePoolAutoscalingList is a list of values of the 'machine_pool_autoscaling' type.
type MachinePoolAutoscalingList struct {
	href  string
	link  bool
	items []*MachinePoolAutoscaling
}

// Kind returns the name of the type of the object.
func (l *MachinePoolAutoscalingList) Kind() string {
	if l == nil {
		return MachinePoolAutoscalingListNilKind
	}
	if l.link {
		return MachinePoolAutoscalingListLinkKind
	}
	return MachinePoolAutoscalingListKind
}

// Link returns true iif this is a link.
func (l *MachinePoolAutoscalingList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *MachinePoolAutoscalingList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *MachinePoolAutoscalingList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *MachinePoolAutoscalingList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *MachinePoolAutoscalingList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MachinePoolAutoscalingList) Get(i int) *MachinePoolAutoscaling {
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
func (l *MachinePoolAutoscalingList) Slice() []*MachinePoolAutoscaling {
	var slice []*MachinePoolAutoscaling
	if l == nil {
		slice = make([]*MachinePoolAutoscaling, 0)
	} else {
		slice = make([]*MachinePoolAutoscaling, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MachinePoolAutoscalingList) Each(f func(item *MachinePoolAutoscaling) bool) {
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
func (l *MachinePoolAutoscalingList) Range(f func(index int, item *MachinePoolAutoscaling) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
