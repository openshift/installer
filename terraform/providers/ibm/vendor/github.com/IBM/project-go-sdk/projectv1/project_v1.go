/**
 * (C) Copyright IBM Corp. 2023.
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
 * IBM OpenAPI SDK Code Generator Version: 3.71.0-316eb5da-20230504-195406
 */

// Package projectv1 : Operations and models for the ProjectV1 service
package projectv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/project-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// ProjectV1 : This document is the **REST API specification** for the Projects Service. The Projects service provides
// the capability to manage Infrastructure as Code in IBM Cloud.
//
// API Version: 1.0.0
type ProjectV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://projects.api.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "project"

// ProjectV1Options : Service options
type ProjectV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewProjectV1UsingExternalConfig : constructs an instance of ProjectV1 with passed in options and external configuration.
func NewProjectV1UsingExternalConfig(options *ProjectV1Options) (project *ProjectV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	project, err = NewProjectV1(options)
	if err != nil {
		return
	}

	err = project.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = project.Service.SetServiceURL(options.URL)
	}
	return
}

// NewProjectV1 : constructs an instance of ProjectV1 with passed in options.
func NewProjectV1(options *ProjectV1Options) (service *ProjectV1, err error) {
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

	service = &ProjectV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "project" suitable for processing requests.
func (project *ProjectV1) Clone() *ProjectV1 {
	if core.IsNil(project) {
		return nil
	}
	clone := *project
	clone.Service = project.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (project *ProjectV1) SetServiceURL(url string) error {
	return project.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (project *ProjectV1) GetServiceURL() string {
	return project.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (project *ProjectV1) SetDefaultHeaders(headers http.Header) {
	project.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (project *ProjectV1) SetEnableGzipCompression(enableGzip bool) {
	project.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (project *ProjectV1) GetEnableGzipCompression() bool {
	return project.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (project *ProjectV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	project.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (project *ProjectV1) DisableRetries() {
	project.Service.DisableRetries()
}

// CreateProject : Create a project
// Create a new project and asynchronously setup the tools to manage it. Add a deployable architecture by customizing
// the configuration. After the changes are validated and approved, deploy the resources that the project configures.
func (project *ProjectV1) CreateProject(createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	return project.CreateProjectWithContext(context.Background(), createProjectOptions)
}

// CreateProjectWithContext is an alternate form of the CreateProject method which supports a Context parameter
func (project *ProjectV1) CreateProjectWithContext(ctx context.Context, createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProjectOptions, "createProjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProjectOptions, "createProjectOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "CreateProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	builder.AddQuery("resource_group", fmt.Sprint(*createProjectOptions.ResourceGroup))
	builder.AddQuery("location", fmt.Sprint(*createProjectOptions.Location))

	body := make(map[string]interface{})
	if createProjectOptions.Name != nil {
		body["name"] = createProjectOptions.Name
	}
	if createProjectOptions.Description != nil {
		body["description"] = createProjectOptions.Description
	}
	if createProjectOptions.Configs != nil {
		body["configs"] = createProjectOptions.Configs
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProjects : List projects
// List existing projects. Projects are sorted by ID.
func (project *ProjectV1) ListProjects(listProjectsOptions *ListProjectsOptions) (result *ProjectCollection, response *core.DetailedResponse, err error) {
	return project.ListProjectsWithContext(context.Background(), listProjectsOptions)
}

// ListProjectsWithContext is an alternate form of the ListProjects method which supports a Context parameter
func (project *ProjectV1) ListProjectsWithContext(ctx context.Context, listProjectsOptions *ListProjectsOptions) (result *ProjectCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProjectsOptions, "listProjectsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProjectsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ListProjects")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProjectsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listProjectsOptions.Start))
	}
	if listProjectsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProjectsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProject : Get a project
// Get information about a project.
func (project *ProjectV1) GetProject(getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	return project.GetProjectWithContext(context.Background(), getProjectOptions)
}

// GetProjectWithContext is an alternate form of the GetProject method which supports a Context parameter
func (project *ProjectV1) GetProjectWithContext(ctx context.Context, getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectOptions, "getProjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProjectOptions, "getProjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getProjectOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateProject : Update a project
// Update a project by the ID.
func (project *ProjectV1) UpdateProject(updateProjectOptions *UpdateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	return project.UpdateProjectWithContext(context.Background(), updateProjectOptions)
}

// UpdateProjectWithContext is an alternate form of the UpdateProject method which supports a Context parameter
func (project *ProjectV1) UpdateProjectWithContext(ctx context.Context, updateProjectOptions *UpdateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProjectOptions, "updateProjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProjectOptions, "updateProjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateProjectOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UpdateProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateProjectOptions.Name != nil {
		body["name"] = updateProjectOptions.Name
	}
	if updateProjectOptions.Description != nil {
		body["description"] = updateProjectOptions.Description
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProject : Delete a project
// Delete a project document by the ID. A project can only be deleted after deleting all of its artifacts.
func (project *ProjectV1) DeleteProject(deleteProjectOptions *DeleteProjectOptions) (response *core.DetailedResponse, err error) {
	return project.DeleteProjectWithContext(context.Background(), deleteProjectOptions)
}

// DeleteProjectWithContext is an alternate form of the DeleteProject method which supports a Context parameter
func (project *ProjectV1) DeleteProjectWithContext(ctx context.Context, deleteProjectOptions *DeleteProjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProjectOptions, "deleteProjectOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProjectOptions, "deleteProjectOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteProjectOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteProjectOptions.Destroy != nil {
		builder.AddQuery("destroy", fmt.Sprint(*deleteProjectOptions.Destroy))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = project.Service.Request(request, nil)

	return
}

// CreateConfig : Add a new configuration
// Add a new configuration to a project.
func (project *ProjectV1) CreateConfig(createConfigOptions *CreateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.CreateConfigWithContext(context.Background(), createConfigOptions)
}

// CreateConfigWithContext is an alternate form of the CreateConfig method which supports a Context parameter
func (project *ProjectV1) CreateConfigWithContext(ctx context.Context, createConfigOptions *CreateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigOptions, "createConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createConfigOptions, "createConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createConfigOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "CreateConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createConfigOptions.Name != nil {
		body["name"] = createConfigOptions.Name
	}
	if createConfigOptions.LocatorID != nil {
		body["locator_id"] = createConfigOptions.LocatorID
	}
	if createConfigOptions.ID != nil {
		body["id"] = createConfigOptions.ID
	}
	if createConfigOptions.Labels != nil {
		body["labels"] = createConfigOptions.Labels
	}
	if createConfigOptions.Description != nil {
		body["description"] = createConfigOptions.Description
	}
	if createConfigOptions.Input != nil {
		body["input"] = createConfigOptions.Input
	}
	if createConfigOptions.Setting != nil {
		body["setting"] = createConfigOptions.Setting
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListConfigs : List all project configurations
// The collection of configurations that are returned.
func (project *ProjectV1) ListConfigs(listConfigsOptions *ListConfigsOptions) (result *ProjectConfigCollection, response *core.DetailedResponse, err error) {
	return project.ListConfigsWithContext(context.Background(), listConfigsOptions)
}

// ListConfigsWithContext is an alternate form of the ListConfigs method which supports a Context parameter
func (project *ProjectV1) ListConfigsWithContext(ctx context.Context, listConfigsOptions *ListConfigsOptions) (result *ProjectConfigCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigsOptions, "listConfigsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listConfigsOptions, "listConfigsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listConfigsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listConfigsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ListConfigs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listConfigsOptions.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*listConfigsOptions.Version))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfig : Get a project configuration
// Returns the specified project configuration in a specific project.
func (project *ProjectV1) GetConfig(getConfigOptions *GetConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.GetConfigWithContext(context.Background(), getConfigOptions)
}

// GetConfigWithContext is an alternate form of the GetConfig method which supports a Context parameter
func (project *ProjectV1) GetConfigWithContext(ctx context.Context, getConfigOptions *GetConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigOptions, "getConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigOptions, "getConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getConfigOptions.ProjectID,
		"id": *getConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getConfigOptions.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*getConfigOptions.Version))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateConfig : Update a configuration
// Update a configuration in a project by the ID.
func (project *ProjectV1) UpdateConfig(updateConfigOptions *UpdateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.UpdateConfigWithContext(context.Background(), updateConfigOptions)
}

// UpdateConfigWithContext is an alternate form of the UpdateConfig method which supports a Context parameter
func (project *ProjectV1) UpdateConfigWithContext(ctx context.Context, updateConfigOptions *UpdateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateConfigOptions, "updateConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateConfigOptions, "updateConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateConfigOptions.ProjectID,
		"id": *updateConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UpdateConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(updateConfigOptions.ProjectConfig)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfig : Delete a configuration in a project by ID
// Delete a configuration in a project. Deleting the configuration will also destroy all the resources deployed by the
// configuration if the query parameter `destroy` is specified.
func (project *ProjectV1) DeleteConfig(deleteConfigOptions *DeleteConfigOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	return project.DeleteConfigWithContext(context.Background(), deleteConfigOptions)
}

// DeleteConfigWithContext is an alternate form of the DeleteConfig method which supports a Context parameter
func (project *ProjectV1) DeleteConfigWithContext(ctx context.Context, deleteConfigOptions *DeleteConfigOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigOptions, "deleteConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteConfigOptions, "deleteConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteConfigOptions.ProjectID,
		"id": *deleteConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteConfigOptions.DraftOnly != nil {
		builder.AddQuery("draft_only", fmt.Sprint(*deleteConfigOptions.DraftOnly))
	}
	if deleteConfigOptions.Destroy != nil {
		builder.AddQuery("destroy", fmt.Sprint(*deleteConfigOptions.Destroy))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigDelete)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfigDiff : Get a diff summary of a project configuration
// Returns a diff summary of the specified project configuration between its current draft and active version of a
// specific project.
func (project *ProjectV1) GetConfigDiff(getConfigDiffOptions *GetConfigDiffOptions) (result *ProjectConfigDiff, response *core.DetailedResponse, err error) {
	return project.GetConfigDiffWithContext(context.Background(), getConfigDiffOptions)
}

// GetConfigDiffWithContext is an alternate form of the GetConfigDiff method which supports a Context parameter
func (project *ProjectV1) GetConfigDiffWithContext(ctx context.Context, getConfigDiffOptions *GetConfigDiffOptions) (result *ProjectConfigDiff, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigDiffOptions, "getConfigDiffOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigDiffOptions, "getConfigDiffOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getConfigDiffOptions.ProjectID,
		"id": *getConfigDiffOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/diff`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigDiffOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetConfigDiff")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigDiff)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ForceApprove : Force approve project configuration
// Force approve configuration edits to the main configuration with an approving comment.
func (project *ProjectV1) ForceApprove(forceApproveOptions *ForceApproveOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.ForceApproveWithContext(context.Background(), forceApproveOptions)
}

// ForceApproveWithContext is an alternate form of the ForceApprove method which supports a Context parameter
func (project *ProjectV1) ForceApproveWithContext(ctx context.Context, forceApproveOptions *ForceApproveOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(forceApproveOptions, "forceApproveOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(forceApproveOptions, "forceApproveOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *forceApproveOptions.ProjectID,
		"id": *forceApproveOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/force_approve`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range forceApproveOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ForceApprove")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if forceApproveOptions.Comment != nil {
		body["comment"] = forceApproveOptions.Comment
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Approve : Approve and merge a configuration draft
// Approve and merge configuration edits to the main configuration.
func (project *ProjectV1) Approve(approveOptions *ApproveOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.ApproveWithContext(context.Background(), approveOptions)
}

// ApproveWithContext is an alternate form of the Approve method which supports a Context parameter
func (project *ProjectV1) ApproveWithContext(ctx context.Context, approveOptions *ApproveOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(approveOptions, "approveOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(approveOptions, "approveOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *approveOptions.ProjectID,
		"id": *approveOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/approve`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range approveOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "Approve")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if approveOptions.Comment != nil {
		body["comment"] = approveOptions.Comment
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CheckConfig : Run a validation check
// Run a validation check on a given configuration in project. The check includes creating or updating the associated
// schematics workspace with a plan job, running the CRA scans, and cost estimatation.
func (project *ProjectV1) CheckConfig(checkConfigOptions *CheckConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.CheckConfigWithContext(context.Background(), checkConfigOptions)
}

// CheckConfigWithContext is an alternate form of the CheckConfig method which supports a Context parameter
func (project *ProjectV1) CheckConfigWithContext(ctx context.Context, checkConfigOptions *CheckConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(checkConfigOptions, "checkConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(checkConfigOptions, "checkConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *checkConfigOptions.ProjectID,
		"id": *checkConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/check`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range checkConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "CheckConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if checkConfigOptions.XAuthRefreshToken != nil {
		builder.AddHeader("X-Auth-Refresh-Token", fmt.Sprint(*checkConfigOptions.XAuthRefreshToken))
	}

	if checkConfigOptions.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*checkConfigOptions.Version))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// InstallConfig : Deploy a configuration
// Deploy a project's configuration. It's an asynchronous operation that can be tracked using the get project
// configuration API with full metadata.
func (project *ProjectV1) InstallConfig(installConfigOptions *InstallConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	return project.InstallConfigWithContext(context.Background(), installConfigOptions)
}

// InstallConfigWithContext is an alternate form of the InstallConfig method which supports a Context parameter
func (project *ProjectV1) InstallConfigWithContext(ctx context.Context, installConfigOptions *InstallConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(installConfigOptions, "installConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(installConfigOptions, "installConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *installConfigOptions.ProjectID,
		"id": *installConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/install`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range installConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "InstallConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UninstallConfig : Destroy configuration resources
// Destroy a project's configuration resources. The operation destroys all the resources that are deployed with the
// specific configuration. You can track it by using the get project configuration API with full metadata.
func (project *ProjectV1) UninstallConfig(uninstallConfigOptions *UninstallConfigOptions) (response *core.DetailedResponse, err error) {
	return project.UninstallConfigWithContext(context.Background(), uninstallConfigOptions)
}

// UninstallConfigWithContext is an alternate form of the UninstallConfig method which supports a Context parameter
func (project *ProjectV1) UninstallConfigWithContext(ctx context.Context, uninstallConfigOptions *UninstallConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(uninstallConfigOptions, "uninstallConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(uninstallConfigOptions, "uninstallConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *uninstallConfigOptions.ProjectID,
		"id": *uninstallConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/uninstall`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uninstallConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UninstallConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = project.Service.Request(request, nil)

	return
}

// GetSchematicsJob : View the latest schematics job
// Fetch and find the latest schematics job that corresponds to a plan, deploy, or destroy configuration resource
// action.
func (project *ProjectV1) GetSchematicsJob(getSchematicsJobOptions *GetSchematicsJobOptions) (result *ActionJob, response *core.DetailedResponse, err error) {
	return project.GetSchematicsJobWithContext(context.Background(), getSchematicsJobOptions)
}

// GetSchematicsJobWithContext is an alternate form of the GetSchematicsJob method which supports a Context parameter
func (project *ProjectV1) GetSchematicsJobWithContext(ctx context.Context, getSchematicsJobOptions *GetSchematicsJobOptions) (result *ActionJob, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSchematicsJobOptions, "getSchematicsJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSchematicsJobOptions, "getSchematicsJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getSchematicsJobOptions.ProjectID,
		"id": *getSchematicsJobOptions.ID,
		"action": *getSchematicsJobOptions.Action,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/job/{action}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSchematicsJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetSchematicsJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getSchematicsJobOptions.Since != nil {
		builder.AddQuery("since", fmt.Sprint(*getSchematicsJobOptions.Since))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalActionJob)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetCostEstimate : Get the cost estimate
// Retrieve the cost estimate for a configuraton.
func (project *ProjectV1) GetCostEstimate(getCostEstimateOptions *GetCostEstimateOptions) (result *CostEstimate, response *core.DetailedResponse, err error) {
	return project.GetCostEstimateWithContext(context.Background(), getCostEstimateOptions)
}

// GetCostEstimateWithContext is an alternate form of the GetCostEstimate method which supports a Context parameter
func (project *ProjectV1) GetCostEstimateWithContext(ctx context.Context, getCostEstimateOptions *GetCostEstimateOptions) (result *CostEstimate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCostEstimateOptions, "getCostEstimateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCostEstimateOptions, "getCostEstimateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getCostEstimateOptions.ProjectID,
		"id": *getCostEstimateOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/cost_estimate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCostEstimateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetCostEstimate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCostEstimateOptions.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*getCostEstimateOptions.Version))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCostEstimate)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostCrnToken : Creates a project CRN token
// Refreshes a project CRN token by creating a new one.
func (project *ProjectV1) PostCrnToken(postCrnTokenOptions *PostCrnTokenOptions) (result *ProjectCRNTokenResponse, response *core.DetailedResponse, err error) {
	return project.PostCrnTokenWithContext(context.Background(), postCrnTokenOptions)
}

// PostCrnTokenWithContext is an alternate form of the PostCrnToken method which supports a Context parameter
func (project *ProjectV1) PostCrnTokenWithContext(ctx context.Context, postCrnTokenOptions *PostCrnTokenOptions) (result *ProjectCRNTokenResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postCrnTokenOptions, "postCrnTokenOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postCrnTokenOptions, "postCrnTokenOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *postCrnTokenOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/token`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postCrnTokenOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "PostCrnToken")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectCRNTokenResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostNotification : Add notifications
// Creates a notification event to be stored on the project definition.
func (project *ProjectV1) PostNotification(postNotificationOptions *PostNotificationOptions) (result *NotificationsPrototypePostResponse, response *core.DetailedResponse, err error) {
	return project.PostNotificationWithContext(context.Background(), postNotificationOptions)
}

// PostNotificationWithContext is an alternate form of the PostNotification method which supports a Context parameter
func (project *ProjectV1) PostNotificationWithContext(ctx context.Context, postNotificationOptions *PostNotificationOptions) (result *NotificationsPrototypePostResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postNotificationOptions, "postNotificationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postNotificationOptions, "postNotificationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *postNotificationOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postNotificationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "PostNotification")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postNotificationOptions.Notifications != nil {
		body["notifications"] = postNotificationOptions.Notifications
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsPrototypePostResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetNotifications : Get events by project ID
// Get all the notification events from a specific project ID.
func (project *ProjectV1) GetNotifications(getNotificationsOptions *GetNotificationsOptions) (result *NotificationsGetResponse, response *core.DetailedResponse, err error) {
	return project.GetNotificationsWithContext(context.Background(), getNotificationsOptions)
}

// GetNotificationsWithContext is an alternate form of the GetNotifications method which supports a Context parameter
func (project *ProjectV1) GetNotificationsWithContext(ctx context.Context, getNotificationsOptions *GetNotificationsOptions) (result *NotificationsGetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getNotificationsOptions, "getNotificationsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getNotificationsOptions, "getNotificationsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getNotificationsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getNotificationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetNotifications")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsGetResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostEventNotificationsIntegration : Connect to a event notifications instance
// Connects a project instance to an event notifications instance.
func (project *ProjectV1) PostEventNotificationsIntegration(postEventNotificationsIntegrationOptions *PostEventNotificationsIntegrationOptions) (result *NotificationsIntegrationPostResponse, response *core.DetailedResponse, err error) {
	return project.PostEventNotificationsIntegrationWithContext(context.Background(), postEventNotificationsIntegrationOptions)
}

// PostEventNotificationsIntegrationWithContext is an alternate form of the PostEventNotificationsIntegration method which supports a Context parameter
func (project *ProjectV1) PostEventNotificationsIntegrationWithContext(ctx context.Context, postEventNotificationsIntegrationOptions *PostEventNotificationsIntegrationOptions) (result *NotificationsIntegrationPostResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postEventNotificationsIntegrationOptions, "postEventNotificationsIntegrationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postEventNotificationsIntegrationOptions, "postEventNotificationsIntegrationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *postEventNotificationsIntegrationOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event_notifications`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postEventNotificationsIntegrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "PostEventNotificationsIntegration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postEventNotificationsIntegrationOptions.InstanceCrn != nil {
		body["instance_crn"] = postEventNotificationsIntegrationOptions.InstanceCrn
	}
	if postEventNotificationsIntegrationOptions.Description != nil {
		body["description"] = postEventNotificationsIntegrationOptions.Description
	}
	if postEventNotificationsIntegrationOptions.EventNotificationsSourceName != nil {
		body["event_notifications_source_name"] = postEventNotificationsIntegrationOptions.EventNotificationsSourceName
	}
	if postEventNotificationsIntegrationOptions.Enabled != nil {
		body["enabled"] = postEventNotificationsIntegrationOptions.Enabled
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsIntegrationPostResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetEventNotificationsIntegration : Get event notification source details by project ID
// Gets the source details of the project from the connect event notifications instance.
func (project *ProjectV1) GetEventNotificationsIntegration(getEventNotificationsIntegrationOptions *GetEventNotificationsIntegrationOptions) (result *NotificationsIntegrationGetResponse, response *core.DetailedResponse, err error) {
	return project.GetEventNotificationsIntegrationWithContext(context.Background(), getEventNotificationsIntegrationOptions)
}

// GetEventNotificationsIntegrationWithContext is an alternate form of the GetEventNotificationsIntegration method which supports a Context parameter
func (project *ProjectV1) GetEventNotificationsIntegrationWithContext(ctx context.Context, getEventNotificationsIntegrationOptions *GetEventNotificationsIntegrationOptions) (result *NotificationsIntegrationGetResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEventNotificationsIntegrationOptions, "getEventNotificationsIntegrationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEventNotificationsIntegrationOptions, "getEventNotificationsIntegrationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getEventNotificationsIntegrationOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event_notifications`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEventNotificationsIntegrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetEventNotificationsIntegration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsIntegrationGetResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteEventNotificationsIntegration : Delete an event notifications connection
// Deletes the event notifications integration if that is where the project was onboarded to.
func (project *ProjectV1) DeleteEventNotificationsIntegration(deleteEventNotificationsIntegrationOptions *DeleteEventNotificationsIntegrationOptions) (response *core.DetailedResponse, err error) {
	return project.DeleteEventNotificationsIntegrationWithContext(context.Background(), deleteEventNotificationsIntegrationOptions)
}

// DeleteEventNotificationsIntegrationWithContext is an alternate form of the DeleteEventNotificationsIntegration method which supports a Context parameter
func (project *ProjectV1) DeleteEventNotificationsIntegrationWithContext(ctx context.Context, deleteEventNotificationsIntegrationOptions *DeleteEventNotificationsIntegrationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEventNotificationsIntegrationOptions, "deleteEventNotificationsIntegrationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEventNotificationsIntegrationOptions, "deleteEventNotificationsIntegrationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteEventNotificationsIntegrationOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event_notifications`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEventNotificationsIntegrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteEventNotificationsIntegration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = project.Service.Request(request, nil)

	return
}

// PostTestEventNotification : Send notification to event notifications instance
// Sends a notification to the event notifications instance.
func (project *ProjectV1) PostTestEventNotification(postTestEventNotificationOptions *PostTestEventNotificationOptions) (result *NotificationsIntegrationTestPostResponse, response *core.DetailedResponse, err error) {
	return project.PostTestEventNotificationWithContext(context.Background(), postTestEventNotificationOptions)
}

// PostTestEventNotificationWithContext is an alternate form of the PostTestEventNotification method which supports a Context parameter
func (project *ProjectV1) PostTestEventNotificationWithContext(ctx context.Context, postTestEventNotificationOptions *PostTestEventNotificationOptions) (result *NotificationsIntegrationTestPostResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postTestEventNotificationOptions, "postTestEventNotificationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postTestEventNotificationOptions, "postTestEventNotificationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *postTestEventNotificationOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{id}/event_notifications/test`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postTestEventNotificationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "PostTestEventNotification")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if postTestEventNotificationOptions.Ibmendefaultlong != nil {
		body["ibmendefaultlong"] = postTestEventNotificationOptions.Ibmendefaultlong
	}
	if postTestEventNotificationOptions.Ibmendefaultshort != nil {
		body["ibmendefaultshort"] = postTestEventNotificationOptions.Ibmendefaultshort
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalNotificationsIntegrationTestPostResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActionJob : The response of a fetching an action job.
type ActionJob struct {
	// The unique ID of a project.
	ID *string `json:"id,omitempty"`
}

// UnmarshalActionJob unmarshals an instance of ActionJob from the specified map of raw messages.
func UnmarshalActionJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJob)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApproveOptions : The Approve options.
type ApproveOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Notes on the project draft action.
	Comment *string `json:"comment,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewApproveOptions : Instantiate ApproveOptions
func (*ProjectV1) NewApproveOptions(projectID string, id string) *ApproveOptions {
	return &ApproveOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ApproveOptions) SetProjectID(projectID string) *ApproveOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ApproveOptions) SetID(id string) *ApproveOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetComment : Allow user to set Comment
func (_options *ApproveOptions) SetComment(comment string) *ApproveOptions {
	_options.Comment = core.StringPtr(comment)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ApproveOptions) SetHeaders(param map[string]string) *ApproveOptions {
	options.Headers = param
	return options
}

// CheckConfigOptions : The CheckConfig options.
type CheckConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The IAM refresh token.
	XAuthRefreshToken *string `json:"X-Auth-Refresh-Token,omitempty"`

	// The version of the configuration that the validation check should trigger against.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCheckConfigOptions : Instantiate CheckConfigOptions
func (*ProjectV1) NewCheckConfigOptions(projectID string, id string) *CheckConfigOptions {
	return &CheckConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CheckConfigOptions) SetProjectID(projectID string) *CheckConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CheckConfigOptions) SetID(id string) *CheckConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetXAuthRefreshToken : Allow user to set XAuthRefreshToken
func (_options *CheckConfigOptions) SetXAuthRefreshToken(xAuthRefreshToken string) *CheckConfigOptions {
	_options.XAuthRefreshToken = core.StringPtr(xAuthRefreshToken)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CheckConfigOptions) SetVersion(version string) *CheckConfigOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CheckConfigOptions) SetHeaders(param map[string]string) *CheckConfigOptions {
	options.Headers = param
	return options
}

// CostEstimate : The cost estimate for the given configuration.
type CostEstimate struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of CostEstimate
func (o *CostEstimate) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of CostEstimate
func (o *CostEstimate) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of CostEstimate
func (o *CostEstimate) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of CostEstimate
func (o *CostEstimate) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of CostEstimate
func (o *CostEstimate) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalCostEstimate unmarshals an instance of CostEstimate from the specified map of raw messages.
func UnmarshalCostEstimate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CostEstimate)
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

// CreateConfigOptions : The CreateConfig options.
type CreateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The configuration name.
	Name *string `json:"name" validate:"required"`

	// A dotted value of catalogID.versionID.
	LocatorID *string `json:"locator_id" validate:"required"`

	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id,omitempty"`

	// A collection of configuration labels.
	Labels []string `json:"labels,omitempty"`

	// The project configuration description.
	Description *string `json:"description,omitempty"`

	// The input values to use to deploy the configuration.
	Input []ProjectConfigInputVariable `json:"input,omitempty"`

	// Schematics environment variables to use to deploy the configuration.
	Setting []ProjectConfigSettingCollection `json:"setting,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateConfigOptions : Instantiate CreateConfigOptions
func (*ProjectV1) NewCreateConfigOptions(projectID string, name string, locatorID string) *CreateConfigOptions {
	return &CreateConfigOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		LocatorID: core.StringPtr(locatorID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateConfigOptions) SetProjectID(projectID string) *CreateConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateConfigOptions) SetName(name string) *CreateConfigOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetLocatorID : Allow user to set LocatorID
func (_options *CreateConfigOptions) SetLocatorID(locatorID string) *CreateConfigOptions {
	_options.LocatorID = core.StringPtr(locatorID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateConfigOptions) SetID(id string) *CreateConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *CreateConfigOptions) SetLabels(labels []string) *CreateConfigOptions {
	_options.Labels = labels
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateConfigOptions) SetDescription(description string) *CreateConfigOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetInput : Allow user to set Input
func (_options *CreateConfigOptions) SetInput(input []ProjectConfigInputVariable) *CreateConfigOptions {
	_options.Input = input
	return _options
}

// SetSetting : Allow user to set Setting
func (_options *CreateConfigOptions) SetSetting(setting []ProjectConfigSettingCollection) *CreateConfigOptions {
	_options.Setting = setting
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigOptions) SetHeaders(param map[string]string) *CreateConfigOptions {
	options.Headers = param
	return options
}

// CreateProjectOptions : The CreateProject options.
type CreateProjectOptions struct {
	// The resource group where the project's data and tools are created.
	ResourceGroup *string `json:"resource_group" validate:"required"`

	// The location where the project's data and tools are created.
	Location *string `json:"location" validate:"required"`

	// The project name.
	Name *string `json:"name" validate:"required"`

	// A project's descriptive text.
	Description *string `json:"description,omitempty"`

	// The project configurations.
	Configs []ProjectConfigPrototype `json:"configs,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProjectOptions : Instantiate CreateProjectOptions
func (*ProjectV1) NewCreateProjectOptions(resourceGroup string, location string, name string) *CreateProjectOptions {
	return &CreateProjectOptions{
		ResourceGroup: core.StringPtr(resourceGroup),
		Location: core.StringPtr(location),
		Name: core.StringPtr(name),
	}
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateProjectOptions) SetResourceGroup(resourceGroup string) *CreateProjectOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateProjectOptions) SetLocation(location string) *CreateProjectOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateProjectOptions) SetName(name string) *CreateProjectOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateProjectOptions) SetDescription(description string) *CreateProjectOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetConfigs : Allow user to set Configs
func (_options *CreateProjectOptions) SetConfigs(configs []ProjectConfigPrototype) *CreateProjectOptions {
	_options.Configs = configs
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProjectOptions) SetHeaders(param map[string]string) *CreateProjectOptions {
	options.Headers = param
	return options
}

// CumulativeNeedsAttention : CumulativeNeedsAttention struct
type CumulativeNeedsAttention struct {
	// The event name.
	Event *string `json:"event,omitempty"`

	// The unique ID of a project.
	EventID *string `json:"event_id,omitempty"`

	// The unique ID of a project.
	ConfigID *string `json:"config_id,omitempty"`

	// The version number of the configuration.
	ConfigVersion *int64 `json:"config_version,omitempty"`
}

// UnmarshalCumulativeNeedsAttention unmarshals an instance of CumulativeNeedsAttention from the specified map of raw messages.
func UnmarshalCumulativeNeedsAttention(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CumulativeNeedsAttention)
	err = core.UnmarshalPrimitive(m, "event", &obj.Event)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_id", &obj.EventID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "config_id", &obj.ConfigID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "config_version", &obj.ConfigVersion)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteConfigOptions : The DeleteConfig options.
type DeleteConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The flag to determine if only the draft version should be deleted.
	DraftOnly *bool `json:"draft_only,omitempty"`

	// The flag that indicates if the resources deployed by schematics should be destroyed.
	Destroy *bool `json:"destroy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteConfigOptions : Instantiate DeleteConfigOptions
func (*ProjectV1) NewDeleteConfigOptions(projectID string, id string) *DeleteConfigOptions {
	return &DeleteConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteConfigOptions) SetProjectID(projectID string) *DeleteConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteConfigOptions) SetID(id string) *DeleteConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetDraftOnly : Allow user to set DraftOnly
func (_options *DeleteConfigOptions) SetDraftOnly(draftOnly bool) *DeleteConfigOptions {
	_options.DraftOnly = core.BoolPtr(draftOnly)
	return _options
}

// SetDestroy : Allow user to set Destroy
func (_options *DeleteConfigOptions) SetDestroy(destroy bool) *DeleteConfigOptions {
	_options.Destroy = core.BoolPtr(destroy)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigOptions) SetHeaders(param map[string]string) *DeleteConfigOptions {
	options.Headers = param
	return options
}

// DeleteEventNotificationsIntegrationOptions : The DeleteEventNotificationsIntegration options.
type DeleteEventNotificationsIntegrationOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEventNotificationsIntegrationOptions : Instantiate DeleteEventNotificationsIntegrationOptions
func (*ProjectV1) NewDeleteEventNotificationsIntegrationOptions(id string) *DeleteEventNotificationsIntegrationOptions {
	return &DeleteEventNotificationsIntegrationOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteEventNotificationsIntegrationOptions) SetID(id string) *DeleteEventNotificationsIntegrationOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEventNotificationsIntegrationOptions) SetHeaders(param map[string]string) *DeleteEventNotificationsIntegrationOptions {
	options.Headers = param
	return options
}

// DeleteProjectOptions : The DeleteProject options.
type DeleteProjectOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// The flag that indicates if the resources deployed by schematics should be destroyed.
	Destroy *bool `json:"destroy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProjectOptions : Instantiate DeleteProjectOptions
func (*ProjectV1) NewDeleteProjectOptions(id string) *DeleteProjectOptions {
	return &DeleteProjectOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteProjectOptions) SetID(id string) *DeleteProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetDestroy : Allow user to set Destroy
func (_options *DeleteProjectOptions) SetDestroy(destroy bool) *DeleteProjectOptions {
	_options.Destroy = core.BoolPtr(destroy)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProjectOptions) SetHeaders(param map[string]string) *DeleteProjectOptions {
	options.Headers = param
	return options
}

// ForceApproveOptions : The ForceApprove options.
type ForceApproveOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Notes on the project draft action.
	Comment *string `json:"comment,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewForceApproveOptions : Instantiate ForceApproveOptions
func (*ProjectV1) NewForceApproveOptions(projectID string, id string) *ForceApproveOptions {
	return &ForceApproveOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ForceApproveOptions) SetProjectID(projectID string) *ForceApproveOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ForceApproveOptions) SetID(id string) *ForceApproveOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetComment : Allow user to set Comment
func (_options *ForceApproveOptions) SetComment(comment string) *ForceApproveOptions {
	_options.Comment = core.StringPtr(comment)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ForceApproveOptions) SetHeaders(param map[string]string) *ForceApproveOptions {
	options.Headers = param
	return options
}

// GetConfigDiffOptions : The GetConfigDiff options.
type GetConfigDiffOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetConfigDiffOptions : Instantiate GetConfigDiffOptions
func (*ProjectV1) NewGetConfigDiffOptions(projectID string, id string) *GetConfigDiffOptions {
	return &GetConfigDiffOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetConfigDiffOptions) SetProjectID(projectID string) *GetConfigDiffOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetConfigDiffOptions) SetID(id string) *GetConfigDiffOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigDiffOptions) SetHeaders(param map[string]string) *GetConfigDiffOptions {
	options.Headers = param
	return options
}

// GetConfigOptions : The GetConfig options.
type GetConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The version of the configuration to return.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetConfigOptions : Instantiate GetConfigOptions
func (*ProjectV1) NewGetConfigOptions(projectID string, id string) *GetConfigOptions {
	return &GetConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetConfigOptions) SetProjectID(projectID string) *GetConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetConfigOptions) SetID(id string) *GetConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetConfigOptions) SetVersion(version string) *GetConfigOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigOptions) SetHeaders(param map[string]string) *GetConfigOptions {
	options.Headers = param
	return options
}

// GetCostEstimateOptions : The GetCostEstimate options.
type GetCostEstimateOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The version of the configuration that the cost estimate will fetch.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCostEstimateOptions : Instantiate GetCostEstimateOptions
func (*ProjectV1) NewGetCostEstimateOptions(projectID string, id string) *GetCostEstimateOptions {
	return &GetCostEstimateOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetCostEstimateOptions) SetProjectID(projectID string) *GetCostEstimateOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetCostEstimateOptions) SetID(id string) *GetCostEstimateOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetCostEstimateOptions) SetVersion(version string) *GetCostEstimateOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCostEstimateOptions) SetHeaders(param map[string]string) *GetCostEstimateOptions {
	options.Headers = param
	return options
}

// GetEventNotificationsIntegrationOptions : The GetEventNotificationsIntegration options.
type GetEventNotificationsIntegrationOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEventNotificationsIntegrationOptions : Instantiate GetEventNotificationsIntegrationOptions
func (*ProjectV1) NewGetEventNotificationsIntegrationOptions(id string) *GetEventNotificationsIntegrationOptions {
	return &GetEventNotificationsIntegrationOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetEventNotificationsIntegrationOptions) SetID(id string) *GetEventNotificationsIntegrationOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetEventNotificationsIntegrationOptions) SetHeaders(param map[string]string) *GetEventNotificationsIntegrationOptions {
	options.Headers = param
	return options
}

// GetNotificationsOptions : The GetNotifications options.
type GetNotificationsOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetNotificationsOptions : Instantiate GetNotificationsOptions
func (*ProjectV1) NewGetNotificationsOptions(id string) *GetNotificationsOptions {
	return &GetNotificationsOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetNotificationsOptions) SetID(id string) *GetNotificationsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetNotificationsOptions) SetHeaders(param map[string]string) *GetNotificationsOptions {
	options.Headers = param
	return options
}

// GetProjectOptions : The GetProject options.
type GetProjectOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProjectOptions : Instantiate GetProjectOptions
func (*ProjectV1) NewGetProjectOptions(id string) *GetProjectOptions {
	return &GetProjectOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetProjectOptions) SetID(id string) *GetProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProjectOptions) SetHeaders(param map[string]string) *GetProjectOptions {
	options.Headers = param
	return options
}

// GetSchematicsJobOptions : The GetSchematicsJob options.
type GetSchematicsJobOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The triggered action.
	Action *string `json:"action" validate:"required,ne="`

	// The timestamp of when the action was triggered.
	Since *int64 `json:"since,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetSchematicsJobOptions.Action property.
// The triggered action.
const (
	GetSchematicsJobOptions_Action_Install = "install"
	GetSchematicsJobOptions_Action_Plan = "plan"
	GetSchematicsJobOptions_Action_Uninstall = "uninstall"
)

// NewGetSchematicsJobOptions : Instantiate GetSchematicsJobOptions
func (*ProjectV1) NewGetSchematicsJobOptions(projectID string, id string, action string) *GetSchematicsJobOptions {
	return &GetSchematicsJobOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Action: core.StringPtr(action),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetSchematicsJobOptions) SetProjectID(projectID string) *GetSchematicsJobOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetSchematicsJobOptions) SetID(id string) *GetSchematicsJobOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAction : Allow user to set Action
func (_options *GetSchematicsJobOptions) SetAction(action string) *GetSchematicsJobOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetSince : Allow user to set Since
func (_options *GetSchematicsJobOptions) SetSince(since int64) *GetSchematicsJobOptions {
	_options.Since = core.Int64Ptr(since)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSchematicsJobOptions) SetHeaders(param map[string]string) *GetSchematicsJobOptions {
	options.Headers = param
	return options
}

// InputVariable : InputVariable struct
type InputVariable struct {
	// The variable name.
	Name *string `json:"name" validate:"required"`

	// The variable type.
	Type *string `json:"type" validate:"required"`

	// Can be any value - a string, number, boolean, array, or object.
	Value interface{} `json:"value,omitempty"`

	// Whether the variable is required or not.
	Required *bool `json:"required,omitempty"`
}

// Constants associated with the InputVariable.Type property.
// The variable type.
const (
	InputVariable_Type_Array = "array"
	InputVariable_Type_Boolean = "boolean"
	InputVariable_Type_Float = "float"
	InputVariable_Type_Int = "int"
	InputVariable_Type_Number = "number"
	InputVariable_Type_Object = "object"
	InputVariable_Type_Password = "password"
	InputVariable_Type_String = "string"
)

// UnmarshalInputVariable unmarshals an instance of InputVariable from the specified map of raw messages.
func UnmarshalInputVariable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(InputVariable)
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
	err = core.UnmarshalPrimitive(m, "required", &obj.Required)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// InstallConfigOptions : The InstallConfig options.
type InstallConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewInstallConfigOptions : Instantiate InstallConfigOptions
func (*ProjectV1) NewInstallConfigOptions(projectID string, id string) *InstallConfigOptions {
	return &InstallConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *InstallConfigOptions) SetProjectID(projectID string) *InstallConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *InstallConfigOptions) SetID(id string) *InstallConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *InstallConfigOptions) SetHeaders(param map[string]string) *InstallConfigOptions {
	options.Headers = param
	return options
}

// ListConfigsOptions : The ListConfigs options.
type ListConfigsOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The version of configuration to return.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListConfigsOptions.Version property.
// The version of configuration to return.
const (
	ListConfigsOptions_Version_Active = "active"
	ListConfigsOptions_Version_Draft = "draft"
	ListConfigsOptions_Version_Mixed = "mixed"
)

// NewListConfigsOptions : Instantiate ListConfigsOptions
func (*ProjectV1) NewListConfigsOptions(projectID string) *ListConfigsOptions {
	return &ListConfigsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListConfigsOptions) SetProjectID(projectID string) *ListConfigsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ListConfigsOptions) SetVersion(version string) *ListConfigsOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigsOptions) SetHeaders(param map[string]string) *ListConfigsOptions {
	options.Headers = param
	return options
}

// ListProjectsOptions : The ListProjects options.
type ListProjectsOptions struct {
	// Marks the last entry that is returned on the page. The server uses this parameter to determine the first entry that
	// is returned on the next page. If this parameter is not specified, the logical first page is returned.
	Start *string `json:"start,omitempty"`

	// Determine the maximum number of resources to return. The number of resources that are returned is the same, with the
	// exception of the last page.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProjectsOptions : Instantiate ListProjectsOptions
func (*ProjectV1) NewListProjectsOptions() *ListProjectsOptions {
	return &ListProjectsOptions{}
}

// SetStart : Allow user to set Start
func (_options *ListProjectsOptions) SetStart(start string) *ListProjectsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProjectsOptions) SetLimit(limit int64) *ListProjectsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProjectsOptions) SetHeaders(param map[string]string) *ListProjectsOptions {
	options.Headers = param
	return options
}

// NotificationEvent : NotificationEvent struct
type NotificationEvent struct {
	// The type of event.
	Event *string `json:"event" validate:"required"`

	// The target of the event.
	Target *string `json:"target" validate:"required"`

	// The source of the event.
	Source *string `json:"source,omitempty"`

	// The user that triggered the flow that posted the event.
	TriggeredBy *string `json:"triggered_by,omitempty"`

	// An actionable URL that users can access in response to the event.
	ActionURL *string `json:"action_url,omitempty"`

	// Any relevant metadata to be stored.
	Data map[string]interface{} `json:"data,omitempty"`
}

// NewNotificationEvent : Instantiate NotificationEvent (Generic Model Constructor)
func (*ProjectV1) NewNotificationEvent(event string, target string) (_model *NotificationEvent, err error) {
	_model = &NotificationEvent{
		Event: core.StringPtr(event),
		Target: core.StringPtr(target),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalNotificationEvent unmarshals an instance of NotificationEvent from the specified map of raw messages.
func UnmarshalNotificationEvent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationEvent)
	err = core.UnmarshalPrimitive(m, "event", &obj.Event)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "triggered_by", &obj.TriggeredBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_url", &obj.ActionURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationEventWithID : NotificationEventWithID struct
type NotificationEventWithID struct {
	// The type of event.
	Event *string `json:"event" validate:"required"`

	// The target of the event.
	Target *string `json:"target" validate:"required"`

	// The source of the event.
	Source *string `json:"source,omitempty"`

	// The user that triggered the flow that posted the event.
	TriggeredBy *string `json:"triggered_by,omitempty"`

	// An actionable URL that users can access in response to the event.
	ActionURL *string `json:"action_url,omitempty"`

	// Any relevant metadata to be stored.
	Data map[string]interface{} `json:"data,omitempty"`

	// The unique ID of a project.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalNotificationEventWithID unmarshals an instance of NotificationEventWithID from the specified map of raw messages.
func UnmarshalNotificationEventWithID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationEventWithID)
	err = core.UnmarshalPrimitive(m, "event", &obj.Event)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "triggered_by", &obj.TriggeredBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_url", &obj.ActionURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
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

// NotificationEventWithStatus : NotificationEventWithStatus struct
type NotificationEventWithStatus struct {
	// The type of event.
	Event *string `json:"event" validate:"required"`

	// The target of the event.
	Target *string `json:"target" validate:"required"`

	// The source of the event.
	Source *string `json:"source,omitempty"`

	// The user that triggered the flow that posted the event.
	TriggeredBy *string `json:"triggered_by,omitempty"`

	// An actionable URL that users can access in response to the event.
	ActionURL *string `json:"action_url,omitempty"`

	// Any relevant metadata to be stored.
	Data map[string]interface{} `json:"data,omitempty"`

	// The unique ID of a project.
	ID *string `json:"id" validate:"required"`

	// Whether or not the event successfully posted.
	Status *string `json:"status,omitempty"`

	// The reasons for the status of an event.
	Reasons []map[string]interface{} `json:"reasons,omitempty"`
}

// UnmarshalNotificationEventWithStatus unmarshals an instance of NotificationEventWithStatus from the specified map of raw messages.
func UnmarshalNotificationEventWithStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationEventWithStatus)
	err = core.UnmarshalPrimitive(m, "event", &obj.Event)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "triggered_by", &obj.TriggeredBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_url", &obj.ActionURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
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
	err = core.UnmarshalPrimitive(m, "reasons", &obj.Reasons)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationsGetResponse : The response from fetching notifications.
type NotificationsGetResponse struct {
	// Collection of the notification events with an ID.
	Notifications []NotificationEventWithID `json:"notifications,omitempty"`
}

// UnmarshalNotificationsGetResponse unmarshals an instance of NotificationsGetResponse from the specified map of raw messages.
func UnmarshalNotificationsGetResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsGetResponse)
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalNotificationEventWithID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationsIntegrationGetResponse : The resulting response of getting the source details of the event notifications integration.
type NotificationsIntegrationGetResponse struct {
	// A description of the instance of the event.
	Description *string `json:"description,omitempty"`

	// The name of the instance of the event.
	Name *string `json:"name,omitempty"`

	// The status of the instance of the event.
	Enabled *bool `json:"enabled,omitempty"`

	// A unique ID of the instance of the event.
	ID *string `json:"id,omitempty"`

	// The type of the instance of event.
	Type *string `json:"type,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The topic count of the instance of the event.
	TopicCount *int64 `json:"topic_count,omitempty"`

	// The topic names of the instance of the event.
	TopicNames []string `json:"topic_names,omitempty"`
}

// UnmarshalNotificationsIntegrationGetResponse unmarshals an instance of NotificationsIntegrationGetResponse from the specified map of raw messages.
func UnmarshalNotificationsIntegrationGetResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsIntegrationGetResponse)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
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
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "topic_count", &obj.TopicCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "topic_names", &obj.TopicNames)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationsIntegrationPostResponse : The resulting response of connecting a project to an event notifications instance.
type NotificationsIntegrationPostResponse struct {
	// A description of the instance of the event.
	Description *string `json:"description,omitempty"`

	// A name of the instance of the event.
	Name *string `json:"name,omitempty"`

	// A status of the instance of the event.
	Enabled *bool `json:"enabled,omitempty"`

	// A unique ID of the instance of the event.
	ID *string `json:"id,omitempty"`

	// The type of instance of the event.
	Type *string `json:"type,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`
}

