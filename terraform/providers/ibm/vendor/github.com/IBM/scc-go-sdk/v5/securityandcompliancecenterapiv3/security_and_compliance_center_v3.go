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
 * IBM OpenAPI SDK Code Generator Version: 3.75.0-726bc7e3-20230713-221716
 */

// Package securityandcompliancecenterapiv3 : Operations and models for the SecurityAndComplianceCenterV3 service
package securityandcompliancecenterapiv3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/scc-go-sdk/v5/common"
	"github.com/go-openapi/strfmt"
)

// SecurityAndComplianceCenterApiV3 : The Security and Compliance Center API reference.
type SecurityAndComplianceCenterApiV3 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.compliance.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "security_and_compliance_center"

// SecurityAndComplianceCenterApiV3Options : Service options
type SecurityAndComplianceCenterApiV3Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSecurityAndComplianceCenterV3 : constructs an instance of SecurityAndComplianceCenterV3 with passed in options.
func NewSecurityAndComplianceCenterV3(options *SecurityAndComplianceCenterApiV3Options) (service *SecurityAndComplianceCenterApiV3, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	if serviceOptions.Authenticator == nil {
		serviceOptions.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		return
	}

	err = baseService.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			return
		}
	}

	service = &SecurityAndComplianceCenterApiV3{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	endpoints := map[string]string{
		"us-south": "https://us-south.compliance.cloud.ibm.com", // Dallas region
		"eu-de":    "https://eu-de.compliance.cloud.ibm.com",    // Frankfurt region
		"eu-fr2":   "https://eu-fr2.compliance.cloud.ibm.com",   // Frankfurt region(Restricted)
		"ca-tor":   "https://ca-tor.compliance.cloud.ibm.com",   // Toronto region
		"au-syd":   "https://au-syd.compliance.cloud.ibm.com",   // Sydney region
		"eu-es":    "https://eu-es.compliance.cloud.ibm.com",    // Madrid region
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "securityAndComplianceCenter" suitable for processing requests.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) Clone() *SecurityAndComplianceCenterApiV3 {
	if core.IsNil(securityAndComplianceCenter) {
		return nil
	}
	clone := *securityAndComplianceCenter
	clone.Service = securityAndComplianceCenter.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) SetServiceURL(url string) error {
	return securityAndComplianceCenter.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetServiceURL() string {
	return securityAndComplianceCenter.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) SetDefaultHeaders(headers http.Header) {
	securityAndComplianceCenter.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) SetEnableGzipCompression(enableGzip bool) {
	securityAndComplianceCenter.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetEnableGzipCompression() bool {
	return securityAndComplianceCenter.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	securityAndComplianceCenter.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DisableRetries() {
	securityAndComplianceCenter.Service.DisableRetries()
}

// GetSettings : List settings
// Retrieve the settings of your service instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetSettings(getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetSettingsWithContext(context.Background(), getSettingsOptions)
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSettingsOptions, "getSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getSettingsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSettings : Update settings
// Update the settings of your service instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateSettings(updateSettingsOptions *UpdateSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.UpdateSettingsWithContext(context.Background(), updateSettingsOptions)
}

// UpdateSettingsWithContext is an alternate form of the UpdateSettings method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateSettingsWithContext(ctx context.Context, updateSettingsOptions *UpdateSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSettingsOptions, "updateSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSettingsOptions, "updateSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateSettingsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "UpdateSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSettingsOptions.ObjectStorage != nil {
		body["object_storage"] = updateSettingsOptions.ObjectStorage
	}
	if updateSettingsOptions.EventNotifications != nil {
		body["event_notifications"] = updateSettingsOptions.EventNotifications
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PostTestEvent : Create a test event
// Send a test event to your Event Notifications instance to ensure that the events that are generated  by Security and
// Compliance Center are being forwarded to Event Notifications. For more information, see [Enabling event
// notifications](/docs/security-compliance?topic=security-compliance-event-notifications#event-notifications-test-api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) PostTestEvent(postTestEventOptions *PostTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.PostTestEventWithContext(context.Background(), postTestEventOptions)
}

// PostTestEventWithContext is an alternate form of the PostTestEvent method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) PostTestEventWithContext(ctx context.Context, postTestEventOptions *PostTestEventOptions) (result *TestEvent, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(postTestEventOptions, "postTestEventOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(postTestEventOptions, "postTestEventOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *postTestEventOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/test_event`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range postTestEventOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "PostTestEvent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTestEvent)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListInstanceAttachments : Get all instance attachments
// Retrieve all instance attachments.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListInstanceAttachments(listInstanceAttachmentsOptions *ListInstanceAttachmentsOptions) (result *ProfileAttachmentCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListInstanceAttachmentsWithContext(context.Background(), listInstanceAttachmentsOptions)
}

// ListInstanceAttachmentsWithContext is an alternate form of the ListInstanceAttachments method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListInstanceAttachmentsWithContext(ctx context.Context, listInstanceAttachmentsOptions *ListInstanceAttachmentsOptions) (result *ProfileAttachmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listInstanceAttachmentsOptions, "listInstanceAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listInstanceAttachmentsOptions, "listInstanceAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listInstanceAttachmentsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listInstanceAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListInstanceAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listInstanceAttachmentsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listInstanceAttachmentsOptions.AccountID))
	}
	if listInstanceAttachmentsOptions.VersionGroupLabel != nil {
		builder.AddQuery("version_group_label", fmt.Sprint(*listInstanceAttachmentsOptions.VersionGroupLabel))
	}
	if listInstanceAttachmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listInstanceAttachmentsOptions.Limit))
	}
	if listInstanceAttachmentsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listInstanceAttachmentsOptions.Sort))
	}
	if listInstanceAttachmentsOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listInstanceAttachmentsOptions.Direction))
	}
	if listInstanceAttachmentsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listInstanceAttachmentsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachmentCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateProfileAttachment : Create a profile attachment
// Create an attachment to link to a profile.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProfileAttachment(createProfileAttachmentOptions *CreateProfileAttachmentOptions) (result *ProfileAttachmentResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateProfileAttachmentWithContext(context.Background(), createProfileAttachmentOptions)
}

// CreateProfileAttachmentWithContext is an alternate form of the CreateProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProfileAttachmentWithContext(ctx context.Context, createProfileAttachmentOptions *CreateProfileAttachmentOptions) (result *ProfileAttachmentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileAttachmentOptions, "createProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProfileAttachmentOptions, "createProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createProfileAttachmentOptions.InstanceID,
		"profile_id":  *createProfileAttachmentOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createProfileAttachmentOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createProfileAttachmentOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createProfileAttachmentOptions.Attachments != nil {
		body["attachments"] = createProfileAttachmentOptions.Attachments
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachmentResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProfileAttachments : Get all attachments tied to a profile
// Retrieve all attachments that are linked to a profile.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfileAttachments(listProfileAttachmentsOptions *ListProfileAttachmentsOptions) (result *ProfileAttachmentCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListProfileAttachmentsWithContext(context.Background(), listProfileAttachmentsOptions)
}

// ListProfileAttachmentsWithContext is an alternate form of the ListProfileAttachments method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfileAttachmentsWithContext(ctx context.Context, listProfileAttachmentsOptions *ListProfileAttachmentsOptions) (result *ProfileAttachmentCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProfileAttachmentsOptions, "listProfileAttachmentsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProfileAttachmentsOptions, "listProfileAttachmentsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listProfileAttachmentsOptions.InstanceID,
		"profile_id":  *listProfileAttachmentsOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProfileAttachmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListProfileAttachments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProfileAttachmentsOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listProfileAttachmentsOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachmentCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfileAttachment : Get an attachment for a profile
// Retrieve an attachment that is linked to a profile by specifying the attachment ID.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProfileAttachment(getProfileAttachmentOptions *GetProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetProfileAttachmentWithContext(context.Background(), getProfileAttachmentOptions)
}

// GetProfileAttachmentWithContext is an alternate form of the GetProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProfileAttachmentWithContext(ctx context.Context, getProfileAttachmentOptions *GetProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileAttachmentOptions, "getProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileAttachmentOptions, "getProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":   *getProfileAttachmentOptions.InstanceID,
		"profile_id":    *getProfileAttachmentOptions.ProfileID,
		"attachment_id": *getProfileAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getProfileAttachmentOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getProfileAttachmentOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceProfileAttachment : Update an attachment
// Update an attachment that is linked to a profile.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfileAttachment(replaceProfileAttachmentOptions *ReplaceProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceProfileAttachmentWithContext(context.Background(), replaceProfileAttachmentOptions)
}

// ReplaceProfileAttachmentWithContext is an alternate form of the ReplaceProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfileAttachmentWithContext(ctx context.Context, replaceProfileAttachmentOptions *ReplaceProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceProfileAttachmentOptions, "replaceProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceProfileAttachmentOptions, "replaceProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":   *replaceProfileAttachmentOptions.InstanceID,
		"profile_id":    *replaceProfileAttachmentOptions.ProfileID,
		"attachment_id": *replaceProfileAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if replaceProfileAttachmentOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*replaceProfileAttachmentOptions.AccountID))
	}

	body := make(map[string]interface{})
	if replaceProfileAttachmentOptions.AttachmentParameters != nil {
		body["attachment_parameters"] = replaceProfileAttachmentOptions.AttachmentParameters
	}
	if replaceProfileAttachmentOptions.Description != nil {
		body["description"] = replaceProfileAttachmentOptions.Description
	}
	if replaceProfileAttachmentOptions.Name != nil {
		body["name"] = replaceProfileAttachmentOptions.Name
	}
	if replaceProfileAttachmentOptions.Notifications != nil {
		body["notifications"] = replaceProfileAttachmentOptions.Notifications
	}
	if replaceProfileAttachmentOptions.Schedule != nil {
		body["schedule"] = replaceProfileAttachmentOptions.Schedule
	}
	if replaceProfileAttachmentOptions.Scope != nil {
		body["scope"] = replaceProfileAttachmentOptions.Scope
	}
	if replaceProfileAttachmentOptions.Status != nil {
		body["status"] = replaceProfileAttachmentOptions.Status
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProfileAttachment : Delete an attachment
// Delete an attachment that is linked to a profile.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteProfileAttachment(deleteProfileAttachmentOptions *DeleteProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteProfileAttachmentWithContext(context.Background(), deleteProfileAttachmentOptions)
}

// DeleteProfileAttachmentWithContext is an alternate form of the DeleteProfileAttachment method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteProfileAttachmentWithContext(ctx context.Context, deleteProfileAttachmentOptions *DeleteProfileAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProfileAttachmentOptions, "deleteProfileAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProfileAttachmentOptions, "deleteProfileAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":   *deleteProfileAttachmentOptions.InstanceID,
		"profile_id":    *deleteProfileAttachmentOptions.ProfileID,
		"attachment_id": *deleteProfileAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments/{attachment_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProfileAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteProfileAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteProfileAttachmentOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteProfileAttachmentOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpgradeAttachment : Upgrade an attachment
// Upgrade an attachment to the latest version of a profile.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule or you can initiate a
// scan at any time. To evaluate your resources, you create an attachment.  An attachment is the association between the
// set of resources that you want to evaluate  and a profile that contains the specific controls that you want to use.
// For more information, see [Running an evaluation for IBM
// Cloud](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpgradeAttachment(upgradeAttachmentOptions *UpgradeAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.UpgradeAttachmentWithContext(context.Background(), upgradeAttachmentOptions)
}

// UpgradeAttachmentWithContext is an alternate form of the UpgradeAttachment method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpgradeAttachmentWithContext(ctx context.Context, upgradeAttachmentOptions *UpgradeAttachmentOptions) (result *ProfileAttachment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(upgradeAttachmentOptions, "upgradeAttachmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(upgradeAttachmentOptions, "upgradeAttachmentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":   *upgradeAttachmentOptions.InstanceID,
		"profile_id":    *upgradeAttachmentOptions.ProfileID,
		"attachment_id": *upgradeAttachmentOptions.AttachmentID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/attachments/{attachment_id}/upgrade`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range upgradeAttachmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "UpgradeAttachment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if upgradeAttachmentOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*upgradeAttachmentOptions.AccountID))
	}

	body := make(map[string]interface{})
	if upgradeAttachmentOptions.AttachmentParameters != nil {
		body["attachment_parameters"] = upgradeAttachmentOptions.AttachmentParameters
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileAttachment)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateCustomControlLibrary : Create a custom control library
// Create a custom control library that is specific to your organization's needs.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateCustomControlLibrary(createCustomControlLibraryOptions *CreateCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateCustomControlLibraryWithContext(context.Background(), createCustomControlLibraryOptions)
}

// CreateCustomControlLibraryWithContext is an alternate form of the CreateCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateCustomControlLibraryWithContext(ctx context.Context, createCustomControlLibraryOptions *CreateCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCustomControlLibraryOptions, "createCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCustomControlLibraryOptions, "createCustomControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createCustomControlLibraryOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/control_libraries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createCustomControlLibraryOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createCustomControlLibraryOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createCustomControlLibraryOptions.ControlLibraryName != nil {
		body["control_library_name"] = createCustomControlLibraryOptions.ControlLibraryName
	}
	if createCustomControlLibraryOptions.ControlLibraryDescription != nil {
		body["control_library_description"] = createCustomControlLibraryOptions.ControlLibraryDescription
	}
	if createCustomControlLibraryOptions.ControlLibraryType != nil {
		body["control_library_type"] = createCustomControlLibraryOptions.ControlLibraryType
	}
	if createCustomControlLibraryOptions.ControlLibraryVersion != nil {
		body["control_library_version"] = createCustomControlLibraryOptions.ControlLibraryVersion
	}
	if createCustomControlLibraryOptions.Controls != nil {
		body["controls"] = createCustomControlLibraryOptions.Controls
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListControlLibraries : Get all control libraries
// Retrieve all the control libraries, including predefined, and custom libraries.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListControlLibraries(listControlLibrariesOptions *ListControlLibrariesOptions) (result *ControlLibraryCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListControlLibrariesWithContext(context.Background(), listControlLibrariesOptions)
}

// ListControlLibrariesWithContext is an alternate form of the ListControlLibraries method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListControlLibrariesWithContext(ctx context.Context, listControlLibrariesOptions *ListControlLibrariesOptions) (result *ControlLibraryCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listControlLibrariesOptions, "listControlLibrariesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listControlLibrariesOptions, "listControlLibrariesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listControlLibrariesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/control_libraries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listControlLibrariesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListControlLibraries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listControlLibrariesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listControlLibrariesOptions.AccountID))
	}
	if listControlLibrariesOptions.ControlLibraryType != nil {
		builder.AddQuery("control_library_type", fmt.Sprint(*listControlLibrariesOptions.ControlLibraryType))
	}
	if listControlLibrariesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listControlLibrariesOptions.Limit))
	}
	if listControlLibrariesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listControlLibrariesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibraryCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceCustomControlLibrary : Update a custom control library
// Update a custom control library by specifying its ID.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceCustomControlLibrary(replaceCustomControlLibraryOptions *ReplaceCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceCustomControlLibraryWithContext(context.Background(), replaceCustomControlLibraryOptions)
}

// ReplaceCustomControlLibraryWithContext is an alternate form of the ReplaceCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceCustomControlLibraryWithContext(ctx context.Context, replaceCustomControlLibraryOptions *ReplaceCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceCustomControlLibraryOptions, "replaceCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceCustomControlLibraryOptions, "replaceCustomControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":        *replaceCustomControlLibraryOptions.InstanceID,
		"control_library_id": *replaceCustomControlLibraryOptions.ControlLibraryID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/control_libraries/{control_library_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if replaceCustomControlLibraryOptions.BssAccount != nil {
		builder.AddQuery("bss_account", fmt.Sprint(*replaceCustomControlLibraryOptions.BssAccount))
	}

	body := make(map[string]interface{})
	if replaceCustomControlLibraryOptions.ControlLibraryName != nil {
		body["control_library_name"] = replaceCustomControlLibraryOptions.ControlLibraryName
	}
	if replaceCustomControlLibraryOptions.ControlLibraryDescription != nil {
		body["control_library_description"] = replaceCustomControlLibraryOptions.ControlLibraryDescription
	}
	if replaceCustomControlLibraryOptions.ControlLibraryType != nil {
		body["control_library_type"] = replaceCustomControlLibraryOptions.ControlLibraryType
	}
	if replaceCustomControlLibraryOptions.ControlLibraryVersion != nil {
		body["control_library_version"] = replaceCustomControlLibraryOptions.ControlLibraryVersion
	}
	if replaceCustomControlLibraryOptions.Controls != nil {
		body["controls"] = replaceCustomControlLibraryOptions.Controls
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetControlLibrary : Get a control library
// View the details of a control library by specifying its ID.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetControlLibrary(getControlLibraryOptions *GetControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetControlLibraryWithContext(context.Background(), getControlLibraryOptions)
}

// GetControlLibraryWithContext is an alternate form of the GetControlLibrary method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetControlLibraryWithContext(ctx context.Context, getControlLibraryOptions *GetControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getControlLibraryOptions, "getControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getControlLibraryOptions, "getControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":        *getControlLibraryOptions.InstanceID,
		"control_library_id": *getControlLibraryOptions.ControlLibraryID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/control_libraries/{control_library_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getControlLibraryOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getControlLibraryOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomControlLibrary : Delete a custom control library
// Delete a custom control library by specifying its ID.
//
// With Security and Compliance Center, you can create a custom control library that is specific to your organization's
// needs.  You define the controls and specifications before you map previously created assessments. Each control has
// several specifications  and assessments that are mapped to it. A specification is a defined requirement that is
// specific to a component. An assessment, or several,  are mapped to each specification with a detailed evaluation that
// is done to check whether the specification is compliant. For more information, see [Creating custom
// libraries](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-custom-library).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteCustomControlLibrary(deleteCustomControlLibraryOptions *DeleteCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteCustomControlLibraryWithContext(context.Background(), deleteCustomControlLibraryOptions)
}

// DeleteCustomControlLibraryWithContext is an alternate form of the DeleteCustomControlLibrary method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteCustomControlLibraryWithContext(ctx context.Context, deleteCustomControlLibraryOptions *DeleteCustomControlLibraryOptions) (result *ControlLibrary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomControlLibraryOptions, "deleteCustomControlLibraryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCustomControlLibraryOptions, "deleteCustomControlLibraryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":        *deleteCustomControlLibraryOptions.InstanceID,
		"control_library_id": *deleteCustomControlLibraryOptions.ControlLibraryID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/control_libraries/{control_library_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCustomControlLibraryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteCustomControlLibrary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteCustomControlLibraryOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteCustomControlLibraryOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalControlLibrary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateProfile : Create a custom profile
// Create a user-defined custom profile.
//
// With Security and Compliance Center, you can create  a profile that is specific to your usecase, by using an existing
// library as a starting point.  For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProfile(createProfileOptions *CreateProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateProfileWithContext(context.Background(), createProfileOptions)
}

// CreateProfileWithContext is an alternate form of the CreateProfile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProfileWithContext(ctx context.Context, createProfileOptions *CreateProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProfileOptions, "createProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProfileOptions, "createProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createProfileOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createProfileOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createProfileOptions.ProfileName != nil {
		body["profile_name"] = createProfileOptions.ProfileName
	}
	if createProfileOptions.ProfileDescription != nil {
		body["profile_description"] = createProfileOptions.ProfileDescription
	}
	if createProfileOptions.ProfileVersion != nil {
		body["profile_version"] = createProfileOptions.ProfileVersion
	}
	if createProfileOptions.Latest != nil {
		body["latest"] = createProfileOptions.Latest
	}
	if createProfileOptions.VersionGroupLabel != nil {
		body["version_group_label"] = createProfileOptions.VersionGroupLabel
	}
	if createProfileOptions.Controls != nil {
		body["controls"] = createProfileOptions.Controls
	}
	if createProfileOptions.DefaultParameters != nil {
		body["default_parameters"] = createProfileOptions.DefaultParameters
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProfiles : Get all profiles
// Retrieve all profiles, including predefined and custom profiles.
//
// With Security and Compliance Center, you can take advantage of predefined profiles  that are curated based on
// industry standards. Or you can choose  to create one that is specific to your usecase by using an existing library as
// a starting point. For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfiles(listProfilesOptions *ListProfilesOptions) (result *ProfileCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListProfilesWithContext(context.Background(), listProfilesOptions)
}

// ListProfilesWithContext is an alternate form of the ListProfiles method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfilesWithContext(ctx context.Context, listProfilesOptions *ListProfilesOptions) (result *ProfileCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProfilesOptions, "listProfilesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProfilesOptions, "listProfilesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listProfilesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProfilesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listProfilesOptions.AccountID))
	}
	if listProfilesOptions.ProfileType != nil {
		builder.AddQuery("profile_type", fmt.Sprint(*listProfilesOptions.ProfileType))
	}
	if listProfilesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProfilesOptions.Limit))
	}
	if listProfilesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listProfilesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceProfile : Update a custom profile
// Update the details of a user-defined profile.
//
// With Security and Compliance Center, you can create  a profile that is specific to your usecase, by using an existing
// library as a starting point.  For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfile(replaceProfileOptions *ReplaceProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceProfileWithContext(context.Background(), replaceProfileOptions)
}

// ReplaceProfileWithContext is an alternate form of the ReplaceProfile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfileWithContext(ctx context.Context, replaceProfileOptions *ReplaceProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceProfileOptions, "replaceProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceProfileOptions, "replaceProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *replaceProfileOptions.InstanceID,
		"profile_id":  *replaceProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if replaceProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*replaceProfileOptions.AccountID))
	}

	body := make(map[string]interface{})
	if replaceProfileOptions.ProfileName != nil {
		body["profile_name"] = replaceProfileOptions.ProfileName
	}
	if replaceProfileOptions.ProfileDescription != nil {
		body["profile_description"] = replaceProfileOptions.ProfileDescription
	}
	if replaceProfileOptions.ProfileVersion != nil {
		body["profile_version"] = replaceProfileOptions.ProfileVersion
	}
	if replaceProfileOptions.ProfileType != nil {
		body["profile_type"] = replaceProfileOptions.ProfileType
	}
	if replaceProfileOptions.Latest != nil {
		body["latest"] = replaceProfileOptions.Latest
	}
	if replaceProfileOptions.VersionGroupLabel != nil {
		body["version_group_label"] = replaceProfileOptions.VersionGroupLabel
	}
	if replaceProfileOptions.Controls != nil {
		body["controls"] = replaceProfileOptions.Controls
	}
	if replaceProfileOptions.DefaultParameters != nil {
		body["default_parameters"] = replaceProfileOptions.DefaultParameters
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProfile : Get a profile
// Retrieve a profile by specifying the profile ID.
//
// With Security and Compliance Center, you can utilize predefined profiles  that are curated based on industry
// standards. Or you can choose  to create one that is specific to your usecase, by using an existing library as a
// starting point. For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProfile(getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetProfileWithContext(context.Background(), getProfileOptions)
}

// GetProfileWithContext is an alternate form of the GetProfile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProfileWithContext(ctx context.Context, getProfileOptions *GetProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProfileOptions, "getProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProfileOptions, "getProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getProfileOptions.InstanceID,
		"profile_id":  *getProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*getProfileOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCustomProfile : Delete a custom profile
// Delete a custom profile by specifying the profile ID.
//
// With Security and Compliance Center, you can utilize predefined profiles  that are curated based on industry
// standards. Or you can choose  to create one that is specific to your usecase, by using an existing library as a
// starting point. For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteCustomProfile(deleteCustomProfileOptions *DeleteCustomProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteCustomProfileWithContext(context.Background(), deleteCustomProfileOptions)
}

// DeleteCustomProfileWithContext is an alternate form of the DeleteCustomProfile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteCustomProfileWithContext(ctx context.Context, deleteCustomProfileOptions *DeleteCustomProfileOptions) (result *Profile, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomProfileOptions, "deleteCustomProfileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCustomProfileOptions, "deleteCustomProfileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteCustomProfileOptions.InstanceID,
		"profile_id":  *deleteCustomProfileOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCustomProfileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteCustomProfile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if deleteCustomProfileOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*deleteCustomProfileOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfile)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceProfileParameters : Update custom profile parameters
// Update the parameters of a custom profile.
//
// With Security and Compliance Center, you can utilize predefined profiles  that are curated based on industry
// standards. Or you can choose  to create one that is specific to your usecase, by using an existing library as a
// starting point. For more information, see [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfileParameters(replaceProfileParametersOptions *ReplaceProfileParametersOptions) (result *ProfileDefaultParametersResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceProfileParametersWithContext(context.Background(), replaceProfileParametersOptions)
}

// ReplaceProfileParametersWithContext is an alternate form of the ReplaceProfileParameters method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceProfileParametersWithContext(ctx context.Context, replaceProfileParametersOptions *ReplaceProfileParametersOptions) (result *ProfileDefaultParametersResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceProfileParametersOptions, "replaceProfileParametersOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceProfileParametersOptions, "replaceProfileParametersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *replaceProfileParametersOptions.InstanceID,
		"profile_id":  *replaceProfileParametersOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/parameters`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceProfileParametersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceProfileParameters")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if replaceProfileParametersOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*replaceProfileParametersOptions.AccountID))
	}

	body := make(map[string]interface{})
	if replaceProfileParametersOptions.ID != nil {
		body["id"] = replaceProfileParametersOptions.ID
	}
	if replaceProfileParametersOptions.DefaultParameters != nil {
		body["default_parameters"] = replaceProfileParametersOptions.DefaultParameters
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileDefaultParametersResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProfileParameters : List profile parameters for a given profile
// List the parameters used in the Profile.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfileParameters(listProfileParametersOptions *ListProfileParametersOptions) (result *ProfileDefaultParametersResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListProfileParametersWithContext(context.Background(), listProfileParametersOptions)
}

// ListProfileParametersWithContext is an alternate form of the ListProfileParameters method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProfileParametersWithContext(ctx context.Context, listProfileParametersOptions *ListProfileParametersOptions) (result *ProfileDefaultParametersResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProfileParametersOptions, "listProfileParametersOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProfileParametersOptions, "listProfileParametersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listProfileParametersOptions.InstanceID,
		"profile_id":  *listProfileParametersOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/parameters`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProfileParametersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListProfileParameters")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProfileDefaultParametersResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CompareProfiles : Compare profiles
// Compare the version of the profile that you're currently using with your attachment to the most recent profile
// version.  By comparing them, you can view what controls were added, removed, or modified. For more information, see
// [Building custom
// profiles](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-build-custom-profiles&interface=api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CompareProfiles(compareProfilesOptions *CompareProfilesOptions) (result *ComparePredefinedProfilesResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CompareProfilesWithContext(context.Background(), compareProfilesOptions)
}

// CompareProfilesWithContext is an alternate form of the CompareProfiles method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CompareProfilesWithContext(ctx context.Context, compareProfilesOptions *CompareProfilesOptions) (result *ComparePredefinedProfilesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(compareProfilesOptions, "compareProfilesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(compareProfilesOptions, "compareProfilesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *compareProfilesOptions.InstanceID,
		"profile_id":  *compareProfilesOptions.ProfileID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/profiles/{profile_id}/compare`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range compareProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CompareProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if compareProfilesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*compareProfilesOptions.AccountID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalComparePredefinedProfilesResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateScope : Create a scope
// Create a scope.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScope(createScopeOptions *CreateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateScopeWithContext(context.Background(), createScopeOptions)
}

// CreateScopeWithContext is an alternate form of the CreateScope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScopeWithContext(ctx context.Context, createScopeOptions *CreateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createScopeOptions, "createScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createScopeOptions, "createScopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createScopeOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createScopeOptions.Name != nil {
		body["name"] = createScopeOptions.Name
	}
	if createScopeOptions.Description != nil {
		body["description"] = createScopeOptions.Description
	}
	if createScopeOptions.Environment != nil {
		body["environment"] = createScopeOptions.Environment
	}
	if createScopeOptions.Properties != nil {
		body["properties"] = createScopeOptions.Properties
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListScopes : Get all scopes
// Get all scopes.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListScopes(listScopesOptions *ListScopesOptions) (result *ScopeCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListScopesWithContext(context.Background(), listScopesOptions)
}

// ListScopesWithContext is an alternate form of the ListScopes method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListScopesWithContext(ctx context.Context, listScopesOptions *ListScopesOptions) (result *ScopeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listScopesOptions, "listScopesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listScopesOptions, "listScopesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listScopesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listScopesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListScopes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listScopesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listScopesOptions.Limit))
	}
	if listScopesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listScopesOptions.Start))
	}
	if listScopesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listScopesOptions.Name))
	}
	if listScopesOptions.Description != nil {
		builder.AddQuery("description", fmt.Sprint(*listScopesOptions.Description))
	}
	if listScopesOptions.Environment != nil {
		builder.AddQuery("environment", fmt.Sprint(*listScopesOptions.Environment))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScopeCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateScope : Update a scope
// Update the details of a scope.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateScope(updateScopeOptions *UpdateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.UpdateScopeWithContext(context.Background(), updateScopeOptions)
}

// UpdateScopeWithContext is an alternate form of the UpdateScope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateScopeWithContext(ctx context.Context, updateScopeOptions *UpdateScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateScopeOptions, "updateScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateScopeOptions, "updateScopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateScopeOptions.InstanceID,
		"scope_id":    *updateScopeOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "UpdateScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateScopeOptions.Name != nil {
		body["name"] = updateScopeOptions.Name
	}
	if updateScopeOptions.Description != nil {
		body["description"] = updateScopeOptions.Description
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScope : Get a scope
// Get a scope by specifying the scope ID.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScope(getScopeOptions *GetScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetScopeWithContext(context.Background(), getScopeOptions)
}

// GetScopeWithContext is an alternate form of the GetScope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScopeWithContext(ctx context.Context, getScopeOptions *GetScopeOptions) (result *Scope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScopeOptions, "getScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScopeOptions, "getScopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getScopeOptions.InstanceID,
		"scope_id":    *getScopeOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteScope : Delete a scope
// Delete a scope by specifying the scope ID.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteScope(deleteScopeOptions *DeleteScopeOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteScopeWithContext(context.Background(), deleteScopeOptions)
}

// DeleteScopeWithContext is an alternate form of the DeleteScope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteScopeWithContext(ctx context.Context, deleteScopeOptions *DeleteScopeOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteScopeOptions, "deleteScopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteScopeOptions, "deleteScopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteScopeOptions.InstanceID,
		"scope_id":    *deleteScopeOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteScopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteScope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, nil)

	return
}

// CreateSubscope : Create a subscope
// Create a subscope.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateSubscope(createSubscopeOptions *CreateSubscopeOptions) (result *SubScopeResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateSubscopeWithContext(context.Background(), createSubscopeOptions)
}

// CreateSubscopeWithContext is an alternate form of the CreateSubscope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateSubscopeWithContext(ctx context.Context, createSubscopeOptions *CreateSubscopeOptions) (result *SubScopeResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSubscopeOptions, "createSubscopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSubscopeOptions, "createSubscopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createSubscopeOptions.InstanceID,
		"scope_id":    *createSubscopeOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}/subscopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSubscopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateSubscope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSubscopeOptions.Subscopes != nil {
		body["subscopes"] = createSubscopeOptions.Subscopes
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSubScopeResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListSubscopes : Get all subscopes
// Get all subscopes.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListSubscopes(listSubscopesOptions *ListSubscopesOptions) (result *SubScopeCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListSubscopesWithContext(context.Background(), listSubscopesOptions)
}

// ListSubscopesWithContext is an alternate form of the ListSubscopes method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListSubscopesWithContext(ctx context.Context, listSubscopesOptions *ListSubscopesOptions) (result *SubScopeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSubscopesOptions, "listSubscopesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSubscopesOptions, "listSubscopesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listSubscopesOptions.InstanceID,
		"scope_id":    *listSubscopesOptions.ScopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}/subscopes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSubscopesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListSubscopes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSubscopesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSubscopesOptions.Limit))
	}
	if listSubscopesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listSubscopesOptions.Start))
	}
	if listSubscopesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listSubscopesOptions.Name))
	}
	if listSubscopesOptions.Description != nil {
		builder.AddQuery("description", fmt.Sprint(*listSubscopesOptions.Description))
	}
	if listSubscopesOptions.Environment != nil {
		builder.AddQuery("environment", fmt.Sprint(*listSubscopesOptions.Environment))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSubScopeCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetSubscope : Get a subscope
// Get the subscope details.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetSubscope(getSubscopeOptions *GetSubscopeOptions) (result *SubScope, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetSubscopeWithContext(context.Background(), getSubscopeOptions)
}

// GetSubscopeWithContext is an alternate form of the GetSubscope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetSubscopeWithContext(ctx context.Context, getSubscopeOptions *GetSubscopeOptions) (result *SubScope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSubscopeOptions, "getSubscopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSubscopeOptions, "getSubscopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getSubscopeOptions.InstanceID,
		"scope_id":    *getSubscopeOptions.ScopeID,
		"subscope_id": *getSubscopeOptions.SubscopeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}/subscopes/{subscope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSubscopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetSubscope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSubScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateSubscope : Update a subscope
// Update the subscope details.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateSubscope(updateSubscopeOptions *UpdateSubscopeOptions) (result *SubScope, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.UpdateSubscopeWithContext(context.Background(), updateSubscopeOptions)
}

// UpdateSubscopeWithContext is an alternate form of the UpdateSubscope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateSubscopeWithContext(ctx context.Context, updateSubscopeOptions *UpdateSubscopeOptions) (result *SubScope, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSubscopeOptions, "updateSubscopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSubscopeOptions, "updateSubscopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *updateSubscopeOptions.InstanceID,
		"scope_id":    *updateSubscopeOptions.ScopeID,
		"subscope_id": *updateSubscopeOptions.SubscopeID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}/subscopes/{subscope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSubscopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "UpdateSubscope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSubscopeOptions.Name != nil {
		body["name"] = updateSubscopeOptions.Name
	}
	if updateSubscopeOptions.Description != nil {
		body["description"] = updateSubscopeOptions.Description
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSubScope)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteSubscope : Delete a subscope
// Delete the subscope by specifying the subscope ID.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteSubscope(deleteSubscopeOptions *DeleteSubscopeOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteSubscopeWithContext(context.Background(), deleteSubscopeOptions)
}

// DeleteSubscopeWithContext is an alternate form of the DeleteSubscope method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteSubscopeWithContext(ctx context.Context, deleteSubscopeOptions *DeleteSubscopeOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSubscopeOptions, "deleteSubscopeOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSubscopeOptions, "deleteSubscopeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteSubscopeOptions.InstanceID,
		"scope_id":    *deleteSubscopeOptions.ScopeID,
		"subscope_id": *deleteSubscopeOptions.SubscopeID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scopes/{scope_id}/subscopes/{subscope_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSubscopeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteSubscope")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, nil)

	return
}

