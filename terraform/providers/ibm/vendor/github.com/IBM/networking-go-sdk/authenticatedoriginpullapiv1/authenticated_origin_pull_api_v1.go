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
 * IBM OpenAPI SDK Code Generator Version: 3.43.0-49eab5c7-20211117-152138
 */

// Package authenticatedoriginpullapiv1 : Operations and models for the AuthenticatedOriginPullApiV1 service
package authenticatedoriginpullapiv1

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

// AuthenticatedOriginPullApiV1 : Authenticated Origin Pull
//
// API Version: 1.0.0
type AuthenticatedOriginPullApiV1 struct {
	Service *core.BaseService

	// cloud resource name.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "authenticated_origin_pull_api"

// AuthenticatedOriginPullApiV1Options : Service options
type AuthenticatedOriginPullApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// cloud resource name.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewAuthenticatedOriginPullApiV1UsingExternalConfig : constructs an instance of AuthenticatedOriginPullApiV1 with passed in options and external configuration.
func NewAuthenticatedOriginPullApiV1UsingExternalConfig(options *AuthenticatedOriginPullApiV1Options) (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	authenticatedOriginPullApi, err = NewAuthenticatedOriginPullApiV1(options)
	if err != nil {
		return
	}

	err = authenticatedOriginPullApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = authenticatedOriginPullApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAuthenticatedOriginPullApiV1 : constructs an instance of AuthenticatedOriginPullApiV1 with passed in options.
func NewAuthenticatedOriginPullApiV1(options *AuthenticatedOriginPullApiV1Options) (service *AuthenticatedOriginPullApiV1, err error) {
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

	service = &AuthenticatedOriginPullApiV1{
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

// Clone makes a copy of "authenticatedOriginPullApi" suitable for processing requests.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) Clone() *AuthenticatedOriginPullApiV1 {
	if core.IsNil(authenticatedOriginPullApi) {
		return nil
	}
	clone := *authenticatedOriginPullApi
	clone.Service = authenticatedOriginPullApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetServiceURL(url string) error {
	return authenticatedOriginPullApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetServiceURL() string {
	return authenticatedOriginPullApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetDefaultHeaders(headers http.Header) {
	authenticatedOriginPullApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetEnableGzipCompression(enableGzip bool) {
	authenticatedOriginPullApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetEnableGzipCompression() bool {
	return authenticatedOriginPullApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	authenticatedOriginPullApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) DisableRetries() {
	authenticatedOriginPullApi.Service.DisableRetries()
}

// GetZoneOriginPullSettings : Get Zone level Authenticated Origin Pull Settings
// Get whether zone-level authenticated origin pulls is enabled or not. It is false by default.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetZoneOriginPullSettings(getZoneOriginPullSettingsOptions *GetZoneOriginPullSettingsOptions) (result *GetZoneOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.GetZoneOriginPullSettingsWithContext(context.Background(), getZoneOriginPullSettingsOptions)
}

// GetZoneOriginPullSettingsWithContext is an alternate form of the GetZoneOriginPullSettings method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetZoneOriginPullSettingsWithContext(ctx context.Context, getZoneOriginPullSettingsOptions *GetZoneOriginPullSettingsOptions) (result *GetZoneOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getZoneOriginPullSettingsOptions, "getZoneOriginPullSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneOriginPullSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "GetZoneOriginPullSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getZoneOriginPullSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getZoneOriginPullSettingsOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetZoneOriginPullSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetZoneOriginPullSettings : Set Zone level Authenticated Origin Pull Settings
// Enable or disable zone-level authenticated origin pulls. 'enabled' should be set true either before/after the
// certificate is uploaded to see the certificate in use.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetZoneOriginPullSettings(setZoneOriginPullSettingsOptions *SetZoneOriginPullSettingsOptions) (result *GetZoneOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.SetZoneOriginPullSettingsWithContext(context.Background(), setZoneOriginPullSettingsOptions)
}

// SetZoneOriginPullSettingsWithContext is an alternate form of the SetZoneOriginPullSettings method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetZoneOriginPullSettingsWithContext(ctx context.Context, setZoneOriginPullSettingsOptions *SetZoneOriginPullSettingsOptions) (result *GetZoneOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(setZoneOriginPullSettingsOptions, "setZoneOriginPullSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/settings`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setZoneOriginPullSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "SetZoneOriginPullSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if setZoneOriginPullSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*setZoneOriginPullSettingsOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if setZoneOriginPullSettingsOptions.Enabled != nil {
		body["enabled"] = setZoneOriginPullSettingsOptions.Enabled
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
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetZoneOriginPullSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListZoneOriginPullCertificates : List Zone level Authenticated Origin Pull Certificates
// List zone-level authenticated origin pulls certificates.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) ListZoneOriginPullCertificates(listZoneOriginPullCertificatesOptions *ListZoneOriginPullCertificatesOptions) (result *ListZoneOriginPullCertificatesResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.ListZoneOriginPullCertificatesWithContext(context.Background(), listZoneOriginPullCertificatesOptions)
}

// ListZoneOriginPullCertificatesWithContext is an alternate form of the ListZoneOriginPullCertificates method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) ListZoneOriginPullCertificatesWithContext(ctx context.Context, listZoneOriginPullCertificatesOptions *ListZoneOriginPullCertificatesOptions) (result *ListZoneOriginPullCertificatesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listZoneOriginPullCertificatesOptions, "listZoneOriginPullCertificatesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listZoneOriginPullCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "ListZoneOriginPullCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listZoneOriginPullCertificatesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listZoneOriginPullCertificatesOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListZoneOriginPullCertificatesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UploadZoneOriginPullCertificate : Upload Zone level Authenticated Origin Pull Certificate
// Upload your own certificate you want Cloudflare to use for edge-to-origin communication to override the shared
// certificate Please note that it is important to keep only one certificate active. Also, make sure to enable
// zone-level authenticated  origin pulls by making a PUT call to settings endpoint to see the uploaded certificate in
// use.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) UploadZoneOriginPullCertificate(uploadZoneOriginPullCertificateOptions *UploadZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.UploadZoneOriginPullCertificateWithContext(context.Background(), uploadZoneOriginPullCertificateOptions)
}

// UploadZoneOriginPullCertificateWithContext is an alternate form of the UploadZoneOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) UploadZoneOriginPullCertificateWithContext(ctx context.Context, uploadZoneOriginPullCertificateOptions *UploadZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(uploadZoneOriginPullCertificateOptions, "uploadZoneOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadZoneOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "UploadZoneOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if uploadZoneOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*uploadZoneOriginPullCertificateOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if uploadZoneOriginPullCertificateOptions.Certificate != nil {
		body["certificate"] = uploadZoneOriginPullCertificateOptions.Certificate
	}
	if uploadZoneOriginPullCertificateOptions.PrivateKey != nil {
		body["private_key"] = uploadZoneOriginPullCertificateOptions.PrivateKey
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
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetZoneOriginPullCertificate : Get a Zone level Authenticated Origin Pull Certificate
// Get a zone-level authenticated origin pulls certificate.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetZoneOriginPullCertificate(getZoneOriginPullCertificateOptions *GetZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.GetZoneOriginPullCertificateWithContext(context.Background(), getZoneOriginPullCertificateOptions)
}

// GetZoneOriginPullCertificateWithContext is an alternate form of the GetZoneOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetZoneOriginPullCertificateWithContext(ctx context.Context, getZoneOriginPullCertificateOptions *GetZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneOriginPullCertificateOptions, "getZoneOriginPullCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneOriginPullCertificateOptions, "getZoneOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
		"cert_identifier": *getZoneOriginPullCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "GetZoneOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getZoneOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getZoneOriginPullCertificateOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteZoneOriginPullCertificate : Delete a Zone level Authenticated Origin Pull Certificate
// Delete a zone-level authenticated origin pulls certificate.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) DeleteZoneOriginPullCertificate(deleteZoneOriginPullCertificateOptions *DeleteZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.DeleteZoneOriginPullCertificateWithContext(context.Background(), deleteZoneOriginPullCertificateOptions)
}

// DeleteZoneOriginPullCertificateWithContext is an alternate form of the DeleteZoneOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) DeleteZoneOriginPullCertificateWithContext(ctx context.Context, deleteZoneOriginPullCertificateOptions *DeleteZoneOriginPullCertificateOptions) (result *ZoneOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneOriginPullCertificateOptions, "deleteZoneOriginPullCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneOriginPullCertificateOptions, "deleteZoneOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
		"cert_identifier": *deleteZoneOriginPullCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "DeleteZoneOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteZoneOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteZoneOriginPullCertificateOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// SetHostnameOriginPullSettings : Set Hostname level Authenticated Origin Pull Settings
// Associate a hostname to a certificate and enable, disable or invalidate the association. If disabled, client
// certificate will not be sent to the hostname even if activated at the zone level. 100 maximum associations on a
// single certificate are allowed.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetHostnameOriginPullSettings(setHostnameOriginPullSettingsOptions *SetHostnameOriginPullSettingsOptions) (result *ListHostnameOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.SetHostnameOriginPullSettingsWithContext(context.Background(), setHostnameOriginPullSettingsOptions)
}

// SetHostnameOriginPullSettingsWithContext is an alternate form of the SetHostnameOriginPullSettings method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) SetHostnameOriginPullSettingsWithContext(ctx context.Context, setHostnameOriginPullSettingsOptions *SetHostnameOriginPullSettingsOptions) (result *ListHostnameOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(setHostnameOriginPullSettingsOptions, "setHostnameOriginPullSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/hostnames`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range setHostnameOriginPullSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "SetHostnameOriginPullSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if setHostnameOriginPullSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*setHostnameOriginPullSettingsOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if setHostnameOriginPullSettingsOptions.Config != nil {
		body["config"] = setHostnameOriginPullSettingsOptions.Config
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
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListHostnameOriginPullSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetHostnameOriginPullSettings : Get Hostname level Authenticated Origin Pull Settings
// Get hostname-level authenticated origin pulls settings.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetHostnameOriginPullSettings(getHostnameOriginPullSettingsOptions *GetHostnameOriginPullSettingsOptions) (result *GetHostnameOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.GetHostnameOriginPullSettingsWithContext(context.Background(), getHostnameOriginPullSettingsOptions)
}

