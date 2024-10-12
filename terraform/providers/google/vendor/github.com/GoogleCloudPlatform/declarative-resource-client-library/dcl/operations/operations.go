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
package operations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// StandardGCPOperation can be parsed from the returned API operation and waited on.
// This is the typical GCP operation.
type StandardGCPOperation struct {
	Name     string                     `json:"name"`
	Error    *StandardGCPOperationError `json:"error"`
	Done     bool                       `json:"done"`
	Response map[string]interface{}     `json:"response"`
	// other irrelevant fields omitted

	config   *dcl.Config
	basePath string
	verb     string

	response map[string]interface{}
}

// StandardGCPOperationError is the GCP operation's Error body.
type StandardGCPOperationError struct {
	Errors []*StandardGCPOperationErrorError `json:"errors"`

	StandardGCPOperationErrorError
}

// String formats the StandardGCPOperationError as an error string.
func (e *StandardGCPOperationError) String() string {
	if e == nil {
		return "nil"
	}
	var b strings.Builder
	for _, err := range e.Errors {
		fmt.Fprintf(&b, "error code %q, message: %s, details: %+v\n", err.Code, err.Message, err.Details)
	}

	if e.Code != "" {
		fmt.Fprintf(&b, "error code %q, message: %s, details: %+v\n", e.Code, e.Message, e.Details)
	}

	return b.String()
}

// StandardGCPOperationErrorError is a singular error in a GCP operation.
type StandardGCPOperationErrorError struct {
	Code    json.Number              `json:"code"`
	Message string                   `json:"message"`
	Details []map[string]interface{} `json:"details"`
}

// Wait waits for an StandardGCPOperation to complete by fetching the operation until it completes.
func (op *StandardGCPOperation) Wait(ctx context.Context, c *dcl.Config, basePath, verb string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c
	op.basePath = basePath
	op.verb = verb

	if len(op.Response) != 0 {
		op.response = op.Response
	}
	if op.Done {
		c.Logger.Infof("Completed operation: %v", op)
		return nil
	}

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *StandardGCPOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := dcl.URL(op.Name, op.basePath, op.config.BasePath, nil)
	resp, err := dcl.SendRequest(ctx, op.config, op.verb, u, &bytes.Buffer{}, nil)
	if err != nil {
		// Since we don't know when this operation started, we will assume the
		// context's timeout applies to all request errors.
		if dcl.IsRetryableRequestError(op.config, err, false, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}

	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}

	if !op.Done {
		return nil, dcl.OperationNotDone{}
	}

	if op.Error != nil {
		return nil, fmt.Errorf("operation received error: %+v details: %v", op.Error, op.Response)
	}

	if len(op.response) == 0 && len(op.Response) != 0 {
		op.response = op.Response
	}

	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *StandardGCPOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}
