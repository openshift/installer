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

// Representation of a machine pool in a cluster.
type MachinePoolBuilder struct {
	fieldSet_            []bool
	id                   string
	href                 string
	aws                  *AWSMachinePoolBuilder
	gcp                  *GCPMachinePoolBuilder
	autoscaling          *MachinePoolAutoscalingBuilder
	availabilityZones    []string
	instanceType         string
	labels               map[string]string
	replicas             int
	rootVolume           *RootVolumeBuilder
	securityGroupFilters []*MachinePoolSecurityGroupFilterBuilder
	subnets              []string
	taints               []*TaintBuilder
}

// NewMachinePool creates a new builder of 'machine_pool' objects.
func NewMachinePool() *MachinePoolBuilder {
	return &MachinePoolBuilder{
		fieldSet_: make([]bool, 14),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *MachinePoolBuilder) Link(value bool) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *MachinePoolBuilder) ID(value string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *MachinePoolBuilder) HREF(value string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MachinePoolBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// Representation of aws machine pool specific parameters.
func (b *MachinePoolBuilder) AWS(value *AWSMachinePoolBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.aws = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Representation of gcp machine pool specific parameters.
func (b *MachinePoolBuilder) GCP(value *GCPMachinePoolBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.gcp = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Autoscaling sets the value of the 'autoscaling' attribute to the given value.
//
// Representation of a autoscaling in a machine pool.
func (b *MachinePoolBuilder) Autoscaling(value *MachinePoolAutoscalingBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.autoscaling = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// AvailabilityZones sets the value of the 'availability_zones' attribute to the given values.
func (b *MachinePoolBuilder) AvailabilityZones(values ...string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.availabilityZones = make([]string, len(values))
	copy(b.availabilityZones, values)
	b.fieldSet_[6] = true
	return b
}

// InstanceType sets the value of the 'instance_type' attribute to the given value.
func (b *MachinePoolBuilder) InstanceType(value string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.instanceType = value
	b.fieldSet_[7] = true
	return b
}

// Labels sets the value of the 'labels' attribute to the given value.
func (b *MachinePoolBuilder) Labels(value map[string]string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.labels = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Replicas sets the value of the 'replicas' attribute to the given value.
func (b *MachinePoolBuilder) Replicas(value int) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.replicas = value
	b.fieldSet_[9] = true
	return b
}

// RootVolume sets the value of the 'root_volume' attribute to the given value.
//
// Root volume capabilities.
func (b *MachinePoolBuilder) RootVolume(value *RootVolumeBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.rootVolume = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// SecurityGroupFilters sets the value of the 'security_group_filters' attribute to the given values.
func (b *MachinePoolBuilder) SecurityGroupFilters(values ...*MachinePoolSecurityGroupFilterBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.securityGroupFilters = make([]*MachinePoolSecurityGroupFilterBuilder, len(values))
	copy(b.securityGroupFilters, values)
	b.fieldSet_[11] = true
	return b
}

// Subnets sets the value of the 'subnets' attribute to the given values.
func (b *MachinePoolBuilder) Subnets(values ...string) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.subnets = make([]string, len(values))
	copy(b.subnets, values)
	b.fieldSet_[12] = true
	return b
}

// Taints sets the value of the 'taints' attribute to the given values.
func (b *MachinePoolBuilder) Taints(values ...*TaintBuilder) *MachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.taints = make([]*TaintBuilder, len(values))
	copy(b.taints, values)
	b.fieldSet_[13] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MachinePoolBuilder) Copy(object *MachinePool) *MachinePoolBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.aws != nil {
		b.aws = NewAWSMachinePool().Copy(object.aws)
	} else {
		b.aws = nil
	}
	if object.gcp != nil {
		b.gcp = NewGCPMachinePool().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	if object.autoscaling != nil {
		b.autoscaling = NewMachinePoolAutoscaling().Copy(object.autoscaling)
	} else {
		b.autoscaling = nil
	}
	if object.availabilityZones != nil {
		b.availabilityZones = make([]string, len(object.availabilityZones))
		copy(b.availabilityZones, object.availabilityZones)
	} else {
		b.availabilityZones = nil
	}
	b.instanceType = object.instanceType
	if len(object.labels) > 0 {
		b.labels = map[string]string{}
		for k, v := range object.labels {
			b.labels[k] = v
		}
	} else {
		b.labels = nil
	}
	b.replicas = object.replicas
	if object.rootVolume != nil {
		b.rootVolume = NewRootVolume().Copy(object.rootVolume)
	} else {
		b.rootVolume = nil
	}
	if object.securityGroupFilters != nil {
		b.securityGroupFilters = make([]*MachinePoolSecurityGroupFilterBuilder, len(object.securityGroupFilters))
		for i, v := range object.securityGroupFilters {
			b.securityGroupFilters[i] = NewMachinePoolSecurityGroupFilter().Copy(v)
		}
	} else {
		b.securityGroupFilters = nil
	}
	if object.subnets != nil {
		b.subnets = make([]string, len(object.subnets))
		copy(b.subnets, object.subnets)
	} else {
		b.subnets = nil
	}
	if object.taints != nil {
		b.taints = make([]*TaintBuilder, len(object.taints))
		for i, v := range object.taints {
			b.taints[i] = NewTaint().Copy(v)
		}
	} else {
		b.taints = nil
	}
	return b
}

// Build creates a 'machine_pool' object using the configuration stored in the builder.
func (b *MachinePoolBuilder) Build() (object *MachinePool, err error) {
	object = new(MachinePool)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.aws != nil {
		object.aws, err = b.aws.Build()
		if err != nil {
			return
		}
	}
	if b.gcp != nil {
		object.gcp, err = b.gcp.Build()
		if err != nil {
			return
		}
	}
	if b.autoscaling != nil {
		object.autoscaling, err = b.autoscaling.Build()
		if err != nil {
			return
		}
	}
	if b.availabilityZones != nil {
		object.availabilityZones = make([]string, len(b.availabilityZones))
		copy(object.availabilityZones, b.availabilityZones)
	}
	object.instanceType = b.instanceType
	if b.labels != nil {
		object.labels = make(map[string]string)
		for k, v := range b.labels {
			object.labels[k] = v
		}
	}
	object.replicas = b.replicas
	if b.rootVolume != nil {
		object.rootVolume, err = b.rootVolume.Build()
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
	if b.subnets != nil {
		object.subnets = make([]string, len(b.subnets))
		copy(object.subnets, b.subnets)
	}
	if b.taints != nil {
		object.taints = make([]*Taint, len(b.taints))
		for i, v := range b.taints {
			object.taints[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
