/**
 * (C) Copyright IBM Corp. 2019.
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

// Package resourceconfigurationv1 : Operations and models for the ResourceConfigurationV1 service
package resourceconfigurationv1

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/go-openapi/strfmt"
	common "github.com/IBM/ibm-cos-sdk-go-config/common"
)

// ResourceConfigurationV1 : REST API used to configure Cloud Object Storage buckets.  This version of the API only
// supports reading bucket metadata and setting IP access controls.
//
// Version: 1.0.0
type ResourceConfigurationV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://config.cloud-object-storage.cloud.ibm.com/v1"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "resource_configuration"

// ResourceConfigurationV1Options : Service options
type ResourceConfigurationV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewResourceConfigurationV1UsingExternalConfig : constructs an instance of ResourceConfigurationV1 with passed in options and external configuration.
func NewResourceConfigurationV1UsingExternalConfig(options *ResourceConfigurationV1Options) (resourceConfiguration *ResourceConfigurationV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	resourceConfiguration, err = NewResourceConfigurationV1(options)
	if err != nil {
		return
	}

	err = resourceConfiguration.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		resourceConfiguration.Service.SetServiceURL(options.URL)
	}
	return
}

// NewResourceConfigurationV1 : constructs an instance of ResourceConfigurationV1 with passed in options.
func NewResourceConfigurationV1(options *ResourceConfigurationV1Options) (service *ResourceConfigurationV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:             DefaultServiceURL,
		Authenticator:   options.Authenticator,
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		return
	}

	if options.URL != "" {
		baseService.SetServiceURL(options.URL)
	}

	service = &ResourceConfigurationV1{
		Service: baseService,
	}

	return
}

// SetServiceURL sets the service URL
func (resourceConfiguration *ResourceConfigurationV1) SetServiceURL(url string) error {
	return resourceConfiguration.Service.SetServiceURL(url)
}


// GetBucketConfig : Returns metadata for the specified bucket
// Returns metadata for the specified bucket.
func (resourceConfiguration *ResourceConfigurationV1) GetBucketConfig(getBucketConfigOptions *GetBucketConfigOptions) (result *Bucket, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBucketConfigOptions, "getBucketConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getBucketConfigOptions, "getBucketConfigOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"b"}
	pathParameters := []string{*getBucketConfigOptions.Bucket}

	builder := core.NewRequestBuilder(core.GET)
	_, err = builder.ConstructHTTPURL(resourceConfiguration.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBucketConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "GetBucketConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = resourceConfiguration.Service.Request(request, new(Bucket))
	if err == nil {
		var ok bool
		result, ok = response.Result.(*Bucket)
		if !ok {
			err = fmt.Errorf("An error occurred while processing the operation response.")
		}
	}

	return
}


// UpdateBucketConfig : Make changes to a bucket's configuration
// Updates a bucket using [JSON Merge Patch](https://tools.ietf.org/html/rfc7396). This request is used to add
// functionality (like an IP access filter) or to update existing parameters.  **Primitives are overwritten and replaced
// in their entirety. It is not possible to append a new (or to delete a specific) value to an array.**  Arrays can be
// cleared by updating the parameter with an empty array `[]`. Only updates specified mutable fields. Please don't use
// `PATCH` trying to update the number of objects in a bucket, any timestamps, or other non-mutable fields.
func (resourceConfiguration *ResourceConfigurationV1) UpdateBucketConfig(updateBucketConfigOptions *UpdateBucketConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateBucketConfigOptions, "updateBucketConfigOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateBucketConfigOptions, "updateBucketConfigOptions")
	if err != nil {
		return
	}

	pathSegments := []string{"b"}
	pathParameters := []string{*updateBucketConfigOptions.Bucket}

	builder := core.NewRequestBuilder(core.PATCH)
	_, err = builder.ConstructHTTPURL(resourceConfiguration.Service.Options.URL, pathSegments, pathParameters)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateBucketConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "UpdateBucketConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	builder.AddHeader("Content-Type", "application/json")
	if updateBucketConfigOptions.IfMatch != nil {
		builder.AddHeader("if-match", fmt.Sprint(*updateBucketConfigOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if updateBucketConfigOptions.Firewall != nil {
		body["firewall"] = updateBucketConfigOptions.Firewall
	}
	if updateBucketConfigOptions.ActivityTracking != nil {
		body["activity_tracking"] = updateBucketConfigOptions.ActivityTracking
	}
	if updateBucketConfigOptions.MetricsMonitoring != nil {
		body["metrics_monitoring"] = updateBucketConfigOptions.MetricsMonitoring
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = resourceConfiguration.Service.Request(request, nil)

	return
}


// ActivityTracking : Enables sending log data to Activity Tracker and LogDNA to provide visibility into object read and write events. All
// object events are sent to the activity tracker instance defined in the `activity_tracker_crn` field.
type ActivityTracking struct {

	// If set to `true`, all object read events (i.e. downloads) will be sent to Activity Tracker.
	ReadDataEvents *bool `json:"read_data_events,omitempty"`

	// If set to `true`, all object write events (i.e. uploads) will be sent to Activity Tracker.
	WriteDataEvents *bool `json:"write_data_events,omitempty"`

	// Required the first time `activity_tracking` is configured. The instance of Activity Tracker that will recieve object
	// event data. The format is "crn:v1:bluemix:public:logdnaat:{bucket location}:a/{storage account}:{activity tracker
	// service instance}::".
	ActivityTrackerCrn *string `json:"activity_tracker_crn,omitempty"`
}

// Bucket : A bucket.
type Bucket struct {

	// The name of the bucket. Non-mutable.
	Name *string `json:"name,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	Crn *string `json:"crn,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`

	// The service instance that holds the bucket. Non-mutable.
	ServiceInstanceCrn *string `json:"service_instance_crn,omitempty"`

	// The creation time of the bucket in RFC 3339 format. Non-mutable.
	TimeCreated *strfmt.DateTime `json:"time_created,omitempty"`

	// The modification time of the bucket in RFC 3339 format. Non-mutable.
	TimeUpdated *strfmt.DateTime `json:"time_updated,omitempty"`

	// Total number of objects in the bucket. Non-mutable.
	ObjectCount *int64 `json:"object_count,omitempty"`

	// Total size of all objects in the bucket. Non-mutable.
	BytesUsed *int64 `json:"bytes_used,omitempty"`

	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to Activity Tracker and LogDNA to provide visibility into object read and write events. All
	// object events are sent to the activity tracker instance defined in the `activity_tracker_crn` field.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring. All metrics are sent to the IBM Cloud Monitoring instance defined
	// in the `monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`
}

// Firewall : An access control mechanism based on the network (IP address) where request originated. Requests not originating from
// IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including public
// access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the requester
// to have the `manager` role.
type Firewall struct {

	// List of IPv4 or IPv6 addresses in CIDR notation to be affected by firewall in CIDR notation is supported. Passing an
	// empty array will lift the IP address filter.  The `allowed_ip` array can contain a maximum of 1000 items.
	AllowedIp []string `json:"allowed_ip"`
}

// GetBucketConfigOptions : The GetBucketConfig options.
type GetBucketConfigOptions struct {

	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// NewGetBucketConfigOptions : Instantiate GetBucketConfigOptions
func (_ *ResourceConfigurationV1) NewGetBucketConfigOptions(bucket string) *GetBucketConfigOptions {
	return &GetBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (options *GetBucketConfigOptions) SetBucket(bucket string) *GetBucketConfigOptions {
	options.Bucket = core.StringPtr(bucket)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetBucketConfigOptions) SetHeaders(param map[string]string) *GetBucketConfigOptions {
	options.Headers = param
	return options
}

// MetricsMonitoring : Enables sending metrics to IBM Cloud Monitoring. All metrics are sent to the IBM Cloud Monitoring instance defined in
// the `monitoring_crn` field.
type MetricsMonitoring struct {

	// If set to `true`, all usage metrics (i.e. `bytes_used`) will be sent to the monitoring service.
	UsageMetricsEnabled *bool `json:"usage_metrics_enabled,omitempty"`

	// If set to `true`, all request metrics (i.e. `rest.object.head`) will be sent to the monitoring service.
	RequestMetricsEnabled *bool `json:"request_metrics_enabled,omitempty"`

	// Required the first time `metrics_monitoring` is configured. The instance of IBM Cloud Monitoring that will receive
	// the bucket metrics. The format is "crn:v1:bluemix:public:logdnaat:{bucket location}:a/{storage account}:{monitoring
	// service instance}::".
	MetricsMonitoringCrn *string `json:"metrics_monitoring_crn,omitempty"`
}

// UpdateBucketConfigOptions : The UpdateBucketConfig options.
type UpdateBucketConfigOptions struct {

	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required"`

	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to Activity Tracker and LogDNA to provide visibility into object read and write events. All
	// object events are sent to the activity tracker instance defined in the `activity_tracker_crn` field.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring. All metrics are sent to the IBM Cloud Monitoring instance defined
	// in the `monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// An Etag previously returned in a header when fetching or updating a bucket's metadata. If this value does not match
	// the active Etag, the request will fail.
	IfMatch *string `json:"if-match,omitempty"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// NewUpdateBucketConfigOptions : Instantiate UpdateBucketConfigOptions
func (_ *ResourceConfigurationV1) NewUpdateBucketConfigOptions(bucket string) *UpdateBucketConfigOptions {
	return &UpdateBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (options *UpdateBucketConfigOptions) SetBucket(bucket string) *UpdateBucketConfigOptions {
	options.Bucket = core.StringPtr(bucket)
	return options
}

// SetFirewall : Allow user to set Firewall
func (options *UpdateBucketConfigOptions) SetFirewall(firewall *Firewall) *UpdateBucketConfigOptions {
	options.Firewall = firewall
	return options
}

// SetActivityTracking : Allow user to set ActivityTracking
func (options *UpdateBucketConfigOptions) SetActivityTracking(activityTracking *ActivityTracking) *UpdateBucketConfigOptions {
	options.ActivityTracking = activityTracking
	return options
}

// SetMetricsMonitoring : Allow user to set MetricsMonitoring
func (options *UpdateBucketConfigOptions) SetMetricsMonitoring(metricsMonitoring *MetricsMonitoring) *UpdateBucketConfigOptions {
	options.MetricsMonitoring = metricsMonitoring
	return options
}

// SetIfMatch : Allow user to set IfMatch
func (options *UpdateBucketConfigOptions) SetIfMatch(ifMatch string) *UpdateBucketConfigOptions {
	options.IfMatch = core.StringPtr(ifMatch)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBucketConfigOptions) SetHeaders(param map[string]string) *UpdateBucketConfigOptions {
	options.Headers = param
	return options
}
