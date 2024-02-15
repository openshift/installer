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
 * IBM OpenAPI SDK Code Generator Version: 3.83.0-adaf0721-20231212-210453
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

// ProjectV1 : Manage infrastructure as code in IBM Cloud.
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

	body := make(map[string]interface{})
	if createProjectOptions.Definition != nil {
		body["definition"] = createProjectOptions.Definition
	}
	if createProjectOptions.Location != nil {
		body["location"] = createProjectOptions.Location
	}
	if createProjectOptions.ResourceGroup != nil {
		body["resource_group"] = createProjectOptions.ResourceGroup
	}
	if createProjectOptions.Configs != nil {
		body["configs"] = createProjectOptions.Configs
	}
	if createProjectOptions.Environments != nil {
		body["environments"] = createProjectOptions.Environments
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
	if updateProjectOptions.Definition != nil {
		body["definition"] = updateProjectOptions.Definition
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
// Delete a project document by the ID. A project can only be deleted after deleting all of its resources.
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = project.Service.Request(request, nil)

	return
}

// CreateProjectEnvironment : Create an environment
// Create an environment.
func (project *ProjectV1) CreateProjectEnvironment(createProjectEnvironmentOptions *CreateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return project.CreateProjectEnvironmentWithContext(context.Background(), createProjectEnvironmentOptions)
}

// CreateProjectEnvironmentWithContext is an alternate form of the CreateProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) CreateProjectEnvironmentWithContext(ctx context.Context, createProjectEnvironmentOptions *CreateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProjectEnvironmentOptions, "createProjectEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProjectEnvironmentOptions, "createProjectEnvironmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createProjectEnvironmentOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/environments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProjectEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "CreateProjectEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProjectEnvironmentOptions.Definition != nil {
		body["definition"] = createProjectEnvironmentOptions.Definition
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProjectEnvironments : List environments
// Returns all environments.
func (project *ProjectV1) ListProjectEnvironments(listProjectEnvironmentsOptions *ListProjectEnvironmentsOptions) (result *EnvironmentCollection, response *core.DetailedResponse, err error) {
	return project.ListProjectEnvironmentsWithContext(context.Background(), listProjectEnvironmentsOptions)
}

// ListProjectEnvironmentsWithContext is an alternate form of the ListProjectEnvironments method which supports a Context parameter
func (project *ProjectV1) ListProjectEnvironmentsWithContext(ctx context.Context, listProjectEnvironmentsOptions *ListProjectEnvironmentsOptions) (result *EnvironmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProjectEnvironmentsOptions, "listProjectEnvironmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProjectEnvironmentsOptions, "listProjectEnvironmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listProjectEnvironmentsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/environments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProjectEnvironmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ListProjectEnvironments")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProjectEnvironment : Get an environment
// Returns an environment.
func (project *ProjectV1) GetProjectEnvironment(getProjectEnvironmentOptions *GetProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return project.GetProjectEnvironmentWithContext(context.Background(), getProjectEnvironmentOptions)
}

// GetProjectEnvironmentWithContext is an alternate form of the GetProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) GetProjectEnvironmentWithContext(ctx context.Context, getProjectEnvironmentOptions *GetProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectEnvironmentOptions, "getProjectEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProjectEnvironmentOptions, "getProjectEnvironmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getProjectEnvironmentOptions.ProjectID,
		"id": *getProjectEnvironmentOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/environments/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProjectEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetProjectEnvironment")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateProjectEnvironment : Update an environment
// Update an environment by the ID.
func (project *ProjectV1) UpdateProjectEnvironment(updateProjectEnvironmentOptions *UpdateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return project.UpdateProjectEnvironmentWithContext(context.Background(), updateProjectEnvironmentOptions)
}

// UpdateProjectEnvironmentWithContext is an alternate form of the UpdateProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) UpdateProjectEnvironmentWithContext(ctx context.Context, updateProjectEnvironmentOptions *UpdateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProjectEnvironmentOptions, "updateProjectEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProjectEnvironmentOptions, "updateProjectEnvironmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateProjectEnvironmentOptions.ProjectID,
		"id": *updateProjectEnvironmentOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/environments/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateProjectEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UpdateProjectEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateProjectEnvironmentOptions.Definition != nil {
		body["definition"] = updateProjectEnvironmentOptions.Definition
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProjectEnvironment : Delete an environment
// Delete an environment in a project by ID.
func (project *ProjectV1) DeleteProjectEnvironment(deleteProjectEnvironmentOptions *DeleteProjectEnvironmentOptions) (result *EnvironmentDeleteResponse, response *core.DetailedResponse, err error) {
	return project.DeleteProjectEnvironmentWithContext(context.Background(), deleteProjectEnvironmentOptions)
}

// DeleteProjectEnvironmentWithContext is an alternate form of the DeleteProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) DeleteProjectEnvironmentWithContext(ctx context.Context, deleteProjectEnvironmentOptions *DeleteProjectEnvironmentOptions) (result *EnvironmentDeleteResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProjectEnvironmentOptions, "deleteProjectEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProjectEnvironmentOptions, "deleteProjectEnvironmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteProjectEnvironmentOptions.ProjectID,
		"id": *deleteProjectEnvironmentOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/environments/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProjectEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteProjectEnvironment")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentDeleteResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

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
	if createConfigOptions.Definition != nil {
		body["definition"] = createConfigOptions.Definition
	}
	if createConfigOptions.Schematics != nil {
		body["schematics"] = createConfigOptions.Schematics
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

	body := make(map[string]interface{})
	if updateConfigOptions.Definition != nil {
		body["definition"] = updateConfigOptions.Definition
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

// DeleteConfig : Delete a configuration in a project by ID
// Delete a configuration in a project.
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

// ForceApprove : Force approve project configuration
// Force approve configuration edits to the main configuration with an approving comment.
func (project *ProjectV1) ForceApprove(forceApproveOptions *ForceApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	return project.ForceApproveWithContext(context.Background(), forceApproveOptions)
}

// ForceApproveWithContext is an alternate form of the ForceApprove method which supports a Context parameter
func (project *ProjectV1) ForceApproveWithContext(ctx context.Context, forceApproveOptions *ForceApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Approve : Approve and merge a configuration draft
// Approve and merge configuration edits to the main configuration.
func (project *ProjectV1) Approve(approveOptions *ApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	return project.ApproveWithContext(context.Background(), approveOptions)
}

// ApproveWithContext is an alternate form of the Approve method which supports a Context parameter
func (project *ProjectV1) ApproveWithContext(ctx context.Context, approveOptions *ApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ValidateConfig : Run a validation check
// Run a validation check on a given configuration in project. The check includes creating or updating the associated
// schematics workspace with a plan job, running the CRA scans, and cost estimatation.
func (project *ProjectV1) ValidateConfig(validateConfigOptions *ValidateConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	return project.ValidateConfigWithContext(context.Background(), validateConfigOptions)
}

// ValidateConfigWithContext is an alternate form of the ValidateConfig method which supports a Context parameter
func (project *ProjectV1) ValidateConfigWithContext(ctx context.Context, validateConfigOptions *ValidateConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateConfigOptions, "validateConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(validateConfigOptions, "validateConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *validateConfigOptions.ProjectID,
		"id": *validateConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/validate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range validateConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ValidateConfig")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeployConfig : Deploy a configuration
// Deploy a project's configuration. It's an asynchronous operation that can be tracked using the get project
// configuration API with full metadata.
func (project *ProjectV1) DeployConfig(deployConfigOptions *DeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	return project.DeployConfigWithContext(context.Background(), deployConfigOptions)
}

// DeployConfigWithContext is an alternate form of the DeployConfig method which supports a Context parameter
func (project *ProjectV1) DeployConfigWithContext(ctx context.Context, deployConfigOptions *DeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deployConfigOptions, "deployConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deployConfigOptions, "deployConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deployConfigOptions.ProjectID,
		"id": *deployConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/deploy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deployConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeployConfig")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UndeployConfig : Undeploy configuration resources
// Undeploy a project's configuration resources. The operation undeploys all the resources that are deployed with the
// specific configuration. You can track it by using the get project configuration API with full metadata.
func (project *ProjectV1) UndeployConfig(undeployConfigOptions *UndeployConfigOptions) (response *core.DetailedResponse, err error) {
	return project.UndeployConfigWithContext(context.Background(), undeployConfigOptions)
}

// UndeployConfigWithContext is an alternate form of the UndeployConfig method which supports a Context parameter
func (project *ProjectV1) UndeployConfigWithContext(ctx context.Context, undeployConfigOptions *UndeployConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(undeployConfigOptions, "undeployConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(undeployConfigOptions, "undeployConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *undeployConfigOptions.ProjectID,
		"id": *undeployConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/undeploy`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range undeployConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UndeployConfig")
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

// SyncConfig : Sync a project configuration
// Sync a project configuration by analyzing the associated pipeline runs and schematics workspace logs to get the
// configuration back to a working state.
func (project *ProjectV1) SyncConfig(syncConfigOptions *SyncConfigOptions) (response *core.DetailedResponse, err error) {
	return project.SyncConfigWithContext(context.Background(), syncConfigOptions)
}

// SyncConfigWithContext is an alternate form of the SyncConfig method which supports a Context parameter
func (project *ProjectV1) SyncConfigWithContext(ctx context.Context, syncConfigOptions *SyncConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(syncConfigOptions, "syncConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(syncConfigOptions, "syncConfigOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *syncConfigOptions.ProjectID,
		"id": *syncConfigOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/sync`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range syncConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "SyncConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if syncConfigOptions.Schematics != nil {
		body["schematics"] = syncConfigOptions.Schematics
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = project.Service.Request(request, nil)

	return
}

// ListConfigResources : List the resources deployed by a configuration
// A list of resources deployed by a configuraton.
func (project *ProjectV1) ListConfigResources(listConfigResourcesOptions *ListConfigResourcesOptions) (result *ProjectConfigResourceCollection, response *core.DetailedResponse, err error) {
	return project.ListConfigResourcesWithContext(context.Background(), listConfigResourcesOptions)
}

// ListConfigResourcesWithContext is an alternate form of the ListConfigResources method which supports a Context parameter
func (project *ProjectV1) ListConfigResourcesWithContext(ctx context.Context, listConfigResourcesOptions *ListConfigResourcesOptions) (result *ProjectConfigResourceCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigResourcesOptions, "listConfigResourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listConfigResourcesOptions, "listConfigResourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listConfigResourcesOptions.ProjectID,
		"id": *listConfigResourcesOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listConfigResourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ListConfigResources")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigResourceCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListConfigVersions : Get a list of versions of a project configuration
// Returns a list of previous and current versions of a project configuration in a specific project.
func (project *ProjectV1) ListConfigVersions(listConfigVersionsOptions *ListConfigVersionsOptions) (result *ProjectConfigVersionSummaryCollection, response *core.DetailedResponse, err error) {
	return project.ListConfigVersionsWithContext(context.Background(), listConfigVersionsOptions)
}

// ListConfigVersionsWithContext is an alternate form of the ListConfigVersions method which supports a Context parameter
func (project *ProjectV1) ListConfigVersionsWithContext(ctx context.Context, listConfigVersionsOptions *ListConfigVersionsOptions) (result *ProjectConfigVersionSummaryCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigVersionsOptions, "listConfigVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listConfigVersionsOptions, "listConfigVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listConfigVersionsOptions.ProjectID,
		"id": *listConfigVersionsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listConfigVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ListConfigVersions")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersionSummaryCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetConfigVersion : Get a specific version of a project configuration
// Returns a specific version of a project configuration in a specific project.
func (project *ProjectV1) GetConfigVersion(getConfigVersionOptions *GetConfigVersionOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	return project.GetConfigVersionWithContext(context.Background(), getConfigVersionOptions)
}

// GetConfigVersionWithContext is an alternate form of the GetConfigVersion method which supports a Context parameter
func (project *ProjectV1) GetConfigVersionWithContext(ctx context.Context, getConfigVersionOptions *GetConfigVersionOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigVersionOptions, "getConfigVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getConfigVersionOptions, "getConfigVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getConfigVersionOptions.ProjectID,
		"id": *getConfigVersionOptions.ID,
		"version": fmt.Sprint(*getConfigVersionOptions.Version),
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/versions/{version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getConfigVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetConfigVersion")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfigVersion : Delete a configuration for the specified project ID and version
// Delete a configuration in a project.
func (project *ProjectV1) DeleteConfigVersion(deleteConfigVersionOptions *DeleteConfigVersionOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	return project.DeleteConfigVersionWithContext(context.Background(), deleteConfigVersionOptions)
}

// DeleteConfigVersionWithContext is an alternate form of the DeleteConfigVersion method which supports a Context parameter
func (project *ProjectV1) DeleteConfigVersionWithContext(ctx context.Context, deleteConfigVersionOptions *DeleteConfigVersionOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigVersionOptions, "deleteConfigVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteConfigVersionOptions, "deleteConfigVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteConfigVersionOptions.ProjectID,
		"id": *deleteConfigVersionOptions.ID,
		"version": fmt.Sprint(*deleteConfigVersionOptions.Version),
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/versions/{version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteConfigVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteConfigVersion")
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
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigDelete)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActionJobApplyMessagesSummary : The messages of apply jobs on the configuration.
type ActionJobApplyMessagesSummary struct {
	// The collection of error messages.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`

	// The collection of success messages.
	SucessMessage []TerraformLogAnalyzerSuccessMessage `json:"sucess_message,omitempty"`
}

// UnmarshalActionJobApplyMessagesSummary unmarshals an instance of ActionJobApplyMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobApplyMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobApplyMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sucess_message", &obj.SucessMessage, UnmarshalTerraformLogAnalyzerSuccessMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobApplySummary : The summary of the apply jobs on the configuration.
type ActionJobApplySummary struct {
	// The number of applied resources.
	Success *int64 `json:"success,omitempty"`

	// The number of failed resources.
	Failed *int64 `json:"failed,omitempty"`

	// The collection of successfully applied resources.
	SuccessResources []string `json:"success_resources,omitempty"`

	// The collection of failed applied resources.
	FailedResources []string `json:"failed_resources,omitempty"`
}

// UnmarshalActionJobApplySummary unmarshals an instance of ActionJobApplySummary from the specified map of raw messages.
func UnmarshalActionJobApplySummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobApplySummary)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success_resources", &obj.SuccessResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_resources", &obj.FailedResources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroyMessagesSummary : The messages of destroy jobs on the configuration.
type ActionJobDestroyMessagesSummary struct {
	// The collection of error messages.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`
}

// UnmarshalActionJobDestroyMessagesSummary unmarshals an instance of ActionJobDestroyMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobDestroyMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroyMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroySummary : The summary of the destroy jobs on the configuration.
type ActionJobDestroySummary struct {
	// The number of destroyed resources.
	Success *int64 `json:"success,omitempty"`

	// The number of failed resources.
	Failed *int64 `json:"failed,omitempty"`

	// The number of tainted resources.
	Tainted *int64 `json:"tainted,omitempty"`

	// The destroy resources results from the job.
	Resources *ActionJobDestroySummaryResources `json:"resources,omitempty"`
}

// UnmarshalActionJobDestroySummary unmarshals an instance of ActionJobDestroySummary from the specified map of raw messages.
func UnmarshalActionJobDestroySummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroySummary)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tainted", &obj.Tainted)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalActionJobDestroySummaryResources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroySummaryResources : The destroy resources results from the job.
type ActionJobDestroySummaryResources struct {
	// The collection of destroyed resources.
	Success []string `json:"success,omitempty"`

	// The collection of failed resources.
	Failed []string `json:"failed,omitempty"`

	// The collection of tainted resources.
	Tainted []string `json:"tainted,omitempty"`
}

// UnmarshalActionJobDestroySummaryResources unmarshals an instance of ActionJobDestroySummaryResources from the specified map of raw messages.
func UnmarshalActionJobDestroySummaryResources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroySummaryResources)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tainted", &obj.Tainted)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobMessageSummary : The message summaries of jobs on the configuration.
type ActionJobMessageSummary struct {
	// The number of info messages.
	Info *int64 `json:"info,omitempty"`

	// The number of debug messages.
	Debug *int64 `json:"debug,omitempty"`

	// The number of error messages.
	Error *int64 `json:"error,omitempty"`
}

// UnmarshalActionJobMessageSummary unmarshals an instance of ActionJobMessageSummary from the specified map of raw messages.
func UnmarshalActionJobMessageSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobMessageSummary)
	err = core.UnmarshalPrimitive(m, "info", &obj.Info)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "debug", &obj.Debug)
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

// ActionJobPlanMessagesSummary : The plan messages on the configuration.
type ActionJobPlanMessagesSummary struct {
	// The collection of error messages.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`

	// The collection of success messages.
	SucessMessage []string `json:"sucess_message,omitempty"`

	// The collection of update messages.
	UpdateMessage []string `json:"update_message,omitempty"`

	// The collection of destroy messages.
	DestroyMessage []string `json:"destroy_message,omitempty"`
}

// UnmarshalActionJobPlanMessagesSummary unmarshals an instance of ActionJobPlanMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobPlanMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobPlanMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sucess_message", &obj.SucessMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_message", &obj.UpdateMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_message", &obj.DestroyMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobPlanSummary : The summary of the plan jobs on the configuration.
type ActionJobPlanSummary struct {
	// The number of resources to be added.
	Add *int64 `json:"add,omitempty"`

	// The number of resources that failed during the plan job.
	Failed *int64 `json:"failed,omitempty"`

	// The number of resources to be updated.
	Update *int64 `json:"update,omitempty"`

	// The number of resources to be destroyed.
	Destroy *int64 `json:"destroy,omitempty"`

	// The collection of planned added resources.
	AddResources []string `json:"add_resources,omitempty"`

	// The collection of failed planned resources.
	FailedResources []string `json:"failed_resources,omitempty"`

	// The collection of planned updated resources.
	UpdatedResources []string `json:"updated_resources,omitempty"`

	// The collection of planned destroy resources.
	DestroyResources []string `json:"destroy_resources,omitempty"`
}

// UnmarshalActionJobPlanSummary unmarshals an instance of ActionJobPlanSummary from the specified map of raw messages.
func UnmarshalActionJobPlanSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobPlanSummary)
	err = core.UnmarshalPrimitive(m, "add", &obj.Add)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update", &obj.Update)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy", &obj.Destroy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "add_resources", &obj.AddResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_resources", &obj.FailedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_resources", &obj.UpdatedResources)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_resources", &obj.DestroyResources)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobSummary : The summaries of jobs that were performed on the configuration.
type ActionJobSummary struct {
	// The summary of the plan jobs on the configuration.
	PlanSummary *ActionJobPlanSummary `json:"plan_summary" validate:"required"`

	// The summary of the apply jobs on the configuration.
	ApplySummary *ActionJobApplySummary `json:"apply_summary" validate:"required"`

	// The summary of the destroy jobs on the configuration.
	DestroySummary *ActionJobDestroySummary `json:"destroy_summary" validate:"required"`

	// The message summaries of jobs on the configuration.
	MessageSummary *ActionJobMessageSummary `json:"message_summary" validate:"required"`

	// The plan messages on the configuration.
	PlanMessages *ActionJobPlanMessagesSummary `json:"plan_messages" validate:"required"`

	// The messages of apply jobs on the configuration.
	ApplyMessages *ActionJobApplyMessagesSummary `json:"apply_messages" validate:"required"`

	// The messages of destroy jobs on the configuration.
	DestroyMessages *ActionJobDestroyMessagesSummary `json:"destroy_messages" validate:"required"`
}

// UnmarshalActionJobSummary unmarshals an instance of ActionJobSummary from the specified map of raw messages.
func UnmarshalActionJobSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobSummary)
	err = core.UnmarshalModel(m, "plan_summary", &obj.PlanSummary, UnmarshalActionJobPlanSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "apply_summary", &obj.ApplySummary, UnmarshalActionJobApplySummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "destroy_summary", &obj.DestroySummary, UnmarshalActionJobDestroySummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "message_summary", &obj.MessageSummary, UnmarshalActionJobMessageSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "plan_messages", &obj.PlanMessages, UnmarshalActionJobPlanMessagesSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "apply_messages", &obj.ApplyMessages, UnmarshalActionJobApplyMessagesSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "destroy_messages", &obj.DestroyMessages, UnmarshalActionJobDestroyMessagesSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobWithIdAndSummary : A brief summary of an action.
type ActionJobWithIdAndSummary struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// The summaries of jobs that were performed on the configuration.
	Summary *ActionJobSummary `json:"summary" validate:"required"`
}

// UnmarshalActionJobWithIdAndSummary unmarshals an instance of ActionJobWithIdAndSummary from the specified map of raw messages.
func UnmarshalActionJobWithIdAndSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobWithIdAndSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalActionJobSummary)
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

	// Notes on the project draft action. If this is a forced approve on the draft configuration, a non-empty comment is
	// required.
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

// CodeRiskAnalyzerLogsSummary : The Code Risk Analyzer logs summary of the configuration.
type CodeRiskAnalyzerLogsSummary struct {
	// The total number of Code Risk Analyzer rules that were applied in the scan.
	Total *string `json:"total,omitempty"`

	// The number of Code Risk Analyzer rules that passed in the scan.
	Passed *string `json:"passed,omitempty"`

	// The number of Code Risk Analyzer rules that failed in the scan.
	Failed *string `json:"failed,omitempty"`

	// The number of Code Risk Analyzer rules that were skipped in the scan.
	Skipped *string `json:"skipped,omitempty"`
}

// UnmarshalCodeRiskAnalyzerLogsSummary unmarshals an instance of CodeRiskAnalyzerLogsSummary from the specified map of raw messages.
func UnmarshalCodeRiskAnalyzerLogsSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CodeRiskAnalyzerLogsSummary)
	err = core.UnmarshalPrimitive(m, "total", &obj.Total)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "passed", &obj.Passed)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateConfigOptions : The CreateConfig options.
type CreateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	Definition ProjectConfigPrototypeDefinitionBlockIntf `json:"definition" validate:"required"`

	// A Schematics workspace to use for deploying this configuration.
	// Either schematics.workspace_crn, definition.locator_id, or both must be specified.
	Schematics *SchematicsWorkspace `json:"schematics,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateConfigOptions : Instantiate CreateConfigOptions
func (*ProjectV1) NewCreateConfigOptions(projectID string, definition ProjectConfigPrototypeDefinitionBlockIntf) *CreateConfigOptions {
	return &CreateConfigOptions{
		ProjectID: core.StringPtr(projectID),
		Definition: definition,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateConfigOptions) SetProjectID(projectID string) *CreateConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetDefinition : Allow user to set Definition
func (_options *CreateConfigOptions) SetDefinition(definition ProjectConfigPrototypeDefinitionBlockIntf) *CreateConfigOptions {
	_options.Definition = definition
	return _options
}

// SetSchematics : Allow user to set Schematics
func (_options *CreateConfigOptions) SetSchematics(schematics *SchematicsWorkspace) *CreateConfigOptions {
	_options.Schematics = schematics
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigOptions) SetHeaders(param map[string]string) *CreateConfigOptions {
	options.Headers = param
	return options
}

// CreateProjectEnvironmentOptions : The CreateProjectEnvironment options.
type CreateProjectEnvironmentOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The environment definition.
	Definition *EnvironmentDefinitionRequiredProperties `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProjectEnvironmentOptions : Instantiate CreateProjectEnvironmentOptions
func (*ProjectV1) NewCreateProjectEnvironmentOptions(projectID string, definition *EnvironmentDefinitionRequiredProperties) *CreateProjectEnvironmentOptions {
	return &CreateProjectEnvironmentOptions{
		ProjectID: core.StringPtr(projectID),
		Definition: definition,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateProjectEnvironmentOptions) SetProjectID(projectID string) *CreateProjectEnvironmentOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetDefinition : Allow user to set Definition
func (_options *CreateProjectEnvironmentOptions) SetDefinition(definition *EnvironmentDefinitionRequiredProperties) *CreateProjectEnvironmentOptions {
	_options.Definition = definition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProjectEnvironmentOptions) SetHeaders(param map[string]string) *CreateProjectEnvironmentOptions {
	options.Headers = param
	return options
}

// CreateProjectOptions : The CreateProject options.
type CreateProjectOptions struct {
	// The definition of the project.
	Definition *ProjectPrototypeDefinition `json:"definition" validate:"required"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location" validate:"required"`

	// The resource group name where the project's data and tools are created.
	ResourceGroup *string `json:"resource_group" validate:"required"`

	// The project configurations. These configurations are only included in the response of creating a project if a
	// configs array is specified in the request payload.
	Configs []ProjectConfigPrototype `json:"configs,omitempty"`

	// The project environments. These environments are only included in the response of creating a project if a
	// environments array is specified in the request payload.
	Environments []EnvironmentPrototype `json:"environments,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProjectOptions : Instantiate CreateProjectOptions
func (*ProjectV1) NewCreateProjectOptions(definition *ProjectPrototypeDefinition, location string, resourceGroup string) *CreateProjectOptions {
	return &CreateProjectOptions{
		Definition: definition,
		Location: core.StringPtr(location),
		ResourceGroup: core.StringPtr(resourceGroup),
	}
}

// SetDefinition : Allow user to set Definition
func (_options *CreateProjectOptions) SetDefinition(definition *ProjectPrototypeDefinition) *CreateProjectOptions {
	_options.Definition = definition
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateProjectOptions) SetLocation(location string) *CreateProjectOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateProjectOptions) SetResourceGroup(resourceGroup string) *CreateProjectOptions {
	_options.ResourceGroup = core.StringPtr(resourceGroup)
	return _options
}

// SetConfigs : Allow user to set Configs
func (_options *CreateProjectOptions) SetConfigs(configs []ProjectConfigPrototype) *CreateProjectOptions {
	_options.Configs = configs
	return _options
}

// SetEnvironments : Allow user to set Environments
func (_options *CreateProjectOptions) SetEnvironments(environments []EnvironmentPrototype) *CreateProjectOptions {
	_options.Environments = environments
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

	// A unique ID for that individual event.
	EventID *string `json:"event_id,omitempty"`

	// A unique ID for the configuration.
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

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigOptions) SetHeaders(param map[string]string) *DeleteConfigOptions {
	options.Headers = param
	return options
}

// DeleteConfigVersionOptions : The DeleteConfigVersion options.
type DeleteConfigVersionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The configuration version.
	Version *int64 `json:"version" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteConfigVersionOptions : Instantiate DeleteConfigVersionOptions
func (*ProjectV1) NewDeleteConfigVersionOptions(projectID string, id string, version int64) *DeleteConfigVersionOptions {
	return &DeleteConfigVersionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Version: core.Int64Ptr(version),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteConfigVersionOptions) SetProjectID(projectID string) *DeleteConfigVersionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteConfigVersionOptions) SetID(id string) *DeleteConfigVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *DeleteConfigVersionOptions) SetVersion(version int64) *DeleteConfigVersionOptions {
	_options.Version = core.Int64Ptr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigVersionOptions) SetHeaders(param map[string]string) *DeleteConfigVersionOptions {
	options.Headers = param
	return options
}

// DeleteProjectEnvironmentOptions : The DeleteProjectEnvironment options.
type DeleteProjectEnvironmentOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The environment ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProjectEnvironmentOptions : Instantiate DeleteProjectEnvironmentOptions
func (*ProjectV1) NewDeleteProjectEnvironmentOptions(projectID string, id string) *DeleteProjectEnvironmentOptions {
	return &DeleteProjectEnvironmentOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteProjectEnvironmentOptions) SetProjectID(projectID string) *DeleteProjectEnvironmentOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteProjectEnvironmentOptions) SetID(id string) *DeleteProjectEnvironmentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProjectEnvironmentOptions) SetHeaders(param map[string]string) *DeleteProjectEnvironmentOptions {
	options.Headers = param
	return options
}

// DeleteProjectOptions : The DeleteProject options.
type DeleteProjectOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

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

// SetHeaders : Allow user to set Headers
func (options *DeleteProjectOptions) SetHeaders(param map[string]string) *DeleteProjectOptions {
	options.Headers = param
	return options
}

// DeployConfigOptions : The DeployConfig options.
type DeployConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeployConfigOptions : Instantiate DeployConfigOptions
func (*ProjectV1) NewDeployConfigOptions(projectID string, id string) *DeployConfigOptions {
	return &DeployConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeployConfigOptions) SetProjectID(projectID string) *DeployConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeployConfigOptions) SetID(id string) *DeployConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeployConfigOptions) SetHeaders(param map[string]string) *DeployConfigOptions {
	options.Headers = param
	return options
}

