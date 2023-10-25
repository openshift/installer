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

// SubscriptionMetricsBuilder contains the data and logic needed to build 'subscription_metrics' objects.
//
// Each field is a metric fetched for a specific Subscription's cluster.
type SubscriptionMetricsBuilder struct {
	bitmap_                      uint32
	cloudProvider                string
	computeNodesCpu              *ClusterResourceBuilder
	computeNodesMemory           *ClusterResourceBuilder
	computeNodesSockets          *ClusterResourceBuilder
	consoleUrl                   string
	cpu                          *ClusterResourceBuilder
	criticalAlertsFiring         float64
	healthState                  string
	memory                       *ClusterResourceBuilder
	nodes                        *ClusterMetricsNodesBuilder
	openshiftVersion             string
	operatingSystem              string
	operatorsConditionFailing    float64
	region                       string
	sockets                      *ClusterResourceBuilder
	state                        string
	stateDescription             string
	storage                      *ClusterResourceBuilder
	subscriptionCpuTotal         float64
	subscriptionObligationExists float64
	subscriptionSocketTotal      float64
	upgrade                      *ClusterUpgradeBuilder
}

// NewSubscriptionMetrics creates a new builder of 'subscription_metrics' objects.
func NewSubscriptionMetrics() *SubscriptionMetricsBuilder {
	return &SubscriptionMetricsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubscriptionMetricsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *SubscriptionMetricsBuilder) CloudProvider(value string) *SubscriptionMetricsBuilder {
	b.cloudProvider = value
	b.bitmap_ |= 1
	return b
}

// ComputeNodesCpu sets the value of the 'compute_nodes_cpu' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesCpu(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.computeNodesCpu = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// ComputeNodesMemory sets the value of the 'compute_nodes_memory' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesMemory(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.computeNodesMemory = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// ComputeNodesSockets sets the value of the 'compute_nodes_sockets' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesSockets(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.computeNodesSockets = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// ConsoleUrl sets the value of the 'console_url' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ConsoleUrl(value string) *SubscriptionMetricsBuilder {
	b.consoleUrl = value
	b.bitmap_ |= 16
	return b
}

// Cpu sets the value of the 'cpu' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Cpu(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.cpu = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// CriticalAlertsFiring sets the value of the 'critical_alerts_firing' attribute to the given value.
func (b *SubscriptionMetricsBuilder) CriticalAlertsFiring(value float64) *SubscriptionMetricsBuilder {
	b.criticalAlertsFiring = value
	b.bitmap_ |= 64
	return b
}

// HealthState sets the value of the 'health_state' attribute to the given value.
func (b *SubscriptionMetricsBuilder) HealthState(value string) *SubscriptionMetricsBuilder {
	b.healthState = value
	b.bitmap_ |= 128
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Memory(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.memory = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// Nodes sets the value of the 'nodes' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Nodes(value *ClusterMetricsNodesBuilder) *SubscriptionMetricsBuilder {
	b.nodes = value
	if value != nil {
		b.bitmap_ |= 512
	} else {
		b.bitmap_ &^= 512
	}
	return b
}

// OpenshiftVersion sets the value of the 'openshift_version' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OpenshiftVersion(value string) *SubscriptionMetricsBuilder {
	b.openshiftVersion = value
	b.bitmap_ |= 1024
	return b
}

// OperatingSystem sets the value of the 'operating_system' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OperatingSystem(value string) *SubscriptionMetricsBuilder {
	b.operatingSystem = value
	b.bitmap_ |= 2048
	return b
}

// OperatorsConditionFailing sets the value of the 'operators_condition_failing' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OperatorsConditionFailing(value float64) *SubscriptionMetricsBuilder {
	b.operatorsConditionFailing = value
	b.bitmap_ |= 4096
	return b
}

// Region sets the value of the 'region' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Region(value string) *SubscriptionMetricsBuilder {
	b.region = value
	b.bitmap_ |= 8192
	return b
}

// Sockets sets the value of the 'sockets' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Sockets(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.sockets = value
	if value != nil {
		b.bitmap_ |= 16384
	} else {
		b.bitmap_ &^= 16384
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *SubscriptionMetricsBuilder) State(value string) *SubscriptionMetricsBuilder {
	b.state = value
	b.bitmap_ |= 32768
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *SubscriptionMetricsBuilder) StateDescription(value string) *SubscriptionMetricsBuilder {
	b.stateDescription = value
	b.bitmap_ |= 65536
	return b
}

// Storage sets the value of the 'storage' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Storage(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	b.storage = value
	if value != nil {
		b.bitmap_ |= 131072
	} else {
		b.bitmap_ &^= 131072
	}
	return b
}

// SubscriptionCpuTotal sets the value of the 'subscription_cpu_total' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionCpuTotal(value float64) *SubscriptionMetricsBuilder {
	b.subscriptionCpuTotal = value
	b.bitmap_ |= 262144
	return b
}

