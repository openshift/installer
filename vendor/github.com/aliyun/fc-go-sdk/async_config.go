package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AsyncConfig represents the async configuration.
type AsyncConfig struct {
	DestinationConfig         *DestinationConfig `json:"destinationConfig"`
	MaxAsyncEventAgeInSeconds *int64             `json:"maxAsyncEventAgeInSeconds"`
	MaxAsyncRetryAttempts     *int64             `json:"maxAsyncRetryAttempts"`
}

// AsyncConfigResponse defines the detail async config object
type AsyncConfigResponse struct {
	Service                   *string            `json:"service"`
	Function                  *string            `json:"function"`
	CreatedTime               *string            `json:"createdTime"`
	Qualifier                 *string            `json:"qualifier"`
	LastModifiedTime          *string            `json:"lastModifiedTime"`
	AsyncConfig
}

// DestinationConfig represents the destination configuration.
type DestinationConfig struct {
	OnSuccess *Destination `json:"onSuccess"`
	OnFailure *Destination `json:"onFailure"`
}

// Destination represents the destination arn.
type Destination struct {
	Destination *string `json:"destination"`
}

// PutFunctionAsyncInvokeConfigInput defines function creation input
type PutFunctionAsyncInvokeConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	AsyncConfig
}

func NewPutFunctionAsyncInvokeConfigInput(serviceName, functionName string) *PutFunctionAsyncInvokeConfigInput {
	return &PutFunctionAsyncInvokeConfigInput{ServiceName: &serviceName, FunctionName: &functionName}
}

func (i *PutFunctionAsyncInvokeConfigInput) WithQualifier(qualifier string) *PutFunctionAsyncInvokeConfigInput {
	i.Qualifier = &qualifier
	return i
}

func (i *PutFunctionAsyncInvokeConfigInput) WithAsyncConfig(config AsyncConfig) *PutFunctionAsyncInvokeConfigInput {
	i.AsyncConfig = config
	return i
}

func (i *PutFunctionAsyncInvokeConfigInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *PutFunctionAsyncInvokeConfigInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(asyncConfigWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
	}
	return fmt.Sprintf(asyncConfigPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
}

func (i *PutFunctionAsyncInvokeConfigInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *PutFunctionAsyncInvokeConfigInput) GetPayload() interface{} {
	return i.AsyncConfig
}

func (i *PutFunctionAsyncInvokeConfigInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

// PutFunctionAsyncInvokeConfigOutput define get async config response
type PutFunctionAsyncInvokeConfigOutput struct {
	Header http.Header
	AsyncConfigResponse
}

func (o PutFunctionAsyncInvokeConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o PutFunctionAsyncInvokeConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

// GetFunctionAsyncInvokeConfigInput defines function creation input
type GetFunctionAsyncInvokeConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
}

func NewGetFunctionAsyncInvokeConfigInput(serviceName, functionName string) *GetFunctionAsyncInvokeConfigInput {
	return &GetFunctionAsyncInvokeConfigInput{ServiceName: &serviceName, FunctionName: &functionName}
}

func (i *GetFunctionAsyncInvokeConfigInput) WithQualifier(qualifier string) *GetFunctionAsyncInvokeConfigInput {
	i.Qualifier = &qualifier
	return i
}

func (i *GetFunctionAsyncInvokeConfigInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *GetFunctionAsyncInvokeConfigInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(asyncConfigWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
	}
	return fmt.Sprintf(asyncConfigPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
}

func (i *GetFunctionAsyncInvokeConfigInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *GetFunctionAsyncInvokeConfigInput) GetPayload() interface{} {
	return nil
}

func (i *GetFunctionAsyncInvokeConfigInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

// GetFunctionAsyncInvokeConfigOutput define get data response
type GetFunctionAsyncInvokeConfigOutput struct {
	Header http.Header
	AsyncConfigResponse
}

func (o GetFunctionAsyncInvokeConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetFunctionAsyncInvokeConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}


// DeleteFunctionAsyncInvokeConfigInput defines function creation input
type DeleteFunctionAsyncInvokeConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
}

func NewDeleteFunctionAsyncInvokeConfigInput(serviceName, functionName string) *DeleteFunctionAsyncInvokeConfigInput {
	return &DeleteFunctionAsyncInvokeConfigInput{ServiceName: &serviceName, FunctionName: &functionName}
}

func (i *DeleteFunctionAsyncInvokeConfigInput) WithQualifier(qualifier string) *DeleteFunctionAsyncInvokeConfigInput {
	i.Qualifier = &qualifier
	return i
}

func (i *DeleteFunctionAsyncInvokeConfigInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *DeleteFunctionAsyncInvokeConfigInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(asyncConfigWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
	}
	return fmt.Sprintf(asyncConfigPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
}

func (i *DeleteFunctionAsyncInvokeConfigInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *DeleteFunctionAsyncInvokeConfigInput) GetPayload() interface{} {
	return nil
}

func (i *DeleteFunctionAsyncInvokeConfigInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

// DeleteFunctionAsyncInvokeConfigOutput define delete data response
type DeleteFunctionAsyncInvokeConfigOutput struct {
	Header http.Header
}

func (o DeleteFunctionAsyncInvokeConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o DeleteFunctionAsyncInvokeConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

// ListFunctionAsyncInvokeConfigsOutput defines ListFunctionAsyncInvokeConfigsOutput result
type ListFunctionAsyncInvokeConfigsOutput struct {
	Header    http.Header
	Configs  []*AsyncConfigResponse `json:"configs"`
	NextToken *string               `json:"nextToken,omitempty"`
}

func (o ListFunctionAsyncInvokeConfigsOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}
func (o ListFunctionAsyncInvokeConfigsOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type ListFunctionAsyncInvokeConfigsInput struct {
	ServiceName  *string `json:"serviceName"`
	FunctionName *string `json:"functionName"`
	NextToken *string
	Limit     *int32
}

func NewListFunctionAsyncInvokeConfigsInput(serviceName, functionName string) *ListFunctionAsyncInvokeConfigsInput {
	return &ListFunctionAsyncInvokeConfigsInput{
		ServiceName: &serviceName,
		FunctionName: &functionName,
	}
}

func (i *ListFunctionAsyncInvokeConfigsInput) WithNextToken(nextToken string) *ListFunctionAsyncInvokeConfigsInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListFunctionAsyncInvokeConfigsInput) WithLimit(limit int32) *ListFunctionAsyncInvokeConfigsInput {
	i.Limit = &limit
	return i
}

func (i *ListFunctionAsyncInvokeConfigsInput) GetQueryParams() url.Values {
	out := url.Values{}

	if i.NextToken != nil {
		out.Set("nextToken", *i.NextToken)
	}

	if i.Limit != nil {
		out.Set("limit", strconv.FormatInt(int64(*i.Limit), 10))
	}

	return out
}

func (i *ListFunctionAsyncInvokeConfigsInput) GetPath() string {
	return fmt.Sprintf(listAsyncConfigsPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
}

func (i *ListFunctionAsyncInvokeConfigsInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListFunctionAsyncInvokeConfigsInput) GetPayload() interface{} {
	return nil
}

func (i *ListFunctionAsyncInvokeConfigsInput) Validate() error {
	return nil
}
