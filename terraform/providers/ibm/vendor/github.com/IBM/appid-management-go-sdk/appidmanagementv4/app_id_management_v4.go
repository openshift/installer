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
 * IBM OpenAPI SDK Code Generator Version: 3.32.0-4c6a3129-20210514-210323
 */

// Package appidmanagementv4 : Operations and models for the AppIDManagementV4 service
package appidmanagementv4

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/appid-management-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// AppIDManagementV4 : You can use the following APIs to configure your instances of IBM Cloud App ID. To define fine
// grain access policies, you must have an instance of App ID that was created after March 15, 2018.</br> New to the
// APIs? Try them out by using the <a href="https://github.com/ibm-cloud-security/appid-postman">App ID Postman
// collection</a>!</br> </br> <b>Important:</b> You must have an <a
// href="https://cloud.ibm.com/docs/account?topic=account-iamoverview">IBM Cloud Identity and Access Management</a>
// token to access the APIs. For help obtaining a token, check out <a
// href="https://cloud.ibm.com/docs/account?topic=account-iamtoken_from_apikey">Getting an IAM token with an API
// key</a>.
//
// Version: 4
type AppIDManagementV4 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://app-id-management.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "app_id_management"

// AppIDManagementV4Options : Service options
type AppIDManagementV4Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewAppIDManagementV4UsingExternalConfig : constructs an instance of AppIDManagementV4 with passed in options and external configuration.
func NewAppIDManagementV4UsingExternalConfig(options *AppIDManagementV4Options) (appIdManagement *AppIDManagementV4, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	appIdManagement, err = NewAppIDManagementV4(options)
	if err != nil {
		return
	}

	err = appIdManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = appIdManagement.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAppIDManagementV4 : constructs an instance of AppIDManagementV4 with passed in options.
func NewAppIDManagementV4(options *AppIDManagementV4Options) (service *AppIDManagementV4, err error) {
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

	service = &AppIDManagementV4{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "appIdManagement" suitable for processing requests.
func (appIdManagement *AppIDManagementV4) Clone() *AppIDManagementV4 {
	if core.IsNil(appIdManagement) {
		return nil
	}
	clone := *appIdManagement
	clone.Service = appIdManagement.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (appIdManagement *AppIDManagementV4) SetServiceURL(url string) error {
	return appIdManagement.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (appIdManagement *AppIDManagementV4) GetServiceURL() string {
	return appIdManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (appIdManagement *AppIDManagementV4) SetDefaultHeaders(headers http.Header) {
	appIdManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (appIdManagement *AppIDManagementV4) SetEnableGzipCompression(enableGzip bool) {
	appIdManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (appIdManagement *AppIDManagementV4) GetEnableGzipCompression() bool {
	return appIdManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (appIdManagement *AppIDManagementV4) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	appIdManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (appIdManagement *AppIDManagementV4) DisableRetries() {
	appIdManagement.Service.DisableRetries()
}

// ListApplications : List applications
// Returns a list of all applications registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) ListApplications(listApplicationsOptions *ListApplicationsOptions) (result *ApplicationsList, response *core.DetailedResponse, err error) {
	return appIdManagement.ListApplicationsWithContext(context.Background(), listApplicationsOptions)
}

// ListApplicationsWithContext is an alternate form of the ListApplications method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ListApplicationsWithContext(ctx context.Context, listApplicationsOptions *ListApplicationsOptions) (result *ApplicationsList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listApplicationsOptions, "listApplicationsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listApplicationsOptions, "listApplicationsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *listApplicationsOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listApplicationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ListApplications")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplicationsList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// RegisterApplication : Create application
// Register a new application with the App ID instance.
func (appIdManagement *AppIDManagementV4) RegisterApplication(registerApplicationOptions *RegisterApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	return appIdManagement.RegisterApplicationWithContext(context.Background(), registerApplicationOptions)
}

// RegisterApplicationWithContext is an alternate form of the RegisterApplication method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) RegisterApplicationWithContext(ctx context.Context, registerApplicationOptions *RegisterApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(registerApplicationOptions, "registerApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(registerApplicationOptions, "registerApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *registerApplicationOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range registerApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "RegisterApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if registerApplicationOptions.Name != nil {
		body["name"] = registerApplicationOptions.Name
	}
	if registerApplicationOptions.Type != nil {
		body["type"] = registerApplicationOptions.Type
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplication)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetApplication : Get application
// Returns a specific application registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) GetApplication(getApplicationOptions *GetApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	return appIdManagement.GetApplicationWithContext(context.Background(), getApplicationOptions)
}

// GetApplicationWithContext is an alternate form of the GetApplication method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetApplicationWithContext(ctx context.Context, getApplicationOptions *GetApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getApplicationOptions, "getApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getApplicationOptions, "getApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getApplicationOptions.TenantID,
		"clientId": *getApplicationOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplication)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateApplication : Update application
// Update an application registered with the App ID instance.
func (appIdManagement *AppIDManagementV4) UpdateApplication(updateApplicationOptions *UpdateApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateApplicationWithContext(context.Background(), updateApplicationOptions)
}

// UpdateApplicationWithContext is an alternate form of the UpdateApplication method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateApplicationWithContext(ctx context.Context, updateApplicationOptions *UpdateApplicationOptions) (result *Application, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateApplicationOptions, "updateApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateApplicationOptions, "updateApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateApplicationOptions.TenantID,
		"clientId": *updateApplicationOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateApplicationOptions.Name != nil {
		body["name"] = updateApplicationOptions.Name
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplication)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteApplication : Delete application
// Delete an application registered with the App ID instance. Note: This action cannot be undone.
func (appIdManagement *AppIDManagementV4) DeleteApplication(deleteApplicationOptions *DeleteApplicationOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.DeleteApplicationWithContext(context.Background(), deleteApplicationOptions)
}

// DeleteApplicationWithContext is an alternate form of the DeleteApplication method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) DeleteApplicationWithContext(ctx context.Context, deleteApplicationOptions *DeleteApplicationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteApplicationOptions, "deleteApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteApplicationOptions, "deleteApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *deleteApplicationOptions.TenantID,
		"clientId": *deleteApplicationOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "DeleteApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetApplicationScopes : Get application scopes
// View the defined scopes for an application that is registered with an App ID instance.
func (appIdManagement *AppIDManagementV4) GetApplicationScopes(getApplicationScopesOptions *GetApplicationScopesOptions) (result *GetScopesForApplication, response *core.DetailedResponse, err error) {
	return appIdManagement.GetApplicationScopesWithContext(context.Background(), getApplicationScopesOptions)
}

// GetApplicationScopesWithContext is an alternate form of the GetApplicationScopes method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetApplicationScopesWithContext(ctx context.Context, getApplicationScopesOptions *GetApplicationScopesOptions) (result *GetScopesForApplication, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getApplicationScopesOptions, "getApplicationScopesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getApplicationScopesOptions, "getApplicationScopesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getApplicationScopesOptions.TenantID,
		"clientId": *getApplicationScopesOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}/scopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getApplicationScopesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetApplicationScopes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetScopesForApplication)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutApplicationsScopes : Add application scope
// Update the scopes for a registered application.</br> <b>Important</b>: Removing a scope from an array deletes it from
// any roles that it is associated with and the action cannot be undone.
func (appIdManagement *AppIDManagementV4) PutApplicationsScopes(putApplicationsScopesOptions *PutApplicationsScopesOptions) (result *GetScopesForApplication, response *core.DetailedResponse, err error) {
	return appIdManagement.PutApplicationsScopesWithContext(context.Background(), putApplicationsScopesOptions)
}

// PutApplicationsScopesWithContext is an alternate form of the PutApplicationsScopes method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PutApplicationsScopesWithContext(ctx context.Context, putApplicationsScopesOptions *PutApplicationsScopesOptions) (result *GetScopesForApplication, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putApplicationsScopesOptions, "putApplicationsScopesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putApplicationsScopesOptions, "putApplicationsScopesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *putApplicationsScopesOptions.TenantID,
		"clientId": *putApplicationsScopesOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}/scopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putApplicationsScopesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PutApplicationsScopes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putApplicationsScopesOptions.Scopes != nil {
		body["scopes"] = putApplicationsScopesOptions.Scopes
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetScopesForApplication)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetApplicationRoles : Get application roles
// View the defined roles for an application that is registered with an App ID instance.
func (appIdManagement *AppIDManagementV4) GetApplicationRoles(getApplicationRolesOptions *GetApplicationRolesOptions) (result *GetUserRolesResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetApplicationRolesWithContext(context.Background(), getApplicationRolesOptions)
}

// GetApplicationRolesWithContext is an alternate form of the GetApplicationRoles method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetApplicationRolesWithContext(ctx context.Context, getApplicationRolesOptions *GetApplicationRolesOptions) (result *GetUserRolesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getApplicationRolesOptions, "getApplicationRolesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getApplicationRolesOptions, "getApplicationRolesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getApplicationRolesOptions.TenantID,
		"clientId": *getApplicationRolesOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getApplicationRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetApplicationRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUserRolesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutApplicationsRoles : Add application role
// Update the roles for a registered application.</br>.
func (appIdManagement *AppIDManagementV4) PutApplicationsRoles(putApplicationsRolesOptions *PutApplicationsRolesOptions) (result *AssignRoleToUser, response *core.DetailedResponse, err error) {
	return appIdManagement.PutApplicationsRolesWithContext(context.Background(), putApplicationsRolesOptions)
}

// PutApplicationsRolesWithContext is an alternate form of the PutApplicationsRoles method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PutApplicationsRolesWithContext(ctx context.Context, putApplicationsRolesOptions *PutApplicationsRolesOptions) (result *AssignRoleToUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putApplicationsRolesOptions, "putApplicationsRolesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putApplicationsRolesOptions, "putApplicationsRolesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *putApplicationsRolesOptions.TenantID,
		"clientId": *putApplicationsRolesOptions.ClientID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/applications/{clientId}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putApplicationsRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PutApplicationsRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putApplicationsRolesOptions.Roles != nil {
		body["roles"] = putApplicationsRolesOptions.Roles
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAssignRoleToUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListCloudDirectoryUsers : List Cloud Directory users
// Get the list of Cloud Directory users. <a href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) ListCloudDirectoryUsers(listCloudDirectoryUsersOptions *ListCloudDirectoryUsersOptions) (result *UsersList, response *core.DetailedResponse, err error) {
	return appIdManagement.ListCloudDirectoryUsersWithContext(context.Background(), listCloudDirectoryUsersOptions)
}

// ListCloudDirectoryUsersWithContext is an alternate form of the ListCloudDirectoryUsers method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ListCloudDirectoryUsersWithContext(ctx context.Context, listCloudDirectoryUsersOptions *ListCloudDirectoryUsersOptions) (result *UsersList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listCloudDirectoryUsersOptions, "listCloudDirectoryUsersOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listCloudDirectoryUsersOptions, "listCloudDirectoryUsersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *listCloudDirectoryUsersOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCloudDirectoryUsersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ListCloudDirectoryUsers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listCloudDirectoryUsersOptions.StartIndex != nil {
		builder.AddQuery("startIndex", fmt.Sprint(*listCloudDirectoryUsersOptions.StartIndex))
	}
	if listCloudDirectoryUsersOptions.Count != nil {
		builder.AddQuery("count", fmt.Sprint(*listCloudDirectoryUsersOptions.Count))
	}
	if listCloudDirectoryUsersOptions.Query != nil {
		builder.AddQuery("query", fmt.Sprint(*listCloudDirectoryUsersOptions.Query))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUsersList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateCloudDirectoryUser : Create a Cloud Directory user
// Create a new record for Cloud Directory (no verification email is sent, and no profile is created).</br> To create a
// new Cloud Directory user use the  <a href="/swagger-ui/#/Management API - Cloud Directory Workflows/mgmt.startSignUp"
// target="_blank">sign_up</a> API. <a href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) CreateCloudDirectoryUser(createCloudDirectoryUserOptions *CreateCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.CreateCloudDirectoryUserWithContext(context.Background(), createCloudDirectoryUserOptions)
}

// CreateCloudDirectoryUserWithContext is an alternate form of the CreateCloudDirectoryUser method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CreateCloudDirectoryUserWithContext(ctx context.Context, createCloudDirectoryUserOptions *CreateCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCloudDirectoryUserOptions, "createCloudDirectoryUserOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCloudDirectoryUserOptions, "createCloudDirectoryUserOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *createCloudDirectoryUserOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCloudDirectoryUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CreateCloudDirectoryUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createCloudDirectoryUserOptions.Emails != nil {
		body["emails"] = createCloudDirectoryUserOptions.Emails
	}
	if createCloudDirectoryUserOptions.Password != nil {
		body["password"] = createCloudDirectoryUserOptions.Password
	}
	if createCloudDirectoryUserOptions.Active != nil {
		body["active"] = createCloudDirectoryUserOptions.Active
	}
	if createCloudDirectoryUserOptions.LockedUntil != nil {
		body["lockedUntil"] = createCloudDirectoryUserOptions.LockedUntil
	}
	if createCloudDirectoryUserOptions.DisplayName != nil {
		body["displayName"] = createCloudDirectoryUserOptions.DisplayName
	}
	if createCloudDirectoryUserOptions.UserName != nil {
		body["userName"] = createCloudDirectoryUserOptions.UserName
	}
	if createCloudDirectoryUserOptions.Status != nil {
		body["status"] = createCloudDirectoryUserOptions.Status
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCloudDirectoryUser : Get a Cloud Directory user
// Returns the requested Cloud Directory user object. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryUser(getCloudDirectoryUserOptions *GetCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryUserWithContext(context.Background(), getCloudDirectoryUserOptions)
}

// GetCloudDirectoryUserWithContext is an alternate form of the GetCloudDirectoryUser method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryUserWithContext(ctx context.Context, getCloudDirectoryUserOptions *GetCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryUserOptions, "getCloudDirectoryUserOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryUserOptions, "getCloudDirectoryUserOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryUserOptions.TenantID,
		"userId":   *getCloudDirectoryUserOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users/{userId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateCloudDirectoryUser : Update a Cloud Directory user
// Updates an existing Cloud Directory user. <a href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateCloudDirectoryUser(updateCloudDirectoryUserOptions *UpdateCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateCloudDirectoryUserWithContext(context.Background(), updateCloudDirectoryUserOptions)
}

// UpdateCloudDirectoryUserWithContext is an alternate form of the UpdateCloudDirectoryUser method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateCloudDirectoryUserWithContext(ctx context.Context, updateCloudDirectoryUserOptions *UpdateCloudDirectoryUserOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCloudDirectoryUserOptions, "updateCloudDirectoryUserOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateCloudDirectoryUserOptions, "updateCloudDirectoryUserOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateCloudDirectoryUserOptions.TenantID,
		"userId":   *updateCloudDirectoryUserOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users/{userId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCloudDirectoryUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateCloudDirectoryUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCloudDirectoryUserOptions.Emails != nil {
		body["emails"] = updateCloudDirectoryUserOptions.Emails
	}
	if updateCloudDirectoryUserOptions.Status != nil {
		body["status"] = updateCloudDirectoryUserOptions.Status
	}
	if updateCloudDirectoryUserOptions.DisplayName != nil {
		body["displayName"] = updateCloudDirectoryUserOptions.DisplayName
	}
	if updateCloudDirectoryUserOptions.UserName != nil {
		body["userName"] = updateCloudDirectoryUserOptions.UserName
	}
	if updateCloudDirectoryUserOptions.Password != nil {
		body["password"] = updateCloudDirectoryUserOptions.Password
	}
	if updateCloudDirectoryUserOptions.Active != nil {
		body["active"] = updateCloudDirectoryUserOptions.Active
	}
	if updateCloudDirectoryUserOptions.LockedUntil != nil {
		body["lockedUntil"] = updateCloudDirectoryUserOptions.LockedUntil
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCloudDirectoryUser : Delete a Cloud Directory user
// Deletes an existing Cloud Directory recored (without removing the associated profile). <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>.</br> To remove a Cloud
// Directory user use the <a href="/swagger-ui/#/Management API - Cloud Directory Workflows/mgmt.cloud_directory_remove"
// target="_blank">remove</a> API. </br> <b>Note: This action cannot be undone</b>.
func (appIdManagement *AppIDManagementV4) DeleteCloudDirectoryUser(deleteCloudDirectoryUserOptions *DeleteCloudDirectoryUserOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.DeleteCloudDirectoryUserWithContext(context.Background(), deleteCloudDirectoryUserOptions)
}

// DeleteCloudDirectoryUserWithContext is an alternate form of the DeleteCloudDirectoryUser method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) DeleteCloudDirectoryUserWithContext(ctx context.Context, deleteCloudDirectoryUserOptions *DeleteCloudDirectoryUserOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCloudDirectoryUserOptions, "deleteCloudDirectoryUserOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCloudDirectoryUserOptions, "deleteCloudDirectoryUserOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *deleteCloudDirectoryUserOptions.TenantID,
		"userId":   *deleteCloudDirectoryUserOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users/{userId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCloudDirectoryUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "DeleteCloudDirectoryUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// InvalidateUserSSOSessions : Invalidate all SSO sessions
// Invalidate all the user's SSO sessions. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-sso#ending-all-sessions-for-a-user" target="_blank">Learn
// more</a>.
func (appIdManagement *AppIDManagementV4) InvalidateUserSSOSessions(invalidateUserSSOSessionsOptions *InvalidateUserSSOSessionsOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.InvalidateUserSSOSessionsWithContext(context.Background(), invalidateUserSSOSessionsOptions)
}

// InvalidateUserSSOSessionsWithContext is an alternate form of the InvalidateUserSSOSessions method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) InvalidateUserSSOSessionsWithContext(ctx context.Context, invalidateUserSSOSessionsOptions *InvalidateUserSSOSessionsOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(invalidateUserSSOSessionsOptions, "invalidateUserSSOSessionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(invalidateUserSSOSessionsOptions, "invalidateUserSSOSessionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *invalidateUserSSOSessionsOptions.TenantID,
		"userId":   *invalidateUserSSOSessionsOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/Users/{userId}/sso/logout`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range invalidateUserSSOSessionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "InvalidateUserSSOSessions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// CloudDirectoryExport : Export Cloud Directory users
// Exports Cloud Directory users with their profile attributes and hashed passwords. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) CloudDirectoryExport(cloudDirectoryExportOptions *CloudDirectoryExportOptions) (result *ExportUser, response *core.DetailedResponse, err error) {
	return appIdManagement.CloudDirectoryExportWithContext(context.Background(), cloudDirectoryExportOptions)
}

// CloudDirectoryExportWithContext is an alternate form of the CloudDirectoryExport method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CloudDirectoryExportWithContext(ctx context.Context, cloudDirectoryExportOptions *CloudDirectoryExportOptions) (result *ExportUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cloudDirectoryExportOptions, "cloudDirectoryExportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(cloudDirectoryExportOptions, "cloudDirectoryExportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *cloudDirectoryExportOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/export`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range cloudDirectoryExportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CloudDirectoryExport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("encryption_secret", fmt.Sprint(*cloudDirectoryExportOptions.EncryptionSecret))
	if cloudDirectoryExportOptions.StartIndex != nil {
		builder.AddQuery("startIndex", fmt.Sprint(*cloudDirectoryExportOptions.StartIndex))
	}
	if cloudDirectoryExportOptions.Count != nil {
		builder.AddQuery("count", fmt.Sprint(*cloudDirectoryExportOptions.Count))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExportUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CloudDirectoryImport : Import Cloud Directory users
// Imports Cloud Directory users list that was exported using the /export API. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) CloudDirectoryImport(cloudDirectoryImportOptions *CloudDirectoryImportOptions) (result *ImportResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.CloudDirectoryImportWithContext(context.Background(), cloudDirectoryImportOptions)
}

// CloudDirectoryImportWithContext is an alternate form of the CloudDirectoryImport method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CloudDirectoryImportWithContext(ctx context.Context, cloudDirectoryImportOptions *CloudDirectoryImportOptions) (result *ImportResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cloudDirectoryImportOptions, "cloudDirectoryImportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(cloudDirectoryImportOptions, "cloudDirectoryImportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *cloudDirectoryImportOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/import`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range cloudDirectoryImportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CloudDirectoryImport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("encryption_secret", fmt.Sprint(*cloudDirectoryImportOptions.EncryptionSecret))

	body := make(map[string]interface{})
	if cloudDirectoryImportOptions.Users != nil {
		body["users"] = cloudDirectoryImportOptions.Users
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImportResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CloudDirectoryGetUserinfo : Get Cloud Directory SCIM and Attributes
// Returns the Cloud Directory user SCIM and the Profile related to it. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) CloudDirectoryGetUserinfo(cloudDirectoryGetUserinfoOptions *CloudDirectoryGetUserinfoOptions) (result *GetUserAndProfile, response *core.DetailedResponse, err error) {
	return appIdManagement.CloudDirectoryGetUserinfoWithContext(context.Background(), cloudDirectoryGetUserinfoOptions)
}

// CloudDirectoryGetUserinfoWithContext is an alternate form of the CloudDirectoryGetUserinfo method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CloudDirectoryGetUserinfoWithContext(ctx context.Context, cloudDirectoryGetUserinfoOptions *CloudDirectoryGetUserinfoOptions) (result *GetUserAndProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cloudDirectoryGetUserinfoOptions, "cloudDirectoryGetUserinfoOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(cloudDirectoryGetUserinfoOptions, "cloudDirectoryGetUserinfoOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *cloudDirectoryGetUserinfoOptions.TenantID,
		"userId":   *cloudDirectoryGetUserinfoOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/{userId}/userinfo`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range cloudDirectoryGetUserinfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CloudDirectoryGetUserinfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUserAndProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// StartSignUp : Sign up
// Start the sign up process <a href="https://cloud.ibm.com/docs/appid?topic=appid-branded" target="_blank">Learn
// more</a>.
func (appIdManagement *AppIDManagementV4) StartSignUp(startSignUpOptions *StartSignUpOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.StartSignUpWithContext(context.Background(), startSignUpOptions)
}

// StartSignUpWithContext is an alternate form of the StartSignUp method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) StartSignUpWithContext(ctx context.Context, startSignUpOptions *StartSignUpOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(startSignUpOptions, "startSignUpOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(startSignUpOptions, "startSignUpOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *startSignUpOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/sign_up`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range startSignUpOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "StartSignUp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("shouldCreateProfile", fmt.Sprint(*startSignUpOptions.ShouldCreateProfile))
	if startSignUpOptions.Language != nil {
		builder.AddQuery("language", fmt.Sprint(*startSignUpOptions.Language))
	}

	body := make(map[string]interface{})
	if startSignUpOptions.Emails != nil {
		body["emails"] = startSignUpOptions.Emails
	}
	if startSignUpOptions.Password != nil {
		body["password"] = startSignUpOptions.Password
	}
	if startSignUpOptions.Active != nil {
		body["active"] = startSignUpOptions.Active
	}
	if startSignUpOptions.LockedUntil != nil {
		body["lockedUntil"] = startSignUpOptions.LockedUntil
	}
	if startSignUpOptions.DisplayName != nil {
		body["displayName"] = startSignUpOptions.DisplayName
	}
	if startSignUpOptions.UserName != nil {
		body["userName"] = startSignUpOptions.UserName
	}
	if startSignUpOptions.Status != nil {
		body["status"] = startSignUpOptions.Status
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UserVerificationResult : Get signup confirmation result
// Returns the sign up confirmation result. <a href="https://cloud.ibm.com/docs/appid?topic=appid-branded"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UserVerificationResult(userVerificationResultOptions *UserVerificationResultOptions) (result *ConfirmationResultOk, response *core.DetailedResponse, err error) {
	return appIdManagement.UserVerificationResultWithContext(context.Background(), userVerificationResultOptions)
}

// UserVerificationResultWithContext is an alternate form of the UserVerificationResult method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UserVerificationResultWithContext(ctx context.Context, userVerificationResultOptions *UserVerificationResultOptions) (result *ConfirmationResultOk, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(userVerificationResultOptions, "userVerificationResultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(userVerificationResultOptions, "userVerificationResultOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *userVerificationResultOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/sign_up/confirmation_result`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range userVerificationResultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UserVerificationResult")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if userVerificationResultOptions.Context != nil {
		body["context"] = userVerificationResultOptions.Context
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfirmationResultOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// StartForgotPassword : Forgot password
// Starts the forgot password process. <a href="https://cloud.ibm.com/docs/appid?topic=appid-branded"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) StartForgotPassword(startForgotPasswordOptions *StartForgotPasswordOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.StartForgotPasswordWithContext(context.Background(), startForgotPasswordOptions)
}

// StartForgotPasswordWithContext is an alternate form of the StartForgotPassword method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) StartForgotPasswordWithContext(ctx context.Context, startForgotPasswordOptions *StartForgotPasswordOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(startForgotPasswordOptions, "startForgotPasswordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(startForgotPasswordOptions, "startForgotPasswordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *startForgotPasswordOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/forgot_password`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range startForgotPasswordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "StartForgotPassword")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if startForgotPasswordOptions.Language != nil {
		builder.AddQuery("language", fmt.Sprint(*startForgotPasswordOptions.Language))
	}

	body := make(map[string]interface{})
	if startForgotPasswordOptions.User != nil {
		body["user"] = startForgotPasswordOptions.User
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ForgotPasswordResult : Forgot password confirmation result
// Returns the forgot password flow confirmation result. <a href="https://cloud.ibm.com/docs/appid?topic=appid-branded"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) ForgotPasswordResult(forgotPasswordResultOptions *ForgotPasswordResultOptions) (result *ConfirmationResultOk, response *core.DetailedResponse, err error) {
	return appIdManagement.ForgotPasswordResultWithContext(context.Background(), forgotPasswordResultOptions)
}

// ForgotPasswordResultWithContext is an alternate form of the ForgotPasswordResult method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ForgotPasswordResultWithContext(ctx context.Context, forgotPasswordResultOptions *ForgotPasswordResultOptions) (result *ConfirmationResultOk, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(forgotPasswordResultOptions, "forgotPasswordResultOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(forgotPasswordResultOptions, "forgotPasswordResultOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *forgotPasswordResultOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/forgot_password/confirmation_result`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range forgotPasswordResultOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ForgotPasswordResult")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if forgotPasswordResultOptions.Context != nil {
		body["context"] = forgotPasswordResultOptions.Context
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfirmationResultOk)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ChangePassword : Change password
// Changes the Cloud Directory user password. <a href="https://cloud.ibm.com/docs/appid?topic=appid-branded"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) ChangePassword(changePasswordOptions *ChangePasswordOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	return appIdManagement.ChangePasswordWithContext(context.Background(), changePasswordOptions)
}

// ChangePasswordWithContext is an alternate form of the ChangePassword method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ChangePasswordWithContext(ctx context.Context, changePasswordOptions *ChangePasswordOptions) (result *GetUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(changePasswordOptions, "changePasswordOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(changePasswordOptions, "changePasswordOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *changePasswordOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/change_password`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changePasswordOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ChangePassword")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if changePasswordOptions.Language != nil {
		builder.AddQuery("language", fmt.Sprint(*changePasswordOptions.Language))
	}

	body := make(map[string]interface{})
	if changePasswordOptions.NewPassword != nil {
		body["newPassword"] = changePasswordOptions.NewPassword
	}
	if changePasswordOptions.UUID != nil {
		body["uuid"] = changePasswordOptions.UUID
	}
	if changePasswordOptions.ChangedIPAddress != nil {
		body["changedIpAddress"] = changePasswordOptions.ChangedIPAddress
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ResendNotification : Resend user notifications
// Resend user email notifications (e.g. resend user verification email). <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-branded" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) ResendNotification(resendNotificationOptions *ResendNotificationOptions) (result *ResendNotificationResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.ResendNotificationWithContext(context.Background(), resendNotificationOptions)
}

// ResendNotificationWithContext is an alternate form of the ResendNotification method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ResendNotificationWithContext(ctx context.Context, resendNotificationOptions *ResendNotificationOptions) (result *ResendNotificationResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(resendNotificationOptions, "resendNotificationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(resendNotificationOptions, "resendNotificationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId":     *resendNotificationOptions.TenantID,
		"templateName": *resendNotificationOptions.TemplateName,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/resend/{templateName}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range resendNotificationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ResendNotification")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if resendNotificationOptions.Language != nil {
		builder.AddQuery("language", fmt.Sprint(*resendNotificationOptions.Language))
	}

	body := make(map[string]interface{})
	if resendNotificationOptions.UUID != nil {
		body["uuid"] = resendNotificationOptions.UUID
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResendNotificationResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CloudDirectoryRemove : Delete Cloud Directory User and Profile
// Deletes an existing Cloud Directory user and the Profile related to it. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>. <b>Note: This action
// cannot be undone</b>.
func (appIdManagement *AppIDManagementV4) CloudDirectoryRemove(cloudDirectoryRemoveOptions *CloudDirectoryRemoveOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.CloudDirectoryRemoveWithContext(context.Background(), cloudDirectoryRemoveOptions)
}

// CloudDirectoryRemoveWithContext is an alternate form of the CloudDirectoryRemove method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CloudDirectoryRemoveWithContext(ctx context.Context, cloudDirectoryRemoveOptions *CloudDirectoryRemoveOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cloudDirectoryRemoveOptions, "cloudDirectoryRemoveOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(cloudDirectoryRemoveOptions, "cloudDirectoryRemoveOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *cloudDirectoryRemoveOptions.TenantID,
		"userId":   *cloudDirectoryRemoveOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/cloud_directory/remove/{userId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range cloudDirectoryRemoveOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CloudDirectoryRemove")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetTokensConfig : Get tokens configuration
// Returns the token configuration. <a href="https://cloud.ibm.com/docs/appid?topic=appid-key-concepts"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetTokensConfig(getTokensConfigOptions *GetTokensConfigOptions) (result *TokensConfigResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetTokensConfigWithContext(context.Background(), getTokensConfigOptions)
}

// GetTokensConfigWithContext is an alternate form of the GetTokensConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetTokensConfigWithContext(ctx context.Context, getTokensConfigOptions *GetTokensConfigOptions) (result *TokensConfigResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTokensConfigOptions, "getTokensConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTokensConfigOptions, "getTokensConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getTokensConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/tokens`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTokensConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetTokensConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTokensConfigResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PutTokensConfig : Update tokens configuration
// Update the tokens' configuration to fine-tune the expiration times of access, id and refresh tokens, to
// enable/disable refresh and anonymous tokens, and to configure custom claims. When a token config object is not
// included in the set, its value will be reset back to default. <br> For more information, check out the <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-key-concepts" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) PutTokensConfig(putTokensConfigOptions *PutTokensConfigOptions) (result *TokensConfigResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.PutTokensConfigWithContext(context.Background(), putTokensConfigOptions)
}

// PutTokensConfigWithContext is an alternate form of the PutTokensConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PutTokensConfigWithContext(ctx context.Context, putTokensConfigOptions *PutTokensConfigOptions) (result *TokensConfigResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putTokensConfigOptions, "putTokensConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(putTokensConfigOptions, "putTokensConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *putTokensConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/tokens`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range putTokensConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PutTokensConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putTokensConfigOptions.IDTokenClaims != nil {
		body["idTokenClaims"] = putTokensConfigOptions.IDTokenClaims
	}
	if putTokensConfigOptions.AccessTokenClaims != nil {
		body["accessTokenClaims"] = putTokensConfigOptions.AccessTokenClaims
	}
	if putTokensConfigOptions.Access != nil {
		body["access"] = putTokensConfigOptions.Access
	}
	if putTokensConfigOptions.Refresh != nil {
		body["refresh"] = putTokensConfigOptions.Refresh
	}
	if putTokensConfigOptions.AnonymousAccess != nil {
		body["anonymousAccess"] = putTokensConfigOptions.AnonymousAccess
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTokensConfigResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRedirectUris : Get redirect URIs
// Returns the list of the redirect URIs that can be used as callbacks of App ID authentication flow. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-managing-idp#add-redirect-uri" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetRedirectUris(getRedirectUrisOptions *GetRedirectUrisOptions) (result *RedirectURIResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetRedirectUrisWithContext(context.Background(), getRedirectUrisOptions)
}

// GetRedirectUrisWithContext is an alternate form of the GetRedirectUris method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetRedirectUrisWithContext(ctx context.Context, getRedirectUrisOptions *GetRedirectUrisOptions) (result *RedirectURIResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRedirectUrisOptions, "getRedirectUrisOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRedirectUrisOptions, "getRedirectUrisOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getRedirectUrisOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/redirect_uris`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRedirectUrisOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetRedirectUris")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRedirectURIResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateRedirectUris : Update redirect URIs
// Update the list of the redirect URIs that can be used as callbacks of App ID authentication flow. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-managing-idp#add-redirect-uri" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateRedirectUris(updateRedirectUrisOptions *UpdateRedirectUrisOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateRedirectUrisWithContext(context.Background(), updateRedirectUrisOptions)
}

// UpdateRedirectUrisWithContext is an alternate form of the UpdateRedirectUris method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateRedirectUrisWithContext(ctx context.Context, updateRedirectUrisOptions *UpdateRedirectUrisOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRedirectUrisOptions, "updateRedirectUrisOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRedirectUrisOptions, "updateRedirectUrisOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateRedirectUrisOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/redirect_uris`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRedirectUrisOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateRedirectUris")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(updateRedirectUrisOptions.RedirectUrisArray)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetUserProfilesConfig : Get user profiles configuration
// A user profile is an entity that is stored and maintained by App ID. The profile holds a user's attributes and
// identity. It can be anonymous or linked to an identity that is managed by an identity provider. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetUserProfilesConfig(getUserProfilesConfigOptions *GetUserProfilesConfigOptions) (result *GetUserProfilesConfigResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetUserProfilesConfigWithContext(context.Background(), getUserProfilesConfigOptions)
}

// GetUserProfilesConfigWithContext is an alternate form of the GetUserProfilesConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetUserProfilesConfigWithContext(ctx context.Context, getUserProfilesConfigOptions *GetUserProfilesConfigOptions) (result *GetUserProfilesConfigResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getUserProfilesConfigOptions, "getUserProfilesConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getUserProfilesConfigOptions, "getUserProfilesConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getUserProfilesConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/users_profile`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUserProfilesConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetUserProfilesConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUserProfilesConfigResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateUserProfilesConfig : Update user profiles configuration
// A user profile is an entity that is stored and maintained by App ID. The profile holds a user's attributes and
// identity. It can be anonymous or linked to an identity that is managed by an identity provider. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateUserProfilesConfig(updateUserProfilesConfigOptions *UpdateUserProfilesConfigOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateUserProfilesConfigWithContext(context.Background(), updateUserProfilesConfigOptions)
}

// UpdateUserProfilesConfigWithContext is an alternate form of the UpdateUserProfilesConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateUserProfilesConfigWithContext(ctx context.Context, updateUserProfilesConfigOptions *UpdateUserProfilesConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateUserProfilesConfigOptions, "updateUserProfilesConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateUserProfilesConfigOptions, "updateUserProfilesConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateUserProfilesConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/users_profile`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateUserProfilesConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateUserProfilesConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateUserProfilesConfigOptions.IsActive != nil {
		body["isActive"] = updateUserProfilesConfigOptions.IsActive
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetThemeText : Get widget texts
// Get the theme texts of the App ID login widget. <a href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetThemeText(getThemeTextOptions *GetThemeTextOptions) (result *GetThemeTextResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetThemeTextWithContext(context.Background(), getThemeTextOptions)
}

// GetThemeTextWithContext is an alternate form of the GetThemeText method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetThemeTextWithContext(ctx context.Context, getThemeTextOptions *GetThemeTextOptions) (result *GetThemeTextResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getThemeTextOptions, "getThemeTextOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getThemeTextOptions, "getThemeTextOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getThemeTextOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/theme_text`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getThemeTextOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetThemeText")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetThemeTextResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostThemeText : Update widget texts
// Update the texts of the App ID login widget. <a href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) PostThemeText(postThemeTextOptions *PostThemeTextOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.PostThemeTextWithContext(context.Background(), postThemeTextOptions)
}

