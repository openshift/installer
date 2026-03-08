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

// CloudWatch configuration for log forwarding.
type LogForwarderCloudWatchConfigBuilder struct {
	fieldSet_              []bool
	logDistributionRoleArn string
	logGroupName           string
}

// NewLogForwarderCloudWatchConfig creates a new builder of 'log_forwarder_cloud_watch_config' objects.
func NewLogForwarderCloudWatchConfig() *LogForwarderCloudWatchConfigBuilder {
	return &LogForwarderCloudWatchConfigBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogForwarderCloudWatchConfigBuilder) Empty() bool {
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

// LogDistributionRoleArn sets the value of the 'log_distribution_role_arn' attribute to the given value.
func (b *LogForwarderCloudWatchConfigBuilder) LogDistributionRoleArn(value string) *LogForwarderCloudWatchConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.logDistributionRoleArn = value
	b.fieldSet_[0] = true
	return b
}

// LogGroupName sets the value of the 'log_group_name' attribute to the given value.
func (b *LogForwarderCloudWatchConfigBuilder) LogGroupName(value string) *LogForwarderCloudWatchConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.logGroupName = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogForwarderCloudWatchConfigBuilder) Copy(object *LogForwarderCloudWatchConfig) *LogForwarderCloudWatchConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.logDistributionRoleArn = object.logDistributionRoleArn
	b.logGroupName = object.logGroupName
	return b
}

// Build creates a 'log_forwarder_cloud_watch_config' object using the configuration stored in the builder.
func (b *LogForwarderCloudWatchConfigBuilder) Build() (object *LogForwarderCloudWatchConfig, err error) {
	object = new(LogForwarderCloudWatchConfig)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.logDistributionRoleArn = b.logDistributionRoleArn
	object.logGroupName = b.logGroupName
	return
}
