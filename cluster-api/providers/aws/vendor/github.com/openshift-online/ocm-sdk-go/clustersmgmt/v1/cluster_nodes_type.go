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

// ClusterNodes represents the values of the 'cluster_nodes' type.
//
// Counts of different classes of nodes inside a cluster.
type ClusterNodes struct {
	bitmap_              uint32
	autoscaleCompute     *MachinePoolAutoscaling
	availabilityZones    []string
	compute              int
	computeLabels        map[string]string
	computeMachineType   *MachineType
	computeRootVolume    *RootVolume
	infra                int
	infraMachineType     *MachineType
	master               int
	masterMachineType    *MachineType
	securityGroupFilters []*MachinePoolSecurityGroupFilter
	total                int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterNodes) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AutoscaleCompute returns the value of the 'autoscale_compute' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details for auto-scaling the compute machine pool.
// Compute and AutoscaleCompute cannot be used together.
func (o *ClusterNodes) AutoscaleCompute() *MachinePoolAutoscaling {
	if o != nil && o.bitmap_&1 != 0 {
		return o.autoscaleCompute
	}
	return nil
}

// GetAutoscaleCompute returns the value of the 'autoscale_compute' attribute and
// a flag indicating if the attribute has a value.
//
// Details for auto-scaling the compute machine pool.
// Compute and AutoscaleCompute cannot be used together.
func (o *ClusterNodes) GetAutoscaleCompute() (value *MachinePoolAutoscaling, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.autoscaleCompute
	}
	return
}

// AvailabilityZones returns the value of the 'availability_zones' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The availability zones upon which the nodes are created.
func (o *ClusterNodes) AvailabilityZones() []string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.availabilityZones
	}
	return nil
}

// GetAvailabilityZones returns the value of the 'availability_zones' attribute and
// a flag indicating if the attribute has a value.
//
// The availability zones upon which the nodes are created.
func (o *ClusterNodes) GetAvailabilityZones() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.availabilityZones
	}
	return
}

// Compute returns the value of the 'compute' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of compute nodes of the cluster.
// Compute and AutoscaleCompute cannot be used together.
func (o *ClusterNodes) Compute() int {
	if o != nil && o.bitmap_&4 != 0 {
		return o.compute
	}
	return 0
}

// GetCompute returns the value of the 'compute' attribute and
// a flag indicating if the attribute has a value.
//
// Number of compute nodes of the cluster.
// Compute and AutoscaleCompute cannot be used together.
func (o *ClusterNodes) GetCompute() (value int, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.compute
	}
	return
}

// ComputeLabels returns the value of the 'compute_labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The labels set on the "default" compute machine pool.
func (o *ClusterNodes) ComputeLabels() map[string]string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.computeLabels
	}
	return nil
}

// GetComputeLabels returns the value of the 'compute_labels' attribute and
// a flag indicating if the attribute has a value.
//
// The labels set on the "default" compute machine pool.
func (o *ClusterNodes) GetComputeLabels() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.computeLabels
	}
	return
}

// ComputeMachineType returns the value of the 'compute_machine_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The compute machine type to use, for example `r5.xlarge`.
func (o *ClusterNodes) ComputeMachineType() *MachineType {
	if o != nil && o.bitmap_&16 != 0 {
		return o.computeMachineType
	}
	return nil
}

// GetComputeMachineType returns the value of the 'compute_machine_type' attribute and
// a flag indicating if the attribute has a value.
//
// The compute machine type to use, for example `r5.xlarge`.
func (o *ClusterNodes) GetComputeMachineType() (value *MachineType, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.computeMachineType
	}
	return
}

// ComputeRootVolume returns the value of the 'compute_root_volume' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The compute machine root volume capabilities.
func (o *ClusterNodes) ComputeRootVolume() *RootVolume {
	if o != nil && o.bitmap_&32 != 0 {
		return o.computeRootVolume
	}
	return nil
}

// GetComputeRootVolume returns the value of the 'compute_root_volume' attribute and
// a flag indicating if the attribute has a value.
//
// The compute machine root volume capabilities.
func (o *ClusterNodes) GetComputeRootVolume() (value *RootVolume, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.computeRootVolume
	}
	return
}

