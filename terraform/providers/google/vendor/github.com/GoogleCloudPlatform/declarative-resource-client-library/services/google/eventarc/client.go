// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package eventarc defines operations in the declarative SDK.
package eventarc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// The Client is the base struct of all operations.  This will receive the
// Get, Delete, List, and Apply operations on all resources.
type Client struct {
	Config *dcl.Config
}

// NewClient creates a client that retries all operations a few times each.
func NewClient(c *dcl.Config) *Client {
	return &Client{
		Config: c,
	}
}
