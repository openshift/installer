/**
 * (C) Copyright IBM Corp. 2024.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * IBM OpenAPI SDK Code Generator Version: 3.89.1-ed9d96f4-20240417-193115
 */

// Package enterprisemanagementv1 : Operations and models for the EnterpriseManagementV1 service
package enterprisemanagementv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// EnterpriseManagementV1 : The Enterprise Management API enables you to create and manage an enterprise, account
// groups, and accounts within the enterprise.
//
// API Version: 1.0
type EnterpriseManagementV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://enterprise.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "enterprise_management"

// EnterpriseManagementV1Options : Service options
type EnterpriseManagementV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewEnterpriseManagementV1UsingExternalConfig : constructs an instance of EnterpriseManagementV1 with passed in options and external configuration.
func NewEnterpriseManagementV1UsingExternalConfig(options *EnterpriseManagementV1Options) (enterpriseManagement *EnterpriseManagementV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	enterpriseManagement, err = NewEnterpriseManagementV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = enterpriseManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = enterpriseManagement.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewEnterpriseManagementV1 : constructs an instance of EnterpriseManagementV1 with passed in options.
func NewEnterpriseManagementV1(options *EnterpriseManagementV1Options) (service *EnterpriseManagementV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "new-base-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-url-error", common.GetComponentInfo())
			return
		}
	}

	service = &EnterpriseManagementV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "enterpriseManagement" suitable for processing requests.
