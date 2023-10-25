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

// MonitoringStackResourcesBuilder contains the data and logic needed to build 'monitoring_stack_resources' objects.
//
// Representation of Monitoring Stack Resources
type MonitoringStackResourcesBuilder struct {
	bitmap_  uint32
	limits   *MonitoringStackResourceBuilder
	requests *MonitoringStackResourceBuilder
}

// NewMonitoringStackResources creates a new builder of 'monitoring_stack_resources' objects.
func NewMonitoringStackResources() *MonitoringStackResourcesBuilder {
	return &MonitoringStackResourcesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MonitoringStackResourcesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Limits sets the value of the 'limits' attribute to the given value.
//
// Representation of Monitoring Stack Resource
func (b *MonitoringStackResourcesBuilder) Limits(value *MonitoringStackResourceBuilder) *MonitoringStackResourcesBuilder {
	b.limits = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Requests sets the value of the 'requests' attribute to the given value.
//
// Representation of Monitoring Stack Resource
func (b *MonitoringStackResourcesBuilder) Requests(value *MonitoringStackResourceBuilder) *MonitoringStackResourcesBuilder {
	b.requests = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MonitoringStackResourcesBuilder) Copy(object *MonitoringStackResources) *MonitoringStackResourcesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
