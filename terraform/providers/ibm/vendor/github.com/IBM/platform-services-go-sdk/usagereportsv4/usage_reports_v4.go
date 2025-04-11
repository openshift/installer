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
 * IBM OpenAPI SDK Code Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

// Package usagereportsv4 : Operations and models for the UsageReportsV4 service
package usagereportsv4

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

// UsageReportsV4 : Usage reports for IBM Cloud accounts
//
// API Version: 4.0.6
type UsageReportsV4 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://billing.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "usage_reports"

// UsageReportsV4Options : Service options
type UsageReportsV4Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewUsageReportsV4UsingExternalConfig : constructs an instance of UsageReportsV4 with passed in options and external configuration.
func NewUsageReportsV4UsingExternalConfig(options *UsageReportsV4Options) (usageReports *UsageReportsV4, err error) {
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

	usageReports, err = NewUsageReportsV4(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = usageReports.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = usageReports.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewUsageReportsV4 : constructs an instance of UsageReportsV4 with passed in options.
func NewUsageReportsV4(options *UsageReportsV4Options) (service *UsageReportsV4, err error) {
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

	service = &UsageReportsV4{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "usageReports" suitable for processing requests.
func (usageReports *UsageReportsV4) Clone() *UsageReportsV4 {
	if core.IsNil(usageReports) {
		return nil
	}
	clone := *usageReports
	clone.Service = usageReports.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (usageReports *UsageReportsV4) SetServiceURL(url string) error {
	err := usageReports.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (usageReports *UsageReportsV4) GetServiceURL() string {
	return usageReports.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (usageReports *UsageReportsV4) SetDefaultHeaders(headers http.Header) {
	usageReports.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (usageReports *UsageReportsV4) SetEnableGzipCompression(enableGzip bool) {
	usageReports.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (usageReports *UsageReportsV4) GetEnableGzipCompression() bool {
	return usageReports.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (usageReports *UsageReportsV4) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	usageReports.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (usageReports *UsageReportsV4) DisableRetries() {
	usageReports.Service.DisableRetries()
}

// GetAccountSummary : Get account summary
// Returns the summary for the account for a given month. Account billing managers are authorized to access this report.
func (usageReports *UsageReportsV4) GetAccountSummary(getAccountSummaryOptions *GetAccountSummaryOptions) (result *AccountSummary, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetAccountSummaryWithContext(context.Background(), getAccountSummaryOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountSummaryWithContext is an alternate form of the GetAccountSummary method which supports a Context parameter
func (usageReports *UsageReportsV4) GetAccountSummaryWithContext(ctx context.Context, getAccountSummaryOptions *GetAccountSummaryOptions) (result *AccountSummary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSummaryOptions, "getAccountSummaryOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountSummaryOptions, "getAccountSummaryOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountSummaryOptions.AccountID,
		"billingmonth": *getAccountSummaryOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/summary/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountSummaryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetAccountSummary")
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
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_summary", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSummary)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAccountUsage : Get account usage
// Usage for all the resources and plans in an account for a given month. Account billing managers are authorized to
// access this report.
func (usageReports *UsageReportsV4) GetAccountUsage(getAccountUsageOptions *GetAccountUsageOptions) (result *AccountUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetAccountUsageWithContext(context.Background(), getAccountUsageOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountUsageWithContext is an alternate form of the GetAccountUsage method which supports a Context parameter
func (usageReports *UsageReportsV4) GetAccountUsageWithContext(ctx context.Context, getAccountUsageOptions *GetAccountUsageOptions) (result *AccountUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountUsageOptions, "getAccountUsageOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountUsageOptions, "getAccountUsageOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountUsageOptions.AccountID,
		"billingmonth": *getAccountUsageOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountUsageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetAccountUsage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getAccountUsageOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getAccountUsageOptions.AcceptLanguage))
	}

	if getAccountUsageOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getAccountUsageOptions.Names))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_usage", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceGroupUsage : Get resource group usage
// Usage for all the resources and plans in a resource group in a given month. Account billing managers or resource
// group billing managers are authorized to access this report.
func (usageReports *UsageReportsV4) GetResourceGroupUsage(getResourceGroupUsageOptions *GetResourceGroupUsageOptions) (result *ResourceGroupUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetResourceGroupUsageWithContext(context.Background(), getResourceGroupUsageOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceGroupUsageWithContext is an alternate form of the GetResourceGroupUsage method which supports a Context parameter
func (usageReports *UsageReportsV4) GetResourceGroupUsageWithContext(ctx context.Context, getResourceGroupUsageOptions *GetResourceGroupUsageOptions) (result *ResourceGroupUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceGroupUsageOptions, "getResourceGroupUsageOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceGroupUsageOptions, "getResourceGroupUsageOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getResourceGroupUsageOptions.AccountID,
		"resource_group_id": *getResourceGroupUsageOptions.ResourceGroupID,
		"billingmonth": *getResourceGroupUsageOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/resource_groups/{resource_group_id}/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceGroupUsageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetResourceGroupUsage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getResourceGroupUsageOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getResourceGroupUsageOptions.AcceptLanguage))
	}

	if getResourceGroupUsageOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getResourceGroupUsageOptions.Names))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_group_usage", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroupUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceUsageAccount : Get resource instance usage in an account
// Query for resource instance usage in an account. Filter the results with query parameters. Account billing
// administrator is authorized to access this report.
func (usageReports *UsageReportsV4) GetResourceUsageAccount(getResourceUsageAccountOptions *GetResourceUsageAccountOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetResourceUsageAccountWithContext(context.Background(), getResourceUsageAccountOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceUsageAccountWithContext is an alternate form of the GetResourceUsageAccount method which supports a Context parameter
func (usageReports *UsageReportsV4) GetResourceUsageAccountWithContext(ctx context.Context, getResourceUsageAccountOptions *GetResourceUsageAccountOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceUsageAccountOptions, "getResourceUsageAccountOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceUsageAccountOptions, "getResourceUsageAccountOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getResourceUsageAccountOptions.AccountID,
		"billingmonth": *getResourceUsageAccountOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/resource_instances/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceUsageAccountOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetResourceUsageAccount")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getResourceUsageAccountOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getResourceUsageAccountOptions.AcceptLanguage))
	}

	if getResourceUsageAccountOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getResourceUsageAccountOptions.Names))
	}
	if getResourceUsageAccountOptions.Tags != nil {
		builder.AddQuery("_tags", fmt.Sprint(*getResourceUsageAccountOptions.Tags))
	}
	if getResourceUsageAccountOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getResourceUsageAccountOptions.Limit))
	}
	if getResourceUsageAccountOptions.Start != nil {
		builder.AddQuery("_start", fmt.Sprint(*getResourceUsageAccountOptions.Start))
	}
	if getResourceUsageAccountOptions.ResourceGroupID != nil {
		builder.AddQuery("resource_group_id", fmt.Sprint(*getResourceUsageAccountOptions.ResourceGroupID))
	}
	if getResourceUsageAccountOptions.OrganizationID != nil {
		builder.AddQuery("organization_id", fmt.Sprint(*getResourceUsageAccountOptions.OrganizationID))
	}
	if getResourceUsageAccountOptions.ResourceInstanceID != nil {
		builder.AddQuery("resource_instance_id", fmt.Sprint(*getResourceUsageAccountOptions.ResourceInstanceID))
	}
	if getResourceUsageAccountOptions.ResourceID != nil {
		builder.AddQuery("resource_id", fmt.Sprint(*getResourceUsageAccountOptions.ResourceID))
	}
	if getResourceUsageAccountOptions.PlanID != nil {
		builder.AddQuery("plan_id", fmt.Sprint(*getResourceUsageAccountOptions.PlanID))
	}
	if getResourceUsageAccountOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*getResourceUsageAccountOptions.Region))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_usage_account", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInstancesUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceUsageResourceGroup : Get resource instance usage in a resource group
// Query for resource instance usage in a resource group. Filter the results with query parameters. Account billing
// administrator and resource group billing administrators are authorized to access this report.
func (usageReports *UsageReportsV4) GetResourceUsageResourceGroup(getResourceUsageResourceGroupOptions *GetResourceUsageResourceGroupOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetResourceUsageResourceGroupWithContext(context.Background(), getResourceUsageResourceGroupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceUsageResourceGroupWithContext is an alternate form of the GetResourceUsageResourceGroup method which supports a Context parameter
func (usageReports *UsageReportsV4) GetResourceUsageResourceGroupWithContext(ctx context.Context, getResourceUsageResourceGroupOptions *GetResourceUsageResourceGroupOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceUsageResourceGroupOptions, "getResourceUsageResourceGroupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceUsageResourceGroupOptions, "getResourceUsageResourceGroupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getResourceUsageResourceGroupOptions.AccountID,
		"resource_group_id": *getResourceUsageResourceGroupOptions.ResourceGroupID,
		"billingmonth": *getResourceUsageResourceGroupOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/resource_groups/{resource_group_id}/resource_instances/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceUsageResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetResourceUsageResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getResourceUsageResourceGroupOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getResourceUsageResourceGroupOptions.AcceptLanguage))
	}

	if getResourceUsageResourceGroupOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getResourceUsageResourceGroupOptions.Names))
	}
	if getResourceUsageResourceGroupOptions.Tags != nil {
		builder.AddQuery("_tags", fmt.Sprint(*getResourceUsageResourceGroupOptions.Tags))
	}
	if getResourceUsageResourceGroupOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getResourceUsageResourceGroupOptions.Limit))
	}
	if getResourceUsageResourceGroupOptions.Start != nil {
		builder.AddQuery("_start", fmt.Sprint(*getResourceUsageResourceGroupOptions.Start))
	}
	if getResourceUsageResourceGroupOptions.ResourceInstanceID != nil {
		builder.AddQuery("resource_instance_id", fmt.Sprint(*getResourceUsageResourceGroupOptions.ResourceInstanceID))
	}
	if getResourceUsageResourceGroupOptions.ResourceID != nil {
		builder.AddQuery("resource_id", fmt.Sprint(*getResourceUsageResourceGroupOptions.ResourceID))
	}
	if getResourceUsageResourceGroupOptions.PlanID != nil {
		builder.AddQuery("plan_id", fmt.Sprint(*getResourceUsageResourceGroupOptions.PlanID))
	}
	if getResourceUsageResourceGroupOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*getResourceUsageResourceGroupOptions.Region))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_usage_resource_group", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInstancesUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceUsageOrg : Get resource instance usage in an organization
// Query for resource instance usage in an organization. Filter the results with query parameters. Account billing
// administrator and organization billing administrators are authorized to access this report.
func (usageReports *UsageReportsV4) GetResourceUsageOrg(getResourceUsageOrgOptions *GetResourceUsageOrgOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetResourceUsageOrgWithContext(context.Background(), getResourceUsageOrgOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceUsageOrgWithContext is an alternate form of the GetResourceUsageOrg method which supports a Context parameter
func (usageReports *UsageReportsV4) GetResourceUsageOrgWithContext(ctx context.Context, getResourceUsageOrgOptions *GetResourceUsageOrgOptions) (result *InstancesUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceUsageOrgOptions, "getResourceUsageOrgOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceUsageOrgOptions, "getResourceUsageOrgOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getResourceUsageOrgOptions.AccountID,
		"organization_id": *getResourceUsageOrgOptions.OrganizationID,
		"billingmonth": *getResourceUsageOrgOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/organizations/{organization_id}/resource_instances/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceUsageOrgOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetResourceUsageOrg")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getResourceUsageOrgOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getResourceUsageOrgOptions.AcceptLanguage))
	}

	if getResourceUsageOrgOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getResourceUsageOrgOptions.Names))
	}
	if getResourceUsageOrgOptions.Tags != nil {
		builder.AddQuery("_tags", fmt.Sprint(*getResourceUsageOrgOptions.Tags))
	}
	if getResourceUsageOrgOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getResourceUsageOrgOptions.Limit))
	}
	if getResourceUsageOrgOptions.Start != nil {
		builder.AddQuery("_start", fmt.Sprint(*getResourceUsageOrgOptions.Start))
	}
	if getResourceUsageOrgOptions.ResourceInstanceID != nil {
		builder.AddQuery("resource_instance_id", fmt.Sprint(*getResourceUsageOrgOptions.ResourceInstanceID))
	}
	if getResourceUsageOrgOptions.ResourceID != nil {
		builder.AddQuery("resource_id", fmt.Sprint(*getResourceUsageOrgOptions.ResourceID))
	}
	if getResourceUsageOrgOptions.PlanID != nil {
		builder.AddQuery("plan_id", fmt.Sprint(*getResourceUsageOrgOptions.PlanID))
	}
	if getResourceUsageOrgOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*getResourceUsageOrgOptions.Region))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_usage_org", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInstancesUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetOrgUsage : Get organization usage