func (enterpriseManagement *EnterpriseManagementV1) Clone() *EnterpriseManagementV1 {
	if core.IsNil(enterpriseManagement) {
		return nil
	}
	clone := *enterpriseManagement
	clone.Service = enterpriseManagement.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (enterpriseManagement *EnterpriseManagementV1) SetServiceURL(url string) error {
	err := enterpriseManagement.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (enterpriseManagement *EnterpriseManagementV1) GetServiceURL() string {
	return enterpriseManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (enterpriseManagement *EnterpriseManagementV1) SetDefaultHeaders(headers http.Header) {
	enterpriseManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (enterpriseManagement *EnterpriseManagementV1) SetEnableGzipCompression(enableGzip bool) {
	enterpriseManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (enterpriseManagement *EnterpriseManagementV1) GetEnableGzipCompression() bool {
	return enterpriseManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (enterpriseManagement *EnterpriseManagementV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	enterpriseManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (enterpriseManagement *EnterpriseManagementV1) DisableRetries() {
	enterpriseManagement.Service.DisableRetries()
}

// CreateEnterprise : Create an enterprise
// Create a new enterprise, which you can use to centrally manage multiple accounts. To create an enterprise, you must
// have an active Subscription account. <br/><br/>The API creates an enterprise entity, which is the root of the
// enterprise hierarchy. It also creates a new enterprise account that is used to manage the enterprise. All
// subscriptions, support entitlements, credits, and discounts from the source subscription account are migrated to the
// enterprise account, and the source account becomes a child account in the hierarchy. The user that you assign as the
// enterprise primary contact is also assigned as the owner of the enterprise account.
func (enterpriseManagement *EnterpriseManagementV1) CreateEnterprise(createEnterpriseOptions *CreateEnterpriseOptions) (result *CreateEnterpriseResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.CreateEnterpriseWithContext(context.Background(), createEnterpriseOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateEnterpriseWithContext is an alternate form of the CreateEnterprise method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) CreateEnterpriseWithContext(ctx context.Context, createEnterpriseOptions *CreateEnterpriseOptions) (result *CreateEnterpriseResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createEnterpriseOptions, "createEnterpriseOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createEnterpriseOptions, "createEnterpriseOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/enterprises`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createEnterpriseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "CreateEnterprise")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createEnterpriseOptions.SourceAccountID != nil {
		body["source_account_id"] = createEnterpriseOptions.SourceAccountID
	}
	if createEnterpriseOptions.Name != nil {
		body["name"] = createEnterpriseOptions.Name
	}
	if createEnterpriseOptions.PrimaryContactIamID != nil {
		body["primary_contact_iam_id"] = createEnterpriseOptions.PrimaryContactIamID
	}
	if createEnterpriseOptions.Domain != nil {
		body["domain"] = createEnterpriseOptions.Domain
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_enterprise", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateEnterpriseResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListEnterprises : List enterprises
// Retrieve all enterprises for a given ID by passing the IDs on query parameters. If no ID is passed, the enterprises
// for which the calling identity is the primary contact are returned. You can use pagination parameters to filter the
// results. <br/><br/>This method ensures that only the enterprises that the user has access to are returned. Access can
// be controlled either through a policy on a specific enterprise, or account-level platform services access roles, such
// as Administrator, Editor, Operator, or Viewer. When you call the method with the `enterprise_account_id` or
// `account_id` query parameter, the account ID in the token is compared with that in the query parameter. If these
// account IDs match, authentication isn't performed and the enterprise information is returned. If the account IDs
// don't match, authentication is performed and only then is the enterprise information returned in the response.
func (enterpriseManagement *EnterpriseManagementV1) ListEnterprises(listEnterprisesOptions *ListEnterprisesOptions) (result *ListEnterprisesResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.ListEnterprisesWithContext(context.Background(), listEnterprisesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListEnterprisesWithContext is an alternate form of the ListEnterprises method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) ListEnterprisesWithContext(ctx context.Context, listEnterprisesOptions *ListEnterprisesOptions) (result *ListEnterprisesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listEnterprisesOptions, "listEnterprisesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/enterprises`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listEnterprisesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "ListEnterprises")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listEnterprisesOptions.EnterpriseAccountID != nil {
		builder.AddQuery("enterprise_account_id", fmt.Sprint(*listEnterprisesOptions.EnterpriseAccountID))
	}
	if listEnterprisesOptions.AccountGroupID != nil {
		builder.AddQuery("account_group_id", fmt.Sprint(*listEnterprisesOptions.AccountGroupID))
	}
	if listEnterprisesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listEnterprisesOptions.AccountID))
	}
	if listEnterprisesOptions.NextDocid != nil {
		builder.AddQuery("next_docid", fmt.Sprint(*listEnterprisesOptions.NextDocid))
	}
	if listEnterprisesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listEnterprisesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_enterprises", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListEnterprisesResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetEnterprise : Get enterprise by ID
// Retrieve an enterprise by the `enterprise_id` parameter. All data related to the enterprise is returned only if the
// caller has access to retrieve the enterprise.
func (enterpriseManagement *EnterpriseManagementV1) GetEnterprise(getEnterpriseOptions *GetEnterpriseOptions) (result *Enterprise, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.GetEnterpriseWithContext(context.Background(), getEnterpriseOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetEnterpriseWithContext is an alternate form of the GetEnterprise method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) GetEnterpriseWithContext(ctx context.Context, getEnterpriseOptions *GetEnterpriseOptions) (result *Enterprise, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEnterpriseOptions, "getEnterpriseOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getEnterpriseOptions, "getEnterpriseOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"enterprise_id": *getEnterpriseOptions.EnterpriseID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/enterprises/{enterprise_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getEnterpriseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "GetEnterprise")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_enterprise", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnterprise)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateEnterprise : Update an enterprise
// Update the name, domain, or IAM ID of the primary contact for an existing enterprise. The new primary contact must
// already be a user in the enterprise account.
func (enterpriseManagement *EnterpriseManagementV1) UpdateEnterprise(updateEnterpriseOptions *UpdateEnterpriseOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.UpdateEnterpriseWithContext(context.Background(), updateEnterpriseOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateEnterpriseWithContext is an alternate form of the UpdateEnterprise method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) UpdateEnterpriseWithContext(ctx context.Context, updateEnterpriseOptions *UpdateEnterpriseOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEnterpriseOptions, "updateEnterpriseOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateEnterpriseOptions, "updateEnterpriseOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"enterprise_id": *updateEnterpriseOptions.EnterpriseID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/enterprises/{enterprise_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateEnterpriseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "UpdateEnterprise")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateEnterpriseOptions.Name != nil {
		body["name"] = updateEnterpriseOptions.Name
	}
	if updateEnterpriseOptions.Domain != nil {
		body["domain"] = updateEnterpriseOptions.Domain
	}
	if updateEnterpriseOptions.PrimaryContactIamID != nil {
		body["primary_contact_iam_id"] = updateEnterpriseOptions.PrimaryContactIamID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_enterprise", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ImportAccountToEnterprise : Import an account into an enterprise
// Import an existing stand-alone account into an enterprise. The existing account can be any type: trial (`TRIAL`),
// Lite (`STANDARD`), Pay-As-You-Go (`PAYG`), or Subscription (`SUBSCRIPTION`). In the case of a `SUBSCRIPTION` account,
// the credits, promotional offers, and discounts are migrated to the billing unit of the enterprise. For a billable
// account (`PAYG` or `SUBSCRIPTION`), the country and currency code of the existing account and the billing unit of the
// enterprise must match. The API returns a `202` response and performs asynchronous operations to import the account
// into the enterprise. <br/></br>For more information about impacts to the account, see [Adding accounts to an
// enterprise](https://{DomainName}/docs/account?topic=account-enterprise-add).
func (enterpriseManagement *EnterpriseManagementV1) ImportAccountToEnterprise(importAccountToEnterpriseOptions *ImportAccountToEnterpriseOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.ImportAccountToEnterpriseWithContext(context.Background(), importAccountToEnterpriseOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ImportAccountToEnterpriseWithContext is an alternate form of the ImportAccountToEnterprise method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) ImportAccountToEnterpriseWithContext(ctx context.Context, importAccountToEnterpriseOptions *ImportAccountToEnterpriseOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importAccountToEnterpriseOptions, "importAccountToEnterpriseOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(importAccountToEnterpriseOptions, "importAccountToEnterpriseOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"enterprise_id": *importAccountToEnterpriseOptions.EnterpriseID,
		"account_id":    *importAccountToEnterpriseOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/enterprises/{enterprise_id}/import/accounts/{account_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range importAccountToEnterpriseOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "ImportAccountToEnterprise")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if importAccountToEnterpriseOptions.Parent != nil {
		body["parent"] = importAccountToEnterpriseOptions.Parent
	}
	if importAccountToEnterpriseOptions.BillingUnitID != nil {
		body["billing_unit_id"] = importAccountToEnterpriseOptions.BillingUnitID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "import_account_to_enterprise", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateAccount : Create a new account in an enterprise
// Create a new account as a part of an existing enterprise. The API creates an account entity under the parent that is
// specified in the payload of the request. The request also takes in the name and the owner of this new account. The
// owner must have a valid IBMid that's registered with IBM Cloud, but they don't need to be a user in the enterprise
// account.
func (enterpriseManagement *EnterpriseManagementV1) CreateAccount(createAccountOptions *CreateAccountOptions) (result *CreateAccountResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.CreateAccountWithContext(context.Background(), createAccountOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAccountWithContext is an alternate form of the CreateAccount method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) CreateAccountWithContext(ctx context.Context, createAccountOptions *CreateAccountOptions) (result *CreateAccountResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccountOptions, "createAccountOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAccountOptions, "createAccountOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/accounts`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "CreateAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccountOptions.Parent != nil {
		body["parent"] = createAccountOptions.Parent
	}
	if createAccountOptions.Name != nil {
		body["name"] = createAccountOptions.Name
	}
	if createAccountOptions.OwnerIamID != nil {
		body["owner_iam_id"] = createAccountOptions.OwnerIamID
	}
	if createAccountOptions.Traits != nil {
		body["traits"] = createAccountOptions.Traits
	}
	if createAccountOptions.Options != nil {
		body["options"] = createAccountOptions.Options
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_account", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateAccountResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAccounts : List accounts
// Retrieve all accounts based on the values that are passed in the query parameters. If no query parameter is passed,
// all of the accounts in the enterprise for which the calling identity has access are returned. <br/><br/>You can use
// pagination parameters to filter the results. The `limit` field can be used to limit the number of results that are
// displayed for this method.<br/><br/>This method ensures that only the accounts that the user has access to are
// returned. Access can be controlled either through a policy on a specific account, or account-level platform services
// access roles, such as Administrator, Editor, Operator, or Viewer. When you call the method with the `enterprise_id`,
// `account_group_id` or `parent` query parameter, all of the accounts that are immediate children of this entity are
// returned. Authentication is performed on all the accounts before they are returned to the user to ensure that only
// those accounts are returned to which the calling identity has access to.
func (enterpriseManagement *EnterpriseManagementV1) ListAccounts(listAccountsOptions *ListAccountsOptions) (result *ListAccountsResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.ListAccountsWithContext(context.Background(), listAccountsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAccountsWithContext is an alternate form of the ListAccounts method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) ListAccountsWithContext(ctx context.Context, listAccountsOptions *ListAccountsOptions) (result *ListAccountsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAccountsOptions, "listAccountsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/accounts`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAccountsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "ListAccounts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAccountsOptions.EnterpriseID != nil {
		builder.AddQuery("enterprise_id", fmt.Sprint(*listAccountsOptions.EnterpriseID))
	}
	if listAccountsOptions.AccountGroupID != nil {
		builder.AddQuery("account_group_id", fmt.Sprint(*listAccountsOptions.AccountGroupID))
	}
	if listAccountsOptions.NextDocid != nil {
		builder.AddQuery("next_docid", fmt.Sprint(*listAccountsOptions.NextDocid))
	}
	if listAccountsOptions.Parent != nil {
		builder.AddQuery("parent", fmt.Sprint(*listAccountsOptions.Parent))
	}
	if listAccountsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAccountsOptions.Limit))
	}
	if listAccountsOptions.IncludeDeleted != nil {
		builder.AddQuery("include_deleted", fmt.Sprint(*listAccountsOptions.IncludeDeleted))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_accounts", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAccountsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAccount : Get account by ID
// Retrieve an account by the `account_id` parameter. All data related to the account is returned only if the caller has
// access to retrieve the account.
func (enterpriseManagement *EnterpriseManagementV1) GetAccount(getAccountOptions *GetAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.GetAccountWithContext(context.Background(), getAccountOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountWithContext is an alternate form of the GetAccount method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) GetAccountWithContext(ctx context.Context, getAccountOptions *GetAccountOptions) (result *Account, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountOptions, "getAccountOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountOptions, "getAccountOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/accounts/{account_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "GetAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccount)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccount : Move an account within the enterprise
// Move an account to a different parent within the same enterprise.
func (enterpriseManagement *EnterpriseManagementV1) UpdateAccount(updateAccountOptions *UpdateAccountOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.UpdateAccountWithContext(context.Background(), updateAccountOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAccountWithContext is an alternate form of the UpdateAccount method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) UpdateAccountWithContext(ctx context.Context, updateAccountOptions *UpdateAccountOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountOptions, "updateAccountOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAccountOptions, "updateAccountOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *updateAccountOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/accounts/{account_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "UpdateAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccountOptions.Parent != nil {
		body["parent"] = updateAccountOptions.Parent
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_account", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DeleteAccount : Remove an account from its enterprise
// Remove an account from the enterprise its currently in. After an account is removed, it will be canceled and cannot
// be reactivated.
func (enterpriseManagement *EnterpriseManagementV1) DeleteAccount(deleteAccountOptions *DeleteAccountOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.DeleteAccountWithContext(context.Background(), deleteAccountOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAccountWithContext is an alternate form of the DeleteAccount method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) DeleteAccountWithContext(ctx context.Context, deleteAccountOptions *DeleteAccountOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccountOptions, "deleteAccountOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAccountOptions, "deleteAccountOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *deleteAccountOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/accounts/{account_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "DeleteAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_account", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateAccountGroup : Create an account group
// Create a new account group, which can be used to group together multiple accounts. To create an account group, you
// must have an existing enterprise. The API creates an account group entity under the parent that is specified in the
// payload of the request. The request also takes in the name and the primary contact of this new account group.
func (enterpriseManagement *EnterpriseManagementV1) CreateAccountGroup(createAccountGroupOptions *CreateAccountGroupOptions) (result *CreateAccountGroupResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.CreateAccountGroupWithContext(context.Background(), createAccountGroupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAccountGroupWithContext is an alternate form of the CreateAccountGroup method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) CreateAccountGroupWithContext(ctx context.Context, createAccountGroupOptions *CreateAccountGroupOptions) (result *CreateAccountGroupResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccountGroupOptions, "createAccountGroupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAccountGroupOptions, "createAccountGroupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/account-groups`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAccountGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "CreateAccountGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccountGroupOptions.Parent != nil {
		body["parent"] = createAccountGroupOptions.Parent
	}
	if createAccountGroupOptions.Name != nil {
		body["name"] = createAccountGroupOptions.Name
	}
	if createAccountGroupOptions.PrimaryContactIamID != nil {
		body["primary_contact_iam_id"] = createAccountGroupOptions.PrimaryContactIamID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_account_group", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateAccountGroupResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAccountGroups : List account groups
// Retrieve all account groups based on the values that are passed in the query parameters. If no query parameter is
// passed, all of the account groups in the enterprise for which the calling identity has access are returned.
// <br/><br/>You can use pagination parameters to filter the results. The `limit` field can be used to limit the number
// of results that are displayed for this method.<br/><br/>This method ensures that only the account groups that the
// user has access to are returned. Access can be controlled either through a policy on a specific account group, or
// account-level platform services access roles, such as Administrator, Editor, Operator, or Viewer. When you call the
// method with the `enterprise_id`, `parent_account_group_id` or `parent` query parameter, all of the account groups
// that are immediate children of this entity are returned. Authentication is performed on all account groups before
// they are returned to the user to ensure that only those account groups are returned to which the calling identity has
// access.
func (enterpriseManagement *EnterpriseManagementV1) ListAccountGroups(listAccountGroupsOptions *ListAccountGroupsOptions) (result *ListAccountGroupsResponse, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.ListAccountGroupsWithContext(context.Background(), listAccountGroupsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAccountGroupsWithContext is an alternate form of the ListAccountGroups method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) ListAccountGroupsWithContext(ctx context.Context, listAccountGroupsOptions *ListAccountGroupsOptions) (result *ListAccountGroupsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAccountGroupsOptions, "listAccountGroupsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/account-groups`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAccountGroupsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "ListAccountGroups")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAccountGroupsOptions.EnterpriseID != nil {
		builder.AddQuery("enterprise_id", fmt.Sprint(*listAccountGroupsOptions.EnterpriseID))
	}
	if listAccountGroupsOptions.ParentAccountGroupID != nil {
		builder.AddQuery("parent_account_group_id", fmt.Sprint(*listAccountGroupsOptions.ParentAccountGroupID))
	}
	if listAccountGroupsOptions.NextDocid != nil {
		builder.AddQuery("next_docid", fmt.Sprint(*listAccountGroupsOptions.NextDocid))
	}
	if listAccountGroupsOptions.Parent != nil {
		builder.AddQuery("parent", fmt.Sprint(*listAccountGroupsOptions.Parent))
	}
	if listAccountGroupsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAccountGroupsOptions.Limit))
	}
	if listAccountGroupsOptions.IncludeDeleted != nil {
		builder.AddQuery("include_deleted", fmt.Sprint(*listAccountGroupsOptions.IncludeDeleted))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_account_groups", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAccountGroupsResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAccountGroup : Get account group by ID
// Retrieve an account by the `account_group_id` parameter. All data related to the account group is returned only if
// the caller has access to retrieve the account group.
func (enterpriseManagement *EnterpriseManagementV1) GetAccountGroup(getAccountGroupOptions *GetAccountGroupOptions) (result *AccountGroup, response *core.DetailedResponse, err error) {
	result, response, err = enterpriseManagement.GetAccountGroupWithContext(context.Background(), getAccountGroupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountGroupWithContext is an alternate form of the GetAccountGroup method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) GetAccountGroupWithContext(ctx context.Context, getAccountGroupOptions *GetAccountGroupOptions) (result *AccountGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountGroupOptions, "getAccountGroupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountGroupOptions, "getAccountGroupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_group_id": *getAccountGroupOptions.AccountGroupID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/account-groups/{account_group_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "GetAccountGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = enterpriseManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_group", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountGroup)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccountGroup : Update an account group
// Update the name or IAM ID of the primary contact for an existing account group. The new primary contact must already
// be a user in the enterprise account.
func (enterpriseManagement *EnterpriseManagementV1) UpdateAccountGroup(updateAccountGroupOptions *UpdateAccountGroupOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.UpdateAccountGroupWithContext(context.Background(), updateAccountGroupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAccountGroupWithContext is an alternate form of the UpdateAccountGroup method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) UpdateAccountGroupWithContext(ctx context.Context, updateAccountGroupOptions *UpdateAccountGroupOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccountGroupOptions, "updateAccountGroupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAccountGroupOptions, "updateAccountGroupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_group_id": *updateAccountGroupOptions.AccountGroupID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/account-groups/{account_group_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAccountGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "UpdateAccountGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccountGroupOptions.Name != nil {
		body["name"] = updateAccountGroupOptions.Name
	}
	if updateAccountGroupOptions.PrimaryContactIamID != nil {
		body["primary_contact_iam_id"] = updateAccountGroupOptions.PrimaryContactIamID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_account_group", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DeleteAccountGroup : Delete an account group from the enterprise
// Delete an existing account group from the enterprise. You can't delete an account group that has child account
// groups, the delete request will fail. This API doesn't perform a recursive delete on the child account groups, it
// only deletes the current account group.
func (enterpriseManagement *EnterpriseManagementV1) DeleteAccountGroup(deleteAccountGroupOptions *DeleteAccountGroupOptions) (response *core.DetailedResponse, err error) {
	response, err = enterpriseManagement.DeleteAccountGroupWithContext(context.Background(), deleteAccountGroupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAccountGroupWithContext is an alternate form of the DeleteAccountGroup method which supports a Context parameter
func (enterpriseManagement *EnterpriseManagementV1) DeleteAccountGroupWithContext(ctx context.Context, deleteAccountGroupOptions *DeleteAccountGroupOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccountGroupOptions, "deleteAccountGroupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAccountGroupOptions, "deleteAccountGroupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_group_id": *deleteAccountGroupOptions.AccountGroupID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = enterpriseManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(enterpriseManagement.Service.Options.URL, `/account-groups/{account_group_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAccountGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("enterprise_management", "V1", "DeleteAccountGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = enterpriseManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_account_group", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0")
}

// Account : An account resource.
type Account struct {
	// The URL of the account.
	URL *string `json:"url,omitempty"`

	// The account ID.
	ID *string `json:"id,omitempty"`

	// The Cloud Resource Name (CRN) of the account.
	CRN *string `json:"crn,omitempty"`

	// The CRN of the parent of the account.
	Parent *string `json:"parent,omitempty"`

	// The enterprise account ID.
	EnterpriseAccountID *string `json:"enterprise_account_id,omitempty"`

	// The enterprise ID that the account is a part of.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// The path from the enterprise to this particular account.
	EnterprisePath *string `json:"enterprise_path,omitempty"`

	// The name of the account.
	Name *string `json:"name,omitempty"`

	// The state of the account.
	State *string `json:"state,omitempty"`

	// The IAM ID of the owner of the account.
	OwnerIamID *string `json:"owner_iam_id,omitempty"`

	// The type of account - whether it is free or paid.
	Paid *bool `json:"paid,omitempty"`

	// The email address of the owner of the account.
	OwnerEmail *string `json:"owner_email,omitempty"`

	// The flag to indicate whether the account is an enterprise account or not.
	IsEnterpriseAccount *bool `json:"is_enterprise_account,omitempty"`

	// The time stamp at which the account was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The IAM ID of the user or service that created the account.
	CreatedBy *string `json:"created_by,omitempty"`

	// The time stamp at which the account was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The IAM ID of the user or service that updated the account.
	UpdatedBy *string `json:"updated_by,omitempty"`
}

// UnmarshalAccount unmarshals an instance of Account from the specified map of raw messages.
func UnmarshalAccount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Account)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent", &obj.Parent)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_account_id", &obj.EnterpriseAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_id", &obj.EnterpriseID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_path", &obj.EnterprisePath)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "owner_iam_id", &obj.OwnerIamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "owner_iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "paid", &obj.Paid)
	if err != nil {
		err = core.SDKErrorf(err, "", "paid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "owner_email", &obj.OwnerEmail)
	if err != nil {
		err = core.SDKErrorf(err, "", "owner_email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_enterprise_account", &obj.IsEnterpriseAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_enterprise_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountGroup : An account group resource.
type AccountGroup struct {
	// The URL of the account group.
	URL *string `json:"url,omitempty"`

	// The account group ID.
	ID *string `json:"id,omitempty"`

	// The Cloud Resource Name (CRN) of the account group.
	CRN *string `json:"crn,omitempty"`

	// The CRN of the parent of the account group.
	Parent *string `json:"parent,omitempty"`

	// The enterprise account ID.
	EnterpriseAccountID *string `json:"enterprise_account_id,omitempty"`

	// The enterprise ID that the account group is a part of.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// The path from the enterprise to this particular account group.
	EnterprisePath *string `json:"enterprise_path,omitempty"`

	// The name of the account group.
	Name *string `json:"name,omitempty"`

	// The state of the account group.
	State *string `json:"state,omitempty"`

	// The IAM ID of the primary contact of the account group.
	PrimaryContactIamID *string `json:"primary_contact_iam_id,omitempty"`

	// The email address of the primary contact of the account group.
	PrimaryContactEmail *string `json:"primary_contact_email,omitempty"`

	// The time stamp at which the account group was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The IAM ID of the user or service that created the account group.
	CreatedBy *string `json:"created_by,omitempty"`

	// The time stamp at which the account group was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The IAM ID of the user or service that updated the account group.
	UpdatedBy *string `json:"updated_by,omitempty"`
}

// UnmarshalAccountGroup unmarshals an instance of AccountGroup from the specified map of raw messages.
func UnmarshalAccountGroup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountGroup)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent", &obj.Parent)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_account_id", &obj.EnterpriseAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_id", &obj.EnterpriseID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_path", &obj.EnterprisePath)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_contact_iam_id", &obj.PrimaryContactIamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact_iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_contact_email", &obj.PrimaryContactEmail)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact_email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAccountGroupOptions : The CreateAccountGroup options.
type CreateAccountGroupOptions struct {
	// The CRN of the parent under which the account group will be created. The parent can be an existing account group or
	// the enterprise itself.
	Parent *string `json:"parent" validate:"required"`

	// The name of the account group. This field must have 3 - 60 characters.
	Name *string `json:"name" validate:"required"`

	// The IAM ID of the primary contact for this account group, such as `IBMid-0123ABC`. The IAM ID must already exist.
	PrimaryContactIamID *string `json:"primary_contact_iam_id" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAccountGroupOptions : Instantiate CreateAccountGroupOptions
func (*EnterpriseManagementV1) NewCreateAccountGroupOptions(parent string, name string, primaryContactIamID string) *CreateAccountGroupOptions {
	return &CreateAccountGroupOptions{
		Parent:              core.StringPtr(parent),
		Name:                core.StringPtr(name),
		PrimaryContactIamID: core.StringPtr(primaryContactIamID),
	}
}

// SetParent : Allow user to set Parent
func (_options *CreateAccountGroupOptions) SetParent(parent string) *CreateAccountGroupOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccountGroupOptions) SetName(name string) *CreateAccountGroupOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPrimaryContactIamID : Allow user to set PrimaryContactIamID
func (_options *CreateAccountGroupOptions) SetPrimaryContactIamID(primaryContactIamID string) *CreateAccountGroupOptions {
	_options.PrimaryContactIamID = core.StringPtr(primaryContactIamID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccountGroupOptions) SetHeaders(param map[string]string) *CreateAccountGroupOptions {
	options.Headers = param
	return options
}

// CreateAccountGroupResponse : A newly-created account group.
type CreateAccountGroupResponse struct {
	// The ID of the account group entity that was created.
	AccountGroupID *string `json:"account_group_id,omitempty"`
}

// UnmarshalCreateAccountGroupResponse unmarshals an instance of CreateAccountGroupResponse from the specified map of raw messages.
func UnmarshalCreateAccountGroupResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccountGroupResponse)
	err = core.UnmarshalPrimitive(m, "account_group_id", &obj.AccountGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_group_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAccountOptions : The CreateAccount options.
type CreateAccountOptions struct {
	// The CRN of the parent under which the account will be created. The parent can be an existing account group or the
	// enterprise itself.
	Parent *string `json:"parent" validate:"required"`

	// The name of the account. This field must have 3 - 60 characters.
	Name *string `json:"name" validate:"required"`

	// The IAM ID of the account owner, such as `IBMid-0123ABC`. The IAM ID must already exist.
	OwnerIamID *string `json:"owner_iam_id" validate:"required"`

	// The traits object can be used to set properties on child accounts of an enterprise. You can pass a field to opt-out
	// of the default multi-factor authentication setting or enable enterprise-managed IAM when creating a child account in
	// the enterprise. This is an optional field.
	Traits *CreateAccountRequestTraits `json:"traits,omitempty"`

	// The options object can be used to set properties on child accounts of an enterprise. You can pass a field to to
	// create IAM service id with IAM api key when creating a child account in the enterprise. This is an optional field.
	Options *CreateAccountRequestOptions `json:"options,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAccountOptions : Instantiate CreateAccountOptions
func (*EnterpriseManagementV1) NewCreateAccountOptions(parent string, name string, ownerIamID string) *CreateAccountOptions {
	return &CreateAccountOptions{
		Parent:     core.StringPtr(parent),
		Name:       core.StringPtr(name),
		OwnerIamID: core.StringPtr(ownerIamID),
	}
}

// SetParent : Allow user to set Parent
func (_options *CreateAccountOptions) SetParent(parent string) *CreateAccountOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccountOptions) SetName(name string) *CreateAccountOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOwnerIamID : Allow user to set OwnerIamID
func (_options *CreateAccountOptions) SetOwnerIamID(ownerIamID string) *CreateAccountOptions {
	_options.OwnerIamID = core.StringPtr(ownerIamID)
	return _options
}

// SetTraits : Allow user to set Traits
func (_options *CreateAccountOptions) SetTraits(traits *CreateAccountRequestTraits) *CreateAccountOptions {
	_options.Traits = traits
	return _options
}

// SetOptions : Allow user to set Options
func (_options *CreateAccountOptions) SetOptions(options *CreateAccountRequestOptions) *CreateAccountOptions {
	_options.Options = options
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccountOptions) SetHeaders(param map[string]string) *CreateAccountOptions {
	options.Headers = param
	return options
}

// CreateAccountRequestOptions : The options object can be used to set properties on child accounts of an enterprise. You can pass a field to to
// create IAM service id with IAM api key when creating a child account in the enterprise. This is an optional field.
type CreateAccountRequestOptions struct {
	// By default create_iam_service_id_with_apikey_and_owner_policies is turned off for a newly created child account. You
	// can enable this property by passing 'true' in this boolean field. IAM service id has account owner IAM policies and
	// the API key associated with it can generate a token and setup resources in the account. This is an optional field.
	CreateIamServiceIDWithApikeyAndOwnerPolicies *bool `json:"create_iam_service_id_with_apikey_and_owner_policies,omitempty"`
}

// UnmarshalCreateAccountRequestOptions unmarshals an instance of CreateAccountRequestOptions from the specified map of raw messages.
func UnmarshalCreateAccountRequestOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccountRequestOptions)
	err = core.UnmarshalPrimitive(m, "create_iam_service_id_with_apikey_and_owner_policies", &obj.CreateIamServiceIDWithApikeyAndOwnerPolicies)
	if err != nil {
		err = core.SDKErrorf(err, "", "create_iam_service_id_with_apikey_and_owner_policies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAccountRequestTraits : The traits object can be used to set properties on child accounts of an enterprise. You can pass a field to opt-out
// of the default multi-factor authentication setting or enable enterprise-managed IAM when creating a child account in
// the enterprise. This is an optional field.
type CreateAccountRequestTraits struct {
	// By default MFA is set to `NONE_NO_ROPC` on a child account, which disables CLI logins with only a password. To opt
	// out, pass the traits object with the mfa field set to empty string. This is an optional field.
	Mfa *string `json:"mfa,omitempty"`

	// By default enterprise-managed IAM is turned off for a newly created child account. You can enable this property by
	// passing 'true' in this boolean field. Enabling enterprise-managed IAM allows the enterprise account to assign IAM
	// resources, like access groups, trusted profiles, and account settings, to the child account. This is an optional
	// field.
	EnterpriseIamManaged *bool `json:"enterprise_iam_managed,omitempty"`
}

// UnmarshalCreateAccountRequestTraits unmarshals an instance of CreateAccountRequestTraits from the specified map of raw messages.
func UnmarshalCreateAccountRequestTraits(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccountRequestTraits)
	err = core.UnmarshalPrimitive(m, "mfa", &obj.Mfa)
	if err != nil {
		err = core.SDKErrorf(err, "", "mfa-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_iam_managed", &obj.EnterpriseIamManaged)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_iam_managed-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAccountResponse : A newly-created account.
type CreateAccountResponse struct {
	// The ID of the account entity that was created.
	AccountID *string `json:"account_id,omitempty"`

	// The iam_service_id of the account entity that was created.
	IamServiceID *string `json:"iam_service_id,omitempty"`

	// The iam_apikey_id of the account entity that was created.
	IamApikeyID *string `json:"iam_apikey_id,omitempty"`

	// The iam_apikey of the account entity with owner iam policies that was created.
	IamApikey *string `json:"iam_apikey,omitempty"`
}

// UnmarshalCreateAccountResponse unmarshals an instance of CreateAccountResponse from the specified map of raw messages.
func UnmarshalCreateAccountResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccountResponse)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_service_id", &obj.IamServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_apikey_id", &obj.IamApikeyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_apikey_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_apikey", &obj.IamApikey)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_apikey-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateEnterpriseOptions : The CreateEnterprise options.
type CreateEnterpriseOptions struct {
	// The ID of the account that is used to create the enterprise.
	SourceAccountID *string `json:"source_account_id" validate:"required"`

	// The name of the enterprise. This field must have 3 - 60 characters.
	Name *string `json:"name" validate:"required"`

	// The IAM ID of the enterprise primary contact, such as `IBMid-0123ABC`. The IAM ID must already exist.
	PrimaryContactIamID *string `json:"primary_contact_iam_id" validate:"required"`

	// A domain or subdomain for the enterprise, such as `example.com` or `my.example.com`.
	Domain *string `json:"domain,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateEnterpriseOptions : Instantiate CreateEnterpriseOptions
func (*EnterpriseManagementV1) NewCreateEnterpriseOptions(sourceAccountID string, name string, primaryContactIamID string) *CreateEnterpriseOptions {
	return &CreateEnterpriseOptions{
		SourceAccountID:     core.StringPtr(sourceAccountID),
		Name:                core.StringPtr(name),
		PrimaryContactIamID: core.StringPtr(primaryContactIamID),
	}
}

// SetSourceAccountID : Allow user to set SourceAccountID
func (_options *CreateEnterpriseOptions) SetSourceAccountID(sourceAccountID string) *CreateEnterpriseOptions {
	_options.SourceAccountID = core.StringPtr(sourceAccountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateEnterpriseOptions) SetName(name string) *CreateEnterpriseOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPrimaryContactIamID : Allow user to set PrimaryContactIamID
func (_options *CreateEnterpriseOptions) SetPrimaryContactIamID(primaryContactIamID string) *CreateEnterpriseOptions {
	_options.PrimaryContactIamID = core.StringPtr(primaryContactIamID)
	return _options
}

// SetDomain : Allow user to set Domain
func (_options *CreateEnterpriseOptions) SetDomain(domain string) *CreateEnterpriseOptions {
	_options.Domain = core.StringPtr(domain)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateEnterpriseOptions) SetHeaders(param map[string]string) *CreateEnterpriseOptions {
	options.Headers = param
	return options
}

// CreateEnterpriseResponse : The response from calling create enterprise.
type CreateEnterpriseResponse struct {
	// The ID of the enterprise entity that was created. This entity is the root of the hierarchy.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// The ID of the enterprise account that was created. The enterprise account is used to manage billing and access to
	// the enterprise management.
	EnterpriseAccountID *string `json:"enterprise_account_id,omitempty"`
}

// UnmarshalCreateEnterpriseResponse unmarshals an instance of CreateEnterpriseResponse from the specified map of raw messages.
func UnmarshalCreateEnterpriseResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateEnterpriseResponse)
	err = core.UnmarshalPrimitive(m, "enterprise_id", &obj.EnterpriseID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_account_id", &obj.EnterpriseAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_account_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccountGroupOptions : The DeleteAccountGroup options.
type DeleteAccountGroupOptions struct {
	// The ID of the account group to retrieve.
	AccountGroupID *string `json:"account_group_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAccountGroupOptions : Instantiate DeleteAccountGroupOptions
func (*EnterpriseManagementV1) NewDeleteAccountGroupOptions(accountGroupID string) *DeleteAccountGroupOptions {
	return &DeleteAccountGroupOptions{
		AccountGroupID: core.StringPtr(accountGroupID),
	}
}

// SetAccountGroupID : Allow user to set AccountGroupID
func (_options *DeleteAccountGroupOptions) SetAccountGroupID(accountGroupID string) *DeleteAccountGroupOptions {
	_options.AccountGroupID = core.StringPtr(accountGroupID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccountGroupOptions) SetHeaders(param map[string]string) *DeleteAccountGroupOptions {
	options.Headers = param
	return options
}

// DeleteAccountOptions : The DeleteAccount options.
type DeleteAccountOptions struct {
	// The ID of the target account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAccountOptions : Instantiate DeleteAccountOptions
func (*EnterpriseManagementV1) NewDeleteAccountOptions(accountID string) *DeleteAccountOptions {
	return &DeleteAccountOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteAccountOptions) SetAccountID(accountID string) *DeleteAccountOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccountOptions) SetHeaders(param map[string]string) *DeleteAccountOptions {
	options.Headers = param
	return options
}

// Enterprise : An enterprise resource.
type Enterprise struct {
	// The URL of the enterprise.
	URL *string `json:"url,omitempty"`

	// The enterprise ID.
	ID *string `json:"id,omitempty"`

	// The enterprise account ID.
	EnterpriseAccountID *string `json:"enterprise_account_id,omitempty"`

	// The Cloud Resource Name (CRN) of the enterprise.
	CRN *string `json:"crn,omitempty"`

	// The name of the enterprise.
	Name *string `json:"name,omitempty"`

	// The domain of the enterprise.
	Domain *string `json:"domain,omitempty"`

	// The state of the enterprise.
	State *string `json:"state,omitempty"`

	// The IAM ID of the primary contact of the enterprise, such as `IBMid-0123ABC`.
	PrimaryContactIamID *string `json:"primary_contact_iam_id,omitempty"`

	// The email of the primary contact of the enterprise.
	PrimaryContactEmail *string `json:"primary_contact_email,omitempty"`

	// The ID of the account that is used to create the enterprise.
	SourceAccountID *string `json:"source_account_id,omitempty"`

	// The time stamp at which the enterprise was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The IAM ID of the user or service that created the enterprise.
	CreatedBy *string `json:"created_by,omitempty"`

	// The time stamp at which the enterprise was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The IAM ID of the user or service that updated the enterprise.
	UpdatedBy *string `json:"updated_by,omitempty"`
}

// UnmarshalEnterprise unmarshals an instance of Enterprise from the specified map of raw messages.
func UnmarshalEnterprise(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Enterprise)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enterprise_account_id", &obj.EnterpriseAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "enterprise_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "domain", &obj.Domain)
	if err != nil {
		err = core.SDKErrorf(err, "", "domain-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_contact_iam_id", &obj.PrimaryContactIamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact_iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_contact_email", &obj.PrimaryContactEmail)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact_email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_account_id", &obj.SourceAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccountGroupOptions : The GetAccountGroup options.
type GetAccountGroupOptions struct {
	// The ID of the account group to retrieve.
	AccountGroupID *string `json:"account_group_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccountGroupOptions : Instantiate GetAccountGroupOptions
func (*EnterpriseManagementV1) NewGetAccountGroupOptions(accountGroupID string) *GetAccountGroupOptions {
	return &GetAccountGroupOptions{
		AccountGroupID: core.StringPtr(accountGroupID),
	}
}

// SetAccountGroupID : Allow user to set AccountGroupID
func (_options *GetAccountGroupOptions) SetAccountGroupID(accountGroupID string) *GetAccountGroupOptions {
	_options.AccountGroupID = core.StringPtr(accountGroupID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountGroupOptions) SetHeaders(param map[string]string) *GetAccountGroupOptions {
	options.Headers = param
	return options
}

// GetAccountOptions : The GetAccount options.
type GetAccountOptions struct {
	// The ID of the target account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccountOptions : Instantiate GetAccountOptions
func (*EnterpriseManagementV1) NewGetAccountOptions(accountID string) *GetAccountOptions {
	return &GetAccountOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountOptions) SetAccountID(accountID string) *GetAccountOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountOptions) SetHeaders(param map[string]string) *GetAccountOptions {
	options.Headers = param
	return options
}

// GetEnterpriseOptions : The GetEnterprise options.
type GetEnterpriseOptions struct {
	// The ID of the enterprise to retrieve.
	EnterpriseID *string `json:"enterprise_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEnterpriseOptions : Instantiate GetEnterpriseOptions
func (*EnterpriseManagementV1) NewGetEnterpriseOptions(enterpriseID string) *GetEnterpriseOptions {
	return &GetEnterpriseOptions{
		EnterpriseID: core.StringPtr(enterpriseID),
	}
}

// SetEnterpriseID : Allow user to set EnterpriseID
func (_options *GetEnterpriseOptions) SetEnterpriseID(enterpriseID string) *GetEnterpriseOptions {
	_options.EnterpriseID = core.StringPtr(enterpriseID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetEnterpriseOptions) SetHeaders(param map[string]string) *GetEnterpriseOptions {
	options.Headers = param
	return options
}

// ImportAccountToEnterpriseOptions : The ImportAccountToEnterprise options.
type ImportAccountToEnterpriseOptions struct {
	// The ID of the enterprise to import the stand-alone account into.
	EnterpriseID *string `json:"enterprise_id" validate:"required,ne="`

	// The ID of the existing stand-alone account to be imported.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The CRN of the expected parent of the imported account. The parent is the enterprise or account group that the
	// account is added to.
	Parent *string `json:"parent,omitempty"`

	// The ID of the [billing unit](/apidocs/enterprise-apis/billing-unit) to use for billing this account in the
	// enterprise.
	BillingUnitID *string `json:"billing_unit_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewImportAccountToEnterpriseOptions : Instantiate ImportAccountToEnterpriseOptions
func (*EnterpriseManagementV1) NewImportAccountToEnterpriseOptions(enterpriseID string, accountID string) *ImportAccountToEnterpriseOptions {
	return &ImportAccountToEnterpriseOptions{
		EnterpriseID: core.StringPtr(enterpriseID),
		AccountID:    core.StringPtr(accountID),
	}
}

// SetEnterpriseID : Allow user to set EnterpriseID
func (_options *ImportAccountToEnterpriseOptions) SetEnterpriseID(enterpriseID string) *ImportAccountToEnterpriseOptions {
	_options.EnterpriseID = core.StringPtr(enterpriseID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ImportAccountToEnterpriseOptions) SetAccountID(accountID string) *ImportAccountToEnterpriseOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetParent : Allow user to set Parent
func (_options *ImportAccountToEnterpriseOptions) SetParent(parent string) *ImportAccountToEnterpriseOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetBillingUnitID : Allow user to set BillingUnitID
func (_options *ImportAccountToEnterpriseOptions) SetBillingUnitID(billingUnitID string) *ImportAccountToEnterpriseOptions {
	_options.BillingUnitID = core.StringPtr(billingUnitID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ImportAccountToEnterpriseOptions) SetHeaders(param map[string]string) *ImportAccountToEnterpriseOptions {
	options.Headers = param
	return options
}

// ListAccountGroupsOptions : The ListAccountGroups options.
type ListAccountGroupsOptions struct {
	// Get account groups that are either immediate children or are a part of the hierarchy for a given enterprise ID.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// Get account groups that are either immediate children or are a part of the hierarchy for a given account group ID.
	ParentAccountGroupID *string `json:"parent_account_group_id,omitempty"`

	// The first item to be returned in the page of results. This value can be obtained from the next_url property from the
	// previous call of the operation. If not specified, then the first page of results is returned.
	NextDocid *string `json:"next_docid,omitempty"`

	// Get account groups that are either immediate children or are a part of the hierarchy for a given parent CRN.
	Parent *string `json:"parent,omitempty"`

	// Return results up to this limit. Valid values are between `0` and `100`.
	Limit *int64 `json:"limit,omitempty"`

	// Include the deleted account groups from an enterprise when used in conjunction with other query parameters.
	IncludeDeleted *bool `json:"include_deleted,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAccountGroupsOptions : Instantiate ListAccountGroupsOptions
func (*EnterpriseManagementV1) NewListAccountGroupsOptions() *ListAccountGroupsOptions {
	return &ListAccountGroupsOptions{}
}

// SetEnterpriseID : Allow user to set EnterpriseID
func (_options *ListAccountGroupsOptions) SetEnterpriseID(enterpriseID string) *ListAccountGroupsOptions {
	_options.EnterpriseID = core.StringPtr(enterpriseID)
	return _options
}

// SetParentAccountGroupID : Allow user to set ParentAccountGroupID
func (_options *ListAccountGroupsOptions) SetParentAccountGroupID(parentAccountGroupID string) *ListAccountGroupsOptions {
	_options.ParentAccountGroupID = core.StringPtr(parentAccountGroupID)
	return _options
}

// SetNextDocid : Allow user to set NextDocid
func (_options *ListAccountGroupsOptions) SetNextDocid(nextDocid string) *ListAccountGroupsOptions {
	_options.NextDocid = core.StringPtr(nextDocid)
	return _options
}

// SetParent : Allow user to set Parent
func (_options *ListAccountGroupsOptions) SetParent(parent string) *ListAccountGroupsOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAccountGroupsOptions) SetLimit(limit int64) *ListAccountGroupsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetIncludeDeleted : Allow user to set IncludeDeleted
func (_options *ListAccountGroupsOptions) SetIncludeDeleted(includeDeleted bool) *ListAccountGroupsOptions {
	_options.IncludeDeleted = core.BoolPtr(includeDeleted)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccountGroupsOptions) SetHeaders(param map[string]string) *ListAccountGroupsOptions {
	options.Headers = param
	return options
}

// ListAccountGroupsResponse : The list_account_groups operation response.
type ListAccountGroupsResponse struct {
	// The number of enterprises returned from calling list account groups.
	RowsCount *int64 `json:"rows_count,omitempty"`

	// A string that represents the link to the next page of results.
	NextURL *string `json:"next_url,omitempty"`

	// A list of account groups.
	Resources []AccountGroup `json:"resources,omitempty"`
}

// UnmarshalListAccountGroupsResponse unmarshals an instance of ListAccountGroupsResponse from the specified map of raw messages.
func UnmarshalListAccountGroupsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAccountGroupsResponse)
	err = core.UnmarshalPrimitive(m, "rows_count", &obj.RowsCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "rows_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next_url", &obj.NextURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "next_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalAccountGroup)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListAccountGroupsResponse) GetNextNextDocid() (*string, error) {
	if core.IsNil(resp.NextURL) {
		return nil, nil
	}
	next_docid, err := core.GetQueryParam(resp.NextURL, "next_docid")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if next_docid == nil {
		return nil, nil
	}
	return next_docid, nil
}

// ListAccountsOptions : The ListAccounts options.
type ListAccountsOptions struct {
	// Get accounts that are either immediate children or are a part of the hierarchy for a given enterprise ID.
	EnterpriseID *string `json:"enterprise_id,omitempty"`

	// Get accounts that are either immediate children or are a part of the hierarchy for a given account group ID.
	AccountGroupID *string `json:"account_group_id,omitempty"`

	// The first item to be returned in the page of results. This value can be obtained from the next_url property from the
	// previous call of the operation. If not specified, then the first page of results is returned.
	NextDocid *string `json:"next_docid,omitempty"`

	// Get accounts that are either immediate children or are a part of the hierarchy for a given parent CRN.
	Parent *string `json:"parent,omitempty"`

	// Return results up to this limit. Valid values are between `0` and `100`.
	Limit *int64 `json:"limit,omitempty"`

	// Include the deleted accounts from an enterprise when used in conjunction with enterprise_id.
	IncludeDeleted *bool `json:"include_deleted,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAccountsOptions : Instantiate ListAccountsOptions
func (*EnterpriseManagementV1) NewListAccountsOptions() *ListAccountsOptions {
	return &ListAccountsOptions{}
}

// SetEnterpriseID : Allow user to set EnterpriseID
func (_options *ListAccountsOptions) SetEnterpriseID(enterpriseID string) *ListAccountsOptions {
	_options.EnterpriseID = core.StringPtr(enterpriseID)
	return _options
}

// SetAccountGroupID : Allow user to set AccountGroupID
func (_options *ListAccountsOptions) SetAccountGroupID(accountGroupID string) *ListAccountsOptions {
	_options.AccountGroupID = core.StringPtr(accountGroupID)
	return _options
}

// SetNextDocid : Allow user to set NextDocid
func (_options *ListAccountsOptions) SetNextDocid(nextDocid string) *ListAccountsOptions {
	_options.NextDocid = core.StringPtr(nextDocid)
	return _options
}

// SetParent : Allow user to set Parent
func (_options *ListAccountsOptions) SetParent(parent string) *ListAccountsOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAccountsOptions) SetLimit(limit int64) *ListAccountsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetIncludeDeleted : Allow user to set IncludeDeleted
func (_options *ListAccountsOptions) SetIncludeDeleted(includeDeleted bool) *ListAccountsOptions {
	_options.IncludeDeleted = core.BoolPtr(includeDeleted)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccountsOptions) SetHeaders(param map[string]string) *ListAccountsOptions {
	options.Headers = param
	return options
}

// ListAccountsResponse : The list_accounts operation response.
type ListAccountsResponse struct {
	// The number of enterprises returned from calling list accounts.
	RowsCount *int64 `json:"rows_count,omitempty"`

	// A string that represents the link to the next page of results.
	NextURL *string `json:"next_url,omitempty"`

	// A list of accounts.
	Resources []Account `json:"resources,omitempty"`
}

// UnmarshalListAccountsResponse unmarshals an instance of ListAccountsResponse from the specified map of raw messages.
func UnmarshalListAccountsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAccountsResponse)
	err = core.UnmarshalPrimitive(m, "rows_count", &obj.RowsCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "rows_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next_url", &obj.NextURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "next_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListAccountsResponse) GetNextNextDocid() (*string, error) {
	if core.IsNil(resp.NextURL) {
		return nil, nil
	}
	next_docid, err := core.GetQueryParam(resp.NextURL, "next_docid")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if next_docid == nil {
		return nil, nil
	}
	return next_docid, nil
}

// ListEnterprisesOptions : The ListEnterprises options.
type ListEnterprisesOptions struct {
	// Get enterprises for a given enterprise account ID.
	EnterpriseAccountID *string `json:"enterprise_account_id,omitempty"`

	// Get enterprises for a given account group ID.
	AccountGroupID *string `json:"account_group_id,omitempty"`

	// Get enterprises for a given account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The first item to be returned in the page of results. This value can be obtained from the next_url property from the
	// previous call of the operation. If not specified, then the first page of results is returned.
	NextDocid *string `json:"next_docid,omitempty"`

	// Return results up to this limit. Valid values are between `0` and `100`.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListEnterprisesOptions : Instantiate ListEnterprisesOptions
func (*EnterpriseManagementV1) NewListEnterprisesOptions() *ListEnterprisesOptions {
	return &ListEnterprisesOptions{}
}

// SetEnterpriseAccountID : Allow user to set EnterpriseAccountID
func (_options *ListEnterprisesOptions) SetEnterpriseAccountID(enterpriseAccountID string) *ListEnterprisesOptions {
	_options.EnterpriseAccountID = core.StringPtr(enterpriseAccountID)
	return _options
}

// SetAccountGroupID : Allow user to set AccountGroupID
func (_options *ListEnterprisesOptions) SetAccountGroupID(accountGroupID string) *ListEnterprisesOptions {
	_options.AccountGroupID = core.StringPtr(accountGroupID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListEnterprisesOptions) SetAccountID(accountID string) *ListEnterprisesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetNextDocid : Allow user to set NextDocid
func (_options *ListEnterprisesOptions) SetNextDocid(nextDocid string) *ListEnterprisesOptions {
	_options.NextDocid = core.StringPtr(nextDocid)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListEnterprisesOptions) SetLimit(limit int64) *ListEnterprisesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListEnterprisesOptions) SetHeaders(param map[string]string) *ListEnterprisesOptions {
	options.Headers = param
	return options
}

// ListEnterprisesResponse : The response from calling list enterprises.
type ListEnterprisesResponse struct {
	// The number of enterprises returned from calling list enterprise.
	RowsCount *int64 `json:"rows_count,omitempty"`

	// A string that represents the link to the next page of results.
	NextURL *string `json:"next_url,omitempty"`

	// A list of enterprise objects.
	Resources []Enterprise `json:"resources,omitempty"`
}

// UnmarshalListEnterprisesResponse unmarshals an instance of ListEnterprisesResponse from the specified map of raw messages.
func UnmarshalListEnterprisesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListEnterprisesResponse)
	err = core.UnmarshalPrimitive(m, "rows_count", &obj.RowsCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "rows_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next_url", &obj.NextURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "next_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalEnterprise)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ListEnterprisesResponse) GetNextNextDocid() (*string, error) {
	if core.IsNil(resp.NextURL) {
		return nil, nil
	}
	next_docid, err := core.GetQueryParam(resp.NextURL, "next_docid")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if next_docid == nil {
		return nil, nil
	}
	return next_docid, nil
}

// UpdateAccountGroupOptions : The UpdateAccountGroup options.
type UpdateAccountGroupOptions struct {
	// The ID of the account group to retrieve.
	AccountGroupID *string `json:"account_group_id" validate:"required,ne="`

	// The new name of the account group. This field must have 3 - 60 characters.
	Name *string `json:"name,omitempty"`

	// The IAM ID of the user to be the new primary contact for the account group.
	PrimaryContactIamID *string `json:"primary_contact_iam_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccountGroupOptions : Instantiate UpdateAccountGroupOptions
func (*EnterpriseManagementV1) NewUpdateAccountGroupOptions(accountGroupID string) *UpdateAccountGroupOptions {
	return &UpdateAccountGroupOptions{
		AccountGroupID: core.StringPtr(accountGroupID),
	}
}

// SetAccountGroupID : Allow user to set AccountGroupID
func (_options *UpdateAccountGroupOptions) SetAccountGroupID(accountGroupID string) *UpdateAccountGroupOptions {
	_options.AccountGroupID = core.StringPtr(accountGroupID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAccountGroupOptions) SetName(name string) *UpdateAccountGroupOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPrimaryContactIamID : Allow user to set PrimaryContactIamID
func (_options *UpdateAccountGroupOptions) SetPrimaryContactIamID(primaryContactIamID string) *UpdateAccountGroupOptions {
	_options.PrimaryContactIamID = core.StringPtr(primaryContactIamID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountGroupOptions) SetHeaders(param map[string]string) *UpdateAccountGroupOptions {
	options.Headers = param
	return options
}

// UpdateAccountOptions : The UpdateAccount options.
type UpdateAccountOptions struct {
	// The ID of the target account.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The CRN of the new parent within the enterprise.
	Parent *string `json:"parent" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccountOptions : Instantiate UpdateAccountOptions
func (*EnterpriseManagementV1) NewUpdateAccountOptions(accountID string, parent string) *UpdateAccountOptions {
	return &UpdateAccountOptions{
		AccountID: core.StringPtr(accountID),
		Parent:    core.StringPtr(parent),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateAccountOptions) SetAccountID(accountID string) *UpdateAccountOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetParent : Allow user to set Parent
func (_options *UpdateAccountOptions) SetParent(parent string) *UpdateAccountOptions {
	_options.Parent = core.StringPtr(parent)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccountOptions) SetHeaders(param map[string]string) *UpdateAccountOptions {
	options.Headers = param
	return options
}

// UpdateEnterpriseOptions : The UpdateEnterprise options.
type UpdateEnterpriseOptions struct {
	// The ID of the enterprise to retrieve.
	EnterpriseID *string `json:"enterprise_id" validate:"required,ne="`

	// The new name of the enterprise. This field must have 3 - 60 characters.
	Name *string `json:"name,omitempty"`

	// The new domain of the enterprise. This field has a limit of 60 characters.
	Domain *string `json:"domain,omitempty"`

	// The IAM ID of the user to be the new primary contact for the enterprise.
	PrimaryContactIamID *string `json:"primary_contact_iam_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEnterpriseOptions : Instantiate UpdateEnterpriseOptions
func (*EnterpriseManagementV1) NewUpdateEnterpriseOptions(enterpriseID string) *UpdateEnterpriseOptions {
	return &UpdateEnterpriseOptions{
		EnterpriseID: core.StringPtr(enterpriseID),
	}
}

// SetEnterpriseID : Allow user to set EnterpriseID
func (_options *UpdateEnterpriseOptions) SetEnterpriseID(enterpriseID string) *UpdateEnterpriseOptions {
	_options.EnterpriseID = core.StringPtr(enterpriseID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateEnterpriseOptions) SetName(name string) *UpdateEnterpriseOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDomain : Allow user to set Domain
func (_options *UpdateEnterpriseOptions) SetDomain(domain string) *UpdateEnterpriseOptions {
	_options.Domain = core.StringPtr(domain)
	return _options
}

// SetPrimaryContactIamID : Allow user to set PrimaryContactIamID
func (_options *UpdateEnterpriseOptions) SetPrimaryContactIamID(primaryContactIamID string) *UpdateEnterpriseOptions {
	_options.PrimaryContactIamID = core.StringPtr(primaryContactIamID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEnterpriseOptions) SetHeaders(param map[string]string) *UpdateEnterpriseOptions {
	options.Headers = param
	return options
}

// EnterprisesPager can be used to simplify the use of the "ListEnterprises" method.
type EnterprisesPager struct {
	hasNext     bool
	options     *ListEnterprisesOptions
	client      *EnterpriseManagementV1
	pageContext struct {
		next *string
	}
}

// NewEnterprisesPager returns a new EnterprisesPager instance.
func (enterpriseManagement *EnterpriseManagementV1) NewEnterprisesPager(options *ListEnterprisesOptions) (pager *EnterprisesPager, err error) {
	if options.NextDocid != nil && *options.NextDocid != "" {
		err = core.SDKErrorf(nil, "the 'options.NextDocid' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListEnterprisesOptions = *options
	pager = &EnterprisesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  enterpriseManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *EnterprisesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *EnterprisesPager) GetNextWithContext(ctx context.Context) (page []Enterprise, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.NextDocid = pager.pageContext.next

	result, _, err := pager.client.ListEnterprisesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.NextURL != nil {
		var next_docid *string
		next_docid, err = core.GetQueryParam(result.NextURL, "next_docid")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'next_docid' query parameter from URL '%s': %s", *result.NextURL, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = next_docid
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *EnterprisesPager) GetAllWithContext(ctx context.Context) (allItems []Enterprise, err error) {
	for pager.HasNext() {
		var nextPage []Enterprise
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *EnterprisesPager) GetNext() (page []Enterprise, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *EnterprisesPager) GetAll() (allItems []Enterprise, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AccountsPager can be used to simplify the use of the "ListAccounts" method.
type AccountsPager struct {
	hasNext     bool
	options     *ListAccountsOptions
	client      *EnterpriseManagementV1
	pageContext struct {
		next *string
	}
}

// NewAccountsPager returns a new AccountsPager instance.
func (enterpriseManagement *EnterpriseManagementV1) NewAccountsPager(options *ListAccountsOptions) (pager *AccountsPager, err error) {
	if options.NextDocid != nil && *options.NextDocid != "" {
		err = core.SDKErrorf(nil, "the 'options.NextDocid' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAccountsOptions = *options
	pager = &AccountsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  enterpriseManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AccountsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AccountsPager) GetNextWithContext(ctx context.Context) (page []Account, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.NextDocid = pager.pageContext.next

	result, _, err := pager.client.ListAccountsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.NextURL != nil {
		var next_docid *string
		next_docid, err = core.GetQueryParam(result.NextURL, "next_docid")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'next_docid' query parameter from URL '%s': %s", *result.NextURL, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = next_docid
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AccountsPager) GetAllWithContext(ctx context.Context) (allItems []Account, err error) {
	for pager.HasNext() {
		var nextPage []Account
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AccountsPager) GetNext() (page []Account, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AccountsPager) GetAll() (allItems []Account, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AccountGroupsPager can be used to simplify the use of the "ListAccountGroups" method.
type AccountGroupsPager struct {
	hasNext     bool
	options     *ListAccountGroupsOptions
	client      *EnterpriseManagementV1
	pageContext struct {
		next *string
	}
}

// NewAccountGroupsPager returns a new AccountGroupsPager instance.
func (enterpriseManagement *EnterpriseManagementV1) NewAccountGroupsPager(options *ListAccountGroupsOptions) (pager *AccountGroupsPager, err error) {
	if options.NextDocid != nil && *options.NextDocid != "" {
		err = core.SDKErrorf(nil, "the 'options.NextDocid' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAccountGroupsOptions = *options
	pager = &AccountGroupsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  enterpriseManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AccountGroupsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AccountGroupsPager) GetNextWithContext(ctx context.Context) (page []AccountGroup, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.NextDocid = pager.pageContext.next

	result, _, err := pager.client.ListAccountGroupsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.NextURL != nil {
		var next_docid *string
		next_docid, err = core.GetQueryParam(result.NextURL, "next_docid")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'next_docid' query parameter from URL '%s': %s", *result.NextURL, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = next_docid
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AccountGroupsPager) GetAllWithContext(ctx context.Context) (allItems []AccountGroup, err error) {
	for pager.HasNext() {
		var nextPage []AccountGroup
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *AccountGroupsPager) GetNext() (page []AccountGroup, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AccountGroupsPager) GetAll() (allItems []AccountGroup, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