// PostThemeTextWithContext is an alternate form of the PostThemeText method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostThemeTextWithContext(ctx context.Context, postThemeTextOptions *PostThemeTextOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postThemeTextOptions, "postThemeTextOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postThemeTextOptions, "postThemeTextOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postThemeTextOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/theme_text`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postThemeTextOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostThemeText")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postThemeTextOptions.TabTitle != nil {
		body["tabTitle"] = postThemeTextOptions.TabTitle
	}
	if postThemeTextOptions.Footnote != nil {
		body["footnote"] = postThemeTextOptions.Footnote
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetThemeColor : Get widget colors
// Get the colors of the App ID login widget. <a href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetThemeColor(getThemeColorOptions *GetThemeColorOptions) (result *GetThemeColorResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetThemeColorWithContext(context.Background(), getThemeColorOptions)
}

// GetThemeColorWithContext is an alternate form of the GetThemeColor method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetThemeColorWithContext(ctx context.Context, getThemeColorOptions *GetThemeColorOptions) (result *GetThemeColorResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getThemeColorOptions, "getThemeColorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getThemeColorOptions, "getThemeColorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getThemeColorOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/theme_color`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getThemeColorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetThemeColor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetThemeColorResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostThemeColor : Update widget colors
// Update the colors of the App ID login widget. <a href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) PostThemeColor(postThemeColorOptions *PostThemeColorOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.PostThemeColorWithContext(context.Background(), postThemeColorOptions)
}

// PostThemeColorWithContext is an alternate form of the PostThemeColor method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostThemeColorWithContext(ctx context.Context, postThemeColorOptions *PostThemeColorOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postThemeColorOptions, "postThemeColorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postThemeColorOptions, "postThemeColorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postThemeColorOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/theme_color`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postThemeColorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostThemeColor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postThemeColorOptions.HeaderColor != nil {
		body["headerColor"] = postThemeColorOptions.HeaderColor
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetMedia : Get widget logo
// Returns the link to the custom logo image of the login widget. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetMedia(getMediaOptions *GetMediaOptions) (result *GetMediaResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetMediaWithContext(context.Background(), getMediaOptions)
}

// GetMediaWithContext is an alternate form of the GetMedia method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetMediaWithContext(ctx context.Context, getMediaOptions *GetMediaOptions) (result *GetMediaResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMediaOptions, "getMediaOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getMediaOptions, "getMediaOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getMediaOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/media`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getMediaOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetMedia")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetMediaResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostMedia : Update widget logo
// You can update the image file shown in the login widget. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-login-widget" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) PostMedia(postMediaOptions *PostMediaOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.PostMediaWithContext(context.Background(), postMediaOptions)
}

// PostMediaWithContext is an alternate form of the PostMedia method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostMediaWithContext(ctx context.Context, postMediaOptions *PostMediaOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postMediaOptions, "postMediaOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postMediaOptions, "postMediaOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postMediaOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/media`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postMediaOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostMedia")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddQuery("mediaType", fmt.Sprint(*postMediaOptions.MediaType))

	builder.AddFormData("file", "filename",
		core.StringNilMapper(postMediaOptions.FileContentType), postMediaOptions.File)

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetSAMLMetadata : Get the SAML metadata
// Returns the SAML metadata required in order to integrate App ID with a SAML identity provider. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-enterprise" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetSAMLMetadata(getSAMLMetadataOptions *GetSAMLMetadataOptions) (result *string, response *core.DetailedResponse, err error) {
	return appIdManagement.GetSAMLMetadataWithContext(context.Background(), getSAMLMetadataOptions)
}

// GetSAMLMetadataWithContext is an alternate form of the GetSAMLMetadata method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetSAMLMetadataWithContext(ctx context.Context, getSAMLMetadataOptions *GetSAMLMetadataOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSAMLMetadataOptions, "getSAMLMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSAMLMetadataOptions, "getSAMLMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getSAMLMetadataOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/saml_metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSAMLMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetSAMLMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/xml")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, &result)

	return
}

// GetTemplate : Get an email template
// Returns the content of a custom email template or the default template in case it wasn't customized. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-types" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetTemplate(getTemplateOptions *GetTemplateOptions) (result *GetTemplate, response *core.DetailedResponse, err error) {
	return appIdManagement.GetTemplateWithContext(context.Background(), getTemplateOptions)
}

// GetTemplateWithContext is an alternate form of the GetTemplate method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetTemplateWithContext(ctx context.Context, getTemplateOptions *GetTemplateOptions) (result *GetTemplate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTemplateOptions, "getTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTemplateOptions, "getTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId":     *getTemplateOptions.TenantID,
		"templateName": *getTemplateOptions.TemplateName,
		"language":     *getTemplateOptions.Language,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/templates/{templateName}/{language}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetTemplate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateTemplate : Update an email template
// Updates the Cloud Directory email template. <a href="https://cloud.ibm.com/docs/appid?topic=appid-cd-types"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateTemplate(updateTemplateOptions *UpdateTemplateOptions) (result *GetTemplate, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateTemplateWithContext(context.Background(), updateTemplateOptions)
}

// UpdateTemplateWithContext is an alternate form of the UpdateTemplate method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateTemplateWithContext(ctx context.Context, updateTemplateOptions *UpdateTemplateOptions) (result *GetTemplate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTemplateOptions, "updateTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateTemplateOptions, "updateTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId":     *updateTemplateOptions.TenantID,
		"templateName": *updateTemplateOptions.TemplateName,
		"language":     *updateTemplateOptions.Language,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/templates/{templateName}/{language}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateTemplateOptions.Subject != nil {
		body["subject"] = updateTemplateOptions.Subject
	}
	if updateTemplateOptions.HTMLBody != nil {
		body["html_body"] = updateTemplateOptions.HTMLBody
	}
	if updateTemplateOptions.Base64EncodedHTMLBody != nil {
		body["base64_encoded_html_body"] = updateTemplateOptions.Base64EncodedHTMLBody
	}
	if updateTemplateOptions.PlainTextBody != nil {
		body["plain_text_body"] = updateTemplateOptions.PlainTextBody
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetTemplate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTemplate : Delete an email template
// Delete the customized email template and reverts to App ID default template. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-users" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) DeleteTemplate(deleteTemplateOptions *DeleteTemplateOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.DeleteTemplateWithContext(context.Background(), deleteTemplateOptions)
}

// DeleteTemplateWithContext is an alternate form of the DeleteTemplate method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) DeleteTemplateWithContext(ctx context.Context, deleteTemplateOptions *DeleteTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTemplateOptions, "deleteTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTemplateOptions, "deleteTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId":     *deleteTemplateOptions.TenantID,
		"templateName": *deleteTemplateOptions.TemplateName,
		"language":     *deleteTemplateOptions.Language,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/templates/{templateName}/{language}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "DeleteTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetLocalization : Get languages
// Returns the list of languages that can be used to customize email templates for Cloud Directory.
func (appIdManagement *AppIDManagementV4) GetLocalization(getLocalizationOptions *GetLocalizationOptions) (result *GetLanguages, response *core.DetailedResponse, err error) {
	return appIdManagement.GetLocalizationWithContext(context.Background(), getLocalizationOptions)
}

// GetLocalizationWithContext is an alternate form of the GetLocalization method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetLocalizationWithContext(ctx context.Context, getLocalizationOptions *GetLocalizationOptions) (result *GetLanguages, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLocalizationOptions, "getLocalizationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLocalizationOptions, "getLocalizationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getLocalizationOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/languages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLocalizationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetLocalization")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetLanguages)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateLocalization : Update languages
// Update the list of languages that can be used to customize email templates for Cloud Directory.
func (appIdManagement *AppIDManagementV4) UpdateLocalization(updateLocalizationOptions *UpdateLocalizationOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateLocalizationWithContext(context.Background(), updateLocalizationOptions)
}

// UpdateLocalizationWithContext is an alternate form of the UpdateLocalization method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateLocalizationWithContext(ctx context.Context, updateLocalizationOptions *UpdateLocalizationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLocalizationOptions, "updateLocalizationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLocalizationOptions, "updateLocalizationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateLocalizationOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/ui/languages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateLocalizationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateLocalization")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateLocalizationOptions.Languages != nil {
		body["languages"] = updateLocalizationOptions.Languages
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetCloudDirectorySenderDetails : Get sender details
// Returns the sender details configuration that is used by Cloud Directory when sending emails. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-types" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetCloudDirectorySenderDetails(getCloudDirectorySenderDetailsOptions *GetCloudDirectorySenderDetailsOptions) (result *CloudDirectorySenderDetails, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectorySenderDetailsWithContext(context.Background(), getCloudDirectorySenderDetailsOptions)
}

// GetCloudDirectorySenderDetailsWithContext is an alternate form of the GetCloudDirectorySenderDetails method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectorySenderDetailsWithContext(ctx context.Context, getCloudDirectorySenderDetailsOptions *GetCloudDirectorySenderDetailsOptions) (result *CloudDirectorySenderDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectorySenderDetailsOptions, "getCloudDirectorySenderDetailsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectorySenderDetailsOptions, "getCloudDirectorySenderDetailsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectorySenderDetailsOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/sender_details`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectorySenderDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectorySenderDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCloudDirectorySenderDetails)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectorySenderDetails : Update the sender details
// Updates the sender details configuration that is used by Cloud Directory when sending emails. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-types" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetCloudDirectorySenderDetails(setCloudDirectorySenderDetailsOptions *SetCloudDirectorySenderDetailsOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectorySenderDetailsWithContext(context.Background(), setCloudDirectorySenderDetailsOptions)
}

// SetCloudDirectorySenderDetailsWithContext is an alternate form of the SetCloudDirectorySenderDetails method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectorySenderDetailsWithContext(ctx context.Context, setCloudDirectorySenderDetailsOptions *SetCloudDirectorySenderDetailsOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectorySenderDetailsOptions, "setCloudDirectorySenderDetailsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectorySenderDetailsOptions, "setCloudDirectorySenderDetailsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectorySenderDetailsOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/sender_details`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectorySenderDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectorySenderDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectorySenderDetailsOptions.SenderDetails != nil {
		body["senderDetails"] = setCloudDirectorySenderDetailsOptions.SenderDetails
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetCloudDirectoryActionURL : Get action url
// Get the custom url to redirect to when <b>action</b> is executed. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryActionURL(getCloudDirectoryActionURLOptions *GetCloudDirectoryActionURLOptions) (result *ActionURLResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryActionURLWithContext(context.Background(), getCloudDirectoryActionURLOptions)
}

// GetCloudDirectoryActionURLWithContext is an alternate form of the GetCloudDirectoryActionURL method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryActionURLWithContext(ctx context.Context, getCloudDirectoryActionURLOptions *GetCloudDirectoryActionURLOptions) (result *ActionURLResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryActionURLOptions, "getCloudDirectoryActionURLOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryActionURLOptions, "getCloudDirectoryActionURLOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryActionURLOptions.TenantID,
		"action":   *getCloudDirectoryActionURLOptions.Action,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/action_url/{action}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryActionURLOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryActionURL")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActionURLResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectoryAction : Update action url
// Updates the custom url to redirect to when <b>action</b> is executed. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryAction(setCloudDirectoryActionOptions *SetCloudDirectoryActionOptions) (result *ActionURLResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectoryActionWithContext(context.Background(), setCloudDirectoryActionOptions)
}

// SetCloudDirectoryActionWithContext is an alternate form of the SetCloudDirectoryAction method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryActionWithContext(ctx context.Context, setCloudDirectoryActionOptions *SetCloudDirectoryActionOptions) (result *ActionURLResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectoryActionOptions, "setCloudDirectoryActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectoryActionOptions, "setCloudDirectoryActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectoryActionOptions.TenantID,
		"action":   *setCloudDirectoryActionOptions.Action,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/action_url/{action}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectoryActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectoryAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectoryActionOptions.ActionURL != nil {
		body["actionUrl"] = setCloudDirectoryActionOptions.ActionURL
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActionURLResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteActionURL : Delete action url
// Delete the custom url to redirect to when <b>action</b> is executed. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) DeleteActionURL(deleteActionURLOptions *DeleteActionURLOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.DeleteActionURLWithContext(context.Background(), deleteActionURLOptions)
}

// DeleteActionURLWithContext is an alternate form of the DeleteActionURL method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) DeleteActionURLWithContext(ctx context.Context, deleteActionURLOptions *DeleteActionURLOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteActionURLOptions, "deleteActionURLOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteActionURLOptions, "deleteActionURLOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *deleteActionURLOptions.TenantID,
		"action":   *deleteActionURLOptions.Action,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/action_url/{action}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteActionURLOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "DeleteActionURL")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetCloudDirectoryPasswordRegex : Get password regex
// Returns the regular expression used by App ID for password strength validation. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-strength" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryPasswordRegex(getCloudDirectoryPasswordRegexOptions *GetCloudDirectoryPasswordRegexOptions) (result *PasswordRegexConfigParamsGet, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryPasswordRegexWithContext(context.Background(), getCloudDirectoryPasswordRegexOptions)
}

// GetCloudDirectoryPasswordRegexWithContext is an alternate form of the GetCloudDirectoryPasswordRegex method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryPasswordRegexWithContext(ctx context.Context, getCloudDirectoryPasswordRegexOptions *GetCloudDirectoryPasswordRegexOptions) (result *PasswordRegexConfigParamsGet, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryPasswordRegexOptions, "getCloudDirectoryPasswordRegexOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryPasswordRegexOptions, "getCloudDirectoryPasswordRegexOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryPasswordRegexOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/password_regex`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryPasswordRegexOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryPasswordRegex")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPasswordRegexConfigParamsGet)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectoryPasswordRegex : Update password regex
// Updates the regular expression used by App ID for password strength validation.<br />For example, the regular
// expression: <code>"^[A-Za-z\d]*$"</code> should be passed as:<br /><code>{<br />&nbsp;&nbsp;"base64_encoded_regex":
// "XltBLVphLXpcZF0qJA==", <br />&nbsp;&nbsp;"error_message": "Must only contain letters and digits"<br />}</code> <br
// /><br /> <a href="https://cloud.ibm.com/docs/appid?topic=appid-cd-strength" target="_blank" rel="noopener">Learn
// more</a>.
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryPasswordRegex(setCloudDirectoryPasswordRegexOptions *SetCloudDirectoryPasswordRegexOptions) (result *PasswordRegexConfigParamsGet, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectoryPasswordRegexWithContext(context.Background(), setCloudDirectoryPasswordRegexOptions)
}

// SetCloudDirectoryPasswordRegexWithContext is an alternate form of the SetCloudDirectoryPasswordRegex method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryPasswordRegexWithContext(ctx context.Context, setCloudDirectoryPasswordRegexOptions *SetCloudDirectoryPasswordRegexOptions) (result *PasswordRegexConfigParamsGet, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectoryPasswordRegexOptions, "setCloudDirectoryPasswordRegexOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectoryPasswordRegexOptions, "setCloudDirectoryPasswordRegexOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectoryPasswordRegexOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/password_regex`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectoryPasswordRegexOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectoryPasswordRegex")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectoryPasswordRegexOptions.Regex != nil {
		body["regex"] = setCloudDirectoryPasswordRegexOptions.Regex
	}
	if setCloudDirectoryPasswordRegexOptions.Base64EncodedRegex != nil {
		body["base64_encoded_regex"] = setCloudDirectoryPasswordRegexOptions.Base64EncodedRegex
	}
	if setCloudDirectoryPasswordRegexOptions.ErrorMessage != nil {
		body["error_message"] = setCloudDirectoryPasswordRegexOptions.ErrorMessage
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPasswordRegexConfigParamsGet)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCloudDirectoryEmailDispatcher : Get email dispatcher configuration
// Get the configuration of email dispatcher that is used by Cloud Directory when sending emails.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryEmailDispatcher(getCloudDirectoryEmailDispatcherOptions *GetCloudDirectoryEmailDispatcherOptions) (result *EmailDispatcherParams, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryEmailDispatcherWithContext(context.Background(), getCloudDirectoryEmailDispatcherOptions)
}

// GetCloudDirectoryEmailDispatcherWithContext is an alternate form of the GetCloudDirectoryEmailDispatcher method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryEmailDispatcherWithContext(ctx context.Context, getCloudDirectoryEmailDispatcherOptions *GetCloudDirectoryEmailDispatcherOptions) (result *EmailDispatcherParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryEmailDispatcherOptions, "getCloudDirectoryEmailDispatcherOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryEmailDispatcherOptions, "getCloudDirectoryEmailDispatcherOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryEmailDispatcherOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/email_dispatcher`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryEmailDispatcherOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryEmailDispatcher")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEmailDispatcherParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectoryEmailDispatcher : Update email dispatcher configuration
// App ID allows you to use your own email provider. You can use your own Sendgrid account by providing your Sendgrind
// API key. Alternatively, you can define a custom email dispatcher by providing App ID with URL. The URL is called for
// sending emails. Optionally, you can determine a specific authorization method – either basic, such as a username and
// password, or a custom value. By default, App ID's email provider will be used.
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryEmailDispatcher(setCloudDirectoryEmailDispatcherOptions *SetCloudDirectoryEmailDispatcherOptions) (result *EmailDispatcherParams, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectoryEmailDispatcherWithContext(context.Background(), setCloudDirectoryEmailDispatcherOptions)
}

// SetCloudDirectoryEmailDispatcherWithContext is an alternate form of the SetCloudDirectoryEmailDispatcher method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryEmailDispatcherWithContext(ctx context.Context, setCloudDirectoryEmailDispatcherOptions *SetCloudDirectoryEmailDispatcherOptions) (result *EmailDispatcherParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectoryEmailDispatcherOptions, "setCloudDirectoryEmailDispatcherOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectoryEmailDispatcherOptions, "setCloudDirectoryEmailDispatcherOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectoryEmailDispatcherOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/email_dispatcher`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectoryEmailDispatcherOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectoryEmailDispatcher")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectoryEmailDispatcherOptions.Provider != nil {
		body["provider"] = setCloudDirectoryEmailDispatcherOptions.Provider
	}
	if setCloudDirectoryEmailDispatcherOptions.Sendgrid != nil {
		body["sendgrid"] = setCloudDirectoryEmailDispatcherOptions.Sendgrid
	}
	if setCloudDirectoryEmailDispatcherOptions.Custom != nil {
		body["custom"] = setCloudDirectoryEmailDispatcherOptions.Custom
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEmailDispatcherParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// EmailSettingTest : Test the email provider configuration
// You can send a message to a specific email to test your settings.
func (appIdManagement *AppIDManagementV4) EmailSettingTest(emailSettingTestOptions *EmailSettingTestOptions) (result *RespEmailSettingsTest, response *core.DetailedResponse, err error) {
	return appIdManagement.EmailSettingTestWithContext(context.Background(), emailSettingTestOptions)
}

// EmailSettingTestWithContext is an alternate form of the EmailSettingTest method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) EmailSettingTestWithContext(ctx context.Context, emailSettingTestOptions *EmailSettingTestOptions) (result *RespEmailSettingsTest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(emailSettingTestOptions, "emailSettingTestOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(emailSettingTestOptions, "emailSettingTestOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *emailSettingTestOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/email_settings/test`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range emailSettingTestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "EmailSettingTest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if emailSettingTestOptions.EmailTo != nil {
		body["emailTo"] = emailSettingTestOptions.EmailTo
	}
	if emailSettingTestOptions.EmailSettings != nil {
		body["emailSettings"] = emailSettingTestOptions.EmailSettings
	}
	if emailSettingTestOptions.SenderDetails != nil {
		body["senderDetails"] = emailSettingTestOptions.SenderDetails
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRespEmailSettingsTest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostEmailDispatcherTest : Test the email dispatcher configuration
// You can send a message to a specific email to test your configuration.
func (appIdManagement *AppIDManagementV4) PostEmailDispatcherTest(postEmailDispatcherTestOptions *PostEmailDispatcherTestOptions) (result *RespCustomEmailDisParams, response *core.DetailedResponse, err error) {
	return appIdManagement.PostEmailDispatcherTestWithContext(context.Background(), postEmailDispatcherTestOptions)
}

// PostEmailDispatcherTestWithContext is an alternate form of the PostEmailDispatcherTest method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostEmailDispatcherTestWithContext(ctx context.Context, postEmailDispatcherTestOptions *PostEmailDispatcherTestOptions) (result *RespCustomEmailDisParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postEmailDispatcherTestOptions, "postEmailDispatcherTestOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postEmailDispatcherTestOptions, "postEmailDispatcherTestOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postEmailDispatcherTestOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/email_dispatcher/test`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postEmailDispatcherTestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostEmailDispatcherTest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postEmailDispatcherTestOptions.Email != nil {
		body["email"] = postEmailDispatcherTestOptions.Email
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRespCustomEmailDisParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostSMSDispatcherTest : Test the MFA SMS dispatcher configuration
// You can send a message to a specific phone number to test your MFA SMS configuration.
func (appIdManagement *AppIDManagementV4) PostSMSDispatcherTest(postSMSDispatcherTestOptions *PostSMSDispatcherTestOptions) (result *RespSMSDisParams, response *core.DetailedResponse, err error) {
	return appIdManagement.PostSMSDispatcherTestWithContext(context.Background(), postSMSDispatcherTestOptions)
}

// PostSMSDispatcherTestWithContext is an alternate form of the PostSMSDispatcherTest method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostSMSDispatcherTestWithContext(ctx context.Context, postSMSDispatcherTestOptions *PostSMSDispatcherTestOptions) (result *RespSMSDisParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postSMSDispatcherTestOptions, "postSMSDispatcherTestOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postSMSDispatcherTestOptions, "postSMSDispatcherTestOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postSMSDispatcherTestOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/sms_dispatcher/test`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postSMSDispatcherTestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostSMSDispatcherTest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postSMSDispatcherTestOptions.PhoneNumber != nil {
		body["phone_number"] = postSMSDispatcherTestOptions.PhoneNumber
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRespSMSDisParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCloudDirectoryAdvancedPasswordManagement : Get APM configuration
// Get the configuration of the advanced password management.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryAdvancedPasswordManagement(getCloudDirectoryAdvancedPasswordManagementOptions *GetCloudDirectoryAdvancedPasswordManagementOptions) (result *ApmSchema, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryAdvancedPasswordManagementWithContext(context.Background(), getCloudDirectoryAdvancedPasswordManagementOptions)
}

// GetCloudDirectoryAdvancedPasswordManagementWithContext is an alternate form of the GetCloudDirectoryAdvancedPasswordManagement method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryAdvancedPasswordManagementWithContext(ctx context.Context, getCloudDirectoryAdvancedPasswordManagementOptions *GetCloudDirectoryAdvancedPasswordManagementOptions) (result *ApmSchema, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryAdvancedPasswordManagementOptions, "getCloudDirectoryAdvancedPasswordManagementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryAdvancedPasswordManagementOptions, "getCloudDirectoryAdvancedPasswordManagementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryAdvancedPasswordManagementOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/advanced_password_management`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryAdvancedPasswordManagementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryAdvancedPasswordManagement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApmSchema)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectoryAdvancedPasswordManagement : Update APM configuration
// Updates the advanced password management configuration for the provided tenantId. By turning this on, any
// authentication event is also charged as advanced security event.
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryAdvancedPasswordManagement(setCloudDirectoryAdvancedPasswordManagementOptions *SetCloudDirectoryAdvancedPasswordManagementOptions) (result *ApmSchema, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectoryAdvancedPasswordManagementWithContext(context.Background(), setCloudDirectoryAdvancedPasswordManagementOptions)
}

// SetCloudDirectoryAdvancedPasswordManagementWithContext is an alternate form of the SetCloudDirectoryAdvancedPasswordManagement method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryAdvancedPasswordManagementWithContext(ctx context.Context, setCloudDirectoryAdvancedPasswordManagementOptions *SetCloudDirectoryAdvancedPasswordManagementOptions) (result *ApmSchema, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectoryAdvancedPasswordManagementOptions, "setCloudDirectoryAdvancedPasswordManagementOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectoryAdvancedPasswordManagementOptions, "setCloudDirectoryAdvancedPasswordManagementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectoryAdvancedPasswordManagementOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/advanced_password_management`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectoryAdvancedPasswordManagementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectoryAdvancedPasswordManagement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectoryAdvancedPasswordManagementOptions.AdvancedPasswordManagement != nil {
		body["advancedPasswordManagement"] = setCloudDirectoryAdvancedPasswordManagementOptions.AdvancedPasswordManagement
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApmSchema)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAuditStatus : Get tenant audit status
// Returns a JSON object containing the auditing status of the tenant.
func (appIdManagement *AppIDManagementV4) GetAuditStatus(getAuditStatusOptions *GetAuditStatusOptions) (result *GetAuditStatusResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetAuditStatusWithContext(context.Background(), getAuditStatusOptions)
}

// GetAuditStatusWithContext is an alternate form of the GetAuditStatus method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetAuditStatusWithContext(ctx context.Context, getAuditStatusOptions *GetAuditStatusOptions) (result *GetAuditStatusResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAuditStatusOptions, "getAuditStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAuditStatusOptions, "getAuditStatusOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getAuditStatusOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/capture_runtime_activity`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAuditStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetAuditStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetAuditStatusResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetAuditStatus : Update tenant audit status
// Capture app user sign-in, sign-up and other runtime events in Activity Tracker for you to search, analyse and report.
// By turning this On, any authentication event is also charged as advanced security event. Activity Tracker with LogDNA
// is available in select regions. <a href="https://cloud.ibm.com/docs/appid?topic=appid-at-events">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetAuditStatus(setAuditStatusOptions *SetAuditStatusOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.SetAuditStatusWithContext(context.Background(), setAuditStatusOptions)
}

// SetAuditStatusWithContext is an alternate form of the SetAuditStatus method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetAuditStatusWithContext(ctx context.Context, setAuditStatusOptions *SetAuditStatusOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setAuditStatusOptions, "setAuditStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setAuditStatusOptions, "setAuditStatusOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setAuditStatusOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/capture_runtime_activity`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setAuditStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetAuditStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setAuditStatusOptions.IsActive != nil {
		body["isActive"] = setAuditStatusOptions.IsActive
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// ListChannels : List channels
// Returns all MFA channels registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) ListChannels(listChannelsOptions *ListChannelsOptions) (result *MFAChannelsList, response *core.DetailedResponse, err error) {
	return appIdManagement.ListChannelsWithContext(context.Background(), listChannelsOptions)
}

// ListChannelsWithContext is an alternate form of the ListChannels method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ListChannelsWithContext(ctx context.Context, listChannelsOptions *ListChannelsOptions) (result *MFAChannelsList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listChannelsOptions, "listChannelsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listChannelsOptions, "listChannelsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *listChannelsOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/channels`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listChannelsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ListChannels")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMFAChannelsList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetChannel : Get channel
// Returns a specific MFA channel registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) GetChannel(getChannelOptions *GetChannelOptions) (result *GetSMSChannel, response *core.DetailedResponse, err error) {
	return appIdManagement.GetChannelWithContext(context.Background(), getChannelOptions)
}

// GetChannelWithContext is an alternate form of the GetChannel method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetChannelWithContext(ctx context.Context, getChannelOptions *GetChannelOptions) (result *GetSMSChannel, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getChannelOptions, "getChannelOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getChannelOptions, "getChannelOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getChannelOptions.TenantID,
		"channel":  *getChannelOptions.Channel,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/channels/{channel}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getChannelOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetChannel")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSMSChannel)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateChannel : Update channel
// Enable or disable a registered MFA channel on the App ID instance.
func (appIdManagement *AppIDManagementV4) UpdateChannel(updateChannelOptions *UpdateChannelOptions) (result *GetSMSChannel, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateChannelWithContext(context.Background(), updateChannelOptions)
}

// UpdateChannelWithContext is an alternate form of the UpdateChannel method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateChannelWithContext(ctx context.Context, updateChannelOptions *UpdateChannelOptions) (result *GetSMSChannel, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateChannelOptions, "updateChannelOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateChannelOptions, "updateChannelOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateChannelOptions.TenantID,
		"channel":  *updateChannelOptions.Channel,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/channels/{channel}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateChannelOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateChannel")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateChannelOptions.IsActive != nil {
		body["isActive"] = updateChannelOptions.IsActive
	}
	if updateChannelOptions.Config != nil {
		body["config"] = updateChannelOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetSMSChannel)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetExtensionConfig : Get an extension configuration
// View a registered extension's configuration for an instance of App ID. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-mfa#cd-mfa-extensions" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetExtensionConfig(getExtensionConfigOptions *GetExtensionConfigOptions) (result *UpdateExtensionConfig, response *core.DetailedResponse, err error) {
	return appIdManagement.GetExtensionConfigWithContext(context.Background(), getExtensionConfigOptions)
}

// GetExtensionConfigWithContext is an alternate form of the GetExtensionConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetExtensionConfigWithContext(ctx context.Context, getExtensionConfigOptions *GetExtensionConfigOptions) (result *UpdateExtensionConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getExtensionConfigOptions, "getExtensionConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getExtensionConfigOptions, "getExtensionConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getExtensionConfigOptions.TenantID,
		"name":     *getExtensionConfigOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/extensions/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getExtensionConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetExtensionConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdateExtensionConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateExtensionConfig : Update an extension configuration
// Set or update a registered extension's configuration for an instance of App ID. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-mfa#cd-mfa-extensions" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateExtensionConfig(updateExtensionConfigOptions *UpdateExtensionConfigOptions) (result *UpdateExtensionConfig, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateExtensionConfigWithContext(context.Background(), updateExtensionConfigOptions)
}

// UpdateExtensionConfigWithContext is an alternate form of the UpdateExtensionConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateExtensionConfigWithContext(ctx context.Context, updateExtensionConfigOptions *UpdateExtensionConfigOptions) (result *UpdateExtensionConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateExtensionConfigOptions, "updateExtensionConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateExtensionConfigOptions, "updateExtensionConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateExtensionConfigOptions.TenantID,
		"name":     *updateExtensionConfigOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/extensions/{name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateExtensionConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateExtensionConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateExtensionConfigOptions.IsActive != nil {
		body["isActive"] = updateExtensionConfigOptions.IsActive
	}
	if updateExtensionConfigOptions.Config != nil {
		body["config"] = updateExtensionConfigOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdateExtensionConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateExtensionActive : Enable or disable an extension
// Update the status of a registered extension for an instance of App ID to enabled or disabled. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cd-mfa#cd-mfa-extensions" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UpdateExtensionActive(updateExtensionActiveOptions *UpdateExtensionActiveOptions) (result *ExtensionActive, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateExtensionActiveWithContext(context.Background(), updateExtensionActiveOptions)
}

// UpdateExtensionActiveWithContext is an alternate form of the UpdateExtensionActive method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateExtensionActiveWithContext(ctx context.Context, updateExtensionActiveOptions *UpdateExtensionActiveOptions) (result *ExtensionActive, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateExtensionActiveOptions, "updateExtensionActiveOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateExtensionActiveOptions, "updateExtensionActiveOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateExtensionActiveOptions.TenantID,
		"name":     *updateExtensionActiveOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/extensions/{name}/active`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateExtensionActiveOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateExtensionActive")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateExtensionActiveOptions.IsActive != nil {
		body["isActive"] = updateExtensionActiveOptions.IsActive
	}
	if updateExtensionActiveOptions.Config != nil {
		body["config"] = updateExtensionActiveOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExtensionActive)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostExtensionsTest : Test the extension configuration
// Test an extension configuration. <a href="https://cloud.ibm.com/docs/appid?topic=appid-cd-mfa#cd-mfa-extensions"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) PostExtensionsTest(postExtensionsTestOptions *PostExtensionsTestOptions) (result *ExtensionTest, response *core.DetailedResponse, err error) {
	return appIdManagement.PostExtensionsTestWithContext(context.Background(), postExtensionsTestOptions)
}

// PostExtensionsTestWithContext is an alternate form of the PostExtensionsTest method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) PostExtensionsTestWithContext(ctx context.Context, postExtensionsTestOptions *PostExtensionsTestOptions) (result *ExtensionTest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postExtensionsTestOptions, "postExtensionsTestOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postExtensionsTestOptions, "postExtensionsTestOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *postExtensionsTestOptions.TenantID,
		"name":     *postExtensionsTestOptions.Name,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa/extensions/{name}/test`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postExtensionsTestOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "PostExtensionsTest")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExtensionTest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetMFAConfig : Get MFA configuration
// Returns MFA configuration registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) GetMFAConfig(getMFAConfigOptions *GetMFAConfigOptions) (result *GetMFAConfiguration, response *core.DetailedResponse, err error) {
	return appIdManagement.GetMFAConfigWithContext(context.Background(), getMFAConfigOptions)
}

// GetMFAConfigWithContext is an alternate form of the GetMFAConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetMFAConfigWithContext(ctx context.Context, getMFAConfigOptions *GetMFAConfigOptions) (result *GetMFAConfiguration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getMFAConfigOptions, "getMFAConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getMFAConfigOptions, "getMFAConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getMFAConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getMFAConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetMFAConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetMFAConfiguration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateMFAConfig : Update MFA configuration
// Update MFA configuration on the App ID instance.
func (appIdManagement *AppIDManagementV4) UpdateMFAConfig(updateMFAConfigOptions *UpdateMFAConfigOptions) (result *GetMFAConfiguration, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateMFAConfigWithContext(context.Background(), updateMFAConfigOptions)
}

// UpdateMFAConfigWithContext is an alternate form of the UpdateMFAConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateMFAConfigWithContext(ctx context.Context, updateMFAConfigOptions *UpdateMFAConfigOptions) (result *GetMFAConfiguration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateMFAConfigOptions, "updateMFAConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateMFAConfigOptions, "updateMFAConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateMFAConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/mfa`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateMFAConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateMFAConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateMFAConfigOptions.IsActive != nil {
		body["isActive"] = updateMFAConfigOptions.IsActive
	}
	if updateMFAConfigOptions.Config != nil {
		body["config"] = updateMFAConfigOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetMFAConfiguration)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSSOConfig : Get SSO configuration
// Returns SSO configuration registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) GetSSOConfig(getSSOConfigOptions *GetSSOConfigOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.GetSSOConfigWithContext(context.Background(), getSSOConfigOptions)
}

// GetSSOConfigWithContext is an alternate form of the GetSSOConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetSSOConfigWithContext(ctx context.Context, getSSOConfigOptions *GetSSOConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSSOConfigOptions, "getSSOConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSSOConfigOptions, "getSSOConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getSSOConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/sso`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSSOConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetSSOConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UpdateSSOConfig : Update SSO configuration
// Update SSO configuration on the App ID instance.
func (appIdManagement *AppIDManagementV4) UpdateSSOConfig(updateSSOConfigOptions *UpdateSSOConfigOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateSSOConfigWithContext(context.Background(), updateSSOConfigOptions)
}