// Environment : The definition of a project environment.
type Environment struct {
	// The environment id as a friendly name.
	ID *string `json:"id" validate:"required"`

	// The project referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The target account ID derived from the authentication block values.
	TargetAccount *string `json:"target_account,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The environment definition.
	Definition *EnvironmentDefinitionRequiredProperties `json:"definition" validate:"required"`
}

// UnmarshalEnvironment unmarshals an instance of Environment from the specified map of raw messages.
func UnmarshalEnvironment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Environment)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_account", &obj.TargetAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalEnvironmentDefinitionRequiredProperties)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentCollection : The list environment response.
type EnvironmentCollection struct {
	// The environments.
	Environments []Environment `json:"environments,omitempty"`
}

// UnmarshalEnvironmentCollection unmarshals an instance of EnvironmentCollection from the specified map of raw messages.
func UnmarshalEnvironmentCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentCollection)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalEnvironment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDefinitionNameDescription : The environment definition used in the project collection.
type EnvironmentDefinitionNameDescription struct {
	// The name of the environment.  It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The description of the environment.
	Description *string `json:"description,omitempty"`
}

// UnmarshalEnvironmentDefinitionNameDescription unmarshals an instance of EnvironmentDefinitionNameDescription from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionNameDescription(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionNameDescription)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// EnvironmentDefinitionProperties : The environment definition used for updates.
type EnvironmentDefinitionProperties struct {
	// The name of the environment.  It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The description of the environment.
	Description *string `json:"description,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`
}

