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

// Package sslcertificateapiv1 : Operations and models for the SslCertificateApiV1 service
package sslcertificateapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// SslCertificateApiV1 : SSL Certificate
//
// API Version: 1.0.0
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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// OrderCertificate : Order dedicated certificate
// Order a dedicated certificate for a given zone. The zone should be active before placing an order of a dedicated
// certificate. Deprecated, please use Advanced Certificate Pack.
// Deprecated: this method is deprecated and may be removed in a future release.
func (sslCertificateApi *SslCertificateApiV1) OrderCertificate(orderCertificateOptions *OrderCertificateOptions) (result *DedicatedCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.OrderCertificateWithContext(context.Background(), orderCertificateOptions)
}

// OrderCertificateWithContext is an alternate form of the OrderCertificate method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (sslCertificateApi *SslCertificateApiV1) OrderCertificateWithContext(ctx context.Context, orderCertificateOptions *OrderCertificateOptions) (result *DedicatedCertificateResp, response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: OrderCertificate")
	err = core.ValidateStruct(orderCertificateOptions, "orderCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDedicatedCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCertificate : Delete a certificate
// Delete a given certificate. Deprecated, please use Advanced Certificate Pack.
// Deprecated: this method is deprecated and may be removed in a future release.
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificate(deleteCertificateOptions *DeleteCertificateOptions) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.DeleteCertificateWithContext(context.Background(), deleteCertificateOptions)
}

// DeleteCertificateWithContext is an alternate form of the DeleteCertificate method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificateWithContext(ctx context.Context, deleteCertificateOptions *DeleteCertificateOptions) (response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: DeleteCertificate")
	err = core.ValidateNotNil(deleteCertificateOptions, "deleteCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCertificateOptions, "deleteCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSslSettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSslSettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCustomCertsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id":  *getCustomCertificateOptions.CustomCertID,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id":  *updateCustomCertificateOptions.CustomCertID,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomCertResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"custom_cert_id":  *deleteCustomCertificateOptions.CustomCertID,
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
		"crn":             *sslCertificateApi.Crn,
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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUniversalSettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls12SettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls12SettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls13SettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

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
		"crn":             *sslCertificateApi.Crn,
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTls13SettingResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// OrderAdvancedCertificate : Order advanced certificate
// Order an advanced certificate pack for a given zone. The zone should be active before ordering of an advanced
// certificate pack.
func (sslCertificateApi *SslCertificateApiV1) OrderAdvancedCertificate(orderAdvancedCertificateOptions *OrderAdvancedCertificateOptions) (result *AdvancedCertInitResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.OrderAdvancedCertificateWithContext(context.Background(), orderAdvancedCertificateOptions)
}

// OrderAdvancedCertificateWithContext is an alternate form of the OrderAdvancedCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) OrderAdvancedCertificateWithContext(ctx context.Context, orderAdvancedCertificateOptions *OrderAdvancedCertificateOptions) (result *AdvancedCertInitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(orderAdvancedCertificateOptions, "orderAdvancedCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v2/{crn}/zones/{zone_identifier}/ssl/certificate_packs/order`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range orderAdvancedCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "OrderAdvancedCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if orderAdvancedCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*orderAdvancedCertificateOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if orderAdvancedCertificateOptions.Type != nil {
		body["type"] = orderAdvancedCertificateOptions.Type
	}
	if orderAdvancedCertificateOptions.Hosts != nil {
		body["hosts"] = orderAdvancedCertificateOptions.Hosts
	}
	if orderAdvancedCertificateOptions.ValidationMethod != nil {
		body["validation_method"] = orderAdvancedCertificateOptions.ValidationMethod
	}
	if orderAdvancedCertificateOptions.ValidityDays != nil {
		body["validity_days"] = orderAdvancedCertificateOptions.ValidityDays
	}
	if orderAdvancedCertificateOptions.CertificateAuthority != nil {
		body["certificate_authority"] = orderAdvancedCertificateOptions.CertificateAuthority
	}
	if orderAdvancedCertificateOptions.CloudflareBranding != nil {
		body["cloudflare_branding"] = orderAdvancedCertificateOptions.CloudflareBranding
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAdvancedCertInitResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PatchCertificate : Restart validation for an advanced certificate pack
// Restart validation for an advanced certificate pack. This is only a validation operation for a Certificate Pack in a
// validation_timed_out status.
func (sslCertificateApi *SslCertificateApiV1) PatchCertificate(patchCertificateOptions *PatchCertificateOptions) (result *AdvancedCertInitResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.PatchCertificateWithContext(context.Background(), patchCertificateOptions)
}

// PatchCertificateWithContext is an alternate form of the PatchCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) PatchCertificateWithContext(ctx context.Context, patchCertificateOptions *PatchCertificateOptions) (result *AdvancedCertInitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(patchCertificateOptions, "patchCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(patchCertificateOptions, "patchCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"cert_identifier": *patchCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v2/{crn}/zones/{zone_identifier}/ssl/certificate_packs/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range patchCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "PatchCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if patchCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*patchCertificateOptions.XCorrelationID))
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAdvancedCertInitResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteCertificateV2 : Delete a certificate
// Delete an advanced certificate pack.
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificateV2(deleteCertificateV2Options *DeleteCertificateV2Options) (response *core.DetailedResponse, err error) {
	return sslCertificateApi.DeleteCertificateV2WithContext(context.Background(), deleteCertificateV2Options)
}

// DeleteCertificateV2WithContext is an alternate form of the DeleteCertificateV2 method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) DeleteCertificateV2WithContext(ctx context.Context, deleteCertificateV2Options *DeleteCertificateV2Options) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCertificateV2Options, "deleteCertificateV2Options cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCertificateV2Options, "deleteCertificateV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
		"cert_identifier": *deleteCertificateV2Options.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v2/{crn}/zones/{zone_identifier}/ssl/certificate_packs/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteCertificateV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "DeleteCertificateV2")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteCertificateV2Options.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteCertificateV2Options.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = sslCertificateApi.Service.Request(request, nil)

	return
}

