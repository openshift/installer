/**
 * (C) Copyright IBM Corp. 2020.
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
 * IBM OpenAPI SDK Code Generator Version: 3.20.0-debb9f29-20201203-202043
 */
 

// Package zonefirewallaccessrulesv1 : Operations and models for the ZoneFirewallAccessRulesV1 service
package zonefirewallaccessrulesv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"net/http"
	"reflect"
	"time"
)

// ZoneFirewallAccessRulesV1 : Zone Firewall Access Rules
//
// Version: 1.0.1
type ZoneFirewallAccessRulesV1 struct {
	Service *core.BaseService

	// Full crn of the service instance.
	Crn *string

	// Identifier of zone whose access rule is to be deleted.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "zone_firewall_access_rules"

// ZoneFirewallAccessRulesV1Options : Service options
type ZoneFirewallAccessRulesV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full crn of the service instance.
	Crn *string `validate:"required"`

	// Identifier of zone whose access rule is to be deleted.
	ZoneIdentifier *string `validate:"required"`
}

// NewZoneFirewallAccessRulesV1UsingExternalConfig : constructs an instance of ZoneFirewallAccessRulesV1 with passed in options and external configuration.
func NewZoneFirewallAccessRulesV1UsingExternalConfig(options *ZoneFirewallAccessRulesV1Options) (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	zoneFirewallAccessRules, err = NewZoneFirewallAccessRulesV1(options)
	if err != nil {
		return
	}

	err = zoneFirewallAccessRules.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = zoneFirewallAccessRules.Service.SetServiceURL(options.URL)
	}
	return
}

// NewZoneFirewallAccessRulesV1 : constructs an instance of ZoneFirewallAccessRulesV1 with passed in options.
func NewZoneFirewallAccessRulesV1(options *ZoneFirewallAccessRulesV1Options) (service *ZoneFirewallAccessRulesV1, err error) {
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

	service = &ZoneFirewallAccessRulesV1{
		Service: baseService,
		Crn: options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "zoneFirewallAccessRules" suitable for processing requests.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) Clone() *ZoneFirewallAccessRulesV1 {
	if core.IsNil(zoneFirewallAccessRules) {
		return nil
	}
	clone := *zoneFirewallAccessRules
	clone.Service = zoneFirewallAccessRules.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) SetServiceURL(url string) error {
	return zoneFirewallAccessRules.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) GetServiceURL() string {
	return zoneFirewallAccessRules.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) SetDefaultHeaders(headers http.Header) {
	zoneFirewallAccessRules.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) SetEnableGzipCompression(enableGzip bool) {
	zoneFirewallAccessRules.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) GetEnableGzipCompression() bool {
	return zoneFirewallAccessRules.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	zoneFirewallAccessRules.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) DisableRetries() {
	zoneFirewallAccessRules.Service.DisableRetries()
}

// ListAllZoneAccessRules : List all firewall access rules
// List all firewall access rules for a zone.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) ListAllZoneAccessRules(listAllZoneAccessRulesOptions *ListAllZoneAccessRulesOptions) (result *ListZoneAccessRulesResp, response *core.DetailedResponse, err error) {
	return zoneFirewallAccessRules.ListAllZoneAccessRulesWithContext(context.Background(), listAllZoneAccessRulesOptions)
}