// CreateTarget : Create a target
// Creates a target to scan against.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateTarget(createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateTargetWithContext(context.Background(), createTargetOptions)
}

// CreateTargetWithContext is an alternate form of the CreateTarget method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateTargetWithContext(ctx context.Context, createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTargetOptions, "createTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTargetOptions, "createTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createTargetOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/targets`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTargetOptions.AccountID != nil {
		body["account_id"] = createTargetOptions.AccountID
	}
	if createTargetOptions.TrustedProfileID != nil {
		body["trusted_profile_id"] = createTargetOptions.TrustedProfileID
	}
	if createTargetOptions.Name != nil {
		body["name"] = createTargetOptions.Name
	}
	if createTargetOptions.Credentials != nil {
		body["credentials"] = createTargetOptions.Credentials
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListTargets : Get a list of targets with pagination
// Returns a list of targets.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListTargets(listTargetsOptions *ListTargetsOptions) (result *TargetCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListTargetsWithContext(context.Background(), listTargetsOptions)
}

// ListTargetsWithContext is an alternate form of the ListTargets method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListTargetsWithContext(ctx context.Context, listTargetsOptions *ListTargetsOptions) (result *TargetCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTargetsOptions, "listTargetsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTargetsOptions, "listTargetsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listTargetsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/targets`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listTargetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListTargets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTargetCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTarget : Get a target by ID
// Retrieves a target by its ID association.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetTarget(getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetTargetWithContext(context.Background(), getTargetOptions)
}

// GetTargetWithContext is an alternate form of the GetTarget method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetTargetWithContext(ctx context.Context, getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTargetOptions, "getTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTargetOptions, "getTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getTargetOptions.InstanceID,
		"target_id":   *getTargetOptions.TargetID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/targets/{target_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTarget : replace a target by ID
// Updates a target by its ID.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceTarget(replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceTargetWithContext(context.Background(), replaceTargetOptions)
}

// ReplaceTargetWithContext is an alternate form of the ReplaceTarget method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceTargetWithContext(ctx context.Context, replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTargetOptions, "replaceTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceTargetOptions, "replaceTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *replaceTargetOptions.InstanceID,
		"target_id":   *replaceTargetOptions.TargetID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/targets/{target_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTargetOptions.AccountID != nil {
		body["account_id"] = replaceTargetOptions.AccountID
	}
	if replaceTargetOptions.TrustedProfileID != nil {
		body["trusted_profile_id"] = replaceTargetOptions.TrustedProfileID
	}
	if replaceTargetOptions.Name != nil {
		body["name"] = replaceTargetOptions.Name
	}
	if replaceTargetOptions.Credentials != nil {
		body["credentials"] = replaceTargetOptions.Credentials
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTarget : Delete a target by ID
// Deletes a target by the ID.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteTarget(deleteTargetOptions *DeleteTargetOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteTargetWithContext(context.Background(), deleteTargetOptions)
}

// DeleteTargetWithContext is an alternate form of the DeleteTarget method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteTargetWithContext(ctx context.Context, deleteTargetOptions *DeleteTargetOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTargetOptions, "deleteTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTargetOptions, "deleteTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteTargetOptions.InstanceID,
		"target_id":   *deleteTargetOptions.TargetID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/targets/{target_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, nil)

	return
}

// CreateProviderTypeInstance : Create a provider type instance
// Create an instance of a provider type. For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProviderTypeInstance(createProviderTypeInstanceOptions *CreateProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateProviderTypeInstanceWithContext(context.Background(), createProviderTypeInstanceOptions)
}

// CreateProviderTypeInstanceWithContext is an alternate form of the CreateProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateProviderTypeInstanceWithContext(ctx context.Context, createProviderTypeInstanceOptions *CreateProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProviderTypeInstanceOptions, "createProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createProviderTypeInstanceOptions, "createProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":      *createProviderTypeInstanceOptions.InstanceID,
		"provider_type_id": *createProviderTypeInstanceOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}/provider_type_instances`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProviderTypeInstanceOptions.Name != nil {
		body["name"] = createProviderTypeInstanceOptions.Name
	}
	if createProviderTypeInstanceOptions.Attributes != nil {
		body["attributes"] = createProviderTypeInstanceOptions.Attributes
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListProviderTypeInstances : List instances of a specific provider type
// Retrieve all instances of a provider type. For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProviderTypeInstances(listProviderTypeInstancesOptions *ListProviderTypeInstancesOptions) (result *ProviderTypeInstanceCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListProviderTypeInstancesWithContext(context.Background(), listProviderTypeInstancesOptions)
}

// ListProviderTypeInstancesWithContext is an alternate form of the ListProviderTypeInstances method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProviderTypeInstancesWithContext(ctx context.Context, listProviderTypeInstancesOptions *ListProviderTypeInstancesOptions) (result *ProviderTypeInstanceCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProviderTypeInstancesOptions, "listProviderTypeInstancesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProviderTypeInstancesOptions, "listProviderTypeInstancesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":      *listProviderTypeInstancesOptions.InstanceID,
		"provider_type_id": *listProviderTypeInstancesOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}/provider_type_instances`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProviderTypeInstancesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListProviderTypeInstances")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstanceCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProviderTypeInstance : Get a provider type instance
// Retrieve a provider type instance by specifying the provider type ID, and Security and Compliance Center instance ID.
// For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProviderTypeInstance(getProviderTypeInstanceOptions *GetProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetProviderTypeInstanceWithContext(context.Background(), getProviderTypeInstanceOptions)
}

// GetProviderTypeInstanceWithContext is an alternate form of the GetProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProviderTypeInstanceWithContext(ctx context.Context, getProviderTypeInstanceOptions *GetProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProviderTypeInstanceOptions, "getProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProviderTypeInstanceOptions, "getProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":               *getProviderTypeInstanceOptions.InstanceID,
		"provider_type_id":          *getProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *getProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateProviderTypeInstance : Update a provider type instance
// Update a provider type instance. For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateProviderTypeInstance(updateProviderTypeInstanceOptions *UpdateProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.UpdateProviderTypeInstanceWithContext(context.Background(), updateProviderTypeInstanceOptions)
}

// UpdateProviderTypeInstanceWithContext is an alternate form of the UpdateProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) UpdateProviderTypeInstanceWithContext(ctx context.Context, updateProviderTypeInstanceOptions *UpdateProviderTypeInstanceOptions) (result *ProviderTypeInstance, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateProviderTypeInstanceOptions, "updateProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateProviderTypeInstanceOptions, "updateProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":               *updateProviderTypeInstanceOptions.InstanceID,
		"provider_type_id":          *updateProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *updateProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "UpdateProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateProviderTypeInstanceOptions.Name != nil {
		body["name"] = updateProviderTypeInstanceOptions.Name
	}
	if updateProviderTypeInstanceOptions.Attributes != nil {
		body["attributes"] = updateProviderTypeInstanceOptions.Attributes
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeInstance)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteProviderTypeInstance : Delete a provider type instance
// Remove a provider type instance. For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteProviderTypeInstance(deleteProviderTypeInstanceOptions *DeleteProviderTypeInstanceOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteProviderTypeInstanceWithContext(context.Background(), deleteProviderTypeInstanceOptions)
}

// DeleteProviderTypeInstanceWithContext is an alternate form of the DeleteProviderTypeInstance method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteProviderTypeInstanceWithContext(ctx context.Context, deleteProviderTypeInstanceOptions *DeleteProviderTypeInstanceOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProviderTypeInstanceOptions, "deleteProviderTypeInstanceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteProviderTypeInstanceOptions, "deleteProviderTypeInstanceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":               *deleteProviderTypeInstanceOptions.InstanceID,
		"provider_type_id":          *deleteProviderTypeInstanceOptions.ProviderTypeID,
		"provider_type_instance_id": *deleteProviderTypeInstanceOptions.ProviderTypeInstanceID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}/provider_type_instances/{provider_type_instance_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteProviderTypeInstanceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteProviderTypeInstance")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, nil)

	return
}

// ListProviderTypes : List provider types
// List all the registered provider types or integrations such as Workload Protection available to connect to Security
// and Compliance Center.  For more information about connecting Workload Protection with the Security and Compliance
// Center, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProviderTypes(listProviderTypesOptions *ListProviderTypesOptions) (result *ProviderTypeCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListProviderTypesWithContext(context.Background(), listProviderTypesOptions)
}

// ListProviderTypesWithContext is an alternate form of the ListProviderTypes method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListProviderTypesWithContext(ctx context.Context, listProviderTypesOptions *ListProviderTypesOptions) (result *ProviderTypeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listProviderTypesOptions, "listProviderTypesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listProviderTypesOptions, "listProviderTypesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listProviderTypesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listProviderTypesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListProviderTypes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderTypeCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetProviderTypeByID : Get a provider type
// Retrieve a provider type by specifying its ID. For more information about integrations, see [Connecting Workload
// Protection](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-setup-workload-protection).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProviderTypeByID(getProviderTypeByIDOptions *GetProviderTypeByIDOptions) (result *ProviderType, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetProviderTypeByIDWithContext(context.Background(), getProviderTypeByIDOptions)
}

// GetProviderTypeByIDWithContext is an alternate form of the GetProviderTypeByID method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetProviderTypeByIDWithContext(ctx context.Context, getProviderTypeByIDOptions *GetProviderTypeByIDOptions) (result *ProviderType, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProviderTypeByIDOptions, "getProviderTypeByIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getProviderTypeByIDOptions, "getProviderTypeByIDOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id":      *getProviderTypeByIDOptions.InstanceID,
		"provider_type_id": *getProviderTypeByIDOptions.ProviderTypeID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/provider_types/{provider_type_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getProviderTypeByIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetProviderTypeByID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProviderType)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScanReport : Get a scan report
// Retrieve the scan report by specifying the ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScanReport(getScanReportOptions *GetScanReportOptions) (result *ScanReport, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetScanReportWithContext(context.Background(), getScanReportOptions)
}

// GetScanReportWithContext is an alternate form of the GetScanReport method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScanReportWithContext(ctx context.Context, getScanReportOptions *GetScanReportOptions) (result *ScanReport, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScanReportOptions, "getScanReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScanReportOptions, "getScanReportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getScanReportOptions.InstanceID,
		"report_id":   *getScanReportOptions.ReportID,
		"job_id":      *getScanReportOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/scan_reports/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScanReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetScanReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScanReport)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetScanReportDownloadFile : Get a scan report details
// Download the scan report with evaluation details for the specified ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScanReportDownloadFile(getScanReportDownloadFileOptions *GetScanReportDownloadFileOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetScanReportDownloadFileWithContext(context.Background(), getScanReportDownloadFileOptions)
}

// GetScanReportDownloadFileWithContext is an alternate form of the GetScanReportDownloadFile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetScanReportDownloadFileWithContext(ctx context.Context, getScanReportDownloadFileOptions *GetScanReportDownloadFileOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getScanReportDownloadFileOptions, "getScanReportDownloadFileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getScanReportDownloadFileOptions, "getScanReportDownloadFileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getScanReportDownloadFileOptions.InstanceID,
		"report_id":   *getScanReportDownloadFileOptions.ReportID,
		"job_id":      *getScanReportDownloadFileOptions.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/scan_reports/{job_id}/download`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getScanReportDownloadFileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetScanReportDownloadFile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/csv")
	if getScanReportDownloadFileOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getScanReportDownloadFileOptions.Accept))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, &result)

	return
}

// GetLatestReports : List latest reports
// Retrieve the latest reports, which are grouped by profile ID, scope ID, and attachment ID. For more information, see
// [Viewing results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetLatestReports(getLatestReportsOptions *GetLatestReportsOptions) (result *ReportLatest, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetLatestReportsWithContext(context.Background(), getLatestReportsOptions)
}

// GetLatestReportsWithContext is an alternate form of the GetLatestReports method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetLatestReportsWithContext(ctx context.Context, getLatestReportsOptions *GetLatestReportsOptions) (result *ReportLatest, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLatestReportsOptions, "getLatestReportsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLatestReportsOptions, "getLatestReportsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getLatestReportsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/latest`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLatestReportsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetLatestReports")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getLatestReportsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*getLatestReportsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportLatest)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReports : List reports
// Retrieve a page of reports that are filtered by the specified parameters. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReports(listReportsOptions *ListReportsOptions) (result *ReportCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListReportsWithContext(context.Background(), listReportsOptions)
}

// ListReportsWithContext is an alternate form of the ListReports method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReportsWithContext(ctx context.Context, listReportsOptions *ListReportsOptions) (result *ReportCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listReportsOptions, "listReportsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listReportsOptions, "listReportsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listReportsOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListReports")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listReportsOptions.AttachmentID != nil {
		builder.AddQuery("attachment_id", fmt.Sprint(*listReportsOptions.AttachmentID))
	}
	if listReportsOptions.GroupID != nil {
		builder.AddQuery("group_id", fmt.Sprint(*listReportsOptions.GroupID))
	}
	if listReportsOptions.ProfileID != nil {
		builder.AddQuery("profile_id", fmt.Sprint(*listReportsOptions.ProfileID))
	}
	if listReportsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listReportsOptions.Type))
	}
	if listReportsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportsOptions.Start))
	}
	if listReportsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportsOptions.Limit))
	}
	if listReportsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listReportsOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReport : Get a report
// Retrieve a specified report and filter the report details by the specified scope ID and/or subscope ID. For more
// information, see [Viewing results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReport(getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportWithContext(context.Background(), getReportOptions)
}

// GetReportWithContext is an alternate form of the GetReport method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportWithContext(ctx context.Context, getReportOptions *GetReportOptions) (result *Report, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportOptions, "getReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportOptions, "getReportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"report_id":   *getReportOptions.ReportID,
		"instance_id": *getReportOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getReportOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*getReportOptions.ScopeID))
	}
	if getReportOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*getReportOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReport)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportSummary : Get a report summary
// Retrieve the complete summarized information for a report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportSummary(getReportSummaryOptions *GetReportSummaryOptions) (result *ReportSummary, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportSummaryWithContext(context.Background(), getReportSummaryOptions)
}

// GetReportSummaryWithContext is an alternate form of the GetReportSummary method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportSummaryWithContext(ctx context.Context, getReportSummaryOptions *GetReportSummaryOptions) (result *ReportSummary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportSummaryOptions, "getReportSummaryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportSummaryOptions, "getReportSummaryOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportSummaryOptions.InstanceID,
		"report_id":   *getReportSummaryOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/summary`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportSummaryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportSummary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportSummary)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportDownloadFile : Get report evaluation details
// Download a .csv file to inspect the evaluation details of a specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportDownloadFile(getReportDownloadFileOptions *GetReportDownloadFileOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportDownloadFileWithContext(context.Background(), getReportDownloadFileOptions)
}

// GetReportDownloadFileWithContext is an alternate form of the GetReportDownloadFile method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportDownloadFileWithContext(ctx context.Context, getReportDownloadFileOptions *GetReportDownloadFileOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportDownloadFileOptions, "getReportDownloadFileOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportDownloadFileOptions, "getReportDownloadFileOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportDownloadFileOptions.InstanceID,
		"report_id":   *getReportDownloadFileOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/download`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportDownloadFileOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportDownloadFile")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/csv")
	if getReportDownloadFileOptions.Accept != nil {
		builder.AddHeader("Accept", fmt.Sprint(*getReportDownloadFileOptions.Accept))
	}

	if getReportDownloadFileOptions.ExcludeSummary != nil {
		builder.AddQuery("exclude_summary", fmt.Sprint(*getReportDownloadFileOptions.ExcludeSummary))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, &result)

	return
}

// GetReportControls : Get report controls
// Retrieve a sorted and filtered list of controls for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportControls(getReportControlsOptions *GetReportControlsOptions) (result *ReportControls, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportControlsWithContext(context.Background(), getReportControlsOptions)
}

// GetReportControlsWithContext is an alternate form of the GetReportControls method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportControlsWithContext(ctx context.Context, getReportControlsOptions *GetReportControlsOptions) (result *ReportControls, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportControlsOptions, "getReportControlsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportControlsOptions, "getReportControlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportControlsOptions.InstanceID,
		"report_id":   *getReportControlsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/controls`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportControlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportControls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getReportControlsOptions.ControlID != nil {
		builder.AddQuery("control_id", fmt.Sprint(*getReportControlsOptions.ControlID))
	}
	if getReportControlsOptions.ControlName != nil {
		builder.AddQuery("control_name", fmt.Sprint(*getReportControlsOptions.ControlName))
	}
	if getReportControlsOptions.ControlDescription != nil {
		builder.AddQuery("control_description", fmt.Sprint(*getReportControlsOptions.ControlDescription))
	}
	if getReportControlsOptions.ControlCategory != nil {
		builder.AddQuery("control_category", fmt.Sprint(*getReportControlsOptions.ControlCategory))
	}
	if getReportControlsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*getReportControlsOptions.Status))
	}
	if getReportControlsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*getReportControlsOptions.Sort))
	}
	if getReportControlsOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*getReportControlsOptions.ScopeID))
	}
	if getReportControlsOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*getReportControlsOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportControls)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportRule : Get a report rule
// Retrieve the rule by specifying the report ID and rule ID. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportRule(getReportRuleOptions *GetReportRuleOptions) (result *RuleInfo, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportRuleWithContext(context.Background(), getReportRuleOptions)
}

// GetReportRuleWithContext is an alternate form of the GetReportRule method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportRuleWithContext(ctx context.Context, getReportRuleOptions *GetReportRuleOptions) (result *RuleInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportRuleOptions, "getReportRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportRuleOptions, "getReportRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportRuleOptions.InstanceID,
		"report_id":   *getReportRuleOptions.ReportID,
		"rule_id":     *getReportRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleInfo)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReportEvaluations : List report evaluations
// Get a paginated list of evaluations for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReportEvaluations(listReportEvaluationsOptions *ListReportEvaluationsOptions) (result *EvaluationPage, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListReportEvaluationsWithContext(context.Background(), listReportEvaluationsOptions)
}

// ListReportEvaluationsWithContext is an alternate form of the ListReportEvaluations method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReportEvaluationsWithContext(ctx context.Context, listReportEvaluationsOptions *ListReportEvaluationsOptions) (result *EvaluationPage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listReportEvaluationsOptions, "listReportEvaluationsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listReportEvaluationsOptions, "listReportEvaluationsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listReportEvaluationsOptions.InstanceID,
		"report_id":   *listReportEvaluationsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/evaluations`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportEvaluationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListReportEvaluations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listReportEvaluationsOptions.AssessmentID != nil {
		builder.AddQuery("assessment_id", fmt.Sprint(*listReportEvaluationsOptions.AssessmentID))
	}
	if listReportEvaluationsOptions.AssessmentMethod != nil {
		builder.AddQuery("assessment_method", fmt.Sprint(*listReportEvaluationsOptions.AssessmentMethod))
	}
	if listReportEvaluationsOptions.ComponentID != nil {
		builder.AddQuery("component_id", fmt.Sprint(*listReportEvaluationsOptions.ComponentID))
	}
	if listReportEvaluationsOptions.TargetID != nil {
		builder.AddQuery("target_id", fmt.Sprint(*listReportEvaluationsOptions.TargetID))
	}
	if listReportEvaluationsOptions.TargetEnv != nil {
		builder.AddQuery("target_env", fmt.Sprint(*listReportEvaluationsOptions.TargetEnv))
	}
	if listReportEvaluationsOptions.TargetName != nil {
		builder.AddQuery("target_name", fmt.Sprint(*listReportEvaluationsOptions.TargetName))
	}
	if listReportEvaluationsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listReportEvaluationsOptions.Status))
	}
	if listReportEvaluationsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportEvaluationsOptions.Start))
	}
	if listReportEvaluationsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportEvaluationsOptions.Limit))
	}
	if listReportEvaluationsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listReportEvaluationsOptions.Sort))
	}
	if listReportEvaluationsOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*listReportEvaluationsOptions.ScopeID))
	}
	if listReportEvaluationsOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*listReportEvaluationsOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEvaluationPage)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListReportResources : List report resources
// Get a paginated list of resources for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReportResources(listReportResourcesOptions *ListReportResourcesOptions) (result *ResourcePage, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListReportResourcesWithContext(context.Background(), listReportResourcesOptions)
}

// ListReportResourcesWithContext is an alternate form of the ListReportResources method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListReportResourcesWithContext(ctx context.Context, listReportResourcesOptions *ListReportResourcesOptions) (result *ResourcePage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listReportResourcesOptions, "listReportResourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listReportResourcesOptions, "listReportResourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listReportResourcesOptions.InstanceID,
		"report_id":   *listReportResourcesOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/resources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listReportResourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListReportResources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listReportResourcesOptions.ID != nil {
		builder.AddQuery("id", fmt.Sprint(*listReportResourcesOptions.ID))
	}
	if listReportResourcesOptions.ResourceName != nil {
		builder.AddQuery("resource_name", fmt.Sprint(*listReportResourcesOptions.ResourceName))
	}
	if listReportResourcesOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*listReportResourcesOptions.AccountID))
	}
	if listReportResourcesOptions.ComponentID != nil {
		builder.AddQuery("component_id", fmt.Sprint(*listReportResourcesOptions.ComponentID))
	}
	if listReportResourcesOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listReportResourcesOptions.Status))
	}
	if listReportResourcesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listReportResourcesOptions.Sort))
	}
	if listReportResourcesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listReportResourcesOptions.Start))
	}
	if listReportResourcesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listReportResourcesOptions.Limit))
	}
	if listReportResourcesOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*listReportResourcesOptions.ScopeID))
	}
	if listReportResourcesOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*listReportResourcesOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalResourcePage)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportTags : List report tags
// Retrieve a list of tags for the specified report. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportTags(getReportTagsOptions *GetReportTagsOptions) (result *ReportTags, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportTagsWithContext(context.Background(), getReportTagsOptions)
}

// GetReportTagsWithContext is an alternate form of the GetReportTags method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportTagsWithContext(ctx context.Context, getReportTagsOptions *GetReportTagsOptions) (result *ReportTags, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportTagsOptions, "getReportTagsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportTagsOptions, "getReportTagsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportTagsOptions.InstanceID,
		"report_id":   *getReportTagsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/tags`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportTagsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportTags")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportTags)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetReportViolationsDrift : Get report violations drift
// Get a list of report violation data points for the specified report and time frame. For more information, see
// [Viewing results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportViolationsDrift(getReportViolationsDriftOptions *GetReportViolationsDriftOptions) (result *ReportViolationsDrift, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetReportViolationsDriftWithContext(context.Background(), getReportViolationsDriftOptions)
}

// GetReportViolationsDriftWithContext is an alternate form of the GetReportViolationsDrift method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetReportViolationsDriftWithContext(ctx context.Context, getReportViolationsDriftOptions *GetReportViolationsDriftOptions) (result *ReportViolationsDrift, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getReportViolationsDriftOptions, "getReportViolationsDriftOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getReportViolationsDriftOptions, "getReportViolationsDriftOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getReportViolationsDriftOptions.InstanceID,
		"report_id":   *getReportViolationsDriftOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/violations_drift`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getReportViolationsDriftOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetReportViolationsDrift")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getReportViolationsDriftOptions.ScanTimeDuration != nil {
		builder.AddQuery("scan_time_duration", fmt.Sprint(*getReportViolationsDriftOptions.ScanTimeDuration))
	}
	if getReportViolationsDriftOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*getReportViolationsDriftOptions.ScopeID))
	}
	if getReportViolationsDriftOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*getReportViolationsDriftOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalReportViolationsDrift)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListScanReports : List scan reports
// Get a list of scan reports and view the status of report generation in progress. For more information, see [Viewing
// results](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-results).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListScanReports(listScanReportsOptions *ListScanReportsOptions) (result *ScanReportCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListScanReportsWithContext(context.Background(), listScanReportsOptions)
}

// ListScanReportsWithContext is an alternate form of the ListScanReports method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListScanReportsWithContext(ctx context.Context, listScanReportsOptions *ListScanReportsOptions) (result *ScanReportCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listScanReportsOptions, "listScanReportsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listScanReportsOptions, "listScanReportsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listScanReportsOptions.InstanceID,
		"report_id":   *listScanReportsOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/scan_reports`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listScanReportsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListScanReports")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listScanReportsOptions.ScopeID != nil {
		builder.AddQuery("scope_id", fmt.Sprint(*listScanReportsOptions.ScopeID))
	}
	if listScanReportsOptions.SubscopeID != nil {
		builder.AddQuery("subscope_id", fmt.Sprint(*listScanReportsOptions.SubscopeID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalScanReportCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateScanReport : Create a scan report
// Create a scan report for a specific scope or sub-scope. For more information, see [Defining custom
// rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScanReport(createScanReportOptions *CreateScanReportOptions) (result *CreateScanReport, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateScanReportWithContext(context.Background(), createScanReportOptions)
}

// CreateScanReportWithContext is an alternate form of the CreateScanReport method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScanReportWithContext(ctx context.Context, createScanReportOptions *CreateScanReportOptions) (result *CreateScanReport, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createScanReportOptions, "createScanReportOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createScanReportOptions, "createScanReportOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createScanReportOptions.InstanceID,
		"report_id":   *createScanReportOptions.ReportID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/reports/{report_id}/scan_reports`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createScanReportOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateScanReport")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createScanReportOptions.Format != nil {
		body["format"] = createScanReportOptions.Format
	}
	if createScanReportOptions.ScopeID != nil {
		body["scope_id"] = createScanReportOptions.ScopeID
	}
	if createScanReportOptions.SubscopeID != nil {
		body["subscope_id"] = createScanReportOptions.SubscopeID
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateScanReport)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateScan : Create a scan
// Create a scan to evaluate your resources.
//
// With Security and Compliance Center, you can evaluate your resources  on a recurring schedule. If your attachment
// exists, but you don't want to wait for the next  scan to see your posture, you can initiate an on-demand scan. For
// more information, see [Running a scan on
// demand](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scan-resources#scan-ondemand-api).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScan(createScanOptions *CreateScanOptions) (result *CreateScanResponse, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateScanWithContext(context.Background(), createScanOptions)
}

// CreateScanWithContext is an alternate form of the CreateScan method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateScanWithContext(ctx context.Context, createScanOptions *CreateScanOptions) (result *CreateScanResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createScanOptions, "createScanOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createScanOptions, "createScanOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createScanOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/scans`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createScanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateScan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createScanOptions.AccountID != nil {
		builder.AddQuery("account_id", fmt.Sprint(*createScanOptions.AccountID))
	}

	body := make(map[string]interface{})
	if createScanOptions.AttachmentID != nil {
		body["attachment_id"] = createScanOptions.AttachmentID
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateScanResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRules : Get all rules
// Retrieve all the rules that you use to target the exact configuration properties  that you need to ensure are
// compliant. For more information, see [Defining custom
// rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListRules(listRulesOptions *ListRulesOptions) (result *RuleCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListRulesWithContext(context.Background(), listRulesOptions)
}

// ListRulesWithContext is an alternate form of the ListRules method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListRulesWithContext(ctx context.Context, listRulesOptions *ListRulesOptions) (result *RuleCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRulesOptions, "listRulesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listRulesOptions, "listRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *listRulesOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listRulesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listRulesOptions.Limit))
	}
	if listRulesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listRulesOptions.Start))
	}
	if listRulesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listRulesOptions.Type))
	}
	if listRulesOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listRulesOptions.Search))
	}
	if listRulesOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listRulesOptions.ServiceName))
	}
	if listRulesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listRulesOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateRule : Create a custom rule
// Create a custom rule to to target the exact configuration properties  that you need to evaluate your resources for
// compliance. For more information, see [Defining custom
// rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateRule(createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.CreateRuleWithContext(context.Background(), createRuleOptions)
}

// CreateRuleWithContext is an alternate form of the CreateRule method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) CreateRuleWithContext(ctx context.Context, createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRuleOptions, "createRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRuleOptions, "createRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *createRuleOptions.InstanceID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "CreateRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRuleOptions.Description != nil {
		body["description"] = createRuleOptions.Description
	}
	if createRuleOptions.Target != nil {
		body["target"] = createRuleOptions.Target
	}
	if createRuleOptions.RequiredConfig != nil {
		body["required_config"] = createRuleOptions.RequiredConfig
	}
	if createRuleOptions.Version != nil {
		body["version"] = createRuleOptions.Version
	}
	if createRuleOptions.Import != nil {
		body["import"] = createRuleOptions.Import
	}
	if createRuleOptions.Labels != nil {
		body["labels"] = createRuleOptions.Labels
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
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

// GetRule : Get a custom rule
// Retrieve a rule that you created to evaluate your resources.  For more information, see [Defining custom
// rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetRule(getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetRuleWithContext(context.Background(), getRuleOptions)
}

// GetRuleWithContext is an alternate form of the GetRule method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetRuleWithContext(ctx context.Context, getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRuleOptions, "getRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRuleOptions, "getRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *getRuleOptions.InstanceID,
		"rule_id":     *getRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
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

// ReplaceRule : Update a custom rule
// Update a custom rule that you use to target the exact configuration properties  that you need to evaluate your
// resources for compliance. For more information, see [Defining custom
// rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceRule(replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ReplaceRuleWithContext(context.Background(), replaceRuleOptions)
}

// ReplaceRuleWithContext is an alternate form of the ReplaceRule method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ReplaceRuleWithContext(ctx context.Context, replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRuleOptions, "replaceRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceRuleOptions, "replaceRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *replaceRuleOptions.InstanceID,
		"rule_id":     *replaceRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ReplaceRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceRuleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceRuleOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replaceRuleOptions.Description != nil {
		body["description"] = replaceRuleOptions.Description
	}
	if replaceRuleOptions.Target != nil {
		body["target"] = replaceRuleOptions.Target
	}
	if replaceRuleOptions.RequiredConfig != nil {
		body["required_config"] = replaceRuleOptions.RequiredConfig
	}
	if replaceRuleOptions.Version != nil {
		body["version"] = replaceRuleOptions.Version
	}
	if replaceRuleOptions.Import != nil {
		body["import"] = replaceRuleOptions.Import
	}
	if replaceRuleOptions.Labels != nil {
		body["labels"] = replaceRuleOptions.Labels
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
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
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

// DeleteRule : Delete a custom rule
// Delete a custom rule that you no longer require to evaluate your resources. For more information, see [Defining
// custom rules](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-rules-define).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteRule(deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.DeleteRuleWithContext(context.Background(), deleteRuleOptions)
}

// DeleteRuleWithContext is an alternate form of the DeleteRule method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) DeleteRuleWithContext(ctx context.Context, deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRuleOptions, "deleteRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRuleOptions, "deleteRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"instance_id": *deleteRuleOptions.InstanceID,
		"rule_id":     *deleteRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/instances/{instance_id}/v3/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "DeleteRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = securityAndComplianceCenter.Service.Request(request, nil)

	return
}

// ListServices : List services
// List all the services that you use to evaluate the configuration of your resources for security and compliance.
// [Learn more](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scannable-components).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListServices(listServicesOptions *ListServicesOptions) (result *ServiceCollection, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.ListServicesWithContext(context.Background(), listServicesOptions)
}

// ListServicesWithContext is an alternate form of the ListServices method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) ListServicesWithContext(ctx context.Context, listServicesOptions *ListServicesOptions) (result *ServiceCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listServicesOptions, "listServicesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/v3/services`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listServicesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "ListServices")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceCollection)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetService : Get a service
// Retrieve a service configuration that you monitor. [Learn
// more](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scannable-components).
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetService(getServiceOptions *GetServiceOptions) (result *Service, response *core.DetailedResponse, err error) {
	return securityAndComplianceCenter.GetServiceWithContext(context.Background(), getServiceOptions)
}

