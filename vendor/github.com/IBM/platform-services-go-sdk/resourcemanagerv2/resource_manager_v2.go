/**
 * (C) Copyright IBM Corp. 2021.
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
 * IBM OpenAPI SDK Code Generator Version: 3.41.0-f1ef0102-20211018-193503
 */

// Package resourcemanagerv2 : Operations and models for the ResourceManagerV2 service
package resourcemanagerv2

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

// ResourceManagerV2 : Manage lifecycle of your Cloud resource groups using Resource Manager APIs.
//
// API Version: 2.0
type ResourceManagerV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://resource-controller.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "resource_manager"

// ResourceManagerV2Options : Service options
type ResourceManagerV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewResourceManagerV2UsingExternalConfig : constructs an instance of ResourceManagerV2 with passed in options and external configuration.
func NewResourceManagerV2UsingExternalConfig(options *ResourceManagerV2Options) (resourceManager *ResourceManagerV2, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	resourceManager, err = NewResourceManagerV2(options)
	if err != nil {
		return
	}

	err = resourceManager.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = resourceManager.Service.SetServiceURL(options.URL)
	}
	return
}

// NewResourceManagerV2 : constructs an instance of ResourceManagerV2 with passed in options.
func NewResourceManagerV2(options *ResourceManagerV2Options) (service *ResourceManagerV2, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			return
		}
	}

	service = &ResourceManagerV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "resourceManager" suitable for processing requests.
