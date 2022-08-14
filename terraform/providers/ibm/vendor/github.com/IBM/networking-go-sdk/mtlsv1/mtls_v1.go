/**
 * (C) Copyright IBM Corp. 2022.
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
 * IBM OpenAPI SDK Code Generator Version: 3.48.0-e80b60a1-20220414-145125
 */

// Package mtlsv1 : Operations and models for the MtlsV1 service
package mtlsv1

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

// MtlsV1 : MTLS
//
// API Version: 1.0.0
type MtlsV1 struct {
	Service *core.BaseService

	// Cloud resource name.
	Crn *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "mtls"

// MtlsV1Options : Service options
type MtlsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Cloud resource name.
	Crn *string `validate:"required"`
}

// NewMtlsV1UsingExternalConfig : constructs an instance of MtlsV1 with passed in options and external configuration.
func NewMtlsV1UsingExternalConfig(options *MtlsV1Options) (mtls *MtlsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	mtls, err = NewMtlsV1(options)
	if err != nil {
		return
	}

	err = mtls.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = mtls.Service.SetServiceURL(options.URL)
	}
	return
}

// NewMtlsV1 : constructs an instance of MtlsV1 with passed in options.
func NewMtlsV1(options *MtlsV1Options) (service *MtlsV1, err error) {
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

	service = &MtlsV1{
		Service: baseService,
		Crn: options.Crn,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "mtls" suitable for processing requests.
func (mtls *MtlsV1) Clone() *MtlsV1 {
	if core.IsNil(mtls) {
		return nil
	}
	clone := *mtls
	clone.Service = mtls.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (mtls *MtlsV1) SetServiceURL(url string) error {
	return mtls.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (mtls *MtlsV1) GetServiceURL() string {
	return mtls.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (mtls *MtlsV1) SetDefaultHeaders(headers http.Header) {
	mtls.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (mtls *MtlsV1) SetEnableGzipCompression(enableGzip bool) {
	mtls.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (mtls *MtlsV1) GetEnableGzipCompression() bool {
	return mtls.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (mtls *MtlsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	mtls.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (mtls *MtlsV1) DisableRetries() {
	mtls.Service.DisableRetries()
}

// ListAccessCertificates : List access certificates
// List access certificates.
func (mtls *MtlsV1) ListAccessCertificates(listAccessCertificatesOptions *ListAccessCertificatesOptions) (result *ListAccessCertsResp, response *core.DetailedResponse, err error) {
	return mtls.ListAccessCertificatesWithContext(context.Background(), listAccessCertificatesOptions)
}

// ListAccessCertificatesWithContext is an alternate form of the ListAccessCertificates method which supports a Context parameter
func (mtls *MtlsV1) ListAccessCertificatesWithContext(ctx context.Context, listAccessCertificatesOptions *ListAccessCertificatesOptions) (result *ListAccessCertsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAccessCertificatesOptions, "listAccessCertificatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAccessCertificatesOptions, "listAccessCertificatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *listAccessCertificatesOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAccessCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "ListAccessCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAccessCertsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAccessCertificate : Create access certificate
// Create access certificate.
func (mtls *MtlsV1) CreateAccessCertificate(createAccessCertificateOptions *CreateAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	return mtls.CreateAccessCertificateWithContext(context.Background(), createAccessCertificateOptions)
}

// CreateAccessCertificateWithContext is an alternate form of the CreateAccessCertificate method which supports a Context parameter
func (mtls *MtlsV1) CreateAccessCertificateWithContext(ctx context.Context, createAccessCertificateOptions *CreateAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccessCertificateOptions, "createAccessCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createAccessCertificateOptions, "createAccessCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *createAccessCertificateOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAccessCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "CreateAccessCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccessCertificateOptions.Name != nil {
		body["name"] = createAccessCertificateOptions.Name
	}
	if createAccessCertificateOptions.Certificate != nil {
		body["certificate"] = createAccessCertificateOptions.Certificate
	}
	if createAccessCertificateOptions.AssociatedHostnames != nil {
		body["associated_hostnames"] = createAccessCertificateOptions.AssociatedHostnames
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAccessCertificate : Get access certificate
// Get access certificate.
func (mtls *MtlsV1) GetAccessCertificate(getAccessCertificateOptions *GetAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	return mtls.GetAccessCertificateWithContext(context.Background(), getAccessCertificateOptions)
}

// GetAccessCertificateWithContext is an alternate form of the GetAccessCertificate method which supports a Context parameter
func (mtls *MtlsV1) GetAccessCertificateWithContext(ctx context.Context, getAccessCertificateOptions *GetAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccessCertificateOptions, "getAccessCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAccessCertificateOptions, "getAccessCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *getAccessCertificateOptions.ZoneID,
		"cert_id": *getAccessCertificateOptions.CertID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates/{cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAccessCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "GetAccessCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccessCertificate : Update access certificate
// Update access certificate.
func (mtls *MtlsV1) UpdateAccessCertificate(updateAccessCertificateOptions *UpdateAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	return mtls.UpdateAccessCertificateWithContext(context.Background(), updateAccessCertificateOptions)
}

// UpdateAccessCertificateWithContext is an alternate form of the UpdateAccessCertificate method which supports a Context parameter
func (mtls *MtlsV1) UpdateAccessCertificateWithContext(ctx context.Context, updateAccessCertificateOptions *UpdateAccessCertificateOptions) (result *AccessCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccessCertificateOptions, "updateAccessCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAccessCertificateOptions, "updateAccessCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *updateAccessCertificateOptions.ZoneID,
		"cert_id": *updateAccessCertificateOptions.CertID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates/{cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAccessCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "UpdateAccessCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccessCertificateOptions.Name != nil {
		body["name"] = updateAccessCertificateOptions.Name
	}
	if updateAccessCertificateOptions.AssociatedHostnames != nil {
		body["associated_hostnames"] = updateAccessCertificateOptions.AssociatedHostnames
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAccessCertificate : Delete access certificate
// Delete access certificate.
func (mtls *MtlsV1) DeleteAccessCertificate(deleteAccessCertificateOptions *DeleteAccessCertificateOptions) (result *DeleteAccessCertResp, response *core.DetailedResponse, err error) {
	return mtls.DeleteAccessCertificateWithContext(context.Background(), deleteAccessCertificateOptions)
}

// DeleteAccessCertificateWithContext is an alternate form of the DeleteAccessCertificate method which supports a Context parameter
func (mtls *MtlsV1) DeleteAccessCertificateWithContext(ctx context.Context, deleteAccessCertificateOptions *DeleteAccessCertificateOptions) (result *DeleteAccessCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccessCertificateOptions, "deleteAccessCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAccessCertificateOptions, "deleteAccessCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *deleteAccessCertificateOptions.ZoneID,
		"cert_id": *deleteAccessCertificateOptions.CertID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates/{cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteAccessCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "DeleteAccessCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteAccessCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAccessApplications : List access applications
// List access applications.
func (mtls *MtlsV1) ListAccessApplications(listAccessApplicationsOptions *ListAccessApplicationsOptions) (result *ListAccessAppsResp, response *core.DetailedResponse, err error) {
	return mtls.ListAccessApplicationsWithContext(context.Background(), listAccessApplicationsOptions)
}

// ListAccessApplicationsWithContext is an alternate form of the ListAccessApplications method which supports a Context parameter
func (mtls *MtlsV1) ListAccessApplicationsWithContext(ctx context.Context, listAccessApplicationsOptions *ListAccessApplicationsOptions) (result *ListAccessAppsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAccessApplicationsOptions, "listAccessApplicationsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAccessApplicationsOptions, "listAccessApplicationsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *listAccessApplicationsOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAccessApplicationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "ListAccessApplications")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAccessAppsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAccessApplication : Create access application
// Create access application.
func (mtls *MtlsV1) CreateAccessApplication(createAccessApplicationOptions *CreateAccessApplicationOptions) (result *CreateAccessAppResp, response *core.DetailedResponse, err error) {
	return mtls.CreateAccessApplicationWithContext(context.Background(), createAccessApplicationOptions)
}

// CreateAccessApplicationWithContext is an alternate form of the CreateAccessApplication method which supports a Context parameter
func (mtls *MtlsV1) CreateAccessApplicationWithContext(ctx context.Context, createAccessApplicationOptions *CreateAccessApplicationOptions) (result *CreateAccessAppResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccessApplicationOptions, "createAccessApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createAccessApplicationOptions, "createAccessApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *createAccessApplicationOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAccessApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "CreateAccessApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccessApplicationOptions.Name != nil {
		body["name"] = createAccessApplicationOptions.Name
	}
	if createAccessApplicationOptions.Domain != nil {
		body["domain"] = createAccessApplicationOptions.Domain
	}
	if createAccessApplicationOptions.SessionDuration != nil {
		body["session_duration"] = createAccessApplicationOptions.SessionDuration
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateAccessAppResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAccessApplication : Get access application
// Get access application.
func (mtls *MtlsV1) GetAccessApplication(getAccessApplicationOptions *GetAccessApplicationOptions) (result *AccessAppResp, response *core.DetailedResponse, err error) {
	return mtls.GetAccessApplicationWithContext(context.Background(), getAccessApplicationOptions)
}

// GetAccessApplicationWithContext is an alternate form of the GetAccessApplication method which supports a Context parameter
func (mtls *MtlsV1) GetAccessApplicationWithContext(ctx context.Context, getAccessApplicationOptions *GetAccessApplicationOptions) (result *AccessAppResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccessApplicationOptions, "getAccessApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAccessApplicationOptions, "getAccessApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *getAccessApplicationOptions.ZoneID,
		"app_id": *getAccessApplicationOptions.AppID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAccessApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "GetAccessApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessAppResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccessApplication : Update access application
// Update access application.
func (mtls *MtlsV1) UpdateAccessApplication(updateAccessApplicationOptions *UpdateAccessApplicationOptions) (result *AccessAppResp, response *core.DetailedResponse, err error) {
	return mtls.UpdateAccessApplicationWithContext(context.Background(), updateAccessApplicationOptions)
}

// UpdateAccessApplicationWithContext is an alternate form of the UpdateAccessApplication method which supports a Context parameter
func (mtls *MtlsV1) UpdateAccessApplicationWithContext(ctx context.Context, updateAccessApplicationOptions *UpdateAccessApplicationOptions) (result *AccessAppResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccessApplicationOptions, "updateAccessApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAccessApplicationOptions, "updateAccessApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *updateAccessApplicationOptions.ZoneID,
		"app_id": *updateAccessApplicationOptions.AppID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAccessApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "UpdateAccessApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccessApplicationOptions.Name != nil {
		body["name"] = updateAccessApplicationOptions.Name
	}
	if updateAccessApplicationOptions.Domain != nil {
		body["domain"] = updateAccessApplicationOptions.Domain
	}
	if updateAccessApplicationOptions.SessionDuration != nil {
		body["session_duration"] = updateAccessApplicationOptions.SessionDuration
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessAppResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAccessApplication : Delete access application
// Delete access application.
func (mtls *MtlsV1) DeleteAccessApplication(deleteAccessApplicationOptions *DeleteAccessApplicationOptions) (result *DeleteAccessAppResp, response *core.DetailedResponse, err error) {
	return mtls.DeleteAccessApplicationWithContext(context.Background(), deleteAccessApplicationOptions)
}

// DeleteAccessApplicationWithContext is an alternate form of the DeleteAccessApplication method which supports a Context parameter
func (mtls *MtlsV1) DeleteAccessApplicationWithContext(ctx context.Context, deleteAccessApplicationOptions *DeleteAccessApplicationOptions) (result *DeleteAccessAppResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccessApplicationOptions, "deleteAccessApplicationOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAccessApplicationOptions, "deleteAccessApplicationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *deleteAccessApplicationOptions.ZoneID,
		"app_id": *deleteAccessApplicationOptions.AppID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteAccessApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "DeleteAccessApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteAccessAppResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListAccessPolicies : List access policies
// List access policies.
func (mtls *MtlsV1) ListAccessPolicies(listAccessPoliciesOptions *ListAccessPoliciesOptions) (result *ListAccessPoliciesResp, response *core.DetailedResponse, err error) {
	return mtls.ListAccessPoliciesWithContext(context.Background(), listAccessPoliciesOptions)
}

// ListAccessPoliciesWithContext is an alternate form of the ListAccessPolicies method which supports a Context parameter
func (mtls *MtlsV1) ListAccessPoliciesWithContext(ctx context.Context, listAccessPoliciesOptions *ListAccessPoliciesOptions) (result *ListAccessPoliciesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAccessPoliciesOptions, "listAccessPoliciesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listAccessPoliciesOptions, "listAccessPoliciesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *listAccessPoliciesOptions.ZoneID,
		"app_id": *listAccessPoliciesOptions.AppID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAccessPoliciesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "ListAccessPolicies")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListAccessPoliciesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAccessPolicy : Create access policy
// Create access policy.
func (mtls *MtlsV1) CreateAccessPolicy(createAccessPolicyOptions *CreateAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	return mtls.CreateAccessPolicyWithContext(context.Background(), createAccessPolicyOptions)
}

// CreateAccessPolicyWithContext is an alternate form of the CreateAccessPolicy method which supports a Context parameter
func (mtls *MtlsV1) CreateAccessPolicyWithContext(ctx context.Context, createAccessPolicyOptions *CreateAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAccessPolicyOptions, "createAccessPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createAccessPolicyOptions, "createAccessPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *createAccessPolicyOptions.ZoneID,
		"app_id": *createAccessPolicyOptions.AppID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}/policies`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAccessPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "CreateAccessPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccessPolicyOptions.Name != nil {
		body["name"] = createAccessPolicyOptions.Name
	}
	if createAccessPolicyOptions.Decision != nil {
		body["decision"] = createAccessPolicyOptions.Decision
	}
	if createAccessPolicyOptions.Include != nil {
		body["include"] = createAccessPolicyOptions.Include
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessPolicyResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAccessPolicy : Get access policy
// Get access policy.
func (mtls *MtlsV1) GetAccessPolicy(getAccessPolicyOptions *GetAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	return mtls.GetAccessPolicyWithContext(context.Background(), getAccessPolicyOptions)
}

// GetAccessPolicyWithContext is an alternate form of the GetAccessPolicy method which supports a Context parameter
func (mtls *MtlsV1) GetAccessPolicyWithContext(ctx context.Context, getAccessPolicyOptions *GetAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccessPolicyOptions, "getAccessPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAccessPolicyOptions, "getAccessPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *getAccessPolicyOptions.ZoneID,
		"app_id": *getAccessPolicyOptions.AppID,
		"policy_id": *getAccessPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAccessPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "GetAccessPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessPolicyResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccessPolicy : Update access policy
// Update access policy.
func (mtls *MtlsV1) UpdateAccessPolicy(updateAccessPolicyOptions *UpdateAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	return mtls.UpdateAccessPolicyWithContext(context.Background(), updateAccessPolicyOptions)
}

// UpdateAccessPolicyWithContext is an alternate form of the UpdateAccessPolicy method which supports a Context parameter
func (mtls *MtlsV1) UpdateAccessPolicyWithContext(ctx context.Context, updateAccessPolicyOptions *UpdateAccessPolicyOptions) (result *AccessPolicyResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccessPolicyOptions, "updateAccessPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAccessPolicyOptions, "updateAccessPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *updateAccessPolicyOptions.ZoneID,
		"app_id": *updateAccessPolicyOptions.AppID,
		"policy_id": *updateAccessPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAccessPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "UpdateAccessPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccessPolicyOptions.Name != nil {
		body["name"] = updateAccessPolicyOptions.Name
	}
	if updateAccessPolicyOptions.Decision != nil {
		body["decision"] = updateAccessPolicyOptions.Decision
	}
	if updateAccessPolicyOptions.Include != nil {
		body["include"] = updateAccessPolicyOptions.Include
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessPolicyResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteAccessPolicy : Delete access policy
// Delete access policy.
func (mtls *MtlsV1) DeleteAccessPolicy(deleteAccessPolicyOptions *DeleteAccessPolicyOptions) (result *DeleteAccessPolicyResp, response *core.DetailedResponse, err error) {
	return mtls.DeleteAccessPolicyWithContext(context.Background(), deleteAccessPolicyOptions)
}

// DeleteAccessPolicyWithContext is an alternate form of the DeleteAccessPolicy method which supports a Context parameter
func (mtls *MtlsV1) DeleteAccessPolicyWithContext(ctx context.Context, deleteAccessPolicyOptions *DeleteAccessPolicyOptions) (result *DeleteAccessPolicyResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAccessPolicyOptions, "deleteAccessPolicyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteAccessPolicyOptions, "deleteAccessPolicyOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *deleteAccessPolicyOptions.ZoneID,
		"app_id": *deleteAccessPolicyOptions.AppID,
		"policy_id": *deleteAccessPolicyOptions.PolicyID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/apps/{app_id}/policies/{policy_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteAccessPolicyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "DeleteAccessPolicy")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteAccessPolicyResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetAccessCertSettings : Get access certificates settings
// Get access certificates settings.
func (mtls *MtlsV1) GetAccessCertSettings(getAccessCertSettingsOptions *GetAccessCertSettingsOptions) (result *AccessCertSettingsResp, response *core.DetailedResponse, err error) {
	return mtls.GetAccessCertSettingsWithContext(context.Background(), getAccessCertSettingsOptions)
}

// GetAccessCertSettingsWithContext is an alternate form of the GetAccessCertSettings method which supports a Context parameter
func (mtls *MtlsV1) GetAccessCertSettingsWithContext(ctx context.Context, getAccessCertSettingsOptions *GetAccessCertSettingsOptions) (result *AccessCertSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccessCertSettingsOptions, "getAccessCertSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAccessCertSettingsOptions, "getAccessCertSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *getAccessCertSettingsOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAccessCertSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "GetAccessCertSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessCertSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateAccessCertSettings : Update access certificates settings
// Update access certificates settings.
func (mtls *MtlsV1) UpdateAccessCertSettings(updateAccessCertSettingsOptions *UpdateAccessCertSettingsOptions) (result *AccessCertSettingsResp, response *core.DetailedResponse, err error) {
	return mtls.UpdateAccessCertSettingsWithContext(context.Background(), updateAccessCertSettingsOptions)
}

// UpdateAccessCertSettingsWithContext is an alternate form of the UpdateAccessCertSettings method which supports a Context parameter
func (mtls *MtlsV1) UpdateAccessCertSettingsWithContext(ctx context.Context, updateAccessCertSettingsOptions *UpdateAccessCertSettingsOptions) (result *AccessCertSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAccessCertSettingsOptions, "updateAccessCertSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateAccessCertSettingsOptions, "updateAccessCertSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
		"zone_id": *updateAccessCertSettingsOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/access/certificates/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateAccessCertSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "UpdateAccessCertSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateAccessCertSettingsOptions.Settings != nil {
		body["settings"] = updateAccessCertSettingsOptions.Settings
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessCertSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateAccessOrganization : Create access organization
// Create access organization.
func (mtls *MtlsV1) CreateAccessOrganization(createAccessOrganizationOptions *CreateAccessOrganizationOptions) (result *AccessOrgResp, response *core.DetailedResponse, err error) {
	return mtls.CreateAccessOrganizationWithContext(context.Background(), createAccessOrganizationOptions)
}

// CreateAccessOrganizationWithContext is an alternate form of the CreateAccessOrganization method which supports a Context parameter
func (mtls *MtlsV1) CreateAccessOrganizationWithContext(ctx context.Context, createAccessOrganizationOptions *CreateAccessOrganizationOptions) (result *AccessOrgResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createAccessOrganizationOptions, "createAccessOrganizationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *mtls.Crn,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mtls.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mtls.Service.Options.URL, `/v1/{crn}/access/organizations`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createAccessOrganizationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mtls", "V1", "CreateAccessOrganization")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createAccessOrganizationOptions.Name != nil {
		body["name"] = createAccessOrganizationOptions.Name
	}
	if createAccessOrganizationOptions.AuthDomain != nil {
		body["auth_domain"] = createAccessOrganizationOptions.AuthDomain
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
	response, err = mtls.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccessOrgResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AccessOrgRespResult : AccessOrgRespResult struct
type AccessOrgRespResult struct {
	AuthDomain *string `json:"auth_domain,omitempty"`

	Name *string `json:"name,omitempty"`

	LoginDesign interface{} `json:"login_design,omitempty"`

	CreatedAt *string `json:"created_at,omitempty"`

	UpdatedAt *string `json:"updated_at,omitempty"`
}

// UnmarshalAccessOrgRespResult unmarshals an instance of AccessOrgRespResult from the specified map of raw messages.
func UnmarshalAccessOrgRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessOrgRespResult)
	err = core.UnmarshalPrimitive(m, "auth_domain", &obj.AuthDomain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "login_design", &obj.LoginDesign)
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

// CreateAccessAppRespResult : Access application details.
type CreateAccessAppRespResult struct {
	ID *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`

	Domain *string `json:"domain,omitempty"`

	Aud *string `json:"aud,omitempty"`

	Policies []interface{} `json:"policies,omitempty"`

	AllowedIdps []string `json:"allowed_idps,omitempty"`

	AutoRedirectToIdentity *bool `json:"auto_redirect_to_identity,omitempty"`

	SessionDuration *string `json:"session_duration,omitempty"`

	Type *string `json:"type,omitempty"`

	Uid *string `json:"uid,omitempty"`

	CreatedAt *string `json:"created_at,omitempty"`

	UpdatedAt *string `json:"updated_at,omitempty"`
}

// UnmarshalCreateAccessAppRespResult unmarshals an instance of CreateAccessAppRespResult from the specified map of raw messages.
func UnmarshalCreateAccessAppRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccessAppRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "domain", &obj.Domain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aud", &obj.Aud)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "policies", &obj.Policies)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_idps", &obj.AllowedIdps)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_redirect_to_identity", &obj.AutoRedirectToIdentity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_duration", &obj.SessionDuration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uid", &obj.Uid)
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

// CreateAccessApplicationOptions : The CreateAccessApplication options.
type CreateAccessApplicationOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Application name.
	Name *string `json:"name,omitempty"`

	// The domain and path that Access blocks.
	Domain *string `json:"domain,omitempty"`

	// The amount of time that the tokens issued for this application are valid.
	SessionDuration *string `json:"session_duration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAccessApplicationOptions : Instantiate CreateAccessApplicationOptions
func (*MtlsV1) NewCreateAccessApplicationOptions(zoneID string) *CreateAccessApplicationOptions {
	return &CreateAccessApplicationOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *CreateAccessApplicationOptions) SetZoneID(zoneID string) *CreateAccessApplicationOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccessApplicationOptions) SetName(name string) *CreateAccessApplicationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDomain : Allow user to set Domain
func (_options *CreateAccessApplicationOptions) SetDomain(domain string) *CreateAccessApplicationOptions {
	_options.Domain = core.StringPtr(domain)
	return _options
}

// SetSessionDuration : Allow user to set SessionDuration
func (_options *CreateAccessApplicationOptions) SetSessionDuration(sessionDuration string) *CreateAccessApplicationOptions {
	_options.SessionDuration = core.StringPtr(sessionDuration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccessApplicationOptions) SetHeaders(param map[string]string) *CreateAccessApplicationOptions {
	options.Headers = param
	return options
}

// CreateAccessCertificateOptions : The CreateAccessCertificate options.
type CreateAccessCertificateOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access certificate name.
	Name *string `json:"name,omitempty"`

	// Access certificate.
	Certificate *string `json:"certificate,omitempty"`

	// The hostnames that are prompted for this certificate.
	AssociatedHostnames []string `json:"associated_hostnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAccessCertificateOptions : Instantiate CreateAccessCertificateOptions
func (*MtlsV1) NewCreateAccessCertificateOptions(zoneID string) *CreateAccessCertificateOptions {
	return &CreateAccessCertificateOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *CreateAccessCertificateOptions) SetZoneID(zoneID string) *CreateAccessCertificateOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccessCertificateOptions) SetName(name string) *CreateAccessCertificateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCertificate : Allow user to set Certificate
func (_options *CreateAccessCertificateOptions) SetCertificate(certificate string) *CreateAccessCertificateOptions {
	_options.Certificate = core.StringPtr(certificate)
	return _options
}

// SetAssociatedHostnames : Allow user to set AssociatedHostnames
func (_options *CreateAccessCertificateOptions) SetAssociatedHostnames(associatedHostnames []string) *CreateAccessCertificateOptions {
	_options.AssociatedHostnames = associatedHostnames
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccessCertificateOptions) SetHeaders(param map[string]string) *CreateAccessCertificateOptions {
	options.Headers = param
	return options
}

// CreateAccessOrganizationOptions : The CreateAccessOrganization options.
type CreateAccessOrganizationOptions struct {
	// Name of the access organization.
	Name *string `json:"name,omitempty"`

	// The domain that you are redirected to on Access login attempts.
	AuthDomain *string `json:"auth_domain,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateAccessOrganizationOptions : Instantiate CreateAccessOrganizationOptions
func (*MtlsV1) NewCreateAccessOrganizationOptions() *CreateAccessOrganizationOptions {
	return &CreateAccessOrganizationOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateAccessOrganizationOptions) SetName(name string) *CreateAccessOrganizationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAuthDomain : Allow user to set AuthDomain
func (_options *CreateAccessOrganizationOptions) SetAuthDomain(authDomain string) *CreateAccessOrganizationOptions {
	_options.AuthDomain = core.StringPtr(authDomain)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccessOrganizationOptions) SetHeaders(param map[string]string) *CreateAccessOrganizationOptions {
	options.Headers = param
	return options
}

// CreateAccessPolicyOptions : The CreateAccessPolicy options.
type CreateAccessPolicyOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Policy name.
	Name *string `json:"name,omitempty"`

	// Defines the action Access takes if the policy matches the user.
	Decision *string `json:"decision,omitempty"`

	// The include policy works like an OR logical operator. The user must satisfy one of the rules in includes.
	Include []PolicyRuleIntf `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateAccessPolicyOptions.Decision property.
// Defines the action Access takes if the policy matches the user.
const (
	CreateAccessPolicyOptions_Decision_NonIdentity = "non_identity"
)

// NewCreateAccessPolicyOptions : Instantiate CreateAccessPolicyOptions
func (*MtlsV1) NewCreateAccessPolicyOptions(zoneID string, appID string) *CreateAccessPolicyOptions {
	return &CreateAccessPolicyOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *CreateAccessPolicyOptions) SetZoneID(zoneID string) *CreateAccessPolicyOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *CreateAccessPolicyOptions) SetAppID(appID string) *CreateAccessPolicyOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAccessPolicyOptions) SetName(name string) *CreateAccessPolicyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDecision : Allow user to set Decision
func (_options *CreateAccessPolicyOptions) SetDecision(decision string) *CreateAccessPolicyOptions {
	_options.Decision = core.StringPtr(decision)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *CreateAccessPolicyOptions) SetInclude(include []PolicyRuleIntf) *CreateAccessPolicyOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAccessPolicyOptions) SetHeaders(param map[string]string) *CreateAccessPolicyOptions {
	options.Headers = param
	return options
}

// DeleteAccessAppRespResult : DeleteAccessAppRespResult struct
type DeleteAccessAppRespResult struct {
	// Application ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalDeleteAccessAppRespResult unmarshals an instance of DeleteAccessAppRespResult from the specified map of raw messages.
func UnmarshalDeleteAccessAppRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessAppRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccessApplicationOptions : The DeleteAccessApplication options.
type DeleteAccessApplicationOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAccessApplicationOptions : Instantiate DeleteAccessApplicationOptions
func (*MtlsV1) NewDeleteAccessApplicationOptions(zoneID string, appID string) *DeleteAccessApplicationOptions {
	return &DeleteAccessApplicationOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *DeleteAccessApplicationOptions) SetZoneID(zoneID string) *DeleteAccessApplicationOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *DeleteAccessApplicationOptions) SetAppID(appID string) *DeleteAccessApplicationOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccessApplicationOptions) SetHeaders(param map[string]string) *DeleteAccessApplicationOptions {
	options.Headers = param
	return options
}

// DeleteAccessCertRespResult : DeleteAccessCertRespResult struct
type DeleteAccessCertRespResult struct {
	// Certificate ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalDeleteAccessCertRespResult unmarshals an instance of DeleteAccessCertRespResult from the specified map of raw messages.
func UnmarshalDeleteAccessCertRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessCertRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccessCertificateOptions : The DeleteAccessCertificate options.
type DeleteAccessCertificateOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access certificate ID.
	CertID *string `json:"cert_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAccessCertificateOptions : Instantiate DeleteAccessCertificateOptions
func (*MtlsV1) NewDeleteAccessCertificateOptions(zoneID string, certID string) *DeleteAccessCertificateOptions {
	return &DeleteAccessCertificateOptions{
		ZoneID: core.StringPtr(zoneID),
		CertID: core.StringPtr(certID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *DeleteAccessCertificateOptions) SetZoneID(zoneID string) *DeleteAccessCertificateOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetCertID : Allow user to set CertID
func (_options *DeleteAccessCertificateOptions) SetCertID(certID string) *DeleteAccessCertificateOptions {
	_options.CertID = core.StringPtr(certID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccessCertificateOptions) SetHeaders(param map[string]string) *DeleteAccessCertificateOptions {
	options.Headers = param
	return options
}

// DeleteAccessPolicyOptions : The DeleteAccessPolicy options.
type DeleteAccessPolicyOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Access policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteAccessPolicyOptions : Instantiate DeleteAccessPolicyOptions
func (*MtlsV1) NewDeleteAccessPolicyOptions(zoneID string, appID string, policyID string) *DeleteAccessPolicyOptions {
	return &DeleteAccessPolicyOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
		PolicyID: core.StringPtr(policyID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *DeleteAccessPolicyOptions) SetZoneID(zoneID string) *DeleteAccessPolicyOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *DeleteAccessPolicyOptions) SetAppID(appID string) *DeleteAccessPolicyOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetPolicyID : Allow user to set PolicyID
func (_options *DeleteAccessPolicyOptions) SetPolicyID(policyID string) *DeleteAccessPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAccessPolicyOptions) SetHeaders(param map[string]string) *DeleteAccessPolicyOptions {
	options.Headers = param
	return options
}

// DeleteAccessPolicyRespResult : DeleteAccessPolicyRespResult struct
type DeleteAccessPolicyRespResult struct {
	// Policy ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalDeleteAccessPolicyRespResult unmarshals an instance of DeleteAccessPolicyRespResult from the specified map of raw messages.
func UnmarshalDeleteAccessPolicyRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessPolicyRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAccessApplicationOptions : The GetAccessApplication options.
type GetAccessApplicationOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccessApplicationOptions : Instantiate GetAccessApplicationOptions
func (*MtlsV1) NewGetAccessApplicationOptions(zoneID string, appID string) *GetAccessApplicationOptions {
	return &GetAccessApplicationOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *GetAccessApplicationOptions) SetZoneID(zoneID string) *GetAccessApplicationOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *GetAccessApplicationOptions) SetAppID(appID string) *GetAccessApplicationOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccessApplicationOptions) SetHeaders(param map[string]string) *GetAccessApplicationOptions {
	options.Headers = param
	return options
}

// GetAccessCertSettingsOptions : The GetAccessCertSettings options.
type GetAccessCertSettingsOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccessCertSettingsOptions : Instantiate GetAccessCertSettingsOptions
func (*MtlsV1) NewGetAccessCertSettingsOptions(zoneID string) *GetAccessCertSettingsOptions {
	return &GetAccessCertSettingsOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *GetAccessCertSettingsOptions) SetZoneID(zoneID string) *GetAccessCertSettingsOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccessCertSettingsOptions) SetHeaders(param map[string]string) *GetAccessCertSettingsOptions {
	options.Headers = param
	return options
}

// GetAccessCertificateOptions : The GetAccessCertificate options.
type GetAccessCertificateOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access certificate ID.
	CertID *string `json:"cert_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccessCertificateOptions : Instantiate GetAccessCertificateOptions
func (*MtlsV1) NewGetAccessCertificateOptions(zoneID string, certID string) *GetAccessCertificateOptions {
	return &GetAccessCertificateOptions{
		ZoneID: core.StringPtr(zoneID),
		CertID: core.StringPtr(certID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *GetAccessCertificateOptions) SetZoneID(zoneID string) *GetAccessCertificateOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetCertID : Allow user to set CertID
func (_options *GetAccessCertificateOptions) SetCertID(certID string) *GetAccessCertificateOptions {
	_options.CertID = core.StringPtr(certID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccessCertificateOptions) SetHeaders(param map[string]string) *GetAccessCertificateOptions {
	options.Headers = param
	return options
}

// GetAccessPolicyOptions : The GetAccessPolicy options.
type GetAccessPolicyOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Access policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccessPolicyOptions : Instantiate GetAccessPolicyOptions
func (*MtlsV1) NewGetAccessPolicyOptions(zoneID string, appID string, policyID string) *GetAccessPolicyOptions {
	return &GetAccessPolicyOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
		PolicyID: core.StringPtr(policyID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *GetAccessPolicyOptions) SetZoneID(zoneID string) *GetAccessPolicyOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *GetAccessPolicyOptions) SetAppID(appID string) *GetAccessPolicyOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetPolicyID : Allow user to set PolicyID
func (_options *GetAccessPolicyOptions) SetPolicyID(policyID string) *GetAccessPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccessPolicyOptions) SetHeaders(param map[string]string) *GetAccessPolicyOptions {
	options.Headers = param
	return options
}

// ListAccessApplicationsOptions : The ListAccessApplications options.
type ListAccessApplicationsOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAccessApplicationsOptions : Instantiate ListAccessApplicationsOptions
func (*MtlsV1) NewListAccessApplicationsOptions(zoneID string) *ListAccessApplicationsOptions {
	return &ListAccessApplicationsOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *ListAccessApplicationsOptions) SetZoneID(zoneID string) *ListAccessApplicationsOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccessApplicationsOptions) SetHeaders(param map[string]string) *ListAccessApplicationsOptions {
	options.Headers = param
	return options
}

// ListAccessCertificatesOptions : The ListAccessCertificates options.
type ListAccessCertificatesOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAccessCertificatesOptions : Instantiate ListAccessCertificatesOptions
func (*MtlsV1) NewListAccessCertificatesOptions(zoneID string) *ListAccessCertificatesOptions {
	return &ListAccessCertificatesOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *ListAccessCertificatesOptions) SetZoneID(zoneID string) *ListAccessCertificatesOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccessCertificatesOptions) SetHeaders(param map[string]string) *ListAccessCertificatesOptions {
	options.Headers = param
	return options
}

// ListAccessPoliciesOptions : The ListAccessPolicies options.
type ListAccessPoliciesOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAccessPoliciesOptions : Instantiate ListAccessPoliciesOptions
func (*MtlsV1) NewListAccessPoliciesOptions(zoneID string, appID string) *ListAccessPoliciesOptions {
	return &ListAccessPoliciesOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *ListAccessPoliciesOptions) SetZoneID(zoneID string) *ListAccessPoliciesOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *ListAccessPoliciesOptions) SetAppID(appID string) *ListAccessPoliciesOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAccessPoliciesOptions) SetHeaders(param map[string]string) *ListAccessPoliciesOptions {
	options.Headers = param
	return options
}

// PolicyCnRuleCommonName : PolicyCnRuleCommonName struct
type PolicyCnRuleCommonName struct {
	// Common name of client certificate.
	CommonName *string `json:"common_name" validate:"required"`
}

// NewPolicyCnRuleCommonName : Instantiate PolicyCnRuleCommonName (Generic Model Constructor)
func (*MtlsV1) NewPolicyCnRuleCommonName(commonName string) (_model *PolicyCnRuleCommonName, err error) {
	_model = &PolicyCnRuleCommonName{
		CommonName: core.StringPtr(commonName),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalPolicyCnRuleCommonName unmarshals an instance of PolicyCnRuleCommonName from the specified map of raw messages.
func UnmarshalPolicyCnRuleCommonName(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyCnRuleCommonName)
	err = core.UnmarshalPrimitive(m, "common_name", &obj.CommonName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAccessApplicationOptions : The UpdateAccessApplication options.
type UpdateAccessApplicationOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Application name.
	Name *string `json:"name,omitempty"`

	// The domain and path that Access blocks.
	Domain *string `json:"domain,omitempty"`

	// The amount of time that the tokens issued for this application are valid.
	SessionDuration *string `json:"session_duration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccessApplicationOptions : Instantiate UpdateAccessApplicationOptions
func (*MtlsV1) NewUpdateAccessApplicationOptions(zoneID string, appID string) *UpdateAccessApplicationOptions {
	return &UpdateAccessApplicationOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *UpdateAccessApplicationOptions) SetZoneID(zoneID string) *UpdateAccessApplicationOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *UpdateAccessApplicationOptions) SetAppID(appID string) *UpdateAccessApplicationOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAccessApplicationOptions) SetName(name string) *UpdateAccessApplicationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDomain : Allow user to set Domain
func (_options *UpdateAccessApplicationOptions) SetDomain(domain string) *UpdateAccessApplicationOptions {
	_options.Domain = core.StringPtr(domain)
	return _options
}

// SetSessionDuration : Allow user to set SessionDuration
func (_options *UpdateAccessApplicationOptions) SetSessionDuration(sessionDuration string) *UpdateAccessApplicationOptions {
	_options.SessionDuration = core.StringPtr(sessionDuration)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccessApplicationOptions) SetHeaders(param map[string]string) *UpdateAccessApplicationOptions {
	options.Headers = param
	return options
}

// UpdateAccessCertSettingsOptions : The UpdateAccessCertSettings options.
type UpdateAccessCertSettingsOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	Settings []AccessCertSettingsInputArray `json:"settings,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccessCertSettingsOptions : Instantiate UpdateAccessCertSettingsOptions
func (*MtlsV1) NewUpdateAccessCertSettingsOptions(zoneID string) *UpdateAccessCertSettingsOptions {
	return &UpdateAccessCertSettingsOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *UpdateAccessCertSettingsOptions) SetZoneID(zoneID string) *UpdateAccessCertSettingsOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetSettings : Allow user to set Settings
func (_options *UpdateAccessCertSettingsOptions) SetSettings(settings []AccessCertSettingsInputArray) *UpdateAccessCertSettingsOptions {
	_options.Settings = settings
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccessCertSettingsOptions) SetHeaders(param map[string]string) *UpdateAccessCertSettingsOptions {
	options.Headers = param
	return options
}

// UpdateAccessCertificateOptions : The UpdateAccessCertificate options.
type UpdateAccessCertificateOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access certificate ID.
	CertID *string `json:"cert_id" validate:"required,ne="`

	// Access certificate name.
	Name *string `json:"name,omitempty"`

	// The hostnames that are prompted for this certificate.
	AssociatedHostnames []string `json:"associated_hostnames,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateAccessCertificateOptions : Instantiate UpdateAccessCertificateOptions
func (*MtlsV1) NewUpdateAccessCertificateOptions(zoneID string, certID string) *UpdateAccessCertificateOptions {
	return &UpdateAccessCertificateOptions{
		ZoneID: core.StringPtr(zoneID),
		CertID: core.StringPtr(certID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *UpdateAccessCertificateOptions) SetZoneID(zoneID string) *UpdateAccessCertificateOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetCertID : Allow user to set CertID
func (_options *UpdateAccessCertificateOptions) SetCertID(certID string) *UpdateAccessCertificateOptions {
	_options.CertID = core.StringPtr(certID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAccessCertificateOptions) SetName(name string) *UpdateAccessCertificateOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAssociatedHostnames : Allow user to set AssociatedHostnames
func (_options *UpdateAccessCertificateOptions) SetAssociatedHostnames(associatedHostnames []string) *UpdateAccessCertificateOptions {
	_options.AssociatedHostnames = associatedHostnames
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccessCertificateOptions) SetHeaders(param map[string]string) *UpdateAccessCertificateOptions {
	options.Headers = param
	return options
}

// UpdateAccessPolicyOptions : The UpdateAccessPolicy options.
type UpdateAccessPolicyOptions struct {
	// Zone ID.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// Access application ID.
	AppID *string `json:"app_id" validate:"required,ne="`

	// Access policy ID.
	PolicyID *string `json:"policy_id" validate:"required,ne="`

	// Policy name.
	Name *string `json:"name,omitempty"`

	// Defines the action Access takes if the policy matches the user.
	Decision *string `json:"decision,omitempty"`

	// The include policy works like an OR logical operator. The user must satisfy one of the rules in includes.
	Include []PolicyRuleIntf `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateAccessPolicyOptions.Decision property.
// Defines the action Access takes if the policy matches the user.
const (
	UpdateAccessPolicyOptions_Decision_NonIdentity = "non_identity"
)

// NewUpdateAccessPolicyOptions : Instantiate UpdateAccessPolicyOptions
func (*MtlsV1) NewUpdateAccessPolicyOptions(zoneID string, appID string, policyID string) *UpdateAccessPolicyOptions {
	return &UpdateAccessPolicyOptions{
		ZoneID: core.StringPtr(zoneID),
		AppID: core.StringPtr(appID),
		PolicyID: core.StringPtr(policyID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *UpdateAccessPolicyOptions) SetZoneID(zoneID string) *UpdateAccessPolicyOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetAppID : Allow user to set AppID
func (_options *UpdateAccessPolicyOptions) SetAppID(appID string) *UpdateAccessPolicyOptions {
	_options.AppID = core.StringPtr(appID)
	return _options
}

// SetPolicyID : Allow user to set PolicyID
func (_options *UpdateAccessPolicyOptions) SetPolicyID(policyID string) *UpdateAccessPolicyOptions {
	_options.PolicyID = core.StringPtr(policyID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAccessPolicyOptions) SetName(name string) *UpdateAccessPolicyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDecision : Allow user to set Decision
func (_options *UpdateAccessPolicyOptions) SetDecision(decision string) *UpdateAccessPolicyOptions {
	_options.Decision = core.StringPtr(decision)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *UpdateAccessPolicyOptions) SetInclude(include []PolicyRuleIntf) *UpdateAccessPolicyOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAccessPolicyOptions) SetHeaders(param map[string]string) *UpdateAccessPolicyOptions {
	options.Headers = param
	return options
}

// AccessAppResp : Access application response.
type AccessAppResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	// Access application details.
	Result *AppResult `json:"result,omitempty"`
}

// UnmarshalAccessAppResp unmarshals an instance of AccessAppResp from the specified map of raw messages.
func UnmarshalAccessAppResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessAppResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalAppResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessCertResp : Access certificate response.
type AccessCertResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	// Access certificate details.
	Result *CertResult `json:"result,omitempty"`
}

// UnmarshalAccessCertResp unmarshals an instance of AccessCertResp from the specified map of raw messages.
func UnmarshalAccessCertResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessCertResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCertResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessCertSettingsInputArray : AccessCertSettingsInputArray struct
type AccessCertSettingsInputArray struct {
	Hostname *string `json:"hostname" validate:"required"`

	// Whether to forward the client certificate.
	ClientCertificateForwarding *bool `json:"client_certificate_forwarding" validate:"required"`
}

// NewAccessCertSettingsInputArray : Instantiate AccessCertSettingsInputArray (Generic Model Constructor)
func (*MtlsV1) NewAccessCertSettingsInputArray(hostname string, clientCertificateForwarding bool) (_model *AccessCertSettingsInputArray, err error) {
	_model = &AccessCertSettingsInputArray{
		Hostname: core.StringPtr(hostname),
		ClientCertificateForwarding: core.BoolPtr(clientCertificateForwarding),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalAccessCertSettingsInputArray unmarshals an instance of AccessCertSettingsInputArray from the specified map of raw messages.
func UnmarshalAccessCertSettingsInputArray(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessCertSettingsInputArray)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_certificate_forwarding", &obj.ClientCertificateForwarding)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessCertSettingsResp : Access certificates settings response.
type AccessCertSettingsResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result []CertSettingsResult `json:"result,omitempty"`
}

// UnmarshalAccessCertSettingsResp unmarshals an instance of AccessCertSettingsResp from the specified map of raw messages.
func UnmarshalAccessCertSettingsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessCertSettingsResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCertSettingsResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessOrgResp : Access organization response.
type AccessOrgResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result *AccessOrgRespResult `json:"result,omitempty"`
}

// UnmarshalAccessOrgResp unmarshals an instance of AccessOrgResp from the specified map of raw messages.
func UnmarshalAccessOrgResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessOrgResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalAccessOrgRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccessPolicyResp : Access policy response.
type AccessPolicyResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	// Access policies information.
	Result *PolicyResult `json:"result,omitempty"`
}

// UnmarshalAccessPolicyResp unmarshals an instance of AccessPolicyResp from the specified map of raw messages.
func UnmarshalAccessPolicyResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccessPolicyResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPolicyResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AppResult : Access application details.
type AppResult struct {
	// Application ID.
	ID *string `json:"id,omitempty"`

	// Application name.
	Name *string `json:"name,omitempty"`

	// The domain and path that Access blocks.
	Domain *string `json:"domain,omitempty"`

	Aud *string `json:"aud,omitempty"`

	// Policies of the application.
	Policies []PolicyResult `json:"policies,omitempty"`

	// The identity providers selected for application.
	AllowedIdps []string `json:"allowed_idps,omitempty"`

	// Option to skip identity provider selection if only one is configured in allowed_idps.
	AutoRedirectToIdentity *bool `json:"auto_redirect_to_identity,omitempty"`

	// The amount of time that the tokens issued for this application are valid.
	SessionDuration *string `json:"session_duration,omitempty"`

	// Application type.
	Type *string `json:"type,omitempty"`

	// UUID, same as ID.
	Uid *string `json:"uid,omitempty"`

	// Created time of the application.
	CreatedAt *string `json:"created_at,omitempty"`

	// Updated time of the application.
	UpdatedAt *string `json:"updated_at,omitempty"`
}

// UnmarshalAppResult unmarshals an instance of AppResult from the specified map of raw messages.
func UnmarshalAppResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "domain", &obj.Domain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "aud", &obj.Aud)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "policies", &obj.Policies, UnmarshalPolicyResult)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_idps", &obj.AllowedIdps)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auto_redirect_to_identity", &obj.AutoRedirectToIdentity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_duration", &obj.SessionDuration)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uid", &obj.Uid)
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

// CertResult : Access certificate details.
type CertResult struct {
	// Access certificate ID.
	ID *string `json:"id,omitempty"`

	// access certificate name.
	Name *string `json:"name,omitempty"`

	// Fingerprint of the certificate.
	Fingerprint *string `json:"fingerprint,omitempty"`

	// The hostnames that are prompted for this certificate.
	AssociatedHostnames []string `json:"associated_hostnames,omitempty"`

	// Created time of the access certificate.
	CreatedAt *string `json:"created_at,omitempty"`

	// Updated time of the access certificate.
	UpdatedAt *string `json:"updated_at,omitempty"`

	// Expire time of the access certificate.
	ExpiresOn *string `json:"expires_on,omitempty"`
}

// UnmarshalCertResult unmarshals an instance of CertResult from the specified map of raw messages.
func UnmarshalCertResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "fingerprint", &obj.Fingerprint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "associated_hostnames", &obj.AssociatedHostnames)
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
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CertSettingsResult : CertSettingsResult struct
type CertSettingsResult struct {
	Hostname *string `json:"hostname,omitempty"`

	ChinaNetwork *bool `json:"china_network,omitempty"`

	ClientCertificateForwarding *bool `json:"client_certificate_forwarding,omitempty"`
}

// UnmarshalCertSettingsResult unmarshals an instance of CertSettingsResult from the specified map of raw messages.
func UnmarshalCertSettingsResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertSettingsResult)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "china_network", &obj.ChinaNetwork)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_certificate_forwarding", &obj.ClientCertificateForwarding)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAccessAppResp : Create access application response.
type CreateAccessAppResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	// Access application details.
	Result *CreateAccessAppRespResult `json:"result,omitempty"`
}

// UnmarshalCreateAccessAppResp unmarshals an instance of CreateAccessAppResp from the specified map of raw messages.
func UnmarshalCreateAccessAppResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateAccessAppResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCreateAccessAppRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccessAppResp : Delete access application response.
type DeleteAccessAppResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result *DeleteAccessAppRespResult `json:"result,omitempty"`
}

// UnmarshalDeleteAccessAppResp unmarshals an instance of DeleteAccessAppResp from the specified map of raw messages.
func UnmarshalDeleteAccessAppResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessAppResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteAccessAppRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccessCertResp : Delete access certificate response.
type DeleteAccessCertResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result *DeleteAccessCertRespResult `json:"result,omitempty"`
}

// UnmarshalDeleteAccessCertResp unmarshals an instance of DeleteAccessCertResp from the specified map of raw messages.
func UnmarshalDeleteAccessCertResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessCertResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteAccessCertRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteAccessPolicyResp : Delete access policy response.
type DeleteAccessPolicyResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result *DeleteAccessPolicyRespResult `json:"result,omitempty"`
}

// UnmarshalDeleteAccessPolicyResp unmarshals an instance of DeleteAccessPolicyResp from the specified map of raw messages.
func UnmarshalDeleteAccessPolicyResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteAccessPolicyResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteAccessPolicyRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAccessAppsResp : List access applications response.
type ListAccessAppsResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result []AppResult `json:"result,omitempty"`
}

// UnmarshalListAccessAppsResp unmarshals an instance of ListAccessAppsResp from the specified map of raw messages.
func UnmarshalListAccessAppsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAccessAppsResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalAppResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAccessCertsResp : List access certificate response.
type ListAccessCertsResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result []CertResult `json:"result,omitempty"`
}

// UnmarshalListAccessCertsResp unmarshals an instance of ListAccessCertsResp from the specified map of raw messages.
func UnmarshalListAccessCertsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAccessCertsResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCertResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAccessPoliciesResp : List access policies response.
type ListAccessPoliciesResp struct {
	// Was operation successful.
	Success *bool `json:"success,omitempty"`

	// Array of errors encountered.
	Errors [][]string `json:"errors,omitempty"`

	// Array of messages returned.
	Messages [][]string `json:"messages,omitempty"`

	Result []PolicyResult `json:"result,omitempty"`
}

// UnmarshalListAccessPoliciesResp unmarshals an instance of ListAccessPoliciesResp from the specified map of raw messages.
func UnmarshalListAccessPoliciesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListAccessPoliciesResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPolicyResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyResult : Access policies information.
type PolicyResult struct {
	// Policy ID.
	ID *string `json:"id,omitempty"`

	// Policy name.
	Name *string `json:"name,omitempty"`

	// The action Access takes if the policy matches the user.
	Decision *string `json:"decision,omitempty"`

	// The include policy works like an OR logical operator.
	Include []PolicyRuleIntf `json:"include,omitempty"`

	// The exclude policy works like a NOT logical operator.
	Exclude []PolicyRuleIntf `json:"exclude,omitempty"`

	// The unique precedence for policies on a single application.
	Precedence *int64 `json:"precedence,omitempty"`

	// The require policy works like a AND logical operator.
	Require []PolicyRuleIntf `json:"require,omitempty"`

	// UUID, same as ID.
	Uid *string `json:"uid,omitempty"`

	// Created time of the policy.
	CreatedAt *string `json:"created_at,omitempty"`

	// Updated time of the policy.
	UpdatedAt *string `json:"updated_at,omitempty"`
}

// UnmarshalPolicyResult unmarshals an instance of PolicyResult from the specified map of raw messages.
func UnmarshalPolicyResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "decision", &obj.Decision)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "include", &obj.Include, UnmarshalPolicyRule)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "exclude", &obj.Exclude, UnmarshalPolicyRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "precedence", &obj.Precedence)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "require", &obj.Require, UnmarshalPolicyRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uid", &obj.Uid)
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

// PolicyRule : Policy rule.
// Models which "extend" this model:
// - PolicyRulePolicyCertRule
// - PolicyRulePolicyCnRule
type PolicyRule struct {
	Certificate interface{} `json:"certificate,omitempty"`

	CommonName *PolicyCnRuleCommonName `json:"common_name,omitempty"`
}
func (*PolicyRule) isaPolicyRule() bool {
	return true
}

type PolicyRuleIntf interface {
	isaPolicyRule() bool
}

// UnmarshalPolicyRule unmarshals an instance of PolicyRule from the specified map of raw messages.
func UnmarshalPolicyRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyRule)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "common_name", &obj.CommonName, UnmarshalPolicyCnRuleCommonName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyRulePolicyCertRule : Policy rule of certificate.
// This model "extends" PolicyRule
type PolicyRulePolicyCertRule struct {
	Certificate interface{} `json:"certificate,omitempty"`
}

