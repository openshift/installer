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
 * IBM OpenAPI SDK Code Generator Version: 3.68.2-ac7def68-20230310-195410
 */

// Package iampolicymanagementv1 : Operations and models for the IamPolicyManagementV1 service
package iampolicymanagementv1

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

// IamPolicyManagementV1 : IAM Policy Management API
//
// API Version: 1.0.1
type IamPolicyManagementV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://iam.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "iam_policy_management"

// IamPolicyManagementV1Options : Service options
type IamPolicyManagementV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewIamPolicyManagementV1UsingExternalConfig : constructs an instance of IamPolicyManagementV1 with passed in options and external configuration.
func NewIamPolicyManagementV1UsingExternalConfig(options *IamPolicyManagementV1Options) (iamPolicyManagement *IamPolicyManagementV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	iamPolicyManagement, err = NewIamPolicyManagementV1(options)
	if err != nil {
		return
	}

	err = iamPolicyManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = iamPolicyManagement.Service.SetServiceURL(options.URL)
	}
	return
}

// NewIamPolicyManagementV1 : constructs an instance of IamPolicyManagementV1 with passed in options.
func NewIamPolicyManagementV1(options *IamPolicyManagementV1Options) (service *IamPolicyManagementV1, err error) {
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

	service = &IamPolicyManagementV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "iamPolicyManagement" suitable for processing requests.
func (iamPolicyManagement *IamPolicyManagementV1) Clone() *IamPolicyManagementV1 {
	if core.IsNil(iamPolicyManagement) {
		return nil
	}
	clone := *iamPolicyManagement
	clone.Service = iamPolicyManagement.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (iamPolicyManagement *IamPolicyManagementV1) SetServiceURL(url string) error {
	return iamPolicyManagement.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (iamPolicyManagement *IamPolicyManagementV1) GetServiceURL() string {
	return iamPolicyManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (iamPolicyManagement *IamPolicyManagementV1) SetDefaultHeaders(headers http.Header) {
	iamPolicyManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (iamPolicyManagement *IamPolicyManagementV1) SetEnableGzipCompression(enableGzip bool) {
	iamPolicyManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (iamPolicyManagement *IamPolicyManagementV1) GetEnableGzipCompression() bool {
	return iamPolicyManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (iamPolicyManagement *IamPolicyManagementV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	iamPolicyManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (iamPolicyManagement *IamPolicyManagementV1) DisableRetries() {
	iamPolicyManagement.Service.DisableRetries()
}

// ListPolicies : Get policies by attributes
// Get policies and filter by attributes. While managing policies, you might want to retrieve policies in the account
// and filter by attribute values. This can be done through query parameters. The following attributes are supported:
// account_id, iam_id, access_group_id, type, service_type, sort, format and state. account_id is a required query
// parameter. Only policies that have the specified attributes and that the caller has read access to are returned. If
// the caller does not have read access to any policies an empty array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicies(listPoliciesOptions *ListPoliciesOptions) (result *PolicyList, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ListPoliciesWithContext(context.Background(), listPoliciesOptions)
}

// ListPoliciesWithContext is an alternate form of the ListPolicies method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListPoliciesWithContext(ctx context.Context, listPoliciesOptions *ListPoliciesOptions) (result *PolicyList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPoliciesOptions, "listPoliciesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listPoliciesOptions, "listPoliciesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listPoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListPolicies")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listPoliciesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listPoliciesOptions.AcceptLanguage))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listPoliciesOptions.AccountID))
	if listPoliciesOptions.IamID != nil {
		builder.AddQuery("iam_id", fmt.Sprint(*listPoliciesOptions.IamID))
	}
	if listPoliciesOptions.AccessGroupID != nil {
		builder.AddQuery("access_group_id", fmt.Sprint(*listPoliciesOptions.AccessGroupID))
	}
	if listPoliciesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listPoliciesOptions.Type))
	}
	if listPoliciesOptions.ServiceType != nil {
		builder.AddQuery("service_type", fmt.Sprint(*listPoliciesOptions.ServiceType))
	}
	if listPoliciesOptions.TagName != nil {
		builder.AddQuery("tag_name", fmt.Sprint(*listPoliciesOptions.TagName))
	}
	if listPoliciesOptions.TagValue != nil {
		builder.AddQuery("tag_value", fmt.Sprint(*listPoliciesOptions.TagValue))
	}
	if listPoliciesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listPoliciesOptions.Sort))
	}
	if listPoliciesOptions.Format != nil {
		builder.AddQuery("format", fmt.Sprint(*listPoliciesOptions.Format))
	}
	if listPoliciesOptions.State != nil {
		builder.AddQuery("state", fmt.Sprint(*listPoliciesOptions.State))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreatePolicy : Create a policy
// Creates a policy to grant access between a subject and a resource. There are two types of policies: **access** and
// **authorization**. A policy administrator might want to create an access policy which grants access to a user,
// service-id, or an access group. They might also want to create an authorization policy and setup access between
// services.
//
// ### Access
//
// To create an access policy, use **`"type": "access"`** in the body. The possible subject attributes are **`iam_id`**
// and **`access_group_id`**. Use the **`iam_id`** subject attribute for assigning access for a user or service-id. Use
// the **`access_group_id`** subject attribute for assigning access for an access group. Assign roles that are supported
// by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). Use only the resource attributes supported by the
// service. To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs).
// The policy resource must include either the **`serviceType`**, **`serviceName`**, **`resourceGroupId`** or
// **`service_group_id`** attribute and the **`accountId`** attribute. The IAM Services group (`IAM`) is a subset of
// account management services that includes the IAM platform services IAM Identity, IAM Access Management, IAM Users
// Management, IAM Groups, and future IAM services. If the subject is a locked service-id, the request will fail.
//
// ### Authorization
//
// Authorization policies are supported by services on a case by case basis. Refer to service documentation to verify
// their support of authorization policies. To create an authorization policy, use **`"type": "authorization"`** in the
// body. The subject attributes must match the supported authorization subjects of the resource. Multiple subject
// attributes might be provided. The following attributes are supported:
//   serviceName, serviceInstance, region, resourceType, resource, accountId Assign roles that are supported by the
// service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`serviceName`** and **`accountId`** attributes.
//
// ### Attribute Operators
//
// Currently, only the `stringEquals` and the `stringMatch` operators are available. Resource attributes may support one
// or both operators. For more information, see [Assigning access by using wildcard
// policies](https://cloud.ibm.com/docs/account?topic=account-wildcard).
//
// ### Attribute Validations
//
// Policy attribute values must be between 1 and 1,000 characters in length. If location related attributes like
// geography, country, metro, region, satellite, and locationvalues are supported by the service, they are validated
// against Global Catalog locations.
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicy(createPolicyOptions *CreatePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.CreatePolicyWithContext(context.Background(), createPolicyOptions)
}

// CreatePolicyWithContext is an alternate form of the CreatePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyWithContext(ctx context.Context, createPolicyOptions *CreatePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPolicyOptions, "createPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createPolicyOptions, "createPolicyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreatePolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createPolicyOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createPolicyOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createPolicyOptions.Type != nil {
		body["type"] = createPolicyOptions.Type
	}
	if createPolicyOptions.Subjects != nil {
		body["subjects"] = createPolicyOptions.Subjects
	}
	if createPolicyOptions.Roles != nil {
		body["roles"] = createPolicyOptions.Roles
	}
	if createPolicyOptions.Resources != nil {
		body["resources"] = createPolicyOptions.Resources
	}
	if createPolicyOptions.Description != nil {
		body["description"] = createPolicyOptions.Description
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplacePolicy : Update a policy
// Update a policy to grant access between a subject and a resource. A policy administrator might want to update an
// existing policy. The policy type cannot be changed (You cannot change an access policy to an authorization policy).
//
// ### Access
//
// To update an access policy, use **`"type": "access"`** in the body. The possible subject attributes are **`iam_id`**
// and **`access_group_id`**. Use the **`iam_id`** subject attribute for assigning access for a user or service-id. Use
// the **`access_group_id`** subject attribute for assigning access for an access group. Assign roles that are supported
// by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). Use only the resource attributes supported by the
// service. To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs).
// The policy resource must include either the **`serviceType`**, **`serviceName`**,  or **`resourceGroupId`** attribute
// and the **`accountId`** attribute.` If the subject is a locked service-id, the request will fail.
//
// ### Authorization
//
// To update an authorization policy, use **`"type": "authorization"`** in the body. The subject attributes must match
// the supported authorization subjects of the resource. Multiple subject attributes might be provided. The following
// attributes are supported:
//   serviceName, serviceInstance, region, resourceType, resource, accountId Assign roles that are supported by the
// service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`serviceName`** and **`accountId`** attributes.
//
// ### Attribute Operators
//
// Currently, only the `stringEquals` and the `stringMatch` operators are available. Resource attributes might support
// one or both operators. For more information, see [Assigning access by using wildcard
// policies](https://cloud.ibm.com/docs/account?topic=account-wildcard).
//
// ### Attribute Validations
//
// Policy attribute values must be between 1 and 1,000 characters in length. If location related attributes like
// geography, country, metro, region, satellite, and locationvalues are supported by the service, they are validated
// against Global Catalog locations.
func (iamPolicyManagement *IamPolicyManagementV1) ReplacePolicy(replacePolicyOptions *ReplacePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ReplacePolicyWithContext(context.Background(), replacePolicyOptions)
}

// ReplacePolicyWithContext is an alternate form of the ReplacePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplacePolicyWithContext(ctx context.Context, replacePolicyOptions *ReplacePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replacePolicyOptions, "replacePolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replacePolicyOptions, "replacePolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"policy_id": *replacePolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replacePolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ReplacePolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replacePolicyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replacePolicyOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replacePolicyOptions.Type != nil {
		body["type"] = replacePolicyOptions.Type
	}
	if replacePolicyOptions.Subjects != nil {
		body["subjects"] = replacePolicyOptions.Subjects
	}
	if replacePolicyOptions.Roles != nil {
		body["roles"] = replacePolicyOptions.Roles
	}
	if replacePolicyOptions.Resources != nil {
		body["resources"] = replacePolicyOptions.Resources
	}
	if replacePolicyOptions.Description != nil {
		body["description"] = replacePolicyOptions.Description
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetPolicy : Retrieve a policy by ID
// Retrieve a policy by providing a policy ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicy(getPolicyOptions *GetPolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.GetPolicyWithContext(context.Background(), getPolicyOptions)
}

// GetPolicyWithContext is an alternate form of the GetPolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyWithContext(ctx context.Context, getPolicyOptions *GetPolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyOptions, "getPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPolicyOptions, "getPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"policy_id": *getPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeletePolicy : Delete a policy by ID
// Delete a policy by providing a policy ID. A policy cannot be deleted if the subject ID contains a locked service ID.
// If the subject of the policy is a locked service-id, the request will fail.
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicy(deletePolicyOptions *DeletePolicyOptions) (response *core.DetailedResponse, err error) {
	return iamPolicyManagement.DeletePolicyWithContext(context.Background(), deletePolicyOptions)
}

// DeletePolicyWithContext is an alternate form of the DeletePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyWithContext(ctx context.Context, deletePolicyOptions *DeletePolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePolicyOptions, "deletePolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deletePolicyOptions, "deletePolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"policy_id": *deletePolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deletePolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeletePolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)

	return
}