// UnmarshalEnvironmentDefinitionProperties unmarshals an instance of EnvironmentDefinitionProperties from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDefinitionRequiredProperties : The environment definition.
type EnvironmentDefinitionRequiredProperties struct {
	// The name of the environment.  It is unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The description of the environment.
	Description *string `json:"description,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`
}

// NewEnvironmentDefinitionRequiredProperties : Instantiate EnvironmentDefinitionRequiredProperties (Generic Model Constructor)
func (*ProjectV1) NewEnvironmentDefinitionRequiredProperties(name string) (_model *EnvironmentDefinitionRequiredProperties, err error) {
	_model = &EnvironmentDefinitionRequiredProperties{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalEnvironmentDefinitionRequiredProperties unmarshals an instance of EnvironmentDefinitionRequiredProperties from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionRequiredProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionRequiredProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDeleteResponse : The delete environment response.
type EnvironmentDeleteResponse struct {
	// The environment id as a friendly name.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalEnvironmentDeleteResponse unmarshals an instance of EnvironmentDeleteResponse from the specified map of raw messages.
func UnmarshalEnvironmentDeleteResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDeleteResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentPrototype : The definition of a project environment.
type EnvironmentPrototype struct {
	// The environment definition.
	Definition *EnvironmentDefinitionRequiredProperties `json:"definition" validate:"required"`
}

// NewEnvironmentPrototype : Instantiate EnvironmentPrototype (Generic Model Constructor)
func (*ProjectV1) NewEnvironmentPrototype(definition *EnvironmentDefinitionRequiredProperties) (_model *EnvironmentPrototype, err error) {
	_model = &EnvironmentPrototype{
		Definition: definition,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalEnvironmentPrototype unmarshals an instance of EnvironmentPrototype from the specified map of raw messages.
func UnmarshalEnvironmentPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentPrototype)
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalEnvironmentDefinitionRequiredProperties)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ForceApproveOptions : The ForceApprove options.
type ForceApproveOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Notes on the project draft action. If this is a forced approve on the draft configuration, a non-empty comment is
	// required.
	Comment *string `json:"comment" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewForceApproveOptions : Instantiate ForceApproveOptions
func (*ProjectV1) NewForceApproveOptions(projectID string, id string, comment string) *ForceApproveOptions {
	return &ForceApproveOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Comment: core.StringPtr(comment),
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

// GetConfigOptions : The GetConfig options.
type GetConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

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

// SetHeaders : Allow user to set Headers
func (options *GetConfigOptions) SetHeaders(param map[string]string) *GetConfigOptions {
	options.Headers = param
	return options
}

// GetConfigVersionOptions : The GetConfigVersion options.
type GetConfigVersionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// The configuration version.
	Version *int64 `json:"version" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetConfigVersionOptions : Instantiate GetConfigVersionOptions
func (*ProjectV1) NewGetConfigVersionOptions(projectID string, id string, version int64) *GetConfigVersionOptions {
	return &GetConfigVersionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Version: core.Int64Ptr(version),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetConfigVersionOptions) SetProjectID(projectID string) *GetConfigVersionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetConfigVersionOptions) SetID(id string) *GetConfigVersionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetConfigVersionOptions) SetVersion(version int64) *GetConfigVersionOptions {
	_options.Version = core.Int64Ptr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigVersionOptions) SetHeaders(param map[string]string) *GetConfigVersionOptions {
	options.Headers = param
	return options
}

// GetProjectEnvironmentOptions : The GetProjectEnvironment options.
type GetProjectEnvironmentOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The environment ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProjectEnvironmentOptions : Instantiate GetProjectEnvironmentOptions
func (*ProjectV1) NewGetProjectEnvironmentOptions(projectID string, id string) *GetProjectEnvironmentOptions {
	return &GetProjectEnvironmentOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetProjectEnvironmentOptions) SetProjectID(projectID string) *GetProjectEnvironmentOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetProjectEnvironmentOptions) SetID(id string) *GetProjectEnvironmentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProjectEnvironmentOptions) SetHeaders(param map[string]string) *GetProjectEnvironmentOptions {
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

// LastActionWithSummary : The action job performed on the project configuration.
type LastActionWithSummary struct {
	// A URL.
	Href *string `json:"href" validate:"required"`

	// The result of the last action.
	Result *string `json:"result,omitempty"`

	// A brief summary of a pre/post action.
	PreJob *PrePostActionJobWithIdAndSummary `json:"pre_job,omitempty"`

	// A brief summary of a pre/post action.
	PostJob *PrePostActionJobWithIdAndSummary `json:"post_job,omitempty"`

	// A brief summary of an action.
	Job *ActionJobWithIdAndSummary `json:"job,omitempty"`
}

// Constants associated with the LastActionWithSummary.Result property.
// The result of the last action.
const (
	LastActionWithSummary_Result_Failed = "failed"
	LastActionWithSummary_Result_Passed = "passed"
)

// UnmarshalLastActionWithSummary unmarshals an instance of LastActionWithSummary from the specified map of raw messages.
func UnmarshalLastActionWithSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastActionWithSummary)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "pre_job", &obj.PreJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "post_job", &obj.PostJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "job", &obj.Job, UnmarshalActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastValidatedActionWithSummary : The action job performed on the project configuration.
type LastValidatedActionWithSummary struct {
	// A URL.
	Href *string `json:"href" validate:"required"`

	// The result of the last action.
	Result *string `json:"result,omitempty"`

	// A brief summary of a pre/post action.
	PreJob *PrePostActionJobWithIdAndSummary `json:"pre_job,omitempty"`

	// A brief summary of a pre/post action.
	PostJob *PrePostActionJobWithIdAndSummary `json:"post_job,omitempty"`

	// A brief summary of an action.
	Job *ActionJobWithIdAndSummary `json:"job,omitempty"`

	// The cost estimate of the configuration.
	// It only exists after the first configuration validation.
	CostEstimate *ProjectConfigMetadataCostEstimate `json:"cost_estimate,omitempty"`

	// The Code Risk Analyzer logs of the configuration.
	CraLogs *ProjectConfigMetadataCodeRiskAnalyzerLogs `json:"cra_logs,omitempty"`
}

