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
 * IBM OpenAPI SDK Code Generator Version: 3.98.0-8be2046a-20241205-162752
 */

// Package vmwarev1 : Operations and models for the VmwareV1 service
package vmwarev1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/vmware-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// VmwareV1 : IBM Cloud for VMware Cloud Foundation as a Service API
//
// API Version: 1.2.0
type VmwareV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.us-south.vmware.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "vmware"

const ParameterizedServiceURL = "https://api.{region}.vmware.cloud.ibm.com/v1"

var defaultUrlVariables = map[string]string{
	"region": "us-south",
}

// VmwareV1Options : Service options
type VmwareV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewVmwareV1UsingExternalConfig : constructs an instance of VmwareV1 with passed in options and external configuration.
func NewVmwareV1UsingExternalConfig(options *VmwareV1Options) (vmware *VmwareV1, err error) {
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

	vmware, err = NewVmwareV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = vmware.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = vmware.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewVmwareV1 : constructs an instance of VmwareV1 with passed in options.
func NewVmwareV1(options *VmwareV1Options) (service *VmwareV1, err error) {
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

	service = &VmwareV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "vmware" suitable for processing requests.
func (vmware *VmwareV1) Clone() *VmwareV1 {
	if core.IsNil(vmware) {
		return nil
	}
	clone := *vmware
	clone.Service = vmware.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (vmware *VmwareV1) SetServiceURL(url string) error {
	err := vmware.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (vmware *VmwareV1) GetServiceURL() string {
	return vmware.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (vmware *VmwareV1) SetDefaultHeaders(headers http.Header) {
	vmware.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (vmware *VmwareV1) SetEnableGzipCompression(enableGzip bool) {
	vmware.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (vmware *VmwareV1) GetEnableGzipCompression() bool {
	return vmware.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (vmware *VmwareV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	vmware.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (vmware *VmwareV1) DisableRetries() {
	vmware.Service.DisableRetries()
}

// CreateDirectorSites : Create a Cloud Director site instance
// Create an instance of a Cloud Director site with specified configurations. The Cloud Director site instance is the
// infrastructure and associated VMware software stack, which consists of VMware vCenter Server, VMware NSX-T, and
// VMware Cloud Director. VMware platform management and operations are performed with Cloud Director. The minimum
// initial order size is 2 hosts (2-Socket 32 Cores, 192 GB RAM) with 24 TB of 2.0 IOPS/GB storage.
func (vmware *VmwareV1) CreateDirectorSites(createDirectorSitesOptions *CreateDirectorSitesOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateDirectorSitesWithContext(context.Background(), createDirectorSitesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDirectorSitesWithContext is an alternate form of the CreateDirectorSites method which supports a Context parameter
func (vmware *VmwareV1) CreateDirectorSitesWithContext(ctx context.Context, createDirectorSitesOptions *CreateDirectorSitesOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDirectorSitesOptions, "createDirectorSitesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDirectorSitesOptions, "createDirectorSitesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDirectorSitesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateDirectorSites")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDirectorSitesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createDirectorSitesOptions.AcceptLanguage))
	}
	if createDirectorSitesOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*createDirectorSitesOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if createDirectorSitesOptions.Name != nil {
		body["name"] = createDirectorSitesOptions.Name
	}
	if createDirectorSitesOptions.Pvdcs != nil {
		body["pvdcs"] = createDirectorSitesOptions.Pvdcs
	}
	if createDirectorSitesOptions.ResourceGroup != nil {
		body["resource_group"] = createDirectorSitesOptions.ResourceGroup
	}
	if createDirectorSitesOptions.Services != nil {
		body["services"] = createDirectorSitesOptions.Services
	}
	if createDirectorSitesOptions.PrivateOnly != nil {
		body["private_only"] = createDirectorSitesOptions.PrivateOnly
	}
	if createDirectorSitesOptions.ConsoleConnectionType != nil {
		body["console_connection_type"] = createDirectorSitesOptions.ConsoleConnectionType
	}
	if createDirectorSitesOptions.IpAllowList != nil {
		body["ip_allow_list"] = createDirectorSitesOptions.IpAllowList
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_director_sites", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSite)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListDirectorSites : List Cloud Director site instances
// List all VMware Cloud Director site instances that the user can access in the cloud account.
func (vmware *VmwareV1) ListDirectorSites(listDirectorSitesOptions *ListDirectorSitesOptions) (result *DirectorSiteCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListDirectorSitesWithContext(context.Background(), listDirectorSitesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDirectorSitesWithContext is an alternate form of the ListDirectorSites method which supports a Context parameter
func (vmware *VmwareV1) ListDirectorSitesWithContext(ctx context.Context, listDirectorSitesOptions *ListDirectorSitesOptions) (result *DirectorSiteCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listDirectorSitesOptions, "listDirectorSitesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDirectorSitesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListDirectorSites")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDirectorSitesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listDirectorSitesOptions.AcceptLanguage))
	}
	if listDirectorSitesOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listDirectorSitesOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_director_sites", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSiteCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDirectorSite : Get a Cloud Director site instance
// Get a Cloud Director site instance by specifying the instance ID.
func (vmware *VmwareV1) GetDirectorSite(getDirectorSiteOptions *GetDirectorSiteOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	result, response, err = vmware.GetDirectorSiteWithContext(context.Background(), getDirectorSiteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDirectorSiteWithContext is an alternate form of the GetDirectorSite method which supports a Context parameter
func (vmware *VmwareV1) GetDirectorSiteWithContext(ctx context.Context, getDirectorSiteOptions *GetDirectorSiteOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDirectorSiteOptions, "getDirectorSiteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDirectorSiteOptions, "getDirectorSiteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getDirectorSiteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDirectorSiteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "GetDirectorSite")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDirectorSiteOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getDirectorSiteOptions.AcceptLanguage))
	}
	if getDirectorSiteOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*getDirectorSiteOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_director_site", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSite)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDirectorSite : Delete a Cloud Director site instance
// Delete a Cloud Director site instance by specifying the instance ID.
func (vmware *VmwareV1) DeleteDirectorSite(deleteDirectorSiteOptions *DeleteDirectorSiteOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	result, response, err = vmware.DeleteDirectorSiteWithContext(context.Background(), deleteDirectorSiteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDirectorSiteWithContext is an alternate form of the DeleteDirectorSite method which supports a Context parameter
func (vmware *VmwareV1) DeleteDirectorSiteWithContext(ctx context.Context, deleteDirectorSiteOptions *DeleteDirectorSiteOptions) (result *DirectorSite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDirectorSiteOptions, "deleteDirectorSiteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDirectorSiteOptions, "deleteDirectorSiteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteDirectorSiteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDirectorSiteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "DeleteDirectorSite")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDirectorSiteOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*deleteDirectorSiteOptions.AcceptLanguage))
	}
	if deleteDirectorSiteOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*deleteDirectorSiteOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_director_site", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSite)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// EnableVeeamOnPvdcsList : Enable or disable Veeam on a Cloud Director site
// Enable or disable Veeam on a Cloud Director site.
func (vmware *VmwareV1) EnableVeeamOnPvdcsList(enableVeeamOnPvdcsListOptions *EnableVeeamOnPvdcsListOptions) (result *ServiceEnabled, response *core.DetailedResponse, err error) {
	result, response, err = vmware.EnableVeeamOnPvdcsListWithContext(context.Background(), enableVeeamOnPvdcsListOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// EnableVeeamOnPvdcsListWithContext is an alternate form of the EnableVeeamOnPvdcsList method which supports a Context parameter
func (vmware *VmwareV1) EnableVeeamOnPvdcsListWithContext(ctx context.Context, enableVeeamOnPvdcsListOptions *EnableVeeamOnPvdcsListOptions) (result *ServiceEnabled, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(enableVeeamOnPvdcsListOptions, "enableVeeamOnPvdcsListOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(enableVeeamOnPvdcsListOptions, "enableVeeamOnPvdcsListOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *enableVeeamOnPvdcsListOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/action/enable_veeam`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range enableVeeamOnPvdcsListOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "EnableVeeamOnPvdcsList")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if enableVeeamOnPvdcsListOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*enableVeeamOnPvdcsListOptions.AcceptLanguage))
	}
	if enableVeeamOnPvdcsListOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*enableVeeamOnPvdcsListOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if enableVeeamOnPvdcsListOptions.Enable != nil {
		body["enable"] = enableVeeamOnPvdcsListOptions.Enable
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "enable_veeam_on_pvdcs_list", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceEnabled)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// EnableVcdaOnDataCenter : Enable or disable VCDA on a Cloud Director site
// Enable or disable VMware Cloud Director Availability (VCDA) on a Cloud Director site.
func (vmware *VmwareV1) EnableVcdaOnDataCenter(enableVcdaOnDataCenterOptions *EnableVcdaOnDataCenterOptions) (result *ServiceEnabled, response *core.DetailedResponse, err error) {
	result, response, err = vmware.EnableVcdaOnDataCenterWithContext(context.Background(), enableVcdaOnDataCenterOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// EnableVcdaOnDataCenterWithContext is an alternate form of the EnableVcdaOnDataCenter method which supports a Context parameter
func (vmware *VmwareV1) EnableVcdaOnDataCenterWithContext(ctx context.Context, enableVcdaOnDataCenterOptions *EnableVcdaOnDataCenterOptions) (result *ServiceEnabled, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(enableVcdaOnDataCenterOptions, "enableVcdaOnDataCenterOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(enableVcdaOnDataCenterOptions, "enableVcdaOnDataCenterOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *enableVcdaOnDataCenterOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/action/enable_vcda`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range enableVcdaOnDataCenterOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "EnableVcdaOnDataCenter")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if enableVcdaOnDataCenterOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*enableVcdaOnDataCenterOptions.AcceptLanguage))
	}
	if enableVcdaOnDataCenterOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*enableVcdaOnDataCenterOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if enableVcdaOnDataCenterOptions.Enable != nil {
		body["enable"] = enableVcdaOnDataCenterOptions.Enable
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "enable_vcda_on_data_center", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceEnabled)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDirectorSitesVcdaConnectionEndpoints : Create a VCDA connection
// Create a VMware Cloud Director Availability (VCDA) connection in the Cloud Director site identified by {site_id}.
func (vmware *VmwareV1) CreateDirectorSitesVcdaConnectionEndpoints(createDirectorSitesVcdaConnectionEndpointsOptions *CreateDirectorSitesVcdaConnectionEndpointsOptions) (result *VcdaConnection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateDirectorSitesVcdaConnectionEndpointsWithContext(context.Background(), createDirectorSitesVcdaConnectionEndpointsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDirectorSitesVcdaConnectionEndpointsWithContext is an alternate form of the CreateDirectorSitesVcdaConnectionEndpoints method which supports a Context parameter
func (vmware *VmwareV1) CreateDirectorSitesVcdaConnectionEndpointsWithContext(ctx context.Context, createDirectorSitesVcdaConnectionEndpointsOptions *CreateDirectorSitesVcdaConnectionEndpointsOptions) (result *VcdaConnection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDirectorSitesVcdaConnectionEndpointsOptions, "createDirectorSitesVcdaConnectionEndpointsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDirectorSitesVcdaConnectionEndpointsOptions, "createDirectorSitesVcdaConnectionEndpointsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *createDirectorSitesVcdaConnectionEndpointsOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/vcda/connection_endpoints`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDirectorSitesVcdaConnectionEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateDirectorSitesVcdaConnectionEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage))
	}
	if createDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*createDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if createDirectorSitesVcdaConnectionEndpointsOptions.Type != nil {
		body["type"] = createDirectorSitesVcdaConnectionEndpointsOptions.Type
	}
	if createDirectorSitesVcdaConnectionEndpointsOptions.DataCenterName != nil {
		body["data_center_name"] = createDirectorSitesVcdaConnectionEndpointsOptions.DataCenterName
	}
	if createDirectorSitesVcdaConnectionEndpointsOptions.AllowList != nil {
		body["allow_list"] = createDirectorSitesVcdaConnectionEndpointsOptions.AllowList
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_director_sites_vcda_connection_endpoints", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVcdaConnection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDirectorSitesVcdaConnectionEndpoints : Delete a VCDA connection
// Delete a VMware Cloud Director Availability (VCDA) connection in the Cloud Director site identified by {site_id} and
// {vcda_connections_id}.
func (vmware *VmwareV1) DeleteDirectorSitesVcdaConnectionEndpoints(deleteDirectorSitesVcdaConnectionEndpointsOptions *DeleteDirectorSitesVcdaConnectionEndpointsOptions) (result *VcdaConnection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.DeleteDirectorSitesVcdaConnectionEndpointsWithContext(context.Background(), deleteDirectorSitesVcdaConnectionEndpointsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDirectorSitesVcdaConnectionEndpointsWithContext is an alternate form of the DeleteDirectorSitesVcdaConnectionEndpoints method which supports a Context parameter
func (vmware *VmwareV1) DeleteDirectorSitesVcdaConnectionEndpointsWithContext(ctx context.Context, deleteDirectorSitesVcdaConnectionEndpointsOptions *DeleteDirectorSitesVcdaConnectionEndpointsOptions) (result *VcdaConnection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDirectorSitesVcdaConnectionEndpointsOptions, "deleteDirectorSitesVcdaConnectionEndpointsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDirectorSitesVcdaConnectionEndpointsOptions, "deleteDirectorSitesVcdaConnectionEndpointsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *deleteDirectorSitesVcdaConnectionEndpointsOptions.SiteID,
		"id": *deleteDirectorSitesVcdaConnectionEndpointsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/services/vcda/connection_endpoints/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDirectorSitesVcdaConnectionEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "DeleteDirectorSitesVcdaConnectionEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*deleteDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage))
	}
	if deleteDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*deleteDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_director_sites_vcda_connection_endpoints", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVcdaConnection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDirectorSitesVcdaConnectionEndpoints : Update VCDA connection allowlist
// Update the allowlist for a private connection to a specific VCDA instance.
func (vmware *VmwareV1) UpdateDirectorSitesVcdaConnectionEndpoints(updateDirectorSitesVcdaConnectionEndpointsOptions *UpdateDirectorSitesVcdaConnectionEndpointsOptions) (result *UpdatedVcdaConnection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.UpdateDirectorSitesVcdaConnectionEndpointsWithContext(context.Background(), updateDirectorSitesVcdaConnectionEndpointsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDirectorSitesVcdaConnectionEndpointsWithContext is an alternate form of the UpdateDirectorSitesVcdaConnectionEndpoints method which supports a Context parameter
func (vmware *VmwareV1) UpdateDirectorSitesVcdaConnectionEndpointsWithContext(ctx context.Context, updateDirectorSitesVcdaConnectionEndpointsOptions *UpdateDirectorSitesVcdaConnectionEndpointsOptions) (result *UpdatedVcdaConnection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDirectorSitesVcdaConnectionEndpointsOptions, "updateDirectorSitesVcdaConnectionEndpointsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDirectorSitesVcdaConnectionEndpointsOptions, "updateDirectorSitesVcdaConnectionEndpointsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *updateDirectorSitesVcdaConnectionEndpointsOptions.SiteID,
		"id": *updateDirectorSitesVcdaConnectionEndpointsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/services/vcda/connection_endpoints/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDirectorSitesVcdaConnectionEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "UpdateDirectorSitesVcdaConnectionEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*updateDirectorSitesVcdaConnectionEndpointsOptions.AcceptLanguage))
	}
	if updateDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*updateDirectorSitesVcdaConnectionEndpointsOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if updateDirectorSitesVcdaConnectionEndpointsOptions.AllowList != nil {
		body["allow_list"] = updateDirectorSitesVcdaConnectionEndpointsOptions.AllowList
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_director_sites_vcda_connection_endpoints", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdatedVcdaConnection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDirectorSitesVcdaC2cConnection : Create a VCDA cloud-to-cloud connection
// Create a VCDA cloud-to-cloud connection in the Cloud Director site identified by {site_id}.
func (vmware *VmwareV1) CreateDirectorSitesVcdaC2cConnection(createDirectorSitesVcdaC2cConnectionOptions *CreateDirectorSitesVcdaC2cConnectionOptions) (result *VcdaC2c, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateDirectorSitesVcdaC2cConnectionWithContext(context.Background(), createDirectorSitesVcdaC2cConnectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDirectorSitesVcdaC2cConnectionWithContext is an alternate form of the CreateDirectorSitesVcdaC2cConnection method which supports a Context parameter
func (vmware *VmwareV1) CreateDirectorSitesVcdaC2cConnectionWithContext(ctx context.Context, createDirectorSitesVcdaC2cConnectionOptions *CreateDirectorSitesVcdaC2cConnectionOptions) (result *VcdaC2c, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDirectorSitesVcdaC2cConnectionOptions, "createDirectorSitesVcdaC2cConnectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDirectorSitesVcdaC2cConnectionOptions, "createDirectorSitesVcdaC2cConnectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *createDirectorSitesVcdaC2cConnectionOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/services/vcda/c2c_connections`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDirectorSitesVcdaC2cConnectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateDirectorSitesVcdaC2cConnection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage))
	}
	if createDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*createDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if createDirectorSitesVcdaC2cConnectionOptions.LocalDataCenterName != nil {
		body["local_data_center_name"] = createDirectorSitesVcdaC2cConnectionOptions.LocalDataCenterName
	}
	if createDirectorSitesVcdaC2cConnectionOptions.LocalSiteName != nil {
		body["local_site_name"] = createDirectorSitesVcdaC2cConnectionOptions.LocalSiteName
	}
	if createDirectorSitesVcdaC2cConnectionOptions.PeerSiteName != nil {
		body["peer_site_name"] = createDirectorSitesVcdaC2cConnectionOptions.PeerSiteName
	}
	if createDirectorSitesVcdaC2cConnectionOptions.PeerRegion != nil {
		body["peer_region"] = createDirectorSitesVcdaC2cConnectionOptions.PeerRegion
	}
	if createDirectorSitesVcdaC2cConnectionOptions.Note != nil {
		body["note"] = createDirectorSitesVcdaC2cConnectionOptions.Note
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_director_sites_vcda_c2c_connection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVcdaC2c)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDirectorSitesVcdaC2cConnection : Delete a VCDA cloud-to-cloud connection
// Delete a VCDA cloud-to-cloud connection in the Cloud Director site identified by {site_id}.
func (vmware *VmwareV1) DeleteDirectorSitesVcdaC2cConnection(deleteDirectorSitesVcdaC2cConnectionOptions *DeleteDirectorSitesVcdaC2cConnectionOptions) (result *VcdaC2c, response *core.DetailedResponse, err error) {
	result, response, err = vmware.DeleteDirectorSitesVcdaC2cConnectionWithContext(context.Background(), deleteDirectorSitesVcdaC2cConnectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDirectorSitesVcdaC2cConnectionWithContext is an alternate form of the DeleteDirectorSitesVcdaC2cConnection method which supports a Context parameter
func (vmware *VmwareV1) DeleteDirectorSitesVcdaC2cConnectionWithContext(ctx context.Context, deleteDirectorSitesVcdaC2cConnectionOptions *DeleteDirectorSitesVcdaC2cConnectionOptions) (result *VcdaC2c, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDirectorSitesVcdaC2cConnectionOptions, "deleteDirectorSitesVcdaC2cConnectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDirectorSitesVcdaC2cConnectionOptions, "deleteDirectorSitesVcdaC2cConnectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *deleteDirectorSitesVcdaC2cConnectionOptions.SiteID,
		"id": *deleteDirectorSitesVcdaC2cConnectionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/services/vcda/c2c_connections/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDirectorSitesVcdaC2cConnectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "DeleteDirectorSitesVcdaC2cConnection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*deleteDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage))
	}
	if deleteDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*deleteDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_director_sites_vcda_c2c_connection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVcdaC2c)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDirectorSitesVcdaC2cConnection : Update note in the cloud-to-cloud connection
// Update the note in the VCDA cloud-to-cloud connection in the Cloud Director site identified by {site_id}.
func (vmware *VmwareV1) UpdateDirectorSitesVcdaC2cConnection(updateDirectorSitesVcdaC2cConnectionOptions *UpdateDirectorSitesVcdaC2cConnectionOptions) (result *UpdatedVcdaC2c, response *core.DetailedResponse, err error) {
	result, response, err = vmware.UpdateDirectorSitesVcdaC2cConnectionWithContext(context.Background(), updateDirectorSitesVcdaC2cConnectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDirectorSitesVcdaC2cConnectionWithContext is an alternate form of the UpdateDirectorSitesVcdaC2cConnection method which supports a Context parameter
func (vmware *VmwareV1) UpdateDirectorSitesVcdaC2cConnectionWithContext(ctx context.Context, updateDirectorSitesVcdaC2cConnectionOptions *UpdateDirectorSitesVcdaC2cConnectionOptions) (result *UpdatedVcdaC2c, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDirectorSitesVcdaC2cConnectionOptions, "updateDirectorSitesVcdaC2cConnectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDirectorSitesVcdaC2cConnectionOptions, "updateDirectorSitesVcdaC2cConnectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *updateDirectorSitesVcdaC2cConnectionOptions.SiteID,
		"id": *updateDirectorSitesVcdaC2cConnectionOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/services/vcda/c2c_connections/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDirectorSitesVcdaC2cConnectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "UpdateDirectorSitesVcdaC2cConnection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*updateDirectorSitesVcdaC2cConnectionOptions.AcceptLanguage))
	}
	if updateDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*updateDirectorSitesVcdaC2cConnectionOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if updateDirectorSitesVcdaC2cConnectionOptions.Note != nil {
		body["note"] = updateDirectorSitesVcdaC2cConnectionOptions.Note
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_director_sites_vcda_c2c_connection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdatedVcdaC2c)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetOidcConfiguration : Get an OIDC configuration
// Return the details of an OpenID Connect (OIDC) configuration on a Cloud Director site.
func (vmware *VmwareV1) GetOidcConfiguration(getOidcConfigurationOptions *GetOidcConfigurationOptions) (result *OIDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.GetOidcConfigurationWithContext(context.Background(), getOidcConfigurationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetOidcConfigurationWithContext is an alternate form of the GetOidcConfiguration method which supports a Context parameter
func (vmware *VmwareV1) GetOidcConfigurationWithContext(ctx context.Context, getOidcConfigurationOptions *GetOidcConfigurationOptions) (result *OIDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOidcConfigurationOptions, "getOidcConfigurationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getOidcConfigurationOptions, "getOidcConfigurationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *getOidcConfigurationOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/oidc_configuration`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getOidcConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "GetOidcConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getOidcConfigurationOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getOidcConfigurationOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_oidc_configuration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOIDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SetOidcConfiguration : Set an OIDC configuration
// Request to configure OpenID Connect (OIDC) on a Cloud Director site.
func (vmware *VmwareV1) SetOidcConfiguration(setOidcConfigurationOptions *SetOidcConfigurationOptions) (result *OIDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.SetOidcConfigurationWithContext(context.Background(), setOidcConfigurationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetOidcConfigurationWithContext is an alternate form of the SetOidcConfiguration method which supports a Context parameter
func (vmware *VmwareV1) SetOidcConfigurationWithContext(ctx context.Context, setOidcConfigurationOptions *SetOidcConfigurationOptions) (result *OIDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setOidcConfigurationOptions, "setOidcConfigurationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(setOidcConfigurationOptions, "setOidcConfigurationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *setOidcConfigurationOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/oidc_configuration`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range setOidcConfigurationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "SetOidcConfiguration")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if setOidcConfigurationOptions.ContentLength != nil {
		builder.AddHeader("Content-Length", fmt.Sprint(*setOidcConfigurationOptions.ContentLength))
	}
	if setOidcConfigurationOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*setOidcConfigurationOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_oidc_configuration", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOIDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListDirectorSitesPvdcs : List the resource pools in a Cloud Director site instance
// List the resource pools in a specified Cloud Director site.
func (vmware *VmwareV1) ListDirectorSitesPvdcs(listDirectorSitesPvdcsOptions *ListDirectorSitesPvdcsOptions) (result *PVDCCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListDirectorSitesPvdcsWithContext(context.Background(), listDirectorSitesPvdcsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDirectorSitesPvdcsWithContext is an alternate form of the ListDirectorSitesPvdcs method which supports a Context parameter
func (vmware *VmwareV1) ListDirectorSitesPvdcsWithContext(ctx context.Context, listDirectorSitesPvdcsOptions *ListDirectorSitesPvdcsOptions) (result *PVDCCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDirectorSitesPvdcsOptions, "listDirectorSitesPvdcsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listDirectorSitesPvdcsOptions, "listDirectorSitesPvdcsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *listDirectorSitesPvdcsOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDirectorSitesPvdcsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListDirectorSitesPvdcs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDirectorSitesPvdcsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listDirectorSitesPvdcsOptions.AcceptLanguage))
	}
	if listDirectorSitesPvdcsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listDirectorSitesPvdcsOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_director_sites_pvdcs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPVDCCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDirectorSitesPvdcs : Create a resource pool instance in a specified Cloud Director site
// Create an instance of a resource pool with specified configurations. The Cloud Director site instance is the
// infrastructure and associated VMware software stack, which consists of VMware vCenter Server, VMware NSX-T, and
// VMware Cloud Director. VMware platform management and operations are performed with Cloud Director. The minimum
// initial order size is 2 hosts (2-Socket 32 Cores, 192 GB RAM) with 24 TB of 2.0 IOPS/GB storage.
func (vmware *VmwareV1) CreateDirectorSitesPvdcs(createDirectorSitesPvdcsOptions *CreateDirectorSitesPvdcsOptions) (result *PVDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateDirectorSitesPvdcsWithContext(context.Background(), createDirectorSitesPvdcsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDirectorSitesPvdcsWithContext is an alternate form of the CreateDirectorSitesPvdcs method which supports a Context parameter
func (vmware *VmwareV1) CreateDirectorSitesPvdcsWithContext(ctx context.Context, createDirectorSitesPvdcsOptions *CreateDirectorSitesPvdcsOptions) (result *PVDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDirectorSitesPvdcsOptions, "createDirectorSitesPvdcsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDirectorSitesPvdcsOptions, "createDirectorSitesPvdcsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *createDirectorSitesPvdcsOptions.SiteID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDirectorSitesPvdcsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateDirectorSitesPvdcs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDirectorSitesPvdcsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createDirectorSitesPvdcsOptions.AcceptLanguage))
	}
	if createDirectorSitesPvdcsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*createDirectorSitesPvdcsOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if createDirectorSitesPvdcsOptions.Name != nil {
		body["name"] = createDirectorSitesPvdcsOptions.Name
	}
	if createDirectorSitesPvdcsOptions.DataCenterName != nil {
		body["data_center_name"] = createDirectorSitesPvdcsOptions.DataCenterName
	}
	if createDirectorSitesPvdcsOptions.Clusters != nil {
		body["clusters"] = createDirectorSitesPvdcsOptions.Clusters
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_director_sites_pvdcs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDirectorSitesPvdcs : Get the specified resource pool in a Cloud Director site instance
// Get the specified resource pools in a specified Cloud Director site.
func (vmware *VmwareV1) GetDirectorSitesPvdcs(getDirectorSitesPvdcsOptions *GetDirectorSitesPvdcsOptions) (result *PVDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.GetDirectorSitesPvdcsWithContext(context.Background(), getDirectorSitesPvdcsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDirectorSitesPvdcsWithContext is an alternate form of the GetDirectorSitesPvdcs method which supports a Context parameter
func (vmware *VmwareV1) GetDirectorSitesPvdcsWithContext(ctx context.Context, getDirectorSitesPvdcsOptions *GetDirectorSitesPvdcsOptions) (result *PVDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDirectorSitesPvdcsOptions, "getDirectorSitesPvdcsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDirectorSitesPvdcsOptions, "getDirectorSitesPvdcsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *getDirectorSitesPvdcsOptions.SiteID,
		"id": *getDirectorSitesPvdcsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDirectorSitesPvdcsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "GetDirectorSitesPvdcs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDirectorSitesPvdcsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getDirectorSitesPvdcsOptions.AcceptLanguage))
	}
	if getDirectorSitesPvdcsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*getDirectorSitesPvdcsOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_director_sites_pvdcs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListDirectorSitesPvdcsClusters : List clusters
// List all VMware clusters of a Cloud Director site instance by specifying the instance ID.
func (vmware *VmwareV1) ListDirectorSitesPvdcsClusters(listDirectorSitesPvdcsClustersOptions *ListDirectorSitesPvdcsClustersOptions) (result *ClusterCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListDirectorSitesPvdcsClustersWithContext(context.Background(), listDirectorSitesPvdcsClustersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDirectorSitesPvdcsClustersWithContext is an alternate form of the ListDirectorSitesPvdcsClusters method which supports a Context parameter
func (vmware *VmwareV1) ListDirectorSitesPvdcsClustersWithContext(ctx context.Context, listDirectorSitesPvdcsClustersOptions *ListDirectorSitesPvdcsClustersOptions) (result *ClusterCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDirectorSitesPvdcsClustersOptions, "listDirectorSitesPvdcsClustersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listDirectorSitesPvdcsClustersOptions, "listDirectorSitesPvdcsClustersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *listDirectorSitesPvdcsClustersOptions.SiteID,
		"pvdc_id": *listDirectorSitesPvdcsClustersOptions.PvdcID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{pvdc_id}/clusters`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDirectorSitesPvdcsClustersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListDirectorSitesPvdcsClusters")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDirectorSitesPvdcsClustersOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listDirectorSitesPvdcsClustersOptions.AcceptLanguage))
	}
	if listDirectorSitesPvdcsClustersOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listDirectorSitesPvdcsClustersOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_director_sites_pvdcs_clusters", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalClusterCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDirectorSitesPvdcsClusters : Create a cluster
// Create a VMware cluster under a specified resource pool in a Cloud Director site instance.
func (vmware *VmwareV1) CreateDirectorSitesPvdcsClusters(createDirectorSitesPvdcsClustersOptions *CreateDirectorSitesPvdcsClustersOptions) (result *Cluster, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateDirectorSitesPvdcsClustersWithContext(context.Background(), createDirectorSitesPvdcsClustersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDirectorSitesPvdcsClustersWithContext is an alternate form of the CreateDirectorSitesPvdcsClusters method which supports a Context parameter
func (vmware *VmwareV1) CreateDirectorSitesPvdcsClustersWithContext(ctx context.Context, createDirectorSitesPvdcsClustersOptions *CreateDirectorSitesPvdcsClustersOptions) (result *Cluster, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDirectorSitesPvdcsClustersOptions, "createDirectorSitesPvdcsClustersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDirectorSitesPvdcsClustersOptions, "createDirectorSitesPvdcsClustersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *createDirectorSitesPvdcsClustersOptions.SiteID,
		"pvdc_id": *createDirectorSitesPvdcsClustersOptions.PvdcID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{pvdc_id}/clusters`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDirectorSitesPvdcsClustersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateDirectorSitesPvdcsClusters")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createDirectorSitesPvdcsClustersOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createDirectorSitesPvdcsClustersOptions.AcceptLanguage))
	}
	if createDirectorSitesPvdcsClustersOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*createDirectorSitesPvdcsClustersOptions.XGlobalTransactionID))
	}

	body := make(map[string]interface{})
	if createDirectorSitesPvdcsClustersOptions.Name != nil {
		body["name"] = createDirectorSitesPvdcsClustersOptions.Name
	}
	if createDirectorSitesPvdcsClustersOptions.HostCount != nil {
		body["host_count"] = createDirectorSitesPvdcsClustersOptions.HostCount
	}
	if createDirectorSitesPvdcsClustersOptions.HostProfile != nil {
		body["host_profile"] = createDirectorSitesPvdcsClustersOptions.HostProfile
	}
	if createDirectorSitesPvdcsClustersOptions.FileShares != nil {
		body["file_shares"] = createDirectorSitesPvdcsClustersOptions.FileShares
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_director_sites_pvdcs_clusters", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCluster)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDirectorInstancesPvdcsCluster : Get a cluster
// Get a specific VMware cluster from the resource pool in a Cloud Director site instance.
func (vmware *VmwareV1) GetDirectorInstancesPvdcsCluster(getDirectorInstancesPvdcsClusterOptions *GetDirectorInstancesPvdcsClusterOptions) (result *Cluster, response *core.DetailedResponse, err error) {
	result, response, err = vmware.GetDirectorInstancesPvdcsClusterWithContext(context.Background(), getDirectorInstancesPvdcsClusterOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDirectorInstancesPvdcsClusterWithContext is an alternate form of the GetDirectorInstancesPvdcsCluster method which supports a Context parameter
func (vmware *VmwareV1) GetDirectorInstancesPvdcsClusterWithContext(ctx context.Context, getDirectorInstancesPvdcsClusterOptions *GetDirectorInstancesPvdcsClusterOptions) (result *Cluster, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDirectorInstancesPvdcsClusterOptions, "getDirectorInstancesPvdcsClusterOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDirectorInstancesPvdcsClusterOptions, "getDirectorInstancesPvdcsClusterOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *getDirectorInstancesPvdcsClusterOptions.SiteID,
		"id": *getDirectorInstancesPvdcsClusterOptions.ID,
		"pvdc_id": *getDirectorInstancesPvdcsClusterOptions.PvdcID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{pvdc_id}/clusters/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDirectorInstancesPvdcsClusterOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "GetDirectorInstancesPvdcsCluster")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getDirectorInstancesPvdcsClusterOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getDirectorInstancesPvdcsClusterOptions.AcceptLanguage))
	}
	if getDirectorInstancesPvdcsClusterOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*getDirectorInstancesPvdcsClusterOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_director_instances_pvdcs_cluster", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCluster)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDirectorSitesPvdcsCluster : Delete a cluster
// Delete a cluster from a resource pool in a Cloud Director site instance by specifying the instance ID.
func (vmware *VmwareV1) DeleteDirectorSitesPvdcsCluster(deleteDirectorSitesPvdcsClusterOptions *DeleteDirectorSitesPvdcsClusterOptions) (result *ClusterSummary, response *core.DetailedResponse, err error) {
	result, response, err = vmware.DeleteDirectorSitesPvdcsClusterWithContext(context.Background(), deleteDirectorSitesPvdcsClusterOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDirectorSitesPvdcsClusterWithContext is an alternate form of the DeleteDirectorSitesPvdcsCluster method which supports a Context parameter
func (vmware *VmwareV1) DeleteDirectorSitesPvdcsClusterWithContext(ctx context.Context, deleteDirectorSitesPvdcsClusterOptions *DeleteDirectorSitesPvdcsClusterOptions) (result *ClusterSummary, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDirectorSitesPvdcsClusterOptions, "deleteDirectorSitesPvdcsClusterOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDirectorSitesPvdcsClusterOptions, "deleteDirectorSitesPvdcsClusterOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *deleteDirectorSitesPvdcsClusterOptions.SiteID,
		"id": *deleteDirectorSitesPvdcsClusterOptions.ID,
		"pvdc_id": *deleteDirectorSitesPvdcsClusterOptions.PvdcID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{pvdc_id}/clusters/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDirectorSitesPvdcsClusterOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "DeleteDirectorSitesPvdcsCluster")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteDirectorSitesPvdcsClusterOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*deleteDirectorSitesPvdcsClusterOptions.AcceptLanguage))
	}
	if deleteDirectorSitesPvdcsClusterOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*deleteDirectorSitesPvdcsClusterOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_director_sites_pvdcs_cluster", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalClusterSummary)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateDirectorSitesPvdcsCluster : Update a cluster
// Update the number of hosts or file storage shares of a specific cluster in a specific Cloud Director site instance.
// VMware clusters must have between [2-25] hosts.
func (vmware *VmwareV1) UpdateDirectorSitesPvdcsCluster(updateDirectorSitesPvdcsClusterOptions *UpdateDirectorSitesPvdcsClusterOptions) (result *UpdateCluster, response *core.DetailedResponse, err error) {
	result, response, err = vmware.UpdateDirectorSitesPvdcsClusterWithContext(context.Background(), updateDirectorSitesPvdcsClusterOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDirectorSitesPvdcsClusterWithContext is an alternate form of the UpdateDirectorSitesPvdcsCluster method which supports a Context parameter
func (vmware *VmwareV1) UpdateDirectorSitesPvdcsClusterWithContext(ctx context.Context, updateDirectorSitesPvdcsClusterOptions *UpdateDirectorSitesPvdcsClusterOptions) (result *UpdateCluster, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDirectorSitesPvdcsClusterOptions, "updateDirectorSitesPvdcsClusterOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDirectorSitesPvdcsClusterOptions, "updateDirectorSitesPvdcsClusterOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"site_id": *updateDirectorSitesPvdcsClusterOptions.SiteID,
		"id": *updateDirectorSitesPvdcsClusterOptions.ID,
		"pvdc_id": *updateDirectorSitesPvdcsClusterOptions.PvdcID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_sites/{site_id}/pvdcs/{pvdc_id}/clusters/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDirectorSitesPvdcsClusterOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "UpdateDirectorSitesPvdcsCluster")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateDirectorSitesPvdcsClusterOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*updateDirectorSitesPvdcsClusterOptions.AcceptLanguage))
	}
	if updateDirectorSitesPvdcsClusterOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*updateDirectorSitesPvdcsClusterOptions.XGlobalTransactionID))
	}

	_, err = builder.SetBodyContentJSON(updateDirectorSitesPvdcsClusterOptions.Body)
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_director_sites_pvdcs_cluster", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdateCluster)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListDirectorSiteRegions : List regions
// List all IBM Cloud regions enabled for users to create a new Cloud Director site instance.
func (vmware *VmwareV1) ListDirectorSiteRegions(listDirectorSiteRegionsOptions *ListDirectorSiteRegionsOptions) (result *DirectorSiteRegionCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListDirectorSiteRegionsWithContext(context.Background(), listDirectorSiteRegionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDirectorSiteRegionsWithContext is an alternate form of the ListDirectorSiteRegions method which supports a Context parameter
func (vmware *VmwareV1) ListDirectorSiteRegionsWithContext(ctx context.Context, listDirectorSiteRegionsOptions *ListDirectorSiteRegionsOptions) (result *DirectorSiteRegionCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listDirectorSiteRegionsOptions, "listDirectorSiteRegionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_site_regions`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDirectorSiteRegionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListDirectorSiteRegions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDirectorSiteRegionsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listDirectorSiteRegionsOptions.AcceptLanguage))
	}
	if listDirectorSiteRegionsOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listDirectorSiteRegionsOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_director_site_regions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSiteRegionCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListMultitenantDirectorSites : Get all multitenant Cloud Director sites
// Retrieve a collection of multitenant Cloud Director sites.
func (vmware *VmwareV1) ListMultitenantDirectorSites(listMultitenantDirectorSitesOptions *ListMultitenantDirectorSitesOptions) (result *MultitenantDirectorSiteCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListMultitenantDirectorSitesWithContext(context.Background(), listMultitenantDirectorSitesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListMultitenantDirectorSitesWithContext is an alternate form of the ListMultitenantDirectorSites method which supports a Context parameter
func (vmware *VmwareV1) ListMultitenantDirectorSitesWithContext(ctx context.Context, listMultitenantDirectorSitesOptions *ListMultitenantDirectorSitesOptions) (result *MultitenantDirectorSiteCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listMultitenantDirectorSitesOptions, "listMultitenantDirectorSitesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/multitenant_director_sites`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listMultitenantDirectorSitesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListMultitenantDirectorSites")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listMultitenantDirectorSitesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listMultitenantDirectorSitesOptions.AcceptLanguage))
	}
	if listMultitenantDirectorSitesOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listMultitenantDirectorSitesOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_multitenant_director_sites", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMultitenantDirectorSiteCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListDirectorSiteHostProfiles : List host profiles
// List available host profiles that can be used when you create a Cloud Director site instance. IBM Cloud offers
// several different host types. Typically, the host type is selected based on the properties of the workload to be run
// in the VMware cluster.
func (vmware *VmwareV1) ListDirectorSiteHostProfiles(listDirectorSiteHostProfilesOptions *ListDirectorSiteHostProfilesOptions) (result *DirectorSiteHostProfileCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListDirectorSiteHostProfilesWithContext(context.Background(), listDirectorSiteHostProfilesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDirectorSiteHostProfilesWithContext is an alternate form of the ListDirectorSiteHostProfiles method which supports a Context parameter
func (vmware *VmwareV1) ListDirectorSiteHostProfilesWithContext(ctx context.Context, listDirectorSiteHostProfilesOptions *ListDirectorSiteHostProfilesOptions) (result *DirectorSiteHostProfileCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listDirectorSiteHostProfilesOptions, "listDirectorSiteHostProfilesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/director_site_host_profiles`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDirectorSiteHostProfilesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListDirectorSiteHostProfiles")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listDirectorSiteHostProfilesOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listDirectorSiteHostProfilesOptions.AcceptLanguage))
	}
	if listDirectorSiteHostProfilesOptions.XGlobalTransactionID != nil {
		builder.AddHeader("X-Global-Transaction-ID", fmt.Sprint(*listDirectorSiteHostProfilesOptions.XGlobalTransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_director_site_host_profiles", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDirectorSiteHostProfileCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListVdcs : List virtual data centers
// List all virtual data centers (VDCs) that user has access to in the cloud account.
func (vmware *VmwareV1) ListVdcs(listVdcsOptions *ListVdcsOptions) (result *VDCCollection, response *core.DetailedResponse, err error) {
	result, response, err = vmware.ListVdcsWithContext(context.Background(), listVdcsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListVdcsWithContext is an alternate form of the ListVdcs method which supports a Context parameter
func (vmware *VmwareV1) ListVdcsWithContext(ctx context.Context, listVdcsOptions *ListVdcsOptions) (result *VDCCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listVdcsOptions, "listVdcsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listVdcsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "ListVdcs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listVdcsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*listVdcsOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_vdcs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVDCCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateVdc : Create a virtual data center
// Create a virtual data center (VDC) with specified configurations.
func (vmware *VmwareV1) CreateVdc(createVdcOptions *CreateVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.CreateVdcWithContext(context.Background(), createVdcOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateVdcWithContext is an alternate form of the CreateVdc method which supports a Context parameter
func (vmware *VmwareV1) CreateVdcWithContext(ctx context.Context, createVdcOptions *CreateVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createVdcOptions, "createVdcOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createVdcOptions, "createVdcOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createVdcOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "CreateVdc")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createVdcOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*createVdcOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createVdcOptions.Name != nil {
		body["name"] = createVdcOptions.Name
	}
	if createVdcOptions.DirectorSite != nil {
		body["director_site"] = createVdcOptions.DirectorSite
	}
	if createVdcOptions.Edge != nil {
		body["edge"] = createVdcOptions.Edge
	}
	if createVdcOptions.FastProvisioningEnabled != nil {
		body["fast_provisioning_enabled"] = createVdcOptions.FastProvisioningEnabled
	}
	if createVdcOptions.ResourceGroup != nil {
		body["resource_group"] = createVdcOptions.ResourceGroup
	}
	if createVdcOptions.Cpu != nil {
		body["cpu"] = createVdcOptions.Cpu
	}
	if createVdcOptions.Ram != nil {
		body["ram"] = createVdcOptions.Ram
	}
	if createVdcOptions.RhelByol != nil {
		body["rhel_byol"] = createVdcOptions.RhelByol
	}
	if createVdcOptions.WindowsByol != nil {
		body["windows_byol"] = createVdcOptions.WindowsByol
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_vdc", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetVdc : Get a virtual data center
// Get details about a virtual data center (VDC) by specifying the VDC ID.
func (vmware *VmwareV1) GetVdc(getVdcOptions *GetVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.GetVdcWithContext(context.Background(), getVdcOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetVdcWithContext is an alternate form of the GetVdc method which supports a Context parameter
func (vmware *VmwareV1) GetVdcWithContext(ctx context.Context, getVdcOptions *GetVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getVdcOptions, "getVdcOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getVdcOptions, "getVdcOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getVdcOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getVdcOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "GetVdc")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getVdcOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*getVdcOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_vdc", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteVdc : Delete a virtual data center
// Delete a virtual data center (VDC) by specifying the VDC ID.
func (vmware *VmwareV1) DeleteVdc(deleteVdcOptions *DeleteVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.DeleteVdcWithContext(context.Background(), deleteVdcOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteVdcWithContext is an alternate form of the DeleteVdc method which supports a Context parameter
func (vmware *VmwareV1) DeleteVdcWithContext(ctx context.Context, deleteVdcOptions *DeleteVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteVdcOptions, "deleteVdcOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteVdcOptions, "deleteVdcOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteVdcOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteVdcOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "DeleteVdc")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteVdcOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*deleteVdcOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_vdc", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateVdc : Update a virtual data center
// Update a virtual data center with the specified ID.
func (vmware *VmwareV1) UpdateVdc(updateVdcOptions *UpdateVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	result, response, err = vmware.UpdateVdcWithContext(context.Background(), updateVdcOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateVdcWithContext is an alternate form of the UpdateVdc method which supports a Context parameter
func (vmware *VmwareV1) UpdateVdcWithContext(ctx context.Context, updateVdcOptions *UpdateVdcOptions) (result *VDC, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateVdcOptions, "updateVdcOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateVdcOptions, "updateVdcOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateVdcOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateVdcOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "UpdateVdc")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateVdcOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*updateVdcOptions.AcceptLanguage))
	}

	_, err = builder.SetBodyContentJSON(updateVdcOptions.VDCPatch)
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_vdc", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVDC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// AddTransitGatewayConnections : Add IBM Transit Gateway connections to edge
// Add IBM Transit Gateway connections to an edge and virtual data center.
func (vmware *VmwareV1) AddTransitGatewayConnections(addTransitGatewayConnectionsOptions *AddTransitGatewayConnectionsOptions) (result *TransitGateway, response *core.DetailedResponse, err error) {
	result, response, err = vmware.AddTransitGatewayConnectionsWithContext(context.Background(), addTransitGatewayConnectionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// AddTransitGatewayConnectionsWithContext is an alternate form of the AddTransitGatewayConnections method which supports a Context parameter
func (vmware *VmwareV1) AddTransitGatewayConnectionsWithContext(ctx context.Context, addTransitGatewayConnectionsOptions *AddTransitGatewayConnectionsOptions) (result *TransitGateway, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(addTransitGatewayConnectionsOptions, "addTransitGatewayConnectionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(addTransitGatewayConnectionsOptions, "addTransitGatewayConnectionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"vdc_id": *addTransitGatewayConnectionsOptions.VdcID,
		"edge_id": *addTransitGatewayConnectionsOptions.EdgeID,
		"id": *addTransitGatewayConnectionsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs/{vdc_id}/edges/{edge_id}/transit_gateways/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range addTransitGatewayConnectionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "AddTransitGatewayConnections")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if addTransitGatewayConnectionsOptions.ContentLength != nil {
		builder.AddHeader("Content-Length", fmt.Sprint(*addTransitGatewayConnectionsOptions.ContentLength))
	}
	if addTransitGatewayConnectionsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*addTransitGatewayConnectionsOptions.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if addTransitGatewayConnectionsOptions.Region != nil {
		body["region"] = addTransitGatewayConnectionsOptions.Region
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
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "add_transit_gateway_connections", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTransitGateway)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// RemoveTransitGatewayConnections : Remove IBM Transit Gateway connections from edge
// Remove IBM Transit Gateway connections from an edge and virtual data center.
func (vmware *VmwareV1) RemoveTransitGatewayConnections(removeTransitGatewayConnectionsOptions *RemoveTransitGatewayConnectionsOptions) (result *TransitGateway, response *core.DetailedResponse, err error) {
	result, response, err = vmware.RemoveTransitGatewayConnectionsWithContext(context.Background(), removeTransitGatewayConnectionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// RemoveTransitGatewayConnectionsWithContext is an alternate form of the RemoveTransitGatewayConnections method which supports a Context parameter
func (vmware *VmwareV1) RemoveTransitGatewayConnectionsWithContext(ctx context.Context, removeTransitGatewayConnectionsOptions *RemoveTransitGatewayConnectionsOptions) (result *TransitGateway, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(removeTransitGatewayConnectionsOptions, "removeTransitGatewayConnectionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(removeTransitGatewayConnectionsOptions, "removeTransitGatewayConnectionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"vdc_id": *removeTransitGatewayConnectionsOptions.VdcID,
		"edge_id": *removeTransitGatewayConnectionsOptions.EdgeID,
		"id": *removeTransitGatewayConnectionsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = vmware.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(vmware.Service.Options.URL, `/vdcs/{vdc_id}/edges/{edge_id}/transit_gateways/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range removeTransitGatewayConnectionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("vmware", "V1", "RemoveTransitGatewayConnections")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if removeTransitGatewayConnectionsOptions.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*removeTransitGatewayConnectionsOptions.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = vmware.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "remove_transit_gateway_connections", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTransitGateway)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.2.0")
}

// AddTransitGatewayConnectionsOptions : The AddTransitGatewayConnections options.
type AddTransitGatewayConnectionsOptions struct {
	// A unique ID for a virtual data center.
	VdcID *string `json:"vdc_id" validate:"required,ne="`

	// A unique ID for an edge.
	EdgeID *string `json:"edge_id" validate:"required,ne="`

	// A unique ID for an IBM Transit Gateway.
	ID *string `json:"id" validate:"required,ne="`

	// Size of the message body in bytes.
	ContentLength *int64 `json:"Content-Length" validate:"required"`

	// The region where the IBM Transit Gateway is deployed.
	Region *string `json:"region,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewAddTransitGatewayConnectionsOptions : Instantiate AddTransitGatewayConnectionsOptions
func (*VmwareV1) NewAddTransitGatewayConnectionsOptions(vdcID string, edgeID string, id string, contentLength int64) *AddTransitGatewayConnectionsOptions {
	return &AddTransitGatewayConnectionsOptions{
		VdcID: core.StringPtr(vdcID),
		EdgeID: core.StringPtr(edgeID),
		ID: core.StringPtr(id),
		ContentLength: core.Int64Ptr(contentLength),
	}
}

// SetVdcID : Allow user to set VdcID
func (_options *AddTransitGatewayConnectionsOptions) SetVdcID(vdcID string) *AddTransitGatewayConnectionsOptions {
	_options.VdcID = core.StringPtr(vdcID)
	return _options
}

// SetEdgeID : Allow user to set EdgeID
func (_options *AddTransitGatewayConnectionsOptions) SetEdgeID(edgeID string) *AddTransitGatewayConnectionsOptions {
	_options.EdgeID = core.StringPtr(edgeID)
	return _options
}

// SetID : Allow user to set ID
func (_options *AddTransitGatewayConnectionsOptions) SetID(id string) *AddTransitGatewayConnectionsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetContentLength : Allow user to set ContentLength
func (_options *AddTransitGatewayConnectionsOptions) SetContentLength(contentLength int64) *AddTransitGatewayConnectionsOptions {
	_options.ContentLength = core.Int64Ptr(contentLength)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *AddTransitGatewayConnectionsOptions) SetRegion(region string) *AddTransitGatewayConnectionsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *AddTransitGatewayConnectionsOptions) SetAcceptLanguage(acceptLanguage string) *AddTransitGatewayConnectionsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *AddTransitGatewayConnectionsOptions) SetHeaders(param map[string]string) *AddTransitGatewayConnectionsOptions {
	options.Headers = param
	return options
}

// Cluster : A cluster resource.
type Cluster struct {
	// The cluster ID.
	ID *string `json:"id" validate:"required"`

	// The cluster name.
	Name *string `json:"name" validate:"required"`

	// The hyperlink of the cluster resource.
	Href *string `json:"href" validate:"required"`

	// The time that the cluster is ordered.
	OrderedAt *strfmt.DateTime `json:"ordered_at" validate:"required"`

	// The time that the cluster is provisioned and available to use.
	ProvisionedAt *strfmt.DateTime `json:"provisioned_at,omitempty"`

	// The number of hosts in the cluster.
	HostCount *int64 `json:"host_count" validate:"required"`

	// The status of the Cloud Director site cluster.
	Status *string `json:"status" validate:"required"`

	// The location of deployed cluster.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// Back link to associated Cloud Director site resource.
	DirectorSite *DirectorSiteReference `json:"director_site" validate:"required"`

	// The name of the host profile.
	HostProfile *string `json:"host_profile" validate:"required"`

	// The storage type of the cluster.
	StorageType *string `json:"storage_type" validate:"required"`

	// The billing plan for the cluster.
	BillingPlan *string `json:"billing_plan" validate:"required"`

	// Chosen storage policies and their sizes.
	FileShares *FileShares `json:"file_shares" validate:"required"`
}

// Constants associated with the Cluster.StorageType property.
// The storage type of the cluster.
const (
	Cluster_StorageType_Nfs = "nfs"
)

// Constants associated with the Cluster.BillingPlan property.
// The billing plan for the cluster.
const (
	Cluster_BillingPlan_Monthly = "monthly"
)

// UnmarshalCluster unmarshals an instance of Cluster from the specified map of raw messages.
func UnmarshalCluster(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Cluster)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_at", &obj.OrderedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "ordered_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioned_at", &obj.ProvisionedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioned_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_count", &obj.HostCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "director_site", &obj.DirectorSite, UnmarshalDirectorSiteReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_site-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_profile", &obj.HostProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "storage_type", &obj.StorageType)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billing_plan", &obj.BillingPlan)
	if err != nil {
		err = core.SDKErrorf(err, "", "billing_plan-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "file_shares", &obj.FileShares, UnmarshalFileShares)
	if err != nil {
		err = core.SDKErrorf(err, "", "file_shares-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ClusterCollection : Return all clusters instances.
type ClusterCollection struct {
	// List of cluster objects.
	Clusters []Cluster `json:"clusters" validate:"required"`
}

// UnmarshalClusterCollection unmarshals an instance of ClusterCollection from the specified map of raw messages.
func UnmarshalClusterCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClusterCollection)
	err = core.UnmarshalModel(m, "clusters", &obj.Clusters, UnmarshalCluster)
	if err != nil {
		err = core.SDKErrorf(err, "", "clusters-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ClusterPatch : The cluster patch. Currently, specifying both file_shares and host_count in one call is not supported.
type ClusterPatch struct {
	// Chosen storage policies and their sizes.
	FileShares *FileSharesPrototype `json:"file_shares,omitempty"`

	// Number of hosts to add to or remove from the cluster.
	HostCount *int64 `json:"host_count,omitempty"`
}

// UnmarshalClusterPatch unmarshals an instance of ClusterPatch from the specified map of raw messages.
func UnmarshalClusterPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClusterPatch)
	err = core.UnmarshalModel(m, "file_shares", &obj.FileShares, UnmarshalFileSharesPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "file_shares-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_count", &obj.HostCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_count-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ClusterPatch
func (clusterPatch *ClusterPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(clusterPatch.FileShares) {
		_patch["file_shares"] = clusterPatch.FileShares.asPatch()
	}
	if !core.IsNil(clusterPatch.HostCount) {
		_patch["host_count"] = clusterPatch.HostCount
	}

	return
}

// ClusterPrototype : VMware Cluster order information. Clusters form VMware workload availability boundaries.
type ClusterPrototype struct {
	// Name of the VMware cluster. Cluster names must be unique per Cloud Director site instance. Cluster names cannot be
	// changed after creation.
	Name *string `json:"name" validate:"required"`

	// Number of hosts in the VMware cluster.
	HostCount *int64 `json:"host_count" validate:"required"`

	// The host type. IBM Cloud offers several different host types. Typically, the host type is selected based on the
	// properties of the workload to be run in the VMware cluster.
	HostProfile *string `json:"host_profile" validate:"required"`

	// Chosen storage policies and their sizes.
	FileShares *FileSharesPrototype `json:"file_shares" validate:"required"`
}

// NewClusterPrototype : Instantiate ClusterPrototype (Generic Model Constructor)
func (*VmwareV1) NewClusterPrototype(name string, hostCount int64, hostProfile string, fileShares *FileSharesPrototype) (_model *ClusterPrototype, err error) {
	_model = &ClusterPrototype{
		Name: core.StringPtr(name),
		HostCount: core.Int64Ptr(hostCount),
		HostProfile: core.StringPtr(hostProfile),
		FileShares: fileShares,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalClusterPrototype unmarshals an instance of ClusterPrototype from the specified map of raw messages.
func UnmarshalClusterPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClusterPrototype)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_count", &obj.HostCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_profile", &obj.HostProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "file_shares", &obj.FileShares, UnmarshalFileSharesPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "file_shares-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ClusterSummary : VMware Cluster basic information.
type ClusterSummary struct {
	// The cluster name.
	Name *string `json:"name" validate:"required"`

	// Number of hosts in the VMware cluster.
	HostCount *int64 `json:"host_count" validate:"required"`

	// The host type. IBM Cloud offers several different host types. Typically, the host type is selected based on the
	// properties of the workload to be run in the VMware cluster.
	HostProfile *string `json:"host_profile" validate:"required"`

	// The cluster ID.
	ID *string `json:"id" validate:"required"`

	// The location of the deployed cluster.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// The status of the cluster.
	Status *string `json:"status" validate:"required"`

	// The hyperlink of the cluster resource.
	Href *string `json:"href" validate:"required"`

	// Chosen storage policies and their sizes.
	FileShares *FileShares `json:"file_shares" validate:"required"`
}

// UnmarshalClusterSummary unmarshals an instance of ClusterSummary from the specified map of raw messages.
func UnmarshalClusterSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClusterSummary)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_count", &obj.HostCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_profile", &obj.HostProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "file_shares", &obj.FileShares, UnmarshalFileShares)
	if err != nil {
		err = core.SDKErrorf(err, "", "file_shares-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateDirectorSitesOptions : The CreateDirectorSites options.
type CreateDirectorSitesOptions struct {
	// Name of the Cloud Director site instance. Use a name that is unique to your region and meaningful. Names cannot be
	// changed after initial creation.
	Name *string `json:"name" validate:"required"`

	// List of VMware resource pools to deploy on the instance.
	Pvdcs []PVDCPrototype `json:"pvdcs" validate:"required"`

	// The resource group to associate with the resource instance.
	// If not specified, the default resource group in the account is used.
	ResourceGroup *ResourceGroupIdentity `json:"resource_group,omitempty"`

	// List of services to deploy on the instance.
	Services []ServiceIdentity `json:"services,omitempty"`

	// Indicates whether the site is private only.
	PrivateOnly *bool `json:"private_only,omitempty"`

	// Type of console connection.
	ConsoleConnectionType *string `json:"console_connection_type,omitempty"`

	// List of allowed IP addresses.
	IpAllowList []string `json:"ip_allow_list,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateDirectorSitesOptions.ConsoleConnectionType property.
// Type of console connection.
const (
	CreateDirectorSitesOptions_ConsoleConnectionType_Private = "private"
	CreateDirectorSitesOptions_ConsoleConnectionType_Public = "public"
)

// NewCreateDirectorSitesOptions : Instantiate CreateDirectorSitesOptions
func (*VmwareV1) NewCreateDirectorSitesOptions(name string, pvdcs []PVDCPrototype) *CreateDirectorSitesOptions {
	return &CreateDirectorSitesOptions{
		Name: core.StringPtr(name),
		Pvdcs: pvdcs,
	}
}

// SetName : Allow user to set Name
func (_options *CreateDirectorSitesOptions) SetName(name string) *CreateDirectorSitesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPvdcs : Allow user to set Pvdcs
func (_options *CreateDirectorSitesOptions) SetPvdcs(pvdcs []PVDCPrototype) *CreateDirectorSitesOptions {
	_options.Pvdcs = pvdcs
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateDirectorSitesOptions) SetResourceGroup(resourceGroup *ResourceGroupIdentity) *CreateDirectorSitesOptions {
	_options.ResourceGroup = resourceGroup
	return _options
}

// SetServices : Allow user to set Services
func (_options *CreateDirectorSitesOptions) SetServices(services []ServiceIdentity) *CreateDirectorSitesOptions {
	_options.Services = services
	return _options
}

// SetPrivateOnly : Allow user to set PrivateOnly
func (_options *CreateDirectorSitesOptions) SetPrivateOnly(privateOnly bool) *CreateDirectorSitesOptions {
	_options.PrivateOnly = core.BoolPtr(privateOnly)
	return _options
}

// SetConsoleConnectionType : Allow user to set ConsoleConnectionType
func (_options *CreateDirectorSitesOptions) SetConsoleConnectionType(consoleConnectionType string) *CreateDirectorSitesOptions {
	_options.ConsoleConnectionType = core.StringPtr(consoleConnectionType)
	return _options
}

// SetIpAllowList : Allow user to set IpAllowList
func (_options *CreateDirectorSitesOptions) SetIpAllowList(ipAllowList []string) *CreateDirectorSitesOptions {
	_options.IpAllowList = ipAllowList
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateDirectorSitesOptions) SetAcceptLanguage(acceptLanguage string) *CreateDirectorSitesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *CreateDirectorSitesOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *CreateDirectorSitesOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDirectorSitesOptions) SetHeaders(param map[string]string) *CreateDirectorSitesOptions {
	options.Headers = param
	return options
}

// CreateDirectorSitesPvdcsClustersOptions : The CreateDirectorSitesPvdcsClusters options.
type CreateDirectorSitesPvdcsClustersOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	PvdcID *string `json:"pvdc_id" validate:"required,ne="`

	// Name of the VMware cluster. Cluster names must be unique per Cloud Director site instance. Cluster names cannot be
	// changed after creation.
	Name *string `json:"name" validate:"required"`

	// Number of hosts in the VMware cluster.
	HostCount *int64 `json:"host_count" validate:"required"`

	// The host type. IBM Cloud offers several different host types. Typically, the host type is selected based on the
	// properties of the workload to be run in the VMware cluster.
	HostProfile *string `json:"host_profile" validate:"required"`

	// Chosen storage policies and their sizes.
	FileShares *FileSharesPrototype `json:"file_shares" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateDirectorSitesPvdcsClustersOptions : Instantiate CreateDirectorSitesPvdcsClustersOptions
func (*VmwareV1) NewCreateDirectorSitesPvdcsClustersOptions(siteID string, pvdcID string, name string, hostCount int64, hostProfile string, fileShares *FileSharesPrototype) *CreateDirectorSitesPvdcsClustersOptions {
	return &CreateDirectorSitesPvdcsClustersOptions{
		SiteID: core.StringPtr(siteID),
		PvdcID: core.StringPtr(pvdcID),
		Name: core.StringPtr(name),
		HostCount: core.Int64Ptr(hostCount),
		HostProfile: core.StringPtr(hostProfile),
		FileShares: fileShares,
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetSiteID(siteID string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetPvdcID : Allow user to set PvdcID
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetPvdcID(pvdcID string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.PvdcID = core.StringPtr(pvdcID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetName(name string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHostCount : Allow user to set HostCount
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetHostCount(hostCount int64) *CreateDirectorSitesPvdcsClustersOptions {
	_options.HostCount = core.Int64Ptr(hostCount)
	return _options
}

// SetHostProfile : Allow user to set HostProfile
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetHostProfile(hostProfile string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.HostProfile = core.StringPtr(hostProfile)
	return _options
}

// SetFileShares : Allow user to set FileShares
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetFileShares(fileShares *FileSharesPrototype) *CreateDirectorSitesPvdcsClustersOptions {
	_options.FileShares = fileShares
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetAcceptLanguage(acceptLanguage string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *CreateDirectorSitesPvdcsClustersOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *CreateDirectorSitesPvdcsClustersOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDirectorSitesPvdcsClustersOptions) SetHeaders(param map[string]string) *CreateDirectorSitesPvdcsClustersOptions {
	options.Headers = param
	return options
}

// CreateDirectorSitesPvdcsOptions : The CreateDirectorSitesPvdcs options.
type CreateDirectorSitesPvdcsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Name of the resource pool. Resource pool names must be unique per Cloud Director site instance and they cannot be
	// changed after creation.
	Name *string `json:"name" validate:"required"`

	// Data center location to deploy the cluster. See `GET /director_site_regions` for supported data center locations.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// List of VMware clusters to deploy on the instance. Clusters form VMware workload availability boundaries.
	Clusters []ClusterPrototype `json:"clusters" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateDirectorSitesPvdcsOptions : Instantiate CreateDirectorSitesPvdcsOptions
func (*VmwareV1) NewCreateDirectorSitesPvdcsOptions(siteID string, name string, dataCenterName string, clusters []ClusterPrototype) *CreateDirectorSitesPvdcsOptions {
	return &CreateDirectorSitesPvdcsOptions{
		SiteID: core.StringPtr(siteID),
		Name: core.StringPtr(name),
		DataCenterName: core.StringPtr(dataCenterName),
		Clusters: clusters,
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *CreateDirectorSitesPvdcsOptions) SetSiteID(siteID string) *CreateDirectorSitesPvdcsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateDirectorSitesPvdcsOptions) SetName(name string) *CreateDirectorSitesPvdcsOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDataCenterName : Allow user to set DataCenterName
func (_options *CreateDirectorSitesPvdcsOptions) SetDataCenterName(dataCenterName string) *CreateDirectorSitesPvdcsOptions {
	_options.DataCenterName = core.StringPtr(dataCenterName)
	return _options
}

// SetClusters : Allow user to set Clusters
func (_options *CreateDirectorSitesPvdcsOptions) SetClusters(clusters []ClusterPrototype) *CreateDirectorSitesPvdcsOptions {
	_options.Clusters = clusters
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateDirectorSitesPvdcsOptions) SetAcceptLanguage(acceptLanguage string) *CreateDirectorSitesPvdcsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *CreateDirectorSitesPvdcsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *CreateDirectorSitesPvdcsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDirectorSitesPvdcsOptions) SetHeaders(param map[string]string) *CreateDirectorSitesPvdcsOptions {
	options.Headers = param
	return options
}

// CreateDirectorSitesVcdaC2cConnectionOptions : The CreateDirectorSitesVcdaC2cConnection options.
type CreateDirectorSitesVcdaC2cConnectionOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Local data center name.
	LocalDataCenterName *string `json:"local_data_center_name" validate:"required"`

	// Local site name.
	LocalSiteName *string `json:"local_site_name" validate:"required"`

	// Peer site name.
	PeerSiteName *string `json:"peer_site_name" validate:"required"`

	// Peer region.
	PeerRegion *string `json:"peer_region" validate:"required"`

	// Note.
	Note *string `json:"note,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateDirectorSitesVcdaC2cConnectionOptions : Instantiate CreateDirectorSitesVcdaC2cConnectionOptions
func (*VmwareV1) NewCreateDirectorSitesVcdaC2cConnectionOptions(siteID string, localDataCenterName string, localSiteName string, peerSiteName string, peerRegion string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	return &CreateDirectorSitesVcdaC2cConnectionOptions{
		SiteID: core.StringPtr(siteID),
		LocalDataCenterName: core.StringPtr(localDataCenterName),
		LocalSiteName: core.StringPtr(localSiteName),
		PeerSiteName: core.StringPtr(peerSiteName),
		PeerRegion: core.StringPtr(peerRegion),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetSiteID(siteID string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetLocalDataCenterName : Allow user to set LocalDataCenterName
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetLocalDataCenterName(localDataCenterName string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.LocalDataCenterName = core.StringPtr(localDataCenterName)
	return _options
}

// SetLocalSiteName : Allow user to set LocalSiteName
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetLocalSiteName(localSiteName string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.LocalSiteName = core.StringPtr(localSiteName)
	return _options
}

// SetPeerSiteName : Allow user to set PeerSiteName
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetPeerSiteName(peerSiteName string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.PeerSiteName = core.StringPtr(peerSiteName)
	return _options
}

// SetPeerRegion : Allow user to set PeerRegion
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetPeerRegion(peerRegion string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.PeerRegion = core.StringPtr(peerRegion)
	return _options
}

// SetNote : Allow user to set Note
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetNote(note string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.Note = core.StringPtr(note)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetAcceptLanguage(acceptLanguage string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *CreateDirectorSitesVcdaC2cConnectionOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDirectorSitesVcdaC2cConnectionOptions) SetHeaders(param map[string]string) *CreateDirectorSitesVcdaC2cConnectionOptions {
	options.Headers = param
	return options
}

// CreateDirectorSitesVcdaConnectionEndpointsOptions : The CreateDirectorSitesVcdaConnectionEndpoints options.
type CreateDirectorSitesVcdaConnectionEndpointsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Connection type.
	Type *string `json:"type" validate:"required"`

	// Where to deploy the cluster.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// List of IP addresses allowed in the public connection.
	AllowList []string `json:"allow_list,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateDirectorSitesVcdaConnectionEndpointsOptions.Type property.
// Connection type.
const (
	CreateDirectorSitesVcdaConnectionEndpointsOptions_Type_Private = "private"
	CreateDirectorSitesVcdaConnectionEndpointsOptions_Type_Public = "public"
)

// NewCreateDirectorSitesVcdaConnectionEndpointsOptions : Instantiate CreateDirectorSitesVcdaConnectionEndpointsOptions
func (*VmwareV1) NewCreateDirectorSitesVcdaConnectionEndpointsOptions(siteID string, typeVar string, dataCenterName string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	return &CreateDirectorSitesVcdaConnectionEndpointsOptions{
		SiteID: core.StringPtr(siteID),
		Type: core.StringPtr(typeVar),
		DataCenterName: core.StringPtr(dataCenterName),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetSiteID(siteID string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetType(typeVar string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetDataCenterName : Allow user to set DataCenterName
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetDataCenterName(dataCenterName string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.DataCenterName = core.StringPtr(dataCenterName)
	return _options
}

// SetAllowList : Allow user to set AllowList
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetAllowList(allowList []string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.AllowList = allowList
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetAcceptLanguage(acceptLanguage string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDirectorSitesVcdaConnectionEndpointsOptions) SetHeaders(param map[string]string) *CreateDirectorSitesVcdaConnectionEndpointsOptions {
	options.Headers = param
	return options
}

// CreateVdcOptions : The CreateVdc options.
type CreateVdcOptions struct {
	// A human readable ID for the virtual data center (VDC). Use a name that is unique to your region.
	Name *string `json:"name" validate:"required"`

	// The Cloud Director site in which to deploy the virtual data center (VDC).
	DirectorSite *VDCDirectorSitePrototype `json:"director_site" validate:"required"`

	// The networking edge to be deployed on the virtual data center (VDC).
	Edge *VDCEdgePrototype `json:"edge,omitempty"`

	// Indicates whether to enable or not fast provisioning.
	FastProvisioningEnabled *bool `json:"fast_provisioning_enabled,omitempty"`

	// The resource group to associate with the resource instance.
	// If not specified, the default resource group in the account is used.
	ResourceGroup *ResourceGroupIdentity `json:"resource_group,omitempty"`

	// The vCPU usage limit on the virtual data center (VDC). Supported for VDCs deployed on a multitenant Cloud Director
	// site. This property is required when the resource pool type is reserved.
	Cpu *int64 `json:"cpu,omitempty"`

	// The RAM usage limit on the virtual data center (VDC) in GB (1024^3 bytes). Supported for VDCs deployed on a
	// multitenant Cloud Director site. This property is required when the resource pool type is reserved.
	Ram *int64 `json:"ram,omitempty"`

	// Indicates if the RHEL VMs will be using the license from IBM or the customer will use their own license (BYOL).
	RhelByol *bool `json:"rhel_byol,omitempty"`

	// Indicates if the Microsoft Windows VMs will be using the license from IBM or the customer will use their own license
	// (BYOL).
	WindowsByol *bool `json:"windows_byol,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateVdcOptions : Instantiate CreateVdcOptions
func (*VmwareV1) NewCreateVdcOptions(name string, directorSite *VDCDirectorSitePrototype) *CreateVdcOptions {
	return &CreateVdcOptions{
		Name: core.StringPtr(name),
		DirectorSite: directorSite,
	}
}

// SetName : Allow user to set Name
func (_options *CreateVdcOptions) SetName(name string) *CreateVdcOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDirectorSite : Allow user to set DirectorSite
func (_options *CreateVdcOptions) SetDirectorSite(directorSite *VDCDirectorSitePrototype) *CreateVdcOptions {
	_options.DirectorSite = directorSite
	return _options
}

// SetEdge : Allow user to set Edge
func (_options *CreateVdcOptions) SetEdge(edge *VDCEdgePrototype) *CreateVdcOptions {
	_options.Edge = edge
	return _options
}

// SetFastProvisioningEnabled : Allow user to set FastProvisioningEnabled
func (_options *CreateVdcOptions) SetFastProvisioningEnabled(fastProvisioningEnabled bool) *CreateVdcOptions {
	_options.FastProvisioningEnabled = core.BoolPtr(fastProvisioningEnabled)
	return _options
}

// SetResourceGroup : Allow user to set ResourceGroup
func (_options *CreateVdcOptions) SetResourceGroup(resourceGroup *ResourceGroupIdentity) *CreateVdcOptions {
	_options.ResourceGroup = resourceGroup
	return _options
}

// SetCpu : Allow user to set Cpu
func (_options *CreateVdcOptions) SetCpu(cpu int64) *CreateVdcOptions {
	_options.Cpu = core.Int64Ptr(cpu)
	return _options
}

// SetRam : Allow user to set Ram
func (_options *CreateVdcOptions) SetRam(ram int64) *CreateVdcOptions {
	_options.Ram = core.Int64Ptr(ram)
	return _options
}

// SetRhelByol : Allow user to set RhelByol
func (_options *CreateVdcOptions) SetRhelByol(rhelByol bool) *CreateVdcOptions {
	_options.RhelByol = core.BoolPtr(rhelByol)
	return _options
}

// SetWindowsByol : Allow user to set WindowsByol
func (_options *CreateVdcOptions) SetWindowsByol(windowsByol bool) *CreateVdcOptions {
	_options.WindowsByol = core.BoolPtr(windowsByol)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *CreateVdcOptions) SetAcceptLanguage(acceptLanguage string) *CreateVdcOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateVdcOptions) SetHeaders(param map[string]string) *CreateVdcOptions {
	options.Headers = param
	return options
}

// DataCenter : Details of the data center.
type DataCenter struct {
	// The display name of the data center.
	DisplayName *string `json:"display_name" validate:"required"`

	// The name of the data center.
	Name *string `json:"name" validate:"required"`

	// The speed available per data center.
	UplinkSpeed *string `json:"uplink_speed" validate:"required"`
}

// UnmarshalDataCenter unmarshals an instance of DataCenter from the specified map of raw messages.
func UnmarshalDataCenter(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DataCenter)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "uplink_speed", &obj.UplinkSpeed)
	if err != nil {
		err = core.SDKErrorf(err, "", "uplink_speed-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteDirectorSiteOptions : The DeleteDirectorSite options.
type DeleteDirectorSiteOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDirectorSiteOptions : Instantiate DeleteDirectorSiteOptions
func (*VmwareV1) NewDeleteDirectorSiteOptions(id string) *DeleteDirectorSiteOptions {
	return &DeleteDirectorSiteOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteDirectorSiteOptions) SetID(id string) *DeleteDirectorSiteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *DeleteDirectorSiteOptions) SetAcceptLanguage(acceptLanguage string) *DeleteDirectorSiteOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *DeleteDirectorSiteOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *DeleteDirectorSiteOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDirectorSiteOptions) SetHeaders(param map[string]string) *DeleteDirectorSiteOptions {
	options.Headers = param
	return options
}

// DeleteDirectorSitesPvdcsClusterOptions : The DeleteDirectorSitesPvdcsCluster options.
type DeleteDirectorSitesPvdcsClusterOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// The cluster to query.
	ID *string `json:"id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	PvdcID *string `json:"pvdc_id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDirectorSitesPvdcsClusterOptions : Instantiate DeleteDirectorSitesPvdcsClusterOptions
func (*VmwareV1) NewDeleteDirectorSitesPvdcsClusterOptions(siteID string, id string, pvdcID string) *DeleteDirectorSitesPvdcsClusterOptions {
	return &DeleteDirectorSitesPvdcsClusterOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
		PvdcID: core.StringPtr(pvdcID),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *DeleteDirectorSitesPvdcsClusterOptions) SetSiteID(siteID string) *DeleteDirectorSitesPvdcsClusterOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteDirectorSitesPvdcsClusterOptions) SetID(id string) *DeleteDirectorSitesPvdcsClusterOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetPvdcID : Allow user to set PvdcID
func (_options *DeleteDirectorSitesPvdcsClusterOptions) SetPvdcID(pvdcID string) *DeleteDirectorSitesPvdcsClusterOptions {
	_options.PvdcID = core.StringPtr(pvdcID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *DeleteDirectorSitesPvdcsClusterOptions) SetAcceptLanguage(acceptLanguage string) *DeleteDirectorSitesPvdcsClusterOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *DeleteDirectorSitesPvdcsClusterOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *DeleteDirectorSitesPvdcsClusterOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDirectorSitesPvdcsClusterOptions) SetHeaders(param map[string]string) *DeleteDirectorSitesPvdcsClusterOptions {
	options.Headers = param
	return options
}

// DeleteDirectorSitesVcdaC2cConnectionOptions : The DeleteDirectorSitesVcdaC2cConnection options.
type DeleteDirectorSitesVcdaC2cConnectionOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the cloud-to-cloud connections in the relationship Cloud Director site.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDirectorSitesVcdaC2cConnectionOptions : Instantiate DeleteDirectorSitesVcdaC2cConnectionOptions
func (*VmwareV1) NewDeleteDirectorSitesVcdaC2cConnectionOptions(siteID string, id string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	return &DeleteDirectorSitesVcdaC2cConnectionOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *DeleteDirectorSitesVcdaC2cConnectionOptions) SetSiteID(siteID string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteDirectorSitesVcdaC2cConnectionOptions) SetID(id string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *DeleteDirectorSitesVcdaC2cConnectionOptions) SetAcceptLanguage(acceptLanguage string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *DeleteDirectorSitesVcdaC2cConnectionOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDirectorSitesVcdaC2cConnectionOptions) SetHeaders(param map[string]string) *DeleteDirectorSitesVcdaC2cConnectionOptions {
	options.Headers = param
	return options
}

// DeleteDirectorSitesVcdaConnectionEndpointsOptions : The DeleteDirectorSitesVcdaConnectionEndpoints options.
type DeleteDirectorSitesVcdaConnectionEndpointsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the VCDA connections in the relationship Cloud Director site.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDirectorSitesVcdaConnectionEndpointsOptions : Instantiate DeleteDirectorSitesVcdaConnectionEndpointsOptions
func (*VmwareV1) NewDeleteDirectorSitesVcdaConnectionEndpointsOptions(siteID string, id string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	return &DeleteDirectorSitesVcdaConnectionEndpointsOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *DeleteDirectorSitesVcdaConnectionEndpointsOptions) SetSiteID(siteID string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteDirectorSitesVcdaConnectionEndpointsOptions) SetID(id string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *DeleteDirectorSitesVcdaConnectionEndpointsOptions) SetAcceptLanguage(acceptLanguage string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *DeleteDirectorSitesVcdaConnectionEndpointsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDirectorSitesVcdaConnectionEndpointsOptions) SetHeaders(param map[string]string) *DeleteDirectorSitesVcdaConnectionEndpointsOptions {
	options.Headers = param
	return options
}

// DeleteVdcOptions : The DeleteVdc options.
type DeleteVdcOptions struct {
	// A unique ID for a specified virtual data center.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteVdcOptions : Instantiate DeleteVdcOptions
func (*VmwareV1) NewDeleteVdcOptions(id string) *DeleteVdcOptions {
	return &DeleteVdcOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteVdcOptions) SetID(id string) *DeleteVdcOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *DeleteVdcOptions) SetAcceptLanguage(acceptLanguage string) *DeleteVdcOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteVdcOptions) SetHeaders(param map[string]string) *DeleteVdcOptions {
	options.Headers = param
	return options
}

// DirectorSite : A Cloud Director site resource. The Cloud Director site instance is the infrastructure and the associated VMware
// software stack, which consists of VMware vCenter Server, VMware NSX-T, and VMware Cloud Director.
type DirectorSite struct {
	// A unique ID for the Cloud Director site in IBM Cloud.
	Crn *string `json:"crn" validate:"required"`

	// The hyperlink of the Cloud Director site resource.
	Href *string `json:"href" validate:"required"`

	// ID of the Cloud Director site.
	ID *string `json:"id" validate:"required"`

	// The time that the Cloud Director site is ordered.
	OrderedAt *strfmt.DateTime `json:"ordered_at" validate:"required"`

	// The time that the Cloud Director site is provisioned and available to use.
	ProvisionedAt *strfmt.DateTime `json:"provisioned_at,omitempty"`

	// The name of Cloud Director site. The name of the Cloud Director site cannot be changed after creation.
	Name *string `json:"name" validate:"required"`

	// The status of Cloud Director site.
	Status *string `json:"status" validate:"required"`

	// The resource group information to associate with the resource instance.
	ResourceGroup *ResourceGroupReference `json:"resource_group" validate:"required"`

	// List of VMware resource pools to deploy on the instance.
	Pvdcs []PVDC `json:"pvdcs" validate:"required"`

	// Director site type.
	Type *string `json:"type" validate:"required"`

	// Services on the Cloud Director site.
	Services []Service `json:"services" validate:"required"`

	// RHEL activation key. This property is applicable when type is multitenant.
	RhelVmActivationKey *string `json:"rhel_vm_activation_key,omitempty"`

	// Type of director console connection.
	ConsoleConnectionType *string `json:"console_connection_type" validate:"required"`

	// Status of director console connection.
	ConsoleConnectionStatus *string `json:"console_connection_status" validate:"required"`

	// Director console IP allowlist.
	IpAllowList []string `json:"ip_allow_list" validate:"required"`
}

// Constants associated with the DirectorSite.Status property.
// The status of Cloud Director site.
const (
	DirectorSite_Status_Creating = "creating"
	DirectorSite_Status_Deleted = "deleted"
	DirectorSite_Status_Deleting = "deleting"
	DirectorSite_Status_ReadyToUse = "ready_to_use"
	DirectorSite_Status_Updating = "updating"
)

// Constants associated with the DirectorSite.Type property.
// Director site type.
const (
	DirectorSite_Type_Multitenant = "multitenant"
	DirectorSite_Type_SingleTenant = "single_tenant"
)

// Constants associated with the DirectorSite.ConsoleConnectionType property.
// Type of director console connection.
const (
	DirectorSite_ConsoleConnectionType_Private = "private"
	DirectorSite_ConsoleConnectionType_Public = "public"
)

// Constants associated with the DirectorSite.ConsoleConnectionStatus property.
// Status of director console connection.
const (
	DirectorSite_ConsoleConnectionStatus_Creating = "creating"
	DirectorSite_ConsoleConnectionStatus_Deleted = "deleted"
	DirectorSite_ConsoleConnectionStatus_Deleting = "deleting"
	DirectorSite_ConsoleConnectionStatus_ReadyToUse = "ready_to_use"
	DirectorSite_ConsoleConnectionStatus_Updating = "updating"
)

// UnmarshalDirectorSite unmarshals an instance of DirectorSite from the specified map of raw messages.
func UnmarshalDirectorSite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSite)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_at", &obj.OrderedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "ordered_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioned_at", &obj.ProvisionedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioned_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_group", &obj.ResourceGroup, UnmarshalResourceGroupReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pvdcs", &obj.Pvdcs, UnmarshalPVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "pvdcs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "services", &obj.Services, UnmarshalService)
	if err != nil {
		err = core.SDKErrorf(err, "", "services-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rhel_vm_activation_key", &obj.RhelVmActivationKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "rhel_vm_activation_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "console_connection_type", &obj.ConsoleConnectionType)
	if err != nil {
		err = core.SDKErrorf(err, "", "console_connection_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "console_connection_status", &obj.ConsoleConnectionStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "console_connection_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_allow_list", &obj.IpAllowList)
	if err != nil {
		err = core.SDKErrorf(err, "", "ip_allow_list-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteCollection : Return all Cloud Director site instances.
type DirectorSiteCollection struct {
	// List of Cloud Director site instances.
	DirectorSites []DirectorSite `json:"director_sites" validate:"required"`
}

// UnmarshalDirectorSiteCollection unmarshals an instance of DirectorSiteCollection from the specified map of raw messages.
func UnmarshalDirectorSiteCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteCollection)
	err = core.UnmarshalModel(m, "director_sites", &obj.DirectorSites, UnmarshalDirectorSite)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_sites-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteHostProfile : Host profile template.
type DirectorSiteHostProfile struct {
	// The ID for this host profile.
	ID *string `json:"id" validate:"required"`

	// The number CPU cores for this host profile.
	Cpu *int64 `json:"cpu" validate:"required"`

	// The CPU family for this host profile.
	Family *string `json:"family" validate:"required"`

	// The CPU type for this host profile.
	Processor *string `json:"processor" validate:"required"`

	// The RAM for this host profile in GB (1024^3 bytes).
	Ram *int64 `json:"ram" validate:"required"`

	// The number of CPU sockets available for this host profile.
	Socket *int64 `json:"socket" validate:"required"`

	// The CPU clock speed.
	Speed *string `json:"speed" validate:"required"`

	// The manufacturer for this host profile.
	Manufacturer *string `json:"manufacturer" validate:"required"`

	// Additional features for this host profile.
	Features []string `json:"features" validate:"required"`
}

// UnmarshalDirectorSiteHostProfile unmarshals an instance of DirectorSiteHostProfile from the specified map of raw messages.
func UnmarshalDirectorSiteHostProfile(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteHostProfile)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cpu", &obj.Cpu)
	if err != nil {
		err = core.SDKErrorf(err, "", "cpu-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "family", &obj.Family)
	if err != nil {
		err = core.SDKErrorf(err, "", "family-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "processor", &obj.Processor)
	if err != nil {
		err = core.SDKErrorf(err, "", "processor-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ram", &obj.Ram)
	if err != nil {
		err = core.SDKErrorf(err, "", "ram-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "socket", &obj.Socket)
	if err != nil {
		err = core.SDKErrorf(err, "", "socket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "speed", &obj.Speed)
	if err != nil {
		err = core.SDKErrorf(err, "", "speed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "manufacturer", &obj.Manufacturer)
	if err != nil {
		err = core.SDKErrorf(err, "", "manufacturer-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "features", &obj.Features)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteHostProfileCollection : Success. The request was successfully processed.
type DirectorSiteHostProfileCollection struct {
	// The list of available host profiles.
	DirectorSiteHostProfiles []DirectorSiteHostProfile `json:"director_site_host_profiles" validate:"required"`
}

// UnmarshalDirectorSiteHostProfileCollection unmarshals an instance of DirectorSiteHostProfileCollection from the specified map of raw messages.
func UnmarshalDirectorSiteHostProfileCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteHostProfileCollection)
	err = core.UnmarshalModel(m, "director_site_host_profiles", &obj.DirectorSiteHostProfiles, UnmarshalDirectorSiteHostProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_site_host_profiles-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSitePVDC : The resource pool within the Director Site in which to deploy the virtual data center (VDC).
type DirectorSitePVDC struct {
	// A unique ID for the resource pool.
	ID *string `json:"id" validate:"required"`

	// Determines how resources are made available to the virtual data center (VDC). Required for VDCs deployed on a
	// multitenant Cloud Director site.
	ProviderType *VDCProviderType `json:"provider_type,omitempty"`
}

// NewDirectorSitePVDC : Instantiate DirectorSitePVDC (Generic Model Constructor)
func (*VmwareV1) NewDirectorSitePVDC(id string) (_model *DirectorSitePVDC, err error) {
	_model = &DirectorSitePVDC{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalDirectorSitePVDC unmarshals an instance of DirectorSitePVDC from the specified map of raw messages.
func UnmarshalDirectorSitePVDC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSitePVDC)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "provider_type", &obj.ProviderType, UnmarshalVDCProviderType)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteReference : Back link to associated Cloud Director site resource.
type DirectorSiteReference struct {
	// A unique ID for the Cloud Director site in IBM Cloud.
	Crn *string `json:"crn" validate:"required"`

	// The hyperlink of the Cloud Director site resource.
	Href *string `json:"href" validate:"required"`

	// ID of the Cloud Director site.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalDirectorSiteReference unmarshals an instance of DirectorSiteReference from the specified map of raw messages.
func UnmarshalDirectorSiteReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteReference)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteRegion : The region details.
type DirectorSiteRegion struct {
	// The region name.
	Name *string `json:"name" validate:"required"`

	// The data center details.
	DataCenters []DataCenter `json:"data_centers" validate:"required"`

	// Accessible endpoint of the region.
	Endpoint *string `json:"endpoint" validate:"required"`
}

// UnmarshalDirectorSiteRegion unmarshals an instance of DirectorSiteRegion from the specified map of raw messages.
func UnmarshalDirectorSiteRegion(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteRegion)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "data_centers", &obj.DataCenters, UnmarshalDataCenter)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_centers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DirectorSiteRegionCollection : Success. The request was successfully processed.
type DirectorSiteRegionCollection struct {
	// Regions of Cloud Director sites.
	DirectorSiteRegions []DirectorSiteRegion `json:"director_site_regions" validate:"required"`
}

// UnmarshalDirectorSiteRegionCollection unmarshals an instance of DirectorSiteRegionCollection from the specified map of raw messages.
func UnmarshalDirectorSiteRegionCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DirectorSiteRegionCollection)
	err = core.UnmarshalModel(m, "director_site_regions", &obj.DirectorSiteRegions, UnmarshalDirectorSiteRegion)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_site_regions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Edge : A networking edge deployed on a virtual data center (VDC). Networking edges are based on NSX-T and used for bridging
// virtualize networking to the physical public-internet and IBM private networking.
type Edge struct {
	// A unique ID for the edge.
	ID *string `json:"id" validate:"required"`

	// The public IP addresses assigned to the edge.
	PublicIps []string `json:"public_ips" validate:"required"`

	// The private IP addresses assigned to the edge.
	PrivateIps []string `json:"private_ips" validate:"required"`

	// Indicates whether the edge is private only. The default value is True for a private Cloud Director site and False
	// for a public Cloud Director site.
	PrivateOnly *bool `json:"private_only,omitempty"`

	// The size of the edge.
	//
	// The size can be specified only for performance edges. Larger sizes require more capacity from the Cloud Director
	// site in which the virtual data center (VDC) was created to be deployed.
	Size *string `json:"size" validate:"required"`

	// Determines the state of the edge.
	Status *string `json:"status" validate:"required"`

	// Connected IBM Transit Gateways.
	TransitGateways []TransitGateway `json:"transit_gateways" validate:"required"`

	// The type of edge to be deployed.
	//
	// Efficiency edges allow for multiple VDCs to share some edge resources. Performance edges do not share resources
	// between VDCs.
	Type *string `json:"type" validate:"required"`

	// The edge version.
	Version *string `json:"version" validate:"required"`
}

// Constants associated with the Edge.Size property.
// The size of the edge.
//
// The size can be specified only for performance edges. Larger sizes require more capacity from the Cloud Director site
// in which the virtual data center (VDC) was created to be deployed.
const (
	Edge_Size_ExtraLarge = "extra_large"
	Edge_Size_Large = "large"
	Edge_Size_Medium = "medium"
)

// Constants associated with the Edge.Status property.
// Determines the state of the edge.
const (
	Edge_Status_Creating = "creating"
	Edge_Status_Deleted = "deleted"
	Edge_Status_Deleting = "deleting"
	Edge_Status_ReadyToUse = "ready_to_use"
)

// Constants associated with the Edge.Type property.
// The type of edge to be deployed.
//
// Efficiency edges allow for multiple VDCs to share some edge resources. Performance edges do not share resources
// between VDCs.
const (
	Edge_Type_Efficiency = "efficiency"
	Edge_Type_Performance = "performance"
)

// UnmarshalEdge unmarshals an instance of Edge from the specified map of raw messages.
func UnmarshalEdge(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Edge)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "public_ips", &obj.PublicIps)
	if err != nil {
		err = core.SDKErrorf(err, "", "public_ips-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_ips", &obj.PrivateIps)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_ips-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_only", &obj.PrivateOnly)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_only-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "transit_gateways", &obj.TransitGateways, UnmarshalTransitGateway)
	if err != nil {
		err = core.SDKErrorf(err, "", "transit_gateways-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnableVcdaOnDataCenterOptions : The EnableVcdaOnDataCenter options.
type EnableVcdaOnDataCenterOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Indicates whether it is required to enable or disable a service.
	Enable *bool `json:"enable" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewEnableVcdaOnDataCenterOptions : Instantiate EnableVcdaOnDataCenterOptions
func (*VmwareV1) NewEnableVcdaOnDataCenterOptions(siteID string, enable bool) *EnableVcdaOnDataCenterOptions {
	return &EnableVcdaOnDataCenterOptions{
		SiteID: core.StringPtr(siteID),
		Enable: core.BoolPtr(enable),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *EnableVcdaOnDataCenterOptions) SetSiteID(siteID string) *EnableVcdaOnDataCenterOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetEnable : Allow user to set Enable
func (_options *EnableVcdaOnDataCenterOptions) SetEnable(enable bool) *EnableVcdaOnDataCenterOptions {
	_options.Enable = core.BoolPtr(enable)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *EnableVcdaOnDataCenterOptions) SetAcceptLanguage(acceptLanguage string) *EnableVcdaOnDataCenterOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *EnableVcdaOnDataCenterOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *EnableVcdaOnDataCenterOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *EnableVcdaOnDataCenterOptions) SetHeaders(param map[string]string) *EnableVcdaOnDataCenterOptions {
	options.Headers = param
	return options
}

// EnableVeeamOnPvdcsListOptions : The EnableVeeamOnPvdcsList options.
type EnableVeeamOnPvdcsListOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Indicates whether it is required to enable or disable a service.
	Enable *bool `json:"enable" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewEnableVeeamOnPvdcsListOptions : Instantiate EnableVeeamOnPvdcsListOptions
func (*VmwareV1) NewEnableVeeamOnPvdcsListOptions(siteID string, enable bool) *EnableVeeamOnPvdcsListOptions {
	return &EnableVeeamOnPvdcsListOptions{
		SiteID: core.StringPtr(siteID),
		Enable: core.BoolPtr(enable),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *EnableVeeamOnPvdcsListOptions) SetSiteID(siteID string) *EnableVeeamOnPvdcsListOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetEnable : Allow user to set Enable
func (_options *EnableVeeamOnPvdcsListOptions) SetEnable(enable bool) *EnableVeeamOnPvdcsListOptions {
	_options.Enable = core.BoolPtr(enable)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *EnableVeeamOnPvdcsListOptions) SetAcceptLanguage(acceptLanguage string) *EnableVeeamOnPvdcsListOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *EnableVeeamOnPvdcsListOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *EnableVeeamOnPvdcsListOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *EnableVeeamOnPvdcsListOptions) SetHeaders(param map[string]string) *EnableVeeamOnPvdcsListOptions {
	options.Headers = param
	return options
}

// FileShares : Chosen storage policies and their sizes.
type FileShares struct {
	// The amount of 0.25 IOPS/GB storage in GB (1024^3 bytes).
	STORAGEPOINTTWOFIVEIOPSGB *int64 `json:"STORAGE_POINT_TWO_FIVE_IOPS_GB,omitempty"`

	// The amount of 2 IOPS/GB storage in GB (1024^3 bytes).
	STORAGETWOIOPSGB *int64 `json:"STORAGE_TWO_IOPS_GB,omitempty"`

	// The amount of 4 IOPS/GB storage in GB (1024^3 bytes).
	STORAGEFOURIOPSGB *int64 `json:"STORAGE_FOUR_IOPS_GB,omitempty"`

	// The amount of 10 IOPS/GB storage in GB (1024^3 bytes).
	STORAGETENIOPSGB *int64 `json:"STORAGE_TEN_IOPS_GB,omitempty"`
}

// UnmarshalFileShares unmarshals an instance of FileShares from the specified map of raw messages.
func UnmarshalFileShares(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FileShares)
	err = core.UnmarshalPrimitive(m, "STORAGE_POINT_TWO_FIVE_IOPS_GB", &obj.STORAGEPOINTTWOFIVEIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_POINT_TWO_FIVE_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_TWO_IOPS_GB", &obj.STORAGETWOIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_TWO_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_FOUR_IOPS_GB", &obj.STORAGEFOURIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_FOUR_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_TEN_IOPS_GB", &obj.STORAGETENIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_TEN_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FileSharesPrototype : Chosen storage policies and their sizes.
type FileSharesPrototype struct {
	// The amount of 0.25 IOPS/GB storage in GB (1024^3 bytes).
	STORAGEPOINTTWOFIVEIOPSGB *int64 `json:"STORAGE_POINT_TWO_FIVE_IOPS_GB,omitempty"`

	// The amount of 2 IOPS/GB storage in GB (1024^3 bytes).
	STORAGETWOIOPSGB *int64 `json:"STORAGE_TWO_IOPS_GB,omitempty"`

	// The amount of 4 IOPS/GB storage in GB (1024^3 bytes).
	STORAGEFOURIOPSGB *int64 `json:"STORAGE_FOUR_IOPS_GB,omitempty"`

	// The amount of 10 IOPS/GB storage in GB (1024^3 bytes).
	STORAGETENIOPSGB *int64 `json:"STORAGE_TEN_IOPS_GB,omitempty"`
}

// UnmarshalFileSharesPrototype unmarshals an instance of FileSharesPrototype from the specified map of raw messages.
func UnmarshalFileSharesPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FileSharesPrototype)
	err = core.UnmarshalPrimitive(m, "STORAGE_POINT_TWO_FIVE_IOPS_GB", &obj.STORAGEPOINTTWOFIVEIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_POINT_TWO_FIVE_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_TWO_IOPS_GB", &obj.STORAGETWOIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_TWO_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_FOUR_IOPS_GB", &obj.STORAGEFOURIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_FOUR_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "STORAGE_TEN_IOPS_GB", &obj.STORAGETENIOPSGB)
	if err != nil {
		err = core.SDKErrorf(err, "", "STORAGE_TEN_IOPS_GB-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the FileSharesPrototype
func (fileSharesPrototype *FileSharesPrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(fileSharesPrototype.STORAGEPOINTTWOFIVEIOPSGB) {
		_patch["STORAGE_POINT_TWO_FIVE_IOPS_GB"] = fileSharesPrototype.STORAGEPOINTTWOFIVEIOPSGB
	}
	if !core.IsNil(fileSharesPrototype.STORAGETWOIOPSGB) {
		_patch["STORAGE_TWO_IOPS_GB"] = fileSharesPrototype.STORAGETWOIOPSGB
	}
	if !core.IsNil(fileSharesPrototype.STORAGEFOURIOPSGB) {
		_patch["STORAGE_FOUR_IOPS_GB"] = fileSharesPrototype.STORAGEFOURIOPSGB
	}
	if !core.IsNil(fileSharesPrototype.STORAGETENIOPSGB) {
		_patch["STORAGE_TEN_IOPS_GB"] = fileSharesPrototype.STORAGETENIOPSGB
	}

	return
}

// GetDirectorInstancesPvdcsClusterOptions : The GetDirectorInstancesPvdcsCluster options.
type GetDirectorInstancesPvdcsClusterOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// The cluster to query.
	ID *string `json:"id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	PvdcID *string `json:"pvdc_id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDirectorInstancesPvdcsClusterOptions : Instantiate GetDirectorInstancesPvdcsClusterOptions
func (*VmwareV1) NewGetDirectorInstancesPvdcsClusterOptions(siteID string, id string, pvdcID string) *GetDirectorInstancesPvdcsClusterOptions {
	return &GetDirectorInstancesPvdcsClusterOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
		PvdcID: core.StringPtr(pvdcID),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *GetDirectorInstancesPvdcsClusterOptions) SetSiteID(siteID string) *GetDirectorInstancesPvdcsClusterOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetDirectorInstancesPvdcsClusterOptions) SetID(id string) *GetDirectorInstancesPvdcsClusterOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetPvdcID : Allow user to set PvdcID
func (_options *GetDirectorInstancesPvdcsClusterOptions) SetPvdcID(pvdcID string) *GetDirectorInstancesPvdcsClusterOptions {
	_options.PvdcID = core.StringPtr(pvdcID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetDirectorInstancesPvdcsClusterOptions) SetAcceptLanguage(acceptLanguage string) *GetDirectorInstancesPvdcsClusterOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *GetDirectorInstancesPvdcsClusterOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *GetDirectorInstancesPvdcsClusterOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDirectorInstancesPvdcsClusterOptions) SetHeaders(param map[string]string) *GetDirectorInstancesPvdcsClusterOptions {
	options.Headers = param
	return options
}

// GetDirectorSiteOptions : The GetDirectorSite options.
type GetDirectorSiteOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDirectorSiteOptions : Instantiate GetDirectorSiteOptions
func (*VmwareV1) NewGetDirectorSiteOptions(id string) *GetDirectorSiteOptions {
	return &GetDirectorSiteOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetDirectorSiteOptions) SetID(id string) *GetDirectorSiteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetDirectorSiteOptions) SetAcceptLanguage(acceptLanguage string) *GetDirectorSiteOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *GetDirectorSiteOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *GetDirectorSiteOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDirectorSiteOptions) SetHeaders(param map[string]string) *GetDirectorSiteOptions {
	options.Headers = param
	return options
}

// GetDirectorSitesPvdcsOptions : The GetDirectorSitesPvdcs options.
type GetDirectorSitesPvdcsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDirectorSitesPvdcsOptions : Instantiate GetDirectorSitesPvdcsOptions
func (*VmwareV1) NewGetDirectorSitesPvdcsOptions(siteID string, id string) *GetDirectorSitesPvdcsOptions {
	return &GetDirectorSitesPvdcsOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *GetDirectorSitesPvdcsOptions) SetSiteID(siteID string) *GetDirectorSitesPvdcsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetDirectorSitesPvdcsOptions) SetID(id string) *GetDirectorSitesPvdcsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetDirectorSitesPvdcsOptions) SetAcceptLanguage(acceptLanguage string) *GetDirectorSitesPvdcsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *GetDirectorSitesPvdcsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *GetDirectorSitesPvdcsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDirectorSitesPvdcsOptions) SetHeaders(param map[string]string) *GetDirectorSitesPvdcsOptions {
	options.Headers = param
	return options
}

// GetOidcConfigurationOptions : The GetOidcConfiguration options.
type GetOidcConfigurationOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetOidcConfigurationOptions : Instantiate GetOidcConfigurationOptions
func (*VmwareV1) NewGetOidcConfigurationOptions(siteID string) *GetOidcConfigurationOptions {
	return &GetOidcConfigurationOptions{
		SiteID: core.StringPtr(siteID),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *GetOidcConfigurationOptions) SetSiteID(siteID string) *GetOidcConfigurationOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetOidcConfigurationOptions) SetAcceptLanguage(acceptLanguage string) *GetOidcConfigurationOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOidcConfigurationOptions) SetHeaders(param map[string]string) *GetOidcConfigurationOptions {
	options.Headers = param
	return options
}

// GetVdcOptions : The GetVdc options.
type GetVdcOptions struct {
	// A unique ID for a specified virtual data center.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetVdcOptions : Instantiate GetVdcOptions
func (*VmwareV1) NewGetVdcOptions(id string) *GetVdcOptions {
	return &GetVdcOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetVdcOptions) SetID(id string) *GetVdcOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *GetVdcOptions) SetAcceptLanguage(acceptLanguage string) *GetVdcOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetVdcOptions) SetHeaders(param map[string]string) *GetVdcOptions {
	options.Headers = param
	return options
}

// ListDirectorSiteHostProfilesOptions : The ListDirectorSiteHostProfiles options.
type ListDirectorSiteHostProfilesOptions struct {
	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDirectorSiteHostProfilesOptions : Instantiate ListDirectorSiteHostProfilesOptions
func (*VmwareV1) NewListDirectorSiteHostProfilesOptions() *ListDirectorSiteHostProfilesOptions {
	return &ListDirectorSiteHostProfilesOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListDirectorSiteHostProfilesOptions) SetAcceptLanguage(acceptLanguage string) *ListDirectorSiteHostProfilesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListDirectorSiteHostProfilesOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListDirectorSiteHostProfilesOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDirectorSiteHostProfilesOptions) SetHeaders(param map[string]string) *ListDirectorSiteHostProfilesOptions {
	options.Headers = param
	return options
}

// ListDirectorSiteRegionsOptions : The ListDirectorSiteRegions options.
type ListDirectorSiteRegionsOptions struct {
	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDirectorSiteRegionsOptions : Instantiate ListDirectorSiteRegionsOptions
func (*VmwareV1) NewListDirectorSiteRegionsOptions() *ListDirectorSiteRegionsOptions {
	return &ListDirectorSiteRegionsOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListDirectorSiteRegionsOptions) SetAcceptLanguage(acceptLanguage string) *ListDirectorSiteRegionsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListDirectorSiteRegionsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListDirectorSiteRegionsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDirectorSiteRegionsOptions) SetHeaders(param map[string]string) *ListDirectorSiteRegionsOptions {
	options.Headers = param
	return options
}

// ListDirectorSitesOptions : The ListDirectorSites options.
type ListDirectorSitesOptions struct {
	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDirectorSitesOptions : Instantiate ListDirectorSitesOptions
func (*VmwareV1) NewListDirectorSitesOptions() *ListDirectorSitesOptions {
	return &ListDirectorSitesOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListDirectorSitesOptions) SetAcceptLanguage(acceptLanguage string) *ListDirectorSitesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListDirectorSitesOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListDirectorSitesOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDirectorSitesOptions) SetHeaders(param map[string]string) *ListDirectorSitesOptions {
	options.Headers = param
	return options
}

// ListDirectorSitesPvdcsClustersOptions : The ListDirectorSitesPvdcsClusters options.
type ListDirectorSitesPvdcsClustersOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	PvdcID *string `json:"pvdc_id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDirectorSitesPvdcsClustersOptions : Instantiate ListDirectorSitesPvdcsClustersOptions
func (*VmwareV1) NewListDirectorSitesPvdcsClustersOptions(siteID string, pvdcID string) *ListDirectorSitesPvdcsClustersOptions {
	return &ListDirectorSitesPvdcsClustersOptions{
		SiteID: core.StringPtr(siteID),
		PvdcID: core.StringPtr(pvdcID),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *ListDirectorSitesPvdcsClustersOptions) SetSiteID(siteID string) *ListDirectorSitesPvdcsClustersOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetPvdcID : Allow user to set PvdcID
func (_options *ListDirectorSitesPvdcsClustersOptions) SetPvdcID(pvdcID string) *ListDirectorSitesPvdcsClustersOptions {
	_options.PvdcID = core.StringPtr(pvdcID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListDirectorSitesPvdcsClustersOptions) SetAcceptLanguage(acceptLanguage string) *ListDirectorSitesPvdcsClustersOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListDirectorSitesPvdcsClustersOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListDirectorSitesPvdcsClustersOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDirectorSitesPvdcsClustersOptions) SetHeaders(param map[string]string) *ListDirectorSitesPvdcsClustersOptions {
	options.Headers = param
	return options
}

// ListDirectorSitesPvdcsOptions : The ListDirectorSitesPvdcs options.
type ListDirectorSitesPvdcsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDirectorSitesPvdcsOptions : Instantiate ListDirectorSitesPvdcsOptions
func (*VmwareV1) NewListDirectorSitesPvdcsOptions(siteID string) *ListDirectorSitesPvdcsOptions {
	return &ListDirectorSitesPvdcsOptions{
		SiteID: core.StringPtr(siteID),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *ListDirectorSitesPvdcsOptions) SetSiteID(siteID string) *ListDirectorSitesPvdcsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListDirectorSitesPvdcsOptions) SetAcceptLanguage(acceptLanguage string) *ListDirectorSitesPvdcsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListDirectorSitesPvdcsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListDirectorSitesPvdcsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDirectorSitesPvdcsOptions) SetHeaders(param map[string]string) *ListDirectorSitesPvdcsOptions {
	options.Headers = param
	return options
}

// ListMultitenantDirectorSitesOptions : The ListMultitenantDirectorSites options.
type ListMultitenantDirectorSitesOptions struct {
	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListMultitenantDirectorSitesOptions : Instantiate ListMultitenantDirectorSitesOptions
func (*VmwareV1) NewListMultitenantDirectorSitesOptions() *ListMultitenantDirectorSitesOptions {
	return &ListMultitenantDirectorSitesOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListMultitenantDirectorSitesOptions) SetAcceptLanguage(acceptLanguage string) *ListMultitenantDirectorSitesOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *ListMultitenantDirectorSitesOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *ListMultitenantDirectorSitesOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListMultitenantDirectorSitesOptions) SetHeaders(param map[string]string) *ListMultitenantDirectorSitesOptions {
	options.Headers = param
	return options
}

// ListVdcsOptions : The ListVdcs options.
type ListVdcsOptions struct {
	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListVdcsOptions : Instantiate ListVdcsOptions
func (*VmwareV1) NewListVdcsOptions() *ListVdcsOptions {
	return &ListVdcsOptions{}
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *ListVdcsOptions) SetAcceptLanguage(acceptLanguage string) *ListVdcsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListVdcsOptions) SetHeaders(param map[string]string) *ListVdcsOptions {
	options.Headers = param
	return options
}

// MultitenantDirectorSite : Multitenant Cloud Director site details.
type MultitenantDirectorSite struct {
	// Multitenant Cloud Director site name.
	Name *string `json:"name" validate:"required"`

	// Multitenant Cloud Director site display name.
	DisplayName *string `json:"display_name" validate:"required"`

	// Multitenant Cloud Director site ID.
	ID *string `json:"id" validate:"required"`

	// Indicates whether the site is private only.
	PrivateOnly *bool `json:"private_only" validate:"required"`

	// Multitenant Cloud Director site region name.
	Region *string `json:"region" validate:"required"`

	// Resource pool details.
	Pvdcs []MultitenantPVDC `json:"pvdcs" validate:"required"`

	// Installed services.
	Services []string `json:"services" validate:"required"`
}

// Constants associated with the MultitenantDirectorSite.Services property.
const (
	MultitenantDirectorSite_Services_Vcda = "vcda"
	MultitenantDirectorSite_Services_Veeam = "veeam"
)

// UnmarshalMultitenantDirectorSite unmarshals an instance of MultitenantDirectorSite from the specified map of raw messages.
func UnmarshalMultitenantDirectorSite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MultitenantDirectorSite)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_only", &obj.PrivateOnly)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_only-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pvdcs", &obj.Pvdcs, UnmarshalMultitenantPVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "pvdcs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "services", &obj.Services)
	if err != nil {
		err = core.SDKErrorf(err, "", "services-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MultitenantDirectorSiteCollection : List of multitenant Cloud Director sites.
type MultitenantDirectorSiteCollection struct {
	// Multitenant Cloud Director sites.
	MultitenantDirectorSites []MultitenantDirectorSite `json:"multitenant_director_sites" validate:"required"`
}

// UnmarshalMultitenantDirectorSiteCollection unmarshals an instance of MultitenantDirectorSiteCollection from the specified map of raw messages.
func UnmarshalMultitenantDirectorSiteCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MultitenantDirectorSiteCollection)
	err = core.UnmarshalModel(m, "multitenant_director_sites", &obj.MultitenantDirectorSites, UnmarshalMultitenantDirectorSite)
	if err != nil {
		err = core.SDKErrorf(err, "", "multitenant_director_sites-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MultitenantPVDC : Multitenant resource pool details.
type MultitenantPVDC struct {
	// Resource pool name.
	Name *string `json:"name" validate:"required"`

	// Resource pool ID.
	ID *string `json:"id" validate:"required"`

	// Data center name.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// Indicates whether the resource pool is private only.
	PrivateOnly *bool `json:"private_only" validate:"required"`

	// List of resource pool types.
	ProviderTypes []ProviderType `json:"provider_types" validate:"required"`
}

// UnmarshalMultitenantPVDC unmarshals an instance of MultitenantPVDC from the specified map of raw messages.
func UnmarshalMultitenantPVDC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MultitenantPVDC)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_only", &obj.PrivateOnly)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_only-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "provider_types", &obj.ProviderTypes, UnmarshalProviderType)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OIDC : Details of the OIDC configuration on a Cloud Director site.
type OIDC struct {
	// Status of the OIDC configuration on a Cloud Director site.
	Status *string `json:"status" validate:"required"`

	// The time after which the OIDC configuration is considered enabled.
	LastSetAt *strfmt.DateTime `json:"last_set_at,omitempty"`
}

// Constants associated with the OIDC.Status property.
// Status of the OIDC configuration on a Cloud Director site.
const (
	OIDC_Status_Deleted = "deleted"
	OIDC_Status_Pending = "pending"
	OIDC_Status_ReadyToUse = "ready_to_use"
)

// UnmarshalOIDC unmarshals an instance of OIDC from the specified map of raw messages.
func UnmarshalOIDC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OIDC)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_set_at", &obj.LastSetAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_set_at-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PVDC : VMware resource pool information.
type PVDC struct {
	// Name of the resource pool. Resource pool names must be unique per Cloud Director site instance and they cannot be
	// changed after creation.
	Name *string `json:"name" validate:"required"`

	// Data center location to deploy the cluster. See `GET /director_site_regions` for supported data center locations.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// The resource pool ID.
	ID *string `json:"id" validate:"required"`

	// The hyperlink of the resource pool resource.
	Href *string `json:"href" validate:"required"`

	// List of VMware clusters to deploy on the instance. Clusters form VMware workload availability boundaries.
	Clusters []ClusterSummary `json:"clusters" validate:"required"`

	// The status of the resource pool.
	Status *string `json:"status" validate:"required"`

	// List of resource pool types.
	ProviderTypes []ProviderType `json:"provider_types" validate:"required"`
}

// Constants associated with the PVDC.Status property.
// The status of the resource pool.
const (
	PVDC_Status_Creating = "creating"
	PVDC_Status_Deleted = "deleted"
	PVDC_Status_Deleting = "deleting"
	PVDC_Status_ReadyToUse = "ready_to_use"
	PVDC_Status_Updating = "updating"
)

// UnmarshalPVDC unmarshals an instance of PVDC from the specified map of raw messages.
func UnmarshalPVDC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PVDC)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "clusters", &obj.Clusters, UnmarshalClusterSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "clusters-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "provider_types", &obj.ProviderTypes, UnmarshalProviderType)
	if err != nil {
		err = core.SDKErrorf(err, "", "provider_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PVDCCollection : Return all resource pool instances.
type PVDCCollection struct {
	// List of resource pool instances.
	Pvdcs []PVDC `json:"pvdcs" validate:"required"`
}

// UnmarshalPVDCCollection unmarshals an instance of PVDCCollection from the specified map of raw messages.
func UnmarshalPVDCCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PVDCCollection)
	err = core.UnmarshalModel(m, "pvdcs", &obj.Pvdcs, UnmarshalPVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "pvdcs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PVDCPrototype : VMware resource pool order information.
type PVDCPrototype struct {
	// Name of the resource pool. Resource pool names must be unique per Cloud Director site instance and they cannot be
	// changed after creation.
	Name *string `json:"name" validate:"required"`

	// Data center location to deploy the cluster. See `GET /director_site_regions` for supported data center locations.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// List of VMware clusters to deploy on the instance. Clusters form VMware workload availability boundaries.
	Clusters []ClusterPrototype `json:"clusters" validate:"required"`
}

// NewPVDCPrototype : Instantiate PVDCPrototype (Generic Model Constructor)
func (*VmwareV1) NewPVDCPrototype(name string, dataCenterName string, clusters []ClusterPrototype) (_model *PVDCPrototype, err error) {
	_model = &PVDCPrototype{
		Name: core.StringPtr(name),
		DataCenterName: core.StringPtr(dataCenterName),
		Clusters: clusters,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPVDCPrototype unmarshals an instance of PVDCPrototype from the specified map of raw messages.
func UnmarshalPVDCPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PVDCPrototype)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "clusters", &obj.Clusters, UnmarshalClusterPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "clusters-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProviderType : Resource pool type.
type ProviderType struct {
	// The name of the resource pool type.
	Name *string `json:"name" validate:"required"`
}

// Constants associated with the ProviderType.Name property.
// The name of the resource pool type.
const (
	ProviderType_Name_OnDemand = "on_demand"
	ProviderType_Name_Reserved = "reserved"
)

// UnmarshalProviderType unmarshals an instance of ProviderType from the specified map of raw messages.
func UnmarshalProviderType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProviderType)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RemoveTransitGatewayConnectionsOptions : The RemoveTransitGatewayConnections options.
type RemoveTransitGatewayConnectionsOptions struct {
	// A unique ID for a virtual data center.
	VdcID *string `json:"vdc_id" validate:"required,ne="`

	// A unique ID for an edge.
	EdgeID *string `json:"edge_id" validate:"required,ne="`

	// A unique ID for an IBM Transit Gateway.
	ID *string `json:"id" validate:"required,ne="`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewRemoveTransitGatewayConnectionsOptions : Instantiate RemoveTransitGatewayConnectionsOptions
func (*VmwareV1) NewRemoveTransitGatewayConnectionsOptions(vdcID string, edgeID string, id string) *RemoveTransitGatewayConnectionsOptions {
	return &RemoveTransitGatewayConnectionsOptions{
		VdcID: core.StringPtr(vdcID),
		EdgeID: core.StringPtr(edgeID),
		ID: core.StringPtr(id),
	}
}

// SetVdcID : Allow user to set VdcID
func (_options *RemoveTransitGatewayConnectionsOptions) SetVdcID(vdcID string) *RemoveTransitGatewayConnectionsOptions {
	_options.VdcID = core.StringPtr(vdcID)
	return _options
}

// SetEdgeID : Allow user to set EdgeID
func (_options *RemoveTransitGatewayConnectionsOptions) SetEdgeID(edgeID string) *RemoveTransitGatewayConnectionsOptions {
	_options.EdgeID = core.StringPtr(edgeID)
	return _options
}

// SetID : Allow user to set ID
func (_options *RemoveTransitGatewayConnectionsOptions) SetID(id string) *RemoveTransitGatewayConnectionsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *RemoveTransitGatewayConnectionsOptions) SetAcceptLanguage(acceptLanguage string) *RemoveTransitGatewayConnectionsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RemoveTransitGatewayConnectionsOptions) SetHeaders(param map[string]string) *RemoveTransitGatewayConnectionsOptions {
	options.Headers = param
	return options
}

// ResourceGroupIdentity : The resource group to associate with the resource instance. If not specified, the default resource group in the
// account is used.
type ResourceGroupIdentity struct {
	// A unique ID for the resource group.
	ID *string `json:"id" validate:"required"`
}

// NewResourceGroupIdentity : Instantiate ResourceGroupIdentity (Generic Model Constructor)
func (*VmwareV1) NewResourceGroupIdentity(id string) (_model *ResourceGroupIdentity, err error) {
	_model = &ResourceGroupIdentity{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResourceGroupIdentity unmarshals an instance of ResourceGroupIdentity from the specified map of raw messages.
func UnmarshalResourceGroupIdentity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupIdentity)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceGroupReference : The resource group information to associate with the resource instance.
type ResourceGroupReference struct {
	// A unique ID for the resource group.
	ID *string `json:"id" validate:"required"`

	// The name of the resource group.
	Name *string `json:"name" validate:"required"`

	// The cloud reference name for the resource group.
	Crn *string `json:"crn" validate:"required"`
}

// UnmarshalResourceGroupReference unmarshals an instance of ResourceGroupReference from the specified map of raw messages.
func UnmarshalResourceGroupReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Service : Service response body.
type Service struct {
	// The service name.
	Name *string `json:"name" validate:"required"`

	// A unique ID for the service.
	ID *string `json:"id" validate:"required"`

	// The time that the service instance is ordered.
	OrderedAt *strfmt.DateTime `json:"ordered_at" validate:"required"`

	// The time that the service instance is provisioned and available to use.
	ProvisionedAt *strfmt.DateTime `json:"provisioned_at,omitempty"`

	// The service instance status.
	Status *string `json:"status" validate:"required"`

	// Service console URL. This property is applicable when the service name is veeam.
	ConsoleURL *string `json:"console_url,omitempty"`

	// Replicators for the VCDA instance.
	Replicators *int64 `json:"replicators,omitempty"`

	// Connection on a VCDA instance.
	Connections []VcdaConnection `json:"connections" validate:"required"`

	// Scale-out backup repositories created on the Veeam service instance.
	Sobrs []Sobr `json:"sobrs" validate:"required"`
}

// Constants associated with the Service.Name property.
// The service name.
const (
	Service_Name_Vcda = "vcda"
	Service_Name_Veeam = "veeam"
)

// Constants associated with the Service.Status property.
// The service instance status.
const (
	Service_Status_Creating = "creating"
	Service_Status_Deleted = "deleted"
	Service_Status_Deleting = "deleting"
	Service_Status_ReadyToUse = "ready_to_use"
	Service_Status_Updating = "updating"
)

// UnmarshalService unmarshals an instance of Service from the specified map of raw messages.
func UnmarshalService(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Service)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_at", &obj.OrderedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "ordered_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioned_at", &obj.ProvisionedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioned_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "console_url", &obj.ConsoleURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "console_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "replicators", &obj.Replicators)
	if err != nil {
		err = core.SDKErrorf(err, "", "replicators-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "connections", &obj.Connections, UnmarshalVcdaConnection)
	if err != nil {
		err = core.SDKErrorf(err, "", "connections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "sobrs", &obj.Sobrs, UnmarshalSobr)
	if err != nil {
		err = core.SDKErrorf(err, "", "sobrs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceEnabled : Enable Service response for accepted request.
type ServiceEnabled struct {
	// Request status.
	Message *string `json:"message,omitempty"`
}

// UnmarshalServiceEnabled unmarshals an instance of ServiceEnabled from the specified map of raw messages.
func UnmarshalServiceEnabled(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceEnabled)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceIdentity : Create Service request body.
type ServiceIdentity struct {
	// The service name.
	Name *string `json:"name" validate:"required"`
}

// Constants associated with the ServiceIdentity.Name property.
// The service name.
const (
	ServiceIdentity_Name_Vcda = "vcda"
	ServiceIdentity_Name_Veeam = "veeam"
)

// NewServiceIdentity : Instantiate ServiceIdentity (Generic Model Constructor)
func (*VmwareV1) NewServiceIdentity(name string) (_model *ServiceIdentity, err error) {
	_model = &ServiceIdentity{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalServiceIdentity unmarshals an instance of ServiceIdentity from the specified map of raw messages.
func UnmarshalServiceIdentity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceIdentity)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetOidcConfigurationOptions : The SetOidcConfiguration options.
type SetOidcConfigurationOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// Size of the message body in bytes.
	ContentLength *int64 `json:"Content-Length" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewSetOidcConfigurationOptions : Instantiate SetOidcConfigurationOptions
func (*VmwareV1) NewSetOidcConfigurationOptions(siteID string, contentLength int64) *SetOidcConfigurationOptions {
	return &SetOidcConfigurationOptions{
		SiteID: core.StringPtr(siteID),
		ContentLength: core.Int64Ptr(contentLength),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *SetOidcConfigurationOptions) SetSiteID(siteID string) *SetOidcConfigurationOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetContentLength : Allow user to set ContentLength
func (_options *SetOidcConfigurationOptions) SetContentLength(contentLength int64) *SetOidcConfigurationOptions {
	_options.ContentLength = core.Int64Ptr(contentLength)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *SetOidcConfigurationOptions) SetAcceptLanguage(acceptLanguage string) *SetOidcConfigurationOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetOidcConfigurationOptions) SetHeaders(param map[string]string) *SetOidcConfigurationOptions {
	options.Headers = param
	return options
}

// StatusReason : Information about why a request cannot be completed or why a resource cannot be created.
type StatusReason struct {
	// An error code specific to the error encountered.
	Code *string `json:"code" validate:"required"`

	// A message that describes why the error ocurred.
	Message *string `json:"message" validate:"required"`

	// A URL that links to a page with more information about this error.
	MoreInfo *string `json:"more_info,omitempty"`
}

// Constants associated with the StatusReason.Code property.
// An error code specific to the error encountered.
const (
	StatusReason_Code_InsufficentCpu = "insufficent_cpu"
	StatusReason_Code_InsufficentCpuAndRam = "insufficent_cpu_and_ram"
	StatusReason_Code_InsufficentRam = "insufficent_ram"
)

// UnmarshalStatusReason unmarshals an instance of StatusReason from the specified map of raw messages.
func UnmarshalStatusReason(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StatusReason)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "more_info", &obj.MoreInfo)
	if err != nil {
		err = core.SDKErrorf(err, "", "more_info-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TransitGateway : An IBM Transit Gateway.
type TransitGateway struct {
	// A unique ID for an IBM Transit Gateway.
	ID *string `json:"id" validate:"required"`

	// IBM Transit Gateway connections.
	Connections []TransitGatewayConnection `json:"connections" validate:"required"`

	// Determines the state of the IBM Transit Gateway based on its connections.
	Status *string `json:"status" validate:"required"`

	// The region where the IBM Transit Gateway is deployed.
	Region *string `json:"region" validate:"required"`
}

// Constants associated with the TransitGateway.Status property.
// Determines the state of the IBM Transit Gateway based on its connections.
const (
	TransitGateway_Status_Creating = "creating"
	TransitGateway_Status_Deleting = "deleting"
	TransitGateway_Status_Pending = "pending"
	TransitGateway_Status_ReadyToUse = "ready_to_use"
)

// UnmarshalTransitGateway unmarshals an instance of TransitGateway from the specified map of raw messages.
func UnmarshalTransitGateway(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TransitGateway)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "connections", &obj.Connections, UnmarshalTransitGatewayConnection)
	if err != nil {
		err = core.SDKErrorf(err, "", "connections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
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

// TransitGatewayConnection : A connection to an IBM Transit Gateway.
type TransitGatewayConnection struct {
	// The autogenerated name for this connection.
	Name *string `json:"name" validate:"required"`

	// The user-defined name of the connection created on the IBM Transit Gateway.
	TransitGatewayConnectionName *string `json:"transit_gateway_connection_name,omitempty"`

	// Determines the state of the connection.
	Status *string `json:"status" validate:"required"`

	// Local gateway IP address for the connection.
	LocalGatewayIp *string `json:"local_gateway_ip,omitempty"`

	// Remote gateway IP address for the connection.
	RemoteGatewayIp *string `json:"remote_gateway_ip,omitempty"`

	// Local tunnel IP address for the connection.
	LocalTunnelIp *string `json:"local_tunnel_ip,omitempty"`

	// Remote tunnel IP address for the connection.
	RemoteTunnelIp *string `json:"remote_tunnel_ip,omitempty"`

	// Local network BGP ASN for the connection.
	LocalBgpAsn *int64 `json:"local_bgp_asn,omitempty"`

	// Remote network BGP ASN for the connection.
	RemoteBgpAsn *int64 `json:"remote_bgp_asn,omitempty"`

	// The ID of the account that owns the connected network.
	NetworkAccountID *string `json:"network_account_id" validate:"required"`

	// The type of the network that is connected through this connection. Only "unbound_gre_tunnel" is supported.
	NetworkType *string `json:"network_type" validate:"required"`

	// The type of the network that the unbound GRE tunnel is targeting. Only "classic" is supported.
	BaseNetworkType *string `json:"base_network_type" validate:"required"`

	// The location of the connection.
	Zone *string `json:"zone" validate:"required"`
}

// Constants associated with the TransitGatewayConnection.Status property.
// Determines the state of the connection.
const (
	TransitGatewayConnection_Status_Creating = "creating"
	TransitGatewayConnection_Status_Deleting = "deleting"
	TransitGatewayConnection_Status_Detached = "detached"
	TransitGatewayConnection_Status_Pending = "pending"
	TransitGatewayConnection_Status_ReadyToUse = "ready_to_use"
)

// UnmarshalTransitGatewayConnection unmarshals an instance of TransitGatewayConnection from the specified map of raw messages.
func UnmarshalTransitGatewayConnection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TransitGatewayConnection)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "transit_gateway_connection_name", &obj.TransitGatewayConnectionName)
	if err != nil {
		err = core.SDKErrorf(err, "", "transit_gateway_connection_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "local_gateway_ip", &obj.LocalGatewayIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "local_gateway_ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remote_gateway_ip", &obj.RemoteGatewayIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "remote_gateway_ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "local_tunnel_ip", &obj.LocalTunnelIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "local_tunnel_ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remote_tunnel_ip", &obj.RemoteTunnelIp)
	if err != nil {
		err = core.SDKErrorf(err, "", "remote_tunnel_ip-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "local_bgp_asn", &obj.LocalBgpAsn)
	if err != nil {
		err = core.SDKErrorf(err, "", "local_bgp_asn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "remote_bgp_asn", &obj.RemoteBgpAsn)
	if err != nil {
		err = core.SDKErrorf(err, "", "remote_bgp_asn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "network_account_id", &obj.NetworkAccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "network_account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "network_type", &obj.NetworkType)
	if err != nil {
		err = core.SDKErrorf(err, "", "network_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "base_network_type", &obj.BaseNetworkType)
	if err != nil {
		err = core.SDKErrorf(err, "", "base_network_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone", &obj.Zone)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateCluster : Response of cluster update.
type UpdateCluster struct {
	// The cluster ID.
	ID *string `json:"id" validate:"required"`

	// The cluster name.
	Name *string `json:"name" validate:"required"`

	// The hyperlink of the cluster resource.
	Href *string `json:"href" validate:"required"`

	// The time that the cluster is ordered.
	OrderedAt *strfmt.DateTime `json:"ordered_at" validate:"required"`

	// The time that the cluster is provisioned and available to use.
	ProvisionedAt *strfmt.DateTime `json:"provisioned_at,omitempty"`

	// The number of hosts in the cluster.
	HostCount *int64 `json:"host_count" validate:"required"`

	// The status of the Cloud Director site cluster.
	Status *string `json:"status" validate:"required"`

	// The location of deployed cluster.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// Back link to associated Cloud Director site resource.
	DirectorSite *DirectorSiteReference `json:"director_site" validate:"required"`

	// The name of the host profile.
	HostProfile *string `json:"host_profile" validate:"required"`

	// The storage type of the cluster.
	StorageType *string `json:"storage_type" validate:"required"`

	// The billing plan for the cluster.
	BillingPlan *string `json:"billing_plan" validate:"required"`

	// Chosen storage policies and their sizes.
	FileShares *FileShares `json:"file_shares" validate:"required"`

	// Information of request accepted.
	Message *string `json:"message" validate:"required"`

	// ID to track the update operation of the cluster.
	OperationID *string `json:"operation_id" validate:"required"`
}

// Constants associated with the UpdateCluster.StorageType property.
// The storage type of the cluster.
const (
	UpdateCluster_StorageType_Nfs = "nfs"
)

// Constants associated with the UpdateCluster.BillingPlan property.
// The billing plan for the cluster.
const (
	UpdateCluster_BillingPlan_Monthly = "monthly"
)

// UnmarshalUpdateCluster unmarshals an instance of UpdateCluster from the specified map of raw messages.
func UnmarshalUpdateCluster(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateCluster)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_at", &obj.OrderedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "ordered_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioned_at", &obj.ProvisionedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioned_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_count", &obj.HostCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "director_site", &obj.DirectorSite, UnmarshalDirectorSiteReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_site-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "host_profile", &obj.HostProfile)
	if err != nil {
		err = core.SDKErrorf(err, "", "host_profile-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "storage_type", &obj.StorageType)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "billing_plan", &obj.BillingPlan)
	if err != nil {
		err = core.SDKErrorf(err, "", "billing_plan-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "file_shares", &obj.FileShares, UnmarshalFileShares)
	if err != nil {
		err = core.SDKErrorf(err, "", "file_shares-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operation_id", &obj.OperationID)
	if err != nil {
		err = core.SDKErrorf(err, "", "operation_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateDirectorSitesPvdcsClusterOptions : The UpdateDirectorSitesPvdcsCluster options.
type UpdateDirectorSitesPvdcsClusterOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// The cluster to query.
	ID *string `json:"id" validate:"required,ne="`

	// A unique ID for the resource pool in a Cloud Director site.
	PvdcID *string `json:"pvdc_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_director_sites_pvdcs_cluster.
	Body map[string]interface{} `json:"body" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateDirectorSitesPvdcsClusterOptions : Instantiate UpdateDirectorSitesPvdcsClusterOptions
func (*VmwareV1) NewUpdateDirectorSitesPvdcsClusterOptions(siteID string, id string, pvdcID string, body map[string]interface{}) *UpdateDirectorSitesPvdcsClusterOptions {
	return &UpdateDirectorSitesPvdcsClusterOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
		PvdcID: core.StringPtr(pvdcID),
		Body: body,
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetSiteID(siteID string) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetID(id string) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetPvdcID : Allow user to set PvdcID
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetPvdcID(pvdcID string) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.PvdcID = core.StringPtr(pvdcID)
	return _options
}

// SetBody : Allow user to set Body
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetBody(body map[string]interface{}) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.Body = body
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetAcceptLanguage(acceptLanguage string) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *UpdateDirectorSitesPvdcsClusterOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *UpdateDirectorSitesPvdcsClusterOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDirectorSitesPvdcsClusterOptions) SetHeaders(param map[string]string) *UpdateDirectorSitesPvdcsClusterOptions {
	options.Headers = param
	return options
}

// UpdateDirectorSitesVcdaC2cConnectionOptions : The UpdateDirectorSitesVcdaC2cConnection options.
type UpdateDirectorSitesVcdaC2cConnectionOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the cloud-to-cloud connections in the relationship Cloud Director site.
	ID *string `json:"id" validate:"required,ne="`

	// Note.
	Note *string `json:"note" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateDirectorSitesVcdaC2cConnectionOptions : Instantiate UpdateDirectorSitesVcdaC2cConnectionOptions
func (*VmwareV1) NewUpdateDirectorSitesVcdaC2cConnectionOptions(siteID string, id string, note string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	return &UpdateDirectorSitesVcdaC2cConnectionOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
		Note: core.StringPtr(note),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetSiteID(siteID string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetID(id string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetNote : Allow user to set Note
func (_options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetNote(note string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	_options.Note = core.StringPtr(note)
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetAcceptLanguage(acceptLanguage string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDirectorSitesVcdaC2cConnectionOptions) SetHeaders(param map[string]string) *UpdateDirectorSitesVcdaC2cConnectionOptions {
	options.Headers = param
	return options
}

// UpdateDirectorSitesVcdaConnectionEndpointsOptions : The UpdateDirectorSitesVcdaConnectionEndpoints options.
type UpdateDirectorSitesVcdaConnectionEndpointsOptions struct {
	// A unique ID for the Cloud Director site in which the virtual data center was created.
	SiteID *string `json:"site_id" validate:"required,ne="`

	// A unique ID for the VCDA connections in the relationship Cloud Director site.
	ID *string `json:"id" validate:"required,ne="`

	// List of allowed IP addresses.
	AllowList []string `json:"allow_list,omitempty"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Transaction ID.
	XGlobalTransactionID *string `json:"X-Global-Transaction-ID,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateDirectorSitesVcdaConnectionEndpointsOptions : Instantiate UpdateDirectorSitesVcdaConnectionEndpointsOptions
func (*VmwareV1) NewUpdateDirectorSitesVcdaConnectionEndpointsOptions(siteID string, id string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	return &UpdateDirectorSitesVcdaConnectionEndpointsOptions{
		SiteID: core.StringPtr(siteID),
		ID: core.StringPtr(id),
	}
}

// SetSiteID : Allow user to set SiteID
func (_options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetSiteID(siteID string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.SiteID = core.StringPtr(siteID)
	return _options
}

// SetID : Allow user to set ID
func (_options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetID(id string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetAllowList : Allow user to set AllowList
func (_options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetAllowList(allowList []string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.AllowList = allowList
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetAcceptLanguage(acceptLanguage string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetXGlobalTransactionID : Allow user to set XGlobalTransactionID
func (_options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetXGlobalTransactionID(xGlobalTransactionID string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	_options.XGlobalTransactionID = core.StringPtr(xGlobalTransactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDirectorSitesVcdaConnectionEndpointsOptions) SetHeaders(param map[string]string) *UpdateDirectorSitesVcdaConnectionEndpointsOptions {
	options.Headers = param
	return options
}

// UpdateVdcOptions : The UpdateVdc options.
type UpdateVdcOptions struct {
	// A unique ID for a specified virtual data center.
	ID *string `json:"id" validate:"required,ne="`

	// JSON Merge-Patch content for update_vdc.
	VDCPatch map[string]interface{} `json:"VDC_patch" validate:"required"`

	// Language.
	AcceptLanguage *string `json:"Accept-Language,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateVdcOptions : Instantiate UpdateVdcOptions
func (*VmwareV1) NewUpdateVdcOptions(id string, vDCPatch map[string]interface{}) *UpdateVdcOptions {
	return &UpdateVdcOptions{
		ID: core.StringPtr(id),
		VDCPatch: vDCPatch,
	}
}

// SetID : Allow user to set ID
func (_options *UpdateVdcOptions) SetID(id string) *UpdateVdcOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetVDCPatch : Allow user to set VDCPatch
func (_options *UpdateVdcOptions) SetVDCPatch(vDCPatch map[string]interface{}) *UpdateVdcOptions {
	_options.VDCPatch = vDCPatch
	return _options
}

// SetAcceptLanguage : Allow user to set AcceptLanguage
func (_options *UpdateVdcOptions) SetAcceptLanguage(acceptLanguage string) *UpdateVdcOptions {
	_options.AcceptLanguage = core.StringPtr(acceptLanguage)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateVdcOptions) SetHeaders(param map[string]string) *UpdateVdcOptions {
	options.Headers = param
	return options
}

// UpdatedVcdaC2c : Updated VCDA cloud-to-cloud connection note.
type UpdatedVcdaC2c struct {
	// ID of VCDA connection on the workload domain.
	ID *string `json:"id" validate:"required"`

	// Note.
	Note *string `json:"note" validate:"required"`
}

// UnmarshalUpdatedVcdaC2c unmarshals an instance of UpdatedVcdaC2c from the specified map of raw messages.
func UnmarshalUpdatedVcdaC2c(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedVcdaC2c)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "note", &obj.Note)
	if err != nil {
		err = core.SDKErrorf(err, "", "note-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedVcdaConnection : Update private connection.
type UpdatedVcdaConnection struct {
	// ID of the VCDA connection.
	ID *string `json:"id,omitempty"`

	// Status of the VCDA connection after accepting the request.
	Status *string `json:"status,omitempty"`
}

// UnmarshalUpdatedVcdaConnection unmarshals an instance of UpdatedVcdaConnection from the specified map of raw messages.
func UnmarshalUpdatedVcdaConnection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedVcdaConnection)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VDC : A VMware virtual data center (VDC). VMware VDCs are used to deploy and run VMware virtualized networking and run
// VMware workloads. VMware VDCs form loose boundaries of networking and workload where networking and workload can be
// shared or optionally isolated between VDCs. You can deploy one or more VDCs in an instance except when you are using
// the minimal instance configuration, which consists of 2 hosts (2-Socket 32 Cores, 192 GB RAM). With the minimal
// instance configuration, you can start with just one VDC and a performance network edge of medium size until
// additional hosts are added to the cluster.
type VDC struct {
	// The URL of this virtual data center (VDC).
	Href *string `json:"href" validate:"required"`

	// A unique ID for the virtual data center (VDC).
	ID *string `json:"id" validate:"required"`

	// The time that the virtual data center (VDC) is provisioned and available to use.
	ProvisionedAt *strfmt.DateTime `json:"provisioned_at,omitempty"`

	// The vCPU usage limit on the virtual data center (VDC). Supported for VDCs deployed on a multitenant Cloud Director
	// site. This property is applicable when the resource pool type is reserved.
	Cpu *int64 `json:"cpu,omitempty"`

	// A unique ID for the virtual data center (VDC) in IBM Cloud.
	Crn *string `json:"crn" validate:"required"`

	// The time that the virtual data center (VDC) is deleted.
	DeletedAt *strfmt.DateTime `json:"deleted_at,omitempty"`

	// The Cloud Director site in which to deploy the virtual data center (VDC).
	DirectorSite *VDCDirectorSite `json:"director_site" validate:"required"`

	// The VMware NSX-T networking edges deployed on the virtual data center (VDC). NSX-T edges are used for bridging
	// virtualization networking to the physical public-internet and IBM private networking.
	Edges []Edge `json:"edges" validate:"required"`

	// Information about why the request to create the virtual data center (VDC) cannot be completed.
	StatusReasons []StatusReason `json:"status_reasons" validate:"required"`

	// A human readable ID for the virtual data center (VDC).
	Name *string `json:"name" validate:"required"`

	// The time that the virtual data center (VDC) is ordered.
	OrderedAt *strfmt.DateTime `json:"ordered_at" validate:"required"`

	// The URL of the organization that owns the VDC.
	OrgHref *string `json:"org_href" validate:"required"`

	// The name of the VMware Cloud Director organization that contains this virtual data center (VDC). VMware Cloud
	// Director organizations are used to create strong boundaries between VDCs. There is a complete isolation of user
	// administration, networking, workloads, and VMware Cloud Director catalogs between different Director organizations.
	OrgName *string `json:"org_name" validate:"required"`

	// The RAM usage limit on the virtual data center (VDC) in GB (1024^3 bytes). Supported for VDCs deployed on a
	// multitenant Cloud Director site. This property is applicable when the resource pool type is reserved.
	Ram *int64 `json:"ram,omitempty"`

	// Determines the state of the virtual data center.
	Status *string `json:"status" validate:"required"`

	// Determines whether this virtual data center is in a single-tenant or multitenant Cloud Director site.
	Type *string `json:"type" validate:"required"`

	// Determines whether this virtual data center has fast provisioning enabled or not.
	FastProvisioningEnabled *bool `json:"fast_provisioning_enabled" validate:"required"`

	// Indicates if the RHEL VMs will be using the license from IBM or the customer will use their own license (BYOL).
	RhelByol *bool `json:"rhel_byol" validate:"required"`

	// Indicates if the Microsoft Windows VMs will be using the license from IBM or the customer will use their own license
	// (BYOL).
	WindowsByol *bool `json:"windows_byol" validate:"required"`
}

// Constants associated with the VDC.Status property.
// Determines the state of the virtual data center.
const (
	VDC_Status_Creating = "creating"
	VDC_Status_Deleted = "deleted"
	VDC_Status_Deleting = "deleting"
	VDC_Status_Failed = "failed"
	VDC_Status_Modifying = "modifying"
	VDC_Status_ReadyToUse = "ready_to_use"
)

// Constants associated with the VDC.Type property.
// Determines whether this virtual data center is in a single-tenant or multitenant Cloud Director site.
const (
	VDC_Type_Multitenant = "multitenant"
	VDC_Type_SingleTenant = "single_tenant"
)

// UnmarshalVDC unmarshals an instance of VDC from the specified map of raw messages.
func UnmarshalVDC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDC)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "provisioned_at", &obj.ProvisionedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "provisioned_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cpu", &obj.Cpu)
	if err != nil {
		err = core.SDKErrorf(err, "", "cpu-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deleted_at", &obj.DeletedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "deleted_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "director_site", &obj.DirectorSite, UnmarshalVDCDirectorSite)
	if err != nil {
		err = core.SDKErrorf(err, "", "director_site-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "edges", &obj.Edges, UnmarshalEdge)
	if err != nil {
		err = core.SDKErrorf(err, "", "edges-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_reasons", &obj.StatusReasons, UnmarshalStatusReason)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_reasons-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ordered_at", &obj.OrderedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "ordered_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "org_href", &obj.OrgHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "org_href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "org_name", &obj.OrgName)
	if err != nil {
		err = core.SDKErrorf(err, "", "org_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ram", &obj.Ram)
	if err != nil {
		err = core.SDKErrorf(err, "", "ram-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fast_provisioning_enabled", &obj.FastProvisioningEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "fast_provisioning_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rhel_byol", &obj.RhelByol)
	if err != nil {
		err = core.SDKErrorf(err, "", "rhel_byol-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "windows_byol", &obj.WindowsByol)
	if err != nil {
		err = core.SDKErrorf(err, "", "windows_byol-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VDCCollection : A list of virtual data centers (VDCs).
type VDCCollection struct {
	// A list of virtual data centers (VDCs).
	Vdcs []VDC `json:"vdcs" validate:"required"`
}

// UnmarshalVDCCollection unmarshals an instance of VDCCollection from the specified map of raw messages.
func UnmarshalVDCCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCCollection)
	err = core.UnmarshalModel(m, "vdcs", &obj.Vdcs, UnmarshalVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "vdcs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VDCDirectorSite : The Cloud Director site in which to deploy the virtual data center (VDC).
type VDCDirectorSite struct {
	// A unique ID for the Cloud Director site.
	ID *string `json:"id" validate:"required"`

	// The resource pool within the Director Site in which to deploy the virtual data center (VDC).
	Pvdc *DirectorSitePVDC `json:"pvdc" validate:"required"`

	// The URL of the VMware Cloud Director tenant portal where this virtual data center (VDC) can be managed.
	URL *string `json:"url" validate:"required"`
}

// UnmarshalVDCDirectorSite unmarshals an instance of VDCDirectorSite from the specified map of raw messages.
func UnmarshalVDCDirectorSite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCDirectorSite)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pvdc", &obj.Pvdc, UnmarshalDirectorSitePVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "pvdc-error", common.GetComponentInfo())
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

// VDCDirectorSitePrototype : The Cloud Director site in which to deploy the virtual data center (VDC).
type VDCDirectorSitePrototype struct {
	// A unique ID for the Cloud Director site.
	ID *string `json:"id" validate:"required"`

	// The resource pool within the Director Site in which to deploy the virtual data center (VDC).
	Pvdc *DirectorSitePVDC `json:"pvdc" validate:"required"`
}

// NewVDCDirectorSitePrototype : Instantiate VDCDirectorSitePrototype (Generic Model Constructor)
func (*VmwareV1) NewVDCDirectorSitePrototype(id string, pvdc *DirectorSitePVDC) (_model *VDCDirectorSitePrototype, err error) {
	_model = &VDCDirectorSitePrototype{
		ID: core.StringPtr(id),
		Pvdc: pvdc,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVDCDirectorSitePrototype unmarshals an instance of VDCDirectorSitePrototype from the specified map of raw messages.
func UnmarshalVDCDirectorSitePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCDirectorSitePrototype)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pvdc", &obj.Pvdc, UnmarshalDirectorSitePVDC)
	if err != nil {
		err = core.SDKErrorf(err, "", "pvdc-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VDCEdgePrototype : The networking edge to be deployed on the virtual data center (VDC).
type VDCEdgePrototype struct {
	// The size of the edge. Only used for edges of type performance.
	Size *string `json:"size,omitempty"`

	// The type of edge to be deployed on the virtual data center (VDC).
	Type *string `json:"type" validate:"required"`

	// Indicates whether the edge is private only. The default value is True for a private Cloud Director site and False
	// for a public Cloud Director site.
	PrivateOnly *bool `json:"private_only,omitempty"`
}

// Constants associated with the VDCEdgePrototype.Size property.
// The size of the edge. Only used for edges of type performance.
const (
	VDCEdgePrototype_Size_ExtraLarge = "extra_large"
	VDCEdgePrototype_Size_Large = "large"
	VDCEdgePrototype_Size_Medium = "medium"
)

// Constants associated with the VDCEdgePrototype.Type property.
// The type of edge to be deployed on the virtual data center (VDC).
const (
	VDCEdgePrototype_Type_Efficiency = "efficiency"
	VDCEdgePrototype_Type_Performance = "performance"
)

// NewVDCEdgePrototype : Instantiate VDCEdgePrototype (Generic Model Constructor)
func (*VmwareV1) NewVDCEdgePrototype(typeVar string) (_model *VDCEdgePrototype, err error) {
	_model = &VDCEdgePrototype{
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVDCEdgePrototype unmarshals an instance of VDCEdgePrototype from the specified map of raw messages.
func UnmarshalVDCEdgePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCEdgePrototype)
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_only", &obj.PrivateOnly)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_only-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VDCPatch : Information required to update a virtual data center (VDC).
type VDCPatch struct {
	// The vCPU usage limit on the virtual data center (VDC). Supported for VDCs deployed on a multitenant Cloud Director
	// site. This property is required when the resource pool type is reserved.
	Cpu *int64 `json:"cpu,omitempty"`

	// Indicates whether to enable or not fast provisioning.
	FastProvisioningEnabled *bool `json:"fast_provisioning_enabled,omitempty"`

	// The RAM usage limit on the virtual data center (VDC) in GB (1024^3 bytes). Supported for VDCs deployed on a
	// multitenant Cloud Director site. This property is required when the resource pool type is reserved.
	Ram *int64 `json:"ram,omitempty"`
}

// UnmarshalVDCPatch unmarshals an instance of VDCPatch from the specified map of raw messages.
func UnmarshalVDCPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCPatch)
	err = core.UnmarshalPrimitive(m, "cpu", &obj.Cpu)
	if err != nil {
		err = core.SDKErrorf(err, "", "cpu-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fast_provisioning_enabled", &obj.FastProvisioningEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "fast_provisioning_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ram", &obj.Ram)
	if err != nil {
		err = core.SDKErrorf(err, "", "ram-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the VDCPatch
func (vDCPatch *VDCPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(vDCPatch.Cpu) {
		_patch["cpu"] = vDCPatch.Cpu
	}
	if !core.IsNil(vDCPatch.FastProvisioningEnabled) {
		_patch["fast_provisioning_enabled"] = vDCPatch.FastProvisioningEnabled
	}
	if !core.IsNil(vDCPatch.Ram) {
		_patch["ram"] = vDCPatch.Ram
	}

	return
}

// VDCProviderType : Determines how resources are made available to the virtual data center (VDC). Required for VDCs deployed on a
// multitenant Cloud Director site.
type VDCProviderType struct {
	// The name of the resource pool type.
	Name *string `json:"name" validate:"required"`
}

// Constants associated with the VDCProviderType.Name property.
// The name of the resource pool type.
const (
	VDCProviderType_Name_OnDemand = "on_demand"
	VDCProviderType_Name_Paygo = "paygo"
	VDCProviderType_Name_Reserved = "reserved"
)

// NewVDCProviderType : Instantiate VDCProviderType (Generic Model Constructor)
func (*VmwareV1) NewVDCProviderType(name string) (_model *VDCProviderType, err error) {
	_model = &VDCProviderType{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVDCProviderType unmarshals an instance of VDCProviderType from the specified map of raw messages.
func UnmarshalVDCProviderType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VDCProviderType)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VcdaC2c : Response to create cloud-to-cloud connection.
type VcdaC2c struct {
	// ID of VCDA connection on the workload domain.
	ID *string `json:"id" validate:"required"`

	// Status of the VCDA connection.
	Status *string `json:"status" validate:"required"`

	// The offering name of the peer site, "vmware_aas" or "vmware_shared".
	PeerOffering *string `json:"peer_offering" validate:"required"`

	// Where to deploy the cluster.
	LocalDataCenterName *string `json:"local_data_center_name" validate:"required"`

	// Site name.
	LocalSiteName *string `json:"local_site_name" validate:"required"`

	// Peer site name.
	PeerSiteName *string `json:"peer_site_name" validate:"required"`

	// Peer region.
	PeerRegion *string `json:"peer_region" validate:"required"`

	// Note.
	Note *string `json:"note" validate:"required"`
}

// Constants associated with the VcdaC2c.Status property.
// Status of the VCDA connection.
const (
	VcdaC2c_Status_Creating = "creating"
	VcdaC2c_Status_Deleted = "deleted"
	VcdaC2c_Status_Deleting = "deleting"
	VcdaC2c_Status_ReadyToUse = "ready_to_use"
	VcdaC2c_Status_Updating = "updating"
)

// UnmarshalVcdaC2c unmarshals an instance of VcdaC2c from the specified map of raw messages.
func UnmarshalVcdaC2c(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VcdaC2c)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "peer_offering", &obj.PeerOffering)
	if err != nil {
		err = core.SDKErrorf(err, "", "peer_offering-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "local_data_center_name", &obj.LocalDataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "local_data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "local_site_name", &obj.LocalSiteName)
	if err != nil {
		err = core.SDKErrorf(err, "", "local_site_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "peer_site_name", &obj.PeerSiteName)
	if err != nil {
		err = core.SDKErrorf(err, "", "peer_site_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "peer_region", &obj.PeerRegion)
	if err != nil {
		err = core.SDKErrorf(err, "", "peer_region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "note", &obj.Note)
	if err != nil {
		err = core.SDKErrorf(err, "", "note-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VcdaConnection : Created VCDA connection.
type VcdaConnection struct {
	// ID of the VCDA connection on the Cloud Director site.
	ID *string `json:"id" validate:"required"`

	// Status of the VCDA connection.
	Status *string `json:"status" validate:"required"`

	// Connection type.
	Type *string `json:"type" validate:"required"`

	// Connection speed.
	Speed *string `json:"speed" validate:"required"`

	// Where to deploy the cluster.
	DataCenterName *string `json:"data_center_name" validate:"required"`

	// List of IP addresses allowed in the public connection.
	AllowList []string `json:"allow_list" validate:"required"`
}

// Constants associated with the VcdaConnection.Status property.
// Status of the VCDA connection.
const (
	VcdaConnection_Status_Creating = "creating"
	VcdaConnection_Status_Deleted = "deleted"
	VcdaConnection_Status_Deleting = "deleting"
	VcdaConnection_Status_ReadyToUse = "ready_to_use"
	VcdaConnection_Status_Updating = "updating"
)

// Constants associated with the VcdaConnection.Type property.
// Connection type.
const (
	VcdaConnection_Type_Private = "private"
	VcdaConnection_Type_Public = "public"
)

// Constants associated with the VcdaConnection.Speed property.
// Connection speed.
const (
	VcdaConnection_Speed_Speed20g = "speed_20g"
)

// UnmarshalVcdaConnection unmarshals an instance of VcdaConnection from the specified map of raw messages.
func UnmarshalVcdaConnection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VcdaConnection)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "speed", &obj.Speed)
	if err != nil {
		err = core.SDKErrorf(err, "", "speed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center_name", &obj.DataCenterName)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_list", &obj.AllowList)
	if err != nil {
		err = core.SDKErrorf(err, "", "allow_list-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Sobr : Configuration details of the scale-out backup repository.
type Sobr struct {
	// The ID of the scale-out backup repository.
	ID *string `json:"id,omitempty"`

	// The name of the scale-out backup repository.
	Name *string `json:"name,omitempty"`

	// The size of the scale-out backup repository.
	Size *int64 `json:"size,omitempty"`

	// The data center location where to create the scale-out backup repository.
	DataCenter *string `json:"data_center,omitempty"`

	// The immutability time of the backup files stored in the scale-out backup repository.
	ImmutabilityTime *int64 `json:"immutability_time,omitempty"`

	// The type of storage for the scale-out backup repository.
	StorageType *string `json:"storage_type" validate:"required"`

	// The type of scale-out backup repository.
	Type *string `json:"type" validate:"required"`

	// The ID of the Veeam organization configuration.
	VeeamOrgConfigID *string `json:"veeam_org_config_id,omitempty"`

	// The status of the scale-out backup repository on the Veeam service instance.
	Status *string `json:"status" validate:"required"`

	// The date and time when the scale-out backup repository is ordered.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`
}

// Constants associated with the Sobr.StorageType property.
// The type of storage for the scale-out backup repository.
const (
	Sobr_StorageType_Cos = "cos"
	Sobr_StorageType_Hybrid = "hybrid"
	Sobr_StorageType_Vsan = "vsan"
)

// Constants associated with the Sobr.Type property.
// The type of scale-out backup repository.
const (
	Sobr_Type_Custom = "custom"
	Sobr_Type_Default = "default"
)

// Constants associated with the Sobr.Status property.
// The status of the scale-out backup repository on the Veeam service instance.
const (
	Sobr_Status_Creating = "creating"
	Sobr_Status_Deleted = "deleted"
	Sobr_Status_Deleting = "deleting"
	Sobr_Status_ReadyToUse = "ready_to_use"
	Sobr_Status_Updating = "updating"
)

// UnmarshalSobr unmarshals an instance of Sobr from the specified map of raw messages.
func UnmarshalSobr(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Sobr)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data_center", &obj.DataCenter)
	if err != nil {
		err = core.SDKErrorf(err, "", "data_center-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "immutability_time", &obj.ImmutabilityTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "immutability_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "storage_type", &obj.StorageType)
	if err != nil {
		err = core.SDKErrorf(err, "", "storage_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "veeam_org_config_id", &obj.VeeamOrgConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "veeam_org_config_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
