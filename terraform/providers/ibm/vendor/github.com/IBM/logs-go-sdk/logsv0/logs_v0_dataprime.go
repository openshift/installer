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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
	common "github.com/IBM/logs-go-sdk/common"
)

const (
	UNKNOWN_DATA_FORMAT = "error while reading event data. unknown format data=%s"
)

// QueryCallBack interface
type QueryCallBack interface {
	OnClose()                      // called when the connection is closed from the server
	OnData(*core.DetailedResponse) // called when the data is received from the server
	OnError(error)                 // called when there is an error in receiving or processing the received data
	OnKeepAlive()                  // called when the server sends empty response to keep the connection alive
}

// Query : Query
// Run dataprime query.
func (logs *LogsV0) Query(queryOptions *QueryOptions, callBack QueryCallBack) {
	logs.QueryWithContext(context.Background(), queryOptions, callBack)
}

// QueryWithContext is an alternate form of the Query method which supports a Context parameter
func (logs *LogsV0) QueryWithContext(ctx context.Context, queryOptions *QueryOptions, callBack QueryCallBack) {
	err := core.ValidateNotNil(queryOptions, "queryOptions cannot be nil")
	if err != nil {
		callBack.OnError(err)
		return
	}

	err = core.ValidateStruct(queryOptions, "queryOptions")
	if err != nil {
		callBack.OnError(err)
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logs.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logs.Service.Options.URL, `/v1/query`, nil)
	if err != nil {
		callBack.OnError(err)
		return
	}

	for headerName, headerValue := range queryOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logs", "V0", "Query")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "text/event-stream")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if queryOptions.Query != nil {
		body["query"] = queryOptions.Query
	}
	if queryOptions.Metadata != nil {
		body["metadata"] = queryOptions.Metadata
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		callBack.OnError(err)
		return
	}

	request, err := builder.Build()
	if err != nil {
		callBack.OnError(err)
		return
	}

	var rawResponse io.ReadCloser
	response, err := logs.Service.Request(request, &rawResponse)
	if err != nil {
		callBack.OnError(err)
		return
	}

	reader := bufio.NewReader(response.Result.(io.ReadCloser))

	queryListener := &QueryListener{
		closed:   make(chan bool, 1),
		callback: callBack,
	}

	go queryListener.readEventLoop(ctx, reader, response)

	queryListener.OnClose()
}

// QueryListener
type QueryListener struct {
	closed   chan bool
	callback QueryCallBack
}

// OnClose blocks on closed chan and calls callback's OnClose()
func (queryListener *QueryListener) OnClose() {
	<-queryListener.closed
	queryListener.callback.OnClose()
}

// OnError calls callback's OnError()
func (queryListener *QueryListener) OnError(err error) {
	queryListener.callback.OnError(err)
	queryListener.closed <- true
}

// hasPrefix
func hasPrefix(s []byte, prefix string) bool {
	return bytes.HasPrefix(s, []byte(prefix))
}

// decodeJSON
func decodeJSON(data []byte, result interface{}) (interface{}, error) {
	decodeErr := json.NewDecoder(bytes.NewReader(data)).Decode(result)
	if decodeErr != nil {
		// Error decoding the response body.
		// Return the response body in RawResult, along with an error.
		return nil, decodeErr
	}

	value := reflect.ValueOf(result).Elem().Interface()
	return value, nil
}

