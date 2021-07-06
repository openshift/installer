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
 * IBM OpenAPI SDK Code Generator Version: 3.17.0-8d569e8f-20201030-142059
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

	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/schematics-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// SchematicsV1 : Schematics Service is to provide the capability to manage resources  of (cloud) provider
// infrastructure using file based configurations.  With the Schematics service the customer is able to specify the
// required set of resources and their configuration in ''config files'',  and then pass these config files to the
// service to fulfill it by  calling the necessary actions on the infrastructure.  This principle is also known as
// Infrastructure as Code.  For more information refer to
// https://cloud.ibm.com/docs/schematics?topic=schematics-getting-started'
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
// List supported schematics locations.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSchematicsLocations)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListResourceGroup : List of resource groups in the Account
// List of resource groups in the Account.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourceGroupResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSchematicsVersion : Get schematics version
// Get schematics version.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVersionResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListWorkspaces : List all workspace definitions
// List all workspace definitions.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponseList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateWorkspace : Create workspace definition
// Create workspace definition.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspace : Get workspace definition
// Get workspace definition.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ReplaceWorkspace : Replace the workspace definition
// Replace the workspace definition.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteWorkspace : Delete a workspace definition
// Delete a workspace definition.  Use destroy_resource='true' to destroy the related cloud resource.
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

// UpdateWorkspace : Update the workspace definition
// Update the workspace definition.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UploadTemplateTar : Upload template tar file for the workspace
// Upload template tar file for the workspace.
func (schematics *SchematicsV1) UploadTemplateTar(uploadTemplateTarOptions *UploadTemplateTarOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	return schematics.UploadTemplateTarWithContext(context.Background(), uploadTemplateTarOptions)
}