// UpdatePolicyState : Restore a deleted policy by ID
// Restore a policy that has recently been deleted. A policy administrator might want to restore a deleted policy. To
// restore a policy, use **`"state": "active"`** in the body.
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyState(updatePolicyStateOptions *UpdatePolicyStateOptions) (result *Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.UpdatePolicyStateWithContext(context.Background(), updatePolicyStateOptions)
}

// UpdatePolicyStateWithContext is an alternate form of the UpdatePolicyState method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyStateWithContext(ctx context.Context, updatePolicyStateOptions *UpdatePolicyStateOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePolicyStateOptions, "updatePolicyStateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePolicyStateOptions, "updatePolicyStateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"policy_id": *updatePolicyStateOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updatePolicyStateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "UpdatePolicyState")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updatePolicyStateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updatePolicyStateOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updatePolicyStateOptions.State != nil {
		body["state"] = updatePolicyStateOptions.State
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRoles : Get roles by filters
// Get roles based on the filters. While managing roles, you may want to retrieve roles and filter by usages. This can
// be done through query parameters. Currently, we only support the following attributes: account_id, service_name,
// service_group_id, source_service_name and policy_type. Both service_name and service_group_id attributes are mutually
// exclusive. Only roles that match the filter and that the caller has read access to are returned. If the caller does
// not have read access to any roles an empty array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListRoles(listRolesOptions *ListRolesOptions) (result *RoleList, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ListRolesWithContext(context.Background(), listRolesOptions)
}

// ListRolesWithContext is an alternate form of the ListRoles method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListRolesWithContext(ctx context.Context, listRolesOptions *ListRolesOptions) (result *RoleList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRolesOptions, "listRolesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRolesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListRoles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listRolesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listRolesOptions.AcceptLanguage))
	}

	if listRolesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listRolesOptions.AccountID))
	}
	if listRolesOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listRolesOptions.ServiceName))
	}
	if listRolesOptions.SourceServiceName != nil {
		builder.AddQuery("source_service_name", fmt.Sprint(*listRolesOptions.SourceServiceName))
	}
	if listRolesOptions.PolicyType != nil {
		builder.AddQuery("policy_type", fmt.Sprint(*listRolesOptions.PolicyType))
	}
	if listRolesOptions.ServiceGroupID != nil {
		builder.AddQuery("service_group_id", fmt.Sprint(*listRolesOptions.ServiceGroupID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoleList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateRole : Create a role
// Creates a custom role for a specific service within the account. An account owner or a user assigned the
// Administrator role on the Role management service can create a custom role. Any number of actions for a single
// service can be mapped to the new role, but there must be at least one service-defined action to successfully create
// the new role.
func (iamPolicyManagement *IamPolicyManagementV1) CreateRole(createRoleOptions *CreateRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.CreateRoleWithContext(context.Background(), createRoleOptions)
}

// CreateRoleWithContext is an alternate form of the CreateRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreateRoleWithContext(ctx context.Context, createRoleOptions *CreateRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRoleOptions, "createRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRoleOptions, "createRoleOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreateRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createRoleOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createRoleOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createRoleOptions.DisplayName != nil {
		body["display_name"] = createRoleOptions.DisplayName
	}
	if createRoleOptions.Actions != nil {
		body["actions"] = createRoleOptions.Actions
	}
	if createRoleOptions.Name != nil {
		body["name"] = createRoleOptions.Name
	}
	if createRoleOptions.AccountID != nil {
		body["account_id"] = createRoleOptions.AccountID
	}
	if createRoleOptions.ServiceName != nil {
		body["service_name"] = createRoleOptions.ServiceName
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceRole : Update a role
// Update a custom role. A role administrator might want to update an existing role by updating the display name,
// description, or the actions that are mapped to the role. The name, account_id, and service_name can't be changed.
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceRole(replaceRoleOptions *ReplaceRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ReplaceRoleWithContext(context.Background(), replaceRoleOptions)
}

// ReplaceRoleWithContext is an alternate form of the ReplaceRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceRoleWithContext(ctx context.Context, replaceRoleOptions *ReplaceRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRoleOptions, "replaceRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceRoleOptions, "replaceRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"role_id": *replaceRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles/{role_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ReplaceRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceRoleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceRoleOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replaceRoleOptions.DisplayName != nil {
		body["display_name"] = replaceRoleOptions.DisplayName
	}
	if replaceRoleOptions.Actions != nil {
		body["actions"] = replaceRoleOptions.Actions
	}
	if replaceRoleOptions.Description != nil {
		body["description"] = replaceRoleOptions.Description
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRole : Retrieve a role by ID
// Retrieve a role by providing a role ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetRole(getRoleOptions *GetRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.GetRoleWithContext(context.Background(), getRoleOptions)
}

// GetRoleWithContext is an alternate form of the GetRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetRoleWithContext(ctx context.Context, getRoleOptions *GetRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRoleOptions, "getRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRoleOptions, "getRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"role_id": *getRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles/{role_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRole : Delete a role by ID
// Delete a role by providing a role ID.
func (iamPolicyManagement *IamPolicyManagementV1) DeleteRole(deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	return iamPolicyManagement.DeleteRoleWithContext(context.Background(), deleteRoleOptions)
}

// DeleteRoleWithContext is an alternate form of the DeleteRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeleteRoleWithContext(ctx context.Context, deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRoleOptions, "deleteRoleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRoleOptions, "deleteRoleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"role_id": *deleteRoleOptions.RoleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles/{role_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRoleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeleteRole")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)

	return
}

// ListV2Policies : Get policies by attributes
// Get policies and filter by attributes. While managing policies, you might want to retrieve policies in the account
// and filter by attribute values. This can be done through query parameters. The following attributes are supported:
// account_id, iam_id, access_group_id, type, service_type, sort, format and state. account_id is a required query
// parameter. Only policies that have the specified attributes and that the caller has read access to are returned. If
// the caller does not have read access to any policies an empty array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListV2Policies(listV2PoliciesOptions *ListV2PoliciesOptions) (result *V2PolicyCollection, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ListV2PoliciesWithContext(context.Background(), listV2PoliciesOptions)
}