func (resourceManager *ResourceManagerV2) Clone() *ResourceManagerV2 {
	if core.IsNil(resourceManager) {
		return nil
	}
	clone := *resourceManager
	clone.Service = resourceManager.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (resourceManager *ResourceManagerV2) SetServiceURL(url string) error {
	return resourceManager.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (resourceManager *ResourceManagerV2) GetServiceURL() string {
	return resourceManager.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (resourceManager *ResourceManagerV2) SetDefaultHeaders(headers http.Header) {
	resourceManager.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (resourceManager *ResourceManagerV2) SetEnableGzipCompression(enableGzip bool) {
	resourceManager.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (resourceManager *ResourceManagerV2) GetEnableGzipCompression() bool {
	return resourceManager.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (resourceManager *ResourceManagerV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	resourceManager.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (resourceManager *ResourceManagerV2) DisableRetries() {
	resourceManager.Service.DisableRetries()
}

// ListResourceGroups : Get a list of all resource groups
// Call this method to retrieve information about all resource groups and associated quotas in an account. The `id`
// returned in the response can be used to [create a resource instance
// later](https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=java#create-resource-instance). The
// response can be filtered based on queryParams such as `account_id`, `name`, `default`, and more to narrow your
// search.Users need to be assigned IAM policies with the Viewer role or higher on the targeted resource groups.
func (resourceManager *ResourceManagerV2) ListResourceGroups(listResourceGroupsOptions *ListResourceGroupsOptions) (result *ResourceGroupList, response *core.DetailedResponse, err error) {
	return resourceManager.ListResourceGroupsWithContext(context.Background(), listResourceGroupsOptions)
}

// ListResourceGroupsWithContext is an alternate form of the ListResourceGroups method which supports a Context parameter
func (resourceManager *ResourceManagerV2) ListResourceGroupsWithContext(ctx context.Context, listResourceGroupsOptions *ListResourceGroupsOptions) (result *ResourceGroupList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listResourceGroupsOptions, "listResourceGroupsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/resource_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listResourceGroupsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "ListResourceGroups")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listResourceGroupsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listResourceGroupsOptions.AccountID))
	}
	if listResourceGroupsOptions.Date != nil {
		builder.AddQuery("date", fmt.Sprint(*listResourceGroupsOptions.Date))
	}
	if listResourceGroupsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listResourceGroupsOptions.Name))
	}
	if listResourceGroupsOptions.Default != nil {
		builder.AddQuery("default", fmt.Sprint(*listResourceGroupsOptions.Default))
	}
	if listResourceGroupsOptions.IncludeDeleted != nil {
		builder.AddQuery("include_deleted", fmt.Sprint(*listResourceGroupsOptions.IncludeDeleted))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroupList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateResourceGroup : Create a resource group
// Create a resource group in an account to organize your account resources in customizable groupings so that you can
// quickly assign users access to more than one resource at a time. To learn what makes a good resource group strategy,
// see [Best practices for organizing resources](https://cloud.ibm.com/docs/account?topic=account-account_setup). A
// default resource group is created when an account is created. If you have a Lite account or 30-day trial, you cannot
// create extra resource groups, but you can rename your default resource group. If you have a Pay-As-You-Go or
// Subscription account, you can create multiple resource groups. You must be assigned an IAM policy with the
// Administrator role on All Account Management services to create extra resource groups.
func (resourceManager *ResourceManagerV2) CreateResourceGroup(createResourceGroupOptions *CreateResourceGroupOptions) (result *ResCreateResourceGroup, response *core.DetailedResponse, err error) {
	return resourceManager.CreateResourceGroupWithContext(context.Background(), createResourceGroupOptions)
}

// CreateResourceGroupWithContext is an alternate form of the CreateResourceGroup method which supports a Context parameter
func (resourceManager *ResourceManagerV2) CreateResourceGroupWithContext(ctx context.Context, createResourceGroupOptions *CreateResourceGroupOptions) (result *ResCreateResourceGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createResourceGroupOptions, "createResourceGroupOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/resource_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "CreateResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createResourceGroupOptions.Name != nil {
		body["name"] = createResourceGroupOptions.Name
	}
	if createResourceGroupOptions.AccountID != nil {
		body["account_id"] = createResourceGroupOptions.AccountID
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResCreateResourceGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetResourceGroup : Get a resource group
// Retrieve a resource group by alias ID. Call this method to get details about a particular resource group, like the
// name of the resource group, associated quotas, whether the state is active, the resource group ID and the CRN. The
// `id` returned in the response can be used to [create a resource instance
// later](https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=java#create-resource-instance).
// Users need to be assigned an IAM policy with the Viewer role or higher on the targeted resource group.
func (resourceManager *ResourceManagerV2) GetResourceGroup(getResourceGroupOptions *GetResourceGroupOptions) (result *ResourceGroup, response *core.DetailedResponse, err error) {
	return resourceManager.GetResourceGroupWithContext(context.Background(), getResourceGroupOptions)
}

// GetResourceGroupWithContext is an alternate form of the GetResourceGroup method which supports a Context parameter
func (resourceManager *ResourceManagerV2) GetResourceGroupWithContext(ctx context.Context, getResourceGroupOptions *GetResourceGroupOptions) (result *ResourceGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceGroupOptions, "getResourceGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getResourceGroupOptions, "getResourceGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getResourceGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/resource_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "GetResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateResourceGroup : Update a resource group
// Update a resource group by the alias ID. Call this method to update information about an existing resource group. You
// can rename a resource group and activate or suspend a particular resource group. To update a resource group, users
// need to be assigned with IAM policies with the Editor role or higher.
func (resourceManager *ResourceManagerV2) UpdateResourceGroup(updateResourceGroupOptions *UpdateResourceGroupOptions) (result *ResourceGroup, response *core.DetailedResponse, err error) {
	return resourceManager.UpdateResourceGroupWithContext(context.Background(), updateResourceGroupOptions)
}

// UpdateResourceGroupWithContext is an alternate form of the UpdateResourceGroup method which supports a Context parameter
func (resourceManager *ResourceManagerV2) UpdateResourceGroupWithContext(ctx context.Context, updateResourceGroupOptions *UpdateResourceGroupOptions) (result *ResourceGroup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateResourceGroupOptions, "updateResourceGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateResourceGroupOptions, "updateResourceGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateResourceGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/resource_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "UpdateResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateResourceGroupOptions.Name != nil {
		body["name"] = updateResourceGroupOptions.Name
	}
	if updateResourceGroupOptions.State != nil {
		body["state"] = updateResourceGroupOptions.State
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroup)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteResourceGroup : Delete a resource group
// Delete a resource group by the alias ID. You can delete a resource group only if the targeted resource group does not
// contain any resources or if it is not a default resource group. When a user creates an account, a default resource
// group is created in the account. If you want to delete a resource group that contains resources, first [delete the
// resource
// instances](https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=java#delete-resource-instance).
// Then, delete the resource group when all resource instances in the group are deleted. Users need to be assigned an
// IAM policy with the Editor role or higher on the targeted resource group.
func (resourceManager *ResourceManagerV2) DeleteResourceGroup(deleteResourceGroupOptions *DeleteResourceGroupOptions) (response *core.DetailedResponse, err error) {
	return resourceManager.DeleteResourceGroupWithContext(context.Background(), deleteResourceGroupOptions)
}

// DeleteResourceGroupWithContext is an alternate form of the DeleteResourceGroup method which supports a Context parameter
func (resourceManager *ResourceManagerV2) DeleteResourceGroupWithContext(ctx context.Context, deleteResourceGroupOptions *DeleteResourceGroupOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourceGroupOptions, "deleteResourceGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteResourceGroupOptions, "deleteResourceGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteResourceGroupOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/resource_groups/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "DeleteResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = resourceManager.Service.Request(request, nil)

	return
}

// ListQuotaDefinitions : List quota definitions
// Get a list of all quota definitions. Quotas for a resource group limit the number of apps, instances, and memory
// allowed for that specific resource group. Each resource group that you have on your account has a specific set of
// quotas. Standard quotas are for resource groups that are created by users with a Lite account, and Pay-As-You-Go
// quotas are for resource groups that are created with a Pay-As-You-Go account. This method provides list of all
// available quota definitions. No specific IAM policy needed.
func (resourceManager *ResourceManagerV2) ListQuotaDefinitions(listQuotaDefinitionsOptions *ListQuotaDefinitionsOptions) (result *QuotaDefinitionList, response *core.DetailedResponse, err error) {
	return resourceManager.ListQuotaDefinitionsWithContext(context.Background(), listQuotaDefinitionsOptions)
}

// ListQuotaDefinitionsWithContext is an alternate form of the ListQuotaDefinitions method which supports a Context parameter
func (resourceManager *ResourceManagerV2) ListQuotaDefinitionsWithContext(ctx context.Context, listQuotaDefinitionsOptions *ListQuotaDefinitionsOptions) (result *QuotaDefinitionList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listQuotaDefinitionsOptions, "listQuotaDefinitionsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/quota_definitions`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listQuotaDefinitionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "ListQuotaDefinitions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQuotaDefinitionList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetQuotaDefinition : Get a quota definition
// Call this method to retrieve information about a particular quota by passing the quota ID. The response can be used
// to identify the quota type, Standard or Paid. Information about available resources, such as number of apps, number
// of service instances, and memory, are returned in the response. Quotas for a resource group limit the number of apps,
// instances, and memory allowed for that specific resource group. Each resource group that you have on your account has
// a specific set of quotas. Standard quotas are for resource groups that are created by users with a Lite account, and
// Pay-As-You-Go quotas are for resource groups that are created with a Pay-As-You-Go account. No specific IAM policy
// needed.
func (resourceManager *ResourceManagerV2) GetQuotaDefinition(getQuotaDefinitionOptions *GetQuotaDefinitionOptions) (result *QuotaDefinition, response *core.DetailedResponse, err error) {
	return resourceManager.GetQuotaDefinitionWithContext(context.Background(), getQuotaDefinitionOptions)
}

// GetQuotaDefinitionWithContext is an alternate form of the GetQuotaDefinition method which supports a Context parameter
func (resourceManager *ResourceManagerV2) GetQuotaDefinitionWithContext(ctx context.Context, getQuotaDefinitionOptions *GetQuotaDefinitionOptions) (result *QuotaDefinition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getQuotaDefinitionOptions, "getQuotaDefinitionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getQuotaDefinitionOptions, "getQuotaDefinitionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getQuotaDefinitionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceManager.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceManager.Service.Options.URL, `/v2/quota_definitions/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getQuotaDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_manager", "V2", "GetQuotaDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceManager.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQuotaDefinition)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateResourceGroupOptions : The CreateResourceGroup options.
type CreateResourceGroupOptions struct {
	// The new name of the resource group.
	Name *string `json:"name,omitempty"`

	// The account id of the resource group.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateResourceGroupOptions : Instantiate CreateResourceGroupOptions
func (*ResourceManagerV2) NewCreateResourceGroupOptions() *CreateResourceGroupOptions {
	return &CreateResourceGroupOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateResourceGroupOptions) SetName(name string) *CreateResourceGroupOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateResourceGroupOptions) SetAccountID(accountID string) *CreateResourceGroupOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateResourceGroupOptions) SetHeaders(param map[string]string) *CreateResourceGroupOptions {
	options.Headers = param
	return options
}

// DeleteResourceGroupOptions : The DeleteResourceGroup options.
type DeleteResourceGroupOptions struct {
	// The short or long ID of the alias.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteResourceGroupOptions : Instantiate DeleteResourceGroupOptions
func (*ResourceManagerV2) NewDeleteResourceGroupOptions(id string) *DeleteResourceGroupOptions {
	return &DeleteResourceGroupOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteResourceGroupOptions) SetID(id string) *DeleteResourceGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteResourceGroupOptions) SetHeaders(param map[string]string) *DeleteResourceGroupOptions {
	options.Headers = param
	return options
}

// GetQuotaDefinitionOptions : The GetQuotaDefinition options.
type GetQuotaDefinitionOptions struct {
	// The id of the quota.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetQuotaDefinitionOptions : Instantiate GetQuotaDefinitionOptions
func (*ResourceManagerV2) NewGetQuotaDefinitionOptions(id string) *GetQuotaDefinitionOptions {
	return &GetQuotaDefinitionOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetQuotaDefinitionOptions) SetID(id string) *GetQuotaDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetQuotaDefinitionOptions) SetHeaders(param map[string]string) *GetQuotaDefinitionOptions {
	options.Headers = param
	return options
}

// GetResourceGroupOptions : The GetResourceGroup options.
type GetResourceGroupOptions struct {
	// The short or long ID of the alias.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetResourceGroupOptions : Instantiate GetResourceGroupOptions
func (*ResourceManagerV2) NewGetResourceGroupOptions(id string) *GetResourceGroupOptions {
	return &GetResourceGroupOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetResourceGroupOptions) SetID(id string) *GetResourceGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceGroupOptions) SetHeaders(param map[string]string) *GetResourceGroupOptions {
	options.Headers = param
	return options
}

// ListQuotaDefinitionsOptions : The ListQuotaDefinitions options.
type ListQuotaDefinitionsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListQuotaDefinitionsOptions : Instantiate ListQuotaDefinitionsOptions
func (*ResourceManagerV2) NewListQuotaDefinitionsOptions() *ListQuotaDefinitionsOptions {
	return &ListQuotaDefinitionsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListQuotaDefinitionsOptions) SetHeaders(param map[string]string) *ListQuotaDefinitionsOptions {
	options.Headers = param
	return options
}

// ListResourceGroupsOptions : The ListResourceGroups options.
type ListResourceGroupsOptions struct {
	// The ID of the account that contains the resource groups that you want to get.
	AccountID *string `json:"account_id,omitempty"`

	// The date in the format of YYYY-MM which returns resource groups. Deleted resource groups will be excluded before
	// this month.
	Date *string `json:"date,omitempty"`

	// The name of the resource group.
	Name *string `json:"name,omitempty"`

	// Boolean value to specify whether or not to list default resource groups.
	Default *bool `json:"default,omitempty"`

	// Boolean value to specify whether or not to list default resource groups.
	IncludeDeleted *bool `json:"include_deleted,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListResourceGroupsOptions : Instantiate ListResourceGroupsOptions
func (*ResourceManagerV2) NewListResourceGroupsOptions() *ListResourceGroupsOptions {
	return &ListResourceGroupsOptions{}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListResourceGroupsOptions) SetAccountID(accountID string) *ListResourceGroupsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetDate : Allow user to set Date
func (_options *ListResourceGroupsOptions) SetDate(date string) *ListResourceGroupsOptions {
	_options.Date = core.StringPtr(date)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListResourceGroupsOptions) SetName(name string) *ListResourceGroupsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDefault : Allow user to set Default
func (_options *ListResourceGroupsOptions) SetDefault(defaultVar bool) *ListResourceGroupsOptions {
	_options.Default = core.BoolPtr(defaultVar)
	return _options
}

// SetIncludeDeleted : Allow user to set IncludeDeleted
func (_options *ListResourceGroupsOptions) SetIncludeDeleted(includeDeleted bool) *ListResourceGroupsOptions {
	_options.IncludeDeleted = core.BoolPtr(includeDeleted)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListResourceGroupsOptions) SetHeaders(param map[string]string) *ListResourceGroupsOptions {
	options.Headers = param
	return options
}

// QuotaDefinition : A returned quota definition.
type QuotaDefinition struct {
	// An alpha-numeric value identifying the quota.
	ID *string `json:"id,omitempty"`

	// The human-readable name of the quota.
	Name *string `json:"name,omitempty"`

	// The type of the quota.
	Type *string `json:"type,omitempty"`

	// The total app limit.
	NumberOfApps *float64 `json:"number_of_apps,omitempty"`

	// The total service instances limit per app.
	NumberOfServiceInstances *float64 `json:"number_of_service_instances,omitempty"`

	// Default number of instances per lite plan.
	DefaultNumberOfInstancesPerLitePlan *float64 `json:"default_number_of_instances_per_lite_plan,omitempty"`

	// The total instances limit per app.
	InstancesPerApp *float64 `json:"instances_per_app,omitempty"`

	// The total memory of app instance.
	InstanceMemory *string `json:"instance_memory,omitempty"`

	// The total app memory capacity.
	TotalAppMemory *string `json:"total_app_memory,omitempty"`

	// The VSI limit.
	VsiLimit *float64 `json:"vsi_limit,omitempty"`

	// The resource quotas associated with a quota definition.
	ResourceQuotas []ResourceQuota `json:"resource_quotas,omitempty"`

	// The date when the quota was initially created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The date when the quota was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalQuotaDefinition unmarshals an instance of QuotaDefinition from the specified map of raw messages.
func UnmarshalQuotaDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QuotaDefinition)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "number_of_apps", &obj.NumberOfApps)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "number_of_service_instances", &obj.NumberOfServiceInstances)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_number_of_instances_per_lite_plan", &obj.DefaultNumberOfInstancesPerLitePlan)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instances_per_app", &obj.InstancesPerApp)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_memory", &obj.InstanceMemory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_app_memory", &obj.TotalAppMemory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "vsi_limit", &obj.VsiLimit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_quotas", &obj.ResourceQuotas, UnmarshalResourceQuota)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QuotaDefinitionList : A list of quota definitions.
type QuotaDefinitionList struct {
	// The list of quota definitions.
	Resources []QuotaDefinition `json:"resources" validate:"required"`
}

// UnmarshalQuotaDefinitionList unmarshals an instance of QuotaDefinitionList from the specified map of raw messages.
func UnmarshalQuotaDefinitionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QuotaDefinitionList)
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalQuotaDefinition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResCreateResourceGroup : A newly-created resource group.
type ResCreateResourceGroup struct {
	// An alpha-numeric value identifying the resource group.
	ID *string `json:"id,omitempty"`

	// The full CRN (cloud resource name) associated with the resource group. For more on this format, see [Cloud Resource
	// Names](https://cloud.ibm.com/docs/account?topic=account-crn).
	CRN *string `json:"crn,omitempty"`
}

// UnmarshalResCreateResourceGroup unmarshals an instance of ResCreateResourceGroup from the specified map of raw messages.
func UnmarshalResCreateResourceGroup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResCreateResourceGroup)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceGroup : A resource group.
type ResourceGroup struct {
	// An alpha-numeric value identifying the resource group.
	ID *string `json:"id,omitempty"`

	// The full CRN (cloud resource name) associated with the resource group. For more on this format, see [Cloud Resource
	// Names](https://cloud.ibm.com/docs/account?topic=account-crn).
	CRN *string `json:"crn,omitempty"`

	// An alpha-numeric value identifying the account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The human-readable name of the resource group.
	Name *string `json:"name,omitempty"`

	// The state of the resource group.
	State *string `json:"state,omitempty"`

	// Identify if this resource group is default of the account or not.
	Default *bool `json:"default,omitempty"`

	// An alpha-numeric value identifying the quota ID associated with the resource group.
	QuotaID *string `json:"quota_id,omitempty"`

	// The URL to access the quota details that associated with the resource group.
	QuotaURL *string `json:"quota_url,omitempty"`

	// The URL to access the payment methods details that associated with the resource group.
	PaymentMethodsURL *string `json:"payment_methods_url,omitempty"`

	// An array of the resources that linked to the resource group.
	ResourceLinkages []interface{} `json:"resource_linkages,omitempty"`

	// The URL to access the team details that associated with the resource group.
	TeamsURL *string `json:"teams_url,omitempty"`

	// The date when the resource group was initially created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The date when the resource group was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalResourceGroup unmarshals an instance of ResourceGroup from the specified map of raw messages.
func UnmarshalResourceGroup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroup)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "quota_id", &obj.QuotaID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "quota_url", &obj.QuotaURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "payment_methods_url", &obj.PaymentMethodsURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_linkages", &obj.ResourceLinkages)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "teams_url", &obj.TeamsURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceGroupList : A list of resource groups.
type ResourceGroupList struct {
	// The list of resource groups.
	Resources []ResourceGroup `json:"resources" validate:"required"`
}

// UnmarshalResourceGroupList unmarshals an instance of ResourceGroupList from the specified map of raw messages.
func UnmarshalResourceGroupList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupList)
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResourceGroup)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQuota : A resource quota.
type ResourceQuota struct {
	// An alpha-numeric value identifying the quota.
	ID *string `json:"_id,omitempty"`

	// The human-readable name of the quota.
	ResourceID *string `json:"resource_id,omitempty"`

	// The full CRN (cloud resource name) associated with the quota. For more on this format, see
	// https://cloud.ibm.com/docs/account?topic=account-crn.
	CRN *string `json:"crn,omitempty"`

	// The limit number of this resource.
	Limit *float64 `json:"limit,omitempty"`
}

// UnmarshalResourceQuota unmarshals an instance of ResourceQuota from the specified map of raw messages.
func UnmarshalResourceQuota(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQuota)
	err = core.UnmarshalPrimitive(m, "_id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_id", &obj.ResourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateResourceGroupOptions : The UpdateResourceGroup options.
type UpdateResourceGroupOptions struct {
	// The short or long ID of the alias.
	ID *string `json:"id" validate:"required,ne="`

	// The new name of the resource group.
	Name *string `json:"name,omitempty"`

	// The state of the resource group.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateResourceGroupOptions : Instantiate UpdateResourceGroupOptions
func (*ResourceManagerV2) NewUpdateResourceGroupOptions(id string) *UpdateResourceGroupOptions {
	return &UpdateResourceGroupOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateResourceGroupOptions) SetID(id string) *UpdateResourceGroupOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateResourceGroupOptions) SetName(name string) *UpdateResourceGroupOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetState : Allow user to set State
func (_options *UpdateResourceGroupOptions) SetState(state string) *UpdateResourceGroupOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateResourceGroupOptions) SetHeaders(param map[string]string) *UpdateResourceGroupOptions {
	options.Headers = param
	return options
}