// Infra returns the value of the 'infra' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of infrastructure nodes of the cluster.
func (o *ClusterNodes) Infra() int {
	if o != nil && o.bitmap_&64 != 0 {
		return o.infra
	}
	return 0
}

// GetInfra returns the value of the 'infra' attribute and
// a flag indicating if the attribute has a value.
//
// Number of infrastructure nodes of the cluster.
func (o *ClusterNodes) GetInfra() (value int, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.infra
	}
	return
}

// InfraMachineType returns the value of the 'infra_machine_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The infra machine type to use, for example `r5.xlarge` (Optional).
func (o *ClusterNodes) InfraMachineType() *MachineType {
	if o != nil && o.bitmap_&128 != 0 {
		return o.infraMachineType
	}
	return nil
}

// GetInfraMachineType returns the value of the 'infra_machine_type' attribute and
// a flag indicating if the attribute has a value.
//
// The infra machine type to use, for example `r5.xlarge` (Optional).
func (o *ClusterNodes) GetInfraMachineType() (value *MachineType, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.infraMachineType
	}
	return
}

// Master returns the value of the 'master' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of master nodes of the cluster.
func (o *ClusterNodes) Master() int {
	if o != nil && o.bitmap_&256 != 0 {
		return o.master
	}
	return 0
}

// GetMaster returns the value of the 'master' attribute and
// a flag indicating if the attribute has a value.
//
// Number of master nodes of the cluster.
func (o *ClusterNodes) GetMaster() (value int, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.master
	}
	return
}

// MasterMachineType returns the value of the 'master_machine_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The master machine type to use, for example `r5.xlarge` (Optional).
func (o *ClusterNodes) MasterMachineType() *MachineType {
	if o != nil && o.bitmap_&512 != 0 {
		return o.masterMachineType
	}
	return nil
}

// GetMasterMachineType returns the value of the 'master_machine_type' attribute and
// a flag indicating if the attribute has a value.
//
// The master machine type to use, for example `r5.xlarge` (Optional).
func (o *ClusterNodes) GetMasterMachineType() (value *MachineType, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.masterMachineType
	}
	return
}

// SecurityGroupFilters returns the value of the 'security_group_filters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of security groups to be applied to nodes (Optional).
func (o *ClusterNodes) SecurityGroupFilters() []*MachinePoolSecurityGroupFilter {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.securityGroupFilters
	}
	return nil
}

// GetSecurityGroupFilters returns the value of the 'security_group_filters' attribute and
// a flag indicating if the attribute has a value.
//
// List of security groups to be applied to nodes (Optional).
func (o *ClusterNodes) GetSecurityGroupFilters() (value []*MachinePoolSecurityGroupFilter, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.securityGroupFilters
	}
	return
}

// Total returns the value of the 'total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Total number of nodes of the cluster.
func (o *ClusterNodes) Total() int {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.total
	}
	return 0
}

// GetTotal returns the value of the 'total' attribute and
// a flag indicating if the attribute has a value.
//
// Total number of nodes of the cluster.
func (o *ClusterNodes) GetTotal() (value int, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.total
	}
	return
}

// ClusterNodesListKind is the name of the type used to represent list of objects of
// type 'cluster_nodes'.
const ClusterNodesListKind = "ClusterNodesList"

// ClusterNodesListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_nodes'.
const ClusterNodesListLinkKind = "ClusterNodesListLink"

// ClusterNodesNilKind is the name of the type used to nil lists of objects of
// type 'cluster_nodes'.
const ClusterNodesListNilKind = "ClusterNodesListNil"

// ClusterNodesList is a list of values of the 'cluster_nodes' type.
type ClusterNodesList struct {
	href  string
	link  bool
	items []*ClusterNodes
}

// Len returns the length of the list.
func (l *ClusterNodesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterNodesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterNodesList) Get(i int) *ClusterNodes {
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
func (l *ClusterNodesList) Slice() []*ClusterNodes {
	var slice []*ClusterNodes
	if l == nil {
		slice = make([]*ClusterNodes, 0)
	} else {
		slice = make([]*ClusterNodes, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterNodesList) Each(f func(item *ClusterNodes) bool) {
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
func (l *ClusterNodesList) Range(f func(index int, item *ClusterNodes) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