// UnmarshalNotificationsIntegrationPostResponse unmarshals an instance of NotificationsIntegrationPostResponse from the specified map of raw messages.
func UnmarshalNotificationsIntegrationPostResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsIntegrationPostResponse)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
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
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NotificationsIntegrationTestPostResponse : The response for posting a test notification to the event notifications instance.
type NotificationsIntegrationTestPostResponse struct {
	// The data content type of the instance of the event.
	Datacontenttype *string `json:"datacontenttype,omitempty"`

	// The IBM default long message for the instance of the event.
	Ibmendefaultlong *string `json:"ibmendefaultlong,omitempty"`

	// The IBM default short message for the instance of the event.
	Ibmendefaultshort *string `json:"ibmendefaultshort,omitempty"`

	// The IBM source ID for the instance of the event.
	Ibmensourceid *string `json:"ibmensourceid,omitempty"`

	// A unique ID of the instance of the event.
	ID *string `json:"id" validate:"required"`

	// The source of the instance of the event.
	Source *string `json:"source" validate:"required"`

	// The spec version of the instance of the event.
	Specversion *string `json:"specversion,omitempty"`

	// The type of instance of the event.
	Type *string `json:"type,omitempty"`
}

// UnmarshalNotificationsIntegrationTestPostResponse unmarshals an instance of NotificationsIntegrationTestPostResponse from the specified map of raw messages.
func UnmarshalNotificationsIntegrationTestPostResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsIntegrationTestPostResponse)
	err = core.UnmarshalPrimitive(m, "datacontenttype", &obj.Datacontenttype)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmendefaultlong", &obj.Ibmendefaultlong)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmendefaultshort", &obj.Ibmendefaultshort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmensourceid", &obj.Ibmensourceid)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "specversion", &obj.Specversion)
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

