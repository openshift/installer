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
 * IBM OpenAPI SDK Code Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

// Package db2saasv1 : Operations and models for the Db2saasV1 service
package db2saasv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/cloud-db2-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
)

// Db2saasV1 : Manage lifecycle of your Db2 on Cloud resources using the  APIs.
//
// API Version: 1.0.0
type Db2saasV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.db2.saas.ibm.com/dbapi/v4"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "db2saas"

const ParameterizedServiceURL = "https://{region}.db2.saas.ibm.com/dbapi/v4"

var defaultUrlVariables = map[string]string{
	"region": "us-south",
}

// Db2saasV1Options : Service options
type Db2saasV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewDb2saasV1UsingExternalConfig : constructs an instance of Db2saasV1 with passed in options and external configuration.
func NewDb2saasV1UsingExternalConfig(options *Db2saasV1Options) (db2saas *Db2saasV1, err error) {
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

	db2saas, err = NewDb2saasV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = db2saas.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = db2saas.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewDb2saasV1 : constructs an instance of Db2saasV1 with passed in options.
func NewDb2saasV1(options *Db2saasV1Options) (service *Db2saasV1, err error) {
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

	service = &Db2saasV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "db2saas" suitable for processing requests.
func (db2saas *Db2saasV1) Clone() *Db2saasV1 {
	if core.IsNil(db2saas) {
		return nil
	}
	clone := *db2saas
	clone.Service = db2saas.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (db2saas *Db2saasV1) SetServiceURL(url string) error {
	err := db2saas.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (db2saas *Db2saasV1) GetServiceURL() string {
	return db2saas.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (db2saas *Db2saasV1) SetDefaultHeaders(headers http.Header) {
	db2saas.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (db2saas *Db2saasV1) SetEnableGzipCompression(enableGzip bool) {
	db2saas.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (db2saas *Db2saasV1) GetEnableGzipCompression() bool {
	return db2saas.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (db2saas *Db2saasV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	db2saas.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (db2saas *Db2saasV1) DisableRetries() {
	db2saas.Service.DisableRetries()
}

// GetDb2SaasConnectionInfo : Get Db2 connection information
func (db2saas *Db2saasV1) GetDb2SaasConnectionInfo(getDb2SaasConnectionInfoOptions *GetDb2SaasConnectionInfoOptions) (result *SuccessConnectionInfo, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasConnectionInfoWithContext(context.Background(), getDb2SaasConnectionInfoOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasConnectionInfoWithContext is an alternate form of the GetDb2SaasConnectionInfo method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasConnectionInfoWithContext(ctx context.Context, getDb2SaasConnectionInfoOptions *GetDb2SaasConnectionInfoOptions) (result *SuccessConnectionInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDb2SaasConnectionInfoOptions, "getDb2SaasConnectionInfoOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDb2SaasConnectionInfoOptions, "getDb2SaasConnectionInfoOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"deployment_id": *getDb2SaasConnectionInfoOptions.DeploymentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/connectioninfo/{deployment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasConnectionInfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasConnectionInfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDb2SaasConnectionInfoOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*getDb2SaasConnectionInfoOptions.XDeploymentID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_connection_info", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessConnectionInfo)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PostDb2SaasAllowlist : Allow listing of new IPs
func (db2saas *Db2saasV1) PostDb2SaasAllowlist(postDb2SaasAllowlistOptions *PostDb2SaasAllowlistOptions) (result *SuccessPostAllowedlistIPs, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.PostDb2SaasAllowlistWithContext(context.Background(), postDb2SaasAllowlistOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PostDb2SaasAllowlistWithContext is an alternate form of the PostDb2SaasAllowlist method which supports a Context parameter
func (db2saas *Db2saasV1) PostDb2SaasAllowlistWithContext(ctx context.Context, postDb2SaasAllowlistOptions *PostDb2SaasAllowlistOptions) (result *SuccessPostAllowedlistIPs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDb2SaasAllowlistOptions, "postDb2SaasAllowlistOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(postDb2SaasAllowlistOptions, "postDb2SaasAllowlistOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/dbsettings/whitelistips`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range postDb2SaasAllowlistOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "PostDb2SaasAllowlist")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postDb2SaasAllowlistOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*postDb2SaasAllowlistOptions.XDeploymentID))
	}

	body := make(map[string]interface{})
	if postDb2SaasAllowlistOptions.IpAddresses != nil {
		body["ip_addresses"] = postDb2SaasAllowlistOptions.IpAddresses
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
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "post_db2_saas_allowlist", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessPostAllowedlistIPs)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDb2SaasAllowlist : Get allowed list of IPs
func (db2saas *Db2saasV1) GetDb2SaasAllowlist(getDb2SaasAllowlistOptions *GetDb2SaasAllowlistOptions) (result *SuccessGetAllowlistIPs, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasAllowlistWithContext(context.Background(), getDb2SaasAllowlistOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasAllowlistWithContext is an alternate form of the GetDb2SaasAllowlist method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasAllowlistWithContext(ctx context.Context, getDb2SaasAllowlistOptions *GetDb2SaasAllowlistOptions) (result *SuccessGetAllowlistIPs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDb2SaasAllowlistOptions, "getDb2SaasAllowlistOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDb2SaasAllowlistOptions, "getDb2SaasAllowlistOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/dbsettings/whitelistips`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasAllowlistOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasAllowlist")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDb2SaasAllowlistOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*getDb2SaasAllowlistOptions.XDeploymentID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_allowlist", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessGetAllowlistIPs)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PostDb2SaasUser : Create new user ( available only for platform users)
func (db2saas *Db2saasV1) PostDb2SaasUser(postDb2SaasUserOptions *PostDb2SaasUserOptions) (result *SuccessUserResponse, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.PostDb2SaasUserWithContext(context.Background(), postDb2SaasUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PostDb2SaasUserWithContext is an alternate form of the PostDb2SaasUser method which supports a Context parameter
func (db2saas *Db2saasV1) PostDb2SaasUserWithContext(ctx context.Context, postDb2SaasUserOptions *PostDb2SaasUserOptions) (result *SuccessUserResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDb2SaasUserOptions, "postDb2SaasUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(postDb2SaasUserOptions, "postDb2SaasUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/users`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range postDb2SaasUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "PostDb2SaasUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postDb2SaasUserOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*postDb2SaasUserOptions.XDeploymentID))
	}

	body := make(map[string]interface{})
	if postDb2SaasUserOptions.ID != nil {
		body["id"] = postDb2SaasUserOptions.ID
	}
	if postDb2SaasUserOptions.Iam != nil {
		body["iam"] = postDb2SaasUserOptions.Iam
	}
	if postDb2SaasUserOptions.Ibmid != nil {
		body["ibmid"] = postDb2SaasUserOptions.Ibmid
	}
	if postDb2SaasUserOptions.Name != nil {
		body["name"] = postDb2SaasUserOptions.Name
	}
	if postDb2SaasUserOptions.Password != nil {
		body["password"] = postDb2SaasUserOptions.Password
	}
	if postDb2SaasUserOptions.Role != nil {
		body["role"] = postDb2SaasUserOptions.Role
	}
	if postDb2SaasUserOptions.Email != nil {
		body["email"] = postDb2SaasUserOptions.Email
	}
	if postDb2SaasUserOptions.Locked != nil {
		body["locked"] = postDb2SaasUserOptions.Locked
	}
	if postDb2SaasUserOptions.Authentication != nil {
		body["authentication"] = postDb2SaasUserOptions.Authentication
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
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "post_db2_saas_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessUserResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDb2SaasUser : Get the list of Users
func (db2saas *Db2saasV1) GetDb2SaasUser(getDb2SaasUserOptions *GetDb2SaasUserOptions) (result *SuccessGetUserInfo, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasUserWithContext(context.Background(), getDb2SaasUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasUserWithContext is an alternate form of the GetDb2SaasUser method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasUserWithContext(ctx context.Context, getDb2SaasUserOptions *GetDb2SaasUserOptions) (result *SuccessGetUserInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDb2SaasUserOptions, "getDb2SaasUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDb2SaasUserOptions, "getDb2SaasUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/users`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDb2SaasUserOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*getDb2SaasUserOptions.XDeploymentID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessGetUserInfo)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDb2SaasUser : Delete a user (only platform admin)
func (db2saas *Db2saasV1) DeleteDb2SaasUser(deleteDb2SaasUserOptions *DeleteDb2SaasUserOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.DeleteDb2SaasUserWithContext(context.Background(), deleteDb2SaasUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDb2SaasUserWithContext is an alternate form of the DeleteDb2SaasUser method which supports a Context parameter
func (db2saas *Db2saasV1) DeleteDb2SaasUserWithContext(ctx context.Context, deleteDb2SaasUserOptions *DeleteDb2SaasUserOptions) (result map[string]interface{}, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDb2SaasUserOptions, "deleteDb2SaasUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDb2SaasUserOptions, "deleteDb2SaasUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteDb2SaasUserOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/users/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDb2SaasUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "DeleteDb2SaasUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDb2SaasUserOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*deleteDb2SaasUserOptions.XDeploymentID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = db2saas.Service.Request(request, &result)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_db2_saas_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetbyidDb2SaasUser : Get specific user by Id
func (db2saas *Db2saasV1) GetbyidDb2SaasUser(getbyidDb2SaasUserOptions *GetbyidDb2SaasUserOptions) (result *SuccessGetUserByID, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetbyidDb2SaasUserWithContext(context.Background(), getbyidDb2SaasUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetbyidDb2SaasUserWithContext is an alternate form of the GetbyidDb2SaasUser method which supports a Context parameter
func (db2saas *Db2saasV1) GetbyidDb2SaasUserWithContext(ctx context.Context, getbyidDb2SaasUserOptions *GetbyidDb2SaasUserOptions) (result *SuccessGetUserByID, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getbyidDb2SaasUserOptions, "getbyidDb2SaasUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getbyidDb2SaasUserOptions, "getbyidDb2SaasUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/users/bluadmin`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getbyidDb2SaasUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetbyidDb2SaasUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getbyidDb2SaasUserOptions.XDeploymentID != nil {
		builder.AddHeader("x-deployment-id", fmt.Sprint(*getbyidDb2SaasUserOptions.XDeploymentID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "getbyid_db2_saas_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessGetUserByID)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PutDb2SaasAutoscale : Update auto scaling configuration
func (db2saas *Db2saasV1) PutDb2SaasAutoscale(putDb2SaasAutoscaleOptions *PutDb2SaasAutoscaleOptions) (result *SuccessUpdateAutoScale, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.PutDb2SaasAutoscaleWithContext(context.Background(), putDb2SaasAutoscaleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PutDb2SaasAutoscaleWithContext is an alternate form of the PutDb2SaasAutoscale method which supports a Context parameter
func (db2saas *Db2saasV1) PutDb2SaasAutoscaleWithContext(ctx context.Context, putDb2SaasAutoscaleOptions *PutDb2SaasAutoscaleOptions) (result *SuccessUpdateAutoScale, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putDb2SaasAutoscaleOptions, "putDb2SaasAutoscaleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(putDb2SaasAutoscaleOptions, "putDb2SaasAutoscaleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/scaling/auto`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range putDb2SaasAutoscaleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "PutDb2SaasAutoscale")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if putDb2SaasAutoscaleOptions.XDbProfile != nil {
		builder.AddHeader("x-db-profile", fmt.Sprint(*putDb2SaasAutoscaleOptions.XDbProfile))
	}

	body := make(map[string]interface{})
	if putDb2SaasAutoscaleOptions.AutoScalingEnabled != nil {
		body["auto_scaling_enabled"] = putDb2SaasAutoscaleOptions.AutoScalingEnabled
	}
	if putDb2SaasAutoscaleOptions.AutoScalingThreshold != nil {
		body["auto_scaling_threshold"] = putDb2SaasAutoscaleOptions.AutoScalingThreshold
	}
	if putDb2SaasAutoscaleOptions.AutoScalingOverTimePeriod != nil {
		body["auto_scaling_over_time_period"] = putDb2SaasAutoscaleOptions.AutoScalingOverTimePeriod
	}
	if putDb2SaasAutoscaleOptions.AutoScalingPauseLimit != nil {
		body["auto_scaling_pause_limit"] = putDb2SaasAutoscaleOptions.AutoScalingPauseLimit
	}
	if putDb2SaasAutoscaleOptions.AutoScalingAllowPlanLimit != nil {
		body["auto_scaling_allow_plan_limit"] = putDb2SaasAutoscaleOptions.AutoScalingAllowPlanLimit
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
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "put_db2_saas_autoscale", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessUpdateAutoScale)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDb2SaasAutoscale : Get auto scaling info
func (db2saas *Db2saasV1) GetDb2SaasAutoscale(getDb2SaasAutoscaleOptions *GetDb2SaasAutoscaleOptions) (result *SuccessAutoScaling, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasAutoscaleWithContext(context.Background(), getDb2SaasAutoscaleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasAutoscaleWithContext is an alternate form of the GetDb2SaasAutoscale method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasAutoscaleWithContext(ctx context.Context, getDb2SaasAutoscaleOptions *GetDb2SaasAutoscaleOptions) (result *SuccessAutoScaling, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDb2SaasAutoscaleOptions, "getDb2SaasAutoscaleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDb2SaasAutoscaleOptions, "getDb2SaasAutoscaleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/scaling/auto`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasAutoscaleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasAutoscale")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDb2SaasAutoscaleOptions.XDbProfile != nil {
		builder.AddHeader("x-db-profile", fmt.Sprint(*getDb2SaasAutoscaleOptions.XDbProfile))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_autoscale", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessAutoScaling)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PostDb2SaasDbConfiguration : Set database and database manager configuration
func (db2saas *Db2saasV1) PostDb2SaasDbConfiguration(postDb2SaasDbConfigurationOptions *PostDb2SaasDbConfigurationOptions) (result *SuccessPostCustomSettings, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.PostDb2SaasDbConfigurationWithContext(context.Background(), postDb2SaasDbConfigurationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PostDb2SaasDbConfigurationWithContext is an alternate form of the PostDb2SaasDbConfiguration method which supports a Context parameter
func (db2saas *Db2saasV1) PostDb2SaasDbConfigurationWithContext(ctx context.Context, postDb2SaasDbConfigurationOptions *PostDb2SaasDbConfigurationOptions) (result *SuccessPostCustomSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDb2SaasDbConfigurationOptions, "postDb2SaasDbConfigurationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(postDb2SaasDbConfigurationOptions, "postDb2SaasDbConfigurationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/deployments/custom_setting`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range postDb2SaasDbConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "PostDb2SaasDbConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if postDb2SaasDbConfigurationOptions.XDbProfile != nil {
		builder.AddHeader("x-db-profile", fmt.Sprint(*postDb2SaasDbConfigurationOptions.XDbProfile))
	}

	body := make(map[string]interface{})
	if postDb2SaasDbConfigurationOptions.Registry != nil {
		body["registry"] = postDb2SaasDbConfigurationOptions.Registry
	}
	if postDb2SaasDbConfigurationOptions.Db != nil {
		body["db"] = postDb2SaasDbConfigurationOptions.Db
	}
	if postDb2SaasDbConfigurationOptions.Dbm != nil {
		body["dbm"] = postDb2SaasDbConfigurationOptions.Dbm
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
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "post_db2_saas_db_configuration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessPostCustomSettings)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDb2SaasTuneableParam : Retrieves the values of tunable parameters of the DB2 instance
func (db2saas *Db2saasV1) GetDb2SaasTuneableParam(getDb2SaasTuneableParamOptions *GetDb2SaasTuneableParamOptions) (result *SuccessTuneableParams, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasTuneableParamWithContext(context.Background(), getDb2SaasTuneableParamOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasTuneableParamWithContext is an alternate form of the GetDb2SaasTuneableParam method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasTuneableParamWithContext(ctx context.Context, getDb2SaasTuneableParamOptions *GetDb2SaasTuneableParamOptions) (result *SuccessTuneableParams, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getDb2SaasTuneableParamOptions, "getDb2SaasTuneableParamOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/tuneable_param`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasTuneableParamOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasTuneableParam")
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
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_tuneable_param", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessTuneableParams)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDb2SaasBackup : Get Db2 instance backup information
func (db2saas *Db2saasV1) GetDb2SaasBackup(getDb2SaasBackupOptions *GetDb2SaasBackupOptions) (result *SuccessGetBackups, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.GetDb2SaasBackupWithContext(context.Background(), getDb2SaasBackupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDb2SaasBackupWithContext is an alternate form of the GetDb2SaasBackup method which supports a Context parameter
func (db2saas *Db2saasV1) GetDb2SaasBackupWithContext(ctx context.Context, getDb2SaasBackupOptions *GetDb2SaasBackupOptions) (result *SuccessGetBackups, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDb2SaasBackupOptions, "getDb2SaasBackupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDb2SaasBackupOptions, "getDb2SaasBackupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/backups`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDb2SaasBackupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "GetDb2SaasBackup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDb2SaasBackupOptions.XDbProfile != nil {
		builder.AddHeader("x-db-profile", fmt.Sprint(*getDb2SaasBackupOptions.XDbProfile))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_db2_saas_backup", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessGetBackups)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PostDb2SaasBackup : Create backup of an instance
func (db2saas *Db2saasV1) PostDb2SaasBackup(postDb2SaasBackupOptions *PostDb2SaasBackupOptions) (result *SuccessCreateBackup, response *core.DetailedResponse, err error) {
	result, response, err = db2saas.PostDb2SaasBackupWithContext(context.Background(), postDb2SaasBackupOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PostDb2SaasBackupWithContext is an alternate form of the PostDb2SaasBackup method which supports a Context parameter
func (db2saas *Db2saasV1) PostDb2SaasBackupWithContext(ctx context.Context, postDb2SaasBackupOptions *PostDb2SaasBackupOptions) (result *SuccessCreateBackup, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postDb2SaasBackupOptions, "postDb2SaasBackupOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(postDb2SaasBackupOptions, "postDb2SaasBackupOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = db2saas.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(db2saas.Service.Options.URL, `/manage/backups/backup`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range postDb2SaasBackupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("db2saas", "V1", "PostDb2SaasBackup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if postDb2SaasBackupOptions.XDbProfile != nil {
		builder.AddHeader("x-db-profile", fmt.Sprint(*postDb2SaasBackupOptions.XDbProfile))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = db2saas.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "post_db2_saas_backup", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSuccessCreateBackup)
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

// Backup : Info of backup.
type Backup struct {
	// CRN of the db2 instance.
	ID *string `json:"id" validate:"required"`

	// Defines the type of execution of backup.
	Type *string `json:"type" validate:"required"`

	// Status of the backup.
	Status *string `json:"status" validate:"required"`

	// Timestamp of the backup created.
	CreatedAt *string `json:"created_at" validate:"required"`

	// Size of the backup or data set.
	Size *int64 `json:"size" validate:"required"`

	// The duration of the backup operation in seconds.
	Duration *int64 `json:"duration" validate:"required"`
}

// UnmarshalBackup unmarshals an instance of Backup from the specified map of raw messages.
func UnmarshalBackup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Backup)
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
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "duration", &obj.Duration)
	if err != nil {
		err = core.SDKErrorf(err, "", "duration-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCustomSettingsDb : Container for general database settings.
type CreateCustomSettingsDb struct {
	// Configures the sort memory limit for DB2.
	ACTSORTMEMLIMIT *string `json:"ACT_SORTMEM_LIMIT,omitempty"`

	// Configures the collation sequence.
	ALTCOLLATE *string `json:"ALT_COLLATE,omitempty"`

	// Sets the application group memory size.
	APPGROUPMEMSZ *string `json:"APPGROUP_MEM_SZ,omitempty"`

	// Configures the application heap size.
	APPLHEAPSZ *string `json:"APPLHEAPSZ,omitempty"`

	// Configures the application memory allocation.
	APPLMEMORY *string `json:"APPL_MEMORY,omitempty"`

	// Configures the application control heap size.
	APPCTLHEAPSZ *string `json:"APP_CTL_HEAP_SZ,omitempty"`

	// Configures the archive retry delay time.
	ARCHRETRYDELAY *string `json:"ARCHRETRYDELAY,omitempty"`

	// Configures the authentication cache duration.
	AUTHNCACHEDURATION *string `json:"AUTHN_CACHE_DURATION,omitempty"`

	// Configures whether the database will automatically restart.
	AUTORESTART *string `json:"AUTORESTART,omitempty"`

	// Configures whether auto collection of CG statistics is enabled.
	AUTOCGSTATS *string `json:"AUTO_CG_STATS,omitempty"`

	// Configures automatic maintenance for the database.
	AUTOMAINT *string `json:"AUTO_MAINT,omitempty"`

	// Configures automatic reorganization for the database.
	AUTOREORG *string `json:"AUTO_REORG,omitempty"`

	// Configures the auto refresh or revalidation method.
	AUTOREVAL *string `json:"AUTO_REVAL,omitempty"`

	// Configures automatic collection of run-time statistics.
	AUTORUNSTATS *string `json:"AUTO_RUNSTATS,omitempty"`

	// Configures whether auto-sampling is enabled.
	AUTOSAMPLING *string `json:"AUTO_SAMPLING,omitempty"`

	// Configures automatic collection of statistics on views.
	AUTOSTATSVIEWS *string `json:"AUTO_STATS_VIEWS,omitempty"`

	// Configures automatic collection of statement-level statistics.
	AUTOSTMTSTATS *string `json:"AUTO_STMT_STATS,omitempty"`

	// Configures automatic table maintenance.
	AUTOTBLMAINT *string `json:"AUTO_TBL_MAINT,omitempty"`

	// Average number of applications.
	AVGAPPLS *string `json:"AVG_APPLS,omitempty"`

	// Configures the catalog cache size.
	CATALOGCACHESZ *string `json:"CATALOGCACHE_SZ,omitempty"`

	// Configures the change pages threshold percentage.
	CHNGPGSTHRESH *string `json:"CHNGPGS_THRESH,omitempty"`

	// Configures the commit behavior.
	CURCOMMIT *string `json:"CUR_COMMIT,omitempty"`

	// Configures the database memory management.
	DATABASEMEMORY *string `json:"DATABASE_MEMORY,omitempty"`

	// Configures the database heap size.
	DBHEAP *string `json:"DBHEAP,omitempty"`

	// Specifies the database collation name.
	DBCOLLNAME *string `json:"DB_COLLNAME,omitempty"`

	// Configures the memory threshold percentage for database.
	DBMEMTHRESH *string `json:"DB_MEM_THRESH,omitempty"`

	// Defines the default DDL compression behavior.
	DDLCOMPRESSIONDEF *string `json:"DDL_COMPRESSION_DEF,omitempty"`

	// Defines the default constraint behavior in DDL.
	DDLCONSTRAINTDEF *string `json:"DDL_CONSTRAINT_DEF,omitempty"`

	// Configures the decimal floating-point rounding method.
	DECFLTROUNDING *string `json:"DECFLT_ROUNDING,omitempty"`

	// Configures the default arithmetic for decimal operations.
	DECARITHMETIC *string `json:"DEC_ARITHMETIC,omitempty"`

	// Configures the decimal-to-character conversion format.
	DECTOCHARFMT *string `json:"DEC_TO_CHAR_FMT,omitempty"`

	// Configures the default degree for parallelism.
	DFTDEGREE *string `json:"DFT_DEGREE,omitempty"`

	// Configures the default extent size for tables.
	DFTEXTENTSZ *string `json:"DFT_EXTENT_SZ,omitempty"`

	// Configures the default load record session count.
	DFTLOADRECSES *string `json:"DFT_LOADREC_SES,omitempty"`

	// Configures the default MTTB (multi-table table scan) types.
	DFTMTTBTYPES *string `json:"DFT_MTTB_TYPES,omitempty"`

	// Configures the default prefetch size for queries.
	DFTPREFETCHSZ *string `json:"DFT_PREFETCH_SZ,omitempty"`

	// Configures the default query optimization level.
	DFTQUERYOPT *string `json:"DFT_QUERYOPT,omitempty"`

	// Configures the default refresh age for views.
	DFTREFRESHAGE *string `json:"DFT_REFRESH_AGE,omitempty"`

	// Configures whether DCC (database control center) is enabled for schemas.
	DFTSCHEMASDCC *string `json:"DFT_SCHEMAS_DCC,omitempty"`

	// Configures whether SQL math warnings are enabled.
	DFTSQLMATHWARN *string `json:"DFT_SQLMATHWARN,omitempty"`

	// Configures the default table organization (ROW or COLUMN).
	DFTTABLEORG *string `json:"DFT_TABLE_ORG,omitempty"`

	// Configures the deadlock check time in milliseconds.
	DLCHKTIME *string `json:"DLCHKTIME,omitempty"`

	// Configures whether XML character support is enabled.
	ENABLEXMLCHAR *string `json:"ENABLE_XMLCHAR,omitempty"`

	// Configures whether extended row size is enabled.
	EXTENDEDROWSZ *string `json:"EXTENDED_ROW_SZ,omitempty"`

	// Configures the heap ratio for group heap memory.
	GROUPHEAPRATIO *string `json:"GROUPHEAP_RATIO,omitempty"`

	// Configures the index recovery method.
	INDEXREC *string `json:"INDEXREC,omitempty"`

	// Configures whether large aggregation is enabled.
	LARGEAGGREGATION *string `json:"LARGE_AGGREGATION,omitempty"`

	// Configures the lock list memory size.
	LOCKLIST *string `json:"LOCKLIST,omitempty"`

	// Configures the lock timeout duration.
	LOCKTIMEOUT *string `json:"LOCKTIMEOUT,omitempty"`

	// Configures whether index builds are logged.
	LOGINDEXBUILD *string `json:"LOGINDEXBUILD,omitempty"`

	// Configures whether application information is logged.
	LOGAPPLINFO *string `json:"LOG_APPL_INFO,omitempty"`

	// Configures whether DDL statements are logged.
	LOGDDLSTMTS *string `json:"LOG_DDL_STMTS,omitempty"`

	// Configures the disk capacity log setting.
	LOGDISKCAP *string `json:"LOG_DISK_CAP,omitempty"`

	// Configures the maximum number of applications.
	MAXAPPLS *string `json:"MAXAPPLS,omitempty"`

	// Configures the maximum number of file operations.
	MAXFILOP *string `json:"MAXFILOP,omitempty"`

	// Configures the maximum number of locks.
	MAXLOCKS *string `json:"MAXLOCKS,omitempty"`

	// Configures whether decimal division by 3 should be handled.
	MINDECDIV3 *string `json:"MIN_DEC_DIV_3,omitempty"`

	// Configures the level of activity metrics to be monitored.
	MONACTMETRICS *string `json:"MON_ACT_METRICS,omitempty"`

	// Configures deadlock monitoring settings.
	MONDEADLOCK *string `json:"MON_DEADLOCK,omitempty"`

	// Configures the lock message level for monitoring.
	MONLCKMSGLVL *string `json:"MON_LCK_MSG_LVL,omitempty"`

	// Configures lock timeout monitoring settings.
	MONLOCKTIMEOUT *string `json:"MON_LOCKTIMEOUT,omitempty"`

	// Configures lock wait monitoring settings.
	MONLOCKWAIT *string `json:"MON_LOCKWAIT,omitempty"`

	// Configures the lightweight threshold for monitoring.
	MONLWTHRESH *string `json:"MON_LW_THRESH,omitempty"`

	// Configures the object metrics level for monitoring.
	MONOBJMETRICS *string `json:"MON_OBJ_METRICS,omitempty"`

	// Configures the package list size for monitoring.
	MONPKGLISTSZ *string `json:"MON_PKGLIST_SZ,omitempty"`

	// Configures the request metrics level for monitoring.
	MONREQMETRICS *string `json:"MON_REQ_METRICS,omitempty"`

	// Configures the level of return data for monitoring.
	MONRTNDATA *string `json:"MON_RTN_DATA,omitempty"`

	// Configures whether stored procedure execution list is monitored.
	MONRTNEXECLIST *string `json:"MON_RTN_EXECLIST,omitempty"`

	// Configures the level of unit of work (UOW) data for monitoring.
	MONUOWDATA *string `json:"MON_UOW_DATA,omitempty"`

	// Configures whether UOW execution list is monitored.
	MONUOWEXECLIST *string `json:"MON_UOW_EXECLIST,omitempty"`

	// Configures whether UOW package list is monitored.
	MONUOWPKGLIST *string `json:"MON_UOW_PKGLIST,omitempty"`

	// Configures the mapping of NCHAR character types.
	NCHARMAPPING *string `json:"NCHAR_MAPPING,omitempty"`

	// Configures the number of frequent values for optimization.
	NUMFREQVALUES *string `json:"NUM_FREQVALUES,omitempty"`

	// Configures the number of IO cleaners.
	NUMIOCLEANERS *string `json:"NUM_IOCLEANERS,omitempty"`

	// Configures the number of IO servers.
	NUMIOSERVERS *string `json:"NUM_IOSERVERS,omitempty"`

	// Configures the number of log spans.
	NUMLOGSPAN *string `json:"NUM_LOG_SPAN,omitempty"`

	// Configures the number of quantiles for optimizations.
	NUMQUANTILES *string `json:"NUM_QUANTILES,omitempty"`

	// Configures the buffer page optimization setting.
	OPTBUFFPAGE *string `json:"OPT_BUFFPAGE,omitempty"`

	// Configures the direct workload optimization setting.
	OPTDIRECTWRKLD *string `json:"OPT_DIRECT_WRKLD,omitempty"`

	// Configures the lock list optimization setting.
	OPTLOCKLIST *string `json:"OPT_LOCKLIST,omitempty"`

	// Configures the max locks optimization setting.
	OPTMAXLOCKS *string `json:"OPT_MAXLOCKS,omitempty"`

	// Configures the sort heap optimization setting.
	OPTSORTHEAP *string `json:"OPT_SORTHEAP,omitempty"`

	// Configures the page age target for garbage collection.
	PAGEAGETRGTGCR *string `json:"PAGE_AGE_TRGT_GCR,omitempty"`

	// Configures the page age target for memory collection.
	PAGEAGETRGTMCR *string `json:"PAGE_AGE_TRGT_MCR,omitempty"`

	// Configures the package cache size.
	PCKCACHESZ *string `json:"PCKCACHESZ,omitempty"`

	// Configures the level of stack trace logging for stored procedures.
	PLSTACKTRACE *string `json:"PL_STACK_TRACE,omitempty"`

	// Configures whether self-tuning memory is enabled.
	SELFTUNINGMEM *string `json:"SELF_TUNING_MEM,omitempty"`

	// Configures sequence detection for queries.
	SEQDETECT *string `json:"SEQDETECT,omitempty"`

	// Configures the shared heap threshold size.
	SHEAPTHRESSHR *string `json:"SHEAPTHRES_SHR,omitempty"`

	// Configures the soft max setting.
	SOFTMAX *string `json:"SOFTMAX,omitempty"`

	// Configures the sort heap memory size.
	SORTHEAP *string `json:"SORTHEAP,omitempty"`

	// Configures the SQL compiler flags.
	SQLCCFLAGS *string `json:"SQL_CCFLAGS,omitempty"`

	// Configures the statistics heap size.
	STATHEAPSZ *string `json:"STAT_HEAP_SZ,omitempty"`

	// Configures the statement heap size.
	STMTHEAP *string `json:"STMTHEAP,omitempty"`

	// Configures the statement concurrency.
	STMTCONC *string `json:"STMT_CONC,omitempty"`

	// Configures the string unit settings.
	STRINGUNITS *string `json:"STRING_UNITS,omitempty"`

	// Configures whether system time period adjustments are enabled.
	SYSTIMEPERIODADJ *string `json:"SYSTIME_PERIOD_ADJ,omitempty"`

	// Configures whether modifications to tracked objects are logged.
	TRACKMOD *string `json:"TRACKMOD,omitempty"`

	// Configures the utility heap size.
	UTILHEAPSZ *string `json:"UTIL_HEAP_SZ,omitempty"`

	// Configures whether WLM (Workload Management) admission control is enabled.
	WLMADMISSIONCTRL *string `json:"WLM_ADMISSION_CTRL,omitempty"`

	// Configures the WLM agent load target.
	WLMAGENTLOADTRGT *string `json:"WLM_AGENT_LOAD_TRGT,omitempty"`

	// Configures the CPU limit for WLM workloads.
	WLMCPULIMIT *string `json:"WLM_CPU_LIMIT,omitempty"`

	// Configures the CPU share count for WLM workloads.
	WLMCPUSHARES *string `json:"WLM_CPU_SHARES,omitempty"`

	// Configures the mode of CPU shares for WLM workloads.
	WLMCPUSHAREMODE *string `json:"WLM_CPU_SHARE_MODE,omitempty"`
}

// Constants associated with the CreateCustomSettingsDb.ACTSORTMEMLIMIT property.
// Configures the sort memory limit for DB2.
const (
	CreateCustomSettingsDb_ACTSORTMEMLIMIT_None = "NONE"
	CreateCustomSettingsDb_ACTSORTMEMLIMIT_Range10100 = "range(10, 100)"
)

// Constants associated with the CreateCustomSettingsDb.ALTCOLLATE property.
// Configures the collation sequence.
const (
	CreateCustomSettingsDb_ALTCOLLATE_Identity16bit = "IDENTITY_16BIT"
	CreateCustomSettingsDb_ALTCOLLATE_Null = "NULL"
)

// Constants associated with the CreateCustomSettingsDb.APPGROUPMEMSZ property.
// Sets the application group memory size.
const (
	CreateCustomSettingsDb_APPGROUPMEMSZ_Range11000000 = "range(1, 1000000)"
)

// Constants associated with the CreateCustomSettingsDb.APPLHEAPSZ property.
// Configures the application heap size.
const (
	CreateCustomSettingsDb_APPLHEAPSZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_APPLHEAPSZ_Range162147483647 = "range(16, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.APPLMEMORY property.
// Configures the application memory allocation.
const (
	CreateCustomSettingsDb_APPLMEMORY_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_APPLMEMORY_Range1284294967295 = "range(128, 4294967295)"
)

// Constants associated with the CreateCustomSettingsDb.APPCTLHEAPSZ property.
// Configures the application control heap size.
const (
	CreateCustomSettingsDb_APPCTLHEAPSZ_Range164000 = "range(1, 64000)"
)

// Constants associated with the CreateCustomSettingsDb.ARCHRETRYDELAY property.
// Configures the archive retry delay time.
const (
	CreateCustomSettingsDb_ARCHRETRYDELAY_Range065535 = "range(0, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.AUTHNCACHEDURATION property.
// Configures the authentication cache duration.
const (
	CreateCustomSettingsDb_AUTHNCACHEDURATION_Range110000 = "range(1,10000)"
)

// Constants associated with the CreateCustomSettingsDb.AUTORESTART property.
// Configures whether the database will automatically restart.
const (
	CreateCustomSettingsDb_AUTORESTART_Off = "OFF"
	CreateCustomSettingsDb_AUTORESTART_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOCGSTATS property.
// Configures whether auto collection of CG statistics is enabled.
const (
	CreateCustomSettingsDb_AUTOCGSTATS_Off = "OFF"
	CreateCustomSettingsDb_AUTOCGSTATS_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOMAINT property.
// Configures automatic maintenance for the database.
const (
	CreateCustomSettingsDb_AUTOMAINT_Off = "OFF"
	CreateCustomSettingsDb_AUTOMAINT_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOREORG property.
// Configures automatic reorganization for the database.
const (
	CreateCustomSettingsDb_AUTOREORG_Off = "OFF"
	CreateCustomSettingsDb_AUTOREORG_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOREVAL property.
// Configures the auto refresh or revalidation method.
const (
	CreateCustomSettingsDb_AUTOREVAL_Deferred = "DEFERRED"
	CreateCustomSettingsDb_AUTOREVAL_DeferredForce = "DEFERRED_FORCE"
	CreateCustomSettingsDb_AUTOREVAL_Disabled = "DISABLED"
	CreateCustomSettingsDb_AUTOREVAL_Immediate = "IMMEDIATE"
)

// Constants associated with the CreateCustomSettingsDb.AUTORUNSTATS property.
// Configures automatic collection of run-time statistics.
const (
	CreateCustomSettingsDb_AUTORUNSTATS_Off = "OFF"
	CreateCustomSettingsDb_AUTORUNSTATS_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOSAMPLING property.
// Configures whether auto-sampling is enabled.
const (
	CreateCustomSettingsDb_AUTOSAMPLING_Off = "OFF"
	CreateCustomSettingsDb_AUTOSAMPLING_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOSTATSVIEWS property.
// Configures automatic collection of statistics on views.
const (
	CreateCustomSettingsDb_AUTOSTATSVIEWS_Off = "OFF"
	CreateCustomSettingsDb_AUTOSTATSVIEWS_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOSTMTSTATS property.
// Configures automatic collection of statement-level statistics.
const (
	CreateCustomSettingsDb_AUTOSTMTSTATS_Off = "OFF"
	CreateCustomSettingsDb_AUTOSTMTSTATS_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.AUTOTBLMAINT property.
// Configures automatic table maintenance.
const (
	CreateCustomSettingsDb_AUTOTBLMAINT_Off = "OFF"
	CreateCustomSettingsDb_AUTOTBLMAINT_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.CHNGPGSTHRESH property.
// Configures the change pages threshold percentage.
const (
	CreateCustomSettingsDb_CHNGPGSTHRESH_Range599 = "range(5,99)"
)

// Constants associated with the CreateCustomSettingsDb.CURCOMMIT property.
// Configures the commit behavior.
const (
	CreateCustomSettingsDb_CURCOMMIT_Available = "AVAILABLE"
	CreateCustomSettingsDb_CURCOMMIT_Disabled = "DISABLED"
	CreateCustomSettingsDb_CURCOMMIT_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.DATABASEMEMORY property.
// Configures the database memory management.
const (
	CreateCustomSettingsDb_DATABASEMEMORY_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_DATABASEMEMORY_Computed = "COMPUTED"
	CreateCustomSettingsDb_DATABASEMEMORY_Range04294967295 = "range(0, 4294967295)"
)

// Constants associated with the CreateCustomSettingsDb.DBHEAP property.
// Configures the database heap size.
const (
	CreateCustomSettingsDb_DBHEAP_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_DBHEAP_Range322147483647 = "range(32, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.DBMEMTHRESH property.
// Configures the memory threshold percentage for database.
const (
	CreateCustomSettingsDb_DBMEMTHRESH_Range0100 = "range(0, 100)"
)

// Constants associated with the CreateCustomSettingsDb.DDLCOMPRESSIONDEF property.
// Defines the default DDL compression behavior.
const (
	CreateCustomSettingsDb_DDLCOMPRESSIONDEF_No = "NO"
	CreateCustomSettingsDb_DDLCOMPRESSIONDEF_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.DDLCONSTRAINTDEF property.
// Defines the default constraint behavior in DDL.
const (
	CreateCustomSettingsDb_DDLCONSTRAINTDEF_No = "NO"
	CreateCustomSettingsDb_DDLCONSTRAINTDEF_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.DECFLTROUNDING property.
// Configures the decimal floating-point rounding method.
const (
	CreateCustomSettingsDb_DECFLTROUNDING_RoundCeiling = "ROUND_CEILING"
	CreateCustomSettingsDb_DECFLTROUNDING_RoundDown = "ROUND_DOWN"
	CreateCustomSettingsDb_DECFLTROUNDING_RoundFloor = "ROUND_FLOOR"
	CreateCustomSettingsDb_DECFLTROUNDING_RoundHalfEven = "ROUND_HALF_EVEN"
	CreateCustomSettingsDb_DECFLTROUNDING_RoundHalfUp = "ROUND_HALF_UP"
)

// Constants associated with the CreateCustomSettingsDb.DECTOCHARFMT property.
// Configures the decimal-to-character conversion format.
const (
	CreateCustomSettingsDb_DECTOCHARFMT_New = "NEW"
	CreateCustomSettingsDb_DECTOCHARFMT_V95 = "V95"
)

// Constants associated with the CreateCustomSettingsDb.DFTDEGREE property.
// Configures the default degree for parallelism.
const (
	CreateCustomSettingsDb_DFTDEGREE_Any = "ANY"
	CreateCustomSettingsDb_DFTDEGREE_Range132767 = "range(1, 32767)"
)

// Constants associated with the CreateCustomSettingsDb.DFTEXTENTSZ property.
// Configures the default extent size for tables.
const (
	CreateCustomSettingsDb_DFTEXTENTSZ_Range2256 = "range(2, 256)"
)

// Constants associated with the CreateCustomSettingsDb.DFTLOADRECSES property.
// Configures the default load record session count.
const (
	CreateCustomSettingsDb_DFTLOADRECSES_Range130000 = "range(1, 30000)"
)

// Constants associated with the CreateCustomSettingsDb.DFTPREFETCHSZ property.
// Configures the default prefetch size for queries.
const (
	CreateCustomSettingsDb_DFTPREFETCHSZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_DFTPREFETCHSZ_Range032767 = "range(0, 32767)"
)

// Constants associated with the CreateCustomSettingsDb.DFTQUERYOPT property.
// Configures the default query optimization level.
const (
	CreateCustomSettingsDb_DFTQUERYOPT_Range09 = "range(0, 9)"
)

// Constants associated with the CreateCustomSettingsDb.DFTSCHEMASDCC property.
// Configures whether DCC (database control center) is enabled for schemas.
const (
	CreateCustomSettingsDb_DFTSCHEMASDCC_No = "NO"
	CreateCustomSettingsDb_DFTSCHEMASDCC_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.DFTSQLMATHWARN property.
// Configures whether SQL math warnings are enabled.
const (
	CreateCustomSettingsDb_DFTSQLMATHWARN_No = "NO"
	CreateCustomSettingsDb_DFTSQLMATHWARN_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.DFTTABLEORG property.
// Configures the default table organization (ROW or COLUMN).
const (
	CreateCustomSettingsDb_DFTTABLEORG_Column = "COLUMN"
	CreateCustomSettingsDb_DFTTABLEORG_Row = "ROW"
)

// Constants associated with the CreateCustomSettingsDb.DLCHKTIME property.
// Configures the deadlock check time in milliseconds.
const (
	CreateCustomSettingsDb_DLCHKTIME_Range1000600000 = "range(1000, 600000)"
)

// Constants associated with the CreateCustomSettingsDb.ENABLEXMLCHAR property.
// Configures whether XML character support is enabled.
const (
	CreateCustomSettingsDb_ENABLEXMLCHAR_No = "NO"
	CreateCustomSettingsDb_ENABLEXMLCHAR_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.EXTENDEDROWSZ property.
// Configures whether extended row size is enabled.
const (
	CreateCustomSettingsDb_EXTENDEDROWSZ_Disable = "DISABLE"
	CreateCustomSettingsDb_EXTENDEDROWSZ_Enable = "ENABLE"
)

// Constants associated with the CreateCustomSettingsDb.GROUPHEAPRATIO property.
// Configures the heap ratio for group heap memory.
const (
	CreateCustomSettingsDb_GROUPHEAPRATIO_Range199 = "range(1, 99)"
)

// Constants associated with the CreateCustomSettingsDb.INDEXREC property.
// Configures the index recovery method.
const (
	CreateCustomSettingsDb_INDEXREC_Access = "ACCESS"
	CreateCustomSettingsDb_INDEXREC_AccessNoRedo = "ACCESS_NO_REDO"
	CreateCustomSettingsDb_INDEXREC_Restart = "RESTART"
	CreateCustomSettingsDb_INDEXREC_RestartNoRedo = "RESTART_NO_REDO"
	CreateCustomSettingsDb_INDEXREC_System = "SYSTEM"
)

// Constants associated with the CreateCustomSettingsDb.LARGEAGGREGATION property.
// Configures whether large aggregation is enabled.
const (
	CreateCustomSettingsDb_LARGEAGGREGATION_No = "NO"
	CreateCustomSettingsDb_LARGEAGGREGATION_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.LOCKLIST property.
// Configures the lock list memory size.
const (
	CreateCustomSettingsDb_LOCKLIST_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_LOCKLIST_Range4134217728 = "range(4, 134217728)"
)

// Constants associated with the CreateCustomSettingsDb.LOCKTIMEOUT property.
// Configures the lock timeout duration.
const (
	CreateCustomSettingsDb_LOCKTIMEOUT_Range032767 = "range(0, 32767)"
)

// Constants associated with the CreateCustomSettingsDb.LOGINDEXBUILD property.
// Configures whether index builds are logged.
const (
	CreateCustomSettingsDb_LOGINDEXBUILD_Off = "OFF"
	CreateCustomSettingsDb_LOGINDEXBUILD_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.LOGAPPLINFO property.
// Configures whether application information is logged.
const (
	CreateCustomSettingsDb_LOGAPPLINFO_No = "NO"
	CreateCustomSettingsDb_LOGAPPLINFO_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.LOGDDLSTMTS property.
// Configures whether DDL statements are logged.
const (
	CreateCustomSettingsDb_LOGDDLSTMTS_No = "NO"
	CreateCustomSettingsDb_LOGDDLSTMTS_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.LOGDISKCAP property.
// Configures the disk capacity log setting.
const (
	CreateCustomSettingsDb_LOGDISKCAP_Range12147483647 = "range(1, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.MAXAPPLS property.
// Configures the maximum number of applications.
const (
	CreateCustomSettingsDb_MAXAPPLS_Range160000 = "range(1, 60000)"
)

// Constants associated with the CreateCustomSettingsDb.MAXFILOP property.
// Configures the maximum number of file operations.
const (
	CreateCustomSettingsDb_MAXFILOP_Range6461440 = "range(64, 61440)"
)

// Constants associated with the CreateCustomSettingsDb.MAXLOCKS property.
// Configures the maximum number of locks.
const (
	CreateCustomSettingsDb_MAXLOCKS_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_MAXLOCKS_Range1100 = "range(1, 100)"
)

// Constants associated with the CreateCustomSettingsDb.MINDECDIV3 property.
// Configures whether decimal division by 3 should be handled.
const (
	CreateCustomSettingsDb_MINDECDIV3_No = "NO"
	CreateCustomSettingsDb_MINDECDIV3_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.MONACTMETRICS property.
// Configures the level of activity metrics to be monitored.
const (
	CreateCustomSettingsDb_MONACTMETRICS_Base = "BASE"
	CreateCustomSettingsDb_MONACTMETRICS_Extended = "EXTENDED"
	CreateCustomSettingsDb_MONACTMETRICS_None = "NONE"
)

// Constants associated with the CreateCustomSettingsDb.MONDEADLOCK property.
// Configures deadlock monitoring settings.
const (
	CreateCustomSettingsDb_MONDEADLOCK_HistAndValues = "HIST_AND_VALUES"
	CreateCustomSettingsDb_MONDEADLOCK_History = "HISTORY"
	CreateCustomSettingsDb_MONDEADLOCK_None = "NONE"
	CreateCustomSettingsDb_MONDEADLOCK_WithoutHist = "WITHOUT_HIST"
)

// Constants associated with the CreateCustomSettingsDb.MONLCKMSGLVL property.
// Configures the lock message level for monitoring.
const (
	CreateCustomSettingsDb_MONLCKMSGLVL_Range03 = "range(0, 3)"
)

// Constants associated with the CreateCustomSettingsDb.MONLOCKTIMEOUT property.
// Configures lock timeout monitoring settings.
const (
	CreateCustomSettingsDb_MONLOCKTIMEOUT_HistAndValues = "HIST_AND_VALUES"
	CreateCustomSettingsDb_MONLOCKTIMEOUT_History = "HISTORY"
	CreateCustomSettingsDb_MONLOCKTIMEOUT_None = "NONE"
	CreateCustomSettingsDb_MONLOCKTIMEOUT_WithoutHist = "WITHOUT_HIST"
)

// Constants associated with the CreateCustomSettingsDb.MONLOCKWAIT property.
// Configures lock wait monitoring settings.
const (
	CreateCustomSettingsDb_MONLOCKWAIT_HistAndValues = "HIST_AND_VALUES"
	CreateCustomSettingsDb_MONLOCKWAIT_History = "HISTORY"
	CreateCustomSettingsDb_MONLOCKWAIT_None = "NONE"
	CreateCustomSettingsDb_MONLOCKWAIT_WithoutHist = "WITHOUT_HIST"
)

// Constants associated with the CreateCustomSettingsDb.MONLWTHRESH property.
// Configures the lightweight threshold for monitoring.
const (
	CreateCustomSettingsDb_MONLWTHRESH_Range10004294967295 = "range(1000, 4294967295)"
)

// Constants associated with the CreateCustomSettingsDb.MONOBJMETRICS property.
// Configures the object metrics level for monitoring.
const (
	CreateCustomSettingsDb_MONOBJMETRICS_Base = "BASE"
	CreateCustomSettingsDb_MONOBJMETRICS_Extended = "EXTENDED"
	CreateCustomSettingsDb_MONOBJMETRICS_None = "NONE"
)

// Constants associated with the CreateCustomSettingsDb.MONPKGLISTSZ property.
// Configures the package list size for monitoring.
const (
	CreateCustomSettingsDb_MONPKGLISTSZ_Range01024 = "range(0, 1024)"
)

// Constants associated with the CreateCustomSettingsDb.MONREQMETRICS property.
// Configures the request metrics level for monitoring.
const (
	CreateCustomSettingsDb_MONREQMETRICS_Base = "BASE"
	CreateCustomSettingsDb_MONREQMETRICS_Extended = "EXTENDED"
	CreateCustomSettingsDb_MONREQMETRICS_None = "NONE"
)

// Constants associated with the CreateCustomSettingsDb.MONRTNDATA property.
// Configures the level of return data for monitoring.
const (
	CreateCustomSettingsDb_MONRTNDATA_Base = "BASE"
	CreateCustomSettingsDb_MONRTNDATA_None = "NONE"
)

// Constants associated with the CreateCustomSettingsDb.MONRTNEXECLIST property.
// Configures whether stored procedure execution list is monitored.
const (
	CreateCustomSettingsDb_MONRTNEXECLIST_Off = "OFF"
	CreateCustomSettingsDb_MONRTNEXECLIST_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.MONUOWDATA property.
// Configures the level of unit of work (UOW) data for monitoring.
const (
	CreateCustomSettingsDb_MONUOWDATA_Base = "BASE"
	CreateCustomSettingsDb_MONUOWDATA_None = "NONE"
)

// Constants associated with the CreateCustomSettingsDb.MONUOWEXECLIST property.
// Configures whether UOW execution list is monitored.
const (
	CreateCustomSettingsDb_MONUOWEXECLIST_Off = "OFF"
	CreateCustomSettingsDb_MONUOWEXECLIST_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.MONUOWPKGLIST property.
// Configures whether UOW package list is monitored.
const (
	CreateCustomSettingsDb_MONUOWPKGLIST_Off = "OFF"
	CreateCustomSettingsDb_MONUOWPKGLIST_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.NCHARMAPPING property.
// Configures the mapping of NCHAR character types.
const (
	CreateCustomSettingsDb_NCHARMAPPING_CharCu32 = "CHAR_CU32"
	CreateCustomSettingsDb_NCHARMAPPING_GraphicCu16 = "GRAPHIC_CU16"
	CreateCustomSettingsDb_NCHARMAPPING_GraphicCu32 = "GRAPHIC_CU32"
	CreateCustomSettingsDb_NCHARMAPPING_NotApplicable = "NOT APPLICABLE"
)

// Constants associated with the CreateCustomSettingsDb.NUMFREQVALUES property.
// Configures the number of frequent values for optimization.
const (
	CreateCustomSettingsDb_NUMFREQVALUES_Range032767 = "range(0, 32767)"
)

// Constants associated with the CreateCustomSettingsDb.NUMIOCLEANERS property.
// Configures the number of IO cleaners.
const (
	CreateCustomSettingsDb_NUMIOCLEANERS_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_NUMIOCLEANERS_Range0255 = "range(0, 255)"
)

// Constants associated with the CreateCustomSettingsDb.NUMIOSERVERS property.
// Configures the number of IO servers.
const (
	CreateCustomSettingsDb_NUMIOSERVERS_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_NUMIOSERVERS_Range1255 = "range(1, 255)"
)

// Constants associated with the CreateCustomSettingsDb.NUMLOGSPAN property.
// Configures the number of log spans.
const (
	CreateCustomSettingsDb_NUMLOGSPAN_Range065535 = "range(0, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.NUMQUANTILES property.
// Configures the number of quantiles for optimizations.
const (
	CreateCustomSettingsDb_NUMQUANTILES_Range032767 = "range(0, 32767)"
)

// Constants associated with the CreateCustomSettingsDb.OPTDIRECTWRKLD property.
// Configures the direct workload optimization setting.
const (
	CreateCustomSettingsDb_OPTDIRECTWRKLD_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_OPTDIRECTWRKLD_No = "NO"
	CreateCustomSettingsDb_OPTDIRECTWRKLD_Off = "OFF"
	CreateCustomSettingsDb_OPTDIRECTWRKLD_On = "ON"
	CreateCustomSettingsDb_OPTDIRECTWRKLD_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.PAGEAGETRGTGCR property.
// Configures the page age target for garbage collection.
const (
	CreateCustomSettingsDb_PAGEAGETRGTGCR_Range165535 = "range(1, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.PAGEAGETRGTMCR property.
// Configures the page age target for memory collection.
const (
	CreateCustomSettingsDb_PAGEAGETRGTMCR_Range165535 = "range(1, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.PCKCACHESZ property.
// Configures the package cache size.
const (
	CreateCustomSettingsDb_PCKCACHESZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_PCKCACHESZ_Range322147483646 = "range(32, 2147483646)"
)

// Constants associated with the CreateCustomSettingsDb.PLSTACKTRACE property.
// Configures the level of stack trace logging for stored procedures.
const (
	CreateCustomSettingsDb_PLSTACKTRACE_All = "ALL"
	CreateCustomSettingsDb_PLSTACKTRACE_None = "NONE"
	CreateCustomSettingsDb_PLSTACKTRACE_Unhandled = "UNHANDLED"
)

// Constants associated with the CreateCustomSettingsDb.SELFTUNINGMEM property.
// Configures whether self-tuning memory is enabled.
const (
	CreateCustomSettingsDb_SELFTUNINGMEM_Off = "OFF"
	CreateCustomSettingsDb_SELFTUNINGMEM_On = "ON"
)

// Constants associated with the CreateCustomSettingsDb.SEQDETECT property.
// Configures sequence detection for queries.
const (
	CreateCustomSettingsDb_SEQDETECT_No = "NO"
	CreateCustomSettingsDb_SEQDETECT_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.SHEAPTHRESSHR property.
// Configures the shared heap threshold size.
const (
	CreateCustomSettingsDb_SHEAPTHRESSHR_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_SHEAPTHRESSHR_Range2502147483647 = "range(250, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.SORTHEAP property.
// Configures the sort heap memory size.
const (
	CreateCustomSettingsDb_SORTHEAP_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_SORTHEAP_Range164294967295 = "range(16, 4294967295)"
)

// Constants associated with the CreateCustomSettingsDb.STATHEAPSZ property.
// Configures the statistics heap size.
const (
	CreateCustomSettingsDb_STATHEAPSZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_STATHEAPSZ_Range10962147483647 = "range(1096, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.STMTHEAP property.
// Configures the statement heap size.
const (
	CreateCustomSettingsDb_STMTHEAP_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_STMTHEAP_Range1282147483647 = "range(128, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.STMTCONC property.
// Configures the statement concurrency.
const (
	CreateCustomSettingsDb_STMTCONC_CommLit = "COMM_LIT"
	CreateCustomSettingsDb_STMTCONC_Comments = "COMMENTS"
	CreateCustomSettingsDb_STMTCONC_Literals = "LITERALS"
	CreateCustomSettingsDb_STMTCONC_Off = "OFF"
)

// Constants associated with the CreateCustomSettingsDb.STRINGUNITS property.
// Configures the string unit settings.
const (
	CreateCustomSettingsDb_STRINGUNITS_Codeunits32 = "CODEUNITS32"
	CreateCustomSettingsDb_STRINGUNITS_System = "SYSTEM"
)

// Constants associated with the CreateCustomSettingsDb.SYSTIMEPERIODADJ property.
// Configures whether system time period adjustments are enabled.
const (
	CreateCustomSettingsDb_SYSTIMEPERIODADJ_No = "NO"
	CreateCustomSettingsDb_SYSTIMEPERIODADJ_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.TRACKMOD property.
// Configures whether modifications to tracked objects are logged.
const (
	CreateCustomSettingsDb_TRACKMOD_No = "NO"
	CreateCustomSettingsDb_TRACKMOD_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.UTILHEAPSZ property.
// Configures the utility heap size.
const (
	CreateCustomSettingsDb_UTILHEAPSZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_UTILHEAPSZ_Range162147483647 = "range(16, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDb.WLMADMISSIONCTRL property.
// Configures whether WLM (Workload Management) admission control is enabled.
const (
	CreateCustomSettingsDb_WLMADMISSIONCTRL_No = "NO"
	CreateCustomSettingsDb_WLMADMISSIONCTRL_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDb.WLMAGENTLOADTRGT property.
// Configures the WLM agent load target.
const (
	CreateCustomSettingsDb_WLMAGENTLOADTRGT_Automatic = "AUTOMATIC"
	CreateCustomSettingsDb_WLMAGENTLOADTRGT_Range165535 = "range(1, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.WLMCPULIMIT property.
// Configures the CPU limit for WLM workloads.
const (
	CreateCustomSettingsDb_WLMCPULIMIT_Range0100 = "range(0, 100)"
)

// Constants associated with the CreateCustomSettingsDb.WLMCPUSHARES property.
// Configures the CPU share count for WLM workloads.
const (
	CreateCustomSettingsDb_WLMCPUSHARES_Range165535 = "range(1, 65535)"
)

// Constants associated with the CreateCustomSettingsDb.WLMCPUSHAREMODE property.
// Configures the mode of CPU shares for WLM workloads.
const (
	CreateCustomSettingsDb_WLMCPUSHAREMODE_Hard = "HARD"
	CreateCustomSettingsDb_WLMCPUSHAREMODE_Soft = "SOFT"
)

// UnmarshalCreateCustomSettingsDb unmarshals an instance of CreateCustomSettingsDb from the specified map of raw messages.
func UnmarshalCreateCustomSettingsDb(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateCustomSettingsDb)
	err = core.UnmarshalPrimitive(m, "ACT_SORTMEM_LIMIT", &obj.ACTSORTMEMLIMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "ACT_SORTMEM_LIMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ALT_COLLATE", &obj.ALTCOLLATE)
	if err != nil {
		err = core.SDKErrorf(err, "", "ALT_COLLATE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPGROUP_MEM_SZ", &obj.APPGROUPMEMSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPGROUP_MEM_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPLHEAPSZ", &obj.APPLHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPLHEAPSZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPL_MEMORY", &obj.APPLMEMORY)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPL_MEMORY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APP_CTL_HEAP_SZ", &obj.APPCTLHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APP_CTL_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ARCHRETRYDELAY", &obj.ARCHRETRYDELAY)
	if err != nil {
		err = core.SDKErrorf(err, "", "ARCHRETRYDELAY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTHN_CACHE_DURATION", &obj.AUTHNCACHEDURATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTHN_CACHE_DURATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTORESTART", &obj.AUTORESTART)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTORESTART-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_CG_STATS", &obj.AUTOCGSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_CG_STATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_MAINT", &obj.AUTOMAINT)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_MAINT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_REORG", &obj.AUTOREORG)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_REORG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_REVAL", &obj.AUTOREVAL)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_REVAL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_RUNSTATS", &obj.AUTORUNSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_RUNSTATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_SAMPLING", &obj.AUTOSAMPLING)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_SAMPLING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_STATS_VIEWS", &obj.AUTOSTATSVIEWS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_STATS_VIEWS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_STMT_STATS", &obj.AUTOSTMTSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_STMT_STATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_TBL_MAINT", &obj.AUTOTBLMAINT)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_TBL_MAINT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AVG_APPLS", &obj.AVGAPPLS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AVG_APPLS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CATALOGCACHE_SZ", &obj.CATALOGCACHESZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "CATALOGCACHE_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CHNGPGS_THRESH", &obj.CHNGPGSTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "CHNGPGS_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CUR_COMMIT", &obj.CURCOMMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "CUR_COMMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DATABASE_MEMORY", &obj.DATABASEMEMORY)
	if err != nil {
		err = core.SDKErrorf(err, "", "DATABASE_MEMORY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DBHEAP", &obj.DBHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "DBHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB_COLLNAME", &obj.DBCOLLNAME)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB_COLLNAME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB_MEM_THRESH", &obj.DBMEMTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB_MEM_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DDL_COMPRESSION_DEF", &obj.DDLCOMPRESSIONDEF)
	if err != nil {
		err = core.SDKErrorf(err, "", "DDL_COMPRESSION_DEF-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DDL_CONSTRAINT_DEF", &obj.DDLCONSTRAINTDEF)
	if err != nil {
		err = core.SDKErrorf(err, "", "DDL_CONSTRAINT_DEF-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DECFLT_ROUNDING", &obj.DECFLTROUNDING)
	if err != nil {
		err = core.SDKErrorf(err, "", "DECFLT_ROUNDING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DEC_ARITHMETIC", &obj.DECARITHMETIC)
	if err != nil {
		err = core.SDKErrorf(err, "", "DEC_ARITHMETIC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DEC_TO_CHAR_FMT", &obj.DECTOCHARFMT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DEC_TO_CHAR_FMT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_DEGREE", &obj.DFTDEGREE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_DEGREE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_EXTENT_SZ", &obj.DFTEXTENTSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_EXTENT_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_LOADREC_SES", &obj.DFTLOADRECSES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_LOADREC_SES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MTTB_TYPES", &obj.DFTMTTBTYPES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MTTB_TYPES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_PREFETCH_SZ", &obj.DFTPREFETCHSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_PREFETCH_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_QUERYOPT", &obj.DFTQUERYOPT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_QUERYOPT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_REFRESH_AGE", &obj.DFTREFRESHAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_REFRESH_AGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_SCHEMAS_DCC", &obj.DFTSCHEMASDCC)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_SCHEMAS_DCC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_SQLMATHWARN", &obj.DFTSQLMATHWARN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_SQLMATHWARN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_TABLE_ORG", &obj.DFTTABLEORG)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_TABLE_ORG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DLCHKTIME", &obj.DLCHKTIME)
	if err != nil {
		err = core.SDKErrorf(err, "", "DLCHKTIME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ENABLE_XMLCHAR", &obj.ENABLEXMLCHAR)
	if err != nil {
		err = core.SDKErrorf(err, "", "ENABLE_XMLCHAR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "EXTENDED_ROW_SZ", &obj.EXTENDEDROWSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "EXTENDED_ROW_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "GROUPHEAP_RATIO", &obj.GROUPHEAPRATIO)
	if err != nil {
		err = core.SDKErrorf(err, "", "GROUPHEAP_RATIO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INDEXREC", &obj.INDEXREC)
	if err != nil {
		err = core.SDKErrorf(err, "", "INDEXREC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LARGE_AGGREGATION", &obj.LARGEAGGREGATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "LARGE_AGGREGATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOCKLIST", &obj.LOCKLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOCKLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOCKTIMEOUT", &obj.LOCKTIMEOUT)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOCKTIMEOUT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOGINDEXBUILD", &obj.LOGINDEXBUILD)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOGINDEXBUILD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_APPL_INFO", &obj.LOGAPPLINFO)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_APPL_INFO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_DDL_STMTS", &obj.LOGDDLSTMTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_DDL_STMTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_DISK_CAP", &obj.LOGDISKCAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_DISK_CAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXAPPLS", &obj.MAXAPPLS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXAPPLS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXFILOP", &obj.MAXFILOP)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXFILOP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXLOCKS", &obj.MAXLOCKS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXLOCKS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MIN_DEC_DIV_3", &obj.MINDECDIV3)
	if err != nil {
		err = core.SDKErrorf(err, "", "MIN_DEC_DIV_3-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_ACT_METRICS", &obj.MONACTMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_ACT_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_DEADLOCK", &obj.MONDEADLOCK)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_DEADLOCK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LCK_MSG_LVL", &obj.MONLCKMSGLVL)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LCK_MSG_LVL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LOCKTIMEOUT", &obj.MONLOCKTIMEOUT)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LOCKTIMEOUT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LOCKWAIT", &obj.MONLOCKWAIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LOCKWAIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LW_THRESH", &obj.MONLWTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LW_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_OBJ_METRICS", &obj.MONOBJMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_OBJ_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_PKGLIST_SZ", &obj.MONPKGLISTSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_PKGLIST_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_REQ_METRICS", &obj.MONREQMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_REQ_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_RTN_DATA", &obj.MONRTNDATA)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_RTN_DATA-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_RTN_EXECLIST", &obj.MONRTNEXECLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_RTN_EXECLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_DATA", &obj.MONUOWDATA)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_DATA-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_EXECLIST", &obj.MONUOWEXECLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_EXECLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_PKGLIST", &obj.MONUOWPKGLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_PKGLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NCHAR_MAPPING", &obj.NCHARMAPPING)
	if err != nil {
		err = core.SDKErrorf(err, "", "NCHAR_MAPPING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_FREQVALUES", &obj.NUMFREQVALUES)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_FREQVALUES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_IOCLEANERS", &obj.NUMIOCLEANERS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_IOCLEANERS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_IOSERVERS", &obj.NUMIOSERVERS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_IOSERVERS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_LOG_SPAN", &obj.NUMLOGSPAN)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_LOG_SPAN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_QUANTILES", &obj.NUMQUANTILES)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_QUANTILES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_BUFFPAGE", &obj.OPTBUFFPAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_BUFFPAGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_DIRECT_WRKLD", &obj.OPTDIRECTWRKLD)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_DIRECT_WRKLD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_LOCKLIST", &obj.OPTLOCKLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_LOCKLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_MAXLOCKS", &obj.OPTMAXLOCKS)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_MAXLOCKS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_SORTHEAP", &obj.OPTSORTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_SORTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PAGE_AGE_TRGT_GCR", &obj.PAGEAGETRGTGCR)
	if err != nil {
		err = core.SDKErrorf(err, "", "PAGE_AGE_TRGT_GCR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PAGE_AGE_TRGT_MCR", &obj.PAGEAGETRGTMCR)
	if err != nil {
		err = core.SDKErrorf(err, "", "PAGE_AGE_TRGT_MCR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PCKCACHESZ", &obj.PCKCACHESZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "PCKCACHESZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PL_STACK_TRACE", &obj.PLSTACKTRACE)
	if err != nil {
		err = core.SDKErrorf(err, "", "PL_STACK_TRACE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SELF_TUNING_MEM", &obj.SELFTUNINGMEM)
	if err != nil {
		err = core.SDKErrorf(err, "", "SELF_TUNING_MEM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SEQDETECT", &obj.SEQDETECT)
	if err != nil {
		err = core.SDKErrorf(err, "", "SEQDETECT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SHEAPTHRES_SHR", &obj.SHEAPTHRESSHR)
	if err != nil {
		err = core.SDKErrorf(err, "", "SHEAPTHRES_SHR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SOFTMAX", &obj.SOFTMAX)
	if err != nil {
		err = core.SDKErrorf(err, "", "SOFTMAX-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SORTHEAP", &obj.SORTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "SORTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SQL_CCFLAGS", &obj.SQLCCFLAGS)
	if err != nil {
		err = core.SDKErrorf(err, "", "SQL_CCFLAGS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STAT_HEAP_SZ", &obj.STATHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "STAT_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STMTHEAP", &obj.STMTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "STMTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STMT_CONC", &obj.STMTCONC)
	if err != nil {
		err = core.SDKErrorf(err, "", "STMT_CONC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STRING_UNITS", &obj.STRINGUNITS)
	if err != nil {
		err = core.SDKErrorf(err, "", "STRING_UNITS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SYSTIME_PERIOD_ADJ", &obj.SYSTIMEPERIODADJ)
	if err != nil {
		err = core.SDKErrorf(err, "", "SYSTIME_PERIOD_ADJ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "TRACKMOD", &obj.TRACKMOD)
	if err != nil {
		err = core.SDKErrorf(err, "", "TRACKMOD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "UTIL_HEAP_SZ", &obj.UTILHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "UTIL_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_ADMISSION_CTRL", &obj.WLMADMISSIONCTRL)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_ADMISSION_CTRL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_AGENT_LOAD_TRGT", &obj.WLMAGENTLOADTRGT)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_AGENT_LOAD_TRGT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_LIMIT", &obj.WLMCPULIMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_LIMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_SHARES", &obj.WLMCPUSHARES)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_SHARES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_SHARE_MODE", &obj.WLMCPUSHAREMODE)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_SHARE_MODE-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCustomSettingsDbm : Container for general database management settings.
type CreateCustomSettingsDbm struct {
	// Configures the communication bandwidth for the database manager.
	COMMBANDWIDTH *string `json:"COMM_BANDWIDTH,omitempty"`

	// Configures the CPU speed for the database manager.
	CPUSPEED *string `json:"CPUSPEED,omitempty"`

	// Configures whether the buffer pool is monitored by default.
	DFTMONBUFPOOL *string `json:"DFT_MON_BUFPOOL,omitempty"`

	// Configures whether lock monitoring is enabled by default.
	DFTMONLOCK *string `json:"DFT_MON_LOCK,omitempty"`

	// Configures whether sort operations are monitored by default.
	DFTMONSORT *string `json:"DFT_MON_SORT,omitempty"`

	// Configures whether statement execution is monitored by default.
	DFTMONSTMT *string `json:"DFT_MON_STMT,omitempty"`

	// Configures whether table operations are monitored by default.
	DFTMONTABLE *string `json:"DFT_MON_TABLE,omitempty"`

	// Configures whether timestamp monitoring is enabled by default.
	DFTMONTIMESTAMP *string `json:"DFT_MON_TIMESTAMP,omitempty"`

	// Configures whether unit of work (UOW) monitoring is enabled by default.
	DFTMONUOW *string `json:"DFT_MON_UOW,omitempty"`

	// Configures the diagnostic level for the database manager.
	DIAGLEVEL *string `json:"DIAGLEVEL,omitempty"`

	// Configures whether federated asynchronous mode is enabled.
	FEDERATEDASYNC *string `json:"FEDERATED_ASYNC,omitempty"`

	// Configures the type of indexing to be used in the database manager.
	INDEXREC *string `json:"INDEXREC,omitempty"`

	// Configures the parallelism settings for intra-query parallelism.
	INTRAPARALLEL *string `json:"INTRA_PARALLEL,omitempty"`

	// Configures whether fenced routines are kept in memory.
	KEEPFENCED *string `json:"KEEPFENCED,omitempty"`

	// Configures the maximum number of connection retries.
	MAXCONNRETRIES *string `json:"MAX_CONNRETRIES,omitempty"`

	// Configures the maximum degree of parallelism for queries.
	MAXQUERYDEGREE *string `json:"MAX_QUERYDEGREE,omitempty"`

	// Configures the size of the monitoring heap.
	MONHEAPSZ *string `json:"MON_HEAP_SZ,omitempty"`

	// Configures the size of multipart queries in MB.
	MULTIPARTSIZEMB *string `json:"MULTIPARTSIZEMB,omitempty"`

	// Configures the level of notifications for the database manager.
	NOTIFYLEVEL *string `json:"NOTIFYLEVEL,omitempty"`

	// Configures the number of initial agents in the database manager.
	NUMINITAGENTS *string `json:"NUM_INITAGENTS,omitempty"`

	// Configures the number of initial fenced routines.
	NUMINITFENCED *string `json:"NUM_INITFENCED,omitempty"`

	// Configures the number of pool agents.
	NUMPOOLAGENTS *string `json:"NUM_POOLAGENTS,omitempty"`

	// Configures the interval between resync operations.
	RESYNCINTERVAL *string `json:"RESYNC_INTERVAL,omitempty"`

	// Configures the request/response I/O block size.
	RQRIOBLK *string `json:"RQRIOBLK,omitempty"`

	// Configures the time in minutes for start/stop operations.
	STARTSTOPTIME *string `json:"START_STOP_TIME,omitempty"`

	// Configures the utility impact limit.
	UTILIMPACTLIM *string `json:"UTIL_IMPACT_LIM,omitempty"`

	// Configures whether the WLM (Workload Management) dispatcher is enabled.
	WLMDISPATCHER *string `json:"WLM_DISPATCHER,omitempty"`

	// Configures the concurrency level for the WLM dispatcher.
	WLMDISPCONCUR *string `json:"WLM_DISP_CONCUR,omitempty"`

	// Configures whether CPU shares are used for WLM dispatcher.
	WLMDISPCPUSHARES *string `json:"WLM_DISP_CPU_SHARES,omitempty"`

	// Configures the minimum utility threshold for WLM dispatcher.
	WLMDISPMINUTIL *string `json:"WLM_DISP_MIN_UTIL,omitempty"`
}

// Constants associated with the CreateCustomSettingsDbm.COMMBANDWIDTH property.
// Configures the communication bandwidth for the database manager.
const (
	CreateCustomSettingsDbm_COMMBANDWIDTH_Range01100000 = "range(0.1, 100000)"
)

// Constants associated with the CreateCustomSettingsDbm.CPUSPEED property.
// Configures the CPU speed for the database manager.
const (
	CreateCustomSettingsDbm_CPUSPEED_Range000000000011 = "range(0.0000000001, 1)"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONBUFPOOL property.
// Configures whether the buffer pool is monitored by default.
const (
	CreateCustomSettingsDbm_DFTMONBUFPOOL_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONBUFPOOL_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONLOCK property.
// Configures whether lock monitoring is enabled by default.
const (
	CreateCustomSettingsDbm_DFTMONLOCK_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONLOCK_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONSORT property.
// Configures whether sort operations are monitored by default.
const (
	CreateCustomSettingsDbm_DFTMONSORT_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONSORT_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONSTMT property.
// Configures whether statement execution is monitored by default.
const (
	CreateCustomSettingsDbm_DFTMONSTMT_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONSTMT_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONTABLE property.
// Configures whether table operations are monitored by default.
const (
	CreateCustomSettingsDbm_DFTMONTABLE_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONTABLE_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONTIMESTAMP property.
// Configures whether timestamp monitoring is enabled by default.
const (
	CreateCustomSettingsDbm_DFTMONTIMESTAMP_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONTIMESTAMP_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DFTMONUOW property.
// Configures whether unit of work (UOW) monitoring is enabled by default.
const (
	CreateCustomSettingsDbm_DFTMONUOW_Off = "OFF"
	CreateCustomSettingsDbm_DFTMONUOW_On = "ON"
)

// Constants associated with the CreateCustomSettingsDbm.DIAGLEVEL property.
// Configures the diagnostic level for the database manager.
const (
	CreateCustomSettingsDbm_DIAGLEVEL_Range04 = "range(0, 4)"
)

// Constants associated with the CreateCustomSettingsDbm.FEDERATEDASYNC property.
// Configures whether federated asynchronous mode is enabled.
const (
	CreateCustomSettingsDbm_FEDERATEDASYNC_Any = "ANY"
	CreateCustomSettingsDbm_FEDERATEDASYNC_Range032767 = "range(0, 32767)"
)

// Constants associated with the CreateCustomSettingsDbm.INDEXREC property.
// Configures the type of indexing to be used in the database manager.
const (
	CreateCustomSettingsDbm_INDEXREC_Access = "ACCESS"
	CreateCustomSettingsDbm_INDEXREC_AccessNoRedo = "ACCESS_NO_REDO"
	CreateCustomSettingsDbm_INDEXREC_Restart = "RESTART"
	CreateCustomSettingsDbm_INDEXREC_RestartNoRedo = "RESTART_NO_REDO"
)

// Constants associated with the CreateCustomSettingsDbm.INTRAPARALLEL property.
// Configures the parallelism settings for intra-query parallelism.
const (
	CreateCustomSettingsDbm_INTRAPARALLEL_No = "NO"
	CreateCustomSettingsDbm_INTRAPARALLEL_System = "SYSTEM"
	CreateCustomSettingsDbm_INTRAPARALLEL_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDbm.KEEPFENCED property.
// Configures whether fenced routines are kept in memory.
const (
	CreateCustomSettingsDbm_KEEPFENCED_No = "NO"
	CreateCustomSettingsDbm_KEEPFENCED_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDbm.MAXCONNRETRIES property.
// Configures the maximum number of connection retries.
const (
	CreateCustomSettingsDbm_MAXCONNRETRIES_Range0100 = "range(0, 100)"
)

// Constants associated with the CreateCustomSettingsDbm.MAXQUERYDEGREE property.
// Configures the maximum degree of parallelism for queries.
const (
	CreateCustomSettingsDbm_MAXQUERYDEGREE_Any = "ANY"
	CreateCustomSettingsDbm_MAXQUERYDEGREE_Range132767 = "range(1, 32767)"
)

// Constants associated with the CreateCustomSettingsDbm.MONHEAPSZ property.
// Configures the size of the monitoring heap.
const (
	CreateCustomSettingsDbm_MONHEAPSZ_Automatic = "AUTOMATIC"
	CreateCustomSettingsDbm_MONHEAPSZ_Range02147483647 = "range(0, 2147483647)"
)

// Constants associated with the CreateCustomSettingsDbm.MULTIPARTSIZEMB property.
// Configures the size of multipart queries in MB.
const (
	CreateCustomSettingsDbm_MULTIPARTSIZEMB_Range55120 = "range(5, 5120)"
)

// Constants associated with the CreateCustomSettingsDbm.NOTIFYLEVEL property.
// Configures the level of notifications for the database manager.
const (
	CreateCustomSettingsDbm_NOTIFYLEVEL_Range04 = "range(0, 4)"
)

// Constants associated with the CreateCustomSettingsDbm.NUMINITAGENTS property.
// Configures the number of initial agents in the database manager.
const (
	CreateCustomSettingsDbm_NUMINITAGENTS_Range064000 = "range(0, 64000)"
)

// Constants associated with the CreateCustomSettingsDbm.NUMINITFENCED property.
// Configures the number of initial fenced routines.
const (
	CreateCustomSettingsDbm_NUMINITFENCED_Range064000 = "range(0, 64000)"
)

// Constants associated with the CreateCustomSettingsDbm.NUMPOOLAGENTS property.
// Configures the number of pool agents.
const (
	CreateCustomSettingsDbm_NUMPOOLAGENTS_Range064000 = "range(0, 64000)"
)

// Constants associated with the CreateCustomSettingsDbm.RESYNCINTERVAL property.
// Configures the interval between resync operations.
const (
	CreateCustomSettingsDbm_RESYNCINTERVAL_Range160000 = "range(1, 60000)"
)

// Constants associated with the CreateCustomSettingsDbm.RQRIOBLK property.
// Configures the request/response I/O block size.
const (
	CreateCustomSettingsDbm_RQRIOBLK_Range409665535 = "range(4096, 65535)"
)

// Constants associated with the CreateCustomSettingsDbm.STARTSTOPTIME property.
// Configures the time in minutes for start/stop operations.
const (
	CreateCustomSettingsDbm_STARTSTOPTIME_Range11440 = "range(1, 1440)"
)

// Constants associated with the CreateCustomSettingsDbm.UTILIMPACTLIM property.
// Configures the utility impact limit.
const (
	CreateCustomSettingsDbm_UTILIMPACTLIM_Range1100 = "range(1, 100)"
)

// Constants associated with the CreateCustomSettingsDbm.WLMDISPATCHER property.
// Configures whether the WLM (Workload Management) dispatcher is enabled.
const (
	CreateCustomSettingsDbm_WLMDISPATCHER_No = "NO"
	CreateCustomSettingsDbm_WLMDISPATCHER_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDbm.WLMDISPCONCUR property.
// Configures the concurrency level for the WLM dispatcher.
const (
	CreateCustomSettingsDbm_WLMDISPCONCUR_Computed = "COMPUTED"
	CreateCustomSettingsDbm_WLMDISPCONCUR_Range132767 = "range(1, 32767)"
)

// Constants associated with the CreateCustomSettingsDbm.WLMDISPCPUSHARES property.
// Configures whether CPU shares are used for WLM dispatcher.
const (
	CreateCustomSettingsDbm_WLMDISPCPUSHARES_No = "NO"
	CreateCustomSettingsDbm_WLMDISPCPUSHARES_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsDbm.WLMDISPMINUTIL property.
// Configures the minimum utility threshold for WLM dispatcher.
const (
	CreateCustomSettingsDbm_WLMDISPMINUTIL_Range0100 = "range(0, 100)"
)

// UnmarshalCreateCustomSettingsDbm unmarshals an instance of CreateCustomSettingsDbm from the specified map of raw messages.
func UnmarshalCreateCustomSettingsDbm(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateCustomSettingsDbm)
	err = core.UnmarshalPrimitive(m, "COMM_BANDWIDTH", &obj.COMMBANDWIDTH)
	if err != nil {
		err = core.SDKErrorf(err, "", "COMM_BANDWIDTH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CPUSPEED", &obj.CPUSPEED)
	if err != nil {
		err = core.SDKErrorf(err, "", "CPUSPEED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_BUFPOOL", &obj.DFTMONBUFPOOL)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_BUFPOOL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_LOCK", &obj.DFTMONLOCK)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_LOCK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_SORT", &obj.DFTMONSORT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_SORT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_STMT", &obj.DFTMONSTMT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_STMT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_TABLE", &obj.DFTMONTABLE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_TABLE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_TIMESTAMP", &obj.DFTMONTIMESTAMP)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_TIMESTAMP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_UOW", &obj.DFTMONUOW)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_UOW-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DIAGLEVEL", &obj.DIAGLEVEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "DIAGLEVEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "FEDERATED_ASYNC", &obj.FEDERATEDASYNC)
	if err != nil {
		err = core.SDKErrorf(err, "", "FEDERATED_ASYNC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INDEXREC", &obj.INDEXREC)
	if err != nil {
		err = core.SDKErrorf(err, "", "INDEXREC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INTRA_PARALLEL", &obj.INTRAPARALLEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "INTRA_PARALLEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "KEEPFENCED", &obj.KEEPFENCED)
	if err != nil {
		err = core.SDKErrorf(err, "", "KEEPFENCED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAX_CONNRETRIES", &obj.MAXCONNRETRIES)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAX_CONNRETRIES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAX_QUERYDEGREE", &obj.MAXQUERYDEGREE)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAX_QUERYDEGREE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_HEAP_SZ", &obj.MONHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MULTIPARTSIZEMB", &obj.MULTIPARTSIZEMB)
	if err != nil {
		err = core.SDKErrorf(err, "", "MULTIPARTSIZEMB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NOTIFYLEVEL", &obj.NOTIFYLEVEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "NOTIFYLEVEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_INITAGENTS", &obj.NUMINITAGENTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_INITAGENTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_INITFENCED", &obj.NUMINITFENCED)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_INITFENCED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_POOLAGENTS", &obj.NUMPOOLAGENTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_POOLAGENTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "RESYNC_INTERVAL", &obj.RESYNCINTERVAL)
	if err != nil {
		err = core.SDKErrorf(err, "", "RESYNC_INTERVAL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "RQRIOBLK", &obj.RQRIOBLK)
	if err != nil {
		err = core.SDKErrorf(err, "", "RQRIOBLK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "START_STOP_TIME", &obj.STARTSTOPTIME)
	if err != nil {
		err = core.SDKErrorf(err, "", "START_STOP_TIME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "UTIL_IMPACT_LIM", &obj.UTILIMPACTLIM)
	if err != nil {
		err = core.SDKErrorf(err, "", "UTIL_IMPACT_LIM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISPATCHER", &obj.WLMDISPATCHER)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISPATCHER-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_CONCUR", &obj.WLMDISPCONCUR)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_CONCUR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_CPU_SHARES", &obj.WLMDISPCPUSHARES)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_CPU_SHARES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_MIN_UTIL", &obj.WLMDISPMINUTIL)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_MIN_UTIL-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCustomSettingsRegistry : registry for db2 related configuration settings/configurations.
type CreateCustomSettingsRegistry struct {
	// Configures the bidi (bidirectional) support for DB2.
	DB2BIDI *string `json:"DB2BIDI,omitempty"`

	// Configures the DB2 component options (not specified in values).
	DB2COMPOPT *string `json:"DB2COMPOPT,omitempty"`

	// Configures the DB2 lock timeout behavior.
	DB2LOCKTORB *string `json:"DB2LOCK_TO_RB,omitempty"`

	// Configures whether DB2's self-tuning memory manager (STMM) is enabled.
	DB2STMM *string `json:"DB2STMM,omitempty"`

	// Configures the alternate authorization behavior for DB2.
	DB2ALTERNATEAUTHZBEHAVIOUR *string `json:"DB2_ALTERNATE_AUTHZ_BEHAVIOUR,omitempty"`

	// Configures how DB2 handles anti-joins.
	DB2ANTIJOIN *string `json:"DB2_ANTIJOIN,omitempty"`

	// Configures whether DB2 asynchronous table scanning (ATS) is enabled.
	DB2ATSENABLE *string `json:"DB2_ATS_ENABLE,omitempty"`

	// Configures whether deferred prepare semantics are enabled in DB2.
	DB2DEFERREDPREPARESEMANTICS *string `json:"DB2_DEFERRED_PREPARE_SEMANTICS,omitempty"`

	// Configures whether uncommitted data is evaluated by DB2.
	DB2EVALUNCOMMITTED *string `json:"DB2_EVALUNCOMMITTED,omitempty"`

	// Configures extended optimization in DB2 (not specified in values).
	DB2EXTENDEDOPTIMIZATION *string `json:"DB2_EXTENDED_OPTIMIZATION,omitempty"`

	// Configures the default percentage of free space for DB2 indexes.
	DB2INDEXPCTFREEDEFAULT *string `json:"DB2_INDEX_PCTFREE_DEFAULT,omitempty"`

	// Configures whether in-list queries are converted to nested loop joins.
	DB2INLISTTONLJN *string `json:"DB2_INLIST_TO_NLJN,omitempty"`

	// Configures whether DB2 minimizes list prefetching for queries.
	DB2MINIMIZELISTPREFETCH *string `json:"DB2_MINIMIZE_LISTPREFETCH,omitempty"`

	// Configures the number of entries for DB2 object tables.
	DB2OBJECTTABLEENTRIES *string `json:"DB2_OBJECT_TABLE_ENTRIES,omitempty"`

	// Configures whether DB2's optimizer profile is enabled.
	DB2OPTPROFILE *string `json:"DB2_OPTPROFILE,omitempty"`

	// Configures the logging of optimizer statistics (not specified in values).
	DB2OPTSTATSLOG *string `json:"DB2_OPTSTATS_LOG,omitempty"`

	// Configures the maximum temporary space size for DB2 optimizer.
	DB2OPTMAXTEMPSIZE *string `json:"DB2_OPT_MAX_TEMP_SIZE,omitempty"`

	// Configures parallel I/O behavior in DB2 (not specified in values).
	DB2PARALLELIO *string `json:"DB2_PARALLEL_IO,omitempty"`

	// Configures whether reduced optimization is applied in DB2 (not specified in values).
	DB2REDUCEDOPTIMIZATION *string `json:"DB2_REDUCED_OPTIMIZATION,omitempty"`

	// Configures the selectivity behavior for DB2 queries.
	DB2SELECTIVITY *string `json:"DB2_SELECTIVITY,omitempty"`

	// Configures whether DB2 skips deleted rows during query processing.
	DB2SKIPDELETED *string `json:"DB2_SKIPDELETED,omitempty"`

	// Configures whether DB2 skips inserted rows during query processing.
	DB2SKIPINSERTED *string `json:"DB2_SKIPINSERTED,omitempty"`

	// Configures whether DB2 synchronizes lock release attributes.
	DB2SYNCRELEASELOCKATTRIBUTES *string `json:"DB2_SYNC_RELEASE_LOCK_ATTRIBUTES,omitempty"`

	// Configures the types of operations that reuse storage after truncation.
	DB2TRUNCATEREUSESTORAGE *string `json:"DB2_TRUNCATE_REUSESTORAGE,omitempty"`

	// Configures whether DB2 uses alternate page cleaning methods.
	DB2USEALTERNATEPAGECLEANING *string `json:"DB2_USE_ALTERNATE_PAGE_CLEANING,omitempty"`

	// Configures whether DB2 view reoptimization values are used.
	DB2VIEWREOPTVALUES *string `json:"DB2_VIEW_REOPT_VALUES,omitempty"`

	// Configures the WLM (Workload Management) settings for DB2 (not specified in values).
	DB2WLMSETTINGS *string `json:"DB2_WLM_SETTINGS,omitempty"`

	// Configures the DB2 workload type.
	DB2WORKLOAD *string `json:"DB2_WORKLOAD,omitempty"`
}

// Constants associated with the CreateCustomSettingsRegistry.DB2BIDI property.
// Configures the bidi (bidirectional) support for DB2.
const (
	CreateCustomSettingsRegistry_DB2BIDI_No = "NO"
	CreateCustomSettingsRegistry_DB2BIDI_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2LOCKTORB property.
// Configures the DB2 lock timeout behavior.
const (
	CreateCustomSettingsRegistry_DB2LOCKTORB_Statement = "STATEMENT"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2STMM property.
// Configures whether DB2's self-tuning memory manager (STMM) is enabled.
const (
	CreateCustomSettingsRegistry_DB2STMM_No = "NO"
	CreateCustomSettingsRegistry_DB2STMM_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2ALTERNATEAUTHZBEHAVIOUR property.
// Configures the alternate authorization behavior for DB2.
const (
	CreateCustomSettingsRegistry_DB2ALTERNATEAUTHZBEHAVIOUR_ExternalRoutineDbadm = "EXTERNAL_ROUTINE_DBADM"
	CreateCustomSettingsRegistry_DB2ALTERNATEAUTHZBEHAVIOUR_ExternalRoutineDbauth = "EXTERNAL_ROUTINE_DBAUTH"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2ANTIJOIN property.
// Configures how DB2 handles anti-joins.
const (
	CreateCustomSettingsRegistry_DB2ANTIJOIN_Extend = "EXTEND"
	CreateCustomSettingsRegistry_DB2ANTIJOIN_No = "NO"
	CreateCustomSettingsRegistry_DB2ANTIJOIN_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2ATSENABLE property.
// Configures whether DB2 asynchronous table scanning (ATS) is enabled.
const (
	CreateCustomSettingsRegistry_DB2ATSENABLE_No = "NO"
	CreateCustomSettingsRegistry_DB2ATSENABLE_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2DEFERREDPREPARESEMANTICS property.
// Configures whether deferred prepare semantics are enabled in DB2.
const (
	CreateCustomSettingsRegistry_DB2DEFERREDPREPARESEMANTICS_No = "NO"
	CreateCustomSettingsRegistry_DB2DEFERREDPREPARESEMANTICS_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2EVALUNCOMMITTED property.
// Configures whether uncommitted data is evaluated by DB2.
const (
	CreateCustomSettingsRegistry_DB2EVALUNCOMMITTED_No = "NO"
	CreateCustomSettingsRegistry_DB2EVALUNCOMMITTED_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2INDEXPCTFREEDEFAULT property.
// Configures the default percentage of free space for DB2 indexes.
const (
	CreateCustomSettingsRegistry_DB2INDEXPCTFREEDEFAULT_Range099 = "range(0, 99)"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2INLISTTONLJN property.
// Configures whether in-list queries are converted to nested loop joins.
const (
	CreateCustomSettingsRegistry_DB2INLISTTONLJN_No = "NO"
	CreateCustomSettingsRegistry_DB2INLISTTONLJN_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2MINIMIZELISTPREFETCH property.
// Configures whether DB2 minimizes list prefetching for queries.
const (
	CreateCustomSettingsRegistry_DB2MINIMIZELISTPREFETCH_No = "NO"
	CreateCustomSettingsRegistry_DB2MINIMIZELISTPREFETCH_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2OBJECTTABLEENTRIES property.
// Configures the number of entries for DB2 object tables.
const (
	CreateCustomSettingsRegistry_DB2OBJECTTABLEENTRIES_Range065532 = "range(0, 65532)"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2OPTPROFILE property.
// Configures whether DB2's optimizer profile is enabled.
const (
	CreateCustomSettingsRegistry_DB2OPTPROFILE_No = "NO"
	CreateCustomSettingsRegistry_DB2OPTPROFILE_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2SELECTIVITY property.
// Configures the selectivity behavior for DB2 queries.
const (
	CreateCustomSettingsRegistry_DB2SELECTIVITY_All = "ALL"
	CreateCustomSettingsRegistry_DB2SELECTIVITY_No = "NO"
	CreateCustomSettingsRegistry_DB2SELECTIVITY_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2SKIPDELETED property.
// Configures whether DB2 skips deleted rows during query processing.
const (
	CreateCustomSettingsRegistry_DB2SKIPDELETED_No = "NO"
	CreateCustomSettingsRegistry_DB2SKIPDELETED_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2SKIPINSERTED property.
// Configures whether DB2 skips inserted rows during query processing.
const (
	CreateCustomSettingsRegistry_DB2SKIPINSERTED_No = "NO"
	CreateCustomSettingsRegistry_DB2SKIPINSERTED_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2SYNCRELEASELOCKATTRIBUTES property.
// Configures whether DB2 synchronizes lock release attributes.
const (
	CreateCustomSettingsRegistry_DB2SYNCRELEASELOCKATTRIBUTES_No = "NO"
	CreateCustomSettingsRegistry_DB2SYNCRELEASELOCKATTRIBUTES_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2TRUNCATEREUSESTORAGE property.
// Configures the types of operations that reuse storage after truncation.
const (
	CreateCustomSettingsRegistry_DB2TRUNCATEREUSESTORAGE_Import = "IMPORT"
	CreateCustomSettingsRegistry_DB2TRUNCATEREUSESTORAGE_Load = "LOAD"
	CreateCustomSettingsRegistry_DB2TRUNCATEREUSESTORAGE_Truncate = "TRUNCATE"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2USEALTERNATEPAGECLEANING property.
// Configures whether DB2 uses alternate page cleaning methods.
const (
	CreateCustomSettingsRegistry_DB2USEALTERNATEPAGECLEANING_Off = "OFF"
	CreateCustomSettingsRegistry_DB2USEALTERNATEPAGECLEANING_On = "ON"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2VIEWREOPTVALUES property.
// Configures whether DB2 view reoptimization values are used.
const (
	CreateCustomSettingsRegistry_DB2VIEWREOPTVALUES_No = "NO"
	CreateCustomSettingsRegistry_DB2VIEWREOPTVALUES_Yes = "YES"
)

// Constants associated with the CreateCustomSettingsRegistry.DB2WORKLOAD property.
// Configures the DB2 workload type.
const (
	CreateCustomSettingsRegistry_DB2WORKLOAD_Analytics = "ANALYTICS"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Cm = "CM"
	CreateCustomSettingsRegistry_DB2WORKLOAD_CognosCs = "COGNOS_CS"
	CreateCustomSettingsRegistry_DB2WORKLOAD_FilenetCm = "FILENET_CM"
	CreateCustomSettingsRegistry_DB2WORKLOAD_InforErpLn = "INFOR_ERP_LN"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Maximo = "MAXIMO"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Mdm = "MDM"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Sap = "SAP"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Tpm = "TPM"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Was = "WAS"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Wc = "WC"
	CreateCustomSettingsRegistry_DB2WORKLOAD_Wp = "WP"
)

// UnmarshalCreateCustomSettingsRegistry unmarshals an instance of CreateCustomSettingsRegistry from the specified map of raw messages.
func UnmarshalCreateCustomSettingsRegistry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateCustomSettingsRegistry)
	err = core.UnmarshalPrimitive(m, "DB2BIDI", &obj.DB2BIDI)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2BIDI-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2COMPOPT", &obj.DB2COMPOPT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2COMPOPT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2LOCK_TO_RB", &obj.DB2LOCKTORB)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2LOCK_TO_RB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2STMM", &obj.DB2STMM)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2STMM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ALTERNATE_AUTHZ_BEHAVIOUR", &obj.DB2ALTERNATEAUTHZBEHAVIOUR)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ALTERNATE_AUTHZ_BEHAVIOUR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ANTIJOIN", &obj.DB2ANTIJOIN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ANTIJOIN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ATS_ENABLE", &obj.DB2ATSENABLE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ATS_ENABLE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_DEFERRED_PREPARE_SEMANTICS", &obj.DB2DEFERREDPREPARESEMANTICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_DEFERRED_PREPARE_SEMANTICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_EVALUNCOMMITTED", &obj.DB2EVALUNCOMMITTED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_EVALUNCOMMITTED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_EXTENDED_OPTIMIZATION", &obj.DB2EXTENDEDOPTIMIZATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_EXTENDED_OPTIMIZATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_INDEX_PCTFREE_DEFAULT", &obj.DB2INDEXPCTFREEDEFAULT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_INDEX_PCTFREE_DEFAULT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_INLIST_TO_NLJN", &obj.DB2INLISTTONLJN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_INLIST_TO_NLJN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_MINIMIZE_LISTPREFETCH", &obj.DB2MINIMIZELISTPREFETCH)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_MINIMIZE_LISTPREFETCH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OBJECT_TABLE_ENTRIES", &obj.DB2OBJECTTABLEENTRIES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OBJECT_TABLE_ENTRIES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPTPROFILE", &obj.DB2OPTPROFILE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPTPROFILE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPTSTATS_LOG", &obj.DB2OPTSTATSLOG)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPTSTATS_LOG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPT_MAX_TEMP_SIZE", &obj.DB2OPTMAXTEMPSIZE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPT_MAX_TEMP_SIZE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_PARALLEL_IO", &obj.DB2PARALLELIO)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_PARALLEL_IO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_REDUCED_OPTIMIZATION", &obj.DB2REDUCEDOPTIMIZATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_REDUCED_OPTIMIZATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SELECTIVITY", &obj.DB2SELECTIVITY)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SELECTIVITY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SKIPDELETED", &obj.DB2SKIPDELETED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SKIPDELETED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SKIPINSERTED", &obj.DB2SKIPINSERTED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SKIPINSERTED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SYNC_RELEASE_LOCK_ATTRIBUTES", &obj.DB2SYNCRELEASELOCKATTRIBUTES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SYNC_RELEASE_LOCK_ATTRIBUTES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_TRUNCATE_REUSESTORAGE", &obj.DB2TRUNCATEREUSESTORAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_TRUNCATE_REUSESTORAGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_USE_ALTERNATE_PAGE_CLEANING", &obj.DB2USEALTERNATEPAGECLEANING)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_USE_ALTERNATE_PAGE_CLEANING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_VIEW_REOPT_VALUES", &obj.DB2VIEWREOPTVALUES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_VIEW_REOPT_VALUES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_WLM_SETTINGS", &obj.DB2WLMSETTINGS)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_WLM_SETTINGS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_WORKLOAD", &obj.DB2WORKLOAD)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_WORKLOAD-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateUserAuthentication : CreateUserAuthentication struct
type CreateUserAuthentication struct {
	// Authentication method.
	Method *string `json:"method" validate:"required"`

	// Authentication policy ID.
	PolicyID *string `json:"policy_id" validate:"required"`
}

// NewCreateUserAuthentication : Instantiate CreateUserAuthentication (Generic Model Constructor)
func (*Db2saasV1) NewCreateUserAuthentication(method string, policyID string) (_model *CreateUserAuthentication, err error) {
	_model = &CreateUserAuthentication{
		Method: core.StringPtr(method),
		PolicyID: core.StringPtr(policyID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCreateUserAuthentication unmarshals an instance of CreateUserAuthentication from the specified map of raw messages.
func UnmarshalCreateUserAuthentication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateUserAuthentication)
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_id", &obj.PolicyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteDb2SaasUserOptions : The DeleteDb2SaasUser options.
type DeleteDb2SaasUserOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// id of the user.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDb2SaasUserOptions : Instantiate DeleteDb2SaasUserOptions
func (*Db2saasV1) NewDeleteDb2SaasUserOptions(xDeploymentID string, id string) *DeleteDb2SaasUserOptions {
	return &DeleteDb2SaasUserOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
		ID: core.StringPtr(id),
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *DeleteDb2SaasUserOptions) SetXDeploymentID(xDeploymentID string) *DeleteDb2SaasUserOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteDb2SaasUserOptions) SetID(id string) *DeleteDb2SaasUserOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDb2SaasUserOptions) SetHeaders(param map[string]string) *DeleteDb2SaasUserOptions {
	options.Headers = param
	return options
}

// GetDb2SaasAllowlistOptions : The GetDb2SaasAllowlist options.
type GetDb2SaasAllowlistOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasAllowlistOptions : Instantiate GetDb2SaasAllowlistOptions
func (*Db2saasV1) NewGetDb2SaasAllowlistOptions(xDeploymentID string) *GetDb2SaasAllowlistOptions {
	return &GetDb2SaasAllowlistOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *GetDb2SaasAllowlistOptions) SetXDeploymentID(xDeploymentID string) *GetDb2SaasAllowlistOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasAllowlistOptions) SetHeaders(param map[string]string) *GetDb2SaasAllowlistOptions {
	options.Headers = param
	return options
}

// GetDb2SaasAutoscaleOptions : The GetDb2SaasAutoscale options.
type GetDb2SaasAutoscaleOptions struct {
	// Encoded CRN deployment id.
	XDbProfile *string `json:"x-db-profile" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasAutoscaleOptions : Instantiate GetDb2SaasAutoscaleOptions
func (*Db2saasV1) NewGetDb2SaasAutoscaleOptions(xDbProfile string) *GetDb2SaasAutoscaleOptions {
	return &GetDb2SaasAutoscaleOptions{
		XDbProfile: core.StringPtr(xDbProfile),
	}
}

// SetXDbProfile : Allow user to set XDbProfile
func (_options *GetDb2SaasAutoscaleOptions) SetXDbProfile(xDbProfile string) *GetDb2SaasAutoscaleOptions {
	_options.XDbProfile = core.StringPtr(xDbProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasAutoscaleOptions) SetHeaders(param map[string]string) *GetDb2SaasAutoscaleOptions {
	options.Headers = param
	return options
}

// GetDb2SaasBackupOptions : The GetDb2SaasBackup options.
type GetDb2SaasBackupOptions struct {
	// Encoded CRN deployment id.
	XDbProfile *string `json:"x-db-profile" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasBackupOptions : Instantiate GetDb2SaasBackupOptions
func (*Db2saasV1) NewGetDb2SaasBackupOptions(xDbProfile string) *GetDb2SaasBackupOptions {
	return &GetDb2SaasBackupOptions{
		XDbProfile: core.StringPtr(xDbProfile),
	}
}

// SetXDbProfile : Allow user to set XDbProfile
func (_options *GetDb2SaasBackupOptions) SetXDbProfile(xDbProfile string) *GetDb2SaasBackupOptions {
	_options.XDbProfile = core.StringPtr(xDbProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasBackupOptions) SetHeaders(param map[string]string) *GetDb2SaasBackupOptions {
	options.Headers = param
	return options
}

// GetDb2SaasConnectionInfoOptions : The GetDb2SaasConnectionInfo options.
type GetDb2SaasConnectionInfoOptions struct {
	// Encoded CRN deployment id.
	DeploymentID *string `json:"deployment_id" validate:"required,ne="`

	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasConnectionInfoOptions : Instantiate GetDb2SaasConnectionInfoOptions
func (*Db2saasV1) NewGetDb2SaasConnectionInfoOptions(deploymentID string, xDeploymentID string) *GetDb2SaasConnectionInfoOptions {
	return &GetDb2SaasConnectionInfoOptions{
		DeploymentID: core.StringPtr(deploymentID),
		XDeploymentID: core.StringPtr(xDeploymentID),
	}
}

// SetDeploymentID : Allow user to set DeploymentID
func (_options *GetDb2SaasConnectionInfoOptions) SetDeploymentID(deploymentID string) *GetDb2SaasConnectionInfoOptions {
	_options.DeploymentID = core.StringPtr(deploymentID)
	return _options
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *GetDb2SaasConnectionInfoOptions) SetXDeploymentID(xDeploymentID string) *GetDb2SaasConnectionInfoOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasConnectionInfoOptions) SetHeaders(param map[string]string) *GetDb2SaasConnectionInfoOptions {
	options.Headers = param
	return options
}

// GetDb2SaasTuneableParamOptions : The GetDb2SaasTuneableParam options.
type GetDb2SaasTuneableParamOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasTuneableParamOptions : Instantiate GetDb2SaasTuneableParamOptions
func (*Db2saasV1) NewGetDb2SaasTuneableParamOptions() *GetDb2SaasTuneableParamOptions {
	return &GetDb2SaasTuneableParamOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasTuneableParamOptions) SetHeaders(param map[string]string) *GetDb2SaasTuneableParamOptions {
	options.Headers = param
	return options
}

// GetDb2SaasUserOptions : The GetDb2SaasUser options.
type GetDb2SaasUserOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDb2SaasUserOptions : Instantiate GetDb2SaasUserOptions
func (*Db2saasV1) NewGetDb2SaasUserOptions(xDeploymentID string) *GetDb2SaasUserOptions {
	return &GetDb2SaasUserOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *GetDb2SaasUserOptions) SetXDeploymentID(xDeploymentID string) *GetDb2SaasUserOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDb2SaasUserOptions) SetHeaders(param map[string]string) *GetDb2SaasUserOptions {
	options.Headers = param
	return options
}

// GetbyidDb2SaasUserOptions : The GetbyidDb2SaasUser options.
type GetbyidDb2SaasUserOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetbyidDb2SaasUserOptions : Instantiate GetbyidDb2SaasUserOptions
func (*Db2saasV1) NewGetbyidDb2SaasUserOptions(xDeploymentID string) *GetbyidDb2SaasUserOptions {
	return &GetbyidDb2SaasUserOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *GetbyidDb2SaasUserOptions) SetXDeploymentID(xDeploymentID string) *GetbyidDb2SaasUserOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetbyidDb2SaasUserOptions) SetHeaders(param map[string]string) *GetbyidDb2SaasUserOptions {
	options.Headers = param
	return options
}

// IpAddress : Details of an IP address.
type IpAddress struct {
	// The IP address, in IPv4/ipv6 format.
	Address *string `json:"address" validate:"required"`

	// Description of the IP address.
	Description *string `json:"description" validate:"required"`
}

// NewIpAddress : Instantiate IpAddress (Generic Model Constructor)
func (*Db2saasV1) NewIpAddress(address string, description string) (_model *IpAddress, err error) {
	_model = &IpAddress{
		Address: core.StringPtr(address),
		Description: core.StringPtr(description),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalIpAddress unmarshals an instance of IpAddress from the specified map of raw messages.
func UnmarshalIpAddress(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IpAddress)
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		err = core.SDKErrorf(err, "", "address-error", common.GetComponentInfo())
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

// PostDb2SaasAllowlistOptions : The PostDb2SaasAllowlist options.
type PostDb2SaasAllowlistOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// List of IP addresses.
	IpAddresses []IpAddress `json:"ip_addresses" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewPostDb2SaasAllowlistOptions : Instantiate PostDb2SaasAllowlistOptions
func (*Db2saasV1) NewPostDb2SaasAllowlistOptions(xDeploymentID string, ipAddresses []IpAddress) *PostDb2SaasAllowlistOptions {
	return &PostDb2SaasAllowlistOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
		IpAddresses: ipAddresses,
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *PostDb2SaasAllowlistOptions) SetXDeploymentID(xDeploymentID string) *PostDb2SaasAllowlistOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetIpAddresses : Allow user to set IpAddresses
func (_options *PostDb2SaasAllowlistOptions) SetIpAddresses(ipAddresses []IpAddress) *PostDb2SaasAllowlistOptions {
	_options.IpAddresses = ipAddresses
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDb2SaasAllowlistOptions) SetHeaders(param map[string]string) *PostDb2SaasAllowlistOptions {
	options.Headers = param
	return options
}

// PostDb2SaasBackupOptions : The PostDb2SaasBackup options.
type PostDb2SaasBackupOptions struct {
	// Encoded CRN deployment id.
	XDbProfile *string `json:"x-db-profile" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewPostDb2SaasBackupOptions : Instantiate PostDb2SaasBackupOptions
func (*Db2saasV1) NewPostDb2SaasBackupOptions(xDbProfile string) *PostDb2SaasBackupOptions {
	return &PostDb2SaasBackupOptions{
		XDbProfile: core.StringPtr(xDbProfile),
	}
}

// SetXDbProfile : Allow user to set XDbProfile
func (_options *PostDb2SaasBackupOptions) SetXDbProfile(xDbProfile string) *PostDb2SaasBackupOptions {
	_options.XDbProfile = core.StringPtr(xDbProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDb2SaasBackupOptions) SetHeaders(param map[string]string) *PostDb2SaasBackupOptions {
	options.Headers = param
	return options
}

// PostDb2SaasDbConfigurationOptions : The PostDb2SaasDbConfiguration options.
type PostDb2SaasDbConfigurationOptions struct {
	// Encoded CRN deployment id.
	XDbProfile *string `json:"x-db-profile" validate:"required"`

	// registry for db2 related configuration settings/configurations.
	Registry *CreateCustomSettingsRegistry `json:"registry,omitempty"`

	// Container for general database settings.
	Db *CreateCustomSettingsDb `json:"db,omitempty"`

	// Container for general database management settings.
	Dbm *CreateCustomSettingsDbm `json:"dbm,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewPostDb2SaasDbConfigurationOptions : Instantiate PostDb2SaasDbConfigurationOptions
func (*Db2saasV1) NewPostDb2SaasDbConfigurationOptions(xDbProfile string) *PostDb2SaasDbConfigurationOptions {
	return &PostDb2SaasDbConfigurationOptions{
		XDbProfile: core.StringPtr(xDbProfile),
	}
}

// SetXDbProfile : Allow user to set XDbProfile
func (_options *PostDb2SaasDbConfigurationOptions) SetXDbProfile(xDbProfile string) *PostDb2SaasDbConfigurationOptions {
	_options.XDbProfile = core.StringPtr(xDbProfile)
	return _options
}

// SetRegistry : Allow user to set Registry
func (_options *PostDb2SaasDbConfigurationOptions) SetRegistry(registry *CreateCustomSettingsRegistry) *PostDb2SaasDbConfigurationOptions {
	_options.Registry = registry
	return _options
}

// SetDb : Allow user to set Db
func (_options *PostDb2SaasDbConfigurationOptions) SetDb(db *CreateCustomSettingsDb) *PostDb2SaasDbConfigurationOptions {
	_options.Db = db
	return _options
}

// SetDbm : Allow user to set Dbm
func (_options *PostDb2SaasDbConfigurationOptions) SetDbm(dbm *CreateCustomSettingsDbm) *PostDb2SaasDbConfigurationOptions {
	_options.Dbm = dbm
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDb2SaasDbConfigurationOptions) SetHeaders(param map[string]string) *PostDb2SaasDbConfigurationOptions {
	options.Headers = param
	return options
}

// PostDb2SaasUserOptions : The PostDb2SaasUser options.
type PostDb2SaasUserOptions struct {
	// CRN deployment id.
	XDeploymentID *string `json:"x-deployment-id" validate:"required"`

	// The id of the User.
	ID *string `json:"id" validate:"required"`

	// Indicates if IAM is enabled.
	Iam *bool `json:"iam" validate:"required"`

	// IBM ID of the User.
	Ibmid *string `json:"ibmid" validate:"required"`

	// The name of the User.
	Name *string `json:"name" validate:"required"`

	// Password of the User.
	Password *string `json:"password" validate:"required"`

	// Role of the User.
	Role *string `json:"role" validate:"required"`

	// Email of the User.
	Email *string `json:"email" validate:"required"`

	// Indicates if the account is locked.
	Locked *string `json:"locked" validate:"required"`

	Authentication *CreateUserAuthentication `json:"authentication" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the PostDb2SaasUserOptions.Role property.
// Role of the User.
const (
	PostDb2SaasUserOptions_Role_Bluadmin = "bluadmin"
	PostDb2SaasUserOptions_Role_Bluuser = "bluuser"
)

// Constants associated with the PostDb2SaasUserOptions.Locked property.
// Indicates if the account is locked.
const (
	PostDb2SaasUserOptions_Locked_No = "no"
	PostDb2SaasUserOptions_Locked_Yes = "yes"
)

// NewPostDb2SaasUserOptions : Instantiate PostDb2SaasUserOptions
func (*Db2saasV1) NewPostDb2SaasUserOptions(xDeploymentID string, id string, iam bool, ibmid string, name string, password string, role string, email string, locked string, authentication *CreateUserAuthentication) *PostDb2SaasUserOptions {
	return &PostDb2SaasUserOptions{
		XDeploymentID: core.StringPtr(xDeploymentID),
		ID: core.StringPtr(id),
		Iam: core.BoolPtr(iam),
		Ibmid: core.StringPtr(ibmid),
		Name: core.StringPtr(name),
		Password: core.StringPtr(password),
		Role: core.StringPtr(role),
		Email: core.StringPtr(email),
		Locked: core.StringPtr(locked),
		Authentication: authentication,
	}
}

// SetXDeploymentID : Allow user to set XDeploymentID
func (_options *PostDb2SaasUserOptions) SetXDeploymentID(xDeploymentID string) *PostDb2SaasUserOptions {
	_options.XDeploymentID = core.StringPtr(xDeploymentID)
	return _options
}

// SetID : Allow user to set ID
func (_options *PostDb2SaasUserOptions) SetID(id string) *PostDb2SaasUserOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIam : Allow user to set Iam
func (_options *PostDb2SaasUserOptions) SetIam(iam bool) *PostDb2SaasUserOptions {
	_options.Iam = core.BoolPtr(iam)
	return _options
}

// SetIbmid : Allow user to set Ibmid
func (_options *PostDb2SaasUserOptions) SetIbmid(ibmid string) *PostDb2SaasUserOptions {
	_options.Ibmid = core.StringPtr(ibmid)
	return _options
}

// SetName : Allow user to set Name
func (_options *PostDb2SaasUserOptions) SetName(name string) *PostDb2SaasUserOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPassword : Allow user to set Password
func (_options *PostDb2SaasUserOptions) SetPassword(password string) *PostDb2SaasUserOptions {
	_options.Password = core.StringPtr(password)
	return _options
}

// SetRole : Allow user to set Role
func (_options *PostDb2SaasUserOptions) SetRole(role string) *PostDb2SaasUserOptions {
	_options.Role = core.StringPtr(role)
	return _options
}

// SetEmail : Allow user to set Email
func (_options *PostDb2SaasUserOptions) SetEmail(email string) *PostDb2SaasUserOptions {
	_options.Email = core.StringPtr(email)
	return _options
}

// SetLocked : Allow user to set Locked
func (_options *PostDb2SaasUserOptions) SetLocked(locked string) *PostDb2SaasUserOptions {
	_options.Locked = core.StringPtr(locked)
	return _options
}

// SetAuthentication : Allow user to set Authentication
func (_options *PostDb2SaasUserOptions) SetAuthentication(authentication *CreateUserAuthentication) *PostDb2SaasUserOptions {
	_options.Authentication = authentication
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostDb2SaasUserOptions) SetHeaders(param map[string]string) *PostDb2SaasUserOptions {
	options.Headers = param
	return options
}

// PutDb2SaasAutoscaleOptions : The PutDb2SaasAutoscale options.
type PutDb2SaasAutoscaleOptions struct {
	// Encoded CRN deployment id.
	XDbProfile *string `json:"x-db-profile" validate:"required"`

	// Indicates if automatic scaling is enabled or not.
	AutoScalingEnabled *string `json:"auto_scaling_enabled,omitempty"`

	// Specifies the resource utilization level that triggers an auto-scaling.
	AutoScalingThreshold *int64 `json:"auto_scaling_threshold,omitempty"`

	// Defines the time period over which auto-scaling adjustments are monitored and applied.
	AutoScalingOverTimePeriod *float64 `json:"auto_scaling_over_time_period,omitempty"`

	// Specifies the duration to pause auto-scaling actions after a scaling event has occurred.
	AutoScalingPauseLimit *int64 `json:"auto_scaling_pause_limit,omitempty"`

	// Indicates the maximum number of scaling actions that are allowed within a specified time period.
	AutoScalingAllowPlanLimit *string `json:"auto_scaling_allow_plan_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the PutDb2SaasAutoscaleOptions.AutoScalingEnabled property.
// Indicates if automatic scaling is enabled or not.
const (
	PutDb2SaasAutoscaleOptions_AutoScalingEnabled_False = "false"
	PutDb2SaasAutoscaleOptions_AutoScalingEnabled_True = "true"
)

// Constants associated with the PutDb2SaasAutoscaleOptions.AutoScalingAllowPlanLimit property.
// Indicates the maximum number of scaling actions that are allowed within a specified time period.
const (
	PutDb2SaasAutoscaleOptions_AutoScalingAllowPlanLimit_No = "NO"
	PutDb2SaasAutoscaleOptions_AutoScalingAllowPlanLimit_Yes = "YES"
)

// NewPutDb2SaasAutoscaleOptions : Instantiate PutDb2SaasAutoscaleOptions
func (*Db2saasV1) NewPutDb2SaasAutoscaleOptions(xDbProfile string) *PutDb2SaasAutoscaleOptions {
	return &PutDb2SaasAutoscaleOptions{
		XDbProfile: core.StringPtr(xDbProfile),
	}
}

// SetXDbProfile : Allow user to set XDbProfile
func (_options *PutDb2SaasAutoscaleOptions) SetXDbProfile(xDbProfile string) *PutDb2SaasAutoscaleOptions {
	_options.XDbProfile = core.StringPtr(xDbProfile)
	return _options
}

// SetAutoScalingEnabled : Allow user to set AutoScalingEnabled
func (_options *PutDb2SaasAutoscaleOptions) SetAutoScalingEnabled(autoScalingEnabled string) *PutDb2SaasAutoscaleOptions {
	_options.AutoScalingEnabled = core.StringPtr(autoScalingEnabled)
	return _options
}

// SetAutoScalingThreshold : Allow user to set AutoScalingThreshold
func (_options *PutDb2SaasAutoscaleOptions) SetAutoScalingThreshold(autoScalingThreshold int64) *PutDb2SaasAutoscaleOptions {
	_options.AutoScalingThreshold = core.Int64Ptr(autoScalingThreshold)
	return _options
}

// SetAutoScalingOverTimePeriod : Allow user to set AutoScalingOverTimePeriod
func (_options *PutDb2SaasAutoscaleOptions) SetAutoScalingOverTimePeriod(autoScalingOverTimePeriod float64) *PutDb2SaasAutoscaleOptions {
	_options.AutoScalingOverTimePeriod = core.Float64Ptr(autoScalingOverTimePeriod)
	return _options
}

// SetAutoScalingPauseLimit : Allow user to set AutoScalingPauseLimit
func (_options *PutDb2SaasAutoscaleOptions) SetAutoScalingPauseLimit(autoScalingPauseLimit int64) *PutDb2SaasAutoscaleOptions {
	_options.AutoScalingPauseLimit = core.Int64Ptr(autoScalingPauseLimit)
	return _options
}

// SetAutoScalingAllowPlanLimit : Allow user to set AutoScalingAllowPlanLimit
func (_options *PutDb2SaasAutoscaleOptions) SetAutoScalingAllowPlanLimit(autoScalingAllowPlanLimit string) *PutDb2SaasAutoscaleOptions {
	_options.AutoScalingAllowPlanLimit = core.StringPtr(autoScalingAllowPlanLimit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutDb2SaasAutoscaleOptions) SetHeaders(param map[string]string) *PutDb2SaasAutoscaleOptions {
	options.Headers = param
	return options
}

// SuccessAutoScaling : The details of the autoscale.
type SuccessAutoScaling struct {
	// Indicates the maximum number of scaling actions that are allowed within a specified time period.
	AutoScalingAllowPlanLimit *bool `json:"auto_scaling_allow_plan_limit" validate:"required"`

	// Indicates if automatic scaling is enabled or not.
	AutoScalingEnabled *bool `json:"auto_scaling_enabled" validate:"required"`

	// The maximum limit for automatically increasing storage capacity to handle growing data needs.
	AutoScalingMaxStorage *int64 `json:"auto_scaling_max_storage" validate:"required"`

	// Defines the time period over which auto-scaling adjustments are monitored and applied.
	AutoScalingOverTimePeriod *int64 `json:"auto_scaling_over_time_period" validate:"required"`

	// Specifies the duration to pause auto-scaling actions after a scaling event has occurred.
	AutoScalingPauseLimit *int64 `json:"auto_scaling_pause_limit" validate:"required"`

	// Specifies the resource utilization level that triggers an auto-scaling.
	AutoScalingThreshold *int64 `json:"auto_scaling_threshold" validate:"required"`

	// Specifies the unit of measurement for storage capacity.
	StorageUnit *string `json:"storage_unit" validate:"required"`

	// Represents the percentage of total storage capacity currently in use.
	StorageUtilizationPercentage *int64 `json:"storage_utilization_percentage" validate:"required"`

	// Indicates whether a system or service can automatically adjust resources based on demand.
	SupportAutoScaling *bool `json:"support_auto_scaling" validate:"required"`
}

// UnmarshalSuccessAutoScaling unmarshals an instance of SuccessAutoScaling from the specified map of raw messages.
func UnmarshalSuccessAutoScaling(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessAutoScaling)
	err = core.UnmarshalPrimitive(m, "auto_scaling_allow_plan_limit", &obj.AutoScalingAllowPlanLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_allow_plan_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_scaling_enabled", &obj.AutoScalingEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_scaling_max_storage", &obj.AutoScalingMaxStorage)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_max_storage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_scaling_over_time_period", &obj.AutoScalingOverTimePeriod)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_over_time_period-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_scaling_pause_limit", &obj.AutoScalingPauseLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_pause_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_scaling_threshold", &obj.AutoScalingThreshold)
	if err != nil {
		err = core.SDKErrorf(err, "", "auto_scaling_threshold-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "storage_unit", &obj.StorageUnit)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_unit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "storage_utilization_percentage", &obj.StorageUtilizationPercentage)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_utilization_percentage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "support_auto_scaling", &obj.SupportAutoScaling)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_auto_scaling-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessConnectionInfo : Responds with JSON of the connection information for the Db2 SaaS Instance.
type SuccessConnectionInfo struct {
	Public *SuccessConnectionInfoPublic `json:"public,omitempty"`

	Private *SuccessConnectionInfoPrivate `json:"private,omitempty"`
}

// UnmarshalSuccessConnectionInfo unmarshals an instance of SuccessConnectionInfo from the specified map of raw messages.
func UnmarshalSuccessConnectionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessConnectionInfo)
	err = core.UnmarshalModel(m, "public", &obj.Public, UnmarshalSuccessConnectionInfoPublic)
	if err != nil {
		err = core.SDKErrorf(err, "", "public-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "private", &obj.Private, UnmarshalSuccessConnectionInfoPrivate)
	if err != nil {
		err = core.SDKErrorf(err, "", "private-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessConnectionInfoPrivate : SuccessConnectionInfoPrivate struct
type SuccessConnectionInfoPrivate struct {
	Hostname *string `json:"hostname,omitempty"`

	DatabaseName *string `json:"databaseName,omitempty"`

	SslPort *string `json:"sslPort,omitempty"`

	Ssl *bool `json:"ssl,omitempty"`

	DatabaseVersion *string `json:"databaseVersion,omitempty"`

	PrivateServiceName *string `json:"private_serviceName,omitempty"`

	CloudServiceOffering *string `json:"cloud_service_offering,omitempty"`

	VpeServiceCrn *string `json:"vpe_service_crn,omitempty"`

	DbVpcEndpointService *string `json:"db_vpc_endpoint_service,omitempty"`
}

// UnmarshalSuccessConnectionInfoPrivate unmarshals an instance of SuccessConnectionInfoPrivate from the specified map of raw messages.
func UnmarshalSuccessConnectionInfoPrivate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessConnectionInfoPrivate)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		err = core.SDKErrorf(err, "", "hostname-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "databaseName", &obj.DatabaseName)
	if err != nil {
		err = core.SDKErrorf(err, "", "databaseName-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sslPort", &obj.SslPort)
	if err != nil {
		err = core.SDKErrorf(err, "", "sslPort-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ssl", &obj.Ssl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ssl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "databaseVersion", &obj.DatabaseVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "databaseVersion-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_serviceName", &obj.PrivateServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_serviceName-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cloud_service_offering", &obj.CloudServiceOffering)
	if err != nil {
		err = core.SDKErrorf(err, "", "cloud_service_offering-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "vpe_service_crn", &obj.VpeServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpe_service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "db_vpc_endpoint_service", &obj.DbVpcEndpointService)
	if err != nil {
		err = core.SDKErrorf(err, "", "db_vpc_endpoint_service-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessConnectionInfoPublic : SuccessConnectionInfoPublic struct
type SuccessConnectionInfoPublic struct {
	Hostname *string `json:"hostname,omitempty"`

	DatabaseName *string `json:"databaseName,omitempty"`

	SslPort *string `json:"sslPort,omitempty"`

	Ssl *bool `json:"ssl,omitempty"`

	DatabaseVersion *string `json:"databaseVersion,omitempty"`
}

// UnmarshalSuccessConnectionInfoPublic unmarshals an instance of SuccessConnectionInfoPublic from the specified map of raw messages.
func UnmarshalSuccessConnectionInfoPublic(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessConnectionInfoPublic)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		err = core.SDKErrorf(err, "", "hostname-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "databaseName", &obj.DatabaseName)
	if err != nil {
		err = core.SDKErrorf(err, "", "databaseName-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sslPort", &obj.SslPort)
	if err != nil {
		err = core.SDKErrorf(err, "", "sslPort-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ssl", &obj.Ssl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ssl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "databaseVersion", &obj.DatabaseVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "databaseVersion-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessCreateBackup : Success response of post backup.
type SuccessCreateBackup struct {
	Task *SuccessCreateBackupTask `json:"task" validate:"required"`
}

// UnmarshalSuccessCreateBackup unmarshals an instance of SuccessCreateBackup from the specified map of raw messages.
func UnmarshalSuccessCreateBackup(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessCreateBackup)
	err = core.UnmarshalModel(m, "task", &obj.Task, UnmarshalSuccessCreateBackupTask)
	if err != nil {
		err = core.SDKErrorf(err, "", "task-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessCreateBackupTask : SuccessCreateBackupTask struct
type SuccessCreateBackupTask struct {
	// CRN of the instance.
	ID *string `json:"id,omitempty"`
}

// UnmarshalSuccessCreateBackupTask unmarshals an instance of SuccessCreateBackupTask from the specified map of raw messages.
func UnmarshalSuccessCreateBackupTask(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessCreateBackupTask)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetAllowlistIPs : Success response of get allowlist IPs.
type SuccessGetAllowlistIPs struct {
	// List of IP addresses.
	IpAddresses []IpAddress `json:"ip_addresses" validate:"required"`
}

// UnmarshalSuccessGetAllowlistIPs unmarshals an instance of SuccessGetAllowlistIPs from the specified map of raw messages.
func UnmarshalSuccessGetAllowlistIPs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetAllowlistIPs)
	err = core.UnmarshalModel(m, "ip_addresses", &obj.IpAddresses, UnmarshalIpAddress)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip_addresses-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetBackups : The details of the backups.
type SuccessGetBackups struct {
	Backups []Backup `json:"backups" validate:"required"`
}

// UnmarshalSuccessGetBackups unmarshals an instance of SuccessGetBackups from the specified map of raw messages.
func UnmarshalSuccessGetBackups(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetBackups)
	err = core.UnmarshalModel(m, "backups", &obj.Backups, UnmarshalBackup)
	if err != nil {
		err = core.SDKErrorf(err, "", "backups-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetUserByID : The details of the users.
type SuccessGetUserByID struct {
	// User's DV role.
	DvRole *string `json:"dvRole" validate:"required"`

	// Metadata associated with the user.
	Metadata map[string]interface{} `json:"metadata" validate:"required"`

	// Formatted IBM ID.
	FormatedIbmid *string `json:"formatedIbmid" validate:"required"`

	// Role assigned to the user.
	Role *string `json:"role" validate:"required"`

	// IAM ID for the user.
	Iamid *string `json:"iamid" validate:"required"`

	// List of allowed actions of the user.
	PermittedActions []string `json:"permittedActions" validate:"required"`

	// Indicates if the user account has no issues.
	AllClean *bool `json:"allClean" validate:"required"`

	// User's password.
	Password *string `json:"password" validate:"required"`

	// Indicates if IAM is enabled or not.
	Iam *bool `json:"iam" validate:"required"`

	// The display name of the user.
	Name *string `json:"name" validate:"required"`

	// IBM ID of the user.
	Ibmid *string `json:"ibmid" validate:"required"`

	// Unique identifier for the user.
	ID *string `json:"id" validate:"required"`

	// Account lock status for the user.
	Locked *string `json:"locked" validate:"required"`

	// Initial error message.
	InitErrorMsg *string `json:"initErrorMsg" validate:"required"`

	// Email address of the user.
	Email *string `json:"email" validate:"required"`

	// Authentication details for the user.
	Authentication *SuccessGetUserByIDAuthentication `json:"authentication" validate:"required"`
}

// Constants associated with the SuccessGetUserByID.Role property.
// Role assigned to the user.
const (
	SuccessGetUserByID_Role_Bluadmin = "bluadmin"
	SuccessGetUserByID_Role_Bluuser = "bluuser"
)

// Constants associated with the SuccessGetUserByID.Locked property.
// Account lock status for the user.
const (
	SuccessGetUserByID_Locked_No = "no"
	SuccessGetUserByID_Locked_Yes = "yes"
)

// UnmarshalSuccessGetUserByID unmarshals an instance of SuccessGetUserByID from the specified map of raw messages.
func UnmarshalSuccessGetUserByID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetUserByID)
	err = core.UnmarshalPrimitive(m, "dvRole", &obj.DvRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "dvRole-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "formatedIbmid", &obj.FormatedIbmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "formatedIbmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "role", &obj.Role)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iamid", &obj.Iamid)
	if err != nil {
		err = core.SDKErrorf(err, "", "iamid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "permittedActions", &obj.PermittedActions)
	if err != nil {
		err = core.SDKErrorf(err, "", "permittedActions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allClean", &obj.AllClean)
	if err != nil {
		err = core.SDKErrorf(err, "", "allClean-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam", &obj.Iam)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmid", &obj.Ibmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initErrorMsg", &obj.InitErrorMsg)
	if err != nil {
		err = core.SDKErrorf(err, "", "initErrorMsg-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authentication", &obj.Authentication, UnmarshalSuccessGetUserByIDAuthentication)
	if err != nil {
		err = core.SDKErrorf(err, "", "authentication-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetUserByIDAuthentication : Authentication details for the user.
type SuccessGetUserByIDAuthentication struct {
	// Authentication method.
	Method *string `json:"method" validate:"required"`

	// Policy ID of authentication.
	PolicyID *string `json:"policy_id" validate:"required"`
}

// UnmarshalSuccessGetUserByIDAuthentication unmarshals an instance of SuccessGetUserByIDAuthentication from the specified map of raw messages.
func UnmarshalSuccessGetUserByIDAuthentication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetUserByIDAuthentication)
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_id", &obj.PolicyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetUserInfo : Success response of get user.
type SuccessGetUserInfo struct {
	// The total number of resources.
	Count *int64 `json:"count" validate:"required"`

	// A list of user resource.
	Resources []SuccessGetUserInfoResourcesItem `json:"resources" validate:"required"`
}

// UnmarshalSuccessGetUserInfo unmarshals an instance of SuccessGetUserInfo from the specified map of raw messages.
func UnmarshalSuccessGetUserInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetUserInfo)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalSuccessGetUserInfoResourcesItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetUserInfoResourcesItem : SuccessGetUserInfoResourcesItem struct
type SuccessGetUserInfoResourcesItem struct {
	// User's DV role.
	DvRole *string `json:"dvRole,omitempty"`

	// Metadata associated with the user.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Formatted IBM ID.
	FormatedIbmid *string `json:"formatedIbmid,omitempty"`

	// Role assigned to the user.
	Role *string `json:"role,omitempty"`

	// IAM ID for the user.
	Iamid *string `json:"iamid,omitempty"`

	// List of allowed actions of the user.
	PermittedActions []string `json:"permittedActions,omitempty"`

	// Indicates if the user account has no issues.
	AllClean *bool `json:"allClean,omitempty"`

	// User's password.
	Password *string `json:"password,omitempty"`

	// Indicates if IAM is enabled or not.
	Iam *bool `json:"iam,omitempty"`

	// The display name of the user.
	Name *string `json:"name,omitempty"`

	// IBM ID of the user.
	Ibmid *string `json:"ibmid,omitempty"`

	// Unique identifier for the user.
	ID *string `json:"id,omitempty"`

	// Account lock status for the user.
	Locked *string `json:"locked,omitempty"`

	// Initial error message.
	InitErrorMsg *string `json:"initErrorMsg,omitempty"`

	// Email address of the user.
	Email *string `json:"email,omitempty"`

	// Authentication details for the user.
	Authentication *SuccessGetUserInfoResourcesItemAuthentication `json:"authentication,omitempty"`
}

// Constants associated with the SuccessGetUserInfoResourcesItem.Role property.
// Role assigned to the user.
const (
	SuccessGetUserInfoResourcesItem_Role_Bluadmin = "bluadmin"
	SuccessGetUserInfoResourcesItem_Role_Bluuser = "bluuser"
)

// Constants associated with the SuccessGetUserInfoResourcesItem.Locked property.
// Account lock status for the user.
const (
	SuccessGetUserInfoResourcesItem_Locked_No = "no"
	SuccessGetUserInfoResourcesItem_Locked_Yes = "yes"
)

// UnmarshalSuccessGetUserInfoResourcesItem unmarshals an instance of SuccessGetUserInfoResourcesItem from the specified map of raw messages.
func UnmarshalSuccessGetUserInfoResourcesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetUserInfoResourcesItem)
	err = core.UnmarshalPrimitive(m, "dvRole", &obj.DvRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "dvRole-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "formatedIbmid", &obj.FormatedIbmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "formatedIbmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "role", &obj.Role)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iamid", &obj.Iamid)
	if err != nil {
		err = core.SDKErrorf(err, "", "iamid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "permittedActions", &obj.PermittedActions)
	if err != nil {
		err = core.SDKErrorf(err, "", "permittedActions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allClean", &obj.AllClean)
	if err != nil {
		err = core.SDKErrorf(err, "", "allClean-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam", &obj.Iam)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmid", &obj.Ibmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initErrorMsg", &obj.InitErrorMsg)
	if err != nil {
		err = core.SDKErrorf(err, "", "initErrorMsg-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authentication", &obj.Authentication, UnmarshalSuccessGetUserInfoResourcesItemAuthentication)
	if err != nil {
		err = core.SDKErrorf(err, "", "authentication-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessGetUserInfoResourcesItemAuthentication : Authentication details for the user.
type SuccessGetUserInfoResourcesItemAuthentication struct {
	// Authentication method.
	Method *string `json:"method" validate:"required"`

	// Policy ID of authentication.
	PolicyID *string `json:"policy_id" validate:"required"`
}

// UnmarshalSuccessGetUserInfoResourcesItemAuthentication unmarshals an instance of SuccessGetUserInfoResourcesItemAuthentication from the specified map of raw messages.
func UnmarshalSuccessGetUserInfoResourcesItemAuthentication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessGetUserInfoResourcesItemAuthentication)
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_id", &obj.PolicyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessPostAllowedlistIPs : Success response of post allowlist IPs.
type SuccessPostAllowedlistIPs struct {
	// status of the post allowlist IPs request.
	Status *string `json:"status" validate:"required"`
}

// UnmarshalSuccessPostAllowedlistIPs unmarshals an instance of SuccessPostAllowedlistIPs from the specified map of raw messages.
func UnmarshalSuccessPostAllowedlistIPs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessPostAllowedlistIPs)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessPostCustomSettings : The details of created custom settings of db2.
type SuccessPostCustomSettings struct {
	// Describes the operation done.
	Description *string `json:"description" validate:"required"`

	// CRN of the db2 instance.
	ID *string `json:"id" validate:"required"`

	// Defines the status of the instance.
	Status *string `json:"status" validate:"required"`
}

// UnmarshalSuccessPostCustomSettings unmarshals an instance of SuccessPostCustomSettings from the specified map of raw messages.
func UnmarshalSuccessPostCustomSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessPostCustomSettings)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
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

// SuccessTuneableParams : Response of tuneable params of the Db2 instance.
type SuccessTuneableParams struct {
	TuneableParam *SuccessTuneableParamsTuneableParam `json:"tuneable_param,omitempty"`
}

// UnmarshalSuccessTuneableParams unmarshals an instance of SuccessTuneableParams from the specified map of raw messages.
func UnmarshalSuccessTuneableParams(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessTuneableParams)
	err = core.UnmarshalModel(m, "tuneable_param", &obj.TuneableParam, UnmarshalSuccessTuneableParamsTuneableParam)
	if err != nil {
		err = core.SDKErrorf(err, "", "tuneable_param-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessTuneableParamsTuneableParam : SuccessTuneableParamsTuneableParam struct
type SuccessTuneableParamsTuneableParam struct {
	// Tunable parameters related to the Db2 database instance.
	Db *SuccessTuneableParamsTuneableParamDb `json:"db,omitempty"`

	// Tunable parameters related to the Db2 instance manager (dbm).
	Dbm *SuccessTuneableParamsTuneableParamDbm `json:"dbm,omitempty"`

	// Tunable parameters related to the Db2 registry.
	Registry *SuccessTuneableParamsTuneableParamRegistry `json:"registry,omitempty"`
}

// UnmarshalSuccessTuneableParamsTuneableParam unmarshals an instance of SuccessTuneableParamsTuneableParam from the specified map of raw messages.
func UnmarshalSuccessTuneableParamsTuneableParam(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessTuneableParamsTuneableParam)
	err = core.UnmarshalModel(m, "db", &obj.Db, UnmarshalSuccessTuneableParamsTuneableParamDb)
	if err != nil {
		err = core.SDKErrorf(err, "", "db-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "dbm", &obj.Dbm, UnmarshalSuccessTuneableParamsTuneableParamDbm)
	if err != nil {
		err = core.SDKErrorf(err, "", "dbm-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "registry", &obj.Registry, UnmarshalSuccessTuneableParamsTuneableParamRegistry)
	if err != nil {
		err = core.SDKErrorf(err, "", "registry-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessTuneableParamsTuneableParamDb : Tunable parameters related to the Db2 database instance.
type SuccessTuneableParamsTuneableParamDb struct {
	ACTSORTMEMLIMIT *string `json:"ACT_SORTMEM_LIMIT,omitempty"`

	ALTCOLLATE *string `json:"ALT_COLLATE,omitempty"`

	APPGROUPMEMSZ *string `json:"APPGROUP_MEM_SZ,omitempty"`

	APPLHEAPSZ *string `json:"APPLHEAPSZ,omitempty"`

	APPLMEMORY *string `json:"APPL_MEMORY,omitempty"`

	APPCTLHEAPSZ *string `json:"APP_CTL_HEAP_SZ,omitempty"`

	ARCHRETRYDELAY *string `json:"ARCHRETRYDELAY,omitempty"`

	AUTHNCACHEDURATION *string `json:"AUTHN_CACHE_DURATION,omitempty"`

	AUTORESTART *string `json:"AUTORESTART,omitempty"`

	AUTOCGSTATS *string `json:"AUTO_CG_STATS,omitempty"`

	AUTOMAINT *string `json:"AUTO_MAINT,omitempty"`

	AUTOREORG *string `json:"AUTO_REORG,omitempty"`

	AUTOREVAL *string `json:"AUTO_REVAL,omitempty"`

	AUTORUNSTATS *string `json:"AUTO_RUNSTATS,omitempty"`

	AUTOSAMPLING *string `json:"AUTO_SAMPLING,omitempty"`

	AUTOSTATSVIEWS *string `json:"AUTO_STATS_VIEWS,omitempty"`

	AUTOSTMTSTATS *string `json:"AUTO_STMT_STATS,omitempty"`

	AUTOTBLMAINT *string `json:"AUTO_TBL_MAINT,omitempty"`

	AVGAPPLS *string `json:"AVG_APPLS,omitempty"`

	CATALOGCACHESZ *string `json:"CATALOGCACHE_SZ,omitempty"`

	CHNGPGSTHRESH *string `json:"CHNGPGS_THRESH,omitempty"`

	CURCOMMIT *string `json:"CUR_COMMIT,omitempty"`

	DATABASEMEMORY *string `json:"DATABASE_MEMORY,omitempty"`

	DBHEAP *string `json:"DBHEAP,omitempty"`

	DBCOLLNAME *string `json:"DB_COLLNAME,omitempty"`

	DBMEMTHRESH *string `json:"DB_MEM_THRESH,omitempty"`

	DDLCOMPRESSIONDEF *string `json:"DDL_COMPRESSION_DEF,omitempty"`

	DDLCONSTRAINTDEF *string `json:"DDL_CONSTRAINT_DEF,omitempty"`

	DECFLTROUNDING *string `json:"DECFLT_ROUNDING,omitempty"`

	DECARITHMETIC *string `json:"DEC_ARITHMETIC,omitempty"`

	DECTOCHARFMT *string `json:"DEC_TO_CHAR_FMT,omitempty"`

	DFTDEGREE *string `json:"DFT_DEGREE,omitempty"`

	DFTEXTENTSZ *string `json:"DFT_EXTENT_SZ,omitempty"`

	DFTLOADRECSES *string `json:"DFT_LOADREC_SES,omitempty"`

	DFTMTTBTYPES *string `json:"DFT_MTTB_TYPES,omitempty"`

	DFTPREFETCHSZ *string `json:"DFT_PREFETCH_SZ,omitempty"`

	DFTQUERYOPT *string `json:"DFT_QUERYOPT,omitempty"`

	DFTREFRESHAGE *string `json:"DFT_REFRESH_AGE,omitempty"`

	DFTSCHEMASDCC *string `json:"DFT_SCHEMAS_DCC,omitempty"`

	DFTSQLMATHWARN *string `json:"DFT_SQLMATHWARN,omitempty"`

	DFTTABLEORG *string `json:"DFT_TABLE_ORG,omitempty"`

	DLCHKTIME *string `json:"DLCHKTIME,omitempty"`

	ENABLEXMLCHAR *string `json:"ENABLE_XMLCHAR,omitempty"`

	EXTENDEDROWSZ *string `json:"EXTENDED_ROW_SZ,omitempty"`

	GROUPHEAPRATIO *string `json:"GROUPHEAP_RATIO,omitempty"`

	INDEXREC *string `json:"INDEXREC,omitempty"`

	LARGEAGGREGATION *string `json:"LARGE_AGGREGATION,omitempty"`

	LOCKLIST *string `json:"LOCKLIST,omitempty"`

	LOCKTIMEOUT *string `json:"LOCKTIMEOUT,omitempty"`

	LOGINDEXBUILD *string `json:"LOGINDEXBUILD,omitempty"`

	LOGAPPLINFO *string `json:"LOG_APPL_INFO,omitempty"`

	LOGDDLSTMTS *string `json:"LOG_DDL_STMTS,omitempty"`

	LOGDISKCAP *string `json:"LOG_DISK_CAP,omitempty"`

	MAXAPPLS *string `json:"MAXAPPLS,omitempty"`

	MAXFILOP *string `json:"MAXFILOP,omitempty"`

	MAXLOCKS *string `json:"MAXLOCKS,omitempty"`

	MINDECDIV3 *string `json:"MIN_DEC_DIV_3,omitempty"`

	MONACTMETRICS *string `json:"MON_ACT_METRICS,omitempty"`

	MONDEADLOCK *string `json:"MON_DEADLOCK,omitempty"`

	MONLCKMSGLVL *string `json:"MON_LCK_MSG_LVL,omitempty"`

	MONLOCKTIMEOUT *string `json:"MON_LOCKTIMEOUT,omitempty"`

	MONLOCKWAIT *string `json:"MON_LOCKWAIT,omitempty"`

	MONLWTHRESH *string `json:"MON_LW_THRESH,omitempty"`

	MONOBJMETRICS *string `json:"MON_OBJ_METRICS,omitempty"`

	MONPKGLISTSZ *string `json:"MON_PKGLIST_SZ,omitempty"`

	MONREQMETRICS *string `json:"MON_REQ_METRICS,omitempty"`

	MONRTNDATA *string `json:"MON_RTN_DATA,omitempty"`

	MONRTNEXECLIST *string `json:"MON_RTN_EXECLIST,omitempty"`

	MONUOWDATA *string `json:"MON_UOW_DATA,omitempty"`

	MONUOWEXECLIST *string `json:"MON_UOW_EXECLIST,omitempty"`

	MONUOWPKGLIST *string `json:"MON_UOW_PKGLIST,omitempty"`

	NCHARMAPPING *string `json:"NCHAR_MAPPING,omitempty"`

	NUMFREQVALUES *string `json:"NUM_FREQVALUES,omitempty"`

	NUMIOCLEANERS *string `json:"NUM_IOCLEANERS,omitempty"`

	NUMIOSERVERS *string `json:"NUM_IOSERVERS,omitempty"`

	NUMLOGSPAN *string `json:"NUM_LOG_SPAN,omitempty"`

	NUMQUANTILES *string `json:"NUM_QUANTILES,omitempty"`

	OPTBUFFPAGE *string `json:"OPT_BUFFPAGE,omitempty"`

	OPTDIRECTWRKLD *string `json:"OPT_DIRECT_WRKLD,omitempty"`

	OPTLOCKLIST *string `json:"OPT_LOCKLIST,omitempty"`

	OPTMAXLOCKS *string `json:"OPT_MAXLOCKS,omitempty"`

	OPTSORTHEAP *string `json:"OPT_SORTHEAP,omitempty"`

	PAGEAGETRGTGCR *string `json:"PAGE_AGE_TRGT_GCR,omitempty"`

	PAGEAGETRGTMCR *string `json:"PAGE_AGE_TRGT_MCR,omitempty"`

	PCKCACHESZ *string `json:"PCKCACHESZ,omitempty"`

	PLSTACKTRACE *string `json:"PL_STACK_TRACE,omitempty"`

	SELFTUNINGMEM *string `json:"SELF_TUNING_MEM,omitempty"`

	SEQDETECT *string `json:"SEQDETECT,omitempty"`

	SHEAPTHRESSHR *string `json:"SHEAPTHRES_SHR,omitempty"`

	SOFTMAX *string `json:"SOFTMAX,omitempty"`

	SORTHEAP *string `json:"SORTHEAP,omitempty"`

	SQLCCFLAGS *string `json:"SQL_CCFLAGS,omitempty"`

	STATHEAPSZ *string `json:"STAT_HEAP_SZ,omitempty"`

	STMTHEAP *string `json:"STMTHEAP,omitempty"`

	STMTCONC *string `json:"STMT_CONC,omitempty"`

	STRINGUNITS *string `json:"STRING_UNITS,omitempty"`

	SYSTIMEPERIODADJ *string `json:"SYSTIME_PERIOD_ADJ,omitempty"`

	TRACKMOD *string `json:"TRACKMOD,omitempty"`

	UTILHEAPSZ *string `json:"UTIL_HEAP_SZ,omitempty"`

	WLMADMISSIONCTRL *string `json:"WLM_ADMISSION_CTRL,omitempty"`

	WLMAGENTLOADTRGT *string `json:"WLM_AGENT_LOAD_TRGT,omitempty"`

	WLMCPULIMIT *string `json:"WLM_CPU_LIMIT,omitempty"`

	WLMCPUSHARES *string `json:"WLM_CPU_SHARES,omitempty"`

	WLMCPUSHAREMODE *string `json:"WLM_CPU_SHARE_MODE,omitempty"`
}

// UnmarshalSuccessTuneableParamsTuneableParamDb unmarshals an instance of SuccessTuneableParamsTuneableParamDb from the specified map of raw messages.
func UnmarshalSuccessTuneableParamsTuneableParamDb(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessTuneableParamsTuneableParamDb)
	err = core.UnmarshalPrimitive(m, "ACT_SORTMEM_LIMIT", &obj.ACTSORTMEMLIMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "ACT_SORTMEM_LIMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ALT_COLLATE", &obj.ALTCOLLATE)
	if err != nil {
		err = core.SDKErrorf(err, "", "ALT_COLLATE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPGROUP_MEM_SZ", &obj.APPGROUPMEMSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPGROUP_MEM_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPLHEAPSZ", &obj.APPLHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPLHEAPSZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APPL_MEMORY", &obj.APPLMEMORY)
	if err != nil {
		err = core.SDKErrorf(err, "", "APPL_MEMORY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "APP_CTL_HEAP_SZ", &obj.APPCTLHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "APP_CTL_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ARCHRETRYDELAY", &obj.ARCHRETRYDELAY)
	if err != nil {
		err = core.SDKErrorf(err, "", "ARCHRETRYDELAY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTHN_CACHE_DURATION", &obj.AUTHNCACHEDURATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTHN_CACHE_DURATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTORESTART", &obj.AUTORESTART)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTORESTART-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_CG_STATS", &obj.AUTOCGSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_CG_STATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_MAINT", &obj.AUTOMAINT)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_MAINT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_REORG", &obj.AUTOREORG)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_REORG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_REVAL", &obj.AUTOREVAL)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_REVAL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_RUNSTATS", &obj.AUTORUNSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_RUNSTATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_SAMPLING", &obj.AUTOSAMPLING)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_SAMPLING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_STATS_VIEWS", &obj.AUTOSTATSVIEWS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_STATS_VIEWS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_STMT_STATS", &obj.AUTOSTMTSTATS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_STMT_STATS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AUTO_TBL_MAINT", &obj.AUTOTBLMAINT)
	if err != nil {
		err = core.SDKErrorf(err, "", "AUTO_TBL_MAINT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "AVG_APPLS", &obj.AVGAPPLS)
	if err != nil {
		err = core.SDKErrorf(err, "", "AVG_APPLS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CATALOGCACHE_SZ", &obj.CATALOGCACHESZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "CATALOGCACHE_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CHNGPGS_THRESH", &obj.CHNGPGSTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "CHNGPGS_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CUR_COMMIT", &obj.CURCOMMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "CUR_COMMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DATABASE_MEMORY", &obj.DATABASEMEMORY)
	if err != nil {
		err = core.SDKErrorf(err, "", "DATABASE_MEMORY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DBHEAP", &obj.DBHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "DBHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB_COLLNAME", &obj.DBCOLLNAME)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB_COLLNAME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB_MEM_THRESH", &obj.DBMEMTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB_MEM_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DDL_COMPRESSION_DEF", &obj.DDLCOMPRESSIONDEF)
	if err != nil {
		err = core.SDKErrorf(err, "", "DDL_COMPRESSION_DEF-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DDL_CONSTRAINT_DEF", &obj.DDLCONSTRAINTDEF)
	if err != nil {
		err = core.SDKErrorf(err, "", "DDL_CONSTRAINT_DEF-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DECFLT_ROUNDING", &obj.DECFLTROUNDING)
	if err != nil {
		err = core.SDKErrorf(err, "", "DECFLT_ROUNDING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DEC_ARITHMETIC", &obj.DECARITHMETIC)
	if err != nil {
		err = core.SDKErrorf(err, "", "DEC_ARITHMETIC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DEC_TO_CHAR_FMT", &obj.DECTOCHARFMT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DEC_TO_CHAR_FMT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_DEGREE", &obj.DFTDEGREE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_DEGREE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_EXTENT_SZ", &obj.DFTEXTENTSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_EXTENT_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_LOADREC_SES", &obj.DFTLOADRECSES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_LOADREC_SES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MTTB_TYPES", &obj.DFTMTTBTYPES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MTTB_TYPES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_PREFETCH_SZ", &obj.DFTPREFETCHSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_PREFETCH_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_QUERYOPT", &obj.DFTQUERYOPT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_QUERYOPT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_REFRESH_AGE", &obj.DFTREFRESHAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_REFRESH_AGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_SCHEMAS_DCC", &obj.DFTSCHEMASDCC)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_SCHEMAS_DCC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_SQLMATHWARN", &obj.DFTSQLMATHWARN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_SQLMATHWARN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_TABLE_ORG", &obj.DFTTABLEORG)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_TABLE_ORG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DLCHKTIME", &obj.DLCHKTIME)
	if err != nil {
		err = core.SDKErrorf(err, "", "DLCHKTIME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ENABLE_XMLCHAR", &obj.ENABLEXMLCHAR)
	if err != nil {
		err = core.SDKErrorf(err, "", "ENABLE_XMLCHAR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "EXTENDED_ROW_SZ", &obj.EXTENDEDROWSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "EXTENDED_ROW_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "GROUPHEAP_RATIO", &obj.GROUPHEAPRATIO)
	if err != nil {
		err = core.SDKErrorf(err, "", "GROUPHEAP_RATIO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INDEXREC", &obj.INDEXREC)
	if err != nil {
		err = core.SDKErrorf(err, "", "INDEXREC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LARGE_AGGREGATION", &obj.LARGEAGGREGATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "LARGE_AGGREGATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOCKLIST", &obj.LOCKLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOCKLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOCKTIMEOUT", &obj.LOCKTIMEOUT)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOCKTIMEOUT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOGINDEXBUILD", &obj.LOGINDEXBUILD)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOGINDEXBUILD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_APPL_INFO", &obj.LOGAPPLINFO)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_APPL_INFO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_DDL_STMTS", &obj.LOGDDLSTMTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_DDL_STMTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "LOG_DISK_CAP", &obj.LOGDISKCAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "LOG_DISK_CAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXAPPLS", &obj.MAXAPPLS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXAPPLS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXFILOP", &obj.MAXFILOP)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXFILOP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAXLOCKS", &obj.MAXLOCKS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAXLOCKS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MIN_DEC_DIV_3", &obj.MINDECDIV3)
	if err != nil {
		err = core.SDKErrorf(err, "", "MIN_DEC_DIV_3-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_ACT_METRICS", &obj.MONACTMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_ACT_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_DEADLOCK", &obj.MONDEADLOCK)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_DEADLOCK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LCK_MSG_LVL", &obj.MONLCKMSGLVL)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LCK_MSG_LVL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LOCKTIMEOUT", &obj.MONLOCKTIMEOUT)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LOCKTIMEOUT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LOCKWAIT", &obj.MONLOCKWAIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LOCKWAIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_LW_THRESH", &obj.MONLWTHRESH)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_LW_THRESH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_OBJ_METRICS", &obj.MONOBJMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_OBJ_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_PKGLIST_SZ", &obj.MONPKGLISTSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_PKGLIST_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_REQ_METRICS", &obj.MONREQMETRICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_REQ_METRICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_RTN_DATA", &obj.MONRTNDATA)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_RTN_DATA-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_RTN_EXECLIST", &obj.MONRTNEXECLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_RTN_EXECLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_DATA", &obj.MONUOWDATA)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_DATA-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_EXECLIST", &obj.MONUOWEXECLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_EXECLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_UOW_PKGLIST", &obj.MONUOWPKGLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_UOW_PKGLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NCHAR_MAPPING", &obj.NCHARMAPPING)
	if err != nil {
		err = core.SDKErrorf(err, "", "NCHAR_MAPPING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_FREQVALUES", &obj.NUMFREQVALUES)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_FREQVALUES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_IOCLEANERS", &obj.NUMIOCLEANERS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_IOCLEANERS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_IOSERVERS", &obj.NUMIOSERVERS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_IOSERVERS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_LOG_SPAN", &obj.NUMLOGSPAN)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_LOG_SPAN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_QUANTILES", &obj.NUMQUANTILES)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_QUANTILES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_BUFFPAGE", &obj.OPTBUFFPAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_BUFFPAGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_DIRECT_WRKLD", &obj.OPTDIRECTWRKLD)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_DIRECT_WRKLD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_LOCKLIST", &obj.OPTLOCKLIST)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_LOCKLIST-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_MAXLOCKS", &obj.OPTMAXLOCKS)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_MAXLOCKS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "OPT_SORTHEAP", &obj.OPTSORTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "OPT_SORTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PAGE_AGE_TRGT_GCR", &obj.PAGEAGETRGTGCR)
	if err != nil {
		err = core.SDKErrorf(err, "", "PAGE_AGE_TRGT_GCR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PAGE_AGE_TRGT_MCR", &obj.PAGEAGETRGTMCR)
	if err != nil {
		err = core.SDKErrorf(err, "", "PAGE_AGE_TRGT_MCR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PCKCACHESZ", &obj.PCKCACHESZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "PCKCACHESZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "PL_STACK_TRACE", &obj.PLSTACKTRACE)
	if err != nil {
		err = core.SDKErrorf(err, "", "PL_STACK_TRACE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SELF_TUNING_MEM", &obj.SELFTUNINGMEM)
	if err != nil {
		err = core.SDKErrorf(err, "", "SELF_TUNING_MEM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SEQDETECT", &obj.SEQDETECT)
	if err != nil {
		err = core.SDKErrorf(err, "", "SEQDETECT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SHEAPTHRES_SHR", &obj.SHEAPTHRESSHR)
	if err != nil {
		err = core.SDKErrorf(err, "", "SHEAPTHRES_SHR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SOFTMAX", &obj.SOFTMAX)
	if err != nil {
		err = core.SDKErrorf(err, "", "SOFTMAX-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SORTHEAP", &obj.SORTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "SORTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SQL_CCFLAGS", &obj.SQLCCFLAGS)
	if err != nil {
		err = core.SDKErrorf(err, "", "SQL_CCFLAGS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STAT_HEAP_SZ", &obj.STATHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "STAT_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STMTHEAP", &obj.STMTHEAP)
	if err != nil {
		err = core.SDKErrorf(err, "", "STMTHEAP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STMT_CONC", &obj.STMTCONC)
	if err != nil {
		err = core.SDKErrorf(err, "", "STMT_CONC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STRING_UNITS", &obj.STRINGUNITS)
	if err != nil {
		err = core.SDKErrorf(err, "", "STRING_UNITS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "SYSTIME_PERIOD_ADJ", &obj.SYSTIMEPERIODADJ)
	if err != nil {
		err = core.SDKErrorf(err, "", "SYSTIME_PERIOD_ADJ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "TRACKMOD", &obj.TRACKMOD)
	if err != nil {
		err = core.SDKErrorf(err, "", "TRACKMOD-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "UTIL_HEAP_SZ", &obj.UTILHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "UTIL_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_ADMISSION_CTRL", &obj.WLMADMISSIONCTRL)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_ADMISSION_CTRL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_AGENT_LOAD_TRGT", &obj.WLMAGENTLOADTRGT)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_AGENT_LOAD_TRGT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_LIMIT", &obj.WLMCPULIMIT)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_LIMIT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_SHARES", &obj.WLMCPUSHARES)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_SHARES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_CPU_SHARE_MODE", &obj.WLMCPUSHAREMODE)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_CPU_SHARE_MODE-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessTuneableParamsTuneableParamDbm : Tunable parameters related to the Db2 instance manager (dbm).
type SuccessTuneableParamsTuneableParamDbm struct {
	COMMBANDWIDTH *string `json:"COMM_BANDWIDTH,omitempty"`

	CPUSPEED *string `json:"CPUSPEED,omitempty"`

	DFTMONBUFPOOL *string `json:"DFT_MON_BUFPOOL,omitempty"`

	DFTMONLOCK *string `json:"DFT_MON_LOCK,omitempty"`

	DFTMONSORT *string `json:"DFT_MON_SORT,omitempty"`

	DFTMONSTMT *string `json:"DFT_MON_STMT,omitempty"`

	DFTMONTABLE *string `json:"DFT_MON_TABLE,omitempty"`

	DFTMONTIMESTAMP *string `json:"DFT_MON_TIMESTAMP,omitempty"`

	DFTMONUOW *string `json:"DFT_MON_UOW,omitempty"`

	DIAGLEVEL *string `json:"DIAGLEVEL,omitempty"`

	FEDERATEDASYNC *string `json:"FEDERATED_ASYNC,omitempty"`

	INDEXREC *string `json:"INDEXREC,omitempty"`

	INTRAPARALLEL *string `json:"INTRA_PARALLEL,omitempty"`

	KEEPFENCED *string `json:"KEEPFENCED,omitempty"`

	MAXCONNRETRIES *string `json:"MAX_CONNRETRIES,omitempty"`

	MAXQUERYDEGREE *string `json:"MAX_QUERYDEGREE,omitempty"`

	MONHEAPSZ *string `json:"MON_HEAP_SZ,omitempty"`

	MULTIPARTSIZEMB *string `json:"MULTIPARTSIZEMB,omitempty"`

	NOTIFYLEVEL *string `json:"NOTIFYLEVEL,omitempty"`

	NUMINITAGENTS *string `json:"NUM_INITAGENTS,omitempty"`

	NUMINITFENCED *string `json:"NUM_INITFENCED,omitempty"`

	NUMPOOLAGENTS *string `json:"NUM_POOLAGENTS,omitempty"`

	RESYNCINTERVAL *string `json:"RESYNC_INTERVAL,omitempty"`

	RQRIOBLK *string `json:"RQRIOBLK,omitempty"`

	STARTSTOPTIME *string `json:"START_STOP_TIME,omitempty"`

	UTILIMPACTLIM *string `json:"UTIL_IMPACT_LIM,omitempty"`

	WLMDISPATCHER *string `json:"WLM_DISPATCHER,omitempty"`

	WLMDISPCONCUR *string `json:"WLM_DISP_CONCUR,omitempty"`

	WLMDISPCPUSHARES *string `json:"WLM_DISP_CPU_SHARES,omitempty"`

	WLMDISPMINUTIL *string `json:"WLM_DISP_MIN_UTIL,omitempty"`
}

// UnmarshalSuccessTuneableParamsTuneableParamDbm unmarshals an instance of SuccessTuneableParamsTuneableParamDbm from the specified map of raw messages.
func UnmarshalSuccessTuneableParamsTuneableParamDbm(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessTuneableParamsTuneableParamDbm)
	err = core.UnmarshalPrimitive(m, "COMM_BANDWIDTH", &obj.COMMBANDWIDTH)
	if err != nil {
		err = core.SDKErrorf(err, "", "COMM_BANDWIDTH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "CPUSPEED", &obj.CPUSPEED)
	if err != nil {
		err = core.SDKErrorf(err, "", "CPUSPEED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_BUFPOOL", &obj.DFTMONBUFPOOL)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_BUFPOOL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_LOCK", &obj.DFTMONLOCK)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_LOCK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_SORT", &obj.DFTMONSORT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_SORT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_STMT", &obj.DFTMONSTMT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_STMT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_TABLE", &obj.DFTMONTABLE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_TABLE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_TIMESTAMP", &obj.DFTMONTIMESTAMP)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_TIMESTAMP-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DFT_MON_UOW", &obj.DFTMONUOW)
	if err != nil {
		err = core.SDKErrorf(err, "", "DFT_MON_UOW-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DIAGLEVEL", &obj.DIAGLEVEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "DIAGLEVEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "FEDERATED_ASYNC", &obj.FEDERATEDASYNC)
	if err != nil {
		err = core.SDKErrorf(err, "", "FEDERATED_ASYNC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INDEXREC", &obj.INDEXREC)
	if err != nil {
		err = core.SDKErrorf(err, "", "INDEXREC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "INTRA_PARALLEL", &obj.INTRAPARALLEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "INTRA_PARALLEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "KEEPFENCED", &obj.KEEPFENCED)
	if err != nil {
		err = core.SDKErrorf(err, "", "KEEPFENCED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAX_CONNRETRIES", &obj.MAXCONNRETRIES)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAX_CONNRETRIES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MAX_QUERYDEGREE", &obj.MAXQUERYDEGREE)
	if err != nil {
		err = core.SDKErrorf(err, "", "MAX_QUERYDEGREE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MON_HEAP_SZ", &obj.MONHEAPSZ)
	if err != nil {
		err = core.SDKErrorf(err, "", "MON_HEAP_SZ-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "MULTIPARTSIZEMB", &obj.MULTIPARTSIZEMB)
	if err != nil {
		err = core.SDKErrorf(err, "", "MULTIPARTSIZEMB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NOTIFYLEVEL", &obj.NOTIFYLEVEL)
	if err != nil {
		err = core.SDKErrorf(err, "", "NOTIFYLEVEL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_INITAGENTS", &obj.NUMINITAGENTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_INITAGENTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_INITFENCED", &obj.NUMINITFENCED)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_INITFENCED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "NUM_POOLAGENTS", &obj.NUMPOOLAGENTS)
	if err != nil {
		err = core.SDKErrorf(err, "", "NUM_POOLAGENTS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "RESYNC_INTERVAL", &obj.RESYNCINTERVAL)
	if err != nil {
		err = core.SDKErrorf(err, "", "RESYNC_INTERVAL-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "RQRIOBLK", &obj.RQRIOBLK)
	if err != nil {
		err = core.SDKErrorf(err, "", "RQRIOBLK-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "START_STOP_TIME", &obj.STARTSTOPTIME)
	if err != nil {
		err = core.SDKErrorf(err, "", "START_STOP_TIME-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "UTIL_IMPACT_LIM", &obj.UTILIMPACTLIM)
	if err != nil {
		err = core.SDKErrorf(err, "", "UTIL_IMPACT_LIM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISPATCHER", &obj.WLMDISPATCHER)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISPATCHER-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_CONCUR", &obj.WLMDISPCONCUR)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_CONCUR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_CPU_SHARES", &obj.WLMDISPCPUSHARES)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_CPU_SHARES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "WLM_DISP_MIN_UTIL", &obj.WLMDISPMINUTIL)
	if err != nil {
		err = core.SDKErrorf(err, "", "WLM_DISP_MIN_UTIL-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessTuneableParamsTuneableParamRegistry : Tunable parameters related to the Db2 registry.
type SuccessTuneableParamsTuneableParamRegistry struct {
	DB2BIDI *string `json:"DB2BIDI,omitempty"`

	DB2COMPOPT *string `json:"DB2COMPOPT,omitempty"`

	DB2LOCKTORB *string `json:"DB2LOCK_TO_RB,omitempty"`

	DB2STMM *string `json:"DB2STMM,omitempty"`

	DB2ALTERNATEAUTHZBEHAVIOUR *string `json:"DB2_ALTERNATE_AUTHZ_BEHAVIOUR,omitempty"`

	DB2ANTIJOIN *string `json:"DB2_ANTIJOIN,omitempty"`

	DB2ATSENABLE *string `json:"DB2_ATS_ENABLE,omitempty"`

	DB2DEFERREDPREPARESEMANTICS *string `json:"DB2_DEFERRED_PREPARE_SEMANTICS,omitempty"`

	DB2EVALUNCOMMITTED *string `json:"DB2_EVALUNCOMMITTED,omitempty"`

	DB2EXTENDEDOPTIMIZATION *string `json:"DB2_EXTENDED_OPTIMIZATION,omitempty"`

	DB2INDEXPCTFREEDEFAULT *string `json:"DB2_INDEX_PCTFREE_DEFAULT,omitempty"`

	DB2INLISTTONLJN *string `json:"DB2_INLIST_TO_NLJN,omitempty"`

	DB2MINIMIZELISTPREFETCH *string `json:"DB2_MINIMIZE_LISTPREFETCH,omitempty"`

	DB2OBJECTTABLEENTRIES *string `json:"DB2_OBJECT_TABLE_ENTRIES,omitempty"`

	DB2OPTPROFILE *string `json:"DB2_OPTPROFILE,omitempty"`

	DB2OPTSTATSLOG *string `json:"DB2_OPTSTATS_LOG,omitempty"`

	DB2OPTMAXTEMPSIZE *string `json:"DB2_OPT_MAX_TEMP_SIZE,omitempty"`

	DB2PARALLELIO *string `json:"DB2_PARALLEL_IO,omitempty"`

	DB2REDUCEDOPTIMIZATION *string `json:"DB2_REDUCED_OPTIMIZATION,omitempty"`

	DB2SELECTIVITY *string `json:"DB2_SELECTIVITY,omitempty"`

	DB2SKIPDELETED *string `json:"DB2_SKIPDELETED,omitempty"`

	DB2SKIPINSERTED *string `json:"DB2_SKIPINSERTED,omitempty"`

	DB2SYNCRELEASELOCKATTRIBUTES *string `json:"DB2_SYNC_RELEASE_LOCK_ATTRIBUTES,omitempty"`

	DB2TRUNCATEREUSESTORAGE *string `json:"DB2_TRUNCATE_REUSESTORAGE,omitempty"`

	DB2USEALTERNATEPAGECLEANING *string `json:"DB2_USE_ALTERNATE_PAGE_CLEANING,omitempty"`

	DB2VIEWREOPTVALUES *string `json:"DB2_VIEW_REOPT_VALUES,omitempty"`

	DB2WLMSETTINGS *string `json:"DB2_WLM_SETTINGS,omitempty"`

	DB2WORKLOAD *string `json:"DB2_WORKLOAD,omitempty"`
}

// UnmarshalSuccessTuneableParamsTuneableParamRegistry unmarshals an instance of SuccessTuneableParamsTuneableParamRegistry from the specified map of raw messages.
func UnmarshalSuccessTuneableParamsTuneableParamRegistry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessTuneableParamsTuneableParamRegistry)
	err = core.UnmarshalPrimitive(m, "DB2BIDI", &obj.DB2BIDI)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2BIDI-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2COMPOPT", &obj.DB2COMPOPT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2COMPOPT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2LOCK_TO_RB", &obj.DB2LOCKTORB)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2LOCK_TO_RB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2STMM", &obj.DB2STMM)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2STMM-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ALTERNATE_AUTHZ_BEHAVIOUR", &obj.DB2ALTERNATEAUTHZBEHAVIOUR)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ALTERNATE_AUTHZ_BEHAVIOUR-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ANTIJOIN", &obj.DB2ANTIJOIN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ANTIJOIN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_ATS_ENABLE", &obj.DB2ATSENABLE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_ATS_ENABLE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_DEFERRED_PREPARE_SEMANTICS", &obj.DB2DEFERREDPREPARESEMANTICS)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_DEFERRED_PREPARE_SEMANTICS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_EVALUNCOMMITTED", &obj.DB2EVALUNCOMMITTED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_EVALUNCOMMITTED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_EXTENDED_OPTIMIZATION", &obj.DB2EXTENDEDOPTIMIZATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_EXTENDED_OPTIMIZATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_INDEX_PCTFREE_DEFAULT", &obj.DB2INDEXPCTFREEDEFAULT)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_INDEX_PCTFREE_DEFAULT-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_INLIST_TO_NLJN", &obj.DB2INLISTTONLJN)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_INLIST_TO_NLJN-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_MINIMIZE_LISTPREFETCH", &obj.DB2MINIMIZELISTPREFETCH)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_MINIMIZE_LISTPREFETCH-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OBJECT_TABLE_ENTRIES", &obj.DB2OBJECTTABLEENTRIES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OBJECT_TABLE_ENTRIES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPTPROFILE", &obj.DB2OPTPROFILE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPTPROFILE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPTSTATS_LOG", &obj.DB2OPTSTATSLOG)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPTSTATS_LOG-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_OPT_MAX_TEMP_SIZE", &obj.DB2OPTMAXTEMPSIZE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_OPT_MAX_TEMP_SIZE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_PARALLEL_IO", &obj.DB2PARALLELIO)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_PARALLEL_IO-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_REDUCED_OPTIMIZATION", &obj.DB2REDUCEDOPTIMIZATION)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_REDUCED_OPTIMIZATION-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SELECTIVITY", &obj.DB2SELECTIVITY)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SELECTIVITY-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SKIPDELETED", &obj.DB2SKIPDELETED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SKIPDELETED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SKIPINSERTED", &obj.DB2SKIPINSERTED)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SKIPINSERTED-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_SYNC_RELEASE_LOCK_ATTRIBUTES", &obj.DB2SYNCRELEASELOCKATTRIBUTES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_SYNC_RELEASE_LOCK_ATTRIBUTES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_TRUNCATE_REUSESTORAGE", &obj.DB2TRUNCATEREUSESTORAGE)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_TRUNCATE_REUSESTORAGE-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_USE_ALTERNATE_PAGE_CLEANING", &obj.DB2USEALTERNATEPAGECLEANING)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_USE_ALTERNATE_PAGE_CLEANING-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_VIEW_REOPT_VALUES", &obj.DB2VIEWREOPTVALUES)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_VIEW_REOPT_VALUES-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_WLM_SETTINGS", &obj.DB2WLMSETTINGS)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_WLM_SETTINGS-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "DB2_WORKLOAD", &obj.DB2WORKLOAD)
	if err != nil {
		err = core.SDKErrorf(err, "", "DB2_WORKLOAD-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessUpdateAutoScale : Response of successful updation of scaling configurations.
type SuccessUpdateAutoScale struct {
	// Indicates the message of the updation.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalSuccessUpdateAutoScale unmarshals an instance of SuccessUpdateAutoScale from the specified map of raw messages.
func UnmarshalSuccessUpdateAutoScale(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessUpdateAutoScale)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessUserResponse : The details of the users.
type SuccessUserResponse struct {
	// User's DV role.
	DvRole *string `json:"dvRole" validate:"required"`

	// Metadata associated with the user.
	Metadata map[string]interface{} `json:"metadata" validate:"required"`

	// Formatted IBM ID.
	FormatedIbmid *string `json:"formatedIbmid" validate:"required"`

	// Role assigned to the user.
	Role *string `json:"role" validate:"required"`

	// IAM ID for the user.
	Iamid *string `json:"iamid" validate:"required"`

	// List of allowed actions of the user.
	PermittedActions []string `json:"permittedActions" validate:"required"`

	// Indicates if the user account has no issues.
	AllClean *bool `json:"allClean" validate:"required"`

	// User's password.
	Password *string `json:"password" validate:"required"`

	// Indicates if IAM is enabled or not.
	Iam *bool `json:"iam" validate:"required"`

	// The display name of the user.
	Name *string `json:"name" validate:"required"`

	// IBM ID of the user.
	Ibmid *string `json:"ibmid" validate:"required"`

	// Unique identifier for the user.
	ID *string `json:"id" validate:"required"`

	// Account lock status for the user.
	Locked *string `json:"locked" validate:"required"`

	// Initial error message.
	InitErrorMsg *string `json:"initErrorMsg" validate:"required"`

	// Email address of the user.
	Email *string `json:"email" validate:"required"`

	// Authentication details for the user.
	Authentication *SuccessUserResponseAuthentication `json:"authentication" validate:"required"`
}

// Constants associated with the SuccessUserResponse.Role property.
// Role assigned to the user.
const (
	SuccessUserResponse_Role_Bluadmin = "bluadmin"
	SuccessUserResponse_Role_Bluuser = "bluuser"
)

// Constants associated with the SuccessUserResponse.Locked property.
// Account lock status for the user.
const (
	SuccessUserResponse_Locked_No = "no"
	SuccessUserResponse_Locked_Yes = "yes"
)

// UnmarshalSuccessUserResponse unmarshals an instance of SuccessUserResponse from the specified map of raw messages.
func UnmarshalSuccessUserResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessUserResponse)
	err = core.UnmarshalPrimitive(m, "dvRole", &obj.DvRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "dvRole-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata", &obj.Metadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "formatedIbmid", &obj.FormatedIbmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "formatedIbmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "role", &obj.Role)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iamid", &obj.Iamid)
	if err != nil {
		err = core.SDKErrorf(err, "", "iamid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "permittedActions", &obj.PermittedActions)
	if err != nil {
		err = core.SDKErrorf(err, "", "permittedActions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allClean", &obj.AllClean)
	if err != nil {
		err = core.SDKErrorf(err, "", "allClean-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam", &obj.Iam)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ibmid", &obj.Ibmid)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initErrorMsg", &obj.InitErrorMsg)
	if err != nil {
		err = core.SDKErrorf(err, "", "initErrorMsg-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "authentication", &obj.Authentication, UnmarshalSuccessUserResponseAuthentication)
	if err != nil {
		err = core.SDKErrorf(err, "", "authentication-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SuccessUserResponseAuthentication : Authentication details for the user.
type SuccessUserResponseAuthentication struct {
	// Authentication method.
	Method *string `json:"method" validate:"required"`

	// Policy ID of authentication.
	PolicyID *string `json:"policy_id" validate:"required"`
}

// UnmarshalSuccessUserResponseAuthentication unmarshals an instance of SuccessUserResponseAuthentication from the specified map of raw messages.
func UnmarshalSuccessUserResponseAuthentication(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SuccessUserResponseAuthentication)
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		err = core.SDKErrorf(err, "", "method-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_id", &obj.PolicyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
