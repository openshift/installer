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

type K8sServiceAccountOperatorIdentityRequirementBuilder struct {
	fieldSet_ []bool
	name      string
	namespace string
}

// NewK8sServiceAccountOperatorIdentityRequirement creates a new builder of 'K8s_service_account_operator_identity_requirement' objects.
func NewK8sServiceAccountOperatorIdentityRequirement() *K8sServiceAccountOperatorIdentityRequirementBuilder {
	return &K8sServiceAccountOperatorIdentityRequirementBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *K8sServiceAccountOperatorIdentityRequirementBuilder) Empty() bool {
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
func (b *K8sServiceAccountOperatorIdentityRequirementBuilder) Name(value string) *K8sServiceAccountOperatorIdentityRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.name = value
	b.fieldSet_[0] = true
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *K8sServiceAccountOperatorIdentityRequirementBuilder) Namespace(value string) *K8sServiceAccountOperatorIdentityRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.namespace = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *K8sServiceAccountOperatorIdentityRequirementBuilder) Copy(object *K8sServiceAccountOperatorIdentityRequirement) *K8sServiceAccountOperatorIdentityRequirementBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.name = object.name
	b.namespace = object.namespace
	return b
}

// Build creates a 'K8s_service_account_operator_identity_requirement' object using the configuration stored in the builder.
func (b *K8sServiceAccountOperatorIdentityRequirementBuilder) Build() (object *K8sServiceAccountOperatorIdentityRequirement, err error) {
	object = new(K8sServiceAccountOperatorIdentityRequirement)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	object.namespace = b.namespace
	return
}
