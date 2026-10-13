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

// ClusterAutoscalerKind is the name of the type used to represent objects
// of type 'cluster_autoscaler'.
const ClusterAutoscalerKind = "ClusterAutoscaler"

// ClusterAutoscalerLinkKind is the name of the type used to represent links
// to objects of type 'cluster_autoscaler'.
const ClusterAutoscalerLinkKind = "ClusterAutoscalerLink"

// ClusterAutoscalerNilKind is the name of the type used to nil references
// to objects of type 'cluster_autoscaler'.
const ClusterAutoscalerNilKind = "ClusterAutoscalerNil"

// ClusterAutoscaler represents the values of the 'cluster_autoscaler' type.
//
// Cluster-wide autoscaling configuration.
type ClusterAutoscaler struct {
	fieldSet_                   []bool
	id                          string
	href                        string
	balancingIgnoredLabels      []string
	logVerbosity                int
	maxNodeProvisionTime        string
	maxPodGracePeriod           int
	podPriorityThreshold        int
	resourceLimits              *AutoscalerResourceLimits
	scaleDown                   *AutoscalerScaleDownConfig
	balanceSimilarNodeGroups    bool
	ignoreDaemonsetsUtilization bool
	skipNodesWithLocalStorage   bool
}

// Kind returns the name of the type of the object.
func (o *ClusterAutoscaler) Kind() string {
	if o == nil {
		return ClusterAutoscalerNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ClusterAutoscalerLinkKind
	}
	return ClusterAutoscalerKind
}

// Link returns true if this is a link.
func (o *ClusterAutoscaler) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ClusterAutoscaler) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ClusterAutoscaler) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ClusterAutoscaler) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ClusterAutoscaler) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterAutoscaler) Empty() bool {
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

// BalanceSimilarNodeGroups returns the value of the 'balance_similar_node_groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// BalanceSimilarNodeGroups enables/disables the
// `--balance-similar-node-groups` cluster-autoscaler feature.
// This feature will automatically identify node groups with
// the same instance type and the same set of labels and try
// to keep the respective sizes of those node groups balanced.
func (o *ClusterAutoscaler) BalanceSimilarNodeGroups() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.balanceSimilarNodeGroups
	}
	return false
}

// GetBalanceSimilarNodeGroups returns the value of the 'balance_similar_node_groups' attribute and
// a flag indicating if the attribute has a value.
//
// BalanceSimilarNodeGroups enables/disables the
// `--balance-similar-node-groups` cluster-autoscaler feature.
// This feature will automatically identify node groups with
// the same instance type and the same set of labels and try
// to keep the respective sizes of those node groups balanced.
func (o *ClusterAutoscaler) GetBalanceSimilarNodeGroups() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.balanceSimilarNodeGroups
	}
	return
}

// BalancingIgnoredLabels returns the value of the 'balancing_ignored_labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// This option specifies labels that cluster autoscaler should ignore when considering node group similarity.
// For example, if you have nodes with "topology.ebs.csi.aws.com/zone" label, you can add name of this label here
// to prevent cluster autoscaler from splitting nodes into different node groups based on its value.
func (o *ClusterAutoscaler) BalancingIgnoredLabels() []string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.balancingIgnoredLabels
	}
	return nil
}

// GetBalancingIgnoredLabels returns the value of the 'balancing_ignored_labels' attribute and
// a flag indicating if the attribute has a value.
//
// This option specifies labels that cluster autoscaler should ignore when considering node group similarity.
// For example, if you have nodes with "topology.ebs.csi.aws.com/zone" label, you can add name of this label here
// to prevent cluster autoscaler from splitting nodes into different node groups based on its value.
func (o *ClusterAutoscaler) GetBalancingIgnoredLabels() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.balancingIgnoredLabels
	}
	return
}

// IgnoreDaemonsetsUtilization returns the value of the 'ignore_daemonsets_utilization' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Should CA ignore DaemonSet pods when calculating resource utilization for scaling down. false by default.
func (o *ClusterAutoscaler) IgnoreDaemonsetsUtilization() bool {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.ignoreDaemonsetsUtilization
	}
	return false
}

// GetIgnoreDaemonsetsUtilization returns the value of the 'ignore_daemonsets_utilization' attribute and
// a flag indicating if the attribute has a value.
//
// Should CA ignore DaemonSet pods when calculating resource utilization for scaling down. false by default.
func (o *ClusterAutoscaler) GetIgnoreDaemonsetsUtilization() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.ignoreDaemonsetsUtilization
	}
	return
}

// LogVerbosity returns the value of the 'log_verbosity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Sets the autoscaler log level.
// Default value is 1, level 4 is recommended for DEBUGGING and level 6 will enable almost everything.
func (o *ClusterAutoscaler) LogVerbosity() int {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.logVerbosity
	}
	return 0
}

// GetLogVerbosity returns the value of the 'log_verbosity' attribute and
// a flag indicating if the attribute has a value.
//
// Sets the autoscaler log level.
// Default value is 1, level 4 is recommended for DEBUGGING and level 6 will enable almost everything.
func (o *ClusterAutoscaler) GetLogVerbosity() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.logVerbosity
	}
	return
}

// MaxNodeProvisionTime returns the value of the 'max_node_provision_time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Maximum time CA waits for node to be provisioned.
func (o *ClusterAutoscaler) MaxNodeProvisionTime() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.maxNodeProvisionTime
	}
	return ""
}

