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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// SubscriptionMetrics represents the values of the 'subscription_metrics' type.
//
// Each field is a metric fetched for a specific Subscription's cluster.
type SubscriptionMetrics struct {
	bitmap_                      uint32
	cloudProvider                string
	computeNodesCpu              *ClusterResource
	computeNodesMemory           *ClusterResource
	computeNodesSockets          *ClusterResource
	consoleUrl                   string
	cpu                          *ClusterResource
	criticalAlertsFiring         float64
	healthState                  string
	memory                       *ClusterResource
	nodes                        *ClusterMetricsNodes
	openshiftVersion             string
	operatingSystem              string
	operatorsConditionFailing    float64
	region                       string
	sockets                      *ClusterResource
	state                        string
	stateDescription             string
	storage                      *ClusterResource
	subscriptionCpuTotal         float64
	subscriptionObligationExists float64
	subscriptionSocketTotal      float64
	upgrade                      *ClusterUpgrade
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SubscriptionMetrics) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) CloudProvider() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cloudProvider
	}
	return ""
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetCloudProvider() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// ComputeNodesCpu returns the value of the 'compute_nodes_cpu' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) ComputeNodesCpu() *ClusterResource {
	if o != nil && o.bitmap_&2 != 0 {
		return o.computeNodesCpu
	}
	return nil
}

// GetComputeNodesCpu returns the value of the 'compute_nodes_cpu' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetComputeNodesCpu() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.computeNodesCpu
	}
	return
}

// ComputeNodesMemory returns the value of the 'compute_nodes_memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) ComputeNodesMemory() *ClusterResource {
	if o != nil && o.bitmap_&4 != 0 {
		return o.computeNodesMemory
	}
	return nil
}

// GetComputeNodesMemory returns the value of the 'compute_nodes_memory' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetComputeNodesMemory() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.computeNodesMemory
	}
	return
}

// ComputeNodesSockets returns the value of the 'compute_nodes_sockets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) ComputeNodesSockets() *ClusterResource {
	if o != nil && o.bitmap_&8 != 0 {
		return o.computeNodesSockets
	}
	return nil
}

// GetComputeNodesSockets returns the value of the 'compute_nodes_sockets' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetComputeNodesSockets() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.computeNodesSockets
	}
	return
}

// ConsoleUrl returns the value of the 'console_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) ConsoleUrl() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.consoleUrl
	}
	return ""
}

// GetConsoleUrl returns the value of the 'console_url' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetConsoleUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.consoleUrl
	}
	return
}

// Cpu returns the value of the 'cpu' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Cpu() *ClusterResource {
	if o != nil && o.bitmap_&32 != 0 {
		return o.cpu
	}
	return nil
}

// GetCpu returns the value of the 'cpu' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetCpu() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.cpu
	}
	return
}

// CriticalAlertsFiring returns the value of the 'critical_alerts_firing' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) CriticalAlertsFiring() float64 {
	if o != nil && o.bitmap_&64 != 0 {
		return o.criticalAlertsFiring
	}
	return 0.0
}

// GetCriticalAlertsFiring returns the value of the 'critical_alerts_firing' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetCriticalAlertsFiring() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.criticalAlertsFiring
	}
	return
}

// HealthState returns the value of the 'health_state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) HealthState() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.healthState
	}
	return ""
}

// GetHealthState returns the value of the 'health_state' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetHealthState() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.healthState
	}
	return
}

// Memory returns the value of the 'memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Memory() *ClusterResource {
	if o != nil && o.bitmap_&256 != 0 {
		return o.memory
	}
	return nil
}

// GetMemory returns the value of the 'memory' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetMemory() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.memory
	}
	return
}

// Nodes returns the value of the 'nodes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Nodes() *ClusterMetricsNodes {
	if o != nil && o.bitmap_&512 != 0 {
		return o.nodes
	}
	return nil
}

// GetNodes returns the value of the 'nodes' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetNodes() (value *ClusterMetricsNodes, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.nodes
	}
	return
}

// OpenshiftVersion returns the value of the 'openshift_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) OpenshiftVersion() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.openshiftVersion
	}
	return ""
}

// GetOpenshiftVersion returns the value of the 'openshift_version' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetOpenshiftVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.openshiftVersion
	}
	return
}

// OperatingSystem returns the value of the 'operating_system' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) OperatingSystem() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.operatingSystem
	}
	return ""
}

// GetOperatingSystem returns the value of the 'operating_system' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetOperatingSystem() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.operatingSystem
	}
	return
}

// OperatorsConditionFailing returns the value of the 'operators_condition_failing' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) OperatorsConditionFailing() float64 {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.operatorsConditionFailing
	}
	return 0.0
}