// GetHostnameOriginPullSettingsWithContext is an alternate form of the GetHostnameOriginPullSettings method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetHostnameOriginPullSettingsWithContext(ctx context.Context, getHostnameOriginPullSettingsOptions *GetHostnameOriginPullSettingsOptions) (result *GetHostnameOriginPullSettingsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getHostnameOriginPullSettingsOptions, "getHostnameOriginPullSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getHostnameOriginPullSettingsOptions, "getHostnameOriginPullSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
		"hostname":        *getHostnameOriginPullSettingsOptions.Hostname,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/hostnames/{hostname}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getHostnameOriginPullSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "GetHostnameOriginPullSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getHostnameOriginPullSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getHostnameOriginPullSettingsOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetHostnameOriginPullSettingsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UploadHostnameOriginPullCertificate : Upload Hostname level Authenticated Origin Pull Certificate
// Upload a certificate to be used for client authentication on a hostname. 10 hostname certificates per zone are
// allowed.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) UploadHostnameOriginPullCertificate(uploadHostnameOriginPullCertificateOptions *UploadHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.UploadHostnameOriginPullCertificateWithContext(context.Background(), uploadHostnameOriginPullCertificateOptions)
}

// UploadHostnameOriginPullCertificateWithContext is an alternate form of the UploadHostnameOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) UploadHostnameOriginPullCertificateWithContext(ctx context.Context, uploadHostnameOriginPullCertificateOptions *UploadHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(uploadHostnameOriginPullCertificateOptions, "uploadHostnameOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/hostnames/certificates`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadHostnameOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "UploadHostnameOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if uploadHostnameOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*uploadHostnameOriginPullCertificateOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if uploadHostnameOriginPullCertificateOptions.Certificate != nil {
		body["certificate"] = uploadHostnameOriginPullCertificateOptions.Certificate
	}
	if uploadHostnameOriginPullCertificateOptions.PrivateKey != nil {
		body["private_key"] = uploadHostnameOriginPullCertificateOptions.PrivateKey
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
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHostnameOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetHostnameOriginPullCertificate : Get a Hostname level Authenticated Origin Pull Certificate
// Get the certificate by ID to be used for client authentication on a hostname.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetHostnameOriginPullCertificate(getHostnameOriginPullCertificateOptions *GetHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.GetHostnameOriginPullCertificateWithContext(context.Background(), getHostnameOriginPullCertificateOptions)
}