// NotificationsPrototypePostResponse : The result of a notification post.
type NotificationsPrototypePostResponse struct {
	// The collection of the notification events with status.
	Notifications []NotificationEventWithStatus `json:"notifications,omitempty"`
}

// UnmarshalNotificationsPrototypePostResponse unmarshals an instance of NotificationsPrototypePostResponse from the specified map of raw messages.
func UnmarshalNotificationsPrototypePostResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NotificationsPrototypePostResponse)
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalNotificationEventWithStatus)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OutputValue : OutputValue struct
type OutputValue struct {
	// The variable name.
	Name *string `json:"name" validate:"required"`

	// A short explanation of the output value.
	Description *string `json:"description,omitempty"`

	// Can be any value - a string, number, boolean, array, or object.
	Value interface{} `json:"value,omitempty"`
}

// UnmarshalOutputValue unmarshals an instance of OutputValue from the specified map of raw messages.
func UnmarshalOutputValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OutputValue)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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

// PaginationLink : A pagination link.
type PaginationLink struct {
	// The URL of the pull request, which uniquely identifies it.
	Href *string `json:"href" validate:"required"`

	// A pagination token.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPaginationLink unmarshals an instance of PaginationLink from the specified map of raw messages.
func UnmarshalPaginationLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginationLink)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostCrnTokenOptions : The PostCrnToken options.
type PostCrnTokenOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostCrnTokenOptions : Instantiate PostCrnTokenOptions
func (*ProjectV1) NewPostCrnTokenOptions(id string) *PostCrnTokenOptions {
	return &PostCrnTokenOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *PostCrnTokenOptions) SetID(id string) *PostCrnTokenOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostCrnTokenOptions) SetHeaders(param map[string]string) *PostCrnTokenOptions {
	options.Headers = param
	return options
}

