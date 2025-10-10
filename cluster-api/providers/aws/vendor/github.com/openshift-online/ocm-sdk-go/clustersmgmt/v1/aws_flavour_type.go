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

// AWSFlavour represents the values of the 'AWS_flavour' type.
//
// Specification for different classes of nodes inside a flavour.
type AWSFlavour struct {
	bitmap_             uint32
	computeInstanceType string
	infraInstanceType   string
	infraVolume         *AWSVolume
	masterInstanceType  string
	masterVolume        *AWSVolume
	workerVolume        *AWSVolume
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSFlavour) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ComputeInstanceType returns the value of the 'compute_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS default instance type for the worker volume.
//
// User can be overridden specifying in the cluster itself a type for compute node.
func (o *AWSFlavour) ComputeInstanceType() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.computeInstanceType
	}
	return ""
}

// GetComputeInstanceType returns the value of the 'compute_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// AWS default instance type for the worker volume.
//
// User can be overridden specifying in the cluster itself a type for compute node.
func (o *AWSFlavour) GetComputeInstanceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.computeInstanceType
	}
	return
}

// InfraInstanceType returns the value of the 'infra_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS default instance type for the infra volume.
func (o *AWSFlavour) InfraInstanceType() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.infraInstanceType
	}
	return ""
}

// GetInfraInstanceType returns the value of the 'infra_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// AWS default instance type for the infra volume.
func (o *AWSFlavour) GetInfraInstanceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.infraInstanceType
	}
	return
}

// InfraVolume returns the value of the 'infra_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Infra volume specification.
func (o *AWSFlavour) InfraVolume() *AWSVolume {
	if o != nil && o.bitmap_&4 != 0 {
		return o.infraVolume
	}
	return nil
}

// GetInfraVolume returns the value of the 'infra_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Infra volume specification.
func (o *AWSFlavour) GetInfraVolume() (value *AWSVolume, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.infraVolume
	}
	return
}

// MasterInstanceType returns the value of the 'master_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS default instance type for the master volume.
func (o *AWSFlavour) MasterInstanceType() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.masterInstanceType
	}
	return ""
}

// GetMasterInstanceType returns the value of the 'master_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// AWS default instance type for the master volume.
func (o *AWSFlavour) GetMasterInstanceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.masterInstanceType
	}
	return
}

// MasterVolume returns the value of the 'master_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Master volume specification.
func (o *AWSFlavour) MasterVolume() *AWSVolume {
	if o != nil && o.bitmap_&16 != 0 {
		return o.masterVolume
	}
	return nil
}

// GetMasterVolume returns the value of the 'master_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Master volume specification.
func (o *AWSFlavour) GetMasterVolume() (value *AWSVolume, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.masterVolume
	}
	return
}

// WorkerVolume returns the value of the 'worker_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Worker volume specification.
func (o *AWSFlavour) WorkerVolume() *AWSVolume {
	if o != nil && o.bitmap_&32 != 0 {
		return o.workerVolume
	}
	return nil
}

// GetWorkerVolume returns the value of the 'worker_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Worker volume specification.
func (o *AWSFlavour) GetWorkerVolume() (value *AWSVolume, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.workerVolume
	}
	return
}

// AWSFlavourListKind is the name of the type used to represent list of objects of
// type 'AWS_flavour'.
const AWSFlavourListKind = "AWSFlavourList"

// AWSFlavourListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_flavour'.
const AWSFlavourListLinkKind = "AWSFlavourListLink"

// AWSFlavourNilKind is the name of the type used to nil lists of objects of
// type 'AWS_flavour'.
const AWSFlavourListNilKind = "AWSFlavourListNil"

// AWSFlavourList is a list of values of the 'AWS_flavour' type.
type AWSFlavourList struct {
	href  string
	link  bool
	items []*AWSFlavour
}

// Len returns the length of the list.
func (l *AWSFlavourList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSFlavourList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSFlavourList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSFlavourList) SetItems(items []*AWSFlavour) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSFlavourList) Items() []*AWSFlavour {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSFlavourList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSFlavourList) Get(i int) *AWSFlavour {
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
func (l *AWSFlavourList) Slice() []*AWSFlavour {
	var slice []*AWSFlavour
	if l == nil {
		slice = make([]*AWSFlavour, 0)
	} else {
		slice = make([]*AWSFlavour, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSFlavourList) Each(f func(item *AWSFlavour) bool) {
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
func (l *AWSFlavourList) Range(f func(index int, item *AWSFlavour) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
