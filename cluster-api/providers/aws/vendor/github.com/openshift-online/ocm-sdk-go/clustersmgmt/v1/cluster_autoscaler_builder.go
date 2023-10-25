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

// ClusterAutoscalerBuilder contains the data and logic needed to build 'cluster_autoscaler' objects.
//
// Cluster-wide autoscaling configuration.
type ClusterAutoscalerBuilder struct {
	bitmap_                     uint32
	id                          string
	href                        string
	balancingIgnoredLabels      []string
	logVerbosity                int
	maxNodeProvisionTime        string
	maxPodGracePeriod           int
	podPriorityThreshold        int
	resourceLimits              *AutoscalerResourceLimitsBuilder
	scaleDown                   *AutoscalerScaleDownConfigBuilder
	balanceSimilarNodeGroups    bool
	ignoreDaemonsetsUtilization bool
	skipNodesWithLocalStorage   bool
}

// NewClusterAutoscaler creates a new builder of 'cluster_autoscaler' objects.
func NewClusterAutoscaler() *ClusterAutoscalerBuilder {
	return &ClusterAutoscalerBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterAutoscalerBuilder) Link(value bool) *ClusterAutoscalerBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ClusterAutoscalerBuilder) ID(value string) *ClusterAutoscalerBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ClusterAutoscalerBuilder) HREF(value string) *ClusterAutoscalerBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAutoscalerBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BalanceSimilarNodeGroups sets the value of the 'balance_similar_node_groups' attribute to the given value.
func (b *ClusterAutoscalerBuilder) BalanceSimilarNodeGroups(value bool) *ClusterAutoscalerBuilder {
	b.balanceSimilarNodeGroups = value
	b.bitmap_ |= 8
	return b
}

// BalancingIgnoredLabels sets the value of the 'balancing_ignored_labels' attribute to the given values.
func (b *ClusterAutoscalerBuilder) BalancingIgnoredLabels(values ...string) *ClusterAutoscalerBuilder {
	b.balancingIgnoredLabels = make([]string, len(values))
	copy(b.balancingIgnoredLabels, values)
	b.bitmap_ |= 16
	return b
}

// IgnoreDaemonsetsUtilization sets the value of the 'ignore_daemonsets_utilization' attribute to the given value.
func (b *ClusterAutoscalerBuilder) IgnoreDaemonsetsUtilization(value bool) *ClusterAutoscalerBuilder {
	b.ignoreDaemonsetsUtilization = value
	b.bitmap_ |= 32
	return b
}

// LogVerbosity sets the value of the 'log_verbosity' attribute to the given value.
func (b *ClusterAutoscalerBuilder) LogVerbosity(value int) *ClusterAutoscalerBuilder {
	b.logVerbosity = value
	b.bitmap_ |= 64
	return b
}

// MaxNodeProvisionTime sets the value of the 'max_node_provision_time' attribute to the given value.
func (b *ClusterAutoscalerBuilder) MaxNodeProvisionTime(value string) *ClusterAutoscalerBuilder {
	b.maxNodeProvisionTime = value
	b.bitmap_ |= 128
	return b
}

// MaxPodGracePeriod sets the value of the 'max_pod_grace_period' attribute to the given value.
func (b *ClusterAutoscalerBuilder) MaxPodGracePeriod(value int) *ClusterAutoscalerBuilder {
	b.maxPodGracePeriod = value
	b.bitmap_ |= 256
	return b
}

// PodPriorityThreshold sets the value of the 'pod_priority_threshold' attribute to the given value.
func (b *ClusterAutoscalerBuilder) PodPriorityThreshold(value int) *ClusterAutoscalerBuilder {
	b.podPriorityThreshold = value
	b.bitmap_ |= 512
	return b
}

// ResourceLimits sets the value of the 'resource_limits' attribute to the given value.
func (b *ClusterAutoscalerBuilder) ResourceLimits(value *AutoscalerResourceLimitsBuilder) *ClusterAutoscalerBuilder {
	b.resourceLimits = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// ScaleDown sets the value of the 'scale_down' attribute to the given value.
func (b *ClusterAutoscalerBuilder) ScaleDown(value *AutoscalerScaleDownConfigBuilder) *ClusterAutoscalerBuilder {
	b.scaleDown = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// SkipNodesWithLocalStorage sets the value of the 'skip_nodes_with_local_storage' attribute to the given value.
func (b *ClusterAutoscalerBuilder) SkipNodesWithLocalStorage(value bool) *ClusterAutoscalerBuilder {
	b.skipNodesWithLocalStorage = value
	b.bitmap_ |= 4096
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAutoscalerBuilder) Copy(object *ClusterAutoscaler) *ClusterAutoscalerBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.balanceSimilarNodeGroups = object.balanceSimilarNodeGroups
	if object.balancingIgnoredLabels != nil {
		b.balancingIgnoredLabels = make([]string, len(object.balancingIgnoredLabels))
		copy(b.balancingIgnoredLabels, object.balancingIgnoredLabels)
	} else {
		b.balancingIgnoredLabels = nil
	}
	b.ignoreDaemonsetsUtilization = object.ignoreDaemonsetsUtilization
	b.logVerbosity = object.logVerbosity
	b.maxNodeProvisionTime = object.maxNodeProvisionTime
	b.maxPodGracePeriod = object.maxPodGracePeriod
	b.podPriorityThreshold = object.podPriorityThreshold
	if object.resourceLimits != nil {
		b.resourceLimits = NewAutoscalerResourceLimits().Copy(object.resourceLimits)
	} else {
		b.resourceLimits = nil
	}
	if object.scaleDown != nil {
		b.scaleDown = NewAutoscalerScaleDownConfig().Copy(object.scaleDown)
	} else {
		b.scaleDown = nil
	}
	b.skipNodesWithLocalStorage = object.skipNodesWithLocalStorage
	return b
}

// Build creates a 'cluster_autoscaler' object using the configuration stored in the builder.
func (b *ClusterAutoscalerBuilder) Build() (object *ClusterAutoscaler, err error) {
	object = new(ClusterAutoscaler)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.balanceSimilarNodeGroups = b.balanceSimilarNodeGroups
	if b.balancingIgnoredLabels != nil {
		object.balancingIgnoredLabels = make([]string, len(b.balancingIgnoredLabels))
		copy(object.balancingIgnoredLabels, b.balancingIgnoredLabels)
	}
	object.ignoreDaemonsetsUtilization = b.ignoreDaemonsetsUtilization
	object.logVerbosity = b.logVerbosity
	object.maxNodeProvisionTime = b.maxNodeProvisionTime
	object.maxPodGracePeriod = b.maxPodGracePeriod
	object.podPriorityThreshold = b.podPriorityThreshold
	if b.resourceLimits != nil {
		object.resourceLimits, err = b.resourceLimits.Build()
		if err != nil {
			return
		}
	}
	if b.scaleDown != nil {
		object.scaleDown, err = b.scaleDown.Build()
		if err != nil {
			return
		}
	}
	object.skipNodesWithLocalStorage = b.skipNodesWithLocalStorage
	return
}
