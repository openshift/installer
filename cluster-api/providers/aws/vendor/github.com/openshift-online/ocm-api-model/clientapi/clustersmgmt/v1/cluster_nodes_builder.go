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

// Counts of different classes of nodes inside a cluster.
type ClusterNodesBuilder struct {
	fieldSet_            []bool
	autoscaleCompute     *MachinePoolAutoscalingBuilder
	availabilityZones    []string
	compute              int
	computeLabels        map[string]string
	computeMachineType   *MachineTypeBuilder
	computeRootVolume    *RootVolumeBuilder
	infra                int
	infraMachineType     *MachineTypeBuilder
	master               int
	masterMachineType    *MachineTypeBuilder
	securityGroupFilters []*MachinePoolSecurityGroupFilterBuilder
	total                int
}

// NewClusterNodes creates a new builder of 'cluster_nodes' objects.
func NewClusterNodes() *ClusterNodesBuilder {
	return &ClusterNodesBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterNodesBuilder) Empty() bool {
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

// AutoscaleCompute sets the value of the 'autoscale_compute' attribute to the given value.
//
// Representation of a autoscaling in a machine pool.
func (b *ClusterNodesBuilder) AutoscaleCompute(value *MachinePoolAutoscalingBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.autoscaleCompute = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// AvailabilityZones sets the value of the 'availability_zones' attribute to the given values.
func (b *ClusterNodesBuilder) AvailabilityZones(values ...string) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.availabilityZones = make([]string, len(values))
	copy(b.availabilityZones, values)
	b.fieldSet_[1] = true
	return b
}

// Compute sets the value of the 'compute' attribute to the given value.
func (b *ClusterNodesBuilder) Compute(value int) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.compute = value
	b.fieldSet_[2] = true
	return b
}

// ComputeLabels sets the value of the 'compute_labels' attribute to the given value.
func (b *ClusterNodesBuilder) ComputeLabels(value map[string]string) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.computeLabels = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// ComputeMachineType sets the value of the 'compute_machine_type' attribute to the given value.
//
// Machine type.
func (b *ClusterNodesBuilder) ComputeMachineType(value *MachineTypeBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.computeMachineType = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// ComputeRootVolume sets the value of the 'compute_root_volume' attribute to the given value.
//
// Root volume capabilities.
func (b *ClusterNodesBuilder) ComputeRootVolume(value *RootVolumeBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.computeRootVolume = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Infra sets the value of the 'infra' attribute to the given value.
func (b *ClusterNodesBuilder) Infra(value int) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.infra = value
	b.fieldSet_[6] = true
	return b
}

// InfraMachineType sets the value of the 'infra_machine_type' attribute to the given value.
//
// Machine type.
func (b *ClusterNodesBuilder) InfraMachineType(value *MachineTypeBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.infraMachineType = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Master sets the value of the 'master' attribute to the given value.
func (b *ClusterNodesBuilder) Master(value int) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.master = value
	b.fieldSet_[8] = true
	return b
}

// MasterMachineType sets the value of the 'master_machine_type' attribute to the given value.
//
// Machine type.
func (b *ClusterNodesBuilder) MasterMachineType(value *MachineTypeBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.masterMachineType = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// SecurityGroupFilters sets the value of the 'security_group_filters' attribute to the given values.
func (b *ClusterNodesBuilder) SecurityGroupFilters(values ...*MachinePoolSecurityGroupFilterBuilder) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.securityGroupFilters = make([]*MachinePoolSecurityGroupFilterBuilder, len(values))
	copy(b.securityGroupFilters, values)
	b.fieldSet_[10] = true
	return b
}

// Total sets the value of the 'total' attribute to the given value.
func (b *ClusterNodesBuilder) Total(value int) *ClusterNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.total = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterNodesBuilder) Copy(object *ClusterNodes) *ClusterNodesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.autoscaleCompute != nil {
		b.autoscaleCompute = NewMachinePoolAutoscaling().Copy(object.autoscaleCompute)
	} else {
		b.autoscaleCompute = nil
	}
	if object.availabilityZones != nil {
		b.availabilityZones = make([]string, len(object.availabilityZones))
		copy(b.availabilityZones, object.availabilityZones)
	} else {
		b.availabilityZones = nil
	}
	b.compute = object.compute
	if len(object.computeLabels) > 0 {
		b.computeLabels = map[string]string{}
		for k, v := range object.computeLabels {
			b.computeLabels[k] = v
		}
	} else {
		b.computeLabels = nil
	}
	if object.computeMachineType != nil {
		b.computeMachineType = NewMachineType().Copy(object.computeMachineType)
	} else {
		b.computeMachineType = nil
	}
	if object.computeRootVolume != nil {
		b.computeRootVolume = NewRootVolume().Copy(object.computeRootVolume)
	} else {
		b.computeRootVolume = nil
	}
	b.infra = object.infra
	if object.infraMachineType != nil {
		b.infraMachineType = NewMachineType().Copy(object.infraMachineType)
	} else {
		b.infraMachineType = nil
	}
	b.master = object.master
	if object.masterMachineType != nil {
		b.masterMachineType = NewMachineType().Copy(object.masterMachineType)
	} else {
		b.masterMachineType = nil
	}
	if object.securityGroupFilters != nil {
		b.securityGroupFilters = make([]*MachinePoolSecurityGroupFilterBuilder, len(object.securityGroupFilters))
		for i, v := range object.securityGroupFilters {
			b.securityGroupFilters[i] = NewMachinePoolSecurityGroupFilter().Copy(v)
		}
	} else {
		b.securityGroupFilters = nil
	}
	b.total = object.total
	return b
}

// Build creates a 'cluster_nodes' object using the configuration stored in the builder.
func (b *ClusterNodesBuilder) Build() (object *ClusterNodes, err error) {
	object = new(ClusterNodes)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.autoscaleCompute != nil {
		object.autoscaleCompute, err = b.autoscaleCompute.Build()
		if err != nil {
			return
		}
	}
	if b.availabilityZones != nil {
		object.availabilityZones = make([]string, len(b.availabilityZones))
		copy(object.availabilityZones, b.availabilityZones)
	}
	object.compute = b.compute
	if b.computeLabels != nil {
		object.computeLabels = make(map[string]string)
		for k, v := range b.computeLabels {
			object.computeLabels[k] = v
		}
	}
	if b.computeMachineType != nil {
		object.computeMachineType, err = b.computeMachineType.Build()
		if err != nil {
			return
		}
	}
	if b.computeRootVolume != nil {
		object.computeRootVolume, err = b.computeRootVolume.Build()
		if err != nil {
			return
		}
	}
	object.infra = b.infra
	if b.infraMachineType != nil {
		object.infraMachineType, err = b.infraMachineType.Build()
		if err != nil {
			return
		}
	}
	object.master = b.master
	if b.masterMachineType != nil {
		object.masterMachineType, err = b.masterMachineType.Build()
		if err != nil {
			return
		}
	}
	if b.securityGroupFilters != nil {
		object.securityGroupFilters = make([]*MachinePoolSecurityGroupFilter, len(b.securityGroupFilters))
		for i, v := range b.securityGroupFilters {
			object.securityGroupFilters[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.total = b.total
	return
}
