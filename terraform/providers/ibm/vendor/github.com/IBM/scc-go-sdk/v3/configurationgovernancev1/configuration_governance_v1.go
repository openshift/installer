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
 * IBM OpenAPI SDK Code Generator Version: 3.34.0-e2a502a2-20210616-185634
 */

// Package configurationgovernancev1 : Operations and models for the ConfigurationGovernanceV1 service
package configurationgovernancev1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/v3/common"
	"github.com/go-openapi/strfmt"
)

// ConfigurationGovernanceV1 : API specification for the Configuration Governance service.
//
// Version: 1.0.0
type ConfigurationGovernanceV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us.compliance.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "configuration_governance"

// ConfigurationGovernanceV1Options : Service options
type ConfigurationGovernanceV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewConfigurationGovernanceV1UsingExternalConfig : constructs an instance of ConfigurationGovernanceV1 with passed in options and external configuration.
func NewConfigurationGovernanceV1UsingExternalConfig(options *ConfigurationGovernanceV1Options) (configurationGovernance *ConfigurationGovernanceV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	configurationGovernance, err = NewConfigurationGovernanceV1(options)
	if err != nil {
		return
	}

	err = configurationGovernance.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = configurationGovernance.Service.SetServiceURL(options.URL)
	}
	return
}

// NewConfigurationGovernanceV1 : constructs an instance of ConfigurationGovernanceV1 with passed in options.
func NewConfigurationGovernanceV1(options *ConfigurationGovernanceV1Options) (service *ConfigurationGovernanceV1, err error) {
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

	service = &ConfigurationGovernanceV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://us.compliance.cloud.ibm.com",
		"us-east":  "https://us.compliance.cloud.ibm.com",
		"eu-de":    "https://eu.compliance.cloud.ibm.com",
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "configurationGovernance" suitable for processing requests.
func (configurationGovernance *ConfigurationGovernanceV1) Clone() *ConfigurationGovernanceV1 {
	if core.IsNil(configurationGovernance) {
		return nil
	}
	clone := *configurationGovernance
	clone.Service = configurationGovernance.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (configurationGovernance *ConfigurationGovernanceV1) SetServiceURL(url string) error {
	return configurationGovernance.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (configurationGovernance *ConfigurationGovernanceV1) GetServiceURL() string {
	return configurationGovernance.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (configurationGovernance *ConfigurationGovernanceV1) SetDefaultHeaders(headers http.Header) {
	configurationGovernance.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (configurationGovernance *ConfigurationGovernanceV1) SetEnableGzipCompression(enableGzip bool) {
	configurationGovernance.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (configurationGovernance *ConfigurationGovernanceV1) GetEnableGzipCompression() bool {
	return configurationGovernance.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (configurationGovernance *ConfigurationGovernanceV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	configurationGovernance.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (configurationGovernance *ConfigurationGovernanceV1) DisableRetries() {
	configurationGovernance.Service.DisableRetries()
}

// CreateRules : Create rules
// Creates one or more rules that you can use to govern the way that IBM Cloud resources can be provisioned and
// configured.
//
// A successful `POST /config/rules` request defines a rule based on the target, conditions, and enforcement actions
// that you specify. The response returns the ID value for your rule, along with other metadata.
//
// To learn more about rules, check out the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule).
func (configurationGovernance *ConfigurationGovernanceV1) CreateRules(createRulesOptions *CreateRulesOptions) (result *CreateRulesResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.CreateRulesWithContext(context.Background(), createRulesOptions)
}

// CreateRulesWithContext is an alternate form of the CreateRules method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) CreateRulesWithContext(ctx context.Context, createRulesOptions *CreateRulesOptions) (result *CreateRulesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRulesOptions, "createRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRulesOptions, "createRulesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "CreateRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createRulesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createRulesOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createRulesOptions.Rules != nil {
		body["rules"] = createRulesOptions.Rules
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateRulesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRules : List rules
// Retrieves a list of the rules that are available in your account.
func (configurationGovernance *ConfigurationGovernanceV1) ListRules(listRulesOptions *ListRulesOptions) (result *RuleList, response *core.DetailedResponse, err error) {
	return configurationGovernance.ListRulesWithContext(context.Background(), listRulesOptions)
}

// ListRulesWithContext is an alternate form of the ListRules method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) ListRulesWithContext(ctx context.Context, listRulesOptions *ListRulesOptions) (result *RuleList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRulesOptions, "listRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listRulesOptions, "listRulesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "ListRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listRulesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listRulesOptions.TransactionID))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listRulesOptions.AccountID))
	if listRulesOptions.Attached != nil {
		builder.AddQuery("attached", fmt.Sprint(*listRulesOptions.Attached))
	}
	if listRulesOptions.Labels != nil {
		builder.AddQuery("labels", fmt.Sprint(*listRulesOptions.Labels))
	}
	if listRulesOptions.Scopes != nil {
		builder.AddQuery("scopes", fmt.Sprint(*listRulesOptions.Scopes))
	}
	if listRulesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listRulesOptions.Limit))
	}
	if listRulesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listRulesOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRule : Get a rule
// Retrieves an existing rule and its details.
func (configurationGovernance *ConfigurationGovernanceV1) GetRule(getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return configurationGovernance.GetRuleWithContext(context.Background(), getRuleOptions)
}

// GetRuleWithContext is an alternate form of the GetRule method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) GetRuleWithContext(ctx context.Context, getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRuleOptions, "getRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRuleOptions, "getRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *getRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "GetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getRuleOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateRule : Update a rule
// Updates an existing rule based on the properties that you specify.
func (configurationGovernance *ConfigurationGovernanceV1) UpdateRule(updateRuleOptions *UpdateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return configurationGovernance.UpdateRuleWithContext(context.Background(), updateRuleOptions)
}

// UpdateRuleWithContext is an alternate form of the UpdateRule method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) UpdateRuleWithContext(ctx context.Context, updateRuleOptions *UpdateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRuleOptions, "updateRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRuleOptions, "updateRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *updateRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "UpdateRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateRuleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateRuleOptions.IfMatch))
	}
	if updateRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateRuleOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateRuleOptions.Name != nil {
		body["name"] = updateRuleOptions.Name
	}
	if updateRuleOptions.Description != nil {
		body["description"] = updateRuleOptions.Description
	}
	if updateRuleOptions.Target != nil {
		body["target"] = updateRuleOptions.Target
	}
	if updateRuleOptions.RequiredConfig != nil {
		body["required_config"] = updateRuleOptions.RequiredConfig
	}
	if updateRuleOptions.EnforcementActions != nil {
		body["enforcement_actions"] = updateRuleOptions.EnforcementActions
	}
	if updateRuleOptions.AccountID != nil {
		body["account_id"] = updateRuleOptions.AccountID
	}
	if updateRuleOptions.RuleType != nil {
		body["rule_type"] = updateRuleOptions.RuleType
	}
	if updateRuleOptions.Labels != nil {
		body["labels"] = updateRuleOptions.Labels
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRule : Delete a rule
// Deletes an existing rule.
func (configurationGovernance *ConfigurationGovernanceV1) DeleteRule(deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	return configurationGovernance.DeleteRuleWithContext(context.Background(), deleteRuleOptions)
}

// DeleteRuleWithContext is an alternate form of the DeleteRule method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) DeleteRuleWithContext(ctx context.Context, deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRuleOptions, "deleteRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRuleOptions, "deleteRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *deleteRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "DeleteRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteRuleOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = configurationGovernance.Service.Request(request, nil)

	return
}

// CreateRuleAttachments : Create attachments
// Creates one or more scope attachments for an existing rule.
//
// You can attach an existing rule to a scope, such as a specific IBM Cloud account, to start evaluating the rule for
// compliance. A successful
// `POST /config/v1/rules/{rule_id}/attachments` returns the ID value for the attachment, along with other metadata.
func (configurationGovernance *ConfigurationGovernanceV1) CreateRuleAttachments(createRuleAttachmentsOptions *CreateRuleAttachmentsOptions) (result *CreateRuleAttachmentsResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.CreateRuleAttachmentsWithContext(context.Background(), createRuleAttachmentsOptions)
}

// CreateRuleAttachmentsWithContext is an alternate form of the CreateRuleAttachments method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) CreateRuleAttachmentsWithContext(ctx context.Context, createRuleAttachmentsOptions *CreateRuleAttachmentsOptions) (result *CreateRuleAttachmentsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRuleAttachmentsOptions, "createRuleAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRuleAttachmentsOptions, "createRuleAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *createRuleAttachmentsOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRuleAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "CreateRuleAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createRuleAttachmentsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createRuleAttachmentsOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createRuleAttachmentsOptions.Attachments != nil {
		body["attachments"] = createRuleAttachmentsOptions.Attachments
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateRuleAttachmentsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRuleAttachments : List attachments
// Retrieves a list of scope attachments that are associated with the specified rule.
func (configurationGovernance *ConfigurationGovernanceV1) ListRuleAttachments(listRuleAttachmentsOptions *ListRuleAttachmentsOptions) (result *RuleAttachmentList, response *core.DetailedResponse, err error) {
	return configurationGovernance.ListRuleAttachmentsWithContext(context.Background(), listRuleAttachmentsOptions)
}

// ListRuleAttachmentsWithContext is an alternate form of the ListRuleAttachments method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) ListRuleAttachmentsWithContext(ctx context.Context, listRuleAttachmentsOptions *ListRuleAttachmentsOptions) (result *RuleAttachmentList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRuleAttachmentsOptions, "listRuleAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listRuleAttachmentsOptions, "listRuleAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *listRuleAttachmentsOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRuleAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "ListRuleAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listRuleAttachmentsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listRuleAttachmentsOptions.TransactionID))
	}

	if listRuleAttachmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listRuleAttachmentsOptions.Limit))
	}
	if listRuleAttachmentsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listRuleAttachmentsOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleAttachmentList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRuleAttachment : Get an attachment
// Retrieves an existing scope attachment for a rule.
func (configurationGovernance *ConfigurationGovernanceV1) GetRuleAttachment(getRuleAttachmentOptions *GetRuleAttachmentOptions) (result *RuleAttachment, response *core.DetailedResponse, err error) {
	return configurationGovernance.GetRuleAttachmentWithContext(context.Background(), getRuleAttachmentOptions)
}

// GetRuleAttachmentWithContext is an alternate form of the GetRuleAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) GetRuleAttachmentWithContext(ctx context.Context, getRuleAttachmentOptions *GetRuleAttachmentOptions) (result *RuleAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRuleAttachmentOptions, "getRuleAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRuleAttachmentOptions, "getRuleAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id":       *getRuleAttachmentOptions.RuleID,
		"attachment_id": *getRuleAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRuleAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "GetRuleAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getRuleAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getRuleAttachmentOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateRuleAttachment : Update an attachment
// Updates an existing scope attachment based on the properties that you specify.
func (configurationGovernance *ConfigurationGovernanceV1) UpdateRuleAttachment(updateRuleAttachmentOptions *UpdateRuleAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	return configurationGovernance.UpdateRuleAttachmentWithContext(context.Background(), updateRuleAttachmentOptions)
}

