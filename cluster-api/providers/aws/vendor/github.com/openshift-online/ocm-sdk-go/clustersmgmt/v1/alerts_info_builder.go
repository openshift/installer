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

// AlertsInfoBuilder contains the data and logic needed to build 'alerts_info' objects.
//
// Provides information about the alerts firing on the cluster.
type AlertsInfoBuilder struct {
	bitmap_ uint32
	alerts  []*AlertInfoBuilder
}

// NewAlertsInfo creates a new builder of 'alerts_info' objects.
func NewAlertsInfo() *AlertsInfoBuilder {
	return &AlertsInfoBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AlertsInfoBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Alerts sets the value of the 'alerts' attribute to the given values.
func (b *AlertsInfoBuilder) Alerts(values ...*AlertInfoBuilder) *AlertsInfoBuilder {
	b.alerts = make([]*AlertInfoBuilder, len(values))
	copy(b.alerts, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AlertsInfoBuilder) Copy(object *AlertsInfo) *AlertsInfoBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.alerts != nil {
		b.alerts = make([]*AlertInfoBuilder, len(object.alerts))
		for i, v := range object.alerts {
			b.alerts[i] = NewAlertInfo().Copy(v)
		}
	} else {
		b.alerts = nil
	}
	return b
}

// Build creates a 'alerts_info' object using the configuration stored in the builder.
func (b *AlertsInfoBuilder) Build() (object *AlertsInfo, err error) {
	object = new(AlertsInfo)
	object.bitmap_ = b.bitmap_
	if b.alerts != nil {
		object.alerts = make([]*AlertInfo, len(b.alerts))
		for i, v := range b.alerts {
			object.alerts[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
