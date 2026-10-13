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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1"
)

// SubscriptionStatus represents the values of the 'subscription_status' enumerated type.
type SubscriptionStatus = api_v1.SubscriptionStatus

const (
	// Active state subscriptions have their related resources currently running and reporting an active state
	// Whether a subscription is active is determined depending on the plan of the subscription
	// For example, OCP subscriptions are active if the OCP cluster is successfully reporting metrics
	//              RHOSR subscriptions are active if the service-registry service determines they are, the service updates the subscription as necessary
	SubscriptionStatusActive SubscriptionStatus = api_v1.SubscriptionStatusActive
	// Subscriptions move to Archived when the resources are no longer visibile to OCM and suspected removed
	// Users can also move some disconnected subscriptions to archived state
	// If a subscription in Archived state's resources start reporting again, the subscription may move back to Active
	SubscriptionStatusArchived SubscriptionStatus = api_v1.SubscriptionStatusArchived
	// Deprovisioned subscriptions can be considered completely deleted. As of this writing, only managed plan subscriptions are completely
	// deleted. Instead of actual DB row deletion, subscriptions are moved to Deprovisioned status and all associated resources (quota,
	// roles, etc) are _actually_ deleted. This allows us to keep track of what subscriptions existed and when.
	SubscriptionStatusDeprovisioned SubscriptionStatus = api_v1.SubscriptionStatusDeprovisioned
	// Disconnected subscriptions are Active subscriptions that are intentionally not reporting an active state. There may be some
	// desire by the subscription owner not to connect the resources to OCM. This status allows the subscription to stay in OCM without
	// automatically moving to Stale or Archived.
	SubscriptionStatusDisconnected SubscriptionStatus = api_v1.SubscriptionStatusDisconnected
	// Reserved subscriptions are created during the resource installation process. A reserved subscription represents a subscription
	// whose resources do not yet exist, but are expected to exist soon. Creating a reserved subscription allows services to reserve quota
	// for resources that are in the process of creation. Services are expected to update the status to Active or Deprovisioned once
	// the creation process completes, or fails.
	SubscriptionStatusReserved SubscriptionStatus = api_v1.SubscriptionStatusReserved
	// Stale subscriptions are active subscriptions who have stopped reporting an active state. Once reports cease, the subscription
	// is moved to stale to indicate to users that OCM can no longer see the Active resources. Subscriptions in stale state will automatically
	// transition back to active if the resources stat reporting again. They will also transition to Disconnected or Archived if the
	// resources never resume reporting.
	SubscriptionStatusStale SubscriptionStatus = api_v1.SubscriptionStatusStale
)
