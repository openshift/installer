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

// MachinePoolAutoscalingBuilder contains the data and logic needed to build 'machine_pool_autoscaling' objects.
//
// Representation of a autoscaling in a machine pool.
type MachinePoolAutoscalingBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	maxReplicas int
	minReplicas int
}

// NewMachinePoolAutoscaling creates a new builder of 'machine_pool_autoscaling' objects.
func NewMachinePoolAutoscaling() *MachinePoolAutoscalingBuilder {
	return &MachinePoolAutoscalingBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *MachinePoolAutoscalingBuilder) Link(value bool) *MachinePoolAutoscalingBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *MachinePoolAutoscalingBuilder) ID(value string) *MachinePoolAutoscalingBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *MachinePoolAutoscalingBuilder) HREF(value string) *MachinePoolAutoscalingBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MachinePoolAutoscalingBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// MaxReplicas sets the value of the 'max_replicas' attribute to the given value.
func (b *MachinePoolAutoscalingBuilder) MaxReplicas(value int) *MachinePoolAutoscalingBuilder {
	b.maxReplicas = value
	b.bitmap_ |= 8
	return b
}

// MinReplicas sets the value of the 'min_replicas' attribute to the given value.
func (b *MachinePoolAutoscalingBuilder) MinReplicas(value int) *MachinePoolAutoscalingBuilder {
	b.minReplicas = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MachinePoolAutoscalingBuilder) Copy(object *MachinePoolAutoscaling) *MachinePoolAutoscalingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.maxReplicas = object.maxReplicas
	b.minReplicas = object.minReplicas
	return b
}

// Build creates a 'machine_pool_autoscaling' object using the configuration stored in the builder.
func (b *MachinePoolAutoscalingBuilder) Build() (object *MachinePoolAutoscaling, err error) {
	object = new(MachinePoolAutoscaling)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.maxReplicas = b.maxReplicas
	object.minReplicas = b.minReplicas
	return
}