// UpdateRuleAttachmentWithContext is an alternate form of the UpdateRuleAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) UpdateRuleAttachmentWithContext(ctx context.Context, updateRuleAttachmentOptions *UpdateRuleAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRuleAttachmentOptions, "updateRuleAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRuleAttachmentOptions, "updateRuleAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id":       *updateRuleAttachmentOptions.RuleID,
		"attachment_id": *updateRuleAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRuleAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "UpdateRuleAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateRuleAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateRuleAttachmentOptions.IfMatch))
	}
	if updateRuleAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateRuleAttachmentOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateRuleAttachmentOptions.AccountID != nil {
		body["account_id"] = updateRuleAttachmentOptions.AccountID
	}
	if updateRuleAttachmentOptions.IncludedScope != nil {
		body["included_scope"] = updateRuleAttachmentOptions.IncludedScope
	}
	if updateRuleAttachmentOptions.ExcludedScopes != nil {
		body["excluded_scopes"] = updateRuleAttachmentOptions.ExcludedScopes
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRuleAttachment : Delete an attachment
// Deletes an existing scope attachment.
func (configurationGovernance *ConfigurationGovernanceV1) DeleteRuleAttachment(deleteRuleAttachmentOptions *DeleteRuleAttachmentOptions) (response *core.DetailedResponse, err error) {
	return configurationGovernance.DeleteRuleAttachmentWithContext(context.Background(), deleteRuleAttachmentOptions)
}

// DeleteRuleAttachmentWithContext is an alternate form of the DeleteRuleAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) DeleteRuleAttachmentWithContext(ctx context.Context, deleteRuleAttachmentOptions *DeleteRuleAttachmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRuleAttachmentOptions, "deleteRuleAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRuleAttachmentOptions, "deleteRuleAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"rule_id":       *deleteRuleAttachmentOptions.RuleID,
		"attachment_id": *deleteRuleAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/rules/{rule_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRuleAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "DeleteRuleAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteRuleAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteRuleAttachmentOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = configurationGovernance.Service.Request(request, nil)

	return
}

// CreateTemplates : Create templates
// Creates one or more templates that you can use to define your preferred property values for IBM Cloud resources.
//
// A successful `POST /config/templates` request defines a template based on the target and customized defaults that you
// specify. The response returns the ID value for your template, along with other metadata.
//
// To learn more about templates, check out the
// [docs](/docs/security-compliance?topic=security-compliance-what-is-template).
func (configurationGovernance *ConfigurationGovernanceV1) CreateTemplates(createTemplatesOptions *CreateTemplatesOptions) (result *CreateTemplatesResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.CreateTemplatesWithContext(context.Background(), createTemplatesOptions)
}

// CreateTemplatesWithContext is an alternate form of the CreateTemplates method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) CreateTemplatesWithContext(ctx context.Context, createTemplatesOptions *CreateTemplatesOptions) (result *CreateTemplatesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTemplatesOptions, "createTemplatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTemplatesOptions, "createTemplatesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "CreateTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createTemplatesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createTemplatesOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createTemplatesOptions.Templates != nil {
		body["templates"] = createTemplatesOptions.Templates
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateTemplatesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListTemplates : List templates
// Retrieves a list of the templates that are available in your account.
func (configurationGovernance *ConfigurationGovernanceV1) ListTemplates(listTemplatesOptions *ListTemplatesOptions) (result *TemplateList, response *core.DetailedResponse, err error) {
	return configurationGovernance.ListTemplatesWithContext(context.Background(), listTemplatesOptions)
}

// ListTemplatesWithContext is an alternate form of the ListTemplates method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) ListTemplatesWithContext(ctx context.Context, listTemplatesOptions *ListTemplatesOptions) (result *TemplateList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTemplatesOptions, "listTemplatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTemplatesOptions, "listTemplatesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listTemplatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "ListTemplates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listTemplatesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listTemplatesOptions.TransactionID))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listTemplatesOptions.AccountID))
	if listTemplatesOptions.Attached != nil {
		builder.AddQuery("attached", fmt.Sprint(*listTemplatesOptions.Attached))
	}
	if listTemplatesOptions.Scopes != nil {
		builder.AddQuery("scopes", fmt.Sprint(*listTemplatesOptions.Scopes))
	}
	if listTemplatesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTemplatesOptions.Limit))
	}
	if listTemplatesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listTemplatesOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTemplate : Get a template
// Retrieves an existing template and its details.
func (configurationGovernance *ConfigurationGovernanceV1) GetTemplate(getTemplateOptions *GetTemplateOptions) (result *TemplateResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.GetTemplateWithContext(context.Background(), getTemplateOptions)
}

// GetTemplateWithContext is an alternate form of the GetTemplate method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) GetTemplateWithContext(ctx context.Context, getTemplateOptions *GetTemplateOptions) (result *TemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTemplateOptions, "getTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTemplateOptions, "getTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *getTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "GetTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getTemplateOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getTemplateOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateTemplate : Update a template
// Updates an existing template based on the properties that you specify.
func (configurationGovernance *ConfigurationGovernanceV1) UpdateTemplate(updateTemplateOptions *UpdateTemplateOptions) (result *TemplateResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.UpdateTemplateWithContext(context.Background(), updateTemplateOptions)
}

// UpdateTemplateWithContext is an alternate form of the UpdateTemplate method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) UpdateTemplateWithContext(ctx context.Context, updateTemplateOptions *UpdateTemplateOptions) (result *TemplateResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTemplateOptions, "updateTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateTemplateOptions, "updateTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *updateTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "UpdateTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateTemplateOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateTemplateOptions.IfMatch))
	}
	if updateTemplateOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateTemplateOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateTemplateOptions.AccountID != nil {
		body["account_id"] = updateTemplateOptions.AccountID
	}
	if updateTemplateOptions.Name != nil {
		body["name"] = updateTemplateOptions.Name
	}
	if updateTemplateOptions.Description != nil {
		body["description"] = updateTemplateOptions.Description
	}
	if updateTemplateOptions.Target != nil {
		body["target"] = updateTemplateOptions.Target
	}
	if updateTemplateOptions.CustomizedDefaults != nil {
		body["customized_defaults"] = updateTemplateOptions.CustomizedDefaults
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTemplate : Delete a template
// Deletes an existing template.
func (configurationGovernance *ConfigurationGovernanceV1) DeleteTemplate(deleteTemplateOptions *DeleteTemplateOptions) (response *core.DetailedResponse, err error) {
	return configurationGovernance.DeleteTemplateWithContext(context.Background(), deleteTemplateOptions)
}

// DeleteTemplateWithContext is an alternate form of the DeleteTemplate method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) DeleteTemplateWithContext(ctx context.Context, deleteTemplateOptions *DeleteTemplateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTemplateOptions, "deleteTemplateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTemplateOptions, "deleteTemplateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *deleteTemplateOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteTemplateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "DeleteTemplate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteTemplateOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteTemplateOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = configurationGovernance.Service.Request(request, nil)

	return
}

// CreateTemplateAttachments : Create attachments
// Creates one or more scope attachments for an existing template.
//
// You can attach an existing template to a scope, such as a specific IBM Cloud account, to start using the template for
// setting default values. A successful
// `POST /config/v1/templates/{template_id}/attachments` returns the ID value for the attachment, along with other
// metadata.
func (configurationGovernance *ConfigurationGovernanceV1) CreateTemplateAttachments(createTemplateAttachmentsOptions *CreateTemplateAttachmentsOptions) (result *CreateTemplateAttachmentsResponse, response *core.DetailedResponse, err error) {
	return configurationGovernance.CreateTemplateAttachmentsWithContext(context.Background(), createTemplateAttachmentsOptions)
}

// CreateTemplateAttachmentsWithContext is an alternate form of the CreateTemplateAttachments method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) CreateTemplateAttachmentsWithContext(ctx context.Context, createTemplateAttachmentsOptions *CreateTemplateAttachmentsOptions) (result *CreateTemplateAttachmentsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTemplateAttachmentsOptions, "createTemplateAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTemplateAttachmentsOptions, "createTemplateAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *createTemplateAttachmentsOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createTemplateAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "CreateTemplateAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createTemplateAttachmentsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createTemplateAttachmentsOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createTemplateAttachmentsOptions.Attachments != nil {
		body["attachments"] = createTemplateAttachmentsOptions.Attachments
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateTemplateAttachmentsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListTemplateAttachments : List attachments
// Retrieves a list of scope attachments that are associated with the specified template.
func (configurationGovernance *ConfigurationGovernanceV1) ListTemplateAttachments(listTemplateAttachmentsOptions *ListTemplateAttachmentsOptions) (result *TemplateAttachmentList, response *core.DetailedResponse, err error) {
	return configurationGovernance.ListTemplateAttachmentsWithContext(context.Background(), listTemplateAttachmentsOptions)
}

// ListTemplateAttachmentsWithContext is an alternate form of the ListTemplateAttachments method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) ListTemplateAttachmentsWithContext(ctx context.Context, listTemplateAttachmentsOptions *ListTemplateAttachmentsOptions) (result *TemplateAttachmentList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTemplateAttachmentsOptions, "listTemplateAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTemplateAttachmentsOptions, "listTemplateAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id": *listTemplateAttachmentsOptions.TemplateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listTemplateAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "ListTemplateAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listTemplateAttachmentsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listTemplateAttachmentsOptions.TransactionID))
	}

	if listTemplateAttachmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTemplateAttachmentsOptions.Limit))
	}
	if listTemplateAttachmentsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listTemplateAttachmentsOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAttachmentList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTemplateAttachment : Get an attachment
// Retrieves an existing scope attachment for a template.
func (configurationGovernance *ConfigurationGovernanceV1) GetTemplateAttachment(getTemplateAttachmentOptions *GetTemplateAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	return configurationGovernance.GetTemplateAttachmentWithContext(context.Background(), getTemplateAttachmentOptions)
}

// GetTemplateAttachmentWithContext is an alternate form of the GetTemplateAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) GetTemplateAttachmentWithContext(ctx context.Context, getTemplateAttachmentOptions *GetTemplateAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTemplateAttachmentOptions, "getTemplateAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTemplateAttachmentOptions, "getTemplateAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id":   *getTemplateAttachmentOptions.TemplateID,
		"attachment_id": *getTemplateAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTemplateAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "GetTemplateAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getTemplateAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getTemplateAttachmentOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateTemplateAttachment : Update an attachment
// Updates an existing scope attachment based on the properties that you specify.
func (configurationGovernance *ConfigurationGovernanceV1) UpdateTemplateAttachment(updateTemplateAttachmentOptions *UpdateTemplateAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	return configurationGovernance.UpdateTemplateAttachmentWithContext(context.Background(), updateTemplateAttachmentOptions)
}

