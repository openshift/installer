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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AlertInfoBuilder contains the data and logic needed to build 'alert_info' objects.
//
// Provides information about a single alert firing on the cluster.
type AlertInfoBuilder struct {
	bitmap_  uint32
	name     string
	severity AlertSeverity
}

// NewAlertInfo creates a new builder of 'alert_info' objects.
func NewAlertInfo() *AlertInfoBuilder {
	return &AlertInfoBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AlertInfoBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AlertInfoBuilder) Name(value string) *AlertInfoBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
//
// Severity of a cluster alert received via telemetry.
func (b *AlertInfoBuilder) Severity(value AlertSeverity) *AlertInfoBuilder {
	b.severity = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AlertInfoBuilder) Copy(object *AlertInfo) *AlertInfoBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	b.severity = object.severity
	return b
}

// Build creates a 'alert_info' object using the configuration stored in the builder.
func (b *AlertInfoBuilder) Build() (object *AlertInfo, err error) {
	object = new(AlertInfo)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.severity = b.severity
	return
}
