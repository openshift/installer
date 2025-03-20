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
package networkconnectivity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Hub struct {
	Name        *string           `json:"name"`
	CreateTime  *string           `json:"createTime"`
	UpdateTime  *string           `json:"updateTime"`
	Labels      map[string]string `json:"labels"`
	Description *string           `json:"description"`
	UniqueId    *string           `json:"uniqueId"`
	State       *HubStateEnum     `json:"state"`
	Project     *string           `json:"project"`
	RoutingVpcs []HubRoutingVpcs  `json:"routingVpcs"`
}

func (r *Hub) String() string {
	return dcl.SprintResource(r)
}

// The enum HubStateEnum.
type HubStateEnum string

// HubStateEnumRef returns a *HubStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func HubStateEnumRef(s string) *HubStateEnum {
	v := HubStateEnum(s)
	return &v
}

func (v HubStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "CREATING", "ACTIVE", "DELETING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "HubStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type HubRoutingVpcs struct {
	empty bool    `json:"-"`
	Uri   *string `json:"uri"`
}

type jsonHubRoutingVpcs HubRoutingVpcs

func (r *HubRoutingVpcs) UnmarshalJSON(data []byte) error {
	var res jsonHubRoutingVpcs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyHubRoutingVpcs
	} else {

		r.Uri = res.Uri

	}
	return nil
}

// This object is used to assert a desired state where this HubRoutingVpcs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyHubRoutingVpcs *HubRoutingVpcs = &HubRoutingVpcs{empty: true}

func (r *HubRoutingVpcs) Empty() bool {
	return r.empty
}

func (r *HubRoutingVpcs) String() string {
	return dcl.SprintResource(r)
}

func (r *HubRoutingVpcs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Hub) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "network_connectivity",
		Type:    "Hub",
		Version: "networkconnectivity",
	}
}

func (r *Hub) ID() (string, error) {
	if err := extractHubFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":         dcl.ValueOrEmptyString(nr.Name),
		"create_time":  dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":  dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":       dcl.ValueOrEmptyString(nr.Labels),
		"description":  dcl.ValueOrEmptyString(nr.Description),
		"unique_id":    dcl.ValueOrEmptyString(nr.UniqueId),
		"state":        dcl.ValueOrEmptyString(nr.State),
		"project":      dcl.ValueOrEmptyString(nr.Project),
		"routing_vpcs": dcl.ValueOrEmptyString(nr.RoutingVpcs),
	}
	return dcl.Nprintf("projects/{{project}}/locations/global/hubs/{{name}}", params), nil
}

const HubMaxPage = -1

type HubList struct {
	Items []*Hub

	nextToken string

	pageSize int32

	resource *Hub
}

func (l *HubList) HasNext() bool {
	return l.nextToken != ""
}

func (l *HubList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listHub(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListHub(ctx context.Context, project string) (*HubList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListHubWithMaxResults(ctx, project, HubMaxPage)

}

func (c *Client) ListHubWithMaxResults(ctx context.Context, project string, pageSize int32) (*HubList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Hub{
		Project: &project,
	}
	items, token, err := c.listHub(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &HubList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetHub(ctx context.Context, r *Hub) (*Hub, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractHubFields(r)

	b, err := c.getHubRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalHub(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeHubNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractHubFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteHub(ctx context.Context, r *Hub) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Hub resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Hub...")
	deleteOp := deleteHubOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllHub deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllHub(ctx context.Context, project string, filter func(*Hub) bool) error {
	listObj, err := c.ListHub(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllHub(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllHub(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyHub(ctx context.Context, rawDesired *Hub, opts ...dcl.ApplyOption) (*Hub, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Hub
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyHubHelper(c, ctx, rawDesired, opts...)
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

func applyHubHelper(c *Client, ctx context.Context, rawDesired *Hub, opts ...dcl.ApplyOption) (*Hub, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyHub...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractHubFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.hubDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToHubDiffs(c.Config, fieldDiffs, opts)
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
	var ops []hubApiOperation
	if create {
		ops = append(ops, &createHubOperation{})
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
	return applyHubDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyHubDiff(c *Client, ctx context.Context, desired *Hub, rawDesired *Hub, ops []hubApiOperation, opts ...dcl.ApplyOption) (*Hub, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetHub(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createHubOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapHub(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeHubNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeHubNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeHubDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractHubFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractHubFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffHub(c, newDesired, newState)
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