// UpdateTemplateAttachmentWithContext is an alternate form of the UpdateTemplateAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) UpdateTemplateAttachmentWithContext(ctx context.Context, updateTemplateAttachmentOptions *UpdateTemplateAttachmentOptions) (result *TemplateAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTemplateAttachmentOptions, "updateTemplateAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateTemplateAttachmentOptions, "updateTemplateAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id":   *updateTemplateAttachmentOptions.TemplateID,
		"attachment_id": *updateTemplateAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateTemplateAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "UpdateTemplateAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateTemplateAttachmentOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateTemplateAttachmentOptions.IfMatch))
	}
	if updateTemplateAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*updateTemplateAttachmentOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if updateTemplateAttachmentOptions.AccountID != nil {
		body["account_id"] = updateTemplateAttachmentOptions.AccountID
	}
	if updateTemplateAttachmentOptions.IncludedScope != nil {
		body["included_scope"] = updateTemplateAttachmentOptions.IncludedScope
	}
	if updateTemplateAttachmentOptions.ExcludedScopes != nil {
		body["excluded_scopes"] = updateTemplateAttachmentOptions.ExcludedScopes
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
	response, err = configurationGovernance.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTemplateAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTemplateAttachment : Delete an attachment
// Deletes an existing scope attachment.
func (configurationGovernance *ConfigurationGovernanceV1) DeleteTemplateAttachment(deleteTemplateAttachmentOptions *DeleteTemplateAttachmentOptions) (response *core.DetailedResponse, err error) {
	return configurationGovernance.DeleteTemplateAttachmentWithContext(context.Background(), deleteTemplateAttachmentOptions)
}

// DeleteTemplateAttachmentWithContext is an alternate form of the DeleteTemplateAttachment method which supports a Context parameter
func (configurationGovernance *ConfigurationGovernanceV1) DeleteTemplateAttachmentWithContext(ctx context.Context, deleteTemplateAttachmentOptions *DeleteTemplateAttachmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTemplateAttachmentOptions, "deleteTemplateAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTemplateAttachmentOptions, "deleteTemplateAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"template_id":   *deleteTemplateAttachmentOptions.TemplateID,
		"attachment_id": *deleteTemplateAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = configurationGovernance.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(configurationGovernance.Service.Options.URL, `/config/v1/templates/{template_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteTemplateAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("configuration_governance", "V1", "DeleteTemplateAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteTemplateAttachmentOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteTemplateAttachmentOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = configurationGovernance.Service.Request(request, nil)

	return
}

// BaseTargetAttribute : The attributes that are associated with a rule or template target.
type BaseTargetAttribute struct {
	// The name of the additional attribute that you want to use to further qualify the target.
	//
	// Options differ depending on the service or resource that you are targeting with a rule or template. For more
	// information, refer to the service documentation.
	Name *string `json:"name" validate:"required"`

	// The value that you want to apply to `name` field.
	//
	// Options differ depending on the rule or template that you configure. For more information, refer to the service
	// documentation.
	Value *string `json:"value" validate:"required"`
}

// NewBaseTargetAttribute : Instantiate BaseTargetAttribute (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewBaseTargetAttribute(name string, value string) (_model *BaseTargetAttribute, err error) {
	_model = &BaseTargetAttribute{
		Name:  core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalBaseTargetAttribute unmarshals an instance of BaseTargetAttribute from the specified map of raw messages.
func UnmarshalBaseTargetAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BaseTargetAttribute)
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

// CreateRuleAttachmentsOptions : The CreateRuleAttachments options.
type CreateRuleAttachmentsOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	Attachments []RuleAttachmentRequest `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRuleAttachmentsOptions : Instantiate CreateRuleAttachmentsOptions
func (*ConfigurationGovernanceV1) NewCreateRuleAttachmentsOptions(ruleID string, attachments []RuleAttachmentRequest) *CreateRuleAttachmentsOptions {
	return &CreateRuleAttachmentsOptions{
		RuleID:      core.StringPtr(ruleID),
		Attachments: attachments,
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *CreateRuleAttachmentsOptions) SetRuleID(ruleID string) *CreateRuleAttachmentsOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *CreateRuleAttachmentsOptions) SetAttachments(attachments []RuleAttachmentRequest) *CreateRuleAttachmentsOptions {
	_options.Attachments = attachments
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateRuleAttachmentsOptions) SetTransactionID(transactionID string) *CreateRuleAttachmentsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRuleAttachmentsOptions) SetHeaders(param map[string]string) *CreateRuleAttachmentsOptions {
	options.Headers = param
	return options
}

// CreateRuleAttachmentsResponse : CreateRuleAttachmentsResponse struct
type CreateRuleAttachmentsResponse struct {
	Attachments []RuleAttachment `json:"attachments" validate:"required"`
}

// UnmarshalCreateRuleAttachmentsResponse unmarshals an instance of CreateRuleAttachmentsResponse from the specified map of raw messages.
func UnmarshalCreateRuleAttachmentsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateRuleAttachmentsResponse)
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalRuleAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRuleRequest : A rule to be created.
type CreateRuleRequest struct {
	// A field that you can use in bulk operations to store a custom identifier for an individual request. If you omit this
	// field, the service generates and sends a `request_id` string for each new rule. The generated string corresponds
	// with the numerical order of the rules request array. For example, `"request_id": "1"`, `"request_id": "2"`.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `request_id` with
	// each request.
	RequestID *string `json:"request_id,omitempty"`

	// Properties that you can associate with a rule.
	Rule *RuleRequest `json:"rule" validate:"required"`
}

// NewCreateRuleRequest : Instantiate CreateRuleRequest (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewCreateRuleRequest(rule *RuleRequest) (_model *CreateRuleRequest, err error) {
	_model = &CreateRuleRequest{
		Rule: rule,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCreateRuleRequest unmarshals an instance of CreateRuleRequest from the specified map of raw messages.
func UnmarshalCreateRuleRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateRuleRequest)
	err = core.UnmarshalPrimitive(m, "request_id", &obj.RequestID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalRuleRequest)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRuleResponse : Response information for a rule request.
//
// If the 'status_code' property indicates success, the 'request_id' and 'rule' properties are returned in the response.
// If the 'status_code' property indicates an error, the 'request_id', 'errors', and 'trace' fields are returned.
type CreateRuleResponse struct {
	// The identifier that is used to correlate an individual request.
	//
	// To assist with debugging, you can use this ID to identify and inspect only one request that was made as part of a
	// bulk operation.
	RequestID *string `json:"request_id,omitempty"`

	// The HTTP response status code.
	StatusCode *int64 `json:"status_code,omitempty"`

	// Information about a newly-created rule.
	//
	// This field is present for successful requests.
	Rule *Rule `json:"rule,omitempty"`

	// The error contents of the multi-status response.
	//
	// This field is present for unsuccessful requests.
	Errors []RuleResponseError `json:"errors,omitempty"`

	// The UUID that uniquely identifies the request.
	//
	// This field is present for unsuccessful requests.
	Trace *string `json:"trace,omitempty"`
}

// UnmarshalCreateRuleResponse unmarshals an instance of CreateRuleResponse from the specified map of raw messages.
func UnmarshalCreateRuleResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateRuleResponse)
	err = core.UnmarshalPrimitive(m, "request_id", &obj.RequestID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rule", &obj.Rule, UnmarshalRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalRuleResponseError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRulesOptions : The CreateRules options.
type CreateRulesOptions struct {
	// A list of rules to be created.
	Rules []CreateRuleRequest `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRulesOptions : Instantiate CreateRulesOptions
func (*ConfigurationGovernanceV1) NewCreateRulesOptions(rules []CreateRuleRequest) *CreateRulesOptions {
	return &CreateRulesOptions{
		Rules: rules,
	}
}

// SetRules : Allow user to set Rules
func (_options *CreateRulesOptions) SetRules(rules []CreateRuleRequest) *CreateRulesOptions {
	_options.Rules = rules
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateRulesOptions) SetTransactionID(transactionID string) *CreateRulesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRulesOptions) SetHeaders(param map[string]string) *CreateRulesOptions {
	options.Headers = param
	return options
}

// CreateRulesResponse : The response associated with a request to create one or more rules.
type CreateRulesResponse struct {
	// An array of rule responses.
	Rules []CreateRuleResponse `json:"rules" validate:"required"`
}

// UnmarshalCreateRulesResponse unmarshals an instance of CreateRulesResponse from the specified map of raw messages.
func UnmarshalCreateRulesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateRulesResponse)
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalCreateRuleResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTemplateAttachmentsOptions : The CreateTemplateAttachments options.
type CreateTemplateAttachmentsOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	Attachments []TemplateAttachmentRequest `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTemplateAttachmentsOptions : Instantiate CreateTemplateAttachmentsOptions
func (*ConfigurationGovernanceV1) NewCreateTemplateAttachmentsOptions(templateID string, attachments []TemplateAttachmentRequest) *CreateTemplateAttachmentsOptions {
	return &CreateTemplateAttachmentsOptions{
		TemplateID:  core.StringPtr(templateID),
		Attachments: attachments,
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *CreateTemplateAttachmentsOptions) SetTemplateID(templateID string) *CreateTemplateAttachmentsOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *CreateTemplateAttachmentsOptions) SetAttachments(attachments []TemplateAttachmentRequest) *CreateTemplateAttachmentsOptions {
	_options.Attachments = attachments
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateTemplateAttachmentsOptions) SetTransactionID(transactionID string) *CreateTemplateAttachmentsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTemplateAttachmentsOptions) SetHeaders(param map[string]string) *CreateTemplateAttachmentsOptions {
	options.Headers = param
	return options
}

// CreateTemplateAttachmentsResponse : CreateTemplateAttachmentsResponse struct
type CreateTemplateAttachmentsResponse struct {
	Attachments []TemplateAttachment `json:"attachments" validate:"required"`
}

// UnmarshalCreateTemplateAttachmentsResponse unmarshals an instance of CreateTemplateAttachmentsResponse from the specified map of raw messages.
func UnmarshalCreateTemplateAttachmentsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTemplateAttachmentsResponse)
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalTemplateAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTemplateRequest : A template to be created.
type CreateTemplateRequest struct {
	// A field that you can use in bulk operations to store a custom identifier for an individual request. If you omit this
	// field, the service generates and sends a `request_id` string for each new template. The generated string corresponds
	// with the numerical order of the templates request array. For example, `"request_id": "1"`,
	// `"request_id": "2"`.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `request_id` with
	// each request.
	RequestID *string `json:"request_id,omitempty"`

	// Properties that you can associate with a template.
	Template *Template `json:"template" validate:"required"`
}

