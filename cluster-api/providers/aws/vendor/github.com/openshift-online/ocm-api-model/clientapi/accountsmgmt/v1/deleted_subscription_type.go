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

// DeletedSubscriptionKind is the name of the type used to represent objects
// of type 'deleted_subscription'.
const DeletedSubscriptionKind = "DeletedSubscription"

// DeletedSubscriptionLinkKind is the name of the type used to represent links
// to objects of type 'deleted_subscription'.
const DeletedSubscriptionLinkKind = "DeletedSubscriptionLink"

// DeletedSubscriptionNilKind is the name of the type used to nil references
// to objects of type 'deleted_subscription'.
const DeletedSubscriptionNilKind = "DeletedSubscriptionNil"

// DeletedSubscription represents the values of the 'deleted_subscription' type.
type DeletedSubscription struct {
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

// Kind returns the name of the type of the object.
func (o *DeletedSubscription) Kind() string {
	if o == nil {
		return DeletedSubscriptionNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return DeletedSubscriptionLinkKind
	}
	return DeletedSubscriptionKind
}

// Link returns true if this is a link.
func (o *DeletedSubscription) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *DeletedSubscription) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *DeletedSubscription) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *DeletedSubscription) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *DeletedSubscription) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *DeletedSubscription) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// BillingExpirationDate returns the value of the 'billing_expiration_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) BillingExpirationDate() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.billingExpirationDate
	}
	return time.Time{}
}

// GetBillingExpirationDate returns the value of the 'billing_expiration_date' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetBillingExpirationDate() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.billingExpirationDate
	}
	return
}

// BillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) BillingMarketplaceAccount() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.billingMarketplaceAccount
	}
	return ""
}

// GetBillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetBillingMarketplaceAccount() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.billingMarketplaceAccount
	}
	return
}

// CloudAccountID returns the value of the 'cloud_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) CloudAccountID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.cloudAccountID
	}
	return ""
}

// GetCloudAccountID returns the value of the 'cloud_account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetCloudAccountID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.cloudAccountID
	}
	return
}

// CloudProviderID returns the value of the 'cloud_provider_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) CloudProviderID() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.cloudProviderID
	}
	return ""
}

// GetCloudProviderID returns the value of the 'cloud_provider_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetCloudProviderID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.cloudProviderID
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.clusterID
	}
	return
}

// ClusterBillingModel returns the value of the 'cluster_billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ClusterBillingModel() BillingModel {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.clusterBillingModel
	}
	return BillingModel("")
}

// GetClusterBillingModel returns the value of the 'cluster_billing_model' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetClusterBillingModel() (value BillingModel, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.clusterBillingModel
	}
	return
}

// ConsoleURL returns the value of the 'console_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ConsoleURL() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.consoleURL
	}
	return ""
}

// GetConsoleURL returns the value of the 'console_URL' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetConsoleURL() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.consoleURL
	}
	return
}

// ConsumerUUID returns the value of the 'consumer_UUID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ConsumerUUID() string {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.consumerUUID
	}
	return ""
}

// GetConsumerUUID returns the value of the 'consumer_UUID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetConsumerUUID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.consumerUUID
	}
	return
}

// CpuTotal returns the value of the 'cpu_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) CpuTotal() int {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.cpuTotal
	}
	return 0
}

// GetCpuTotal returns the value of the 'cpu_total' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetCpuTotal() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.cpuTotal
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.createdAt
	}
	return
}

// CreatorId returns the value of the 'creator_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) CreatorId() string {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.creatorId
	}
	return ""
}

// GetCreatorId returns the value of the 'creator_id' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetCreatorId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.creatorId
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) DisplayName() string {
	if o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14] {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetDisplayName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14]
	if ok {
		value = o.displayName
	}
	return
}

// ExternalClusterID returns the value of the 'external_cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ExternalClusterID() string {
	if o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15] {
		return o.externalClusterID
	}
	return ""
}

// GetExternalClusterID returns the value of the 'external_cluster_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetExternalClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15]
	if ok {
		value = o.externalClusterID
	}
	return
}

