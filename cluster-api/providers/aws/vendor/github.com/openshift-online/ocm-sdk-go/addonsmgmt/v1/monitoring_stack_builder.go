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

// MonitoringStackBuilder contains the data and logic needed to build 'monitoring_stack' objects.
//
// Representation of Monitoring Stack
type MonitoringStackBuilder struct {
	bitmap_   uint32
	resources *MonitoringStackResourcesBuilder
	enabled   bool
}

// NewMonitoringStack creates a new builder of 'monitoring_stack' objects.
func NewMonitoringStack() *MonitoringStackBuilder {
	return &MonitoringStackBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *MonitoringStackBuilder) Enabled(value bool) *MonitoringStackBuilder {
	b.enabled = value
	b.bitmap_ |= 1
	return b
}

// Resources sets the value of the 'resources' attribute to the given value.
//
// Representation of Monitoring Stack Resources
func (b *MonitoringStackBuilder) Resources(value *MonitoringStackResourcesBuilder) *MonitoringStackBuilder {
	b.resources = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackBuilder) Copy(object *MonitoringStack) *MonitoringStackBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.enabled = object.enabled
	if object.resources != nil {
		b.resources = NewMonitoringStackResources().Copy(object.resources)
	} else {
		b.resources = nil
	}
	return b
}

// Build creates a 'monitoring_stack' object using the configuration stored in the builder.
func (b *MonitoringStackBuilder) Build() (object *MonitoringStack, err error) {
	object = new(MonitoringStack)
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	if b.resources != nil {
		object.resources, err = b.resources.Build()
		if err != nil {
			return
		}
	}
	return
}
