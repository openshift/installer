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

// Representation of Monitoring Stack
type MonitoringStackBuilder struct {
	fieldSet_ []bool
	resources *MonitoringStackResourcesBuilder
	enabled   bool
}

// NewMonitoringStack creates a new builder of 'monitoring_stack' objects.
func NewMonitoringStack() *MonitoringStackBuilder {
	return &MonitoringStackBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackBuilder) Empty() bool {
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

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *MonitoringStackBuilder) Enabled(value bool) *MonitoringStackBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.enabled = value
	b.fieldSet_[0] = true
	return b
}

// Resources sets the value of the 'resources' attribute to the given value.
//
// Representation of Monitoring Stack Resources
func (b *MonitoringStackBuilder) Resources(value *MonitoringStackResourcesBuilder) *MonitoringStackBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.resources = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackBuilder) Copy(object *MonitoringStack) *MonitoringStackBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	if b.resources != nil {
		object.resources, err = b.resources.Build()
		if err != nil {
			return
		}
	}
	return
}
