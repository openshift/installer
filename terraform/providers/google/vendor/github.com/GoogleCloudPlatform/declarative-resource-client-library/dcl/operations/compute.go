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
// Package operations contains all of the Operations used by the DCL.
package operations

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// ComputeOperation can be parsed from the returned API operation and waited on.
// Based on https://cloud.google.com/compute/docs/reference/rest/v1/regionOperations
type ComputeOperation struct {
	ID         string                 `json:"id"`
	Error      *ComputeOperationError `json:"error"`
	SelfLink   string                 `json:"selfLink"`
	Status     string                 `json:"status"`
	TargetLink string                 `json:"targetLink"`
	TargetID   string                 `json:"targetId"`
	// other irrelevant fields omitted

	config *dcl.Config
}

// ComputeOperationError is the GCE operation's Error body.
type ComputeOperationError struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Errors  []*ComputeOperationErrorError `json:"errors"`
}

// String formats the OperationError as an error string.
func (e *ComputeOperationError) String() string {
	if e == nil {
		return "nil"
	}
	var b strings.Builder
	for _, err := range e.Errors {
		fmt.Fprintf(&b, "error code %q, message: %s\n", err.Code, err.Message)
	}
	if e.Code != 0 || e.Message != "" {
		fmt.Fprintf(&b, "error code %d, message: %s\n", e.Code, e.Message)
	}

	return b.String()
}

// ComputeOperationErrorError is a singular error in a GCE operation.
type ComputeOperationErrorError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Wait waits for an ComputeOperation to complete by fetching the operation until it completes.
func (op *ComputeOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *ComputeOperation) handleResponse(resp *dcl.RetryDetails, err error) (*dcl.RetryDetails, error) {
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, false, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}

	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}

	if op.Status != "DONE" {
		return nil, dcl.OperationNotDone{}
	}

	if op.Error != nil {
		return nil, fmt.Errorf("operation received error: %v", op.Error)
	}

	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *ComputeOperation) FirstResponse() (map[string]interface{}, bool) {
	return make(map[string]interface{}), false
}

func (op *ComputeOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	return op.handleResponse(dcl.SendRequest(ctx, op.config, "GET", op.SelfLink, &bytes.Buffer{}, nil))
}

// ComputeGlobalOrganizationOperation can be parsed from the returned API operation and waited on.
// Based on https://cloud.google.com/compute/docs/reference/rest/v1/globalOrganizationOperations
type ComputeGlobalOrganizationOperation struct {
	BaseOperation ComputeOperation
	Parent        string
}

func (op *ComputeGlobalOrganizationOperation) Wait(ctx context.Context, c *dcl.Config, parent *string) error {
	c.Logger.Infof("Waiting on: %v", op)
	op.BaseOperation.config = c

	op.Parent = *parent

	return dcl.Do(ctx, op.operate, c.RetryProvider)
}

func (op *ComputeGlobalOrganizationOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	return op.BaseOperation.handleResponse(dcl.SendRequest(ctx, op.BaseOperation.config, "GET", op.BaseOperation.SelfLink+"?parentId="+op.Parent, &bytes.Buffer{}, nil))
}