// GetHostnameOriginPullCertificateWithContext is an alternate form of the GetHostnameOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) GetHostnameOriginPullCertificateWithContext(ctx context.Context, getHostnameOriginPullCertificateOptions *GetHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getHostnameOriginPullCertificateOptions, "getHostnameOriginPullCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getHostnameOriginPullCertificateOptions, "getHostnameOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
		"cert_identifier": *getHostnameOriginPullCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/hostnames/certificates/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getHostnameOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "GetHostnameOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getHostnameOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getHostnameOriginPullCertificateOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHostnameOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteHostnameOriginPullCertificate : Delete a Hostname level Authenticated Origin Pull Certificate
// Delete the certificate by ID to be used for client authentication on a hostname.
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) DeleteHostnameOriginPullCertificate(deleteHostnameOriginPullCertificateOptions *DeleteHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	return authenticatedOriginPullApi.DeleteHostnameOriginPullCertificateWithContext(context.Background(), deleteHostnameOriginPullCertificateOptions)
}

// DeleteHostnameOriginPullCertificateWithContext is an alternate form of the DeleteHostnameOriginPullCertificate method which supports a Context parameter
func (authenticatedOriginPullApi *AuthenticatedOriginPullApiV1) DeleteHostnameOriginPullCertificateWithContext(ctx context.Context, deleteHostnameOriginPullCertificateOptions *DeleteHostnameOriginPullCertificateOptions) (result *HostnameOriginPullCertificateResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteHostnameOriginPullCertificateOptions, "deleteHostnameOriginPullCertificateOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteHostnameOriginPullCertificateOptions, "deleteHostnameOriginPullCertificateOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *authenticatedOriginPullApi.Crn,
		"zone_identifier": *authenticatedOriginPullApi.ZoneIdentifier,
		"cert_identifier": *deleteHostnameOriginPullCertificateOptions.CertIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = authenticatedOriginPullApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(authenticatedOriginPullApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/origin_tls_client_auth/hostnames/certificates/{cert_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteHostnameOriginPullCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("authenticated_origin_pull_api", "V1", "DeleteHostnameOriginPullCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteHostnameOriginPullCertificateOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteHostnameOriginPullCertificateOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = authenticatedOriginPullApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalHostnameOriginPullCertificateResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteHostnameOriginPullCertificateOptions : The DeleteHostnameOriginPullCertificate options.
type DeleteHostnameOriginPullCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteHostnameOriginPullCertificateOptions : Instantiate DeleteHostnameOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewDeleteHostnameOriginPullCertificateOptions(certIdentifier string) *DeleteHostnameOriginPullCertificateOptions {
	return &DeleteHostnameOriginPullCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *DeleteHostnameOriginPullCertificateOptions) SetCertIdentifier(certIdentifier string) *DeleteHostnameOriginPullCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteHostnameOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *DeleteHostnameOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteHostnameOriginPullCertificateOptions) SetHeaders(param map[string]string) *DeleteHostnameOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// DeleteZoneOriginPullCertificateOptions : The DeleteZoneOriginPullCertificate options.
type DeleteZoneOriginPullCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneOriginPullCertificateOptions : Instantiate DeleteZoneOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewDeleteZoneOriginPullCertificateOptions(certIdentifier string) *DeleteZoneOriginPullCertificateOptions {
	return &DeleteZoneOriginPullCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *DeleteZoneOriginPullCertificateOptions) SetCertIdentifier(certIdentifier string) *DeleteZoneOriginPullCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteZoneOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *DeleteZoneOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneOriginPullCertificateOptions) SetHeaders(param map[string]string) *DeleteZoneOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// GetHostnameOriginPullCertificateOptions : The GetHostnameOriginPullCertificate options.
type GetHostnameOriginPullCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetHostnameOriginPullCertificateOptions : Instantiate GetHostnameOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewGetHostnameOriginPullCertificateOptions(certIdentifier string) *GetHostnameOriginPullCertificateOptions {
	return &GetHostnameOriginPullCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *GetHostnameOriginPullCertificateOptions) SetCertIdentifier(certIdentifier string) *GetHostnameOriginPullCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetHostnameOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *GetHostnameOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetHostnameOriginPullCertificateOptions) SetHeaders(param map[string]string) *GetHostnameOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// GetHostnameOriginPullSettingsOptions : The GetHostnameOriginPullSettings options.
type GetHostnameOriginPullSettingsOptions struct {
	// the hostname on the origin for which the client certificate associate.
	Hostname *string `json:"hostname" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetHostnameOriginPullSettingsOptions : Instantiate GetHostnameOriginPullSettingsOptions
func (*AuthenticatedOriginPullApiV1) NewGetHostnameOriginPullSettingsOptions(hostname string) *GetHostnameOriginPullSettingsOptions {
	return &GetHostnameOriginPullSettingsOptions{
		Hostname: core.StringPtr(hostname),
	}
}

// SetHostname : Allow user to set Hostname
func (_options *GetHostnameOriginPullSettingsOptions) SetHostname(hostname string) *GetHostnameOriginPullSettingsOptions {
	_options.Hostname = core.StringPtr(hostname)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetHostnameOriginPullSettingsOptions) SetXCorrelationID(xCorrelationID string) *GetHostnameOriginPullSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetHostnameOriginPullSettingsOptions) SetHeaders(param map[string]string) *GetHostnameOriginPullSettingsOptions {
	options.Headers = param
	return options
}

// GetZoneOriginPullCertificateOptions : The GetZoneOriginPullCertificate options.
type GetZoneOriginPullCertificateOptions struct {
	// cedrtificate identifier.
	CertIdentifier *string `json:"cert_identifier" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneOriginPullCertificateOptions : Instantiate GetZoneOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewGetZoneOriginPullCertificateOptions(certIdentifier string) *GetZoneOriginPullCertificateOptions {
	return &GetZoneOriginPullCertificateOptions{
		CertIdentifier: core.StringPtr(certIdentifier),
	}
}

// SetCertIdentifier : Allow user to set CertIdentifier
func (_options *GetZoneOriginPullCertificateOptions) SetCertIdentifier(certIdentifier string) *GetZoneOriginPullCertificateOptions {
	_options.CertIdentifier = core.StringPtr(certIdentifier)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetZoneOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *GetZoneOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneOriginPullCertificateOptions) SetHeaders(param map[string]string) *GetZoneOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// GetZoneOriginPullSettingsOptions : The GetZoneOriginPullSettings options.
type GetZoneOriginPullSettingsOptions struct {
	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneOriginPullSettingsOptions : Instantiate GetZoneOriginPullSettingsOptions
func (*AuthenticatedOriginPullApiV1) NewGetZoneOriginPullSettingsOptions() *GetZoneOriginPullSettingsOptions {
	return &GetZoneOriginPullSettingsOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetZoneOriginPullSettingsOptions) SetXCorrelationID(xCorrelationID string) *GetZoneOriginPullSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneOriginPullSettingsOptions) SetHeaders(param map[string]string) *GetZoneOriginPullSettingsOptions {
	options.Headers = param
	return options
}

// GetZoneOriginPullSettingsRespResult : result.
type GetZoneOriginPullSettingsRespResult struct {
	// enabled.
	Enabled *bool `json:"enabled" validate:"required"`
}

// UnmarshalGetZoneOriginPullSettingsRespResult unmarshals an instance of GetZoneOriginPullSettingsRespResult from the specified map of raw messages.
func UnmarshalGetZoneOriginPullSettingsRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetZoneOriginPullSettingsRespResult)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListZoneOriginPullCertificatesOptions : The ListZoneOriginPullCertificates options.
type ListZoneOriginPullCertificatesOptions struct {
	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListZoneOriginPullCertificatesOptions : Instantiate ListZoneOriginPullCertificatesOptions
func (*AuthenticatedOriginPullApiV1) NewListZoneOriginPullCertificatesOptions() *ListZoneOriginPullCertificatesOptions {
	return &ListZoneOriginPullCertificatesOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListZoneOriginPullCertificatesOptions) SetXCorrelationID(xCorrelationID string) *ListZoneOriginPullCertificatesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListZoneOriginPullCertificatesOptions) SetHeaders(param map[string]string) *ListZoneOriginPullCertificatesOptions {
	options.Headers = param
	return options
}

// SetHostnameOriginPullSettingsOptions : The SetHostnameOriginPullSettings options.
type SetHostnameOriginPullSettingsOptions struct {
	// An array with items in the settings request.
	Config []HostnameOriginPullSettings `json:"config,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetHostnameOriginPullSettingsOptions : Instantiate SetHostnameOriginPullSettingsOptions
func (*AuthenticatedOriginPullApiV1) NewSetHostnameOriginPullSettingsOptions() *SetHostnameOriginPullSettingsOptions {
	return &SetHostnameOriginPullSettingsOptions{}
}

// SetConfig : Allow user to set Config
func (_options *SetHostnameOriginPullSettingsOptions) SetConfig(config []HostnameOriginPullSettings) *SetHostnameOriginPullSettingsOptions {
	_options.Config = config
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *SetHostnameOriginPullSettingsOptions) SetXCorrelationID(xCorrelationID string) *SetHostnameOriginPullSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetHostnameOriginPullSettingsOptions) SetHeaders(param map[string]string) *SetHostnameOriginPullSettingsOptions {
	options.Headers = param
	return options
}

// SetZoneOriginPullSettingsOptions : The SetZoneOriginPullSettings options.
type SetZoneOriginPullSettingsOptions struct {
	// enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewSetZoneOriginPullSettingsOptions : Instantiate SetZoneOriginPullSettingsOptions
func (*AuthenticatedOriginPullApiV1) NewSetZoneOriginPullSettingsOptions() *SetZoneOriginPullSettingsOptions {
	return &SetZoneOriginPullSettingsOptions{}
}

// SetEnabled : Allow user to set Enabled
func (_options *SetZoneOriginPullSettingsOptions) SetEnabled(enabled bool) *SetZoneOriginPullSettingsOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *SetZoneOriginPullSettingsOptions) SetXCorrelationID(xCorrelationID string) *SetZoneOriginPullSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetZoneOriginPullSettingsOptions) SetHeaders(param map[string]string) *SetZoneOriginPullSettingsOptions {
	options.Headers = param
	return options
}

// UploadHostnameOriginPullCertificateOptions : The UploadHostnameOriginPullCertificate options.
type UploadHostnameOriginPullCertificateOptions struct {
	// the zone's leaf certificate.
	Certificate *string `json:"certificate,omitempty"`

	// the zone's private key.
	PrivateKey *string `json:"private_key,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUploadHostnameOriginPullCertificateOptions : Instantiate UploadHostnameOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewUploadHostnameOriginPullCertificateOptions() *UploadHostnameOriginPullCertificateOptions {
	return &UploadHostnameOriginPullCertificateOptions{}
}

// SetCertificate : Allow user to set Certificate
func (_options *UploadHostnameOriginPullCertificateOptions) SetCertificate(certificate string) *UploadHostnameOriginPullCertificateOptions {
	_options.Certificate = core.StringPtr(certificate)
	return _options
}

// SetPrivateKey : Allow user to set PrivateKey
func (_options *UploadHostnameOriginPullCertificateOptions) SetPrivateKey(privateKey string) *UploadHostnameOriginPullCertificateOptions {
	_options.PrivateKey = core.StringPtr(privateKey)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UploadHostnameOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *UploadHostnameOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UploadHostnameOriginPullCertificateOptions) SetHeaders(param map[string]string) *UploadHostnameOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// UploadZoneOriginPullCertificateOptions : The UploadZoneOriginPullCertificate options.
type UploadZoneOriginPullCertificateOptions struct {
	// the zone's leaf certificate.
	Certificate *string `json:"certificate,omitempty"`

	// the zone's private key.
	PrivateKey *string `json:"private_key,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUploadZoneOriginPullCertificateOptions : Instantiate UploadZoneOriginPullCertificateOptions
func (*AuthenticatedOriginPullApiV1) NewUploadZoneOriginPullCertificateOptions() *UploadZoneOriginPullCertificateOptions {
	return &UploadZoneOriginPullCertificateOptions{}
}

// SetCertificate : Allow user to set Certificate
func (_options *UploadZoneOriginPullCertificateOptions) SetCertificate(certificate string) *UploadZoneOriginPullCertificateOptions {
	_options.Certificate = core.StringPtr(certificate)
	return _options
}

// SetPrivateKey : Allow user to set PrivateKey
func (_options *UploadZoneOriginPullCertificateOptions) SetPrivateKey(privateKey string) *UploadZoneOriginPullCertificateOptions {
	_options.PrivateKey = core.StringPtr(privateKey)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *UploadZoneOriginPullCertificateOptions) SetXCorrelationID(xCorrelationID string) *UploadZoneOriginPullCertificateOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UploadZoneOriginPullCertificateOptions) SetHeaders(param map[string]string) *UploadZoneOriginPullCertificateOptions {
	options.Headers = param
	return options
}

// CertificatePack : certificate pack.
type CertificatePack struct {
	// certificate identifier tag.
	ID *string `json:"id,omitempty"`

	// the zone's leaf certificate.
	Certificate *string `json:"certificate,omitempty"`

	// the certificate authority that issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// the type of hash used for the certificate.
	Signature *string `json:"signature,omitempty"`

	// status of the certificate activation.
	Status *string `json:"status,omitempty"`

	// when the certificate from the authority expires.
	ExpiresOn *string `json:"expires_on,omitempty"`

	// the time the certificate was uploaded.
	UploadedOn *string `json:"uploaded_on,omitempty"`
}

// UnmarshalCertificatePack unmarshals an instance of CertificatePack from the specified map of raw messages.
func UnmarshalCertificatePack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificatePack)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
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
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uploaded_on", &obj.UploadedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetHostnameOriginPullSettingsResp : detail of the hostname level authenticated origin pull settings response.
type GetHostnameOriginPullSettingsResp struct {
	// hostname level authenticated origin pull settings response.
	Result *HostnameSettingsResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalGetHostnameOriginPullSettingsResp unmarshals an instance of GetHostnameOriginPullSettingsResp from the specified map of raw messages.
func UnmarshalGetHostnameOriginPullSettingsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetHostnameOriginPullSettingsResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalHostnameSettingsResp)
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

// GetZoneOriginPullSettingsResp : zone level authenticated origin pull settings response.
type GetZoneOriginPullSettingsResp struct {
	// result.
	Result *GetZoneOriginPullSettingsRespResult `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalGetZoneOriginPullSettingsResp unmarshals an instance of GetZoneOriginPullSettingsResp from the specified map of raw messages.
func UnmarshalGetZoneOriginPullSettingsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetZoneOriginPullSettingsResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalGetZoneOriginPullSettingsRespResult)
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

// HostnameCertificatePack : certificate pack.
type HostnameCertificatePack struct {
	// certificate identifier tag.
	ID *string `json:"id,omitempty"`

	// the zone's leaf certificate.
	Certificate *string `json:"certificate,omitempty"`

	// the certificate authority that issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// the type of hash used for the certificate.
	Signature *string `json:"signature,omitempty"`

	// the serial number on the uploaded certificate.
	SerialNumber *string `json:"serial_number,omitempty"`

	// status of the certificate activation.
	Status *string `json:"status,omitempty"`

	// when the certificate from the authority expires.
	ExpiresOn *string `json:"expires_on,omitempty"`

	// the time the certificate was uploaded.
	UploadedOn *string `json:"uploaded_on,omitempty"`
}

// UnmarshalHostnameCertificatePack unmarshals an instance of HostnameCertificatePack from the specified map of raw messages.
func UnmarshalHostnameCertificatePack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostnameCertificatePack)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
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
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expires_on", &obj.ExpiresOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "uploaded_on", &obj.UploadedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HostnameOriginPullCertificateResp : certificate response.
type HostnameOriginPullCertificateResp struct {
	// certificate pack.
	Result *HostnameCertificatePack `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalHostnameOriginPullCertificateResp unmarshals an instance of HostnameOriginPullCertificateResp from the specified map of raw messages.
func UnmarshalHostnameOriginPullCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostnameOriginPullCertificateResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalHostnameCertificatePack)
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

// HostnameOriginPullSettings : hostname-level authenticated origin pull settings request.
type HostnameOriginPullSettings struct {
	// the hostname on the origin for which the client certificate uploaded will be used.
	Hostname *string `json:"hostname" validate:"required"`

	// certificate identifier tag.
	CertID *string `json:"cert_id" validate:"required"`

	// enabled.
	Enabled *bool `json:"enabled" validate:"required"`
}

// NewHostnameOriginPullSettings : Instantiate HostnameOriginPullSettings (Generic Model Constructor)
func (*AuthenticatedOriginPullApiV1) NewHostnameOriginPullSettings(hostname string, certID string, enabled bool) (_model *HostnameOriginPullSettings, err error) {
	_model = &HostnameOriginPullSettings{
		Hostname: core.StringPtr(hostname),
		CertID:   core.StringPtr(certID),
		Enabled:  core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalHostnameOriginPullSettings unmarshals an instance of HostnameOriginPullSettings from the specified map of raw messages.
func UnmarshalHostnameOriginPullSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostnameOriginPullSettings)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cert_id", &obj.CertID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// HostnameSettingsResp : hostname level authenticated origin pull settings response.
type HostnameSettingsResp struct {
	// the hostname on the origin for which the client certificate uploaded will be used.
	Hostname *string `json:"hostname,omitempty"`

	// certificate identifier tag.
	CertID *string `json:"cert_id,omitempty"`

	// enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// status of the certificate activation.
	Status *string `json:"status,omitempty"`

	// the time when the certificate was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// the time when the certificate was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// status of the certificate or the association.
	CertStatus *string `json:"cert_status,omitempty"`

	// the certificate authority that issued the certificate.
	Issuer *string `json:"issuer,omitempty"`

	// the type of hash used for the certificate.
	Signature *string `json:"signature,omitempty"`

	// the serial number on the uploaded certificate.
	SerialNumber *string `json:"serial_number,omitempty"`

	// the zone's leaf certificate.
	Certificate *string `json:"certificate,omitempty"`

	// the time the certificate was uploaded.
	CertUploadedOn *strfmt.DateTime `json:"cert_uploaded_on,omitempty"`

	// the time when the certificate was updated.
	CertUpdatedAt *strfmt.DateTime `json:"cert_updated_at,omitempty"`

	// the date when the certificate expires.
	ExpiresOn *strfmt.DateTime `json:"expires_on,omitempty"`
}

// UnmarshalHostnameSettingsResp unmarshals an instance of HostnameSettingsResp from the specified map of raw messages.
func UnmarshalHostnameSettingsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(HostnameSettingsResp)
	err = core.UnmarshalPrimitive(m, "hostname", &obj.Hostname)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cert_id", &obj.CertID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
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
	err = core.UnmarshalPrimitive(m, "cert_status", &obj.CertStatus)
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
	err = core.UnmarshalPrimitive(m, "serial_number", &obj.SerialNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate", &obj.Certificate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cert_uploaded_on", &obj.CertUploadedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cert_updated_at", &obj.CertUpdatedAt)
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

// ListHostnameOriginPullSettingsResp : array of hostname level authenticated origin pull settings response.
type ListHostnameOriginPullSettingsResp struct {
	// array of hostname settings response.
	Result []HostnameSettingsResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalListHostnameOriginPullSettingsResp unmarshals an instance of ListHostnameOriginPullSettingsResp from the specified map of raw messages.
func UnmarshalListHostnameOriginPullSettingsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListHostnameOriginPullSettingsResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalHostnameSettingsResp)
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

// ListZoneOriginPullCertificatesResp : certificate response.
type ListZoneOriginPullCertificatesResp struct {
	// list certificate packs.
	Result []CertificatePack `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalListZoneOriginPullCertificatesResp unmarshals an instance of ListZoneOriginPullCertificatesResp from the specified map of raw messages.
func UnmarshalListZoneOriginPullCertificatesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListZoneOriginPullCertificatesResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCertificatePack)
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

// ZoneOriginPullCertificateResp : zone level authenticated origin pull certificate response.
type ZoneOriginPullCertificateResp struct {
	// certificate pack.
	Result *CertificatePack `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	Errors []string `json:"errors,omitempty"`

	Messages []string `json:"messages,omitempty"`
}

// UnmarshalZoneOriginPullCertificateResp unmarshals an instance of ZoneOriginPullCertificateResp from the specified map of raw messages.
func UnmarshalZoneOriginPullCertificateResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneOriginPullCertificateResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCertificatePack)
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