// UploadTemplateTarWithContext is an alternate form of the UploadTemplateTar method which supports a Context parameter
func (schematics *SchematicsV1) UploadTemplateTarWithContext(ctx context.Context, uploadTemplateTarOptions *UploadTemplateTarOptions) (result *TemplateRepoTarUploadResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(uploadTemplateTarOptions, "uploadTemplateTarOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(uploadTemplateTarOptions, "uploadTemplateTarOptions")
	if err != nil {
		return
	}
	if uploadTemplateTarOptions.File == nil {
		err = fmt.Errorf("at least one of  or file must be supplied")
		return
	}

	pathParamsMap := map[string]string{
		"w_id": *uploadTemplateTarOptions.WID,
		"t_id": *uploadTemplateTarOptions.TID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v1/workspaces/{w_id}/template_data/{t_id}/template_repo_upload`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadTemplateTarOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "UploadTemplateTar")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if uploadTemplateTarOptions.File != nil {
		builder.AddFormData("file", "filename",
			core.StringNilMapper(uploadTemplateTarOptions.FileContentType), uploadTemplateTarOptions.File)
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateRepoTarUploadResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceReadme : Get the workspace readme
// Get the workspace readme.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateReadme)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListWorkspaceActivities : List all workspace activities
// List all workspace activities.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivities)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceActivity : Get workspace activity details
// Get workspace activity details.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivity)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteWorkspaceActivity : Stop the workspace activity
// Stop the workspace activity.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityApplyResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// RunWorkspaceCommands : Run terraform Commands
// Run terraform Commands on workspaces.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityCommandResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ApplyWorkspaceCommand : Run schematics workspace 'apply' activity
// Run schematics workspace 'apply' activity.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityApplyResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DestroyWorkspaceCommand : Run workspace 'destroy' activity
// Run workspace 'destroy' activity,  to destroy all the resources associated with the workspace.  WARNING: This action
// cannot be reversed.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityDestroyResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// PlanWorkspaceCommand : Run workspace 'plan' activity,
// Run schematics workspace 'plan' activity,  to preview the change before running an 'apply' activity.
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityPlanResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// RefreshWorkspaceCommand : Run workspace 'refresh' activity
// Run workspace 'refresh' activity.
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityRefreshResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceInputs : Get the input values of the workspace
// Get the input values of the workspace.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateValues)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ReplaceWorkspaceInputs : Replace the input values for the workspace
// Replace the input values for the workspace.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserValues)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetAllWorkspaceInputs : Get all the input values of the workspace
// Get all the input values of the workspace.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceTemplateValuesResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceInputMetadata : Get the input metadata of the workspace
// Get the input metadata of the workspace.
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

// GetWorkspaceOutputs : Get all the output values of the workspace
// Get all the output values from your workspace; (ex. result of terraform output command).
func (schematics *SchematicsV1) GetWorkspaceOutputs(getWorkspaceOutputsOptions *GetWorkspaceOutputsOptions) (result []OutputValuesItem, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceOutputsWithContext(context.Background(), getWorkspaceOutputsOptions)
}

// GetWorkspaceOutputsWithContext is an alternate form of the GetWorkspaceOutputs method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceOutputsWithContext(ctx context.Context, getWorkspaceOutputsOptions *GetWorkspaceOutputsOptions) (result []OutputValuesItem, response *core.DetailedResponse, err error) {
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOutputValuesItem)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceResources : Get all the resources created by the workspace
// Get all the resources created by the workspace.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateResources)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceState : Get the workspace state
// Get the workspace state.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStateStoreResponseList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceTemplateState : Get the template state
// Get the template state.
func (schematics *SchematicsV1) GetWorkspaceTemplateState(getWorkspaceTemplateStateOptions *GetWorkspaceTemplateStateOptions) (result *TemplateStateStore, response *core.DetailedResponse, err error) {
	return schematics.GetWorkspaceTemplateStateWithContext(context.Background(), getWorkspaceTemplateStateOptions)
}

// GetWorkspaceTemplateStateWithContext is an alternate form of the GetWorkspaceTemplateState method which supports a Context parameter
func (schematics *SchematicsV1) GetWorkspaceTemplateStateWithContext(ctx context.Context, getWorkspaceTemplateStateOptions *GetWorkspaceTemplateStateOptions) (result *TemplateStateStore, response *core.DetailedResponse, err error) {
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

	// The response is the terraform statefile and the structure
	// can change between versions. So unmarshalling using a fixed
	// schema is impossible. Hence the result is sent as json.RawMessage
	var rawResponse json.RawMessage
	response, err = schematics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	/*err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateStateStore)
	if err != nil {
		return
	}*/
	response.Result = rawResponse

	return
}

// GetWorkspaceActivityLogs : Get the workspace activity log urls
// View an activity log for Terraform actions that ran against your workspace.  You can view logs for plan, apply, and
// destroy actions.      operationId: get_activity_log_urls.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceActivityLogs)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWorkspaceLogUrls : Get all workspace log urls
// Get all workspace log urls.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogStoreResponseList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetTemplateLogs : Get all template logs
// Get all template logs.
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

// GetTemplateActivityLog : Get the template activity logs
// View an activity log for Terraform actions that ran for a template against your workspace.  You can view logs for
// plan, apply, and destroy actions.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceBulkDeleteResponse)
	if err != nil {
		return
	}
	response.Result = result

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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWorkspaceJobResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateAction : Create an Action definition
// Create a new Action definition.
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
	if createActionOptions.Bastion != nil {
		body["bastion"] = createActionOptions.Bastion
	}
	if createActionOptions.TargetsIni != nil {
		body["targets_ini"] = createActionOptions.TargetsIni
	}
	if createActionOptions.Credentials != nil {
		body["credentials"] = createActionOptions.Credentials
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
	if createActionOptions.TriggerRecordID != nil {
		body["trigger_record_id"] = createActionOptions.TriggerRecordID
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListActions : Get all the Action definitions
// Get all the Action definitions.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActionList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetAction : Get the Action definition
// Get the Action definition.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteAction : Delete the Action
// Delete the Action definition.
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

// UpdateAction : Update the Action definition
// Update the Action definition.
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
	if updateActionOptions.Bastion != nil {
		body["bastion"] = updateActionOptions.Bastion
	}
	if updateActionOptions.TargetsIni != nil {
		body["targets_ini"] = updateActionOptions.TargetsIni
	}
	if updateActionOptions.Credentials != nil {
		body["credentials"] = updateActionOptions.Credentials
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
	if updateActionOptions.TriggerRecordID != nil {
		body["trigger_record_id"] = updateActionOptions.TriggerRecordID
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAction)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateJob : Create a Job record and launch the Job
// Creare a Job record and launch the Job.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListJobs : Get all the Job records
// Get all the Job records.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ReplaceJob : Clone the Job-record, and relaunch the Job
// Clone the Job-record, and relaunch the Job.
func (schematics *SchematicsV1) ReplaceJob(replaceJobOptions *ReplaceJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	return schematics.ReplaceJobWithContext(context.Background(), replaceJobOptions)
}

// ReplaceJobWithContext is an alternate form of the ReplaceJob method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceJobWithContext(ctx context.Context, replaceJobOptions *ReplaceJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceJobOptions, "replaceJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceJobOptions, "replaceJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *replaceJobOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceJobOptions.RefreshToken != nil {
		builder.AddHeader("refresh_token", fmt.Sprint(*replaceJobOptions.RefreshToken))
	}

	body := make(map[string]interface{})
	if replaceJobOptions.CommandObject != nil {
		body["command_object"] = replaceJobOptions.CommandObject
	}
	if replaceJobOptions.CommandObjectID != nil {
		body["command_object_id"] = replaceJobOptions.CommandObjectID
	}
	if replaceJobOptions.CommandName != nil {
		body["command_name"] = replaceJobOptions.CommandName
	}
	if replaceJobOptions.CommandParameter != nil {
		body["command_parameter"] = replaceJobOptions.CommandParameter
	}
	if replaceJobOptions.CommandOptions != nil {
		body["command_options"] = replaceJobOptions.CommandOptions
	}
	if replaceJobOptions.Inputs != nil {
		body["inputs"] = replaceJobOptions.Inputs
	}
	if replaceJobOptions.Settings != nil {
		body["settings"] = replaceJobOptions.Settings
	}
	if replaceJobOptions.Tags != nil {
		body["tags"] = replaceJobOptions.Tags
	}
	if replaceJobOptions.Location != nil {
		body["location"] = replaceJobOptions.Location
	}
	if replaceJobOptions.Status != nil {
		body["status"] = replaceJobOptions.Status
	}
	if replaceJobOptions.Data != nil {
		body["data"] = replaceJobOptions.Data
	}
	if replaceJobOptions.Bastion != nil {
		body["bastion"] = replaceJobOptions.Bastion
	}
	if replaceJobOptions.LogSummary != nil {
		body["log_summary"] = replaceJobOptions.LogSummary
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteJob : Stop the running Job, and delete the Job-record
// Stop the running Job, and delete the Job-record.
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

// GetJob : Get the Job record
// Get the Job record.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListJobLogs : Get log-file from the Job record
// Get log-file from the Job record.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobLog)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListJobStates : Get state-data from the Job record
// Get state-data from the Job record.
func (schematics *SchematicsV1) ListJobStates(listJobStatesOptions *ListJobStatesOptions) (result *JobStateData, response *core.DetailedResponse, err error) {
	return schematics.ListJobStatesWithContext(context.Background(), listJobStatesOptions)
}

// ListJobStatesWithContext is an alternate form of the ListJobStates method which supports a Context parameter
func (schematics *SchematicsV1) ListJobStatesWithContext(ctx context.Context, listJobStatesOptions *ListJobStatesOptions) (result *JobStateData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listJobStatesOptions, "listJobStatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listJobStatesOptions, "listJobStatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"job_id": *listJobStatesOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/jobs/{job_id}/states`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listJobStatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListJobStates")
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobStateData)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListSharedDatasets : List all shared datasets
// List all shared datasets.
func (schematics *SchematicsV1) ListSharedDatasets(listSharedDatasetsOptions *ListSharedDatasetsOptions) (result *SharedDatasetResponseList, response *core.DetailedResponse, err error) {
	return schematics.ListSharedDatasetsWithContext(context.Background(), listSharedDatasetsOptions)
}

// ListSharedDatasetsWithContext is an alternate form of the ListSharedDatasets method which supports a Context parameter
func (schematics *SchematicsV1) ListSharedDatasetsWithContext(ctx context.Context, listSharedDatasetsOptions *ListSharedDatasetsOptions) (result *SharedDatasetResponseList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSharedDatasetsOptions, "listSharedDatasetsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/shared_datasets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSharedDatasetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ListSharedDatasets")
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSharedDatasetResponseList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateSharedDataset : Create a shared dataset definition
// Create a shared dataset definition.
func (schematics *SchematicsV1) CreateSharedDataset(createSharedDatasetOptions *CreateSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	return schematics.CreateSharedDatasetWithContext(context.Background(), createSharedDatasetOptions)
}

// CreateSharedDatasetWithContext is an alternate form of the CreateSharedDataset method which supports a Context parameter
func (schematics *SchematicsV1) CreateSharedDatasetWithContext(ctx context.Context, createSharedDatasetOptions *CreateSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSharedDatasetOptions, "createSharedDatasetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSharedDatasetOptions, "createSharedDatasetOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/shared_datasets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSharedDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "CreateSharedDataset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSharedDatasetOptions.AutoPropagateChange != nil {
		body["auto_propagate_change"] = createSharedDatasetOptions.AutoPropagateChange
	}
	if createSharedDatasetOptions.Description != nil {
		body["description"] = createSharedDatasetOptions.Description
	}
	if createSharedDatasetOptions.EffectedWorkspaceIds != nil {
		body["effected_workspace_ids"] = createSharedDatasetOptions.EffectedWorkspaceIds
	}
	if createSharedDatasetOptions.ResourceGroup != nil {
		body["resource_group"] = createSharedDatasetOptions.ResourceGroup
	}
	if createSharedDatasetOptions.SharedDatasetData != nil {
		body["shared_dataset_data"] = createSharedDatasetOptions.SharedDatasetData
	}
	if createSharedDatasetOptions.SharedDatasetName != nil {
		body["shared_dataset_name"] = createSharedDatasetOptions.SharedDatasetName
	}
	if createSharedDatasetOptions.SharedDatasetSourceName != nil {
		body["shared_dataset_source_name"] = createSharedDatasetOptions.SharedDatasetSourceName
	}
	if createSharedDatasetOptions.SharedDatasetType != nil {
		body["shared_dataset_type"] = createSharedDatasetOptions.SharedDatasetType
	}
	if createSharedDatasetOptions.Tags != nil {
		body["tags"] = createSharedDatasetOptions.Tags
	}
	if createSharedDatasetOptions.Version != nil {
		body["version"] = createSharedDatasetOptions.Version
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSharedDatasetResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSharedDataset : Get the shared dataset
// Get the shared dataset.
func (schematics *SchematicsV1) GetSharedDataset(getSharedDatasetOptions *GetSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	return schematics.GetSharedDatasetWithContext(context.Background(), getSharedDatasetOptions)
}

// GetSharedDatasetWithContext is an alternate form of the GetSharedDataset method which supports a Context parameter
func (schematics *SchematicsV1) GetSharedDatasetWithContext(ctx context.Context, getSharedDatasetOptions *GetSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSharedDatasetOptions, "getSharedDatasetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSharedDatasetOptions, "getSharedDatasetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"sd_id": *getSharedDatasetOptions.SdID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/shared_datasets/{sd_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSharedDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetSharedDataset")
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSharedDatasetResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ReplaceSharedDataset : Replace the shared dataset
// Replace the shared dataset.
func (schematics *SchematicsV1) ReplaceSharedDataset(replaceSharedDatasetOptions *ReplaceSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	return schematics.ReplaceSharedDatasetWithContext(context.Background(), replaceSharedDatasetOptions)
}

// ReplaceSharedDatasetWithContext is an alternate form of the ReplaceSharedDataset method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceSharedDatasetWithContext(ctx context.Context, replaceSharedDatasetOptions *ReplaceSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceSharedDatasetOptions, "replaceSharedDatasetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceSharedDatasetOptions, "replaceSharedDatasetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"sd_id": *replaceSharedDatasetOptions.SdID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/shared_datasets/{sd_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceSharedDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceSharedDataset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceSharedDatasetOptions.AutoPropagateChange != nil {
		body["auto_propagate_change"] = replaceSharedDatasetOptions.AutoPropagateChange
	}
	if replaceSharedDatasetOptions.Description != nil {
		body["description"] = replaceSharedDatasetOptions.Description
	}
	if replaceSharedDatasetOptions.EffectedWorkspaceIds != nil {
		body["effected_workspace_ids"] = replaceSharedDatasetOptions.EffectedWorkspaceIds
	}
	if replaceSharedDatasetOptions.ResourceGroup != nil {
		body["resource_group"] = replaceSharedDatasetOptions.ResourceGroup
	}
	if replaceSharedDatasetOptions.SharedDatasetData != nil {
		body["shared_dataset_data"] = replaceSharedDatasetOptions.SharedDatasetData
	}
	if replaceSharedDatasetOptions.SharedDatasetName != nil {
		body["shared_dataset_name"] = replaceSharedDatasetOptions.SharedDatasetName
	}
	if replaceSharedDatasetOptions.SharedDatasetSourceName != nil {
		body["shared_dataset_source_name"] = replaceSharedDatasetOptions.SharedDatasetSourceName
	}
	if replaceSharedDatasetOptions.SharedDatasetType != nil {
		body["shared_dataset_type"] = replaceSharedDatasetOptions.SharedDatasetType
	}
	if replaceSharedDatasetOptions.Tags != nil {
		body["tags"] = replaceSharedDatasetOptions.Tags
	}
	if replaceSharedDatasetOptions.Version != nil {
		body["version"] = replaceSharedDatasetOptions.Version
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSharedDatasetResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteSharedDataset : Delete the shared dataset
// Replace the shared dataset.
func (schematics *SchematicsV1) DeleteSharedDataset(deleteSharedDatasetOptions *DeleteSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	return schematics.DeleteSharedDatasetWithContext(context.Background(), deleteSharedDatasetOptions)
}

// DeleteSharedDatasetWithContext is an alternate form of the DeleteSharedDataset method which supports a Context parameter
func (schematics *SchematicsV1) DeleteSharedDatasetWithContext(ctx context.Context, deleteSharedDatasetOptions *DeleteSharedDatasetOptions) (result *SharedDatasetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSharedDatasetOptions, "deleteSharedDatasetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSharedDatasetOptions, "deleteSharedDatasetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"sd_id": *deleteSharedDatasetOptions.SdID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = schematics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(schematics.Service.Options.URL, `/v2/shared_datasets/{sd_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSharedDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "DeleteSharedDataset")
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSharedDatasetResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetKmsSettings : Get the KMS settings for customer account
// Get the KMS settings for customer account.
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSSettings)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ReplaceKmsSettings : Set the KMS settings for customer account
// Set the KMS settings for customer account.
func (schematics *SchematicsV1) ReplaceKmsSettings(replaceKmsSettingsOptions *ReplaceKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	return schematics.ReplaceKmsSettingsWithContext(context.Background(), replaceKmsSettingsOptions)
}

// ReplaceKmsSettingsWithContext is an alternate form of the ReplaceKmsSettings method which supports a Context parameter
func (schematics *SchematicsV1) ReplaceKmsSettingsWithContext(ctx context.Context, replaceKmsSettingsOptions *ReplaceKmsSettingsOptions) (result *KMSSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceKmsSettingsOptions, "replaceKmsSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceKmsSettingsOptions, "replaceKmsSettingsOptions")
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

	for headerName, headerValue := range replaceKmsSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "ReplaceKmsSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceKmsSettingsOptions.Location != nil {
		body["location"] = replaceKmsSettingsOptions.Location
	}
	if replaceKmsSettingsOptions.EncryptionScheme != nil {
		body["encryption_scheme"] = replaceKmsSettingsOptions.EncryptionScheme
	}
	if replaceKmsSettingsOptions.ResourceGroup != nil {
		body["resource_group"] = replaceKmsSettingsOptions.ResourceGroup
	}
	if replaceKmsSettingsOptions.PrimaryCrk != nil {
		body["primary_crk"] = replaceKmsSettingsOptions.PrimaryCrk
	}
	if replaceKmsSettingsOptions.SecondaryCrk != nil {
		body["secondary_crk"] = replaceKmsSettingsOptions.SecondaryCrk
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSSettings)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetDiscoveredKmsInstances : Discover the KMS instances in the account
// Discover the KMS instances in the account.
func (schematics *SchematicsV1) GetDiscoveredKmsInstances(getDiscoveredKmsInstancesOptions *GetDiscoveredKmsInstancesOptions) (result *KMSDiscovery, response *core.DetailedResponse, err error) {
	return schematics.GetDiscoveredKmsInstancesWithContext(context.Background(), getDiscoveredKmsInstancesOptions)
}

// GetDiscoveredKmsInstancesWithContext is an alternate form of the GetDiscoveredKmsInstances method which supports a Context parameter
func (schematics *SchematicsV1) GetDiscoveredKmsInstancesWithContext(ctx context.Context, getDiscoveredKmsInstancesOptions *GetDiscoveredKmsInstancesOptions) (result *KMSDiscovery, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDiscoveredKmsInstancesOptions, "getDiscoveredKmsInstancesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getDiscoveredKmsInstancesOptions, "getDiscoveredKmsInstancesOptions")
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

	for headerName, headerValue := range getDiscoveredKmsInstancesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("schematics", "V1", "GetDiscoveredKmsInstances")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("encryption_scheme", fmt.Sprint(*getDiscoveredKmsInstancesOptions.EncryptionScheme))
	builder.AddQuery("location", fmt.Sprint(*getDiscoveredKmsInstancesOptions.Location))
	if getDiscoveredKmsInstancesOptions.ResourceGroup != nil {
		builder.AddQuery("resource_group", fmt.Sprint(*getDiscoveredKmsInstancesOptions.ResourceGroup))
	}
	if getDiscoveredKmsInstancesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*getDiscoveredKmsInstancesOptions.Limit))
	}
	if getDiscoveredKmsInstancesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*getDiscoveredKmsInstancesOptions.Sort))
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
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKMSDiscovery)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// Action : Complete Action details with user inputs and system generated data.
type Action struct {
	// Action name (unique for an account).
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Action.  By default, Action will be created in Default Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the README file, for the source.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

	// Inventory of host and host group for the playbook, in .ini file format.
	TargetsIni *string `json:"targets_ini,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Id to the Trigger.
	TriggerRecordID *string `json:"trigger_record_id,omitempty"`

	// Action Id.
	ID *string `json:"id,omitempty"`

	// Action Cloud Resource Name.
	Crn *string `json:"crn,omitempty"`

	// Action account id.
	Account *string `json:"account,omitempty"`

	// Action Playbook Source creation time.
	SourceCreatedAt *strfmt.DateTime `json:"source_created_at,omitempty"`

	// Email address of user who created the Action Playbook Source.
	SourceCreatedBy *string `json:"source_created_by,omitempty"`

	// Action Playbook updation time.
	SourceUpdatedAt *strfmt.DateTime `json:"source_updated_at,omitempty"`

	// Email address of user who updated the Action Playbook Source.
	SourceUpdatedBy *string `json:"source_updated_by,omitempty"`

	// Action creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Email address of user who created the action.
	CreatedBy *string `json:"created_by,omitempty"`

	// Action updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Email address of user who updated the action.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// name of the namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// Playbook names retrieved from repo.
	PlaybookNames []string `json:"playbook_names,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`
}

// Constants associated with the Action.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	Action_Location_EuDe    = "eu_de"
	Action_Location_EuGb    = "eu_gb"
	Action_Location_UsEast  = "us_east"
	Action_Location_UsSouth = "us_south"
)

// Constants associated with the Action.SourceType property.
// Type of source for the Template.
const (
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
	err = core.UnmarshalModel(m, "bastion", &obj.Bastion, UnmarshalTargetResourceset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "targets_ini", &obj.TargetsIni)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalVariableData)
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
	err = core.UnmarshalPrimitive(m, "trigger_record_id", &obj.TriggerRecordID)
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
	err = core.UnmarshalPrimitive(m, "namespace", &obj.Namespace)
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

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
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
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	ActionLite_Location_EuDe    = "eu_de"
	ActionLite_Location_EuGb    = "eu_gb"
	ActionLite_Location_UsEast  = "us_east"
	ActionLite_Location_UsSouth = "us_south"
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
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// Action Options Template ...
	ActionOptions *WorkspaceActivityOptionsTemplate `json:"action_options,omitempty"`

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
func (options *ApplyWorkspaceCommandOptions) SetWID(wID string) *ApplyWorkspaceCommandOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *ApplyWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *ApplyWorkspaceCommandOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetActionOptions : Allow user to set ActionOptions
func (options *ApplyWorkspaceCommandOptions) SetActionOptions(actionOptions *WorkspaceActivityOptionsTemplate) *ApplyWorkspaceCommandOptions {
	options.ActionOptions = actionOptions
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ApplyWorkspaceCommandOptions) SetHeaders(param map[string]string) *ApplyWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// CatalogRef : CatalogRef -.
type CatalogRef struct {
	// Dry run.
	DryRun *bool `json:"dry_run,omitempty"`

	// Catalog item icon url.
	ItemIconURL *string `json:"item_icon_url,omitempty"`

	// Catalog item id.
	ItemID *string `json:"item_id,omitempty"`

	// Catalog item name.
	ItemName *string `json:"item_name,omitempty"`

	// Catalog item readme url.
	ItemReadmeURL *string `json:"item_readme_url,omitempty"`

	// Catalog item url.
	ItemURL *string `json:"item_url,omitempty"`

	// Catalog item launch url.
	LaunchURL *string `json:"launch_url,omitempty"`

	// Catalog item offering version.
	OfferingVersion *string `json:"offering_version,omitempty"`
}

// UnmarshalCatalogRef unmarshals an instance of CatalogRef from the specified map of raw messages.
func UnmarshalCatalogRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogRef)
	err = core.UnmarshalPrimitive(m, "dry_run", &obj.DryRun)
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
	// Action name (unique for an account).
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Action.  By default, Action will be created in Default Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the README file, for the source.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

	// Inventory of host and host group for the playbook, in .ini file format.
	TargetsIni *string `json:"targets_ini,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Id to the Trigger.
	TriggerRecordID *string `json:"trigger_record_id,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// The github token associated with the GIT. Required for cloning of repo.
	XGithubToken *string `json:"X-Github-token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateActionOptions.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	CreateActionOptions_Location_EuDe    = "eu_de"
	CreateActionOptions_Location_EuGb    = "eu_gb"
	CreateActionOptions_Location_UsEast  = "us_east"
	CreateActionOptions_Location_UsSouth = "us_south"
)

// Constants associated with the CreateActionOptions.SourceType property.
// Type of source for the Template.
const (
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
func (options *CreateActionOptions) SetName(name string) *CreateActionOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateActionOptions) SetDescription(description string) *CreateActionOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetLocation : Allow user to set Location
func (options *CreateActionOptions) SetLocation(location string) *CreateActionOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *CreateActionOptions) SetResourceGroup(resourceGroup string) *CreateActionOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateActionOptions) SetTags(tags []string) *CreateActionOptions {
	options.Tags = tags
	return options
}

// SetUserState : Allow user to set UserState
func (options *CreateActionOptions) SetUserState(userState *UserState) *CreateActionOptions {
	options.UserState = userState
	return options
}

// SetSourceReadmeURL : Allow user to set SourceReadmeURL
func (options *CreateActionOptions) SetSourceReadmeURL(sourceReadmeURL string) *CreateActionOptions {
	options.SourceReadmeURL = core.StringPtr(sourceReadmeURL)
	return options
}

// SetSource : Allow user to set Source
func (options *CreateActionOptions) SetSource(source *ExternalSource) *CreateActionOptions {
	options.Source = source
	return options
}

// SetSourceType : Allow user to set SourceType
func (options *CreateActionOptions) SetSourceType(sourceType string) *CreateActionOptions {
	options.SourceType = core.StringPtr(sourceType)
	return options
}

// SetCommandParameter : Allow user to set CommandParameter
func (options *CreateActionOptions) SetCommandParameter(commandParameter string) *CreateActionOptions {
	options.CommandParameter = core.StringPtr(commandParameter)
	return options
}

// SetBastion : Allow user to set Bastion
func (options *CreateActionOptions) SetBastion(bastion *TargetResourceset) *CreateActionOptions {
	options.Bastion = bastion
	return options
}

// SetTargetsIni : Allow user to set TargetsIni
func (options *CreateActionOptions) SetTargetsIni(targetsIni string) *CreateActionOptions {
	options.TargetsIni = core.StringPtr(targetsIni)
	return options
}

// SetCredentials : Allow user to set Credentials
func (options *CreateActionOptions) SetCredentials(credentials []VariableData) *CreateActionOptions {
	options.Credentials = credentials
	return options
}

// SetInputs : Allow user to set Inputs
func (options *CreateActionOptions) SetInputs(inputs []VariableData) *CreateActionOptions {
	options.Inputs = inputs
	return options
}

// SetOutputs : Allow user to set Outputs
func (options *CreateActionOptions) SetOutputs(outputs []VariableData) *CreateActionOptions {
	options.Outputs = outputs
	return options
}

// SetSettings : Allow user to set Settings
func (options *CreateActionOptions) SetSettings(settings []VariableData) *CreateActionOptions {
	options.Settings = settings
	return options
}

// SetTriggerRecordID : Allow user to set TriggerRecordID
func (options *CreateActionOptions) SetTriggerRecordID(triggerRecordID string) *CreateActionOptions {
	options.TriggerRecordID = core.StringPtr(triggerRecordID)
	return options
}

// SetState : Allow user to set State
func (options *CreateActionOptions) SetState(state *ActionState) *CreateActionOptions {
	options.State = state
	return options
}

// SetSysLock : Allow user to set SysLock
func (options *CreateActionOptions) SetSysLock(sysLock *SystemLock) *CreateActionOptions {
	options.SysLock = sysLock
	return options
}

// SetXGithubToken : Allow user to set XGithubToken
func (options *CreateActionOptions) SetXGithubToken(xGithubToken string) *CreateActionOptions {
	options.XGithubToken = core.StringPtr(xGithubToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateActionOptions) SetHeaders(param map[string]string) *CreateActionOptions {
	options.Headers = param
	return options
}

// CreateJobOptions : The CreateJob options.
type CreateJobOptions struct {
	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id or control-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job data.
	Data *JobData `json:"data,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateJobOptions.CommandObject property.
// Name of the Schematics automation resource.
const (
	CreateJobOptions_CommandObject_Action    = "action"
	CreateJobOptions_CommandObject_Workspace = "workspace"
)

// Constants associated with the CreateJobOptions.CommandName property.
// Schematics job command name.
const (
	CreateJobOptions_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	CreateJobOptions_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	CreateJobOptions_CommandName_HelmInstall          = "helm_install"
	CreateJobOptions_CommandName_HelmList             = "helm_list"
	CreateJobOptions_CommandName_HelmShow             = "helm_show"
	CreateJobOptions_CommandName_OpaEvaluate          = "opa_evaluate"
	CreateJobOptions_CommandName_TerraformInit        = "terraform_init"
	CreateJobOptions_CommandName_TerrformApply        = "terrform_apply"
	CreateJobOptions_CommandName_TerrformDestroy      = "terrform_destroy"
	CreateJobOptions_CommandName_TerrformPlan         = "terrform_plan"
	CreateJobOptions_CommandName_TerrformRefresh      = "terrform_refresh"
	CreateJobOptions_CommandName_TerrformShow         = "terrform_show"
	CreateJobOptions_CommandName_TerrformTaint        = "terrform_taint"
	CreateJobOptions_CommandName_WorkspaceApplyFlow   = "workspace_apply_flow"
	CreateJobOptions_CommandName_WorkspaceCustomFlow  = "workspace_custom_flow"
	CreateJobOptions_CommandName_WorkspaceDestroyFlow = "workspace_destroy_flow"
	CreateJobOptions_CommandName_WorkspaceInitFlow    = "workspace_init_flow"
	CreateJobOptions_CommandName_WorkspacePlanFlow    = "workspace_plan_flow"
	CreateJobOptions_CommandName_WorkspaceRefreshFlow = "workspace_refresh_flow"
	CreateJobOptions_CommandName_WorkspaceShowFlow    = "workspace_show_flow"
)

// Constants associated with the CreateJobOptions.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	CreateJobOptions_Location_EuDe    = "eu_de"
	CreateJobOptions_Location_EuGb    = "eu_gb"
	CreateJobOptions_Location_UsEast  = "us_east"
	CreateJobOptions_Location_UsSouth = "us_south"
)

// NewCreateJobOptions : Instantiate CreateJobOptions
func (*SchematicsV1) NewCreateJobOptions(refreshToken string) *CreateJobOptions {
	return &CreateJobOptions{
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *CreateJobOptions) SetRefreshToken(refreshToken string) *CreateJobOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetCommandObject : Allow user to set CommandObject
func (options *CreateJobOptions) SetCommandObject(commandObject string) *CreateJobOptions {
	options.CommandObject = core.StringPtr(commandObject)
	return options
}

// SetCommandObjectID : Allow user to set CommandObjectID
func (options *CreateJobOptions) SetCommandObjectID(commandObjectID string) *CreateJobOptions {
	options.CommandObjectID = core.StringPtr(commandObjectID)
	return options
}

// SetCommandName : Allow user to set CommandName
func (options *CreateJobOptions) SetCommandName(commandName string) *CreateJobOptions {
	options.CommandName = core.StringPtr(commandName)
	return options
}

// SetCommandParameter : Allow user to set CommandParameter
func (options *CreateJobOptions) SetCommandParameter(commandParameter string) *CreateJobOptions {
	options.CommandParameter = core.StringPtr(commandParameter)
	return options
}

// SetCommandOptions : Allow user to set CommandOptions
func (options *CreateJobOptions) SetCommandOptions(commandOptions []string) *CreateJobOptions {
	options.CommandOptions = commandOptions
	return options
}

// SetInputs : Allow user to set Inputs
func (options *CreateJobOptions) SetInputs(inputs []VariableData) *CreateJobOptions {
	options.Inputs = inputs
	return options
}

// SetSettings : Allow user to set Settings
func (options *CreateJobOptions) SetSettings(settings []VariableData) *CreateJobOptions {
	options.Settings = settings
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateJobOptions) SetTags(tags []string) *CreateJobOptions {
	options.Tags = tags
	return options
}

// SetLocation : Allow user to set Location
func (options *CreateJobOptions) SetLocation(location string) *CreateJobOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetStatus : Allow user to set Status
func (options *CreateJobOptions) SetStatus(status *JobStatus) *CreateJobOptions {
	options.Status = status
	return options
}

// SetData : Allow user to set Data
func (options *CreateJobOptions) SetData(data *JobData) *CreateJobOptions {
	options.Data = data
	return options
}

// SetBastion : Allow user to set Bastion
func (options *CreateJobOptions) SetBastion(bastion *TargetResourceset) *CreateJobOptions {
	options.Bastion = bastion
	return options
}

// SetLogSummary : Allow user to set LogSummary
func (options *CreateJobOptions) SetLogSummary(logSummary *JobLogSummary) *CreateJobOptions {
	options.LogSummary = logSummary
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateJobOptions) SetHeaders(param map[string]string) *CreateJobOptions {
	options.Headers = param
	return options
}

// CreateSharedDatasetOptions : The CreateSharedDataset options.
type CreateSharedDatasetOptions struct {
	// Automatically propagate changes to consumers.
	AutoPropagateChange *bool `json:"auto_propagate_change,omitempty"`

	// Dataset description.
	Description *string `json:"description,omitempty"`

	// Affected workspaces.
	EffectedWorkspaceIds []string `json:"effected_workspace_ids,omitempty"`

	// Resource group name.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Shared dataset data.
	SharedDatasetData []SharedDatasetData `json:"shared_dataset_data,omitempty"`

	// Shared dataset name.
	SharedDatasetName *string `json:"shared_dataset_name,omitempty"`

	// Shared dataset source name.
	SharedDatasetSourceName *string `json:"shared_dataset_source_name,omitempty"`

	// Shared dataset type.
	SharedDatasetType []string `json:"shared_dataset_type,omitempty"`

	// Shared dataset tags.
	Tags []string `json:"tags,omitempty"`

	// Shared dataset version.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSharedDatasetOptions : Instantiate CreateSharedDatasetOptions
func (*SchematicsV1) NewCreateSharedDatasetOptions() *CreateSharedDatasetOptions {
	return &CreateSharedDatasetOptions{}
}

// SetAutoPropagateChange : Allow user to set AutoPropagateChange
func (options *CreateSharedDatasetOptions) SetAutoPropagateChange(autoPropagateChange bool) *CreateSharedDatasetOptions {
	options.AutoPropagateChange = core.BoolPtr(autoPropagateChange)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateSharedDatasetOptions) SetDescription(description string) *CreateSharedDatasetOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEffectedWorkspaceIds : Allow user to set EffectedWorkspaceIds
func (options *CreateSharedDatasetOptions) SetEffectedWorkspaceIds(effectedWorkspaceIds []string) *CreateSharedDatasetOptions {
	options.EffectedWorkspaceIds = effectedWorkspaceIds
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *CreateSharedDatasetOptions) SetResourceGroup(resourceGroup string) *CreateSharedDatasetOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetSharedDatasetData : Allow user to set SharedDatasetData
func (options *CreateSharedDatasetOptions) SetSharedDatasetData(sharedDatasetData []SharedDatasetData) *CreateSharedDatasetOptions {
	options.SharedDatasetData = sharedDatasetData
	return options
}

// SetSharedDatasetName : Allow user to set SharedDatasetName
func (options *CreateSharedDatasetOptions) SetSharedDatasetName(sharedDatasetName string) *CreateSharedDatasetOptions {
	options.SharedDatasetName = core.StringPtr(sharedDatasetName)
	return options
}

// SetSharedDatasetSourceName : Allow user to set SharedDatasetSourceName
func (options *CreateSharedDatasetOptions) SetSharedDatasetSourceName(sharedDatasetSourceName string) *CreateSharedDatasetOptions {
	options.SharedDatasetSourceName = core.StringPtr(sharedDatasetSourceName)
	return options
}

// SetSharedDatasetType : Allow user to set SharedDatasetType
func (options *CreateSharedDatasetOptions) SetSharedDatasetType(sharedDatasetType []string) *CreateSharedDatasetOptions {
	options.SharedDatasetType = sharedDatasetType
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateSharedDatasetOptions) SetTags(tags []string) *CreateSharedDatasetOptions {
	options.Tags = tags
	return options
}

// SetVersion : Allow user to set Version
func (options *CreateSharedDatasetOptions) SetVersion(version string) *CreateSharedDatasetOptions {
	options.Version = core.StringPtr(version)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSharedDatasetOptions) SetHeaders(param map[string]string) *CreateSharedDatasetOptions {
	options.Headers = param
	return options
}

// CreateWorkspaceDeletionJobOptions : The CreateWorkspaceDeletionJob options.
type CreateWorkspaceDeletionJobOptions struct {
	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// True to delete workspace.
	NewDeleteWorkspaces *bool `json:"new_delete_workspaces,omitempty"`

	// True to destroy the resources managed by this workspace.
	NewDestroyResources *bool `json:"new_destroy_resources,omitempty"`

	// Workspace deletion job name.
	NewJob *string `json:"new_job,omitempty"`

	// Version.
	NewVersion *string `json:"new_version,omitempty"`

	// List of workspaces to be deleted.
	NewWorkspaces []string `json:"new_workspaces,omitempty"`

	// true or 1 - to destroy resources before deleting workspace;  If this is true, refresh_token is mandatory.
	DestroyResources *string `json:"destroy_resources,omitempty"`

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
func (options *CreateWorkspaceDeletionJobOptions) SetRefreshToken(refreshToken string) *CreateWorkspaceDeletionJobOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetNewDeleteWorkspaces : Allow user to set NewDeleteWorkspaces
func (options *CreateWorkspaceDeletionJobOptions) SetNewDeleteWorkspaces(newDeleteWorkspaces bool) *CreateWorkspaceDeletionJobOptions {
	options.NewDeleteWorkspaces = core.BoolPtr(newDeleteWorkspaces)
	return options
}

// SetNewDestroyResources : Allow user to set NewDestroyResources
func (options *CreateWorkspaceDeletionJobOptions) SetNewDestroyResources(newDestroyResources bool) *CreateWorkspaceDeletionJobOptions {
	options.NewDestroyResources = core.BoolPtr(newDestroyResources)
	return options
}

// SetNewJob : Allow user to set NewJob
func (options *CreateWorkspaceDeletionJobOptions) SetNewJob(newJob string) *CreateWorkspaceDeletionJobOptions {
	options.NewJob = core.StringPtr(newJob)
	return options
}

// SetNewVersion : Allow user to set NewVersion
func (options *CreateWorkspaceDeletionJobOptions) SetNewVersion(newVersion string) *CreateWorkspaceDeletionJobOptions {
	options.NewVersion = core.StringPtr(newVersion)
	return options
}

// SetNewWorkspaces : Allow user to set NewWorkspaces
func (options *CreateWorkspaceDeletionJobOptions) SetNewWorkspaces(newWorkspaces []string) *CreateWorkspaceDeletionJobOptions {
	options.NewWorkspaces = newWorkspaces
	return options
}

// SetDestroyResources : Allow user to set DestroyResources
func (options *CreateWorkspaceDeletionJobOptions) SetDestroyResources(destroyResources string) *CreateWorkspaceDeletionJobOptions {
	options.DestroyResources = core.StringPtr(destroyResources)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateWorkspaceDeletionJobOptions) SetHeaders(param map[string]string) *CreateWorkspaceDeletionJobOptions {
	options.Headers = param
	return options
}

// CreateWorkspaceOptions : The CreateWorkspace options.
type CreateWorkspaceOptions struct {
	// List of applied shared dataset id.
	AppliedShareddataIds []string `json:"applied_shareddata_ids,omitempty"`

	// CatalogRef -.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// Workspace description.
	Description *string `json:"description,omitempty"`

	// Workspace location.
	Location *string `json:"location,omitempty"`

	// Workspace name.
	Name *string `json:"name,omitempty"`

	// Workspace resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// SharedTargetData -.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// Workspace tags.
	Tags []string `json:"tags,omitempty"`

	// TemplateData -.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// Workspace template ref.
	TemplateRef *string `json:"template_ref,omitempty"`

	// TemplateRepoRequest -.
	TemplateRepo *TemplateRepoRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// WorkspaceStatusRequest -.
	WorkspaceStatus *WorkspaceStatusRequest `json:"workspace_status,omitempty"`

	// The github token associated with the GIT. Required for cloning of repo.
	XGithubToken *string `json:"X-Github-token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateWorkspaceOptions : Instantiate CreateWorkspaceOptions
func (*SchematicsV1) NewCreateWorkspaceOptions() *CreateWorkspaceOptions {
	return &CreateWorkspaceOptions{}
}

// SetAppliedShareddataIds : Allow user to set AppliedShareddataIds
func (options *CreateWorkspaceOptions) SetAppliedShareddataIds(appliedShareddataIds []string) *CreateWorkspaceOptions {
	options.AppliedShareddataIds = appliedShareddataIds
	return options
}

// SetCatalogRef : Allow user to set CatalogRef
func (options *CreateWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *CreateWorkspaceOptions {
	options.CatalogRef = catalogRef
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateWorkspaceOptions) SetDescription(description string) *CreateWorkspaceOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetLocation : Allow user to set Location
func (options *CreateWorkspaceOptions) SetLocation(location string) *CreateWorkspaceOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetName : Allow user to set Name
func (options *CreateWorkspaceOptions) SetName(name string) *CreateWorkspaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *CreateWorkspaceOptions) SetResourceGroup(resourceGroup string) *CreateWorkspaceOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetSharedData : Allow user to set SharedData
func (options *CreateWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *CreateWorkspaceOptions {
	options.SharedData = sharedData
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateWorkspaceOptions) SetTags(tags []string) *CreateWorkspaceOptions {
	options.Tags = tags
	return options
}

// SetTemplateData : Allow user to set TemplateData
func (options *CreateWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *CreateWorkspaceOptions {
	options.TemplateData = templateData
	return options
}

// SetTemplateRef : Allow user to set TemplateRef
func (options *CreateWorkspaceOptions) SetTemplateRef(templateRef string) *CreateWorkspaceOptions {
	options.TemplateRef = core.StringPtr(templateRef)
	return options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (options *CreateWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoRequest) *CreateWorkspaceOptions {
	options.TemplateRepo = templateRepo
	return options
}

// SetType : Allow user to set Type
func (options *CreateWorkspaceOptions) SetType(typeVar []string) *CreateWorkspaceOptions {
	options.Type = typeVar
	return options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (options *CreateWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusRequest) *CreateWorkspaceOptions {
	options.WorkspaceStatus = workspaceStatus
	return options
}

// SetXGithubToken : Allow user to set XGithubToken
func (options *CreateWorkspaceOptions) SetXGithubToken(xGithubToken string) *CreateWorkspaceOptions {
	options.XGithubToken = core.StringPtr(xGithubToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateWorkspaceOptions) SetHeaders(param map[string]string) *CreateWorkspaceOptions {
	options.Headers = param
	return options
}

// DeleteActionOptions : The DeleteAction options.
type DeleteActionOptions struct {
	// Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.
	ActionID *string `json:"action_id" validate:"required,ne="`

	// Equivalent to -force options in the command line.
	Force *bool `json:"force,omitempty"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"propagate,omitempty"`

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
func (options *DeleteActionOptions) SetActionID(actionID string) *DeleteActionOptions {
	options.ActionID = core.StringPtr(actionID)
	return options
}

// SetForce : Allow user to set Force
func (options *DeleteActionOptions) SetForce(force bool) *DeleteActionOptions {
	options.Force = core.BoolPtr(force)
	return options
}

// SetPropagate : Allow user to set Propagate
func (options *DeleteActionOptions) SetPropagate(propagate bool) *DeleteActionOptions {
	options.Propagate = core.BoolPtr(propagate)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteActionOptions) SetHeaders(param map[string]string) *DeleteActionOptions {
	options.Headers = param
	return options
}

// DeleteJobOptions : The DeleteJob options.
type DeleteJobOptions struct {
	// Job Id. Use GET /jobs API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"job_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// Equivalent to -force options in the command line.
	Force *bool `json:"force,omitempty"`

	// Auto propagate the chaange or deletion to the dependent resources.
	Propagate *bool `json:"propagate,omitempty"`

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
func (options *DeleteJobOptions) SetJobID(jobID string) *DeleteJobOptions {
	options.JobID = core.StringPtr(jobID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *DeleteJobOptions) SetRefreshToken(refreshToken string) *DeleteJobOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetForce : Allow user to set Force
func (options *DeleteJobOptions) SetForce(force bool) *DeleteJobOptions {
	options.Force = core.BoolPtr(force)
	return options
}

// SetPropagate : Allow user to set Propagate
func (options *DeleteJobOptions) SetPropagate(propagate bool) *DeleteJobOptions {
	options.Propagate = core.BoolPtr(propagate)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteJobOptions) SetHeaders(param map[string]string) *DeleteJobOptions {
	options.Headers = param
	return options
}

// DeleteSharedDatasetOptions : The DeleteSharedDataset options.
type DeleteSharedDatasetOptions struct {
	// The shared dataset ID Use the GET /shared_datasets to look up the shared dataset IDs  in your IBM Cloud account.
	SdID *string `json:"sd_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSharedDatasetOptions : Instantiate DeleteSharedDatasetOptions
func (*SchematicsV1) NewDeleteSharedDatasetOptions(sdID string) *DeleteSharedDatasetOptions {
	return &DeleteSharedDatasetOptions{
		SdID: core.StringPtr(sdID),
	}
}

// SetSdID : Allow user to set SdID
func (options *DeleteSharedDatasetOptions) SetSdID(sdID string) *DeleteSharedDatasetOptions {
	options.SdID = core.StringPtr(sdID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSharedDatasetOptions) SetHeaders(param map[string]string) *DeleteSharedDatasetOptions {
	options.Headers = param
	return options
}

// DeleteWorkspaceActivityOptions : The DeleteWorkspaceActivity options.
type DeleteWorkspaceActivityOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The activity ID that you want to see additional details.
	ActivityID *string `json:"activity_id" validate:"required,ne="`

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
func (options *DeleteWorkspaceActivityOptions) SetWID(wID string) *DeleteWorkspaceActivityOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetActivityID : Allow user to set ActivityID
func (options *DeleteWorkspaceActivityOptions) SetActivityID(activityID string) *DeleteWorkspaceActivityOptions {
	options.ActivityID = core.StringPtr(activityID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWorkspaceActivityOptions) SetHeaders(param map[string]string) *DeleteWorkspaceActivityOptions {
	options.Headers = param
	return options
}

// DeleteWorkspaceOptions : The DeleteWorkspace options.
type DeleteWorkspaceOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// true or 1 - to destroy resources before deleting workspace;  If this is true, refresh_token is mandatory.
	DestroyResources *string `json:"destroy_resources,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteWorkspaceOptions : Instantiate DeleteWorkspaceOptions
func (*SchematicsV1) NewDeleteWorkspaceOptions(wID string, refreshToken string) *DeleteWorkspaceOptions {
	return &DeleteWorkspaceOptions{
		WID:          core.StringPtr(wID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetWID : Allow user to set WID
func (options *DeleteWorkspaceOptions) SetWID(wID string) *DeleteWorkspaceOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *DeleteWorkspaceOptions) SetRefreshToken(refreshToken string) *DeleteWorkspaceOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetDestroyResources : Allow user to set DestroyResources
func (options *DeleteWorkspaceOptions) SetDestroyResources(destroyResources string) *DeleteWorkspaceOptions {
	options.DestroyResources = core.StringPtr(destroyResources)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWorkspaceOptions) SetHeaders(param map[string]string) *DeleteWorkspaceOptions {
	options.Headers = param
	return options
}

// DestroyWorkspaceCommandOptions : The DestroyWorkspaceCommand options.
type DestroyWorkspaceCommandOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// Action Options Template ...
	ActionOptions *WorkspaceActivityOptionsTemplate `json:"action_options,omitempty"`

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
func (options *DestroyWorkspaceCommandOptions) SetWID(wID string) *DestroyWorkspaceCommandOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *DestroyWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *DestroyWorkspaceCommandOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetActionOptions : Allow user to set ActionOptions
func (options *DestroyWorkspaceCommandOptions) SetActionOptions(actionOptions *WorkspaceActivityOptionsTemplate) *DestroyWorkspaceCommandOptions {
	options.ActionOptions = actionOptions
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DestroyWorkspaceCommandOptions) SetHeaders(param map[string]string) *DestroyWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// EnvVariableResponse : EnvVariableResponse -.
type EnvVariableResponse struct {
	// Env variable is hidden.
	Hidden *bool `json:"hidden,omitempty"`

	// Env variable name.
	Name *string `json:"name,omitempty"`

	// Env variable is secure.
	Secure *bool `json:"secure,omitempty"`

	// Value for env variable.
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

// ExternalSource : Source of templates, playbooks, or controls.
type ExternalSource struct {
	// Type of source for the Template.
	SourceType *string `json:"source_type" validate:"required"`

	// Connection details to Git source.
	Git *ExternalSourceGit `json:"git,omitempty"`
}

// Constants associated with the ExternalSource.SourceType property.
// Type of source for the Template.
const (
	ExternalSource_SourceType_ExternalScm      = "external_scm"
	ExternalSource_SourceType_GitHub           = "git_hub"
	ExternalSource_SourceType_GitHubEnterprise = "git_hub_enterprise"
	ExternalSource_SourceType_GitLab           = "git_lab"
	ExternalSource_SourceType_IbmCloudCatalog  = "ibm_cloud_catalog"
	ExternalSource_SourceType_IbmGitLab        = "ibm_git_lab"
	ExternalSource_SourceType_Local            = "local"
)

// NewExternalSource : Instantiate ExternalSource (Generic Model Constructor)
func (*SchematicsV1) NewExternalSource(sourceType string) (model *ExternalSource, err error) {
	model = &ExternalSource{
		SourceType: core.StringPtr(sourceType),
	}
	err = core.ValidateStruct(model, "required parameters")
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalSourceGit : Connection details to Git source.
type ExternalSourceGit struct {
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
	ActionID *string `json:"action_id" validate:"required,ne="`

	// Level of details returned by the get method.
	Profile *string `json:"profile,omitempty"`

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
func (options *GetActionOptions) SetActionID(actionID string) *GetActionOptions {
	options.ActionID = core.StringPtr(actionID)
	return options
}

// SetProfile : Allow user to set Profile
func (options *GetActionOptions) SetProfile(profile string) *GetActionOptions {
	options.Profile = core.StringPtr(profile)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetActionOptions) SetHeaders(param map[string]string) *GetActionOptions {
	options.Headers = param
	return options
}

// GetAllWorkspaceInputsOptions : The GetAllWorkspaceInputs options.
type GetAllWorkspaceInputsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetAllWorkspaceInputsOptions) SetWID(wID string) *GetAllWorkspaceInputsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllWorkspaceInputsOptions) SetHeaders(param map[string]string) *GetAllWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// GetDiscoveredKmsInstancesOptions : The GetDiscoveredKmsInstances options.
type GetDiscoveredKmsInstancesOptions struct {
	// The encryption scheme to be used.
	EncryptionScheme *string `json:"encryption_scheme" validate:"required"`

	// The location of the Resource.
	Location *string `json:"location" validate:"required"`

	// The resource group (by default, fetch from all resource groups).
	ResourceGroup *string `json:"resource_group,omitempty"`

	// The numbers of items to return.
	Limit *int64 `json:"limit,omitempty"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDiscoveredKmsInstancesOptions : Instantiate GetDiscoveredKmsInstancesOptions
func (*SchematicsV1) NewGetDiscoveredKmsInstancesOptions(encryptionScheme string, location string) *GetDiscoveredKmsInstancesOptions {
	return &GetDiscoveredKmsInstancesOptions{
		EncryptionScheme: core.StringPtr(encryptionScheme),
		Location:         core.StringPtr(location),
	}
}

// SetEncryptionScheme : Allow user to set EncryptionScheme
func (options *GetDiscoveredKmsInstancesOptions) SetEncryptionScheme(encryptionScheme string) *GetDiscoveredKmsInstancesOptions {
	options.EncryptionScheme = core.StringPtr(encryptionScheme)
	return options
}

// SetLocation : Allow user to set Location
func (options *GetDiscoveredKmsInstancesOptions) SetLocation(location string) *GetDiscoveredKmsInstancesOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *GetDiscoveredKmsInstancesOptions) SetResourceGroup(resourceGroup string) *GetDiscoveredKmsInstancesOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetLimit : Allow user to set Limit
func (options *GetDiscoveredKmsInstancesOptions) SetLimit(limit int64) *GetDiscoveredKmsInstancesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetSort : Allow user to set Sort
func (options *GetDiscoveredKmsInstancesOptions) SetSort(sort string) *GetDiscoveredKmsInstancesOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetDiscoveredKmsInstancesOptions) SetHeaders(param map[string]string) *GetDiscoveredKmsInstancesOptions {
	options.Headers = param
	return options
}

// GetJobOptions : The GetJob options.
type GetJobOptions struct {
	// Job Id. Use GET /jobs API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Level of details returned by the get method.
	Profile *string `json:"profile,omitempty"`

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
func (options *GetJobOptions) SetJobID(jobID string) *GetJobOptions {
	options.JobID = core.StringPtr(jobID)
	return options
}

// SetProfile : Allow user to set Profile
func (options *GetJobOptions) SetProfile(profile string) *GetJobOptions {
	options.Profile = core.StringPtr(profile)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetJobOptions) SetHeaders(param map[string]string) *GetJobOptions {
	options.Headers = param
	return options
}

// GetKmsSettingsOptions : The GetKmsSettings options.
type GetKmsSettingsOptions struct {
	// The location of the Resource.
	Location *string `json:"location" validate:"required"`

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
func (options *GetKmsSettingsOptions) SetLocation(location string) *GetKmsSettingsOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetKmsSettingsOptions) SetHeaders(param map[string]string) *GetKmsSettingsOptions {
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

// GetSharedDatasetOptions : The GetSharedDataset options.
type GetSharedDatasetOptions struct {
	// The shared dataset ID Use the GET /shared_datasets to look up the shared dataset IDs  in your IBM Cloud account.
	SdID *string `json:"sd_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSharedDatasetOptions : Instantiate GetSharedDatasetOptions
func (*SchematicsV1) NewGetSharedDatasetOptions(sdID string) *GetSharedDatasetOptions {
	return &GetSharedDatasetOptions{
		SdID: core.StringPtr(sdID),
	}
}

// SetSdID : Allow user to set SdID
func (options *GetSharedDatasetOptions) SetSdID(sdID string) *GetSharedDatasetOptions {
	options.SdID = core.StringPtr(sdID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSharedDatasetOptions) SetHeaders(param map[string]string) *GetSharedDatasetOptions {
	options.Headers = param
	return options
}

// GetTemplateActivityLogOptions : The GetTemplateActivityLog options.
type GetTemplateActivityLogOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

	// The activity ID that you want to see additional details.
	ActivityID *string `json:"activity_id" validate:"required,ne="`

	// `false` will hide the terraform command header in the logs.
	LogTfCmd *bool `json:"log_tf_cmd,omitempty"`

	// `false` will hide all the terraform command prefix in the log statements.
	LogTfPrefix *bool `json:"log_tf_prefix,omitempty"`

	// `false` will hide all the null resource prefix in the log statements.
	LogTfNullResource *bool `json:"log_tf_null_resource,omitempty"`

	// `true` will format all logs to withhold the original format  of ansible output in the log statements.
	LogTfAnsible *bool `json:"log_tf_ansible,omitempty"`

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
func (options *GetTemplateActivityLogOptions) SetWID(wID string) *GetTemplateActivityLogOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *GetTemplateActivityLogOptions) SetTID(tID string) *GetTemplateActivityLogOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetActivityID : Allow user to set ActivityID
func (options *GetTemplateActivityLogOptions) SetActivityID(activityID string) *GetTemplateActivityLogOptions {
	options.ActivityID = core.StringPtr(activityID)
	return options
}

// SetLogTfCmd : Allow user to set LogTfCmd
func (options *GetTemplateActivityLogOptions) SetLogTfCmd(logTfCmd bool) *GetTemplateActivityLogOptions {
	options.LogTfCmd = core.BoolPtr(logTfCmd)
	return options
}

// SetLogTfPrefix : Allow user to set LogTfPrefix
func (options *GetTemplateActivityLogOptions) SetLogTfPrefix(logTfPrefix bool) *GetTemplateActivityLogOptions {
	options.LogTfPrefix = core.BoolPtr(logTfPrefix)
	return options
}

// SetLogTfNullResource : Allow user to set LogTfNullResource
func (options *GetTemplateActivityLogOptions) SetLogTfNullResource(logTfNullResource bool) *GetTemplateActivityLogOptions {
	options.LogTfNullResource = core.BoolPtr(logTfNullResource)
	return options
}

// SetLogTfAnsible : Allow user to set LogTfAnsible
func (options *GetTemplateActivityLogOptions) SetLogTfAnsible(logTfAnsible bool) *GetTemplateActivityLogOptions {
	options.LogTfAnsible = core.BoolPtr(logTfAnsible)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateActivityLogOptions) SetHeaders(param map[string]string) *GetTemplateActivityLogOptions {
	options.Headers = param
	return options
}

// GetTemplateLogsOptions : The GetTemplateLogs options.
type GetTemplateLogsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

	// `false` will hide the terraform command header in the logs.
	LogTfCmd *bool `json:"log_tf_cmd,omitempty"`

	// `false` will hide all the terraform command prefix in the log statements.
	LogTfPrefix *bool `json:"log_tf_prefix,omitempty"`

	// `false` will hide all the null resource prefix in the log statements.
	LogTfNullResource *bool `json:"log_tf_null_resource,omitempty"`

	// `true` will format all logs to withhold the original format  of ansible output in the log statements.
	LogTfAnsible *bool `json:"log_tf_ansible,omitempty"`

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
func (options *GetTemplateLogsOptions) SetWID(wID string) *GetTemplateLogsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *GetTemplateLogsOptions) SetTID(tID string) *GetTemplateLogsOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetLogTfCmd : Allow user to set LogTfCmd
func (options *GetTemplateLogsOptions) SetLogTfCmd(logTfCmd bool) *GetTemplateLogsOptions {
	options.LogTfCmd = core.BoolPtr(logTfCmd)
	return options
}

// SetLogTfPrefix : Allow user to set LogTfPrefix
func (options *GetTemplateLogsOptions) SetLogTfPrefix(logTfPrefix bool) *GetTemplateLogsOptions {
	options.LogTfPrefix = core.BoolPtr(logTfPrefix)
	return options
}

// SetLogTfNullResource : Allow user to set LogTfNullResource
func (options *GetTemplateLogsOptions) SetLogTfNullResource(logTfNullResource bool) *GetTemplateLogsOptions {
	options.LogTfNullResource = core.BoolPtr(logTfNullResource)
	return options
}

// SetLogTfAnsible : Allow user to set LogTfAnsible
func (options *GetTemplateLogsOptions) SetLogTfAnsible(logTfAnsible bool) *GetTemplateLogsOptions {
	options.LogTfAnsible = core.BoolPtr(logTfAnsible)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateLogsOptions) SetHeaders(param map[string]string) *GetTemplateLogsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceActivityLogsOptions : The GetWorkspaceActivityLogs options.
type GetWorkspaceActivityLogsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The activity ID that you want to see additional details.
	ActivityID *string `json:"activity_id" validate:"required,ne="`

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
func (options *GetWorkspaceActivityLogsOptions) SetWID(wID string) *GetWorkspaceActivityLogsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetActivityID : Allow user to set ActivityID
func (options *GetWorkspaceActivityLogsOptions) SetActivityID(activityID string) *GetWorkspaceActivityLogsOptions {
	options.ActivityID = core.StringPtr(activityID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceActivityLogsOptions) SetHeaders(param map[string]string) *GetWorkspaceActivityLogsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceActivityOptions : The GetWorkspaceActivity options.
type GetWorkspaceActivityOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The activity ID that you want to see additional details.
	ActivityID *string `json:"activity_id" validate:"required,ne="`

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
func (options *GetWorkspaceActivityOptions) SetWID(wID string) *GetWorkspaceActivityOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetActivityID : Allow user to set ActivityID
func (options *GetWorkspaceActivityOptions) SetActivityID(activityID string) *GetWorkspaceActivityOptions {
	options.ActivityID = core.StringPtr(activityID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceActivityOptions) SetHeaders(param map[string]string) *GetWorkspaceActivityOptions {
	options.Headers = param
	return options
}

// GetWorkspaceDeletionJobStatusOptions : The GetWorkspaceDeletionJobStatus options.
type GetWorkspaceDeletionJobStatusOptions struct {
	// The workspace job deletion ID.
	WjID *string `json:"wj_id" validate:"required,ne="`

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
func (options *GetWorkspaceDeletionJobStatusOptions) SetWjID(wjID string) *GetWorkspaceDeletionJobStatusOptions {
	options.WjID = core.StringPtr(wjID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceDeletionJobStatusOptions) SetHeaders(param map[string]string) *GetWorkspaceDeletionJobStatusOptions {
	options.Headers = param
	return options
}

// GetWorkspaceInputMetadataOptions : The GetWorkspaceInputMetadata options.
type GetWorkspaceInputMetadataOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

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
func (options *GetWorkspaceInputMetadataOptions) SetWID(wID string) *GetWorkspaceInputMetadataOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *GetWorkspaceInputMetadataOptions) SetTID(tID string) *GetWorkspaceInputMetadataOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceInputMetadataOptions) SetHeaders(param map[string]string) *GetWorkspaceInputMetadataOptions {
	options.Headers = param
	return options
}

// GetWorkspaceInputsOptions : The GetWorkspaceInputs options.
type GetWorkspaceInputsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

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
func (options *GetWorkspaceInputsOptions) SetWID(wID string) *GetWorkspaceInputsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *GetWorkspaceInputsOptions) SetTID(tID string) *GetWorkspaceInputsOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceInputsOptions) SetHeaders(param map[string]string) *GetWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceLogUrlsOptions : The GetWorkspaceLogUrls options.
type GetWorkspaceLogUrlsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetWorkspaceLogUrlsOptions) SetWID(wID string) *GetWorkspaceLogUrlsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceLogUrlsOptions) SetHeaders(param map[string]string) *GetWorkspaceLogUrlsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceOptions : The GetWorkspace options.
type GetWorkspaceOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetWorkspaceOptions) SetWID(wID string) *GetWorkspaceOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceOptions) SetHeaders(param map[string]string) *GetWorkspaceOptions {
	options.Headers = param
	return options
}

// GetWorkspaceOutputsOptions : The GetWorkspaceOutputs options.
type GetWorkspaceOutputsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetWorkspaceOutputsOptions) SetWID(wID string) *GetWorkspaceOutputsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceOutputsOptions) SetHeaders(param map[string]string) *GetWorkspaceOutputsOptions {
	options.Headers = param
	return options
}

// GetWorkspaceReadmeOptions : The GetWorkspaceReadme options.
type GetWorkspaceReadmeOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The name of the commit/branch/tag.  Default, the repository’s default branch (usually master).
	Ref *string `json:"ref,omitempty"`

	// The format of the readme file.  Value ''markdown'' will give markdown, otherwise html.
	Formatted *string `json:"formatted,omitempty"`

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
func (options *GetWorkspaceReadmeOptions) SetWID(wID string) *GetWorkspaceReadmeOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRef : Allow user to set Ref
func (options *GetWorkspaceReadmeOptions) SetRef(ref string) *GetWorkspaceReadmeOptions {
	options.Ref = core.StringPtr(ref)
	return options
}

// SetFormatted : Allow user to set Formatted
func (options *GetWorkspaceReadmeOptions) SetFormatted(formatted string) *GetWorkspaceReadmeOptions {
	options.Formatted = core.StringPtr(formatted)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceReadmeOptions) SetHeaders(param map[string]string) *GetWorkspaceReadmeOptions {
	options.Headers = param
	return options
}

// GetWorkspaceResourcesOptions : The GetWorkspaceResources options.
type GetWorkspaceResourcesOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetWorkspaceResourcesOptions) SetWID(wID string) *GetWorkspaceResourcesOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceResourcesOptions) SetHeaders(param map[string]string) *GetWorkspaceResourcesOptions {
	options.Headers = param
	return options
}

// GetWorkspaceStateOptions : The GetWorkspaceState options.
type GetWorkspaceStateOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

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
func (options *GetWorkspaceStateOptions) SetWID(wID string) *GetWorkspaceStateOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceStateOptions) SetHeaders(param map[string]string) *GetWorkspaceStateOptions {
	options.Headers = param
	return options
}

// GetWorkspaceTemplateStateOptions : The GetWorkspaceTemplateState options.
type GetWorkspaceTemplateStateOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

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
func (options *GetWorkspaceTemplateStateOptions) SetWID(wID string) *GetWorkspaceTemplateStateOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *GetWorkspaceTemplateStateOptions) SetTID(tID string) *GetWorkspaceTemplateStateOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWorkspaceTemplateStateOptions) SetHeaders(param map[string]string) *GetWorkspaceTemplateStateOptions {
	options.Headers = param
	return options
}

// Job : Complete Job with user inputs and system generated data.
type Job struct {
	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id or control-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// Job ID.
	ID *string `json:"id,omitempty"`

	// Job name, uniquely derived from the related Action.
	Name *string `json:"name,omitempty"`

	// Job description derived from the related Action.
	Description *string `json:"description,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name derived from the related Action.
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

	// Inventory of host and host group for the playbook, in .ini file format.
	TargetsIni *string `json:"targets_ini,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

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
	Job_CommandObject_Action    = "action"
	Job_CommandObject_Workspace = "workspace"
)

// Constants associated with the Job.CommandName property.
// Schematics job command name.
const (
	Job_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	Job_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	Job_CommandName_HelmInstall          = "helm_install"
	Job_CommandName_HelmList             = "helm_list"
	Job_CommandName_HelmShow             = "helm_show"
	Job_CommandName_OpaEvaluate          = "opa_evaluate"
	Job_CommandName_TerraformInit        = "terraform_init"
	Job_CommandName_TerrformApply        = "terrform_apply"
	Job_CommandName_TerrformDestroy      = "terrform_destroy"
	Job_CommandName_TerrformPlan         = "terrform_plan"
	Job_CommandName_TerrformRefresh      = "terrform_refresh"
	Job_CommandName_TerrformShow         = "terrform_show"
	Job_CommandName_TerrformTaint        = "terrform_taint"
	Job_CommandName_WorkspaceApplyFlow   = "workspace_apply_flow"
	Job_CommandName_WorkspaceCustomFlow  = "workspace_custom_flow"
	Job_CommandName_WorkspaceDestroyFlow = "workspace_destroy_flow"
	Job_CommandName_WorkspaceInitFlow    = "workspace_init_flow"
	Job_CommandName_WorkspacePlanFlow    = "workspace_plan_flow"
	Job_CommandName_WorkspaceRefreshFlow = "workspace_refresh_flow"
	Job_CommandName_WorkspaceShowFlow    = "workspace_show_flow"
)

// Constants associated with the Job.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	Job_Location_EuDe    = "eu_de"
	Job_Location_EuGb    = "eu_gb"
	Job_Location_UsEast  = "us_east"
	Job_Location_UsSouth = "us_south"
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
	err = core.UnmarshalPrimitive(m, "targets_ini", &obj.TargetsIni)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bastion", &obj.Bastion, UnmarshalTargetResourceset)
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

	// Action Job data.
	ActionJobData *JobDataAction `json:"action_job_data,omitempty"`
}

// Constants associated with the JobData.JobType property.
// Type of Job.
const (
	JobData_JobType_ActionJob       = "action_job"
	JobData_JobType_RepoDownloadJob = "repo_download_job"
)

// NewJobData : Instantiate JobData (Generic Model Constructor)
func (*SchematicsV1) NewJobData(jobType string) (model *JobData, err error) {
	model = &JobData{
		JobType: core.StringPtr(jobType),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalJobData unmarshals an instance of JobData from the specified map of raw messages.
func UnmarshalJobData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobData)
	err = core.UnmarshalPrimitive(m, "job_type", &obj.JobType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_job_data", &obj.ActionJobData, UnmarshalJobDataAction)
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

	// Job name, uniquely derived from the related Action.
	Name *string `json:"name,omitempty"`

	// Job description derived from the related Action.
	Description *string `json:"description,omitempty"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (action-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name derived from the related Action,.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Inventory of host and host group for the playbook, in .ini file format.
	TargetsIni *string `json:"targets_ini,omitempty"`

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
	JobLite_CommandObject_Action    = "action"
	JobLite_CommandObject_Workspace = "workspace"
)

// Constants associated with the JobLite.CommandName property.
// Schematics job command name.
const (
	JobLite_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	JobLite_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	JobLite_CommandName_HelmInstall          = "helm_install"
	JobLite_CommandName_HelmList             = "helm_list"
	JobLite_CommandName_HelmShow             = "helm_show"
	JobLite_CommandName_OpaEvaluate          = "opa_evaluate"
	JobLite_CommandName_TerraformInit        = "terraform_init"
	JobLite_CommandName_TerrformApply        = "terrform_apply"
	JobLite_CommandName_TerrformDestroy      = "terrform_destroy"
	JobLite_CommandName_TerrformPlan         = "terrform_plan"
	JobLite_CommandName_TerrformRefresh      = "terrform_refresh"
	JobLite_CommandName_TerrformShow         = "terrform_show"
	JobLite_CommandName_TerrformTaint        = "terrform_taint"
	JobLite_CommandName_WorkspaceApplyFlow   = "workspace_apply_flow"
	JobLite_CommandName_WorkspaceCustomFlow  = "workspace_custom_flow"
	JobLite_CommandName_WorkspaceDestroyFlow = "workspace_destroy_flow"
	JobLite_CommandName_WorkspaceInitFlow    = "workspace_init_flow"
	JobLite_CommandName_WorkspacePlanFlow    = "workspace_plan_flow"
	JobLite_CommandName_WorkspaceRefreshFlow = "workspace_refresh_flow"
	JobLite_CommandName_WorkspaceShowFlow    = "workspace_show_flow"
)

// Constants associated with the JobLite.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	JobLite_Location_EuDe    = "eu_de"
	JobLite_Location_EuGb    = "eu_gb"
	JobLite_Location_UsEast  = "us_east"
	JobLite_Location_UsSouth = "us_south"
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
	err = core.UnmarshalPrimitive(m, "targets_ini", &obj.TargetsIni)
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

	// Job name, uniquely derived from the related Action.
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
	LogErrors []JobLogSummaryLogErrorsItem `json:"log_errors,omitempty"`

	// Repo download Job log summary.
	RepoDownloadJob *JobLogSummaryRepoDownloadJob `json:"repo_download_job,omitempty"`

	// Flow Job log summary.
	ActionJob *JobLogSummaryActionJob `json:"action_job,omitempty"`
}

// Constants associated with the JobLogSummary.JobType property.
// Type of Job.
const (
	JobLogSummary_JobType_ActionJob       = "action_job"
	JobLogSummary_JobType_CapsuleJob      = "capsule_job"
	JobLogSummary_JobType_ControlsJob     = "controls_job"
	JobLogSummary_JobType_RepoDownloadJob = "repo_download_job"
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
	err = core.UnmarshalModel(m, "log_errors", &obj.LogErrors, UnmarshalJobLogSummaryLogErrorsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "repo_download_job", &obj.RepoDownloadJob, UnmarshalJobLogSummaryRepoDownloadJob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_job", &obj.ActionJob, UnmarshalJobLogSummaryActionJob)
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

// JobLogSummaryLogErrorsItem : JobLogSummaryLogErrorsItem struct
type JobLogSummaryLogErrorsItem struct {
	// Error code in the Log.
	ErrorCode *string `json:"error_code,omitempty"`

	// Summary error message in the log.
	ErrorMsg *string `json:"error_msg,omitempty"`

	// Number of occurrence.
	ErrorCount *float64 `json:"error_count,omitempty"`
}

// UnmarshalJobLogSummaryLogErrorsItem unmarshals an instance of JobLogSummaryLogErrorsItem from the specified map of raw messages.
func UnmarshalJobLogSummaryLogErrorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobLogSummaryLogErrorsItem)
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

// JobStateData : Workspace Job state-file.
type JobStateData struct {
	// Job Id.
	JobID *string `json:"job_id,omitempty"`

	// Job name, uniquely derived from the related Action.
	JobName *string `json:"job_name,omitempty"`

	// Job state summary.
	Summary []JobStateDataSummaryItem `json:"summary,omitempty"`

	// Format of the State data (eg. tfstate).
	Format *string `json:"format,omitempty"`

	// State data file.
	Details *[]byte `json:"details,omitempty"`

	// Job status updation timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalJobStateData unmarshals an instance of JobStateData from the specified map of raw messages.
func UnmarshalJobStateData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStateData)
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "job_name", &obj.JobName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalJobStateDataSummaryItem)
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

// JobStateDataSummaryItem : JobStateDataSummaryItem struct
type JobStateDataSummaryItem struct {
	// State summary feature name.
	Name *string `json:"name,omitempty"`

	// State summary feature type.
	Type *string `json:"type,omitempty"`

	// State summary feature value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the JobStateDataSummaryItem.Type property.
// State summary feature type.
const (
	JobStateDataSummaryItem_Type_Number = "number"
	JobStateDataSummaryItem_Type_String = "string"
)

// UnmarshalJobStateDataSummaryItem unmarshals an instance of JobStateDataSummaryItem from the specified map of raw messages.
func UnmarshalJobStateDataSummaryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStateDataSummaryItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// JobStatus : Job Status.
type JobStatus struct {
	// Action Job Status.
	ActionJobStatus *JobStatusAction `json:"action_job_status,omitempty"`
}

// UnmarshalJobStatus unmarshals an instance of JobStatus from the specified map of raw messages.
func UnmarshalJobStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatus)
	err = core.UnmarshalModel(m, "action_job_status", &obj.ActionJobStatus, UnmarshalJobStatusAction)
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
	JobStatusAction_StatusCode_IobFinished   = "iob_finished"
	JobStatusAction_StatusCode_JobCancelled  = "job_cancelled"
	JobStatusAction_StatusCode_JobFailed     = "job_failed"
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

// JobStatusType : JobStatusType -.
type JobStatusType struct {
	// List of failed workspace jobs.
	Failed []string `json:"failed,omitempty"`

	// List of in_progress workspace jobs.
	InProgress []string `json:"in_progress,omitempty"`

	// List of successful workspace jobs.
	Success []string `json:"success,omitempty"`

	// Workspace job status updated at.
	LastUpdatedOn *strfmt.DateTime `json:"last_updated_on,omitempty"`
}

// UnmarshalJobStatusType unmarshals an instance of JobStatusType from the specified map of raw messages.
func UnmarshalJobStatusType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobStatusType)
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
	Keys []KMSInstancesKeysItem `json:"keys,omitempty"`
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
	err = core.UnmarshalModel(m, "keys", &obj.Keys, UnmarshalKMSInstancesKeysItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KMSInstancesKeysItem : KMSInstancesKeysItem struct
type KMSInstancesKeysItem struct {
	// Key name.
	Name *string `json:"name,omitempty"`

	// CRN of the Key.
	Crn *string `json:"crn,omitempty"`

	// Error message.
	Error *string `json:"error,omitempty"`
}

// UnmarshalKMSInstancesKeysItem unmarshals an instance of KMSInstancesKeysItem from the specified map of raw messages.
func UnmarshalKMSInstancesKeysItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KMSInstancesKeysItem)
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
	// The number of items to skip before starting to collect the result set.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *int64 `json:"limit,omitempty"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"sort,omitempty"`

	// Level of details returned by the get method.
	Profile *string `json:"profile,omitempty"`

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
func (options *ListActionsOptions) SetOffset(offset int64) *ListActionsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListActionsOptions) SetLimit(limit int64) *ListActionsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListActionsOptions) SetSort(sort string) *ListActionsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetProfile : Allow user to set Profile
func (options *ListActionsOptions) SetProfile(profile string) *ListActionsOptions {
	options.Profile = core.StringPtr(profile)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListActionsOptions) SetHeaders(param map[string]string) *ListActionsOptions {
	options.Headers = param
	return options
}

// ListJobLogsOptions : The ListJobLogs options.
type ListJobLogsOptions struct {
	// Job Id. Use GET /jobs API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"job_id" validate:"required,ne="`

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
func (options *ListJobLogsOptions) SetJobID(jobID string) *ListJobLogsOptions {
	options.JobID = core.StringPtr(jobID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobLogsOptions) SetHeaders(param map[string]string) *ListJobLogsOptions {
	options.Headers = param
	return options
}

// ListJobStatesOptions : The ListJobStates options.
type ListJobStatesOptions struct {
	// Job Id. Use GET /jobs API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListJobStatesOptions : Instantiate ListJobStatesOptions
func (*SchematicsV1) NewListJobStatesOptions(jobID string) *ListJobStatesOptions {
	return &ListJobStatesOptions{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (options *ListJobStatesOptions) SetJobID(jobID string) *ListJobStatesOptions {
	options.JobID = core.StringPtr(jobID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobStatesOptions) SetHeaders(param map[string]string) *ListJobStatesOptions {
	options.Headers = param
	return options
}

// ListJobsOptions : The ListJobs options.
type ListJobsOptions struct {
	// The number of items to skip before starting to collect the result set.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *int64 `json:"limit,omitempty"`

	// Name of the field to sort-by;  Use the '.' character to delineate sub-resources and sub-fields (eg.
	// owner.last_name). Prepend the field with '+' or '-', indicating 'ascending' or 'descending' (default is ascending)
	// Ignore unrecognized or unsupported sort field.
	Sort *string `json:"sort,omitempty"`

	// Level of details returned by the get method.
	Profile *string `json:"profile,omitempty"`

	// Name of the resource (workspace, actions or controls).
	Resource *string `json:"resource,omitempty"`

	// Action Id.
	ActionID *string `json:"action_id,omitempty"`

	// list jobs.
	List *string `json:"list,omitempty"`

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
	ListJobsOptions_Resource_Actions    = "actions"
	ListJobsOptions_Resource_Controls   = "controls"
	ListJobsOptions_Resource_Workspaces = "workspaces"
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
func (options *ListJobsOptions) SetOffset(offset int64) *ListJobsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListJobsOptions) SetLimit(limit int64) *ListJobsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListJobsOptions) SetSort(sort string) *ListJobsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetProfile : Allow user to set Profile
func (options *ListJobsOptions) SetProfile(profile string) *ListJobsOptions {
	options.Profile = core.StringPtr(profile)
	return options
}

// SetResource : Allow user to set Resource
func (options *ListJobsOptions) SetResource(resource string) *ListJobsOptions {
	options.Resource = core.StringPtr(resource)
	return options
}

// SetActionID : Allow user to set ActionID
func (options *ListJobsOptions) SetActionID(actionID string) *ListJobsOptions {
	options.ActionID = core.StringPtr(actionID)
	return options
}

// SetList : Allow user to set List
func (options *ListJobsOptions) SetList(list string) *ListJobsOptions {
	options.List = core.StringPtr(list)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobsOptions) SetHeaders(param map[string]string) *ListJobsOptions {
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

// ListSharedDatasetsOptions : The ListSharedDatasets options.
type ListSharedDatasetsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSharedDatasetsOptions : Instantiate ListSharedDatasetsOptions
func (*SchematicsV1) NewListSharedDatasetsOptions() *ListSharedDatasetsOptions {
	return &ListSharedDatasetsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListSharedDatasetsOptions) SetHeaders(param map[string]string) *ListSharedDatasetsOptions {
	options.Headers = param
	return options
}

// ListWorkspaceActivitiesOptions : The ListWorkspaceActivities options.
type ListWorkspaceActivitiesOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The number of items to skip before starting to collect the result set.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *int64 `json:"limit,omitempty"`

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
func (options *ListWorkspaceActivitiesOptions) SetWID(wID string) *ListWorkspaceActivitiesOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListWorkspaceActivitiesOptions) SetOffset(offset int64) *ListWorkspaceActivitiesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListWorkspaceActivitiesOptions) SetLimit(limit int64) *ListWorkspaceActivitiesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListWorkspaceActivitiesOptions) SetHeaders(param map[string]string) *ListWorkspaceActivitiesOptions {
	options.Headers = param
	return options
}

// ListWorkspacesOptions : The ListWorkspaces options.
type ListWorkspacesOptions struct {
	// The number of items to skip before starting to collect the result set.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListWorkspacesOptions : Instantiate ListWorkspacesOptions
func (*SchematicsV1) NewListWorkspacesOptions() *ListWorkspacesOptions {
	return &ListWorkspacesOptions{}
}

// SetOffset : Allow user to set Offset
func (options *ListWorkspacesOptions) SetOffset(offset int64) *ListWorkspacesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListWorkspacesOptions) SetLimit(limit int64) *ListWorkspacesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListWorkspacesOptions) SetHeaders(param map[string]string) *ListWorkspacesOptions {
	options.Headers = param
	return options
}

// LogStoreResponse : LogStoreResponse -.
type LogStoreResponse struct {
	// Engine name.
	EngineName *string `json:"engine_name,omitempty"`

	// Engine version.
	EngineVersion *string `json:"engine_version,omitempty"`

	// Engine id.
	ID *string `json:"id,omitempty"`

	// Log store url.
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

// LogStoreResponseList : LogStoreResponseList -.
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

// LogSummary : LogSummary ...
type LogSummary struct {
	// WorkspaceActivityStatus activity status type.
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

	// Time takemn to perform activity.
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

// OutputValuesItem : OutputValuesItem struct
type OutputValuesItem struct {
	// Output variable name.
	Folder *string `json:"folder,omitempty"`

	// Output variable id.
	ID *string `json:"id,omitempty"`

	// List of Output values.
	OutputValues []interface{} `json:"output_values,omitempty"`

	// Output variable type.
	ValueType *string `json:"value_type,omitempty"`
}

// UnmarshalOutputValuesItem unmarshals an instance of OutputValuesItem from the specified map of raw messages.
func UnmarshalOutputValuesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OutputValuesItem)
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
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

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
func (options *PlanWorkspaceCommandOptions) SetWID(wID string) *PlanWorkspaceCommandOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *PlanWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *PlanWorkspaceCommandOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PlanWorkspaceCommandOptions) SetHeaders(param map[string]string) *PlanWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// RefreshWorkspaceCommandOptions : The RefreshWorkspaceCommand options.
type RefreshWorkspaceCommandOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

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
func (options *RefreshWorkspaceCommandOptions) SetWID(wID string) *RefreshWorkspaceCommandOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *RefreshWorkspaceCommandOptions) SetRefreshToken(refreshToken string) *RefreshWorkspaceCommandOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *RefreshWorkspaceCommandOptions) SetHeaders(param map[string]string) *RefreshWorkspaceCommandOptions {
	options.Headers = param
	return options
}

// ReplaceJobOptions : The ReplaceJob options.
type ReplaceJobOptions struct {
	// Job Id. Use GET /jobs API to look up the Job Ids in your IBM Cloud account.
	JobID *string `json:"job_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// Name of the Schematics automation resource.
	CommandObject *string `json:"command_object,omitempty"`

	// Job command object id (workspace-id, action-id or control-id).
	CommandObjectID *string `json:"command_object_id,omitempty"`

	// Schematics job command name.
	CommandName *string `json:"command_name,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Command line options for the command.
	CommandOptions []string `json:"command_options,omitempty"`

	// Job inputs used by Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Environment variables used by the Job while performing Action.
	Settings []VariableData `json:"settings,omitempty"`

	// User defined tags, while running the job.
	Tags []string `json:"tags,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Job Status.
	Status *JobStatus `json:"status,omitempty"`

	// Job data.
	Data *JobData `json:"data,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

	// Job log summary record.
	LogSummary *JobLogSummary `json:"log_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceJobOptions.CommandObject property.
// Name of the Schematics automation resource.
const (
	ReplaceJobOptions_CommandObject_Action    = "action"
	ReplaceJobOptions_CommandObject_Workspace = "workspace"
)

// Constants associated with the ReplaceJobOptions.CommandName property.
// Schematics job command name.
const (
	ReplaceJobOptions_CommandName_AnsiblePlaybookCheck = "ansible_playbook_check"
	ReplaceJobOptions_CommandName_AnsiblePlaybookRun   = "ansible_playbook_run"
	ReplaceJobOptions_CommandName_HelmInstall          = "helm_install"
	ReplaceJobOptions_CommandName_HelmList             = "helm_list"
	ReplaceJobOptions_CommandName_HelmShow             = "helm_show"
	ReplaceJobOptions_CommandName_OpaEvaluate          = "opa_evaluate"
	ReplaceJobOptions_CommandName_TerraformInit        = "terraform_init"
	ReplaceJobOptions_CommandName_TerrformApply        = "terrform_apply"
	ReplaceJobOptions_CommandName_TerrformDestroy      = "terrform_destroy"
	ReplaceJobOptions_CommandName_TerrformPlan         = "terrform_plan"
	ReplaceJobOptions_CommandName_TerrformRefresh      = "terrform_refresh"
	ReplaceJobOptions_CommandName_TerrformShow         = "terrform_show"
	ReplaceJobOptions_CommandName_TerrformTaint        = "terrform_taint"
	ReplaceJobOptions_CommandName_WorkspaceApplyFlow   = "workspace_apply_flow"
	ReplaceJobOptions_CommandName_WorkspaceCustomFlow  = "workspace_custom_flow"
	ReplaceJobOptions_CommandName_WorkspaceDestroyFlow = "workspace_destroy_flow"
	ReplaceJobOptions_CommandName_WorkspaceInitFlow    = "workspace_init_flow"
	ReplaceJobOptions_CommandName_WorkspacePlanFlow    = "workspace_plan_flow"
	ReplaceJobOptions_CommandName_WorkspaceRefreshFlow = "workspace_refresh_flow"
	ReplaceJobOptions_CommandName_WorkspaceShowFlow    = "workspace_show_flow"
)

// Constants associated with the ReplaceJobOptions.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	ReplaceJobOptions_Location_EuDe    = "eu_de"
	ReplaceJobOptions_Location_EuGb    = "eu_gb"
	ReplaceJobOptions_Location_UsEast  = "us_east"
	ReplaceJobOptions_Location_UsSouth = "us_south"
)

// NewReplaceJobOptions : Instantiate ReplaceJobOptions
func (*SchematicsV1) NewReplaceJobOptions(jobID string, refreshToken string) *ReplaceJobOptions {
	return &ReplaceJobOptions{
		JobID:        core.StringPtr(jobID),
		RefreshToken: core.StringPtr(refreshToken),
	}
}

// SetJobID : Allow user to set JobID
func (options *ReplaceJobOptions) SetJobID(jobID string) *ReplaceJobOptions {
	options.JobID = core.StringPtr(jobID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *ReplaceJobOptions) SetRefreshToken(refreshToken string) *ReplaceJobOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetCommandObject : Allow user to set CommandObject
func (options *ReplaceJobOptions) SetCommandObject(commandObject string) *ReplaceJobOptions {
	options.CommandObject = core.StringPtr(commandObject)
	return options
}

// SetCommandObjectID : Allow user to set CommandObjectID
func (options *ReplaceJobOptions) SetCommandObjectID(commandObjectID string) *ReplaceJobOptions {
	options.CommandObjectID = core.StringPtr(commandObjectID)
	return options
}

// SetCommandName : Allow user to set CommandName
func (options *ReplaceJobOptions) SetCommandName(commandName string) *ReplaceJobOptions {
	options.CommandName = core.StringPtr(commandName)
	return options
}

// SetCommandParameter : Allow user to set CommandParameter
func (options *ReplaceJobOptions) SetCommandParameter(commandParameter string) *ReplaceJobOptions {
	options.CommandParameter = core.StringPtr(commandParameter)
	return options
}

// SetCommandOptions : Allow user to set CommandOptions
func (options *ReplaceJobOptions) SetCommandOptions(commandOptions []string) *ReplaceJobOptions {
	options.CommandOptions = commandOptions
	return options
}

// SetInputs : Allow user to set Inputs
func (options *ReplaceJobOptions) SetInputs(inputs []VariableData) *ReplaceJobOptions {
	options.Inputs = inputs
	return options
}

// SetSettings : Allow user to set Settings
func (options *ReplaceJobOptions) SetSettings(settings []VariableData) *ReplaceJobOptions {
	options.Settings = settings
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceJobOptions) SetTags(tags []string) *ReplaceJobOptions {
	options.Tags = tags
	return options
}

// SetLocation : Allow user to set Location
func (options *ReplaceJobOptions) SetLocation(location string) *ReplaceJobOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetStatus : Allow user to set Status
func (options *ReplaceJobOptions) SetStatus(status *JobStatus) *ReplaceJobOptions {
	options.Status = status
	return options
}

// SetData : Allow user to set Data
func (options *ReplaceJobOptions) SetData(data *JobData) *ReplaceJobOptions {
	options.Data = data
	return options
}

// SetBastion : Allow user to set Bastion
func (options *ReplaceJobOptions) SetBastion(bastion *TargetResourceset) *ReplaceJobOptions {
	options.Bastion = bastion
	return options
}

// SetLogSummary : Allow user to set LogSummary
func (options *ReplaceJobOptions) SetLogSummary(logSummary *JobLogSummary) *ReplaceJobOptions {
	options.LogSummary = logSummary
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceJobOptions) SetHeaders(param map[string]string) *ReplaceJobOptions {
	options.Headers = param
	return options
}

// ReplaceKmsSettingsOptions : The ReplaceKmsSettings options.
type ReplaceKmsSettingsOptions struct {
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

// NewReplaceKmsSettingsOptions : Instantiate ReplaceKmsSettingsOptions
func (*SchematicsV1) NewReplaceKmsSettingsOptions() *ReplaceKmsSettingsOptions {
	return &ReplaceKmsSettingsOptions{}
}

// SetLocation : Allow user to set Location
func (options *ReplaceKmsSettingsOptions) SetLocation(location string) *ReplaceKmsSettingsOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetEncryptionScheme : Allow user to set EncryptionScheme
func (options *ReplaceKmsSettingsOptions) SetEncryptionScheme(encryptionScheme string) *ReplaceKmsSettingsOptions {
	options.EncryptionScheme = core.StringPtr(encryptionScheme)
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *ReplaceKmsSettingsOptions) SetResourceGroup(resourceGroup string) *ReplaceKmsSettingsOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetPrimaryCrk : Allow user to set PrimaryCrk
func (options *ReplaceKmsSettingsOptions) SetPrimaryCrk(primaryCrk *KMSSettingsPrimaryCrk) *ReplaceKmsSettingsOptions {
	options.PrimaryCrk = primaryCrk
	return options
}

// SetSecondaryCrk : Allow user to set SecondaryCrk
func (options *ReplaceKmsSettingsOptions) SetSecondaryCrk(secondaryCrk *KMSSettingsSecondaryCrk) *ReplaceKmsSettingsOptions {
	options.SecondaryCrk = secondaryCrk
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceKmsSettingsOptions) SetHeaders(param map[string]string) *ReplaceKmsSettingsOptions {
	options.Headers = param
	return options
}

// ReplaceSharedDatasetOptions : The ReplaceSharedDataset options.
type ReplaceSharedDatasetOptions struct {
	// The shared dataset ID Use the GET /shared_datasets to look up the shared dataset IDs  in your IBM Cloud account.
	SdID *string `json:"sd_id" validate:"required,ne="`

	// Automatically propagate changes to consumers.
	AutoPropagateChange *bool `json:"auto_propagate_change,omitempty"`

	// Dataset description.
	Description *string `json:"description,omitempty"`

	// Affected workspaces.
	EffectedWorkspaceIds []string `json:"effected_workspace_ids,omitempty"`

	// Resource group name.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Shared dataset data.
	SharedDatasetData []SharedDatasetData `json:"shared_dataset_data,omitempty"`

	// Shared dataset name.
	SharedDatasetName *string `json:"shared_dataset_name,omitempty"`

	// Shared dataset source name.
	SharedDatasetSourceName *string `json:"shared_dataset_source_name,omitempty"`

	// Shared dataset type.
	SharedDatasetType []string `json:"shared_dataset_type,omitempty"`

	// Shared dataset tags.
	Tags []string `json:"tags,omitempty"`

	// Shared dataset version.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceSharedDatasetOptions : Instantiate ReplaceSharedDatasetOptions
func (*SchematicsV1) NewReplaceSharedDatasetOptions(sdID string) *ReplaceSharedDatasetOptions {
	return &ReplaceSharedDatasetOptions{
		SdID: core.StringPtr(sdID),
	}
}

// SetSdID : Allow user to set SdID
func (options *ReplaceSharedDatasetOptions) SetSdID(sdID string) *ReplaceSharedDatasetOptions {
	options.SdID = core.StringPtr(sdID)
	return options
}

// SetAutoPropagateChange : Allow user to set AutoPropagateChange
func (options *ReplaceSharedDatasetOptions) SetAutoPropagateChange(autoPropagateChange bool) *ReplaceSharedDatasetOptions {
	options.AutoPropagateChange = core.BoolPtr(autoPropagateChange)
	return options
}

// SetDescription : Allow user to set Description
func (options *ReplaceSharedDatasetOptions) SetDescription(description string) *ReplaceSharedDatasetOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEffectedWorkspaceIds : Allow user to set EffectedWorkspaceIds
func (options *ReplaceSharedDatasetOptions) SetEffectedWorkspaceIds(effectedWorkspaceIds []string) *ReplaceSharedDatasetOptions {
	options.EffectedWorkspaceIds = effectedWorkspaceIds
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *ReplaceSharedDatasetOptions) SetResourceGroup(resourceGroup string) *ReplaceSharedDatasetOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetSharedDatasetData : Allow user to set SharedDatasetData
func (options *ReplaceSharedDatasetOptions) SetSharedDatasetData(sharedDatasetData []SharedDatasetData) *ReplaceSharedDatasetOptions {
	options.SharedDatasetData = sharedDatasetData
	return options
}

// SetSharedDatasetName : Allow user to set SharedDatasetName
func (options *ReplaceSharedDatasetOptions) SetSharedDatasetName(sharedDatasetName string) *ReplaceSharedDatasetOptions {
	options.SharedDatasetName = core.StringPtr(sharedDatasetName)
	return options
}

// SetSharedDatasetSourceName : Allow user to set SharedDatasetSourceName
func (options *ReplaceSharedDatasetOptions) SetSharedDatasetSourceName(sharedDatasetSourceName string) *ReplaceSharedDatasetOptions {
	options.SharedDatasetSourceName = core.StringPtr(sharedDatasetSourceName)
	return options
}

// SetSharedDatasetType : Allow user to set SharedDatasetType
func (options *ReplaceSharedDatasetOptions) SetSharedDatasetType(sharedDatasetType []string) *ReplaceSharedDatasetOptions {
	options.SharedDatasetType = sharedDatasetType
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceSharedDatasetOptions) SetTags(tags []string) *ReplaceSharedDatasetOptions {
	options.Tags = tags
	return options
}

// SetVersion : Allow user to set Version
func (options *ReplaceSharedDatasetOptions) SetVersion(version string) *ReplaceSharedDatasetOptions {
	options.Version = core.StringPtr(version)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceSharedDatasetOptions) SetHeaders(param map[string]string) *ReplaceSharedDatasetOptions {
	options.Headers = param
	return options
}

// ReplaceWorkspaceInputsOptions : The ReplaceWorkspaceInputs options.
type ReplaceWorkspaceInputsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

	// EnvVariableRequest ..
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
func (options *ReplaceWorkspaceInputsOptions) SetWID(wID string) *ReplaceWorkspaceInputsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *ReplaceWorkspaceInputsOptions) SetTID(tID string) *ReplaceWorkspaceInputsOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetEnvValues : Allow user to set EnvValues
func (options *ReplaceWorkspaceInputsOptions) SetEnvValues(envValues []interface{}) *ReplaceWorkspaceInputsOptions {
	options.EnvValues = envValues
	return options
}

// SetValues : Allow user to set Values
func (options *ReplaceWorkspaceInputsOptions) SetValues(values string) *ReplaceWorkspaceInputsOptions {
	options.Values = core.StringPtr(values)
	return options
}

// SetVariablestore : Allow user to set Variablestore
func (options *ReplaceWorkspaceInputsOptions) SetVariablestore(variablestore []WorkspaceVariableRequest) *ReplaceWorkspaceInputsOptions {
	options.Variablestore = variablestore
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceWorkspaceInputsOptions) SetHeaders(param map[string]string) *ReplaceWorkspaceInputsOptions {
	options.Headers = param
	return options
}

// ReplaceWorkspaceOptions : The ReplaceWorkspace options.
type ReplaceWorkspaceOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// CatalogRef -.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// Workspace description.
	Description *string `json:"description,omitempty"`

	// Workspace name.
	Name *string `json:"name,omitempty"`

	// SharedTargetData -.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// Tags -.
	Tags []string `json:"tags,omitempty"`

	// TemplateData -.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// TemplateRepoUpdateRequest -.
	TemplateRepo *TemplateRepoUpdateRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// WorkspaceStatusUpdateRequest -.
	WorkspaceStatus *WorkspaceStatusUpdateRequest `json:"workspace_status,omitempty"`

	// WorkspaceStatusMessage -.
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
func (options *ReplaceWorkspaceOptions) SetWID(wID string) *ReplaceWorkspaceOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetCatalogRef : Allow user to set CatalogRef
func (options *ReplaceWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *ReplaceWorkspaceOptions {
	options.CatalogRef = catalogRef
	return options
}

// SetDescription : Allow user to set Description
func (options *ReplaceWorkspaceOptions) SetDescription(description string) *ReplaceWorkspaceOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetName : Allow user to set Name
func (options *ReplaceWorkspaceOptions) SetName(name string) *ReplaceWorkspaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetSharedData : Allow user to set SharedData
func (options *ReplaceWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *ReplaceWorkspaceOptions {
	options.SharedData = sharedData
	return options
}

// SetTags : Allow user to set Tags
func (options *ReplaceWorkspaceOptions) SetTags(tags []string) *ReplaceWorkspaceOptions {
	options.Tags = tags
	return options
}

// SetTemplateData : Allow user to set TemplateData
func (options *ReplaceWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *ReplaceWorkspaceOptions {
	options.TemplateData = templateData
	return options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (options *ReplaceWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoUpdateRequest) *ReplaceWorkspaceOptions {
	options.TemplateRepo = templateRepo
	return options
}

// SetType : Allow user to set Type
func (options *ReplaceWorkspaceOptions) SetType(typeVar []string) *ReplaceWorkspaceOptions {
	options.Type = typeVar
	return options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (options *ReplaceWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusUpdateRequest) *ReplaceWorkspaceOptions {
	options.WorkspaceStatus = workspaceStatus
	return options
}

// SetWorkspaceStatusMsg : Allow user to set WorkspaceStatusMsg
func (options *ReplaceWorkspaceOptions) SetWorkspaceStatusMsg(workspaceStatusMsg *WorkspaceStatusMessage) *ReplaceWorkspaceOptions {
	options.WorkspaceStatusMsg = workspaceStatusMsg
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceWorkspaceOptions) SetHeaders(param map[string]string) *ReplaceWorkspaceOptions {
	options.Headers = param
	return options
}

// ResourceGroupResponse : ResourceGroupResponse -.
type ResourceGroupResponse struct {
	// Account id.
	AccountID *string `json:"account_id,omitempty"`

	// CRN.
	Crn *string `json:"crn,omitempty"`

	// default.
	Default *bool `json:"default,omitempty"`

	// Resource group name.
	Name *string `json:"name,omitempty"`

	// Resource group id.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Resource group state.
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

// RunWorkspaceCommandsOptions : The RunWorkspaceCommands options.
type RunWorkspaceCommandsOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The IAM refresh token associated with the IBM Cloud account.
	RefreshToken *string `json:"refresh_token" validate:"required"`

	// List of commands.
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
func (options *RunWorkspaceCommandsOptions) SetWID(wID string) *RunWorkspaceCommandsOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetRefreshToken : Allow user to set RefreshToken
func (options *RunWorkspaceCommandsOptions) SetRefreshToken(refreshToken string) *RunWorkspaceCommandsOptions {
	options.RefreshToken = core.StringPtr(refreshToken)
	return options
}

// SetCommands : Allow user to set Commands
func (options *RunWorkspaceCommandsOptions) SetCommands(commands []TerraformCommand) *RunWorkspaceCommandsOptions {
	options.Commands = commands
	return options
}

// SetOperationName : Allow user to set OperationName
func (options *RunWorkspaceCommandsOptions) SetOperationName(operationName string) *RunWorkspaceCommandsOptions {
	options.OperationName = core.StringPtr(operationName)
	return options
}

// SetDescription : Allow user to set Description
func (options *RunWorkspaceCommandsOptions) SetDescription(description string) *RunWorkspaceCommandsOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *RunWorkspaceCommandsOptions) SetHeaders(param map[string]string) *RunWorkspaceCommandsOptions {
	options.Headers = param
	return options
}

// SchematicsLocations : Schematics locations.
type SchematicsLocations struct {
	// Country.
	Country *string `json:"country,omitempty"`

	// Geography.
	Geography *string `json:"geography,omitempty"`

	// Location id.
	ID *string `json:"id,omitempty"`

	// Kind.
	Kind *string `json:"kind,omitempty"`

	// Metro.
	Metro *string `json:"metro,omitempty"`

	// Multizone metro.
	MultizoneMetro *string `json:"multizone_metro,omitempty"`

	// Location name.
	Name *string `json:"name,omitempty"`
}

// UnmarshalSchematicsLocations unmarshals an instance of SchematicsLocations from the specified map of raw messages.
func UnmarshalSchematicsLocations(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsLocations)
	err = core.UnmarshalPrimitive(m, "country", &obj.Country)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "geography", &obj.Geography)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedDatasetData : SharedDatasetData ...
type SharedDatasetData struct {
	// Default values.
	DefaultValue *string `json:"default_value,omitempty"`

	// Data description.
	Description *string `json:"description,omitempty"`

	// Data is hidden.
	Hidden *bool `json:"hidden,omitempty"`

	// Data is readonly.
	Immutable *bool `json:"immutable,omitempty"`

	// Data is matches regular expression.
	Matches *string `json:"matches,omitempty"`

	// Max value of the data.
	MaxValue *string `json:"max_value,omitempty"`

	// Max string length of the data.
	MaxValueLen *string `json:"max_value_len,omitempty"`

	// Min value of the data.
	MinValue *string `json:"min_value,omitempty"`

	// Min string length of the data.
	MinValueLen *string `json:"min_value_len,omitempty"`

	// Possible options for the Data.
	Options []string `json:"options,omitempty"`

	// Override value for the Data.
	OverrideValue *string `json:"override_value,omitempty"`

	// Data is secure.
	Secure *bool `json:"secure,omitempty"`

	// Alias strings for the variable names.
	VarAliases []string `json:"var_aliases,omitempty"`

	// Variable name.
	VarName *string `json:"var_name,omitempty"`

	// Variable reference.
	VarRef *string `json:"var_ref,omitempty"`

	// Variable type.
	VarType *string `json:"var_type,omitempty"`
}

// UnmarshalSharedDatasetData unmarshals an instance of SharedDatasetData from the specified map of raw messages.
func UnmarshalSharedDatasetData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SharedDatasetData)
	err = core.UnmarshalPrimitive(m, "default_value", &obj.DefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "immutable", &obj.Immutable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "matches", &obj.Matches)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_value", &obj.MaxValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "max_value_len", &obj.MaxValueLen)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "min_value", &obj.MinValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "min_value_len", &obj.MinValueLen)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "options", &obj.Options)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "override_value", &obj.OverrideValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secure", &obj.Secure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "var_aliases", &obj.VarAliases)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "var_name", &obj.VarName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "var_ref", &obj.VarRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "var_type", &obj.VarType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedDatasetResponse : SharedDatasetResponse - request returned by create.
type SharedDatasetResponse struct {
	// Account id.
	Account *string `json:"account,omitempty"`

	// Dataset created at.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Dataset created by.
	CreatedBy *string `json:"created_by,omitempty"`

	// Dataset description.
	Description *string `json:"description,omitempty"`

	// Affected workspace id.
	EffectedWorkspaceIds []string `json:"effected_workspace_ids,omitempty"`

	// Resource group name.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Shared dataset data.
	SharedDatasetData []SharedDatasetData `json:"shared_dataset_data,omitempty"`

	// Shared dataset id.
	SharedDatasetID *string `json:"shared_dataset_id,omitempty"`

	// Shared dataset name.
	SharedDatasetName *string `json:"shared_dataset_name,omitempty"`

	// Shared dataset type.
	SharedDatasetType []string `json:"shared_dataset_type,omitempty"`

	// shareddata variable status type.
	State *string `json:"state,omitempty"`

	// Shared dataset tags.
	Tags []string `json:"tags,omitempty"`

	// Shared dataset updated at.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Shared dataset updated by.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// Shared dataset version.
	Version *string `json:"version,omitempty"`
}

// UnmarshalSharedDatasetResponse unmarshals an instance of SharedDatasetResponse from the specified map of raw messages.
func UnmarshalSharedDatasetResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SharedDatasetResponse)
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "effected_workspace_ids", &obj.EffectedWorkspaceIds)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "shared_dataset_data", &obj.SharedDatasetData, UnmarshalSharedDatasetData)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "shared_dataset_id", &obj.SharedDatasetID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "shared_dataset_name", &obj.SharedDatasetName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "shared_dataset_type", &obj.SharedDatasetType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
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
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedDatasetResponseList : SharedDatasetResponseList -.
type SharedDatasetResponseList struct {
	// Shared dataset count.
	Count *int64 `json:"count,omitempty"`

	// List of datasets.
	SharedDatasets []SharedDatasetResponse `json:"shared_datasets,omitempty"`
}

// UnmarshalSharedDatasetResponseList unmarshals an instance of SharedDatasetResponseList from the specified map of raw messages.
func UnmarshalSharedDatasetResponseList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SharedDatasetResponseList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "shared_datasets", &obj.SharedDatasets, UnmarshalSharedDatasetResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SharedTargetData : SharedTargetData -.
type SharedTargetData struct {
	// Cluster created on.
	ClusterCreatedOn *string `json:"cluster_created_on,omitempty"`

	// Cluster id.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Cluster name.
	ClusterName *string `json:"cluster_name,omitempty"`

	// Cluster type.
	ClusterType *string `json:"cluster_type,omitempty"`

	// Entitlement keys.
	EntitlementKeys []interface{} `json:"entitlement_keys,omitempty"`

	// Target namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Target region.
	Region *string `json:"region,omitempty"`

	// Target resource group id.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Cluster worker count.
	WorkerCount *int64 `json:"worker_count,omitempty"`

	// Cluster worker type.
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

// SharedTargetDataResponse : SharedTargetDataResponse -.
type SharedTargetDataResponse struct {
	// Target cluster id.
	ClusterID *string `json:"cluster_id,omitempty"`

	// Target cluster name.
	ClusterName *string `json:"cluster_name,omitempty"`

	// Entitlement keys.
	EntitlementKeys []interface{} `json:"entitlement_keys,omitempty"`

	// Target namespace.
	Namespace *string `json:"namespace,omitempty"`

	// Target region.
	Region *string `json:"region,omitempty"`

	// Target resource group id.
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

// StateStoreResponse : StateStoreResponse -.
type StateStoreResponse struct {
	// Engine name.
	EngineName *string `json:"engine_name,omitempty"`

	// Engine version.
	EngineVersion *string `json:"engine_version,omitempty"`

	// State store id.
	ID *string `json:"id,omitempty"`

	// State store url.
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

// StateStoreResponseList : StateStoreResponseList -.
type StateStoreResponseList struct {
	// List of state stores.
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
	// Is the Workspace locked by a Schematic action ?.
	SysLocked *bool `json:"sys_locked,omitempty"`

	// Name of the User who performed the action, that lead to the locking of the Workspace.
	SysLockedBy *string `json:"sys_locked_by,omitempty"`

	// When the User performed the action that lead to locking of the Workspace ?.
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

// TargetResourceset : Complete Target details with user inputs and system generated data.
type TargetResourceset struct {
	// Target name.
	Name *string `json:"name,omitempty"`

	// Target type (cluster, vsi, icd, vpc).
	Type *string `json:"type,omitempty"`

	// Target description.
	Description *string `json:"description,omitempty"`

	// Resource selection query string.
	ResourceQuery *string `json:"resource_query,omitempty"`

	// Override credential for each resource.  Reference to credentials values, used by all resources.
	CredentialRef *string `json:"credential_ref,omitempty"`

	// Target id.
	ID *string `json:"id,omitempty"`

	// Targets creation time.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Email address of user who created the Targets.
	CreatedBy *string `json:"created_by,omitempty"`

	// Targets updation time.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Email address of user who updated the Targets.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// Array of resource ids.
	ResourceIds []string `json:"resource_ids,omitempty"`
}

// UnmarshalTargetResourceset unmarshals an instance of TargetResourceset from the specified map of raw messages.
func UnmarshalTargetResourceset(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetResourceset)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_query", &obj.ResourceQuery)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "credential_ref", &obj.CredentialRef)
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
	err = core.UnmarshalModel(m, "sys_lock", &obj.SysLock, UnmarshalSystemLock)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_ids", &obj.ResourceIds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateReadme : TemplateReadme -.
type TemplateReadme struct {
	// Readme string.
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

// TemplateRepoRequest : TemplateRepoRequest -.
type TemplateRepoRequest struct {
	// Repo branch.
	Branch *string `json:"branch,omitempty"`

	// Repo release.
	Release *string `json:"release,omitempty"`

	// Repo SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// Repo URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// Source URL.
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

// TemplateRepoResponse : TemplateRepoResponse -.
type TemplateRepoResponse struct {
	// Repo branch.
	Branch *string `json:"branch,omitempty"`

	// Full repo URL.
	FullURL *string `json:"full_url,omitempty"`

	// Has uploaded git repo tar.
	HasUploadedgitrepotar *bool `json:"has_uploadedgitrepotar,omitempty"`

	// Repo release.
	Release *string `json:"release,omitempty"`

	// Repo SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// Repo URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// Source URL.
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

	// Has received tar file.
	HasReceivedFile *bool `json:"has_received_file,omitempty"`

	// Template id.
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

// TemplateRepoUpdateRequest : TemplateRepoUpdateRequest -.
type TemplateRepoUpdateRequest struct {
	// Repo branch.
	Branch *string `json:"branch,omitempty"`

	// Repo release.
	Release *string `json:"release,omitempty"`

	// Repo SHA value.
	RepoShaValue *string `json:"repo_sha_value,omitempty"`

	// Repo URL.
	RepoURL *string `json:"repo_url,omitempty"`

	// Source URL.
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

// TemplateResources : TemplateResources -.
type TemplateResources struct {
	// Template folder name.
	Folder *string `json:"folder,omitempty"`

	// Template id.
	ID *string `json:"id,omitempty"`

	// List of null resources.
	NullResources []interface{} `json:"null_resources,omitempty"`

	// List of related resources.
	RelatedResources []interface{} `json:"related_resources,omitempty"`

	// List of resources.
	Resources []interface{} `json:"resources,omitempty"`

	// Number of resources.
	ResourcesCount *int64 `json:"resources_count,omitempty"`

	// Type of templaes.
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

// TemplateRunTimeDataResponse : TemplateRunTimeDataResponse -.
type TemplateRunTimeDataResponse struct {
	// Engine command.
	EngineCmd *string `json:"engine_cmd,omitempty"`

	// Engine name.
	EngineName *string `json:"engine_name,omitempty"`

	// Engine version.
	EngineVersion *string `json:"engine_version,omitempty"`

	// Template id.
	ID *string `json:"id,omitempty"`

	// Log store url.
	LogStoreURL *string `json:"log_store_url,omitempty"`

	// List of Output values.
	OutputValues []interface{} `json:"output_values,omitempty"`

	// List of resources.
	Resources [][]interface{} `json:"resources,omitempty"`

	// State store URL.
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

// TemplateSourceDataRequest : TemplateSourceDataRequest -.
type TemplateSourceDataRequest struct {
	// EnvVariableRequest ..
	EnvValues []interface{} `json:"env_values,omitempty"`

	// Folder name.
	Folder *string `json:"folder,omitempty"`

	// Init state file.
	InitStateFile *string `json:"init_state_file,omitempty"`

	// Template type.
	Type *string `json:"type,omitempty"`

	// Uninstall script name.
	UninstallScriptName *string `json:"uninstall_script_name,omitempty"`

	// Value.
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

// TemplateSourceDataResponse : TemplateSourceDataResponse -.
type TemplateSourceDataResponse struct {
	// List of environment values.
	EnvValues []EnvVariableResponse `json:"env_values,omitempty"`

	// Folder name.
	Folder *string `json:"folder,omitempty"`

	// Has github token.
	HasGithubtoken *bool `json:"has_githubtoken,omitempty"`

	// Template id.
	ID *string `json:"id,omitempty"`

	// Template tyoe.
	Type *string `json:"type,omitempty"`

	// Uninstall script name.
	UninstallScriptName *string `json:"uninstall_script_name,omitempty"`

	// Values.
	Values *string `json:"values,omitempty"`

	// List of values metadata.
	ValuesMetadata []interface{} `json:"values_metadata,omitempty"`

	// Values URL.
	ValuesURL *string `json:"values_url,omitempty"`

	// VariablesResponse -.
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

// TemplateStateStore : TemplateStateStore -.
type TemplateStateStore struct {
	Version *float64 `json:"version,omitempty"`

	TerraformVersion *string `json:"terraform_version,omitempty"`

	Serial *float64 `json:"serial,omitempty"`

	Lineage *string `json:"lineage,omitempty"`

	Modules []interface{} `json:"modules,omitempty"`
}

// UnmarshalTemplateStateStore unmarshals an instance of TemplateStateStore from the specified map of raw messages.
func UnmarshalTemplateStateStore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateStateStore)
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "terraform_version", &obj.TerraformVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "serial", &obj.Serial)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "lineage", &obj.Lineage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modules", &obj.Modules)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateValues : TemplateValues -.
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

// TerraformCommand : TerraformCommand -.
type TerraformCommand struct {
	// Command to execute.
	Command *string `json:"command,omitempty"`

	// Command Parameters.
	CommandParams *string `json:"command_params,omitempty"`

	// Command name.
	CommandName *string `json:"command_name,omitempty"`

	// Command description.
	CommandDesc *string `json:"command_desc,omitempty"`

	// Instruction to continue or break in case of error.
	CommandOnError *string `json:"command_onError,omitempty"`

	// Dependency on previous commands.
	CommandDependsOn *string `json:"command_dependsOn,omitempty"`

	// Command status.
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
	err = core.UnmarshalPrimitive(m, "command_onError", &obj.CommandOnError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "command_dependsOn", &obj.CommandDependsOn)
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
	ActionID *string `json:"action_id" validate:"required,ne="`

	// Action name (unique for an account).
	Name *string `json:"name,omitempty"`

	// Action description.
	Description *string `json:"description,omitempty"`

	// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of
	// the resources provisioned using Schematics.
	Location *string `json:"location,omitempty"`

	// Resource-group name for the Action.  By default, Action will be created in Default Resource Group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Action tags.
	Tags []string `json:"tags,omitempty"`

	// User defined status of the Schematics object.
	UserState *UserState `json:"user_state,omitempty"`

	// URL of the README file, for the source.
	SourceReadmeURL *string `json:"source_readme_url,omitempty"`

	// Source of templates, playbooks, or controls.
	Source *ExternalSource `json:"source,omitempty"`

	// Type of source for the Template.
	SourceType *string `json:"source_type,omitempty"`

	// Schematics job command parameter (playbook-name, capsule-name or flow-name).
	CommandParameter *string `json:"command_parameter,omitempty"`

	// Complete Target details with user inputs and system generated data.
	Bastion *TargetResourceset `json:"bastion,omitempty"`

	// Inventory of host and host group for the playbook, in .ini file format.
	TargetsIni *string `json:"targets_ini,omitempty"`

	// credentials of the Action.
	Credentials []VariableData `json:"credentials,omitempty"`

	// Input variables for the Action.
	Inputs []VariableData `json:"inputs,omitempty"`

	// Output variables for the Action.
	Outputs []VariableData `json:"outputs,omitempty"`

	// Environment variables for the Action.
	Settings []VariableData `json:"settings,omitempty"`

	// Id to the Trigger.
	TriggerRecordID *string `json:"trigger_record_id,omitempty"`

	// Computed state of the Action.
	State *ActionState `json:"state,omitempty"`

	// System lock status.
	SysLock *SystemLock `json:"sys_lock,omitempty"`

	// The github token associated with the GIT. Required for cloning of repo.
	XGithubToken *string `json:"X-Github-token,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateActionOptions.Location property.
// List of workspace locations supported by IBM Cloud Schematics service.  Note, this does not limit the location of the
// resources provisioned using Schematics.
const (
	UpdateActionOptions_Location_EuDe    = "eu_de"
	UpdateActionOptions_Location_EuGb    = "eu_gb"
	UpdateActionOptions_Location_UsEast  = "us_east"
	UpdateActionOptions_Location_UsSouth = "us_south"
)

// Constants associated with the UpdateActionOptions.SourceType property.
// Type of source for the Template.
const (
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
func (options *UpdateActionOptions) SetActionID(actionID string) *UpdateActionOptions {
	options.ActionID = core.StringPtr(actionID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateActionOptions) SetName(name string) *UpdateActionOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateActionOptions) SetDescription(description string) *UpdateActionOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetLocation : Allow user to set Location
func (options *UpdateActionOptions) SetLocation(location string) *UpdateActionOptions {
	options.Location = core.StringPtr(location)
	return options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (options *UpdateActionOptions) SetResourceGroup(resourceGroup string) *UpdateActionOptions {
	options.ResourceGroup = core.StringPtr(resourceGroup)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateActionOptions) SetTags(tags []string) *UpdateActionOptions {
	options.Tags = tags
	return options
}

// SetUserState : Allow user to set UserState
func (options *UpdateActionOptions) SetUserState(userState *UserState) *UpdateActionOptions {
	options.UserState = userState
	return options
}

// SetSourceReadmeURL : Allow user to set SourceReadmeURL
func (options *UpdateActionOptions) SetSourceReadmeURL(sourceReadmeURL string) *UpdateActionOptions {
	options.SourceReadmeURL = core.StringPtr(sourceReadmeURL)
	return options
}

// SetSource : Allow user to set Source
func (options *UpdateActionOptions) SetSource(source *ExternalSource) *UpdateActionOptions {
	options.Source = source
	return options
}

// SetSourceType : Allow user to set SourceType
func (options *UpdateActionOptions) SetSourceType(sourceType string) *UpdateActionOptions {
	options.SourceType = core.StringPtr(sourceType)
	return options
}

// SetCommandParameter : Allow user to set CommandParameter
func (options *UpdateActionOptions) SetCommandParameter(commandParameter string) *UpdateActionOptions {
	options.CommandParameter = core.StringPtr(commandParameter)
	return options
}

// SetBastion : Allow user to set Bastion
func (options *UpdateActionOptions) SetBastion(bastion *TargetResourceset) *UpdateActionOptions {
	options.Bastion = bastion
	return options
}

// SetTargetsIni : Allow user to set TargetsIni
func (options *UpdateActionOptions) SetTargetsIni(targetsIni string) *UpdateActionOptions {
	options.TargetsIni = core.StringPtr(targetsIni)
	return options
}

// SetCredentials : Allow user to set Credentials
func (options *UpdateActionOptions) SetCredentials(credentials []VariableData) *UpdateActionOptions {
	options.Credentials = credentials
	return options
}

// SetInputs : Allow user to set Inputs
func (options *UpdateActionOptions) SetInputs(inputs []VariableData) *UpdateActionOptions {
	options.Inputs = inputs
	return options
}

// SetOutputs : Allow user to set Outputs
func (options *UpdateActionOptions) SetOutputs(outputs []VariableData) *UpdateActionOptions {
	options.Outputs = outputs
	return options
}

// SetSettings : Allow user to set Settings
func (options *UpdateActionOptions) SetSettings(settings []VariableData) *UpdateActionOptions {
	options.Settings = settings
	return options
}

// SetTriggerRecordID : Allow user to set TriggerRecordID
func (options *UpdateActionOptions) SetTriggerRecordID(triggerRecordID string) *UpdateActionOptions {
	options.TriggerRecordID = core.StringPtr(triggerRecordID)
	return options
}

// SetState : Allow user to set State
func (options *UpdateActionOptions) SetState(state *ActionState) *UpdateActionOptions {
	options.State = state
	return options
}

// SetSysLock : Allow user to set SysLock
func (options *UpdateActionOptions) SetSysLock(sysLock *SystemLock) *UpdateActionOptions {
	options.SysLock = sysLock
	return options
}

// SetXGithubToken : Allow user to set XGithubToken
func (options *UpdateActionOptions) SetXGithubToken(xGithubToken string) *UpdateActionOptions {
	options.XGithubToken = core.StringPtr(xGithubToken)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateActionOptions) SetHeaders(param map[string]string) *UpdateActionOptions {
	options.Headers = param
	return options
}

// UpdateWorkspaceOptions : The UpdateWorkspace options.
type UpdateWorkspaceOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// CatalogRef -.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// Workspace description.
	Description *string `json:"description,omitempty"`

	// Workspace name.
	Name *string `json:"name,omitempty"`

	// SharedTargetData -.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// Tags -.
	Tags []string `json:"tags,omitempty"`

	// TemplateData -.
	TemplateData []TemplateSourceDataRequest `json:"template_data,omitempty"`

	// TemplateRepoUpdateRequest -.
	TemplateRepo *TemplateRepoUpdateRequest `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// WorkspaceStatusUpdateRequest -.
	WorkspaceStatus *WorkspaceStatusUpdateRequest `json:"workspace_status,omitempty"`

	// WorkspaceStatusMessage -.
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
func (options *UpdateWorkspaceOptions) SetWID(wID string) *UpdateWorkspaceOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetCatalogRef : Allow user to set CatalogRef
func (options *UpdateWorkspaceOptions) SetCatalogRef(catalogRef *CatalogRef) *UpdateWorkspaceOptions {
	options.CatalogRef = catalogRef
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateWorkspaceOptions) SetDescription(description string) *UpdateWorkspaceOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateWorkspaceOptions) SetName(name string) *UpdateWorkspaceOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetSharedData : Allow user to set SharedData
func (options *UpdateWorkspaceOptions) SetSharedData(sharedData *SharedTargetData) *UpdateWorkspaceOptions {
	options.SharedData = sharedData
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateWorkspaceOptions) SetTags(tags []string) *UpdateWorkspaceOptions {
	options.Tags = tags
	return options
}

// SetTemplateData : Allow user to set TemplateData
func (options *UpdateWorkspaceOptions) SetTemplateData(templateData []TemplateSourceDataRequest) *UpdateWorkspaceOptions {
	options.TemplateData = templateData
	return options
}

// SetTemplateRepo : Allow user to set TemplateRepo
func (options *UpdateWorkspaceOptions) SetTemplateRepo(templateRepo *TemplateRepoUpdateRequest) *UpdateWorkspaceOptions {
	options.TemplateRepo = templateRepo
	return options
}

// SetType : Allow user to set Type
func (options *UpdateWorkspaceOptions) SetType(typeVar []string) *UpdateWorkspaceOptions {
	options.Type = typeVar
	return options
}

// SetWorkspaceStatus : Allow user to set WorkspaceStatus
func (options *UpdateWorkspaceOptions) SetWorkspaceStatus(workspaceStatus *WorkspaceStatusUpdateRequest) *UpdateWorkspaceOptions {
	options.WorkspaceStatus = workspaceStatus
	return options
}

// SetWorkspaceStatusMsg : Allow user to set WorkspaceStatusMsg
func (options *UpdateWorkspaceOptions) SetWorkspaceStatusMsg(workspaceStatusMsg *WorkspaceStatusMessage) *UpdateWorkspaceOptions {
	options.WorkspaceStatusMsg = workspaceStatusMsg
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateWorkspaceOptions) SetHeaders(param map[string]string) *UpdateWorkspaceOptions {
	options.Headers = param
	return options
}

// UploadTemplateTarOptions : The UploadTemplateTar options.
type UploadTemplateTarOptions struct {
	// The workspace ID for the workspace that you want to query.  You can run the GET /workspaces call if you need to look
	// up the  workspace IDs in your IBM Cloud account.
	WID *string `json:"w_id" validate:"required,ne="`

	// The Template ID for which you want to get the values.  Use the GET /workspaces to look up the workspace IDs  or
	// template IDs in your IBM Cloud account.
	TID *string `json:"t_id" validate:"required,ne="`

	// Template tar file.
	File io.ReadCloser `json:"file,omitempty"`

	// The content type of file.
	FileContentType *string `json:"file_content_type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUploadTemplateTarOptions : Instantiate UploadTemplateTarOptions
func (*SchematicsV1) NewUploadTemplateTarOptions(wID string, tID string) *UploadTemplateTarOptions {
	return &UploadTemplateTarOptions{
		WID: core.StringPtr(wID),
		TID: core.StringPtr(tID),
	}
}

// SetWID : Allow user to set WID
func (options *UploadTemplateTarOptions) SetWID(wID string) *UploadTemplateTarOptions {
	options.WID = core.StringPtr(wID)
	return options
}

// SetTID : Allow user to set TID
func (options *UploadTemplateTarOptions) SetTID(tID string) *UploadTemplateTarOptions {
	options.TID = core.StringPtr(tID)
	return options
}

// SetFile : Allow user to set File
func (options *UploadTemplateTarOptions) SetFile(file io.ReadCloser) *UploadTemplateTarOptions {
	options.File = file
	return options
}

// SetFileContentType : Allow user to set FileContentType
func (options *UploadTemplateTarOptions) SetFileContentType(fileContentType string) *UploadTemplateTarOptions {
	options.FileContentType = core.StringPtr(fileContentType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UploadTemplateTarOptions) SetHeaders(param map[string]string) *UploadTemplateTarOptions {
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
	// EnvVariableRequest ..
	EnvValues []interface{} `json:"env_values,omitempty"`

	// User values.
	Values *string `json:"values,omitempty"`

	// VariablesResponse -.
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

// VersionResponse : VersionResponse -.
type VersionResponse struct {
	// Build data.
	Builddate *string `json:"builddate,omitempty"`

	// Build number.
	Buildno *string `json:"buildno,omitempty"`

	// Commit SHA.
	Commitsha *string `json:"commitsha,omitempty"`

	// Version number of 'Helm provider for Terraform'.
	HelmProviderVersion *string `json:"helm_provider_version,omitempty"`

	// Helm Version.
	HelmVersion *string `json:"helm_version,omitempty"`

	// Supported template types.
	SupportedTemplateTypes interface{} `json:"supported_template_types,omitempty"`

	// Terraform provider versions.
	TerraformProviderVersion *string `json:"terraform_provider_version,omitempty"`

	// Terraform versions.
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

// WorkspaceActivities : WorkspaceActivities -.
type WorkspaceActivities struct {
	// List of workspace activities.
	Actions []WorkspaceActivity `json:"actions,omitempty"`

	// Workspace id.
	WorkspaceID *string `json:"workspace_id,omitempty"`

	// Workspace name.
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

// WorkspaceActivity : WorkspaceActivity -.
type WorkspaceActivity struct {
	// Activity id.
	ActionID *string `json:"action_id,omitempty"`

	// StatusMessages -.
	Message []string `json:"message,omitempty"`

	// WorkspaceActivityAction activity action type.
	Name *string `json:"name,omitempty"`

	// Activity performed at.
	PerformedAt *strfmt.DateTime `json:"performed_at,omitempty"`

	// Activity performed by.
	PerformedBy *string `json:"performed_by,omitempty"`

	// WorkspaceActivityStatus activity status type.
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

// WorkspaceActivityApplyResult : WorkspaceActivityApplyResult -.
type WorkspaceActivityApplyResult struct {
	// Activity id.
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

// WorkspaceActivityCommandResult : WorkspaceActivityCommandResult -.
type WorkspaceActivityCommandResult struct {
	// Activity id.
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

// WorkspaceActivityDestroyResult : WorkspaceActivityDestroyResult -.
type WorkspaceActivityDestroyResult struct {
	// Activity id.
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

// WorkspaceActivityLogs : WorkspaceActivityLogs -.
type WorkspaceActivityLogs struct {
	// Activity id.
	ActionID *string `json:"action_id,omitempty"`

	// WorkspaceActivityAction activity action type.
	Name *string `json:"name,omitempty"`

	// List of activity logs.
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

// WorkspaceActivityOptionsTemplate : Action Options Template ...
type WorkspaceActivityOptionsTemplate struct {
	// Action targets.
	Target []string `json:"target,omitempty"`

	// Action tfvars.
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

// WorkspaceActivityPlanResult : WorkspaceActivityPlanResult -.
type WorkspaceActivityPlanResult struct {
	// Activity id.
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

// WorkspaceActivityRefreshResult : WorkspaceActivityRefreshResult -.
type WorkspaceActivityRefreshResult struct {
	// Activity id.
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

// WorkspaceActivityTemplate : WorkspaceActivityTemplate -.
type WorkspaceActivityTemplate struct {
	// End time for the activity.
	EndTime *strfmt.DateTime `json:"end_time,omitempty"`

	// LogSummary ...
	LogSummary *LogSummary `json:"log_summary,omitempty"`

	// Log URL.
	LogURL *string `json:"log_url,omitempty"`

	// Message.
	Message *string `json:"message,omitempty"`

	// Activity start time.
	StartTime *strfmt.DateTime `json:"start_time,omitempty"`

	// WorkspaceActivityStatus activity status type.
	Status *string `json:"status,omitempty"`

	// Template id.
	TemplateID *string `json:"template_id,omitempty"`

	// Template type.
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

// WorkspaceActivityTemplateLogs : WorkspaceActivityTemplateLogs -.
type WorkspaceActivityTemplateLogs struct {
	// Log URL.
	LogURL *string `json:"log_url,omitempty"`

	// Template id.
	TemplateID *string `json:"template_id,omitempty"`

	// Template type.
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

// WorkspaceBulkDeleteResponse : WorkspaceBulkDeleteResponse -.
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

// WorkspaceJobResponse : WorkspaceJobResponse -.
type WorkspaceJobResponse struct {
	// JobStatusType -.
	JobStatus *JobStatusType `json:"job_status,omitempty"`
}

// UnmarshalWorkspaceJobResponse unmarshals an instance of WorkspaceJobResponse from the specified map of raw messages.
func UnmarshalWorkspaceJobResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkspaceJobResponse)
	err = core.UnmarshalModel(m, "job_status", &obj.JobStatus, UnmarshalJobStatusType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WorkspaceResponse : WorkspaceResponse - request returned by create.
type WorkspaceResponse struct {
	// List of applied shared dataset id.
	AppliedShareddataIds []string `json:"applied_shareddata_ids,omitempty"`

	// CatalogRef -.
	CatalogRef *CatalogRef `json:"catalog_ref,omitempty"`

	// Workspace created at.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Workspace created by.
	CreatedBy *string `json:"created_by,omitempty"`

	// Workspace CRN.
	Crn *string `json:"crn,omitempty"`

	// Workspace description.
	Description *string `json:"description,omitempty"`

	// Workspace id.
	ID *string `json:"id,omitempty"`

	// Last health checked at.
	LastHealthCheckAt *strfmt.DateTime `json:"last_health_check_at,omitempty"`

	// Workspace location.
	Location *string `json:"location,omitempty"`

	// Workspace name.
	Name *string `json:"name,omitempty"`

	// Workspace resource group.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// Workspace runtime data.
	RuntimeData []TemplateRunTimeDataResponse `json:"runtime_data,omitempty"`

	// SharedTargetDataResponse -.
	SharedData *SharedTargetDataResponse `json:"shared_data,omitempty"`

	// Workspace status type.
	Status *string `json:"status,omitempty"`

	// Workspace tags.
	Tags []string `json:"tags,omitempty"`

	// Workspace template data.
	TemplateData []TemplateSourceDataResponse `json:"template_data,omitempty"`

	// Workspace template ref.
	TemplateRef *string `json:"template_ref,omitempty"`

	// TemplateRepoResponse -.
	TemplateRepo *TemplateRepoResponse `json:"template_repo,omitempty"`

	// List of Workspace type.
	Type []string `json:"type,omitempty"`

	// Workspace updated at.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// Workspace updated by.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// WorkspaceStatusResponse -.
	WorkspaceStatus *WorkspaceStatusResponse `json:"workspace_status,omitempty"`

	// WorkspaceStatusMessage -.
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

// WorkspaceResponseList : WorkspaceResponseList -.
type WorkspaceResponseList struct {
	// Total number of workspaces.
	Count *int64 `json:"count,omitempty"`

	// Limit for the list.
	Limit *int64 `json:"limit" validate:"required"`

	// Offset for the list.
	Offset *int64 `json:"offset" validate:"required"`

	// List of Workspaces.
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

// WorkspaceStatusMessage : WorkspaceStatusMessage -.
type WorkspaceStatusMessage struct {
	// Status code.
	StatusCode *string `json:"status_code,omitempty"`

	// Status message.
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
	// Frozen status.
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

// WorkspaceStatusResponse : WorkspaceStatusResponse -.
type WorkspaceStatusResponse struct {
	// Frozen status.
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

// WorkspaceStatusUpdateRequest : WorkspaceStatusUpdateRequest -.
type WorkspaceStatusUpdateRequest struct {
	// Frozen status.
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

// WorkspaceTemplateValuesResponse : WorkspaceTemplateValuesResponse -.
type WorkspaceTemplateValuesResponse struct {
	// List of runtime data.
	RuntimeData []TemplateRunTimeDataResponse `json:"runtime_data,omitempty"`

	// SharedTargetData -.
	SharedData *SharedTargetData `json:"shared_data,omitempty"`

	// List of source data.
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

// WorkspaceVariableRequest : WorkspaceVariableRequest -.
type WorkspaceVariableRequest struct {
	// Variable description.
	Description *string `json:"description,omitempty"`

	// Variable name.
	Name *string `json:"name,omitempty"`

	// Variable is secure.
	Secure *bool `json:"secure,omitempty"`

	// Variable type.
	Type *string `json:"type,omitempty"`

	// Variable uses default value; and is not over-ridden.
	UseDefault *bool `json:"use_default,omitempty"`

	// Value of the Variable.
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

// WorkspaceVariableResponse : WorkspaceVariableResponse -.
type WorkspaceVariableResponse struct {
	// Variable descrption.
	Description *string `json:"description,omitempty"`

	// Variable name.
	Name *string `json:"name,omitempty"`

	// Variable is secure.
	Secure *bool `json:"secure,omitempty"`

	// Variable type.
	Type *string `json:"type,omitempty"`

	// Value of the Variable.
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
