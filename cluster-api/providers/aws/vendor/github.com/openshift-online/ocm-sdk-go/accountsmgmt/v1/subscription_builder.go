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

// SubscriptionBuilder contains the data and logic needed to build 'subscription' objects.
type SubscriptionBuilder struct {
	bitmap_                   uint64
	id                        string
	href                      string
	billingMarketplaceAccount string
	capabilities              []*CapabilityBuilder
	cloudAccountID            string
	cloudProviderID           string
	clusterID                 string
	clusterBillingModel       BillingModel
	consoleURL                string
	consumerUUID              string
	cpuTotal                  int
	createdAt                 time.Time
	creator                   *AccountBuilder
	displayName               string
	externalClusterID         string
	labels                    []*LabelBuilder
	lastReconcileDate         time.Time
	lastReleasedAt            time.Time
	lastTelemetryDate         time.Time
	metrics                   []*SubscriptionMetricsBuilder
	notificationContacts      []*AccountBuilder
	organizationID            string
	plan                      *PlanBuilder
	productBundle             string
	provenance                string
	regionID                  string
	serviceLevel              string
	socketTotal               int
	status                    string
	supportLevel              string
	systemUnits               string
	trialEndDate              time.Time
	updatedAt                 time.Time
	usage                     string
	managed                   bool
	released                  bool
}

// NewSubscription creates a new builder of 'subscription' objects.
func NewSubscription() *SubscriptionBuilder {
	return &SubscriptionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SubscriptionBuilder) Link(value bool) *SubscriptionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SubscriptionBuilder) ID(value string) *SubscriptionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SubscriptionBuilder) HREF(value string) *SubscriptionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubscriptionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *SubscriptionBuilder) BillingMarketplaceAccount(value string) *SubscriptionBuilder {
	b.billingMarketplaceAccount = value
	b.bitmap_ |= 8
	return b
}

// Capabilities sets the value of the 'capabilities' attribute to the given values.
func (b *SubscriptionBuilder) Capabilities(values ...*CapabilityBuilder) *SubscriptionBuilder {
	b.capabilities = make([]*CapabilityBuilder, len(values))
	copy(b.capabilities, values)
	b.bitmap_ |= 16
	return b
}

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *SubscriptionBuilder) CloudAccountID(value string) *SubscriptionBuilder {
	b.cloudAccountID = value
	b.bitmap_ |= 32
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *SubscriptionBuilder) CloudProviderID(value string) *SubscriptionBuilder {
	b.cloudProviderID = value
	b.bitmap_ |= 64
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *SubscriptionBuilder) ClusterID(value string) *SubscriptionBuilder {
	b.clusterID = value
	b.bitmap_ |= 128
	return b
}

// ClusterBillingModel sets the value of the 'cluster_billing_model' attribute to the given value.
//
// Billing model for subscripiton and reserved_resource resources.
func (b *SubscriptionBuilder) ClusterBillingModel(value BillingModel) *SubscriptionBuilder {
	b.clusterBillingModel = value
	b.bitmap_ |= 256
	return b
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *SubscriptionBuilder) ConsoleURL(value string) *SubscriptionBuilder {
	b.consoleURL = value
	b.bitmap_ |= 512
	return b
}

// ConsumerUUID sets the value of the 'consumer_UUID' attribute to the given value.
func (b *SubscriptionBuilder) ConsumerUUID(value string) *SubscriptionBuilder {
	b.consumerUUID = value
	b.bitmap_ |= 1024
	return b
}

// CpuTotal sets the value of the 'cpu_total' attribute to the given value.
func (b *SubscriptionBuilder) CpuTotal(value int) *SubscriptionBuilder {
	b.cpuTotal = value
	b.bitmap_ |= 2048
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *SubscriptionBuilder) CreatedAt(value time.Time) *SubscriptionBuilder {
	b.createdAt = value
	b.bitmap_ |= 4096
	return b
}

// Creator sets the value of the 'creator' attribute to the given value.
func (b *SubscriptionBuilder) Creator(value *AccountBuilder) *SubscriptionBuilder {
	b.creator = value
	if value != nil {
		b.bitmap_ |= 8192
	} else {
		b.bitmap_ &^= 8192
	}
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *SubscriptionBuilder) DisplayName(value string) *SubscriptionBuilder {
	b.displayName = value
	b.bitmap_ |= 16384
	return b
}

// ExternalClusterID sets the value of the 'external_cluster_ID' attribute to the given value.
func (b *SubscriptionBuilder) ExternalClusterID(value string) *SubscriptionBuilder {
	b.externalClusterID = value
	b.bitmap_ |= 32768
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *SubscriptionBuilder) Labels(values ...*LabelBuilder) *SubscriptionBuilder {
	b.labels = make([]*LabelBuilder, len(values))
	copy(b.labels, values)
	b.bitmap_ |= 65536
	return b
}