// readEventLoop reads and processes the event data
func (queryListener *QueryListener) readEventLoop(ctx context.Context, reader *bufio.Reader, response *core.DetailedResponse) {
	var buf bytes.Buffer

	for {
		select {
		case <-ctx.Done():
			queryListener.closed <- true
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				queryListener.closed <- true
				return
			}

			if err != nil {
				queryListener.OnError(err)
				return
			}

			// stream received from the dataprime API will be in the following format
			// : success
			// data: {"query_id":{"query_id":"009e11a2-78fb-4fa9-bb00-2cec624e7855"}}
			switch {
			// handles ": success" message
			// do nothing
			case hasPrefix(line, ":"):

			// handles "data: {}" message
			case hasPrefix(line, "data: "):
				buf.Write(bytes.TrimPrefix(line, []byte("data: ")))

			// end of event
			case bytes.Equal(line, []byte("\n")):
				// there will be empty response sent from the server for keeping the connection alive
				// this response will only have ":" in the output. not handling this case will lead to
				// error when we try to decode as JSON
				if buf.Len() == 0 {
					queryListener.callback.OnKeepAlive()
					continue
				}

				b := buf.Bytes()

				// we have to decode the JSON to "map[string]json.RawMessage" type
				// so that we can use "core.UnmarshalModel" to unmarshal the response to QueryResponseStreamItem type
				// QueryResponseStreamItem has a field which is type of "interface{}" so the normal unmarshalling will fail.
				// core package has UnmarshalModel function which supports unmarshalling struts which containes fields which has interface{} type
				var rawResponse map[string]json.RawMessage
				result, err := decodeJSON(b, &rawResponse)
				if err != nil {
					queryListener.OnError(err)
					return
				}

				var queryResponse *QueryResponseStreamItem
				if result != nil {
					err := core.UnmarshalModel(result, "", &queryResponse, UnmarshalQueryResponseStreamItem)
					if err != nil {
						queryListener.OnError(err)
						return
					}
				}
				response.Result = queryResponse
				response.RawResult = b
				buf.Reset()
				queryListener.callback.OnData(response)

			default:
				queryListener.OnError(fmt.Errorf(UNKNOWN_DATA_FORMAT, line))
				return
			}
		}
	}
}

// ApisDataprimeV1BlocksLimitWarning : ApisDataprimeV1BlocksLimitWarning struct
type ApisDataprimeV1BlocksLimitWarning struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of ApisDataprimeV1BlocksLimitWarning
func (o *ApisDataprimeV1BlocksLimitWarning) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of ApisDataprimeV1BlocksLimitWarning
func (o *ApisDataprimeV1BlocksLimitWarning) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of ApisDataprimeV1BlocksLimitWarning
func (o *ApisDataprimeV1BlocksLimitWarning) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of ApisDataprimeV1BlocksLimitWarning
func (o *ApisDataprimeV1BlocksLimitWarning) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of ApisDataprimeV1BlocksLimitWarning
func (o *ApisDataprimeV1BlocksLimitWarning) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalApisDataprimeV1BlocksLimitWarning unmarshals an instance of ApisDataprimeV1BlocksLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1BlocksLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1BlocksLimitWarning)
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1BytesScannedLimitWarning : ApisDataprimeV1BytesScannedLimitWarning struct
type ApisDataprimeV1BytesScannedLimitWarning struct {

	// Allows users to set arbitrary properties
	additionalProperties map[string]interface{}
}

