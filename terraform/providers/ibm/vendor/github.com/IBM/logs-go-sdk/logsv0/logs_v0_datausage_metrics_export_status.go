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

// Package logsv0 : Operations and models for the LogsV0 service

package logsv0

import (
	"context"
	"encoding/json"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/logs-go-sdk/common"
)

// GetDataUsageMetricsExportStatusOptions : The GetDataUsageMetricsExportStatus options.
type GetDataUsageMetricsExportStatusOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDataUsageMetricsExportStatusOptions : Instantiate GetDataUsageMetricsExportStatusOptions
func (*LogsV0) NewGetDataUsageMetricsExportStatusOptions() *GetDataUsageMetricsExportStatusOptions {
	return &GetDataUsageMetricsExportStatusOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetDataUsageMetricsExportStatusOptions) SetHeaders(param map[string]string) *GetDataUsageMetricsExportStatusOptions {
	options.Headers = param
	return options
}

// GetDataUsageMetricsExportStatus : Get data usage metrics export status
// Get data usage metrics export status.
func (logs *LogsV0) GetDataUsageMetricsExportStatus(getDataUsageMetricsExportStatusOptions *GetDataUsageMetricsExportStatusOptions) (result *DataUsageMetricsExportStatus, response *core.DetailedResponse, err error) {
	result, response, err = logs.GetDataUsageMetricsExportStatusWithContext(context.Background(), getDataUsageMetricsExportStatusOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDataUsageMetricsExportStatusWithContext is an alternate form of the GetDataUsageMetricsExportStatus method which supports a Context parameter
func (logs *LogsV0) GetDataUsageMetricsExportStatusWithContext(ctx context.Context, getDataUsageMetricsExportStatusOptions *GetDataUsageMetricsExportStatusOptions) (result *DataUsageMetricsExportStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getDataUsageMetricsExportStatusOptions, "getDataUsageMetricsExportStatusOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logs.Service.Options.URL, `/v1/data_usage`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDataUsageMetricsExportStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logs", "V0", "GetDataUsageMetricsExportStatus")
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
	response, err = logs.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_data_usage_metrics_export_status", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDataUsageMetricsExportStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
