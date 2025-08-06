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

// NodePoolBuilder contains the data and logic needed to build 'node_pool' objects.
//
// Representation of a node pool in a cluster.
type NodePoolBuilder struct {
	bitmap_              uint32
	id                   string
	href                 string
	awsNodePool          *AWSNodePoolBuilder
	autoscaling          *NodePoolAutoscalingBuilder
	availabilityZone     string
	azureNodePool        *AzureNodePoolBuilder
	kubeletConfigs       []string
	labels               map[string]string
	managementUpgrade    *NodePoolManagementUpgradeBuilder
	nodeDrainGracePeriod *ValueBuilder
	replicas             int
	status               *NodePoolStatusBuilder
	subnet               string
	taints               []*TaintBuilder
	tuningConfigs        []string
	version              *VersionBuilder
	autoRepair           bool
}

// NewNodePool creates a new builder of 'node_pool' objects.
func NewNodePool() *NodePoolBuilder {
	return &NodePoolBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolBuilder) Link(value bool) *NodePoolBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolBuilder) ID(value string) *NodePoolBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *NodePoolBuilder) HREF(value string) *NodePoolBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AWSNodePool sets the value of the 'AWS_node_pool' attribute to the given value.
//
// Representation of aws node pool specific parameters.
func (b *NodePoolBuilder) AWSNodePool(value *AWSNodePoolBuilder) *NodePoolBuilder {
	b.awsNodePool = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// AutoRepair sets the value of the 'auto_repair' attribute to the given value.
func (b *NodePoolBuilder) AutoRepair(value bool) *NodePoolBuilder {
	b.autoRepair = value
	b.bitmap_ |= 16
	return b
}

// Autoscaling sets the value of the 'autoscaling' attribute to the given value.
//
// Representation of a autoscaling in a node pool.
func (b *NodePoolBuilder) Autoscaling(value *NodePoolAutoscalingBuilder) *NodePoolBuilder {
	b.autoscaling = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *NodePoolBuilder) AvailabilityZone(value string) *NodePoolBuilder {
	b.availabilityZone = value
	b.bitmap_ |= 64
	return b
}

// AzureNodePool sets the value of the 'azure_node_pool' attribute to the given value.
//
// Representation of azure node pool specific parameters.
func (b *NodePoolBuilder) AzureNodePool(value *AzureNodePoolBuilder) *NodePoolBuilder {
	b.azureNodePool = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// KubeletConfigs sets the value of the 'kubelet_configs' attribute to the given values.
func (b *NodePoolBuilder) KubeletConfigs(values ...string) *NodePoolBuilder {
	b.kubeletConfigs = make([]string, len(values))
	copy(b.kubeletConfigs, values)
	b.bitmap_ |= 256
	return b
}

// Labels sets the value of the 'labels' attribute to the given value.
func (b *NodePoolBuilder) Labels(value map[string]string) *NodePoolBuilder {
	b.labels = value
	if value != nil {
		b.bitmap_ |= 512
	} else {
		b.bitmap_ &^= 512
	}
	return b
}

// ManagementUpgrade sets the value of the 'management_upgrade' attribute to the given value.
//
// Representation of node pool management.
func (b *NodePoolBuilder) ManagementUpgrade(value *NodePoolManagementUpgradeBuilder) *NodePoolBuilder {
	b.managementUpgrade = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// NodeDrainGracePeriod sets the value of the 'node_drain_grace_period' attribute to the given value.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
func (b *NodePoolBuilder) NodeDrainGracePeriod(value *ValueBuilder) *NodePoolBuilder {
	b.nodeDrainGracePeriod = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// Replicas sets the value of the 'replicas' attribute to the given value.
func (b *NodePoolBuilder) Replicas(value int) *NodePoolBuilder {
	b.replicas = value
	b.bitmap_ |= 4096
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of the status of a node pool.
func (b *NodePoolBuilder) Status(value *NodePoolStatusBuilder) *NodePoolBuilder {
	b.status = value
	if value != nil {
		b.bitmap_ |= 8192
	} else {
		b.bitmap_ &^= 8192
	}
	return b
}

// Subnet sets the value of the 'subnet' attribute to the given value.
func (b *NodePoolBuilder) Subnet(value string) *NodePoolBuilder {
	b.subnet = value
	b.bitmap_ |= 16384
	return b
}

// Taints sets the value of the 'taints' attribute to the given values.
func (b *NodePoolBuilder) Taints(values ...*TaintBuilder) *NodePoolBuilder {
	b.taints = make([]*TaintBuilder, len(values))
	copy(b.taints, values)
	b.bitmap_ |= 32768
	return b
}

// TuningConfigs sets the value of the 'tuning_configs' attribute to the given values.
func (b *NodePoolBuilder) TuningConfigs(values ...string) *NodePoolBuilder {
	b.tuningConfigs = make([]string, len(values))
	copy(b.tuningConfigs, values)
	b.bitmap_ |= 65536
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an _OpenShift_ version.
func (b *NodePoolBuilder) Version(value *VersionBuilder) *NodePoolBuilder {
	b.version = value
	if value != nil {
		b.bitmap_ |= 131072
	} else {
		b.bitmap_ &^= 131072
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolBuilder) Copy(object *NodePool) *NodePoolBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.awsNodePool != nil {
		b.awsNodePool = NewAWSNodePool().Copy(object.awsNodePool)
	} else {
		b.awsNodePool = nil
	}
	b.autoRepair = object.autoRepair
	if object.autoscaling != nil {
		b.autoscaling = NewNodePoolAutoscaling().Copy(object.autoscaling)
	} else {
		b.autoscaling = nil
	}
	b.availabilityZone = object.availabilityZone
	if object.azureNodePool != nil {
		b.azureNodePool = NewAzureNodePool().Copy(object.azureNodePool)
	} else {
		b.azureNodePool = nil
	}
	if object.kubeletConfigs != nil {
		b.kubeletConfigs = make([]string, len(object.kubeletConfigs))
		copy(b.kubeletConfigs, object.kubeletConfigs)
	} else {
		b.kubeletConfigs = nil
	}
	if len(object.labels) > 0 {
		b.labels = map[string]string{}
		for k, v := range object.labels {
			b.labels[k] = v
		}
	} else {
		b.labels = nil
	}
	if object.managementUpgrade != nil {
		b.managementUpgrade = NewNodePoolManagementUpgrade().Copy(object.managementUpgrade)
	} else {
		b.managementUpgrade = nil
	}
	if object.nodeDrainGracePeriod != nil {
		b.nodeDrainGracePeriod = NewValue().Copy(object.nodeDrainGracePeriod)
	} else {
		b.nodeDrainGracePeriod = nil
	}
	b.replicas = object.replicas
	if object.status != nil {
		b.status = NewNodePoolStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	b.subnet = object.subnet
	if object.taints != nil {
		b.taints = make([]*TaintBuilder, len(object.taints))
		for i, v := range object.taints {
			b.taints[i] = NewTaint().Copy(v)
		}
	} else {
		b.taints = nil
	}
	if object.tuningConfigs != nil {
		b.tuningConfigs = make([]string, len(object.tuningConfigs))
		copy(b.tuningConfigs, object.tuningConfigs)
	} else {
		b.tuningConfigs = nil
	}
	if object.version != nil {
		b.version = NewVersion().Copy(object.version)
	} else {
		b.version = nil
	}
	return b
}

// Build creates a 'node_pool' object using the configuration stored in the builder.
func (b *NodePoolBuilder) Build() (object *NodePool, err error) {
	object = new(NodePool)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.awsNodePool != nil {
		object.awsNodePool, err = b.awsNodePool.Build()
		if err != nil {
			return
		}
	}
	object.autoRepair = b.autoRepair
	if b.autoscaling != nil {
		object.autoscaling, err = b.autoscaling.Build()
		if err != nil {
			return
		}
	}
	object.availabilityZone = b.availabilityZone
	if b.azureNodePool != nil {
		object.azureNodePool, err = b.azureNodePool.Build()
		if err != nil {
			return
		}
	}
	if b.kubeletConfigs != nil {
		object.kubeletConfigs = make([]string, len(b.kubeletConfigs))
		copy(object.kubeletConfigs, b.kubeletConfigs)
	}
	if b.labels != nil {
		object.labels = make(map[string]string)
		for k, v := range b.labels {
			object.labels[k] = v
		}
	}
	if b.managementUpgrade != nil {
		object.managementUpgrade, err = b.managementUpgrade.Build()
		if err != nil {
			return
		}
	}
	if b.nodeDrainGracePeriod != nil {
		object.nodeDrainGracePeriod, err = b.nodeDrainGracePeriod.Build()
		if err != nil {
			return
		}
	}
	object.replicas = b.replicas
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	object.subnet = b.subnet
	if b.taints != nil {
		object.taints = make([]*Taint, len(b.taints))
		for i, v := range b.taints {
			object.taints[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.tuningConfigs != nil {
		object.tuningConfigs = make([]string, len(b.tuningConfigs))
		copy(object.tuningConfigs, b.tuningConfigs)
	}
	if b.version != nil {
		object.version, err = b.version.Build()
		if err != nil {
			return
		}
	}
	return
}
