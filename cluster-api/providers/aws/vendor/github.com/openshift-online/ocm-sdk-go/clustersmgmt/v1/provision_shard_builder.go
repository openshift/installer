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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	time "time"
)

// ProvisionShardBuilder contains the data and logic needed to build 'provision_shard' objects.
//
// Contains the properties of the provision shard, including AWS and GCP related configurations
type ProvisionShardBuilder struct {
	bitmap_                  uint32
	id                       string
	href                     string
	awsAccountOperatorConfig *ServerConfigBuilder
	awsBaseDomain            string
	gcpBaseDomain            string
	gcpProjectOperator       *ServerConfigBuilder
	cloudProvider            *CloudProviderBuilder
	creationTimestamp        time.Time
	hiveConfig               *ServerConfigBuilder
	hypershiftConfig         *ServerConfigBuilder
	lastUpdateTimestamp      time.Time
	managementCluster        string
	region                   *CloudRegionBuilder
	status                   string
}

// NewProvisionShard creates a new builder of 'provision_shard' objects.
func NewProvisionShard() *ProvisionShardBuilder {
	return &ProvisionShardBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ProvisionShardBuilder) Link(value bool) *ProvisionShardBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ProvisionShardBuilder) ID(value string) *ProvisionShardBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ProvisionShardBuilder) HREF(value string) *ProvisionShardBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProvisionShardBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AWSAccountOperatorConfig sets the value of the 'AWS_account_operator_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) AWSAccountOperatorConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	b.awsAccountOperatorConfig = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// AWSBaseDomain sets the value of the 'AWS_base_domain' attribute to the given value.
func (b *ProvisionShardBuilder) AWSBaseDomain(value string) *ProvisionShardBuilder {
	b.awsBaseDomain = value
	b.bitmap_ |= 16
	return b
}

// GCPBaseDomain sets the value of the 'GCP_base_domain' attribute to the given value.
func (b *ProvisionShardBuilder) GCPBaseDomain(value string) *ProvisionShardBuilder {
	b.gcpBaseDomain = value
	b.bitmap_ |= 32
	return b
}

// GCPProjectOperator sets the value of the 'GCP_project_operator' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) GCPProjectOperator(value *ServerConfigBuilder) *ProvisionShardBuilder {
	b.gcpProjectOperator = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *ProvisionShardBuilder) CloudProvider(value *CloudProviderBuilder) *ProvisionShardBuilder {
	b.cloudProvider = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) CreationTimestamp(value time.Time) *ProvisionShardBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 256
	return b
}

// HiveConfig sets the value of the 'hive_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) HiveConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	b.hiveConfig = value
	if value != nil {
		b.bitmap_ |= 512
	} else {
		b.bitmap_ &^= 512
	}
	return b
}

// HypershiftConfig sets the value of the 'hypershift_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) HypershiftConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	b.hypershiftConfig = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) LastUpdateTimestamp(value time.Time) *ProvisionShardBuilder {
	b.lastUpdateTimestamp = value
	b.bitmap_ |= 2048
	return b
}

// ManagementCluster sets the value of the 'management_cluster' attribute to the given value.
func (b *ProvisionShardBuilder) ManagementCluster(value string) *ProvisionShardBuilder {
	b.managementCluster = value
	b.bitmap_ |= 4096
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *ProvisionShardBuilder) Region(value *CloudRegionBuilder) *ProvisionShardBuilder {
	b.region = value
	if value != nil {
		b.bitmap_ |= 8192
	} else {
		b.bitmap_ &^= 8192
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *ProvisionShardBuilder) Status(value string) *ProvisionShardBuilder {
	b.status = value
	b.bitmap_ |= 16384
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProvisionShardBuilder) Copy(object *ProvisionShard) *ProvisionShardBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.awsAccountOperatorConfig != nil {
		b.awsAccountOperatorConfig = NewServerConfig().Copy(object.awsAccountOperatorConfig)
	} else {
		b.awsAccountOperatorConfig = nil
	}
	b.awsBaseDomain = object.awsBaseDomain
	b.gcpBaseDomain = object.gcpBaseDomain
	if object.gcpProjectOperator != nil {
		b.gcpProjectOperator = NewServerConfig().Copy(object.gcpProjectOperator)
	} else {
		b.gcpProjectOperator = nil
	}
	if object.cloudProvider != nil {
		b.cloudProvider = NewCloudProvider().Copy(object.cloudProvider)
	} else {
		b.cloudProvider = nil
	}
	b.creationTimestamp = object.creationTimestamp
	if object.hiveConfig != nil {
		b.hiveConfig = NewServerConfig().Copy(object.hiveConfig)
	} else {
		b.hiveConfig = nil
	}
	if object.hypershiftConfig != nil {
		b.hypershiftConfig = NewServerConfig().Copy(object.hypershiftConfig)
	} else {
		b.hypershiftConfig = nil
	}
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	b.managementCluster = object.managementCluster
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	b.status = object.status
	return b
}

// Build creates a 'provision_shard' object using the configuration stored in the builder.
func (b *ProvisionShardBuilder) Build() (object *ProvisionShard, err error) {
	object = new(ProvisionShard)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.awsAccountOperatorConfig != nil {
		object.awsAccountOperatorConfig, err = b.awsAccountOperatorConfig.Build()
		if err != nil {
			return
		}
	}
	object.awsBaseDomain = b.awsBaseDomain
	object.gcpBaseDomain = b.gcpBaseDomain
	if b.gcpProjectOperator != nil {
		object.gcpProjectOperator, err = b.gcpProjectOperator.Build()
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
	if b.hiveConfig != nil {
		object.hiveConfig, err = b.hiveConfig.Build()
		if err != nil {
			return
		}
	}
	if b.hypershiftConfig != nil {
		object.hypershiftConfig, err = b.hypershiftConfig.Build()
		if err != nil {
			return
		}
	}
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	object.managementCluster = b.managementCluster
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	object.status = b.status
	return
}
