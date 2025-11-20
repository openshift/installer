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
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1"
)

// BillingModel represents the values of the 'billing_model' enumerated type.
type BillingModel = api_v1.BillingModel

const (
	// BillingModelMarketplace Legacy Marketplace billing model. Currently only used for tests. Use cloud-provider specific billing models instead.
	BillingModelMarketplace BillingModel = api_v1.BillingModelMarketplace
	// AWS Marketplace billing model.
	BillingModelMarketplaceAWS BillingModel = api_v1.BillingModelMarketplaceAWS
	// GCP Marketplace billing model.
	BillingModelMarketplaceGCP BillingModel = api_v1.BillingModelMarketplaceGCP
	// RH Marketplace billing model.
	BillingModelMarketplaceRHM BillingModel = api_v1.BillingModelMarketplaceRHM
	// Azure Marketplace billing model.
	BillingModelMarketplaceAzure BillingModel = api_v1.BillingModelMarketplaceAzure
	// Standard. This is the default billing model
	BillingModelStandard BillingModel = api_v1.BillingModelStandard
)
