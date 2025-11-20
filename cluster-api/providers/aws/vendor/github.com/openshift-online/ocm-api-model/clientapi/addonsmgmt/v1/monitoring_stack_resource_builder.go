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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// Representation of Monitoring Stack Resource
type MonitoringStackResourceBuilder struct {
	fieldSet_ []bool
	cpu       string
	memory    string
}

// NewMonitoringStackResource creates a new builder of 'monitoring_stack_resource' objects.
func NewMonitoringStackResource() *MonitoringStackResourceBuilder {
	return &MonitoringStackResourceBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackResourceBuilder) Empty() bool {
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

// Cpu sets the value of the 'cpu' attribute to the given value.
func (b *MonitoringStackResourceBuilder) Cpu(value string) *MonitoringStackResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.cpu = value
	b.fieldSet_[0] = true
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *MonitoringStackResourceBuilder) Memory(value string) *MonitoringStackResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.memory = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackResourceBuilder) Copy(object *MonitoringStackResource) *MonitoringStackResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.cpu = object.cpu
	b.memory = object.memory
	return b
}

// Build creates a 'monitoring_stack_resource' object using the configuration stored in the builder.
func (b *MonitoringStackResourceBuilder) Build() (object *MonitoringStackResource, err error) {
	object = new(MonitoringStackResource)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.cpu = b.cpu
	object.memory = b.memory
	return
}
