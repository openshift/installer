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
 * IBM OpenAPI SDK Code Generator Version: 3.84.0-a4533f12-20240103-170852
 */

// Package rulesetsv1 : Operations and models for the RulesetsV1 service
package rulesetsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// RulesetsV1 : Rulesets Engine
//
// API Version: 1.0.1
type RulesetsV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "rulesets"

// RulesetsV1Options : Service options
type RulesetsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewRulesetsV1UsingExternalConfig : constructs an instance of RulesetsV1 with passed in options and external configuration.
func NewRulesetsV1UsingExternalConfig(options *RulesetsV1Options) (rulesets *RulesetsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	rulesets, err = NewRulesetsV1(options)
	if err != nil {
		return
	}

	err = rulesets.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = rulesets.Service.SetServiceURL(options.URL)
	}
	return
}

// NewRulesetsV1 : constructs an instance of RulesetsV1 with passed in options.
func NewRulesetsV1(options *RulesetsV1Options) (service *RulesetsV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	err = core.ValidateStruct(options, "options")
	if err != nil {
		return
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

	service = &RulesetsV1{
		Service:        baseService,
		Crn:            options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "rulesets" suitable for processing requests.
func (rulesets *RulesetsV1) Clone() *RulesetsV1 {
	if core.IsNil(rulesets) {
		return nil
	}
	clone := *rulesets
	clone.Service = rulesets.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (rulesets *RulesetsV1) SetServiceURL(url string) error {
	return rulesets.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (rulesets *RulesetsV1) GetServiceURL() string {
	return rulesets.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (rulesets *RulesetsV1) SetDefaultHeaders(headers http.Header) {
	rulesets.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (rulesets *RulesetsV1) SetEnableGzipCompression(enableGzip bool) {
	rulesets.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (rulesets *RulesetsV1) GetEnableGzipCompression() bool {
	return rulesets.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (rulesets *RulesetsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	rulesets.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (rulesets *RulesetsV1) DisableRetries() {
	rulesets.Service.DisableRetries()
}

// GetInstanceRulesets : List Instance rulesets
// List all rulesets at the instance level.
func (rulesets *RulesetsV1) GetInstanceRulesets(getInstanceRulesetsOptions *GetInstanceRulesetsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceRulesetsWithContext(context.Background(), getInstanceRulesetsOptions)
}

// GetInstanceRulesetsWithContext is an alternate form of the GetInstanceRulesets method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceRulesetsWithContext(ctx context.Context, getInstanceRulesetsOptions *GetInstanceRulesetsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getInstanceRulesetsOptions, "getInstanceRulesetsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rulesets.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceRulesetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceRulesets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInstanceRuleset : Get an instance ruleset
// View a specific instance ruleset.
func (rulesets *RulesetsV1) GetInstanceRuleset(getInstanceRulesetOptions *GetInstanceRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceRulesetWithContext(context.Background(), getInstanceRulesetOptions)
}

// GetInstanceRulesetWithContext is an alternate form of the GetInstanceRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceRulesetWithContext(ctx context.Context, getInstanceRulesetOptions *GetInstanceRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceRulesetOptions, "getInstanceRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceRulesetOptions, "getInstanceRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *getInstanceRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateInstanceRuleset : Update an instance ruleset
// Update a specific instance ruleset.
func (rulesets *RulesetsV1) UpdateInstanceRuleset(updateInstanceRulesetOptions *UpdateInstanceRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateInstanceRulesetWithContext(context.Background(), updateInstanceRulesetOptions)
}

// UpdateInstanceRulesetWithContext is an alternate form of the UpdateInstanceRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateInstanceRulesetWithContext(ctx context.Context, updateInstanceRulesetOptions *UpdateInstanceRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateInstanceRulesetOptions, "updateInstanceRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateInstanceRulesetOptions, "updateInstanceRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *updateInstanceRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateInstanceRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateInstanceRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateInstanceRulesetOptions.Description != nil {
		body["description"] = updateInstanceRulesetOptions.Description
	}
	if updateInstanceRulesetOptions.Kind != nil {
		body["kind"] = updateInstanceRulesetOptions.Kind
	}
	if updateInstanceRulesetOptions.Name != nil {
		body["name"] = updateInstanceRulesetOptions.Name
	}
	if updateInstanceRulesetOptions.Phase != nil {
		body["phase"] = updateInstanceRulesetOptions.Phase
	}
	if updateInstanceRulesetOptions.Rules != nil {
		body["rules"] = updateInstanceRulesetOptions.Rules
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteInstanceRuleset : Delete an instance ruleset
// Delete a specific instance ruleset.
func (rulesets *RulesetsV1) DeleteInstanceRuleset(deleteInstanceRulesetOptions *DeleteInstanceRulesetOptions) (response *core.DetailedResponse, err error) {
	return rulesets.DeleteInstanceRulesetWithContext(context.Background(), deleteInstanceRulesetOptions)
}

// DeleteInstanceRulesetWithContext is an alternate form of the DeleteInstanceRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteInstanceRulesetWithContext(ctx context.Context, deleteInstanceRulesetOptions *DeleteInstanceRulesetOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteInstanceRulesetOptions, "deleteInstanceRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteInstanceRulesetOptions, "deleteInstanceRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *deleteInstanceRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteInstanceRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteInstanceRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = rulesets.Service.Request(request, nil)

	return
}

// GetInstanceRulesetVersions : List version of an instance ruleset
// List all versions of a specific instance ruleset.
func (rulesets *RulesetsV1) GetInstanceRulesetVersions(getInstanceRulesetVersionsOptions *GetInstanceRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceRulesetVersionsWithContext(context.Background(), getInstanceRulesetVersionsOptions)
}

// GetInstanceRulesetVersionsWithContext is an alternate form of the GetInstanceRulesetVersions method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceRulesetVersionsWithContext(ctx context.Context, getInstanceRulesetVersionsOptions *GetInstanceRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceRulesetVersionsOptions, "getInstanceRulesetVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceRulesetVersionsOptions, "getInstanceRulesetVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *getInstanceRulesetVersionsOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceRulesetVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceRulesetVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInstanceRulesetVersion : Get a specific version of an instance ruleset
// View a specific version of a specific instance ruleset.
func (rulesets *RulesetsV1) GetInstanceRulesetVersion(getInstanceRulesetVersionOptions *GetInstanceRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceRulesetVersionWithContext(context.Background(), getInstanceRulesetVersionOptions)
}

// GetInstanceRulesetVersionWithContext is an alternate form of the GetInstanceRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceRulesetVersionWithContext(ctx context.Context, getInstanceRulesetVersionOptions *GetInstanceRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceRulesetVersionOptions, "getInstanceRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceRulesetVersionOptions, "getInstanceRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"ruleset_id":      *getInstanceRulesetVersionOptions.RulesetID,
		"ruleset_version": *getInstanceRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteInstanceRulesetVersion : Delete a specific version of an instance ruleset
// Delete a specific version of a specific instance ruleset.
func (rulesets *RulesetsV1) DeleteInstanceRulesetVersion(deleteInstanceRulesetVersionOptions *DeleteInstanceRulesetVersionOptions) (response *core.DetailedResponse, err error) {
	return rulesets.DeleteInstanceRulesetVersionWithContext(context.Background(), deleteInstanceRulesetVersionOptions)
}

// DeleteInstanceRulesetVersionWithContext is an alternate form of the DeleteInstanceRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteInstanceRulesetVersionWithContext(ctx context.Context, deleteInstanceRulesetVersionOptions *DeleteInstanceRulesetVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteInstanceRulesetVersionOptions, "deleteInstanceRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteInstanceRulesetVersionOptions, "deleteInstanceRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"ruleset_id":      *deleteInstanceRulesetVersionOptions.RulesetID,
		"ruleset_version": *deleteInstanceRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteInstanceRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteInstanceRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = rulesets.Service.Request(request, nil)

	return
}

// GetInstanceEntrypointRuleset : Get an instance entrypoint ruleset
// Get the instance ruleset for the given phase's entrypoint.
func (rulesets *RulesetsV1) GetInstanceEntrypointRuleset(getInstanceEntrypointRulesetOptions *GetInstanceEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceEntrypointRulesetWithContext(context.Background(), getInstanceEntrypointRulesetOptions)
}

// GetInstanceEntrypointRulesetWithContext is an alternate form of the GetInstanceEntrypointRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceEntrypointRulesetWithContext(ctx context.Context, getInstanceEntrypointRulesetOptions *GetInstanceEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceEntrypointRulesetOptions, "getInstanceEntrypointRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceEntrypointRulesetOptions, "getInstanceEntrypointRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":           *rulesets.Crn,
		"ruleset_phase": *getInstanceEntrypointRulesetOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/phases/{ruleset_phase}/entrypoint`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceEntrypointRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceEntrypointRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateInstanceEntrypointRuleset : Update an instance entrypoint ruleset
// Updates the instance ruleset for the given phase's entry point.
func (rulesets *RulesetsV1) UpdateInstanceEntrypointRuleset(updateInstanceEntrypointRulesetOptions *UpdateInstanceEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateInstanceEntrypointRulesetWithContext(context.Background(), updateInstanceEntrypointRulesetOptions)
}

// UpdateInstanceEntrypointRulesetWithContext is an alternate form of the UpdateInstanceEntrypointRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateInstanceEntrypointRulesetWithContext(ctx context.Context, updateInstanceEntrypointRulesetOptions *UpdateInstanceEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateInstanceEntrypointRulesetOptions, "updateInstanceEntrypointRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateInstanceEntrypointRulesetOptions, "updateInstanceEntrypointRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":           *rulesets.Crn,
		"ruleset_phase": *updateInstanceEntrypointRulesetOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/phases/{ruleset_phase}/entrypoint`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateInstanceEntrypointRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateInstanceEntrypointRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateInstanceEntrypointRulesetOptions.Description != nil {
		body["description"] = updateInstanceEntrypointRulesetOptions.Description
	}
	if updateInstanceEntrypointRulesetOptions.Kind != nil {
		body["kind"] = updateInstanceEntrypointRulesetOptions.Kind
	}
	if updateInstanceEntrypointRulesetOptions.Name != nil {
		body["name"] = updateInstanceEntrypointRulesetOptions.Name
	}
	if updateInstanceEntrypointRulesetOptions.Phase != nil {
		body["phase"] = updateInstanceEntrypointRulesetOptions.Phase
	}
	if updateInstanceEntrypointRulesetOptions.Rules != nil {
		body["rules"] = updateInstanceEntrypointRulesetOptions.Rules
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInstanceEntryPointRulesetVersions : List an instance entry point ruleset's versions
// Lists the instance ruleset versions for the given phase's entry point.
func (rulesets *RulesetsV1) GetInstanceEntryPointRulesetVersions(getInstanceEntryPointRulesetVersionsOptions *GetInstanceEntryPointRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceEntryPointRulesetVersionsWithContext(context.Background(), getInstanceEntryPointRulesetVersionsOptions)
}

// GetInstanceEntryPointRulesetVersionsWithContext is an alternate form of the GetInstanceEntryPointRulesetVersions method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceEntryPointRulesetVersionsWithContext(ctx context.Context, getInstanceEntryPointRulesetVersionsOptions *GetInstanceEntryPointRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceEntryPointRulesetVersionsOptions, "getInstanceEntryPointRulesetVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceEntryPointRulesetVersionsOptions, "getInstanceEntryPointRulesetVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":           *rulesets.Crn,
		"ruleset_phase": *getInstanceEntryPointRulesetVersionsOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/phases/{ruleset_phase}/entrypoint/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceEntryPointRulesetVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceEntryPointRulesetVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInstanceEntryPointRulesetVersion : Get an instance entry point ruleset version
// Fetches a specific version of an instance entry point ruleset.
func (rulesets *RulesetsV1) GetInstanceEntryPointRulesetVersion(getInstanceEntryPointRulesetVersionOptions *GetInstanceEntryPointRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceEntryPointRulesetVersionWithContext(context.Background(), getInstanceEntryPointRulesetVersionOptions)
}

// GetInstanceEntryPointRulesetVersionWithContext is an alternate form of the GetInstanceEntryPointRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceEntryPointRulesetVersionWithContext(ctx context.Context, getInstanceEntryPointRulesetVersionOptions *GetInstanceEntryPointRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceEntryPointRulesetVersionOptions, "getInstanceEntryPointRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceEntryPointRulesetVersionOptions, "getInstanceEntryPointRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"ruleset_phase":   *getInstanceEntryPointRulesetVersionOptions.RulesetPhase,
		"ruleset_version": *getInstanceEntryPointRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/phases/{ruleset_phase}/entrypoint/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceEntryPointRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceEntryPointRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateInstanceRulesetRule : Create an instance ruleset rule
// Create an instance ruleset rule.
func (rulesets *RulesetsV1) CreateInstanceRulesetRule(createInstanceRulesetRuleOptions *CreateInstanceRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.CreateInstanceRulesetRuleWithContext(context.Background(), createInstanceRulesetRuleOptions)
}

// CreateInstanceRulesetRuleWithContext is an alternate form of the CreateInstanceRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) CreateInstanceRulesetRuleWithContext(ctx context.Context, createInstanceRulesetRuleOptions *CreateInstanceRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createInstanceRulesetRuleOptions, "createInstanceRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createInstanceRulesetRuleOptions, "createInstanceRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *createInstanceRulesetRuleOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createInstanceRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "CreateInstanceRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createInstanceRulesetRuleOptions.Action != nil {
		body["action"] = createInstanceRulesetRuleOptions.Action
	}
	if createInstanceRulesetRuleOptions.ActionParameters != nil {
		body["action_parameters"] = createInstanceRulesetRuleOptions.ActionParameters
	}
	if createInstanceRulesetRuleOptions.Description != nil {
		body["description"] = createInstanceRulesetRuleOptions.Description
	}
	if createInstanceRulesetRuleOptions.Enabled != nil {
		body["enabled"] = createInstanceRulesetRuleOptions.Enabled
	}
	if createInstanceRulesetRuleOptions.Expression != nil {
		body["expression"] = createInstanceRulesetRuleOptions.Expression
	}
	if createInstanceRulesetRuleOptions.ID != nil {
		body["id"] = createInstanceRulesetRuleOptions.ID
	}
	if createInstanceRulesetRuleOptions.Logging != nil {
		body["logging"] = createInstanceRulesetRuleOptions.Logging
	}
	if createInstanceRulesetRuleOptions.Ref != nil {
		body["ref"] = createInstanceRulesetRuleOptions.Ref
	}
	if createInstanceRulesetRuleOptions.Position != nil {
		body["position"] = createInstanceRulesetRuleOptions.Position
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateInstanceRulesetRule : Update an instance ruleset rule
// Update an instance ruleset rule.
func (rulesets *RulesetsV1) UpdateInstanceRulesetRule(updateInstanceRulesetRuleOptions *UpdateInstanceRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateInstanceRulesetRuleWithContext(context.Background(), updateInstanceRulesetRuleOptions)
}

// UpdateInstanceRulesetRuleWithContext is an alternate form of the UpdateInstanceRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateInstanceRulesetRuleWithContext(ctx context.Context, updateInstanceRulesetRuleOptions *UpdateInstanceRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateInstanceRulesetRuleOptions, "updateInstanceRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateInstanceRulesetRuleOptions, "updateInstanceRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *updateInstanceRulesetRuleOptions.RulesetID,
		"rule_id":    *updateInstanceRulesetRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateInstanceRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateInstanceRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateInstanceRulesetRuleOptions.Action != nil {
		body["action"] = updateInstanceRulesetRuleOptions.Action
	}
	if updateInstanceRulesetRuleOptions.ActionParameters != nil {
		body["action_parameters"] = updateInstanceRulesetRuleOptions.ActionParameters
	}
	if updateInstanceRulesetRuleOptions.Description != nil {
		body["description"] = updateInstanceRulesetRuleOptions.Description
	}
	if updateInstanceRulesetRuleOptions.Enabled != nil {
		body["enabled"] = updateInstanceRulesetRuleOptions.Enabled
	}
	if updateInstanceRulesetRuleOptions.Expression != nil {
		body["expression"] = updateInstanceRulesetRuleOptions.Expression
	}
	if updateInstanceRulesetRuleOptions.ID != nil {
		body["id"] = updateInstanceRulesetRuleOptions.ID
	}
	if updateInstanceRulesetRuleOptions.Logging != nil {
		body["logging"] = updateInstanceRulesetRuleOptions.Logging
	}
	if updateInstanceRulesetRuleOptions.Ref != nil {
		body["ref"] = updateInstanceRulesetRuleOptions.Ref
	}
	if updateInstanceRulesetRuleOptions.Position != nil {
		body["position"] = updateInstanceRulesetRuleOptions.Position
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteInstanceRulesetRule : Delete an instance ruleset rule
// Delete an instance ruleset rule.
func (rulesets *RulesetsV1) DeleteInstanceRulesetRule(deleteInstanceRulesetRuleOptions *DeleteInstanceRulesetRuleOptions) (result *RuleResp, response *core.DetailedResponse, err error) {
	return rulesets.DeleteInstanceRulesetRuleWithContext(context.Background(), deleteInstanceRulesetRuleOptions)
}

// DeleteInstanceRulesetRuleWithContext is an alternate form of the DeleteInstanceRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteInstanceRulesetRuleWithContext(ctx context.Context, deleteInstanceRulesetRuleOptions *DeleteInstanceRulesetRuleOptions) (result *RuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteInstanceRulesetRuleOptions, "deleteInstanceRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteInstanceRulesetRuleOptions, "deleteInstanceRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":        *rulesets.Crn,
		"ruleset_id": *deleteInstanceRulesetRuleOptions.RulesetID,
		"rule_id":    *deleteInstanceRulesetRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteInstanceRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteInstanceRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetInstanceRulesetVersionByTag : List an instance ruleset verion's rules by tag
// Lists rules by tag for a specific version of an instance ruleset.
func (rulesets *RulesetsV1) GetInstanceRulesetVersionByTag(getInstanceRulesetVersionByTagOptions *GetInstanceRulesetVersionByTagOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetInstanceRulesetVersionByTagWithContext(context.Background(), getInstanceRulesetVersionByTagOptions)
}

// GetInstanceRulesetVersionByTagWithContext is an alternate form of the GetInstanceRulesetVersionByTag method which supports a Context parameter
func (rulesets *RulesetsV1) GetInstanceRulesetVersionByTagWithContext(ctx context.Context, getInstanceRulesetVersionByTagOptions *GetInstanceRulesetVersionByTagOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceRulesetVersionByTagOptions, "getInstanceRulesetVersionByTagOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceRulesetVersionByTagOptions, "getInstanceRulesetVersionByTagOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"ruleset_id":      *getInstanceRulesetVersionByTagOptions.RulesetID,
		"ruleset_version": *getInstanceRulesetVersionByTagOptions.RulesetVersion,
		"rule_tag":        *getInstanceRulesetVersionByTagOptions.RuleTag,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/rulesets/{ruleset_id}/versions/{ruleset_version}/by_tag/{rule_tag}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceRulesetVersionByTagOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetInstanceRulesetVersionByTag")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneRulesets : List zone rulesets
// List all rulesets at the zone level.
func (rulesets *RulesetsV1) GetZoneRulesets(getZoneRulesetsOptions *GetZoneRulesetsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneRulesetsWithContext(context.Background(), getZoneRulesetsOptions)
}

// GetZoneRulesetsWithContext is an alternate form of the GetZoneRulesets method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneRulesetsWithContext(ctx context.Context, getZoneRulesetsOptions *GetZoneRulesetsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getZoneRulesetsOptions, "getZoneRulesetsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneRulesetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneRulesets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneRuleset : Get a zone ruleset
// View a specific zone ruleset.
func (rulesets *RulesetsV1) GetZoneRuleset(getZoneRulesetOptions *GetZoneRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneRulesetWithContext(context.Background(), getZoneRulesetOptions)
}

// GetZoneRulesetWithContext is an alternate form of the GetZoneRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneRulesetWithContext(ctx context.Context, getZoneRulesetOptions *GetZoneRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneRulesetOptions, "getZoneRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneRulesetOptions, "getZoneRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *getZoneRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateZoneRuleset : Update a zone ruleset
// Update a specific zone ruleset.
func (rulesets *RulesetsV1) UpdateZoneRuleset(updateZoneRulesetOptions *UpdateZoneRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateZoneRulesetWithContext(context.Background(), updateZoneRulesetOptions)
}

// UpdateZoneRulesetWithContext is an alternate form of the UpdateZoneRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateZoneRulesetWithContext(ctx context.Context, updateZoneRulesetOptions *UpdateZoneRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateZoneRulesetOptions, "updateZoneRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateZoneRulesetOptions, "updateZoneRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *updateZoneRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateZoneRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateZoneRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateZoneRulesetOptions.Description != nil {
		body["description"] = updateZoneRulesetOptions.Description
	}
	if updateZoneRulesetOptions.Kind != nil {
		body["kind"] = updateZoneRulesetOptions.Kind
	}
	if updateZoneRulesetOptions.Name != nil {
		body["name"] = updateZoneRulesetOptions.Name
	}
	if updateZoneRulesetOptions.Phase != nil {
		body["phase"] = updateZoneRulesetOptions.Phase
	}
	if updateZoneRulesetOptions.Rules != nil {
		body["rules"] = updateZoneRulesetOptions.Rules
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteZoneRuleset : Delete a zone ruleset
// Delete a specific zone ruleset.
func (rulesets *RulesetsV1) DeleteZoneRuleset(deleteZoneRulesetOptions *DeleteZoneRulesetOptions) (response *core.DetailedResponse, err error) {
	return rulesets.DeleteZoneRulesetWithContext(context.Background(), deleteZoneRulesetOptions)
}

// DeleteZoneRulesetWithContext is an alternate form of the DeleteZoneRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteZoneRulesetWithContext(ctx context.Context, deleteZoneRulesetOptions *DeleteZoneRulesetOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneRulesetOptions, "deleteZoneRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneRulesetOptions, "deleteZoneRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *deleteZoneRulesetOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteZoneRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = rulesets.Service.Request(request, nil)

	return
}

// GetZoneRulesetVersions : List version of a zone ruleset
// List all versions of a specific zone ruleset.
func (rulesets *RulesetsV1) GetZoneRulesetVersions(getZoneRulesetVersionsOptions *GetZoneRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneRulesetVersionsWithContext(context.Background(), getZoneRulesetVersionsOptions)
}

// GetZoneRulesetVersionsWithContext is an alternate form of the GetZoneRulesetVersions method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneRulesetVersionsWithContext(ctx context.Context, getZoneRulesetVersionsOptions *GetZoneRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneRulesetVersionsOptions, "getZoneRulesetVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneRulesetVersionsOptions, "getZoneRulesetVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *getZoneRulesetVersionsOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneRulesetVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneRulesetVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneRulesetVersion : Get a specific version of a zone ruleset
// View a specific version of a specific zone ruleset.
func (rulesets *RulesetsV1) GetZoneRulesetVersion(getZoneRulesetVersionOptions *GetZoneRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneRulesetVersionWithContext(context.Background(), getZoneRulesetVersionOptions)
}

// GetZoneRulesetVersionWithContext is an alternate form of the GetZoneRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneRulesetVersionWithContext(ctx context.Context, getZoneRulesetVersionOptions *GetZoneRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneRulesetVersionOptions, "getZoneRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneRulesetVersionOptions, "getZoneRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *getZoneRulesetVersionOptions.RulesetID,
		"ruleset_version": *getZoneRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteZoneRulesetVersion : Delete a specific version of a zone ruleset
// Delete a specific version of a specific zone ruleset.
func (rulesets *RulesetsV1) DeleteZoneRulesetVersion(deleteZoneRulesetVersionOptions *DeleteZoneRulesetVersionOptions) (response *core.DetailedResponse, err error) {
	return rulesets.DeleteZoneRulesetVersionWithContext(context.Background(), deleteZoneRulesetVersionOptions)
}

// DeleteZoneRulesetVersionWithContext is an alternate form of the DeleteZoneRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteZoneRulesetVersionWithContext(ctx context.Context, deleteZoneRulesetVersionOptions *DeleteZoneRulesetVersionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneRulesetVersionOptions, "deleteZoneRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneRulesetVersionOptions, "deleteZoneRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *deleteZoneRulesetVersionOptions.RulesetID,
		"ruleset_version": *deleteZoneRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteZoneRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = rulesets.Service.Request(request, nil)

	return
}

// GetZoneEntrypointRuleset : Get a zone entrypoint ruleset
// Get the zone ruleset for the given phase's entrypoint.
func (rulesets *RulesetsV1) GetZoneEntrypointRuleset(getZoneEntrypointRulesetOptions *GetZoneEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneEntrypointRulesetWithContext(context.Background(), getZoneEntrypointRulesetOptions)
}

// GetZoneEntrypointRulesetWithContext is an alternate form of the GetZoneEntrypointRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneEntrypointRulesetWithContext(ctx context.Context, getZoneEntrypointRulesetOptions *GetZoneEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneEntrypointRulesetOptions, "getZoneEntrypointRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneEntrypointRulesetOptions, "getZoneEntrypointRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_phase":   *getZoneEntrypointRulesetOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/phases/{ruleset_phase}/entrypoint`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneEntrypointRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneEntrypointRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateZoneEntrypointRuleset : Update a zone entrypoint ruleset
// Updates the instance ruleset for the given phase's entry point.
func (rulesets *RulesetsV1) UpdateZoneEntrypointRuleset(updateZoneEntrypointRulesetOptions *UpdateZoneEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateZoneEntrypointRulesetWithContext(context.Background(), updateZoneEntrypointRulesetOptions)
}

// UpdateZoneEntrypointRulesetWithContext is an alternate form of the UpdateZoneEntrypointRuleset method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateZoneEntrypointRulesetWithContext(ctx context.Context, updateZoneEntrypointRulesetOptions *UpdateZoneEntrypointRulesetOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateZoneEntrypointRulesetOptions, "updateZoneEntrypointRulesetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateZoneEntrypointRulesetOptions, "updateZoneEntrypointRulesetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_phase":   *updateZoneEntrypointRulesetOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/phases/{ruleset_phase}/entrypoint`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateZoneEntrypointRulesetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateZoneEntrypointRuleset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateZoneEntrypointRulesetOptions.Description != nil {
		body["description"] = updateZoneEntrypointRulesetOptions.Description
	}
	if updateZoneEntrypointRulesetOptions.Kind != nil {
		body["kind"] = updateZoneEntrypointRulesetOptions.Kind
	}
	if updateZoneEntrypointRulesetOptions.Name != nil {
		body["name"] = updateZoneEntrypointRulesetOptions.Name
	}
	if updateZoneEntrypointRulesetOptions.Phase != nil {
		body["phase"] = updateZoneEntrypointRulesetOptions.Phase
	}
	if updateZoneEntrypointRulesetOptions.Rules != nil {
		body["rules"] = updateZoneEntrypointRulesetOptions.Rules
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneEntryPointRulesetVersions : List a zone entry point ruleset's versions
// Lists the zone ruleset versions for the given phase's entry point.
func (rulesets *RulesetsV1) GetZoneEntryPointRulesetVersions(getZoneEntryPointRulesetVersionsOptions *GetZoneEntryPointRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneEntryPointRulesetVersionsWithContext(context.Background(), getZoneEntryPointRulesetVersionsOptions)
}

// GetZoneEntryPointRulesetVersionsWithContext is an alternate form of the GetZoneEntryPointRulesetVersions method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneEntryPointRulesetVersionsWithContext(ctx context.Context, getZoneEntryPointRulesetVersionsOptions *GetZoneEntryPointRulesetVersionsOptions) (result *ListRulesetsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneEntryPointRulesetVersionsOptions, "getZoneEntryPointRulesetVersionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneEntryPointRulesetVersionsOptions, "getZoneEntryPointRulesetVersionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_phase":   *getZoneEntryPointRulesetVersionsOptions.RulesetPhase,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/phases/{ruleset_phase}/entrypoint/versions`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneEntryPointRulesetVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneEntryPointRulesetVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRulesetsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneEntryPointRulesetVersion : Get a zone entry point ruleset version
// Fetches a specific version of a zone entry point ruleset.
func (rulesets *RulesetsV1) GetZoneEntryPointRulesetVersion(getZoneEntryPointRulesetVersionOptions *GetZoneEntryPointRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.GetZoneEntryPointRulesetVersionWithContext(context.Background(), getZoneEntryPointRulesetVersionOptions)
}

// GetZoneEntryPointRulesetVersionWithContext is an alternate form of the GetZoneEntryPointRulesetVersion method which supports a Context parameter
func (rulesets *RulesetsV1) GetZoneEntryPointRulesetVersionWithContext(ctx context.Context, getZoneEntryPointRulesetVersionOptions *GetZoneEntryPointRulesetVersionOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneEntryPointRulesetVersionOptions, "getZoneEntryPointRulesetVersionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneEntryPointRulesetVersionOptions, "getZoneEntryPointRulesetVersionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_phase":   *getZoneEntryPointRulesetVersionOptions.RulesetPhase,
		"ruleset_version": *getZoneEntryPointRulesetVersionOptions.RulesetVersion,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/phases/{ruleset_phase}/entrypoint/versions/{ruleset_version}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneEntryPointRulesetVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "GetZoneEntryPointRulesetVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateZoneRulesetRule : Create a zone ruleset rule
// Create a zone ruleset rule.
func (rulesets *RulesetsV1) CreateZoneRulesetRule(createZoneRulesetRuleOptions *CreateZoneRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.CreateZoneRulesetRuleWithContext(context.Background(), createZoneRulesetRuleOptions)
}

// CreateZoneRulesetRuleWithContext is an alternate form of the CreateZoneRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) CreateZoneRulesetRuleWithContext(ctx context.Context, createZoneRulesetRuleOptions *CreateZoneRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createZoneRulesetRuleOptions, "createZoneRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createZoneRulesetRuleOptions, "createZoneRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *createZoneRulesetRuleOptions.RulesetID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createZoneRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "CreateZoneRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createZoneRulesetRuleOptions.Action != nil {
		body["action"] = createZoneRulesetRuleOptions.Action
	}
	if createZoneRulesetRuleOptions.ActionParameters != nil {
		body["action_parameters"] = createZoneRulesetRuleOptions.ActionParameters
	}
	if createZoneRulesetRuleOptions.Description != nil {
		body["description"] = createZoneRulesetRuleOptions.Description
	}
	if createZoneRulesetRuleOptions.Enabled != nil {
		body["enabled"] = createZoneRulesetRuleOptions.Enabled
	}
	if createZoneRulesetRuleOptions.Expression != nil {
		body["expression"] = createZoneRulesetRuleOptions.Expression
	}
	if createZoneRulesetRuleOptions.ID != nil {
		body["id"] = createZoneRulesetRuleOptions.ID
	}
	if createZoneRulesetRuleOptions.Logging != nil {
		body["logging"] = createZoneRulesetRuleOptions.Logging
	}
	if createZoneRulesetRuleOptions.Ref != nil {
		body["ref"] = createZoneRulesetRuleOptions.Ref
	}
	if createZoneRulesetRuleOptions.Position != nil {
		body["position"] = createZoneRulesetRuleOptions.Position
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateZoneRulesetRule : Update a zone ruleset rule
// Update a zone ruleset rule.
func (rulesets *RulesetsV1) UpdateZoneRulesetRule(updateZoneRulesetRuleOptions *UpdateZoneRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	return rulesets.UpdateZoneRulesetRuleWithContext(context.Background(), updateZoneRulesetRuleOptions)
}

// UpdateZoneRulesetRuleWithContext is an alternate form of the UpdateZoneRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) UpdateZoneRulesetRuleWithContext(ctx context.Context, updateZoneRulesetRuleOptions *UpdateZoneRulesetRuleOptions) (result *RulesetResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateZoneRulesetRuleOptions, "updateZoneRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateZoneRulesetRuleOptions, "updateZoneRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *updateZoneRulesetRuleOptions.RulesetID,
		"rule_id":         *updateZoneRulesetRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateZoneRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "UpdateZoneRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateZoneRulesetRuleOptions.Action != nil {
		body["action"] = updateZoneRulesetRuleOptions.Action
	}
	if updateZoneRulesetRuleOptions.ActionParameters != nil {
		body["action_parameters"] = updateZoneRulesetRuleOptions.ActionParameters
	}
	if updateZoneRulesetRuleOptions.Description != nil {
		body["description"] = updateZoneRulesetRuleOptions.Description
	}
	if updateZoneRulesetRuleOptions.Enabled != nil {
		body["enabled"] = updateZoneRulesetRuleOptions.Enabled
	}
	if updateZoneRulesetRuleOptions.Expression != nil {
		body["expression"] = updateZoneRulesetRuleOptions.Expression
	}
	if updateZoneRulesetRuleOptions.ID != nil {
		body["id"] = updateZoneRulesetRuleOptions.ID
	}
	if updateZoneRulesetRuleOptions.Logging != nil {
		body["logging"] = updateZoneRulesetRuleOptions.Logging
	}
	if updateZoneRulesetRuleOptions.Ref != nil {
		body["ref"] = updateZoneRulesetRuleOptions.Ref
	}
	if updateZoneRulesetRuleOptions.Position != nil {
		body["position"] = updateZoneRulesetRuleOptions.Position
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
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRulesetResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteZoneRulesetRule : Delete a zone ruleset rule
// Delete an instance ruleset rule.
func (rulesets *RulesetsV1) DeleteZoneRulesetRule(deleteZoneRulesetRuleOptions *DeleteZoneRulesetRuleOptions) (result *RuleResp, response *core.DetailedResponse, err error) {
	return rulesets.DeleteZoneRulesetRuleWithContext(context.Background(), deleteZoneRulesetRuleOptions)
}

// DeleteZoneRulesetRuleWithContext is an alternate form of the DeleteZoneRulesetRule method which supports a Context parameter
func (rulesets *RulesetsV1) DeleteZoneRulesetRuleWithContext(ctx context.Context, deleteZoneRulesetRuleOptions *DeleteZoneRulesetRuleOptions) (result *RuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneRulesetRuleOptions, "deleteZoneRulesetRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneRulesetRuleOptions, "deleteZoneRulesetRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *rulesets.Crn,
		"zone_identifier": *rulesets.ZoneIdentifier,
		"ruleset_id":      *deleteZoneRulesetRuleOptions.RulesetID,
		"rule_id":         *deleteZoneRulesetRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rulesets.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rulesets.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rulesets/{ruleset_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneRulesetRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("rulesets", "V1", "DeleteZoneRulesetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rulesets.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ActionParametersResponse : ActionParametersResponse struct
type ActionParametersResponse struct {
	// the content to return.
	Content *string `json:"content" validate:"required"`

	ContentType *string `json:"content_type" validate:"required"`

	// The status code to return.
	StatusCode *int64 `json:"status_code" validate:"required"`
}

// NewActionParametersResponse : Instantiate ActionParametersResponse (Generic Model Constructor)
func (*RulesetsV1) NewActionParametersResponse(content string, contentType string, statusCode int64) (_model *ActionParametersResponse, err error) {
	_model = &ActionParametersResponse{
		Content:     core.StringPtr(content),
		ContentType: core.StringPtr(contentType),
		StatusCode:  core.Int64Ptr(statusCode),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalActionParametersResponse unmarshals an instance of ActionParametersResponse from the specified map of raw messages.
func UnmarshalActionParametersResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionParametersResponse)
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateInstanceRulesetRuleOptions : The CreateInstanceRulesetRule options.
type CreateInstanceRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	Description *string `json:"description,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression,omitempty"`

	ID *string `json:"id,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Position *Position `json:"position,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateInstanceRulesetRuleOptions : Instantiate CreateInstanceRulesetRuleOptions
func (*RulesetsV1) NewCreateInstanceRulesetRuleOptions(rulesetID string) *CreateInstanceRulesetRuleOptions {
	return &CreateInstanceRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *CreateInstanceRulesetRuleOptions) SetRulesetID(rulesetID string) *CreateInstanceRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *CreateInstanceRulesetRuleOptions) SetAction(action string) *CreateInstanceRulesetRuleOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetActionParameters : Allow user to set ActionParameters
func (_options *CreateInstanceRulesetRuleOptions) SetActionParameters(actionParameters *ActionParameters) *CreateInstanceRulesetRuleOptions {
	_options.ActionParameters = actionParameters
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateInstanceRulesetRuleOptions) SetDescription(description string) *CreateInstanceRulesetRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateInstanceRulesetRuleOptions) SetEnabled(enabled bool) *CreateInstanceRulesetRuleOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetExpression : Allow user to set Expression
func (_options *CreateInstanceRulesetRuleOptions) SetExpression(expression string) *CreateInstanceRulesetRuleOptions {
	_options.Expression = core.StringPtr(expression)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateInstanceRulesetRuleOptions) SetID(id string) *CreateInstanceRulesetRuleOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLogging : Allow user to set Logging
func (_options *CreateInstanceRulesetRuleOptions) SetLogging(logging *Logging) *CreateInstanceRulesetRuleOptions {
	_options.Logging = logging
	return _options
}

// SetRef : Allow user to set Ref
func (_options *CreateInstanceRulesetRuleOptions) SetRef(ref string) *CreateInstanceRulesetRuleOptions {
	_options.Ref = core.StringPtr(ref)
	return _options
}

// SetPosition : Allow user to set Position
func (_options *CreateInstanceRulesetRuleOptions) SetPosition(position *Position) *CreateInstanceRulesetRuleOptions {
	_options.Position = position
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateInstanceRulesetRuleOptions) SetHeaders(param map[string]string) *CreateInstanceRulesetRuleOptions {
	options.Headers = param
	return options
}

// CreateZoneRulesetRuleOptions : The CreateZoneRulesetRule options.
type CreateZoneRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	Description *string `json:"description,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression,omitempty"`

	ID *string `json:"id,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Position *Position `json:"position,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateZoneRulesetRuleOptions : Instantiate CreateZoneRulesetRuleOptions
func (*RulesetsV1) NewCreateZoneRulesetRuleOptions(rulesetID string) *CreateZoneRulesetRuleOptions {
	return &CreateZoneRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *CreateZoneRulesetRuleOptions) SetRulesetID(rulesetID string) *CreateZoneRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *CreateZoneRulesetRuleOptions) SetAction(action string) *CreateZoneRulesetRuleOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetActionParameters : Allow user to set ActionParameters
func (_options *CreateZoneRulesetRuleOptions) SetActionParameters(actionParameters *ActionParameters) *CreateZoneRulesetRuleOptions {
	_options.ActionParameters = actionParameters
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateZoneRulesetRuleOptions) SetDescription(description string) *CreateZoneRulesetRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateZoneRulesetRuleOptions) SetEnabled(enabled bool) *CreateZoneRulesetRuleOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetExpression : Allow user to set Expression
func (_options *CreateZoneRulesetRuleOptions) SetExpression(expression string) *CreateZoneRulesetRuleOptions {
	_options.Expression = core.StringPtr(expression)
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateZoneRulesetRuleOptions) SetID(id string) *CreateZoneRulesetRuleOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLogging : Allow user to set Logging
func (_options *CreateZoneRulesetRuleOptions) SetLogging(logging *Logging) *CreateZoneRulesetRuleOptions {
	_options.Logging = logging
	return _options
}

// SetRef : Allow user to set Ref
func (_options *CreateZoneRulesetRuleOptions) SetRef(ref string) *CreateZoneRulesetRuleOptions {
	_options.Ref = core.StringPtr(ref)
	return _options
}

// SetPosition : Allow user to set Position
func (_options *CreateZoneRulesetRuleOptions) SetPosition(position *Position) *CreateZoneRulesetRuleOptions {
	_options.Position = position
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneRulesetRuleOptions) SetHeaders(param map[string]string) *CreateZoneRulesetRuleOptions {
	options.Headers = param
	return options
}

// DeleteInstanceRulesetOptions : The DeleteInstanceRuleset options.
type DeleteInstanceRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteInstanceRulesetOptions : Instantiate DeleteInstanceRulesetOptions
func (*RulesetsV1) NewDeleteInstanceRulesetOptions(rulesetID string) *DeleteInstanceRulesetOptions {
	return &DeleteInstanceRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteInstanceRulesetOptions) SetRulesetID(rulesetID string) *DeleteInstanceRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteInstanceRulesetOptions) SetHeaders(param map[string]string) *DeleteInstanceRulesetOptions {
	options.Headers = param
	return options
}

// DeleteInstanceRulesetRuleOptions : The DeleteInstanceRulesetRule options.
type DeleteInstanceRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// ID of a specific rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteInstanceRulesetRuleOptions : Instantiate DeleteInstanceRulesetRuleOptions
func (*RulesetsV1) NewDeleteInstanceRulesetRuleOptions(rulesetID string, ruleID string) *DeleteInstanceRulesetRuleOptions {
	return &DeleteInstanceRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteInstanceRulesetRuleOptions) SetRulesetID(rulesetID string) *DeleteInstanceRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteInstanceRulesetRuleOptions) SetRuleID(ruleID string) *DeleteInstanceRulesetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteInstanceRulesetRuleOptions) SetHeaders(param map[string]string) *DeleteInstanceRulesetRuleOptions {
	options.Headers = param
	return options
}

// DeleteInstanceRulesetVersionOptions : The DeleteInstanceRulesetVersion options.
type DeleteInstanceRulesetVersionOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteInstanceRulesetVersionOptions : Instantiate DeleteInstanceRulesetVersionOptions
func (*RulesetsV1) NewDeleteInstanceRulesetVersionOptions(rulesetID string, rulesetVersion string) *DeleteInstanceRulesetVersionOptions {
	return &DeleteInstanceRulesetVersionOptions{
		RulesetID:      core.StringPtr(rulesetID),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteInstanceRulesetVersionOptions) SetRulesetID(rulesetID string) *DeleteInstanceRulesetVersionOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *DeleteInstanceRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *DeleteInstanceRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteInstanceRulesetVersionOptions) SetHeaders(param map[string]string) *DeleteInstanceRulesetVersionOptions {
	options.Headers = param
	return options
}

// DeleteZoneRulesetOptions : The DeleteZoneRuleset options.
type DeleteZoneRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneRulesetOptions : Instantiate DeleteZoneRulesetOptions
func (*RulesetsV1) NewDeleteZoneRulesetOptions(rulesetID string) *DeleteZoneRulesetOptions {
	return &DeleteZoneRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteZoneRulesetOptions) SetRulesetID(rulesetID string) *DeleteZoneRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneRulesetOptions) SetHeaders(param map[string]string) *DeleteZoneRulesetOptions {
	options.Headers = param
	return options
}

// DeleteZoneRulesetRuleOptions : The DeleteZoneRulesetRule options.
type DeleteZoneRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// ID of a specific rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneRulesetRuleOptions : Instantiate DeleteZoneRulesetRuleOptions
func (*RulesetsV1) NewDeleteZoneRulesetRuleOptions(rulesetID string, ruleID string) *DeleteZoneRulesetRuleOptions {
	return &DeleteZoneRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteZoneRulesetRuleOptions) SetRulesetID(rulesetID string) *DeleteZoneRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteZoneRulesetRuleOptions) SetRuleID(ruleID string) *DeleteZoneRulesetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneRulesetRuleOptions) SetHeaders(param map[string]string) *DeleteZoneRulesetRuleOptions {
	options.Headers = param
	return options
}

// DeleteZoneRulesetVersionOptions : The DeleteZoneRulesetVersion options.
type DeleteZoneRulesetVersionOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneRulesetVersionOptions : Instantiate DeleteZoneRulesetVersionOptions
func (*RulesetsV1) NewDeleteZoneRulesetVersionOptions(rulesetID string, rulesetVersion string) *DeleteZoneRulesetVersionOptions {
	return &DeleteZoneRulesetVersionOptions{
		RulesetID:      core.StringPtr(rulesetID),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *DeleteZoneRulesetVersionOptions) SetRulesetID(rulesetID string) *DeleteZoneRulesetVersionOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *DeleteZoneRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *DeleteZoneRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneRulesetVersionOptions) SetHeaders(param map[string]string) *DeleteZoneRulesetVersionOptions {
	options.Headers = param
	return options
}

// GetInstanceEntryPointRulesetVersionOptions : The GetInstanceEntryPointRulesetVersion options.
type GetInstanceEntryPointRulesetVersionOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetInstanceEntryPointRulesetVersionOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetInstanceEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetInstanceEntryPointRulesetVersionOptions : Instantiate GetInstanceEntryPointRulesetVersionOptions
func (*RulesetsV1) NewGetInstanceEntryPointRulesetVersionOptions(rulesetPhase string, rulesetVersion string) *GetInstanceEntryPointRulesetVersionOptions {
	return &GetInstanceEntryPointRulesetVersionOptions{
		RulesetPhase:   core.StringPtr(rulesetPhase),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetInstanceEntryPointRulesetVersionOptions) SetRulesetPhase(rulesetPhase string) *GetInstanceEntryPointRulesetVersionOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *GetInstanceEntryPointRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *GetInstanceEntryPointRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceEntryPointRulesetVersionOptions) SetHeaders(param map[string]string) *GetInstanceEntryPointRulesetVersionOptions {
	options.Headers = param
	return options
}

// GetInstanceEntryPointRulesetVersionsOptions : The GetInstanceEntryPointRulesetVersions options.
type GetInstanceEntryPointRulesetVersionsOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetInstanceEntryPointRulesetVersionsOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetInstanceEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetInstanceEntryPointRulesetVersionsOptions : Instantiate GetInstanceEntryPointRulesetVersionsOptions
func (*RulesetsV1) NewGetInstanceEntryPointRulesetVersionsOptions(rulesetPhase string) *GetInstanceEntryPointRulesetVersionsOptions {
	return &GetInstanceEntryPointRulesetVersionsOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetInstanceEntryPointRulesetVersionsOptions) SetRulesetPhase(rulesetPhase string) *GetInstanceEntryPointRulesetVersionsOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceEntryPointRulesetVersionsOptions) SetHeaders(param map[string]string) *GetInstanceEntryPointRulesetVersionsOptions {
	options.Headers = param
	return options
}

// GetInstanceEntrypointRulesetOptions : The GetInstanceEntrypointRuleset options.
type GetInstanceEntrypointRulesetOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetInstanceEntrypointRulesetOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetInstanceEntrypointRulesetOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetInstanceEntrypointRulesetOptions : Instantiate GetInstanceEntrypointRulesetOptions
func (*RulesetsV1) NewGetInstanceEntrypointRulesetOptions(rulesetPhase string) *GetInstanceEntrypointRulesetOptions {
	return &GetInstanceEntrypointRulesetOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetInstanceEntrypointRulesetOptions) SetRulesetPhase(rulesetPhase string) *GetInstanceEntrypointRulesetOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceEntrypointRulesetOptions) SetHeaders(param map[string]string) *GetInstanceEntrypointRulesetOptions {
	options.Headers = param
	return options
}

// GetInstanceRulesetOptions : The GetInstanceRuleset options.
type GetInstanceRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInstanceRulesetOptions : Instantiate GetInstanceRulesetOptions
func (*RulesetsV1) NewGetInstanceRulesetOptions(rulesetID string) *GetInstanceRulesetOptions {
	return &GetInstanceRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetInstanceRulesetOptions) SetRulesetID(rulesetID string) *GetInstanceRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceRulesetOptions) SetHeaders(param map[string]string) *GetInstanceRulesetOptions {
	options.Headers = param
	return options
}

// GetInstanceRulesetVersionByTagOptions : The GetInstanceRulesetVersionByTag options.
type GetInstanceRulesetVersionByTagOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// A category of the rule.
	RuleTag *string `json:"rule_tag" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInstanceRulesetVersionByTagOptions : Instantiate GetInstanceRulesetVersionByTagOptions
func (*RulesetsV1) NewGetInstanceRulesetVersionByTagOptions(rulesetID string, rulesetVersion string, ruleTag string) *GetInstanceRulesetVersionByTagOptions {
	return &GetInstanceRulesetVersionByTagOptions{
		RulesetID:      core.StringPtr(rulesetID),
		RulesetVersion: core.StringPtr(rulesetVersion),
		RuleTag:        core.StringPtr(ruleTag),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetInstanceRulesetVersionByTagOptions) SetRulesetID(rulesetID string) *GetInstanceRulesetVersionByTagOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *GetInstanceRulesetVersionByTagOptions) SetRulesetVersion(rulesetVersion string) *GetInstanceRulesetVersionByTagOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetRuleTag : Allow user to set RuleTag
func (_options *GetInstanceRulesetVersionByTagOptions) SetRuleTag(ruleTag string) *GetInstanceRulesetVersionByTagOptions {
	_options.RuleTag = core.StringPtr(ruleTag)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceRulesetVersionByTagOptions) SetHeaders(param map[string]string) *GetInstanceRulesetVersionByTagOptions {
	options.Headers = param
	return options
}

// GetInstanceRulesetVersionOptions : The GetInstanceRulesetVersion options.
type GetInstanceRulesetVersionOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInstanceRulesetVersionOptions : Instantiate GetInstanceRulesetVersionOptions
func (*RulesetsV1) NewGetInstanceRulesetVersionOptions(rulesetID string, rulesetVersion string) *GetInstanceRulesetVersionOptions {
	return &GetInstanceRulesetVersionOptions{
		RulesetID:      core.StringPtr(rulesetID),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetInstanceRulesetVersionOptions) SetRulesetID(rulesetID string) *GetInstanceRulesetVersionOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *GetInstanceRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *GetInstanceRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceRulesetVersionOptions) SetHeaders(param map[string]string) *GetInstanceRulesetVersionOptions {
	options.Headers = param
	return options
}

// GetInstanceRulesetVersionsOptions : The GetInstanceRulesetVersions options.
type GetInstanceRulesetVersionsOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInstanceRulesetVersionsOptions : Instantiate GetInstanceRulesetVersionsOptions
func (*RulesetsV1) NewGetInstanceRulesetVersionsOptions(rulesetID string) *GetInstanceRulesetVersionsOptions {
	return &GetInstanceRulesetVersionsOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetInstanceRulesetVersionsOptions) SetRulesetID(rulesetID string) *GetInstanceRulesetVersionsOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceRulesetVersionsOptions) SetHeaders(param map[string]string) *GetInstanceRulesetVersionsOptions {
	options.Headers = param
	return options
}

// GetInstanceRulesetsOptions : The GetInstanceRulesets options.
type GetInstanceRulesetsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetInstanceRulesetsOptions : Instantiate GetInstanceRulesetsOptions
func (*RulesetsV1) NewGetInstanceRulesetsOptions() *GetInstanceRulesetsOptions {
	return &GetInstanceRulesetsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceRulesetsOptions) SetHeaders(param map[string]string) *GetInstanceRulesetsOptions {
	options.Headers = param
	return options
}

// GetZoneEntryPointRulesetVersionOptions : The GetZoneEntryPointRulesetVersion options.
type GetZoneEntryPointRulesetVersionOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetZoneEntryPointRulesetVersionOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetZoneEntryPointRulesetVersionOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetZoneEntryPointRulesetVersionOptions : Instantiate GetZoneEntryPointRulesetVersionOptions
func (*RulesetsV1) NewGetZoneEntryPointRulesetVersionOptions(rulesetPhase string, rulesetVersion string) *GetZoneEntryPointRulesetVersionOptions {
	return &GetZoneEntryPointRulesetVersionOptions{
		RulesetPhase:   core.StringPtr(rulesetPhase),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetZoneEntryPointRulesetVersionOptions) SetRulesetPhase(rulesetPhase string) *GetZoneEntryPointRulesetVersionOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *GetZoneEntryPointRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *GetZoneEntryPointRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneEntryPointRulesetVersionOptions) SetHeaders(param map[string]string) *GetZoneEntryPointRulesetVersionOptions {
	options.Headers = param
	return options
}

// GetZoneEntryPointRulesetVersionsOptions : The GetZoneEntryPointRulesetVersions options.
type GetZoneEntryPointRulesetVersionsOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetZoneEntryPointRulesetVersionsOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetZoneEntryPointRulesetVersionsOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetZoneEntryPointRulesetVersionsOptions : Instantiate GetZoneEntryPointRulesetVersionsOptions
func (*RulesetsV1) NewGetZoneEntryPointRulesetVersionsOptions(rulesetPhase string) *GetZoneEntryPointRulesetVersionsOptions {
	return &GetZoneEntryPointRulesetVersionsOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetZoneEntryPointRulesetVersionsOptions) SetRulesetPhase(rulesetPhase string) *GetZoneEntryPointRulesetVersionsOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneEntryPointRulesetVersionsOptions) SetHeaders(param map[string]string) *GetZoneEntryPointRulesetVersionsOptions {
	options.Headers = param
	return options
}

// GetZoneEntrypointRulesetOptions : The GetZoneEntrypointRuleset options.
type GetZoneEntrypointRulesetOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetZoneEntrypointRulesetOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	GetZoneEntrypointRulesetOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	GetZoneEntrypointRulesetOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	GetZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewGetZoneEntrypointRulesetOptions : Instantiate GetZoneEntrypointRulesetOptions
func (*RulesetsV1) NewGetZoneEntrypointRulesetOptions(rulesetPhase string) *GetZoneEntrypointRulesetOptions {
	return &GetZoneEntrypointRulesetOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *GetZoneEntrypointRulesetOptions) SetRulesetPhase(rulesetPhase string) *GetZoneEntrypointRulesetOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneEntrypointRulesetOptions) SetHeaders(param map[string]string) *GetZoneEntrypointRulesetOptions {
	options.Headers = param
	return options
}

// GetZoneRulesetOptions : The GetZoneRuleset options.
type GetZoneRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneRulesetOptions : Instantiate GetZoneRulesetOptions
func (*RulesetsV1) NewGetZoneRulesetOptions(rulesetID string) *GetZoneRulesetOptions {
	return &GetZoneRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetZoneRulesetOptions) SetRulesetID(rulesetID string) *GetZoneRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneRulesetOptions) SetHeaders(param map[string]string) *GetZoneRulesetOptions {
	options.Headers = param
	return options
}

// GetZoneRulesetVersionOptions : The GetZoneRulesetVersion options.
type GetZoneRulesetVersionOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// The version of the ruleset.
	RulesetVersion *string `json:"ruleset_version" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneRulesetVersionOptions : Instantiate GetZoneRulesetVersionOptions
func (*RulesetsV1) NewGetZoneRulesetVersionOptions(rulesetID string, rulesetVersion string) *GetZoneRulesetVersionOptions {
	return &GetZoneRulesetVersionOptions{
		RulesetID:      core.StringPtr(rulesetID),
		RulesetVersion: core.StringPtr(rulesetVersion),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetZoneRulesetVersionOptions) SetRulesetID(rulesetID string) *GetZoneRulesetVersionOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRulesetVersion : Allow user to set RulesetVersion
func (_options *GetZoneRulesetVersionOptions) SetRulesetVersion(rulesetVersion string) *GetZoneRulesetVersionOptions {
	_options.RulesetVersion = core.StringPtr(rulesetVersion)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneRulesetVersionOptions) SetHeaders(param map[string]string) *GetZoneRulesetVersionOptions {
	options.Headers = param
	return options
}

// GetZoneRulesetVersionsOptions : The GetZoneRulesetVersions options.
type GetZoneRulesetVersionsOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneRulesetVersionsOptions : Instantiate GetZoneRulesetVersionsOptions
func (*RulesetsV1) NewGetZoneRulesetVersionsOptions(rulesetID string) *GetZoneRulesetVersionsOptions {
	return &GetZoneRulesetVersionsOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *GetZoneRulesetVersionsOptions) SetRulesetID(rulesetID string) *GetZoneRulesetVersionsOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneRulesetVersionsOptions) SetHeaders(param map[string]string) *GetZoneRulesetVersionsOptions {
	options.Headers = param
	return options
}

// GetZoneRulesetsOptions : The GetZoneRulesets options.
type GetZoneRulesetsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneRulesetsOptions : Instantiate GetZoneRulesetsOptions
func (*RulesetsV1) NewGetZoneRulesetsOptions() *GetZoneRulesetsOptions {
	return &GetZoneRulesetsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneRulesetsOptions) SetHeaders(param map[string]string) *GetZoneRulesetsOptions {
	options.Headers = param
	return options
}

// MessageSource : The source of this message.
type MessageSource struct {
	// A JSON pointer to the field that is the source of the message.
	Pointer *string `json:"pointer" validate:"required"`
}

// UnmarshalMessageSource unmarshals an instance of MessageSource from the specified map of raw messages.
func UnmarshalMessageSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MessageSource)
	err = core.UnmarshalPrimitive(m, "pointer", &obj.Pointer)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateInstanceEntrypointRulesetOptions : The UpdateInstanceEntrypointRuleset options.
type UpdateInstanceEntrypointRulesetOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// description of the ruleset.
	Description *string `json:"description,omitempty"`

	Kind *string `json:"kind,omitempty"`

	// human readable name of the ruleset.
	Name *string `json:"name,omitempty"`

	// The phase of the ruleset.
	Phase *string `json:"phase,omitempty"`

	Rules []RuleCreate `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateInstanceEntrypointRulesetOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateInstanceEntrypointRulesetOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// Constants associated with the UpdateInstanceEntrypointRulesetOptions.Kind property.
const (
	UpdateInstanceEntrypointRulesetOptions_Kind_Custom  = "custom"
	UpdateInstanceEntrypointRulesetOptions_Kind_Managed = "managed"
	UpdateInstanceEntrypointRulesetOptions_Kind_Root    = "root"
	UpdateInstanceEntrypointRulesetOptions_Kind_Zone    = "zone"
)

// Constants associated with the UpdateInstanceEntrypointRulesetOptions.Phase property.
// The phase of the ruleset.
const (
	UpdateInstanceEntrypointRulesetOptions_Phase_DdosL4                         = "ddos_l4"
	UpdateInstanceEntrypointRulesetOptions_Phase_DdosL7                         = "ddos_l7"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpConfigSettings             = "http_config_settings"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpCustomErrors               = "http_custom_errors"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRatelimit                  = "http_ratelimit"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestOrigin              = "http_request_origin"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestRedirect            = "http_request_redirect"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpRequestTransform           = "http_request_transform"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpResponseCompression        = "http_response_compression"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateInstanceEntrypointRulesetOptions_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewUpdateInstanceEntrypointRulesetOptions : Instantiate UpdateInstanceEntrypointRulesetOptions
func (*RulesetsV1) NewUpdateInstanceEntrypointRulesetOptions(rulesetPhase string) *UpdateInstanceEntrypointRulesetOptions {
	return &UpdateInstanceEntrypointRulesetOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *UpdateInstanceEntrypointRulesetOptions) SetRulesetPhase(rulesetPhase string) *UpdateInstanceEntrypointRulesetOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateInstanceEntrypointRulesetOptions) SetDescription(description string) *UpdateInstanceEntrypointRulesetOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateInstanceEntrypointRulesetOptions) SetKind(kind string) *UpdateInstanceEntrypointRulesetOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateInstanceEntrypointRulesetOptions) SetName(name string) *UpdateInstanceEntrypointRulesetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPhase : Allow user to set Phase
func (_options *UpdateInstanceEntrypointRulesetOptions) SetPhase(phase string) *UpdateInstanceEntrypointRulesetOptions {
	_options.Phase = core.StringPtr(phase)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *UpdateInstanceEntrypointRulesetOptions) SetRules(rules []RuleCreate) *UpdateInstanceEntrypointRulesetOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateInstanceEntrypointRulesetOptions) SetHeaders(param map[string]string) *UpdateInstanceEntrypointRulesetOptions {
	options.Headers = param
	return options
}

// UpdateInstanceRulesetOptions : The UpdateInstanceRuleset options.
type UpdateInstanceRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// description of the ruleset.
	Description *string `json:"description,omitempty"`

	Kind *string `json:"kind,omitempty"`

	// human readable name of the ruleset.
	Name *string `json:"name,omitempty"`

	// The phase of the ruleset.
	Phase *string `json:"phase,omitempty"`

	Rules []RuleCreate `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateInstanceRulesetOptions.Kind property.
