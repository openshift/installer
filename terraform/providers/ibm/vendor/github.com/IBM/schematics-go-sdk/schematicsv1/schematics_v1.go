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
 * IBM OpenAPI SDK Code Generator Version: 3.37.0-a85661cd-20210802-190136
 */

// Package schematicsv1 : Operations and models for the SchematicsV1 service
package schematicsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/schematics-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// SchematicsV1 : IBM Cloud Schematics service is to provide the capability to manage resources  of cloud provider
// infrastructure by using file based configurations.  With the IBM Cloud Schematics service you can specify the
// required set of resources and the configuration in `config files`,  and then pass the config files to the service to
// fulfill it by  calling the necessary actions on the infrastructure.  This principle is known as Infrastructure as
// Code.  For more information, refer to [Getting started with IBM Cloud Schematics]
// (https://cloud.ibm.com/docs/schematics?topic=schematics-getting-started).
//
// Version: 1.0
type SchematicsV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://schematics-dev.containers.appdomain.cloud"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "schematics"

// SchematicsV1Options : Service options
type SchematicsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSchematicsV1UsingExternalConfig : constructs an instance of SchematicsV1 with passed in options and external configuration.
func NewSchematicsV1UsingExternalConfig(options *SchematicsV1Options) (schematics *SchematicsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	schematics, err = NewSchematicsV1(options)
	if err != nil {
		return
	}

	err = schematics.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = schematics.Service.SetServiceURL(options.URL)
	}
	return
}

// NewSchematicsV1 : constructs an instance of SchematicsV1 with passed in options.
func NewSchematicsV1(options *SchematicsV1Options) (service *SchematicsV1, err error) {
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

	service = &SchematicsV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "schematics" suitable for processing requests.
func (schematics *SchematicsV1) Clone() *SchematicsV1 {
	if core.IsNil(schematics) {
		return nil
	}
	clone := *schematics
	clone.Service = schematics.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (schematics *SchematicsV1) SetServiceURL(url string) error {
	return schematics.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (schematics *SchematicsV1) GetServiceURL() string {
	return schematics.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (schematics *SchematicsV1) SetDefaultHeaders(headers http.Header) {
	schematics.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (schematics *SchematicsV1) SetEnableGzipCompression(enableGzip bool) {
	schematics.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (schematics *SchematicsV1) GetEnableGzipCompression() bool {
	return schematics.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (schematics *SchematicsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	schematics.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (schematics *SchematicsV1) DisableRetries() {
	schematics.Service.DisableRetries()
}

// ListSchematicsLocation : List supported schematics locations
// Retrieve a list of IBM Cloud locations where you can create Schematics workspaces.
func (schematics *SchematicsV1) ListSchematicsLocation(listSchematicsLocationOptions *ListSchematicsLocationOptions) (result []SchematicsLocations, response *core.DetailedResponse, err error) {
	return schematics.ListSchematicsLocationWithContext(context.Background(), listSchematicsLocationOptions)
}

// ListSchematicsLocationWithContext is an alternate form of the ListSchematicsLocation method which supports a Context parameter
func (schematics *SchematicsV1) ListSchematicsLocationWithContext(ctx context.Context, listSchematicsLocationOptions *ListSchematicsLocationOptions) (result []SchematicsLocations, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSchematicsLocationOptions, "listSchematicsLocationOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/locations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSchematicsLocationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListSchematicsLocation")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchematicsLocations)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLocations : List supported schematics locations
// Retrieve a list of IBM Cloud locations where you can work with Schematics objects.
func (schematics *SchematicsV1) ListLocations(listLocationsOptions *ListLocationsOptions) (result *SchematicsLocationsList, response *core.DetailedResponse, err error) {
	return schematics.ListLocationsWithContext(context.Background(), listLocationsOptions)
}

// ListLocationsWithContext is an alternate form of the ListLocations method which supports a Context parameter
func (schematics *SchematicsV1) ListLocationsWithContext(ctx context.Context, listLocationsOptions *ListLocationsOptions) (result *SchematicsLocationsList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLocationsOptions, "listLocationsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/locations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listLocationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListLocations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchematicsLocationsList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListResourceGroup : List resource groups
// Retrieve a list of IBM Cloud resource groups that your account has access to.
func (schematics *SchematicsV1) ListResourceGroup(listResourceGroupOptions *ListResourceGroupOptions) (result []ResourceGroupResponse, response *core.DetailedResponse, err error) {
	return schematics.ListResourceGroupWithContext(context.Background(), listResourceGroupOptions)
}

// ListResourceGroupWithContext is an alternate form of the ListResourceGroup method which supports a Context parameter
func (schematics *SchematicsV1) ListResourceGroupWithContext(ctx context.Context, listResourceGroupOptions *ListResourceGroupOptions) (result []ResourceGroupResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listResourceGroupOptions, "listResourceGroupOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/resource_groups`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listResourceGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListResourceGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroupResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSchematicsVersion : Get Schematics API information
// Retrieve detailed information about the IBM Cloud Schematics API version and the version of the provider plug-ins
// that the API uses.
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to the
// workspace ID and the resource group. For more information, about Schematics access and permissions, see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetSchematicsVersion(getSchematicsVersionOptions *GetSchematicsVersionOptions) (result *VersionResponse, response *core.DetailedResponse, err error) {
	return schematics.GetSchematicsVersionWithContext(context.Background(), getSchematicsVersionOptions)
}

// GetSchematicsVersionWithContext is an alternate form of the GetSchematicsVersion method which supports a Context parameter
func (schematics *SchematicsV1) GetSchematicsVersionWithContext(ctx context.Context, getSchematicsVersionOptions *GetSchematicsVersionOptions) (result *VersionResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSchematicsVersionOptions, "getSchematicsVersionOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/version`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchematicsVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetSchematicsVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVersionResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ProcessTemplateMetaData : Create metadata by  processing the template
// Create a Template metadata definition.
func (schematics *SchematicsV1) ProcessTemplateMetaData(processTemplateMetaDataOptions *ProcessTemplateMetaDataOptions) (result *TemplateMetaDataResponse, response *core.DetailedResponse, err error) {
	return schematics.ProcessTemplateMetaDataWithContext(context.Background(), processTemplateMetaDataOptions)
}

// ProcessTemplateMetaDataWithContext is an alternate form of the ProcessTemplateMetaData method which supports a Context parameter
func (schematics *SchematicsV1) ProcessTemplateMetaDataWithContext(ctx context.Context, processTemplateMetaDataOptions *ProcessTemplateMetaDataOptions) (result *TemplateMetaDataResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(processTemplateMetaDataOptions, "processTemplateMetaDataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(processTemplateMetaDataOptions, "processTemplateMetaDataOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/template_metadata_processor`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range processTemplateMetaDataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ProcessTemplateMetaData")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if processTemplateMetaDataOptions.XGithubToken != nil {
		builder.AddHeader("X-Github-token", fmt.Sprint(*processTemplateMetaDataOptions.XGithubToken))
	}

	body := make(map[string]interface{})
	if processTemplateMetaDataOptions.TemplateType != nil {
		body["template_type"] = processTemplateMetaDataOptions.TemplateType
	}
	if processTemplateMetaDataOptions.Source != nil {
		body["source"] = processTemplateMetaDataOptions.Source
	}
	if processTemplateMetaDataOptions.Region != nil {
		body["region"] = processTemplateMetaDataOptions.Region
	}
	if processTemplateMetaDataOptions.SourceType != nil {
		body["source_type"] = processTemplateMetaDataOptions.SourceType
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateMetaDataResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListWorkspaces : List workspaces
// Retrieve a list of Schematics workspaces from your IBM Cloud account that you have access to. The list of workspaces
// that is returned depends on the API endpoint that you use. For example, if you use an API endpoint for a geography,
// such as North America, only workspaces that are created in `us-south` or `us-east` are returned. For more information
// about supported API endpoints, see [API endpoints](/apidocs/schematics#api-endpoints).
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to the
// workspace ID and the resource group. For more information, about Schematics access and permissions, see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ListWorkspaces(listWorkspacesOptions *ListWorkspacesOptions) (result *WorkspaceResponseList, response *core.DetailedResponse, err error) {
	return schematics.ListWorkspacesWithContext(context.Background(), listWorkspacesOptions)
}

// ListWorkspacesWithContext is an alternate form of the ListWorkspaces method which supports a Context parameter
func (schematics *SchematicsV1) ListWorkspacesWithContext(ctx context.Context, listWorkspacesOptions *ListWorkspacesOptions) (result *WorkspaceResponseList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listWorkspacesOptions, "listWorkspacesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listWorkspacesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListWorkspaces")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listWorkspacesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listWorkspacesOptions.Offset))
	}
	if listWorkspacesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listWorkspacesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponseList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateWorkspace : Create a workspace
// Create an IBM Cloud Schematics workspace that points to the source repository where your Terraform template or the
// IBM Cloud software template is stored. You can decide to create your workspace without connecting it to a GitHub or
// GitLab repository. Your workspace is then created with a **Draft** state. To later connect your workspace to a GitHub
// or GitLab repository, you must use the `PUT /v1/workspaces/{id}` API to update the workspace or use the
// `/v1/workspaces/{id}/templates/{template_id}/template_repo_upload` API to upload a TAR file instead.
//
//  **Getting API endpoint**:-
//
//  * The Schematics API endpoint that you use to create the workspace determines where your Schematics actions run and
// your data is stored. See [API endpoints](/apidocs/schematics#api-endpoints) for more information.
//  * If you use the API endpoint for a geography and not a specific location, such as North America, you can specify
// the location in your API request body.
//  * If you do not specify the location in the request body, Schematics determines your workspace location based on
// availability.
//  * If you use an API endpoint for a specific location, such as Frankfurt, the location that you enter in your API
// request body must match your API endpoint.
//  * You also have the option to not specify a location in your API request body if you use a location-specific API
// endpoint.
//
//  **Getting IAM access token** :-
//  * Before you create Schematics workspace, you need to create the IAM access token for your IBM Cloud Account.
//  * To create IAM access token, use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>` and execute `curl -X POST
// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type= application/x-www-form-urlencoded" -d
// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`. For more information, about
// creating IAM access token and API Docs, see [IAM access token](/apidocs/iam-identity-token-api#gettoken-password) and
// [Create API key](/apidocs/iam-identity-token-api#create-api-key).
//  * You can set the environment values  `export ACCESS_TOKEN=<access_token>` and `export
// REFRESH_TOKEN=<refresh_token>`.
//  * You can use the obtained IAM access token in create workspace `curl` command.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) CreateWorkspace(createWorkspaceOptions *CreateWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	return schematics.CreateWorkspaceWithContext(context.Background(), createWorkspaceOptions)
}

// CreateWorkspaceWithContext is an alternate form of the CreateWorkspace method which supports a Context parameter
func (schematics *SchematicsV1) CreateWorkspaceWithContext(ctx context.Context, createWorkspaceOptions *CreateWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createWorkspaceOptions, "createWorkspaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createWorkspaceOptions, "createWorkspaceOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createWorkspaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateWorkspace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createWorkspaceOptions.XGithubToken != nil {
		builder.AddHeader("X-Github-token", fmt.Sprint(*createWorkspaceOptions.XGithubToken))
	}

	body := make(map[string]interface{})
	if createWorkspaceOptions.AppliedShareddataIds != nil {
		body["applied_shareddata_ids"] = createWorkspaceOptions.AppliedShareddataIds
	}
	if createWorkspaceOptions.CatalogRef != nil {
		body["catalog_ref"] = createWorkspaceOptions.CatalogRef
	}
	if createWorkspaceOptions.Description != nil {
		body["description"] = createWorkspaceOptions.Description
	}
	if createWorkspaceOptions.Location != nil {
		body["location"] = createWorkspaceOptions.Location
	}
	if createWorkspaceOptions.Name != nil {
		body["name"] = createWorkspaceOptions.Name
	}
	if createWorkspaceOptions.ResourceGroup != nil {
		body["resource_group"] = createWorkspaceOptions.ResourceGroup
	}
	if createWorkspaceOptions.SharedData != nil {
		body["shared_data"] = createWorkspaceOptions.SharedData
	}
	if createWorkspaceOptions.Tags != nil {
		body["tags"] = createWorkspaceOptions.Tags
	}
	if createWorkspaceOptions.TemplateData != nil {
		body["template_data"] = createWorkspaceOptions.TemplateData
	}
	if createWorkspaceOptions.TemplateRef != nil {
		body["template_ref"] = createWorkspaceOptions.TemplateRef
	}
	if createWorkspaceOptions.TemplateRepo != nil {
		body["template_repo"] = createWorkspaceOptions.TemplateRepo
	}
	if createWorkspaceOptions.Type != nil {
		body["type"] = createWorkspaceOptions.Type
	}
	if createWorkspaceOptions.WorkspaceStatus != nil {
		body["workspace_status"] = createWorkspaceOptions.WorkspaceStatus
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspace : Get workspace details
// Retrieve detailed information for a workspace in your IBM Cloud account.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform
//  access to the workspace ID and the resource group. For more information,
//  about Schematics access and permissions, see [Schematics service access
//  roles and required permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetWorkspace(getWorkspaceOptions *GetWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceWithContext(context.Background(), getWorkspaceOptions)
}

// GetWorkspaceWithContext is an alternate form of the GetWorkspace method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceWithContext(ctx context.Context, getWorkspaceOptions *GetWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceOptions, "getWorkspaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceOptions, "getWorkspaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceWorkspace : Update workspace
// Use this API to update or replace the entire workspace, including the Terraform template (`template_repo`) or IBM
// Cloud catalog software template (`catalog_ref`) that your workspace points to.
//
//  **Tip**:- If you want to update workspace metadata, use the `PATCH /v1/workspaces/{id}` API.
//  To update workspace variables, use the `PUT /v1/workspaces/{id}/template_data/{template_id}/values` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ReplaceWorkspace(replaceWorkspaceOptions *ReplaceWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	return schematics.ReplaceWorkspaceWithContext(context.Background(), replaceWorkspaceOptions)
}

// ReplaceWorkspaceWithContext is an alternate form of the ReplaceWorkspace method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceWorkspaceWithContext(ctx context.Context, replaceWorkspaceOptions *ReplaceWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceWorkspaceOptions, "replaceWorkspaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceWorkspaceOptions, "replaceWorkspaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *replaceWorkspaceOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceWorkspaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceWorkspace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceWorkspaceOptions.CatalogRef != nil {
		body["catalog_ref"] = replaceWorkspaceOptions.CatalogRef
	}
	if replaceWorkspaceOptions.Description != nil {
		body["description"] = replaceWorkspaceOptions.Description
	}
	if replaceWorkspaceOptions.Name != nil {
		body["name"] = replaceWorkspaceOptions.Name
	}
	if replaceWorkspaceOptions.SharedData != nil {
		body["shared_data"] = replaceWorkspaceOptions.SharedData
	}
	if replaceWorkspaceOptions.Tags != nil {
		body["tags"] = replaceWorkspaceOptions.Tags
	}
	if replaceWorkspaceOptions.TemplateData != nil {
		body["template_data"] = replaceWorkspaceOptions.TemplateData
	}
	if replaceWorkspaceOptions.TemplateRepo != nil {
		body["template_repo"] = replaceWorkspaceOptions.TemplateRepo
	}
	if replaceWorkspaceOptions.Type != nil {
		body["type"] = replaceWorkspaceOptions.Type
	}
	if replaceWorkspaceOptions.WorkspaceStatus != nil {
		body["workspace_status"] = replaceWorkspaceOptions.WorkspaceStatus
	}
	if replaceWorkspaceOptions.WorkspaceStatusMsg != nil {
		body["workspace_status_msg"] = replaceWorkspaceOptions.WorkspaceStatusMsg
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteWorkspace : Delete a workspace
// Deletes a workspace from IBM Cloud Schematics. Deleting a workspace does not automatically remove the IBM Cloud
// resources that the workspace manages. To remove all resources that are associated with the workspace, use the `DELETE
// /v1/workspaces/{id}?destroy_resources=true` API.
//
//  **Note**: If you delete a workspace without deleting the resources,
//  you must manage your resources with the resource dashboard or CLI afterwards.
//  You cannot use IBM Cloud Schematics anymore to manage your resources.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) DeleteWorkspace(deleteWorkspaceOptions *DeleteWorkspaceOptions) (result *string, response *core.DetailedResponse, err error) {
	return schematics.DeleteWorkspaceWithContext(context.Background(), deleteWorkspaceOptions)
}

// DeleteWorkspaceWithContext is an alternate form of the DeleteWorkspace method which supports a Context parameter
func (schematics *SchematicsV1) DeleteWorkspaceWithContext(ctx context.Context, deleteWorkspaceOptions *DeleteWorkspaceOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteWorkspaceOptions, "deleteWorkspaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteWorkspaceOptions, "deleteWorkspaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *deleteWorkspaceOptions.WID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteWorkspaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteWorkspace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteWorkspaceOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*deleteWorkspaceOptions.RefreshToken))
	}

	if deleteWorkspaceOptions.DestroyResources != nil {
		builder.AddQuery("destroy_resources", fmt.Sprint(*deleteWorkspaceOptions.DestroyResources))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, &result)

	return
}

// UpdateWorkspace : Update workspace metadata
// Use this API to update the following workspace metadata:
//
//  * Workspace name (`name`) - **Note**: Updating the workspace name does not update the ID of the workspace.
//  * Workspace description (`description`)
//  * Tags (`tags[]`)
//  * Resource group (`resource_group`)
//  * Workspace status (`workspace_status.frozen`)
//
//
//  **Tip**: If you want to update information about the Terraform template
//  or IBM Cloud catalog software template that your workspace points to,
//  use the `PUT /v1/workspaces/{id}` API. To update workspace variables,
//  use the `PUT /v1/workspaces/{id}/template_data/{template_id}/values` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) UpdateWorkspace(updateWorkspaceOptions *UpdateWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	return schematics.UpdateWorkspaceWithContext(context.Background(), updateWorkspaceOptions)
}

// UpdateWorkspaceWithContext is an alternate form of the UpdateWorkspace method which supports a Context parameter
func (schematics *SchematicsV1) UpdateWorkspaceWithContext(ctx context.Context, updateWorkspaceOptions *UpdateWorkspaceOptions) (result *WorkspaceResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateWorkspaceOptions, "updateWorkspaceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateWorkspaceOptions, "updateWorkspaceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *updateWorkspaceOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateWorkspaceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UpdateWorkspace")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateWorkspaceOptions.CatalogRef != nil {
		body["catalog_ref"] = updateWorkspaceOptions.CatalogRef
	}
	if updateWorkspaceOptions.Description != nil {
		body["description"] = updateWorkspaceOptions.Description
	}
	if updateWorkspaceOptions.Name != nil {
		body["name"] = updateWorkspaceOptions.Name
	}
	if updateWorkspaceOptions.SharedData != nil {
		body["shared_data"] = updateWorkspaceOptions.SharedData
	}
	if updateWorkspaceOptions.Tags != nil {
		body["tags"] = updateWorkspaceOptions.Tags
	}
	if updateWorkspaceOptions.TemplateData != nil {
		body["template_data"] = updateWorkspaceOptions.TemplateData
	}
	if updateWorkspaceOptions.TemplateRepo != nil {
		body["template_repo"] = updateWorkspaceOptions.TemplateRepo
	}
	if updateWorkspaceOptions.Type != nil {
		body["type"] = updateWorkspaceOptions.Type
	}
	if updateWorkspaceOptions.WorkspaceStatus != nil {
		body["workspace_status"] = updateWorkspaceOptions.WorkspaceStatus
	}
	if updateWorkspaceOptions.WorkspaceStatusMsg != nil {
		body["workspace_status_msg"] = updateWorkspaceOptions.WorkspaceStatusMsg
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceReadme : Show workspace template readme
// Retrieve the `README.md` file of the Terraform of IBM Cloud catalog template that your workspace points to.
func (schematics *SchematicsV1) GetWorkspaceReadme(getWorkspaceReadmeOptions *GetWorkspaceReadmeOptions) (result *TemplateReadme, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceReadmeWithContext(context.Background(), getWorkspaceReadmeOptions)
}

// GetWorkspaceReadmeWithContext is an alternate form of the GetWorkspaceReadme method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceReadmeWithContext(ctx context.Context, getWorkspaceReadmeOptions *GetWorkspaceReadmeOptions) (result *TemplateReadme, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceReadmeOptions, "getWorkspaceReadmeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceReadmeOptions, "getWorkspaceReadmeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceReadmeOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/templates/readme`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceReadmeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceReadme")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getWorkspaceReadmeOptions.Ref != nil {
		builder.AddQuery("ref", fmt.Sprint(*getWorkspaceReadmeOptions.Ref))
	}
	if getWorkspaceReadmeOptions.Formatted != nil {
		builder.AddQuery("formatted", fmt.Sprint(*getWorkspaceReadmeOptions.Formatted))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateReadme)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// TemplateRepoUpload : Upload a TAR file to your workspace
// Provide your Terraform template by uploading a TAR file from your local machine. Before you use this API, you must
// create a workspace without a link to a GitHub or GitLab repository with the `POST /v1/workspaces` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access to the workspace ID and the
// resource group. For more information, about Schematics access and permissions, see [Schematics service access roles
// and required permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) TemplateRepoUpload(templateRepoUploadOptions *TemplateRepoUploadOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	return schematics.TemplateRepoUploadWithContext(context.Background(), templateRepoUploadOptions)
}

// TemplateRepoUploadWithContext is an alternate form of the TemplateRepoUpload method which supports a Context parameter
func (schematics *SchematicsV1) TemplateRepoUploadWithContext(ctx context.Context, templateRepoUploadOptions *TemplateRepoUploadOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(templateRepoUploadOptions, "templateRepoUploadOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(templateRepoUploadOptions, "templateRepoUploadOptions")
	if err != nil {
		return
	}
	if templateRepoUploadOptions.File == nil {
		err = fmt.Errorf("at least one of  or file must be supplied")
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *templateRepoUploadOptions.WID,
		"t_id": *templateRepoUploadOptions.TID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/template_data/{t_id}/template_repo_upload`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range templateRepoUploadOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "TemplateRepoUpload")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if templateRepoUploadOptions.File != nil {
		builder.AddFormData("file", "filename",
			core.StringNilMapper(templateRepoUploadOptions.FileContentType), templateRepoUploadOptions.File)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateRepoTarUploadResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceInputs : List workspace input variables
// Retrieve a list of input variables that are declared in your Terraform or IBM Cloud catalog template.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetWorkspaceInputs(getWorkspaceInputsOptions *GetWorkspaceInputsOptions) (result *TemplateValues, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceInputsWithContext(context.Background(), getWorkspaceInputsOptions)
}

// GetWorkspaceInputsWithContext is an alternate form of the GetWorkspaceInputs method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceInputsWithContext(ctx context.Context, getWorkspaceInputsOptions *GetWorkspaceInputsOptions) (result *TemplateValues, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceInputsOptions, "getWorkspaceInputsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceInputsOptions, "getWorkspaceInputsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceInputsOptions.WID,
		"t_id": *getWorkspaceInputsOptions.TID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/template_data/{t_id}/values`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceInputsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceInputs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateValues)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceWorkspaceInputs : Replace workspace input variables
// Replace or Update the input variables for the template that your workspace points to.
func (schematics *SchematicsV1) ReplaceWorkspaceInputs(replaceWorkspaceInputsOptions *ReplaceWorkspaceInputsOptions) (result *UserValues, response *core.DetailedResponse, err error) {
	return schematics.ReplaceWorkspaceInputsWithContext(context.Background(), replaceWorkspaceInputsOptions)
}

// ReplaceWorkspaceInputsWithContext is an alternate form of the ReplaceWorkspaceInputs method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceWorkspaceInputsWithContext(ctx context.Context, replaceWorkspaceInputsOptions *ReplaceWorkspaceInputsOptions) (result *UserValues, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceWorkspaceInputsOptions, "replaceWorkspaceInputsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceWorkspaceInputsOptions, "replaceWorkspaceInputsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *replaceWorkspaceInputsOptions.WID,
		"t_id": *replaceWorkspaceInputsOptions.TID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/template_data/{t_id}/values`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceWorkspaceInputsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceWorkspaceInputs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceWorkspaceInputsOptions.EnvValues != nil {
		body["env_values"] = replaceWorkspaceInputsOptions.EnvValues
	}
	if replaceWorkspaceInputsOptions.Values != nil {
		body["values"] = replaceWorkspaceInputsOptions.Values
	}
	if replaceWorkspaceInputsOptions.Variablestore != nil {
		body["variablestore"] = replaceWorkspaceInputsOptions.Variablestore
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserValues)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAllWorkspaceInputs : Get workspace template details
// Retrieve detailed information about the Terraform template that your workspace points to.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform
//  access to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetAllWorkspaceInputs(getAllWorkspaceInputsOptions *GetAllWorkspaceInputsOptions) (result *WorkspaceTemplateValuesResponse, response *core.DetailedResponse, err error) {
	return schematics.GetAllWorkspaceInputsWithContext(context.Background(), getAllWorkspaceInputsOptions)
}

// GetAllWorkspaceInputsWithContext is an alternate form of the GetAllWorkspaceInputs method which supports a Context parameter
func (schematics *SchematicsV1) GetAllWorkspaceInputsWithContext(ctx context.Context, getAllWorkspaceInputsOptions *GetAllWorkspaceInputsOptions) (result *WorkspaceTemplateValuesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAllWorkspaceInputsOptions, "getAllWorkspaceInputsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAllWorkspaceInputsOptions, "getAllWorkspaceInputsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getAllWorkspaceInputsOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/templates/values`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAllWorkspaceInputsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetAllWorkspaceInputs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceTemplateValuesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceInputMetadata : List workspace variable metadata
// Retrieve the metadata for all the workspace input variables that are declared in the template that your workspace
// points to.
func (schematics *SchematicsV1) GetWorkspaceInputMetadata(getWorkspaceInputMetadataOptions *GetWorkspaceInputMetadataOptions) (result []interface{}, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceInputMetadataWithContext(context.Background(), getWorkspaceInputMetadataOptions)
}

// GetWorkspaceInputMetadataWithContext is an alternate form of the GetWorkspaceInputMetadata method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceInputMetadataWithContext(ctx context.Context, getWorkspaceInputMetadataOptions *GetWorkspaceInputMetadataOptions) (result []interface{}, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceInputMetadataOptions, "getWorkspaceInputMetadataOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceInputMetadataOptions, "getWorkspaceInputMetadataOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceInputMetadataOptions.WID,
		"t_id": *getWorkspaceInputMetadataOptions.TID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/template_data/{t_id}/values_metadata`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceInputMetadataOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceInputMetadata")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, &result)

	return
}

// GetWorkspaceOutputs : List workspace output values
// Retrieve a list of Terraform output variables. You define output values in your Terraform template to include
// information that you want to make accessible for other Terraform templates.
func (schematics *SchematicsV1) GetWorkspaceOutputs(getWorkspaceOutputsOptions *GetWorkspaceOutputsOptions) (result []OutputValuesInner, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceOutputsWithContext(context.Background(), getWorkspaceOutputsOptions)
}

// GetWorkspaceOutputsWithContext is an alternate form of the GetWorkspaceOutputs method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceOutputsWithContext(ctx context.Context, getWorkspaceOutputsOptions *GetWorkspaceOutputsOptions) (result []OutputValuesInner, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceOutputsOptions, "getWorkspaceOutputsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceOutputsOptions, "getWorkspaceOutputsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceOutputsOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/output_values`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceOutputsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceOutputs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOutputValuesInner)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceResources : List workspace resources
// Retrieve a list of IBM Cloud resources that you created with your workspace.
func (schematics *SchematicsV1) GetWorkspaceResources(getWorkspaceResourcesOptions *GetWorkspaceResourcesOptions) (result []TemplateResources, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceResourcesWithContext(context.Background(), getWorkspaceResourcesOptions)
}

// GetWorkspaceResourcesWithContext is an alternate form of the GetWorkspaceResources method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceResourcesWithContext(ctx context.Context, getWorkspaceResourcesOptions *GetWorkspaceResourcesOptions) (result []TemplateResources, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceResourcesOptions, "getWorkspaceResourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceResourcesOptions, "getWorkspaceResourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceResourcesOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceResourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceResources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse []json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateResources)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceState : Get Terraform statefile URL
// Retrieve the URL to the Terraform statefile (`terraform.tfstate`). You use the URL to access the Terraform statefile.
// The Terraform statefile includes detailed information about the IBM Cloud resources that you provisioned with IBM
// Cloud Schematics and Schematics uses the file to determine future create, modify, or delete actions for your
// resources. To show the content of the Terraform statefile, use the `GET
// /v1/workspaces/{id}/runtime_data/{template_id}/state_store` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetWorkspaceState(getWorkspaceStateOptions *GetWorkspaceStateOptions) (result *StateStoreResponseList, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceStateWithContext(context.Background(), getWorkspaceStateOptions)
}

// GetWorkspaceStateWithContext is an alternate form of the GetWorkspaceState method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceStateWithContext(ctx context.Context, getWorkspaceStateOptions *GetWorkspaceStateOptions) (result *StateStoreResponseList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceStateOptions, "getWorkspaceStateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceStateOptions, "getWorkspaceStateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceStateOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/state_stores`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceStateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceState")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStateStoreResponseList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceTemplateState : Show Terraform statefile content
// Show the content of the Terraform statefile (`terraform.tfstate`) that was created when your Terraform template was
// applied in IBM Cloud. The statefile holds detailed information about the IBM Cloud resources that were created by IBM
// Cloud Schematics and Schematics uses the file to determine future create, modify, or delete actions for your
// resources.
func (schematics *SchematicsV1) GetWorkspaceTemplateState(getWorkspaceTemplateStateOptions *GetWorkspaceTemplateStateOptions) (result []interface{}, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceTemplateStateWithContext(context.Background(), getWorkspaceTemplateStateOptions)
}

// GetWorkspaceTemplateStateWithContext is an alternate form of the GetWorkspaceTemplateState method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceTemplateStateWithContext(ctx context.Context, getWorkspaceTemplateStateOptions *GetWorkspaceTemplateStateOptions) (result []interface{}, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceTemplateStateOptions, "getWorkspaceTemplateStateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceTemplateStateOptions, "getWorkspaceTemplateStateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceTemplateStateOptions.WID,
		"t_id": *getWorkspaceTemplateStateOptions.TID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/runtime_data/{t_id}/state_store`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceTemplateStateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceTemplateState")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, &result)

	return
}

// GetWorkspaceActivityLogs : Get workspace job log URL
// Get the Terraform log file URL for a workspace job. You can retrieve the log URL for jobs that were created with the
// `PUT /v1/workspaces/{id}/apply`, `POST /v1/workspaces/{id}/plan`, or `DELETE /v1/workspaces/{id}/destroy` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetWorkspaceActivityLogs(getWorkspaceActivityLogsOptions *GetWorkspaceActivityLogsOptions) (result *WorkspaceActivityLogs, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceActivityLogsWithContext(context.Background(), getWorkspaceActivityLogsOptions)
}

// GetWorkspaceActivityLogsWithContext is an alternate form of the GetWorkspaceActivityLogs method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceActivityLogsWithContext(ctx context.Context, getWorkspaceActivityLogsOptions *GetWorkspaceActivityLogsOptions) (result *WorkspaceActivityLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceActivityLogsOptions, "getWorkspaceActivityLogsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceActivityLogsOptions, "getWorkspaceActivityLogsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id":        *getWorkspaceActivityLogsOptions.WID,
		"activity_id": *getWorkspaceActivityLogsOptions.ActivityID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/actions/{activity_id}/logs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceActivityLogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceActivityLogs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceLogUrls : Get latest workspace job log URL for all workspace templates
// Retrieve the log file URL for the latest job of a template that ran against your workspace. You use this URL to
// retrieve detailed logs for the latest job.
func (schematics *SchematicsV1) GetWorkspaceLogUrls(getWorkspaceLogUrlsOptions *GetWorkspaceLogUrlsOptions) (result *LogStoreResponseList, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceLogUrlsWithContext(context.Background(), getWorkspaceLogUrlsOptions)
}

// GetWorkspaceLogUrlsWithContext is an alternate form of the GetWorkspaceLogUrls method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceLogUrlsWithContext(ctx context.Context, getWorkspaceLogUrlsOptions *GetWorkspaceLogUrlsOptions) (result *LogStoreResponseList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceLogUrlsOptions, "getWorkspaceLogUrlsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceLogUrlsOptions, "getWorkspaceLogUrlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getWorkspaceLogUrlsOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/log_stores`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceLogUrlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceLogUrls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogStoreResponseList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTemplateLogs : Show logs for the latest action for a workspace template
// Show the Terraform logs for the most recent job of a template that ran against your workspace.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access to the workspace ID and the
// resource group. For more information, about Schematics access and permissions, see [Schematics service access roles
// and required permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetTemplateLogs(getTemplateLogsOptions *GetTemplateLogsOptions) (result *string, response *core.DetailedResponse, err error) {
	return schematics.GetTemplateLogsWithContext(context.Background(), getTemplateLogsOptions)
}

// GetTemplateLogsWithContext is an alternate form of the GetTemplateLogs method which supports a Context parameter
func (schematics *SchematicsV1) GetTemplateLogsWithContext(ctx context.Context, getTemplateLogsOptions *GetTemplateLogsOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTemplateLogsOptions, "getTemplateLogsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTemplateLogsOptions, "getTemplateLogsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *getTemplateLogsOptions.WID,
		"t_id": *getTemplateLogsOptions.TID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/runtime_data/{t_id}/log_store`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTemplateLogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetTemplateLogs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getTemplateLogsOptions.LogTfCmd != nil {
		builder.AddQuery("log_tf_cmd", fmt.Sprint(*getTemplateLogsOptions.LogTfCmd))
	}
	if getTemplateLogsOptions.LogTfPrefix != nil {
		builder.AddQuery("log_tf_prefix", fmt.Sprint(*getTemplateLogsOptions.LogTfPrefix))
	}
	if getTemplateLogsOptions.LogTfNullResource != nil {
		builder.AddQuery("log_tf_null_resource", fmt.Sprint(*getTemplateLogsOptions.LogTfNullResource))
	}
	if getTemplateLogsOptions.LogTfAnsible != nil {
		builder.AddQuery("log_tf_ansible", fmt.Sprint(*getTemplateLogsOptions.LogTfAnsible))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, &result)

	return
}

// GetTemplateActivityLog : Show logs for a workspace job
// Show the Terraform logs for an job that ran against your workspace.
func (schematics *SchematicsV1) GetTemplateActivityLog(getTemplateActivityLogOptions *GetTemplateActivityLogOptions) (result *string, response *core.DetailedResponse, err error) {
	return schematics.GetTemplateActivityLogWithContext(context.Background(), getTemplateActivityLogOptions)
}

// GetTemplateActivityLogWithContext is an alternate form of the GetTemplateActivityLog method which supports a Context parameter
func (schematics *SchematicsV1) GetTemplateActivityLogWithContext(ctx context.Context, getTemplateActivityLogOptions *GetTemplateActivityLogOptions) (result *string, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTemplateActivityLogOptions, "getTemplateActivityLogOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTemplateActivityLogOptions, "getTemplateActivityLogOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id":        *getTemplateActivityLogOptions.WID,
		"t_id":        *getTemplateActivityLogOptions.TID,
		"activity_id": *getTemplateActivityLogOptions.ActivityID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/runtime_data/{t_id}/log_store/actions/{activity_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTemplateActivityLogOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetTemplateActivityLog")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getTemplateActivityLogOptions.LogTfCmd != nil {
		builder.AddQuery("log_tf_cmd", fmt.Sprint(*getTemplateActivityLogOptions.LogTfCmd))
	}
	if getTemplateActivityLogOptions.LogTfPrefix != nil {
		builder.AddQuery("log_tf_prefix", fmt.Sprint(*getTemplateActivityLogOptions.LogTfPrefix))
	}
	if getTemplateActivityLogOptions.LogTfNullResource != nil {
		builder.AddQuery("log_tf_null_resource", fmt.Sprint(*getTemplateActivityLogOptions.LogTfNullResource))
	}
	if getTemplateActivityLogOptions.LogTfAnsible != nil {
		builder.AddQuery("log_tf_ansible", fmt.Sprint(*getTemplateActivityLogOptions.LogTfAnsible))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, &result)

	return
}

// ListActions : List actions
// Retrieve a list of all Schematics actions that depends on the API endpoint that you have access. For example, if you
// use an API endpoint for a geography, such as North America, only actions that are created in `us-south` or `us-east`
// are retrieved.
//
//  For more information, about supported API endpoints, see
// [API endpoints](/apidocs/schematics#api-endpoints).
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to an action ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ListActions(listActionsOptions *ListActionsOptions) (result *ActionList, response *core.DetailedResponse, err error) {
	return schematics.ListActionsWithContext(context.Background(), listActionsOptions)
}

// ListActionsWithContext is an alternate form of the ListActions method which supports a Context parameter
func (schematics *SchematicsV1) ListActionsWithContext(ctx context.Context, listActionsOptions *ListActionsOptions) (result *ActionList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listActionsOptions, "listActionsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listActionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListActions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listActionsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listActionsOptions.Offset))
	}
	if listActionsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listActionsOptions.Limit))
	}
	if listActionsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listActionsOptions.Sort))
	}
	if listActionsOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*listActionsOptions.Profile))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActionList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAction : Create an action
// Create an IBM Cloud Schematics action to run on a single target or groups of target hosts, roles, policies, or steps
// to deploy your resources in the target hosts. You can run the IBM Cloud resources the order in which you want to
// execute them. **Note** If your Git repository already contains a host file. Schematics does not overwrite the host
// file already present in your Git repository. For sample templates, see IBM Cloud Automation
// [templates](https://github.com/Cloud-Schematics).
//
//  For more information, about the Schematics create action,
//  see [ibmcloud schematics action
// create](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-create-action).
//  **Note** you cannot update the location and region once an action is created.
//  Also, make sure your IP addresses are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to an action ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions.
func (schematics *SchematicsV1) CreateAction(createActionOptions *CreateActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	return schematics.CreateActionWithContext(context.Background(), createActionOptions)
}

// CreateActionWithContext is an alternate form of the CreateAction method which supports a Context parameter
func (schematics *SchematicsV1) CreateActionWithContext(ctx context.Context, createActionOptions *CreateActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createActionOptions, "createActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createActionOptions, "createActionOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createActionOptions.XGithubToken != nil {
		builder.AddHeader("X-Github-token", fmt.Sprint(*createActionOptions.XGithubToken))
	}

	body := make(map[string]interface{})
	if createActionOptions.Name != nil {
		body["name"] = createActionOptions.Name
	}
	if createActionOptions.Description != nil {
		body["description"] = createActionOptions.Description
	}
	if createActionOptions.Location != nil {
		body["location"] = createActionOptions.Location
	}
	if createActionOptions.ResourceGroup != nil {
		body["resource_group"] = createActionOptions.ResourceGroup
	}
	if createActionOptions.Tags != nil {
		body["tags"] = createActionOptions.Tags
	}
	if createActionOptions.UserState != nil {
		body["user_state"] = createActionOptions.UserState
	}
	if createActionOptions.SourceReadmeURL != nil {
		body["source_readme_url"] = createActionOptions.SourceReadmeURL
	}
	if createActionOptions.Source != nil {
		body["source"] = createActionOptions.Source
	}
	if createActionOptions.SourceType != nil {
		body["source_type"] = createActionOptions.SourceType
	}
	if createActionOptions.CommandParameter != nil {
		body["command_parameter"] = createActionOptions.CommandParameter
	}
	if createActionOptions.Inventory != nil {
		body["inventory"] = createActionOptions.Inventory
	}
	if createActionOptions.Credentials != nil {
		body["credentials"] = createActionOptions.Credentials
	}
	if createActionOptions.Bastion != nil {
		body["bastion"] = createActionOptions.Bastion
	}
	if createActionOptions.BastionCredential != nil {
		body["bastion_credential"] = createActionOptions.BastionCredential
	}
	if createActionOptions.TargetsIni != nil {
		body["targets_ini"] = createActionOptions.TargetsIni
	}
	if createActionOptions.Inputs != nil {
		body["inputs"] = createActionOptions.Inputs
	}
	if createActionOptions.Outputs != nil {
		body["outputs"] = createActionOptions.Outputs
	}
	if createActionOptions.Settings != nil {
		body["settings"] = createActionOptions.Settings
	}
	if createActionOptions.State != nil {
		body["state"] = createActionOptions.State
	}
	if createActionOptions.SysLock != nil {
		body["sys_lock"] = createActionOptions.SysLock
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAction : Get action details
// Retrieve the detailed information of an actions from your IBM Cloud account.  This API returns a URL to the log file
// that you can retrieve by using  the `GET /v2/actions/{action_id}/logs` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to an action ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) GetAction(getActionOptions *GetActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	return schematics.GetActionWithContext(context.Background(), getActionOptions)
}

// GetActionWithContext is an alternate form of the GetAction method which supports a Context parameter
func (schematics *SchematicsV1) GetActionWithContext(ctx context.Context, getActionOptions *GetActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getActionOptions, "getActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getActionOptions, "getActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"action_id": *getActionOptions.ActionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions/{action_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getActionOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*getActionOptions.Profile))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAction : Delete an action
// Delete a Schematics action and specify the Ansible playbook that you want to run against your IBM Cloud resources.
// **Note** you cannot delete or stop the job activity from an ongoing execution of an action defined in the playbook.
// You can repeat the execution of same job, whenever you patch the actions. For more information, about the Schematics
// action state, see  [Schematics action state
// diagram](https://cloud.ibm.com/docs/schematics?topic=schematics-action-setup#action-state-diagram).
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to an action ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions.
func (schematics *SchematicsV1) DeleteAction(deleteActionOptions *DeleteActionOptions) (response *core.DetailedResponse, err error) {
	return schematics.DeleteActionWithContext(context.Background(), deleteActionOptions)
}

// DeleteActionWithContext is an alternate form of the DeleteAction method which supports a Context parameter
func (schematics *SchematicsV1) DeleteActionWithContext(ctx context.Context, deleteActionOptions *DeleteActionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteActionOptions, "deleteActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteActionOptions, "deleteActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"action_id": *deleteActionOptions.ActionID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions/{action_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteActionOptions.Force != nil {
		builder.AddHeader("force", fmt.Sprint(*deleteActionOptions.Force))
	}
	if deleteActionOptions.Propagate != nil {
		builder.AddHeader("propagate", fmt.Sprint(*deleteActionOptions.Propagate))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, nil)

	return
}

// UpdateAction : Update an action
// Update or replace an action to change the action state from the critical state to normal state, or pending state to
// the normal state for a successful execution.  For more information, about the Schematics action state, see
// [Schematics action state
// diagram](https://cloud.ibm.com/docs/schematics?topic=schematics-action-setup#action-state-diagram).
//
//  **Note** you cannot update the location and region once an action is created.
//  Also, make sure your IP addresses are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses].
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to an action ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions.
func (schematics *SchematicsV1) UpdateAction(updateActionOptions *UpdateActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	return schematics.UpdateActionWithContext(context.Background(), updateActionOptions)
}

// UpdateActionWithContext is an alternate form of the UpdateAction method which supports a Context parameter
func (schematics *SchematicsV1) UpdateActionWithContext(ctx context.Context, updateActionOptions *UpdateActionOptions) (result *Action, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateActionOptions, "updateActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateActionOptions, "updateActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"action_id": *updateActionOptions.ActionID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions/{action_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UpdateAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateActionOptions.XGithubToken != nil {
		builder.AddHeader("X-Github-token", fmt.Sprint(*updateActionOptions.XGithubToken))
	}

	body := make(map[string]interface{})
	if updateActionOptions.Name != nil {
		body["name"] = updateActionOptions.Name
	}
	if updateActionOptions.Description != nil {
		body["description"] = updateActionOptions.Description
	}
	if updateActionOptions.Location != nil {
		body["location"] = updateActionOptions.Location
	}
	if updateActionOptions.ResourceGroup != nil {
		body["resource_group"] = updateActionOptions.ResourceGroup
	}
	if updateActionOptions.Tags != nil {
		body["tags"] = updateActionOptions.Tags
	}
	if updateActionOptions.UserState != nil {
		body["user_state"] = updateActionOptions.UserState
	}
	if updateActionOptions.SourceReadmeURL != nil {
		body["source_readme_url"] = updateActionOptions.SourceReadmeURL
	}
	if updateActionOptions.Source != nil {
		body["source"] = updateActionOptions.Source
	}
	if updateActionOptions.SourceType != nil {
		body["source_type"] = updateActionOptions.SourceType
	}
	if updateActionOptions.CommandParameter != nil {
		body["command_parameter"] = updateActionOptions.CommandParameter
	}
	if updateActionOptions.Inventory != nil {
		body["inventory"] = updateActionOptions.Inventory
	}
	if updateActionOptions.Credentials != nil {
		body["credentials"] = updateActionOptions.Credentials
	}
	if updateActionOptions.Bastion != nil {
		body["bastion"] = updateActionOptions.Bastion
	}
	if updateActionOptions.BastionCredential != nil {
		body["bastion_credential"] = updateActionOptions.BastionCredential
	}
	if updateActionOptions.TargetsIni != nil {
		body["targets_ini"] = updateActionOptions.TargetsIni
	}
	if updateActionOptions.Inputs != nil {
		body["inputs"] = updateActionOptions.Inputs
	}
	if updateActionOptions.Outputs != nil {
		body["outputs"] = updateActionOptions.Outputs
	}
	if updateActionOptions.Settings != nil {
		body["settings"] = updateActionOptions.Settings
	}
	if updateActionOptions.State != nil {
		body["state"] = updateActionOptions.State
	}
	if updateActionOptions.SysLock != nil {
		body["sys_lock"] = updateActionOptions.SysLock
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UploadTemplateTarAction : Upload a TAR file to an action
// Update your template by uploading tape archive file (.tar) file from  your local machine. Before you use this API,
// you must create an action  without a link to a GitHub or GitLab repository with the `POST /v2/actions` API.
//
// <h3>Authorization</h3>
//   Schematics support generic authorization such as service access or  platform access to an action ID and the
// resource group.  For more information, about Schematics access and permissions, see  [Schematics service access roles
// and required permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions.
func (schematics *SchematicsV1) UploadTemplateTarAction(uploadTemplateTarActionOptions *UploadTemplateTarActionOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	return schematics.UploadTemplateTarActionWithContext(context.Background(), uploadTemplateTarActionOptions)
}

// UploadTemplateTarActionWithContext is an alternate form of the UploadTemplateTarAction method which supports a Context parameter
func (schematics *SchematicsV1) UploadTemplateTarActionWithContext(ctx context.Context, uploadTemplateTarActionOptions *UploadTemplateTarActionOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(uploadTemplateTarActionOptions, "uploadTemplateTarActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(uploadTemplateTarActionOptions, "uploadTemplateTarActionOptions")
	if err != nil {
		return
	}
	if uploadTemplateTarActionOptions.File == nil {
		err = fmt.Errorf("at least one of  or file must be supplied")
		return
	}

	pathParamsMap := map[string]string{
		"action_id": *uploadTemplateTarActionOptions.ActionID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/actions/{action_id}/template_repo_upload`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadTemplateTarActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UploadTemplateTarAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if uploadTemplateTarActionOptions.File != nil {
		builder.AddFormData("file", "filename",
			core.StringNilMapper(uploadTemplateTarActionOptions.FileContentType), uploadTemplateTarActionOptions.File)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateRepoTarUploadResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListWorkspaceActivities : List all workspace jobs
// Retrieve a list of all jobs that ran against a workspace. Jobs are generated when you use the `apply`, `plan`,
// `destroy`, and `refresh`,   command API.
func (schematics *SchematicsV1) ListWorkspaceActivities(listWorkspaceActivitiesOptions *ListWorkspaceActivitiesOptions) (result *WorkspaceActivities, response *core.DetailedResponse, err error) {
	return schematics.ListWorkspaceActivitiesWithContext(context.Background(), listWorkspaceActivitiesOptions)
}

// ListWorkspaceActivitiesWithContext is an alternate form of the ListWorkspaceActivities method which supports a Context parameter
func (schematics *SchematicsV1) ListWorkspaceActivitiesWithContext(ctx context.Context, listWorkspaceActivitiesOptions *ListWorkspaceActivitiesOptions) (result *WorkspaceActivities, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listWorkspaceActivitiesOptions, "listWorkspaceActivitiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listWorkspaceActivitiesOptions, "listWorkspaceActivitiesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *listWorkspaceActivitiesOptions.WID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/actions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listWorkspaceActivitiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListWorkspaceActivities")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listWorkspaceActivitiesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listWorkspaceActivitiesOptions.Offset))
	}
	if listWorkspaceActivitiesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listWorkspaceActivitiesOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivities)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceActivity : Get workspace job details
// Get the details for a workspace job that ran against the workspace. This API returns the job status and a URL to the
// log file that you can  retrieve by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
func (schematics *SchematicsV1) GetWorkspaceActivity(getWorkspaceActivityOptions *GetWorkspaceActivityOptions) (result *WorkspaceActivity, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceActivityWithContext(context.Background(), getWorkspaceActivityOptions)
}

// GetWorkspaceActivityWithContext is an alternate form of the GetWorkspaceActivity method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceActivityWithContext(ctx context.Context, getWorkspaceActivityOptions *GetWorkspaceActivityOptions) (result *WorkspaceActivity, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceActivityOptions, "getWorkspaceActivityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceActivityOptions, "getWorkspaceActivityOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id":        *getWorkspaceActivityOptions.WID,
		"activity_id": *getWorkspaceActivityOptions.ActivityID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/actions/{activity_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceActivityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceActivity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivity)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteWorkspaceActivity : Stop the workspace job
// Stop an ongoing schematics job that runs against your workspace.
// **Note**: If you remove the Schematics apply job that runs against your workspace,  any changes to your IBM Cloud
// resources that are already applied are not reverted.  If a creation, update, or deletion is currently in progress,
// Schematics waits for  the job to be completed first. Then, any other resource creations, updates, or  deletions that
// are included in your Terraform template file are ignored.
// <h3>Authorization</h3>  Schematics supports generic authorization such as service access or platform access  to the
// workspace ID and the resource group. For more information, about Schematics  access and permissions, see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) DeleteWorkspaceActivity(deleteWorkspaceActivityOptions *DeleteWorkspaceActivityOptions) (result *WorkspaceActivityApplyResult, response *core.DetailedResponse, err error) {
	return schematics.DeleteWorkspaceActivityWithContext(context.Background(), deleteWorkspaceActivityOptions)
}

// DeleteWorkspaceActivityWithContext is an alternate form of the DeleteWorkspaceActivity method which supports a Context parameter
func (schematics *SchematicsV1) DeleteWorkspaceActivityWithContext(ctx context.Context, deleteWorkspaceActivityOptions *DeleteWorkspaceActivityOptions) (result *WorkspaceActivityApplyResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteWorkspaceActivityOptions, "deleteWorkspaceActivityOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteWorkspaceActivityOptions, "deleteWorkspaceActivityOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id":        *deleteWorkspaceActivityOptions.WID,
		"activity_id": *deleteWorkspaceActivityOptions.ActivityID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/actions/{activity_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteWorkspaceActivityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteWorkspaceActivity")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityApplyResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// RunWorkspaceCommands : Run Terraform Commands
// Run Terraform state commands to modify the workspace state file, by using the IBM Cloud Schematics API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group. For more information, about Schematics
//  access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) RunWorkspaceCommands(runWorkspaceCommandsOptions *RunWorkspaceCommandsOptions) (result *WorkspaceActivityCommandResult, response *core.DetailedResponse, err error) {
	return schematics.RunWorkspaceCommandsWithContext(context.Background(), runWorkspaceCommandsOptions)
}

// RunWorkspaceCommandsWithContext is an alternate form of the RunWorkspaceCommands method which supports a Context parameter
func (schematics *SchematicsV1) RunWorkspaceCommandsWithContext(ctx context.Context, runWorkspaceCommandsOptions *RunWorkspaceCommandsOptions) (result *WorkspaceActivityCommandResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(runWorkspaceCommandsOptions, "runWorkspaceCommandsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(runWorkspaceCommandsOptions, "runWorkspaceCommandsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *runWorkspaceCommandsOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/commands`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range runWorkspaceCommandsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "RunWorkspaceCommands")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if runWorkspaceCommandsOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*runWorkspaceCommandsOptions.RefreshToken))
	}

	body := make(map[string]interface{})
	if runWorkspaceCommandsOptions.Commands != nil {
		body["commands"] = runWorkspaceCommandsOptions.Commands
	}
	if runWorkspaceCommandsOptions.OperationName != nil {
		body["operation_name"] = runWorkspaceCommandsOptions.OperationName
	}
	if runWorkspaceCommandsOptions.Description != nil {
		body["description"] = runWorkspaceCommandsOptions.Description
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityCommandResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ApplyWorkspaceCommand : Perform a Schematics `apply` job
// Run a Schematics `apply` job against your workspace. An `apply` job provisions, modifies, or removes the IBM Cloud
// resources that you described in the Terraform template that your workspace points to. Depending on the type and
// number of resources that you want to provision or modify, this process might take a few minutes, or even up to hours
// to complete. During this time, you cannot make changes to your workspace. After all updates are applied, the state of
// the files is [persisted](https://cloud.ibm.com/docs/schematics?topic=schematics-persist-files) to determine what
// resources exist in your IBM Cloud account.
//
//
//  **Important**: Your workspace must be in an `Inactive`, `Active`, `Failed`, or
//  `Stopped` state to perform a Schematics `apply` job. After all updates are applied,
//  the state of the files is [persisted](https://cloud.ibm.com/docs/schematics?topic=schematics-persist-files)
//  to determine what resources exist in your IBM Cloud account.
//
//
//  **Note**: This API returns an activity or job ID that you use to retrieve the
//  log URL with the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
//
//
//  **Important:** Applying a template might incur costs. Make sure to review
//  the pricing information for the resources that you specified in your
//  templates before you apply the template in IBM Cloud.
//  To find a summary of job that Schematics is about to perform,
//  create a Terraform execution plan with the `POST /v1/workspaces/{id}/plan` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or
//  platform access to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ApplyWorkspaceCommand(applyWorkspaceCommandOptions *ApplyWorkspaceCommandOptions) (result *WorkspaceActivityApplyResult, response *core.DetailedResponse, err error) {
	return schematics.ApplyWorkspaceCommandWithContext(context.Background(), applyWorkspaceCommandOptions)
}

// ApplyWorkspaceCommandWithContext is an alternate form of the ApplyWorkspaceCommand method which supports a Context parameter
func (schematics *SchematicsV1) ApplyWorkspaceCommandWithContext(ctx context.Context, applyWorkspaceCommandOptions *ApplyWorkspaceCommandOptions) (result *WorkspaceActivityApplyResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(applyWorkspaceCommandOptions, "applyWorkspaceCommandOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(applyWorkspaceCommandOptions, "applyWorkspaceCommandOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *applyWorkspaceCommandOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/apply`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range applyWorkspaceCommandOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ApplyWorkspaceCommand")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if applyWorkspaceCommandOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*applyWorkspaceCommandOptions.RefreshToken))
	}
	if applyWorkspaceCommandOptions.DelegatedToken != nil {
		builder.AddHeader("delegated_token", fmt.Sprint(*applyWorkspaceCommandOptions.DelegatedToken))
	}

	body := make(map[string]interface{})
	if applyWorkspaceCommandOptions.ActionOptions != nil {
		body["action_options"] = applyWorkspaceCommandOptions.ActionOptions
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityApplyResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DestroyWorkspaceCommand : Perform a Schematics `destroy` job
// Run a Schematics `destroy` job against your workspace. A `destroy` job removes all IBM Cloud resources that are
// associated with your workspace. Removing your resources does not delete the Schematics workspace. To delete the
// workspace, use the `DELETE /v1/workspaces/{id}` API. This API returns an activity or job ID that you use to retrieve
// the URL to the log file with the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
//
//
//  **Important**: Your workspace must be in an `Active`, `Failed`, or `Stopped` state to perform a Schematics `destroy`
// job.
//
//
//  **Note**: Deleting IBM Cloud resources cannot be undone. Make sure that you back up any required data before you
// remove your resources.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//   to the workspace ID and the resource group.
//   For more information, about Schematics access and permissions,
//   see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) DestroyWorkspaceCommand(destroyWorkspaceCommandOptions *DestroyWorkspaceCommandOptions) (result *WorkspaceActivityDestroyResult, response *core.DetailedResponse, err error) {
	return schematics.DestroyWorkspaceCommandWithContext(context.Background(), destroyWorkspaceCommandOptions)
}

// DestroyWorkspaceCommandWithContext is an alternate form of the DestroyWorkspaceCommand method which supports a Context parameter
func (schematics *SchematicsV1) DestroyWorkspaceCommandWithContext(ctx context.Context, destroyWorkspaceCommandOptions *DestroyWorkspaceCommandOptions) (result *WorkspaceActivityDestroyResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(destroyWorkspaceCommandOptions, "destroyWorkspaceCommandOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(destroyWorkspaceCommandOptions, "destroyWorkspaceCommandOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *destroyWorkspaceCommandOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/destroy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range destroyWorkspaceCommandOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DestroyWorkspaceCommand")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if destroyWorkspaceCommandOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*destroyWorkspaceCommandOptions.RefreshToken))
	}
	if destroyWorkspaceCommandOptions.DelegatedToken != nil {
		builder.AddHeader("delegated_token", fmt.Sprint(*destroyWorkspaceCommandOptions.DelegatedToken))
	}

	body := make(map[string]interface{})
	if destroyWorkspaceCommandOptions.ActionOptions != nil {
		body["action_options"] = destroyWorkspaceCommandOptions.ActionOptions
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityDestroyResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PlanWorkspaceCommand : Perform a Schematics `plan` job
// Run a Schematics `plan` job against your workspace. The `plan` job creates a summary of IBM Cloud resources that must
// be created, modified, or deleted to achieve the state that is described in the Terraform or IBM Cloud catalog
// template that your workspace points to. During this time, you cannot make changes to your workspace. You can use the
// summary to verify your changes before you apply the template in IBM Cloud.
//
//
//  **Important**: Your workspace must be in an `Inactive`, `Active`, `Failed`, or `Stopped` state to perform a
// Schematics `plan` job.
//
//
//  **Note**: This API returns an activity or job ID that you use to retrieve the URL to the log file with the `GET
// /v1/workspaces/{id}/actions/{action_id}/logs` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) PlanWorkspaceCommand(planWorkspaceCommandOptions *PlanWorkspaceCommandOptions) (result *WorkspaceActivityPlanResult, response *core.DetailedResponse, err error) {
	return schematics.PlanWorkspaceCommandWithContext(context.Background(), planWorkspaceCommandOptions)
}

// PlanWorkspaceCommandWithContext is an alternate form of the PlanWorkspaceCommand method which supports a Context parameter
func (schematics *SchematicsV1) PlanWorkspaceCommandWithContext(ctx context.Context, planWorkspaceCommandOptions *PlanWorkspaceCommandOptions) (result *WorkspaceActivityPlanResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(planWorkspaceCommandOptions, "planWorkspaceCommandOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(planWorkspaceCommandOptions, "planWorkspaceCommandOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *planWorkspaceCommandOptions.WID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/plan`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range planWorkspaceCommandOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "PlanWorkspaceCommand")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if planWorkspaceCommandOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*planWorkspaceCommandOptions.RefreshToken))
	}
	if planWorkspaceCommandOptions.DelegatedToken != nil {
		builder.AddHeader("delegated_token", fmt.Sprint(*planWorkspaceCommandOptions.DelegatedToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityPlanResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// RefreshWorkspaceCommand : Perform a Schematics `refresh` job
// Run a Schematics `refresh` job against your workspace. A `refresh` job validates the IBM Cloud resources in your
// account against the state that is stored in the Terraform statefile of your workspace. If differences are found, the
// Terraform statefile is updated accordingly. This API returns an activity or job ID that you use to retrieve the URL
// to the log file with the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
//
//  <h3>Authorization</h3>
//
//  Schematics support generic authorization such as service access or platform access
//  to the workspace ID and the resource group.
//  For more information, about Schematics access and permissions,
//  see [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) RefreshWorkspaceCommand(refreshWorkspaceCommandOptions *RefreshWorkspaceCommandOptions) (result *WorkspaceActivityRefreshResult, response *core.DetailedResponse, err error) {
	return schematics.RefreshWorkspaceCommandWithContext(context.Background(), refreshWorkspaceCommandOptions)
}

// RefreshWorkspaceCommandWithContext is an alternate form of the RefreshWorkspaceCommand method which supports a Context parameter
func (schematics *SchematicsV1) RefreshWorkspaceCommandWithContext(ctx context.Context, refreshWorkspaceCommandOptions *RefreshWorkspaceCommandOptions) (result *WorkspaceActivityRefreshResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(refreshWorkspaceCommandOptions, "refreshWorkspaceCommandOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(refreshWorkspaceCommandOptions, "refreshWorkspaceCommandOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *refreshWorkspaceCommandOptions.WID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/refresh`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range refreshWorkspaceCommandOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "RefreshWorkspaceCommand")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if refreshWorkspaceCommandOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*refreshWorkspaceCommandOptions.RefreshToken))
	}
	if refreshWorkspaceCommandOptions.DelegatedToken != nil {
		builder.AddHeader("delegated_token", fmt.Sprint(*refreshWorkspaceCommandOptions.DelegatedToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityRefreshResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListJobs : List jobs
// Retrieve a list of all Schematics jobs.  The job displays a list of jobs with the status as `pending`, `in_progess`,
// `success`, or `failed`. Jobs are generated when you use the  `POST /v2/jobs`, `PUT /v2/jobs/{job_id}`, or `DELETE
// /v2/jobs/{job_id}`.
//
//  <h3>Authorization</h3>
//  Schematics support generic authorization such as service access or
//  platform access to the job ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ListJobs(listJobsOptions *ListJobsOptions) (result *JobList, response *core.DetailedResponse, err error) {
	return schematics.ListJobsWithContext(context.Background(), listJobsOptions)
}

// ListJobsWithContext is an alternate form of the ListJobs method which supports a Context parameter
func (schematics *SchematicsV1) ListJobsWithContext(ctx context.Context, listJobsOptions *ListJobsOptions) (result *JobList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listJobsOptions, "listJobsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listJobsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListJobs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listJobsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listJobsOptions.Offset))
	}
	if listJobsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listJobsOptions.Limit))
	}
	if listJobsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listJobsOptions.Sort))
	}
	if listJobsOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*listJobsOptions.Profile))
	}
	if listJobsOptions.Resource != nil {
		builder.AddQuery("resource", fmt.Sprint(*listJobsOptions.Resource))
	}
	if listJobsOptions.ResourceID != nil {
		builder.AddQuery("resource_id", fmt.Sprint(*listJobsOptions.ResourceID))
	}
	if listJobsOptions.ActionID != nil {
		builder.AddQuery("action_id", fmt.Sprint(*listJobsOptions.ActionID))
	}
	if listJobsOptions.List != nil {
		builder.AddQuery("list", fmt.Sprint(*listJobsOptions.List))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateJob : Create a job
// Create & launch the Schematics job. It can be used to launch an Ansible playbook against a target hosts.  The job
// displays a list of jobs with the status as `pending`, `in_progess`, `success`, or `failed`.
func (schematics *SchematicsV1) CreateJob(createJobOptions *CreateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	return schematics.CreateJobWithContext(context.Background(), createJobOptions)
}

// CreateJobWithContext is an alternate form of the CreateJob method which supports a Context parameter
func (schematics *SchematicsV1) CreateJobWithContext(ctx context.Context, createJobOptions *CreateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createJobOptions, "createJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createJobOptions, "createJobOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createJobOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*createJobOptions.RefreshToken))
	}

	body := make(map[string]interface{})
	if createJobOptions.CommandObject != nil {
		body["command_object"] = createJobOptions.CommandObject
	}
	if createJobOptions.CommandObjectID != nil {
		body["command_object_id"] = createJobOptions.CommandObjectID
	}
	if createJobOptions.CommandName != nil {
		body["command_name"] = createJobOptions.CommandName
	}
	if createJobOptions.CommandParameter != nil {
		body["command_parameter"] = createJobOptions.CommandParameter
	}
	if createJobOptions.CommandOptions != nil {
		body["command_options"] = createJobOptions.CommandOptions
	}
	if createJobOptions.Inputs != nil {
		body["inputs"] = createJobOptions.Inputs
	}
	if createJobOptions.Settings != nil {
		body["settings"] = createJobOptions.Settings
	}
	if createJobOptions.Tags != nil {
		body["tags"] = createJobOptions.Tags
	}
	if createJobOptions.Location != nil {
		body["location"] = createJobOptions.Location
	}
	if createJobOptions.Status != nil {
		body["status"] = createJobOptions.Status
	}
	if createJobOptions.Data != nil {
		body["data"] = createJobOptions.Data
	}
	if createJobOptions.Bastion != nil {
		body["bastion"] = createJobOptions.Bastion
	}
	if createJobOptions.LogSummary != nil {
		body["log_summary"] = createJobOptions.LogSummary
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetJob : Get a job
// Retrieve the detailed information of Job
//
//  <h3>Authorization</h3>
//  Schematics support generic authorization such as service access or
//  platform access to the job ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetJob(getJobOptions *GetJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	return schematics.GetJobWithContext(context.Background(), getJobOptions)
}

// GetJobWithContext is an alternate form of the GetJob method which supports a Context parameter
func (schematics *SchematicsV1) GetJobWithContext(ctx context.Context, getJobOptions *GetJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getJobOptions, "getJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getJobOptions, "getJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *getJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getJobOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*getJobOptions.Profile))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateJob : Update a job
// Creates a copy of the Schematics job and relaunches an existing job  by updating the information of an existing
// Schematics job.
//
//  <h3>Authorization</h3>
//  Schematics support generic authorization such as service access or
//  platform access to the job ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) UpdateJob(updateJobOptions *UpdateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	return schematics.UpdateJobWithContext(context.Background(), updateJobOptions)
}

// UpdateJobWithContext is an alternate form of the UpdateJob method which supports a Context parameter
func (schematics *SchematicsV1) UpdateJobWithContext(ctx context.Context, updateJobOptions *UpdateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateJobOptions, "updateJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateJobOptions, "updateJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *updateJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UpdateJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateJobOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*updateJobOptions.RefreshToken))
	}

	body := make(map[string]interface{})
	if updateJobOptions.CommandObject != nil {
		body["command_object"] = updateJobOptions.CommandObject
	}
	if updateJobOptions.CommandObjectID != nil {
		body["command_object_id"] = updateJobOptions.CommandObjectID
	}
	if updateJobOptions.CommandName != nil {
		body["command_name"] = updateJobOptions.CommandName
	}
	if updateJobOptions.CommandParameter != nil {
		body["command_parameter"] = updateJobOptions.CommandParameter
	}
	if updateJobOptions.CommandOptions != nil {
		body["command_options"] = updateJobOptions.CommandOptions
	}
	if updateJobOptions.Inputs != nil {
		body["inputs"] = updateJobOptions.Inputs
	}
	if updateJobOptions.Settings != nil {
		body["settings"] = updateJobOptions.Settings
	}
	if updateJobOptions.Tags != nil {
		body["tags"] = updateJobOptions.Tags
	}
	if updateJobOptions.Location != nil {
		body["location"] = updateJobOptions.Location
	}
	if updateJobOptions.Status != nil {
		body["status"] = updateJobOptions.Status
	}
	if updateJobOptions.Data != nil {
		body["data"] = updateJobOptions.Data
	}
	if updateJobOptions.Bastion != nil {
		body["bastion"] = updateJobOptions.Bastion
	}
	if updateJobOptions.LogSummary != nil {
		body["log_summary"] = updateJobOptions.LogSummary
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteJob : Stop the running Job, and delete the Job
// Stop the running Job, and delete the Job.  **Note** You cannot delete or stop the job activity from an ongoing
// execution of an action defined in the playbook.  You can repeat the execution of same job, whenever you patch or
// update the action or workspace.
//
//  <h3>Authorization</h3>
//  Schematics support generic authorization such as service access or
//  platform access to the job ID and the resource group.
//  For more information, about Schematics access and permissions, see
//  [Schematics service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) DeleteJob(deleteJobOptions *DeleteJobOptions) (response *core.DetailedResponse, err error) {
	return schematics.DeleteJobWithContext(context.Background(), deleteJobOptions)
}

// DeleteJobWithContext is an alternate form of the DeleteJob method which supports a Context parameter
func (schematics *SchematicsV1) DeleteJobWithContext(ctx context.Context, deleteJobOptions *DeleteJobOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteJobOptions, "deleteJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteJobOptions, "deleteJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *deleteJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteJobOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*deleteJobOptions.RefreshToken))
	}
	if deleteJobOptions.Force != nil {
		builder.AddHeader("force", fmt.Sprint(*deleteJobOptions.Force))
	}
	if deleteJobOptions.Propagate != nil {
		builder.AddHeader("propagate", fmt.Sprint(*deleteJobOptions.Propagate))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, nil)

	return
}

// ListJobLogs : Get job logs
// Retrieve the job logs
// <h3>Authorization</h3> Schematics support generic authorization such as service access or  platform access to the
// action ID and the resource group.  For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ListJobLogs(listJobLogsOptions *ListJobLogsOptions) (result *JobLog, response *core.DetailedResponse, err error) {
	return schematics.ListJobLogsWithContext(context.Background(), listJobLogsOptions)
}

// ListJobLogsWithContext is an alternate form of the ListJobLogs method which supports a Context parameter
func (schematics *SchematicsV1) ListJobLogsWithContext(ctx context.Context, listJobLogsOptions *ListJobLogsOptions) (result *JobLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listJobLogsOptions, "listJobLogsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listJobLogsOptions, "listJobLogsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *listJobLogsOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}/logs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listJobLogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListJobLogs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateWorkspaceDeletionJob : Delete multiple workspaces
// Delete multiple workspaces.  Use ?destroy_resource="true" to destroy the related cloud resources,  otherwise the
// resources must be managed outside of Schematics.
func (schematics *SchematicsV1) CreateWorkspaceDeletionJob(createWorkspaceDeletionJobOptions *CreateWorkspaceDeletionJobOptions) (result *WorkspaceBulkDeleteResponse, response *core.DetailedResponse, err error) {
	return schematics.CreateWorkspaceDeletionJobWithContext(context.Background(), createWorkspaceDeletionJobOptions)
}

// CreateWorkspaceDeletionJobWithContext is an alternate form of the CreateWorkspaceDeletionJob method which supports a Context parameter
func (schematics *SchematicsV1) CreateWorkspaceDeletionJobWithContext(ctx context.Context, createWorkspaceDeletionJobOptions *CreateWorkspaceDeletionJobOptions) (result *WorkspaceBulkDeleteResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createWorkspaceDeletionJobOptions, "createWorkspaceDeletionJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createWorkspaceDeletionJobOptions, "createWorkspaceDeletionJobOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspace_jobs`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createWorkspaceDeletionJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateWorkspaceDeletionJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createWorkspaceDeletionJobOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*createWorkspaceDeletionJobOptions.RefreshToken))
	}

	if createWorkspaceDeletionJobOptions.DestroyResources != nil {
		builder.AddQuery("destroy_resources", fmt.Sprint(*createWorkspaceDeletionJobOptions.DestroyResources))
	}

	body := make(map[string]interface{})
	if createWorkspaceDeletionJobOptions.NewDeleteWorkspaces != nil {
		body["delete_workspaces"] = createWorkspaceDeletionJobOptions.NewDeleteWorkspaces
	}
	if createWorkspaceDeletionJobOptions.NewDestroyResources != nil {
		body["destroy_resources"] = createWorkspaceDeletionJobOptions.NewDestroyResources
	}
	if createWorkspaceDeletionJobOptions.NewJob != nil {
		body["job"] = createWorkspaceDeletionJobOptions.NewJob
	}
	if createWorkspaceDeletionJobOptions.NewVersion != nil {
		body["version"] = createWorkspaceDeletionJobOptions.NewVersion
	}
	if createWorkspaceDeletionJobOptions.NewWorkspaces != nil {
		body["workspaces"] = createWorkspaceDeletionJobOptions.NewWorkspaces
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceBulkDeleteResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetWorkspaceDeletionJobStatus : Get the workspace deletion job status
// Get the workspace deletion job status.
func (schematics *SchematicsV1) GetWorkspaceDeletionJobStatus(getWorkspaceDeletionJobStatusOptions *GetWorkspaceDeletionJobStatusOptions) (result *WorkspaceJobResponse, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceDeletionJobStatusWithContext(context.Background(), getWorkspaceDeletionJobStatusOptions)
}

// GetWorkspaceDeletionJobStatusWithContext is an alternate form of the GetWorkspaceDeletionJobStatus method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceDeletionJobStatusWithContext(ctx context.Context, getWorkspaceDeletionJobStatusOptions *GetWorkspaceDeletionJobStatusOptions) (result *WorkspaceJobResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWorkspaceDeletionJobStatusOptions, "getWorkspaceDeletionJobStatusOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWorkspaceDeletionJobStatusOptions, "getWorkspaceDeletionJobStatusOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"wj_id": *getWorkspaceDeletionJobStatusOptions.WjID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspace_jobs/{wj_id}/status`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWorkspaceDeletionJobStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetWorkspaceDeletionJobStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceJobResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListInventories : List all inventory definitions
// Retrieve a list of all Schematics inventories that depends on the API endpoint that you have access. For example, if
// you use an API endpoint for a geography, such as North America, only inventories that are created in `us-south` or
// `us-east` are retrieved. For more information, about supported API endpoints, see
// [APIendpoints](/apidocs/schematics#api-endpoints).
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) ListInventories(listInventoriesOptions *ListInventoriesOptions) (result *InventoryResourceRecordList, response *core.DetailedResponse, err error) {
	return schematics.ListInventoriesWithContext(context.Background(), listInventoriesOptions)
}

// ListInventoriesWithContext is an alternate form of the ListInventories method which supports a Context parameter
func (schematics *SchematicsV1) ListInventoriesWithContext(ctx context.Context, listInventoriesOptions *ListInventoriesOptions) (result *InventoryResourceRecordList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listInventoriesOptions, "listInventoriesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listInventoriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListInventories")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listInventoriesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listInventoriesOptions.Offset))
	}
	if listInventoriesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listInventoriesOptions.Limit))
	}
	if listInventoriesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listInventoriesOptions.Sort))
	}
	if listInventoriesOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*listInventoriesOptions.Profile))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInventoryResourceRecordList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateInventory : Create an inventory definition
// Create an IBM Cloud Schematics inventory as a single IBM Cloud resource where you want to run Ansible playbook by
// using Schematics actions. For more information, about inventory host groups, refer to [creating static and dynamic
// inventory for Schematics actions](https://cloud.ibm.com/docs/schematics?topic=schematics-inventories-setup). **Note**
// you cannot update the location and region, resource group once an action is created.  Also, make sure your IP
// addresses are in the [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).  If
// your Git repository already contains a host file. Schematics does not  overwrite the host file already present in
// your Git repository.
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) CreateInventory(createInventoryOptions *CreateInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	return schematics.CreateInventoryWithContext(context.Background(), createInventoryOptions)
}

// CreateInventoryWithContext is an alternate form of the CreateInventory method which supports a Context parameter
func (schematics *SchematicsV1) CreateInventoryWithContext(ctx context.Context, createInventoryOptions *CreateInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createInventoryOptions, "createInventoryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createInventoryOptions, "createInventoryOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createInventoryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateInventory")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createInventoryOptions.Name != nil {
		body["name"] = createInventoryOptions.Name
	}
	if createInventoryOptions.Description != nil {
		body["description"] = createInventoryOptions.Description
	}
	if createInventoryOptions.Location != nil {
		body["location"] = createInventoryOptions.Location
	}
	if createInventoryOptions.ResourceGroup != nil {
		body["resource_group"] = createInventoryOptions.ResourceGroup
	}
	if createInventoryOptions.InventoriesIni != nil {
		body["inventories_ini"] = createInventoryOptions.InventoriesIni
	}
	if createInventoryOptions.ResourceQueries != nil {
		body["resource_queries"] = createInventoryOptions.ResourceQueries
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInventoryResourceRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInventory : Get the inventory definition
// Use this API to retrieve the detailed information for a resource inventory definition used to target an action in
// your IBM Cloud account. For more information, about inventory get, refer to [ibmcloud schematics inventory
// get](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-get-inv).
//  **Note** you can fetch only the location and region, resource group from where your inventory is created.
//  Also, make sure your IP addresses are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).
//
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions, see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) GetInventory(getInventoryOptions *GetInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	return schematics.GetInventoryWithContext(context.Background(), getInventoryOptions)
}

// GetInventoryWithContext is an alternate form of the GetInventory method which supports a Context parameter
func (schematics *SchematicsV1) GetInventoryWithContext(ctx context.Context, getInventoryOptions *GetInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInventoryOptions, "getInventoryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInventoryOptions, "getInventoryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"inventory_id": *getInventoryOptions.InventoryID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories/{inventory_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInventoryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetInventory")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInventoryResourceRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceInventory : Replace an inventory definition
// Use this API to update the inventory definition resource used to target an action. For more information, about
// inventory update, refer to [ibmcloud schematics inventory
// update](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-update-inv).
//  **Note** you cannot update the location and region, resource group once an action is created.
//  Also, make sure your IP addresses are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).
//
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) ReplaceInventory(replaceInventoryOptions *ReplaceInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	return schematics.ReplaceInventoryWithContext(context.Background(), replaceInventoryOptions)
}

// ReplaceInventoryWithContext is an alternate form of the ReplaceInventory method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceInventoryWithContext(ctx context.Context, replaceInventoryOptions *ReplaceInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceInventoryOptions, "replaceInventoryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceInventoryOptions, "replaceInventoryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"inventory_id": *replaceInventoryOptions.InventoryID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories/{inventory_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceInventoryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceInventory")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceInventoryOptions.Name != nil {
		body["name"] = replaceInventoryOptions.Name
	}
	if replaceInventoryOptions.Description != nil {
		body["description"] = replaceInventoryOptions.Description
	}
	if replaceInventoryOptions.Location != nil {
		body["location"] = replaceInventoryOptions.Location
	}
	if replaceInventoryOptions.ResourceGroup != nil {
		body["resource_group"] = replaceInventoryOptions.ResourceGroup
	}
	if replaceInventoryOptions.InventoriesIni != nil {
		body["inventories_ini"] = replaceInventoryOptions.InventoriesIni
	}
	if replaceInventoryOptions.ResourceQueries != nil {
		body["resource_queries"] = replaceInventoryOptions.ResourceQueries
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInventoryResourceRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteInventory : Delete inventory definition
// Use this API to delete the resource inventory definition by using the inventory ID that you want to run against. For
// more information, about inventory delete, refer to [ibmcloud schematics inventory
// delete](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-delete-inventory).
//  **Note** you cannot delete the location and region, resource group from where your inventory is created.
//  Also, make sure your IP addresses are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).
//
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) DeleteInventory(deleteInventoryOptions *DeleteInventoryOptions) (response *core.DetailedResponse, err error) {
	return schematics.DeleteInventoryWithContext(context.Background(), deleteInventoryOptions)
}

// DeleteInventoryWithContext is an alternate form of the DeleteInventory method which supports a Context parameter
func (schematics *SchematicsV1) DeleteInventoryWithContext(ctx context.Context, deleteInventoryOptions *DeleteInventoryOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteInventoryOptions, "deleteInventoryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteInventoryOptions, "deleteInventoryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"inventory_id": *deleteInventoryOptions.InventoryID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories/{inventory_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteInventoryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteInventory")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteInventoryOptions.Force != nil {
		builder.AddHeader("force", fmt.Sprint(*deleteInventoryOptions.Force))
	}
	if deleteInventoryOptions.Propagate != nil {
		builder.AddHeader("propagate", fmt.Sprint(*deleteInventoryOptions.Propagate))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, nil)

	return
}

// UpdateInventory : Update the inventory definition
// Update the resource inventory definition.
func (schematics *SchematicsV1) UpdateInventory(updateInventoryOptions *UpdateInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	return schematics.UpdateInventoryWithContext(context.Background(), updateInventoryOptions)
}

// UpdateInventoryWithContext is an alternate form of the UpdateInventory method which supports a Context parameter
func (schematics *SchematicsV1) UpdateInventoryWithContext(ctx context.Context, updateInventoryOptions *UpdateInventoryOptions) (result *InventoryResourceRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateInventoryOptions, "updateInventoryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateInventoryOptions, "updateInventoryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"inventory_id": *updateInventoryOptions.InventoryID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/inventories/{inventory_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateInventoryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UpdateInventory")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateInventoryOptions.Name != nil {
		body["name"] = updateInventoryOptions.Name
	}
	if updateInventoryOptions.Description != nil {
		body["description"] = updateInventoryOptions.Description
	}
	if updateInventoryOptions.Location != nil {
		body["location"] = updateInventoryOptions.Location
	}
	if updateInventoryOptions.ResourceGroup != nil {
		body["resource_group"] = updateInventoryOptions.ResourceGroup
	}
	if updateInventoryOptions.InventoriesIni != nil {
		body["inventories_ini"] = updateInventoryOptions.InventoriesIni
	}
	if updateInventoryOptions.ResourceQueries != nil {
		body["resource_queries"] = updateInventoryOptions.ResourceQueries
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalInventoryResourceRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListResourceQuery : List resource queries
// Retrieve the list of resource query definitions that you have access to.  The list of resource queries that is
// returned depends on the API  endpoint that you use. For example, if you use an API endpoint for a geography, such as
// North America, only resource query definitions that are created in `us-south` or `us-east` are retrieved. For more
// information, about supported API endpoints, see [API endpoints](/apidocs/schematics#api-endpoints).
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions,  see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) ListResourceQuery(listResourceQueryOptions *ListResourceQueryOptions) (result *ResourceQueryRecordList, response *core.DetailedResponse, err error) {
	return schematics.ListResourceQueryWithContext(context.Background(), listResourceQueryOptions)
}

// ListResourceQueryWithContext is an alternate form of the ListResourceQuery method which supports a Context parameter
func (schematics *SchematicsV1) ListResourceQueryWithContext(ctx context.Context, listResourceQueryOptions *ListResourceQueryOptions) (result *ResourceQueryRecordList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listResourceQueryOptions, "listResourceQueryOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listResourceQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListResourceQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listResourceQueryOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listResourceQueryOptions.Offset))
	}
	if listResourceQueryOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listResourceQueryOptions.Limit))
	}
	if listResourceQueryOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listResourceQueryOptions.Sort))
	}
	if listResourceQueryOptions.Profile != nil {
		builder.AddQuery("profile", fmt.Sprint(*listResourceQueryOptions.Profile))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceQueryRecordList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateResourceQuery : Create resource query
// Use this API to create a resource query definition that will be used to select an IBM Cloud resource or a group of
// resources as the dynamic inventory for the Schematics Actions.  For more information, about resource query commands,
// refer to  [ibmcloud schematics resource query
// create](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-create-rq).
// **Note** you cannot update the location and region, resource group  once an action is created. Also, make sure your
// IP addresses are  in the [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).  If
// your Git repository already contains a host file.  Schematics does not overwrite the host file already present in
// your Git repository.
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions,  see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) CreateResourceQuery(createResourceQueryOptions *CreateResourceQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	return schematics.CreateResourceQueryWithContext(context.Background(), createResourceQueryOptions)
}

// CreateResourceQueryWithContext is an alternate form of the CreateResourceQuery method which supports a Context parameter
func (schematics *SchematicsV1) CreateResourceQueryWithContext(ctx context.Context, createResourceQueryOptions *CreateResourceQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createResourceQueryOptions, "createResourceQueryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createResourceQueryOptions, "createResourceQueryOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createResourceQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateResourceQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createResourceQueryOptions.Type != nil {
		body["type"] = createResourceQueryOptions.Type
	}
	if createResourceQueryOptions.Name != nil {
		body["name"] = createResourceQueryOptions.Name
	}
	if createResourceQueryOptions.Queries != nil {
		body["queries"] = createResourceQueryOptions.Queries
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceQueryRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetResourcesQuery : Get resources query
// Use this API to retrieve the information resource query by Id.  For more information, about resource query commands,
// refer to  [ibmcloud schematics resource query
// get](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-get-rq).
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions,  see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) GetResourcesQuery(getResourcesQueryOptions *GetResourcesQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	return schematics.GetResourcesQueryWithContext(context.Background(), getResourcesQueryOptions)
}

// GetResourcesQueryWithContext is an alternate form of the GetResourcesQuery method which supports a Context parameter
func (schematics *SchematicsV1) GetResourcesQueryWithContext(ctx context.Context, getResourcesQueryOptions *GetResourcesQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourcesQueryOptions, "getResourcesQueryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getResourcesQueryOptions, "getResourcesQueryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"query_id": *getResourcesQueryOptions.QueryID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query/{query_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getResourcesQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetResourcesQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceQueryRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceResourcesQuery : Replace resources query definition
// Use this API to replace the resource query definition used to build  the dynamic inventory for the Schematics Action.
//  For more information, about resource query commands, refer to [ibmcloud schematics resource query
// update](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-update-rq).
// **Note** you cannot update the location and region, resource group  once a resource query is created. Also, make sure
// your IP addresses  are in the
// [allowlist](https://cloud.ibm.com/docs/schematics?topic=schematics-allowed-ipaddresses).
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions,  see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) ReplaceResourcesQuery(replaceResourcesQueryOptions *ReplaceResourcesQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	return schematics.ReplaceResourcesQueryWithContext(context.Background(), replaceResourcesQueryOptions)
}

// ReplaceResourcesQueryWithContext is an alternate form of the ReplaceResourcesQuery method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceResourcesQueryWithContext(ctx context.Context, replaceResourcesQueryOptions *ReplaceResourcesQueryOptions) (result *ResourceQueryRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceResourcesQueryOptions, "replaceResourcesQueryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceResourcesQueryOptions, "replaceResourcesQueryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"query_id": *replaceResourcesQueryOptions.QueryID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query/{query_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceResourcesQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceResourcesQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceResourcesQueryOptions.Type != nil {
		body["type"] = replaceResourcesQueryOptions.Type
	}
	if replaceResourcesQueryOptions.Name != nil {
		body["name"] = replaceResourcesQueryOptions.Name
	}
	if replaceResourcesQueryOptions.Queries != nil {
		body["queries"] = replaceResourcesQueryOptions.Queries
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceQueryRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ExecuteResourceQuery : Run the resource query
// Run the resource query.
func (schematics *SchematicsV1) ExecuteResourceQuery(executeResourceQueryOptions *ExecuteResourceQueryOptions) (result *ResourceQueryResponseRecord, response *core.DetailedResponse, err error) {
	return schematics.ExecuteResourceQueryWithContext(context.Background(), executeResourceQueryOptions)
}

// ExecuteResourceQueryWithContext is an alternate form of the ExecuteResourceQuery method which supports a Context parameter
func (schematics *SchematicsV1) ExecuteResourceQueryWithContext(ctx context.Context, executeResourceQueryOptions *ExecuteResourceQueryOptions) (result *ResourceQueryResponseRecord, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(executeResourceQueryOptions, "executeResourceQueryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(executeResourceQueryOptions, "executeResourceQueryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"query_id": *executeResourceQueryOptions.QueryID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query/{query_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range executeResourceQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ExecuteResourceQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceQueryResponseRecord)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteResourcesQuery : Delete resources query
// Use this API to delete the resource query definition by Id.  For more information, about resource query commands,
// refer to  [ibmcloud schematics resource query
// delete](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-delete-resource-query).
//
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to an
// action ID and the resource group. For more information, about Schematics access and permissions,  see [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#action-permissions).
func (schematics *SchematicsV1) DeleteResourcesQuery(deleteResourcesQueryOptions *DeleteResourcesQueryOptions) (response *core.DetailedResponse, err error) {
	return schematics.DeleteResourcesQueryWithContext(context.Background(), deleteResourcesQueryOptions)
}

// DeleteResourcesQueryWithContext is an alternate form of the DeleteResourcesQuery method which supports a Context parameter
func (schematics *SchematicsV1) DeleteResourcesQueryWithContext(ctx context.Context, deleteResourcesQueryOptions *DeleteResourcesQueryOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourcesQueryOptions, "deleteResourcesQueryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteResourcesQueryOptions, "deleteResourcesQueryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"query_id": *deleteResourcesQueryOptions.QueryID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/resources_query/{query_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteResourcesQueryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteResourcesQuery")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteResourcesQueryOptions.Force != nil {
		builder.AddHeader("force", fmt.Sprint(*deleteResourcesQueryOptions.Force))
	}
	if deleteResourcesQueryOptions.Propagate != nil {
		builder.AddHeader("propagate", fmt.Sprint(*deleteResourcesQueryOptions.Propagate))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = schematics.Service.Request(request, nil)

	return
}

// GetKmsSettings : Get KMS settings
// Retrieve the KMS on the API endpoint that you have access. For example, if you use an API endpoint for a geography,
// such as North America, only Schematics resource that are created in `us-south` or `us-east` are retrieved.
// <h3>Authorization</h3>
//
//
//  Schematics support generic authorization such as service access or platform access to the action ID and the resource
// group. For more information, about Schematics access and permissions, see [Schematics service access roles and
// required permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) GetKmsSettings(getKmsSettingsOptions *GetKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	return schematics.GetKmsSettingsWithContext(context.Background(), getKmsSettingsOptions)
}

// GetKmsSettingsWithContext is an alternate form of the GetKmsSettings method which supports a Context parameter
func (schematics *SchematicsV1) GetKmsSettingsWithContext(ctx context.Context, getKmsSettingsOptions *GetKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKmsSettingsOptions, "getKmsSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getKmsSettingsOptions, "getKmsSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/settings/kms`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getKmsSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetKmsSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("location", fmt.Sprint(*getKmsSettingsOptions.Location))

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateKmsSettings : Replace KMS settings
// Replace or Update the KMS setting for your location, by using your private endpoint, `CRN`, primary `CRK`, and
// secondary `CRK`. **Note** you can update the KMS settings only once. For example, if you use an API endpoint for a
// geography, such as North America, only Schematics resource that are created in `us-south` or `us-east` are retrieved.
//
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to the
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) UpdateKmsSettings(updateKmsSettingsOptions *UpdateKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	return schematics.UpdateKmsSettingsWithContext(context.Background(), updateKmsSettingsOptions)
}

// UpdateKmsSettingsWithContext is an alternate form of the UpdateKmsSettings method which supports a Context parameter
func (schematics *SchematicsV1) UpdateKmsSettingsWithContext(ctx context.Context, updateKmsSettingsOptions *UpdateKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateKmsSettingsOptions, "updateKmsSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateKmsSettingsOptions, "updateKmsSettingsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/settings/kms`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateKmsSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UpdateKmsSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateKmsSettingsOptions.Location != nil {
		body["location"] = updateKmsSettingsOptions.Location
	}
	if updateKmsSettingsOptions.EncryptionScheme != nil {
		body["encryption_scheme"] = updateKmsSettingsOptions.EncryptionScheme
	}
	if updateKmsSettingsOptions.ResourceGroup != nil {
		body["resource_group"] = updateKmsSettingsOptions.ResourceGroup
	}
	if updateKmsSettingsOptions.PrimaryCrk != nil {
		body["primary_crk"] = updateKmsSettingsOptions.PrimaryCrk
	}
	if updateKmsSettingsOptions.SecondaryCrk != nil {
		body["secondary_crk"] = updateKmsSettingsOptions.SecondaryCrk
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
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListKms : List KMS instances
// Lists the KMS instances of your IBM Cloud account to find your Key Protect or Hyper Protect Crypto Services by using
// the location and encrypted scheme such as KYOK or BYOK.
// <h3>Authorization</h3> Schematics support generic authorization such as service access or platform access to the
// action ID and the resource group. For more information, about Schematics access and permissions, see  [Schematics
// service access roles and required
// permissions](https://cloud.ibm.com/docs/schematics?topic=schematics-access#access-roles).
func (schematics *SchematicsV1) ListKms(listKmsOptions *ListKmsOptions) (result *KMSDiscovery, response *core.DetailedResponse, err error) {
	return schematics.ListKmsWithContext(context.Background(), listKmsOptions)
}

// ListKmsWithContext is an alternate form of the ListKms method which supports a Context parameter
func (schematics *SchematicsV1) ListKmsWithContext(ctx context.Context, listKmsOptions *ListKmsOptions) (result *KMSDiscovery, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listKmsOptions, "listKmsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listKmsOptions, "listKmsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/settings/kms_instances`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listKmsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListKms")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("encryption_scheme", fmt.Sprint(*listKmsOptions.EncryptionScheme))
	builder.AddQuery("location", fmt.Sprint(*listKmsOptions.Location))
	if listKmsOptions.ResourceGroup != nil {
		builder.AddQuery("resource_group", fmt.Sprint(*listKmsOptions.ResourceGroup))
	}
	if listKmsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listKmsOptions.Limit))
	}
	if listKmsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listKmsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSDiscovery)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Action : Complete Action details with user inputs and system generated data.
type Action struct {
	// The unique name of your action. The name can be up to 128 characters long and can include alphanumeric characters,
	// spaces, dashes, and underscores. **Example** you can use the name to stop action.
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for an action.  By default, action is created in default resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the `README` file, for the source URL.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Target inventory record ID, used by the action or ansible playbook.
	Inventory *string `json:"inventory,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// User editable variable data & system generated reference to value.
	BastionCredential *VariableData `json:"bastion_credential,omitempty"`

	// Inventory of host and host group for the playbook in `INI` file format. For example, `"targets_ini":
	// "[webserverhost]
	//  172.22.192.6
	//  [dbhost]
	//  172.22.192.5"`. For more information, about an inventory host group syntax, see [Inventory host
	// groups](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).
	TargetsIni *string `json:"targets_ini,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Action ID.
	ID *string `json:"id,omitempty"`

	// Action Cloud Resource Name.
	Crn *string `json:"crn,omitempty"`

	// Action account ID.
	Account *string `json:"account,omitempty"`

	// Action Playbook Source creation time.
	SourceCreatedAt *strfmt.DateTime `json:"source_created_at,omitempty"`

	// E-mail address of user who created the Action Playbook Source.
	SourceCreatedBy *string `json:"source_created_by,omitempty"`

	// The action playbook updation time.
	SourceUpdatedAt *strfmt.DateTime `json:"source_updated_at,omitempty"`

	// E-mail address of user who updated the action playbook source.
	SourceUpdatedBy *string `json:"source_updated_by,omitempty"`

	// Action creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// E-mail address of the user who created an action.
	CreatedBy *string `json:"created_by,omitempty"`

	// Action updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// E-mail address of the user who updated an action.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// Playbook names retrieved from the respository.
	PlaybookNames []string `json:"playbook_names,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`
}

// Constants associated with the Action.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	Action_Location_EuDe    = "eu-de"
	Action_Location_EuGb    = "eu-gb"
	Action_Location_UsEast  = "us-east"
	Action_Location_UsSouth = "us-south"
)

// Constants associated with the Action.SourceType property.
// Type of source for the Template.
const (
	Action_SourceType_CosBucket        = "cos_bucket"
	Action_SourceType_ExternalScm      = "external_scm"
	Action_SourceType_GitHub           = "git_hub"
	Action_SourceType_GitHubEnterprise = "git_hub_enterprise"
	Action_SourceType_GitLab           = "git_lab"
	Action_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	Action_SourceType_IbmGitLab        = "ibm_git_lab"
	Action_SourceType_Local            = "local"
)

// UnmarshalAction unmarshals an instance of Action from the specified map of raw messages.
func UnmarshalAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Action)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "user_state", &obj.UserState, UnmarshalUserState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_readme_url", &obj.SourceReadmeURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalExternalSource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_parameter", &obj.CommandParameter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inventory", &obj.Inventory)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bastion", &obj.Bastion, UnmarshalBastionResourceDefinition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bastion_credential", &obj.BastionCredential, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "targets_ini", &obj.TargetsIni)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_created_at", &obj.SourceCreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_created_by", &obj.SourceCreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_updated_at", &obj.SourceUpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_updated_by", &obj.SourceUpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "state", &obj.State, UnmarshalActionState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "playbook_names", &obj.PlaybookNames)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sys_lock", &obj.SysLock, UnmarshalSystemLock)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionList : List of Action definition response.
type ActionList struct {
	// Total number of records.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// List of action records.
	Actions []ActionLite `json:"actions,omitempty"`
}

// UnmarshalActionList unmarshals an instance of ActionList from the specified map of raw messages.
func UnmarshalActionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionList)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalActionLite)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionLite : Action summary profile with user inputs and system generated data.
type ActionLite struct {
	// Action name (unique for an account).
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// Action Id.
	ID *string `json:"id,omitempty"`

	// Action Cloud Resource Name.
	Crn *string `json:"crn,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Action.  By default, Action will be created in Default Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// name of the namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// Name of the selected playbook.
	PlaybookName *string `json:"playbook_name,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// Computed state of the Action.
	State *ActionLiteState `json:"state,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// Action creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Email address of user who created the action.
	CreatedBy *string `json:"created_by,omitempty"`

	// Action updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Email address of user who updated the action.
	UpdatedBy *string `json:"updated_by,omitempty"`
}

// Constants associated with the ActionLite.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	ActionLite_Location_EuDe    = "eu-de"
	ActionLite_Location_EuGb    = "eu-gb"
	ActionLite_Location_UsEast  = "us-east"
	ActionLite_Location_UsSouth = "us-south"
)

// UnmarshalActionLite unmarshals an instance of ActionLite from the specified map of raw messages.
func UnmarshalActionLite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionLite)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "playbook_name", &obj.PlaybookName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "user_state", &obj.UserState, UnmarshalUserState)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "state", &obj.State, UnmarshalActionLiteState)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sys_lock", &obj.SysLock, UnmarshalSystemLock)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionLiteState : Computed state of the Action.
type ActionLiteState struct {
	// Status of automation (workspace or action).
	StatusCode *string `json:"status_code,omitempty"`

	// Automation status message - to be displayed along with the status_code.
	StatusMessage *string `json:"status_message,omitempty"`
}

// Constants associated with the ActionLiteState.StatusCode property.
// Status of automation (workspace or action).
const (
	ActionLiteState_StatusCode_Critical = "critical"
	ActionLiteState_StatusCode_Disabled = "disabled"
	ActionLiteState_StatusCode_Normal   = "normal"
	ActionLiteState_StatusCode_Pending  = "pending"
)

// UnmarshalActionLiteState unmarshals an instance of ActionLiteState from the specified map of raw messages.
func UnmarshalActionLiteState(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionLiteState)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionState : Computed state of the Action.
type ActionState struct {
	// Status of automation (workspace or action).
	StatusCode *string `json:"status_code,omitempty"`

	// Job id reference for this status.
	StatusJobID *string `json:"status_job_id,omitempty"`

	// Automation status message - to be displayed along with the status_code.
	StatusMessage *string `json:"status_message,omitempty"`
}

// Constants associated with the ActionState.StatusCode property.
// Status of automation (workspace or action).
const (
	ActionState_StatusCode_Critical = "critical"
	ActionState_StatusCode_Disabled = "disabled"
	ActionState_StatusCode_Normal   = "normal"
	ActionState_StatusCode_Pending  = "pending"
)

// UnmarshalActionState unmarshals an instance of ActionState from the specified map of raw messages.
func UnmarshalActionState(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionState)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_job_id", &obj.StatusJobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApplyWorkspaceCommandOptions : The ApplyWorkspaceCommand options.
type ApplyWorkspaceCommandOptions struct {
	// The ID of the workspace for which you want to run a Schematics `apply` job.  To find the workspace ID, use the `GET
	// /workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// Workspace job options template.
	ActionOptions *WorkspaceActivityOptionsTemplate `json:"action_options,omitempty"`

	// The IAM delegated token for your IBM Cloud account.  This token is required for requests that are sent via the UI
	// only.
	DelegatedToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewApplyWorkspaceCommandOptions : Instantiate ApplyWorkspaceCommandOptions
func (*SchematicsV1) NewApplyWorkspaceCommandOptions(wID string, refreshToken string) *ApplyWorkspaceCommandOptions {
	return &ApplyWorkspaceCommandOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (_options *ApplyWorkspaceCommandOptions) SetWID(wID string) *ApplyWorkspaceCommandOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *ApplyWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *ApplyWorkspaceCommandOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetActionOptions : Allow user to set ActionOptions
func (_options *ApplyWorkspaceCommandOptions) SetActionOptions(actionOptions *WorkspaceActivityOptionsTemplate) *ApplyWorkspaceCommandOptions {
	_options.ActionOptions = actionOptions
	return _options
}

// SetDelegatedToken : Allow user to set DelegatedToken
func (_options *ApplyWorkspaceCommandOptions) SetDelegatedToken(delegatedToken string) *ApplyWorkspaceCommandOptions {
	_options.DelegatedToken = core.StringPtr(delegatedToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ApplyWorkspaceCommandOptions) SetHeaders(param map[string]string) *ApplyWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// BastionResourceDefinition : Describes a bastion resource.
type BastionResourceDefinition struct {
	// Bastion Name(Unique).
	Name *string `json:"name,omitempty"`

	// Reference to the Inventory resource definition.
	Host *string `json:"host,omitempty"`
}

// UnmarshalBastionResourceDefinition unmarshals an instance of BastionResourceDefinition from the specified map of raw messages.
func UnmarshalBastionResourceDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BastionResourceDefinition)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "host", &obj.Host)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CatalogRef : Information about the software template that you chose from the IBM Cloud catalog. This information is returned for
// IBM Cloud catalog offerings only.
type CatalogRef struct {
	// Dry run.
	DryRun *bool `json:"dry_run,omitempty"`

	// Owning account ID of the catalog.
	OwningAccount *string `json:"owning_account,omitempty"`

	// The URL to the icon of the software template in the IBM Cloud catalog.
	ItemIconURL *string `json:"item_icon_url,omitempty"`

	// The ID of the software template that you chose to install from the IBM Cloud catalog. This software is provisioned
	// with Schematics.
	ItemID *string `json:"item_id,omitempty"`

	// The name of the software that you chose to install from the IBM Cloud catalog.
	ItemName *string `json:"item_name,omitempty"`

	// The URL to the readme file of the software template in the IBM Cloud catalog.
	ItemReadmeURL *string `json:"item_readme_url,omitempty"`

	// The URL to the software template in the IBM Cloud catalog.
	ItemURL *string `json:"item_url,omitempty"`

	// The URL to the dashboard to access your software.
	LaunchURL *string `json:"launch_url,omitempty"`

	// The version of the software template that you chose to install from the IBM Cloud catalog.
	OfferingVersion *string `json:"offering_version,omitempty"`
}

// UnmarshalCatalogRef unmarshals an instance of CatalogRef from the specified map of raw messages.
func UnmarshalCatalogRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogRef)
	err = core.UnmarshalPrimitive(m, "dry_run", &obj.DryRun)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "owning_account", &obj.OwningAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "item_icon_url", &obj.ItemIconURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "item_id", &obj.ItemID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "item_name", &obj.ItemName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "item_readme_url", &obj.ItemReadmeURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "item_url", &obj.ItemURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "launch_url", &obj.LaunchURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_version", &obj.OfferingVersion)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateActionOptions : The CreateAction options.
type CreateActionOptions struct {
	// The unique name of your action. The name can be up to 128 characters long and can include alphanumeric characters,
	// spaces, dashes, and underscores. **Example** you can use the name to stop action.
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for an action.  By default, action is created in default resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the `README` file, for the source URL.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Target inventory record ID, used by the action or ansible playbook.
	Inventory *string `json:"inventory,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// User editable variable data & system generated reference to value.
	BastionCredential *VariableData `json:"bastion_credential,omitempty"`

	// Inventory of host and host group for the playbook in `INI` file format. For example, `"targets_ini":
	// "[webserverhost]
	//  172.22.192.6
	//  [dbhost]
	//  172.22.192.5"`. For more information, about an inventory host group syntax, see [Inventory host
	// groups](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).
	TargetsIni *string `json:"targets_ini,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform
	// template.
	XGithubToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateActionOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	CreateActionOptions_Location_EuDe    = "eu-de"
	CreateActionOptions_Location_EuGb    = "eu-gb"
	CreateActionOptions_Location_UsEast  = "us-east"
	CreateActionOptions_Location_UsSouth = "us-south"
)

// Constants associated with the CreateActionOptions.SourceType property.
// Type of source for the Template.
const (
	CreateActionOptions_SourceType_CosBucket        = "cos_bucket"
	CreateActionOptions_SourceType_ExternalScm      = "external_scm"
	CreateActionOptions_SourceType_GitHub           = "git_hub"
	CreateActionOptions_SourceType_GitHubEnterprise = "git_hub_enterprise"
	CreateActionOptions_SourceType_GitLab           = "git_lab"
	CreateActionOptions_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	CreateActionOptions_SourceType_IbmGitLab        = "ibm_git_lab"
	CreateActionOptions_SourceType_Local            = "local"
)

// NewCreateActionOptions : Instantiate CreateActionOptions
func (*SchematicsV1) NewCreateActionOptions() *CreateActionOptions {
	return &CreateActionOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateActionOptions) SetName(name string) *CreateActionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateActionOptions) SetDescription(description string) *CreateActionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateActionOptions) SetLocation(location string) *CreateActionOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateActionOptions) SetResourceGroup(resourceGroup string) *CreateActionOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateActionOptions) SetTags(tags []string) *CreateActionOptions {
	_options.Tags = tags
	return _options
}

// SetUserState : Allow user to set UserState
func (_options *CreateActionOptions) SetUserState(userState *UserState) *CreateActionOptions {
	_options.UserState = userState
	return _options
}

// SetSourceReadmeURL : Allow user to set SourceReadmeURL
func (_options *CreateActionOptions) SetSourceReadmeURL(sourceReadmeURL string) *CreateActionOptions {
	_options.SourceReadmeURL = core.StringPtr(sourceReadmeURL)
	return _options
}

// SetSource : Allow user to set Source
func (_options *CreateActionOptions) SetSource(source *ExternalSource) *CreateActionOptions {
	_options.Source = source
	return _options
}

// SetSourceType : Allow user to set SourceType
func (_options *CreateActionOptions) SetSourceType(sourceType string) *CreateActionOptions {
	_options.SourceType = core.StringPtr(sourceType)
	return _options
}

// SetCommandParameter : Allow user to set CommandParameter
func (_options *CreateActionOptions) SetCommandParameter(commandParameter string) *CreateActionOptions {
	_options.CommandParameter = core.StringPtr(commandParameter)
	return _options
}

// SetInventory : Allow user to set Inventory
func (_options *CreateActionOptions) SetInventory(inventory string) *CreateActionOptions {
	_options.Inventory = core.StringPtr(inventory)
	return _options
}

// SetCredentials : Allow user to set Credentials
func (_options *CreateActionOptions) SetCredentials(credentials []VariableData) *CreateActionOptions {
	_options.Credentials = credentials
	return _options
}

// SetBastion : Allow user to set Bastion
func (_options *CreateActionOptions) SetBastion(bastion *BastionResourceDefinition) *CreateActionOptions {
	_options.Bastion = bastion
	return _options
}

// SetBastionCredential : Allow user to set BastionCredential
func (_options *CreateActionOptions) SetBastionCredential(bastionCredential *VariableData) *CreateActionOptions {
	_options.BastionCredential = bastionCredential
	return _options
}

// SetTargetsIni : Allow user to set TargetsIni
func (_options *CreateActionOptions) SetTargetsIni(targetsIni string) *CreateActionOptions {
	_options.TargetsIni = core.StringPtr(targetsIni)
	return _options
}

// SetInputs : Allow user to set Inputs
func (_options *CreateActionOptions) SetInputs(inputs []VariableData) *CreateActionOptions {
	_options.Inputs = inputs
	return _options
}

// SetOutputs : Allow user to set Outputs
func (_options *CreateActionOptions) SetOutputs(outputs []VariableData) *CreateActionOptions {
	_options.Outputs = outputs
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *CreateActionOptions) SetSettings(settings []VariableData) *CreateActionOptions {
	_options.Settings = settings
	return _options
}

// SetState : Allow user to set State
func (_options *CreateActionOptions) SetState(state *ActionState) *CreateActionOptions {
	_options.State = state
	return _options
}

// SetSysLock : Allow user to set SysLock
func (_options *CreateActionOptions) SetSysLock(sysLock *SystemLock) *CreateActionOptions {
	_options.SysLock = sysLock
	return _options
}

// SetXGithubToken : Allow user to set XGithubToken
func (_options *CreateActionOptions) SetXGithubToken(xGithubToken string) *CreateActionOptions {
	_options.XGithubToken = core.StringPtr(xGithubToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateActionOptions) SetHeaders(param map[string]string) *CreateActionOptions {
	options.Headers = param
	return options
}

// CreateInventoryOptions : The CreateInventory options.
type CreateInventoryOptions struct {
	// The unique name of your Inventory definition. The name can be up to 128 characters long and can include alphanumeric
	// characters, spaces, dashes, and underscores.
	Name *string `json:"name,omitempty"`

	// The description of your Inventory definition. The description can be up to 2048 characters long in size.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Inventory definition.   By default, Inventory definition will be created in Default
	// Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Input inventory of host and host group for the playbook, in the `.ini` file format.
	InventoriesIni *string `json:"inventories_ini,omitempty"`

	// Input resource query definitions that is used to dynamically generate the inventory of host and host group for the
	// playbook.
	ResourceQueries []string `json:"resource_queries,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateInventoryOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	CreateInventoryOptions_Location_EuDe    = "eu-de"
	CreateInventoryOptions_Location_EuGb    = "eu-gb"
	CreateInventoryOptions_Location_UsEast  = "us-east"
	CreateInventoryOptions_Location_UsSouth = "us-south"
)

// NewCreateInventoryOptions : Instantiate CreateInventoryOptions
func (*SchematicsV1) NewCreateInventoryOptions() *CreateInventoryOptions {
	return &CreateInventoryOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateInventoryOptions) SetName(name string) *CreateInventoryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateInventoryOptions) SetDescription(description string) *CreateInventoryOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateInventoryOptions) SetLocation(location string) *CreateInventoryOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateInventoryOptions) SetResourceGroup(resourceGroup string) *CreateInventoryOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetInventoriesIni : Allow user to set InventoriesIni
func (_options *CreateInventoryOptions) SetInventoriesIni(inventoriesIni string) *CreateInventoryOptions {
	_options.InventoriesIni = core.StringPtr(inventoriesIni)
	return _options
}

// SetResourceQueries : Allow user to set ResourceQueries
func (_options *CreateInventoryOptions) SetResourceQueries(resourceQueries []string) *CreateInventoryOptions {
	_options.ResourceQueries = resourceQueries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateInventoryOptions) SetHeaders(param map[string]string) *CreateInventoryOptions {
	options.Headers = param
	return options
}

// CreateJobOptions : The CreateJob options.
type CreateJobOptions struct {
	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action or Workspace.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action or Workspace.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job data.
	Data *JobData `json:"data,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateJobOptions.CommandObject property.
// Name of the Schematics automation resource.
const (
	CreateJobOptions_CommandObject_Action      = "action"
	CreateJobOptions_CommandObject_Environment = "environment"
	CreateJobOptions_CommandObject_System      = "system"
	CreateJobOptions_CommandObject_Workspace   = "workspace"
)

// Constants associated with the CreateJobOptions.CommandName property.
// Schematics job command name.
const (
	CreateJobOptions_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	CreateJobOptions_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	CreateJobOptions_CommandName_CreateAction         = "create_action"
	CreateJobOptions_CommandName_CreateCart           = "create_cart"
	CreateJobOptions_CommandName_CreateEnvironment    = "create_environment"
	CreateJobOptions_CommandName_CreateWorkspace      = "create_workspace"
	CreateJobOptions_CommandName_DeleteAction         = "delete_action"
	CreateJobOptions_CommandName_DeleteEnvironment    = "delete_environment"
	CreateJobOptions_CommandName_DeleteWorkspace      = "delete_workspace"
	CreateJobOptions_CommandName_EnvironmentInit      = "environment_init"
	CreateJobOptions_CommandName_EnvironmentInstall   = "environment_install"
	CreateJobOptions_CommandName_EnvironmentUninstall = "environment_uninstall"
	CreateJobOptions_CommandName_PatchAction          = "patch_action"
	CreateJobOptions_CommandName_PatchWorkspace       = "patch_workspace"
	CreateJobOptions_CommandName_PutAction            = "put_action"
	CreateJobOptions_CommandName_PutEnvironment       = "put_environment"
	CreateJobOptions_CommandName_PutWorkspace         = "put_workspace"
	CreateJobOptions_CommandName_RepositoryProcess    = "repository_process"
	CreateJobOptions_CommandName_SystemKeyDelete      = "system_key_delete"
	CreateJobOptions_CommandName_SystemKeyDisable     = "system_key_disable"
	CreateJobOptions_CommandName_SystemKeyEnable      = "system_key_enable"
	CreateJobOptions_CommandName_SystemKeyRestore     = "system_key_restore"
	CreateJobOptions_CommandName_SystemKeyRotate      = "system_key_rotate"
	CreateJobOptions_CommandName_WorkspaceApply       = "workspace_apply"
	CreateJobOptions_CommandName_WorkspaceDestroy     = "workspace_destroy"
	CreateJobOptions_CommandName_WorkspacePlan        = "workspace_plan"
	CreateJobOptions_CommandName_WorkspaceRefresh     = "workspace_refresh"
)

// Constants associated with the CreateJobOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	CreateJobOptions_Location_EuDe    = "eu-de"
	CreateJobOptions_Location_EuGb    = "eu-gb"
	CreateJobOptions_Location_UsEast  = "us-east"
	CreateJobOptions_Location_UsSouth = "us-south"
)

// NewCreateJobOptions : Instantiate CreateJobOptions
func (*SchematicsV1) NewCreateJobOptions(refreshToken string) *CreateJobOptions {
	return &CreateJobOptions{
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *CreateJobOptions) SetRefreshToken(refreshToken string) *CreateJobOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetCommandObject : Allow user to set CommandObject
func (_options *CreateJobOptions) SetCommandObject(commandObject string) *CreateJobOptions {
	_options.CommandObject = core.StringPtr(commandObject)
	return _options
}

// SetCommandObjectID : Allow user to set CommandObjectID
func (_options *CreateJobOptions) SetCommandObjectID(commandObjectID string) *CreateJobOptions {
	_options.CommandObjectID = core.StringPtr(commandObjectID)
	return _options
}

// SetCommandName : Allow user to set CommandName
func (_options *CreateJobOptions) SetCommandName(commandName string) *CreateJobOptions {
	_options.CommandName = core.StringPtr(commandName)
	return _options
}

// SetCommandParameter : Allow user to set CommandParameter
func (_options *CreateJobOptions) SetCommandParameter(commandParameter string) *CreateJobOptions {
	_options.CommandParameter = core.StringPtr(commandParameter)
	return _options
}

// SetCommandOptions : Allow user to set CommandOptions
func (_options *CreateJobOptions) SetCommandOptions(commandOptions []string) *CreateJobOptions {
	_options.CommandOptions = commandOptions
	return _options
}

// SetInputs : Allow user to set Inputs
func (_options *CreateJobOptions) SetInputs(inputs []VariableData) *CreateJobOptions {
	_options.Inputs = inputs
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *CreateJobOptions) SetSettings(settings []VariableData) *CreateJobOptions {
	_options.Settings = settings
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateJobOptions) SetTags(tags []string) *CreateJobOptions {
	_options.Tags = tags
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateJobOptions) SetLocation(location string) *CreateJobOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *CreateJobOptions) SetStatus(status *JobStatus) *CreateJobOptions {
	_options.Status = status
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateJobOptions) SetData(data *JobData) *CreateJobOptions {
	_options.Data = data
	return _options
}

// SetBastion : Allow user to set Bastion
func (_options *CreateJobOptions) SetBastion(bastion *BastionResourceDefinition) *CreateJobOptions {
	_options.Bastion = bastion
	return _options
}

// SetLogSummary : Allow user to set LogSummary
func (_options *CreateJobOptions) SetLogSummary(logSummary *JobLogSummary) *CreateJobOptions {
	_options.LogSummary = logSummary
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateJobOptions) SetHeaders(param map[string]string) *CreateJobOptions {
	options.Headers = param
	return options
}

// CreateResourceQueryOptions : The CreateResourceQuery options.
type CreateResourceQueryOptions struct {
	// Resource type (cluster, vsi, icd, vpc).
	Type *string `json:"type,omitempty"`

	// Resource query name.
	Name *string `json:"name,omitempty"`

	Queries []ResourceQuery `json:"queries,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateResourceQueryOptions.Type property.
// Resource type (cluster, vsi, icd, vpc).
const (
	CreateResourceQueryOptions_Type_Vsi = "vsi"
)

// NewCreateResourceQueryOptions : Instantiate CreateResourceQueryOptions
func (*SchematicsV1) NewCreateResourceQueryOptions() *CreateResourceQueryOptions {
	return &CreateResourceQueryOptions{}
}

// SetType : Allow user to set Type
func (_options *CreateResourceQueryOptions) SetType(typeVar string) *CreateResourceQueryOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateResourceQueryOptions) SetName(name string) *CreateResourceQueryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetQueries : Allow user to set Queries
func (_options *CreateResourceQueryOptions) SetQueries(queries []ResourceQuery) *CreateResourceQueryOptions {
	_options.Queries = queries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateResourceQueryOptions) SetHeaders(param map[string]string) *CreateResourceQueryOptions {
	options.Headers = param
	return options
}

// CreateWorkspaceDeletionJobOptions : The CreateWorkspaceDeletionJob options.
type CreateWorkspaceDeletionJobOptions struct {
	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// True to delete workspace.
	NewDeleteWorkspaces *bool `json:"delete_workspaces,omitempty"`

	// True to destroy the resources managed by this workspace.
	NewDestroyResources *bool `json:"destroy_resources,omitempty"`

	// Workspace deletion job name.
	NewJob *string `json:"job,omitempty"`

	// Version of the terraform template.
	NewVersion *string `json:"version,omitempty"`

	// List of workspaces to be deleted.
	NewWorkspaces []string `json:"workspaces,omitempty"`

	// If set to `true`, refresh_token header configuration is required to delete all the Terraform resources, and the
	// Schematics workspace. If set to `false`, you can remove only the workspace. Your Terraform resources are still
	// available and must be managed with the resource dashboard or CLI.
	DestroyResources *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateWorkspaceDeletionJobOptions : Instantiate CreateWorkspaceDeletionJobOptions
func (*SchematicsV1) NewCreateWorkspaceDeletionJobOptions(refreshToken string) *CreateWorkspaceDeletionJobOptions {
	return &CreateWorkspaceDeletionJobOptions{
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *CreateWorkspaceDeletionJobOptions) SetRefreshToken(refreshToken string) *CreateWorkspaceDeletionJobOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetNewDeleteWorkspaces : Allow user to set NewDeleteWorkspaces
func (_options *CreateWorkspaceDeletionJobOptions) SetNewDeleteWorkspaces(newDeleteWorkspaces bool) *CreateWorkspaceDeletionJobOptions {
	_options.NewDeleteWorkspaces = core.BoolPtr(newDeleteWorkspaces)
	return _options
}

// SetNewDestroyResources : Allow user to set NewDestroyResources
func (_options *CreateWorkspaceDeletionJobOptions) SetNewDestroyResources(newDestroyResources bool) *CreateWorkspaceDeletionJobOptions {
	_options.NewDestroyResources = core.BoolPtr(newDestroyResources)
	return _options
}

// SetNewJob : Allow user to set NewJob
func (_options *CreateWorkspaceDeletionJobOptions) SetNewJob(newJob string) *CreateWorkspaceDeletionJobOptions {
	_options.NewJob = core.StringPtr(newJob)
	return _options
}

// SetNewVersion : Allow user to set NewVersion
func (_options *CreateWorkspaceDeletionJobOptions) SetNewVersion(newVersion string) *CreateWorkspaceDeletionJobOptions {
	_options.NewVersion = core.StringPtr(newVersion)
	return _options
}

// SetNewWorkspaces : Allow user to set NewWorkspaces
func (_options *CreateWorkspaceDeletionJobOptions) SetNewWorkspaces(newWorkspaces []string) *CreateWorkspaceDeletionJobOptions {
	_options.NewWorkspaces = newWorkspaces
	return _options
}

// SetDestroyResources : Allow user to set DestroyResources
func (_options *CreateWorkspaceDeletionJobOptions) SetDestroyResources(destroyResources string) *CreateWorkspaceDeletionJobOptions {
	_options.DestroyResources = core.StringPtr(destroyResources)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateWorkspaceDeletionJobOptions) SetHeaders(param map[string]string) *CreateWorkspaceDeletionJobOptions {
	options.Headers = param
	return options
}

// CreateWorkspaceOptions : The CreateWorkspace options.
type CreateWorkspaceOptions struct {
	// List of applied shared dataset ID.
	AppliedShareddataIds []string `json:"applied_shareddata_ids,omitempty"`

	// Information about the software template that you chose from the IBM Cloud catalog. This information is returned for
	// IBM Cloud catalog offerings only.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// The description of the workspace.
	Description *string `json:"description,omitempty"`

	// The location where you want to create your Schematics workspace and run the Schematics jobs. The location that you
	// enter must match the API endpoint that you use. For example, if you use the Frankfurt API endpoint, you must specify
	// `eu-de` as your location. If you use an API endpoint for a geography and you do not specify a location, Schematics
	// determines the location based on availability.
	Location *string `json:"location,omitempty"`

	// The name of your workspace. The name can be up to 128 characters long and can include alphanumeric characters,
	// spaces, dashes, and underscores. When you create a workspace for your own Terraform template, consider including the
	// microservice component that you set up with your Terraform template and the IBM Cloud environment where you want to
	// deploy your resources in your name.
	Name *string `json:"name,omitempty"`

	// The ID of the resource group where you want to provision the workspace.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Information about the Target used by the templates originating from the  IBM Cloud catalog offerings. This
	// information is not relevant for workspace created using your own Terraform template.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// A list of tags that are associated with the workspace.
	Tags []string `json:"tags,omitempty"`

	// Input data for the Template.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// Workspace template ref.
	TemplateRef *string `json:"template_ref,omitempty"`

	// Input variables for the Template repoository, while creating a workspace.
	TemplateRepo *TemplateRepoRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// WorkspaceStatusRequest -.
	WorkspaceStatus *WorkspaceStatusRequest `json:"workspace_status,omitempty"`

	// The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform
	// template.
	XGithubToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateWorkspaceOptions : Instantiate CreateWorkspaceOptions
func (*SchematicsV1) NewCreateWorkspaceOptions() *CreateWorkspaceOptions {
	return &CreateWorkspaceOptions{}
}

// SetAppliedShareddataIds : Allow user to set AppliedShareddataIds
func (_options *CreateWorkspaceOptions) SetAppliedShareddataIds(appliedShareddataIds []string) *CreateWorkspaceOptions {
	_options.AppliedShareddataIds = appliedShareddataIds
	return _options
}

// SetCatalogRef : Allow user to set CatalogRef
func (_options *CreateWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *CreateWorkspaceOptions {
	_options.CatalogRef = catalogRef
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateWorkspaceOptions) SetDescription(description string) *CreateWorkspaceOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateWorkspaceOptions) SetLocation(location string) *CreateWorkspaceOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateWorkspaceOptions) SetName(name string) *CreateWorkspaceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateWorkspaceOptions) SetResourceGroup(resourceGroup string) *CreateWorkspaceOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetSharedData : Allow user to set SharedData
func (_options *CreateWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *CreateWorkspaceOptions {
	_options.SharedData = sharedData
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateWorkspaceOptions) SetTags(tags []string) *CreateWorkspaceOptions {
	_options.Tags = tags
	return _options
}

// SetTemplateData : Allow user to set TemplateData
func (_options *CreateWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *CreateWorkspaceOptions {
	_options.TemplateData = templateData
	return _options
}

// SetTemplateRef : Allow user to set TemplateRef
func (_options *CreateWorkspaceOptions) SetTemplateRef(templateRef string) *CreateWorkspaceOptions {
	_options.TemplateRef = core.StringPtr(templateRef)
	return _options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (_options *CreateWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoRequest) *CreateWorkspaceOptions {
	_options.TemplateRepo = templateRepo
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateWorkspaceOptions) SetType(typeVar []string) *CreateWorkspaceOptions {
	_options.Type = typeVar
	return _options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (_options *CreateWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusRequest) *CreateWorkspaceOptions {
	_options.WorkspaceStatus = workspaceStatus
	return _options
}

// SetXGithubToken : Allow user to set XGithubToken
func (_options *CreateWorkspaceOptions) SetXGithubToken(xGithubToken string) *CreateWorkspaceOptions {
	_options.XGithubToken = core.StringPtr(xGithubToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateWorkspaceOptions) SetHeaders(param map[string]string) *CreateWorkspaceOptions {
	options.Headers = param
	return options
}

// DeleteActionOptions : The DeleteAction options.
type DeleteActionOptions struct {
	// Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.
	ActionID *string `json:"-" validate:"required,ne="`

	// Equivalent to -force options in the command line.
	Force *bool `json:"-"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteActionOptions : Instantiate DeleteActionOptions
func (*SchematicsV1) NewDeleteActionOptions(actionID string) *DeleteActionOptions {
	return &DeleteActionOptions{
		ActionID: core.StringPtr(actionID),
	}
}

// SetActionID : Allow user to set ActionID
func (_options *DeleteActionOptions) SetActionID(actionID string) *DeleteActionOptions {
	_options.ActionID = core.StringPtr(actionID)
	return _options
}

// SetForce : Allow user to set Force
func (_options *DeleteActionOptions) SetForce(force bool) *DeleteActionOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetPropagate : Allow user to set Propagate
func (_options *DeleteActionOptions) SetPropagate(propagate bool) *DeleteActionOptions {
	_options.Propagate = core.BoolPtr(propagate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteActionOptions) SetHeaders(param map[string]string) *DeleteActionOptions {
	options.Headers = param
	return options
}

// DeleteInventoryOptions : The DeleteInventory options.
type DeleteInventoryOptions struct {
	// Resource Inventory Id.  Use `GET /v2/inventories` API to look up the Resource Inventory definition Ids  in your IBM
	// Cloud account.
	InventoryID *string `json:"-" validate:"required,ne="`

	// Equivalent to -force options in the command line.
	Force *bool `json:"-"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteInventoryOptions : Instantiate DeleteInventoryOptions
func (*SchematicsV1) NewDeleteInventoryOptions(inventoryID string) *DeleteInventoryOptions {
	return &DeleteInventoryOptions{
		InventoryID: core.StringPtr(inventoryID),
	}
}

// SetInventoryID : Allow user to set InventoryID
func (_options *DeleteInventoryOptions) SetInventoryID(inventoryID string) *DeleteInventoryOptions {
	_options.InventoryID = core.StringPtr(inventoryID)
	return _options
}

// SetForce : Allow user to set Force
func (_options *DeleteInventoryOptions) SetForce(force bool) *DeleteInventoryOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetPropagate : Allow user to set Propagate
func (_options *DeleteInventoryOptions) SetPropagate(propagate bool) *DeleteInventoryOptions {
	_options.Propagate = core.BoolPtr(propagate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteInventoryOptions) SetHeaders(param map[string]string) *DeleteInventoryOptions {
	options.Headers = param
	return options
}

// DeleteJobOptions : The DeleteJob options.
type DeleteJobOptions struct {
	// Job Id. Use `GET /v2/jobs` API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// Equivalent to -force options in the command line.
	Force *bool `json:"-"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteJobOptions : Instantiate DeleteJobOptions
func (*SchematicsV1) NewDeleteJobOptions(jobID string, refreshToken string) *DeleteJobOptions {
	return &DeleteJobOptions{
		JobID:        core.StringPtr(jobID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetJobID : Allow user to set JobID
func (_options *DeleteJobOptions) SetJobID(jobID string) *DeleteJobOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *DeleteJobOptions) SetRefreshToken(refreshToken string) *DeleteJobOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetForce : Allow user to set Force
func (_options *DeleteJobOptions) SetForce(force bool) *DeleteJobOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetPropagate : Allow user to set Propagate
func (_options *DeleteJobOptions) SetPropagate(propagate bool) *DeleteJobOptions {
	_options.Propagate = core.BoolPtr(propagate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteJobOptions) SetHeaders(param map[string]string) *DeleteJobOptions {
	options.Headers = param
	return options
}

// DeleteResourcesQueryOptions : The DeleteResourcesQuery options.
type DeleteResourcesQueryOptions struct {
	// Resource query Id.  Use `GET /v2/resource_query` API to look up the Resource query definition Ids  in your IBM Cloud
	// account.
	QueryID *string `json:"-" validate:"required,ne="`

	// Equivalent to -force options in the command line.
	Force *bool `json:"-"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteResourcesQueryOptions : Instantiate DeleteResourcesQueryOptions
func (*SchematicsV1) NewDeleteResourcesQueryOptions(queryID string) *DeleteResourcesQueryOptions {
	return &DeleteResourcesQueryOptions{
		QueryID: core.StringPtr(queryID),
	}
}

// SetQueryID : Allow user to set QueryID
func (_options *DeleteResourcesQueryOptions) SetQueryID(queryID string) *DeleteResourcesQueryOptions {
	_options.QueryID = core.StringPtr(queryID)
	return _options
}

// SetForce : Allow user to set Force
func (_options *DeleteResourcesQueryOptions) SetForce(force bool) *DeleteResourcesQueryOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetPropagate : Allow user to set Propagate
func (_options *DeleteResourcesQueryOptions) SetPropagate(propagate bool) *DeleteResourcesQueryOptions {
	_options.Propagate = core.BoolPtr(propagate)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteResourcesQueryOptions) SetHeaders(param map[string]string) *DeleteResourcesQueryOptions {
	options.Headers = param
	return options
}

// DeleteWorkspaceActivityOptions : The DeleteWorkspaceActivity options.
type DeleteWorkspaceActivityOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the activity or job, for which you want to retrieve details.  To find the job ID, use the `GET
	// /v1/workspaces/{id}/actions` API.
	ActivityID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteWorkspaceActivityOptions : Instantiate DeleteWorkspaceActivityOptions
func (*SchematicsV1) NewDeleteWorkspaceActivityOptions(wID string, activityID string) *DeleteWorkspaceActivityOptions {
	return &DeleteWorkspaceActivityOptions{
		WID:        core.StringPtr(wID),
		ActivityID: core.StringPtr(activityID),
	}
}

// SetWID : Allow user to set WID
func (_options *DeleteWorkspaceActivityOptions) SetWID(wID string) *DeleteWorkspaceActivityOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetActivityID : Allow user to set ActivityID
func (_options *DeleteWorkspaceActivityOptions) SetActivityID(activityID string) *DeleteWorkspaceActivityOptions {
	_options.ActivityID = core.StringPtr(activityID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWorkspaceActivityOptions) SetHeaders(param map[string]string) *DeleteWorkspaceActivityOptions {
	options.Headers = param
	return options
}

// DeleteWorkspaceOptions : The DeleteWorkspace options.
type DeleteWorkspaceOptions struct {
	// The IAM refresh token for the user or service identity. The IAM refresh token is required only if you want to
	// destroy the Terraform resources before deleting the Schematics workspace. If you want to delete the workspace only
	// and keep all your Terraform resources, refresh token is not required.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// If set to `true`, refresh_token header configuration is required to delete all the Terraform resources, and the
	// Schematics workspace. If set to `false`, you can remove only the workspace. Your Terraform resources are still
	// available and must be managed with the resource dashboard or CLI.
	DestroyResources *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteWorkspaceOptions : Instantiate DeleteWorkspaceOptions
func (*SchematicsV1) NewDeleteWorkspaceOptions(refreshToken string, wID string) *DeleteWorkspaceOptions {
	return &DeleteWorkspaceOptions{
		RefreshToken: core.StringPtr(refreshToken),
		WID:          core.StringPtr(wID),
	}
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *DeleteWorkspaceOptions) SetRefreshToken(refreshToken string) *DeleteWorkspaceOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetWID : Allow user to set WID
func (_options *DeleteWorkspaceOptions) SetWID(wID string) *DeleteWorkspaceOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetDestroyResources : Allow user to set DestroyResources
func (_options *DeleteWorkspaceOptions) SetDestroyResources(destroyResources string) *DeleteWorkspaceOptions {
	_options.DestroyResources = core.StringPtr(destroyResources)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWorkspaceOptions) SetHeaders(param map[string]string) *DeleteWorkspaceOptions {
	options.Headers = param
	return options
}

// DestroyWorkspaceCommandOptions : The DestroyWorkspaceCommand options.
type DestroyWorkspaceCommandOptions struct {
	// The ID of the workspace for which you want to perform a Schematics `destroy` job.  To find the workspace ID, use the
	// `GET /workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// Workspace job options template.
	ActionOptions *WorkspaceActivityOptionsTemplate `json:"action_options,omitempty"`

	// The IAM delegated token for your IBM Cloud account.  This token is required for requests that are sent via the UI
	// only.
	DelegatedToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDestroyWorkspaceCommandOptions : Instantiate DestroyWorkspaceCommandOptions
func (*SchematicsV1) NewDestroyWorkspaceCommandOptions(wID string, refreshToken string) *DestroyWorkspaceCommandOptions {
	return &DestroyWorkspaceCommandOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (_options *DestroyWorkspaceCommandOptions) SetWID(wID string) *DestroyWorkspaceCommandOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *DestroyWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *DestroyWorkspaceCommandOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetActionOptions : Allow user to set ActionOptions
func (_options *DestroyWorkspaceCommandOptions) SetActionOptions(actionOptions *WorkspaceActivityOptionsTemplate) *DestroyWorkspaceCommandOptions {
	_options.ActionOptions = actionOptions
	return _options
}

// SetDelegatedToken : Allow user to set DelegatedToken
func (_options *DestroyWorkspaceCommandOptions) SetDelegatedToken(delegatedToken string) *DestroyWorkspaceCommandOptions {
	_options.DelegatedToken = core.StringPtr(delegatedToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DestroyWorkspaceCommandOptions) SetHeaders(param map[string]string) *DestroyWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// EnvVariableResponse : List of environment values.
type EnvVariableResponse struct {
	// Environment variable is hidden.
	Hidden *bool `json:"hidden,omitempty"`

	// Environment variable name.
	Name *string `json:"name,omitempty"`

	// Environment variable is secure.
	Secure *bool `json:"secure,omitempty"`

	// Value for environment variable.
	Value *string `json:"value,omitempty"`
}

// UnmarshalEnvVariableResponse unmarshals an instance of EnvVariableResponse from the specified map of raw messages.
func UnmarshalEnvVariableResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvVariableResponse)
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExecuteResourceQueryOptions : The ExecuteResourceQuery options.
type ExecuteResourceQueryOptions struct {
	// Resource query Id.  Use `GET /v2/resource_query` API to look up the Resource query definition Ids  in your IBM Cloud
	// account.
	QueryID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewExecuteResourceQueryOptions : Instantiate ExecuteResourceQueryOptions
func (*SchematicsV1) NewExecuteResourceQueryOptions(queryID string) *ExecuteResourceQueryOptions {
	return &ExecuteResourceQueryOptions{
		QueryID: core.StringPtr(queryID),
	}
}

// SetQueryID : Allow user to set QueryID
func (_options *ExecuteResourceQueryOptions) SetQueryID(queryID string) *ExecuteResourceQueryOptions {
	_options.QueryID = core.StringPtr(queryID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ExecuteResourceQueryOptions) SetHeaders(param map[string]string) *ExecuteResourceQueryOptions {
	options.Headers = param
	return options
}

// ExternalSource : Source of templates, playbooks, or controls.
type ExternalSource struct {
	// Type of source for the Template.
	SourceType *string `json:"source_type" validate:"required"`

	// Connection details to Git source.
	Git *ExternalSourceGit `json:"git,omitempty"`

	// Connection details to IBM Cloud Catalog source.
	Catalog *ExternalSourceCatalog `json:"catalog,omitempty"`

	// Connection details to a IBM Cloud Object Storage bucket.
	CosBucket *ExternalSourceCosBucket `json:"cos_bucket,omitempty"`
}

// Constants associated with the ExternalSource.SourceType property.
// Type of source for the Template.
const (
	ExternalSource_SourceType_CosBucket        = "cos_bucket"
	ExternalSource_SourceType_ExternalScm      = "external_scm"
	ExternalSource_SourceType_GitHub           = "git_hub"
	ExternalSource_SourceType_GitHubEnterprise = "git_hub_enterprise"
	ExternalSource_SourceType_GitLab           = "git_lab"
	ExternalSource_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	ExternalSource_SourceType_IbmGitLab        = "ibm_git_lab"
	ExternalSource_SourceType_Local            = "local"
)

// NewExternalSource : Instantiate ExternalSource (Generic Model Constructor)
func (*SchematicsV1) NewExternalSource(sourceType string) (_model *ExternalSource, err error) {
	_model = &ExternalSource{
		SourceType: core.StringPtr(sourceType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalExternalSource unmarshals an instance of ExternalSource from the specified map of raw messages.
func UnmarshalExternalSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalSource)
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "git", &obj.Git, UnmarshalExternalSourceGit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "catalog", &obj.Catalog, UnmarshalExternalSourceCatalog)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cos_bucket", &obj.CosBucket, UnmarshalExternalSourceCosBucket)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalSourceCatalog : Connection details to IBM Cloud Catalog source.
type ExternalSourceCatalog struct {
	// name of the private catalog.
	CatalogName *string `json:"catalog_name,omitempty"`

	// Name of the offering in the IBM Catalog.
	OfferingName *string `json:"offering_name,omitempty"`

	// Version string of the offering in the IBM Catalog.
	OfferingVersion *string `json:"offering_version,omitempty"`

	// Type of the offering, in the IBM Catalog.
	OfferingKind *string `json:"offering_kind,omitempty"`

	// Id of the offering the IBM Catalog.
	OfferingID *string `json:"offering_id,omitempty"`

	// Id of the offering version the IBM Catalog.
	OfferingVersionID *string `json:"offering_version_id,omitempty"`

	// Repo Url of the offering, in the IBM Catalog.
	OfferingRepoURL *string `json:"offering_repo_url,omitempty"`
}

// UnmarshalExternalSourceCatalog unmarshals an instance of ExternalSourceCatalog from the specified map of raw messages.
func UnmarshalExternalSourceCatalog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalSourceCatalog)
	err = core.UnmarshalPrimitive(m, "catalog_name", &obj.CatalogName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_name", &obj.OfferingName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_version", &obj.OfferingVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_kind", &obj.OfferingKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_id", &obj.OfferingID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_version_id", &obj.OfferingVersionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offering_repo_url", &obj.OfferingRepoURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalSourceCosBucket : Connection details to a IBM Cloud Object Storage bucket.
type ExternalSourceCosBucket struct {
	// COS Bucket Url.
	CosBucketURL *string `json:"cos_bucket_url,omitempty"`
}

// UnmarshalExternalSourceCosBucket unmarshals an instance of ExternalSourceCosBucket from the specified map of raw messages.
func UnmarshalExternalSourceCosBucket(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalSourceCosBucket)
	err = core.UnmarshalPrimitive(m, "cos_bucket_url", &obj.CosBucketURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalSourceGit : Connection details to Git source.
type ExternalSourceGit struct {
	// The Complete URL which is computed by git_repo_url, git_repo_folder and branch.
	ComputedGitRepoURL *string `json:"computed_git_repo_url,omitempty"`

	// URL to the GIT Repo that can be used to clone the template.
	GitRepoURL *string `json:"git_repo_url,omitempty"`

	// Personal Access Token to connect to Git URLs.
	GitToken *string `json:"git_token,omitempty"`

	// Name of the folder in the Git Repo, that contains the template.
	GitRepoFolder *string `json:"git_repo_folder,omitempty"`

	// Name of the release tag, used to fetch the Git Repo.
	GitRelease *string `json:"git_release,omitempty"`

	// Name of the branch, used to fetch the Git Repo.
	GitBranch *string `json:"git_branch,omitempty"`
}

// UnmarshalExternalSourceGit unmarshals an instance of ExternalSourceGit from the specified map of raw messages.
func UnmarshalExternalSourceGit(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalSourceGit)
	err = core.UnmarshalPrimitive(m, "computed_git_repo_url", &obj.ComputedGitRepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "git_repo_url", &obj.GitRepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "git_token", &obj.GitToken)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "git_repo_folder", &obj.GitRepoFolder)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "git_release", &obj.GitRelease)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "git_branch", &obj.GitBranch)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetActionOptions : The GetAction options.
type GetActionOptions struct {
	// Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.
	ActionID *string `json:"-" validate:"required,ne="`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetActionOptions.Profile property.
// Level of details returned by the get method.
const (
	GetActionOptions_Profile_Detailed = "detailed"
	GetActionOptions_Profile_Summary  = "summary"
)

// NewGetActionOptions : Instantiate GetActionOptions
func (*SchematicsV1) NewGetActionOptions(actionID string) *GetActionOptions {
	return &GetActionOptions{
		ActionID: core.StringPtr(actionID),
	}
}

// SetActionID : Allow user to set ActionID
func (_options *GetActionOptions) SetActionID(actionID string) *GetActionOptions {
	_options.ActionID = core.StringPtr(actionID)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *GetActionOptions) SetProfile(profile string) *GetActionOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetActionOptions) SetHeaders(param map[string]string) *GetActionOptions {
	options.Headers = param
	return options
}

// GetAllWorkspaceInputsOptions : The GetAllWorkspaceInputs options.
type GetAllWorkspaceInputsOptions struct {
	// The ID of the workspace for which you want to retrieve input parameters and  values. To find the workspace ID, use
	// the `GET /workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAllWorkspaceInputsOptions : Instantiate GetAllWorkspaceInputsOptions
func (*SchematicsV1) NewGetAllWorkspaceInputsOptions(wID string) *GetAllWorkspaceInputsOptions {
	return &GetAllWorkspaceInputsOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetAllWorkspaceInputsOptions) SetWID(wID string) *GetAllWorkspaceInputsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllWorkspaceInputsOptions) SetHeaders(param map[string]string) *GetAllWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// GetInventoryOptions : The GetInventory options.
type GetInventoryOptions struct {
	// Resource Inventory Id.  Use `GET /v2/inventories` API to look up the Resource Inventory definition Ids  in your IBM
	// Cloud account.
	InventoryID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInventoryOptions : Instantiate GetInventoryOptions
func (*SchematicsV1) NewGetInventoryOptions(inventoryID string) *GetInventoryOptions {
	return &GetInventoryOptions{
		InventoryID: core.StringPtr(inventoryID),
	}
}

// SetInventoryID : Allow user to set InventoryID
func (_options *GetInventoryOptions) SetInventoryID(inventoryID string) *GetInventoryOptions {
	_options.InventoryID = core.StringPtr(inventoryID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInventoryOptions) SetHeaders(param map[string]string) *GetInventoryOptions {
	options.Headers = param
	return options
}

// GetJobOptions : The GetJob options.
type GetJobOptions struct {
	// Job Id. Use `GET /v2/jobs` API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"-" validate:"required,ne="`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetJobOptions.Profile property.
// Level of details returned by the get method.
const (
	GetJobOptions_Profile_Detailed = "detailed"
	GetJobOptions_Profile_Summary  = "summary"
)

// NewGetJobOptions : Instantiate GetJobOptions
func (*SchematicsV1) NewGetJobOptions(jobID string) *GetJobOptions {
	return &GetJobOptions{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *GetJobOptions) SetJobID(jobID string) *GetJobOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *GetJobOptions) SetProfile(profile string) *GetJobOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetJobOptions) SetHeaders(param map[string]string) *GetJobOptions {
	options.Headers = param
	return options
}

// GetKmsSettingsOptions : The GetKmsSettings options.
type GetKmsSettingsOptions struct {
	// The location of the Resource.
	Location *string `json:"-" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetKmsSettingsOptions : Instantiate GetKmsSettingsOptions
func (*SchematicsV1) NewGetKmsSettingsOptions(location string) *GetKmsSettingsOptions {
	return &GetKmsSettingsOptions{
		Location: core.StringPtr(location),
	}
}

// SetLocation : Allow user to set Location
func (_options *GetKmsSettingsOptions) SetLocation(location string) *GetKmsSettingsOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKmsSettingsOptions) SetHeaders(param map[string]string) *GetKmsSettingsOptions {
	options.Headers = param
	return options
}

// GetResourcesQueryOptions : The GetResourcesQuery options.
type GetResourcesQueryOptions struct {
	// Resource query Id.  Use `GET /v2/resource_query` API to look up the Resource query definition Ids  in your IBM Cloud
	// account.
	QueryID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetResourcesQueryOptions : Instantiate GetResourcesQueryOptions
func (*SchematicsV1) NewGetResourcesQueryOptions(queryID string) *GetResourcesQueryOptions {
	return &GetResourcesQueryOptions{
		QueryID: core.StringPtr(queryID),
	}
}

// SetQueryID : Allow user to set QueryID
func (_options *GetResourcesQueryOptions) SetQueryID(queryID string) *GetResourcesQueryOptions {
	_options.QueryID = core.StringPtr(queryID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourcesQueryOptions) SetHeaders(param map[string]string) *GetResourcesQueryOptions {
	options.Headers = param
	return options
}

// GetSchematicsVersionOptions : The GetSchematicsVersion options.
type GetSchematicsVersionOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSchematicsVersionOptions : Instantiate GetSchematicsVersionOptions
func (*SchematicsV1) NewGetSchematicsVersionOptions() *GetSchematicsVersionOptions {
	return &GetSchematicsVersionOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSchematicsVersionOptions) SetHeaders(param map[string]string) *GetSchematicsVersionOptions {
	options.Headers = param
	return options
}

// GetTemplateActivityLogOptions : The GetTemplateActivityLog options.
type GetTemplateActivityLogOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template or IBM Cloud catalog software template in the workspace.  Use the `GET
	// /v1/workspaces` to look up the workspace IDs and template IDs or `template_data.id`.
	TID *string `json:"-" validate:"required,ne="`

	// The ID of the activity or job, for which you want to retrieve details.  To find the job ID, use the `GET
	// /v1/workspaces/{id}/actions` API.
	ActivityID *string `json:"-" validate:"required,ne="`

	// Enter false to replace the first line in each Terraform command section, such as Terraform INIT or Terraform PLAN,
	// with Schematics INIT (Schematics PLAN) in your log output.  In addition, the log lines Starting command: terraform
	// init -input=false -no-color and Starting command: terraform apply -state=terraform.tfstate
	// -var-file=schematics.tfvars -auto-approve -no-color are suppressed.  All subsequent command lines still use the
	// Terraform command prefix. To remove this prefix, use the log_tf_prefix option.
	LogTfCmd *bool `json:"-"`

	// `false` will hide all the terraform command prefix in the log statements.
	LogTfPrefix *bool `json:"-"`

	// `false` will hide all the null resource prefix in the log statements.
	LogTfNullResource *bool `json:"-"`

	// `true` will format all logs to withhold the original format  of ansible output in the log statements.
	LogTfAnsible *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTemplateActivityLogOptions : Instantiate GetTemplateActivityLogOptions
func (*SchematicsV1) NewGetTemplateActivityLogOptions(wID string, tID string, activityID string) *GetTemplateActivityLogOptions {
	return &GetTemplateActivityLogOptions{
		WID:        core.StringPtr(wID),
		TID:        core.StringPtr(tID),
		ActivityID: core.StringPtr(activityID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetTemplateActivityLogOptions) SetWID(wID string) *GetTemplateActivityLogOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *GetTemplateActivityLogOptions) SetTID(tID string) *GetTemplateActivityLogOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetActivityID : Allow user to set ActivityID
func (_options *GetTemplateActivityLogOptions) SetActivityID(activityID string) *GetTemplateActivityLogOptions {
	_options.ActivityID = core.StringPtr(activityID)
	return _options
}

// SetLogTfCmd : Allow user to set LogTfCmd
func (_options *GetTemplateActivityLogOptions) SetLogTfCmd(logTfCmd bool) *GetTemplateActivityLogOptions {
	_options.LogTfCmd = core.BoolPtr(logTfCmd)
	return _options
}

// SetLogTfPrefix : Allow user to set LogTfPrefix
func (_options *GetTemplateActivityLogOptions) SetLogTfPrefix(logTfPrefix bool) *GetTemplateActivityLogOptions {
	_options.LogTfPrefix = core.BoolPtr(logTfPrefix)
	return _options
}

// SetLogTfNullResource : Allow user to set LogTfNullResource
func (_options *GetTemplateActivityLogOptions) SetLogTfNullResource(logTfNullResource bool) *GetTemplateActivityLogOptions {
	_options.LogTfNullResource = core.BoolPtr(logTfNullResource)
	return _options
}

// SetLogTfAnsible : Allow user to set LogTfAnsible
func (_options *GetTemplateActivityLogOptions) SetLogTfAnsible(logTfAnsible bool) *GetTemplateActivityLogOptions {
	_options.LogTfAnsible = core.BoolPtr(logTfAnsible)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateActivityLogOptions) SetHeaders(param map[string]string) *GetTemplateActivityLogOptions {
	options.Headers = param
	return options
}

// GetTemplateLogsOptions : The GetTemplateLogs options.
type GetTemplateLogsOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template or IBM Cloud catalog software template in the workspace.  Use the `GET
	// /v1/workspaces` to look up the workspace IDs and template IDs or `template_data.id`.
	TID *string `json:"-" validate:"required,ne="`

	// Enter false to replace the first line in each Terraform command section, such as Terraform INIT or Terraform PLAN,
	// with Schematics INIT (Schematics PLAN) in your log output.  In addition, the log lines Starting command: terraform
	// init -input=false -no-color and Starting command: terraform apply -state=terraform.tfstate
	// -var-file=schematics.tfvars -auto-approve -no-color are suppressed.  All subsequent command lines still use the
	// Terraform command prefix. To remove this prefix, use the log_tf_prefix option.
	LogTfCmd *bool `json:"-"`

	// `false` will hide all the terraform command prefix in the log statements.
	LogTfPrefix *bool `json:"-"`

	// `false` will hide all the null resource prefix in the log statements.
	LogTfNullResource *bool `json:"-"`

	// `true` will format all logs to withhold the original format  of ansible output in the log statements.
	LogTfAnsible *bool `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTemplateLogsOptions : Instantiate GetTemplateLogsOptions
func (*SchematicsV1) NewGetTemplateLogsOptions(wID string, tID string) *GetTemplateLogsOptions {
	return &GetTemplateLogsOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetTemplateLogsOptions) SetWID(wID string) *GetTemplateLogsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *GetTemplateLogsOptions) SetTID(tID string) *GetTemplateLogsOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetLogTfCmd : Allow user to set LogTfCmd
func (_options *GetTemplateLogsOptions) SetLogTfCmd(logTfCmd bool) *GetTemplateLogsOptions {
	_options.LogTfCmd = core.BoolPtr(logTfCmd)
	return _options
}

// SetLogTfPrefix : Allow user to set LogTfPrefix
func (_options *GetTemplateLogsOptions) SetLogTfPrefix(logTfPrefix bool) *GetTemplateLogsOptions {
	_options.LogTfPrefix = core.BoolPtr(logTfPrefix)
	return _options
}

// SetLogTfNullResource : Allow user to set LogTfNullResource
func (_options *GetTemplateLogsOptions) SetLogTfNullResource(logTfNullResource bool) *GetTemplateLogsOptions {
	_options.LogTfNullResource = core.BoolPtr(logTfNullResource)
	return _options
}

// SetLogTfAnsible : Allow user to set LogTfAnsible
func (_options *GetTemplateLogsOptions) SetLogTfAnsible(logTfAnsible bool) *GetTemplateLogsOptions {
	_options.LogTfAnsible = core.BoolPtr(logTfAnsible)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateLogsOptions) SetHeaders(param map[string]string) *GetTemplateLogsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceActivityLogsOptions : The GetWorkspaceActivityLogs options.
type GetWorkspaceActivityLogsOptions struct {
	// The ID of the workspace for which you want to retrieve the Terraform statefile.  To find the workspace ID, use the
	// `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the activity or job, for which you want to retrieve details.  To find the job ID, use the `GET
	// /v1/workspaces/{id}/actions` API.
	ActivityID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceActivityLogsOptions : Instantiate GetWorkspaceActivityLogsOptions
func (*SchematicsV1) NewGetWorkspaceActivityLogsOptions(wID string, activityID string) *GetWorkspaceActivityLogsOptions {
	return &GetWorkspaceActivityLogsOptions{
		WID:        core.StringPtr(wID),
		ActivityID: core.StringPtr(activityID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceActivityLogsOptions) SetWID(wID string) *GetWorkspaceActivityLogsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetActivityID : Allow user to set ActivityID
func (_options *GetWorkspaceActivityLogsOptions) SetActivityID(activityID string) *GetWorkspaceActivityLogsOptions {
	_options.ActivityID = core.StringPtr(activityID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceActivityLogsOptions) SetHeaders(param map[string]string) *GetWorkspaceActivityLogsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceActivityOptions : The GetWorkspaceActivity options.
type GetWorkspaceActivityOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the activity or job, for which you want to retrieve details.  To find the job ID, use the `GET
	// /v1/workspaces/{id}/actions` API.
	ActivityID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceActivityOptions : Instantiate GetWorkspaceActivityOptions
func (*SchematicsV1) NewGetWorkspaceActivityOptions(wID string, activityID string) *GetWorkspaceActivityOptions {
	return &GetWorkspaceActivityOptions{
		WID:        core.StringPtr(wID),
		ActivityID: core.StringPtr(activityID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceActivityOptions) SetWID(wID string) *GetWorkspaceActivityOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetActivityID : Allow user to set ActivityID
func (_options *GetWorkspaceActivityOptions) SetActivityID(activityID string) *GetWorkspaceActivityOptions {
	_options.ActivityID = core.StringPtr(activityID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceActivityOptions) SetHeaders(param map[string]string) *GetWorkspaceActivityOptions {
	options.Headers = param
	return options
}

// GetWorkspaceDeletionJobStatusOptions : The GetWorkspaceDeletionJobStatus options.
type GetWorkspaceDeletionJobStatusOptions struct {
	// The workspace job ID.
	WjID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceDeletionJobStatusOptions : Instantiate GetWorkspaceDeletionJobStatusOptions
func (*SchematicsV1) NewGetWorkspaceDeletionJobStatusOptions(wjID string) *GetWorkspaceDeletionJobStatusOptions {
	return &GetWorkspaceDeletionJobStatusOptions{
		WjID: core.StringPtr(wjID),
	}
}

// SetWjID : Allow user to set WjID
func (_options *GetWorkspaceDeletionJobStatusOptions) SetWjID(wjID string) *GetWorkspaceDeletionJobStatusOptions {
	_options.WjID = core.StringPtr(wjID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceDeletionJobStatusOptions) SetHeaders(param map[string]string) *GetWorkspaceDeletionJobStatusOptions {
	options.Headers = param
	return options
}

// GetWorkspaceInputMetadataOptions : The GetWorkspaceInputMetadata options.
type GetWorkspaceInputMetadataOptions struct {
	// The ID of the workspace for which you want to retrieve the metadata of the input variables that are declared in the
	// template. To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template for which you want to retrieve the metadata of your input variables. When you
	// create a workspace, the Terraform template that your workspace points to is assigned a unique ID. To find this ID,
	// use the `GET /v1/workspaces` API and review the `template_data.id` value.
	TID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceInputMetadataOptions : Instantiate GetWorkspaceInputMetadataOptions
func (*SchematicsV1) NewGetWorkspaceInputMetadataOptions(wID string, tID string) *GetWorkspaceInputMetadataOptions {
	return &GetWorkspaceInputMetadataOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceInputMetadataOptions) SetWID(wID string) *GetWorkspaceInputMetadataOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *GetWorkspaceInputMetadataOptions) SetTID(tID string) *GetWorkspaceInputMetadataOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceInputMetadataOptions) SetHeaders(param map[string]string) *GetWorkspaceInputMetadataOptions {
	options.Headers = param
	return options
}

// GetWorkspaceInputsOptions : The GetWorkspaceInputs options.
type GetWorkspaceInputsOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template in your workspace.  When you create a workspace, the Terraform template that  your
	// workspace points to is assigned a unique ID. Use the `GET /v1/workspaces` to look up the workspace IDs  and template
	// IDs or `template_data.id` in your IBM Cloud account.
	TID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceInputsOptions : Instantiate GetWorkspaceInputsOptions
func (*SchematicsV1) NewGetWorkspaceInputsOptions(wID string, tID string) *GetWorkspaceInputsOptions {
	return &GetWorkspaceInputsOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceInputsOptions) SetWID(wID string) *GetWorkspaceInputsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *GetWorkspaceInputsOptions) SetTID(tID string) *GetWorkspaceInputsOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceInputsOptions) SetHeaders(param map[string]string) *GetWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceLogUrlsOptions : The GetWorkspaceLogUrls options.
type GetWorkspaceLogUrlsOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceLogUrlsOptions : Instantiate GetWorkspaceLogUrlsOptions
func (*SchematicsV1) NewGetWorkspaceLogUrlsOptions(wID string) *GetWorkspaceLogUrlsOptions {
	return &GetWorkspaceLogUrlsOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceLogUrlsOptions) SetWID(wID string) *GetWorkspaceLogUrlsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceLogUrlsOptions) SetHeaders(param map[string]string) *GetWorkspaceLogUrlsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceOptions : The GetWorkspace options.
type GetWorkspaceOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceOptions : Instantiate GetWorkspaceOptions
func (*SchematicsV1) NewGetWorkspaceOptions(wID string) *GetWorkspaceOptions {
	return &GetWorkspaceOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceOptions) SetWID(wID string) *GetWorkspaceOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceOptions) SetHeaders(param map[string]string) *GetWorkspaceOptions {
	options.Headers = param
	return options
}

// GetWorkspaceOutputsOptions : The GetWorkspaceOutputs options.
type GetWorkspaceOutputsOptions struct {
	// The ID of the workspace for which you want to retrieve output parameters and  values. To find the workspace ID, use
	// the `GET /workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceOutputsOptions : Instantiate GetWorkspaceOutputsOptions
func (*SchematicsV1) NewGetWorkspaceOutputsOptions(wID string) *GetWorkspaceOutputsOptions {
	return &GetWorkspaceOutputsOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceOutputsOptions) SetWID(wID string) *GetWorkspaceOutputsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceOutputsOptions) SetHeaders(param map[string]string) *GetWorkspaceOutputsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceReadmeOptions : The GetWorkspaceReadme options.
type GetWorkspaceReadmeOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The GitHub or GitLab branch where the `README.md` file is stored,  or the commit ID or tag that references the
	// `README.md` file that you want to retrieve.  If you do not specify this option, the `README.md` file is retrieved
	// from the master branch by default.
	Ref *string `json:"-"`

	// The format of the readme file.  Value ''markdown'' will give markdown, otherwise html.
	Formatted *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetWorkspaceReadmeOptions.Formatted property.
// The format of the readme file.  Value ''markdown'' will give markdown, otherwise html.
const (
	GetWorkspaceReadmeOptions_Formatted_HTML     = "html"
	GetWorkspaceReadmeOptions_Formatted_Markdown = "markdown"
)

// NewGetWorkspaceReadmeOptions : Instantiate GetWorkspaceReadmeOptions
func (*SchematicsV1) NewGetWorkspaceReadmeOptions(wID string) *GetWorkspaceReadmeOptions {
	return &GetWorkspaceReadmeOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceReadmeOptions) SetWID(wID string) *GetWorkspaceReadmeOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRef : Allow user to set Ref
func (_options *GetWorkspaceReadmeOptions) SetRef(ref string) *GetWorkspaceReadmeOptions {
	_options.Ref = core.StringPtr(ref)
	return _options
}

// SetFormatted : Allow user to set Formatted
func (_options *GetWorkspaceReadmeOptions) SetFormatted(formatted string) *GetWorkspaceReadmeOptions {
	_options.Formatted = core.StringPtr(formatted)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceReadmeOptions) SetHeaders(param map[string]string) *GetWorkspaceReadmeOptions {
	options.Headers = param
	return options
}

// GetWorkspaceResourcesOptions : The GetWorkspaceResources options.
type GetWorkspaceResourcesOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceResourcesOptions : Instantiate GetWorkspaceResourcesOptions
func (*SchematicsV1) NewGetWorkspaceResourcesOptions(wID string) *GetWorkspaceResourcesOptions {
	return &GetWorkspaceResourcesOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceResourcesOptions) SetWID(wID string) *GetWorkspaceResourcesOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceResourcesOptions) SetHeaders(param map[string]string) *GetWorkspaceResourcesOptions {
	options.Headers = param
	return options
}

// GetWorkspaceStateOptions : The GetWorkspaceState options.
type GetWorkspaceStateOptions struct {
	// The ID of the workspace for which you want to retrieve the Terraform statefile.  To find the workspace ID, use the
	// `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceStateOptions : Instantiate GetWorkspaceStateOptions
func (*SchematicsV1) NewGetWorkspaceStateOptions(wID string) *GetWorkspaceStateOptions {
	return &GetWorkspaceStateOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceStateOptions) SetWID(wID string) *GetWorkspaceStateOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceStateOptions) SetHeaders(param map[string]string) *GetWorkspaceStateOptions {
	options.Headers = param
	return options
}

// GetWorkspaceTemplateStateOptions : The GetWorkspaceTemplateState options.
type GetWorkspaceTemplateStateOptions struct {
	// The ID of the workspace for which you want to retrieve the Terraform statefile.  To find the workspace ID, use the
	// `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template for which you want to retrieve the Terraform statefile.  When you create a
	// workspace, the Terraform template that your workspace points to is assigned a unique ID.  To find this ID, use the
	// `GET /v1/workspaces` API and review the template_data.id value.
	TID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWorkspaceTemplateStateOptions : Instantiate GetWorkspaceTemplateStateOptions
func (*SchematicsV1) NewGetWorkspaceTemplateStateOptions(wID string, tID string) *GetWorkspaceTemplateStateOptions {
	return &GetWorkspaceTemplateStateOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *GetWorkspaceTemplateStateOptions) SetWID(wID string) *GetWorkspaceTemplateStateOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *GetWorkspaceTemplateStateOptions) SetTID(tID string) *GetWorkspaceTemplateStateOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceTemplateStateOptions) SetHeaders(param map[string]string) *GetWorkspaceTemplateStateOptions {
	options.Headers = param
	return options
}

// InventoryResourceRecord : Complete inventory resource details with user inputs and system generated data.
type InventoryResourceRecord struct {
	// The unique name of your Inventory.  The name can be up to 128 characters long and can include alphanumeric
	// characters, spaces, dashes, and underscores.
	Name *string `json:"name,omitempty"`

	// Inventory id.
	ID *string `json:"id,omitempty"`

	// The description of your Inventory.  The description can be up to 2048 characters long in size.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Inventory definition.  By default, Inventory will be created in Default Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Inventory creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Email address of user who created the Inventory.
	CreatedBy *string `json:"created_by,omitempty"`

	// Inventory updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Email address of user who updated the Inventory.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// Input inventory of host and host group for the playbook,  in the .ini file format.
	InventoriesIni *string `json:"inventories_ini,omitempty"`

	// Input resource queries that is used to dynamically generate  the inventory of host and host group for the playbook.
	ResourceQueries []string `json:"resource_queries,omitempty"`
}

// Constants associated with the InventoryResourceRecord.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	InventoryResourceRecord_Location_EuDe    = "eu-de"
	InventoryResourceRecord_Location_EuGb    = "eu-gb"
	InventoryResourceRecord_Location_UsEast  = "us-east"
	InventoryResourceRecord_Location_UsSouth = "us-south"
)

// UnmarshalInventoryResourceRecord unmarshals an instance of InventoryResourceRecord from the specified map of raw messages.
func UnmarshalInventoryResourceRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InventoryResourceRecord)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inventories_ini", &obj.InventoriesIni)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_queries", &obj.ResourceQueries)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InventoryResourceRecordList : List of Inventory definition records.
type InventoryResourceRecordList struct {
	// Total number of records.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// List of inventory definition records.
	Inventories []InventoryResourceRecord `json:"inventories,omitempty"`
}

// UnmarshalInventoryResourceRecordList unmarshals an instance of InventoryResourceRecordList from the specified map of raw messages.
func UnmarshalInventoryResourceRecordList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InventoryResourceRecordList)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inventories", &obj.Inventories, UnmarshalInventoryResourceRecord)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Job : Complete Job with user inputs and system generated data.
type Job struct {
	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action or Workspace.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action or Workspace.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// Job ID.
	ID *string `json:"id,omitempty"`

	// Job name, uniquely derived from the related Workspace or Action.
	Name *string `json:"name,omitempty"`

	// The description of your job is derived from the related action or workspace.  The description can be up to 2048
	// characters long in size.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name derived from the related Workspace or Action.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Job submission time.
	SubmittedAt *strfmt.DateTime `json:"submitted_at,omitempty"`

	// Email address of user who submitted the job.
	SubmittedBy *string `json:"submitted_by,omitempty"`

	// Job start time.
	StartAt *strfmt.DateTime `json:"start_at,omitempty"`

	// Job end time.
	EndAt *strfmt.DateTime `json:"end_at,omitempty"`

	// Duration of job execution; example 40 sec.
	Duration *string `json:"duration,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job data.
	Data *JobData `json:"data,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Job log store URL.
	LogStoreURL *string `json:"log_store_url,omitempty"`

	// Job state store URL.
	StateStoreURL *string `json:"state_store_url,omitempty"`

	// Job results store URL.
	ResultsURL *string `json:"results_url,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the Job.CommandObject property.
// Name of the Schematics automation resource.
const (
	Job_CommandObject_Action      = "action"
	Job_CommandObject_Environment = "environment"
	Job_CommandObject_System      = "system"
	Job_CommandObject_Workspace   = "workspace"
)

// Constants associated with the Job.CommandName property.
// Schematics job command name.
const (
	Job_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	Job_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	Job_CommandName_CreateAction         = "create_action"
	Job_CommandName_CreateCart           = "create_cart"
	Job_CommandName_CreateEnvironment    = "create_environment"
	Job_CommandName_CreateWorkspace      = "create_workspace"
	Job_CommandName_DeleteAction         = "delete_action"
	Job_CommandName_DeleteEnvironment    = "delete_environment"
	Job_CommandName_DeleteWorkspace      = "delete_workspace"
	Job_CommandName_EnvironmentInit      = "environment_init"
	Job_CommandName_EnvironmentInstall   = "environment_install"
	Job_CommandName_EnvironmentUninstall = "environment_uninstall"
	Job_CommandName_PatchAction          = "patch_action"
	Job_CommandName_PatchWorkspace       = "patch_workspace"
	Job_CommandName_PutAction            = "put_action"
	Job_CommandName_PutEnvironment       = "put_environment"
	Job_CommandName_PutWorkspace         = "put_workspace"
	Job_CommandName_RepositoryProcess    = "repository_process"
	Job_CommandName_SystemKeyDelete      = "system_key_delete"
	Job_CommandName_SystemKeyDisable     = "system_key_disable"
	Job_CommandName_SystemKeyEnable      = "system_key_enable"
	Job_CommandName_SystemKeyRestore     = "system_key_restore"
	Job_CommandName_SystemKeyRotate      = "system_key_rotate"
	Job_CommandName_WorkspaceApply       = "workspace_apply"
	Job_CommandName_WorkspaceDestroy     = "workspace_destroy"
	Job_CommandName_WorkspacePlan        = "workspace_plan"
	Job_CommandName_WorkspaceRefresh     = "workspace_refresh"
)

// Constants associated with the Job.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	Job_Location_EuDe    = "eu-de"
	Job_Location_EuGb    = "eu-gb"
	Job_Location_UsEast  = "us-east"
	Job_Location_UsSouth = "us-south"
)

// UnmarshalJob unmarshals an instance of Job from the specified map of raw messages.
func UnmarshalJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Job)
	err = core.UnmarshalPrimitive(m, "command_object", &obj.CommandObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_object_id", &obj.CommandObjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_name", &obj.CommandName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_parameter", &obj.CommandParameter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_options", &obj.CommandOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
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
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "submitted_at", &obj.SubmittedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "submitted_by", &obj.SubmittedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_at", &obj.StartAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_at", &obj.EndAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "duration", &obj.Duration)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "status", &obj.Status, UnmarshalJobStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data", &obj.Data, UnmarshalJobData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bastion", &obj.Bastion, UnmarshalBastionResourceDefinition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "log_summary", &obj.LogSummary, UnmarshalJobLogSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_store_url", &obj.LogStoreURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_store_url", &obj.StateStoreURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "results_url", &obj.ResultsURL)
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

// JobData : Job data.
type JobData struct {
	// Type of Job.
	JobType *string `json:"job_type" validate:"required"`

	// Workspace Job data.
	WorkspaceJobData *JobDataWorkspace `json:"workspace_job_data,omitempty"`

	// Action Job data.
	ActionJobData *JobDataAction `json:"action_job_data,omitempty"`

	// Controls Job data.
	SystemJobData *JobDataSystem `json:"system_job_data,omitempty"`

	// Flow Job data.
	FlowJobData *JobDataFlow `json:"flow_job_data,omitempty"`
}

// Constants associated with the JobData.JobType property.
// Type of Job.
const (
	JobData_JobType_ActionJob       = "action_job"
	JobData_JobType_FlowJob         = "flow-job"
	JobData_JobType_RepoDownloadJob = "repo_download_job"
	JobData_JobType_SystemJob       = "system_job"
	JobData_JobType_WorkspaceJob    = "workspace_job"
)

// NewJobData : Instantiate JobData (Generic Model Constructor)
func (*SchematicsV1) NewJobData(jobType string) (_model *JobData, err error) {
	_model = &JobData{
		JobType: core.StringPtr(jobType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalJobData unmarshals an instance of JobData from the specified map of raw messages.
func UnmarshalJobData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobData)
	err = core.UnmarshalPrimitive(m, "job_type", &obj.JobType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workspace_job_data", &obj.WorkspaceJobData, UnmarshalJobDataWorkspace)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_job_data", &obj.ActionJobData, UnmarshalJobDataAction)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "system_job_data", &obj.SystemJobData, UnmarshalJobDataSystem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "flow_job_data", &obj.FlowJobData, UnmarshalJobDataFlow)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobDataAction : Action Job data.
type JobDataAction struct {
	// Flow name.
	ActionName *string `json:"action_name,omitempty"`

	// Input variables data used by the Action Job.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables data from the Action Job.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables used by all the templates in the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Complete inventory resource details with user inputs and system generated data.
	InventoryRecord *InventoryResourceRecord `json:"inventory_record,omitempty"`

	// Materialized inventory details used by the Action Job, in .ini format.
	MaterializedInventory *string `json:"materialized_inventory,omitempty"`
}

// UnmarshalJobDataAction unmarshals an instance of JobDataAction from the specified map of raw messages.
func UnmarshalJobDataAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataAction)
	err = core.UnmarshalPrimitive(m, "action_name", &obj.ActionName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inventory_record", &obj.InventoryRecord, UnmarshalInventoryResourceRecord)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "materialized_inventory", &obj.MaterializedInventory)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobDataFlow : Flow Job data.
type JobDataFlow struct {
	// Flow ID.
	FlowID *string `json:"flow_id,omitempty"`

	// Flow Name.
	FlowName *string `json:"flow_name,omitempty"`

	// Job data used by each workitem Job.
	Workitems []JobDataWorkItem `json:"workitems,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalJobDataFlow unmarshals an instance of JobDataFlow from the specified map of raw messages.
func UnmarshalJobDataFlow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataFlow)
	err = core.UnmarshalPrimitive(m, "flow_id", &obj.FlowID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_name", &obj.FlowName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workitems", &obj.Workitems, UnmarshalJobDataWorkItem)
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

// JobDataSystem : Controls Job data.
type JobDataSystem struct {
	// Key ID for which key event is generated.
	KeyID *string `json:"key_id,omitempty"`

	// List of the schematics resource id.
	SchematicsResourceID []string `json:"schematics_resource_id,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalJobDataSystem unmarshals an instance of JobDataSystem from the specified map of raw messages.
func UnmarshalJobDataSystem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataSystem)
	err = core.UnmarshalPrimitive(m, "key_id", &obj.KeyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schematics_resource_id", &obj.SchematicsResourceID)
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

// JobDataTemplate : Template Job data.
type JobDataTemplate struct {
	// Template Id.
	TemplateID *string `json:"template_id,omitempty"`

	// Template name.
	TemplateName *string `json:"template_name,omitempty"`

	// Index of the template in the Flow.
	FlowIndex *int64 `json:"flow_index,omitempty"`

	// Job inputs used by the Templates.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Job output from the Templates.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables used by the template.
	Settings []VariableData `json:"settings,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalJobDataTemplate unmarshals an instance of JobDataTemplate from the specified map of raw messages.
func UnmarshalJobDataTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataTemplate)
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_name", &obj.TemplateName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_index", &obj.FlowIndex)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
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

// JobDataWorkItem : Environment work items.
type JobDataWorkItem struct {
	// command object id.
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// command object name.
	CommandObjectName *string `json:"command_object_name,omitempty"`

	// layer name.
	Layers *string `json:"layers,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Input variables data for the workItem used in FlowJob.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the workItem.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the workItem.
	Settings []VariableData `json:"settings,omitempty"`

	// Status of the last job executed by the workitem.
	LastJob *JobDataWorkItemLastJob `json:"last_job,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobDataWorkItem.SourceType property.
// Type of source for the Template.
const (
	JobDataWorkItem_SourceType_CosBucket        = "cos_bucket"
	JobDataWorkItem_SourceType_ExternalScm      = "external_scm"
	JobDataWorkItem_SourceType_GitHub           = "git_hub"
	JobDataWorkItem_SourceType_GitHubEnterprise = "git_hub_enterprise"
	JobDataWorkItem_SourceType_GitLab           = "git_lab"
	JobDataWorkItem_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	JobDataWorkItem_SourceType_IbmGitLab        = "ibm_git_lab"
	JobDataWorkItem_SourceType_Local            = "local"
)

// UnmarshalJobDataWorkItem unmarshals an instance of JobDataWorkItem from the specified map of raw messages.
func UnmarshalJobDataWorkItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataWorkItem)
	err = core.UnmarshalPrimitive(m, "command_object_id", &obj.CommandObjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_object_name", &obj.CommandObjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "layers", &obj.Layers)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalExternalSource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_job", &obj.LastJob, UnmarshalJobDataWorkItemLastJob)
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

// JobDataWorkItemLastJob : Status of the last job executed by the workitem.
type JobDataWorkItemLastJob struct {
	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// command object name (workspace_name/action_name).
	CommandObjectName *string `json:"command_object_name,omitempty"`

	// Workitem command object id, maps to workspace_id or action_id.
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Workspace job id.
	JobID *string `json:"job_id,omitempty"`

	// Status of Jobs.
	JobStatus *string `json:"job_status,omitempty"`
}

// Constants associated with the JobDataWorkItemLastJob.CommandObject property.
// Name of the Schematics automation resource.
const (
	JobDataWorkItemLastJob_CommandObject_Action      = "action"
	JobDataWorkItemLastJob_CommandObject_Environment = "environment"
	JobDataWorkItemLastJob_CommandObject_System      = "system"
	JobDataWorkItemLastJob_CommandObject_Workspace   = "workspace"
)

// Constants associated with the JobDataWorkItemLastJob.CommandName property.
// Schematics job command name.
const (
	JobDataWorkItemLastJob_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	JobDataWorkItemLastJob_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	JobDataWorkItemLastJob_CommandName_CreateAction         = "create_action"
	JobDataWorkItemLastJob_CommandName_CreateCart           = "create_cart"
	JobDataWorkItemLastJob_CommandName_CreateEnvironment    = "create_environment"
	JobDataWorkItemLastJob_CommandName_CreateWorkspace      = "create_workspace"
	JobDataWorkItemLastJob_CommandName_DeleteAction         = "delete_action"
	JobDataWorkItemLastJob_CommandName_DeleteEnvironment    = "delete_environment"
	JobDataWorkItemLastJob_CommandName_DeleteWorkspace      = "delete_workspace"
	JobDataWorkItemLastJob_CommandName_EnvironmentInit      = "environment_init"
	JobDataWorkItemLastJob_CommandName_EnvironmentInstall   = "environment_install"
	JobDataWorkItemLastJob_CommandName_EnvironmentUninstall = "environment_uninstall"
	JobDataWorkItemLastJob_CommandName_PatchAction          = "patch_action"
	JobDataWorkItemLastJob_CommandName_PatchWorkspace       = "patch_workspace"
	JobDataWorkItemLastJob_CommandName_PutAction            = "put_action"
	JobDataWorkItemLastJob_CommandName_PutEnvironment       = "put_environment"
	JobDataWorkItemLastJob_CommandName_PutWorkspace         = "put_workspace"
	JobDataWorkItemLastJob_CommandName_RepositoryProcess    = "repository_process"
	JobDataWorkItemLastJob_CommandName_SystemKeyDelete      = "system_key_delete"
	JobDataWorkItemLastJob_CommandName_SystemKeyDisable     = "system_key_disable"
	JobDataWorkItemLastJob_CommandName_SystemKeyEnable      = "system_key_enable"
	JobDataWorkItemLastJob_CommandName_SystemKeyRestore     = "system_key_restore"
	JobDataWorkItemLastJob_CommandName_SystemKeyRotate      = "system_key_rotate"
	JobDataWorkItemLastJob_CommandName_WorkspaceApply       = "workspace_apply"
	JobDataWorkItemLastJob_CommandName_WorkspaceDestroy     = "workspace_destroy"
	JobDataWorkItemLastJob_CommandName_WorkspacePlan        = "workspace_plan"
	JobDataWorkItemLastJob_CommandName_WorkspaceRefresh     = "workspace_refresh"
)

// Constants associated with the JobDataWorkItemLastJob.JobStatus property.
// Status of Jobs.
const (
	JobDataWorkItemLastJob_JobStatus_JobCancelled  = "job_cancelled"
	JobDataWorkItemLastJob_JobStatus_JobFailed     = "job_failed"
	JobDataWorkItemLastJob_JobStatus_JobFinished   = "job_finished"
	JobDataWorkItemLastJob_JobStatus_JobInProgress = "job_in_progress"
	JobDataWorkItemLastJob_JobStatus_JobPending    = "job_pending"
)

// UnmarshalJobDataWorkItemLastJob unmarshals an instance of JobDataWorkItemLastJob from the specified map of raw messages.
func UnmarshalJobDataWorkItemLastJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataWorkItemLastJob)
	err = core.UnmarshalPrimitive(m, "command_object", &obj.CommandObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_object_name", &obj.CommandObjectName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_object_id", &obj.CommandObjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_name", &obj.CommandName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_status", &obj.JobStatus)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobDataWorkspace : Workspace Job data.
type JobDataWorkspace struct {
	// Workspace name.
	WorkspaceName *string `json:"workspace_name,omitempty"`

	// Flow Id.
	FlowID *string `json:"flow_id,omitempty"`

	// Flow name.
	FlowName *string `json:"flow_name,omitempty"`

	// Input variables data used by the Workspace Job.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables data from the Workspace Job.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables used by all the templates in the Workspace.
	Settings []VariableData `json:"settings,omitempty"`

	// Input / output data of the Template in the Workspace Job.
	TemplateData []JobDataTemplate `json:"template_data,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalJobDataWorkspace unmarshals an instance of JobDataWorkspace from the specified map of raw messages.
func UnmarshalJobDataWorkspace(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobDataWorkspace)
	err = core.UnmarshalPrimitive(m, "workspace_name", &obj.WorkspaceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_id", &obj.FlowID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_name", &obj.FlowName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "settings", &obj.Settings, UnmarshalVariableData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template_data", &obj.TemplateData, UnmarshalJobDataTemplate)
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

// JobList : List of Job details.
type JobList struct {
	// Total number of records.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// List of job records.
	Jobs []JobLite `json:"jobs,omitempty"`
}

// UnmarshalJobList unmarshals an instance of JobList from the specified map of raw messages.
func UnmarshalJobList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobList)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "jobs", &obj.Jobs, UnmarshalJobLite)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLite : Job summary profile with system generated data.
type JobLite struct {
	// Job ID.
	ID *string `json:"id,omitempty"`

	// Job name, uniquely derived from the related Workspace or Action.
	Name *string `json:"name,omitempty"`

	// Job description derived from the related Workspace or Action.
	Description *string `json:"description,omitempty"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name derived from the related Workspace or Action.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Job submission time.
	SubmittedAt *strfmt.DateTime `json:"submitted_at,omitempty"`

	// Email address of user who submitted the job.
	SubmittedBy *string `json:"submitted_by,omitempty"`

	// Duration of job execution; example 40 sec.
	Duration *string `json:"duration,omitempty"`

	// Job start time.
	StartAt *strfmt.DateTime `json:"start_at,omitempty"`

	// Job end time.
	EndAt *strfmt.DateTime `json:"end_at,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobLite.CommandObject property.
// Name of the Schematics automation resource.
const (
	JobLite_CommandObject_Action      = "action"
	JobLite_CommandObject_Environment = "environment"
	JobLite_CommandObject_System      = "system"
	JobLite_CommandObject_Workspace   = "workspace"
)

// Constants associated with the JobLite.CommandName property.
// Schematics job command name.
const (
	JobLite_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	JobLite_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	JobLite_CommandName_CreateAction         = "create_action"
	JobLite_CommandName_CreateCart           = "create_cart"
	JobLite_CommandName_CreateEnvironment    = "create_environment"
	JobLite_CommandName_CreateWorkspace      = "create_workspace"
	JobLite_CommandName_DeleteAction         = "delete_action"
	JobLite_CommandName_DeleteEnvironment    = "delete_environment"
	JobLite_CommandName_DeleteWorkspace      = "delete_workspace"
	JobLite_CommandName_EnvironmentInit      = "environment_init"
	JobLite_CommandName_EnvironmentInstall   = "environment_install"
	JobLite_CommandName_EnvironmentUninstall = "environment_uninstall"
	JobLite_CommandName_PatchAction          = "patch_action"
	JobLite_CommandName_PatchWorkspace       = "patch_workspace"
	JobLite_CommandName_PutAction            = "put_action"
	JobLite_CommandName_PutEnvironment       = "put_environment"
	JobLite_CommandName_PutWorkspace         = "put_workspace"
	JobLite_CommandName_RepositoryProcess    = "repository_process"
	JobLite_CommandName_SystemKeyDelete      = "system_key_delete"
	JobLite_CommandName_SystemKeyDisable     = "system_key_disable"
	JobLite_CommandName_SystemKeyEnable      = "system_key_enable"
	JobLite_CommandName_SystemKeyRestore     = "system_key_restore"
	JobLite_CommandName_SystemKeyRotate      = "system_key_rotate"
	JobLite_CommandName_WorkspaceApply       = "workspace_apply"
	JobLite_CommandName_WorkspaceDestroy     = "workspace_destroy"
	JobLite_CommandName_WorkspacePlan        = "workspace_plan"
	JobLite_CommandName_WorkspaceRefresh     = "workspace_refresh"
)

// Constants associated with the JobLite.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	JobLite_Location_EuDe    = "eu-de"
	JobLite_Location_EuGb    = "eu-gb"
	JobLite_Location_UsEast  = "us-east"
	JobLite_Location_UsSouth = "us-south"
)

// UnmarshalJobLite unmarshals an instance of JobLite from the specified map of raw messages.
func UnmarshalJobLite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLite)
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
	err = core.UnmarshalPrimitive(m, "command_object", &obj.CommandObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_object_id", &obj.CommandObjectID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_name", &obj.CommandName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "submitted_at", &obj.SubmittedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "submitted_by", &obj.SubmittedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "duration", &obj.Duration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_at", &obj.StartAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_at", &obj.EndAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "status", &obj.Status, UnmarshalJobStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "log_summary", &obj.LogSummary, UnmarshalJobLogSummary)
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

// JobLog : Job Log details.
type JobLog struct {
	// Job Id.
	JobID *string `json:"job_id,omitempty"`

	// Job name, uniquely derived from the related Workspace, Action or Controls.
	JobName *string `json:"job_name,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Format of the Log text.
	Format *string `json:"format,omitempty"`

	// Log text, generated by the Job.
	Details *[]byte `json:"details,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobLog.Format property.
// Format of the Log text.
const (
	JobLog_Format_HTML     = "html"
	JobLog_Format_JSON     = "json"
	JobLog_Format_Markdown = "markdown"
	JobLog_Format_Rtf      = "rtf"
)

// UnmarshalJobLog unmarshals an instance of JobLog from the specified map of raw messages.
func UnmarshalJobLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLog)
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_name", &obj.JobName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "log_summary", &obj.LogSummary, UnmarshalJobLogSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "details", &obj.Details)
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

// JobLogSummary : Job log summary record.
type JobLogSummary struct {
	// Workspace Id.
	JobID *string `json:"job_id,omitempty"`

	// Type of Job.
	JobType *string `json:"job_type,omitempty"`

	// Job log start timestamp.
	LogStartAt *strfmt.DateTime `json:"log_start_at,omitempty"`

	// Job log update timestamp.
	LogAnalyzedTill *strfmt.DateTime `json:"log_analyzed_till,omitempty"`

	// Job log elapsed time (log_analyzed_till - log_start_at).
	ElapsedTime *float64 `json:"elapsed_time,omitempty"`

	// Job log errors.
	LogErrors []JobLogSummaryLogErrors `json:"log_errors,omitempty"`

	// Repo download Job log summary.
	RepoDownloadJob *JobLogSummaryRepoDownloadJob `json:"repo_download_job,omitempty"`

	// Workspace Job log summary.
	WorkspaceJob *JobLogSummaryWorkspaceJob `json:"workspace_job,omitempty"`

	// Flow Job log summary.
	FlowJob *JobLogSummaryFlowJob `json:"flow_job,omitempty"`

	// Flow Job log summary.
	ActionJob *JobLogSummaryActionJob `json:"action_job,omitempty"`

	// System Job log summary.
	SystemJob *JobLogSummarySystemJob `json:"system_job,omitempty"`
}

// Constants associated with the JobLogSummary.JobType property.
// Type of Job.
const (
	JobLogSummary_JobType_ActionJob       = "action_job"
	JobLogSummary_JobType_FlowJob         = "flow_job"
	JobLogSummary_JobType_RepoDownloadJob = "repo_download_job"
	JobLogSummary_JobType_SystemJob       = "system_job"
	JobLogSummary_JobType_WorkspaceJob    = "workspace_job"
)

// UnmarshalJobLogSummary unmarshals an instance of JobLogSummary from the specified map of raw messages.
func UnmarshalJobLogSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummary)
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_type", &obj.JobType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_start_at", &obj.LogStartAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_analyzed_till", &obj.LogAnalyzedTill)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "elapsed_time", &obj.ElapsedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "log_errors", &obj.LogErrors, UnmarshalJobLogSummaryLogErrors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "repo_download_job", &obj.RepoDownloadJob, UnmarshalJobLogSummaryRepoDownloadJob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workspace_job", &obj.WorkspaceJob, UnmarshalJobLogSummaryWorkspaceJob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "flow_job", &obj.FlowJob, UnmarshalJobLogSummaryFlowJob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_job", &obj.ActionJob, UnmarshalJobLogSummaryActionJob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "system_job", &obj.SystemJob, UnmarshalJobLogSummarySystemJob)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryWorkitems : Job log summary of the flow workitem.
type JobLogSummaryWorkitems struct {
	// workspace ID.
	WorkspaceID *string `json:"workspace_id,omitempty"`

	// workspace JOB ID.
	JobID *string `json:"job_id,omitempty"`

	// Number of resources add.
	ResourcesAdd *float64 `json:"resources_add,omitempty"`

	// Number of resources modify.
	ResourcesModify *float64 `json:"resources_modify,omitempty"`

	// Number of resources destroy.
	ResourcesDestroy *float64 `json:"resources_destroy,omitempty"`

	// Log url for job.
	LogURL *string `json:"log_url,omitempty"`
}

// UnmarshalJobLogSummaryWorkitems unmarshals an instance of JobLogSummaryWorkitems from the specified map of raw messages.
func UnmarshalJobLogSummaryWorkitems(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryWorkitems)
	err = core.UnmarshalPrimitive(m, "workspace_id", &obj.WorkspaceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_add", &obj.ResourcesAdd)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_modify", &obj.ResourcesModify)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_destroy", &obj.ResourcesDestroy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_url", &obj.LogURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryActionJob : Flow Job log summary.
type JobLogSummaryActionJob struct {
	// number of targets or hosts.
	TargetCount *float64 `json:"target_count,omitempty"`

	// number of tasks in playbook.
	TaskCount *float64 `json:"task_count,omitempty"`

	// number of plays in playbook.
	PlayCount *float64 `json:"play_count,omitempty"`

	// Recap records.
	Recap *JobLogSummaryActionJobRecap `json:"recap,omitempty"`
}

// UnmarshalJobLogSummaryActionJob unmarshals an instance of JobLogSummaryActionJob from the specified map of raw messages.
func UnmarshalJobLogSummaryActionJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryActionJob)
	err = core.UnmarshalPrimitive(m, "target_count", &obj.TargetCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "task_count", &obj.TaskCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "play_count", &obj.PlayCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "recap", &obj.Recap, UnmarshalJobLogSummaryActionJobRecap)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryActionJobRecap : Recap records.
type JobLogSummaryActionJobRecap struct {
	// List of target or host name.
	Target []string `json:"target,omitempty"`

	// Number of OK.
	Ok *float64 `json:"ok,omitempty"`

	// Number of changed.
	Changed *float64 `json:"changed,omitempty"`

	// Number of failed.
	Failed *float64 `json:"failed,omitempty"`

	// Number of skipped.
	Skipped *float64 `json:"skipped,omitempty"`

	// Number of unreachable.
	Unreachable *float64 `json:"unreachable,omitempty"`
}

// UnmarshalJobLogSummaryActionJobRecap unmarshals an instance of JobLogSummaryActionJobRecap from the specified map of raw messages.
func UnmarshalJobLogSummaryActionJobRecap(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryActionJobRecap)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "changed", &obj.Changed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped", &obj.Skipped)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unreachable", &obj.Unreachable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryFlowJob : Flow Job log summary.
type JobLogSummaryFlowJob struct {
	// Number of workitems completed successfully.
	WorkitemsCompleted *float64 `json:"workitems_completed,omitempty"`

	// Number of workitems pending in the flow.
	WorkitemsPending *float64 `json:"workitems_pending,omitempty"`

	// Number of workitems failed.
	WorkitemsFailed *float64 `json:"workitems_failed,omitempty"`

	Workitems []JobLogSummaryWorkitems `json:"workitems,omitempty"`
}

// UnmarshalJobLogSummaryFlowJob unmarshals an instance of JobLogSummaryFlowJob from the specified map of raw messages.
func UnmarshalJobLogSummaryFlowJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryFlowJob)
	err = core.UnmarshalPrimitive(m, "workitems_completed", &obj.WorkitemsCompleted)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workitems_pending", &obj.WorkitemsPending)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workitems_failed", &obj.WorkitemsFailed)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workitems", &obj.Workitems, UnmarshalJobLogSummaryWorkitems)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryLogErrors : JobLogSummaryLogErrors struct
type JobLogSummaryLogErrors struct {
	// Error code in the Log.
	ErrorCode *string `json:"error_code,omitempty"`

	// Summary error message in the log.
	ErrorMsg *string `json:"error_msg,omitempty"`

	// Number of occurrence.
	ErrorCount *float64 `json:"error_count,omitempty"`
}

// UnmarshalJobLogSummaryLogErrors unmarshals an instance of JobLogSummaryLogErrors from the specified map of raw messages.
func UnmarshalJobLogSummaryLogErrors(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryLogErrors)
	err = core.UnmarshalPrimitive(m, "error_code", &obj.ErrorCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_msg", &obj.ErrorMsg)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryRepoDownloadJob : Repo download Job log summary.
type JobLogSummaryRepoDownloadJob struct {
	// Number of files scanned.
	ScannedFileCount *float64 `json:"scanned_file_count,omitempty"`

	// Number of files quarantined.
	QuarantinedFileCount *float64 `json:"quarantined_file_count,omitempty"`

	// Detected template or data file type.
	DetectedFiletype *string `json:"detected_filetype,omitempty"`

	// Number of inputs detected.
	InputsCount *string `json:"inputs_count,omitempty"`

	// Number of outputs detected.
	OutputsCount *string `json:"outputs_count,omitempty"`
}

// UnmarshalJobLogSummaryRepoDownloadJob unmarshals an instance of JobLogSummaryRepoDownloadJob from the specified map of raw messages.
func UnmarshalJobLogSummaryRepoDownloadJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryRepoDownloadJob)
	err = core.UnmarshalPrimitive(m, "scanned_file_count", &obj.ScannedFileCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "quarantined_file_count", &obj.QuarantinedFileCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "detected_filetype", &obj.DetectedFiletype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs_count", &obj.InputsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "outputs_count", &obj.OutputsCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummarySystemJob : System Job log summary.
type JobLogSummarySystemJob struct {
	// number of targets or hosts.
	TargetCount *float64 `json:"target_count,omitempty"`

	// Number of passed.
	Success *float64 `json:"success,omitempty"`

	// Number of failed.
	Failed *float64 `json:"failed,omitempty"`
}

// UnmarshalJobLogSummarySystemJob unmarshals an instance of JobLogSummarySystemJob from the specified map of raw messages.
func UnmarshalJobLogSummarySystemJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummarySystemJob)
	err = core.UnmarshalPrimitive(m, "target_count", &obj.TargetCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobLogSummaryWorkspaceJob : Workspace Job log summary.
type JobLogSummaryWorkspaceJob struct {
	// Number of resources add.
	ResourcesAdd *float64 `json:"resources_add,omitempty"`

	// Number of resources modify.
	ResourcesModify *float64 `json:"resources_modify,omitempty"`

	// Number of resources destroy.
	ResourcesDestroy *float64 `json:"resources_destroy,omitempty"`
}

// UnmarshalJobLogSummaryWorkspaceJob unmarshals an instance of JobLogSummaryWorkspaceJob from the specified map of raw messages.
func UnmarshalJobLogSummaryWorkspaceJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryWorkspaceJob)
	err = core.UnmarshalPrimitive(m, "resources_add", &obj.ResourcesAdd)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_modify", &obj.ResourcesModify)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_destroy", &obj.ResourcesDestroy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobStatus : Job Status.
type JobStatus struct {
	// Workspace Job Status.
	WorkspaceJobStatus *JobStatusWorkspace `json:"workspace_job_status,omitempty"`

	// Action Job Status.
	ActionJobStatus *JobStatusAction `json:"action_job_status,omitempty"`

	// System Job Status.
	SystemJobStatus *JobStatusSystem `json:"system_job_status,omitempty"`

	// Environment Flow JOB Status.
	FlowJobStatus *JobStatusFlow `json:"flow_job_status,omitempty"`
}

// UnmarshalJobStatus unmarshals an instance of JobStatus from the specified map of raw messages.
func UnmarshalJobStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatus)
	err = core.UnmarshalModel(m, "workspace_job_status", &obj.WorkspaceJobStatus, UnmarshalJobStatusWorkspace)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_job_status", &obj.ActionJobStatus, UnmarshalJobStatusAction)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "system_job_status", &obj.SystemJobStatus, UnmarshalJobStatusSystem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "flow_job_status", &obj.FlowJobStatus, UnmarshalJobStatusFlow)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobStatusAction : Action Job Status.
type JobStatusAction struct {
	// Action name.
	ActionName *string `json:"action_name,omitempty"`

	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// Action Job status message - to be displayed along with the action_status_code.
	StatusMessage *string `json:"status_message,omitempty"`

	// Status of Resources.
	BastionStatusCode *string `json:"bastion_status_code,omitempty"`

	// Bastion status message - to be displayed along with the bastion_status_code;.
	BastionStatusMessage *string `json:"bastion_status_message,omitempty"`

	// Status of Resources.
	TargetsStatusCode *string `json:"targets_status_code,omitempty"`

	// Aggregated status message for all target resources,  to be displayed along with the targets_status_code;.
	TargetsStatusMessage *string `json:"targets_status_message,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusAction.StatusCode property.
// Status of Jobs.
const (
	JobStatusAction_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusAction_StatusCode_JobFailed     = "job_failed"
	JobStatusAction_StatusCode_JobFinished   = "job_finished"
	JobStatusAction_StatusCode_JobInProgress = "job_in_progress"
	JobStatusAction_StatusCode_JobPending    = "job_pending"
)

// Constants associated with the JobStatusAction.BastionStatusCode property.
// Status of Resources.
const (
	JobStatusAction_BastionStatusCode_Error      = "error"
	JobStatusAction_BastionStatusCode_None       = "none"
	JobStatusAction_BastionStatusCode_Processing = "processing"
	JobStatusAction_BastionStatusCode_Ready      = "ready"
)

// Constants associated with the JobStatusAction.TargetsStatusCode property.
// Status of Resources.
const (
	JobStatusAction_TargetsStatusCode_Error      = "error"
	JobStatusAction_TargetsStatusCode_None       = "none"
	JobStatusAction_TargetsStatusCode_Processing = "processing"
	JobStatusAction_TargetsStatusCode_Ready      = "ready"
)

// UnmarshalJobStatusAction unmarshals an instance of JobStatusAction from the specified map of raw messages.
func UnmarshalJobStatusAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusAction)
	err = core.UnmarshalPrimitive(m, "action_name", &obj.ActionName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bastion_status_code", &obj.BastionStatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bastion_status_message", &obj.BastionStatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "targets_status_code", &obj.TargetsStatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "targets_status_message", &obj.TargetsStatusMessage)
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

// JobStatusFlow : Environment Flow JOB Status.
type JobStatusFlow struct {
	// flow id.
	FlowID *string `json:"flow_id,omitempty"`

	// flow name.
	FlowName *string `json:"flow_name,omitempty"`

	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// Flow Job status message - to be displayed along with the status_code;.
	StatusMessage *string `json:"status_message,omitempty"`

	// Environment's individual workItem status details;.
	Workitems []JobStatusWorkitem `json:"workitems,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusFlow.StatusCode property.
// Status of Jobs.
const (
	JobStatusFlow_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusFlow_StatusCode_JobFailed     = "job_failed"
	JobStatusFlow_StatusCode_JobFinished   = "job_finished"
	JobStatusFlow_StatusCode_JobInProgress = "job_in_progress"
	JobStatusFlow_StatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusFlow unmarshals an instance of JobStatusFlow from the specified map of raw messages.
func UnmarshalJobStatusFlow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusFlow)
	err = core.UnmarshalPrimitive(m, "flow_id", &obj.FlowID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_name", &obj.FlowName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workitems", &obj.Workitems, UnmarshalJobStatusWorkitem)
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

// JobStatusSchematicsResources : schematics Resources Job Status.
type JobStatusSchematicsResources struct {
	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// system job status message.
	StatusMessage *string `json:"status_message,omitempty"`

	// id for each resource which is targeted as a part of system job.
	SchematicsResourceID *string `json:"schematics_resource_id,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusSchematicsResources.StatusCode property.
// Status of Jobs.
const (
	JobStatusSchematicsResources_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusSchematicsResources_StatusCode_JobFailed     = "job_failed"
	JobStatusSchematicsResources_StatusCode_JobFinished   = "job_finished"
	JobStatusSchematicsResources_StatusCode_JobInProgress = "job_in_progress"
	JobStatusSchematicsResources_StatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusSchematicsResources unmarshals an instance of JobStatusSchematicsResources from the specified map of raw messages.
func UnmarshalJobStatusSchematicsResources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusSchematicsResources)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schematics_resource_id", &obj.SchematicsResourceID)
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

// JobStatusSystem : System Job Status.
type JobStatusSystem struct {
	// System job message.
	SystemStatusMessage *string `json:"system_status_message,omitempty"`

	// Status of Jobs.
	SystemStatusCode *string `json:"system_status_code,omitempty"`

	// job staus for each schematics resource.
	SchematicsResourceStatus []JobStatusSchematicsResources `json:"schematics_resource_status,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusSystem.SystemStatusCode property.
// Status of Jobs.
const (
	JobStatusSystem_SystemStatusCode_JobCancelled  = "job_cancelled"
	JobStatusSystem_SystemStatusCode_JobFailed     = "job_failed"
	JobStatusSystem_SystemStatusCode_JobFinished   = "job_finished"
	JobStatusSystem_SystemStatusCode_JobInProgress = "job_in_progress"
	JobStatusSystem_SystemStatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusSystem unmarshals an instance of JobStatusSystem from the specified map of raw messages.
func UnmarshalJobStatusSystem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusSystem)
	err = core.UnmarshalPrimitive(m, "system_status_message", &obj.SystemStatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "system_status_code", &obj.SystemStatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "schematics_resource_status", &obj.SchematicsResourceStatus, UnmarshalJobStatusSchematicsResources)
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

// JobStatusTemplate : Template Job Status.
type JobStatusTemplate struct {
	// Template Id.
	TemplateID *string `json:"template_id,omitempty"`

	// Template name.
	TemplateName *string `json:"template_name,omitempty"`

	// Index of the template in the Flow.
	FlowIndex *int64 `json:"flow_index,omitempty"`

	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// Template job status message (eg. VPCt1_Apply_Pending, for a 'VPCt1' Template).
	StatusMessage *string `json:"status_message,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusTemplate.StatusCode property.
// Status of Jobs.
const (
	JobStatusTemplate_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusTemplate_StatusCode_JobFailed     = "job_failed"
	JobStatusTemplate_StatusCode_JobFinished   = "job_finished"
	JobStatusTemplate_StatusCode_JobInProgress = "job_in_progress"
	JobStatusTemplate_StatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusTemplate unmarshals an instance of JobStatusTemplate from the specified map of raw messages.
func UnmarshalJobStatusTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusTemplate)
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_name", &obj.TemplateName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "flow_index", &obj.FlowIndex)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
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

// JobStatusWorkitem : Individual workitem status info.
type JobStatusWorkitem struct {
	// Workspace id.
	WorkspaceID *string `json:"workspace_id,omitempty"`

	// workspace name.
	WorkspaceName *string `json:"workspace_name,omitempty"`

	// workspace job id.
	JobID *string `json:"job_id,omitempty"`

	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// workitem job status message;.
	StatusMessage *string `json:"status_message,omitempty"`

	// workitem job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusWorkitem.StatusCode property.
// Status of Jobs.
const (
	JobStatusWorkitem_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusWorkitem_StatusCode_JobFailed     = "job_failed"
	JobStatusWorkitem_StatusCode_JobFinished   = "job_finished"
	JobStatusWorkitem_StatusCode_JobInProgress = "job_in_progress"
	JobStatusWorkitem_StatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusWorkitem unmarshals an instance of JobStatusWorkitem from the specified map of raw messages.
func UnmarshalJobStatusWorkitem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusWorkitem)
	err = core.UnmarshalPrimitive(m, "workspace_id", &obj.WorkspaceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workspace_name", &obj.WorkspaceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
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

// JobStatusWorkspace : Workspace Job Status.
type JobStatusWorkspace struct {
	// Workspace name.
	WorkspaceName *string `json:"workspace_name,omitempty"`

	// Status of Jobs.
	StatusCode *string `json:"status_code,omitempty"`

	// Workspace job status message (eg. App1_Setup_Pending, for a 'Setup' flow in the 'App1' Workspace).
	StatusMessage *string `json:"status_message,omitempty"`

	// Environment Flow JOB Status.
	FlowStatus *JobStatusFlow `json:"flow_status,omitempty"`

	// Workspace Flow Template job status.
	TemplateStatus []JobStatusTemplate `json:"template_status,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// Constants associated with the JobStatusWorkspace.StatusCode property.
// Status of Jobs.
const (
	JobStatusWorkspace_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusWorkspace_StatusCode_JobFailed     = "job_failed"
	JobStatusWorkspace_StatusCode_JobFinished   = "job_finished"
	JobStatusWorkspace_StatusCode_JobInProgress = "job_in_progress"
	JobStatusWorkspace_StatusCode_JobPending    = "job_pending"
)

// UnmarshalJobStatusWorkspace unmarshals an instance of JobStatusWorkspace from the specified map of raw messages.
func UnmarshalJobStatusWorkspace(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusWorkspace)
	err = core.UnmarshalPrimitive(m, "workspace_name", &obj.WorkspaceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_message", &obj.StatusMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "flow_status", &obj.FlowStatus, UnmarshalJobStatusFlow)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template_status", &obj.TemplateStatus, UnmarshalJobStatusTemplate)
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

// KMSDiscovery : Discovered KMS instances.
type KMSDiscovery struct {
	// Total number of records.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// List of KMS instances.
	KmsInstances []KMSInstances `json:"kms_instances,omitempty"`
}

// UnmarshalKMSDiscovery unmarshals an instance of KMSDiscovery from the specified map of raw messages.
func UnmarshalKMSDiscovery(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSDiscovery)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "kms_instances", &obj.KmsInstances, UnmarshalKMSInstances)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KMSInstances : KMS Instances.
type KMSInstances struct {
	// Location.
	Location *string `json:"location,omitempty"`

	// Encryption schema.
	EncryptionScheme *string `json:"encryption_scheme,omitempty"`

	// Resource groups.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// KMS CRN.
	KmsCrn *string `json:"kms_crn,omitempty"`

	// KMS Name.
	KmsName *string `json:"kms_name,omitempty"`

	// KMS private endpoint.
	KmsPrivateEndpoint *string `json:"kms_private_endpoint,omitempty"`

	// KMS public endpoint.
	KmsPublicEndpoint *string `json:"kms_public_endpoint,omitempty"`

	// List of keys.
	Keys []KMSInstancesKeys `json:"keys,omitempty"`
}

// UnmarshalKMSInstances unmarshals an instance of KMSInstances from the specified map of raw messages.
func UnmarshalKMSInstances(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSInstances)
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "encryption_scheme", &obj.EncryptionScheme)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_crn", &obj.KmsCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_name", &obj.KmsName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_private_endpoint", &obj.KmsPrivateEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_public_endpoint", &obj.KmsPublicEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "keys", &obj.Keys, UnmarshalKMSInstancesKeys)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KMSInstancesKeys : KMSInstancesKeys struct
type KMSInstancesKeys struct {
	// Key name.
	Name *string `json:"name,omitempty"`

	// CRN of the Key.
	Crn *string `json:"crn,omitempty"`

	// Error message.
	Error *string `json:"error,omitempty"`
}

// UnmarshalKMSInstancesKeys unmarshals an instance of KMSInstancesKeys from the specified map of raw messages.
func UnmarshalKMSInstancesKeys(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSInstancesKeys)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
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

// KMSSettings : User defined KMS Settings details.
type KMSSettings struct {
	// Location.
	Location *string `json:"location,omitempty"`

	// Encryption scheme.
	EncryptionScheme *string `json:"encryption_scheme,omitempty"`

	// Resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Primary CRK details.
	PrimaryCrk *KMSSettingsPrimaryCrk `json:"primary_crk,omitempty"`

	// Secondary CRK details.
	SecondaryCrk *KMSSettingsSecondaryCrk `json:"secondary_crk,omitempty"`
}

// UnmarshalKMSSettings unmarshals an instance of KMSSettings from the specified map of raw messages.
func UnmarshalKMSSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSSettings)
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "encryption_scheme", &obj.EncryptionScheme)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "primary_crk", &obj.PrimaryCrk, UnmarshalKMSSettingsPrimaryCrk)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secondary_crk", &obj.SecondaryCrk, UnmarshalKMSSettingsSecondaryCrk)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KMSSettingsPrimaryCrk : Primary CRK details.
type KMSSettingsPrimaryCrk struct {
	// Primary KMS name.
	KmsName *string `json:"kms_name,omitempty"`

	// Primary KMS endpoint.
	KmsPrivateEndpoint *string `json:"kms_private_endpoint,omitempty"`

	// CRN of the Primary Key.
	KeyCrn *string `json:"key_crn,omitempty"`
}

// UnmarshalKMSSettingsPrimaryCrk unmarshals an instance of KMSSettingsPrimaryCrk from the specified map of raw messages.
func UnmarshalKMSSettingsPrimaryCrk(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSSettingsPrimaryCrk)
	err = core.UnmarshalPrimitive(m, "kms_name", &obj.KmsName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_private_endpoint", &obj.KmsPrivateEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_crn", &obj.KeyCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KMSSettingsSecondaryCrk : Secondary CRK details.
type KMSSettingsSecondaryCrk struct {
	// Secondary KMS name.
	KmsName *string `json:"kms_name,omitempty"`

	// Secondary KMS endpoint.
	KmsPrivateEndpoint *string `json:"kms_private_endpoint,omitempty"`

	// CRN of the Secondary Key.
	KeyCrn *string `json:"key_crn,omitempty"`
}

// UnmarshalKMSSettingsSecondaryCrk unmarshals an instance of KMSSettingsSecondaryCrk from the specified map of raw messages.
func UnmarshalKMSSettingsSecondaryCrk(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSSettingsSecondaryCrk)
	err = core.UnmarshalPrimitive(m, "kms_name", &obj.KmsName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kms_private_endpoint", &obj.KmsPrivateEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_crn", &obj.KeyCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListActionsOptions : The ListActions options.
type ListActionsOptions struct {
	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"-"`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListActionsOptions.Profile property.
// Level of details returned by the get method.
const (
	ListActionsOptions_Profile_Ids     = "ids"
	ListActionsOptions_Profile_Summary = "summary"
)

// NewListActionsOptions : Instantiate ListActionsOptions
func (*SchematicsV1) NewListActionsOptions() *ListActionsOptions {
	return &ListActionsOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListActionsOptions) SetOffset(offset int64) *ListActionsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListActionsOptions) SetLimit(limit int64) *ListActionsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListActionsOptions) SetSort(sort string) *ListActionsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *ListActionsOptions) SetProfile(profile string) *ListActionsOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListActionsOptions) SetHeaders(param map[string]string) *ListActionsOptions {
	options.Headers = param
	return options
}

// ListInventoriesOptions : The ListInventories options.
type ListInventoriesOptions struct {
	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"-"`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListInventoriesOptions.Profile property.
// Level of details returned by the get method.
const (
	ListInventoriesOptions_Profile_Ids     = "ids"
	ListInventoriesOptions_Profile_Summary = "summary"
)

// NewListInventoriesOptions : Instantiate ListInventoriesOptions
func (*SchematicsV1) NewListInventoriesOptions() *ListInventoriesOptions {
	return &ListInventoriesOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListInventoriesOptions) SetOffset(offset int64) *ListInventoriesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListInventoriesOptions) SetLimit(limit int64) *ListInventoriesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListInventoriesOptions) SetSort(sort string) *ListInventoriesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *ListInventoriesOptions) SetProfile(profile string) *ListInventoriesOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListInventoriesOptions) SetHeaders(param map[string]string) *ListInventoriesOptions {
	options.Headers = param
	return options
}

// ListJobLogsOptions : The ListJobLogs options.
type ListJobLogsOptions struct {
	// Job Id. Use `GET /v2/jobs` API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListJobLogsOptions : Instantiate ListJobLogsOptions
func (*SchematicsV1) NewListJobLogsOptions(jobID string) *ListJobLogsOptions {
	return &ListJobLogsOptions{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *ListJobLogsOptions) SetJobID(jobID string) *ListJobLogsOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobLogsOptions) SetHeaders(param map[string]string) *ListJobLogsOptions {
	options.Headers = param
	return options
}

// ListJobsOptions : The ListJobs options.
type ListJobsOptions struct {
	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"-"`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Name of the resource (workspace, actions or controls).
	Resource *string `json:"-"`

	// The Resource Id. It could be an Action-id or Workspace-id.
	ResourceID *string `json:"-"`

	// Action Id.
	ActionID *string `json:"-"`

	// list jobs.
	List *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListJobsOptions.Profile property.
// Level of details returned by the get method.
const (
	ListJobsOptions_Profile_Ids     = "ids"
	ListJobsOptions_Profile_Summary = "summary"
)

// Constants associated with the ListJobsOptions.Resource property.
// Name of the resource (workspace, actions or controls).
const (
	ListJobsOptions_Resource_Action    = "action"
	ListJobsOptions_Resource_Workspace = "workspace"
)

// Constants associated with the ListJobsOptions.List property.
// list jobs.
const (
	ListJobsOptions_List_All = "all"
)

// NewListJobsOptions : Instantiate ListJobsOptions
func (*SchematicsV1) NewListJobsOptions() *ListJobsOptions {
	return &ListJobsOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListJobsOptions) SetOffset(offset int64) *ListJobsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListJobsOptions) SetLimit(limit int64) *ListJobsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListJobsOptions) SetSort(sort string) *ListJobsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *ListJobsOptions) SetProfile(profile string) *ListJobsOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetResource : Allow user to set Resource
func (_options *ListJobsOptions) SetResource(resource string) *ListJobsOptions {
	_options.Resource = core.StringPtr(resource)
	return _options
}

// SetResourceID : Allow user to set ResourceID
func (_options *ListJobsOptions) SetResourceID(resourceID string) *ListJobsOptions {
	_options.ResourceID = core.StringPtr(resourceID)
	return _options
}

// SetActionID : Allow user to set ActionID
func (_options *ListJobsOptions) SetActionID(actionID string) *ListJobsOptions {
	_options.ActionID = core.StringPtr(actionID)
	return _options
}

// SetList : Allow user to set List
func (_options *ListJobsOptions) SetList(list string) *ListJobsOptions {
	_options.List = core.StringPtr(list)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobsOptions) SetHeaders(param map[string]string) *ListJobsOptions {
	options.Headers = param
	return options
}

// ListKmsOptions : The ListKms options.
type ListKmsOptions struct {
	// The encryption scheme to be used.
	EncryptionScheme *string `json:"-" validate:"required"`

	// The location of the Resource.
	Location *string `json:"-" validate:"required"`

	// The resource group (by default, fetch from all resource groups).
	ResourceGroup *string `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListKmsOptions : Instantiate ListKmsOptions
func (*SchematicsV1) NewListKmsOptions(encryptionScheme string, location string) *ListKmsOptions {
	return &ListKmsOptions{
		EncryptionScheme: core.StringPtr(encryptionScheme),
		Location:         core.StringPtr(location),
	}
}

// SetEncryptionScheme : Allow user to set EncryptionScheme
func (_options *ListKmsOptions) SetEncryptionScheme(encryptionScheme string) *ListKmsOptions {
	_options.EncryptionScheme = core.StringPtr(encryptionScheme)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *ListKmsOptions) SetLocation(location string) *ListKmsOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *ListKmsOptions) SetResourceGroup(resourceGroup string) *ListKmsOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListKmsOptions) SetLimit(limit int64) *ListKmsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListKmsOptions) SetSort(sort string) *ListKmsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListKmsOptions) SetHeaders(param map[string]string) *ListKmsOptions {
	options.Headers = param
	return options
}

// ListLocationsOptions : The ListLocations options.
type ListLocationsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListLocationsOptions : Instantiate ListLocationsOptions
func (*SchematicsV1) NewListLocationsOptions() *ListLocationsOptions {
	return &ListLocationsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListLocationsOptions) SetHeaders(param map[string]string) *ListLocationsOptions {
	options.Headers = param
	return options
}

// ListResourceGroupOptions : The ListResourceGroup options.
type ListResourceGroupOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListResourceGroupOptions : Instantiate ListResourceGroupOptions
func (*SchematicsV1) NewListResourceGroupOptions() *ListResourceGroupOptions {
	return &ListResourceGroupOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListResourceGroupOptions) SetHeaders(param map[string]string) *ListResourceGroupOptions {
	options.Headers = param
	return options
}

// ListResourceQueryOptions : The ListResourceQuery options.
type ListResourceQueryOptions struct {
	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"-"`

	// Level of details returned by the get method.
	Profile *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListResourceQueryOptions.Profile property.
// Level of details returned by the get method.
const (
	ListResourceQueryOptions_Profile_Ids     = "ids"
	ListResourceQueryOptions_Profile_Summary = "summary"
)

// NewListResourceQueryOptions : Instantiate ListResourceQueryOptions
func (*SchematicsV1) NewListResourceQueryOptions() *ListResourceQueryOptions {
	return &ListResourceQueryOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListResourceQueryOptions) SetOffset(offset int64) *ListResourceQueryOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListResourceQueryOptions) SetLimit(limit int64) *ListResourceQueryOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListResourceQueryOptions) SetSort(sort string) *ListResourceQueryOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetProfile : Allow user to set Profile
func (_options *ListResourceQueryOptions) SetProfile(profile string) *ListResourceQueryOptions {
	_options.Profile = core.StringPtr(profile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListResourceQueryOptions) SetHeaders(param map[string]string) *ListResourceQueryOptions {
	options.Headers = param
	return options
}

// ListSchematicsLocationOptions : The ListSchematicsLocation options.
type ListSchematicsLocationOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSchematicsLocationOptions : Instantiate ListSchematicsLocationOptions
func (*SchematicsV1) NewListSchematicsLocationOptions() *ListSchematicsLocationOptions {
	return &ListSchematicsLocationOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListSchematicsLocationOptions) SetHeaders(param map[string]string) *ListSchematicsLocationOptions {
	options.Headers = param
	return options
}

// ListWorkspaceActivitiesOptions : The ListWorkspaceActivities options.
type ListWorkspaceActivitiesOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListWorkspaceActivitiesOptions : Instantiate ListWorkspaceActivitiesOptions
func (*SchematicsV1) NewListWorkspaceActivitiesOptions(wID string) *ListWorkspaceActivitiesOptions {
	return &ListWorkspaceActivitiesOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *ListWorkspaceActivitiesOptions) SetWID(wID string) *ListWorkspaceActivitiesOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListWorkspaceActivitiesOptions) SetOffset(offset int64) *ListWorkspaceActivitiesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListWorkspaceActivitiesOptions) SetLimit(limit int64) *ListWorkspaceActivitiesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListWorkspaceActivitiesOptions) SetHeaders(param map[string]string) *ListWorkspaceActivitiesOptions {
	options.Headers = param
	return options
}

// ListWorkspacesOptions : The ListWorkspaces options.
type ListWorkspacesOptions struct {
	// The starting position of the item in the list of items. For example, if you have three workspaces in your account,
	// the first workspace is assigned position number 0, the second workspace is assigned position number 1, and so forth.
	// If you have 6 workspaces and you want to list the details for workspaces 2-6, enter 1. To limit the number of
	// workspaces that is returned, use the `limit` option in addition to the `offset` option. Negative numbers are not
	// supported and are ignored.
	Offset *int64 `json:"-"`

	// The maximum number of items that you want to list. The number must be a positive integer between 1 and 2000. If no
	// value is provided, 100 is used by default.
	Limit *int64 `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListWorkspacesOptions : Instantiate ListWorkspacesOptions
func (*SchematicsV1) NewListWorkspacesOptions() *ListWorkspacesOptions {
	return &ListWorkspacesOptions{}
}

// SetOffset : Allow user to set Offset
func (_options *ListWorkspacesOptions) SetOffset(offset int64) *ListWorkspacesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListWorkspacesOptions) SetLimit(limit int64) *ListWorkspacesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListWorkspacesOptions) SetHeaders(param map[string]string) *ListWorkspacesOptions {
	options.Headers = param
	return options
}

// LogStoreResponse : Log file URL for job that ran against your workspace.
type LogStoreResponse struct {
	// The provisioning engine that was used for the job.
	EngineName *string `json:"engine_name,omitempty"`

	// The version of the provisioning engine that was used for the job.
	EngineVersion *string `json:"engine_version,omitempty"`

	// The ID that was assigned to your Terraform template of IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// The URL to access the logs that were created during the plan, apply, or destroy job.
	LogStoreURL *string `json:"log_store_url,omitempty"`
}

// UnmarshalLogStoreResponse unmarshals an instance of LogStoreResponse from the specified map of raw messages.
func UnmarshalLogStoreResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogStoreResponse)
	err = core.UnmarshalPrimitive(m, "engine_name", &obj.EngineName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "engine_version", &obj.EngineVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_store_url", &obj.LogStoreURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogStoreResponseList : List of log file URL that ran against your workspace.
type LogStoreResponseList struct {
	// Runtime data.
	RuntimeData []LogStoreResponse `json:"runtime_data,omitempty"`
}

// UnmarshalLogStoreResponseList unmarshals an instance of LogStoreResponseList from the specified map of raw messages.
func UnmarshalLogStoreResponseList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogStoreResponseList)
	err = core.UnmarshalModel(m, "runtime_data", &obj.RuntimeData, UnmarshalLogStoreResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogSummary : Summary information extracted from the job logs.
type LogSummary struct {
	// The status of your activity or job. To retrieve the URL to your job logs, use the GET
	// /v1/workspaces/{id}/actions/{action_id}/logs API.
	//
	// * **COMPLETED**: The job completed successfully.
	// * **CREATED**: The job was created, but the provisioning, modification, or removal of IBM Cloud resources has not
	// started yet.
	// * **FAILED**: An error occurred during the plan, apply, or destroy job. Use the job ID to retrieve the URL to the
	// log files for your job.
	// * **IN PROGRESS**: The job is in progress. You can use the log_url to access the logs.
	ActivityStatus *string `json:"activity_status,omitempty"`

	// Template detected type.
	DetectedTemplateType *string `json:"detected_template_type,omitempty"`

	// Numner of discarded files.
	DiscardedFiles *int64 `json:"discarded_files,omitempty"`

	// Numner of errors in log.
	Error *string `json:"error,omitempty"`

	// Numner of resources added.
	ResourcesAdded *int64 `json:"resources_added,omitempty"`

	// Numner of resources destroyed.
	ResourcesDestroyed *int64 `json:"resources_destroyed,omitempty"`

	// Numner of resources modified.
	ResourcesModified *int64 `json:"resources_modified,omitempty"`

	// Numner of filed scanned.
	ScannedFiles *int64 `json:"scanned_files,omitempty"`

	// Numner of template variables.
	TemplateVariableCount *int64 `json:"template_variable_count,omitempty"`

	// Elapsed time to run the job.
	TimeTaken *float64 `json:"time_taken,omitempty"`
}

// UnmarshalLogSummary unmarshals an instance of LogSummary from the specified map of raw messages.
func UnmarshalLogSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogSummary)
	err = core.UnmarshalPrimitive(m, "activity_status", &obj.ActivityStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "detected_template_type", &obj.DetectedTemplateType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "discarded_files", &obj.DiscardedFiles)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_added", &obj.ResourcesAdded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_destroyed", &obj.ResourcesDestroyed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_modified", &obj.ResourcesModified)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scanned_files", &obj.ScannedFiles)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_variable_count", &obj.TemplateVariableCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "time_taken", &obj.TimeTaken)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OutputValuesInner : OutputValuesInner struct
type OutputValuesInner struct {
	// The subfolder in the GitHub or GitLab repository where your Terraform template is stored. If the template is stored
	// in the root directory, `.` is returned.
	Folder *string `json:"folder,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// A list of Terraform output values.
	OutputValues []interface{} `json:"output_values,omitempty"`

	// The Terraform version that was used to apply your template.
	ValueType *string `json:"value_type,omitempty"`
}

// UnmarshalOutputValuesInner unmarshals an instance of OutputValuesInner from the specified map of raw messages.
func UnmarshalOutputValuesInner(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OutputValuesInner)
	err = core.UnmarshalPrimitive(m, "folder", &obj.Folder)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "output_values", &obj.OutputValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value_type", &obj.ValueType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PlanWorkspaceCommandOptions : The PlanWorkspaceCommand options.
type PlanWorkspaceCommandOptions struct {
	// The ID of the workspace, for which you want to run a Schematics `plan` job.  To find the ID of your workspace, use
	// the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// The IAM delegated token for your IBM Cloud account.  This token is required for requests that are sent via the UI
	// only.
	DelegatedToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPlanWorkspaceCommandOptions : Instantiate PlanWorkspaceCommandOptions
func (*SchematicsV1) NewPlanWorkspaceCommandOptions(wID string, refreshToken string) *PlanWorkspaceCommandOptions {
	return &PlanWorkspaceCommandOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (_options *PlanWorkspaceCommandOptions) SetWID(wID string) *PlanWorkspaceCommandOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *PlanWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *PlanWorkspaceCommandOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetDelegatedToken : Allow user to set DelegatedToken
func (_options *PlanWorkspaceCommandOptions) SetDelegatedToken(delegatedToken string) *PlanWorkspaceCommandOptions {
	_options.DelegatedToken = core.StringPtr(delegatedToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PlanWorkspaceCommandOptions) SetHeaders(param map[string]string) *PlanWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// ProcessTemplateMetaDataOptions : The ProcessTemplateMetaData options.
type ProcessTemplateMetaDataOptions struct {
	// Template type (terraform, ansible, helm, cloudpak, bash script).
	TemplateType *string `json:"template_type" validate:"required"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source" validate:"required"`

	// Region to which request should go. Applicable only on global endpoint.
	Region *string `json:"region,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform
	// template.
	XGithubToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ProcessTemplateMetaDataOptions.SourceType property.
// Type of source for the Template.
const (
	ProcessTemplateMetaDataOptions_SourceType_CosBucket        = "cos_bucket"
	ProcessTemplateMetaDataOptions_SourceType_ExternalScm      = "external_scm"
	ProcessTemplateMetaDataOptions_SourceType_GitHub           = "git_hub"
	ProcessTemplateMetaDataOptions_SourceType_GitHubEnterprise = "git_hub_enterprise"
	ProcessTemplateMetaDataOptions_SourceType_GitLab           = "git_lab"
	ProcessTemplateMetaDataOptions_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	ProcessTemplateMetaDataOptions_SourceType_IbmGitLab        = "ibm_git_lab"
	ProcessTemplateMetaDataOptions_SourceType_Local            = "local"
)

// NewProcessTemplateMetaDataOptions : Instantiate ProcessTemplateMetaDataOptions
func (*SchematicsV1) NewProcessTemplateMetaDataOptions(templateType string, source *ExternalSource) *ProcessTemplateMetaDataOptions {
	return &ProcessTemplateMetaDataOptions{
		TemplateType: core.StringPtr(templateType),
		Source:       source,
	}
}

// SetTemplateType : Allow user to set TemplateType
func (_options *ProcessTemplateMetaDataOptions) SetTemplateType(templateType string) *ProcessTemplateMetaDataOptions {
	_options.TemplateType = core.StringPtr(templateType)
	return _options
}

// SetSource : Allow user to set Source
func (_options *ProcessTemplateMetaDataOptions) SetSource(source *ExternalSource) *ProcessTemplateMetaDataOptions {
	_options.Source = source
	return _options
}

// SetRegion : Allow user to set Region
func (_options *ProcessTemplateMetaDataOptions) SetRegion(region string) *ProcessTemplateMetaDataOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetSourceType : Allow user to set SourceType
func (_options *ProcessTemplateMetaDataOptions) SetSourceType(sourceType string) *ProcessTemplateMetaDataOptions {
	_options.SourceType = core.StringPtr(sourceType)
	return _options
}

// SetXGithubToken : Allow user to set XGithubToken
func (_options *ProcessTemplateMetaDataOptions) SetXGithubToken(xGithubToken string) *ProcessTemplateMetaDataOptions {
	_options.XGithubToken = core.StringPtr(xGithubToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ProcessTemplateMetaDataOptions) SetHeaders(param map[string]string) *ProcessTemplateMetaDataOptions {
	options.Headers = param
	return options
}

// RefreshWorkspaceCommandOptions : The RefreshWorkspaceCommand options.
type RefreshWorkspaceCommandOptions struct {
	// The ID of the workspace, for which you want to run a Schematics `refresh` job.  To find the ID of your workspace,
	// use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// The IAM delegated token for your IBM Cloud account.  This token is required for requests that are sent via the UI
	// only.
	DelegatedToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRefreshWorkspaceCommandOptions : Instantiate RefreshWorkspaceCommandOptions
func (*SchematicsV1) NewRefreshWorkspaceCommandOptions(wID string, refreshToken string) *RefreshWorkspaceCommandOptions {
	return &RefreshWorkspaceCommandOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (_options *RefreshWorkspaceCommandOptions) SetWID(wID string) *RefreshWorkspaceCommandOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *RefreshWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *RefreshWorkspaceCommandOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetDelegatedToken : Allow user to set DelegatedToken
func (_options *RefreshWorkspaceCommandOptions) SetDelegatedToken(delegatedToken string) *RefreshWorkspaceCommandOptions {
	_options.DelegatedToken = core.StringPtr(delegatedToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RefreshWorkspaceCommandOptions) SetHeaders(param map[string]string) *RefreshWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// ReplaceInventoryOptions : The ReplaceInventory options.
type ReplaceInventoryOptions struct {
	// Resource Inventory Id.  Use `GET /v2/inventories` API to look up the Resource Inventory definition Ids  in your IBM
	// Cloud account.
	InventoryID *string `json:"-" validate:"required,ne="`

	// The unique name of your Inventory definition. The name can be up to 128 characters long and can include alphanumeric
	// characters, spaces, dashes, and underscores.
	Name *string `json:"name,omitempty"`

	// The description of your Inventory definition. The description can be up to 2048 characters long in size.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Inventory definition.   By default, Inventory definition will be created in Default
	// Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Input inventory of host and host group for the playbook, in the `.ini` file format.
	InventoriesIni *string `json:"inventories_ini,omitempty"`

	// Input resource query definitions that is used to dynamically generate the inventory of host and host group for the
	// playbook.
	ResourceQueries []string `json:"resource_queries,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceInventoryOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	ReplaceInventoryOptions_Location_EuDe    = "eu-de"
	ReplaceInventoryOptions_Location_EuGb    = "eu-gb"
	ReplaceInventoryOptions_Location_UsEast  = "us-east"
	ReplaceInventoryOptions_Location_UsSouth = "us-south"
)

// NewReplaceInventoryOptions : Instantiate ReplaceInventoryOptions
func (*SchematicsV1) NewReplaceInventoryOptions(inventoryID string) *ReplaceInventoryOptions {
	return &ReplaceInventoryOptions{
		InventoryID: core.StringPtr(inventoryID),
	}
}

// SetInventoryID : Allow user to set InventoryID
func (_options *ReplaceInventoryOptions) SetInventoryID(inventoryID string) *ReplaceInventoryOptions {
	_options.InventoryID = core.StringPtr(inventoryID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceInventoryOptions) SetName(name string) *ReplaceInventoryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceInventoryOptions) SetDescription(description string) *ReplaceInventoryOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *ReplaceInventoryOptions) SetLocation(location string) *ReplaceInventoryOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *ReplaceInventoryOptions) SetResourceGroup(resourceGroup string) *ReplaceInventoryOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetInventoriesIni : Allow user to set InventoriesIni
func (_options *ReplaceInventoryOptions) SetInventoriesIni(inventoriesIni string) *ReplaceInventoryOptions {
	_options.InventoriesIni = core.StringPtr(inventoriesIni)
	return _options
}

// SetResourceQueries : Allow user to set ResourceQueries
func (_options *ReplaceInventoryOptions) SetResourceQueries(resourceQueries []string) *ReplaceInventoryOptions {
	_options.ResourceQueries = resourceQueries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceInventoryOptions) SetHeaders(param map[string]string) *ReplaceInventoryOptions {
	options.Headers = param
	return options
}

// ReplaceResourcesQueryOptions : The ReplaceResourcesQuery options.
type ReplaceResourcesQueryOptions struct {
	// Resource query Id.  Use `GET /v2/resource_query` API to look up the Resource query definition Ids  in your IBM Cloud
	// account.
	QueryID *string `json:"-" validate:"required,ne="`

	// Resource type (cluster, vsi, icd, vpc).
	Type *string `json:"type,omitempty"`

	// Resource query name.
	Name *string `json:"name,omitempty"`

	Queries []ResourceQuery `json:"queries,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceResourcesQueryOptions.Type property.
// Resource type (cluster, vsi, icd, vpc).
const (
	ReplaceResourcesQueryOptions_Type_Vsi = "vsi"
)

// NewReplaceResourcesQueryOptions : Instantiate ReplaceResourcesQueryOptions
func (*SchematicsV1) NewReplaceResourcesQueryOptions(queryID string) *ReplaceResourcesQueryOptions {
	return &ReplaceResourcesQueryOptions{
		QueryID: core.StringPtr(queryID),
	}
}

// SetQueryID : Allow user to set QueryID
func (_options *ReplaceResourcesQueryOptions) SetQueryID(queryID string) *ReplaceResourcesQueryOptions {
	_options.QueryID = core.StringPtr(queryID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceResourcesQueryOptions) SetType(typeVar string) *ReplaceResourcesQueryOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceResourcesQueryOptions) SetName(name string) *ReplaceResourcesQueryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetQueries : Allow user to set Queries
func (_options *ReplaceResourcesQueryOptions) SetQueries(queries []ResourceQuery) *ReplaceResourcesQueryOptions {
	_options.Queries = queries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceResourcesQueryOptions) SetHeaders(param map[string]string) *ReplaceResourcesQueryOptions {
	options.Headers = param
	return options
}

// ReplaceWorkspaceInputsOptions : The ReplaceWorkspaceInputs options.
type ReplaceWorkspaceInputsOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template in your workspace.  When you create a workspace, the Terraform template that  your
	// workspace points to is assigned a unique ID. Use the `GET /v1/workspaces` to look up the workspace IDs  and template
	// IDs or `template_data.id` in your IBM Cloud account.
	TID *string `json:"-" validate:"required,ne="`

	// A list of environment variables that you want to apply during the execution of a bash script or Terraform job. This
	// field must be provided as a list of key-value pairs, for example, **TF_LOG=debug**. Each entry will be a map with
	// one entry where `key is the environment variable name and value is value`. You can define environment variables for
	// IBM Cloud catalog offerings that are provisioned by using a bash script. See [example to use special environment
	// variable](https://cloud.ibm.com/docs/schematics?topic=schematics-set-parallelism#parallelism-example)  that are
	// supported by Schematics.
	EnvValues []interface{} `json:"env_values,omitempty"`

	// User values.
	Values *string `json:"values,omitempty"`

	// VariablesRequest -.
	Variablestore []WorkspaceVariableRequest `json:"variablestore,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceWorkspaceInputsOptions : Instantiate ReplaceWorkspaceInputsOptions
func (*SchematicsV1) NewReplaceWorkspaceInputsOptions(wID string, tID string) *ReplaceWorkspaceInputsOptions {
	return &ReplaceWorkspaceInputsOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *ReplaceWorkspaceInputsOptions) SetWID(wID string) *ReplaceWorkspaceInputsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *ReplaceWorkspaceInputsOptions) SetTID(tID string) *ReplaceWorkspaceInputsOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetEnvValues : Allow user to set EnvValues
func (_options *ReplaceWorkspaceInputsOptions) SetEnvValues(envValues []interface{}) *ReplaceWorkspaceInputsOptions {
	_options.EnvValues = envValues
	return _options
}

// SetValues : Allow user to set Values
func (_options *ReplaceWorkspaceInputsOptions) SetValues(values string) *ReplaceWorkspaceInputsOptions {
	_options.Values = core.StringPtr(values)
	return _options
}

// SetVariablestore : Allow user to set Variablestore
func (_options *ReplaceWorkspaceInputsOptions) SetVariablestore(variablestore []WorkspaceVariableRequest) *ReplaceWorkspaceInputsOptions {
	_options.Variablestore = variablestore
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceWorkspaceInputsOptions) SetHeaders(param map[string]string) *ReplaceWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// ReplaceWorkspaceOptions : The ReplaceWorkspace options.
type ReplaceWorkspaceOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Information about the software template that you chose from the IBM Cloud catalog. This information is returned for
	// IBM Cloud catalog offerings only.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// The description of the workspace.
	Description *string `json:"description,omitempty"`

	// The name of the workspace.
	Name *string `json:"name,omitempty"`

	// Information about the Target used by the templates originating from the  IBM Cloud catalog offerings. This
	// information is not relevant for workspace created using your own Terraform template.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// A list of tags that you want to associate with your workspace.
	Tags []string `json:"tags,omitempty"`

	// Input data for the Template.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// Input to update the template repository data.
	TemplateRepo *TemplateRepoUpdateRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// Input to update the workspace status.
	WorkspaceStatus *WorkspaceStatusUpdateRequest `json:"workspace_status,omitempty"`

	// Information about the last job that ran against the workspace. -.
	WorkspaceStatusMsg *WorkspaceStatusMessage `json:"workspace_status_msg,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceWorkspaceOptions : Instantiate ReplaceWorkspaceOptions
func (*SchematicsV1) NewReplaceWorkspaceOptions(wID string) *ReplaceWorkspaceOptions {
	return &ReplaceWorkspaceOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *ReplaceWorkspaceOptions) SetWID(wID string) *ReplaceWorkspaceOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetCatalogRef : Allow user to set CatalogRef
func (_options *ReplaceWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *ReplaceWorkspaceOptions {
	_options.CatalogRef = catalogRef
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceWorkspaceOptions) SetDescription(description string) *ReplaceWorkspaceOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceWorkspaceOptions) SetName(name string) *ReplaceWorkspaceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSharedData : Allow user to set SharedData
func (_options *ReplaceWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *ReplaceWorkspaceOptions {
	_options.SharedData = sharedData
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ReplaceWorkspaceOptions) SetTags(tags []string) *ReplaceWorkspaceOptions {
	_options.Tags = tags
	return _options
}

// SetTemplateData : Allow user to set TemplateData
func (_options *ReplaceWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *ReplaceWorkspaceOptions {
	_options.TemplateData = templateData
	return _options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (_options *ReplaceWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoUpdateRequest) *ReplaceWorkspaceOptions {
	_options.TemplateRepo = templateRepo
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceWorkspaceOptions) SetType(typeVar []string) *ReplaceWorkspaceOptions {
	_options.Type = typeVar
	return _options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (_options *ReplaceWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusUpdateRequest) *ReplaceWorkspaceOptions {
	_options.WorkspaceStatus = workspaceStatus
	return _options
}

// SetWorkspaceStatusMsg : Allow user to set WorkspaceStatusMsg
func (_options *ReplaceWorkspaceOptions) SetWorkspaceStatusMsg(workspaceStatusMsg *WorkspaceStatusMessage) *ReplaceWorkspaceOptions {
	_options.WorkspaceStatusMsg = workspaceStatusMsg
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceWorkspaceOptions) SetHeaders(param map[string]string) *ReplaceWorkspaceOptions {
	options.Headers = param
	return options
}

// ResourceGroupResponse : A list of resource groups that your account has access to.
type ResourceGroupResponse struct {
	// The ID of the account for which you listed the resource groups.
	AccountID *string `json:"account_id,omitempty"`

	// The CRN of the resource group.
	Crn *string `json:"crn,omitempty"`

	// If set to **true**, the resource group is used as the default resource group for your account. If set to **false**,
	// the resource group is not used as the default resource group in your account.
	Default *bool `json:"default,omitempty"`

	// The name of the resource group.
	Name *string `json:"name,omitempty"`

	// The ID of the resource group.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// The state of the resource group.
	State *string `json:"state,omitempty"`
}

// UnmarshalResourceGroupResponse unmarshals an instance of ResourceGroupResponse from the specified map of raw messages.
func UnmarshalResourceGroupResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupResponse)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQuery : Describe resource query.
type ResourceQuery struct {
	// Type of the query(workspaces).
	QueryType *string `json:"query_type,omitempty"`

	QueryCondition []ResourceQueryParam `json:"query_condition,omitempty"`

	// List of query selection parameters.
	QuerySelect []string `json:"query_select,omitempty"`
}

// Constants associated with the ResourceQuery.QueryType property.
// Type of the query(workspaces).
const (
	ResourceQuery_QueryType_Workspaces = "workspaces"
)

// UnmarshalResourceQuery unmarshals an instance of ResourceQuery from the specified map of raw messages.
func UnmarshalResourceQuery(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQuery)
	err = core.UnmarshalPrimitive(m, "query_type", &obj.QueryType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "query_condition", &obj.QueryCondition, UnmarshalResourceQueryParam)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "query_select", &obj.QuerySelect)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryParam : Describe resource query param.
type ResourceQueryParam struct {
	// Name of the resource query param.
	Name *string `json:"name,omitempty"`

	// Value of the resource query param.
	Value *string `json:"value,omitempty"`

	// Description of resource query param variable.
	Description *string `json:"description,omitempty"`
}

// UnmarshalResourceQueryParam unmarshals an instance of ResourceQueryParam from the specified map of raw messages.
func UnmarshalResourceQueryParam(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryParam)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryRecord : Describe resource query record.
type ResourceQueryRecord struct {
	// Resource type (cluster, vsi, icd, vpc).
	Type *string `json:"type,omitempty"`

	// Resource query name.
	Name *string `json:"name,omitempty"`

	// Resource Query id.
	ID *string `json:"id,omitempty"`

	// Resource query creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Email address of user who created the Resource query.
	CreatedBy *string `json:"created_by,omitempty"`

	// Resource query updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Email address of user who updated the Resource query.
	UpdatedBy *string `json:"updated_by,omitempty"`

	Queries []ResourceQuery `json:"queries,omitempty"`
}

// Constants associated with the ResourceQueryRecord.Type property.
// Resource type (cluster, vsi, icd, vpc).
const (
	ResourceQueryRecord_Type_Vsi = "vsi"
)

// UnmarshalResourceQueryRecord unmarshals an instance of ResourceQueryRecord from the specified map of raw messages.
func UnmarshalResourceQueryRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryRecord)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "queries", &obj.Queries, UnmarshalResourceQuery)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryRecordList : List of Resource query records.
type ResourceQueryRecordList struct {
	// Total number of records.
	TotalCount *int64 `json:"total_count,omitempty"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// List of resource query records.
	ResourceQueries []ResourceQueryRecord `json:"resource_queries,omitempty"`
}

// UnmarshalResourceQueryRecordList unmarshals an instance of ResourceQueryRecordList from the specified map of raw messages.
func UnmarshalResourceQueryRecordList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryRecordList)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_queries", &obj.ResourceQueries, UnmarshalResourceQueryRecord)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryResponseRecord : Describe resource query.
type ResourceQueryResponseRecord struct {
	Response []ResourceQueryResponseRecordResponse `json:"response,omitempty"`
}

// UnmarshalResourceQueryResponseRecord unmarshals an instance of ResourceQueryResponseRecord from the specified map of raw messages.
func UnmarshalResourceQueryResponseRecord(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryResponseRecord)
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalResourceQueryResponseRecordResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryResponseRecordQueryOutput : List of query output values.
type ResourceQueryResponseRecordQueryOutput struct {
	// Name of the output param.
	Name *string `json:"name,omitempty"`

	// value of the output param.
	Value *string `json:"value,omitempty"`
}

// UnmarshalResourceQueryResponseRecordQueryOutput unmarshals an instance of ResourceQueryResponseRecordQueryOutput from the specified map of raw messages.
func UnmarshalResourceQueryResponseRecordQueryOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryResponseRecordQueryOutput)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceQueryResponseRecordResponse : ResourceQueryResponseRecordResponse struct
type ResourceQueryResponseRecordResponse struct {
	// Type of the query(workspaces).
	QueryType *string `json:"query_type,omitempty"`

	QueryCondition []ResourceQueryParam `json:"query_condition,omitempty"`

	// List of query selection parameters.
	QuerySelect []string `json:"query_select,omitempty"`

	QueryOutput []ResourceQueryResponseRecordQueryOutput `json:"query_output,omitempty"`
}

// Constants associated with the ResourceQueryResponseRecordResponse.QueryType property.
// Type of the query(workspaces).
const (
	ResourceQueryResponseRecordResponse_QueryType_Workspaces = "workspaces"
)

// UnmarshalResourceQueryResponseRecordResponse unmarshals an instance of ResourceQueryResponseRecordResponse from the specified map of raw messages.
func UnmarshalResourceQueryResponseRecordResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceQueryResponseRecordResponse)
	err = core.UnmarshalPrimitive(m, "query_type", &obj.QueryType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "query_condition", &obj.QueryCondition, UnmarshalResourceQueryParam)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "query_select", &obj.QuerySelect)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "query_output", &obj.QueryOutput, UnmarshalResourceQueryResponseRecordQueryOutput)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunWorkspaceCommandsOptions : The RunWorkspaceCommands options.
type RunWorkspaceCommandsOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// List of commands.  You can execute single set of commands or multiple commands.  For more information, about the
	// payload of the multiple commands,  refer to
	// [Commands](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#commands).
	Commands []TerraformCommand `json:"commands,omitempty"`

	// Command name.
	OperationName *string `json:"operation_name,omitempty"`

	// Command description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRunWorkspaceCommandsOptions : Instantiate RunWorkspaceCommandsOptions
func (*SchematicsV1) NewRunWorkspaceCommandsOptions(wID string, refreshToken string) *RunWorkspaceCommandsOptions {
	return &RunWorkspaceCommandsOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (_options *RunWorkspaceCommandsOptions) SetWID(wID string) *RunWorkspaceCommandsOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *RunWorkspaceCommandsOptions) SetRefreshToken(refreshToken string) *RunWorkspaceCommandsOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetCommands : Allow user to set Commands
func (_options *RunWorkspaceCommandsOptions) SetCommands(commands []TerraformCommand) *RunWorkspaceCommandsOptions {
	_options.Commands = commands
	return _options
}

// SetOperationName : Allow user to set OperationName
func (_options *RunWorkspaceCommandsOptions) SetOperationName(operationName string) *RunWorkspaceCommandsOptions {
	_options.OperationName = core.StringPtr(operationName)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *RunWorkspaceCommandsOptions) SetDescription(description string) *RunWorkspaceCommandsOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RunWorkspaceCommandsOptions) SetHeaders(param map[string]string) *RunWorkspaceCommandsOptions {
	options.Headers = param
	return options
}

// SchematicsLocations : Information about the location.
type SchematicsLocations struct {
	// The name of the location.
	Name *string `json:"name,omitempty"`

	// The ID of the location.
	ID *string `json:"id,omitempty"`

	// The country where the location is located.
	Country *string `json:"country,omitempty"`

	// The geography that the location belongs to.
	Geography *string `json:"geography,omitempty"`

	// Geographical continent locations code having the data centres of IBM Cloud Schematics service.
	GeographyCode *string `json:"geography_code,omitempty"`

	// The metro area that the location belongs to.
	Metro *string `json:"metro,omitempty"`

	// The multizone metro area that the location belongs to.
	MultizoneMetro *string `json:"multizone_metro,omitempty"`

	// The kind of location.
	Kind *string `json:"kind,omitempty"`

	// List of paired regions used by Schematics.
	PairedRegion []string `json:"paired_region,omitempty"`

	// Restricted Region.
	Restricted *bool `json:"restricted,omitempty"`
}

// UnmarshalSchematicsLocations unmarshals an instance of SchematicsLocations from the specified map of raw messages.
func UnmarshalSchematicsLocations(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsLocations)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "geography", &obj.Geography)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "geography_code", &obj.GeographyCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metro", &obj.Metro)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "multizone_metro", &obj.MultizoneMetro)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paired_region", &obj.PairedRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "restricted", &obj.Restricted)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsLocationsList : List of Locations details.
type SchematicsLocationsList struct {
	// List of Locations.
	Locations []SchematicsLocationsLite `json:"locations,omitempty"`
}

// UnmarshalSchematicsLocationsList unmarshals an instance of SchematicsLocationsList from the specified map of raw messages.
func UnmarshalSchematicsLocationsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsLocationsList)
	err = core.UnmarshalModel(m, "locations", &obj.Locations, UnmarshalSchematicsLocationsLite)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsLocationsLite : individual location details.
type SchematicsLocationsLite struct {
	// Geographical Region code having the data centres of IBM Cloud Schematics service.
	Region *string `json:"region,omitempty"`

	// Geographical city locations having the data centres of IBM Cloud Schematics service.
	Metro *string `json:"metro,omitempty"`

	// Geographical continent locations code having the data centres of IBM Cloud Schematics service.
	GeographyCode *string `json:"geography_code,omitempty"`

	// Geographical continent locations having the data centres of IBM Cloud Schematics service.
	Geography *string `json:"geography,omitempty"`

	// Country locations having the data centres of IBM Cloud Schematics service.
	Country *string `json:"country,omitempty"`

	// The kind of location.
	Kind *string `json:"kind,omitempty"`

	// List of paired regions used by Schematics.
	PairedRegion []string `json:"paired_region,omitempty"`

	// Restricted Region.
	Restricted *bool `json:"restricted,omitempty"`
}

// UnmarshalSchematicsLocationsLite unmarshals an instance of SchematicsLocationsLite from the specified map of raw messages.
func UnmarshalSchematicsLocationsLite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsLocationsLite)
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "metro", &obj.Metro)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "geography_code", &obj.GeographyCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "geography", &obj.Geography)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paired_region", &obj.PairedRegion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "restricted", &obj.Restricted)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedTargetData : Information about the Target used by the templates originating from the  IBM Cloud catalog offerings. This
// information is not relevant for workspace created using your own Terraform template.
type SharedTargetData struct {
	// Cluster created on.
	ClusterCreatedOn *string `json:"cluster_created_on,omitempty"`

	// The ID of the cluster where you want to provision the resources of all IBM Cloud catalog templates that are included
	// in the catalog offering.
	ClusterID *string `json:"cluster_id,omitempty"`

	// The cluster name.
	ClusterName *string `json:"cluster_name,omitempty"`

	// The cluster type.
	ClusterType *string `json:"cluster_type,omitempty"`

	// The entitlement key that you want to use to install IBM Cloud entitled software.
	EntitlementKeys []interface{} `json:"entitlement_keys,omitempty"`

	// The Kubernetes namespace or OpenShift project where the resources of all IBM Cloud catalog templates that are
	// included in the catalog offering are deployed into.
	Namespace *string `json:"namespace,omitempty"`

	// The IBM Cloud region that you want to use for the resources of all IBM Cloud catalog templates that are included in
	// the catalog offering.
	Region *string `json:"region,omitempty"`

	// The ID of the resource group that you want to use for the resources of all IBM Cloud catalog templates that are
	// included in the catalog offering.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// The cluster worker count.
	WorkerCount *int64 `json:"worker_count,omitempty"`

	// The cluster worker type.
	WorkerMachineType *string `json:"worker_machine_type,omitempty"`
}

// UnmarshalSharedTargetData unmarshals an instance of SharedTargetData from the specified map of raw messages.
func UnmarshalSharedTargetData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SharedTargetData)
	err = core.UnmarshalPrimitive(m, "cluster_created_on", &obj.ClusterCreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_id", &obj.ClusterID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_name", &obj.ClusterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_type", &obj.ClusterType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entitlement_keys", &obj.EntitlementKeys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "worker_count", &obj.WorkerCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "worker_machine_type", &obj.WorkerMachineType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedTargetDataResponse : Information about the Target used by the templates originating from IBM Cloud catalog offerings. This information is
// not relevant when you create a workspace from your own Terraform template.
type SharedTargetDataResponse struct {
	// The ID of the cluster where you want to provision the resources of all IBM Cloud catalog templates that are included
	// in the catalog offering.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Target cluster name.
	ClusterName *string `json:"cluster_name,omitempty"`

	// The entitlement key that you want to use to install IBM Cloud entitled software.
	EntitlementKeys []interface{} `json:"entitlement_keys,omitempty"`

	// The Kubernetes namespace or OpenShift project where the resources of all IBM Cloud catalog templates that are
	// included in the catalog offering are deployed into.
	Namespace *string `json:"namespace,omitempty"`

	// The IBM Cloud region that you want to use for the resources of all IBM Cloud catalog templates that are included in
	// the catalog offering.
	Region *string `json:"region,omitempty"`

	// The ID of the resource group that you want to use for the resources of all IBM Cloud catalog templates that are
	// included in the catalog offering.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`
}

// UnmarshalSharedTargetDataResponse unmarshals an instance of SharedTargetDataResponse from the specified map of raw messages.
func UnmarshalSharedTargetDataResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SharedTargetDataResponse)
	err = core.UnmarshalPrimitive(m, "cluster_id", &obj.ClusterID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cluster_name", &obj.ClusterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "entitlement_keys", &obj.EntitlementKeys)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StateStoreResponse : Information about workspace runtime data.
type StateStoreResponse struct {
	// The provisioning engine that was used to apply the Terraform template or IBM Cloud catalog software template.
	EngineName *string `json:"engine_name,omitempty"`

	// The version of the provisioning engine that was used.
	EngineVersion *string `json:"engine_version,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// The URL where the Terraform statefile (`terraform.tfstate`) is stored. You can use the statefile to find an overview
	// of IBM Cloud resources that were created by Schematics. Schematics uses the statefile as an inventory list to
	// determine future create, update, or deletion jobs.
	StateStoreURL *string `json:"state_store_url,omitempty"`
}

// UnmarshalStateStoreResponse unmarshals an instance of StateStoreResponse from the specified map of raw messages.
func UnmarshalStateStoreResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StateStoreResponse)
	err = core.UnmarshalPrimitive(m, "engine_name", &obj.EngineName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "engine_version", &obj.EngineVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_store_url", &obj.StateStoreURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StateStoreResponseList : Information about the Terraform statefile URL.
type StateStoreResponseList struct {
	// Information about workspace runtime data.
	RuntimeData []StateStoreResponse `json:"runtime_data,omitempty"`
}

// UnmarshalStateStoreResponseList unmarshals an instance of StateStoreResponseList from the specified map of raw messages.
func UnmarshalStateStoreResponseList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StateStoreResponseList)
	err = core.UnmarshalModel(m, "runtime_data", &obj.RuntimeData, UnmarshalStateStoreResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SystemLock : System lock status.
type SystemLock struct {
	// Is the automation locked by a Schematic job ?.
	SysLocked *bool `json:"sys_locked,omitempty"`

	// Name of the User who performed the job, that lead to the locking of the automation.
	SysLockedBy *string `json:"sys_locked_by,omitempty"`

	// When the User performed the job that lead to locking of the automation ?.
	SysLockedAt *strfmt.DateTime `json:"sys_locked_at,omitempty"`
}

// UnmarshalSystemLock unmarshals an instance of SystemLock from the specified map of raw messages.
func UnmarshalSystemLock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SystemLock)
	err = core.UnmarshalPrimitive(m, "sys_locked", &obj.SysLocked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sys_locked_by", &obj.SysLockedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sys_locked_at", &obj.SysLockedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateMetaDataResponse : Template metadata response.
type TemplateMetaDataResponse struct {
	// Template type (terraform, ansible, helm, cloudpak, bash script).
	Type *string `json:"type,omitempty"`

	// List of variables and its metadata.
	Variables []VariableData `json:"variables" validate:"required"`
}

// UnmarshalTemplateMetaDataResponse unmarshals an instance of TemplateMetaDataResponse from the specified map of raw messages.
func UnmarshalTemplateMetaDataResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateMetaDataResponse)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "variables", &obj.Variables, UnmarshalVariableData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateReadme : The `README.md` file for the template used by the workspace.
type TemplateReadme struct {
	// The `README.md` file for the template used by the workspace.
	Readme *string `json:"readme,omitempty"`
}

// UnmarshalTemplateReadme unmarshals an instance of TemplateReadme from the specified map of raw messages.
func UnmarshalTemplateReadme(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateReadme)
	err = core.UnmarshalPrimitive(m, "readme", &obj.Readme)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRepoRequest : Input variables for the Template repoository, while creating a workspace.
type TemplateRepoRequest struct {
	// The repository branch.
	Branch *string `json:"branch,omitempty"`

	// The repository release.
	Release *string `json:"release,omitempty"`

	// The repository SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// The repository URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// The source URL.
	URL *string `json:"url,omitempty"`
}

// UnmarshalTemplateRepoRequest unmarshals an instance of TemplateRepoRequest from the specified map of raw messages.
func UnmarshalTemplateRepoRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateRepoRequest)
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "release", &obj.Release)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_sha_value", &obj.RepoShaValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_url", &obj.RepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRepoResponse : Information about the Template repository used by the workspace.
type TemplateRepoResponse struct {
	// The repository branch.
	Branch *string `json:"branch,omitempty"`

	// Full repository URL.
	FullURL *string `json:"full_url,omitempty"`

	// Has uploaded Git repository tar.
	HasUploadedgitrepotar *bool `json:"has_uploadedgitrepotar,omitempty"`

	// The repository release.
	Release *string `json:"release,omitempty"`

	// The repository SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// The repository URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// The source URL.
	URL *string `json:"url,omitempty"`
}

// UnmarshalTemplateRepoResponse unmarshals an instance of TemplateRepoResponse from the specified map of raw messages.
func UnmarshalTemplateRepoResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateRepoResponse)
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "full_url", &obj.FullURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "has_uploadedgitrepotar", &obj.HasUploadedgitrepotar)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "release", &obj.Release)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_sha_value", &obj.RepoShaValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_url", &obj.RepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRepoTarUploadResponse : TemplateRepoTarUploadResponse -.
type TemplateRepoTarUploadResponse struct {
	// Tar file value.
	FileValue *string `json:"file_value,omitempty"`

	// Has received tar file?.
	HasReceivedFile *bool `json:"has_received_file,omitempty"`

	// Template ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalTemplateRepoTarUploadResponse unmarshals an instance of TemplateRepoTarUploadResponse from the specified map of raw messages.
func UnmarshalTemplateRepoTarUploadResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateRepoTarUploadResponse)
	err = core.UnmarshalPrimitive(m, "file_value", &obj.FileValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "has_received_file", &obj.HasReceivedFile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRepoUpdateRequest : Input to update the template repository data.
type TemplateRepoUpdateRequest struct {
	// The repository branch.
	Branch *string `json:"branch,omitempty"`

	// The repository release.
	Release *string `json:"release,omitempty"`

	// The repository SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// The repository URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// The source URL.
	URL *string `json:"url,omitempty"`
}

// UnmarshalTemplateRepoUpdateRequest unmarshals an instance of TemplateRepoUpdateRequest from the specified map of raw messages.
func UnmarshalTemplateRepoUpdateRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateRepoUpdateRequest)
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "release", &obj.Release)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_sha_value", &obj.RepoShaValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "repo_url", &obj.RepoURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRepoUploadOptions : The TemplateRepoUpload options.
type TemplateRepoUploadOptions struct {
	// The ID of the workspace where you want to upload your `.tar` file. To find the workspace ID, use the `GET
	// /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// The ID of the Terraform template in your workspace. When you create a workspace, a unique ID is assigned to your
	// Terraform template, even if no template was provided during workspace creation. To find this ID, use the `GET
	// /v1/workspaces` API and review the `template_data.id` value.
	TID *string `json:"-" validate:"required,ne="`

	// Template tar file.
	File io.ReadCloser `json:"-"`

	// The content type of file.
	FileContentType *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewTemplateRepoUploadOptions : Instantiate TemplateRepoUploadOptions
func (*SchematicsV1) NewTemplateRepoUploadOptions(wID string, tID string) *TemplateRepoUploadOptions {
	return &TemplateRepoUploadOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (_options *TemplateRepoUploadOptions) SetWID(wID string) *TemplateRepoUploadOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetTID : Allow user to set TID
func (_options *TemplateRepoUploadOptions) SetTID(tID string) *TemplateRepoUploadOptions {
	_options.TID = core.StringPtr(tID)
	return _options
}

// SetFile : Allow user to set File
func (_options *TemplateRepoUploadOptions) SetFile(file io.ReadCloser) *TemplateRepoUploadOptions {
	_options.File = file
	return _options
}

// SetFileContentType : Allow user to set FileContentType
func (_options *TemplateRepoUploadOptions) SetFileContentType(fileContentType string) *TemplateRepoUploadOptions {
	_options.FileContentType = core.StringPtr(fileContentType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *TemplateRepoUploadOptions) SetHeaders(param map[string]string) *TemplateRepoUploadOptions {
	options.Headers = param
	return options
}

// TemplateResources : Information about the resources provisioned by the Terraform template.
type TemplateResources struct {
	// The subfolder in GitHub or GitLab where your Terraform templates are stored.  If your template is stored in the root
	// directory, `.` is returned.
	Folder *string `json:"folder,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// List of null resources.
	NullResources []interface{} `json:"null_resources,omitempty"`

	// Information about the IBM Cloud resources that are associated with your workspace.
	RelatedResources []interface{} `json:"related_resources,omitempty"`

	// Information about the IBM Cloud resources that are associated with your workspace. **Note** The `resource_tainted`
	// flag marks `true` when an instance is times out after few hours, if your resource provisioning takes longer
	// duration. When you rerun the apply plan, based on the `resource_taint` flag result the provisioning continues from
	// the state where the provisioning has stopped.
	Resources []interface{} `json:"resources,omitempty"`

	// Number of resources.
	ResourcesCount *int64 `json:"resources_count,omitempty"`

	// The Terraform version that was used to apply your template.
	TemplateType *string `json:"template_type,omitempty"`
}

// UnmarshalTemplateResources unmarshals an instance of TemplateResources from the specified map of raw messages.
func UnmarshalTemplateResources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateResources)
	err = core.UnmarshalPrimitive(m, "folder", &obj.Folder)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "null_resources", &obj.NullResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "related_resources", &obj.RelatedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources", &obj.Resources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_count", &obj.ResourcesCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_type", &obj.TemplateType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateRunTimeDataResponse : Information about the provisioning engine, state file, and runtime logs.
type TemplateRunTimeDataResponse struct {
	// The command that was used to apply the Terraform template or IBM Cloud catalog software template.
	EngineCmd *string `json:"engine_cmd,omitempty"`

	// The provisioning engine that was used to apply the Terraform template or IBM Cloud catalog software template.
	EngineName *string `json:"engine_name,omitempty"`

	// The version of the provisioning engine that was used.
	EngineVersion *string `json:"engine_version,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// The URL to access the logs that were created during the creation, update, or deletion of your IBM Cloud resources.
	LogStoreURL *string `json:"log_store_url,omitempty"`

	// List of Output values.
	OutputValues []interface{} `json:"output_values,omitempty"`

	// List of resources.
	Resources [][]interface{} `json:"resources,omitempty"`

	// The URL where the Terraform statefile (`terraform.tfstate`) is stored. You can use the statefile to find an overview
	// of IBM Cloud resources that were created by Schematics. Schematics uses the statefile as an inventory list to
	// determine future create, update, or deletion jobs.
	StateStoreURL *string `json:"state_store_url,omitempty"`
}

// UnmarshalTemplateRunTimeDataResponse unmarshals an instance of TemplateRunTimeDataResponse from the specified map of raw messages.
func UnmarshalTemplateRunTimeDataResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateRunTimeDataResponse)
	err = core.UnmarshalPrimitive(m, "engine_cmd", &obj.EngineCmd)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "engine_name", &obj.EngineName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "engine_version", &obj.EngineVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_store_url", &obj.LogStoreURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "output_values", &obj.OutputValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources", &obj.Resources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state_store_url", &obj.StateStoreURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateSourceDataRequest : Input parameters to define input variables for your Terraform template.
type TemplateSourceDataRequest struct {
	// A list of environment variables that you want to apply during the execution of a bash script or Terraform job. This
	// field must be provided as a list of key-value pairs, for example, **TF_LOG=debug**. Each entry will be a map with
	// one entry where `key is the environment variable name and value is value`. You can define environment variables for
	// IBM Cloud catalog offerings that are provisioned by using a bash script. See [example to use special environment
	// variable](https://cloud.ibm.com/docs/schematics?topic=schematics-set-parallelism#parallelism-example)  that are
	// supported by Schematics.
	EnvValues []interface{} `json:"env_values,omitempty"`

	// The subfolder in your GitHub or GitLab repository where your Terraform template is stored.
	Folder *string `json:"folder,omitempty"`

	// True, to use the files from the specified folder & subfolder in your GitHub or GitLab repository and ignore the
	// other folders in the repository.
	Compact *bool `json:"compact,omitempty"`

	// The content of an existing Terraform statefile that you want to import in to your workspace. To get the content of a
	// Terraform statefile for a specific Terraform template in an existing workspace, run `ibmcloud terraform state pull
	// --id <workspace_id> --template <template_id>`.
	InitStateFile *string `json:"init_state_file,omitempty"`

	// The Terraform version that you want to use to run your Terraform code. Enter `terraform_v0.12` to use Terraform
	// version 0.12, and `terraform_v0.11` to use Terraform version 0.11. The Terraform config files are run with Terraform
	// version 0.11. This is a required variable. Make sure that your Terraform config files are compatible with the
	// Terraform version that you select.
	Type *string `json:"type,omitempty"`

	// Uninstall script name.
	UninstallScriptName *string `json:"uninstall_script_name,omitempty"`

	// A list of variable values that you want to apply during the Helm chart installation. The list must be provided in
	// JSON format, such as `"autoscaling: enabled: true minReplicas: 2"`. The values that you define here override the
	// default Helm chart values. This field is supported only for IBM Cloud catalog offerings that are provisioned by
	// using the Terraform Helm provider.
	Values *string `json:"values,omitempty"`

	// List of values metadata.
	ValuesMetadata []interface{} `json:"values_metadata,omitempty"`

	// VariablesRequest -.
	Variablestore []WorkspaceVariableRequest `json:"variablestore,omitempty"`
}

// UnmarshalTemplateSourceDataRequest unmarshals an instance of TemplateSourceDataRequest from the specified map of raw messages.
func UnmarshalTemplateSourceDataRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateSourceDataRequest)
	err = core.UnmarshalPrimitive(m, "env_values", &obj.EnvValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "folder", &obj.Folder)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compact", &obj.Compact)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "init_state_file", &obj.InitStateFile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uninstall_script_name", &obj.UninstallScriptName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values_metadata", &obj.ValuesMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "variablestore", &obj.Variablestore, UnmarshalWorkspaceVariableRequest)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateSourceDataResponse : Information about the input variables that are used in the template.
type TemplateSourceDataResponse struct {
	// List of environment values.
	EnvValues []EnvVariableResponse `json:"env_values,omitempty"`

	// The subfolder in your GitHub or GitLab repository where your Terraform template is stored. If your template is
	// stored in the root directory, `.` is returned.
	Folder *string `json:"folder,omitempty"`

	// True, to use the files from the specified folder & subfolder in your GitHub or GitLab repository and ignore the
	// other folders in the repository.
	Compact *bool `json:"compact,omitempty"`

	// Has github token.
	HasGithubtoken *bool `json:"has_githubtoken,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	ID *string `json:"id,omitempty"`

	// The Terraform version that was used to run your Terraform code.
	Type *string `json:"type,omitempty"`

	// Uninstall script name.
	UninstallScriptName *string `json:"uninstall_script_name,omitempty"`

	// A list of variable values that you want to apply during the Helm chart installation. The list must be provided in
	// JSON format, such as `"autoscaling: enabled: true minReplicas: 2"`. The values that you define here override the
	// default Helm chart values. This field is supported only for IBM Cloud catalog offerings that are provisioned by
	// using the Terraform Helm provider.
	Values *string `json:"values,omitempty"`

	// A list of input variables that are associated with the workspace.
	ValuesMetadata []interface{} `json:"values_metadata,omitempty"`

	// The API endpoint to access the input variables that you defined for your template.
	ValuesURL *string `json:"values_url,omitempty"`

	// Information about the input variables that your template uses.
	Variablestore []WorkspaceVariableResponse `json:"variablestore,omitempty"`
}

// UnmarshalTemplateSourceDataResponse unmarshals an instance of TemplateSourceDataResponse from the specified map of raw messages.
func UnmarshalTemplateSourceDataResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateSourceDataResponse)
	err = core.UnmarshalModel(m, "env_values", &obj.EnvValues, UnmarshalEnvVariableResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "folder", &obj.Folder)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compact", &obj.Compact)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "has_githubtoken", &obj.HasGithubtoken)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uninstall_script_name", &obj.UninstallScriptName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values_metadata", &obj.ValuesMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values_url", &obj.ValuesURL)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "variablestore", &obj.Variablestore, UnmarshalWorkspaceVariableResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateValues : Information about the input variables that are declared in the template that your workspace points to.
type TemplateValues struct {
	ValuesMetadata []interface{} `json:"values_metadata,omitempty"`
}

// UnmarshalTemplateValues unmarshals an instance of TemplateValues from the specified map of raw messages.
func UnmarshalTemplateValues(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateValues)
	err = core.UnmarshalPrimitive(m, "values_metadata", &obj.ValuesMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TerraformCommand : Inputs for running a Terraform command on the workspace.
type TerraformCommand struct {
	// You must provide the command to execute. Supported commands are `show`,`taint`, `untaint`, `state`, `import`,
	// `output`.
	Command *string `json:"command,omitempty"`

	// The required address parameters for the command name. You can send the option flag and address parameter in the
	// payload. **Syntax ** "command_params": "<option>=<flag>", "<address>". **Example ** "command_params":
	// "-allow-missing=true", "-lock=true", "data.template_file.test".
	CommandParams *string `json:"command_params,omitempty"`

	// The optional name for the command block.
	CommandName *string `json:"command_name,omitempty"`

	// The optional text to describe the command block.
	CommandDesc *string `json:"command_desc,omitempty"`

	// Instruction to continue or break in case of error.
	CommandOnError *string `json:"command_on_error,omitempty"`

	// Dependency on previous commands.
	CommandDependsOn *string `json:"command_depends_on,omitempty"`

	// Displays the command executed status, either `success` or `failure`.
	CommandStatus *string `json:"command_status,omitempty"`
}

// UnmarshalTerraformCommand unmarshals an instance of TerraformCommand from the specified map of raw messages.
func UnmarshalTerraformCommand(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TerraformCommand)
	err = core.UnmarshalPrimitive(m, "command", &obj.Command)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_params", &obj.CommandParams)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_name", &obj.CommandName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_desc", &obj.CommandDesc)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_on_error", &obj.CommandOnError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_depends_on", &obj.CommandDependsOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_status", &obj.CommandStatus)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateActionOptions : The UpdateAction options.
type UpdateActionOptions struct {
	// Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.
	ActionID *string `json:"-" validate:"required,ne="`

	// The unique name of your action. The name can be up to 128 characters long and can include alphanumeric characters,
	// spaces, dashes, and underscores. **Example** you can use the name to stop action.
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for an action.  By default, action is created in default resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the `README` file, for the source URL.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Target inventory record ID, used by the action or ansible playbook.
	Inventory *string `json:"inventory,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// User editable variable data & system generated reference to value.
	BastionCredential *VariableData `json:"bastion_credential,omitempty"`

	// Inventory of host and host group for the playbook in `INI` file format. For example, `"targets_ini":
	// "[webserverhost]
	//  172.22.192.6
	//  [dbhost]
	//  172.22.192.5"`. For more information, about an inventory host group syntax, see [Inventory host
	// groups](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).
	TargetsIni *string `json:"targets_ini,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform
	// template.
	XGithubToken *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateActionOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	UpdateActionOptions_Location_EuDe    = "eu-de"
	UpdateActionOptions_Location_EuGb    = "eu-gb"
	UpdateActionOptions_Location_UsEast  = "us-east"
	UpdateActionOptions_Location_UsSouth = "us-south"
)

// Constants associated with the UpdateActionOptions.SourceType property.
// Type of source for the Template.
const (
	UpdateActionOptions_SourceType_CosBucket        = "cos_bucket"
	UpdateActionOptions_SourceType_ExternalScm      = "external_scm"
	UpdateActionOptions_SourceType_GitHub           = "git_hub"
	UpdateActionOptions_SourceType_GitHubEnterprise = "git_hub_enterprise"
	UpdateActionOptions_SourceType_GitLab           = "git_lab"
	UpdateActionOptions_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	UpdateActionOptions_SourceType_IbmGitLab        = "ibm_git_lab"
	UpdateActionOptions_SourceType_Local            = "local"
)

// NewUpdateActionOptions : Instantiate UpdateActionOptions
func (*SchematicsV1) NewUpdateActionOptions(actionID string) *UpdateActionOptions {
	return &UpdateActionOptions{
		ActionID: core.StringPtr(actionID),
	}
}

// SetActionID : Allow user to set ActionID
func (_options *UpdateActionOptions) SetActionID(actionID string) *UpdateActionOptions {
	_options.ActionID = core.StringPtr(actionID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateActionOptions) SetName(name string) *UpdateActionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateActionOptions) SetDescription(description string) *UpdateActionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *UpdateActionOptions) SetLocation(location string) *UpdateActionOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *UpdateActionOptions) SetResourceGroup(resourceGroup string) *UpdateActionOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateActionOptions) SetTags(tags []string) *UpdateActionOptions {
	_options.Tags = tags
	return _options
}

// SetUserState : Allow user to set UserState
func (_options *UpdateActionOptions) SetUserState(userState *UserState) *UpdateActionOptions {
	_options.UserState = userState
	return _options
}

// SetSourceReadmeURL : Allow user to set SourceReadmeURL
func (_options *UpdateActionOptions) SetSourceReadmeURL(sourceReadmeURL string) *UpdateActionOptions {
	_options.SourceReadmeURL = core.StringPtr(sourceReadmeURL)
	return _options
}

// SetSource : Allow user to set Source
func (_options *UpdateActionOptions) SetSource(source *ExternalSource) *UpdateActionOptions {
	_options.Source = source
	return _options
}

// SetSourceType : Allow user to set SourceType
func (_options *UpdateActionOptions) SetSourceType(sourceType string) *UpdateActionOptions {
	_options.SourceType = core.StringPtr(sourceType)
	return _options
}

// SetCommandParameter : Allow user to set CommandParameter
func (_options *UpdateActionOptions) SetCommandParameter(commandParameter string) *UpdateActionOptions {
	_options.CommandParameter = core.StringPtr(commandParameter)
	return _options
}

// SetInventory : Allow user to set Inventory
func (_options *UpdateActionOptions) SetInventory(inventory string) *UpdateActionOptions {
	_options.Inventory = core.StringPtr(inventory)
	return _options
}

// SetCredentials : Allow user to set Credentials
func (_options *UpdateActionOptions) SetCredentials(credentials []VariableData) *UpdateActionOptions {
	_options.Credentials = credentials
	return _options
}

// SetBastion : Allow user to set Bastion
func (_options *UpdateActionOptions) SetBastion(bastion *BastionResourceDefinition) *UpdateActionOptions {
	_options.Bastion = bastion
	return _options
}

// SetBastionCredential : Allow user to set BastionCredential
func (_options *UpdateActionOptions) SetBastionCredential(bastionCredential *VariableData) *UpdateActionOptions {
	_options.BastionCredential = bastionCredential
	return _options
}

// SetTargetsIni : Allow user to set TargetsIni
func (_options *UpdateActionOptions) SetTargetsIni(targetsIni string) *UpdateActionOptions {
	_options.TargetsIni = core.StringPtr(targetsIni)
	return _options
}

// SetInputs : Allow user to set Inputs
func (_options *UpdateActionOptions) SetInputs(inputs []VariableData) *UpdateActionOptions {
	_options.Inputs = inputs
	return _options
}

// SetOutputs : Allow user to set Outputs
func (_options *UpdateActionOptions) SetOutputs(outputs []VariableData) *UpdateActionOptions {
	_options.Outputs = outputs
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *UpdateActionOptions) SetSettings(settings []VariableData) *UpdateActionOptions {
	_options.Settings = settings
	return _options
}

// SetState : Allow user to set State
func (_options *UpdateActionOptions) SetState(state *ActionState) *UpdateActionOptions {
	_options.State = state
	return _options
}

// SetSysLock : Allow user to set SysLock
func (_options *UpdateActionOptions) SetSysLock(sysLock *SystemLock) *UpdateActionOptions {
	_options.SysLock = sysLock
	return _options
}

// SetXGithubToken : Allow user to set XGithubToken
func (_options *UpdateActionOptions) SetXGithubToken(xGithubToken string) *UpdateActionOptions {
	_options.XGithubToken = core.StringPtr(xGithubToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateActionOptions) SetHeaders(param map[string]string) *UpdateActionOptions {
	options.Headers = param
	return options
}

// UpdateInventoryOptions : The UpdateInventory options.
type UpdateInventoryOptions struct {
	// Resource Inventory Id.  Use `GET /v2/inventories` API to look up the Resource Inventory definition Ids  in your IBM
	// Cloud account.
	InventoryID *string `json:"-" validate:"required,ne="`

	// The unique name of your Inventory definition. The name can be up to 128 characters long and can include alphanumeric
	// characters, spaces, dashes, and underscores.
	Name *string `json:"name,omitempty"`

	// The description of your Inventory definition. The description can be up to 2048 characters long in size.
	Description *string `json:"description,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Inventory definition.   By default, Inventory definition will be created in Default
	// Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Input inventory of host and host group for the playbook, in the `.ini` file format.
	InventoriesIni *string `json:"inventories_ini,omitempty"`

	// Input resource query definitions that is used to dynamically generate the inventory of host and host group for the
	// playbook.
	ResourceQueries []string `json:"resource_queries,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateInventoryOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	UpdateInventoryOptions_Location_EuDe    = "eu-de"
	UpdateInventoryOptions_Location_EuGb    = "eu-gb"
	UpdateInventoryOptions_Location_UsEast  = "us-east"
	UpdateInventoryOptions_Location_UsSouth = "us-south"
)

// NewUpdateInventoryOptions : Instantiate UpdateInventoryOptions
func (*SchematicsV1) NewUpdateInventoryOptions(inventoryID string) *UpdateInventoryOptions {
	return &UpdateInventoryOptions{
		InventoryID: core.StringPtr(inventoryID),
	}
}

// SetInventoryID : Allow user to set InventoryID
func (_options *UpdateInventoryOptions) SetInventoryID(inventoryID string) *UpdateInventoryOptions {
	_options.InventoryID = core.StringPtr(inventoryID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateInventoryOptions) SetName(name string) *UpdateInventoryOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateInventoryOptions) SetDescription(description string) *UpdateInventoryOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *UpdateInventoryOptions) SetLocation(location string) *UpdateInventoryOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *UpdateInventoryOptions) SetResourceGroup(resourceGroup string) *UpdateInventoryOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetInventoriesIni : Allow user to set InventoriesIni
func (_options *UpdateInventoryOptions) SetInventoriesIni(inventoriesIni string) *UpdateInventoryOptions {
	_options.InventoriesIni = core.StringPtr(inventoriesIni)
	return _options
}

// SetResourceQueries : Allow user to set ResourceQueries
func (_options *UpdateInventoryOptions) SetResourceQueries(resourceQueries []string) *UpdateInventoryOptions {
	_options.ResourceQueries = resourceQueries
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateInventoryOptions) SetHeaders(param map[string]string) *UpdateInventoryOptions {
	options.Headers = param
	return options
}

// UpdateJobOptions : The UpdateJob options.
type UpdateJobOptions struct {
	// Job Id. Use `GET /v2/jobs` API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"-" validate:"required,ne="`

	// The IAM refresh token for the user or service identity.
	//
	//   **Retrieving refresh token**:
	//   * Use `export IBMCLOUD_API_KEY=<ibmcloud_api_key>`, and execute `curl -X POST
	// "https://iam.cloud.ibm.com/identity/token" -H "Content-Type: application/x-www-form-urlencoded" -d
	// "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$IBMCLOUD_API_KEY" -u bx:bx`.
	//   * For more information, about creating IAM access token and API Docs, refer, [IAM access
	// token](/apidocs/iam-identity-token-api#gettoken-password) and [Create API
	// key](/apidocs/iam-identity-token-api#create-api-key).
	//
	//   **Limitation**:
	//   * If the token is expired, you can use `refresh token` to get a new IAM access token.
	//   * The `refresh_token` parameter cannot be used to retrieve a new IAM access token.
	//   * When the IAM access token is about to expire, use the API key to create a new access token.
	RefreshToken *string `json:"-" validate:"required"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action or Workspace.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action or Workspace.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
	// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
	// provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job data.
	Data *JobData `json:"data,omitempty"`

	// Describes a bastion resource.
	Bastion *BastionResourceDefinition `json:"bastion,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateJobOptions.CommandObject property.
// Name of the Schematics automation resource.
const (
	UpdateJobOptions_CommandObject_Action      = "action"
	UpdateJobOptions_CommandObject_Environment = "environment"
	UpdateJobOptions_CommandObject_System      = "system"
	UpdateJobOptions_CommandObject_Workspace   = "workspace"
)

// Constants associated with the UpdateJobOptions.CommandName property.
// Schematics job command name.
const (
	UpdateJobOptions_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	UpdateJobOptions_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	UpdateJobOptions_CommandName_CreateAction         = "create_action"
	UpdateJobOptions_CommandName_CreateCart           = "create_cart"
	UpdateJobOptions_CommandName_CreateEnvironment    = "create_environment"
	UpdateJobOptions_CommandName_CreateWorkspace      = "create_workspace"
	UpdateJobOptions_CommandName_DeleteAction         = "delete_action"
	UpdateJobOptions_CommandName_DeleteEnvironment    = "delete_environment"
	UpdateJobOptions_CommandName_DeleteWorkspace      = "delete_workspace"
	UpdateJobOptions_CommandName_EnvironmentInit      = "environment_init"
	UpdateJobOptions_CommandName_EnvironmentInstall   = "environment_install"
	UpdateJobOptions_CommandName_EnvironmentUninstall = "environment_uninstall"
	UpdateJobOptions_CommandName_PatchAction          = "patch_action"
	UpdateJobOptions_CommandName_PatchWorkspace       = "patch_workspace"
	UpdateJobOptions_CommandName_PutAction            = "put_action"
	UpdateJobOptions_CommandName_PutEnvironment       = "put_environment"
	UpdateJobOptions_CommandName_PutWorkspace         = "put_workspace"
	UpdateJobOptions_CommandName_RepositoryProcess    = "repository_process"
	UpdateJobOptions_CommandName_SystemKeyDelete      = "system_key_delete"
	UpdateJobOptions_CommandName_SystemKeyDisable     = "system_key_disable"
	UpdateJobOptions_CommandName_SystemKeyEnable      = "system_key_enable"
	UpdateJobOptions_CommandName_SystemKeyRestore     = "system_key_restore"
	UpdateJobOptions_CommandName_SystemKeyRotate      = "system_key_rotate"
	UpdateJobOptions_CommandName_WorkspaceApply       = "workspace_apply"
	UpdateJobOptions_CommandName_WorkspaceDestroy     = "workspace_destroy"
	UpdateJobOptions_CommandName_WorkspacePlan        = "workspace_plan"
	UpdateJobOptions_CommandName_WorkspaceRefresh     = "workspace_refresh"
)

// Constants associated with the UpdateJobOptions.Location property.
// List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the
// right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources,
// provisioned using Schematics.
const (
	UpdateJobOptions_Location_EuDe    = "eu-de"
	UpdateJobOptions_Location_EuGb    = "eu-gb"
	UpdateJobOptions_Location_UsEast  = "us-east"
	UpdateJobOptions_Location_UsSouth = "us-south"
)

// NewUpdateJobOptions : Instantiate UpdateJobOptions
func (*SchematicsV1) NewUpdateJobOptions(jobID string, refreshToken string) *UpdateJobOptions {
	return &UpdateJobOptions{
		JobID:        core.StringPtr(jobID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetJobID : Allow user to set JobID
func (_options *UpdateJobOptions) SetJobID(jobID string) *UpdateJobOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetRefreshToken : Allow user to set RefreshToken
func (_options *UpdateJobOptions) SetRefreshToken(refreshToken string) *UpdateJobOptions {
	_options.RefreshToken = core.StringPtr(refreshToken)
	return _options
}

// SetCommandObject : Allow user to set CommandObject
func (_options *UpdateJobOptions) SetCommandObject(commandObject string) *UpdateJobOptions {
	_options.CommandObject = core.StringPtr(commandObject)
	return _options
}

// SetCommandObjectID : Allow user to set CommandObjectID
func (_options *UpdateJobOptions) SetCommandObjectID(commandObjectID string) *UpdateJobOptions {
	_options.CommandObjectID = core.StringPtr(commandObjectID)
	return _options
}

// SetCommandName : Allow user to set CommandName
func (_options *UpdateJobOptions) SetCommandName(commandName string) *UpdateJobOptions {
	_options.CommandName = core.StringPtr(commandName)
	return _options
}

// SetCommandParameter : Allow user to set CommandParameter
func (_options *UpdateJobOptions) SetCommandParameter(commandParameter string) *UpdateJobOptions {
	_options.CommandParameter = core.StringPtr(commandParameter)
	return _options
}

// SetCommandOptions : Allow user to set CommandOptions
func (_options *UpdateJobOptions) SetCommandOptions(commandOptions []string) *UpdateJobOptions {
	_options.CommandOptions = commandOptions
	return _options
}

// SetInputs : Allow user to set Inputs
func (_options *UpdateJobOptions) SetInputs(inputs []VariableData) *UpdateJobOptions {
	_options.Inputs = inputs
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *UpdateJobOptions) SetSettings(settings []VariableData) *UpdateJobOptions {
	_options.Settings = settings
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateJobOptions) SetTags(tags []string) *UpdateJobOptions {
	_options.Tags = tags
	return _options
}

// SetLocation : Allow user to set Location
func (_options *UpdateJobOptions) SetLocation(location string) *UpdateJobOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *UpdateJobOptions) SetStatus(status *JobStatus) *UpdateJobOptions {
	_options.Status = status
	return _options
}

// SetData : Allow user to set Data
func (_options *UpdateJobOptions) SetData(data *JobData) *UpdateJobOptions {
	_options.Data = data
	return _options
}

// SetBastion : Allow user to set Bastion
func (_options *UpdateJobOptions) SetBastion(bastion *BastionResourceDefinition) *UpdateJobOptions {
	_options.Bastion = bastion
	return _options
}

// SetLogSummary : Allow user to set LogSummary
func (_options *UpdateJobOptions) SetLogSummary(logSummary *JobLogSummary) *UpdateJobOptions {
	_options.LogSummary = logSummary
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateJobOptions) SetHeaders(param map[string]string) *UpdateJobOptions {
	options.Headers = param
	return options
}

// UpdateKmsSettingsOptions : The UpdateKmsSettings options.
type UpdateKmsSettingsOptions struct {
	// Location.
	Location *string `json:"location,omitempty"`

	// Encryption scheme.
	EncryptionScheme *string `json:"encryption_scheme,omitempty"`

	// Resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Primary CRK details.
	PrimaryCrk *KMSSettingsPrimaryCrk `json:"primary_crk,omitempty"`

	// Secondary CRK details.
	SecondaryCrk *KMSSettingsSecondaryCrk `json:"secondary_crk,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateKmsSettingsOptions : Instantiate UpdateKmsSettingsOptions
func (*SchematicsV1) NewUpdateKmsSettingsOptions() *UpdateKmsSettingsOptions {
	return &UpdateKmsSettingsOptions{}
}

// SetLocation : Allow user to set Location
func (_options *UpdateKmsSettingsOptions) SetLocation(location string) *UpdateKmsSettingsOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetEncryptionScheme : Allow user to set EncryptionScheme
func (_options *UpdateKmsSettingsOptions) SetEncryptionScheme(encryptionScheme string) *UpdateKmsSettingsOptions {
	_options.EncryptionScheme = core.StringPtr(encryptionScheme)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *UpdateKmsSettingsOptions) SetResourceGroup(resourceGroup string) *UpdateKmsSettingsOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetPrimaryCrk : Allow user to set PrimaryCrk
func (_options *UpdateKmsSettingsOptions) SetPrimaryCrk(primaryCrk *KMSSettingsPrimaryCrk) *UpdateKmsSettingsOptions {
	_options.PrimaryCrk = primaryCrk
	return _options
}

// SetSecondaryCrk : Allow user to set SecondaryCrk
func (_options *UpdateKmsSettingsOptions) SetSecondaryCrk(secondaryCrk *KMSSettingsSecondaryCrk) *UpdateKmsSettingsOptions {
	_options.SecondaryCrk = secondaryCrk
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateKmsSettingsOptions) SetHeaders(param map[string]string) *UpdateKmsSettingsOptions {
	options.Headers = param
	return options
}

// UpdateWorkspaceOptions : The UpdateWorkspace options.
type UpdateWorkspaceOptions struct {
	// The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.
	WID *string `json:"-" validate:"required,ne="`

	// Information about the software template that you chose from the IBM Cloud catalog. This information is returned for
	// IBM Cloud catalog offerings only.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// The description of the workspace.
	Description *string `json:"description,omitempty"`

	// The name of the workspace.
	Name *string `json:"name,omitempty"`

	// Information about the Target used by the templates originating from the  IBM Cloud catalog offerings. This
	// information is not relevant for workspace created using your own Terraform template.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// A list of tags that you want to associate with your workspace.
	Tags []string `json:"tags,omitempty"`

	// Input data for the Template.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// Input to update the template repository data.
	TemplateRepo *TemplateRepoUpdateRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// Input to update the workspace status.
	WorkspaceStatus *WorkspaceStatusUpdateRequest `json:"workspace_status,omitempty"`

	// Information about the last job that ran against the workspace. -.
	WorkspaceStatusMsg *WorkspaceStatusMessage `json:"workspace_status_msg,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateWorkspaceOptions : Instantiate UpdateWorkspaceOptions
func (*SchematicsV1) NewUpdateWorkspaceOptions(wID string) *UpdateWorkspaceOptions {
	return &UpdateWorkspaceOptions{
		WID: core.StringPtr(wID),
	}
}

// SetWID : Allow user to set WID
func (_options *UpdateWorkspaceOptions) SetWID(wID string) *UpdateWorkspaceOptions {
	_options.WID = core.StringPtr(wID)
	return _options
}

// SetCatalogRef : Allow user to set CatalogRef
func (_options *UpdateWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *UpdateWorkspaceOptions {
	_options.CatalogRef = catalogRef
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateWorkspaceOptions) SetDescription(description string) *UpdateWorkspaceOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateWorkspaceOptions) SetName(name string) *UpdateWorkspaceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSharedData : Allow user to set SharedData
func (_options *UpdateWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *UpdateWorkspaceOptions {
	_options.SharedData = sharedData
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateWorkspaceOptions) SetTags(tags []string) *UpdateWorkspaceOptions {
	_options.Tags = tags
	return _options
}

// SetTemplateData : Allow user to set TemplateData
func (_options *UpdateWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *UpdateWorkspaceOptions {
	_options.TemplateData = templateData
	return _options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (_options *UpdateWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoUpdateRequest) *UpdateWorkspaceOptions {
	_options.TemplateRepo = templateRepo
	return _options
}

// SetType : Allow user to set Type
func (_options *UpdateWorkspaceOptions) SetType(typeVar []string) *UpdateWorkspaceOptions {
	_options.Type = typeVar
	return _options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (_options *UpdateWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusUpdateRequest) *UpdateWorkspaceOptions {
	_options.WorkspaceStatus = workspaceStatus
	return _options
}

// SetWorkspaceStatusMsg : Allow user to set WorkspaceStatusMsg
func (_options *UpdateWorkspaceOptions) SetWorkspaceStatusMsg(workspaceStatusMsg *WorkspaceStatusMessage) *UpdateWorkspaceOptions {
	_options.WorkspaceStatusMsg = workspaceStatusMsg
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateWorkspaceOptions) SetHeaders(param map[string]string) *UpdateWorkspaceOptions {
	options.Headers = param
	return options
}

// UploadTemplateTarActionOptions : The UploadTemplateTarAction options.
type UploadTemplateTarActionOptions struct {
	// Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.
	ActionID *string `json:"-" validate:"required,ne="`

	// Template tar file.
	File io.ReadCloser `json:"-"`

	// The content type of file.
	FileContentType *string `json:"-"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUploadTemplateTarActionOptions : Instantiate UploadTemplateTarActionOptions
func (*SchematicsV1) NewUploadTemplateTarActionOptions(actionID string) *UploadTemplateTarActionOptions {
	return &UploadTemplateTarActionOptions{
		ActionID: core.StringPtr(actionID),
	}
}

// SetActionID : Allow user to set ActionID
func (_options *UploadTemplateTarActionOptions) SetActionID(actionID string) *UploadTemplateTarActionOptions {
	_options.ActionID = core.StringPtr(actionID)
	return _options
}

// SetFile : Allow user to set File
func (_options *UploadTemplateTarActionOptions) SetFile(file io.ReadCloser) *UploadTemplateTarActionOptions {
	_options.File = file
	return _options
}

// SetFileContentType : Allow user to set FileContentType
func (_options *UploadTemplateTarActionOptions) SetFileContentType(fileContentType string) *UploadTemplateTarActionOptions {
	_options.FileContentType = core.StringPtr(fileContentType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UploadTemplateTarActionOptions) SetHeaders(param map[string]string) *UploadTemplateTarActionOptions {
	options.Headers = param
	return options
}

// UserState : User defined status of the Schematics object.
type UserState struct {
	// User-defined states
	//   * `draft` Object can be modified; can be used by Jobs run by the author, during execution
	//   * `live` Object can be modified; can be used by Jobs during execution
	//   * `locked` Object cannot be modified; can be used by Jobs during execution
	//   * `disable` Object can be modified. cannot be used by Jobs during execution.
	State *string `json:"state,omitempty"`

	// Name of the User who set the state of the Object.
	SetBy *string `json:"set_by,omitempty"`

	// When the User who set the state of the Object.
	SetAt *strfmt.DateTime `json:"set_at,omitempty"`
}

// Constants associated with the UserState.State property.
// User-defined states
//   * `draft` Object can be modified; can be used by Jobs run by the author, during execution
//   * `live` Object can be modified; can be used by Jobs during execution
//   * `locked` Object cannot be modified; can be used by Jobs during execution
//   * `disable` Object can be modified. cannot be used by Jobs during execution.
const (
	UserState_State_Disable = "disable"
	UserState_State_Draft   = "draft"
	UserState_State_Live    = "live"
	UserState_State_Locked  = "locked"
)

// UnmarshalUserState unmarshals an instance of UserState from the specified map of raw messages.
func UnmarshalUserState(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserState)
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "set_by", &obj.SetBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "set_at", &obj.SetAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserValues : UserValues -.
type UserValues struct {
	// A list of environment variables that you want to apply during the execution of a bash script or Terraform job. This
	// field must be provided as a list of key-value pairs, for example, **TF_LOG=debug**. Each entry will be a map with
	// one entry where `key is the environment variable name and value is value`. You can define environment variables for
	// IBM Cloud catalog offerings that are provisioned by using a bash script. See [example to use special environment
	// variable](https://cloud.ibm.com/docs/schematics?topic=schematics-set-parallelism#parallelism-example)  that are
	// supported by Schematics.
	EnvValues []interface{} `json:"env_values,omitempty"`

	// User values.
	Values *string `json:"values,omitempty"`

	// Information about the input variables that your template uses.
	Variablestore []WorkspaceVariableResponse `json:"variablestore,omitempty"`
}

// UnmarshalUserValues unmarshals an instance of UserValues from the specified map of raw messages.
func UnmarshalUserValues(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserValues)
	err = core.UnmarshalPrimitive(m, "env_values", &obj.EnvValues)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "variablestore", &obj.Variablestore, UnmarshalWorkspaceVariableResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VariableData : User editable variable data & system generated reference to value.
type VariableData struct {
	// Name of the variable.
	Name *string `json:"name,omitempty"`

	// Value for the variable or reference to the value.
	Value *string `json:"value,omitempty"`

	// User editable metadata for the variables.
	Metadata *VariableMetadata `json:"metadata,omitempty"`

	// Reference link to the variable value By default the expression will point to self.value.
	Link *string `json:"link,omitempty"`
}

// UnmarshalVariableData unmarshals an instance of VariableData from the specified map of raw messages.
func UnmarshalVariableData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VariableData)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalVariableMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "link", &obj.Link)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VariableMetadata : User editable metadata for the variables.
type VariableMetadata struct {
	// Type of the variable.
	Type *string `json:"type,omitempty"`

	// List of aliases for the variable name.
	Aliases []string `json:"aliases,omitempty"`

	// Description of the meta data.
	Description *string `json:"description,omitempty"`

	// Default value for the variable, if the override value is not specified.
	DefaultValue *string `json:"default_value,omitempty"`

	// Is the variable secure or sensitive ?.
	Secure *bool `json:"secure,omitempty"`

	// Is the variable readonly ?.
	Immutable *bool `json:"immutable,omitempty"`

	// If true, the variable will not be displayed on UI or CLI.
	Hidden *bool `json:"hidden,omitempty"`

	// List of possible values for this variable.  If type is integer or date, then the array of string will be  converted
	// to array of integers or date during runtime.
	Options []string `json:"options,omitempty"`

	// Minimum value of the variable. Applicable for integer type.
	MinValue *int64 `json:"min_value,omitempty"`

	// Maximum value of the variable. Applicable for integer type.
	MaxValue *int64 `json:"max_value,omitempty"`

	// Minimum length of the variable value. Applicable for string type.
	MinLength *int64 `json:"min_length,omitempty"`

	// Maximum length of the variable value. Applicable for string type.
	MaxLength *int64 `json:"max_length,omitempty"`

	// Regex for the variable value.
	Matches *string `json:"matches,omitempty"`

	// Relative position of this variable in a list.
	Position *int64 `json:"position,omitempty"`

	// Display name of the group this variable belongs to.
	GroupBy *string `json:"group_by,omitempty"`

	// Source of this meta-data.
	Source *string `json:"source,omitempty"`
}

// Constants associated with the VariableMetadata.Type property.
// Type of the variable.
const (
	VariableMetadata_Type_Array   = "array"
	VariableMetadata_Type_Boolean = "boolean"
	VariableMetadata_Type_Complex = "complex"
	VariableMetadata_Type_Date    = "date"
	VariableMetadata_Type_Integer = "integer"
	VariableMetadata_Type_List    = "list"
	VariableMetadata_Type_Map     = "map"
	VariableMetadata_Type_String  = "string"
)

// UnmarshalVariableMetadata unmarshals an instance of VariableMetadata from the specified map of raw messages.
func UnmarshalVariableMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VariableMetadata)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aliases", &obj.Aliases)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_value", &obj.DefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "immutable", &obj.Immutable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "options", &obj.Options)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "min_value", &obj.MinValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_value", &obj.MaxValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "min_length", &obj.MinLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_length", &obj.MaxLength)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "matches", &obj.Matches)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "position", &obj.Position)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_by", &obj.GroupBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VersionResponse : Successful response when you retrieve detailed information about the IBM Cloud Schematics API.
type VersionResponse struct {
	// The date when the API version was built.
	Builddate *string `json:"builddate,omitempty"`

	// The build number that the API is based on.
	Buildno *string `json:"buildno,omitempty"`

	// The SHA value for the Git commit that represents the latest version of the API.
	Commitsha *string `json:"commitsha,omitempty"`

	// The Terraform Helm provider version that is used when you install Helm charts with Schematics.
	HelmProviderVersion *string `json:"helm_provider_version,omitempty"`

	// The Helm version that is used when you install Helm charts with Schematics.
	HelmVersion *string `json:"helm_version,omitempty"`

	// Supported template types.
	SupportedTemplateTypes interface{} `json:"supported_template_types,omitempty"`

	// The version of the IBM Cloud Terraform provider plug-in that is used when you apply Terraform templates with
	// Schematics.
	TerraformProviderVersion *string `json:"terraform_provider_version,omitempty"`

	// The Terraform version that is used when you apply Terraform templates with Schematics.
	TerraformVersion *string `json:"terraform_version,omitempty"`
}

// UnmarshalVersionResponse unmarshals an instance of VersionResponse from the specified map of raw messages.
func UnmarshalVersionResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VersionResponse)
	err = core.UnmarshalPrimitive(m, "builddate", &obj.Builddate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "buildno", &obj.Buildno)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "commitsha", &obj.Commitsha)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "helm_provider_version", &obj.HelmProviderVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "helm_version", &obj.HelmVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "supported_template_types", &obj.SupportedTemplateTypes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "terraform_provider_version", &obj.TerraformProviderVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "terraform_version", &obj.TerraformVersion)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivities : List of workspace jobs.
type WorkspaceActivities struct {
	// List of workspace jobs.
	Actions []WorkspaceActivity `json:"actions,omitempty"`

	// The ID of the workspace.
	WorkspaceID *string `json:"workspace_id,omitempty"`

	// The name of the workspace.
	WorkspaceName *string `json:"workspace_name,omitempty"`
}

// UnmarshalWorkspaceActivities unmarshals an instance of WorkspaceActivities from the specified map of raw messages.
func UnmarshalWorkspaceActivities(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivities)
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalWorkspaceActivity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workspace_id", &obj.WorkspaceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "workspace_name", &obj.WorkspaceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivity : Information about the workspace jobs.
type WorkspaceActivity struct {
	// The ID of the activity or job.  You can use the ID to retrieve the logs for that job by using the `GET
	// /v1/workspaces/{id}/actions/{action_id}/logs` API.
	ActionID *string `json:"action_id,omitempty"`

	// Information about the success or failure of your job,  including a success or error code and the timestamp when the
	// job succeeded or failed.
	Message []string `json:"message,omitempty"`

	// The type of actovoty or job that ran against your workspace.
	//
	//  * **APPLY**: The apply job was created when you used the `PUT /v1/workspaces/{id}/apply` API to apply a Terraform
	// template in IBM Cloud.
	//  * **DESTROY**: The destroy job was created when you used the `DELETE /v1/workspaces/{id}/destroy` API to remove all
	// resources that are associated with your workspace.
	//  * **PLAN**: The plan job was created when you used the `POST /v1/workspaces/{id}/plan` API to create a Terraform
	// execution plan.
	Name *string `json:"name,omitempty"`

	// The timestamp when the job was initiated.
	PerformedAt *strfmt.DateTime `json:"performed_at,omitempty"`

	// The user ID who initiated the job.
	PerformedBy *string `json:"performed_by,omitempty"`

	// The status of your activity or job. To retrieve the URL to your job logs, use the GET
	// /v1/workspaces/{id}/actions/{action_id}/logs API.
	//
	// * **COMPLETED**: The job completed successfully.
	// * **CREATED**: The job was created, but the provisioning, modification, or removal of IBM Cloud resources has not
	// started yet.
	// * **FAILED**: An error occurred during the plan, apply, or destroy job. Use the job ID to retrieve the URL to the
	// log files for your job.
	// * **IN PROGRESS**: The job is in progress. You can use the log_url to access the logs.
	Status *string `json:"status,omitempty"`

	// List of template activities.
	Templates []WorkspaceActivityTemplate `json:"templates,omitempty"`
}

// UnmarshalWorkspaceActivity unmarshals an instance of WorkspaceActivity from the specified map of raw messages.
func UnmarshalWorkspaceActivity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivity)
	err = core.UnmarshalPrimitive(m, "action_id", &obj.ActionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "performed_at", &obj.PerformedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "performed_by", &obj.PerformedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "templates", &obj.Templates, UnmarshalWorkspaceActivityTemplate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityApplyResult : Response after successfully initiating a request to `apply` the Terraform template in IBM Cloud.
type WorkspaceActivityApplyResult struct {
	// The ID of the activity or job that was created when you initiated a request to `apply` a Terraform template.  You
	// can use the ID to retrieve log file by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
	Activityid *string `json:"activityid,omitempty"`
}

// UnmarshalWorkspaceActivityApplyResult unmarshals an instance of WorkspaceActivityApplyResult from the specified map of raw messages.
func UnmarshalWorkspaceActivityApplyResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityApplyResult)
	err = core.UnmarshalPrimitive(m, "activityid", &obj.Activityid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityCommandResult : Response after successfully initiating a request to run a workspace command on the stack of resources provisioned
// using Terraform.
type WorkspaceActivityCommandResult struct {
	// The ID of the job that was created when you initiated a request to `apply` a Terraform template.  You can use the ID
	// to retrieve log file by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
	Activityid *string `json:"activityid,omitempty"`
}

// UnmarshalWorkspaceActivityCommandResult unmarshals an instance of WorkspaceActivityCommandResult from the specified map of raw messages.
func UnmarshalWorkspaceActivityCommandResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityCommandResult)
	err = core.UnmarshalPrimitive(m, "activityid", &obj.Activityid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityDestroyResult : Response after successfully initiating a request to `destroy` the stack of resources provisioned using Terraform.
type WorkspaceActivityDestroyResult struct {
	// The ID of the activity or job that was created when you initiated a request to `destroy` a Terraform template.  You
	// can use the ID to retrieve log file by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
	Activityid *string `json:"activityid,omitempty"`
}

// UnmarshalWorkspaceActivityDestroyResult unmarshals an instance of WorkspaceActivityDestroyResult from the specified map of raw messages.
func UnmarshalWorkspaceActivityDestroyResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityDestroyResult)
	err = core.UnmarshalPrimitive(m, "activityid", &obj.Activityid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityLogs : Workspace job logs for all the templates in the workspace.
type WorkspaceActivityLogs struct {
	// The ID of the activity or job that ran against your workspace.
	ActionID *string `json:"action_id,omitempty"`

	// The type of actovoty or job that ran against your workspace.
	//
	//  * **APPLY**: The apply job was created when you used the `PUT /v1/workspaces/{id}/apply` API to apply a Terraform
	// template in IBM Cloud.
	//  * **DESTROY**: The destroy job was created when you used the `DELETE /v1/workspaces/{id}/destroy` API to remove all
	// resources that are associated with your workspace.
	//  * **PLAN**: The plan job was created when you used the `POST /v1/workspaces/{id}/plan` API to create a Terraform
	// execution plan.
	Name *string `json:"name,omitempty"`

	// List of templates in the workspace.
	Templates []WorkspaceActivityTemplateLogs `json:"templates,omitempty"`
}

// UnmarshalWorkspaceActivityLogs unmarshals an instance of WorkspaceActivityLogs from the specified map of raw messages.
func UnmarshalWorkspaceActivityLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityLogs)
	err = core.UnmarshalPrimitive(m, "action_id", &obj.ActionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "templates", &obj.Templates, UnmarshalWorkspaceActivityTemplateLogs)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityOptionsTemplate : Workspace job options template.
type WorkspaceActivityOptionsTemplate struct {
	// A list of Terraform resources to target.
	Target []string `json:"target,omitempty"`

	// Terraform variables for the workspace job options.
	TfVars []string `json:"tf_vars,omitempty"`
}

// UnmarshalWorkspaceActivityOptionsTemplate unmarshals an instance of WorkspaceActivityOptionsTemplate from the specified map of raw messages.
func UnmarshalWorkspaceActivityOptionsTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityOptionsTemplate)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tf_vars", &obj.TfVars)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityPlanResult : Response after successfully initiating a request to `plan` the Terraform template in IBM Cloud.
type WorkspaceActivityPlanResult struct {
	// The ID of the activity or job that was created when you initiated a request to `plan` a Terraform template.  You can
	// use the ID to retrieve log file by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
	Activityid *string `json:"activityid,omitempty"`
}

// UnmarshalWorkspaceActivityPlanResult unmarshals an instance of WorkspaceActivityPlanResult from the specified map of raw messages.
func UnmarshalWorkspaceActivityPlanResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityPlanResult)
	err = core.UnmarshalPrimitive(m, "activityid", &obj.Activityid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityRefreshResult : Response after successfully initiating a request to `refresh` the Terraform template in IBM Cloud.
type WorkspaceActivityRefreshResult struct {
	// The ID of the activity or job that was created for your workspace `refresh` activity or job.  You can use the ID to
	// retrieve the log file by using the `GET /v1/workspaces/{id}/actions/{action_id}/logs` API.
	Activityid *string `json:"activityid,omitempty"`
}

// UnmarshalWorkspaceActivityRefreshResult unmarshals an instance of WorkspaceActivityRefreshResult from the specified map of raw messages.
func UnmarshalWorkspaceActivityRefreshResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityRefreshResult)
	err = core.UnmarshalPrimitive(m, "activityid", &obj.Activityid)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityTemplate : Information about the template in the workspace.
type WorkspaceActivityTemplate struct {
	// End time for the job.
	EndTime *strfmt.DateTime `json:"end_time,omitempty"`

	// Summary information extracted from the job logs.
	LogSummary *LogSummary `json:"log_summary,omitempty"`

	// Log URL.
	LogURL *string `json:"log_url,omitempty"`

	// Message.
	Message *string `json:"message,omitempty"`

	// Job start time.
	StartTime *strfmt.DateTime `json:"start_time,omitempty"`

	// The status of your activity or job. To retrieve the URL to your job logs, use the GET
	// /v1/workspaces/{id}/actions/{action_id}/logs API.
	//
	// * **COMPLETED**: The job completed successfully.
	// * **CREATED**: The job was created, but the provisioning, modification, or removal of IBM Cloud resources has not
	// started yet.
	// * **FAILED**: An error occurred during the plan, apply, or destroy job. Use the job ID to retrieve the URL to the
	// log files for your job.
	// * **IN PROGRESS**: The job is in progress. You can use the log_url to access the logs.
	Status *string `json:"status,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	TemplateID *string `json:"template_id,omitempty"`

	// The type of template.
	TemplateType *string `json:"template_type,omitempty"`
}

// UnmarshalWorkspaceActivityTemplate unmarshals an instance of WorkspaceActivityTemplate from the specified map of raw messages.
func UnmarshalWorkspaceActivityTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityTemplate)
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "log_summary", &obj.LogSummary, UnmarshalLogSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "log_url", &obj.LogURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_type", &obj.TemplateType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceActivityTemplateLogs : Information about the log URL for a job that ran for a template against your workspace.
type WorkspaceActivityTemplateLogs struct {
	// The URL to access the logs that were created during the plan, apply, or destroy job.
	LogURL *string `json:"log_url,omitempty"`

	// The ID that was assigned to your Terraform template or IBM Cloud catalog software template.
	TemplateID *string `json:"template_id,omitempty"`

	// The type of template.
	TemplateType *string `json:"template_type,omitempty"`
}

// UnmarshalWorkspaceActivityTemplateLogs unmarshals an instance of WorkspaceActivityTemplateLogs from the specified map of raw messages.
func UnmarshalWorkspaceActivityTemplateLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceActivityTemplateLogs)
	err = core.UnmarshalPrimitive(m, "log_url", &obj.LogURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_type", &obj.TemplateType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceBulkDeleteResponse : Response after successfully initiating the bulk job to delete multiple workspaces.
type WorkspaceBulkDeleteResponse struct {
	// Workspace deletion job name.
	Job *string `json:"job,omitempty"`

	// Workspace deletion job id.
	JobID *string `json:"job_id,omitempty"`
}

// UnmarshalWorkspaceBulkDeleteResponse unmarshals an instance of WorkspaceBulkDeleteResponse from the specified map of raw messages.
func UnmarshalWorkspaceBulkDeleteResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceBulkDeleteResponse)
	err = core.UnmarshalPrimitive(m, "job", &obj.Job)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceJobResponse : Response from the workspace bulk job status.
type WorkspaceJobResponse struct {
	// Status of the workspace bulk job.
	JobStatus *WorkspaceJobStatusType `json:"job_status,omitempty"`
}

// UnmarshalWorkspaceJobResponse unmarshals an instance of WorkspaceJobResponse from the specified map of raw messages.
func UnmarshalWorkspaceJobResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceJobResponse)
	err = core.UnmarshalModel(m, "job_status", &obj.JobStatus, UnmarshalWorkspaceJobStatusType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceJobStatusType : Status of the workspace bulk job.
type WorkspaceJobStatusType struct {
	// List of failed workspace jobs.
	Failed []string `json:"failed,omitempty"`

	// List of in_progress workspace jobs.
	InProgress []string `json:"in_progress,omitempty"`

	// List of successful workspace jobs.
	Success []string `json:"success,omitempty"`

	// Workspace job status updated at.
	LastUpdatedOn *strfmt.DateTime `json:"last_updated_on,omitempty"`
}

// UnmarshalWorkspaceJobStatusType unmarshals an instance of WorkspaceJobStatusType from the specified map of raw messages.
func UnmarshalWorkspaceJobStatusType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceJobStatusType)
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "in_progress", &obj.InProgress)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_updated_on", &obj.LastUpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceResponse : Workspace details.
type WorkspaceResponse struct {
	// List of applied shared dataset ID.
	AppliedShareddataIds []string `json:"applied_shareddata_ids,omitempty"`

	// Information about the software template that you chose from the IBM Cloud catalog. This information is returned for
	// IBM Cloud catalog offerings only.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// The timestamp when the workspace was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The user ID that created the workspace.
	CreatedBy *string `json:"created_by,omitempty"`

	// The workspace CRN.
	Crn *string `json:"crn,omitempty"`

	// The description of the workspace.
	Description *string `json:"description,omitempty"`

	// The unique identifier of the workspace.
	ID *string `json:"id,omitempty"`

	// The timestamp when the last health check was performed by Schematics.
	LastHealthCheckAt *strfmt.DateTime `json:"last_health_check_at,omitempty"`

	// The IBM Cloud location where your workspace was provisioned.
	Location *string `json:"location,omitempty"`

	// The name of the workspace.
	Name *string `json:"name,omitempty"`

	// The resource group the workspace was provisioned in.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Information about the provisioning engine, state file, and runtime logs.
	RuntimeData []TemplateRunTimeDataResponse `json:"runtime_data,omitempty"`

	// Information about the Target used by the templates originating from IBM Cloud catalog offerings. This information is
	// not relevant when you create a workspace from your own Terraform template.
	SharedData *SharedTargetDataResponse `json:"shared_data,omitempty"`

	// The status of the workspace.
	//
	//   **Active**: After you successfully ran your infrastructure code by applying your Terraform execution plan, the
	// state of your workspace changes to `Active`.
	//
	//   **Connecting**: Schematics tries to connect to the template in your source repo. If successfully connected, the
	// template is downloaded and metadata, such as input parameters, is extracted. After the template is downloaded, the
	// state of the workspace changes to `Scanning`.
	//
	//   **Draft**: The workspace is created without a reference to a GitHub or GitLab repository.
	//
	//   **Failed**: If errors occur during the execution of your infrastructure code in IBM Cloud Schematics, your
	// workspace status is set to `Failed`.
	//
	//   **Inactive**: The Terraform template was scanned successfully and the workspace creation is complete. You can now
	// start running Schematics plan and apply jobs to provision the IBM Cloud resources that you specified in your
	// template. If you have an `Active` workspace and decide to remove all your resources, your workspace is set to
	// `Inactive` after all your resources are removed.
	//
	//   **In progress**: When you instruct IBM Cloud Schematics to run your infrastructure code by applying your Terraform
	// execution plan, the status of our workspace changes to `In progress`.
	//
	//   **Scanning**: The download of the Terraform template is complete and vulnerability scanning started. If the scan
	// is successful, the workspace state changes to `Inactive`. If errors in your template are found, the state changes to
	// `Template Error`.
	//
	//   **Stopped**: The Schematics plan, apply, or destroy job was cancelled manually.
	//
	//   **Template Error**: The Schematics template contains errors and cannot be processed.
	Status *string `json:"status,omitempty"`

	// A list of tags that are associated with the workspace.
	Tags []string `json:"tags,omitempty"`

	// Information about the Terraform or IBM Cloud software template that you want to use.
	TemplateData []TemplateSourceDataResponse `json:"template_data,omitempty"`

	// Workspace template reference.
	TemplateRef *string `json:"template_ref,omitempty"`

	// Information about the Template repository used by the workspace.
	TemplateRepo *TemplateRepoResponse `json:"template_repo,omitempty"`

	// The Terraform version that was used to run your Terraform code.
	Type []string `json:"type,omitempty"`

	// The timestamp when the workspace was last updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The user ID that updated the workspace.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// Response that indicate the status of the workspace as either frozen or locked.
	WorkspaceStatus *WorkspaceStatusResponse `json:"workspace_status,omitempty"`

	// Information about the last job that ran against the workspace. -.
	WorkspaceStatusMsg *WorkspaceStatusMessage `json:"workspace_status_msg,omitempty"`
}

// UnmarshalWorkspaceResponse unmarshals an instance of WorkspaceResponse from the specified map of raw messages.
func UnmarshalWorkspaceResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceResponse)
	err = core.UnmarshalPrimitive(m, "applied_shareddata_ids", &obj.AppliedShareddataIds)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "catalog_ref", &obj.CatalogRef, UnmarshalCatalogRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_health_check_at", &obj.LastHealthCheckAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "runtime_data", &obj.RuntimeData, UnmarshalTemplateRunTimeDataResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "shared_data", &obj.SharedData, UnmarshalSharedTargetDataResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template_data", &obj.TemplateData, UnmarshalTemplateSourceDataResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_ref", &obj.TemplateRef)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template_repo", &obj.TemplateRepo, UnmarshalTemplateRepoResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workspace_status", &obj.WorkspaceStatus, UnmarshalWorkspaceStatusResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workspace_status_msg", &obj.WorkspaceStatusMsg, UnmarshalWorkspaceStatusMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceResponseList : List of workspaces.
type WorkspaceResponseList struct {
	// The number of workspaces in the IBM Cloud account that you have access to and that matched your search criteria.
	Count *int64 `json:"count,omitempty"`

	// The `limit` value that you set in your API request and that represents the maximum number of workspaces that you
	// wanted to list.
	Limit *int64 `json:"limit" validate:"required"`

	// The `offset` value that you set in your API request. The offset value represents the position number of the
	// workspace from which you wanted to start listing your workspaces.
	Offset *int64 `json:"offset" validate:"required"`

	// The list of workspaces that was included in your API response.
	Workspaces []WorkspaceResponse `json:"workspaces,omitempty"`
}

// UnmarshalWorkspaceResponseList unmarshals an instance of WorkspaceResponseList from the specified map of raw messages.
func UnmarshalWorkspaceResponseList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceResponseList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "workspaces", &obj.Workspaces, UnmarshalWorkspaceResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceStatusMessage : Information about the last job that ran against the workspace. -.
type WorkspaceStatusMessage struct {
	// The success or error code that was returned for the last plan, apply, or destroy job that ran against your
	// workspace.
	StatusCode *string `json:"status_code,omitempty"`

	// The success or error message that was returned for the last plan, apply, or destroy job that ran against your
	// workspace.
	StatusMsg *string `json:"status_msg,omitempty"`
}

// UnmarshalWorkspaceStatusMessage unmarshals an instance of WorkspaceStatusMessage from the specified map of raw messages.
func UnmarshalWorkspaceStatusMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceStatusMessage)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_msg", &obj.StatusMsg)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceStatusRequest : WorkspaceStatusRequest -.
type WorkspaceStatusRequest struct {
	// If set to true, the workspace is frozen and changes to the workspace are disabled.
	Frozen *bool `json:"frozen,omitempty"`

	// The timestamp when the workspace was frozen.
	FrozenAt *strfmt.DateTime `json:"frozen_at,omitempty"`

	// The user ID that froze the workspace.
	FrozenBy *string `json:"frozen_by,omitempty"`

	// If set to true, the workspace is locked and disabled for changes.
	Locked *bool `json:"locked,omitempty"`

	// The user ID that initiated a resource-related job, such as applying or destroying resources, that locked the
	// workspace.
	LockedBy *string `json:"locked_by,omitempty"`

	// The timestamp when the workspace was locked.
	LockedTime *strfmt.DateTime `json:"locked_time,omitempty"`
}

// UnmarshalWorkspaceStatusRequest unmarshals an instance of WorkspaceStatusRequest from the specified map of raw messages.
func UnmarshalWorkspaceStatusRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceStatusRequest)
	err = core.UnmarshalPrimitive(m, "frozen", &obj.Frozen)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_at", &obj.FrozenAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_by", &obj.FrozenBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_by", &obj.LockedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_time", &obj.LockedTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceStatusResponse : Response that indicate the status of the workspace as either frozen or locked.
type WorkspaceStatusResponse struct {
	// If set to true, the workspace is frozen and changes to the workspace are disabled.
	Frozen *bool `json:"frozen,omitempty"`

	// The timestamp when the workspace was frozen.
	FrozenAt *strfmt.DateTime `json:"frozen_at,omitempty"`

	// The user ID that froze the workspace.
	FrozenBy *string `json:"frozen_by,omitempty"`

	// If set to true, the workspace is locked and disabled for changes.
	Locked *bool `json:"locked,omitempty"`

	// The user ID that initiated a resource-related job, such as applying or destroying resources, that locked the
	// workspace.
	LockedBy *string `json:"locked_by,omitempty"`

	// The timestamp when the workspace was locked.
	LockedTime *strfmt.DateTime `json:"locked_time,omitempty"`
}

// UnmarshalWorkspaceStatusResponse unmarshals an instance of WorkspaceStatusResponse from the specified map of raw messages.
func UnmarshalWorkspaceStatusResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceStatusResponse)
	err = core.UnmarshalPrimitive(m, "frozen", &obj.Frozen)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_at", &obj.FrozenAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_by", &obj.FrozenBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_by", &obj.LockedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_time", &obj.LockedTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceStatusUpdateRequest : Input to update the workspace status.
type WorkspaceStatusUpdateRequest struct {
	// If set to true, the workspace is frozen and changes to the workspace are disabled.
	Frozen *bool `json:"frozen,omitempty"`

	// Frozen at.
	FrozenAt *strfmt.DateTime `json:"frozen_at,omitempty"`

	// Frozen by.
	FrozenBy *string `json:"frozen_by,omitempty"`

	// Locked status.
	Locked *bool `json:"locked,omitempty"`

	// Locked by.
	LockedBy *string `json:"locked_by,omitempty"`

	// Locked at.
	LockedTime *strfmt.DateTime `json:"locked_time,omitempty"`
}

// UnmarshalWorkspaceStatusUpdateRequest unmarshals an instance of WorkspaceStatusUpdateRequest from the specified map of raw messages.
func UnmarshalWorkspaceStatusUpdateRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceStatusUpdateRequest)
	err = core.UnmarshalPrimitive(m, "frozen", &obj.Frozen)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_at", &obj.FrozenAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frozen_by", &obj.FrozenBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_by", &obj.LockedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locked_time", &obj.LockedTime)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceTemplateValuesResponse : Response with the template details in your workspace.
type WorkspaceTemplateValuesResponse struct {
	// Information about the provisioning engine, state file, and runtime logs.
	RuntimeData []TemplateRunTimeDataResponse `json:"runtime_data,omitempty"`

	// Information about the Target used by the templates originating from the  IBM Cloud catalog offerings. This
	// information is not relevant for workspace created using your own Terraform template.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// Information about the input variables that are used in the template.
	TemplateData []TemplateSourceDataResponse `json:"template_data,omitempty"`
}

// UnmarshalWorkspaceTemplateValuesResponse unmarshals an instance of WorkspaceTemplateValuesResponse from the specified map of raw messages.
func UnmarshalWorkspaceTemplateValuesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceTemplateValuesResponse)
	err = core.UnmarshalModel(m, "runtime_data", &obj.RuntimeData, UnmarshalTemplateRunTimeDataResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "shared_data", &obj.SharedData, UnmarshalSharedTargetData)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template_data", &obj.TemplateData, UnmarshalTemplateSourceDataResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceVariableRequest : Input variables for your workspace.
type WorkspaceVariableRequest struct {
	// The description of your input variable.
	Description *string `json:"description,omitempty"`

	// The name of the variable.
	Name *string `json:"name,omitempty"`

	// If set to `true`, the value of your input variable is protected and not returned in your API response.
	Secure *bool `json:"secure,omitempty"`

	// `Terraform v0.11` supports `string`, `list`, `map` data type. For more information, about the syntax, see
	// [Configuring input variables](https://www.terraform.io/docs/configuration-0-11/variables.html).
	// <br> `Terraform v0.12` additionally, supports `bool`, `number` and complex data types such as `list(type)`,
	// `map(type)`,
	// `object({attribute name=type,..})`, `set(type)`, `tuple([type])`. For more information, about the syntax to use the
	// complex data type, see [Configuring
	// variables](https://www.terraform.io/docs/configuration/variables.html#type-constraints).
	Type *string `json:"type,omitempty"`

	// Variable uses default value; and is not over-ridden.
	UseDefault *bool `json:"use_default,omitempty"`

	// Enter the value as a string for the primitive types such as `bool`, `number`, `string`, and `HCL` format for the
	// complex variables, as you provide in a `.tfvars` file. **You need to enter escaped string of `HCL` format for the
	// complex variable value**. For more information, about how to declare variables in a terraform configuration file and
	// provide value to schematics, see [Providing values for the declared
	// variables](https://cloud.ibm.com/docs/schematics?topic=schematics-create-tf-config#declare-variable).
	Value *string `json:"value,omitempty"`
}

// UnmarshalWorkspaceVariableRequest unmarshals an instance of WorkspaceVariableRequest from the specified map of raw messages.
func UnmarshalWorkspaceVariableRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceVariableRequest)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_default", &obj.UseDefault)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceVariableResponse : The description of your input variable.
type WorkspaceVariableResponse struct {
	// The description of your input variable.
	Description *string `json:"description,omitempty"`

	// The name of the variable.
	Name *string `json:"name,omitempty"`

	// If set to `true`, the value of your input variable is protected and not returned in your API response.
	Secure *bool `json:"secure,omitempty"`

	// `Terraform v0.11` supports `string`, `list`, `map` data type. For more information, about the syntax, see
	// [Configuring input variables](https://www.terraform.io/docs/configuration-0-11/variables.html).
	// <br> `Terraform v0.12` additionally, supports `bool`, `number` and complex data types such as `list(type)`,
	// `map(type)`,
	// `object({attribute name=type,..})`, `set(type)`, `tuple([type])`. For more information, about the syntax to use the
	// complex data type, see [Configuring
	// variables](https://www.terraform.io/docs/configuration/variables.html#type-constraints).
	Type *string `json:"type,omitempty"`

	// Enter the value as a string for the primitive types such as `bool`, `number`, `string`, and `HCL` format for the
	// complex variables, as you provide in a `.tfvars` file. **You need to enter escaped string of `HCL` format for the
	// complex variable value**. For more information, about how to declare variables in a terraform configuration file and
	// provide value to schematics, see [Providing values for the declared
	// variables](https://cloud.ibm.com/docs/schematics?topic=schematics-create-tf-config#declare-variable).
	Value *string `json:"value,omitempty"`
}

// UnmarshalWorkspaceVariableResponse unmarshals an instance of WorkspaceVariableResponse from the specified map of raw messages.
func UnmarshalWorkspaceVariableResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceVariableResponse)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
