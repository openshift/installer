package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type provisionTarget struct {
	Resource *string `json:"resource"`
	Target   *int64  `json:"target"`
}

type provisionConfig struct {
	Resource *string `json:"resource"`
	Target   *int64  `json:"target"`
	Current  *int64  `json:"current"`
}

type PutProvisionConfigObject struct {
	Target *int64 `json:"target"`
}

type PutProvisionConfigInput struct {
	PutProvisionConfigObject
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	IfMatch      *string
}

func NewPutProvisionConfigInput(serviceName, qualifier, functionName string) *PutProvisionConfigInput {
	return &PutProvisionConfigInput{
		ServiceName:  &serviceName,
		Qualifier:    &qualifier,
		FunctionName: &functionName,
	}
}

func (i *PutProvisionConfigInput) WithTarget(target int64) *PutProvisionConfigInput {
	i.Target = &target
	return i
}

func (i *PutProvisionConfigInput) WithIfMatch(ifMatch string) *PutProvisionConfigInput {
	i.IfMatch = &ifMatch
	return i
}

func (i *PutProvisionConfigInput) GetPath() string {
	return fmt.Sprintf(provisionConfigWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
}

func (i *PutProvisionConfigInput) GetHeaders() Header {
	header := make(Header)
	if i.IfMatch != nil {
		header[ifMatch] = *i.IfMatch
	}
	return header
}

func (i *PutProvisionConfigInput) GetPayload() interface{} {
	return i.PutProvisionConfigObject
}

func (i *PutProvisionConfigInput) GetQueryParams() url.Values {
	return url.Values{}
}

func (i *PutProvisionConfigInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.Qualifier) {
		return fmt.Errorf("Qualifier is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

type PutProvisionConfigOutput struct {
	Header http.Header
	provisionTarget
}

func (o PutProvisionConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o PutProvisionConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o PutProvisionConfigOutput) GetEtag() string {
	return GetEtag(o.Header)
}

type GetProvisionConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
}

func NewGetProvisionConfigInput(serviceName, qualifier, functionName string) *GetProvisionConfigInput {
	return &GetProvisionConfigInput{
		ServiceName:  &serviceName,
		Qualifier:    &qualifier,
		FunctionName: &functionName,
	}
}

func (i *GetProvisionConfigInput) GetPath() string {
	return fmt.Sprintf(provisionConfigWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
}

func (i *GetProvisionConfigInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *GetProvisionConfigInput) GetPayload() interface{} {
	return nil
}

func (i *GetProvisionConfigInput) GetQueryParams() url.Values {
	return url.Values{}
}

func (i *GetProvisionConfigInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.Qualifier) {
		return fmt.Errorf("Qualifier is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

type GetProvisionConfigOutput struct {
	Header http.Header
	provisionConfig
}

func (o GetProvisionConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetProvisionConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o GetProvisionConfigOutput) GetEtag() string {
	return GetEtag(o.Header)
}

type ListProvisionConfigsInput struct {
	ServiceName *string `json:"serviceName"`
	Qualifier   *string `json:"qualifier"`
	NextToken   *string `json:"nextToken"`
	Limit       *int32  `json:"limit"`
}

func NewListProvisionConfigsInput() *ListProvisionConfigsInput {
	return &ListProvisionConfigsInput{}
}

func (i *ListProvisionConfigsInput) WithServiceName(serviceName string) *ListProvisionConfigsInput {
	i.ServiceName = &serviceName
	return i
}

func (i *ListProvisionConfigsInput) WithQualifier(qualifier string) *ListProvisionConfigsInput {
	i.Qualifier = &qualifier
	return i
}

func (i *ListProvisionConfigsInput) WithNextToken(nextToken string) *ListProvisionConfigsInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListProvisionConfigsInput) WithLimit(limit int32) *ListProvisionConfigsInput {
	i.Limit = &limit
	return i
}

func (i *ListProvisionConfigsInput) GetQueryParams() url.Values {
	out := url.Values{}
	if i.ServiceName != nil {
		out.Set("serviceName", *i.ServiceName)
	}

	if i.Qualifier != nil {
		out.Set("qualifier", *i.Qualifier)
	}

	if i.NextToken != nil {
		out.Set("nextToken", *i.NextToken)
	}

	if i.Limit != nil {
		out.Set("limit", strconv.FormatInt(int64(*i.Limit), 10))
	}

	return out
}

func (i *ListProvisionConfigsInput) GetPath() string {
	return provisionConfigPath
}

func (i *ListProvisionConfigsInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListProvisionConfigsInput) GetPayload() interface{} {
	return nil
}

func (i *ListProvisionConfigsInput) Validate() error {
	if IsBlank(i.ServiceName) && !IsBlank(i.Qualifier) {
		return fmt.Errorf("Service name is required if qualifier is provided")
	}
	return nil
}

type ListProvisionConfigsOutput struct {
	Header           http.Header
	ProvisionConfigs []*provisionConfig `json:"provisionConfigs"`
	NextToken        *string            `json:"nextToken,omitempty"`
}

func (o ListProvisionConfigsOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListProvisionConfigsOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}