// LastReconcileDate returns the value of the 'last_reconcile_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) LastReconcileDate() time.Time {
	if o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16] {
		return o.lastReconcileDate
	}
	return time.Time{}
}

// GetLastReconcileDate returns the value of the 'last_reconcile_date' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetLastReconcileDate() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16]
	if ok {
		value = o.lastReconcileDate
	}
	return
}

// LastReleasedAt returns the value of the 'last_released_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) LastReleasedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17] {
		return o.lastReleasedAt
	}
	return time.Time{}
}

// GetLastReleasedAt returns the value of the 'last_released_at' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetLastReleasedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17]
	if ok {
		value = o.lastReleasedAt
	}
	return
}

// LastTelemetryDate returns the value of the 'last_telemetry_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) LastTelemetryDate() time.Time {
	if o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18] {
		return o.lastTelemetryDate
	}
	return time.Time{}
}

// GetLastTelemetryDate returns the value of the 'last_telemetry_date' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetLastTelemetryDate() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18]
	if ok {
		value = o.lastTelemetryDate
	}
	return
}

// Managed returns the value of the 'managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Managed() bool {
	if o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19] {
		return o.managed
	}
	return false
}

// GetManaged returns the value of the 'managed' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetManaged() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19]
	if ok {
		value = o.managed
	}
	return
}

// Metrics returns the value of the 'metrics' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Metrics() string {
	if o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20] {
		return o.metrics
	}
	return ""
}

// GetMetrics returns the value of the 'metrics' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetMetrics() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20]
	if ok {
		value = o.metrics
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) OrganizationID() string {
	if o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21] {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21]
	if ok {
		value = o.organizationID
	}
	return
}

// PlanID returns the value of the 'plan_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) PlanID() string {
	if o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22] {
		return o.planID
	}
	return ""
}

// GetPlanID returns the value of the 'plan_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetPlanID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22]
	if ok {
		value = o.planID
	}
	return
}

// ProductBundle returns the value of the 'product_bundle' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ProductBundle() string {
	if o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23] {
		return o.productBundle
	}
	return ""
}

// GetProductBundle returns the value of the 'product_bundle' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetProductBundle() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23]
	if ok {
		value = o.productBundle
	}
	return
}

// Provenance returns the value of the 'provenance' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Provenance() string {
	if o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24] {
		return o.provenance
	}
	return ""
}

// GetProvenance returns the value of the 'provenance' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetProvenance() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24]
	if ok {
		value = o.provenance
	}
	return
}

// QueryTimestamp returns the value of the 'query_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) QueryTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25] {
		return o.queryTimestamp
	}
	return time.Time{}
}

// GetQueryTimestamp returns the value of the 'query_timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetQueryTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25]
	if ok {
		value = o.queryTimestamp
	}
	return
}

// RegionID returns the value of the 'region_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) RegionID() string {
	if o != nil && len(o.fieldSet_) > 26 && o.fieldSet_[26] {
		return o.regionID
	}
	return ""
}

// GetRegionID returns the value of the 'region_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetRegionID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 26 && o.fieldSet_[26]
	if ok {
		value = o.regionID
	}
	return
}

// Released returns the value of the 'released' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Released() bool {
	if o != nil && len(o.fieldSet_) > 27 && o.fieldSet_[27] {
		return o.released
	}
	return false
}

// GetReleased returns the value of the 'released' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetReleased() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 27 && o.fieldSet_[27]
	if ok {
		value = o.released
	}
	return
}

// ServiceLevel returns the value of the 'service_level' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) ServiceLevel() string {
	if o != nil && len(o.fieldSet_) > 28 && o.fieldSet_[28] {
		return o.serviceLevel
	}
	return ""
}

// GetServiceLevel returns the value of the 'service_level' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetServiceLevel() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 28 && o.fieldSet_[28]
	if ok {
		value = o.serviceLevel
	}
	return
}

// SocketTotal returns the value of the 'socket_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) SocketTotal() int {
	if o != nil && len(o.fieldSet_) > 29 && o.fieldSet_[29] {
		return o.socketTotal
	}
	return 0
}