// PostEventNotificationsIntegrationOptions : The PostEventNotificationsIntegration options.
type PostEventNotificationsIntegrationOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// A CRN of the instance of the event.
	InstanceCrn *string `json:"instance_crn" validate:"required"`

	// A description of the instance of the event.
	Description *string `json:"description,omitempty"`

	// The name of the project source for the event notifications instance.
	EventNotificationsSourceName *string `json:"event_notifications_source_name,omitempty"`

	// A status of the instance of the event.
	Enabled *bool `json:"enabled,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostEventNotificationsIntegrationOptions : Instantiate PostEventNotificationsIntegrationOptions
func (*ProjectV1) NewPostEventNotificationsIntegrationOptions(id string, instanceCrn string) *PostEventNotificationsIntegrationOptions {
	return &PostEventNotificationsIntegrationOptions{
		ID: core.StringPtr(id),
		InstanceCrn: core.StringPtr(instanceCrn),
	}
}

// SetID : Allow user to set ID
func (_options *PostEventNotificationsIntegrationOptions) SetID(id string) *PostEventNotificationsIntegrationOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetInstanceCrn : Allow user to set InstanceCrn
func (_options *PostEventNotificationsIntegrationOptions) SetInstanceCrn(instanceCrn string) *PostEventNotificationsIntegrationOptions {
	_options.InstanceCrn = core.StringPtr(instanceCrn)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *PostEventNotificationsIntegrationOptions) SetDescription(description string) *PostEventNotificationsIntegrationOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEventNotificationsSourceName : Allow user to set EventNotificationsSourceName
func (_options *PostEventNotificationsIntegrationOptions) SetEventNotificationsSourceName(eventNotificationsSourceName string) *PostEventNotificationsIntegrationOptions {
	_options.EventNotificationsSourceName = core.StringPtr(eventNotificationsSourceName)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *PostEventNotificationsIntegrationOptions) SetEnabled(enabled bool) *PostEventNotificationsIntegrationOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostEventNotificationsIntegrationOptions) SetHeaders(param map[string]string) *PostEventNotificationsIntegrationOptions {
	options.Headers = param
	return options
}

// PostNotificationOptions : The PostNotification options.
type PostNotificationOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// Collection of the notification events to post.
	Notifications []NotificationEvent `json:"notifications,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostNotificationOptions : Instantiate PostNotificationOptions
