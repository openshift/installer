// Copyright 2023 Google LLC. All Rights Reserved.
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
// Package osconfig defines types and functions for managing osconfig GCP resources.
package osconfig

import (
	"bytes"
	"context"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

// Returns true if m and n represent the same floating point value of seconds.
func canonicalizeOSPolicyAssignmentRolloutMinWaitDuration(m, n interface{}) bool {
	mStr := dcl.ValueOrEmptyString(m)
	nStr := dcl.ValueOrEmptyString(n)
	if mStr == "" && nStr == "" {
		return true
	}
	if mStr == "" || nStr == "" {
		return false
	}
	mDuration, err := time.ParseDuration(mStr)
	if err != nil {
		return false
	}
	nDuration, err := time.ParseDuration(nStr)
	if err != nil {
		return false
	}
	return mDuration == nDuration
}

// Waits for os policy assignment to be done reconciling before deletion.
func (r *OSPolicyAssignment) waitForNotReconciling(ctx context.Context, client *Client) error {
	return dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		nr, err := client.GetOSPolicyAssignment(ctx, r)
		if err != nil {
			return nil, err
		}
		if dcl.ValueOrEmptyBool(nr.Reconciling) {
			return &dcl.RetryDetails{}, dcl.OperationNotDone{}
		}
		return nil, nil
	}, client.Config.RetryProvider)
}

func (op *createOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(req), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	if !dcl.ValueOrEmptyBool(r.SkipAwaitRollout) {
		// wait for object to be created.
		var o operations.StandardGCPOperation
		if err := dcl.ParseResponse(resp.Response, &o); err != nil {
			return err
		}
		if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
			c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
			return err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
		op.response, _ = o.FirstResponse()
	}

	if _, err := c.GetOSPolicyAssignment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (op *updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
	_, err := c.GetOSPolicyAssignment(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateOSPolicyAssignment")
	if err != nil {
		return err
	}
	diffs := make([]*dcl.FieldDiff, 0)
	for _, d := range op.FieldDiffs {
		// skipAwaitUpdate is a custom field not available in the API and should not be included in an update mask
		if d.FieldName != "SkipAwaitRollout" {
			diffs = append(diffs, d)
		}
	}
	if len(diffs) == 0 {
		// Only diff was skipAwaitUpdate, return success
		return nil
	}
	mask := dcl.TopLevelUpdateMask(diffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	if !dcl.ValueOrEmptyBool(r.SkipAwaitRollout) {
		var o operations.StandardGCPOperation
		if err := dcl.ParseResponse(resp.Response, &o); err != nil {
			return err
		}
		err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

		if err != nil {
			return err
		}
	}

	return nil
}

func (op *deleteOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
	r, err := c.GetOSPolicyAssignment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "OSPolicyAssignment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetOSPolicyAssignment checking for existence. error: %v", err)
		return err
	}
	err = r.waitForNotReconciling(ctx, c)
	if err != nil {
		return err
	}
	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	if !dcl.ValueOrEmptyBool(r.SkipAwaitRollout) {
		// wait for object to be deleted.
		var o operations.StandardGCPOperation
		if err := dcl.ParseResponse(resp.Response, &o); err != nil {
			return err
		}
		if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
			return err
		}
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetOSPolicyAssignment(ctx, r)
		if dcl.IsNotFound(err) {
			return nil, nil
		}
		if retriesRemaining > 0 {
			retriesRemaining--
			return &dcl.RetryDetails{}, dcl.OperationNotDone{}
		}
		return nil, dcl.NotDeletedError{ExistingResource: r}
	}, c.Config.RetryProvider)
	return nil
}
