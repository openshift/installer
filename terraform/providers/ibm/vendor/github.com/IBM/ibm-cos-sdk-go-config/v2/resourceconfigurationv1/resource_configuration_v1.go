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
 * IBM OpenAPI SDK Code Generator Version: 3.91.0-d9755c53-20240605-153412
 */

// Package resourceconfigurationv1 : Operations and models for the ResourceConfigurationV1 service
package resourceconfigurationv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/ibm-cos-sdk-go-config/v2/common"
	"github.com/go-openapi/strfmt"
)

// ResourceConfigurationV1 : REST API used to configure Cloud Object Storage buckets.
//
// API Version: 1.0.0
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
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	resourceConfiguration, err = NewResourceConfigurationV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = resourceConfiguration.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = resourceConfiguration.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewResourceConfigurationV1 : constructs an instance of ResourceConfigurationV1 with passed in options.
func NewResourceConfigurationV1(options *ResourceConfigurationV1Options) (service *ResourceConfigurationV1, err error) {
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

	service = &ResourceConfigurationV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "resourceConfiguration" suitable for processing requests.
func (resourceConfiguration *ResourceConfigurationV1) Clone() *ResourceConfigurationV1 {
	if core.IsNil(resourceConfiguration) {
		return nil
	}
	clone := *resourceConfiguration
	clone.Service = resourceConfiguration.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (resourceConfiguration *ResourceConfigurationV1) SetServiceURL(url string) error {
	err := resourceConfiguration.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (resourceConfiguration *ResourceConfigurationV1) GetServiceURL() string {
	return resourceConfiguration.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (resourceConfiguration *ResourceConfigurationV1) SetDefaultHeaders(headers http.Header) {
	resourceConfiguration.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (resourceConfiguration *ResourceConfigurationV1) SetEnableGzipCompression(enableGzip bool) {
	resourceConfiguration.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (resourceConfiguration *ResourceConfigurationV1) GetEnableGzipCompression() bool {
	return resourceConfiguration.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (resourceConfiguration *ResourceConfigurationV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	resourceConfiguration.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (resourceConfiguration *ResourceConfigurationV1) DisableRetries() {
	resourceConfiguration.Service.DisableRetries()
}

// GetBucketConfig : Returns metadata for the specified bucket
// Returns metadata for the specified bucket.
func (resourceConfiguration *ResourceConfigurationV1) GetBucketConfig(getBucketConfigOptions *GetBucketConfigOptions) (result *Bucket, response *core.DetailedResponse, err error) {
	result, response, err = resourceConfiguration.GetBucketConfigWithContext(context.Background(), getBucketConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBucketConfigWithContext is an alternate form of the GetBucketConfig method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) GetBucketConfigWithContext(ctx context.Context, getBucketConfigOptions *GetBucketConfigOptions) (result *Bucket, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBucketConfigOptions, "getBucketConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBucketConfigOptions, "getBucketConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *getBucketConfigOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/b/{bucket}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = resourceConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "getBucketConfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBucket)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateBucketConfig : Make changes to a bucket's configuration
// Updates a bucket using [JSON Merge Patch](https://tools.ietf.org/html/rfc7396). This request is used to add
// functionality (like an IP access filter) or to update existing parameters.  **Primitives are overwritten and replaced
// in their entirety. It is not possible to append a new (or to delete a specific) value to an array.**  Arrays can be
// cleared by updating the parameter with an empty array `[]`. A `PATCH` operation only updates specified mutable
// fields. Please don't use `PATCH` trying to update the number of objects in a bucket, any timestamps, or other
// non-mutable fields.
func (resourceConfiguration *ResourceConfigurationV1) UpdateBucketConfig(updateBucketConfigOptions *UpdateBucketConfigOptions) (response *core.DetailedResponse, err error) {
	response, err = resourceConfiguration.UpdateBucketConfigWithContext(context.Background(), updateBucketConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateBucketConfigWithContext is an alternate form of the UpdateBucketConfig method which supports a Context parameter
func (resourceConfiguration *ResourceConfigurationV1) UpdateBucketConfigWithContext(ctx context.Context, updateBucketConfigOptions *UpdateBucketConfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateBucketConfigOptions, "updateBucketConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateBucketConfigOptions, "updateBucketConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"bucket": *updateBucketConfigOptions.Bucket,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = resourceConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(resourceConfiguration.Service.Options.URL, `/b/{bucket}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateBucketConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("resource_configuration", "V1", "UpdateBucketConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateBucketConfigOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateBucketConfigOptions.IfMatch))
	}

	if updateBucketConfigOptions.BucketPatch != nil {
		_, err = builder.SetBodyContentJSON(updateBucketConfigOptions.BucketPatch)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = resourceConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "updateBucketConfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// ActivityTracking : Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events are
// sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in the
// Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
type ActivityTracking struct {
	// If set to `true`, all object read events (i.e. downloads) will be sent to Activity Tracker Event Routing.
	ReadDataEvents *bool `json:"read_data_events,omitempty"`

	// If set to `true`, all object write events (i.e. uploads) will be sent to Activity Tracker Event Routing.
	WriteDataEvents *bool `json:"write_data_events,omitempty"`

	// When the `activity_tracker_crn` is not populated, then enabled events are sent to the Activity Tracker Event Routing
	// instance associated to the container's location unless otherwise specified in the Activity Tracker Event Routing
	// Event Routing service configuration. If `activity_tracker_crn` is populated, then enabled events are sent to the
	// Activity Tracker Event Routing instance specified and bucket management events are always enabled.
	ActivityTrackerCrn *string `json:"activity_tracker_crn,omitempty"`

	// This field only applies if `activity_tracker_crn` is not populated. If set to `true`, all bucket management events
	// will be sent to Activity Tracker Event Routing.
	ManagementEvents *bool `json:"management_events,omitempty"`
}

// UnmarshalActivityTracking unmarshals an instance of ActivityTracking from the specified map of raw messages.
func UnmarshalActivityTracking(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActivityTracking)
	err = core.UnmarshalPrimitive(m, "read_data_events", &obj.ReadDataEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "read_data_events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "write_data_events", &obj.WriteDataEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "write_data_events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "activity_tracker_crn", &obj.ActivityTrackerCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracker_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "management_events", &obj.ManagementEvents)
	if err != nil {
		err = core.SDKErrorf(err, "", "management_events-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
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

	// Number of non-current object versions in the bucket. Non-mutable.
	NoncurrentObjectCount *int64 `json:"noncurrent_object_count,omitempty"`

	// Total size of all non-current object versions in the bucket. Non-mutable.
	NoncurrentBytesUsed *int64 `json:"noncurrent_bytes_used,omitempty"`

	// Total number of delete markers in the bucket. Non-mutable.
	DeleteMarkerCount *int64 `json:"delete_marker_count,omitempty"`

	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
	// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events
	// are sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in
	// the Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
	// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
	// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
	// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
	// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
	// `metrics_monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// Maximum bytes for this bucket.
	HardQuota *int64 `json:"hard_quota,omitempty"`

	// Data structure holding protection management response.
	ProtectionManagement *ProtectionManagementResponse `json:"protection_management,omitempty"`
}

// UnmarshalBucket unmarshals an instance of Bucket from the specified map of raw messages.
func UnmarshalBucket(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Bucket)
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
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_crn", &obj.ServiceInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_created", &obj.TimeCreated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "time_updated", &obj.TimeUpdated)
	if err != nil {
		err = core.SDKErrorf(err, "", "time_updated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "object_count", &obj.ObjectCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "object_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bytes_used", &obj.BytesUsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "bytes_used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "noncurrent_object_count", &obj.NoncurrentObjectCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "noncurrent_object_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "noncurrent_bytes_used", &obj.NoncurrentBytesUsed)
	if err != nil {
		err = core.SDKErrorf(err, "", "noncurrent_bytes_used-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "delete_marker_count", &obj.DeleteMarkerCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "delete_marker_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "firewall", &obj.Firewall, UnmarshalFirewall)
	if err != nil {
		err = core.SDKErrorf(err, "", "firewall-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hard_quota", &obj.HardQuota)
	if err != nil {
		err = core.SDKErrorf(err, "", "hard_quota-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "protection_management", &obj.ProtectionManagement, UnmarshalProtectionManagementResponse)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BucketPatch : An object containing new bucket metadata.
type BucketPatch struct {
	// An access control mechanism based on the network (IP address) where request originated. Requests not originating
	// from IP addresses listed in the `allowed_ip` field will be denied regardless of any access policies (including
	// public access) that might otherwise permit the request.  Viewing or updating the `Firewall` element requires the
	// requester to have the `manager` role.
	Firewall *Firewall `json:"firewall,omitempty"`

	// Enables sending log data to IBM Cloud Activity Tracker Event Routing to provide visibility into bucket management,
	// object read and write events. (Recommended) When the `activity_tracker_crn` is not populated, then enabled events
	// are sent to the Activity Tracker Event Routing instance at the container's location unless otherwise specified in
	// the Activity Tracker Event Routing Event Routing service configuration. (Legacy) When the `activity_tracker_crn` is
	// populated, then enabled events are sent to the Activity Tracker Event Routing instance specified.
	ActivityTracking *ActivityTracking `json:"activity_tracking,omitempty"`

	// Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
	// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
	// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
	// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
	// `metrics_monitoring_crn` field.
	MetricsMonitoring *MetricsMonitoring `json:"metrics_monitoring,omitempty"`

	// Maximum bytes for this bucket.
	HardQuota *int64 `json:"hard_quota,omitempty"`

	// Data structure holding protection management operations.
	ProtectionManagement *ProtectionManagement `json:"protection_management,omitempty"`
}

// UnmarshalBucketPatch unmarshals an instance of BucketPatch from the specified map of raw messages.
func UnmarshalBucketPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BucketPatch)
	err = core.UnmarshalModel(m, "firewall", &obj.Firewall, UnmarshalFirewall)
	if err != nil {
		err = core.SDKErrorf(err, "", "firewall-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "activity_tracking", &obj.ActivityTracking, UnmarshalActivityTracking)
	if err != nil {
		err = core.SDKErrorf(err, "", "activity_tracking-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "metrics_monitoring", &obj.MetricsMonitoring, UnmarshalMetricsMonitoring)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hard_quota", &obj.HardQuota)
	if err != nil {
		err = core.SDKErrorf(err, "", "hard_quota-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "protection_management", &obj.ProtectionManagement, UnmarshalProtectionManagement)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the BucketPatch
func (bucketPatch *BucketPatch) AsPatch() (_patch map[string]interface{}, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(bucketPatch)
	if err == nil {
		err = json.Unmarshal(jsonData, &_patch)
	}
	if err != nil {
		err = core.SDKErrorf(err, "", "unmarshal-patch-data-error", common.GetComponentInfo())
	}
	return
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


// UnmarshalFirewall unmarshals an instance of Firewall from the specified map of raw messages.
func UnmarshalFirewall(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Firewall)
	err = core.UnmarshalPrimitive(m, "allowed_ip", &obj.AllowedIp)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetBucketConfigOptions : The GetBucketConfig options.
type GetBucketConfigOptions struct {
	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBucketConfigOptions : Instantiate GetBucketConfigOptions
func (*ResourceConfigurationV1) NewGetBucketConfigOptions(bucket string) *GetBucketConfigOptions {
	return &GetBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *GetBucketConfigOptions) SetBucket(bucket string) *GetBucketConfigOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBucketConfigOptions) SetHeaders(param map[string]string) *GetBucketConfigOptions {
	options.Headers = param
	return options
}

// MetricsMonitoring : Enables sending metrics to IBM Cloud Monitoring.  All metrics are opt-in. (Recommended) When the
// `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the Monitoring instance at the
// container's location unless otherwise specified in the Metrics Router service configuration. (Legacy) When the
// `metrics_monitoring_crn` is populated, then enabled metrics are sent to the Monitoring instance defined in the
// `metrics_monitoring_crn` field.
type MetricsMonitoring struct {
	// If set to `true`, all usage metrics (i.e. `bytes_used`) will be sent to the monitoring service.
	UsageMetricsEnabled *bool `json:"usage_metrics_enabled,omitempty"`

	// If set to `true`, all request metrics (i.e. `rest.object.head`) will be sent to the monitoring service.
	RequestMetricsEnabled *bool `json:"request_metrics_enabled,omitempty"`

	// When the `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the monitoring instance
	// associated to the container's location unless otherwise specified in the Metrics Router service configuration. If
	// `metrics_monitoring_crn` is populated, then enabled events are sent to the Metrics Monitoring instance specified.
	MetricsMonitoringCrn *string `json:"metrics_monitoring_crn,omitempty"`
}

// UnmarshalMetricsMonitoring unmarshals an instance of MetricsMonitoring from the specified map of raw messages.
func UnmarshalMetricsMonitoring(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MetricsMonitoring)
	err = core.UnmarshalPrimitive(m, "usage_metrics_enabled", &obj.UsageMetricsEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "usage_metrics_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "request_metrics_enabled", &obj.RequestMetricsEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "request_metrics_enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metrics_monitoring_crn", &obj.MetricsMonitoringCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "metrics_monitoring_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProtectionManagement : Data structure holding protection management operations.
type ProtectionManagement struct {
	// If set to `activate`, protection management action on the bucket is being activated.
	RequestedState *string `json:"requested_state,omitempty"`

	// This field is required when using requested_state\:`activate` and holds a JWT that is provided by the Cloud
	// Operator. This should be the encoded JWT.
	ProtectionManagementToken *string `json:"protection_management_token,omitempty"`
}

// Constants associated with the ProtectionManagement.RequestedState property.
// If set to `activate`, protection management action on the bucket is being activated.
const (
	ProtectionManagement_RequestedState_Activate = "activate"
	ProtectionManagement_RequestedState_Deactivate = "deactivate"
)

// UnmarshalProtectionManagement unmarshals an instance of ProtectionManagement from the specified map of raw messages.
func UnmarshalProtectionManagement(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagement)
	err = core.UnmarshalPrimitive(m, "requested_state", &obj.RequestedState)
	if err != nil {
		err = core.SDKErrorf(err, "", "requested_state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "protection_management_token", &obj.ProtectionManagementToken)
	if err != nil {
		err = core.SDKErrorf(err, "", "protection_management_token-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProtectionManagementResponse : Data structure holding protection management response.
type ProtectionManagementResponse struct {
	// Indicates the X number of protection management tokens that have been applied to the bucket in its lifetime.
	TokenAppliedCounter *string `json:"token_applied_counter,omitempty"`

	// The 'protection management token list' holding a recent list of applied tokens. This list may contain a subset of
	// all tokens applied to the bucket, as indicated by the counter.
	TokenEntries []ProtectionManagementResponseTokenEntry `json:"token_entries,omitempty"`
}

// UnmarshalProtectionManagementResponse unmarshals an instance of ProtectionManagementResponse from the specified map of raw messages.
func UnmarshalProtectionManagementResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagementResponse)
	err = core.UnmarshalPrimitive(m, "token_applied_counter", &obj.TokenAppliedCounter)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_applied_counter-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "token_entries", &obj.TokenEntries, UnmarshalProtectionManagementResponseTokenEntry)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_entries-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProtectionManagementResponseTokenEntry : Data structure holding protection management token.
type ProtectionManagementResponseTokenEntry struct {
	TokenID *string `json:"token_id,omitempty"`

	TokenExpirationTime *string `json:"token_expiration_time,omitempty"`

	TokenReferenceID *string `json:"token_reference_id,omitempty"`

	AppliedTime *string `json:"applied_time,omitempty"`

	InvalidatedTime *string `json:"invalidated_time,omitempty"`

	ExpirationTime *string `json:"expiration_time,omitempty"`

	ShortenRetentionFlag *bool `json:"shorten_retention_flag,omitempty"`
}

// UnmarshalProtectionManagementResponseTokenEntry unmarshals an instance of ProtectionManagementResponseTokenEntry from the specified map of raw messages.
func UnmarshalProtectionManagementResponseTokenEntry(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProtectionManagementResponseTokenEntry)
	err = core.UnmarshalPrimitive(m, "token_id", &obj.TokenID)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "token_expiration_time", &obj.TokenExpirationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_expiration_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "token_reference_id", &obj.TokenReferenceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "token_reference_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "applied_time", &obj.AppliedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "applied_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "invalidated_time", &obj.InvalidatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "invalidated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiration_time", &obj.ExpirationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiration_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "shorten_retention_flag", &obj.ShortenRetentionFlag)
	if err != nil {
		err = core.SDKErrorf(err, "", "shorten_retention_flag-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateBucketConfigOptions : The UpdateBucketConfig options.
type UpdateBucketConfigOptions struct {
	// Name of a bucket.
	Bucket *string `json:"bucket" validate:"required,ne="`

	// An object containing new configuration metadata.
	BucketPatch map[string]interface{} `json:"Bucket_patch,omitempty"`

	// An Etag previously returned in a header when fetching or updating a bucket's metadata. If this value does not match
	// the active Etag, the request will fail.
	IfMatch *string `json:"If-Match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateBucketConfigOptions : Instantiate UpdateBucketConfigOptions
func (*ResourceConfigurationV1) NewUpdateBucketConfigOptions(bucket string) *UpdateBucketConfigOptions {
	return &UpdateBucketConfigOptions{
		Bucket: core.StringPtr(bucket),
	}
}

// SetBucket : Allow user to set Bucket
func (_options *UpdateBucketConfigOptions) SetBucket(bucket string) *UpdateBucketConfigOptions {
	_options.Bucket = core.StringPtr(bucket)
	return _options
}

// SetBucketPatch : Allow user to set BucketPatch
func (_options *UpdateBucketConfigOptions) SetBucketPatch(bucketPatch map[string]interface{}) *UpdateBucketConfigOptions {
	_options.BucketPatch = bucketPatch
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateBucketConfigOptions) SetIfMatch(ifMatch string) *UpdateBucketConfigOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBucketConfigOptions) SetHeaders(param map[string]string) *UpdateBucketConfigOptions {
	options.Headers = param
	return options
}
