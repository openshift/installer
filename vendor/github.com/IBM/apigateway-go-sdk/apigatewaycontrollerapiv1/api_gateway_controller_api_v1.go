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

// Package apigatewaycontrollerapiv1 : Operations and models for the ApiGatewayControllerApiV1 service
package apigatewaycontrollerapiv1

import (
	"fmt"

	common "github.com/IBM/apigateway-go-sdk/common"
	"github.com/IBM/go-sdk-core/v3/core"
)

// ApiGatewayControllerApiV1 : Primary REST API for creating and managing APIs within the IBM Cloud API Gateway service.
//
// Version: 1.0.0
type ApiGatewayControllerApiV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.us-south.apigw.cloud.ibm.com/controller"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "api_gateway_controller_api"

// ApiGatewayControllerApiV1Options : Service options
type ApiGatewayControllerApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewApiGatewayControllerApiV1UsingExternalConfig : constructs an instance of ApiGatewayControllerApiV1 with passed in options and external configuration.
func NewApiGatewayControllerApiV1UsingExternalConfig(options *ApiGatewayControllerApiV1Options) (apiGatewayControllerApi *ApiGatewayControllerApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	apiGatewayControllerApi, err = NewApiGatewayControllerApiV1(options)
	if err != nil {
		return
	}

	err = apiGatewayControllerApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = apiGatewayControllerApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewApiGatewayControllerApiV1 : constructs an instance of ApiGatewayControllerApiV1 with passed in options.
func NewApiGatewayControllerApiV1(options *ApiGatewayControllerApiV1Options) (service *ApiGatewayControllerApiV1, err error) {
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

	service = &ApiGatewayControllerApiV1{
		Service: baseService,
	}

	return
}

// SetServiceURL sets the service URL
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) SetServiceURL(url string) error {
	return apiGatewayControllerApi.Service.SetServiceURL(url)
}

// GetAllEndpoints : Get details for all Endpoints
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetAllEndpoints(getAllEndpointsOptions *GetAllEndpointsOptions) (result *[]V2Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAllEndpointsOptions, "getAllEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAllEndpointsOptions, "getAllEndpointsOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAllEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetAllEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if getAllEndpointsOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*getAllEndpointsOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*getAllEndpointsOptions.ServiceInstanceCrn))
	if getAllEndpointsOptions.ProviderID != nil {
		builder.AddQuery("provider_id", fmt.Sprint(*getAllEndpointsOptions.ProviderID))
	}
	if getAllEndpointsOptions.Shared != nil {
		builder.AddQuery("shared", fmt.Sprint(*getAllEndpointsOptions.Shared))
	}
	if getAllEndpointsOptions.Managed != nil {
		builder.AddQuery("managed", fmt.Sprint(*getAllEndpointsOptions.Managed))
	}
	if getAllEndpointsOptions.Swagger != nil {
		builder.AddQuery("swagger", fmt.Sprint(*getAllEndpointsOptions.Swagger))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make([]map[string]interface{}, 1))
	if err == nil {
		s, ok := response.Result.([]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		slice, e := UnmarshalV2EndpointSlice(s)
		result = &slice
		err = e
		response.Result = result
	}

	return
}

// CreateEndpoint : Create an Endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) CreateEndpoint(createEndpointOptions *CreateEndpointOptions) (result *V2Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createEndpointOptions, "createEndpointOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createEndpointOptions, "createEndpointOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.POST)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range createEndpointOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "CreateEndpoint")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createEndpointOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*createEndpointOptions.Authorization))
	}

	body := make(map[string]interface{})
	if createEndpointOptions.ArtifactID != nil {
		body["artifact_id"] = createEndpointOptions.ArtifactID
	}
	if createEndpointOptions.ParentCrn != nil {
		body["parent_crn"] = createEndpointOptions.ParentCrn
	}
	if createEndpointOptions.ServiceInstanceCrn != nil {
		body["service_instance_crn"] = createEndpointOptions.ServiceInstanceCrn
	}
	if createEndpointOptions.Name != nil {
		body["name"] = createEndpointOptions.Name
	}
	if createEndpointOptions.Routes != nil {
		body["routes"] = createEndpointOptions.Routes
	}
	if createEndpointOptions.Managed != nil {
		body["managed"] = createEndpointOptions.Managed
	}
	if createEndpointOptions.Metadata != nil {
		body["metadata"] = createEndpointOptions.Metadata
	}
	if createEndpointOptions.OpenApiDoc != nil {
		body["open_api_doc"] = createEndpointOptions.OpenApiDoc
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Endpoint(m)
		response.Result = result
	}

	return
}

// GetEndpoint : Get details for a given Endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetEndpoint(getEndpointOptions *GetEndpointOptions) (result *V2Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEndpointOptions, "getEndpointOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEndpointOptions, "getEndpointOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints"}
	pathParameters := []string{*getEndpointOptions.ID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEndpointOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetEndpoint")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if getEndpointOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*getEndpointOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*getEndpointOptions.ServiceInstanceCrn))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Endpoint(m)
		response.Result = result
	}

	return
}

// UpdateEndpoint : Update an endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) UpdateEndpoint(updateEndpointOptions *UpdateEndpointOptions) (result *V2Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEndpointOptions, "updateEndpointOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEndpointOptions, "updateEndpointOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints"}
	pathParameters := []string{*updateEndpointOptions.ID}

	builder := core.NewRequestBuilder(core.PUT)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateEndpointOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "UpdateEndpoint")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateEndpointOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*updateEndpointOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*updateEndpointOptions.ServiceInstanceCrn))

	body := make(map[string]interface{})
	if updateEndpointOptions.NewArtifactID != nil {
		body["artifact_id"] = updateEndpointOptions.NewArtifactID
	}
	if updateEndpointOptions.NewParentCrn != nil {
		body["parent_crn"] = updateEndpointOptions.NewParentCrn
	}
	if updateEndpointOptions.NewServiceInstanceCrn != nil {
		body["service_instance_crn"] = updateEndpointOptions.NewServiceInstanceCrn
	}
	if updateEndpointOptions.NewName != nil {
		body["name"] = updateEndpointOptions.NewName
	}
	if updateEndpointOptions.NewRoutes != nil {
		body["routes"] = updateEndpointOptions.NewRoutes
	}
	if updateEndpointOptions.NewManaged != nil {
		body["managed"] = updateEndpointOptions.NewManaged
	}
	if updateEndpointOptions.NewMetadata != nil {
		body["metadata"] = updateEndpointOptions.NewMetadata
	}
	if updateEndpointOptions.NewOpenApiDoc != nil {
		body["open_api_doc"] = updateEndpointOptions.NewOpenApiDoc
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Endpoint(m)
		response.Result = result
	}

	return
}

// DeleteEndpoint : Delete an Endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) DeleteEndpoint(deleteEndpointOptions *DeleteEndpointOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEndpointOptions, "deleteEndpointOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEndpointOptions, "deleteEndpointOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints"}
	pathParameters := []string{*deleteEndpointOptions.ID}

	builder := core.NewRequestBuilder(core.DELETE)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEndpointOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "DeleteEndpoint")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteEndpointOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*deleteEndpointOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*deleteEndpointOptions.ServiceInstanceCrn))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, nil)

	return
}

