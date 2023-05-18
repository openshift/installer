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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type MetricsScope struct {
	Name              *string                         `json:"name"`
	CreateTime        *string                         `json:"createTime"`
	UpdateTime        *string                         `json:"updateTime"`
	MonitoredProjects []MetricsScopeMonitoredProjects `json:"monitoredProjects"`
}

func (r *MetricsScope) String() string {
	return dcl.SprintResource(r)
}

type MetricsScopeMonitoredProjects struct {
	empty      bool    `json:"-"`
	Name       *string `json:"name"`
	CreateTime *string `json:"createTime"`
}

type jsonMetricsScopeMonitoredProjects MetricsScopeMonitoredProjects

func (r *MetricsScopeMonitoredProjects) UnmarshalJSON(data []byte) error {
	var res jsonMetricsScopeMonitoredProjects
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyMetricsScopeMonitoredProjects
	} else {

		r.Name = res.Name

		r.CreateTime = res.CreateTime

	}
	return nil
}

// This object is used to assert a desired state where this MetricsScopeMonitoredProjects is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyMetricsScopeMonitoredProjects *MetricsScopeMonitoredProjects = &MetricsScopeMonitoredProjects{empty: true}

func (r *MetricsScopeMonitoredProjects) Empty() bool {
	return r.empty
}

func (r *MetricsScopeMonitoredProjects) String() string {
	return dcl.SprintResource(r)
}

func (r *MetricsScopeMonitoredProjects) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *MetricsScope) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "monitoring",
		Type:    "MetricsScope",
		Version: "monitoring",
	}
}

func (r *MetricsScope) ID() (string, error) {
	if err := extractMetricsScopeFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":               dcl.ValueOrEmptyString(nr.Name),
		"create_time":        dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":        dcl.ValueOrEmptyString(nr.UpdateTime),
		"monitored_projects": dcl.ValueOrEmptyString(nr.MonitoredProjects),
	}
	return dcl.Nprintf("locations/global/metricsScopes/{{name}}", params), nil
}

const MetricsScopeMaxPage = -1

type MetricsScopeList struct {
	Items []*MetricsScope

	nextToken string

	resource *MetricsScope
}

func (c *Client) GetMetricsScope(ctx context.Context, r *MetricsScope) (*MetricsScope, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractMetricsScopeFields(r)

	b, err := c.getMetricsScopeRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalMetricsScope(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeMetricsScopeNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractMetricsScopeFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) ApplyMetricsScope(ctx context.Context, rawDesired *MetricsScope, opts ...dcl.ApplyOption) (*MetricsScope, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *MetricsScope
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyMetricsScopeHelper(c, ctx, rawDesired, opts...)
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

func applyMetricsScopeHelper(c *Client, ctx context.Context, rawDesired *MetricsScope, opts ...dcl.ApplyOption) (*MetricsScope, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyMetricsScope...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractMetricsScopeFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.metricsScopeDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToMetricsScopeDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		return nil, dcl.ApplyInfeasibleError{Message: "No initial state found for singleton resource."}
	} else {
		for _, d := range diffs {
			if d.UpdateOp == nil {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) no update method found for field", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}
	var ops []metricsScopeApiOperation
	for _, d := range diffs {
		ops = append(ops, d.UpdateOp)
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
	return applyMetricsScopeDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyMetricsScopeDiff(c *Client, ctx context.Context, desired *MetricsScope, rawDesired *MetricsScope, ops []metricsScopeApiOperation, opts ...dcl.ApplyOption) (*MetricsScope, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetMetricsScope(ctx, desired)
	if err != nil {
		return nil, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeMetricsScopeNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeMetricsScopeDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractMetricsScopeFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractMetricsScopeFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffMetricsScope(c, newDesired, newState)
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