// GetServiceWithContext is an alternate form of the GetService method which supports a Context parameter
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) GetServiceWithContext(ctx context.Context, getServiceOptions *GetServiceOptions) (result *Service, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getServiceOptions, "getServiceOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getServiceOptions, "getServiceOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"services_name": *getServiceOptions.ServicesName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = securityAndComplianceCenter.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(securityAndComplianceCenter.Service.Options.URL, `/v3/services/{services_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getServiceOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("security_and_compliance_center", "V3", "GetService")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = securityAndComplianceCenter.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalService)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// Account : The account that is associated with a report.
type Account struct {
	// The account ID.
	ID *string `json:"id,omitempty"`

	// The account name.
	Name *string `json:"name,omitempty"`

	// The account type.
	Type *string `json:"type,omitempty"`
}

// UnmarshalAccount unmarshals an instance of Account from the specified map of raw messages.
func UnmarshalAccount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Account)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// AdditionalProperty : AdditionalProperty struct
type AdditionalProperty struct {
	// An additional property that indicates the type of the attribute in various formats (text, url, secret, label,
	// masked).
	Type *string `json:"type" validate:"required"`

	// The name of the attribute that is displayed in the UI.
	DisplayName *string `json:"display_name" validate:"required"`
}

// Constants associated with the AdditionalProperty.Type property.
// An additional property that indicates the type of the attribute in various formats (text, url, secret, label,
// masked).
const (
	AdditionalPropertyTypeLabelConst  = "label"
	AdditionalPropertyTypeMaskedConst = "masked"
	AdditionalPropertyTypeSecretConst = "secret"
	AdditionalPropertyTypeTextConst   = "text"
	AdditionalPropertyTypeURLConst    = "url"
)

// UnmarshalAdditionalProperty unmarshals an instance of AdditionalProperty from the specified map of raw messages.
func UnmarshalAdditionalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalProperty)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalTargetAttribute : The additional target attribute of the service.
type AdditionalTargetAttribute struct {
	// The additional target attribute name.
	Name *string `json:"name,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the AdditionalTargetAttribute.Operator property.
// The operator.
const (
	AdditionalTargetAttributeOperatorDaysLessThanConst         = "days_less_than"
	AdditionalTargetAttributeOperatorIpsEqualsConst            = "ips_equals"
	AdditionalTargetAttributeOperatorIpsInRangeConst           = "ips_in_range"
	AdditionalTargetAttributeOperatorIpsNotEqualsConst         = "ips_not_equals"
	AdditionalTargetAttributeOperatorIsEmptyConst              = "is_empty"
	AdditionalTargetAttributeOperatorIsFalseConst              = "is_false"
	AdditionalTargetAttributeOperatorIsNotEmptyConst           = "is_not_empty"
	AdditionalTargetAttributeOperatorIsTrueConst               = "is_true"
	AdditionalTargetAttributeOperatorNumEqualsConst            = "num_equals"
	AdditionalTargetAttributeOperatorNumGreaterThanConst       = "num_greater_than"
	AdditionalTargetAttributeOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	AdditionalTargetAttributeOperatorNumLessThanConst          = "num_less_than"
	AdditionalTargetAttributeOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	AdditionalTargetAttributeOperatorNumNotEqualsConst         = "num_not_equals"
	AdditionalTargetAttributeOperatorStringContainsConst       = "string_contains"
	AdditionalTargetAttributeOperatorStringEqualsConst         = "string_equals"
	AdditionalTargetAttributeOperatorStringMatchConst          = "string_match"
	AdditionalTargetAttributeOperatorStringNotContainsConst    = "string_not_contains"
	AdditionalTargetAttributeOperatorStringNotEqualsConst      = "string_not_equals"
	AdditionalTargetAttributeOperatorStringNotMatchConst       = "string_not_match"
	AdditionalTargetAttributeOperatorStringsAllowedConst       = "strings_allowed"
	AdditionalTargetAttributeOperatorStringsInListConst        = "strings_in_list"
	AdditionalTargetAttributeOperatorStringsRequiredConst      = "strings_required"
)

// UnmarshalAdditionalTargetAttribute unmarshals an instance of AdditionalTargetAttribute from the specified map of raw messages.
func UnmarshalAdditionalTargetAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalTargetAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// Assessment : The control specification assessment.
type Assessment struct {
	// The assessment ID.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The assessment type.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The assessment method.
	AssessmentMethod *string `json:"assessment_method,omitempty"`

	// The assessment description.
	AssessmentDescription *string `json:"assessment_description,omitempty"`

	// The number of parameters of this assessment.
	ParameterCount *int64 `json:"parameter_count,omitempty"`

	// The list of parameters of this assessment.
	Parameters []Parameter `json:"parameters,omitempty"`
}

// UnmarshalAssessment unmarshals an instance of Assessment from the specified map of raw messages.
func UnmarshalAssessment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Assessment)
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_method", &obj.AssessmentMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_description", &obj.AssessmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_count", &obj.ParameterCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalParameter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AssessmentWithStats : The control specification assessment.
type AssessmentWithStats struct {
	// The assessment ID.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The assessment type.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The assessment method.
	AssessmentMethod *string `json:"assessment_method,omitempty"`

	// The assessment description.
	AssessmentDescription *string `json:"assessment_description,omitempty"`

	// The number of parameters of this assessment.
	ParameterCount *int64 `json:"parameter_count,omitempty"`

	// The list of parameters of this assessment.
	Parameters []Parameter `json:"parameters,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// UnmarshalAssessmentWithStats unmarshals an instance of AssessmentWithStats from the specified map of raw messages.
func UnmarshalAssessmentWithStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AssessmentWithStats)
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_method", &obj.AssessmentMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_description", &obj.AssessmentDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_count", &obj.ParameterCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalParameter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Attachment : The attachment that is associated with a report.
type Attachment struct {
	// The attachment ID.
	ID *string `json:"id,omitempty"`

	// The name of the attachment.
	Name *string `json:"name,omitempty"`

	// The description of the attachment.
	Description *string `json:"description,omitempty"`

	// The attachment schedule.
	Schedule *string `json:"schedule,omitempty"`

	// The report's scope from backwards compatiblity
	Scope Scope `json:"scope,omitempty"`

	// The report's scopes based on the caller's access permissions.
	Scopes []Scope `json:"scopes,omitempty"`

	// The notification configuration of the attachment.
	Notifications *AttachmentNotifications `json:"notifications,omitempty"`
}

// UnmarshalAttachment unmarshals an instance of Attachment from the specified map of raw messages.
func UnmarshalAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Attachment)
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
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scopes", &obj.Scopes, UnmarshalScope)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalAttachmentNotifications)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentNotifications : The notification configuration of the attachment.
type AttachmentNotifications struct {
	// Shows if the notification is enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The controls associated with an AttachmentNotification.
	Controls *AttachmentNotificationsControls `json:"controls,omitempty"`
}

// UnmarshalAttachmentNotifications unmarshals an instance of AttachmentNotifications from the specified map of raw messages.
func UnmarshalAttachmentNotifications(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentNotifications)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalAttachmentNotificationsControls)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AttachmentNotificationsControls : The controls associated with an AttachmentNotification.
type AttachmentNotificationsControls struct {
	// The maximum number of not compliant controls before a notification is triggered.
	ThresholdLimit *int64 `json:"threshold_limit,omitempty"`

	// List of controls that triggers a notification should a scan fail.
	FailedControlIds []string `json:"failed_control_ids,omitempty"`
}

// UnmarshalAttachmentNotificationsControls unmarshals an instance of AttachmentNotificationsControls from the specified map of raw messages.
func UnmarshalAttachmentNotificationsControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AttachmentNotificationsControls)
	err = core.UnmarshalPrimitive(m, "threshold_limit", &obj.ThresholdLimit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_control_ids", &obj.FailedControlIds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComparePredefinedProfilesResponse : The predefined profile comparison response.
type ComparePredefinedProfilesResponse struct {
	// Shows a change in the Profile.
	CurrentPredefinedVersion *CompareProfileResponse `json:"current_predefined_version,omitempty"`

	// Shows a change in the Profile.
	LatestPredefinedVersion *CompareProfileResponse `json:"latest_predefined_version,omitempty"`

	// Shows details of the controls that were changed.
	ControlsChanges *ControlChanges `json:"controls_changes,omitempty"`

	// Shows details of the parameters that were changed.
	DefaultParametersChanges *DefaultParametersChanges `json:"default_parameters_changes,omitempty"`
}

// UnmarshalComparePredefinedProfilesResponse unmarshals an instance of ComparePredefinedProfilesResponse from the specified map of raw messages.
func UnmarshalComparePredefinedProfilesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComparePredefinedProfilesResponse)
	err = core.UnmarshalModel(m, "current_predefined_version", &obj.CurrentPredefinedVersion, UnmarshalCompareProfileResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "latest_predefined_version", &obj.LatestPredefinedVersion, UnmarshalCompareProfileResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls_changes", &obj.ControlsChanges, UnmarshalControlChanges)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "default_parameters_changes", &obj.DefaultParametersChanges, UnmarshalDefaultParametersChanges)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CompareProfileResponse : Shows a change in the Profile.
type CompareProfileResponse struct {
	// The ID of the profile.
	ID *string `json:"id,omitempty"`

	// The name of the profile.
	ProfileName *string `json:"profile_name,omitempty"`

	// A description of what the profile should represent.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The type of profile, either predefined or custom.
	ProfileType *string `json:"profile_type,omitempty"`

	// The version of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// The unique identifier of the profile revision.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// Determines if the profile is up to date with the latest revisions.
	Latest *bool `json:"latest,omitempty"`

	// User who created the profile.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the profile was created, in date-time format.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// User who made the latest changes to the profile.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date when the profile was last updated, in date-time format.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// Number of controls in the profile.
	ControlsCount *int64 `json:"controls_count,omitempty"`
}

// Constants associated with the CompareProfileResponse.ProfileType property.
// The type of profile, either predefined or custom.
const (
	CompareProfileResponseProfileTypeCustomConst     = "custom"
	CompareProfileResponseProfileTypePredefinedConst = "predefined"
)

// UnmarshalCompareProfileResponse unmarshals an instance of CompareProfileResponse from the specified map of raw messages.
func UnmarshalCompareProfileResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CompareProfileResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_description", &obj.ProfileDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_type", &obj.ProfileType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_version", &obj.ProfileVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CompareProfilesOptions : The CompareProfiles options.
type CompareProfilesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCompareProfilesOptions : Instantiate CompareProfilesOptions
func (*SecurityAndComplianceCenterApiV3) NewCompareProfilesOptions(instanceID string, profileID string) *CompareProfilesOptions {
	return &CompareProfilesOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CompareProfilesOptions) SetInstanceID(instanceID string) *CompareProfilesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *CompareProfilesOptions) SetProfileID(profileID string) *CompareProfilesOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CompareProfilesOptions) SetAccountID(accountID string) *CompareProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CompareProfilesOptions) SetHeaders(param map[string]string) *CompareProfilesOptions {
	options.Headers = param
	return options
}

// ComplianceScore : The compliance score.
type ComplianceScore struct {
	// The number of successful evaluations.
	Passed *int64 `json:"passed,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The percentage of successful evaluations.
	Percent *int64 `json:"percent,omitempty"`
}

// UnmarshalComplianceScore unmarshals an instance of ComplianceScore from the specified map of raw messages.
func UnmarshalComplianceScore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceScore)
	err = core.UnmarshalPrimitive(m, "passed", &obj.Passed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "percent", &obj.Percent)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceStats : The compliance stats.
type ComplianceStats struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The number of not applicable (with no evaluations) checks.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`
}

// Constants associated with the ComplianceStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ComplianceStatsStatusCompliantConst              = "compliant"
	ComplianceStatsStatusNotApplicableConst          = "not_applicable"
	ComplianceStatsStatusNotCompliantConst           = "not_compliant"
	ComplianceStatsStatusUnableToPerformConst        = "unable_to_perform"
	ComplianceStatsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalComplianceStats unmarshals an instance of ComplianceStats from the specified map of raw messages.
func UnmarshalComplianceStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceStats)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComplianceStatsWithNonCompliant : The compliance stats.
type ComplianceStatsWithNonCompliant struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The number of not applicable (with no evaluations) checks.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`

	// The list of non compliant controls.
	NotCompliantControls []ControlSummary `json:"not_compliant_controls,omitempty"`
}

// Constants associated with the ComplianceStatsWithNonCompliant.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ComplianceStatsWithNonCompliantStatusCompliantConst              = "compliant"
	ComplianceStatsWithNonCompliantStatusNotApplicableConst          = "not_applicable"
	ComplianceStatsWithNonCompliantStatusNotCompliantConst           = "not_compliant"
	ComplianceStatsWithNonCompliantStatusUnableToPerformConst        = "unable_to_perform"
	ComplianceStatsWithNonCompliantStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalComplianceStatsWithNonCompliant unmarshals an instance of ComplianceStatsWithNonCompliant from the specified map of raw messages.
func UnmarshalComplianceStatsWithNonCompliant(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComplianceStatsWithNonCompliant)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "not_compliant_controls", &obj.NotCompliantControls, UnmarshalControlSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItem : ConditionItem struct
// Models which "extend" this model:
// - ConditionItemConditionBase
// - ConditionItemConditionList
// - ConditionItemConditionSubRule
type ConditionItem struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	Value interface{} `json:"value,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []ConditionItemIntf `json:"or,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []ConditionItemIntf `json:"and,omitempty"`

	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`

	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

// Constants associated with the ConditionItem.Operator property.
// The operator.
const (
	ConditionItemOperatorDaysLessThanConst         = "days_less_than"
	ConditionItemOperatorIpsEqualsConst            = "ips_equals"
	ConditionItemOperatorIpsInRangeConst           = "ips_in_range"
	ConditionItemOperatorIpsNotEqualsConst         = "ips_not_equals"
	ConditionItemOperatorIsEmptyConst              = "is_empty"
	ConditionItemOperatorIsFalseConst              = "is_false"
	ConditionItemOperatorIsNotEmptyConst           = "is_not_empty"
	ConditionItemOperatorIsTrueConst               = "is_true"
	ConditionItemOperatorNumEqualsConst            = "num_equals"
	ConditionItemOperatorNumGreaterThanConst       = "num_greater_than"
	ConditionItemOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	ConditionItemOperatorNumLessThanConst          = "num_less_than"
	ConditionItemOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	ConditionItemOperatorNumNotEqualsConst         = "num_not_equals"
	ConditionItemOperatorStringContainsConst       = "string_contains"
	ConditionItemOperatorStringEqualsConst         = "string_equals"
	ConditionItemOperatorStringMatchConst          = "string_match"
	ConditionItemOperatorStringNotContainsConst    = "string_not_contains"
	ConditionItemOperatorStringNotEqualsConst      = "string_not_equals"
	ConditionItemOperatorStringNotMatchConst       = "string_not_match"
	ConditionItemOperatorStringsAllowedConst       = "strings_allowed"
	ConditionItemOperatorStringsInListConst        = "strings_in_list"
	ConditionItemOperatorStringsRequiredConst      = "strings_required"
)

func (*ConditionItem) isaConditionItem() bool {
	return true
}

type ConditionItemIntf interface {
	isaConditionItem() bool
}

// UnmarshalConditionItem unmarshals an instance of ConditionItem from the specified map of raw messages.
func UnmarshalConditionItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItem)
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
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalConditionItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalConditionItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigurationInformationPoints : The service configuration information.
type ConfigurationInformationPoints struct {
	// The information type.
	Type *string `json:"type,omitempty"`

	// The service configurations endpoints.
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// UnmarshalConfigurationInformationPoints unmarshals an instance of ConfigurationInformationPoints from the specified map of raw messages.
func UnmarshalConfigurationInformationPoints(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigurationInformationPoints)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "endpoints", &obj.Endpoints, UnmarshalEndpoint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlChanges : Shows details of the controls that were changed.
type ControlChanges struct {
	// How many controls were added.
	TotalAdded *int64 `json:"total_added,omitempty"`

	// How many controls were removed.
	TotalRemoved *int64 `json:"total_removed,omitempty"`

	// How many controls were updated.
	TotalUpdated *int64 `json:"total_updated,omitempty"`

	// A list of controls that were added.
	Added []ProfileControlsInResponse `json:"added,omitempty"`

	// A list of controls that were removed.
	Removed []ProfileControlsInResponse `json:"removed,omitempty"`

	// A list of controls that were updated.
	Updated []ControlChangesUpdated `json:"updated,omitempty"`
}

// UnmarshalControlChanges unmarshals an instance of ControlChanges from the specified map of raw messages.
func UnmarshalControlChanges(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlChanges)
	err = core.UnmarshalPrimitive(m, "total_added", &obj.TotalAdded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_removed", &obj.TotalRemoved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_updated", &obj.TotalUpdated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "added", &obj.Added, UnmarshalProfileControlsInResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "removed", &obj.Removed, UnmarshalProfileControlsInResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "updated", &obj.Updated, UnmarshalControlChangesUpdated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlChangesUpdated : Shows the difference in the Controls.
type ControlChangesUpdated struct {
	// The control details for a profile.
	Current *ProfileControlsInResponse `json:"current,omitempty"`

	// The control details for a profile.
	Latest *ProfileControlsInResponse `json:"latest,omitempty"`
}

// UnmarshalControlChangesUpdated unmarshals an instance of ControlChangesUpdated from the specified map of raw messages.
func UnmarshalControlChangesUpdated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlChangesUpdated)
	err = core.UnmarshalModel(m, "current", &obj.Current, UnmarshalProfileControlsInResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "latest", &obj.Latest, UnmarshalProfileControlsInResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlDoc : References to a control documentation.
type ControlDoc struct {
	// The ID of the control doc.
	ControlDocsID *string `json:"control_docs_id,omitempty"`

	// The type of the control doc.
	ControlDocsType *string `json:"control_docs_type,omitempty"`
}

// UnmarshalControlDoc unmarshals an instance of ControlDoc from the specified map of raw messages.
func UnmarshalControlDoc(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlDoc)
	err = core.UnmarshalPrimitive(m, "control_docs_id", &obj.ControlDocsID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_docs_type", &obj.ControlDocsType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibrary : A Control Library.
type ControlLibrary struct {
	// The name of the control library.
	ControlLibraryName *string `json:"control_library_name,omitempty"`

	// Details of the control library.
	ControlLibraryDescription *string `json:"control_library_description,omitempty"`

	// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// The revision number of the control library.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The list of rules that the control library attempts to adhere to.
	Controls []Control `json:"controls,omitempty"`

	// The ID of the control library.
	ID *string `json:"id,omitempty"`

	// The ID of the account associated with the creation of the control library.
	AccountID *string `json:"account_id,omitempty"`

	// The ETag or version of the Control Library.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// Shows if the Control Library is the latest.
	Latest *bool `json:"latest,omitempty"`

	// The ID of the creator of the Control Library.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date-time of the creation.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The ID of the user who made the last update.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date-time of the update.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// Determines if the control library has any hierarchy.
	HierarchyEnabled *bool `json:"hierarchy_enabled,omitempty"`

	// The count of controls tied to the control library.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// THe count of control parents in the control library.
	ControlParentsCount *int64 `json:"control_parents_count,omitempty"`
}

// Constants associated with the ControlLibrary.ControlLibraryType property.
// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
const (
	ControlLibraryControlLibraryTypeCustomConst     = "custom"
	ControlLibraryControlLibraryTypePredefinedConst = "predefined"
)

// UnmarshalControlLibrary unmarshals an instance of ControlLibrary from the specified map of raw messages.
func UnmarshalControlLibrary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibrary)
	err = core.UnmarshalPrimitive(m, "control_library_name", &obj.ControlLibraryName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_description", &obj.ControlLibraryDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_type", &obj.ControlLibraryType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControl)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hierarchy_enabled", &obj.HierarchyEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parents_count", &obj.ControlParentsCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibraryCollection : A list of control libraries.
type ControlLibraryCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The list of control libraries.
	ControlLibraries []ControlLibrary `json:"control_libraries,omitempty"`
}

// UnmarshalControlLibraryCollection unmarshals an instance of ControlLibraryCollection from the specified map of raw messages.
func UnmarshalControlLibraryCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlLibraryCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_libraries", &obj.ControlLibraries, UnmarshalControlLibrary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Control : The assesment to abide to.
type Control struct {
	// The ID of the control library that contains the profile.
	ControlID *string `json:"control_id,omitempty"`

	// The Name of the control
	ControlName *string `json:"control_name,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The association of the control.
	ControlCategory *string `json:"control_category,omitempty"`

	// true if the control can be automated, false if the control cannot.
	ControlRequirement *bool `json:"control_requirement,omitempty"`

	// The ID of the parent control.
	ControlParent *string `json:"control_parent,omitempty"`

	// The path of the control
	ControlPath *string `json:"control_path,omitempty"`

	// Number of control specifications associated with the control.
	ControlSpecificationCount *int64 `json:"control_specification_count,omitempty"`

	// List of control specifications associated with the control.
	ControlSpecifications []ControlSpecification `json:"control_specifications,omitempty"`

	// List of Tags associated with the control
	ControlTags []string `json:"control_tags,omitempty"`

	// References to a control documentation.
	ControlDocs *ControlDoc `json:"control_docs,omitempty"`

	// Details if a control library can be used or not.
	Status *string `json:"status,omitempty"`
}

// UnmarshalControl unmarshals an instance of Control from the specified map of raw messages.
func UnmarshalControl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Control)
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_requirement", &obj.ControlRequirement)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parent", &obj.ControlParent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_path", &obj.ControlPath)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_count", &obj.ControlSpecificationCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecification)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_docs", &obj.ControlDocs, UnmarshalControlDoc)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_tags", &obj.ControlTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NewControl: Instantiate Control(Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewControl(controlName string, controlCategory string, controlRequirement bool, controlSpecifications []ControlSpecification) (_model *Control, err error) {
	_model = &Control{
		ControlName:           core.StringPtr(controlName),
		ControlCategory:       core.StringPtr(controlCategory),
		ControlRequirement:    core.BoolPtr(controlRequirement),
		ControlSpecifications: controlSpecifications,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// ControlSpecification : A statement that defines a security/privacy requirement for a Control.
type ControlSpecification struct {
	// The ID of the control.
	ID *string `json:"control_specification_id,omitempty"`

	// The Name of the control specification
	Name *string `json:"control_specification_name,omitempty"`

	// Details which party is responsible for the implementation of a specification.
	Responsibility *string `json:"responsibility,omitempty"`

	// The ID of the component.
	ComponentID *string `json:"component_id,omitempty"`

	// The name of the component.
	ComponentName *string `json:"component_name,omitempty"`

	// The type of the component.
	ComponentType *string `json:"component_type,omitempty"`

	// The cloud provider the specification is targeting.
	Environment *string `json:"environment,omitempty"`

	// Information about the Control Specification.
	Description *string `json:"control_specification_description,omitempty"`

	// The number of rules tied to the specification.
	AssessmentsCount *int64 `json:"assessments_count,omitempty"`

	// The detailed list of rules associated with the Specification.
	Assessments []Assessment `json:"assessments,omitempty"`
}

// UnmarshalControlSpecification unmarshals an instance of ControlSpecification from the specified map of raw messages.
func UnmarshalControlSpecification(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlSpecification)
	err = core.UnmarshalPrimitive(m, "control_specification_id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "responsibility", &obj.Responsibility)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_type", &obj.ComponentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessments_count", &obj.AssessmentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessments", &obj.Assessments, UnmarshalAssessment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlSpecificationWithStats : The control specification with compliance stats.
type ControlSpecificationWithStats struct {
	// The control specification ID.
	ControlSpecificationID *string `json:"control_specification_id,omitempty"`

	// The component description.
	ControlSpecificationDescription *string `json:"control_specification_description,omitempty"`

	// The component ID.
	ComponentID *string `json:"component_id,omitempty"`

	// The components name.
	ComponentName *string `json:"component_name,omitempty"`

	// The environment.
	Environment *string `json:"environment,omitempty"`

	// The responsibility for managing control specifications.
	Responsibility *string `json:"responsibility,omitempty"`

	// The list of assessments.
	Assessments []AssessmentWithStats `json:"assessments,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The number of not applicable (with no evaluations) checks.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`
}

// Constants associated with the ControlSpecificationWithStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ControlSpecificationWithStatsStatusCompliantConst              = "compliant"
	ControlSpecificationWithStatsStatusNotApplicableConst          = "not_applicable"
	ControlSpecificationWithStatsStatusNotCompliantConst           = "not_compliant"
	ControlSpecificationWithStatsStatusUnableToPerformConst        = "unable_to_perform"
	ControlSpecificationWithStatsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalControlSpecificationWithStats unmarshals an instance of ControlSpecificationWithStats from the specified map of raw messages.
func UnmarshalControlSpecificationWithStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlSpecificationWithStats)
	err = core.UnmarshalPrimitive(m, "control_specification_id", &obj.ControlSpecificationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specification_description", &obj.ControlSpecificationDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "responsibility", &obj.Responsibility)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessments", &obj.Assessments, UnmarshalAssessmentWithStats)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlSummary : The summary of the control.
type ControlSummary struct {
	// The controls ID.
	ID *string `json:"id,omitempty"`

	// The controls name.
	ControlName *string `json:"control_name,omitempty"`

	// The controls description.
	ControlDescription *string `json:"control_description,omitempty"`
}

// UnmarshalControlSummary unmarshals an instance of ControlSummary from the specified map of raw messages.
func UnmarshalControlSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlWithStats : The control with compliance stats.
type ControlWithStats struct {
	// The report ID.
	ReportID *string `json:"report_id,omitempty"`

	// The home account ID.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The control ID.
	ID *string `json:"id,omitempty"`

	// The control library ID.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The control name.
	ControlName *string `json:"control_name,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The control category.
	ControlCategory *string `json:"control_category,omitempty"`

	// The list of specifications that are on the page.
	ControlSpecifications []ControlSpecificationWithStats `json:"control_specifications,omitempty"`

	// The collection of different types of tags.
	ResourceTags *Tags `json:"resource_tags,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The number of not applicable (with no evaluations) checks.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`
}

// Constants associated with the ControlWithStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ControlWithStatsStatusCompliantConst              = "compliant"
	ControlWithStatsStatusNotApplicableConst          = "not_applicable"
	ControlWithStatsStatusNotCompliantConst           = "not_compliant"
	ControlWithStatsStatusUnableToPerformConst        = "unable_to_perform"
	ControlWithStatsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalControlWithStats unmarshals an instance of ControlWithStats from the specified map of raw messages.
func UnmarshalControlWithStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ControlWithStats)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecificationWithStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resource_tags", &obj.ResourceTags, UnmarshalTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCustomControlLibraryOptions : The CreateCustomControlLibrary options.
type CreateCustomControlLibraryOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The name of the control library.
	ControlLibraryName *string `json:"control_library_name" validate:"required"`

	// Details of the control library.
	ControlLibraryDescription *string `json:"control_library_description" validate:"required"`

	// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
	ControlLibraryType *string `json:"control_library_type" validate:"required"`

	// The revision number of the control library.
	ControlLibraryVersion *string `json:"control_library_version" validate:"required"`

	// The list of rules that the control library attempts to adhere to.
	Controls []Control `json:"controls" validate:"required"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The version group label
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// Set to determine if the latest is true
	Latest *bool `json:"latest,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateCustomControlLibraryOptions.ControlLibraryType property.
// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
const (
	CreateCustomControlLibraryOptionsControlLibraryTypeCustomConst = "custom"
)

// NewCreateCustomControlLibraryOptions : Instantiate CreateCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateCustomControlLibraryOptions(instanceID string, controlLibraryName string, controlLibraryDescription string, controlLibraryType string, controlLibraryVersion string, controls []Control) *CreateCustomControlLibraryOptions {
	return &CreateCustomControlLibraryOptions{
		InstanceID:                core.StringPtr(instanceID),
		ControlLibraryName:        core.StringPtr(controlLibraryName),
		ControlLibraryDescription: core.StringPtr(controlLibraryDescription),
		ControlLibraryType:        core.StringPtr(controlLibraryType),
		ControlLibraryVersion:     core.StringPtr(controlLibraryVersion),
		Controls:                  controls,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateCustomControlLibraryOptions) SetInstanceID(instanceID string) *CreateCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibraryName : Allow user to set ControlLibraryName
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryName(controlLibraryName string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryName = core.StringPtr(controlLibraryName)
	return _options
}

// SetControlLibraryDescription : Allow user to set ControlLibraryDescription
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryDescription(controlLibraryDescription string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryDescription = core.StringPtr(controlLibraryDescription)
	return _options
}

// SetControlLibraryType : Allow user to set ControlLibraryType
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryType(controlLibraryType string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryType = core.StringPtr(controlLibraryType)
	return _options
}

// SetControlLibraryVersion : Allow user to set ControlLibraryVersion
func (_options *CreateCustomControlLibraryOptions) SetControlLibraryVersion(controlLibraryVersion string) *CreateCustomControlLibraryOptions {
	_options.ControlLibraryVersion = core.StringPtr(controlLibraryVersion)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *CreateCustomControlLibraryOptions) SetControls(controls []Control) *CreateCustomControlLibraryOptions {
	_options.Controls = controls
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateCustomControlLibraryOptions) SetAccountID(accountID string) *CreateCustomControlLibraryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

func (_options *CreateCustomControlLibraryOptions) SetVersionGroupLabel(versionGroupLabel string) *CreateCustomControlLibraryOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

func (_options *CreateCustomControlLibraryOptions) SetLatest(latest bool) *CreateCustomControlLibraryOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCustomControlLibraryOptions) SetHeaders(param map[string]string) *CreateCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// CreateProfileAttachmentOptions : The CreateProfileAttachment options.
type CreateProfileAttachmentOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The Prototype to create a profile attachment.
	Attachments []ProfileAttachmentBase `json:"attachments,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProfileAttachmentOptions : Instantiate CreateProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateProfileAttachmentOptions(instanceID string, profileID string) *CreateProfileAttachmentOptions {
	return &CreateProfileAttachmentOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateProfileAttachmentOptions) SetInstanceID(instanceID string) *CreateProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *CreateProfileAttachmentOptions) SetProfileID(profileID string) *CreateProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachments : Allow user to set Attachments
func (_options *CreateProfileAttachmentOptions) SetAttachments(attachments []ProfileAttachmentBase) *CreateProfileAttachmentOptions {
	_options.Attachments = attachments
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateProfileAttachmentOptions) SetAccountID(accountID string) *CreateProfileAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileAttachmentOptions) SetHeaders(param map[string]string) *CreateProfileAttachmentOptions {
	options.Headers = param
	return options
}

// CreateProfileOptions : The CreateProfile options.
type CreateProfileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The name of the profile.
	ProfileName *string `json:"profile_name,omitempty"`

	// A description of what the profile should represent.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The version of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// Determines if the profile is up to date with the latest revisions.
	Latest *bool `json:"latest,omitempty"`

	// The unique identifier of the revision.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// List of controls associated with the profile.
	Controls []ProfileControlsPrototype `json:"controls,omitempty"`

	// The default values when using the profile.
	DefaultParameters []DefaultParametersPrototype `json:"default_parameters,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProfileOptions : Instantiate CreateProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateProfileOptions(instanceID string) *CreateProfileOptions {
	return &CreateProfileOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateProfileOptions) SetInstanceID(instanceID string) *CreateProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileName : Allow user to set ProfileName
func (_options *CreateProfileOptions) SetProfileName(profileName string) *CreateProfileOptions {
	_options.ProfileName = core.StringPtr(profileName)
	return _options
}

// SetProfileDescription : Allow user to set ProfileDescription
func (_options *CreateProfileOptions) SetProfileDescription(profileDescription string) *CreateProfileOptions {
	_options.ProfileDescription = core.StringPtr(profileDescription)
	return _options
}

// SetProfileVersion : Allow user to set ProfileVersion
func (_options *CreateProfileOptions) SetProfileVersion(profileVersion string) *CreateProfileOptions {
	_options.ProfileVersion = core.StringPtr(profileVersion)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *CreateProfileOptions) SetLatest(latest bool) *CreateProfileOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetVersionGroupLabel : Allow user to set VersionGroupLabel
func (_options *CreateProfileOptions) SetVersionGroupLabel(versionGroupLabel string) *CreateProfileOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *CreateProfileOptions) SetControls(controls []ProfileControlsPrototype) *CreateProfileOptions {
	_options.Controls = controls
	return _options
}

// SetDefaultParameters : Allow user to set DefaultParameters
func (_options *CreateProfileOptions) SetDefaultParameters(defaultParameters []DefaultParametersPrototype) *CreateProfileOptions {
	_options.DefaultParameters = defaultParameters
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateProfileOptions) SetAccountID(accountID string) *CreateProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProfileOptions) SetHeaders(param map[string]string) *CreateProfileOptions {
	options.Headers = param
	return options
}

// CreateProviderTypeInstanceOptions : The CreateProviderTypeInstance options.
type CreateProviderTypeInstanceOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance name.
	Name *string `json:"name" validate:"required,ne="`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateProviderTypeInstanceOptions : Instantiate CreateProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceName string) *CreateProviderTypeInstanceOptions {
	return &CreateProviderTypeInstanceOptions{
		InstanceID:     core.StringPtr(instanceID),
		ProviderTypeID: core.StringPtr(providerTypeID),
		Name:           core.StringPtr(providerTypeInstanceName),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateProviderTypeInstanceOptions) SetInstanceID(instanceID string) *CreateProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *CreateProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *CreateProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateProviderTypeInstanceOptions) SetName(name string) *CreateProviderTypeInstanceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAttributes : Allow user to set Attributes
func (_options *CreateProviderTypeInstanceOptions) SetAttributes(attributes map[string]interface{}) *CreateProviderTypeInstanceOptions {
	_options.Attributes = attributes
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProviderTypeInstanceOptions) SetHeaders(param map[string]string) *CreateProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// CreateRuleOptions : The CreateRule options.
type CreateRuleOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The rule description.
	Description *string `json:"description" validate:"required"`

	// The rule target.
	Target *RuleTargetPrototype `json:"target" validate:"required"`

	// The required configurations for a Rule.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The rule version number.
	Version *string `json:"version,omitempty"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The list of labels that correspond to a rule.
	Labels []string `json:"labels,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRuleOptions : Instantiate CreateRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateRuleOptions(instanceID string, description string, target *RuleTargetPrototype, requiredConfig RequiredConfigIntf) *CreateRuleOptions {
	return &CreateRuleOptions{
		InstanceID:     core.StringPtr(instanceID),
		Description:    core.StringPtr(description),
		Target:         target,
		RequiredConfig: requiredConfig,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateRuleOptions) SetInstanceID(instanceID string) *CreateRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateRuleOptions) SetDescription(description string) *CreateRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *CreateRuleOptions) SetTarget(target *RuleTargetPrototype) *CreateRuleOptions {
	_options.Target = target
	return _options
}

// SetRequiredConfig : Allow user to set RequiredConfig
func (_options *CreateRuleOptions) SetRequiredConfig(requiredConfig RequiredConfigIntf) *CreateRuleOptions {
	_options.RequiredConfig = requiredConfig
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CreateRuleOptions) SetVersion(version string) *CreateRuleOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetImport : Allow user to set Import
func (_options *CreateRuleOptions) SetImport(importVar *Import) *CreateRuleOptions {
	_options.Import = importVar
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *CreateRuleOptions) SetLabels(labels []string) *CreateRuleOptions {
	_options.Labels = labels
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRuleOptions) SetHeaders(param map[string]string) *CreateRuleOptions {
	options.Headers = param
	return options
}

// CreateScanOptions : The CreateScan options.
type CreateScanOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the profile attachment.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateScanOptions : Instantiate CreateScanOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateScanOptions(instanceID string) *CreateScanOptions {
	return &CreateScanOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateScanOptions) SetInstanceID(instanceID string) *CreateScanOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *CreateScanOptions) SetAttachmentID(attachmentID string) *CreateScanOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateScanOptions) SetAccountID(accountID string) *CreateScanOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateScanOptions) SetHeaders(param map[string]string) *CreateScanOptions {
	options.Headers = param
	return options
}