// GetEndpointSwagger : Get the OpenAPI doc for a given Endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetEndpointSwagger(getEndpointSwaggerOptions *GetEndpointSwaggerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEndpointSwaggerOptions, "getEndpointSwaggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEndpointSwaggerOptions, "getEndpointSwaggerOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints", "swagger"}
	pathParameters := []string{*getEndpointSwaggerOptions.ID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEndpointSwaggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetEndpointSwagger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if getEndpointSwaggerOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*getEndpointSwaggerOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*getEndpointSwaggerOptions.ServiceInstanceCrn))
	if getEndpointSwaggerOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*getEndpointSwaggerOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))

	return
}

// EndpointActions : Execute actions for a given Endpoint
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) EndpointActions(endpointActionsOptions *EndpointActionsOptions) (result *V2Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(endpointActionsOptions, "endpointActionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(endpointActionsOptions, "endpointActionsOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints", "actions"}
	pathParameters := []string{*endpointActionsOptions.ID}

	builder := core.NewRequestBuilder(core.POST)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range endpointActionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "EndpointActions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if endpointActionsOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*endpointActionsOptions.Authorization))
	}

	builder.AddQuery("service_instance_crn", fmt.Sprint(*endpointActionsOptions.ServiceInstanceCrn))
	builder.AddQuery("provider_id", fmt.Sprint(*endpointActionsOptions.ProviderID))

	body := make(map[string]interface{})
	if endpointActionsOptions.Type != nil {
		body["type"] = endpointActionsOptions.Type
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Endpoint(m)
		response.Result = result
	}

	return
}

// EndpointSummary : Get provider sorted summary about all Endpoints
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) EndpointSummary(endpointSummaryOptions *EndpointSummaryOptions) (result *[]V2EndpointSummary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(endpointSummaryOptions, "endpointSummaryOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(endpointSummaryOptions, "endpointSummaryOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/endpoints/summary"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range endpointSummaryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "EndpointSummary")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if endpointSummaryOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*endpointSummaryOptions.Authorization))
	}

	builder.AddQuery("account_id", fmt.Sprint(*endpointSummaryOptions.AccountID))
	if endpointSummaryOptions.ServiceInstanceCrn != nil {
		builder.AddQuery("service_instance_crn", fmt.Sprint(*endpointSummaryOptions.ServiceInstanceCrn))
	}
	if endpointSummaryOptions.Swagger != nil {
		builder.AddQuery("swagger", fmt.Sprint(*endpointSummaryOptions.Swagger))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make([]map[string]interface{}, 1))
	if err == nil {
		s, ok := response.Result.([]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		slice, e := UnmarshalV2EndpointSummarySlice(s)
		result = &slice
		err = e
		response.Result = result
	}

	return
}

// GetAllSubscriptions : Get all subscriptions tied to a given artifact
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetAllSubscriptions(getAllSubscriptionsOptions *GetAllSubscriptionsOptions) (result *[]V2Subscription, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAllSubscriptionsOptions, "getAllSubscriptionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getAllSubscriptionsOptions, "getAllSubscriptionsOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getAllSubscriptionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetAllSubscriptions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if getAllSubscriptionsOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*getAllSubscriptionsOptions.Authorization))
	}

	builder.AddQuery("artifact_id", fmt.Sprint(*getAllSubscriptionsOptions.ArtifactID))
	if getAllSubscriptionsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*getAllSubscriptionsOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make([]map[string]interface{}, 1))
	if err == nil {
		s, ok := response.Result.([]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		slice, e := UnmarshalV2SubscriptionSlice(s)
		result = &slice
		err = e
		response.Result = result
	}

	return
}

// CreateSubscription : Create a subscription for an artifact
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) CreateSubscription(createSubscriptionOptions *CreateSubscriptionOptions) (result *V2Subscription, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSubscriptionOptions, "createSubscriptionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSubscriptionOptions, "createSubscriptionOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.POST)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSubscriptionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "CreateSubscription")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createSubscriptionOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*createSubscriptionOptions.Authorization))
	}

	body := make(map[string]interface{})
	if createSubscriptionOptions.ClientID != nil {
		body["client_id"] = createSubscriptionOptions.ClientID
	}
	if createSubscriptionOptions.ArtifactID != nil {
		body["artifact_id"] = createSubscriptionOptions.ArtifactID
	}
	if createSubscriptionOptions.ClientSecret != nil {
		body["client_secret"] = createSubscriptionOptions.ClientSecret
	}
	if createSubscriptionOptions.GenerateSecret != nil {
		body["generate_secret"] = createSubscriptionOptions.GenerateSecret
	}
	if createSubscriptionOptions.AccountID != nil {
		body["account_id"] = createSubscriptionOptions.AccountID
	}
	if createSubscriptionOptions.Name != nil {
		body["name"] = createSubscriptionOptions.Name
	}
	if createSubscriptionOptions.Type != nil {
		body["type"] = createSubscriptionOptions.Type
	}
	if createSubscriptionOptions.Metadata != nil {
		body["metadata"] = createSubscriptionOptions.Metadata
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Subscription(m)
		response.Result = result
	}

	return
}

// GetSubscription : Get subscription for a given clientid
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetSubscription(getSubscriptionOptions *GetSubscriptionOptions) (result *V2Subscription, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSubscriptionOptions, "getSubscriptionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSubscriptionOptions, "getSubscriptionOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions"}
	pathParameters := []string{*getSubscriptionOptions.ID}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSubscriptionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetSubscription")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	if getSubscriptionOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*getSubscriptionOptions.Authorization))
	}

	builder.AddQuery("artifact_id", fmt.Sprint(*getSubscriptionOptions.ArtifactID))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Subscription(m)
		response.Result = result
	}

	return
}

// UpdateSubscription : Update a subscription
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) UpdateSubscription(updateSubscriptionOptions *UpdateSubscriptionOptions) (result *V2Subscription, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSubscriptionOptions, "updateSubscriptionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSubscriptionOptions, "updateSubscriptionOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions"}
	pathParameters := []string{*updateSubscriptionOptions.ID}

	builder := core.NewRequestBuilder(core.PUT)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSubscriptionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "UpdateSubscription")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateSubscriptionOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*updateSubscriptionOptions.Authorization))
	}

	builder.AddQuery("artifact_id", fmt.Sprint(*updateSubscriptionOptions.ArtifactID))

	body := make(map[string]interface{})
	if updateSubscriptionOptions.NewClientID != nil {
		body["client_id"] = updateSubscriptionOptions.NewClientID
	}
	if updateSubscriptionOptions.NewClientSecret != nil {
		body["client_secret"] = updateSubscriptionOptions.NewClientSecret
	}
	if updateSubscriptionOptions.NewArtifactID != nil {
		body["artifact_id"] = updateSubscriptionOptions.NewArtifactID
	}
	if updateSubscriptionOptions.NewAccountID != nil {
		body["account_id"] = updateSubscriptionOptions.NewAccountID
	}
	if updateSubscriptionOptions.NewName != nil {
		body["name"] = updateSubscriptionOptions.NewName
	}
	if updateSubscriptionOptions.NewType != nil {
		body["type"] = updateSubscriptionOptions.NewType
	}
	if updateSubscriptionOptions.NewMetadata != nil {
		body["metadata"] = updateSubscriptionOptions.NewMetadata
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Subscription(m)
		response.Result = result
	}

	return
}

// DeleteSubscription : Delete a subscription
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) DeleteSubscription(deleteSubscriptionOptions *DeleteSubscriptionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSubscriptionOptions, "deleteSubscriptionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSubscriptionOptions, "deleteSubscriptionOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions"}
	pathParameters := []string{*deleteSubscriptionOptions.ID}

	builder := core.NewRequestBuilder(core.DELETE)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSubscriptionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "DeleteSubscription")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if deleteSubscriptionOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*deleteSubscriptionOptions.Authorization))
	}

	builder.AddQuery("artifact_id", fmt.Sprint(*deleteSubscriptionOptions.ArtifactID))

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, nil)

	return
}