// ListV2PoliciesWithContext is an alternate form of the ListV2Policies method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListV2PoliciesWithContext(ctx context.Context, listV2PoliciesOptions *ListV2PoliciesOptions) (result *V2PolicyCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listV2PoliciesOptions, "listV2PoliciesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listV2PoliciesOptions, "listV2PoliciesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listV2PoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListV2Policies")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listV2PoliciesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listV2PoliciesOptions.AcceptLanguage))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listV2PoliciesOptions.AccountID))
	if listV2PoliciesOptions.IamID != nil {
		builder.AddQuery("iam_id", fmt.Sprint(*listV2PoliciesOptions.IamID))
	}
	if listV2PoliciesOptions.AccessGroupID != nil {
		builder.AddQuery("access_group_id", fmt.Sprint(*listV2PoliciesOptions.AccessGroupID))
	}
	if listV2PoliciesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listV2PoliciesOptions.Type))
	}
	if listV2PoliciesOptions.ServiceType != nil {
		builder.AddQuery("service_type", fmt.Sprint(*listV2PoliciesOptions.ServiceType))
	}
	if listV2PoliciesOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listV2PoliciesOptions.ServiceName))
	}
	if listV2PoliciesOptions.ServiceGroupID != nil {
		builder.AddQuery("service_group_id", fmt.Sprint(*listV2PoliciesOptions.ServiceGroupID))
	}
	if listV2PoliciesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listV2PoliciesOptions.Sort))
	}
	if listV2PoliciesOptions.Format != nil {
		builder.AddQuery("format", fmt.Sprint(*listV2PoliciesOptions.Format))
	}
	if listV2PoliciesOptions.State != nil {
		builder.AddQuery("state", fmt.Sprint(*listV2PoliciesOptions.State))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2PolicyCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateV2Policy : Create a policy
// Creates a policy to grant access between a subject and a resource. Currently, there is one type of a v2/policy:
// **access**. A policy administrator might want to create an access policy that grants access to a user, service-id, or
// an access group.
//
// ### Access
//
// To create an access policy, use **`"type": "access"`** in the body. The supported subject attributes are **`iam_id`**
// and **`access_group_id`**. Use the **`iam_id`** subject attribute to assign access to a user or service-id. Use the
// **`access_group_id`** subject attribute to assign access to an access group. Assign roles that are supported by the
// service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). Use only the resource attributes supported by the
// service. To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs).
// The policy resource must include either the **`serviceType`**, **`serviceName`**, **`resourceGroupId`** or
// **`service_group_id`** attribute and the **`accountId`** attribute. In the rule field, you can specify a single
// condition by using **`key`**, **`value`**, and condition **`operator`**, or a set of **`conditions`** with a
// combination **`operator`**.  The possible combination operators are **`and`** and **`or`**. Combine conditions to
// specify a time-based restriction (e.g., access only during business hours, during the Monday-Friday work week). For
// example, a policy can grant access Monday-Friday, 9:00am-5:00pm using the following rule:
// ```json
//   "rule": {
//     "operator": "and",
//     "conditions": [{
//       "key": "{{environment.attributes.day_of_week}}",
//       "operator": "dayOfWeekAnyOf",
//       "value": ["1+00:00", "2+00:00", "3+00:00", "4+00:00", "5+00:00"]
//     },
//       "key": "{{environment.attributes.current_time}}",
//       "operator": "timeGreaterThanOrEquals",
//       "value": "09:00:00+00:00"
//     },
//       "key": "{{environment.attributes.current_time}}",
//       "operator": "timeLessThanOrEquals",
//       "value": "17:00:00+00:00"
//     }]
//   }
// ``` You can use the following operators in the **`key`** and **`value`** pair:
// ```
//   'timeLessThan', 'timeLessThanOrEquals', 'timeGreaterThan', 'timeGreaterThanOrEquals',
//   'dateTimeLessThan', 'dateTimeLessThanOrEquals', 'dateTimeGreaterThan', 'dateTimeGreaterThanOrEquals',
//   'dayOfWeekEquals', 'dayOfWeekAnyOf',
// ```
//
// The pattern field that matches the rule is required when rule is provided. For the business hour rule example above,
// the **`pattern`** is **`"time-based-conditions:weekly"`**. For more information, see [Time-based conditions
// operators](https://cloud.ibm.com/docs/account?topic=account-iam-condition-properties&interface=ui#policy-condition-properties)
// and
// [Limiting access with time-based
// conditions](https://cloud.ibm.com/docs/account?topic=account-iam-time-based&interface=ui). If the subject is a locked
// service-id, the request will fail.
//
// ### Attribute Operators
//
// Currently, only the `stringEquals`, `stringMatch`, and `stringEquals` operators are available. For more information,
// see [Assigning access by using wildcard policies](https://cloud.ibm.com/docs/account?topic=account-wildcard).
//
// ### Attribute Validations
//
// Policy attribute values must be between 1 and 1,000 characters in length. If location related attributes like
// geography, country, metro, region, satellite, and locationvalues are supported by the service, they are validated
// against Global Catalog locations.
func (iamPolicyManagement *IamPolicyManagementV1) CreateV2Policy(createV2PolicyOptions *CreateV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.CreateV2PolicyWithContext(context.Background(), createV2PolicyOptions)
}

// CreateV2PolicyWithContext is an alternate form of the CreateV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreateV2PolicyWithContext(ctx context.Context, createV2PolicyOptions *CreateV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createV2PolicyOptions, "createV2PolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createV2PolicyOptions, "createV2PolicyOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createV2PolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreateV2Policy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createV2PolicyOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createV2PolicyOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createV2PolicyOptions.Control != nil {
		body["control"] = createV2PolicyOptions.Control
	}
	if createV2PolicyOptions.Type != nil {
		body["type"] = createV2PolicyOptions.Type
	}
	if createV2PolicyOptions.Description != nil {
		body["description"] = createV2PolicyOptions.Description
	}
	if createV2PolicyOptions.Subject != nil {
		body["subject"] = createV2PolicyOptions.Subject
	}
	if createV2PolicyOptions.Resource != nil {
		body["resource"] = createV2PolicyOptions.Resource
	}
	if createV2PolicyOptions.Pattern != nil {
		body["pattern"] = createV2PolicyOptions.Pattern
	}
	if createV2PolicyOptions.Rule != nil {
		body["rule"] = createV2PolicyOptions.Rule
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2Policy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceV2Policy : Update a policy
// Update a policy to grant access between a subject and a resource. A policy administrator might want to update an
// existing policy.
//
// ### Access
//
// To create an access policy, use **`"type": "access"`** in the body. The supported subject attributes are **`iam_id`**
// and **`access_group_id`**. Use the **`iam_id`** subject attribute to assign access to a user or service-id. Use the
// **`access_group_id`** subject attribute to assign access to an access group. Assign roles that are supported by the
// service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). Use only the resource attributes supported by the
// service. To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs).
// The policy resource must include either the **`serviceType`**, **`serviceName`**, **`resourceGroupId`** or
// **`service_group_id`** attribute and the **`accountId`** attribute. In the rule field, you can specify a single
// condition by using **`key`**, **`value`**, and condition **`operator`**, or a set of **`conditions`** with a
// combination **`operator`**.  The possible combination operators are **`and`** and **`or`**. Combine conditions to
// specify a time-based restriction (e.g., access only during business hours, during the Monday-Friday work week). For
// example, a policy can grant access Monday-Friday, 9:00am-5:00pm using the following rule:
// ```json
//   "rule": {
//     "operator": "and",
//     "conditions": [{
//       "key": "{{environment.attributes.day_of_week}}",
//       "operator": "dayOfWeekAnyOf",
//       "value": ["1+00:00", "2+00:00", "3+00:00", "4+00:00", "5+00:00"]
//     },
//       "key": "{{environment.attributes.current_time}}",
//       "operator": "timeGreaterThanOrEquals",
//       "value": "09:00:00+00:00"
//     },
//       "key": "{{environment.attributes.current_time}}",
//       "operator": "timeLessThanOrEquals",
//       "value": "17:00:00+00:00"
//     }]
//   }
// ``` You can use the following operators in the **`key`**, **`value`** pair:
// ```
//   'timeLessThan', 'timeLessThanOrEquals', 'timeGreaterThan', 'timeGreaterThanOrEquals',
//   'dateTimeLessThan', 'dateTimeLessThanOrEquals', 'dateTimeGreaterThan', 'dateTimeGreaterThanOrEquals',
//   'dayOfWeekEquals', 'dayOfWeekAnyOf',
// ``` The pattern field that matches the rule is required when rule is provided. For the business hour rule example
// above, the **`pattern`** is **`"time-based-conditions:weekly"`**. For more information, see [Time-based conditions
// operators](https://cloud.ibm.com/docs/account?topic=account-iam-condition-properties&interface=ui#policy-condition-properties)
// and
// [Limiting access with time-based
// conditions](https://cloud.ibm.com/docs/account?topic=account-iam-time-based&interface=ui).
// ### Attribute Operators
//
// Currently, only the `stringEquals`, `stringMatch`, and `stringEquals` operators are available. For more information,
// see [Assigning access by using wildcard policies](https://cloud.ibm.com/docs/account?topic=account-wildcard).
//
// ### Attribute Validations
//
// Policy attribute values must be between 1 and 1,000 characters in length. If location related attributes like
// geography, country, metro, region, satellite, and locationvalues are supported by the service, they are validated
// against Global Catalog locations.
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceV2Policy(replaceV2PolicyOptions *ReplaceV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.ReplaceV2PolicyWithContext(context.Background(), replaceV2PolicyOptions)
}

// ReplaceV2PolicyWithContext is an alternate form of the ReplaceV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceV2PolicyWithContext(ctx context.Context, replaceV2PolicyOptions *ReplaceV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceV2PolicyOptions, "replaceV2PolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceV2PolicyOptions, "replaceV2PolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *replaceV2PolicyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceV2PolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ReplaceV2Policy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceV2PolicyOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceV2PolicyOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replaceV2PolicyOptions.Control != nil {
		body["control"] = replaceV2PolicyOptions.Control
	}
	if replaceV2PolicyOptions.Type != nil {
		body["type"] = replaceV2PolicyOptions.Type
	}
	if replaceV2PolicyOptions.Description != nil {
		body["description"] = replaceV2PolicyOptions.Description
	}
	if replaceV2PolicyOptions.Subject != nil {
		body["subject"] = replaceV2PolicyOptions.Subject
	}
	if replaceV2PolicyOptions.Resource != nil {
		body["resource"] = replaceV2PolicyOptions.Resource
	}
	if replaceV2PolicyOptions.Pattern != nil {
		body["pattern"] = replaceV2PolicyOptions.Pattern
	}
	if replaceV2PolicyOptions.Rule != nil {
		body["rule"] = replaceV2PolicyOptions.Rule
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2Policy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetV2Policy : Retrieve a policy by ID
// Retrieve a policy by providing a policy ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetV2Policy(getV2PolicyOptions *GetV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	return iamPolicyManagement.GetV2PolicyWithContext(context.Background(), getV2PolicyOptions)
}

// GetV2PolicyWithContext is an alternate form of the GetV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetV2PolicyWithContext(ctx context.Context, getV2PolicyOptions *GetV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getV2PolicyOptions, "getV2PolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getV2PolicyOptions, "getV2PolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getV2PolicyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getV2PolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetV2Policy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getV2PolicyOptions.Format != nil {
		builder.AddQuery("format", fmt.Sprint(*getV2PolicyOptions.Format))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2Policy)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteV2Policy : Delete a policy by ID
// Delete a policy by providing a policy ID. A policy cannot be deleted if the subject ID contains a locked service ID.
// If the subject of the policy is a locked service-id, the request will fail.
func (iamPolicyManagement *IamPolicyManagementV1) DeleteV2Policy(deleteV2PolicyOptions *DeleteV2PolicyOptions) (response *core.DetailedResponse, err error) {
	return iamPolicyManagement.DeleteV2PolicyWithContext(context.Background(), deleteV2PolicyOptions)
}

// DeleteV2PolicyWithContext is an alternate form of the DeleteV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeleteV2PolicyWithContext(ctx context.Context, deleteV2PolicyOptions *DeleteV2PolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteV2PolicyOptions, "deleteV2PolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteV2PolicyOptions, "deleteV2PolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteV2PolicyOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteV2PolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeleteV2Policy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)

	return
}

// Control : Specifies the type of access granted by the policy.
type Control struct {
	// Permission granted by the policy.
	Grant *Grant `json:"grant" validate:"required"`
}

// NewControl : Instantiate Control (Generic Model Constructor)
func (*IamPolicyManagementV1) NewControl(grant *Grant) (_model *Control, err error) {
	_model = &Control{
		Grant: grant,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalControl unmarshals an instance of Control from the specified map of raw messages.
func UnmarshalControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Control)
	err = core.UnmarshalModel(m, "grant", &obj.Grant, UnmarshalGrant)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlResponse : ControlResponse struct
// Models which "extend" this model:
// - ControlResponseControl
// - ControlResponseControlWithEnrichedRoles
type ControlResponse struct {
	// Permission granted by the policy.
	Grant *Grant `json:"grant,omitempty"`
}
func (*ControlResponse) isaControlResponse() bool {
	return true
}

type ControlResponseIntf interface {
	isaControlResponse() bool
}

// UnmarshalControlResponse unmarshals an instance of ControlResponse from the specified map of raw messages.
func UnmarshalControlResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlResponse)
	err = core.UnmarshalModel(m, "grant", &obj.Grant, UnmarshalGrant)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreatePolicyOptions : The CreatePolicy options.
type CreatePolicyOptions struct {
	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// The subjects associated with a policy.
	Subjects []PolicySubject `json:"subjects" validate:"required"`

	// A set of role cloud resource names (CRNs) granted by the policy.
	Roles []PolicyRole `json:"roles" validate:"required"`

	// The resources associated with a policy.
	Resources []PolicyResource `json:"resources" validate:"required"`

	// Customer-defined description.
	Description *string `json:"description,omitempty"`

	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreatePolicyOptions : Instantiate CreatePolicyOptions
func (*IamPolicyManagementV1) NewCreatePolicyOptions(typeVar string, subjects []PolicySubject, roles []PolicyRole, resources []PolicyResource) *CreatePolicyOptions {
	return &CreatePolicyOptions{
		Type: core.StringPtr(typeVar),
		Subjects: subjects,
		Roles: roles,
		Resources: resources,
	}
}

// SetType : Allow user to set Type
func (_options *CreatePolicyOptions) SetType(typeVar string) *CreatePolicyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSubjects : Allow user to set Subjects
func (_options *CreatePolicyOptions) SetSubjects(subjects []PolicySubject) *CreatePolicyOptions {
	_options.Subjects = subjects
	return _options
}

// SetRoles : Allow user to set Roles
func (_options *CreatePolicyOptions) SetRoles(roles []PolicyRole) *CreatePolicyOptions {
	_options.Roles = roles
	return _options
}

// SetResources : Allow user to set Resources
func (_options *CreatePolicyOptions) SetResources(resources []PolicyResource) *CreatePolicyOptions {
	_options.Resources = resources
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreatePolicyOptions) SetDescription(description string) *CreatePolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreatePolicyOptions) SetAcceptLanguage(acceptLanguage string) *CreatePolicyOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePolicyOptions) SetHeaders(param map[string]string) *CreatePolicyOptions {
	options.Headers = param
	return options
}