// UpdateSSOConfigWithContext is an alternate form of the UpdateSSOConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateSSOConfigWithContext(ctx context.Context, updateSSOConfigOptions *UpdateSSOConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSSOConfigOptions, "updateSSOConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSSOConfigOptions, "updateSSOConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateSSOConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/sso`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSSOConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateSSOConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSSOConfigOptions.IsActive != nil {
		body["isActive"] = updateSSOConfigOptions.IsActive
	}
	if updateSSOConfigOptions.InactivityTimeoutSeconds != nil {
		body["inactivityTimeoutSeconds"] = updateSSOConfigOptions.InactivityTimeoutSeconds
	}
	if updateSSOConfigOptions.LogoutRedirectUris != nil {
		body["logoutRedirectUris"] = updateSSOConfigOptions.LogoutRedirectUris
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetRateLimitConfig : Get the rate limit configuration
// Returns the rate limit configuration registered with the App ID Instance.
func (appIdManagement *AppIDManagementV4) GetRateLimitConfig(getRateLimitConfigOptions *GetRateLimitConfigOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.GetRateLimitConfigWithContext(context.Background(), getRateLimitConfigOptions)
}

// GetRateLimitConfigWithContext is an alternate form of the GetRateLimitConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetRateLimitConfigWithContext(ctx context.Context, getRateLimitConfigOptions *GetRateLimitConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRateLimitConfigOptions, "getRateLimitConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRateLimitConfigOptions, "getRateLimitConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getRateLimitConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/rate_limit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRateLimitConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetRateLimitConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UpdateRateLimitConfig : Update the rate limit configuration
// Update the rate limit configuration on the App ID instance.
func (appIdManagement *AppIDManagementV4) UpdateRateLimitConfig(updateRateLimitConfigOptions *UpdateRateLimitConfigOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateRateLimitConfigWithContext(context.Background(), updateRateLimitConfigOptions)
}

// UpdateRateLimitConfigWithContext is an alternate form of the UpdateRateLimitConfig method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateRateLimitConfigWithContext(ctx context.Context, updateRateLimitConfigOptions *UpdateRateLimitConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRateLimitConfigOptions, "updateRateLimitConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRateLimitConfigOptions, "updateRateLimitConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateRateLimitConfigOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/cloud_directory/rate_limit`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRateLimitConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateRateLimitConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateRateLimitConfigOptions.SignUpLimitPerMinute != nil {
		body["signUpLimitPerMinute"] = updateRateLimitConfigOptions.SignUpLimitPerMinute
	}
	if updateRateLimitConfigOptions.SignInLimitPerMinute != nil {
		body["signInLimitPerMinute"] = updateRateLimitConfigOptions.SignInLimitPerMinute
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetFacebookIDP : Get Facebook IDP configuration
// Returns the Facebook identity provider configuration.
func (appIdManagement *AppIDManagementV4) GetFacebookIDP(getFacebookIDPOptions *GetFacebookIDPOptions) (result *FacebookConfigParams, response *core.DetailedResponse, err error) {
	return appIdManagement.GetFacebookIDPWithContext(context.Background(), getFacebookIDPOptions)
}

// GetFacebookIDPWithContext is an alternate form of the GetFacebookIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetFacebookIDPWithContext(ctx context.Context, getFacebookIDPOptions *GetFacebookIDPOptions) (result *FacebookConfigParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getFacebookIDPOptions, "getFacebookIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getFacebookIDPOptions, "getFacebookIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getFacebookIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/facebook`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getFacebookIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetFacebookIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFacebookConfigParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetFacebookIDP : Update Facebook IDP configuration
// Configure Facebook to set up a single sign-on experience for your users. By using Facebook, users are able to sign in
// with credentials with which they are already familiar. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-social#facebook" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetFacebookIDP(setFacebookIDPOptions *SetFacebookIDPOptions) (result *FacebookConfigParamsPut, response *core.DetailedResponse, err error) {
	return appIdManagement.SetFacebookIDPWithContext(context.Background(), setFacebookIDPOptions)
}

// SetFacebookIDPWithContext is an alternate form of the SetFacebookIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetFacebookIDPWithContext(ctx context.Context, setFacebookIDPOptions *SetFacebookIDPOptions) (result *FacebookConfigParamsPut, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setFacebookIDPOptions, "setFacebookIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setFacebookIDPOptions, "setFacebookIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setFacebookIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/facebook`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setFacebookIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetFacebookIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(setFacebookIDPOptions.IDP)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFacebookConfigParamsPut)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetGoogleIDP : Get Google IDP configuration
// Returns the Google identity provider configuration.
func (appIdManagement *AppIDManagementV4) GetGoogleIDP(getGoogleIDPOptions *GetGoogleIDPOptions) (result *GoogleConfigParams, response *core.DetailedResponse, err error) {
	return appIdManagement.GetGoogleIDPWithContext(context.Background(), getGoogleIDPOptions)
}

// GetGoogleIDPWithContext is an alternate form of the GetGoogleIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetGoogleIDPWithContext(ctx context.Context, getGoogleIDPOptions *GetGoogleIDPOptions) (result *GoogleConfigParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGoogleIDPOptions, "getGoogleIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getGoogleIDPOptions, "getGoogleIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getGoogleIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/google`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getGoogleIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetGoogleIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGoogleConfigParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetGoogleIDP : Update Google IDP configuration
// Configure Google to set up a single sign-on experience for your users. By using Google, users are able to sign in
// with credentials with which they are already familiar. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-social#google" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetGoogleIDP(setGoogleIDPOptions *SetGoogleIDPOptions) (result *GoogleConfigParamsPut, response *core.DetailedResponse, err error) {
	return appIdManagement.SetGoogleIDPWithContext(context.Background(), setGoogleIDPOptions)
}

// SetGoogleIDPWithContext is an alternate form of the SetGoogleIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetGoogleIDPWithContext(ctx context.Context, setGoogleIDPOptions *SetGoogleIDPOptions) (result *GoogleConfigParamsPut, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setGoogleIDPOptions, "setGoogleIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setGoogleIDPOptions, "setGoogleIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setGoogleIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/google`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setGoogleIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetGoogleIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(setGoogleIDPOptions.IDP)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGoogleConfigParamsPut)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCustomIDP : Returns the Custom identity configuration
func (appIdManagement *AppIDManagementV4) GetCustomIDP(getCustomIDPOptions *GetCustomIDPOptions) (result *CustomIDPConfigParams, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCustomIDPWithContext(context.Background(), getCustomIDPOptions)
}

// GetCustomIDPWithContext is an alternate form of the GetCustomIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCustomIDPWithContext(ctx context.Context, getCustomIDPOptions *GetCustomIDPOptions) (result *CustomIDPConfigParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCustomIDPOptions, "getCustomIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCustomIDPOptions, "getCustomIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCustomIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/custom`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCustomIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCustomIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomIDPConfigParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCustomIDP : Update or change the configuration of the Custom identity
// Configure App ID Custom identity to allow users to sign-in using your own identity provider.
func (appIdManagement *AppIDManagementV4) SetCustomIDP(setCustomIDPOptions *SetCustomIDPOptions) (result *CustomIDPConfigParams, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCustomIDPWithContext(context.Background(), setCustomIDPOptions)
}

// SetCustomIDPWithContext is an alternate form of the SetCustomIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCustomIDPWithContext(ctx context.Context, setCustomIDPOptions *SetCustomIDPOptions) (result *CustomIDPConfigParams, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCustomIDPOptions, "setCustomIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCustomIDPOptions, "setCustomIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCustomIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/custom`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCustomIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCustomIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCustomIDPOptions.IsActive != nil {
		body["isActive"] = setCustomIDPOptions.IsActive
	}
	if setCustomIDPOptions.Config != nil {
		body["config"] = setCustomIDPOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomIDPConfigParams)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCloudDirectoryIDP : Get Cloud Directory IDP configuration
// Returns the Cloud Directory identity provider configuration. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryIDP(getCloudDirectoryIDPOptions *GetCloudDirectoryIDPOptions) (result *CloudDirectoryResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetCloudDirectoryIDPWithContext(context.Background(), getCloudDirectoryIDPOptions)
}

// GetCloudDirectoryIDPWithContext is an alternate form of the GetCloudDirectoryIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetCloudDirectoryIDPWithContext(ctx context.Context, getCloudDirectoryIDPOptions *GetCloudDirectoryIDPOptions) (result *CloudDirectoryResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCloudDirectoryIDPOptions, "getCloudDirectoryIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCloudDirectoryIDPOptions, "getCloudDirectoryIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getCloudDirectoryIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/cloud_directory`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCloudDirectoryIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetCloudDirectoryIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCloudDirectoryResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetCloudDirectoryIDP : Update Cloud Directory IDP configuration
// Configure Cloud Directory to set up a single sign-on experience for your users. With Cloud Directory users can use
// their email and a password of their choice to log in to your applications. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-cloud-directory" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryIDP(setCloudDirectoryIDPOptions *SetCloudDirectoryIDPOptions) (result *CloudDirectoryResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.SetCloudDirectoryIDPWithContext(context.Background(), setCloudDirectoryIDPOptions)
}

// SetCloudDirectoryIDPWithContext is an alternate form of the SetCloudDirectoryIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetCloudDirectoryIDPWithContext(ctx context.Context, setCloudDirectoryIDPOptions *SetCloudDirectoryIDPOptions) (result *CloudDirectoryResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCloudDirectoryIDPOptions, "setCloudDirectoryIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setCloudDirectoryIDPOptions, "setCloudDirectoryIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setCloudDirectoryIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/cloud_directory`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setCloudDirectoryIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetCloudDirectoryIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setCloudDirectoryIDPOptions.IsActive != nil {
		body["isActive"] = setCloudDirectoryIDPOptions.IsActive
	}
	if setCloudDirectoryIDPOptions.Config != nil {
		body["config"] = setCloudDirectoryIDPOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCloudDirectoryResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSAMLIDP : Get SAML IDP configuration
// Returns the SAML identity provider configuration, including status and credentials. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-enterprise" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) GetSAMLIDP(getSAMLIDPOptions *GetSAMLIDPOptions) (result *SAMLResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetSAMLIDPWithContext(context.Background(), getSAMLIDPOptions)
}

// GetSAMLIDPWithContext is an alternate form of the GetSAMLIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetSAMLIDPWithContext(ctx context.Context, getSAMLIDPOptions *GetSAMLIDPOptions) (result *SAMLResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSAMLIDPOptions, "getSAMLIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSAMLIDPOptions, "getSAMLIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getSAMLIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/saml`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSAMLIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetSAMLIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSAMLResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetSAMLIDP : Update SAML IDP configuration
// Configure SAML to set up a single sign-on experience for your users. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-enterprise" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) SetSAMLIDP(setSAMLIDPOptions *SetSAMLIDPOptions) (result *SAMLResponseWithValidationData, response *core.DetailedResponse, err error) {
	return appIdManagement.SetSAMLIDPWithContext(context.Background(), setSAMLIDPOptions)
}

// SetSAMLIDPWithContext is an alternate form of the SetSAMLIDP method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) SetSAMLIDPWithContext(ctx context.Context, setSAMLIDPOptions *SetSAMLIDPOptions) (result *SAMLResponseWithValidationData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setSAMLIDPOptions, "setSAMLIDPOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(setSAMLIDPOptions, "setSAMLIDPOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *setSAMLIDPOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/config/idps/saml`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setSAMLIDPOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "SetSAMLIDP")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if setSAMLIDPOptions.IsActive != nil {
		body["isActive"] = setSAMLIDPOptions.IsActive
	}
	if setSAMLIDPOptions.Config != nil {
		body["config"] = setSAMLIDPOptions.Config
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSAMLResponseWithValidationData)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRoles : List all roles
// Obtain a list of the roles that are associated with your registered application.
func (appIdManagement *AppIDManagementV4) ListRoles(listRolesOptions *ListRolesOptions) (result *RolesList, response *core.DetailedResponse, err error) {
	return appIdManagement.ListRolesWithContext(context.Background(), listRolesOptions)
}

// ListRolesWithContext is an alternate form of the ListRoles method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) ListRolesWithContext(ctx context.Context, listRolesOptions *ListRolesOptions) (result *RolesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRolesOptions, "listRolesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listRolesOptions, "listRolesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *listRolesOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "ListRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRolesList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateRole : Create a role
// Create a new role for a registered application.
func (appIdManagement *AppIDManagementV4) CreateRole(createRoleOptions *CreateRoleOptions) (result *CreateRoleResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.CreateRoleWithContext(context.Background(), createRoleOptions)
}

// CreateRoleWithContext is an alternate form of the CreateRole method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) CreateRoleWithContext(ctx context.Context, createRoleOptions *CreateRoleOptions) (result *CreateRoleResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRoleOptions, "createRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRoleOptions, "createRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *createRoleOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "CreateRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRoleOptions.Name != nil {
		body["name"] = createRoleOptions.Name
	}
	if createRoleOptions.Access != nil {
		body["access"] = createRoleOptions.Access
	}
	if createRoleOptions.Description != nil {
		body["description"] = createRoleOptions.Description
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateRoleResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRole : View a specific role
// By using the role ID, obtain the information for a specific role that is associated with a registered application.
func (appIdManagement *AppIDManagementV4) GetRole(getRoleOptions *GetRoleOptions) (result *GetRoleResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetRoleWithContext(context.Background(), getRoleOptions)
}

// GetRoleWithContext is an alternate form of the GetRole method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetRoleWithContext(ctx context.Context, getRoleOptions *GetRoleOptions) (result *GetRoleResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRoleOptions, "getRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRoleOptions, "getRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getRoleOptions.TenantID,
		"roleId":   *getRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/roles/{roleId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetRoleResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateRole : Update a role
// Update an existing role.
func (appIdManagement *AppIDManagementV4) UpdateRole(updateRoleOptions *UpdateRoleOptions) (result *UpdateRoleResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateRoleWithContext(context.Background(), updateRoleOptions)
}

// UpdateRoleWithContext is an alternate form of the UpdateRole method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateRoleWithContext(ctx context.Context, updateRoleOptions *UpdateRoleOptions) (result *UpdateRoleResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRoleOptions, "updateRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRoleOptions, "updateRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateRoleOptions.TenantID,
		"roleId":   *updateRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/roles/{roleId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateRoleOptions.Name != nil {
		body["name"] = updateRoleOptions.Name
	}
	if updateRoleOptions.Access != nil {
		body["access"] = updateRoleOptions.Access
	}
	if updateRoleOptions.Description != nil {
		body["description"] = updateRoleOptions.Description
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdateRoleResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRole : Delete a role
// Delete an existing role.
func (appIdManagement *AppIDManagementV4) DeleteRole(deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.DeleteRoleWithContext(context.Background(), deleteRoleOptions)
}

// DeleteRoleWithContext is an alternate form of the DeleteRole method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) DeleteRoleWithContext(ctx context.Context, deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRoleOptions, "deleteRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRoleOptions, "deleteRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *deleteRoleOptions.TenantID,
		"roleId":   *deleteRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/roles/{roleId}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "DeleteRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UsersSearchUserProfile : Search users
// Returns list of users, if given email/id returns only users which match the email/id - not including anonymous
// profiles. <a href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UsersSearchUserProfile(usersSearchUserProfileOptions *UsersSearchUserProfileOptions) (result *UserSearchResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.UsersSearchUserProfileWithContext(context.Background(), usersSearchUserProfileOptions)
}

// UsersSearchUserProfileWithContext is an alternate form of the UsersSearchUserProfile method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersSearchUserProfileWithContext(ctx context.Context, usersSearchUserProfileOptions *UsersSearchUserProfileOptions) (result *UserSearchResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersSearchUserProfileOptions, "usersSearchUserProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersSearchUserProfileOptions, "usersSearchUserProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersSearchUserProfileOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersSearchUserProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersSearchUserProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("dataScope", fmt.Sprint(*usersSearchUserProfileOptions.DataScope))
	if usersSearchUserProfileOptions.Email != nil {
		builder.AddQuery("email", fmt.Sprint(*usersSearchUserProfileOptions.Email))
	}
	if usersSearchUserProfileOptions.ID != nil {
		builder.AddQuery("id", fmt.Sprint(*usersSearchUserProfileOptions.ID))
	}
	if usersSearchUserProfileOptions.StartIndex != nil {
		builder.AddQuery("startIndex", fmt.Sprint(*usersSearchUserProfileOptions.StartIndex))
	}
	if usersSearchUserProfileOptions.Count != nil {
		builder.AddQuery("count", fmt.Sprint(*usersSearchUserProfileOptions.Count))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserSearchResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UsersNominateUser : Pre-register a user profile
// Create a profile for a user that you know needs access to your app before they sign in to your app for the first
// time. <a href="https://cloud.ibm.com/docs/appid?topic=appid-preregister" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UsersNominateUser(usersNominateUserOptions *UsersNominateUserOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UsersNominateUserWithContext(context.Background(), usersNominateUserOptions)
}

// UsersNominateUserWithContext is an alternate form of the UsersNominateUser method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersNominateUserWithContext(ctx context.Context, usersNominateUserOptions *UsersNominateUserOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersNominateUserOptions, "usersNominateUserOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersNominateUserOptions, "usersNominateUserOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersNominateUserOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersNominateUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersNominateUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if usersNominateUserOptions.IDP != nil {
		body["idp"] = usersNominateUserOptions.IDP
	}
	if usersNominateUserOptions.IDPIdentity != nil {
		body["idp-identity"] = usersNominateUserOptions.IDPIdentity
	}
	if usersNominateUserOptions.Profile != nil {
		body["profile"] = usersNominateUserOptions.Profile
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UserProfilesExport : Export user profiles
// Exports App ID user profiles, not including Cloud Directory and anonymous users.
func (appIdManagement *AppIDManagementV4) UserProfilesExport(userProfilesExportOptions *UserProfilesExportOptions) (result *ExportUserProfile, response *core.DetailedResponse, err error) {
	return appIdManagement.UserProfilesExportWithContext(context.Background(), userProfilesExportOptions)
}

// UserProfilesExportWithContext is an alternate form of the UserProfilesExport method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UserProfilesExportWithContext(ctx context.Context, userProfilesExportOptions *UserProfilesExportOptions) (result *ExportUserProfile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(userProfilesExportOptions, "userProfilesExportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(userProfilesExportOptions, "userProfilesExportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *userProfilesExportOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/export`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range userProfilesExportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UserProfilesExport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if userProfilesExportOptions.StartIndex != nil {
		builder.AddQuery("startIndex", fmt.Sprint(*userProfilesExportOptions.StartIndex))
	}
	if userProfilesExportOptions.Count != nil {
		builder.AddQuery("count", fmt.Sprint(*userProfilesExportOptions.Count))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExportUserProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UserProfilesImport : Import user profiles
// Imports App ID user profiles, not including Cloud Directory and anonymous users.
func (appIdManagement *AppIDManagementV4) UserProfilesImport(userProfilesImportOptions *UserProfilesImportOptions) (result *ImportProfilesResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.UserProfilesImportWithContext(context.Background(), userProfilesImportOptions)
}

// UserProfilesImportWithContext is an alternate form of the UserProfilesImport method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UserProfilesImportWithContext(ctx context.Context, userProfilesImportOptions *UserProfilesImportOptions) (result *ImportProfilesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(userProfilesImportOptions, "userProfilesImportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(userProfilesImportOptions, "userProfilesImportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *userProfilesImportOptions.TenantID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/import`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range userProfilesImportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UserProfilesImport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if userProfilesImportOptions.Users != nil {
		body["users"] = userProfilesImportOptions.Users
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImportProfilesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UsersDeleteUserProfile : Delete user
// Deletes a user by id. <a href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UsersDeleteUserProfile(usersDeleteUserProfileOptions *UsersDeleteUserProfileOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UsersDeleteUserProfileWithContext(context.Background(), usersDeleteUserProfileOptions)
}

// UsersDeleteUserProfileWithContext is an alternate form of the UsersDeleteUserProfile method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersDeleteUserProfileWithContext(ctx context.Context, usersDeleteUserProfileOptions *UsersDeleteUserProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersDeleteUserProfileOptions, "usersDeleteUserProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersDeleteUserProfileOptions, "usersDeleteUserProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersDeleteUserProfileOptions.TenantID,
		"id":       *usersDeleteUserProfileOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersDeleteUserProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersDeleteUserProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UsersRevokeRefreshToken : Revoke refresh token
// Revokes all the refresh tokens issued for the given user. <a
// href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UsersRevokeRefreshToken(usersRevokeRefreshTokenOptions *UsersRevokeRefreshTokenOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UsersRevokeRefreshTokenWithContext(context.Background(), usersRevokeRefreshTokenOptions)
}

// UsersRevokeRefreshTokenWithContext is an alternate form of the UsersRevokeRefreshToken method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersRevokeRefreshTokenWithContext(ctx context.Context, usersRevokeRefreshTokenOptions *UsersRevokeRefreshTokenOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersRevokeRefreshTokenOptions, "usersRevokeRefreshTokenOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersRevokeRefreshTokenOptions, "usersRevokeRefreshTokenOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersRevokeRefreshTokenOptions.TenantID,
		"id":       *usersRevokeRefreshTokenOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}/revoke_refresh_token`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersRevokeRefreshTokenOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersRevokeRefreshToken")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UsersGetUserProfile : Get user profile
// Returns the profile of a given user. <a href="https://cloud.ibm.com/docs/appid?topic=appid-profiles"
// target="_blank">Learn more</a>.
func (appIdManagement *AppIDManagementV4) UsersGetUserProfile(usersGetUserProfileOptions *UsersGetUserProfileOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UsersGetUserProfileWithContext(context.Background(), usersGetUserProfileOptions)
}

// UsersGetUserProfileWithContext is an alternate form of the UsersGetUserProfile method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersGetUserProfileWithContext(ctx context.Context, usersGetUserProfileOptions *UsersGetUserProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersGetUserProfileOptions, "usersGetUserProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersGetUserProfileOptions, "usersGetUserProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersGetUserProfileOptions.TenantID,
		"id":       *usersGetUserProfileOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}/profile`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersGetUserProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersGetUserProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// UsersSetUserProfile : Update user profile
// Updates a user profile. <a href="https://cloud.ibm.com/docs/appid?topic=appid-profiles" target="_blank">Learn
// more</a>.
func (appIdManagement *AppIDManagementV4) UsersSetUserProfile(usersSetUserProfileOptions *UsersSetUserProfileOptions) (response *core.DetailedResponse, err error) {
	return appIdManagement.UsersSetUserProfileWithContext(context.Background(), usersSetUserProfileOptions)
}

// UsersSetUserProfileWithContext is an alternate form of the UsersSetUserProfile method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UsersSetUserProfileWithContext(ctx context.Context, usersSetUserProfileOptions *UsersSetUserProfileOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(usersSetUserProfileOptions, "usersSetUserProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(usersSetUserProfileOptions, "usersSetUserProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *usersSetUserProfileOptions.TenantID,
		"id":       *usersSetUserProfileOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}/profile`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range usersSetUserProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UsersSetUserProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if usersSetUserProfileOptions.Attributes != nil {
		body["attributes"] = usersSetUserProfileOptions.Attributes
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = appIdManagement.Service.Request(request, nil)

	return
}

// GetUserRoles : Get a user's roles
// View a list of roles that are associated with a specific user.
func (appIdManagement *AppIDManagementV4) GetUserRoles(getUserRolesOptions *GetUserRolesOptions) (result *GetUserRolesResponse, response *core.DetailedResponse, err error) {
	return appIdManagement.GetUserRolesWithContext(context.Background(), getUserRolesOptions)
}

// GetUserRolesWithContext is an alternate form of the GetUserRoles method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) GetUserRolesWithContext(ctx context.Context, getUserRolesOptions *GetUserRolesOptions) (result *GetUserRolesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getUserRolesOptions, "getUserRolesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getUserRolesOptions, "getUserRolesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *getUserRolesOptions.TenantID,
		"id":       *getUserRolesOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUserRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "GetUserRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetUserRolesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateUserRoles : Update a user's roles
// Update which roles are associated with a specific user or assign a role to a user for the first time.
func (appIdManagement *AppIDManagementV4) UpdateUserRoles(updateUserRolesOptions *UpdateUserRolesOptions) (result *AssignRoleToUser, response *core.DetailedResponse, err error) {
	return appIdManagement.UpdateUserRolesWithContext(context.Background(), updateUserRolesOptions)
}

// UpdateUserRolesWithContext is an alternate form of the UpdateUserRoles method which supports a Context parameter
func (appIdManagement *AppIDManagementV4) UpdateUserRolesWithContext(ctx context.Context, updateUserRolesOptions *UpdateUserRolesOptions) (result *AssignRoleToUser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateUserRolesOptions, "updateUserRolesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateUserRolesOptions, "updateUserRolesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"tenantId": *updateUserRolesOptions.TenantID,
		"id":       *updateUserRolesOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appIdManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appIdManagement.Service.Options.URL, `/management/v4/{tenantId}/users/{id}/roles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateUserRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_id_management", "V4", "UpdateUserRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateUserRolesOptions.Roles != nil {
		body["roles"] = updateUserRolesOptions.Roles
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
	response, err = appIdManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAssignRoleToUser)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ApmSchemaAdvancedPasswordManagement : ApmSchemaAdvancedPasswordManagement struct
type ApmSchemaAdvancedPasswordManagement struct {
	Enabled *bool `json:"enabled" validate:"required"`

	PasswordReuse *ApmSchemaAdvancedPasswordManagementPasswordReuse `json:"passwordReuse" validate:"required"`

	PreventPasswordWithUsername *ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername `json:"preventPasswordWithUsername" validate:"required"`

	PasswordExpiration *ApmSchemaAdvancedPasswordManagementPasswordExpiration `json:"passwordExpiration" validate:"required"`

	LockOutPolicy *ApmSchemaAdvancedPasswordManagementLockOutPolicy `json:"lockOutPolicy" validate:"required"`

	MinPasswordChangeInterval *ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval `json:"minPasswordChangeInterval,omitempty"`
}

// NewApmSchemaAdvancedPasswordManagement : Instantiate ApmSchemaAdvancedPasswordManagement (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagement(enabled bool, passwordReuse *ApmSchemaAdvancedPasswordManagementPasswordReuse, preventPasswordWithUsername *ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername, passwordExpiration *ApmSchemaAdvancedPasswordManagementPasswordExpiration, lockOutPolicy *ApmSchemaAdvancedPasswordManagementLockOutPolicy) (model *ApmSchemaAdvancedPasswordManagement, err error) {
	model = &ApmSchemaAdvancedPasswordManagement{
		Enabled:                     core.BoolPtr(enabled),
		PasswordReuse:               passwordReuse,
		PreventPasswordWithUsername: preventPasswordWithUsername,
		PasswordExpiration:          passwordExpiration,
		LockOutPolicy:               lockOutPolicy,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagement unmarshals an instance of ApmSchemaAdvancedPasswordManagement from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagement)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "passwordReuse", &obj.PasswordReuse, UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "preventPasswordWithUsername", &obj.PreventPasswordWithUsername, UnmarshalApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "passwordExpiration", &obj.PasswordExpiration, UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpiration)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "lockOutPolicy", &obj.LockOutPolicy, UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "minPasswordChangeInterval", &obj.MinPasswordChangeInterval, UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementLockOutPolicy : ApmSchemaAdvancedPasswordManagementLockOutPolicy struct
type ApmSchemaAdvancedPasswordManagementLockOutPolicy struct {
	Enabled *bool `json:"enabled" validate:"required"`

	Config *ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig `json:"config,omitempty"`
}

// NewApmSchemaAdvancedPasswordManagementLockOutPolicy : Instantiate ApmSchemaAdvancedPasswordManagementLockOutPolicy (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementLockOutPolicy(enabled bool) (model *ApmSchemaAdvancedPasswordManagementLockOutPolicy, err error) {
	model = &ApmSchemaAdvancedPasswordManagementLockOutPolicy{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicy unmarshals an instance of ApmSchemaAdvancedPasswordManagementLockOutPolicy from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementLockOutPolicy)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicyConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig : ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig struct
type ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig struct {
	LockOutTimeSec *int64 `json:"lockOutTimeSec" validate:"required"`

	NumOfAttempts *int64 `json:"numOfAttempts" validate:"required"`
}

// NewApmSchemaAdvancedPasswordManagementLockOutPolicyConfig : Instantiate ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementLockOutPolicyConfig(lockOutTimeSec int64, numOfAttempts int64) (model *ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig, err error) {
	model = &ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig{
		LockOutTimeSec: core.Int64Ptr(lockOutTimeSec),
		NumOfAttempts:  core.Int64Ptr(numOfAttempts),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicyConfig unmarshals an instance of ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementLockOutPolicyConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig)
	err = core.UnmarshalPrimitive(m, "lockOutTimeSec", &obj.LockOutTimeSec)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "numOfAttempts", &obj.NumOfAttempts)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval : ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval struct
type ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval struct {
	Enabled *bool `json:"enabled" validate:"required"`

	Config *ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig `json:"config,omitempty"`
}

// NewApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval : Instantiate ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval(enabled bool) (model *ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval, err error) {
	model = &ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval unmarshals an instance of ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig : ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig struct
type ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig struct {
	MinHoursToChangePassword *int64 `json:"minHoursToChangePassword" validate:"required"`
}

// NewApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig : Instantiate ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig(minHoursToChangePassword int64) (model *ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig, err error) {
	model = &ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig{
		MinHoursToChangePassword: core.Int64Ptr(minHoursToChangePassword),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig unmarshals an instance of ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig)
	err = core.UnmarshalPrimitive(m, "minHoursToChangePassword", &obj.MinHoursToChangePassword)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementPasswordExpiration : ApmSchemaAdvancedPasswordManagementPasswordExpiration struct
type ApmSchemaAdvancedPasswordManagementPasswordExpiration struct {
	Enabled *bool `json:"enabled" validate:"required"`

	Config *ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig `json:"config,omitempty"`
}

// NewApmSchemaAdvancedPasswordManagementPasswordExpiration : Instantiate ApmSchemaAdvancedPasswordManagementPasswordExpiration (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementPasswordExpiration(enabled bool) (model *ApmSchemaAdvancedPasswordManagementPasswordExpiration, err error) {
	model = &ApmSchemaAdvancedPasswordManagementPasswordExpiration{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpiration unmarshals an instance of ApmSchemaAdvancedPasswordManagementPasswordExpiration from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpiration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementPasswordExpiration)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpirationConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig : ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig struct
type ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig struct {
	DaysToExpire *int64 `json:"daysToExpire" validate:"required"`
}

// NewApmSchemaAdvancedPasswordManagementPasswordExpirationConfig : Instantiate ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementPasswordExpirationConfig(daysToExpire int64) (model *ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig, err error) {
	model = &ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig{
		DaysToExpire: core.Int64Ptr(daysToExpire),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpirationConfig unmarshals an instance of ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementPasswordExpirationConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig)
	err = core.UnmarshalPrimitive(m, "daysToExpire", &obj.DaysToExpire)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementPasswordReuse : ApmSchemaAdvancedPasswordManagementPasswordReuse struct
type ApmSchemaAdvancedPasswordManagementPasswordReuse struct {
	Enabled *bool `json:"enabled" validate:"required"`

	Config *ApmSchemaAdvancedPasswordManagementPasswordReuseConfig `json:"config,omitempty"`
}

// NewApmSchemaAdvancedPasswordManagementPasswordReuse : Instantiate ApmSchemaAdvancedPasswordManagementPasswordReuse (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementPasswordReuse(enabled bool) (model *ApmSchemaAdvancedPasswordManagementPasswordReuse, err error) {
	model = &ApmSchemaAdvancedPasswordManagementPasswordReuse{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuse unmarshals an instance of ApmSchemaAdvancedPasswordManagementPasswordReuse from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementPasswordReuse)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuseConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementPasswordReuseConfig : ApmSchemaAdvancedPasswordManagementPasswordReuseConfig struct
type ApmSchemaAdvancedPasswordManagementPasswordReuseConfig struct {
	MaxPasswordReuse *int64 `json:"maxPasswordReuse" validate:"required"`
}

// NewApmSchemaAdvancedPasswordManagementPasswordReuseConfig : Instantiate ApmSchemaAdvancedPasswordManagementPasswordReuseConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementPasswordReuseConfig(maxPasswordReuse int64) (model *ApmSchemaAdvancedPasswordManagementPasswordReuseConfig, err error) {
	model = &ApmSchemaAdvancedPasswordManagementPasswordReuseConfig{
		MaxPasswordReuse: core.Int64Ptr(maxPasswordReuse),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuseConfig unmarshals an instance of ApmSchemaAdvancedPasswordManagementPasswordReuseConfig from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementPasswordReuseConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementPasswordReuseConfig)
	err = core.UnmarshalPrimitive(m, "maxPasswordReuse", &obj.MaxPasswordReuse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername : ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername struct
type ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

// NewApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername : Instantiate ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername(enabled bool) (model *ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername, err error) {
	model = &ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername unmarshals an instance of ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername from the specified map of raw messages.
func UnmarshalApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Application : Application struct
type Application struct {
	// The application clientId.
	ClientID *string `json:"clientId,omitempty"`

	// The service tenantId.
	TenantID *string `json:"tenantId,omitempty"`

	Secret *string `json:"secret,omitempty"`

	// The application name.
	Name *string `json:"name,omitempty"`

	OAuthServerURL *string `json:"oAuthServerUrl,omitempty"`

	ProfilesURL *string `json:"profilesUrl,omitempty"`

	DiscoveryEndpoint *string `json:"discoveryEndpoint,omitempty"`

	// The type of application. Allowed types are regularwebapp and singlepageapp.
	Type *string `json:"type,omitempty"`
}

// UnmarshalApplication unmarshals an instance of Application from the specified map of raw messages.
func UnmarshalApplication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Application)
	err = core.UnmarshalPrimitive(m, "clientId", &obj.ClientID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tenantId", &obj.TenantID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "oAuthServerUrl", &obj.OAuthServerURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profilesUrl", &obj.ProfilesURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discoveryEndpoint", &obj.DiscoveryEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApplicationsList : ApplicationsList struct
type ApplicationsList struct {
	Applications []Application `json:"applications" validate:"required"`
}

// UnmarshalApplicationsList unmarshals an instance of ApplicationsList from the specified map of raw messages.
func UnmarshalApplicationsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicationsList)
	err = core.UnmarshalModel(m, "applications", &obj.Applications, UnmarshalApplication)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssignRoleToUserRolesItem : AssignRoleToUserRolesItem struct