// GetOperatorsConditionFailing returns the value of the 'operators_condition_failing' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetOperatorsConditionFailing() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.operatorsConditionFailing
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Region() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.region
	}
	return ""
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetRegion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.region
	}
	return
}

// Sockets returns the value of the 'sockets' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Sockets() *ClusterResource {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.sockets
	}
	return nil
}

// GetSockets returns the value of the 'sockets' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetSockets() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.sockets
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) State() string {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.state
	}
	return ""
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetState() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.state
	}
	return
}

// StateDescription returns the value of the 'state_description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) StateDescription() string {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.stateDescription
	}
	return ""
}

// GetStateDescription returns the value of the 'state_description' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetStateDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.stateDescription
	}
	return
}

// Storage returns the value of the 'storage' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Storage() *ClusterResource {
	if o != nil && o.bitmap_&131072 != 0 {
		return o.storage
	}
	return nil
}

// GetStorage returns the value of the 'storage' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetStorage() (value *ClusterResource, ok bool) {
	ok = o != nil && o.bitmap_&131072 != 0
	if ok {
		value = o.storage
	}
	return
}

// SubscriptionCpuTotal returns the value of the 'subscription_cpu_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) SubscriptionCpuTotal() float64 {
	if o != nil && o.bitmap_&262144 != 0 {
		return o.subscriptionCpuTotal
	}
	return 0.0
}

// GetSubscriptionCpuTotal returns the value of the 'subscription_cpu_total' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetSubscriptionCpuTotal() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&262144 != 0
	if ok {
		value = o.subscriptionCpuTotal
	}
	return
}

// SubscriptionObligationExists returns the value of the 'subscription_obligation_exists' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) SubscriptionObligationExists() float64 {
	if o != nil && o.bitmap_&524288 != 0 {
		return o.subscriptionObligationExists
	}
	return 0.0
}

// GetSubscriptionObligationExists returns the value of the 'subscription_obligation_exists' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetSubscriptionObligationExists() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&524288 != 0
	if ok {
		value = o.subscriptionObligationExists
	}
	return
}

// SubscriptionSocketTotal returns the value of the 'subscription_socket_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) SubscriptionSocketTotal() float64 {
	if o != nil && o.bitmap_&1048576 != 0 {
		return o.subscriptionSocketTotal
	}
	return 0.0
}

// GetSubscriptionSocketTotal returns the value of the 'subscription_socket_total' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetSubscriptionSocketTotal() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&1048576 != 0
	if ok {
		value = o.subscriptionSocketTotal
	}
	return
}

// Upgrade returns the value of the 'upgrade' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SubscriptionMetrics) Upgrade() *ClusterUpgrade {
	if o != nil && o.bitmap_&2097152 != 0 {
		return o.upgrade
	}
	return nil
}

// GetUpgrade returns the value of the 'upgrade' attribute and
// a flag indicating if the attribute has a value.
func (o *SubscriptionMetrics) GetUpgrade() (value *ClusterUpgrade, ok bool) {
	ok = o != nil && o.bitmap_&2097152 != 0
	if ok {
		value = o.upgrade
	}
	return
}

// SubscriptionMetricsListKind is the name of the type used to represent list of objects of
// type 'subscription_metrics'.
const SubscriptionMetricsListKind = "SubscriptionMetricsList"

// SubscriptionMetricsListLinkKind is the name of the type used to represent links to list
// of objects of type 'subscription_metrics'.
const SubscriptionMetricsListLinkKind = "SubscriptionMetricsListLink"

// SubscriptionMetricsNilKind is the name of the type used to nil lists of objects of
// type 'subscription_metrics'.
const SubscriptionMetricsListNilKind = "SubscriptionMetricsListNil"

// SubscriptionMetricsList is a list of values of the 'subscription_metrics' type.
type SubscriptionMetricsList struct {
	href  string
	link  bool
	items []*SubscriptionMetrics
}

// Len returns the length of the list.
func (l *SubscriptionMetricsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SubscriptionMetricsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SubscriptionMetricsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SubscriptionMetricsList) SetItems(items []*SubscriptionMetrics) {
	l.items = items
}

// Items returns the items of the list.
func (l *SubscriptionMetricsList) Items() []*SubscriptionMetrics {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SubscriptionMetricsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SubscriptionMetricsList) Get(i int) *SubscriptionMetrics {
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
func (l *SubscriptionMetricsList) Slice() []*SubscriptionMetrics {
	var slice []*SubscriptionMetrics
	if l == nil {
		slice = make([]*SubscriptionMetrics, 0)
	} else {
		slice = make([]*SubscriptionMetrics, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SubscriptionMetricsList) Each(f func(item *SubscriptionMetrics) bool) {
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
func (l *SubscriptionMetricsList) Range(f func(index int, item *SubscriptionMetrics) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