// Constants associated with the LastValidatedActionWithSummary.Result property.
// The result of the last action.
const (
	LastValidatedActionWithSummary_Result_Failed = "failed"
	LastValidatedActionWithSummary_Result_Passed = "passed"
)

// UnmarshalLastValidatedActionWithSummary unmarshals an instance of LastValidatedActionWithSummary from the specified map of raw messages.
func UnmarshalLastValidatedActionWithSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastValidatedActionWithSummary)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "pre_job", &obj.PreJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "post_job", &obj.PostJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "job", &obj.Job, UnmarshalActionJobWithIdAndSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cost_estimate", &obj.CostEstimate, UnmarshalProjectConfigMetadataCostEstimate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cra_logs", &obj.CraLogs, UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListConfigResourcesOptions : The ListConfigResources options.
type ListConfigResourcesOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListConfigResourcesOptions : Instantiate ListConfigResourcesOptions
func (*ProjectV1) NewListConfigResourcesOptions(projectID string, id string) *ListConfigResourcesOptions {
	return &ListConfigResourcesOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListConfigResourcesOptions) SetProjectID(projectID string) *ListConfigResourcesOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListConfigResourcesOptions) SetID(id string) *ListConfigResourcesOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigResourcesOptions) SetHeaders(param map[string]string) *ListConfigResourcesOptions {
	options.Headers = param
	return options
}

// ListConfigVersionsOptions : The ListConfigVersions options.
type ListConfigVersionsOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListConfigVersionsOptions : Instantiate ListConfigVersionsOptions
func (*ProjectV1) NewListConfigVersionsOptions(projectID string, id string) *ListConfigVersionsOptions {
	return &ListConfigVersionsOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListConfigVersionsOptions) SetProjectID(projectID string) *ListConfigVersionsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListConfigVersionsOptions) SetID(id string) *ListConfigVersionsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigVersionsOptions) SetHeaders(param map[string]string) *ListConfigVersionsOptions {
	options.Headers = param
	return options
}

// ListConfigsOptions : The ListConfigs options.
type ListConfigsOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

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

// SetHeaders : Allow user to set Headers
func (options *ListConfigsOptions) SetHeaders(param map[string]string) *ListConfigsOptions {
	options.Headers = param
	return options
}

// ListProjectEnvironmentsOptions : The ListProjectEnvironments options.
type ListProjectEnvironmentsOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProjectEnvironmentsOptions : Instantiate ListProjectEnvironmentsOptions
func (*ProjectV1) NewListProjectEnvironmentsOptions(projectID string) *ListProjectEnvironmentsOptions {
	return &ListProjectEnvironmentsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListProjectEnvironmentsOptions) SetProjectID(projectID string) *ListProjectEnvironmentsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProjectEnvironmentsOptions) SetHeaders(param map[string]string) *ListProjectEnvironmentsOptions {
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

// OutputValue : OutputValue struct
type OutputValue struct {
	// The variable name.
	Name *string `json:"name" validate:"required"`

	// A short explanation of the output value.
	Description *string `json:"description,omitempty"`

	// Can be any value - a string, number, boolean, array, or object.
	Value map[string]interface{} `json:"value,omitempty"`
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
	// A URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginationLink unmarshals an instance of PaginationLink from the specified map of raw messages.
func UnmarshalPaginationLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginationLink)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobSummary : A brief summary of a pre/post action job.
type PrePostActionJobSummary struct {
	// The ID of the Schematics action job that ran as part of the pre/post job.
	JobID *string `json:"job_id" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	StartTime *strfmt.DateTime `json:"start_time,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	EndTime *strfmt.DateTime `json:"end_time,omitempty"`

	// The number of tasks run in the job.
	Tasks *int64 `json:"tasks,omitempty"`

	// The number of tasks that successfully ran in the job.
	Ok *int64 `json:"ok,omitempty"`

	// The number of tasks that failed in the job.
	Failed *int64 `json:"failed,omitempty"`

	// The number of tasks that were skipped in the job.
	Skipped *int64 `json:"skipped,omitempty"`

	// The number of tasks that were changed in the job.
	Changed *int64 `json:"changed,omitempty"`

	// A system-level error from the pipeline that ran for this specific pre- and post-job.
	ProjectError *PrePostActionJobSystemError `json:"project_error,omitempty"`
}

// UnmarshalPrePostActionJobSummary unmarshals an instance of PrePostActionJobSummary from the specified map of raw messages.
func UnmarshalPrePostActionJobSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrePostActionJobSummary)
	err = core.UnmarshalPrimitive(m, "job_id", &obj.JobID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tasks", &obj.Tasks)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
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
	err = core.UnmarshalPrimitive(m, "changed", &obj.Changed)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project_error", &obj.ProjectError, UnmarshalPrePostActionJobSystemError)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobSystemError : System level error captured in the Projects Pipelines for pre/post job.
type PrePostActionJobSystemError struct {
	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	Timestamp *strfmt.DateTime `json:"timestamp" validate:"required"`

	// Id of user that triggered pipeline that ran pre/post job.
	UserID *string `json:"user_id" validate:"required"`

	// HTTP status code for the error.
	StatusCode *string `json:"status_code" validate:"required"`

	// Summary description of the error.
	Description *string `json:"description" validate:"required"`

	// Detailed message from the source error.
	ErrorResponse *string `json:"error_response,omitempty"`
}

// UnmarshalPrePostActionJobSystemError unmarshals an instance of PrePostActionJobSystemError from the specified map of raw messages.
func UnmarshalPrePostActionJobSystemError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrePostActionJobSystemError)
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_response", &obj.ErrorResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobWithIdAndSummary : A brief summary of a pre/post action.
type PrePostActionJobWithIdAndSummary struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// A brief summary of a pre/post action job.
	Summary *PrePostActionJobSummary `json:"summary" validate:"required"`
}

