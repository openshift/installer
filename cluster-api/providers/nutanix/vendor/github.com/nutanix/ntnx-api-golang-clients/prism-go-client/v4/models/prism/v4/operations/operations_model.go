/*
 * Generated file models/prism/v4/operations/operations_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix Prism Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module prism.v4.operations of Nutanix Prism Versioned APIs
*/
package operations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/response"
	import3 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/error"
	"time"
)

/*
The batch request has been accepted for processing.
*/
type ActionType int

const (
	ACTIONTYPE_UNKNOWN  ActionType = 0
	ACTIONTYPE_REDACTED ActionType = 1
	ACTIONTYPE_CREATE   ActionType = 2
	ACTIONTYPE_MODIFY   ActionType = 3
	ACTIONTYPE_ACTION   ActionType = 4
	ACTIONTYPE_DELETE   ActionType = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ActionType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CREATE",
		"MODIFY",
		"ACTION",
		"DELETE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ActionType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CREATE",
		"MODIFY",
		"ACTION",
		"DELETE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ActionType) index(name string) ActionType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CREATE",
		"MODIFY",
		"ACTION",
		"DELETE",
	}
	for idx := range names {
		if names[idx] == name {
			return ActionType(idx)
		}
	}
	return ACTIONTYPE_UNKNOWN
}

func (e *ActionType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ActionType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ActionType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ActionType) Ref() *ActionType {
	return &e
}