// GetSocketTotal returns the value of the 'socket_total' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetSocketTotal() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 29 && o.fieldSet_[29]
	if ok {
		value = o.socketTotal
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Status() string {
	if o != nil && len(o.fieldSet_) > 30 && o.fieldSet_[30] {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetStatus() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 30 && o.fieldSet_[30]
	if ok {
		value = o.status
	}
	return
}

// SupportLevel returns the value of the 'support_level' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) SupportLevel() string {
	if o != nil && len(o.fieldSet_) > 31 && o.fieldSet_[31] {
		return o.supportLevel
	}
	return ""
}

// GetSupportLevel returns the value of the 'support_level' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetSupportLevel() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 31 && o.fieldSet_[31]
	if ok {
		value = o.supportLevel
	}
	return
}

// SystemUnits returns the value of the 'system_units' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) SystemUnits() string {
	if o != nil && len(o.fieldSet_) > 32 && o.fieldSet_[32] {
		return o.systemUnits
	}
	return ""
}

// GetSystemUnits returns the value of the 'system_units' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetSystemUnits() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 32 && o.fieldSet_[32]
	if ok {
		value = o.systemUnits
	}
	return
}

// TrialEndDate returns the value of the 'trial_end_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) TrialEndDate() time.Time {
	if o != nil && len(o.fieldSet_) > 33 && o.fieldSet_[33] {
		return o.trialEndDate
	}
	return time.Time{}
}

// GetTrialEndDate returns the value of the 'trial_end_date' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetTrialEndDate() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 33 && o.fieldSet_[33]
	if ok {
		value = o.trialEndDate
	}
	return
}

// Usage returns the value of the 'usage' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *DeletedSubscription) Usage() string {
	if o != nil && len(o.fieldSet_) > 34 && o.fieldSet_[34] {
		return o.usage
	}
	return ""
}

// GetUsage returns the value of the 'usage' attribute and
// a flag indicating if the attribute has a value.
func (o *DeletedSubscription) GetUsage() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 34 && o.fieldSet_[34]
	if ok {
		value = o.usage
	}
	return
}

// DeletedSubscriptionListKind is the name of the type used to represent list of objects of
// type 'deleted_subscription'.
const DeletedSubscriptionListKind = "DeletedSubscriptionList"

// DeletedSubscriptionListLinkKind is the name of the type used to represent links to list
// of objects of type 'deleted_subscription'.
const DeletedSubscriptionListLinkKind = "DeletedSubscriptionListLink"

// DeletedSubscriptionNilKind is the name of the type used to nil lists of objects of
// type 'deleted_subscription'.
const DeletedSubscriptionListNilKind = "DeletedSubscriptionListNil"

// DeletedSubscriptionList is a list of values of the 'deleted_subscription' type.
type DeletedSubscriptionList struct {
	href  string
	link  bool
	items []*DeletedSubscription
}

// Kind returns the name of the type of the object.
func (l *DeletedSubscriptionList) Kind() string {
	if l == nil {
		return DeletedSubscriptionListNilKind
	}
	if l.link {
		return DeletedSubscriptionListLinkKind
	}
	return DeletedSubscriptionListKind
}

// Link returns true iif this is a link.
func (l *DeletedSubscriptionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *DeletedSubscriptionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *DeletedSubscriptionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *DeletedSubscriptionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *DeletedSubscriptionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *DeletedSubscriptionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *DeletedSubscriptionList) SetItems(items []*DeletedSubscription) {
	l.items = items
}

// Items returns the items of the list.
func (l *DeletedSubscriptionList) Items() []*DeletedSubscription {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *DeletedSubscriptionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DeletedSubscriptionList) Get(i int) *DeletedSubscription {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *DeletedSubscriptionList) Slice() []*DeletedSubscription {
	var slice []*DeletedSubscription
	if l == nil {
		slice = make([]*DeletedSubscription, 0)
	} else {
		slice = make([]*DeletedSubscription, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DeletedSubscriptionList) Each(f func(item *DeletedSubscription) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *DeletedSubscriptionList) Range(f func(index int, item *DeletedSubscription) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