// GetSubscriptionArtifact : Get artifact associated to a subscription
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) GetSubscriptionArtifact(getSubscriptionArtifactOptions *GetSubscriptionArtifactOptions) (result *InlineResponse200, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSubscriptionArtifactOptions, "getSubscriptionArtifactOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSubscriptionArtifactOptions, "getSubscriptionArtifactOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions/artifact"}
	pathParameters := []string{}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getSubscriptionArtifactOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "GetSubscriptionArtifact")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("artifact_id", fmt.Sprint(*getSubscriptionArtifactOptions.ArtifactID))
	if getSubscriptionArtifactOptions.ClientID != nil {
		builder.AddQuery("client_id", fmt.Sprint(*getSubscriptionArtifactOptions.ClientID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalInlineResponse200(m)
		response.Result = result
	}

	return
}

// AddSubscriptionSecret : Add a Subscription Secret
func (apiGatewayControllerApi *ApiGatewayControllerApiV1) AddSubscriptionSecret(addSubscriptionSecretOptions *AddSubscriptionSecretOptions) (result *V2Subscription, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(addSubscriptionSecretOptions, "addSubscriptionSecretOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(addSubscriptionSecretOptions, "addSubscriptionSecretOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"v2/subscriptions", "secret"}
	pathParameters := []string{*addSubscriptionSecretOptions.ID}

	builder := core.NewRequestBuilder(core.PUT)
	_, err = builder.ConstructHTTPURL(apiGatewayControllerApi.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range addSubscriptionSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("api_gateway_controller_api", "V1", "AddSubscriptionSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if addSubscriptionSecretOptions.Authorization != nil {
		builder.AddHeader("authorization", fmt.Sprint(*addSubscriptionSecretOptions.Authorization))
	}

	builder.AddQuery("artifact_id", fmt.Sprint(*addSubscriptionSecretOptions.ArtifactID))

	body := make(map[string]interface{})
	if addSubscriptionSecretOptions.ClientSecret != nil {
		body["client_secret"] = addSubscriptionSecretOptions.ClientSecret
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = apiGatewayControllerApi.Service.Request(request, make(map[string]interface{}))
	if err == nil {
		m, ok := response.Result.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("an error occurred while processing the operation response")
			return
		}
		result, err = UnmarshalV2Subscription(m)
		response.Result = result
	}

	return
}

// AddSubscriptionSecretOptions : The AddSubscriptionSecret options.
type AddSubscriptionSecretOptions struct {
	// Client Id.
	ID *string `json:"id" validate:"required"`

	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// Client id.
	Authorization *string `json:"authorization" validate:"required"`

	ClientSecret *string `json:"client_secret,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewAddSubscriptionSecretOptions : Instantiate AddSubscriptionSecretOptions
func (*ApiGatewayControllerApiV1) NewAddSubscriptionSecretOptions(id string, artifactID string, authorization string) *AddSubscriptionSecretOptions {
	return &AddSubscriptionSecretOptions{
		ID:            core.StringPtr(id),
		ArtifactID:    core.StringPtr(artifactID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *AddSubscriptionSecretOptions) SetID(id string) *AddSubscriptionSecretOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *AddSubscriptionSecretOptions) SetArtifactID(artifactID string) *AddSubscriptionSecretOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *AddSubscriptionSecretOptions) SetAuthorization(authorization string) *AddSubscriptionSecretOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetClientSecret : Allow user to set ClientSecret
func (options *AddSubscriptionSecretOptions) SetClientSecret(clientSecret string) *AddSubscriptionSecretOptions {
	options.ClientSecret = core.StringPtr(clientSecret)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *AddSubscriptionSecretOptions) SetHeaders(param map[string]string) *AddSubscriptionSecretOptions {
	options.Headers = param
	return options
}

// CreateEndpointOptions : The CreateEndpoint options.
type CreateEndpointOptions struct {
	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// The endpoint ID.
	ArtifactID *string `json:"artifact_id",omitempt`

	// The API Gateway service instance CRN.
	ParentCrn *string `json:"parent_crn" validate:"required"`

	// The API Gateway service instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// The endpoint's name.
	Name *string `json:"name validate:"required"`

	// Invokable endpoint routes.
	Routes []string `json:"routes,omitempty"`

	// Is the endpoint managed?.
	Managed *bool `json:"managed,omitempty"`

	Metadata interface{} `json:"metadata,omitempty"`

	// The OpenAPI document.
	OpenApiDoc interface{} `json:"open_api_doc validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateEndpointOptions : Instantiate CreateEndpointOptions
func (*ApiGatewayControllerApiV1) NewCreateEndpointOptions(authorization string, artifactID string, parentCrn string, serviceInstanceCrn string) *CreateEndpointOptions {
	return &CreateEndpointOptions{
		Authorization:      core.StringPtr(authorization),
		ArtifactID:         core.StringPtr(artifactID),
		ParentCrn:          core.StringPtr(parentCrn),
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
	}
}

// SetAuthorization : Allow user to set Authorization
func (options *CreateEndpointOptions) SetAuthorization(authorization string) *CreateEndpointOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *CreateEndpointOptions) SetArtifactID(artifactID string) *CreateEndpointOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetParentCrn : Allow user to set ParentCrn
func (options *CreateEndpointOptions) SetParentCrn(parentCrn string) *CreateEndpointOptions {
	options.ParentCrn = core.StringPtr(parentCrn)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *CreateEndpointOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *CreateEndpointOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetName : Allow user to set Name
func (options *CreateEndpointOptions) SetName(name string) *CreateEndpointOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetRoutes : Allow user to set Routes
func (options *CreateEndpointOptions) SetRoutes(routes []string) *CreateEndpointOptions {
	options.Routes = routes
	return options
}

// SetManaged : Allow user to set Managed
func (options *CreateEndpointOptions) SetManaged(managed bool) *CreateEndpointOptions {
	options.Managed = core.BoolPtr(managed)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *CreateEndpointOptions) SetMetadata(metadata interface{}) *CreateEndpointOptions {
	options.Metadata = metadata
	return options
}

// SetOpenApiDoc : Allow user to set OpenApiDoc
func (options *CreateEndpointOptions) SetOpenApiDoc(openApiDoc interface{}) *CreateEndpointOptions {
	options.OpenApiDoc = openApiDoc
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateEndpointOptions) SetHeaders(param map[string]string) *CreateEndpointOptions {
	options.Headers = param
	return options
}

// CreateSubscriptionOptions : The CreateSubscription options.
type CreateSubscriptionOptions struct {
	// User bearer token.
	Authorization *string `json:"authorization" validate:"required"`

	ClientID *string `json:"client_id",omitempty`

	ArtifactID *string `json:"artifact_id" validate:"required"`

	ClientSecret *string `json:"client_secret,omitempty"`

	GenerateSecret *bool `json:"generate_secret, omitempty"`

	AccountID *string `json:"account_id,omitempty"`

	Name *string `json:"name,omitempty"`

	Type *string `json:"type,omitempty"`

	Metadata interface{} `json:"metadata,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSubscriptionOptions : Instantiate CreateSubscriptionOptions
func (*ApiGatewayControllerApiV1) NewCreateSubscriptionOptions(authorization string, clientID string, artifactID string) *CreateSubscriptionOptions {
	return &CreateSubscriptionOptions{
		Authorization: core.StringPtr(authorization),
		ClientID:      core.StringPtr(clientID),
		ArtifactID:    core.StringPtr(artifactID),
	}
}

// SetAuthorization : Allow user to set Authorization
func (options *CreateSubscriptionOptions) SetAuthorization(authorization string) *CreateSubscriptionOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *CreateSubscriptionOptions) SetClientID(clientID string) *CreateSubscriptionOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *CreateSubscriptionOptions) SetArtifactID(artifactID string) *CreateSubscriptionOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetClientSecret : Allow user to set ClientSecret
func (options *CreateSubscriptionOptions) SetClientSecret(clientSecret string) *CreateSubscriptionOptions {
	options.ClientSecret = core.StringPtr(clientSecret)
	return options
}

// SetAccountID : Allow user to set AccountID
func (options *CreateSubscriptionOptions) SetAccountID(accountID string) *CreateSubscriptionOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateSubscriptionOptions) SetName(name string) *CreateSubscriptionOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetType : Allow user to set Type
func (options *CreateSubscriptionOptions) SetType(typeVar string) *CreateSubscriptionOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetMetadata : Allow user to set Metadata
func (options *CreateSubscriptionOptions) SetMetadata(metadata interface{}) *CreateSubscriptionOptions {
	options.Metadata = metadata
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSubscriptionOptions) SetHeaders(param map[string]string) *CreateSubscriptionOptions {
	options.Headers = param
	return options
}