// GetSslVerification : Get SSL Verification Info for a Zone
// Get SSL Verification Info for a Zone.
func (sslCertificateApi *SslCertificateApiV1) GetSslVerification(getSslVerificationOptions *GetSslVerificationOptions) (result *SslVerificationResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetSslVerificationWithContext(context.Background(), getSslVerificationOptions)
}

// GetSslVerificationWithContext is an alternate form of the GetSslVerification method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetSslVerificationWithContext(ctx context.Context, getSslVerificationOptions *GetSslVerificationOptions) (result *SslVerificationResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSslVerificationOptions, "getSslVerificationOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *sslCertificateApi.Crn,
		"zone_identifier": *sslCertificateApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v2/{crn}/zones/{zone_identifier}/ssl/verification`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSslVerificationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetSslVerification")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getSslVerificationOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getSslVerificationOptions.XCorrelationID))
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSslVerificationResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListOriginCertificates : List all Origin Certificates
// List all existing CIS-issued Certificates for a given domain.
func (sslCertificateApi *SslCertificateApiV1) ListOriginCertificates(listOriginCertificatesOptions *ListOriginCertificatesOptions) (result *ListOriginCertificatesResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.ListOriginCertificatesWithContext(context.Background(), listOriginCertificatesOptions)
}

// ListOriginCertificatesWithContext is an alternate form of the ListOriginCertificates method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) ListOriginCertificatesWithContext(ctx context.Context, listOriginCertificatesOptions *ListOriginCertificatesOptions) (result *ListOriginCertificatesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listOriginCertificatesOptions, "listOriginCertificatesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listOriginCertificatesOptions, "listOriginCertificatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *listOriginCertificatesOptions.Crn,
		"zone_identifier": *listOriginCertificatesOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listOriginCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "ListOriginCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listOriginCertificatesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listOriginCertificatesOptions.XCorrelationID))
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListOriginCertificatesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateOriginCertificate : Create a CIS-signed certificate
// Create a CIS-signed certificate.
func (sslCertificateApi *SslCertificateApiV1) CreateOriginCertificate(createOriginCertificateOptions *CreateOriginCertificateOptions) (result *OriginCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.CreateOriginCertificateWithContext(context.Background(), createOriginCertificateOptions)
}

// CreateOriginCertificateWithContext is an alternate form of the CreateOriginCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) CreateOriginCertificateWithContext(ctx context.Context, createOriginCertificateOptions *CreateOriginCertificateOptions) (result *OriginCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createOriginCertificateOptions, "createOriginCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createOriginCertificateOptions, "createOriginCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *createOriginCertificateOptions.Crn,
		"zone_identifier": *createOriginCertificateOptions.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createOriginCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "CreateOriginCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createOriginCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createOriginCertificateOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createOriginCertificateOptions.Hostnames != nil {
		body["hostnames"] = createOriginCertificateOptions.Hostnames
	}
	if createOriginCertificateOptions.RequestType != nil {
		body["request_type"] = createOriginCertificateOptions.RequestType
	}
	if createOriginCertificateOptions.RequestedValidity != nil {
		body["requested_validity"] = createOriginCertificateOptions.RequestedValidity
	}
	if createOriginCertificateOptions.Csr != nil {
		body["csr"] = createOriginCertificateOptions.Csr
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOriginCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// RevokeOriginCertificate : Revoke a created Origin Certificate for a domain
// Revoke a created Origin certificate.
func (sslCertificateApi *SslCertificateApiV1) RevokeOriginCertificate(revokeOriginCertificateOptions *RevokeOriginCertificateOptions) (result *RevokeOriginCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.RevokeOriginCertificateWithContext(context.Background(), revokeOriginCertificateOptions)
}

// RevokeOriginCertificateWithContext is an alternate form of the RevokeOriginCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) RevokeOriginCertificateWithContext(ctx context.Context, revokeOriginCertificateOptions *RevokeOriginCertificateOptions) (result *RevokeOriginCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(revokeOriginCertificateOptions, "revokeOriginCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(revokeOriginCertificateOptions, "revokeOriginCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *revokeOriginCertificateOptions.Crn,
		"zone_identifier": *revokeOriginCertificateOptions.ZoneIdentifier,
		"cert_identifier": *revokeOriginCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_certificates/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range revokeOriginCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "RevokeOriginCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if revokeOriginCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*revokeOriginCertificateOptions.XCorrelationID))
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRevokeOriginCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetOriginCertificate : Get an existing Origin certificate
// Get an existing Origin certificate by its serial number.
func (sslCertificateApi *SslCertificateApiV1) GetOriginCertificate(getOriginCertificateOptions *GetOriginCertificateOptions) (result *OriginCertificateResp, response *core.DetailedResponse, err error) {
	return sslCertificateApi.GetOriginCertificateWithContext(context.Background(), getOriginCertificateOptions)
}

// GetOriginCertificateWithContext is an alternate form of the GetOriginCertificate method which supports a Context parameter
func (sslCertificateApi *SslCertificateApiV1) GetOriginCertificateWithContext(ctx context.Context, getOriginCertificateOptions *GetOriginCertificateOptions) (result *OriginCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOriginCertificateOptions, "getOriginCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getOriginCertificateOptions, "getOriginCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *getOriginCertificateOptions.Crn,
		"zone_identifier": *getOriginCertificateOptions.ZoneIdentifier,
		"cert_identifier": *getOriginCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = sslCertificateApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(sslCertificateApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_certificates/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getOriginCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("ssl_certificate_api", "V1", "GetOriginCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getOriginCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getOriginCertificateOptions.XCorrelationID))
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
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOriginCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// AdvancedCertInitRespResult : result of ordering advanced certificate pack.
type AdvancedCertInitRespResult struct {
	// advanced certificate pack ID.
	ID *string `json:"id,omitempty"`

	// certificate type.
	Type *string `json:"type,omitempty"`

	// host name.
	Hosts []string `json:"hosts,omitempty"`

	// validation Method selected for the order.
	ValidationMethod *string `json:"validation_method,omitempty"`

	// validity Days selected for the order.
	ValidityDays *int64 `json:"validity_days,omitempty"`

	// Certificate Authority selected for the order.
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// whether or not to add Cloudflare Branding for the order.
	CloudflareBranding *bool `json:"cloudflare_branding,omitempty"`

	// certificate status.
	Status *string `json:"status,omitempty"`
}

// UnmarshalAdvancedCertInitRespResult unmarshals an instance of AdvancedCertInitRespResult from the specified map of raw messages.
func UnmarshalAdvancedCertInitRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdvancedCertInitRespResult)
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
	err = core.UnmarshalPrimitive(m, "validation_method", &obj.ValidationMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "validity_days", &obj.ValidityDays)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_authority", &obj.CertificateAuthority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cloudflare_branding", &obj.CloudflareBranding)
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

// CertPriorityReqCertificatesItem : certificate items.
type CertPriorityReqCertificatesItem struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// certificate priority.
	Priority *int64 `json:"priority" validate:"required"`
}

// NewCertPriorityReqCertificatesItem : Instantiate CertPriorityReqCertificatesItem (Generic Model Constructor)
func (*SslCertificateApiV1) NewCertPriorityReqCertificatesItem(id string, priority int64) (_model *CertPriorityReqCertificatesItem, err error) {
	_model = &CertPriorityReqCertificatesItem{
		ID:       core.StringPtr(id),
		Priority: core.Int64Ptr(priority),
	}
	err = core.ValidateStruct(_model, "required parameters")
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
func (_options *ChangeCertificatePriorityOptions) SetCertificates(certificates []CertPriorityReqCertificatesItem) *ChangeCertificatePriorityOptions {
	_options.Certificates = certificates
	return _options
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
	ChangeSslSettingOptions_Value_Full     = "full"
	ChangeSslSettingOptions_Value_Off      = "off"
	ChangeSslSettingOptions_Value_Strict   = "strict"
)

// NewChangeSslSettingOptions : Instantiate ChangeSslSettingOptions
func (*SslCertificateApiV1) NewChangeSslSettingOptions() *ChangeSslSettingOptions {
	return &ChangeSslSettingOptions{}
}

// SetValue : Allow user to set Value
func (_options *ChangeSslSettingOptions) SetValue(value string) *ChangeSslSettingOptions {
	_options.Value = core.StringPtr(value)
	return _options
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
	ChangeTls12SettingOptions_Value_On  = "on"
)

// NewChangeTls12SettingOptions : Instantiate ChangeTls12SettingOptions
func (*SslCertificateApiV1) NewChangeTls12SettingOptions() *ChangeTls12SettingOptions {
	return &ChangeTls12SettingOptions{}
}

// SetValue : Allow user to set Value
func (_options *ChangeTls12SettingOptions) SetValue(value string) *ChangeTls12SettingOptions {
	_options.Value = core.StringPtr(value)
	return _options
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
	ChangeTls13SettingOptions_Value_On  = "on"
)

// NewChangeTls13SettingOptions : Instantiate ChangeTls13SettingOptions
func (*SslCertificateApiV1) NewChangeTls13SettingOptions() *ChangeTls13SettingOptions {
	return &ChangeTls13SettingOptions{}
}

// SetValue : Allow user to set Value
func (_options *ChangeTls13SettingOptions) SetValue(value string) *ChangeTls13SettingOptions {
	_options.Value = core.StringPtr(value)
	return _options
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
func (_options *ChangeUniversalCertificateSettingOptions) SetEnabled(enabled bool) *ChangeUniversalCertificateSettingOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ChangeUniversalCertificateSettingOptions) SetHeaders(param map[string]string) *ChangeUniversalCertificateSettingOptions {
	options.Headers = param
	return options
}

// CreateOriginCertificateOptions : The CreateOriginCertificate options.
type CreateOriginCertificateOptions struct {
	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `json:"crn" validate:"required,ne="`

	// zone identifier.
	ZoneIdentifier *string `json:"zone_identifier" validate:"required,ne="`

	// Array of hostnames or wildcard names (e.g., *.example.com) bound to the certificate.
	Hostnames []string `json:"hostnames,omitempty"`

	// Signature type desired on certificate.
	RequestType *string `json:"request_type,omitempty"`

	// The number of days for which the certificate should be valid.
	RequestedValidity *int64 `json:"requested_validity,omitempty"`

	// The Certificate Signing Request (CSR).
	Csr *string `json:"csr,omitempty"`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateOriginCertificateOptions.RequestType property.
// Signature type desired on certificate.
const (
	CreateOriginCertificateOptions_RequestType_OriginEcc = "origin-ecc"
	CreateOriginCertificateOptions_RequestType_OriginRsa = "origin-rsa"
)

// NewCreateOriginCertificateOptions : Instantiate CreateOriginCertificateOptions
func (*SslCertificateApiV1) NewCreateOriginCertificateOptions(crn string, zoneIdentifier string) *CreateOriginCertificateOptions {
	return &CreateOriginCertificateOptions{
		Crn:            core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
	}
}

// SetCrn : Allow user to set Crn
func (_options *CreateOriginCertificateOptions) SetCrn(crn string) *CreateOriginCertificateOptions {
	_options.Crn = core.StringPtr(crn)
	return _options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (_options *CreateOriginCertificateOptions) SetZoneIdentifier(zoneIdentifier string) *CreateOriginCertificateOptions {
	_options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return _options
}

// SetHostnames : Allow user to set Hostnames
func (_options *CreateOriginCertificateOptions) SetHostnames(hostnames []string) *CreateOriginCertificateOptions {
	_options.Hostnames = hostnames
	return _options
}

// SetRequestType : Allow user to set RequestType
func (_options *CreateOriginCertificateOptions) SetRequestType(requestType string) *CreateOriginCertificateOptions {
	_options.RequestType = core.StringPtr(requestType)
	return _options
}

// SetRequestedValidity : Allow user to set RequestedValidity
func (_options *CreateOriginCertificateOptions) SetRequestedValidity(requestedValidity int64) *CreateOriginCertificateOptions {
	_options.RequestedValidity = core.Int64Ptr(requestedValidity)
	return _options
}

// SetCsr : Allow user to set Csr
func (_options *CreateOriginCertificateOptions) SetCsr(csr string) *CreateOriginCertificateOptions {
	_options.Csr = core.StringPtr(csr)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateOriginCertificateOptions) SetXCorrelationID(xCorrelationID string) *CreateOriginCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOriginCertificateOptions) SetHeaders(param map[string]string) *CreateOriginCertificateOptions {
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
	CustomCertReqGeoRestrictions_Label_Eu              = "eu"
	CustomCertReqGeoRestrictions_Label_HighestSecurity = "highest_security"
	CustomCertReqGeoRestrictions_Label_Us              = "us"
)

// NewCustomCertReqGeoRestrictions : Instantiate CustomCertReqGeoRestrictions (Generic Model Constructor)
func (*SslCertificateApiV1) NewCustomCertReqGeoRestrictions(label string) (_model *CustomCertReqGeoRestrictions, err error) {
	_model = &CustomCertReqGeoRestrictions{
		Label: core.StringPtr(label),
	}
	err = core.ValidateStruct(_model, "required parameters")
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
func (_options *DeleteCertificateOptions) SetCertIdentifier(certIdentifier string) *DeleteCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCertificateOptions) SetXCorrelationID(xCorrelationID string) *DeleteCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCertificateOptions) SetHeaders(param map[string]string) *DeleteCertificateOptions {
	options.Headers = param
	return options
}

// DeleteCertificateV2Options : The DeleteCertificateV2 options.
type DeleteCertificateV2Options struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCertificateV2Options : Instantiate DeleteCertificateV2Options
func (*SslCertificateApiV1) NewDeleteCertificateV2Options(certIdentifier string) *DeleteCertificateV2Options {
	return &DeleteCertificateV2Options{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *DeleteCertificateV2Options) SetCertIdentifier(certIdentifier string) *DeleteCertificateV2Options {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteCertificateV2Options) SetXCorrelationID(xCorrelationID string) *DeleteCertificateV2Options {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCertificateV2Options) SetHeaders(param map[string]string) *DeleteCertificateV2Options {
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
func (_options *DeleteCustomCertificateOptions) SetCustomCertID(customCertID string) *DeleteCustomCertificateOptions {
	_options.CustomCertID = core.StringPtr(customCertID)
	return _options
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
func (_options *GetCustomCertificateOptions) SetCustomCertID(customCertID string) *GetCustomCertificateOptions {
	_options.CustomCertID = core.StringPtr(customCertID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCustomCertificateOptions) SetHeaders(param map[string]string) *GetCustomCertificateOptions {
	options.Headers = param
	return options
}

// GetOriginCertificateOptions : The GetOriginCertificate options.
type GetOriginCertificateOptions struct {
	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `json:"crn" validate:"required,ne="`

	// zone identifier.
	ZoneIdentifier *string `json:"zone_identifier" validate:"required,ne="`

	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetOriginCertificateOptions : Instantiate GetOriginCertificateOptions
func (*SslCertificateApiV1) NewGetOriginCertificateOptions(crn string, zoneIdentifier string, certIdentifier string) *GetOriginCertificateOptions {
	return &GetOriginCertificateOptions{
		Crn:            core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCrn : Allow user to set Crn
func (_options *GetOriginCertificateOptions) SetCrn(crn string) *GetOriginCertificateOptions {
	_options.Crn = core.StringPtr(crn)
	return _options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (_options *GetOriginCertificateOptions) SetZoneIdentifier(zoneIdentifier string) *GetOriginCertificateOptions {
	_options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return _options
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *GetOriginCertificateOptions) SetCertIdentifier(certIdentifier string) *GetOriginCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetOriginCertificateOptions) SetXCorrelationID(xCorrelationID string) *GetOriginCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOriginCertificateOptions) SetHeaders(param map[string]string) *GetOriginCertificateOptions {
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

// GetSslVerificationOptions : The GetSslVerification options.
type GetSslVerificationOptions struct {
	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSslVerificationOptions : Instantiate GetSslVerificationOptions
func (*SslCertificateApiV1) NewGetSslVerificationOptions() *GetSslVerificationOptions {
	return &GetSslVerificationOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetSslVerificationOptions) SetXCorrelationID(xCorrelationID string) *GetSslVerificationOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSslVerificationOptions) SetHeaders(param map[string]string) *GetSslVerificationOptions {
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
func (_options *ListCertificatesOptions) SetXCorrelationID(xCorrelationID string) *ListCertificatesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
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

// ListOriginCertificatesOptions : The ListOriginCertificates options.
type ListOriginCertificatesOptions struct {
	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `json:"crn" validate:"required,ne="`

	// zone identifier.
	ZoneIdentifier *string `json:"zone_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListOriginCertificatesOptions : Instantiate ListOriginCertificatesOptions
