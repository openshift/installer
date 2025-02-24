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
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// KNativeOperation can be parsed from the returned Service.
type KNativeOperation struct {
	Status   KNativeOperationStatus   `json:"status"`
	Metadata KNativeOperationMetadata `json:"metadata"`
	// other irrelevant fields omitted

	config   *dcl.Config
	basePath string
	verb     string
	location string
}

// KNativeOperationMetadata contains the Labels block.
type KNativeOperationMetadata struct {
	SelfLink string            `json:"selfLink"`
	Labels   map[string]string `json:"labels"`
}

// KNativeOperationStatus contains the Conditions block.
type KNativeOperationStatus struct {
	Conditions []KNativeOperationCondition `json:"conditions"`
}

// KNativeOperationCondition contains the
type KNativeOperationCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

// Wait waits for an DNSOperation to complete by fetching the operation until it completes.
func (op *KNativeOperation) Wait(ctx context.Context, c *dcl.Config, basePath, verb string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c
	op.basePath = basePath
	op.verb = verb

	location, ok := op.Metadata.Labels["cloud.googleapis.com/location"]
	if !ok {
		return fmt.Errorf("no location found")
	}
	op.location = location

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *KNativeOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := fmt.Sprintf("https://%s-run.googleapis.com/%s", op.location, op.Metadata.SelfLink)
	resp, err := dcl.SendRequest(ctx, op.config, "GET", u, &bytes.Buffer{}, nil)
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, false, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}
	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}

	for _, condition := range op.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return resp, nil
		}
	}
	return nil, dcl.OperationNotDone{}
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *KNativeOperation) FirstResponse() (map[string]interface{}, bool) {
	return make(map[string]interface{}), false
}
