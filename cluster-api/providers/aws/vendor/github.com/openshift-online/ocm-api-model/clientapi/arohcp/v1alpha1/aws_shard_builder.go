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

// Config for AWS provision shards
type AWSShardBuilder struct {
	fieldSet_         []bool
	ecrRepositoryURLs []string
	backupConfigs     []*AWSBackupConfigBuilder
}

// NewAWSShard creates a new builder of 'AWS_shard' objects.
func NewAWSShard() *AWSShardBuilder {
	return &AWSShardBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSShardBuilder) Empty() bool {
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

// ECRRepositoryURLs sets the value of the 'ECR_repository_URLs' attribute to the given values.
func (b *AWSShardBuilder) ECRRepositoryURLs(values ...string) *AWSShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.ecrRepositoryURLs = make([]string, len(values))
	copy(b.ecrRepositoryURLs, values)
	b.fieldSet_[0] = true
	return b
}

// BackupConfigs sets the value of the 'backup_configs' attribute to the given values.
func (b *AWSShardBuilder) BackupConfigs(values ...*AWSBackupConfigBuilder) *AWSShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.backupConfigs = make([]*AWSBackupConfigBuilder, len(values))
	copy(b.backupConfigs, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSShardBuilder) Copy(object *AWSShard) *AWSShardBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.ecrRepositoryURLs != nil {
		b.ecrRepositoryURLs = make([]string, len(object.ecrRepositoryURLs))
		copy(b.ecrRepositoryURLs, object.ecrRepositoryURLs)
	} else {
		b.ecrRepositoryURLs = nil
	}
	if object.backupConfigs != nil {
		b.backupConfigs = make([]*AWSBackupConfigBuilder, len(object.backupConfigs))
		for i, v := range object.backupConfigs {
			b.backupConfigs[i] = NewAWSBackupConfig().Copy(v)
		}
	} else {
		b.backupConfigs = nil
	}
	return b
}

// Build creates a 'AWS_shard' object using the configuration stored in the builder.
func (b *AWSShardBuilder) Build() (object *AWSShard, err error) {
	object = new(AWSShard)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.ecrRepositoryURLs != nil {
		object.ecrRepositoryURLs = make([]string, len(b.ecrRepositoryURLs))
		copy(object.ecrRepositoryURLs, b.ecrRepositoryURLs)
	}
	if b.backupConfigs != nil {
		object.backupConfigs = make([]*AWSBackupConfig, len(b.backupConfigs))
		for i, v := range b.backupConfigs {
			object.backupConfigs[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
