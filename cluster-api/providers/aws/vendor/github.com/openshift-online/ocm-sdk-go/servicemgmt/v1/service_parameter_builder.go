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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// ServiceParameterBuilder contains the data and logic needed to build 'service_parameter' objects.
type ServiceParameterBuilder struct {
	bitmap_ uint32
	id      string
	value   string
}

// NewServiceParameter creates a new builder of 'service_parameter' objects.
func NewServiceParameter() *ServiceParameterBuilder {
	return &ServiceParameterBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceParameterBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *ServiceParameterBuilder) ID(value string) *ServiceParameterBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *ServiceParameterBuilder) Value(value string) *ServiceParameterBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceParameterBuilder) Copy(object *ServiceParameter) *ServiceParameterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.value = object.value
	return b
}

// Build creates a 'service_parameter' object using the configuration stored in the builder.
func (b *ServiceParameterBuilder) Build() (object *ServiceParameter, err error) {
	object = new(ServiceParameter)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.value = b.value
	return
}
