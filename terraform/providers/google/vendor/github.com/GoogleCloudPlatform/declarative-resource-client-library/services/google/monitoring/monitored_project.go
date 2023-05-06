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
package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type MonitoredProject struct {
	Name         *string `json:"name"`
	CreateTime   *string `json:"createTime"`
	MetricsScope *string `json:"metricsScope"`
}

func (r *MonitoredProject) String() string {
	return dcl.SprintResource(r)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *MonitoredProject) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "monitoring",
		Type:    "MonitoredProject",
		Version: "monitoring",
	}
}

func (r *MonitoredProject) ID() (string, error) {
	if err := extractMonitoredProjectFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":          dcl.ValueOrEmptyString(nr.Name),
		"create_time":   dcl.ValueOrEmptyString(nr.CreateTime),
		"metrics_scope": dcl.ValueOrEmptyString(nr.MetricsScope),
	}
	return dcl.Nprintf("locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}", params), nil
}

const MonitoredProjectMaxPage = -1

type MonitoredProjectList struct {
	Items []*MonitoredProject

	nextToken string

	pageSize int32

	resource *MonitoredProject
}

func (l *MonitoredProjectList) HasNext() bool {
	return l.nextToken != ""
}

func (l *MonitoredProjectList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listMonitoredProject(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListMonitoredProject(ctx context.Context, metricsScope string) (*MonitoredProjectList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "The caller does not have permission",
			Timeout:   120000000000,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListMonitoredProjectWithMaxResults(ctx, metricsScope, MonitoredProjectMaxPage)

}

func (c *Client) ListMonitoredProjectWithMaxResults(ctx context.Context, metricsScope string, pageSize int32) (*MonitoredProjectList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &MonitoredProject{
		MetricsScope: &metricsScope,
	}
	items, token, err := c.listMonitoredProject(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &MonitoredProjectList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) DeleteMonitoredProject(ctx context.Context, r *MonitoredProject) error {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "The caller does not have permission",
			Timeout:   120000000000,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("MonitoredProject resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting MonitoredProject...")
	deleteOp := deleteMonitoredProjectOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllMonitoredProject deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllMonitoredProject(ctx context.Context, metricsScope string, filter func(*MonitoredProject) bool) error {
	listObj, err := c.ListMonitoredProject(ctx, metricsScope)
	if err != nil {
		return err
	}

	err = c.deleteAllMonitoredProject(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllMonitoredProject(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyMonitoredProject(ctx context.Context, rawDesired *MonitoredProject, opts ...dcl.ApplyOption) (*MonitoredProject, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "The caller does not have permission",
			Timeout:   120000000000,
		},
	})))
	var resultNewState *MonitoredProject
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyMonitoredProjectHelper(c, ctx, rawDesired, opts...)
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

func applyMonitoredProjectHelper(c *Client, ctx context.Context, rawDesired *MonitoredProject, opts ...dcl.ApplyOption) (*MonitoredProject, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyMonitoredProject...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractMonitoredProjectFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.monitoredProjectDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToMonitoredProjectDiffs(c.Config, fieldDiffs, opts)
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
	var ops []monitoredProjectApiOperation
	if create {
		ops = append(ops, &createMonitoredProjectOperation{})
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
	return applyMonitoredProjectDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyMonitoredProjectDiff(c *Client, ctx context.Context, desired *MonitoredProject, rawDesired *MonitoredProject, ops []monitoredProjectApiOperation, opts ...dcl.ApplyOption) (*MonitoredProject, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetMonitoredProject(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createMonitoredProjectOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapMonitoredProject(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeMonitoredProjectNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeMonitoredProjectNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeMonitoredProjectDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractMonitoredProjectFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractMonitoredProjectFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffMonitoredProject(c, newDesired, newState)
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
