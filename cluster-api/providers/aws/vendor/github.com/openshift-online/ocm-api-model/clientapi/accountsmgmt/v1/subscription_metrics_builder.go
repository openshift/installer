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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// Each field is a metric fetched for a specific Subscription's cluster.
type SubscriptionMetricsBuilder struct {
	fieldSet_                    []bool
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
	return &SubscriptionMetricsBuilder{
		fieldSet_: make([]bool, 22),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubscriptionMetricsBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *SubscriptionMetricsBuilder) CloudProvider(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.cloudProvider = value
	b.fieldSet_[0] = true
	return b
}

// ComputeNodesCpu sets the value of the 'compute_nodes_cpu' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesCpu(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.computeNodesCpu = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// ComputeNodesMemory sets the value of the 'compute_nodes_memory' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesMemory(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.computeNodesMemory = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// ComputeNodesSockets sets the value of the 'compute_nodes_sockets' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ComputeNodesSockets(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.computeNodesSockets = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// ConsoleUrl sets the value of the 'console_url' attribute to the given value.
func (b *SubscriptionMetricsBuilder) ConsoleUrl(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.consoleUrl = value
	b.fieldSet_[4] = true
	return b
}

// Cpu sets the value of the 'cpu' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Cpu(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.cpu = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// CriticalAlertsFiring sets the value of the 'critical_alerts_firing' attribute to the given value.
func (b *SubscriptionMetricsBuilder) CriticalAlertsFiring(value float64) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.criticalAlertsFiring = value
	b.fieldSet_[6] = true
	return b
}

// HealthState sets the value of the 'health_state' attribute to the given value.
func (b *SubscriptionMetricsBuilder) HealthState(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.healthState = value
	b.fieldSet_[7] = true
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Memory(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.memory = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Nodes sets the value of the 'nodes' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Nodes(value *ClusterMetricsNodesBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.nodes = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// OpenshiftVersion sets the value of the 'openshift_version' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OpenshiftVersion(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.openshiftVersion = value
	b.fieldSet_[10] = true
	return b
}

// OperatingSystem sets the value of the 'operating_system' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OperatingSystem(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.operatingSystem = value
	b.fieldSet_[11] = true
	return b
}

// OperatorsConditionFailing sets the value of the 'operators_condition_failing' attribute to the given value.
func (b *SubscriptionMetricsBuilder) OperatorsConditionFailing(value float64) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.operatorsConditionFailing = value
	b.fieldSet_[12] = true
	return b
}

// Region sets the value of the 'region' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Region(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.region = value
	b.fieldSet_[13] = true
	return b
}

// Sockets sets the value of the 'sockets' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Sockets(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.sockets = value
	if value != nil {
		b.fieldSet_[14] = true
	} else {
		b.fieldSet_[14] = false
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *SubscriptionMetricsBuilder) State(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.state = value
	b.fieldSet_[15] = true
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *SubscriptionMetricsBuilder) StateDescription(value string) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.stateDescription = value
	b.fieldSet_[16] = true
	return b
}

// Storage sets the value of the 'storage' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Storage(value *ClusterResourceBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.storage = value
	if value != nil {
		b.fieldSet_[17] = true
	} else {
		b.fieldSet_[17] = false
	}
	return b
}

// SubscriptionCpuTotal sets the value of the 'subscription_cpu_total' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionCpuTotal(value float64) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.subscriptionCpuTotal = value
	b.fieldSet_[18] = true
	return b
}

// SubscriptionObligationExists sets the value of the 'subscription_obligation_exists' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionObligationExists(value float64) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.subscriptionObligationExists = value
	b.fieldSet_[19] = true
	return b
}

// SubscriptionSocketTotal sets the value of the 'subscription_socket_total' attribute to the given value.
func (b *SubscriptionMetricsBuilder) SubscriptionSocketTotal(value float64) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.subscriptionSocketTotal = value
	b.fieldSet_[20] = true
	return b
}

// Upgrade sets the value of the 'upgrade' attribute to the given value.
func (b *SubscriptionMetricsBuilder) Upgrade(value *ClusterUpgradeBuilder) *SubscriptionMetricsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 22)
	}
	b.upgrade = value
	if value != nil {
		b.fieldSet_[21] = true
	} else {
		b.fieldSet_[21] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubscriptionMetricsBuilder) Copy(object *SubscriptionMetrics) *SubscriptionMetricsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
