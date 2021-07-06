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
 

// Package sslcertificateapiv1 : Operations and models for the SslCertificateApiV1 service
package sslcertificateapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
	"net/http"
	"reflect"
	"time"
)

// SslCertificateApiV1 : SSL Certificate
//
// Version: 1.0.0
type SslCertificateApiV1 struct {
	Service *core.BaseService

	// cloud resource name.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "ssl_certificate_api"

// SslCertificateApiV1Options : Service options
type SslCertificateApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// cloud resource name.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewSslCertificateApiV1UsingExternalConfig : constructs an instance of SslCertificateApiV1 with passed in options and external configuration.
func NewSslCertificateApiV1UsingExternalConfig(options *SslCertificateApiV1Options) (sslCertificateApi *SslCertificateApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	sslCertificateApi, err = NewSslCertificateApiV1(options)
	if err != nil {
		return
	}

	err = sslCertificateApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = sslCertificateApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewSslCertificateApiV1 : constructs an instance of SslCertificateApiV1 with passed in options.
func NewSslCertificateApiV1(options *SslCertificateApiV1Options) (service *SslCertificateApiV1, err error) {
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

	service = &SslCertificateApiV1{
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

// Clone makes a copy of "sslCertificateApi" suitable for processing requests.
func (sslCertificateApi *SslCertificateApiV1) Clone() *SslCertificateApiV1 {
	if core.IsNil(sslCertificateApi) {
		return nil
	}
	clone := *sslCertificateApi
	clone.Service = sslCertificateApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (sslCertificateApi *SslCertificateApiV1) SetServiceURL(url string) error {
	return sslCertificateApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (sslCertificateApi *SslCertificateApiV1) GetServiceURL() string {
	return sslCertificateApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (sslCertificateApi *SslCertificateApiV1) SetDefaultHeaders(headers http.Header) {
	sslCertificateApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (sslCertificateApi *SslCertificateApiV1) SetEnableGzipCompression(enableGzip bool) {
	sslCertificateApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (sslCertificateApi *SslCertificateApiV1) GetEnableGzipCompression() bool {
	return sslCertificateApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (sslCertificateApi *SslCertificateApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	sslCertificateApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (sslCertificateApi *SslCertificateApiV1) DisableRetries() {
	sslCertificateApi.Service.DisableRetries()
}

// ListCertificates : List all certificates
// CIS automatically add an active DNS zone to a universal SSL certificate, shared among multiple customers. Customer
// may order dedicated certificates for the owning zones. This API list all certificates for a given zone, including
// shared and dedicated certificates.
func (sslCertificateApi *SslCertificateApiV1) ListCertificates(listCertificatesOptions *ListCertificatesOptions) (result *ListCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ListCertificatesWithContext(context.Background(), listCertificatesOptions)
}

// ListCertificatesWithContext is an alternate form of the ListCertificates method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ListCertificatesWithContext(ctx context.Context, listCertificatesOptions *ListCertificatesOptions) (result *ListCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCertificatesOptions, "listCertificatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/ssl/certificate_packs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ListCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listCertificatesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listCertificatesOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCertificateResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// OrderCertificate : Order dedicated certificate
// Order a dedicated certificate for a given zone. The zone should be active before placing an order of a dedicated
// certificate.
func (sslCertificateApi *SslCertificateApiV1) OrderCertificate(orderCertificateOptions *OrderCertificateOptions) (result *DedicatedCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.OrderCertificateWithContext(context.Background(), orderCertificateOptions)
}

// OrderCertificateWithContext is an alternate form of the OrderCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) OrderCertificateWithContext(ctx context.Context, orderCertificateOptions *OrderCertificateOptions) (result *DedicatedCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(orderCertificateOptions, "orderCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/ssl/certificate_packs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range orderCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "OrderCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if orderCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*orderCertificateOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if orderCertificateOptions.Type != nil {
		body["type"] = orderCertificateOptions.Type
	}
	if orderCertificateOptions.Hosts != nil {
		body["hosts"] = orderCertificateOptions.Hosts
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDedicatedCertificateResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteCertificate : Delete a certificate
// Delete a given certificate.
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificate(deleteCertificateOptions *DeleteCertificateOptions) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.DeleteCertificateWithContext(context.Background(), deleteCertificateOptions)
}

// DeleteCertificateWithContext is an alternate form of the DeleteCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificateWithContext(ctx context.Context, deleteCertificateOptions *DeleteCertificateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCertificateOptions, "deleteCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCertificateOptions, "deleteCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"cert_identifier": *deleteCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/ssl/certificate_packs/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "DeleteCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCertificateOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = sslCertificateApi.Service.Request(request, nil)

	return
}

// GetSslSetting : Get SSL setting
// For a given zone identifier, get SSL setting.
func (sslCertificateApi *SslCertificateApiV1) GetSslSetting(getSslSettingOptions *GetSslSettingOptions) (result *SslSettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetSslSettingWithContext(context.Background(), getSslSettingOptions)
}

// GetSslSettingWithContext is an alternate form of the GetSslSetting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetSslSettingWithContext(ctx context.Context, getSslSettingOptions *GetSslSettingOptions) (result *SslSettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSslSettingOptions, "getSslSettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/ssl`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSslSettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetSslSetting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSslSettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ChangeSslSetting : Change SSL setting
// For a given zone identifier, change SSL setting.
func (sslCertificateApi *SslCertificateApiV1) ChangeSslSetting(changeSslSettingOptions *ChangeSslSettingOptions) (result *SslSettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ChangeSslSettingWithContext(context.Background(), changeSslSettingOptions)
}

// ChangeSslSettingWithContext is an alternate form of the ChangeSslSetting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ChangeSslSettingWithContext(ctx context.Context, changeSslSettingOptions *ChangeSslSettingOptions) (result *SslSettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(changeSslSettingOptions, "changeSslSettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/ssl`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changeSslSettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ChangeSslSetting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changeSslSettingOptions.Value != nil {
		body["value"] = changeSslSettingOptions.Value
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSslSettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListCustomCertificates : List all custom certificates
// For a given zone identifier, list all custom certificates.
func (sslCertificateApi *SslCertificateApiV1) ListCustomCertificates(listCustomCertificatesOptions *ListCustomCertificatesOptions) (result *ListCustomCertsResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ListCustomCertificatesWithContext(context.Background(), listCustomCertificatesOptions)
}

// ListCustomCertificatesWithContext is an alternate form of the ListCustomCertificates method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ListCustomCertificatesWithContext(ctx context.Context, listCustomCertificatesOptions *ListCustomCertificatesOptions) (result *ListCustomCertsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCustomCertificatesOptions, "listCustomCertificatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listCustomCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ListCustomCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCustomCertsResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UploadCustomCertificate : Upload a custom certificate
// For a given zone identifier, upload a custom certificates.
func (sslCertificateApi *SslCertificateApiV1) UploadCustomCertificate(uploadCustomCertificateOptions *UploadCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.UploadCustomCertificateWithContext(context.Background(), uploadCustomCertificateOptions)
}

// UploadCustomCertificateWithContext is an alternate form of the UploadCustomCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) UploadCustomCertificateWithContext(ctx context.Context, uploadCustomCertificateOptions *UploadCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(uploadCustomCertificateOptions, "uploadCustomCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadCustomCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "UploadCustomCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if uploadCustomCertificateOptions.Certificate != nil {
		body["certificate"] = uploadCustomCertificateOptions.Certificate
	}
	if uploadCustomCertificateOptions.PrivateKey != nil {
		body["private_key"] = uploadCustomCertificateOptions.PrivateKey
	}
	if uploadCustomCertificateOptions.BundleMethod != nil {
		body["bundle_method"] = uploadCustomCertificateOptions.BundleMethod
	}
	if uploadCustomCertificateOptions.GeoRestrictions != nil {
		body["geo_restrictions"] = uploadCustomCertificateOptions.GeoRestrictions
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetCustomCertificate : Get custom certificate
// For a given zone identifier, get a custom certificates.
func (sslCertificateApi *SslCertificateApiV1) GetCustomCertificate(getCustomCertificateOptions *GetCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetCustomCertificateWithContext(context.Background(), getCustomCertificateOptions)
}

// GetCustomCertificateWithContext is an alternate form of the GetCustomCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetCustomCertificateWithContext(ctx context.Context, getCustomCertificateOptions *GetCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCustomCertificateOptions, "getCustomCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCustomCertificateOptions, "getCustomCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id": *getCustomCertificateOptions.CustomCertID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates/{custom_cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCustomCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetCustomCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateCustomCertificate : Update custom certificate
// For a given zone identifier, update a custom certificates.
func (sslCertificateApi *SslCertificateApiV1) UpdateCustomCertificate(updateCustomCertificateOptions *UpdateCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.UpdateCustomCertificateWithContext(context.Background(), updateCustomCertificateOptions)
}

// UpdateCustomCertificateWithContext is an alternate form of the UpdateCustomCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) UpdateCustomCertificateWithContext(ctx context.Context, updateCustomCertificateOptions *UpdateCustomCertificateOptions) (result *CustomCertResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCustomCertificateOptions, "updateCustomCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateCustomCertificateOptions, "updateCustomCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id": *updateCustomCertificateOptions.CustomCertID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates/{custom_cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCustomCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "UpdateCustomCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCustomCertificateOptions.Certificate != nil {
		body["certificate"] = updateCustomCertificateOptions.Certificate
	}
	if updateCustomCertificateOptions.PrivateKey != nil {
		body["private_key"] = updateCustomCertificateOptions.PrivateKey
	}
	if updateCustomCertificateOptions.BundleMethod != nil {
		body["bundle_method"] = updateCustomCertificateOptions.BundleMethod
	}
	if updateCustomCertificateOptions.GeoRestrictions != nil {
		body["geo_restrictions"] = updateCustomCertificateOptions.GeoRestrictions
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteCustomCertificate : Delete custom certificate
// For a given zone identifier, delete a custom certificates.
func (sslCertificateApi *SslCertificateApiV1) DeleteCustomCertificate(deleteCustomCertificateOptions *DeleteCustomCertificateOptions) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.DeleteCustomCertificateWithContext(context.Background(), deleteCustomCertificateOptions)
}

// DeleteCustomCertificateWithContext is an alternate form of the DeleteCustomCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) DeleteCustomCertificateWithContext(ctx context.Context, deleteCustomCertificateOptions *DeleteCustomCertificateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCustomCertificateOptions, "deleteCustomCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCustomCertificateOptions, "deleteCustomCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id": *deleteCustomCertificateOptions.CustomCertID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates/{custom_cert_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCustomCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "DeleteCustomCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = sslCertificateApi.Service.Request(request, nil)

	return
}

// ChangeCertificatePriority : Set certificate priority
// For a given zone identifier, set priority of certificates.
func (sslCertificateApi *SslCertificateApiV1) ChangeCertificatePriority(changeCertificatePriorityOptions *ChangeCertificatePriorityOptions) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.ChangeCertificatePriorityWithContext(context.Background(), changeCertificatePriorityOptions)
}

// ChangeCertificatePriorityWithContext is an alternate form of the ChangeCertificatePriority method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ChangeCertificatePriorityWithContext(ctx context.Context, changeCertificatePriorityOptions *ChangeCertificatePriorityOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(changeCertificatePriorityOptions, "changeCertificatePriorityOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_certificates/prioritize`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changeCertificatePriorityOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ChangeCertificatePriority")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changeCertificatePriorityOptions.Certificates != nil {
		body["certificates"] = changeCertificatePriorityOptions.Certificates
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = sslCertificateApi.Service.Request(request, nil)

	return
}

// GetUniversalCertificateSetting : Get details of universal certificate
// For a given zone identifier, get universal certificate.
func (sslCertificateApi *SslCertificateApiV1) GetUniversalCertificateSetting(getUniversalCertificateSettingOptions *GetUniversalCertificateSettingOptions) (result *UniversalSettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetUniversalCertificateSettingWithContext(context.Background(), getUniversalCertificateSettingOptions)
}

// GetUniversalCertificateSettingWithContext is an alternate form of the GetUniversalCertificateSetting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetUniversalCertificateSettingWithContext(ctx context.Context, getUniversalCertificateSettingOptions *GetUniversalCertificateSettingOptions) (result *UniversalSettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getUniversalCertificateSettingOptions, "getUniversalCertificateSettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/ssl/universal/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUniversalCertificateSettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetUniversalCertificateSetting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUniversalSettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ChangeUniversalCertificateSetting : Enable or Disable universal certificate
// change universal certificate setting.
func (sslCertificateApi *SslCertificateApiV1) ChangeUniversalCertificateSetting(changeUniversalCertificateSettingOptions *ChangeUniversalCertificateSettingOptions) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.ChangeUniversalCertificateSettingWithContext(context.Background(), changeUniversalCertificateSettingOptions)
}

// ChangeUniversalCertificateSettingWithContext is an alternate form of the ChangeUniversalCertificateSetting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ChangeUniversalCertificateSettingWithContext(ctx context.Context, changeUniversalCertificateSettingOptions *ChangeUniversalCertificateSettingOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(changeUniversalCertificateSettingOptions, "changeUniversalCertificateSettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/ssl/universal/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changeUniversalCertificateSettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ChangeUniversalCertificateSetting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changeUniversalCertificateSettingOptions.Enabled != nil {
		body["enabled"] = changeUniversalCertificateSettingOptions.Enabled
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = sslCertificateApi.Service.Request(request, nil)

	return
}

// GetTls12Setting : Get TLS 1.2 only setting
// For a given zone identifier, get TLS 1.2 only setting.
func (sslCertificateApi *SslCertificateApiV1) GetTls12Setting(getTls12SettingOptions *GetTls12SettingOptions) (result *Tls12SettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetTls12SettingWithContext(context.Background(), getTls12SettingOptions)
}

// GetTls12SettingWithContext is an alternate form of the GetTls12Setting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetTls12SettingWithContext(ctx context.Context, getTls12SettingOptions *GetTls12SettingOptions) (result *Tls12SettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getTls12SettingOptions, "getTls12SettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/tls_1_2_only`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTls12SettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetTls12Setting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls12SettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ChangeTls12Setting : Set TLS 1.2 setting
// For a given zone identifier, set TLS 1.2 setting.
func (sslCertificateApi *SslCertificateApiV1) ChangeTls12Setting(changeTls12SettingOptions *ChangeTls12SettingOptions) (result *Tls12SettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ChangeTls12SettingWithContext(context.Background(), changeTls12SettingOptions)
}

// ChangeTls12SettingWithContext is an alternate form of the ChangeTls12Setting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ChangeTls12SettingWithContext(ctx context.Context, changeTls12SettingOptions *ChangeTls12SettingOptions) (result *Tls12SettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(changeTls12SettingOptions, "changeTls12SettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/tls_1_2_only`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changeTls12SettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ChangeTls12Setting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changeTls12SettingOptions.Value != nil {
		body["value"] = changeTls12SettingOptions.Value
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls12SettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetTls13Setting : Get TLS 1.3 setting
// For a given zone identifier, get TLS 1.3 setting.
func (sslCertificateApi *SslCertificateApiV1) GetTls13Setting(getTls13SettingOptions *GetTls13SettingOptions) (result *Tls13SettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetTls13SettingWithContext(context.Background(), getTls13SettingOptions)
}

// GetTls13SettingWithContext is an alternate form of the GetTls13Setting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetTls13SettingWithContext(ctx context.Context, getTls13SettingOptions *GetTls13SettingOptions) (result *Tls13SettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getTls13SettingOptions, "getTls13SettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/tls_1_3`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTls13SettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetTls13Setting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls13SettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ChangeTls13Setting : Set TLS 1.3 setting
// For a given zone identifier, set TLS 1.3 setting.
func (sslCertificateApi *SslCertificateApiV1) ChangeTls13Setting(changeTls13SettingOptions *ChangeTls13SettingOptions) (result *Tls13SettingResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ChangeTls13SettingWithContext(context.Background(), changeTls13SettingOptions)
}

// ChangeTls13SettingWithContext is an alternate form of the ChangeTls13Setting method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ChangeTls13SettingWithContext(ctx context.Context, changeTls13SettingOptions *ChangeTls13SettingOptions) (result *Tls13SettingResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(changeTls13SettingOptions, "changeTls13SettingOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/settings/tls_1_3`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changeTls13SettingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ChangeTls13Setting")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changeTls13SettingOptions.Value != nil {
		body["value"] = changeTls13SettingOptions.Value
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
	response, err = sslCertificateApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls13SettingResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CertPriorityReqCertificatesItem : certificate items.
type CertPriorityReqCertificatesItem struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// certificate priority.
	Priority *int64 `json:"priority" validate:"required"`
}


// NewCertPriorityReqCertificatesItem : Instantiate CertPriorityReqCertificatesItem (Generic Model Constructor)
func (*SslCertificateApiV1) NewCertPriorityReqCertificatesItem(id string, priority int64) (model *CertPriorityReqCertificatesItem, err error) {
	model = &CertPriorityReqCertificatesItem{
		ID: core.StringPtr(id),
		Priority: core.Int64Ptr(priority),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCertPriorityReqCertificatesItem unmarshals an instance of CertPriorityReqCertificatesItem from the specified map of raw messages.
func UnmarshalCertPriorityReqCertificatesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertPriorityReqCertificatesItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChangeCertificatePriorityOptions : The ChangeCertificatePriority options.
type ChangeCertificatePriorityOptions struct {
	// certificates array.
	Certificates []CertPriorityReqCertificatesItem `json:"certificates,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewChangeCertificatePriorityOptions : Instantiate ChangeCertificatePriorityOptions
func (*SslCertificateApiV1) NewChangeCertificatePriorityOptions() *ChangeCertificatePriorityOptions {
	return &ChangeCertificatePriorityOptions{}
}

// SetCertificates : Allow user to set Certificates
func (options *ChangeCertificatePriorityOptions) SetCertificates(certificates []CertPriorityReqCertificatesItem) *ChangeCertificatePriorityOptions {
	options.Certificates = certificates
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeCertificatePriorityOptions) SetHeaders(param map[string]string) *ChangeCertificatePriorityOptions {
	options.Headers = param
	return options
}

// ChangeSslSettingOptions : The ChangeSslSetting options.
type ChangeSslSettingOptions struct {
	// value.
	Value *string `json:"value,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ChangeSslSettingOptions.Value property.
// value.
const (
	ChangeSslSettingOptions_Value_Flexible = "flexible"
	ChangeSslSettingOptions_Value_Full = "full"
	ChangeSslSettingOptions_Value_Off = "off"
	ChangeSslSettingOptions_Value_Strict = "strict"
)

// NewChangeSslSettingOptions : Instantiate ChangeSslSettingOptions
func (*SslCertificateApiV1) NewChangeSslSettingOptions() *ChangeSslSettingOptions {
	return &ChangeSslSettingOptions{}
}

// SetValue : Allow user to set Value
func (options *ChangeSslSettingOptions) SetValue(value string) *ChangeSslSettingOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeSslSettingOptions) SetHeaders(param map[string]string) *ChangeSslSettingOptions {
	options.Headers = param
	return options
}

// ChangeTls12SettingOptions : The ChangeTls12Setting options.
type ChangeTls12SettingOptions struct {
	// value.
	Value *string `json:"value,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ChangeTls12SettingOptions.Value property.
// value.
const (
	ChangeTls12SettingOptions_Value_Off = "off"
	ChangeTls12SettingOptions_Value_On = "on"
)

// NewChangeTls12SettingOptions : Instantiate ChangeTls12SettingOptions
func (*SslCertificateApiV1) NewChangeTls12SettingOptions() *ChangeTls12SettingOptions {
	return &ChangeTls12SettingOptions{}
}

// SetValue : Allow user to set Value
func (options *ChangeTls12SettingOptions) SetValue(value string) *ChangeTls12SettingOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeTls12SettingOptions) SetHeaders(param map[string]string) *ChangeTls12SettingOptions {
	options.Headers = param
	return options
}

// ChangeTls13SettingOptions : The ChangeTls13Setting options.
type ChangeTls13SettingOptions struct {
	// value.
	Value *string `json:"value,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ChangeTls13SettingOptions.Value property.
// value.
const (
	ChangeTls13SettingOptions_Value_Off = "off"
	ChangeTls13SettingOptions_Value_On = "on"
)

// NewChangeTls13SettingOptions : Instantiate ChangeTls13SettingOptions
func (*SslCertificateApiV1) NewChangeTls13SettingOptions() *ChangeTls13SettingOptions {
	return &ChangeTls13SettingOptions{}
}

// SetValue : Allow user to set Value
func (options *ChangeTls13SettingOptions) SetValue(value string) *ChangeTls13SettingOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeTls13SettingOptions) SetHeaders(param map[string]string) *ChangeTls13SettingOptions {
	options.Headers = param
	return options
}

// ChangeUniversalCertificateSettingOptions : The ChangeUniversalCertificateSetting options.
type ChangeUniversalCertificateSettingOptions struct {
	// enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewChangeUniversalCertificateSettingOptions : Instantiate ChangeUniversalCertificateSettingOptions
func (*SslCertificateApiV1) NewChangeUniversalCertificateSettingOptions() *ChangeUniversalCertificateSettingOptions {
	return &ChangeUniversalCertificateSettingOptions{}
}

// SetEnabled : Allow user to set Enabled
func (options *ChangeUniversalCertificateSettingOptions) SetEnabled(enabled bool) *ChangeUniversalCertificateSettingOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeUniversalCertificateSettingOptions) SetHeaders(param map[string]string) *ChangeUniversalCertificateSettingOptions {
	options.Headers = param
	return options
}

// CustomCertReqGeoRestrictions : geo restrictions.
type CustomCertReqGeoRestrictions struct {
	// properties.
	Label *string `json:"label" validate:"required"`
}

// Constants associated with the CustomCertReqGeoRestrictions.Label property.
// properties.
const (
	CustomCertReqGeoRestrictions_Label_Eu = "eu"
	CustomCertReqGeoRestrictions_Label_HighestSecurity = "highest_security"
	CustomCertReqGeoRestrictions_Label_Us = "us"
)


// NewCustomCertReqGeoRestrictions : Instantiate CustomCertReqGeoRestrictions (Generic Model Constructor)
func (*SslCertificateApiV1) NewCustomCertReqGeoRestrictions(label string) (model *CustomCertReqGeoRestrictions, err error) {
	model = &CustomCertReqGeoRestrictions{
		Label: core.StringPtr(label),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCustomCertReqGeoRestrictions unmarshals an instance of CustomCertReqGeoRestrictions from the specified map of raw messages.
func UnmarshalCustomCertReqGeoRestrictions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomCertReqGeoRestrictions)
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteCertificateOptions : The DeleteCertificate options.
type DeleteCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCertificateOptions : Instantiate DeleteCertificateOptions
func (*SslCertificateApiV1) NewDeleteCertificateOptions(certIdentifier string) *DeleteCertificateOptions {
	return &DeleteCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (options *DeleteCertificateOptions) SetCertIdentifier(certIdentifier string) *DeleteCertificateOptions {
	options.CertIdentifier = core.StringPtr(certIdentifier)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteCertificateOptions) SetXCorrelationID(xCorrelationID string) *DeleteCertificateOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCertificateOptions) SetHeaders(param map[string]string) *DeleteCertificateOptions {
	options.Headers = param
	return options
}

// DeleteCustomCertificateOptions : The DeleteCustomCertificate options.
type DeleteCustomCertificateOptions struct {
	// custom certificate id.
	CustomCertID *string `json:"custom_cert_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCustomCertificateOptions : Instantiate DeleteCustomCertificateOptions
func (*SslCertificateApiV1) NewDeleteCustomCertificateOptions(customCertID string) *DeleteCustomCertificateOptions {
	return &DeleteCustomCertificateOptions{
		CustomCertID: core.StringPtr(customCertID),
	}
}

// SetCustomCertID : Allow user to set CustomCertID
func (options *DeleteCustomCertificateOptions) SetCustomCertID(customCertID string) *DeleteCustomCertificateOptions {
	options.CustomCertID = core.StringPtr(customCertID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCustomCertificateOptions) SetHeaders(param map[string]string) *DeleteCustomCertificateOptions {
	options.Headers = param
	return options
}

// GetCustomCertificateOptions : The GetCustomCertificate options.
type GetCustomCertificateOptions struct {
	// custom certificate id.
	CustomCertID *string `json:"custom_cert_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCustomCertificateOptions : Instantiate GetCustomCertificateOptions
func (*SslCertificateApiV1) NewGetCustomCertificateOptions(customCertID string) *GetCustomCertificateOptions {
	return &GetCustomCertificateOptions{
		CustomCertID: core.StringPtr(customCertID),
	}
}

// SetCustomCertID : Allow user to set CustomCertID
func (options *GetCustomCertificateOptions) SetCustomCertID(customCertID string) *GetCustomCertificateOptions {
	options.CustomCertID = core.StringPtr(customCertID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetCustomCertificateOptions) SetHeaders(param map[string]string) *GetCustomCertificateOptions {
	options.Headers = param
	return options
}

// GetSslSettingOptions : The GetSslSetting options.
type GetSslSettingOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSslSettingOptions : Instantiate GetSslSettingOptions
func (*SslCertificateApiV1) NewGetSslSettingOptions() *GetSslSettingOptions {
	return &GetSslSettingOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSslSettingOptions) SetHeaders(param map[string]string) *GetSslSettingOptions {
	options.Headers = param
	return options
}

// GetTls12SettingOptions : The GetTls12Setting options.
type GetTls12SettingOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTls12SettingOptions : Instantiate GetTls12SettingOptions
func (*SslCertificateApiV1) NewGetTls12SettingOptions() *GetTls12SettingOptions {
	return &GetTls12SettingOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetTls12SettingOptions) SetHeaders(param map[string]string) *GetTls12SettingOptions {
	options.Headers = param
	return options
}

// GetTls13SettingOptions : The GetTls13Setting options.
type GetTls13SettingOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTls13SettingOptions : Instantiate GetTls13SettingOptions
func (*SslCertificateApiV1) NewGetTls13SettingOptions() *GetTls13SettingOptions {
	return &GetTls13SettingOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetTls13SettingOptions) SetHeaders(param map[string]string) *GetTls13SettingOptions {
	options.Headers = param
	return options
}

// GetUniversalCertificateSettingOptions : The GetUniversalCertificateSetting options.
type GetUniversalCertificateSettingOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUniversalCertificateSettingOptions : Instantiate GetUniversalCertificateSettingOptions
func (*SslCertificateApiV1) NewGetUniversalCertificateSettingOptions() *GetUniversalCertificateSettingOptions {
	return &GetUniversalCertificateSettingOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetUniversalCertificateSettingOptions) SetHeaders(param map[string]string) *GetUniversalCertificateSettingOptions {
	options.Headers = param
	return options
}

// ListCertificatesOptions : The ListCertificates options.
type ListCertificatesOptions struct {
	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCertificatesOptions : Instantiate ListCertificatesOptions
func (*SslCertificateApiV1) NewListCertificatesOptions() *ListCertificatesOptions {
	return &ListCertificatesOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListCertificatesOptions) SetXCorrelationID(xCorrelationID string) *ListCertificatesOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListCertificatesOptions) SetHeaders(param map[string]string) *ListCertificatesOptions {
	options.Headers = param
	return options
}

// ListCustomCertificatesOptions : The ListCustomCertificates options.
type ListCustomCertificatesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListCustomCertificatesOptions : Instantiate ListCustomCertificatesOptions
func (*SslCertificateApiV1) NewListCustomCertificatesOptions() *ListCustomCertificatesOptions {
	return &ListCustomCertificatesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListCustomCertificatesOptions) SetHeaders(param map[string]string) *ListCustomCertificatesOptions {
	options.Headers = param
	return options
}

// OrderCertificateOptions : The OrderCertificate options.
type OrderCertificateOptions struct {
	// priorities.
	Type *string `json:"type,omitempty"`

	// host name.
	Hosts []string `json:"hosts,omitempty"`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the OrderCertificateOptions.Type property.
// priorities.
const (
	OrderCertificateOptions_Type_Dedicated = "dedicated"
)

// NewOrderCertificateOptions : Instantiate OrderCertificateOptions
func (*SslCertificateApiV1) NewOrderCertificateOptions() *OrderCertificateOptions {
	return &OrderCertificateOptions{}
}

// SetType : Allow user to set Type
func (options *OrderCertificateOptions) SetType(typeVar string) *OrderCertificateOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHosts : Allow user to set Hosts
func (options *OrderCertificateOptions) SetHosts(hosts []string) *OrderCertificateOptions {
	options.Hosts = hosts
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *OrderCertificateOptions) SetXCorrelationID(xCorrelationID string) *OrderCertificateOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *OrderCertificateOptions) SetHeaders(param map[string]string) *OrderCertificateOptions {
	options.Headers = param
	return options
}

// Tls12SettingRespMessagesItem : Tls12SettingRespMessagesItem struct
type Tls12SettingRespMessagesItem struct {
	// status.
	Status *string `json:"status,omitempty"`
}


// UnmarshalTls12SettingRespMessagesItem unmarshals an instance of Tls12SettingRespMessagesItem from the specified map of raw messages.
func UnmarshalTls12SettingRespMessagesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tls12SettingRespMessagesItem)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tls12SettingRespResult : result.
type Tls12SettingRespResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value" validate:"required"`

	// editable.
	Editable *bool `json:"editable" validate:"required"`

	// modified date.
	ModifiedOn *strfmt.DateTime `json:"modified_on" validate:"required"`
}

// Constants associated with the Tls12SettingRespResult.ID property.
// identifier.
const (
	Tls12SettingRespResult_ID_Tls12Only = "tls_1_2_only"
)


// UnmarshalTls12SettingRespResult unmarshals an instance of Tls12SettingRespResult from the specified map of raw messages.
func UnmarshalTls12SettingRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tls12SettingRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
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

// Tls13SettingRespResult : result.
type Tls13SettingRespResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value" validate:"required"`

	// editable.
	Editable *bool `json:"editable" validate:"required"`

	// modified date.
	ModifiedOn *strfmt.DateTime `json:"modified_on" validate:"required"`
}

// Constants associated with the Tls13SettingRespResult.ID property.
// identifier.
const (
	Tls13SettingRespResult_ID_Tls13 = "tls_1_3"
)


// UnmarshalTls13SettingRespResult unmarshals an instance of Tls13SettingRespResult from the specified map of raw messages.
func UnmarshalTls13SettingRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tls13SettingRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
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

// UniversalSettingRespResult : result.
type UniversalSettingRespResult struct {
	// enabled.
	Enabled *bool `json:"enabled" validate:"required"`
}


// UnmarshalUniversalSettingRespResult unmarshals an instance of UniversalSettingRespResult from the specified map of raw messages.
func UnmarshalUniversalSettingRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UniversalSettingRespResult)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCustomCertificateOptions : The UpdateCustomCertificate options.
type UpdateCustomCertificateOptions struct {
	// custom certificate id.
	CustomCertID *string `json:"custom_cert_id" validate:"required,ne="`

	// certificates.
	Certificate *string `json:"certificate,omitempty"`

	// private key.
	PrivateKey *string `json:"private_key,omitempty"`

	// Methods shown in UI mapping to API: Compatible(ubiquitous), Modern(optimal), User Defined(force).
	BundleMethod *string `json:"bundle_method,omitempty"`

	// geo restrictions.
	GeoRestrictions *CustomCertReqGeoRestrictions `json:"geo_restrictions,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateCustomCertificateOptions.BundleMethod property.
// Methods shown in UI mapping to API: Compatible(ubiquitous), Modern(optimal), User Defined(force).
const (
	UpdateCustomCertificateOptions_BundleMethod_Force = "force"
	UpdateCustomCertificateOptions_BundleMethod_Optimal = "optimal"
	UpdateCustomCertificateOptions_BundleMethod_Ubiquitous = "ubiquitous"
)

// NewUpdateCustomCertificateOptions : Instantiate UpdateCustomCertificateOptions
func (*SslCertificateApiV1) NewUpdateCustomCertificateOptions(customCertID string) *UpdateCustomCertificateOptions {
	return &UpdateCustomCertificateOptions{
		CustomCertID: core.StringPtr(customCertID),
	}
}

// SetCustomCertID : Allow user to set CustomCertID
func (options *UpdateCustomCertificateOptions) SetCustomCertID(customCertID string) *UpdateCustomCertificateOptions {
	options.CustomCertID = core.StringPtr(customCertID)
	return options
}

// SetCertificate : Allow user to set Certificate
func (options *UpdateCustomCertificateOptions) SetCertificate(certificate string) *UpdateCustomCertificateOptions {
	options.Certificate = core.StringPtr(certificate)
	return options
}

// SetPrivateKey : Allow user to set PrivateKey
func (options *UpdateCustomCertificateOptions) SetPrivateKey(privateKey string) *UpdateCustomCertificateOptions {
	options.PrivateKey = core.StringPtr(privateKey)
	return options
}

// SetBundleMethod : Allow user to set BundleMethod
func (options *UpdateCustomCertificateOptions) SetBundleMethod(bundleMethod string) *UpdateCustomCertificateOptions {
	options.BundleMethod = core.StringPtr(bundleMethod)
	return options
}

// SetGeoRestrictions : Allow user to set GeoRestrictions
func (options *UpdateCustomCertificateOptions) SetGeoRestrictions(geoRestrictions *CustomCertReqGeoRestrictions) *UpdateCustomCertificateOptions {
	options.GeoRestrictions = geoRestrictions
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCustomCertificateOptions) SetHeaders(param map[string]string) *UpdateCustomCertificateOptions {
	options.Headers = param
	return options
}

// UploadCustomCertificateOptions : The UploadCustomCertificate options.
type UploadCustomCertificateOptions struct {
	// certificates.
	Certificate *string `json:"certificate,omitempty"`

	// private key.
	PrivateKey *string `json:"private_key,omitempty"`

	// Methods shown in UI mapping to API: Compatible(ubiquitous), Modern(optimal), User Defined(force).
	BundleMethod *string `json:"bundle_method,omitempty"`

	// geo restrictions.
	GeoRestrictions *CustomCertReqGeoRestrictions `json:"geo_restrictions,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UploadCustomCertificateOptions.BundleMethod property.
// Methods shown in UI mapping to API: Compatible(ubiquitous), Modern(optimal), User Defined(force).
const (
	UploadCustomCertificateOptions_BundleMethod_Force = "force"
	UploadCustomCertificateOptions_BundleMethod_Optimal = "optimal"
	UploadCustomCertificateOptions_BundleMethod_Ubiquitous = "ubiquitous"
)

// NewUploadCustomCertificateOptions : Instantiate UploadCustomCertificateOptions
func (*SslCertificateApiV1) NewUploadCustomCertificateOptions() *UploadCustomCertificateOptions {
	return &UploadCustomCertificateOptions{}
}

// SetCertificate : Allow user to set Certificate
func (options *UploadCustomCertificateOptions) SetCertificate(certificate string) *UploadCustomCertificateOptions {
	options.Certificate = core.StringPtr(certificate)
	return options
}

// SetPrivateKey : Allow user to set PrivateKey
func (options *UploadCustomCertificateOptions) SetPrivateKey(privateKey string) *UploadCustomCertificateOptions {
	options.PrivateKey = core.StringPtr(privateKey)
	return options
}

// SetBundleMethod : Allow user to set BundleMethod
func (options *UploadCustomCertificateOptions) SetBundleMethod(bundleMethod string) *UploadCustomCertificateOptions {
	options.BundleMethod = core.StringPtr(bundleMethod)
	return options
}

// SetGeoRestrictions : Allow user to set GeoRestrictions
func (options *UploadCustomCertificateOptions) SetGeoRestrictions(geoRestrictions *CustomCertReqGeoRestrictions) *UploadCustomCertificateOptions {
	options.GeoRestrictions = geoRestrictions
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UploadCustomCertificateOptions) SetHeaders(param map[string]string) *UploadCustomCertificateOptions {
	options.Headers = param
	return options
}

// Certificate : certificate.
type Certificate struct {
	// identifier.
	ID interface{} `json:"id" validate:"required"`

	// host name.
	Hosts []string `json:"hosts" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`
}


// UnmarshalCertificate unmarshals an instance of Certificate from the specified map of raw messages.
func UnmarshalCertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Certificate)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hosts", &obj.Hosts)
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

// CustomCertPack : custom certificate pack.
type CustomCertPack struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// host name.
	Hosts []string `json:"hosts" validate:"required"`

	// issuer.
	Issuer *string `json:"issuer" validate:"required"`

	// signature.
	Signature *string `json:"signature" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`

	// bundle method.
	BundleMethod *string `json:"bundle_method" validate:"required"`

	// zone identifier.
	ZoneID *string `json:"zone_id" validate:"required"`

	// uploaded date.
	UploadedOn *string `json:"uploaded_on" validate:"required"`

	// modified date.
	ModifiedOn *string `json:"modified_on" validate:"required"`

	// expire date.
	ExpiresOn *string `json:"expires_on" validate:"required"`

	// priority.
	Priority *float64 `json:"priority" validate:"required"`
}


// UnmarshalCustomCertPack unmarshals an instance of CustomCertPack from the specified map of raw messages.
func UnmarshalCustomCertPack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomCertPack)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hosts", &obj.Hosts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signature", &obj.Signature)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bundle_method", &obj.BundleMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_id", &obj.ZoneID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uploaded_on", &obj.UploadedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CustomCertResp : custom certificate response.
type CustomCertResp struct {
	// custom certificate pack.
	Result *CustomCertPack `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalCustomCertResp unmarshals an instance of CustomCertResp from the specified map of raw messages.
func UnmarshalCustomCertResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomCertResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCustomCertPack)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DedicatedCertificatePack : dedicated certificate packs.
type DedicatedCertificatePack struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// certificate type.
	Type *string `json:"type" validate:"required"`

	// host name.
	Hosts []string `json:"hosts" validate:"required"`

	// certificates.
	Certificates []Certificate `json:"certificates" validate:"required"`

	// primary certificate.
	PrimaryCertificate interface{} `json:"primary_certificate" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`
}


// UnmarshalDedicatedCertificatePack unmarshals an instance of DedicatedCertificatePack from the specified map of raw messages.
func UnmarshalDedicatedCertificatePack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DedicatedCertificatePack)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hosts", &obj.Hosts)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certificates", &obj.Certificates, UnmarshalCertificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "primary_certificate", &obj.PrimaryCertificate)
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

// DedicatedCertificateResp : certificate response.
type DedicatedCertificateResp struct {
	// dedicated certificate packs.
	Result *DedicatedCertificatePack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalDedicatedCertificateResp unmarshals an instance of DedicatedCertificateResp from the specified map of raw messages.
func UnmarshalDedicatedCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DedicatedCertificateResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDedicatedCertificatePack)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListCertificateResp : certificate response.
type ListCertificateResp struct {
	// certificate packs.
	Result []DedicatedCertificatePack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalListCertificateResp unmarshals an instance of ListCertificateResp from the specified map of raw messages.
func UnmarshalListCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListCertificateResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDedicatedCertificatePack)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListCustomCertsResp : custom certificate response.
type ListCustomCertsResp struct {
	// custom certificate packs.
	Result []CustomCertPack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalListCustomCertsResp unmarshals an instance of ListCustomCertsResp from the specified map of raw messages.
func UnmarshalListCustomCertsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListCustomCertsResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCustomCertPack)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResultInfo : result information.
type ResultInfo struct {
	// page number.
	Page *int64 `json:"page" validate:"required"`

	// per page count.
	PerPage *int64 `json:"per_page" validate:"required"`

	// count.
	Count *int64 `json:"count" validate:"required"`

	// total count.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalResultInfo unmarshals an instance of ResultInfo from the specified map of raw messages.
func UnmarshalResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResultInfo)
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

// SslSetting : ssl setting.
type SslSetting struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value" validate:"required"`

	// editable.
	Editable *bool `json:"editable" validate:"required"`

	// modified date.
	ModifiedOn *string `json:"modified_on" validate:"required"`
}


