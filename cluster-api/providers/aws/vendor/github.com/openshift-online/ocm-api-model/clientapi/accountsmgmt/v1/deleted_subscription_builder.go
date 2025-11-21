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

import (
	time "time"
)

type DeletedSubscriptionBuilder struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	billingExpirationDate     time.Time
	billingMarketplaceAccount string
	cloudAccountID            string
	cloudProviderID           string
	clusterID                 string
	clusterBillingModel       BillingModel
	consoleURL                string
	consumerUUID              string
	cpuTotal                  int
	createdAt                 time.Time
	creatorId                 string
	displayName               string
	externalClusterID         string
	lastReconcileDate         time.Time
	lastReleasedAt            time.Time
	lastTelemetryDate         time.Time
	metrics                   string
	organizationID            string
	planID                    string
	productBundle             string
	provenance                string
	queryTimestamp            time.Time
	regionID                  string
	serviceLevel              string
	socketTotal               int
	status                    string
	supportLevel              string
	systemUnits               string
	trialEndDate              time.Time
	usage                     string
	managed                   bool
	released                  bool
}

// NewDeletedSubscription creates a new builder of 'deleted_subscription' objects.
func NewDeletedSubscription() *DeletedSubscriptionBuilder {
	return &DeletedSubscriptionBuilder{
		fieldSet_: make([]bool, 35),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *DeletedSubscriptionBuilder) Link(value bool) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *DeletedSubscriptionBuilder) ID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *DeletedSubscriptionBuilder) HREF(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DeletedSubscriptionBuilder) Empty() bool {
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

// BillingExpirationDate sets the value of the 'billing_expiration_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) BillingExpirationDate(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.billingExpirationDate = value
	b.fieldSet_[3] = true
	return b
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *DeletedSubscriptionBuilder) BillingMarketplaceAccount(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.billingMarketplaceAccount = value
	b.fieldSet_[4] = true
	return b
}

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CloudAccountID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.cloudAccountID = value
	b.fieldSet_[5] = true
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CloudProviderID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.cloudProviderID = value
	b.fieldSet_[6] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ClusterID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.clusterID = value
	b.fieldSet_[7] = true
	return b
}

// ClusterBillingModel sets the value of the 'cluster_billing_model' attribute to the given value.
//
// Billing model for subscripiton and reserved_resource resources.
func (b *DeletedSubscriptionBuilder) ClusterBillingModel(value BillingModel) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.clusterBillingModel = value
	b.fieldSet_[8] = true
	return b
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ConsoleURL(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.consoleURL = value
	b.fieldSet_[9] = true
	return b
}

// ConsumerUUID sets the value of the 'consumer_UUID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ConsumerUUID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.consumerUUID = value
	b.fieldSet_[10] = true
	return b
}

// CpuTotal sets the value of the 'cpu_total' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CpuTotal(value int) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.cpuTotal = value
	b.fieldSet_[11] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CreatedAt(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.createdAt = value
	b.fieldSet_[12] = true
	return b
}

// CreatorId sets the value of the 'creator_id' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CreatorId(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.creatorId = value
	b.fieldSet_[13] = true
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *DeletedSubscriptionBuilder) DisplayName(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.displayName = value
	b.fieldSet_[14] = true
	return b
}

// ExternalClusterID sets the value of the 'external_cluster_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ExternalClusterID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.externalClusterID = value
	b.fieldSet_[15] = true
	return b
}

// LastReconcileDate sets the value of the 'last_reconcile_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastReconcileDate(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.lastReconcileDate = value
	b.fieldSet_[16] = true
	return b
}

// LastReleasedAt sets the value of the 'last_released_at' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastReleasedAt(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.lastReleasedAt = value
	b.fieldSet_[17] = true
	return b
}

// LastTelemetryDate sets the value of the 'last_telemetry_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastTelemetryDate(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.lastTelemetryDate = value
	b.fieldSet_[18] = true
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Managed(value bool) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.managed = value
	b.fieldSet_[19] = true
	return b
}

// Metrics sets the value of the 'metrics' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Metrics(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.metrics = value
	b.fieldSet_[20] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) OrganizationID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.organizationID = value
	b.fieldSet_[21] = true
	return b
}

// PlanID sets the value of the 'plan_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) PlanID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.planID = value
	b.fieldSet_[22] = true
	return b
}

// ProductBundle sets the value of the 'product_bundle' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ProductBundle(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.productBundle = value
	b.fieldSet_[23] = true
	return b
}