func (*ProjectV1) NewPostNotificationOptions(id string) *PostNotificationOptions {
	return &PostNotificationOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *PostNotificationOptions) SetID(id string) *PostNotificationOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetNotifications : Allow user to set Notifications
func (_options *PostNotificationOptions) SetNotifications(notifications []NotificationEvent) *PostNotificationOptions {
	_options.Notifications = notifications
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostNotificationOptions) SetHeaders(param map[string]string) *PostNotificationOptions {
	options.Headers = param
	return options
}

// PostTestEventNotificationOptions : The PostTestEventNotification options.
type PostTestEventNotificationOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// The IBM default long message for the instance of an event.
	Ibmendefaultlong *string `json:"ibmendefaultlong,omitempty"`

	// The IBM default long message for the instance of an event.
	Ibmendefaultshort *string `json:"ibmendefaultshort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostTestEventNotificationOptions : Instantiate PostTestEventNotificationOptions
func (*ProjectV1) NewPostTestEventNotificationOptions(id string) *PostTestEventNotificationOptions {
	return &PostTestEventNotificationOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *PostTestEventNotificationOptions) SetID(id string) *PostTestEventNotificationOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIbmendefaultlong : Allow user to set Ibmendefaultlong
func (_options *PostTestEventNotificationOptions) SetIbmendefaultlong(ibmendefaultlong string) *PostTestEventNotificationOptions {
	_options.Ibmendefaultlong = core.StringPtr(ibmendefaultlong)
	return _options
}

// SetIbmendefaultshort : Allow user to set Ibmendefaultshort
func (_options *PostTestEventNotificationOptions) SetIbmendefaultshort(ibmendefaultshort string) *PostTestEventNotificationOptions {
	_options.Ibmendefaultshort = core.StringPtr(ibmendefaultshort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostTestEventNotificationOptions) SetHeaders(param map[string]string) *PostTestEventNotificationOptions {
	options.Headers = param
	return options
}

// Project : The project returned in the response body.
type Project struct {
	// The project name.
	Name *string `json:"name" validate:"required"`

	// A project descriptive text.
	Description *string `json:"description,omitempty"`

	// The unique ID of a project.
	ID *string `json:"id,omitempty"`

	// An IBM Cloud resource name, which uniquely identifies a resource.
	Crn *string `json:"crn,omitempty"`

	// The project configurations.
	Configs []ProjectConfig `json:"configs,omitempty"`

	// The metadata of the project.
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
}

// UnmarshalProject unmarshals an instance of Project from the specified map of raw messages.
func UnmarshalProject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Project)
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
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalProjectMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectCRNTokenResponse : The project CRN token.
type ProjectCRNTokenResponse struct {
	// The IAM access token.
	AccesToken *string `json:"acces_token,omitempty"`

	// Number of seconds counted since January 1st, 1970, until the IAM access token will expire.
	Expiration *int64 `json:"expiration,omitempty"`
}

// UnmarshalProjectCRNTokenResponse unmarshals an instance of ProjectCRNTokenResponse from the specified map of raw messages.
func UnmarshalProjectCRNTokenResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectCRNTokenResponse)
	err = core.UnmarshalPrimitive(m, "acces_token", &obj.AccesToken)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration", &obj.Expiration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectCollection : Projects list.
type ProjectCollection struct {
	// A pagination limit.
	Limit *int64 `json:"limit" validate:"required"`

	// Get the occurrencies of the total projects.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A pagination link.
	First *PaginationLink `json:"first" validate:"required"`

	// A pagination link.
	Last *PaginationLink `json:"last,omitempty"`

	// A pagination link.
	Previous *PaginationLink `json:"previous,omitempty"`

	// A pagination link.
	Next *PaginationLink `json:"next,omitempty"`

	// An array of projects.
	Projects []ProjectCollectionMemberWithMetadata `json:"projects,omitempty"`
}

// UnmarshalProjectCollection unmarshals an instance of ProjectCollection from the specified map of raw messages.
func UnmarshalProjectCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "projects", &obj.Projects, UnmarshalProjectCollectionMemberWithMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProjectCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ProjectCollectionMemberWithMetadata : ProjectCollectionMemberWithMetadata struct
type ProjectCollectionMemberWithMetadata struct {
	// The unique ID of a project.
	ID *string `json:"id,omitempty"`

	// The project name.
	Name *string `json:"name,omitempty"`

	// The project description.
	Description *string `json:"description,omitempty"`

	// The metadata of the project.
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
}

// UnmarshalProjectCollectionMemberWithMetadata unmarshals an instance of ProjectCollectionMemberWithMetadata from the specified map of raw messages.
func UnmarshalProjectCollectionMemberWithMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectCollectionMemberWithMetadata)
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
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalProjectMetadata)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfig : The project configuration.
type ProjectConfig struct {
	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id,omitempty"`

	// The configuration name.
	Name *string `json:"name" validate:"required"`

	// A collection of configuration labels.
	Labels []string `json:"labels,omitempty"`

	// The project configuration description.
	Description *string `json:"description,omitempty"`

	// A dotted value of catalogID.versionID.
	LocatorID *string `json:"locator_id" validate:"required"`

	// The type of a project configuration manual property.
	Type *string `json:"type" validate:"required"`

	// The outputs of a Schematics template property.
	Input []InputVariable `json:"input,omitempty"`

	// The outputs of a Schematics template property.
	Output []OutputValue `json:"output,omitempty"`

	// Schematics environment variables to use to deploy the configuration.
	Setting []ProjectConfigSettingCollection `json:"setting,omitempty"`
}

