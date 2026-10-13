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

// Backup configuration for AWS clusters
type AWSBackupConfigBuilder struct {
	fieldSet_           []bool
	s3Bucket            string
	accountId           string
	identityProviderArn string
	managementCluster   string
	roleArn             string
}

// NewAWSBackupConfig creates a new builder of 'AWS_backup_config' objects.
func NewAWSBackupConfig() *AWSBackupConfigBuilder {
	return &AWSBackupConfigBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSBackupConfigBuilder) Empty() bool {
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

// S3Bucket sets the value of the 'S3_bucket' attribute to the given value.
func (b *AWSBackupConfigBuilder) S3Bucket(value string) *AWSBackupConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.s3Bucket = value
	b.fieldSet_[0] = true
	return b
}

// AccountId sets the value of the 'account_id' attribute to the given value.
func (b *AWSBackupConfigBuilder) AccountId(value string) *AWSBackupConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.accountId = value
	b.fieldSet_[1] = true
	return b
}

// IdentityProviderArn sets the value of the 'identity_provider_arn' attribute to the given value.
func (b *AWSBackupConfigBuilder) IdentityProviderArn(value string) *AWSBackupConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.identityProviderArn = value
	b.fieldSet_[2] = true
	return b
}

// ManagementCluster sets the value of the 'management_cluster' attribute to the given value.
func (b *AWSBackupConfigBuilder) ManagementCluster(value string) *AWSBackupConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.managementCluster = value
	b.fieldSet_[3] = true
	return b
}

// RoleArn sets the value of the 'role_arn' attribute to the given value.
func (b *AWSBackupConfigBuilder) RoleArn(value string) *AWSBackupConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.roleArn = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSBackupConfigBuilder) Copy(object *AWSBackupConfig) *AWSBackupConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.s3Bucket = object.s3Bucket
	b.accountId = object.accountId
	b.identityProviderArn = object.identityProviderArn
	b.managementCluster = object.managementCluster
	b.roleArn = object.roleArn
	return b
}

// Build creates a 'AWS_backup_config' object using the configuration stored in the builder.
func (b *AWSBackupConfigBuilder) Build() (object *AWSBackupConfig, err error) {
	object = new(AWSBackupConfig)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.s3Bucket = b.s3Bucket
	object.accountId = b.accountId
	object.identityProviderArn = b.identityProviderArn
	object.managementCluster = b.managementCluster
	object.roleArn = b.roleArn
	return
}