const (
	UpdateInstanceRulesetOptions_Kind_Custom  = "custom"
	UpdateInstanceRulesetOptions_Kind_Managed = "managed"
	UpdateInstanceRulesetOptions_Kind_Root    = "root"
	UpdateInstanceRulesetOptions_Kind_Zone    = "zone"
)

// Constants associated with the UpdateInstanceRulesetOptions.Phase property.
// The phase of the ruleset.
const (
	UpdateInstanceRulesetOptions_Phase_DdosL4                         = "ddos_l4"
	UpdateInstanceRulesetOptions_Phase_DdosL7                         = "ddos_l7"
	UpdateInstanceRulesetOptions_Phase_HttpConfigSettings             = "http_config_settings"
	UpdateInstanceRulesetOptions_Phase_HttpCustomErrors               = "http_custom_errors"
	UpdateInstanceRulesetOptions_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateInstanceRulesetOptions_Phase_HttpRatelimit                  = "http_ratelimit"
	UpdateInstanceRulesetOptions_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateInstanceRulesetOptions_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateInstanceRulesetOptions_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateInstanceRulesetOptions_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateInstanceRulesetOptions_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateInstanceRulesetOptions_Phase_HttpRequestOrigin              = "http_request_origin"
	UpdateInstanceRulesetOptions_Phase_HttpRequestRedirect            = "http_request_redirect"
	UpdateInstanceRulesetOptions_Phase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateInstanceRulesetOptions_Phase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateInstanceRulesetOptions_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateInstanceRulesetOptions_Phase_HttpRequestTransform           = "http_request_transform"
	UpdateInstanceRulesetOptions_Phase_HttpResponseCompression        = "http_response_compression"
	UpdateInstanceRulesetOptions_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateInstanceRulesetOptions_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewUpdateInstanceRulesetOptions : Instantiate UpdateInstanceRulesetOptions
func (*RulesetsV1) NewUpdateInstanceRulesetOptions(rulesetID string) *UpdateInstanceRulesetOptions {
	return &UpdateInstanceRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *UpdateInstanceRulesetOptions) SetRulesetID(rulesetID string) *UpdateInstanceRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateInstanceRulesetOptions) SetDescription(description string) *UpdateInstanceRulesetOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateInstanceRulesetOptions) SetKind(kind string) *UpdateInstanceRulesetOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateInstanceRulesetOptions) SetName(name string) *UpdateInstanceRulesetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPhase : Allow user to set Phase