// CreateScanReport : The scan report ID.
type CreateScanReport struct {
	// The scan report ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalCreateScanReport unmarshals an instance of CreateScanReport from the specified map of raw messages.
func UnmarshalCreateScanReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateScanReport)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateScanReportOptions : The CreateScanReport options.
type CreateScanReportOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The enum of different report format types.
	Format *string `json:"format" validate:"required"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the sub-scope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateScanReportOptions.Format property.
// The enum of different report format types.
const (
	CreateScanReportOptionsFormatCSVConst = "csv"
	CreateScanReportOptionsFormatPDFConst = "pdf"
)

// NewCreateScanReportOptions : Instantiate CreateScanReportOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateScanReportOptions(instanceID string, reportID string, format string) *CreateScanReportOptions {
	return &CreateScanReportOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
		Format:     core.StringPtr(format),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateScanReportOptions) SetInstanceID(instanceID string) *CreateScanReportOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *CreateScanReportOptions) SetReportID(reportID string) *CreateScanReportOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *CreateScanReportOptions) SetFormat(format string) *CreateScanReportOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *CreateScanReportOptions) SetScopeID(scopeID string) *CreateScanReportOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *CreateScanReportOptions) SetSubscopeID(subscopeID string) *CreateScanReportOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateScanReportOptions) SetHeaders(param map[string]string) *CreateScanReportOptions {
	options.Headers = param
	return options
}

// CreateScanResponse : The response that details the whether starting a scan was successful.
type CreateScanResponse struct {
	// The ID of the scan generated.
	ID *string `json:"id,omitempty"`

	// The ID of the account associated with the scan.
	AccountID *string `json:"account_id,omitempty"`

	// The ID of the profile attachment associated with the scan.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The ID of the report associated with the scan.
	ReportID *string `json:"report_id,omitempty"`

	// Details the state of a scan.
	Status *string `json:"status,omitempty"`

	// The last time a scan was performed.
	LastScanTime *strfmt.DateTime `json:"last_scan_time,omitempty"`

	// The next time a scan will be triggered.
	NextScanTime *strfmt.DateTime `json:"next_scan_time,omitempty"`

	// Shows how a scan gets triggered.
	ScanType *string `json:"scan_type,omitempty"`

	// The number of times the scan appeared.
	Occurence *int64 `json:"occurence,omitempty"`
}

// UnmarshalCreateScanResponse unmarshals an instance of CreateScanResponse from the specified map of raw messages.
func UnmarshalCreateScanResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateScanResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachment_id", &obj.AttachmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_scan_time", &obj.LastScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_scan_time", &obj.NextScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_type", &obj.ScanType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "occurence", &obj.Occurence)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateScopeOptions : The CreateScope options.
type CreateScopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The scope name.
	Name *string `json:"name,omitempty"`

	// The scope description.
	Description *string `json:"description,omitempty"`

	// The scope environment.
	Environment *string `json:"environment,omitempty"`

	// The properties that are supported for scoping by this environment.
	Properties []ScopePropertyIntf `json:"properties,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateScopeOptions : Instantiate CreateScopeOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateScopeOptions(instanceID string) *CreateScopeOptions {
	return &CreateScopeOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateScopeOptions) SetInstanceID(instanceID string) *CreateScopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateScopeOptions) SetName(name string) *CreateScopeOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateScopeOptions) SetDescription(description string) *CreateScopeOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnvironment : Allow user to set Environment
func (_options *CreateScopeOptions) SetEnvironment(environment string) *CreateScopeOptions {
	_options.Environment = core.StringPtr(environment)
	return _options
}

// SetProperties : Allow user to set Properties
func (_options *CreateScopeOptions) SetProperties(properties []ScopePropertyIntf) *CreateScopeOptions {
	_options.Properties = properties
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateScopeOptions) SetHeaders(param map[string]string) *CreateScopeOptions {
	options.Headers = param
	return options
}

// CreateSubscopeOptions : The CreateSubscope options.
type CreateSubscopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The array of subscopes.
	Subscopes []ScopePrototype `json:"subscopes,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSubscopeOptions : Instantiate CreateSubscopeOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateSubscopeOptions(instanceID string, scopeID string) *CreateSubscopeOptions {
	return &CreateSubscopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateSubscopeOptions) SetInstanceID(instanceID string) *CreateSubscopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *CreateSubscopeOptions) SetScopeID(scopeID string) *CreateSubscopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopes : Allow user to set Subscopes
func (_options *CreateSubscopeOptions) SetSubscopes(subscopes []ScopePrototype) *CreateSubscopeOptions {
	_options.Subscopes = subscopes
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSubscopeOptions) SetHeaders(param map[string]string) *CreateSubscopeOptions {
	options.Headers = param
	return options
}

// CreateTargetOptions : The CreateTarget options.
type CreateTargetOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The target account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The trusted profile ID.
	TrustedProfileID *string `json:"trusted_profile_id" validate:"required"`

	// The target name.
	Name *string `json:"name" validate:"required"`

	// Customer credential to access for a specific service to scan.
	Credentials []Credential `json:"credentials,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTargetOptions : Instantiate CreateTargetOptions
func (*SecurityAndComplianceCenterApiV3) NewCreateTargetOptions(instanceID string, accountID string, trustedProfileID string, name string) *CreateTargetOptions {
	return &CreateTargetOptions{
		InstanceID:       core.StringPtr(instanceID),
		AccountID:        core.StringPtr(accountID),
		TrustedProfileID: core.StringPtr(trustedProfileID),
		Name:             core.StringPtr(name),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *CreateTargetOptions) SetInstanceID(instanceID string) *CreateTargetOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateTargetOptions) SetAccountID(accountID string) *CreateTargetOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTrustedProfileID : Allow user to set TrustedProfileID
func (_options *CreateTargetOptions) SetTrustedProfileID(trustedProfileID string) *CreateTargetOptions {
	_options.TrustedProfileID = core.StringPtr(trustedProfileID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateTargetOptions) SetName(name string) *CreateTargetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCredentials : Allow user to set Credentials
func (_options *CreateTargetOptions) SetCredentials(credentials []Credential) *CreateTargetOptions {
	_options.Credentials = credentials
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTargetOptions) SetHeaders(param map[string]string) *CreateTargetOptions {
	options.Headers = param
	return options
}

// Credential : Credential struct
type Credential struct {
	// The CRN of the secret.
	SecretCRN *string `json:"secret_crn" validate:"required"`

	// The type of the credential.
	Type *string `json:"type,omitempty"`

	// The type of credentials
	SecretName *string `json:"secret_name,omitempty"`

	// Credential having service name and instance crn.
	Resources []Resource `json:"resources" validate:"required"`
}

// NewCredential : Instantiate Credential (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewCredential(secretCRN string, resources []Resource) (_model *Credential, err error) {
	_model = &Credential{
		SecretCRN: core.StringPtr(secretCRN),
		Resources: resources,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCredential unmarshals an instance of Credential from the specified map of raw messages.
func UnmarshalCredential(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Credential)
	err = core.UnmarshalPrimitive(m, "secret_crn", &obj.SecretCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CredentialResponse : CredentialResponse struct
type CredentialResponse struct {
	// The type of the credential.
	Type *string `json:"type" validate:"required"`

	// The CRN of the secret.
	SecretCRN *string `json:"secret_crn" validate:"required"`

	// The name of the secret.
	SecretName *string `json:"secret_name,omitempty"`

	// Credential having service name and instance crn.
	Resources []Resource `json:"resources" validate:"required"`
}

// UnmarshalCredentialResponse unmarshals an instance of CredentialResponse from the specified map of raw messages.
func UnmarshalCredentialResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CredentialResponse)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_crn", &obj.SecretCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefaultParameters : The parameters of the profile that are inherently set by the profile.
type DefaultParameters struct {
	// The type of the implementation.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The implementation ID of the parameter.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The parameter's name.
	ParameterName *string `json:"parameter_name,omitempty"`

	ParameterDefaultValue interface{} `json:"parameter_default_value,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`
}

// UnmarshalDefaultParameters unmarshals an instance of DefaultParameters from the specified map of raw messages.
func UnmarshalDefaultParameters(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefaultParameters)
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_default_value", &obj.ParameterDefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefaultParametersChanges : Shows details of the parameters that were changed.
type DefaultParametersChanges struct {
	// Number of parameters added.
	TotalAdded *int64 `json:"total_added,omitempty"`

	// Number of parameters removed.
	TotalRemoved *int64 `json:"total_removed,omitempty"`

	// Number of parameters updated.
	TotalUpdated *int64 `json:"total_updated,omitempty"`

	// List of parameters added.
	Added []DefaultParameters `json:"added,omitempty"`

	// List of parameters removed.
	Removed []DefaultParameters `json:"removed,omitempty"`

	// List of parameters updated.
	Updated []DefaultParametersDifference `json:"updated,omitempty"`
}

// UnmarshalDefaultParametersChanges unmarshals an instance of DefaultParametersChanges from the specified map of raw messages.
func UnmarshalDefaultParametersChanges(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefaultParametersChanges)
	err = core.UnmarshalPrimitive(m, "total_added", &obj.TotalAdded)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_removed", &obj.TotalRemoved)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_updated", &obj.TotalUpdated)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "added", &obj.Added, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "removed", &obj.Removed, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "updated", &obj.Updated, UnmarshalDefaultParametersDifference)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefaultParametersDifference : Details the difference between the current parameters compared to the latest.
type DefaultParametersDifference struct {
	// The parameters of the profile that are inherently set by the profile.
	Current *DefaultParameters `json:"current,omitempty"`

	// The parameters of the profile that are inherently set by the profile.
	Latest *DefaultParameters `json:"latest,omitempty"`
}

// UnmarshalDefaultParametersDifference unmarshals an instance of DefaultParametersDifference from the specified map of raw messages.
func UnmarshalDefaultParametersDifference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefaultParametersDifference)
	err = core.UnmarshalModel(m, "current", &obj.Current, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "latest", &obj.Latest, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefaultParametersPrototype : The parameters of the profile to inherently set.
type DefaultParametersPrototype struct {
	// The type of the implementation.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The implementation ID of the parameter.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The parameter's name.
	ParameterName *string `json:"parameter_name,omitempty"`

	ParameterDefaultValue interface{} `json:"parameter_default_value,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`
}

// UnmarshalDefaultParametersPrototype unmarshals an instance of DefaultParametersPrototype from the specified map of raw messages.
func UnmarshalDefaultParametersPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefaultParametersPrototype)
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_default_value", &obj.ParameterDefaultValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteCustomControlLibraryOptions : The DeleteCustomControlLibrary options.
type DeleteCustomControlLibraryOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library ID.
	ControlLibraryID *string `json:"control_library_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCustomControlLibraryOptions : Instantiate DeleteCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteCustomControlLibraryOptions(instanceID string, controlLibraryID string) *DeleteCustomControlLibraryOptions {
	return &DeleteCustomControlLibraryOptions{
		InstanceID:       core.StringPtr(instanceID),
		ControlLibraryID: core.StringPtr(controlLibraryID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomControlLibraryOptions) SetInstanceID(instanceID string) *DeleteCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibraryID : Allow user to set ControlLibraryID
func (_options *DeleteCustomControlLibraryOptions) SetControlLibraryID(controlLibraryID string) *DeleteCustomControlLibraryOptions {
	_options.ControlLibraryID = core.StringPtr(controlLibraryID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteCustomControlLibraryOptions) SetAccountID(accountID string) *DeleteCustomControlLibraryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomControlLibraryOptions) SetHeaders(param map[string]string) *DeleteCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// DeleteCustomProfileOptions : The DeleteCustomProfile options.
type DeleteCustomProfileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCustomProfileOptions : Instantiate DeleteCustomProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteCustomProfileOptions(instanceID string, profileID string) *DeleteCustomProfileOptions {
	return &DeleteCustomProfileOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteCustomProfileOptions) SetInstanceID(instanceID string) *DeleteCustomProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteCustomProfileOptions) SetProfileID(profileID string) *DeleteCustomProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteCustomProfileOptions) SetAccountID(accountID string) *DeleteCustomProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomProfileOptions) SetHeaders(param map[string]string) *DeleteCustomProfileOptions {
	options.Headers = param
	return options
}

// DeleteProfileAttachmentOptions : The DeleteProfileAttachment options.
type DeleteProfileAttachmentOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProfileAttachmentOptions : Instantiate DeleteProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteProfileAttachmentOptions(instanceID string, profileID string, attachmentID string) *DeleteProfileAttachmentOptions {
	return &DeleteProfileAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		ProfileID:    core.StringPtr(profileID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteProfileAttachmentOptions) SetInstanceID(instanceID string) *DeleteProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *DeleteProfileAttachmentOptions) SetProfileID(profileID string) *DeleteProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *DeleteProfileAttachmentOptions) SetAttachmentID(attachmentID string) *DeleteProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *DeleteProfileAttachmentOptions) SetAccountID(accountID string) *DeleteProfileAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProfileAttachmentOptions) SetHeaders(param map[string]string) *DeleteProfileAttachmentOptions {
	options.Headers = param
	return options
}

// DeleteProviderTypeInstanceOptions : The DeleteProviderTypeInstance options.
type DeleteProviderTypeInstanceOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteProviderTypeInstanceOptions : Instantiate DeleteProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string) *DeleteProviderTypeInstanceOptions {
	return &DeleteProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteProviderTypeInstanceOptions) SetInstanceID(instanceID string) *DeleteProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *DeleteProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *DeleteProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *DeleteProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *DeleteProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProviderTypeInstanceOptions) SetHeaders(param map[string]string) *DeleteProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// DeleteRuleOptions : The DeleteRule options.
type DeleteRuleOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of a rule/assessment.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRuleOptions : Instantiate DeleteRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteRuleOptions(instanceID string, ruleID string) *DeleteRuleOptions {
	return &DeleteRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteRuleOptions) SetInstanceID(instanceID string) *DeleteRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteRuleOptions) SetRuleID(ruleID string) *DeleteRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRuleOptions) SetHeaders(param map[string]string) *DeleteRuleOptions {
	options.Headers = param
	return options
}

// DeleteScopeOptions : The DeleteScope options.
type DeleteScopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteScopeOptions : Instantiate DeleteScopeOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteScopeOptions(instanceID string, scopeID string) *DeleteScopeOptions {
	return &DeleteScopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteScopeOptions) SetInstanceID(instanceID string) *DeleteScopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *DeleteScopeOptions) SetScopeID(scopeID string) *DeleteScopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteScopeOptions) SetHeaders(param map[string]string) *DeleteScopeOptions {
	options.Headers = param
	return options
}

// DeleteSubscopeOptions : The DeleteSubscope options.
type DeleteSubscopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	SubscopeID *string `json:"subscope_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSubscopeOptions : Instantiate DeleteSubscopeOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteSubscopeOptions(instanceID string, scopeID string, subscopeID string) *DeleteSubscopeOptions {
	return &DeleteSubscopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
		SubscopeID: core.StringPtr(subscopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteSubscopeOptions) SetInstanceID(instanceID string) *DeleteSubscopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *DeleteSubscopeOptions) SetScopeID(scopeID string) *DeleteSubscopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *DeleteSubscopeOptions) SetSubscopeID(subscopeID string) *DeleteSubscopeOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSubscopeOptions) SetHeaders(param map[string]string) *DeleteSubscopeOptions {
	options.Headers = param
	return options
}

// DeleteTargetOptions : The DeleteTarget options.
type DeleteTargetOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The target ID.
	TargetID *string `json:"target_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteTargetOptions : Instantiate DeleteTargetOptions
func (*SecurityAndComplianceCenterApiV3) NewDeleteTargetOptions(instanceID string, targetID string) *DeleteTargetOptions {
	return &DeleteTargetOptions{
		InstanceID: core.StringPtr(instanceID),
		TargetID:   core.StringPtr(targetID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *DeleteTargetOptions) SetInstanceID(instanceID string) *DeleteTargetOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetTargetID : Allow user to set TargetID
func (_options *DeleteTargetOptions) SetTargetID(targetID string) *DeleteTargetOptions {
	_options.TargetID = core.StringPtr(targetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTargetOptions) SetHeaders(param map[string]string) *DeleteTargetOptions {
	options.Headers = param
	return options
}

// Endpoint : The service configurations endpoint.
type Endpoint struct {
	// The endpoint host.
	Host *string `json:"host,omitempty"`

	// The endpoint path.
	Path *string `json:"path,omitempty"`

	// The endpoint region.
	Region *string `json:"region,omitempty"`

	// The endpoints advisory call limit.
	AdvisoryCallLimit *int64 `json:"advisory_call_limit,omitempty"`
}

// UnmarshalEndpoint unmarshals an instance of Endpoint from the specified map of raw messages.
func UnmarshalEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Endpoint)
	err = core.UnmarshalPrimitive(m, "host", &obj.Host)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "advisory_call_limit", &obj.AdvisoryCallLimit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvalStats : The evaluation stats.
type EvalStats struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The number of assessments with no corresponding evaluations.
	SkippedCount *int64 `json:"skipped_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// Constants associated with the EvalStats.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	EvalStatsStatusCompliantConst              = "compliant"
	EvalStatsStatusNotApplicableConst          = "not_applicable"
	EvalStatsStatusNotCompliantConst           = "not_compliant"
	EvalStatsStatusUnableToPerformConst        = "unable_to_perform"
	EvalStatsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalEvalStats unmarshals an instance of EvalStats from the specified map of raw messages.
func UnmarshalEvalStats(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvalStats)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped_count", &obj.SkippedCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Evaluation : The evaluation of a control specification assessment.
type Evaluation struct {
	// The ID of the report that is associated to the evaluation.
	ReportID *string `json:"report_id,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The component ID.
	ComponentID *string `json:"component_id,omitempty"`

	// The components name.
	ComponentName *string `json:"component_name,omitempty"`

	// The control specification assessment.
	Assessment *Assessment `json:"assessment,omitempty"`

	// The time when the evaluation was made.
	EvaluateTime *string `json:"evaluate_time,omitempty"`

	// The evaluation target.
	Target *TargetInfo `json:"target,omitempty"`

	// The allowed values of an evaluation status.
	Status *string `json:"status,omitempty"`

	// The reason for the evaluation failure.
	Reason *string `json:"reason,omitempty"`

	// A list of details related to the Evaluation.
	Details *EvaluationDetails `json:"details,omitempty"`

	// By whom the evaluation was made for erictree results.
	EvaluatedBy *string `json:"evaluated_by,omitempty"`
}

// Constants associated with the Evaluation.Status property.
// The allowed values of an evaluation status.
const (
	EvaluationStatusErrorConst   = "error"
	EvaluationStatusFailureConst = "failure"
	EvaluationStatusPassConst    = "pass"
	EvaluationStatusSkippedConst = "skipped"
)

// UnmarshalEvaluation unmarshals an instance of Evaluation from the specified map of raw messages.
func UnmarshalEvaluation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Evaluation)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessment", &obj.Assessment, UnmarshalAssessment)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluate_time", &obj.EvaluateTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalTargetInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "details", &obj.Details, UnmarshalEvaluationDetails)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluated_by", &obj.EvaluatedBy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvaluationDetails : A list of details related to the Evaluation.
type EvaluationDetails struct {
	// Details the evaluations that were incorrect.
	Properties []EvaluationProperty `json:"properties,omitempty"`

	// The source provider of the evaluation.
	ProviderInfo *EvaluationProviderInfo `json:"provider_info,omitempty"`
}

// UnmarshalEvaluationDetails unmarshals an instance of EvaluationDetails from the specified map of raw messages.
func UnmarshalEvaluationDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvaluationDetails)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalEvaluationProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "provider_info", &obj.ProviderInfo, UnmarshalEvaluationProviderInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvaluationPage : The page of assessment evaluations.
type EvaluationPage struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The list of evaluations that are on the page.
	Evaluations []Evaluation `json:"evaluations,omitempty"`
}

// UnmarshalEvaluationPage unmarshals an instance of EvaluationPage from the specified map of raw messages.
func UnmarshalEvaluationPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvaluationPage)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations", &obj.Evaluations, UnmarshalEvaluation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *EvaluationPage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// EvaluationProperty : An aspect of the evaluation.
type EvaluationProperty struct {
	// The attribute of the resource.
	Property *string `json:"property,omitempty"`

	// An explanation of the resourcer.
	PropertyDescription *string `json:"property_description,omitempty"`

	// The operator used during the evaluation.
	Operator *string `json:"operator,omitempty"`

	ExpectedValue interface{} `json:"expected_value,omitempty"`

	FoundValue interface{} `json:"found_value,omitempty"`
}

// Constants associated with the EvaluationProperty.Operator property.
// The operator used during the evaluation.
const (
	EvaluationPropertyOperatorDaysLessThanConst         = "days_less_than"
	EvaluationPropertyOperatorIpsEqualsConst            = "ips_equals"
	EvaluationPropertyOperatorIpsInRangeConst           = "ips_in_range"
	EvaluationPropertyOperatorIpsNotEqualsConst         = "ips_not_equals"
	EvaluationPropertyOperatorIsEmptyConst              = "is_empty"
	EvaluationPropertyOperatorIsFalseConst              = "is_false"
	EvaluationPropertyOperatorIsNotEmptyConst           = "is_not_empty"
	EvaluationPropertyOperatorIsTrueConst               = "is_true"
	EvaluationPropertyOperatorNumEqualsConst            = "num_equals"
	EvaluationPropertyOperatorNumGreaterThanConst       = "num_greater_than"
	EvaluationPropertyOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	EvaluationPropertyOperatorNumLessThanConst          = "num_less_than"
	EvaluationPropertyOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	EvaluationPropertyOperatorNumNotEqualsConst         = "num_not_equals"
	EvaluationPropertyOperatorStringContainsConst       = "string_contains"
	EvaluationPropertyOperatorStringEqualsConst         = "string_equals"
	EvaluationPropertyOperatorStringMatchConst          = "string_match"
	EvaluationPropertyOperatorStringNotContainsConst    = "string_not_contains"
	EvaluationPropertyOperatorStringNotEqualsConst      = "string_not_equals"
	EvaluationPropertyOperatorStringNotMatchConst       = "string_not_match"
	EvaluationPropertyOperatorStringsAllowedConst       = "strings_allowed"
	EvaluationPropertyOperatorStringsInListConst        = "strings_in_list"
	EvaluationPropertyOperatorStringsRequiredConst      = "strings_required"
)

// UnmarshalEvaluationProperty unmarshals an instance of EvaluationProperty from the specified map of raw messages.
func UnmarshalEvaluationProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvaluationProperty)
	err = core.UnmarshalPrimitive(m, "property", &obj.Property)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property_description", &obj.PropertyDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_value", &obj.ExpectedValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "found_value", &obj.FoundValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EvaluationProviderInfo : The source provider of the evaluation.
type EvaluationProviderInfo struct {
	// Details the source of the evaluation.
	ProviderType *string `json:"provider_type,omitempty"`
}

// UnmarshalEvaluationProviderInfo unmarshals an instance of EvaluationProviderInfo from the specified map of raw messages.
func UnmarshalEvaluationProviderInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EvaluationProviderInfo)
	err = core.UnmarshalPrimitive(m, "provider_type", &obj.ProviderType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EventNotifications : The Event Notifications settings.
type EventNotifications struct {
	// The Event Notifications instance CRN.
	InstanceCRN *string `json:"instance_crn,omitempty"`

	// The date when the Event Notifications connection was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The connected Security and Compliance Center instance CRN.
	SourceID *string `json:"source_id,omitempty"`

	// The description of the source of the Event Notifications.
	SourceDescription *string `json:"source_description,omitempty"`

	// The name of the source of the Event Notifications.
	SourceName *string `json:"source_name,omitempty"`
}

// UnmarshalEventNotifications unmarshals an instance of EventNotifications from the specified map of raw messages.
func UnmarshalEventNotifications(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventNotifications)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_id", &obj.SourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_description", &obj.SourceDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_name", &obj.SourceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EventNotificationsPrototype : The payload to connect an Event Notification instance with a Security and Compliance Center instance.
type EventNotificationsPrototype struct {
	// The CRN of the Event Notification instance to connect.
	InstanceCRN *string `json:"instance_crn,omitempty"`

	// The description of the source of the Event Notifications.
	SourceDescription *string `json:"source_description,omitempty"`

	// The name of the source of the Event Notifications.
	SourceName *string `json:"source_name,omitempty"`
}

// UnmarshalEventNotificationsPrototype unmarshals an instance of EventNotificationsPrototype from the specified map of raw messages.
func UnmarshalEventNotificationsPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventNotificationsPrototype)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_description", &obj.SourceDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_name", &obj.SourceName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetControlLibraryOptions : The GetControlLibrary options.
type GetControlLibraryOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library ID.
	ControlLibraryID *string `json:"control_library_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetControlLibraryOptions : Instantiate GetControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewGetControlLibraryOptions(instanceID string, controlLibraryID string) *GetControlLibraryOptions {
	return &GetControlLibraryOptions{
		InstanceID:       core.StringPtr(instanceID),
		ControlLibraryID: core.StringPtr(controlLibraryID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetControlLibraryOptions) SetInstanceID(instanceID string) *GetControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibraryID : Allow user to set ControlLibraryID
func (_options *GetControlLibraryOptions) SetControlLibraryID(controlLibraryID string) *GetControlLibraryOptions {
	_options.ControlLibraryID = core.StringPtr(controlLibraryID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetControlLibraryOptions) SetAccountID(accountID string) *GetControlLibraryOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetControlLibraryOptions) SetHeaders(param map[string]string) *GetControlLibraryOptions {
	options.Headers = param
	return options
}

// GetLatestReportsOptions : The GetLatestReports options.
type GetLatestReportsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// This field sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLatestReportsOptions : Instantiate GetLatestReportsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetLatestReportsOptions(instanceID string) *GetLatestReportsOptions {
	return &GetLatestReportsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetLatestReportsOptions) SetInstanceID(instanceID string) *GetLatestReportsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *GetLatestReportsOptions) SetSort(sort string) *GetLatestReportsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLatestReportsOptions) SetHeaders(param map[string]string) *GetLatestReportsOptions {
	options.Headers = param
	return options
}

// GetProfileAttachmentOptions : The GetProfileAttachment options.
type GetProfileAttachmentOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProfileAttachmentOptions : Instantiate GetProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProfileAttachmentOptions(instanceID string, profileID string, attachmentID string) *GetProfileAttachmentOptions {
	return &GetProfileAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		ProfileID:    core.StringPtr(profileID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProfileAttachmentOptions) SetInstanceID(instanceID string) *GetProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileAttachmentOptions) SetProfileID(profileID string) *GetProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *GetProfileAttachmentOptions) SetAttachmentID(attachmentID string) *GetProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetProfileAttachmentOptions) SetAccountID(accountID string) *GetProfileAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileAttachmentOptions) SetHeaders(param map[string]string) *GetProfileAttachmentOptions {
	options.Headers = param
	return options
}

// GetProfileOptions : The GetProfile options.
type GetProfileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProfileOptions : Instantiate GetProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProfileOptions(instanceID string, profileID string) *GetProfileOptions {
	return &GetProfileOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProfileOptions) SetInstanceID(instanceID string) *GetProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *GetProfileOptions) SetProfileID(profileID string) *GetProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *GetProfileOptions) SetAccountID(accountID string) *GetProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProfileOptions) SetHeaders(param map[string]string) *GetProfileOptions {
	options.Headers = param
	return options
}

// GetProviderTypeByIDOptions : The GetProviderTypeByID options.
type GetProviderTypeByIDOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProviderTypeByIDOptions : Instantiate GetProviderTypeByIDOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProviderTypeByIDOptions(instanceID string, providerTypeID string) *GetProviderTypeByIDOptions {
	return &GetProviderTypeByIDOptions{
		InstanceID:     core.StringPtr(instanceID),
		ProviderTypeID: core.StringPtr(providerTypeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProviderTypeByIDOptions) SetInstanceID(instanceID string) *GetProviderTypeByIDOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *GetProviderTypeByIDOptions) SetProviderTypeID(providerTypeID string) *GetProviderTypeByIDOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProviderTypeByIDOptions) SetHeaders(param map[string]string) *GetProviderTypeByIDOptions {
	options.Headers = param
	return options
}

// GetProviderTypeInstanceOptions : The GetProviderTypeInstance options.
type GetProviderTypeInstanceOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetProviderTypeInstanceOptions : Instantiate GetProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewGetProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string) *GetProviderTypeInstanceOptions {
	return &GetProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetProviderTypeInstanceOptions) SetInstanceID(instanceID string) *GetProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *GetProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *GetProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *GetProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *GetProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProviderTypeInstanceOptions) SetHeaders(param map[string]string) *GetProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// GetReportControlsOptions : The GetReportControls options.
type GetReportControlsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the control.
	ControlID *string `json:"control_id,omitempty"`

	// The name of the control.
	ControlName *string `json:"control_name,omitempty"`

	// The description of the control.
	ControlDescription *string `json:"control_description,omitempty"`

	// A control category value.
	ControlCategory *string `json:"control_category,omitempty"`

	// The compliance status value.
	Status *string `json:"status,omitempty"`

	// This field sorts controls by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetReportControlsOptions.Status property.
// The compliance status value.
const (
	GetReportControlsOptionsStatusCompliantConst              = "compliant"
	GetReportControlsOptionsStatusNotApplicableConst          = "not_applicable"
	GetReportControlsOptionsStatusNotCompliantConst           = "not_compliant"
	GetReportControlsOptionsStatusUnableToPerformConst        = "unable_to_perform"
	GetReportControlsOptionsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// Constants associated with the GetReportControlsOptions.Sort property.
// This field sorts controls by using a valid sort field. To learn more, see
// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
const (
	GetReportControlsOptionsSortControlCategoryConst = "control_category"
	GetReportControlsOptionsSortControlNameConst     = "control_name"
	GetReportControlsOptionsSortStatusConst          = "status"
)

// NewGetReportControlsOptions : Instantiate GetReportControlsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportControlsOptions(instanceID string, reportID string) *GetReportControlsOptions {
	return &GetReportControlsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportControlsOptions) SetInstanceID(instanceID string) *GetReportControlsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportControlsOptions) SetReportID(reportID string) *GetReportControlsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetControlID : Allow user to set ControlID
func (_options *GetReportControlsOptions) SetControlID(controlID string) *GetReportControlsOptions {
	_options.ControlID = core.StringPtr(controlID)
	return _options
}

// SetControlName : Allow user to set ControlName
func (_options *GetReportControlsOptions) SetControlName(controlName string) *GetReportControlsOptions {
	_options.ControlName = core.StringPtr(controlName)
	return _options
}

// SetControlDescription : Allow user to set ControlDescription
func (_options *GetReportControlsOptions) SetControlDescription(controlDescription string) *GetReportControlsOptions {
	_options.ControlDescription = core.StringPtr(controlDescription)
	return _options
}

// SetControlCategory : Allow user to set ControlCategory
func (_options *GetReportControlsOptions) SetControlCategory(controlCategory string) *GetReportControlsOptions {
	_options.ControlCategory = core.StringPtr(controlCategory)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *GetReportControlsOptions) SetStatus(status string) *GetReportControlsOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *GetReportControlsOptions) SetSort(sort string) *GetReportControlsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetReportControlsOptions) SetScopeID(scopeID string) *GetReportControlsOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *GetReportControlsOptions) SetSubscopeID(subscopeID string) *GetReportControlsOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportControlsOptions) SetHeaders(param map[string]string) *GetReportControlsOptions {
	options.Headers = param
	return options
}

// GetReportDownloadFileOptions : The GetReportDownloadFile options.
type GetReportDownloadFileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The type of the response: application/csv or application/pdf.
	Accept *string `json:"Accept,omitempty"`

	// The indication of whether report summary metadata must be excluded.
	ExcludeSummary *bool `json:"exclude_summary,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportDownloadFileOptions : Instantiate GetReportDownloadFileOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportDownloadFileOptions(instanceID string, reportID string) *GetReportDownloadFileOptions {
	return &GetReportDownloadFileOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportDownloadFileOptions) SetInstanceID(instanceID string) *GetReportDownloadFileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportDownloadFileOptions) SetReportID(reportID string) *GetReportDownloadFileOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetReportDownloadFileOptions) SetAccept(accept string) *GetReportDownloadFileOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetExcludeSummary : Allow user to set ExcludeSummary
func (_options *GetReportDownloadFileOptions) SetExcludeSummary(excludeSummary bool) *GetReportDownloadFileOptions {
	_options.ExcludeSummary = core.BoolPtr(excludeSummary)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportDownloadFileOptions) SetHeaders(param map[string]string) *GetReportDownloadFileOptions {
	options.Headers = param
	return options
}

// GetReportOptions : The GetReport options.
type GetReportOptions struct {
	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportOptions : Instantiate GetReportOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportOptions(reportID string, instanceID string) *GetReportOptions {
	return &GetReportOptions{
		ReportID:   core.StringPtr(reportID),
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportOptions) SetReportID(reportID string) *GetReportOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportOptions) SetInstanceID(instanceID string) *GetReportOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetReportOptions) SetScopeID(scopeID string) *GetReportOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *GetReportOptions) SetSubscopeID(subscopeID string) *GetReportOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportOptions) SetHeaders(param map[string]string) *GetReportOptions {
	options.Headers = param
	return options
}

// GetReportRuleOptions : The GetReportRule options.
type GetReportRuleOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of a rule/assessment.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportRuleOptions : Instantiate GetReportRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportRuleOptions(instanceID string, reportID string, ruleID string) *GetReportRuleOptions {
	return &GetReportRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportRuleOptions) SetInstanceID(instanceID string) *GetReportRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportRuleOptions) SetReportID(reportID string) *GetReportRuleOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetReportRuleOptions) SetRuleID(ruleID string) *GetReportRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportRuleOptions) SetHeaders(param map[string]string) *GetReportRuleOptions {
	options.Headers = param
	return options
}

// GetReportSummaryOptions : The GetReportSummary options.
type GetReportSummaryOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportSummaryOptions : Instantiate GetReportSummaryOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportSummaryOptions(instanceID string, reportID string) *GetReportSummaryOptions {
	return &GetReportSummaryOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportSummaryOptions) SetInstanceID(instanceID string) *GetReportSummaryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportSummaryOptions) SetReportID(reportID string) *GetReportSummaryOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportSummaryOptions) SetHeaders(param map[string]string) *GetReportSummaryOptions {
	options.Headers = param
	return options
}

// GetReportTagsOptions : The GetReportTags options.
type GetReportTagsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportTagsOptions : Instantiate GetReportTagsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportTagsOptions(instanceID string, reportID string) *GetReportTagsOptions {
	return &GetReportTagsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportTagsOptions) SetInstanceID(instanceID string) *GetReportTagsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportTagsOptions) SetReportID(reportID string) *GetReportTagsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportTagsOptions) SetHeaders(param map[string]string) *GetReportTagsOptions {
	options.Headers = param
	return options
}

// GetReportViolationsDriftOptions : The GetReportViolationsDrift options.
type GetReportViolationsDriftOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The duration of the `scan_time` timestamp in number of days.
	ScanTimeDuration *int64 `json:"scan_time_duration,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetReportViolationsDriftOptions : Instantiate GetReportViolationsDriftOptions
func (*SecurityAndComplianceCenterApiV3) NewGetReportViolationsDriftOptions(instanceID string, reportID string) *GetReportViolationsDriftOptions {
	return &GetReportViolationsDriftOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetReportViolationsDriftOptions) SetInstanceID(instanceID string) *GetReportViolationsDriftOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetReportViolationsDriftOptions) SetReportID(reportID string) *GetReportViolationsDriftOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetScanTimeDuration : Allow user to set ScanTimeDuration
func (_options *GetReportViolationsDriftOptions) SetScanTimeDuration(scanTimeDuration int64) *GetReportViolationsDriftOptions {
	_options.ScanTimeDuration = core.Int64Ptr(scanTimeDuration)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetReportViolationsDriftOptions) SetScopeID(scopeID string) *GetReportViolationsDriftOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *GetReportViolationsDriftOptions) SetSubscopeID(subscopeID string) *GetReportViolationsDriftOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetReportViolationsDriftOptions) SetHeaders(param map[string]string) *GetReportViolationsDriftOptions {
	options.Headers = param
	return options
}