// DeleteEndpointOptions : The DeleteEndpoint options.
type DeleteEndpointOptions struct {
	// Endpoint id.
	ID *string `json:"id" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEndpointOptions : Instantiate DeleteEndpointOptions
func (*ApiGatewayControllerApiV1) NewDeleteEndpointOptions(id string, serviceInstanceCrn string, authorization string) *DeleteEndpointOptions {
	return &DeleteEndpointOptions{
		ID:                 core.StringPtr(id),
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
		Authorization:      core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *DeleteEndpointOptions) SetID(id string) *DeleteEndpointOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *DeleteEndpointOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *DeleteEndpointOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *DeleteEndpointOptions) SetAuthorization(authorization string) *DeleteEndpointOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEndpointOptions) SetHeaders(param map[string]string) *DeleteEndpointOptions {
	options.Headers = param
	return options
}

// DeleteSubscriptionOptions : The DeleteSubscription options.
type DeleteSubscriptionOptions struct {
	// Client Id.
	ID *string `json:"id" validate:"required"`

	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// User bearer token.
	Authorization *string `json:"authorization" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSubscriptionOptions : Instantiate DeleteSubscriptionOptions
func (*ApiGatewayControllerApiV1) NewDeleteSubscriptionOptions(id string, artifactID string, authorization string) *DeleteSubscriptionOptions {
	return &DeleteSubscriptionOptions{
		ID:            core.StringPtr(id),
		ArtifactID:    core.StringPtr(artifactID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *DeleteSubscriptionOptions) SetID(id string) *DeleteSubscriptionOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *DeleteSubscriptionOptions) SetArtifactID(artifactID string) *DeleteSubscriptionOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *DeleteSubscriptionOptions) SetAuthorization(authorization string) *DeleteSubscriptionOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSubscriptionOptions) SetHeaders(param map[string]string) *DeleteSubscriptionOptions {
	options.Headers = param
	return options
}

