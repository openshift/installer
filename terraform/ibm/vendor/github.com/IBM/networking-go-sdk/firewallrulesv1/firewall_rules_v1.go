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

// Package firewallrulesv1 : Operations and models for the FirewallRulesV1 service
package firewallrulesv1

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

// FirewallRulesV1 : Firewall rules
//
// Version: 1.0.1
type FirewallRulesV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "firewall_rules"

// FirewallRulesV1Options : Service options
type FirewallRulesV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewFirewallRulesV1UsingExternalConfig : constructs an instance of FirewallRulesV1 with passed in options and external configuration.
func NewFirewallRulesV1UsingExternalConfig(options *FirewallRulesV1Options) (firewallRules *FirewallRulesV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	firewallRules, err = NewFirewallRulesV1(options)
	if err != nil {
		return
	}

	err = firewallRules.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = firewallRules.Service.SetServiceURL(options.URL)
	}
	return
}

// NewFirewallRulesV1 : constructs an instance of FirewallRulesV1 with passed in options.
func NewFirewallRulesV1(options *FirewallRulesV1Options) (service *FirewallRulesV1, err error) {
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

	service = &FirewallRulesV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "firewallRules" suitable for processing requests.
func (firewallRules *FirewallRulesV1) Clone() *FirewallRulesV1 {
	if core.IsNil(firewallRules) {
		return nil
	}
	clone := *firewallRules
	clone.Service = firewallRules.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (firewallRules *FirewallRulesV1) SetServiceURL(url string) error {
	return firewallRules.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (firewallRules *FirewallRulesV1) GetServiceURL() string {
	return firewallRules.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (firewallRules *FirewallRulesV1) SetDefaultHeaders(headers http.Header) {
	firewallRules.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (firewallRules *FirewallRulesV1) SetEnableGzipCompression(enableGzip bool) {
	firewallRules.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (firewallRules *FirewallRulesV1) GetEnableGzipCompression() bool {
	return firewallRules.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (firewallRules *FirewallRulesV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	firewallRules.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (firewallRules *FirewallRulesV1) DisableRetries() {
	firewallRules.Service.DisableRetries()
}

// ListAllFirewallRules : List all firewall rules for a zone
// List all firewall rules for a zone.
func (firewallRules *FirewallRulesV1) ListAllFirewallRules(listAllFirewallRulesOptions *ListAllFirewallRulesOptions) (result *ListFirewallRulesResp, response *core.DetailedResponse, err error) {
	return firewallRules.ListAllFirewallRulesWithContext(context.Background(), listAllFirewallRulesOptions)
}

// ListAllFirewallRulesWithContext is an alternate form of the ListAllFirewallRules method which supports a Context parameter
func (firewallRules *FirewallRulesV1) ListAllFirewallRulesWithContext(ctx context.Context, listAllFirewallRulesOptions *ListAllFirewallRulesOptions) (result *ListFirewallRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAllFirewallRulesOptions, "listAllFirewallRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAllFirewallRulesOptions, "listAllFirewallRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *listAllFirewallRulesOptions.Crn,
		"zone_identifier": *listAllFirewallRulesOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllFirewallRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "ListAllFirewallRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAllFirewallRulesOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*listAllFirewallRulesOptions.XAuthUserToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListFirewallRulesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateFirewallRules : Create firewall rules for a zone
// Create new firewall rules for a given zone under a service instance.
func (firewallRules *FirewallRulesV1) CreateFirewallRules(createFirewallRulesOptions *CreateFirewallRulesOptions) (result *FirewallRulesResp, response *core.DetailedResponse, err error) {
	return firewallRules.CreateFirewallRulesWithContext(context.Background(), createFirewallRulesOptions)
}

// CreateFirewallRulesWithContext is an alternate form of the CreateFirewallRules method which supports a Context parameter
func (firewallRules *FirewallRulesV1) CreateFirewallRulesWithContext(ctx context.Context, createFirewallRulesOptions *CreateFirewallRulesOptions) (result *FirewallRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createFirewallRulesOptions, "createFirewallRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createFirewallRulesOptions, "createFirewallRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *createFirewallRulesOptions.Crn,
		"zone_identifier": *createFirewallRulesOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createFirewallRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "CreateFirewallRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createFirewallRulesOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*createFirewallRulesOptions.XAuthUserToken))
	}

	if createFirewallRulesOptions.FirewallRuleInputWithFilterID != nil {
		_, err = builder.SetBodyContentJSON(createFirewallRulesOptions.FirewallRuleInputWithFilterID)
		if err != nil {
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFirewallRulesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateFirewllRules : Update firewall rules
// Update existing firewall rules for a given zone under a given service instance.
func (firewallRules *FirewallRulesV1) UpdateFirewllRules(updateFirewllRulesOptions *UpdateFirewllRulesOptions) (result *FirewallRulesResp, response *core.DetailedResponse, err error) {
	return firewallRules.UpdateFirewllRulesWithContext(context.Background(), updateFirewllRulesOptions)
}

// UpdateFirewllRulesWithContext is an alternate form of the UpdateFirewllRules method which supports a Context parameter
func (firewallRules *FirewallRulesV1) UpdateFirewllRulesWithContext(ctx context.Context, updateFirewllRulesOptions *UpdateFirewllRulesOptions) (result *FirewallRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFirewllRulesOptions, "updateFirewllRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateFirewllRulesOptions, "updateFirewllRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *updateFirewllRulesOptions.Crn,
		"zone_identifier": *updateFirewllRulesOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateFirewllRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "UpdateFirewllRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateFirewllRulesOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*updateFirewllRulesOptions.XAuthUserToken))
	}

	if updateFirewllRulesOptions.FirewallRulesUpdateInputItem != nil {
		_, err = builder.SetBodyContentJSON(updateFirewllRulesOptions.FirewallRulesUpdateInputItem)
		if err != nil {
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFirewallRulesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteFirewallRules : Delete firewall rules
// Delete firewall rules by filter ids.
func (firewallRules *FirewallRulesV1) DeleteFirewallRules(deleteFirewallRulesOptions *DeleteFirewallRulesOptions) (result *DeleteFirewallRulesResp, response *core.DetailedResponse, err error) {
	return firewallRules.DeleteFirewallRulesWithContext(context.Background(), deleteFirewallRulesOptions)
}

// DeleteFirewallRulesWithContext is an alternate form of the DeleteFirewallRules method which supports a Context parameter
func (firewallRules *FirewallRulesV1) DeleteFirewallRulesWithContext(ctx context.Context, deleteFirewallRulesOptions *DeleteFirewallRulesOptions) (result *DeleteFirewallRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteFirewallRulesOptions, "deleteFirewallRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteFirewallRulesOptions, "deleteFirewallRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *deleteFirewallRulesOptions.Crn,
		"zone_identifier": *deleteFirewallRulesOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteFirewallRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "DeleteFirewallRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteFirewallRulesOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*deleteFirewallRulesOptions.XAuthUserToken))
	}

	builder.AddQuery("id", fmt.Sprint(*deleteFirewallRulesOptions.ID))

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteFirewallRulesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteFirewallRule : Delete a firewall rule
// Delete a firewall rule given its id.
func (firewallRules *FirewallRulesV1) DeleteFirewallRule(deleteFirewallRuleOptions *DeleteFirewallRuleOptions) (result *DeleteFirewallRuleResp, response *core.DetailedResponse, err error) {
	return firewallRules.DeleteFirewallRuleWithContext(context.Background(), deleteFirewallRuleOptions)
}

// DeleteFirewallRuleWithContext is an alternate form of the DeleteFirewallRule method which supports a Context parameter
func (firewallRules *FirewallRulesV1) DeleteFirewallRuleWithContext(ctx context.Context, deleteFirewallRuleOptions *DeleteFirewallRuleOptions) (result *DeleteFirewallRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteFirewallRuleOptions, "deleteFirewallRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteFirewallRuleOptions, "deleteFirewallRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *deleteFirewallRuleOptions.Crn,
		"zone_identifier": *deleteFirewallRuleOptions.ZoneIdentifier,
		"firewall_rule_identifier": *deleteFirewallRuleOptions.FirewallRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules/{firewall_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteFirewallRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "DeleteFirewallRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteFirewallRuleOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*deleteFirewallRuleOptions.XAuthUserToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteFirewallRuleResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetFirewallRule : Get firewall rule details by id
// Get the details of a firewall rule for a given zone under a given service instance.
func (firewallRules *FirewallRulesV1) GetFirewallRule(getFirewallRuleOptions *GetFirewallRuleOptions) (result *FirewallRuleResp, response *core.DetailedResponse, err error) {
	return firewallRules.GetFirewallRuleWithContext(context.Background(), getFirewallRuleOptions)
}

// GetFirewallRuleWithContext is an alternate form of the GetFirewallRule method which supports a Context parameter
func (firewallRules *FirewallRulesV1) GetFirewallRuleWithContext(ctx context.Context, getFirewallRuleOptions *GetFirewallRuleOptions) (result *FirewallRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getFirewallRuleOptions, "getFirewallRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getFirewallRuleOptions, "getFirewallRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *getFirewallRuleOptions.Crn,
		"zone_identifier": *getFirewallRuleOptions.ZoneIdentifier,
		"firewall_rule_identifier": *getFirewallRuleOptions.FirewallRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules/{firewall_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getFirewallRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "GetFirewallRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getFirewallRuleOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*getFirewallRuleOptions.XAuthUserToken))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFirewallRuleResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateFirewallRule : Update a firewall rule
// Update an existing firewall rule for a given zone under a given service instance.
func (firewallRules *FirewallRulesV1) UpdateFirewallRule(updateFirewallRuleOptions *UpdateFirewallRuleOptions) (result *FirewallRuleResp, response *core.DetailedResponse, err error) {
	return firewallRules.UpdateFirewallRuleWithContext(context.Background(), updateFirewallRuleOptions)
}

// UpdateFirewallRuleWithContext is an alternate form of the UpdateFirewallRule method which supports a Context parameter
func (firewallRules *FirewallRulesV1) UpdateFirewallRuleWithContext(ctx context.Context, updateFirewallRuleOptions *UpdateFirewallRuleOptions) (result *FirewallRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFirewallRuleOptions, "updateFirewallRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateFirewallRuleOptions, "updateFirewallRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *updateFirewallRuleOptions.Crn,
		"zone_identifier": *updateFirewallRuleOptions.ZoneIdentifier,
		"firewall_rule_identifier": *updateFirewallRuleOptions.FirewallRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = firewallRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(firewallRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/rules/{firewall_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateFirewallRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("firewall_rules", "V1", "UpdateFirewallRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateFirewallRuleOptions.XAuthUserToken != nil {
		builder.AddHeader("X-Auth-User-Token", fmt.Sprint(*updateFirewallRuleOptions.XAuthUserToken))
	}

	body := make(map[string]interface{})
	if updateFirewallRuleOptions.Action != nil {
		body["action"] = updateFirewallRuleOptions.Action
	}
	if updateFirewallRuleOptions.Paused != nil {
		body["paused"] = updateFirewallRuleOptions.Paused
	}
	if updateFirewallRuleOptions.Description != nil {
		body["description"] = updateFirewallRuleOptions.Description
	}
	if updateFirewallRuleOptions.Filter != nil {
		body["filter"] = updateFirewallRuleOptions.Filter
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
	response, err = firewallRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFirewallRuleResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateFirewallRulesOptions : The CreateFirewallRules options.
type CreateFirewallRulesOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `validate:"required,ne="`

	// Zone identifier of the zone for which firewall rules are created.
	ZoneIdentifier *string `validate:"required,ne="`

	// Json objects which are used to create firewall rules.
	FirewallRuleInputWithFilterID []FirewallRuleInputWithFilterID

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateFirewallRulesOptions : Instantiate CreateFirewallRulesOptions
func (*FirewallRulesV1) NewCreateFirewallRulesOptions(xAuthUserToken string, crn string, zoneIdentifier string) *CreateFirewallRulesOptions {
	return &CreateFirewallRulesOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *CreateFirewallRulesOptions) SetXAuthUserToken(xAuthUserToken string) *CreateFirewallRulesOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *CreateFirewallRulesOptions) SetCrn(crn string) *CreateFirewallRulesOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *CreateFirewallRulesOptions) SetZoneIdentifier(zoneIdentifier string) *CreateFirewallRulesOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetFirewallRuleInputWithFilterID : Allow user to set FirewallRuleInputWithFilterID
func (options *CreateFirewallRulesOptions) SetFirewallRuleInputWithFilterID(firewallRuleInputWithFilterID []FirewallRuleInputWithFilterID) *CreateFirewallRulesOptions {
	options.FirewallRuleInputWithFilterID = firewallRuleInputWithFilterID
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateFirewallRulesOptions) SetHeaders(param map[string]string) *CreateFirewallRulesOptions {
	options.Headers = param
	return options
}

// DeleteFirewallRuleOptions : The DeleteFirewallRule options.
type DeleteFirewallRuleOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full crn of the service instance.
	Crn *string `validate:"required,ne="`

	// Identifier of zone whose firewall rule is to be deleted.
	ZoneIdentifier *string `validate:"required,ne="`

	// Identifier of the firewall rule to be deleted.
	FirewallRuleIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteFirewallRuleOptions : Instantiate DeleteFirewallRuleOptions
func (*FirewallRulesV1) NewDeleteFirewallRuleOptions(xAuthUserToken string, crn string, zoneIdentifier string, firewallRuleIdentifier string) *DeleteFirewallRuleOptions {
	return &DeleteFirewallRuleOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		FirewallRuleIdentifier: core.StringPtr(firewallRuleIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *DeleteFirewallRuleOptions) SetXAuthUserToken(xAuthUserToken string) *DeleteFirewallRuleOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *DeleteFirewallRuleOptions) SetCrn(crn string) *DeleteFirewallRuleOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *DeleteFirewallRuleOptions) SetZoneIdentifier(zoneIdentifier string) *DeleteFirewallRuleOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetFirewallRuleIdentifier : Allow user to set FirewallRuleIdentifier
func (options *DeleteFirewallRuleOptions) SetFirewallRuleIdentifier(firewallRuleIdentifier string) *DeleteFirewallRuleOptions {
	options.FirewallRuleIdentifier = core.StringPtr(firewallRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteFirewallRuleOptions) SetHeaders(param map[string]string) *DeleteFirewallRuleOptions {
	options.Headers = param
	return options
}

// DeleteFirewallRuleRespResult : Container for response information.
type DeleteFirewallRuleRespResult struct {
	// Identifier of the firewall rule.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalDeleteFirewallRuleRespResult unmarshals an instance of DeleteFirewallRuleRespResult from the specified map of raw messages.
func UnmarshalDeleteFirewallRuleRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteFirewallRuleRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteFirewallRulesOptions : The DeleteFirewallRules options.
type DeleteFirewallRulesOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full crn of the service instance.
	Crn *string `validate:"required,ne="`

	// Identifier of zone whose firewall rules are to be deleted.
	ZoneIdentifier *string `validate:"required,ne="`

	// ids of firewall rules which will be deleted.
	ID *string `validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteFirewallRulesOptions : Instantiate DeleteFirewallRulesOptions
func (*FirewallRulesV1) NewDeleteFirewallRulesOptions(xAuthUserToken string, crn string, zoneIdentifier string, id string) *DeleteFirewallRulesOptions {
	return &DeleteFirewallRulesOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		ID: core.StringPtr(id),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *DeleteFirewallRulesOptions) SetXAuthUserToken(xAuthUserToken string) *DeleteFirewallRulesOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *DeleteFirewallRulesOptions) SetCrn(crn string) *DeleteFirewallRulesOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *DeleteFirewallRulesOptions) SetZoneIdentifier(zoneIdentifier string) *DeleteFirewallRulesOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *DeleteFirewallRulesOptions) SetID(id string) *DeleteFirewallRulesOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteFirewallRulesOptions) SetHeaders(param map[string]string) *DeleteFirewallRulesOptions {
	options.Headers = param
	return options
}

// DeleteFirewallRulesRespResultItem : DeleteFirewallRulesRespResultItem struct
type DeleteFirewallRulesRespResultItem struct {
	// Identifier of firewall rules.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalDeleteFirewallRulesRespResultItem unmarshals an instance of DeleteFirewallRulesRespResultItem from the specified map of raw messages.
func UnmarshalDeleteFirewallRulesRespResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteFirewallRulesRespResultItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRuleInputWithFilterIdFilter : An existing filter.
type FirewallRuleInputWithFilterIdFilter struct {
	// Identifier of the filter.
	ID *string `json:"id" validate:"required"`
}

// NewFirewallRuleInputWithFilterIdFilter : Instantiate FirewallRuleInputWithFilterIdFilter (Generic Model Constructor)
func (*FirewallRulesV1) NewFirewallRuleInputWithFilterIdFilter(id string) (model *FirewallRuleInputWithFilterIdFilter, err error) {
	model = &FirewallRuleInputWithFilterIdFilter{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFirewallRuleInputWithFilterIdFilter unmarshals an instance of FirewallRuleInputWithFilterIdFilter from the specified map of raw messages.
func UnmarshalFirewallRuleInputWithFilterIdFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleInputWithFilterIdFilter)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRuleObjectFilter : An existing filter.
type FirewallRuleObjectFilter struct {
	// Identifier of the filter.
	ID *string `json:"id" validate:"required"`

	// Indicates if the filter is active.
	Paused *bool `json:"paused" validate:"required"`

	// To briefly describe the filter, omitted from object if empty.
	Description *string `json:"description" validate:"required"`

	// A filter expression.
	Expression *string `json:"expression" validate:"required"`
}

// UnmarshalFirewallRuleObjectFilter unmarshals an instance of FirewallRuleObjectFilter from the specified map of raw messages.
func UnmarshalFirewallRuleObjectFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleObjectFilter)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paused", &obj.Paused)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRuleUpdateInputFilter : An existing filter.
type FirewallRuleUpdateInputFilter struct {
	// Identifier of the filter.
	ID *string `json:"id" validate:"required"`
}

// NewFirewallRuleUpdateInputFilter : Instantiate FirewallRuleUpdateInputFilter (Generic Model Constructor)
func (*FirewallRulesV1) NewFirewallRuleUpdateInputFilter(id string) (model *FirewallRuleUpdateInputFilter, err error) {
	model = &FirewallRuleUpdateInputFilter{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFirewallRuleUpdateInputFilter unmarshals an instance of FirewallRuleUpdateInputFilter from the specified map of raw messages.
func UnmarshalFirewallRuleUpdateInputFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleUpdateInputFilter)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRulesUpdateInputItem : FirewallRulesUpdateInputItem struct
type FirewallRulesUpdateInputItem struct {
	// Identifier of the firewall rule.
	ID *string `json:"id" validate:"required"`

	// The firewall action to perform, "log" action is only available for enterprise plan instances.
	Action *string `json:"action" validate:"required"`

	// Indicates if the firewall rule is active.
	Paused *bool `json:"paused,omitempty"`

	// To briefly describe the firewall rule, omitted from object if empty.
	Description *string `json:"description,omitempty"`

	// An existing filter.
	Filter *FirewallRulesUpdateInputItemFilter `json:"filter,omitempty"`
}

// Constants associated with the FirewallRulesUpdateInputItem.Action property.
// The firewall action to perform, "log" action is only available for enterprise plan instances.
const (
	FirewallRulesUpdateInputItem_Action_Allow = "allow"
	FirewallRulesUpdateInputItem_Action_Block = "block"
	FirewallRulesUpdateInputItem_Action_Challenge = "challenge"
	FirewallRulesUpdateInputItem_Action_JsChallenge = "js_challenge"
	FirewallRulesUpdateInputItem_Action_Log = "log"
)

// NewFirewallRulesUpdateInputItem : Instantiate FirewallRulesUpdateInputItem (Generic Model Constructor)
func (*FirewallRulesV1) NewFirewallRulesUpdateInputItem(id string, action string) (model *FirewallRulesUpdateInputItem, err error) {
	model = &FirewallRulesUpdateInputItem{
		ID: core.StringPtr(id),
		Action: core.StringPtr(action),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFirewallRulesUpdateInputItem unmarshals an instance of FirewallRulesUpdateInputItem from the specified map of raw messages.
func UnmarshalFirewallRulesUpdateInputItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRulesUpdateInputItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paused", &obj.Paused)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "filter", &obj.Filter, UnmarshalFirewallRulesUpdateInputItemFilter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRulesUpdateInputItemFilter : An existing filter.
type FirewallRulesUpdateInputItemFilter struct {
	// Identifier of the filter.
	ID *string `json:"id" validate:"required"`
}

// NewFirewallRulesUpdateInputItemFilter : Instantiate FirewallRulesUpdateInputItemFilter (Generic Model Constructor)
func (*FirewallRulesV1) NewFirewallRulesUpdateInputItemFilter(id string) (model *FirewallRulesUpdateInputItemFilter, err error) {
	model = &FirewallRulesUpdateInputItemFilter{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFirewallRulesUpdateInputItemFilter unmarshals an instance of FirewallRulesUpdateInputItemFilter from the specified map of raw messages.
func UnmarshalFirewallRulesUpdateInputItemFilter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRulesUpdateInputItemFilter)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetFirewallRuleOptions : The GetFirewallRule options.
type GetFirewallRuleOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full crn of the service instance.
	Crn *string `validate:"required,ne="`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required,ne="`

	// Identifier of firewall rule for the given zone.
	FirewallRuleIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetFirewallRuleOptions : Instantiate GetFirewallRuleOptions
func (*FirewallRulesV1) NewGetFirewallRuleOptions(xAuthUserToken string, crn string, zoneIdentifier string, firewallRuleIdentifier string) *GetFirewallRuleOptions {
	return &GetFirewallRuleOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		FirewallRuleIdentifier: core.StringPtr(firewallRuleIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *GetFirewallRuleOptions) SetXAuthUserToken(xAuthUserToken string) *GetFirewallRuleOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *GetFirewallRuleOptions) SetCrn(crn string) *GetFirewallRuleOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *GetFirewallRuleOptions) SetZoneIdentifier(zoneIdentifier string) *GetFirewallRuleOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetFirewallRuleIdentifier : Allow user to set FirewallRuleIdentifier
func (options *GetFirewallRuleOptions) SetFirewallRuleIdentifier(firewallRuleIdentifier string) *GetFirewallRuleOptions {
	options.FirewallRuleIdentifier = core.StringPtr(firewallRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetFirewallRuleOptions) SetHeaders(param map[string]string) *GetFirewallRuleOptions {
	options.Headers = param
	return options
}

// ListAllFirewallRulesOptions : The ListAllFirewallRules options.
type ListAllFirewallRulesOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `validate:"required,ne="`

	// Zone identifier of the zone for which firewall rules are listed.
	ZoneIdentifier *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllFirewallRulesOptions : Instantiate ListAllFirewallRulesOptions
func (*FirewallRulesV1) NewListAllFirewallRulesOptions(xAuthUserToken string, crn string, zoneIdentifier string) *ListAllFirewallRulesOptions {
	return &ListAllFirewallRulesOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *ListAllFirewallRulesOptions) SetXAuthUserToken(xAuthUserToken string) *ListAllFirewallRulesOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *ListAllFirewallRulesOptions) SetCrn(crn string) *ListAllFirewallRulesOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *ListAllFirewallRulesOptions) SetZoneIdentifier(zoneIdentifier string) *ListAllFirewallRulesOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllFirewallRulesOptions) SetHeaders(param map[string]string) *ListAllFirewallRulesOptions {
	options.Headers = param
	return options
}

// ListFirewallRulesRespResultInfo : Statistics of results.
type ListFirewallRulesRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// UnmarshalListFirewallRulesRespResultInfo unmarshals an instance of ListFirewallRulesRespResultInfo from the specified map of raw messages.
func UnmarshalListFirewallRulesRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListFirewallRulesRespResultInfo)
	err = core.UnmarshalPrimitive(m, "page", &obj.Page)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "per_page", &obj.PerPage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateFirewallRuleOptions : The UpdateFirewallRule options.
type UpdateFirewallRuleOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full crn of the service instance.
	Crn *string `validate:"required,ne="`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required,ne="`

	// Identifier of firewall rule.
	FirewallRuleIdentifier *string `validate:"required,ne="`

	// The firewall action to perform, "log" action is only available for enterprise plan instances.
	Action *string

	// Indicates if the firewall rule is active.
	Paused *bool

	// To briefly describe the firewall rule, omitted from object if empty.
	Description *string

	// An existing filter.
	Filter *FirewallRuleUpdateInputFilter

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateFirewallRuleOptions.Action property.
// The firewall action to perform, "log" action is only available for enterprise plan instances.
const (
	UpdateFirewallRuleOptions_Action_Allow = "allow"
	UpdateFirewallRuleOptions_Action_Block = "block"
	UpdateFirewallRuleOptions_Action_Challenge = "challenge"
	UpdateFirewallRuleOptions_Action_JsChallenge = "js_challenge"
	UpdateFirewallRuleOptions_Action_Log = "log"
)

// NewUpdateFirewallRuleOptions : Instantiate UpdateFirewallRuleOptions
func (*FirewallRulesV1) NewUpdateFirewallRuleOptions(xAuthUserToken string, crn string, zoneIdentifier string, firewallRuleIdentifier string) *UpdateFirewallRuleOptions {
	return &UpdateFirewallRuleOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		FirewallRuleIdentifier: core.StringPtr(firewallRuleIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *UpdateFirewallRuleOptions) SetXAuthUserToken(xAuthUserToken string) *UpdateFirewallRuleOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *UpdateFirewallRuleOptions) SetCrn(crn string) *UpdateFirewallRuleOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *UpdateFirewallRuleOptions) SetZoneIdentifier(zoneIdentifier string) *UpdateFirewallRuleOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetFirewallRuleIdentifier : Allow user to set FirewallRuleIdentifier
func (options *UpdateFirewallRuleOptions) SetFirewallRuleIdentifier(firewallRuleIdentifier string) *UpdateFirewallRuleOptions {
	options.FirewallRuleIdentifier = core.StringPtr(firewallRuleIdentifier)
	return options
}

// SetAction : Allow user to set Action
func (options *UpdateFirewallRuleOptions) SetAction(action string) *UpdateFirewallRuleOptions {
	options.Action = core.StringPtr(action)
	return options
}

// SetPaused : Allow user to set Paused
func (options *UpdateFirewallRuleOptions) SetPaused(paused bool) *UpdateFirewallRuleOptions {
	options.Paused = core.BoolPtr(paused)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateFirewallRuleOptions) SetDescription(description string) *UpdateFirewallRuleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetFilter : Allow user to set Filter
func (options *UpdateFirewallRuleOptions) SetFilter(filter *FirewallRuleUpdateInputFilter) *UpdateFirewallRuleOptions {
	options.Filter = filter
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFirewallRuleOptions) SetHeaders(param map[string]string) *UpdateFirewallRuleOptions {
	options.Headers = param
	return options
}

// UpdateFirewllRulesOptions : The UpdateFirewllRules options.
type UpdateFirewllRulesOptions struct {
	// IBM Cloud user IAM token.
	XAuthUserToken *string `validate:"required"`

	// Full crn of the service instance.
	Crn *string `validate:"required,ne="`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required,ne="`

	FirewallRulesUpdateInputItem []FirewallRulesUpdateInputItem

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateFirewllRulesOptions : Instantiate UpdateFirewllRulesOptions
func (*FirewallRulesV1) NewUpdateFirewllRulesOptions(xAuthUserToken string, crn string, zoneIdentifier string) *UpdateFirewllRulesOptions {
	return &UpdateFirewllRulesOptions{
		XAuthUserToken: core.StringPtr(xAuthUserToken),
		Crn: core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
	}
}

// SetXAuthUserToken : Allow user to set XAuthUserToken
func (options *UpdateFirewllRulesOptions) SetXAuthUserToken(xAuthUserToken string) *UpdateFirewllRulesOptions {
	options.XAuthUserToken = core.StringPtr(xAuthUserToken)
	return options
}

// SetCrn : Allow user to set Crn
func (options *UpdateFirewllRulesOptions) SetCrn(crn string) *UpdateFirewllRulesOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (options *UpdateFirewllRulesOptions) SetZoneIdentifier(zoneIdentifier string) *UpdateFirewllRulesOptions {
	options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return options
}

// SetFirewallRulesUpdateInputItem : Allow user to set FirewallRulesUpdateInputItem
func (options *UpdateFirewllRulesOptions) SetFirewallRulesUpdateInputItem(firewallRulesUpdateInputItem []FirewallRulesUpdateInputItem) *UpdateFirewllRulesOptions {
	options.FirewallRulesUpdateInputItem = firewallRulesUpdateInputItem
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFirewllRulesOptions) SetHeaders(param map[string]string) *UpdateFirewllRulesOptions {
	options.Headers = param
	return options
}

// DeleteFirewallRuleResp : DeleteFirewallRuleResp struct
type DeleteFirewallRuleResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *DeleteFirewallRuleRespResult `json:"result" validate:"required"`
}

// UnmarshalDeleteFirewallRuleResp unmarshals an instance of DeleteFirewallRuleResp from the specified map of raw messages.
func UnmarshalDeleteFirewallRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteFirewallRuleResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteFirewallRuleRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteFirewallRulesResp : DeleteFirewallRulesResp struct
type DeleteFirewallRulesResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []DeleteFirewallRulesRespResultItem `json:"result" validate:"required"`
}

// UnmarshalDeleteFirewallRulesResp unmarshals an instance of DeleteFirewallRulesResp from the specified map of raw messages.
func UnmarshalDeleteFirewallRulesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteFirewallRulesResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteFirewallRulesRespResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRuleInputWithFilterID : Json objects which are used to create firewall rule.
type FirewallRuleInputWithFilterID struct {
	// An existing filter.
	Filter *FirewallRuleInputWithFilterIdFilter `json:"filter" validate:"required"`

	// The firewall action to perform, "log" action is only available for enterprise plan instances.
	Action *string `json:"action" validate:"required"`

	// To briefly describe the firewall rule, omitted from object if empty.
	Description *string `json:"description,omitempty"`
}

// Constants associated with the FirewallRuleInputWithFilterID.Action property.
// The firewall action to perform, "log" action is only available for enterprise plan instances.
const (
	FirewallRuleInputWithFilterID_Action_Allow = "allow"
	FirewallRuleInputWithFilterID_Action_Block = "block"
	FirewallRuleInputWithFilterID_Action_Challenge = "challenge"
	FirewallRuleInputWithFilterID_Action_JsChallenge = "js_challenge"
	FirewallRuleInputWithFilterID_Action_Log = "log"
)

// NewFirewallRuleInputWithFilterID : Instantiate FirewallRuleInputWithFilterID (Generic Model Constructor)
func (*FirewallRulesV1) NewFirewallRuleInputWithFilterID(filter *FirewallRuleInputWithFilterIdFilter, action string) (model *FirewallRuleInputWithFilterID, err error) {
	model = &FirewallRuleInputWithFilterID{
		Filter: filter,
		Action: core.StringPtr(action),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFirewallRuleInputWithFilterID unmarshals an instance of FirewallRuleInputWithFilterID from the specified map of raw messages.
func UnmarshalFirewallRuleInputWithFilterID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleInputWithFilterID)
	err = core.UnmarshalModel(m, "filter", &obj.Filter, UnmarshalFirewallRuleInputWithFilterIdFilter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
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

// FirewallRuleObject : FirewallRuleObject struct
type FirewallRuleObject struct {
	// Identifier of the firewall rule.
	ID *string `json:"id" validate:"required"`

	// Indicates if the firewall rule is active.
	Paused *bool `json:"paused" validate:"required"`

	// To briefly describe the firewall rule, omitted from object if empty.
	Description *string `json:"description" validate:"required"`

	// The firewall action to perform, "log" action is only available for enterprise plan instances.
	Action *string `json:"action" validate:"required"`

	// An existing filter.
	Filter *FirewallRuleObjectFilter `json:"filter" validate:"required"`

	// The creation date-time of the filter.
	CreatedOn *string `json:"created_on" validate:"required"`

	// The modification date-time of the filter.
	ModifiedOn *string `json:"modified_on" validate:"required"`
}

// Constants associated with the FirewallRuleObject.Action property.
// The firewall action to perform, "log" action is only available for enterprise plan instances.
const (
	FirewallRuleObject_Action_Allow = "allow"
	FirewallRuleObject_Action_Block = "block"
	FirewallRuleObject_Action_Challenge = "challenge"
	FirewallRuleObject_Action_JsChallenge = "js_challenge"
	FirewallRuleObject_Action_Log = "log"
)

// UnmarshalFirewallRuleObject unmarshals an instance of FirewallRuleObject from the specified map of raw messages.
func UnmarshalFirewallRuleObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paused", &obj.Paused)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "filter", &obj.Filter, UnmarshalFirewallRuleObjectFilter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRuleResp : FirewallRuleResp struct
type FirewallRuleResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	Result *FirewallRuleObject `json:"result" validate:"required"`
}

// UnmarshalFirewallRuleResp unmarshals an instance of FirewallRuleResp from the specified map of raw messages.
func UnmarshalFirewallRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRuleResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalFirewallRuleObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FirewallRulesResp : FirewallRulesResp struct
type FirewallRulesResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []FirewallRuleObject `json:"result" validate:"required"`
}

// UnmarshalFirewallRulesResp unmarshals an instance of FirewallRulesResp from the specified map of raw messages.
func UnmarshalFirewallRulesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FirewallRulesResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalFirewallRuleObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListFirewallRulesResp : ListFirewallRulesResp struct
type ListFirewallRulesResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []FirewallRuleObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListFirewallRulesRespResultInfo `json:"result_info" validate:"required"`
}

// UnmarshalListFirewallRulesResp unmarshals an instance of ListFirewallRulesResp from the specified map of raw messages.
func UnmarshalListFirewallRulesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListFirewallRulesResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalFirewallRuleObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListFirewallRulesRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