func (_options *UpdateInstanceRulesetOptions) SetPhase(phase string) *UpdateInstanceRulesetOptions {
	_options.Phase = core.StringPtr(phase)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *UpdateInstanceRulesetOptions) SetRules(rules []RuleCreate) *UpdateInstanceRulesetOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateInstanceRulesetOptions) SetHeaders(param map[string]string) *UpdateInstanceRulesetOptions {
	options.Headers = param
	return options
}

// UpdateInstanceRulesetRuleOptions : The UpdateInstanceRulesetRule options.
type UpdateInstanceRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// ID of a specific rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	Description *string `json:"description,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression,omitempty"`

	ID *string `json:"id,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Position *Position `json:"position,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateInstanceRulesetRuleOptions : Instantiate UpdateInstanceRulesetRuleOptions
func (*RulesetsV1) NewUpdateInstanceRulesetRuleOptions(rulesetID string, ruleID string) *UpdateInstanceRulesetRuleOptions {
	return &UpdateInstanceRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *UpdateInstanceRulesetRuleOptions) SetRulesetID(rulesetID string) *UpdateInstanceRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateInstanceRulesetRuleOptions) SetRuleID(ruleID string) *UpdateInstanceRulesetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *UpdateInstanceRulesetRuleOptions) SetAction(action string) *UpdateInstanceRulesetRuleOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetActionParameters : Allow user to set ActionParameters