// SubscriptionObligationExists sets the value of the 'subscription_obligation_exists' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionObligationExists(value float64) *SubscriptionMetricsBuilder {
	b.subscriptionObligationExists = value
	b.bitmap_ |= 524288
	return b
}

// SubscriptionSocketTotal sets the value of the 'subscription_socket_total' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionSocketTotal(value float64) *SubscriptionMetricsBuilder {
	b.subscriptionSocketTotal = value
	b.bitmap_ |= 1048576
	return b
}

// Upgrade sets the value of the 'upgrade' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Upgrade(value *ClusterUpgradeBuilder) *SubscriptionMetricsBuilder {
	b.upgrade = value
	if value != nil {
		b.bitmap_ |= 2097152
	} else {
		b.bitmap_ &^= 2097152
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubscriptionMetricsBuilder) Copy(object *SubscriptionMetrics) *SubscriptionMetricsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.cloudProvider = object.cloudProvider
	if object.computeNodesCpu != nil {
		b.computeNodesCpu = NewClusterResource().Copy(object.computeNodesCpu)
	} else {
		b.computeNodesCpu = nil
	}
	if object.computeNodesMemory != nil {
		b.computeNodesMemory = NewClusterResource().Copy(object.computeNodesMemory)
	} else {
		b.computeNodesMemory = nil
	}
	if object.computeNodesSockets != nil {
		b.computeNodesSockets = NewClusterResource().Copy(object.computeNodesSockets)
	} else {
		b.computeNodesSockets = nil
	}
	b.consoleUrl = object.consoleUrl
	if object.cpu != nil {
		b.cpu = NewClusterResource().Copy(object.cpu)
	} else {
		b.cpu = nil
	}
	b.criticalAlertsFiring = object.criticalAlertsFiring
	b.healthState = object.healthState
	if object.memory != nil {
		b.memory = NewClusterResource().Copy(object.memory)
	} else {
		b.memory = nil
	}
	if object.nodes != nil {
		b.nodes = NewClusterMetricsNodes().Copy(object.nodes)
	} else {
		b.nodes = nil
	}
	b.openshiftVersion = object.openshiftVersion
	b.operatingSystem = object.operatingSystem
	b.operatorsConditionFailing = object.operatorsConditionFailing
	b.region = object.region
	if object.sockets != nil {
		b.sockets = NewClusterResource().Copy(object.sockets)
	} else {
		b.sockets = nil
	}
	b.state = object.state
	b.stateDescription = object.stateDescription
	if object.storage != nil {
		b.storage = NewClusterResource().Copy(object.storage)
	} else {
		b.storage = nil
	}
	b.subscriptionCpuTotal = object.subscriptionCpuTotal
	b.subscriptionObligationExists = object.subscriptionObligationExists
	b.subscriptionSocketTotal = object.subscriptionSocketTotal
	if object.upgrade != nil {
		b.upgrade = NewClusterUpgrade().Copy(object.upgrade)
	} else {
		b.upgrade = nil
	}
	return b
}

// Build creates a 'subscription_metrics' object using the configuration stored in the builder.
func (b *SubscriptionMetricsBuilder) Build() (object *SubscriptionMetrics, err error) {
	object = new(SubscriptionMetrics)
	object.bitmap_ = b.bitmap_
	object.cloudProvider = b.cloudProvider
	if b.computeNodesCpu != nil {
		object.computeNodesCpu, err = b.computeNodesCpu.Build()
		if err != nil {
			return
		}
	}
	if b.computeNodesMemory != nil {
		object.computeNodesMemory, err = b.computeNodesMemory.Build()
		if err != nil {
			return
		}
	}
	if b.computeNodesSockets != nil {
		object.computeNodesSockets, err = b.computeNodesSockets.Build()
		if err != nil {
			return
		}
	}
	object.consoleUrl = b.consoleUrl
	if b.cpu != nil {
		object.cpu, err = b.cpu.Build()
		if err != nil {
			return
		}
	}
	object.criticalAlertsFiring = b.criticalAlertsFiring
	object.healthState = b.healthState
	if b.memory != nil {
		object.memory, err = b.memory.Build()
		if err != nil {
			return
		}
	}
	if b.nodes != nil {
		object.nodes, err = b.nodes.Build()
		if err != nil {
			return
		}
	}
	object.openshiftVersion = b.openshiftVersion
	object.operatingSystem = b.operatingSystem
	object.operatorsConditionFailing = b.operatorsConditionFailing
	object.region = b.region
	if b.sockets != nil {
		object.sockets, err = b.sockets.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	object.stateDescription = b.stateDescription
	if b.storage != nil {
		object.storage, err = b.storage.Build()
		if err != nil {
			return
		}
	}
	object.subscriptionCpuTotal = b.subscriptionCpuTotal
	object.subscriptionObligationExists = b.subscriptionObligationExists
	object.subscriptionSocketTotal = b.subscriptionSocketTotal
	if b.upgrade != nil {
		object.upgrade, err = b.upgrade.Build()
		if err != nil {
			return
		}
	}
	return
}
