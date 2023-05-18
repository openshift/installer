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
package logging

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type LogBucket struct {
	Name           *string                      `json:"name"`
	Description    *string                      `json:"description"`
	CreateTime     *string                      `json:"createTime"`
	UpdateTime     *string                      `json:"updateTime"`
	RetentionDays  *int64                       `json:"retentionDays"`
	Locked         *bool                        `json:"locked"`
	LifecycleState *LogBucketLifecycleStateEnum `json:"lifecycleState"`
	Parent         *string                      `json:"parent"`
	Location       *string                      `json:"location"`
}

func (r *LogBucket) String() string {
	return dcl.SprintResource(r)
}

// The enum LogBucketLifecycleStateEnum.
type LogBucketLifecycleStateEnum string

// LogBucketLifecycleStateEnumRef returns a *LogBucketLifecycleStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func LogBucketLifecycleStateEnumRef(s string) *LogBucketLifecycleStateEnum {
	v := LogBucketLifecycleStateEnum(s)
	return &v
}

func (v LogBucketLifecycleStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LIFECYCLE_STATE_UNSPECIFIED", "ACTIVE", "DELETE_REQUESTED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LogBucketLifecycleStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *LogBucket) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "logging",
		Type:    "LogBucket",
		Version: "logging",
	}
}

func (r *LogBucket) ID() (string, error) {
	if err := extractLogBucketFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":            dcl.ValueOrEmptyString(nr.Name),
		"description":     dcl.ValueOrEmptyString(nr.Description),
		"create_time":     dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":     dcl.ValueOrEmptyString(nr.UpdateTime),
		"retention_days":  dcl.ValueOrEmptyString(nr.RetentionDays),
		"locked":          dcl.ValueOrEmptyString(nr.Locked),
		"lifecycle_state": dcl.ValueOrEmptyString(nr.LifecycleState),
		"parent":          dcl.ValueOrEmptyString(nr.Parent),
		"location":        dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("{{parent}}/locations/{{location}}/buckets/{{name}}", params), nil
}

const LogBucketMaxPage = -1

type LogBucketList struct {
	Items []*LogBucket

	nextToken string

	pageSize int32

	resource *LogBucket
}

func (l *LogBucketList) HasNext() bool {
	return l.nextToken != ""
}

func (l *LogBucketList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listLogBucket(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListLogBucket(ctx context.Context, location, parent string) (*LogBucketList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListLogBucketWithMaxResults(ctx, location, parent, LogBucketMaxPage)

}

func (c *Client) ListLogBucketWithMaxResults(ctx context.Context, location, parent string, pageSize int32) (*LogBucketList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &LogBucket{
		Location: &location,
		Parent:   &parent,
	}
	items, token, err := c.listLogBucket(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &LogBucketList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetLogBucket(ctx context.Context, r *LogBucket) (*LogBucket, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractLogBucketFields(r)

	b, err := c.getLogBucketRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalLogBucket(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Location = r.Location
	result.Parent = r.Parent
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeLogBucketNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractLogBucketFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteLogBucket(ctx context.Context, r *LogBucket) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("LogBucket resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting LogBucket...")
	deleteOp := deleteLogBucketOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllLogBucket deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllLogBucket(ctx context.Context, location, parent string, filter func(*LogBucket) bool) error {
	listObj, err := c.ListLogBucket(ctx, location, parent)
	if err != nil {
		return err
	}

	err = c.deleteAllLogBucket(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllLogBucket(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyLogBucket(ctx context.Context, rawDesired *LogBucket, opts ...dcl.ApplyOption) (*LogBucket, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *LogBucket
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyLogBucketHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyLogBucketHelper(c *Client, ctx context.Context, rawDesired *LogBucket, opts ...dcl.ApplyOption) (*LogBucket, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyLogBucket...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractLogBucketFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.logBucketDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToLogBucketDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []logBucketApiOperation
	if create {
		ops = append(ops, &createLogBucketOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyLogBucketDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyLogBucketDiff(c *Client, ctx context.Context, desired *LogBucket, rawDesired *LogBucket, ops []logBucketApiOperation, opts ...dcl.ApplyOption) (*LogBucket, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetLogBucket(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createLogBucketOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapLogBucket(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeLogBucketNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeLogBucketNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeLogBucketDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractLogBucketFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractLogBucketFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffLogBucket(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}