// UnmarshalSslSetting unmarshals an instance of SslSetting from the specified map of raw messages.
func UnmarshalSslSetting(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SslSetting)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
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

// SslSettingResp : ssl setting response.
type SslSettingResp struct {
	// success.
	Success *bool `json:"success" validate:"required"`

	// ssl setting.
	Result *SslSetting `json:"result" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalSslSettingResp unmarshals an instance of SslSettingResp from the specified map of raw messages.
func UnmarshalSslSettingResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SslSettingResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalSslSetting)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tls12SettingResp : tls 1.2 setting response.
type Tls12SettingResp struct {
	// result.
	Result *Tls12SettingRespResult `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalTls12SettingResp unmarshals an instance of Tls12SettingResp from the specified map of raw messages.
func UnmarshalTls12SettingResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tls12SettingResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalTls12SettingRespResult)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tls13SettingResp : tls 1.3 setting response.
type Tls13SettingResp struct {
	// result.
	Result *Tls13SettingRespResult `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalTls13SettingResp unmarshals an instance of Tls13SettingResp from the specified map of raw messages.
func UnmarshalTls13SettingResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tls13SettingResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalTls13SettingRespResult)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UniversalSettingResp : universal setting response.
type UniversalSettingResp struct {
	// result.
	Result *UniversalSettingRespResult `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages []Tls12SettingRespMessagesItem `json:"messages" validate:"required"`
}


// UnmarshalUniversalSettingResp unmarshals an instance of UniversalSettingResp from the specified map of raw messages.
func UnmarshalUniversalSettingResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UniversalSettingResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalUniversalSettingRespResult)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "messages", &obj.Messages, UnmarshalTls12SettingRespMessagesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
