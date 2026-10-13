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

type WifGcpBuilder struct {
	fieldSet_              []bool
	federatedProjectId     string
	federatedProjectNumber string
	impersonatorEmail      string
	projectId              string
	projectNumber          string
	rolePrefix             string
	serviceAccounts        []*WifServiceAccountBuilder
	support                *WifSupportBuilder
	workloadIdentityPool   *WifPoolBuilder
}

// NewWifGcp creates a new builder of 'wif_gcp' objects.
func NewWifGcp() *WifGcpBuilder {
	return &WifGcpBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifGcpBuilder) Empty() bool {
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

// FederatedProjectId sets the value of the 'federated_project_id' attribute to the given value.
func (b *WifGcpBuilder) FederatedProjectId(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.federatedProjectId = value
	b.fieldSet_[0] = true
	return b
}

// FederatedProjectNumber sets the value of the 'federated_project_number' attribute to the given value.
func (b *WifGcpBuilder) FederatedProjectNumber(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.federatedProjectNumber = value
	b.fieldSet_[1] = true
	return b
}

// ImpersonatorEmail sets the value of the 'impersonator_email' attribute to the given value.
func (b *WifGcpBuilder) ImpersonatorEmail(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.impersonatorEmail = value
	b.fieldSet_[2] = true
	return b
}

// ProjectId sets the value of the 'project_id' attribute to the given value.
func (b *WifGcpBuilder) ProjectId(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.projectId = value
	b.fieldSet_[3] = true
	return b
}

// ProjectNumber sets the value of the 'project_number' attribute to the given value.
func (b *WifGcpBuilder) ProjectNumber(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.projectNumber = value
	b.fieldSet_[4] = true
	return b
}

// RolePrefix sets the value of the 'role_prefix' attribute to the given value.
func (b *WifGcpBuilder) RolePrefix(value string) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.rolePrefix = value
	b.fieldSet_[5] = true
	return b
}

// ServiceAccounts sets the value of the 'service_accounts' attribute to the given values.
func (b *WifGcpBuilder) ServiceAccounts(values ...*WifServiceAccountBuilder) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.serviceAccounts = make([]*WifServiceAccountBuilder, len(values))
	copy(b.serviceAccounts, values)
	b.fieldSet_[6] = true
	return b
}

// Support sets the value of the 'support' attribute to the given value.
func (b *WifGcpBuilder) Support(value *WifSupportBuilder) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.support = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// WorkloadIdentityPool sets the value of the 'workload_identity_pool' attribute to the given value.
func (b *WifGcpBuilder) WorkloadIdentityPool(value *WifPoolBuilder) *WifGcpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.workloadIdentityPool = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifGcpBuilder) Copy(object *WifGcp) *WifGcpBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.federatedProjectId = object.federatedProjectId
	b.federatedProjectNumber = object.federatedProjectNumber
	b.impersonatorEmail = object.impersonatorEmail
	b.projectId = object.projectId
	b.projectNumber = object.projectNumber
	b.rolePrefix = object.rolePrefix
	if object.serviceAccounts != nil {
		b.serviceAccounts = make([]*WifServiceAccountBuilder, len(object.serviceAccounts))
		for i, v := range object.serviceAccounts {
			b.serviceAccounts[i] = NewWifServiceAccount().Copy(v)
		}
	} else {
		b.serviceAccounts = nil
	}
	if object.support != nil {
		b.support = NewWifSupport().Copy(object.support)
	} else {
		b.support = nil
	}
	if object.workloadIdentityPool != nil {
		b.workloadIdentityPool = NewWifPool().Copy(object.workloadIdentityPool)
	} else {
		b.workloadIdentityPool = nil
	}
	return b
}

// Build creates a 'wif_gcp' object using the configuration stored in the builder.
func (b *WifGcpBuilder) Build() (object *WifGcp, err error) {
	object = new(WifGcp)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.federatedProjectId = b.federatedProjectId
	object.federatedProjectNumber = b.federatedProjectNumber
	object.impersonatorEmail = b.impersonatorEmail
	object.projectId = b.projectId
	object.projectNumber = b.projectNumber
	object.rolePrefix = b.rolePrefix
	if b.serviceAccounts != nil {
		object.serviceAccounts = make([]*WifServiceAccount, len(b.serviceAccounts))
		for i, v := range b.serviceAccounts {
			object.serviceAccounts[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.support != nil {
		object.support, err = b.support.Build()
		if err != nil {
			return
		}
	}
	if b.workloadIdentityPool != nil {
		object.workloadIdentityPool, err = b.workloadIdentityPool.Build()
		if err != nil {
			return
		}
	}
	return
}