// GetRuleOptions : The GetRule options.
type GetRuleOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of a rule/assessment.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRuleOptions : Instantiate GetRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewGetRuleOptions(instanceID string, ruleID string) *GetRuleOptions {
	return &GetRuleOptions{
		InstanceID: core.StringPtr(instanceID),
		RuleID:     core.StringPtr(ruleID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetRuleOptions) SetInstanceID(instanceID string) *GetRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *GetRuleOptions) SetRuleID(ruleID string) *GetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRuleOptions) SetHeaders(param map[string]string) *GetRuleOptions {
	options.Headers = param
	return options
}

// GetScanReportDownloadFileOptions : The GetScanReportDownloadFile options.
type GetScanReportDownloadFileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the scan_result.
	JobID *string `json:"job_id" validate:"required,ne="`

	// The type of the response: application/csv or application/pdf.
	Accept *string `json:"Accept,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScanReportDownloadFileOptions : Instantiate GetScanReportDownloadFileOptions
func (*SecurityAndComplianceCenterApiV3) NewGetScanReportDownloadFileOptions(instanceID string, reportID string, jobID string) *GetScanReportDownloadFileOptions {
	return &GetScanReportDownloadFileOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
		JobID:      core.StringPtr(jobID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetScanReportDownloadFileOptions) SetInstanceID(instanceID string) *GetScanReportDownloadFileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetScanReportDownloadFileOptions) SetReportID(reportID string) *GetScanReportDownloadFileOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetJobID : Allow user to set JobID
func (_options *GetScanReportDownloadFileOptions) SetJobID(jobID string) *GetScanReportDownloadFileOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetAccept : Allow user to set Accept
func (_options *GetScanReportDownloadFileOptions) SetAccept(accept string) *GetScanReportDownloadFileOptions {
	_options.Accept = core.StringPtr(accept)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScanReportDownloadFileOptions) SetHeaders(param map[string]string) *GetScanReportDownloadFileOptions {
	options.Headers = param
	return options
}

// GetScanReportOptions : The GetScanReport options.
type GetScanReportOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the scan_result.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScanReportOptions : Instantiate GetScanReportOptions
func (*SecurityAndComplianceCenterApiV3) NewGetScanReportOptions(instanceID string, reportID string, jobID string) *GetScanReportOptions {
	return &GetScanReportOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
		JobID:      core.StringPtr(jobID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetScanReportOptions) SetInstanceID(instanceID string) *GetScanReportOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *GetScanReportOptions) SetReportID(reportID string) *GetScanReportOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetJobID : Allow user to set JobID
func (_options *GetScanReportOptions) SetJobID(jobID string) *GetScanReportOptions {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScanReportOptions) SetHeaders(param map[string]string) *GetScanReportOptions {
	options.Headers = param
	return options
}

// GetScopeOptions : The GetScope options.
type GetScopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetScopeOptions : Instantiate GetScopeOptions
func (*SecurityAndComplianceCenterApiV3) NewGetScopeOptions(instanceID string, scopeID string) *GetScopeOptions {
	return &GetScopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetScopeOptions) SetInstanceID(instanceID string) *GetScopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetScopeOptions) SetScopeID(scopeID string) *GetScopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetScopeOptions) SetHeaders(param map[string]string) *GetScopeOptions {
	options.Headers = param
	return options
}

// GetServiceOptions : The GetService options.
type GetServiceOptions struct {
	// The name of the corresponding service.
	ServicesName *string `json:"services_name" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetServiceOptions : Instantiate GetServiceOptions
func (*SecurityAndComplianceCenterApiV3) NewGetServiceOptions(servicesName string) *GetServiceOptions {
	return &GetServiceOptions{
		ServicesName: core.StringPtr(servicesName),
	}
}

// SetServicesName : Allow user to set ServicesName
func (_options *GetServiceOptions) SetServicesName(servicesName string) *GetServiceOptions {
	_options.ServicesName = core.StringPtr(servicesName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetServiceOptions) SetHeaders(param map[string]string) *GetServiceOptions {
	options.Headers = param
	return options
}

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*SecurityAndComplianceCenterApiV3) NewGetSettingsOptions(instanceID string) *GetSettingsOptions {
	return &GetSettingsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetSettingsOptions) SetInstanceID(instanceID string) *GetSettingsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// GetSubscopeOptions : The GetSubscope options.
type GetSubscopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	SubscopeID *string `json:"subscope_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSubscopeOptions : Instantiate GetSubscopeOptions
func (*SecurityAndComplianceCenterApiV3) NewGetSubscopeOptions(instanceID string, scopeID string, subscopeID string) *GetSubscopeOptions {
	return &GetSubscopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
		SubscopeID: core.StringPtr(subscopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetSubscopeOptions) SetInstanceID(instanceID string) *GetSubscopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *GetSubscopeOptions) SetScopeID(scopeID string) *GetSubscopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *GetSubscopeOptions) SetSubscopeID(subscopeID string) *GetSubscopeOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSubscopeOptions) SetHeaders(param map[string]string) *GetSubscopeOptions {
	options.Headers = param
	return options
}

// GetTargetOptions : The GetTarget options.
type GetTargetOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The target ID.
	TargetID *string `json:"target_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTargetOptions : Instantiate GetTargetOptions
func (*SecurityAndComplianceCenterApiV3) NewGetTargetOptions(instanceID string, targetID string) *GetTargetOptions {
	return &GetTargetOptions{
		InstanceID: core.StringPtr(instanceID),
		TargetID:   core.StringPtr(targetID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *GetTargetOptions) SetInstanceID(instanceID string) *GetTargetOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetTargetID : Allow user to set TargetID
func (_options *GetTargetOptions) SetTargetID(targetID string) *GetTargetOptions {
	_options.TargetID = core.StringPtr(targetID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTargetOptions) SetHeaders(param map[string]string) *GetTargetOptions {
	options.Headers = param
	return options
}

// Import : The collection of import parameters.
type Import struct {
	// The list of import parameters.
	Parameters []RuleParameter `json:"parameters,omitempty"`
}

// UnmarshalImport unmarshals an instance of Import from the specified map of raw messages.
func UnmarshalImport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Import)
	err = core.UnmarshalModel(m, "parameters", &obj.Parameters, UnmarshalRuleParameter)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LabelType : The label that is associated with the provider type.
type LabelType struct {
	// The text of the label.
	Text *string `json:"text,omitempty"`

	// The text to be shown when user hover overs the label.
	Tip *string `json:"tip,omitempty"`
}

// UnmarshalLabelType unmarshals an instance of LabelType from the specified map of raw messages.
func UnmarshalLabelType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LabelType)
	err = core.UnmarshalPrimitive(m, "text", &obj.Text)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tip", &obj.Tip)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LastScan : The last scan performed on a profile attachment.
type LastScan struct {
	// The ID of the last scan.
	ID *string `json:"id,omitempty"`

	// Details the state of the last scan.
	Status *string `json:"status,omitempty"`

	// The last time the scan ran.
	Time *strfmt.DateTime `json:"time,omitempty"`
}

// UnmarshalLastScan unmarshals an instance of LastScan from the specified map of raw messages.
func UnmarshalLastScan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LastScan)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "time", &obj.Time)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListControlLibrariesOptions : The ListControlLibraries options.
type ListControlLibrariesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The field that indicates how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// The field that indicate how you want the resources to be filtered by.
	ControlLibraryType *string `json:"control_library_type,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListControlLibrariesOptions : Instantiate ListControlLibrariesOptions
func (*SecurityAndComplianceCenterApiV3) NewListControlLibrariesOptions(instanceID string) *ListControlLibrariesOptions {
	return &ListControlLibrariesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListControlLibrariesOptions) SetInstanceID(instanceID string) *ListControlLibrariesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListControlLibrariesOptions) SetAccountID(accountID string) *ListControlLibrariesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetControlLibraryType : Allow user to set the ControlLibraryType
func (_options *ListControlLibrariesOptions) SetControlLibraryType(controlLibraryType string) *ListControlLibrariesOptions {
	_options.ControlLibraryType = &controlLibraryType
	return _options
}

// SetStart : Allow user to set the Start
func (_options *ListControlLibrariesOptions) SetStart(controlLibraryType string) *ListControlLibrariesOptions {
	_options.Start = &controlLibraryType
	return _options
}

// SetLimit : Allow user to set the Limit
func (_options *ListControlLibrariesOptions) SetLimit(limit int64) *ListControlLibrariesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListControlLibrariesOptions) SetHeaders(param map[string]string) *ListControlLibrariesOptions {
	options.Headers = param
	return options
}

// ListInstanceAttachmentsOptions : The ListInstanceAttachments options.
type ListInstanceAttachmentsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The profile version group label.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The number of items that are retrieved in a collection.
	Limit *int64 `json:"limit,omitempty"`

	// The sorted collection of attachments. The available values are `created_on` and `scope_type`.
	Sort *string `json:"sort,omitempty"`

	// The collection of attachments that is sorted in ascending order. To sort the collection in descending order, use the
	// `DESC` schema.
	Direction *string `json:"direction,omitempty"`

	// The reference to the first item in the results page. Take the value from the `next` field that is in the response
	// from the previous page.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListInstanceAttachmentsOptions.Sort property.
// The sorted collection of attachments. The available values are `created_on` and `scope_type`.
const (
	ListInstanceAttachmentsOptionsSortCreatedOnConst = "created_on"
	ListInstanceAttachmentsOptionsSortScopeTypeConst = "scope_type"
)

// Constants associated with the ListInstanceAttachmentsOptions.Direction property.
// The collection of attachments that is sorted in ascending order. To sort the collection in descending order, use the
// `DESC` schema.
const (
	ListInstanceAttachmentsOptionsDirectionAscConst  = "asc"
	ListInstanceAttachmentsOptionsDirectionDescConst = "desc"
)

// NewListInstanceAttachmentsOptions : Instantiate ListInstanceAttachmentsOptions
func (*SecurityAndComplianceCenterApiV3) NewListInstanceAttachmentsOptions(instanceID string) *ListInstanceAttachmentsOptions {
	return &ListInstanceAttachmentsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListInstanceAttachmentsOptions) SetInstanceID(instanceID string) *ListInstanceAttachmentsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListInstanceAttachmentsOptions) SetAccountID(accountID string) *ListInstanceAttachmentsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetVersionGroupLabel : Allow user to set VersionGroupLabel
func (_options *ListInstanceAttachmentsOptions) SetVersionGroupLabel(versionGroupLabel string) *ListInstanceAttachmentsOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListInstanceAttachmentsOptions) SetLimit(limit int64) *ListInstanceAttachmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListInstanceAttachmentsOptions) SetSort(sort string) *ListInstanceAttachmentsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetDirection : Allow user to set Direction
func (_options *ListInstanceAttachmentsOptions) SetDirection(direction string) *ListInstanceAttachmentsOptions {
	_options.Direction = core.StringPtr(direction)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListInstanceAttachmentsOptions) SetStart(start string) *ListInstanceAttachmentsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListInstanceAttachmentsOptions) SetHeaders(param map[string]string) *ListInstanceAttachmentsOptions {
	options.Headers = param
	return options
}

// ListProfileAttachmentsOptions : The ListProfileAttachments options.
type ListProfileAttachmentsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProfileAttachmentsOptions : Instantiate ListProfileAttachmentsOptions
func (*SecurityAndComplianceCenterApiV3) NewListProfileAttachmentsOptions(instanceID string, profileID string) *ListProfileAttachmentsOptions {
	return &ListProfileAttachmentsOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProfileAttachmentsOptions) SetInstanceID(instanceID string) *ListProfileAttachmentsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListProfileAttachmentsOptions) SetProfileID(profileID string) *ListProfileAttachmentsOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListProfileAttachmentsOptions) SetAccountID(accountID string) *ListProfileAttachmentsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfileAttachmentsOptions) SetHeaders(param map[string]string) *ListProfileAttachmentsOptions {
	options.Headers = param
	return options
}

// ListProfileParametersOptions : The ListProfileParameters options.
type ListProfileParametersOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProfileParametersOptions : Instantiate ListProfileParametersOptions
func (*SecurityAndComplianceCenterApiV3) NewListProfileParametersOptions(instanceID string, profileID string) *ListProfileParametersOptions {
	return &ListProfileParametersOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProfileParametersOptions) SetInstanceID(instanceID string) *ListProfileParametersOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListProfileParametersOptions) SetProfileID(profileID string) *ListProfileParametersOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfileParametersOptions) SetHeaders(param map[string]string) *ListProfileParametersOptions {
	options.Headers = param
	return options
}

// ListProfilesOptions : The ListProfiles options.
type ListProfilesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string

	// The field that indicates how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// The field that indicate how you want the resources to be filtered by.
	ProfileType *string `json:"profile_type,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`
}

// NewListProfilesOptions : Instantiate ListProfilesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProfilesOptions(instanceID string) *ListProfilesOptions {
	return &ListProfilesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProfilesOptions) SetInstanceID(instanceID string) *ListProfilesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListProfilesOptions) SetAccountID(accountID string) *ListProfilesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *ListProfilesOptions) SetProfileType(typeVar string) *ListProfilesOptions {
	_options.ProfileType = core.StringPtr(typeVar)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListProfilesOptions) SetStart(start string) *ListProfilesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListProfilesOptions) SetLimit(limit int64) *ListProfilesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProfilesOptions) SetHeaders(param map[string]string) *ListProfilesOptions {
	options.Headers = param
	return options
}

// ListProviderTypeInstancesOptions : The ListProviderTypeInstances options.
type ListProviderTypeInstancesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProviderTypeInstancesOptions : Instantiate ListProviderTypeInstancesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProviderTypeInstancesOptions(instanceID string, providerTypeID string) *ListProviderTypeInstancesOptions {
	return &ListProviderTypeInstancesOptions{
		InstanceID:     core.StringPtr(instanceID),
		ProviderTypeID: core.StringPtr(providerTypeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProviderTypeInstancesOptions) SetInstanceID(instanceID string) *ListProviderTypeInstancesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *ListProviderTypeInstancesOptions) SetProviderTypeID(providerTypeID string) *ListProviderTypeInstancesOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProviderTypeInstancesOptions) SetHeaders(param map[string]string) *ListProviderTypeInstancesOptions {
	options.Headers = param
	return options
}

// ListProviderTypesOptions : The ListProviderTypes options.
type ListProviderTypesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListProviderTypesOptions : Instantiate ListProviderTypesOptions
func (*SecurityAndComplianceCenterApiV3) NewListProviderTypesOptions(instanceID string) *ListProviderTypesOptions {
	return &ListProviderTypesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListProviderTypesOptions) SetInstanceID(instanceID string) *ListProviderTypesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProviderTypesOptions) SetHeaders(param map[string]string) *ListProviderTypesOptions {
	options.Headers = param
	return options
}

// ListReportEvaluationsOptions : The ListReportEvaluations options.
type ListReportEvaluationsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the assessment.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The assessment method.
	AssessmentMethod *string `json:"assessment_method,omitempty"`

	// The ID of component.
	ComponentID *string `json:"component_id,omitempty"`

	// The ID of the evaluation target.
	TargetID *string `json:"target_id,omitempty"`

	// The environment of the evaluation target.
	TargetEnv *string `json:"target_env,omitempty"`

	// The name of the evaluation target.
	TargetName *string `json:"target_name,omitempty"`

	// The evaluation status value.
	Status *string `json:"status,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// This field sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportEvaluationsOptions.Status property.
// The evaluation status value.
const (
	ListReportEvaluationsOptionsStatusErrorConst   = "error"
	ListReportEvaluationsOptionsStatusFailureConst = "failure"
	ListReportEvaluationsOptionsStatusPassConst    = "pass"
	ListReportEvaluationsOptionsStatusSkippedConst = "skipped"
)

// NewListReportEvaluationsOptions : Instantiate ListReportEvaluationsOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportEvaluationsOptions(instanceID string, reportID string) *ListReportEvaluationsOptions {
	return &ListReportEvaluationsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportEvaluationsOptions) SetInstanceID(instanceID string) *ListReportEvaluationsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *ListReportEvaluationsOptions) SetReportID(reportID string) *ListReportEvaluationsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetAssessmentID : Allow user to set AssessmentID
func (_options *ListReportEvaluationsOptions) SetAssessmentID(assessmentID string) *ListReportEvaluationsOptions {
	_options.AssessmentID = core.StringPtr(assessmentID)
	return _options
}

// SetAssessmentMethod : Allow user to set AssessmentMethod
func (_options *ListReportEvaluationsOptions) SetAssessmentMethod(assessmentMethod string) *ListReportEvaluationsOptions {
	_options.AssessmentMethod = core.StringPtr(assessmentMethod)
	return _options
}

// SetComponentID : Allow user to set ComponentID
func (_options *ListReportEvaluationsOptions) SetComponentID(componentID string) *ListReportEvaluationsOptions {
	_options.ComponentID = core.StringPtr(componentID)
	return _options
}

// SetTargetID : Allow user to set TargetID
func (_options *ListReportEvaluationsOptions) SetTargetID(targetID string) *ListReportEvaluationsOptions {
	_options.TargetID = core.StringPtr(targetID)
	return _options
}

// SetTargetEnv : Allow user to set TargetEnv
func (_options *ListReportEvaluationsOptions) SetTargetEnv(targetEnv string) *ListReportEvaluationsOptions {
	_options.TargetEnv = core.StringPtr(targetEnv)
	return _options
}

// SetTargetName : Allow user to set TargetName
func (_options *ListReportEvaluationsOptions) SetTargetName(targetName string) *ListReportEvaluationsOptions {
	_options.TargetName = core.StringPtr(targetName)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ListReportEvaluationsOptions) SetStatus(status string) *ListReportEvaluationsOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportEvaluationsOptions) SetStart(start string) *ListReportEvaluationsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportEvaluationsOptions) SetLimit(limit int64) *ListReportEvaluationsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListReportEvaluationsOptions) SetSort(sort string) *ListReportEvaluationsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *ListReportEvaluationsOptions) SetScopeID(scopeID string) *ListReportEvaluationsOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *ListReportEvaluationsOptions) SetSubscopeID(subscopeID string) *ListReportEvaluationsOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportEvaluationsOptions) SetHeaders(param map[string]string) *ListReportEvaluationsOptions {
	options.Headers = param
	return options
}

// ListReportResourcesOptions : The ListReportResources options.
type ListReportResourcesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The ID of component.
	ComponentID *string `json:"component_id,omitempty"`

	// The compliance status value.
	Status *string `json:"status,omitempty"`

	// This field sorts resources by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportResourcesOptions.Status property.
// The compliance status value.
const (
	ListReportResourcesOptionsStatusCompliantConst              = "compliant"
	ListReportResourcesOptionsStatusNotApplicableConst          = "not_applicable"
	ListReportResourcesOptionsStatusNotCompliantConst           = "not_compliant"
	ListReportResourcesOptionsStatusUnableToPerformConst        = "unable_to_perform"
	ListReportResourcesOptionsStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// Constants associated with the ListReportResourcesOptions.Sort property.
// This field sorts resources by using a valid sort field. To learn more, see
// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
const (
	ListReportResourcesOptionsSortAccountIDConst    = "account_id"
	ListReportResourcesOptionsSortComponentIDConst  = "component_id"
	ListReportResourcesOptionsSortResourceNameConst = "resource_name"
	ListReportResourcesOptionsSortStatusConst       = "status"
)

// NewListReportResourcesOptions : Instantiate ListReportResourcesOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportResourcesOptions(instanceID string, reportID string) *ListReportResourcesOptions {
	return &ListReportResourcesOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportResourcesOptions) SetInstanceID(instanceID string) *ListReportResourcesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *ListReportResourcesOptions) SetReportID(reportID string) *ListReportResourcesOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ListReportResourcesOptions) SetID(id string) *ListReportResourcesOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetResourceName : Allow user to set ResourceName
func (_options *ListReportResourcesOptions) SetResourceName(resourceName string) *ListReportResourcesOptions {
	_options.ResourceName = core.StringPtr(resourceName)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ListReportResourcesOptions) SetAccountID(accountID string) *ListReportResourcesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetComponentID : Allow user to set ComponentID
func (_options *ListReportResourcesOptions) SetComponentID(componentID string) *ListReportResourcesOptions {
	_options.ComponentID = core.StringPtr(componentID)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ListReportResourcesOptions) SetStatus(status string) *ListReportResourcesOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListReportResourcesOptions) SetSort(sort string) *ListReportResourcesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportResourcesOptions) SetStart(start string) *ListReportResourcesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportResourcesOptions) SetLimit(limit int64) *ListReportResourcesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *ListReportResourcesOptions) SetScopeID(scopeID string) *ListReportResourcesOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *ListReportResourcesOptions) SetSubscopeID(subscopeID string) *ListReportResourcesOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportResourcesOptions) SetHeaders(param map[string]string) *ListReportResourcesOptions {
	options.Headers = param
	return options
}

// ListReportsOptions : The ListReports options.
type ListReportsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the attachment.
	AttachmentID *string `json:"attachment_id,omitempty"`

	// The report group ID.
	GroupID *string `json:"group_id,omitempty"`

	// The ID of the profile.
	ProfileID *string `json:"profile_id,omitempty"`

	// The type of the scan.
	Type *string `json:"type,omitempty"`

	// The indication of what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The indication of many resources to return, unless the response is  the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// This field sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListReportsOptions.Type property.
// The type of the scan.
const (
	ListReportsOptionsTypeOndemandConst  = "ondemand"
	ListReportsOptionsTypeScheduledConst = "scheduled"
)

// NewListReportsOptions : Instantiate ListReportsOptions
func (*SecurityAndComplianceCenterApiV3) NewListReportsOptions(instanceID string) *ListReportsOptions {
	return &ListReportsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListReportsOptions) SetInstanceID(instanceID string) *ListReportsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *ListReportsOptions) SetAttachmentID(attachmentID string) *ListReportsOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetGroupID : Allow user to set GroupID
func (_options *ListReportsOptions) SetGroupID(groupID string) *ListReportsOptions {
	_options.GroupID = core.StringPtr(groupID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ListReportsOptions) SetProfileID(profileID string) *ListReportsOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListReportsOptions) SetType(typeVar string) *ListReportsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListReportsOptions) SetStart(start string) *ListReportsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListReportsOptions) SetLimit(limit int64) *ListReportsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListReportsOptions) SetSort(sort string) *ListReportsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListReportsOptions) SetHeaders(param map[string]string) *ListReportsOptions {
	options.Headers = param
	return options
}

// ListRulesOptions : The ListRules options.
type ListRulesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// The list of only user-defined, or system-defined rules.
	Type *string `json:"type,omitempty"`

	// The indication of whether to search for rules with a specific string string in the name, description, or labels.
	Search *string `json:"search,omitempty"`

	// Searches for rules targeting corresponding service.
	ServiceName *string `json:"service_name,omitempty"`

	// Field used to sort rules. Rules can be sorted in ascending or descending order.
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListRulesOptions.Type property.
// The list of only user-defined, or system-defined rules.
const (
	ListRulesOptionsTypeSystemDefinedConst = "system_defined"
	ListRulesOptionsTypeUserDefinedConst   = "user_defined"
)

// Constants associated with the ListRulesOptions.Sort property.
// Field used to sort rules. Rules can be sorted in ascending or descending order.
const (
	ListRulesOptionsSortDescriptionConst        = "description"
	ListRulesOptionsSortServiceDisplayNameConst = "service_display_name"
	ListRulesOptionsSortTypeConst               = "type"
	ListRulesOptionsSortUpdatedOnConst          = "updated_on"
)

// NewListRulesOptions : Instantiate ListRulesOptions
func (*SecurityAndComplianceCenterApiV3) NewListRulesOptions(instanceID string) *ListRulesOptions {
	return &ListRulesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListRulesOptions) SetInstanceID(instanceID string) *ListRulesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListRulesOptions) SetLimit(limit int64) *ListRulesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListRulesOptions) SetStart(start string) *ListRulesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListRulesOptions) SetType(typeVar string) *ListRulesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListRulesOptions) SetSearch(search string) *ListRulesOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListRulesOptions) SetServiceName(serviceName string) *ListRulesOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListRulesOptions) SetSort(sort string) *ListRulesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRulesOptions) SetHeaders(param map[string]string) *ListRulesOptions {
	options.Headers = param
	return options
}

// ListScanReportsOptions : The ListScanReports options.
type ListScanReportsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scan that is associated with a report.
	ReportID *string `json:"report_id" validate:"required,ne="`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListScanReportsOptions : Instantiate ListScanReportsOptions
func (*SecurityAndComplianceCenterApiV3) NewListScanReportsOptions(instanceID string, reportID string) *ListScanReportsOptions {
	return &ListScanReportsOptions{
		InstanceID: core.StringPtr(instanceID),
		ReportID:   core.StringPtr(reportID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListScanReportsOptions) SetInstanceID(instanceID string) *ListScanReportsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetReportID : Allow user to set ReportID
func (_options *ListScanReportsOptions) SetReportID(reportID string) *ListScanReportsOptions {
	_options.ReportID = core.StringPtr(reportID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *ListScanReportsOptions) SetScopeID(scopeID string) *ListScanReportsOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *ListScanReportsOptions) SetSubscopeID(subscopeID string) *ListScanReportsOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListScanReportsOptions) SetHeaders(param map[string]string) *ListScanReportsOptions {
	options.Headers = param
	return options
}

// ListScopesOptions : The ListScopes options.
type ListScopesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// determine name of scope returned in response.
	Name *string `json:"name,omitempty"`

	// determine descriptions of scope returned in response.
	Description *string `json:"description,omitempty"`

	// determine environment of scopes returned in response.
	Environment *string `json:"environment,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListScopesOptions : Instantiate ListScopesOptions
func (*SecurityAndComplianceCenterApiV3) NewListScopesOptions(instanceID string) *ListScopesOptions {
	return &ListScopesOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListScopesOptions) SetInstanceID(instanceID string) *ListScopesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListScopesOptions) SetLimit(limit int64) *ListScopesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListScopesOptions) SetStart(start string) *ListScopesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListScopesOptions) SetName(name string) *ListScopesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ListScopesOptions) SetDescription(description string) *ListScopesOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnvironment : Allow user to set Environment
func (_options *ListScopesOptions) SetEnvironment(environment string) *ListScopesOptions {
	_options.Environment = core.StringPtr(environment)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListScopesOptions) SetHeaders(param map[string]string) *ListScopesOptions {
	options.Headers = param
	return options
}

// ListServicesOptions : The ListServices options.
type ListServicesOptions struct {
	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListServicesOptions : Instantiate ListServicesOptions
func (*SecurityAndComplianceCenterApiV3) NewListServicesOptions() *ListServicesOptions {
	return &ListServicesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListServicesOptions) SetHeaders(param map[string]string) *ListServicesOptions {
	options.Headers = param
	return options
}

// ListSubscopesOptions : The ListSubscopes options.
type ListSubscopesOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The indication of how many resources to return, unless the response is the last page of resources.
	Limit *int64 `json:"limit,omitempty"`

	// Determine what resource to start the page on or after.
	Start *string `json:"start,omitempty"`

	// determine name of subscope returned in response.
	Name *string `json:"name,omitempty"`

	// determine descriptions of subscopes returned in response.
	Description *string `json:"description,omitempty"`

	// determine environment of subscopes returned in response.
	Environment *string `json:"environment,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSubscopesOptions : Instantiate ListSubscopesOptions
func (*SecurityAndComplianceCenterApiV3) NewListSubscopesOptions(instanceID string, scopeID string) *ListSubscopesOptions {
	return &ListSubscopesOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListSubscopesOptions) SetInstanceID(instanceID string) *ListSubscopesOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *ListSubscopesOptions) SetScopeID(scopeID string) *ListSubscopesOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSubscopesOptions) SetLimit(limit int64) *ListSubscopesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListSubscopesOptions) SetStart(start string) *ListSubscopesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListSubscopesOptions) SetName(name string) *ListSubscopesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ListSubscopesOptions) SetDescription(description string) *ListSubscopesOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnvironment : Allow user to set Environment
func (_options *ListSubscopesOptions) SetEnvironment(environment string) *ListSubscopesOptions {
	_options.Environment = core.StringPtr(environment)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSubscopesOptions) SetHeaders(param map[string]string) *ListSubscopesOptions {
	options.Headers = param
	return options
}

// ListTargetsOptions : The ListTargets options.
type ListTargetsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListTargetsOptions : Instantiate ListTargetsOptions
func (*SecurityAndComplianceCenterApiV3) NewListTargetsOptions(instanceID string) *ListTargetsOptions {
	return &ListTargetsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ListTargetsOptions) SetInstanceID(instanceID string) *ListTargetsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTargetsOptions) SetHeaders(param map[string]string) *ListTargetsOptions {
	options.Headers = param
	return options
}

// MultiCloudScopePayload : MultiCloudScopePayload struct
type MultiCloudScopePayload struct {
	// The ID of the scope.
	ID *string `json:"id,omitempty"`

	// The environment that relates to this scope.
	Environment *string `json:"environment,omitempty"`

	// The properties supported for scoping by this environment.
	Properties []ScopePropertyIntf `json:"properties,omitempty"`
}

// UnmarshalMultiCloudScopePayload unmarshals an instance of MultiCloudScopePayload from the specified map of raw messages.
func UnmarshalMultiCloudScopePayload(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MultiCloudScopePayload)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalScopeProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectStorage : The Cloud Object Storage settings.
type ObjectStorage struct {
	// The connected Cloud Object Storage instance CRN.
	InstanceCRN *string `json:"instance_crn,omitempty"`

	// The connected Cloud Object Storage bucket name.
	Bucket *string `json:"bucket,omitempty"`

	// The connected Cloud Object Storage bucket location.
	BucketLocation *string `json:"bucket_location,omitempty"`

	// The connected Cloud Object Storage bucket endpoint.
	BucketEndpoint *string `json:"bucket_endpoint,omitempty"`

	// The date when the bucket connection was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`
}

// UnmarshalObjectStorage unmarshals an instance of ObjectStorage from the specified map of raw messages.
func UnmarshalObjectStorage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectStorage)
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket_location", &obj.BucketLocation)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket_endpoint", &obj.BucketEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ObjectStoragePrototype : The payload to connect a Cloud Object Storage instance to an Security and Compliance Center instance.
type ObjectStoragePrototype struct {
	// The connected Cloud Object Storage bucket name.
	Bucket *string `json:"bucket,omitempty"`

	// The connected Cloud Object Storage instance CRN.
	InstanceCRN *string `json:"instance_crn,omitempty"`
}

// UnmarshalObjectStoragePrototype unmarshals an instance of ObjectStoragePrototype from the specified map of raw messages.
func UnmarshalObjectStoragePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ObjectStoragePrototype)
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCRN)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageHRefFirst : A page reference.
type PageHRefFirst struct {
	// The URL for the first page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPageHRefFirst unmarshals an instance of PageHRefFirst from the specified map of raw messages.
func UnmarshalPageHRefFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHRefFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageHRefNext : A page reference.
type PageHRefNext struct {
	// The URL for the next page.
	Href *string `json:"href" validate:"required"`

	// The token of the next page, when it's present.
	Start *string `json:"start,omitempty"`
}

// UnmarshalPageHRefNext unmarshals an instance of PageHRefNext from the specified map of raw messages.
func UnmarshalPageHRefNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHRefNext)
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

// Parameter : The details of a parameter used during an assessment.
type Parameter struct {
	// The type of evaluation.
	AssessmentType *string `json:"assessment_type,omitempty"`

	// The ID of the assessment.
	AssessmentID *string `json:"assessment_id,omitempty"`

	// The parameter name.
	ParameterName *string `json:"parameter_name,omitempty"`

	// The parameter display name.
	ParameterDisplayName *string `json:"parameter_display_name,omitempty"`

	// The parameter type.
	ParameterType *string `json:"parameter_type,omitempty"`

	ParameterValue interface{} `json:"parameter_value,omitempty"`
}

// UnmarshalParameter unmarshals an instance of Parameter from the specified map of raw messages.
func UnmarshalParameter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Parameter)
	err = core.UnmarshalPrimitive(m, "assessment_type", &obj.AssessmentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_name", &obj.ParameterName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_display_name", &obj.ParameterDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_type", &obj.ParameterType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameter_value", &obj.ParameterValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostTestEventOptions : The PostTestEvent options.
type PostTestEventOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPostTestEventOptions : Instantiate PostTestEventOptions
func (*SecurityAndComplianceCenterApiV3) NewPostTestEventOptions(instanceID string) *PostTestEventOptions {
	return &PostTestEventOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *PostTestEventOptions) SetInstanceID(instanceID string) *PostTestEventOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PostTestEventOptions) SetHeaders(param map[string]string) *PostTestEventOptions {
	options.Headers = param
	return options
}

// Profile : A group of controls that are related to a compliance objective.
type Profile struct {
	// The ID of the profile.
	ID *string `json:"id,omitempty"`

	// The name of the profile.
	ProfileName *string `json:"profile_name,omitempty"`

	// A description of what the profile should represent.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The type of profile, either predefined or custom.
	ProfileType *string `json:"profile_type,omitempty"`

	// The version of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// The unique identifier of the revision.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// Determines if the profile is up to date with the latest revisions.
	Latest *bool `json:"latest,omitempty"`

	// User who created the profile.
	CreatedBy *string `json:"created_by,omitempty"`

	// The ID associated with the profile.
	InstanceID *string `json:"instance_id,omitempty"`

	// Determines if an heirarchy is enabled.
	HierarchyEnabled *bool `json:"heirarchy_enabled,omitempty"`

	// The date when the profile was created, in date-time format.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// User who made the latest changes to the profile.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date when the profile was last updated, in date-time format.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The number of controls contained in the profile.
	ControlsCount *int64 `json:"controls_count,omitempty"`

	// The number of parent controls contained in the profile.
	ControlParentsCount *int64 `json:"control_parents_count,omitempty"`

	// The number of attachments associated with the profile.
	AttachmentsCount *int64 `json:"attachments_count,omitempty"`

	// The list of controls.
	Controls []ProfileControlsInResponse `json:"controls,omitempty"`

	// The default parameters of the profile.
	DefaultParameters []DefaultParameters `json:"default_parameters,omitempty"`
}

// Constants associated with the Profile.ProfileType property.
// The type of profile, either predefined or custom.
const (
	ProfileProfileTypeCustomConst     = "custom"
	ProfileProfileTypePredefinedConst = "predefined"
)

