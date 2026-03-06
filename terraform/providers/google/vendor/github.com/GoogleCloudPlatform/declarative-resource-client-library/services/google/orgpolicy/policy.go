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
package orgpolicy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Policy struct {
	Name       *string           `json:"name"`
	Spec       *PolicySpec       `json:"spec"`
	DryRunSpec *PolicyDryRunSpec `json:"dryRunSpec"`
	Etag       *string           `json:"etag"`
	Parent     *string           `json:"parent"`
}

func (r *Policy) String() string {
	return dcl.SprintResource(r)
}

type PolicySpec struct {
	empty             bool              `json:"-"`
	Etag              *string           `json:"etag"`
	UpdateTime        *string           `json:"updateTime"`
	Rules             []PolicySpecRules `json:"rules"`
	InheritFromParent *bool             `json:"inheritFromParent"`
	Reset             *bool             `json:"reset"`
}

type jsonPolicySpec PolicySpec

func (r *PolicySpec) UnmarshalJSON(data []byte) error {
	var res jsonPolicySpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicySpec
	} else {

		r.Etag = res.Etag

		r.UpdateTime = res.UpdateTime

		r.Rules = res.Rules

		r.InheritFromParent = res.InheritFromParent

		r.Reset = res.Reset

	}
	return nil
}

// This object is used to assert a desired state where this PolicySpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicySpec *PolicySpec = &PolicySpec{empty: true}

func (r *PolicySpec) Empty() bool {
	return r.empty
}

func (r *PolicySpec) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicySpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicySpecRules struct {
	empty     bool                      `json:"-"`
	Values    *PolicySpecRulesValues    `json:"values"`
	AllowAll  *bool                     `json:"allowAll"`
	DenyAll   *bool                     `json:"denyAll"`
	Enforce   *bool                     `json:"enforce"`
	Condition *PolicySpecRulesCondition `json:"condition"`
}

type jsonPolicySpecRules PolicySpecRules

func (r *PolicySpecRules) UnmarshalJSON(data []byte) error {
	var res jsonPolicySpecRules
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicySpecRules
	} else {

		r.Values = res.Values

		r.AllowAll = res.AllowAll

		r.DenyAll = res.DenyAll

		r.Enforce = res.Enforce

		r.Condition = res.Condition

	}
	return nil
}

// This object is used to assert a desired state where this PolicySpecRules is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicySpecRules *PolicySpecRules = &PolicySpecRules{empty: true}

func (r *PolicySpecRules) Empty() bool {
	return r.empty
}

func (r *PolicySpecRules) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicySpecRules) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicySpecRulesValues struct {
	empty         bool     `json:"-"`
	AllowedValues []string `json:"allowedValues"`
	DeniedValues  []string `json:"deniedValues"`
}

type jsonPolicySpecRulesValues PolicySpecRulesValues

func (r *PolicySpecRulesValues) UnmarshalJSON(data []byte) error {
	var res jsonPolicySpecRulesValues
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicySpecRulesValues
	} else {

		r.AllowedValues = res.AllowedValues

		r.DeniedValues = res.DeniedValues

	}
	return nil
}

// This object is used to assert a desired state where this PolicySpecRulesValues is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicySpecRulesValues *PolicySpecRulesValues = &PolicySpecRulesValues{empty: true}

func (r *PolicySpecRulesValues) Empty() bool {
	return r.empty
}

