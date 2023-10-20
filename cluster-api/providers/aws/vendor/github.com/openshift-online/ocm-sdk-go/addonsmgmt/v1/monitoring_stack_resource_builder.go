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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// MonitoringStackResourceBuilder contains the data and logic needed to build 'monitoring_stack_resource' objects.
//
// Representation of Monitoring Stack Resource
type MonitoringStackResourceBuilder struct {
	bitmap_ uint32
	cpu     string
	memory  string
}

// NewMonitoringStackResource creates a new builder of 'monitoring_stack_resource' objects.
func NewMonitoringStackResource() *MonitoringStackResourceBuilder {
	return &MonitoringStackResourceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackResourceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Cpu sets the value of the 'cpu' attribute to the given value.
func (b *MonitoringStackResourceBuilder) Cpu(value string) *MonitoringStackResourceBuilder {
	b.cpu = value
	b.bitmap_ |= 1
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *MonitoringStackResourceBuilder) Memory(value string) *MonitoringStackResourceBuilder {
	b.memory = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackResourceBuilder) Copy(object *MonitoringStackResource) *MonitoringStackResourceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.cpu = object.cpu
	b.memory = object.memory
	return b
}

// Build creates a 'monitoring_stack_resource' object using the configuration stored in the builder.
func (b *MonitoringStackResourceBuilder) Build() (object *MonitoringStackResource, err error) {
	object = new(MonitoringStackResource)
	object.bitmap_ = b.bitmap_
	object.cpu = b.cpu
	object.memory = b.memory
	return
}
