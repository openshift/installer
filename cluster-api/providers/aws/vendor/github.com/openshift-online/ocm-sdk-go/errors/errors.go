/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package errors // github.com/openshift-online/ocm-sdk-go/errors

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// Error kind is the name of the type used to represent errors.
const ErrorKind = "Error"

// ErrorNilKind is the name of the type used to nil errors.
const ErrorNilKind = "ErrorNil"

// ErrorBuilder is a builder for the error type.
type ErrorBuilder struct {
	bitmap_     uint32
	status      int
	id          string
	href        string
	code        string
	reason      string
	details     interface{}
	operationID string
}

// Error represents errors.
type Error struct {
	bitmap_     uint32
	status      int
	id          string
	href        string
	code        string
	reason      string
	details     interface{}
	operationID string
}

// NewError creates a new builder that can then be used to create error objects.
func NewError() *ErrorBuilder {
	return &ErrorBuilder{}
}

// Status sets the HTTP status code.
func (b *ErrorBuilder) Status(value int) *ErrorBuilder {
	b.status = value
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the error.
func (b *ErrorBuilder) ID(value string) *ErrorBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link of the error.
func (b *ErrorBuilder) HREF(value string) *ErrorBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Code sets the code of the error.
func (b *ErrorBuilder) Code(value string) *ErrorBuilder {
	b.code = value
	b.bitmap_ |= 8
	return b
}

// Reason sets the reason of the error.
func (b *ErrorBuilder) Reason(value string) *ErrorBuilder {
	b.reason = value
	b.bitmap_ |= 16
	return b
}

// OperationID sets the identifier of the operation that caused the error.
func (b *ErrorBuilder) OperationID(value string) *ErrorBuilder {
	b.operationID = value
	b.bitmap_ |= 32
	return b
}

// Details sets additional details of the error.
func (b *ErrorBuilder) Details(value interface{}) *ErrorBuilder {
	b.details = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given error into this
// builder, discarding any previous values.
func (b *ErrorBuilder) Copy(object *Error) *ErrorBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.status = object.status
	b.id = object.id
	b.href = object.href
	b.code = object.code
	b.reason = object.reason
	b.details = object.details
	b.operationID = object.operationID
	return b
}

// Build uses the information stored in the builder to create a new error object.
func (b *ErrorBuilder) Build() (result *Error, err error) {
	result = &Error{
		status:      b.status,
		id:          b.id,
		href:        b.href,
		code:        b.code,
		reason:      b.reason,
		details:     b.details,
		operationID: b.operationID,
		bitmap_:     b.bitmap_,
	}
	return
}

// Kind returns the name of the type of the error.
func (e *Error) Kind() string {
	if e == nil {
		return ErrorNilKind
	}
	return ErrorKind
}

// Status returns the HTTP status code.
func (e *Error) Status() int {
	if e != nil && e.bitmap_&1 != 0 {
		return e.status
	}
	return 0
}

// GetStatus returns the HTTP status code of the error and a flag indicating
// if the status has a value.
func (e *Error) GetStatus() (value int, ok bool) {
	ok = e != nil && e.bitmap_&1 != 0
	if ok {
		value = e.status
	}
	return
}

// ID returns the identifier of the error.
func (e *Error) ID() string {
	if e != nil && e.bitmap_&2 != 0 {
		return e.id
	}
	return ""
}

// GetID returns the identifier of the error and a flag indicating if the
// identifier has a value.
func (e *Error) GetID() (value string, ok bool) {
	ok = e != nil && e.bitmap_&2 != 0
	if ok {
		value = e.id
	}
	return
}

// HREF returns the link to the error.
func (e *Error) HREF() string {
	if e != nil && e.bitmap_&4 != 0 {
		return e.href
	}
	return ""
}

// GetHREF returns the link of the error and a flag indicating if the
// link has a value.
func (e *Error) GetHREF() (value string, ok bool) {
	ok = e != nil && e.bitmap_&4 != 0
	if ok {
		value = e.href
	}
	return
}

// Code returns the code of the error.
func (e *Error) Code() string {
	if e != nil && e.bitmap_&8 != 0 {
		return e.code
	}
	return ""
}

// GetCode returns the link of the error and a flag indicating if the
// code has a value.
func (e *Error) GetCode() (value string, ok bool) {
	ok = e != nil && e.bitmap_&8 != 0
	if ok {
		value = e.code
	}
	return
}

// Reason returns the reason of the error.
func (e *Error) Reason() string {
	if e != nil && e.bitmap_&16 != 0 {
		return e.reason
	}
	return ""
}

// GetReason returns the link of the error and a flag indicating if the
// reason has a value.
func (e *Error) GetReason() (value string, ok bool) {
	ok = e != nil && e.bitmap_&16 != 0
	if ok {
		value = e.reason
	}
	return
}

// OperationID returns the identifier of the operation that caused the error.
func (e *Error) OperationID() string {
	if e != nil && e.bitmap_&32 != 0 {
		return e.operationID
	}
	return ""
}

// GetOperationID returns the identifier of the operation that caused the error and
// a flag indicating if that identifier does have a value.
func (e *Error) GetOperationID() (value string, ok bool) {
	ok = e != nil && e.bitmap_&32 != 0
	if ok {
		value = e.operationID
	}
	return
}

// Details returns the details of the error
func (e *Error) Details() interface{} {
	if e != nil && e.bitmap_&64 != 0 {
		return e.details
	}
	return nil
}

// GetDetails returns the details of the error and a flag
// indicating if the details have a value.
func (e *Error) GetDetails() (value interface{}, ok bool) {
	ok = e != nil && e.bitmap_&64 != 0
	if ok {
		value = e.details
	}
	return
}

// Error is the implementation of the error interface.
func (e *Error) Error() string {
	chunks := make([]string, 0, 3)
	if e.bitmap_&1 != 0 {
		chunks = append(chunks, fmt.Sprintf("status is %d", e.status))
	}
	if e.bitmap_&2 != 0 {
		chunks = append(chunks, fmt.Sprintf("identifier is '%s'", e.id))
	}
	if e.bitmap_&8 != 0 {
		chunks = append(chunks, fmt.Sprintf("code is '%s'", e.code))
	}
	if e.bitmap_&32 != 0 {
		chunks = append(chunks, fmt.Sprintf("operation identifier is '%s'", e.operationID))
	}
	var result string
	size := len(chunks)
	if size == 1 {
		result = chunks[0]
	} else if size > 1 {
		result = strings.Join(chunks[0:size-1], ", ") + " and " + chunks[size-1]
	}
	if e.bitmap_&16 != 0 {
		if result != "" {
			result = result + ": "
		}
		result = result + e.reason
	}
	if result == "" {
		result = "unknown error"
	}
	return result
}

// String returns a string representing the error.
func (e *Error) String() string {
	return e.Error()
}

// UnmarshalError reads an error from the given source which can be an slice of
// bytes, a string, a reader or a JSON decoder.
func UnmarshalError(source interface{}) (object *Error, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readError(iterator)
	err = iterator.Error
	return
}

// UnmarshalErrorStatus reads an error from the given source and sets
// the given status code.
func UnmarshalErrorStatus(source interface{}, status int) (object *Error, err error) {
	object, err = UnmarshalError(source)
	if err != nil {
		return
	}
	object.status = status
	object.bitmap_ |= 1
	return
}
func readError(iterator *jsoniter.Iterator) *Error {
	object := &Error{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "status":
			object.status = iterator.ReadInt()
			object.bitmap_ |= 1
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "code":
			object.code = iterator.ReadString()
			object.bitmap_ |= 8
		case "reason":
			object.reason = iterator.ReadString()
			object.bitmap_ |= 16
		case "operation_id":
			object.operationID = iterator.ReadString()
			object.bitmap_ |= 32
		case "details":
			object.details = iterator.ReadAny().GetInterface()
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}

// MarshalError writes an error to the given writer.
func MarshalError(e *Error, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeError(e, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}
func writeError(e *Error, stream *jsoniter.Stream) {
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	stream.WriteString(ErrorKind)
	if e.bitmap_&1 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("status")
		stream.WriteInt(e.status)
	}
	if e.bitmap_&2 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("id")
		stream.WriteString(e.id)
	}
	if e.bitmap_&4 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("href")
		stream.WriteString(e.href)
	}
	if e.bitmap_&8 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("code")
		stream.WriteString(e.code)
	}
	if e.bitmap_&16 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("reason")
		stream.WriteString(e.reason)
	}
	if e.bitmap_&32 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("operation_id")
		stream.WriteString(e.operationID)
	}
	if e.bitmap_&64 != 0 {
		stream.WriteMore()
		stream.WriteObjectField("details")
		stream.WriteVal(e.details)
	}
	stream.WriteObjectEnd()
}