type AssignRoleToUserRolesItem struct {
	ID *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`
}

// UnmarshalAssignRoleToUserRolesItem unmarshals an instance of AssignRoleToUserRolesItem from the specified map of raw messages.
func UnmarshalAssignRoleToUserRolesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssignRoleToUserRolesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChangePasswordOptions : The ChangePassword options.
type ChangePasswordOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The new password.
	NewPassword *string `validate:"required"`

	// The Cloud Directory unique user Id.
	UUID *string `validate:"required"`

	// The ip address the password changed from.
	ChangedIPAddress *string

	// Preferred language for resource. Format as described at RFC5646.
	Language *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewChangePasswordOptions : Instantiate ChangePasswordOptions
func (*AppIDManagementV4) NewChangePasswordOptions(tenantID string, newPassword string, uuid string) *ChangePasswordOptions {
	return &ChangePasswordOptions{
		TenantID:    core.StringPtr(tenantID),
		NewPassword: core.StringPtr(newPassword),
		UUID:        core.StringPtr(uuid),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ChangePasswordOptions) SetTenantID(tenantID string) *ChangePasswordOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetNewPassword : Allow user to set NewPassword
func (options *ChangePasswordOptions) SetNewPassword(newPassword string) *ChangePasswordOptions {
	options.NewPassword = core.StringPtr(newPassword)
	return options
}

// SetUUID : Allow user to set UUID
func (options *ChangePasswordOptions) SetUUID(uuid string) *ChangePasswordOptions {
	options.UUID = core.StringPtr(uuid)
	return options
}

// SetChangedIPAddress : Allow user to set ChangedIPAddress
func (options *ChangePasswordOptions) SetChangedIPAddress(changedIPAddress string) *ChangePasswordOptions {
	options.ChangedIPAddress = core.StringPtr(changedIPAddress)
	return options
}

// SetLanguage : Allow user to set Language
func (options *ChangePasswordOptions) SetLanguage(language string) *ChangePasswordOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangePasswordOptions) SetHeaders(param map[string]string) *ChangePasswordOptions {
	options.Headers = param
	return options
}

// CloudDirectoryConfigParamsInteractions : CloudDirectoryConfigParamsInteractions struct
type CloudDirectoryConfigParamsInteractions struct {
	IdentityConfirmation *CloudDirectoryConfigParamsInteractionsIdentityConfirmation `json:"identityConfirmation" validate:"required"`

	WelcomeEnabled *bool `json:"welcomeEnabled" validate:"required"`

	ResetPasswordEnabled *bool `json:"resetPasswordEnabled" validate:"required"`

	ResetPasswordNotificationEnable *bool `json:"resetPasswordNotificationEnable" validate:"required"`
}

// NewCloudDirectoryConfigParamsInteractions : Instantiate CloudDirectoryConfigParamsInteractions (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectoryConfigParamsInteractions(identityConfirmation *CloudDirectoryConfigParamsInteractionsIdentityConfirmation, welcomeEnabled bool, resetPasswordEnabled bool, resetPasswordNotificationEnable bool) (model *CloudDirectoryConfigParamsInteractions, err error) {
	model = &CloudDirectoryConfigParamsInteractions{
		IdentityConfirmation:            identityConfirmation,
		WelcomeEnabled:                  core.BoolPtr(welcomeEnabled),
		ResetPasswordEnabled:            core.BoolPtr(resetPasswordEnabled),
		ResetPasswordNotificationEnable: core.BoolPtr(resetPasswordNotificationEnable),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectoryConfigParamsInteractions unmarshals an instance of CloudDirectoryConfigParamsInteractions from the specified map of raw messages.
func UnmarshalCloudDirectoryConfigParamsInteractions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectoryConfigParamsInteractions)
	err = core.UnmarshalModel(m, "identityConfirmation", &obj.IdentityConfirmation, UnmarshalCloudDirectoryConfigParamsInteractionsIdentityConfirmation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "welcomeEnabled", &obj.WelcomeEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resetPasswordEnabled", &obj.ResetPasswordEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resetPasswordNotificationEnable", &obj.ResetPasswordNotificationEnable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectoryConfigParamsInteractionsIdentityConfirmation : CloudDirectoryConfigParamsInteractionsIdentityConfirmation struct
type CloudDirectoryConfigParamsInteractionsIdentityConfirmation struct {
	AccessMode *string `json:"accessMode" validate:"required"`

	Methods []string `json:"methods,omitempty"`
}

// Constants associated with the CloudDirectoryConfigParamsInteractionsIdentityConfirmation.AccessMode property.
const (
	CloudDirectoryConfigParamsInteractionsIdentityConfirmationAccessModeFalseConst       = "false"
	CloudDirectoryConfigParamsInteractionsIdentityConfirmationAccessModeFullConst        = "FULL"
	CloudDirectoryConfigParamsInteractionsIdentityConfirmationAccessModeRestrictiveConst = "RESTRICTIVE"
)

// Constants associated with the CloudDirectoryConfigParamsInteractionsIdentityConfirmation.Methods property.
const (
	CloudDirectoryConfigParamsInteractionsIdentityConfirmationMethodsEmailConst = "email"
)

// NewCloudDirectoryConfigParamsInteractionsIdentityConfirmation : Instantiate CloudDirectoryConfigParamsInteractionsIdentityConfirmation (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectoryConfigParamsInteractionsIdentityConfirmation(accessMode string) (model *CloudDirectoryConfigParamsInteractionsIdentityConfirmation, err error) {
	model = &CloudDirectoryConfigParamsInteractionsIdentityConfirmation{
		AccessMode: core.StringPtr(accessMode),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectoryConfigParamsInteractionsIdentityConfirmation unmarshals an instance of CloudDirectoryConfigParamsInteractionsIdentityConfirmation from the specified map of raw messages.
func UnmarshalCloudDirectoryConfigParamsInteractionsIdentityConfirmation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectoryConfigParamsInteractionsIdentityConfirmation)
	err = core.UnmarshalPrimitive(m, "accessMode", &obj.AccessMode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "methods", &obj.Methods)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectoryExportOptions : The CloudDirectoryExport options.
type CloudDirectoryExportOptions struct {
	// A custom string that will be use to encrypt and decrypt the users hashed password.
	EncryptionSecret *string `validate:"required"`

	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The first result in a set list of results.
	StartIndex *int64

	// The maximum number of results per page. Limit to 50 users per request.
	Count *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCloudDirectoryExportOptions : Instantiate CloudDirectoryExportOptions
func (*AppIDManagementV4) NewCloudDirectoryExportOptions(encryptionSecret string, tenantID string) *CloudDirectoryExportOptions {
	return &CloudDirectoryExportOptions{
		EncryptionSecret: core.StringPtr(encryptionSecret),
		TenantID:         core.StringPtr(tenantID),
	}
}

// SetEncryptionSecret : Allow user to set EncryptionSecret
func (options *CloudDirectoryExportOptions) SetEncryptionSecret(encryptionSecret string) *CloudDirectoryExportOptions {
	options.EncryptionSecret = core.StringPtr(encryptionSecret)
	return options
}

// SetTenantID : Allow user to set TenantID
func (options *CloudDirectoryExportOptions) SetTenantID(tenantID string) *CloudDirectoryExportOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetStartIndex : Allow user to set StartIndex
func (options *CloudDirectoryExportOptions) SetStartIndex(startIndex int64) *CloudDirectoryExportOptions {
	options.StartIndex = core.Int64Ptr(startIndex)
	return options
}

// SetCount : Allow user to set Count
func (options *CloudDirectoryExportOptions) SetCount(count int64) *CloudDirectoryExportOptions {
	options.Count = core.Int64Ptr(count)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CloudDirectoryExportOptions) SetHeaders(param map[string]string) *CloudDirectoryExportOptions {
	options.Headers = param
	return options
}

// CloudDirectoryGetUserinfoOptions : The CloudDirectoryGetUserinfo options.
type CloudDirectoryGetUserinfoOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCloudDirectoryGetUserinfoOptions : Instantiate CloudDirectoryGetUserinfoOptions
func (*AppIDManagementV4) NewCloudDirectoryGetUserinfoOptions(tenantID string, userID string) *CloudDirectoryGetUserinfoOptions {
	return &CloudDirectoryGetUserinfoOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *CloudDirectoryGetUserinfoOptions) SetTenantID(tenantID string) *CloudDirectoryGetUserinfoOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *CloudDirectoryGetUserinfoOptions) SetUserID(userID string) *CloudDirectoryGetUserinfoOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CloudDirectoryGetUserinfoOptions) SetHeaders(param map[string]string) *CloudDirectoryGetUserinfoOptions {
	options.Headers = param
	return options
}

// CloudDirectoryImportOptions : The CloudDirectoryImport options.
type CloudDirectoryImportOptions struct {
	// A custom string that will be use to encrypt and decrypt the users hashed password.
	EncryptionSecret *string `validate:"required"`

	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Users []ExportUserUsersItem `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCloudDirectoryImportOptions : Instantiate CloudDirectoryImportOptions
func (*AppIDManagementV4) NewCloudDirectoryImportOptions(encryptionSecret string, tenantID string, users []ExportUserUsersItem) *CloudDirectoryImportOptions {
	return &CloudDirectoryImportOptions{
		EncryptionSecret: core.StringPtr(encryptionSecret),
		TenantID:         core.StringPtr(tenantID),
		Users:            users,
	}
}

// SetEncryptionSecret : Allow user to set EncryptionSecret
func (options *CloudDirectoryImportOptions) SetEncryptionSecret(encryptionSecret string) *CloudDirectoryImportOptions {
	options.EncryptionSecret = core.StringPtr(encryptionSecret)
	return options
}

// SetTenantID : Allow user to set TenantID
func (options *CloudDirectoryImportOptions) SetTenantID(tenantID string) *CloudDirectoryImportOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUsers : Allow user to set Users
func (options *CloudDirectoryImportOptions) SetUsers(users []ExportUserUsersItem) *CloudDirectoryImportOptions {
	options.Users = users
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CloudDirectoryImportOptions) SetHeaders(param map[string]string) *CloudDirectoryImportOptions {
	options.Headers = param
	return options
}

// CloudDirectoryRemoveOptions : The CloudDirectoryRemove options.
type CloudDirectoryRemoveOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCloudDirectoryRemoveOptions : Instantiate CloudDirectoryRemoveOptions
func (*AppIDManagementV4) NewCloudDirectoryRemoveOptions(tenantID string, userID string) *CloudDirectoryRemoveOptions {
	return &CloudDirectoryRemoveOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *CloudDirectoryRemoveOptions) SetTenantID(tenantID string) *CloudDirectoryRemoveOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *CloudDirectoryRemoveOptions) SetUserID(userID string) *CloudDirectoryRemoveOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CloudDirectoryRemoveOptions) SetHeaders(param map[string]string) *CloudDirectoryRemoveOptions {
	options.Headers = param
	return options
}

// CloudDirectorySenderDetailsSenderDetails : CloudDirectorySenderDetailsSenderDetails struct
type CloudDirectorySenderDetailsSenderDetails struct {
	From *CloudDirectorySenderDetailsSenderDetailsFrom `json:"from" validate:"required"`

	ReplyTo *CloudDirectorySenderDetailsSenderDetailsReplyTo `json:"reply_to,omitempty"`

	LinkExpirationSec *int64 `json:"linkExpirationSec,omitempty"`
}

// NewCloudDirectorySenderDetailsSenderDetails : Instantiate CloudDirectorySenderDetailsSenderDetails (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectorySenderDetailsSenderDetails(from *CloudDirectorySenderDetailsSenderDetailsFrom) (model *CloudDirectorySenderDetailsSenderDetails, err error) {
	model = &CloudDirectorySenderDetailsSenderDetails{
		From: from,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectorySenderDetailsSenderDetails unmarshals an instance of CloudDirectorySenderDetailsSenderDetails from the specified map of raw messages.
func UnmarshalCloudDirectorySenderDetailsSenderDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectorySenderDetailsSenderDetails)
	err = core.UnmarshalModel(m, "from", &obj.From, UnmarshalCloudDirectorySenderDetailsSenderDetailsFrom)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reply_to", &obj.ReplyTo, UnmarshalCloudDirectorySenderDetailsSenderDetailsReplyTo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "linkExpirationSec", &obj.LinkExpirationSec)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectorySenderDetailsSenderDetailsFrom : CloudDirectorySenderDetailsSenderDetailsFrom struct
type CloudDirectorySenderDetailsSenderDetailsFrom struct {
	Name *string `json:"name,omitempty"`

	Email *string `json:"email" validate:"required"`
}

// NewCloudDirectorySenderDetailsSenderDetailsFrom : Instantiate CloudDirectorySenderDetailsSenderDetailsFrom (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectorySenderDetailsSenderDetailsFrom(email string) (model *CloudDirectorySenderDetailsSenderDetailsFrom, err error) {
	model = &CloudDirectorySenderDetailsSenderDetailsFrom{
		Email: core.StringPtr(email),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectorySenderDetailsSenderDetailsFrom unmarshals an instance of CloudDirectorySenderDetailsSenderDetailsFrom from the specified map of raw messages.
func UnmarshalCloudDirectorySenderDetailsSenderDetailsFrom(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectorySenderDetailsSenderDetailsFrom)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectorySenderDetailsSenderDetailsReplyTo : CloudDirectorySenderDetailsSenderDetailsReplyTo struct
type CloudDirectorySenderDetailsSenderDetailsReplyTo struct {
	Name *string `json:"name,omitempty"`

	Email *string `json:"email,omitempty"`
}

// UnmarshalCloudDirectorySenderDetailsSenderDetailsReplyTo unmarshals an instance of CloudDirectorySenderDetailsSenderDetailsReplyTo from the specified map of raw messages.
func UnmarshalCloudDirectorySenderDetailsSenderDetailsReplyTo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectorySenderDetailsSenderDetailsReplyTo)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCloudDirectoryUserOptions : The CreateCloudDirectoryUser options.
type CreateCloudDirectoryUserOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Emails []CreateNewUserEmailsItem `validate:"required"`

	Password *string `validate:"required"`

	Active *bool

	LockedUntil *int64

	DisplayName *string

	UserName *string

	// Accepted values "PENDING" or "CONFIRMED".
	Status *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateCloudDirectoryUserOptions : Instantiate CreateCloudDirectoryUserOptions
func (*AppIDManagementV4) NewCreateCloudDirectoryUserOptions(tenantID string, emails []CreateNewUserEmailsItem, password string) *CreateCloudDirectoryUserOptions {
	return &CreateCloudDirectoryUserOptions{
		TenantID: core.StringPtr(tenantID),
		Emails:   emails,
		Password: core.StringPtr(password),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *CreateCloudDirectoryUserOptions) SetTenantID(tenantID string) *CreateCloudDirectoryUserOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetEmails : Allow user to set Emails
func (options *CreateCloudDirectoryUserOptions) SetEmails(emails []CreateNewUserEmailsItem) *CreateCloudDirectoryUserOptions {
	options.Emails = emails
	return options
}

// SetPassword : Allow user to set Password
func (options *CreateCloudDirectoryUserOptions) SetPassword(password string) *CreateCloudDirectoryUserOptions {
	options.Password = core.StringPtr(password)
	return options
}

// SetActive : Allow user to set Active
func (options *CreateCloudDirectoryUserOptions) SetActive(active bool) *CreateCloudDirectoryUserOptions {
	options.Active = core.BoolPtr(active)
	return options
}

// SetLockedUntil : Allow user to set LockedUntil
func (options *CreateCloudDirectoryUserOptions) SetLockedUntil(lockedUntil int64) *CreateCloudDirectoryUserOptions {
	options.LockedUntil = core.Int64Ptr(lockedUntil)
	return options
}

// SetDisplayName : Allow user to set DisplayName
func (options *CreateCloudDirectoryUserOptions) SetDisplayName(displayName string) *CreateCloudDirectoryUserOptions {
	options.DisplayName = core.StringPtr(displayName)
	return options
}

// SetUserName : Allow user to set UserName
func (options *CreateCloudDirectoryUserOptions) SetUserName(userName string) *CreateCloudDirectoryUserOptions {
	options.UserName = core.StringPtr(userName)
	return options
}

// SetStatus : Allow user to set Status
func (options *CreateCloudDirectoryUserOptions) SetStatus(status string) *CreateCloudDirectoryUserOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCloudDirectoryUserOptions) SetHeaders(param map[string]string) *CreateCloudDirectoryUserOptions {
	options.Headers = param
	return options
}

// CreateNewUserEmailsItem : CreateNewUserEmailsItem struct
type CreateNewUserEmailsItem struct {
	Value *string `json:"value" validate:"required"`

	Primary *bool `json:"primary,omitempty"`
}

// NewCreateNewUserEmailsItem : Instantiate CreateNewUserEmailsItem (Generic Model Constructor)
func (*AppIDManagementV4) NewCreateNewUserEmailsItem(value string) (model *CreateNewUserEmailsItem, err error) {
	model = &CreateNewUserEmailsItem{
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCreateNewUserEmailsItem unmarshals an instance of CreateNewUserEmailsItem from the specified map of raw messages.
func UnmarshalCreateNewUserEmailsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateNewUserEmailsItem)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "primary", &obj.Primary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRoleOptions : The CreateRole options.
type CreateRoleOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Name *string `validate:"required"`

	Access []RoleAccessItem `validate:"required"`

	Description *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRoleOptions : Instantiate CreateRoleOptions
func (*AppIDManagementV4) NewCreateRoleOptions(tenantID string, name string, access []RoleAccessItem) *CreateRoleOptions {
	return &CreateRoleOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
		Access:   access,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *CreateRoleOptions) SetTenantID(tenantID string) *CreateRoleOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateRoleOptions) SetName(name string) *CreateRoleOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetAccess : Allow user to set Access
func (options *CreateRoleOptions) SetAccess(access []RoleAccessItem) *CreateRoleOptions {
	options.Access = access
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateRoleOptions) SetDescription(description string) *CreateRoleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRoleOptions) SetHeaders(param map[string]string) *CreateRoleOptions {
	options.Headers = param
	return options
}

// CustomIDPConfigParamsConfig : CustomIDPConfigParamsConfig struct
type CustomIDPConfigParamsConfig struct {
	PublicKey *string `json:"publicKey,omitempty"`
}

// UnmarshalCustomIDPConfigParamsConfig unmarshals an instance of CustomIDPConfigParamsConfig from the specified map of raw messages.
func UnmarshalCustomIDPConfigParamsConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomIDPConfigParamsConfig)
	err = core.UnmarshalPrimitive(m, "publicKey", &obj.PublicKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteActionURLOptions : The DeleteActionURL options.
type DeleteActionURLOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
	// your custom reset password page.
	Action *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteActionURLOptions.Action property.
// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
// your custom reset password page.
const (
	DeleteActionURLOptionsActionOnResetPasswordConst = "on_reset_password"
	DeleteActionURLOptionsActionOnUserVerifiedConst  = "on_user_verified"
)

// NewDeleteActionURLOptions : Instantiate DeleteActionURLOptions
func (*AppIDManagementV4) NewDeleteActionURLOptions(tenantID string, action string) *DeleteActionURLOptions {
	return &DeleteActionURLOptions{
		TenantID: core.StringPtr(tenantID),
		Action:   core.StringPtr(action),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *DeleteActionURLOptions) SetTenantID(tenantID string) *DeleteActionURLOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetAction : Allow user to set Action
func (options *DeleteActionURLOptions) SetAction(action string) *DeleteActionURLOptions {
	options.Action = core.StringPtr(action)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteActionURLOptions) SetHeaders(param map[string]string) *DeleteActionURLOptions {
	options.Headers = param
	return options
}

// DeleteApplicationOptions : The DeleteApplication options.
type DeleteApplicationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteApplicationOptions : Instantiate DeleteApplicationOptions
func (*AppIDManagementV4) NewDeleteApplicationOptions(tenantID string, clientID string) *DeleteApplicationOptions {
	return &DeleteApplicationOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *DeleteApplicationOptions) SetTenantID(tenantID string) *DeleteApplicationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *DeleteApplicationOptions) SetClientID(clientID string) *DeleteApplicationOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteApplicationOptions) SetHeaders(param map[string]string) *DeleteApplicationOptions {
	options.Headers = param
	return options
}

// DeleteCloudDirectoryUserOptions : The DeleteCloudDirectoryUser options.
type DeleteCloudDirectoryUserOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCloudDirectoryUserOptions : Instantiate DeleteCloudDirectoryUserOptions
func (*AppIDManagementV4) NewDeleteCloudDirectoryUserOptions(tenantID string, userID string) *DeleteCloudDirectoryUserOptions {
	return &DeleteCloudDirectoryUserOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *DeleteCloudDirectoryUserOptions) SetTenantID(tenantID string) *DeleteCloudDirectoryUserOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *DeleteCloudDirectoryUserOptions) SetUserID(userID string) *DeleteCloudDirectoryUserOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCloudDirectoryUserOptions) SetHeaders(param map[string]string) *DeleteCloudDirectoryUserOptions {
	options.Headers = param
	return options
}

// DeleteRoleOptions : The DeleteRole options.
type DeleteRoleOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The role identifier.
	RoleID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRoleOptions : Instantiate DeleteRoleOptions
func (*AppIDManagementV4) NewDeleteRoleOptions(tenantID string, roleID string) *DeleteRoleOptions {
	return &DeleteRoleOptions{
		TenantID: core.StringPtr(tenantID),
		RoleID:   core.StringPtr(roleID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *DeleteRoleOptions) SetTenantID(tenantID string) *DeleteRoleOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetRoleID : Allow user to set RoleID
func (options *DeleteRoleOptions) SetRoleID(roleID string) *DeleteRoleOptions {
	options.RoleID = core.StringPtr(roleID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRoleOptions) SetHeaders(param map[string]string) *DeleteRoleOptions {
	options.Headers = param
	return options
}

// DeleteTemplateOptions : The DeleteTemplate options.
type DeleteTemplateOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
	// "MFA_VERIFICATION".
	TemplateName *string `validate:"required,ne="`

	// Preferred language for resource. Format as described at RFC5646. According to the configured languages codes
	// returned from the `GET /management/v4/{tenantId}/config/ui/languages` API.
	Language *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the DeleteTemplateOptions.TemplateName property.
// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
// "MFA_VERIFICATION".
const (
	DeleteTemplateOptionsTemplateNameMFAVerificationConst  = "MFA_VERIFICATION"
	DeleteTemplateOptionsTemplateNamePasswordChangedConst  = "PASSWORD_CHANGED"
	DeleteTemplateOptionsTemplateNameResetPasswordConst    = "RESET_PASSWORD"
	DeleteTemplateOptionsTemplateNameUserVerificationConst = "USER_VERIFICATION"
	DeleteTemplateOptionsTemplateNameWelcomeConst          = "WELCOME"
)

// NewDeleteTemplateOptions : Instantiate DeleteTemplateOptions
func (*AppIDManagementV4) NewDeleteTemplateOptions(tenantID string, templateName string, language string) *DeleteTemplateOptions {
	return &DeleteTemplateOptions{
		TenantID:     core.StringPtr(tenantID),
		TemplateName: core.StringPtr(templateName),
		Language:     core.StringPtr(language),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *DeleteTemplateOptions) SetTenantID(tenantID string) *DeleteTemplateOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetTemplateName : Allow user to set TemplateName
func (options *DeleteTemplateOptions) SetTemplateName(templateName string) *DeleteTemplateOptions {
	options.TemplateName = core.StringPtr(templateName)
	return options
}

// SetLanguage : Allow user to set Language
func (options *DeleteTemplateOptions) SetLanguage(language string) *DeleteTemplateOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTemplateOptions) SetHeaders(param map[string]string) *DeleteTemplateOptions {
	options.Headers = param
	return options
}

// EmailDispatcherParamsCustom : EmailDispatcherParamsCustom struct
type EmailDispatcherParamsCustom struct {
	URL *string `json:"url" validate:"required"`

	Authorization *EmailDispatcherParamsCustomAuthorization `json:"authorization" validate:"required"`
}

// NewEmailDispatcherParamsCustom : Instantiate EmailDispatcherParamsCustom (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailDispatcherParamsCustom(url string, authorization *EmailDispatcherParamsCustomAuthorization) (model *EmailDispatcherParamsCustom, err error) {
	model = &EmailDispatcherParamsCustom{
		URL:           core.StringPtr(url),
		Authorization: authorization,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailDispatcherParamsCustom unmarshals an instance of EmailDispatcherParamsCustom from the specified map of raw messages.
func UnmarshalEmailDispatcherParamsCustom(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailDispatcherParamsCustom)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorization", &obj.Authorization, UnmarshalEmailDispatcherParamsCustomAuthorization)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailDispatcherParamsCustomAuthorization : EmailDispatcherParamsCustomAuthorization struct
type EmailDispatcherParamsCustomAuthorization struct {
	Type *string `json:"type" validate:"required"`

	Value *string `json:"value,omitempty"`

	Username *string `json:"username,omitempty"`

	Password *string `json:"password,omitempty"`
}

// Constants associated with the EmailDispatcherParamsCustomAuthorization.Type property.
const (
	EmailDispatcherParamsCustomAuthorizationTypeBasicConst = "basic"
	EmailDispatcherParamsCustomAuthorizationTypeNoneConst  = "none"
	EmailDispatcherParamsCustomAuthorizationTypeValueConst = "value"
)

// NewEmailDispatcherParamsCustomAuthorization : Instantiate EmailDispatcherParamsCustomAuthorization (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailDispatcherParamsCustomAuthorization(typeVar string) (model *EmailDispatcherParamsCustomAuthorization, err error) {
	model = &EmailDispatcherParamsCustomAuthorization{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailDispatcherParamsCustomAuthorization unmarshals an instance of EmailDispatcherParamsCustomAuthorization from the specified map of raw messages.
func UnmarshalEmailDispatcherParamsCustomAuthorization(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailDispatcherParamsCustomAuthorization)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailDispatcherParamsSendgrid : EmailDispatcherParamsSendgrid struct
type EmailDispatcherParamsSendgrid struct {
	APIKey *string `json:"apiKey" validate:"required"`
}

// NewEmailDispatcherParamsSendgrid : Instantiate EmailDispatcherParamsSendgrid (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailDispatcherParamsSendgrid(apiKey string) (model *EmailDispatcherParamsSendgrid, err error) {
	model = &EmailDispatcherParamsSendgrid{
		APIKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailDispatcherParamsSendgrid unmarshals an instance of EmailDispatcherParamsSendgrid from the specified map of raw messages.
func UnmarshalEmailDispatcherParamsSendgrid(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailDispatcherParamsSendgrid)
	err = core.UnmarshalPrimitive(m, "apiKey", &obj.APIKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingTestOptions : The EmailSettingTest options.
type EmailSettingTestOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	EmailTo *string `validate:"required"`

	EmailSettings *EmailSettingsTestParamsEmailSettings `validate:"required"`

	SenderDetails *EmailSettingsTestParamsSenderDetails `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewEmailSettingTestOptions : Instantiate EmailSettingTestOptions
func (*AppIDManagementV4) NewEmailSettingTestOptions(tenantID string, emailTo string, emailSettings *EmailSettingsTestParamsEmailSettings, senderDetails *EmailSettingsTestParamsSenderDetails) *EmailSettingTestOptions {
	return &EmailSettingTestOptions{
		TenantID:      core.StringPtr(tenantID),
		EmailTo:       core.StringPtr(emailTo),
		EmailSettings: emailSettings,
		SenderDetails: senderDetails,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *EmailSettingTestOptions) SetTenantID(tenantID string) *EmailSettingTestOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetEmailTo : Allow user to set EmailTo
func (options *EmailSettingTestOptions) SetEmailTo(emailTo string) *EmailSettingTestOptions {
	options.EmailTo = core.StringPtr(emailTo)
	return options
}

// SetEmailSettings : Allow user to set EmailSettings
func (options *EmailSettingTestOptions) SetEmailSettings(emailSettings *EmailSettingsTestParamsEmailSettings) *EmailSettingTestOptions {
	options.EmailSettings = emailSettings
	return options
}

// SetSenderDetails : Allow user to set SenderDetails
func (options *EmailSettingTestOptions) SetSenderDetails(senderDetails *EmailSettingsTestParamsSenderDetails) *EmailSettingTestOptions {
	options.SenderDetails = senderDetails
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EmailSettingTestOptions) SetHeaders(param map[string]string) *EmailSettingTestOptions {
	options.Headers = param
	return options
}

// EmailSettingsTestParamsEmailSettings : EmailSettingsTestParamsEmailSettings struct
type EmailSettingsTestParamsEmailSettings struct {
	Provider *string `json:"provider" validate:"required"`

	Sendgrid *EmailSettingsTestParamsEmailSettingsSendgrid `json:"sendgrid,omitempty"`

	Custom *EmailSettingsTestParamsEmailSettingsCustom `json:"custom,omitempty"`
}

// Constants associated with the EmailSettingsTestParamsEmailSettings.Provider property.
const (
	EmailSettingsTestParamsEmailSettingsProviderCustomConst   = "custom"
	EmailSettingsTestParamsEmailSettingsProviderSendgridConst = "sendgrid"
)

// NewEmailSettingsTestParamsEmailSettings : Instantiate EmailSettingsTestParamsEmailSettings (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsEmailSettings(provider string) (model *EmailSettingsTestParamsEmailSettings, err error) {
	model = &EmailSettingsTestParamsEmailSettings{
		Provider: core.StringPtr(provider),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsEmailSettings unmarshals an instance of EmailSettingsTestParamsEmailSettings from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsEmailSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsEmailSettings)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sendgrid", &obj.Sendgrid, UnmarshalEmailSettingsTestParamsEmailSettingsSendgrid)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "custom", &obj.Custom, UnmarshalEmailSettingsTestParamsEmailSettingsCustom)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsEmailSettingsCustom : EmailSettingsTestParamsEmailSettingsCustom struct
type EmailSettingsTestParamsEmailSettingsCustom struct {
	URL *string `json:"url" validate:"required"`

	Authorization *EmailSettingsTestParamsEmailSettingsCustomAuthorization `json:"authorization" validate:"required"`
}

// NewEmailSettingsTestParamsEmailSettingsCustom : Instantiate EmailSettingsTestParamsEmailSettingsCustom (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsEmailSettingsCustom(url string, authorization *EmailSettingsTestParamsEmailSettingsCustomAuthorization) (model *EmailSettingsTestParamsEmailSettingsCustom, err error) {
	model = &EmailSettingsTestParamsEmailSettingsCustom{
		URL:           core.StringPtr(url),
		Authorization: authorization,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsEmailSettingsCustom unmarshals an instance of EmailSettingsTestParamsEmailSettingsCustom from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsEmailSettingsCustom(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsEmailSettingsCustom)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorization", &obj.Authorization, UnmarshalEmailSettingsTestParamsEmailSettingsCustomAuthorization)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsEmailSettingsCustomAuthorization : EmailSettingsTestParamsEmailSettingsCustomAuthorization struct
type EmailSettingsTestParamsEmailSettingsCustomAuthorization struct {
	Type *string `json:"type" validate:"required"`

	Value *string `json:"value,omitempty"`

	Username *string `json:"username,omitempty"`

	Password *string `json:"password,omitempty"`
}

// Constants associated with the EmailSettingsTestParamsEmailSettingsCustomAuthorization.Type property.
const (
	EmailSettingsTestParamsEmailSettingsCustomAuthorizationTypeBasicConst = "basic"
	EmailSettingsTestParamsEmailSettingsCustomAuthorizationTypeNoneConst  = "none"
	EmailSettingsTestParamsEmailSettingsCustomAuthorizationTypeValueConst = "value"
)

// NewEmailSettingsTestParamsEmailSettingsCustomAuthorization : Instantiate EmailSettingsTestParamsEmailSettingsCustomAuthorization (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsEmailSettingsCustomAuthorization(typeVar string) (model *EmailSettingsTestParamsEmailSettingsCustomAuthorization, err error) {
	model = &EmailSettingsTestParamsEmailSettingsCustomAuthorization{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsEmailSettingsCustomAuthorization unmarshals an instance of EmailSettingsTestParamsEmailSettingsCustomAuthorization from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsEmailSettingsCustomAuthorization(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsEmailSettingsCustomAuthorization)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsEmailSettingsSendgrid : EmailSettingsTestParamsEmailSettingsSendgrid struct
type EmailSettingsTestParamsEmailSettingsSendgrid struct {
	APIKey *string `json:"apiKey" validate:"required"`
}

// NewEmailSettingsTestParamsEmailSettingsSendgrid : Instantiate EmailSettingsTestParamsEmailSettingsSendgrid (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsEmailSettingsSendgrid(apiKey string) (model *EmailSettingsTestParamsEmailSettingsSendgrid, err error) {
	model = &EmailSettingsTestParamsEmailSettingsSendgrid{
		APIKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsEmailSettingsSendgrid unmarshals an instance of EmailSettingsTestParamsEmailSettingsSendgrid from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsEmailSettingsSendgrid(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsEmailSettingsSendgrid)
	err = core.UnmarshalPrimitive(m, "apiKey", &obj.APIKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsSenderDetails : EmailSettingsTestParamsSenderDetails struct
type EmailSettingsTestParamsSenderDetails struct {
	From *EmailSettingsTestParamsSenderDetailsFrom `json:"from" validate:"required"`

	ReplyTo *EmailSettingsTestParamsSenderDetailsReplyTo `json:"reply_to,omitempty"`
}

// NewEmailSettingsTestParamsSenderDetails : Instantiate EmailSettingsTestParamsSenderDetails (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsSenderDetails(from *EmailSettingsTestParamsSenderDetailsFrom) (model *EmailSettingsTestParamsSenderDetails, err error) {
	model = &EmailSettingsTestParamsSenderDetails{
		From: from,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsSenderDetails unmarshals an instance of EmailSettingsTestParamsSenderDetails from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsSenderDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsSenderDetails)
	err = core.UnmarshalModel(m, "from", &obj.From, UnmarshalEmailSettingsTestParamsSenderDetailsFrom)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reply_to", &obj.ReplyTo, UnmarshalEmailSettingsTestParamsSenderDetailsReplyTo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsSenderDetailsFrom : EmailSettingsTestParamsSenderDetailsFrom struct
type EmailSettingsTestParamsSenderDetailsFrom struct {
	Email *string `json:"email" validate:"required"`

	Name *string `json:"name,omitempty"`
}

// NewEmailSettingsTestParamsSenderDetailsFrom : Instantiate EmailSettingsTestParamsSenderDetailsFrom (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsSenderDetailsFrom(email string) (model *EmailSettingsTestParamsSenderDetailsFrom, err error) {
	model = &EmailSettingsTestParamsSenderDetailsFrom{
		Email: core.StringPtr(email),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsSenderDetailsFrom unmarshals an instance of EmailSettingsTestParamsSenderDetailsFrom from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsSenderDetailsFrom(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsSenderDetailsFrom)
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailSettingsTestParamsSenderDetailsReplyTo : EmailSettingsTestParamsSenderDetailsReplyTo struct
type EmailSettingsTestParamsSenderDetailsReplyTo struct {
	Email *string `json:"email" validate:"required"`

	Name *string `json:"name,omitempty"`
}

// NewEmailSettingsTestParamsSenderDetailsReplyTo : Instantiate EmailSettingsTestParamsSenderDetailsReplyTo (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailSettingsTestParamsSenderDetailsReplyTo(email string) (model *EmailSettingsTestParamsSenderDetailsReplyTo, err error) {
	model = &EmailSettingsTestParamsSenderDetailsReplyTo{
		Email: core.StringPtr(email),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailSettingsTestParamsSenderDetailsReplyTo unmarshals an instance of EmailSettingsTestParamsSenderDetailsReplyTo from the specified map of raw messages.
func UnmarshalEmailSettingsTestParamsSenderDetailsReplyTo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailSettingsTestParamsSenderDetailsReplyTo)
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUserProfileUsersItem : ExportUserProfileUsersItem struct
type ExportUserProfileUsersItem struct {
	ID *string `json:"id" validate:"required"`

	Identities []ExportUserProfileUsersItemIdentitiesItem `json:"identities" validate:"required"`

	Attributes interface{} `json:"attributes" validate:"required"`

	Name *string `json:"name,omitempty"`

	Email *string `json:"email,omitempty"`

	Picture *string `json:"picture,omitempty"`

	Gender *string `json:"gender,omitempty"`

	Locale *string `json:"locale,omitempty"`

	PreferredUsername *string `json:"preferred_username,omitempty"`

	IDP *string `json:"idp,omitempty"`

	HashedIDPID *string `json:"hashedIdpId,omitempty"`

	HashedEmail *string `json:"hashedEmail,omitempty"`

	Roles []string `json:"roles,omitempty"`
}

// NewExportUserProfileUsersItem : Instantiate ExportUserProfileUsersItem (Generic Model Constructor)
func (*AppIDManagementV4) NewExportUserProfileUsersItem(id string, identities []ExportUserProfileUsersItemIdentitiesItem, attributes interface{}) (model *ExportUserProfileUsersItem, err error) {
	model = &ExportUserProfileUsersItem{
		ID:         core.StringPtr(id),
		Identities: identities,
		Attributes: attributes,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExportUserProfileUsersItem unmarshals an instance of ExportUserProfileUsersItem from the specified map of raw messages.
func UnmarshalExportUserProfileUsersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUserProfileUsersItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalExportUserProfileUsersItemIdentitiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "picture", &obj.Picture)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "gender", &obj.Gender)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locale", &obj.Locale)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "preferred_username", &obj.PreferredUsername)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "idp", &obj.IDP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hashedIdpId", &obj.HashedIDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hashedEmail", &obj.HashedEmail)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUserProfileUsersItemIdentitiesItem : ExportUserProfileUsersItemIdentitiesItem struct
type ExportUserProfileUsersItemIdentitiesItem struct {
	Provider *string `json:"provider,omitempty"`

	ID *string `json:"id,omitempty"`

	IDPUserInfo interface{} `json:"idpUserInfo,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of ExportUserProfileUsersItemIdentitiesItem
func (o *ExportUserProfileUsersItemIdentitiesItem) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of ExportUserProfileUsersItemIdentitiesItem
func (o *ExportUserProfileUsersItemIdentitiesItem) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of ExportUserProfileUsersItemIdentitiesItem
func (o *ExportUserProfileUsersItemIdentitiesItem) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of ExportUserProfileUsersItemIdentitiesItem
func (o *ExportUserProfileUsersItemIdentitiesItem) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Provider != nil {
		m["provider"] = o.Provider
	}
	if o.ID != nil {
		m["id"] = o.ID
	}
	if o.IDPUserInfo != nil {
		m["idpUserInfo"] = o.IDPUserInfo
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalExportUserProfileUsersItemIdentitiesItem unmarshals an instance of ExportUserProfileUsersItemIdentitiesItem from the specified map of raw messages.
func UnmarshalExportUserProfileUsersItemIdentitiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUserProfileUsersItemIdentitiesItem)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	delete(m, "provider")
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	delete(m, "id")
	err = core.UnmarshalPrimitive(m, "idpUserInfo", &obj.IDPUserInfo)
	if err != nil {
		return
	}
	delete(m, "idpUserInfo")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUserUsersItem : ExportUserUsersItem struct
type ExportUserUsersItem struct {
	ScimUser interface{} `json:"scimUser" validate:"required"`

	PasswordHash *string `json:"passwordHash" validate:"required"`

	PasswordHashAlg *string `json:"passwordHashAlg" validate:"required"`

	Profile *ExportUserUsersItemProfile `json:"profile" validate:"required"`

	Roles []string `json:"roles" validate:"required"`
}

// NewExportUserUsersItem : Instantiate ExportUserUsersItem (Generic Model Constructor)
func (*AppIDManagementV4) NewExportUserUsersItem(scimUser interface{}, passwordHash string, passwordHashAlg string, profile *ExportUserUsersItemProfile, roles []string) (model *ExportUserUsersItem, err error) {
	model = &ExportUserUsersItem{
		ScimUser:        scimUser,
		PasswordHash:    core.StringPtr(passwordHash),
		PasswordHashAlg: core.StringPtr(passwordHashAlg),
		Profile:         profile,
		Roles:           roles,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExportUserUsersItem unmarshals an instance of ExportUserUsersItem from the specified map of raw messages.
func UnmarshalExportUserUsersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUserUsersItem)
	err = core.UnmarshalPrimitive(m, "scimUser", &obj.ScimUser)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "passwordHash", &obj.PasswordHash)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "passwordHashAlg", &obj.PasswordHashAlg)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profile", &obj.Profile, UnmarshalExportUserUsersItemProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUserUsersItemProfile : ExportUserUsersItemProfile struct
type ExportUserUsersItemProfile struct {
	Attributes interface{} `json:"attributes" validate:"required"`
}

// NewExportUserUsersItemProfile : Instantiate ExportUserUsersItemProfile (Generic Model Constructor)
func (*AppIDManagementV4) NewExportUserUsersItemProfile(attributes interface{}) (model *ExportUserUsersItemProfile, err error) {
	model = &ExportUserUsersItemProfile{
		Attributes: attributes,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExportUserUsersItemProfile unmarshals an instance of ExportUserUsersItemProfile from the specified map of raw messages.
func UnmarshalExportUserUsersItemProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUserUsersItemProfile)
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookConfigParamsConfig : FacebookConfigParamsConfig struct
type FacebookConfigParamsConfig struct {
	IDPID *string `json:"idpId" validate:"required"`

	Secret *string `json:"secret" validate:"required"`
}

// UnmarshalFacebookConfigParamsConfig unmarshals an instance of FacebookConfigParamsConfig from the specified map of raw messages.
func UnmarshalFacebookConfigParamsConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookConfigParamsConfig)
	err = core.UnmarshalPrimitive(m, "idpId", &obj.IDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookConfigParamsPutConfig : FacebookConfigParamsPutConfig struct
type FacebookConfigParamsPutConfig struct {
	IDPID *string `json:"idpId" validate:"required"`

	Secret *string `json:"secret" validate:"required"`
}

// UnmarshalFacebookConfigParamsPutConfig unmarshals an instance of FacebookConfigParamsPutConfig from the specified map of raw messages.
func UnmarshalFacebookConfigParamsPutConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookConfigParamsPutConfig)
	err = core.UnmarshalPrimitive(m, "idpId", &obj.IDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookGoogleConfigParamsConfig : FacebookGoogleConfigParamsConfig struct
type FacebookGoogleConfigParamsConfig struct {
	IDPID *string `json:"idpId" validate:"required"`

	Secret *string `json:"secret" validate:"required"`
}

// NewFacebookGoogleConfigParamsConfig : Instantiate FacebookGoogleConfigParamsConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewFacebookGoogleConfigParamsConfig(idpID string, secret string) (model *FacebookGoogleConfigParamsConfig, err error) {
	model = &FacebookGoogleConfigParamsConfig{
		IDPID:  core.StringPtr(idpID),
		Secret: core.StringPtr(secret),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFacebookGoogleConfigParamsConfig unmarshals an instance of FacebookGoogleConfigParamsConfig from the specified map of raw messages.
func UnmarshalFacebookGoogleConfigParamsConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookGoogleConfigParamsConfig)
	err = core.UnmarshalPrimitive(m, "idpId", &obj.IDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForgotPasswordResultOptions : The ForgotPasswordResult options.
type ForgotPasswordResultOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The context that will be use to get the forgot password confirmation result.
	Context *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewForgotPasswordResultOptions : Instantiate ForgotPasswordResultOptions
func (*AppIDManagementV4) NewForgotPasswordResultOptions(tenantID string, context string) *ForgotPasswordResultOptions {
	return &ForgotPasswordResultOptions{
		TenantID: core.StringPtr(tenantID),
		Context:  core.StringPtr(context),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ForgotPasswordResultOptions) SetTenantID(tenantID string) *ForgotPasswordResultOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetContext : Allow user to set Context
func (options *ForgotPasswordResultOptions) SetContext(context string) *ForgotPasswordResultOptions {
	options.Context = core.StringPtr(context)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ForgotPasswordResultOptions) SetHeaders(param map[string]string) *ForgotPasswordResultOptions {
	options.Headers = param
	return options
}

// GetApplicationOptions : The GetApplication options.
type GetApplicationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetApplicationOptions : Instantiate GetApplicationOptions
func (*AppIDManagementV4) NewGetApplicationOptions(tenantID string, clientID string) *GetApplicationOptions {
	return &GetApplicationOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetApplicationOptions) SetTenantID(tenantID string) *GetApplicationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *GetApplicationOptions) SetClientID(clientID string) *GetApplicationOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetApplicationOptions) SetHeaders(param map[string]string) *GetApplicationOptions {
	options.Headers = param
	return options
}

// GetApplicationRolesOptions : The GetApplicationRoles options.
type GetApplicationRolesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetApplicationRolesOptions : Instantiate GetApplicationRolesOptions
func (*AppIDManagementV4) NewGetApplicationRolesOptions(tenantID string, clientID string) *GetApplicationRolesOptions {
	return &GetApplicationRolesOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetApplicationRolesOptions) SetTenantID(tenantID string) *GetApplicationRolesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *GetApplicationRolesOptions) SetClientID(clientID string) *GetApplicationRolesOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetApplicationRolesOptions) SetHeaders(param map[string]string) *GetApplicationRolesOptions {
	options.Headers = param
	return options
}

// GetApplicationScopesOptions : The GetApplicationScopes options.
type GetApplicationScopesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetApplicationScopesOptions : Instantiate GetApplicationScopesOptions
func (*AppIDManagementV4) NewGetApplicationScopesOptions(tenantID string, clientID string) *GetApplicationScopesOptions {
	return &GetApplicationScopesOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetApplicationScopesOptions) SetTenantID(tenantID string) *GetApplicationScopesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *GetApplicationScopesOptions) SetClientID(clientID string) *GetApplicationScopesOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetApplicationScopesOptions) SetHeaders(param map[string]string) *GetApplicationScopesOptions {
	options.Headers = param
	return options
}

// GetAuditStatusOptions : The GetAuditStatus options.
type GetAuditStatusOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAuditStatusOptions : Instantiate GetAuditStatusOptions
func (*AppIDManagementV4) NewGetAuditStatusOptions(tenantID string) *GetAuditStatusOptions {
	return &GetAuditStatusOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetAuditStatusOptions) SetTenantID(tenantID string) *GetAuditStatusOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAuditStatusOptions) SetHeaders(param map[string]string) *GetAuditStatusOptions {
	options.Headers = param
	return options
}

// GetAuditStatusResponse : GetAuditStatusResponse struct
type GetAuditStatusResponse struct {
	IsActive *bool `json:"isActive,omitempty"`
}

// UnmarshalGetAuditStatusResponse unmarshals an instance of GetAuditStatusResponse from the specified map of raw messages.
func UnmarshalGetAuditStatusResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetAuditStatusResponse)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetChannelOptions : The GetChannel options.
type GetChannelOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The MFA channel.
	Channel *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetChannelOptions.Channel property.
// The MFA channel.
const (
	GetChannelOptionsChannelEmailConst = "email"
	GetChannelOptionsChannelNexmoConst = "nexmo"
)

// NewGetChannelOptions : Instantiate GetChannelOptions
func (*AppIDManagementV4) NewGetChannelOptions(tenantID string, channel string) *GetChannelOptions {
	return &GetChannelOptions{
		TenantID: core.StringPtr(tenantID),
		Channel:  core.StringPtr(channel),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetChannelOptions) SetTenantID(tenantID string) *GetChannelOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetChannel : Allow user to set Channel
func (options *GetChannelOptions) SetChannel(channel string) *GetChannelOptions {
	options.Channel = core.StringPtr(channel)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetChannelOptions) SetHeaders(param map[string]string) *GetChannelOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryActionURLOptions : The GetCloudDirectoryActionURL options.
type GetCloudDirectoryActionURLOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
	// your custom reset password page.
	Action *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetCloudDirectoryActionURLOptions.Action property.
// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
// your custom reset password page.
const (
	GetCloudDirectoryActionURLOptionsActionOnResetPasswordConst = "on_reset_password"
	GetCloudDirectoryActionURLOptionsActionOnUserVerifiedConst  = "on_user_verified"
)

// NewGetCloudDirectoryActionURLOptions : Instantiate GetCloudDirectoryActionURLOptions
func (*AppIDManagementV4) NewGetCloudDirectoryActionURLOptions(tenantID string, action string) *GetCloudDirectoryActionURLOptions {
	return &GetCloudDirectoryActionURLOptions{
		TenantID: core.StringPtr(tenantID),
		Action:   core.StringPtr(action),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryActionURLOptions) SetTenantID(tenantID string) *GetCloudDirectoryActionURLOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetAction : Allow user to set Action
func (options *GetCloudDirectoryActionURLOptions) SetAction(action string) *GetCloudDirectoryActionURLOptions {
	options.Action = core.StringPtr(action)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryActionURLOptions) SetHeaders(param map[string]string) *GetCloudDirectoryActionURLOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryAdvancedPasswordManagementOptions : The GetCloudDirectoryAdvancedPasswordManagement options.
type GetCloudDirectoryAdvancedPasswordManagementOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectoryAdvancedPasswordManagementOptions : Instantiate GetCloudDirectoryAdvancedPasswordManagementOptions
func (*AppIDManagementV4) NewGetCloudDirectoryAdvancedPasswordManagementOptions(tenantID string) *GetCloudDirectoryAdvancedPasswordManagementOptions {
	return &GetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryAdvancedPasswordManagementOptions) SetTenantID(tenantID string) *GetCloudDirectoryAdvancedPasswordManagementOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryAdvancedPasswordManagementOptions) SetHeaders(param map[string]string) *GetCloudDirectoryAdvancedPasswordManagementOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryEmailDispatcherOptions : The GetCloudDirectoryEmailDispatcher options.
type GetCloudDirectoryEmailDispatcherOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectoryEmailDispatcherOptions : Instantiate GetCloudDirectoryEmailDispatcherOptions
func (*AppIDManagementV4) NewGetCloudDirectoryEmailDispatcherOptions(tenantID string) *GetCloudDirectoryEmailDispatcherOptions {
	return &GetCloudDirectoryEmailDispatcherOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryEmailDispatcherOptions) SetTenantID(tenantID string) *GetCloudDirectoryEmailDispatcherOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryEmailDispatcherOptions) SetHeaders(param map[string]string) *GetCloudDirectoryEmailDispatcherOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryIDPOptions : The GetCloudDirectoryIDP options.
type GetCloudDirectoryIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectoryIDPOptions : Instantiate GetCloudDirectoryIDPOptions
func (*AppIDManagementV4) NewGetCloudDirectoryIDPOptions(tenantID string) *GetCloudDirectoryIDPOptions {
	return &GetCloudDirectoryIDPOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryIDPOptions) SetTenantID(tenantID string) *GetCloudDirectoryIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryIDPOptions) SetHeaders(param map[string]string) *GetCloudDirectoryIDPOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryPasswordRegexOptions : The GetCloudDirectoryPasswordRegex options.
type GetCloudDirectoryPasswordRegexOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectoryPasswordRegexOptions : Instantiate GetCloudDirectoryPasswordRegexOptions
func (*AppIDManagementV4) NewGetCloudDirectoryPasswordRegexOptions(tenantID string) *GetCloudDirectoryPasswordRegexOptions {
	return &GetCloudDirectoryPasswordRegexOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryPasswordRegexOptions) SetTenantID(tenantID string) *GetCloudDirectoryPasswordRegexOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryPasswordRegexOptions) SetHeaders(param map[string]string) *GetCloudDirectoryPasswordRegexOptions {
	options.Headers = param
	return options
}

// GetCloudDirectorySenderDetailsOptions : The GetCloudDirectorySenderDetails options.
type GetCloudDirectorySenderDetailsOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectorySenderDetailsOptions : Instantiate GetCloudDirectorySenderDetailsOptions
func (*AppIDManagementV4) NewGetCloudDirectorySenderDetailsOptions(tenantID string) *GetCloudDirectorySenderDetailsOptions {
	return &GetCloudDirectorySenderDetailsOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectorySenderDetailsOptions) SetTenantID(tenantID string) *GetCloudDirectorySenderDetailsOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectorySenderDetailsOptions) SetHeaders(param map[string]string) *GetCloudDirectorySenderDetailsOptions {
	options.Headers = param
	return options
}

// GetCloudDirectoryUserOptions : The GetCloudDirectoryUser options.
type GetCloudDirectoryUserOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCloudDirectoryUserOptions : Instantiate GetCloudDirectoryUserOptions
func (*AppIDManagementV4) NewGetCloudDirectoryUserOptions(tenantID string, userID string) *GetCloudDirectoryUserOptions {
	return &GetCloudDirectoryUserOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCloudDirectoryUserOptions) SetTenantID(tenantID string) *GetCloudDirectoryUserOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *GetCloudDirectoryUserOptions) SetUserID(userID string) *GetCloudDirectoryUserOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCloudDirectoryUserOptions) SetHeaders(param map[string]string) *GetCloudDirectoryUserOptions {
	options.Headers = param
	return options
}

// GetCustomIDPOptions : The GetCustomIDP options.
type GetCustomIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCustomIDPOptions : Instantiate GetCustomIDPOptions
func (*AppIDManagementV4) NewGetCustomIDPOptions(tenantID string) *GetCustomIDPOptions {
	return &GetCustomIDPOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetCustomIDPOptions) SetTenantID(tenantID string) *GetCustomIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCustomIDPOptions) SetHeaders(param map[string]string) *GetCustomIDPOptions {
	options.Headers = param
	return options
}

// GetExtensionConfigOptions : The GetExtensionConfig options.
type GetExtensionConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The name of the extension.
	Name *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetExtensionConfigOptions.Name property.
// The name of the extension.
const (
	GetExtensionConfigOptionsNamePostmfaConst = "postmfa"
	GetExtensionConfigOptionsNamePremfaConst  = "premfa"
)

// NewGetExtensionConfigOptions : Instantiate GetExtensionConfigOptions
func (*AppIDManagementV4) NewGetExtensionConfigOptions(tenantID string, name string) *GetExtensionConfigOptions {
	return &GetExtensionConfigOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetExtensionConfigOptions) SetTenantID(tenantID string) *GetExtensionConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *GetExtensionConfigOptions) SetName(name string) *GetExtensionConfigOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetExtensionConfigOptions) SetHeaders(param map[string]string) *GetExtensionConfigOptions {
	options.Headers = param
	return options
}

// GetFacebookIDPOptions : The GetFacebookIDP options.
type GetFacebookIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetFacebookIDPOptions : Instantiate GetFacebookIDPOptions
func (*AppIDManagementV4) NewGetFacebookIDPOptions(tenantID string) *GetFacebookIDPOptions {
	return &GetFacebookIDPOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetFacebookIDPOptions) SetTenantID(tenantID string) *GetFacebookIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetFacebookIDPOptions) SetHeaders(param map[string]string) *GetFacebookIDPOptions {
	options.Headers = param
	return options
}

// GetGoogleIDPOptions : The GetGoogleIDP options.
type GetGoogleIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetGoogleIDPOptions : Instantiate GetGoogleIDPOptions
func (*AppIDManagementV4) NewGetGoogleIDPOptions(tenantID string) *GetGoogleIDPOptions {
	return &GetGoogleIDPOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetGoogleIDPOptions) SetTenantID(tenantID string) *GetGoogleIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetGoogleIDPOptions) SetHeaders(param map[string]string) *GetGoogleIDPOptions {
	options.Headers = param
	return options
}

// GetLocalizationOptions : The GetLocalization options.
type GetLocalizationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLocalizationOptions : Instantiate GetLocalizationOptions
func (*AppIDManagementV4) NewGetLocalizationOptions(tenantID string) *GetLocalizationOptions {
	return &GetLocalizationOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetLocalizationOptions) SetTenantID(tenantID string) *GetLocalizationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLocalizationOptions) SetHeaders(param map[string]string) *GetLocalizationOptions {
	options.Headers = param
	return options
}

// GetMFAConfigOptions : The GetMFAConfig options.
type GetMFAConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetMFAConfigOptions : Instantiate GetMFAConfigOptions
func (*AppIDManagementV4) NewGetMFAConfigOptions(tenantID string) *GetMFAConfigOptions {
	return &GetMFAConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetMFAConfigOptions) SetTenantID(tenantID string) *GetMFAConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetMFAConfigOptions) SetHeaders(param map[string]string) *GetMFAConfigOptions {
	options.Headers = param
	return options
}

// GetMediaOptions : The GetMedia options.
type GetMediaOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetMediaOptions : Instantiate GetMediaOptions
func (*AppIDManagementV4) NewGetMediaOptions(tenantID string) *GetMediaOptions {
	return &GetMediaOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetMediaOptions) SetTenantID(tenantID string) *GetMediaOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetMediaOptions) SetHeaders(param map[string]string) *GetMediaOptions {
	options.Headers = param
	return options
}

// GetMediaResponse : GetMediaResponse struct
type GetMediaResponse struct {
	Image *string `json:"image" validate:"required"`
}

// UnmarshalGetMediaResponse unmarshals an instance of GetMediaResponse from the specified map of raw messages.
func UnmarshalGetMediaResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetMediaResponse)
	err = core.UnmarshalPrimitive(m, "image", &obj.Image)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetRateLimitConfigOptions : The GetRateLimitConfig options.
type GetRateLimitConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRateLimitConfigOptions : Instantiate GetRateLimitConfigOptions
func (*AppIDManagementV4) NewGetRateLimitConfigOptions(tenantID string) *GetRateLimitConfigOptions {
	return &GetRateLimitConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetRateLimitConfigOptions) SetTenantID(tenantID string) *GetRateLimitConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRateLimitConfigOptions) SetHeaders(param map[string]string) *GetRateLimitConfigOptions {
	options.Headers = param
	return options
}

// GetRedirectUrisOptions : The GetRedirectUris options.
type GetRedirectUrisOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRedirectUrisOptions : Instantiate GetRedirectUrisOptions
func (*AppIDManagementV4) NewGetRedirectUrisOptions(tenantID string) *GetRedirectUrisOptions {
	return &GetRedirectUrisOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetRedirectUrisOptions) SetTenantID(tenantID string) *GetRedirectUrisOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRedirectUrisOptions) SetHeaders(param map[string]string) *GetRedirectUrisOptions {
	options.Headers = param
	return options
}

// GetRoleOptions : The GetRole options.
type GetRoleOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The role identifier.
	RoleID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRoleOptions : Instantiate GetRoleOptions
func (*AppIDManagementV4) NewGetRoleOptions(tenantID string, roleID string) *GetRoleOptions {
	return &GetRoleOptions{
		TenantID: core.StringPtr(tenantID),
		RoleID:   core.StringPtr(roleID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetRoleOptions) SetTenantID(tenantID string) *GetRoleOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetRoleID : Allow user to set RoleID
func (options *GetRoleOptions) SetRoleID(roleID string) *GetRoleOptions {
	options.RoleID = core.StringPtr(roleID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRoleOptions) SetHeaders(param map[string]string) *GetRoleOptions {
	options.Headers = param
	return options
}

// GetSMSChannelConfig : GetSMSChannelConfig struct
type GetSMSChannelConfig struct {
	Key *string `json:"key,omitempty"`

	Secret *string `json:"secret,omitempty"`

	From *string `json:"from,omitempty"`

	Provider *string `json:"provider,omitempty"`
}

// UnmarshalGetSMSChannelConfig unmarshals an instance of GetSMSChannelConfig from the specified map of raw messages.
func UnmarshalGetSMSChannelConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSMSChannelConfig)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "from", &obj.From)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSSOConfigOptions : The GetSSOConfig options.
type GetSSOConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSSOConfigOptions : Instantiate GetSSOConfigOptions
func (*AppIDManagementV4) NewGetSSOConfigOptions(tenantID string) *GetSSOConfigOptions {
	return &GetSSOConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetSSOConfigOptions) SetTenantID(tenantID string) *GetSSOConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSSOConfigOptions) SetHeaders(param map[string]string) *GetSSOConfigOptions {
	options.Headers = param
	return options
}

// GetSAMLIDPOptions : The GetSAMLIDP options.
type GetSAMLIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSAMLIDPOptions : Instantiate GetSAMLIDPOptions
func (*AppIDManagementV4) NewGetSAMLIDPOptions(tenantID string) *GetSAMLIDPOptions {
	return &GetSAMLIDPOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetSAMLIDPOptions) SetTenantID(tenantID string) *GetSAMLIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSAMLIDPOptions) SetHeaders(param map[string]string) *GetSAMLIDPOptions {
	options.Headers = param
	return options
}

// GetSAMLMetadataOptions : The GetSAMLMetadata options.
type GetSAMLMetadataOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSAMLMetadataOptions : Instantiate GetSAMLMetadataOptions
func (*AppIDManagementV4) NewGetSAMLMetadataOptions(tenantID string) *GetSAMLMetadataOptions {
	return &GetSAMLMetadataOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetSAMLMetadataOptions) SetTenantID(tenantID string) *GetSAMLMetadataOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSAMLMetadataOptions) SetHeaders(param map[string]string) *GetSAMLMetadataOptions {
	options.Headers = param
	return options
}

// GetTemplateOptions : The GetTemplate options.
type GetTemplateOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
	// "MFA_VERIFICATION".
	TemplateName *string `validate:"required,ne="`

	// Preferred language for resource. Format as described at RFC5646. According to the configured languages codes
	// returned from the `GET /management/v4/{tenantId}/config/ui/languages` API.
	Language *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetTemplateOptions.TemplateName property.
// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
// "MFA_VERIFICATION".
const (
	GetTemplateOptionsTemplateNameMFAVerificationConst  = "MFA_VERIFICATION"
	GetTemplateOptionsTemplateNamePasswordChangedConst  = "PASSWORD_CHANGED"
	GetTemplateOptionsTemplateNameResetPasswordConst    = "RESET_PASSWORD"
	GetTemplateOptionsTemplateNameUserVerificationConst = "USER_VERIFICATION"
	GetTemplateOptionsTemplateNameWelcomeConst          = "WELCOME"
)

// NewGetTemplateOptions : Instantiate GetTemplateOptions
func (*AppIDManagementV4) NewGetTemplateOptions(tenantID string, templateName string, language string) *GetTemplateOptions {
	return &GetTemplateOptions{
		TenantID:     core.StringPtr(tenantID),
		TemplateName: core.StringPtr(templateName),
		Language:     core.StringPtr(language),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetTemplateOptions) SetTenantID(tenantID string) *GetTemplateOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetTemplateName : Allow user to set TemplateName
func (options *GetTemplateOptions) SetTemplateName(templateName string) *GetTemplateOptions {
	options.TemplateName = core.StringPtr(templateName)
	return options
}

// SetLanguage : Allow user to set Language
func (options *GetTemplateOptions) SetLanguage(language string) *GetTemplateOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateOptions) SetHeaders(param map[string]string) *GetTemplateOptions {
	options.Headers = param
	return options
}

// GetThemeColorOptions : The GetThemeColor options.
type GetThemeColorOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetThemeColorOptions : Instantiate GetThemeColorOptions
func (*AppIDManagementV4) NewGetThemeColorOptions(tenantID string) *GetThemeColorOptions {
	return &GetThemeColorOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetThemeColorOptions) SetTenantID(tenantID string) *GetThemeColorOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetThemeColorOptions) SetHeaders(param map[string]string) *GetThemeColorOptions {
	options.Headers = param
	return options
}

// GetThemeColorResponse : GetThemeColorResponse struct
type GetThemeColorResponse struct {
	HeaderColor *string `json:"headerColor" validate:"required"`
}

// UnmarshalGetThemeColorResponse unmarshals an instance of GetThemeColorResponse from the specified map of raw messages.
func UnmarshalGetThemeColorResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetThemeColorResponse)
	err = core.UnmarshalPrimitive(m, "headerColor", &obj.HeaderColor)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetThemeTextOptions : The GetThemeText options.
type GetThemeTextOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetThemeTextOptions : Instantiate GetThemeTextOptions
func (*AppIDManagementV4) NewGetThemeTextOptions(tenantID string) *GetThemeTextOptions {
	return &GetThemeTextOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetThemeTextOptions) SetTenantID(tenantID string) *GetThemeTextOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetThemeTextOptions) SetHeaders(param map[string]string) *GetThemeTextOptions {
	options.Headers = param
	return options
}

// GetThemeTextResponse : GetThemeTextResponse struct
type GetThemeTextResponse struct {
	Footnote *string `json:"footnote" validate:"required"`

	TabTitle *string `json:"tabTitle" validate:"required"`
}

// UnmarshalGetThemeTextResponse unmarshals an instance of GetThemeTextResponse from the specified map of raw messages.
func UnmarshalGetThemeTextResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetThemeTextResponse)
	err = core.UnmarshalPrimitive(m, "footnote", &obj.Footnote)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tabTitle", &obj.TabTitle)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetTokensConfigOptions : The GetTokensConfig options.
type GetTokensConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTokensConfigOptions : Instantiate GetTokensConfigOptions
func (*AppIDManagementV4) NewGetTokensConfigOptions(tenantID string) *GetTokensConfigOptions {
	return &GetTokensConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetTokensConfigOptions) SetTenantID(tenantID string) *GetTokensConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetTokensConfigOptions) SetHeaders(param map[string]string) *GetTokensConfigOptions {
	options.Headers = param
	return options
}

// GetUserAndProfileIdentitiesItem : GetUserAndProfileIdentitiesItem struct
type GetUserAndProfileIdentitiesItem struct {
	Provider *string `json:"provider,omitempty"`

	ID *string `json:"id,omitempty"`

	IDPUserInfo *GetUserAndProfileIdentitiesItemIDPUserInfo `json:"idpUserInfo,omitempty"`
}

// UnmarshalGetUserAndProfileIdentitiesItem unmarshals an instance of GetUserAndProfileIdentitiesItem from the specified map of raw messages.
func UnmarshalGetUserAndProfileIdentitiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserAndProfileIdentitiesItem)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "idpUserInfo", &obj.IDPUserInfo, UnmarshalGetUserAndProfileIdentitiesItemIDPUserInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserAndProfileIdentitiesItemIDPUserInfo : GetUserAndProfileIdentitiesItemIDPUserInfo struct
type GetUserAndProfileIdentitiesItemIDPUserInfo struct {
	DisplayName *string `json:"displayName,omitempty"`

	Active *bool `json:"active" validate:"required"`

	LockedUntil *int64 `json:"lockedUntil,omitempty"`

	Emails []GetUserEmailsItem `json:"emails" validate:"required"`

	Meta *GetUserMeta `json:"meta" validate:"required"`

	Schemas []string `json:"schemas,omitempty"`

	Name *GetUserName `json:"name,omitempty"`

	UserName *string `json:"userName,omitempty"`

	ID *string `json:"id" validate:"required"`

	Status *string `json:"status" validate:"required"`
}

// UnmarshalGetUserAndProfileIdentitiesItemIDPUserInfo unmarshals an instance of GetUserAndProfileIdentitiesItemIDPUserInfo from the specified map of raw messages.
func UnmarshalGetUserAndProfileIdentitiesItemIDPUserInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserAndProfileIdentitiesItemIDPUserInfo)
	err = core.UnmarshalPrimitive(m, "displayName", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lockedUntil", &obj.LockedUntil)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "emails", &obj.Emails, UnmarshalGetUserEmailsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "meta", &obj.Meta, UnmarshalGetUserMeta)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schemas", &obj.Schemas)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "name", &obj.Name, UnmarshalGetUserName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "userName", &obj.UserName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserEmailsItem : GetUserEmailsItem struct
type GetUserEmailsItem struct {
	Value *string `json:"value" validate:"required"`

	Primary *bool `json:"primary,omitempty"`
}

// UnmarshalGetUserEmailsItem unmarshals an instance of GetUserEmailsItem from the specified map of raw messages.
func UnmarshalGetUserEmailsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserEmailsItem)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "primary", &obj.Primary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserMeta : GetUserMeta struct
type GetUserMeta struct {
	Created *strfmt.DateTime `json:"created,omitempty"`

	LastModified *strfmt.DateTime `json:"lastModified,omitempty"`

	ResourceType *string `json:"resourceType,omitempty"`
}

// UnmarshalGetUserMeta unmarshals an instance of GetUserMeta from the specified map of raw messages.
func UnmarshalGetUserMeta(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserMeta)
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lastModified", &obj.LastModified)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resourceType", &obj.ResourceType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserName : GetUserName struct
type GetUserName struct {
	GivenName *string `json:"givenName,omitempty"`

	FamilyName *string `json:"familyName,omitempty"`

	Formatted *string `json:"formatted,omitempty"`
}

// UnmarshalGetUserName unmarshals an instance of GetUserName from the specified map of raw messages.
func UnmarshalGetUserName(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserName)
	err = core.UnmarshalPrimitive(m, "givenName", &obj.GivenName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "familyName", &obj.FamilyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "formatted", &obj.Formatted)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserProfilesConfigOptions : The GetUserProfilesConfig options.
type GetUserProfilesConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUserProfilesConfigOptions : Instantiate GetUserProfilesConfigOptions
func (*AppIDManagementV4) NewGetUserProfilesConfigOptions(tenantID string) *GetUserProfilesConfigOptions {
	return &GetUserProfilesConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetUserProfilesConfigOptions) SetTenantID(tenantID string) *GetUserProfilesConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetUserProfilesConfigOptions) SetHeaders(param map[string]string) *GetUserProfilesConfigOptions {
	options.Headers = param
	return options
}

// GetUserProfilesConfigResponse : GetUserProfilesConfigResponse struct
type GetUserProfilesConfigResponse struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

// UnmarshalGetUserProfilesConfigResponse unmarshals an instance of GetUserProfilesConfigResponse from the specified map of raw messages.
func UnmarshalGetUserProfilesConfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserProfilesConfigResponse)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserRolesOptions : The GetUserRoles options.
type GetUserRolesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUserRolesOptions : Instantiate GetUserRolesOptions
func (*AppIDManagementV4) NewGetUserRolesOptions(tenantID string, id string) *GetUserRolesOptions {
	return &GetUserRolesOptions{
		TenantID: core.StringPtr(tenantID),
		ID:       core.StringPtr(id),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *GetUserRolesOptions) SetTenantID(tenantID string) *GetUserRolesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *GetUserRolesOptions) SetID(id string) *GetUserRolesOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetUserRolesOptions) SetHeaders(param map[string]string) *GetUserRolesOptions {
	options.Headers = param
	return options
}

