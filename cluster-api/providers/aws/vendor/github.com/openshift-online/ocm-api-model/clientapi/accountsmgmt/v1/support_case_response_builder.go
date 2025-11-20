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

type SupportCaseResponseBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	uri            string
	caseNumber     string
	clusterId      string
	clusterUuid    string
	description    string
	severity       string
	status         string
	subscriptionId string
	summary        string
}

// NewSupportCaseResponse creates a new builder of 'support_case_response' objects.
func NewSupportCaseResponse() *SupportCaseResponseBuilder {
	return &SupportCaseResponseBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *SupportCaseResponseBuilder) Link(value bool) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *SupportCaseResponseBuilder) ID(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *SupportCaseResponseBuilder) HREF(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SupportCaseResponseBuilder) Empty() bool {
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

// URI sets the value of the 'URI' attribute to the given value.
func (b *SupportCaseResponseBuilder) URI(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.uri = value
	b.fieldSet_[3] = true
	return b
}

// CaseNumber sets the value of the 'case_number' attribute to the given value.
func (b *SupportCaseResponseBuilder) CaseNumber(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.caseNumber = value
	b.fieldSet_[4] = true
	return b
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *SupportCaseResponseBuilder) ClusterId(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.clusterId = value
	b.fieldSet_[5] = true
	return b
}

// ClusterUuid sets the value of the 'cluster_uuid' attribute to the given value.
func (b *SupportCaseResponseBuilder) ClusterUuid(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.clusterUuid = value
	b.fieldSet_[6] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *SupportCaseResponseBuilder) Description(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.description = value
	b.fieldSet_[7] = true
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
func (b *SupportCaseResponseBuilder) Severity(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.severity = value
	b.fieldSet_[8] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *SupportCaseResponseBuilder) Status(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.status = value
	b.fieldSet_[9] = true
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *SupportCaseResponseBuilder) SubscriptionId(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.subscriptionId = value
	b.fieldSet_[10] = true
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *SupportCaseResponseBuilder) Summary(value string) *SupportCaseResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.summary = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SupportCaseResponseBuilder) Copy(object *SupportCaseResponse) *SupportCaseResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.uri = object.uri
	b.caseNumber = object.caseNumber
	b.clusterId = object.clusterId
	b.clusterUuid = object.clusterUuid
	b.description = object.description
	b.severity = object.severity
	b.status = object.status
	b.subscriptionId = object.subscriptionId
	b.summary = object.summary
	return b
}

// Build creates a 'support_case_response' object using the configuration stored in the builder.
func (b *SupportCaseResponseBuilder) Build() (object *SupportCaseResponse, err error) {
	object = new(SupportCaseResponse)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.uri = b.uri
	object.caseNumber = b.caseNumber
	object.clusterId = b.clusterId
	object.clusterUuid = b.clusterUuid
	object.description = b.description
	object.severity = b.severity
	object.status = b.status
	object.subscriptionId = b.subscriptionId
	object.summary = b.summary
	return
}
