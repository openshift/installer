package msgraphgocore

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	abs "github.com/microsoft/kiota-abstractions-go"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

const BatchRequestErrorRegistryKey = "BATCH_REQUEST_ERROR_REGISTRY_KEY"

// RequestHeader is a type alias for http request headers
type RequestHeader map[string]string

// Serialize serializes information the current object
func (br RequestHeader) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

// GetFieldDeserializers the deserialization information for the current model
func (br RequestHeader) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

// RequestBody is a type alias for http request bodies
type RequestBody map[string]interface{}

// Serialize serializes information the current object
func (br RequestBody) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

// GetFieldDeserializers the deserialization information for the current model
func (br RequestBody) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

type batchRequest struct {
	requests []BatchItem
}

// NewBatchRequest creates an instance of BatchRequest
func NewBatchRequest() BatchRequest {
	return &batchRequest{}
}

// BatchRequest models all the properties of a batch request
type BatchRequest interface {
	serialization.Parsable
	GetRequests() []BatchItem
	SetRequests(requests []BatchItem)
	AddBatchRequestStep(reqInfo abstractions.RequestInformation) (BatchItem, error)
	Send(ctx context.Context, adapter abstractions.RequestAdapter) (BatchResponse, error)
}

// GetRequests return all the Items in the batch request
func (br *batchRequest) GetRequests() []BatchItem {
	return br.requests
}

// SetRequests add a collection of requests to the batch Items
func (br *batchRequest) SetRequests(requests []BatchItem) {
	br.requests = requests
}

// Serialize serializes information the current object
func (br *batchRequest) Serialize(writer serialization.SerializationWriter) error {
	{
		cast := abs.CollectionApply(br.requests, func(v BatchItem) serialization.Parsable {
			return v.(serialization.Parsable)
		})
		err := writer.WriteCollectionOfObjectValues("requests", cast)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFieldDeserializers the deserialization information for the current model
func (br *batchRequest) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

// AddBatchRequestStep converts RequestInformation to a BatchItem and adds it to a BatchRequest
//
// You can add upto 20 BatchItems to a BatchRequest
func (br *batchRequest) AddBatchRequestStep(reqInfo abstractions.RequestInformation) (BatchItem, error) {
	if len(br.GetRequests()) > 19 {
		return nil, errors.New("batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := toBatchItem(reqInfo)
	if err != nil {
		return nil, err
	}

	br.SetRequests(append(br.GetRequests(), batchItem))
	return batchItem, nil
}

func toBatchItem(requestInfo abstractions.RequestInformation) (BatchItem, error) {
	uri, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(requestInfo.Content, &body)
	if err != nil {
		return nil, err
	}

	newID := uuid.NewString()
	method := requestInfo.Method.String()

	request := NewBatchItem()
	request.SetId(&newID)
	request.SetMethod(&method)
	request.SetBody(body)
	request.SetHeaders(requestInfo.Headers)
	request.SetUrl(&uri.Path)

	return request, nil
}

// Send serializes and sends the batch request to the server
func (br *batchRequest) Send(ctx context.Context, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	baseUrl, err := getBaseUrl(adapter)
	if err != nil {
		return nil, err
	}

	requestInfo, err := buildRequestInfo(ctx, adapter, br, baseUrl)
	if err != nil {
		return nil, err
	}
	return sendBatchRequest(ctx, requestInfo, adapter)
}

func getBaseUrl(adapter abstractions.RequestAdapter) (*url.URL, error) {
	return url.Parse(adapter.GetBaseUrl())
}

func buildRequestInfo(ctx context.Context, adapter abstractions.RequestAdapter, body BatchRequest, baseUrl *url.URL) (*abstractions.RequestInformation, error) {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.POST
	requestInfo.UrlTemplate = "{+baseurl}/$batch"
	requestInfo.SetUri(*baseUrl)
	err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", body)
	if err != nil {
		return nil, err
	}
	requestInfo.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	return requestInfo, nil
}

func getResponsePrimaryContentType(responseItem BatchItem) string {
	header := responseItem.GetHeaders()
	if header == nil {
		return ""
	}
	rawType := header["Content-Type"]
	splat := strings.Split(rawType, ";")
	return strings.ToLower(splat[0])
}

func getRootParseNode(responseItem BatchItem) (absser.ParseNode, error) {
	contentType := getResponsePrimaryContentType(responseItem)
	if contentType == "" {
		return nil, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(responseItem.GetBody())
	if err != nil {
		return nil, err
	}
	return serialization.DefaultParseNodeFactoryInstance.GetRootParseNode(contentType, buf.Bytes())
}

func throwErrors(responseItem BatchItem, typeName string) error {
	errorMappings := getErrorMapper(typeName)
	if errorMappings == nil {
		errorMappings = getErrorMapper(BatchRequestErrorRegistryKey)
	}
	responseStatus := *responseItem.GetStatus()

	statusAsString := strconv.Itoa(int(responseStatus))
	var errorCtor absser.ParsableFactory = nil
	if len(errorMappings) != 0 {
		if responseStatus >= 400 && responseStatus < 500 && errorMappings["4XX"] != nil {
			errorCtor = errorMappings["4XX"]
		} else if responseStatus >= 500 && responseStatus < 600 && errorMappings["5XX"] != nil {
			errorCtor = errorMappings["5XX"]
		}
	}

	if errorCtor == nil {
		return &abstractions.ApiError{
			Message: "The server returned an unexpected status code and no error factory is registered for this code: " + statusAsString,
		}
	}

	rootNode, err := getRootParseNode(responseItem)
	if err != nil {
		return err
	}
	if rootNode == nil {
		return &abstractions.ApiError{
			Message: "The server returned an unexpected status code with no response body: " + statusAsString,
		}
	}

	errValue, err := rootNode.GetObjectValue(errorCtor)
	if err != nil {
		return err
	}

	return errValue.(error)
}

// GetBatchResponseById returns the response of the batch request item with the given id.
func GetBatchResponseById[T serialization.Parsable](resp BatchResponse, itemId string) (*T, error) {
	var res T
	item := resp.GetResponseById(itemId)

	if *item.GetStatus() >= 400 {
		return nil, throwErrors(item, reflect.TypeOf(res).Name())
	}

	jsonStr, err := json.Marshal(item.GetBody())
	if err != nil {
		return &res, err
	}
	err = json.Unmarshal(jsonStr, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

func getErrorMapper(key string) abstractions.ErrorMappings {
	errorMapperSrc, found := GetErrorFactoryFromRegistry(key)
	if found {
		return errorMapperSrc
	}
	return nil
}

func sendBatchRequest(ctx context.Context, requestInfo *abstractions.RequestInformation, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}

	response, err := adapter.SendAsync(ctx, requestInfo, CreateBatchResponseDiscriminator, getErrorMapper(BatchRequestErrorRegistryKey))
	if err != nil {
		return nil, err
	}

	return response.(BatchResponse), nil
}
