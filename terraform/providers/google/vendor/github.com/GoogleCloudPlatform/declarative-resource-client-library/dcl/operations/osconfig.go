// Copyright 2021 Google LLC. All Rights Reserved.
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

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// OSPolicyAssignmentDeleteOperation can be parsed from the returned API operation and waited on.
type OSPolicyAssignmentDeleteOperation struct {
	Name string `json:"name"`

	config *dcl.Config
}

// Wait waits for an OSPolicyAssignmentDeleteOperation to complete by waiting until the operation returns a 404.
func (op *OSPolicyAssignmentDeleteOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	c.Logger.Infof("Waiting on: %q", op.Name)
	op.config = c

	return dcl.Do(ctx, op.operate, c.RetryProvider)
}

func (op *OSPolicyAssignmentDeleteOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := dcl.URL(op.Name, "https://osconfig.googleapis.com/v1alpha", op.config.BasePath, nil)
	resp, err := dcl.SendRequest(ctx, op.config, "GET", u, &bytes.Buffer{}, nil)
	if dcl.IsNotFound(err) {
		return nil, nil
	}
	return resp, dcl.OperationNotDone{}
}