// Usage for all the resources and plans in an organization in a given month. Account billing managers or organization
// billing managers are authorized to access this report.
func (usageReports *UsageReportsV4) GetOrgUsage(getOrgUsageOptions *GetOrgUsageOptions) (result *OrgUsage, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetOrgUsageWithContext(context.Background(), getOrgUsageOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetOrgUsageWithContext is an alternate form of the GetOrgUsage method which supports a Context parameter
func (usageReports *UsageReportsV4) GetOrgUsageWithContext(ctx context.Context, getOrgUsageOptions *GetOrgUsageOptions) (result *OrgUsage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOrgUsageOptions, "getOrgUsageOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getOrgUsageOptions, "getOrgUsageOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getOrgUsageOptions.AccountID,
		"organization_id": *getOrgUsageOptions.OrganizationID,
		"billingmonth": *getOrgUsageOptions.Billingmonth,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v4/accounts/{account_id}/organizations/{organization_id}/usage/{billingmonth}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getOrgUsageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetOrgUsage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getOrgUsageOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getOrgUsageOptions.AcceptLanguage))
	}

	if getOrgUsageOptions.Names != nil {
		builder.AddQuery("_names", fmt.Sprint(*getOrgUsageOptions.Names))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_org_usage", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOrgUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateReportsSnapshotConfig : Setup the snapshot configuration
// Snapshots of the billing reports would be taken on a periodic interval and stored based on the configuration setup by
// the customer for the given Account Id.
func (usageReports *UsageReportsV4) CreateReportsSnapshotConfig(createReportsSnapshotConfigOptions *CreateReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.CreateReportsSnapshotConfigWithContext(context.Background(), createReportsSnapshotConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateReportsSnapshotConfigWithContext is an alternate form of the CreateReportsSnapshotConfig method which supports a Context parameter
func (usageReports *UsageReportsV4) CreateReportsSnapshotConfigWithContext(ctx context.Context, createReportsSnapshotConfigOptions *CreateReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createReportsSnapshotConfigOptions, "createReportsSnapshotConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createReportsSnapshotConfigOptions, "createReportsSnapshotConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshot-config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createReportsSnapshotConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "CreateReportsSnapshotConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createReportsSnapshotConfigOptions.AccountID != nil {
		body["account_id"] = createReportsSnapshotConfigOptions.AccountID
	}
	if createReportsSnapshotConfigOptions.Interval != nil {
		body["interval"] = createReportsSnapshotConfigOptions.Interval
	}
	if createReportsSnapshotConfigOptions.CosBucket != nil {
		body["cos_bucket"] = createReportsSnapshotConfigOptions.CosBucket
	}
	if createReportsSnapshotConfigOptions.CosLocation != nil {
		body["cos_location"] = createReportsSnapshotConfigOptions.CosLocation
	}
	if createReportsSnapshotConfigOptions.CosReportsFolder != nil {
		body["cos_reports_folder"] = createReportsSnapshotConfigOptions.CosReportsFolder
	}
	if createReportsSnapshotConfigOptions.ReportTypes != nil {
		body["report_types"] = createReportsSnapshotConfigOptions.ReportTypes
	}
	if createReportsSnapshotConfigOptions.Versioning != nil {
		body["versioning"] = createReportsSnapshotConfigOptions.Versioning
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
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_reports_snapshot_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSnapshotConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetReportsSnapshotConfig : Fetch the snapshot configuration
// Returns the configuration of snapshot of the billing reports setup by the customer for the given Account Id.
func (usageReports *UsageReportsV4) GetReportsSnapshotConfig(getReportsSnapshotConfigOptions *GetReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetReportsSnapshotConfigWithContext(context.Background(), getReportsSnapshotConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetReportsSnapshotConfigWithContext is an alternate form of the GetReportsSnapshotConfig method which supports a Context parameter
func (usageReports *UsageReportsV4) GetReportsSnapshotConfigWithContext(ctx context.Context, getReportsSnapshotConfigOptions *GetReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportsSnapshotConfigOptions, "getReportsSnapshotConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getReportsSnapshotConfigOptions, "getReportsSnapshotConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshot-config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getReportsSnapshotConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetReportsSnapshotConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("account_id", fmt.Sprint(*getReportsSnapshotConfigOptions.AccountID))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_reports_snapshot_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSnapshotConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateReportsSnapshotConfig : Update the snapshot configuration
// Updates the configuration of snapshot of the billing reports setup by the customer for the given Account Id.
func (usageReports *UsageReportsV4) UpdateReportsSnapshotConfig(updateReportsSnapshotConfigOptions *UpdateReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.UpdateReportsSnapshotConfigWithContext(context.Background(), updateReportsSnapshotConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateReportsSnapshotConfigWithContext is an alternate form of the UpdateReportsSnapshotConfig method which supports a Context parameter
func (usageReports *UsageReportsV4) UpdateReportsSnapshotConfigWithContext(ctx context.Context, updateReportsSnapshotConfigOptions *UpdateReportsSnapshotConfigOptions) (result *SnapshotConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateReportsSnapshotConfigOptions, "updateReportsSnapshotConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateReportsSnapshotConfigOptions, "updateReportsSnapshotConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshot-config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateReportsSnapshotConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "UpdateReportsSnapshotConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateReportsSnapshotConfigOptions.AccountID != nil {
		body["account_id"] = updateReportsSnapshotConfigOptions.AccountID
	}
	if updateReportsSnapshotConfigOptions.Interval != nil {
		body["interval"] = updateReportsSnapshotConfigOptions.Interval
	}
	if updateReportsSnapshotConfigOptions.CosBucket != nil {
		body["cos_bucket"] = updateReportsSnapshotConfigOptions.CosBucket
	}
	if updateReportsSnapshotConfigOptions.CosLocation != nil {
		body["cos_location"] = updateReportsSnapshotConfigOptions.CosLocation
	}
	if updateReportsSnapshotConfigOptions.CosReportsFolder != nil {
		body["cos_reports_folder"] = updateReportsSnapshotConfigOptions.CosReportsFolder
	}
	if updateReportsSnapshotConfigOptions.ReportTypes != nil {
		body["report_types"] = updateReportsSnapshotConfigOptions.ReportTypes
	}
	if updateReportsSnapshotConfigOptions.Versioning != nil {
		body["versioning"] = updateReportsSnapshotConfigOptions.Versioning
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
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_reports_snapshot_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSnapshotConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteReportsSnapshotConfig : Delete the snapshot configuration
// Delete the configuration of snapshot of the billing reports setup by the customer for the given Account Id.
func (usageReports *UsageReportsV4) DeleteReportsSnapshotConfig(deleteReportsSnapshotConfigOptions *DeleteReportsSnapshotConfigOptions) (response *core.DetailedResponse, err error) {
	response, err = usageReports.DeleteReportsSnapshotConfigWithContext(context.Background(), deleteReportsSnapshotConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteReportsSnapshotConfigWithContext is an alternate form of the DeleteReportsSnapshotConfig method which supports a Context parameter
func (usageReports *UsageReportsV4) DeleteReportsSnapshotConfigWithContext(ctx context.Context, deleteReportsSnapshotConfigOptions *DeleteReportsSnapshotConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteReportsSnapshotConfigOptions, "deleteReportsSnapshotConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteReportsSnapshotConfigOptions, "deleteReportsSnapshotConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshot-config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteReportsSnapshotConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "DeleteReportsSnapshotConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddQuery("account_id", fmt.Sprint(*deleteReportsSnapshotConfigOptions.AccountID))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = usageReports.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_reports_snapshot_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ValidateReportsSnapshotConfig : Verify billing to COS authorization
// Verify billing service to COS bucket authorization for the given account_id. If COS bucket information is not
// provided, COS bucket information is retrieved from the configuration file.
func (usageReports *UsageReportsV4) ValidateReportsSnapshotConfig(validateReportsSnapshotConfigOptions *ValidateReportsSnapshotConfigOptions) (result *SnapshotConfigValidateResponse, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.ValidateReportsSnapshotConfigWithContext(context.Background(), validateReportsSnapshotConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ValidateReportsSnapshotConfigWithContext is an alternate form of the ValidateReportsSnapshotConfig method which supports a Context parameter
func (usageReports *UsageReportsV4) ValidateReportsSnapshotConfigWithContext(ctx context.Context, validateReportsSnapshotConfigOptions *ValidateReportsSnapshotConfigOptions) (result *SnapshotConfigValidateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateReportsSnapshotConfigOptions, "validateReportsSnapshotConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(validateReportsSnapshotConfigOptions, "validateReportsSnapshotConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshot-config/validate`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range validateReportsSnapshotConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "ValidateReportsSnapshotConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if validateReportsSnapshotConfigOptions.AccountID != nil {
		body["account_id"] = validateReportsSnapshotConfigOptions.AccountID
	}
	if validateReportsSnapshotConfigOptions.Interval != nil {
		body["interval"] = validateReportsSnapshotConfigOptions.Interval
	}
	if validateReportsSnapshotConfigOptions.CosBucket != nil {
		body["cos_bucket"] = validateReportsSnapshotConfigOptions.CosBucket
	}
	if validateReportsSnapshotConfigOptions.CosLocation != nil {
		body["cos_location"] = validateReportsSnapshotConfigOptions.CosLocation
	}
	if validateReportsSnapshotConfigOptions.CosReportsFolder != nil {
		body["cos_reports_folder"] = validateReportsSnapshotConfigOptions.CosReportsFolder
	}
	if validateReportsSnapshotConfigOptions.ReportTypes != nil {
		body["report_types"] = validateReportsSnapshotConfigOptions.ReportTypes
	}
	if validateReportsSnapshotConfigOptions.Versioning != nil {
		body["versioning"] = validateReportsSnapshotConfigOptions.Versioning
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
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "validate_reports_snapshot_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSnapshotConfigValidateResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetReportsSnapshot : Fetch the current or past snapshots
// Returns the billing reports snapshots captured for the given Account Id in the specific time period.
func (usageReports *UsageReportsV4) GetReportsSnapshot(getReportsSnapshotOptions *GetReportsSnapshotOptions) (result *SnapshotList, response *core.DetailedResponse, err error) {
	result, response, err = usageReports.GetReportsSnapshotWithContext(context.Background(), getReportsSnapshotOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetReportsSnapshotWithContext is an alternate form of the GetReportsSnapshot method which supports a Context parameter
func (usageReports *UsageReportsV4) GetReportsSnapshotWithContext(ctx context.Context, getReportsSnapshotOptions *GetReportsSnapshotOptions) (result *SnapshotList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportsSnapshotOptions, "getReportsSnapshotOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getReportsSnapshotOptions, "getReportsSnapshotOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = usageReports.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(usageReports.Service.Options.URL, `/v1/billing-reports-snapshots`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getReportsSnapshotOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("usage_reports", "V4", "GetReportsSnapshot")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("account_id", fmt.Sprint(*getReportsSnapshotOptions.AccountID))
	builder.AddQuery("month", fmt.Sprint(*getReportsSnapshotOptions.Month))
	if getReportsSnapshotOptions.DateFrom != nil {
		builder.AddQuery("date_from", fmt.Sprint(*getReportsSnapshotOptions.DateFrom))
	}
	if getReportsSnapshotOptions.DateTo != nil {
		builder.AddQuery("date_to", fmt.Sprint(*getReportsSnapshotOptions.DateTo))
	}
	if getReportsSnapshotOptions.Limit != nil {
		builder.AddQuery("_limit", fmt.Sprint(*getReportsSnapshotOptions.Limit))
	}
	if getReportsSnapshotOptions.Start != nil {
		builder.AddQuery("_start", fmt.Sprint(*getReportsSnapshotOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = usageReports.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_reports_snapshot", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSnapshotList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "4.0.6")
}

// AccountSummary : A summary of charges and credits for an account.
type AccountSummary struct {
	// The ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// The list of account resources for the month.
	AccountResources []Resource `json:"account_resources,omitempty"`

	// The month in which usages were incurred. Represented in yyyy-mm format.
	Month *string `json:"month" validate:"required"`

	// Country.
	BillingCountryCode *string `json:"billing_country_code" validate:"required"`

	// The currency in which the account is billed.
	BillingCurrencyCode *string `json:"billing_currency_code" validate:"required"`

	// Charges related to cloud resources.
	Resources *ResourcesSummary `json:"resources" validate:"required"`

	// The list of offers applicable for the account for the month.
	Offers []Offer `json:"offers" validate:"required"`

	// Support-related charges.
	Support []SupportSummary `json:"support" validate:"required"`

	// The list of support resources for the month.
	SupportResources []interface{} `json:"support_resources,omitempty"`

	// A summary of charges and credits related to a subscription.
	Subscription *SubscriptionSummary `json:"subscription" validate:"required"`
}

// UnmarshalAccountSummary unmarshals an instance of AccountSummary from the specified map of raw messages.
func UnmarshalAccountSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSummary)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "account_resources", &obj.AccountResources, UnmarshalResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billing_country_code", &obj.BillingCountryCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "billing_country_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billing_currency_code", &obj.BillingCurrencyCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "billing_currency_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResourcesSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "offers", &obj.Offers, UnmarshalOffer)
	if err != nil {
		err = core.SDKErrorf(err, "", "offers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "support", &obj.Support, UnmarshalSupportSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "support-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "support_resources", &obj.SupportResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subscription", &obj.Subscription, UnmarshalSubscriptionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "subscription-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountUsage : The aggregated usage and charges for all the plans in the account.
type AccountUsage struct {
	// The ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// The target country pricing that should be used.
	PricingCountry *string `json:"pricing_country" validate:"required"`

	// The currency for the cost fields in the resources, plans and metrics.
	CurrencyCode *string `json:"currency_code" validate:"required"`

	// The month.
	Month *string `json:"month" validate:"required"`

	// All the resource used in the account.
	Resources []Resource `json:"resources" validate:"required"`

	// The value of the account's currency in USD.
	CurrencyRate *float64 `json:"currency_rate,omitempty"`
}

// UnmarshalAccountUsage unmarshals an instance of AccountUsage from the specified map of raw messages.
func UnmarshalAccountUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountUsage)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_country", &obj.PricingCountry)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_country-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_code", &obj.CurrencyCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_rate", &obj.CurrencyRate)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_rate-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateReportsSnapshotConfigOptions : The CreateReportsSnapshotConfig options.
type CreateReportsSnapshotConfigOptions struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id" validate:"required"`

	// Frequency of taking the snapshot of the billing reports.
	Interval *string `json:"interval" validate:"required"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket" validate:"required"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location" validate:"required"`

	// The billing reports root folder to store the billing reports snapshots. Defaults to "IBMCloud-Billing-Reports".
	CosReportsFolder *string `json:"cos_reports_folder,omitempty"`

	// The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes []string `json:"report_types,omitempty"`

	// A new version of report is created or the existing report version is overwritten with every update. Defaults to
	// "new".
	Versioning *string `json:"versioning,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateReportsSnapshotConfigOptions.Interval property.
// Frequency of taking the snapshot of the billing reports.
const (
	CreateReportsSnapshotConfigOptionsIntervalDailyConst = "daily"
)

// Constants associated with the CreateReportsSnapshotConfigOptions.ReportTypes property.
const (
	CreateReportsSnapshotConfigOptionsReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	CreateReportsSnapshotConfigOptionsReportTypesAccountSummaryConst = "account_summary"
	CreateReportsSnapshotConfigOptionsReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// Constants associated with the CreateReportsSnapshotConfigOptions.Versioning property.
// A new version of report is created or the existing report version is overwritten with every update. Defaults to
// "new".
const (
	CreateReportsSnapshotConfigOptionsVersioningNewConst = "new"
	CreateReportsSnapshotConfigOptionsVersioningOverwriteConst = "overwrite"
)

// NewCreateReportsSnapshotConfigOptions : Instantiate CreateReportsSnapshotConfigOptions
func (*UsageReportsV4) NewCreateReportsSnapshotConfigOptions(accountID string, interval string, cosBucket string, cosLocation string) *CreateReportsSnapshotConfigOptions {
	return &CreateReportsSnapshotConfigOptions{
		AccountID: core.StringPtr(accountID),
		Interval: core.StringPtr(interval),
		CosBucket: core.StringPtr(cosBucket),
		CosLocation: core.StringPtr(cosLocation),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateReportsSnapshotConfigOptions) SetAccountID(accountID string) *CreateReportsSnapshotConfigOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *CreateReportsSnapshotConfigOptions) SetInterval(interval string) *CreateReportsSnapshotConfigOptions {
	_options.Interval = core.StringPtr(interval)
	return _options
}

// SetCosBucket : Allow user to set CosBucket
func (_options *CreateReportsSnapshotConfigOptions) SetCosBucket(cosBucket string) *CreateReportsSnapshotConfigOptions {
	_options.CosBucket = core.StringPtr(cosBucket)
	return _options
}

// SetCosLocation : Allow user to set CosLocation
func (_options *CreateReportsSnapshotConfigOptions) SetCosLocation(cosLocation string) *CreateReportsSnapshotConfigOptions {
	_options.CosLocation = core.StringPtr(cosLocation)
	return _options
}

// SetCosReportsFolder : Allow user to set CosReportsFolder
func (_options *CreateReportsSnapshotConfigOptions) SetCosReportsFolder(cosReportsFolder string) *CreateReportsSnapshotConfigOptions {
	_options.CosReportsFolder = core.StringPtr(cosReportsFolder)
	return _options
}

// SetReportTypes : Allow user to set ReportTypes
func (_options *CreateReportsSnapshotConfigOptions) SetReportTypes(reportTypes []string) *CreateReportsSnapshotConfigOptions {
	_options.ReportTypes = reportTypes
	return _options
}

// SetVersioning : Allow user to set Versioning
func (_options *CreateReportsSnapshotConfigOptions) SetVersioning(versioning string) *CreateReportsSnapshotConfigOptions {
	_options.Versioning = core.StringPtr(versioning)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateReportsSnapshotConfigOptions) SetHeaders(param map[string]string) *CreateReportsSnapshotConfigOptions {
	options.Headers = param
	return options
}

// DeleteReportsSnapshotConfigOptions : The DeleteReportsSnapshotConfig options.
type DeleteReportsSnapshotConfigOptions struct {
	// Account ID for which the billing report snapshot is configured.
	AccountID *string `json:"account_id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteReportsSnapshotConfigOptions : Instantiate DeleteReportsSnapshotConfigOptions
func (*UsageReportsV4) NewDeleteReportsSnapshotConfigOptions(accountID string) *DeleteReportsSnapshotConfigOptions {
	return &DeleteReportsSnapshotConfigOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteReportsSnapshotConfigOptions) SetAccountID(accountID string) *DeleteReportsSnapshotConfigOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteReportsSnapshotConfigOptions) SetHeaders(param map[string]string) *DeleteReportsSnapshotConfigOptions {
	options.Headers = param
	return options
}

// Discount : Information about a discount that is associated with a metric.
type Discount struct {
	// The reference ID of the discount.
	Ref *string `json:"ref" validate:"required"`

	// The name of the discount indicating category.
	Name *string `json:"name,omitempty"`

	// The name of the discount.
	DisplayName *string `json:"display_name,omitempty"`

	// The discount percentage.
	Discount *float64 `json:"discount" validate:"required"`
}

// UnmarshalDiscount unmarshals an instance of Discount from the specified map of raw messages.
func UnmarshalDiscount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Discount)
	err = core.UnmarshalPrimitive(m, "ref", &obj.Ref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ref-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "discount", &obj.Discount)
	if err != nil {
		err = core.SDKErrorf(err, "", "discount-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccountSummaryOptions : The GetAccountSummary options.
type GetAccountSummaryOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAccountSummaryOptions : Instantiate GetAccountSummaryOptions
func (*UsageReportsV4) NewGetAccountSummaryOptions(accountID string, billingmonth string) *GetAccountSummaryOptions {
	return &GetAccountSummaryOptions{
		AccountID: core.StringPtr(accountID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountSummaryOptions) SetAccountID(accountID string) *GetAccountSummaryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetAccountSummaryOptions) SetBillingmonth(billingmonth string) *GetAccountSummaryOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSummaryOptions) SetHeaders(param map[string]string) *GetAccountSummaryOptions {
	options.Headers = param
	return options
}

// GetAccountUsageOptions : The GetAccountUsage options.
type GetAccountUsageOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAccountUsageOptions : Instantiate GetAccountUsageOptions
func (*UsageReportsV4) NewGetAccountUsageOptions(accountID string, billingmonth string) *GetAccountUsageOptions {
	return &GetAccountUsageOptions{
		AccountID: core.StringPtr(accountID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountUsageOptions) SetAccountID(accountID string) *GetAccountUsageOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetAccountUsageOptions) SetBillingmonth(billingmonth string) *GetAccountUsageOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetAccountUsageOptions) SetNames(names bool) *GetAccountUsageOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetAccountUsageOptions) SetAcceptLanguage(acceptLanguage string) *GetAccountUsageOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountUsageOptions) SetHeaders(param map[string]string) *GetAccountUsageOptions {
	options.Headers = param
	return options
}

// GetOrgUsageOptions : The GetOrgUsage options.
type GetOrgUsageOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// ID of the organization.
	OrganizationID *string `json:"organization_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetOrgUsageOptions : Instantiate GetOrgUsageOptions
func (*UsageReportsV4) NewGetOrgUsageOptions(accountID string, organizationID string, billingmonth string) *GetOrgUsageOptions {
	return &GetOrgUsageOptions{
		AccountID: core.StringPtr(accountID),
		OrganizationID: core.StringPtr(organizationID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetOrgUsageOptions) SetAccountID(accountID string) *GetOrgUsageOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetOrganizationID : Allow user to set OrganizationID
func (_options *GetOrgUsageOptions) SetOrganizationID(organizationID string) *GetOrgUsageOptions {
	_options.OrganizationID = core.StringPtr(organizationID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetOrgUsageOptions) SetBillingmonth(billingmonth string) *GetOrgUsageOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetOrgUsageOptions) SetNames(names bool) *GetOrgUsageOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetOrgUsageOptions) SetAcceptLanguage(acceptLanguage string) *GetOrgUsageOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOrgUsageOptions) SetHeaders(param map[string]string) *GetOrgUsageOptions {
	options.Headers = param
	return options
}

// GetReportsSnapshotConfigOptions : The GetReportsSnapshotConfig options.
type GetReportsSnapshotConfigOptions struct {
	// Account ID for which the billing report snapshot is configured.
	AccountID *string `json:"account_id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetReportsSnapshotConfigOptions : Instantiate GetReportsSnapshotConfigOptions
func (*UsageReportsV4) NewGetReportsSnapshotConfigOptions(accountID string) *GetReportsSnapshotConfigOptions {
	return &GetReportsSnapshotConfigOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetReportsSnapshotConfigOptions) SetAccountID(accountID string) *GetReportsSnapshotConfigOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportsSnapshotConfigOptions) SetHeaders(param map[string]string) *GetReportsSnapshotConfigOptions {
	options.Headers = param
	return options
}

// GetReportsSnapshotOptions : The GetReportsSnapshot options.
type GetReportsSnapshotOptions struct {
	// Account ID for which the billing report snapshot is requested.
	AccountID *string `json:"account_id" validate:"required"`

	// The month for which billing report snapshot is requested.  Format is yyyy-mm.
	Month *string `json:"month" validate:"required"`

	// Timestamp in milliseconds for which billing report snapshot is requested.
	DateFrom *int64 `json:"date_from,omitempty"`

	// Timestamp in milliseconds for which billing report snapshot is requested.
	DateTo *int64 `json:"date_to,omitempty"`

	// Number of usage records returned. The default value is 30. Maximum value is 200.
	Limit *int64 `json:"_limit,omitempty"`

	// The offset from which the records must be fetched. Offset information is included in the response.
	Start *string `json:"_start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetReportsSnapshotOptions : Instantiate GetReportsSnapshotOptions
func (*UsageReportsV4) NewGetReportsSnapshotOptions(accountID string, month string) *GetReportsSnapshotOptions {
	return &GetReportsSnapshotOptions{
		AccountID: core.StringPtr(accountID),
		Month: core.StringPtr(month),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetReportsSnapshotOptions) SetAccountID(accountID string) *GetReportsSnapshotOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetMonth : Allow user to set Month
func (_options *GetReportsSnapshotOptions) SetMonth(month string) *GetReportsSnapshotOptions {
	_options.Month = core.StringPtr(month)
	return _options
}

// SetDateFrom : Allow user to set DateFrom
func (_options *GetReportsSnapshotOptions) SetDateFrom(dateFrom int64) *GetReportsSnapshotOptions {
	_options.DateFrom = core.Int64Ptr(dateFrom)
	return _options
}

// SetDateTo : Allow user to set DateTo
func (_options *GetReportsSnapshotOptions) SetDateTo(dateTo int64) *GetReportsSnapshotOptions {
	_options.DateTo = core.Int64Ptr(dateTo)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetReportsSnapshotOptions) SetLimit(limit int64) *GetReportsSnapshotOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetReportsSnapshotOptions) SetStart(start string) *GetReportsSnapshotOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportsSnapshotOptions) SetHeaders(param map[string]string) *GetReportsSnapshotOptions {
	options.Headers = param
	return options
}

// GetResourceGroupUsageOptions : The GetResourceGroupUsage options.
type GetResourceGroupUsageOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Resource group for which the usage report is requested.
	ResourceGroupID *string `json:"resource_group_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceGroupUsageOptions : Instantiate GetResourceGroupUsageOptions
func (*UsageReportsV4) NewGetResourceGroupUsageOptions(accountID string, resourceGroupID string, billingmonth string) *GetResourceGroupUsageOptions {
	return &GetResourceGroupUsageOptions{
		AccountID: core.StringPtr(accountID),
		ResourceGroupID: core.StringPtr(resourceGroupID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetResourceGroupUsageOptions) SetAccountID(accountID string) *GetResourceGroupUsageOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *GetResourceGroupUsageOptions) SetResourceGroupID(resourceGroupID string) *GetResourceGroupUsageOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetResourceGroupUsageOptions) SetBillingmonth(billingmonth string) *GetResourceGroupUsageOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetResourceGroupUsageOptions) SetNames(names bool) *GetResourceGroupUsageOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetResourceGroupUsageOptions) SetAcceptLanguage(acceptLanguage string) *GetResourceGroupUsageOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceGroupUsageOptions) SetHeaders(param map[string]string) *GetResourceGroupUsageOptions {
	options.Headers = param
	return options
}

// GetResourceUsageAccountOptions : The GetResourceUsageAccount options.
type GetResourceUsageAccountOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Include the tags associated with every resource instance. By default it is always `true`.
	Tags *bool `json:"_tags,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Number of usage records returned. The default value is 30. Maximum value is 200.
	Limit *int64 `json:"_limit,omitempty"`

	// The offset from which the records must be fetched. Offset information is included in the response.
	Start *string `json:"_start,omitempty"`

	// Filter by resource group.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Filter by organization_id.
	OrganizationID *string `json:"organization_id,omitempty"`

	// Filter by resource instance_id.
	ResourceInstanceID *string `json:"resource_instance_id,omitempty"`

	// Filter by resource_id.
	ResourceID *string `json:"resource_id,omitempty"`

	// Filter by plan_id.
	PlanID *string `json:"plan_id,omitempty"`

	// Region in which the resource instance is provisioned.
	Region *string `json:"region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceUsageAccountOptions : Instantiate GetResourceUsageAccountOptions
func (*UsageReportsV4) NewGetResourceUsageAccountOptions(accountID string, billingmonth string) *GetResourceUsageAccountOptions {
	return &GetResourceUsageAccountOptions{
		AccountID: core.StringPtr(accountID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetResourceUsageAccountOptions) SetAccountID(accountID string) *GetResourceUsageAccountOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetResourceUsageAccountOptions) SetBillingmonth(billingmonth string) *GetResourceUsageAccountOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetResourceUsageAccountOptions) SetNames(names bool) *GetResourceUsageAccountOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *GetResourceUsageAccountOptions) SetTags(tags bool) *GetResourceUsageAccountOptions {
	_options.Tags = core.BoolPtr(tags)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetResourceUsageAccountOptions) SetAcceptLanguage(acceptLanguage string) *GetResourceUsageAccountOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetResourceUsageAccountOptions) SetLimit(limit int64) *GetResourceUsageAccountOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetResourceUsageAccountOptions) SetStart(start string) *GetResourceUsageAccountOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *GetResourceUsageAccountOptions) SetResourceGroupID(resourceGroupID string) *GetResourceUsageAccountOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetOrganizationID : Allow user to set OrganizationID
func (_options *GetResourceUsageAccountOptions) SetOrganizationID(organizationID string) *GetResourceUsageAccountOptions {
	_options.OrganizationID = core.StringPtr(organizationID)
	return _options
}

// SetResourceInstanceID : Allow user to set ResourceInstanceID
func (_options *GetResourceUsageAccountOptions) SetResourceInstanceID(resourceInstanceID string) *GetResourceUsageAccountOptions {
	_options.ResourceInstanceID = core.StringPtr(resourceInstanceID)
	return _options
}

// SetResourceID : Allow user to set ResourceID
func (_options *GetResourceUsageAccountOptions) SetResourceID(resourceID string) *GetResourceUsageAccountOptions {
	_options.ResourceID = core.StringPtr(resourceID)
	return _options
}

// SetPlanID : Allow user to set PlanID
func (_options *GetResourceUsageAccountOptions) SetPlanID(planID string) *GetResourceUsageAccountOptions {
	_options.PlanID = core.StringPtr(planID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetResourceUsageAccountOptions) SetRegion(region string) *GetResourceUsageAccountOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceUsageAccountOptions) SetHeaders(param map[string]string) *GetResourceUsageAccountOptions {
	options.Headers = param
	return options
}

// GetResourceUsageOrgOptions : The GetResourceUsageOrg options.
type GetResourceUsageOrgOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// ID of the organization.
	OrganizationID *string `json:"organization_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Include the tags associated with every resource instance. By default it is always `true`.
	Tags *bool `json:"_tags,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Number of usage records returned. The default value is 30. Maximum value is 200.
	Limit *int64 `json:"_limit,omitempty"`

	// The offset from which the records must be fetched. Offset information is included in the response.
	Start *string `json:"_start,omitempty"`

	// Filter by resource instance id.
	ResourceInstanceID *string `json:"resource_instance_id,omitempty"`

	// Filter by resource_id.
	ResourceID *string `json:"resource_id,omitempty"`

	// Filter by plan_id.
	PlanID *string `json:"plan_id,omitempty"`

	// Region in which the resource instance is provisioned.
	Region *string `json:"region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceUsageOrgOptions : Instantiate GetResourceUsageOrgOptions
func (*UsageReportsV4) NewGetResourceUsageOrgOptions(accountID string, organizationID string, billingmonth string) *GetResourceUsageOrgOptions {
	return &GetResourceUsageOrgOptions{
		AccountID: core.StringPtr(accountID),
		OrganizationID: core.StringPtr(organizationID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetResourceUsageOrgOptions) SetAccountID(accountID string) *GetResourceUsageOrgOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetOrganizationID : Allow user to set OrganizationID
func (_options *GetResourceUsageOrgOptions) SetOrganizationID(organizationID string) *GetResourceUsageOrgOptions {
	_options.OrganizationID = core.StringPtr(organizationID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetResourceUsageOrgOptions) SetBillingmonth(billingmonth string) *GetResourceUsageOrgOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetResourceUsageOrgOptions) SetNames(names bool) *GetResourceUsageOrgOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *GetResourceUsageOrgOptions) SetTags(tags bool) *GetResourceUsageOrgOptions {
	_options.Tags = core.BoolPtr(tags)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetResourceUsageOrgOptions) SetAcceptLanguage(acceptLanguage string) *GetResourceUsageOrgOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetResourceUsageOrgOptions) SetLimit(limit int64) *GetResourceUsageOrgOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetResourceUsageOrgOptions) SetStart(start string) *GetResourceUsageOrgOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetResourceInstanceID : Allow user to set ResourceInstanceID
func (_options *GetResourceUsageOrgOptions) SetResourceInstanceID(resourceInstanceID string) *GetResourceUsageOrgOptions {
	_options.ResourceInstanceID = core.StringPtr(resourceInstanceID)
	return _options
}

// SetResourceID : Allow user to set ResourceID
func (_options *GetResourceUsageOrgOptions) SetResourceID(resourceID string) *GetResourceUsageOrgOptions {
	_options.ResourceID = core.StringPtr(resourceID)
	return _options
}

// SetPlanID : Allow user to set PlanID
func (_options *GetResourceUsageOrgOptions) SetPlanID(planID string) *GetResourceUsageOrgOptions {
	_options.PlanID = core.StringPtr(planID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetResourceUsageOrgOptions) SetRegion(region string) *GetResourceUsageOrgOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceUsageOrgOptions) SetHeaders(param map[string]string) *GetResourceUsageOrgOptions {
	options.Headers = param
	return options
}

// GetResourceUsageResourceGroupOptions : The GetResourceUsageResourceGroup options.
type GetResourceUsageResourceGroupOptions struct {
	// Account ID for which the usage report is requested.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// Resource group for which the usage report is requested.
	ResourceGroupID *string `json:"resource_group_id" validate:"required,ne="`

	// The billing month for which the usage report is requested.  Format is yyyy-mm.
	Billingmonth *string `json:"billingmonth" validate:"required,ne="`

	// Include the name of every resource, plan, resource instance, organization, and resource group.
	Names *bool `json:"_names,omitempty"`

	// Include the tags associated with every resource instance. By default it is always `true`.
	Tags *bool `json:"_tags,omitempty"`

	// Prioritize the names returned in the order of the specified languages. Language will default to English.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Number of usage records returned. The default value is 30. Maximum value is 200.
	Limit *int64 `json:"_limit,omitempty"`

	// The offset from which the records must be fetched. Offset information is included in the response.
	Start *string `json:"_start,omitempty"`

	// Filter by resource instance id.
	ResourceInstanceID *string `json:"resource_instance_id,omitempty"`

	// Filter by resource_id.
	ResourceID *string `json:"resource_id,omitempty"`

	// Filter by plan_id.
	PlanID *string `json:"plan_id,omitempty"`

	// Region in which the resource instance is provisioned.
	Region *string `json:"region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceUsageResourceGroupOptions : Instantiate GetResourceUsageResourceGroupOptions
func (*UsageReportsV4) NewGetResourceUsageResourceGroupOptions(accountID string, resourceGroupID string, billingmonth string) *GetResourceUsageResourceGroupOptions {
	return &GetResourceUsageResourceGroupOptions{
		AccountID: core.StringPtr(accountID),
		ResourceGroupID: core.StringPtr(resourceGroupID),
		Billingmonth: core.StringPtr(billingmonth),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetResourceUsageResourceGroupOptions) SetAccountID(accountID string) *GetResourceUsageResourceGroupOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *GetResourceUsageResourceGroupOptions) SetResourceGroupID(resourceGroupID string) *GetResourceUsageResourceGroupOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetBillingmonth : Allow user to set Billingmonth
func (_options *GetResourceUsageResourceGroupOptions) SetBillingmonth(billingmonth string) *GetResourceUsageResourceGroupOptions {
	_options.Billingmonth = core.StringPtr(billingmonth)
	return _options
}

// SetNames : Allow user to set Names
func (_options *GetResourceUsageResourceGroupOptions) SetNames(names bool) *GetResourceUsageResourceGroupOptions {
	_options.Names = core.BoolPtr(names)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *GetResourceUsageResourceGroupOptions) SetTags(tags bool) *GetResourceUsageResourceGroupOptions {
	_options.Tags = core.BoolPtr(tags)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetResourceUsageResourceGroupOptions) SetAcceptLanguage(acceptLanguage string) *GetResourceUsageResourceGroupOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *GetResourceUsageResourceGroupOptions) SetLimit(limit int64) *GetResourceUsageResourceGroupOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *GetResourceUsageResourceGroupOptions) SetStart(start string) *GetResourceUsageResourceGroupOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetResourceInstanceID : Allow user to set ResourceInstanceID
func (_options *GetResourceUsageResourceGroupOptions) SetResourceInstanceID(resourceInstanceID string) *GetResourceUsageResourceGroupOptions {
	_options.ResourceInstanceID = core.StringPtr(resourceInstanceID)
	return _options
}

// SetResourceID : Allow user to set ResourceID
func (_options *GetResourceUsageResourceGroupOptions) SetResourceID(resourceID string) *GetResourceUsageResourceGroupOptions {
	_options.ResourceID = core.StringPtr(resourceID)
	return _options
}

// SetPlanID : Allow user to set PlanID
func (_options *GetResourceUsageResourceGroupOptions) SetPlanID(planID string) *GetResourceUsageResourceGroupOptions {
	_options.PlanID = core.StringPtr(planID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *GetResourceUsageResourceGroupOptions) SetRegion(region string) *GetResourceUsageResourceGroupOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceUsageResourceGroupOptions) SetHeaders(param map[string]string) *GetResourceUsageResourceGroupOptions {
	options.Headers = param
	return options
}

// InstanceUsage : The aggregated usage and charges for an instance.
type InstanceUsage struct {
	// The ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// The ID of the resource instance.
	ResourceInstanceID *string `json:"resource_instance_id" validate:"required"`

	// The name of the resource instance.
	ResourceInstanceName *string `json:"resource_instance_name,omitempty"`

	// The ID of the resource.
	ResourceID *string `json:"resource_id" validate:"required"`

	// The catalog ID of the resource.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty"`

	// The ID of the resource group.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// The name of the resource group.
	ResourceGroupName *string `json:"resource_group_name,omitempty"`

	// The ID of the organization.
	OrganizationID *string `json:"organization_id,omitempty"`

	// The name of the organization.
	OrganizationName *string `json:"organization_name,omitempty"`

	// The ID of the space.
	SpaceID *string `json:"space_id,omitempty"`

	// The name of the space.
	SpaceName *string `json:"space_name,omitempty"`

	// The ID of the consumer.
	ConsumerID *string `json:"consumer_id,omitempty"`

	// The region where instance was provisioned.
	Region *string `json:"region,omitempty"`

	// The pricing region where the usage that was submitted was rated.
	PricingRegion *string `json:"pricing_region,omitempty"`

	// The target country pricing that should be used.
	PricingCountry *string `json:"pricing_country" validate:"required"`

	// The currency for the cost fields in the resources, plans and metrics.
	CurrencyCode *string `json:"currency_code" validate:"required"`

	// Is the cost charged to the account.
	Billable *bool `json:"billable" validate:"required"`

	// The resource instance id of the parent resource associated with this instance.
	ParentResourceInstanceID *string `json:"parent_resource_instance_id,omitempty"`

	// The ID of the plan where the instance was provisioned and rated.
	PlanID *string `json:"plan_id" validate:"required"`

	// The name of the plan where the instance was provisioned and rated.
	PlanName *string `json:"plan_name,omitempty"`

	// The ID of the pricing plan used to rate the usage.
	PricingPlanID *string `json:"pricing_plan_id,omitempty"`

	// The month.
	Month *string `json:"month" validate:"required"`

	// All the resource used in the account.
	Usage []Metric `json:"usage" validate:"required"`

	// Pending charge from classic infrastructure.
	Pending *bool `json:"pending,omitempty"`

	// The value of the account's currency in USD.
	CurrencyRate *float64 `json:"currency_rate,omitempty"`

	// The user tags associated with a resource instance.
	Tags []interface{} `json:"tags,omitempty"`

	// The service tags associated with a resource instance.
	ServiceTags []interface{} `json:"service_tags,omitempty"`
}

// UnmarshalInstanceUsage unmarshals an instance of InstanceUsage from the specified map of raw messages.
func UnmarshalInstanceUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstanceUsage)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_instance_id", &obj.ResourceInstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_instance_name", &obj.ResourceInstanceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_instance_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "organization_id", &obj.OrganizationID)
	if err != nil {
		err = core.SDKErrorf(err, "", "organization_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "organization_name", &obj.OrganizationName)
	if err != nil {
		err = core.SDKErrorf(err, "", "organization_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "space_id", &obj.SpaceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "space_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "space_name", &obj.SpaceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "space_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "consumer_id", &obj.ConsumerID)
	if err != nil {
		err = core.SDKErrorf(err, "", "consumer_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_region", &obj.PricingRegion)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_country", &obj.PricingCountry)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_country-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_code", &obj.CurrencyCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billable", &obj.Billable)
	if err != nil {
		err = core.SDKErrorf(err, "", "billable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent_resource_instance_id", &obj.ParentResourceInstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent_resource_instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_id", &obj.PlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_name", &obj.PlanName)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_plan_id", &obj.PricingPlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_plan_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "usage", &obj.Usage, UnmarshalMetric)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		err = core.SDKErrorf(err, "", "pending-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_rate", &obj.CurrencyRate)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_rate-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_tags", &obj.ServiceTags)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstancesUsageFirst : The link to the first page of the search query.
type InstancesUsageFirst struct {
	// A link to a page of query results.
	Href *string `json:"href,omitempty"`
}

// UnmarshalInstancesUsageFirst unmarshals an instance of InstancesUsageFirst from the specified map of raw messages.
func UnmarshalInstancesUsageFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstancesUsageFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstancesUsageNext : The link to the next page of the search query.
type InstancesUsageNext struct {
	// A link to a page of query results.
	Href *string `json:"href,omitempty"`

	// The value of the `_start` query parameter to fetch the next page.
	Offset *string `json:"offset,omitempty"`
}

// UnmarshalInstancesUsageNext unmarshals an instance of InstancesUsageNext from the specified map of raw messages.
func UnmarshalInstancesUsageNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstancesUsageNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstancesUsage : The list of instance usage reports.
type InstancesUsage struct {
	// The max number of reports in the response.
	Limit *int64 `json:"limit,omitempty"`

	// The number of reports in the response.
	Count *int64 `json:"count,omitempty"`

	// The link to the first page of the search query.
	First *InstancesUsageFirst `json:"first,omitempty"`

	// The link to the next page of the search query.
	Next *InstancesUsageNext `json:"next,omitempty"`

	// The list of instance usage reports.
	Resources []InstanceUsage `json:"resources,omitempty"`
}

// UnmarshalInstancesUsage unmarshals an instance of InstancesUsage from the specified map of raw messages.
func UnmarshalInstancesUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InstancesUsage)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalInstancesUsageFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalInstancesUsageNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalInstanceUsage)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *InstancesUsage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	_start, err := core.GetQueryParam(resp.Next.Href, "_start")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if _start == nil {
		return nil, nil
	}
	return _start, nil
}

// Metric : Information about a metric.
type Metric struct {
	// The ID of the metric.
	Metric *string `json:"metric" validate:"required"`

	// The name of the metric.
	MetricName *string `json:"metric_name,omitempty"`

	// The aggregated value for the metric.
	Quantity *float64 `json:"quantity" validate:"required"`

	// The quantity that is used for calculating charges.
	RateableQuantity *float64 `json:"rateable_quantity,omitempty"`

	// The cost incurred by the metric.
	Cost *float64 `json:"cost" validate:"required"`

	// Pre-discounted cost incurred by the metric.
	RatedCost *float64 `json:"rated_cost" validate:"required"`

	// The price with which the cost was calculated.
	Price []interface{} `json:"price,omitempty"`

	// The unit that qualifies the quantity.
	Unit *string `json:"unit,omitempty"`

	// The name of the unit.
	UnitName *string `json:"unit_name,omitempty"`

	// When set to `true`, the cost is for informational purpose and is not included while calculating the plan charges.
	NonChargeable *bool `json:"non_chargeable,omitempty"`

	// All the discounts applicable to the metric.
	Discounts []Discount `json:"discounts" validate:"required"`

	// This percentage reflects the reduction to the original cost that you receive under a volume based pricing structure.
	VolumeDiscount *float64 `json:"volume_discount,omitempty"`

	// The original cost adjusted for volume based discounts that are applied at the account level.
	VolumeCost *float64 `json:"volume_cost,omitempty"`
}

// UnmarshalMetric unmarshals an instance of Metric from the specified map of raw messages.
func UnmarshalMetric(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Metric)
	err = core.UnmarshalPrimitive(m, "metric", &obj.Metric)
	if err != nil {
		err = core.SDKErrorf(err, "", "metric-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metric_name", &obj.MetricName)
	if err != nil {
		err = core.SDKErrorf(err, "", "metric_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "quantity", &obj.Quantity)
	if err != nil {
		err = core.SDKErrorf(err, "", "quantity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rateable_quantity", &obj.RateableQuantity)
	if err != nil {
		err = core.SDKErrorf(err, "", "rateable_quantity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cost", &obj.Cost)
	if err != nil {
		err = core.SDKErrorf(err, "", "cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rated_cost", &obj.RatedCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "rated_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "price", &obj.Price)
	if err != nil {
		err = core.SDKErrorf(err, "", "price-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unit", &obj.Unit)
	if err != nil {
		err = core.SDKErrorf(err, "", "unit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unit_name", &obj.UnitName)
	if err != nil {
		err = core.SDKErrorf(err, "", "unit_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "non_chargeable", &obj.NonChargeable)
	if err != nil {
		err = core.SDKErrorf(err, "", "non_chargeable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "discounts", &obj.Discounts, UnmarshalDiscount)
	if err != nil {
		err = core.SDKErrorf(err, "", "discounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "volume_discount", &obj.VolumeDiscount)
	if err != nil {
		err = core.SDKErrorf(err, "", "volume_discount-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "volume_cost", &obj.VolumeCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "volume_cost-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Offer : Information about an individual offer.
type Offer struct {
	// The ID of the offer.
	OfferID *string `json:"offer_id" validate:"required"`

	// The total credits before applying the offer.
	CreditsTotal *float64 `json:"credits_total" validate:"required"`

	// The template with which the offer was generated.
	OfferTemplate *string `json:"offer_template" validate:"required"`

	// The date from which the offer is valid.
	ValidFrom *strfmt.DateTime `json:"valid_from" validate:"required"`

	// The offer's creator's email id.
	CreatedByEmailID *string `json:"created_by_email_id" validate:"required"`

	// The date until the offer is valid.
	ExpiresOn *strfmt.DateTime `json:"expires_on" validate:"required"`

	// Credit information related to an offer.
	Credits *OfferCredits `json:"credits" validate:"required"`
}

// UnmarshalOffer unmarshals an instance of Offer from the specified map of raw messages.
func UnmarshalOffer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Offer)
	err = core.UnmarshalPrimitive(m, "offer_id", &obj.OfferID)
	if err != nil {
		err = core.SDKErrorf(err, "", "offer_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "credits_total", &obj.CreditsTotal)
	if err != nil {
		err = core.SDKErrorf(err, "", "credits_total-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offer_template", &obj.OfferTemplate)
	if err != nil {
		err = core.SDKErrorf(err, "", "offer_template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "valid_from", &obj.ValidFrom)
	if err != nil {
		err = core.SDKErrorf(err, "", "valid_from-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_email_id", &obj.CreatedByEmailID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_email_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "expires_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "credits", &obj.Credits, UnmarshalOfferCredits)
	if err != nil {
		err = core.SDKErrorf(err, "", "credits-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OfferCredits : Credit information related to an offer.
type OfferCredits struct {
	// The available credits in the offer at the beginning of the month.
	StartingBalance *float64 `json:"starting_balance" validate:"required"`

	// The credits used in this month.
	Used *float64 `json:"used" validate:"required"`

	// The remaining credits in the offer.
	Balance *float64 `json:"balance" validate:"required"`
}

// UnmarshalOfferCredits unmarshals an instance of OfferCredits from the specified map of raw messages.
func UnmarshalOfferCredits(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OfferCredits)
	err = core.UnmarshalPrimitive(m, "starting_balance", &obj.StartingBalance)
	if err != nil {
		err = core.SDKErrorf(err, "", "starting_balance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "used", &obj.Used)
	if err != nil {
		err = core.SDKErrorf(err, "", "used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "balance", &obj.Balance)
	if err != nil {
		err = core.SDKErrorf(err, "", "balance-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OrgUsage : The aggregated usage and charges for all the plans in the org.
type OrgUsage struct {
	// The ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// The ID of the organization.
	OrganizationID *string `json:"organization_id" validate:"required"`

	// The name of the organization.
	OrganizationName *string `json:"organization_name,omitempty"`

	// The target country pricing that should be used.
	PricingCountry *string `json:"pricing_country" validate:"required"`

	// The currency for the cost fields in the resources, plans and metrics.
	CurrencyCode *string `json:"currency_code" validate:"required"`

	// The month.
	Month *string `json:"month" validate:"required"`

	// All the resource used in the account.
	Resources []Resource `json:"resources" validate:"required"`

	// The value of the account's currency in USD.
	CurrencyRate *float64 `json:"currency_rate,omitempty"`
}

// UnmarshalOrgUsage unmarshals an instance of OrgUsage from the specified map of raw messages.
func UnmarshalOrgUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OrgUsage)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "organization_id", &obj.OrganizationID)
	if err != nil {
		err = core.SDKErrorf(err, "", "organization_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "organization_name", &obj.OrganizationName)
	if err != nil {
		err = core.SDKErrorf(err, "", "organization_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_country", &obj.PricingCountry)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_country-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_code", &obj.CurrencyCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_rate", &obj.CurrencyRate)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_rate-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Plan : The aggregated values for the plan.
type Plan struct {
	// The ID of the plan.
	PlanID *string `json:"plan_id" validate:"required"`

	// The name of the plan.
	PlanName *string `json:"plan_name,omitempty"`

	// The pricing region for the plan.
	PricingRegion *string `json:"pricing_region,omitempty"`

	// The ID of the pricing plan used to rate the usage.
	PricingPlanID *string `json:"pricing_plan_id,omitempty"`

	// Indicates if the plan charges are billed to the customer.
	Billable *bool `json:"billable" validate:"required"`

	// The total cost incurred by the plan.
	Cost *float64 `json:"cost" validate:"required"`

	// Total pre-discounted cost incurred by the plan.
	RatedCost *float64 `json:"rated_cost" validate:"required"`

	// All the metrics in the plan.
	Usage []Metric `json:"usage" validate:"required"`

	// All the discounts applicable to the plan.
	Discounts []Discount `json:"discounts" validate:"required"`

	// Pending charge from classic infrastructure.
	Pending *bool `json:"pending,omitempty"`
}

// UnmarshalPlan unmarshals an instance of Plan from the specified map of raw messages.
func UnmarshalPlan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Plan)
	err = core.UnmarshalPrimitive(m, "plan_id", &obj.PlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_name", &obj.PlanName)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_region", &obj.PricingRegion)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_plan_id", &obj.PricingPlanID)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_plan_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billable", &obj.Billable)
	if err != nil {
		err = core.SDKErrorf(err, "", "billable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cost", &obj.Cost)
	if err != nil {
		err = core.SDKErrorf(err, "", "cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rated_cost", &obj.RatedCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "rated_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "usage", &obj.Usage, UnmarshalMetric)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "discounts", &obj.Discounts, UnmarshalDiscount)
	if err != nil {
		err = core.SDKErrorf(err, "", "discounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		err = core.SDKErrorf(err, "", "pending-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Resource : The container for all the plans in the resource.
type Resource struct {
	// The ID of the resource.
	ResourceID *string `json:"resource_id" validate:"required"`

	// The catalog ID of the resource.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty"`

	// The billable charges for the account.
	BillableCost *float64 `json:"billable_cost" validate:"required"`

	// The pre-discounted billable charges for the account.
	BillableRatedCost *float64 `json:"billable_rated_cost" validate:"required"`

	// The non-billable charges for the account.
	NonBillableCost *float64 `json:"non_billable_cost" validate:"required"`

	// The pre-discounted non-billable charges for the account.
	NonBillableRatedCost *float64 `json:"non_billable_rated_cost" validate:"required"`

	// All the plans in the resource.
	Plans []Plan `json:"plans" validate:"required"`

	// All the discounts applicable to the resource.
	Discounts []Discount `json:"discounts" validate:"required"`
}

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billable_cost", &obj.BillableCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "billable_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billable_rated_cost", &obj.BillableRatedCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "billable_rated_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "non_billable_cost", &obj.NonBillableCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "non_billable_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "non_billable_rated_cost", &obj.NonBillableRatedCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "non_billable_rated_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plans", &obj.Plans, UnmarshalPlan)
	if err != nil {
		err = core.SDKErrorf(err, "", "plans-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "discounts", &obj.Discounts, UnmarshalDiscount)
	if err != nil {
		err = core.SDKErrorf(err, "", "discounts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceGroupUsage : The aggregated usage and charges for all the plans in the resource group.
type ResourceGroupUsage struct {
	// The ID of the account.
	AccountID *string `json:"account_id" validate:"required"`

	// The ID of the resource group.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The name of the resource group.
	ResourceGroupName *string `json:"resource_group_name,omitempty"`

	// The target country pricing that should be used.
	PricingCountry *string `json:"pricing_country" validate:"required"`

	// The currency for the cost fields in the resources, plans and metrics.
	CurrencyCode *string `json:"currency_code" validate:"required"`

	// The month.
	Month *string `json:"month" validate:"required"`

	// All the resource used in the account.
	Resources []Resource `json:"resources" validate:"required"`

	// The value of the account's currency in USD.
	CurrencyRate *float64 `json:"currency_rate,omitempty"`
}

// UnmarshalResourceGroupUsage unmarshals an instance of ResourceGroupUsage from the specified map of raw messages.
func UnmarshalResourceGroupUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupUsage)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pricing_country", &obj.PricingCountry)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing_country-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_code", &obj.CurrencyCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency_rate", &obj.CurrencyRate)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency_rate-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourcesSummary : Charges related to cloud resources.
type ResourcesSummary struct {
	// The billable charges for all cloud resources used in the account.
	BillableCost *float64 `json:"billable_cost" validate:"required"`

	// Non-billable charges for all cloud resources used in the account.
	NonBillableCost *float64 `json:"non_billable_cost" validate:"required"`
}

// UnmarshalResourcesSummary unmarshals an instance of ResourcesSummary from the specified map of raw messages.
func UnmarshalResourcesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourcesSummary)
	err = core.UnmarshalPrimitive(m, "billable_cost", &obj.BillableCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "billable_cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "non_billable_cost", &obj.NonBillableCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "non_billable_cost-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotConfigHistoryItem : SnapshotConfigHistoryItem struct
type SnapshotConfigHistoryItem struct {
	// Timestamp in milliseconds when the snapshot configuration was created.
	StartTime *int64 `json:"start_time,omitempty"`

	// Timestamp in milliseconds when the snapshot configuration ends.
	EndTime *int64 `json:"end_time,omitempty"`

	// Account that updated the billing snapshot configuration.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id,omitempty"`

	// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
	State *string `json:"state,omitempty"`

	// Type of account. Possible values [enterprise, account].
	AccountType *string `json:"account_type,omitempty"`

	// Frequency of taking the snapshot of the billing reports.
	Interval *string `json:"interval,omitempty"`

	// A new version of report is created or the existing report version is overwritten with every update.
	Versioning *string `json:"versioning,omitempty"`

	// The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes []string `json:"report_types,omitempty"`

	// Compression format of the snapshot report.
	Compression *string `json:"compression,omitempty"`

	// Type of content stored in snapshot report.
	ContentType *string `json:"content_type,omitempty"`

	// The billing reports root folder to store the billing reports snapshots. Defaults to "IBMCloud-Billing-Reports".
	CosReportsFolder *string `json:"cos_reports_folder,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket,omitempty"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location,omitempty"`

	// The endpoint of the COS instance.
	CosEndpoint *string `json:"cos_endpoint,omitempty"`
}

// Constants associated with the SnapshotConfigHistoryItem.State property.
// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
const (
	SnapshotConfigHistoryItemStateDisabledConst = "disabled"
	SnapshotConfigHistoryItemStateEnabledConst = "enabled"
)

// Constants associated with the SnapshotConfigHistoryItem.AccountType property.
// Type of account. Possible values [enterprise, account].
const (
	SnapshotConfigHistoryItemAccountTypeAccountConst = "account"
	SnapshotConfigHistoryItemAccountTypeEnterpriseConst = "enterprise"
)

// Constants associated with the SnapshotConfigHistoryItem.Interval property.
// Frequency of taking the snapshot of the billing reports.
const (
	SnapshotConfigHistoryItemIntervalDailyConst = "daily"
)

// Constants associated with the SnapshotConfigHistoryItem.Versioning property.
// A new version of report is created or the existing report version is overwritten with every update.
const (
	SnapshotConfigHistoryItemVersioningNewConst = "new"
	SnapshotConfigHistoryItemVersioningOverwriteConst = "overwrite"
)

// Constants associated with the SnapshotConfigHistoryItem.ReportTypes property.
const (
	SnapshotConfigHistoryItemReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	SnapshotConfigHistoryItemReportTypesAccountSummaryConst = "account_summary"
	SnapshotConfigHistoryItemReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// UnmarshalSnapshotConfigHistoryItem unmarshals an instance of SnapshotConfigHistoryItem from the specified map of raw messages.
func UnmarshalSnapshotConfigHistoryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotConfigHistoryItem)
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "end_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_type", &obj.AccountType)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		err = core.SDKErrorf(err, "", "interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "versioning", &obj.Versioning)
	if err != nil {
		err = core.SDKErrorf(err, "", "versioning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_types", &obj.ReportTypes)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_types-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "compression", &obj.Compression)
	if err != nil {
		err = core.SDKErrorf(err, "", "compression-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		err = core.SDKErrorf(err, "", "content_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_reports_folder", &obj.CosReportsFolder)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_reports_folder-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_bucket", &obj.CosBucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_location", &obj.CosLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_endpoint", &obj.CosEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_endpoint-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotList : List of billing reports snapshots.
type SnapshotList struct {
	// Number of total snapshots.
	Count *int64 `json:"count,omitempty"`

	// Reference to the first page of the search query.
	First *SnapshotListFirst `json:"first,omitempty"`

	// Reference to the next page of the search query if any.
	Next *SnapshotListNext `json:"next,omitempty"`

	Snapshots []SnapshotListSnapshotsItem `json:"snapshots,omitempty"`
}

// UnmarshalSnapshotList unmarshals an instance of SnapshotList from the specified map of raw messages.
func UnmarshalSnapshotList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalSnapshotListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalSnapshotListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "snapshots", &obj.Snapshots, UnmarshalSnapshotListSnapshotsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "snapshots-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SnapshotList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	_start, err := core.GetQueryParam(resp.Next.Href, "_start")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if _start == nil {
		return nil, nil
	}
	return _start, nil
}

// SnapshotListFirst : Reference to the first page of the search query.
type SnapshotListFirst struct {
	Href *string `json:"href,omitempty"`
}

// UnmarshalSnapshotListFirst unmarshals an instance of SnapshotListFirst from the specified map of raw messages.
func UnmarshalSnapshotListFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotListNext : Reference to the next page of the search query if any.
type SnapshotListNext struct {
	Href *string `json:"href,omitempty"`

	// The value of the `_start` query parameter to fetch the next page.
	Offset *string `json:"offset,omitempty"`
}

// UnmarshalSnapshotListNext unmarshals an instance of SnapshotListNext from the specified map of raw messages.
func UnmarshalSnapshotListNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotListSnapshotsItem : Snapshot Schema.
type SnapshotListSnapshotsItem struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id,omitempty"`

	// Month of captured snapshot.
	Month *string `json:"month,omitempty"`

	// Type of account. Possible values are [enterprise, account].
	AccountType *string `json:"account_type,omitempty"`

	// Timestamp of snapshot processed.
	ExpectedProcessedAt *int64 `json:"expected_processed_at,omitempty"`

	// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
	State *string `json:"state,omitempty"`

	// Period of billing in snapshot.
	BillingPeriod *SnapshotListSnapshotsItemBillingPeriod `json:"billing_period,omitempty"`

	// Id of the snapshot captured.
	SnapshotID *string `json:"snapshot_id,omitempty"`

	// Character encoding used.
	Charset *string `json:"charset,omitempty"`

	// Compression format of the snapshot report.
	Compression *string `json:"compression,omitempty"`

	// Type of content stored in snapshot report.
	ContentType *string `json:"content_type,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	Bucket *string `json:"bucket,omitempty"`

	// Version of the snapshot.
	Version *string `json:"version,omitempty"`

	// Date and time of creation of snapshot.
	CreatedOn *string `json:"created_on,omitempty"`

	// List of report types configured for the snapshot.
	ReportTypes []SnapshotListSnapshotsItemReportTypesItem `json:"report_types,omitempty"`

	// List of location of reports.
	Files []SnapshotListSnapshotsItemFilesItem `json:"files,omitempty"`

	// Timestamp at which snapshot is captured.
	ProcessedAt *int64 `json:"processed_at,omitempty"`
}

// Constants associated with the SnapshotListSnapshotsItem.AccountType property.
// Type of account. Possible values are [enterprise, account].
const (
	SnapshotListSnapshotsItemAccountTypeAccountConst = "account"
	SnapshotListSnapshotsItemAccountTypeEnterpriseConst = "enterprise"
)

// Constants associated with the SnapshotListSnapshotsItem.State property.
// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
const (
	SnapshotListSnapshotsItemStateDisabledConst = "disabled"
	SnapshotListSnapshotsItemStateEnabledConst = "enabled"
)

// UnmarshalSnapshotListSnapshotsItem unmarshals an instance of SnapshotListSnapshotsItem from the specified map of raw messages.
func UnmarshalSnapshotListSnapshotsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListSnapshotsItem)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "month", &obj.Month)
	if err != nil {
		err = core.SDKErrorf(err, "", "month-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_type", &obj.AccountType)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_processed_at", &obj.ExpectedProcessedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "expected_processed_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "billing_period", &obj.BillingPeriod, UnmarshalSnapshotListSnapshotsItemBillingPeriod)
	if err != nil {
		err = core.SDKErrorf(err, "", "billing_period-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "snapshot_id", &obj.SnapshotID)
	if err != nil {
		err = core.SDKErrorf(err, "", "snapshot_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charset", &obj.Charset)
	if err != nil {
		err = core.SDKErrorf(err, "", "charset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "compression", &obj.Compression)
	if err != nil {
		err = core.SDKErrorf(err, "", "compression-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		err = core.SDKErrorf(err, "", "content_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_on-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "report_types", &obj.ReportTypes, UnmarshalSnapshotListSnapshotsItemReportTypesItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_types-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "files", &obj.Files, UnmarshalSnapshotListSnapshotsItemFilesItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "files-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "processed_at", &obj.ProcessedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "processed_at-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotListSnapshotsItemBillingPeriod : Period of billing in snapshot.
type SnapshotListSnapshotsItemBillingPeriod struct {
	// Date and time of start of billing in the respective snapshot.
	Start *string `json:"start,omitempty"`

	// Date and time of end of billing in the respective snapshot.
	End *string `json:"end,omitempty"`
}

// UnmarshalSnapshotListSnapshotsItemBillingPeriod unmarshals an instance of SnapshotListSnapshotsItemBillingPeriod from the specified map of raw messages.
func UnmarshalSnapshotListSnapshotsItemBillingPeriod(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListSnapshotsItemBillingPeriod)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end", &obj.End)
	if err != nil {
		err = core.SDKErrorf(err, "", "end-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotListSnapshotsItemFilesItem : SnapshotListSnapshotsItemFilesItem struct
type SnapshotListSnapshotsItemFilesItem struct {
	// The type of billing report stored. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes *string `json:"report_types,omitempty"`

	// Absolute path of the billing report in the COS instance.
	Location *string `json:"location,omitempty"`

	// Account ID for which billing report is captured.
	AccountID *string `json:"account_id,omitempty"`
}

// Constants associated with the SnapshotListSnapshotsItemFilesItem.ReportTypes property.
// The type of billing report stored. Possible values are [account_summary, enterprise_summary,
// account_resource_instance_usage].
const (
	SnapshotListSnapshotsItemFilesItemReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	SnapshotListSnapshotsItemFilesItemReportTypesAccountSummaryConst = "account_summary"
	SnapshotListSnapshotsItemFilesItemReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// UnmarshalSnapshotListSnapshotsItemFilesItem unmarshals an instance of SnapshotListSnapshotsItemFilesItem from the specified map of raw messages.
func UnmarshalSnapshotListSnapshotsItemFilesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListSnapshotsItemFilesItem)
	err = core.UnmarshalPrimitive(m, "report_types", &obj.ReportTypes)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_types-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotListSnapshotsItemReportTypesItem : SnapshotListSnapshotsItemReportTypesItem struct
type SnapshotListSnapshotsItemReportTypesItem struct {
	// The type of billing report of the snapshot. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	Type *string `json:"type,omitempty"`

	// Version of the snapshot.
	Version *string `json:"version,omitempty"`
}

// Constants associated with the SnapshotListSnapshotsItemReportTypesItem.Type property.
// The type of billing report of the snapshot. Possible values are [account_summary, enterprise_summary,
// account_resource_instance_usage].
const (
	SnapshotListSnapshotsItemReportTypesItemTypeAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	SnapshotListSnapshotsItemReportTypesItemTypeAccountSummaryConst = "account_summary"
	SnapshotListSnapshotsItemReportTypesItemTypeEnterpriseSummaryConst = "enterprise_summary"
)

// UnmarshalSnapshotListSnapshotsItemReportTypesItem unmarshals an instance of SnapshotListSnapshotsItemReportTypesItem from the specified map of raw messages.
func UnmarshalSnapshotListSnapshotsItemReportTypesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotListSnapshotsItemReportTypesItem)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotConfig : Billing reports snapshot configuration.
type SnapshotConfig struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id,omitempty"`

	// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
	State *string `json:"state,omitempty"`

	// Type of account. Possible values are [enterprise, account].
	AccountType *string `json:"account_type,omitempty"`

	// Frequency of taking the snapshot of the billing reports.
	Interval *string `json:"interval,omitempty"`

	// A new version of report is created or the existing report version is overwritten with every update.
	Versioning *string `json:"versioning,omitempty"`

	// The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes []string `json:"report_types,omitempty"`

	// Compression format of the snapshot report.
	Compression *string `json:"compression,omitempty"`

	// Type of content stored in snapshot report.
	ContentType *string `json:"content_type,omitempty"`

	// The billing reports root folder to store the billing reports snapshots. Defaults to "IBMCloud-Billing-Reports".
	CosReportsFolder *string `json:"cos_reports_folder,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket,omitempty"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location,omitempty"`

	// The endpoint of the COS instance.
	CosEndpoint *string `json:"cos_endpoint,omitempty"`

	// Timestamp in milliseconds when the snapshot configuration was created.
	CreatedAt *int64 `json:"created_at,omitempty"`

	// Timestamp in milliseconds when the snapshot configuration was last updated.
	LastUpdatedAt *int64 `json:"last_updated_at,omitempty"`

	// List of previous versions of the snapshot configurations.
	History []SnapshotConfigHistoryItem `json:"history,omitempty"`
}

// Constants associated with the SnapshotConfig.State property.
// Status of the billing snapshot configuration. Possible values are [enabled, disabled].
const (
	SnapshotConfigStateDisabledConst = "disabled"
	SnapshotConfigStateEnabledConst = "enabled"
)

// Constants associated with the SnapshotConfig.AccountType property.
// Type of account. Possible values are [enterprise, account].
const (
	SnapshotConfigAccountTypeAccountConst = "account"
	SnapshotConfigAccountTypeEnterpriseConst = "enterprise"
)

// Constants associated with the SnapshotConfig.Interval property.
// Frequency of taking the snapshot of the billing reports.
const (
	SnapshotConfigIntervalDailyConst = "daily"
)

// Constants associated with the SnapshotConfig.Versioning property.
// A new version of report is created or the existing report version is overwritten with every update.
const (
	SnapshotConfigVersioningNewConst = "new"
	SnapshotConfigVersioningOverwriteConst = "overwrite"
)

// Constants associated with the SnapshotConfig.ReportTypes property.
const (
	SnapshotConfigReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	SnapshotConfigReportTypesAccountSummaryConst = "account_summary"
	SnapshotConfigReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// UnmarshalSnapshotConfig unmarshals an instance of SnapshotConfig from the specified map of raw messages.
func UnmarshalSnapshotConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotConfig)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_type", &obj.AccountType)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		err = core.SDKErrorf(err, "", "interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "versioning", &obj.Versioning)
	if err != nil {
		err = core.SDKErrorf(err, "", "versioning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "report_types", &obj.ReportTypes)
	if err != nil {
		err = core.SDKErrorf(err, "", "report_types-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "compression", &obj.Compression)
	if err != nil {
		err = core.SDKErrorf(err, "", "compression-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		err = core.SDKErrorf(err, "", "content_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_reports_folder", &obj.CosReportsFolder)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_reports_folder-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_bucket", &obj.CosBucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_location", &obj.CosLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_endpoint", &obj.CosEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_updated_at", &obj.LastUpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "history", &obj.History, UnmarshalSnapshotConfigHistoryItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "history-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SnapshotConfigValidateResponse : Validated billing service to COS bucket authorization.
type SnapshotConfigValidateResponse struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket,omitempty"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location,omitempty"`
}

// UnmarshalSnapshotConfigValidateResponse unmarshals an instance of SnapshotConfigValidateResponse from the specified map of raw messages.
func UnmarshalSnapshotConfigValidateResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotConfigValidateResponse)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_bucket", &obj.CosBucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_location", &obj.CosLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_location-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Subscription : Subscription struct
type Subscription struct {
	// The ID of the subscription.
	SubscriptionID *string `json:"subscription_id" validate:"required"`

	// The charge agreement number of the subsciption.
	ChargeAgreementNumber *string `json:"charge_agreement_number" validate:"required"`

	// Type of the subscription.
	Type *string `json:"type" validate:"required"`

	// The credits available in the subscription for the month.
	SubscriptionAmount *float64 `json:"subscription_amount" validate:"required"`

	// The date from which the subscription was active.
	Start *strfmt.DateTime `json:"start" validate:"required"`

	// The date until which the subscription is active. End time is unavailable for PayGO accounts.
	End *strfmt.DateTime `json:"end,omitempty"`

	// The total credits available in the subscription.
	CreditsTotal *float64 `json:"credits_total" validate:"required"`

	// The terms through which the subscription is split into.
	Terms []SubscriptionTerm `json:"terms" validate:"required"`
}

// UnmarshalSubscription unmarshals an instance of Subscription from the specified map of raw messages.
func UnmarshalSubscription(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Subscription)
	err = core.UnmarshalPrimitive(m, "subscription_id", &obj.SubscriptionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "subscription_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "charge_agreement_number", &obj.ChargeAgreementNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "charge_agreement_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subscription_amount", &obj.SubscriptionAmount)
	if err != nil {
		err = core.SDKErrorf(err, "", "subscription_amount-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end", &obj.End)
	if err != nil {
		err = core.SDKErrorf(err, "", "end-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "credits_total", &obj.CreditsTotal)
	if err != nil {
		err = core.SDKErrorf(err, "", "credits_total-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "terms", &obj.Terms, UnmarshalSubscriptionTerm)
	if err != nil {
		err = core.SDKErrorf(err, "", "terms-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubscriptionSummary : A summary of charges and credits related to a subscription.
type SubscriptionSummary struct {
	// The charges after exhausting subscription credits and offers credits.
	Overage *float64 `json:"overage,omitempty"`

	// The list of subscriptions applicable for the month.
	Subscriptions []Subscription `json:"subscriptions,omitempty"`
}

// UnmarshalSubscriptionSummary unmarshals an instance of SubscriptionSummary from the specified map of raw messages.
func UnmarshalSubscriptionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubscriptionSummary)
	err = core.UnmarshalPrimitive(m, "overage", &obj.Overage)
	if err != nil {
		err = core.SDKErrorf(err, "", "overage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subscriptions", &obj.Subscriptions, UnmarshalSubscription)
	if err != nil {
		err = core.SDKErrorf(err, "", "subscriptions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubscriptionTerm : SubscriptionTerm struct
type SubscriptionTerm struct {
	// The start date of the term.
	Start *strfmt.DateTime `json:"start" validate:"required"`

	// The end date of the term.
	End *strfmt.DateTime `json:"end" validate:"required"`

	// Information about credits related to a subscription.
	Credits *SubscriptionTermCredits `json:"credits" validate:"required"`
}

// UnmarshalSubscriptionTerm unmarshals an instance of SubscriptionTerm from the specified map of raw messages.
func UnmarshalSubscriptionTerm(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubscriptionTerm)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end", &obj.End)
	if err != nil {
		err = core.SDKErrorf(err, "", "end-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "credits", &obj.Credits, UnmarshalSubscriptionTermCredits)
	if err != nil {
		err = core.SDKErrorf(err, "", "credits-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubscriptionTermCredits : Information about credits related to a subscription.
type SubscriptionTermCredits struct {
	// The total credits available for the term.
	Total *float64 `json:"total" validate:"required"`

	// The unused credits in the term at the beginning of the month.
	StartingBalance *float64 `json:"starting_balance" validate:"required"`

	// The credits used in this month.
	Used *float64 `json:"used" validate:"required"`

	// The remaining credits in this term.
	Balance *float64 `json:"balance" validate:"required"`
}

// UnmarshalSubscriptionTermCredits unmarshals an instance of SubscriptionTermCredits from the specified map of raw messages.
func UnmarshalSubscriptionTermCredits(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubscriptionTermCredits)
	err = core.UnmarshalPrimitive(m, "total", &obj.Total)
	if err != nil {
		err = core.SDKErrorf(err, "", "total-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "starting_balance", &obj.StartingBalance)
	if err != nil {
		err = core.SDKErrorf(err, "", "starting_balance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "used", &obj.Used)
	if err != nil {
		err = core.SDKErrorf(err, "", "used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "balance", &obj.Balance)
	if err != nil {
		err = core.SDKErrorf(err, "", "balance-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportSummary : SupportSummary struct
type SupportSummary struct {
	// The monthly support cost.
	Cost *float64 `json:"cost" validate:"required"`

	// The type of support.
	Type *string `json:"type" validate:"required"`

	// Additional support cost for the month.
	Overage *float64 `json:"overage" validate:"required"`
}

// UnmarshalSupportSummary unmarshals an instance of SupportSummary from the specified map of raw messages.
func UnmarshalSupportSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportSummary)
	err = core.UnmarshalPrimitive(m, "cost", &obj.Cost)
	if err != nil {
		err = core.SDKErrorf(err, "", "cost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "overage", &obj.Overage)
	if err != nil {
		err = core.SDKErrorf(err, "", "overage-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateReportsSnapshotConfigOptions : The UpdateReportsSnapshotConfig options.
type UpdateReportsSnapshotConfigOptions struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id" validate:"required"`

	// Frequency of taking the snapshot of the billing reports.
	Interval *string `json:"interval,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket,omitempty"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location,omitempty"`

	// The billing reports root folder to store the billing reports snapshots.
	CosReportsFolder *string `json:"cos_reports_folder,omitempty"`

	// The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes []string `json:"report_types,omitempty"`

	// A new version of report is created or the existing report version is overwritten with every update.
	Versioning *string `json:"versioning,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the UpdateReportsSnapshotConfigOptions.Interval property.
// Frequency of taking the snapshot of the billing reports.
const (
	UpdateReportsSnapshotConfigOptionsIntervalDailyConst = "daily"
)

// Constants associated with the UpdateReportsSnapshotConfigOptions.ReportTypes property.
const (
	UpdateReportsSnapshotConfigOptionsReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	UpdateReportsSnapshotConfigOptionsReportTypesAccountSummaryConst = "account_summary"
	UpdateReportsSnapshotConfigOptionsReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// Constants associated with the UpdateReportsSnapshotConfigOptions.Versioning property.
// A new version of report is created or the existing report version is overwritten with every update.
const (
	UpdateReportsSnapshotConfigOptionsVersioningNewConst = "new"
	UpdateReportsSnapshotConfigOptionsVersioningOverwriteConst = "overwrite"
)

// NewUpdateReportsSnapshotConfigOptions : Instantiate UpdateReportsSnapshotConfigOptions
func (*UsageReportsV4) NewUpdateReportsSnapshotConfigOptions(accountID string) *UpdateReportsSnapshotConfigOptions {
	return &UpdateReportsSnapshotConfigOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateReportsSnapshotConfigOptions) SetAccountID(accountID string) *UpdateReportsSnapshotConfigOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *UpdateReportsSnapshotConfigOptions) SetInterval(interval string) *UpdateReportsSnapshotConfigOptions {
	_options.Interval = core.StringPtr(interval)
	return _options
}

// SetCosBucket : Allow user to set CosBucket
func (_options *UpdateReportsSnapshotConfigOptions) SetCosBucket(cosBucket string) *UpdateReportsSnapshotConfigOptions {
	_options.CosBucket = core.StringPtr(cosBucket)
	return _options
}

// SetCosLocation : Allow user to set CosLocation
func (_options *UpdateReportsSnapshotConfigOptions) SetCosLocation(cosLocation string) *UpdateReportsSnapshotConfigOptions {
	_options.CosLocation = core.StringPtr(cosLocation)
	return _options
}

// SetCosReportsFolder : Allow user to set CosReportsFolder
func (_options *UpdateReportsSnapshotConfigOptions) SetCosReportsFolder(cosReportsFolder string) *UpdateReportsSnapshotConfigOptions {
	_options.CosReportsFolder = core.StringPtr(cosReportsFolder)
	return _options
}

// SetReportTypes : Allow user to set ReportTypes
func (_options *UpdateReportsSnapshotConfigOptions) SetReportTypes(reportTypes []string) *UpdateReportsSnapshotConfigOptions {
	_options.ReportTypes = reportTypes
	return _options
}

// SetVersioning : Allow user to set Versioning
func (_options *UpdateReportsSnapshotConfigOptions) SetVersioning(versioning string) *UpdateReportsSnapshotConfigOptions {
	_options.Versioning = core.StringPtr(versioning)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateReportsSnapshotConfigOptions) SetHeaders(param map[string]string) *UpdateReportsSnapshotConfigOptions {
	options.Headers = param
	return options
}

// ValidateReportsSnapshotConfigOptions : The ValidateReportsSnapshotConfig options.
type ValidateReportsSnapshotConfigOptions struct {
	// Account ID for which billing report snapshot is configured.
	AccountID *string `json:"account_id" validate:"required"`

	// Frequency of taking the snapshot of the billing reports.
	Interval *string `json:"interval,omitempty"`

	// The name of the COS bucket to store the snapshot of the billing reports.
	CosBucket *string `json:"cos_bucket,omitempty"`

	// Region of the COS instance.
	CosLocation *string `json:"cos_location,omitempty"`

	// The billing reports root folder to store the billing reports snapshots. Defaults to "IBMCloud-Billing-Reports".
	CosReportsFolder *string `json:"cos_reports_folder,omitempty"`

	// The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary,
	// account_resource_instance_usage].
	ReportTypes []string `json:"report_types,omitempty"`

	// A new version of report is created or the existing report version is overwritten with every update. Defaults to
	// "new".
	Versioning *string `json:"versioning,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ValidateReportsSnapshotConfigOptions.Interval property.
// Frequency of taking the snapshot of the billing reports.
const (
	ValidateReportsSnapshotConfigOptionsIntervalDailyConst = "daily"
)

// Constants associated with the ValidateReportsSnapshotConfigOptions.ReportTypes property.
const (
	ValidateReportsSnapshotConfigOptionsReportTypesAccountResourceInstanceUsageConst = "account_resource_instance_usage"
	ValidateReportsSnapshotConfigOptionsReportTypesAccountSummaryConst = "account_summary"
	ValidateReportsSnapshotConfigOptionsReportTypesEnterpriseSummaryConst = "enterprise_summary"
)

// Constants associated with the ValidateReportsSnapshotConfigOptions.Versioning property.
// A new version of report is created or the existing report version is overwritten with every update. Defaults to
// "new".
const (
	ValidateReportsSnapshotConfigOptionsVersioningNewConst = "new"
	ValidateReportsSnapshotConfigOptionsVersioningOverwriteConst = "overwrite"
)

// NewValidateReportsSnapshotConfigOptions : Instantiate ValidateReportsSnapshotConfigOptions
func (*UsageReportsV4) NewValidateReportsSnapshotConfigOptions(accountID string) *ValidateReportsSnapshotConfigOptions {
	return &ValidateReportsSnapshotConfigOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ValidateReportsSnapshotConfigOptions) SetAccountID(accountID string) *ValidateReportsSnapshotConfigOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetInterval : Allow user to set Interval
func (_options *ValidateReportsSnapshotConfigOptions) SetInterval(interval string) *ValidateReportsSnapshotConfigOptions {
	_options.Interval = core.StringPtr(interval)
	return _options
}

// SetCosBucket : Allow user to set CosBucket
func (_options *ValidateReportsSnapshotConfigOptions) SetCosBucket(cosBucket string) *ValidateReportsSnapshotConfigOptions {
	_options.CosBucket = core.StringPtr(cosBucket)
	return _options
}

// SetCosLocation : Allow user to set CosLocation
func (_options *ValidateReportsSnapshotConfigOptions) SetCosLocation(cosLocation string) *ValidateReportsSnapshotConfigOptions {
	_options.CosLocation = core.StringPtr(cosLocation)
	return _options
}

// SetCosReportsFolder : Allow user to set CosReportsFolder
func (_options *ValidateReportsSnapshotConfigOptions) SetCosReportsFolder(cosReportsFolder string) *ValidateReportsSnapshotConfigOptions {
	_options.CosReportsFolder = core.StringPtr(cosReportsFolder)
	return _options
}

// SetReportTypes : Allow user to set ReportTypes
func (_options *ValidateReportsSnapshotConfigOptions) SetReportTypes(reportTypes []string) *ValidateReportsSnapshotConfigOptions {
	_options.ReportTypes = reportTypes
	return _options
}

// SetVersioning : Allow user to set Versioning
func (_options *ValidateReportsSnapshotConfigOptions) SetVersioning(versioning string) *ValidateReportsSnapshotConfigOptions {
	_options.Versioning = core.StringPtr(versioning)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateReportsSnapshotConfigOptions) SetHeaders(param map[string]string) *ValidateReportsSnapshotConfigOptions {
	options.Headers = param
	return options
}

//
// GetResourceUsageAccountPager can be used to simplify the use of the "GetResourceUsageAccount" method.
//
type GetResourceUsageAccountPager struct {
	hasNext bool
	options *GetResourceUsageAccountOptions
	client  *UsageReportsV4
	pageContext struct {
		next *string
	}
}

// NewGetResourceUsageAccountPager returns a new GetResourceUsageAccountPager instance.
func (usageReports *UsageReportsV4) NewGetResourceUsageAccountPager(options *GetResourceUsageAccountOptions) (pager *GetResourceUsageAccountPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy GetResourceUsageAccountOptions = *options
	pager = &GetResourceUsageAccountPager{
		hasNext: true,
		options: &optionsCopy,
		client:  usageReports,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetResourceUsageAccountPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetResourceUsageAccountPager) GetNextWithContext(ctx context.Context) (page []InstanceUsage, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetResourceUsageAccountWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var _start *string
		_start, err = core.GetQueryParam(result.Next.Href, "_start")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving '_start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = _start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetResourceUsageAccountPager) GetAllWithContext(ctx context.Context) (allItems []InstanceUsage, err error) {
	for pager.HasNext() {
		var nextPage []InstanceUsage
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
func (pager *GetResourceUsageAccountPager) GetNext() (page []InstanceUsage, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetResourceUsageAccountPager) GetAll() (allItems []InstanceUsage, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// GetResourceUsageResourceGroupPager can be used to simplify the use of the "GetResourceUsageResourceGroup" method.
//
type GetResourceUsageResourceGroupPager struct {
	hasNext bool
	options *GetResourceUsageResourceGroupOptions
	client  *UsageReportsV4
	pageContext struct {
		next *string
	}
}

// NewGetResourceUsageResourceGroupPager returns a new GetResourceUsageResourceGroupPager instance.
func (usageReports *UsageReportsV4) NewGetResourceUsageResourceGroupPager(options *GetResourceUsageResourceGroupOptions) (pager *GetResourceUsageResourceGroupPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy GetResourceUsageResourceGroupOptions = *options
	pager = &GetResourceUsageResourceGroupPager{
		hasNext: true,
		options: &optionsCopy,
		client:  usageReports,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetResourceUsageResourceGroupPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetResourceUsageResourceGroupPager) GetNextWithContext(ctx context.Context) (page []InstanceUsage, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetResourceUsageResourceGroupWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var _start *string
		_start, err = core.GetQueryParam(result.Next.Href, "_start")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving '_start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = _start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetResourceUsageResourceGroupPager) GetAllWithContext(ctx context.Context) (allItems []InstanceUsage, err error) {
	for pager.HasNext() {
		var nextPage []InstanceUsage
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
func (pager *GetResourceUsageResourceGroupPager) GetNext() (page []InstanceUsage, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetResourceUsageResourceGroupPager) GetAll() (allItems []InstanceUsage, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// GetResourceUsageOrgPager can be used to simplify the use of the "GetResourceUsageOrg" method.
//
type GetResourceUsageOrgPager struct {
	hasNext bool
	options *GetResourceUsageOrgOptions
	client  *UsageReportsV4
	pageContext struct {
		next *string
	}
}

// NewGetResourceUsageOrgPager returns a new GetResourceUsageOrgPager instance.
func (usageReports *UsageReportsV4) NewGetResourceUsageOrgPager(options *GetResourceUsageOrgOptions) (pager *GetResourceUsageOrgPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy GetResourceUsageOrgOptions = *options
	pager = &GetResourceUsageOrgPager{
		hasNext: true,
		options: &optionsCopy,
		client:  usageReports,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetResourceUsageOrgPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetResourceUsageOrgPager) GetNextWithContext(ctx context.Context) (page []InstanceUsage, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetResourceUsageOrgWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var _start *string
		_start, err = core.GetQueryParam(result.Next.Href, "_start")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving '_start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = _start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetResourceUsageOrgPager) GetAllWithContext(ctx context.Context) (allItems []InstanceUsage, err error) {
	for pager.HasNext() {
		var nextPage []InstanceUsage
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
func (pager *GetResourceUsageOrgPager) GetNext() (page []InstanceUsage, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetResourceUsageOrgPager) GetAll() (allItems []InstanceUsage, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// GetReportsSnapshotPager can be used to simplify the use of the "GetReportsSnapshot" method.
//
type GetReportsSnapshotPager struct {
	hasNext bool
	options *GetReportsSnapshotOptions
	client  *UsageReportsV4
	pageContext struct {
		next *string
	}
}

// NewGetReportsSnapshotPager returns a new GetReportsSnapshotPager instance.
func (usageReports *UsageReportsV4) NewGetReportsSnapshotPager(options *GetReportsSnapshotOptions) (pager *GetReportsSnapshotPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy GetReportsSnapshotOptions = *options
	pager = &GetReportsSnapshotPager{
		hasNext: true,
		options: &optionsCopy,
		client:  usageReports,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *GetReportsSnapshotPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *GetReportsSnapshotPager) GetNextWithContext(ctx context.Context) (page []SnapshotListSnapshotsItem, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.GetReportsSnapshotWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var _start *string
		_start, err = core.GetQueryParam(result.Next.Href, "_start")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving '_start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = _start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Snapshots

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *GetReportsSnapshotPager) GetAllWithContext(ctx context.Context) (allItems []SnapshotListSnapshotsItem, err error) {
	for pager.HasNext() {
		var nextPage []SnapshotListSnapshotsItem
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
func (pager *GetReportsSnapshotPager) GetNext() (page []SnapshotListSnapshotsItem, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *GetReportsSnapshotPager) GetAll() (allItems []SnapshotListSnapshotsItem, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