// CreateRoleOptions : The CreateRole options.
type CreateRoleOptions struct {
	// The display name of the role that is shown in the console.
	DisplayName *string `json:"display_name" validate:"required"`

	// The actions of the role. For more information, see [IAM roles and
	// actions](https://cloud.ibm.com/docs/account?topic=account-iam-service-roles-actions).
	Actions []string `json:"actions" validate:"required"`

	// The name of the role that is used in the CRN. Can only be alphanumeric and has to be capitalized.
	Name *string `json:"name" validate:"required"`

	// The account GUID.
	AccountID *string `json:"account_id" validate:"required"`

	// The service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The description of the role.
	Description *string `json:"description,omitempty"`

	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRoleOptions : Instantiate CreateRoleOptions
func (*IamPolicyManagementV1) NewCreateRoleOptions(displayName string, actions []string, name string, accountID string, serviceName string) *CreateRoleOptions {
	return &CreateRoleOptions{
		DisplayName: core.StringPtr(displayName),
		Actions: actions,
		Name: core.StringPtr(name),
		AccountID: core.StringPtr(accountID),
		ServiceName: core.StringPtr(serviceName),
	}
}

// SetDisplayName : Allow user to set DisplayName
func (_options *CreateRoleOptions) SetDisplayName(displayName string) *CreateRoleOptions {
	_options.DisplayName = core.StringPtr(displayName)
	return _options
}

// SetActions : Allow user to set Actions
func (_options *CreateRoleOptions) SetActions(actions []string) *CreateRoleOptions {
	_options.Actions = actions
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateRoleOptions) SetName(name string) *CreateRoleOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateRoleOptions) SetAccountID(accountID string) *CreateRoleOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *CreateRoleOptions) SetServiceName(serviceName string) *CreateRoleOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateRoleOptions) SetDescription(description string) *CreateRoleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateRoleOptions) SetAcceptLanguage(acceptLanguage string) *CreateRoleOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRoleOptions) SetHeaders(param map[string]string) *CreateRoleOptions {
	options.Headers = param
	return options
}

// CreateV2PolicyOptions : The CreateV2Policy options.
type CreateV2PolicyOptions struct {
	// Specifies the type of access granted by the policy.
	Control *Control `json:"control" validate:"required"`

	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Allows the customer to use their own words to record the purpose/context related to a policy.
	Description *string `json:"description,omitempty"`

	// The subject attributes for whom the policy grants access.
	Subject *V2PolicySubject `json:"subject,omitempty"`

	// The resource attributes to which the policy grants access.
	Resource *V2PolicyResource `json:"resource,omitempty"`

	// Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or
	// 'time-based-conditions:weekly:custom-hours'.
	Pattern *string `json:"pattern,omitempty"`

	// Additional access conditions associated with the policy.
	Rule V2PolicyRuleIntf `json:"rule,omitempty"`

	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateV2PolicyOptions.Type property.
// The policy type; either 'access' or 'authorization'.
const (
	CreateV2PolicyOptionsTypeAccessConst = "access"
	CreateV2PolicyOptionsTypeAuthorizationConst = "authorization"
)

// NewCreateV2PolicyOptions : Instantiate CreateV2PolicyOptions
func (*IamPolicyManagementV1) NewCreateV2PolicyOptions(control *Control, typeVar string) *CreateV2PolicyOptions {
	return &CreateV2PolicyOptions{
		Control: control,
		Type: core.StringPtr(typeVar),
	}
}

// SetControl : Allow user to set Control
func (_options *CreateV2PolicyOptions) SetControl(control *Control) *CreateV2PolicyOptions {
	_options.Control = control
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateV2PolicyOptions) SetType(typeVar string) *CreateV2PolicyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateV2PolicyOptions) SetDescription(description string) *CreateV2PolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetSubject : Allow user to set Subject
func (_options *CreateV2PolicyOptions) SetSubject(subject *V2PolicySubject) *CreateV2PolicyOptions {
	_options.Subject = subject
	return _options
}

// SetResource : Allow user to set Resource
func (_options *CreateV2PolicyOptions) SetResource(resource *V2PolicyResource) *CreateV2PolicyOptions {
	_options.Resource = resource
	return _options
}

// SetPattern : Allow user to set Pattern
func (_options *CreateV2PolicyOptions) SetPattern(pattern string) *CreateV2PolicyOptions {
	_options.Pattern = core.StringPtr(pattern)
	return _options
}

// SetRule : Allow user to set Rule
func (_options *CreateV2PolicyOptions) SetRule(rule V2PolicyRuleIntf) *CreateV2PolicyOptions {
	_options.Rule = rule
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateV2PolicyOptions) SetAcceptLanguage(acceptLanguage string) *CreateV2PolicyOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateV2PolicyOptions) SetHeaders(param map[string]string) *CreateV2PolicyOptions {
	options.Headers = param
	return options
}

// DeletePolicyOptions : The DeletePolicy options.
type DeletePolicyOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeletePolicyOptions : Instantiate DeletePolicyOptions
func (*IamPolicyManagementV1) NewDeletePolicyOptions(policyID string) *DeletePolicyOptions {
	return &DeletePolicyOptions{
		PolicyID: core.StringPtr(policyID),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *DeletePolicyOptions) SetPolicyID(policyID string) *DeletePolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePolicyOptions) SetHeaders(param map[string]string) *DeletePolicyOptions {
	options.Headers = param
	return options
}

