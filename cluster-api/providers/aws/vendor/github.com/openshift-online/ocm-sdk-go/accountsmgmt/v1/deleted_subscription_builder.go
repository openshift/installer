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

import (
	time "time"
)

// DeletedSubscriptionBuilder contains the data and logic needed to build 'deleted_subscription' objects.
type DeletedSubscriptionBuilder struct {
	bitmap_                   uint64
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
	return &DeletedSubscriptionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *DeletedSubscriptionBuilder) Link(value bool) *DeletedSubscriptionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *DeletedSubscriptionBuilder) ID(value string) *DeletedSubscriptionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *DeletedSubscriptionBuilder) HREF(value string) *DeletedSubscriptionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DeletedSubscriptionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BillingExpirationDate sets the value of the 'billing_expiration_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) BillingExpirationDate(value time.Time) *DeletedSubscriptionBuilder {
	b.billingExpirationDate = value
	b.bitmap_ |= 8
	return b
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *DeletedSubscriptionBuilder) BillingMarketplaceAccount(value string) *DeletedSubscriptionBuilder {
	b.billingMarketplaceAccount = value
	b.bitmap_ |= 16
	return b
}

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CloudAccountID(value string) *DeletedSubscriptionBuilder {
	b.cloudAccountID = value
	b.bitmap_ |= 32
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CloudProviderID(value string) *DeletedSubscriptionBuilder {
	b.cloudProviderID = value
	b.bitmap_ |= 64
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ClusterID(value string) *DeletedSubscriptionBuilder {
	b.clusterID = value
	b.bitmap_ |= 128
	return b
}

// ClusterBillingModel sets the value of the 'cluster_billing_model' attribute to the given value.
//
// Billing model for subscripiton and reserved_resource resources.
func (b *DeletedSubscriptionBuilder) ClusterBillingModel(value BillingModel) *DeletedSubscriptionBuilder {
	b.clusterBillingModel = value
	b.bitmap_ |= 256
	return b
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ConsoleURL(value string) *DeletedSubscriptionBuilder {
	b.consoleURL = value
	b.bitmap_ |= 512
	return b
}

// ConsumerUUID sets the value of the 'consumer_UUID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ConsumerUUID(value string) *DeletedSubscriptionBuilder {
	b.consumerUUID = value
	b.bitmap_ |= 1024
	return b
}

// CpuTotal sets the value of the 'cpu_total' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CpuTotal(value int) *DeletedSubscriptionBuilder {
	b.cpuTotal = value
	b.bitmap_ |= 2048
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CreatedAt(value time.Time) *DeletedSubscriptionBuilder {
	b.createdAt = value
	b.bitmap_ |= 4096
	return b
}

// CreatorId sets the value of the 'creator_id' attribute to the given value.
func (b *DeletedSubscriptionBuilder) CreatorId(value string) *DeletedSubscriptionBuilder {
	b.creatorId = value
	b.bitmap_ |= 8192
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *DeletedSubscriptionBuilder) DisplayName(value string) *DeletedSubscriptionBuilder {
	b.displayName = value
	b.bitmap_ |= 16384
	return b
}

// ExternalClusterID sets the value of the 'external_cluster_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ExternalClusterID(value string) *DeletedSubscriptionBuilder {
	b.externalClusterID = value
	b.bitmap_ |= 32768
	return b
}

// LastReconcileDate sets the value of the 'last_reconcile_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastReconcileDate(value time.Time) *DeletedSubscriptionBuilder {
	b.lastReconcileDate = value
	b.bitmap_ |= 65536
	return b
}

// LastReleasedAt sets the value of the 'last_released_at' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastReleasedAt(value time.Time) *DeletedSubscriptionBuilder {
	b.lastReleasedAt = value
	b.bitmap_ |= 131072
	return b
}

// LastTelemetryDate sets the value of the 'last_telemetry_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) LastTelemetryDate(value time.Time) *DeletedSubscriptionBuilder {
	b.lastTelemetryDate = value
	b.bitmap_ |= 262144
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Managed(value bool) *DeletedSubscriptionBuilder {
	b.managed = value
	b.bitmap_ |= 524288
	return b
}

// Metrics sets the value of the 'metrics' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Metrics(value string) *DeletedSubscriptionBuilder {
	b.metrics = value
	b.bitmap_ |= 1048576
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) OrganizationID(value string) *DeletedSubscriptionBuilder {
	b.organizationID = value
	b.bitmap_ |= 2097152
	return b
}

// PlanID sets the value of the 'plan_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) PlanID(value string) *DeletedSubscriptionBuilder {
	b.planID = value
	b.bitmap_ |= 4194304
	return b
}

// ProductBundle sets the value of the 'product_bundle' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ProductBundle(value string) *DeletedSubscriptionBuilder {
	b.productBundle = value
	b.bitmap_ |= 8388608
	return b
}

// Provenance sets the value of the 'provenance' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Provenance(value string) *DeletedSubscriptionBuilder {
	b.provenance = value
	b.bitmap_ |= 16777216
	return b
}

// QueryTimestamp sets the value of the 'query_timestamp' attribute to the given value.
func (b *DeletedSubscriptionBuilder) QueryTimestamp(value time.Time) *DeletedSubscriptionBuilder {
	b.queryTimestamp = value
	b.bitmap_ |= 33554432
	return b
}

// RegionID sets the value of the 'region_ID' attribute to the given value.
func (b *DeletedSubscriptionBuilder) RegionID(value string) *DeletedSubscriptionBuilder {
	b.regionID = value
	b.bitmap_ |= 67108864
	return b
}

// Released sets the value of the 'released' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Released(value bool) *DeletedSubscriptionBuilder {
	b.released = value
	b.bitmap_ |= 134217728
	return b
}

// ServiceLevel sets the value of the 'service_level' attribute to the given value.
func (b *DeletedSubscriptionBuilder) ServiceLevel(value string) *DeletedSubscriptionBuilder {
	b.serviceLevel = value
	b.bitmap_ |= 268435456
	return b
}

// SocketTotal sets the value of the 'socket_total' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SocketTotal(value int) *DeletedSubscriptionBuilder {
	b.socketTotal = value
	b.bitmap_ |= 536870912
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Status(value string) *DeletedSubscriptionBuilder {
	b.status = value
	b.bitmap_ |= 1073741824
	return b
}

// SupportLevel sets the value of the 'support_level' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SupportLevel(value string) *DeletedSubscriptionBuilder {
	b.supportLevel = value
	b.bitmap_ |= 2147483648
	return b
}

// SystemUnits sets the value of the 'system_units' attribute to the given value.
func (b *DeletedSubscriptionBuilder) SystemUnits(value string) *DeletedSubscriptionBuilder {
	b.systemUnits = value
	b.bitmap_ |= 4294967296
	return b
}

// TrialEndDate sets the value of the 'trial_end_date' attribute to the given value.
func (b *DeletedSubscriptionBuilder) TrialEndDate(value time.Time) *DeletedSubscriptionBuilder {
	b.trialEndDate = value
	b.bitmap_ |= 8589934592
	return b
}

// Usage sets the value of the 'usage' attribute to the given value.
func (b *DeletedSubscriptionBuilder) Usage(value string) *DeletedSubscriptionBuilder {
	b.usage = value
	b.bitmap_ |= 17179869184
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DeletedSubscriptionBuilder) Copy(object *DeletedSubscription) *DeletedSubscriptionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