func (*SslCertificateApiV1) NewListOriginCertificatesOptions(crn string, zoneIdentifier string) *ListOriginCertificatesOptions {
	return &ListOriginCertificatesOptions{
		Crn:            core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
	}
}

// SetCrn : Allow user to set Crn
func (_options *ListOriginCertificatesOptions) SetCrn(crn string) *ListOriginCertificatesOptions {
	_options.Crn = core.StringPtr(crn)
	return _options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (_options *ListOriginCertificatesOptions) SetZoneIdentifier(zoneIdentifier string) *ListOriginCertificatesOptions {
	_options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListOriginCertificatesOptions) SetXCorrelationID(xCorrelationID string) *ListOriginCertificatesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListOriginCertificatesOptions) SetHeaders(param map[string]string) *ListOriginCertificatesOptions {
	options.Headers = param
	return options
}

// ListOriginCertificatesRespResultInfo : Statistics of results.
type ListOriginCertificatesRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// UnmarshalListOriginCertificatesRespResultInfo unmarshals an instance of ListOriginCertificatesRespResultInfo from the specified map of raw messages.
func UnmarshalListOriginCertificatesRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListOriginCertificatesRespResultInfo)
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

// OrderAdvancedCertificateOptions : The OrderAdvancedCertificate options.
type OrderAdvancedCertificateOptions struct {
	// certificate type.
	Type *string `json:"type,omitempty"`

	// host name.
	Hosts []string `json:"hosts,omitempty"`

	// validation Method selected for the order.
	ValidationMethod *string `json:"validation_method,omitempty"`

	// validity Days selected for the order.
	ValidityDays *int64 `json:"validity_days,omitempty"`

	// Certificate Authority selected for the order. Selecting Let's Encrypt will reduce customization of other fields:
	// validation_method must be 'txt', validity_days must be 90, cloudflare_branding must be omitted, and hosts must
	// contain only 2 entries, one for the zone name and one for the subdomain wildcard of the zone name (e.g. example.com,
	// *.example.com).
	CertificateAuthority *string `json:"certificate_authority,omitempty"`

	// Whether or not to add Cloudflare Branding for the order. This will add sni.cloudflaressl.com as the Common Name if
	// set true.
	CloudflareBranding *bool `json:"cloudflare_branding,omitempty"`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the OrderAdvancedCertificateOptions.Type property.
// certificate type.
const (
	OrderAdvancedCertificateOptions_Type_Advanced = "advanced"
)

// Constants associated with the OrderAdvancedCertificateOptions.ValidationMethod property.
// validation Method selected for the order.
const (
	OrderAdvancedCertificateOptions_ValidationMethod_Email = "email"
	OrderAdvancedCertificateOptions_ValidationMethod_Http  = "http"
	OrderAdvancedCertificateOptions_ValidationMethod_Txt   = "txt"
)

// Constants associated with the OrderAdvancedCertificateOptions.CertificateAuthority property.
// Certificate Authority selected for the order. Selecting Let's Encrypt will reduce customization of other fields:
// validation_method must be 'txt', validity_days must be 90, cloudflare_branding must be omitted, and hosts must
// contain only 2 entries, one for the zone name and one for the subdomain wildcard of the zone name (e.g. example.com,
// *.example.com).
const (
	OrderAdvancedCertificateOptions_CertificateAuthority_Google      = "google"
	OrderAdvancedCertificateOptions_CertificateAuthority_LetsEncrypt = "lets_encrypt"
)

// NewOrderAdvancedCertificateOptions : Instantiate OrderAdvancedCertificateOptions
func (*SslCertificateApiV1) NewOrderAdvancedCertificateOptions() *OrderAdvancedCertificateOptions {
	return &OrderAdvancedCertificateOptions{}
}

// SetType : Allow user to set Type
func (_options *OrderAdvancedCertificateOptions) SetType(typeVar string) *OrderAdvancedCertificateOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetHosts : Allow user to set Hosts
func (_options *OrderAdvancedCertificateOptions) SetHosts(hosts []string) *OrderAdvancedCertificateOptions {
	_options.Hosts = hosts
	return _options
}

// SetValidationMethod : Allow user to set ValidationMethod
func (_options *OrderAdvancedCertificateOptions) SetValidationMethod(validationMethod string) *OrderAdvancedCertificateOptions {
	_options.ValidationMethod = core.StringPtr(validationMethod)
	return _options
}

// SetValidityDays : Allow user to set ValidityDays
func (_options *OrderAdvancedCertificateOptions) SetValidityDays(validityDays int64) *OrderAdvancedCertificateOptions {
	_options.ValidityDays = core.Int64Ptr(validityDays)
	return _options
}

// SetCertificateAuthority : Allow user to set CertificateAuthority
func (_options *OrderAdvancedCertificateOptions) SetCertificateAuthority(certificateAuthority string) *OrderAdvancedCertificateOptions {
	_options.CertificateAuthority = core.StringPtr(certificateAuthority)
	return _options
}

// SetCloudflareBranding : Allow user to set CloudflareBranding
func (_options *OrderAdvancedCertificateOptions) SetCloudflareBranding(cloudflareBranding bool) *OrderAdvancedCertificateOptions {
	_options.CloudflareBranding = core.BoolPtr(cloudflareBranding)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *OrderAdvancedCertificateOptions) SetXCorrelationID(xCorrelationID string) *OrderAdvancedCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *OrderAdvancedCertificateOptions) SetHeaders(param map[string]string) *OrderAdvancedCertificateOptions {
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
func (_options *OrderCertificateOptions) SetType(typeVar string) *OrderCertificateOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetHosts : Allow user to set Hosts
func (_options *OrderCertificateOptions) SetHosts(hosts []string) *OrderCertificateOptions {
	_options.Hosts = hosts
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *OrderCertificateOptions) SetXCorrelationID(xCorrelationID string) *OrderCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *OrderCertificateOptions) SetHeaders(param map[string]string) *OrderCertificateOptions {
	options.Headers = param
	return options
}

