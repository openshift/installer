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

type WifPoolBuilder struct {
	fieldSet_        []bool
	identityProvider *WifIdentityProviderBuilder
	poolId           string
	poolName         string
}

// NewWifPool creates a new builder of 'wif_pool' objects.
func NewWifPool() *WifPoolBuilder {
	return &WifPoolBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifPoolBuilder) Empty() bool {
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

// IdentityProvider sets the value of the 'identity_provider' attribute to the given value.
func (b *WifPoolBuilder) IdentityProvider(value *WifIdentityProviderBuilder) *WifPoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.identityProvider = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// PoolId sets the value of the 'pool_id' attribute to the given value.
func (b *WifPoolBuilder) PoolId(value string) *WifPoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.poolId = value
	b.fieldSet_[1] = true
	return b
}

// PoolName sets the value of the 'pool_name' attribute to the given value.
func (b *WifPoolBuilder) PoolName(value string) *WifPoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.poolName = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifPoolBuilder) Copy(object *WifPool) *WifPoolBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.identityProvider != nil {
		b.identityProvider = NewWifIdentityProvider().Copy(object.identityProvider)
	} else {
		b.identityProvider = nil
	}
	b.poolId = object.poolId
	b.poolName = object.poolName
	return b
}

// Build creates a 'wif_pool' object using the configuration stored in the builder.
func (b *WifPoolBuilder) Build() (object *WifPool, err error) {
	object = new(WifPool)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.identityProvider != nil {
		object.identityProvider, err = b.identityProvider.Build()
		if err != nil {
			return
		}
	}
	object.poolId = b.poolId
	object.poolName = b.poolName
	return
}
