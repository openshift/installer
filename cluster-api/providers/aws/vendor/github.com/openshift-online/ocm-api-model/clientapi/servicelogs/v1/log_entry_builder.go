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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicelogs/v1

import (
	time "time"
)

type LogEntryBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	clusterID      string
	clusterUUID    string
	createdAt      time.Time
	createdBy      string
	description    string
	docReferences  []string
	eventStreamID  string
	logType        LogType
	serviceName    string
	severity       Severity
	subscriptionID string
	summary        string
	timestamp      time.Time
	username       string
	internalOnly   bool
}

// NewLogEntry creates a new builder of 'log_entry' objects.
func NewLogEntry() *LogEntryBuilder {
	return &LogEntryBuilder{
		fieldSet_: make([]bool, 18),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *LogEntryBuilder) Link(value bool) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *LogEntryBuilder) ID(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *LogEntryBuilder) HREF(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogEntryBuilder) Empty() bool {
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

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *LogEntryBuilder) ClusterID(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *LogEntryBuilder) ClusterUUID(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.clusterUUID = value
	b.fieldSet_[4] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *LogEntryBuilder) CreatedAt(value time.Time) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.createdAt = value
	b.fieldSet_[5] = true
	return b
}

// CreatedBy sets the value of the 'created_by' attribute to the given value.
func (b *LogEntryBuilder) CreatedBy(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.createdBy = value
	b.fieldSet_[6] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *LogEntryBuilder) Description(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.description = value
	b.fieldSet_[7] = true
	return b
}

// DocReferences sets the value of the 'doc_references' attribute to the given values.
func (b *LogEntryBuilder) DocReferences(values ...string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.docReferences = make([]string, len(values))
	copy(b.docReferences, values)
	b.fieldSet_[8] = true
	return b
}

// EventStreamID sets the value of the 'event_stream_ID' attribute to the given value.
func (b *LogEntryBuilder) EventStreamID(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.eventStreamID = value
	b.fieldSet_[9] = true
	return b
}

// InternalOnly sets the value of the 'internal_only' attribute to the given value.
func (b *LogEntryBuilder) InternalOnly(value bool) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.internalOnly = value
	b.fieldSet_[10] = true
	return b
}

// LogType sets the value of the 'log_type' attribute to the given value.
//
// Representation of the log type field used in cluster log type model.
func (b *LogEntryBuilder) LogType(value LogType) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.logType = value
	b.fieldSet_[11] = true
	return b
}

// ServiceName sets the value of the 'service_name' attribute to the given value.
func (b *LogEntryBuilder) ServiceName(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.serviceName = value
	b.fieldSet_[12] = true
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
func (b *LogEntryBuilder) Severity(value Severity) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.severity = value
	b.fieldSet_[13] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *LogEntryBuilder) SubscriptionID(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.subscriptionID = value
	b.fieldSet_[14] = true
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *LogEntryBuilder) Summary(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.summary = value
	b.fieldSet_[15] = true
	return b
}

// Timestamp sets the value of the 'timestamp' attribute to the given value.
func (b *LogEntryBuilder) Timestamp(value time.Time) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.timestamp = value
	b.fieldSet_[16] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *LogEntryBuilder) Username(value string) *LogEntryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.username = value
	b.fieldSet_[17] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogEntryBuilder) Copy(object *LogEntry) *LogEntryBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.createdAt = object.createdAt
	b.createdBy = object.createdBy
	b.description = object.description
	if object.docReferences != nil {
		b.docReferences = make([]string, len(object.docReferences))
		copy(b.docReferences, object.docReferences)
	} else {
		b.docReferences = nil
	}
	b.eventStreamID = object.eventStreamID
	b.internalOnly = object.internalOnly
	b.logType = object.logType
	b.serviceName = object.serviceName
	b.severity = object.severity
	b.subscriptionID = object.subscriptionID
	b.summary = object.summary
	b.timestamp = object.timestamp
	b.username = object.username
	return b
}

// Build creates a 'log_entry' object using the configuration stored in the builder.
func (b *LogEntryBuilder) Build() (object *LogEntry, err error) {
	object = new(LogEntry)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.createdAt = b.createdAt
	object.createdBy = b.createdBy
	object.description = b.description
	if b.docReferences != nil {
		object.docReferences = make([]string, len(b.docReferences))
		copy(object.docReferences, b.docReferences)
	}
	object.eventStreamID = b.eventStreamID
	object.internalOnly = b.internalOnly
	object.logType = b.logType
	object.serviceName = b.serviceName
	object.severity = b.severity
	object.subscriptionID = b.subscriptionID
	object.summary = b.summary
	object.timestamp = b.timestamp
	object.username = b.username
	return
}
