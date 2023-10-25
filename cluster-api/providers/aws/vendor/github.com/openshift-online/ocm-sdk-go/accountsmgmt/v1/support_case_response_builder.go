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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// SupportCaseResponseBuilder contains the data and logic needed to build 'support_case_response' objects.
type SupportCaseResponseBuilder struct {
	bitmap_        uint32
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
	return &SupportCaseResponseBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SupportCaseResponseBuilder) Link(value bool) *SupportCaseResponseBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SupportCaseResponseBuilder) ID(value string) *SupportCaseResponseBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SupportCaseResponseBuilder) HREF(value string) *SupportCaseResponseBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SupportCaseResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// URI sets the value of the 'URI' attribute to the given value.
func (b *SupportCaseResponseBuilder) URI(value string) *SupportCaseResponseBuilder {
	b.uri = value
	b.bitmap_ |= 8
	return b
}

// CaseNumber sets the value of the 'case_number' attribute to the given value.
func (b *SupportCaseResponseBuilder) CaseNumber(value string) *SupportCaseResponseBuilder {
	b.caseNumber = value
	b.bitmap_ |= 16
	return b
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *SupportCaseResponseBuilder) ClusterId(value string) *SupportCaseResponseBuilder {
	b.clusterId = value
	b.bitmap_ |= 32
	return b
}

// ClusterUuid sets the value of the 'cluster_uuid' attribute to the given value.
func (b *SupportCaseResponseBuilder) ClusterUuid(value string) *SupportCaseResponseBuilder {
	b.clusterUuid = value
	b.bitmap_ |= 64
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *SupportCaseResponseBuilder) Description(value string) *SupportCaseResponseBuilder {
	b.description = value
	b.bitmap_ |= 128
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
func (b *SupportCaseResponseBuilder) Severity(value string) *SupportCaseResponseBuilder {
	b.severity = value
	b.bitmap_ |= 256
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *SupportCaseResponseBuilder) Status(value string) *SupportCaseResponseBuilder {
	b.status = value
	b.bitmap_ |= 512
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *SupportCaseResponseBuilder) SubscriptionId(value string) *SupportCaseResponseBuilder {
	b.subscriptionId = value
	b.bitmap_ |= 1024
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *SupportCaseResponseBuilder) Summary(value string) *SupportCaseResponseBuilder {
	b.summary = value
	b.bitmap_ |= 2048
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SupportCaseResponseBuilder) Copy(object *SupportCaseResponse) *SupportCaseResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
