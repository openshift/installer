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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

type SupportCaseRequestBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	clusterId      string
	clusterUuid    string
	description    string
	eventStreamId  string
	severity       string
	subscriptionId string
	summary        string
}

// NewSupportCaseRequest creates a new builder of 'support_case_request' objects.
func NewSupportCaseRequest() *SupportCaseRequestBuilder {
	return &SupportCaseRequestBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *SupportCaseRequestBuilder) Link(value bool) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *SupportCaseRequestBuilder) ID(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *SupportCaseRequestBuilder) HREF(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SupportCaseRequestBuilder) Empty() bool {
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

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *SupportCaseRequestBuilder) ClusterId(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.clusterId = value
	b.fieldSet_[3] = true
	return b
}

// ClusterUuid sets the value of the 'cluster_uuid' attribute to the given value.
func (b *SupportCaseRequestBuilder) ClusterUuid(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.clusterUuid = value
	b.fieldSet_[4] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *SupportCaseRequestBuilder) Description(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.description = value
	b.fieldSet_[5] = true
	return b
}

// EventStreamId sets the value of the 'event_stream_id' attribute to the given value.
func (b *SupportCaseRequestBuilder) EventStreamId(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.eventStreamId = value
	b.fieldSet_[6] = true
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
func (b *SupportCaseRequestBuilder) Severity(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.severity = value
	b.fieldSet_[7] = true
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *SupportCaseRequestBuilder) SubscriptionId(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.subscriptionId = value
	b.fieldSet_[8] = true
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *SupportCaseRequestBuilder) Summary(value string) *SupportCaseRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.summary = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SupportCaseRequestBuilder) Copy(object *SupportCaseRequest) *SupportCaseRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.clusterId = object.clusterId
	b.clusterUuid = object.clusterUuid
	b.description = object.description
	b.eventStreamId = object.eventStreamId
	b.severity = object.severity
	b.subscriptionId = object.subscriptionId
	b.summary = object.summary
	return b
}

// Build creates a 'support_case_request' object using the configuration stored in the builder.
func (b *SupportCaseRequestBuilder) Build() (object *SupportCaseRequest, err error) {
	object = new(SupportCaseRequest)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterId = b.clusterId
	object.clusterUuid = b.clusterUuid
	object.description = b.description
	object.eventStreamId = b.eventStreamId
	object.severity = b.severity
	object.subscriptionId = b.subscriptionId
	object.summary = b.summary
	return
}
