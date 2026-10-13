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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Security Group Filter object, containing name of the filter tag and value of the filter tag
type MachinePoolSecurityGroupFilterBuilder struct {
	fieldSet_ []bool
	name      string
	value     string
}

// NewMachinePoolSecurityGroupFilter creates a new builder of 'machine_pool_security_group_filter' objects.
func NewMachinePoolSecurityGroupFilter() *MachinePoolSecurityGroupFilterBuilder {
	return &MachinePoolSecurityGroupFilterBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *MachinePoolSecurityGroupFilterBuilder) Empty() bool {
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

// Name sets the value of the 'name' attribute to the given value.
func (b *MachinePoolSecurityGroupFilterBuilder) Name(value string) *MachinePoolSecurityGroupFilterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.name = value
	b.fieldSet_[0] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *MachinePoolSecurityGroupFilterBuilder) Value(value string) *MachinePoolSecurityGroupFilterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.value = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *MachinePoolSecurityGroupFilterBuilder) Copy(object *MachinePoolSecurityGroupFilter) *MachinePoolSecurityGroupFilterBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'machine_pool_security_group_filter' object using the configuration stored in the builder.
func (b *MachinePoolSecurityGroupFilterBuilder) Build() (object *MachinePoolSecurityGroupFilter, err error) {
	object = new(MachinePoolSecurityGroupFilter)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	object.value = b.value
	return
}