// DeleteRoleOptions : The DeleteRole options.
type DeleteRoleOptions struct {
	// The role ID.
	RoleID *string `json:"role_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRoleOptions : Instantiate DeleteRoleOptions
func (*IamPolicyManagementV1) NewDeleteRoleOptions(roleID string) *DeleteRoleOptions {
	return &DeleteRoleOptions{
		RoleID: core.StringPtr(roleID),
	}
}

// SetRoleID : Allow user to set RoleID
func (_options *DeleteRoleOptions) SetRoleID(roleID string) *DeleteRoleOptions {
	_options.RoleID = core.StringPtr(roleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRoleOptions) SetHeaders(param map[string]string) *DeleteRoleOptions {
	options.Headers = param
	return options
}

// DeleteV2PolicyOptions : The DeleteV2Policy options.
type DeleteV2PolicyOptions struct {
	// The policy ID.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteV2PolicyOptions : Instantiate DeleteV2PolicyOptions
func (*IamPolicyManagementV1) NewDeleteV2PolicyOptions(id string) *DeleteV2PolicyOptions {
	return &DeleteV2PolicyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteV2PolicyOptions) SetID(id string) *DeleteV2PolicyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteV2PolicyOptions) SetHeaders(param map[string]string) *DeleteV2PolicyOptions {
	options.Headers = param
	return options
}

// EnrichedRoles : A role associated with a policy with additional information (display_name, description, actions) when
// `format=display`.
type EnrichedRoles struct {
	// The role Cloud Resource Name (CRN) granted by the policy. Example CRN: 'crn:v1:bluemix:public:iam::::role:Editor'.
	RoleID *string `json:"role_id" validate:"required"`

	// The service defined (or user defined if a custom role) display name of the role.
	DisplayName *string `json:"display_name,omitempty"`

	// The service defined (or user defined if a custom role) description of the role.
	Description *string `json:"description,omitempty"`

	// The actions of the role. For more information, see [IAM roles and
	// actions](https://cloud.ibm.com/docs/account?topic=account-iam-service-roles-actions).
	Actions []RoleAction `json:"actions" validate:"required"`
}

// UnmarshalEnrichedRoles unmarshals an instance of EnrichedRoles from the specified map of raw messages.
func UnmarshalEnrichedRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnrichedRoles)
	err = core.UnmarshalPrimitive(m, "role_id", &obj.RoleID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalRoleAction)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetPolicyOptions : The GetPolicy options.
type GetPolicyOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPolicyOptions : Instantiate GetPolicyOptions
func (*IamPolicyManagementV1) NewGetPolicyOptions(policyID string) *GetPolicyOptions {
	return &GetPolicyOptions{
		PolicyID: core.StringPtr(policyID),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *GetPolicyOptions) SetPolicyID(policyID string) *GetPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyOptions) SetHeaders(param map[string]string) *GetPolicyOptions {
	options.Headers = param
	return options
}

// GetRoleOptions : The GetRole options.
type GetRoleOptions struct {
	// The role ID.
	RoleID *string `json:"role_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRoleOptions : Instantiate GetRoleOptions
func (*IamPolicyManagementV1) NewGetRoleOptions(roleID string) *GetRoleOptions {
	return &GetRoleOptions{
		RoleID: core.StringPtr(roleID),
	}
}

// SetRoleID : Allow user to set RoleID
func (_options *GetRoleOptions) SetRoleID(roleID string) *GetRoleOptions {
	_options.RoleID = core.StringPtr(roleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRoleOptions) SetHeaders(param map[string]string) *GetRoleOptions {
	options.Headers = param
	return options
}

// GetV2PolicyOptions : The GetV2Policy options.
type GetV2PolicyOptions struct {
	// The policy ID.
	ID *string `json:"id" validate:"required,ne="`

	// Include additional data for policy returned
	// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
	// it has done so
	// * `display` - returns the list of all actions included in each of the policy roles and translations for all relevant
	// fields.
	Format *string `json:"format,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetV2PolicyOptions.Format property.
// Include additional data for policy returned
// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
// it has done so
// * `display` - returns the list of all actions included in each of the policy roles and translations for all relevant
// fields.
const (
	GetV2PolicyOptionsFormatDisplayConst = "display"
	GetV2PolicyOptionsFormatIncludeLastPermitConst = "include_last_permit"
)

// NewGetV2PolicyOptions : Instantiate GetV2PolicyOptions
func (*IamPolicyManagementV1) NewGetV2PolicyOptions(id string) *GetV2PolicyOptions {
	return &GetV2PolicyOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetV2PolicyOptions) SetID(id string) *GetV2PolicyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *GetV2PolicyOptions) SetFormat(format string) *GetV2PolicyOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetV2PolicyOptions) SetHeaders(param map[string]string) *GetV2PolicyOptions {
	options.Headers = param
	return options
}

// Grant : Permission granted by the policy.
type Grant struct {
	// A set of role cloud resource names (CRNs) granted by the policy.
	Roles []Roles `json:"roles" validate:"required"`
}

// NewGrant : Instantiate Grant (Generic Model Constructor)
func (*IamPolicyManagementV1) NewGrant(roles []Roles) (_model *Grant, err error) {
	_model = &Grant{
		Roles: roles,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalGrant unmarshals an instance of Grant from the specified map of raw messages.
func UnmarshalGrant(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Grant)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalRoles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GrantWithEnrichedRoles : Permission granted by the policy with translated roles and additional role information.
type GrantWithEnrichedRoles struct {
	// A set of roles granted by the policy.
	Roles []EnrichedRoles `json:"roles" validate:"required"`
}

// UnmarshalGrantWithEnrichedRoles unmarshals an instance of GrantWithEnrichedRoles from the specified map of raw messages.
func UnmarshalGrantWithEnrichedRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GrantWithEnrichedRoles)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalEnrichedRoles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListPoliciesOptions : The ListPolicies options.
type ListPoliciesOptions struct {
	// The account GUID that the policies belong to.
	AccountID *string `json:"account_id" validate:"required"`

	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Optional IAM ID used to identify the subject.
	IamID *string `json:"iam_id,omitempty"`

	// Optional access group id.
	AccessGroupID *string `json:"access_group_id,omitempty"`

	// Optional type of policy.
	Type *string `json:"type,omitempty"`

	// Optional type of service.
	ServiceType *string `json:"service_type,omitempty"`

	// Optional name of the access tag in the policy.
	TagName *string `json:"tag_name,omitempty"`

	// Optional value of the access tag in the policy.
	TagValue *string `json:"tag_value,omitempty"`

	// Optional top level policy field to sort results. Ascending sort is default. Descending sort available by prepending
	// '-' to field. Example '-last_modified_at'.
	Sort *string `json:"sort,omitempty"`

	// Include additional data per policy returned
	// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
	// it has done so
	// * `display` - returns the list of all actions included in each of the policy roles.
	Format *string `json:"format,omitempty"`

	// The state of the policy.
	// * `active` - returns active policies
	// * `deleted` - returns non-active policies.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListPoliciesOptions.Type property.
// Optional type of policy.
const (
	ListPoliciesOptionsTypeAccessConst = "access"
	ListPoliciesOptionsTypeAuthorizationConst = "authorization"
)

// Constants associated with the ListPoliciesOptions.ServiceType property.
// Optional type of service.
const (
	ListPoliciesOptionsServiceTypePlatformServiceConst = "platform_service"
	ListPoliciesOptionsServiceTypeServiceConst = "service"
)

// Constants associated with the ListPoliciesOptions.Sort property.
// Optional top level policy field to sort results. Ascending sort is default. Descending sort available by prepending
// '-' to field. Example '-last_modified_at'.
const (
	ListPoliciesOptionsSortCreatedAtConst = "created_at"
	ListPoliciesOptionsSortCreatedByIDConst = "created_by_id"
	ListPoliciesOptionsSortHrefConst = "href"
	ListPoliciesOptionsSortIDConst = "id"
	ListPoliciesOptionsSortLastModifiedAtConst = "last_modified_at"
	ListPoliciesOptionsSortLastModifiedByIDConst = "last_modified_by_id"
	ListPoliciesOptionsSortStateConst = "state"
	ListPoliciesOptionsSortTypeConst = "type"
)

// Constants associated with the ListPoliciesOptions.Format property.
// Include additional data per policy returned
// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
// it has done so
// * `display` - returns the list of all actions included in each of the policy roles.
const (
	ListPoliciesOptionsFormatDisplayConst = "display"
	ListPoliciesOptionsFormatIncludeLastPermitConst = "include_last_permit"
)

// Constants associated with the ListPoliciesOptions.State property.
// The state of the policy.
// * `active` - returns active policies
// * `deleted` - returns non-active policies.
const (
	ListPoliciesOptionsStateActiveConst = "active"
	ListPoliciesOptionsStateDeletedConst = "deleted"
)

// NewListPoliciesOptions : Instantiate ListPoliciesOptions
func (*IamPolicyManagementV1) NewListPoliciesOptions(accountID string) *ListPoliciesOptions {
	return &ListPoliciesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListPoliciesOptions) SetAccountID(accountID string) *ListPoliciesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListPoliciesOptions) SetAcceptLanguage(acceptLanguage string) *ListPoliciesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetIamID : Allow user to set IamID
func (_options *ListPoliciesOptions) SetIamID(iamID string) *ListPoliciesOptions {
	_options.IamID = core.StringPtr(iamID)
	return _options
}

// SetAccessGroupID : Allow user to set AccessGroupID
func (_options *ListPoliciesOptions) SetAccessGroupID(accessGroupID string) *ListPoliciesOptions {
	_options.AccessGroupID = core.StringPtr(accessGroupID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListPoliciesOptions) SetType(typeVar string) *ListPoliciesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetServiceType : Allow user to set ServiceType
func (_options *ListPoliciesOptions) SetServiceType(serviceType string) *ListPoliciesOptions {
	_options.ServiceType = core.StringPtr(serviceType)
	return _options
}

// SetTagName : Allow user to set TagName
func (_options *ListPoliciesOptions) SetTagName(tagName string) *ListPoliciesOptions {
	_options.TagName = core.StringPtr(tagName)
	return _options
}

// SetTagValue : Allow user to set TagValue
func (_options *ListPoliciesOptions) SetTagValue(tagValue string) *ListPoliciesOptions {
	_options.TagValue = core.StringPtr(tagValue)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListPoliciesOptions) SetSort(sort string) *ListPoliciesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *ListPoliciesOptions) SetFormat(format string) *ListPoliciesOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetState : Allow user to set State
func (_options *ListPoliciesOptions) SetState(state string) *ListPoliciesOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPoliciesOptions) SetHeaders(param map[string]string) *ListPoliciesOptions {
	options.Headers = param
	return options
}

// ListRolesOptions : The ListRoles options.
type ListRolesOptions struct {
	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Optional account GUID in which the roles belong to.
	AccountID *string `json:"account_id,omitempty"`

	// Optional name of IAM enabled service.
	ServiceName *string `json:"service_name,omitempty"`

	// Optional name of source IAM enabled service.
	SourceServiceName *string `json:"source_service_name,omitempty"`

	// Optional Policy Type.
	PolicyType *string `json:"policy_type,omitempty"`

	// Optional id of service group.
	ServiceGroupID *string `json:"service_group_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRolesOptions : Instantiate ListRolesOptions
func (*IamPolicyManagementV1) NewListRolesOptions() *ListRolesOptions {
	return &ListRolesOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListRolesOptions) SetAcceptLanguage(acceptLanguage string) *ListRolesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListRolesOptions) SetAccountID(accountID string) *ListRolesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListRolesOptions) SetServiceName(serviceName string) *ListRolesOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetSourceServiceName : Allow user to set SourceServiceName
func (_options *ListRolesOptions) SetSourceServiceName(sourceServiceName string) *ListRolesOptions {
	_options.SourceServiceName = core.StringPtr(sourceServiceName)
	return _options
}

// SetPolicyType : Allow user to set PolicyType
func (_options *ListRolesOptions) SetPolicyType(policyType string) *ListRolesOptions {
	_options.PolicyType = core.StringPtr(policyType)
	return _options
}

// SetServiceGroupID : Allow user to set ServiceGroupID
func (_options *ListRolesOptions) SetServiceGroupID(serviceGroupID string) *ListRolesOptions {
	_options.ServiceGroupID = core.StringPtr(serviceGroupID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRolesOptions) SetHeaders(param map[string]string) *ListRolesOptions {
	options.Headers = param
	return options
}

// ListV2PoliciesOptions : The ListV2Policies options.
type ListV2PoliciesOptions struct {
	// The account GUID in which the policies belong to.
	AccountID *string `json:"account_id" validate:"required"`

	// Language code for translations
	// * `default` - English
	// * `de` -  German (Standard)
	// * `en` - English
	// * `es` - Spanish (Spain)
	// * `fr` - French (Standard)
	// * `it` - Italian (Standard)
	// * `ja` - Japanese
	// * `ko` - Korean
	// * `pt-br` - Portuguese (Brazil)
	// * `zh-cn` - Chinese (Simplified, PRC)
	// * `zh-tw` - (Chinese, Taiwan).
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Optional IAM ID used to identify the subject.
	IamID *string `json:"iam_id,omitempty"`

	// Optional access group id.
	AccessGroupID *string `json:"access_group_id,omitempty"`

	// Optional type of policy.
	Type *string `json:"type,omitempty"`

	// Optional type of service.
	ServiceType *string `json:"service_type,omitempty"`

	// Optional name of service.
	ServiceName *string `json:"service_name,omitempty"`

	// Optional ID of service group.
	ServiceGroupID *string `json:"service_group_id,omitempty"`

	// Optional top level policy field to sort results. Ascending sort is default. Descending sort available by prepending
	// '-' to field, for example, '-last_modified_at'. Note that last permit information is only included when
	// 'format=include_last_permit', for example, "format=include_last_permit&sort=last_permit_at" Example fields that can
	// be sorted on:
	//   - 'id'
	//   - 'type'
	//   - 'href'
	//   - 'created_at'
	//   - 'created_by_id'
	//   - 'last_modified_at'
	//   - 'last_modified_by_id'
	//   - 'state'
	//   - 'last_permit_at'
	//   - 'last_permit_frequency'.
	Sort *string `json:"sort,omitempty"`

	// Include additional data per policy returned
	// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
	// it has done so
	// * `display` - returns the list of all actions included in each of the policy roles and translations for all relevant
	// fields.
	Format *string `json:"format,omitempty"`

	// The state of the policy.
	// * `active` - returns active policies
	// * `deleted` - returns non-active policies.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListV2PoliciesOptions.Type property.
// Optional type of policy.
const (
	ListV2PoliciesOptionsTypeAccessConst = "access"
	ListV2PoliciesOptionsTypeAuthorizationConst = "authorization"
)

// Constants associated with the ListV2PoliciesOptions.ServiceType property.
// Optional type of service.
const (
	ListV2PoliciesOptionsServiceTypePlatformServiceConst = "platform_service"
	ListV2PoliciesOptionsServiceTypeServiceConst = "service"
)

// Constants associated with the ListV2PoliciesOptions.Format property.
// Include additional data per policy returned
// * `include_last_permit` - returns details of when the policy last granted a permit decision and the number of times
// it has done so
// * `display` - returns the list of all actions included in each of the policy roles and translations for all relevant
// fields.
const (
	ListV2PoliciesOptionsFormatDisplayConst = "display"
	ListV2PoliciesOptionsFormatIncludeLastPermitConst = "include_last_permit"
)

// Constants associated with the ListV2PoliciesOptions.State property.
// The state of the policy.
// * `active` - returns active policies
// * `deleted` - returns non-active policies.
const (
	ListV2PoliciesOptionsStateActiveConst = "active"
	ListV2PoliciesOptionsStateDeletedConst = "deleted"
)

// NewListV2PoliciesOptions : Instantiate ListV2PoliciesOptions
func (*IamPolicyManagementV1) NewListV2PoliciesOptions(accountID string) *ListV2PoliciesOptions {
	return &ListV2PoliciesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListV2PoliciesOptions) SetAccountID(accountID string) *ListV2PoliciesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListV2PoliciesOptions) SetAcceptLanguage(acceptLanguage string) *ListV2PoliciesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetIamID : Allow user to set IamID
func (_options *ListV2PoliciesOptions) SetIamID(iamID string) *ListV2PoliciesOptions {
	_options.IamID = core.StringPtr(iamID)
	return _options
}

// SetAccessGroupID : Allow user to set AccessGroupID
func (_options *ListV2PoliciesOptions) SetAccessGroupID(accessGroupID string) *ListV2PoliciesOptions {
	_options.AccessGroupID = core.StringPtr(accessGroupID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListV2PoliciesOptions) SetType(typeVar string) *ListV2PoliciesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetServiceType : Allow user to set ServiceType
func (_options *ListV2PoliciesOptions) SetServiceType(serviceType string) *ListV2PoliciesOptions {
	_options.ServiceType = core.StringPtr(serviceType)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListV2PoliciesOptions) SetServiceName(serviceName string) *ListV2PoliciesOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetServiceGroupID : Allow user to set ServiceGroupID
func (_options *ListV2PoliciesOptions) SetServiceGroupID(serviceGroupID string) *ListV2PoliciesOptions {
	_options.ServiceGroupID = core.StringPtr(serviceGroupID)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListV2PoliciesOptions) SetSort(sort string) *ListV2PoliciesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *ListV2PoliciesOptions) SetFormat(format string) *ListV2PoliciesOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetState : Allow user to set State
func (_options *ListV2PoliciesOptions) SetState(state string) *ListV2PoliciesOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListV2PoliciesOptions) SetHeaders(param map[string]string) *ListV2PoliciesOptions {
	options.Headers = param
	return options
}

// PolicyRole : A role associated with a policy.
type PolicyRole struct {
	// The role Cloud Resource Name (CRN) granted by the policy. Example CRN: 'crn:v1:bluemix:public:iam::::role:Editor'.
	RoleID *string `json:"role_id" validate:"required"`

	// The display name of the role.
	DisplayName *string `json:"display_name,omitempty"`

	// The description of the role.
	Description *string `json:"description,omitempty"`
}

// NewPolicyRole : Instantiate PolicyRole (Generic Model Constructor)
func (*IamPolicyManagementV1) NewPolicyRole(roleID string) (_model *PolicyRole, err error) {
	_model = &PolicyRole{
		RoleID: core.StringPtr(roleID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalPolicyRole unmarshals an instance of PolicyRole from the specified map of raw messages.
func UnmarshalPolicyRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyRole)
	err = core.UnmarshalPrimitive(m, "role_id", &obj.RoleID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
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

// ReplacePolicyOptions : The ReplacePolicy options.
type ReplacePolicyOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// The revision number for updating a policy and must match the ETag value of the existing policy. The Etag can be
	// retrieved using the GET /v1/policies/{policy_id} API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// The subjects associated with a policy.
	Subjects []PolicySubject `json:"subjects" validate:"required"`

	// A set of role cloud resource names (CRNs) granted by the policy.
	Roles []PolicyRole `json:"roles" validate:"required"`

	// The resources associated with a policy.
	Resources []PolicyResource `json:"resources" validate:"required"`

	// Customer-defined description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplacePolicyOptions : Instantiate ReplacePolicyOptions
func (*IamPolicyManagementV1) NewReplacePolicyOptions(policyID string, ifMatch string, typeVar string, subjects []PolicySubject, roles []PolicyRole, resources []PolicyResource) *ReplacePolicyOptions {
	return &ReplacePolicyOptions{
		PolicyID: core.StringPtr(policyID),
		IfMatch: core.StringPtr(ifMatch),
		Type: core.StringPtr(typeVar),
		Subjects: subjects,
		Roles: roles,
		Resources: resources,
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *ReplacePolicyOptions) SetPolicyID(policyID string) *ReplacePolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplacePolicyOptions) SetIfMatch(ifMatch string) *ReplacePolicyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplacePolicyOptions) SetType(typeVar string) *ReplacePolicyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSubjects : Allow user to set Subjects
func (_options *ReplacePolicyOptions) SetSubjects(subjects []PolicySubject) *ReplacePolicyOptions {
	_options.Subjects = subjects
	return _options
}

// SetRoles : Allow user to set Roles
func (_options *ReplacePolicyOptions) SetRoles(roles []PolicyRole) *ReplacePolicyOptions {
	_options.Roles = roles
	return _options
}

// SetResources : Allow user to set Resources
func (_options *ReplacePolicyOptions) SetResources(resources []PolicyResource) *ReplacePolicyOptions {
	_options.Resources = resources
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplacePolicyOptions) SetDescription(description string) *ReplacePolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplacePolicyOptions) SetHeaders(param map[string]string) *ReplacePolicyOptions {
	options.Headers = param
	return options
}

// ReplaceRoleOptions : The ReplaceRole options.
type ReplaceRoleOptions struct {
	// The role ID.
	RoleID *string `json:"role_id" validate:"required,ne="`

	// The revision number for updating a role and must match the ETag value of the existing role. The Etag can be
	// retrieved using the GET /v2/roles/{role_id} API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The display name of the role that is shown in the console.
	DisplayName *string `json:"display_name" validate:"required"`

	// The actions of the role. For more information, see [IAM roles and
	// actions](https://cloud.ibm.com/docs/account?topic=account-iam-service-roles-actions).
	Actions []string `json:"actions" validate:"required"`

	// The description of the role.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceRoleOptions : Instantiate ReplaceRoleOptions
func (*IamPolicyManagementV1) NewReplaceRoleOptions(roleID string, ifMatch string, displayName string, actions []string) *ReplaceRoleOptions {
	return &ReplaceRoleOptions{
		RoleID: core.StringPtr(roleID),
		IfMatch: core.StringPtr(ifMatch),
		DisplayName: core.StringPtr(displayName),
		Actions: actions,
	}
}

// SetRoleID : Allow user to set RoleID
func (_options *ReplaceRoleOptions) SetRoleID(roleID string) *ReplaceRoleOptions {
	_options.RoleID = core.StringPtr(roleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceRoleOptions) SetIfMatch(ifMatch string) *ReplaceRoleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetDisplayName : Allow user to set DisplayName
func (_options *ReplaceRoleOptions) SetDisplayName(displayName string) *ReplaceRoleOptions {
	_options.DisplayName = core.StringPtr(displayName)
	return _options
}

// SetActions : Allow user to set Actions
func (_options *ReplaceRoleOptions) SetActions(actions []string) *ReplaceRoleOptions {
	_options.Actions = actions
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceRoleOptions) SetDescription(description string) *ReplaceRoleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceRoleOptions) SetHeaders(param map[string]string) *ReplaceRoleOptions {
	options.Headers = param
	return options
}

// ReplaceV2PolicyOptions : The ReplaceV2Policy options.
type ReplaceV2PolicyOptions struct {
	// The policy ID.
	ID *string `json:"id" validate:"required,ne="`

	// The revision number for updating a policy and must match the ETag value of the existing policy. The Etag can be
	// retrieved using the GET /v2/policies/{id} API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Specifies the type of access granted by the policy.
	Control *Control `json:"control" validate:"required"`

	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Allows the customer to use their own words to record the purpose/context related to a policy.
	Description *string `json:"description,omitempty"`

	// The subject attributes for whom the policy grants access.
	Subject *V2PolicySubject `json:"subject,omitempty"`

	// The resource attributes to which the policy grants access.
	Resource *V2PolicyResource `json:"resource,omitempty"`

	// Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or
	// 'time-based-conditions:weekly:custom-hours'.
	Pattern *string `json:"pattern,omitempty"`

	// Additional access conditions associated with the policy.
	Rule V2PolicyRuleIntf `json:"rule,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceV2PolicyOptions.Type property.
// The policy type; either 'access' or 'authorization'.
const (
	ReplaceV2PolicyOptionsTypeAccessConst = "access"
	ReplaceV2PolicyOptionsTypeAuthorizationConst = "authorization"
)

// NewReplaceV2PolicyOptions : Instantiate ReplaceV2PolicyOptions
func (*IamPolicyManagementV1) NewReplaceV2PolicyOptions(id string, ifMatch string, control *Control, typeVar string) *ReplaceV2PolicyOptions {
	return &ReplaceV2PolicyOptions{
		ID: core.StringPtr(id),
		IfMatch: core.StringPtr(ifMatch),
		Control: control,
		Type: core.StringPtr(typeVar),
	}
}

// SetID : Allow user to set ID
func (_options *ReplaceV2PolicyOptions) SetID(id string) *ReplaceV2PolicyOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceV2PolicyOptions) SetIfMatch(ifMatch string) *ReplaceV2PolicyOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetControl : Allow user to set Control
func (_options *ReplaceV2PolicyOptions) SetControl(control *Control) *ReplaceV2PolicyOptions {
	_options.Control = control
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceV2PolicyOptions) SetType(typeVar string) *ReplaceV2PolicyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceV2PolicyOptions) SetDescription(description string) *ReplaceV2PolicyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetSubject : Allow user to set Subject
func (_options *ReplaceV2PolicyOptions) SetSubject(subject *V2PolicySubject) *ReplaceV2PolicyOptions {
	_options.Subject = subject
	return _options
}

// SetResource : Allow user to set Resource
func (_options *ReplaceV2PolicyOptions) SetResource(resource *V2PolicyResource) *ReplaceV2PolicyOptions {
	_options.Resource = resource
	return _options
}

// SetPattern : Allow user to set Pattern
func (_options *ReplaceV2PolicyOptions) SetPattern(pattern string) *ReplaceV2PolicyOptions {
	_options.Pattern = core.StringPtr(pattern)
	return _options
}

// SetRule : Allow user to set Rule
func (_options *ReplaceV2PolicyOptions) SetRule(rule V2PolicyRuleIntf) *ReplaceV2PolicyOptions {
	_options.Rule = rule
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceV2PolicyOptions) SetHeaders(param map[string]string) *ReplaceV2PolicyOptions {
	options.Headers = param
	return options
}

// RoleAction : An action that can be performed by the policy subject when assigned role.
type RoleAction struct {
	// Unique identifier for action with structure service.resource.action e.g., cbr.rule.read.
	ID *string `json:"id" validate:"required"`

	// Service defined display name for action.
	DisplayName *string `json:"display_name" validate:"required"`

	// Service defined description for action.
	Description *string `json:"description" validate:"required"`
}

// UnmarshalRoleAction unmarshals an instance of RoleAction from the specified map of raw messages.
func UnmarshalRoleAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleAction)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
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

// Roles : A role associated with a policy.
type Roles struct {
	// The role Cloud Resource Name (CRN) granted by the policy. Example CRN: 'crn:v1:bluemix:public:iam::::role:Editor'.
	RoleID *string `json:"role_id" validate:"required"`
}

// NewRoles : Instantiate Roles (Generic Model Constructor)
func (*IamPolicyManagementV1) NewRoles(roleID string) (_model *Roles, err error) {
	_model = &Roles{
		RoleID: core.StringPtr(roleID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRoles unmarshals an instance of Roles from the specified map of raw messages.
func UnmarshalRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Roles)
	err = core.UnmarshalPrimitive(m, "role_id", &obj.RoleID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleAttribute : Rule that specifies additional access granted (e.g., time-based condition).
type RuleAttribute struct {
	// The name of an attribute.
	Key *string `json:"key" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator" validate:"required"`

	// The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an
	// array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the RuleAttribute.Operator property.
// The operator of an attribute.
const (
	RuleAttributeOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	RuleAttributeOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	RuleAttributeOperatorDatetimelessthanConst = "dateTimeLessThan"
	RuleAttributeOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	RuleAttributeOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	RuleAttributeOperatorDayofweekequalsConst = "dayOfWeekEquals"
	RuleAttributeOperatorTimegreaterthanConst = "timeGreaterThan"
	RuleAttributeOperatorTimegreaterthanorequalsConst = "timeGreaterThanOrEquals"
	RuleAttributeOperatorTimelessthanConst = "timeLessThan"
	RuleAttributeOperatorTimelessthanorequalsConst = "timeLessThanOrEquals"
)

// NewRuleAttribute : Instantiate RuleAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewRuleAttribute(key string, operator string, value interface{}) (_model *RuleAttribute, err error) {
	_model = &RuleAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleAttribute unmarshals an instance of RuleAttribute from the specified map of raw messages.
func UnmarshalRuleAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// UpdatePolicyStateOptions : The UpdatePolicyState options.
type UpdatePolicyStateOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// The revision number for updating a policy and must match the ETag value of the existing policy. The Etag can be
	// retrieved using the GET /v1/policies/{policy_id} API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The policy state.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdatePolicyStateOptions.State property.
// The policy state.
const (
	UpdatePolicyStateOptionsStateActiveConst = "active"
	UpdatePolicyStateOptionsStateDeletedConst = "deleted"
)

// NewUpdatePolicyStateOptions : Instantiate UpdatePolicyStateOptions
func (*IamPolicyManagementV1) NewUpdatePolicyStateOptions(policyID string, ifMatch string) *UpdatePolicyStateOptions {
	return &UpdatePolicyStateOptions{
		PolicyID: core.StringPtr(policyID),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetPolicyID : Allow user to set PolicyID
func (_options *UpdatePolicyStateOptions) SetPolicyID(policyID string) *UpdatePolicyStateOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdatePolicyStateOptions) SetIfMatch(ifMatch string) *UpdatePolicyStateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetState : Allow user to set State
func (_options *UpdatePolicyStateOptions) SetState(state string) *UpdatePolicyStateOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePolicyStateOptions) SetHeaders(param map[string]string) *UpdatePolicyStateOptions {
	options.Headers = param
	return options
}

// V2Policy : The core set of properties associated with the policy.
type V2Policy struct {
	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Allows the customer to use their own words to record the purpose/context related to a policy.
	Description *string `json:"description,omitempty"`

	// The subject attributes for whom the policy grants access.
	Subject *V2PolicySubject `json:"subject,omitempty"`

	// The resource attributes to which the policy grants access.
	Resource *V2PolicyResource `json:"resource,omitempty"`

	// Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or
	// 'time-based-conditions:weekly:custom-hours'.
	Pattern *string `json:"pattern,omitempty"`

	// Additional access conditions associated with the policy.
	Rule V2PolicyRuleIntf `json:"rule,omitempty"`

	// The policy ID.
	ID *string `json:"id,omitempty"`

	// The href URL that links to the policies API by policy ID.
	Href *string `json:"href,omitempty"`

	Control ControlResponseIntf `json:"control" validate:"required"`

	// The UTC timestamp when the policy was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// The policy state, either 'deleted' or 'active'.
	State *string `json:"state" validate:"required"`

	// The optional last permit time of policy, when passing query parameter format=include_last_permit.
	LastPermitAt *string `json:"last_permit_at,omitempty"`

	// The optional count of times that policy has provided a permit, when passing query parameter
	// format=include_last_permit.
	LastPermitFrequency *int64 `json:"last_permit_frequency,omitempty"`
}

// Constants associated with the V2Policy.Type property.
// The policy type; either 'access' or 'authorization'.
const (
	V2PolicyTypeAccessConst = "access"
	V2PolicyTypeAuthorizationConst = "authorization"
)

// Constants associated with the V2Policy.State property.
// The policy state, either 'deleted' or 'active'.
const (
	V2PolicyStateActiveConst = "active"
	V2PolicyStateDeletedConst = "deleted"
)

// UnmarshalV2Policy unmarshals an instance of V2Policy from the specified map of raw messages.
func UnmarshalV2Policy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2Policy)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalV2PolicySubject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource", &obj.Resource, UnmarshalV2PolicyResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalV2PolicyRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control", &obj.Control, UnmarshalControlResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_at", &obj.LastPermitAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_frequency", &obj.LastPermitFrequency)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyCollection : A collection of policies.
type V2PolicyCollection struct {
	// List of policies.
	Policies []V2Policy `json:"policies,omitempty"`
}

// UnmarshalV2PolicyCollection unmarshals an instance of V2PolicyCollection from the specified map of raw messages.
func UnmarshalV2PolicyCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyCollection)
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalV2Policy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyResource : The resource attributes to which the policy grants access.
type V2PolicyResource struct {
	// List of resource attributes to which the policy grants access.
	Attributes []V2PolicyResourceAttribute `json:"attributes" validate:"required"`

	// Optional list of resource tags to which the policy grants access.
	Tags []V2PolicyResourceTag `json:"tags,omitempty"`
}

// NewV2PolicyResource : Instantiate V2PolicyResource (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyResource(attributes []V2PolicyResourceAttribute) (_model *V2PolicyResource, err error) {
	_model = &V2PolicyResource{
		Attributes: attributes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalV2PolicyResource unmarshals an instance of V2PolicyResource from the specified map of raw messages.
func UnmarshalV2PolicyResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResource)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalV2PolicyResourceAttribute)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalV2PolicyResourceTag)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyResourceAttribute : Resource attribute to which the policy grants access.
type V2PolicyResourceAttribute struct {
	// The name of a resource attribute.
	Key *string `json:"key" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator" validate:"required"`

	// The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an
	// array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the V2PolicyResourceAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicyResourceAttributeOperatorStringequalsConst = "stringEquals"
	V2PolicyResourceAttributeOperatorStringexistsConst = "stringExists"
	V2PolicyResourceAttributeOperatorStringmatchConst = "stringMatch"
)