// UnmarshalPrePostActionJobWithIdAndSummary unmarshals an instance of PrePostActionJobWithIdAndSummary from the specified map of raw messages.
func UnmarshalPrePostActionJobWithIdAndSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrePostActionJobWithIdAndSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalPrePostActionJobSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Project : The canonical schema of a project.
type Project struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The cumulative list of needs attention items for a project. If the view is successfully retrieved, an array which
	// could be empty is returned.
	CumulativeNeedsAttentionView []CumulativeNeedsAttention `json:"cumulative_needs_attention_view,omitempty"`

	// True indicates that the fetch of the needs attention items failed. It only exists if there was an error while
	// retrieving the cumulative needs attention view.
	CumulativeNeedsAttentionViewError *bool `json:"cumulative_needs_attention_view_error,omitempty"`

	// The unique project ID.
	ID *string `json:"id" validate:"required"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location" validate:"required"`

	// The resource group id where the project's data and tools are created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The project status value.
	State *string `json:"state" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The resource group name where the project's data and tools are created.
	ResourceGroup *string `json:"resource_group" validate:"required"`

	// The CRN of the event notifications instance if one is connected to this project.
	EventNotificationsCrn *string `json:"event_notifications_crn,omitempty"`

	// The project configurations. These configurations are only included in the response of creating a project if a
	// configs array is specified in the request payload.
	Configs []ProjectConfigSummary `json:"configs,omitempty"`

	// The project environments. These environments are only included in the response if project environments were created
	// on the project.
	Environments []ProjectEnvironmentSummary `json:"environments,omitempty"`

	// The definition of the project.
	Definition *ProjectDefinitionProperties `json:"definition" validate:"required"`
}

// Constants associated with the Project.State property.
// The project status value.
const (
	Project_State_Deleting = "deleting"
	Project_State_DeletingFailed = "deleting_failed"
	Project_State_Ready = "ready"
)

// UnmarshalProject unmarshals an instance of Project from the specified map of raw messages.
func UnmarshalProject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Project)
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
	err = core.UnmarshalPrimitive(m, "cumulative_needs_attention_view_error", &obj.CumulativeNeedsAttentionViewError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
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
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_notifications_crn", &obj.EventNotificationsCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfigSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalProjectEnvironmentSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionProperties)
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

	// A pagination link.
	First *PaginationLink `json:"first" validate:"required"`

	// A pagination link.
	Next *PaginationLink `json:"next,omitempty"`

	// An array of projects.
	Projects []ProjectSummary `json:"projects,omitempty"`
}

// UnmarshalProjectCollection unmarshals an instance of ProjectCollection from the specified map of raw messages.
func UnmarshalProjectCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "projects", &obj.Projects, UnmarshalProjectSummary)
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
	start, err := core.GetQueryParam(resp.Next.Href, "start")
	if err != nil || start == nil {
		return nil, err
	}
	return start, nil
}

// ProjectComplianceProfile : The profile required for compliance.
type ProjectComplianceProfile struct {
	// The unique ID for that compliance profile.
	ID *string `json:"id,omitempty"`

	// A unique ID for an instance of a compliance profile.
	InstanceID *string `json:"instance_id,omitempty"`

	// The location of the compliance instance.
	InstanceLocation *string `json:"instance_location,omitempty"`

	// A unique ID for the attachment to a compliance profile.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The name of the compliance profile.
	ProfileName *string `json:"profile_name,omitempty"`
}

// UnmarshalProjectComplianceProfile unmarshals an instance of ProjectComplianceProfile from the specified map of raw messages.
func UnmarshalProjectComplianceProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectComplianceProfile)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_location", &obj.InstanceLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfig : The canonical schema of a project configuration.
type ProjectConfig struct {
	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id" validate:"required"`

	// The version of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// The flag that indicates whether the version of the configuration is draft, or active.
	IsDraft *bool `json:"is_draft" validate:"required"`

	// The needs attention state of a configuration.
	NeedsAttentionState []map[string]interface{} `json:"needs_attention_state,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// The last approved metadata of the configuration.
	LastApproved *ProjectConfigMetadataLastApproved `json:"last_approved,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	LastSavedAt *strfmt.DateTime `json:"last_saved_at,omitempty"`

	// The action job performed on the project configuration.
	LastValidated *LastValidatedActionWithSummary `json:"last_validated,omitempty"`

	// The action job performed on the project configuration.
	LastDeployed *LastActionWithSummary `json:"last_deployed,omitempty"`

	// The action job performed on the project configuration.
	LastUndeployed *LastActionWithSummary `json:"last_undeployed,omitempty"`

	// The outputs of a Schematics template property.
	Outputs []OutputValue `json:"outputs,omitempty"`

	// The project referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// The references used in the config to resolve input values.
	References map[string]interface{} `json:"references,omitempty"`

	// A schematics workspace associated to a project configuration, with scripts.
	Schematics *SchematicsMetadata `json:"schematics,omitempty"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// The flag that indicates whether a configuration update is available.
	UpdateAvailable *bool `json:"update_available,omitempty"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	Definition ProjectConfigResponseDefinitionIntf `json:"definition" validate:"required"`

	// The project configuration version.
	ApprovedVersion *ProjectConfigVersionSummary `json:"approved_version,omitempty"`

	// The project configuration version.
	DeployedVersion *ProjectConfigVersionSummary `json:"deployed_version,omitempty"`
}

// Constants associated with the ProjectConfig.State property.
// The state of the configuration.
const (
	ProjectConfig_State_Applied = "applied"
	ProjectConfig_State_ApplyFailed = "apply_failed"
	ProjectConfig_State_Approved = "approved"
	ProjectConfig_State_Deleted = "deleted"
	ProjectConfig_State_Deleting = "deleting"
	ProjectConfig_State_DeletingFailed = "deleting_failed"
	ProjectConfig_State_Deployed = "deployed"
	ProjectConfig_State_Deploying = "deploying"
	ProjectConfig_State_DeployingFailed = "deploying_failed"
	ProjectConfig_State_Discarded = "discarded"
	ProjectConfig_State_Draft = "draft"
	ProjectConfig_State_Superseded = "superseded"
	ProjectConfig_State_Undeploying = "undeploying"
	ProjectConfig_State_UndeployingFailed = "undeploying_failed"
	ProjectConfig_State_Validated = "validated"
	ProjectConfig_State_Validating = "validating"
	ProjectConfig_State_ValidatingFailed = "validating_failed"
)

// UnmarshalProjectConfig unmarshals an instance of ProjectConfig from the specified map of raw messages.
func UnmarshalProjectConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfig)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_draft", &obj.IsDraft)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "needs_attention_state", &obj.NeedsAttentionState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_approved", &obj.LastApproved, UnmarshalProjectConfigMetadataLastApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_saved_at", &obj.LastSavedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_validated", &obj.LastValidated, UnmarshalLastValidatedActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_deployed", &obj.LastDeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_undeployed", &obj.LastUndeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalOutputValue)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "references", &obj.References)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_available", &obj.UpdateAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigResponseDefinition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "approved_version", &obj.ApprovedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deployed_version", &obj.DeployedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigAuth : The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
type ProjectConfigAuth struct {
	// The trusted profile ID.
	TrustedProfileID *string `json:"trusted_profile_id,omitempty"`

	// The authorization method. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Method *string `json:"method,omitempty"`

	// The IBM Cloud API Key.
	ApiKey *string `json:"api_key,omitempty"`
}

// Constants associated with the ProjectConfigAuth.Method property.
// The authorization method. You can authorize by using a trusted profile or an API key in Secrets Manager.
const (
	ProjectConfigAuth_Method_ApiKey = "api_key"
	ProjectConfigAuth_Method_TrustedProfile = "trusted_profile"
)