// NewCreateTemplateRequest : Instantiate CreateTemplateRequest (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewCreateTemplateRequest(template *Template) (_model *CreateTemplateRequest, err error) {
	_model = &CreateTemplateRequest{
		Template: template,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCreateTemplateRequest unmarshals an instance of CreateTemplateRequest from the specified map of raw messages.
func UnmarshalCreateTemplateRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTemplateRequest)
	err = core.UnmarshalPrimitive(m, "request_id", &obj.RequestID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTemplateResponse : Response information for a template request.
//
// If the 'status_code' property indicates success, the 'request_id' and
// 'template' properties is returned in the response. If the 'status_code' property indicates an error, the
// 'request_id', 'errors', and 'trace' fields are returned.
type CreateTemplateResponse struct {
	// The identifier that is used to correlate an individual request.
	//
	// To assist with debugging, you can use this ID to identify and inspect only one request that was made as part of a
	// bulk operation.
	RequestID *string `json:"request_id,omitempty"`

	// The HTTP response status code.
	StatusCode *int64 `json:"status_code,omitempty"`

	// Information about a newly-created template.
	//
	// This field is present for successful requests.
	Template *Template `json:"template,omitempty"`

	// The error contents of the multi-status response.
	//
	// This field is present for unsuccessful requests.
	Errors []TemplateResponseError `json:"errors,omitempty"`

	// The UUID that uniquely identifies the request.
	//
	// This field is present for unsuccessful requests.
	Trace *string `json:"trace,omitempty"`
}

// UnmarshalCreateTemplateResponse unmarshals an instance of CreateTemplateResponse from the specified map of raw messages.
func UnmarshalCreateTemplateResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTemplateResponse)
	err = core.UnmarshalPrimitive(m, "request_id", &obj.RequestID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "template", &obj.Template, UnmarshalTemplate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalTemplateResponseError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTemplatesOptions : The CreateTemplates options.
type CreateTemplatesOptions struct {
	// A list of templates to be created.
	Templates []CreateTemplateRequest `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTemplatesOptions : Instantiate CreateTemplatesOptions
func (*ConfigurationGovernanceV1) NewCreateTemplatesOptions(templates []CreateTemplateRequest) *CreateTemplatesOptions {
	return &CreateTemplatesOptions{
		Templates: templates,
	}
}

// SetTemplates : Allow user to set Templates
func (_options *CreateTemplatesOptions) SetTemplates(templates []CreateTemplateRequest) *CreateTemplatesOptions {
	_options.Templates = templates
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *CreateTemplatesOptions) SetTransactionID(transactionID string) *CreateTemplatesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTemplatesOptions) SetHeaders(param map[string]string) *CreateTemplatesOptions {
	options.Headers = param
	return options
}

// CreateTemplatesResponse : The response associated with a request to create one or more templates.
type CreateTemplatesResponse struct {
	// An array of template responses.
	Templates []CreateTemplateResponse `json:"templates" validate:"required"`
}

// UnmarshalCreateTemplatesResponse unmarshals an instance of CreateTemplatesResponse from the specified map of raw messages.
func UnmarshalCreateTemplatesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateTemplatesResponse)
	err = core.UnmarshalModel(m, "templates", &obj.Templates, UnmarshalCreateTemplateResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteRuleAttachmentOptions : The DeleteRuleAttachment options.
type DeleteRuleAttachmentOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRuleAttachmentOptions : Instantiate DeleteRuleAttachmentOptions
func (*ConfigurationGovernanceV1) NewDeleteRuleAttachmentOptions(ruleID string, attachmentID string) *DeleteRuleAttachmentOptions {
	return &DeleteRuleAttachmentOptions{
		RuleID:       core.StringPtr(ruleID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteRuleAttachmentOptions) SetRuleID(ruleID string) *DeleteRuleAttachmentOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *DeleteRuleAttachmentOptions) SetAttachmentID(attachmentID string) *DeleteRuleAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteRuleAttachmentOptions) SetTransactionID(transactionID string) *DeleteRuleAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRuleAttachmentOptions) SetHeaders(param map[string]string) *DeleteRuleAttachmentOptions {
	options.Headers = param
	return options
}

// DeleteRuleOptions : The DeleteRule options.
type DeleteRuleOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRuleOptions : Instantiate DeleteRuleOptions
func (*ConfigurationGovernanceV1) NewDeleteRuleOptions(ruleID string) *DeleteRuleOptions {
	return &DeleteRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteRuleOptions) SetRuleID(ruleID string) *DeleteRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteRuleOptions) SetTransactionID(transactionID string) *DeleteRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRuleOptions) SetHeaders(param map[string]string) *DeleteRuleOptions {
	options.Headers = param
	return options
}

// DeleteTemplateAttachmentOptions : The DeleteTemplateAttachment options.
type DeleteTemplateAttachmentOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteTemplateAttachmentOptions : Instantiate DeleteTemplateAttachmentOptions
func (*ConfigurationGovernanceV1) NewDeleteTemplateAttachmentOptions(templateID string, attachmentID string) *DeleteTemplateAttachmentOptions {
	return &DeleteTemplateAttachmentOptions{
		TemplateID:   core.StringPtr(templateID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteTemplateAttachmentOptions) SetTemplateID(templateID string) *DeleteTemplateAttachmentOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *DeleteTemplateAttachmentOptions) SetAttachmentID(attachmentID string) *DeleteTemplateAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteTemplateAttachmentOptions) SetTransactionID(transactionID string) *DeleteTemplateAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTemplateAttachmentOptions) SetHeaders(param map[string]string) *DeleteTemplateAttachmentOptions {
	options.Headers = param
	return options
}

// DeleteTemplateOptions : The DeleteTemplate options.
type DeleteTemplateOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteTemplateOptions : Instantiate DeleteTemplateOptions
func (*ConfigurationGovernanceV1) NewDeleteTemplateOptions(templateID string) *DeleteTemplateOptions {
	return &DeleteTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *DeleteTemplateOptions) SetTemplateID(templateID string) *DeleteTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *DeleteTemplateOptions) SetTransactionID(transactionID string) *DeleteTemplateOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTemplateOptions) SetHeaders(param map[string]string) *DeleteTemplateOptions {
	options.Headers = param
	return options
}

// EnforcementAction : EnforcementAction struct
type EnforcementAction struct {
	// To block a request from completing, use `disallow`. To log the request to Activity Tracker with LogDNA, use
	// `audit_log`.
	Action *string `json:"action" validate:"required"`
}

// Constants associated with the EnforcementAction.Action property.
// To block a request from completing, use `disallow`. To log the request to Activity Tracker with LogDNA, use
// `audit_log`.
// To block a request from completing, use `disallow`.
const (
	EnforcementActionActionDisallowConst = "audit_log"
)

// NewEnforcementAction : Instantiate EnforcementAction (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewEnforcementAction(action string) (_model *EnforcementAction, err error) {
	_model = &EnforcementAction{
		Action: core.StringPtr(action),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalEnforcementAction unmarshals an instance of EnforcementAction from the specified map of raw messages.
func UnmarshalEnforcementAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnforcementAction)
	err = core.UnmarshalPrimitive(m, "action", &obj.Action)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetRuleAttachmentOptions : The GetRuleAttachment options.
type GetRuleAttachmentOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRuleAttachmentOptions : Instantiate GetRuleAttachmentOptions
func (*ConfigurationGovernanceV1) NewGetRuleAttachmentOptions(ruleID string, attachmentID string) *GetRuleAttachmentOptions {
	return &GetRuleAttachmentOptions{
		RuleID:       core.StringPtr(ruleID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *GetRuleAttachmentOptions) SetRuleID(ruleID string) *GetRuleAttachmentOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *GetRuleAttachmentOptions) SetAttachmentID(attachmentID string) *GetRuleAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetRuleAttachmentOptions) SetTransactionID(transactionID string) *GetRuleAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRuleAttachmentOptions) SetHeaders(param map[string]string) *GetRuleAttachmentOptions {
	options.Headers = param
	return options
}

// GetRuleOptions : The GetRule options.
type GetRuleOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRuleOptions : Instantiate GetRuleOptions
func (*ConfigurationGovernanceV1) NewGetRuleOptions(ruleID string) *GetRuleOptions {
	return &GetRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *GetRuleOptions) SetRuleID(ruleID string) *GetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetRuleOptions) SetTransactionID(transactionID string) *GetRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRuleOptions) SetHeaders(param map[string]string) *GetRuleOptions {
	options.Headers = param
	return options
}

// GetTemplateAttachmentOptions : The GetTemplateAttachment options.
type GetTemplateAttachmentOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTemplateAttachmentOptions : Instantiate GetTemplateAttachmentOptions
func (*ConfigurationGovernanceV1) NewGetTemplateAttachmentOptions(templateID string, attachmentID string) *GetTemplateAttachmentOptions {
	return &GetTemplateAttachmentOptions{
		TemplateID:   core.StringPtr(templateID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetTemplateAttachmentOptions) SetTemplateID(templateID string) *GetTemplateAttachmentOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *GetTemplateAttachmentOptions) SetAttachmentID(attachmentID string) *GetTemplateAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetTemplateAttachmentOptions) SetTransactionID(transactionID string) *GetTemplateAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateAttachmentOptions) SetHeaders(param map[string]string) *GetTemplateAttachmentOptions {
	options.Headers = param
	return options
}

// GetTemplateOptions : The GetTemplate options.
type GetTemplateOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTemplateOptions : Instantiate GetTemplateOptions
func (*ConfigurationGovernanceV1) NewGetTemplateOptions(templateID string) *GetTemplateOptions {
	return &GetTemplateOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *GetTemplateOptions) SetTemplateID(templateID string) *GetTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *GetTemplateOptions) SetTransactionID(transactionID string) *GetTemplateOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTemplateOptions) SetHeaders(param map[string]string) *GetTemplateOptions {
	options.Headers = param
	return options
}

// Link : A link that is used to paginate through available resources.
type Link struct {
	// The URL for the first, previous, next, or last page of resources.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalLink unmarshals an instance of Link from the specified map of raw messages.
func UnmarshalLink(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Link)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListRuleAttachmentsOptions : The ListRuleAttachments options.
type ListRuleAttachmentsOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// The number of resources to retrieve. By default, list operations return the first 100 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 rules, and you want to retrieve only the first 5 rules, use
	// `../rules?account_id={account_id}&limit=5`.
	Limit *int64

	// The number of resources to skip. By specifying `offset`, you retrieve a subset of resources that starts with the
	// `offset` value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 rules, and you want to retrieve rules 26 through 50, use
	// `../rules?account_id={account_id}&offset=25&limit=5`.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRuleAttachmentsOptions : Instantiate ListRuleAttachmentsOptions
func (*ConfigurationGovernanceV1) NewListRuleAttachmentsOptions(ruleID string) *ListRuleAttachmentsOptions {
	return &ListRuleAttachmentsOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *ListRuleAttachmentsOptions) SetRuleID(ruleID string) *ListRuleAttachmentsOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListRuleAttachmentsOptions) SetTransactionID(transactionID string) *ListRuleAttachmentsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListRuleAttachmentsOptions) SetLimit(limit int64) *ListRuleAttachmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListRuleAttachmentsOptions) SetOffset(offset int64) *ListRuleAttachmentsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRuleAttachmentsOptions) SetHeaders(param map[string]string) *ListRuleAttachmentsOptions {
	options.Headers = param
	return options
}

// ListRulesOptions : The ListRules options.
type ListRulesOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Retrieves a list of rules that have scope attachments.
	Attached *bool

	// Retrieves a list of rules that match the labels that you specify.
	Labels *string

	// Retrieves a list of rules that match the scope ID that you specify.
	Scopes *string

	// The number of resources to retrieve. By default, list operations return the first 100 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 rules, and you want to retrieve only the first 5 rules, use
	// `../rules?account_id={account_id}&limit=5`.
	Limit *int64

	// The number of resources to skip. By specifying `offset`, you retrieve a subset of resources that starts with the
	// `offset` value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 rules, and you want to retrieve rules 26 through 50, use
	// `../rules?account_id={account_id}&offset=25&limit=5`.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRulesOptions : Instantiate ListRulesOptions
func (*ConfigurationGovernanceV1) NewListRulesOptions(accountID string) *ListRulesOptions {
	return &ListRulesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListRulesOptions) SetAccountID(accountID string) *ListRulesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListRulesOptions) SetTransactionID(transactionID string) *ListRulesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetAttached : Allow user to set Attached
func (_options *ListRulesOptions) SetAttached(attached bool) *ListRulesOptions {
	_options.Attached = core.BoolPtr(attached)
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *ListRulesOptions) SetLabels(labels string) *ListRulesOptions {
	_options.Labels = core.StringPtr(labels)
	return _options
}

// SetScopes : Allow user to set Scopes
func (_options *ListRulesOptions) SetScopes(scopes string) *ListRulesOptions {
	_options.Scopes = core.StringPtr(scopes)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListRulesOptions) SetLimit(limit int64) *ListRulesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListRulesOptions) SetOffset(offset int64) *ListRulesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRulesOptions) SetHeaders(param map[string]string) *ListRulesOptions {
	options.Headers = param
	return options
}

// ListTemplateAttachmentsOptions : The ListTemplateAttachments options.
type ListTemplateAttachmentsOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// The number of resources to retrieve. By default, list operations return the first 100 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 rules, and you want to retrieve only the first 5 rules, use
	// `../rules?account_id={account_id}&limit=5`.
	Limit *int64

	// The number of resources to skip. By specifying `offset`, you retrieve a subset of resources that starts with the
	// `offset` value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 rules, and you want to retrieve rules 26 through 50, use
	// `../rules?account_id={account_id}&offset=25&limit=5`.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListTemplateAttachmentsOptions : Instantiate ListTemplateAttachmentsOptions
func (*ConfigurationGovernanceV1) NewListTemplateAttachmentsOptions(templateID string) *ListTemplateAttachmentsOptions {
	return &ListTemplateAttachmentsOptions{
		TemplateID: core.StringPtr(templateID),
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *ListTemplateAttachmentsOptions) SetTemplateID(templateID string) *ListTemplateAttachmentsOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListTemplateAttachmentsOptions) SetTransactionID(transactionID string) *ListTemplateAttachmentsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListTemplateAttachmentsOptions) SetLimit(limit int64) *ListTemplateAttachmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListTemplateAttachmentsOptions) SetOffset(offset int64) *ListTemplateAttachmentsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTemplateAttachmentsOptions) SetHeaders(param map[string]string) *ListTemplateAttachmentsOptions {
	options.Headers = param
	return options
}

// ListTemplatesOptions : The ListTemplates options.
type ListTemplatesOptions struct {
	// Your IBM Cloud account ID.
	AccountID *string `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Retrieves a list of templates that have scope attachments.
	Attached *bool

	// Retrieves a list of templates that match the scope ID that you specify.
	Scopes *string

	// The number of resources to retrieve. By default, list operations return the first 100 items. To retrieve a different
	// set of items, use `limit` with `offset` to page through your available resources.
	//
	// **Usage:** If you have 20 rules, and you want to retrieve only the first 5 rules, use
	// `../rules?account_id={account_id}&limit=5`.
	Limit *int64

	// The number of resources to skip. By specifying `offset`, you retrieve a subset of resources that starts with the
	// `offset` value. Use `offset` with `limit` to page through your available resources.
	//
	// **Usage:** If you have 100 rules, and you want to retrieve rules 26 through 50, use
	// `../rules?account_id={account_id}&offset=25&limit=5`.
	Offset *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListTemplatesOptions : Instantiate ListTemplatesOptions
func (*ConfigurationGovernanceV1) NewListTemplatesOptions(accountID string) *ListTemplatesOptions {
	return &ListTemplatesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListTemplatesOptions) SetAccountID(accountID string) *ListTemplatesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *ListTemplatesOptions) SetTransactionID(transactionID string) *ListTemplatesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetAttached : Allow user to set Attached
func (_options *ListTemplatesOptions) SetAttached(attached bool) *ListTemplatesOptions {
	_options.Attached = core.BoolPtr(attached)
	return _options
}

// SetScopes : Allow user to set Scopes
func (_options *ListTemplatesOptions) SetScopes(scopes string) *ListTemplatesOptions {
	_options.Scopes = core.StringPtr(scopes)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListTemplatesOptions) SetLimit(limit int64) *ListTemplatesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListTemplatesOptions) SetOffset(offset int64) *ListTemplatesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTemplatesOptions) SetHeaders(param map[string]string) *ListTemplatesOptions {
	options.Headers = param
	return options
}

// Rule : Properties associated with a rule, including both user-defined and server-populated properties.
type Rule struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id,omitempty"`

	// A human-readable alias to assign to your rule.
	Name *string `json:"name" validate:"required"`

	// An extended description of your rule.
	Description *string `json:"description" validate:"required"`

	// The type of rule. Rules that you create are `user_defined`.
	RuleType *string `json:"rule_type,omitempty"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *TargetResource `json:"target" validate:"required"`

	RequiredConfig RuleRequiredConfigIntf `json:"required_config" validate:"required"`

	// The actions that the service must run on your behalf when a request to create or modify the target resource does not
	// comply with your conditions.
	EnforcementActions []EnforcementAction `json:"enforcement_actions" validate:"required"`

	// Labels that you can use to group and search for similar rules, such as those that help you to meet a specific
	// organization guideline.
	Labels []string `json:"labels,omitempty"`

	// The UUID that uniquely identifies the rule.
	RuleID *string `json:"rule_id,omitempty"`

	// The date the resource was created.
	CreationDate *strfmt.DateTime `json:"creation_date,omitempty"`

	// The unique identifier for the user or application that created the resource.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date the resource was last modified.
	ModificationDate *strfmt.DateTime `json:"modification_date,omitempty"`

	// The unique identifier for the user or application that last modified the resource.
	ModifiedBy *string `json:"modified_by,omitempty"`

	// The number of scope attachments that are associated with the rule.
	NumberOfAttachments *int64 `json:"number_of_attachments,omitempty"`
}

// Constants associated with the Rule.RuleType property.
// The type of rule. Rules that you create are `user_defined`.
const (
	RuleRuleTypeUserDefinedConst = "user_defined"
)

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "rule_type", &obj.RuleType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTargetResource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_config", &obj.RequiredConfig, UnmarshalRuleRequiredConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "enforcement_actions", &obj.EnforcementActions, UnmarshalEnforcementAction)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rule_id", &obj.RuleID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modification_date", &obj.ModificationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_by", &obj.ModifiedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "number_of_attachments", &obj.NumberOfAttachments)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleAttachment : The scopes to attach to a rule.
type RuleAttachment struct {
	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `json:"attachment_id" validate:"required"`

	// The UUID that uniquely identifies the rule.
	RuleID *string `json:"rule_id" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The extent at which the rule can be attached across your accounts.
	IncludedScope *RuleScope `json:"included_scope" validate:"required"`

	// The extent at which the rule can be excluded from the included scope.
	ExcludedScopes []RuleScope `json:"excluded_scopes,omitempty"`
}