// LastReconcileDate sets the value of the 'last_reconcile_date' attribute to the given value.
func (b *SubscriptionBuilder) LastReconcileDate(value time.Time) *SubscriptionBuilder {
	b.lastReconcileDate = value
	b.bitmap_ |= 131072
	return b
}

// LastReleasedAt sets the value of the 'last_released_at' attribute to the given value.
func (b *SubscriptionBuilder) LastReleasedAt(value time.Time) *SubscriptionBuilder {
	b.lastReleasedAt = value
	b.bitmap_ |= 262144
	return b
}

// LastTelemetryDate sets the value of the 'last_telemetry_date' attribute to the given value.
func (b *SubscriptionBuilder) LastTelemetryDate(value time.Time) *SubscriptionBuilder {
	b.lastTelemetryDate = value
	b.bitmap_ |= 524288
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *SubscriptionBuilder) Managed(value bool) *SubscriptionBuilder {
	b.managed = value
	b.bitmap_ |= 1048576
	return b
}

// Metrics sets the value of the 'metrics' attribute to the given values.
func (b *SubscriptionBuilder) Metrics(values ...*SubscriptionMetricsBuilder) *SubscriptionBuilder {
	b.metrics = make([]*SubscriptionMetricsBuilder, len(values))
	copy(b.metrics, values)
	b.bitmap_ |= 2097152
	return b
}

// NotificationContacts sets the value of the 'notification_contacts' attribute to the given values.
func (b *SubscriptionBuilder) NotificationContacts(values ...*AccountBuilder) *SubscriptionBuilder {
	b.notificationContacts = make([]*AccountBuilder, len(values))
	copy(b.notificationContacts, values)
	b.bitmap_ |= 4194304
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *SubscriptionBuilder) OrganizationID(value string) *SubscriptionBuilder {
	b.organizationID = value
	b.bitmap_ |= 8388608
	return b
}

// Plan sets the value of the 'plan' attribute to the given value.
func (b *SubscriptionBuilder) Plan(value *PlanBuilder) *SubscriptionBuilder {
	b.plan = value
	if value != nil {
		b.bitmap_ |= 16777216
	} else {
		b.bitmap_ &^= 16777216
	}
	return b
}

// ProductBundle sets the value of the 'product_bundle' attribute to the given value.
func (b *SubscriptionBuilder) ProductBundle(value string) *SubscriptionBuilder {
	b.productBundle = value
	b.bitmap_ |= 33554432
	return b
}

// Provenance sets the value of the 'provenance' attribute to the given value.
func (b *SubscriptionBuilder) Provenance(value string) *SubscriptionBuilder {
	b.provenance = value
	b.bitmap_ |= 67108864
	return b
}

// RegionID sets the value of the 'region_ID' attribute to the given value.
func (b *SubscriptionBuilder) RegionID(value string) *SubscriptionBuilder {
	b.regionID = value
	b.bitmap_ |= 134217728
	return b
}

// Released sets the value of the 'released' attribute to the given value.
func (b *SubscriptionBuilder) Released(value bool) *SubscriptionBuilder {
	b.released = value
	b.bitmap_ |= 268435456
	return b
}

// ServiceLevel sets the value of the 'service_level' attribute to the given value.
func (b *SubscriptionBuilder) ServiceLevel(value string) *SubscriptionBuilder {
	b.serviceLevel = value
	b.bitmap_ |= 536870912
	return b
}

// SocketTotal sets the value of the 'socket_total' attribute to the given value.
func (b *SubscriptionBuilder) SocketTotal(value int) *SubscriptionBuilder {
	b.socketTotal = value
	b.bitmap_ |= 1073741824
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *SubscriptionBuilder) Status(value string) *SubscriptionBuilder {
	b.status = value
	b.bitmap_ |= 2147483648
	return b
}

// SupportLevel sets the value of the 'support_level' attribute to the given value.
func (b *SubscriptionBuilder) SupportLevel(value string) *SubscriptionBuilder {
	b.supportLevel = value
	b.bitmap_ |= 4294967296
	return b
}

// SystemUnits sets the value of the 'system_units' attribute to the given value.
func (b *SubscriptionBuilder) SystemUnits(value string) *SubscriptionBuilder {
	b.systemUnits = value
	b.bitmap_ |= 8589934592
	return b
}

// TrialEndDate sets the value of the 'trial_end_date' attribute to the given value.
func (b *SubscriptionBuilder) TrialEndDate(value time.Time) *SubscriptionBuilder {
	b.trialEndDate = value
	b.bitmap_ |= 17179869184
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *SubscriptionBuilder) UpdatedAt(value time.Time) *SubscriptionBuilder {
	b.updatedAt = value
	b.bitmap_ |= 34359738368
	return b
}

