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

// Representation of Monitoring Stack Resources
type MonitoringStackResourcesBuilder struct {
	fieldSet_ []bool
	limits    *MonitoringStackResourceBuilder
	requests  *MonitoringStackResourceBuilder
}

// NewMonitoringStackResources creates a new builder of 'monitoring_stack_resources' objects.
func NewMonitoringStackResources() *MonitoringStackResourcesBuilder {
	return &MonitoringStackResourcesBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackResourcesBuilder) Empty() bool {
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

// Limits sets the value of the 'limits' attribute to the given value.
//
// Representation of Monitoring Stack Resource
func (b *MonitoringStackResourcesBuilder) Limits(value *MonitoringStackResourceBuilder) *MonitoringStackResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.limits = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Requests sets the value of the 'requests' attribute to the given value.
//
// Representation of Monitoring Stack Resource
func (b *MonitoringStackResourcesBuilder) Requests(value *MonitoringStackResourceBuilder) *MonitoringStackResourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.requests = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackResourcesBuilder) Copy(object *MonitoringStackResources) *MonitoringStackResourcesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.limits != nil {
		b.limits = NewMonitoringStackResource().Copy(object.limits)
	} else {
		b.limits = nil
	}
	if object.requests != nil {
		b.requests = NewMonitoringStackResource().Copy(object.requests)
	} else {
		b.requests = nil
	}
	return b
}

// Build creates a 'monitoring_stack_resources' object using the configuration stored in the builder.
func (b *MonitoringStackResourcesBuilder) Build() (object *MonitoringStackResources, err error) {
	object = new(MonitoringStackResources)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.limits != nil {
		object.limits, err = b.limits.Build()
		if err != nil {
			return
		}
	}
	if b.requests != nil {
		object.requests, err = b.requests.Build()
		if err != nil {
			return
		}
	}
	return
}
