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
package compute

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Route struct {
	Id               *int64         `json:"id"`
	Name             *string        `json:"name"`
	Description      *string        `json:"description"`
	Network          *string        `json:"network"`
	Tag              []string       `json:"tag"`
	DestRange        *string        `json:"destRange"`
	Priority         *int64         `json:"priority"`
	NextHopInstance  *string        `json:"nextHopInstance"`
	NextHopIP        *string        `json:"nextHopIP"`
	NextHopNetwork   *string        `json:"nextHopNetwork"`
	NextHopGateway   *string        `json:"nextHopGateway"`
	NextHopPeering   *string        `json:"nextHopPeering"`
	NextHopIlb       *string        `json:"nextHopIlb"`
	Warning          []RouteWarning `json:"warning"`
	NextHopVpnTunnel *string        `json:"nextHopVpnTunnel"`
	SelfLink         *string        `json:"selfLink"`
	Project          *string        `json:"project"`
}

func (r *Route) String() string {
	return dcl.SprintResource(r)
}

// The enum RouteWarningCodeEnum.
type RouteWarningCodeEnum string

// RouteWarningCodeEnumRef returns a *RouteWarningCodeEnum with the value of string s
// If the empty string is provided, nil is returned.
func RouteWarningCodeEnumRef(s string) *RouteWarningCodeEnum {
	v := RouteWarningCodeEnum(s)
	return &v
}

func (v RouteWarningCodeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"BAD_REQUEST", "FORBIDDEN", "NOT_FOUND", "CONFLICT", "GONE", "PRECONDITION_FAILED", "INTERNAL_ERROR", "SERVICE_UNAVAILABLE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "RouteWarningCodeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type RouteWarning struct {
	empty   bool                  `json:"-"`
	Code    *RouteWarningCodeEnum `json:"code"`
	Message *string               `json:"message"`
	Data    map[string]string     `json:"data"`
}

type jsonRouteWarning RouteWarning

func (r *RouteWarning) UnmarshalJSON(data []byte) error {
	var res jsonRouteWarning
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyRouteWarning
	} else {

		r.Code = res.Code

		r.Message = res.Message

		r.Data = res.Data

	}
	return nil
}

// This object is used to assert a desired state where this RouteWarning is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyRouteWarning *RouteWarning = &RouteWarning{empty: true}

func (r *RouteWarning) Empty() bool {
	return r.empty
}

func (r *RouteWarning) String() string {
	return dcl.SprintResource(r)
}

func (r *RouteWarning) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Route) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "Route",
		Version: "compute",
	}
}

func (r *Route) ID() (string, error) {
	if err := extractRouteFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"id":                  dcl.ValueOrEmptyString(nr.Id),
		"name":                dcl.ValueOrEmptyString(nr.Name),
		"description":         dcl.ValueOrEmptyString(nr.Description),
		"network":             dcl.ValueOrEmptyString(nr.Network),
		"tag":                 dcl.ValueOrEmptyString(nr.Tag),
		"dest_range":          dcl.ValueOrEmptyString(nr.DestRange),
		"priority":            dcl.ValueOrEmptyString(nr.Priority),
		"next_hop_instance":   dcl.ValueOrEmptyString(nr.NextHopInstance),
		"next_hop_ip":         dcl.ValueOrEmptyString(nr.NextHopIP),
		"next_hop_network":    dcl.ValueOrEmptyString(nr.NextHopNetwork),
		"next_hop_gateway":    dcl.ValueOrEmptyString(nr.NextHopGateway),
		"next_hop_peering":    dcl.ValueOrEmptyString(nr.NextHopPeering),
		"next_hop_ilb":        dcl.ValueOrEmptyString(nr.NextHopIlb),
		"warning":             dcl.ValueOrEmptyString(nr.Warning),
		"next_hop_vpn_tunnel": dcl.ValueOrEmptyString(nr.NextHopVpnTunnel),
		"self_link":           dcl.ValueOrEmptyString(nr.SelfLink),
		"project":             dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/global/routes/{{name}}", params), nil
}

const RouteMaxPage = -1

type RouteList struct {
	Items []*Route

	nextToken string

	pageSize int32

	resource *Route
}

func (l *RouteList) HasNext() bool {
	return l.nextToken != ""
}

func (l *RouteList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listRoute(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListRoute(ctx context.Context, project string) (*RouteList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListRouteWithMaxResults(ctx, project, RouteMaxPage)

}

func (c *Client) ListRouteWithMaxResults(ctx context.Context, project string, pageSize int32) (*RouteList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Route{
		Project: &project,
	}
	items, token, err := c.listRoute(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &RouteList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetRoute(ctx context.Context, r *Route) (*Route, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractRouteFields(r)

	b, err := c.getRouteRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalRoute(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name
	if dcl.IsZeroValue(result.Priority) {
		result.Priority = dcl.Int64(1000)
	}

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeRouteNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractRouteFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteRoute(ctx context.Context, r *Route) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Route resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Route...")
	deleteOp := deleteRouteOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllRoute deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllRoute(ctx context.Context, project string, filter func(*Route) bool) error {
	listObj, err := c.ListRoute(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllRoute(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllRoute(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyRoute(ctx context.Context, rawDesired *Route, opts ...dcl.ApplyOption) (*Route, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Route
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyRouteHelper(c, ctx, rawDesired, opts...)
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

func applyRouteHelper(c *Client, ctx context.Context, rawDesired *Route, opts ...dcl.ApplyOption) (*Route, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyRoute...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractRouteFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.routeDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToRouteDiffs(c.Config, fieldDiffs, opts)
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
	var ops []routeApiOperation
	if create {
		ops = append(ops, &createRouteOperation{})
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
	return applyRouteDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyRouteDiff(c *Client, ctx context.Context, desired *Route, rawDesired *Route, ops []routeApiOperation, opts ...dcl.ApplyOption) (*Route, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetRoute(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createRouteOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapRoute(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeRouteNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeRouteNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeRouteDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractRouteFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractRouteFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffRoute(c, newDesired, newState)
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
