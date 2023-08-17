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

type NetworkFirewallPolicyRule struct {
	Description           *string                                     `json:"description"`
	RuleName              *string                                     `json:"ruleName"`
	Priority              *int64                                      `json:"priority"`
	Location              *string                                     `json:"location"`
	Match                 *NetworkFirewallPolicyRuleMatch             `json:"match"`
	Action                *string                                     `json:"action"`
	Direction             *NetworkFirewallPolicyRuleDirectionEnum     `json:"direction"`
	EnableLogging         *bool                                       `json:"enableLogging"`
	RuleTupleCount        *int64                                      `json:"ruleTupleCount"`
	TargetServiceAccounts []string                                    `json:"targetServiceAccounts"`
	TargetSecureTags      []NetworkFirewallPolicyRuleTargetSecureTags `json:"targetSecureTags"`
	Disabled              *bool                                       `json:"disabled"`
	Kind                  *string                                     `json:"kind"`
	FirewallPolicy        *string                                     `json:"firewallPolicy"`
	Project               *string                                     `json:"project"`
}

func (r *NetworkFirewallPolicyRule) String() string {
	return dcl.SprintResource(r)
}

// The enum NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum.
type NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum string

// NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumRef returns a *NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumRef(s string) *NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum {
	v := NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum(s)
	return &v
}

func (v NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"EFFECTIVE", "INEFFECTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum NetworkFirewallPolicyRuleDirectionEnum.
type NetworkFirewallPolicyRuleDirectionEnum string

// NetworkFirewallPolicyRuleDirectionEnumRef returns a *NetworkFirewallPolicyRuleDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func NetworkFirewallPolicyRuleDirectionEnumRef(s string) *NetworkFirewallPolicyRuleDirectionEnum {
	v := NetworkFirewallPolicyRuleDirectionEnum(s)
	return &v
}

func (v NetworkFirewallPolicyRuleDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"INGRESS", "EGRESS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NetworkFirewallPolicyRuleDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum NetworkFirewallPolicyRuleTargetSecureTagsStateEnum.
type NetworkFirewallPolicyRuleTargetSecureTagsStateEnum string

// NetworkFirewallPolicyRuleTargetSecureTagsStateEnumRef returns a *NetworkFirewallPolicyRuleTargetSecureTagsStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func NetworkFirewallPolicyRuleTargetSecureTagsStateEnumRef(s string) *NetworkFirewallPolicyRuleTargetSecureTagsStateEnum {
	v := NetworkFirewallPolicyRuleTargetSecureTagsStateEnum(s)
	return &v
}

func (v NetworkFirewallPolicyRuleTargetSecureTagsStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"EFFECTIVE", "INEFFECTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "NetworkFirewallPolicyRuleTargetSecureTagsStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type NetworkFirewallPolicyRuleMatch struct {
	empty                   bool                                          `json:"-"`
	SrcIPRanges             []string                                      `json:"srcIPRanges"`
	DestIPRanges            []string                                      `json:"destIPRanges"`
	Layer4Configs           []NetworkFirewallPolicyRuleMatchLayer4Configs `json:"layer4Configs"`
	SrcSecureTags           []NetworkFirewallPolicyRuleMatchSrcSecureTags `json:"srcSecureTags"`
	SrcRegionCodes          []string                                      `json:"srcRegionCodes"`
	DestRegionCodes         []string                                      `json:"destRegionCodes"`
	SrcThreatIntelligences  []string                                      `json:"srcThreatIntelligences"`
	DestThreatIntelligences []string                                      `json:"destThreatIntelligences"`
	SrcFqdns                []string                                      `json:"srcFqdns"`
	DestFqdns               []string                                      `json:"destFqdns"`
	SrcAddressGroups        []string                                      `json:"srcAddressGroups"`
	DestAddressGroups       []string                                      `json:"destAddressGroups"`
}

type jsonNetworkFirewallPolicyRuleMatch NetworkFirewallPolicyRuleMatch

func (r *NetworkFirewallPolicyRuleMatch) UnmarshalJSON(data []byte) error {
	var res jsonNetworkFirewallPolicyRuleMatch
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNetworkFirewallPolicyRuleMatch
	} else {

		r.SrcIPRanges = res.SrcIPRanges

		r.DestIPRanges = res.DestIPRanges

		r.Layer4Configs = res.Layer4Configs

		r.SrcSecureTags = res.SrcSecureTags

		r.SrcRegionCodes = res.SrcRegionCodes

		r.DestRegionCodes = res.DestRegionCodes

		r.SrcThreatIntelligences = res.SrcThreatIntelligences

		r.DestThreatIntelligences = res.DestThreatIntelligences

		r.SrcFqdns = res.SrcFqdns

		r.DestFqdns = res.DestFqdns

		r.SrcAddressGroups = res.SrcAddressGroups

		r.DestAddressGroups = res.DestAddressGroups

	}
	return nil
}

// This object is used to assert a desired state where this NetworkFirewallPolicyRuleMatch is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNetworkFirewallPolicyRuleMatch *NetworkFirewallPolicyRuleMatch = &NetworkFirewallPolicyRuleMatch{empty: true}