// UnmarshalProjectConfigAuth unmarshals an instance of ProjectConfigAuth from the specified map of raw messages.
func UnmarshalProjectConfigAuth(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigAuth)
	err = core.UnmarshalPrimitive(m, "trusted_profile_id", &obj.TrustedProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigCollection : The project configuration list.
type ProjectConfigCollection struct {
	// The collection list operation response schema that should define the array property with the name "configs".
	Configs []ProjectConfigSummary `json:"configs,omitempty"`
}

// UnmarshalProjectConfigCollection unmarshals an instance of ProjectConfigCollection from the specified map of raw messages.
func UnmarshalProjectConfigCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigCollection)
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfigSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionNameDescription : The name and description of a project configuration.
type ProjectConfigDefinitionNameDescription struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalProjectConfigDefinitionNameDescription unmarshals an instance of ProjectConfigDefinitionNameDescription from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionNameDescription(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionNameDescription)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// ProjectConfigDelete : Deletes the configuration response.
type ProjectConfigDelete struct {
	// The unique configuration ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalProjectConfigDelete unmarshals an instance of ProjectConfigDelete from the specified map of raw messages.
func UnmarshalProjectConfigDelete(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDelete)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataCodeRiskAnalyzerLogs : The Code Risk Analyzer logs of the configuration.
type ProjectConfigMetadataCodeRiskAnalyzerLogs struct {
	// The version of the Code Risk Analyzer logs of the configuration. This will change as the Code Risk Analyzer is
	// updated.
	CraVersion *string `json:"cra_version,omitempty"`

	// The schema version of Code Risk Analyzer logs of the configuration.
	SchemaVersion *string `json:"schema_version,omitempty"`

	// The status of the Code Risk Analyzer logs of the configuration.
	Status *string `json:"status,omitempty"`

	// The Code Risk Analyzer logs summary of the configuration.
	Summary *CodeRiskAnalyzerLogsSummary `json:"summary,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	Timestamp *strfmt.DateTime `json:"timestamp,omitempty"`
}

// Constants associated with the ProjectConfigMetadataCodeRiskAnalyzerLogs.Status property.
// The status of the Code Risk Analyzer logs of the configuration.
const (
	ProjectConfigMetadataCodeRiskAnalyzerLogs_Status_Failed = "failed"
	ProjectConfigMetadataCodeRiskAnalyzerLogs_Status_Passed = "passed"
)

// UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs unmarshals an instance of ProjectConfigMetadataCodeRiskAnalyzerLogs from the specified map of raw messages.
func UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigMetadataCodeRiskAnalyzerLogs)
	err = core.UnmarshalPrimitive(m, "cra_version", &obj.CraVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schema_version", &obj.SchemaVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalCodeRiskAnalyzerLogsSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataCostEstimate : The cost estimate of the configuration. It only exists after the first configuration validation.
type ProjectConfigMetadataCostEstimate struct {
	// The version of the cost estimate of the configuration.
	Version *string `json:"version,omitempty"`

	// The currency of the cost estimate of the configuration.
	Currency *string `json:"currency,omitempty"`

	// The total hourly cost estimate of the configuration.
	TotalHourlyCost *string `json:"totalHourlyCost,omitempty"`

	// The total monthly cost estimate of the configuration.
	TotalMonthlyCost *string `json:"totalMonthlyCost,omitempty"`

	// The past total hourly cost estimate of the configuration.
	PastTotalHourlyCost *string `json:"pastTotalHourlyCost,omitempty"`

	// The past total monthly cost estimate of the configuration.
	PastTotalMonthlyCost *string `json:"pastTotalMonthlyCost,omitempty"`

	// The difference between current and past total hourly cost estimates of the configuration.
	DiffTotalHourlyCost *string `json:"diffTotalHourlyCost,omitempty"`

	// The difference between current and past total monthly cost estimates of the configuration.
	DiffTotalMonthlyCost *string `json:"diffTotalMonthlyCost,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	TimeGenerated *strfmt.DateTime `json:"timeGenerated,omitempty"`

	// The unique ID.
	UserID *string `json:"user_id,omitempty"`
}

// UnmarshalProjectConfigMetadataCostEstimate unmarshals an instance of ProjectConfigMetadataCostEstimate from the specified map of raw messages.
func UnmarshalProjectConfigMetadataCostEstimate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigMetadataCostEstimate)
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "currency", &obj.Currency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalHourlyCost", &obj.TotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "totalMonthlyCost", &obj.TotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalHourlyCost", &obj.PastTotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalMonthlyCost", &obj.PastTotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalHourlyCost", &obj.DiffTotalHourlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalMonthlyCost", &obj.DiffTotalMonthlyCost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeGenerated", &obj.TimeGenerated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataLastApproved : The last approved metadata of the configuration.
type ProjectConfigMetadataLastApproved struct {
	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	At *strfmt.DateTime `json:"at" validate:"required"`

	// The comment left by the user who approved the configuration.
	Comment *string `json:"comment,omitempty"`

	// The flag that indicates whether the approval was forced approved.
	IsForced *bool `json:"is_forced" validate:"required"`

	// The unique ID.
	UserID *string `json:"user_id" validate:"required"`
}

// UnmarshalProjectConfigMetadataLastApproved unmarshals an instance of ProjectConfigMetadataLastApproved from the specified map of raw messages.
func UnmarshalProjectConfigMetadataLastApproved(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigMetadataLastApproved)
	err = core.UnmarshalPrimitive(m, "at", &obj.At)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "comment", &obj.Comment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_forced", &obj.IsForced)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPatchDefinitionBlock : ProjectConfigPatchDefinitionBlock struct
// Models which "extend" this model:
// - ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties
// - ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties
type ProjectConfigPatchDefinitionBlock struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}
func (*ProjectConfigPatchDefinitionBlock) isaProjectConfigPatchDefinitionBlock() bool {
	return true
}

type ProjectConfigPatchDefinitionBlockIntf interface {
	isaProjectConfigPatchDefinitionBlock() bool
}

// UnmarshalProjectConfigPatchDefinitionBlock unmarshals an instance of ProjectConfigPatchDefinitionBlock from the specified map of raw messages.
func UnmarshalProjectConfigPatchDefinitionBlock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchDefinitionBlock)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototype : The input of a project configuration.
type ProjectConfigPrototype struct {
	Definition ProjectConfigPrototypeDefinitionBlockIntf `json:"definition" validate:"required"`

	// A Schematics workspace to use for deploying this configuration.
	// Either schematics.workspace_crn, definition.locator_id, or both must be specified.
	Schematics *SchematicsWorkspace `json:"schematics,omitempty"`
}

// NewProjectConfigPrototype : Instantiate ProjectConfigPrototype (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigPrototype(definition ProjectConfigPrototypeDefinitionBlockIntf) (_model *ProjectConfigPrototype, err error) {
	_model = &ProjectConfigPrototype{
		Definition: definition,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProjectConfigPrototype unmarshals an instance of ProjectConfigPrototype from the specified map of raw messages.
func UnmarshalProjectConfigPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototype)
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigPrototypeDefinitionBlock)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsWorkspace)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototypeDefinitionBlock : ProjectConfigPrototypeDefinitionBlock struct
// Models which "extend" this model:
// - ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties
// - ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties
type ProjectConfigPrototypeDefinitionBlock struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}
func (*ProjectConfigPrototypeDefinitionBlock) isaProjectConfigPrototypeDefinitionBlock() bool {
	return true
}

type ProjectConfigPrototypeDefinitionBlockIntf interface {
	isaProjectConfigPrototypeDefinitionBlock() bool
}

// UnmarshalProjectConfigPrototypeDefinitionBlock unmarshals an instance of ProjectConfigPrototypeDefinitionBlock from the specified map of raw messages.
func UnmarshalProjectConfigPrototypeDefinitionBlock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototypeDefinitionBlock)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResource : ProjectConfigResource struct
type ProjectConfigResource struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	ResourceCrn *string `json:"resource_crn,omitempty"`

	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty"`

	// The resource type.
	ResourceType *string `json:"resource_type,omitempty"`

	// The flag that indicates whether the status of the resource is tainted.
	ResourceTainted *bool `json:"resource_tainted,omitempty"`

	// The resource group of the resource.
	ResourceGroupName *string `json:"resource_group_name,omitempty"`
}

// UnmarshalProjectConfigResource unmarshals an instance of ProjectConfigResource from the specified map of raw messages.
func UnmarshalProjectConfigResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResource)
	err = core.UnmarshalPrimitive(m, "resource_crn", &obj.ResourceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_tainted", &obj.ResourceTainted)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResourceCollection : The project configuration resource list.
type ProjectConfigResourceCollection struct {
	// The collection list operation response schema that defines the array property with the name `resources`.
	Resources []ProjectConfigResource `json:"resources,omitempty"`

	// The total number of resources deployed by the configuration.
	ResourcesCount *int64 `json:"resources_count" validate:"required"`
}

// UnmarshalProjectConfigResourceCollection unmarshals an instance of ProjectConfigResourceCollection from the specified map of raw messages.
func UnmarshalProjectConfigResourceCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResourceCollection)
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalProjectConfigResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_count", &obj.ResourcesCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResponseDefinition : ProjectConfigResponseDefinition struct
// Models which "extend" this model:
// - ProjectConfigResponseDefinitionDAConfigDefinitionProperties
// - ProjectConfigResponseDefinitionResourceConfigDefinitionProperties
type ProjectConfigResponseDefinition struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}
func (*ProjectConfigResponseDefinition) isaProjectConfigResponseDefinition() bool {
	return true
}

type ProjectConfigResponseDefinitionIntf interface {
	isaProjectConfigResponseDefinition() bool
}

// UnmarshalProjectConfigResponseDefinition unmarshals an instance of ProjectConfigResponseDefinition from the specified map of raw messages.
func UnmarshalProjectConfigResponseDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResponseDefinition)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigSummary : ProjectConfigSummary struct
type ProjectConfigSummary struct {
	// The project configuration version.
	ApprovedVersion *ProjectConfigVersionSummary `json:"approved_version,omitempty"`

	// The project configuration version.
	DeployedVersion *ProjectConfigVersionSummary `json:"deployed_version,omitempty"`

	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id" validate:"required"`

	// The version of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The name and description of a project configuration.
	Definition *ProjectConfigDefinitionNameDescription `json:"definition" validate:"required"`

	// The project referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// The configuration type.
	DeploymentModel *string `json:"deployment_model,omitempty"`
}

// Constants associated with the ProjectConfigSummary.State property.
// The state of the configuration.
const (
	ProjectConfigSummary_State_Applied = "applied"
	ProjectConfigSummary_State_ApplyFailed = "apply_failed"
	ProjectConfigSummary_State_Approved = "approved"
	ProjectConfigSummary_State_Deleted = "deleted"
	ProjectConfigSummary_State_Deleting = "deleting"
	ProjectConfigSummary_State_DeletingFailed = "deleting_failed"
	ProjectConfigSummary_State_Deployed = "deployed"
	ProjectConfigSummary_State_Deploying = "deploying"
	ProjectConfigSummary_State_DeployingFailed = "deploying_failed"
	ProjectConfigSummary_State_Discarded = "discarded"
	ProjectConfigSummary_State_Draft = "draft"
	ProjectConfigSummary_State_Superseded = "superseded"
	ProjectConfigSummary_State_Undeploying = "undeploying"
	ProjectConfigSummary_State_UndeployingFailed = "undeploying_failed"
	ProjectConfigSummary_State_Validated = "validated"
	ProjectConfigSummary_State_Validating = "validating"
	ProjectConfigSummary_State_ValidatingFailed = "validating_failed"
)

// Constants associated with the ProjectConfigSummary.DeploymentModel property.
// The configuration type.
const (
	ProjectConfigSummary_DeploymentModel_ProjectDeployed = "project_deployed"
	ProjectConfigSummary_DeploymentModel_UserDeployed = "user_deployed"
)

// UnmarshalProjectConfigSummary unmarshals an instance of ProjectConfigSummary from the specified map of raw messages.
func UnmarshalProjectConfigSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigSummary)
	err = core.UnmarshalModel(m, "approved_version", &obj.ApprovedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deployed_version", &obj.DeployedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigDefinitionNameDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_model", &obj.DeploymentModel)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersion : A specific version of a project configuration.
