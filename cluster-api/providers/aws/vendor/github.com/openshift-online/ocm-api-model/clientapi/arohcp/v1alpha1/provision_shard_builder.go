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

	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// Contains the properties of the provision shard, including AWS and GCP related configurations
type ProvisionShardBuilder struct {
	fieldSet_                []bool
	id                       string
	href                     string
	awsAccountOperatorConfig *ServerConfigBuilder
	awsBaseDomain            string
	gcpBaseDomain            string
	gcpProjectOperator       *ServerConfigBuilder
	cloudProvider            *v1.CloudProviderBuilder
	creationTimestamp        time.Time
	hiveConfig               *ServerConfigBuilder
	hypershiftConfig         *ServerConfigBuilder
	lastUpdateTimestamp      time.Time
	managementCluster        string
	region                   *v1.CloudRegionBuilder
	status                   string
}

// NewProvisionShard creates a new builder of 'provision_shard' objects.
func NewProvisionShard() *ProvisionShardBuilder {
	return &ProvisionShardBuilder{
		fieldSet_: make([]bool, 15),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ProvisionShardBuilder) Link(value bool) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ProvisionShardBuilder) ID(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ProvisionShardBuilder) HREF(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
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

// AWSAccountOperatorConfig sets the value of the 'AWS_account_operator_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) AWSAccountOperatorConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.awsAccountOperatorConfig = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// AWSBaseDomain sets the value of the 'AWS_base_domain' attribute to the given value.
func (b *ProvisionShardBuilder) AWSBaseDomain(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.awsBaseDomain = value
	b.fieldSet_[4] = true
	return b
}

// GCPBaseDomain sets the value of the 'GCP_base_domain' attribute to the given value.
func (b *ProvisionShardBuilder) GCPBaseDomain(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.gcpBaseDomain = value
	b.fieldSet_[5] = true
	return b
}

// GCPProjectOperator sets the value of the 'GCP_project_operator' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) GCPProjectOperator(value *ServerConfigBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.gcpProjectOperator = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *ProvisionShardBuilder) CloudProvider(value *v1.CloudProviderBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.cloudProvider = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) CreationTimestamp(value time.Time) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.creationTimestamp = value
	b.fieldSet_[8] = true
	return b
}

// HiveConfig sets the value of the 'hive_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) HiveConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.hiveConfig = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// HypershiftConfig sets the value of the 'hypershift_config' attribute to the given value.
//
// Representation of a server config
func (b *ProvisionShardBuilder) HypershiftConfig(value *ServerConfigBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.hypershiftConfig = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ProvisionShardBuilder) LastUpdateTimestamp(value time.Time) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[11] = true
	return b
}

// ManagementCluster sets the value of the 'management_cluster' attribute to the given value.
func (b *ProvisionShardBuilder) ManagementCluster(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.managementCluster = value
	b.fieldSet_[12] = true
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *ProvisionShardBuilder) Region(value *v1.CloudRegionBuilder) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.region = value
	if value != nil {
		b.fieldSet_[13] = true
	} else {
		b.fieldSet_[13] = false
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *ProvisionShardBuilder) Status(value string) *ProvisionShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 15)
	}
	b.status = value
	b.fieldSet_[14] = true
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
		b.cloudProvider = v1.NewCloudProvider().Copy(object.cloudProvider)
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
		b.region = v1.NewCloudRegion().Copy(object.region)
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