// EndpointActionsOptions : The EndpointActions options.
type EndpointActionsOptions struct {
	// Endpoint Id.
	ID *string `json:"id" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// Provider Id.
	ProviderID *string `json:"provider_id" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	Type *string `json:"type" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewEndpointActionsOptions : Instantiate EndpointActionsOptions
func (*ApiGatewayControllerApiV1) NewEndpointActionsOptions(id string, serviceInstanceCrn string, providerID string, authorization string, typeVar string) *EndpointActionsOptions {
	return &EndpointActionsOptions{
		ID:                 core.StringPtr(id),
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
		ProviderID:         core.StringPtr(providerID),
		Authorization:      core.StringPtr(authorization),
		Type:               core.StringPtr(typeVar),
	}
}

// SetID : Allow user to set ID
func (options *EndpointActionsOptions) SetID(id string) *EndpointActionsOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *EndpointActionsOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *EndpointActionsOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetProviderID : Allow user to set ProviderID
func (options *EndpointActionsOptions) SetProviderID(providerID string) *EndpointActionsOptions {
	options.ProviderID = core.StringPtr(providerID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *EndpointActionsOptions) SetAuthorization(authorization string) *EndpointActionsOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetType : Allow user to set Type
func (options *EndpointActionsOptions) SetType(typeVar string) *EndpointActionsOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EndpointActionsOptions) SetHeaders(param map[string]string) *EndpointActionsOptions {
	options.Headers = param
	return options
}

// EndpointSummaryOptions : The EndpointSummary options.
type EndpointSummaryOptions struct {
	// User account ID.
	AccountID *string `json:"account_id" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn,omitempty"`

	// Return OpenAPI doc with list results. Possible values are ['provider', 'consumer']. Defaults to 'provider'.
	Swagger *string `json:"swagger,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewEndpointSummaryOptions : Instantiate EndpointSummaryOptions
func (*ApiGatewayControllerApiV1) NewEndpointSummaryOptions(accountID string, authorization string) *EndpointSummaryOptions {
	return &EndpointSummaryOptions{
		AccountID:     core.StringPtr(accountID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetAccountID : Allow user to set AccountID
func (options *EndpointSummaryOptions) SetAccountID(accountID string) *EndpointSummaryOptions {
	options.AccountID = core.StringPtr(accountID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *EndpointSummaryOptions) SetAuthorization(authorization string) *EndpointSummaryOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *EndpointSummaryOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *EndpointSummaryOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetSwagger : Allow user to set Swagger
func (options *EndpointSummaryOptions) SetSwagger(swagger string) *EndpointSummaryOptions {
	options.Swagger = core.StringPtr(swagger)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EndpointSummaryOptions) SetHeaders(param map[string]string) *EndpointSummaryOptions {
	options.Headers = param
	return options
}

// GetAllEndpointsOptions : The GetAllEndpoints options.
type GetAllEndpointsOptions struct {
	// The API Gateway service instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// Your IBM Cloud Identity and Access Management (IAM) token. To retrieve your IAM token, run `ibmcloud iam
	// oauth-tokens`.
	Authorization *string `json:"authorization" validate:"required"`

	// Provider Id.
	ProviderID *string `json:"provider_id,omitempty"`

	// Only return shared endpoints.
	Shared *bool `json:"shared,omitempty"`

	// Only return managed endpoints.
	Managed *bool `json:"managed,omitempty"`

	// Return OpenAPI doc with list results. Possible values are ['provider', 'consumer']. Defaults to 'provider'.
	Swagger *string `json:"swagger,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAllEndpointsOptions : Instantiate GetAllEndpointsOptions
func (*ApiGatewayControllerApiV1) NewGetAllEndpointsOptions(serviceInstanceCrn string, authorization string) *GetAllEndpointsOptions {
	return &GetAllEndpointsOptions{
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
		Authorization:      core.StringPtr(authorization),
	}
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *GetAllEndpointsOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *GetAllEndpointsOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *GetAllEndpointsOptions) SetAuthorization(authorization string) *GetAllEndpointsOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetProviderID : Allow user to set ProviderID
func (options *GetAllEndpointsOptions) SetProviderID(providerID string) *GetAllEndpointsOptions {
	options.ProviderID = core.StringPtr(providerID)
	return options
}

// SetShared : Allow user to set Shared
func (options *GetAllEndpointsOptions) SetShared(shared bool) *GetAllEndpointsOptions {
	options.Shared = core.BoolPtr(shared)
	return options
}

// SetManaged : Allow user to set Managed
func (options *GetAllEndpointsOptions) SetManaged(managed bool) *GetAllEndpointsOptions {
	options.Managed = core.BoolPtr(managed)
	return options
}

// SetSwagger : Allow user to set Swagger
func (options *GetAllEndpointsOptions) SetSwagger(swagger string) *GetAllEndpointsOptions {
	options.Swagger = core.StringPtr(swagger)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllEndpointsOptions) SetHeaders(param map[string]string) *GetAllEndpointsOptions {
	options.Headers = param
	return options
}

// GetAllSubscriptionsOptions : The GetAllSubscriptions options.
type GetAllSubscriptionsOptions struct {
	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// User bearer token.
	Authorization *string `json:"authorization" validate:"required"`

	// Subscription type.
	Type *string `json:"type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAllSubscriptionsOptions : Instantiate GetAllSubscriptionsOptions
func (*ApiGatewayControllerApiV1) NewGetAllSubscriptionsOptions(artifactID string, authorization string) *GetAllSubscriptionsOptions {
	return &GetAllSubscriptionsOptions{
		ArtifactID:    core.StringPtr(artifactID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetArtifactID : Allow user to set ArtifactID
func (options *GetAllSubscriptionsOptions) SetArtifactID(artifactID string) *GetAllSubscriptionsOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *GetAllSubscriptionsOptions) SetAuthorization(authorization string) *GetAllSubscriptionsOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetType : Allow user to set Type
func (options *GetAllSubscriptionsOptions) SetType(typeVar string) *GetAllSubscriptionsOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllSubscriptionsOptions) SetHeaders(param map[string]string) *GetAllSubscriptionsOptions {
	options.Headers = param
	return options
}

// GetEndpointOptions : The GetEndpoint options.
type GetEndpointOptions struct {
	// Endpoint Id.
	ID *string `json:"id" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEndpointOptions : Instantiate GetEndpointOptions
func (*ApiGatewayControllerApiV1) NewGetEndpointOptions(id string, serviceInstanceCrn string, authorization string) *GetEndpointOptions {
	return &GetEndpointOptions{
		ID:                 core.StringPtr(id),
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
		Authorization:      core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *GetEndpointOptions) SetID(id string) *GetEndpointOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *GetEndpointOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *GetEndpointOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *GetEndpointOptions) SetAuthorization(authorization string) *GetEndpointOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEndpointOptions) SetHeaders(param map[string]string) *GetEndpointOptions {
	options.Headers = param
	return options
}

// GetEndpointSwaggerOptions : The GetEndpointSwagger options.
type GetEndpointSwaggerOptions struct {
	// Endpoint Id.
	ID *string `json:"id" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// Type of swagger to retrieve.
	Type *string `json:"type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetEndpointSwaggerOptions.Type property.
// Type of swagger to retrieve.
const (
	GetEndpointSwaggerOptions_Type_JSON = "json"
	GetEndpointSwaggerOptions_Type_Yaml = "yaml"
)

// NewGetEndpointSwaggerOptions : Instantiate GetEndpointSwaggerOptions
func (*ApiGatewayControllerApiV1) NewGetEndpointSwaggerOptions(id string, serviceInstanceCrn string, authorization string) *GetEndpointSwaggerOptions {
	return &GetEndpointSwaggerOptions{
		ID:                 core.StringPtr(id),
		ServiceInstanceCrn: core.StringPtr(serviceInstanceCrn),
		Authorization:      core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *GetEndpointSwaggerOptions) SetID(id string) *GetEndpointSwaggerOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *GetEndpointSwaggerOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *GetEndpointSwaggerOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *GetEndpointSwaggerOptions) SetAuthorization(authorization string) *GetEndpointSwaggerOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetType : Allow user to set Type
func (options *GetEndpointSwaggerOptions) SetType(typeVar string) *GetEndpointSwaggerOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEndpointSwaggerOptions) SetHeaders(param map[string]string) *GetEndpointSwaggerOptions {
	options.Headers = param
	return options
}

// GetSubscriptionArtifactOptions : The GetSubscriptionArtifact options.
type GetSubscriptionArtifactOptions struct {
	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// Client Id.
	ClientID *string `json:"client_id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSubscriptionArtifactOptions : Instantiate GetSubscriptionArtifactOptions
func (*ApiGatewayControllerApiV1) NewGetSubscriptionArtifactOptions(artifactID string) *GetSubscriptionArtifactOptions {
	return &GetSubscriptionArtifactOptions{
		ArtifactID: core.StringPtr(artifactID),
	}
}

// SetArtifactID : Allow user to set ArtifactID
func (options *GetSubscriptionArtifactOptions) SetArtifactID(artifactID string) *GetSubscriptionArtifactOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetClientID : Allow user to set ClientID
func (options *GetSubscriptionArtifactOptions) SetClientID(clientID string) *GetSubscriptionArtifactOptions {
	options.ClientID = core.StringPtr(clientID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSubscriptionArtifactOptions) SetHeaders(param map[string]string) *GetSubscriptionArtifactOptions {
	options.Headers = param
	return options
}

// GetSubscriptionOptions : The GetSubscription options.
type GetSubscriptionOptions struct {
	// Client Id.
	ID *string `json:"id" validate:"required"`

	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// User bearer token.
	Authorization *string `json:"authorization" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetSubscriptionOptions : Instantiate GetSubscriptionOptions
func (*ApiGatewayControllerApiV1) NewGetSubscriptionOptions(id string, artifactID string, authorization string) *GetSubscriptionOptions {
	return &GetSubscriptionOptions{
		ID:            core.StringPtr(id),
		ArtifactID:    core.StringPtr(artifactID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *GetSubscriptionOptions) SetID(id string) *GetSubscriptionOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *GetSubscriptionOptions) SetArtifactID(artifactID string) *GetSubscriptionOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *GetSubscriptionOptions) SetAuthorization(authorization string) *GetSubscriptionOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSubscriptionOptions) SetHeaders(param map[string]string) *GetSubscriptionOptions {
	options.Headers = param
	return options
}

// UpdateEndpointOptions : The UpdateEndpoint options.
type UpdateEndpointOptions struct {
	// Endpoint Id.
	ID *string `json:"id" validate:"required"`

	// Service Instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn" validate:"required"`

	// User IAM token.
	Authorization *string `json:"authorization" validate:"required"`

	// The endpoint ID.
	NewArtifactID *string `json:"new_artifact_id" validate:"required"`

	// The API Gateway service instance CRN.
	NewParentCrn *string `json:"new_parent_crn" validate:"required"`

	// The API Gateway service instance CRN.
	NewServiceInstanceCrn *string `json:"new_service_instance_crn" validate:"required"`

	// The endpoint's name.
	NewName *string `json:"new_name,omitempty"`

	// Invokable endpoint routes.
	NewRoutes []string `json:"new_routes,omitempty"`

	// Is the endpoint managed?.
	NewManaged *bool `json:"new_managed,omitempty"`

	NewMetadata interface{} `json:"new_metadata,omitempty"`

	// The OpenAPI document.
	NewOpenApiDoc interface{} `json:"new_open_api_doc,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEndpointOptions : Instantiate UpdateEndpointOptions
func (*ApiGatewayControllerApiV1) NewUpdateEndpointOptions(id string, serviceInstanceCrn string, authorization string, newArtifactID string, newParentCrn string, newServiceInstanceCrn string) *UpdateEndpointOptions {
	return &UpdateEndpointOptions{
		ID:                    core.StringPtr(id),
		ServiceInstanceCrn:    core.StringPtr(serviceInstanceCrn),
		Authorization:         core.StringPtr(authorization),
		NewArtifactID:         core.StringPtr(newArtifactID),
		NewParentCrn:          core.StringPtr(newParentCrn),
		NewServiceInstanceCrn: core.StringPtr(newServiceInstanceCrn),
	}
}

// SetID : Allow user to set ID
func (options *UpdateEndpointOptions) SetID(id string) *UpdateEndpointOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetServiceInstanceCrn : Allow user to set ServiceInstanceCrn
func (options *UpdateEndpointOptions) SetServiceInstanceCrn(serviceInstanceCrn string) *UpdateEndpointOptions {
	options.ServiceInstanceCrn = core.StringPtr(serviceInstanceCrn)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *UpdateEndpointOptions) SetAuthorization(authorization string) *UpdateEndpointOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetNewArtifactID : Allow user to set NewArtifactID
func (options *UpdateEndpointOptions) SetNewArtifactID(newArtifactID string) *UpdateEndpointOptions {
	options.NewArtifactID = core.StringPtr(newArtifactID)
	return options
}

// SetNewParentCrn : Allow user to set NewParentCrn
func (options *UpdateEndpointOptions) SetNewParentCrn(newParentCrn string) *UpdateEndpointOptions {
	options.NewParentCrn = core.StringPtr(newParentCrn)
	return options
}

// SetNewServiceInstanceCrn : Allow user to set NewServiceInstanceCrn
func (options *UpdateEndpointOptions) SetNewServiceInstanceCrn(newServiceInstanceCrn string) *UpdateEndpointOptions {
	options.NewServiceInstanceCrn = core.StringPtr(newServiceInstanceCrn)
	return options
}

// SetNewName : Allow user to set NewName
func (options *UpdateEndpointOptions) SetNewName(newName string) *UpdateEndpointOptions {
	options.NewName = core.StringPtr(newName)
	return options
}

// SetNewRoutes : Allow user to set NewRoutes
func (options *UpdateEndpointOptions) SetNewRoutes(newRoutes []string) *UpdateEndpointOptions {
	options.NewRoutes = newRoutes
	return options
}

// SetNewManaged : Allow user to set NewManaged
func (options *UpdateEndpointOptions) SetNewManaged(newManaged bool) *UpdateEndpointOptions {
	options.NewManaged = core.BoolPtr(newManaged)
	return options
}

// SetNewMetadata : Allow user to set NewMetadata
func (options *UpdateEndpointOptions) SetNewMetadata(newMetadata interface{}) *UpdateEndpointOptions {
	options.NewMetadata = newMetadata
	return options
}

// SetNewOpenApiDoc : Allow user to set NewOpenApiDoc
func (options *UpdateEndpointOptions) SetNewOpenApiDoc(newOpenApiDoc interface{}) *UpdateEndpointOptions {
	options.NewOpenApiDoc = newOpenApiDoc
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEndpointOptions) SetHeaders(param map[string]string) *UpdateEndpointOptions {
	options.Headers = param
	return options
}

// UpdateSubscriptionOptions : The UpdateSubscription options.
type UpdateSubscriptionOptions struct {
	// Client Id.
	ID *string `json:"id" validate:"required"`

	// Artifact Id.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// User bearer token.
	Authorization *string `json:"authorization" validate:"required"`

	NewClientID *string `json:"new_client_id,omitempty"`

	NewClientSecret *string `json:"new_client_secret,omitempty"`

	NewArtifactID *string `json:"new_artifact_id,omitempty"`

	NewAccountID *string `json:"new_account_id,omitempty"`

	NewName *string `json:"new_name,omitempty"`

	NewType *string `json:"new_type,omitempty"`

	NewMetadata interface{} `json:"new_metadata,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSubscriptionOptions : Instantiate UpdateSubscriptionOptions
func (*ApiGatewayControllerApiV1) NewUpdateSubscriptionOptions(id string, artifactID string, authorization string) *UpdateSubscriptionOptions {
	return &UpdateSubscriptionOptions{
		ID:            core.StringPtr(id),
		ArtifactID:    core.StringPtr(artifactID),
		Authorization: core.StringPtr(authorization),
	}
}

// SetID : Allow user to set ID
func (options *UpdateSubscriptionOptions) SetID(id string) *UpdateSubscriptionOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetArtifactID : Allow user to set ArtifactID
func (options *UpdateSubscriptionOptions) SetArtifactID(artifactID string) *UpdateSubscriptionOptions {
	options.ArtifactID = core.StringPtr(artifactID)
	return options
}

// SetAuthorization : Allow user to set Authorization
func (options *UpdateSubscriptionOptions) SetAuthorization(authorization string) *UpdateSubscriptionOptions {
	options.Authorization = core.StringPtr(authorization)
	return options
}

// SetNewClientID : Allow user to set NewClientID
func (options *UpdateSubscriptionOptions) SetNewClientID(newClientID string) *UpdateSubscriptionOptions {
	options.NewClientID = core.StringPtr(newClientID)
	return options
}

// SetNewClientSecret : Allow user to set NewClientSecret
func (options *UpdateSubscriptionOptions) SetNewClientSecret(newClientSecret string) *UpdateSubscriptionOptions {
	options.NewClientSecret = core.StringPtr(newClientSecret)
	return options
}

// SetNewArtifactID : Allow user to set NewArtifactID
func (options *UpdateSubscriptionOptions) SetNewArtifactID(newArtifactID string) *UpdateSubscriptionOptions {
	options.NewArtifactID = core.StringPtr(newArtifactID)
	return options
}

// SetNewAccountID : Allow user to set NewAccountID
func (options *UpdateSubscriptionOptions) SetNewAccountID(newAccountID string) *UpdateSubscriptionOptions {
	options.NewAccountID = core.StringPtr(newAccountID)
	return options
}

// SetNewName : Allow user to set NewName
func (options *UpdateSubscriptionOptions) SetNewName(newName string) *UpdateSubscriptionOptions {
	options.NewName = core.StringPtr(newName)
	return options
}

// SetNewType : Allow user to set NewType
func (options *UpdateSubscriptionOptions) SetNewType(newType string) *UpdateSubscriptionOptions {
	options.NewType = core.StringPtr(newType)
	return options
}

// SetNewMetadata : Allow user to set NewMetadata
func (options *UpdateSubscriptionOptions) SetNewMetadata(newMetadata interface{}) *UpdateSubscriptionOptions {
	options.NewMetadata = newMetadata
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSubscriptionOptions) SetHeaders(param map[string]string) *UpdateSubscriptionOptions {
	options.Headers = param
	return options
}

// InlineResponse200 : Subscription artifact requested.
type InlineResponse200 struct {
	OpenApiDoc interface{} `json:"open_api_doc,omitempty"`

	ManagedURL *string `json:"managed_url,omitempty"`
}

// UnmarshalInlineResponse200 constructs an instance of InlineResponse200 from the specified map.
func UnmarshalInlineResponse200(m map[string]interface{}) (result *InlineResponse200, err error) {
	obj := new(InlineResponse200)
	obj.OpenApiDoc, err = core.UnmarshalObject(m, "open_api_doc")
	if err != nil {
		return
	}
	obj.ManagedURL, err = core.UnmarshalString(m, "managed_url")
	if err != nil {
		return
	}
	result = obj
	return
}

// UnmarshalInlineResponse200Slice unmarshals a slice of InlineResponse200 instances from the specified list of maps.
func UnmarshalInlineResponse200Slice(s []interface{}) (slice []InlineResponse200, err error) {
	for _, v := range s {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("slice element should be a map containing an instance of 'InlineResponse200'")
			return
		}
		obj, e := UnmarshalInlineResponse200(objMap)
		if e != nil {
			err = e
			return
		}
		slice = append(slice, *obj)
	}
	return
}

// UnmarshalInlineResponse200AsProperty unmarshals an instance of InlineResponse200 that is stored as a property
// within the specified map.
func UnmarshalInlineResponse200AsProperty(m map[string]interface{}, propertyName string) (result *InlineResponse200, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a map containing an instance of 'InlineResponse200'", propertyName)
			return
		}
		result, err = UnmarshalInlineResponse200(objMap)
	}
	return
}

// UnmarshalInlineResponse200SliceAsProperty unmarshals a slice of InlineResponse200 instances that are stored as a property
// within the specified map.
func UnmarshalInlineResponse200SliceAsProperty(m map[string]interface{}, propertyName string) (slice []InlineResponse200, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a slice of maps, each containing an instance of 'InlineResponse200'", propertyName)
			return
		}
		slice, err = UnmarshalInlineResponse200Slice(vSlice)
	}
	return
}

