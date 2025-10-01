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

// NodePoolKind is the name of the type used to represent objects
// of type 'node_pool'.
const NodePoolKind = "NodePool"

// NodePoolLinkKind is the name of the type used to represent links
// to objects of type 'node_pool'.
const NodePoolLinkKind = "NodePoolLink"

// NodePoolNilKind is the name of the type used to nil references
// to objects of type 'node_pool'.
const NodePoolNilKind = "NodePoolNil"

// NodePool represents the values of the 'node_pool' type.
//
// Representation of a node pool in a cluster.
type NodePool struct {
	bitmap_              uint32
	id                   string
	href                 string
	awsNodePool          *AWSNodePool
	autoscaling          *NodePoolAutoscaling
	availabilityZone     string
	azureNodePool        *AzureNodePool
	kubeletConfigs       []string
	labels               map[string]string
	managementUpgrade    *NodePoolManagementUpgrade
	nodeDrainGracePeriod *Value
	replicas             int
	status               *NodePoolStatus
	subnet               string
	taints               []*Taint
	tuningConfigs        []string
	version              *Version
	autoRepair           bool
}

// Kind returns the name of the type of the object.
func (o *NodePool) Kind() string {
	if o == nil {
		return NodePoolNilKind
	}
	if o.bitmap_&1 != 0 {
		return NodePoolLinkKind
	}
	return NodePoolKind
}

// Link returns true if this is a link.
func (o *NodePool) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *NodePool) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NodePool) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NodePool) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NodePool) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodePool) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AWSNodePool returns the value of the 'AWS_node_pool' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS specific parameters (Optional).
func (o *NodePool) AWSNodePool() *AWSNodePool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.awsNodePool
	}
	return nil
}

// GetAWSNodePool returns the value of the 'AWS_node_pool' attribute and
// a flag indicating if the attribute has a value.
//
// AWS specific parameters (Optional).
func (o *NodePool) GetAWSNodePool() (value *AWSNodePool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.awsNodePool
	}
	return
}

// AutoRepair returns the value of the 'auto_repair' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specifies whether health checks should be enabled for machines in the NodePool.
func (o *NodePool) AutoRepair() bool {
	if o != nil && o.bitmap_&16 != 0 {
		return o.autoRepair
	}
	return false
}

// GetAutoRepair returns the value of the 'auto_repair' attribute and
// a flag indicating if the attribute has a value.
//
// Specifies whether health checks should be enabled for machines in the NodePool.
func (o *NodePool) GetAutoRepair() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.autoRepair
	}
	return
}

// Autoscaling returns the value of the 'autoscaling' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details for auto-scaling the machine pool.
// Replicas and autoscaling cannot be used together.
func (o *NodePool) Autoscaling() *NodePoolAutoscaling {
	if o != nil && o.bitmap_&32 != 0 {
		return o.autoscaling
	}
	return nil
}

// GetAutoscaling returns the value of the 'autoscaling' attribute and
// a flag indicating if the attribute has a value.
//
// Details for auto-scaling the machine pool.
// Replicas and autoscaling cannot be used together.
func (o *NodePool) GetAutoscaling() (value *NodePoolAutoscaling, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.autoscaling
	}
	return
}

// AvailabilityZone returns the value of the 'availability_zone' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The availability zone upon which the node is created.
func (o *NodePool) AvailabilityZone() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.availabilityZone
	}
	return ""
}

// GetAvailabilityZone returns the value of the 'availability_zone' attribute and
// a flag indicating if the attribute has a value.
//
// The availability zone upon which the node is created.
func (o *NodePool) GetAvailabilityZone() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.availabilityZone
	}
	return
}

// AzureNodePool returns the value of the 'azure_node_pool' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Azure specific parameters.
func (o *NodePool) AzureNodePool() *AzureNodePool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.azureNodePool
	}
	return nil
}

// GetAzureNodePool returns the value of the 'azure_node_pool' attribute and
// a flag indicating if the attribute has a value.
//
// Azure specific parameters.
func (o *NodePool) GetAzureNodePool() (value *AzureNodePool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.azureNodePool
	}
	return
}

// KubeletConfigs returns the value of the 'kubelet_configs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The names of the KubeletConfigs for this node pool.
func (o *NodePool) KubeletConfigs() []string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.kubeletConfigs
	}
	return nil
}

// GetKubeletConfigs returns the value of the 'kubelet_configs' attribute and
// a flag indicating if the attribute has a value.
//
// The names of the KubeletConfigs for this node pool.
func (o *NodePool) GetKubeletConfigs() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.kubeletConfigs
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The labels set on the Nodes created.
func (o *NodePool) Labels() map[string]string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
//
// The labels set on the Nodes created.
func (o *NodePool) GetLabels() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.labels
	}
	return
}

