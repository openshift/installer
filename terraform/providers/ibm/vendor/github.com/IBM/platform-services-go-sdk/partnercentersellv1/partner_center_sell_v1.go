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
 * IBM OpenAPI SDK Code Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

// Package partnercentersellv1 : Operations and models for the PartnerCenterSellV1 service
package partnercentersellv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// PartnerCenterSellV1 : This API is experimental and is likely to change in the future.
//
// API Version: 1.5.0
type PartnerCenterSellV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://product-lifecycle.cloud.ibm.com/openapi/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "partner_center_sell"

// PartnerCenterSellV1Options : Service options
type PartnerCenterSellV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewPartnerCenterSellV1UsingExternalConfig : constructs an instance of PartnerCenterSellV1 with passed in options and external configuration.
func NewPartnerCenterSellV1UsingExternalConfig(options *PartnerCenterSellV1Options) (partnerCenterSell *PartnerCenterSellV1, err error) {
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

	partnerCenterSell, err = NewPartnerCenterSellV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = partnerCenterSell.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = partnerCenterSell.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewPartnerCenterSellV1 : constructs an instance of PartnerCenterSellV1 with passed in options.
func NewPartnerCenterSellV1(options *PartnerCenterSellV1Options) (service *PartnerCenterSellV1, err error) {
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

	service = &PartnerCenterSellV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "partnerCenterSell" suitable for processing requests.
func (partnerCenterSell *PartnerCenterSellV1) Clone() *PartnerCenterSellV1 {
	if core.IsNil(partnerCenterSell) {
		return nil
	}
	clone := *partnerCenterSell
	clone.Service = partnerCenterSell.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (partnerCenterSell *PartnerCenterSellV1) SetServiceURL(url string) error {
	err := partnerCenterSell.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (partnerCenterSell *PartnerCenterSellV1) GetServiceURL() string {
	return partnerCenterSell.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (partnerCenterSell *PartnerCenterSellV1) SetDefaultHeaders(headers http.Header) {
	partnerCenterSell.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (partnerCenterSell *PartnerCenterSellV1) SetEnableGzipCompression(enableGzip bool) {
	partnerCenterSell.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (partnerCenterSell *PartnerCenterSellV1) GetEnableGzipCompression() bool {
	return partnerCenterSell.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (partnerCenterSell *PartnerCenterSellV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	partnerCenterSell.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (partnerCenterSell *PartnerCenterSellV1) DisableRetries() {
	partnerCenterSell.Service.DisableRetries()
}

// CreateRegistration : Register your account in Partner Center - Sell
// Register your account in Partner Center - Sell to onboard your product.
func (partnerCenterSell *PartnerCenterSellV1) CreateRegistration(createRegistrationOptions *CreateRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateRegistrationWithContext(context.Background(), createRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateRegistrationWithContext is an alternate form of the CreateRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateRegistrationWithContext(ctx context.Context, createRegistrationOptions *CreateRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRegistrationOptions, "createRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createRegistrationOptions, "createRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/registration`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRegistrationOptions.AccountID != nil {
		body["account_id"] = createRegistrationOptions.AccountID
	}
	if createRegistrationOptions.CompanyName != nil {
		body["company_name"] = createRegistrationOptions.CompanyName
	}
	if createRegistrationOptions.PrimaryContact != nil {
		body["primary_contact"] = createRegistrationOptions.PrimaryContact
	}
	if createRegistrationOptions.DefaultPrivateCatalogID != nil {
		body["default_private_catalog_id"] = createRegistrationOptions.DefaultPrivateCatalogID
	}
	if createRegistrationOptions.ProviderAccessGroup != nil {
		body["provider_access_group"] = createRegistrationOptions.ProviderAccessGroup
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetRegistration : Retrieve a Partner Center - Sell registration
// Retrieve a Partner Center - Sell registration.
func (partnerCenterSell *PartnerCenterSellV1) GetRegistration(getRegistrationOptions *GetRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetRegistrationWithContext(context.Background(), getRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetRegistrationWithContext is an alternate form of the GetRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetRegistrationWithContext(ctx context.Context, getRegistrationOptions *GetRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRegistrationOptions, "getRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getRegistrationOptions, "getRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"registration_id": *getRegistrationOptions.RegistrationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/registration/{registration_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetRegistration")
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateRegistration : Update a Partner Center - Sell registration
// Update your registration in Partner Center - Sell.
func (partnerCenterSell *PartnerCenterSellV1) UpdateRegistration(updateRegistrationOptions *UpdateRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateRegistrationWithContext(context.Background(), updateRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateRegistrationWithContext is an alternate form of the UpdateRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateRegistrationWithContext(ctx context.Context, updateRegistrationOptions *UpdateRegistrationOptions) (result *Registration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRegistrationOptions, "updateRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateRegistrationOptions, "updateRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"registration_id": *updateRegistrationOptions.RegistrationID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/registration/{registration_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateRegistrationOptions.RegistrationPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteRegistration : Delete your registration in Partner - Center Sell
// Delete a Partner Center - Sell registration.
func (partnerCenterSell *PartnerCenterSellV1) DeleteRegistration(deleteRegistrationOptions *DeleteRegistrationOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteRegistrationWithContext(context.Background(), deleteRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteRegistrationWithContext is an alternate form of the DeleteRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteRegistrationWithContext(ctx context.Context, deleteRegistrationOptions *DeleteRegistrationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRegistrationOptions, "deleteRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteRegistrationOptions, "deleteRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"registration_id": *deleteRegistrationOptions.RegistrationID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/registration/{registration_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateOnboardingProduct : Create a product to onboard
// Create an onboarding product in Partner Center. An onboarding product is Partner Center - Sell's object representing
// a product.
func (partnerCenterSell *PartnerCenterSellV1) CreateOnboardingProduct(createOnboardingProductOptions *CreateOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateOnboardingProductWithContext(context.Background(), createOnboardingProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateOnboardingProductWithContext is an alternate form of the CreateOnboardingProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateOnboardingProductWithContext(ctx context.Context, createOnboardingProductOptions *CreateOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createOnboardingProductOptions, "createOnboardingProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createOnboardingProductOptions, "createOnboardingProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createOnboardingProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateOnboardingProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createOnboardingProductOptions.Type != nil {
		body["type"] = createOnboardingProductOptions.Type
	}
	if createOnboardingProductOptions.PrimaryContact != nil {
		body["primary_contact"] = createOnboardingProductOptions.PrimaryContact
	}
	if createOnboardingProductOptions.EccnNumber != nil {
		body["eccn_number"] = createOnboardingProductOptions.EccnNumber
	}
	if createOnboardingProductOptions.EroClass != nil {
		body["ero_class"] = createOnboardingProductOptions.EroClass
	}
	if createOnboardingProductOptions.Unspsc != nil {
		body["unspsc"] = createOnboardingProductOptions.Unspsc
	}
	if createOnboardingProductOptions.TaxAssessment != nil {
		body["tax_assessment"] = createOnboardingProductOptions.TaxAssessment
	}
	if createOnboardingProductOptions.Support != nil {
		body["support"] = createOnboardingProductOptions.Support
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_onboarding_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOnboardingProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetOnboardingProduct : Get an onboarding product
// Retrieve the details of a product in Partner Center.
func (partnerCenterSell *PartnerCenterSellV1) GetOnboardingProduct(getOnboardingProductOptions *GetOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetOnboardingProductWithContext(context.Background(), getOnboardingProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetOnboardingProductWithContext is an alternate form of the GetOnboardingProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetOnboardingProductWithContext(ctx context.Context, getOnboardingProductOptions *GetOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOnboardingProductOptions, "getOnboardingProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getOnboardingProductOptions, "getOnboardingProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id": *getOnboardingProductOptions.ProductID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getOnboardingProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetOnboardingProduct")
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_onboarding_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOnboardingProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateOnboardingProduct : Update an onboarding product
// Update the details of a product in Partner Center.
func (partnerCenterSell *PartnerCenterSellV1) UpdateOnboardingProduct(updateOnboardingProductOptions *UpdateOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateOnboardingProductWithContext(context.Background(), updateOnboardingProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateOnboardingProductWithContext is an alternate form of the UpdateOnboardingProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateOnboardingProductWithContext(ctx context.Context, updateOnboardingProductOptions *UpdateOnboardingProductOptions) (result *OnboardingProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateOnboardingProductOptions, "updateOnboardingProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateOnboardingProductOptions, "updateOnboardingProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id": *updateOnboardingProductOptions.ProductID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateOnboardingProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateOnboardingProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateOnboardingProductOptions.OnboardingProductPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_onboarding_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOnboardingProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteOnboardingProduct : Delete an onboarding product
// Delete a product in Partner Center.
func (partnerCenterSell *PartnerCenterSellV1) DeleteOnboardingProduct(deleteOnboardingProductOptions *DeleteOnboardingProductOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteOnboardingProductWithContext(context.Background(), deleteOnboardingProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteOnboardingProductWithContext is an alternate form of the DeleteOnboardingProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteOnboardingProductWithContext(ctx context.Context, deleteOnboardingProductOptions *DeleteOnboardingProductOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteOnboardingProductOptions, "deleteOnboardingProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteOnboardingProductOptions, "deleteOnboardingProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id": *deleteOnboardingProductOptions.ProductID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteOnboardingProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteOnboardingProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_onboarding_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateCatalogProduct : Create a global catalog product
// Create a product for global catalog and link it to a product in Partner Center. A global catalog product is the data
// used by the platform for service metadata definition.
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogProduct(createCatalogProductOptions *CreateCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateCatalogProductWithContext(context.Background(), createCatalogProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCatalogProductWithContext is an alternate form of the CreateCatalogProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogProductWithContext(ctx context.Context, createCatalogProductOptions *CreateCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCatalogProductOptions, "createCatalogProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCatalogProductOptions, "createCatalogProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id": *createCatalogProductOptions.ProductID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCatalogProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateCatalogProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createCatalogProductOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*createCatalogProductOptions.Env))
	}

	body := make(map[string]interface{})
	if createCatalogProductOptions.Name != nil {
		body["name"] = createCatalogProductOptions.Name
	}
	if createCatalogProductOptions.Active != nil {
		body["active"] = createCatalogProductOptions.Active
	}
	if createCatalogProductOptions.Disabled != nil {
		body["disabled"] = createCatalogProductOptions.Disabled
	}
	if createCatalogProductOptions.Kind != nil {
		body["kind"] = createCatalogProductOptions.Kind
	}
	if createCatalogProductOptions.Tags != nil {
		body["tags"] = createCatalogProductOptions.Tags
	}
	if createCatalogProductOptions.ObjectProvider != nil {
		body["object_provider"] = createCatalogProductOptions.ObjectProvider
	}
	if createCatalogProductOptions.ID != nil {
		body["id"] = createCatalogProductOptions.ID
	}
	if createCatalogProductOptions.ObjectID != nil {
		body["object_id"] = createCatalogProductOptions.ObjectID
	}
	if createCatalogProductOptions.OverviewUi != nil {
		body["overview_ui"] = createCatalogProductOptions.OverviewUi
	}
	if createCatalogProductOptions.Images != nil {
		body["images"] = createCatalogProductOptions.Images
	}
	if createCatalogProductOptions.Metadata != nil {
		body["metadata"] = createCatalogProductOptions.Metadata
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_catalog_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogProduct : Get a global catalog product
// Retrieve the details of a product in the global catalog.
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogProduct(getCatalogProductOptions *GetCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetCatalogProductWithContext(context.Background(), getCatalogProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCatalogProductWithContext is an alternate form of the GetCatalogProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogProductWithContext(ctx context.Context, getCatalogProductOptions *GetCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogProductOptions, "getCatalogProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCatalogProductOptions, "getCatalogProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *getCatalogProductOptions.ProductID,
		"catalog_product_id": *getCatalogProductOptions.CatalogProductID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCatalogProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetCatalogProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCatalogProductOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*getCatalogProductOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_catalog_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCatalogProduct : Update a global catalog product
// Update the details of a product in global catalog.
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogProduct(updateCatalogProductOptions *UpdateCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateCatalogProductWithContext(context.Background(), updateCatalogProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCatalogProductWithContext is an alternate form of the UpdateCatalogProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogProductWithContext(ctx context.Context, updateCatalogProductOptions *UpdateCatalogProductOptions) (result *GlobalCatalogProduct, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCatalogProductOptions, "updateCatalogProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCatalogProductOptions, "updateCatalogProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *updateCatalogProductOptions.ProductID,
		"catalog_product_id": *updateCatalogProductOptions.CatalogProductID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCatalogProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateCatalogProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateCatalogProductOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*updateCatalogProductOptions.Env))
	}

	_, err = builder.SetBodyContentJSON(updateCatalogProductOptions.GlobalCatalogProductPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_catalog_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogProduct)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCatalogProduct : Delete a global catalog product
// Delete a product from the global catalog.
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogProduct(deleteCatalogProductOptions *DeleteCatalogProductOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteCatalogProductWithContext(context.Background(), deleteCatalogProductOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCatalogProductWithContext is an alternate form of the DeleteCatalogProduct method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogProductWithContext(ctx context.Context, deleteCatalogProductOptions *DeleteCatalogProductOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCatalogProductOptions, "deleteCatalogProductOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCatalogProductOptions, "deleteCatalogProductOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *deleteCatalogProductOptions.ProductID,
		"catalog_product_id": *deleteCatalogProductOptions.CatalogProductID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCatalogProductOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteCatalogProduct")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteCatalogProductOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*deleteCatalogProductOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_catalog_product", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateCatalogPlan : Create a pricing plan in global catalog
// Create pricing plan in global catalog. A catalog plan is the data used by the platform for plan metadata definition.
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogPlan(createCatalogPlanOptions *CreateCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateCatalogPlanWithContext(context.Background(), createCatalogPlanOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCatalogPlanWithContext is an alternate form of the CreateCatalogPlan method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogPlanWithContext(ctx context.Context, createCatalogPlanOptions *CreateCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCatalogPlanOptions, "createCatalogPlanOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCatalogPlanOptions, "createCatalogPlanOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *createCatalogPlanOptions.ProductID,
		"catalog_product_id": *createCatalogPlanOptions.CatalogProductID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCatalogPlanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateCatalogPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createCatalogPlanOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*createCatalogPlanOptions.Env))
	}

	body := make(map[string]interface{})
	if createCatalogPlanOptions.Name != nil {
		body["name"] = createCatalogPlanOptions.Name
	}
	if createCatalogPlanOptions.Active != nil {
		body["active"] = createCatalogPlanOptions.Active
	}
	if createCatalogPlanOptions.Disabled != nil {
		body["disabled"] = createCatalogPlanOptions.Disabled
	}
	if createCatalogPlanOptions.Kind != nil {
		body["kind"] = createCatalogPlanOptions.Kind
	}
	if createCatalogPlanOptions.Tags != nil {
		body["tags"] = createCatalogPlanOptions.Tags
	}
	if createCatalogPlanOptions.ObjectProvider != nil {
		body["object_provider"] = createCatalogPlanOptions.ObjectProvider
	}
	if createCatalogPlanOptions.ID != nil {
		body["id"] = createCatalogPlanOptions.ID
	}
	if createCatalogPlanOptions.ObjectID != nil {
		body["object_id"] = createCatalogPlanOptions.ObjectID
	}
	if createCatalogPlanOptions.OverviewUi != nil {
		body["overview_ui"] = createCatalogPlanOptions.OverviewUi
	}
	if createCatalogPlanOptions.Metadata != nil {
		body["metadata"] = createCatalogPlanOptions.Metadata
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_catalog_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogPlan)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogPlan : Get a global catalog pricing plan
// Retrieve the details of a global catalog pricing plan.
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogPlan(getCatalogPlanOptions *GetCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetCatalogPlanWithContext(context.Background(), getCatalogPlanOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCatalogPlanWithContext is an alternate form of the GetCatalogPlan method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogPlanWithContext(ctx context.Context, getCatalogPlanOptions *GetCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogPlanOptions, "getCatalogPlanOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCatalogPlanOptions, "getCatalogPlanOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *getCatalogPlanOptions.ProductID,
		"catalog_product_id": *getCatalogPlanOptions.CatalogProductID,
		"catalog_plan_id":    *getCatalogPlanOptions.CatalogPlanID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCatalogPlanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetCatalogPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCatalogPlanOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*getCatalogPlanOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_catalog_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogPlan)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCatalogPlan : Update a global catalog plan
// Update the details of a global catalog pricing plan.
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogPlan(updateCatalogPlanOptions *UpdateCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateCatalogPlanWithContext(context.Background(), updateCatalogPlanOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCatalogPlanWithContext is an alternate form of the UpdateCatalogPlan method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogPlanWithContext(ctx context.Context, updateCatalogPlanOptions *UpdateCatalogPlanOptions) (result *GlobalCatalogPlan, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCatalogPlanOptions, "updateCatalogPlanOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCatalogPlanOptions, "updateCatalogPlanOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *updateCatalogPlanOptions.ProductID,
		"catalog_product_id": *updateCatalogPlanOptions.CatalogProductID,
		"catalog_plan_id":    *updateCatalogPlanOptions.CatalogPlanID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCatalogPlanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateCatalogPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateCatalogPlanOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*updateCatalogPlanOptions.Env))
	}

	_, err = builder.SetBodyContentJSON(updateCatalogPlanOptions.GlobalCatalogPlanPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_catalog_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogPlan)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCatalogPlan : Delete a global catalog pricing plan
// Delete a global catalog pricing plan.
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogPlan(deleteCatalogPlanOptions *DeleteCatalogPlanOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteCatalogPlanWithContext(context.Background(), deleteCatalogPlanOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCatalogPlanWithContext is an alternate form of the DeleteCatalogPlan method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogPlanWithContext(ctx context.Context, deleteCatalogPlanOptions *DeleteCatalogPlanOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCatalogPlanOptions, "deleteCatalogPlanOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCatalogPlanOptions, "deleteCatalogPlanOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *deleteCatalogPlanOptions.ProductID,
		"catalog_product_id": *deleteCatalogPlanOptions.CatalogProductID,
		"catalog_plan_id":    *deleteCatalogPlanOptions.CatalogPlanID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCatalogPlanOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteCatalogPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteCatalogPlanOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*deleteCatalogPlanOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_catalog_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateCatalogDeployment : Create a global catalog deployment
// Create a global catalog deployment. A catalog deployment is the data used by the platform for deployment metadata
// definition.
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogDeployment(createCatalogDeploymentOptions *CreateCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateCatalogDeploymentWithContext(context.Background(), createCatalogDeploymentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCatalogDeploymentWithContext is an alternate form of the CreateCatalogDeployment method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateCatalogDeploymentWithContext(ctx context.Context, createCatalogDeploymentOptions *CreateCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCatalogDeploymentOptions, "createCatalogDeploymentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCatalogDeploymentOptions, "createCatalogDeploymentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":         *createCatalogDeploymentOptions.ProductID,
		"catalog_product_id": *createCatalogDeploymentOptions.CatalogProductID,
		"catalog_plan_id":    *createCatalogDeploymentOptions.CatalogPlanID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}/catalog_deployments`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCatalogDeploymentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateCatalogDeployment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createCatalogDeploymentOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*createCatalogDeploymentOptions.Env))
	}

	body := make(map[string]interface{})
	if createCatalogDeploymentOptions.Name != nil {
		body["name"] = createCatalogDeploymentOptions.Name
	}
	if createCatalogDeploymentOptions.Active != nil {
		body["active"] = createCatalogDeploymentOptions.Active
	}
	if createCatalogDeploymentOptions.Disabled != nil {
		body["disabled"] = createCatalogDeploymentOptions.Disabled
	}
	if createCatalogDeploymentOptions.Kind != nil {
		body["kind"] = createCatalogDeploymentOptions.Kind
	}
	if createCatalogDeploymentOptions.Tags != nil {
		body["tags"] = createCatalogDeploymentOptions.Tags
	}
	if createCatalogDeploymentOptions.ObjectProvider != nil {
		body["object_provider"] = createCatalogDeploymentOptions.ObjectProvider
	}
	if createCatalogDeploymentOptions.ID != nil {
		body["id"] = createCatalogDeploymentOptions.ID
	}
	if createCatalogDeploymentOptions.ObjectID != nil {
		body["object_id"] = createCatalogDeploymentOptions.ObjectID
	}
	if createCatalogDeploymentOptions.OverviewUi != nil {
		body["overview_ui"] = createCatalogDeploymentOptions.OverviewUi
	}
	if createCatalogDeploymentOptions.Metadata != nil {
		body["metadata"] = createCatalogDeploymentOptions.Metadata
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_catalog_deployment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogDeployment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetCatalogDeployment : Get a global catalog deployment
// Retrieve the details of a global catalog deployment.
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogDeployment(getCatalogDeploymentOptions *GetCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetCatalogDeploymentWithContext(context.Background(), getCatalogDeploymentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCatalogDeploymentWithContext is an alternate form of the GetCatalogDeployment method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetCatalogDeploymentWithContext(ctx context.Context, getCatalogDeploymentOptions *GetCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCatalogDeploymentOptions, "getCatalogDeploymentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCatalogDeploymentOptions, "getCatalogDeploymentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":            *getCatalogDeploymentOptions.ProductID,
		"catalog_product_id":    *getCatalogDeploymentOptions.CatalogProductID,
		"catalog_plan_id":       *getCatalogDeploymentOptions.CatalogPlanID,
		"catalog_deployment_id": *getCatalogDeploymentOptions.CatalogDeploymentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}/catalog_deployments/{catalog_deployment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCatalogDeploymentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetCatalogDeployment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCatalogDeploymentOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*getCatalogDeploymentOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_catalog_deployment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogDeployment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCatalogDeployment : Update a global catalog deployment
// Update the details of a global catalog deployment.
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogDeployment(updateCatalogDeploymentOptions *UpdateCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateCatalogDeploymentWithContext(context.Background(), updateCatalogDeploymentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCatalogDeploymentWithContext is an alternate form of the UpdateCatalogDeployment method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateCatalogDeploymentWithContext(ctx context.Context, updateCatalogDeploymentOptions *UpdateCatalogDeploymentOptions) (result *GlobalCatalogDeployment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCatalogDeploymentOptions, "updateCatalogDeploymentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCatalogDeploymentOptions, "updateCatalogDeploymentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":            *updateCatalogDeploymentOptions.ProductID,
		"catalog_product_id":    *updateCatalogDeploymentOptions.CatalogProductID,
		"catalog_plan_id":       *updateCatalogDeploymentOptions.CatalogPlanID,
		"catalog_deployment_id": *updateCatalogDeploymentOptions.CatalogDeploymentID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}/catalog_deployments/{catalog_deployment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCatalogDeploymentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateCatalogDeployment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateCatalogDeploymentOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*updateCatalogDeploymentOptions.Env))
	}

	_, err = builder.SetBodyContentJSON(updateCatalogDeploymentOptions.GlobalCatalogDeploymentPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_catalog_deployment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGlobalCatalogDeployment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCatalogDeployment : Delete a global catalog deployment
// Delete a global catalog deployment.
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogDeployment(deleteCatalogDeploymentOptions *DeleteCatalogDeploymentOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteCatalogDeploymentWithContext(context.Background(), deleteCatalogDeploymentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCatalogDeploymentWithContext is an alternate form of the DeleteCatalogDeployment method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteCatalogDeploymentWithContext(ctx context.Context, deleteCatalogDeploymentOptions *DeleteCatalogDeploymentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCatalogDeploymentOptions, "deleteCatalogDeploymentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCatalogDeploymentOptions, "deleteCatalogDeploymentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":            *deleteCatalogDeploymentOptions.ProductID,
		"catalog_product_id":    *deleteCatalogDeploymentOptions.CatalogProductID,
		"catalog_plan_id":       *deleteCatalogDeploymentOptions.CatalogPlanID,
		"catalog_deployment_id": *deleteCatalogDeploymentOptions.CatalogDeploymentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/catalog_products/{catalog_product_id}/catalog_plans/{catalog_plan_id}/catalog_deployments/{catalog_deployment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCatalogDeploymentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteCatalogDeployment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteCatalogDeploymentOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*deleteCatalogDeploymentOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_catalog_deployment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateIamRegistration : Create IAM registration for your service
// Create IAM registration based on your product.
func (partnerCenterSell *PartnerCenterSellV1) CreateIamRegistration(createIamRegistrationOptions *CreateIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateIamRegistrationWithContext(context.Background(), createIamRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateIamRegistrationWithContext is an alternate form of the CreateIamRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateIamRegistrationWithContext(ctx context.Context, createIamRegistrationOptions *CreateIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createIamRegistrationOptions, "createIamRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createIamRegistrationOptions, "createIamRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id": *createIamRegistrationOptions.ProductID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/iam_registration`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createIamRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateIamRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createIamRegistrationOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*createIamRegistrationOptions.Env))
	}

	body := make(map[string]interface{})
	if createIamRegistrationOptions.Name != nil {
		body["name"] = createIamRegistrationOptions.Name
	}
	if createIamRegistrationOptions.Enabled != nil {
		body["enabled"] = createIamRegistrationOptions.Enabled
	}
	if createIamRegistrationOptions.ServiceType != nil {
		body["service_type"] = createIamRegistrationOptions.ServiceType
	}
	if createIamRegistrationOptions.Actions != nil {
		body["actions"] = createIamRegistrationOptions.Actions
	}
	if createIamRegistrationOptions.AdditionalPolicyScopes != nil {
		body["additional_policy_scopes"] = createIamRegistrationOptions.AdditionalPolicyScopes
	}
	if createIamRegistrationOptions.DisplayName != nil {
		body["display_name"] = createIamRegistrationOptions.DisplayName
	}
	if createIamRegistrationOptions.ParentIds != nil {
		body["parent_ids"] = createIamRegistrationOptions.ParentIds
	}
	if createIamRegistrationOptions.ResourceHierarchyAttribute != nil {
		body["resource_hierarchy_attribute"] = createIamRegistrationOptions.ResourceHierarchyAttribute
	}
	if createIamRegistrationOptions.SupportedAnonymousAccesses != nil {
		body["supported_anonymous_accesses"] = createIamRegistrationOptions.SupportedAnonymousAccesses
	}
	if createIamRegistrationOptions.SupportedAttributes != nil {
		body["supported_attributes"] = createIamRegistrationOptions.SupportedAttributes
	}
	if createIamRegistrationOptions.SupportedAuthorizationSubjects != nil {
		body["supported_authorization_subjects"] = createIamRegistrationOptions.SupportedAuthorizationSubjects
	}
	if createIamRegistrationOptions.SupportedRoles != nil {
		body["supported_roles"] = createIamRegistrationOptions.SupportedRoles
	}
	if createIamRegistrationOptions.SupportedNetwork != nil {
		body["supported_network"] = createIamRegistrationOptions.SupportedNetwork
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_iam_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalIamServiceRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateIamRegistration : Update IAM registration for your service
// Update your service IAM registration based on your product.
func (partnerCenterSell *PartnerCenterSellV1) UpdateIamRegistration(updateIamRegistrationOptions *UpdateIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateIamRegistrationWithContext(context.Background(), updateIamRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateIamRegistrationWithContext is an alternate form of the UpdateIamRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateIamRegistrationWithContext(ctx context.Context, updateIamRegistrationOptions *UpdateIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateIamRegistrationOptions, "updateIamRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateIamRegistrationOptions, "updateIamRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":        *updateIamRegistrationOptions.ProductID,
		"programmatic_name": *updateIamRegistrationOptions.ProgrammaticName,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/iam_registration/{programmatic_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateIamRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateIamRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateIamRegistrationOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*updateIamRegistrationOptions.Env))
	}

	_, err = builder.SetBodyContentJSON(updateIamRegistrationOptions.IamRegistrationPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_iam_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalIamServiceRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteIamRegistration : Delete IAM registration for your service
// Delete your service IAM registration based on your product.
func (partnerCenterSell *PartnerCenterSellV1) DeleteIamRegistration(deleteIamRegistrationOptions *DeleteIamRegistrationOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteIamRegistrationWithContext(context.Background(), deleteIamRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteIamRegistrationWithContext is an alternate form of the DeleteIamRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteIamRegistrationWithContext(ctx context.Context, deleteIamRegistrationOptions *DeleteIamRegistrationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteIamRegistrationOptions, "deleteIamRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteIamRegistrationOptions, "deleteIamRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":        *deleteIamRegistrationOptions.ProductID,
		"programmatic_name": *deleteIamRegistrationOptions.ProgrammaticName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/iam_registration/{programmatic_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteIamRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteIamRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteIamRegistrationOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*deleteIamRegistrationOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_iam_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetIamRegistration : Get IAM registration for your service
// This method gets your service IAM registration.
func (partnerCenterSell *PartnerCenterSellV1) GetIamRegistration(getIamRegistrationOptions *GetIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetIamRegistrationWithContext(context.Background(), getIamRegistrationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetIamRegistrationWithContext is an alternate form of the GetIamRegistration method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetIamRegistrationWithContext(ctx context.Context, getIamRegistrationOptions *GetIamRegistrationOptions) (result *IamServiceRegistration, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getIamRegistrationOptions, "getIamRegistrationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getIamRegistrationOptions, "getIamRegistrationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"product_id":        *getIamRegistrationOptions.ProductID,
		"programmatic_name": *getIamRegistrationOptions.ProgrammaticName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/products/{product_id}/iam_registration/{programmatic_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getIamRegistrationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetIamRegistration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getIamRegistrationOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*getIamRegistrationOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_iam_registration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalIamServiceRegistration)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateResourceBroker : Create a broker
// Create a new broker that manages the lifecycle of your service and its metering integration.
func (partnerCenterSell *PartnerCenterSellV1) CreateResourceBroker(createResourceBrokerOptions *CreateResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.CreateResourceBrokerWithContext(context.Background(), createResourceBrokerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateResourceBrokerWithContext is an alternate form of the CreateResourceBroker method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) CreateResourceBrokerWithContext(ctx context.Context, createResourceBrokerOptions *CreateResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createResourceBrokerOptions, "createResourceBrokerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createResourceBrokerOptions, "createResourceBrokerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/brokers`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createResourceBrokerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "CreateResourceBroker")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createResourceBrokerOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*createResourceBrokerOptions.Env))
	}

	body := make(map[string]interface{})
	if createResourceBrokerOptions.AuthScheme != nil {
		body["auth_scheme"] = createResourceBrokerOptions.AuthScheme
	}
	if createResourceBrokerOptions.Name != nil {
		body["name"] = createResourceBrokerOptions.Name
	}
	if createResourceBrokerOptions.BrokerURL != nil {
		body["broker_url"] = createResourceBrokerOptions.BrokerURL
	}
	if createResourceBrokerOptions.Type != nil {
		body["type"] = createResourceBrokerOptions.Type
	}
	if createResourceBrokerOptions.AuthUsername != nil {
		body["auth_username"] = createResourceBrokerOptions.AuthUsername
	}
	if createResourceBrokerOptions.AuthPassword != nil {
		body["auth_password"] = createResourceBrokerOptions.AuthPassword
	}
	if createResourceBrokerOptions.ResourceGroupCrn != nil {
		body["resource_group_crn"] = createResourceBrokerOptions.ResourceGroupCrn
	}
	if createResourceBrokerOptions.State != nil {
		body["state"] = createResourceBrokerOptions.State
	}
	if createResourceBrokerOptions.AllowContextUpdates != nil {
		body["allow_context_updates"] = createResourceBrokerOptions.AllowContextUpdates
	}
	if createResourceBrokerOptions.CatalogType != nil {
		body["catalog_type"] = createResourceBrokerOptions.CatalogType
	}
	if createResourceBrokerOptions.Region != nil {
		body["region"] = createResourceBrokerOptions.Region
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_resource_broker", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBroker)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateResourceBroker : Update broker details
// Update your service broker details, including the username, password, name, or the URL of the broker.
func (partnerCenterSell *PartnerCenterSellV1) UpdateResourceBroker(updateResourceBrokerOptions *UpdateResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.UpdateResourceBrokerWithContext(context.Background(), updateResourceBrokerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateResourceBrokerWithContext is an alternate form of the UpdateResourceBroker method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) UpdateResourceBrokerWithContext(ctx context.Context, updateResourceBrokerOptions *UpdateResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateResourceBrokerOptions, "updateResourceBrokerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateResourceBrokerOptions, "updateResourceBrokerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"broker_id": *updateResourceBrokerOptions.BrokerID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/brokers/{broker_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateResourceBrokerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "UpdateResourceBroker")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateResourceBrokerOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*updateResourceBrokerOptions.Env))
	}

	_, err = builder.SetBodyContentJSON(updateResourceBrokerOptions.BrokerPatch)
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_resource_broker", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBroker)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetResourceBroker : Get a broker
// Get a resource broker.
func (partnerCenterSell *PartnerCenterSellV1) GetResourceBroker(getResourceBrokerOptions *GetResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetResourceBrokerWithContext(context.Background(), getResourceBrokerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetResourceBrokerWithContext is an alternate form of the GetResourceBroker method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetResourceBrokerWithContext(ctx context.Context, getResourceBrokerOptions *GetResourceBrokerOptions) (result *Broker, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getResourceBrokerOptions, "getResourceBrokerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getResourceBrokerOptions, "getResourceBrokerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"broker_id": *getResourceBrokerOptions.BrokerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/brokers/{broker_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getResourceBrokerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetResourceBroker")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getResourceBrokerOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*getResourceBrokerOptions.Env))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_resource_broker", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBroker)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteResourceBroker : Remove a broker
// Remove a broker from the account.
func (partnerCenterSell *PartnerCenterSellV1) DeleteResourceBroker(deleteResourceBrokerOptions *DeleteResourceBrokerOptions) (response *core.DetailedResponse, err error) {
	response, err = partnerCenterSell.DeleteResourceBrokerWithContext(context.Background(), deleteResourceBrokerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteResourceBrokerWithContext is an alternate form of the DeleteResourceBroker method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) DeleteResourceBrokerWithContext(ctx context.Context, deleteResourceBrokerOptions *DeleteResourceBrokerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteResourceBrokerOptions, "deleteResourceBrokerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteResourceBrokerOptions, "deleteResourceBrokerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"broker_id": *deleteResourceBrokerOptions.BrokerID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/brokers/{broker_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteResourceBrokerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "DeleteResourceBroker")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteResourceBrokerOptions.Env != nil {
		builder.AddQuery("env", fmt.Sprint(*deleteResourceBrokerOptions.Env))
	}
	if deleteResourceBrokerOptions.RemoveFromAccount != nil {
		builder.AddQuery("remove_from_account", fmt.Sprint(*deleteResourceBrokerOptions.RemoveFromAccount))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = partnerCenterSell.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_resource_broker", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListProductBadges : List badges
// List all available badges that a product can be validated against.
func (partnerCenterSell *PartnerCenterSellV1) ListProductBadges(listProductBadgesOptions *ListProductBadgesOptions) (result *ProductBadgeCollection, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.ListProductBadgesWithContext(context.Background(), listProductBadgesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProductBadgesWithContext is an alternate form of the ListProductBadges method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) ListProductBadgesWithContext(ctx context.Context, listProductBadgesOptions *ListProductBadgesOptions) (result *ProductBadgeCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProductBadgesOptions, "listProductBadgesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/product_badges`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listProductBadgesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "ListProductBadges")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProductBadgesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProductBadgesOptions.Limit))
	}
	if listProductBadgesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listProductBadgesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_product_badges", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProductBadgeCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProductBadge : Get badge
// Retrieve the details of a badge.
func (partnerCenterSell *PartnerCenterSellV1) GetProductBadge(getProductBadgeOptions *GetProductBadgeOptions) (result *ProductBadge, response *core.DetailedResponse, err error) {
	result, response, err = partnerCenterSell.GetProductBadgeWithContext(context.Background(), getProductBadgeOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProductBadgeWithContext is an alternate form of the GetProductBadge method which supports a Context parameter
func (partnerCenterSell *PartnerCenterSellV1) GetProductBadgeWithContext(ctx context.Context, getProductBadgeOptions *GetProductBadgeOptions) (result *ProductBadge, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProductBadgeOptions, "getProductBadgeOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProductBadgeOptions, "getProductBadgeOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"badge_id": fmt.Sprint(*getProductBadgeOptions.BadgeID),
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = partnerCenterSell.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(partnerCenterSell.Service.Options.URL, `/product_badges/{badge_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProductBadgeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("partner_center_sell", "V1", "GetProductBadge")
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
	response, err = partnerCenterSell.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_product_badge", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProductBadge)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.5.0")
}

// Bookmark : The page reference information.
type Bookmark struct {
	// The URL of the next or previous page.
	Href *string `json:"href,omitempty"`

	// The reference ID of the first item on the page.
	Start *strfmt.UUID `json:"start,omitempty"`
}

// UnmarshalBookmark unmarshals an instance of Bookmark from the specified map of raw messages.
func UnmarshalBookmark(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Bookmark)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Broker : A resource broker registration in the IBM Cloud resource controller system.
type Broker struct {
	// The authentication username to reach the broker.
	AuthUsername *string `json:"auth_username,omitempty"`

	// The authentication password to reach the broker.
	AuthPassword *string `json:"auth_password,omitempty"`

	// The supported authentication scheme for the broker.
	AuthScheme *string `json:"auth_scheme,omitempty"`

	// The cloud resource name of the resource group.
	ResourceGroupCrn *string `json:"resource_group_crn,omitempty"`

	// The state of the broker.
	State *string `json:"state,omitempty"`

	// The URL associated with the broker application.
	BrokerURL *string `json:"broker_url,omitempty"`

	// Whether the resource controller will call the broker for any context changes to the instance. Currently, the only
	// context related change is an instance name update.
	AllowContextUpdates *bool `json:"allow_context_updates,omitempty"`

	// To enable the provisioning of your broker, set this parameter value to `service`.
	CatalogType *string `json:"catalog_type,omitempty"`

	// The type of the provisioning model.
	Type *string `json:"type,omitempty"`

	// The name of the broker.
	Name *string `json:"name,omitempty"`

	// The region where the pricing plan is available.
	Region *string `json:"region,omitempty"`

	// The ID of the account in which you manage the broker.
	AccountID *string `json:"account_id,omitempty"`

	// The cloud resource name (CRN) of the broker.
	Crn *string `json:"crn,omitempty"`

	// The time when the service broker was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// The time when the service broker was updated.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// The time when the service broker was deleted.
	DeletedAt *strfmt.DateTime `json:"deleted_at,omitempty"`

	// The details of the user who created this broker.
	CreatedBy *BrokerEventCreatedByUser `json:"created_by,omitempty"`

	// The details of the user who updated this broker.
	UpdatedBy *BrokerEventUpdatedByUser `json:"updated_by,omitempty"`

	// The details of the user who deleted this broker.
	DeletedBy *BrokerEventDeletedByUser `json:"deleted_by,omitempty"`

	// The globally unique identifier of the broker.
	Guid *string `json:"guid,omitempty"`

	// The identifier of the broker.
	ID *string `json:"id,omitempty"`

	// The URL associated with the broker.
	URL *string `json:"url,omitempty"`
}

// Constants associated with the Broker.AuthUsername property.
// The authentication username to reach the broker.
const (
	Broker_AuthUsername_Apikey = "apikey"
)

// Constants associated with the Broker.AuthScheme property.
// The supported authentication scheme for the broker.
const (
	Broker_AuthScheme_Bearer    = "bearer"
	Broker_AuthScheme_BearerCrn = "bearer-crn"
)

// Constants associated with the Broker.State property.
// The state of the broker.
const (
	Broker_State_Active  = "active"
	Broker_State_Removed = "removed"
)

// Constants associated with the Broker.Type property.
// The type of the provisioning model.
const (
	Broker_Type_ProvisionBehind  = "provision_behind"
	Broker_Type_ProvisionThrough = "provision_through"
)

// UnmarshalBroker unmarshals an instance of Broker from the specified map of raw messages.
func UnmarshalBroker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Broker)
	err = core.UnmarshalPrimitive(m, "auth_username", &obj.AuthUsername)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_password", &obj.AuthPassword)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_scheme", &obj.AuthScheme)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_scheme-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_crn", &obj.ResourceGroupCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "broker_url", &obj.BrokerURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_context_updates", &obj.AllowContextUpdates)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_context_updates-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_type", &obj.CatalogType)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deleted_at", &obj.DeletedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "deleted_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "created_by", &obj.CreatedBy, UnmarshalBrokerEventCreatedByUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "updated_by", &obj.UpdatedBy, UnmarshalBrokerEventUpdatedByUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deleted_by", &obj.DeletedBy, UnmarshalBrokerEventDeletedByUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "deleted_by-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "guid", &obj.Guid)
	if err != nil {
		err = core.SDKErrorf(err, "", "guid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BrokerEventCreatedByUser : The details of the user who created this broker.
type BrokerEventCreatedByUser struct {
	// The ID of the user who dispatched this action.
	UserID *string `json:"user_id,omitempty"`

	// The username of the user who dispatched this action.
	UserName *string `json:"user_name,omitempty"`
}

// UnmarshalBrokerEventCreatedByUser unmarshals an instance of BrokerEventCreatedByUser from the specified map of raw messages.
func UnmarshalBrokerEventCreatedByUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrokerEventCreatedByUser)
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_name", &obj.UserName)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BrokerEventDeletedByUser : The details of the user who deleted this broker.
type BrokerEventDeletedByUser struct {
	// The ID of the user who dispatched this action.
	UserID *string `json:"user_id,omitempty"`

	// The username of the user who dispatched this action.
	UserName *string `json:"user_name,omitempty"`
}

// UnmarshalBrokerEventDeletedByUser unmarshals an instance of BrokerEventDeletedByUser from the specified map of raw messages.
func UnmarshalBrokerEventDeletedByUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrokerEventDeletedByUser)
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_name", &obj.UserName)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BrokerEventUpdatedByUser : The details of the user who updated this broker.
type BrokerEventUpdatedByUser struct {
	// The ID of the user who dispatched this action.
	UserID *string `json:"user_id,omitempty"`

	// The username of the user who dispatched this action.
	UserName *string `json:"user_name,omitempty"`
}

// UnmarshalBrokerEventUpdatedByUser unmarshals an instance of BrokerEventUpdatedByUser from the specified map of raw messages.
func UnmarshalBrokerEventUpdatedByUser(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrokerEventUpdatedByUser)
	err = core.UnmarshalPrimitive(m, "user_id", &obj.UserID)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_name", &obj.UserName)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BrokerPatch : Request body for updating a Resource Controller broker registration.
type BrokerPatch struct {
	// The authentication username to reach the broker.
	AuthUsername *string `json:"auth_username,omitempty"`

	// The authentication password to reach the broker.
	AuthPassword *string `json:"auth_password,omitempty"`

	// The supported authentication scheme for the broker.
	AuthScheme *string `json:"auth_scheme,omitempty"`

	// The cloud resource name of the resource group.
	ResourceGroupCrn *string `json:"resource_group_crn,omitempty"`

	// The state of the broker.
	State *string `json:"state,omitempty"`

	// The URL associated with the broker application.
	BrokerURL *string `json:"broker_url,omitempty"`

	// Whether the resource controller will call the broker for any context changes to the instance. Currently, the only
	// context related change is an instance name update.
	AllowContextUpdates *bool `json:"allow_context_updates,omitempty"`

	// To enable the provisioning of your broker, set this parameter value to `service`.
	CatalogType *string `json:"catalog_type,omitempty"`

	// The type of the provisioning model.
	Type *string `json:"type,omitempty"`

	// The region where the pricing plan is available.
	Region *string `json:"region,omitempty"`
}

// Constants associated with the BrokerPatch.AuthUsername property.
// The authentication username to reach the broker.
const (
	BrokerPatch_AuthUsername_Apikey = "apikey"
)

// Constants associated with the BrokerPatch.AuthScheme property.
// The supported authentication scheme for the broker.
const (
	BrokerPatch_AuthScheme_Bearer    = "bearer"
	BrokerPatch_AuthScheme_BearerCrn = "bearer-crn"
)

// Constants associated with the BrokerPatch.State property.
// The state of the broker.
const (
	BrokerPatch_State_Active  = "active"
	BrokerPatch_State_Removed = "removed"
)

// Constants associated with the BrokerPatch.Type property.
// The type of the provisioning model.
const (
	BrokerPatch_Type_ProvisionBehind  = "provision_behind"
	BrokerPatch_Type_ProvisionThrough = "provision_through"
)

// UnmarshalBrokerPatch unmarshals an instance of BrokerPatch from the specified map of raw messages.
func UnmarshalBrokerPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrokerPatch)
	err = core.UnmarshalPrimitive(m, "auth_username", &obj.AuthUsername)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_password", &obj.AuthPassword)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_scheme", &obj.AuthScheme)
	if err != nil {
		err = core.SDKErrorf(err, "", "auth_scheme-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_crn", &obj.ResourceGroupCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "broker_url", &obj.BrokerURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_context_updates", &obj.AllowContextUpdates)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_context_updates-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_type", &obj.CatalogType)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the BrokerPatch
func (brokerPatch *BrokerPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(brokerPatch.AuthUsername) {
		_patch["auth_username"] = brokerPatch.AuthUsername
	}
	if !core.IsNil(brokerPatch.AuthPassword) {
		_patch["auth_password"] = brokerPatch.AuthPassword
	}
	if !core.IsNil(brokerPatch.AuthScheme) {
		_patch["auth_scheme"] = brokerPatch.AuthScheme
	}
	if !core.IsNil(brokerPatch.ResourceGroupCrn) {
		_patch["resource_group_crn"] = brokerPatch.ResourceGroupCrn
	}
	if !core.IsNil(brokerPatch.State) {
		_patch["state"] = brokerPatch.State
	}
	if !core.IsNil(brokerPatch.BrokerURL) {
		_patch["broker_url"] = brokerPatch.BrokerURL
	}
	if !core.IsNil(brokerPatch.AllowContextUpdates) {
		_patch["allow_context_updates"] = brokerPatch.AllowContextUpdates
	}
	if !core.IsNil(brokerPatch.CatalogType) {
		_patch["catalog_type"] = brokerPatch.CatalogType
	}
	if !core.IsNil(brokerPatch.Type) {
		_patch["type"] = brokerPatch.Type
	}
	if !core.IsNil(brokerPatch.Region) {
		_patch["region"] = brokerPatch.Region
	}

	return
}

// CatalogHighlightItem : The attributes of the product that differentiate it in the market.
type CatalogHighlightItem struct {
	// The description about the features of the product.
	Description *string `json:"description,omitempty"`

	// The description about the features of the product in translation.
	DescriptionI18n map[string]string `json:"description_i18n,omitempty"`

	// The descriptive title for the feature.
	Title *string `json:"title,omitempty"`

	// The descriptive title for the feature in translation.
	TitleI18n map[string]string `json:"title_i18n,omitempty"`
}

// UnmarshalCatalogHighlightItem unmarshals an instance of CatalogHighlightItem from the specified map of raw messages.
func UnmarshalCatalogHighlightItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogHighlightItem)
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description_i18n", &obj.DescriptionI18n)
	if err != nil {
		err = core.SDKErrorf(err, "", "description_i18n-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "title", &obj.Title)
	if err != nil {
		err = core.SDKErrorf(err, "", "title-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "title_i18n", &obj.TitleI18n)
	if err != nil {
		err = core.SDKErrorf(err, "", "title_i18n-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the CatalogHighlightItem
func (catalogHighlightItem *CatalogHighlightItem) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(catalogHighlightItem.Description) {
		_patch["description"] = catalogHighlightItem.Description
	}
	if !core.IsNil(catalogHighlightItem.DescriptionI18n) {
		_patch["description_i18n"] = catalogHighlightItem.DescriptionI18n
	}
	if !core.IsNil(catalogHighlightItem.Title) {
		_patch["title"] = catalogHighlightItem.Title
	}
	if !core.IsNil(catalogHighlightItem.TitleI18n) {
		_patch["title_i18n"] = catalogHighlightItem.TitleI18n
	}

	return
}

// CatalogProductMediaItem : CatalogProductMediaItem struct
type CatalogProductMediaItem struct {
	// Provide a descriptive caption that indicates what the media illustrates. This caption is displayed in the catalog.
	Caption *string `json:"caption" validate:"required"`

	// The brief explanation for your images and videos in translation.
	CaptionI18n map[string]string `json:"caption_i18n,omitempty"`

	// The reduced-size version of your images and videos.
	Thumbnail *string `json:"thumbnail,omitempty"`

	// The type of the media.
	Type *string `json:"type" validate:"required"`

	// The URL that links to the media that shows off the product.
	URL *string `json:"url" validate:"required"`
}

// Constants associated with the CatalogProductMediaItem.Type property.
// The type of the media.
const (
	CatalogProductMediaItem_Type_Image     = "image"
	CatalogProductMediaItem_Type_VideoMp4  = "video_mp_4"
	CatalogProductMediaItem_Type_VideoWebm = "video_webm"
	CatalogProductMediaItem_Type_Youtube   = "youtube"
)

// NewCatalogProductMediaItem : Instantiate CatalogProductMediaItem (Generic Model Constructor)
func (*PartnerCenterSellV1) NewCatalogProductMediaItem(caption string, typeVar string, url string) (_model *CatalogProductMediaItem, err error) {
	_model = &CatalogProductMediaItem{
		Caption: core.StringPtr(caption),
		Type:    core.StringPtr(typeVar),
		URL:     core.StringPtr(url),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCatalogProductMediaItem unmarshals an instance of CatalogProductMediaItem from the specified map of raw messages.
func UnmarshalCatalogProductMediaItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogProductMediaItem)
	err = core.UnmarshalPrimitive(m, "caption", &obj.Caption)
	if err != nil {
		err = core.SDKErrorf(err, "", "caption-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "caption_i18n", &obj.CaptionI18n)
	if err != nil {
		err = core.SDKErrorf(err, "", "caption_i18n-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "thumbnail", &obj.Thumbnail)
	if err != nil {
		err = core.SDKErrorf(err, "", "thumbnail-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the CatalogProductMediaItem
func (catalogProductMediaItem *CatalogProductMediaItem) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(catalogProductMediaItem.Caption) {
		_patch["caption"] = catalogProductMediaItem.Caption
	}
	if !core.IsNil(catalogProductMediaItem.CaptionI18n) {
		_patch["caption_i18n"] = catalogProductMediaItem.CaptionI18n
	}
	if !core.IsNil(catalogProductMediaItem.Thumbnail) {
		_patch["thumbnail"] = catalogProductMediaItem.Thumbnail
	}
	if !core.IsNil(catalogProductMediaItem.Type) {
		_patch["type"] = catalogProductMediaItem.Type
	}
	if !core.IsNil(catalogProductMediaItem.URL) {
		_patch["url"] = catalogProductMediaItem.URL
	}

	return
}

// CatalogProductProvider : The provider or owner of the product.
type CatalogProductProvider struct {
	// The name of the provider.
	Name *string `json:"name,omitempty"`

	// The email address of the provider.
	Email *string `json:"email,omitempty"`
}

// UnmarshalCatalogProductProvider unmarshals an instance of CatalogProductProvider from the specified map of raw messages.
func UnmarshalCatalogProductProvider(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CatalogProductProvider)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the CatalogProductProvider
func (catalogProductProvider *CatalogProductProvider) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(catalogProductProvider.Name) {
		_patch["name"] = catalogProductProvider.Name
	}
	if !core.IsNil(catalogProductProvider.Email) {
		_patch["email"] = catalogProductProvider.Email
	}

	return
}

// CreateCatalogDeploymentOptions : The CreateCatalogDeployment options.
type CreateCatalogDeploymentOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The programmatic name of this deployment.
	Name *string `json:"name" validate:"required"`

	// Whether the service is active.
	Active *bool `json:"active" validate:"required"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// The kind of the global catalog object.
	Kind *string `json:"kind" validate:"required"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags" validate:"required"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider" validate:"required"`

	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// Global catalog deployment metadata.
	Metadata *GlobalCatalogDeploymentMetadata `json:"metadata,omitempty"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateCatalogDeploymentOptions.Kind property.
// The kind of the global catalog object.
const (
	CreateCatalogDeploymentOptions_Kind_Deployment = "deployment"
)

// NewCreateCatalogDeploymentOptions : Instantiate CreateCatalogDeploymentOptions
func (*PartnerCenterSellV1) NewCreateCatalogDeploymentOptions(productID string, catalogProductID string, catalogPlanID string, name string, active bool, disabled bool, kind string, tags []string, objectProvider *CatalogProductProvider) *CreateCatalogDeploymentOptions {
	return &CreateCatalogDeploymentOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
		CatalogPlanID:    core.StringPtr(catalogPlanID),
		Name:             core.StringPtr(name),
		Active:           core.BoolPtr(active),
		Disabled:         core.BoolPtr(disabled),
		Kind:             core.StringPtr(kind),
		Tags:             tags,
		ObjectProvider:   objectProvider,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *CreateCatalogDeploymentOptions) SetProductID(productID string) *CreateCatalogDeploymentOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *CreateCatalogDeploymentOptions) SetCatalogProductID(catalogProductID string) *CreateCatalogDeploymentOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *CreateCatalogDeploymentOptions) SetCatalogPlanID(catalogPlanID string) *CreateCatalogDeploymentOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateCatalogDeploymentOptions) SetName(name string) *CreateCatalogDeploymentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetActive : Allow user to set Active
func (_options *CreateCatalogDeploymentOptions) SetActive(active bool) *CreateCatalogDeploymentOptions {
	_options.Active = core.BoolPtr(active)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateCatalogDeploymentOptions) SetDisabled(disabled bool) *CreateCatalogDeploymentOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateCatalogDeploymentOptions) SetKind(kind string) *CreateCatalogDeploymentOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCatalogDeploymentOptions) SetTags(tags []string) *CreateCatalogDeploymentOptions {
	_options.Tags = tags
	return _options
}

// SetObjectProvider : Allow user to set ObjectProvider
func (_options *CreateCatalogDeploymentOptions) SetObjectProvider(objectProvider *CatalogProductProvider) *CreateCatalogDeploymentOptions {
	_options.ObjectProvider = objectProvider
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateCatalogDeploymentOptions) SetID(id string) *CreateCatalogDeploymentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetObjectID : Allow user to set ObjectID
func (_options *CreateCatalogDeploymentOptions) SetObjectID(objectID string) *CreateCatalogDeploymentOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetOverviewUi : Allow user to set OverviewUi
func (_options *CreateCatalogDeploymentOptions) SetOverviewUi(overviewUi *GlobalCatalogOverviewUI) *CreateCatalogDeploymentOptions {
	_options.OverviewUi = overviewUi
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateCatalogDeploymentOptions) SetMetadata(metadata *GlobalCatalogDeploymentMetadata) *CreateCatalogDeploymentOptions {
	_options.Metadata = metadata
	return _options
}

// SetEnv : Allow user to set Env
func (_options *CreateCatalogDeploymentOptions) SetEnv(env string) *CreateCatalogDeploymentOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogDeploymentOptions) SetHeaders(param map[string]string) *CreateCatalogDeploymentOptions {
	options.Headers = param
	return options
}

// CreateCatalogPlanOptions : The CreateCatalogPlan options.
type CreateCatalogPlanOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The programmatic name of this plan.
	Name *string `json:"name" validate:"required"`

	// Whether the service is active.
	Active *bool `json:"active" validate:"required"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// The kind of the global catalog object.
	Kind *string `json:"kind" validate:"required"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags" validate:"required"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider" validate:"required"`

	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// Global catalog plan metadata.
	Metadata *GlobalCatalogPlanMetadata `json:"metadata,omitempty"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateCatalogPlanOptions.Kind property.
// The kind of the global catalog object.
const (
	CreateCatalogPlanOptions_Kind_Plan = "plan"
)

// NewCreateCatalogPlanOptions : Instantiate CreateCatalogPlanOptions
func (*PartnerCenterSellV1) NewCreateCatalogPlanOptions(productID string, catalogProductID string, name string, active bool, disabled bool, kind string, tags []string, objectProvider *CatalogProductProvider) *CreateCatalogPlanOptions {
	return &CreateCatalogPlanOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
		Name:             core.StringPtr(name),
		Active:           core.BoolPtr(active),
		Disabled:         core.BoolPtr(disabled),
		Kind:             core.StringPtr(kind),
		Tags:             tags,
		ObjectProvider:   objectProvider,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *CreateCatalogPlanOptions) SetProductID(productID string) *CreateCatalogPlanOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *CreateCatalogPlanOptions) SetCatalogProductID(catalogProductID string) *CreateCatalogPlanOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateCatalogPlanOptions) SetName(name string) *CreateCatalogPlanOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetActive : Allow user to set Active
func (_options *CreateCatalogPlanOptions) SetActive(active bool) *CreateCatalogPlanOptions {
	_options.Active = core.BoolPtr(active)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateCatalogPlanOptions) SetDisabled(disabled bool) *CreateCatalogPlanOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateCatalogPlanOptions) SetKind(kind string) *CreateCatalogPlanOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCatalogPlanOptions) SetTags(tags []string) *CreateCatalogPlanOptions {
	_options.Tags = tags
	return _options
}

// SetObjectProvider : Allow user to set ObjectProvider
func (_options *CreateCatalogPlanOptions) SetObjectProvider(objectProvider *CatalogProductProvider) *CreateCatalogPlanOptions {
	_options.ObjectProvider = objectProvider
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateCatalogPlanOptions) SetID(id string) *CreateCatalogPlanOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetObjectID : Allow user to set ObjectID
func (_options *CreateCatalogPlanOptions) SetObjectID(objectID string) *CreateCatalogPlanOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetOverviewUi : Allow user to set OverviewUi
func (_options *CreateCatalogPlanOptions) SetOverviewUi(overviewUi *GlobalCatalogOverviewUI) *CreateCatalogPlanOptions {
	_options.OverviewUi = overviewUi
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateCatalogPlanOptions) SetMetadata(metadata *GlobalCatalogPlanMetadata) *CreateCatalogPlanOptions {
	_options.Metadata = metadata
	return _options
}

// SetEnv : Allow user to set Env
func (_options *CreateCatalogPlanOptions) SetEnv(env string) *CreateCatalogPlanOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogPlanOptions) SetHeaders(param map[string]string) *CreateCatalogPlanOptions {
	options.Headers = param
	return options
}

// CreateCatalogProductOptions : The CreateCatalogProduct options.
type CreateCatalogProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The programmatic name of this product.
	Name *string `json:"name" validate:"required"`

	// Whether the service is active.
	Active *bool `json:"active" validate:"required"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// The kind of the global catalog object.
	Kind *string `json:"kind" validate:"required"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags" validate:"required"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider" validate:"required"`

	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// Images from the global catalog entry that help illustrate the service.
	Images *GlobalCatalogProductImages `json:"images,omitempty"`

	// The global catalog service metadata object.
	Metadata *GlobalCatalogProductMetadata `json:"metadata,omitempty"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateCatalogProductOptions.Kind property.
// The kind of the global catalog object.
const (
	CreateCatalogProductOptions_Kind_Composite       = "composite"
	CreateCatalogProductOptions_Kind_PlatformService = "platform_service"
	CreateCatalogProductOptions_Kind_Service         = "service"
)

// NewCreateCatalogProductOptions : Instantiate CreateCatalogProductOptions
func (*PartnerCenterSellV1) NewCreateCatalogProductOptions(productID string, name string, active bool, disabled bool, kind string, tags []string, objectProvider *CatalogProductProvider) *CreateCatalogProductOptions {
	return &CreateCatalogProductOptions{
		ProductID:      core.StringPtr(productID),
		Name:           core.StringPtr(name),
		Active:         core.BoolPtr(active),
		Disabled:       core.BoolPtr(disabled),
		Kind:           core.StringPtr(kind),
		Tags:           tags,
		ObjectProvider: objectProvider,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *CreateCatalogProductOptions) SetProductID(productID string) *CreateCatalogProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateCatalogProductOptions) SetName(name string) *CreateCatalogProductOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetActive : Allow user to set Active
func (_options *CreateCatalogProductOptions) SetActive(active bool) *CreateCatalogProductOptions {
	_options.Active = core.BoolPtr(active)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *CreateCatalogProductOptions) SetDisabled(disabled bool) *CreateCatalogProductOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetKind : Allow user to set Kind
func (_options *CreateCatalogProductOptions) SetKind(kind string) *CreateCatalogProductOptions {
	_options.Kind = core.StringPtr(kind)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCatalogProductOptions) SetTags(tags []string) *CreateCatalogProductOptions {
	_options.Tags = tags
	return _options
}

// SetObjectProvider : Allow user to set ObjectProvider
func (_options *CreateCatalogProductOptions) SetObjectProvider(objectProvider *CatalogProductProvider) *CreateCatalogProductOptions {
	_options.ObjectProvider = objectProvider
	return _options
}

// SetID : Allow user to set ID
func (_options *CreateCatalogProductOptions) SetID(id string) *CreateCatalogProductOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetObjectID : Allow user to set ObjectID
func (_options *CreateCatalogProductOptions) SetObjectID(objectID string) *CreateCatalogProductOptions {
	_options.ObjectID = core.StringPtr(objectID)
	return _options
}

// SetOverviewUi : Allow user to set OverviewUi
func (_options *CreateCatalogProductOptions) SetOverviewUi(overviewUi *GlobalCatalogOverviewUI) *CreateCatalogProductOptions {
	_options.OverviewUi = overviewUi
	return _options
}

// SetImages : Allow user to set Images
func (_options *CreateCatalogProductOptions) SetImages(images *GlobalCatalogProductImages) *CreateCatalogProductOptions {
	_options.Images = images
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *CreateCatalogProductOptions) SetMetadata(metadata *GlobalCatalogProductMetadata) *CreateCatalogProductOptions {
	_options.Metadata = metadata
	return _options
}

// SetEnv : Allow user to set Env
func (_options *CreateCatalogProductOptions) SetEnv(env string) *CreateCatalogProductOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCatalogProductOptions) SetHeaders(param map[string]string) *CreateCatalogProductOptions {
	options.Headers = param
	return options
}

// CreateIamRegistrationOptions : The CreateIamRegistration options.
type CreateIamRegistrationOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The IAM registration name, which must be the programmatic name of the product.
	Name *string `json:"name" validate:"required"`

	// Whether the service is enabled or disabled for IAM.
	Enabled *bool `json:"enabled,omitempty"`

	// The type of the service.
	ServiceType *string `json:"service_type,omitempty"`

	// The product access management action.
	Actions []IamServiceRegistrationAction `json:"actions,omitempty"`

	// List of additional policy scopes.
	AdditionalPolicyScopes []string `json:"additional_policy_scopes,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// The list of parent IDs for product access management.
	ParentIds []string `json:"parent_ids,omitempty"`

	// The resource hierarchy key-value pair for composite services.
	ResourceHierarchyAttribute *IamServiceRegistrationResourceHierarchyAttribute `json:"resource_hierarchy_attribute,omitempty"`

	// The list of supported anonymous accesses.
	SupportedAnonymousAccesses []IamServiceRegistrationSupportedAnonymousAccess `json:"supported_anonymous_accesses,omitempty"`

	// The list of supported attributes.
	SupportedAttributes []IamServiceRegistrationSupportedAttribute `json:"supported_attributes,omitempty"`

	// The list of supported authorization subjects.
	SupportedAuthorizationSubjects []IamServiceRegistrationSupportedAuthorizationSubject `json:"supported_authorization_subjects,omitempty"`

	// The list of roles that you can use to assign access.
	SupportedRoles []IamServiceRegistrationSupportedRole `json:"supported_roles,omitempty"`

	// The registration of set of endpoint types that are supported by your service in the `networkType` environment
	// attribute. This constrains the context-based restriction rules specific to the service such that they describe
	// access restrictions on only this set of endpoints.
	SupportedNetwork *IamServiceRegistrationSupportedNetwork `json:"supported_network,omitempty"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateIamRegistrationOptions.ServiceType property.
// The type of the service.
const (
	CreateIamRegistrationOptions_ServiceType_PlatformService = "platform_service"
	CreateIamRegistrationOptions_ServiceType_Service         = "service"
)

// NewCreateIamRegistrationOptions : Instantiate CreateIamRegistrationOptions
func (*PartnerCenterSellV1) NewCreateIamRegistrationOptions(productID string, name string) *CreateIamRegistrationOptions {
	return &CreateIamRegistrationOptions{
		ProductID: core.StringPtr(productID),
		Name:      core.StringPtr(name),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *CreateIamRegistrationOptions) SetProductID(productID string) *CreateIamRegistrationOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateIamRegistrationOptions) SetName(name string) *CreateIamRegistrationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateIamRegistrationOptions) SetEnabled(enabled bool) *CreateIamRegistrationOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetServiceType : Allow user to set ServiceType
func (_options *CreateIamRegistrationOptions) SetServiceType(serviceType string) *CreateIamRegistrationOptions {
	_options.ServiceType = core.StringPtr(serviceType)
	return _options
}

// SetActions : Allow user to set Actions
func (_options *CreateIamRegistrationOptions) SetActions(actions []IamServiceRegistrationAction) *CreateIamRegistrationOptions {
	_options.Actions = actions
	return _options
}

// SetAdditionalPolicyScopes : Allow user to set AdditionalPolicyScopes
func (_options *CreateIamRegistrationOptions) SetAdditionalPolicyScopes(additionalPolicyScopes []string) *CreateIamRegistrationOptions {
	_options.AdditionalPolicyScopes = additionalPolicyScopes
	return _options
}

// SetDisplayName : Allow user to set DisplayName
func (_options *CreateIamRegistrationOptions) SetDisplayName(displayName *IamServiceRegistrationDisplayNameObject) *CreateIamRegistrationOptions {
	_options.DisplayName = displayName
	return _options
}

// SetParentIds : Allow user to set ParentIds
func (_options *CreateIamRegistrationOptions) SetParentIds(parentIds []string) *CreateIamRegistrationOptions {
	_options.ParentIds = parentIds
	return _options
}

// SetResourceHierarchyAttribute : Allow user to set ResourceHierarchyAttribute
func (_options *CreateIamRegistrationOptions) SetResourceHierarchyAttribute(resourceHierarchyAttribute *IamServiceRegistrationResourceHierarchyAttribute) *CreateIamRegistrationOptions {
	_options.ResourceHierarchyAttribute = resourceHierarchyAttribute
	return _options
}

// SetSupportedAnonymousAccesses : Allow user to set SupportedAnonymousAccesses
func (_options *CreateIamRegistrationOptions) SetSupportedAnonymousAccesses(supportedAnonymousAccesses []IamServiceRegistrationSupportedAnonymousAccess) *CreateIamRegistrationOptions {
	_options.SupportedAnonymousAccesses = supportedAnonymousAccesses
	return _options
}

// SetSupportedAttributes : Allow user to set SupportedAttributes
func (_options *CreateIamRegistrationOptions) SetSupportedAttributes(supportedAttributes []IamServiceRegistrationSupportedAttribute) *CreateIamRegistrationOptions {
	_options.SupportedAttributes = supportedAttributes
	return _options
}

// SetSupportedAuthorizationSubjects : Allow user to set SupportedAuthorizationSubjects
func (_options *CreateIamRegistrationOptions) SetSupportedAuthorizationSubjects(supportedAuthorizationSubjects []IamServiceRegistrationSupportedAuthorizationSubject) *CreateIamRegistrationOptions {
	_options.SupportedAuthorizationSubjects = supportedAuthorizationSubjects
	return _options
}

// SetSupportedRoles : Allow user to set SupportedRoles
func (_options *CreateIamRegistrationOptions) SetSupportedRoles(supportedRoles []IamServiceRegistrationSupportedRole) *CreateIamRegistrationOptions {
	_options.SupportedRoles = supportedRoles
	return _options
}

// SetSupportedNetwork : Allow user to set SupportedNetwork
func (_options *CreateIamRegistrationOptions) SetSupportedNetwork(supportedNetwork *IamServiceRegistrationSupportedNetwork) *CreateIamRegistrationOptions {
	_options.SupportedNetwork = supportedNetwork
	return _options
}

// SetEnv : Allow user to set Env
func (_options *CreateIamRegistrationOptions) SetEnv(env string) *CreateIamRegistrationOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateIamRegistrationOptions) SetHeaders(param map[string]string) *CreateIamRegistrationOptions {
	options.Headers = param
	return options
}

// CreateOnboardingProductOptions : The CreateOnboardingProduct options.
type CreateOnboardingProductOptions struct {
	// The type of the product.
	Type *string `json:"type" validate:"required"`

	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact" validate:"required"`

	// The Export Control Classification Number of your product.
	EccnNumber *string `json:"eccn_number,omitempty"`

	// The ERO class of your product.
	EroClass *string `json:"ero_class,omitempty"`

	// The United Nations Standard Products and Services Code of your product.
	Unspsc *float64 `json:"unspsc,omitempty"`

	// The tax assessment type of your product.
	TaxAssessment *string `json:"tax_assessment,omitempty"`

	// The support information that is not displayed in the catalog, but available in ServiceNow.
	Support *OnboardingProductSupport `json:"support,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateOnboardingProductOptions.Type property.
// The type of the product.
const (
	CreateOnboardingProductOptions_Type_ProfessionalService = "professional_service"
	CreateOnboardingProductOptions_Type_Service             = "service"
	CreateOnboardingProductOptions_Type_Software            = "software"
)

// NewCreateOnboardingProductOptions : Instantiate CreateOnboardingProductOptions
func (*PartnerCenterSellV1) NewCreateOnboardingProductOptions(typeVar string, primaryContact *PrimaryContact) *CreateOnboardingProductOptions {
	return &CreateOnboardingProductOptions{
		Type:           core.StringPtr(typeVar),
		PrimaryContact: primaryContact,
	}
}

// SetType : Allow user to set Type
func (_options *CreateOnboardingProductOptions) SetType(typeVar string) *CreateOnboardingProductOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPrimaryContact : Allow user to set PrimaryContact
func (_options *CreateOnboardingProductOptions) SetPrimaryContact(primaryContact *PrimaryContact) *CreateOnboardingProductOptions {
	_options.PrimaryContact = primaryContact
	return _options
}

// SetEccnNumber : Allow user to set EccnNumber
func (_options *CreateOnboardingProductOptions) SetEccnNumber(eccnNumber string) *CreateOnboardingProductOptions {
	_options.EccnNumber = core.StringPtr(eccnNumber)
	return _options
}

// SetEroClass : Allow user to set EroClass
func (_options *CreateOnboardingProductOptions) SetEroClass(eroClass string) *CreateOnboardingProductOptions {
	_options.EroClass = core.StringPtr(eroClass)
	return _options
}

// SetUnspsc : Allow user to set Unspsc
func (_options *CreateOnboardingProductOptions) SetUnspsc(unspsc float64) *CreateOnboardingProductOptions {
	_options.Unspsc = core.Float64Ptr(unspsc)
	return _options
}

// SetTaxAssessment : Allow user to set TaxAssessment
func (_options *CreateOnboardingProductOptions) SetTaxAssessment(taxAssessment string) *CreateOnboardingProductOptions {
	_options.TaxAssessment = core.StringPtr(taxAssessment)
	return _options
}

// SetSupport : Allow user to set Support
func (_options *CreateOnboardingProductOptions) SetSupport(support *OnboardingProductSupport) *CreateOnboardingProductOptions {
	_options.Support = support
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateOnboardingProductOptions) SetHeaders(param map[string]string) *CreateOnboardingProductOptions {
	options.Headers = param
	return options
}

// CreateRegistrationOptions : The CreateRegistration options.
type CreateRegistrationOptions struct {
	// The ID of your account.
	AccountID *string `json:"account_id" validate:"required"`

	// The name of your company that is displayed in the IBM Cloud catalog.
	CompanyName *string `json:"company_name" validate:"required"`

	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact" validate:"required"`

	// The default private catalog in which products are created.
	DefaultPrivateCatalogID *string `json:"default_private_catalog_id,omitempty"`

	// The onboarding access group for your team.
	ProviderAccessGroup *string `json:"provider_access_group,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateRegistrationOptions : Instantiate CreateRegistrationOptions
func (*PartnerCenterSellV1) NewCreateRegistrationOptions(accountID string, companyName string, primaryContact *PrimaryContact) *CreateRegistrationOptions {
	return &CreateRegistrationOptions{
		AccountID:      core.StringPtr(accountID),
		CompanyName:    core.StringPtr(companyName),
		PrimaryContact: primaryContact,
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateRegistrationOptions) SetAccountID(accountID string) *CreateRegistrationOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetCompanyName : Allow user to set CompanyName
func (_options *CreateRegistrationOptions) SetCompanyName(companyName string) *CreateRegistrationOptions {
	_options.CompanyName = core.StringPtr(companyName)
	return _options
}

// SetPrimaryContact : Allow user to set PrimaryContact
func (_options *CreateRegistrationOptions) SetPrimaryContact(primaryContact *PrimaryContact) *CreateRegistrationOptions {
	_options.PrimaryContact = primaryContact
	return _options
}

// SetDefaultPrivateCatalogID : Allow user to set DefaultPrivateCatalogID
func (_options *CreateRegistrationOptions) SetDefaultPrivateCatalogID(defaultPrivateCatalogID string) *CreateRegistrationOptions {
	_options.DefaultPrivateCatalogID = core.StringPtr(defaultPrivateCatalogID)
	return _options
}

// SetProviderAccessGroup : Allow user to set ProviderAccessGroup
func (_options *CreateRegistrationOptions) SetProviderAccessGroup(providerAccessGroup string) *CreateRegistrationOptions {
	_options.ProviderAccessGroup = core.StringPtr(providerAccessGroup)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRegistrationOptions) SetHeaders(param map[string]string) *CreateRegistrationOptions {
	options.Headers = param
	return options
}

// CreateResourceBrokerOptions : The CreateResourceBroker options.
type CreateResourceBrokerOptions struct {
	// The supported authentication scheme for the broker.
	AuthScheme *string `json:"auth_scheme" validate:"required"`

	// The name of the broker.
	Name *string `json:"name" validate:"required"`

	// The URL associated with the broker application.
	BrokerURL *string `json:"broker_url" validate:"required"`

	// The type of the provisioning model.
	Type *string `json:"type" validate:"required"`

	// The authentication username to reach the broker.
	AuthUsername *string `json:"auth_username,omitempty"`

	// The authentication password to reach the broker.
	AuthPassword *string `json:"auth_password,omitempty"`

	// The cloud resource name of the resource group.
	ResourceGroupCrn *string `json:"resource_group_crn,omitempty"`

	// The state of the broker.
	State *string `json:"state,omitempty"`

	// Whether the resource controller will call the broker for any context changes to the instance. Currently, the only
	// context related change is an instance name update.
	AllowContextUpdates *bool `json:"allow_context_updates,omitempty"`

	// To enable the provisioning of your broker, set this parameter value to `service`.
	CatalogType *string `json:"catalog_type,omitempty"`

	// The region where the pricing plan is available.
	Region *string `json:"region,omitempty"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateResourceBrokerOptions.AuthScheme property.
// The supported authentication scheme for the broker.
const (
	CreateResourceBrokerOptions_AuthScheme_Bearer    = "bearer"
	CreateResourceBrokerOptions_AuthScheme_BearerCrn = "bearer-crn"
)

// Constants associated with the CreateResourceBrokerOptions.Type property.
// The type of the provisioning model.
const (
	CreateResourceBrokerOptions_Type_ProvisionBehind  = "provision_behind"
	CreateResourceBrokerOptions_Type_ProvisionThrough = "provision_through"
)

// Constants associated with the CreateResourceBrokerOptions.AuthUsername property.
// The authentication username to reach the broker.
const (
	CreateResourceBrokerOptions_AuthUsername_Apikey = "apikey"
)

// Constants associated with the CreateResourceBrokerOptions.State property.
// The state of the broker.
const (
	CreateResourceBrokerOptions_State_Active  = "active"
	CreateResourceBrokerOptions_State_Removed = "removed"
)

// NewCreateResourceBrokerOptions : Instantiate CreateResourceBrokerOptions
func (*PartnerCenterSellV1) NewCreateResourceBrokerOptions(authScheme string, name string, brokerURL string, typeVar string) *CreateResourceBrokerOptions {
	return &CreateResourceBrokerOptions{
		AuthScheme: core.StringPtr(authScheme),
		Name:       core.StringPtr(name),
		BrokerURL:  core.StringPtr(brokerURL),
		Type:       core.StringPtr(typeVar),
	}
}

// SetAuthScheme : Allow user to set AuthScheme
func (_options *CreateResourceBrokerOptions) SetAuthScheme(authScheme string) *CreateResourceBrokerOptions {
	_options.AuthScheme = core.StringPtr(authScheme)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateResourceBrokerOptions) SetName(name string) *CreateResourceBrokerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetBrokerURL : Allow user to set BrokerURL
func (_options *CreateResourceBrokerOptions) SetBrokerURL(brokerURL string) *CreateResourceBrokerOptions {
	_options.BrokerURL = core.StringPtr(brokerURL)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateResourceBrokerOptions) SetType(typeVar string) *CreateResourceBrokerOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetAuthUsername : Allow user to set AuthUsername
func (_options *CreateResourceBrokerOptions) SetAuthUsername(authUsername string) *CreateResourceBrokerOptions {
	_options.AuthUsername = core.StringPtr(authUsername)
	return _options
}

// SetAuthPassword : Allow user to set AuthPassword
func (_options *CreateResourceBrokerOptions) SetAuthPassword(authPassword string) *CreateResourceBrokerOptions {
	_options.AuthPassword = core.StringPtr(authPassword)
	return _options
}

// SetResourceGroupCrn : Allow user to set ResourceGroupCrn
func (_options *CreateResourceBrokerOptions) SetResourceGroupCrn(resourceGroupCrn string) *CreateResourceBrokerOptions {
	_options.ResourceGroupCrn = core.StringPtr(resourceGroupCrn)
	return _options
}

// SetState : Allow user to set State
func (_options *CreateResourceBrokerOptions) SetState(state string) *CreateResourceBrokerOptions {
	_options.State = core.StringPtr(state)
	return _options
}

// SetAllowContextUpdates : Allow user to set AllowContextUpdates
func (_options *CreateResourceBrokerOptions) SetAllowContextUpdates(allowContextUpdates bool) *CreateResourceBrokerOptions {
	_options.AllowContextUpdates = core.BoolPtr(allowContextUpdates)
	return _options
}

// SetCatalogType : Allow user to set CatalogType
func (_options *CreateResourceBrokerOptions) SetCatalogType(catalogType string) *CreateResourceBrokerOptions {
	_options.CatalogType = core.StringPtr(catalogType)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *CreateResourceBrokerOptions) SetRegion(region string) *CreateResourceBrokerOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *CreateResourceBrokerOptions) SetEnv(env string) *CreateResourceBrokerOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateResourceBrokerOptions) SetHeaders(param map[string]string) *CreateResourceBrokerOptions {
	options.Headers = param
	return options
}

// DeleteCatalogDeploymentOptions : The DeleteCatalogDeployment options.
type DeleteCatalogDeploymentOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The unique ID of this global catalog deployment.
	CatalogDeploymentID *string `json:"catalog_deployment_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCatalogDeploymentOptions : Instantiate DeleteCatalogDeploymentOptions
func (*PartnerCenterSellV1) NewDeleteCatalogDeploymentOptions(productID string, catalogProductID string, catalogPlanID string, catalogDeploymentID string) *DeleteCatalogDeploymentOptions {
	return &DeleteCatalogDeploymentOptions{
		ProductID:           core.StringPtr(productID),
		CatalogProductID:    core.StringPtr(catalogProductID),
		CatalogPlanID:       core.StringPtr(catalogPlanID),
		CatalogDeploymentID: core.StringPtr(catalogDeploymentID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *DeleteCatalogDeploymentOptions) SetProductID(productID string) *DeleteCatalogDeploymentOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *DeleteCatalogDeploymentOptions) SetCatalogProductID(catalogProductID string) *DeleteCatalogDeploymentOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *DeleteCatalogDeploymentOptions) SetCatalogPlanID(catalogPlanID string) *DeleteCatalogDeploymentOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetCatalogDeploymentID : Allow user to set CatalogDeploymentID
func (_options *DeleteCatalogDeploymentOptions) SetCatalogDeploymentID(catalogDeploymentID string) *DeleteCatalogDeploymentOptions {
	_options.CatalogDeploymentID = core.StringPtr(catalogDeploymentID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *DeleteCatalogDeploymentOptions) SetEnv(env string) *DeleteCatalogDeploymentOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogDeploymentOptions) SetHeaders(param map[string]string) *DeleteCatalogDeploymentOptions {
	options.Headers = param
	return options
}

// DeleteCatalogPlanOptions : The DeleteCatalogPlan options.
type DeleteCatalogPlanOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCatalogPlanOptions : Instantiate DeleteCatalogPlanOptions
func (*PartnerCenterSellV1) NewDeleteCatalogPlanOptions(productID string, catalogProductID string, catalogPlanID string) *DeleteCatalogPlanOptions {
	return &DeleteCatalogPlanOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
		CatalogPlanID:    core.StringPtr(catalogPlanID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *DeleteCatalogPlanOptions) SetProductID(productID string) *DeleteCatalogPlanOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *DeleteCatalogPlanOptions) SetCatalogProductID(catalogProductID string) *DeleteCatalogPlanOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *DeleteCatalogPlanOptions) SetCatalogPlanID(catalogPlanID string) *DeleteCatalogPlanOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *DeleteCatalogPlanOptions) SetEnv(env string) *DeleteCatalogPlanOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogPlanOptions) SetHeaders(param map[string]string) *DeleteCatalogPlanOptions {
	options.Headers = param
	return options
}

// DeleteCatalogProductOptions : The DeleteCatalogProduct options.
type DeleteCatalogProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCatalogProductOptions : Instantiate DeleteCatalogProductOptions
func (*PartnerCenterSellV1) NewDeleteCatalogProductOptions(productID string, catalogProductID string) *DeleteCatalogProductOptions {
	return &DeleteCatalogProductOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *DeleteCatalogProductOptions) SetProductID(productID string) *DeleteCatalogProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *DeleteCatalogProductOptions) SetCatalogProductID(catalogProductID string) *DeleteCatalogProductOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *DeleteCatalogProductOptions) SetEnv(env string) *DeleteCatalogProductOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCatalogProductOptions) SetHeaders(param map[string]string) *DeleteCatalogProductOptions {
	options.Headers = param
	return options
}

// DeleteIamRegistrationOptions : The DeleteIamRegistration options.
type DeleteIamRegistrationOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The approved programmatic name of the product.
	ProgrammaticName *string `json:"programmatic_name" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteIamRegistrationOptions : Instantiate DeleteIamRegistrationOptions
func (*PartnerCenterSellV1) NewDeleteIamRegistrationOptions(productID string, programmaticName string) *DeleteIamRegistrationOptions {
	return &DeleteIamRegistrationOptions{
		ProductID:        core.StringPtr(productID),
		ProgrammaticName: core.StringPtr(programmaticName),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *DeleteIamRegistrationOptions) SetProductID(productID string) *DeleteIamRegistrationOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetProgrammaticName : Allow user to set ProgrammaticName
func (_options *DeleteIamRegistrationOptions) SetProgrammaticName(programmaticName string) *DeleteIamRegistrationOptions {
	_options.ProgrammaticName = core.StringPtr(programmaticName)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *DeleteIamRegistrationOptions) SetEnv(env string) *DeleteIamRegistrationOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteIamRegistrationOptions) SetHeaders(param map[string]string) *DeleteIamRegistrationOptions {
	options.Headers = param
	return options
}

// DeleteOnboardingProductOptions : The DeleteOnboardingProduct options.
type DeleteOnboardingProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteOnboardingProductOptions : Instantiate DeleteOnboardingProductOptions
func (*PartnerCenterSellV1) NewDeleteOnboardingProductOptions(productID string) *DeleteOnboardingProductOptions {
	return &DeleteOnboardingProductOptions{
		ProductID: core.StringPtr(productID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *DeleteOnboardingProductOptions) SetProductID(productID string) *DeleteOnboardingProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteOnboardingProductOptions) SetHeaders(param map[string]string) *DeleteOnboardingProductOptions {
	options.Headers = param
	return options
}

// DeleteRegistrationOptions : The DeleteRegistration options.
type DeleteRegistrationOptions struct {
	// The unique ID of your registration.
	RegistrationID *string `json:"registration_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteRegistrationOptions : Instantiate DeleteRegistrationOptions
func (*PartnerCenterSellV1) NewDeleteRegistrationOptions(registrationID string) *DeleteRegistrationOptions {
	return &DeleteRegistrationOptions{
		RegistrationID: core.StringPtr(registrationID),
	}
}

// SetRegistrationID : Allow user to set RegistrationID
func (_options *DeleteRegistrationOptions) SetRegistrationID(registrationID string) *DeleteRegistrationOptions {
	_options.RegistrationID = core.StringPtr(registrationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRegistrationOptions) SetHeaders(param map[string]string) *DeleteRegistrationOptions {
	options.Headers = param
	return options
}

// DeleteResourceBrokerOptions : The DeleteResourceBroker options.
type DeleteResourceBrokerOptions struct {
	// The unique identifier of the broker.
	BrokerID *string `json:"broker_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Remove a broker with the ID that was provided in the API call URL.
	RemoveFromAccount *bool `json:"remove_from_account,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteResourceBrokerOptions : Instantiate DeleteResourceBrokerOptions
func (*PartnerCenterSellV1) NewDeleteResourceBrokerOptions(brokerID string) *DeleteResourceBrokerOptions {
	return &DeleteResourceBrokerOptions{
		BrokerID: core.StringPtr(brokerID),
	}
}

// SetBrokerID : Allow user to set BrokerID
func (_options *DeleteResourceBrokerOptions) SetBrokerID(brokerID string) *DeleteResourceBrokerOptions {
	_options.BrokerID = core.StringPtr(brokerID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *DeleteResourceBrokerOptions) SetEnv(env string) *DeleteResourceBrokerOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetRemoveFromAccount : Allow user to set RemoveFromAccount
func (_options *DeleteResourceBrokerOptions) SetRemoveFromAccount(removeFromAccount bool) *DeleteResourceBrokerOptions {
	_options.RemoveFromAccount = core.BoolPtr(removeFromAccount)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteResourceBrokerOptions) SetHeaders(param map[string]string) *DeleteResourceBrokerOptions {
	options.Headers = param
	return options
}

// EnvironmentAttribute : EnvironmentAttribute struct
type EnvironmentAttribute struct {
	// The name of the key.
	Key *string `json:"key,omitempty"`

	// The list of values that belong to the key.
	Values []string `json:"values,omitempty"`

	// The list of options for supported networks.
	Options *EnvironmentAttributeOptions `json:"options,omitempty"`
}

// UnmarshalEnvironmentAttribute unmarshals an instance of EnvironmentAttribute from the specified map of raw messages.
func UnmarshalEnvironmentAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		err = core.SDKErrorf(err, "", "values-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "options", &obj.Options, UnmarshalEnvironmentAttributeOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "options-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the EnvironmentAttribute
func (environmentAttribute *EnvironmentAttribute) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(environmentAttribute.Key) {
		_patch["key"] = environmentAttribute.Key
	}
	if !core.IsNil(environmentAttribute.Values) {
		_patch["values"] = environmentAttribute.Values
	}
	if !core.IsNil(environmentAttribute.Options) {
		_patch["options"] = environmentAttribute.Options.asPatch()
	}

	return
}

// EnvironmentAttributeOptions : The list of options for supported networks.
type EnvironmentAttributeOptions struct {
	// Whether the attribute is hidden or not.
	Hidden *bool `json:"hidden,omitempty"`
}

// UnmarshalEnvironmentAttributeOptions unmarshals an instance of EnvironmentAttributeOptions from the specified map of raw messages.
func UnmarshalEnvironmentAttributeOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentAttributeOptions)
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the EnvironmentAttributeOptions
func (environmentAttributeOptions *EnvironmentAttributeOptions) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(environmentAttributeOptions.Hidden) {
		_patch["hidden"] = environmentAttributeOptions.Hidden
	}

	return
}

// GetCatalogDeploymentOptions : The GetCatalogDeployment options.
type GetCatalogDeploymentOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The unique ID of this global catalog deployment.
	CatalogDeploymentID *string `json:"catalog_deployment_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCatalogDeploymentOptions : Instantiate GetCatalogDeploymentOptions
func (*PartnerCenterSellV1) NewGetCatalogDeploymentOptions(productID string, catalogProductID string, catalogPlanID string, catalogDeploymentID string) *GetCatalogDeploymentOptions {
	return &GetCatalogDeploymentOptions{
		ProductID:           core.StringPtr(productID),
		CatalogProductID:    core.StringPtr(catalogProductID),
		CatalogPlanID:       core.StringPtr(catalogPlanID),
		CatalogDeploymentID: core.StringPtr(catalogDeploymentID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *GetCatalogDeploymentOptions) SetProductID(productID string) *GetCatalogDeploymentOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *GetCatalogDeploymentOptions) SetCatalogProductID(catalogProductID string) *GetCatalogDeploymentOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *GetCatalogDeploymentOptions) SetCatalogPlanID(catalogPlanID string) *GetCatalogDeploymentOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetCatalogDeploymentID : Allow user to set CatalogDeploymentID
func (_options *GetCatalogDeploymentOptions) SetCatalogDeploymentID(catalogDeploymentID string) *GetCatalogDeploymentOptions {
	_options.CatalogDeploymentID = core.StringPtr(catalogDeploymentID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *GetCatalogDeploymentOptions) SetEnv(env string) *GetCatalogDeploymentOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogDeploymentOptions) SetHeaders(param map[string]string) *GetCatalogDeploymentOptions {
	options.Headers = param
	return options
}

// GetCatalogPlanOptions : The GetCatalogPlan options.
type GetCatalogPlanOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCatalogPlanOptions : Instantiate GetCatalogPlanOptions
func (*PartnerCenterSellV1) NewGetCatalogPlanOptions(productID string, catalogProductID string, catalogPlanID string) *GetCatalogPlanOptions {
	return &GetCatalogPlanOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
		CatalogPlanID:    core.StringPtr(catalogPlanID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *GetCatalogPlanOptions) SetProductID(productID string) *GetCatalogPlanOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *GetCatalogPlanOptions) SetCatalogProductID(catalogProductID string) *GetCatalogPlanOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *GetCatalogPlanOptions) SetCatalogPlanID(catalogPlanID string) *GetCatalogPlanOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *GetCatalogPlanOptions) SetEnv(env string) *GetCatalogPlanOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogPlanOptions) SetHeaders(param map[string]string) *GetCatalogPlanOptions {
	options.Headers = param
	return options
}

// GetCatalogProductOptions : The GetCatalogProduct options.
type GetCatalogProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCatalogProductOptions : Instantiate GetCatalogProductOptions
func (*PartnerCenterSellV1) NewGetCatalogProductOptions(productID string, catalogProductID string) *GetCatalogProductOptions {
	return &GetCatalogProductOptions{
		ProductID:        core.StringPtr(productID),
		CatalogProductID: core.StringPtr(catalogProductID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *GetCatalogProductOptions) SetProductID(productID string) *GetCatalogProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *GetCatalogProductOptions) SetCatalogProductID(catalogProductID string) *GetCatalogProductOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *GetCatalogProductOptions) SetEnv(env string) *GetCatalogProductOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCatalogProductOptions) SetHeaders(param map[string]string) *GetCatalogProductOptions {
	options.Headers = param
	return options
}

// GetIamRegistrationOptions : The GetIamRegistration options.
type GetIamRegistrationOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The approved programmatic name of the product.
	ProgrammaticName *string `json:"programmatic_name" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetIamRegistrationOptions : Instantiate GetIamRegistrationOptions
func (*PartnerCenterSellV1) NewGetIamRegistrationOptions(productID string, programmaticName string) *GetIamRegistrationOptions {
	return &GetIamRegistrationOptions{
		ProductID:        core.StringPtr(productID),
		ProgrammaticName: core.StringPtr(programmaticName),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *GetIamRegistrationOptions) SetProductID(productID string) *GetIamRegistrationOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetProgrammaticName : Allow user to set ProgrammaticName
func (_options *GetIamRegistrationOptions) SetProgrammaticName(programmaticName string) *GetIamRegistrationOptions {
	_options.ProgrammaticName = core.StringPtr(programmaticName)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *GetIamRegistrationOptions) SetEnv(env string) *GetIamRegistrationOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetIamRegistrationOptions) SetHeaders(param map[string]string) *GetIamRegistrationOptions {
	options.Headers = param
	return options
}

// GetOnboardingProductOptions : The GetOnboardingProduct options.
type GetOnboardingProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetOnboardingProductOptions : Instantiate GetOnboardingProductOptions
func (*PartnerCenterSellV1) NewGetOnboardingProductOptions(productID string) *GetOnboardingProductOptions {
	return &GetOnboardingProductOptions{
		ProductID: core.StringPtr(productID),
	}
}

// SetProductID : Allow user to set ProductID
func (_options *GetOnboardingProductOptions) SetProductID(productID string) *GetOnboardingProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOnboardingProductOptions) SetHeaders(param map[string]string) *GetOnboardingProductOptions {
	options.Headers = param
	return options
}

// GetProductBadgeOptions : The GetProductBadge options.
type GetProductBadgeOptions struct {
	// The unique ID of the badge. This ID can be obtained by calling the list badges method.
	BadgeID *strfmt.UUID `json:"badge_id" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProductBadgeOptions : Instantiate GetProductBadgeOptions
func (*PartnerCenterSellV1) NewGetProductBadgeOptions(badgeID *strfmt.UUID) *GetProductBadgeOptions {
	return &GetProductBadgeOptions{
		BadgeID: badgeID,
	}
}

// SetBadgeID : Allow user to set BadgeID
func (_options *GetProductBadgeOptions) SetBadgeID(badgeID *strfmt.UUID) *GetProductBadgeOptions {
	_options.BadgeID = badgeID
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProductBadgeOptions) SetHeaders(param map[string]string) *GetProductBadgeOptions {
	options.Headers = param
	return options
}

// GetRegistrationOptions : The GetRegistration options.
type GetRegistrationOptions struct {
	// The unique ID of your registration.
	RegistrationID *string `json:"registration_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetRegistrationOptions : Instantiate GetRegistrationOptions
func (*PartnerCenterSellV1) NewGetRegistrationOptions(registrationID string) *GetRegistrationOptions {
	return &GetRegistrationOptions{
		RegistrationID: core.StringPtr(registrationID),
	}
}

// SetRegistrationID : Allow user to set RegistrationID
func (_options *GetRegistrationOptions) SetRegistrationID(registrationID string) *GetRegistrationOptions {
	_options.RegistrationID = core.StringPtr(registrationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRegistrationOptions) SetHeaders(param map[string]string) *GetRegistrationOptions {
	options.Headers = param
	return options
}

// GetResourceBrokerOptions : The GetResourceBroker options.
type GetResourceBrokerOptions struct {
	// The unique identifier of the broker.
	BrokerID *string `json:"broker_id" validate:"required,ne="`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetResourceBrokerOptions : Instantiate GetResourceBrokerOptions
func (*PartnerCenterSellV1) NewGetResourceBrokerOptions(brokerID string) *GetResourceBrokerOptions {
	return &GetResourceBrokerOptions{
		BrokerID: core.StringPtr(brokerID),
	}
}

// SetBrokerID : Allow user to set BrokerID
func (_options *GetResourceBrokerOptions) SetBrokerID(brokerID string) *GetResourceBrokerOptions {
	_options.BrokerID = core.StringPtr(brokerID)
	return _options
}

// SetEnv : Allow user to set Env
func (_options *GetResourceBrokerOptions) SetEnv(env string) *GetResourceBrokerOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetResourceBrokerOptions) SetHeaders(param map[string]string) *GetResourceBrokerOptions {
	options.Headers = param
	return options
}

// GlobalCatalogDeployment : The object defining a global catalog deployment.
type GlobalCatalogDeployment struct {
	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The programmatic name of this deployment.
	Name *string `json:"name,omitempty"`

	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The kind of the global catalog object.
	Kind *string `json:"kind,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// The global catalog URL of your product.
	URL *string `json:"url,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// Global catalog deployment metadata.
	Metadata *GlobalCatalogDeploymentMetadata `json:"metadata,omitempty"`
}

// Constants associated with the GlobalCatalogDeployment.Kind property.
// The kind of the global catalog object.
const (
	GlobalCatalogDeployment_Kind_Deployment = "deployment"
)

// UnmarshalGlobalCatalogDeployment unmarshals an instance of GlobalCatalogDeployment from the specified map of raw messages.
func UnmarshalGlobalCatalogDeployment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogDeployment)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "object_id", &obj.ObjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogDeploymentMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GlobalCatalogDeploymentMetadata : Global catalog deployment metadata.
type GlobalCatalogDeploymentMetadata struct {
	// Whether the object is compatible with the resource controller service.
	RcCompatible *bool `json:"rc_compatible,omitempty"`

	// The UI metadata of this service.
	Ui *GlobalCatalogMetadataUI `json:"ui,omitempty"`

	// The global catalog metadata of the service.
	Service *GlobalCatalogMetadataService `json:"service,omitempty"`

	// The global catalog metadata of the deployment.
	Deployment *GlobalCatalogMetadataDeployment `json:"deployment,omitempty"`
}

// UnmarshalGlobalCatalogDeploymentMetadata unmarshals an instance of GlobalCatalogDeploymentMetadata from the specified map of raw messages.
func UnmarshalGlobalCatalogDeploymentMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogDeploymentMetadata)
	err = core.UnmarshalPrimitive(m, "rc_compatible", &obj.RcCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.Ui, UnmarshalGlobalCatalogMetadataUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalGlobalCatalogMetadataService)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "deployment", &obj.Deployment, UnmarshalGlobalCatalogMetadataDeployment)
	if err != nil {
		err = core.SDKErrorf(err, "", "deployment-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogDeploymentMetadata
func (globalCatalogDeploymentMetadata *GlobalCatalogDeploymentMetadata) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogDeploymentMetadata.RcCompatible) {
		_patch["rc_compatible"] = globalCatalogDeploymentMetadata.RcCompatible
	}
	if !core.IsNil(globalCatalogDeploymentMetadata.Ui) {
		_patch["ui"] = globalCatalogDeploymentMetadata.Ui.asPatch()
	}
	if !core.IsNil(globalCatalogDeploymentMetadata.Service) {
		_patch["service"] = globalCatalogDeploymentMetadata.Service.asPatch()
	}
	if !core.IsNil(globalCatalogDeploymentMetadata.Deployment) {
		_patch["deployment"] = globalCatalogDeploymentMetadata.Deployment.asPatch()
	}

	return
}

// GlobalCatalogDeploymentPatch : The request body for updating a global catalog deployment.
type GlobalCatalogDeploymentPatch struct {
	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// Global catalog deployment metadata.
	Metadata *GlobalCatalogDeploymentMetadata `json:"metadata,omitempty"`
}

// UnmarshalGlobalCatalogDeploymentPatch unmarshals an instance of GlobalCatalogDeploymentPatch from the specified map of raw messages.
func UnmarshalGlobalCatalogDeploymentPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogDeploymentPatch)
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogDeploymentMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the GlobalCatalogDeploymentPatch
func (globalCatalogDeploymentPatch *GlobalCatalogDeploymentPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogDeploymentPatch.Active) {
		_patch["active"] = globalCatalogDeploymentPatch.Active
	}
	if !core.IsNil(globalCatalogDeploymentPatch.Disabled) {
		_patch["disabled"] = globalCatalogDeploymentPatch.Disabled
	}
	if !core.IsNil(globalCatalogDeploymentPatch.OverviewUi) {
		_patch["overview_ui"] = globalCatalogDeploymentPatch.OverviewUi.asPatch()
	}
	if !core.IsNil(globalCatalogDeploymentPatch.Tags) {
		_patch["tags"] = globalCatalogDeploymentPatch.Tags
	}
	if !core.IsNil(globalCatalogDeploymentPatch.ObjectProvider) {
		_patch["object_provider"] = globalCatalogDeploymentPatch.ObjectProvider.asPatch()
	}
	if !core.IsNil(globalCatalogDeploymentPatch.Metadata) {
		_patch["metadata"] = globalCatalogDeploymentPatch.Metadata.asPatch()
	}

	return
}

// GlobalCatalogMetadataDeployment : The global catalog metadata of the deployment.
type GlobalCatalogMetadataDeployment struct {
	// The global catalog metadata of the deployment.
	Broker *GlobalCatalogMetadataDeploymentBroker `json:"broker,omitempty"`

	// The global catalog deployment location.
	Location *string `json:"location,omitempty"`

	// The global catalog deployment URL of location.
	LocationURL *string `json:"location_url,omitempty"`

	// Region crn.
	TargetCrn *string `json:"target_crn,omitempty"`
}

// UnmarshalGlobalCatalogMetadataDeployment unmarshals an instance of GlobalCatalogMetadataDeployment from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataDeployment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataDeployment)
	err = core.UnmarshalModel(m, "broker", &obj.Broker, UnmarshalGlobalCatalogMetadataDeploymentBroker)
	if err != nil {
		err = core.SDKErrorf(err, "", "broker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location_url", &obj.LocationURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "location_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataDeployment
func (globalCatalogMetadataDeployment *GlobalCatalogMetadataDeployment) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataDeployment.Broker) {
		_patch["broker"] = globalCatalogMetadataDeployment.Broker.asPatch()
	}
	if !core.IsNil(globalCatalogMetadataDeployment.Location) {
		_patch["location"] = globalCatalogMetadataDeployment.Location
	}
	if !core.IsNil(globalCatalogMetadataDeployment.LocationURL) {
		_patch["location_url"] = globalCatalogMetadataDeployment.LocationURL
	}
	if !core.IsNil(globalCatalogMetadataDeployment.TargetCrn) {
		_patch["target_crn"] = globalCatalogMetadataDeployment.TargetCrn
	}

	return
}

// GlobalCatalogMetadataDeploymentBroker : The global catalog metadata of the deployment.
type GlobalCatalogMetadataDeploymentBroker struct {
	// The name of the resource broker.
	Name *string `json:"name,omitempty"`

	// Crn or guid of the resource broker.
	Guid *string `json:"guid,omitempty"`
}

// UnmarshalGlobalCatalogMetadataDeploymentBroker unmarshals an instance of GlobalCatalogMetadataDeploymentBroker from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataDeploymentBroker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataDeploymentBroker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "guid", &obj.Guid)
	if err != nil {
		err = core.SDKErrorf(err, "", "guid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataDeploymentBroker
func (globalCatalogMetadataDeploymentBroker *GlobalCatalogMetadataDeploymentBroker) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataDeploymentBroker.Name) {
		_patch["name"] = globalCatalogMetadataDeploymentBroker.Name
	}
	if !core.IsNil(globalCatalogMetadataDeploymentBroker.Guid) {
		_patch["guid"] = globalCatalogMetadataDeploymentBroker.Guid
	}

	return
}

// GlobalCatalogMetadataPricing : The pricing metadata of this object.
type GlobalCatalogMetadataPricing struct {
	// The type of the pricing plan.
	Type *string `json:"type,omitempty"`

	// The source of the pricing information: global_catalog or pricing_catalog.
	Origin *string `json:"origin,omitempty"`
}

// Constants associated with the GlobalCatalogMetadataPricing.Type property.
// The type of the pricing plan.
const (
	GlobalCatalogMetadataPricing_Type_Free         = "free"
	GlobalCatalogMetadataPricing_Type_Paid         = "paid"
	GlobalCatalogMetadataPricing_Type_Subscription = "subscription"
)

// Constants associated with the GlobalCatalogMetadataPricing.Origin property.
// The source of the pricing information: global_catalog or pricing_catalog.
const (
	GlobalCatalogMetadataPricing_Origin_GlobalCatalog  = "global_catalog"
	GlobalCatalogMetadataPricing_Origin_PricingCatalog = "pricing_catalog"
)

// UnmarshalGlobalCatalogMetadataPricing unmarshals an instance of GlobalCatalogMetadataPricing from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataPricing(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataPricing)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "origin", &obj.Origin)
	if err != nil {
		err = core.SDKErrorf(err, "", "origin-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataPricing
func (globalCatalogMetadataPricing *GlobalCatalogMetadataPricing) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataPricing.Type) {
		_patch["type"] = globalCatalogMetadataPricing.Type
	}
	if !core.IsNil(globalCatalogMetadataPricing.Origin) {
		_patch["origin"] = globalCatalogMetadataPricing.Origin
	}

	return
}

// GlobalCatalogMetadataService : The global catalog metadata of the service.
type GlobalCatalogMetadataService struct {
	// Whether the service is provisionable by the resource controller service.
	RcProvisionable *bool `json:"rc_provisionable,omitempty"`

	// Whether the service is compatible with the IAM service.
	IamCompatible *bool `json:"iam_compatible,omitempty"`

	// Deprecated. Controls the Connections tab on the Resource Details page.
	Bindable *bool `json:"bindable,omitempty"`

	// Indicates plan update support and controls the Plan tab on the Resource Details page.
	PlanUpdateable *bool `json:"plan_updateable,omitempty"`

	// Indicates service credentials support and controls the Service Credential tab on Resource Details page.
	ServiceKeySupported *bool `json:"service_key_supported,omitempty"`
}

// UnmarshalGlobalCatalogMetadataService unmarshals an instance of GlobalCatalogMetadataService from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataService(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataService)
	err = core.UnmarshalPrimitive(m, "rc_provisionable", &obj.RcProvisionable)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_provisionable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_compatible", &obj.IamCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bindable", &obj.Bindable)
	if err != nil {
		err = core.SDKErrorf(err, "", "bindable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "plan_updateable", &obj.PlanUpdateable)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan_updateable-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_key_supported", &obj.ServiceKeySupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_key_supported-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataService
func (globalCatalogMetadataService *GlobalCatalogMetadataService) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataService.RcProvisionable) {
		_patch["rc_provisionable"] = globalCatalogMetadataService.RcProvisionable
	}
	if !core.IsNil(globalCatalogMetadataService.IamCompatible) {
		_patch["iam_compatible"] = globalCatalogMetadataService.IamCompatible
	}
	if !core.IsNil(globalCatalogMetadataService.Bindable) {
		_patch["bindable"] = globalCatalogMetadataService.Bindable
	}
	if !core.IsNil(globalCatalogMetadataService.PlanUpdateable) {
		_patch["plan_updateable"] = globalCatalogMetadataService.PlanUpdateable
	}
	if !core.IsNil(globalCatalogMetadataService.ServiceKeySupported) {
		_patch["service_key_supported"] = globalCatalogMetadataService.ServiceKeySupported
	}

	return
}

// GlobalCatalogMetadataUI : The UI metadata of this service.
type GlobalCatalogMetadataUI struct {
	// The data strings.
	Strings *GlobalCatalogMetadataUIStrings `json:"strings,omitempty"`

	// Metadata with URLs related to a service.
	Urls *GlobalCatalogMetadataUIUrls `json:"urls,omitempty"`

	// Whether the object is hidden from the consumption catalog.
	Hidden *bool `json:"hidden,omitempty"`

	// When the objects are listed side-by-side, this value controls the ordering.
	SideBySideIndex *float64 `json:"side_by_side_index,omitempty"`
}

// UnmarshalGlobalCatalogMetadataUI unmarshals an instance of GlobalCatalogMetadataUI from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataUI(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataUI)
	err = core.UnmarshalModel(m, "strings", &obj.Strings, UnmarshalGlobalCatalogMetadataUIStrings)
	if err != nil {
		err = core.SDKErrorf(err, "", "strings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "urls", &obj.Urls, UnmarshalGlobalCatalogMetadataUIUrls)
	if err != nil {
		err = core.SDKErrorf(err, "", "urls-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "side_by_side_index", &obj.SideBySideIndex)
	if err != nil {
		err = core.SDKErrorf(err, "", "side_by_side_index-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataUI
func (globalCatalogMetadataUI *GlobalCatalogMetadataUI) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataUI.Strings) {
		_patch["strings"] = globalCatalogMetadataUI.Strings.asPatch()
	}
	if !core.IsNil(globalCatalogMetadataUI.Urls) {
		_patch["urls"] = globalCatalogMetadataUI.Urls.asPatch()
	}
	if !core.IsNil(globalCatalogMetadataUI.Hidden) {
		_patch["hidden"] = globalCatalogMetadataUI.Hidden
	}
	if !core.IsNil(globalCatalogMetadataUI.SideBySideIndex) {
		_patch["side_by_side_index"] = globalCatalogMetadataUI.SideBySideIndex
	}

	return
}

// GlobalCatalogMetadataUIStrings : The data strings.
type GlobalCatalogMetadataUIStrings struct {
	// Translated content of additional information about the service.
	En *GlobalCatalogMetadataUIStringsContent `json:"en,omitempty"`
}

// UnmarshalGlobalCatalogMetadataUIStrings unmarshals an instance of GlobalCatalogMetadataUIStrings from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataUIStrings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataUIStrings)
	err = core.UnmarshalModel(m, "en", &obj.En, UnmarshalGlobalCatalogMetadataUIStringsContent)
	if err != nil {
		err = core.SDKErrorf(err, "", "en-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataUIStrings
func (globalCatalogMetadataUIStrings *GlobalCatalogMetadataUIStrings) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataUIStrings.En) {
		_patch["en"] = globalCatalogMetadataUIStrings.En.asPatch()
	}

	return
}

// GlobalCatalogMetadataUIStringsContent : Translated content of additional information about the service.
type GlobalCatalogMetadataUIStringsContent struct {
	// The list of features that highlights your product's attributes and benefits for users.
	Bullets []CatalogHighlightItem `json:"bullets,omitempty"`

	// The list of supporting media for this product.
	Media []CatalogProductMediaItem `json:"media,omitempty"`

	// On a service kind record this controls if your service has a custom dashboard or Resource Detail page.
	EmbeddableDashboard *string `json:"embeddable_dashboard,omitempty"`
}

// UnmarshalGlobalCatalogMetadataUIStringsContent unmarshals an instance of GlobalCatalogMetadataUIStringsContent from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataUIStringsContent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataUIStringsContent)
	err = core.UnmarshalModel(m, "bullets", &obj.Bullets, UnmarshalCatalogHighlightItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "bullets-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "media", &obj.Media, UnmarshalCatalogProductMediaItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "media-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "embeddable_dashboard", &obj.EmbeddableDashboard)
	if err != nil {
		err = core.SDKErrorf(err, "", "embeddable_dashboard-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataUIStringsContent
func (globalCatalogMetadataUIStringsContent *GlobalCatalogMetadataUIStringsContent) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataUIStringsContent.Bullets) {
		var bulletsPatches []map[string]interface{}
		for _, bullets := range globalCatalogMetadataUIStringsContent.Bullets {
			bulletsPatches = append(bulletsPatches, bullets.asPatch())
		}
		_patch["bullets"] = bulletsPatches
	}
	if !core.IsNil(globalCatalogMetadataUIStringsContent.Media) {
		var mediaPatches []map[string]interface{}
		for _, media := range globalCatalogMetadataUIStringsContent.Media {
			mediaPatches = append(mediaPatches, media.asPatch())
		}
		_patch["media"] = mediaPatches
	}
	if !core.IsNil(globalCatalogMetadataUIStringsContent.EmbeddableDashboard) {
		_patch["embeddable_dashboard"] = globalCatalogMetadataUIStringsContent.EmbeddableDashboard
	}

	return
}

// GlobalCatalogMetadataUIUrls : Metadata with URLs related to a service.
type GlobalCatalogMetadataUIUrls struct {
	// The URL for your product's documentation.
	DocURL *string `json:"doc_url,omitempty"`

	// The URL for your product's API documentation.
	ApidocsURL *string `json:"apidocs_url,omitempty"`

	// The URL for your product's end user license agreement.
	TermsURL *string `json:"terms_url,omitempty"`

	// Controls the Getting Started tab on the Resource Details page. Setting it the content is loaded from the specified
	// URL.
	InstructionsURL *string `json:"instructions_url,omitempty"`

	// Controls the Provisioning page URL, if set the assumption is that this URL is the provisioning URL for your service.
	CatalogDetailsURL *string `json:"catalog_details_url,omitempty"`

	// Controls the Provisioning page URL, if set the assumption is that this URL is the provisioning URL for your service.
	CustomCreatePageURL *string `json:"custom_create_page_url,omitempty"`

	// Controls if your service has a custom dashboard or Resource Detail page.
	Dashboard *string `json:"dashboard,omitempty"`
}

// UnmarshalGlobalCatalogMetadataUIUrls unmarshals an instance of GlobalCatalogMetadataUIUrls from the specified map of raw messages.
func UnmarshalGlobalCatalogMetadataUIUrls(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogMetadataUIUrls)
	err = core.UnmarshalPrimitive(m, "doc_url", &obj.DocURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "doc_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "apidocs_url", &obj.ApidocsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "apidocs_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "terms_url", &obj.TermsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "terms_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "instructions_url", &obj.InstructionsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "instructions_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "catalog_details_url", &obj.CatalogDetailsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "catalog_details_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "custom_create_page_url", &obj.CustomCreatePageURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "custom_create_page_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dashboard", &obj.Dashboard)
	if err != nil {
		err = core.SDKErrorf(err, "", "dashboard-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogMetadataUIUrls
func (globalCatalogMetadataUIUrls *GlobalCatalogMetadataUIUrls) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogMetadataUIUrls.DocURL) {
		_patch["doc_url"] = globalCatalogMetadataUIUrls.DocURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.ApidocsURL) {
		_patch["apidocs_url"] = globalCatalogMetadataUIUrls.ApidocsURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.TermsURL) {
		_patch["terms_url"] = globalCatalogMetadataUIUrls.TermsURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.InstructionsURL) {
		_patch["instructions_url"] = globalCatalogMetadataUIUrls.InstructionsURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.CatalogDetailsURL) {
		_patch["catalog_details_url"] = globalCatalogMetadataUIUrls.CatalogDetailsURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.CustomCreatePageURL) {
		_patch["custom_create_page_url"] = globalCatalogMetadataUIUrls.CustomCreatePageURL
	}
	if !core.IsNil(globalCatalogMetadataUIUrls.Dashboard) {
		_patch["dashboard"] = globalCatalogMetadataUIUrls.Dashboard
	}

	return
}

// GlobalCatalogOverviewUI : The object that contains the service details from the Overview page in global catalog.
type GlobalCatalogOverviewUI struct {
	// Translated details about the service, for example, display name, short description, and long description.
	En *GlobalCatalogOverviewUITranslatedContent `json:"en,omitempty"`
}

// UnmarshalGlobalCatalogOverviewUI unmarshals an instance of GlobalCatalogOverviewUI from the specified map of raw messages.
func UnmarshalGlobalCatalogOverviewUI(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogOverviewUI)
	err = core.UnmarshalModel(m, "en", &obj.En, UnmarshalGlobalCatalogOverviewUITranslatedContent)
	if err != nil {
		err = core.SDKErrorf(err, "", "en-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogOverviewUI
func (globalCatalogOverviewUI *GlobalCatalogOverviewUI) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogOverviewUI.En) {
		_patch["en"] = globalCatalogOverviewUI.En.asPatch()
	}

	return
}

// GlobalCatalogOverviewUITranslatedContent : Translated details about the service, for example, display name, short description, and long description.
type GlobalCatalogOverviewUITranslatedContent struct {
	// The display name of the product.
	DisplayName *string `json:"display_name,omitempty"`

	// The short description of the product that is displayed in your catalog entry.
	Description *string `json:"description,omitempty"`

	// The detailed description of your product that is displayed at the beginning of your product page in the catalog.
	// Markdown markup language is supported.
	LongDescription *string `json:"long_description,omitempty"`
}

// UnmarshalGlobalCatalogOverviewUITranslatedContent unmarshals an instance of GlobalCatalogOverviewUITranslatedContent from the specified map of raw messages.
func UnmarshalGlobalCatalogOverviewUITranslatedContent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogOverviewUITranslatedContent)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "long_description", &obj.LongDescription)
	if err != nil {
		err = core.SDKErrorf(err, "", "long_description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogOverviewUITranslatedContent
func (globalCatalogOverviewUITranslatedContent *GlobalCatalogOverviewUITranslatedContent) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogOverviewUITranslatedContent.DisplayName) {
		_patch["display_name"] = globalCatalogOverviewUITranslatedContent.DisplayName
	}
	if !core.IsNil(globalCatalogOverviewUITranslatedContent.Description) {
		_patch["description"] = globalCatalogOverviewUITranslatedContent.Description
	}
	if !core.IsNil(globalCatalogOverviewUITranslatedContent.LongDescription) {
		_patch["long_description"] = globalCatalogOverviewUITranslatedContent.LongDescription
	}

	return
}

// GlobalCatalogPlan : The object defining a global catalog plan.
type GlobalCatalogPlan struct {
	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The programmatic name of this plan.
	Name *string `json:"name,omitempty"`

	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The kind of the global catalog object.
	Kind *string `json:"kind,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// The global catalog URL of your product.
	URL *string `json:"url,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// Global catalog plan metadata.
	Metadata *GlobalCatalogPlanMetadata `json:"metadata,omitempty"`
}

// Constants associated with the GlobalCatalogPlan.Kind property.
// The kind of the global catalog object.
const (
	GlobalCatalogPlan_Kind_Plan = "plan"
)

// UnmarshalGlobalCatalogPlan unmarshals an instance of GlobalCatalogPlan from the specified map of raw messages.
func UnmarshalGlobalCatalogPlan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogPlan)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "object_id", &obj.ObjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogPlanMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GlobalCatalogPlanMetadata : Global catalog plan metadata.
type GlobalCatalogPlanMetadata struct {
	// Whether the object is compatible with the resource controller service.
	RcCompatible *bool `json:"rc_compatible,omitempty"`

	// The UI metadata of this service.
	Ui *GlobalCatalogMetadataUI `json:"ui,omitempty"`

	// The global catalog metadata of the service.
	Service *GlobalCatalogMetadataService `json:"service,omitempty"`

	// The pricing metadata of this object.
	Pricing *GlobalCatalogMetadataPricing `json:"pricing,omitempty"`

	// Metadata controlling Plan related settings.
	Plan *GlobalCatalogPlanMetadataPlan `json:"plan,omitempty"`
}

// UnmarshalGlobalCatalogPlanMetadata unmarshals an instance of GlobalCatalogPlanMetadata from the specified map of raw messages.
func UnmarshalGlobalCatalogPlanMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogPlanMetadata)
	err = core.UnmarshalPrimitive(m, "rc_compatible", &obj.RcCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.Ui, UnmarshalGlobalCatalogMetadataUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalGlobalCatalogMetadataService)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pricing", &obj.Pricing, UnmarshalGlobalCatalogMetadataPricing)
	if err != nil {
		err = core.SDKErrorf(err, "", "pricing-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "plan", &obj.Plan, UnmarshalGlobalCatalogPlanMetadataPlan)
	if err != nil {
		err = core.SDKErrorf(err, "", "plan-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogPlanMetadata
func (globalCatalogPlanMetadata *GlobalCatalogPlanMetadata) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogPlanMetadata.RcCompatible) {
		_patch["rc_compatible"] = globalCatalogPlanMetadata.RcCompatible
	}
	if !core.IsNil(globalCatalogPlanMetadata.Ui) {
		_patch["ui"] = globalCatalogPlanMetadata.Ui.asPatch()
	}
	if !core.IsNil(globalCatalogPlanMetadata.Service) {
		_patch["service"] = globalCatalogPlanMetadata.Service.asPatch()
	}
	if !core.IsNil(globalCatalogPlanMetadata.Pricing) {
		_patch["pricing"] = globalCatalogPlanMetadata.Pricing.asPatch()
	}
	if !core.IsNil(globalCatalogPlanMetadata.Plan) {
		_patch["plan"] = globalCatalogPlanMetadata.Plan.asPatch()
	}

	return
}

// GlobalCatalogPlanMetadataPlan : Metadata controlling Plan related settings.
type GlobalCatalogPlanMetadataPlan struct {
	// Controls if IBMers are allowed to provision this plan.
	AllowInternalUsers *bool `json:"allow_internal_users,omitempty"`

	// Deprecated. Controls the Connections tab on the Resource Details page.
	Bindable *bool `json:"bindable,omitempty"`
}

// UnmarshalGlobalCatalogPlanMetadataPlan unmarshals an instance of GlobalCatalogPlanMetadataPlan from the specified map of raw messages.
func UnmarshalGlobalCatalogPlanMetadataPlan(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogPlanMetadataPlan)
	err = core.UnmarshalPrimitive(m, "allow_internal_users", &obj.AllowInternalUsers)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_internal_users-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bindable", &obj.Bindable)
	if err != nil {
		err = core.SDKErrorf(err, "", "bindable-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogPlanMetadataPlan
func (globalCatalogPlanMetadataPlan *GlobalCatalogPlanMetadataPlan) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogPlanMetadataPlan.AllowInternalUsers) {
		_patch["allow_internal_users"] = globalCatalogPlanMetadataPlan.AllowInternalUsers
	}
	if !core.IsNil(globalCatalogPlanMetadataPlan.Bindable) {
		_patch["bindable"] = globalCatalogPlanMetadataPlan.Bindable
	}

	return
}

// GlobalCatalogPlanPatch : The request body for updating a global catalog plan.
type GlobalCatalogPlanPatch struct {
	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// Global catalog plan metadata.
	Metadata *GlobalCatalogPlanMetadata `json:"metadata,omitempty"`
}

// UnmarshalGlobalCatalogPlanPatch unmarshals an instance of GlobalCatalogPlanPatch from the specified map of raw messages.
func UnmarshalGlobalCatalogPlanPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogPlanPatch)
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogPlanMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the GlobalCatalogPlanPatch
func (globalCatalogPlanPatch *GlobalCatalogPlanPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogPlanPatch.Active) {
		_patch["active"] = globalCatalogPlanPatch.Active
	}
	if !core.IsNil(globalCatalogPlanPatch.Disabled) {
		_patch["disabled"] = globalCatalogPlanPatch.Disabled
	}
	if !core.IsNil(globalCatalogPlanPatch.OverviewUi) {
		_patch["overview_ui"] = globalCatalogPlanPatch.OverviewUi.asPatch()
	}
	if !core.IsNil(globalCatalogPlanPatch.Tags) {
		_patch["tags"] = globalCatalogPlanPatch.Tags
	}
	if !core.IsNil(globalCatalogPlanPatch.ObjectProvider) {
		_patch["object_provider"] = globalCatalogPlanPatch.ObjectProvider.asPatch()
	}
	if !core.IsNil(globalCatalogPlanPatch.Metadata) {
		_patch["metadata"] = globalCatalogPlanPatch.Metadata.asPatch()
	}

	return
}

// GlobalCatalogProduct : The object defining the global catalog product.
type GlobalCatalogProduct struct {
	// The ID of a global catalog object.
	ID *string `json:"id,omitempty"`

	// The desired ID of the global catalog object.
	ObjectID *string `json:"object_id,omitempty"`

	// The programmatic name of this product.
	Name *string `json:"name,omitempty"`

	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The kind of the global catalog object.
	Kind *string `json:"kind,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// Images from the global catalog entry that help illustrate the service.
	Images *GlobalCatalogProductImages `json:"images,omitempty"`

	// The global catalog URL of your product.
	URL *string `json:"url,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// The global catalog service metadata object.
	Metadata *GlobalCatalogProductMetadata `json:"metadata,omitempty"`
}

// Constants associated with the GlobalCatalogProduct.Kind property.
// The kind of the global catalog object.
const (
	GlobalCatalogProduct_Kind_Composite       = "composite"
	GlobalCatalogProduct_Kind_PlatformService = "platform_service"
	GlobalCatalogProduct_Kind_Service         = "service"
)

// UnmarshalGlobalCatalogProduct unmarshals an instance of GlobalCatalogProduct from the specified map of raw messages.
func UnmarshalGlobalCatalogProduct(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProduct)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "object_id", &obj.ObjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "images", &obj.Images, UnmarshalGlobalCatalogProductImages)
	if err != nil {
		err = core.SDKErrorf(err, "", "images-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogProductMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GlobalCatalogProductImages : Images from the global catalog entry that help illustrate the service.
type GlobalCatalogProductImages struct {
	// The URL for your product logo.
	Image *string `json:"image,omitempty"`
}

// UnmarshalGlobalCatalogProductImages unmarshals an instance of GlobalCatalogProductImages from the specified map of raw messages.
func UnmarshalGlobalCatalogProductImages(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductImages)
	err = core.UnmarshalPrimitive(m, "image", &obj.Image)
	if err != nil {
		err = core.SDKErrorf(err, "", "image-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductImages
func (globalCatalogProductImages *GlobalCatalogProductImages) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductImages.Image) {
		_patch["image"] = globalCatalogProductImages.Image
	}

	return
}

// GlobalCatalogProductMetadata : The global catalog service metadata object.
type GlobalCatalogProductMetadata struct {
	// Whether the object is compatible with the resource controller service.
	RcCompatible *bool `json:"rc_compatible,omitempty"`

	// The UI metadata of this service.
	Ui *GlobalCatalogMetadataUI `json:"ui,omitempty"`

	// The global catalog metadata of the service.
	Service *GlobalCatalogMetadataService `json:"service,omitempty"`

	// The additional metadata of the service in global catalog.
	Other *GlobalCatalogProductMetadataOther `json:"other,omitempty"`
}

// UnmarshalGlobalCatalogProductMetadata unmarshals an instance of GlobalCatalogProductMetadata from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadata)
	err = core.UnmarshalPrimitive(m, "rc_compatible", &obj.RcCompatible)
	if err != nil {
		err = core.SDKErrorf(err, "", "rc_compatible-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.Ui, UnmarshalGlobalCatalogMetadataUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service", &obj.Service, UnmarshalGlobalCatalogMetadataService)
	if err != nil {
		err = core.SDKErrorf(err, "", "service-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "other", &obj.Other, UnmarshalGlobalCatalogProductMetadataOther)
	if err != nil {
		err = core.SDKErrorf(err, "", "other-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadata
func (globalCatalogProductMetadata *GlobalCatalogProductMetadata) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadata.RcCompatible) {
		_patch["rc_compatible"] = globalCatalogProductMetadata.RcCompatible
	}
	if !core.IsNil(globalCatalogProductMetadata.Ui) {
		_patch["ui"] = globalCatalogProductMetadata.Ui.asPatch()
	}
	if !core.IsNil(globalCatalogProductMetadata.Service) {
		_patch["service"] = globalCatalogProductMetadata.Service.asPatch()
	}
	if !core.IsNil(globalCatalogProductMetadata.Other) {
		_patch["other"] = globalCatalogProductMetadata.Other.asPatch()
	}

	return
}

// GlobalCatalogProductMetadataOther : The additional metadata of the service in global catalog.
type GlobalCatalogProductMetadataOther struct {
	// The metadata of the service owned and managed by Partner Center - Sell.
	PC *GlobalCatalogProductMetadataOtherPC `json:"PC,omitempty"`

	// Optional metadata of the service defining it as a composite.
	Composite *GlobalCatalogProductMetadataOtherComposite `json:"composite,omitempty"`
}

// UnmarshalGlobalCatalogProductMetadataOther unmarshals an instance of GlobalCatalogProductMetadataOther from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadataOther(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadataOther)
	err = core.UnmarshalModel(m, "PC", &obj.PC, UnmarshalGlobalCatalogProductMetadataOtherPC)
	if err != nil {
		err = core.SDKErrorf(err, "", "PC-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "composite", &obj.Composite, UnmarshalGlobalCatalogProductMetadataOtherComposite)
	if err != nil {
		err = core.SDKErrorf(err, "", "composite-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadataOther
func (globalCatalogProductMetadataOther *GlobalCatalogProductMetadataOther) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadataOther.PC) {
		_patch["PC"] = globalCatalogProductMetadataOther.PC.asPatch()
	}
	if !core.IsNil(globalCatalogProductMetadataOther.Composite) {
		_patch["composite"] = globalCatalogProductMetadataOther.Composite.asPatch()
	}

	return
}

// GlobalCatalogProductMetadataOtherComposite : Optional metadata of the service defining it as a composite.
type GlobalCatalogProductMetadataOtherComposite struct {
	// The type of the composite service.
	CompositeKind *string `json:"composite_kind,omitempty"`

	// The tag used for the composite parent and its children.
	CompositeTag *string `json:"composite_tag,omitempty"`

	Children []GlobalCatalogProductMetadataOtherCompositeChild `json:"children,omitempty"`
}

// Constants associated with the GlobalCatalogProductMetadataOtherComposite.CompositeKind property.
// The type of the composite service.
const (
	GlobalCatalogProductMetadataOtherComposite_CompositeKind_PlatformService = "platform_service"
	GlobalCatalogProductMetadataOtherComposite_CompositeKind_Service         = "service"
)

// UnmarshalGlobalCatalogProductMetadataOtherComposite unmarshals an instance of GlobalCatalogProductMetadataOtherComposite from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadataOtherComposite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadataOtherComposite)
	err = core.UnmarshalPrimitive(m, "composite_kind", &obj.CompositeKind)
	if err != nil {
		err = core.SDKErrorf(err, "", "composite_kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "composite_tag", &obj.CompositeTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "composite_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "children", &obj.Children, UnmarshalGlobalCatalogProductMetadataOtherCompositeChild)
	if err != nil {
		err = core.SDKErrorf(err, "", "children-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadataOtherComposite
func (globalCatalogProductMetadataOtherComposite *GlobalCatalogProductMetadataOtherComposite) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadataOtherComposite.CompositeKind) {
		_patch["composite_kind"] = globalCatalogProductMetadataOtherComposite.CompositeKind
	}
	if !core.IsNil(globalCatalogProductMetadataOtherComposite.CompositeTag) {
		_patch["composite_tag"] = globalCatalogProductMetadataOtherComposite.CompositeTag
	}
	if !core.IsNil(globalCatalogProductMetadataOtherComposite.Children) {
		var childrenPatches []map[string]interface{}
		for _, children := range globalCatalogProductMetadataOtherComposite.Children {
			childrenPatches = append(childrenPatches, children.asPatch())
		}
		_patch["children"] = childrenPatches
	}

	return
}

// GlobalCatalogProductMetadataOtherCompositeChild : Object defining a composite child of a composite parent.
type GlobalCatalogProductMetadataOtherCompositeChild struct {
	// The type of the composite child.
	Kind *string `json:"kind,omitempty"`

	// The name of the composite child.
	Name *string `json:"name,omitempty"`
}

// Constants associated with the GlobalCatalogProductMetadataOtherCompositeChild.Kind property.
// The type of the composite child.
const (
	GlobalCatalogProductMetadataOtherCompositeChild_Kind_PlatformService = "platform_service"
	GlobalCatalogProductMetadataOtherCompositeChild_Kind_Service         = "service"
)

// UnmarshalGlobalCatalogProductMetadataOtherCompositeChild unmarshals an instance of GlobalCatalogProductMetadataOtherCompositeChild from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadataOtherCompositeChild(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadataOtherCompositeChild)
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadataOtherCompositeChild
func (globalCatalogProductMetadataOtherCompositeChild *GlobalCatalogProductMetadataOtherCompositeChild) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadataOtherCompositeChild.Kind) {
		_patch["kind"] = globalCatalogProductMetadataOtherCompositeChild.Kind
	}
	if !core.IsNil(globalCatalogProductMetadataOtherCompositeChild.Name) {
		_patch["name"] = globalCatalogProductMetadataOtherCompositeChild.Name
	}

	return
}

// GlobalCatalogProductMetadataOtherPC : The metadata of the service owned and managed by Partner Center - Sell.
type GlobalCatalogProductMetadataOtherPC struct {
	// The support metadata of the service.
	Support *GlobalCatalogProductMetadataOtherPCSupport `json:"support,omitempty"`
}

// UnmarshalGlobalCatalogProductMetadataOtherPC unmarshals an instance of GlobalCatalogProductMetadataOtherPC from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadataOtherPC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadataOtherPC)
	err = core.UnmarshalModel(m, "support", &obj.Support, UnmarshalGlobalCatalogProductMetadataOtherPCSupport)
	if err != nil {
		err = core.SDKErrorf(err, "", "support-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadataOtherPC
func (globalCatalogProductMetadataOtherPC *GlobalCatalogProductMetadataOtherPC) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadataOtherPC.Support) {
		_patch["support"] = globalCatalogProductMetadataOtherPC.Support.asPatch()
	}

	return
}

// GlobalCatalogProductMetadataOtherPCSupport : The support metadata of the service.
type GlobalCatalogProductMetadataOtherPCSupport struct {
	// The support site URL where the support for your service is available.
	URL *string `json:"url,omitempty"`

	// The URL where the status of your service is available.
	StatusURL *string `json:"status_url,omitempty"`

	// The countries in which your support is available. Provide a list of country codes.
	Locations []string `json:"locations,omitempty"`

	// The languages in which support is available.
	Languages []string `json:"languages,omitempty"`

	// The description of your support process.
	Process *string `json:"process,omitempty"`

	// The description of your support process.
	ProcessI18n map[string]string `json:"process_i18n,omitempty"`

	// The type of support provided.
	SupportType *string `json:"support_type,omitempty"`

	// The details of the support escalation process.
	SupportEscalation *SupportEscalation `json:"support_escalation,omitempty"`

	// The support options for the service.
	SupportDetails []SupportDetailsItem `json:"support_details,omitempty"`
}

// Constants associated with the GlobalCatalogProductMetadataOtherPCSupport.SupportType property.
// The type of support provided.
const (
	GlobalCatalogProductMetadataOtherPCSupport_SupportType_Community  = "community"
	GlobalCatalogProductMetadataOtherPCSupport_SupportType_Ibm        = "ibm"
	GlobalCatalogProductMetadataOtherPCSupport_SupportType_IbmCloud   = "ibm_cloud"
	GlobalCatalogProductMetadataOtherPCSupport_SupportType_ThirdParty = "third_party"
)

// UnmarshalGlobalCatalogProductMetadataOtherPCSupport unmarshals an instance of GlobalCatalogProductMetadataOtherPCSupport from the specified map of raw messages.
func UnmarshalGlobalCatalogProductMetadataOtherPCSupport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductMetadataOtherPCSupport)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_url", &obj.StatusURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locations", &obj.Locations)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "languages", &obj.Languages)
	if err != nil {
		err = core.SDKErrorf(err, "", "languages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "process", &obj.Process)
	if err != nil {
		err = core.SDKErrorf(err, "", "process-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "process_i18n", &obj.ProcessI18n)
	if err != nil {
		err = core.SDKErrorf(err, "", "process_i18n-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "support_type", &obj.SupportType)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "support_escalation", &obj.SupportEscalation, UnmarshalSupportEscalation)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_escalation-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "support_details", &obj.SupportDetails, UnmarshalSupportDetailsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "support_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GlobalCatalogProductMetadataOtherPCSupport
func (globalCatalogProductMetadataOtherPCSupport *GlobalCatalogProductMetadataOtherPCSupport) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.URL) {
		_patch["url"] = globalCatalogProductMetadataOtherPCSupport.URL
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.StatusURL) {
		_patch["status_url"] = globalCatalogProductMetadataOtherPCSupport.StatusURL
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.Locations) {
		_patch["locations"] = globalCatalogProductMetadataOtherPCSupport.Locations
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.Languages) {
		_patch["languages"] = globalCatalogProductMetadataOtherPCSupport.Languages
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.Process) {
		_patch["process"] = globalCatalogProductMetadataOtherPCSupport.Process
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.ProcessI18n) {
		_patch["process_i18n"] = globalCatalogProductMetadataOtherPCSupport.ProcessI18n
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.SupportType) {
		_patch["support_type"] = globalCatalogProductMetadataOtherPCSupport.SupportType
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.SupportEscalation) {
		_patch["support_escalation"] = globalCatalogProductMetadataOtherPCSupport.SupportEscalation.asPatch()
	}
	if !core.IsNil(globalCatalogProductMetadataOtherPCSupport.SupportDetails) {
		var supportDetailsPatches []map[string]interface{}
		for _, supportDetails := range globalCatalogProductMetadataOtherPCSupport.SupportDetails {
			supportDetailsPatches = append(supportDetailsPatches, supportDetails.asPatch())
		}
		_patch["support_details"] = supportDetailsPatches
	}

	return
}

// GlobalCatalogProductPatch : The request body for updating a product in global catalog.
type GlobalCatalogProductPatch struct {
	// Whether the service is active.
	Active *bool `json:"active,omitempty"`

	// Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are
	// disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// The object that contains the service details from the Overview page in global catalog.
	OverviewUi *GlobalCatalogOverviewUI `json:"overview_ui,omitempty"`

	// A list of tags that carry information about your product. These tags can be used to find your product in the IBM
	// Cloud catalog.
	Tags []string `json:"tags,omitempty"`

	// Images from the global catalog entry that help illustrate the service.
	Images *GlobalCatalogProductImages `json:"images,omitempty"`

	// The provider or owner of the product.
	ObjectProvider *CatalogProductProvider `json:"object_provider,omitempty"`

	// The global catalog service metadata object.
	Metadata *GlobalCatalogProductMetadata `json:"metadata,omitempty"`
}

// UnmarshalGlobalCatalogProductPatch unmarshals an instance of GlobalCatalogProductPatch from the specified map of raw messages.
func UnmarshalGlobalCatalogProductPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GlobalCatalogProductPatch)
	err = core.UnmarshalPrimitive(m, "active", &obj.Active)
	if err != nil {
		err = core.SDKErrorf(err, "", "active-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "overview_ui", &obj.OverviewUi, UnmarshalGlobalCatalogOverviewUI)
	if err != nil {
		err = core.SDKErrorf(err, "", "overview_ui-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "images", &obj.Images, UnmarshalGlobalCatalogProductImages)
	if err != nil {
		err = core.SDKErrorf(err, "", "images-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "object_provider", &obj.ObjectProvider, UnmarshalCatalogProductProvider)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_provider-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalGlobalCatalogProductMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the GlobalCatalogProductPatch
func (globalCatalogProductPatch *GlobalCatalogProductPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(globalCatalogProductPatch.Active) {
		_patch["active"] = globalCatalogProductPatch.Active
	}
	if !core.IsNil(globalCatalogProductPatch.Disabled) {
		_patch["disabled"] = globalCatalogProductPatch.Disabled
	}
	if !core.IsNil(globalCatalogProductPatch.OverviewUi) {
		_patch["overview_ui"] = globalCatalogProductPatch.OverviewUi.asPatch()
	}
	if !core.IsNil(globalCatalogProductPatch.Tags) {
		_patch["tags"] = globalCatalogProductPatch.Tags
	}
	if !core.IsNil(globalCatalogProductPatch.Images) {
		_patch["images"] = globalCatalogProductPatch.Images.asPatch()
	}
	if !core.IsNil(globalCatalogProductPatch.ObjectProvider) {
		_patch["object_provider"] = globalCatalogProductPatch.ObjectProvider.asPatch()
	}
	if !core.IsNil(globalCatalogProductPatch.Metadata) {
		_patch["metadata"] = globalCatalogProductPatch.Metadata.asPatch()
	}

	return
}

// IamServiceRegistration : The IAM service registration.
type IamServiceRegistration struct {
	// The IAM registration name, which must be the programmatic name of the product.
	Name *string `json:"name" validate:"required"`

	// Whether the service is enabled or disabled for IAM.
	Enabled *bool `json:"enabled,omitempty"`

	// The type of the service.
	ServiceType *string `json:"service_type,omitempty"`

	// The product access management action.
	Actions []IamServiceRegistrationAction `json:"actions,omitempty"`

	// List of additional policy scopes.
	AdditionalPolicyScopes []string `json:"additional_policy_scopes,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// The list of parent IDs for product access management.
	ParentIds []string `json:"parent_ids,omitempty"`

	// The resource hierarchy key-value pair for composite services.
	ResourceHierarchyAttribute *IamServiceRegistrationResourceHierarchyAttribute `json:"resource_hierarchy_attribute,omitempty"`

	// The list of supported anonymous accesses.
	SupportedAnonymousAccesses []IamServiceRegistrationSupportedAnonymousAccess `json:"supported_anonymous_accesses,omitempty"`

	// The list of supported attributes.
	SupportedAttributes []IamServiceRegistrationSupportedAttribute `json:"supported_attributes,omitempty"`

	// The list of supported authorization subjects.
	SupportedAuthorizationSubjects []IamServiceRegistrationSupportedAuthorizationSubject `json:"supported_authorization_subjects,omitempty"`

	// The list of roles that you can use to assign access.
	SupportedRoles []IamServiceRegistrationSupportedRole `json:"supported_roles,omitempty"`

	// The registration of set of endpoint types that are supported by your service in the `networkType` environment
	// attribute. This constrains the context-based restriction rules specific to the service such that they describe
	// access restrictions on only this set of endpoints.
	SupportedNetwork *IamServiceRegistrationSupportedNetwork `json:"supported_network,omitempty"`
}

// Constants associated with the IamServiceRegistration.ServiceType property.
// The type of the service.
const (
	IamServiceRegistration_ServiceType_PlatformService = "platform_service"
	IamServiceRegistration_ServiceType_Service         = "service"
)

// UnmarshalIamServiceRegistration unmarshals an instance of IamServiceRegistration from the specified map of raw messages.
func UnmarshalIamServiceRegistration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistration)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_type", &obj.ServiceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalIamServiceRegistrationAction)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "additional_policy_scopes", &obj.AdditionalPolicyScopes)
	if err != nil {
		err = core.SDKErrorf(err, "", "additional_policy_scopes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent_ids", &obj.ParentIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_hierarchy_attribute", &obj.ResourceHierarchyAttribute, UnmarshalIamServiceRegistrationResourceHierarchyAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_hierarchy_attribute-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_anonymous_accesses", &obj.SupportedAnonymousAccesses, UnmarshalIamServiceRegistrationSupportedAnonymousAccess)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_anonymous_accesses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_attributes", &obj.SupportedAttributes, UnmarshalIamServiceRegistrationSupportedAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_authorization_subjects", &obj.SupportedAuthorizationSubjects, UnmarshalIamServiceRegistrationSupportedAuthorizationSubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_authorization_subjects-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_roles", &obj.SupportedRoles, UnmarshalIamServiceRegistrationSupportedRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_network", &obj.SupportedNetwork, UnmarshalIamServiceRegistrationSupportedNetwork)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_network-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// IamServiceRegistrationAction : Action.
type IamServiceRegistrationAction struct {
	// The unique identifier for the action.
	ID *string `json:"id,omitempty"`

	// The list of roles for the action.
	Roles []string `json:"roles,omitempty"`

	// The description for the object.
	Description *IamServiceRegistrationDescriptionObject `json:"description,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// Extra options.
	Options *IamServiceRegistrationActionOptions `json:"options,omitempty"`
}

// UnmarshalIamServiceRegistrationAction unmarshals an instance of IamServiceRegistrationAction from the specified map of raw messages.
func UnmarshalIamServiceRegistrationAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationAction)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "description", &obj.Description, UnmarshalIamServiceRegistrationDescriptionObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "options", &obj.Options, UnmarshalIamServiceRegistrationActionOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "options-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationAction
func (iamServiceRegistrationAction *IamServiceRegistrationAction) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationAction.ID) {
		_patch["id"] = iamServiceRegistrationAction.ID
	}
	if !core.IsNil(iamServiceRegistrationAction.Roles) {
		_patch["roles"] = iamServiceRegistrationAction.Roles
	}
	if !core.IsNil(iamServiceRegistrationAction.Description) {
		_patch["description"] = iamServiceRegistrationAction.Description.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationAction.DisplayName) {
		_patch["display_name"] = iamServiceRegistrationAction.DisplayName.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationAction.Options) {
		_patch["options"] = iamServiceRegistrationAction.Options.asPatch()
	}

	return
}

// IamServiceRegistrationActionOptions : Extra options.
type IamServiceRegistrationActionOptions struct {
	// Optional opt-in if action is hidden from customers.
	Hidden *bool `json:"hidden,omitempty"`
}

// UnmarshalIamServiceRegistrationActionOptions unmarshals an instance of IamServiceRegistrationActionOptions from the specified map of raw messages.
func UnmarshalIamServiceRegistrationActionOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationActionOptions)
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationActionOptions
func (iamServiceRegistrationActionOptions *IamServiceRegistrationActionOptions) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationActionOptions.Hidden) {
		_patch["hidden"] = iamServiceRegistrationActionOptions.Hidden
	}

	return
}

// IamServiceRegistrationDescriptionObject : The description for the object.
type IamServiceRegistrationDescriptionObject struct {
	// The fallback string for the description object.
	Default *string `json:"default,omitempty"`

	// English.
	En *string `json:"en,omitempty"`

	// German.
	De *string `json:"de,omitempty"`

	// Spanish.
	Es *string `json:"es,omitempty"`

	// French.
	Fr *string `json:"fr,omitempty"`

	// Italian.
	It *string `json:"it,omitempty"`

	// Japanese.
	Ja *string `json:"ja,omitempty"`

	// Korean.
	Ko *string `json:"ko,omitempty"`

	// Portuguese (Brazil).
	PtBr *string `json:"pt_br,omitempty"`

	// Traditional Chinese.
	ZhTw *string `json:"zh_tw,omitempty"`

	// Simplified Chinese.
	ZhCn *string `json:"zh_cn,omitempty"`
}

// UnmarshalIamServiceRegistrationDescriptionObject unmarshals an instance of IamServiceRegistrationDescriptionObject from the specified map of raw messages.
func UnmarshalIamServiceRegistrationDescriptionObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationDescriptionObject)
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		err = core.SDKErrorf(err, "", "default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "en", &obj.En)
	if err != nil {
		err = core.SDKErrorf(err, "", "en-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "de", &obj.De)
	if err != nil {
		err = core.SDKErrorf(err, "", "de-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "es", &obj.Es)
	if err != nil {
		err = core.SDKErrorf(err, "", "es-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fr", &obj.Fr)
	if err != nil {
		err = core.SDKErrorf(err, "", "fr-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "it", &obj.It)
	if err != nil {
		err = core.SDKErrorf(err, "", "it-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ja", &obj.Ja)
	if err != nil {
		err = core.SDKErrorf(err, "", "ja-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ko", &obj.Ko)
	if err != nil {
		err = core.SDKErrorf(err, "", "ko-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pt_br", &obj.PtBr)
	if err != nil {
		err = core.SDKErrorf(err, "", "pt_br-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zh_tw", &obj.ZhTw)
	if err != nil {
		err = core.SDKErrorf(err, "", "zh_tw-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zh_cn", &obj.ZhCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "zh_cn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationDescriptionObject
func (iamServiceRegistrationDescriptionObject *IamServiceRegistrationDescriptionObject) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.Default) {
		_patch["default"] = iamServiceRegistrationDescriptionObject.Default
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.En) {
		_patch["en"] = iamServiceRegistrationDescriptionObject.En
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.De) {
		_patch["de"] = iamServiceRegistrationDescriptionObject.De
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.Es) {
		_patch["es"] = iamServiceRegistrationDescriptionObject.Es
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.Fr) {
		_patch["fr"] = iamServiceRegistrationDescriptionObject.Fr
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.It) {
		_patch["it"] = iamServiceRegistrationDescriptionObject.It
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.Ja) {
		_patch["ja"] = iamServiceRegistrationDescriptionObject.Ja
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.Ko) {
		_patch["ko"] = iamServiceRegistrationDescriptionObject.Ko
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.PtBr) {
		_patch["pt_br"] = iamServiceRegistrationDescriptionObject.PtBr
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.ZhTw) {
		_patch["zh_tw"] = iamServiceRegistrationDescriptionObject.ZhTw
	}
	if !core.IsNil(iamServiceRegistrationDescriptionObject.ZhCn) {
		_patch["zh_cn"] = iamServiceRegistrationDescriptionObject.ZhCn
	}

	return
}

// IamServiceRegistrationDisplayNameObject : The display name of the object.
type IamServiceRegistrationDisplayNameObject struct {
	// The fallback string for the description object.
	Default *string `json:"default,omitempty"`

	// English.
	En *string `json:"en,omitempty"`

	// German.
	De *string `json:"de,omitempty"`

	// Spanish.
	Es *string `json:"es,omitempty"`

	// French.
	Fr *string `json:"fr,omitempty"`

	// Italian.
	It *string `json:"it,omitempty"`

	// Japanese.
	Ja *string `json:"ja,omitempty"`

	// Korean.
	Ko *string `json:"ko,omitempty"`

	// Portuguese (Brazil).
	PtBr *string `json:"pt_br,omitempty"`

	// Traditional Chinese.
	ZhTw *string `json:"zh_tw,omitempty"`

	// Simplified Chinese.
	ZhCn *string `json:"zh_cn,omitempty"`
}

// UnmarshalIamServiceRegistrationDisplayNameObject unmarshals an instance of IamServiceRegistrationDisplayNameObject from the specified map of raw messages.
func UnmarshalIamServiceRegistrationDisplayNameObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationDisplayNameObject)
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		err = core.SDKErrorf(err, "", "default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "en", &obj.En)
	if err != nil {
		err = core.SDKErrorf(err, "", "en-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "de", &obj.De)
	if err != nil {
		err = core.SDKErrorf(err, "", "de-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "es", &obj.Es)
	if err != nil {
		err = core.SDKErrorf(err, "", "es-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fr", &obj.Fr)
	if err != nil {
		err = core.SDKErrorf(err, "", "fr-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "it", &obj.It)
	if err != nil {
		err = core.SDKErrorf(err, "", "it-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ja", &obj.Ja)
	if err != nil {
		err = core.SDKErrorf(err, "", "ja-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ko", &obj.Ko)
	if err != nil {
		err = core.SDKErrorf(err, "", "ko-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pt_br", &obj.PtBr)
	if err != nil {
		err = core.SDKErrorf(err, "", "pt_br-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zh_tw", &obj.ZhTw)
	if err != nil {
		err = core.SDKErrorf(err, "", "zh_tw-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zh_cn", &obj.ZhCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "zh_cn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationDisplayNameObject
func (iamServiceRegistrationDisplayNameObject *IamServiceRegistrationDisplayNameObject) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.Default) {
		_patch["default"] = iamServiceRegistrationDisplayNameObject.Default
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.En) {
		_patch["en"] = iamServiceRegistrationDisplayNameObject.En
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.De) {
		_patch["de"] = iamServiceRegistrationDisplayNameObject.De
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.Es) {
		_patch["es"] = iamServiceRegistrationDisplayNameObject.Es
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.Fr) {
		_patch["fr"] = iamServiceRegistrationDisplayNameObject.Fr
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.It) {
		_patch["it"] = iamServiceRegistrationDisplayNameObject.It
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.Ja) {
		_patch["ja"] = iamServiceRegistrationDisplayNameObject.Ja
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.Ko) {
		_patch["ko"] = iamServiceRegistrationDisplayNameObject.Ko
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.PtBr) {
		_patch["pt_br"] = iamServiceRegistrationDisplayNameObject.PtBr
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.ZhTw) {
		_patch["zh_tw"] = iamServiceRegistrationDisplayNameObject.ZhTw
	}
	if !core.IsNil(iamServiceRegistrationDisplayNameObject.ZhCn) {
		_patch["zh_cn"] = iamServiceRegistrationDisplayNameObject.ZhCn
	}

	return
}

// IamServiceRegistrationPatch : The patch object of an IAM service registration.
type IamServiceRegistrationPatch struct {
	// Whether the service is enabled or disabled for IAM.
	Enabled *bool `json:"enabled,omitempty"`

	// The type of the service.
	ServiceType *string `json:"service_type,omitempty"`

	// The product access management action.
	Actions []IamServiceRegistrationAction `json:"actions,omitempty"`

	// List of additional policy scopes.
	AdditionalPolicyScopes []string `json:"additional_policy_scopes,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// The list of parent IDs for product access management.
	ParentIds []string `json:"parent_ids,omitempty"`

	// The resource hierarchy key-value pair for composite services.
	ResourceHierarchyAttribute *IamServiceRegistrationResourceHierarchyAttribute `json:"resource_hierarchy_attribute,omitempty"`

	// The list of supported anonymous accesses.
	SupportedAnonymousAccesses []IamServiceRegistrationSupportedAnonymousAccess `json:"supported_anonymous_accesses,omitempty"`

	// The list of supported attributes.
	SupportedAttributes []IamServiceRegistrationSupportedAttribute `json:"supported_attributes,omitempty"`

	// The list of supported authorization subjects.
	SupportedAuthorizationSubjects []IamServiceRegistrationSupportedAuthorizationSubject `json:"supported_authorization_subjects,omitempty"`

	// The list of roles that you can use to assign access.
	SupportedRoles []IamServiceRegistrationSupportedRole `json:"supported_roles,omitempty"`

	// The registration of set of endpoint types that are supported by your service in the `networkType` environment
	// attribute. This constrains the context-based restriction rules specific to the service such that they describe
	// access restrictions on only this set of endpoints.
	SupportedNetwork *IamServiceRegistrationSupportedNetwork `json:"supported_network,omitempty"`
}

// Constants associated with the IamServiceRegistrationPatch.ServiceType property.
// The type of the service.
const (
	IamServiceRegistrationPatch_ServiceType_PlatformService = "platform_service"
	IamServiceRegistrationPatch_ServiceType_Service         = "service"
)

// UnmarshalIamServiceRegistrationPatch unmarshals an instance of IamServiceRegistrationPatch from the specified map of raw messages.
func UnmarshalIamServiceRegistrationPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationPatch)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_type", &obj.ServiceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalIamServiceRegistrationAction)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "additional_policy_scopes", &obj.AdditionalPolicyScopes)
	if err != nil {
		err = core.SDKErrorf(err, "", "additional_policy_scopes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parent_ids", &obj.ParentIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "parent_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_hierarchy_attribute", &obj.ResourceHierarchyAttribute, UnmarshalIamServiceRegistrationResourceHierarchyAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_hierarchy_attribute-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_anonymous_accesses", &obj.SupportedAnonymousAccesses, UnmarshalIamServiceRegistrationSupportedAnonymousAccess)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_anonymous_accesses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_attributes", &obj.SupportedAttributes, UnmarshalIamServiceRegistrationSupportedAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_authorization_subjects", &obj.SupportedAuthorizationSubjects, UnmarshalIamServiceRegistrationSupportedAuthorizationSubject)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_authorization_subjects-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_roles", &obj.SupportedRoles, UnmarshalIamServiceRegistrationSupportedRole)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_roles-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "supported_network", &obj.SupportedNetwork, UnmarshalIamServiceRegistrationSupportedNetwork)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_network-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the IamServiceRegistrationPatch
func (iamServiceRegistrationPatch *IamServiceRegistrationPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationPatch.Enabled) {
		_patch["enabled"] = iamServiceRegistrationPatch.Enabled
	}
	if !core.IsNil(iamServiceRegistrationPatch.ServiceType) {
		_patch["service_type"] = iamServiceRegistrationPatch.ServiceType
	}
	if !core.IsNil(iamServiceRegistrationPatch.Actions) {
		var actionsPatches []map[string]interface{}
		for _, actions := range iamServiceRegistrationPatch.Actions {
			actionsPatches = append(actionsPatches, actions.asPatch())
		}
		_patch["actions"] = actionsPatches
	}
	if !core.IsNil(iamServiceRegistrationPatch.AdditionalPolicyScopes) {
		_patch["additional_policy_scopes"] = iamServiceRegistrationPatch.AdditionalPolicyScopes
	}
	if !core.IsNil(iamServiceRegistrationPatch.DisplayName) {
		_patch["display_name"] = iamServiceRegistrationPatch.DisplayName.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationPatch.ParentIds) {
		_patch["parent_ids"] = iamServiceRegistrationPatch.ParentIds
	}
	if !core.IsNil(iamServiceRegistrationPatch.ResourceHierarchyAttribute) {
		_patch["resource_hierarchy_attribute"] = iamServiceRegistrationPatch.ResourceHierarchyAttribute.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationPatch.SupportedAnonymousAccesses) {
		var supportedAnonymousAccessesPatches []map[string]interface{}
		for _, supportedAnonymousAccesses := range iamServiceRegistrationPatch.SupportedAnonymousAccesses {
			supportedAnonymousAccessesPatches = append(supportedAnonymousAccessesPatches, supportedAnonymousAccesses.asPatch())
		}
		_patch["supported_anonymous_accesses"] = supportedAnonymousAccessesPatches
	}
	if !core.IsNil(iamServiceRegistrationPatch.SupportedAttributes) {
		var supportedAttributesPatches []map[string]interface{}
		for _, supportedAttributes := range iamServiceRegistrationPatch.SupportedAttributes {
			supportedAttributesPatches = append(supportedAttributesPatches, supportedAttributes.asPatch())
		}
		_patch["supported_attributes"] = supportedAttributesPatches
	}
	if !core.IsNil(iamServiceRegistrationPatch.SupportedAuthorizationSubjects) {
		var supportedAuthorizationSubjectsPatches []map[string]interface{}
		for _, supportedAuthorizationSubjects := range iamServiceRegistrationPatch.SupportedAuthorizationSubjects {
			supportedAuthorizationSubjectsPatches = append(supportedAuthorizationSubjectsPatches, supportedAuthorizationSubjects.asPatch())
		}
		_patch["supported_authorization_subjects"] = supportedAuthorizationSubjectsPatches
	}
	if !core.IsNil(iamServiceRegistrationPatch.SupportedRoles) {
		var supportedRolesPatches []map[string]interface{}
		for _, supportedRoles := range iamServiceRegistrationPatch.SupportedRoles {
			supportedRolesPatches = append(supportedRolesPatches, supportedRoles.asPatch())
		}
		_patch["supported_roles"] = supportedRolesPatches
	}
	if !core.IsNil(iamServiceRegistrationPatch.SupportedNetwork) {
		_patch["supported_network"] = iamServiceRegistrationPatch.SupportedNetwork.asPatch()
	}

	return
}

// IamServiceRegistrationResourceHierarchyAttribute : The resource hierarchy key-value pair for composite services.
type IamServiceRegistrationResourceHierarchyAttribute struct {
	// The resource hierarchy key.
	Key *string `json:"key,omitempty"`

	// The resource hierarchy value.
	Value *string `json:"value,omitempty"`
}

// UnmarshalIamServiceRegistrationResourceHierarchyAttribute unmarshals an instance of IamServiceRegistrationResourceHierarchyAttribute from the specified map of raw messages.
func UnmarshalIamServiceRegistrationResourceHierarchyAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationResourceHierarchyAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationResourceHierarchyAttribute
func (iamServiceRegistrationResourceHierarchyAttribute *IamServiceRegistrationResourceHierarchyAttribute) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationResourceHierarchyAttribute.Key) {
		_patch["key"] = iamServiceRegistrationResourceHierarchyAttribute.Key
	}
	if !core.IsNil(iamServiceRegistrationResourceHierarchyAttribute.Value) {
		_patch["value"] = iamServiceRegistrationResourceHierarchyAttribute.Value
	}

	return
}

// IamServiceRegistrationSupportedAnonymousAccess : Resources within the service that supports anonymous access. Attributes defined in here must be included in anonymous
// access policies for this service.
type IamServiceRegistrationSupportedAnonymousAccess struct {
	// The attributes for anonymous accesses.
	Attributes *IamServiceRegistrationSupportedAnonymousAccessAttributes `json:"attributes,omitempty"`

	// The roles of supported anonymous accesses.
	Roles []string `json:"roles,omitempty"`
}

// UnmarshalIamServiceRegistrationSupportedAnonymousAccess unmarshals an instance of IamServiceRegistrationSupportedAnonymousAccess from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedAnonymousAccess(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedAnonymousAccess)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalIamServiceRegistrationSupportedAnonymousAccessAttributes)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedAnonymousAccess
func (iamServiceRegistrationSupportedAnonymousAccess *IamServiceRegistrationSupportedAnonymousAccess) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccess.Attributes) {
		_patch["attributes"] = iamServiceRegistrationSupportedAnonymousAccess.Attributes.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccess.Roles) {
		_patch["roles"] = iamServiceRegistrationSupportedAnonymousAccess.Roles
	}

	return
}

// IamServiceRegistrationSupportedAnonymousAccessAttributes : The attributes for anonymous accesses.
// This type supports additional properties of type *string.
type IamServiceRegistrationSupportedAnonymousAccessAttributes struct {
	// An account id.
	AccountID *string `json:"account_id" validate:"required"`

	// The name of the service.
	ServiceName *string `json:"service_name" validate:"required"`

	// Additional properties the key must come from supported attributes.
	AdditionalProperties map[string]string `json:"additional_properties" validate:"required"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

// NewIamServiceRegistrationSupportedAnonymousAccessAttributes : Instantiate IamServiceRegistrationSupportedAnonymousAccessAttributes (Generic Model Constructor)
func (*PartnerCenterSellV1) NewIamServiceRegistrationSupportedAnonymousAccessAttributes(accountID string, serviceName string, additionalProperties map[string]string) (_model *IamServiceRegistrationSupportedAnonymousAccessAttributes, err error) {
	_model = &IamServiceRegistrationSupportedAnonymousAccessAttributes{
		AccountID:            core.StringPtr(accountID),
		ServiceName:          core.StringPtr(serviceName),
		AdditionalProperties: additionalProperties,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// SetProperty allows the user to set an arbitrary property on an instance of IamServiceRegistrationSupportedAnonymousAccessAttributes.
func (o *IamServiceRegistrationSupportedAnonymousAccessAttributes) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of IamServiceRegistrationSupportedAnonymousAccessAttributes.
func (o *IamServiceRegistrationSupportedAnonymousAccessAttributes) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of IamServiceRegistrationSupportedAnonymousAccessAttributes.
func (o *IamServiceRegistrationSupportedAnonymousAccessAttributes) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of IamServiceRegistrationSupportedAnonymousAccessAttributes.
func (o *IamServiceRegistrationSupportedAnonymousAccessAttributes) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of IamServiceRegistrationSupportedAnonymousAccessAttributes
func (o *IamServiceRegistrationSupportedAnonymousAccessAttributes) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.AccountID != nil {
		m["account_id"] = o.AccountID
	}
	if o.ServiceName != nil {
		m["service_name"] = o.ServiceName
	}
	if o.AdditionalProperties != nil {
		m["additional_properties"] = o.AdditionalProperties
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalIamServiceRegistrationSupportedAnonymousAccessAttributes unmarshals an instance of IamServiceRegistrationSupportedAnonymousAccessAttributes from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedAnonymousAccessAttributes(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedAnonymousAccessAttributes)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	delete(m, "account_id")
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	delete(m, "service_name")
	err = core.UnmarshalPrimitive(m, "additional_properties", &obj.AdditionalProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "additional_properties-error", common.GetComponentInfo())
		return
	}
	delete(m, "additional_properties")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedAnonymousAccessAttributes
func (iamServiceRegistrationSupportedAnonymousAccessAttributes *IamServiceRegistrationSupportedAnonymousAccessAttributes) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccessAttributes.additionalProperties) {
		for k, v := range iamServiceRegistrationSupportedAnonymousAccessAttributes.additionalProperties {
			_patch[k] = v
		}
	}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccessAttributes.AccountID) {
		_patch["account_id"] = iamServiceRegistrationSupportedAnonymousAccessAttributes.AccountID
	}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccessAttributes.ServiceName) {
		_patch["service_name"] = iamServiceRegistrationSupportedAnonymousAccessAttributes.ServiceName
	}
	if !core.IsNil(iamServiceRegistrationSupportedAnonymousAccessAttributes.AdditionalProperties) {
		_patch["additional_properties"] = iamServiceRegistrationSupportedAnonymousAccessAttributes.AdditionalProperties
	}

	return
}

// IamServiceRegistrationSupportedAttribute : The list of supported attributes.
type IamServiceRegistrationSupportedAttribute struct {
	// The supported attribute key.
	Key *string `json:"key,omitempty"`

	// The list of support attribute options.
	Options *SupportedAttributesOptions `json:"options,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// The description for the object.
	Description *IamServiceRegistrationDescriptionObject `json:"description,omitempty"`

	// The user interface.
	Ui *SupportedAttributeUi `json:"ui,omitempty"`
}

// UnmarshalIamServiceRegistrationSupportedAttribute unmarshals an instance of IamServiceRegistrationSupportedAttribute from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedAttribute)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "options", &obj.Options, UnmarshalSupportedAttributesOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "description", &obj.Description, UnmarshalIamServiceRegistrationDescriptionObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ui", &obj.Ui, UnmarshalSupportedAttributeUi)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedAttribute
func (iamServiceRegistrationSupportedAttribute *IamServiceRegistrationSupportedAttribute) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedAttribute.Key) {
		_patch["key"] = iamServiceRegistrationSupportedAttribute.Key
	}
	if !core.IsNil(iamServiceRegistrationSupportedAttribute.Options) {
		_patch["options"] = iamServiceRegistrationSupportedAttribute.Options.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedAttribute.DisplayName) {
		_patch["display_name"] = iamServiceRegistrationSupportedAttribute.DisplayName.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedAttribute.Description) {
		_patch["description"] = iamServiceRegistrationSupportedAttribute.Description.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedAttribute.Ui) {
		_patch["ui"] = iamServiceRegistrationSupportedAttribute.Ui.asPatch()
	}

	return
}

// IamServiceRegistrationSupportedAuthorizationSubject : IamServiceRegistrationSupportedAuthorizationSubject struct
type IamServiceRegistrationSupportedAuthorizationSubject struct {
	// The list of supported authorization subject properties.
	Attributes *SupportAuthorizationSubjectAttribute `json:"attributes,omitempty"`

	// The list of roles for authorization.
	Roles []string `json:"roles,omitempty"`
}

// UnmarshalIamServiceRegistrationSupportedAuthorizationSubject unmarshals an instance of IamServiceRegistrationSupportedAuthorizationSubject from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedAuthorizationSubject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedAuthorizationSubject)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalSupportAuthorizationSubjectAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "roles", &obj.Roles)
	if err != nil {
		err = core.SDKErrorf(err, "", "roles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedAuthorizationSubject
func (iamServiceRegistrationSupportedAuthorizationSubject *IamServiceRegistrationSupportedAuthorizationSubject) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedAuthorizationSubject.Attributes) {
		_patch["attributes"] = iamServiceRegistrationSupportedAuthorizationSubject.Attributes.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedAuthorizationSubject.Roles) {
		_patch["roles"] = iamServiceRegistrationSupportedAuthorizationSubject.Roles
	}

	return
}

// IamServiceRegistrationSupportedNetwork : The registration of set of endpoint types that are supported by your service in the `networkType` environment
// attribute. This constrains the context-based restriction rules specific to the service such that they describe access
// restrictions on only this set of endpoints.
type IamServiceRegistrationSupportedNetwork struct {
	// The environment attribute for support.
	EnvironmentAttributes []EnvironmentAttribute `json:"environment_attributes,omitempty"`
}

// UnmarshalIamServiceRegistrationSupportedNetwork unmarshals an instance of IamServiceRegistrationSupportedNetwork from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedNetwork(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedNetwork)
	err = core.UnmarshalModel(m, "environment_attributes", &obj.EnvironmentAttributes, UnmarshalEnvironmentAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_attributes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedNetwork
func (iamServiceRegistrationSupportedNetwork *IamServiceRegistrationSupportedNetwork) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedNetwork.EnvironmentAttributes) {
		var environmentAttributesPatches []map[string]interface{}
		for _, environmentAttributes := range iamServiceRegistrationSupportedNetwork.EnvironmentAttributes {
			environmentAttributesPatches = append(environmentAttributesPatches, environmentAttributes.asPatch())
		}
		_patch["environment_attributes"] = environmentAttributesPatches
	}

	return
}

// IamServiceRegistrationSupportedRole : The supported role.
type IamServiceRegistrationSupportedRole struct {
	// The value belonging to the key.
	ID *string `json:"id,omitempty"`

	// The description for the object.
	Description *IamServiceRegistrationDescriptionObject `json:"description,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`

	// The supported role options.
	Options *SupportedRoleOptions `json:"options,omitempty"`
}

// UnmarshalIamServiceRegistrationSupportedRole unmarshals an instance of IamServiceRegistrationSupportedRole from the specified map of raw messages.
func UnmarshalIamServiceRegistrationSupportedRole(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(IamServiceRegistrationSupportedRole)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "description", &obj.Description, UnmarshalIamServiceRegistrationDescriptionObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "options", &obj.Options, UnmarshalSupportedRoleOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "options-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the IamServiceRegistrationSupportedRole
func (iamServiceRegistrationSupportedRole *IamServiceRegistrationSupportedRole) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(iamServiceRegistrationSupportedRole.ID) {
		_patch["id"] = iamServiceRegistrationSupportedRole.ID
	}
	if !core.IsNil(iamServiceRegistrationSupportedRole.Description) {
		_patch["description"] = iamServiceRegistrationSupportedRole.Description.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedRole.DisplayName) {
		_patch["display_name"] = iamServiceRegistrationSupportedRole.DisplayName.asPatch()
	}
	if !core.IsNil(iamServiceRegistrationSupportedRole.Options) {
		_patch["options"] = iamServiceRegistrationSupportedRole.Options.asPatch()
	}

	return
}

// LearnMoreLinks : The collection of URLs where vendors can learn more about the badge.
type LearnMoreLinks struct {
	// The URL where first-party (IBMer) vendors can learn more about the badge.
	FirstParty *string `json:"first_party,omitempty"`

	// The URL where third-party (non-IBMer) vendors can learn more about the badge.
	ThirdParty *string `json:"third_party,omitempty"`
}

// UnmarshalLearnMoreLinks unmarshals an instance of LearnMoreLinks from the specified map of raw messages.
func UnmarshalLearnMoreLinks(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LearnMoreLinks)
	err = core.UnmarshalPrimitive(m, "first_party", &obj.FirstParty)
	if err != nil {
		err = core.SDKErrorf(err, "", "first_party-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "third_party", &obj.ThirdParty)
	if err != nil {
		err = core.SDKErrorf(err, "", "third_party-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListProductBadgesOptions : The ListProductBadges options.
type ListProductBadgesOptions struct {
	// The maximum number of results returned in the response.
	Limit *int64 `json:"limit,omitempty"`

	// The reference ID of the first item on the page.
	Start *strfmt.UUID `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListProductBadgesOptions : Instantiate ListProductBadgesOptions
func (*PartnerCenterSellV1) NewListProductBadgesOptions() *ListProductBadgesOptions {
	return &ListProductBadgesOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *ListProductBadgesOptions) SetLimit(limit int64) *ListProductBadgesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListProductBadgesOptions) SetStart(start *strfmt.UUID) *ListProductBadgesOptions {
	_options.Start = start
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProductBadgesOptions) SetHeaders(param map[string]string) *ListProductBadgesOptions {
	options.Headers = param
	return options
}

// OnboardingProduct : The object defining the onboarding product in Partner Center - Sell.
type OnboardingProduct struct {
	// The ID of a product in Partner Center - Sell.
	ID *string `json:"id,omitempty"`

	// The IBM Cloud account ID of the provider.
	AccountID *string `json:"account_id,omitempty"`

	// The type of the product.
	Type *string `json:"type,omitempty"`

	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact,omitempty"`

	// The ID of the private catalog that contains the product. Only applicable for software type products.
	PrivateCatalogID *string `json:"private_catalog_id,omitempty"`

	// The ID of the linked private catalog product. Only applicable for software type products.
	PrivateCatalogOfferingID *string `json:"private_catalog_offering_id,omitempty"`

	// The ID of a global catalog object.
	GlobalCatalogOfferingID *string `json:"global_catalog_offering_id,omitempty"`

	// The ID of a global catalog object.
	StagingGlobalCatalogOfferingID *string `json:"staging_global_catalog_offering_id,omitempty"`

	// The ID of the approval workflow of your product.
	ApproverResourceID *string `json:"approver_resource_id,omitempty"`

	// IAM registration identifier.
	IamRegistrationID *string `json:"iam_registration_id,omitempty"`

	// The Export Control Classification Number of your product.
	EccnNumber *string `json:"eccn_number,omitempty"`

	// The ERO class of your product.
	EroClass *string `json:"ero_class,omitempty"`

	// The United Nations Standard Products and Services Code of your product.
	Unspsc *float64 `json:"unspsc,omitempty"`

	// The tax assessment type of your product.
	TaxAssessment *string `json:"tax_assessment,omitempty"`

	// The support information that is not displayed in the catalog, but available in ServiceNow.
	Support *OnboardingProductSupport `json:"support,omitempty"`
}

// Constants associated with the OnboardingProduct.Type property.
// The type of the product.
const (
	OnboardingProduct_Type_ProfessionalService = "professional_service"
	OnboardingProduct_Type_Service             = "service"
	OnboardingProduct_Type_Software            = "software"
)

// UnmarshalOnboardingProduct unmarshals an instance of OnboardingProduct from the specified map of raw messages.
func UnmarshalOnboardingProduct(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OnboardingProduct)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "primary_contact", &obj.PrimaryContact, UnmarshalPrimaryContact)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_catalog_id", &obj.PrivateCatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_catalog_offering_id", &obj.PrivateCatalogOfferingID)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_catalog_offering_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "global_catalog_offering_id", &obj.GlobalCatalogOfferingID)
	if err != nil {
		err = core.SDKErrorf(err, "", "global_catalog_offering_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "staging_global_catalog_offering_id", &obj.StagingGlobalCatalogOfferingID)
	if err != nil {
		err = core.SDKErrorf(err, "", "staging_global_catalog_offering_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approver_resource_id", &obj.ApproverResourceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "approver_resource_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "iam_registration_id", &obj.IamRegistrationID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_registration_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "eccn_number", &obj.EccnNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "eccn_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ero_class", &obj.EroClass)
	if err != nil {
		err = core.SDKErrorf(err, "", "ero_class-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unspsc", &obj.Unspsc)
	if err != nil {
		err = core.SDKErrorf(err, "", "unspsc-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tax_assessment", &obj.TaxAssessment)
	if err != nil {
		err = core.SDKErrorf(err, "", "tax_assessment-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "support", &obj.Support, UnmarshalOnboardingProductSupport)
	if err != nil {
		err = core.SDKErrorf(err, "", "support-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OnboardingProductPatch : The request body for updating a product in Partner Center - Sell.
type OnboardingProductPatch struct {
	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact,omitempty"`

	// The Export Control Classification Number of your product.
	EccnNumber *string `json:"eccn_number,omitempty"`

	// The ERO class of your product.
	EroClass *string `json:"ero_class,omitempty"`

	// The United Nations Standard Products and Services Code of your product.
	Unspsc *float64 `json:"unspsc,omitempty"`

	// The tax assessment type of your product.
	TaxAssessment *string `json:"tax_assessment,omitempty"`

	// The support information that is not displayed in the catalog, but available in ServiceNow.
	Support *OnboardingProductSupport `json:"support,omitempty"`
}

// UnmarshalOnboardingProductPatch unmarshals an instance of OnboardingProductPatch from the specified map of raw messages.
func UnmarshalOnboardingProductPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OnboardingProductPatch)
	err = core.UnmarshalModel(m, "primary_contact", &obj.PrimaryContact, UnmarshalPrimaryContact)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "eccn_number", &obj.EccnNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "eccn_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ero_class", &obj.EroClass)
	if err != nil {
		err = core.SDKErrorf(err, "", "ero_class-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unspsc", &obj.Unspsc)
	if err != nil {
		err = core.SDKErrorf(err, "", "unspsc-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tax_assessment", &obj.TaxAssessment)
	if err != nil {
		err = core.SDKErrorf(err, "", "tax_assessment-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "support", &obj.Support, UnmarshalOnboardingProductSupport)
	if err != nil {
		err = core.SDKErrorf(err, "", "support-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the OnboardingProductPatch
func (onboardingProductPatch *OnboardingProductPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(onboardingProductPatch.PrimaryContact) {
		_patch["primary_contact"] = onboardingProductPatch.PrimaryContact.asPatch()
	}
	if !core.IsNil(onboardingProductPatch.EccnNumber) {
		_patch["eccn_number"] = onboardingProductPatch.EccnNumber
	}
	if !core.IsNil(onboardingProductPatch.EroClass) {
		_patch["ero_class"] = onboardingProductPatch.EroClass
	}
	if !core.IsNil(onboardingProductPatch.Unspsc) {
		_patch["unspsc"] = onboardingProductPatch.Unspsc
	}
	if !core.IsNil(onboardingProductPatch.TaxAssessment) {
		_patch["tax_assessment"] = onboardingProductPatch.TaxAssessment
	}
	if !core.IsNil(onboardingProductPatch.Support) {
		_patch["support"] = onboardingProductPatch.Support.asPatch()
	}

	return
}

// OnboardingProductSupport : The support information that is not displayed in the catalog, but available in ServiceNow.
type OnboardingProductSupport struct {
	// The list of contacts in case of support escalations.
	EscalationContacts []OnboardingProductSupportEscalationContactItems `json:"escalation_contacts,omitempty"`
}

// UnmarshalOnboardingProductSupport unmarshals an instance of OnboardingProductSupport from the specified map of raw messages.
func UnmarshalOnboardingProductSupport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OnboardingProductSupport)
	err = core.UnmarshalModel(m, "escalation_contacts", &obj.EscalationContacts, UnmarshalOnboardingProductSupportEscalationContactItems)
	if err != nil {
		err = core.SDKErrorf(err, "", "escalation_contacts-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the OnboardingProductSupport
func (onboardingProductSupport *OnboardingProductSupport) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(onboardingProductSupport.EscalationContacts) {
		var escalationContactsPatches []map[string]interface{}
		for _, escalationContacts := range onboardingProductSupport.EscalationContacts {
			escalationContactsPatches = append(escalationContactsPatches, escalationContacts.asPatch())
		}
		_patch["escalation_contacts"] = escalationContactsPatches
	}

	return
}

// OnboardingProductSupportEscalationContactItems : The details of a support escalation contact.
type OnboardingProductSupportEscalationContactItems struct {
	// The name of the support escalation contact.
	Name *string `json:"name,omitempty"`

	// The email address of the support escalation contact.
	Email *string `json:"email,omitempty"`

	// The role of the support escalation contact.
	Role *string `json:"role,omitempty"`
}

// UnmarshalOnboardingProductSupportEscalationContactItems unmarshals an instance of OnboardingProductSupportEscalationContactItems from the specified map of raw messages.
func UnmarshalOnboardingProductSupportEscalationContactItems(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OnboardingProductSupportEscalationContactItems)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "role", &obj.Role)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the OnboardingProductSupportEscalationContactItems
func (onboardingProductSupportEscalationContactItems *OnboardingProductSupportEscalationContactItems) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(onboardingProductSupportEscalationContactItems.Name) {
		_patch["name"] = onboardingProductSupportEscalationContactItems.Name
	}
	if !core.IsNil(onboardingProductSupportEscalationContactItems.Email) {
		_patch["email"] = onboardingProductSupportEscalationContactItems.Email
	}
	if !core.IsNil(onboardingProductSupportEscalationContactItems.Role) {
		_patch["role"] = onboardingProductSupportEscalationContactItems.Role
	}

	return
}

// PrimaryContact : The primary contact for your product.
type PrimaryContact struct {
	// The name of the primary contact for your product.
	Name *string `json:"name" validate:"required"`

	// The email address of the primary contact for your product.
	Email *string `json:"email" validate:"required"`
}

// NewPrimaryContact : Instantiate PrimaryContact (Generic Model Constructor)
func (*PartnerCenterSellV1) NewPrimaryContact(name string, email string) (_model *PrimaryContact, err error) {
	_model = &PrimaryContact{
		Name:  core.StringPtr(name),
		Email: core.StringPtr(email),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPrimaryContact unmarshals an instance of PrimaryContact from the specified map of raw messages.
func UnmarshalPrimaryContact(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PrimaryContact)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the PrimaryContact
func (primaryContact *PrimaryContact) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(primaryContact.Name) {
		_patch["name"] = primaryContact.Name
	}
	if !core.IsNil(primaryContact.Email) {
		_patch["email"] = primaryContact.Email
	}

	return
}

// ProductBadge : The details of the cloud badge.
type ProductBadge struct {
	// The ID of the badge.
	ID *string `json:"id" validate:"required"`

	// The name of the badge.
	Label *string `json:"label,omitempty"`

	// The description of the badge.
	Description *string `json:"description,omitempty"`

	// The internal description of the badge.
	InternalDescription *string `json:"internal_description,omitempty"`

	// The collection of URLs where vendors can learn more about the badge.
	LearnMoreLinks *LearnMoreLinks `json:"learn_more_links,omitempty"`

	// The URL to get started with the validation against this certification.
	GetStartedLink *string `json:"get_started_link,omitempty"`

	// Deprecated, will be removed.
	Tag *string `json:"tag,omitempty"`
}

// UnmarshalProductBadge unmarshals an instance of ProductBadge from the specified map of raw messages.
func UnmarshalProductBadge(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProductBadge)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "internal_description", &obj.InternalDescription)
	if err != nil {
		err = core.SDKErrorf(err, "", "internal_description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "learn_more_links", &obj.LearnMoreLinks, UnmarshalLearnMoreLinks)
	if err != nil {
		err = core.SDKErrorf(err, "", "learn_more_links-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "get_started_link", &obj.GetStartedLink)
	if err != nil {
		err = core.SDKErrorf(err, "", "get_started_link-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tag", &obj.Tag)
	if err != nil {
		err = core.SDKErrorf(err, "", "tag-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProductBadgeCollection : The list of all the available cloud badges.
type ProductBadgeCollection struct {
	// The maximum number of results returned in this response.
	Limit *int64 `json:"limit" validate:"required"`

	// The maximum number of results returned in this response.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The page reference information.
	First *Bookmark `json:"first,omitempty"`

	// The page reference information.
	Next *Bookmark `json:"next,omitempty"`

	// The page reference information.
	Previous *Bookmark `json:"previous,omitempty"`

	// The page reference information.
	Last *Bookmark `json:"last,omitempty"`

	// The list of all the available cloud badges.
	ProductBadges []ProductBadge `json:"product_badges" validate:"required"`
}

// UnmarshalProductBadgeCollection unmarshals an instance of ProductBadgeCollection from the specified map of raw messages.
func UnmarshalProductBadgeCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProductBadgeCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalBookmark)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalBookmark)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalBookmark)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalBookmark)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "product_badges", &obj.ProductBadges, UnmarshalProductBadge)
	if err != nil {
		err = core.SDKErrorf(err, "", "product_badges-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Registration : The object defining the Partner Center - Sell registration.
type Registration struct {
	// The ID of your registration, which is the same as your billing and metering (BSS) account ID.
	ID *string `json:"id,omitempty"`

	// The ID of your account.
	AccountID *string `json:"account_id" validate:"required"`

	// The name of your company that is displayed in the IBM Cloud catalog.
	CompanyName *string `json:"company_name" validate:"required"`

	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact" validate:"required"`

	// The default private catalog in which products are created.
	DefaultPrivateCatalogID *string `json:"default_private_catalog_id,omitempty"`

	// The onboarding access group for your team.
	ProviderAccessGroup *string `json:"provider_access_group,omitempty"`

	// The time when the registration was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The time when the registration was updated.
	UpdatedAt *string `json:"updated_at,omitempty"`
}

// UnmarshalRegistration unmarshals an instance of Registration from the specified map of raw messages.
func UnmarshalRegistration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Registration)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "company_name", &obj.CompanyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "company_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "primary_contact", &obj.PrimaryContact, UnmarshalPrimaryContact)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "default_private_catalog_id", &obj.DefaultPrivateCatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "default_private_catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provider_access_group", &obj.ProviderAccessGroup)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider_access_group-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RegistrationPatch : The request body for updating a registration in Partner Center - Sell.
type RegistrationPatch struct {
	// The name of your company that is displayed in the IBM Cloud catalog.
	CompanyName *string `json:"company_name,omitempty"`

	// The primary contact for your product.
	PrimaryContact *PrimaryContact `json:"primary_contact,omitempty"`

	// The default private catalog in which products are created.
	DefaultPrivateCatalogID *string `json:"default_private_catalog_id,omitempty"`

	// The onboarding access group for your team.
	ProviderAccessGroup *string `json:"provider_access_group,omitempty"`
}

// UnmarshalRegistrationPatch unmarshals an instance of RegistrationPatch from the specified map of raw messages.
func UnmarshalRegistrationPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RegistrationPatch)
	err = core.UnmarshalPrimitive(m, "company_name", &obj.CompanyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "company_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "primary_contact", &obj.PrimaryContact, UnmarshalPrimaryContact)
	if err != nil {
		err = core.SDKErrorf(err, "", "primary_contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "default_private_catalog_id", &obj.DefaultPrivateCatalogID)
	if err != nil {
		err = core.SDKErrorf(err, "", "default_private_catalog_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provider_access_group", &obj.ProviderAccessGroup)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider_access_group-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the RegistrationPatch
func (registrationPatch *RegistrationPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(registrationPatch.CompanyName) {
		_patch["company_name"] = registrationPatch.CompanyName
	}
	if !core.IsNil(registrationPatch.PrimaryContact) {
		_patch["primary_contact"] = registrationPatch.PrimaryContact.asPatch()
	}
	if !core.IsNil(registrationPatch.DefaultPrivateCatalogID) {
		_patch["default_private_catalog_id"] = registrationPatch.DefaultPrivateCatalogID
	}
	if !core.IsNil(registrationPatch.ProviderAccessGroup) {
		_patch["provider_access_group"] = registrationPatch.ProviderAccessGroup
	}

	return
}

// SupportAuthorizationSubjectAttribute : The list of supported authorization subject properties.
type SupportAuthorizationSubjectAttribute struct {
	// The name of the service.
	ServiceName *string `json:"service_name,omitempty"`

	// The type of the service.
	ResourceType *string `json:"resource_type,omitempty"`
}

// UnmarshalSupportAuthorizationSubjectAttribute unmarshals an instance of SupportAuthorizationSubjectAttribute from the specified map of raw messages.
func UnmarshalSupportAuthorizationSubjectAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportAuthorizationSubjectAttribute)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportAuthorizationSubjectAttribute
func (supportAuthorizationSubjectAttribute *SupportAuthorizationSubjectAttribute) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportAuthorizationSubjectAttribute.ServiceName) {
		_patch["service_name"] = supportAuthorizationSubjectAttribute.ServiceName
	}
	if !core.IsNil(supportAuthorizationSubjectAttribute.ResourceType) {
		_patch["resource_type"] = supportAuthorizationSubjectAttribute.ResourceType
	}

	return
}

// SupportDetailsItem : SupportDetailsItem struct
type SupportDetailsItem struct {
	// The type of support for this support channel.
	Type *string `json:"type,omitempty"`

	// The contact information for this support channel.
	Contact *string `json:"contact,omitempty"`

	// The time interval of providing support in units and values.
	ResponseWaitTime *SupportTimeInterval `json:"response_wait_time,omitempty"`

	// The time period during which support is available for the service.
	Availability *SupportDetailsItemAvailability `json:"availability,omitempty"`
}

// Constants associated with the SupportDetailsItem.Type property.
// The type of support for this support channel.
const (
	SupportDetailsItem_Type_Chat        = "chat"
	SupportDetailsItem_Type_Email       = "email"
	SupportDetailsItem_Type_Other       = "other"
	SupportDetailsItem_Type_Phone       = "phone"
	SupportDetailsItem_Type_Slack       = "slack"
	SupportDetailsItem_Type_SupportSite = "support_site"
)

// UnmarshalSupportDetailsItem unmarshals an instance of SupportDetailsItem from the specified map of raw messages.
func UnmarshalSupportDetailsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportDetailsItem)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "contact", &obj.Contact)
	if err != nil {
		err = core.SDKErrorf(err, "", "contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "response_wait_time", &obj.ResponseWaitTime, UnmarshalSupportTimeInterval)
	if err != nil {
		err = core.SDKErrorf(err, "", "response_wait_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "availability", &obj.Availability, UnmarshalSupportDetailsItemAvailability)
	if err != nil {
		err = core.SDKErrorf(err, "", "availability-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportDetailsItem
func (supportDetailsItem *SupportDetailsItem) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportDetailsItem.Type) {
		_patch["type"] = supportDetailsItem.Type
	}
	if !core.IsNil(supportDetailsItem.Contact) {
		_patch["contact"] = supportDetailsItem.Contact
	}
	if !core.IsNil(supportDetailsItem.ResponseWaitTime) {
		_patch["response_wait_time"] = supportDetailsItem.ResponseWaitTime.asPatch()
	}
	if !core.IsNil(supportDetailsItem.Availability) {
		_patch["availability"] = supportDetailsItem.Availability.asPatch()
	}

	return
}

// SupportDetailsItemAvailability : The time period during which support is available for the service.
type SupportDetailsItemAvailability struct {
	// The support hours available for the service.
	Times []SupportDetailsItemAvailabilityTime `json:"times,omitempty"`

	// The timezones in which support is available. Only relevant if `always_available` is set to false.
	Timezone *string `json:"timezone,omitempty"`

	// Whether the support for the service is always available.
	AlwaysAvailable *bool `json:"always_available,omitempty"`
}

// UnmarshalSupportDetailsItemAvailability unmarshals an instance of SupportDetailsItemAvailability from the specified map of raw messages.
func UnmarshalSupportDetailsItemAvailability(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportDetailsItemAvailability)
	err = core.UnmarshalModel(m, "times", &obj.Times, UnmarshalSupportDetailsItemAvailabilityTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "times-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		err = core.SDKErrorf(err, "", "timezone-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "always_available", &obj.AlwaysAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "always_available-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportDetailsItemAvailability
func (supportDetailsItemAvailability *SupportDetailsItemAvailability) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportDetailsItemAvailability.Times) {
		var timesPatches []map[string]interface{}
		for _, times := range supportDetailsItemAvailability.Times {
			timesPatches = append(timesPatches, times.asPatch())
		}
		_patch["times"] = timesPatches
	}
	if !core.IsNil(supportDetailsItemAvailability.Timezone) {
		_patch["timezone"] = supportDetailsItemAvailability.Timezone
	}
	if !core.IsNil(supportDetailsItemAvailability.AlwaysAvailable) {
		_patch["always_available"] = supportDetailsItemAvailability.AlwaysAvailable
	}

	return
}

// SupportDetailsItemAvailabilityTime : SupportDetailsItemAvailabilityTime struct
type SupportDetailsItemAvailabilityTime struct {
	// The number of days in a week when support is available for the service.
	Day *float64 `json:"day,omitempty"`

	// The time in the day when support starts for the service.
	StartTime *string `json:"start_time,omitempty"`

	// The time in the day when support ends for the service.
	EndTime *string `json:"end_time,omitempty"`
}

// UnmarshalSupportDetailsItemAvailabilityTime unmarshals an instance of SupportDetailsItemAvailabilityTime from the specified map of raw messages.
func UnmarshalSupportDetailsItemAvailabilityTime(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportDetailsItemAvailabilityTime)
	err = core.UnmarshalPrimitive(m, "day", &obj.Day)
	if err != nil {
		err = core.SDKErrorf(err, "", "day-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "end_time", &obj.EndTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "end_time-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportDetailsItemAvailabilityTime
func (supportDetailsItemAvailabilityTime *SupportDetailsItemAvailabilityTime) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportDetailsItemAvailabilityTime.Day) {
		_patch["day"] = supportDetailsItemAvailabilityTime.Day
	}
	if !core.IsNil(supportDetailsItemAvailabilityTime.StartTime) {
		_patch["start_time"] = supportDetailsItemAvailabilityTime.StartTime
	}
	if !core.IsNil(supportDetailsItemAvailabilityTime.EndTime) {
		_patch["end_time"] = supportDetailsItemAvailabilityTime.EndTime
	}

	return
}

// SupportEscalation : The details of the support escalation process.
type SupportEscalation struct {
	// The support contact information of the escalation team.
	Contact *string `json:"contact,omitempty"`

	// The time interval of providing support in units and values.
	EscalationWaitTime *SupportTimeInterval `json:"escalation_wait_time,omitempty"`

	// The time interval of providing support in units and values.
	ResponseWaitTime *SupportTimeInterval `json:"response_wait_time,omitempty"`
}

// UnmarshalSupportEscalation unmarshals an instance of SupportEscalation from the specified map of raw messages.
func UnmarshalSupportEscalation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportEscalation)
	err = core.UnmarshalPrimitive(m, "contact", &obj.Contact)
	if err != nil {
		err = core.SDKErrorf(err, "", "contact-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "escalation_wait_time", &obj.EscalationWaitTime, UnmarshalSupportTimeInterval)
	if err != nil {
		err = core.SDKErrorf(err, "", "escalation_wait_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "response_wait_time", &obj.ResponseWaitTime, UnmarshalSupportTimeInterval)
	if err != nil {
		err = core.SDKErrorf(err, "", "response_wait_time-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportEscalation
func (supportEscalation *SupportEscalation) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportEscalation.Contact) {
		_patch["contact"] = supportEscalation.Contact
	}
	if !core.IsNil(supportEscalation.EscalationWaitTime) {
		_patch["escalation_wait_time"] = supportEscalation.EscalationWaitTime.asPatch()
	}
	if !core.IsNil(supportEscalation.ResponseWaitTime) {
		_patch["response_wait_time"] = supportEscalation.ResponseWaitTime.asPatch()
	}

	return
}

// SupportTimeInterval : The time interval of providing support in units and values.
type SupportTimeInterval struct {
	// The number of time units.
	Value *float64 `json:"value,omitempty"`

	// The unit of the time.
	Type *string `json:"type,omitempty"`
}

// UnmarshalSupportTimeInterval unmarshals an instance of SupportTimeInterval from the specified map of raw messages.
func UnmarshalSupportTimeInterval(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportTimeInterval)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportTimeInterval
func (supportTimeInterval *SupportTimeInterval) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportTimeInterval.Value) {
		_patch["value"] = supportTimeInterval.Value
	}
	if !core.IsNil(supportTimeInterval.Type) {
		_patch["type"] = supportTimeInterval.Type
	}

	return
}

// SupportedAttributeUi : The user interface.
type SupportedAttributeUi struct {
	// The type of the input.
	InputType *string `json:"input_type,omitempty"`

	// The details of the input.
	InputDetails *SupportedAttributeUiInputDetails `json:"input_details,omitempty"`
}

// UnmarshalSupportedAttributeUi unmarshals an instance of SupportedAttributeUi from the specified map of raw messages.
func UnmarshalSupportedAttributeUi(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributeUi)
	err = core.UnmarshalPrimitive(m, "input_type", &obj.InputType)
	if err != nil {
		err = core.SDKErrorf(err, "", "input_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "input_details", &obj.InputDetails, UnmarshalSupportedAttributeUiInputDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "input_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributeUi
func (supportedAttributeUi *SupportedAttributeUi) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributeUi.InputType) {
		_patch["input_type"] = supportedAttributeUi.InputType
	}
	if !core.IsNil(supportedAttributeUi.InputDetails) {
		_patch["input_details"] = supportedAttributeUi.InputDetails.asPatch()
	}

	return
}

// SupportedAttributeUiInputDetails : The details of the input.
type SupportedAttributeUiInputDetails struct {
	// They type of the input details.
	Type *string `json:"type,omitempty"`

	// The provided values of input details.
	Values []SupportedAttributeUiInputValue `json:"values,omitempty"`

	// Required if type is gst.
	Gst *SupportedAttributeUiInputGst `json:"gst,omitempty"`

	// The URL data for user interface.
	URL *SupportedAttributeUiInputURL `json:"url,omitempty"`
}

// UnmarshalSupportedAttributeUiInputDetails unmarshals an instance of SupportedAttributeUiInputDetails from the specified map of raw messages.
func UnmarshalSupportedAttributeUiInputDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributeUiInputDetails)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "values", &obj.Values, UnmarshalSupportedAttributeUiInputValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "values-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "gst", &obj.Gst, UnmarshalSupportedAttributeUiInputGst)
	if err != nil {
		err = core.SDKErrorf(err, "", "gst-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "url", &obj.URL, UnmarshalSupportedAttributeUiInputURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributeUiInputDetails
func (supportedAttributeUiInputDetails *SupportedAttributeUiInputDetails) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributeUiInputDetails.Type) {
		_patch["type"] = supportedAttributeUiInputDetails.Type
	}
	if !core.IsNil(supportedAttributeUiInputDetails.Values) {
		var valuesPatches []map[string]interface{}
		for _, values := range supportedAttributeUiInputDetails.Values {
			valuesPatches = append(valuesPatches, values.asPatch())
		}
		_patch["values"] = valuesPatches
	}
	if !core.IsNil(supportedAttributeUiInputDetails.Gst) {
		_patch["gst"] = supportedAttributeUiInputDetails.Gst.asPatch()
	}
	if !core.IsNil(supportedAttributeUiInputDetails.URL) {
		_patch["url"] = supportedAttributeUiInputDetails.URL.asPatch()
	}

	return
}

// SupportedAttributeUiInputGst : Required if type is gst.
type SupportedAttributeUiInputGst struct {
	// The query to use.
	Query *string `json:"query,omitempty"`

	// The value of the property name.
	ValuePropertyName *string `json:"value_property_name,omitempty"`

	// One of labelPropertyName or inputOptionLabel is required.
	LabelPropertyName *string `json:"label_property_name,omitempty"`

	// The label for option input.
	InputOptionLabel *string `json:"input_option_label,omitempty"`
}

// UnmarshalSupportedAttributeUiInputGst unmarshals an instance of SupportedAttributeUiInputGst from the specified map of raw messages.
func UnmarshalSupportedAttributeUiInputGst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributeUiInputGst)
	err = core.UnmarshalPrimitive(m, "query", &obj.Query)
	if err != nil {
		err = core.SDKErrorf(err, "", "query-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value_property_name", &obj.ValuePropertyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "value_property_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label_property_name", &obj.LabelPropertyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "label_property_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "input_option_label", &obj.InputOptionLabel)
	if err != nil {
		err = core.SDKErrorf(err, "", "input_option_label-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributeUiInputGst
func (supportedAttributeUiInputGst *SupportedAttributeUiInputGst) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributeUiInputGst.Query) {
		_patch["query"] = supportedAttributeUiInputGst.Query
	}
	if !core.IsNil(supportedAttributeUiInputGst.ValuePropertyName) {
		_patch["value_property_name"] = supportedAttributeUiInputGst.ValuePropertyName
	}
	if !core.IsNil(supportedAttributeUiInputGst.LabelPropertyName) {
		_patch["label_property_name"] = supportedAttributeUiInputGst.LabelPropertyName
	}
	if !core.IsNil(supportedAttributeUiInputGst.InputOptionLabel) {
		_patch["input_option_label"] = supportedAttributeUiInputGst.InputOptionLabel
	}

	return
}

// SupportedAttributeUiInputURL : The URL data for user interface.
type SupportedAttributeUiInputURL struct {
	// The URL of the user interface interface.
	UrlEndpoint *string `json:"url_endpoint,omitempty"`

	// The label options for the user interface URL.
	InputOptionLabel *string `json:"input_option_label,omitempty"`
}

// UnmarshalSupportedAttributeUiInputURL unmarshals an instance of SupportedAttributeUiInputURL from the specified map of raw messages.
func UnmarshalSupportedAttributeUiInputURL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributeUiInputURL)
	err = core.UnmarshalPrimitive(m, "url_endpoint", &obj.UrlEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "url_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "input_option_label", &obj.InputOptionLabel)
	if err != nil {
		err = core.SDKErrorf(err, "", "input_option_label-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributeUiInputURL
func (supportedAttributeUiInputURL *SupportedAttributeUiInputURL) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributeUiInputURL.UrlEndpoint) {
		_patch["url_endpoint"] = supportedAttributeUiInputURL.UrlEndpoint
	}
	if !core.IsNil(supportedAttributeUiInputURL.InputOptionLabel) {
		_patch["input_option_label"] = supportedAttributeUiInputURL.InputOptionLabel
	}

	return
}

// SupportedAttributeUiInputValue : SupportedAttributeUiInputValue struct
type SupportedAttributeUiInputValue struct {
	// The values of input details.
	Value *string `json:"value,omitempty"`

	// The display name of the object.
	DisplayName *IamServiceRegistrationDisplayNameObject `json:"display_name,omitempty"`
}

// UnmarshalSupportedAttributeUiInputValue unmarshals an instance of SupportedAttributeUiInputValue from the specified map of raw messages.
func UnmarshalSupportedAttributeUiInputValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributeUiInputValue)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "display_name", &obj.DisplayName, UnmarshalIamServiceRegistrationDisplayNameObject)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributeUiInputValue
func (supportedAttributeUiInputValue *SupportedAttributeUiInputValue) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributeUiInputValue.Value) {
		_patch["value"] = supportedAttributeUiInputValue.Value
	}
	if !core.IsNil(supportedAttributeUiInputValue.DisplayName) {
		_patch["display_name"] = supportedAttributeUiInputValue.DisplayName.asPatch()
	}

	return
}

// SupportedAttributesOptions : The list of support attribute options.
type SupportedAttributesOptions struct {
	// The supported attribute operator.
	Operators []string `json:"operators,omitempty"`

	// Optional opt-in if attribute is hidden from customers (customer can still use it if they found out themselves).
	Hidden *bool `json:"hidden,omitempty"`

	// The list of supported patterns.
	SupportedAttributes []string `json:"supported_attributes,omitempty"`

	// The list of policy types.
	PolicyTypes []string `json:"policy_types,omitempty"`

	// Indicate whether the empty value is supported.
	IsEmptyValueSupported *bool `json:"is_empty_value_supported,omitempty"`

	// Indicate whether the false value is supported for stringExists operator.
	IsStringExistsFalseValueSupported *bool `json:"is_string_exists_false_value_supported,omitempty"`

	// The name of attribute.
	Key *string `json:"key,omitempty"`

	// Resource hierarchy options for composite services.
	ResourceHierarchy *SupportedAttributesOptionsResourceHierarchy `json:"resource_hierarchy,omitempty"`
}

// Constants associated with the SupportedAttributesOptions.Operators property.
// The list of multiple option values.
const (
	SupportedAttributesOptions_Operators_Stringequals      = "stringEquals"
	SupportedAttributesOptions_Operators_Stringequalsanyof = "stringEqualsAnyOf"
	SupportedAttributesOptions_Operators_Stringmatch       = "stringMatch"
	SupportedAttributesOptions_Operators_Stringmatchanyof  = "stringMatchAnyOf"
)

// Constants associated with the SupportedAttributesOptions.PolicyTypes property.
const (
	SupportedAttributesOptions_PolicyTypes_Access        = "access"
	SupportedAttributesOptions_PolicyTypes_Authorization = "authorization"
)

// UnmarshalSupportedAttributesOptions unmarshals an instance of SupportedAttributesOptions from the specified map of raw messages.
func UnmarshalSupportedAttributesOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributesOptions)
	err = core.UnmarshalPrimitive(m, "operators", &obj.Operators)
	if err != nil {
		err = core.SDKErrorf(err, "", "operators-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hidden", &obj.Hidden)
	if err != nil {
		err = core.SDKErrorf(err, "", "hidden-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "supported_attributes", &obj.SupportedAttributes)
	if err != nil {
		err = core.SDKErrorf(err, "", "supported_attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_types", &obj.PolicyTypes)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_types-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_empty_value_supported", &obj.IsEmptyValueSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_empty_value_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_string_exists_false_value_supported", &obj.IsStringExistsFalseValueSupported)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_string_exists_false_value_supported-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_hierarchy", &obj.ResourceHierarchy, UnmarshalSupportedAttributesOptionsResourceHierarchy)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_hierarchy-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributesOptions
func (supportedAttributesOptions *SupportedAttributesOptions) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributesOptions.Operators) {
		_patch["operators"] = supportedAttributesOptions.Operators
	}
	if !core.IsNil(supportedAttributesOptions.Hidden) {
		_patch["hidden"] = supportedAttributesOptions.Hidden
	}
	if !core.IsNil(supportedAttributesOptions.SupportedAttributes) {
		_patch["supported_attributes"] = supportedAttributesOptions.SupportedAttributes
	}
	if !core.IsNil(supportedAttributesOptions.PolicyTypes) {
		_patch["policy_types"] = supportedAttributesOptions.PolicyTypes
	}
	if !core.IsNil(supportedAttributesOptions.IsEmptyValueSupported) {
		_patch["is_empty_value_supported"] = supportedAttributesOptions.IsEmptyValueSupported
	}
	if !core.IsNil(supportedAttributesOptions.IsStringExistsFalseValueSupported) {
		_patch["is_string_exists_false_value_supported"] = supportedAttributesOptions.IsStringExistsFalseValueSupported
	}
	if !core.IsNil(supportedAttributesOptions.Key) {
		_patch["key"] = supportedAttributesOptions.Key
	}
	if !core.IsNil(supportedAttributesOptions.ResourceHierarchy) {
		_patch["resource_hierarchy"] = supportedAttributesOptions.ResourceHierarchy.asPatch()
	}

	return
}

// SupportedAttributesOptionsResourceHierarchy : Resource hierarchy options for composite services.
type SupportedAttributesOptionsResourceHierarchy struct {
	// Hierarchy description key.
	Key *SupportedAttributesOptionsResourceHierarchyKey `json:"key,omitempty"`

	// Hierarchy description value.
	Value *SupportedAttributesOptionsResourceHierarchyValue `json:"value,omitempty"`
}

// UnmarshalSupportedAttributesOptionsResourceHierarchy unmarshals an instance of SupportedAttributesOptionsResourceHierarchy from the specified map of raw messages.
func UnmarshalSupportedAttributesOptionsResourceHierarchy(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributesOptionsResourceHierarchy)
	err = core.UnmarshalModel(m, "key", &obj.Key, UnmarshalSupportedAttributesOptionsResourceHierarchyKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "value", &obj.Value, UnmarshalSupportedAttributesOptionsResourceHierarchyValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributesOptionsResourceHierarchy
func (supportedAttributesOptionsResourceHierarchy *SupportedAttributesOptionsResourceHierarchy) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributesOptionsResourceHierarchy.Key) {
		_patch["key"] = supportedAttributesOptionsResourceHierarchy.Key.asPatch()
	}
	if !core.IsNil(supportedAttributesOptionsResourceHierarchy.Value) {
		_patch["value"] = supportedAttributesOptionsResourceHierarchy.Value.asPatch()
	}

	return
}

// SupportedAttributesOptionsResourceHierarchyKey : Hierarchy description key.
type SupportedAttributesOptionsResourceHierarchyKey struct {
	// Key.
	Key *string `json:"key,omitempty"`

	// Value.
	Value *string `json:"value,omitempty"`
}

// UnmarshalSupportedAttributesOptionsResourceHierarchyKey unmarshals an instance of SupportedAttributesOptionsResourceHierarchyKey from the specified map of raw messages.
func UnmarshalSupportedAttributesOptionsResourceHierarchyKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributesOptionsResourceHierarchyKey)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributesOptionsResourceHierarchyKey
func (supportedAttributesOptionsResourceHierarchyKey *SupportedAttributesOptionsResourceHierarchyKey) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributesOptionsResourceHierarchyKey.Key) {
		_patch["key"] = supportedAttributesOptionsResourceHierarchyKey.Key
	}
	if !core.IsNil(supportedAttributesOptionsResourceHierarchyKey.Value) {
		_patch["value"] = supportedAttributesOptionsResourceHierarchyKey.Value
	}

	return
}

// SupportedAttributesOptionsResourceHierarchyValue : Hierarchy description value.
type SupportedAttributesOptionsResourceHierarchyValue struct {
	// Key.
	Key *string `json:"key,omitempty"`
}

// UnmarshalSupportedAttributesOptionsResourceHierarchyValue unmarshals an instance of SupportedAttributesOptionsResourceHierarchyValue from the specified map of raw messages.
func UnmarshalSupportedAttributesOptionsResourceHierarchyValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedAttributesOptionsResourceHierarchyValue)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedAttributesOptionsResourceHierarchyValue
func (supportedAttributesOptionsResourceHierarchyValue *SupportedAttributesOptionsResourceHierarchyValue) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedAttributesOptionsResourceHierarchyValue.Key) {
		_patch["key"] = supportedAttributesOptionsResourceHierarchyValue.Key
	}

	return
}

// SupportedRoleOptions : The supported role options.
type SupportedRoleOptions struct {
	// Optional opt-in to require access control on the role.
	AccessPolicy *bool `json:"access_policy" validate:"required"`

	// Additional properties for access policy.
	AdditionalPropertiesForAccessPolicy map[string]string `json:"additional_properties_for_access_policy,omitempty"`

	// Optional opt-in to require checking policy type when applying the role.
	PolicyType []string `json:"policy_type,omitempty"`

	// Optional opt-in to require checking account type when applying the role.
	AccountType *string `json:"account_type,omitempty"`
}

// Constants associated with the SupportedRoleOptions.PolicyType property.
// Policy type.
const (
	SupportedRoleOptions_PolicyType_Access                 = "access"
	SupportedRoleOptions_PolicyType_Authorization          = "authorization"
	SupportedRoleOptions_PolicyType_AuthorizationDelegated = "authorization-delegated"
)

// Constants associated with the SupportedRoleOptions.AccountType property.
// Optional opt-in to require checking account type when applying the role.
const (
	SupportedRoleOptions_AccountType_Enterprise = "enterprise"
)

// NewSupportedRoleOptions : Instantiate SupportedRoleOptions (Generic Model Constructor)
func (*PartnerCenterSellV1) NewSupportedRoleOptions(accessPolicy bool) (_model *SupportedRoleOptions, err error) {
	_model = &SupportedRoleOptions{
		AccessPolicy: core.BoolPtr(accessPolicy),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalSupportedRoleOptions unmarshals an instance of SupportedRoleOptions from the specified map of raw messages.
func UnmarshalSupportedRoleOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SupportedRoleOptions)
	err = core.UnmarshalPrimitive(m, "access_policy", &obj.AccessPolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "access_policy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "additional_properties_for_access_policy", &obj.AdditionalPropertiesForAccessPolicy)
	if err != nil {
		err = core.SDKErrorf(err, "", "additional_properties_for_access_policy-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "policy_type", &obj.PolicyType)
	if err != nil {
		err = core.SDKErrorf(err, "", "policy_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_type", &obj.AccountType)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the SupportedRoleOptions
func (supportedRoleOptions *SupportedRoleOptions) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(supportedRoleOptions.AccessPolicy) {
		_patch["access_policy"] = supportedRoleOptions.AccessPolicy
	}
	if !core.IsNil(supportedRoleOptions.AdditionalPropertiesForAccessPolicy) {
		_patch["additional_properties_for_access_policy"] = supportedRoleOptions.AdditionalPropertiesForAccessPolicy
	}
	if !core.IsNil(supportedRoleOptions.PolicyType) {
		_patch["policy_type"] = supportedRoleOptions.PolicyType
	}
	if !core.IsNil(supportedRoleOptions.AccountType) {
		_patch["account_type"] = supportedRoleOptions.AccountType
	}

	return
}

// UpdateCatalogDeploymentOptions : The UpdateCatalogDeployment options.
type UpdateCatalogDeploymentOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// The unique ID of this global catalog deployment.
	CatalogDeploymentID *string `json:"catalog_deployment_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_catalog_deployment.
	GlobalCatalogDeploymentPatch map[string]interface{} `json:"GlobalCatalogDeployment_patch" validate:"required"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCatalogDeploymentOptions : Instantiate UpdateCatalogDeploymentOptions
func (*PartnerCenterSellV1) NewUpdateCatalogDeploymentOptions(productID string, catalogProductID string, catalogPlanID string, catalogDeploymentID string, globalCatalogDeploymentPatch map[string]interface{}) *UpdateCatalogDeploymentOptions {
	return &UpdateCatalogDeploymentOptions{
		ProductID:                    core.StringPtr(productID),
		CatalogProductID:             core.StringPtr(catalogProductID),
		CatalogPlanID:                core.StringPtr(catalogPlanID),
		CatalogDeploymentID:          core.StringPtr(catalogDeploymentID),
		GlobalCatalogDeploymentPatch: globalCatalogDeploymentPatch,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *UpdateCatalogDeploymentOptions) SetProductID(productID string) *UpdateCatalogDeploymentOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *UpdateCatalogDeploymentOptions) SetCatalogProductID(catalogProductID string) *UpdateCatalogDeploymentOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *UpdateCatalogDeploymentOptions) SetCatalogPlanID(catalogPlanID string) *UpdateCatalogDeploymentOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetCatalogDeploymentID : Allow user to set CatalogDeploymentID
func (_options *UpdateCatalogDeploymentOptions) SetCatalogDeploymentID(catalogDeploymentID string) *UpdateCatalogDeploymentOptions {
	_options.CatalogDeploymentID = core.StringPtr(catalogDeploymentID)
	return _options
}

// SetGlobalCatalogDeploymentPatch : Allow user to set GlobalCatalogDeploymentPatch
func (_options *UpdateCatalogDeploymentOptions) SetGlobalCatalogDeploymentPatch(globalCatalogDeploymentPatch map[string]interface{}) *UpdateCatalogDeploymentOptions {
	_options.GlobalCatalogDeploymentPatch = globalCatalogDeploymentPatch
	return _options
}

// SetEnv : Allow user to set Env
func (_options *UpdateCatalogDeploymentOptions) SetEnv(env string) *UpdateCatalogDeploymentOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogDeploymentOptions) SetHeaders(param map[string]string) *UpdateCatalogDeploymentOptions {
	options.Headers = param
	return options
}

// UpdateCatalogPlanOptions : The UpdateCatalogPlan options.
type UpdateCatalogPlanOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// The unique ID of this global catalog plan.
	CatalogPlanID *string `json:"catalog_plan_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_catalog_plan.
	GlobalCatalogPlanPatch map[string]interface{} `json:"GlobalCatalogPlan_patch" validate:"required"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCatalogPlanOptions : Instantiate UpdateCatalogPlanOptions
func (*PartnerCenterSellV1) NewUpdateCatalogPlanOptions(productID string, catalogProductID string, catalogPlanID string, globalCatalogPlanPatch map[string]interface{}) *UpdateCatalogPlanOptions {
	return &UpdateCatalogPlanOptions{
		ProductID:              core.StringPtr(productID),
		CatalogProductID:       core.StringPtr(catalogProductID),
		CatalogPlanID:          core.StringPtr(catalogPlanID),
		GlobalCatalogPlanPatch: globalCatalogPlanPatch,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *UpdateCatalogPlanOptions) SetProductID(productID string) *UpdateCatalogPlanOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *UpdateCatalogPlanOptions) SetCatalogProductID(catalogProductID string) *UpdateCatalogPlanOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetCatalogPlanID : Allow user to set CatalogPlanID
func (_options *UpdateCatalogPlanOptions) SetCatalogPlanID(catalogPlanID string) *UpdateCatalogPlanOptions {
	_options.CatalogPlanID = core.StringPtr(catalogPlanID)
	return _options
}

// SetGlobalCatalogPlanPatch : Allow user to set GlobalCatalogPlanPatch
func (_options *UpdateCatalogPlanOptions) SetGlobalCatalogPlanPatch(globalCatalogPlanPatch map[string]interface{}) *UpdateCatalogPlanOptions {
	_options.GlobalCatalogPlanPatch = globalCatalogPlanPatch
	return _options
}

// SetEnv : Allow user to set Env
func (_options *UpdateCatalogPlanOptions) SetEnv(env string) *UpdateCatalogPlanOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogPlanOptions) SetHeaders(param map[string]string) *UpdateCatalogPlanOptions {
	options.Headers = param
	return options
}

// UpdateCatalogProductOptions : The UpdateCatalogProduct options.
type UpdateCatalogProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The unique ID of this global catalog product.
	CatalogProductID *string `json:"catalog_product_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_catalog_product.
	GlobalCatalogProductPatch map[string]interface{} `json:"GlobalCatalogProduct_patch" validate:"required"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCatalogProductOptions : Instantiate UpdateCatalogProductOptions
func (*PartnerCenterSellV1) NewUpdateCatalogProductOptions(productID string, catalogProductID string, globalCatalogProductPatch map[string]interface{}) *UpdateCatalogProductOptions {
	return &UpdateCatalogProductOptions{
		ProductID:                 core.StringPtr(productID),
		CatalogProductID:          core.StringPtr(catalogProductID),
		GlobalCatalogProductPatch: globalCatalogProductPatch,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *UpdateCatalogProductOptions) SetProductID(productID string) *UpdateCatalogProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetCatalogProductID : Allow user to set CatalogProductID
func (_options *UpdateCatalogProductOptions) SetCatalogProductID(catalogProductID string) *UpdateCatalogProductOptions {
	_options.CatalogProductID = core.StringPtr(catalogProductID)
	return _options
}

// SetGlobalCatalogProductPatch : Allow user to set GlobalCatalogProductPatch
func (_options *UpdateCatalogProductOptions) SetGlobalCatalogProductPatch(globalCatalogProductPatch map[string]interface{}) *UpdateCatalogProductOptions {
	_options.GlobalCatalogProductPatch = globalCatalogProductPatch
	return _options
}

// SetEnv : Allow user to set Env
func (_options *UpdateCatalogProductOptions) SetEnv(env string) *UpdateCatalogProductOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCatalogProductOptions) SetHeaders(param map[string]string) *UpdateCatalogProductOptions {
	options.Headers = param
	return options
}

// UpdateIamRegistrationOptions : The UpdateIamRegistration options.
type UpdateIamRegistrationOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// The approved programmatic name of the product.
	ProgrammaticName *string `json:"programmatic_name" validate:"required,ne="`

	// JSON Merge-Patch content for update_iam_registration.
	IamRegistrationPatch map[string]interface{} `json:"iam-registration-patch" validate:"required"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateIamRegistrationOptions : Instantiate UpdateIamRegistrationOptions
func (*PartnerCenterSellV1) NewUpdateIamRegistrationOptions(productID string, programmaticName string, iamRegistrationPatch map[string]interface{}) *UpdateIamRegistrationOptions {
	return &UpdateIamRegistrationOptions{
		ProductID:            core.StringPtr(productID),
		ProgrammaticName:     core.StringPtr(programmaticName),
		IamRegistrationPatch: iamRegistrationPatch,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *UpdateIamRegistrationOptions) SetProductID(productID string) *UpdateIamRegistrationOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetProgrammaticName : Allow user to set ProgrammaticName
func (_options *UpdateIamRegistrationOptions) SetProgrammaticName(programmaticName string) *UpdateIamRegistrationOptions {
	_options.ProgrammaticName = core.StringPtr(programmaticName)
	return _options
}

// SetIamRegistrationPatch : Allow user to set IamRegistrationPatch
func (_options *UpdateIamRegistrationOptions) SetIamRegistrationPatch(iamRegistrationPatch map[string]interface{}) *UpdateIamRegistrationOptions {
	_options.IamRegistrationPatch = iamRegistrationPatch
	return _options
}

// SetEnv : Allow user to set Env
func (_options *UpdateIamRegistrationOptions) SetEnv(env string) *UpdateIamRegistrationOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateIamRegistrationOptions) SetHeaders(param map[string]string) *UpdateIamRegistrationOptions {
	options.Headers = param
	return options
}

// UpdateOnboardingProductOptions : The UpdateOnboardingProduct options.
type UpdateOnboardingProductOptions struct {
	// The unique ID of the product.
	ProductID *string `json:"product_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_onboarding_product.
	OnboardingProductPatch map[string]interface{} `json:"OnboardingProduct_patch" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateOnboardingProductOptions : Instantiate UpdateOnboardingProductOptions
func (*PartnerCenterSellV1) NewUpdateOnboardingProductOptions(productID string, onboardingProductPatch map[string]interface{}) *UpdateOnboardingProductOptions {
	return &UpdateOnboardingProductOptions{
		ProductID:              core.StringPtr(productID),
		OnboardingProductPatch: onboardingProductPatch,
	}
}

// SetProductID : Allow user to set ProductID
func (_options *UpdateOnboardingProductOptions) SetProductID(productID string) *UpdateOnboardingProductOptions {
	_options.ProductID = core.StringPtr(productID)
	return _options
}

// SetOnboardingProductPatch : Allow user to set OnboardingProductPatch
func (_options *UpdateOnboardingProductOptions) SetOnboardingProductPatch(onboardingProductPatch map[string]interface{}) *UpdateOnboardingProductOptions {
	_options.OnboardingProductPatch = onboardingProductPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateOnboardingProductOptions) SetHeaders(param map[string]string) *UpdateOnboardingProductOptions {
	options.Headers = param
	return options
}

// UpdateRegistrationOptions : The UpdateRegistration options.
type UpdateRegistrationOptions struct {
	// The unique ID of your registration.
	RegistrationID *string `json:"registration_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_registration.
	RegistrationPatch map[string]interface{} `json:"registration-patch" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateRegistrationOptions : Instantiate UpdateRegistrationOptions
func (*PartnerCenterSellV1) NewUpdateRegistrationOptions(registrationID string, registrationPatch map[string]interface{}) *UpdateRegistrationOptions {
	return &UpdateRegistrationOptions{
		RegistrationID:    core.StringPtr(registrationID),
		RegistrationPatch: registrationPatch,
	}
}

// SetRegistrationID : Allow user to set RegistrationID
func (_options *UpdateRegistrationOptions) SetRegistrationID(registrationID string) *UpdateRegistrationOptions {
	_options.RegistrationID = core.StringPtr(registrationID)
	return _options
}

// SetRegistrationPatch : Allow user to set RegistrationPatch
func (_options *UpdateRegistrationOptions) SetRegistrationPatch(registrationPatch map[string]interface{}) *UpdateRegistrationOptions {
	_options.RegistrationPatch = registrationPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRegistrationOptions) SetHeaders(param map[string]string) *UpdateRegistrationOptions {
	options.Headers = param
	return options
}

// UpdateResourceBrokerOptions : The UpdateResourceBroker options.
type UpdateResourceBrokerOptions struct {
	// The unique identifier of the broker.
	BrokerID *string `json:"broker_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_resource_broker.
	BrokerPatch map[string]interface{} `json:"Broker_patch" validate:"required"`

	// The environment to fetch this object from.
	Env *string `json:"env,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateResourceBrokerOptions : Instantiate UpdateResourceBrokerOptions
func (*PartnerCenterSellV1) NewUpdateResourceBrokerOptions(brokerID string, brokerPatch map[string]interface{}) *UpdateResourceBrokerOptions {
	return &UpdateResourceBrokerOptions{
		BrokerID:    core.StringPtr(brokerID),
		BrokerPatch: brokerPatch,
	}
}

// SetBrokerID : Allow user to set BrokerID
func (_options *UpdateResourceBrokerOptions) SetBrokerID(brokerID string) *UpdateResourceBrokerOptions {
	_options.BrokerID = core.StringPtr(brokerID)
	return _options
}

// SetBrokerPatch : Allow user to set BrokerPatch
func (_options *UpdateResourceBrokerOptions) SetBrokerPatch(brokerPatch map[string]interface{}) *UpdateResourceBrokerOptions {
	_options.BrokerPatch = brokerPatch
	return _options
}

// SetEnv : Allow user to set Env
func (_options *UpdateResourceBrokerOptions) SetEnv(env string) *UpdateResourceBrokerOptions {
	_options.Env = core.StringPtr(env)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateResourceBrokerOptions) SetHeaders(param map[string]string) *UpdateResourceBrokerOptions {
	options.Headers = param
	return options
}