// GetUserRolesResponseRolesItem : GetUserRolesResponseRolesItem struct
type GetUserRolesResponseRolesItem struct {
	ID *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`
}

// UnmarshalGetUserRolesResponseRolesItem unmarshals an instance of GetUserRolesResponseRolesItem from the specified map of raw messages.
func UnmarshalGetUserRolesResponseRolesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserRolesResponseRolesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoogleConfigParamsConfig : GoogleConfigParamsConfig struct
type GoogleConfigParamsConfig struct {
	IDPID *string `json:"idpId" validate:"required"`

	Secret *string `json:"secret" validate:"required"`
}

// UnmarshalGoogleConfigParamsConfig unmarshals an instance of GoogleConfigParamsConfig from the specified map of raw messages.
func UnmarshalGoogleConfigParamsConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoogleConfigParamsConfig)
	err = core.UnmarshalPrimitive(m, "idpId", &obj.IDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoogleConfigParamsPutConfig : GoogleConfigParamsPutConfig struct
type GoogleConfigParamsPutConfig struct {
	IDPID *string `json:"idpId" validate:"required"`

	Secret *string `json:"secret" validate:"required"`
}

// UnmarshalGoogleConfigParamsPutConfig unmarshals an instance of GoogleConfigParamsPutConfig from the specified map of raw messages.
func UnmarshalGoogleConfigParamsPutConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoogleConfigParamsPutConfig)
	err = core.UnmarshalPrimitive(m, "idpId", &obj.IDPID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportProfilesResponseFailReasonsItem : ImportProfilesResponseFailReasonsItem struct
type ImportProfilesResponseFailReasonsItem struct {
	OriginalID *string `json:"originalId,omitempty"`

	IDP *string `json:"idp,omitempty"`

	Error interface{} `json:"error,omitempty"`
}

// UnmarshalImportProfilesResponseFailReasonsItem unmarshals an instance of ImportProfilesResponseFailReasonsItem from the specified map of raw messages.
func UnmarshalImportProfilesResponseFailReasonsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportProfilesResponseFailReasonsItem)
	err = core.UnmarshalPrimitive(m, "originalId", &obj.OriginalID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "idp", &obj.IDP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportResponseFailReasonsItem : ImportResponseFailReasonsItem struct
type ImportResponseFailReasonsItem struct {
	OriginalID *string `json:"originalId,omitempty"`

	ID *string `json:"id,omitempty"`

	Email *string `json:"email,omitempty"`

	UserName *string `json:"userName,omitempty"`

	Error interface{} `json:"error,omitempty"`
}

// UnmarshalImportResponseFailReasonsItem unmarshals an instance of ImportResponseFailReasonsItem from the specified map of raw messages.
func UnmarshalImportResponseFailReasonsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportResponseFailReasonsItem)
	err = core.UnmarshalPrimitive(m, "originalId", &obj.OriginalID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "userName", &obj.UserName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InvalidateUserSSOSessionsOptions : The InvalidateUserSSOSessions options.
type InvalidateUserSSOSessionsOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewInvalidateUserSSOSessionsOptions : Instantiate InvalidateUserSSOSessionsOptions
func (*AppIDManagementV4) NewInvalidateUserSSOSessionsOptions(tenantID string, userID string) *InvalidateUserSSOSessionsOptions {
	return &InvalidateUserSSOSessionsOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *InvalidateUserSSOSessionsOptions) SetTenantID(tenantID string) *InvalidateUserSSOSessionsOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *InvalidateUserSSOSessionsOptions) SetUserID(userID string) *InvalidateUserSSOSessionsOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *InvalidateUserSSOSessionsOptions) SetHeaders(param map[string]string) *InvalidateUserSSOSessionsOptions {
	options.Headers = param
	return options
}

// ListApplicationsOptions : The ListApplications options.
type ListApplicationsOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListApplicationsOptions : Instantiate ListApplicationsOptions
func (*AppIDManagementV4) NewListApplicationsOptions(tenantID string) *ListApplicationsOptions {
	return &ListApplicationsOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ListApplicationsOptions) SetTenantID(tenantID string) *ListApplicationsOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListApplicationsOptions) SetHeaders(param map[string]string) *ListApplicationsOptions {
	options.Headers = param
	return options
}

// ListChannelsOptions : The ListChannels options.
type ListChannelsOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListChannelsOptions : Instantiate ListChannelsOptions
func (*AppIDManagementV4) NewListChannelsOptions(tenantID string) *ListChannelsOptions {
	return &ListChannelsOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ListChannelsOptions) SetTenantID(tenantID string) *ListChannelsOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListChannelsOptions) SetHeaders(param map[string]string) *ListChannelsOptions {
	options.Headers = param
	return options
}

// ListCloudDirectoryUsersOptions : The ListCloudDirectoryUsers options.
type ListCloudDirectoryUsersOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The first result in a set list of results.
	StartIndex *int64

	// The maximum number of results per page.
	Count *int64

	// Filter users by identity field.
	Query *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCloudDirectoryUsersOptions : Instantiate ListCloudDirectoryUsersOptions
func (*AppIDManagementV4) NewListCloudDirectoryUsersOptions(tenantID string) *ListCloudDirectoryUsersOptions {
	return &ListCloudDirectoryUsersOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ListCloudDirectoryUsersOptions) SetTenantID(tenantID string) *ListCloudDirectoryUsersOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetStartIndex : Allow user to set StartIndex
func (options *ListCloudDirectoryUsersOptions) SetStartIndex(startIndex int64) *ListCloudDirectoryUsersOptions {
	options.StartIndex = core.Int64Ptr(startIndex)
	return options
}

// SetCount : Allow user to set Count
func (options *ListCloudDirectoryUsersOptions) SetCount(count int64) *ListCloudDirectoryUsersOptions {
	options.Count = core.Int64Ptr(count)
	return options
}

// SetQuery : Allow user to set Query
func (options *ListCloudDirectoryUsersOptions) SetQuery(query string) *ListCloudDirectoryUsersOptions {
	options.Query = core.StringPtr(query)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListCloudDirectoryUsersOptions) SetHeaders(param map[string]string) *ListCloudDirectoryUsersOptions {
	options.Headers = param
	return options
}

// ListRolesOptions : The ListRoles options.
type ListRolesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRolesOptions : Instantiate ListRolesOptions
func (*AppIDManagementV4) NewListRolesOptions(tenantID string) *ListRolesOptions {
	return &ListRolesOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ListRolesOptions) SetTenantID(tenantID string) *ListRolesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListRolesOptions) SetHeaders(param map[string]string) *ListRolesOptions {
	options.Headers = param
	return options
}

// MFAChannelsListChannelsItem : MFAChannelsListChannelsItem struct
type MFAChannelsListChannelsItem struct {
	Type *string `json:"type" validate:"required"`

	IsActive *bool `json:"isActive" validate:"required"`

	Config *MFAChannelsListChannelsItemConfig `json:"config,omitempty"`
}

// UnmarshalMFAChannelsListChannelsItem unmarshals an instance of MFAChannelsListChannelsItem from the specified map of raw messages.
func UnmarshalMFAChannelsListChannelsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MFAChannelsListChannelsItem)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalMFAChannelsListChannelsItemConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MFAChannelsListChannelsItemConfig : MFAChannelsListChannelsItemConfig struct
type MFAChannelsListChannelsItemConfig struct {
	Key *string `json:"key,omitempty"`

	Secret *string `json:"secret,omitempty"`

	From *string `json:"from,omitempty"`

	Provider *string `json:"provider,omitempty"`
}

// UnmarshalMFAChannelsListChannelsItemConfig unmarshals an instance of MFAChannelsListChannelsItemConfig from the specified map of raw messages.
func UnmarshalMFAChannelsListChannelsItemConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MFAChannelsListChannelsItemConfig)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret", &obj.Secret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "from", &obj.From)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MFAChannelsList : MFAChannelsList struct
type MFAChannelsList struct {
	Channels []MFAChannelsListChannelsItem `json:"channels" validate:"required"`
}

// UnmarshalMFAChannelsList unmarshals an instance of MFAChannelsList from the specified map of raw messages.
func UnmarshalMFAChannelsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MFAChannelsList)
	err = core.UnmarshalModel(m, "channels", &obj.Channels, UnmarshalMFAChannelsListChannelsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostEmailDispatcherTestOptions : The PostEmailDispatcherTest options.
type PostEmailDispatcherTestOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The email address where you want to send your test message.
	Email *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostEmailDispatcherTestOptions : Instantiate PostEmailDispatcherTestOptions
func (*AppIDManagementV4) NewPostEmailDispatcherTestOptions(tenantID string, email string) *PostEmailDispatcherTestOptions {
	return &PostEmailDispatcherTestOptions{
		TenantID: core.StringPtr(tenantID),
		Email:    core.StringPtr(email),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostEmailDispatcherTestOptions) SetTenantID(tenantID string) *PostEmailDispatcherTestOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetEmail : Allow user to set Email
func (options *PostEmailDispatcherTestOptions) SetEmail(email string) *PostEmailDispatcherTestOptions {
	options.Email = core.StringPtr(email)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostEmailDispatcherTestOptions) SetHeaders(param map[string]string) *PostEmailDispatcherTestOptions {
	options.Headers = param
	return options
}

// PostExtensionsTestOptions : The PostExtensionsTest options.
type PostExtensionsTestOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The name of the extension.
	Name *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostExtensionsTestOptions.Name property.
// The name of the extension.
const (
	PostExtensionsTestOptionsNamePostmfaConst = "postmfa"
	PostExtensionsTestOptionsNamePremfaConst  = "premfa"
)

// NewPostExtensionsTestOptions : Instantiate PostExtensionsTestOptions
func (*AppIDManagementV4) NewPostExtensionsTestOptions(tenantID string, name string) *PostExtensionsTestOptions {
	return &PostExtensionsTestOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostExtensionsTestOptions) SetTenantID(tenantID string) *PostExtensionsTestOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *PostExtensionsTestOptions) SetName(name string) *PostExtensionsTestOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostExtensionsTestOptions) SetHeaders(param map[string]string) *PostExtensionsTestOptions {
	options.Headers = param
	return options
}

// PostMediaOptions : The PostMedia options.
type PostMediaOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of media. You can upload JPG or PNG files.
	MediaType *string `validate:"required"`

	// The image file. The recommended size is 320x320 px. The maxmimum files size is 100kb.
	File io.ReadCloser `validate:"required"`

	// The content type of file.
	FileContentType *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the PostMediaOptions.MediaType property.
// The type of media. You can upload JPG or PNG files.
const (
	PostMediaOptionsMediaTypeLogoConst = "logo"
)

// NewPostMediaOptions : Instantiate PostMediaOptions
func (*AppIDManagementV4) NewPostMediaOptions(tenantID string, mediaType string, file io.ReadCloser) *PostMediaOptions {
	return &PostMediaOptions{
		TenantID:  core.StringPtr(tenantID),
		MediaType: core.StringPtr(mediaType),
		File:      file,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostMediaOptions) SetTenantID(tenantID string) *PostMediaOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetMediaType : Allow user to set MediaType
func (options *PostMediaOptions) SetMediaType(mediaType string) *PostMediaOptions {
	options.MediaType = core.StringPtr(mediaType)
	return options
}

// SetFile : Allow user to set File
func (options *PostMediaOptions) SetFile(file io.ReadCloser) *PostMediaOptions {
	options.File = file
	return options
}

// SetFileContentType : Allow user to set FileContentType
func (options *PostMediaOptions) SetFileContentType(fileContentType string) *PostMediaOptions {
	options.FileContentType = core.StringPtr(fileContentType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostMediaOptions) SetHeaders(param map[string]string) *PostMediaOptions {
	options.Headers = param
	return options
}

// PostSMSDispatcherTestOptions : The PostSMSDispatcherTest options.
type PostSMSDispatcherTestOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	PhoneNumber *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostSMSDispatcherTestOptions : Instantiate PostSMSDispatcherTestOptions
func (*AppIDManagementV4) NewPostSMSDispatcherTestOptions(tenantID string, phoneNumber string) *PostSMSDispatcherTestOptions {
	return &PostSMSDispatcherTestOptions{
		TenantID:    core.StringPtr(tenantID),
		PhoneNumber: core.StringPtr(phoneNumber),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostSMSDispatcherTestOptions) SetTenantID(tenantID string) *PostSMSDispatcherTestOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetPhoneNumber : Allow user to set PhoneNumber
func (options *PostSMSDispatcherTestOptions) SetPhoneNumber(phoneNumber string) *PostSMSDispatcherTestOptions {
	options.PhoneNumber = core.StringPtr(phoneNumber)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostSMSDispatcherTestOptions) SetHeaders(param map[string]string) *PostSMSDispatcherTestOptions {
	options.Headers = param
	return options
}

// PostThemeColorOptions : The PostThemeColor options.
type PostThemeColorOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	HeaderColor *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostThemeColorOptions : Instantiate PostThemeColorOptions
func (*AppIDManagementV4) NewPostThemeColorOptions(tenantID string) *PostThemeColorOptions {
	return &PostThemeColorOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostThemeColorOptions) SetTenantID(tenantID string) *PostThemeColorOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetHeaderColor : Allow user to set HeaderColor
func (options *PostThemeColorOptions) SetHeaderColor(headerColor string) *PostThemeColorOptions {
	options.HeaderColor = core.StringPtr(headerColor)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostThemeColorOptions) SetHeaders(param map[string]string) *PostThemeColorOptions {
	options.Headers = param
	return options
}

// PostThemeTextOptions : The PostThemeText options.
type PostThemeTextOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	TabTitle *string

	Footnote *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostThemeTextOptions : Instantiate PostThemeTextOptions
func (*AppIDManagementV4) NewPostThemeTextOptions(tenantID string) *PostThemeTextOptions {
	return &PostThemeTextOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PostThemeTextOptions) SetTenantID(tenantID string) *PostThemeTextOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetTabTitle : Allow user to set TabTitle
func (options *PostThemeTextOptions) SetTabTitle(tabTitle string) *PostThemeTextOptions {
	options.TabTitle = core.StringPtr(tabTitle)
	return options
}

// SetFootnote : Allow user to set Footnote
func (options *PostThemeTextOptions) SetFootnote(footnote string) *PostThemeTextOptions {
	options.Footnote = core.StringPtr(footnote)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PostThemeTextOptions) SetHeaders(param map[string]string) *PostThemeTextOptions {
	options.Headers = param
	return options
}

// PutApplicationsRolesOptions : The PutApplicationsRoles options.
type PutApplicationsRolesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	Roles *UpdateUserRolesParamsRoles `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutApplicationsRolesOptions : Instantiate PutApplicationsRolesOptions
func (*AppIDManagementV4) NewPutApplicationsRolesOptions(tenantID string, clientID string, roles *UpdateUserRolesParamsRoles) *PutApplicationsRolesOptions {
	return &PutApplicationsRolesOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
		Roles:    roles,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PutApplicationsRolesOptions) SetTenantID(tenantID string) *PutApplicationsRolesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *PutApplicationsRolesOptions) SetClientID(clientID string) *PutApplicationsRolesOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetRoles : Allow user to set Roles
func (options *PutApplicationsRolesOptions) SetRoles(roles *UpdateUserRolesParamsRoles) *PutApplicationsRolesOptions {
	options.Roles = roles
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutApplicationsRolesOptions) SetHeaders(param map[string]string) *PutApplicationsRolesOptions {
	options.Headers = param
	return options
}

// PutApplicationsScopesOptions : The PutApplicationsScopes options.
type PutApplicationsScopesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	Scopes []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutApplicationsScopesOptions : Instantiate PutApplicationsScopesOptions
func (*AppIDManagementV4) NewPutApplicationsScopesOptions(tenantID string, clientID string) *PutApplicationsScopesOptions {
	return &PutApplicationsScopesOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PutApplicationsScopesOptions) SetTenantID(tenantID string) *PutApplicationsScopesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *PutApplicationsScopesOptions) SetClientID(clientID string) *PutApplicationsScopesOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetScopes : Allow user to set Scopes
func (options *PutApplicationsScopesOptions) SetScopes(scopes []string) *PutApplicationsScopesOptions {
	options.Scopes = scopes
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutApplicationsScopesOptions) SetHeaders(param map[string]string) *PutApplicationsScopesOptions {
	options.Headers = param
	return options
}

// PutTokensConfigOptions : The PutTokensConfig options.
type PutTokensConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IDTokenClaims []TokenClaimMapping

	AccessTokenClaims []TokenClaimMapping

	Access *AccessTokenConfigParams

	Refresh *TokenConfigParams

	AnonymousAccess *TokenConfigParams

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPutTokensConfigOptions : Instantiate PutTokensConfigOptions
func (*AppIDManagementV4) NewPutTokensConfigOptions(tenantID string) *PutTokensConfigOptions {
	return &PutTokensConfigOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *PutTokensConfigOptions) SetTenantID(tenantID string) *PutTokensConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIDTokenClaims : Allow user to set IDTokenClaims
func (options *PutTokensConfigOptions) SetIDTokenClaims(idTokenClaims []TokenClaimMapping) *PutTokensConfigOptions {
	options.IDTokenClaims = idTokenClaims
	return options
}

// SetAccessTokenClaims : Allow user to set AccessTokenClaims
func (options *PutTokensConfigOptions) SetAccessTokenClaims(accessTokenClaims []TokenClaimMapping) *PutTokensConfigOptions {
	options.AccessTokenClaims = accessTokenClaims
	return options
}

// SetAccess : Allow user to set Access
func (options *PutTokensConfigOptions) SetAccess(access *AccessTokenConfigParams) *PutTokensConfigOptions {
	options.Access = access
	return options
}

// SetRefresh : Allow user to set Refresh
func (options *PutTokensConfigOptions) SetRefresh(refresh *TokenConfigParams) *PutTokensConfigOptions {
	options.Refresh = refresh
	return options
}

// SetAnonymousAccess : Allow user to set AnonymousAccess
func (options *PutTokensConfigOptions) SetAnonymousAccess(anonymousAccess *TokenConfigParams) *PutTokensConfigOptions {
	options.AnonymousAccess = anonymousAccess
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PutTokensConfigOptions) SetHeaders(param map[string]string) *PutTokensConfigOptions {
	options.Headers = param
	return options
}

// RegisterApplicationOptions : The RegisterApplication options.
type RegisterApplicationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application name to be registered. Application name cannot exceed 50 characters.
	Name *string `validate:"required"`

	// The type of application to be registered. Allowed types are regularwebapp and singlepageapp.
	Type *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRegisterApplicationOptions : Instantiate RegisterApplicationOptions
func (*AppIDManagementV4) NewRegisterApplicationOptions(tenantID string, name string) *RegisterApplicationOptions {
	return &RegisterApplicationOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *RegisterApplicationOptions) SetTenantID(tenantID string) *RegisterApplicationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *RegisterApplicationOptions) SetName(name string) *RegisterApplicationOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetType : Allow user to set Type
func (options *RegisterApplicationOptions) SetType(typeVar string) *RegisterApplicationOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *RegisterApplicationOptions) SetHeaders(param map[string]string) *RegisterApplicationOptions {
	options.Headers = param
	return options
}

// ResendNotificationOptions : The ResendNotification options.
type ResendNotificationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED" or "RESET_PASSWORD".
	TemplateName *string `validate:"required,ne="`

	// The Cloud Directory unique user Id.
	UUID *string `validate:"required"`

	// Preferred language for resource. Format as described at RFC5646.
	Language *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ResendNotificationOptions.TemplateName property.
// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED" or "RESET_PASSWORD".
const (
	ResendNotificationOptionsTemplateNamePasswordChangedConst  = "PASSWORD_CHANGED"
	ResendNotificationOptionsTemplateNameResetPasswordConst    = "RESET_PASSWORD"
	ResendNotificationOptionsTemplateNameUserVerificationConst = "USER_VERIFICATION"
	ResendNotificationOptionsTemplateNameWelcomeConst          = "WELCOME"
)

// NewResendNotificationOptions : Instantiate ResendNotificationOptions
func (*AppIDManagementV4) NewResendNotificationOptions(tenantID string, templateName string, uuid string) *ResendNotificationOptions {
	return &ResendNotificationOptions{
		TenantID:     core.StringPtr(tenantID),
		TemplateName: core.StringPtr(templateName),
		UUID:         core.StringPtr(uuid),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *ResendNotificationOptions) SetTenantID(tenantID string) *ResendNotificationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetTemplateName : Allow user to set TemplateName
func (options *ResendNotificationOptions) SetTemplateName(templateName string) *ResendNotificationOptions {
	options.TemplateName = core.StringPtr(templateName)
	return options
}

// SetUUID : Allow user to set UUID
func (options *ResendNotificationOptions) SetUUID(uuid string) *ResendNotificationOptions {
	options.UUID = core.StringPtr(uuid)
	return options
}

// SetLanguage : Allow user to set Language
func (options *ResendNotificationOptions) SetLanguage(language string) *ResendNotificationOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ResendNotificationOptions) SetHeaders(param map[string]string) *ResendNotificationOptions {
	options.Headers = param
	return options
}

// ResendNotificationResponse : ResendNotificationResponse struct
type ResendNotificationResponse struct {
	Message *string `json:"message,omitempty"`
}

// UnmarshalResendNotificationResponse unmarshals an instance of ResendNotificationResponse from the specified map of raw messages.
func UnmarshalResendNotificationResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResendNotificationResponse)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RolesList : RolesList struct
type RolesList struct {
	Roles []RolesListRolesItem `json:"roles,omitempty"`
}

// UnmarshalRolesList unmarshals an instance of RolesList from the specified map of raw messages.
func UnmarshalRolesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RolesList)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalRolesListRolesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RolesListRolesItem : RolesListRolesItem struct
type RolesListRolesItem struct {
	ID *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`

	Description *string `json:"description,omitempty"`

	Access []RoleAccessItem `json:"access,omitempty"`
}

// UnmarshalRolesListRolesItem unmarshals an instance of RolesListRolesItem from the specified map of raw messages.
func UnmarshalRolesListRolesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RolesListRolesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "access", &obj.Access, UnmarshalRoleAccessItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLConfigParamsAuthnContext : SAMLConfigParamsAuthnContext struct
type SAMLConfigParamsAuthnContext struct {
	Class []string `json:"class,omitempty"`

	Comparison *string `json:"comparison,omitempty"`
}

// Constants associated with the SAMLConfigParamsAuthnContext.Class property.
const (
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesAuthenticatedtelephonyConst      = "urn:oasis:names:tc:SAML:2.0:ac:classes:AuthenticatedTelephony"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesInternetprotocolConst            = "urn:oasis:names:tc:SAML:2.0:ac:classes:InternetProtocol"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesInternetprotocolpasswordConst    = "urn:oasis:names:tc:SAML:2.0:ac:classes:InternetProtocolPassword"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesKerberosConst                    = "urn:oasis:names:tc:SAML:2.0:ac:classes:Kerberos"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesMobileonefactorcontractConst     = "urn:oasis:names:tc:SAML:2.0:ac:classes:MobileOneFactorContract"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesMobileonefactorunregisteredConst = "urn:oasis:names:tc:SAML:2.0:ac:classes:MobileOneFactorUnregistered"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesMobiletwofactorcontractConst     = "urn:oasis:names:tc:SAML:2.0:ac:classes:MobileTwoFactorContract"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesMobiletwofactorunregisteredConst = "urn:oasis:names:tc:SAML:2.0:ac:classes:MobileTwoFactorUnregistered"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesNomadtelephonyConst              = "urn:oasis:names:tc:SAML:2.0:ac:classes:NomadTelephony"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesPasswordConst                    = "urn:oasis:names:tc:SAML:2.0:ac:classes:Password"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesPasswordprotectedtransportConst  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesPersonaltelephonyConst           = "urn:oasis:names:tc:SAML:2.0:ac:classes:PersonalTelephony"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesPgpConst                         = "urn:oasis:names:tc:SAML:2.0:ac:classes:PGP"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesPrevioussessionConst             = "urn:oasis:names:tc:SAML:2.0:ac:classes:PreviousSession"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesSecureremotepasswordConst        = "urn:oasis:names:tc:SAML:2.0:ac:classes:SecureRemotePassword"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesSmartcardConst                   = "urn:oasis:names:tc:SAML:2.0:ac:classes:Smartcard"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesSmartcardpkiConst                = "urn:oasis:names:tc:SAML:2.0:ac:classes:SmartcardPKI"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesSoftwarepkiConst                 = "urn:oasis:names:tc:SAML:2.0:ac:classes:SoftwarePKI"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesSpkiConst                        = "urn:oasis:names:tc:SAML:2.0:ac:classes:SPKI"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesTelephonyConst                   = "urn:oasis:names:tc:SAML:2.0:ac:classes:Telephony"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesTimesynctokenConst               = "urn:oasis:names:tc:SAML:2.0:ac:classes:TimeSyncToken"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesTlsclientConst                   = "urn:oasis:names:tc:SAML:2.0:ac:classes:TLSClient"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesUnspecifiedConst                 = "urn:oasis:names:tc:SAML:2.0:ac:classes:unspecified"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesX509Const                        = "urn:oasis:names:tc:SAML:2.0:ac:classes:X509"
	SAMLConfigParamsAuthnContextClassUrnOasisNamesTcSaml20AcClassesXmldsigConst                     = "urn:oasis:names:tc:SAML:2.0:ac:classes:XMLDSig"
)

// Constants associated with the SAMLConfigParamsAuthnContext.Comparison property.
const (
	SAMLConfigParamsAuthnContextComparisonBetterConst  = "better"
	SAMLConfigParamsAuthnContextComparisonExactConst   = "exact"
	SAMLConfigParamsAuthnContextComparisonMaximumConst = "maximum"
	SAMLConfigParamsAuthnContextComparisonMinimumConst = "minimum"
)

// UnmarshalSAMLConfigParamsAuthnContext unmarshals an instance of SAMLConfigParamsAuthnContext from the specified map of raw messages.
func UnmarshalSAMLConfigParamsAuthnContext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLConfigParamsAuthnContext)
	err = core.UnmarshalPrimitive(m, "class", &obj.Class)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "comparison", &obj.Comparison)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLResponseWithValidationDataValidationData : SAMLResponseWithValidationDataValidationData struct
type SAMLResponseWithValidationDataValidationData struct {
	Certificates []SAMLResponseWithValidationDataValidationDataCertificatesItem `json:"certificates" validate:"required"`
}

// UnmarshalSAMLResponseWithValidationDataValidationData unmarshals an instance of SAMLResponseWithValidationDataValidationData from the specified map of raw messages.
func UnmarshalSAMLResponseWithValidationDataValidationData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLResponseWithValidationDataValidationData)
	err = core.UnmarshalModel(m, "certificates", &obj.Certificates, UnmarshalSAMLResponseWithValidationDataValidationDataCertificatesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLResponseWithValidationDataValidationDataCertificatesItem : SAMLResponseWithValidationDataValidationDataCertificatesItem struct
type SAMLResponseWithValidationDataValidationDataCertificatesItem struct {
	CertificateIndex *int64 `json:"certificate_index" validate:"required"`

	ExpirationTimestamp *int64 `json:"expiration_timestamp" validate:"required"`

	Warning *string `json:"warning,omitempty"`
}

// UnmarshalSAMLResponseWithValidationDataValidationDataCertificatesItem unmarshals an instance of SAMLResponseWithValidationDataValidationDataCertificatesItem from the specified map of raw messages.
func UnmarshalSAMLResponseWithValidationDataValidationDataCertificatesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLResponseWithValidationDataValidationDataCertificatesItem)
	err = core.UnmarshalPrimitive(m, "certificate_index", &obj.CertificateIndex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_timestamp", &obj.ExpirationTimestamp)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "warning", &obj.Warning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetAuditStatusOptions : The SetAuditStatus options.
type SetAuditStatusOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetAuditStatusOptions : Instantiate SetAuditStatusOptions
func (*AppIDManagementV4) NewSetAuditStatusOptions(tenantID string, isActive bool) *SetAuditStatusOptions {
	return &SetAuditStatusOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetAuditStatusOptions) SetTenantID(tenantID string) *SetAuditStatusOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *SetAuditStatusOptions) SetIsActive(isActive bool) *SetAuditStatusOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetAuditStatusOptions) SetHeaders(param map[string]string) *SetAuditStatusOptions {
	options.Headers = param
	return options
}

// SetCloudDirectoryActionOptions : The SetCloudDirectoryAction options.
type SetCloudDirectoryActionOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
	// your custom reset password page.
	Action *string `validate:"required,ne="`

	// The action URL.
	ActionURL *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the SetCloudDirectoryActionOptions.Action property.
// The type of the action. on_user_verified - the URL of your custom user verified page, on_reset_password - the URL of
// your custom reset password page.
const (
	SetCloudDirectoryActionOptionsActionOnResetPasswordConst = "on_reset_password"
	SetCloudDirectoryActionOptionsActionOnUserVerifiedConst  = "on_user_verified"
)

// NewSetCloudDirectoryActionOptions : Instantiate SetCloudDirectoryActionOptions
func (*AppIDManagementV4) NewSetCloudDirectoryActionOptions(tenantID string, action string, actionURL string) *SetCloudDirectoryActionOptions {
	return &SetCloudDirectoryActionOptions{
		TenantID:  core.StringPtr(tenantID),
		Action:    core.StringPtr(action),
		ActionURL: core.StringPtr(actionURL),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectoryActionOptions) SetTenantID(tenantID string) *SetCloudDirectoryActionOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetAction : Allow user to set Action
func (options *SetCloudDirectoryActionOptions) SetAction(action string) *SetCloudDirectoryActionOptions {
	options.Action = core.StringPtr(action)
	return options
}

// SetActionURL : Allow user to set ActionURL
func (options *SetCloudDirectoryActionOptions) SetActionURL(actionURL string) *SetCloudDirectoryActionOptions {
	options.ActionURL = core.StringPtr(actionURL)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectoryActionOptions) SetHeaders(param map[string]string) *SetCloudDirectoryActionOptions {
	options.Headers = param
	return options
}

// SetCloudDirectoryAdvancedPasswordManagementOptions : The SetCloudDirectoryAdvancedPasswordManagement options.
type SetCloudDirectoryAdvancedPasswordManagementOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	AdvancedPasswordManagement *ApmSchemaAdvancedPasswordManagement `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetCloudDirectoryAdvancedPasswordManagementOptions : Instantiate SetCloudDirectoryAdvancedPasswordManagementOptions
func (*AppIDManagementV4) NewSetCloudDirectoryAdvancedPasswordManagementOptions(tenantID string, advancedPasswordManagement *ApmSchemaAdvancedPasswordManagement) *SetCloudDirectoryAdvancedPasswordManagementOptions {
	return &SetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID:                   core.StringPtr(tenantID),
		AdvancedPasswordManagement: advancedPasswordManagement,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectoryAdvancedPasswordManagementOptions) SetTenantID(tenantID string) *SetCloudDirectoryAdvancedPasswordManagementOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetAdvancedPasswordManagement : Allow user to set AdvancedPasswordManagement
func (options *SetCloudDirectoryAdvancedPasswordManagementOptions) SetAdvancedPasswordManagement(advancedPasswordManagement *ApmSchemaAdvancedPasswordManagement) *SetCloudDirectoryAdvancedPasswordManagementOptions {
	options.AdvancedPasswordManagement = advancedPasswordManagement
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectoryAdvancedPasswordManagementOptions) SetHeaders(param map[string]string) *SetCloudDirectoryAdvancedPasswordManagementOptions {
	options.Headers = param
	return options
}

// SetCloudDirectoryEmailDispatcherOptions : The SetCloudDirectoryEmailDispatcher options.
type SetCloudDirectoryEmailDispatcherOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Provider *string `validate:"required"`

	Sendgrid *EmailDispatcherParamsSendgrid

	Custom *EmailDispatcherParamsCustom

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the SetCloudDirectoryEmailDispatcherOptions.Provider property.
const (
	SetCloudDirectoryEmailDispatcherOptionsProviderAppidConst    = "appid"
	SetCloudDirectoryEmailDispatcherOptionsProviderCustomConst   = "custom"
	SetCloudDirectoryEmailDispatcherOptionsProviderSendgridConst = "sendgrid"
)

// NewSetCloudDirectoryEmailDispatcherOptions : Instantiate SetCloudDirectoryEmailDispatcherOptions
func (*AppIDManagementV4) NewSetCloudDirectoryEmailDispatcherOptions(tenantID string, provider string) *SetCloudDirectoryEmailDispatcherOptions {
	return &SetCloudDirectoryEmailDispatcherOptions{
		TenantID: core.StringPtr(tenantID),
		Provider: core.StringPtr(provider),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectoryEmailDispatcherOptions) SetTenantID(tenantID string) *SetCloudDirectoryEmailDispatcherOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetProvider : Allow user to set Provider
func (options *SetCloudDirectoryEmailDispatcherOptions) SetProvider(provider string) *SetCloudDirectoryEmailDispatcherOptions {
	options.Provider = core.StringPtr(provider)
	return options
}

// SetSendgrid : Allow user to set Sendgrid
func (options *SetCloudDirectoryEmailDispatcherOptions) SetSendgrid(sendgrid *EmailDispatcherParamsSendgrid) *SetCloudDirectoryEmailDispatcherOptions {
	options.Sendgrid = sendgrid
	return options
}

// SetCustom : Allow user to set Custom
func (options *SetCloudDirectoryEmailDispatcherOptions) SetCustom(custom *EmailDispatcherParamsCustom) *SetCloudDirectoryEmailDispatcherOptions {
	options.Custom = custom
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectoryEmailDispatcherOptions) SetHeaders(param map[string]string) *SetCloudDirectoryEmailDispatcherOptions {
	options.Headers = param
	return options
}

// SetCloudDirectoryIDPOptions : The SetCloudDirectoryIDP options.
type SetCloudDirectoryIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config *CloudDirectoryConfigParams `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetCloudDirectoryIDPOptions : Instantiate SetCloudDirectoryIDPOptions
func (*AppIDManagementV4) NewSetCloudDirectoryIDPOptions(tenantID string, isActive bool, config *CloudDirectoryConfigParams) *SetCloudDirectoryIDPOptions {
	return &SetCloudDirectoryIDPOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
		Config:   config,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectoryIDPOptions) SetTenantID(tenantID string) *SetCloudDirectoryIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *SetCloudDirectoryIDPOptions) SetIsActive(isActive bool) *SetCloudDirectoryIDPOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *SetCloudDirectoryIDPOptions) SetConfig(config *CloudDirectoryConfigParams) *SetCloudDirectoryIDPOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectoryIDPOptions) SetHeaders(param map[string]string) *SetCloudDirectoryIDPOptions {
	options.Headers = param
	return options
}

// SetCloudDirectoryPasswordRegexOptions : The SetCloudDirectoryPasswordRegex options.
type SetCloudDirectoryPasswordRegexOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Regex *string

	Base64EncodedRegex *string

	ErrorMessage *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetCloudDirectoryPasswordRegexOptions : Instantiate SetCloudDirectoryPasswordRegexOptions
func (*AppIDManagementV4) NewSetCloudDirectoryPasswordRegexOptions(tenantID string) *SetCloudDirectoryPasswordRegexOptions {
	return &SetCloudDirectoryPasswordRegexOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectoryPasswordRegexOptions) SetTenantID(tenantID string) *SetCloudDirectoryPasswordRegexOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetRegex : Allow user to set Regex
func (options *SetCloudDirectoryPasswordRegexOptions) SetRegex(regex string) *SetCloudDirectoryPasswordRegexOptions {
	options.Regex = core.StringPtr(regex)
	return options
}

// SetBase64EncodedRegex : Allow user to set Base64EncodedRegex
func (options *SetCloudDirectoryPasswordRegexOptions) SetBase64EncodedRegex(base64EncodedRegex string) *SetCloudDirectoryPasswordRegexOptions {
	options.Base64EncodedRegex = core.StringPtr(base64EncodedRegex)
	return options
}

// SetErrorMessage : Allow user to set ErrorMessage
func (options *SetCloudDirectoryPasswordRegexOptions) SetErrorMessage(errorMessage string) *SetCloudDirectoryPasswordRegexOptions {
	options.ErrorMessage = core.StringPtr(errorMessage)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectoryPasswordRegexOptions) SetHeaders(param map[string]string) *SetCloudDirectoryPasswordRegexOptions {
	options.Headers = param
	return options
}

// SetCloudDirectorySenderDetailsOptions : The SetCloudDirectorySenderDetails options.
type SetCloudDirectorySenderDetailsOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	SenderDetails *CloudDirectorySenderDetailsSenderDetails `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetCloudDirectorySenderDetailsOptions : Instantiate SetCloudDirectorySenderDetailsOptions
func (*AppIDManagementV4) NewSetCloudDirectorySenderDetailsOptions(tenantID string, senderDetails *CloudDirectorySenderDetailsSenderDetails) *SetCloudDirectorySenderDetailsOptions {
	return &SetCloudDirectorySenderDetailsOptions{
		TenantID:      core.StringPtr(tenantID),
		SenderDetails: senderDetails,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCloudDirectorySenderDetailsOptions) SetTenantID(tenantID string) *SetCloudDirectorySenderDetailsOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetSenderDetails : Allow user to set SenderDetails
func (options *SetCloudDirectorySenderDetailsOptions) SetSenderDetails(senderDetails *CloudDirectorySenderDetailsSenderDetails) *SetCloudDirectorySenderDetailsOptions {
	options.SenderDetails = senderDetails
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCloudDirectorySenderDetailsOptions) SetHeaders(param map[string]string) *SetCloudDirectorySenderDetailsOptions {
	options.Headers = param
	return options
}

// SetCustomIDPOptions : The SetCustomIDP options.
type SetCustomIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config *CustomIDPConfigParamsConfig

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetCustomIDPOptions : Instantiate SetCustomIDPOptions
func (*AppIDManagementV4) NewSetCustomIDPOptions(tenantID string, isActive bool) *SetCustomIDPOptions {
	return &SetCustomIDPOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetCustomIDPOptions) SetTenantID(tenantID string) *SetCustomIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *SetCustomIDPOptions) SetIsActive(isActive bool) *SetCustomIDPOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *SetCustomIDPOptions) SetConfig(config *CustomIDPConfigParamsConfig) *SetCustomIDPOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetCustomIDPOptions) SetHeaders(param map[string]string) *SetCustomIDPOptions {
	options.Headers = param
	return options
}

// SetFacebookIDPOptions : The SetFacebookIDP options.
type SetFacebookIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The identity provider configuration as a JSON object. If the configuration is not set, IBM default credentials are
	// used.
	IDP *FacebookGoogleConfigParams `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetFacebookIDPOptions : Instantiate SetFacebookIDPOptions
func (*AppIDManagementV4) NewSetFacebookIDPOptions(tenantID string, idp *FacebookGoogleConfigParams) *SetFacebookIDPOptions {
	return &SetFacebookIDPOptions{
		TenantID: core.StringPtr(tenantID),
		IDP:      idp,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetFacebookIDPOptions) SetTenantID(tenantID string) *SetFacebookIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIDP : Allow user to set IDP
func (options *SetFacebookIDPOptions) SetIDP(idp *FacebookGoogleConfigParams) *SetFacebookIDPOptions {
	options.IDP = idp
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetFacebookIDPOptions) SetHeaders(param map[string]string) *SetFacebookIDPOptions {
	options.Headers = param
	return options
}

// SetGoogleIDPOptions : The SetGoogleIDP options.
type SetGoogleIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The identity provider configuration as a JSON object. If the configuration is not set, IBM default credentials are
	// used.
	IDP *FacebookGoogleConfigParams `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetGoogleIDPOptions : Instantiate SetGoogleIDPOptions
func (*AppIDManagementV4) NewSetGoogleIDPOptions(tenantID string, idp *FacebookGoogleConfigParams) *SetGoogleIDPOptions {
	return &SetGoogleIDPOptions{
		TenantID: core.StringPtr(tenantID),
		IDP:      idp,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetGoogleIDPOptions) SetTenantID(tenantID string) *SetGoogleIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIDP : Allow user to set IDP
func (options *SetGoogleIDPOptions) SetIDP(idp *FacebookGoogleConfigParams) *SetGoogleIDPOptions {
	options.IDP = idp
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetGoogleIDPOptions) SetHeaders(param map[string]string) *SetGoogleIDPOptions {
	options.Headers = param
	return options
}

// SetSAMLIDPOptions : The SetSAMLIDP options.
type SetSAMLIDPOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config *SAMLConfigParams

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetSAMLIDPOptions : Instantiate SetSAMLIDPOptions
func (*AppIDManagementV4) NewSetSAMLIDPOptions(tenantID string, isActive bool) *SetSAMLIDPOptions {
	return &SetSAMLIDPOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *SetSAMLIDPOptions) SetTenantID(tenantID string) *SetSAMLIDPOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *SetSAMLIDPOptions) SetIsActive(isActive bool) *SetSAMLIDPOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *SetSAMLIDPOptions) SetConfig(config *SAMLConfigParams) *SetSAMLIDPOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *SetSAMLIDPOptions) SetHeaders(param map[string]string) *SetSAMLIDPOptions {
	options.Headers = param
	return options
}

