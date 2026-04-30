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

import (
	time "time"
)

// Contains the properties of the provision shard
type ProvisionShardBuilder struct {
	fieldSet_           []bool
	id                  string
	href                string
	azureShard          *AzureShardBuilder
	cloudProvider       *CloudProviderBuilder
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	maestroConfig       *ProvisionShardMaestroConfigBuilder
	region              *CloudRegionBuilder
	status              string
	topology            string
}

// NewProvisionShard creates a new builder of 'provision_shard' objects.
func NewProvisionShard() *ProvisionShardBuilder {
	return &ProvisionShardBuilder{
		fieldSet_: make([]bool, 11),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ProvisionShardBuilder) Link(value bool) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ProvisionShardBuilder) ID(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ProvisionShardBuilder) HREF(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProvisionShardBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AzureShard sets the value of the 'azure_shard' attribute to the given value.
//
// The Azure related configuration of the Provision Shard
func (b *ProvisionShardBuilder) AzureShard(value *AzureShardBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.azureShard = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *ProvisionShardBuilder) CloudProvider(value *CloudProviderBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.cloudProvider = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) CreationTimestamp(value time.Time) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.creationTimestamp = value
	b.fieldSet_[5] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) LastUpdateTimestamp(value time.Time) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// MaestroConfig sets the value of the 'maestro_config' attribute to the given value.
//
// The Maestro related configuration of the Provision Shard.
// The combination of `consumer_name` and `rest_api_config.url`
// must be unique across shards.
// The combination of `consumer_name` and `grpc_api_config.url`
// must be unique across shards.
func (b *ProvisionShardBuilder) MaestroConfig(value *ProvisionShardMaestroConfigBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.maestroConfig = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *ProvisionShardBuilder) Region(value *CloudRegionBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.region = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *ProvisionShardBuilder) Status(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.status = value
	b.fieldSet_[9] = true
	return b
}

// Topology sets the value of the 'topology' attribute to the given value.
func (b *ProvisionShardBuilder) Topology(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.topology = value
	b.fieldSet_[10] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProvisionShardBuilder) Copy(object *ProvisionShard) *ProvisionShardBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.azureShard != nil {
		b.azureShard = NewAzureShard().Copy(object.azureShard)
	} else {
		b.azureShard = nil
	}
	if object.cloudProvider != nil {
		b.cloudProvider = NewCloudProvider().Copy(object.cloudProvider)
	} else {
		b.cloudProvider = nil
	}
	b.creationTimestamp = object.creationTimestamp
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	if object.maestroConfig != nil {
		b.maestroConfig = NewProvisionShardMaestroConfig().Copy(object.maestroConfig)
	} else {
		b.maestroConfig = nil
	}
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	b.status = object.status
	b.topology = object.topology
	return b
}

// Build creates a 'provision_shard' object using the configuration stored in the builder.
func (b *ProvisionShardBuilder) Build() (object *ProvisionShard, err error) {
	object = new(ProvisionShard)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.azureShard != nil {
		object.azureShard, err = b.azureShard.Build()
		if err != nil {
			return
		}
	}
	if b.cloudProvider != nil {
		object.cloudProvider, err = b.cloudProvider.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	if b.maestroConfig != nil {
		object.maestroConfig, err = b.maestroConfig.Build()
		if err != nil {
			return
		}
	}
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	object.status = b.status
	object.topology = b.topology
	return
}