type ProjectConfigVersion struct {
	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id" validate:"required"`

	// The version of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// The flag that indicates whether the version of the configuration is draft, or active.
	IsDraft *bool `json:"is_draft" validate:"required"`

	// The needs attention state of a configuration.
	NeedsAttentionState []map[string]interface{} `json:"needs_attention_state,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// The last approved metadata of the configuration.
	LastApproved *ProjectConfigMetadataLastApproved `json:"last_approved,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	LastSavedAt *strfmt.DateTime `json:"last_saved_at,omitempty"`

	// The action job performed on the project configuration.
	LastValidated *LastValidatedActionWithSummary `json:"last_validated,omitempty"`

	// The action job performed on the project configuration.
	LastDeployed *LastActionWithSummary `json:"last_deployed,omitempty"`

	// The action job performed on the project configuration.
	LastUndeployed *LastActionWithSummary `json:"last_undeployed,omitempty"`

	// The outputs of a Schematics template property.
	Outputs []OutputValue `json:"outputs,omitempty"`

	// The project referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// The references used in the config to resolve input values.
	References map[string]interface{} `json:"references,omitempty"`

	// A schematics workspace associated to a project configuration, with scripts.
	Schematics *SchematicsMetadata `json:"schematics,omitempty"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// The flag that indicates whether a configuration update is available.
	UpdateAvailable *bool `json:"update_available,omitempty"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	Definition ProjectConfigResponseDefinitionIntf `json:"definition" validate:"required"`
}

// Constants associated with the ProjectConfigVersion.State property.
// The state of the configuration.
const (
	ProjectConfigVersion_State_Applied = "applied"
	ProjectConfigVersion_State_ApplyFailed = "apply_failed"
	ProjectConfigVersion_State_Approved = "approved"
	ProjectConfigVersion_State_Deleted = "deleted"
	ProjectConfigVersion_State_Deleting = "deleting"
	ProjectConfigVersion_State_DeletingFailed = "deleting_failed"
	ProjectConfigVersion_State_Deployed = "deployed"
	ProjectConfigVersion_State_Deploying = "deploying"
	ProjectConfigVersion_State_DeployingFailed = "deploying_failed"
	ProjectConfigVersion_State_Discarded = "discarded"
	ProjectConfigVersion_State_Draft = "draft"
	ProjectConfigVersion_State_Superseded = "superseded"
	ProjectConfigVersion_State_Undeploying = "undeploying"
	ProjectConfigVersion_State_UndeployingFailed = "undeploying_failed"
	ProjectConfigVersion_State_Validated = "validated"
	ProjectConfigVersion_State_Validating = "validating"
	ProjectConfigVersion_State_ValidatingFailed = "validating_failed"
)

// UnmarshalProjectConfigVersion unmarshals an instance of ProjectConfigVersion from the specified map of raw messages.
func UnmarshalProjectConfigVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "is_draft", &obj.IsDraft)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "needs_attention_state", &obj.NeedsAttentionState)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_approved", &obj.LastApproved, UnmarshalProjectConfigMetadataLastApproved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_saved_at", &obj.LastSavedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_validated", &obj.LastValidated, UnmarshalLastValidatedActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_deployed", &obj.LastDeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_undeployed", &obj.LastUndeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalOutputValue)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "references", &obj.References)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsMetadata)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "update_available", &obj.UpdateAvailable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigResponseDefinition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersionSummary : The project configuration version.
type ProjectConfigVersionSummary struct {
	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// The version number of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`
}

// Constants associated with the ProjectConfigVersionSummary.State property.
// The state of the configuration.
const (
	ProjectConfigVersionSummary_State_Applied = "applied"
	ProjectConfigVersionSummary_State_ApplyFailed = "apply_failed"
	ProjectConfigVersionSummary_State_Approved = "approved"
	ProjectConfigVersionSummary_State_Deleted = "deleted"
	ProjectConfigVersionSummary_State_Deleting = "deleting"
	ProjectConfigVersionSummary_State_DeletingFailed = "deleting_failed"
	ProjectConfigVersionSummary_State_Deployed = "deployed"
	ProjectConfigVersionSummary_State_Deploying = "deploying"
	ProjectConfigVersionSummary_State_DeployingFailed = "deploying_failed"
	ProjectConfigVersionSummary_State_Discarded = "discarded"
	ProjectConfigVersionSummary_State_Draft = "draft"
	ProjectConfigVersionSummary_State_Superseded = "superseded"
	ProjectConfigVersionSummary_State_Undeploying = "undeploying"
	ProjectConfigVersionSummary_State_UndeployingFailed = "undeploying_failed"
	ProjectConfigVersionSummary_State_Validated = "validated"
	ProjectConfigVersionSummary_State_Validating = "validating"
	ProjectConfigVersionSummary_State_ValidatingFailed = "validating_failed"
)

// UnmarshalProjectConfigVersionSummary unmarshals an instance of ProjectConfigVersionSummary from the specified map of raw messages.
func UnmarshalProjectConfigVersionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersionSummary)
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersionSummaryCollection : The project configuration version list.
type ProjectConfigVersionSummaryCollection struct {
	// The collection list operation response schema that defines the array property with the name `versions`.
	Versions []ProjectConfigVersionSummary `json:"versions,omitempty"`
}

// UnmarshalProjectConfigVersionSummaryCollection unmarshals an instance of ProjectConfigVersionSummaryCollection from the specified map of raw messages.
func UnmarshalProjectConfigVersionSummaryCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersionSummaryCollection)
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectDefinitionProperties : The definition of the project.
type ProjectDefinitionProperties struct {
	// The name of the project.  It is unique within the account across regions.
	Name *string `json:"name" validate:"required"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. It is possible to create
	// a project without providing a description.
	Description *string `json:"description,omitempty"`

	// The policy that indicates whether the resources are destroyed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete" validate:"required"`
}

// UnmarshalProjectDefinitionProperties unmarshals an instance of ProjectDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectDefinitionReference : The definition of the project reference.
type ProjectDefinitionReference struct {
	// The name of the project.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalProjectDefinitionReference unmarshals an instance of ProjectDefinitionReference from the specified map of raw messages.
func UnmarshalProjectDefinitionReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectDefinitionReference)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectEnvironmentSummary : The environment metadata.
type ProjectEnvironmentSummary struct {
	// The environment id as a friendly name.
	ID *string `json:"id" validate:"required"`

	// The project referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The environment definition used in the project collection.
	Definition *EnvironmentDefinitionNameDescription `json:"definition" validate:"required"`
}

// UnmarshalProjectEnvironmentSummary unmarshals an instance of ProjectEnvironmentSummary from the specified map of raw messages.
func UnmarshalProjectEnvironmentSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectEnvironmentSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalEnvironmentDefinitionNameDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectPatchDefinitionBlock : The definition of the project.
type ProjectPatchDefinitionBlock struct {
	// The name of the project.  It is unique within the account across regions.
	Name *string `json:"name,omitempty"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. It is possible to create
	// a project without providing a description.
	Description *string `json:"description,omitempty"`

	// The policy that indicates whether the resources are destroyed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete,omitempty"`
}

// UnmarshalProjectPatchDefinitionBlock unmarshals an instance of ProjectPatchDefinitionBlock from the specified map of raw messages.
func UnmarshalProjectPatchDefinitionBlock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectPatchDefinitionBlock)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectPrototypeDefinition : The definition of the project.
type ProjectPrototypeDefinition struct {
	// The name of the project.  It is unique within the account across regions.
	Name *string `json:"name" validate:"required"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. It is possible to create
	// a project without providing a description.
	Description *string `json:"description,omitempty"`

	// The policy that indicates whether the resources are undeployed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete,omitempty"`
}

// NewProjectPrototypeDefinition : Instantiate ProjectPrototypeDefinition (Generic Model Constructor)
func (*ProjectV1) NewProjectPrototypeDefinition(name string) (_model *ProjectPrototypeDefinition, err error) {
	_model = &ProjectPrototypeDefinition{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProjectPrototypeDefinition unmarshals an instance of ProjectPrototypeDefinition from the specified map of raw messages.
func UnmarshalProjectPrototypeDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectPrototypeDefinition)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectReference : The project referenced by this resource.
type ProjectReference struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// The definition of the project reference.
	Definition *ProjectDefinitionReference `json:"definition" validate:"required"`

	// An IBM Cloud resource name, which uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalProjectReference unmarshals an instance of ProjectReference from the specified map of raw messages.
func UnmarshalProjectReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionReference)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectSummary : ProjectSummary struct
type ProjectSummary struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The cumulative list of needs attention items for a project. If the view is successfully retrieved, an array which
	// could be empty is returned.
	CumulativeNeedsAttentionView []CumulativeNeedsAttention `json:"cumulative_needs_attention_view,omitempty"`

	// True indicates that the fetch of the needs attention items failed. It only exists if there was an error while
	// retrieving the cumulative needs attention view.
	CumulativeNeedsAttentionViewError *bool `json:"cumulative_needs_attention_view_error,omitempty"`

	// The unique project ID.
	ID *string `json:"id" validate:"required"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location" validate:"required"`

	// The resource group id where the project's data and tools are created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The project status value.
	State *string `json:"state" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The definition of the project.
	Definition *ProjectDefinitionProperties `json:"definition" validate:"required"`
}

// Constants associated with the ProjectSummary.State property.
// The project status value.
const (
	ProjectSummary_State_Deleting = "deleting"
	ProjectSummary_State_DeletingFailed = "deleting_failed"
	ProjectSummary_State_Ready = "ready"
)

// UnmarshalProjectSummary unmarshals an instance of ProjectSummary from the specified map of raw messages.
func UnmarshalProjectSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectSummary)
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
	err = core.UnmarshalPrimitive(m, "cumulative_needs_attention_view_error", &obj.CumulativeNeedsAttentionViewError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
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
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionProperties)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsMetadata : A schematics workspace associated to a project configuration, with scripts.
type SchematicsMetadata struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	WorkspaceCrn *string `json:"workspace_crn,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	ValidatePreScript *Script `json:"validate_pre_script,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	ValidatePostScript *Script `json:"validate_post_script,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	DeployPreScript *Script `json:"deploy_pre_script,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	DeployPostScript *Script `json:"deploy_post_script,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	UndeployPreScript *Script `json:"undeploy_pre_script,omitempty"`

	// A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
	// undeploy).
	UndeployPostScript *Script `json:"undeploy_post_script,omitempty"`
}

// UnmarshalSchematicsMetadata unmarshals an instance of SchematicsMetadata from the specified map of raw messages.
func UnmarshalSchematicsMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsMetadata)
	err = core.UnmarshalPrimitive(m, "workspace_crn", &obj.WorkspaceCrn)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validate_pre_script", &obj.ValidatePreScript, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "validate_post_script", &obj.ValidatePostScript, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deploy_pre_script", &obj.DeployPreScript, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deploy_post_script", &obj.DeployPostScript, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "undeploy_pre_script", &obj.UndeployPreScript, UnmarshalScript)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "undeploy_post_script", &obj.UndeployPostScript, UnmarshalScript)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsWorkspace : A Schematics workspace to use for deploying this configuration. Either schematics.workspace_crn,
