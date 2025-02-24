package catalogmanagementv1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
)

// Add a plan to an offering.
func (catalogManagement *CatalogManagementV1) AddPlan(catalogID string, offeringID string, plan *Plan, headers map[string]string) (result *Plan, response *core.DetailedResponse, err error) {
	result, response, err = catalogManagement.AddPlanWithContext(context.Background(), catalogID, offeringID, plan, headers)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AddPlanWithContext is an alternate form of the addPlan method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) AddPlanWithContext(ctx context.Context, catalogID string, offeringID string, plan *Plan, headers map[string]string) (result *Plan, response *core.DetailedResponse, err error) {
	pathParamsMap := map[string]string{
		"catalogID":  catalogID,
		"offeringID": offeringID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/catalogs/{catalogID}/offerings/{offeringID}/plans`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "addPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(plan)
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
	response, err = catalogManagement.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "add_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPlan)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// Set a plan as validated.
func (catalogManagement *CatalogManagementV1) SetValidatePlan(planID string, headers map[string]string) (response *core.DetailedResponse, err error) {
	response, err = catalogManagement.SetValidatePlanWithContext(context.Background(), planID, headers)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetValidatePlanWithContext is an alternate form of the setValidatePlan method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SetValidatePlanWithContext(ctx context.Context, planID string, headers map[string]string) (response *core.DetailedResponse, err error) {
	pathParamsMap := map[string]string{
		"planLocID": planID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/plans/{planLocID}/validate/true`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "setValidatePlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
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

	response, err = catalogManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_allow_publish_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// Set a plan as publish approved.
func (catalogManagement *CatalogManagementV1) SetAllowPublishPlan(planID string, headers map[string]string) (response *core.DetailedResponse, err error) {
	response, err = catalogManagement.SetAllowPublishPlanWithContext(context.Background(), planID, headers)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetAllowPublishPlanWithContext is an alternate form of the setAllowPublishPlan method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SetAllowPublishPlanWithContext(ctx context.Context, planID string, headers map[string]string) (response *core.DetailedResponse, err error) {
	pathParamsMap := map[string]string{
		"planLocID": planID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, `/plans/{planLocID}/publish/publish_approved/true`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "setAllowPublishPlan")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
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

	response, err = catalogManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_allow_publish_plan", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// Set allow publish offering.
func (catalogManagement *CatalogManagementV1) SetAllowPublishOffering(catalogID string, offeringID string, approvalType string, setting bool, headers map[string]string) (response *core.DetailedResponse, err error) {
	response, err = catalogManagement.SetAllowPublishOfferingWithContext(context.Background(), catalogID, offeringID, approvalType, setting, headers)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetAllowPublishOfferingWithContext is an alternate form of the setAllowPublishOffering method which supports a Context parameter
func (catalogManagement *CatalogManagementV1) SetAllowPublishOfferingWithContext(ctx context.Context, catalogID string, offeringID string, approvalType string, setting bool, headers map[string]string) (response *core.DetailedResponse, err error) {
	pathParamsMap := map[string]string{
		"catalogID":  catalogID,
		"offeringID": offeringID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = catalogManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(catalogManagement.Service.Options.URL, fmt.Sprintf("/catalogs/{catalogID}/offerings/{offeringID}/publish/%s/%v", approvalType, setting), pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("catalog_management", "V1", "setAllowPublishOffering")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
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

	response, err = catalogManagement.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_allow_publish_offering", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