func (r *PolicySpecRulesValues) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicySpecRulesValues) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicySpecRulesCondition struct {
	empty       bool    `json:"-"`
	Expression  *string `json:"expression"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type jsonPolicySpecRulesCondition PolicySpecRulesCondition

func (r *PolicySpecRulesCondition) UnmarshalJSON(data []byte) error {
	var res jsonPolicySpecRulesCondition
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicySpecRulesCondition
	} else {

		r.Expression = res.Expression

		r.Title = res.Title

		r.Description = res.Description

		r.Location = res.Location

	}
	return nil
}

// This object is used to assert a desired state where this PolicySpecRulesCondition is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicySpecRulesCondition *PolicySpecRulesCondition = &PolicySpecRulesCondition{empty: true}

func (r *PolicySpecRulesCondition) Empty() bool {
	return r.empty
}

func (r *PolicySpecRulesCondition) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicySpecRulesCondition) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicyDryRunSpec struct {
	empty             bool                    `json:"-"`
	Etag              *string                 `json:"etag"`
	UpdateTime        *string                 `json:"updateTime"`
	Rules             []PolicyDryRunSpecRules `json:"rules"`
	InheritFromParent *bool                   `json:"inheritFromParent"`
	Reset             *bool                   `json:"reset"`
}

type jsonPolicyDryRunSpec PolicyDryRunSpec

func (r *PolicyDryRunSpec) UnmarshalJSON(data []byte) error {
	var res jsonPolicyDryRunSpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicyDryRunSpec
	} else {

		r.Etag = res.Etag

		r.UpdateTime = res.UpdateTime

		r.Rules = res.Rules

		r.InheritFromParent = res.InheritFromParent

		r.Reset = res.Reset

	}
	return nil
}

// This object is used to assert a desired state where this PolicyDryRunSpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicyDryRunSpec *PolicyDryRunSpec = &PolicyDryRunSpec{empty: true}

func (r *PolicyDryRunSpec) Empty() bool {
	return r.empty
}

func (r *PolicyDryRunSpec) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicyDryRunSpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicyDryRunSpecRules struct {
	empty     bool                            `json:"-"`
	Values    *PolicyDryRunSpecRulesValues    `json:"values"`
	AllowAll  *bool                           `json:"allowAll"`
	DenyAll   *bool                           `json:"denyAll"`
	Enforce   *bool                           `json:"enforce"`
	Condition *PolicyDryRunSpecRulesCondition `json:"condition"`
}

type jsonPolicyDryRunSpecRules PolicyDryRunSpecRules

func (r *PolicyDryRunSpecRules) UnmarshalJSON(data []byte) error {
	var res jsonPolicyDryRunSpecRules
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicyDryRunSpecRules
	} else {

		r.Values = res.Values

		r.AllowAll = res.AllowAll

		r.DenyAll = res.DenyAll

		r.Enforce = res.Enforce

		r.Condition = res.Condition

	}
	return nil
}

// This object is used to assert a desired state where this PolicyDryRunSpecRules is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicyDryRunSpecRules *PolicyDryRunSpecRules = &PolicyDryRunSpecRules{empty: true}

func (r *PolicyDryRunSpecRules) Empty() bool {
	return r.empty
}

func (r *PolicyDryRunSpecRules) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicyDryRunSpecRules) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicyDryRunSpecRulesValues struct {
	empty         bool     `json:"-"`
	AllowedValues []string `json:"allowedValues"`
	DeniedValues  []string `json:"deniedValues"`
}

type jsonPolicyDryRunSpecRulesValues PolicyDryRunSpecRulesValues

func (r *PolicyDryRunSpecRulesValues) UnmarshalJSON(data []byte) error {
	var res jsonPolicyDryRunSpecRulesValues
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicyDryRunSpecRulesValues
	} else {

		r.AllowedValues = res.AllowedValues

		r.DeniedValues = res.DeniedValues

	}
	return nil
}

// This object is used to assert a desired state where this PolicyDryRunSpecRulesValues is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicyDryRunSpecRulesValues *PolicyDryRunSpecRulesValues = &PolicyDryRunSpecRulesValues{empty: true}

func (r *PolicyDryRunSpecRulesValues) Empty() bool {
	return r.empty
}

func (r *PolicyDryRunSpecRulesValues) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicyDryRunSpecRulesValues) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PolicyDryRunSpecRulesCondition struct {
	empty       bool    `json:"-"`
	Expression  *string `json:"expression"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type jsonPolicyDryRunSpecRulesCondition PolicyDryRunSpecRulesCondition

func (r *PolicyDryRunSpecRulesCondition) UnmarshalJSON(data []byte) error {
	var res jsonPolicyDryRunSpecRulesCondition
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPolicyDryRunSpecRulesCondition
	} else {

		r.Expression = res.Expression

		r.Title = res.Title

		r.Description = res.Description

		r.Location = res.Location

	}
	return nil
}

// This object is used to assert a desired state where this PolicyDryRunSpecRulesCondition is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPolicyDryRunSpecRulesCondition *PolicyDryRunSpecRulesCondition = &PolicyDryRunSpecRulesCondition{empty: true}