// NewV2PolicyResourceAttribute : Instantiate V2PolicyResourceAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyResourceAttribute(key string, operator string, value interface{}) (_model *V2PolicyResourceAttribute, err error) {
	_model = &V2PolicyResourceAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalV2PolicyResourceAttribute unmarshals an instance of V2PolicyResourceAttribute from the specified map of raw messages.
func UnmarshalV2PolicyResourceAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResourceAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// V2PolicyResourceTag : A tag associated with a resource.
type V2PolicyResourceTag struct {
	// The name of an access management tag.
	Key *string `json:"key" validate:"required"`

	// The value of an access management tag.
	Value *string `json:"value" validate:"required"`

	// The operator of an access management tag.
	Operator *string `json:"operator" validate:"required"`
}

// Constants associated with the V2PolicyResourceTag.Operator property.
// The operator of an access management tag.
const (
	V2PolicyResourceTagOperatorStringequalsConst = "stringEquals"
	V2PolicyResourceTagOperatorStringmatchConst = "stringMatch"
)

// NewV2PolicyResourceTag : Instantiate V2PolicyResourceTag (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyResourceTag(key string, value string, operator string) (_model *V2PolicyResourceTag, err error) {
	_model = &V2PolicyResourceTag{
		Key: core.StringPtr(key),
		Value: core.StringPtr(value),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalV2PolicyResourceTag unmarshals an instance of V2PolicyResourceTag from the specified map of raw messages.
func UnmarshalV2PolicyResourceTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResourceTag)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyRule : Additional access conditions associated with the policy.
// Models which "extend" this model:
// - V2PolicyRuleRuleAttribute
// - V2PolicyRuleRuleWithConditions
type V2PolicyRule struct {
	// The name of an attribute.
	Key *string `json:"key,omitempty"`

	// The operator of an attribute.
	Operator *string `json:"operator,omitempty"`

	// The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an
	// array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value,omitempty"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []RuleAttribute `json:"conditions,omitempty"`
}