// UnmarshalRuleAttachment unmarshals an instance of RuleAttachment from the specified map of raw messages.
func UnmarshalRuleAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleAttachment)
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rule_id", &obj.RuleID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "included_scope", &obj.IncludedScope, UnmarshalRuleScope)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "excluded_scopes", &obj.ExcludedScopes, UnmarshalRuleScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleAttachmentList : A list of attachments.
type RuleAttachmentList struct {
	// The requested offset for the returned items.
	Offset *int64 `json:"offset" validate:"required"`

	// The requested limit for the returned items.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of available items.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The first page of available items.
	First *Link `json:"first" validate:"required"`

	// The last page of available items.
	Last *Link `json:"last" validate:"required"`

	Attachments []RuleAttachment `json:"attachments" validate:"required"`
}

// UnmarshalRuleAttachmentList unmarshals an instance of RuleAttachmentList from the specified map of raw messages.
func UnmarshalRuleAttachmentList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleAttachmentList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalRuleAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleAttachmentRequest : The scopes to attach to a rule.
type RuleAttachmentRequest struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The extent at which the rule can be attached across your accounts.
	IncludedScope *RuleScope `json:"included_scope" validate:"required"`

	// The extent at which the rule can be excluded from the included scope.
	ExcludedScopes []RuleScope `json:"excluded_scopes,omitempty"`
}

// NewRuleAttachmentRequest : Instantiate RuleAttachmentRequest (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleAttachmentRequest(accountID string, includedScope *RuleScope) (_model *RuleAttachmentRequest, err error) {
	_model = &RuleAttachmentRequest{
		AccountID:     core.StringPtr(accountID),
		IncludedScope: includedScope,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleAttachmentRequest unmarshals an instance of RuleAttachmentRequest from the specified map of raw messages.
func UnmarshalRuleAttachmentRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleAttachmentRequest)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "included_scope", &obj.IncludedScope, UnmarshalRuleScope)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "excluded_scopes", &obj.ExcludedScopes, UnmarshalRuleScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleCondition : RuleCondition struct
// Models which "extend" this model:
// - RuleConditionSingleProperty
// - RuleConditionOrLvl2
// - RuleConditionAndLvl2
type RuleCondition struct {
	Description *string `json:"description,omitempty"`

	// A resource configuration variable that describes the property that you want to apply to the target resource.
	//
	// Available options depend on the target service and resource.
	Property *string `json:"property,omitempty"`

	// The way in which the `property` field is compared to its value.
	//
	// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
	Operator *string `json:"operator,omitempty"`

	// The way in which you want your property to be applied.
	//
	// Value options differ depending on the rule that you configure. If you use a boolean operator, you do not need to
	// input a value.
	Value *string `json:"value,omitempty"`

	Or []RuleSingleProperty `json:"or,omitempty"`

	And []RuleSingleProperty `json:"and,omitempty"`
}

// Constants associated with the RuleCondition.Operator property.
// The way in which the `property` field is compared to its value.
//
// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
const (
	RuleConditionOperatorIpsInRangeConst           = "ips_in_range"
	RuleConditionOperatorIsEmptyConst              = "is_empty"
	RuleConditionOperatorIsFalseConst              = "is_false"
	RuleConditionOperatorIsNotEmptyConst           = "is_not_empty"
	RuleConditionOperatorIsTrueConst               = "is_true"
	RuleConditionOperatorNumEqualsConst            = "num_equals"
	RuleConditionOperatorNumGreaterThanConst       = "num_greater_than"
	RuleConditionOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RuleConditionOperatorNumLessThanConst          = "num_less_than"
	RuleConditionOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RuleConditionOperatorNumNotEqualsConst         = "num_not_equals"
	RuleConditionOperatorStringEqualsConst         = "string_equals"
	RuleConditionOperatorStringMatchConst          = "string_match"
	RuleConditionOperatorStringNotEqualsConst      = "string_not_equals"
	RuleConditionOperatorStringNotMatchConst       = "string_not_match"
	RuleConditionOperatorStringsInListConst        = "strings_in_list"
)

func (*RuleCondition) isaRuleCondition() bool {
	return true
}

type RuleConditionIntf interface {
	isaRuleCondition() bool
}

// UnmarshalRuleCondition unmarshals an instance of RuleCondition from the specified map of raw messages.
func UnmarshalRuleCondition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleCondition)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRuleSingleProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRuleSingleProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleList : A list of rules.
type RuleList struct {
	// The requested offset for the returned items.
	Offset *int64 `json:"offset" validate:"required"`

	// The requested limit for the returned items.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of available items.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The first page of available items.
	First *Link `json:"first" validate:"required"`

	// The last page of available items.
	Last *Link `json:"last" validate:"required"`

	// An array of rules.
	Rules []Rule `json:"rules" validate:"required"`
}

