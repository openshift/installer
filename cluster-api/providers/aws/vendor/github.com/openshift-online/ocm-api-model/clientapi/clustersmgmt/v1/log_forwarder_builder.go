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

// Representation of a log forwarder configuration for a cluster.
type LogForwarderBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	s3           *LogForwarderS3ConfigBuilder
	applications []string
	cloudwatch   *LogForwarderCloudWatchConfigBuilder
	clusterID    string
	groups       []*LogForwarderGroupBuilder
	status       *LogForwarderStatusBuilder
}

// NewLogForwarder creates a new builder of 'log_forwarder' objects.
func NewLogForwarder() *LogForwarderBuilder {
	return &LogForwarderBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *LogForwarderBuilder) Link(value bool) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *LogForwarderBuilder) ID(value string) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *LogForwarderBuilder) HREF(value string) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogForwarderBuilder) Empty() bool {
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

// S3 sets the value of the 'S3' attribute to the given value.
//
// S3 configuration for log forwarding.
func (b *LogForwarderBuilder) S3(value *LogForwarderS3ConfigBuilder) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.s3 = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Applications sets the value of the 'applications' attribute to the given values.
func (b *LogForwarderBuilder) Applications(values ...string) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.applications = make([]string, len(values))
	copy(b.applications, values)
	b.fieldSet_[4] = true
	return b
}

// Cloudwatch sets the value of the 'cloudwatch' attribute to the given value.
//
// CloudWatch configuration for log forwarding.
func (b *LogForwarderBuilder) Cloudwatch(value *LogForwarderCloudWatchConfigBuilder) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.cloudwatch = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *LogForwarderBuilder) ClusterID(value string) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.clusterID = value
	b.fieldSet_[6] = true
	return b
}

// Groups sets the value of the 'groups' attribute to the given values.
func (b *LogForwarderBuilder) Groups(values ...*LogForwarderGroupBuilder) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.groups = make([]*LogForwarderGroupBuilder, len(values))
	copy(b.groups, values)
	b.fieldSet_[7] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Represents the status of a log forwarder.
func (b *LogForwarderBuilder) Status(value *LogForwarderStatusBuilder) *LogForwarderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogForwarderBuilder) Copy(object *LogForwarder) *LogForwarderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.s3 != nil {
		b.s3 = NewLogForwarderS3Config().Copy(object.s3)
	} else {
		b.s3 = nil
	}
	if object.applications != nil {
		b.applications = make([]string, len(object.applications))
		copy(b.applications, object.applications)
	} else {
		b.applications = nil
	}
	if object.cloudwatch != nil {
		b.cloudwatch = NewLogForwarderCloudWatchConfig().Copy(object.cloudwatch)
	} else {
		b.cloudwatch = nil
	}
	b.clusterID = object.clusterID
	if object.groups != nil {
		b.groups = make([]*LogForwarderGroupBuilder, len(object.groups))
		for i, v := range object.groups {
			b.groups[i] = NewLogForwarderGroup().Copy(v)
		}
	} else {
		b.groups = nil
	}
	if object.status != nil {
		b.status = NewLogForwarderStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	return b
}

// Build creates a 'log_forwarder' object using the configuration stored in the builder.
func (b *LogForwarderBuilder) Build() (object *LogForwarder, err error) {
	object = new(LogForwarder)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.s3 != nil {
		object.s3, err = b.s3.Build()
		if err != nil {
			return
		}
	}
	if b.applications != nil {
		object.applications = make([]string, len(b.applications))
		copy(object.applications, b.applications)
	}
	if b.cloudwatch != nil {
		object.cloudwatch, err = b.cloudwatch.Build()
		if err != nil {
			return
		}
	}
	object.clusterID = b.clusterID
	if b.groups != nil {
		object.groups = make([]*LogForwarderGroup, len(b.groups))
		for i, v := range b.groups {
			object.groups[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	return
}
