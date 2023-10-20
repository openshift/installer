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

package v1 // github.com/openshift-online/ocm-sdk-go/jobqueue/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// JobClient is the client of the 'job' resource.
//
// Manages status of jobs on a job queue.
type JobClient struct {
	transport http.RoundTripper
	path      string
}

// NewJobClient creates a new client for the 'job'
// resource using the given transport to send the requests and receive the
// responses.
func NewJobClient(transport http.RoundTripper, path string) *JobClient {
	return &JobClient{
		transport: transport,
		path:      path,
	}
}

// Failure creates a request for the 'failure' method.
//
// Mark a job as Failed. This method returns '204 No Content'
func (c *JobClient) Failure() *JobFailureRequest {
	return &JobFailureRequest{
		transport: c.transport,
		path:      path.Join(c.path, "failure"),
	}
}

// Success creates a request for the 'success' method.
//
// Mark a job as Successful. This method returns '204 No Content'
func (c *JobClient) Success() *JobSuccessRequest {
	return &JobSuccessRequest{
		transport: c.transport,
		path:      path.Join(c.path, "success"),
	}
}

// JobFailureRequest is the request for the 'failure' method.
type JobFailureRequest struct {
	transport     http.RoundTripper
	path          string
	query         url.Values
	header        http.Header
	failureReason *string
	receiptId     *string
}

// Parameter adds a query parameter.
func (r *JobFailureRequest) Parameter(name string, value interface{}) *JobFailureRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *JobFailureRequest) Header(name string, value interface{}) *JobFailureRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *JobFailureRequest) Impersonate(user string) *JobFailureRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// FailureReason sets the value of the 'failure_reason' parameter.
func (r *JobFailureRequest) FailureReason(value string) *JobFailureRequest {
	r.failureReason = &value
	return r
}

// ReceiptId sets the value of the 'receipt_id' parameter.
//
// A unique ID of a pop'ed job
func (r *JobFailureRequest) ReceiptId(value string) *JobFailureRequest {
	r.receiptId = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *JobFailureRequest) Send() (result *JobFailureResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *JobFailureRequest) SendContext(ctx context.Context) (result *JobFailureResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeJobFailureRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &JobFailureResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// JobFailureResponse is the response for the 'failure' method.
type JobFailureResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *JobFailureResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *JobFailureResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *JobFailureResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// JobSuccessRequest is the request for the 'success' method.
type JobSuccessRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	receiptId *string
}

// Parameter adds a query parameter.
func (r *JobSuccessRequest) Parameter(name string, value interface{}) *JobSuccessRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *JobSuccessRequest) Header(name string, value interface{}) *JobSuccessRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *JobSuccessRequest) Impersonate(user string) *JobSuccessRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// ReceiptId sets the value of the 'receipt_id' parameter.
//
// A unique ID of a pop'ed job
func (r *JobSuccessRequest) ReceiptId(value string) *JobSuccessRequest {
	r.receiptId = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *JobSuccessRequest) Send() (result *JobSuccessResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *JobSuccessRequest) SendContext(ctx context.Context) (result *JobSuccessResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeJobSuccessRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &JobSuccessResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// JobSuccessResponse is the response for the 'success' method.
type JobSuccessResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *JobSuccessResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *JobSuccessResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *JobSuccessResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}