// StartForgotPasswordOptions : The StartForgotPassword options.
type StartForgotPasswordOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user identitier - email or userName based on the identityField property in <a
	// href="#!/Identity_Providers/set_cloud_directory_idp" target="_blank"> cloud directory configuration.</a>.
	User *string `validate:"required"`

	// Preferred language for resource. Format as described at RFC5646.
	Language *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewStartForgotPasswordOptions : Instantiate StartForgotPasswordOptions
func (*AppIDManagementV4) NewStartForgotPasswordOptions(tenantID string, user string) *StartForgotPasswordOptions {
	return &StartForgotPasswordOptions{
		TenantID: core.StringPtr(tenantID),
		User:     core.StringPtr(user),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *StartForgotPasswordOptions) SetTenantID(tenantID string) *StartForgotPasswordOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUser : Allow user to set User
func (options *StartForgotPasswordOptions) SetUser(user string) *StartForgotPasswordOptions {
	options.User = core.StringPtr(user)
	return options
}

// SetLanguage : Allow user to set Language
func (options *StartForgotPasswordOptions) SetLanguage(language string) *StartForgotPasswordOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *StartForgotPasswordOptions) SetHeaders(param map[string]string) *StartForgotPasswordOptions {
	options.Headers = param
	return options
}

// StartSignUpOptions : The StartSignUp options.
type StartSignUpOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// A boolean indication if a profile should be created for the Cloud Directory user.
	ShouldCreateProfile *bool `validate:"required"`

	Emails []CreateNewUserEmailsItem `validate:"required"`

	Password *string `validate:"required"`

	Active *bool

	LockedUntil *int64

	DisplayName *string

	UserName *string

	// Accepted values "PENDING" or "CONFIRMED".
	Status *string

	// Preferred language for resource. Format as described at RFC5646.
	Language *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewStartSignUpOptions : Instantiate StartSignUpOptions
func (*AppIDManagementV4) NewStartSignUpOptions(tenantID string, shouldCreateProfile bool, emails []CreateNewUserEmailsItem, password string) *StartSignUpOptions {
	return &StartSignUpOptions{
		TenantID:            core.StringPtr(tenantID),
		ShouldCreateProfile: core.BoolPtr(shouldCreateProfile),
		Emails:              emails,
		Password:            core.StringPtr(password),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *StartSignUpOptions) SetTenantID(tenantID string) *StartSignUpOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetShouldCreateProfile : Allow user to set ShouldCreateProfile
func (options *StartSignUpOptions) SetShouldCreateProfile(shouldCreateProfile bool) *StartSignUpOptions {
	options.ShouldCreateProfile = core.BoolPtr(shouldCreateProfile)
	return options
}

// SetEmails : Allow user to set Emails
func (options *StartSignUpOptions) SetEmails(emails []CreateNewUserEmailsItem) *StartSignUpOptions {
	options.Emails = emails
	return options
}

// SetPassword : Allow user to set Password
func (options *StartSignUpOptions) SetPassword(password string) *StartSignUpOptions {
	options.Password = core.StringPtr(password)
	return options
}

// SetActive : Allow user to set Active
func (options *StartSignUpOptions) SetActive(active bool) *StartSignUpOptions {
	options.Active = core.BoolPtr(active)
	return options
}

// SetLockedUntil : Allow user to set LockedUntil
func (options *StartSignUpOptions) SetLockedUntil(lockedUntil int64) *StartSignUpOptions {
	options.LockedUntil = core.Int64Ptr(lockedUntil)
	return options
}

// SetDisplayName : Allow user to set DisplayName
func (options *StartSignUpOptions) SetDisplayName(displayName string) *StartSignUpOptions {
	options.DisplayName = core.StringPtr(displayName)
	return options
}

// SetUserName : Allow user to set UserName
func (options *StartSignUpOptions) SetUserName(userName string) *StartSignUpOptions {
	options.UserName = core.StringPtr(userName)
	return options
}

// SetStatus : Allow user to set Status
func (options *StartSignUpOptions) SetStatus(status string) *StartSignUpOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetLanguage : Allow user to set Language
func (options *StartSignUpOptions) SetLanguage(language string) *StartSignUpOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *StartSignUpOptions) SetHeaders(param map[string]string) *StartSignUpOptions {
	options.Headers = param
	return options
}

// UpdateApplicationOptions : The UpdateApplication options.
type UpdateApplicationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The application clientId.
	ClientID *string `validate:"required,ne="`

	// The application name to be updated. Application name cannot exceed 50 characters.
	Name *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateApplicationOptions : Instantiate UpdateApplicationOptions
func (*AppIDManagementV4) NewUpdateApplicationOptions(tenantID string, clientID string, name string) *UpdateApplicationOptions {
	return &UpdateApplicationOptions{
		TenantID: core.StringPtr(tenantID),
		ClientID: core.StringPtr(clientID),
		Name:     core.StringPtr(name),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateApplicationOptions) SetTenantID(tenantID string) *UpdateApplicationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *UpdateApplicationOptions) SetClientID(clientID string) *UpdateApplicationOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateApplicationOptions) SetName(name string) *UpdateApplicationOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateApplicationOptions) SetHeaders(param map[string]string) *UpdateApplicationOptions {
	options.Headers = param
	return options
}

// UpdateChannelOptions : The UpdateChannel options.
type UpdateChannelOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The MFA channel.
	Channel *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config interface{}

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateChannelOptions.Channel property.
// The MFA channel.
const (
	UpdateChannelOptionsChannelEmailConst = "email"
	UpdateChannelOptionsChannelNexmoConst = "nexmo"
)

// NewUpdateChannelOptions : Instantiate UpdateChannelOptions
func (*AppIDManagementV4) NewUpdateChannelOptions(tenantID string, channel string, isActive bool) *UpdateChannelOptions {
	return &UpdateChannelOptions{
		TenantID: core.StringPtr(tenantID),
		Channel:  core.StringPtr(channel),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateChannelOptions) SetTenantID(tenantID string) *UpdateChannelOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetChannel : Allow user to set Channel
func (options *UpdateChannelOptions) SetChannel(channel string) *UpdateChannelOptions {
	options.Channel = core.StringPtr(channel)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateChannelOptions) SetIsActive(isActive bool) *UpdateChannelOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *UpdateChannelOptions) SetConfig(config interface{}) *UpdateChannelOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateChannelOptions) SetHeaders(param map[string]string) *UpdateChannelOptions {
	options.Headers = param
	return options
}

// UpdateCloudDirectoryUserOptions : The UpdateCloudDirectoryUser options.
type UpdateCloudDirectoryUserOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The ID assigned to a user when they sign in by using Cloud Directory.
	UserID *string `validate:"required,ne="`

	Emails []CreateNewUserEmailsItem `validate:"required"`

	Status *string

	DisplayName *string

	UserName *string

	Password *string

	Active *bool

	LockedUntil *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateCloudDirectoryUserOptions : Instantiate UpdateCloudDirectoryUserOptions
func (*AppIDManagementV4) NewUpdateCloudDirectoryUserOptions(tenantID string, userID string, emails []CreateNewUserEmailsItem) *UpdateCloudDirectoryUserOptions {
	return &UpdateCloudDirectoryUserOptions{
		TenantID: core.StringPtr(tenantID),
		UserID:   core.StringPtr(userID),
		Emails:   emails,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateCloudDirectoryUserOptions) SetTenantID(tenantID string) *UpdateCloudDirectoryUserOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUserID : Allow user to set UserID
func (options *UpdateCloudDirectoryUserOptions) SetUserID(userID string) *UpdateCloudDirectoryUserOptions {
	options.UserID = core.StringPtr(userID)
	return options
}

// SetEmails : Allow user to set Emails
func (options *UpdateCloudDirectoryUserOptions) SetEmails(emails []CreateNewUserEmailsItem) *UpdateCloudDirectoryUserOptions {
	options.Emails = emails
	return options
}

// SetStatus : Allow user to set Status
func (options *UpdateCloudDirectoryUserOptions) SetStatus(status string) *UpdateCloudDirectoryUserOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetDisplayName : Allow user to set DisplayName
func (options *UpdateCloudDirectoryUserOptions) SetDisplayName(displayName string) *UpdateCloudDirectoryUserOptions {
	options.DisplayName = core.StringPtr(displayName)
	return options
}

// SetUserName : Allow user to set UserName
func (options *UpdateCloudDirectoryUserOptions) SetUserName(userName string) *UpdateCloudDirectoryUserOptions {
	options.UserName = core.StringPtr(userName)
	return options
}

// SetPassword : Allow user to set Password
func (options *UpdateCloudDirectoryUserOptions) SetPassword(password string) *UpdateCloudDirectoryUserOptions {
	options.Password = core.StringPtr(password)
	return options
}

// SetActive : Allow user to set Active
func (options *UpdateCloudDirectoryUserOptions) SetActive(active bool) *UpdateCloudDirectoryUserOptions {
	options.Active = core.BoolPtr(active)
	return options
}

// SetLockedUntil : Allow user to set LockedUntil
func (options *UpdateCloudDirectoryUserOptions) SetLockedUntil(lockedUntil int64) *UpdateCloudDirectoryUserOptions {
	options.LockedUntil = core.Int64Ptr(lockedUntil)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCloudDirectoryUserOptions) SetHeaders(param map[string]string) *UpdateCloudDirectoryUserOptions {
	options.Headers = param
	return options
}

// UpdateExtensionActiveOptions : The UpdateExtensionActive options.
type UpdateExtensionActiveOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The name of the extension.
	Name *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config interface{}

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateExtensionActiveOptions.Name property.
// The name of the extension.
const (
	UpdateExtensionActiveOptionsNamePostmfaConst = "postmfa"
	UpdateExtensionActiveOptionsNamePremfaConst  = "premfa"
)

// NewUpdateExtensionActiveOptions : Instantiate UpdateExtensionActiveOptions
func (*AppIDManagementV4) NewUpdateExtensionActiveOptions(tenantID string, name string, isActive bool) *UpdateExtensionActiveOptions {
	return &UpdateExtensionActiveOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateExtensionActiveOptions) SetTenantID(tenantID string) *UpdateExtensionActiveOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateExtensionActiveOptions) SetName(name string) *UpdateExtensionActiveOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateExtensionActiveOptions) SetIsActive(isActive bool) *UpdateExtensionActiveOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *UpdateExtensionActiveOptions) SetConfig(config interface{}) *UpdateExtensionActiveOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateExtensionActiveOptions) SetHeaders(param map[string]string) *UpdateExtensionActiveOptions {
	options.Headers = param
	return options
}

// UpdateExtensionConfigConfig : UpdateExtensionConfigConfig struct
type UpdateExtensionConfigConfig struct {
	URL *string `json:"url,omitempty"`

	HeadersVar interface{} `json:"headers,omitempty"`
}

// UnmarshalUpdateExtensionConfigConfig unmarshals an instance of UpdateExtensionConfigConfig from the specified map of raw messages.
func UnmarshalUpdateExtensionConfigConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateExtensionConfigConfig)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "headers", &obj.HeadersVar)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateExtensionConfigOptions : The UpdateExtensionConfig options.
type UpdateExtensionConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The name of the extension.
	Name *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config *UpdateExtensionConfigConfig

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateExtensionConfigOptions.Name property.
// The name of the extension.
const (
	UpdateExtensionConfigOptionsNamePostmfaConst = "postmfa"
	UpdateExtensionConfigOptionsNamePremfaConst  = "premfa"
)

// NewUpdateExtensionConfigOptions : Instantiate UpdateExtensionConfigOptions
func (*AppIDManagementV4) NewUpdateExtensionConfigOptions(tenantID string, name string, isActive bool) *UpdateExtensionConfigOptions {
	return &UpdateExtensionConfigOptions{
		TenantID: core.StringPtr(tenantID),
		Name:     core.StringPtr(name),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateExtensionConfigOptions) SetTenantID(tenantID string) *UpdateExtensionConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateExtensionConfigOptions) SetName(name string) *UpdateExtensionConfigOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateExtensionConfigOptions) SetIsActive(isActive bool) *UpdateExtensionConfigOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *UpdateExtensionConfigOptions) SetConfig(config *UpdateExtensionConfigConfig) *UpdateExtensionConfigOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateExtensionConfigOptions) SetHeaders(param map[string]string) *UpdateExtensionConfigOptions {
	options.Headers = param
	return options
}

// UpdateLocalizationOptions : The UpdateLocalization options.
type UpdateLocalizationOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Languages []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateLocalizationOptions : Instantiate UpdateLocalizationOptions
func (*AppIDManagementV4) NewUpdateLocalizationOptions(tenantID string) *UpdateLocalizationOptions {
	return &UpdateLocalizationOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateLocalizationOptions) SetTenantID(tenantID string) *UpdateLocalizationOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetLanguages : Allow user to set Languages
func (options *UpdateLocalizationOptions) SetLanguages(languages []string) *UpdateLocalizationOptions {
	options.Languages = languages
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLocalizationOptions) SetHeaders(param map[string]string) *UpdateLocalizationOptions {
	options.Headers = param
	return options
}

// UpdateMFAConfigOptions : The UpdateMFAConfig options.
type UpdateMFAConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	Config interface{}

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateMFAConfigOptions : Instantiate UpdateMFAConfigOptions
func (*AppIDManagementV4) NewUpdateMFAConfigOptions(tenantID string, isActive bool) *UpdateMFAConfigOptions {
	return &UpdateMFAConfigOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateMFAConfigOptions) SetTenantID(tenantID string) *UpdateMFAConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateMFAConfigOptions) SetIsActive(isActive bool) *UpdateMFAConfigOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetConfig : Allow user to set Config
func (options *UpdateMFAConfigOptions) SetConfig(config interface{}) *UpdateMFAConfigOptions {
	options.Config = config
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateMFAConfigOptions) SetHeaders(param map[string]string) *UpdateMFAConfigOptions {
	options.Headers = param
	return options
}

// UpdateRateLimitConfigOptions : The UpdateRateLimitConfig options.
type UpdateRateLimitConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	SignUpLimitPerMinute *int64 `validate:"required"`

	SignInLimitPerMinute *int64 `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateRateLimitConfigOptions : Instantiate UpdateRateLimitConfigOptions
func (*AppIDManagementV4) NewUpdateRateLimitConfigOptions(tenantID string, signUpLimitPerMinute int64, signInLimitPerMinute int64) *UpdateRateLimitConfigOptions {
	return &UpdateRateLimitConfigOptions{
		TenantID:             core.StringPtr(tenantID),
		SignUpLimitPerMinute: core.Int64Ptr(signUpLimitPerMinute),
		SignInLimitPerMinute: core.Int64Ptr(signInLimitPerMinute),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateRateLimitConfigOptions) SetTenantID(tenantID string) *UpdateRateLimitConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetSignUpLimitPerMinute : Allow user to set SignUpLimitPerMinute
func (options *UpdateRateLimitConfigOptions) SetSignUpLimitPerMinute(signUpLimitPerMinute int64) *UpdateRateLimitConfigOptions {
	options.SignUpLimitPerMinute = core.Int64Ptr(signUpLimitPerMinute)
	return options
}

// SetSignInLimitPerMinute : Allow user to set SignInLimitPerMinute
func (options *UpdateRateLimitConfigOptions) SetSignInLimitPerMinute(signInLimitPerMinute int64) *UpdateRateLimitConfigOptions {
	options.SignInLimitPerMinute = core.Int64Ptr(signInLimitPerMinute)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRateLimitConfigOptions) SetHeaders(param map[string]string) *UpdateRateLimitConfigOptions {
	options.Headers = param
	return options
}

// UpdateRedirectUrisOptions : The UpdateRedirectUris options.
type UpdateRedirectUrisOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The redirect URIs JSON object. If IBM default credentials are used, the redirect URIs are ignored.
	RedirectUrisArray *RedirectURIConfig `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateRedirectUrisOptions : Instantiate UpdateRedirectUrisOptions
func (*AppIDManagementV4) NewUpdateRedirectUrisOptions(tenantID string, redirectUrisArray *RedirectURIConfig) *UpdateRedirectUrisOptions {
	return &UpdateRedirectUrisOptions{
		TenantID:          core.StringPtr(tenantID),
		RedirectUrisArray: redirectUrisArray,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateRedirectUrisOptions) SetTenantID(tenantID string) *UpdateRedirectUrisOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetRedirectUrisArray : Allow user to set RedirectUrisArray
func (options *UpdateRedirectUrisOptions) SetRedirectUrisArray(redirectUrisArray *RedirectURIConfig) *UpdateRedirectUrisOptions {
	options.RedirectUrisArray = redirectUrisArray
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRedirectUrisOptions) SetHeaders(param map[string]string) *UpdateRedirectUrisOptions {
	options.Headers = param
	return options
}

// UpdateRoleOptions : The UpdateRole options.
type UpdateRoleOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The role identifier.
	RoleID *string `validate:"required,ne="`

	Name *string `validate:"required"`

	Access []RoleAccessItem `validate:"required"`

	Description *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateRoleOptions : Instantiate UpdateRoleOptions
func (*AppIDManagementV4) NewUpdateRoleOptions(tenantID string, roleID string, name string, access []RoleAccessItem) *UpdateRoleOptions {
	return &UpdateRoleOptions{
		TenantID: core.StringPtr(tenantID),
		RoleID:   core.StringPtr(roleID),
		Name:     core.StringPtr(name),
		Access:   access,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateRoleOptions) SetTenantID(tenantID string) *UpdateRoleOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetRoleID : Allow user to set RoleID
func (options *UpdateRoleOptions) SetRoleID(roleID string) *UpdateRoleOptions {
	options.RoleID = core.StringPtr(roleID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateRoleOptions) SetName(name string) *UpdateRoleOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetAccess : Allow user to set Access
func (options *UpdateRoleOptions) SetAccess(access []RoleAccessItem) *UpdateRoleOptions {
	options.Access = access
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateRoleOptions) SetDescription(description string) *UpdateRoleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRoleOptions) SetHeaders(param map[string]string) *UpdateRoleOptions {
	options.Headers = param
	return options
}

// UpdateSSOConfigOptions : The UpdateSSOConfig options.
type UpdateSSOConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	InactivityTimeoutSeconds *int64 `validate:"required"`

	LogoutRedirectUris []string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSSOConfigOptions : Instantiate UpdateSSOConfigOptions
func (*AppIDManagementV4) NewUpdateSSOConfigOptions(tenantID string, isActive bool, inactivityTimeoutSeconds int64, logoutRedirectUris []string) *UpdateSSOConfigOptions {
	return &UpdateSSOConfigOptions{
		TenantID:                 core.StringPtr(tenantID),
		IsActive:                 core.BoolPtr(isActive),
		InactivityTimeoutSeconds: core.Int64Ptr(inactivityTimeoutSeconds),
		LogoutRedirectUris:       logoutRedirectUris,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateSSOConfigOptions) SetTenantID(tenantID string) *UpdateSSOConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateSSOConfigOptions) SetIsActive(isActive bool) *UpdateSSOConfigOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetInactivityTimeoutSeconds : Allow user to set InactivityTimeoutSeconds
func (options *UpdateSSOConfigOptions) SetInactivityTimeoutSeconds(inactivityTimeoutSeconds int64) *UpdateSSOConfigOptions {
	options.InactivityTimeoutSeconds = core.Int64Ptr(inactivityTimeoutSeconds)
	return options
}

// SetLogoutRedirectUris : Allow user to set LogoutRedirectUris
func (options *UpdateSSOConfigOptions) SetLogoutRedirectUris(logoutRedirectUris []string) *UpdateSSOConfigOptions {
	options.LogoutRedirectUris = logoutRedirectUris
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSSOConfigOptions) SetHeaders(param map[string]string) *UpdateSSOConfigOptions {
	options.Headers = param
	return options
}

// UpdateTemplateOptions : The UpdateTemplate options.
type UpdateTemplateOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
	// "MFA_VERIFICATION".
	TemplateName *string `validate:"required,ne="`

	// Preferred language for resource. Format as described at RFC5646. According to the configured languages codes
	// returned from the `GET /management/v4/{tenantId}/config/ui/languages` API.
	Language *string `validate:"required,ne="`

	Subject *string `validate:"required"`

	HTMLBody *string

	Base64EncodedHTMLBody *string

	PlainTextBody *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateTemplateOptions.TemplateName property.
// The type of email template. This can be "USER_VERIFICATION", "WELCOME", "PASSWORD_CHANGED", "RESET_PASSWORD" or
// "MFA_VERIFICATION".
const (
	UpdateTemplateOptionsTemplateNameMFAVerificationConst  = "MFA_VERIFICATION"
	UpdateTemplateOptionsTemplateNamePasswordChangedConst  = "PASSWORD_CHANGED"
	UpdateTemplateOptionsTemplateNameResetPasswordConst    = "RESET_PASSWORD"
	UpdateTemplateOptionsTemplateNameUserVerificationConst = "USER_VERIFICATION"
	UpdateTemplateOptionsTemplateNameWelcomeConst          = "WELCOME"
)

// NewUpdateTemplateOptions : Instantiate UpdateTemplateOptions
func (*AppIDManagementV4) NewUpdateTemplateOptions(tenantID string, templateName string, language string, subject string) *UpdateTemplateOptions {
	return &UpdateTemplateOptions{
		TenantID:     core.StringPtr(tenantID),
		TemplateName: core.StringPtr(templateName),
		Language:     core.StringPtr(language),
		Subject:      core.StringPtr(subject),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateTemplateOptions) SetTenantID(tenantID string) *UpdateTemplateOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetTemplateName : Allow user to set TemplateName
func (options *UpdateTemplateOptions) SetTemplateName(templateName string) *UpdateTemplateOptions {
	options.TemplateName = core.StringPtr(templateName)
	return options
}

// SetLanguage : Allow user to set Language
func (options *UpdateTemplateOptions) SetLanguage(language string) *UpdateTemplateOptions {
	options.Language = core.StringPtr(language)
	return options
}

// SetSubject : Allow user to set Subject
func (options *UpdateTemplateOptions) SetSubject(subject string) *UpdateTemplateOptions {
	options.Subject = core.StringPtr(subject)
	return options
}

// SetHTMLBody : Allow user to set HTMLBody
func (options *UpdateTemplateOptions) SetHTMLBody(htmlBody string) *UpdateTemplateOptions {
	options.HTMLBody = core.StringPtr(htmlBody)
	return options
}

// SetBase64EncodedHTMLBody : Allow user to set Base64EncodedHTMLBody
func (options *UpdateTemplateOptions) SetBase64EncodedHTMLBody(base64EncodedHTMLBody string) *UpdateTemplateOptions {
	options.Base64EncodedHTMLBody = core.StringPtr(base64EncodedHTMLBody)
	return options
}

// SetPlainTextBody : Allow user to set PlainTextBody
func (options *UpdateTemplateOptions) SetPlainTextBody(plainTextBody string) *UpdateTemplateOptions {
	options.PlainTextBody = core.StringPtr(plainTextBody)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTemplateOptions) SetHeaders(param map[string]string) *UpdateTemplateOptions {
	options.Headers = param
	return options
}

// UpdateUserProfilesConfigOptions : The UpdateUserProfilesConfig options.
type UpdateUserProfilesConfigOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IsActive *bool `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateUserProfilesConfigOptions : Instantiate UpdateUserProfilesConfigOptions
func (*AppIDManagementV4) NewUpdateUserProfilesConfigOptions(tenantID string, isActive bool) *UpdateUserProfilesConfigOptions {
	return &UpdateUserProfilesConfigOptions{
		TenantID: core.StringPtr(tenantID),
		IsActive: core.BoolPtr(isActive),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateUserProfilesConfigOptions) SetTenantID(tenantID string) *UpdateUserProfilesConfigOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIsActive : Allow user to set IsActive
func (options *UpdateUserProfilesConfigOptions) SetIsActive(isActive bool) *UpdateUserProfilesConfigOptions {
	options.IsActive = core.BoolPtr(isActive)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateUserProfilesConfigOptions) SetHeaders(param map[string]string) *UpdateUserProfilesConfigOptions {
	options.Headers = param
	return options
}

// UpdateUserRolesOptions : The UpdateUserRoles options.
type UpdateUserRolesOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	Roles *UpdateUserRolesParamsRoles `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateUserRolesOptions : Instantiate UpdateUserRolesOptions
func (*AppIDManagementV4) NewUpdateUserRolesOptions(tenantID string, id string, roles *UpdateUserRolesParamsRoles) *UpdateUserRolesOptions {
	return &UpdateUserRolesOptions{
		TenantID: core.StringPtr(tenantID),
		ID:       core.StringPtr(id),
		Roles:    roles,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UpdateUserRolesOptions) SetTenantID(tenantID string) *UpdateUserRolesOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *UpdateUserRolesOptions) SetID(id string) *UpdateUserRolesOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetRoles : Allow user to set Roles
func (options *UpdateUserRolesOptions) SetRoles(roles *UpdateUserRolesParamsRoles) *UpdateUserRolesOptions {
	options.Roles = roles
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateUserRolesOptions) SetHeaders(param map[string]string) *UpdateUserRolesOptions {
	options.Headers = param
	return options
}

// UpdateUserRolesParamsRoles : UpdateUserRolesParamsRoles struct
type UpdateUserRolesParamsRoles struct {
	Ids []string `json:"ids,omitempty"`
}

// UnmarshalUpdateUserRolesParamsRoles unmarshals an instance of UpdateUserRolesParamsRoles from the specified map of raw messages.
func UnmarshalUpdateUserRolesParamsRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateUserRolesParamsRoles)
	err = core.UnmarshalPrimitive(m, "ids", &obj.Ids)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserProfilesExportOptions : The UserProfilesExport options.
type UserProfilesExportOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The first result in a set list of results.
	StartIndex *int64

	// The maximum number of results per page.
	Count *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUserProfilesExportOptions : Instantiate UserProfilesExportOptions
func (*AppIDManagementV4) NewUserProfilesExportOptions(tenantID string) *UserProfilesExportOptions {
	return &UserProfilesExportOptions{
		TenantID: core.StringPtr(tenantID),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UserProfilesExportOptions) SetTenantID(tenantID string) *UserProfilesExportOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetStartIndex : Allow user to set StartIndex
func (options *UserProfilesExportOptions) SetStartIndex(startIndex int64) *UserProfilesExportOptions {
	options.StartIndex = core.Int64Ptr(startIndex)
	return options
}

// SetCount : Allow user to set Count
func (options *UserProfilesExportOptions) SetCount(count int64) *UserProfilesExportOptions {
	options.Count = core.Int64Ptr(count)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UserProfilesExportOptions) SetHeaders(param map[string]string) *UserProfilesExportOptions {
	options.Headers = param
	return options
}

// UserProfilesImportOptions : The UserProfilesImport options.
type UserProfilesImportOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	Users []ExportUserProfileUsersItem `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUserProfilesImportOptions : Instantiate UserProfilesImportOptions
func (*AppIDManagementV4) NewUserProfilesImportOptions(tenantID string, users []ExportUserProfileUsersItem) *UserProfilesImportOptions {
	return &UserProfilesImportOptions{
		TenantID: core.StringPtr(tenantID),
		Users:    users,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UserProfilesImportOptions) SetTenantID(tenantID string) *UserProfilesImportOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetUsers : Allow user to set Users
func (options *UserProfilesImportOptions) SetUsers(users []ExportUserProfileUsersItem) *UserProfilesImportOptions {
	options.Users = users
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UserProfilesImportOptions) SetHeaders(param map[string]string) *UserProfilesImportOptions {
	options.Headers = param
	return options
}

// UserSearchResponseRequestOptions : UserSearchResponseRequestOptions struct
type UserSearchResponseRequestOptions struct {
	StartIndex *int64 `json:"startIndex,omitempty"`

	Count *int64 `json:"count,omitempty"`
}

// UnmarshalUserSearchResponseRequestOptions unmarshals an instance of UserSearchResponseRequestOptions from the specified map of raw messages.
func UnmarshalUserSearchResponseRequestOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserSearchResponseRequestOptions)
	err = core.UnmarshalPrimitive(m, "startIndex", &obj.StartIndex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserSearchResponseUsersItem : UserSearchResponseUsersItem struct
type UserSearchResponseUsersItem struct {
	ID *string `json:"id,omitempty"`

	IDP *string `json:"idp,omitempty"`

	Email *string `json:"email,omitempty"`
}

// UnmarshalUserSearchResponseUsersItem unmarshals an instance of UserSearchResponseUsersItem from the specified map of raw messages.
func UnmarshalUserSearchResponseUsersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserSearchResponseUsersItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "idp", &obj.IDP)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserVerificationResultOptions : The UserVerificationResult options.
type UserVerificationResultOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The context that will be use to get the verification result.
	Context *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUserVerificationResultOptions : Instantiate UserVerificationResultOptions
func (*AppIDManagementV4) NewUserVerificationResultOptions(tenantID string, context string) *UserVerificationResultOptions {
	return &UserVerificationResultOptions{
		TenantID: core.StringPtr(tenantID),
		Context:  core.StringPtr(context),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UserVerificationResultOptions) SetTenantID(tenantID string) *UserVerificationResultOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetContext : Allow user to set Context
func (options *UserVerificationResultOptions) SetContext(context string) *UserVerificationResultOptions {
	options.Context = core.StringPtr(context)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UserVerificationResultOptions) SetHeaders(param map[string]string) *UserVerificationResultOptions {
	options.Headers = param
	return options
}

// UsersDeleteUserProfileOptions : The UsersDeleteUserProfile options.
type UsersDeleteUserProfileOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUsersDeleteUserProfileOptions : Instantiate UsersDeleteUserProfileOptions
func (*AppIDManagementV4) NewUsersDeleteUserProfileOptions(tenantID string, id string) *UsersDeleteUserProfileOptions {
	return &UsersDeleteUserProfileOptions{
		TenantID: core.StringPtr(tenantID),
		ID:       core.StringPtr(id),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersDeleteUserProfileOptions) SetTenantID(tenantID string) *UsersDeleteUserProfileOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *UsersDeleteUserProfileOptions) SetID(id string) *UsersDeleteUserProfileOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersDeleteUserProfileOptions) SetHeaders(param map[string]string) *UsersDeleteUserProfileOptions {
	options.Headers = param
	return options
}

// UsersGetUserProfileOptions : The UsersGetUserProfile options.
type UsersGetUserProfileOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUsersGetUserProfileOptions : Instantiate UsersGetUserProfileOptions
func (*AppIDManagementV4) NewUsersGetUserProfileOptions(tenantID string, id string) *UsersGetUserProfileOptions {
	return &UsersGetUserProfileOptions{
		TenantID: core.StringPtr(tenantID),
		ID:       core.StringPtr(id),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersGetUserProfileOptions) SetTenantID(tenantID string) *UsersGetUserProfileOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *UsersGetUserProfileOptions) SetID(id string) *UsersGetUserProfileOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersGetUserProfileOptions) SetHeaders(param map[string]string) *UsersGetUserProfileOptions {
	options.Headers = param
	return options
}

// UsersList : UsersList struct
type UsersList struct {
	TotalResults *int64 `json:"totalResults,omitempty"`

	ItemsPerPage *int64 `json:"itemsPerPage,omitempty"`

	Resources []interface{} `json:"Resources" validate:"required"`
}

// UnmarshalUsersList unmarshals an instance of UsersList from the specified map of raw messages.
func UnmarshalUsersList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsersList)
	err = core.UnmarshalPrimitive(m, "totalResults", &obj.TotalResults)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "itemsPerPage", &obj.ItemsPerPage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "Resources", &obj.Resources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsersNominateUserOptions : The UsersNominateUser options.
type UsersNominateUserOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	IDP *string `validate:"required"`

	IDPIdentity *string `validate:"required"`

	Profile *UsersNominateUserParamsProfile

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UsersNominateUserOptions.IDP property.
const (
	UsersNominateUserOptionsIDPAppidCustomConst    = "appid_custom"
	UsersNominateUserOptionsIDPCloudDirectoryConst = "cloud_directory"
	UsersNominateUserOptionsIDPFacebookConst       = "facebook"
	UsersNominateUserOptionsIDPGoogleConst         = "google"
	UsersNominateUserOptionsIDPIbmidConst          = "ibmid"
	UsersNominateUserOptionsIDPSAMLConst           = "saml"
)

// NewUsersNominateUserOptions : Instantiate UsersNominateUserOptions
func (*AppIDManagementV4) NewUsersNominateUserOptions(tenantID string, idp string, idpIdentity string) *UsersNominateUserOptions {
	return &UsersNominateUserOptions{
		TenantID:    core.StringPtr(tenantID),
		IDP:         core.StringPtr(idp),
		IDPIdentity: core.StringPtr(idpIdentity),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersNominateUserOptions) SetTenantID(tenantID string) *UsersNominateUserOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetIDP : Allow user to set IDP
func (options *UsersNominateUserOptions) SetIDP(idp string) *UsersNominateUserOptions {
	options.IDP = core.StringPtr(idp)
	return options
}

// SetIDPIdentity : Allow user to set IDPIdentity
func (options *UsersNominateUserOptions) SetIDPIdentity(idpIdentity string) *UsersNominateUserOptions {
	options.IDPIdentity = core.StringPtr(idpIdentity)
	return options
}

// SetProfile : Allow user to set Profile
func (options *UsersNominateUserOptions) SetProfile(profile *UsersNominateUserParamsProfile) *UsersNominateUserOptions {
	options.Profile = profile
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersNominateUserOptions) SetHeaders(param map[string]string) *UsersNominateUserOptions {
	options.Headers = param
	return options
}

// UsersNominateUserParamsProfile : UsersNominateUserParamsProfile struct
type UsersNominateUserParamsProfile struct {
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// UnmarshalUsersNominateUserParamsProfile unmarshals an instance of UsersNominateUserParamsProfile from the specified map of raw messages.
func UnmarshalUsersNominateUserParamsProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UsersNominateUserParamsProfile)
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UsersRevokeRefreshTokenOptions : The UsersRevokeRefreshToken options.
type UsersRevokeRefreshTokenOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUsersRevokeRefreshTokenOptions : Instantiate UsersRevokeRefreshTokenOptions
func (*AppIDManagementV4) NewUsersRevokeRefreshTokenOptions(tenantID string, id string) *UsersRevokeRefreshTokenOptions {
	return &UsersRevokeRefreshTokenOptions{
		TenantID: core.StringPtr(tenantID),
		ID:       core.StringPtr(id),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersRevokeRefreshTokenOptions) SetTenantID(tenantID string) *UsersRevokeRefreshTokenOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *UsersRevokeRefreshTokenOptions) SetID(id string) *UsersRevokeRefreshTokenOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersRevokeRefreshTokenOptions) SetHeaders(param map[string]string) *UsersRevokeRefreshTokenOptions {
	options.Headers = param
	return options
}

// UsersSearchUserProfileOptions : The UsersSearchUserProfile options.
type UsersSearchUserProfileOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// display user data.
	DataScope *string `validate:"required"`

	// Email (as retrieved from the Identity Provider).
	Email *string

	// The IDP specific user identifier.
	ID *string

	// The first result in a set list of results.
	StartIndex *int64

	// The maximum number of results per page.
	Count *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UsersSearchUserProfileOptions.DataScope property.
// display user data.
const (
	UsersSearchUserProfileOptionsDataScopeFullConst  = "full"
	UsersSearchUserProfileOptionsDataScopeIndexConst = "index"
)

// NewUsersSearchUserProfileOptions : Instantiate UsersSearchUserProfileOptions
func (*AppIDManagementV4) NewUsersSearchUserProfileOptions(tenantID string, dataScope string) *UsersSearchUserProfileOptions {
	return &UsersSearchUserProfileOptions{
		TenantID:  core.StringPtr(tenantID),
		DataScope: core.StringPtr(dataScope),
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersSearchUserProfileOptions) SetTenantID(tenantID string) *UsersSearchUserProfileOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetDataScope : Allow user to set DataScope
func (options *UsersSearchUserProfileOptions) SetDataScope(dataScope string) *UsersSearchUserProfileOptions {
	options.DataScope = core.StringPtr(dataScope)
	return options
}

// SetEmail : Allow user to set Email
func (options *UsersSearchUserProfileOptions) SetEmail(email string) *UsersSearchUserProfileOptions {
	options.Email = core.StringPtr(email)
	return options
}

// SetID : Allow user to set ID
func (options *UsersSearchUserProfileOptions) SetID(id string) *UsersSearchUserProfileOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetStartIndex : Allow user to set StartIndex
func (options *UsersSearchUserProfileOptions) SetStartIndex(startIndex int64) *UsersSearchUserProfileOptions {
	options.StartIndex = core.Int64Ptr(startIndex)
	return options
}

// SetCount : Allow user to set Count
func (options *UsersSearchUserProfileOptions) SetCount(count int64) *UsersSearchUserProfileOptions {
	options.Count = core.Int64Ptr(count)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersSearchUserProfileOptions) SetHeaders(param map[string]string) *UsersSearchUserProfileOptions {
	options.Headers = param
	return options
}

// UsersSetUserProfileOptions : The UsersSetUserProfile options.
type UsersSetUserProfileOptions struct {
	// The service tenantId. The tenantId can be found in the service credentials.
	TenantID *string `validate:"required,ne="`

	// The user's identifier ('subject' in identity token) You can search user in <a
	// href="#!/Users/users_search_user_profile" target="_blank">/users</a>.
	ID *string `validate:"required,ne="`

	Attributes map[string]interface{} `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUsersSetUserProfileOptions : Instantiate UsersSetUserProfileOptions
func (*AppIDManagementV4) NewUsersSetUserProfileOptions(tenantID string, id string, attributes map[string]interface{}) *UsersSetUserProfileOptions {
	return &UsersSetUserProfileOptions{
		TenantID:   core.StringPtr(tenantID),
		ID:         core.StringPtr(id),
		Attributes: attributes,
	}
}

// SetTenantID : Allow user to set TenantID
func (options *UsersSetUserProfileOptions) SetTenantID(tenantID string) *UsersSetUserProfileOptions {
	options.TenantID = core.StringPtr(tenantID)
	return options
}

// SetID : Allow user to set ID
func (options *UsersSetUserProfileOptions) SetID(id string) *UsersSetUserProfileOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetAttributes : Allow user to set Attributes
func (options *UsersSetUserProfileOptions) SetAttributes(attributes map[string]interface{}) *UsersSetUserProfileOptions {
	options.Attributes = attributes
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UsersSetUserProfileOptions) SetHeaders(param map[string]string) *UsersSetUserProfileOptions {
	options.Headers = param
	return options
}

// AccessTokenConfigParams : AccessTokenConfigParams struct
type AccessTokenConfigParams struct {
	ExpiresIn *int64 `json:"expires_in" validate:"required"`
}

// NewAccessTokenConfigParams : Instantiate AccessTokenConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewAccessTokenConfigParams(expiresIn int64) (model *AccessTokenConfigParams, err error) {
	model = &AccessTokenConfigParams{
		ExpiresIn: core.Int64Ptr(expiresIn),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalAccessTokenConfigParams unmarshals an instance of AccessTokenConfigParams from the specified map of raw messages.
func UnmarshalAccessTokenConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessTokenConfigParams)
	err = core.UnmarshalPrimitive(m, "expires_in", &obj.ExpiresIn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionURLResponse : ActionURLResponse struct
type ActionURLResponse struct {
	ActionURL *string `json:"actionUrl" validate:"required"`
}

// UnmarshalActionURLResponse unmarshals an instance of ActionURLResponse from the specified map of raw messages.
func UnmarshalActionURLResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionURLResponse)
	err = core.UnmarshalPrimitive(m, "actionUrl", &obj.ActionURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApmSchema : ApmSchema struct
type ApmSchema struct {
	AdvancedPasswordManagement *ApmSchemaAdvancedPasswordManagement `json:"advancedPasswordManagement" validate:"required"`
}

// NewApmSchema : Instantiate ApmSchema (Generic Model Constructor)
func (*AppIDManagementV4) NewApmSchema(advancedPasswordManagement *ApmSchemaAdvancedPasswordManagement) (model *ApmSchema, err error) {
	model = &ApmSchema{
		AdvancedPasswordManagement: advancedPasswordManagement,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalApmSchema unmarshals an instance of ApmSchema from the specified map of raw messages.
func UnmarshalApmSchema(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApmSchema)
	err = core.UnmarshalModel(m, "advancedPasswordManagement", &obj.AdvancedPasswordManagement, UnmarshalApmSchemaAdvancedPasswordManagement)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssignRoleToUser : AssignRoleToUser struct
type AssignRoleToUser struct {
	Roles []AssignRoleToUserRolesItem `json:"roles,omitempty"`
}

// UnmarshalAssignRoleToUser unmarshals an instance of AssignRoleToUser from the specified map of raw messages.
func UnmarshalAssignRoleToUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssignRoleToUser)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalAssignRoleToUserRolesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectoryConfigParams : CloudDirectoryConfigParams struct
type CloudDirectoryConfigParams struct {
	SelfServiceEnabled *bool `json:"selfServiceEnabled" validate:"required"`

	SignupEnabled *bool `json:"signupEnabled,omitempty"`

	Interactions *CloudDirectoryConfigParamsInteractions `json:"interactions" validate:"required"`

	IdentityField *string `json:"identityField,omitempty"`
}

// Constants associated with the CloudDirectoryConfigParams.IdentityField property.
const (
	CloudDirectoryConfigParamsIdentityFieldEmailConst    = "email"
	CloudDirectoryConfigParamsIdentityFieldUsernameConst = "userName"
)

// NewCloudDirectoryConfigParams : Instantiate CloudDirectoryConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectoryConfigParams(selfServiceEnabled bool, interactions *CloudDirectoryConfigParamsInteractions) (model *CloudDirectoryConfigParams, err error) {
	model = &CloudDirectoryConfigParams{
		SelfServiceEnabled: core.BoolPtr(selfServiceEnabled),
		Interactions:       interactions,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectoryConfigParams unmarshals an instance of CloudDirectoryConfigParams from the specified map of raw messages.
func UnmarshalCloudDirectoryConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectoryConfigParams)
	err = core.UnmarshalPrimitive(m, "selfServiceEnabled", &obj.SelfServiceEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signupEnabled", &obj.SignupEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "interactions", &obj.Interactions, UnmarshalCloudDirectoryConfigParamsInteractions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "identityField", &obj.IdentityField)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectoryResponse : CloudDirectoryResponse struct
type CloudDirectoryResponse struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *CloudDirectoryConfigParams `json:"config,omitempty"`
}

// UnmarshalCloudDirectoryResponse unmarshals an instance of CloudDirectoryResponse from the specified map of raw messages.
func UnmarshalCloudDirectoryResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectoryResponse)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalCloudDirectoryConfigParams)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudDirectorySenderDetails : CloudDirectorySenderDetails struct
type CloudDirectorySenderDetails struct {
	SenderDetails *CloudDirectorySenderDetailsSenderDetails `json:"senderDetails" validate:"required"`
}

// NewCloudDirectorySenderDetails : Instantiate CloudDirectorySenderDetails (Generic Model Constructor)
func (*AppIDManagementV4) NewCloudDirectorySenderDetails(senderDetails *CloudDirectorySenderDetailsSenderDetails) (model *CloudDirectorySenderDetails, err error) {
	model = &CloudDirectorySenderDetails{
		SenderDetails: senderDetails,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCloudDirectorySenderDetails unmarshals an instance of CloudDirectorySenderDetails from the specified map of raw messages.
func UnmarshalCloudDirectorySenderDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudDirectorySenderDetails)
	err = core.UnmarshalModel(m, "senderDetails", &obj.SenderDetails, UnmarshalCloudDirectorySenderDetailsSenderDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfirmationResultOk : ConfirmationResultOk struct
type ConfirmationResultOk struct {
	Success *bool `json:"success" validate:"required"`

	UUID *string `json:"uuid" validate:"required"`
}

// UnmarshalConfirmationResultOk unmarshals an instance of ConfirmationResultOk from the specified map of raw messages.
func UnmarshalConfirmationResultOk(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfirmationResultOk)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uuid", &obj.UUID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRoleResponse : CreateRoleResponse struct
type CreateRoleResponse struct {
	ID *string `json:"id" validate:"required"`

	Name *string `json:"name" validate:"required"`

	Description *string `json:"description,omitempty"`

	Access []RoleAccessItem `json:"access" validate:"required"`
}

// UnmarshalCreateRoleResponse unmarshals an instance of CreateRoleResponse from the specified map of raw messages.
func UnmarshalCreateRoleResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateRoleResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "access", &obj.Access, UnmarshalRoleAccessItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CustomIDPConfigParams : CustomIDPConfigParams struct
type CustomIDPConfigParams struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *CustomIDPConfigParamsConfig `json:"config,omitempty"`
}

// NewCustomIDPConfigParams : Instantiate CustomIDPConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewCustomIDPConfigParams(isActive bool) (model *CustomIDPConfigParams, err error) {
	model = &CustomIDPConfigParams{
		IsActive: core.BoolPtr(isActive),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCustomIDPConfigParams unmarshals an instance of CustomIDPConfigParams from the specified map of raw messages.
func UnmarshalCustomIDPConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomIDPConfigParams)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalCustomIDPConfigParamsConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EmailDispatcherParams : EmailDispatcherParams struct
type EmailDispatcherParams struct {
	Provider *string `json:"provider" validate:"required"`

	Sendgrid *EmailDispatcherParamsSendgrid `json:"sendgrid,omitempty"`

	Custom *EmailDispatcherParamsCustom `json:"custom,omitempty"`
}

// Constants associated with the EmailDispatcherParams.Provider property.
const (
	EmailDispatcherParamsProviderAppidConst    = "appid"
	EmailDispatcherParamsProviderCustomConst   = "custom"
	EmailDispatcherParamsProviderSendgridConst = "sendgrid"
)

// NewEmailDispatcherParams : Instantiate EmailDispatcherParams (Generic Model Constructor)
func (*AppIDManagementV4) NewEmailDispatcherParams(provider string) (model *EmailDispatcherParams, err error) {
	model = &EmailDispatcherParams{
		Provider: core.StringPtr(provider),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEmailDispatcherParams unmarshals an instance of EmailDispatcherParams from the specified map of raw messages.
func UnmarshalEmailDispatcherParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EmailDispatcherParams)
	err = core.UnmarshalPrimitive(m, "provider", &obj.Provider)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sendgrid", &obj.Sendgrid, UnmarshalEmailDispatcherParamsSendgrid)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "custom", &obj.Custom, UnmarshalEmailDispatcherParamsCustom)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUser : ExportUser struct
type ExportUser struct {
	Users []ExportUserUsersItem `json:"users" validate:"required"`
}

// NewExportUser : Instantiate ExportUser (Generic Model Constructor)
func (*AppIDManagementV4) NewExportUser(users []ExportUserUsersItem) (model *ExportUser, err error) {
	model = &ExportUser{
		Users: users,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExportUser unmarshals an instance of ExportUser from the specified map of raw messages.
func UnmarshalExportUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUser)
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalExportUserUsersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportUserProfile : ExportUserProfile struct
type ExportUserProfile struct {
	Users []ExportUserProfileUsersItem `json:"users" validate:"required"`
}

// NewExportUserProfile : Instantiate ExportUserProfile (Generic Model Constructor)
func (*AppIDManagementV4) NewExportUserProfile(users []ExportUserProfileUsersItem) (model *ExportUserProfile, err error) {
	model = &ExportUserProfile{
		Users: users,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExportUserProfile unmarshals an instance of ExportUserProfile from the specified map of raw messages.
func UnmarshalExportUserProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportUserProfile)
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalExportUserProfileUsersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExtensionActive : ExtensionActive struct
type ExtensionActive struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config interface{} `json:"config,omitempty"`
}

// NewExtensionActive : Instantiate ExtensionActive (Generic Model Constructor)
func (*AppIDManagementV4) NewExtensionActive(isActive bool) (model *ExtensionActive, err error) {
	model = &ExtensionActive{
		IsActive: core.BoolPtr(isActive),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalExtensionActive unmarshals an instance of ExtensionActive from the specified map of raw messages.
func UnmarshalExtensionActive(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExtensionActive)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "config", &obj.Config)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExtensionTest : ExtensionTest struct
type ExtensionTest struct {
	StatusCode *int64 `json:"statusCode,omitempty"`

	HeadersVar interface{} `json:"headers,omitempty"`
}

// UnmarshalExtensionTest unmarshals an instance of ExtensionTest from the specified map of raw messages.
func UnmarshalExtensionTest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExtensionTest)
	err = core.UnmarshalPrimitive(m, "statusCode", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "headers", &obj.HeadersVar)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookConfigParams : FacebookConfigParams struct
type FacebookConfigParams struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *FacebookConfigParamsConfig `json:"config,omitempty"`

	RedirectURL *string `json:"redirectURL,omitempty"`
}

// UnmarshalFacebookConfigParams unmarshals an instance of FacebookConfigParams from the specified map of raw messages.
func UnmarshalFacebookConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookConfigParams)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalFacebookConfigParamsConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "redirectURL", &obj.RedirectURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookConfigParamsPut : FacebookConfigParamsPut struct
type FacebookConfigParamsPut struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *FacebookConfigParamsPutConfig `json:"config,omitempty"`

	RedirectURL *string `json:"redirectURL,omitempty"`
}

// UnmarshalFacebookConfigParamsPut unmarshals an instance of FacebookConfigParamsPut from the specified map of raw messages.
func UnmarshalFacebookConfigParamsPut(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookConfigParamsPut)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalFacebookConfigParamsPutConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "redirectURL", &obj.RedirectURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FacebookGoogleConfigParams : FacebookGoogleConfigParams struct
type FacebookGoogleConfigParams struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *FacebookGoogleConfigParamsConfig `json:"config,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// NewFacebookGoogleConfigParams : Instantiate FacebookGoogleConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewFacebookGoogleConfigParams(isActive bool) (model *FacebookGoogleConfigParams, err error) {
	model = &FacebookGoogleConfigParams{
		IsActive: core.BoolPtr(isActive),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// SetProperty allows the user to set an arbitrary property on an instance of FacebookGoogleConfigParams
func (o *FacebookGoogleConfigParams) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of FacebookGoogleConfigParams
func (o *FacebookGoogleConfigParams) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of FacebookGoogleConfigParams
func (o *FacebookGoogleConfigParams) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of FacebookGoogleConfigParams
func (o *FacebookGoogleConfigParams) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.IsActive != nil {
		m["isActive"] = o.IsActive
	}
	if o.Config != nil {
		m["config"] = o.Config
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalFacebookGoogleConfigParams unmarshals an instance of FacebookGoogleConfigParams from the specified map of raw messages.
func UnmarshalFacebookGoogleConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FacebookGoogleConfigParams)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	delete(m, "isActive")
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalFacebookGoogleConfigParamsConfig)
	if err != nil {
		return
	}
	delete(m, "config")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetLanguages : GetLanguages struct
type GetLanguages struct {
	Languages []string `json:"languages" validate:"required"`
}

// NewGetLanguages : Instantiate GetLanguages (Generic Model Constructor)
func (*AppIDManagementV4) NewGetLanguages(languages []string) (model *GetLanguages, err error) {
	model = &GetLanguages{
		Languages: languages,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalGetLanguages unmarshals an instance of GetLanguages from the specified map of raw messages.
func UnmarshalGetLanguages(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetLanguages)
	err = core.UnmarshalPrimitive(m, "languages", &obj.Languages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetMFAConfiguration : GetMFAConfiguration struct
type GetMFAConfiguration struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

// UnmarshalGetMFAConfiguration unmarshals an instance of GetMFAConfiguration from the specified map of raw messages.
func UnmarshalGetMFAConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetMFAConfiguration)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetRoleResponse : GetRoleResponse struct
type GetRoleResponse struct {
	ID *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`

	Description *string `json:"description,omitempty"`

	Access []RoleAccessItem `json:"access,omitempty"`
}

// UnmarshalGetRoleResponse unmarshals an instance of GetRoleResponse from the specified map of raw messages.
func UnmarshalGetRoleResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetRoleResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "access", &obj.Access, UnmarshalRoleAccessItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetSMSChannel : GetSMSChannel struct
type GetSMSChannel struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Type *string `json:"type" validate:"required"`

	Config *GetSMSChannelConfig `json:"config,omitempty"`
}

// UnmarshalGetSMSChannel unmarshals an instance of GetSMSChannel from the specified map of raw messages.
func UnmarshalGetSMSChannel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetSMSChannel)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalGetSMSChannelConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetScopesForApplication : GetScopesForApplication struct
type GetScopesForApplication struct {
	Scopes []string `json:"scopes,omitempty"`
}

