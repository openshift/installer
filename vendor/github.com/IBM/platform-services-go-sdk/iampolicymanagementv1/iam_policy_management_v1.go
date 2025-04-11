/**
 * (C) Copyright IBM Corp. 2025.
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
 * IBM OpenAPI SDK Code Generator Version: 3.100.0-2ad7a784-20250212-162551
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
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	iamPolicyManagement, err = NewIamPolicyManagementV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = iamPolicyManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = iamPolicyManagement.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
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

	service = &IamPolicyManagementV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
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
	err := iamPolicyManagement.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
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
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicies(listPoliciesOptions *ListPoliciesOptions) (result *PolicyCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListPoliciesWithContext(context.Background(), listPoliciesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPoliciesWithContext is an alternate form of the ListPolicies method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListPoliciesWithContext(ctx context.Context, listPoliciesOptions *ListPoliciesOptions) (result *PolicyCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPoliciesOptions, "listPoliciesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPoliciesOptions, "listPoliciesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
	if listPoliciesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPoliciesOptions.Limit))
	}
	if listPoliciesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listPoliciesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_policies", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
//   serviceName, serviceInstance, region, resourceType, resource, accountId, resourceGroupId Assign roles that are
// supported by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`accountId`** attributes. The policy subject must include
// either **`serviceName`** or **`resourceGroupId`** (or both) attributes.
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
	result, response, err = iamPolicyManagement.CreatePolicyWithContext(context.Background(), createPolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePolicyWithContext is an alternate form of the CreatePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyWithContext(ctx context.Context, createPolicyOptions *CreatePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPolicyOptions, "createPolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPolicyOptions, "createPolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policies`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
//   serviceName, serviceInstance, region, resourceType, resource, accountId, resourceGroupId Assign roles that are
// supported by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`accountId`** attributes. The policy subject must include
// either **`serviceName`** or **`resourceGroupId`** (or both) attributes.
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
	result, response, err = iamPolicyManagement.ReplacePolicyWithContext(context.Background(), replacePolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplacePolicyWithContext is an alternate form of the ReplacePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplacePolicyWithContext(ctx context.Context, replacePolicyOptions *ReplacePolicyOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replacePolicyOptions, "replacePolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replacePolicyOptions, "replacePolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetPolicy : Retrieve a policy by ID
// Retrieve a policy by providing a policy ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicy(getPolicyOptions *GetPolicyOptions) (result *PolicyTemplateMetaData, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetPolicyWithContext(context.Background(), getPolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPolicyWithContext is an alternate form of the GetPolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyWithContext(ctx context.Context, getPolicyOptions *GetPolicyOptions) (result *PolicyTemplateMetaData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyOptions, "getPolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPolicyOptions, "getPolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateMetaData)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
	response, err = iamPolicyManagement.DeletePolicyWithContext(context.Background(), deletePolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePolicyWithContext is an alternate form of the DeletePolicy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyWithContext(ctx context.Context, deletePolicyOptions *DeletePolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePolicyOptions, "deletePolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePolicyOptions, "deletePolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdatePolicyState : Restore a deleted policy by ID
// Restore a policy that has recently been deleted. A policy administrator might want to restore a deleted policy. To
// restore a policy, use **`"state": "active"`** in the body.
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyState(updatePolicyStateOptions *UpdatePolicyStateOptions) (result *Policy, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.UpdatePolicyStateWithContext(context.Background(), updatePolicyStateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdatePolicyStateWithContext is an alternate form of the UpdatePolicyState method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyStateWithContext(ctx context.Context, updatePolicyStateOptions *UpdatePolicyStateOptions) (result *Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePolicyStateOptions, "updatePolicyStateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updatePolicyStateOptions, "updatePolicyStateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_policy_state", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
func (iamPolicyManagement *IamPolicyManagementV1) ListRoles(listRolesOptions *ListRolesOptions) (result *RoleCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListRolesWithContext(context.Background(), listRolesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListRolesWithContext is an alternate form of the ListRoles method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListRolesWithContext(ctx context.Context, listRolesOptions *ListRolesOptions) (result *RoleCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRolesOptions, "listRolesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_roles", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoleCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
	result, response, err = iamPolicyManagement.CreateRoleWithContext(context.Background(), createRoleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateRoleWithContext is an alternate form of the CreateRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreateRoleWithContext(ctx context.Context, createRoleOptions *CreateRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRoleOptions, "createRoleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createRoleOptions, "createRoleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/roles`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_role", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
	result, response, err = iamPolicyManagement.ReplaceRoleWithContext(context.Background(), replaceRoleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceRoleWithContext is an alternate form of the ReplaceRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceRoleWithContext(ctx context.Context, replaceRoleOptions *ReplaceRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRoleOptions, "replaceRoleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceRoleOptions, "replaceRoleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_role", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetRole : Retrieve a role by ID
// Retrieve a role by providing a role ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetRole(getRoleOptions *GetRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetRoleWithContext(context.Background(), getRoleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetRoleWithContext is an alternate form of the GetRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetRoleWithContext(ctx context.Context, getRoleOptions *GetRoleOptions) (result *CustomRole, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRoleOptions, "getRoleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getRoleOptions, "getRoleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_role", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomRole)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteRole : Delete a role by ID
// Delete a role by providing a role ID.
func (iamPolicyManagement *IamPolicyManagementV1) DeleteRole(deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	response, err = iamPolicyManagement.DeleteRoleWithContext(context.Background(), deleteRoleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteRoleWithContext is an alternate form of the DeleteRole method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeleteRoleWithContext(ctx context.Context, deleteRoleOptions *DeleteRoleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRoleOptions, "deleteRoleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteRoleOptions, "deleteRoleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_role", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListV2Policies : Get policies by attributes
// Get policies and filter by attributes. While managing policies, you might want to retrieve policies in the account
// and filter by attribute values. This can be done through query parameters. The following attributes are supported:
// account_id, iam_id, access_group_id, type, service_type, sort, format and state. account_id is a required query
// parameter. Only policies that have the specified attributes and that the caller has read access to are returned. If
// the caller does not have read access to any policies an empty array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListV2Policies(listV2PoliciesOptions *ListV2PoliciesOptions) (result *V2PolicyCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListV2PoliciesWithContext(context.Background(), listV2PoliciesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListV2PoliciesWithContext is an alternate form of the ListV2Policies method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListV2PoliciesWithContext(ctx context.Context, listV2PoliciesOptions *ListV2PoliciesOptions) (result *V2PolicyCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listV2PoliciesOptions, "listV2PoliciesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listV2PoliciesOptions, "listV2PoliciesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
	if listV2PoliciesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listV2PoliciesOptions.Limit))
	}
	if listV2PoliciesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listV2PoliciesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_v2_policies", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2PolicyCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
// combination **`operator`**. The possible combination operators are **`and`** and **`or`**.
//
// Currently, we support two types of patterns:
//
// 1. `time-based`: Used to specify a time-based restriction
//
// Combine conditions to specify a time-based restriction (e.g., access only during business hours, during the
// Monday-Friday work week). For example, a policy can grant access Monday-Friday, 9:00am-5:00pm using the following
// rule:
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
//   'dateLessThan', 'dateLessThanOrEquals', 'dateGreaterThan', 'dateGreaterThanOrEquals',
//   'dateTimeLessThan', 'dateTimeLessThanOrEquals', 'dateTimeGreaterThan', 'dateTimeGreaterThanOrEquals',
//   'dayOfWeekEquals', 'dayOfWeekAnyOf'
// ``` The pattern field that matches the rule is required when rule is provided. For the business hour rule example
// above, the **`pattern`** is **`"time-based-conditions:weekly"`**. For more information, see [Time-based conditions
// operators](/docs/account?topic=account-iam-condition-properties&interface=ui#policy-condition-properties) and
// [Limiting access with time-based conditions](/docs/account?topic=account-iam-time-based&interface=ui). If the subject
// is a locked service-id, the request will fail.
//
// 2. `attribute-based`: Used to specify a combination of OR/AND based conditions applied on resource attributes.
//
// Combine conditions to specify an attribute-based condition using AND/OR-based operators.
//
// For example, a policy can grant access based on multiple conditions applied on the resource attributes below:
// ```json
//   "pattern": "attribute-based-condition:resource:literal-and-wildcard"
//   "rule": {
//       "operator": "or",
//       "conditions": [
//         {
//           "operator": "and",
//           "conditions": [
//             {
//               "key": "{{resource.attributes.prefix}}",
//               "operator": "stringEquals",
//               "value": "home/test"
//             },
//             {
//               "key": "{{environment.attributes.delimiter}}",
//               "operator": "stringEquals",
//               "value": "/"
//             }
//           ]
//         },
//         {
//           "key": "{{resource.attributes.path}}",
//           "operator": "stringMatch",
//           "value": "home/David/_*"
//         }
//       ]
//   }
// ```
//
// In addition to satisfying the `resources` section, the policy grants permission only if either the `path` begins with
// `home/David/` **OR**  the `prefix` is `home/test` and the `delimiter` is `/`. This mechanism helps you consolidate
// multiple policies in to a single policy,  making policies easier to administer and stay within the policy limit for
// an account. View the list of operators that can be used in the condition
// [here](/docs/account?topic=account-wildcard#string-comparisons).
//
// ### Authorization
//
// Authorization policies are supported by services on a case by case basis. Refer to service documentation to verify
// their support of authorization policies. To create an authorization policy, use **`"type": "authorization"`** in the
// body. The subject attributes must match the supported authorization subjects of the resource. Multiple subject
// attributes might be provided. The following attributes are supported:
//   serviceName, serviceInstance, region, resourceType, resource, accountId, resourceGroupId Assign roles that are
// supported by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`accountId`** attributes. The policy subject must include
// either **`serviceName`** or **`resourceGroupId`** (or both) attributes.
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
	result, response, err = iamPolicyManagement.CreateV2PolicyWithContext(context.Background(), createV2PolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateV2PolicyWithContext is an alternate form of the CreateV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreateV2PolicyWithContext(ctx context.Context, createV2PolicyOptions *CreateV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createV2PolicyOptions, "createV2PolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createV2PolicyOptions, "createV2PolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v2/policies`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_v2_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2Policy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
// To update an access policy, use **`"type": "access"`** in the body. The supported subject attributes are **`iam_id`**
// and **`access_group_id`**. Use the **`iam_id`** subject attribute to assign access to a user or service-id. Use the
// **`access_group_id`** subject attribute to assign access to an access group. Assign roles that are supported by the
// service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). Use only the resource attributes supported by the
// service. To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs).
// The policy resource must include either the **`serviceType`**, **`serviceName`**, **`resourceGroupId`** or
// **`service_group_id`** attribute and the **`accountId`** attribute. In the rule field, you can specify a single
// condition by using **`key`**, **`value`**, and condition **`operator`**, or a set of **`conditions`** with a
// combination **`operator`**. The possible combination operators are **`and`** and **`or`**.
//
// Currently, we support two types of patterns:
//
// 1. `time-based`: Used to specify a time-based restriction
//
// Combine conditions to specify a time-based restriction (e.g., access only during business hours, during the
// Monday-Friday work week). For example, a policy can grant access Monday-Friday, 9:00am-5:00pm using the following
// rule:
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
//   'dateLessThan', 'dateLessThanOrEquals', 'dateGreaterThan', 'dateGreaterThanOrEquals',
//   'dateTimeLessThan', 'dateTimeLessThanOrEquals', 'dateTimeGreaterThan', 'dateTimeGreaterThanOrEquals',
//   'dayOfWeekEquals', 'dayOfWeekAnyOf'
// ``` The pattern field that matches the rule is required when rule is provided. For the business hour rule example
// above, the **`pattern`** is **`"time-based-conditions:weekly"`**. For more information, see [Time-based conditions
// operators](/docs/account?topic=account-iam-condition-properties&interface=ui#policy-condition-properties) and
// [Limiting access with time-based conditions](/docs/account?topic=account-iam-time-based&interface=ui). If the subject
// is a locked service-id, the request will fail.
//
// 2. `attribute-based`: Used to specify a combination of OR/AND based conditions applied on resource attributes.
//
// Combine conditions to specify an attribute-based condition using AND/OR-based operators.
//
// For example, a policy can grant access based on multiple conditions applied on the resource attributes below:
// ```json
//   "pattern": "attribute-based-condition:resource:literal-and-wildcard"
//   "rule": {
//       "operator": "or",
//       "conditions": [
//         {
//           "operator": "and",
//           "conditions": [
//             {
//               "key": "{{resource.attributes.prefix}}",
//               "operator": "stringEquals",
//               "value": "home/test"
//             },
//             {
//               "key": "{{environment.attributes.delimiter}}",
//               "operator": "stringEquals",
//               "value": "/"
//             }
//           ]
//         },
//         {
//           "key": "{{resource.attributes.path}}",
//           "operator": "stringMatch",
//           "value": "home/David/_*"
//         }
//       ]
//   }
// ```
//
// In addition to satisfying the `resources` section, the policy grants permission only if either the `path` begins with
// `home/David/` **OR**  the `prefix` is `home/test` and the `delimiter` is `/`. This mechanism helps you consolidate
// multiple policies in to a single policy,  making policies easier to administer and stay within the policy limit for
// an account. View the list of operators that can be used in the condition
// [here](/docs/account?topic=account-wildcard#string-comparisons).
//
// ### Authorization
//
// To update an authorization policy, use **`"type": "authorization"`** in the body. The subject attributes must match
// the supported authorization subjects of the resource. Multiple subject attributes might be provided. The following
// attributes are supported:
//   serviceName, serviceInstance, region, resourceType, resource, accountId, resourceGroupId Assign roles that are
// supported by the service or platform roles. For more information, see [IAM roles and
// actions](/docs/account?topic=account-iam-service-roles-actions). The user must also have the same level of access or
// greater to the target resource in order to grant the role. Use only the resource attributes supported by the service.
// To view a service's or the platform's supported attributes, check the [documentation](/docs?tab=all-docs). Both the
// policy subject and the policy resource must include the **`accountId`** attributes. The policy subject must include
// either **`serviceName`** or **`resourceGroupId`** (or both) attributes.
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
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceV2Policy(replaceV2PolicyOptions *ReplaceV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ReplaceV2PolicyWithContext(context.Background(), replaceV2PolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceV2PolicyWithContext is an alternate form of the ReplaceV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplaceV2PolicyWithContext(ctx context.Context, replaceV2PolicyOptions *ReplaceV2PolicyOptions) (result *V2Policy, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceV2PolicyOptions, "replaceV2PolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceV2PolicyOptions, "replaceV2PolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_v2_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2Policy)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetV2Policy : Retrieve a policy by ID
// Retrieve a policy by providing a policy ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetV2Policy(getV2PolicyOptions *GetV2PolicyOptions) (result *V2PolicyTemplateMetaData, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetV2PolicyWithContext(context.Background(), getV2PolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetV2PolicyWithContext is an alternate form of the GetV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetV2PolicyWithContext(ctx context.Context, getV2PolicyOptions *GetV2PolicyOptions) (result *V2PolicyTemplateMetaData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getV2PolicyOptions, "getV2PolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getV2PolicyOptions, "getV2PolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_v2_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalV2PolicyTemplateMetaData)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
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
	response, err = iamPolicyManagement.DeleteV2PolicyWithContext(context.Background(), deleteV2PolicyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteV2PolicyWithContext is an alternate form of the DeleteV2Policy method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeleteV2PolicyWithContext(ctx context.Context, deleteV2PolicyOptions *DeleteV2PolicyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteV2PolicyOptions, "deleteV2PolicyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteV2PolicyOptions, "deleteV2PolicyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_v2_policy", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListPolicyTemplates : List policy templates by attributes
// List policy templates and filter by attributes by using query parameters. The following attributes are supported:
// `account_id`, `policy_service_name`, `policy_service_type`, `policy_service_group_id` and `policy_type`.
// `account_id` is a required query parameter. These attributes `policy_service_name`, `policy_service_type` and
// `policy_service_group_id` are mutually exclusive. Only policy templates that have the specified attributes and that
// the caller has read access to are returned. If the caller does not have read access to any policy templates an empty
// array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyTemplates(listPolicyTemplatesOptions *ListPolicyTemplatesOptions) (result *PolicyTemplateCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListPolicyTemplatesWithContext(context.Background(), listPolicyTemplatesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPolicyTemplatesWithContext is an alternate form of the ListPolicyTemplates method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyTemplatesWithContext(ctx context.Context, listPolicyTemplatesOptions *ListPolicyTemplatesOptions) (result *PolicyTemplateCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPolicyTemplatesOptions, "listPolicyTemplatesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPolicyTemplatesOptions, "listPolicyTemplatesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPolicyTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListPolicyTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listPolicyTemplatesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listPolicyTemplatesOptions.AcceptLanguage))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listPolicyTemplatesOptions.AccountID))
	if listPolicyTemplatesOptions.State != nil {
		builder.AddQuery("state", fmt.Sprint(*listPolicyTemplatesOptions.State))
	}
	if listPolicyTemplatesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listPolicyTemplatesOptions.Name))
	}
	if listPolicyTemplatesOptions.PolicyServiceType != nil {
		builder.AddQuery("policy_service_type", fmt.Sprint(*listPolicyTemplatesOptions.PolicyServiceType))
	}
	if listPolicyTemplatesOptions.PolicyServiceName != nil {
		builder.AddQuery("policy_service_name", fmt.Sprint(*listPolicyTemplatesOptions.PolicyServiceName))
	}
	if listPolicyTemplatesOptions.PolicyServiceGroupID != nil {
		builder.AddQuery("policy_service_group_id", fmt.Sprint(*listPolicyTemplatesOptions.PolicyServiceGroupID))
	}
	if listPolicyTemplatesOptions.PolicyType != nil {
		builder.AddQuery("policy_type", fmt.Sprint(*listPolicyTemplatesOptions.PolicyType))
	}
	if listPolicyTemplatesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPolicyTemplatesOptions.Limit))
	}
	if listPolicyTemplatesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listPolicyTemplatesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_policy_templates", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreatePolicyTemplate : Create a policy template
// Create a policy template. Policy templates define a policy without requiring a subject, and you can use them to grant
// access to multiple subjects.
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplate(createPolicyTemplateOptions *CreatePolicyTemplateOptions) (result *PolicyTemplateLimitData, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.CreatePolicyTemplateWithContext(context.Background(), createPolicyTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePolicyTemplateWithContext is an alternate form of the CreatePolicyTemplate method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplateWithContext(ctx context.Context, createPolicyTemplateOptions *CreatePolicyTemplateOptions) (result *PolicyTemplateLimitData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPolicyTemplateOptions, "createPolicyTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPolicyTemplateOptions, "createPolicyTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPolicyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreatePolicyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createPolicyTemplateOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createPolicyTemplateOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createPolicyTemplateOptions.Name != nil {
		body["name"] = createPolicyTemplateOptions.Name
	}
	if createPolicyTemplateOptions.AccountID != nil {
		body["account_id"] = createPolicyTemplateOptions.AccountID
	}
	if createPolicyTemplateOptions.Policy != nil {
		body["policy"] = createPolicyTemplateOptions.Policy
	}
	if createPolicyTemplateOptions.Description != nil {
		body["description"] = createPolicyTemplateOptions.Description
	}
	if createPolicyTemplateOptions.Committed != nil {
		body["committed"] = createPolicyTemplateOptions.Committed
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_policy_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateLimitData)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetPolicyTemplate : Retrieve latest version of a policy template
// Retrieve the latest version of a policy template by providing a policy template ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyTemplate(getPolicyTemplateOptions *GetPolicyTemplateOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetPolicyTemplateWithContext(context.Background(), getPolicyTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPolicyTemplateWithContext is an alternate form of the GetPolicyTemplate method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyTemplateWithContext(ctx context.Context, getPolicyTemplateOptions *GetPolicyTemplateOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyTemplateOptions, "getPolicyTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPolicyTemplateOptions, "getPolicyTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *getPolicyTemplateOptions.PolicyTemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPolicyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetPolicyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getPolicyTemplateOptions.State != nil {
		builder.AddQuery("state", fmt.Sprint(*getPolicyTemplateOptions.State))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_policy_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplate)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeletePolicyTemplate : Delete a policy template
// Delete a policy template by providing the policy template ID. This deletes all versions of this template. A policy
// template can't be deleted if any version of the template is assigned to one or more child accounts. You must remove
// the policy assignments first.
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyTemplate(deletePolicyTemplateOptions *DeletePolicyTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamPolicyManagement.DeletePolicyTemplateWithContext(context.Background(), deletePolicyTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePolicyTemplateWithContext is an alternate form of the DeletePolicyTemplate method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyTemplateWithContext(ctx context.Context, deletePolicyTemplateOptions *DeletePolicyTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePolicyTemplateOptions, "deletePolicyTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePolicyTemplateOptions, "deletePolicyTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *deletePolicyTemplateOptions.PolicyTemplateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePolicyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeletePolicyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_policy_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreatePolicyTemplateVersion : Create a new policy template version
// Create a new version of a policy template. Use this if you need to make updates to a policy template that is
// committed.
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplateVersion(createPolicyTemplateVersionOptions *CreatePolicyTemplateVersionOptions) (result *PolicyTemplateLimitData, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.CreatePolicyTemplateVersionWithContext(context.Background(), createPolicyTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePolicyTemplateVersionWithContext is an alternate form of the CreatePolicyTemplateVersion method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplateVersionWithContext(ctx context.Context, createPolicyTemplateVersionOptions *CreatePolicyTemplateVersionOptions) (result *PolicyTemplateLimitData, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPolicyTemplateVersionOptions, "createPolicyTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPolicyTemplateVersionOptions, "createPolicyTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *createPolicyTemplateVersionOptions.PolicyTemplateID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPolicyTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreatePolicyTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createPolicyTemplateVersionOptions.Policy != nil {
		body["policy"] = createPolicyTemplateVersionOptions.Policy
	}
	if createPolicyTemplateVersionOptions.Name != nil {
		body["name"] = createPolicyTemplateVersionOptions.Name
	}
	if createPolicyTemplateVersionOptions.Description != nil {
		body["description"] = createPolicyTemplateVersionOptions.Description
	}
	if createPolicyTemplateVersionOptions.Committed != nil {
		body["committed"] = createPolicyTemplateVersionOptions.Committed
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_policy_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateLimitData)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListPolicyTemplateVersions : Retrieve policy template versions
// Retrieve the versions of a policy template by providing a policy template ID.
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyTemplateVersions(listPolicyTemplateVersionsOptions *ListPolicyTemplateVersionsOptions) (result *PolicyTemplateVersionsCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListPolicyTemplateVersionsWithContext(context.Background(), listPolicyTemplateVersionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPolicyTemplateVersionsWithContext is an alternate form of the ListPolicyTemplateVersions method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyTemplateVersionsWithContext(ctx context.Context, listPolicyTemplateVersionsOptions *ListPolicyTemplateVersionsOptions) (result *PolicyTemplateVersionsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPolicyTemplateVersionsOptions, "listPolicyTemplateVersionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPolicyTemplateVersionsOptions, "listPolicyTemplateVersionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *listPolicyTemplateVersionsOptions.PolicyTemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPolicyTemplateVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListPolicyTemplateVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listPolicyTemplateVersionsOptions.State != nil {
		builder.AddQuery("state", fmt.Sprint(*listPolicyTemplateVersionsOptions.State))
	}
	if listPolicyTemplateVersionsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPolicyTemplateVersionsOptions.Limit))
	}
	if listPolicyTemplateVersionsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listPolicyTemplateVersionsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_policy_template_versions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateVersionsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplacePolicyTemplate : Update a policy template version
// Update a specific version of a policy template. You can use this only if the version isn't committed.
func (iamPolicyManagement *IamPolicyManagementV1) ReplacePolicyTemplate(replacePolicyTemplateOptions *ReplacePolicyTemplateOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ReplacePolicyTemplateWithContext(context.Background(), replacePolicyTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplacePolicyTemplateWithContext is an alternate form of the ReplacePolicyTemplate method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ReplacePolicyTemplateWithContext(ctx context.Context, replacePolicyTemplateOptions *ReplacePolicyTemplateOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replacePolicyTemplateOptions, "replacePolicyTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replacePolicyTemplateOptions, "replacePolicyTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *replacePolicyTemplateOptions.PolicyTemplateID,
		"version": *replacePolicyTemplateOptions.Version,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replacePolicyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ReplacePolicyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replacePolicyTemplateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replacePolicyTemplateOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replacePolicyTemplateOptions.Policy != nil {
		body["policy"] = replacePolicyTemplateOptions.Policy
	}
	if replacePolicyTemplateOptions.Name != nil {
		body["name"] = replacePolicyTemplateOptions.Name
	}
	if replacePolicyTemplateOptions.Description != nil {
		body["description"] = replacePolicyTemplateOptions.Description
	}
	if replacePolicyTemplateOptions.Committed != nil {
		body["committed"] = replacePolicyTemplateOptions.Committed
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_policy_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplate)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeletePolicyTemplateVersion : Delete a policy template version
// Delete a specific version of a policy template by providing a policy template ID and version number. You can't delete
// a policy template version that is assigned to one or more child accounts. You must remove the policy assignments
// first.
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyTemplateVersion(deletePolicyTemplateVersionOptions *DeletePolicyTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	response, err = iamPolicyManagement.DeletePolicyTemplateVersionWithContext(context.Background(), deletePolicyTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePolicyTemplateVersionWithContext is an alternate form of the DeletePolicyTemplateVersion method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyTemplateVersionWithContext(ctx context.Context, deletePolicyTemplateVersionOptions *DeletePolicyTemplateVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePolicyTemplateVersionOptions, "deletePolicyTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePolicyTemplateVersionOptions, "deletePolicyTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *deletePolicyTemplateVersionOptions.PolicyTemplateID,
		"version": *deletePolicyTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePolicyTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeletePolicyTemplateVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_policy_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetPolicyTemplateVersion : Retrieve a policy template version
// Retrieve a policy template by providing a policy template ID and version number.
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyTemplateVersion(getPolicyTemplateVersionOptions *GetPolicyTemplateVersionOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetPolicyTemplateVersionWithContext(context.Background(), getPolicyTemplateVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPolicyTemplateVersionWithContext is an alternate form of the GetPolicyTemplateVersion method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyTemplateVersionWithContext(ctx context.Context, getPolicyTemplateVersionOptions *GetPolicyTemplateVersionOptions) (result *PolicyTemplate, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyTemplateVersionOptions, "getPolicyTemplateVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPolicyTemplateVersionOptions, "getPolicyTemplateVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *getPolicyTemplateVersionOptions.PolicyTemplateID,
		"version": *getPolicyTemplateVersionOptions.Version,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions/{version}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPolicyTemplateVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetPolicyTemplateVersion")
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_policy_template_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplate)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CommitPolicyTemplate : Commit a policy template version
// Commit a policy template version. You can make no further changes to the policy template once it's committed. If you
// need to make updates after committing a version, create a new version.
func (iamPolicyManagement *IamPolicyManagementV1) CommitPolicyTemplate(commitPolicyTemplateOptions *CommitPolicyTemplateOptions) (response *core.DetailedResponse, err error) {
	response, err = iamPolicyManagement.CommitPolicyTemplateWithContext(context.Background(), commitPolicyTemplateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CommitPolicyTemplateWithContext is an alternate form of the CommitPolicyTemplate method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CommitPolicyTemplateWithContext(ctx context.Context, commitPolicyTemplateOptions *CommitPolicyTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(commitPolicyTemplateOptions, "commitPolicyTemplateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(commitPolicyTemplateOptions, "commitPolicyTemplateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"policy_template_id": *commitPolicyTemplateOptions.PolicyTemplateID,
		"version": *commitPolicyTemplateOptions.Version,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_templates/{policy_template_id}/versions/{version}/commit`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range commitPolicyTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CommitPolicyTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "commit_policy_template", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListPolicyAssignments : Get policy template assignments
// Get policy template assignments by attributes. The following attributes are supported:
// `account_id`, `template_id`, `template_version`, `sort`.
// `account_id` is a required query parameter. Only policy template assignments that have the specified attributes and
// that the caller has read access to are returned. If the caller does not have read access to any policy template
// assignments an empty array is returned.
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyAssignments(listPolicyAssignmentsOptions *ListPolicyAssignmentsOptions) (result *PolicyTemplateAssignmentCollection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.ListPolicyAssignmentsWithContext(context.Background(), listPolicyAssignmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPolicyAssignmentsWithContext is an alternate form of the ListPolicyAssignments method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) ListPolicyAssignmentsWithContext(ctx context.Context, listPolicyAssignmentsOptions *ListPolicyAssignmentsOptions) (result *PolicyTemplateAssignmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPolicyAssignmentsOptions, "listPolicyAssignmentsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPolicyAssignmentsOptions, "listPolicyAssignmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_assignments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPolicyAssignmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "ListPolicyAssignments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listPolicyAssignmentsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listPolicyAssignmentsOptions.AcceptLanguage))
	}

	builder.AddQuery("version", fmt.Sprint(*listPolicyAssignmentsOptions.Version))
	builder.AddQuery("account_id", fmt.Sprint(*listPolicyAssignmentsOptions.AccountID))
	if listPolicyAssignmentsOptions.TemplateID != nil {
		builder.AddQuery("template_id", fmt.Sprint(*listPolicyAssignmentsOptions.TemplateID))
	}
	if listPolicyAssignmentsOptions.TemplateVersion != nil {
		builder.AddQuery("template_version", fmt.Sprint(*listPolicyAssignmentsOptions.TemplateVersion))
	}
	if listPolicyAssignmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPolicyAssignmentsOptions.Limit))
	}
	if listPolicyAssignmentsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listPolicyAssignmentsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_policy_assignments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateAssignmentCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreatePolicyTemplateAssignment : Create a policy authorization template assignment
// Assign a policy template to child accounts and account groups. This creates the policy in the accounts and account
// groups that you specify.
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplateAssignment(createPolicyTemplateAssignmentOptions *CreatePolicyTemplateAssignmentOptions) (result *PolicyAssignmentV1Collection, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.CreatePolicyTemplateAssignmentWithContext(context.Background(), createPolicyTemplateAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePolicyTemplateAssignmentWithContext is an alternate form of the CreatePolicyTemplateAssignment method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) CreatePolicyTemplateAssignmentWithContext(ctx context.Context, createPolicyTemplateAssignmentOptions *CreatePolicyTemplateAssignmentOptions) (result *PolicyAssignmentV1Collection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPolicyTemplateAssignmentOptions, "createPolicyTemplateAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPolicyTemplateAssignmentOptions, "createPolicyTemplateAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_assignments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPolicyTemplateAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "CreatePolicyTemplateAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createPolicyTemplateAssignmentOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createPolicyTemplateAssignmentOptions.AcceptLanguage))
	}

	builder.AddQuery("version", fmt.Sprint(*createPolicyTemplateAssignmentOptions.Version))

	body := make(map[string]interface{})
	if createPolicyTemplateAssignmentOptions.Target != nil {
		body["target"] = createPolicyTemplateAssignmentOptions.Target
	}
	if createPolicyTemplateAssignmentOptions.Templates != nil {
		body["templates"] = createPolicyTemplateAssignmentOptions.Templates
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_policy_template_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyAssignmentV1Collection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetPolicyAssignment : Retrieve a policy assignment
// Retrieve a policy template assignment by providing a policy assignment ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyAssignment(getPolicyAssignmentOptions *GetPolicyAssignmentOptions) (result PolicyTemplateAssignmentItemsIntf, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetPolicyAssignmentWithContext(context.Background(), getPolicyAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPolicyAssignmentWithContext is an alternate form of the GetPolicyAssignment method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetPolicyAssignmentWithContext(ctx context.Context, getPolicyAssignmentOptions *GetPolicyAssignmentOptions) (result PolicyTemplateAssignmentItemsIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPolicyAssignmentOptions, "getPolicyAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPolicyAssignmentOptions, "getPolicyAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *getPolicyAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPolicyAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetPolicyAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("version", fmt.Sprint(*getPolicyAssignmentOptions.Version))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_policy_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyTemplateAssignmentItems)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdatePolicyAssignment : Update a policy authorization type assignment
// Update a policy assignment by providing a policy assignment ID.
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyAssignment(updatePolicyAssignmentOptions *UpdatePolicyAssignmentOptions) (result *PolicyAssignmentV1, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.UpdatePolicyAssignmentWithContext(context.Background(), updatePolicyAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdatePolicyAssignmentWithContext is an alternate form of the UpdatePolicyAssignment method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) UpdatePolicyAssignmentWithContext(ctx context.Context, updatePolicyAssignmentOptions *UpdatePolicyAssignmentOptions) (result *PolicyAssignmentV1, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePolicyAssignmentOptions, "updatePolicyAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updatePolicyAssignmentOptions, "updatePolicyAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *updatePolicyAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updatePolicyAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "UpdatePolicyAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updatePolicyAssignmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updatePolicyAssignmentOptions.IfMatch))
	}

	builder.AddQuery("version", fmt.Sprint(*updatePolicyAssignmentOptions.Version))

	body := make(map[string]interface{})
	if updatePolicyAssignmentOptions.TemplateVersion != nil {
		body["template_version"] = updatePolicyAssignmentOptions.TemplateVersion
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_policy_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPolicyAssignmentV1)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeletePolicyAssignment : Remove a policy assignment
// Remove a policy template assignment by providing a policy assignment ID. You can't delete a policy assignment if the
// status is "in_progress".
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyAssignment(deletePolicyAssignmentOptions *DeletePolicyAssignmentOptions) (response *core.DetailedResponse, err error) {
	response, err = iamPolicyManagement.DeletePolicyAssignmentWithContext(context.Background(), deletePolicyAssignmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePolicyAssignmentWithContext is an alternate form of the DeletePolicyAssignment method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) DeletePolicyAssignmentWithContext(ctx context.Context, deletePolicyAssignmentOptions *DeletePolicyAssignmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePolicyAssignmentOptions, "deletePolicyAssignmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePolicyAssignmentOptions, "deletePolicyAssignmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"assignment_id": *deletePolicyAssignmentOptions.AssignmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/policy_assignments/{assignment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePolicyAssignmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "DeletePolicyAssignment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = iamPolicyManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_policy_assignment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetSettings : Retrieve Access Management account settings by account ID
// Retrieve Access Management settings for an account by providing the account ID.
func (iamPolicyManagement *IamPolicyManagementV1) GetSettings(getSettingsOptions *GetSettingsOptions) (result *AccountSettingsAccessManagement, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.GetSettingsWithContext(context.Background(), getSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *AccountSettingsAccessManagement, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSettingsOptions, "getSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/accounts/{account_id}/settings/access_management`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getSettingsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getSettingsOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsAccessManagement)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateSettings : Update Access Management account settings by account ID
// Update access management settings for an account.
//
// ### External Account Identity Interaction
//
// Update the way identities within an external account are allowed to interact with the requested account by providing:
// * the `account_id` as a parameter
// * the external account ID(s) and state for the specific identity in the request body
//
// External account identity interaction includes the following `identity_types`: `user` (user identities defined as
// [IBMid's](https://test.cloud.ibm.com/docs/account?topic=account-identity-overview#users-bestpract)), `service_id`
// (defined as [IAM
// ServiceIds](https://test.cloud.ibm.com/docs/account?topic=account-identity-overview#serviceid-bestpract)), `service`
// (defined by a service’s [CRN](https://test.cloud.ibm.com/docs/account?topic=account-crn)). To update an Identity’s
// setting, the `state` and `external_allowed_accounts` fields are required.
//
// Different identity states are:
// * "enabled": An identity type is allowed to access resources in the account provided it has access policies on those
// resources.
// * "limited": An identity type is allowed to access resources in the account provided it has access policies on those
// resources AND it is associated with either the account the resources are in or one of the allowed accounts. This
// setting leverages the "external_allowed_accounts" list.
// * "monitor": Has no direct impact on an Identity’s access. Instead, it creates AT events for access decisions as if
// the account were in a limited “state”.
//
// **Note**: The state "enabled" is a special case. In this case, access is given to all accounts and there is no need
// to specify a particular list. Therefore, when updating "state" to "enabled" for an identity type
// "external_allowed_accounts" should be left empty.
func (iamPolicyManagement *IamPolicyManagementV1) UpdateSettings(updateSettingsOptions *UpdateSettingsOptions) (result *AccountSettingsAccessManagement, response *core.DetailedResponse, err error) {
	result, response, err = iamPolicyManagement.UpdateSettingsWithContext(context.Background(), updateSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateSettingsWithContext is an alternate form of the UpdateSettings method which supports a Context parameter
func (iamPolicyManagement *IamPolicyManagementV1) UpdateSettingsWithContext(ctx context.Context, updateSettingsOptions *UpdateSettingsOptions) (result *AccountSettingsAccessManagement, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSettingsOptions, "updateSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateSettingsOptions, "updateSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *updateSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = iamPolicyManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(iamPolicyManagement.Service.Options.URL, `/v1/accounts/{account_id}/settings/access_management`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("iam_policy_management", "V1", "UpdateSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateSettingsOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateSettingsOptions.IfMatch))
	}
	if updateSettingsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*updateSettingsOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if updateSettingsOptions.ExternalAccountIdentityInteraction != nil {
		body["external_account_identity_interaction"] = updateSettingsOptions.ExternalAccountIdentityInteraction
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
	response, err = iamPolicyManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettingsAccessManagement)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.1")
}

// AccountSettingsAccessManagement : The Access Management Account Settings that are currently set for the requested account.
type AccountSettingsAccessManagement struct {
	// How external accounts can interact in relation to the requested account.
	ExternalAccountIdentityInteraction *ExternalAccountIdentityInteraction `json:"external_account_identity_interaction" validate:"required"`
}

// UnmarshalAccountSettingsAccessManagement unmarshals an instance of AccountSettingsAccessManagement from the specified map of raw messages.
func UnmarshalAccountSettingsAccessManagement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettingsAccessManagement)
	err = core.UnmarshalModel(m, "external_account_identity_interaction", &obj.ExternalAccountIdentityInteraction, UnmarshalExternalAccountIdentityInteraction)
	if err != nil {
		err = core.SDKErrorf(err, "", "external_account_identity_interaction-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssignmentResourceCreated : On success, includes the  policy assigned.
type AssignmentResourceCreated struct {
	// policy id.
	ID *string `json:"id,omitempty"`
}

// UnmarshalAssignmentResourceCreated unmarshals an instance of AssignmentResourceCreated from the specified map of raw messages.
func UnmarshalAssignmentResourceCreated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssignmentResourceCreated)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssignmentTargetDetails : assignment target account and type.
type AssignmentTargetDetails struct {
	// Assignment target type.
	Type *string `json:"type,omitempty"`

	// ID of the target account.
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AssignmentTargetDetails.Type property.
// Assignment target type.
const (
	AssignmentTargetDetailsTypeAccountConst = "Account"
)

// UnmarshalAssignmentTargetDetails unmarshals an instance of AssignmentTargetDetails from the specified map of raw messages.
func UnmarshalAssignmentTargetDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssignmentTargetDetails)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
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

// AssignmentTemplateDetails : policy template details.
type AssignmentTemplateDetails struct {
	// policy template id.
	ID *string `json:"id,omitempty"`

	// policy template version.
	Version *string `json:"version,omitempty"`
}

// UnmarshalAssignmentTemplateDetails unmarshals an instance of AssignmentTemplateDetails from the specified map of raw messages.
func UnmarshalAssignmentTemplateDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssignmentTemplateDetails)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CommitPolicyTemplateOptions : The CommitPolicyTemplate options.
type CommitPolicyTemplateOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template version.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCommitPolicyTemplateOptions : Instantiate CommitPolicyTemplateOptions
func (*IamPolicyManagementV1) NewCommitPolicyTemplateOptions(policyTemplateID string, version string) *CommitPolicyTemplateOptions {
	return &CommitPolicyTemplateOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
		Version: core.StringPtr(version),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *CommitPolicyTemplateOptions) SetPolicyTemplateID(policyTemplateID string) *CommitPolicyTemplateOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CommitPolicyTemplateOptions) SetVersion(version string) *CommitPolicyTemplateOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CommitPolicyTemplateOptions) SetHeaders(param map[string]string) *CommitPolicyTemplateOptions {
	options.Headers = param
	return options
}

// ConflictsWith : Details of conflicting resource.
type ConflictsWith struct {
	// The revision number of the resource.
	Etag *string `json:"etag,omitempty"`

	// The conflicting role id.
	Role *string `json:"role,omitempty"`

	// The conflicting policy id.
	Policy *string `json:"policy,omitempty"`
}

// UnmarshalConflictsWith unmarshals an instance of ConflictsWith from the specified map of raw messages.
func UnmarshalConflictsWith(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConflictsWith)
	err = core.UnmarshalPrimitive(m, "etag", &obj.Etag)
	if err != nil {
		err = core.SDKErrorf(err, "", "etag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "role", &obj.Role)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy", &obj.Policy)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalControl unmarshals an instance of Control from the specified map of raw messages.
func UnmarshalControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Control)
	err = core.UnmarshalModel(m, "grant", &obj.Grant, UnmarshalGrant)
	if err != nil {
		err = core.SDKErrorf(err, "", "grant-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "grant-error", common.GetComponentInfo())
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

	// Allows users to set headers on API requests.
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

// CreatePolicyTemplateAssignmentOptions : The CreatePolicyTemplateAssignment options.
type CreatePolicyTemplateAssignmentOptions struct {
	// specify version of response body format.
	Version *string `json:"version" validate:"required"`

	// assignment target account and type.
	Target *AssignmentTargetDetails `json:"target" validate:"required"`

	// List of template details for policy assignment.
	Templates []AssignmentTemplateDetails `json:"templates" validate:"required"`

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

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreatePolicyTemplateAssignmentOptions : Instantiate CreatePolicyTemplateAssignmentOptions
func (*IamPolicyManagementV1) NewCreatePolicyTemplateAssignmentOptions(version string, target *AssignmentTargetDetails, templates []AssignmentTemplateDetails) *CreatePolicyTemplateAssignmentOptions {
	return &CreatePolicyTemplateAssignmentOptions{
		Version: core.StringPtr(version),
		Target: target,
		Templates: templates,
	}
}

// SetVersion : Allow user to set Version
func (_options *CreatePolicyTemplateAssignmentOptions) SetVersion(version string) *CreatePolicyTemplateAssignmentOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *CreatePolicyTemplateAssignmentOptions) SetTarget(target *AssignmentTargetDetails) *CreatePolicyTemplateAssignmentOptions {
	_options.Target = target
	return _options
}

// SetTemplates : Allow user to set Templates
func (_options *CreatePolicyTemplateAssignmentOptions) SetTemplates(templates []AssignmentTemplateDetails) *CreatePolicyTemplateAssignmentOptions {
	_options.Templates = templates
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreatePolicyTemplateAssignmentOptions) SetAcceptLanguage(acceptLanguage string) *CreatePolicyTemplateAssignmentOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePolicyTemplateAssignmentOptions) SetHeaders(param map[string]string) *CreatePolicyTemplateAssignmentOptions {
	options.Headers = param
	return options
}

// CreatePolicyTemplateOptions : The CreatePolicyTemplate options.
type CreatePolicyTemplateOptions struct {
	// Required field when creating a new template. Otherwise this field is optional. If the field is included it will
	// change the name value for all existing versions of the template.
	Name *string `json:"name" validate:"required"`

	// Enterprise account ID where this template will be created.
	AccountID *string `json:"account_id" validate:"required"`

	// The core set of properties associated with the template's policy objet.
	Policy *TemplatePolicy `json:"policy" validate:"required"`

	// Description of the policy template. This is shown to users in the enterprise account. Use this to describe the
	// purpose or context of the policy for enterprise users managing IAM templates.
	Description *string `json:"description,omitempty"`

	// Committed status of the template.
	Committed *bool `json:"committed,omitempty"`

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

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreatePolicyTemplateOptions : Instantiate CreatePolicyTemplateOptions
func (*IamPolicyManagementV1) NewCreatePolicyTemplateOptions(name string, accountID string, policy *TemplatePolicy) *CreatePolicyTemplateOptions {
	return &CreatePolicyTemplateOptions{
		Name: core.StringPtr(name),
		AccountID: core.StringPtr(accountID),
		Policy: policy,
	}
}

// SetName : Allow user to set Name
func (_options *CreatePolicyTemplateOptions) SetName(name string) *CreatePolicyTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreatePolicyTemplateOptions) SetAccountID(accountID string) *CreatePolicyTemplateOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetPolicy : Allow user to set Policy
func (_options *CreatePolicyTemplateOptions) SetPolicy(policy *TemplatePolicy) *CreatePolicyTemplateOptions {
	_options.Policy = policy
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreatePolicyTemplateOptions) SetDescription(description string) *CreatePolicyTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetCommitted : Allow user to set Committed
func (_options *CreatePolicyTemplateOptions) SetCommitted(committed bool) *CreatePolicyTemplateOptions {
	_options.Committed = core.BoolPtr(committed)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreatePolicyTemplateOptions) SetAcceptLanguage(acceptLanguage string) *CreatePolicyTemplateOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePolicyTemplateOptions) SetHeaders(param map[string]string) *CreatePolicyTemplateOptions {
	options.Headers = param
	return options
}

// CreatePolicyTemplateVersionOptions : The CreatePolicyTemplateVersion options.
type CreatePolicyTemplateVersionOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The core set of properties associated with the template's policy objet.
	Policy *TemplatePolicy `json:"policy" validate:"required"`

	// Required field when creating a new template. Otherwise this field is optional. If the field is included it will
	// change the name value for all existing versions of the template.
	Name *string `json:"name,omitempty"`

	// Description of the policy template. This is shown to users in the enterprise account. Use this to describe the
	// purpose or context of the policy for enterprise users managing IAM templates.
	Description *string `json:"description,omitempty"`

	// Committed status of the template version.
	Committed *bool `json:"committed,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreatePolicyTemplateVersionOptions : Instantiate CreatePolicyTemplateVersionOptions
func (*IamPolicyManagementV1) NewCreatePolicyTemplateVersionOptions(policyTemplateID string, policy *TemplatePolicy) *CreatePolicyTemplateVersionOptions {
	return &CreatePolicyTemplateVersionOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
		Policy: policy,
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *CreatePolicyTemplateVersionOptions) SetPolicyTemplateID(policyTemplateID string) *CreatePolicyTemplateVersionOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetPolicy : Allow user to set Policy
func (_options *CreatePolicyTemplateVersionOptions) SetPolicy(policy *TemplatePolicy) *CreatePolicyTemplateVersionOptions {
	_options.Policy = policy
	return _options
}

// SetName : Allow user to set Name
func (_options *CreatePolicyTemplateVersionOptions) SetName(name string) *CreatePolicyTemplateVersionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreatePolicyTemplateVersionOptions) SetDescription(description string) *CreatePolicyTemplateVersionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetCommitted : Allow user to set Committed
func (_options *CreatePolicyTemplateVersionOptions) SetCommitted(committed bool) *CreatePolicyTemplateVersionOptions {
	_options.Committed = core.BoolPtr(committed)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePolicyTemplateVersionOptions) SetHeaders(param map[string]string) *CreatePolicyTemplateVersionOptions {
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

	// Allows users to set headers on API requests.
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

	// Description of the policy.
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

	// Allows users to set headers on API requests.
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
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "actions", &obj.Actions)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
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

// DeletePolicyAssignmentOptions : The DeletePolicyAssignment options.
type DeletePolicyAssignmentOptions struct {
	// The policy template assignment ID.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePolicyAssignmentOptions : Instantiate DeletePolicyAssignmentOptions
func (*IamPolicyManagementV1) NewDeletePolicyAssignmentOptions(assignmentID string) *DeletePolicyAssignmentOptions {
	return &DeletePolicyAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *DeletePolicyAssignmentOptions) SetAssignmentID(assignmentID string) *DeletePolicyAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePolicyAssignmentOptions) SetHeaders(param map[string]string) *DeletePolicyAssignmentOptions {
	options.Headers = param
	return options
}

// DeletePolicyOptions : The DeletePolicy options.
type DeletePolicyOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
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

// DeletePolicyTemplateOptions : The DeletePolicyTemplate options.
type DeletePolicyTemplateOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePolicyTemplateOptions : Instantiate DeletePolicyTemplateOptions
func (*IamPolicyManagementV1) NewDeletePolicyTemplateOptions(policyTemplateID string) *DeletePolicyTemplateOptions {
	return &DeletePolicyTemplateOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *DeletePolicyTemplateOptions) SetPolicyTemplateID(policyTemplateID string) *DeletePolicyTemplateOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePolicyTemplateOptions) SetHeaders(param map[string]string) *DeletePolicyTemplateOptions {
	options.Headers = param
	return options
}

// DeletePolicyTemplateVersionOptions : The DeletePolicyTemplateVersion options.
type DeletePolicyTemplateVersionOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template version.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePolicyTemplateVersionOptions : Instantiate DeletePolicyTemplateVersionOptions
func (*IamPolicyManagementV1) NewDeletePolicyTemplateVersionOptions(policyTemplateID string, version string) *DeletePolicyTemplateVersionOptions {
	return &DeletePolicyTemplateVersionOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
		Version: core.StringPtr(version),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *DeletePolicyTemplateVersionOptions) SetPolicyTemplateID(policyTemplateID string) *DeletePolicyTemplateVersionOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *DeletePolicyTemplateVersionOptions) SetVersion(version string) *DeletePolicyTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePolicyTemplateVersionOptions) SetHeaders(param map[string]string) *DeletePolicyTemplateVersionOptions {
	options.Headers = param
	return options
}

// DeleteRoleOptions : The DeleteRole options.
type DeleteRoleOptions struct {
	// The role ID.
	RoleID *string `json:"role_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
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

	// Allows users to set headers on API requests.
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
		err = core.SDKErrorf(err, "", "role_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalRoleAction)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ErrorDetails : Additional error details.
type ErrorDetails struct {
	// Details of conflicting resource.
	ConflictsWith *ConflictsWith `json:"conflicts_with,omitempty"`
}

// UnmarshalErrorDetails unmarshals an instance of ErrorDetails from the specified map of raw messages.
func UnmarshalErrorDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ErrorDetails)
	err = core.UnmarshalModel(m, "conflicts_with", &obj.ConflictsWith, UnmarshalConflictsWith)
	if err != nil {
		err = core.SDKErrorf(err, "", "conflicts_with-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ErrorObject : ErrorObject struct
type ErrorObject struct {
	// The API error code for the error.
	Code *string `json:"code" validate:"required"`

	// The error message returned by the API.
	Message *string `json:"message" validate:"required"`

	// Additional error details.
	Details *ErrorDetails `json:"details,omitempty"`

	// Additional info for error.
	MoreInfo *string `json:"more_info,omitempty"`
}

// Constants associated with the ErrorObject.Code property.
// The API error code for the error.
const (
	ErrorObjectCodeInsufficentPermissionsConst = "insufficent_permissions"
	ErrorObjectCodeInvalidBodyConst = "invalid_body"
	ErrorObjectCodeInvalidTokenConst = "invalid_token"
	ErrorObjectCodeMissingRequiredQueryParameterConst = "missing_required_query_parameter"
	ErrorObjectCodeNotFoundConst = "not_found"
	ErrorObjectCodePolicyAssignmentConflictErrorConst = "policy_assignment_conflict_error"
	ErrorObjectCodePolicyAssignmentNotFoundConst = "policy_assignment_not_found"
	ErrorObjectCodePolicyConflictErrorConst = "policy_conflict_error"
	ErrorObjectCodePolicyNotFoundConst = "policy_not_found"
	ErrorObjectCodePolicyTemplateConflictErrorConst = "policy_template_conflict_error"
	ErrorObjectCodePolicyTemplateNotFoundConst = "policy_template_not_found"
	ErrorObjectCodeRequestNotProcessedConst = "request_not_processed"
	ErrorObjectCodeResourceNotFoundConst = "resource_not_found"
	ErrorObjectCodeRoleConflictErrorConst = "role_conflict_error"
	ErrorObjectCodeRoleNotFoundConst = "role_not_found"
	ErrorObjectCodeTooManyRequestsConst = "too_many_requests"
	ErrorObjectCodeUnableToProcessConst = "unable_to_process"
	ErrorObjectCodeUnsupportedContentTypeConst = "unsupported_content_type"
)

// UnmarshalErrorObject unmarshals an instance of ErrorObject from the specified map of raw messages.
func UnmarshalErrorObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ErrorObject)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "details", &obj.Details, UnmarshalErrorDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "details-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "more_info", &obj.MoreInfo)
	if err != nil {
		err = core.SDKErrorf(err, "", "more_info-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ErrorResponse : The error response from API.
type ErrorResponse struct {
	// The unique transaction id for the request.
	Trace *string `json:"trace,omitempty"`

	// The errors encountered during the response.
	Errors []ErrorObject `json:"errors,omitempty"`

	// The http error code of the response.
	StatusCode *int64 `json:"status_code,omitempty"`
}

// UnmarshalErrorResponse unmarshals an instance of ErrorResponse from the specified map of raw messages.
func UnmarshalErrorResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ErrorResponse)
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		err = core.SDKErrorf(err, "", "trace-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalErrorObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_code-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalAccountIdentityInteraction : How external accounts can interact in relation to the requested account.
type ExternalAccountIdentityInteraction struct {
	// The settings for each identity type.
	IdentityTypes *IdentityTypes `json:"identity_types" validate:"required"`
}

// UnmarshalExternalAccountIdentityInteraction unmarshals an instance of ExternalAccountIdentityInteraction from the specified map of raw messages.
func UnmarshalExternalAccountIdentityInteraction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalAccountIdentityInteraction)
	err = core.UnmarshalModel(m, "identity_types", &obj.IdentityTypes, UnmarshalIdentityTypes)
	if err != nil {
		err = core.SDKErrorf(err, "", "identity_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExternalAccountIdentityInteractionPatch : Update to how external accounts can interact in relation to the requested account.
type ExternalAccountIdentityInteractionPatch struct {
	// The settings to apply for each identity type for a request.
	IdentityTypes *IdentityTypesPatch `json:"identity_types,omitempty"`
}

// UnmarshalExternalAccountIdentityInteractionPatch unmarshals an instance of ExternalAccountIdentityInteractionPatch from the specified map of raw messages.
func UnmarshalExternalAccountIdentityInteractionPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalAccountIdentityInteractionPatch)
	err = core.UnmarshalModel(m, "identity_types", &obj.IdentityTypes, UnmarshalIdentityTypesPatch)
	if err != nil {
		err = core.SDKErrorf(err, "", "identity_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// First : Details with href linking to first page of requested collection.
type First struct {
	// The href linking to the page of requested collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalFirst unmarshals an instance of First from the specified map of raw messages.
func UnmarshalFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(First)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetPolicyAssignmentOptions : The GetPolicyAssignment options.
type GetPolicyAssignmentOptions struct {
	// The policy template assignment ID.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// specify version of response body format.
	Version *string `json:"version" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPolicyAssignmentOptions : Instantiate GetPolicyAssignmentOptions
func (*IamPolicyManagementV1) NewGetPolicyAssignmentOptions(assignmentID string, version string) *GetPolicyAssignmentOptions {
	return &GetPolicyAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
		Version: core.StringPtr(version),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *GetPolicyAssignmentOptions) SetAssignmentID(assignmentID string) *GetPolicyAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetPolicyAssignmentOptions) SetVersion(version string) *GetPolicyAssignmentOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyAssignmentOptions) SetHeaders(param map[string]string) *GetPolicyAssignmentOptions {
	options.Headers = param
	return options
}

// GetPolicyOptions : The GetPolicy options.
type GetPolicyOptions struct {
	// The policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
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

// GetPolicyTemplateOptions : The GetPolicyTemplate options.
type GetPolicyTemplateOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template state.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetPolicyTemplateOptions.State property.
// The policy template state.
const (
	GetPolicyTemplateOptionsStateActiveConst = "active"
	GetPolicyTemplateOptionsStateDeletedConst = "deleted"
)

// NewGetPolicyTemplateOptions : Instantiate GetPolicyTemplateOptions
func (*IamPolicyManagementV1) NewGetPolicyTemplateOptions(policyTemplateID string) *GetPolicyTemplateOptions {
	return &GetPolicyTemplateOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *GetPolicyTemplateOptions) SetPolicyTemplateID(policyTemplateID string) *GetPolicyTemplateOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetState : Allow user to set State
func (_options *GetPolicyTemplateOptions) SetState(state string) *GetPolicyTemplateOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyTemplateOptions) SetHeaders(param map[string]string) *GetPolicyTemplateOptions {
	options.Headers = param
	return options
}

// GetPolicyTemplateVersionOptions : The GetPolicyTemplateVersion options.
type GetPolicyTemplateVersionOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template version.
	Version *string `json:"version" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetPolicyTemplateVersionOptions : Instantiate GetPolicyTemplateVersionOptions
func (*IamPolicyManagementV1) NewGetPolicyTemplateVersionOptions(policyTemplateID string, version string) *GetPolicyTemplateVersionOptions {
	return &GetPolicyTemplateVersionOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
		Version: core.StringPtr(version),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *GetPolicyTemplateVersionOptions) SetPolicyTemplateID(policyTemplateID string) *GetPolicyTemplateVersionOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *GetPolicyTemplateVersionOptions) SetVersion(version string) *GetPolicyTemplateVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPolicyTemplateVersionOptions) SetHeaders(param map[string]string) *GetPolicyTemplateVersionOptions {
	options.Headers = param
	return options
}

// GetRoleOptions : The GetRole options.
type GetRoleOptions struct {
	// The role ID.
	RoleID *string `json:"role_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
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

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {
	// The account GUID that the settings belong to.
	AccountID *string `json:"account_id" validate:"required,ne="`

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

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*IamPolicyManagementV1) NewGetSettingsOptions(accountID string) *GetSettingsOptions {
	return &GetSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetSettingsOptions) SetAccountID(accountID string) *GetSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetSettingsOptions) SetAcceptLanguage(acceptLanguage string) *GetSettingsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
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

	// Allows users to set headers on API requests.
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalGrant unmarshals an instance of Grant from the specified map of raw messages.
func UnmarshalGrant(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Grant)
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalRoles)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IdentityTypes : The settings for each identity type.
type IdentityTypes struct {
	// The core set of properties associated with an identity type.
	User *IdentityTypesBase `json:"user" validate:"required"`

	// The core set of properties associated with an identity type.
	ServiceID *IdentityTypesBase `json:"service_id" validate:"required"`

	// The core set of properties associated with an identity type.
	Service *IdentityTypesBase `json:"service" validate:"required"`
}

// UnmarshalIdentityTypes unmarshals an instance of IdentityTypes from the specified map of raw messages.
func UnmarshalIdentityTypes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IdentityTypes)
	err = core.UnmarshalModel(m, "user", &obj.User, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_id", &obj.ServiceID, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IdentityTypesBase : The core set of properties associated with an identity type.
type IdentityTypesBase struct {
	// The state of the identity type.
	State *string `json:"state" validate:"required"`

	// List of accounts that the state applies to for a given identity.
	ExternalAllowedAccounts []string `json:"external_allowed_accounts" validate:"required"`
}

// Constants associated with the IdentityTypesBase.State property.
// The state of the identity type.
const (
	IdentityTypesBaseStateEnabledConst = "enabled"
	IdentityTypesBaseStateLimitedConst = "limited"
	IdentityTypesBaseStateMonitorConst = "monitor"
)

// NewIdentityTypesBase : Instantiate IdentityTypesBase (Generic Model Constructor)
func (*IamPolicyManagementV1) NewIdentityTypesBase(state string, externalAllowedAccounts []string) (_model *IdentityTypesBase, err error) {
	_model = &IdentityTypesBase{
		State: core.StringPtr(state),
		ExternalAllowedAccounts: externalAllowedAccounts,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalIdentityTypesBase unmarshals an instance of IdentityTypesBase from the specified map of raw messages.
func UnmarshalIdentityTypesBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IdentityTypesBase)
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "external_allowed_accounts", &obj.ExternalAllowedAccounts)
	if err != nil {
		err = core.SDKErrorf(err, "", "external_allowed_accounts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IdentityTypesPatch : The settings to apply for each identity type for a request.
type IdentityTypesPatch struct {
	// The core set of properties associated with an identity type.
	User *IdentityTypesBase `json:"user,omitempty"`

	// The core set of properties associated with an identity type.
	ServiceID *IdentityTypesBase `json:"service_id,omitempty"`

	// The core set of properties associated with an identity type.
	Service *IdentityTypesBase `json:"service,omitempty"`
}

// UnmarshalIdentityTypesPatch unmarshals an instance of IdentityTypesPatch from the specified map of raw messages.
func UnmarshalIdentityTypesPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IdentityTypesPatch)
	err = core.UnmarshalModel(m, "user", &obj.User, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_id", &obj.ServiceID, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalIdentityTypesBase)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LimitData : policy template current and limit details with in an account.
type LimitData struct {
	// policy template current count.
	Current *int64 `json:"current,omitempty"`

	// policy template limit count.
	Limit *int64 `json:"limit,omitempty"`
}

// UnmarshalLimitData unmarshals an instance of LimitData from the specified map of raw messages.
func UnmarshalLimitData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LimitData)
	err = core.UnmarshalPrimitive(m, "current", &obj.Current)
	if err != nil {
		err = core.SDKErrorf(err, "", "current-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
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

	// The number of documents to include in collection.
	Limit *int64 `json:"limit,omitempty"`

	// Page token that refers to the page of collection to return.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
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

// SetLimit : Allow user to set Limit
func (_options *ListPoliciesOptions) SetLimit(limit int64) *ListPoliciesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListPoliciesOptions) SetStart(start string) *ListPoliciesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPoliciesOptions) SetHeaders(param map[string]string) *ListPoliciesOptions {
	options.Headers = param
	return options
}

// ListPolicyAssignmentsOptions : The ListPolicyAssignments options.
type ListPolicyAssignmentsOptions struct {
	// specify version of response body format.
	Version *string `json:"version" validate:"required"`

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

	// Optional template id.
	TemplateID *string `json:"template_id,omitempty"`

	// Optional policy template version.
	TemplateVersion *string `json:"template_version,omitempty"`

	// The number of documents to include in collection.
	Limit *int64 `json:"limit,omitempty"`

	// Page token that refers to the page of collection to return.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListPolicyAssignmentsOptions : Instantiate ListPolicyAssignmentsOptions
func (*IamPolicyManagementV1) NewListPolicyAssignmentsOptions(version string, accountID string) *ListPolicyAssignmentsOptions {
	return &ListPolicyAssignmentsOptions{
		Version: core.StringPtr(version),
		AccountID: core.StringPtr(accountID),
	}
}

// SetVersion : Allow user to set Version
func (_options *ListPolicyAssignmentsOptions) SetVersion(version string) *ListPolicyAssignmentsOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListPolicyAssignmentsOptions) SetAccountID(accountID string) *ListPolicyAssignmentsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListPolicyAssignmentsOptions) SetAcceptLanguage(acceptLanguage string) *ListPolicyAssignmentsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListPolicyAssignmentsOptions) SetTemplateID(templateID string) *ListPolicyAssignmentsOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *ListPolicyAssignmentsOptions) SetTemplateVersion(templateVersion string) *ListPolicyAssignmentsOptions {
	_options.TemplateVersion = core.StringPtr(templateVersion)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListPolicyAssignmentsOptions) SetLimit(limit int64) *ListPolicyAssignmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListPolicyAssignmentsOptions) SetStart(start string) *ListPolicyAssignmentsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPolicyAssignmentsOptions) SetHeaders(param map[string]string) *ListPolicyAssignmentsOptions {
	options.Headers = param
	return options
}

// ListPolicyTemplateVersionsOptions : The ListPolicyTemplateVersions options.
type ListPolicyTemplateVersionsOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template state.
	State *string `json:"state,omitempty"`

	// The number of documents to include in collection.
	Limit *int64 `json:"limit,omitempty"`

	// Page token that refers to the page of collection to return.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListPolicyTemplateVersionsOptions.State property.
// The policy template state.
const (
	ListPolicyTemplateVersionsOptionsStateActiveConst = "active"
	ListPolicyTemplateVersionsOptionsStateDeletedConst = "deleted"
)

// NewListPolicyTemplateVersionsOptions : Instantiate ListPolicyTemplateVersionsOptions
func (*IamPolicyManagementV1) NewListPolicyTemplateVersionsOptions(policyTemplateID string) *ListPolicyTemplateVersionsOptions {
	return &ListPolicyTemplateVersionsOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *ListPolicyTemplateVersionsOptions) SetPolicyTemplateID(policyTemplateID string) *ListPolicyTemplateVersionsOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetState : Allow user to set State
func (_options *ListPolicyTemplateVersionsOptions) SetState(state string) *ListPolicyTemplateVersionsOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListPolicyTemplateVersionsOptions) SetLimit(limit int64) *ListPolicyTemplateVersionsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListPolicyTemplateVersionsOptions) SetStart(start string) *ListPolicyTemplateVersionsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPolicyTemplateVersionsOptions) SetHeaders(param map[string]string) *ListPolicyTemplateVersionsOptions {
	options.Headers = param
	return options
}

// ListPolicyTemplatesOptions : The ListPolicyTemplates options.
type ListPolicyTemplatesOptions struct {
	// The account GUID that the policy templates belong to.
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

	// The policy template state.
	State *string `json:"state,omitempty"`

	// The policy template name.
	Name *string `json:"name,omitempty"`

	// Service type, Optional.
	PolicyServiceType *string `json:"policy_service_type,omitempty"`

	// Service name, Optional.
	PolicyServiceName *string `json:"policy_service_name,omitempty"`

	// Service group id, Optional.
	PolicyServiceGroupID *string `json:"policy_service_group_id,omitempty"`

	// Policy type, Optional.
	PolicyType *string `json:"policy_type,omitempty"`

	// The number of documents to include in collection.
	Limit *int64 `json:"limit,omitempty"`

	// Page token that refers to the page of collection to return.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListPolicyTemplatesOptions.State property.
// The policy template state.
const (
	ListPolicyTemplatesOptionsStateActiveConst = "active"
	ListPolicyTemplatesOptionsStateDeletedConst = "deleted"
)

// Constants associated with the ListPolicyTemplatesOptions.PolicyServiceType property.
// Service type, Optional.
const (
	ListPolicyTemplatesOptionsPolicyServiceTypePlatformServiceConst = "platform_service"
	ListPolicyTemplatesOptionsPolicyServiceTypeServiceConst = "service"
)

// Constants associated with the ListPolicyTemplatesOptions.PolicyType property.
// Policy type, Optional.
const (
	ListPolicyTemplatesOptionsPolicyTypeAccessConst = "access"
	ListPolicyTemplatesOptionsPolicyTypeAuthorizationConst = "authorization"
)

// NewListPolicyTemplatesOptions : Instantiate ListPolicyTemplatesOptions
func (*IamPolicyManagementV1) NewListPolicyTemplatesOptions(accountID string) *ListPolicyTemplatesOptions {
	return &ListPolicyTemplatesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListPolicyTemplatesOptions) SetAccountID(accountID string) *ListPolicyTemplatesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListPolicyTemplatesOptions) SetAcceptLanguage(acceptLanguage string) *ListPolicyTemplatesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetState : Allow user to set State
func (_options *ListPolicyTemplatesOptions) SetState(state string) *ListPolicyTemplatesOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListPolicyTemplatesOptions) SetName(name string) *ListPolicyTemplatesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPolicyServiceType : Allow user to set PolicyServiceType
func (_options *ListPolicyTemplatesOptions) SetPolicyServiceType(policyServiceType string) *ListPolicyTemplatesOptions {
	_options.PolicyServiceType = core.StringPtr(policyServiceType)
	return _options
}

// SetPolicyServiceName : Allow user to set PolicyServiceName
func (_options *ListPolicyTemplatesOptions) SetPolicyServiceName(policyServiceName string) *ListPolicyTemplatesOptions {
	_options.PolicyServiceName = core.StringPtr(policyServiceName)
	return _options
}

// SetPolicyServiceGroupID : Allow user to set PolicyServiceGroupID
func (_options *ListPolicyTemplatesOptions) SetPolicyServiceGroupID(policyServiceGroupID string) *ListPolicyTemplatesOptions {
	_options.PolicyServiceGroupID = core.StringPtr(policyServiceGroupID)
	return _options
}

// SetPolicyType : Allow user to set PolicyType
func (_options *ListPolicyTemplatesOptions) SetPolicyType(policyType string) *ListPolicyTemplatesOptions {
	_options.PolicyType = core.StringPtr(policyType)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListPolicyTemplatesOptions) SetLimit(limit int64) *ListPolicyTemplatesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListPolicyTemplatesOptions) SetStart(start string) *ListPolicyTemplatesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPolicyTemplatesOptions) SetHeaders(param map[string]string) *ListPolicyTemplatesOptions {
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

	// Allows users to set headers on API requests.
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

	// The number of documents to include in collection.
	Limit *int64 `json:"limit,omitempty"`

	// Page token that refers to the page of collection to return.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
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

// SetLimit : Allow user to set Limit
func (_options *ListV2PoliciesOptions) SetLimit(limit int64) *ListV2PoliciesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListV2PoliciesOptions) SetStart(start string) *ListV2PoliciesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListV2PoliciesOptions) SetHeaders(param map[string]string) *ListV2PoliciesOptions {
	options.Headers = param
	return options
}

// NestedCondition : Condition that specifies additional conditions or RuleAttribute to grant access.
// Models which "extend" this model:
// - NestedConditionRuleAttribute
// - NestedConditionRuleWithConditions
type NestedCondition struct {
	// The name of an attribute.
	Key *string `json:"key,omitempty"`

	// The operator of an attribute.
	Operator *string `json:"operator,omitempty"`

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value,omitempty"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []RuleAttribute `json:"conditions,omitempty"`
}

// Constants associated with the NestedCondition.Operator property.
// The operator of an attribute.
const (
	NestedConditionOperatorDategreaterthanConst = "dateGreaterThan"
	NestedConditionOperatorDategreaterthanorequalsConst = "dateGreaterThanOrEquals"
	NestedConditionOperatorDatelessthanConst = "dateLessThan"
	NestedConditionOperatorDatelessthanorequalsConst = "dateLessThanOrEquals"
	NestedConditionOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	NestedConditionOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	NestedConditionOperatorDatetimelessthanConst = "dateTimeLessThan"
	NestedConditionOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	NestedConditionOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	NestedConditionOperatorDayofweekequalsConst = "dayOfWeekEquals"
	NestedConditionOperatorStringequalsConst = "stringEquals"
	NestedConditionOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	NestedConditionOperatorStringexistsConst = "stringExists"
	NestedConditionOperatorStringmatchConst = "stringMatch"
	NestedConditionOperatorStringmatchanyofConst = "stringMatchAnyOf"
	NestedConditionOperatorTimegreaterthanConst = "timeGreaterThan"
	NestedConditionOperatorTimegreaterthanorequalsConst = "timeGreaterThanOrEquals"
	NestedConditionOperatorTimelessthanConst = "timeLessThan"
	NestedConditionOperatorTimelessthanorequalsConst = "timeLessThanOrEquals"
)
func (*NestedCondition) isaNestedCondition() bool {
	return true
}

type NestedConditionIntf interface {
	isaNestedCondition() bool
}

// UnmarshalNestedCondition unmarshals an instance of NestedCondition from the specified map of raw messages.
func UnmarshalNestedCondition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NestedCondition)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalRuleAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Next : Details with href linking to following page of requested collection.
type Next struct {
	// The href linking to the page of requested collection.
	Href *string `json:"href,omitempty"`

	// Page token that refers to the page of collection.
	Start *string `json:"start,omitempty"`
}

// UnmarshalNext unmarshals an instance of Next from the specified map of raw messages.
func UnmarshalNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Next)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
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
	err = core.UnmarshalModel(m, "subjects", &obj.Subjects, UnmarshalPolicySubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subjects-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalPolicyRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentResourcePolicy : Set of properties for the assigned resource.
type PolicyAssignmentResourcePolicy struct {
	// On success, includes the  policy assigned.
	ResourceCreated *AssignmentResourceCreated `json:"resource_created,omitempty"`

	// policy status.
	Status *string `json:"status,omitempty"`

	// The error response from API.
	ErrorMessage *ErrorResponse `json:"error_message,omitempty"`
}

// UnmarshalPolicyAssignmentResourcePolicy unmarshals an instance of PolicyAssignmentResourcePolicy from the specified map of raw messages.
func UnmarshalPolicyAssignmentResourcePolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentResourcePolicy)
	err = core.UnmarshalModel(m, "resource_created", &obj.ResourceCreated, UnmarshalAssignmentResourceCreated)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "error_message", &obj.ErrorMessage, UnmarshalErrorResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentResources : The policy assignment resources.
type PolicyAssignmentResources struct {
	// Account ID where resources are assigned.
	Target *string `json:"target,omitempty"`

	// Set of properties for the assigned resource.
	Policy *PolicyAssignmentResourcePolicy `json:"policy,omitempty"`
}

// UnmarshalPolicyAssignmentResources unmarshals an instance of PolicyAssignmentResources from the specified map of raw messages.
func UnmarshalPolicyAssignmentResources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentResources)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy", &obj.Policy, UnmarshalPolicyAssignmentResourcePolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentV1 : The set of properties associated with the policy template assignment.
type PolicyAssignmentV1 struct {
	// assignment target account and type.
	Target *AssignmentTargetDetails `json:"target" validate:"required"`

	// Policy assignment ID.
	ID *string `json:"id,omitempty"`

	// The account GUID that the policies assignments belong to..
	AccountID *string `json:"account_id,omitempty"`

	// The href URL that links to the policies assignments API by policy assignment ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy assignment was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy assignment.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy assignment was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy assignment.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// Object for each account assigned.
	Resources []PolicyAssignmentV1Resources `json:"resources" validate:"required"`

	// subject details of access type assignment.
	Subject *PolicyAssignmentV1Subject `json:"subject,omitempty"`

	// policy template details.
	Template *AssignmentTemplateDetails `json:"template" validate:"required"`

	// The policy assignment status.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the PolicyAssignmentV1.Status property.
// The policy assignment status.
const (
	PolicyAssignmentV1StatusFailedConst = "failed"
	PolicyAssignmentV1StatusInProgressConst = "in_progress"
	PolicyAssignmentV1StatusSucceedWithErrorsConst = "succeed_with_errors"
	PolicyAssignmentV1StatusSucceededConst = "succeeded"
)

// UnmarshalPolicyAssignmentV1 unmarshals an instance of PolicyAssignmentV1 from the specified map of raw messages.
func UnmarshalPolicyAssignmentV1(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentV1)
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalAssignmentTargetDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyAssignmentV1Resources)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalPolicyAssignmentV1Subject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalAssignmentTemplateDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentV1Collection : Policy assignment response.
type PolicyAssignmentV1Collection struct {
	// Response of policy assignments.
	Assignments []PolicyAssignmentV1 `json:"assignments,omitempty"`
}

// UnmarshalPolicyAssignmentV1Collection unmarshals an instance of PolicyAssignmentV1Collection from the specified map of raw messages.
func UnmarshalPolicyAssignmentV1Collection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentV1Collection)
	err = core.UnmarshalModel(m, "assignments", &obj.Assignments, UnmarshalPolicyAssignmentV1)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentV1Resources : The policy assignment resources.
type PolicyAssignmentV1Resources struct {
	// assignment target account and type.
	Target *AssignmentTargetDetails `json:"target,omitempty"`

	// Set of properties for the assigned resource.
	Policy *PolicyAssignmentResourcePolicy `json:"policy,omitempty"`
}

// UnmarshalPolicyAssignmentV1Resources unmarshals an instance of PolicyAssignmentV1Resources from the specified map of raw messages.
func UnmarshalPolicyAssignmentV1Resources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentV1Resources)
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalAssignmentTargetDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy", &obj.Policy, UnmarshalPolicyAssignmentResourcePolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyAssignmentV1Subject : subject details of access type assignment.
type PolicyAssignmentV1Subject struct {
	// The unique identifier of the subject of the assignment.
	ID *string `json:"id,omitempty"`

	// The identity type of the subject of the assignment.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the PolicyAssignmentV1Subject.Type property.
// The identity type of the subject of the assignment.
const (
	PolicyAssignmentV1SubjectTypeAccessGroupIDConst = "access_group_id"
	PolicyAssignmentV1SubjectTypeIamIDConst = "iam_id"
)

// UnmarshalPolicyAssignmentV1Subject unmarshals an instance of PolicyAssignmentV1Subject from the specified map of raw messages.
func UnmarshalPolicyAssignmentV1Subject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyAssignmentV1Subject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyCollection : A collection of policies.
type PolicyCollection struct {
	// The number of documents to include per each page of collection.
	Limit *int64 `json:"limit,omitempty"`

	// Details with href linking to first page of requested collection.
	First *First `json:"first,omitempty"`

	// Details with href linking to following page of requested collection.
	Next *Next `json:"next,omitempty"`

	// Details with href linking to previous page of requested collection.
	Previous *Previous `json:"previous,omitempty"`

	// List of policies.
	Policies []PolicyTemplateMetaData `json:"policies,omitempty"`
}

// UnmarshalPolicyCollection unmarshals an instance of PolicyCollection from the specified map of raw messages.
func UnmarshalPolicyCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalPolicyTemplateMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "policies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PolicyCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
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
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalResourceTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPolicyRole unmarshals an instance of PolicyRole from the specified map of raw messages.
func UnmarshalPolicyRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyRole)
	err = core.UnmarshalPrimitive(m, "role_id", &obj.RoleID)
	if err != nil {
		err = core.SDKErrorf(err, "", "role_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplate : The core set of properties associated with the policy template.
type PolicyTemplate struct {
	// Required field when creating a new template. Otherwise this field is optional. If the field is included it will
	// change the name value for all existing versions of the template.
	Name *string `json:"name" validate:"required"`

	// Description of the policy template. This is shown to users in the enterprise account. Use this to describe the
	// purpose or context of the policy for enterprise users managing IAM templates.
	Description *string `json:"description,omitempty"`

	// Enterprise account ID where this template will be created.
	AccountID *string `json:"account_id" validate:"required"`

	// Template version.
	Version *string `json:"version" validate:"required"`

	// Committed status of the template version.
	Committed *bool `json:"committed,omitempty"`

	// The core set of properties associated with the template's policy objet.
	Policy *TemplatePolicy `json:"policy" validate:"required"`

	// State of policy template.
	State *string `json:"state,omitempty"`

	// The policy template ID.
	ID *string `json:"id,omitempty"`

	// The href URL that links to the policy templates API by policy template ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy template was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy template.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy template was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy template.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`
}

// Constants associated with the PolicyTemplate.State property.
// State of policy template.
const (
	PolicyTemplateStateActiveConst = "active"
	PolicyTemplateStateDeletedConst = "deleted"
)

// UnmarshalPolicyTemplate unmarshals an instance of PolicyTemplate from the specified map of raw messages.
func UnmarshalPolicyTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplate)
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "committed", &obj.Committed)
	if err != nil {
		err = core.SDKErrorf(err, "", "committed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy", &obj.Policy, UnmarshalTemplatePolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateAssignmentCollection : A collection of policies assignments.
type PolicyTemplateAssignmentCollection struct {
	// The number of documents to include per each page of collection.
	Limit *int64 `json:"limit,omitempty"`

	// Details with href linking to first page of requested collection.
	First *First `json:"first,omitempty"`

	// Details with href linking to following page of requested collection.
	Next *Next `json:"next,omitempty"`

	// Details with href linking to previous page of requested collection.
	Previous *Previous `json:"previous,omitempty"`

	// List of policy assignments.
	Assignments []PolicyTemplateAssignmentItemsIntf `json:"assignments,omitempty"`
}

// UnmarshalPolicyTemplateAssignmentCollection unmarshals an instance of PolicyTemplateAssignmentCollection from the specified map of raw messages.
func UnmarshalPolicyTemplateAssignmentCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateAssignmentCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "assignments", &obj.Assignments, UnmarshalPolicyTemplateAssignmentItems)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PolicyTemplateAssignmentCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// PolicyTemplateAssignmentItems : PolicyTemplateAssignmentItems struct
// Models which "extend" this model:
// - PolicyTemplateAssignmentItemsPolicyAssignmentV1
// - PolicyTemplateAssignmentItemsPolicyAssignment
type PolicyTemplateAssignmentItems struct {
	// assignment target account and type.
	Target *AssignmentTargetDetails `json:"target,omitempty"`

	// Policy assignment ID.
	ID *string `json:"id,omitempty"`

	// The account GUID that the policies assignments belong to..
	AccountID *string `json:"account_id,omitempty"`

	// The href URL that links to the policies assignments API by policy assignment ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy assignment was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy assignment.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy assignment was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy assignment.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// Object for each account assigned.
	Resources []PolicyAssignmentV1Resources `json:"resources,omitempty"`

	// subject details of access type assignment.
	Subject *PolicyAssignmentV1Subject `json:"subject,omitempty"`

	// policy template details.
	Template *AssignmentTemplateDetails `json:"template,omitempty"`

	// The policy assignment status.
	Status *string `json:"status,omitempty"`

	// policy template id.
	TemplateID *string `json:"template_id,omitempty"`

	// policy template version.
	TemplateVersion *string `json:"template_version,omitempty"`

	// Passed in value to correlate with other assignments.
	AssignmentID *string `json:"assignment_id,omitempty"`

	// Assignment target type.
	TargetType *string `json:"target_type,omitempty"`
}

// Constants associated with the PolicyTemplateAssignmentItems.Status property.
// The policy assignment status.
const (
	PolicyTemplateAssignmentItemsStatusFailedConst = "failed"
	PolicyTemplateAssignmentItemsStatusInProgressConst = "in_progress"
	PolicyTemplateAssignmentItemsStatusSucceedWithErrorsConst = "succeed_with_errors"
	PolicyTemplateAssignmentItemsStatusSucceededConst = "succeeded"
)

// Constants associated with the PolicyTemplateAssignmentItems.TargetType property.
// Assignment target type.
const (
	PolicyTemplateAssignmentItemsTargetTypeAccountConst = "Account"
)
func (*PolicyTemplateAssignmentItems) isaPolicyTemplateAssignmentItems() bool {
	return true
}

type PolicyTemplateAssignmentItemsIntf interface {
	isaPolicyTemplateAssignmentItems() bool
}

// UnmarshalPolicyTemplateAssignmentItems unmarshals an instance of PolicyTemplateAssignmentItems from the specified map of raw messages.
func UnmarshalPolicyTemplateAssignmentItems(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateAssignmentItems)
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalAssignmentTargetDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyAssignmentV1Resources)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalPolicyAssignmentV1Subject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalAssignmentTemplateDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_version", &obj.TemplateVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "assignment_id", &obj.AssignmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateCollection : A collection of policy Templates.
type PolicyTemplateCollection struct {
	// The number of documents to include per each page of collection.
	Limit *int64 `json:"limit,omitempty"`

	// Details with href linking to first page of requested collection.
	First *First `json:"first,omitempty"`

	// Details with href linking to following page of requested collection.
	Next *Next `json:"next,omitempty"`

	// Details with href linking to previous page of requested collection.
	Previous *Previous `json:"previous,omitempty"`

	// List of policy templates.
	PolicyTemplates []PolicyTemplate `json:"policy_templates,omitempty"`
}

// UnmarshalPolicyTemplateCollection unmarshals an instance of PolicyTemplateCollection from the specified map of raw messages.
func UnmarshalPolicyTemplateCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy_templates", &obj.PolicyTemplates, UnmarshalPolicyTemplate)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_templates-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PolicyTemplateCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// PolicyTemplateLimitData : The core set of properties associated with the policy template.
type PolicyTemplateLimitData struct {
	// Required field when creating a new template. Otherwise this field is optional. If the field is included it will
	// change the name value for all existing versions of the template.
	Name *string `json:"name" validate:"required"`

	// Description of the policy template. This is shown to users in the enterprise account. Use this to describe the
	// purpose or context of the policy for enterprise users managing IAM templates.
	Description *string `json:"description,omitempty"`

	// Enterprise account ID where this template will be created.
	AccountID *string `json:"account_id" validate:"required"`

	// Template version.
	Version *string `json:"version" validate:"required"`

	// Committed status of the template version.
	Committed *bool `json:"committed,omitempty"`

	// The core set of properties associated with the template's policy objet.
	Policy *TemplatePolicy `json:"policy" validate:"required"`

	// State of policy template.
	State *string `json:"state,omitempty"`

	// The policy template ID.
	ID *string `json:"id,omitempty"`

	// The href URL that links to the policy templates API by policy template ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy template was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy template.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy template was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy template.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// policy template count details.
	Counts *TemplateCountData `json:"counts,omitempty"`
}

// Constants associated with the PolicyTemplateLimitData.State property.
// State of policy template.
const (
	PolicyTemplateLimitDataStateActiveConst = "active"
	PolicyTemplateLimitDataStateDeletedConst = "deleted"
)

// UnmarshalPolicyTemplateLimitData unmarshals an instance of PolicyTemplateLimitData from the specified map of raw messages.
func UnmarshalPolicyTemplateLimitData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateLimitData)
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "committed", &obj.Committed)
	if err != nil {
		err = core.SDKErrorf(err, "", "committed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policy", &obj.Policy, UnmarshalTemplatePolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "counts", &obj.Counts, UnmarshalTemplateCountData)
	if err != nil {
		err = core.SDKErrorf(err, "", "counts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateMetaData : The core set of properties associated with a policy.
type PolicyTemplateMetaData struct {
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

	// The details of the IAM template that was used to create an enterprise-managed policy in your account. When returned,
	// this indicates that the policy is created from and managed by a template in the root enterprise account.
	Template *TemplateMetadata `json:"template,omitempty"`
}

// Constants associated with the PolicyTemplateMetaData.State property.
// The policy state.
const (
	PolicyTemplateMetaDataStateActiveConst = "active"
	PolicyTemplateMetaDataStateDeletedConst = "deleted"
)

// UnmarshalPolicyTemplateMetaData unmarshals an instance of PolicyTemplateMetaData from the specified map of raw messages.
func UnmarshalPolicyTemplateMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateMetaData)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
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
	err = core.UnmarshalModel(m, "subjects", &obj.Subjects, UnmarshalPolicySubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subjects-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "roles", &obj.Roles, UnmarshalPolicyRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplateMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateVersionsCollection : A collection of versions for a specific policy template.
type PolicyTemplateVersionsCollection struct {
	// The number of documents to include per each page of collection.
	Limit *int64 `json:"limit,omitempty"`

	// Details with href linking to first page of requested collection.
	First *First `json:"first,omitempty"`

	// Details with href linking to following page of requested collection.
	Next *Next `json:"next,omitempty"`

	// Details with href linking to previous page of requested collection.
	Previous *Previous `json:"previous,omitempty"`

	// List of policy templates versions.
	Versions []PolicyTemplate `json:"versions,omitempty"`
}

// UnmarshalPolicyTemplateVersionsCollection unmarshals an instance of PolicyTemplateVersionsCollection from the specified map of raw messages.
func UnmarshalPolicyTemplateVersionsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateVersionsCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalPolicyTemplate)
	if err != nil {
		err = core.SDKErrorf(err, "", "versions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PolicyTemplateVersionsCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// Previous : Details with href linking to previous page of requested collection.
type Previous struct {
	// The href linking to the page of requested collection.
	Href *string `json:"href,omitempty"`

	// Page token that refers to the page of collection.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPrevious unmarshals an instance of Previous from the specified map of raw messages.
func UnmarshalPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Previous)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
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

	// Allows users to set headers on API requests.
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

// ReplacePolicyTemplateOptions : The ReplacePolicyTemplate options.
type ReplacePolicyTemplateOptions struct {
	// The policy template ID.
	PolicyTemplateID *string `json:"policy_template_id" validate:"required,ne="`

	// The policy template version.
	Version *string `json:"version" validate:"required,ne="`

	// The revision number for updating a policy template version and must match the ETag value of the existing policy
	// template version. The Etag can be retrieved using the GET
	// /v1/policy_templates/{policy_template_id}/versions/{version} API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The core set of properties associated with the template's policy objet.
	Policy *TemplatePolicy `json:"policy" validate:"required"`

	// Required field when creating a new template. Otherwise this field is optional. If the field is included it will
	// change the name value for all existing versions of the template.
	Name *string `json:"name,omitempty"`

	// Description of the policy template. This is shown to users in the enterprise account. Use this to describe the
	// purpose or context of the policy for enterprise users managing IAM templates.
	Description *string `json:"description,omitempty"`

	// Committed status of the template version.
	Committed *bool `json:"committed,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewReplacePolicyTemplateOptions : Instantiate ReplacePolicyTemplateOptions
func (*IamPolicyManagementV1) NewReplacePolicyTemplateOptions(policyTemplateID string, version string, ifMatch string, policy *TemplatePolicy) *ReplacePolicyTemplateOptions {
	return &ReplacePolicyTemplateOptions{
		PolicyTemplateID: core.StringPtr(policyTemplateID),
		Version: core.StringPtr(version),
		IfMatch: core.StringPtr(ifMatch),
		Policy: policy,
	}
}

// SetPolicyTemplateID : Allow user to set PolicyTemplateID
func (_options *ReplacePolicyTemplateOptions) SetPolicyTemplateID(policyTemplateID string) *ReplacePolicyTemplateOptions {
	_options.PolicyTemplateID = core.StringPtr(policyTemplateID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ReplacePolicyTemplateOptions) SetVersion(version string) *ReplacePolicyTemplateOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplacePolicyTemplateOptions) SetIfMatch(ifMatch string) *ReplacePolicyTemplateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetPolicy : Allow user to set Policy
func (_options *ReplacePolicyTemplateOptions) SetPolicy(policy *TemplatePolicy) *ReplacePolicyTemplateOptions {
	_options.Policy = policy
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplacePolicyTemplateOptions) SetName(name string) *ReplacePolicyTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplacePolicyTemplateOptions) SetDescription(description string) *ReplacePolicyTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetCommitted : Allow user to set Committed
func (_options *ReplacePolicyTemplateOptions) SetCommitted(committed bool) *ReplacePolicyTemplateOptions {
	_options.Committed = core.BoolPtr(committed)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplacePolicyTemplateOptions) SetHeaders(param map[string]string) *ReplacePolicyTemplateOptions {
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

	// Allows users to set headers on API requests.
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

	// Description of the policy.
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

	// Allows users to set headers on API requests.
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResourceAttribute unmarshals an instance of ResourceAttribute from the specified map of raw messages.
func UnmarshalResourceAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceAttribute)
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
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResourceTag unmarshals an instance of ResourceTag from the specified map of raw messages.
func UnmarshalResourceTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceTag)
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
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRole unmarshals an instance of Role from the specified map of raw messages.
func UnmarshalRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Role)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "actions", &obj.Actions)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
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

// RoleCollection : A collection of roles returned by the 'list roles' operation.
type RoleCollection struct {
	// List of custom roles.
	CustomRoles []CustomRole `json:"custom_roles,omitempty"`

	// List of service roles.
	ServiceRoles []Role `json:"service_roles,omitempty"`

	// List of system roles.
	SystemRoles []Role `json:"system_roles,omitempty"`
}

// UnmarshalRoleCollection unmarshals an instance of RoleCollection from the specified map of raw messages.
func UnmarshalRoleCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleCollection)
	err = core.UnmarshalModel(m, "custom_roles", &obj.CustomRoles, UnmarshalCustomRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "custom_roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_roles", &obj.ServiceRoles, UnmarshalRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "system_roles", &obj.SystemRoles, UnmarshalRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_roles-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRoles unmarshals an instance of Roles from the specified map of raw messages.
func UnmarshalRoles(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Roles)
	err = core.UnmarshalPrimitive(m, "role_id", &obj.RoleID)
	if err != nil {
		err = core.SDKErrorf(err, "", "role_id-error", common.GetComponentInfo())
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

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the RuleAttribute.Operator property.
// The operator of an attribute.
const (
	RuleAttributeOperatorDategreaterthanConst = "dateGreaterThan"
	RuleAttributeOperatorDategreaterthanorequalsConst = "dateGreaterThanOrEquals"
	RuleAttributeOperatorDatelessthanConst = "dateLessThan"
	RuleAttributeOperatorDatelessthanorequalsConst = "dateLessThanOrEquals"
	RuleAttributeOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	RuleAttributeOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	RuleAttributeOperatorDatetimelessthanConst = "dateTimeLessThan"
	RuleAttributeOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	RuleAttributeOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	RuleAttributeOperatorDayofweekequalsConst = "dayOfWeekEquals"
	RuleAttributeOperatorStringequalsConst = "stringEquals"
	RuleAttributeOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	RuleAttributeOperatorStringexistsConst = "stringExists"
	RuleAttributeOperatorStringmatchConst = "stringMatch"
	RuleAttributeOperatorStringmatchanyofConst = "stringMatchAnyOf"
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRuleAttribute unmarshals an instance of RuleAttribute from the specified map of raw messages.
func UnmarshalRuleAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalSubjectAttribute unmarshals an instance of SubjectAttribute from the specified map of raw messages.
func UnmarshalSubjectAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubjectAttribute)
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

// TemplateCountData : policy template count details.
type TemplateCountData struct {
	// policy template current and limit details with in an account.
	Template *LimitData `json:"template,omitempty"`

	// policy template current and limit details with in an account.
	Version *LimitData `json:"version,omitempty"`
}

// UnmarshalTemplateCountData unmarshals an instance of TemplateCountData from the specified map of raw messages.
func UnmarshalTemplateCountData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateCountData)
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalLimitData)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "version", &obj.Version, UnmarshalLimitData)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateMetadata : The details of the IAM template that was used to create an enterprise-managed policy in your account. When returned,
// this indicates that the policy is created from and managed by a template in the root enterprise account.
type TemplateMetadata struct {
	// The policy template ID.
	ID *string `json:"id,omitempty"`

	// Template version.
	Version *string `json:"version,omitempty"`

	// policy assignment id.
	AssignmentID *string `json:"assignment_id,omitempty"`

	// orchestrator template id.
	RootID *string `json:"root_id,omitempty"`

	// orchestrator template version.
	RootVersion *string `json:"root_version,omitempty"`
}

// UnmarshalTemplateMetadata unmarshals an instance of TemplateMetadata from the specified map of raw messages.
func UnmarshalTemplateMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateMetadata)
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
	err = core.UnmarshalPrimitive(m, "assignment_id", &obj.AssignmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "root_id", &obj.RootID)
	if err != nil {
		err = core.SDKErrorf(err, "", "root_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "root_version", &obj.RootVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "root_version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplatePolicy : The core set of properties associated with the template's policy objet.
type TemplatePolicy struct {
	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Description of the policy. This is shown in child accounts when an access group or trusted profile template uses the
	// policy template to assign access.
	Description *string `json:"description,omitempty"`

	// The resource attributes to which the policy grants access.
	Resource *V2PolicyResource `json:"resource,omitempty"`

	// The subject attributes for whom the policy grants access.
	Subject *V2PolicySubject `json:"subject,omitempty"`

	// Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or
	// 'time-based-conditions:weekly:custom-hours'.
	Pattern *string `json:"pattern,omitempty"`

	// Additional access conditions associated with the policy.
	Rule V2PolicyRuleIntf `json:"rule,omitempty"`

	// Specifies the type of access granted by the policy.
	Control *Control `json:"control,omitempty"`
}

// Constants associated with the TemplatePolicy.Type property.
// The policy type; either 'access' or 'authorization'.
const (
	TemplatePolicyTypeAccessConst = "access"
	TemplatePolicyTypeAuthorizationConst = "authorization"
)

// NewTemplatePolicy : Instantiate TemplatePolicy (Generic Model Constructor)
func (*IamPolicyManagementV1) NewTemplatePolicy(typeVar string) (_model *TemplatePolicy, err error) {
	_model = &TemplatePolicy{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTemplatePolicy unmarshals an instance of TemplatePolicy from the specified map of raw messages.
func UnmarshalTemplatePolicy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplatePolicy)
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
	err = core.UnmarshalModel(m, "resource", &obj.Resource, UnmarshalV2PolicyResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalV2PolicySubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		err = core.SDKErrorf(err, "", "pattern-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalV2PolicyRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rule-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "control", &obj.Control, UnmarshalControl)
	if err != nil {
		err = core.SDKErrorf(err, "", "control-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatePolicyAssignmentOptions : The UpdatePolicyAssignment options.
type UpdatePolicyAssignmentOptions struct {
	// The policy template assignment ID.
	AssignmentID *string `json:"assignment_id" validate:"required,ne="`

	// specify version of response body format.
	Version *string `json:"version" validate:"required"`

	// The revision number for updating a policy assignment and must match the ETag value of the existing policy
	// assignment. The Etag can be retrieved using the GET /v1/policy_assignments/{assignment_id} API and looking at the
	// ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The policy template version to update to.
	TemplateVersion *string `json:"template_version" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdatePolicyAssignmentOptions : Instantiate UpdatePolicyAssignmentOptions
func (*IamPolicyManagementV1) NewUpdatePolicyAssignmentOptions(assignmentID string, version string, ifMatch string, templateVersion string) *UpdatePolicyAssignmentOptions {
	return &UpdatePolicyAssignmentOptions{
		AssignmentID: core.StringPtr(assignmentID),
		Version: core.StringPtr(version),
		IfMatch: core.StringPtr(ifMatch),
		TemplateVersion: core.StringPtr(templateVersion),
	}
}

// SetAssignmentID : Allow user to set AssignmentID
func (_options *UpdatePolicyAssignmentOptions) SetAssignmentID(assignmentID string) *UpdatePolicyAssignmentOptions {
	_options.AssignmentID = core.StringPtr(assignmentID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *UpdatePolicyAssignmentOptions) SetVersion(version string) *UpdatePolicyAssignmentOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdatePolicyAssignmentOptions) SetIfMatch(ifMatch string) *UpdatePolicyAssignmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetTemplateVersion : Allow user to set TemplateVersion
func (_options *UpdatePolicyAssignmentOptions) SetTemplateVersion(templateVersion string) *UpdatePolicyAssignmentOptions {
	_options.TemplateVersion = core.StringPtr(templateVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePolicyAssignmentOptions) SetHeaders(param map[string]string) *UpdatePolicyAssignmentOptions {
	options.Headers = param
	return options
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

	// Allows users to set headers on API requests.
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

// UpdateSettingsOptions : The UpdateSettings options.
type UpdateSettingsOptions struct {
	// The account GUID that the settings belong to.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The revision number for updating Access Management Account Settings and must match the ETag value of the existing
	// Access Management Account Settings. The Etag can be retrieved using the GET
	// /v1/accounts/{account_id}/settings/access_management API and looking at the ETag response header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Update to how external accounts can interact in relation to the requested account.
	ExternalAccountIdentityInteraction *ExternalAccountIdentityInteractionPatch `json:"external_account_identity_interaction,omitempty"`

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

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateSettingsOptions : Instantiate UpdateSettingsOptions
func (*IamPolicyManagementV1) NewUpdateSettingsOptions(accountID string, ifMatch string) *UpdateSettingsOptions {
	return &UpdateSettingsOptions{
		AccountID: core.StringPtr(accountID),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateSettingsOptions) SetAccountID(accountID string) *UpdateSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateSettingsOptions) SetIfMatch(ifMatch string) *UpdateSettingsOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetExternalAccountIdentityInteraction : Allow user to set ExternalAccountIdentityInteraction
func (_options *UpdateSettingsOptions) SetExternalAccountIdentityInteraction(externalAccountIdentityInteraction *ExternalAccountIdentityInteractionPatch) *UpdateSettingsOptions {
	_options.ExternalAccountIdentityInteraction = externalAccountIdentityInteraction
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *UpdateSettingsOptions) SetAcceptLanguage(acceptLanguage string) *UpdateSettingsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSettingsOptions) SetHeaders(param map[string]string) *UpdateSettingsOptions {
	options.Headers = param
	return options
}

// V2Policy : The core set of properties associated with the policy.
type V2Policy struct {
	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Description of the policy.
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
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalV2PolicySubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource", &obj.Resource, UnmarshalV2PolicyResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		err = core.SDKErrorf(err, "", "pattern-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalV2PolicyRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rule-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalModel(m, "control", &obj.Control, UnmarshalControlResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "control-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_at", &obj.LastPermitAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_permit_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_frequency", &obj.LastPermitFrequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_permit_frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyCollection : A collection of policies.
type V2PolicyCollection struct {
	// The number of documents to include per each page of collection.
	Limit *int64 `json:"limit,omitempty"`

	// Details with href linking to first page of requested collection.
	First *First `json:"first,omitempty"`

	// Details with href linking to following page of requested collection.
	Next *Next `json:"next,omitempty"`

	// Details with href linking to previous page of requested collection.
	Previous *Previous `json:"previous,omitempty"`

	// List of policies.
	Policies []V2PolicyTemplateMetaData `json:"policies,omitempty"`
}

// UnmarshalV2PolicyCollection unmarshals an instance of V2PolicyCollection from the specified map of raw messages.
func UnmarshalV2PolicyCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalV2PolicyTemplateMetaData)
	if err != nil {
		err = core.SDKErrorf(err, "", "policies-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *V2PolicyCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalV2PolicyResource unmarshals an instance of V2PolicyResource from the specified map of raw messages.
func UnmarshalV2PolicyResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResource)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalV2PolicyResourceAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalV2PolicyResourceTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
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

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the V2PolicyResourceAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicyResourceAttributeOperatorStringequalsConst = "stringEquals"
	V2PolicyResourceAttributeOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	V2PolicyResourceAttributeOperatorStringexistsConst = "stringExists"
	V2PolicyResourceAttributeOperatorStringmatchConst = "stringMatch"
	V2PolicyResourceAttributeOperatorStringmatchanyofConst = "stringMatchAnyOf"
)

// NewV2PolicyResourceAttribute : Instantiate V2PolicyResourceAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyResourceAttribute(key string, operator string, value interface{}) (_model *V2PolicyResourceAttribute, err error) {
	_model = &V2PolicyResourceAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalV2PolicyResourceAttribute unmarshals an instance of V2PolicyResourceAttribute from the specified map of raw messages.
func UnmarshalV2PolicyResourceAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResourceAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalV2PolicyResourceTag unmarshals an instance of V2PolicyResourceTag from the specified map of raw messages.
func UnmarshalV2PolicyResourceTag(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyResourceTag)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// V2PolicyRule : Additional access conditions associated with the policy.
// Models which "extend" this model:
// - V2PolicyRuleRuleAttribute
// - V2PolicyRuleRuleWithNestedConditions
type V2PolicyRule struct {
	// The name of an attribute.
	Key *string `json:"key,omitempty"`

	// The operator of an attribute.
	Operator *string `json:"operator,omitempty"`

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value,omitempty"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []NestedConditionIntf `json:"conditions,omitempty"`
}

// Constants associated with the V2PolicyRule.Operator property.
// The operator of an attribute.
const (
	V2PolicyRuleOperatorDategreaterthanConst = "dateGreaterThan"
	V2PolicyRuleOperatorDategreaterthanorequalsConst = "dateGreaterThanOrEquals"
	V2PolicyRuleOperatorDatelessthanConst = "dateLessThan"
	V2PolicyRuleOperatorDatelessthanorequalsConst = "dateLessThanOrEquals"
	V2PolicyRuleOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	V2PolicyRuleOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	V2PolicyRuleOperatorDatetimelessthanConst = "dateTimeLessThan"
	V2PolicyRuleOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	V2PolicyRuleOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	V2PolicyRuleOperatorDayofweekequalsConst = "dayOfWeekEquals"
	V2PolicyRuleOperatorStringequalsConst = "stringEquals"
	V2PolicyRuleOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	V2PolicyRuleOperatorStringexistsConst = "stringExists"
	V2PolicyRuleOperatorStringmatchConst = "stringMatch"
	V2PolicyRuleOperatorStringmatchanyofConst = "stringMatchAnyOf"
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
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalNestedCondition)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalV2PolicySubject unmarshals an instance of V2PolicySubject from the specified map of raw messages.
func UnmarshalV2PolicySubject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicySubject)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalV2PolicySubjectAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
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

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the V2PolicySubjectAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicySubjectAttributeOperatorStringequalsConst = "stringEquals"
	V2PolicySubjectAttributeOperatorStringexistsConst = "stringExists"
)

// NewV2PolicySubjectAttribute : Instantiate V2PolicySubjectAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicySubjectAttribute(key string, operator string, value interface{}) (_model *V2PolicySubjectAttribute, err error) {
	_model = &V2PolicySubjectAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalV2PolicySubjectAttribute unmarshals an instance of V2PolicySubjectAttribute from the specified map of raw messages.
func UnmarshalV2PolicySubjectAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicySubjectAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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

// V2PolicyTemplateMetaData : The core set of properties associated with the policy.
type V2PolicyTemplateMetaData struct {
	// The policy type; either 'access' or 'authorization'.
	Type *string `json:"type" validate:"required"`

	// Description of the policy.
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

	// The details of the IAM template that was used to create an enterprise-managed policy in your account. When returned,
	// this indicates that the policy is created from and managed by a template in the root enterprise account.
	Template *TemplateMetadata `json:"template,omitempty"`
}

// Constants associated with the V2PolicyTemplateMetaData.Type property.
// The policy type; either 'access' or 'authorization'.
const (
	V2PolicyTemplateMetaDataTypeAccessConst = "access"
	V2PolicyTemplateMetaDataTypeAuthorizationConst = "authorization"
)

// Constants associated with the V2PolicyTemplateMetaData.State property.
// The policy state, either 'deleted' or 'active'.
const (
	V2PolicyTemplateMetaDataStateActiveConst = "active"
	V2PolicyTemplateMetaDataStateDeletedConst = "deleted"
)

// UnmarshalV2PolicyTemplateMetaData unmarshals an instance of V2PolicyTemplateMetaData from the specified map of raw messages.
func UnmarshalV2PolicyTemplateMetaData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyTemplateMetaData)
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
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalV2PolicySubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource", &obj.Resource, UnmarshalV2PolicyResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		err = core.SDKErrorf(err, "", "pattern-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalV2PolicyRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rule-error", common.GetComponentInfo())
		return
	}
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
	err = core.UnmarshalModel(m, "control", &obj.Control, UnmarshalControlResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "control-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_at", &obj.LastPermitAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_permit_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_permit_frequency", &obj.LastPermitFrequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_permit_frequency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplateMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "grant-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "grant-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NestedConditionRuleAttribute : Rule that specifies additional access granted (e.g., time-based condition).
// This model "extends" NestedCondition
type NestedConditionRuleAttribute struct {
	// The name of an attribute.
	Key *string `json:"key" validate:"required"`

	// The operator of an attribute.
	Operator *string `json:"operator" validate:"required"`

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the NestedConditionRuleAttribute.Operator property.
// The operator of an attribute.
const (
	NestedConditionRuleAttributeOperatorDategreaterthanConst = "dateGreaterThan"
	NestedConditionRuleAttributeOperatorDategreaterthanorequalsConst = "dateGreaterThanOrEquals"
	NestedConditionRuleAttributeOperatorDatelessthanConst = "dateLessThan"
	NestedConditionRuleAttributeOperatorDatelessthanorequalsConst = "dateLessThanOrEquals"
	NestedConditionRuleAttributeOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	NestedConditionRuleAttributeOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	NestedConditionRuleAttributeOperatorDatetimelessthanConst = "dateTimeLessThan"
	NestedConditionRuleAttributeOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	NestedConditionRuleAttributeOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	NestedConditionRuleAttributeOperatorDayofweekequalsConst = "dayOfWeekEquals"
	NestedConditionRuleAttributeOperatorStringequalsConst = "stringEquals"
	NestedConditionRuleAttributeOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	NestedConditionRuleAttributeOperatorStringexistsConst = "stringExists"
	NestedConditionRuleAttributeOperatorStringmatchConst = "stringMatch"
	NestedConditionRuleAttributeOperatorStringmatchanyofConst = "stringMatchAnyOf"
	NestedConditionRuleAttributeOperatorTimegreaterthanConst = "timeGreaterThan"
	NestedConditionRuleAttributeOperatorTimegreaterthanorequalsConst = "timeGreaterThanOrEquals"
	NestedConditionRuleAttributeOperatorTimelessthanConst = "timeLessThan"
	NestedConditionRuleAttributeOperatorTimelessthanorequalsConst = "timeLessThanOrEquals"
)

// NewNestedConditionRuleAttribute : Instantiate NestedConditionRuleAttribute (Generic Model Constructor)
func (*IamPolicyManagementV1) NewNestedConditionRuleAttribute(key string, operator string, value interface{}) (_model *NestedConditionRuleAttribute, err error) {
	_model = &NestedConditionRuleAttribute{
		Key: core.StringPtr(key),
		Operator: core.StringPtr(operator),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*NestedConditionRuleAttribute) isaNestedCondition() bool {
	return true
}

// UnmarshalNestedConditionRuleAttribute unmarshals an instance of NestedConditionRuleAttribute from the specified map of raw messages.
func UnmarshalNestedConditionRuleAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NestedConditionRuleAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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

// NestedConditionRuleWithConditions : Rule that specifies additional access granted (e.g., time-based condition) accross multiple conditions.
// This model "extends" NestedCondition
type NestedConditionRuleWithConditions struct {
	// Operator to evaluate conditions.
	Operator *string `json:"operator" validate:"required"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []RuleAttribute `json:"conditions" validate:"required"`
}

// Constants associated with the NestedConditionRuleWithConditions.Operator property.
// Operator to evaluate conditions.
const (
	NestedConditionRuleWithConditionsOperatorAndConst = "and"
	NestedConditionRuleWithConditionsOperatorOrConst = "or"
)

// NewNestedConditionRuleWithConditions : Instantiate NestedConditionRuleWithConditions (Generic Model Constructor)
func (*IamPolicyManagementV1) NewNestedConditionRuleWithConditions(operator string, conditions []RuleAttribute) (_model *NestedConditionRuleWithConditions, err error) {
	_model = &NestedConditionRuleWithConditions{
		Operator: core.StringPtr(operator),
		Conditions: conditions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*NestedConditionRuleWithConditions) isaNestedCondition() bool {
	return true
}

// UnmarshalNestedConditionRuleWithConditions unmarshals an instance of NestedConditionRuleWithConditions from the specified map of raw messages.
func UnmarshalNestedConditionRuleWithConditions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NestedConditionRuleWithConditions)
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalRuleAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateAssignmentItemsPolicyAssignment : The set of properties associated with the policy template assignment.
// This model "extends" PolicyTemplateAssignmentItems
type PolicyTemplateAssignmentItemsPolicyAssignment struct {
	// policy template id.
	TemplateID *string `json:"template_id,omitempty"`

	// policy template version.
	TemplateVersion *string `json:"template_version,omitempty"`

	// Passed in value to correlate with other assignments.
	AssignmentID *string `json:"assignment_id,omitempty"`

	// Assignment target type.
	TargetType *string `json:"target_type,omitempty"`

	// ID of the target account.
	Target *string `json:"target,omitempty"`

	// Policy assignment ID.
	ID *string `json:"id,omitempty"`

	// The account GUID that the policies assignments belong to..
	AccountID *string `json:"account_id,omitempty"`

	// The href URL that links to the policies assignments API by policy assignment ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy assignment was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy assignment.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy assignment was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy assignment.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// Object for each account assigned.
	Resources []PolicyAssignmentResources `json:"resources,omitempty"`

	// The policy assignment status.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the PolicyTemplateAssignmentItemsPolicyAssignment.TargetType property.
// Assignment target type.
const (
	PolicyTemplateAssignmentItemsPolicyAssignmentTargetTypeAccountConst = "Account"
)

// Constants associated with the PolicyTemplateAssignmentItemsPolicyAssignment.Status property.
// The policy assignment status.
const (
	PolicyTemplateAssignmentItemsPolicyAssignmentStatusFailedConst = "failed"
	PolicyTemplateAssignmentItemsPolicyAssignmentStatusInProgressConst = "in_progress"
	PolicyTemplateAssignmentItemsPolicyAssignmentStatusSucceedWithErrorsConst = "succeed_with_errors"
	PolicyTemplateAssignmentItemsPolicyAssignmentStatusSucceededConst = "succeeded"
)

func (*PolicyTemplateAssignmentItemsPolicyAssignment) isaPolicyTemplateAssignmentItems() bool {
	return true
}

// UnmarshalPolicyTemplateAssignmentItemsPolicyAssignment unmarshals an instance of PolicyTemplateAssignmentItemsPolicyAssignment from the specified map of raw messages.
func UnmarshalPolicyTemplateAssignmentItemsPolicyAssignment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateAssignmentItemsPolicyAssignment)
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "template_version", &obj.TemplateVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "template_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "assignment_id", &obj.AssignmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "assignment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyAssignmentResources)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyTemplateAssignmentItemsPolicyAssignmentV1 : The set of properties associated with the policy template assignment.
// This model "extends" PolicyTemplateAssignmentItems
type PolicyTemplateAssignmentItemsPolicyAssignmentV1 struct {
	// assignment target account and type.
	Target *AssignmentTargetDetails `json:"target" validate:"required"`

	// Policy assignment ID.
	ID *string `json:"id,omitempty"`

	// The account GUID that the policies assignments belong to..
	AccountID *string `json:"account_id,omitempty"`

	// The href URL that links to the policies assignments API by policy assignment ID.
	Href *string `json:"href,omitempty"`

	// The UTC timestamp when the policy assignment was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The iam ID of the entity that created the policy assignment.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The UTC timestamp when the policy assignment was last modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// The iam ID of the entity that last modified the policy assignment.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`

	// Object for each account assigned.
	Resources []PolicyAssignmentV1Resources `json:"resources" validate:"required"`

	// subject details of access type assignment.
	Subject *PolicyAssignmentV1Subject `json:"subject,omitempty"`

	// policy template details.
	Template *AssignmentTemplateDetails `json:"template" validate:"required"`

	// The policy assignment status.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the PolicyTemplateAssignmentItemsPolicyAssignmentV1.Status property.
// The policy assignment status.
const (
	PolicyTemplateAssignmentItemsPolicyAssignmentV1StatusFailedConst = "failed"
	PolicyTemplateAssignmentItemsPolicyAssignmentV1StatusInProgressConst = "in_progress"
	PolicyTemplateAssignmentItemsPolicyAssignmentV1StatusSucceedWithErrorsConst = "succeed_with_errors"
	PolicyTemplateAssignmentItemsPolicyAssignmentV1StatusSucceededConst = "succeeded"
)

func (*PolicyTemplateAssignmentItemsPolicyAssignmentV1) isaPolicyTemplateAssignmentItems() bool {
	return true
}

// UnmarshalPolicyTemplateAssignmentItemsPolicyAssignmentV1 unmarshals an instance of PolicyTemplateAssignmentItemsPolicyAssignmentV1 from the specified map of raw messages.
func UnmarshalPolicyTemplateAssignmentItemsPolicyAssignmentV1(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyTemplateAssignmentItemsPolicyAssignmentV1)
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalAssignmentTargetDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalPolicyAssignmentV1Resources)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "subject", &obj.Subject, UnmarshalPolicyAssignmentV1Subject)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalAssignmentTemplateDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "template-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
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

	// The value of a rule, resource, or subject attribute; can be boolean or string for resource and subject attribute.
	// Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.
	Value interface{} `json:"value" validate:"required"`
}

// Constants associated with the V2PolicyRuleRuleAttribute.Operator property.
// The operator of an attribute.
const (
	V2PolicyRuleRuleAttributeOperatorDategreaterthanConst = "dateGreaterThan"
	V2PolicyRuleRuleAttributeOperatorDategreaterthanorequalsConst = "dateGreaterThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDatelessthanConst = "dateLessThan"
	V2PolicyRuleRuleAttributeOperatorDatelessthanorequalsConst = "dateLessThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDatetimegreaterthanConst = "dateTimeGreaterThan"
	V2PolicyRuleRuleAttributeOperatorDatetimegreaterthanorequalsConst = "dateTimeGreaterThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDatetimelessthanConst = "dateTimeLessThan"
	V2PolicyRuleRuleAttributeOperatorDatetimelessthanorequalsConst = "dateTimeLessThanOrEquals"
	V2PolicyRuleRuleAttributeOperatorDayofweekanyofConst = "dayOfWeekAnyOf"
	V2PolicyRuleRuleAttributeOperatorDayofweekequalsConst = "dayOfWeekEquals"
	V2PolicyRuleRuleAttributeOperatorStringequalsConst = "stringEquals"
	V2PolicyRuleRuleAttributeOperatorStringequalsanyofConst = "stringEqualsAnyOf"
	V2PolicyRuleRuleAttributeOperatorStringexistsConst = "stringExists"
	V2PolicyRuleRuleAttributeOperatorStringmatchConst = "stringMatch"
	V2PolicyRuleRuleAttributeOperatorStringmatchanyofConst = "stringMatchAnyOf"
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
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
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
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
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

// V2PolicyRuleRuleWithNestedConditions : Rule that specifies additional access granted (e.g., time-based condition) accross multiple conditions.
// This model "extends" V2PolicyRule
type V2PolicyRuleRuleWithNestedConditions struct {
	// Operator to evaluate conditions.
	Operator *string `json:"operator" validate:"required"`

	// List of conditions associated with a policy, e.g., time-based conditions that grant access over a certain time
	// period.
	Conditions []NestedConditionIntf `json:"conditions" validate:"required"`
}

// Constants associated with the V2PolicyRuleRuleWithNestedConditions.Operator property.
// Operator to evaluate conditions.
const (
	V2PolicyRuleRuleWithNestedConditionsOperatorAndConst = "and"
	V2PolicyRuleRuleWithNestedConditionsOperatorOrConst = "or"
)

// NewV2PolicyRuleRuleWithNestedConditions : Instantiate V2PolicyRuleRuleWithNestedConditions (Generic Model Constructor)
func (*IamPolicyManagementV1) NewV2PolicyRuleRuleWithNestedConditions(operator string, conditions []NestedConditionIntf) (_model *V2PolicyRuleRuleWithNestedConditions, err error) {
	_model = &V2PolicyRuleRuleWithNestedConditions{
		Operator: core.StringPtr(operator),
		Conditions: conditions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*V2PolicyRuleRuleWithNestedConditions) isaV2PolicyRule() bool {
	return true
}

// UnmarshalV2PolicyRuleRuleWithNestedConditions unmarshals an instance of V2PolicyRuleRuleWithNestedConditions from the specified map of raw messages.
func UnmarshalV2PolicyRuleRuleWithNestedConditions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(V2PolicyRuleRuleWithNestedConditions)
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "conditions", &obj.Conditions, UnmarshalNestedCondition)
	if err != nil {
		err = core.SDKErrorf(err, "", "conditions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

//
// PoliciesPager can be used to simplify the use of the "ListPolicies" method.
//
type PoliciesPager struct {
	hasNext bool
	options *ListPoliciesOptions
	client  *IamPolicyManagementV1
	pageContext struct {
		next *string
	}
}

// NewPoliciesPager returns a new PoliciesPager instance.
func (iamPolicyManagement *IamPolicyManagementV1) NewPoliciesPager(options *ListPoliciesOptions) (pager *PoliciesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPoliciesOptions = *options
	pager = &PoliciesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  iamPolicyManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PoliciesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PoliciesPager) GetNextWithContext(ctx context.Context) (page []PolicyTemplateMetaData, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListPoliciesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Policies

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PoliciesPager) GetAllWithContext(ctx context.Context) (allItems []PolicyTemplateMetaData, err error) {
	for pager.HasNext() {
		var nextPage []PolicyTemplateMetaData
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
func (pager *PoliciesPager) GetNext() (page []PolicyTemplateMetaData, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PoliciesPager) GetAll() (allItems []PolicyTemplateMetaData, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// V2PoliciesPager can be used to simplify the use of the "ListV2Policies" method.
//
type V2PoliciesPager struct {
	hasNext bool
	options *ListV2PoliciesOptions
	client  *IamPolicyManagementV1
	pageContext struct {
		next *string
	}
}

// NewV2PoliciesPager returns a new V2PoliciesPager instance.
func (iamPolicyManagement *IamPolicyManagementV1) NewV2PoliciesPager(options *ListV2PoliciesOptions) (pager *V2PoliciesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListV2PoliciesOptions = *options
	pager = &V2PoliciesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  iamPolicyManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *V2PoliciesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *V2PoliciesPager) GetNextWithContext(ctx context.Context) (page []V2PolicyTemplateMetaData, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListV2PoliciesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Policies

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *V2PoliciesPager) GetAllWithContext(ctx context.Context) (allItems []V2PolicyTemplateMetaData, err error) {
	for pager.HasNext() {
		var nextPage []V2PolicyTemplateMetaData
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
func (pager *V2PoliciesPager) GetNext() (page []V2PolicyTemplateMetaData, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *V2PoliciesPager) GetAll() (allItems []V2PolicyTemplateMetaData, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// PolicyTemplatesPager can be used to simplify the use of the "ListPolicyTemplates" method.
//
type PolicyTemplatesPager struct {
	hasNext bool
	options *ListPolicyTemplatesOptions
	client  *IamPolicyManagementV1
	pageContext struct {
		next *string
	}
}

// NewPolicyTemplatesPager returns a new PolicyTemplatesPager instance.
func (iamPolicyManagement *IamPolicyManagementV1) NewPolicyTemplatesPager(options *ListPolicyTemplatesOptions) (pager *PolicyTemplatesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPolicyTemplatesOptions = *options
	pager = &PolicyTemplatesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  iamPolicyManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PolicyTemplatesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PolicyTemplatesPager) GetNextWithContext(ctx context.Context) (page []PolicyTemplate, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListPolicyTemplatesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.PolicyTemplates

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PolicyTemplatesPager) GetAllWithContext(ctx context.Context) (allItems []PolicyTemplate, err error) {
	for pager.HasNext() {
		var nextPage []PolicyTemplate
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
func (pager *PolicyTemplatesPager) GetNext() (page []PolicyTemplate, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PolicyTemplatesPager) GetAll() (allItems []PolicyTemplate, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// PolicyTemplateVersionsPager can be used to simplify the use of the "ListPolicyTemplateVersions" method.
//
type PolicyTemplateVersionsPager struct {
	hasNext bool
	options *ListPolicyTemplateVersionsOptions
	client  *IamPolicyManagementV1
	pageContext struct {
		next *string
	}
}

// NewPolicyTemplateVersionsPager returns a new PolicyTemplateVersionsPager instance.
func (iamPolicyManagement *IamPolicyManagementV1) NewPolicyTemplateVersionsPager(options *ListPolicyTemplateVersionsOptions) (pager *PolicyTemplateVersionsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPolicyTemplateVersionsOptions = *options
	pager = &PolicyTemplateVersionsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  iamPolicyManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PolicyTemplateVersionsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PolicyTemplateVersionsPager) GetNextWithContext(ctx context.Context) (page []PolicyTemplate, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListPolicyTemplateVersionsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Versions

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PolicyTemplateVersionsPager) GetAllWithContext(ctx context.Context) (allItems []PolicyTemplate, err error) {
	for pager.HasNext() {
		var nextPage []PolicyTemplate
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
func (pager *PolicyTemplateVersionsPager) GetNext() (page []PolicyTemplate, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PolicyTemplateVersionsPager) GetAll() (allItems []PolicyTemplate, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// PolicyAssignmentsPager can be used to simplify the use of the "ListPolicyAssignments" method.
//
type PolicyAssignmentsPager struct {
	hasNext bool
	options *ListPolicyAssignmentsOptions
	client  *IamPolicyManagementV1
	pageContext struct {
		next *string
	}
}

// NewPolicyAssignmentsPager returns a new PolicyAssignmentsPager instance.
func (iamPolicyManagement *IamPolicyManagementV1) NewPolicyAssignmentsPager(options *ListPolicyAssignmentsOptions) (pager *PolicyAssignmentsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPolicyAssignmentsOptions = *options
	pager = &PolicyAssignmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  iamPolicyManagement,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PolicyAssignmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PolicyAssignmentsPager) GetNextWithContext(ctx context.Context) (page []PolicyTemplateAssignmentItemsIntf, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListPolicyAssignmentsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Assignments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PolicyAssignmentsPager) GetAllWithContext(ctx context.Context) (allItems []PolicyTemplateAssignmentItemsIntf, err error) {
	for pager.HasNext() {
		var nextPage []PolicyTemplateAssignmentItemsIntf
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
func (pager *PolicyAssignmentsPager) GetNext() (page []PolicyTemplateAssignmentItemsIntf, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PolicyAssignmentsPager) GetAll() (allItems []PolicyTemplateAssignmentItemsIntf, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
