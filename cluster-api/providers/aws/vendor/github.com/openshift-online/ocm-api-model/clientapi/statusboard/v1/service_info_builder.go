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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

// Definition of a Status Board service info.
type ServiceInfoBuilder struct {
	fieldSet_  []bool
	fullname   string
	statusType string
}

// NewServiceInfo creates a new builder of 'service_info' objects.
func NewServiceInfo() *ServiceInfoBuilder {
	return &ServiceInfoBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceInfoBuilder) Empty() bool {
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

// Fullname sets the value of the 'fullname' attribute to the given value.
func (b *ServiceInfoBuilder) Fullname(value string) *ServiceInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.fullname = value
	b.fieldSet_[0] = true
	return b
}

// StatusType sets the value of the 'status_type' attribute to the given value.
func (b *ServiceInfoBuilder) StatusType(value string) *ServiceInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.statusType = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceInfoBuilder) Copy(object *ServiceInfo) *ServiceInfoBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.fullname = object.fullname
	b.statusType = object.statusType
	return b
}

// Build creates a 'service_info' object using the configuration stored in the builder.
func (b *ServiceInfoBuilder) Build() (object *ServiceInfo, err error) {
	object = new(ServiceInfo)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.fullname = b.fullname
	object.statusType = b.statusType
	return
}