func (*PolicyRulePolicyCertRule) isaPolicyRule() bool {
	return true
}

// UnmarshalPolicyRulePolicyCertRule unmarshals an instance of PolicyRulePolicyCertRule from the specified map of raw messages.
func UnmarshalPolicyRulePolicyCertRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyRulePolicyCertRule)
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PolicyRulePolicyCnRule : Policy rule of common name.
// This model "extends" PolicyRule
type PolicyRulePolicyCnRule struct {
	CommonName *PolicyCnRuleCommonName `json:"common_name" validate:"required"`
}

// NewPolicyRulePolicyCnRule : Instantiate PolicyRulePolicyCnRule (Generic Model Constructor)
func (*MtlsV1) NewPolicyRulePolicyCnRule(commonName *PolicyCnRuleCommonName) (_model *PolicyRulePolicyCnRule, err error) {
	_model = &PolicyRulePolicyCnRule{
		CommonName: commonName,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*PolicyRulePolicyCnRule) isaPolicyRule() bool {
	return true
}

// UnmarshalPolicyRulePolicyCnRule unmarshals an instance of PolicyRulePolicyCnRule from the specified map of raw messages.
func UnmarshalPolicyRulePolicyCnRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PolicyRulePolicyCnRule)
	err = core.UnmarshalModel(m, "common_name", &obj.CommonName, UnmarshalPolicyCnRuleCommonName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