var panicID = "1000"
var panicError, _ = NewError().
	ID(panicID).
	Reason("An unexpected error happened, please check the log of the service " +
		"for details").
	Build()

// SendError writes a given error and status code to a response writer.
// if an error occurred it will log the error and exit.
// This methods is used internaly and no backwards compatibily is guaranteed.
func SendError(w http.ResponseWriter, r *http.Request, object *Error) {
	status, err := strconv.Atoi(object.ID())
	if err != nil {
		SendPanic(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err = MarshalError(object, w)
	if err != nil {
		glog.Errorf("Can't send response body for request '%s'", r.URL.Path)
		return
	}
}

// SendPanic sends a panic error response to the client, but it doesn't end the process.
// This methods is used internaly and no backwards compatibily is guaranteed.
func SendPanic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := MarshalError(panicError, w)
	if err != nil {
		glog.Errorf(
			"Can't send panic response for request '%s': %s",
			r.URL.Path,
			err.Error(),
		)
	}
}

// SendNotFound sends a generic 404 error.
func SendNotFound(w http.ResponseWriter, r *http.Request) {
	reason := fmt.Sprintf(
		"Can't find resource for path '%s'",
		r.URL.Path,
	)
	body, err := NewError().
		ID("404").
		Reason(reason).
		Build()
	if err != nil {
		SendPanic(w, r)
		return
	}
	SendError(w, r, body)
}

// SendMethodNotAllowed sends a generic 405 error.
func SendMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	reason := fmt.Sprintf(
		"Method '%s' isn't supported for path '%s'",
		r.Method, r.URL.Path,
	)
	body, err := NewError().
		ID("405").
		Reason(reason).
		Build()
	if err != nil {
		SendPanic(w, r)
		return
	}
	SendError(w, r, body)
}

// SendInternalServerError sends a generic 500 error.
func SendInternalServerError(w http.ResponseWriter, r *http.Request) {
	reason := fmt.Sprintf(
		"Can't process '%s' request for path '%s' due to an internal"+
			"server error",
		r.Method, r.URL.Path,
	)
	body, err := NewError().
		ID("500").
		Reason(reason).
		Build()
	if err != nil {
		SendPanic(w, r)
		return
	}
	SendError(w, r, body)
}