// definition.locator_id, or both must be specified.
type SchematicsWorkspace struct {
	// An IBM Cloud resource name, which uniquely identifies a resource.
	WorkspaceCrn *string `json:"workspace_crn,omitempty"`
}

// UnmarshalSchematicsWorkspace unmarshals an instance of SchematicsWorkspace from the specified map of raw messages.
func UnmarshalSchematicsWorkspace(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsWorkspace)
	err = core.UnmarshalPrimitive(m, "workspace_crn", &obj.WorkspaceCrn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Script : A script to be run as part of a Project configuration, for a given stage (pre, post) and action (validate, deploy,
// undeploy).
type Script struct {
	// The type of the script.
	Type *string `json:"type,omitempty"`

	// The path to this script within the current version source.
	Path *string `json:"path,omitempty"`

	// The short description for this script.
	ShortDescription *string `json:"short_description,omitempty"`
}

// UnmarshalScript unmarshals an instance of Script from the specified map of raw messages.
func UnmarshalScript(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Script)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyncConfigOptions : The SyncConfig options.
type SyncConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// A Schematics workspace to use for deploying this configuration.
	// Either schematics.workspace_crn, definition.locator_id, or both must be specified.
	Schematics *SchematicsWorkspace `json:"schematics,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSyncConfigOptions : Instantiate SyncConfigOptions
func (*ProjectV1) NewSyncConfigOptions(projectID string, id string) *SyncConfigOptions {
	return &SyncConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *SyncConfigOptions) SetProjectID(projectID string) *SyncConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *SyncConfigOptions) SetID(id string) *SyncConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSchematics : Allow user to set Schematics
func (_options *SyncConfigOptions) SetSchematics(schematics *SchematicsWorkspace) *SyncConfigOptions {
	_options.Schematics = schematics
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SyncConfigOptions) SetHeaders(param map[string]string) *SyncConfigOptions {
	options.Headers = param
	return options
}

// TerraformLogAnalyzerErrorMessage : The error message parsed by the Terraform Log Analyzer.
type TerraformLogAnalyzerErrorMessage struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of TerraformLogAnalyzerErrorMessage
func (o *TerraformLogAnalyzerErrorMessage) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of TerraformLogAnalyzerErrorMessage
func (o *TerraformLogAnalyzerErrorMessage) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of TerraformLogAnalyzerErrorMessage
func (o *TerraformLogAnalyzerErrorMessage) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of TerraformLogAnalyzerErrorMessage
func (o *TerraformLogAnalyzerErrorMessage) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of TerraformLogAnalyzerErrorMessage
func (o *TerraformLogAnalyzerErrorMessage) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalTerraformLogAnalyzerErrorMessage unmarshals an instance of TerraformLogAnalyzerErrorMessage from the specified map of raw messages.
func UnmarshalTerraformLogAnalyzerErrorMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TerraformLogAnalyzerErrorMessage)
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

// TerraformLogAnalyzerSuccessMessage : The success message parsed by the Terraform Log Analyzer.
type TerraformLogAnalyzerSuccessMessage struct {
	// The resource type.
	ResourceType *string `json:"resource_type,omitempty"`

	// The time taken.
	TimeTaken *string `json:"time-taken,omitempty"`

	// The id.
	ID *string `json:"id,omitempty"`
}

// UnmarshalTerraformLogAnalyzerSuccessMessage unmarshals an instance of TerraformLogAnalyzerSuccessMessage from the specified map of raw messages.
func UnmarshalTerraformLogAnalyzerSuccessMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TerraformLogAnalyzerSuccessMessage)
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "time-taken", &obj.TimeTaken)
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

// UndeployConfigOptions : The UndeployConfig options.
type UndeployConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUndeployConfigOptions : Instantiate UndeployConfigOptions
func (*ProjectV1) NewUndeployConfigOptions(projectID string, id string) *UndeployConfigOptions {
	return &UndeployConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UndeployConfigOptions) SetProjectID(projectID string) *UndeployConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UndeployConfigOptions) SetID(id string) *UndeployConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UndeployConfigOptions) SetHeaders(param map[string]string) *UndeployConfigOptions {
	options.Headers = param
	return options
}

// UpdateConfigOptions : The UpdateConfig options.
type UpdateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	Definition ProjectConfigPatchDefinitionBlockIntf `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateConfigOptions : Instantiate UpdateConfigOptions
func (*ProjectV1) NewUpdateConfigOptions(projectID string, id string, definition ProjectConfigPatchDefinitionBlockIntf) *UpdateConfigOptions {
	return &UpdateConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Definition: definition,
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

// SetDefinition : Allow user to set Definition
func (_options *UpdateConfigOptions) SetDefinition(definition ProjectConfigPatchDefinitionBlockIntf) *UpdateConfigOptions {
	_options.Definition = definition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateConfigOptions) SetHeaders(param map[string]string) *UpdateConfigOptions {
	options.Headers = param
	return options
}

// UpdateProjectEnvironmentOptions : The UpdateProjectEnvironment options.
type UpdateProjectEnvironmentOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The environment ID.
	ID *string `json:"id" validate:"required,ne="`

	// The environment definition used for updates.
	Definition *EnvironmentDefinitionProperties `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProjectEnvironmentOptions : Instantiate UpdateProjectEnvironmentOptions
func (*ProjectV1) NewUpdateProjectEnvironmentOptions(projectID string, id string, definition *EnvironmentDefinitionProperties) *UpdateProjectEnvironmentOptions {
	return &UpdateProjectEnvironmentOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Definition: definition,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateProjectEnvironmentOptions) SetProjectID(projectID string) *UpdateProjectEnvironmentOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateProjectEnvironmentOptions) SetID(id string) *UpdateProjectEnvironmentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetDefinition : Allow user to set Definition
func (_options *UpdateProjectEnvironmentOptions) SetDefinition(definition *EnvironmentDefinitionProperties) *UpdateProjectEnvironmentOptions {
	_options.Definition = definition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProjectEnvironmentOptions) SetHeaders(param map[string]string) *UpdateProjectEnvironmentOptions {
	options.Headers = param
	return options
}

// UpdateProjectOptions : The UpdateProject options.
type UpdateProjectOptions struct {
	// The unique project ID.
	ID *string `json:"id" validate:"required,ne="`

	// The definition of the project.
	Definition *ProjectPatchDefinitionBlock `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProjectOptions : Instantiate UpdateProjectOptions
func (*ProjectV1) NewUpdateProjectOptions(id string, definition *ProjectPatchDefinitionBlock) *UpdateProjectOptions {
	return &UpdateProjectOptions{
		ID: core.StringPtr(id),
		Definition: definition,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateProjectOptions) SetID(id string) *UpdateProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetDefinition : Allow user to set Definition
func (_options *UpdateProjectOptions) SetDefinition(definition *ProjectPatchDefinitionBlock) *UpdateProjectOptions {
	_options.Definition = definition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProjectOptions) SetHeaders(param map[string]string) *UpdateProjectOptions {
	options.Headers = param
	return options
}

// ValidateConfigOptions : The ValidateConfig options.
type ValidateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique config ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewValidateConfigOptions : Instantiate ValidateConfigOptions
func (*ProjectV1) NewValidateConfigOptions(projectID string, id string) *ValidateConfigOptions {
	return &ValidateConfigOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ValidateConfigOptions) SetProjectID(projectID string) *ValidateConfigOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ValidateConfigOptions) SetID(id string) *ValidateConfigOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateConfigOptions) SetHeaders(param map[string]string) *ValidateConfigOptions {
	options.Headers = param
	return options
}

// ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigPatchDefinitionBlock
type ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`
}

func (*ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties) isaProjectConfigPatchDefinitionBlock() bool {
	return true
}

// UnmarshalProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties unmarshals an instance of ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchDefinitionBlockDAConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigPatchDefinitionBlock
type ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}

func (*ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties) isaProjectConfigPatchDefinitionBlock() bool {
	return true
}

// UnmarshalProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties unmarshals an instance of ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPatchDefinitionBlockResourceConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigPrototypeDefinitionBlock
type ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`
}

func (*ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties) isaProjectConfigPrototypeDefinitionBlock() bool {
	return true
}

// UnmarshalProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties unmarshals an instance of ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigPrototypeDefinitionBlock
type ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}

func (*ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties) isaProjectConfigPrototypeDefinitionBlock() bool {
	return true
}

// UnmarshalProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties unmarshals an instance of ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResponseDefinitionDAConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigResponseDefinition
type ProjectConfigResponseDefinitionDAConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The profile required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of catalogID.versionID that identifies the DA in the catalog. Either
	// schematics.workspace_crn, definition.locator_id, or both must be specified.
	LocatorID *string `json:"locator_id,omitempty"`
}

func (*ProjectConfigResponseDefinitionDAConfigDefinitionProperties) isaProjectConfigResponseDefinition() bool {
	return true
}

// UnmarshalProjectConfigResponseDefinitionDAConfigDefinitionProperties unmarshals an instance of ProjectConfigResponseDefinitionDAConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigResponseDefinitionDAConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResponseDefinitionDAConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResponseDefinitionResourceConfigDefinitionProperties : The name and description of a project configuration.
// This model "extends" ProjectConfigResponseDefinition
type ProjectConfigResponseDefinitionResourceConfigDefinitionProperties struct {
	// The configuration name. It is unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Schematics environment variables to use to deploy the configuration. Settings are only available if they were
	// specified when the configuration was initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The CRNs of resources associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}

func (*ProjectConfigResponseDefinitionResourceConfigDefinitionProperties) isaProjectConfigResponseDefinition() bool {
	return true
}

// UnmarshalProjectConfigResponseDefinitionResourceConfigDefinitionProperties unmarshals an instance of ProjectConfigResponseDefinitionResourceConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigResponseDefinitionResourceConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResponseDefinitionResourceConfigDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
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
func (pager *ProjectsPager) GetNextWithContext(ctx context.Context) (page []ProjectSummary, err error) {
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
		var start *string
		start, err = core.GetQueryParam(result.Next.Href, "start")
		if err != nil {
			err = fmt.Errorf("error retrieving 'start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			return
		}
		next = start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Projects

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProjectsPager) GetAllWithContext(ctx context.Context) (allItems []ProjectSummary, err error) {
	for pager.HasNext() {
		var nextPage []ProjectSummary
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetNext() (page []ProjectSummary, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetAll() (allItems []ProjectSummary, err error) {
	return pager.GetAllWithContext(context.Background())
}