// PatchCertificateOptions : The PatchCertificate options.
type PatchCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPatchCertificateOptions : Instantiate PatchCertificateOptions
func (*SslCertificateApiV1) NewPatchCertificateOptions(certIdentifier string) *PatchCertificateOptions {
	return &PatchCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *PatchCertificateOptions) SetCertIdentifier(certIdentifier string) *PatchCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *PatchCertificateOptions) SetXCorrelationID(xCorrelationID string) *PatchCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PatchCertificateOptions) SetHeaders(param map[string]string) *PatchCertificateOptions {
	options.Headers = param
	return options
}

// RevokeOriginCertificateOptions : The RevokeOriginCertificate options.
type RevokeOriginCertificateOptions struct {
	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `json:"crn" validate:"required,ne="`

	// zone identifier.
	ZoneIdentifier *string `json:"zone_identifier" validate:"required,ne="`

	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// uuid, identify a session.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewRevokeOriginCertificateOptions : Instantiate RevokeOriginCertificateOptions
func (*SslCertificateApiV1) NewRevokeOriginCertificateOptions(crn string, zoneIdentifier string, certIdentifier string) *RevokeOriginCertificateOptions {
	return &RevokeOriginCertificateOptions{
		Crn:            core.StringPtr(crn),
		ZoneIdentifier: core.StringPtr(zoneIdentifier),
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCrn : Allow user to set Crn
func (_options *RevokeOriginCertificateOptions) SetCrn(crn string) *RevokeOriginCertificateOptions {
	_options.Crn = core.StringPtr(crn)
	return _options
}

// SetZoneIdentifier : Allow user to set ZoneIdentifier
func (_options *RevokeOriginCertificateOptions) SetZoneIdentifier(zoneIdentifier string) *RevokeOriginCertificateOptions {
	_options.ZoneIdentifier = core.StringPtr(zoneIdentifier)
	return _options
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *RevokeOriginCertificateOptions) SetCertIdentifier(certIdentifier string) *RevokeOriginCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *RevokeOriginCertificateOptions) SetXCorrelationID(xCorrelationID string) *RevokeOriginCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RevokeOriginCertificateOptions) SetHeaders(param map[string]string) *RevokeOriginCertificateOptions {
	options.Headers = param
	return options
}

// RevokeOriginCertificateRespResult : Container for response information.
type RevokeOriginCertificateRespResult struct {
	// ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalRevokeOriginCertificateRespResult unmarshals an instance of RevokeOriginCertificateRespResult from the specified map of raw messages.
func UnmarshalRevokeOriginCertificateRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RevokeOriginCertificateRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SslVerificationInfoVerificationInfo : certificate's required verification information.
type SslVerificationInfoVerificationInfo struct {
	// name of CNAME record.
	RecordName *string `json:"record_name,omitempty"`

	// target of CNAME record.
	RecordTarget *string `json:"record_target,omitempty"`
}

// UnmarshalSslVerificationInfoVerificationInfo unmarshals an instance of SslVerificationInfoVerificationInfo from the specified map of raw messages.
func UnmarshalSslVerificationInfoVerificationInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SslVerificationInfoVerificationInfo)
	err = core.UnmarshalPrimitive(m, "record_name", &obj.RecordName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "record_target", &obj.RecordTarget)
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
	UpdateCustomCertificateOptions_BundleMethod_Force      = "force"
	UpdateCustomCertificateOptions_BundleMethod_Optimal    = "optimal"
	UpdateCustomCertificateOptions_BundleMethod_Ubiquitous = "ubiquitous"
)