func (r *NetworkFirewallPolicyRuleMatch) Empty() bool {
	return r.empty
}

func (r *NetworkFirewallPolicyRuleMatch) String() string {
	return dcl.SprintResource(r)
}

func (r *NetworkFirewallPolicyRuleMatch) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NetworkFirewallPolicyRuleMatchLayer4Configs struct {
	empty      bool     `json:"-"`
	IPProtocol *string  `json:"ipProtocol"`
	Ports      []string `json:"ports"`
}

type jsonNetworkFirewallPolicyRuleMatchLayer4Configs NetworkFirewallPolicyRuleMatchLayer4Configs

func (r *NetworkFirewallPolicyRuleMatchLayer4Configs) UnmarshalJSON(data []byte) error {
	var res jsonNetworkFirewallPolicyRuleMatchLayer4Configs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNetworkFirewallPolicyRuleMatchLayer4Configs
	} else {

		r.IPProtocol = res.IPProtocol

		r.Ports = res.Ports

	}
	return nil
}

// This object is used to assert a desired state where this NetworkFirewallPolicyRuleMatchLayer4Configs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNetworkFirewallPolicyRuleMatchLayer4Configs *NetworkFirewallPolicyRuleMatchLayer4Configs = &NetworkFirewallPolicyRuleMatchLayer4Configs{empty: true}

func (r *NetworkFirewallPolicyRuleMatchLayer4Configs) Empty() bool {
	return r.empty
}

func (r *NetworkFirewallPolicyRuleMatchLayer4Configs) String() string {
	return dcl.SprintResource(r)
}

func (r *NetworkFirewallPolicyRuleMatchLayer4Configs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NetworkFirewallPolicyRuleMatchSrcSecureTags struct {
	empty bool                                                  `json:"-"`
	Name  *string                                               `json:"name"`
	State *NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum `json:"state"`
}

type jsonNetworkFirewallPolicyRuleMatchSrcSecureTags NetworkFirewallPolicyRuleMatchSrcSecureTags

func (r *NetworkFirewallPolicyRuleMatchSrcSecureTags) UnmarshalJSON(data []byte) error {
	var res jsonNetworkFirewallPolicyRuleMatchSrcSecureTags
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNetworkFirewallPolicyRuleMatchSrcSecureTags
	} else {

		r.Name = res.Name

		r.State = res.State

	}
	return nil
}

// This object is used to assert a desired state where this NetworkFirewallPolicyRuleMatchSrcSecureTags is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNetworkFirewallPolicyRuleMatchSrcSecureTags *NetworkFirewallPolicyRuleMatchSrcSecureTags = &NetworkFirewallPolicyRuleMatchSrcSecureTags{empty: true}

func (r *NetworkFirewallPolicyRuleMatchSrcSecureTags) Empty() bool {
	return r.empty
}

func (r *NetworkFirewallPolicyRuleMatchSrcSecureTags) String() string {
	return dcl.SprintResource(r)
}

func (r *NetworkFirewallPolicyRuleMatchSrcSecureTags) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type NetworkFirewallPolicyRuleTargetSecureTags struct {
	empty bool                                                `json:"-"`
	Name  *string                                             `json:"name"`
	State *NetworkFirewallPolicyRuleTargetSecureTagsStateEnum `json:"state"`
}

type jsonNetworkFirewallPolicyRuleTargetSecureTags NetworkFirewallPolicyRuleTargetSecureTags

func (r *NetworkFirewallPolicyRuleTargetSecureTags) UnmarshalJSON(data []byte) error {
	var res jsonNetworkFirewallPolicyRuleTargetSecureTags
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyNetworkFirewallPolicyRuleTargetSecureTags
	} else {

		r.Name = res.Name

		r.State = res.State

	}
	return nil
}

// This object is used to assert a desired state where this NetworkFirewallPolicyRuleTargetSecureTags is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyNetworkFirewallPolicyRuleTargetSecureTags *NetworkFirewallPolicyRuleTargetSecureTags = &NetworkFirewallPolicyRuleTargetSecureTags{empty: true}

func (r *NetworkFirewallPolicyRuleTargetSecureTags) Empty() bool {
	return r.empty
}

func (r *NetworkFirewallPolicyRuleTargetSecureTags) String() string {
	return dcl.SprintResource(r)
}

func (r *NetworkFirewallPolicyRuleTargetSecureTags) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *NetworkFirewallPolicyRule) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "NetworkFirewallPolicyRule",
		Version: "compute",
	}
}

