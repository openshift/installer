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

// Representation of Metrics Federation
type MetricsFederationBuilder struct {
	fieldSet_   []bool
	matchLabels map[string]string
	matchNames  []string
	namespace   string
	portName    string
}

// NewMetricsFederation creates a new builder of 'metrics_federation' objects.
func NewMetricsFederation() *MetricsFederationBuilder {
	return &MetricsFederationBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MetricsFederationBuilder) Empty() bool {
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

// MatchLabels sets the value of the 'match_labels' attribute to the given value.
func (b *MetricsFederationBuilder) MatchLabels(value map[string]string) *MetricsFederationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.matchLabels = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// MatchNames sets the value of the 'match_names' attribute to the given values.
func (b *MetricsFederationBuilder) MatchNames(values ...string) *MetricsFederationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.matchNames = make([]string, len(values))
	copy(b.matchNames, values)
	b.fieldSet_[1] = true
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *MetricsFederationBuilder) Namespace(value string) *MetricsFederationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.namespace = value
	b.fieldSet_[2] = true
	return b
}

// PortName sets the value of the 'port_name' attribute to the given value.
func (b *MetricsFederationBuilder) PortName(value string) *MetricsFederationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.portName = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MetricsFederationBuilder) Copy(object *MetricsFederation) *MetricsFederationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if len(object.matchLabels) > 0 {
		b.matchLabels = map[string]string{}
		for k, v := range object.matchLabels {
			b.matchLabels[k] = v
		}
	} else {
		b.matchLabels = nil
	}
	if object.matchNames != nil {
		b.matchNames = make([]string, len(object.matchNames))
		copy(b.matchNames, object.matchNames)
	} else {
		b.matchNames = nil
	}
	b.namespace = object.namespace
	b.portName = object.portName
	return b
}

// Build creates a 'metrics_federation' object using the configuration stored in the builder.
func (b *MetricsFederationBuilder) Build() (object *MetricsFederation, err error) {
	object = new(MetricsFederation)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.matchLabels != nil {
		object.matchLabels = make(map[string]string)
		for k, v := range b.matchLabels {
			object.matchLabels[k] = v
		}
	}
	if b.matchNames != nil {
		object.matchNames = make([]string, len(b.matchNames))
		copy(object.matchNames, b.matchNames)
	}
	object.namespace = b.namespace
	object.portName = b.portName
	return
}
