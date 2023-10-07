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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ResourceReviewClient is the client of the 'resource_review' resource.
//
// Manages resource review.
type ResourceReviewClient struct {
	transport http.RoundTripper
	path      string
}

// NewResourceReviewClient creates a new client for the 'resource_review'
// resource using the given transport to send the requests and receive the
// responses.
func NewResourceReviewClient(transport http.RoundTripper, path string) *ResourceReviewClient {
	return &ResourceReviewClient{
		transport: transport,
		path:      path,
	}
}

// Post creates a request for the 'post' method.
//
// Returns the list of identifiers of the resources that an account can
// perform the specified action upon.
func (c *ResourceReviewClient) Post() *ResourceReviewPostRequest {
	return &ResourceReviewPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ResourceReviewPostRequest is the request for the 'post' method.
type ResourceReviewPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	request   *ResourceReviewRequest
}

// Parameter adds a query parameter.
func (r *ResourceReviewPostRequest) Parameter(name string, value interface{}) *ResourceReviewPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ResourceReviewPostRequest) Header(name string, value interface{}) *ResourceReviewPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ResourceReviewPostRequest) Impersonate(user string) *ResourceReviewPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Request sets the value of the 'request' parameter.
func (r *ResourceReviewPostRequest) Request(value *ResourceReviewRequest) *ResourceReviewPostRequest {
	r.request = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ResourceReviewPostRequest) Send() (result *ResourceReviewPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ResourceReviewPostRequest) SendContext(ctx context.Context) (result *ResourceReviewPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeResourceReviewPostRequest(r, buffer)
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
	result = &ResourceReviewPostResponse{}
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
	err = readResourceReviewPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ResourceReviewPostResponse is the response for the 'post' method.
type ResourceReviewPostResponse struct {
	status int
	header http.Header
	err    *errors.Error
	review *ResourceReview
}

// Status returns the response status code.
func (r *ResourceReviewPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ResourceReviewPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ResourceReviewPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Review returns the value of the 'review' parameter.
func (r *ResourceReviewPostResponse) Review() *ResourceReview {
	if r == nil {
		return nil
	}
	return r.review
}

// GetReview returns the value of the 'review' parameter and
// a flag indicating if the parameter has a value.
func (r *ResourceReviewPostResponse) GetReview() (value *ResourceReview, ok bool) {
	ok = r != nil && r.review != nil
	if ok {
		value = r.review
	}
	return
}