// Usage sets the value of the 'usage' attribute to the given value.
func (b *SubscriptionBuilder) Usage(value string) *SubscriptionBuilder {
	b.usage = value
	b.bitmap_ |= 68719476736
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubscriptionBuilder) Copy(object *Subscription) *SubscriptionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	if object.capabilities != nil {
		b.capabilities = make([]*CapabilityBuilder, len(object.capabilities))
		for i, v := range object.capabilities {
			b.capabilities[i] = NewCapability().Copy(v)
		}
	} else {
		b.capabilities = nil
	}
	b.cloudAccountID = object.cloudAccountID
	b.cloudProviderID = object.cloudProviderID
	b.clusterID = object.clusterID
	b.clusterBillingModel = object.clusterBillingModel
	b.consoleURL = object.consoleURL
	b.consumerUUID = object.consumerUUID
	b.cpuTotal = object.cpuTotal
	b.createdAt = object.createdAt
	if object.creator != nil {
		b.creator = NewAccount().Copy(object.creator)
	} else {
		b.creator = nil
	}
	b.displayName = object.displayName
	b.externalClusterID = object.externalClusterID
	if object.labels != nil {
		b.labels = make([]*LabelBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabel().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.lastReconcileDate = object.lastReconcileDate
	b.lastReleasedAt = object.lastReleasedAt
	b.lastTelemetryDate = object.lastTelemetryDate
	b.managed = object.managed
	if object.metrics != nil {
		b.metrics = make([]*SubscriptionMetricsBuilder, len(object.metrics))
		for i, v := range object.metrics {
			b.metrics[i] = NewSubscriptionMetrics().Copy(v)
		}
	} else {
		b.metrics = nil
	}
	if object.notificationContacts != nil {
		b.notificationContacts = make([]*AccountBuilder, len(object.notificationContacts))
		for i, v := range object.notificationContacts {
			b.notificationContacts[i] = NewAccount().Copy(v)
		}
	} else {
		b.notificationContacts = nil
	}
	b.organizationID = object.organizationID
	if object.plan != nil {
		b.plan = NewPlan().Copy(object.plan)
	} else {
		b.plan = nil
	}
	b.productBundle = object.productBundle
	b.provenance = object.provenance
	b.regionID = object.regionID
	b.released = object.released
	b.serviceLevel = object.serviceLevel
	b.socketTotal = object.socketTotal
	b.status = object.status
	b.supportLevel = object.supportLevel
	b.systemUnits = object.systemUnits
	b.trialEndDate = object.trialEndDate
	b.updatedAt = object.updatedAt
	b.usage = object.usage
	return b
}

// Build creates a 'subscription' object using the configuration stored in the builder.
func (b *SubscriptionBuilder) Build() (object *Subscription, err error) {
	object = new(Subscription)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	if b.capabilities != nil {
		object.capabilities = make([]*Capability, len(b.capabilities))
		for i, v := range b.capabilities {
			object.capabilities[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.cloudAccountID = b.cloudAccountID
	object.cloudProviderID = b.cloudProviderID
	object.clusterID = b.clusterID
	object.clusterBillingModel = b.clusterBillingModel
	object.consoleURL = b.consoleURL
	object.consumerUUID = b.consumerUUID
	object.cpuTotal = b.cpuTotal
	object.createdAt = b.createdAt
	if b.creator != nil {
		object.creator, err = b.creator.Build()
		if err != nil {
			return
		}
	}
	object.displayName = b.displayName
	object.externalClusterID = b.externalClusterID
	if b.labels != nil {
		object.labels = make([]*Label, len(b.labels))
		for i, v := range b.labels {
			object.labels[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.lastReconcileDate = b.lastReconcileDate
	object.lastReleasedAt = b.lastReleasedAt
	object.lastTelemetryDate = b.lastTelemetryDate
	object.managed = b.managed
	if b.metrics != nil {
		object.metrics = make([]*SubscriptionMetrics, len(b.metrics))
		for i, v := range b.metrics {
			object.metrics[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.notificationContacts != nil {
		object.notificationContacts = make([]*Account, len(b.notificationContacts))
		for i, v := range b.notificationContacts {
			object.notificationContacts[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.organizationID = b.organizationID
	if b.plan != nil {
		object.plan, err = b.plan.Build()
		if err != nil {
			return
		}
	}
	object.productBundle = b.productBundle
	object.provenance = b.provenance
	object.regionID = b.regionID
	object.released = b.released
	object.serviceLevel = b.serviceLevel
	object.socketTotal = b.socketTotal
	object.status = b.status
	object.supportLevel = b.supportLevel
	object.systemUnits = b.systemUnits
	object.trialEndDate = b.trialEndDate
	object.updatedAt = b.updatedAt
	object.usage = b.usage
	return
}
