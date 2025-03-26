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
 * IBM OpenAPI SDK Code Generator Version: 3.92.1-44330004-20240620-143510
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
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	project, err = NewProjectV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = project.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = project.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
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

	service = &ProjectV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
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
	err := project.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
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
// For more information, see [Creating a
// project](/docs/secure-enterprise?topic=secure-enterprise-setup-project&interface=ui/docs-draft/secure-enterprise?topic=secure-enterprise-setup-project).
func (project *ProjectV1) CreateProject(createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	result, response, err = project.CreateProjectWithContext(context.Background(), createProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProjectWithContext is an alternate form of the CreateProject method which supports a Context parameter
func (project *ProjectV1) CreateProjectWithContext(ctx context.Context, createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProjectOptions, "createProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProjectOptions, "createProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListProjects : List projects
// List existing projects. Projects are sorted by ID.
func (project *ProjectV1) ListProjects(listProjectsOptions *ListProjectsOptions) (result *ProjectCollection, response *core.DetailedResponse, err error) {
	result, response, err = project.ListProjectsWithContext(context.Background(), listProjectsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProjectsWithContext is an alternate form of the ListProjects method which supports a Context parameter
func (project *ProjectV1) ListProjectsWithContext(ctx context.Context, listProjectsOptions *ListProjectsOptions) (result *ProjectCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProjectsOptions, "listProjectsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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

	if listProjectsOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listProjectsOptions.Token))
	}
	if listProjectsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProjectsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_projects", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProject : Get a project
// Get information about a project.
func (project *ProjectV1) GetProject(getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	result, response, err = project.GetProjectWithContext(context.Background(), getProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProjectWithContext is an alternate form of the GetProject method which supports a Context parameter
func (project *ProjectV1) GetProjectWithContext(ctx context.Context, getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectOptions, "getProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProjectOptions, "getProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateProject : Update a project
// Update a project by specifying its ID.
func (project *ProjectV1) UpdateProject(updateProjectOptions *UpdateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	result, response, err = project.UpdateProjectWithContext(context.Background(), updateProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateProjectWithContext is an alternate form of the UpdateProject method which supports a Context parameter
func (project *ProjectV1) UpdateProjectWithContext(ctx context.Context, updateProjectOptions *UpdateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProjectOptions, "updateProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateProjectOptions, "updateProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProject : Delete a project
// Delete a project document by specifying the ID. A project can be deleted only after you delete all of its resources.
func (project *ProjectV1) DeleteProject(deleteProjectOptions *DeleteProjectOptions) (result *ProjectDeleteResponse, response *core.DetailedResponse, err error) {
	result, response, err = project.DeleteProjectWithContext(context.Background(), deleteProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProjectWithContext is an alternate form of the DeleteProject method which supports a Context parameter
func (project *ProjectV1) DeleteProjectWithContext(ctx context.Context, deleteProjectOptions *DeleteProjectOptions) (result *ProjectDeleteResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProjectOptions, "deleteProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProjectOptions, "deleteProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "DeleteProject")
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectDeleteResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateProjectEnvironment : Create an environment
// Create an environment to group related configurations together and share values across them for easier deployment.
// For more information, see [Creating an environment](/docs/secure-enterprise?topic=secure-enterprise-create-env).
func (project *ProjectV1) CreateProjectEnvironment(createProjectEnvironmentOptions *CreateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = project.CreateProjectEnvironmentWithContext(context.Background(), createProjectEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProjectEnvironmentWithContext is an alternate form of the CreateProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) CreateProjectEnvironmentWithContext(ctx context.Context, createProjectEnvironmentOptions *CreateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProjectEnvironmentOptions, "createProjectEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProjectEnvironmentOptions, "createProjectEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_project_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListProjectEnvironments : List environments
// List all available environments. For more information, see [Creating an
// environment](/docs/secure-enterprise?topic=secure-enterprise-create-env).
func (project *ProjectV1) ListProjectEnvironments(listProjectEnvironmentsOptions *ListProjectEnvironmentsOptions) (result *EnvironmentCollection, response *core.DetailedResponse, err error) {
	result, response, err = project.ListProjectEnvironmentsWithContext(context.Background(), listProjectEnvironmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProjectEnvironmentsWithContext is an alternate form of the ListProjectEnvironments method which supports a Context parameter
func (project *ProjectV1) ListProjectEnvironmentsWithContext(ctx context.Context, listProjectEnvironmentsOptions *ListProjectEnvironmentsOptions) (result *EnvironmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProjectEnvironmentsOptions, "listProjectEnvironmentsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listProjectEnvironmentsOptions, "listProjectEnvironmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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

	if listProjectEnvironmentsOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listProjectEnvironmentsOptions.Token))
	}
	if listProjectEnvironmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProjectEnvironmentsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_project_environments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProjectEnvironment : Get an environment
// Get an environment. [Learn more](/docs/secure-enterprise?topic=secure-enterprise-create-env).
func (project *ProjectV1) GetProjectEnvironment(getProjectEnvironmentOptions *GetProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = project.GetProjectEnvironmentWithContext(context.Background(), getProjectEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProjectEnvironmentWithContext is an alternate form of the GetProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) GetProjectEnvironmentWithContext(ctx context.Context, getProjectEnvironmentOptions *GetProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectEnvironmentOptions, "getProjectEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProjectEnvironmentOptions, "getProjectEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_project_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateProjectEnvironment : Update an environment
// Update an environment by specifying its ID. [Learn more](/docs/secure-enterprise?topic=secure-enterprise-create-env).
func (project *ProjectV1) UpdateProjectEnvironment(updateProjectEnvironmentOptions *UpdateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = project.UpdateProjectEnvironmentWithContext(context.Background(), updateProjectEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateProjectEnvironmentWithContext is an alternate form of the UpdateProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) UpdateProjectEnvironmentWithContext(ctx context.Context, updateProjectEnvironmentOptions *UpdateProjectEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProjectEnvironmentOptions, "updateProjectEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateProjectEnvironmentOptions, "updateProjectEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_project_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProjectEnvironment : Delete an environment
// Delete an environment in a project by specifying its ID.
func (project *ProjectV1) DeleteProjectEnvironment(deleteProjectEnvironmentOptions *DeleteProjectEnvironmentOptions) (result *EnvironmentDeleteResponse, response *core.DetailedResponse, err error) {
	result, response, err = project.DeleteProjectEnvironmentWithContext(context.Background(), deleteProjectEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProjectEnvironmentWithContext is an alternate form of the DeleteProjectEnvironment method which supports a Context parameter
func (project *ProjectV1) DeleteProjectEnvironmentWithContext(ctx context.Context, deleteProjectEnvironmentOptions *DeleteProjectEnvironmentOptions) (result *EnvironmentDeleteResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProjectEnvironmentOptions, "deleteProjectEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProjectEnvironmentOptions, "deleteProjectEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_project_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentDeleteResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateConfig : Add a new configuration
// Add a new configuration to a project.
func (project *ProjectV1) CreateConfig(createConfigOptions *CreateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	result, response, err = project.CreateConfigWithContext(context.Background(), createConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateConfigWithContext is an alternate form of the CreateConfig method which supports a Context parameter
func (project *ProjectV1) CreateConfigWithContext(ctx context.Context, createConfigOptions *CreateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigOptions, "createConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createConfigOptions, "createConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListConfigs : List all project configurations
// Retrieve the collection of configurations.
func (project *ProjectV1) ListConfigs(listConfigsOptions *ListConfigsOptions) (result *ProjectConfigCollection, response *core.DetailedResponse, err error) {
	result, response, err = project.ListConfigsWithContext(context.Background(), listConfigsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListConfigsWithContext is an alternate form of the ListConfigs method which supports a Context parameter
func (project *ProjectV1) ListConfigsWithContext(ctx context.Context, listConfigsOptions *ListConfigsOptions) (result *ProjectConfigCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigsOptions, "listConfigsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listConfigsOptions, "listConfigsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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

	if listConfigsOptions.Token != nil {
		builder.AddQuery("token", fmt.Sprint(*listConfigsOptions.Token))
	}
	if listConfigsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listConfigsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_configs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetConfig : Get a project configuration
// Retrieve the specified project configuration in a specific project. For more information about project
// configurations, see [Monitoring the status of a configuration and its
// resources](/docs/secure-enterprise?topic=secure-enterprise-monitor-status-projects).
func (project *ProjectV1) GetConfig(getConfigOptions *GetConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	result, response, err = project.GetConfigWithContext(context.Background(), getConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetConfigWithContext is an alternate form of the GetConfig method which supports a Context parameter
func (project *ProjectV1) GetConfigWithContext(ctx context.Context, getConfigOptions *GetConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigOptions, "getConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getConfigOptions, "getConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateConfig : Update a configuration
// Update a configuration in a project by specifying the ID. [Learn
// more](/docs/secure-enterprise?topic=secure-enterprise-config-project).
func (project *ProjectV1) UpdateConfig(updateConfigOptions *UpdateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	result, response, err = project.UpdateConfigWithContext(context.Background(), updateConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateConfigWithContext is an alternate form of the UpdateConfig method which supports a Context parameter
func (project *ProjectV1) UpdateConfigWithContext(ctx context.Context, updateConfigOptions *UpdateConfigOptions) (result *ProjectConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateConfigOptions, "updateConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateConfigOptions, "updateConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfig : Delete a configuration
// Delete a configuration in a project by specifying its ID.
func (project *ProjectV1) DeleteConfig(deleteConfigOptions *DeleteConfigOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	result, response, err = project.DeleteConfigWithContext(context.Background(), deleteConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteConfigWithContext is an alternate form of the DeleteConfig method which supports a Context parameter
func (project *ProjectV1) DeleteConfigWithContext(ctx context.Context, deleteConfigOptions *DeleteConfigOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigOptions, "deleteConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteConfigOptions, "deleteConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigDelete)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ForceApprove : Force approve a project configuration
// Force approve configuration edits to the main configuration with an approving comment.
func (project *ProjectV1) ForceApprove(forceApproveOptions *ForceApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.ForceApproveWithContext(context.Background(), forceApproveOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ForceApproveWithContext is an alternate form of the ForceApprove method which supports a Context parameter
func (project *ProjectV1) ForceApproveWithContext(ctx context.Context, forceApproveOptions *ForceApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(forceApproveOptions, "forceApproveOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(forceApproveOptions, "forceApproveOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "force_approve", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Approve : Approve and merge a configuration draft
// Approve and merge configuration edits to the main configuration.
func (project *ProjectV1) Approve(approveOptions *ApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.ApproveWithContext(context.Background(), approveOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ApproveWithContext is an alternate form of the Approve method which supports a Context parameter
func (project *ProjectV1) ApproveWithContext(ctx context.Context, approveOptions *ApproveOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(approveOptions, "approveOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(approveOptions, "approveOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "approve", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ValidateConfig : Run a validation check
// Run a validation check on a specific configuration in the project. The check includes creating or updating the
// associated Schematics workspace with a plan job, running the CRA scans, and cost estimation.
func (project *ProjectV1) ValidateConfig(validateConfigOptions *ValidateConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.ValidateConfigWithContext(context.Background(), validateConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ValidateConfigWithContext is an alternate form of the ValidateConfig method which supports a Context parameter
func (project *ProjectV1) ValidateConfigWithContext(ctx context.Context, validateConfigOptions *ValidateConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateConfigOptions, "validateConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(validateConfigOptions, "validateConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "validate_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeployConfig : Deploy a configuration
// Deploy a project's configuration. This operation is asynchronous and can be tracked by using the get project
// configuration API with full metadata.
func (project *ProjectV1) DeployConfig(deployConfigOptions *DeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.DeployConfigWithContext(context.Background(), deployConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeployConfigWithContext is an alternate form of the DeployConfig method which supports a Context parameter
func (project *ProjectV1) DeployConfigWithContext(ctx context.Context, deployConfigOptions *DeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deployConfigOptions, "deployConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deployConfigOptions, "deployConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "deploy_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UndeployConfig : Undeploy configuration resources
// Undeploy a project's configuration resources. The operation undeploys all the resources that are deployed with the
// specific configuration. You can track it by using the get project configuration API with full metadata.
func (project *ProjectV1) UndeployConfig(undeployConfigOptions *UndeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.UndeployConfigWithContext(context.Background(), undeployConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UndeployConfigWithContext is an alternate form of the UndeployConfig method which supports a Context parameter
func (project *ProjectV1) UndeployConfigWithContext(ctx context.Context, undeployConfigOptions *UndeployConfigOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(undeployConfigOptions, "undeployConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(undeployConfigOptions, "undeployConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range undeployConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UndeployConfig")
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "undeploy_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SyncConfig : Sync a project configuration
// Sync a project configuration by analyzing the associated pipeline runs and Schematics workspace logs to get the
// configuration back to a working state.
func (project *ProjectV1) SyncConfig(syncConfigOptions *SyncConfigOptions) (response *core.DetailedResponse, err error) {
	response, err = project.SyncConfigWithContext(context.Background(), syncConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SyncConfigWithContext is an alternate form of the SyncConfig method which supports a Context parameter
func (project *ProjectV1) SyncConfigWithContext(ctx context.Context, syncConfigOptions *SyncConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(syncConfigOptions, "syncConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(syncConfigOptions, "syncConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = project.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "sync_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListConfigResources : List all deployed resources
// List resources that are deployed by a configuration.
func (project *ProjectV1) ListConfigResources(listConfigResourcesOptions *ListConfigResourcesOptions) (result *ProjectConfigResourceCollection, response *core.DetailedResponse, err error) {
	result, response, err = project.ListConfigResourcesWithContext(context.Background(), listConfigResourcesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListConfigResourcesWithContext is an alternate form of the ListConfigResources method which supports a Context parameter
func (project *ProjectV1) ListConfigResourcesWithContext(ctx context.Context, listConfigResourcesOptions *ListConfigResourcesOptions) (result *ProjectConfigResourceCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigResourcesOptions, "listConfigResourcesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listConfigResourcesOptions, "listConfigResourcesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_config_resources", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigResourceCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateStackDefinition : Create a stack definition
// Defines inputs at the stack level that users need to configure along with input values at the member level. These
// values are included in the catalog entry when the deployable architecture stack is exported to a private catalog and
// are required for the deployable architecture stack to deploy. You can add a reference to a value, or add the value
// explicitly at the member level.
func (project *ProjectV1) CreateStackDefinition(createStackDefinitionOptions *CreateStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	result, response, err = project.CreateStackDefinitionWithContext(context.Background(), createStackDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateStackDefinitionWithContext is an alternate form of the CreateStackDefinition method which supports a Context parameter
func (project *ProjectV1) CreateStackDefinitionWithContext(ctx context.Context, createStackDefinitionOptions *CreateStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createStackDefinitionOptions, "createStackDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createStackDefinitionOptions, "createStackDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createStackDefinitionOptions.ProjectID,
		"id": *createStackDefinitionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/stack_definition`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createStackDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "CreateStackDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createStackDefinitionOptions.StackDefinition != nil {
		body["stack_definition"] = createStackDefinitionOptions.StackDefinition
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_stack_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStackDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetStackDefinition : Get a stack definition
// Retrieve the stack definition that is associated to the configuration.
func (project *ProjectV1) GetStackDefinition(getStackDefinitionOptions *GetStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	result, response, err = project.GetStackDefinitionWithContext(context.Background(), getStackDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetStackDefinitionWithContext is an alternate form of the GetStackDefinition method which supports a Context parameter
func (project *ProjectV1) GetStackDefinitionWithContext(ctx context.Context, getStackDefinitionOptions *GetStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getStackDefinitionOptions, "getStackDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getStackDefinitionOptions, "getStackDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getStackDefinitionOptions.ProjectID,
		"id": *getStackDefinitionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/stack_definition`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getStackDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "GetStackDefinition")
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_stack_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStackDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateStackDefinition : Update a stack definition
// Update the stack definition that is associated with the configuration.
func (project *ProjectV1) UpdateStackDefinition(updateStackDefinitionOptions *UpdateStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	result, response, err = project.UpdateStackDefinitionWithContext(context.Background(), updateStackDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateStackDefinitionWithContext is an alternate form of the UpdateStackDefinition method which supports a Context parameter
func (project *ProjectV1) UpdateStackDefinitionWithContext(ctx context.Context, updateStackDefinitionOptions *UpdateStackDefinitionOptions) (result *StackDefinition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateStackDefinitionOptions, "updateStackDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateStackDefinitionOptions, "updateStackDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateStackDefinitionOptions.ProjectID,
		"id": *updateStackDefinitionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/stack_definition`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateStackDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "UpdateStackDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateStackDefinitionOptions.StackDefinition != nil {
		body["stack_definition"] = updateStackDefinitionOptions.StackDefinition
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_stack_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStackDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ExportStackDefinition : Export a deployable architecture stack to the private catalog
// Exports the deployable architecture stack to a private catalog. All member deployable architectures within the stack
// must be validated and deployed before the stack is exported. The stack definition must also exist before the stack is
// exported. You can export the stack as a new product, or as a new version of an existing product.
func (project *ProjectV1) ExportStackDefinition(exportStackDefinitionOptions *ExportStackDefinitionOptions) (result *StackDefinitionExportResponse, response *core.DetailedResponse, err error) {
	result, response, err = project.ExportStackDefinitionWithContext(context.Background(), exportStackDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ExportStackDefinitionWithContext is an alternate form of the ExportStackDefinition method which supports a Context parameter
func (project *ProjectV1) ExportStackDefinitionWithContext(ctx context.Context, exportStackDefinitionOptions *ExportStackDefinitionOptions) (result *StackDefinitionExportResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(exportStackDefinitionOptions, "exportStackDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(exportStackDefinitionOptions, "exportStackDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *exportStackDefinitionOptions.ProjectID,
		"id": *exportStackDefinitionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = project.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(project.Service.Options.URL, `/v1/projects/{project_id}/configs/{id}/stack_definition/export`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range exportStackDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("project", "V1", "ExportStackDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(exportStackDefinitionOptions.Settings)
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
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "export_stack_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStackDefinitionExportResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListConfigVersions : Get a list of project configuration versions
// Retrieve a list of previous and current versions of a project configuration in a specific project.
func (project *ProjectV1) ListConfigVersions(listConfigVersionsOptions *ListConfigVersionsOptions) (result *ProjectConfigVersionSummaryCollection, response *core.DetailedResponse, err error) {
	result, response, err = project.ListConfigVersionsWithContext(context.Background(), listConfigVersionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListConfigVersionsWithContext is an alternate form of the ListConfigVersions method which supports a Context parameter
func (project *ProjectV1) ListConfigVersionsWithContext(ctx context.Context, listConfigVersionsOptions *ListConfigVersionsOptions) (result *ProjectConfigVersionSummaryCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigVersionsOptions, "listConfigVersionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listConfigVersionsOptions, "listConfigVersionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_config_versions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersionSummaryCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetConfigVersion : Get a specific project configuration version
// Retrieve a specific version of a configuration in a project.
func (project *ProjectV1) GetConfigVersion(getConfigVersionOptions *GetConfigVersionOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	result, response, err = project.GetConfigVersionWithContext(context.Background(), getConfigVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetConfigVersionWithContext is an alternate form of the GetConfigVersion method which supports a Context parameter
func (project *ProjectV1) GetConfigVersionWithContext(ctx context.Context, getConfigVersionOptions *GetConfigVersionOptions) (result *ProjectConfigVersion, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigVersionOptions, "getConfigVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getConfigVersionOptions, "getConfigVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_config_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigVersion)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfigVersion : Delete a project configuration version
// Delete a configuration version by specifying the project ID.
func (project *ProjectV1) DeleteConfigVersion(deleteConfigVersionOptions *DeleteConfigVersionOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	result, response, err = project.DeleteConfigVersionWithContext(context.Background(), deleteConfigVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteConfigVersionWithContext is an alternate form of the DeleteConfigVersion method which supports a Context parameter
func (project *ProjectV1) DeleteConfigVersionWithContext(ctx context.Context, deleteConfigVersionOptions *DeleteConfigVersionOptions) (result *ProjectConfigDelete, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigVersionOptions, "deleteConfigVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteConfigVersionOptions, "deleteConfigVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = project.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_config_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectConfigDelete)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// ActionJobApplyMessagesSummary : The messages of apply jobs on the configuration.
type ActionJobApplyMessagesSummary struct {
	// The collection of error messages. This property is reported only if Schematics triggered a Terraform apply job.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`

	// The collection of success messages. This property is reported only if Schematics triggered a Terraform apply job.
	SuccessMessages []TerraformLogAnalyzerSuccessMessage `json:"success_messages,omitempty"`
}

// UnmarshalActionJobApplyMessagesSummary unmarshals an instance of ActionJobApplyMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobApplyMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobApplyMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "success_messages", &obj.SuccessMessages, UnmarshalTerraformLogAnalyzerSuccessMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "success_messages-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobApplySummary : The summary of the apply jobs on the configuration.
type ActionJobApplySummary struct {
	// The number of applied resources. This property is reported only if Schematics triggered a Terraform apply job.
	Success *int64 `json:"success,omitempty"`

	// The number of failed applied resources. This property is reported only if Schematics triggered a Terraform apply
	// job.
	Failed *int64 `json:"failed,omitempty"`

	// The collection of successfully applied resources. This property is reported only if Schematics triggered a Terraform
	// apply job.
	SuccessResources []string `json:"success_resources,omitempty"`

	// The collection of failed applied resources. This property is reported only if Schematics triggered a Terraform apply
	// job.
	FailedResources []string `json:"failed_resources,omitempty"`
}

// UnmarshalActionJobApplySummary unmarshals an instance of ActionJobApplySummary from the specified map of raw messages.
func UnmarshalActionJobApplySummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobApplySummary)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "success_resources", &obj.SuccessResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "success_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_resources", &obj.FailedResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed_resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroyMessagesSummary : The messages of destroy jobs on the configuration.
type ActionJobDestroyMessagesSummary struct {
	// The collection of error messages. This property is reported only if Schematics triggered a Terraform destroy job.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`
}

// UnmarshalActionJobDestroyMessagesSummary unmarshals an instance of ActionJobDestroyMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobDestroyMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroyMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_messages-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroySummary : The summary of the destroy jobs on the configuration.
type ActionJobDestroySummary struct {
	// The number of destroyed resources. This property is reported only if Schematics triggered a Terraform destroy job.
	Success *int64 `json:"success,omitempty"`

	// The number of failed resources. This property is reported only if Schematics triggered a Terraform destroy job.
	Failed *int64 `json:"failed,omitempty"`

	// The number of tainted resources. This property is reported only if Schematics triggered a Terraform destroy job.
	Tainted *int64 `json:"tainted,omitempty"`

	// The summary of results from destroyed resources from the job. This property is reported only if Schematics triggered
	// a Terraform destroy job.
	Resources *ActionJobDestroySummaryResources `json:"resources,omitempty"`
}

// UnmarshalActionJobDestroySummary unmarshals an instance of ActionJobDestroySummary from the specified map of raw messages.
func UnmarshalActionJobDestroySummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroySummary)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tainted", &obj.Tainted)
	if err != nil {
		err = core.SDKErrorf(err, "", "tainted-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalActionJobDestroySummaryResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobDestroySummaryResources : The summary of results from destroyed resources from the job. This property is reported only if Schematics triggered
// a Terraform destroy job.
type ActionJobDestroySummaryResources struct {
	// The collection of destroyed resources. This property is reported only if Schematics triggered a Terraform destroy
	// job.
	Success []string `json:"success,omitempty"`

	// The collection of failed resources. This property is reported only if Schematics triggered a Terraform destroy job.
	Failed []string `json:"failed,omitempty"`

	// The collection of tainted resources. This property is reported only if Schematics triggered a Terraform destroy job.
	Tainted []string `json:"tainted,omitempty"`
}

// UnmarshalActionJobDestroySummaryResources unmarshals an instance of ActionJobDestroySummaryResources from the specified map of raw messages.
func UnmarshalActionJobDestroySummaryResources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobDestroySummaryResources)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tainted", &obj.Tainted)
	if err != nil {
		err = core.SDKErrorf(err, "", "tainted-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobMessageSummary : The message summaries of jobs on the configuration.
type ActionJobMessageSummary struct {
	// The number of information messages. This property is reported only if Schematics triggered a Terraform job.
	Info *int64 `json:"info,omitempty"`

	// The number of debug messages. This property is reported only if Schematics triggered a Terraform job.
	Debug *int64 `json:"debug,omitempty"`

	// The number of error messages. This property is reported only if Schematics triggered a Terraform job.
	Error *int64 `json:"error,omitempty"`
}

// UnmarshalActionJobMessageSummary unmarshals an instance of ActionJobMessageSummary from the specified map of raw messages.
func UnmarshalActionJobMessageSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobMessageSummary)
	err = core.UnmarshalPrimitive(m, "info", &obj.Info)
	if err != nil {
		err = core.SDKErrorf(err, "", "info-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "debug", &obj.Debug)
	if err != nil {
		err = core.SDKErrorf(err, "", "debug-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "error", &obj.Error)
	if err != nil {
		err = core.SDKErrorf(err, "", "error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobPlanMessagesSummary : The plan messages on the configuration.
type ActionJobPlanMessagesSummary struct {
	// The collection of error messages. This property is reported only if Schematics triggered a Terraform plan job.
	ErrorMessages []TerraformLogAnalyzerErrorMessage `json:"error_messages,omitempty"`

	// The collection of success messages. This property is reported only if Schematics triggered a Terraform plan job.
	SuccessMessages []string `json:"success_messages,omitempty"`

	// The collection of update messages. This property is reported only if Schematics triggered a Terraform plan job.
	UpdateMessages []string `json:"update_messages,omitempty"`

	// The collection of destroy messages. This property is reported only if Schematics triggered a Terraform plan job.
	DestroyMessages []string `json:"destroy_messages,omitempty"`
}

// UnmarshalActionJobPlanMessagesSummary unmarshals an instance of ActionJobPlanMessagesSummary from the specified map of raw messages.
func UnmarshalActionJobPlanMessagesSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobPlanMessagesSummary)
	err = core.UnmarshalModel(m, "error_messages", &obj.ErrorMessages, UnmarshalTerraformLogAnalyzerErrorMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "success_messages", &obj.SuccessMessages)
	if err != nil {
		err = core.SDKErrorf(err, "", "success_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "update_messages", &obj.UpdateMessages)
	if err != nil {
		err = core.SDKErrorf(err, "", "update_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_messages", &obj.DestroyMessages)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_messages-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobPlanSummary : The summary of the plan jobs on the configuration.
type ActionJobPlanSummary struct {
	// The number of resources to be added. This property is reported only if Schematics triggered a terraform plan job.
	Add *int64 `json:"add,omitempty"`

	// The number of resources that failed during the plan job. This property is reported only if Schematics triggered a
	// Terraform plan job.
	Failed *int64 `json:"failed,omitempty"`

	// The number of resources to be updated. This property is reported only if Schematics triggered a Terraform plan job.
	Update *int64 `json:"update,omitempty"`

	// The number of resources to be destroyed. This property is reported only if Schematics triggered a Terraform plan
	// job.
	Destroy *int64 `json:"destroy,omitempty"`

	// The collection of planned added resources. This property is reported only if Schematics triggered a Terraform plan
	// job.
	AddResources []string `json:"add_resources,omitempty"`

	// The collection of failed planned resources. This property is reported only if Schematics triggered a Terraform plan
	// job.
	FailedResources []string `json:"failed_resources,omitempty"`

	// The collection of planned updated resources. This property is reported only if Schematics triggered a Terraform plan
	// job.
	UpdatedResources []string `json:"updated_resources,omitempty"`

	// The collection of planned destroy resources. This property is reported only if Schematics triggered a Terraform plan
	// job.
	DestroyResources []string `json:"destroy_resources,omitempty"`
}

// UnmarshalActionJobPlanSummary unmarshals an instance of ActionJobPlanSummary from the specified map of raw messages.
func UnmarshalActionJobPlanSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionJobPlanSummary)
	err = core.UnmarshalPrimitive(m, "add", &obj.Add)
	if err != nil {
		err = core.SDKErrorf(err, "", "add-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "update", &obj.Update)
	if err != nil {
		err = core.SDKErrorf(err, "", "update-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy", &obj.Destroy)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "add_resources", &obj.AddResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "add_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_resources", &obj.FailedResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_resources", &obj.UpdatedResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_resources", &obj.DestroyResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ActionJobSummary : The summaries of jobs that were performed on the configuration.
type ActionJobSummary struct {
	// The version of the job summary.
	Version *string `json:"version" validate:"required"`

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
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plan_summary", &obj.PlanSummary, UnmarshalActionJobPlanSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "apply_summary", &obj.ApplySummary, UnmarshalActionJobApplySummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "apply_summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "destroy_summary", &obj.DestroySummary, UnmarshalActionJobDestroySummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "message_summary", &obj.MessageSummary, UnmarshalActionJobMessageSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "message_summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plan_messages", &obj.PlanMessages, UnmarshalActionJobPlanMessagesSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "apply_messages", &obj.ApplyMessages, UnmarshalActionJobApplyMessagesSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "apply_messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "destroy_messages", &obj.DestroyMessages, UnmarshalActionJobDestroyMessagesSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_messages-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalActionJobSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "summary-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApproveOptions : The Approve options.
type ApproveOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// Notes on the project draft action. If this action is a force approve on the draft configuration, you must include a
	// nonempty comment.
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

// CodeRiskAnalyzerLogsSummary : The Code Risk Analyzer logs a summary of the configuration.
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
		err = core.SDKErrorf(err, "", "total-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "passed", &obj.Passed)
	if err != nil {
		err = core.SDKErrorf(err, "", "passed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped", &obj.Skipped)
	if err != nil {
		err = core.SDKErrorf(err, "", "skipped-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigDefinitionReference : The definition of the config reference.
type ConfigDefinitionReference struct {
	// The name of the configuration.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalConfigDefinitionReference unmarshals an instance of ConfigDefinitionReference from the specified map of raw messages.
func UnmarshalConfigDefinitionReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigDefinitionReference)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateConfigOptions : The CreateConfig options.
type CreateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	Definition ProjectConfigDefinitionPrototypeIntf `json:"definition" validate:"required"`

	// A Schematics workspace to use for deploying this deployable architecture.
	// > If you are importing data from an existing Schematics workspace that is not backed by cart, then you must provide
	// a `locator_id`. If you are using a Schematics workspace that is backed by cart, a `locator_id` is not required
	// because the Schematics workspace has one.
	// >
	// There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400`code  is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// >
	// For more information, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	Schematics *SchematicsWorkspace `json:"schematics,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateConfigOptions : Instantiate CreateConfigOptions
func (*ProjectV1) NewCreateConfigOptions(projectID string, definition ProjectConfigDefinitionPrototypeIntf) *CreateConfigOptions {
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
func (_options *CreateConfigOptions) SetDefinition(definition ProjectConfigDefinitionPrototypeIntf) *CreateConfigOptions {
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

	// The project configurations. These configurations are included in the response of creating a project only if a
	// configuration array is specified in the request payload.
	Configs []ProjectConfigPrototype `json:"configs,omitempty"`

	// The project environment. These environments are included in the response of creating a project only if an
	// environment array is specified in the request payload.
	Environments []EnvironmentPrototype `json:"environments,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateProjectOptions.Location property.
// The IBM Cloud location where a resource is deployed.
const (
	CreateProjectOptions_Location_CaTor = "ca-tor"
	CreateProjectOptions_Location_EuDe = "eu-de"
	CreateProjectOptions_Location_EuGb = "eu-gb"
	CreateProjectOptions_Location_UsEast = "us-east"
	CreateProjectOptions_Location_UsSouth = "us-south"
)

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

// CreateStackDefinitionOptions : The CreateStackDefinition options.
type CreateStackDefinitionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// The definition block for a stack definition.
	StackDefinition *StackDefinitionBlockPrototype `json:"stack_definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateStackDefinitionOptions : Instantiate CreateStackDefinitionOptions
func (*ProjectV1) NewCreateStackDefinitionOptions(projectID string, id string, stackDefinition *StackDefinitionBlockPrototype) *CreateStackDefinitionOptions {
	return &CreateStackDefinitionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		StackDefinition: stackDefinition,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateStackDefinitionOptions) SetProjectID(projectID string) *CreateStackDefinitionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateStackDefinitionOptions) SetID(id string) *CreateStackDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetStackDefinition : Allow user to set StackDefinition
func (_options *CreateStackDefinitionOptions) SetStackDefinition(stackDefinition *StackDefinitionBlockPrototype) *CreateStackDefinitionOptions {
	_options.StackDefinition = stackDefinition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateStackDefinitionOptions) SetHeaders(param map[string]string) *CreateStackDefinitionOptions {
	options.Headers = param
	return options
}

// CumulativeNeedsAttention : CumulativeNeedsAttention struct
type CumulativeNeedsAttention struct {
	// The event name.
	Event *string `json:"event,omitempty"`

	// A unique ID for this individual event.
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
		err = core.SDKErrorf(err, "", "event-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_id", &obj.EventID)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "config_id", &obj.ConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "config_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "config_version", &obj.ConfigVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "config_version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteConfigOptions : The DeleteConfig options.
type DeleteConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
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

	// The unique configuration ID.
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

	// The unique configuration ID.
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
	// The environment ID as a friendly name.
	ID *string `json:"id" validate:"required"`

	// The project that is referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The target account ID derived from the authentication block values. The target account exists only if the
	// environment currently has an authorization block.
	TargetAccount *string `json:"target_account,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The environment definition.
	Definition *EnvironmentDefinitionRequiredPropertiesResponse `json:"definition" validate:"required"`
}

// UnmarshalEnvironment unmarshals an instance of Environment from the specified map of raw messages.
func UnmarshalEnvironment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Environment)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_account", &obj.TargetAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalEnvironmentDefinitionRequiredPropertiesResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentCollection : The list environment response.
type EnvironmentCollection struct {
	// A pagination limit.
	Limit *int64 `json:"limit" validate:"required"`

	// A pagination link.
	First *PaginationLink `json:"first" validate:"required"`

	// A pagination link.
	Next *PaginationLink `json:"next,omitempty"`

	// The environment.
	Environments []Environment `json:"environments,omitempty"`
}

// UnmarshalEnvironmentCollection unmarshals an instance of EnvironmentCollection from the specified map of raw messages.
func UnmarshalEnvironmentCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalEnvironment)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *EnvironmentCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	token, err := core.GetQueryParam(resp.Next.Href, "token")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if token == nil {
		return nil, nil
	}
	return token, nil
}

// EnvironmentDefinitionPropertiesPatch : The environment definition that is used for updates.
type EnvironmentDefinitionPropertiesPatch struct {
	// The description of the environment.
	Description *string `json:"description,omitempty"`

	// The name of the environment. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`
}

// UnmarshalEnvironmentDefinitionPropertiesPatch unmarshals an instance of EnvironmentDefinitionPropertiesPatch from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionPropertiesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionPropertiesPatch)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDefinitionRequiredProperties : The environment definition.
type EnvironmentDefinitionRequiredProperties struct {
	// The description of the environment.
	Description *string `json:"description,omitempty"`

	// The name of the environment. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`
}

// NewEnvironmentDefinitionRequiredProperties : Instantiate EnvironmentDefinitionRequiredProperties (Generic Model Constructor)
func (*ProjectV1) NewEnvironmentDefinitionRequiredProperties(name string) (_model *EnvironmentDefinitionRequiredProperties, err error) {
	_model = &EnvironmentDefinitionRequiredProperties{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalEnvironmentDefinitionRequiredProperties unmarshals an instance of EnvironmentDefinitionRequiredProperties from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionRequiredProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionRequiredProperties)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDefinitionRequiredPropertiesResponse : The environment definition.
type EnvironmentDefinitionRequiredPropertiesResponse struct {
	// The description of the environment.
	Description *string `json:"description" validate:"required"`

	// The name of the environment. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`
}

// UnmarshalEnvironmentDefinitionRequiredPropertiesResponse unmarshals an instance of EnvironmentDefinitionRequiredPropertiesResponse from the specified map of raw messages.
func UnmarshalEnvironmentDefinitionRequiredPropertiesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDefinitionRequiredPropertiesResponse)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentDeleteResponse : The response to a request to delete an environment.
type EnvironmentDeleteResponse struct {
	// The environment ID as a friendly name.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalEnvironmentDeleteResponse unmarshals an instance of EnvironmentDeleteResponse from the specified map of raw messages.
func UnmarshalEnvironmentDeleteResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentDeleteResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalEnvironmentPrototype unmarshals an instance of EnvironmentPrototype from the specified map of raw messages.
func UnmarshalEnvironmentPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentPrototype)
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalEnvironmentDefinitionRequiredProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportStackDefinitionOptions : The ExportStackDefinition options.
type ExportStackDefinitionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// The payload for the private catalog export request.
	Settings StackDefinitionExportRequestIntf `json:"settings" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewExportStackDefinitionOptions : Instantiate ExportStackDefinitionOptions
func (*ProjectV1) NewExportStackDefinitionOptions(projectID string, id string, settings StackDefinitionExportRequestIntf) *ExportStackDefinitionOptions {
	return &ExportStackDefinitionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		Settings: settings,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ExportStackDefinitionOptions) SetProjectID(projectID string) *ExportStackDefinitionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ExportStackDefinitionOptions) SetID(id string) *ExportStackDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *ExportStackDefinitionOptions) SetSettings(settings StackDefinitionExportRequestIntf) *ExportStackDefinitionOptions {
	_options.Settings = settings
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ExportStackDefinitionOptions) SetHeaders(param map[string]string) *ExportStackDefinitionOptions {
	options.Headers = param
	return options
}

// ForceApproveOptions : The ForceApprove options.
type ForceApproveOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// Notes on the project draft action. If this action is a force approve on the draft configuration, you must include a
	// nonempty comment.
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

	// The unique configuration ID.
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

	// The unique configuration ID.
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

// GetStackDefinitionOptions : The GetStackDefinition options.
type GetStackDefinitionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetStackDefinitionOptions : Instantiate GetStackDefinitionOptions
func (*ProjectV1) NewGetStackDefinitionOptions(projectID string, id string) *GetStackDefinitionOptions {
	return &GetStackDefinitionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetStackDefinitionOptions) SetProjectID(projectID string) *GetStackDefinitionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetStackDefinitionOptions) SetID(id string) *GetStackDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetStackDefinitionOptions) SetHeaders(param map[string]string) *GetStackDefinitionOptions {
	options.Headers = param
	return options
}

// LastActionWithSummary : The href and results from the last action job that is performed on the project configuration.
type LastActionWithSummary struct {
	// A URL.
	Href *string `json:"href" validate:"required"`

	// The result of the last action.
	Result *string `json:"result,omitempty"`

	// A brief summary of an action.
	Job *ActionJobWithIdAndSummary `json:"job,omitempty"`

	// A brief summary of a pre- and post-action.
	PreJob *PrePostActionJobWithIdAndSummary `json:"pre_job,omitempty"`

	// A brief summary of a pre- and post-action.
	PostJob *PrePostActionJobWithIdAndSummary `json:"post_job,omitempty"`
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
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "job", &obj.Job, UnmarshalActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "job-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pre_job", &obj.PreJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "pre_job-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "post_job", &obj.PostJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "post_job-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastDriftDetectionJobSummary : The summary for drift detection jobs that are performed as part of the last monitoring action.
type LastDriftDetectionJobSummary struct {
	// A brief summary of an action.
	Job *ActionJobWithIdAndSummary `json:"job,omitempty"`
}

// UnmarshalLastDriftDetectionJobSummary unmarshals an instance of LastDriftDetectionJobSummary from the specified map of raw messages.
func UnmarshalLastDriftDetectionJobSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastDriftDetectionJobSummary)
	err = core.UnmarshalModel(m, "job", &obj.Job, UnmarshalActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "job-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastMonitoringActionWithSummary : The summary from the last monitoring action job that is performed on the project configuration.
type LastMonitoringActionWithSummary struct {
	// A URL.
	Href *string `json:"href" validate:"required"`

	// The result of the last action.
	Result *string `json:"result,omitempty"`

	// The summary for drift detection jobs that are performed as part of the last monitoring action.
	DriftDetection *LastDriftDetectionJobSummary `json:"drift_detection,omitempty"`
}

// Constants associated with the LastMonitoringActionWithSummary.Result property.
// The result of the last action.
const (
	LastMonitoringActionWithSummary_Result_Failed = "failed"
	LastMonitoringActionWithSummary_Result_Passed = "passed"
)

// UnmarshalLastMonitoringActionWithSummary unmarshals an instance of LastMonitoringActionWithSummary from the specified map of raw messages.
func UnmarshalLastMonitoringActionWithSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastMonitoringActionWithSummary)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "drift_detection", &obj.DriftDetection, UnmarshalLastDriftDetectionJobSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "drift_detection-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastValidatedActionWithSummary : The href and results from the last action job that is performed on the project configuration.
type LastValidatedActionWithSummary struct {
	// A URL.
	Href *string `json:"href" validate:"required"`

	// The result of the last action.
	Result *string `json:"result,omitempty"`

	// A brief summary of an action.
	Job *ActionJobWithIdAndSummary `json:"job,omitempty"`

	// A brief summary of a pre- and post-action.
	PreJob *PrePostActionJobWithIdAndSummary `json:"pre_job,omitempty"`

	// A brief summary of a pre- and post-action.
	PostJob *PrePostActionJobWithIdAndSummary `json:"post_job,omitempty"`

	// The cost estimate of the configuration.
	// This property exists only after the first configuration validation.
	CostEstimate *ProjectConfigMetadataCostEstimate `json:"cost_estimate,omitempty"`

	// The Code Risk Analyzer logs of the configuration. This property is populated only after the validation step when the
	// Code Risk Analyzer is run. Note: `cra` is the abbreviated form of Code Risk Analyzer.
	CraLogs ProjectConfigMetadataCodeRiskAnalyzerLogsIntf `json:"cra_logs,omitempty"`
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
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "job", &obj.Job, UnmarshalActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "job-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pre_job", &obj.PreJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "pre_job-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "post_job", &obj.PostJob, UnmarshalPrePostActionJobWithIdAndSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "post_job-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cost_estimate", &obj.CostEstimate, UnmarshalProjectConfigMetadataCostEstimate)
	if err != nil {
		err = core.SDKErrorf(err, "", "cost_estimate-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cra_logs", &obj.CraLogs, UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs)
	if err != nil {
		err = core.SDKErrorf(err, "", "cra_logs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListConfigResourcesOptions : The ListConfigResources options.
type ListConfigResourcesOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
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

	// The unique configuration ID.
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

	// The server uses this parameter to determine the first entry that is returned on the next page. If this parameter is
	// not specified, the logical first page is returned.
	Token *string `json:"token,omitempty"`

	// The maximum number of resources to return. The number of resources that are returned is the same, except for the
	// last page.
	Limit *int64 `json:"limit,omitempty"`

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

// SetToken : Allow user to set Token
func (_options *ListConfigsOptions) SetToken(token string) *ListConfigsOptions {
	_options.Token = core.StringPtr(token)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListConfigsOptions) SetLimit(limit int64) *ListConfigsOptions {
	_options.Limit = core.Int64Ptr(limit)
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

	// The server uses this parameter to determine the first entry that is returned on the next page. If this parameter is
	// not specified, the logical first page is returned.
	Token *string `json:"token,omitempty"`

	// The maximum number of resources to return. The number of resources that are returned is the same, except for the
	// last page.
	Limit *int64 `json:"limit,omitempty"`

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

// SetToken : Allow user to set Token
func (_options *ListProjectEnvironmentsOptions) SetToken(token string) *ListProjectEnvironmentsOptions {
	_options.Token = core.StringPtr(token)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProjectEnvironmentsOptions) SetLimit(limit int64) *ListProjectEnvironmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProjectEnvironmentsOptions) SetHeaders(param map[string]string) *ListProjectEnvironmentsOptions {
	options.Headers = param
	return options
}

// ListProjectsOptions : The ListProjects options.
type ListProjectsOptions struct {
	// The server uses this parameter to determine the first entry that is returned on the next page. If this parameter is
	// not specified, the logical first page is returned.
	Token *string `json:"token,omitempty"`

	// The maximum number of resources to return. The number of resources that are returned is the same, except for the
	// last page.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProjectsOptions : Instantiate ListProjectsOptions
func (*ProjectV1) NewListProjectsOptions() *ListProjectsOptions {
	return &ListProjectsOptions{}
}

// SetToken : Allow user to set Token
func (_options *ListProjectsOptions) SetToken(token string) *ListProjectsOptions {
	_options.Token = core.StringPtr(token)
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

// MemberOfDefinition : The stack config parent of which this configuration is a member of.
type MemberOfDefinition struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// The definition summary of the stack configuration.
	Definition *StackConfigDefinitionSummary `json:"definition" validate:"required"`

	// The version of the stack configuration.
	Version *int64 `json:"version" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalMemberOfDefinition unmarshals an instance of MemberOfDefinition from the specified map of raw messages.
func UnmarshalMemberOfDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MemberOfDefinition)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalStackConfigDefinitionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
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

	// This property can be any value - a string, number, boolean, array, or object.
	Value interface{} `json:"value,omitempty"`
}

// UnmarshalOutputValue unmarshals an instance of OutputValue from the specified map of raw messages.
func UnmarshalOutputValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OutputValue)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobSummary : A brief summary of a pre- and post-action job. This property is populated only after an action is run as part of a
// validation, deployment, or undeployment.
type PrePostActionJobSummary struct {
	// The ID of the Schematics action job that ran as part of the pre- and post-job.
	JobID *string `json:"job_id" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	StartTime *strfmt.DateTime `json:"start_time,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	EndTime *strfmt.DateTime `json:"end_time,omitempty"`

	// The number of tasks that were run in the job.
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
		err = core.SDKErrorf(err, "", "job_id-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "tasks", &obj.Tasks)
	if err != nil {
		err = core.SDKErrorf(err, "", "tasks-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ok", &obj.Ok)
	if err != nil {
		err = core.SDKErrorf(err, "", "ok-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped", &obj.Skipped)
	if err != nil {
		err = core.SDKErrorf(err, "", "skipped-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "changed", &obj.Changed)
	if err != nil {
		err = core.SDKErrorf(err, "", "changed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project_error", &obj.ProjectError, UnmarshalPrePostActionJobSystemError)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_error-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobSystemError : The system-level error that OS captured in the project pipelines for the pre- and post-job.
type PrePostActionJobSystemError struct {
	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	Timestamp *strfmt.DateTime `json:"timestamp" validate:"required"`

	// The ID of the user that triggered the pipeline that ran the pre- and post-job.
	UserID *string `json:"user_id" validate:"required"`

	// The HTTP status code for the error.
	StatusCode *string `json:"status_code" validate:"required"`

	// The summary description of the error.
	Description *string `json:"description" validate:"required"`

	// The detailed message from the source error.
	ErrorResponse *string `json:"error_response,omitempty"`
}

// UnmarshalPrePostActionJobSystemError unmarshals an instance of PrePostActionJobSystemError from the specified map of raw messages.
func UnmarshalPrePostActionJobSystemError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrePostActionJobSystemError)
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "timestamp-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "error_response", &obj.ErrorResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_response-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PrePostActionJobWithIdAndSummary : A brief summary of a pre- and post-action.
type PrePostActionJobWithIdAndSummary struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// A brief summary of a pre- and post-action job. This property is populated only after an action is run as part of a
	// validation, deployment, or undeployment.
	Summary *PrePostActionJobSummary `json:"summary" validate:"required"`
}

// UnmarshalPrePostActionJobWithIdAndSummary unmarshals an instance of PrePostActionJobWithIdAndSummary from the specified map of raw messages.
func UnmarshalPrePostActionJobWithIdAndSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrePostActionJobWithIdAndSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalPrePostActionJobSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "summary-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Project : The standard schema of a project.
type Project struct {
	// An IBM Cloud resource name that uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The cumulative list of needs attention items for a project. If the view is successfully retrieved, an empty or
	// nonempty array is returned.
	CumulativeNeedsAttentionView []CumulativeNeedsAttention `json:"cumulative_needs_attention_view" validate:"required"`

	// A value of `true` indicates that the fetch of the needs attention items failed. This property only exists if there
	// was an error when you retrieved the cumulative needs attention view.
	CumulativeNeedsAttentionViewError *bool `json:"cumulative_needs_attention_view_error,omitempty"`

	// The unique project ID.
	ID *string `json:"id" validate:"required"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location" validate:"required"`

	// The resource group ID where the project's data and tools are created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The project status value.
	State *string `json:"state" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The resource group name where the project's data and tools are created.
	ResourceGroup *string `json:"resource_group" validate:"required"`

	// The CRN of the Event Notifications instance if one is connected to this project.
	EventNotificationsCrn *string `json:"event_notifications_crn,omitempty"`

	// The project configurations. These configurations are only included in the response of creating a project if a
	// configuration array is specified in the request payload.
	Configs []ProjectConfigSummary `json:"configs" validate:"required"`

	// The project environment. These environments are only included in the response if project environments were created
	// on the project.
	Environments []ProjectEnvironmentSummary `json:"environments" validate:"required"`

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
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cumulative_needs_attention_view", &obj.CumulativeNeedsAttentionView, UnmarshalCumulativeNeedsAttention)
	if err != nil {
		err = core.SDKErrorf(err, "", "cumulative_needs_attention_view-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cumulative_needs_attention_view_error", &obj.CumulativeNeedsAttentionViewError)
	if err != nil {
		err = core.SDKErrorf(err, "", "cumulative_needs_attention_view_error-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group", &obj.ResourceGroup)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_notifications_crn", &obj.EventNotificationsCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_notifications_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfigSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "configs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalProjectEnvironmentSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "projects", &obj.Projects, UnmarshalProjectSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "projects-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProjectCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	token, err := core.GetQueryParam(resp.Next.Href, "token")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if token == nil {
		return nil, nil
	}
	return token, nil
}

// ProjectComplianceProfile : The profile that is required for compliance.
type ProjectComplianceProfile struct {
	// The unique ID for the compliance profile.
	ID *string `json:"id,omitempty"`

	// A unique ID for the instance of a compliance profile.
	InstanceID *string `json:"instance_id,omitempty"`

	// The location of the compliance instance.
	InstanceLocation *string `json:"instance_location,omitempty"`

	// A unique ID for the attachment to a compliance profile.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The name of the compliance profile.
	ProfileName *string `json:"profile_name,omitempty"`
}

// Constants associated with the ProjectComplianceProfile.InstanceLocation property.
// The location of the compliance instance.
const (
	ProjectComplianceProfile_InstanceLocation_CaTor = "ca-tor"
	ProjectComplianceProfile_InstanceLocation_EuDe = "eu-de"
	ProjectComplianceProfile_InstanceLocation_EuGb = "eu-gb"
	ProjectComplianceProfile_InstanceLocation_UsEast = "us-east"
	ProjectComplianceProfile_InstanceLocation_UsSouth = "us-south"
)

// UnmarshalProjectComplianceProfile unmarshals an instance of ProjectComplianceProfile from the specified map of raw messages.
func UnmarshalProjectComplianceProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectComplianceProfile)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_location", &obj.InstanceLocation)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "attachment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		err = core.SDKErrorf(err, "", "profile_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfig : The standard schema of a project configuration.
type ProjectConfig struct {
	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id" validate:"required"`

	// The version of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// The flag that indicates whether the version of the configuration is draft, or active.
	IsDraft *bool `json:"is_draft" validate:"required"`

	// The needs attention state of a configuration.
	NeedsAttentionState []ProjectConfigNeedsAttentionState `json:"needs_attention_state" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// The last approved metadata of the configuration.
	LastApproved *ProjectConfigMetadataLastApproved `json:"last_approved,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	LastSavedAt *strfmt.DateTime `json:"last_saved_at,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastValidated *LastValidatedActionWithSummary `json:"last_validated,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastDeployed *LastActionWithSummary `json:"last_deployed,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastUndeployed *LastActionWithSummary `json:"last_undeployed,omitempty"`

	// The summary from the last monitoring action job that is performed on the project configuration.
	LastMonitoring *LastMonitoringActionWithSummary `json:"last_monitoring,omitempty"`

	// The outputs of a Schematics template property.
	Outputs []OutputValue `json:"outputs" validate:"required"`

	// The project that is referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// The references that are used in the configuration to resolve input values.
	References map[string]interface{} `json:"references,omitempty"`

	// A Schematics workspace that is associated to a project configuration, with scripts.
	Schematics *SchematicsMetadata `json:"schematics,omitempty"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// The flag that indicates whether a configuration update is available.
	UpdateAvailable *bool `json:"update_available,omitempty"`

	// The stack definition identifier.
	TemplateID *string `json:"template_id,omitempty"`

	// The stack config parent of which this configuration is a member of.
	MemberOf *MemberOfDefinition `json:"member_of,omitempty"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The configuration type.
	DeploymentModel *string `json:"deployment_model,omitempty"`

	// Computed state code clarifying the prerequisites for validation for the configuration.
	StateCode *string `json:"state_code,omitempty"`

	Definition ProjectConfigDefinitionResponseIntf `json:"definition" validate:"required"`

	// A summary of a project configuration version.
	ApprovedVersion *ProjectConfigVersionSummary `json:"approved_version,omitempty"`

	// A summary of a project configuration version.
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

// Constants associated with the ProjectConfig.DeploymentModel property.
// The configuration type.
const (
	ProjectConfig_DeploymentModel_ProjectDeployed = "project_deployed"
	ProjectConfig_DeploymentModel_Stack = "stack"
	ProjectConfig_DeploymentModel_UserDeployed = "user_deployed"
)

// Constants associated with the ProjectConfig.StateCode property.
// Computed state code clarifying the prerequisites for validation for the configuration.
const (
	ProjectConfig_StateCode_AwaitingInput = "awaiting_input"
	ProjectConfig_StateCode_AwaitingMemberDeployment = "awaiting_member_deployment"
	ProjectConfig_StateCode_AwaitingPrerequisite = "awaiting_prerequisite"
	ProjectConfig_StateCode_AwaitingStackSetup = "awaiting_stack_setup"
	ProjectConfig_StateCode_AwaitingValidation = "awaiting_validation"
)

// UnmarshalProjectConfig unmarshals an instance of ProjectConfig from the specified map of raw messages.
func UnmarshalProjectConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfig)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_draft", &obj.IsDraft)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_draft-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "needs_attention_state", &obj.NeedsAttentionState, UnmarshalProjectConfigNeedsAttentionState)
	if err != nil {
		err = core.SDKErrorf(err, "", "needs_attention_state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_approved", &obj.LastApproved, UnmarshalProjectConfigMetadataLastApproved)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_approved-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_saved_at", &obj.LastSavedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_saved_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_validated", &obj.LastValidated, UnmarshalLastValidatedActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_validated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_deployed", &obj.LastDeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_deployed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_undeployed", &obj.LastUndeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_undeployed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_monitoring", &obj.LastMonitoring, UnmarshalLastMonitoringActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalOutputValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "outputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "references", &obj.References)
	if err != nil {
		err = core.SDKErrorf(err, "", "references-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "schematics-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "update_available", &obj.UpdateAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "update_available-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "member_of", &obj.MemberOf, UnmarshalMemberOfDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "member_of-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_model", &obj.DeploymentModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_model-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state_code", &obj.StateCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "state_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigDefinitionResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "approved_version", &obj.ApprovedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "approved_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deployed_version", &obj.DeployedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployed_version-error", common.GetComponentInfo())
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

	// The IBM Cloud API Key. It can be either raw or pulled from the catalog via a `CRN` or `JSON` blob.
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
		err = core.SDKErrorf(err, "", "trusted_profile_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigCollection : The project configuration list.
type ProjectConfigCollection struct {
	// A pagination limit.
	Limit *int64 `json:"limit" validate:"required"`

	// A pagination link.
	First *PaginationLink `json:"first" validate:"required"`

	// A pagination link.
	Next *PaginationLink `json:"next,omitempty"`

	// The response schema of the collection list operation that defines the array property with the name `configs`.
	Configs []ProjectConfigSummary `json:"configs,omitempty"`
}

// UnmarshalProjectConfigCollection unmarshals an instance of ProjectConfigCollection from the specified map of raw messages.
func UnmarshalProjectConfigCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginationLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "configs", &obj.Configs, UnmarshalProjectConfigSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "configs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProjectConfigCollection) GetNextToken() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	token, err := core.GetQueryParam(resp.Next.Href, "token")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if token == nil {
		return nil, nil
	}
	return token, nil
}

// ProjectConfigDefinitionPatch : ProjectConfigDefinitionPatch struct
// Models which "extend" this model:
// - ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch
// - ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch
// - ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch
type ProjectConfigDefinitionPatch struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`
}
func (*ProjectConfigDefinitionPatch) isaProjectConfigDefinitionPatch() bool {
	return true
}

type ProjectConfigDefinitionPatchIntf interface {
	isaProjectConfigDefinitionPatch() bool
}

// UnmarshalProjectConfigDefinitionPatch unmarshals an instance of ProjectConfigDefinitionPatch from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPatch)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPrototype : ProjectConfigDefinitionPrototype struct
// Models which "extend" this model:
// - ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype
// - ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties
// - ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype
type ProjectConfigDefinitionPrototype struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`

	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`
}
func (*ProjectConfigDefinitionPrototype) isaProjectConfigDefinitionPrototype() bool {
	return true
}

type ProjectConfigDefinitionPrototypeIntf interface {
	isaProjectConfigDefinitionPrototype() bool
}

// UnmarshalProjectConfigDefinitionPrototype unmarshals an instance of ProjectConfigDefinitionPrototype from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPrototype)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionResponse : ProjectConfigDefinitionResponse struct
// Models which "extend" this model:
// - ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse
// - ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse
// - ProjectConfigDefinitionResponseStackConfigDefinitionProperties
type ProjectConfigDefinitionResponse struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`
}
func (*ProjectConfigDefinitionResponse) isaProjectConfigDefinitionResponse() bool {
	return true
}

type ProjectConfigDefinitionResponseIntf interface {
	isaProjectConfigDefinitionResponse() bool
}

// UnmarshalProjectConfigDefinitionResponse unmarshals an instance of ProjectConfigDefinitionResponse from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionResponse)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDelete : The ID of the deleted configuration.
type ProjectConfigDelete struct {
	// The ID of the deleted project or configuration.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalProjectConfigDelete unmarshals an instance of ProjectConfigDelete from the specified map of raw messages.
func UnmarshalProjectConfigDelete(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDelete)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataCodeRiskAnalyzerLogs : The Code Risk Analyzer logs of the configuration. This property is populated only after the validation step when the
// Code Risk Analyzer is run. Note: `cra` is the abbreviated form of Code Risk Analyzer.
// Models which "extend" this model:
// - ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204
type ProjectConfigMetadataCodeRiskAnalyzerLogs struct {
	// The version of the Code Risk Analyzer logs of the configuration. The metadata for this schema is specific to Code
	// Risk Analyzer version 2.0.4.
	CraVersion *string `json:"cra_version,omitempty"`

	// The schema version of Code Risk Analyzer logs of the configuration.
	SchemaVersion *string `json:"schema_version,omitempty"`

	// The status of the Code Risk Analyzer logs of the configuration.
	Status *string `json:"status,omitempty"`

	// The Code Risk Analyzer logs a summary of the configuration.
	Summary *CodeRiskAnalyzerLogsSummary `json:"summary,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	Timestamp *strfmt.DateTime `json:"timestamp,omitempty"`
}

// Constants associated with the ProjectConfigMetadataCodeRiskAnalyzerLogs.Status property.
// The status of the Code Risk Analyzer logs of the configuration.
const (
	ProjectConfigMetadataCodeRiskAnalyzerLogs_Status_Failed = "failed"
	ProjectConfigMetadataCodeRiskAnalyzerLogs_Status_Passed = "passed"
)
func (*ProjectConfigMetadataCodeRiskAnalyzerLogs) isaProjectConfigMetadataCodeRiskAnalyzerLogs() bool {
	return true
}

type ProjectConfigMetadataCodeRiskAnalyzerLogsIntf interface {
	isaProjectConfigMetadataCodeRiskAnalyzerLogs() bool
}

// UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs unmarshals an instance of ProjectConfigMetadataCodeRiskAnalyzerLogs from the specified map of raw messages.
func UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigMetadataCodeRiskAnalyzerLogs)
	err = core.UnmarshalPrimitive(m, "cra_version", &obj.CraVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "cra_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "schema_version", &obj.SchemaVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "schema_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalCodeRiskAnalyzerLogsSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "timestamp-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataCostEstimate : The cost estimate of the configuration. This property exists only after the first configuration validation.
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

	// The difference between the current and past total hourly cost estimates of the configuration.
	DiffTotalHourlyCost *string `json:"diffTotalHourlyCost,omitempty"`

	// The difference between the current and past total monthly cost estimates of the configuration.
	DiffTotalMonthlyCost *string `json:"diffTotalMonthlyCost,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
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
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "currency", &obj.Currency)
	if err != nil {
		err = core.SDKErrorf(err, "", "currency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "totalHourlyCost", &obj.TotalHourlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "totalHourlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "totalMonthlyCost", &obj.TotalMonthlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "totalMonthlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalHourlyCost", &obj.PastTotalHourlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "pastTotalHourlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pastTotalMonthlyCost", &obj.PastTotalMonthlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "pastTotalMonthlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalHourlyCost", &obj.DiffTotalHourlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "diffTotalHourlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "diffTotalMonthlyCost", &obj.DiffTotalMonthlyCost)
	if err != nil {
		err = core.SDKErrorf(err, "", "diffTotalMonthlyCost-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeGenerated", &obj.TimeGenerated)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeGenerated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataLastApproved : The last approved metadata of the configuration.
type ProjectConfigMetadataLastApproved struct {
	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	At *strfmt.DateTime `json:"at" validate:"required"`

	// The comment that is left by the user who approved the configuration.
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
		err = core.SDKErrorf(err, "", "at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "comment", &obj.Comment)
	if err != nil {
		err = core.SDKErrorf(err, "", "comment-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_forced", &obj.IsForced)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_forced-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigNeedsAttentionState : A needs attention state item shown to users is a specific actionable event that occurs during the lifecycle of a
// configuration.
type ProjectConfigNeedsAttentionState struct {
	// The id of the event.
	EventID *string `json:"event_id" validate:"required"`

	// The name of the event.
	Event *string `json:"event" validate:"required"`

	// The severity of the event. This is a system generated field. For user triggered events the field is not present.
	Severity *string `json:"severity,omitempty"`

	// An actionable URL that users can access in response to the event. This is a system generated field. For user
	// triggered events the field is not present.
	ActionURL *string `json:"action_url,omitempty"`

	// The configuration id and version for which the event occurred. This field is only available for user generated
	// events. For system triggered events the field is not present.
	Target *string `json:"target,omitempty"`

	// The IAM id of the user that triggered the event. This field is only available for user generated events. For system
	// triggered events the field is not present.
	TriggeredBy *string `json:"triggered_by,omitempty"`

	// The timestamp of the event.
	Timestamp *string `json:"timestamp" validate:"required"`
}

// Constants associated with the ProjectConfigNeedsAttentionState.Severity property.
// The severity of the event. This is a system generated field. For user triggered events the field is not present.
const (
	ProjectConfigNeedsAttentionState_Severity_Error = "ERROR"
	ProjectConfigNeedsAttentionState_Severity_Info = "INFO"
	ProjectConfigNeedsAttentionState_Severity_Warning = "WARNING"
)

// UnmarshalProjectConfigNeedsAttentionState unmarshals an instance of ProjectConfigNeedsAttentionState from the specified map of raw messages.
func UnmarshalProjectConfigNeedsAttentionState(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigNeedsAttentionState)
	err = core.UnmarshalPrimitive(m, "event_id", &obj.EventID)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event", &obj.Event)
	if err != nil {
		err = core.SDKErrorf(err, "", "event-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "severity", &obj.Severity)
	if err != nil {
		err = core.SDKErrorf(err, "", "severity-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "action_url", &obj.ActionURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "action_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "triggered_by", &obj.TriggeredBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "triggered_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "timestamp-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigPrototype : The input of a project configuration.
type ProjectConfigPrototype struct {
	Definition ProjectConfigDefinitionPrototypeIntf `json:"definition" validate:"required"`

	// A Schematics workspace to use for deploying this deployable architecture.
	// > If you are importing data from an existing Schematics workspace that is not backed by cart, then you must provide
	// a `locator_id`. If you are using a Schematics workspace that is backed by cart, a `locator_id` is not required
	// because the Schematics workspace has one.
	// >
	// There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400`code  is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// >
	// For more information, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	Schematics *SchematicsWorkspace `json:"schematics,omitempty"`
}

// NewProjectConfigPrototype : Instantiate ProjectConfigPrototype (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigPrototype(definition ProjectConfigDefinitionPrototypeIntf) (_model *ProjectConfigPrototype, err error) {
	_model = &ProjectConfigPrototype{
		Definition: definition,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProjectConfigPrototype unmarshals an instance of ProjectConfigPrototype from the specified map of raw messages.
func UnmarshalProjectConfigPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigPrototype)
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigDefinitionPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsWorkspace)
	if err != nil {
		err = core.SDKErrorf(err, "", "schematics-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResource : ProjectConfigResource struct
type ProjectConfigResource struct {
	// An IBM Cloud resource name that uniquely identifies a resource.
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
		err = core.SDKErrorf(err, "", "resource_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_tainted", &obj.ResourceTainted)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_tainted-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_name", &obj.ResourceGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigResourceCollection : The project configuration resource list.
type ProjectConfigResourceCollection struct {
	// The collection list operation response schema that defines the array property with the name `resources`.
	Resources []ProjectConfigResource `json:"resources" validate:"required"`

	// The total number of resources that are deployed by the configuration.
	ResourcesCount *int64 `json:"resources_count" validate:"required"`
}

// UnmarshalProjectConfigResourceCollection unmarshals an instance of ProjectConfigResourceCollection from the specified map of raw messages.
func UnmarshalProjectConfigResourceCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigResourceCollection)
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalProjectConfigResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resources_count", &obj.ResourcesCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources_count-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigSummary : ProjectConfigSummary struct
type ProjectConfigSummary struct {
	// A summary of a project configuration version.
	ApprovedVersion *ProjectConfigVersionSummary `json:"approved_version,omitempty"`

	// A summary of a project configuration version.
	DeployedVersion *ProjectConfigVersionSummary `json:"deployed_version,omitempty"`

	// The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.
	ID *string `json:"id" validate:"required"`

	// The version of the configuration.
	Version *int64 `json:"version" validate:"required"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The description of a project configuration.
	Definition *ProjectConfigSummaryDefinition `json:"definition" validate:"required"`

	// The project that is referenced by this resource.
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
	ProjectConfigSummary_DeploymentModel_Stack = "stack"
	ProjectConfigSummary_DeploymentModel_UserDeployed = "user_deployed"
)

// UnmarshalProjectConfigSummary unmarshals an instance of ProjectConfigSummary from the specified map of raw messages.
func UnmarshalProjectConfigSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigSummary)
	err = core.UnmarshalModel(m, "approved_version", &obj.ApprovedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "approved_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deployed_version", &obj.DeployedVersion, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployed_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigSummaryDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_model", &obj.DeploymentModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_model-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigSummaryDefinition : The description of a project configuration.
type ProjectConfigSummaryDefinition struct {
	// A project configuration description.
	Description *string `json:"description" validate:"required"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`
}

// UnmarshalProjectConfigSummaryDefinition unmarshals an instance of ProjectConfigSummaryDefinition from the specified map of raw messages.
func UnmarshalProjectConfigSummaryDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigSummaryDefinition)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
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
	NeedsAttentionState []ProjectConfigNeedsAttentionState `json:"needs_attention_state" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// The last approved metadata of the configuration.
	LastApproved *ProjectConfigMetadataLastApproved `json:"last_approved,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	LastSavedAt *strfmt.DateTime `json:"last_saved_at,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastValidated *LastValidatedActionWithSummary `json:"last_validated,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastDeployed *LastActionWithSummary `json:"last_deployed,omitempty"`

	// The href and results from the last action job that is performed on the project configuration.
	LastUndeployed *LastActionWithSummary `json:"last_undeployed,omitempty"`

	// The summary from the last monitoring action job that is performed on the project configuration.
	LastMonitoring *LastMonitoringActionWithSummary `json:"last_monitoring,omitempty"`

	// The outputs of a Schematics template property.
	Outputs []OutputValue `json:"outputs" validate:"required"`

	// The project that is referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// The references that are used in the configuration to resolve input values.
	References map[string]interface{} `json:"references,omitempty"`

	// A Schematics workspace that is associated to a project configuration, with scripts.
	Schematics *SchematicsMetadata `json:"schematics,omitempty"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// The flag that indicates whether a configuration update is available.
	UpdateAvailable *bool `json:"update_available,omitempty"`

	// The stack definition identifier.
	TemplateID *string `json:"template_id,omitempty"`

	// The stack config parent of which this configuration is a member of.
	MemberOf *MemberOfDefinition `json:"member_of,omitempty"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The configuration type.
	DeploymentModel *string `json:"deployment_model,omitempty"`

	// Computed state code clarifying the prerequisites for validation for the configuration.
	StateCode *string `json:"state_code,omitempty"`

	Definition ProjectConfigDefinitionResponseIntf `json:"definition" validate:"required"`
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

// Constants associated with the ProjectConfigVersion.DeploymentModel property.
// The configuration type.
const (
	ProjectConfigVersion_DeploymentModel_ProjectDeployed = "project_deployed"
	ProjectConfigVersion_DeploymentModel_Stack = "stack"
	ProjectConfigVersion_DeploymentModel_UserDeployed = "user_deployed"
)

// Constants associated with the ProjectConfigVersion.StateCode property.
// Computed state code clarifying the prerequisites for validation for the configuration.
const (
	ProjectConfigVersion_StateCode_AwaitingInput = "awaiting_input"
	ProjectConfigVersion_StateCode_AwaitingMemberDeployment = "awaiting_member_deployment"
	ProjectConfigVersion_StateCode_AwaitingPrerequisite = "awaiting_prerequisite"
	ProjectConfigVersion_StateCode_AwaitingStackSetup = "awaiting_stack_setup"
	ProjectConfigVersion_StateCode_AwaitingValidation = "awaiting_validation"
)

// UnmarshalProjectConfigVersion unmarshals an instance of ProjectConfigVersion from the specified map of raw messages.
func UnmarshalProjectConfigVersion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersion)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_draft", &obj.IsDraft)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_draft-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "needs_attention_state", &obj.NeedsAttentionState, UnmarshalProjectConfigNeedsAttentionState)
	if err != nil {
		err = core.SDKErrorf(err, "", "needs_attention_state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_approved", &obj.LastApproved, UnmarshalProjectConfigMetadataLastApproved)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_approved-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_saved_at", &obj.LastSavedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_saved_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_validated", &obj.LastValidated, UnmarshalLastValidatedActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_validated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_deployed", &obj.LastDeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_deployed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_undeployed", &obj.LastUndeployed, UnmarshalLastActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_undeployed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_monitoring", &obj.LastMonitoring, UnmarshalLastMonitoringActionWithSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalOutputValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "outputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "references", &obj.References)
	if err != nil {
		err = core.SDKErrorf(err, "", "references-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "schematics", &obj.Schematics, UnmarshalSchematicsMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "schematics-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "update_available", &obj.UpdateAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "update_available-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "member_of", &obj.MemberOf, UnmarshalMemberOfDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "member_of-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deployment_model", &obj.DeploymentModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment_model-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state_code", &obj.StateCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "state_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigDefinitionResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersionDefinitionSummary : A summary of the definition in a project configuration version.
type ProjectConfigVersionDefinitionSummary struct {
	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`
}

// UnmarshalProjectConfigVersionDefinitionSummary unmarshals an instance of ProjectConfigVersionDefinitionSummary from the specified map of raw messages.
func UnmarshalProjectConfigVersionDefinitionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersionDefinitionSummary)
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersionSummary : A summary of a project configuration version.
type ProjectConfigVersionSummary struct {
	// A summary of the definition in a project configuration version.
	Definition *ProjectConfigVersionDefinitionSummary `json:"definition" validate:"required"`

	// The state of the configuration.
	State *string `json:"state" validate:"required"`

	// Computed state code clarifying the prerequisites for validation for the configuration.
	StateCode *string `json:"state_code,omitempty"`

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

// Constants associated with the ProjectConfigVersionSummary.StateCode property.
// Computed state code clarifying the prerequisites for validation for the configuration.
const (
	ProjectConfigVersionSummary_StateCode_AwaitingInput = "awaiting_input"
	ProjectConfigVersionSummary_StateCode_AwaitingMemberDeployment = "awaiting_member_deployment"
	ProjectConfigVersionSummary_StateCode_AwaitingPrerequisite = "awaiting_prerequisite"
	ProjectConfigVersionSummary_StateCode_AwaitingStackSetup = "awaiting_stack_setup"
	ProjectConfigVersionSummary_StateCode_AwaitingValidation = "awaiting_validation"
)

// UnmarshalProjectConfigVersionSummary unmarshals an instance of ProjectConfigVersionSummary from the specified map of raw messages.
func UnmarshalProjectConfigVersionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersionSummary)
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectConfigVersionDefinitionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state_code", &obj.StateCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "state_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigVersionSummaryCollection : The project configuration version list.
type ProjectConfigVersionSummaryCollection struct {
	// The collection list operation response schema that defines the array property with the name `versions`.
	Versions []ProjectConfigVersionSummary `json:"versions" validate:"required"`
}

// UnmarshalProjectConfigVersionSummaryCollection unmarshals an instance of ProjectConfigVersionSummaryCollection from the specified map of raw messages.
func UnmarshalProjectConfigVersionSummaryCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigVersionSummaryCollection)
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalProjectConfigVersionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "versions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectDefinitionProperties : The definition of the project.
type ProjectDefinitionProperties struct {
	// The name of the project.  It's unique within the account across regions.
	Name *string `json:"name" validate:"required"`

	// The policy that indicates whether the resources are destroyed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete" validate:"required"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project
	// without providing a description.
	Description *string `json:"description" validate:"required"`

	// A boolean flag to enable auto deploy.
	AutoDeploy *bool `json:"auto_deploy,omitempty"`

	// A boolean flag to enable automatic drift detection. Use this field to run a daily check to compare your
	// configurations to your deployed resources to detect any difference.
	MonitoringEnabled *bool `json:"monitoring_enabled,omitempty"`
}

// UnmarshalProjectDefinitionProperties unmarshals an instance of ProjectDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectDefinitionProperties)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_on_delete-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_deploy", &obj.AutoDeploy)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_deploy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "monitoring_enabled", &obj.MonitoringEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "monitoring_enabled-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectDefinitionSummary : The definition of the project.
type ProjectDefinitionSummary struct {
	// The name of the project.  It's unique within the account across regions.
	Name *string `json:"name" validate:"required"`

	// The policy that indicates whether the resources are destroyed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete" validate:"required"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project
	// without providing a description.
	Description *string `json:"description" validate:"required"`
}

// UnmarshalProjectDefinitionSummary unmarshals an instance of ProjectDefinitionSummary from the specified map of raw messages.
func UnmarshalProjectDefinitionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectDefinitionSummary)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_on_delete-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectDeleteResponse : The ID of the deleted project.
type ProjectDeleteResponse struct {
	// The ID of the deleted project or configuration.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalProjectDeleteResponse unmarshals an instance of ProjectDeleteResponse from the specified map of raw messages.
func UnmarshalProjectDeleteResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectDeleteResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectEnvironmentSummary : The environment metadata.
type ProjectEnvironmentSummary struct {
	// The environment ID as a friendly name.
	ID *string `json:"id" validate:"required"`

	// The project that is referenced by this resource.
	Project *ProjectReference `json:"project" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The environment definition that is used in the project collection.
	Definition *ProjectEnvironmentSummaryDefinition `json:"definition" validate:"required"`
}

// UnmarshalProjectEnvironmentSummary unmarshals an instance of ProjectEnvironmentSummary from the specified map of raw messages.
func UnmarshalProjectEnvironmentSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectEnvironmentSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "project", &obj.Project, UnmarshalProjectReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectEnvironmentSummaryDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectEnvironmentSummaryDefinition : The environment definition that is used in the project collection.
type ProjectEnvironmentSummaryDefinition struct {
	// The description of the environment.
	Description *string `json:"description" validate:"required"`

	// The name of the environment. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalProjectEnvironmentSummaryDefinition unmarshals an instance of ProjectEnvironmentSummaryDefinition from the specified map of raw messages.
func UnmarshalProjectEnvironmentSummaryDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectEnvironmentSummaryDefinition)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectPatchDefinitionBlock : The definition of the project.
type ProjectPatchDefinitionBlock struct {
	// The name of the project.  It's unique within the account across regions.
	Name *string `json:"name,omitempty"`

	// The policy that indicates whether the resources are destroyed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete,omitempty"`

	// A boolean flag to enable auto deploy.
	AutoDeploy *bool `json:"auto_deploy,omitempty"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project
	// without providing a description.
	Description *string `json:"description,omitempty"`

	// A boolean flag to enable automatic drift detection. Use this field to run a daily check to compare your
	// configurations to your deployed resources to detect any difference.
	MonitoringEnabled *bool `json:"monitoring_enabled,omitempty"`
}

// UnmarshalProjectPatchDefinitionBlock unmarshals an instance of ProjectPatchDefinitionBlock from the specified map of raw messages.
func UnmarshalProjectPatchDefinitionBlock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectPatchDefinitionBlock)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_on_delete-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_deploy", &obj.AutoDeploy)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_deploy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "monitoring_enabled", &obj.MonitoringEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "monitoring_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectPrototypeDefinition : The definition of the project.
type ProjectPrototypeDefinition struct {
	// The name of the project.  It's unique within the account across regions.
	Name *string `json:"name" validate:"required"`

	// The policy that indicates whether the resources are undeployed or not when a project is deleted.
	DestroyOnDelete *bool `json:"destroy_on_delete,omitempty"`

	// A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project
	// without providing a description.
	Description *string `json:"description,omitempty"`

	// A boolean flag to enable auto deploy.
	AutoDeploy *bool `json:"auto_deploy,omitempty"`

	// A boolean flag to enable automatic drift detection. Use this field to run a daily check to compare your
	// configurations to your deployed resources to detect any difference.
	MonitoringEnabled *bool `json:"monitoring_enabled,omitempty"`
}

// NewProjectPrototypeDefinition : Instantiate ProjectPrototypeDefinition (Generic Model Constructor)
func (*ProjectV1) NewProjectPrototypeDefinition(name string) (_model *ProjectPrototypeDefinition, err error) {
	_model = &ProjectPrototypeDefinition{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProjectPrototypeDefinition unmarshals an instance of ProjectPrototypeDefinition from the specified map of raw messages.
func UnmarshalProjectPrototypeDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectPrototypeDefinition)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destroy_on_delete", &obj.DestroyOnDelete)
	if err != nil {
		err = core.SDKErrorf(err, "", "destroy_on_delete-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_deploy", &obj.AutoDeploy)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_deploy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "monitoring_enabled", &obj.MonitoringEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "monitoring_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectReference : The project that is referenced by this resource.
type ProjectReference struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The definition of the project reference.
	Definition *ProjectDefinitionReference `json:"definition" validate:"required"`

	// An IBM Cloud resource name that uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`
}

// UnmarshalProjectReference unmarshals an instance of ProjectReference from the specified map of raw messages.
func UnmarshalProjectReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectSummary : ProjectSummary struct
type ProjectSummary struct {
	// An IBM Cloud resource name that uniquely identifies a resource.
	Crn *string `json:"crn" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The cumulative list of needs attention items for a project. If the view is successfully retrieved, an empty or
	// nonempty array is returned.
	CumulativeNeedsAttentionView []CumulativeNeedsAttention `json:"cumulative_needs_attention_view" validate:"required"`

	// A value of `true` indicates that the fetch of the needs attention items failed. This property only exists if there
	// was an error when you retrieved the cumulative needs attention view.
	CumulativeNeedsAttentionViewError *bool `json:"cumulative_needs_attention_view_error,omitempty"`

	// The unique project ID.
	ID *string `json:"id" validate:"required"`

	// The IBM Cloud location where a resource is deployed.
	Location *string `json:"location" validate:"required"`

	// The resource group ID where the project's data and tools are created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The project status value.
	State *string `json:"state" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The definition of the project.
	Definition *ProjectDefinitionSummary `json:"definition" validate:"required"`
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
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cumulative_needs_attention_view", &obj.CumulativeNeedsAttentionView, UnmarshalCumulativeNeedsAttention)
	if err != nil {
		err = core.SDKErrorf(err, "", "cumulative_needs_attention_view-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cumulative_needs_attention_view_error", &obj.CumulativeNeedsAttentionViewError)
	if err != nil {
		err = core.SDKErrorf(err, "", "cumulative_needs_attention_view_error-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalProjectDefinitionSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsMetadata : A Schematics workspace that is associated to a project configuration, with scripts.
type SchematicsMetadata struct {
	// An IBM Cloud resource name that uniquely identifies a resource.
	WorkspaceCrn *string `json:"workspace_crn,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	ValidatePreScript *Script `json:"validate_pre_script,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	ValidatePostScript *Script `json:"validate_post_script,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	DeployPreScript *Script `json:"deploy_pre_script,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	DeployPostScript *Script `json:"deploy_post_script,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	UndeployPreScript *Script `json:"undeploy_pre_script,omitempty"`

	// A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
	// deploy, or undeploy).
	UndeployPostScript *Script `json:"undeploy_post_script,omitempty"`
}

// UnmarshalSchematicsMetadata unmarshals an instance of SchematicsMetadata from the specified map of raw messages.
func UnmarshalSchematicsMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsMetadata)
	err = core.UnmarshalPrimitive(m, "workspace_crn", &obj.WorkspaceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "workspace_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "validate_pre_script", &obj.ValidatePreScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "validate_pre_script-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "validate_post_script", &obj.ValidatePostScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "validate_post_script-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deploy_pre_script", &obj.DeployPreScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "deploy_pre_script-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deploy_post_script", &obj.DeployPostScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "deploy_post_script-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "undeploy_pre_script", &obj.UndeployPreScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "undeploy_pre_script-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "undeploy_post_script", &obj.UndeployPostScript, UnmarshalScript)
	if err != nil {
		err = core.SDKErrorf(err, "", "undeploy_post_script-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SchematicsWorkspace : A Schematics workspace to use for deploying this deployable architecture.
// > If you are importing data from an existing Schematics workspace that is not backed by cart, then you must provide a
// `locator_id`. If you are using a Schematics workspace that is backed by cart, a `locator_id` is not required because
// the Schematics workspace has one.
// > There are 3 scenarios:
// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
// existing schematics workspace.
// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400`code  is returned if the specified
// `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
// > For more information, see [Creating workspaces and importing your Terraform
// template](/docs/schematics?topic=schematics-sch-create-wks).
type SchematicsWorkspace struct {
	// An IBM Cloud resource name that uniquely identifies a resource.
	WorkspaceCrn *string `json:"workspace_crn,omitempty"`
}

// UnmarshalSchematicsWorkspace unmarshals an instance of SchematicsWorkspace from the specified map of raw messages.
func UnmarshalSchematicsWorkspace(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SchematicsWorkspace)
	err = core.UnmarshalPrimitive(m, "workspace_crn", &obj.WorkspaceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "workspace_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Script : A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate,
// deploy, or undeploy).
type Script struct {
	// The type of the script.
	Type *string `json:"type,omitempty"`

	// The path to this script is within the current version source.
	Path *string `json:"path,omitempty"`

	// The short description for this script.
	ShortDescription *string `json:"short_description,omitempty"`
}

// UnmarshalScript unmarshals an instance of Script from the specified map of raw messages.
func UnmarshalScript(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Script)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "short_description", &obj.ShortDescription)
	if err != nil {
		err = core.SDKErrorf(err, "", "short_description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackConfigDefinitionSummary : The definition summary of the stack configuration.
type StackConfigDefinitionSummary struct {
	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members" validate:"required"`
}

// UnmarshalStackConfigDefinitionSummary unmarshals an instance of StackConfigDefinitionSummary from the specified map of raw messages.
func UnmarshalStackConfigDefinitionSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackConfigDefinitionSummary)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackConfigMember : A member deployable architecture that is included in your stack.
type StackConfigMember struct {
	// The name matching the alias in the stack definition.
	Name *string `json:"name" validate:"required"`

	// The unique ID.
	ConfigID *string `json:"config_id" validate:"required"`
}

// NewStackConfigMember : Instantiate StackConfigMember (Generic Model Constructor)
func (*ProjectV1) NewStackConfigMember(name string, configID string) (_model *StackConfigMember, err error) {
	_model = &StackConfigMember{
		Name: core.StringPtr(name),
		ConfigID: core.StringPtr(configID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalStackConfigMember unmarshals an instance of StackConfigMember from the specified map of raw messages.
func UnmarshalStackConfigMember(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackConfigMember)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "config_id", &obj.ConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "config_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinition : The stack definition.
type StackDefinition struct {
	// The ID of the stack definition.
	ID *string `json:"id" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	ModifiedAt *strfmt.DateTime `json:"modified_at" validate:"required"`

	// The state for the stack definition.
	State *string `json:"state" validate:"required"`

	// The configuration reference.
	Configuration *StackDefinitionMetadataConfiguration `json:"configuration" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The definition block for a stack definition.
	StackDefinition *StackDefinitionBlock `json:"stack_definition" validate:"required"`
}

// Constants associated with the StackDefinition.State property.
// The state for the stack definition.
const (
	StackDefinition_State_Draft = "draft"
	StackDefinition_State_Published = "published"
)

// UnmarshalStackDefinition unmarshals an instance of StackDefinition from the specified map of raw messages.
func UnmarshalStackDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinition)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_at", &obj.ModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "configuration", &obj.Configuration, UnmarshalStackDefinitionMetadataConfiguration)
	if err != nil {
		err = core.SDKErrorf(err, "", "configuration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "stack_definition", &obj.StackDefinition, UnmarshalStackDefinitionBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "stack_definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionBlock : The definition block for a stack definition.
type StackDefinitionBlock struct {
	// Defines the inputs that users need to configure at the stack level. These inputs are included in the catalog entry
	// when the deployable architecture stack is exported to a private catalog.
	Inputs []StackDefinitionInputVariable `json:"inputs,omitempty"`

	// The outputs associated with this stack definition.
	Outputs []StackDefinitionOutputVariable `json:"outputs,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackDefinitionMember `json:"members,omitempty"`
}

// UnmarshalStackDefinitionBlock unmarshals an instance of StackDefinitionBlock from the specified map of raw messages.
func UnmarshalStackDefinitionBlock(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionBlock)
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalStackDefinitionInputVariable)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalStackDefinitionOutputVariable)
	if err != nil {
		err = core.SDKErrorf(err, "", "outputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackDefinitionMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionBlockPrototype : The definition block for a stack definition.
type StackDefinitionBlockPrototype struct {
	// Defines the inputs that users need to configure at the stack level. These inputs are included in the catalog entry
	// when the deployable architecture stack is exported to a private catalog.
	Inputs []StackDefinitionInputVariable `json:"inputs,omitempty"`

	// The outputs associated with this stack definition.
	Outputs []StackDefinitionOutputVariable `json:"outputs,omitempty"`
}

// UnmarshalStackDefinitionBlockPrototype unmarshals an instance of StackDefinitionBlockPrototype from the specified map of raw messages.
func UnmarshalStackDefinitionBlockPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionBlockPrototype)
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalStackDefinitionInputVariable)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "outputs", &obj.Outputs, UnmarshalStackDefinitionOutputVariable)
	if err != nil {
		err = core.SDKErrorf(err, "", "outputs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionExportRequest : The payload for the stack definition export request.
// Models which "extend" this model:
// - StackDefinitionExportRequestStackDefinitionExportCatalogRequest
// - StackDefinitionExportRequestStackDefinitionExportProductRequest
type StackDefinitionExportRequest struct {
	// The catalog ID to publish.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The semver value of this new version of the product.
	TargetVersion *string `json:"target_version,omitempty"`

	// The variation of this new version of the product.
	Variation *string `json:"variation,omitempty"`

	// The product label.
	Label *string `json:"label,omitempty"`

	// Tags associated with the catalog product.
	Tags []string `json:"tags,omitempty"`

	// The product ID to publish.
	ProductID *string `json:"product_id,omitempty"`
}
func (*StackDefinitionExportRequest) isaStackDefinitionExportRequest() bool {
	return true
}

type StackDefinitionExportRequestIntf interface {
	isaStackDefinitionExportRequest() bool
}

// UnmarshalStackDefinitionExportRequest unmarshals an instance of StackDefinitionExportRequest from the specified map of raw messages.
func UnmarshalStackDefinitionExportRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionExportRequest)
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_version", &obj.TargetVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "variation", &obj.Variation)
	if err != nil {
		err = core.SDKErrorf(err, "", "variation-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "product_id", &obj.ProductID)
	if err != nil {
		err = core.SDKErrorf(err, "", "product_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionExportResponse : The payload for the stack definition export response.
type StackDefinitionExportResponse struct {
	// The catalog ID to publish.
	CatalogID *string `json:"catalog_id,omitempty"`

	// The product ID to publish.
	ProductID *string `json:"product_id,omitempty"`

	// The version locator of the created deployable architecture.
	VersionLocator *string `json:"version_locator,omitempty"`

	// The product target kind value.
	Kind *string `json:"kind,omitempty"`

	// The product format kind value.
	Format *string `json:"format,omitempty"`
}

// UnmarshalStackDefinitionExportResponse unmarshals an instance of StackDefinitionExportResponse from the specified map of raw messages.
func UnmarshalStackDefinitionExportResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionExportResponse)
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "product_id", &obj.ProductID)
	if err != nil {
		err = core.SDKErrorf(err, "", "product_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version_locator", &obj.VersionLocator)
	if err != nil {
		err = core.SDKErrorf(err, "", "version_locator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionInputVariable : The input variables for a stack definition.
type StackDefinitionInputVariable struct {
	// The stack definition input name.
	Name *string `json:"name" validate:"required"`

	// The variable type.
	Type *string `json:"type" validate:"required"`

	// The description of the variable.
	Description *string `json:"description,omitempty"`

	// This property can be any value - a string, number, boolean, array, or object.
	Default interface{} `json:"default" validate:"required"`

	// A boolean value to denote if the property is required.
	Required *bool `json:"required,omitempty"`

	// A boolean value to denote whether the property is hidden, as in not exposed to the user.
	Hidden *bool `json:"hidden,omitempty"`
}

// Constants associated with the StackDefinitionInputVariable.Type property.
// The variable type.
const (
	StackDefinitionInputVariable_Type_Array = "array"
	StackDefinitionInputVariable_Type_Boolean = "boolean"
	StackDefinitionInputVariable_Type_Float = "float"
	StackDefinitionInputVariable_Type_Int = "int"
	StackDefinitionInputVariable_Type_Number = "number"
	StackDefinitionInputVariable_Type_Object = "object"
	StackDefinitionInputVariable_Type_Password = "password"
	StackDefinitionInputVariable_Type_String = "string"
)

// NewStackDefinitionInputVariable : Instantiate StackDefinitionInputVariable (Generic Model Constructor)
func (*ProjectV1) NewStackDefinitionInputVariable(name string, typeVar string, defaultVar interface{}) (_model *StackDefinitionInputVariable, err error) {
	_model = &StackDefinitionInputVariable{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
		Default: defaultVar,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalStackDefinitionInputVariable unmarshals an instance of StackDefinitionInputVariable from the specified map of raw messages.
func UnmarshalStackDefinitionInputVariable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionInputVariable)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		err = core.SDKErrorf(err, "", "default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "required", &obj.Required)
	if err != nil {
		err = core.SDKErrorf(err, "", "required-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionMember : The member definition associated with this stack definition.
type StackDefinitionMember struct {
	// The name matching the alias in the stack definition.
	Name *string `json:"name" validate:"required"`

	// The version locator of the member deployable architecture.
	VersionLocator *string `json:"version_locator" validate:"required"`

	// The member inputs to use for the stack definition.
	Inputs []StackDefinitionMemberInput `json:"inputs,omitempty"`
}

// UnmarshalStackDefinitionMember unmarshals an instance of StackDefinitionMember from the specified map of raw messages.
func UnmarshalStackDefinitionMember(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionMember)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version_locator", &obj.VersionLocator)
	if err != nil {
		err = core.SDKErrorf(err, "", "version_locator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "inputs", &obj.Inputs, UnmarshalStackDefinitionMemberInput)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionMemberInput : StackDefinitionMemberInput struct
type StackDefinitionMemberInput struct {
	// The member input name to use.
	Name *string `json:"name" validate:"required"`

	// The value of the stack definition output.
	Value interface{} `json:"value" validate:"required"`
}

// UnmarshalStackDefinitionMemberInput unmarshals an instance of StackDefinitionMemberInput from the specified map of raw messages.
func UnmarshalStackDefinitionMemberInput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionMemberInput)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionMetadataConfiguration : The configuration reference.
type StackDefinitionMetadataConfiguration struct {
	// The unique ID.
	ID *string `json:"id" validate:"required"`

	// A URL.
	Href *string `json:"href" validate:"required"`

	// The definition of the config reference.
	Definition *ConfigDefinitionReference `json:"definition" validate:"required"`
}

// UnmarshalStackDefinitionMetadataConfiguration unmarshals an instance of StackDefinitionMetadataConfiguration from the specified map of raw messages.
func UnmarshalStackDefinitionMetadataConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionMetadataConfiguration)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalConfigDefinitionReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionOutputVariable : The output variables for a stack definition.
type StackDefinitionOutputVariable struct {
	// The stack definition output name.
	Name *string `json:"name" validate:"required"`

	// The value of the stack definition output.
	Value interface{} `json:"value" validate:"required"`
}

// NewStackDefinitionOutputVariable : Instantiate StackDefinitionOutputVariable (Generic Model Constructor)
func (*ProjectV1) NewStackDefinitionOutputVariable(name string, value interface{}) (_model *StackDefinitionOutputVariable, err error) {
	_model = &StackDefinitionOutputVariable{
		Name: core.StringPtr(name),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalStackDefinitionOutputVariable unmarshals an instance of StackDefinitionOutputVariable from the specified map of raw messages.
func UnmarshalStackDefinitionOutputVariable(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionOutputVariable)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SyncConfigOptions : The SyncConfig options.
type SyncConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// A Schematics workspace to use for deploying this deployable architecture.
	// > If you are importing data from an existing Schematics workspace that is not backed by cart, then you must provide
	// a `locator_id`. If you are using a Schematics workspace that is backed by cart, a `locator_id` is not required
	// because the Schematics workspace has one.
	// >
	// There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400`code  is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// >
	// For more information, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
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

// TerraformLogAnalyzerErrorMessage : The error message that is parsed by the Terraform log analyzer.
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalTerraformLogAnalyzerErrorMessage unmarshals an instance of TerraformLogAnalyzerErrorMessage from the specified map of raw messages.
func UnmarshalTerraformLogAnalyzerErrorMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TerraformLogAnalyzerErrorMessage)
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TerraformLogAnalyzerSuccessMessage : The success message that is parsed by the terraform log analyzer.
type TerraformLogAnalyzerSuccessMessage struct {
	// The resource type.
	ResourceType *string `json:"resource_type,omitempty"`

	// The time that is taken.
	TimeTaken *string `json:"time-taken,omitempty"`

	// The ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalTerraformLogAnalyzerSuccessMessage unmarshals an instance of TerraformLogAnalyzerSuccessMessage from the specified map of raw messages.
func UnmarshalTerraformLogAnalyzerSuccessMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TerraformLogAnalyzerSuccessMessage)
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time-taken", &obj.TimeTaken)
	if err != nil {
		err = core.SDKErrorf(err, "", "time-taken-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UndeployConfigOptions : The UndeployConfig options.
type UndeployConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
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

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	Definition ProjectConfigDefinitionPatchIntf `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateConfigOptions : Instantiate UpdateConfigOptions
func (*ProjectV1) NewUpdateConfigOptions(projectID string, id string, definition ProjectConfigDefinitionPatchIntf) *UpdateConfigOptions {
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
func (_options *UpdateConfigOptions) SetDefinition(definition ProjectConfigDefinitionPatchIntf) *UpdateConfigOptions {
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

	// The environment definition that is used for updates.
	Definition *EnvironmentDefinitionPropertiesPatch `json:"definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProjectEnvironmentOptions : Instantiate UpdateProjectEnvironmentOptions
func (*ProjectV1) NewUpdateProjectEnvironmentOptions(projectID string, id string, definition *EnvironmentDefinitionPropertiesPatch) *UpdateProjectEnvironmentOptions {
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
func (_options *UpdateProjectEnvironmentOptions) SetDefinition(definition *EnvironmentDefinitionPropertiesPatch) *UpdateProjectEnvironmentOptions {
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

// UpdateStackDefinitionOptions : The UpdateStackDefinition options.
type UpdateStackDefinitionOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
	ID *string `json:"id" validate:"required,ne="`

	// The definition block for a stack definition.
	StackDefinition *StackDefinitionBlockPrototype `json:"stack_definition" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateStackDefinitionOptions : Instantiate UpdateStackDefinitionOptions
func (*ProjectV1) NewUpdateStackDefinitionOptions(projectID string, id string, stackDefinition *StackDefinitionBlockPrototype) *UpdateStackDefinitionOptions {
	return &UpdateStackDefinitionOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
		StackDefinition: stackDefinition,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateStackDefinitionOptions) SetProjectID(projectID string) *UpdateStackDefinitionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateStackDefinitionOptions) SetID(id string) *UpdateStackDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetStackDefinition : Allow user to set StackDefinition
func (_options *UpdateStackDefinitionOptions) SetStackDefinition(stackDefinition *StackDefinitionBlockPrototype) *UpdateStackDefinitionOptions {
	_options.StackDefinition = stackDefinition
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateStackDefinitionOptions) SetHeaders(param map[string]string) *UpdateStackDefinitionOptions {
	options.Headers = param
	return options
}

// ValidateConfigOptions : The ValidateConfig options.
type ValidateConfigOptions struct {
	// The unique project ID.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The unique configuration ID.
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

// ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch : The name and description of a project configuration.
// This model "extends" ProjectConfigDefinitionPatch
type ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch) isaProjectConfigDefinitionPatch() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch unmarshals an instance of ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPatchDAConfigDefinitionPropertiesPatch)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch : The name and description of a project configuration.
// This model "extends" ProjectConfigDefinitionPatch
type ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch struct {
	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch) isaProjectConfigDefinitionPatch() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch unmarshals an instance of ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPatchResourceConfigDefinitionPropertiesPatch)
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch : The name and description of a project configuration.
// This model "extends" ProjectConfigDefinitionPatch
type ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name,omitempty"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch) isaProjectConfigDefinitionPatch() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch unmarshals an instance of ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPatchStackConfigDefinitionPropertiesPatch)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionPrototype
type ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// NewProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype : Instantiate ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype(name string) (_model *ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype, err error) {
	_model = &ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype) isaProjectConfigDefinitionPrototype() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype unmarshals an instance of ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionPrototype
type ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype struct {
	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// NewProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype : Instantiate ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype(name string) (_model *ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype, err error) {
	_model = &ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype) isaProjectConfigDefinitionPrototype() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype unmarshals an instance of ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype)
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionPrototype
type ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// NewProjectConfigDefinitionPrototypeStackConfigDefinitionProperties : Instantiate ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties (Generic Model Constructor)
func (*ProjectV1) NewProjectConfigDefinitionPrototypeStackConfigDefinitionProperties(name string) (_model *ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties, err error) {
	_model = &ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties) isaProjectConfigDefinitionPrototype() bool {
	return true
}

// UnmarshalProjectConfigDefinitionPrototypeStackConfigDefinitionProperties unmarshals an instance of ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionPrototypeStackConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionPrototypeStackConfigDefinitionProperties)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionResponse
type ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// A project configuration description.
	Description *string `json:"description" validate:"required"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse) isaProjectConfigDefinitionResponse() bool {
	return true
}

// UnmarshalProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse unmarshals an instance of ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionResponse
type ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse struct {
	// The CRNs of the resources that are associated with this configuration.
	ResourceCrns []string `json:"resource_crns,omitempty"`

	// A project configuration description.
	Description *string `json:"description" validate:"required"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse) isaProjectConfigDefinitionResponse() bool {
	return true
}

// UnmarshalProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse unmarshals an instance of ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse)
	err = core.UnmarshalPrimitive(m, "resource_crns", &obj.ResourceCrns)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_crns-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigDefinitionResponseStackConfigDefinitionProperties : The description of a project configuration.
// This model "extends" ProjectConfigDefinitionResponse
type ProjectConfigDefinitionResponseStackConfigDefinitionProperties struct {
	// The profile that is required for compliance.
	ComplianceProfile *ProjectComplianceProfile `json:"compliance_profile,omitempty"`

	// A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the
	// catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is
	// required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the
	// Schematics workspace has one.
	// > There are 3 scenarios:
	// > 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.
	// > 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the
	// existing schematics workspace.
	// > 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the
	// specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.
	// > For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform
	// template](/docs/schematics?topic=schematics-sch-create-wks).
	LocatorID *string `json:"locator_id,omitempty"`

	// The member deployabe architectures that are included in your stack.
	Members []StackConfigMember `json:"members,omitempty"`

	// A project configuration description.
	Description *string `json:"description,omitempty"`

	// The configuration name. It's unique within the account across projects and regions.
	Name *string `json:"name" validate:"required"`

	// The ID of the project environment.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.
	Authorizations *ProjectConfigAuth `json:"authorizations,omitempty"`

	// The input variables that are used for configuration definition and environment.
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// The Schematics environment variables to use to deploy the configuration. Settings are only available if they are
	// specified when the configuration is initially created.
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func (*ProjectConfigDefinitionResponseStackConfigDefinitionProperties) isaProjectConfigDefinitionResponse() bool {
	return true
}

// UnmarshalProjectConfigDefinitionResponseStackConfigDefinitionProperties unmarshals an instance of ProjectConfigDefinitionResponseStackConfigDefinitionProperties from the specified map of raw messages.
func UnmarshalProjectConfigDefinitionResponseStackConfigDefinitionProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigDefinitionResponseStackConfigDefinitionProperties)
	err = core.UnmarshalModel(m, "compliance_profile", &obj.ComplianceProfile, UnmarshalProjectComplianceProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "compliance_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locator_id", &obj.LocatorID)
	if err != nil {
		err = core.SDKErrorf(err, "", "locator_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "members", &obj.Members, UnmarshalStackConfigMember)
	if err != nil {
		err = core.SDKErrorf(err, "", "members-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authorizations", &obj.Authorizations, UnmarshalProjectConfigAuth)
	if err != nil {
		err = core.SDKErrorf(err, "", "authorizations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "inputs", &obj.Inputs)
	if err != nil {
		err = core.SDKErrorf(err, "", "inputs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "settings", &obj.Settings)
	if err != nil {
		err = core.SDKErrorf(err, "", "settings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204 : The Code Risk Analyzer logs of the configuration based on Code Risk Analyzer version 2.0.4.
// This model "extends" ProjectConfigMetadataCodeRiskAnalyzerLogs
type ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204 struct {
	// The version of the Code Risk Analyzer logs of the configuration. The metadata for this schema is specific to Code
	// Risk Analyzer version 2.0.4.
	CraVersion *string `json:"cra_version,omitempty"`

	// The schema version of Code Risk Analyzer logs of the configuration.
	SchemaVersion *string `json:"schema_version,omitempty"`

	// The status of the Code Risk Analyzer logs of the configuration.
	Status *string `json:"status,omitempty"`

	// The Code Risk Analyzer logs a summary of the configuration.
	Summary *CodeRiskAnalyzerLogsSummary `json:"summary,omitempty"`

	// A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time
	// format as specified by RFC 3339.
	Timestamp *strfmt.DateTime `json:"timestamp,omitempty"`
}

// Constants associated with the ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204.Status property.
// The status of the Code Risk Analyzer logs of the configuration.
const (
	ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204_Status_Failed = "failed"
	ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204_Status_Passed = "passed"
)

func (*ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204) isaProjectConfigMetadataCodeRiskAnalyzerLogs() bool {
	return true
}

// UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogsVersion204 unmarshals an instance of ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204 from the specified map of raw messages.
func UnmarshalProjectConfigMetadataCodeRiskAnalyzerLogsVersion204(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectConfigMetadataCodeRiskAnalyzerLogsVersion204)
	err = core.UnmarshalPrimitive(m, "cra_version", &obj.CraVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "cra_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "schema_version", &obj.SchemaVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "schema_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "summary", &obj.Summary, UnmarshalCodeRiskAnalyzerLogsSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "summary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timestamp", &obj.Timestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "timestamp-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionExportRequestStackDefinitionExportCatalogRequest : The payload for the stack definition export request to create a product.
// This model "extends" StackDefinitionExportRequest
type StackDefinitionExportRequestStackDefinitionExportCatalogRequest struct {
	// The catalog ID to publish.
	CatalogID *string `json:"catalog_id" validate:"required"`

	// The semver value of this new version of the product.
	TargetVersion *string `json:"target_version,omitempty"`

	// The variation of this new version of the product.
	Variation *string `json:"variation,omitempty"`

	// The product label.
	Label *string `json:"label" validate:"required"`

	// Tags associated with the catalog product.
	Tags []string `json:"tags,omitempty"`
}

// NewStackDefinitionExportRequestStackDefinitionExportCatalogRequest : Instantiate StackDefinitionExportRequestStackDefinitionExportCatalogRequest (Generic Model Constructor)
func (*ProjectV1) NewStackDefinitionExportRequestStackDefinitionExportCatalogRequest(catalogID string, label string) (_model *StackDefinitionExportRequestStackDefinitionExportCatalogRequest, err error) {
	_model = &StackDefinitionExportRequestStackDefinitionExportCatalogRequest{
		CatalogID: core.StringPtr(catalogID),
		Label: core.StringPtr(label),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*StackDefinitionExportRequestStackDefinitionExportCatalogRequest) isaStackDefinitionExportRequest() bool {
	return true
}

// UnmarshalStackDefinitionExportRequestStackDefinitionExportCatalogRequest unmarshals an instance of StackDefinitionExportRequestStackDefinitionExportCatalogRequest from the specified map of raw messages.
func UnmarshalStackDefinitionExportRequestStackDefinitionExportCatalogRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionExportRequestStackDefinitionExportCatalogRequest)
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_version", &obj.TargetVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "variation", &obj.Variation)
	if err != nil {
		err = core.SDKErrorf(err, "", "variation-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StackDefinitionExportRequestStackDefinitionExportProductRequest : The payload for the stack definition export request to create a new product version.
// This model "extends" StackDefinitionExportRequest
type StackDefinitionExportRequestStackDefinitionExportProductRequest struct {
	// The catalog ID to publish.
	CatalogID *string `json:"catalog_id" validate:"required"`

	// The semver value of this new version of the product.
	TargetVersion *string `json:"target_version" validate:"required"`

	// The variation of this new version of the product.
	Variation *string `json:"variation,omitempty"`

	// The product ID to publish.
	ProductID *string `json:"product_id" validate:"required"`
}

// NewStackDefinitionExportRequestStackDefinitionExportProductRequest : Instantiate StackDefinitionExportRequestStackDefinitionExportProductRequest (Generic Model Constructor)
func (*ProjectV1) NewStackDefinitionExportRequestStackDefinitionExportProductRequest(catalogID string, targetVersion string, productID string) (_model *StackDefinitionExportRequestStackDefinitionExportProductRequest, err error) {
	_model = &StackDefinitionExportRequestStackDefinitionExportProductRequest{
		CatalogID: core.StringPtr(catalogID),
		TargetVersion: core.StringPtr(targetVersion),
		ProductID: core.StringPtr(productID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*StackDefinitionExportRequestStackDefinitionExportProductRequest) isaStackDefinitionExportRequest() bool {
	return true
}

// UnmarshalStackDefinitionExportRequestStackDefinitionExportProductRequest unmarshals an instance of StackDefinitionExportRequestStackDefinitionExportProductRequest from the specified map of raw messages.
func UnmarshalStackDefinitionExportRequestStackDefinitionExportProductRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StackDefinitionExportRequestStackDefinitionExportProductRequest)
	err = core.UnmarshalPrimitive(m, "catalog_id", &obj.CatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_version", &obj.TargetVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "variation", &obj.Variation)
	if err != nil {
		err = core.SDKErrorf(err, "", "variation-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "product_id", &obj.ProductID)
	if err != nil {
		err = core.SDKErrorf(err, "", "product_id-error", common.GetComponentInfo())
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
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
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

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListProjectsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var token *string
		token, err = core.GetQueryParam(result.Next.Href, "token")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'token' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = token
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
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetNext() (page []ProjectSummary, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetAll() (allItems []ProjectSummary, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// ProjectEnvironmentsPager can be used to simplify the use of the "ListProjectEnvironments" method.
//
type ProjectEnvironmentsPager struct {
	hasNext bool
	options *ListProjectEnvironmentsOptions
	client  *ProjectV1
	pageContext struct {
		next *string
	}
}

// NewProjectEnvironmentsPager returns a new ProjectEnvironmentsPager instance.
func (project *ProjectV1) NewProjectEnvironmentsPager(options *ListProjectEnvironmentsOptions) (pager *ProjectEnvironmentsPager, err error) {
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListProjectEnvironmentsOptions = *options
	pager = &ProjectEnvironmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  project,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ProjectEnvironmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ProjectEnvironmentsPager) GetNextWithContext(ctx context.Context) (page []Environment, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListProjectEnvironmentsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var token *string
		token, err = core.GetQueryParam(result.Next.Href, "token")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'token' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = token
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Environments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProjectEnvironmentsPager) GetAllWithContext(ctx context.Context) (allItems []Environment, err error) {
	for pager.HasNext() {
		var nextPage []Environment
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
func (pager *ProjectEnvironmentsPager) GetNext() (page []Environment, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProjectEnvironmentsPager) GetAll() (allItems []Environment, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// ConfigsPager can be used to simplify the use of the "ListConfigs" method.
//
type ConfigsPager struct {
	hasNext bool
	options *ListConfigsOptions
	client  *ProjectV1
	pageContext struct {
		next *string
	}
}

// NewConfigsPager returns a new ConfigsPager instance.
func (project *ProjectV1) NewConfigsPager(options *ListConfigsOptions) (pager *ConfigsPager, err error) {
	if options.Token != nil && *options.Token != "" {
		err = core.SDKErrorf(nil, "the 'options.Token' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListConfigsOptions = *options
	pager = &ConfigsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  project,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ConfigsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ConfigsPager) GetNextWithContext(ctx context.Context) (page []ProjectConfigSummary, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Token = pager.pageContext.next

	result, _, err := pager.client.ListConfigsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var token *string
		token, err = core.GetQueryParam(result.Next.Href, "token")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'token' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = token
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Configs

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ConfigsPager) GetAllWithContext(ctx context.Context) (allItems []ProjectConfigSummary, err error) {
	for pager.HasNext() {
		var nextPage []ProjectConfigSummary
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
func (pager *ConfigsPager) GetNext() (page []ProjectConfigSummary, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ConfigsPager) GetAll() (allItems []ProjectConfigSummary, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