// UnmarshalRuleList unmarshals an instance of RuleList from the specified map of raw messages.
func UnmarshalRuleList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleRequest : Properties that you can associate with a rule.
type RuleRequest struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id,omitempty"`

	// A human-readable alias to assign to your rule.
	Name *string `json:"name" validate:"required"`

	// An extended description of your rule.
	Description *string `json:"description" validate:"required"`

	// The type of rule. Rules that you create are `user_defined`.
	RuleType *string `json:"rule_type,omitempty"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *TargetResource `json:"target" validate:"required"`

	RequiredConfig RuleRequiredConfigIntf `json:"required_config" validate:"required"`

	// The actions that the service must run on your behalf when a request to create or modify the target resource does not
	// comply with your conditions.
	EnforcementActions []EnforcementAction `json:"enforcement_actions" validate:"required"`

	// Labels that you can use to group and search for similar rules, such as those that help you to meet a specific
	// organization guideline.
	Labels []string `json:"labels,omitempty"`
}

// Constants associated with the RuleRequest.RuleType property.
// The type of rule. Rules that you create are `user_defined`.
const (
	RuleRequestRuleTypeUserDefinedConst = "user_defined"
)

// NewRuleRequest : Instantiate RuleRequest (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleRequest(name string, description string, target *TargetResource, requiredConfig RuleRequiredConfigIntf, enforcementActions []EnforcementAction) (_model *RuleRequest, err error) {
	_model = &RuleRequest{
		Name:               core.StringPtr(name),
		Description:        core.StringPtr(description),
		Target:             target,
		RequiredConfig:     requiredConfig,
		EnforcementActions: enforcementActions,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleRequest unmarshals an instance of RuleRequest from the specified map of raw messages.
func UnmarshalRuleRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequest)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "rule_type", &obj.RuleType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTargetResource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_config", &obj.RequiredConfig, UnmarshalRuleRequiredConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "enforcement_actions", &obj.EnforcementActions, UnmarshalEnforcementAction)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "labels", &obj.Labels)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleRequiredConfig : RuleRequiredConfig struct
// Models which "extend" this model:
// - RuleRequiredConfigSingleProperty
// - RuleRequiredConfigMultipleProperties
type RuleRequiredConfig struct {
	Description *string `json:"description,omitempty"`

	// A resource configuration variable that describes the property that you want to apply to the target resource.
	//
	// Available options depend on the target service and resource.
	Property *string `json:"property,omitempty"`

	// The way in which the `property` field is compared to its value.
	//
	// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
	Operator *string `json:"operator,omitempty"`

	// The way in which you want your property to be applied.
	//
	// Value options differ depending on the rule that you configure. If you use a boolean operator, you do not need to
	// input a value.
	Value *string `json:"value,omitempty"`

	Or []RuleConditionIntf `json:"or,omitempty"`

	And []RuleConditionIntf `json:"and,omitempty"`
}

// Constants associated with the RuleRequiredConfig.Operator property.
// The way in which the `property` field is compared to its value.
//
// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
const (
	RuleRequiredConfigOperatorIpsInRangeConst           = "ips_in_range"
	RuleRequiredConfigOperatorIsEmptyConst              = "is_empty"
	RuleRequiredConfigOperatorIsFalseConst              = "is_false"
	RuleRequiredConfigOperatorIsNotEmptyConst           = "is_not_empty"
	RuleRequiredConfigOperatorIsTrueConst               = "is_true"
	RuleRequiredConfigOperatorNumEqualsConst            = "num_equals"
	RuleRequiredConfigOperatorNumGreaterThanConst       = "num_greater_than"
	RuleRequiredConfigOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RuleRequiredConfigOperatorNumLessThanConst          = "num_less_than"
	RuleRequiredConfigOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RuleRequiredConfigOperatorNumNotEqualsConst         = "num_not_equals"
	RuleRequiredConfigOperatorStringEqualsConst         = "string_equals"
	RuleRequiredConfigOperatorStringMatchConst          = "string_match"
	RuleRequiredConfigOperatorStringNotEqualsConst      = "string_not_equals"
	RuleRequiredConfigOperatorStringNotMatchConst       = "string_not_match"
	RuleRequiredConfigOperatorStringsInListConst        = "strings_in_list"
)

func (*RuleRequiredConfig) isaRuleRequiredConfig() bool {
	return true
}

type RuleRequiredConfigIntf interface {
	isaRuleRequiredConfig() bool
}

// UnmarshalRuleRequiredConfig unmarshals an instance of RuleRequiredConfig from the specified map of raw messages.
func UnmarshalRuleRequiredConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequiredConfig)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleResponseError : RuleResponseError struct
type RuleResponseError struct {
	// Specifies the problem that caused the error.
	Code *string `json:"code" validate:"required"`

	// Describes the problem.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalRuleResponseError unmarshals an instance of RuleResponseError from the specified map of raw messages.
func UnmarshalRuleResponseError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleResponseError)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleScope : The extent at which the rule can be attached across your accounts.
type RuleScope struct {
	// A short description or alias to assign to the scope.
	Note *string `json:"note,omitempty"`

	// The ID of the scope, such as an enterprise, account, or account group, that you want to evaluate.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The type of scope that you want to evaluate.
	ScopeType *string `json:"scope_type" validate:"required"`
}

// Constants associated with the RuleScope.ScopeType property.
// The type of scope that you want to evaluate.
const (
	RuleScopeScopeTypeAccountConst                = "account"
	RuleScopeScopeTypeAccountResourceGroupConst   = "account.resource_group"
	RuleScopeScopeTypeEnterpriseConst             = "enterprise"
	RuleScopeScopeTypeEnterpriseAccountConst      = "enterprise.account"
	RuleScopeScopeTypeEnterpriseAccountGroupConst = "enterprise.account_group"
)

// NewRuleScope : Instantiate RuleScope (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleScope(scopeID string, scopeType string) (_model *RuleScope, err error) {
	_model = &RuleScope{
		ScopeID:   core.StringPtr(scopeID),
		ScopeType: core.StringPtr(scopeType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleScope unmarshals an instance of RuleScope from the specified map of raw messages.
func UnmarshalRuleScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleScope)
	err = core.UnmarshalPrimitive(m, "note", &obj.Note)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_type", &obj.ScopeType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleSingleProperty : The requirement that must be met to determine the resource's level of compliance in accordance with the rule.
//
// To apply a single property check, define a configuration property and the desired value that you want to check
// against.
type RuleSingleProperty struct {
	Description *string `json:"description,omitempty"`

	// A resource configuration variable that describes the property that you want to apply to the target resource.
	//
	// Available options depend on the target service and resource.
	Property *string `json:"property" validate:"required"`

	// The way in which the `property` field is compared to its value.
	//
	// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
	Operator *string `json:"operator" validate:"required"`

	// The way in which you want your property to be applied.
	//
	// Value options differ depending on the rule that you configure. If you use a boolean operator, you do not need to
	// input a value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the RuleSingleProperty.Operator property.
// The way in which the `property` field is compared to its value.
//
// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
const (
	RuleSinglePropertyOperatorIpsInRangeConst           = "ips_in_range"
	RuleSinglePropertyOperatorIsEmptyConst              = "is_empty"
	RuleSinglePropertyOperatorIsFalseConst              = "is_false"
	RuleSinglePropertyOperatorIsNotEmptyConst           = "is_not_empty"
	RuleSinglePropertyOperatorIsTrueConst               = "is_true"
	RuleSinglePropertyOperatorNumEqualsConst            = "num_equals"
	RuleSinglePropertyOperatorNumGreaterThanConst       = "num_greater_than"
	RuleSinglePropertyOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RuleSinglePropertyOperatorNumLessThanConst          = "num_less_than"
	RuleSinglePropertyOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RuleSinglePropertyOperatorNumNotEqualsConst         = "num_not_equals"
	RuleSinglePropertyOperatorStringEqualsConst         = "string_equals"
	RuleSinglePropertyOperatorStringMatchConst          = "string_match"
	RuleSinglePropertyOperatorStringNotEqualsConst      = "string_not_equals"
	RuleSinglePropertyOperatorStringNotMatchConst       = "string_not_match"
	RuleSinglePropertyOperatorStringsInListConst        = "strings_in_list"
)

// NewRuleSingleProperty : Instantiate RuleSingleProperty (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleSingleProperty(property string, operator string) (_model *RuleSingleProperty, err error) {
	_model = &RuleSingleProperty{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleSingleProperty unmarshals an instance of RuleSingleProperty from the specified map of raw messages.
func UnmarshalRuleSingleProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleSingleProperty)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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

// SimpleTargetResource : The properties that describe the resource that you want to target with the rule or template.
type SimpleTargetResource struct {
	// The programmatic name of the IBM Cloud service that you want to target with the rule or template.
	ServiceName *string `json:"service_name" validate:"required"`

	// The type of resource that you want to target.
	ResourceKind *string `json:"resource_kind" validate:"required"`

	// An extra qualifier for the resource kind. When you include additional attributes, only the resources that match the
	// definition are included in the rule or template.
	AdditionalTargetAttributes []BaseTargetAttribute `json:"additional_target_attributes,omitempty"`
}

// NewSimpleTargetResource : Instantiate SimpleTargetResource (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewSimpleTargetResource(serviceName string, resourceKind string) (_model *SimpleTargetResource, err error) {
	_model = &SimpleTargetResource{
		ServiceName:  core.StringPtr(serviceName),
		ResourceKind: core.StringPtr(resourceKind),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalSimpleTargetResource unmarshals an instance of SimpleTargetResource from the specified map of raw messages.
func UnmarshalSimpleTargetResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SimpleTargetResource)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalBaseTargetAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetResource : The properties that describe the resource that you want to target with the rule or template.
type TargetResource struct {
	// The programmatic name of the IBM Cloud service that you want to target with the rule or template.
	ServiceName *string `json:"service_name" validate:"required"`

	// The type of resource that you want to target.
	ResourceKind *string `json:"resource_kind" validate:"required"`

	// An extra qualifier for the resource kind. When you include additional attributes, only the resources that match the
	// definition are included in the rule or template.
	AdditionalTargetAttributes []TargetResourceAdditionalTargetAttributesItem `json:"additional_target_attributes,omitempty"`
}

// NewTargetResource : Instantiate TargetResource (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTargetResource(serviceName string, resourceKind string) (_model *TargetResource, err error) {
	_model = &TargetResource{
		ServiceName:  core.StringPtr(serviceName),
		ResourceKind: core.StringPtr(resourceKind),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTargetResource unmarshals an instance of TargetResource from the specified map of raw messages.
func UnmarshalTargetResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetResource)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalTargetResourceAdditionalTargetAttributesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetResourceAdditionalTargetAttributesItem : The attributes that are associated with a rule or template target.
type TargetResourceAdditionalTargetAttributesItem struct {
	// The name of the additional attribute that you want to use to further qualify the target.
	//
	// Options differ depending on the service or resource that you are targeting with a rule or template. For more
	// information, refer to the service documentation.
	Name *string `json:"name" validate:"required"`

	// The value that you want to apply to `name` field.
	//
	// Options differ depending on the rule or template that you configure. For more information, refer to the service
	// documentation.
	Value *string `json:"value" validate:"required"`

	// The way in which the `name` field is compared to its value.
	//
	// There are three types of operators: string, numeric, and boolean.
	Operator *string `json:"operator" validate:"required"`
}

// Constants associated with the TargetResourceAdditionalTargetAttributesItem.Operator property.
// The way in which the `name` field is compared to its value.
//
// There are three types of operators: string, numeric, and boolean.
const (
	TargetResourceAdditionalTargetAttributesItemOperatorIpsInRangeConst           = "ips_in_range"
	TargetResourceAdditionalTargetAttributesItemOperatorIsEmptyConst              = "is_empty"
	TargetResourceAdditionalTargetAttributesItemOperatorIsFalseConst              = "is_false"
	TargetResourceAdditionalTargetAttributesItemOperatorIsNotEmptyConst           = "is_not_empty"
	TargetResourceAdditionalTargetAttributesItemOperatorIsTrueConst               = "is_true"
	TargetResourceAdditionalTargetAttributesItemOperatorNumEqualsConst            = "num_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorNumGreaterThanConst       = "num_greater_than"
	TargetResourceAdditionalTargetAttributesItemOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorNumLessThanConst          = "num_less_than"
	TargetResourceAdditionalTargetAttributesItemOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorNumNotEqualsConst         = "num_not_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorStringEqualsConst         = "string_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorStringMatchConst          = "string_match"
	TargetResourceAdditionalTargetAttributesItemOperatorStringNotEqualsConst      = "string_not_equals"
	TargetResourceAdditionalTargetAttributesItemOperatorStringNotMatchConst       = "string_not_match"
	TargetResourceAdditionalTargetAttributesItemOperatorStringsInListConst        = "strings_in_list"
)

// NewTargetResourceAdditionalTargetAttributesItem : Instantiate TargetResourceAdditionalTargetAttributesItem (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTargetResourceAdditionalTargetAttributesItem(name string, value string, operator string) (_model *TargetResourceAdditionalTargetAttributesItem, err error) {
	_model = &TargetResourceAdditionalTargetAttributesItem{
		Name:     core.StringPtr(name),
		Value:    core.StringPtr(value),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTargetResourceAdditionalTargetAttributesItem unmarshals an instance of TargetResourceAdditionalTargetAttributesItem from the specified map of raw messages.
func UnmarshalTargetResourceAdditionalTargetAttributesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetResourceAdditionalTargetAttributesItem)
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

// Template : Properties that you can associate with a template.
type Template struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// A human-readablse alias to assign to your template.
	Name *string `json:"name" validate:"required"`

	// An extended description of your template.
	Description *string `json:"description" validate:"required"`

	// The UUID that uniquely identifies the template.
	TemplateID *string `json:"template_id,omitempty"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *SimpleTargetResource `json:"target" validate:"required"`

	// A list of default property values to apply to your template.
	CustomizedDefaults []TemplateCustomizedDefaultProperty `json:"customized_defaults" validate:"required"`
}

// NewTemplate : Instantiate Template (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTemplate(accountID string, name string, description string, target *SimpleTargetResource, customizedDefaults []TemplateCustomizedDefaultProperty) (_model *Template, err error) {
	_model = &Template{
		AccountID:          core.StringPtr(accountID),
		Name:               core.StringPtr(name),
		Description:        core.StringPtr(description),
		Target:             target,
		CustomizedDefaults: customizedDefaults,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTemplate unmarshals an instance of Template from the specified map of raw messages.
func UnmarshalTemplate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Template)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalSimpleTargetResource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "customized_defaults", &obj.CustomizedDefaults, UnmarshalTemplateCustomizedDefaultProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAttachment : The scopes to attach to a template.
type TemplateAttachment struct {
	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `json:"attachment_id" validate:"required"`

	// The UUID that uniquely identifies the template.
	TemplateID *string `json:"template_id" validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The extent at which the template can be attached across your accounts.
	IncludedScope *TemplateScope `json:"included_scope" validate:"required"`

	ExcludedScopes []TemplateScope `json:"excluded_scopes,omitempty"`
}

// UnmarshalTemplateAttachment unmarshals an instance of TemplateAttachment from the specified map of raw messages.
func UnmarshalTemplateAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAttachment)
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "included_scope", &obj.IncludedScope, UnmarshalTemplateScope)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "excluded_scopes", &obj.ExcludedScopes, UnmarshalTemplateScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAttachmentList : A list of attachments.
type TemplateAttachmentList struct {
	// The requested offset for the returned items.
	Offset *int64 `json:"offset" validate:"required"`

	// The requested limit for the returned items.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of available items.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The first page of available items.
	First *Link `json:"first" validate:"required"`

	// The last page of available items.
	Last *Link `json:"last" validate:"required"`

	Attachments []TemplateAttachment `json:"attachments" validate:"required"`
}

// UnmarshalTemplateAttachmentList unmarshals an instance of TemplateAttachmentList from the specified map of raw messages.
func UnmarshalTemplateAttachmentList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAttachmentList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalTemplateAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateAttachmentRequest : The scopes to attach to a template.
type TemplateAttachmentRequest struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The extent at which the template can be attached across your accounts.
	IncludedScope *TemplateScope `json:"included_scope" validate:"required"`

	ExcludedScopes []TemplateScope `json:"excluded_scopes,omitempty"`
}

// NewTemplateAttachmentRequest : Instantiate TemplateAttachmentRequest (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTemplateAttachmentRequest(accountID string, includedScope *TemplateScope) (_model *TemplateAttachmentRequest, err error) {
	_model = &TemplateAttachmentRequest{
		AccountID:     core.StringPtr(accountID),
		IncludedScope: includedScope,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTemplateAttachmentRequest unmarshals an instance of TemplateAttachmentRequest from the specified map of raw messages.
func UnmarshalTemplateAttachmentRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateAttachmentRequest)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "included_scope", &obj.IncludedScope, UnmarshalTemplateScope)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "excluded_scopes", &obj.ExcludedScopes, UnmarshalTemplateScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateCustomizedDefaultProperty : The default property values to customize.
type TemplateCustomizedDefaultProperty struct {
	// The name of the resource property that you want to configure.
	//
	// Property options differ depending on the service or resource that you are targeting with a template. To view a list
	// of properties that are compatible with templates, refer to the service documentation.
	Property *string `json:"property" validate:"required"`

	// The custom value that you want to apply as the default for the resource property in the `name` field.
	//
	// This value is used to to override the default value that is provided by IBM when a resource is created. Value
	// options differ depending on the resource that you are configuring. To learn more about your options, refer to the
	// service documentation.
	Value *string `json:"value" validate:"required"`
}

// NewTemplateCustomizedDefaultProperty : Instantiate TemplateCustomizedDefaultProperty (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTemplateCustomizedDefaultProperty(property string, value string) (_model *TemplateCustomizedDefaultProperty, err error) {
	_model = &TemplateCustomizedDefaultProperty{
		Property: core.StringPtr(property),
		Value:    core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTemplateCustomizedDefaultProperty unmarshals an instance of TemplateCustomizedDefaultProperty from the specified map of raw messages.
func UnmarshalTemplateCustomizedDefaultProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateCustomizedDefaultProperty)
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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

// TemplateList : A list of templates.
type TemplateList struct {
	// The requested offset for the returned items.
	Offset *int64 `json:"offset" validate:"required"`

	// The requested limit for the returned items.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of available items.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The first page of available items.
	First *Link `json:"first" validate:"required"`

	// The last page of available items.
	Last *Link `json:"last" validate:"required"`

	// An array of templates.
	Templates []Template `json:"templates" validate:"required"`
}

// UnmarshalTemplateList unmarshals an instance of TemplateList from the specified map of raw messages.
func UnmarshalTemplateList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateList)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalLink)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "templates", &obj.Templates, UnmarshalTemplate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateResponse : Properties associated with a template, including both user-defined and server-populated properties.
type TemplateResponse struct {
	// Your IBM Cloud account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// A human-readablse alias to assign to your template.
	Name *string `json:"name" validate:"required"`

	// An extended description of your template.
	Description *string `json:"description" validate:"required"`

	// The UUID that uniquely identifies the template.
	TemplateID *string `json:"template_id,omitempty"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *SimpleTargetResource `json:"target" validate:"required"`

	// A list of default property values to apply to your template.
	CustomizedDefaults []TemplateCustomizedDefaultProperty `json:"customized_defaults" validate:"required"`

	// The date the resource was created.
	CreationDate *strfmt.DateTime `json:"creation_date" validate:"required"`

	// The unique identifier for the user or application that created the resource.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date the resource was last modified.
	ModificationDate *strfmt.DateTime `json:"modification_date" validate:"required"`

	// The unique identifier for the user or application that last modified the resource.
	ModifiedBy *string `json:"modified_by,omitempty"`
}