// UnmarshalProfile unmarshals an instance of Profile from the specified map of raw messages.
func UnmarshalProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Profile)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_name", &obj.ProfileName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_description", &obj.ProfileDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_type", &obj.ProfileType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_version", &obj.ProfileVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version_group_label", &obj.VersionGroupLabel)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latest", &obj.Latest)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "controls_count", &obj.ControlsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachments_count", &obj.AttachmentsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalProfileControlsInResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "default_parameters", &obj.DefaultParameters, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileAssessmentPrototype : The attributes needed for the payload.
type ProfileAssessmentPrototype struct {
	// The rule or assessment to target.
	AssessmentID *string `json:"assessment_id,omitempty"`
}

// UnmarshalProfileAssessmentPrototype unmarshals an instance of ProfileAssessmentPrototype from the specified map of raw messages.
func UnmarshalProfileAssessmentPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileAssessmentPrototype)
	err = core.UnmarshalPrimitive(m, "assessment_id", &obj.AssessmentID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileAttachment : The configuration set when starting a scan against a profile.
type ProfileAttachment struct {
	// The parameters associated with the profile attachment.
	AttachmentParameters []Parameter `json:"attachment_parameters" validate:"required"`

	// The details to describe the profile attachment.
	Description *string `json:"description" validate:"required"`

	// The name of the Profile Attachment.
	Name *string `json:"name" validate:"required"`

	// The notification configuration of the attachment.
	Notifications *AttachmentNotifications `json:"notifications" validate:"required"`

	// Details how often a scan from a profile attachment is ran.
	Schedule *string `json:"schedule" validate:"required"`

	// A list of scopes associated with a profile attachment.
	Scope []MultiCloudScopePayload `json:"scope" validate:"required"`

	// Details the state of a profile attachment.
	Status *string `json:"status" validate:"required"`

	// The ID of the account.
	AccountID *string `json:"account_id,omitempty"`

	// User who created the profile attachment.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date-time that the profile attachment was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The ID of the profile attachment.
	ID *string `json:"id,omitempty"`

	// The ID of the associated Security and Compliance Center instance.
	InstanceID *string `json:"instance_id,omitempty"`

	// The last scan performed on a profile attachment.
	LastScan *LastScan `json:"last_scan,omitempty"`

	// The date-time for next scan.
	NextScanTime *strfmt.DateTime `json:"next_scan_time,omitempty"`

	// The ID of the profile.
	ProfileID *string `json:"profile_id,omitempty"`

	// User who made the latest changes to the profile attachment.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The date when the profile was last updated, in date-time format.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`
}

// Constants associated with the ProfileAttachment.Schedule property.
// Details how often a scan from a profile attachment is ran.
const (
	ProfileAttachmentScheduleDailyConst       = "daily"
	ProfileAttachmentScheduleEvery30DaysConst = "every_30_days"
	ProfileAttachmentScheduleEvery7DaysConst  = "every_7_days"
)

// Constants associated with the ProfileAttachment.Status property.
// Details the state of a profile attachment.
const (
	ProfileAttachmentStatusDisabledConst = "disabled"
	ProfileAttachmentStatusEnabledConst  = "enabled"
)

// UnmarshalProfileAttachment unmarshals an instance of ProfileAttachment from the specified map of raw messages.
func UnmarshalProfileAttachment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileAttachment)
	err = core.UnmarshalModel(m, "attachment_parameters", &obj.AttachmentParameters, UnmarshalParameter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalAttachmentNotifications)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalMultiCloudScopePayload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last_scan", &obj.LastScan, UnmarshalLastScan)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "next_scan_time", &obj.NextScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileAttachmentBase : The prototype for creating a Profile attachment.
type ProfileAttachmentBase struct {
	// The parameters associated with the profile attachment.
	AttachmentParameters []Parameter `json:"attachment_parameters" validate:"required"`

	// The details to describe the profile attachment.
	Description *string `json:"description" validate:"required"`

	// The name of the Profile Attachment.
	Name *string `json:"name" validate:"required"`

	// The notification configuration of the attachment.
	Notifications *AttachmentNotifications `json:"notifications" validate:"required"`

	// Details how often a scan from a profile attachment is ran.
	Schedule *string `json:"schedule" validate:"required"`

	// A list of scopes associated with a profile attachment.
	Scope []MultiCloudScopePayload `json:"scope" validate:"required"`

	// Details the state of a profile attachment.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the ProfileAttachmentBase.Schedule property.
// Details how often a scan from a profile attachment is ran.
const (
	ProfileAttachmentBaseScheduleDailyConst       = "daily"
	ProfileAttachmentBaseScheduleEvery30DaysConst = "every_30_days"
	ProfileAttachmentBaseScheduleEvery7DaysConst  = "every_7_days"
)

// Constants associated with the ProfileAttachmentBase.Status property.
// Details the state of a profile attachment.
const (
	ProfileAttachmentBaseStatusDisabledConst = "disabled"
	ProfileAttachmentBaseStatusEnabledConst  = "enabled"
)

// NewProfileAttachmentBase : Instantiate ProfileAttachmentBase (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewProfileAttachmentBase(attachmentParameters []Parameter, description string, name string, notifications *AttachmentNotifications, schedule string, scope []MultiCloudScopePayload, status string) (_model *ProfileAttachmentBase, err error) {
	_model = &ProfileAttachmentBase{
		AttachmentParameters: attachmentParameters,
		Description:          core.StringPtr(description),
		Name:                 core.StringPtr(name),
		Notifications:        notifications,
		Schedule:             core.StringPtr(schedule),
		Scope:                scope,
		Status:               core.StringPtr(status),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProfileAttachmentBase unmarshals an instance of ProfileAttachmentBase from the specified map of raw messages.
func UnmarshalProfileAttachmentBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileAttachmentBase)
	err = core.UnmarshalModel(m, "attachment_parameters", &obj.AttachmentParameters, UnmarshalParameter)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "notifications", &obj.Notifications, UnmarshalAttachmentNotifications)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schedule", &obj.Schedule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalMultiCloudScopePayload)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileAttachmentCollection : A list of ProfileAttachment tied to a profile or instance.
type ProfileAttachmentCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// List of attachments.
	Attachments []ProfileAttachment `json:"attachments,omitempty"`
}

// UnmarshalProfileAttachmentCollection unmarshals an instance of ProfileAttachmentCollection from the specified map of raw messages.
func UnmarshalProfileAttachmentCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileAttachmentCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalProfileAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProfileAttachmentCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ProfileAttachmentResponse : The response coming back from creating a ProfileAttachment.
type ProfileAttachmentResponse struct {
	// The ID of the profile.
	ProfileID *string `json:"profile_id,omitempty"`

	// List of profile attachments associated with profile.
	Attachments []ProfileAttachment `json:"attachments,omitempty"`
}

// UnmarshalProfileAttachmentResponse unmarshals an instance of ProfileAttachmentResponse from the specified map of raw messages.
func UnmarshalProfileAttachmentResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileAttachmentResponse)
	err = core.UnmarshalPrimitive(m, "profile_id", &obj.ProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachments", &obj.Attachments, UnmarshalProfileAttachment)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileCollection : Show a list of Profiles.
type ProfileCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// A list of profiles.
	Profiles []Profile `json:"profiles,omitempty"`
}

// UnmarshalProfileCollection unmarshals an instance of ProfileCollection from the specified map of raw messages.
func UnmarshalProfileCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profiles", &obj.Profiles, UnmarshalProfile)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileControlSpecificationPrototype : The control specification payload.
type ProfileControlSpecificationPrototype struct {
	// The ID of the control specification to target.
	ControlSpecificationID *string `json:"control_specification_id,omitempty"`

	// List of rules to target.
	Assessments []ProfileAssessmentPrototype `json:"assessments,omitempty"`
}

// UnmarshalProfileControlSpecificationPrototype unmarshals an instance of ProfileControlSpecificationPrototype from the specified map of raw messages.
func UnmarshalProfileControlSpecificationPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileControlSpecificationPrototype)
	err = core.UnmarshalPrimitive(m, "control_specification_id", &obj.ControlSpecificationID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "assessments", &obj.Assessments, UnmarshalProfileAssessmentPrototype)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileControlsInResponse : The control details for a profile.
type ProfileControlsInResponse struct {
	// The ID of the control library that contains a profile.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The control ID.
	ControlID *string `json:"control_id,omitempty"`

	// The control library version.
	ControlLibraryVersion *string `json:"control_library_version,omitempty"`

	// The control name.
	ControlName *string `json:"control_name,omitempty"`

	// Determines if the control needs to be satisfied
	ControlRequirement *bool `json:"control_requirement,omitempty"`

	// The control description.
	ControlDescription *string `json:"control_description,omitempty"`

	// The control severity.
	ControlSeverity *string `json:"control_severity,omitempty"`

	// The control category.
	ControlCategory *string `json:"control_category,omitempty"`

	// The control parent.
	ControlParent *string `json:"control_parent,omitempty"`

	// References to a control documentation.
	ControlDocs *ControlDoc `json:"control_docs,omitempty"`

	// The number of control specifications in the control
	ControlSpecificationsCount *int64 `json:"control_specifications_count,omitempty"`

	// List of control specifications in a profile.
	ControlSpecifications []ControlSpecification `json:"control_specifications,omitempty"`
}

// UnmarshalProfileControlsInResponse unmarshals an instance of ProfileControlsInResponse from the specified map of raw messages.
func UnmarshalProfileControlsInResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileControlsInResponse)
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_requirement", &obj.ControlRequirement)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_library_version", &obj.ControlLibraryVersion)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_name", &obj.ControlName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_description", &obj.ControlDescription)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_severity", &obj.ControlSeverity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_category", &obj.ControlCategory)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_parent", &obj.ControlParent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_specifications_count", &obj.ControlSpecificationsCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_docs", &obj.ControlDocs, UnmarshalControlDoc)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalControlSpecification)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileControlsPrototype : The control details of a profile.
type ProfileControlsPrototype struct {
	// The ID of the control library that contains the profile.
	ControlLibraryID *string `json:"control_library_id,omitempty"`

	// The control ID.
	ControlID *string `json:"control_id,omitempty"`

	// List of control specifications in a profile.
	ControlSpecifications []ProfileControlSpecificationPrototype `json:"control_specifications,omitempty"`
}

// UnmarshalProfileControlsPrototype unmarshals an instance of ProfileControlsPrototype from the specified map of raw messages.
func UnmarshalProfileControlsPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileControlsPrototype)
	err = core.UnmarshalPrimitive(m, "control_library_id", &obj.ControlLibraryID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "control_id", &obj.ControlID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "control_specifications", &obj.ControlSpecifications, UnmarshalProfileControlSpecificationPrototype)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileDefaultParametersResponse : The default parameters of a profile.
type ProfileDefaultParametersResponse struct {
	// The ID of the Profile.
	ID *string `json:"id,omitempty"`

	// list of parameters given by default.
	DefaultParameters []DefaultParameters `json:"default_parameters,omitempty"`
}

// UnmarshalProfileDefaultParametersResponse unmarshals an instance of ProfileDefaultParametersResponse from the specified map of raw messages.
func UnmarshalProfileDefaultParametersResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileDefaultParametersResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "default_parameters", &obj.DefaultParameters, UnmarshalDefaultParameters)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProfileInfo : The profile information.
type ProfileInfo struct {
	// The profile ID.
	ID *string `json:"id,omitempty"`

	// The profile name.
	Name *string `json:"name,omitempty"`

	// The profile version.
	Version *string `json:"version,omitempty"`
}

// UnmarshalProfileInfo unmarshals an instance of ProfileInfo from the specified map of raw messages.
func UnmarshalProfileInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProfileInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// ProviderType : The provider type item.
type ProviderType struct {
	// The unique identifier of the provider type.
	ID *string `json:"id" validate:"required"`

	// The type of the provider type.
	Type *string `json:"type" validate:"required"`

	// The name of the provider type.
	Name *string `json:"name" validate:"required"`

	// The provider type description.
	Description *string `json:"description" validate:"required"`

	// A boolean that indicates whether the provider type is s2s-enabled.
	S2sEnabled *bool `json:"s2s_enabled" validate:"required"`

	// The maximum number of instances that can be created for the provider type.
	InstanceLimit *int64 `json:"instance_limit" validate:"required"`

	// The mode that is used to get results from provider (`PUSH` or `PULL`).
	Mode *string `json:"mode" validate:"required"`

	// The format of the results that a provider supports.
	DataType *string `json:"data_type" validate:"required"`

	// The icon of a provider in .svg format that is encoded as a base64 string.
	Icon *string `json:"icon" validate:"required"`

	// The label that is associated with the provider type.
	Label *LabelType `json:"label,omitempty"`

	// The attributes that are required when you're creating an instance of a provider type. The attributes field can have
	// multiple  keys in its value. Each of those keys has a value  object that includes the type, and display name as
	// keys. For example, `{type:"", display_name:""}`.
	// **NOTE;** If the provider type is s2s-enabled, which means that if the `s2s_enabled` field is set to `true`, then a
	// CRN field of type text is required in the attributes value object.
	Attributes map[string]AdditionalProperty `json:"attributes" validate:"required"`

	// Time at which resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Time at which resource was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalProviderType unmarshals an instance of ProviderType from the specified map of raw messages.
func UnmarshalProviderType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderType)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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
	err = core.UnmarshalPrimitive(m, "s2s_enabled", &obj.S2sEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_limit", &obj.InstanceLimit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data_type", &obj.DataType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "icon", &obj.Icon)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "label", &obj.Label, UnmarshalLabelType)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalAdditionalProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypeCollection : The provider types collection.
type ProviderTypeCollection struct {
	// The array of provder type.
	ProviderTypes []ProviderType `json:"provider_types,omitempty"`
}

// UnmarshalProviderTypeCollection unmarshals an instance of ProviderTypeCollection from the specified map of raw messages.
func UnmarshalProviderTypeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeCollection)
	err = core.UnmarshalModel(m, "provider_types", &obj.ProviderTypes, UnmarshalProviderType)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypeInstance : A provider type instance.
type ProviderTypeInstance struct {
	// The unique identifier of the provider type instance.
	ID *string `json:"id,omitempty"`

	// The type of the provider type.
	Type *string `json:"type,omitempty"`

	// The name of the provider type instance.
	Name *string `json:"name,omitempty"`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Time at which resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// Time at which resource was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`
}

// UnmarshalProviderTypeInstance unmarshals an instance of ProviderTypeInstance from the specified map of raw messages.
func UnmarshalProviderTypeInstance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeInstance)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attributes", &obj.Attributes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderTypeInstanceCollection : Provider types instances response.
type ProviderTypeInstanceCollection struct {
	// The array of instances for all provider types.
	ProviderTypeInstances []ProviderTypeInstance `json:"provider_type_instances,omitempty"`
}

// UnmarshalProviderTypeInstanceCollection unmarshals an instance of ProviderTypeInstanceCollection from the specified map of raw messages.
func UnmarshalProviderTypeInstanceCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderTypeInstanceCollection)
	err = core.UnmarshalModel(m, "provider_type_instances", &obj.ProviderTypeInstances, UnmarshalProviderTypeInstance)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceCustomControlLibraryOptions : The ReplaceCustomControlLibrary options.
type ReplaceCustomControlLibraryOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The control library ID.
	ControlLibraryID *string `json:"control_library_id" validate:"required,ne="`

	// The name of the control library.
	ControlLibraryName *string `json:"control_library_name" validate:"required"`

	// Details of the control library.
	ControlLibraryDescription *string `json:"control_library_description" validate:"required"`

	// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
	ControlLibraryType *string `json:"control_library_type" validate:"required"`

	// The revision number of the control library.
	ControlLibraryVersion *string `json:"control_library_version" validate:"required"`

	// The list of rules that the control library attempts to adhere to.
	Controls []Control `json:"controls" validate:"required"`

	// The unique identifier of the revision.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// The account id tied to billing.
	BssAccount *string `json:"bss_account,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceCustomControlLibraryOptions.ControlLibraryType property.
// Details that the control library is a user made(custom) or Security Compliance Center(predefined).
const (
	ReplaceCustomControlLibraryOptionsControlLibraryTypeCustomConst = "custom"
)

// NewReplaceCustomControlLibraryOptions : Instantiate ReplaceCustomControlLibraryOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceCustomControlLibraryOptions(instanceID string, controlLibraryID string, controlLibraryName string, controlLibraryDescription string, controlLibraryType string, controlLibraryVersion string, controls []Control) *ReplaceCustomControlLibraryOptions {
	return &ReplaceCustomControlLibraryOptions{
		InstanceID:                core.StringPtr(instanceID),
		ControlLibraryID:          core.StringPtr(controlLibraryID),
		ControlLibraryName:        core.StringPtr(controlLibraryName),
		ControlLibraryDescription: core.StringPtr(controlLibraryDescription),
		ControlLibraryType:        core.StringPtr(controlLibraryType),
		ControlLibraryVersion:     core.StringPtr(controlLibraryVersion),
		Controls:                  controls,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceCustomControlLibraryOptions) SetInstanceID(instanceID string) *ReplaceCustomControlLibraryOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetControlLibraryID : Allow user to set ControlLibraryID
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryID(controlLibraryID string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryID = core.StringPtr(controlLibraryID)
	return _options
}

// SetControlLibraryName : Allow user to set ControlLibraryName
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryName(controlLibraryName string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryName = core.StringPtr(controlLibraryName)
	return _options
}

// SetControlLibraryDescription : Allow user to set ControlLibraryDescription
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryDescription(controlLibraryDescription string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryDescription = core.StringPtr(controlLibraryDescription)
	return _options
}

// SetControlLibraryType : Allow user to set ControlLibraryType
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryType(controlLibraryType string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryType = core.StringPtr(controlLibraryType)
	return _options
}

// SetControlLibraryVersion : Allow user to set ControlLibraryVersion
func (_options *ReplaceCustomControlLibraryOptions) SetControlLibraryVersion(controlLibraryVersion string) *ReplaceCustomControlLibraryOptions {
	_options.ControlLibraryVersion = core.StringPtr(controlLibraryVersion)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *ReplaceCustomControlLibraryOptions) SetControls(controls []Control) *ReplaceCustomControlLibraryOptions {
	_options.Controls = controls
	return _options
}

// SetVersionGroupLabel: Allows user to set VersionGroupLabel
func (_options *ReplaceCustomControlLibraryOptions) SetVersionGroupLabel(versionGroupLabel string) *ReplaceCustomControlLibraryOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetBssAccount : Allow user to set BssAccount
func (_options *ReplaceCustomControlLibraryOptions) SetBssAccount(bssAccount string) *ReplaceCustomControlLibraryOptions {
	_options.BssAccount = core.StringPtr(bssAccount)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceCustomControlLibraryOptions) SetHeaders(param map[string]string) *ReplaceCustomControlLibraryOptions {
	options.Headers = param
	return options
}

// ReplaceProfileAttachmentOptions : The ReplaceProfileAttachment options.
type ReplaceProfileAttachmentOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The details to describe the profile attachment.
	Description *string `json:"description" validate:"required"`

	// The name of the Profile Attachment.
	Name *string `json:"name" validate:"required"`

	// Details how often a scan from a profile attachment is ran.
	Schedule *string `json:"schedule" validate:"required"`

	// Details the state of a profile attachment.
	Status *string `json:"status" validate:"required"`

	// The notification configuration of the attachment.
	Notifications *AttachmentNotifications `json:"notifications" validate:"required"`

	// The parameters associated with the profile attachment.
	AttachmentParameters []Parameter `json:"attachment_parameters" validate:"required"`

	// A list of scopes associated with a profile attachment.
	Scope []MultiCloudScopePayload `json:"scope" validate:"required"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceProfileAttachmentOptions : Instantiate ReplaceProfileAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceProfileAttachmentOptions(instanceID string, profileID string, attachmentID string, description string, name string, schedule string, status string, notifications *AttachmentNotifications, scope []MultiCloudScopePayload, attachmentParameters []Parameter) *ReplaceProfileAttachmentOptions {
	return &ReplaceProfileAttachmentOptions{
		InstanceID:           core.StringPtr(instanceID),
		ProfileID:            core.StringPtr(profileID),
		AttachmentID:         core.StringPtr(attachmentID),
		Description:          core.StringPtr(description),
		Name:                 core.StringPtr(name),
		Schedule:             core.StringPtr(schedule),
		Status:               core.StringPtr(status),
		Notifications:        notifications,
		AttachmentParameters: attachmentParameters,
		Scope:                scope,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceProfileAttachmentOptions) SetInstanceID(instanceID string) *ReplaceProfileAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ReplaceProfileAttachmentOptions) SetProfileID(profileID string) *ReplaceProfileAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *ReplaceProfileAttachmentOptions) SetAttachmentID(attachmentID string) *ReplaceProfileAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceProfileAttachmentOptions) SetDescription(description string) *ReplaceProfileAttachmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceProfileAttachmentOptions) SetName(name string) *ReplaceProfileAttachmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSchedule : Allow user to set Schedule
func (_options *ReplaceProfileAttachmentOptions) SetSchedule(schedule string) *ReplaceProfileAttachmentOptions {
	_options.Schedule = core.StringPtr(schedule)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ReplaceProfileAttachmentOptions) SetStatus(status string) *ReplaceProfileAttachmentOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetScope : Allow user to set Scope
func (_options *ReplaceProfileAttachmentOptions) SetScope(scope []MultiCloudScopePayload) *ReplaceProfileAttachmentOptions {
	_options.Scope = scope
	return _options
}

// SetNotifiations : Allow user to set Notifications
func (_options *ReplaceProfileAttachmentOptions) SetNotifications(notifications *AttachmentNotifications) *ReplaceProfileAttachmentOptions {
	_options.Notifications = notifications
	return _options
}

// SetAttachmentParameters : Allow user to set AttachmentParameters
func (_options *ReplaceProfileAttachmentOptions) SetAttachmentParameters(parameters []Parameter) *ReplaceProfileAttachmentOptions {
	_options.AttachmentParameters = parameters
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceProfileAttachmentOptions) SetAccountID(accountID string) *ReplaceProfileAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceProfileAttachmentOptions) SetHeaders(param map[string]string) *ReplaceProfileAttachmentOptions {
	options.Headers = param
	return options
}

// ReplaceProfileOptions : The ReplaceProfile options.
type ReplaceProfileOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The name of the profile.
	ProfileName *string `json:"profile_name" validate:"required,ne="`

	// The type of the profile.
	ProfileType *string `json:"profile_type" validate:"required,ne="`

	// A description of what the profile should represent.
	ProfileDescription *string `json:"profile_description,omitempty"`

	// The version of the profile.
	ProfileVersion *string `json:"profile_version,omitempty"`

	// Determines if the profile is up to date with the latest revisions.
	Latest *bool `json:"latest,omitempty"`

	// The unique identifier of the revision.
	VersionGroupLabel *string `json:"version_group_label,omitempty"`

	// List of controls associated with the profile.
	Controls []ProfileControlsPrototype `json:"controls,omitempty"`

	// The default values when using the profile.
	DefaultParameters []DefaultParametersPrototype `json:"default_parameters,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceProfileOptions : Instantiate ReplaceProfileOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceProfileOptions(instanceID string, profileID string, profileName string, profileDescription string) *ReplaceProfileOptions {
	return &ReplaceProfileOptions{
		InstanceID:         core.StringPtr(instanceID),
		ProfileDescription: core.StringPtr(profileDescription),
		ProfileID:          core.StringPtr(profileID),
		ProfileName:        core.StringPtr(profileName),
		ProfileType:        core.StringPtr("custom"),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceProfileOptions) SetInstanceID(instanceID string) *ReplaceProfileOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ReplaceProfileOptions) SetProfileID(profileID string) *ReplaceProfileOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetProfileName : Allow user to set ProfileName
func (_options *ReplaceProfileOptions) SetProfileName(profileName string) *ReplaceProfileOptions {
	_options.ProfileName = core.StringPtr(profileName)
	return _options
}

// SetProfileType : Allow user to set ProfileType
func (_options *ReplaceProfileOptions) SetProfileType(profileType string) *ReplaceProfileOptions {
	_options.ProfileType = core.StringPtr(profileType)
	return _options
}

// SetProfileDescription : Allow user to set ProfileDescription
func (_options *ReplaceProfileOptions) SetProfileDescription(profileDescription string) *ReplaceProfileOptions {
	_options.ProfileDescription = core.StringPtr(profileDescription)
	return _options
}

// SetProfileVersion : Allow user to set ProfileVersion
func (_options *ReplaceProfileOptions) SetProfileVersion(profileVersion string) *ReplaceProfileOptions {
	_options.ProfileVersion = core.StringPtr(profileVersion)
	return _options
}

// SetLatest : Allow user to set Latest
func (_options *ReplaceProfileOptions) SetLatest(latest bool) *ReplaceProfileOptions {
	_options.Latest = core.BoolPtr(latest)
	return _options
}

// SetVersionGroupLabel : Allow user to set VersionGroupLabel
func (_options *ReplaceProfileOptions) SetVersionGroupLabel(versionGroupLabel string) *ReplaceProfileOptions {
	_options.VersionGroupLabel = core.StringPtr(versionGroupLabel)
	return _options
}

// SetControls : Allow user to set Controls
func (_options *ReplaceProfileOptions) SetControls(controls []ProfileControlsPrototype) *ReplaceProfileOptions {
	_options.Controls = controls
	return _options
}

// SetDefaultParameters : Allow user to set DefaultParameters
func (_options *ReplaceProfileOptions) SetDefaultParameters(defaultParameters []DefaultParametersPrototype) *ReplaceProfileOptions {
	_options.DefaultParameters = defaultParameters
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceProfileOptions) SetAccountID(accountID string) *ReplaceProfileOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceProfileOptions) SetHeaders(param map[string]string) *ReplaceProfileOptions {
	options.Headers = param
	return options
}

// ReplaceProfileParametersOptions : The ReplaceProfileParameters options.
type ReplaceProfileParametersOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The ID of the Profile.
	ID *string `json:"id,omitempty"`

	// list of parameters given by default.
	DefaultParameters []DefaultParameters `json:"default_parameters,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceProfileParametersOptions : Instantiate ReplaceProfileParametersOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceProfileParametersOptions(instanceID string, profileID string) *ReplaceProfileParametersOptions {
	return &ReplaceProfileParametersOptions{
		InstanceID: core.StringPtr(instanceID),
		ProfileID:  core.StringPtr(profileID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceProfileParametersOptions) SetInstanceID(instanceID string) *ReplaceProfileParametersOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *ReplaceProfileParametersOptions) SetProfileID(profileID string) *ReplaceProfileParametersOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceProfileParametersOptions) SetID(id string) *ReplaceProfileParametersOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetDefaultParameters : Allow user to set DefaultParameters
func (_options *ReplaceProfileParametersOptions) SetDefaultParameters(defaultParameters []DefaultParameters) *ReplaceProfileParametersOptions {
	_options.DefaultParameters = defaultParameters
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceProfileParametersOptions) SetAccountID(accountID string) *ReplaceProfileParametersOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceProfileParametersOptions) SetHeaders(param map[string]string) *ReplaceProfileParametersOptions {
	options.Headers = param
	return options
}

// ReplaceRuleOptions : The ReplaceRule options.
type ReplaceRuleOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of a rule/assessment.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// This field compares a supplied `Etag` value with the version that is stored for the requested resource. If the
	// values match, the server allows the request method to continue.
	//
	// To find the `Etag` value, run a GET request on the resource that you want to modify, and check the response headers.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The rule description.
	Description *string `json:"description" validate:"required"`

	// The rule target.
	Target *RuleTargetPrototype `json:"target" validate:"required"`

	// The required configurations for a Rule.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The rule version number.
	Version *string `json:"version,omitempty"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The list of labels that correspond to a rule.
	Labels []string `json:"labels,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceRuleOptions : Instantiate ReplaceRuleOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceRuleOptions(instanceID string, ruleID string, ifMatch string, description string, target *RuleTargetPrototype, requiredConfig RequiredConfigIntf) *ReplaceRuleOptions {
	return &ReplaceRuleOptions{
		InstanceID:     core.StringPtr(instanceID),
		RuleID:         core.StringPtr(ruleID),
		IfMatch:        core.StringPtr(ifMatch),
		Description:    core.StringPtr(description),
		Target:         target,
		RequiredConfig: requiredConfig,
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceRuleOptions) SetInstanceID(instanceID string) *ReplaceRuleOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetRuleID : Allow user to set RuleID
func (_options *ReplaceRuleOptions) SetRuleID(ruleID string) *ReplaceRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceRuleOptions) SetIfMatch(ifMatch string) *ReplaceRuleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceRuleOptions) SetDescription(description string) *ReplaceRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTarget : Allow user to set Target
func (_options *ReplaceRuleOptions) SetTarget(target *RuleTargetPrototype) *ReplaceRuleOptions {
	_options.Target = target
	return _options
}

// SetRequiredConfig : Allow user to set RequiredConfig
func (_options *ReplaceRuleOptions) SetRequiredConfig(requiredConfig RequiredConfigIntf) *ReplaceRuleOptions {
	_options.RequiredConfig = requiredConfig
	return _options
}

// SetVersion : Allow user to set Version
func (_options *ReplaceRuleOptions) SetVersion(version string) *ReplaceRuleOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetImport : Allow user to set Import
func (_options *ReplaceRuleOptions) SetImport(importVar *Import) *ReplaceRuleOptions {
	_options.Import = importVar
	return _options
}

// SetLabels : Allow user to set Labels
func (_options *ReplaceRuleOptions) SetLabels(labels []string) *ReplaceRuleOptions {
	_options.Labels = labels
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceRuleOptions) SetHeaders(param map[string]string) *ReplaceRuleOptions {
	options.Headers = param
	return options
}

// ReplaceTargetOptions : The ReplaceTarget options.
type ReplaceTargetOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The target ID.
	TargetID *string `json:"target_id" validate:"required,ne="`

	// The target account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The trusted profile ID.
	TrustedProfileID *string `json:"trusted_profile_id" validate:"required"`

	// The target name.
	Name *string `json:"name" validate:"required"`

	// Customer credential to access for a specific service to scan.
	Credentials []Credential `json:"credentials,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceTargetOptions : Instantiate ReplaceTargetOptions
func (*SecurityAndComplianceCenterApiV3) NewReplaceTargetOptions(instanceID string, targetID string, accountID string, trustedProfileID string, name string) *ReplaceTargetOptions {
	return &ReplaceTargetOptions{
		InstanceID:       core.StringPtr(instanceID),
		TargetID:         core.StringPtr(targetID),
		AccountID:        core.StringPtr(accountID),
		TrustedProfileID: core.StringPtr(trustedProfileID),
		Name:             core.StringPtr(name),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *ReplaceTargetOptions) SetInstanceID(instanceID string) *ReplaceTargetOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetTargetID : Allow user to set TargetID
func (_options *ReplaceTargetOptions) SetTargetID(targetID string) *ReplaceTargetOptions {
	_options.TargetID = core.StringPtr(targetID)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceTargetOptions) SetAccountID(accountID string) *ReplaceTargetOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetTrustedProfileID : Allow user to set TrustedProfileID
func (_options *ReplaceTargetOptions) SetTrustedProfileID(trustedProfileID string) *ReplaceTargetOptions {
	_options.TrustedProfileID = core.StringPtr(trustedProfileID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceTargetOptions) SetName(name string) *ReplaceTargetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCredentials : Allow user to set Credentials
func (_options *ReplaceTargetOptions) SetCredentials(credentials []Credential) *ReplaceTargetOptions {
	_options.Credentials = credentials
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTargetOptions) SetHeaders(param map[string]string) *ReplaceTargetOptions {
	options.Headers = param
	return options
}

// Report : The report.
type Report struct {
	// The ID of the report.
	ID *string `json:"id" validate:"required"`

	// The type of the scan.
	Type *string `json:"type" validate:"required"`

	// The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.
	GroupID *string `json:"group_id" validate:"required"`

	// The date when the report was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The date when the scan was run.
	ScanTime *string `json:"scan_time" validate:"required"`

	// The Cloud Object Storage object that is associated with the report.
	CosObject *string `json:"cos_object" validate:"required"`

	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The account that is associated with a report.
	Account *Account `json:"account" validate:"required"`

	// The profile information.
	Profile *ProfileInfo `json:"profile" validate:"required"`

	// The scope ID that is associated with a report. Attributes for this object will be blank if the report has multiple
	// scopes tied to the report.
	Scope *ScopeID `json:"scope" validate:"required"`

	// The attachment that is associated with a report.
	Attachment *Attachment `json:"attachment" validate:"required"`

	// The compliance stats.
	ControlsSummary *ComplianceStatsWithNonCompliant `json:"controls_summary" validate:"required"`

	// The evaluation stats.
	EvaluationsSummary *EvalStats `json:"evaluations_summary" validate:"required"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags" validate:"required"`

	// The scopes used in the report.
	Scopes []ReportScope `json:"scopes" validate:"required"`
}

// UnmarshalReport unmarshals an instance of Report from the specified map of raw messages.
func UnmarshalReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Report)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "group_id", &obj.GroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_time", &obj.ScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos_object", &obj.CosObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "profile", &obj.Profile, UnmarshalProfileInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scope", &obj.Scope, UnmarshalScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "attachment", &obj.Attachment, UnmarshalAttachment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls_summary", &obj.ControlsSummary, UnmarshalComplianceStatsWithNonCompliant)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations_summary", &obj.EvaluationsSummary, UnmarshalEvalStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scopes", &obj.Scopes, UnmarshalReportScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportCollection : The page of reports.
type ReportCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The list of reports that are on the page.
	Reports []Report `json:"reports,omitempty"`
}

// UnmarshalReportCollection unmarshals an instance of ReportCollection from the specified map of raw messages.
func UnmarshalReportCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reports", &obj.Reports, UnmarshalReport)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ReportCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ReportControls : The list of controls.
type ReportControls struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The list of controls that are in the report.
	Controls []ControlWithStats `json:"controls,omitempty"`
}

// UnmarshalReportControls unmarshals an instance of ReportControls from the specified map of raw messages.
func UnmarshalReportControls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportControls)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalControlWithStats)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportLatest : The response body of the `get_latest_reports` operation.
type ReportLatest struct {
	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The compliance stats.
	ControlsSummary *ComplianceStats `json:"controls_summary,omitempty"`

	// The evaluation stats.
	EvaluationsSummary *EvalStats `json:"evaluations_summary,omitempty"`

	// The compliance score.
	Score *ComplianceScore `json:"score,omitempty"`

	// The list of reports.
	Reports []Report `json:"reports,omitempty"`
}