// Constants associated with the ProjectConfig.Type property.
// The type of a project configuration manual property.
const (
	ProjectConfig_Type_SchematicsBlueprint = "schematics_blueprint"
	ProjectConfig_Type_TerraformTemplate = "terraform_template"
)

// UnmarshalProjectConfig unmarshals an instance of ProjectConfig from the specified map of raw messages.
func UnmarshalProjectConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfig)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalInputVariable)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "output", &obj.Output, UnmarshalOutputValue)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "setting", &obj.Setting, UnmarshalProjectConfigSettingCollection)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigCollection : The project configuration list.
type ProjectConfigCollection struct {
	// The collection list operation response schema that should define the array property with the name "configs".
	Configs []ProjectConfig `json:"configs,omitempty"`
}

// UnmarshalProjectConfigCollection unmarshals an instance of ProjectConfigCollection from the specified map of raw messages.
func UnmarshalProjectConfigCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigCollection)
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDelete : Deletes the configuration response.
type ProjectConfigDelete struct {
	// The unique ID of a project.
	ID *string `json:"id,omitempty"`

	// The name of the configuration being deleted.
	Name *string `json:"name,omitempty"`
}

// UnmarshalProjectConfigDelete unmarshals an instance of ProjectConfigDelete from the specified map of raw messages.
func UnmarshalProjectConfigDelete(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDelete)
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