// Constants associated with the V2PolicyRule.Operator property.
// The operator of an attribute.
const (
	V2PolicyRuleOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	V2PolicyRuleOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	V2PolicyRuleOperatorDatetimelessthanConst = "dateTimeLessThan"
	V2PolicyRuleOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	V2PolicyRuleOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	V2PolicyRuleOperatorDayofweekequalsConst = "dayOfWeekEquals"
	V2PolicyRuleOperatorTimegreaterthanConst = "timeGreaterThan"
	V2PolicyRuleOperatorTimegreaterthanorequalsConst = "timeGreaterThanOrEquals"
	V2PolicyRuleOperatorTimelessthanConst = "timeLessThan"
	V2PolicyRuleOperatorTimelessthanorequalsConst = "timeLessThanOrEquals"
)
func (*V2PolicyRule) isaV2PolicyRule() bool {
	return true
}

type V2PolicyRuleIntf interface {
	isaV2PolicyRule() bool
}

// UnmarshalV2PolicyRule unmarshals an instance of V2PolicyRule from the specified map of raw messages.
func UnmarshalV2PolicyRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyRule)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalRuleAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicySubject : The subject attributes for whom the policy grants access.
type V2PolicySubject struct {
	// List of subject attributes associated with policy/.
	Attributes []V2PolicySubjectAttribute `json:"attributes" validate:"required"`
}

// NewV2PolicySubject : Instantiate V2PolicySubject (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicySubject(attributes []V2PolicySubjectAttribute) (_model *V2PolicySubject, err error) {
	_model = &V2PolicySubject{
		Attributes: attributes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalV2PolicySubject unmarshals an instance of V2PolicySubject from the specified map of raw messages.
func UnmarshalV2PolicySubject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicySubject)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalV2PolicySubjectAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicySubjectAttribute : Subject attribute for whom the policy grants access.
type V2PolicySubjectAttribute struct {
	// The name of a subject attribute, e.g., iam_id, access_group_id.
	Key *string `json:"key" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator" validate:"required"`

	// The value of the ID of the subject, e.g., service ID, access group ID, IAM ID.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the V2PolicySubjectAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicySubjectAttributeOperatorStringequalsConst = "stringEquals"
)

// NewV2PolicySubjectAttribute : Instantiate V2PolicySubjectAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicySubjectAttribute(key string, operator string, value string) (_model *V2PolicySubjectAttribute, err error) {
	_model = &V2PolicySubjectAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalV2PolicySubjectAttribute unmarshals an instance of V2PolicySubjectAttribute from the specified map of raw messages.