// NewUpdateCustomCertificateOptions : Instantiate UpdateCustomCertificateOptions
func (*SslCertificateApiV1) NewUpdateCustomCertificateOptions(customCertID string) *UpdateCustomCertificateOptions {
	return &UpdateCustomCertificateOptions{
		CustomCertID: core.StringPtr(customCertID),
	}
}

// SetCustomCertID : Allow user to set CustomCertID
func (_options *UpdateCustomCertificateOptions) SetCustomCertID(customCertID string) *UpdateCustomCertificateOptions {
	_options.CustomCertID = core.StringPtr(customCertID)
	return _options
}

// SetCertificate : Allow user to set Certificate
func (_options *UpdateCustomCertificateOptions) SetCertificate(certificate string) *UpdateCustomCertificateOptions {
	_options.Certificate = core.StringPtr(certificate)
	return _options
}

// SetPrivateKey : Allow user to set PrivateKey
func (_options *UpdateCustomCertificateOptions) SetPrivateKey(privateKey string) *UpdateCustomCertificateOptions {
	_options.PrivateKey = core.StringPtr(privateKey)
	return _options
}

// SetBundleMethod : Allow user to set BundleMethod
func (_options *UpdateCustomCertificateOptions) SetBundleMethod(bundleMethod string) *UpdateCustomCertificateOptions {
	_options.BundleMethod = core.StringPtr(bundleMethod)
	return _options
}

