/*
Copyright (c) 2019 Red Hat, Inc.

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

// This file contains the implementation of the methods.

package sdk

import (
	"net/http"
)

// Get creates an HTTP GET request. Note that this request won't be sent till the Send method is
// called.
func (c *Connection) Get() *Request {
	request := new(Request)
	request.transport = c
	request.method = http.MethodGet
	return request
}

// Post creates an HTTP POST request. Note that this request won't be sent till the Send method is
// called.
func (c *Connection) Post() *Request {
	request := new(Request)
	request.transport = c
	request.method = http.MethodPost
	return request
}

// Patch creates an HTTP PATCH request. Note that this request won't be sent till the Send method is
// called.
func (c *Connection) Patch() *Request {
	request := new(Request)
	request.transport = c
	request.method = http.MethodPatch
	return request
}

// Put creates an HTTP PUT request. Note that this request won't be sent till the Send method is
// called.
func (c *Connection) Put() *Request {
	request := new(Request)
	request.transport = c
	request.method = http.MethodPut
	return request
}

// Delete creates an HTTP DELETE request. Note that this request won't be sent till the Send method
// is called.
func (c *Connection) Delete() *Request {
	request := new(Request)
	request.transport = c
	request.method = http.MethodDelete
	return request
}