// ListAllZoneAccessRulesWithContext is an alternate form of the ListAllZoneAccessRules method which supports a Context parameter
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) ListAllZoneAccessRulesWithContext(ctx context.Context, listAllZoneAccessRulesOptions *ListAllZoneAccessRulesOptions) (result *ListZoneAccessRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllZoneAccessRulesOptions, "listAllZoneAccessRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneFirewallAccessRules.Crn,
		"zone_identifier": *zoneFirewallAccessRules.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneFirewallAccessRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneFirewallAccessRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/access_rules/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllZoneAccessRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_firewall_access_rules", "V1", "ListAllZoneAccessRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllZoneAccessRulesOptions.Notes != nil {
		builder.AddQuery("notes", fmt.Sprint(*listAllZoneAccessRulesOptions.Notes))
	}
	if listAllZoneAccessRulesOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*listAllZoneAccessRulesOptions.Mode))
	}
	if listAllZoneAccessRulesOptions.ConfigurationTarget != nil {
		builder.AddQuery("configuration.target", fmt.Sprint(*listAllZoneAccessRulesOptions.ConfigurationTarget))
	}
	if listAllZoneAccessRulesOptions.ConfigurationValue != nil {
		builder.AddQuery("configuration.value", fmt.Sprint(*listAllZoneAccessRulesOptions.ConfigurationValue))
	}
	if listAllZoneAccessRulesOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listAllZoneAccessRulesOptions.Page))
	}
	if listAllZoneAccessRulesOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listAllZoneAccessRulesOptions.PerPage))
	}
	if listAllZoneAccessRulesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listAllZoneAccessRulesOptions.Order))
	}
	if listAllZoneAccessRulesOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listAllZoneAccessRulesOptions.Direction))
	}
	if listAllZoneAccessRulesOptions.Match != nil {
		builder.AddQuery("match", fmt.Sprint(*listAllZoneAccessRulesOptions.Match))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneFirewallAccessRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListZoneAccessRulesResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneAccessRule : Create firewall access rule
// Create a new firewall access rule for a given zone under a service instance.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) CreateZoneAccessRule(createZoneAccessRuleOptions *CreateZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	return zoneFirewallAccessRules.CreateZoneAccessRuleWithContext(context.Background(), createZoneAccessRuleOptions)
}