// ProjectConfigDiff : The project configuration diff summary.
type ProjectConfigDiff struct {
	// The additions to configurations in the diff summary.
	Added *ProjectConfigDiffAdded `json:"added,omitempty"`

	// The changes to configurations in the diff summary.
	Changed *ProjectConfigDiffChanged `json:"changed,omitempty"`

	// The deletions to configurations in the diff summary.
	Removed *ProjectConfigDiffRemoved `json:"removed,omitempty"`
}

// UnmarshalProjectConfigDiff unmarshals an instance of ProjectConfigDiff from the specified map of raw messages.
func UnmarshalProjectConfigDiff(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDiff)
	err = core.UnmarshalModel(m, "added", &obj.Added, UnmarshalProjectConfigDiffAdded)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "changed", &obj.Changed, UnmarshalProjectConfigDiffChanged)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "removed", &obj.Removed, UnmarshalProjectConfigDiffRemoved)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDiffAdded : The additions to configurations in the diff summary.
type ProjectConfigDiffAdded struct {
	// The collection of additions to the configurations in the diff summary.
	Input []ProjectConfigDiffInputVariable `json:"input,omitempty"`
}

// UnmarshalProjectConfigDiffAdded unmarshals an instance of ProjectConfigDiffAdded from the specified map of raw messages.
func UnmarshalProjectConfigDiffAdded(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDiffAdded)
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigDiffInputVariable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDiffChanged : The changes to configurations in the diff summary.
type ProjectConfigDiffChanged struct {
	// The collection of changes to configurations in the diff summary.
	Input []ProjectConfigDiffInputVariable `json:"input,omitempty"`
}

// UnmarshalProjectConfigDiffChanged unmarshals an instance of ProjectConfigDiffChanged from the specified map of raw messages.
func UnmarshalProjectConfigDiffChanged(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDiffChanged)
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigDiffInputVariable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDiffInputVariable : ProjectConfigDiffInputVariable struct
type ProjectConfigDiffInputVariable struct {
	// The variable name.
	Name *string `json:"name" validate:"required"`

	// The variable type.
	Type *string `json:"type" validate:"required"`

	// Can be any value - a string, number, boolean, array, or object.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the ProjectConfigDiffInputVariable.Type property.
// The variable type.
const (
	ProjectConfigDiffInputVariable_Type_Array = "array"
	ProjectConfigDiffInputVariable_Type_Boolean = "boolean"
	ProjectConfigDiffInputVariable_Type_Float = "float"
	ProjectConfigDiffInputVariable_Type_Int = "int"
	ProjectConfigDiffInputVariable_Type_Number = "number"
	ProjectConfigDiffInputVariable_Type_Object = "object"
	ProjectConfigDiffInputVariable_Type_Password = "password"
	ProjectConfigDiffInputVariable_Type_String = "string"
)

// UnmarshalProjectConfigDiffInputVariable unmarshals an instance of ProjectConfigDiffInputVariable from the specified map of raw messages.
func UnmarshalProjectConfigDiffInputVariable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDiffInputVariable)
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

// ProjectConfigDiffRemoved : The deletions to configurations in the diff summary.
type ProjectConfigDiffRemoved struct {
	// The collection of deletions to configurations in the diff summary.
	Input []ProjectConfigDiffInputVariable `json:"input,omitempty"`
}

// UnmarshalProjectConfigDiffRemoved unmarshals an instance of ProjectConfigDiffRemoved from the specified map of raw messages.
func UnmarshalProjectConfigDiffRemoved(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDiffRemoved)
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigDiffInputVariable)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigInputVariable : ProjectConfigInputVariable struct
type ProjectConfigInputVariable struct {
	// The variable name.
	Name *string `json:"name" validate:"required"`

	// Can be any value - a string, number, boolean, array, or object.
	Value interface{} `json:"value,omitempty"`
}

// NewProjectConfigInputVariable : Instantiate ProjectConfigInputVariable (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigInputVariable(name string) (_model *ProjectConfigInputVariable, err error) {
	_model = &ProjectConfigInputVariable{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProjectConfigInputVariable unmarshals an instance of ProjectConfigInputVariable from the specified map of raw messages.
func UnmarshalProjectConfigInputVariable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigInputVariable)
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

// ProjectConfigPatchRequest : The project configuration input.
// Models which "extend" this model:
// - ProjectConfigPatchRequestProjectConfigManualProperty
// - ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate
type ProjectConfigPatchRequest struct {
	// The configuration name.
	Name *string `json:"name,omitempty"`

	// The configuration labels.
	Labels []string `json:"labels,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The type of a project configuration manual property.
	Type *string `json:"type,omitempty"`

	// The external resource account ID in project configuration.
	ExternalResourcesAccount *string `json:"external_resources_account,omitempty"`

	// A dotted value of catalogID.versionID.
	LocatorID *string `json:"locator_id,omitempty"`

	// The inputs of a Schematics template property.
	Input []ProjectConfigInputVariable `json:"input,omitempty"`

	// Schematics environment variables to use to deploy the configuration.
	Setting []ProjectConfigSettingCollection `json:"setting,omitempty"`
}

// Constants associated with the ProjectConfigPatchRequest.Type property.
// The type of a project configuration manual property.
const (
	ProjectConfigPatchRequest_Type_Manual = "manual"
)
func (*ProjectConfigPatchRequest) isaProjectConfigPatchRequest() bool {
	return true
}

type ProjectConfigPatchRequestIntf interface {
	isaProjectConfigPatchRequest() bool
}

// UnmarshalProjectConfigPatchRequest unmarshals an instance of ProjectConfigPatchRequest from the specified map of raw messages.
func UnmarshalProjectConfigPatchRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchRequest)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external_resources_account", &obj.ExternalResourcesAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigInputVariable)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "setting", &obj.Setting, UnmarshalProjectConfigSettingCollection)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototype : The input of a project configuration.
type ProjectConfigPrototype struct {
	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id,omitempty"`

	// The configuration name.
	Name *string `json:"name" validate:"required"`

	// A collection of configuration labels.
	Labels []string `json:"labels,omitempty"`

	// The project configuration description.
	Description *string `json:"description,omitempty"`

	// A dotted value of catalogID.versionID.
	LocatorID *string `json:"locator_id" validate:"required"`

	// The input values to use to deploy the configuration.
	Input []ProjectConfigInputVariable `json:"input,omitempty"`

	// Schematics environment variables to use to deploy the configuration.
	Setting []ProjectConfigSettingCollection `json:"setting,omitempty"`
}