// SetProperty allows the user to set an arbitrary property on an instance of ApisDataprimeV1BytesScannedLimitWarning
func (o *ApisDataprimeV1BytesScannedLimitWarning) SetProperty(key string, value interface{}) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]interface{})
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of ApisDataprimeV1BytesScannedLimitWarning
func (o *ApisDataprimeV1BytesScannedLimitWarning) SetProperties(m map[string]interface{}) {
	o.additionalProperties = make(map[string]interface{})
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of ApisDataprimeV1BytesScannedLimitWarning
func (o *ApisDataprimeV1BytesScannedLimitWarning) GetProperty(key string) interface{} {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of ApisDataprimeV1BytesScannedLimitWarning
func (o *ApisDataprimeV1BytesScannedLimitWarning) GetProperties() map[string]interface{} {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of ApisDataprimeV1BytesScannedLimitWarning
func (o *ApisDataprimeV1BytesScannedLimitWarning) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	return
}

// UnmarshalApisDataprimeV1BytesScannedLimitWarning unmarshals an instance of ApisDataprimeV1BytesScannedLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1BytesScannedLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1BytesScannedLimitWarning)
	for k := range m {
		var v interface{}
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = e
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1CompileWarning : ApisDataprimeV1CompileWarning struct
type ApisDataprimeV1CompileWarning struct {
	WarningMessage *string `json:"warning_message,omitempty"`
}

// UnmarshalApisDataprimeV1CompileWarning unmarshals an instance of ApisDataprimeV1CompileWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1CompileWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1CompileWarning)
	err = core.UnmarshalPrimitive(m, "warning_message", &obj.WarningMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeError : ApisDataprimeV1DataprimeError struct
type ApisDataprimeV1DataprimeError struct {
	Message *string `json:"message,omitempty"`
}

// UnmarshalApisDataprimeV1DataprimeError unmarshals an instance of ApisDataprimeV1DataprimeError from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeError)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeResult : ApisDataprimeV1DataprimeResult struct
type ApisDataprimeV1DataprimeResult struct {
	Results []ApisDataprimeV1DataprimeResults `json:"results,omitempty"`
}

// UnmarshalApisDataprimeV1DataprimeResult unmarshals an instance of ApisDataprimeV1DataprimeResult from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeResult)
	err = core.UnmarshalModel(m, "results", &obj.Results, UnmarshalApisDataprimeV1DataprimeResults)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeResults : ApisDataprimeV1DataprimeResults struct
type ApisDataprimeV1DataprimeResults struct {
	Metadata []ApisDataprimeV1DataprimeResultsKeyValue `json:"metadata,omitempty"`

	Labels []ApisDataprimeV1DataprimeResultsKeyValue `json:"labels,omitempty"`

	UserData *string `json:"user_data,omitempty"`
}

// UnmarshalApisDataprimeV1DataprimeResults unmarshals an instance of ApisDataprimeV1DataprimeResults from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeResults(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeResults)
	err = core.UnmarshalModel(m, "metadata", &obj.Metadata, UnmarshalApisDataprimeV1DataprimeResultsKeyValue)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "labels", &obj.Labels, UnmarshalApisDataprimeV1DataprimeResultsKeyValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "user_data", &obj.UserData)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeResultsKeyValue : ApisDataprimeV1DataprimeResultsKeyValue struct
type ApisDataprimeV1DataprimeResultsKeyValue struct {
	Key *string `json:"key,omitempty"`

	Value *string `json:"value,omitempty"`
}

// UnmarshalApisDataprimeV1DataprimeResultsKeyValue unmarshals an instance of ApisDataprimeV1DataprimeResultsKeyValue from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeResultsKeyValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeResultsKeyValue)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarning : ApisDataprimeV1DataprimeWarning struct
// Models which "extend" this model:
// - ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning
// - ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning
// - ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning
// - ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning
// - ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning
// - ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning
type ApisDataprimeV1DataprimeWarning struct {
	CompileWarning *ApisDataprimeV1CompileWarning `json:"compile_warning,omitempty"`

	TimeRangeWarning *ApisDataprimeV1TimeRangeWarning `json:"time_range_warning,omitempty"`

	NumberOfResultsLimitWarning *ApisDataprimeV1NumberOfResultsLimitWarning `json:"number_of_results_limit_warning,omitempty"`

	BytesScannedLimitWarning *ApisDataprimeV1BytesScannedLimitWarning `json:"bytes_scanned_limit_warning,omitempty"`

	DeprecationWarning *ApisDataprimeV1DeprecationWarning `json:"deprecation_warning,omitempty"`

	BlocksLimitWarning *ApisDataprimeV1BlocksLimitWarning `json:"blocks_limit_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

type ApisDataprimeV1DataprimeWarningIntf interface {
	isaApisDataprimeV1DataprimeWarning() bool
}

// UnmarshalApisDataprimeV1DataprimeWarning unmarshals an instance of ApisDataprimeV1DataprimeWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarning)
	err = core.UnmarshalModel(m, "compile_warning", &obj.CompileWarning, UnmarshalApisDataprimeV1CompileWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "time_range_warning", &obj.TimeRangeWarning, UnmarshalApisDataprimeV1TimeRangeWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "number_of_results_limit_warning", &obj.NumberOfResultsLimitWarning, UnmarshalApisDataprimeV1NumberOfResultsLimitWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bytes_scanned_limit_warning", &obj.BytesScannedLimitWarning, UnmarshalApisDataprimeV1BytesScannedLimitWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "deprecation_warning", &obj.DeprecationWarning, UnmarshalApisDataprimeV1DeprecationWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "blocks_limit_warning", &obj.BlocksLimitWarning, UnmarshalApisDataprimeV1BlocksLimitWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DeprecationWarning : ApisDataprimeV1DeprecationWarning struct
type ApisDataprimeV1DeprecationWarning struct {
	WarningMessage *string `json:"warning_message,omitempty"`
}

// UnmarshalApisDataprimeV1DeprecationWarning unmarshals an instance of ApisDataprimeV1DeprecationWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DeprecationWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DeprecationWarning)
	err = core.UnmarshalPrimitive(m, "warning_message", &obj.WarningMessage)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1Metadata : ApisDataprimeV1Metadata struct
type ApisDataprimeV1Metadata struct {
	StartDate *strfmt.DateTime `json:"start_date,omitempty"`

	EndDate *strfmt.DateTime `json:"end_date,omitempty"`

	DefaultSource *string `json:"default_source,omitempty"`

	Tier *string `json:"tier,omitempty"`

	Syntax *string `json:"syntax,omitempty"`

	Limit *int64 `json:"limit,omitempty"`

	StrictFieldsValidation *bool `json:"strict_fields_validation,omitempty"`
}

// Constants associated with the ApisDataprimeV1Metadata.Tier property.
const (
	ApisDataprimeV1Metadata_Tier_Archive        = "archive"
	ApisDataprimeV1Metadata_Tier_FrequentSearch = "frequent_search"
	ApisDataprimeV1Metadata_Tier_Unspecified    = "unspecified"
)

// Constants associated with the ApisDataprimeV1Metadata.Syntax property.
const (
	ApisDataprimeV1Metadata_Syntax_Dataprime   = "dataprime"
	ApisDataprimeV1Metadata_Syntax_Lucene      = "lucene"
	ApisDataprimeV1Metadata_Syntax_Unspecified = "unspecified"
)

// UnmarshalApisDataprimeV1Metadata unmarshals an instance of ApisDataprimeV1Metadata from the specified map of raw messages.
func UnmarshalApisDataprimeV1Metadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1Metadata)
	err = core.UnmarshalPrimitive(m, "start_date", &obj.StartDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_date", &obj.EndDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_source", &obj.DefaultSource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tier", &obj.Tier)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "syntax", &obj.Syntax)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "strict_fields_validation", &obj.StrictFieldsValidation)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1NumberOfResultsLimitWarning : ApisDataprimeV1NumberOfResultsLimitWarning struct
type ApisDataprimeV1NumberOfResultsLimitWarning struct {
	NumberOfResultsLimit *int64 `json:"number_of_results_limit,omitempty"`
}

// UnmarshalApisDataprimeV1NumberOfResultsLimitWarning unmarshals an instance of ApisDataprimeV1NumberOfResultsLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1NumberOfResultsLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1NumberOfResultsLimitWarning)
	err = core.UnmarshalPrimitive(m, "number_of_results_limit", &obj.NumberOfResultsLimit)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1QueryID : ApisDataprimeV1QueryID struct
type ApisDataprimeV1QueryID struct {
	QueryID *string `json:"query_id,omitempty"`
}

// UnmarshalApisDataprimeV1QueryID unmarshals an instance of ApisDataprimeV1QueryID from the specified map of raw messages.
func UnmarshalApisDataprimeV1QueryID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1QueryID)
	err = core.UnmarshalPrimitive(m, "query_id", &obj.QueryID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1TimeRangeWarning : ApisDataprimeV1TimeRangeWarning struct
type ApisDataprimeV1TimeRangeWarning struct {
	WarningMessage *string `json:"warning_message,omitempty"`

	StartDate *strfmt.DateTime `json:"start_date,omitempty"`

	EndDate *strfmt.DateTime `json:"end_date,omitempty"`
}

// UnmarshalApisDataprimeV1TimeRangeWarning unmarshals an instance of ApisDataprimeV1TimeRangeWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1TimeRangeWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1TimeRangeWarning)
	err = core.UnmarshalPrimitive(m, "warning_message", &obj.WarningMessage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "start_date", &obj.StartDate)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "end_date", &obj.EndDate)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OpenapiApiErrorError : OpenapiApiErrorError struct
type OpenapiApiErrorError struct {
	Code *string `json:"code" validate:"required"`

	Message *string `json:"message,omitempty"`

	MoreInfo *string `json:"more_info,omitempty"`
}

// Constants associated with the OpenapiApiErrorError.Code property.
const (
	OpenapiApiErrorError_Code_BadRequestOrUnspecified = "bad_request_or_unspecified"
	OpenapiApiErrorError_Code_Conflict                = "conflict"
	OpenapiApiErrorError_Code_DeadlineExceeded        = "deadline_exceeded"
	OpenapiApiErrorError_Code_Forbidden               = "forbidden"
	OpenapiApiErrorError_Code_MethodInternalError     = "method_internal_error"
	OpenapiApiErrorError_Code_NotFound                = "not_found"
	OpenapiApiErrorError_Code_ResourceExhausted       = "resource_exhausted"
	OpenapiApiErrorError_Code_Unauthenticated         = "unauthenticated"
	OpenapiApiErrorError_Code_Unauthorized            = "unauthorized"
)

// UnmarshalOpenapiApiErrorError unmarshals an instance of OpenapiApiErrorError from the specified map of raw messages.
func UnmarshalOpenapiApiErrorError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OpenapiApiErrorError)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "more_info", &obj.MoreInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryOptions : The Query options.
type QueryOptions struct {
	Query *string `json:"query,omitempty"`

	Metadata *ApisDataprimeV1Metadata `json:"metadata,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewQueryOptions : Instantiate QueryOptions
func (*LogsV0) NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

// SetQuery : Allow user to set Query
func (_options *QueryOptions) SetQuery(query string) *QueryOptions {
	_options.Query = core.StringPtr(query)
	return _options
}

// SetMetadata : Allow user to set Metadata
func (_options *QueryOptions) SetMetadata(metadata *ApisDataprimeV1Metadata) *QueryOptions {
	_options.Metadata = metadata
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *QueryOptions) SetHeaders(param map[string]string) *QueryOptions {
	options.Headers = param
	return options
}

// QueryResponseStreamItem : QueryResponseStreamItem struct
// Models which "extend" this model:
// - QueryResponseStreamItemQueryResponse
// - QueryResponseStreamItemApiError
type QueryResponseStreamItem struct {
	Error *ApisDataprimeV1DataprimeError `json:"error,omitempty"`

	Result *ApisDataprimeV1DataprimeResult `json:"result,omitempty"`

	Warning ApisDataprimeV1DataprimeWarningIntf `json:"warning,omitempty"`

	QueryID *ApisDataprimeV1QueryID `json:"query_id,omitempty"`

	Errors []OpenapiApiErrorError `json:"errors,omitempty"`

	Trace *string `json:"trace,omitempty"`

	StatusCode *int64 `json:"status_code,omitempty"`
}

func (*QueryResponseStreamItem) isaQueryResponseStreamItem() bool {
	return true
}

type QueryResponseStreamItemIntf interface {
	isaQueryResponseStreamItem() bool
}

// UnmarshalQueryResponseStreamItem unmarshals an instance of QueryResponseStreamItem from the specified map of raw messages.
func UnmarshalQueryResponseStreamItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItem)
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalApisDataprimeV1DataprimeError)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalApisDataprimeV1DataprimeResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "warning", &obj.Warning, UnmarshalApisDataprimeV1DataprimeWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "query_id", &obj.QueryID, UnmarshalApisDataprimeV1QueryID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalOpenapiApiErrorError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning : ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning struct {
	BlocksLimitWarning *ApisDataprimeV1BlocksLimitWarning `json:"blocks_limit_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeBlocksLimitWarning)
	err = core.UnmarshalModel(m, "blocks_limit_warning", &obj.BlocksLimitWarning, UnmarshalApisDataprimeV1BlocksLimitWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning : ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning struct {
	BytesScannedLimitWarning *ApisDataprimeV1BytesScannedLimitWarning `json:"bytes_scanned_limit_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeBytesScannedLimitWarning)
	err = core.UnmarshalModel(m, "bytes_scanned_limit_warning", &obj.BytesScannedLimitWarning, UnmarshalApisDataprimeV1BytesScannedLimitWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning : ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning struct {
	CompileWarning *ApisDataprimeV1CompileWarning `json:"compile_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeCompileWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeCompileWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeCompileWarning)
	err = core.UnmarshalModel(m, "compile_warning", &obj.CompileWarning, UnmarshalApisDataprimeV1CompileWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning : ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning struct {
	DeprecationWarning *ApisDataprimeV1DeprecationWarning `json:"deprecation_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeDeprecationWarning)
	err = core.UnmarshalModel(m, "deprecation_warning", &obj.DeprecationWarning, UnmarshalApisDataprimeV1DeprecationWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning : ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning struct {
	NumberOfResultsLimitWarning *ApisDataprimeV1NumberOfResultsLimitWarning `json:"number_of_results_limit_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeNumberOfResultsLimitWarning)
	err = core.UnmarshalModel(m, "number_of_results_limit_warning", &obj.NumberOfResultsLimitWarning, UnmarshalApisDataprimeV1NumberOfResultsLimitWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning : ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning struct
// This model "extends" ApisDataprimeV1DataprimeWarning
type ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning struct {
	TimeRangeWarning *ApisDataprimeV1TimeRangeWarning `json:"time_range_warning,omitempty"`
}

func (*ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning) isaApisDataprimeV1DataprimeWarning() bool {
	return true
}

// UnmarshalApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning unmarshals an instance of ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning from the specified map of raw messages.
func UnmarshalApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApisDataprimeV1DataprimeWarningWarningTypeTimeRangeWarning)
	err = core.UnmarshalModel(m, "time_range_warning", &obj.TimeRangeWarning, UnmarshalApisDataprimeV1TimeRangeWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemApiError : QueryResponseStreamItemApiError struct
// This model "extends" QueryResponseStreamItem
type QueryResponseStreamItemApiError struct {
	Errors []OpenapiApiErrorError `json:"errors" validate:"required"`

	Trace *string `json:"trace" validate:"required"`

	StatusCode *int64 `json:"status_code,omitempty"`
}

func (*QueryResponseStreamItemApiError) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemApiError unmarshals an instance of QueryResponseStreamItemApiError from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemApiError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemApiError)
	err = core.UnmarshalModel(m, "errors", &obj.Errors, UnmarshalOpenapiApiErrorError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemQueryResponse : QueryResponseStreamItemQueryResponse struct
// Models which "extend" this model:
// - QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError
// - QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult
// - QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning
// - QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID
// This model "extends" QueryResponseStreamItem
type QueryResponseStreamItemQueryResponse struct {
	Error *ApisDataprimeV1DataprimeError `json:"error,omitempty"`

	Result *ApisDataprimeV1DataprimeResult `json:"result,omitempty"`

	Warning ApisDataprimeV1DataprimeWarningIntf `json:"warning,omitempty"`

	QueryID *ApisDataprimeV1QueryID `json:"query_id,omitempty"`
}

func (*QueryResponseStreamItemQueryResponse) isaQueryResponseStreamItemQueryResponse() bool {
	return true
}

type QueryResponseStreamItemQueryResponseIntf interface {
	QueryResponseStreamItemIntf
	isaQueryResponseStreamItemQueryResponse() bool
}

func (*QueryResponseStreamItemQueryResponse) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemQueryResponse unmarshals an instance of QueryResponseStreamItemQueryResponse from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemQueryResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemQueryResponse)
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalApisDataprimeV1DataprimeError)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalApisDataprimeV1DataprimeResult)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "warning", &obj.Warning, UnmarshalApisDataprimeV1DataprimeWarning)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "query_id", &obj.QueryID, UnmarshalApisDataprimeV1QueryID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError : QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError struct
// This model "extends" QueryResponseStreamItemQueryResponse
type QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError struct {
	Error *ApisDataprimeV1DataprimeError `json:"error,omitempty"`
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError) isaQueryResponseStreamItemQueryResponse() bool {
	return true
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError unmarshals an instance of QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageError)
	err = core.UnmarshalModel(m, "error", &obj.Error, UnmarshalApisDataprimeV1DataprimeError)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID : QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID struct
// This model "extends" QueryResponseStreamItemQueryResponse
type QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID struct {
	QueryID *ApisDataprimeV1QueryID `json:"query_id,omitempty"`
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID) isaQueryResponseStreamItemQueryResponse() bool {
	return true
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID unmarshals an instance of QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageQueryID)
	err = core.UnmarshalModel(m, "query_id", &obj.QueryID, UnmarshalApisDataprimeV1QueryID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult : QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult struct
// This model "extends" QueryResponseStreamItemQueryResponse
type QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult struct {
	Result *ApisDataprimeV1DataprimeResult `json:"result,omitempty"`
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult) isaQueryResponseStreamItemQueryResponse() bool {
	return true
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult unmarshals an instance of QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageResult)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalApisDataprimeV1DataprimeResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning : QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning struct
// This model "extends" QueryResponseStreamItemQueryResponse
type QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning struct {
	Warning ApisDataprimeV1DataprimeWarningIntf `json:"warning,omitempty"`
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning) isaQueryResponseStreamItemQueryResponse() bool {
	return true
}

func (*QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning) isaQueryResponseStreamItem() bool {
	return true
}

// UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning unmarshals an instance of QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning from the specified map of raw messages.
func UnmarshalQueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueryResponseStreamItemQueryResponseQueryResponseApisDataprimeV1QueryResponseMessageWarning)
	err = core.UnmarshalModel(m, "warning", &obj.Warning, UnmarshalApisDataprimeV1DataprimeWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