func (_options *UpdateInstanceRulesetRuleOptions) SetActionParameters(actionParameters *ActionParameters) *UpdateInstanceRulesetRuleOptions {
	_options.ActionParameters = actionParameters
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateInstanceRulesetRuleOptions) SetDescription(description string) *UpdateInstanceRulesetRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateInstanceRulesetRuleOptions) SetEnabled(enabled bool) *UpdateInstanceRulesetRuleOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetExpression : Allow user to set Expression
func (_options *UpdateInstanceRulesetRuleOptions) SetExpression(expression string) *UpdateInstanceRulesetRuleOptions {
	_options.Expression = core.StringPtr(expression)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateInstanceRulesetRuleOptions) SetID(id string) *UpdateInstanceRulesetRuleOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLogging : Allow user to set Logging
func (_options *UpdateInstanceRulesetRuleOptions) SetLogging(logging *Logging) *UpdateInstanceRulesetRuleOptions {
	_options.Logging = logging
	return _options
}

// SetRef : Allow user to set Ref
func (_options *UpdateInstanceRulesetRuleOptions) SetRef(ref string) *UpdateInstanceRulesetRuleOptions {
	_options.Ref = core.StringPtr(ref)
	return _options
}

// SetPosition : Allow user to set Position
func (_options *UpdateInstanceRulesetRuleOptions) SetPosition(position *Position) *UpdateInstanceRulesetRuleOptions {
	_options.Position = position
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateInstanceRulesetRuleOptions) SetHeaders(param map[string]string) *UpdateInstanceRulesetRuleOptions {
	options.Headers = param
	return options
}