// SetGeoRestrictions : Allow user to set GeoRestrictions
func (_options *UpdateCustomCertificateOptions) SetGeoRestrictions(geoRestrictions *CustomCertReqGeoRestrictions) *UpdateCustomCertificateOptions {
	_options.GeoRestrictions = geoRestrictions
	return _options
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
	UploadCustomCertificateOptions_BundleMethod_Force      = "force"
	UploadCustomCertificateOptions_BundleMethod_Optimal    = "optimal"
	UploadCustomCertificateOptions_BundleMethod_Ubiquitous = "ubiquitous"
)

// NewUploadCustomCertificateOptions : Instantiate UploadCustomCertificateOptions
func (*SslCertificateApiV1) NewUploadCustomCertificateOptions() *UploadCustomCertificateOptions {
	return &UploadCustomCertificateOptions{}
}

// SetCertificate : Allow user to set Certificate
func (_options *UploadCustomCertificateOptions) SetCertificate(certificate string) *UploadCustomCertificateOptions {
	_options.Certificate = core.StringPtr(certificate)
	return _options
}

// SetPrivateKey : Allow user to set PrivateKey
func (_options *UploadCustomCertificateOptions) SetPrivateKey(privateKey string) *UploadCustomCertificateOptions {
	_options.PrivateKey = core.StringPtr(privateKey)
	return _options
}