/*
A model that represents a Batch resource.
*/
type Batch struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	CompletionStatus *BatchCompletionStatus `json:"completionStatus,omitempty"`
	/*
	  The completion time of the batch. The value will be in extended ISO-8601 format. For example, end time of 2022-04-23T01:23:45.678+09:00 would imply that the batch completed its execution at 1:23:45.678  on the 23rd of April 2022. Details around ISO-8601 format can be found at https://www.iso.org/standard/70907.html
	*/
	EndTime *time.Time `json:"endTime,omitempty"`

	ExecutionStatus *BatchExecutionStatus `json:"executionStatus,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The total number of elements that failed to be processed in the batch.
	*/
	FailedCount *int `json:"failedCount,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  An user friendly name of the batch.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  The total number of elements submitted for processing in the batch.
	*/
	Size *int `json:"size,omitempty"`
	/*
	  The execution start time of the batch. The value will be in extended ISO-8601 format. For example, start time of 2022-04-23T01:23:45.678+09:00 would imply that the batch started execution at 1:23:45.678  on the 23rd of April 2022. Details around ISO-8601 format can be found at https://www.iso.org/standard/70907.html
	*/
	StartTime *time.Time `json:"startTime,omitempty"`
	/*
	  The total number of elements successfully processed in the batch.
	*/
	SuccessCount *int `json:"successCount,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewBatch() *Batch {
	p := new(Batch)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.Batch"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
The completion status of the batch.
*/
type BatchCompletionStatus int

const (
	BATCHCOMPLETIONSTATUS_UNKNOWN             BatchCompletionStatus = 0
	BATCHCOMPLETIONSTATUS_REDACTED            BatchCompletionStatus = 1
	BATCHCOMPLETIONSTATUS_SUCCEEDED           BatchCompletionStatus = 2
	BATCHCOMPLETIONSTATUS_FAILED              BatchCompletionStatus = 3
	BATCHCOMPLETIONSTATUS_PARTIALLY_SUCCEEDED BatchCompletionStatus = 4
	BATCHCOMPLETIONSTATUS_CANCELLED           BatchCompletionStatus = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *BatchCompletionStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUCCEEDED",
		"FAILED",
		"PARTIALLY_SUCCEEDED",
		"CANCELLED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e BatchCompletionStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUCCEEDED",
		"FAILED",
		"PARTIALLY_SUCCEEDED",
		"CANCELLED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *BatchCompletionStatus) index(name string) BatchCompletionStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SUCCEEDED",
		"FAILED",
		"PARTIALLY_SUCCEEDED",
		"CANCELLED",
	}
	for idx := range names {
		if names[idx] == name {
			return BatchCompletionStatus(idx)
		}
	}
	return BATCHCOMPLETIONSTATUS_UNKNOWN
}

func (e *BatchCompletionStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for BatchCompletionStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *BatchCompletionStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e BatchCompletionStatus) Ref() *BatchCompletionStatus {
	return &e
}

/*
The completion status of the batch.
*/
type BatchExecutionStatus int

const (
	BATCHEXECUTIONSTATUS_UNKNOWN     BatchExecutionStatus = 0
	BATCHEXECUTIONSTATUS_REDACTED    BatchExecutionStatus = 1
	BATCHEXECUTIONSTATUS_IN_PROGRESS BatchExecutionStatus = 2
	BATCHEXECUTIONSTATUS_COMPLETED   BatchExecutionStatus = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *BatchExecutionStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IN_PROGRESS",
		"COMPLETED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e BatchExecutionStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IN_PROGRESS",
		"COMPLETED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *BatchExecutionStatus) index(name string) BatchExecutionStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IN_PROGRESS",
		"COMPLETED",
	}
	for idx := range names {
		if names[idx] == name {
			return BatchExecutionStatus(idx)
		}
	}
	return BATCHEXECUTIONSTATUS_UNKNOWN
}

func (e *BatchExecutionStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for BatchExecutionStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *BatchExecutionStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e BatchExecutionStatus) Ref() *BatchExecutionStatus {
	return &e
}

/*
The input specification for performing the batch operation.
*/
type BatchSpec struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Metadata *BatchSpecMetadata `json:"metadata,omitempty"`

	Payload []BatchSpecPayload `json:"payload,omitempty"`
}

func NewBatchSpec() *BatchSpec {
	p := new(BatchSpec)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpec"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
The metadata section on the input specification for performing the batch operation.
*/
type BatchSpecMetadata struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Action *ActionType `json:"action"`
	/*
	  The chunk size to use during the batching operation. If not specified a minimum value of 1 would be chosen.
	*/
	ChunkSize *int `json:"chunkSize,omitempty"`
	/*
	  An user friendly name of the batch.
	*/
	Name *string `json:"name"`
	/*
	  A flag indicating whether the batch procession should halt or continue when an error response is received from the server during the execution of a batch chunk
	*/
	StopOnError *bool `json:"stopOnError,omitempty"`
	/*
	  The absolute URI of the API operation on which batching will be performed.
	*/
	Uri *string `json:"uri"`
}

func (p *BatchSpecMetadata) MarshalJSON() ([]byte, error) {
	type BatchSpecMetadataProxy BatchSpecMetadata
	return json.Marshal(struct {
		*BatchSpecMetadataProxy
		Action *ActionType `json:"action,omitempty"`
		Name   *string     `json:"name,omitempty"`
		Uri    *string     `json:"uri,omitempty"`
	}{
		BatchSpecMetadataProxy: (*BatchSpecMetadataProxy)(p),
		Action:                 p.Action,
		Name:                   p.Name,
		Uri:                    p.Uri,
	})
}

func NewBatchSpecMetadata() *BatchSpecMetadata {
	p := new(BatchSpecMetadata)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpecMetadata"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	p.ChunkSize = new(int)
	*p.ChunkSize = 1

	return p
}

/*
The specification corresponding to the actual payload provided as an input to the batch operation.
*/
type BatchSpecPayload struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The data section of the payload provided to the batch operation.
	*/
	Data map[string]interface{} `json:"data,omitempty"`

	Metadata *BatchSpecPayloadMetadata `json:"metadata,omitempty"`
}

func NewBatchSpecPayload() *BatchSpecPayload {
	p := new(BatchSpecPayload)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpecPayload"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
The metadata section on the input specification for performing the batch operation.
*/
type BatchSpecPayloadMetadata struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Headers []BatchSpecPayloadMetadataHeader `json:"headers,omitempty"`

	Path []BatchSpecPayloadMetadataPath `json:"path,omitempty"`
}

func NewBatchSpecPayloadMetadata() *BatchSpecPayloadMetadata {
	p := new(BatchSpecPayloadMetadata)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpecPayloadMetadata"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
The metadata section on the input specification for performing the batch operation.
*/
type BatchSpecPayloadMetadataHeader struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The name of the header parameter.
	*/
	Name *string `json:"name"`
	/*
	  The value of the header parameter.
	*/
	Value *string `json:"value"`
}

func (p *BatchSpecPayloadMetadataHeader) MarshalJSON() ([]byte, error) {
	type BatchSpecPayloadMetadataHeaderProxy BatchSpecPayloadMetadataHeader
	return json.Marshal(struct {
		*BatchSpecPayloadMetadataHeaderProxy
		Name  *string `json:"name,omitempty"`
		Value *string `json:"value,omitempty"`
	}{
		BatchSpecPayloadMetadataHeaderProxy: (*BatchSpecPayloadMetadataHeaderProxy)(p),
		Name:                                p.Name,
		Value:                               p.Value,
	})
}

func NewBatchSpecPayloadMetadataHeader() *BatchSpecPayloadMetadataHeader {
	p := new(BatchSpecPayloadMetadataHeader)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpecPayloadMetadataHeader"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
The metadata section on the input specification for performing the batch operation.
*/
type BatchSpecPayloadMetadataPath struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The name of the  path parameter.
	*/
	Name *string `json:"name"`
	/*
	  The value of the  path parameter.
	*/
	Value *string `json:"value"`
}

func (p *BatchSpecPayloadMetadataPath) MarshalJSON() ([]byte, error) {
	type BatchSpecPayloadMetadataPathProxy BatchSpecPayloadMetadataPath
	return json.Marshal(struct {
		*BatchSpecPayloadMetadataPathProxy
		Name  *string `json:"name,omitempty"`
		Value *string `json:"value,omitempty"`
	}{
		BatchSpecPayloadMetadataPathProxy: (*BatchSpecPayloadMetadataPathProxy)(p),
		Name:                              p.Name,
		Value:                             p.Value,
	})
}

func NewBatchSpecPayloadMetadataPath() *BatchSpecPayloadMetadataPath {
	p := new(BatchSpecPayloadMetadataPath)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.BatchSpecPayloadMetadataPath"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /prism/v4.0.b1/operations/batches/{extId} Get operation
*/
type GetBatchApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetBatchApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetBatchApiResponse() *GetBatchApiResponse {
	p := new(GetBatchApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.GetBatchApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetBatchApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetBatchApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetBatchApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /prism/v4.0.b1/operations/batches Get operation
*/
type ListBatchesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListBatchesApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListBatchesApiResponse() *ListBatchesApiResponse {
	p := new(ListBatchesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.ListBatchesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListBatchesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListBatchesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListBatchesApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /prism/v4.0.b1/operations/$actions/batch Post operation
*/
type SubmitBatchApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfSubmitBatchApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewSubmitBatchApiResponse() *SubmitBatchApiResponse {
	p := new(SubmitBatchApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.operations.SubmitBatchApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *SubmitBatchApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *SubmitBatchApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfSubmitBatchApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

type OneOfGetBatchApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *Batch                 `json:"-"`
}

func NewOneOfGetBatchApiResponseData() *OneOfGetBatchApiResponseData {
	p := new(OneOfGetBatchApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetBatchApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetBatchApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case Batch:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(Batch)
		}
		*p.oneOfType0 = v.(Batch)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetBatchApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfGetBatchApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(Batch)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.operations.Batch" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(Batch)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetBatchApiResponseData"))
}

func (p *OneOfGetBatchApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfGetBatchApiResponseData")
}

type OneOfListBatchesApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    []Batch                `json:"-"`
}

func NewOneOfListBatchesApiResponseData() *OneOfListBatchesApiResponseData {
	p := new(OneOfListBatchesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListBatchesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListBatchesApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []Batch:
		p.oneOfType0 = v.([]Batch)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<prism.v4.operations.Batch>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<prism.v4.operations.Batch>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListBatchesApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<prism.v4.operations.Batch>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfListBatchesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new([]Batch)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "prism.v4.operations.Batch" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<prism.v4.operations.Batch>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<prism.v4.operations.Batch>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListBatchesApiResponseData"))
}

func (p *OneOfListBatchesApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<prism.v4.operations.Batch>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfListBatchesApiResponseData")
}

type OneOfSubmitBatchApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import3.TaskReference `json:"-"`
}

func NewOneOfSubmitBatchApiResponseData() *OneOfSubmitBatchApiResponseData {
	p := new(OneOfSubmitBatchApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfSubmitBatchApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfSubmitBatchApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import3.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import3.TaskReference)
		}
		*p.oneOfType0 = v.(import3.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfSubmitBatchApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfSubmitBatchApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import3.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import3.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfSubmitBatchApiResponseData"))
}

func (p *OneOfSubmitBatchApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfSubmitBatchApiResponseData")
}

type FileDetail struct {
	Path        *string `json:"-"`
	ObjectType_ *string `json:"-"`
}

func NewFileDetail() *FileDetail {
	p := new(FileDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "FileDetail"

	return p
}
