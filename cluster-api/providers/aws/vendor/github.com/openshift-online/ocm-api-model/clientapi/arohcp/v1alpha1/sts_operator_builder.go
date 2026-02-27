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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// Representation of an sts operator
type STSOperatorBuilder struct {
	fieldSet_       []bool
	maxVersion      string
	minVersion      string
	name            string
	namespace       string
	serviceAccounts []string
}

// NewSTSOperator creates a new builder of 'STS_operator' objects.
func NewSTSOperator() *STSOperatorBuilder {
	return &STSOperatorBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *STSOperatorBuilder) Empty() bool {
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

// MaxVersion sets the value of the 'max_version' attribute to the given value.
func (b *STSOperatorBuilder) MaxVersion(value string) *STSOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.maxVersion = value
	b.fieldSet_[0] = true
	return b
}

// MinVersion sets the value of the 'min_version' attribute to the given value.
func (b *STSOperatorBuilder) MinVersion(value string) *STSOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.minVersion = value
	b.fieldSet_[1] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *STSOperatorBuilder) Name(value string) *STSOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.name = value
	b.fieldSet_[2] = true
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *STSOperatorBuilder) Namespace(value string) *STSOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.namespace = value
	b.fieldSet_[3] = true
	return b
}

// ServiceAccounts sets the value of the 'service_accounts' attribute to the given values.
func (b *STSOperatorBuilder) ServiceAccounts(values ...string) *STSOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.serviceAccounts = make([]string, len(values))
	copy(b.serviceAccounts, values)
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSOperatorBuilder) Copy(object *STSOperator) *STSOperatorBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.maxVersion = object.maxVersion
	b.minVersion = object.minVersion
	b.name = object.name
	b.namespace = object.namespace
	if object.serviceAccounts != nil {
		b.serviceAccounts = make([]string, len(object.serviceAccounts))
		copy(b.serviceAccounts, object.serviceAccounts)
	} else {
		b.serviceAccounts = nil
	}
	return b
}

// Build creates a 'STS_operator' object using the configuration stored in the builder.
func (b *STSOperatorBuilder) Build() (object *STSOperator, err error) {
	object = new(STSOperator)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.maxVersion = b.maxVersion
	object.minVersion = b.minVersion
	object.name = b.name
	object.namespace = b.namespace
	if b.serviceAccounts != nil {
		object.serviceAccounts = make([]string, len(b.serviceAccounts))
		copy(object.serviceAccounts, b.serviceAccounts)
	}
	return
}