// SetBundleMethod : Allow user to set BundleMethod
func (_options *UploadCustomCertificateOptions) SetBundleMethod(bundleMethod string) *UploadCustomCertificateOptions {
	_options.BundleMethod = core.StringPtr(bundleMethod)
	return _options
}

// SetGeoRestrictions : Allow user to set GeoRestrictions
func (_options *UploadCustomCertificateOptions) SetGeoRestrictions(geoRestrictions *CustomCertReqGeoRestrictions) *UploadCustomCertificateOptions {
	_options.GeoRestrictions = geoRestrictions
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UploadCustomCertificateOptions) SetHeaders(param map[string]string) *UploadCustomCertificateOptions {
	options.Headers = param
	return options
}

// AdvancedCertInitResp : certificate response.
type AdvancedCertInitResp struct {
	// result of ordering advanced certificate pack.
	Result *AdvancedCertInitRespResult `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
}

// UnmarshalAdvancedCertInitResp unmarshals an instance of AdvancedCertInitResp from the specified map of raw messages.
func UnmarshalAdvancedCertInitResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdvancedCertInitResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalAdvancedCertInitRespResult)
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Certificate : certificate.
type Certificate struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// host name.
	Hosts []string `json:"hosts" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`

	// issuer.
	Issuer *string `json:"issuer" validate:"required"`

	// signature.
	Signature *string `json:"signature" validate:"required"`

	// bundle method.
	BundleMethod *string `json:"bundle_method" validate:"required"`

	// zone ID.
	ZoneID *string `json:"zone_id" validate:"required"`

	// uploaded time.
	UploadedOn *string `json:"uploaded_on" validate:"required"`

	// modified time.
	ModifiedOn *string `json:"modified_on" validate:"required"`

	// expire time.
	ExpiresOn *string `json:"expires_on" validate:"required"`

	// certificate priority.
	Priority *int64 `json:"priority,omitempty"`
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
	err = core.UnmarshalPrimitive(m, "issuer", &obj.Issuer)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "signature", &obj.Signature)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	PrimaryCertificate map[string]interface{} `json:"primary_certificate" validate:"required"`

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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListOriginCertificatesResp : response of list origin certificates.
type ListOriginCertificatesResp struct {
	// Container for response information.
	Result []OriginCertificate `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListOriginCertificatesRespResultInfo `json:"result_info" validate:"required"`

	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`
}

// UnmarshalListOriginCertificatesResp unmarshals an instance of ListOriginCertificatesResp from the specified map of raw messages.
func UnmarshalListOriginCertificatesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListOriginCertificatesResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalOriginCertificate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListOriginCertificatesRespResultInfo)
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OriginCertificate : origin certificate.
type OriginCertificate struct {
	// ID.
	ID *string `json:"id" validate:"required"`

	// a CIS-signed certificatece.
	Certificate *string `json:"certificate" validate:"required"`

	// Array of hostnames or wildcard names (e.g., *.example.com) bound to the certificate.
	Hostnames []string `json:"hostnames" validate:"required"`

	// The expires date for this certificate.
	ExpiresOn *string `json:"expires_on" validate:"required"`

	// Signature type desired on certificate.
	RequestType *string `json:"request_type" validate:"required"`

	// The number of days for which the certificate should be valid.
	RequestedValidity *int64 `json:"requested_validity" validate:"required"`

	// The Certificate Signing Request (CSR).
	Csr *string `json:"csr" validate:"required"`

	// The private key.
	PrivateKey *string `json:"private_key,omitempty"`
}

// Constants associated with the OriginCertificate.RequestType property.
// Signature type desired on certificate.
const (
	OriginCertificate_RequestType_OriginEcc = "origin-ecc"
	OriginCertificate_RequestType_OriginRsa = "origin-rsa"
)

// UnmarshalOriginCertificate unmarshals an instance of OriginCertificate from the specified map of raw messages.
func UnmarshalOriginCertificate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OriginCertificate)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hostnames", &obj.Hostnames)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "request_type", &obj.RequestType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "requested_validity", &obj.RequestedValidity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "csr", &obj.Csr)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_key", &obj.PrivateKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OriginCertificateResp : response of origin certificate.
type OriginCertificateResp struct {
	// origin certificate.
	Result *OriginCertificate `json:"result" validate:"required"`

	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`
}