// UnmarshalTemplateResponse unmarshals an instance of TemplateResponse from the specified map of raw messages.
func UnmarshalTemplateResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateResponse)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "template_id", &obj.TemplateID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalSimpleTargetResource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "customized_defaults", &obj.CustomizedDefaults, UnmarshalTemplateCustomizedDefaultProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "creation_date", &obj.CreationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modification_date", &obj.ModificationDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_by", &obj.ModifiedBy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateResponseError : TemplateResponseError struct
type TemplateResponseError struct {
	// Specifies the problem that caused the error.
	Code *string `json:"code" validate:"required"`

	// Describes the problem.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalTemplateResponseError unmarshals an instance of TemplateResponseError from the specified map of raw messages.
func UnmarshalTemplateResponseError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateResponseError)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TemplateScope : The extent at which the template can be attached across your accounts.
type TemplateScope struct {
	// A short description or alias to assign to the scope.
	Note *string `json:"note,omitempty"`

	// The ID of the scope, such as an enterprise, account, or account group, where you want to apply the customized
	// defaults that are associated with a template.
	ScopeID *string `json:"scope_id" validate:"required"`

	// The type of scope.
	ScopeType *string `json:"scope_type" validate:"required"`
}

// Constants associated with the TemplateScope.ScopeType property.
// The type of scope.
const (
	TemplateScopeScopeTypeAccountConst                = "account"
	TemplateScopeScopeTypeAccountResourceGroupConst   = "account.resource_group"
	TemplateScopeScopeTypeEnterpriseConst             = "enterprise"
	TemplateScopeScopeTypeEnterpriseAccountConst      = "enterprise.account"
	TemplateScopeScopeTypeEnterpriseAccountGroupConst = "enterprise.account_group"
)

// NewTemplateScope : Instantiate TemplateScope (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewTemplateScope(scopeID string, scopeType string) (_model *TemplateScope, err error) {
	_model = &TemplateScope{
		ScopeID:   core.StringPtr(scopeID),
		ScopeType: core.StringPtr(scopeType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTemplateScope unmarshals an instance of TemplateScope from the specified map of raw messages.
func UnmarshalTemplateScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TemplateScope)
	err = core.UnmarshalPrimitive(m, "note", &obj.Note)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_type", &obj.ScopeType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateRuleAttachmentOptions : The UpdateRuleAttachment options.
type UpdateRuleAttachmentOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// Compares a supplied `Etag` value with the version that is stored for the requested resource. If the values match,
	// the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `validate:"required"`

	// The extent at which the rule can be attached across your accounts.
	IncludedScope *RuleScope `validate:"required"`

	// The extent at which the rule can be excluded from the included scope.
	ExcludedScopes []RuleScope

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateRuleAttachmentOptions : Instantiate UpdateRuleAttachmentOptions
func (*ConfigurationGovernanceV1) NewUpdateRuleAttachmentOptions(ruleID string, attachmentID string, ifMatch string, accountID string, includedScope *RuleScope) *UpdateRuleAttachmentOptions {
	return &UpdateRuleAttachmentOptions{
		RuleID:        core.StringPtr(ruleID),
		AttachmentID:  core.StringPtr(attachmentID),
		IfMatch:       core.StringPtr(ifMatch),
		AccountID:     core.StringPtr(accountID),
		IncludedScope: includedScope,
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateRuleAttachmentOptions) SetRuleID(ruleID string) *UpdateRuleAttachmentOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *UpdateRuleAttachmentOptions) SetAttachmentID(attachmentID string) *UpdateRuleAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateRuleAttachmentOptions) SetIfMatch(ifMatch string) *UpdateRuleAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateRuleAttachmentOptions) SetAccountID(accountID string) *UpdateRuleAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIncludedScope : Allow user to set IncludedScope
func (_options *UpdateRuleAttachmentOptions) SetIncludedScope(includedScope *RuleScope) *UpdateRuleAttachmentOptions {
	_options.IncludedScope = includedScope
	return _options
}

// SetExcludedScopes : Allow user to set ExcludedScopes
func (_options *UpdateRuleAttachmentOptions) SetExcludedScopes(excludedScopes []RuleScope) *UpdateRuleAttachmentOptions {
	_options.ExcludedScopes = excludedScopes
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateRuleAttachmentOptions) SetTransactionID(transactionID string) *UpdateRuleAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRuleAttachmentOptions) SetHeaders(param map[string]string) *UpdateRuleAttachmentOptions {
	options.Headers = param
	return options
}

// UpdateRuleOptions : The UpdateRule options.
type UpdateRuleOptions struct {
	// The UUID that uniquely identifies the rule.
	RuleID *string `validate:"required,ne="`

	// Compares a supplied `Etag` value with the version that is stored for the requested resource. If the values match,
	// the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `validate:"required"`

	// A human-readable alias to assign to your rule.
	Name *string `validate:"required"`

	// An extended description of your rule.
	Description *string `validate:"required"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *TargetResource `validate:"required"`

	RequiredConfig RuleRequiredConfigIntf `validate:"required"`

	// The actions that the service must run on your behalf when a request to create or modify the target resource does not
	// comply with your conditions.
	EnforcementActions []EnforcementAction `validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string

	// The type of rule. Rules that you create are `user_defined`.
	RuleType *string

	// Labels that you can use to group and search for similar rules, such as those that help you to meet a specific
	// organization guideline.
	Labels []string

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateRuleOptions.RuleType property.
// The type of rule. Rules that you create are `user_defined`.
const (
	UpdateRuleOptionsRuleTypeUserDefinedConst = "user_defined"
)

// NewUpdateRuleOptions : Instantiate UpdateRuleOptions
func (*ConfigurationGovernanceV1) NewUpdateRuleOptions(ruleID string, ifMatch string, name string, description string, target *TargetResource, requiredConfig RuleRequiredConfigIntf, enforcementActions []EnforcementAction) *UpdateRuleOptions {
	return &UpdateRuleOptions{
		RuleID:             core.StringPtr(ruleID),
		IfMatch:            core.StringPtr(ifMatch),
		Name:               core.StringPtr(name),
		Description:        core.StringPtr(description),
		Target:             target,
		RequiredConfig:     requiredConfig,
		EnforcementActions: enforcementActions,
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *UpdateRuleOptions) SetRuleID(ruleID string) *UpdateRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateRuleOptions) SetIfMatch(ifMatch string) *UpdateRuleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateRuleOptions) SetName(name string) *UpdateRuleOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateRuleOptions) SetDescription(description string) *UpdateRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *UpdateRuleOptions) SetTarget(target *TargetResource) *UpdateRuleOptions {
	_options.Target = target
	return _options
}

// SetRequiredConfig : Allow user to set RequiredConfig
func (_options *UpdateRuleOptions) SetRequiredConfig(requiredConfig RuleRequiredConfigIntf) *UpdateRuleOptions {
	_options.RequiredConfig = requiredConfig
	return _options
}

// SetEnforcementActions : Allow user to set EnforcementActions
func (_options *UpdateRuleOptions) SetEnforcementActions(enforcementActions []EnforcementAction) *UpdateRuleOptions {
	_options.EnforcementActions = enforcementActions
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateRuleOptions) SetAccountID(accountID string) *UpdateRuleOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetRuleType : Allow user to set RuleType
func (_options *UpdateRuleOptions) SetRuleType(ruleType string) *UpdateRuleOptions {
	_options.RuleType = core.StringPtr(ruleType)
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *UpdateRuleOptions) SetLabels(labels []string) *UpdateRuleOptions {
	_options.Labels = labels
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateRuleOptions) SetTransactionID(transactionID string) *UpdateRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRuleOptions) SetHeaders(param map[string]string) *UpdateRuleOptions {
	options.Headers = param
	return options
}

// UpdateTemplateAttachmentOptions : The UpdateTemplateAttachment options.
type UpdateTemplateAttachmentOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// The UUID that uniquely identifies the attachment.
	AttachmentID *string `validate:"required,ne="`

	// Compares a supplied `Etag` value with the version that is stored for the requested resource. If the values match,
	// the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `validate:"required"`

	// The extent at which the template can be attached across your accounts.
	IncludedScope *TemplateScope `validate:"required"`

	ExcludedScopes []TemplateScope

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateTemplateAttachmentOptions : Instantiate UpdateTemplateAttachmentOptions
func (*ConfigurationGovernanceV1) NewUpdateTemplateAttachmentOptions(templateID string, attachmentID string, ifMatch string, accountID string, includedScope *TemplateScope) *UpdateTemplateAttachmentOptions {
	return &UpdateTemplateAttachmentOptions{
		TemplateID:    core.StringPtr(templateID),
		AttachmentID:  core.StringPtr(attachmentID),
		IfMatch:       core.StringPtr(ifMatch),
		AccountID:     core.StringPtr(accountID),
		IncludedScope: includedScope,
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *UpdateTemplateAttachmentOptions) SetTemplateID(templateID string) *UpdateTemplateAttachmentOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *UpdateTemplateAttachmentOptions) SetAttachmentID(attachmentID string) *UpdateTemplateAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateTemplateAttachmentOptions) SetIfMatch(ifMatch string) *UpdateTemplateAttachmentOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateTemplateAttachmentOptions) SetAccountID(accountID string) *UpdateTemplateAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetIncludedScope : Allow user to set IncludedScope
func (_options *UpdateTemplateAttachmentOptions) SetIncludedScope(includedScope *TemplateScope) *UpdateTemplateAttachmentOptions {
	_options.IncludedScope = includedScope
	return _options
}

// SetExcludedScopes : Allow user to set ExcludedScopes
func (_options *UpdateTemplateAttachmentOptions) SetExcludedScopes(excludedScopes []TemplateScope) *UpdateTemplateAttachmentOptions {
	_options.ExcludedScopes = excludedScopes
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateTemplateAttachmentOptions) SetTransactionID(transactionID string) *UpdateTemplateAttachmentOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTemplateAttachmentOptions) SetHeaders(param map[string]string) *UpdateTemplateAttachmentOptions {
	options.Headers = param
	return options
}

// UpdateTemplateOptions : The UpdateTemplate options.
type UpdateTemplateOptions struct {
	// The UUID that uniquely identifies the template.
	TemplateID *string `validate:"required,ne="`

	// Compares a supplied `Etag` value with the version that is stored for the requested resource. If the values match,
	// the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `validate:"required"`

	// Your IBM Cloud account ID.
	AccountID *string `validate:"required"`

	// A human-readablse alias to assign to your template.
	Name *string `validate:"required"`

	// An extended description of your template.
	Description *string `validate:"required"`

	// The properties that describe the resource that you want to target
	// with the rule or template.
	Target *SimpleTargetResource `validate:"required"`

	// A list of default property values to apply to your template.
	CustomizedDefaults []TemplateCustomizedDefaultProperty `validate:"required"`

	// The unique identifier that is used to trace an entire request. If you omit this field, the service generates and
	// sends a transaction ID as a response header of the request. In the case of an error, the transaction ID is set in
	// the `trace` field of the response body.
	//
	// **Note:** To help with debugging logs, it is strongly recommended that you generate and supply a `Transaction-Id`
	// with each request.
	TransactionID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateTemplateOptions : Instantiate UpdateTemplateOptions
func (*ConfigurationGovernanceV1) NewUpdateTemplateOptions(templateID string, ifMatch string, accountID string, name string, description string, target *SimpleTargetResource, customizedDefaults []TemplateCustomizedDefaultProperty) *UpdateTemplateOptions {
	return &UpdateTemplateOptions{
		TemplateID:         core.StringPtr(templateID),
		IfMatch:            core.StringPtr(ifMatch),
		AccountID:          core.StringPtr(accountID),
		Name:               core.StringPtr(name),
		Description:        core.StringPtr(description),
		Target:             target,
		CustomizedDefaults: customizedDefaults,
	}
}

// SetTemplateID : Allow user to set TemplateID
func (_options *UpdateTemplateOptions) SetTemplateID(templateID string) *UpdateTemplateOptions {
	_options.TemplateID = core.StringPtr(templateID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateTemplateOptions) SetIfMatch(ifMatch string) *UpdateTemplateOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpdateTemplateOptions) SetAccountID(accountID string) *UpdateTemplateOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateTemplateOptions) SetName(name string) *UpdateTemplateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateTemplateOptions) SetDescription(description string) *UpdateTemplateOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *UpdateTemplateOptions) SetTarget(target *SimpleTargetResource) *UpdateTemplateOptions {
	_options.Target = target
	return _options
}