// V2DiscoveryConfig : V2DiscoveryConfig struct
type V2DiscoveryConfig struct {
	Headers interface{} `json:"headers,omitempty"`

	BridgeURL *string `json:"bridge_url" validate:"required"`
}

// UnmarshalV2DiscoveryConfig constructs an instance of V2DiscoveryConfig from the specified map.
func UnmarshalV2DiscoveryConfig(m map[string]interface{}) (result *V2DiscoveryConfig, err error) {
	obj := new(V2DiscoveryConfig)
	obj.Headers, err = core.UnmarshalObject(m, "headers")
	if err != nil {
		return
	}
	obj.BridgeURL, err = core.UnmarshalString(m, "bridge_url")
	if err != nil {
		return
	}
	result = obj
	return
}

// UnmarshalV2DiscoveryConfigSlice unmarshals a slice of V2DiscoveryConfig instances from the specified list of maps.
func UnmarshalV2DiscoveryConfigSlice(s []interface{}) (slice []V2DiscoveryConfig, err error) {
	for _, v := range s {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("slice element should be a map containing an instance of 'V2DiscoveryConfig'")
			return
		}
		obj, e := UnmarshalV2DiscoveryConfig(objMap)
		if e != nil {
			err = e
			return
		}
		slice = append(slice, *obj)
	}
	return
}