func (r *PolicyDryRunSpecRulesCondition) Empty() bool {
	return r.empty
}

func (r *PolicyDryRunSpecRulesCondition) String() string {
	return dcl.SprintResource(r)
}

func (r *PolicyDryRunSpecRulesCondition) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Policy) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "org_policy",
		Type:    "Policy",
		Version: "orgpolicy",
	}
}

func (r *Policy) ID() (string, error) {
	if err := extractPolicyFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":         dcl.ValueOrEmptyString(nr.Name),
		"spec":         dcl.ValueOrEmptyString(nr.Spec),
		"dry_run_spec": dcl.ValueOrEmptyString(nr.DryRunSpec),
		"etag":         dcl.ValueOrEmptyString(nr.Etag),
		"parent":       dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.Nprintf("{{parent}}/policies/{{name}}", params), nil
}

const PolicyMaxPage = -1

type PolicyList struct {
	Items []*Policy

	nextToken string

	pageSize int32

	resource *Policy
}

func (l *PolicyList) HasNext() bool {
	return l.nextToken != ""
}

func (l *PolicyList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listPolicy(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListPolicy(ctx context.Context, parent string) (*PolicyList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "Permission 'orgpolicy\\.policy\\.[a-z]*' denied on resource '//orgpolicy\\.googleapis\\.com/(projects|folders)/[a-z0-9-]*/policies/[a-zA-Z.]*' \\(or it may not exist\\)\\.",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListPolicyWithMaxResults(ctx, parent, PolicyMaxPage)

}

func (c *Client) ListPolicyWithMaxResults(ctx context.Context, parent string, pageSize int32) (*PolicyList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Policy{
		Parent: &parent,
	}
	items, token, err := c.listPolicy(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &PolicyList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetPolicy(ctx context.Context, r *Policy) (*Policy, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "Permission 'orgpolicy\\.policy\\.[a-z]*' denied on resource '//orgpolicy\\.googleapis\\.com/(projects|folders)/[a-z0-9-]*/policies/[a-zA-Z.]*' \\(or it may not exist\\)\\.",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractPolicyFields(r)

	b, err := c.getPolicyRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalPolicy(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Parent = r.Parent
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizePolicyNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractPolicyFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeletePolicy(ctx context.Context, r *Policy) error {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "Permission 'orgpolicy\\.policy\\.[a-z]*' denied on resource '//orgpolicy\\.googleapis\\.com/(projects|folders)/[a-z0-9-]*/policies/[a-zA-Z.]*' \\(or it may not exist\\)\\.",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Policy resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Policy...")
	deleteOp := deletePolicyOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllPolicy deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllPolicy(ctx context.Context, parent string, filter func(*Policy) bool) error {
	listObj, err := c.ListPolicy(ctx, parent)
	if err != nil {
		return err
	}

	err = c.deleteAllPolicy(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllPolicy(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyPolicy(ctx context.Context, rawDesired *Policy, opts ...dcl.ApplyOption) (*Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "Permission 'orgpolicy\\.policy\\.[a-z]*' denied on resource '//orgpolicy\\.googleapis\\.com/(projects|folders)/[a-z0-9-]*/policies/[a-zA-Z.]*' \\(or it may not exist\\)\\.",
			Timeout:   0,
		},
	})))
	var resultNewState *Policy
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyPolicyHelper(c, ctx, rawDesired, opts...)
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

func applyPolicyHelper(c *Client, ctx context.Context, rawDesired *Policy, opts ...dcl.ApplyOption) (*Policy, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyPolicy...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractPolicyFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.policyDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToPolicyDiffs(c.Config, fieldDiffs, opts)
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
	var ops []policyApiOperation
	if create {
		ops = append(ops, &createPolicyOperation{})
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
	return applyPolicyDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyPolicyDiff(c *Client, ctx context.Context, desired *Policy, rawDesired *Policy, ops []policyApiOperation, opts ...dcl.ApplyOption) (*Policy, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetPolicy(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createPolicyOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapPolicy(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizePolicyNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizePolicyNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizePolicyDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractPolicyFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractPolicyFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffPolicy(c, newDesired, newState)
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