// Provenance sets the value of the 'provenance' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Provenance(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.provenance = value
	b.fieldSet_[24] = true
	return b
}

// QueryTimestamp sets the value of the 'query_timestamp' attribute to the given value.
func (b *DeletedSubscriptionBuilder) QueryTimestamp(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.queryTimestamp = value
	b.fieldSet_[25] = true
	return b
}

// RegionID sets the value of the 'region_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) RegionID(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.regionID = value
	b.fieldSet_[26] = true
	return b
}

// Released sets the value of the 'released' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Released(value bool) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.released = value
	b.fieldSet_[27] = true
	return b
}

// ServiceLevel sets the value of the 'service_level' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ServiceLevel(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.serviceLevel = value
	b.fieldSet_[28] = true
	return b
}

// SocketTotal sets the value of the 'socket_total' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SocketTotal(value int) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.socketTotal = value
	b.fieldSet_[29] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Status(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.status = value
	b.fieldSet_[30] = true
	return b
}

// SupportLevel sets the value of the 'support_level' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SupportLevel(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.supportLevel = value
	b.fieldSet_[31] = true
	return b
}

// SystemUnits sets the value of the 'system_units' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SystemUnits(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.systemUnits = value
	b.fieldSet_[32] = true
	return b
}

// TrialEndDate sets the value of the 'trial_end_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) TrialEndDate(value time.Time) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.trialEndDate = value
	b.fieldSet_[33] = true
	return b
}

// Usage sets the value of the 'usage' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Usage(value string) *DeletedSubscriptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 35)
	}
	b.usage = value
	b.fieldSet_[34] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DeletedSubscriptionBuilder) Copy(object *DeletedSubscription) *DeletedSubscriptionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.billingExpirationDate = object.billingExpirationDate
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	b.cloudAccountID = object.cloudAccountID
	b.cloudProviderID = object.cloudProviderID
	b.clusterID = object.clusterID
	b.clusterBillingModel = object.clusterBillingModel
	b.consoleURL = object.consoleURL
	b.consumerUUID = object.consumerUUID
	b.cpuTotal = object.cpuTotal
	b.createdAt = object.createdAt
	b.creatorId = object.creatorId
	b.displayName = object.displayName
	b.externalClusterID = object.externalClusterID
	b.lastReconcileDate = object.lastReconcileDate
	b.lastReleasedAt = object.lastReleasedAt
	b.lastTelemetryDate = object.lastTelemetryDate
	b.managed = object.managed
	b.metrics = object.metrics
	b.organizationID = object.organizationID
	b.planID = object.planID
	b.productBundle = object.productBundle
	b.provenance = object.provenance
	b.queryTimestamp = object.queryTimestamp
	b.regionID = object.regionID
	b.released = object.released
	b.serviceLevel = object.serviceLevel
	b.socketTotal = object.socketTotal
	b.status = object.status
	b.supportLevel = object.supportLevel
	b.systemUnits = object.systemUnits
	b.trialEndDate = object.trialEndDate
	b.usage = object.usage
	return b
}

// Build creates a 'deleted_subscription' object using the configuration stored in the builder.
func (b *DeletedSubscriptionBuilder) Build() (object *DeletedSubscription, err error) {
	object = new(DeletedSubscription)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.billingExpirationDate = b.billingExpirationDate
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.cloudAccountID = b.cloudAccountID
	object.cloudProviderID = b.cloudProviderID
	object.clusterID = b.clusterID
	object.clusterBillingModel = b.clusterBillingModel
	object.consoleURL = b.consoleURL
	object.consumerUUID = b.consumerUUID
	object.cpuTotal = b.cpuTotal
	object.createdAt = b.createdAt
	object.creatorId = b.creatorId
	object.displayName = b.displayName
	object.externalClusterID = b.externalClusterID
	object.lastReconcileDate = b.lastReconcileDate
	object.lastReleasedAt = b.lastReleasedAt
	object.lastTelemetryDate = b.lastTelemetryDate
	object.managed = b.managed
	object.metrics = b.metrics
	object.organizationID = b.organizationID
	object.planID = b.planID
	object.productBundle = b.productBundle
	object.provenance = b.provenance
	object.queryTimestamp = b.queryTimestamp
	object.regionID = b.regionID
	object.released = b.released
	object.serviceLevel = b.serviceLevel
	object.socketTotal = b.socketTotal
	object.status = b.status
	object.supportLevel = b.supportLevel
	object.systemUnits = b.systemUnits
	object.trialEndDate = b.trialEndDate
	object.usage = b.usage
	return
}
