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

// Detailed status of a cluster.
type ClusterStatusBuilder struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	configurationMode         ClusterConfigurationMode
	currentCompute            int
	description               string
	limitedSupportReasonCount int
	provisionErrorCode        string
	provisionErrorMessage     string
	state                     ClusterState
	dnsReady                  bool
	oidcReady                 bool
}

// NewClusterStatus creates a new builder of 'cluster_status' objects.
func NewClusterStatus() *ClusterStatusBuilder {
	return &ClusterStatusBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterStatusBuilder) Link(value bool) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ClusterStatusBuilder) ID(value string) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ClusterStatusBuilder) HREF(value string) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterStatusBuilder) Empty() bool {
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

// DNSReady sets the value of the 'DNS_ready' attribute to the given value.
func (b *ClusterStatusBuilder) DNSReady(value bool) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.dnsReady = value
	b.fieldSet_[3] = true
	return b
}

// OIDCReady sets the value of the 'OIDC_ready' attribute to the given value.
func (b *ClusterStatusBuilder) OIDCReady(value bool) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.oidcReady = value
	b.fieldSet_[4] = true
	return b
}

// ConfigurationMode sets the value of the 'configuration_mode' attribute to the given value.
//
// Configuration mode of a cluster.
func (b *ClusterStatusBuilder) ConfigurationMode(value ClusterConfigurationMode) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.configurationMode = value
	b.fieldSet_[5] = true
	return b
}

// CurrentCompute sets the value of the 'current_compute' attribute to the given value.
func (b *ClusterStatusBuilder) CurrentCompute(value int) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.currentCompute = value
	b.fieldSet_[6] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *ClusterStatusBuilder) Description(value string) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.description = value
	b.fieldSet_[7] = true
	return b
}

// LimitedSupportReasonCount sets the value of the 'limited_support_reason_count' attribute to the given value.
func (b *ClusterStatusBuilder) LimitedSupportReasonCount(value int) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.limitedSupportReasonCount = value
	b.fieldSet_[8] = true
	return b
}

// ProvisionErrorCode sets the value of the 'provision_error_code' attribute to the given value.
func (b *ClusterStatusBuilder) ProvisionErrorCode(value string) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.provisionErrorCode = value
	b.fieldSet_[9] = true
	return b
}

// ProvisionErrorMessage sets the value of the 'provision_error_message' attribute to the given value.
func (b *ClusterStatusBuilder) ProvisionErrorMessage(value string) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.provisionErrorMessage = value
	b.fieldSet_[10] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Overall state of a cluster.
func (b *ClusterStatusBuilder) State(value ClusterState) *ClusterStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.state = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterStatusBuilder) Copy(object *ClusterStatus) *ClusterStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.dnsReady = object.dnsReady
	b.oidcReady = object.oidcReady
	b.configurationMode = object.configurationMode
	b.currentCompute = object.currentCompute
	b.description = object.description
	b.limitedSupportReasonCount = object.limitedSupportReasonCount
	b.provisionErrorCode = object.provisionErrorCode
	b.provisionErrorMessage = object.provisionErrorMessage
	b.state = object.state
	return b
}

// Build creates a 'cluster_status' object using the configuration stored in the builder.
func (b *ClusterStatusBuilder) Build() (object *ClusterStatus, err error) {
	object = new(ClusterStatus)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.dnsReady = b.dnsReady
	object.oidcReady = b.oidcReady
	object.configurationMode = b.configurationMode
	object.currentCompute = b.currentCompute
	object.description = b.description
	object.limitedSupportReasonCount = b.limitedSupportReasonCount
	object.provisionErrorCode = b.provisionErrorCode
	object.provisionErrorMessage = b.provisionErrorMessage
	object.state = b.state
	return
}
