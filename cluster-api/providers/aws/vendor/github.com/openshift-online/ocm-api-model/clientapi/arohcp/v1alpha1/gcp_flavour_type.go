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

// GCPFlavour represents the values of the 'GCP_flavour' type.
//
// Specification for different classes of nodes inside a flavour.
type GCPFlavour struct {
	fieldSet_           []bool
	computeInstanceType string
	infraInstanceType   string
	infraVolume         *GCPVolume
	masterInstanceType  string
	masterVolume        *GCPVolume
	workerVolume        *GCPVolume
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPFlavour) Empty() bool {
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

// ComputeInstanceType returns the value of the 'compute_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP default instance type for the worker volume.
//
// User can be overridden specifying in the cluster itself a type for compute node.
func (o *GCPFlavour) ComputeInstanceType() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.computeInstanceType
	}
	return ""
}

// GetComputeInstanceType returns the value of the 'compute_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// GCP default instance type for the worker volume.
//
// User can be overridden specifying in the cluster itself a type for compute node.
func (o *GCPFlavour) GetComputeInstanceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.computeInstanceType
	}
	return
}

// InfraInstanceType returns the value of the 'infra_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP default instance type for the infra volume.
func (o *GCPFlavour) InfraInstanceType() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.infraInstanceType
	}
	return ""
}

// GetInfraInstanceType returns the value of the 'infra_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// GCP default instance type for the infra volume.
func (o *GCPFlavour) GetInfraInstanceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.infraInstanceType
	}
	return
}

// InfraVolume returns the value of the 'infra_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Infra volume specification.
func (o *GCPFlavour) InfraVolume() *GCPVolume {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.infraVolume
	}
	return nil
}

// GetInfraVolume returns the value of the 'infra_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Infra volume specification.
func (o *GCPFlavour) GetInfraVolume() (value *GCPVolume, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.infraVolume
	}
	return
}

// MasterInstanceType returns the value of the 'master_instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP default instance type for the master volume.
func (o *GCPFlavour) MasterInstanceType() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.masterInstanceType
	}
	return ""
}

// GetMasterInstanceType returns the value of the 'master_instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// GCP default instance type for the master volume.
func (o *GCPFlavour) GetMasterInstanceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.masterInstanceType
	}
	return
}

// MasterVolume returns the value of the 'master_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Master volume specification.
func (o *GCPFlavour) MasterVolume() *GCPVolume {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.masterVolume
	}
	return nil
}

// GetMasterVolume returns the value of the 'master_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Master volume specification.
func (o *GCPFlavour) GetMasterVolume() (value *GCPVolume, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.masterVolume
	}
	return
}

// WorkerVolume returns the value of the 'worker_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Worker volume specification.
func (o *GCPFlavour) WorkerVolume() *GCPVolume {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.workerVolume
	}
	return nil
}

// GetWorkerVolume returns the value of the 'worker_volume' attribute and
// a flag indicating if the attribute has a value.
//
// Worker volume specification.
func (o *GCPFlavour) GetWorkerVolume() (value *GCPVolume, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.workerVolume
	}
	return
}

// GCPFlavourListKind is the name of the type used to represent list of objects of
// type 'GCP_flavour'.
const GCPFlavourListKind = "GCPFlavourList"

// GCPFlavourListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_flavour'.
const GCPFlavourListLinkKind = "GCPFlavourListLink"

// GCPFlavourNilKind is the name of the type used to nil lists of objects of
// type 'GCP_flavour'.
const GCPFlavourListNilKind = "GCPFlavourListNil"

// GCPFlavourList is a list of values of the 'GCP_flavour' type.
type GCPFlavourList struct {
	href  string
	link  bool
	items []*GCPFlavour
}

// Len returns the length of the list.
func (l *GCPFlavourList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPFlavourList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPFlavourList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPFlavourList) SetItems(items []*GCPFlavour) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPFlavourList) Items() []*GCPFlavour {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPFlavourList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPFlavourList) Get(i int) *GCPFlavour {
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
func (l *GCPFlavourList) Slice() []*GCPFlavour {
	var slice []*GCPFlavour
	if l == nil {
		slice = make([]*GCPFlavour, 0)
	} else {
		slice = make([]*GCPFlavour, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPFlavourList) Each(f func(item *GCPFlavour) bool) {
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
func (l *GCPFlavourList) Range(f func(index int, item *GCPFlavour) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
