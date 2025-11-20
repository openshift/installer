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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

type ServiceClusterRequestPayloadBuilder struct {
	fieldSet_     []bool
	cloudProvider string
	labels        []*LabelRequestPayloadBuilder
	region        string
}

// NewServiceClusterRequestPayload creates a new builder of 'service_cluster_request_payload' objects.
func NewServiceClusterRequestPayload() *ServiceClusterRequestPayloadBuilder {
	return &ServiceClusterRequestPayloadBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceClusterRequestPayloadBuilder) Empty() bool {
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

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *ServiceClusterRequestPayloadBuilder) CloudProvider(value string) *ServiceClusterRequestPayloadBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.cloudProvider = value
	b.fieldSet_[0] = true
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *ServiceClusterRequestPayloadBuilder) Labels(values ...*LabelRequestPayloadBuilder) *ServiceClusterRequestPayloadBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.labels = make([]*LabelRequestPayloadBuilder, len(values))
	copy(b.labels, values)
	b.fieldSet_[1] = true
	return b
}

// Region sets the value of the 'region' attribute to the given value.
func (b *ServiceClusterRequestPayloadBuilder) Region(value string) *ServiceClusterRequestPayloadBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.region = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceClusterRequestPayloadBuilder) Copy(object *ServiceClusterRequestPayload) *ServiceClusterRequestPayloadBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.cloudProvider = object.cloudProvider
	if object.labels != nil {
		b.labels = make([]*LabelRequestPayloadBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabelRequestPayload().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.region = object.region
	return b
}

// Build creates a 'service_cluster_request_payload' object using the configuration stored in the builder.
func (b *ServiceClusterRequestPayloadBuilder) Build() (object *ServiceClusterRequestPayload, err error) {
	object = new(ServiceClusterRequestPayload)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.cloudProvider = b.cloudProvider
	if b.labels != nil {
		object.labels = make([]*LabelRequestPayload, len(b.labels))
		for i, v := range b.labels {
			object.labels[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.region = b.region
	return
}
