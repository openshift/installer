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

// BillingModel represents the values of the 'billing_model' enumerated type.
type BillingModel string

const (
	// BillingModelMarketplace Legacy Marketplace billing model. Currently only used for tests. Use cloud-provider specific billing models instead.
	BillingModelMarketplace BillingModel = "marketplace"
	// AWS Marketplace billing model.
	BillingModelMarketplaceAWS BillingModel = "marketplace-aws"
	// GCP Marketplace billing model.
	BillingModelMarketplaceGCP BillingModel = "marketplace-gcp"
	// RH Marketplace billing model.
	BillingModelMarketplaceRHM BillingModel = "marketplace-rhm"
	// Azure Marketplace billing model.
	BillingModelMarketplaceAzure BillingModel = "marketplace-azure"
	// Standard. This is the default billing model
	BillingModelStandard BillingModel = "standard"
)