// UnmarshalReportLatest unmarshals an instance of ReportLatest from the specified map of raw messages.
func UnmarshalReportLatest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportLatest)
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls_summary", &obj.ControlsSummary, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations_summary", &obj.EvaluationsSummary, UnmarshalEvalStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "score", &obj.Score, UnmarshalComplianceScore)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "reports", &obj.Reports, UnmarshalReport)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportScope : The scopes used in the report.
type ReportScope struct {
	// The ID of the scope used.
	ID *string `json:"id" validate:"required"`

	// The name of the scope used.
	Name *string `json:"name" validate:"required"`

	// The url to a report concerning the specified scope.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalReportScope unmarshals an instance of ReportScope from the specified map of raw messages.
func UnmarshalReportScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportScope)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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

// ReportSummary : The report summary.
type ReportSummary struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// Instance ID.
	IsntanceID *string `json:"isntance_id,omitempty"`

	// The account that is associated with a report.
	Account *Account `json:"account,omitempty"`

	// The compliance score.
	Score *ComplianceScore `json:"score,omitempty"`

	// The evaluation stats.
	Evaluations *EvalStats `json:"evaluations,omitempty"`

	// The compliance stats.
	Controls *ComplianceStats `json:"controls,omitempty"`

	// The resource summary.
	Resources *ResourceSummary `json:"resources,omitempty"`
}

// UnmarshalReportSummary unmarshals an instance of ReportSummary from the specified map of raw messages.
func UnmarshalReportSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportSummary)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "isntance_id", &obj.IsntanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "score", &obj.Score, UnmarshalComplianceScore)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "evaluations", &obj.Evaluations, UnmarshalEvalStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls", &obj.Controls, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResourceSummary)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportTags : The response body of the `get_tags` operation.
type ReportTags struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`
}

// UnmarshalReportTags unmarshals an instance of ReportTags from the specified map of raw messages.
func UnmarshalReportTags(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportTags)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportViolationDataPoint : The report violation data point.
type ReportViolationDataPoint struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.
	ReportGroupID *string `json:"report_group_id,omitempty"`

	// The date when the scan was run.
	ScanTime *string `json:"scan_time,omitempty"`

	// The compliance stats.
	ControlsSummary *ComplianceStats `json:"controls_summary,omitempty"`
}

// UnmarshalReportViolationDataPoint unmarshals an instance of ReportViolationDataPoint from the specified map of raw messages.
func UnmarshalReportViolationDataPoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportViolationDataPoint)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_group_id", &obj.ReportGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_time", &obj.ScanTime)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "controls_summary", &obj.ControlsSummary, UnmarshalComplianceStats)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReportViolationsDrift : The response body of the `get_report_violations_drift` operation.
type ReportViolationsDrift struct {
	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The ID of the report group.
	ReportGroupID *string `json:"report_group_id,omitempty"`

	// The list of report violations data points.
	DataPoints []ReportViolationDataPoint `json:"data_points,omitempty"`
}

// UnmarshalReportViolationsDrift unmarshals an instance of ReportViolationsDrift from the specified map of raw messages.
func UnmarshalReportViolationsDrift(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ReportViolationsDrift)
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_group_id", &obj.ReportGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "data_points", &obj.DataPoints, UnmarshalReportViolationDataPoint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfig : The required configurations for a Rule.
// Models which "extend" this model:
// - RequiredConfigConditionBase
// - RequiredConfigConditionList
// - RequiredConfigConditionSubRule
type RequiredConfig struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property,omitempty"`

	// The operator.
	Operator *string `json:"operator,omitempty"`

	Value interface{} `json:"value,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []RequiredConfigIntf `json:"or,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []RequiredConfigIntf `json:"and,omitempty"`

	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`

	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

// Constants associated with the RequiredConfig.Operator property.
// The operator.
const (
	RequiredConfigOperatorDaysLessThanConst         = "days_less_than"
	RequiredConfigOperatorIpsEqualsConst            = "ips_equals"
	RequiredConfigOperatorIpsInRangeConst           = "ips_in_range"
	RequiredConfigOperatorIpsNotEqualsConst         = "ips_not_equals"
	RequiredConfigOperatorIsEmptyConst              = "is_empty"
	RequiredConfigOperatorIsFalseConst              = "is_false"
	RequiredConfigOperatorIsNotEmptyConst           = "is_not_empty"
	RequiredConfigOperatorIsTrueConst               = "is_true"
	RequiredConfigOperatorNumEqualsConst            = "num_equals"
	RequiredConfigOperatorNumGreaterThanConst       = "num_greater_than"
	RequiredConfigOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RequiredConfigOperatorNumLessThanConst          = "num_less_than"
	RequiredConfigOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RequiredConfigOperatorNumNotEqualsConst         = "num_not_equals"
	RequiredConfigOperatorStringContainsConst       = "string_contains"
	RequiredConfigOperatorStringEqualsConst         = "string_equals"
	RequiredConfigOperatorStringMatchConst          = "string_match"
	RequiredConfigOperatorStringNotContainsConst    = "string_not_contains"
	RequiredConfigOperatorStringNotEqualsConst      = "string_not_equals"
	RequiredConfigOperatorStringNotMatchConst       = "string_not_match"
	RequiredConfigOperatorStringsAllowedConst       = "strings_allowed"
	RequiredConfigOperatorStringsInListConst        = "strings_in_list"
	RequiredConfigOperatorStringsRequiredConst      = "strings_required"
)

func (*RequiredConfig) isaRequiredConfig() bool {
	return true
}

type RequiredConfigIntf interface {
	isaRequiredConfig() bool
}

// UnmarshalRequiredConfig unmarshals an instance of RequiredConfig from the specified map of raw messages.
func UnmarshalRequiredConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfig)
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
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalRequiredConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalRequiredConfig)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Resource : The resource.
type Resource struct {
	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The resource CRN.
	ID *string `json:"id,omitempty"`

	// The resource name.
	ResourceName *string `json:"resource_name,omitempty"`

	// The account that is associated with a report.
	Account *Account `json:"account,omitempty"`

	// The ID of the component.
	ComponentID *string `json:"component_id,omitempty"`

	// The name of the component.
	ComponentName *string `json:"component_name,omitempty"`

	// The environment.
	Environment *string `json:"environment,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The number of assessments with no corresponding evaluations.
	SkippedCount *int64 `json:"skipped_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`

	// The name of the service.
	ServiceName *string `json:"service_name,omitempty"`

	// The instance CRN.
	InstanceCRN *string `json:"instance_crn,omitempty"`
}

// Constants associated with the Resource.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ResourceStatusCompliantConst              = "compliant"
	ResourceStatusNotApplicableConst          = "not_applicable"
	ResourceStatusNotCompliantConst           = "not_compliant"
	ResourceStatusUnableToPerformConst        = "unable_to_perform"
	ResourceStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "account", &obj.Account, UnmarshalAccount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_id", &obj.ComponentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "component_name", &obj.ComponentName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped_count", &obj.SkippedCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_crn", &obj.InstanceCRN)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourcePage : The page of resource evaluation summaries.
type ResourcePage struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The ID of the report.
	ReportID *string `json:"report_id,omitempty"`

	// The ID of the home account.
	HomeAccountID *string `json:"home_account_id,omitempty"`

	// The list of resource evaluation summaries that are on the page.
	Resources []Resource `json:"resources,omitempty"`
}

// UnmarshalResourcePage unmarshals an instance of ResourcePage from the specified map of raw messages.
func UnmarshalResourcePage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourcePage)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "report_id", &obj.ReportID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "home_account_id", &obj.HomeAccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ResourcePage) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ResourceSummary : The resource summary.
type ResourceSummary struct {
	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of checks.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of compliant checks.
	CompliantCount *int64 `json:"compliant_count,omitempty"`

	// The number of checks that are not compliant.
	NotCompliantCount *int64 `json:"not_compliant_count,omitempty"`

	// The number of checks that are unable to perform.
	UnableToPerformCount *int64 `json:"unable_to_perform_count,omitempty"`

	// The number of checks that require a user evaluation.
	UserEvaluationRequiredCount *int64 `json:"user_evaluation_required_count,omitempty"`

	// The number of not applicable (with no evaluations) checks.
	NotApplicableCount *int64 `json:"not_applicable_count,omitempty"`

	// The top 10 resources that have the most failures.
	TopFailed []ResourceSummaryItem `json:"top_failed,omitempty"`
}

// Constants associated with the ResourceSummary.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ResourceSummaryStatusCompliantConst              = "compliant"
	ResourceSummaryStatusNotApplicableConst          = "not_applicable"
	ResourceSummaryStatusNotCompliantConst           = "not_compliant"
	ResourceSummaryStatusUnableToPerformConst        = "unable_to_perform"
	ResourceSummaryStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalResourceSummary unmarshals an instance of ResourceSummary from the specified map of raw messages.
func UnmarshalResourceSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceSummary)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "compliant_count", &obj.CompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_compliant_count", &obj.NotCompliantCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "unable_to_perform_count", &obj.UnableToPerformCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_evaluation_required_count", &obj.UserEvaluationRequiredCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "not_applicable_count", &obj.NotApplicableCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "top_failed", &obj.TopFailed, UnmarshalResourceSummaryItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceSummaryItem : The resource summary item.
type ResourceSummaryItem struct {
	// The resource ID.
	ID *string `json:"id,omitempty"`

	// The resource name.
	Name *string `json:"name,omitempty"`

	// The account that owns the resource.
	Account *string `json:"account,omitempty"`

	// The service that is managing the resource.
	Service *string `json:"service,omitempty"`

	// The services display name that is managing the resource.
	ServiceDisplayName *string `json:"service_display_name,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`

	// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
	Status *string `json:"status,omitempty"`

	// The total number of evaluations.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The number of passed evaluations.
	PassCount *int64 `json:"pass_count,omitempty"`

	// The number of failed evaluations.
	FailureCount *int64 `json:"failure_count,omitempty"`

	// The number of evaluations that started, but did not finish, and ended with errors.
	ErrorCount *int64 `json:"error_count,omitempty"`

	// The number of assessments with no corresponding evaluations.
	SkippedCount *int64 `json:"skipped_count,omitempty"`

	// The total number of completed evaluations.
	CompletedCount *int64 `json:"completed_count,omitempty"`
}

// Constants associated with the ResourceSummaryItem.Status property.
// The allowed values of an aggregated status for controls, specifications, assessments, and resources.
const (
	ResourceSummaryItemStatusCompliantConst              = "compliant"
	ResourceSummaryItemStatusNotApplicableConst          = "not_applicable"
	ResourceSummaryItemStatusNotCompliantConst           = "not_compliant"
	ResourceSummaryItemStatusUnableToPerformConst        = "unable_to_perform"
	ResourceSummaryItemStatusUserEvaluationRequiredConst = "user_evaluation_required"
)

// UnmarshalResourceSummaryItem unmarshals an instance of ResourceSummaryItem from the specified map of raw messages.
func UnmarshalResourceSummaryItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceSummaryItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account", &obj.Account)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_display_name", &obj.ServiceDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pass_count", &obj.PassCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_count", &obj.FailureCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_count", &obj.ErrorCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "skipped_count", &obj.SkippedCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "completed_count", &obj.CompletedCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rule : The rule response that corresponds to an account instance.
type Rule struct {
	// The date when the rule was created.
	CreatedOn *strfmt.DateTime `json:"created_on" validate:"required"`

	// The user who created the rule.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the rule was modified.
	UpdatedOn *strfmt.DateTime `json:"updated_on" validate:"required"`

	// The user who modified the rule.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The rule ID.
	ID *string `json:"id" validate:"required"`

	// The account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The details of a rule's response.
	Description *string `json:"description" validate:"required"`

	// The rule type (allowable values are `user_defined` or `system_defined`).
	Type *string `json:"type" validate:"required"`

	// The version number of a rule.
	Version *string `json:"version" validate:"required"`

	// The collection of import parameters.
	Import *Import `json:"import,omitempty"`

	// The rule target.
	Target *RuleTarget `json:"target" validate:"required"`

	// The required configurations for a Rule.
	RequiredConfig RequiredConfigIntf `json:"required_config" validate:"required"`

	// The list of labels.
	Labels []string `json:"labels" validate:"required"`
}

// Constants associated with the Rule.Type property.
// The rule type (allowable values are `user_defined` or `system_defined`).
const (
	RuleTypeSystemDefinedConst = "system_defined"
	RuleTypeUserDefinedConst   = "user_defined"
)

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
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
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "import", &obj.Import, UnmarshalImport)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalRuleTarget)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_config", &obj.RequiredConfig, UnmarshalRequiredConfig)
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

// RuleCollection : The page of rules.
type RuleCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The collection of rules that correspond to an account instance. Maximum of 100/500 custom rules per
	// stand-alone/enterprise account.
	Rules []Rule `json:"rules,omitempty"`
}

// UnmarshalRuleCollection unmarshals an instance of RuleCollection from the specified map of raw messages.
func UnmarshalRuleCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
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

// Retrieve the value to be passed to a request to access the next page of results
func (resp *RuleCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// RuleInfo : The rule.
type RuleInfo struct {
	// The rule ID.
	ID *string `json:"id,omitempty"`

	// The rule type.
	Type *string `json:"type,omitempty"`

	// The rule description.
	Description *string `json:"description,omitempty"`

	// The rule version.
	Version *string `json:"version,omitempty"`

	// The rule account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The date when the rule was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The ID of the user who created the rule.
	CreatedBy *string `json:"created_by,omitempty"`

	// The date when the rule was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`

	// The ID of the user who updated the rule.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The rule labels.
	Labels []string `json:"labels,omitempty"`
}

// UnmarshalRuleInfo unmarshals an instance of RuleInfo from the specified map of raw messages.
func UnmarshalRuleInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleInfo)
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
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
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

// RuleParameter : The rule import parameter.
type RuleParameter struct {
	// The import parameter name.
	Name *string `json:"name,omitempty"`

	// The display name of the property.
	DisplayName *string `json:"display_name,omitempty"`

	// The propery description.
	Description *string `json:"description,omitempty"`

	// The property type.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the RuleParameter.Type property.
// The property type.
const (
	RuleParameterTypeBooleanConst    = "boolean"
	RuleParameterTypeGeneralConst    = "general"
	RuleParameterTypeIPListConst     = "ip_list"
	RuleParameterTypeNumericConst    = "numeric"
	RuleParameterTypeStringConst     = "string"
	RuleParameterTypeStringListConst = "string_list"
	RuleParameterTypeTimestampConst  = "timestamp"
)

// UnmarshalRuleParameter unmarshals an instance of RuleParameter from the specified map of raw messages.
func UnmarshalRuleParameter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleParameter)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleProperty : The supported config property.
type RuleProperty struct {
	// The property name.
	Name *string `json:"name,omitempty"`

	// The property description.
	Description *string `json:"description,omitempty"`

	// The operator kind used when evaluating a property.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the RuleProperty.Type property.
// The operator kind used when evaluating a property.
const (
	RulePropertyTypeBooleanConst    = "boolean"
	RulePropertyTypeGeneralConst    = "general"
	RulePropertyTypeIPListConst     = "ip_list"
	RulePropertyTypeNumericConst    = "numeric"
	RulePropertyTypeStringConst     = "string"
	RulePropertyTypeStringListConst = "string_list"
	RulePropertyTypeTimestampConst  = "timestamp"
)

// UnmarshalRuleProperty unmarshals an instance of RuleProperty from the specified map of raw messages.
func UnmarshalRuleProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleProperty)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleTarget : The rule target.
type RuleTarget struct {
	// The target service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The display name of the target service.
	ServiceDisplayName *string `json:"service_display_name,omitempty"`

	// The target resource kind.
	ResourceKind *string `json:"resource_kind" validate:"required"`

	// The reference name used
	Ref *string `json:"ref,omitempty"`

	// The additional target attributes used to filter to a subset of resources.
	AdditionalTargetAttributes []AdditionalTargetAttribute `json:"additional_target_attributes,omitempty"`
}

// NewRuleTarget : Instantiate RuleTarget (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewRuleTarget(serviceName string, resourceKind string) (_model *RuleTarget, err error) {
	_model = &RuleTarget{
		ServiceName:  core.StringPtr(serviceName),
		ResourceKind: core.StringPtr(resourceKind),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleTarget unmarshals an instance of RuleTarget from the specified map of raw messages.
func UnmarshalRuleTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleTarget)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_display_name", &obj.ServiceDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ref", &obj.Ref)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalAdditionalTargetAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleTargetPrototype : The rule target.
type RuleTargetPrototype struct {
	// The target service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The target resource kind.
	ResourceKind *string `json:"resource_kind" validate:"required"`

	// The reference name used
	Ref *string `json:"ref,omitempty"`

	// The additional target attributes used to filter to a subset of resources.
	AdditionalTargetAttributes []AdditionalTargetAttribute `json:"additional_target_attributes,omitempty"`
}

// NewRuleTargetPrototype : Instantiate RuleTargetPrototype (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewRuleTargetPrototype(serviceName string, resourceKind string) (_model *RuleTargetPrototype, err error) {
	_model = &RuleTargetPrototype{
		ServiceName:  core.StringPtr(serviceName),
		ResourceKind: core.StringPtr(resourceKind),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRuleTargetPrototype unmarshals an instance of RuleTargetPrototype from the specified map of raw messages.
func UnmarshalRuleTargetPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleTargetPrototype)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalAdditionalTargetAttribute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScanReport : A report detailing the evaluations related to a specific control.
type ScanReport struct {
	// The ID of the scan report.
	ID *string `json:"id,omitempty"`

	// The ID of the scan.
	ScanID *string `json:"scan_id,omitempty"`

	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id,omitempty"`

	// The ID of the scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The ID of the sub-scope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// The enum of different scan report status.
	Status *string `json:"status,omitempty"`

	// The date when the report was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The file type of the report.
	Format *string `json:"format" validate:"required"`

	// The URL of the scan report.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the ScanReport.Status property.
// The enum of different scan report status.
const (
	ScanReportStatusCompletedConst  = "completed"
	ScanReportStatusDeletedConst    = "deleted"
	ScanReportStatusErrorConst      = "error"
	ScanReportStatusInProgressConst = "in_progress"
	ScanReportStatusPendingConst    = "pending"
)

// Constants associated with the ScanReport.Format property.
// The file type of the report.
const (
	ScanReportFormatCSVConst = "csv"
	ScanReportFormatPDFConst = "pdf"
)

// UnmarshalScanReport unmarshals an instance of ScanReport from the specified map of raw messages.
func UnmarshalScanReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScanReport)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scan_id", &obj.ScanID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "subscope_id", &obj.SubscopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
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

// ScanReportCollection : The page of scan reports.
type ScanReportCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The id of the requested scope.
	ScopeID *string `json:"scope_id,omitempty"`

	// The id of the requested subscope.
	SubscopeID *string `json:"subscope_id,omitempty"`

	// The list of scan reports.
	ScanReports []ScanReport `json:"scan_reports,omitempty"`
}

// UnmarshalScanReportCollection unmarshals an instance of ScanReportCollection from the specified map of raw messages.
func UnmarshalScanReportCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScanReportCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "scope_id", &obj.ScopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "subscope_id", &obj.SubscopeID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scan_reports", &obj.ScanReports, UnmarshalScanReport)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Scope : The group of resources that you want to evaluate. In the new API-based architecture, a scope can be an Enterprise,
// Account group, Account, or Resource group.
type Scope struct {
	// The ID of the scope.
	ID *string `json:"id" validate:"required"`

	// The scope name.
	Name *string `json:"name" validate:"required"`

	// The scope description.
	Description *string `json:"description" validate:"required"`

	// The scope environment. This value details what cloud provider the scope targets.
	Environment *string `json:"environment" validate:"required"`

	// The properties that are supported for scoping by this environment.
	Properties []ScopePropertyIntf `json:"properties" validate:"required"`

	// The ID of the account associated with the scope.
	AccountID *string `json:"account_id" validate:"required"`

	// The ID of the instance associated with the scope.
	InstanceID *string `json:"instance_id" validate:"required"`

	// The identifier of the account or service ID who created the scope.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the scope was created.
	CreatedOn *strfmt.DateTime `json:"created_on" validate:"required"`

	// The ID of the user or service ID who updated the scope.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The date when the scope was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on" validate:"required"`

	// The number of attachments tied to the scope.
	AttachmentCount *float64 `json:"attachment_count" validate:"required"`
}

// UnmarshalScope unmarshals an instance of Scope from the specified map of raw messages.
func UnmarshalScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Scope)
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
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalScopeProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "attachment_count", &obj.AttachmentCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopeCollection : A list of scopes.
type ScopeCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The array of scopes.
	Scopes []Scope `json:"scopes,omitempty"`
}

// UnmarshalScopeCollection unmarshals an instance of ScopeCollection from the specified map of raw messages.
func UnmarshalScopeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scopes", &obj.Scopes, UnmarshalScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ScopeCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	if core.IsNil(resp.Next.Start) {
		return nil, nil
	}
	if *resp.Next.Start == "" {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ScopeID : The scope ID that is associated with a report. Attributes for this object will be blank if the report has multiple
// scopes tied to the report.
type ScopeID struct {
	// The scope ID.
	ID *string `json:"id,omitempty"`

	// The scope type.
	Type *string `json:"type,omitempty"`
}

// UnmarshalScopeID unmarshals an instance of ScopeID from the specified map of raw messages.
func UnmarshalScopeID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeID)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// ScopeProperty : ScopeProperty struct
// Models which "extend" this model:
// - ScopePropertyScopeID
// - ScopePropertyScopeType
// - ScopePropertyExclusions
type ScopeProperty struct {
	// The key for the scope property.
	Name *string `json:"name,omitempty"`

	Value interface{} `json:"value,omitempty"`

	// A list of scopes/targets to exclude from a scope.
	Exclusions []ScopePropertyExclusionItem `json:"exclusions,omitempty"`
}

// Constants associated with the ScopeProperty.Name property.
// The key for the scope property.
const (
	ScopePropertyNameScopeIDConst = "scope_id"
)

func (*ScopeProperty) isaScopeProperty() bool {
	return true
}

type ScopePropertyIntf interface {
	isaScopeProperty() bool
}

// UnmarshalScopeProperty unmarshals an instance of ScopeProperty from the specified map of raw messages.
func UnmarshalScopeProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopeProperty)
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

// ScopePropertyExclusionItem : Any exclusion to be defined in the scope.
type ScopePropertyExclusionItem struct {
	ScopeID *string `json:"scope_id,omitempty"`

	// The type of scope it targets
	//
	// The scope values are as followed:
	// - enterprise: The scope targets an enterprise account
	// - enterprise.account_group: The scope targets an account group within an enterprise
	// - enterprise.account: The scope targets an account within an enterprise
	// - account: The scope targets an account not tied to an enterprise
	// - account.resource_group: The scope targets a resource group within an account.
	ScopeType *string `json:"scope_type,omitempty"`
}

// Constants associated with the ScopePropertyExclusionItem.ScopeType property.
// The type of scope it targets
//
// The scope values are as followed:
// - enterprise: The scope targets an enterprise account
// - enterprise.account_group: The scope targets an account group within an enterprise
// - enterprise.account: The scope targets an account within an enterprise
// - account: The scope targets an account not tied to an enterprise
// - account.resource_group: The scope targets a resource group within an account.
const (
	ScopePropertyExclusionItemScopeTypeAccountConst                = "account"
	ScopePropertyExclusionItemScopeTypeAccountResourceGroupConst   = "account.resource_group"
	ScopePropertyExclusionItemScopeTypeEnterpriseConst             = "enterprise"
	ScopePropertyExclusionItemScopeTypeEnterpriseAccountConst      = "enterprise.account"
	ScopePropertyExclusionItemScopeTypeEnterpriseAccountGroupConst = "enterprise.account_group"
)

// UnmarshalScopePropertyExclusionItem unmarshals an instance of ScopePropertyExclusionItem from the specified map of raw messages.
func UnmarshalScopePropertyExclusionItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopePropertyExclusionItem)
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

// ScopePrototype : The request body to make a Scope.
type ScopePrototype struct {
	// The scope name.
	Name *string `json:"name,omitempty"`

	// The scope description.
	Description *string `json:"description,omitempty"`

	// The scope environment.
	Environment *string `json:"environment,omitempty"`

	// The properties that are supported for scoping by this environment.
	Properties []ScopePropertyIntf `json:"properties,omitempty"`
}

// UnmarshalScopePrototype unmarshals an instance of ScopePrototype from the specified map of raw messages.
func UnmarshalScopePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopePrototype)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalScopeProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Service : The response body for creating a service instance.
type Service struct {
	// The service creation date.
	CreatedOn *strfmt.DateTime `json:"created_on" validate:"required"`

	// The service author.
	CreatedBy *string `json:"created_by" validate:"required"`

	// The date when the service was modified.
	UpdatedOn *strfmt.DateTime `json:"updated_on" validate:"required"`

	// The user who modified the service.
	UpdatedBy *string `json:"updated_by" validate:"required"`

	// The service name.
	ServiceName *string `json:"service_name" validate:"required"`

	// The display name of the service.
	ServiceDisplayName *string `json:"service_display_name,omitempty"`

	// The service description.
	Description *string `json:"description" validate:"required"`

	// The indication of whether monitoring is enabled.
	MonitoringEnabled *bool `json:"monitoring_enabled" validate:"required"`

	// The indication of whether enforcement is enabled.
	EnforcementEnabled *bool `json:"enforcement_enabled" validate:"required"`

	// The indication of whether service listing is enabled.
	ServiceListingEnabled *bool `json:"service_listing_enabled" validate:"required"`

	// The service configuration information.
	ConfigInformationPoint *ConfigurationInformationPoints `json:"config_information_point" validate:"required"`

	// The supported configurations.
	SupportedConfigs []SupportedConfigs `json:"supported_configs" validate:"required"`
}

// UnmarshalService unmarshals an instance of Service from the specified map of raw messages.
func UnmarshalService(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Service)
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_display_name", &obj.ServiceDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monitoring_enabled", &obj.MonitoringEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enforcement_enabled", &obj.EnforcementEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_listing_enabled", &obj.ServiceListingEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "config_information_point", &obj.ConfigInformationPoint, UnmarshalConfigurationInformationPoints)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "supported_configs", &obj.SupportedConfigs, UnmarshalSupportedConfigs)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceCollection : The services.
type ServiceCollection struct {
	// The list of services.
	Services []Service `json:"services,omitempty"`
}

// UnmarshalServiceCollection unmarshals an instance of ServiceCollection from the specified map of raw messages.
func UnmarshalServiceCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceCollection)
	err = core.UnmarshalModel(m, "services", &obj.Services, UnmarshalService)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Settings : The settings.
type Settings struct {
	// The Event Notifications settings.
	EventNotifications *EventNotifications `json:"event_notifications,omitempty"`

	// The Cloud Object Storage settings.
	ObjectStorage *ObjectStorage `json:"object_storage,omitempty"`
}

// UnmarshalSettings unmarshals an instance of Settings from the specified map of raw messages.
func UnmarshalSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Settings)
	err = core.UnmarshalModel(m, "event_notifications", &obj.EventNotifications, UnmarshalEventNotifications)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "object_storage", &obj.ObjectStorage, UnmarshalObjectStorage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubRule : A rule within a rule used in the requiredConfig.
type SubRule struct {
	// The rule target.
	Target *RuleTarget `json:"target,omitempty"`

	// The required configurations for a Rule.
	RequiredConfig RequiredConfigIntf `json:"required_config,omitempty"`
}

// UnmarshalSubRule unmarshals an instance of SubRule from the specified map of raw messages.
func UnmarshalSubRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubRule)
	err = core.UnmarshalModel(m, "target", &obj.Target, UnmarshalRuleTarget)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "required_config", &obj.RequiredConfig, UnmarshalRequiredConfig)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubScope : A segment of a scope. Subscopes are used to ensure that the members of your teams who review results only have access
// to the information regarding the instances that they have access to.
type SubScope struct {
	// The Subscope ID.
	ID *string `json:"id,omitempty"`

	// The name of the Subscope.
	Name *string `json:"name,omitempty"`

	// Text to describe the Subscope.
	Description *string `json:"description,omitempty"`

	// The virtual space where applications can be deployed and managed.
	Environment *string `json:"environment,omitempty"`

	// Additional attributes that are supported for scoping by this environment.
	Properties []ScopePropertyIntf `json:"properties,omitempty"`
}

// UnmarshalSubScope unmarshals an instance of SubScope from the specified map of raw messages.
func UnmarshalSubScope(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubScope)
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
	err = core.UnmarshalPrimitive(m, "environment", &obj.Environment)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalScopeProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SubScopeCollection : The response body of the subscopes.
type SubScopeCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The array of subscopes.
	Subscopes []SubScope `json:"subscopes,omitempty"`
}

// UnmarshalSubScopeCollection unmarshals an instance of SubScopeCollection from the specified map of raw messages.
func UnmarshalSubScopeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubScopeCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "subscopes", &obj.Subscopes, UnmarshalSubScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SubScopeCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// SubScopeResponse : The response body of the subscope.
type SubScopeResponse struct {
	// The array of subscopes.
	Subscopes []SubScope `json:"subscopes,omitempty"`
}

// UnmarshalSubScopeResponse unmarshals an instance of SubScopeResponse from the specified map of raw messages.
func UnmarshalSubScopeResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SubScopeResponse)
	err = core.UnmarshalModel(m, "subscopes", &obj.Subscopes, UnmarshalSubScope)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SupportedConfigs : The location information of supported configurations.
type SupportedConfigs struct {
	// The supported config resource kind.
	ResourceKind *string `json:"resource_kind,omitempty"`

	// The supported config list of additional target attributes.
	AdditionalTargetAttributes []AdditionalTargetAttribute `json:"additional_target_attributes,omitempty"`

	// The supported config list properties.
	Properties []RuleProperty `json:"properties,omitempty"`

	// The supported config description.
	Description *string `json:"description,omitempty"`

	// The indication of whether the configuration information point (CIP) requires a service instance.
	CipRequiresServiceInstance *bool `json:"cip_requires_service_instance,omitempty"`

	// The supported config resource group support.
	ResourceGroupSupport *bool `json:"resource_group_support,omitempty"`

	// The supported config tagging support.
	TaggingSupport *bool `json:"tagging_support,omitempty"`

	// The supported config inherited tags.
	InheritTags *bool `json:"inherit_tags,omitempty"`
}

// UnmarshalSupportedConfigs unmarshals an instance of SupportedConfigs from the specified map of raw messages.
func UnmarshalSupportedConfigs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedConfigs)
	err = core.UnmarshalPrimitive(m, "resource_kind", &obj.ResourceKind)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "additional_target_attributes", &obj.AdditionalTargetAttributes, UnmarshalAdditionalTargetAttribute)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalRuleProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cip_requires_service_instance", &obj.CipRequiresServiceInstance)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_support", &obj.ResourceGroupSupport)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tagging_support", &obj.TaggingSupport)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "inherit_tags", &obj.InheritTags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tags : The collection of different types of tags.
type Tags struct {
	// The collection of user tags.
	User []string `json:"user,omitempty"`

	// The collection of access tags.
	Access []string `json:"access,omitempty"`

	// The collection of service tags.
	Service []string `json:"service,omitempty"`
}

// UnmarshalTags unmarshals an instance of Tags from the specified map of raw messages.
func UnmarshalTags(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tags)
	err = core.UnmarshalPrimitive(m, "user", &obj.User)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "access", &obj.Access)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service", &obj.Service)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Target : The details of the target account.
type Target struct {
	// The UUID of the target.
	ID *string `json:"id" validate:"required"`

	// The target account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// The trusted profile ID.
	TrustedProfileID *string `json:"trusted_profile_id" validate:"required"`

	// The target name.
	Name *string `json:"name" validate:"required"`

	// List of credentials.
	Credentials []CredentialResponse `json:"credentials,omitempty"`

	// The user ID who created the target.
	CreatedBy *string `json:"created_by,omitempty"`

	// The time when the target was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// The user ID who updated the target.
	UpdatedBy *string `json:"updated_by,omitempty"`

	// The time when the target was updated.
	UpdatedOn *strfmt.DateTime `json:"updated_on,omitempty"`
}

// UnmarshalTarget unmarshals an instance of Target from the specified map of raw messages.
func UnmarshalTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Target)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trusted_profile_id", &obj.TrustedProfileID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "credentials", &obj.Credentials, UnmarshalCredentialResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_by", &obj.UpdatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_on", &obj.UpdatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetCollection : The target list collection.
type TargetCollection struct {
	// The requested page limit.
	Limit *int64 `json:"limit" validate:"required"`

	// The total number of resources that are in the collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// A page reference.
	First *PageHRefFirst `json:"first,omitempty"`

	// A page reference.
	Next *PageHRefNext `json:"next,omitempty"`

	// The details of the target account.
	Targets []Target `json:"targets,omitempty"`
}

// UnmarshalTargetCollection unmarshals an instance of TargetCollection from the specified map of raw messages.
func UnmarshalTargetCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHRefFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHRefNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "targets", &obj.Targets, UnmarshalTarget)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetInfo : The evaluation target.
type TargetInfo struct {
	// The target ID.
	ID *string `json:"id,omitempty"`

	// The target account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The target service name.
	ServiceName *string `json:"service_name,omitempty"`

	// The target service display name.
	ServiceDisplayName *string `json:"service_display_name,omitempty"`

	// The target resource CRN.
	ResourceCRN *string `json:"resource_crn,omitempty"`

	// The target resource name.
	ResourceName *string `json:"resource_name,omitempty"`

	// The collection of different types of tags.
	Tags *Tags `json:"tags,omitempty"`
}