// ManagementUpgrade returns the value of the 'management_upgrade' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Management parameters (Optional).
func (o *NodePool) ManagementUpgrade() *NodePoolManagementUpgrade {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.managementUpgrade
	}
	return nil
}

// GetManagementUpgrade returns the value of the 'management_upgrade' attribute and
// a flag indicating if the attribute has a value.
//
// Management parameters (Optional).
func (o *NodePool) GetManagementUpgrade() (value *NodePoolManagementUpgrade, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.managementUpgrade
	}
	return
}

// NodeDrainGracePeriod returns the value of the 'node_drain_grace_period' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Time to wait for a NodePool to drain when it is upgraded or replaced before it is forcibly removed.
func (o *NodePool) NodeDrainGracePeriod() *Value {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.nodeDrainGracePeriod
	}
	return nil
}

// GetNodeDrainGracePeriod returns the value of the 'node_drain_grace_period' attribute and
// a flag indicating if the attribute has a value.
//
// Time to wait for a NodePool to drain when it is upgraded or replaced before it is forcibly removed.
func (o *NodePool) GetNodeDrainGracePeriod() (value *Value, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.nodeDrainGracePeriod
	}
	return
}

// Replicas returns the value of the 'replicas' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The number of Machines (and Nodes) to create.
// Replicas and autoscaling cannot be used together.
func (o *NodePool) Replicas() int {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.replicas
	}
	return 0
}

// GetReplicas returns the value of the 'replicas' attribute and
// a flag indicating if the attribute has a value.
//
// The number of Machines (and Nodes) to create.
// Replicas and autoscaling cannot be used together.
func (o *NodePool) GetReplicas() (value int, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.replicas
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// NodePool status.
func (o *NodePool) Status() *NodePoolStatus {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// NodePool status.
func (o *NodePool) GetStatus() (value *NodePoolStatus, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.status
	}
	return
}

// Subnet returns the value of the 'subnet' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The subnet upon which the nodes are created.
func (o *NodePool) Subnet() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.subnet
	}
	return ""
}

// GetSubnet returns the value of the 'subnet' attribute and
// a flag indicating if the attribute has a value.
//
// The subnet upon which the nodes are created.
func (o *NodePool) GetSubnet() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.subnet
	}
	return
}

// Taints returns the value of the 'taints' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The taints set on the Nodes created.
func (o *NodePool) Taints() []*Taint {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.taints
	}
	return nil
}

// GetTaints returns the value of the 'taints' attribute and
// a flag indicating if the attribute has a value.
//
// The taints set on the Nodes created.
func (o *NodePool) GetTaints() (value []*Taint, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.taints
	}
	return
}

// TuningConfigs returns the value of the 'tuning_configs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The names of the tuning configs for this node pool.
func (o *NodePool) TuningConfigs() []string {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.tuningConfigs
	}
	return nil
}

// GetTuningConfigs returns the value of the 'tuning_configs' attribute and
// a flag indicating if the attribute has a value.
//
// The names of the tuning configs for this node pool.
func (o *NodePool) GetTuningConfigs() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.tuningConfigs
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version of the node pool.
func (o *NodePool) Version() *Version {
	if o != nil && o.bitmap_&131072 != 0 {
		return o.version
	}
	return nil
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Version of the node pool.
func (o *NodePool) GetVersion() (value *Version, ok bool) {
	ok = o != nil && o.bitmap_&131072 != 0
	if ok {
		value = o.version
	}
	return
}

// NodePoolListKind is the name of the type used to represent list of objects of
// type 'node_pool'.
const NodePoolListKind = "NodePoolList"

// NodePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_pool'.
const NodePoolListLinkKind = "NodePoolListLink"

// NodePoolNilKind is the name of the type used to nil lists of objects of
// type 'node_pool'.
const NodePoolListNilKind = "NodePoolListNil"

// NodePoolList is a list of values of the 'node_pool' type.
type NodePoolList struct {
	href  string
	link  bool
	items []*NodePool
}

// Kind returns the name of the type of the object.
func (l *NodePoolList) Kind() string {
	if l == nil {
		return NodePoolListNilKind
	}
	if l.link {
		return NodePoolListLinkKind
	}
	return NodePoolListKind
}

// Link returns true iif this is a link.
func (l *NodePoolList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NodePoolList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NodePoolList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NodePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NodePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NodePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NodePoolList) SetItems(items []*NodePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *NodePoolList) Items() []*NodePool {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NodePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodePoolList) Get(i int) *NodePool {
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
func (l *NodePoolList) Slice() []*NodePool {
	var slice []*NodePool
	if l == nil {
		slice = make([]*NodePool, 0)
	} else {
		slice = make([]*NodePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodePoolList) Each(f func(item *NodePool) bool) {
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
func (l *NodePoolList) Range(f func(index int, item *NodePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
