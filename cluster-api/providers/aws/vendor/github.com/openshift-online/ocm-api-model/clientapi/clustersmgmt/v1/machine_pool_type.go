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

// MachinePoolKind is the name of the type used to represent objects
// of type 'machine_pool'.
const MachinePoolKind = "MachinePool"

// MachinePoolLinkKind is the name of the type used to represent links
// to objects of type 'machine_pool'.
const MachinePoolLinkKind = "MachinePoolLink"

// MachinePoolNilKind is the name of the type used to nil references
// to objects of type 'machine_pool'.
const MachinePoolNilKind = "MachinePoolNil"

// MachinePool represents the values of the 'machine_pool' type.
//
// Representation of a machine pool in a cluster.
type MachinePool struct {
	fieldSet_            []bool
	id                   string
	href                 string
	aws                  *AWSMachinePool
	gcp                  *GCPMachinePool
	autoscaling          *MachinePoolAutoscaling
	availabilityZones    []string
	instanceType         string
	labels               map[string]string
	replicas             int
	rootVolume           *RootVolume
	securityGroupFilters []*MachinePoolSecurityGroupFilter
	subnets              []string
	taints               []*Taint
}

// Kind returns the name of the type of the object.
func (o *MachinePool) Kind() string {
	if o == nil {
		return MachinePoolNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return MachinePoolLinkKind
	}
	return MachinePoolKind
}

// Link returns true if this is a link.
func (o *MachinePool) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *MachinePool) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *MachinePool) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *MachinePool) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *MachinePool) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MachinePool) Empty() bool {
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
// AWS specific parameters (Optional).
func (o *MachinePool) AWS() *AWSMachinePool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
//
// AWS specific parameters (Optional).
func (o *MachinePool) GetAWS() (value *AWSMachinePool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.aws
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP specific parameters (Optional).
func (o *MachinePool) GCP() *GCPMachinePool {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
//
// GCP specific parameters (Optional).
func (o *MachinePool) GetGCP() (value *GCPMachinePool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.gcp
	}
	return
}

// Autoscaling returns the value of the 'autoscaling' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details for auto-scaling the machine pool.
// Replicas and autoscaling cannot be used together.
func (o *MachinePool) Autoscaling() *MachinePoolAutoscaling {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.autoscaling
	}
	return nil
}

// GetAutoscaling returns the value of the 'autoscaling' attribute and
// a flag indicating if the attribute has a value.
//
// Details for auto-scaling the machine pool.
// Replicas and autoscaling cannot be used together.
func (o *MachinePool) GetAutoscaling() (value *MachinePoolAutoscaling, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.autoscaling
	}
	return
}

// AvailabilityZones returns the value of the 'availability_zones' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The availability zones upon which the nodes are created.
func (o *MachinePool) AvailabilityZones() []string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.availabilityZones
	}
	return nil
}

// GetAvailabilityZones returns the value of the 'availability_zones' attribute and
// a flag indicating if the attribute has a value.
//
// The availability zones upon which the nodes are created.
func (o *MachinePool) GetAvailabilityZones() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.availabilityZones
	}
	return
}

// InstanceType returns the value of the 'instance_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The instance type of Nodes to create.
func (o *MachinePool) InstanceType() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.instanceType
	}
	return ""
}

// GetInstanceType returns the value of the 'instance_type' attribute and
// a flag indicating if the attribute has a value.
//
// The instance type of Nodes to create.
func (o *MachinePool) GetInstanceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.instanceType
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The labels set on the Nodes created.
func (o *MachinePool) Labels() map[string]string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// The labels set on the Nodes created.
func (o *MachinePool) GetLabels() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.labels
	}
	return
}

// Replicas returns the value of the 'replicas' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The number of Machines (and Nodes) to create.
// Replicas and autoscaling cannot be used together.
func (o *MachinePool) Replicas() int {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.replicas
	}
	return 0
}

// GetReplicas returns the value of the 'replicas' attribute and
// a flag indicating if the attribute has a value.
//
// The number of Machines (and Nodes) to create.
// Replicas and autoscaling cannot be used together.
func (o *MachinePool) GetReplicas() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.replicas
	}
	return
}

// RootVolume returns the value of the 'root_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The machine root volume capabilities.
func (o *MachinePool) RootVolume() *RootVolume {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.rootVolume
	}
	return nil
}

// GetRootVolume returns the value of the 'root_volume' attribute and
// a flag indicating if the attribute has a value.
//
// The machine root volume capabilities.
func (o *MachinePool) GetRootVolume() (value *RootVolume, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.rootVolume
	}
	return
}

// SecurityGroupFilters returns the value of the 'security_group_filters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of security groups to be applied to MachinePool (Optional)
func (o *MachinePool) SecurityGroupFilters() []*MachinePoolSecurityGroupFilter {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.securityGroupFilters
	}
	return nil
}

// GetSecurityGroupFilters returns the value of the 'security_group_filters' attribute and
// a flag indicating if the attribute has a value.
//
// List of security groups to be applied to MachinePool (Optional)
func (o *MachinePool) GetSecurityGroupFilters() (value []*MachinePoolSecurityGroupFilter, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.securityGroupFilters
	}
	return
}

// Subnets returns the value of the 'subnets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The subnets upon which the nodes are created.
func (o *MachinePool) Subnets() []string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.subnets
	}
	return nil
}

// GetSubnets returns the value of the 'subnets' attribute and
// a flag indicating if the attribute has a value.
//
// The subnets upon which the nodes are created.
func (o *MachinePool) GetSubnets() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.subnets
	}
	return
}

// Taints returns the value of the 'taints' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The taints set on the Nodes created.
func (o *MachinePool) Taints() []*Taint {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.taints
	}
	return nil
}

// GetTaints returns the value of the 'taints' attribute and
// a flag indicating if the attribute has a value.
//
// The taints set on the Nodes created.
func (o *MachinePool) GetTaints() (value []*Taint, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.taints
	}
	return
}

// MachinePoolListKind is the name of the type used to represent list of objects of
// type 'machine_pool'.
const MachinePoolListKind = "MachinePoolList"

// MachinePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'machine_pool'.
const MachinePoolListLinkKind = "MachinePoolListLink"

// MachinePoolNilKind is the name of the type used to nil lists of objects of
// type 'machine_pool'.
const MachinePoolListNilKind = "MachinePoolListNil"

// MachinePoolList is a list of values of the 'machine_pool' type.
type MachinePoolList struct {
	href  string
	link  bool
	items []*MachinePool
}

// Kind returns the name of the type of the object.
func (l *MachinePoolList) Kind() string {
	if l == nil {
		return MachinePoolListNilKind
	}
	if l.link {
		return MachinePoolListLinkKind
	}
	return MachinePoolListKind
}

// Link returns true iif this is a link.
func (l *MachinePoolList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *MachinePoolList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *MachinePoolList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *MachinePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *MachinePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *MachinePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *MachinePoolList) SetItems(items []*MachinePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *MachinePoolList) Items() []*MachinePool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *MachinePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MachinePoolList) Get(i int) *MachinePool {
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
func (l *MachinePoolList) Slice() []*MachinePool {
	var slice []*MachinePool
	if l == nil {
		slice = make([]*MachinePool, 0)
	} else {
		slice = make([]*MachinePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MachinePoolList) Each(f func(item *MachinePool) bool) {
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
func (l *MachinePoolList) Range(f func(index int, item *MachinePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