// UpdateZoneEntrypointRulesetOptions : The UpdateZoneEntrypointRuleset options.
type UpdateZoneEntrypointRulesetOptions struct {
	// The phase of the ruleset.
	RulesetPhase *string `json:"ruleset_phase" validate:"required,ne="`

	// description of the ruleset.
	Description *string `json:"description,omitempty"`

	Kind *string `json:"kind,omitempty"`

	// human readable name of the ruleset.
	Name *string `json:"name,omitempty"`

	// The phase of the ruleset.
	Phase *string `json:"phase,omitempty"`

	Rules []RuleCreate `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateZoneEntrypointRulesetOptions.RulesetPhase property.
// The phase of the ruleset.
const (
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_DdosL4                         = "ddos_l4"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_DdosL7                         = "ddos_l7"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpConfigSettings             = "http_config_settings"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpCustomErrors               = "http_custom_errors"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRatelimit                  = "http_ratelimit"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestOrigin              = "http_request_origin"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestRedirect            = "http_request_redirect"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpRequestTransform           = "http_request_transform"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseCompression        = "http_response_compression"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateZoneEntrypointRulesetOptions_RulesetPhase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// Constants associated with the UpdateZoneEntrypointRulesetOptions.Kind property.
const (
	UpdateZoneEntrypointRulesetOptions_Kind_Custom  = "custom"
	UpdateZoneEntrypointRulesetOptions_Kind_Managed = "managed"
	UpdateZoneEntrypointRulesetOptions_Kind_Root    = "root"
	UpdateZoneEntrypointRulesetOptions_Kind_Zone    = "zone"
)

// Constants associated with the UpdateZoneEntrypointRulesetOptions.Phase property.
// The phase of the ruleset.
const (
	UpdateZoneEntrypointRulesetOptions_Phase_DdosL4                         = "ddos_l4"
	UpdateZoneEntrypointRulesetOptions_Phase_DdosL7                         = "ddos_l7"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpConfigSettings             = "http_config_settings"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpCustomErrors               = "http_custom_errors"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRatelimit                  = "http_ratelimit"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestOrigin              = "http_request_origin"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestRedirect            = "http_request_redirect"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpRequestTransform           = "http_request_transform"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpResponseCompression        = "http_response_compression"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateZoneEntrypointRulesetOptions_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewUpdateZoneEntrypointRulesetOptions : Instantiate UpdateZoneEntrypointRulesetOptions
func (*RulesetsV1) NewUpdateZoneEntrypointRulesetOptions(rulesetPhase string) *UpdateZoneEntrypointRulesetOptions {
	return &UpdateZoneEntrypointRulesetOptions{
		RulesetPhase: core.StringPtr(rulesetPhase),
	}
}

// SetRulesetPhase : Allow user to set RulesetPhase
func (_options *UpdateZoneEntrypointRulesetOptions) SetRulesetPhase(rulesetPhase string) *UpdateZoneEntrypointRulesetOptions {
	_options.RulesetPhase = core.StringPtr(rulesetPhase)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateZoneEntrypointRulesetOptions) SetDescription(description string) *UpdateZoneEntrypointRulesetOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateZoneEntrypointRulesetOptions) SetKind(kind string) *UpdateZoneEntrypointRulesetOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateZoneEntrypointRulesetOptions) SetName(name string) *UpdateZoneEntrypointRulesetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPhase : Allow user to set Phase
func (_options *UpdateZoneEntrypointRulesetOptions) SetPhase(phase string) *UpdateZoneEntrypointRulesetOptions {
	_options.Phase = core.StringPtr(phase)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *UpdateZoneEntrypointRulesetOptions) SetRules(rules []RuleCreate) *UpdateZoneEntrypointRulesetOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateZoneEntrypointRulesetOptions) SetHeaders(param map[string]string) *UpdateZoneEntrypointRulesetOptions {
	options.Headers = param
	return options
}

// UpdateZoneRulesetOptions : The UpdateZoneRuleset options.
type UpdateZoneRulesetOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// description of the ruleset.
	Description *string `json:"description,omitempty"`

	Kind *string `json:"kind,omitempty"`

	// human readable name of the ruleset.
	Name *string `json:"name,omitempty"`

	// The phase of the ruleset.
	Phase *string `json:"phase,omitempty"`

	Rules []RuleCreate `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateZoneRulesetOptions.Kind property.
const (
	UpdateZoneRulesetOptions_Kind_Custom  = "custom"
	UpdateZoneRulesetOptions_Kind_Managed = "managed"
	UpdateZoneRulesetOptions_Kind_Root    = "root"
	UpdateZoneRulesetOptions_Kind_Zone    = "zone"
)

// Constants associated with the UpdateZoneRulesetOptions.Phase property.
// The phase of the ruleset.
const (
	UpdateZoneRulesetOptions_Phase_DdosL4                         = "ddos_l4"
	UpdateZoneRulesetOptions_Phase_DdosL7                         = "ddos_l7"
	UpdateZoneRulesetOptions_Phase_HttpConfigSettings             = "http_config_settings"
	UpdateZoneRulesetOptions_Phase_HttpCustomErrors               = "http_custom_errors"
	UpdateZoneRulesetOptions_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	UpdateZoneRulesetOptions_Phase_HttpRatelimit                  = "http_ratelimit"
	UpdateZoneRulesetOptions_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	UpdateZoneRulesetOptions_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	UpdateZoneRulesetOptions_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	UpdateZoneRulesetOptions_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	UpdateZoneRulesetOptions_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	UpdateZoneRulesetOptions_Phase_HttpRequestOrigin              = "http_request_origin"
	UpdateZoneRulesetOptions_Phase_HttpRequestRedirect            = "http_request_redirect"
	UpdateZoneRulesetOptions_Phase_HttpRequestSanitize            = "http_request_sanitize"
	UpdateZoneRulesetOptions_Phase_HttpRequestSbfm                = "http_request_sbfm"
	UpdateZoneRulesetOptions_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	UpdateZoneRulesetOptions_Phase_HttpRequestTransform           = "http_request_transform"
	UpdateZoneRulesetOptions_Phase_HttpResponseCompression        = "http_response_compression"
	UpdateZoneRulesetOptions_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	UpdateZoneRulesetOptions_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// NewUpdateZoneRulesetOptions : Instantiate UpdateZoneRulesetOptions
func (*RulesetsV1) NewUpdateZoneRulesetOptions(rulesetID string) *UpdateZoneRulesetOptions {
	return &UpdateZoneRulesetOptions{
		RulesetID: core.StringPtr(rulesetID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *UpdateZoneRulesetOptions) SetRulesetID(rulesetID string) *UpdateZoneRulesetOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateZoneRulesetOptions) SetDescription(description string) *UpdateZoneRulesetOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *UpdateZoneRulesetOptions) SetKind(kind string) *UpdateZoneRulesetOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateZoneRulesetOptions) SetName(name string) *UpdateZoneRulesetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPhase : Allow user to set Phase
func (_options *UpdateZoneRulesetOptions) SetPhase(phase string) *UpdateZoneRulesetOptions {
	_options.Phase = core.StringPtr(phase)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *UpdateZoneRulesetOptions) SetRules(rules []RuleCreate) *UpdateZoneRulesetOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateZoneRulesetOptions) SetHeaders(param map[string]string) *UpdateZoneRulesetOptions {
	options.Headers = param
	return options
}