// UnmarshalOriginCertificateResp unmarshals an instance of OriginCertificateResp from the specified map of raw messages.
func UnmarshalOriginCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OriginCertificateResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalOriginCertificate)
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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

	// total pages.
	TotalPages *int64 `json:"total_pages" validate:"required"`

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
	err = core.UnmarshalPrimitive(m, "total_pages", &obj.TotalPages)
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

// RevokeOriginCertificateResp : response of revoke origin certificate.
type RevokeOriginCertificateResp struct {
	// Was the get successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *RevokeOriginCertificateRespResult `json:"result" validate:"required"`
}

// UnmarshalRevokeOriginCertificateResp unmarshals an instance of RevokeOriginCertificateResp from the specified map of raw messages.
func UnmarshalRevokeOriginCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RevokeOriginCertificateResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRevokeOriginCertificateRespResult)
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

// Constants associated with the SslSetting.Value property.
// value.
const (
	SslSetting_Value_False    = "false"
	SslSetting_Value_Flexible = "flexible"
	SslSetting_Value_Full     = "full"
	SslSetting_Value_Strict   = "strict"
)

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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SslVerificationInfo : ssl verification details.
type SslVerificationInfo struct {
	// current status of certificate.
	CertificateStatus *string `json:"certificate_status,omitempty"`

	// validation method in use for a certificate pack order.
	ValidationMethod *string `json:"validation_method,omitempty"`

	// method of certificate verification.
	VerificationType *string `json:"verification_type,omitempty"`

	// certificate pack identifier.
	CertPackUUID *string `json:"cert_pack_uuid,omitempty"`

	// status of the required verification information.
	VerificationStatus *bool `json:"verification_status,omitempty"`

	// certificate's required verification information.
	VerificationInfo *SslVerificationInfoVerificationInfo `json:"verification_info,omitempty"`

	// Wether or not Certificate Authority is manually reviewing the order.
	BrandCheck *bool `json:"brand_check,omitempty"`
}

// UnmarshalSslVerificationInfo unmarshals an instance of SslVerificationInfo from the specified map of raw messages.
func UnmarshalSslVerificationInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SslVerificationInfo)
	err = core.UnmarshalPrimitive(m, "certificate_status", &obj.CertificateStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "validation_method", &obj.ValidationMethod)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "verification_type", &obj.VerificationType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cert_pack_uuid", &obj.CertPackUUID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "verification_status", &obj.VerificationStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "verification_info", &obj.VerificationInfo, UnmarshalSslVerificationInfoVerificationInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "brand_check", &obj.BrandCheck)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SslVerificationResp : ssl verification response.
type SslVerificationResp struct {
	// ssl verification result.
	Result []SslVerificationInfo `json:"result" validate:"required"`

	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
}

// UnmarshalSslVerificationResp unmarshals an instance of SslVerificationResp from the specified map of raw messages.
func UnmarshalSslVerificationResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SslVerificationResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalSslVerificationInfo)
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
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
	Errors []map[string]interface{} `json:"errors" validate:"required"`

	// messages.
	Messages []map[string]interface{} `json:"messages" validate:"required"`
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
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