// UnmarshalTargetInfo unmarshals an instance of TargetInfo from the specified map of raw messages.
func UnmarshalTargetInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetInfo)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
	err = core.UnmarshalPrimitive(m, "service_display_name", &obj.ServiceDisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_crn", &obj.ResourceCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_name", &obj.ResourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalTags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TestEvent : The details of a test event response.
type TestEvent struct {
	// The indication of whether the event was received by Event Notifications.
	Success *bool `json:"success" validate:"required"`
}

// UnmarshalTestEvent unmarshals an instance of TestEvent from the specified map of raw messages.
func UnmarshalTestEvent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TestEvent)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateProviderTypeInstanceOptions : The UpdateProviderTypeInstance options.
type UpdateProviderTypeInstanceOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The provider type ID.
	ProviderTypeID *string `json:"provider_type_id" validate:"required,ne="`

	// The provider type instance ID.
	ProviderTypeInstanceID *string `json:"provider_type_instance_id" validate:"required,ne="`

	// The provider type instance name.
	Name *string `json:"name" validate:"required,ne="`

	// The attributes for connecting to the provider type instance.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateProviderTypeInstanceOptions : Instantiate UpdateProviderTypeInstanceOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateProviderTypeInstanceOptions(instanceID string, providerTypeID string, providerTypeInstanceID string, providerTypeInstanceName string) *UpdateProviderTypeInstanceOptions {
	return &UpdateProviderTypeInstanceOptions{
		InstanceID:             core.StringPtr(instanceID),
		Name:                   core.StringPtr(providerTypeInstanceName),
		ProviderTypeID:         core.StringPtr(providerTypeID),
		ProviderTypeInstanceID: core.StringPtr(providerTypeInstanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateProviderTypeInstanceOptions) SetInstanceID(instanceID string) *UpdateProviderTypeInstanceOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProviderTypeID : Allow user to set ProviderTypeID
func (_options *UpdateProviderTypeInstanceOptions) SetProviderTypeID(providerTypeID string) *UpdateProviderTypeInstanceOptions {
	_options.ProviderTypeID = core.StringPtr(providerTypeID)
	return _options
}

// SetProviderTypeInstanceID : Allow user to set ProviderTypeInstanceID
func (_options *UpdateProviderTypeInstanceOptions) SetProviderTypeInstanceID(providerTypeInstanceID string) *UpdateProviderTypeInstanceOptions {
	_options.ProviderTypeInstanceID = core.StringPtr(providerTypeInstanceID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateProviderTypeInstanceOptions) SetName(name string) *UpdateProviderTypeInstanceOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAttributes : Allow user to set Attributes
func (_options *UpdateProviderTypeInstanceOptions) SetAttributes(attributes map[string]interface{}) *UpdateProviderTypeInstanceOptions {
	_options.Attributes = attributes
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateProviderTypeInstanceOptions) SetHeaders(param map[string]string) *UpdateProviderTypeInstanceOptions {
	options.Headers = param
	return options
}

// UpdateScopeOptions : The UpdateScope options.
type UpdateScopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The scope name.
	Name *string `json:"name,omitempty"`

	// The scope description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateScopeOptions : Instantiate UpdateScopeOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateScopeOptions(instanceID string, scopeID string) *UpdateScopeOptions {
	return &UpdateScopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateScopeOptions) SetInstanceID(instanceID string) *UpdateScopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *UpdateScopeOptions) SetScopeID(scopeID string) *UpdateScopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateScopeOptions) SetName(name string) *UpdateScopeOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateScopeOptions) SetDescription(description string) *UpdateScopeOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateScopeOptions) SetHeaders(param map[string]string) *UpdateScopeOptions {
	options.Headers = param
	return options
}

// UpdateSettingsOptions : The UpdateSettings options.
type UpdateSettingsOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The payload to connect a Cloud Object Storage instance to an Security and Compliance Center instance.
	ObjectStorage *ObjectStoragePrototype `json:"object_storage,omitempty"`

	// The payload to connect an Event Notification instance with a Security and Compliance Center instance.
	EventNotifications *EventNotificationsPrototype `json:"event_notifications,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSettingsOptions : Instantiate UpdateSettingsOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateSettingsOptions(instanceID string) *UpdateSettingsOptions {
	return &UpdateSettingsOptions{
		InstanceID: core.StringPtr(instanceID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateSettingsOptions) SetInstanceID(instanceID string) *UpdateSettingsOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetObjectStorage : Allow user to set ObjectStorage
func (_options *UpdateSettingsOptions) SetObjectStorage(objectStorage *ObjectStoragePrototype) *UpdateSettingsOptions {
	_options.ObjectStorage = objectStorage
	return _options
}

// SetEventNotifications : Allow user to set EventNotifications
func (_options *UpdateSettingsOptions) SetEventNotifications(eventNotifications *EventNotificationsPrototype) *UpdateSettingsOptions {
	_options.EventNotifications = eventNotifications
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSettingsOptions) SetHeaders(param map[string]string) *UpdateSettingsOptions {
	options.Headers = param
	return options
}

// UpdateSubscopeOptions : The UpdateSubscope options.
type UpdateSubscopeOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	ScopeID *string `json:"scope_id" validate:"required,ne="`

	// The ID of the scope being targeted.
	SubscopeID *string `json:"subscope_id" validate:"required,ne="`

	// The scope name.
	Name *string `json:"name,omitempty"`

	// The scope description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSubscopeOptions : Instantiate UpdateSubscopeOptions
func (*SecurityAndComplianceCenterApiV3) NewUpdateSubscopeOptions(instanceID string, scopeID string, subscopeID string) *UpdateSubscopeOptions {
	return &UpdateSubscopeOptions{
		InstanceID: core.StringPtr(instanceID),
		ScopeID:    core.StringPtr(scopeID),
		SubscopeID: core.StringPtr(subscopeID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpdateSubscopeOptions) SetInstanceID(instanceID string) *UpdateSubscopeOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetScopeID : Allow user to set ScopeID
func (_options *UpdateSubscopeOptions) SetScopeID(scopeID string) *UpdateSubscopeOptions {
	_options.ScopeID = core.StringPtr(scopeID)
	return _options
}

// SetSubscopeID : Allow user to set SubscopeID
func (_options *UpdateSubscopeOptions) SetSubscopeID(subscopeID string) *UpdateSubscopeOptions {
	_options.SubscopeID = core.StringPtr(subscopeID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateSubscopeOptions) SetName(name string) *UpdateSubscopeOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateSubscopeOptions) SetDescription(description string) *UpdateSubscopeOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSubscopeOptions) SetHeaders(param map[string]string) *UpdateSubscopeOptions {
	options.Headers = param
	return options
}

// UpgradeAttachmentOptions : The UpgradeAttachment options.
type UpgradeAttachmentOptions struct {
	// The ID of the Security and Compliance Center instance.
	InstanceID *string `json:"instance_id" validate:"required,ne="`

	// The profile ID.
	ProfileID *string `json:"profile_id" validate:"required,ne="`

	// The attachment ID.
	AttachmentID *string `json:"attachment_id" validate:"required,ne="`

	// The attachment_parameters to set for a Profile Attachment.
	AttachmentParameters []Parameter `json:"attachment_parameters,omitempty"`

	// The user account ID.
	AccountID *string `json:"account_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpgradeAttachmentOptions : Instantiate UpgradeAttachmentOptions
func (*SecurityAndComplianceCenterApiV3) NewUpgradeAttachmentOptions(instanceID string, profileID string, attachmentID string) *UpgradeAttachmentOptions {
	return &UpgradeAttachmentOptions{
		InstanceID:   core.StringPtr(instanceID),
		ProfileID:    core.StringPtr(profileID),
		AttachmentID: core.StringPtr(attachmentID),
	}
}

// SetInstanceID : Allow user to set InstanceID
func (_options *UpgradeAttachmentOptions) SetInstanceID(instanceID string) *UpgradeAttachmentOptions {
	_options.InstanceID = core.StringPtr(instanceID)
	return _options
}

// SetProfileID : Allow user to set ProfileID
func (_options *UpgradeAttachmentOptions) SetProfileID(profileID string) *UpgradeAttachmentOptions {
	_options.ProfileID = core.StringPtr(profileID)
	return _options
}

// SetAttachmentID : Allow user to set AttachmentID
func (_options *UpgradeAttachmentOptions) SetAttachmentID(attachmentID string) *UpgradeAttachmentOptions {
	_options.AttachmentID = core.StringPtr(attachmentID)
	return _options
}

// SetAttachmentParameters : Allow user to set AttachmentParameters
func (_options *UpgradeAttachmentOptions) SetAttachmentParameters(attachmentParameters []Parameter) *UpgradeAttachmentOptions {
	_options.AttachmentParameters = attachmentParameters
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *UpgradeAttachmentOptions) SetAccountID(accountID string) *UpgradeAttachmentOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpgradeAttachmentOptions) SetHeaders(param map[string]string) *UpgradeAttachmentOptions {
	options.Headers = param
	return options
}

// ConditionItemConditionBase : The required configuration base object.
// This model "extends" ConditionItem
type ConditionItemConditionBase struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property" validate:"required"`

	// The operator.
	Operator *string `json:"operator" validate:"required"`

	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the ConditionItemConditionBase.Operator property.
// The operator.
const (
	ConditionItemConditionBaseOperatorDaysLessThanConst         = "days_less_than"
	ConditionItemConditionBaseOperatorIpsEqualsConst            = "ips_equals"
	ConditionItemConditionBaseOperatorIpsInRangeConst           = "ips_in_range"
	ConditionItemConditionBaseOperatorIpsNotEqualsConst         = "ips_not_equals"
	ConditionItemConditionBaseOperatorIsEmptyConst              = "is_empty"
	ConditionItemConditionBaseOperatorIsFalseConst              = "is_false"
	ConditionItemConditionBaseOperatorIsNotEmptyConst           = "is_not_empty"
	ConditionItemConditionBaseOperatorIsTrueConst               = "is_true"
	ConditionItemConditionBaseOperatorNumEqualsConst            = "num_equals"
	ConditionItemConditionBaseOperatorNumGreaterThanConst       = "num_greater_than"
	ConditionItemConditionBaseOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	ConditionItemConditionBaseOperatorNumLessThanConst          = "num_less_than"
	ConditionItemConditionBaseOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	ConditionItemConditionBaseOperatorNumNotEqualsConst         = "num_not_equals"
	ConditionItemConditionBaseOperatorStringContainsConst       = "string_contains"
	ConditionItemConditionBaseOperatorStringEqualsConst         = "string_equals"
	ConditionItemConditionBaseOperatorStringMatchConst          = "string_match"
	ConditionItemConditionBaseOperatorStringNotContainsConst    = "string_not_contains"
	ConditionItemConditionBaseOperatorStringNotEqualsConst      = "string_not_equals"
	ConditionItemConditionBaseOperatorStringNotMatchConst       = "string_not_match"
	ConditionItemConditionBaseOperatorStringsAllowedConst       = "strings_allowed"
	ConditionItemConditionBaseOperatorStringsInListConst        = "strings_in_list"
	ConditionItemConditionBaseOperatorStringsRequiredConst      = "strings_required"
)

// NewConditionItemConditionBase : Instantiate ConditionItemConditionBase (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewConditionItemConditionBase(property string, operator string) (_model *ConditionItemConditionBase, err error) {
	_model = &ConditionItemConditionBase{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*ConditionItemConditionBase) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionBase unmarshals an instance of ConditionItemConditionBase from the specified map of raw messages.
func UnmarshalConditionItemConditionBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionBase)
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

// ConditionItemConditionList : A list of required configurations.
// Models which "extend" this model:
// - ConditionItemConditionListConditionListConditionOr
// - ConditionItemConditionListConditionListConditionAnd
// This model "extends" ConditionItem
type ConditionItemConditionList struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []ConditionItemIntf `json:"or,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []ConditionItemIntf `json:"and,omitempty"`
}

func (*ConditionItemConditionList) isaConditionItemConditionList() bool {
	return true
}

type ConditionItemConditionListIntf interface {
	ConditionItemIntf
	isaConditionItemConditionList() bool
}

func (*ConditionItemConditionList) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionList unmarshals an instance of ConditionItemConditionList from the specified map of raw messages.
func UnmarshalConditionItemConditionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionList)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalConditionItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionSubRule : ConditionItemConditionSubRule struct
// Models which "extend" this model:
// - ConditionItemConditionSubRuleConditionSubRuleConditionAny
// - ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf
// - ConditionItemConditionSubRuleConditionSubRuleConditionAll
// - ConditionItemConditionSubRuleConditionSubRuleConditionAllIf
// This model "extends" ConditionItem
type ConditionItemConditionSubRule struct {
	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`

	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

func (*ConditionItemConditionSubRule) isaConditionItemConditionSubRule() bool {
	return true
}

type ConditionItemConditionSubRuleIntf interface {
	ConditionItemIntf
	isaConditionItemConditionSubRule() bool
}

func (*ConditionItemConditionSubRule) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionSubRule unmarshals an instance of ConditionItemConditionSubRule from the specified map of raw messages.
func UnmarshalConditionItemConditionSubRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionSubRule)
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionBase : The required configuration base object.
// This model "extends" RequiredConfig
type RequiredConfigConditionBase struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// The property.
	Property *string `json:"property" validate:"required"`

	// The operator.
	Operator *string `json:"operator" validate:"required"`

	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the RequiredConfigConditionBase.Operator property.
// The operator.
const (
	RequiredConfigConditionBaseOperatorDaysLessThanConst         = "days_less_than"
	RequiredConfigConditionBaseOperatorIpsEqualsConst            = "ips_equals"
	RequiredConfigConditionBaseOperatorIpsInRangeConst           = "ips_in_range"
	RequiredConfigConditionBaseOperatorIpsNotEqualsConst         = "ips_not_equals"
	RequiredConfigConditionBaseOperatorIsEmptyConst              = "is_empty"
	RequiredConfigConditionBaseOperatorIsFalseConst              = "is_false"
	RequiredConfigConditionBaseOperatorIsNotEmptyConst           = "is_not_empty"
	RequiredConfigConditionBaseOperatorIsTrueConst               = "is_true"
	RequiredConfigConditionBaseOperatorNumEqualsConst            = "num_equals"
	RequiredConfigConditionBaseOperatorNumGreaterThanConst       = "num_greater_than"
	RequiredConfigConditionBaseOperatorNumGreaterThanEqualsConst = "num_greater_than_equals"
	RequiredConfigConditionBaseOperatorNumLessThanConst          = "num_less_than"
	RequiredConfigConditionBaseOperatorNumLessThanEqualsConst    = "num_less_than_equals"
	RequiredConfigConditionBaseOperatorNumNotEqualsConst         = "num_not_equals"
	RequiredConfigConditionBaseOperatorStringContainsConst       = "string_contains"
	RequiredConfigConditionBaseOperatorStringEqualsConst         = "string_equals"
	RequiredConfigConditionBaseOperatorStringMatchConst          = "string_match"
	RequiredConfigConditionBaseOperatorStringNotContainsConst    = "string_not_contains"
	RequiredConfigConditionBaseOperatorStringNotEqualsConst      = "string_not_equals"
	RequiredConfigConditionBaseOperatorStringNotMatchConst       = "string_not_match"
	RequiredConfigConditionBaseOperatorStringsAllowedConst       = "strings_allowed"
	RequiredConfigConditionBaseOperatorStringsInListConst        = "strings_in_list"
	RequiredConfigConditionBaseOperatorStringsRequiredConst      = "strings_required"
)

// NewRequiredConfigConditionBase : Instantiate RequiredConfigConditionBase (Generic Model Constructor)
func (*SecurityAndComplianceCenterApiV3) NewRequiredConfigConditionBase(property string, operator string) (_model *RequiredConfigConditionBase, err error) {
	_model = &RequiredConfigConditionBase{
		Property: core.StringPtr(property),
		Operator: core.StringPtr(operator),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*RequiredConfigConditionBase) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionBase unmarshals an instance of RequiredConfigConditionBase from the specified map of raw messages.
func UnmarshalRequiredConfigConditionBase(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionBase)
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

// RequiredConfigConditionList : A list of required configurations.
// Models which "extend" this model:
// - RequiredConfigConditionListConditionListConditionOr
// - RequiredConfigConditionListConditionListConditionAnd
// This model "extends" RequiredConfig
type RequiredConfigConditionList struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []ConditionItemIntf `json:"or,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []ConditionItemIntf `json:"and,omitempty"`
}

func (*RequiredConfigConditionList) isaRequiredConfigConditionList() bool {
	return true
}

type RequiredConfigConditionListIntf interface {
	RequiredConfigIntf
	isaRequiredConfigConditionList() bool
}

func (*RequiredConfigConditionList) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionList unmarshals an instance of RequiredConfigConditionList from the specified map of raw messages.
func UnmarshalRequiredConfigConditionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionList)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalConditionItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionSubRule : RequiredConfigConditionSubRule struct
// Models which "extend" this model:
// - RequiredConfigConditionSubRuleConditionSubRuleConditionAny
// - RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf
// - RequiredConfigConditionSubRuleConditionSubRuleConditionAll
// - RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf
// This model "extends" RequiredConfig
type RequiredConfigConditionSubRule struct {
	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`

	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`

	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

func (*RequiredConfigConditionSubRule) isaRequiredConfigConditionSubRule() bool {
	return true
}

type RequiredConfigConditionSubRuleIntf interface {
	RequiredConfigIntf
	isaRequiredConfigConditionSubRule() bool
}

func (*RequiredConfigConditionSubRule) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionSubRule unmarshals an instance of RequiredConfigConditionSubRule from the specified map of raw messages.
func UnmarshalRequiredConfigConditionSubRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionSubRule)
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopePropertyExclusions : Any exclusions or resources that should not be part of the scope. Has to be the same type as the one specified.
// This model "extends" ScopeProperty
type ScopePropertyExclusions struct {
	Name *string `json:"name,omitempty"`

	// A list of scopes/targets to exclude from a scope.
	Value []ScopePropertyExclusionItem `json:"value,omitempty"`
}

func (*ScopePropertyExclusions) isaScopeProperty() bool {
	return true
}

// UnmarshalScopePropertyExclusions unmarshals an instance of ScopePropertyExclusions from the specified map of raw messages.
func UnmarshalScopePropertyExclusions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopePropertyExclusions)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value", &obj.Value, UnmarshalScopePropertyExclusionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ScopePropertyScopeID : The value of the identifier that correlates to scope type. If ScopePropertyScopeType schema uses the value
// enterprise.account_group, the value should be the identifier or ID of the account_group within the enterprise.
// This model "extends" ScopeProperty
type ScopePropertyScopeID struct {
	// The key for the scope property.
	Name *string `json:"name,omitempty"`

	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the ScopePropertyScopeID.Name property.
// The key for the scope property.
const (
	ScopePropertyScopeIDNameScopeIDConst = "scope_id"
)

func (*ScopePropertyScopeID) isaScopeProperty() bool {
	return true
}

// UnmarshalScopePropertyScopeID unmarshals an instance of ScopePropertyScopeID from the specified map of raw messages.
func UnmarshalScopePropertyScopeID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopePropertyScopeID)
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

// ScopePropertyScopeType : Attribute that details what kind of type of scope.
// This model "extends" ScopeProperty
type ScopePropertyScopeType struct {
	// key to say the attribute targets the scope type.
	Name *string `json:"name,omitempty"`

	// The type of scope it targets
	//
	// The scope values are as followed:
	// - enterprise: The scope targets an enterprise account
	// - enterprise.account_group: The scope targets an account group within an enterprise
	// - enterprise.account: The scope targets an account within an enterprise
	// - account: The scope targets an account not tied to an enterprise
	// - account.resource_group: The scope targets a resource group within an account.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the ScopePropertyScopeType.Name property.
// key to say the attribute targets the scope type.
const (
	ScopePropertyScopeTypeNameScopeTypeConst = "scope_type"
)

// Constants associated with the ScopePropertyScopeType.Value property.
// The type of scope it targets
//
// The scope values are as followed:
// - enterprise: The scope targets an enterprise account
// - enterprise.account_group: The scope targets an account group within an enterprise
// - enterprise.account: The scope targets an account within an enterprise
// - account: The scope targets an account not tied to an enterprise
// - account.resource_group: The scope targets a resource group within an account.
const (
	ScopePropertyScopeTypeValueAccountConst                = "account"
	ScopePropertyScopeTypeValueAccountResourceGroupConst   = "account.resource_group"
	ScopePropertyScopeTypeValueEnterpriseConst             = "enterprise"
	ScopePropertyScopeTypeValueEnterpriseAccountConst      = "enterprise.account"
	ScopePropertyScopeTypeValueEnterpriseAccountGroupConst = "enterprise.account_group"
)

func (*ScopePropertyScopeType) isaScopeProperty() bool {
	return true
}

// UnmarshalScopePropertyScopeType unmarshals an instance of ScopePropertyScopeType from the specified map of raw messages.
func UnmarshalScopePropertyScopeType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ScopePropertyScopeType)
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

// ConditionItemConditionListConditionListConditionAnd : ConditionItemConditionListConditionListConditionAnd struct
// This model "extends" ConditionItemConditionList
type ConditionItemConditionListConditionListConditionAnd struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []ConditionItemIntf `json:"and,omitempty"`
}

func (*ConditionItemConditionListConditionListConditionAnd) isaConditionItemConditionList() bool {
	return true
}

func (*ConditionItemConditionListConditionListConditionAnd) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionListConditionListConditionAnd unmarshals an instance of ConditionItemConditionListConditionListConditionAnd from the specified map of raw messages.
func UnmarshalConditionItemConditionListConditionListConditionAnd(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionListConditionListConditionAnd)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionListConditionListConditionOr : The `OR` required configurations.
// This model "extends" ConditionItemConditionList
type ConditionItemConditionListConditionListConditionOr struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []ConditionItemIntf `json:"or,omitempty"`
}

func (*ConditionItemConditionListConditionListConditionOr) isaConditionItemConditionList() bool {
	return true
}

func (*ConditionItemConditionListConditionListConditionOr) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionListConditionListConditionOr unmarshals an instance of ConditionItemConditionListConditionListConditionOr from the specified map of raw messages.
func UnmarshalConditionItemConditionListConditionListConditionOr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionListConditionListConditionOr)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionSubRuleConditionSubRuleConditionAll : A subrule using the 'all' operator.
// This model "extends" ConditionItemConditionSubRule
type ConditionItemConditionSubRuleConditionSubRuleConditionAll struct {
	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAll) isaConditionItemConditionSubRule() bool {
	return true
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAll) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAll unmarshals an instance of ConditionItemConditionSubRuleConditionSubRuleConditionAll from the specified map of raw messages.
func UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAll(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionSubRuleConditionSubRuleConditionAll)
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionSubRuleConditionSubRuleConditionAllIf : A subrule using the 'all_ifexists' operator.
// This model "extends" ConditionItemConditionSubRule
type ConditionItemConditionSubRuleConditionSubRuleConditionAllIf struct {
	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAllIf) isaConditionItemConditionSubRule() bool {
	return true
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAllIf) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAllIf unmarshals an instance of ConditionItemConditionSubRuleConditionSubRuleConditionAllIf from the specified map of raw messages.
func UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAllIf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionSubRuleConditionSubRuleConditionAllIf)
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionSubRuleConditionSubRuleConditionAny : A subrule using the 'any' operator.
// This model "extends" ConditionItemConditionSubRule
type ConditionItemConditionSubRuleConditionSubRuleConditionAny struct {
	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAny) isaConditionItemConditionSubRule() bool {
	return true
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAny) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAny unmarshals an instance of ConditionItemConditionSubRuleConditionSubRuleConditionAny from the specified map of raw messages.
func UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAny(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionSubRuleConditionSubRuleConditionAny)
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf : A subrule using the 'any_ifexists' operator.
// This model "extends" ConditionItemConditionSubRule
type ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf struct {
	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf) isaConditionItemConditionSubRule() bool {
	return true
}

func (*ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf) isaConditionItem() bool {
	return true
}

// UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAnyIf unmarshals an instance of ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf from the specified map of raw messages.
func UnmarshalConditionItemConditionSubRuleConditionSubRuleConditionAnyIf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConditionItemConditionSubRuleConditionSubRuleConditionAnyIf)
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionListConditionListConditionAnd : RequiredConfigConditionListConditionListConditionAnd struct
// This model "extends" RequiredConfigConditionList
type RequiredConfigConditionListConditionListConditionAnd struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where all items should evaluate to true.
	And []ConditionItemIntf `json:"and,omitempty"`
}

func (*RequiredConfigConditionListConditionListConditionAnd) isaRequiredConfigConditionList() bool {
	return true
}

func (*RequiredConfigConditionListConditionListConditionAnd) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionListConditionListConditionAnd unmarshals an instance of RequiredConfigConditionListConditionListConditionAnd from the specified map of raw messages.
func UnmarshalRequiredConfigConditionListConditionListConditionAnd(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionListConditionListConditionAnd)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "and", &obj.And, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionListConditionListConditionOr : The `OR` required configurations.
// This model "extends" RequiredConfigConditionList
type RequiredConfigConditionListConditionListConditionOr struct {
	// The required config description.
	Description *string `json:"description,omitempty"`

	// A list of required configurations where one item should evaluate to true.
	Or []ConditionItemIntf `json:"or,omitempty"`
}

func (*RequiredConfigConditionListConditionListConditionOr) isaRequiredConfigConditionList() bool {
	return true
}

func (*RequiredConfigConditionListConditionListConditionOr) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionListConditionListConditionOr unmarshals an instance of RequiredConfigConditionListConditionListConditionOr from the specified map of raw messages.
func UnmarshalRequiredConfigConditionListConditionListConditionOr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionListConditionListConditionOr)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "or", &obj.Or, UnmarshalConditionItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionSubRuleConditionSubRuleConditionAll : A subrule using the 'all' operator.
// This model "extends" RequiredConfigConditionSubRule
type RequiredConfigConditionSubRuleConditionSubRuleConditionAll struct {
	// A rule within a rule used in the requiredConfig.
	All *SubRule `json:"all,omitempty"`
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAll) isaRequiredConfigConditionSubRule() bool {
	return true
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAll) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAll unmarshals an instance of RequiredConfigConditionSubRuleConditionSubRuleConditionAll from the specified map of raw messages.
func UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAll(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionSubRuleConditionSubRuleConditionAll)
	err = core.UnmarshalModel(m, "all", &obj.All, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf : A subrule using the 'all_ifexists' operator.
// This model "extends" RequiredConfigConditionSubRule
type RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf struct {
	// A rule within a rule used in the requiredConfig.
	AllIfexists *SubRule `json:"all_ifexists,omitempty"`
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf) isaRequiredConfigConditionSubRule() bool {
	return true
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAllIf unmarshals an instance of RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf from the specified map of raw messages.
func UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAllIf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionSubRuleConditionSubRuleConditionAllIf)
	err = core.UnmarshalModel(m, "all_ifexists", &obj.AllIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionSubRuleConditionSubRuleConditionAny : A subrule using the 'any' operator.
// This model "extends" RequiredConfigConditionSubRule
type RequiredConfigConditionSubRuleConditionSubRuleConditionAny struct {
	// A rule within a rule used in the requiredConfig.
	Any *SubRule `json:"any,omitempty"`
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAny) isaRequiredConfigConditionSubRule() bool {
	return true
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAny) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAny unmarshals an instance of RequiredConfigConditionSubRuleConditionSubRuleConditionAny from the specified map of raw messages.
func UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAny(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionSubRuleConditionSubRuleConditionAny)
	err = core.UnmarshalModel(m, "any", &obj.Any, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf : A subrule using the 'any_ifexists' operator.
// This model "extends" RequiredConfigConditionSubRule
type RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf struct {
	// A rule within a rule used in the requiredConfig.
	AnyIfexists *SubRule `json:"any_ifexists,omitempty"`
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf) isaRequiredConfigConditionSubRule() bool {
	return true
}

func (*RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf) isaRequiredConfig() bool {
	return true
}

// UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf unmarshals an instance of RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf from the specified map of raw messages.
func UnmarshalRequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RequiredConfigConditionSubRuleConditionSubRuleConditionAnyIf)
	err = core.UnmarshalModel(m, "any_ifexists", &obj.AnyIfexists, UnmarshalSubRule)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ControlLibrariesPager can be used to simplify the use of the "ListControlLibraries" method.
type ControlLibrariesPager struct {
	hasNext     bool
	options     *ListControlLibrariesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewControlLibrariesPager returns a new ControlLibrariesPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewControlLibrariesPager(options *ListControlLibrariesOptions) (pager *ControlLibrariesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListControlLibrariesOptions = *options
	pager = &ControlLibrariesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ControlLibrariesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ControlLibrariesPager) GetNextWithContext(ctx context.Context) (page []ControlLibrary, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListControlLibrariesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.ControlLibraries

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ControlLibrariesPager) GetAllWithContext(ctx context.Context) (allItems []ControlLibrary, err error) {
	for pager.HasNext() {
		var nextPage []ControlLibrary
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ControlLibrariesPager) GetNext() (page []ControlLibrary, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ControlLibrariesPager) GetAll() (allItems []ControlLibrary, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ProfilesPager can be used to simplify the use of the "ListProfiles" method.
type ProfilesPager struct {
	hasNext     bool
	options     *ListProfilesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewControlLibrariesPager returns a new ControlLibrariesPager instance.
func (securityAndComplianceCenterApi *SecurityAndComplianceCenterApiV3) NewProfilesPager(options *ListProfilesOptions) (pager *ProfilesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListProfilesOptions = *options
	pager = &ProfilesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenterApi,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ProfilesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ProfilesPager) GetNextWithContext(ctx context.Context) (page []Profile, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListProfilesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Profiles

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProfilesPager) GetAllWithContext(ctx context.Context) (allItems []Profile, err error) {
	allItems = make([]Profile, 0)
	for pager.HasNext() {
		var nextPage []Profile
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ProfilesPager) GetNext() (page []Profile, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProfilesPager) GetAll() (allItems []Profile, err error) {
	return pager.GetAllWithContext(context.Background())
}

// InstanceAttachmentsPager can be used to simplify the use of the "ListInstanceAttachments" method.
type InstanceAttachmentsPager struct {
	hasNext     bool
	options     *ListInstanceAttachmentsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewInstanceAttachmentsPager returns a new InstanceAttachmentsPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewInstanceAttachmentsPager(options *ListInstanceAttachmentsOptions) (pager *InstanceAttachmentsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListInstanceAttachmentsOptions = *options
	pager = &InstanceAttachmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *InstanceAttachmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *InstanceAttachmentsPager) GetNextWithContext(ctx context.Context) (page []ProfileAttachment, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListInstanceAttachmentsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Attachments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *InstanceAttachmentsPager) GetAllWithContext(ctx context.Context) (allItems []ProfileAttachment, err error) {
	for pager.HasNext() {
		var nextPage []ProfileAttachment
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *InstanceAttachmentsPager) GetNext() (page []ProfileAttachment, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *InstanceAttachmentsPager) GetAll() (allItems []ProfileAttachment, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ScopesPager can be used to simplify the use of the "ListScopes" method.
type ScopesPager struct {
	hasNext     bool
	options     *ListScopesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewScopesPager returns a new ScopesPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewScopesPager(options *ListScopesOptions) (pager *ScopesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListScopesOptions = *options
	pager = &ScopesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ScopesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ScopesPager) GetNextWithContext(ctx context.Context) (page []Scope, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListScopesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil && *pager.pageContext.next != "")
	page = result.Scopes

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ScopesPager) GetAllWithContext(ctx context.Context) (allItems []Scope, err error) {
	allItems = make([]Scope, 0)
	for pager.HasNext() {
		var nextPage []Scope
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ScopesPager) GetNext() (page []Scope, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ScopesPager) GetAll() (allItems []Scope, err error) {
	return pager.GetAllWithContext(context.Background())
}

// SubscopesPager can be used to simplify the use of the "ListSubscopes" method.
type SubscopesPager struct {
	hasNext     bool
	options     *ListSubscopesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewSubscopesPager returns a new SubscopesPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewSubscopesPager(options *ListSubscopesOptions) (pager *SubscopesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListSubscopesOptions = *options
	pager = &SubscopesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SubscopesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SubscopesPager) GetNextWithContext(ctx context.Context) (page []SubScope, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListSubscopesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Subscopes

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SubscopesPager) GetAllWithContext(ctx context.Context) (allItems []SubScope, err error) {
	for pager.HasNext() {
		var nextPage []SubScope
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *SubscopesPager) GetNext() (page []SubScope, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SubscopesPager) GetAll() (allItems []SubScope, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportsPager can be used to simplify the use of the "ListReports" method.
type ReportsPager struct {
	hasNext     bool
	options     *ListReportsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportsPager returns a new ReportsPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewReportsPager(options *ListReportsOptions) (pager *ReportsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportsOptions = *options
	pager = &ReportsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportsPager) GetNextWithContext(ctx context.Context) (page []Report, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Reports

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportsPager) GetAllWithContext(ctx context.Context) (allItems []Report, err error) {
	for pager.HasNext() {
		var nextPage []Report
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportsPager) GetNext() (page []Report, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportsPager) GetAll() (allItems []Report, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportEvaluationsPager can be used to simplify the use of the "ListReportEvaluations" method.
type ReportEvaluationsPager struct {
	hasNext     bool
	options     *ListReportEvaluationsOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportEvaluationsPager returns a new ReportEvaluationsPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewReportEvaluationsPager(options *ListReportEvaluationsOptions) (pager *ReportEvaluationsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportEvaluationsOptions = *options
	pager = &ReportEvaluationsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportEvaluationsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportEvaluationsPager) GetNextWithContext(ctx context.Context) (page []Evaluation, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportEvaluationsWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Evaluations

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportEvaluationsPager) GetAllWithContext(ctx context.Context) (allItems []Evaluation, err error) {
	for pager.HasNext() {
		var nextPage []Evaluation
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportEvaluationsPager) GetNext() (page []Evaluation, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportEvaluationsPager) GetAll() (allItems []Evaluation, err error) {
	return pager.GetAllWithContext(context.Background())
}

// ReportResourcesPager can be used to simplify the use of the "ListReportResources" method.
type ReportResourcesPager struct {
	hasNext     bool
	options     *ListReportResourcesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewReportResourcesPager returns a new ReportResourcesPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewReportResourcesPager(options *ListReportResourcesOptions) (pager *ReportResourcesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListReportResourcesOptions = *options
	pager = &ReportResourcesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ReportResourcesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ReportResourcesPager) GetNextWithContext(ctx context.Context) (page []Resource, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListReportResourcesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Resources

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ReportResourcesPager) GetAllWithContext(ctx context.Context) (allItems []Resource, err error) {
	for pager.HasNext() {
		var nextPage []Resource
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *ReportResourcesPager) GetNext() (page []Resource, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ReportResourcesPager) GetAll() (allItems []Resource, err error) {
	return pager.GetAllWithContext(context.Background())
}

// RulesPager can be used to simplify the use of the "ListRules" method.
type RulesPager struct {
	hasNext     bool
	options     *ListRulesOptions
	client      *SecurityAndComplianceCenterApiV3
	pageContext struct {
		next *string
	}
}

// NewRulesPager returns a new RulesPager instance.
func (securityAndComplianceCenter *SecurityAndComplianceCenterApiV3) NewRulesPager(options *ListRulesOptions) (pager *RulesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = fmt.Errorf("the 'options.Start' field should not be set")
		return
	}

	var optionsCopy ListRulesOptions = *options
	pager = &RulesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  securityAndComplianceCenter,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *RulesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *RulesPager) GetNextWithContext(ctx context.Context) (page []Rule, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListRulesWithContext(ctx, pager.options)
	if err != nil {
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Rules

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *RulesPager) GetAllWithContext(ctx context.Context) (allItems []Rule, err error) {
	for pager.HasNext() {
		var nextPage []Rule
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *RulesPager) GetNext() (page []Rule, err error) {
	return pager.GetNextWithContext(context.Background())
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *RulesPager) GetAll() (allItems []Rule, err error) {
	return pager.GetAllWithContext(context.Background())
}