// UpdateZoneRulesetRuleOptions : The UpdateZoneRulesetRule options.
type UpdateZoneRulesetRuleOptions struct {
	// ID of a specific ruleset.
	RulesetID *string `json:"ruleset_id" validate:"required,ne="`

	// ID of a specific rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	Description *string `json:"description,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression,omitempty"`

	ID *string `json:"id,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Position *Position `json:"position,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateZoneRulesetRuleOptions : Instantiate UpdateZoneRulesetRuleOptions
func (*RulesetsV1) NewUpdateZoneRulesetRuleOptions(rulesetID string, ruleID string) *UpdateZoneRulesetRuleOptions {
	return &UpdateZoneRulesetRuleOptions{
		RulesetID: core.StringPtr(rulesetID),
		RuleID:    core.StringPtr(ruleID),
	}
}

// SetRulesetID : Allow user to set RulesetID
func (_options *UpdateZoneRulesetRuleOptions) SetRulesetID(rulesetID string) *UpdateZoneRulesetRuleOptions {
	_options.RulesetID = core.StringPtr(rulesetID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateZoneRulesetRuleOptions) SetRuleID(ruleID string) *UpdateZoneRulesetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *UpdateZoneRulesetRuleOptions) SetAction(action string) *UpdateZoneRulesetRuleOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetActionParameters : Allow user to set ActionParameters
func (_options *UpdateZoneRulesetRuleOptions) SetActionParameters(actionParameters *ActionParameters) *UpdateZoneRulesetRuleOptions {
	_options.ActionParameters = actionParameters
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateZoneRulesetRuleOptions) SetDescription(description string) *UpdateZoneRulesetRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateZoneRulesetRuleOptions) SetEnabled(enabled bool) *UpdateZoneRulesetRuleOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetExpression : Allow user to set Expression
func (_options *UpdateZoneRulesetRuleOptions) SetExpression(expression string) *UpdateZoneRulesetRuleOptions {
	_options.Expression = core.StringPtr(expression)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateZoneRulesetRuleOptions) SetID(id string) *UpdateZoneRulesetRuleOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetLogging : Allow user to set Logging
func (_options *UpdateZoneRulesetRuleOptions) SetLogging(logging *Logging) *UpdateZoneRulesetRuleOptions {
	_options.Logging = logging
	return _options
}

// SetRef : Allow user to set Ref
func (_options *UpdateZoneRulesetRuleOptions) SetRef(ref string) *UpdateZoneRulesetRuleOptions {
	_options.Ref = core.StringPtr(ref)
	return _options
}

// SetPosition : Allow user to set Position
func (_options *UpdateZoneRulesetRuleOptions) SetPosition(position *Position) *UpdateZoneRulesetRuleOptions {
	_options.Position = position
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateZoneRulesetRuleOptions) SetHeaders(param map[string]string) *UpdateZoneRulesetRuleOptions {
	options.Headers = param
	return options
}

// ActionParameters : ActionParameters struct
type ActionParameters struct {
	// unique ID of the ruleset.
	ID *string `json:"id,omitempty"`

	Overrides *Overrides `json:"overrides,omitempty"`

	// The version of the ruleset. Use "latest" to get the latest version.
	Version *string `json:"version,omitempty"`

	// Ruleset ID of the ruleset to apply action to. Use "current" to apply to the current ruleset.
	Ruleset *string `json:"ruleset,omitempty"`

	// List of ruleset ids to apply action to. Use "current" to apply to the current ruleset.
	Rulesets []string `json:"rulesets,omitempty"`

	Response *ActionParametersResponse `json:"response,omitempty"`
}

// UnmarshalActionParameters unmarshals an instance of ActionParameters from the specified map of raw messages.
func UnmarshalActionParameters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionParameters)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "overrides", &obj.Overrides, UnmarshalOverrides)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ruleset", &obj.Ruleset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rulesets", &obj.Rulesets)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalActionParametersResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CategoriesOverride : CategoriesOverride struct
type CategoriesOverride struct {
	// The category tag name to override.
	Category *string `json:"category,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`
}

// UnmarshalCategoriesOverride unmarshals an instance of CategoriesOverride from the specified map of raw messages.
func UnmarshalCategoriesOverride(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CategoriesOverride)
	err = core.UnmarshalPrimitive(m, "category", &obj.Category)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListRulesetsResp : List rulesets response.
type ListRulesetsResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []Message `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []Message `json:"messages" validate:"required"`

	// Container for response information.
	Result []ListedRuleset `json:"result" validate:"required"`
}

// UnmarshalListRulesetsResp unmarshals an instance of ListRulesetsResp from the specified map of raw messages.
func UnmarshalListRulesetsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListRulesetsResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalListedRuleset)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListedRuleset : ListedRuleset struct
type ListedRuleset struct {
	// description of the ruleset.
	Description *string `json:"description" validate:"required"`

	// unique ID of the ruleset.
	ID *string `json:"id" validate:"required"`

	Kind *string `json:"kind" validate:"required"`

	// The timestamp of when the resource was last modified.
	LastUpdated *string `json:"last_updated" validate:"required"`

	// human readable name of the ruleset.
	Name *string `json:"name" validate:"required"`

	// The phase of the ruleset.
	Phase *string `json:"phase" validate:"required"`

	// The version of the ruleset.
	Version *string `json:"version" validate:"required"`
}

// Constants associated with the ListedRuleset.Kind property.
const (
	ListedRuleset_Kind_Custom  = "custom"
	ListedRuleset_Kind_Managed = "managed"
	ListedRuleset_Kind_Root    = "root"
	ListedRuleset_Kind_Zone    = "zone"
)

// Constants associated with the ListedRuleset.Phase property.
// The phase of the ruleset.
const (
	ListedRuleset_Phase_DdosL4                         = "ddos_l4"
	ListedRuleset_Phase_DdosL7                         = "ddos_l7"
	ListedRuleset_Phase_HttpConfigSettings             = "http_config_settings"
	ListedRuleset_Phase_HttpCustomErrors               = "http_custom_errors"
	ListedRuleset_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	ListedRuleset_Phase_HttpRatelimit                  = "http_ratelimit"
	ListedRuleset_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	ListedRuleset_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	ListedRuleset_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	ListedRuleset_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	ListedRuleset_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	ListedRuleset_Phase_HttpRequestOrigin              = "http_request_origin"
	ListedRuleset_Phase_HttpRequestRedirect            = "http_request_redirect"
	ListedRuleset_Phase_HttpRequestSanitize            = "http_request_sanitize"
	ListedRuleset_Phase_HttpRequestSbfm                = "http_request_sbfm"
	ListedRuleset_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	ListedRuleset_Phase_HttpRequestTransform           = "http_request_transform"
	ListedRuleset_Phase_HttpResponseCompression        = "http_response_compression"
	ListedRuleset_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	ListedRuleset_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// UnmarshalListedRuleset unmarshals an instance of ListedRuleset from the specified map of raw messages.
func UnmarshalListedRuleset(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListedRuleset)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	err = core.UnmarshalPrimitive(m, "last_updated", &obj.LastUpdated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "phase", &obj.Phase)
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

// Logging : Logging struct
type Logging struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

// NewLogging : Instantiate Logging (Generic Model Constructor)
func (*RulesetsV1) NewLogging(enabled bool) (_model *Logging, err error) {
	_model = &Logging{
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalLogging unmarshals an instance of Logging from the specified map of raw messages.
func UnmarshalLogging(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Logging)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Message : Message struct
type Message struct {
	// A unique code for this message.
	Code *int64 `json:"code,omitempty"`

	// A text description of this message.
	Message *string `json:"message" validate:"required"`

	// The source of this message.
	Source *MessageSource `json:"source,omitempty"`
}

// UnmarshalMessage unmarshals an instance of Message from the specified map of raw messages.
func UnmarshalMessage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Message)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalMessageSource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Overrides : Overrides struct
type Overrides struct {
	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The sensitivity level of the rule.
	SensitivityLevel *string `json:"sensitivity_level,omitempty"`

	Rules []RulesOverride `json:"rules,omitempty"`

	Categories []CategoriesOverride `json:"categories,omitempty"`
}

// Constants associated with the Overrides.SensitivityLevel property.
// The sensitivity level of the rule.
const (
	Overrides_SensitivityLevel_High   = "high"
	Overrides_SensitivityLevel_Low    = "low"
	Overrides_SensitivityLevel_Medium = "medium"
)

// UnmarshalOverrides unmarshals an instance of Overrides from the specified map of raw messages.
func UnmarshalOverrides(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Overrides)
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sensitivity_level", &obj.SensitivityLevel)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRulesOverride)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "categories", &obj.Categories, UnmarshalCategoriesOverride)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Position : Position struct
type Position struct {
	// The rule ID to place this rule before.
	Before *string `json:"before,omitempty"`

	// The rule ID to place this rule after.
	After *string `json:"after,omitempty"`

	// The index to place this rule at.
	Index *int64 `json:"index,omitempty"`
}

// UnmarshalPosition unmarshals an instance of Position from the specified map of raw messages.
func UnmarshalPosition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Position)
	err = core.UnmarshalPrimitive(m, "before", &obj.Before)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "after", &obj.After)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "index", &obj.Index)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleCreate : RuleCreate struct
type RuleCreate struct {
	// What happens when theres a match for the rule expression.
	Action *string `json:"action" validate:"required"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	Description *string `json:"description,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression" validate:"required"`

	ID *string `json:"id,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Position *Position `json:"position,omitempty"`
}

// NewRuleCreate : Instantiate RuleCreate (Generic Model Constructor)
func (*RulesetsV1) NewRuleCreate(action string, expression string) (_model *RuleCreate, err error) {
	_model = &RuleCreate{
		Action:     core.StringPtr(action),
		Expression: core.StringPtr(expression),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleCreate unmarshals an instance of RuleCreate from the specified map of raw messages.
func UnmarshalRuleCreate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleCreate)
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_parameters", &obj.ActionParameters, UnmarshalActionParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expression", &obj.Expression)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "logging", &obj.Logging, UnmarshalLogging)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ref", &obj.Ref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "position", &obj.Position, UnmarshalPosition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleDetails : RuleDetails struct
type RuleDetails struct {
	// unique ID of rule.
	ID *string `json:"id" validate:"required"`

	// The version of the rule.
	Version *string `json:"version,omitempty"`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	ActionParameters *ActionParameters `json:"action_parameters,omitempty"`

	// List of categories for the rule.
	Categories []string `json:"categories,omitempty"`

	// Is the rule enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// description of the rule.
	Description *string `json:"description,omitempty"`

	// The expression defining which traffic will match the rule.
	Expression *string `json:"expression,omitempty"`

	// The reference of the rule (the rule ID by default).
	Ref *string `json:"ref,omitempty"`

	Logging *Logging `json:"logging,omitempty"`

	// The timestamp of when the resource was last modified.
	LastUpdated *string `json:"last_updated,omitempty"`
}

// UnmarshalRuleDetails unmarshals an instance of RuleDetails from the specified map of raw messages.
func UnmarshalRuleDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleDetails)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action_parameters", &obj.ActionParameters, UnmarshalActionParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "categories", &obj.Categories)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expression", &obj.Expression)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ref", &obj.Ref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "logging", &obj.Logging, UnmarshalLogging)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_updated", &obj.LastUpdated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleResp : List rules response.
type RuleResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []Message `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []Message `json:"messages" validate:"required"`

	Result *RuleDetails `json:"result" validate:"required"`
}

// UnmarshalRuleResp unmarshals an instance of RuleResp from the specified map of raw messages.
func UnmarshalRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRuleDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RulesOverride : RulesOverride struct
type RulesOverride struct {
	ID *string `json:"id,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	// What happens when theres a match for the rule expression.
	Action *string `json:"action,omitempty"`

	// The sensitivity level of the rule.
	SensitivityLevel *string `json:"sensitivity_level,omitempty"`
}

// Constants associated with the RulesOverride.SensitivityLevel property.
// The sensitivity level of the rule.
const (
	RulesOverride_SensitivityLevel_High   = "high"
	RulesOverride_SensitivityLevel_Low    = "low"
	RulesOverride_SensitivityLevel_Medium = "medium"
)

// UnmarshalRulesOverride unmarshals an instance of RulesOverride from the specified map of raw messages.
func UnmarshalRulesOverride(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RulesOverride)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sensitivity_level", &obj.SensitivityLevel)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RulesetDetails : RulesetDetails struct
type RulesetDetails struct {
	// description of the ruleset.
	Description *string `json:"description" validate:"required"`

	// unique ID of the ruleset.
	ID *string `json:"id" validate:"required"`

	Kind *string `json:"kind" validate:"required"`

	// The timestamp of when the resource was last modified.
	LastUpdated *string `json:"last_updated" validate:"required"`

	// human readable name of the ruleset.
	Name *string `json:"name" validate:"required"`

	// The phase of the ruleset.
	Phase *string `json:"phase" validate:"required"`

	// The version of the ruleset.
	Version *string `json:"version" validate:"required"`

	Rules []RuleDetails `json:"rules" validate:"required"`
}

// Constants associated with the RulesetDetails.Kind property.
const (
	RulesetDetails_Kind_Custom  = "custom"
	RulesetDetails_Kind_Managed = "managed"
	RulesetDetails_Kind_Root    = "root"
	RulesetDetails_Kind_Zone    = "zone"
)

// Constants associated with the RulesetDetails.Phase property.
// The phase of the ruleset.
const (
	RulesetDetails_Phase_DdosL4                         = "ddos_l4"
	RulesetDetails_Phase_DdosL7                         = "ddos_l7"
	RulesetDetails_Phase_HttpConfigSettings             = "http_config_settings"
	RulesetDetails_Phase_HttpCustomErrors               = "http_custom_errors"
	RulesetDetails_Phase_HttpLogCustomFields            = "http_log_custom_fields"
	RulesetDetails_Phase_HttpRatelimit                  = "http_ratelimit"
	RulesetDetails_Phase_HttpRequestCacheSettings       = "http_request_cache_settings"
	RulesetDetails_Phase_HttpRequestDynamicRedirect     = "http_request_dynamic_redirect"
	RulesetDetails_Phase_HttpRequestFirewallCustom      = "http_request_firewall_custom"
	RulesetDetails_Phase_HttpRequestFirewallManaged     = "http_request_firewall_managed"
	RulesetDetails_Phase_HttpRequestLateTransform       = "http_request_late_transform"
	RulesetDetails_Phase_HttpRequestOrigin              = "http_request_origin"
	RulesetDetails_Phase_HttpRequestRedirect            = "http_request_redirect"
	RulesetDetails_Phase_HttpRequestSanitize            = "http_request_sanitize"
	RulesetDetails_Phase_HttpRequestSbfm                = "http_request_sbfm"
	RulesetDetails_Phase_HttpRequestSelectConfiguration = "http_request_select_configuration"
	RulesetDetails_Phase_HttpRequestTransform           = "http_request_transform"
	RulesetDetails_Phase_HttpResponseCompression        = "http_response_compression"
	RulesetDetails_Phase_HttpResponseFirewallManaged    = "http_response_firewall_managed"
	RulesetDetails_Phase_HttpResponseHeadersTransform   = "http_response_headers_transform"
)

// UnmarshalRulesetDetails unmarshals an instance of RulesetDetails from the specified map of raw messages.
func UnmarshalRulesetDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RulesetDetails)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	err = core.UnmarshalPrimitive(m, "last_updated", &obj.LastUpdated)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "phase", &obj.Phase)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRuleDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RulesetResp : Ruleset response.
type RulesetResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors []Message `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages []Message `json:"messages" validate:"required"`

	Result *RulesetDetails `json:"result" validate:"required"`
}

// UnmarshalRulesetResp unmarshals an instance of RulesetResp from the specified map of raw messages.
func UnmarshalRulesetResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RulesetResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRulesetDetails)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