// UnmarshalV2DiscoveryConfigAsProperty unmarshals an instance of V2DiscoveryConfig that is stored as a property
// within the specified map.
func UnmarshalV2DiscoveryConfigAsProperty(m map[string]interface{}, propertyName string) (result *V2DiscoveryConfig, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a map containing an instance of 'V2DiscoveryConfig'", propertyName)
			return
		}
		result, err = UnmarshalV2DiscoveryConfig(objMap)
	}
	return
}

// UnmarshalV2DiscoveryConfigSliceAsProperty unmarshals a slice of V2DiscoveryConfig instances that are stored as a property
// within the specified map.
func UnmarshalV2DiscoveryConfigSliceAsProperty(m map[string]interface{}, propertyName string) (slice []V2DiscoveryConfig, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a slice of maps, each containing an instance of 'V2DiscoveryConfig'", propertyName)
			return
		}
		slice, err = UnmarshalV2DiscoveryConfigSlice(vSlice)
	}
	return
}

// V2Endpoint : V2Endpoint struct
type V2Endpoint struct {
	// The endpoint ID.
	ArtifactID *string `json:"artifact_id" validate:"required"`

	// The endpoint CRN.
	Crn *string `json:"crn" validate:"required"`

	// The API Gateway service instance CRN.
	ParentCrn *string `json:"parent_crn" validate:"required"`

	// The API Gateway service instance CRN.
	ServiceInstanceCrn *string `json:"service_instance_crn,omitempty"`

	// The account where the API Gateway service instance was provisioned.
	AccountID *string `json:"account_id,omitempty"`

	// Endpoint metadata.
	Metadata interface{} `json:"metadata,omitempty"`

	// The provider type of this endpoint.
	ProviderID *string `json:"provider_id,omitempty"`

	// THe endpoint's name.
	Name *string `json:"name,omitempty"`

	// Invokable endpoint routes.
	Routes []string `json:"routes,omitempty"`

	// Invoke your endpoint with this URL.
	ManagedURL *string `json:"managed_url,omitempty"`

	// Invoke your endpoint with this alias URL.
	AliasURL *string `json:"alias_url,omitempty"`

	// Is your endpoint shared?.
	Shared *bool `json:"shared,omitempty"`

	// Is your endpoint managed by the API Gateway service instance?.
	Managed *bool `json:"managed,omitempty"`

	// Policies enforced on the endpoint.
	Policies []map[string]interface{} `json:"policies,omitempty"`

	// THe OpenAPI doc representing the endpoint.
	OpenApiDoc map[string]interface{} `json:"open_api_doc,omitempty"`

	// The base path of the endpoint.
	BasePath *string `json:"base_path,omitempty"`
}

// UnmarshalV2Endpoint constructs an instance of V2Endpoint from the specified map.
func UnmarshalV2Endpoint(m map[string]interface{}) (result *V2Endpoint, err error) {
	obj := new(V2Endpoint)
	obj.ArtifactID, err = core.UnmarshalString(m, "artifact_id")
	if err != nil {
		return
	}
	obj.Crn, err = core.UnmarshalString(m, "crn")
	if err != nil {
		return
	}
	obj.ParentCrn, err = core.UnmarshalString(m, "parent_crn")
	if err != nil {
		return
	}
	obj.ServiceInstanceCrn, err = core.UnmarshalString(m, "service_instance_crn")
	if err != nil {
		return
	}
	obj.AccountID, err = core.UnmarshalString(m, "account_id")
	if err != nil {
		return
	}
	obj.Metadata, err = core.UnmarshalObject(m, "metadata")
	if err != nil {
		return
	}
	obj.ProviderID, err = core.UnmarshalString(m, "provider_id")
	if err != nil {
		return
	}
	obj.Name, err = core.UnmarshalString(m, "name")
	if err != nil {
		return
	}
	obj.Routes, err = core.UnmarshalStringSlice(m, "routes")
	if err != nil {
		return
	}
	obj.ManagedURL, err = core.UnmarshalString(m, "managed_url")
	if err != nil {
		return
	}
	obj.AliasURL, err = core.UnmarshalString(m, "alias_url")
	if err != nil {
		return
	}
	obj.Shared, err = core.UnmarshalBool(m, "shared")
	if err != nil {
		return
	}
	obj.Managed, err = core.UnmarshalBool(m, "managed")
	if err != nil {
		return
	}
	obj.Policies, err = core.UnmarshalObjectSlice(m, "policies")
	if err != nil {
		return
	}
	obj.OpenApiDoc, err = core.UnmarshalObject(m, "open_api_doc")
	if err != nil {
		return
	}
	obj.BasePath, err = core.UnmarshalString(m, "base_path")
	if err != nil {
		return
	}
	result = obj
	return
}