func UnmarshalV2PolicySubjectAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicySubjectAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// CustomRole : An additional set of properties associated with a role.
type CustomRole struct {
	// The role ID. Composed of hexadecimal characters.
	ID *string `json:"id,omitempty"`

	// The display name of the role that is shown in the console.
	DisplayName *string `json:"display_name" validate:"required"`

	// The description of the role.
	Description *string `json:"description,omitempty"`

	// The actions of the role. For more information, see [IAM roles and
	// actions](https://cloud.ibm.com/docs/account?topic=account-iam-service-roles-actions).
	Actions []string `json:"actions" validate:"required"`

	// The role Cloud Resource Name (CRN). Example CRN:
	// 'crn:v1:ibmcloud:public:iam-access-management::a/exampleAccountId::customRole:ExampleRoleName'.
	CRN *string `json:"crn,omitempty"`

	// The name of the role that is used in the CRN. Can only be alphanumeric and has to be capitalized.
	Name *string `json:"name" validate:"required"`

	// The account GUID.
	AccountID *string `json:"account_id" validate:"required"`

	// The service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The UTC timestamp when the role was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the role.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the role was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// The href link back to the role.
	Href *string `json:"href,omitempty"`
}

// UnmarshalCustomRole unmarshals an instance of CustomRole from the specified map of raw messages.
func UnmarshalCustomRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomRole)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "actions", &obj.Actions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
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

// Policy : The core set of properties associated with a policy.
type Policy struct {
	// The policy ID.
	ID *string `json:"id,omitempty"`

	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Customer-defined description.
	Description *string `json:"description,omitempty"`

	// The subjects associated with a policy.
	Subjects []PolicySubject `json:"subjects" validate:"required"`

	// A set of role cloud resource names (CRNs) granted by the policy.
	Roles []PolicyRole `json:"roles" validate:"required"`

	// The resources associated with a policy.
	Resources []PolicyResource `json:"resources" validate:"required"`

	// The href link back to the policy.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// The policy state.
	State *string `json:"state,omitempty"`
}

// Constants associated with the Policy.State property.
// The policy state.
const (
	PolicyStateActiveConst = "active"
	PolicyStateDeletedConst = "deleted"
)

// UnmarshalPolicy unmarshals an instance of Policy from the specified map of raw messages.
func UnmarshalPolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Policy)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalModel(m, "subjects", &obj.Subjects, UnmarshalPolicySubject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalPolicyRole)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyResource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
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

// PolicyList : A collection of policies.
type PolicyList struct {
	// List of policies.
	Policies []Policy `json:"policies,omitempty"`
}

// UnmarshalPolicyList unmarshals an instance of PolicyList from the specified map of raw messages.
func UnmarshalPolicyList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyList)
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyResource : The attributes of the resource. Note that only one resource is allowed in a policy.
type PolicyResource struct {
	// List of resource attributes.
	Attributes []ResourceAttribute `json:"attributes,omitempty"`

	// List of access management tags.
	Tags []ResourceTag `json:"tags,omitempty"`
}

// UnmarshalPolicyResource unmarshals an instance of PolicyResource from the specified map of raw messages.
func UnmarshalPolicyResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyResource)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalResourceAttribute)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalResourceTag)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicySubject : The subject attribute values that must match in order for this policy to apply in a permission decision.
type PolicySubject struct {
	// List of subject attributes.
	Attributes []SubjectAttribute `json:"attributes,omitempty"`
}

// UnmarshalPolicySubject unmarshals an instance of PolicySubject from the specified map of raw messages.
func UnmarshalPolicySubject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicySubject)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalSubjectAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceAttribute : An attribute associated with a resource.
type ResourceAttribute struct {
	// The name of an attribute.
	Name *string `json:"name" validate:"required"`

	// The value of an attribute.
	Value *string `json:"value" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator,omitempty"`
}

// NewResourceAttribute : Instantiate ResourceAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewResourceAttribute(name string, value string) (_model *ResourceAttribute, err error) {
	_model = &ResourceAttribute{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalResourceAttribute unmarshals an instance of ResourceAttribute from the specified map of raw messages.
func UnmarshalResourceAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceTag : A tag associated with a resource.
type ResourceTag struct {
	// The name of an access management tag.
	Name *string `json:"name" validate:"required"`

	// The value of an access management tag.
	Value *string `json:"value" validate:"required"`

	// The operator of an access management tag.
	Operator *string `json:"operator,omitempty"`
}

// NewResourceTag : Instantiate ResourceTag (Generic Model Constructor)
func (*IamPolicyManagementV1) NewResourceTag(name string, value string) (_model *ResourceTag, err error) {
	_model = &ResourceTag{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalResourceTag unmarshals an instance of ResourceTag from the specified map of raw messages.
func UnmarshalResourceTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceTag)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Role : A role resource.
type Role struct {
	// The display name of the role that is shown in the console.
	DisplayName *string `json:"display_name" validate:"required"`

	// The description of the role.
	Description *string `json:"description,omitempty"`

	// The actions of the role. For more information, see [IAM roles and
	// actions](https://cloud.ibm.com/docs/account?topic=account-iam-service-roles-actions).
	Actions []string `json:"actions" validate:"required"`

	// The role Cloud Resource Name (CRN). Example CRN:
	// 'crn:v1:ibmcloud:public:iam-access-management::a/exampleAccountId::customRole:ExampleRoleName'.
	CRN *string `json:"crn,omitempty"`
}

// NewRole : Instantiate Role (Generic Model Constructor)
func (*IamPolicyManagementV1) NewRole(displayName string, actions []string) (_model *Role, err error) {
	_model = &Role{
		DisplayName: core.StringPtr(displayName),
		Actions: actions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRole unmarshals an instance of Role from the specified map of raw messages.
func UnmarshalRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Role)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "actions", &obj.Actions)
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

// RoleList : A collection of roles returned by the 'list roles' operation.
type RoleList struct {
	// List of custom roles.
	CustomRoles []CustomRole `json:"custom_roles,omitempty"`

	// List of service roles.
	ServiceRoles []Role `json:"service_roles,omitempty"`

	// List of system roles.
	SystemRoles []Role `json:"system_roles,omitempty"`
}

// UnmarshalRoleList unmarshals an instance of RoleList from the specified map of raw messages.
func UnmarshalRoleList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleList)
	err = core.UnmarshalModel(m, "custom_roles", &obj.CustomRoles, UnmarshalCustomRole)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "service_roles", &obj.ServiceRoles, UnmarshalRole)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "system_roles", &obj.SystemRoles, UnmarshalRole)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubjectAttribute : An attribute associated with a subject.
type SubjectAttribute struct {
	// The name of an attribute.
	Name *string `json:"name" validate:"required"`

	// The value of an attribute.
	Value *string `json:"value" validate:"required"`
}

// NewSubjectAttribute : Instantiate SubjectAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewSubjectAttribute(name string, value string) (_model *SubjectAttribute, err error) {
	_model = &SubjectAttribute{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSubjectAttribute unmarshals an instance of SubjectAttribute from the specified map of raw messages.
func UnmarshalSubjectAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubjectAttribute)
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

// ControlResponseControl : Specifies the type of access granted by the policy.
// This model "extends" ControlResponse
type ControlResponseControl struct {
	// Permission granted by the policy.
	Grant *Grant `json:"grant" validate:"required"`
}

func (*ControlResponseControl) isaControlResponse() bool {
	return true
}

// UnmarshalControlResponseControl unmarshals an instance of ControlResponseControl from the specified map of raw messages.
func UnmarshalControlResponseControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlResponseControl)
	err = core.UnmarshalModel(m, "grant", &obj.Grant, UnmarshalGrant)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlResponseControlWithEnrichedRoles : Specifies the type of access granted by the policy with additional role information.
// This model "extends" ControlResponse
type ControlResponseControlWithEnrichedRoles struct {
	// Permission granted by the policy with translated roles and additional role information.
	Grant *GrantWithEnrichedRoles `json:"grant" validate:"required"`
}

func (*ControlResponseControlWithEnrichedRoles) isaControlResponse() bool {
	return true
}

// UnmarshalControlResponseControlWithEnrichedRoles unmarshals an instance of ControlResponseControlWithEnrichedRoles from the specified map of raw messages.
func UnmarshalControlResponseControlWithEnrichedRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlResponseControlWithEnrichedRoles)
	err = core.UnmarshalModel(m, "grant", &obj.Grant, UnmarshalGrantWithEnrichedRoles)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyRuleRuleAttribute : Rule that specifies additional access granted (e.g., time-based condition).
// This model "extends" V2PolicyRule
type V2PolicyRuleRuleAttribute struct {
	// The name of an attribute.
	Key *string `json:"key" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator" validate:"required"`

	// The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an
	// array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the V2PolicyRuleRuleAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicyRuleRuleAttributeOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	V2PolicyRuleRuleAttributeOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDatetimelessthanConst = "dateTimeLessThan"
	V2PolicyRuleRuleAttributeOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	V2PolicyRuleRuleAttributeOperatorDayofweekequalsConst = "dayOfWeekEquals"
	V2PolicyRuleRuleAttributeOperatorTimegreaterthanConst = "timeGreaterThan"
	V2PolicyRuleRuleAttributeOperatorTimegreaterthanorequalsConst = "timeGreaterThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorTimelessthanConst = "timeLessThan"
	V2PolicyRuleRuleAttributeOperatorTimelessthanorequalsConst = "timeLessThanOrEquals"
)

// NewV2PolicyRuleRuleAttribute : Instantiate V2PolicyRuleRuleAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyRuleRuleAttribute(key string, operator string, value interface{}) (_model *V2PolicyRuleRuleAttribute, err error) {
	_model = &V2PolicyRuleRuleAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*V2PolicyRuleRuleAttribute) isaV2PolicyRule() bool {
	return true
}

// UnmarshalV2PolicyRuleRuleAttribute unmarshals an instance of V2PolicyRuleRuleAttribute from the specified map of raw messages.
func UnmarshalV2PolicyRuleRuleAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyRuleRuleAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// V2PolicyRuleRuleWithConditions : Rule that specifies additional access granted (e.g., time-based condition) accross multiple conditions.
// This model "extends" V2PolicyRule
type V2PolicyRuleRuleWithConditions struct {
	// Operator to evaluate conditions.
	Operator *string `json:"operator" validate:"required"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []RuleAttribute `json:"conditions" validate:"required"`
}

// Constants associated with the V2PolicyRuleRuleWithConditions.Operator property.
// Operator to evaluate conditions.
const (
	V2PolicyRuleRuleWithConditionsOperatorAndConst = "and"
	V2PolicyRuleRuleWithConditionsOperatorOrConst = "or"
)

// NewV2PolicyRuleRuleWithConditions : Instantiate V2PolicyRuleRuleWithConditions (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyRuleRuleWithConditions(operator string, conditions []RuleAttribute) (_model *V2PolicyRuleRuleWithConditions, err error) {
	_model = &V2PolicyRuleRuleWithConditions{
		Operator: core.StringPtr(operator),
		Conditions: conditions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*V2PolicyRuleRuleWithConditions) isaV2PolicyRule() bool {
	return true
}

// UnmarshalV2PolicyRuleRuleWithConditions unmarshals an instance of V2PolicyRuleRuleWithConditions from the specified map of raw messages.
func UnmarshalV2PolicyRuleRuleWithConditions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyRuleRuleWithConditions)
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalRuleAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