func (r *NetworkFirewallPolicyRule) ID() (string, error) {
	if err := extractNetworkFirewallPolicyRuleFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"description":             dcl.ValueOrEmptyString(nr.Description),
		"rule_name":               dcl.ValueOrEmptyString(nr.RuleName),
		"priority":                dcl.ValueOrEmptyString(nr.Priority),
		"location":                dcl.ValueOrEmptyString(nr.Location),
		"match":                   dcl.ValueOrEmptyString(nr.Match),
		"action":                  dcl.ValueOrEmptyString(nr.Action),
		"direction":               dcl.ValueOrEmptyString(nr.Direction),
		"enable_logging":          dcl.ValueOrEmptyString(nr.EnableLogging),
		"rule_tuple_count":        dcl.ValueOrEmptyString(nr.RuleTupleCount),
		"target_service_accounts": dcl.ValueOrEmptyString(nr.TargetServiceAccounts),
		"target_secure_tags":      dcl.ValueOrEmptyString(nr.TargetSecureTags),
		"disabled":                dcl.ValueOrEmptyString(nr.Disabled),
		"kind":                    dcl.ValueOrEmptyString(nr.Kind),
		"firewall_policy":         dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"project":                 dcl.ValueOrEmptyString(nr.Project),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.Nprintf("projects/{{project}}/regions/{{location}}}/firewallPolicies/{{firewall_policy}}/rules/{{priority}}", params), nil
	}

	return dcl.Nprintf("projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}", params), nil
}

const NetworkFirewallPolicyRuleMaxPage = -1

type NetworkFirewallPolicyRuleList struct {
	Items []*NetworkFirewallPolicyRule

	nextToken string

	pageSize int32

	resource *NetworkFirewallPolicyRule
}

func (l *NetworkFirewallPolicyRuleList) HasNext() bool {
	return l.nextToken != ""
}

func (l *NetworkFirewallPolicyRuleList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listNetworkFirewallPolicyRule(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListNetworkFirewallPolicyRule(ctx context.Context, project, location, firewallPolicy string) (*NetworkFirewallPolicyRuleList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListNetworkFirewallPolicyRuleWithMaxResults(ctx, project, location, firewallPolicy, NetworkFirewallPolicyRuleMaxPage)

}

func (c *Client) ListNetworkFirewallPolicyRuleWithMaxResults(ctx context.Context, project, location, firewallPolicy string, pageSize int32) (*NetworkFirewallPolicyRuleList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &NetworkFirewallPolicyRule{
		Project:        &project,
		Location:       &location,
		FirewallPolicy: &firewallPolicy,
	}
	items, token, err := c.listNetworkFirewallPolicyRule(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &NetworkFirewallPolicyRuleList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetNetworkFirewallPolicyRule(ctx context.Context, r *NetworkFirewallPolicyRule) (*NetworkFirewallPolicyRule, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractNetworkFirewallPolicyRuleFields(r)

	b, err := c.getNetworkFirewallPolicyRuleRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalNetworkFirewallPolicyRule(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.FirewallPolicy = r.FirewallPolicy
	result.Priority = r.Priority

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeNetworkFirewallPolicyRuleNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractNetworkFirewallPolicyRuleFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteNetworkFirewallPolicyRule(ctx context.Context, r *NetworkFirewallPolicyRule) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("NetworkFirewallPolicyRule resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting NetworkFirewallPolicyRule...")
	deleteOp := deleteNetworkFirewallPolicyRuleOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllNetworkFirewallPolicyRule deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllNetworkFirewallPolicyRule(ctx context.Context, project, location, firewallPolicy string, filter func(*NetworkFirewallPolicyRule) bool) error {
	listObj, err := c.ListNetworkFirewallPolicyRule(ctx, project, location, firewallPolicy)
	if err != nil {
		return err
	}

	err = c.deleteAllNetworkFirewallPolicyRule(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllNetworkFirewallPolicyRule(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyNetworkFirewallPolicyRule(ctx context.Context, rawDesired *NetworkFirewallPolicyRule, opts ...dcl.ApplyOption) (*NetworkFirewallPolicyRule, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *NetworkFirewallPolicyRule
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyNetworkFirewallPolicyRuleHelper(c, ctx, rawDesired, opts...)
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

func applyNetworkFirewallPolicyRuleHelper(c *Client, ctx context.Context, rawDesired *NetworkFirewallPolicyRule, opts ...dcl.ApplyOption) (*NetworkFirewallPolicyRule, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyNetworkFirewallPolicyRule...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractNetworkFirewallPolicyRuleFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.networkFirewallPolicyRuleDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToNetworkFirewallPolicyRuleDiffs(c.Config, fieldDiffs, opts)
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
	var ops []networkFirewallPolicyRuleApiOperation
	if create {
		ops = append(ops, &createNetworkFirewallPolicyRuleOperation{})
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
	return applyNetworkFirewallPolicyRuleDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyNetworkFirewallPolicyRuleDiff(c *Client, ctx context.Context, desired *NetworkFirewallPolicyRule, rawDesired *NetworkFirewallPolicyRule, ops []networkFirewallPolicyRuleApiOperation, opts ...dcl.ApplyOption) (*NetworkFirewallPolicyRule, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetNetworkFirewallPolicyRule(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createNetworkFirewallPolicyRuleOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapNetworkFirewallPolicyRule(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeNetworkFirewallPolicyRuleNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeNetworkFirewallPolicyRuleNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeNetworkFirewallPolicyRuleDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractNetworkFirewallPolicyRuleFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractNetworkFirewallPolicyRuleFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffNetworkFirewallPolicyRule(c, newDesired, newState)
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