// UnmarshalGetScopesForApplication unmarshals an instance of GetScopesForApplication from the specified map of raw messages.
func UnmarshalGetScopesForApplication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetScopesForApplication)
	err = core.UnmarshalPrimitive(m, "scopes", &obj.Scopes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetTemplate : GetTemplate struct
type GetTemplate struct {
	Subject *string `json:"subject" validate:"required"`

	HTMLBody *string `json:"html_body,omitempty"`

	Base64EncodedHTMLBody *string `json:"base64_encoded_html_body,omitempty"`

	PlainTextBody *string `json:"plain_text_body,omitempty"`
}

// UnmarshalGetTemplate unmarshals an instance of GetTemplate from the specified map of raw messages.
func UnmarshalGetTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetTemplate)
	err = core.UnmarshalPrimitive(m, "subject", &obj.Subject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "html_body", &obj.HTMLBody)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "base64_encoded_html_body", &obj.Base64EncodedHTMLBody)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "plain_text_body", &obj.PlainTextBody)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUser : GetUser struct
type GetUser struct {
	DisplayName *string `json:"displayName,omitempty"`

	Active *bool `json:"active" validate:"required"`

	LockedUntil *int64 `json:"lockedUntil,omitempty"`

	Emails []GetUserEmailsItem `json:"emails" validate:"required"`

	Meta *GetUserMeta `json:"meta" validate:"required"`

	Schemas []string `json:"schemas,omitempty"`

	Name *GetUserName `json:"name,omitempty"`

	UserName *string `json:"userName,omitempty"`

	ID *string `json:"id" validate:"required"`

	Status *string `json:"status" validate:"required"`
}

// UnmarshalGetUser unmarshals an instance of GetUser from the specified map of raw messages.
func UnmarshalGetUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUser)
	err = core.UnmarshalPrimitive(m, "displayName", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lockedUntil", &obj.LockedUntil)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "emails", &obj.Emails, UnmarshalGetUserEmailsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "meta", &obj.Meta, UnmarshalGetUserMeta)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schemas", &obj.Schemas)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "name", &obj.Name, UnmarshalGetUserName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "userName", &obj.UserName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserAndProfile : GetUserAndProfile struct
type GetUserAndProfile struct {
	Sub *string `json:"sub,omitempty"`

	Identities []GetUserAndProfileIdentitiesItem `json:"identities,omitempty"`

	Attributes interface{} `json:"attributes,omitempty"`
}

// UnmarshalGetUserAndProfile unmarshals an instance of GetUserAndProfile from the specified map of raw messages.
func UnmarshalGetUserAndProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserAndProfile)
	err = core.UnmarshalPrimitive(m, "sub", &obj.Sub)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "identities", &obj.Identities, UnmarshalGetUserAndProfileIdentitiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetUserRolesResponse : GetUserRolesResponse struct
type GetUserRolesResponse struct {
	Roles []GetUserRolesResponseRolesItem `json:"roles,omitempty"`
}

// UnmarshalGetUserRolesResponse unmarshals an instance of GetUserRolesResponse from the specified map of raw messages.
func UnmarshalGetUserRolesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetUserRolesResponse)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalGetUserRolesResponseRolesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoogleConfigParams : GoogleConfigParams struct
type GoogleConfigParams struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *GoogleConfigParamsConfig `json:"config,omitempty"`

	RedirectURL *string `json:"redirectURL,omitempty"`
}

// UnmarshalGoogleConfigParams unmarshals an instance of GoogleConfigParams from the specified map of raw messages.
func UnmarshalGoogleConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoogleConfigParams)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalGoogleConfigParamsConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "redirectURL", &obj.RedirectURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GoogleConfigParamsPut : GoogleConfigParamsPut struct
type GoogleConfigParamsPut struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *GoogleConfigParamsPutConfig `json:"config,omitempty"`

	RedirectURL *string `json:"redirectURL,omitempty"`
}

// UnmarshalGoogleConfigParamsPut unmarshals an instance of GoogleConfigParamsPut from the specified map of raw messages.
func UnmarshalGoogleConfigParamsPut(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GoogleConfigParamsPut)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalGoogleConfigParamsPutConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "redirectURL", &obj.RedirectURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportProfilesResponse : ImportProfilesResponse struct
type ImportProfilesResponse struct {
	Added *int64 `json:"added,omitempty"`

	Failed *int64 `json:"failed,omitempty"`

	FailReasons []ImportProfilesResponseFailReasonsItem `json:"failReasons,omitempty"`
}

// UnmarshalImportProfilesResponse unmarshals an instance of ImportProfilesResponse from the specified map of raw messages.
func UnmarshalImportProfilesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportProfilesResponse)
	err = core.UnmarshalPrimitive(m, "added", &obj.Added)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "failReasons", &obj.FailReasons, UnmarshalImportProfilesResponseFailReasonsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportResponse : ImportResponse struct
type ImportResponse struct {
	Added *int64 `json:"added,omitempty"`

	Failed *int64 `json:"failed,omitempty"`

	FailReasons []ImportResponseFailReasonsItem `json:"failReasons,omitempty"`
}

// UnmarshalImportResponse unmarshals an instance of ImportResponse from the specified map of raw messages.
func UnmarshalImportResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportResponse)
	err = core.UnmarshalPrimitive(m, "added", &obj.Added)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "failReasons", &obj.FailReasons, UnmarshalImportResponseFailReasonsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PasswordRegexConfigParamsGet : PasswordRegexConfigParamsGet struct
type PasswordRegexConfigParamsGet struct {
	Regex *string `json:"regex,omitempty"`

	Base64EncodedRegex *string `json:"base64_encoded_regex,omitempty"`

	ErrorMessage *string `json:"error_message,omitempty"`
}

// UnmarshalPasswordRegexConfigParamsGet unmarshals an instance of PasswordRegexConfigParamsGet from the specified map of raw messages.
func UnmarshalPasswordRegexConfigParamsGet(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PasswordRegexConfigParamsGet)
	err = core.UnmarshalPrimitive(m, "regex", &obj.Regex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "base64_encoded_regex", &obj.Base64EncodedRegex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_message", &obj.ErrorMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RedirectURIConfig : RedirectURIConfig struct
type RedirectURIConfig struct {
	RedirectUris []string `json:"redirectUris,omitempty"`

	TrustCloudIAMRedirectUris *bool `json:"trustCloudIAMRedirectUris,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of RedirectURIConfig
func (o *RedirectURIConfig) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of RedirectURIConfig
func (o *RedirectURIConfig) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of RedirectURIConfig
func (o *RedirectURIConfig) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of RedirectURIConfig
func (o *RedirectURIConfig) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.RedirectUris != nil {
		m["redirectUris"] = o.RedirectUris
	}
	if o.TrustCloudIAMRedirectUris != nil {
		m["trustCloudIAMRedirectUris"] = o.TrustCloudIAMRedirectUris
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalRedirectURIConfig unmarshals an instance of RedirectURIConfig from the specified map of raw messages.
func UnmarshalRedirectURIConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RedirectURIConfig)
	err = core.UnmarshalPrimitive(m, "redirectUris", &obj.RedirectUris)
	if err != nil {
		return
	}
	delete(m, "redirectUris")
	err = core.UnmarshalPrimitive(m, "trustCloudIAMRedirectUris", &obj.TrustCloudIAMRedirectUris)
	if err != nil {
		return
	}
	delete(m, "trustCloudIAMRedirectUris")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RedirectURIResponse : RedirectURIResponse struct
type RedirectURIResponse struct {
	RedirectUris []string `json:"redirectUris,omitempty"`
}

// UnmarshalRedirectURIResponse unmarshals an instance of RedirectURIResponse from the specified map of raw messages.
func UnmarshalRedirectURIResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RedirectURIResponse)
	err = core.UnmarshalPrimitive(m, "redirectUris", &obj.RedirectUris)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RespCustomEmailDisParams : RespCustomEmailDisParams struct
type RespCustomEmailDisParams struct {
	StatusCode *int64 `json:"statusCode,omitempty"`

	HeadersVar interface{} `json:"headers,omitempty"`
}

// UnmarshalRespCustomEmailDisParams unmarshals an instance of RespCustomEmailDisParams from the specified map of raw messages.
func UnmarshalRespCustomEmailDisParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RespCustomEmailDisParams)
	err = core.UnmarshalPrimitive(m, "statusCode", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "headers", &obj.HeadersVar)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RespEmailSettingsTest : RespEmailSettingsTest struct
type RespEmailSettingsTest struct {
	Success *bool `json:"success" validate:"required"`

	DispatcherStatusCode *int64 `json:"dispatcherStatusCode" validate:"required"`

	DispatcherResponse interface{} `json:"dispatcherResponse,omitempty"`
}

// UnmarshalRespEmailSettingsTest unmarshals an instance of RespEmailSettingsTest from the specified map of raw messages.
func UnmarshalRespEmailSettingsTest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RespEmailSettingsTest)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dispatcherStatusCode", &obj.DispatcherStatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dispatcherResponse", &obj.DispatcherResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RespSMSDisParams : RespSMSDisParams struct
type RespSMSDisParams struct {
	ConfirmationCode *int64 `json:"confirmationCode,omitempty"`

	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// UnmarshalRespSMSDisParams unmarshals an instance of RespSMSDisParams from the specified map of raw messages.
func UnmarshalRespSMSDisParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RespSMSDisParams)
	err = core.UnmarshalPrimitive(m, "confirmationCode", &obj.ConfirmationCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "phoneNumber", &obj.PhoneNumber)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RoleAccessItem : RoleAccessItem struct
type RoleAccessItem struct {
	ApplicationID *string `json:"application_id" validate:"required"`

	Scopes []string `json:"scopes" validate:"required"`
}

// NewRoleAccessItem : Instantiate RoleAccessItem (Generic Model Constructor)
func (*AppIDManagementV4) NewRoleAccessItem(applicationID string, scopes []string) (model *RoleAccessItem, err error) {
	model = &RoleAccessItem{
		ApplicationID: core.StringPtr(applicationID),
		Scopes:        scopes,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRoleAccessItem unmarshals an instance of RoleAccessItem from the specified map of raw messages.
func UnmarshalRoleAccessItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleAccessItem)
	err = core.UnmarshalPrimitive(m, "application_id", &obj.ApplicationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scopes", &obj.Scopes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLConfigParams : SAMLConfigParams struct
type SAMLConfigParams struct {
	EntityID *string `json:"entityID" validate:"required"`

	SignInURL *string `json:"signInUrl" validate:"required"`

	Certificates []string `json:"certificates" validate:"required"`

	DisplayName *string `json:"displayName,omitempty"`

	AuthnContext *SAMLConfigParamsAuthnContext `json:"authnContext,omitempty"`

	SignRequest *bool `json:"signRequest,omitempty"`

	EncryptResponse *bool `json:"encryptResponse,omitempty"`

	IncludeScoping *bool `json:"includeScoping,omitempty"`

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// NewSAMLConfigParams : Instantiate SAMLConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewSAMLConfigParams(entityID string, signInURL string, certificates []string) (model *SAMLConfigParams, err error) {
	model = &SAMLConfigParams{
		EntityID:     core.StringPtr(entityID),
		SignInURL:    core.StringPtr(signInURL),
		Certificates: certificates,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// SetProperty allows the user to set an arbitrary property on an instance of SAMLConfigParams
func (o *SAMLConfigParams) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SAMLConfigParams
func (o *SAMLConfigParams) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SAMLConfigParams
func (o *SAMLConfigParams) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SAMLConfigParams
func (o *SAMLConfigParams) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.EntityID != nil {
		m["entityID"] = o.EntityID
	}
	if o.SignInURL != nil {
		m["signInUrl"] = o.SignInURL
	}
	if o.Certificates != nil {
		m["certificates"] = o.Certificates
	}
	if o.DisplayName != nil {
		m["displayName"] = o.DisplayName
	}
	if o.AuthnContext != nil {
		m["authnContext"] = o.AuthnContext
	}
	if o.SignRequest != nil {
		m["signRequest"] = o.SignRequest
	}
	if o.EncryptResponse != nil {
		m["encryptResponse"] = o.EncryptResponse
	}
	if o.IncludeScoping != nil {
		m["includeScoping"] = o.IncludeScoping
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalSAMLConfigParams unmarshals an instance of SAMLConfigParams from the specified map of raw messages.
func UnmarshalSAMLConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLConfigParams)
	err = core.UnmarshalPrimitive(m, "entityID", &obj.EntityID)
	if err != nil {
		return
	}
	delete(m, "entityID")
	err = core.UnmarshalPrimitive(m, "signInUrl", &obj.SignInURL)
	if err != nil {
		return
	}
	delete(m, "signInUrl")
	err = core.UnmarshalPrimitive(m, "certificates", &obj.Certificates)
	if err != nil {
		return
	}
	delete(m, "certificates")
	err = core.UnmarshalPrimitive(m, "displayName", &obj.DisplayName)
	if err != nil {
		return
	}
	delete(m, "displayName")
	err = core.UnmarshalModel(m, "authnContext", &obj.AuthnContext, UnmarshalSAMLConfigParamsAuthnContext)
	if err != nil {
		return
	}
	delete(m, "authnContext")
	err = core.UnmarshalPrimitive(m, "signRequest", &obj.SignRequest)
	if err != nil {
		return
	}
	delete(m, "signRequest")
	err = core.UnmarshalPrimitive(m, "encryptResponse", &obj.EncryptResponse)
	if err != nil {
		return
	}
	delete(m, "encryptResponse")
	err = core.UnmarshalPrimitive(m, "includeScoping", &obj.IncludeScoping)
	if err != nil {
		return
	}
	delete(m, "includeScoping")
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLResponse : SAMLResponse struct
type SAMLResponse struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *SAMLConfigParams `json:"config,omitempty"`
}

// UnmarshalSAMLResponse unmarshals an instance of SAMLResponse from the specified map of raw messages.
func UnmarshalSAMLResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLResponse)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalSAMLConfigParams)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SAMLResponseWithValidationData : SAMLResponseWithValidationData struct
type SAMLResponseWithValidationData struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *SAMLConfigParams `json:"config,omitempty"`

	ValidationData *SAMLResponseWithValidationDataValidationData `json:"validation_data,omitempty"`
}

// UnmarshalSAMLResponseWithValidationData unmarshals an instance of SAMLResponseWithValidationData from the specified map of raw messages.
func UnmarshalSAMLResponseWithValidationData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SAMLResponseWithValidationData)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalSAMLConfigParams)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validation_data", &obj.ValidationData, UnmarshalSAMLResponseWithValidationDataValidationData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TokenClaimMapping : TokenClaimMapping struct
type TokenClaimMapping struct {
	Source *string `json:"source" validate:"required"`

	SourceClaim *string `json:"sourceClaim,omitempty"`

	DestinationClaim *string `json:"destinationClaim,omitempty"`
}

// Constants associated with the TokenClaimMapping.Source property.
const (
	TokenClaimMappingSourceAppidCustomConst    = "appid_custom"
	TokenClaimMappingSourceAttributesConst     = "attributes"
	TokenClaimMappingSourceCloudDirectoryConst = "cloud_directory"
	TokenClaimMappingSourceFacebookConst       = "facebook"
	TokenClaimMappingSourceGoogleConst         = "google"
	TokenClaimMappingSourceIbmidConst          = "ibmid"
	TokenClaimMappingSourceRolesConst          = "roles"
	TokenClaimMappingSourceSAMLConst           = "saml"
)

// NewTokenClaimMapping : Instantiate TokenClaimMapping (Generic Model Constructor)
func (*AppIDManagementV4) NewTokenClaimMapping(source string) (model *TokenClaimMapping, err error) {
	model = &TokenClaimMapping{
		Source: core.StringPtr(source),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalTokenClaimMapping unmarshals an instance of TokenClaimMapping from the specified map of raw messages.
func UnmarshalTokenClaimMapping(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TokenClaimMapping)
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sourceClaim", &obj.SourceClaim)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destinationClaim", &obj.DestinationClaim)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TokenConfigParams : TokenConfigParams struct
type TokenConfigParams struct {
	ExpiresIn *int64 `json:"expires_in" validate:"required"`

	Enabled *bool `json:"enabled" validate:"required"`
}

// NewTokenConfigParams : Instantiate TokenConfigParams (Generic Model Constructor)
func (*AppIDManagementV4) NewTokenConfigParams(expiresIn int64, enabled bool) (model *TokenConfigParams, err error) {
	model = &TokenConfigParams{
		ExpiresIn: core.Int64Ptr(expiresIn),
		Enabled:   core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalTokenConfigParams unmarshals an instance of TokenConfigParams from the specified map of raw messages.
func UnmarshalTokenConfigParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TokenConfigParams)
	err = core.UnmarshalPrimitive(m, "expires_in", &obj.ExpiresIn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TokensConfigResponse : TokensConfigResponse struct
type TokensConfigResponse struct {
	IDTokenClaims []TokenClaimMapping `json:"idTokenClaims,omitempty"`

	AccessTokenClaims []TokenClaimMapping `json:"accessTokenClaims,omitempty"`

	Access *AccessTokenConfigParams `json:"access,omitempty"`

	Refresh *TokenConfigParams `json:"refresh,omitempty"`

	AnonymousAccess *TokenConfigParams `json:"anonymousAccess,omitempty"`
}

// UnmarshalTokensConfigResponse unmarshals an instance of TokensConfigResponse from the specified map of raw messages.
func UnmarshalTokensConfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TokensConfigResponse)
	err = core.UnmarshalModel(m, "idTokenClaims", &obj.IDTokenClaims, UnmarshalTokenClaimMapping)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "accessTokenClaims", &obj.AccessTokenClaims, UnmarshalTokenClaimMapping)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "access", &obj.Access, UnmarshalAccessTokenConfigParams)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "refresh", &obj.Refresh, UnmarshalTokenConfigParams)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "anonymousAccess", &obj.AnonymousAccess, UnmarshalTokenConfigParams)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateExtensionConfig : UpdateExtensionConfig struct
type UpdateExtensionConfig struct {
	IsActive *bool `json:"isActive" validate:"required"`

	Config *UpdateExtensionConfigConfig `json:"config,omitempty"`
}

// NewUpdateExtensionConfig : Instantiate UpdateExtensionConfig (Generic Model Constructor)
func (*AppIDManagementV4) NewUpdateExtensionConfig(isActive bool) (model *UpdateExtensionConfig, err error) {
	model = &UpdateExtensionConfig{
		IsActive: core.BoolPtr(isActive),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalUpdateExtensionConfig unmarshals an instance of UpdateExtensionConfig from the specified map of raw messages.
func UnmarshalUpdateExtensionConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateExtensionConfig)
	err = core.UnmarshalPrimitive(m, "isActive", &obj.IsActive)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalUpdateExtensionConfigConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateRoleResponse : UpdateRoleResponse struct
type UpdateRoleResponse struct {
	ID *string `json:"id" validate:"required"`

	Name *string `json:"name" validate:"required"`

	Description *string `json:"description,omitempty"`

	Access []RoleAccessItem `json:"access" validate:"required"`
}

// UnmarshalUpdateRoleResponse unmarshals an instance of UpdateRoleResponse from the specified map of raw messages.
func UnmarshalUpdateRoleResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateRoleResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "access", &obj.Access, UnmarshalRoleAccessItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserSearchResponse : UserSearchResponse struct
type UserSearchResponse struct {
	TotalResults *int64 `json:"totalResults,omitempty"`

	ItemsPerPage *int64 `json:"itemsPerPage,omitempty"`

	RequestOptions *UserSearchResponseRequestOptions `json:"requestOptions,omitempty"`

	Users []UserSearchResponseUsersItem `json:"users,omitempty"`
}

// UnmarshalUserSearchResponse unmarshals an instance of UserSearchResponse from the specified map of raw messages.
func UnmarshalUserSearchResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserSearchResponse)
	err = core.UnmarshalPrimitive(m, "totalResults", &obj.TotalResults)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "itemsPerPage", &obj.ItemsPerPage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "requestOptions", &obj.RequestOptions, UnmarshalUserSearchResponseRequestOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserSearchResponseUsersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