// SetCustomizedDefaults : Allow user to set CustomizedDefaults
func (_options *UpdateTemplateOptions) SetCustomizedDefaults(customizedDefaults []TemplateCustomizedDefaultProperty) *UpdateTemplateOptions {
	_options.CustomizedDefaults = customizedDefaults
	return _options
}

// SetTransactionID : Allow user to set TransactionID
func (_options *UpdateTemplateOptions) SetTransactionID(transactionID string) *UpdateTemplateOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTemplateOptions) SetHeaders(param map[string]string) *UpdateTemplateOptions {
	options.Headers = param
	return options
}

// RuleConditionAndLvl2 : A condition with the `and` logical operator.
// This model "extends" RuleCondition
type RuleConditionAndLvl2 struct {
	Description *string `json:"description,omitempty"`

	And []RuleSingleProperty `json:"and" validate:"required"`
}

// NewRuleConditionAndLvl2 : Instantiate RuleConditionAndLvl2 (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleConditionAndLvl2(and []RuleSingleProperty) (_model *RuleConditionAndLvl2, err error) {
	_model = &RuleConditionAndLvl2{
		And: and,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleConditionAndLvl2) isaRuleCondition() bool {
	return true
}

// UnmarshalRuleConditionAndLvl2 unmarshals an instance of RuleConditionAndLvl2 from the specified map of raw messages.
func UnmarshalRuleConditionAndLvl2(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleConditionAndLvl2)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRuleSingleProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleConditionOrLvl2 : A condition with the `or` logical operator.
// This model "extends" RuleCondition
type RuleConditionOrLvl2 struct {
	Description *string `json:"description,omitempty"`

	Or []RuleSingleProperty `json:"or" validate:"required"`
}

// NewRuleConditionOrLvl2 : Instantiate RuleConditionOrLvl2 (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleConditionOrLvl2(or []RuleSingleProperty) (_model *RuleConditionOrLvl2, err error) {
	_model = &RuleConditionOrLvl2{
		Or: or,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleConditionOrLvl2) isaRuleCondition() bool {
	return true
}

// UnmarshalRuleConditionOrLvl2 unmarshals an instance of RuleConditionOrLvl2 from the specified map of raw messages.
func UnmarshalRuleConditionOrLvl2(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleConditionOrLvl2)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRuleSingleProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleConditionSingleProperty : The requirement that must be met to determine the resource's level of compliance in accordance with the rule.
//
// To apply a single property check, define a configuration property and the desired value that you want to check
// against.
// This model "extends" RuleCondition
type RuleConditionSingleProperty struct {
	Description *string `json:"description,omitempty"`

	// A resource configuration variable that describes the property that you want to apply to the target resource.
	//
	// Available options depend on the target service and resource.
	Property *string `json:"property" validate:"required"`

	// The way in which the `property` field is compared to its value.
	//
	// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
	Operator *string `json:"operator" validate:"required"`

	// The way in which you want your property to be applied.
	//
	// Value options differ depending on the rule that you configure. If you use a boolean operator, you do not need to
	// input a value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the RuleConditionSingleProperty.Operator property.
// The way in which the `property` field is compared to its value.
//
// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
const (
	RuleConditionSinglePropertyOperatorIpsInRangeConst           = "ips_in_range"
	RuleConditionSinglePropertyOperatorIsEmptyConst              = "is_empty"
	RuleConditionSinglePropertyOperatorIsFalseConst              = "is_false"
	RuleConditionSinglePropertyOperatorIsNotEmptyConst           = "is_not_empty"
	RuleConditionSinglePropertyOperatorIsTrueConst               = "is_true"
	RuleConditionSinglePropertyOperatorNumEqualsConst            = "num_equals"
	RuleConditionSinglePropertyOperatorNumGreaterThanConst       = "num_greater_than"
	RuleConditionSinglePropertyOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RuleConditionSinglePropertyOperatorNumLessThanConst          = "num_less_than"
	RuleConditionSinglePropertyOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RuleConditionSinglePropertyOperatorNumNotEqualsConst         = "num_not_equals"
	RuleConditionSinglePropertyOperatorStringEqualsConst         = "string_equals"
	RuleConditionSinglePropertyOperatorStringMatchConst          = "string_match"
	RuleConditionSinglePropertyOperatorStringNotEqualsConst      = "string_not_equals"
	RuleConditionSinglePropertyOperatorStringNotMatchConst       = "string_not_match"
	RuleConditionSinglePropertyOperatorStringsInListConst        = "strings_in_list"
)

// NewRuleConditionSingleProperty : Instantiate RuleConditionSingleProperty (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleConditionSingleProperty(property string, operator string) (_model *RuleConditionSingleProperty, err error) {
	_model = &RuleConditionSingleProperty{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleConditionSingleProperty) isaRuleCondition() bool {
	return true
}

// UnmarshalRuleConditionSingleProperty unmarshals an instance of RuleConditionSingleProperty from the specified map of raw messages.
func UnmarshalRuleConditionSingleProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleConditionSingleProperty)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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

// RuleRequiredConfigMultipleProperties : The requirements that must be met to determine the resource's level of compliance in accordance with the rule.
//
// Use logical operators (`and`/`or`) to define multiple property checks and conditions. To define requirements for a
// rule, list one or more property check objects in the `and` array. To add conditions to a property check, use `or`.
// Models which "extend" this model:
// - RuleRequiredConfigMultiplePropertiesConditionOr
// - RuleRequiredConfigMultiplePropertiesConditionAnd
// This model "extends" RuleRequiredConfig
type RuleRequiredConfigMultipleProperties struct {
	Description *string `json:"description,omitempty"`

	Or []RuleConditionIntf `json:"or,omitempty"`

	And []RuleConditionIntf `json:"and,omitempty"`
}

func (*RuleRequiredConfigMultipleProperties) isaRuleRequiredConfigMultipleProperties() bool {
	return true
}

type RuleRequiredConfigMultiplePropertiesIntf interface {
	RuleRequiredConfigIntf
	isaRuleRequiredConfigMultipleProperties() bool
}

func (*RuleRequiredConfigMultipleProperties) isaRuleRequiredConfig() bool {
	return true
}

// UnmarshalRuleRequiredConfigMultipleProperties unmarshals an instance of RuleRequiredConfigMultipleProperties from the specified map of raw messages.
func UnmarshalRuleRequiredConfigMultipleProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequiredConfigMultipleProperties)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleRequiredConfigSingleProperty : The requirement that must be met to determine the resource's level of compliance in accordance with the rule.
//
// To apply a single property check, define a configuration property and the desired value that you want to check
// against.
// This model "extends" RuleRequiredConfig
type RuleRequiredConfigSingleProperty struct {
	Description *string `json:"description,omitempty"`

	// A resource configuration variable that describes the property that you want to apply to the target resource.
	//
	// Available options depend on the target service and resource.
	Property *string `json:"property" validate:"required"`

	// The way in which the `property` field is compared to its value.
	//
	// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
	Operator *string `json:"operator" validate:"required"`

	// The way in which you want your property to be applied.
	//
	// Value options differ depending on the rule that you configure. If you use a boolean operator, you do not need to
	// input a value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the RuleRequiredConfigSingleProperty.Operator property.
// The way in which the `property` field is compared to its value.
//
// To learn more, see the [docs](/docs/security-compliance?topic=security-compliance-what-is-rule#rule-operators).
const (
	RuleRequiredConfigSinglePropertyOperatorIpsInRangeConst           = "ips_in_range"
	RuleRequiredConfigSinglePropertyOperatorIsEmptyConst              = "is_empty"
	RuleRequiredConfigSinglePropertyOperatorIsFalseConst              = "is_false"
	RuleRequiredConfigSinglePropertyOperatorIsNotEmptyConst           = "is_not_empty"
	RuleRequiredConfigSinglePropertyOperatorIsTrueConst               = "is_true"
	RuleRequiredConfigSinglePropertyOperatorNumEqualsConst            = "num_equals"
	RuleRequiredConfigSinglePropertyOperatorNumGreaterThanConst       = "num_greater_than"
	RuleRequiredConfigSinglePropertyOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RuleRequiredConfigSinglePropertyOperatorNumLessThanConst          = "num_less_than"
	RuleRequiredConfigSinglePropertyOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RuleRequiredConfigSinglePropertyOperatorNumNotEqualsConst         = "num_not_equals"
	RuleRequiredConfigSinglePropertyOperatorStringEqualsConst         = "string_equals"
	RuleRequiredConfigSinglePropertyOperatorStringMatchConst          = "string_match"
	RuleRequiredConfigSinglePropertyOperatorStringNotEqualsConst      = "string_not_equals"
	RuleRequiredConfigSinglePropertyOperatorStringNotMatchConst       = "string_not_match"
	RuleRequiredConfigSinglePropertyOperatorStringsInListConst        = "strings_in_list"
)

// NewRuleRequiredConfigSingleProperty : Instantiate RuleRequiredConfigSingleProperty (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleRequiredConfigSingleProperty(property string, operator string) (_model *RuleRequiredConfigSingleProperty, err error) {
	_model = &RuleRequiredConfigSingleProperty{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleRequiredConfigSingleProperty) isaRuleRequiredConfig() bool {
	return true
}

// UnmarshalRuleRequiredConfigSingleProperty unmarshals an instance of RuleRequiredConfigSingleProperty from the specified map of raw messages.
func UnmarshalRuleRequiredConfigSingleProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequiredConfigSingleProperty)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
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

// RuleRequiredConfigMultiplePropertiesConditionAnd : A condition with the `and` logical operator.
// This model "extends" RuleRequiredConfigMultipleProperties
type RuleRequiredConfigMultiplePropertiesConditionAnd struct {
	Description *string `json:"description,omitempty"`

	And []RuleConditionIntf `json:"and" validate:"required"`
}

// NewRuleRequiredConfigMultiplePropertiesConditionAnd : Instantiate RuleRequiredConfigMultiplePropertiesConditionAnd (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleRequiredConfigMultiplePropertiesConditionAnd(and []RuleConditionIntf) (_model *RuleRequiredConfigMultiplePropertiesConditionAnd, err error) {
	_model = &RuleRequiredConfigMultiplePropertiesConditionAnd{
		And: and,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleRequiredConfigMultiplePropertiesConditionAnd) isaRuleRequiredConfigMultipleProperties() bool {
	return true
}

func (*RuleRequiredConfigMultiplePropertiesConditionAnd) isaRuleRequiredConfig() bool {
	return true
}

// UnmarshalRuleRequiredConfigMultiplePropertiesConditionAnd unmarshals an instance of RuleRequiredConfigMultiplePropertiesConditionAnd from the specified map of raw messages.
func UnmarshalRuleRequiredConfigMultiplePropertiesConditionAnd(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequiredConfigMultiplePropertiesConditionAnd)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleRequiredConfigMultiplePropertiesConditionOr : A condition with the `or` logical operator.
// This model "extends" RuleRequiredConfigMultipleProperties
type RuleRequiredConfigMultiplePropertiesConditionOr struct {
	Description *string `json:"description,omitempty"`

	Or []RuleConditionIntf `json:"or" validate:"required"`
}

// NewRuleRequiredConfigMultiplePropertiesConditionOr : Instantiate RuleRequiredConfigMultiplePropertiesConditionOr (Generic Model Constructor)
func (*ConfigurationGovernanceV1) NewRuleRequiredConfigMultiplePropertiesConditionOr(or []RuleConditionIntf) (_model *RuleRequiredConfigMultiplePropertiesConditionOr, err error) {
	_model = &RuleRequiredConfigMultiplePropertiesConditionOr{
		Or: or,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RuleRequiredConfigMultiplePropertiesConditionOr) isaRuleRequiredConfigMultipleProperties() bool {
	return true
}

func (*RuleRequiredConfigMultiplePropertiesConditionOr) isaRuleRequiredConfig() bool {
	return true
}

// UnmarshalRuleRequiredConfigMultiplePropertiesConditionOr unmarshals an instance of RuleRequiredConfigMultiplePropertiesConditionOr from the specified map of raw messages.
func UnmarshalRuleRequiredConfigMultiplePropertiesConditionOr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleRequiredConfigMultiplePropertiesConditionOr)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRuleCondition)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