// NewProjectConfigPrototype : Instantiate ProjectConfigPrototype (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigPrototype(name string, locatorID string) (_model *ProjectConfigPrototype, err error) {
	_model = &ProjectConfigPrototype{
		Name: core.StringPtr(name),
		LocatorID: core.StringPtr(locatorID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProjectConfigPrototype unmarshals an instance of ProjectConfigPrototype from the specified map of raw messages.
func UnmarshalProjectConfigPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototype)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigInputVariable)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "setting", &obj.Setting, UnmarshalProjectConfigSettingCollection)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigSettingCollection : ProjectConfigSettingCollection struct
type ProjectConfigSettingCollection struct {
	// The name of the configuration setting.
	Name *string `json:"name" validate:"required"`

	// The value of the configuration setting.
	Value *string `json:"value" validate:"required"`
}

// NewProjectConfigSettingCollection : Instantiate ProjectConfigSettingCollection (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigSettingCollection(name string, value string) (_model *ProjectConfigSettingCollection, err error) {
	_model = &ProjectConfigSettingCollection{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProjectConfigSettingCollection unmarshals an instance of ProjectConfigSettingCollection from the specified map of raw messages.
func UnmarshalProjectConfigSettingCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigSettingCollection)
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

// ProjectMetadata : The metadata of the project.
type ProjectMetadata struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	Crn *string `json:"crn,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The cumulative list of needs attention items for a project.
	CumulativeNeedsAttentionView []CumulativeNeedsAttention `json:"cumulative_needs_attention_view,omitempty"`

	// "True" indicates that the fetch of the needs attention items failed.
	CumulativeNeedsAttentionViewErr *string `json:"cumulative_needs_attention_view_err,omitempty"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location,omitempty"`

	// The resource group where the project's data and tools are created.
	ResourceGroup *string `json:"resource_group,omitempty"`

	// The project status value.
	State *string `json:"state,omitempty"`

	// The CRN of the event notifications instance if one is connected to this project.
	EventNotificationsCrn *string `json:"event_notifications_crn,omitempty"`
}

// UnmarshalProjectMetadata unmarshals an instance of ProjectMetadata from the specified map of raw messages.
func UnmarshalProjectMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectMetadata)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cumulative_needs_attention_view", &obj.CumulativeNeedsAttentionView, UnmarshalCumulativeNeedsAttention)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cumulative_needs_attention_view_err", &obj.CumulativeNeedsAttentionViewErr)
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
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_notifications_crn", &obj.EventNotificationsCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UninstallConfigOptions : The UninstallConfig options.
type UninstallConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUninstallConfigOptions : Instantiate UninstallConfigOptions
func (*ProjectV1) NewUninstallConfigOptions(projectID string, id string) *UninstallConfigOptions {
	return &UninstallConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UninstallConfigOptions) SetProjectID(projectID string) *UninstallConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UninstallConfigOptions) SetID(id string) *UninstallConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UninstallConfigOptions) SetHeaders(param map[string]string) *UninstallConfigOptions {
	options.Headers = param
	return options
}

// UpdateConfigOptions : The UpdateConfig options.
type UpdateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The change delta of the project configuration to update.
	ProjectConfig ProjectConfigPatchRequestIntf `json:"project_config" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateConfigOptions : Instantiate UpdateConfigOptions
func (*ProjectV1) NewUpdateConfigOptions(projectID string, id string, projectConfig ProjectConfigPatchRequestIntf) *UpdateConfigOptions {
	return &UpdateConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		ProjectConfig: projectConfig,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateConfigOptions) SetProjectID(projectID string) *UpdateConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateConfigOptions) SetID(id string) *UpdateConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetProjectConfig : Allow user to set ProjectConfig
func (_options *UpdateConfigOptions) SetProjectConfig(projectConfig ProjectConfigPatchRequestIntf) *UpdateConfigOptions {
	_options.ProjectConfig = projectConfig
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateConfigOptions) SetHeaders(param map[string]string) *UpdateConfigOptions {
	options.Headers = param
	return options
}

// UpdateProjectOptions : The UpdateProject options.
type UpdateProjectOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// The project name.
	Name *string `json:"name,omitempty"`

	// The project descriptive text.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProjectOptions : Instantiate UpdateProjectOptions
func (*ProjectV1) NewUpdateProjectOptions(id string) *UpdateProjectOptions {
	return &UpdateProjectOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateProjectOptions) SetID(id string) *UpdateProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProjectOptions) SetName(name string) *UpdateProjectOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateProjectOptions) SetDescription(description string) *UpdateProjectOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProjectOptions) SetHeaders(param map[string]string) *UpdateProjectOptions {
	options.Headers = param
	return options
}

// ProjectConfigPatchRequestProjectConfigManualProperty : The project configuration manual type.
// This model "extends" ProjectConfigPatchRequest
type ProjectConfigPatchRequestProjectConfigManualProperty struct {
	// The configuration name.
	Name *string `json:"name,omitempty"`

	// The configuration labels.
	Labels []string `json:"labels,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The type of a project configuration manual property.
	Type *string `json:"type" validate:"required"`

	// The external resource account ID in project configuration.
	ExternalResourcesAccount *string `json:"external_resources_account,omitempty"`
}

// Constants associated with the ProjectConfigPatchRequestProjectConfigManualProperty.Type property.
// The type of a project configuration manual property.
const (
	ProjectConfigPatchRequestProjectConfigManualProperty_Type_Manual = "manual"
)

// NewProjectConfigPatchRequestProjectConfigManualProperty : Instantiate ProjectConfigPatchRequestProjectConfigManualProperty (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigPatchRequestProjectConfigManualProperty(typeVar string) (_model *ProjectConfigPatchRequestProjectConfigManualProperty, err error) {
	_model = &ProjectConfigPatchRequestProjectConfigManualProperty{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ProjectConfigPatchRequestProjectConfigManualProperty) isaProjectConfigPatchRequest() bool {
	return true
}

// UnmarshalProjectConfigPatchRequestProjectConfigManualProperty unmarshals an instance of ProjectConfigPatchRequestProjectConfigManualProperty from the specified map of raw messages.
func UnmarshalProjectConfigPatchRequestProjectConfigManualProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchRequestProjectConfigManualProperty)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "external_resources_account", &obj.ExternalResourcesAccount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate : The Schematics template property.
// This model "extends" ProjectConfigPatchRequest
type ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate struct {
	// The configuration name.
	Name *string `json:"name,omitempty"`

	// The configuration labels.
	Labels []string `json:"labels,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// A dotted value of catalogID.versionID.
	LocatorID *string `json:"locator_id,omitempty"`

	// The inputs of a Schematics template property.
	Input []ProjectConfigInputVariable `json:"input,omitempty"`

	// Schematics environment variables to use to deploy the configuration.
	Setting []ProjectConfigSettingCollection `json:"setting,omitempty"`
}

func (*ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate) isaProjectConfigPatchRequest() bool {
	return true
}

// UnmarshalProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate unmarshals an instance of ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate from the specified map of raw messages.
func UnmarshalProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchRequestProjectConfigPatchSchematicsTemplate)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "input", &obj.Input, UnmarshalProjectConfigInputVariable)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "setting", &obj.Setting, UnmarshalProjectConfigSettingCollection)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

//
// ProjectsPager can be used to simplify the use of the "ListProjects" method.
//
type ProjectsPager struct {
	hasNext bool
	options *ListProjectsOptions
	client  *ProjectV1
	pageContext struct {
		next *string
	}
}

// NewProjectsPager returns a new ProjectsPager instance.
func (project *ProjectV1) NewProjectsPager(options *ListProjectsOptions) (pager *ProjectsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListProjectsOptions = *options
	pager = &ProjectsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  project,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ProjectsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ProjectsPager) GetNextWithContext(ctx context.Context) (page []ProjectCollectionMemberWithMetadata, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListProjectsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Projects

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProjectsPager) GetAllWithContext(ctx context.Context) (allItems []ProjectCollectionMemberWithMetadata, err error) {
	for pager.HasNext() {
		var nextPage []ProjectCollectionMemberWithMetadata
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetNext() (page []ProjectCollectionMemberWithMetadata, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetAll() (allItems []ProjectCollectionMemberWithMetadata, err error) {
	return pager.GetAllWithContext(context.Background())
}