// UnmarshalV2EndpointSlice unmarshals a slice of V2Endpoint instances from the specified list of maps.
func UnmarshalV2EndpointSlice(s []interface{}) (slice []V2Endpoint, err error) {
	for _, v := range s {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("slice element should be a map containing an instance of 'V2Endpoint'")
			return
		}
		obj, e := UnmarshalV2Endpoint(objMap)
		if e != nil {
			err = e
			return
		}
		slice = append(slice, *obj)
	}
	return
}

// UnmarshalV2EndpointAsProperty unmarshals an instance of V2Endpoint that is stored as a property
// within the specified map.
func UnmarshalV2EndpointAsProperty(m map[string]interface{}, propertyName string) (result *V2Endpoint, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a map containing an instance of 'V2Endpoint'", propertyName)
			return
		}
		result, err = UnmarshalV2Endpoint(objMap)
	}
	return
}

// UnmarshalV2EndpointSliceAsProperty unmarshals a slice of V2Endpoint instances that are stored as a property
// within the specified map.
func UnmarshalV2EndpointSliceAsProperty(m map[string]interface{}, propertyName string) (slice []V2Endpoint, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a slice of maps, each containing an instance of 'V2Endpoint'", propertyName)
			return
		}
		slice, err = UnmarshalV2EndpointSlice(vSlice)
	}
	return
}

// V2EndpointSummary : V2EndpointSummary struct
type V2EndpointSummary struct {
	ID *string `json:"id,omitempty"`

	DisplayName *string `json:"display_name,omitempty"`

	Metadata interface{} `json:"metadata,omitempty"`

	Discoverable *bool `json:"discoverable,omitempty"`

	DiscoveryConfig *V2DiscoveryConfig `json:"discovery_config,omitempty"`

	Endpoints []V2Endpoint `json:"endpoints,omitempty"`
}

// UnmarshalV2EndpointSummary constructs an instance of V2EndpointSummary from the specified map.
func UnmarshalV2EndpointSummary(m map[string]interface{}) (result *V2EndpointSummary, err error) {
	obj := new(V2EndpointSummary)
	obj.ID, err = core.UnmarshalString(m, "id")
	if err != nil {
		return
	}
	obj.DisplayName, err = core.UnmarshalString(m, "display_name")
	if err != nil {
		return
	}
	obj.Metadata, err = core.UnmarshalObject(m, "metadata")
	if err != nil {
		return
	}
	obj.Discoverable, err = core.UnmarshalBool(m, "discoverable")
	if err != nil {
		return
	}
	obj.DiscoveryConfig, err = UnmarshalV2DiscoveryConfigAsProperty(m, "discovery_config")
	if err != nil {
		return
	}
	obj.Endpoints, err = UnmarshalV2EndpointSliceAsProperty(m, "endpoints")
	if err != nil {
		return
	}
	result = obj
	return
}

// UnmarshalV2EndpointSummarySlice unmarshals a slice of V2EndpointSummary instances from the specified list of maps.
func UnmarshalV2EndpointSummarySlice(s []interface{}) (slice []V2EndpointSummary, err error) {
	for _, v := range s {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("slice element should be a map containing an instance of 'V2EndpointSummary'")
			return
		}
		obj, e := UnmarshalV2EndpointSummary(objMap)
		if e != nil {
			err = e
			return
		}
		slice = append(slice, *obj)
	}
	return
}

// UnmarshalV2EndpointSummaryAsProperty unmarshals an instance of V2EndpointSummary that is stored as a property
// within the specified map.
func UnmarshalV2EndpointSummaryAsProperty(m map[string]interface{}, propertyName string) (result *V2EndpointSummary, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a map containing an instance of 'V2EndpointSummary'", propertyName)
			return
		}
		result, err = UnmarshalV2EndpointSummary(objMap)
	}
	return
}

// UnmarshalV2EndpointSummarySliceAsProperty unmarshals a slice of V2EndpointSummary instances that are stored as a property
// within the specified map.
func UnmarshalV2EndpointSummarySliceAsProperty(m map[string]interface{}, propertyName string) (slice []V2EndpointSummary, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a slice of maps, each containing an instance of 'V2EndpointSummary'", propertyName)
			return
		}
		slice, err = UnmarshalV2EndpointSummarySlice(vSlice)
	}
	return
}

// V2Subscription : V2Subscription struct
type V2Subscription struct {
	ClientID *string `json:"client_id" ,omitempty`

	SecretProvided *bool `json:"secret_provided,omitempty"`

	ArtifactID *string `json:"artifact_id" validate:"required"`

	AccountID *string `json:"account_id,omitempty"`

	Name *string `json:"name,omitempty"`

	Type *string `json:"type,omitempty"`

	Metadata interface{} `json:"metadata,omitempty"`
}

// UnmarshalV2Subscription constructs an instance of V2Subscription from the specified map.
func UnmarshalV2Subscription(m map[string]interface{}) (result *V2Subscription, err error) {
	obj := new(V2Subscription)
	obj.ClientID, err = core.UnmarshalString(m, "client_id")
	if err != nil {
		return
	}
	obj.SecretProvided, err = core.UnmarshalBool(m, "secret_provided")
	if err != nil {
		return
	}
	obj.ArtifactID, err = core.UnmarshalString(m, "artifact_id")
	if err != nil {
		return
	}
	obj.AccountID, err = core.UnmarshalString(m, "account_id")
	if err != nil {
		return
	}
	obj.Name, err = core.UnmarshalString(m, "name")
	if err != nil {
		return
	}
	obj.Type, err = core.UnmarshalString(m, "type")
	if err != nil {
		return
	}
	obj.Metadata, err = core.UnmarshalObject(m, "metadata")
	if err != nil {
		return
	}
	result = obj
	return
}

// UnmarshalV2SubscriptionSlice unmarshals a slice of V2Subscription instances from the specified list of maps.
func UnmarshalV2SubscriptionSlice(s []interface{}) (slice []V2Subscription, err error) {
	for _, v := range s {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("slice element should be a map containing an instance of 'V2Subscription'")
			return
		}
		obj, e := UnmarshalV2Subscription(objMap)
		if e != nil {
			err = e
			return
		}
		slice = append(slice, *obj)
	}
	return
}

// UnmarshalV2SubscriptionAsProperty unmarshals an instance of V2Subscription that is stored as a property
// within the specified map.
func UnmarshalV2SubscriptionAsProperty(m map[string]interface{}, propertyName string) (result *V2Subscription, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		objMap, ok := v.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a map containing an instance of 'V2Subscription'", propertyName)
			return
		}
		result, err = UnmarshalV2Subscription(objMap)
	}
	return
}

// UnmarshalV2SubscriptionSliceAsProperty unmarshals a slice of V2Subscription instances that are stored as a property
// within the specified map.
func UnmarshalV2SubscriptionSliceAsProperty(m map[string]interface{}, propertyName string) (slice []V2Subscription, err error) {
	v, foundIt := m[propertyName]
	if foundIt {
		vSlice, ok := v.([]interface{})
		if !ok {
			err = fmt.Errorf("map property '%s' should be a slice of maps, each containing an instance of 'V2Subscription'", propertyName)
			return
		}
		slice, err = UnmarshalV2SubscriptionSlice(vSlice)
	}
	return
}