// CreateZoneAccessRuleWithContext is an alternate form of the CreateZoneAccessRule method which supports a Context parameter
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) CreateZoneAccessRuleWithContext(ctx context.Context, createZoneAccessRuleOptions *CreateZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createZoneAccessRuleOptions, "createZoneAccessRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneFirewallAccessRules.Crn,
		"zone_identifier": *zoneFirewallAccessRules.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneFirewallAccessRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneFirewallAccessRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/access_rules/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createZoneAccessRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_firewall_access_rules", "V1", "CreateZoneAccessRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createZoneAccessRuleOptions.Mode != nil {
		body["mode"] = createZoneAccessRuleOptions.Mode
	}
	if createZoneAccessRuleOptions.Notes != nil {
		body["notes"] = createZoneAccessRuleOptions.Notes
	}
	if createZoneAccessRuleOptions.Configuration != nil {
		body["configuration"] = createZoneAccessRuleOptions.Configuration
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
	response, err = zoneFirewallAccessRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneAccessRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteZoneAccessRule : Delete firewall access rule
// Delete an access rule given its id.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) DeleteZoneAccessRule(deleteZoneAccessRuleOptions *DeleteZoneAccessRuleOptions) (result *DeleteZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	return zoneFirewallAccessRules.DeleteZoneAccessRuleWithContext(context.Background(), deleteZoneAccessRuleOptions)
}

// DeleteZoneAccessRuleWithContext is an alternate form of the DeleteZoneAccessRule method which supports a Context parameter
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) DeleteZoneAccessRuleWithContext(ctx context.Context, deleteZoneAccessRuleOptions *DeleteZoneAccessRuleOptions) (result *DeleteZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneAccessRuleOptions, "deleteZoneAccessRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneAccessRuleOptions, "deleteZoneAccessRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneFirewallAccessRules.Crn,
		"zone_identifier": *zoneFirewallAccessRules.ZoneIdentifier,
		"accessrule_identifier": *deleteZoneAccessRuleOptions.AccessruleIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneFirewallAccessRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneFirewallAccessRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/access_rules/rules/{accessrule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneAccessRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_firewall_access_rules", "V1", "DeleteZoneAccessRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneFirewallAccessRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteZoneAccessRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetZoneAccessRule : Get firewall access rule
// Get the details of a firewall access rule for a given zone under a given service instance.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) GetZoneAccessRule(getZoneAccessRuleOptions *GetZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	return zoneFirewallAccessRules.GetZoneAccessRuleWithContext(context.Background(), getZoneAccessRuleOptions)
}

// GetZoneAccessRuleWithContext is an alternate form of the GetZoneAccessRule method which supports a Context parameter
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) GetZoneAccessRuleWithContext(ctx context.Context, getZoneAccessRuleOptions *GetZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneAccessRuleOptions, "getZoneAccessRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneAccessRuleOptions, "getZoneAccessRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneFirewallAccessRules.Crn,
		"zone_identifier": *zoneFirewallAccessRules.ZoneIdentifier,
		"accessrule_identifier": *getZoneAccessRuleOptions.AccessruleIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneFirewallAccessRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneFirewallAccessRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/access_rules/rules/{accessrule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneAccessRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_firewall_access_rules", "V1", "GetZoneAccessRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneFirewallAccessRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneAccessRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateZoneAccessRule : Update firewall access rule
// Update an existing firewall access rule for a given zone under a given service instance.
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) UpdateZoneAccessRule(updateZoneAccessRuleOptions *UpdateZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	return zoneFirewallAccessRules.UpdateZoneAccessRuleWithContext(context.Background(), updateZoneAccessRuleOptions)
}

// UpdateZoneAccessRuleWithContext is an alternate form of the UpdateZoneAccessRule method which supports a Context parameter
func (zoneFirewallAccessRules *ZoneFirewallAccessRulesV1) UpdateZoneAccessRuleWithContext(ctx context.Context, updateZoneAccessRuleOptions *UpdateZoneAccessRuleOptions) (result *ZoneAccessRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateZoneAccessRuleOptions, "updateZoneAccessRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateZoneAccessRuleOptions, "updateZoneAccessRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneFirewallAccessRules.Crn,
		"zone_identifier": *zoneFirewallAccessRules.ZoneIdentifier,
		"accessrule_identifier": *updateZoneAccessRuleOptions.AccessruleIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneFirewallAccessRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneFirewallAccessRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/access_rules/rules/{accessrule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateZoneAccessRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_firewall_access_rules", "V1", "UpdateZoneAccessRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateZoneAccessRuleOptions.Mode != nil {
		body["mode"] = updateZoneAccessRuleOptions.Mode
	}
	if updateZoneAccessRuleOptions.Notes != nil {
		body["notes"] = updateZoneAccessRuleOptions.Notes
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
	response, err = zoneFirewallAccessRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneAccessRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneAccessRuleOptions : The CreateZoneAccessRule options.
type CreateZoneAccessRuleOptions struct {
	// The action to apply to a matched request.
	Mode *string `json:"mode,omitempty"`

	// A personal note about the rule. Typically used as a reminder or explanation for the rule.
	Notes *string `json:"notes,omitempty"`

	// Configuration object specifying access rule.
	Configuration *ZoneAccessRuleInputConfiguration `json:"configuration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateZoneAccessRuleOptions.Mode property.
// The action to apply to a matched request.
const (
	CreateZoneAccessRuleOptions_Mode_Block = "block"
	CreateZoneAccessRuleOptions_Mode_Challenge = "challenge"
	CreateZoneAccessRuleOptions_Mode_JsChallenge = "js_challenge"
	CreateZoneAccessRuleOptions_Mode_Whitelist = "whitelist"
)

// NewCreateZoneAccessRuleOptions : Instantiate CreateZoneAccessRuleOptions
func (*ZoneFirewallAccessRulesV1) NewCreateZoneAccessRuleOptions() *CreateZoneAccessRuleOptions {
	return &CreateZoneAccessRuleOptions{}
}

// SetMode : Allow user to set Mode
func (options *CreateZoneAccessRuleOptions) SetMode(mode string) *CreateZoneAccessRuleOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetNotes : Allow user to set Notes
func (options *CreateZoneAccessRuleOptions) SetNotes(notes string) *CreateZoneAccessRuleOptions {
	options.Notes = core.StringPtr(notes)
	return options
}

// SetConfiguration : Allow user to set Configuration
func (options *CreateZoneAccessRuleOptions) SetConfiguration(configuration *ZoneAccessRuleInputConfiguration) *CreateZoneAccessRuleOptions {
	options.Configuration = configuration
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneAccessRuleOptions) SetHeaders(param map[string]string) *CreateZoneAccessRuleOptions {
	options.Headers = param
	return options
}

// DeleteZoneAccessRuleOptions : The DeleteZoneAccessRule options.
type DeleteZoneAccessRuleOptions struct {
	// Identifier of the access rule to be deleted.
	AccessruleIdentifier *string `json:"accessrule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneAccessRuleOptions : Instantiate DeleteZoneAccessRuleOptions
func (*ZoneFirewallAccessRulesV1) NewDeleteZoneAccessRuleOptions(accessruleIdentifier string) *DeleteZoneAccessRuleOptions {
	return &DeleteZoneAccessRuleOptions{
		AccessruleIdentifier: core.StringPtr(accessruleIdentifier),
	}
}

// SetAccessruleIdentifier : Allow user to set AccessruleIdentifier
func (options *DeleteZoneAccessRuleOptions) SetAccessruleIdentifier(accessruleIdentifier string) *DeleteZoneAccessRuleOptions {
	options.AccessruleIdentifier = core.StringPtr(accessruleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneAccessRuleOptions) SetHeaders(param map[string]string) *DeleteZoneAccessRuleOptions {
	options.Headers = param
	return options
}

// DeleteZoneAccessRuleRespResult : Container for response information.
type DeleteZoneAccessRuleRespResult struct {
	// ID.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteZoneAccessRuleRespResult unmarshals an instance of DeleteZoneAccessRuleRespResult from the specified map of raw messages.
func UnmarshalDeleteZoneAccessRuleRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteZoneAccessRuleRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetZoneAccessRuleOptions : The GetZoneAccessRule options.
type GetZoneAccessRuleOptions struct {
	// Identifier of firewall access rule for the given zone.
	AccessruleIdentifier *string `json:"accessrule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneAccessRuleOptions : Instantiate GetZoneAccessRuleOptions
func (*ZoneFirewallAccessRulesV1) NewGetZoneAccessRuleOptions(accessruleIdentifier string) *GetZoneAccessRuleOptions {
	return &GetZoneAccessRuleOptions{
		AccessruleIdentifier: core.StringPtr(accessruleIdentifier),
	}
}

// SetAccessruleIdentifier : Allow user to set AccessruleIdentifier
func (options *GetZoneAccessRuleOptions) SetAccessruleIdentifier(accessruleIdentifier string) *GetZoneAccessRuleOptions {
	options.AccessruleIdentifier = core.StringPtr(accessruleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneAccessRuleOptions) SetHeaders(param map[string]string) *GetZoneAccessRuleOptions {
	options.Headers = param
	return options
}

// ListAllZoneAccessRulesOptions : The ListAllZoneAccessRules options.
type ListAllZoneAccessRulesOptions struct {
	// Search access rules by note.(Not case sensitive).
	Notes *string `json:"notes,omitempty"`

	// Search access rules by mode.
	Mode *string `json:"mode,omitempty"`

	// Search access rules by configuration target.
	ConfigurationTarget *string `json:"configuration.target,omitempty"`

	// Search access rules by configuration value which can be IP, IPrange, or country code.
	ConfigurationValue *string `json:"configuration.value,omitempty"`

	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of access rules per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Field by which to order list of access rules.
	Order *string `json:"order,omitempty"`

	// Direction in which to order results [ascending/descending order].
	Direction *string `json:"direction,omitempty"`

	// Whether to match all (all) or atleast one search parameter (any).
	Match *string `json:"match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListAllZoneAccessRulesOptions.Mode property.
// Search access rules by mode.
const (
	ListAllZoneAccessRulesOptions_Mode_Block = "block"
	ListAllZoneAccessRulesOptions_Mode_Challenge = "challenge"
	ListAllZoneAccessRulesOptions_Mode_JsChallenge = "js_challenge"
	ListAllZoneAccessRulesOptions_Mode_Whitelist = "whitelist"
)

// Constants associated with the ListAllZoneAccessRulesOptions.ConfigurationTarget property.
// Search access rules by configuration target.
const (
	ListAllZoneAccessRulesOptions_ConfigurationTarget_Asn = "asn"
	ListAllZoneAccessRulesOptions_ConfigurationTarget_Country = "country"
	ListAllZoneAccessRulesOptions_ConfigurationTarget_Ip = "ip"
	ListAllZoneAccessRulesOptions_ConfigurationTarget_IpRange = "ip_range"
)

// Constants associated with the ListAllZoneAccessRulesOptions.Order property.
// Field by which to order list of access rules.
const (
	ListAllZoneAccessRulesOptions_Order_ConfigurationTarget = "configuration.target"
	ListAllZoneAccessRulesOptions_Order_ConfigurationValue = "configuration.value"
	ListAllZoneAccessRulesOptions_Order_Mode = "mode"
)

// Constants associated with the ListAllZoneAccessRulesOptions.Direction property.
// Direction in which to order results [ascending/descending order].
const (
	ListAllZoneAccessRulesOptions_Direction_Asc = "asc"
	ListAllZoneAccessRulesOptions_Direction_Desc = "desc"
)

// Constants associated with the ListAllZoneAccessRulesOptions.Match property.
// Whether to match all (all) or atleast one search parameter (any).
const (
	ListAllZoneAccessRulesOptions_Match_All = "all"
	ListAllZoneAccessRulesOptions_Match_Any = "any"
)

// NewListAllZoneAccessRulesOptions : Instantiate ListAllZoneAccessRulesOptions
func (*ZoneFirewallAccessRulesV1) NewListAllZoneAccessRulesOptions() *ListAllZoneAccessRulesOptions {
	return &ListAllZoneAccessRulesOptions{}
}

// SetNotes : Allow user to set Notes
func (options *ListAllZoneAccessRulesOptions) SetNotes(notes string) *ListAllZoneAccessRulesOptions {
	options.Notes = core.StringPtr(notes)
	return options
}

// SetMode : Allow user to set Mode
func (options *ListAllZoneAccessRulesOptions) SetMode(mode string) *ListAllZoneAccessRulesOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetConfigurationTarget : Allow user to set ConfigurationTarget
func (options *ListAllZoneAccessRulesOptions) SetConfigurationTarget(configurationTarget string) *ListAllZoneAccessRulesOptions {
	options.ConfigurationTarget = core.StringPtr(configurationTarget)
	return options
}

// SetConfigurationValue : Allow user to set ConfigurationValue
func (options *ListAllZoneAccessRulesOptions) SetConfigurationValue(configurationValue string) *ListAllZoneAccessRulesOptions {
	options.ConfigurationValue = core.StringPtr(configurationValue)
	return options
}

// SetPage : Allow user to set Page
func (options *ListAllZoneAccessRulesOptions) SetPage(page int64) *ListAllZoneAccessRulesOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListAllZoneAccessRulesOptions) SetPerPage(perPage int64) *ListAllZoneAccessRulesOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListAllZoneAccessRulesOptions) SetOrder(order string) *ListAllZoneAccessRulesOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetDirection : Allow user to set Direction
func (options *ListAllZoneAccessRulesOptions) SetDirection(direction string) *ListAllZoneAccessRulesOptions {
	options.Direction = core.StringPtr(direction)
	return options
}

// SetMatch : Allow user to set Match
func (options *ListAllZoneAccessRulesOptions) SetMatch(match string) *ListAllZoneAccessRulesOptions {
	options.Match = core.StringPtr(match)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllZoneAccessRulesOptions) SetHeaders(param map[string]string) *ListAllZoneAccessRulesOptions {
	options.Headers = param
	return options
}

// ListZoneAccessRulesRespResultInfo : Statistics of results.
type ListZoneAccessRulesRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalListZoneAccessRulesRespResultInfo unmarshals an instance of ListZoneAccessRulesRespResultInfo from the specified map of raw messages.
func UnmarshalListZoneAccessRulesRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListZoneAccessRulesRespResultInfo)
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

// UpdateZoneAccessRuleOptions : The UpdateZoneAccessRule options.
type UpdateZoneAccessRuleOptions struct {
	// Identifier of firewall access rule.
	AccessruleIdentifier *string `json:"accessrule_identifier" validate:"required,ne="`

	// The action to apply to a matched request.
	Mode *string `json:"mode,omitempty"`

	// A personal note about the rule. Typically used as a reminder or explanation for the rule.
	Notes *string `json:"notes,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateZoneAccessRuleOptions.Mode property.
// The action to apply to a matched request.
const (
	UpdateZoneAccessRuleOptions_Mode_Block = "block"
	UpdateZoneAccessRuleOptions_Mode_Challenge = "challenge"
	UpdateZoneAccessRuleOptions_Mode_JsChallenge = "js_challenge"
	UpdateZoneAccessRuleOptions_Mode_Whitelist = "whitelist"
)

// NewUpdateZoneAccessRuleOptions : Instantiate UpdateZoneAccessRuleOptions
func (*ZoneFirewallAccessRulesV1) NewUpdateZoneAccessRuleOptions(accessruleIdentifier string) *UpdateZoneAccessRuleOptions {
	return &UpdateZoneAccessRuleOptions{
		AccessruleIdentifier: core.StringPtr(accessruleIdentifier),
	}
}

// SetAccessruleIdentifier : Allow user to set AccessruleIdentifier
func (options *UpdateZoneAccessRuleOptions) SetAccessruleIdentifier(accessruleIdentifier string) *UpdateZoneAccessRuleOptions {
	options.AccessruleIdentifier = core.StringPtr(accessruleIdentifier)
	return options
}

// SetMode : Allow user to set Mode
func (options *UpdateZoneAccessRuleOptions) SetMode(mode string) *UpdateZoneAccessRuleOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetNotes : Allow user to set Notes
func (options *UpdateZoneAccessRuleOptions) SetNotes(notes string) *UpdateZoneAccessRuleOptions {
	options.Notes = core.StringPtr(notes)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateZoneAccessRuleOptions) SetHeaders(param map[string]string) *UpdateZoneAccessRuleOptions {
	options.Headers = param
	return options
}

// ZoneAccessRuleInputConfiguration : Configuration object specifying access rule.
type ZoneAccessRuleInputConfiguration struct {
	// The request property to target.
	Target *string `json:"target" validate:"required"`

	// The value for the selected target.For ip the value is a valid ip address.For ip_range the value specifies ip range
	// limited to /16 and /24. For asn the value is an AS number. For country the value is a country code for the country.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the ZoneAccessRuleInputConfiguration.Target property.
// The request property to target.
const (
	ZoneAccessRuleInputConfiguration_Target_Asn = "asn"
	ZoneAccessRuleInputConfiguration_Target_Country = "country"
	ZoneAccessRuleInputConfiguration_Target_Ip = "ip"
	ZoneAccessRuleInputConfiguration_Target_IpRange = "ip_range"
)


// NewZoneAccessRuleInputConfiguration : Instantiate ZoneAccessRuleInputConfiguration (Generic Model Constructor)
func (*ZoneFirewallAccessRulesV1) NewZoneAccessRuleInputConfiguration(target string, value string) (model *ZoneAccessRuleInputConfiguration, err error) {
	model = &ZoneAccessRuleInputConfiguration{
		Target: core.StringPtr(target),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalZoneAccessRuleInputConfiguration unmarshals an instance of ZoneAccessRuleInputConfiguration from the specified map of raw messages.
func UnmarshalZoneAccessRuleInputConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneAccessRuleInputConfiguration)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// ZoneAccessRuleObjectConfiguration : configuration.
type ZoneAccessRuleObjectConfiguration struct {
	// target.
	Target *string `json:"target" validate:"required"`

	// Value for the given target. For ip the value is a valid ip address.For ip_range the value specifies ip range limited
	// to /16 and /24. For asn the value is an AS number. For country the value is a country code for the country.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the ZoneAccessRuleObjectConfiguration.Target property.
// target.
const (
	ZoneAccessRuleObjectConfiguration_Target_Asn = "asn"
	ZoneAccessRuleObjectConfiguration_Target_Country = "country"
	ZoneAccessRuleObjectConfiguration_Target_Ip = "ip"
	ZoneAccessRuleObjectConfiguration_Target_IpRange = "ip_range"
)


// UnmarshalZoneAccessRuleObjectConfiguration unmarshals an instance of ZoneAccessRuleObjectConfiguration from the specified map of raw messages.
func UnmarshalZoneAccessRuleObjectConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneAccessRuleObjectConfiguration)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// ZoneAccessRuleObjectScope : The scope definition of the access rule.
type ZoneAccessRuleObjectScope struct {
	// The scope of the access rule, indicating if its applicable at zone level("zone") or inherited from instance
	// level("account").
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the ZoneAccessRuleObjectScope.Type property.
// The scope of the access rule, indicating if its applicable at zone level("zone") or inherited from instance
// level("account").
const (
	ZoneAccessRuleObjectScope_Type_Account = "account"
	ZoneAccessRuleObjectScope_Type_Zone = "zone"
)


// UnmarshalZoneAccessRuleObjectScope unmarshals an instance of ZoneAccessRuleObjectScope from the specified map of raw messages.
func UnmarshalZoneAccessRuleObjectScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneAccessRuleObjectScope)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteZoneAccessRuleResp : delete access rule response.
type DeleteZoneAccessRuleResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *DeleteZoneAccessRuleRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteZoneAccessRuleResp unmarshals an instance of DeleteZoneAccessRuleResp from the specified map of raw messages.
func UnmarshalDeleteZoneAccessRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteZoneAccessRuleResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteZoneAccessRuleRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListZoneAccessRulesResp : list access rules response.
type ListZoneAccessRulesResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []ZoneAccessRuleObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListZoneAccessRulesRespResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListZoneAccessRulesResp unmarshals an instance of ListZoneAccessRulesResp from the specified map of raw messages.
func UnmarshalListZoneAccessRulesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListZoneAccessRulesResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalZoneAccessRuleObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListZoneAccessRulesRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ZoneAccessRuleObject : access rule object.
type ZoneAccessRuleObject struct {
	// Identifier of the firewall access rule.
	ID *string `json:"id" validate:"required"`

	// A personal note about the rule. Typically used as a reminder or explanation for the rule.
	Notes *string `json:"notes" validate:"required"`

	// List of modes that are allowed.
	AllowedModes []string `json:"allowed_modes" validate:"required"`

	// The action to be applied to a request matching the access rule.
	Mode *string `json:"mode" validate:"required"`

	// The scope definition of the access rule.
	Scope *ZoneAccessRuleObjectScope `json:"scope,omitempty"`

	// The creation date-time of the firewall access rule.
	CreatedOn *string `json:"created_on" validate:"required"`

	// The modification date-time of the firewall access rule.
	ModifiedOn *string `json:"modified_on" validate:"required"`

	// configuration.
	Configuration *ZoneAccessRuleObjectConfiguration `json:"configuration" validate:"required"`
}

// Constants associated with the ZoneAccessRuleObject.AllowedModes property.
const (
	ZoneAccessRuleObject_AllowedModes_Block = "block"
	ZoneAccessRuleObject_AllowedModes_Challenge = "challenge"
	ZoneAccessRuleObject_AllowedModes_JsChallenge = "js_challenge"
	ZoneAccessRuleObject_AllowedModes_Whitelist = "whitelist"
)

// Constants associated with the ZoneAccessRuleObject.Mode property.
// The action to be applied to a request matching the access rule.
const (
	ZoneAccessRuleObject_Mode_Block = "block"
	ZoneAccessRuleObject_Mode_Challenge = "challenge"
	ZoneAccessRuleObject_Mode_JsChallenge = "js_challenge"
	ZoneAccessRuleObject_Mode_Whitelist = "whitelist"
)


// UnmarshalZoneAccessRuleObject unmarshals an instance of ZoneAccessRuleObject from the specified map of raw messages.
func UnmarshalZoneAccessRuleObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneAccessRuleObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "notes", &obj.Notes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_modes", &obj.AllowedModes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalZoneAccessRuleObjectScope)
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
	err = core.UnmarshalModel(m, "configuration", &obj.Configuration, UnmarshalZoneAccessRuleObjectConfiguration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ZoneAccessRuleResp : access rule response.
type ZoneAccessRuleResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// access rule object.
	Result *ZoneAccessRuleObject `json:"result" validate:"required"`
}


// UnmarshalZoneAccessRuleResp unmarshals an instance of ZoneAccessRuleResp from the specified map of raw messages.
func UnmarshalZoneAccessRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneAccessRuleResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalZoneAccessRuleObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