// GetMaxNodeProvisionTime returns the value of the 'max_node_provision_time' attribute and
// a flag indicating if the attribute has a value.
//
// Maximum time CA waits for node to be provisioned.
func (o *ClusterAutoscaler) GetMaxNodeProvisionTime() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.maxNodeProvisionTime
	}
	return
}

// MaxPodGracePeriod returns the value of the 'max_pod_grace_period' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Gives pods graceful termination time before scaling down.
func (o *ClusterAutoscaler) MaxPodGracePeriod() int {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.maxPodGracePeriod
	}
	return 0
}

// GetMaxPodGracePeriod returns the value of the 'max_pod_grace_period' attribute and
// a flag indicating if the attribute has a value.
//
// Gives pods graceful termination time before scaling down.
func (o *ClusterAutoscaler) GetMaxPodGracePeriod() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.maxPodGracePeriod
	}
	return
}

// PodPriorityThreshold returns the value of the 'pod_priority_threshold' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// To allow users to schedule "best-effort" pods, which shouldn't trigger
// Cluster Autoscaler actions, but only run when there are spare resources available,
// More info: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-cluster-autoscaler-work-with-pod-priority-and-preemption.
func (o *ClusterAutoscaler) PodPriorityThreshold() int {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.podPriorityThreshold
	}
	return 0
}

// GetPodPriorityThreshold returns the value of the 'pod_priority_threshold' attribute and
// a flag indicating if the attribute has a value.
//
// To allow users to schedule "best-effort" pods, which shouldn't trigger
// Cluster Autoscaler actions, but only run when there are spare resources available,
// More info: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-cluster-autoscaler-work-with-pod-priority-and-preemption.
func (o *ClusterAutoscaler) GetPodPriorityThreshold() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.podPriorityThreshold
	}
	return
}

// ResourceLimits returns the value of the 'resource_limits' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Constraints of autoscaling resources.
func (o *ClusterAutoscaler) ResourceLimits() *AutoscalerResourceLimits {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.resourceLimits
	}
	return nil
}

// GetResourceLimits returns the value of the 'resource_limits' attribute and
// a flag indicating if the attribute has a value.
//
// Constraints of autoscaling resources.
func (o *ClusterAutoscaler) GetResourceLimits() (value *AutoscalerResourceLimits, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.resourceLimits
	}
	return
}

// ScaleDown returns the value of the 'scale_down' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Configuration of scale down operation.
func (o *ClusterAutoscaler) ScaleDown() *AutoscalerScaleDownConfig {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.scaleDown
	}
	return nil
}

// GetScaleDown returns the value of the 'scale_down' attribute and
// a flag indicating if the attribute has a value.
//
// Configuration of scale down operation.
func (o *ClusterAutoscaler) GetScaleDown() (value *AutoscalerScaleDownConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.scaleDown
	}
	return
}

// SkipNodesWithLocalStorage returns the value of the 'skip_nodes_with_local_storage' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Enables/Disables `--skip-nodes-with-local-storage` CA feature flag. If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath. true by default at autoscaler.
func (o *ClusterAutoscaler) SkipNodesWithLocalStorage() bool {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.skipNodesWithLocalStorage
	}
	return false
}

// GetSkipNodesWithLocalStorage returns the value of the 'skip_nodes_with_local_storage' attribute and
// a flag indicating if the attribute has a value.
//
// Enables/Disables `--skip-nodes-with-local-storage` CA feature flag. If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath. true by default at autoscaler.
func (o *ClusterAutoscaler) GetSkipNodesWithLocalStorage() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.skipNodesWithLocalStorage
	}
	return
}

// ClusterAutoscalerListKind is the name of the type used to represent list of objects of
// type 'cluster_autoscaler'.
const ClusterAutoscalerListKind = "ClusterAutoscalerList"

// ClusterAutoscalerListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_autoscaler'.
const ClusterAutoscalerListLinkKind = "ClusterAutoscalerListLink"

// ClusterAutoscalerNilKind is the name of the type used to nil lists of objects of
// type 'cluster_autoscaler'.
const ClusterAutoscalerListNilKind = "ClusterAutoscalerListNil"

// ClusterAutoscalerList is a list of values of the 'cluster_autoscaler' type.
type ClusterAutoscalerList struct {
	href  string
	link  bool
	items []*ClusterAutoscaler
}

// Kind returns the name of the type of the object.
func (l *ClusterAutoscalerList) Kind() string {
	if l == nil {
		return ClusterAutoscalerListNilKind
	}
	if l.link {
		return ClusterAutoscalerListLinkKind
	}
	return ClusterAutoscalerListKind
}

// Link returns true iif this is a link.
func (l *ClusterAutoscalerList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ClusterAutoscalerList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ClusterAutoscalerList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ClusterAutoscalerList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterAutoscalerList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterAutoscalerList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterAutoscalerList) SetItems(items []*ClusterAutoscaler) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterAutoscalerList) Items() []*ClusterAutoscaler {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterAutoscalerList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterAutoscalerList) Get(i int) *ClusterAutoscaler {
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
func (l *ClusterAutoscalerList) Slice() []*ClusterAutoscaler {
	var slice []*ClusterAutoscaler
	if l == nil {
		slice = make([]*ClusterAutoscaler, 0)
	} else {
		slice = make([]*ClusterAutoscaler, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterAutoscalerList) Each(f func(item *ClusterAutoscaler) bool) {
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
func (l *ClusterAutoscalerList) Range(f func(index int, item *ClusterAutoscaler) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
