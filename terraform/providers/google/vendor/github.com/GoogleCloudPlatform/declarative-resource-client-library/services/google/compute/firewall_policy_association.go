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
package compute

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type FirewallPolicyAssociation struct {
	Name             *string `json:"name"`
	AttachmentTarget *string `json:"attachmentTarget"`
	FirewallPolicy   *string `json:"firewallPolicy"`
	ShortName        *string `json:"shortName"`
}

func (r *FirewallPolicyAssociation) String() string {
	return dcl.SprintResource(r)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *FirewallPolicyAssociation) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "FirewallPolicyAssociation",
		Version: "compute",
	}
}

func (r *FirewallPolicyAssociation) ID() (string, error) {
	if err := extractFirewallPolicyAssociationFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":              dcl.ValueOrEmptyString(nr.Name),
		"attachment_target": dcl.ValueOrEmptyString(nr.AttachmentTarget),
		"firewall_policy":   dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"short_name":        dcl.ValueOrEmptyString(nr.ShortName),
	}
	return dcl.Nprintf("locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}", params), nil
}

const FirewallPolicyAssociationMaxPage = -1

type FirewallPolicyAssociationList struct {
	Items []*FirewallPolicyAssociation

	nextToken string

	pageSize int32

	resource *FirewallPolicyAssociation
}

func (l *FirewallPolicyAssociationList) HasNext() bool {
	return l.nextToken != ""
}

func (l *FirewallPolicyAssociationList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listFirewallPolicyAssociation(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListFirewallPolicyAssociation(ctx context.Context, firewallPolicy string) (*FirewallPolicyAssociationList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListFirewallPolicyAssociationWithMaxResults(ctx, firewallPolicy, FirewallPolicyAssociationMaxPage)

}

func (c *Client) ListFirewallPolicyAssociationWithMaxResults(ctx context.Context, firewallPolicy string, pageSize int32) (*FirewallPolicyAssociationList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &FirewallPolicyAssociation{
		FirewallPolicy: &firewallPolicy,
	}
	items, token, err := c.listFirewallPolicyAssociation(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &FirewallPolicyAssociationList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetFirewallPolicyAssociation(ctx context.Context, r *FirewallPolicyAssociation) (*FirewallPolicyAssociation, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractFirewallPolicyAssociationFields(r)

	b, err := c.getFirewallPolicyAssociationRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalFirewallPolicyAssociation(b, c, r)
	if err != nil {
		return nil, err
	}
	result.FirewallPolicy = r.FirewallPolicy
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeFirewallPolicyAssociationNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractFirewallPolicyAssociationFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteFirewallPolicyAssociation(ctx context.Context, r *FirewallPolicyAssociation) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("FirewallPolicyAssociation resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting FirewallPolicyAssociation...")
	deleteOp := deleteFirewallPolicyAssociationOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllFirewallPolicyAssociation deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllFirewallPolicyAssociation(ctx context.Context, firewallPolicy string, filter func(*FirewallPolicyAssociation) bool) error {
	listObj, err := c.ListFirewallPolicyAssociation(ctx, firewallPolicy)
	if err != nil {
		return err
	}

	err = c.deleteAllFirewallPolicyAssociation(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllFirewallPolicyAssociation(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyFirewallPolicyAssociation(ctx context.Context, rawDesired *FirewallPolicyAssociation, opts ...dcl.ApplyOption) (*FirewallPolicyAssociation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *FirewallPolicyAssociation
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyFirewallPolicyAssociationHelper(c, ctx, rawDesired, opts...)
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

func applyFirewallPolicyAssociationHelper(c *Client, ctx context.Context, rawDesired *FirewallPolicyAssociation, opts ...dcl.ApplyOption) (*FirewallPolicyAssociation, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyFirewallPolicyAssociation...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractFirewallPolicyAssociationFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.firewallPolicyAssociationDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToFirewallPolicyAssociationDiffs(c.Config, fieldDiffs, opts)
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
	var ops []firewallPolicyAssociationApiOperation
	if create {
		ops = append(ops, &createFirewallPolicyAssociationOperation{})
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
	return applyFirewallPolicyAssociationDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyFirewallPolicyAssociationDiff(c *Client, ctx context.Context, desired *FirewallPolicyAssociation, rawDesired *FirewallPolicyAssociation, ops []firewallPolicyAssociationApiOperation, opts ...dcl.ApplyOption) (*FirewallPolicyAssociation, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetFirewallPolicyAssociation(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createFirewallPolicyAssociationOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapFirewallPolicyAssociation(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeFirewallPolicyAssociationNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeFirewallPolicyAssociationNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeFirewallPolicyAssociationDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractFirewallPolicyAssociationFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractFirewallPolicyAssociationFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffFirewallPolicyAssociation(c, newDesired, newState)
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
