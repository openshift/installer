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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *NetworkFirewallPolicyRule) validate() error {

	if err := dcl.Required(r, "priority"); err != nil {
		return err
	}
	if err := dcl.Required(r, "match"); err != nil {
		return err
	}
	if err := dcl.Required(r, "action"); err != nil {
		return err
	}
	if err := dcl.Required(r, "direction"); err != nil {
		return err
	}
	if err := dcl.Required(r, "firewallPolicy"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Match) {
		if err := r.Match.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *NetworkFirewallPolicyRuleMatch) validate() error {
	if err := dcl.Required(r, "layer4Configs"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicyRuleMatchLayer4Configs) validate() error {
	if err := dcl.Required(r, "ipProtocol"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicyRuleMatchSrcSecureTags) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicyRuleTargetSecureTags) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicyRule) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *NetworkFirewallPolicyRule) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"priority":       dcl.ValueOrEmptyString(nr.Priority),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/getRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/getRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
}

func (r *NetworkFirewallPolicyRule) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicyRule) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/addRule", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/addRule", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicyRule) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"priority":       dcl.ValueOrEmptyString(nr.Priority),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/removeRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/removeRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
}

// networkFirewallPolicyRuleApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type networkFirewallPolicyRuleApiOperation interface {
	do(context.Context, *NetworkFirewallPolicyRule, *Client) error
}

// newUpdateNetworkFirewallPolicyRulePatchRuleRequest creates a request for an
// NetworkFirewallPolicyRule resource's PatchRule update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateNetworkFirewallPolicyRulePatchRuleRequest(ctx context.Context, f *NetworkFirewallPolicyRule, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.RuleName; !dcl.IsEmptyValueIndirect(v) {
		req["ruleName"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleMatch(c, f.Match, res); err != nil {
		return nil, fmt.Errorf("error expanding Match into match: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["match"] = v
	}
	if v := f.Action; !dcl.IsEmptyValueIndirect(v) {
		req["action"] = v
	}
	if v := f.Direction; !dcl.IsEmptyValueIndirect(v) {
		req["direction"] = v
	}
	if v := f.EnableLogging; !dcl.IsEmptyValueIndirect(v) {
		req["enableLogging"] = v
	}
	if v := f.TargetServiceAccounts; v != nil {
		req["targetServiceAccounts"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleTargetSecureTagsSlice(c, f.TargetSecureTags, res); err != nil {
		return nil, fmt.Errorf("error expanding TargetSecureTags into targetSecureTags: %w", err)
	} else if v != nil {
		req["targetSecureTags"] = v
	}
	if v := f.Disabled; !dcl.IsEmptyValueIndirect(v) {
		req["disabled"] = v
	}
	return req, nil
}

// marshalUpdateNetworkFirewallPolicyRulePatchRuleRequest converts the update into
// the final JSON request body.
func marshalUpdateNetworkFirewallPolicyRulePatchRuleRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateNetworkFirewallPolicyRulePatchRuleOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateNetworkFirewallPolicyRulePatchRuleOperation) do(ctx context.Context, r *NetworkFirewallPolicyRule, c *Client) error {
	_, err := c.GetNetworkFirewallPolicyRule(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "PatchRule")
	if err != nil {
		return err
	}

	req, err := newUpdateNetworkFirewallPolicyRulePatchRuleRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateNetworkFirewallPolicyRulePatchRuleRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listNetworkFirewallPolicyRuleRaw(ctx context.Context, r *NetworkFirewallPolicyRule, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != NetworkFirewallPolicyRuleMaxPage {
		m["pageSize"] = fmt.Sprintf("%v", pageSize)
	}

	u, err = dcl.AddQueryParams(u, m)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	return ioutil.ReadAll(resp.Response.Body)
}

type listNetworkFirewallPolicyRuleOperation struct {
	Rules []map[string]interface{} `json:"rules"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listNetworkFirewallPolicyRule(ctx context.Context, r *NetworkFirewallPolicyRule, pageToken string, pageSize int32) ([]*NetworkFirewallPolicyRule, string, error) {
	b, err := c.listNetworkFirewallPolicyRuleRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listNetworkFirewallPolicyRuleOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*NetworkFirewallPolicyRule
	for _, v := range m.Rules {
		res, err := unmarshalMapNetworkFirewallPolicyRule(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.FirewallPolicy = r.FirewallPolicy
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllNetworkFirewallPolicyRule(ctx context.Context, f func(*NetworkFirewallPolicyRule) bool, resources []*NetworkFirewallPolicyRule) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteNetworkFirewallPolicyRule(ctx, res)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", strings.Join(errors, "\n"))
	} else {
		return nil
	}
}

type deleteNetworkFirewallPolicyRuleOperation struct{}

func (op *deleteNetworkFirewallPolicyRuleOperation) do(ctx context.Context, r *NetworkFirewallPolicyRule, c *Client) error {
	r, err := c.GetNetworkFirewallPolicyRule(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.InfoWithContextf(ctx, "NetworkFirewallPolicyRule not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetNetworkFirewallPolicyRule checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		return err
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetNetworkFirewallPolicyRule(ctx, r)
		if dcl.IsNotFoundOrCode(err, 400) {
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

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createNetworkFirewallPolicyRuleOperation struct {
	response map[string]interface{}
}

func (op *createNetworkFirewallPolicyRuleOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createNetworkFirewallPolicyRuleOperation) do(ctx context.Context, r *NetworkFirewallPolicyRule, c *Client) error {
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
	// wait for object to be created.
	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetNetworkFirewallPolicyRule(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getNetworkFirewallPolicyRuleRaw(ctx context.Context, r *NetworkFirewallPolicyRule) ([]byte, error) {

	u, err := r.getURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	b, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) networkFirewallPolicyRuleDiffsForRawDesired(ctx context.Context, rawDesired *NetworkFirewallPolicyRule, opts ...dcl.ApplyOption) (initial, desired *NetworkFirewallPolicyRule, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *NetworkFirewallPolicyRule
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*NetworkFirewallPolicyRule); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected NetworkFirewallPolicyRule, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetNetworkFirewallPolicyRule(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a NetworkFirewallPolicyRule resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve NetworkFirewallPolicyRule resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that NetworkFirewallPolicyRule resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeNetworkFirewallPolicyRuleDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for NetworkFirewallPolicyRule: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for NetworkFirewallPolicyRule: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractNetworkFirewallPolicyRuleFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeNetworkFirewallPolicyRuleInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for NetworkFirewallPolicyRule: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeNetworkFirewallPolicyRuleDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for NetworkFirewallPolicyRule: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffNetworkFirewallPolicyRule(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeNetworkFirewallPolicyRuleInitialState(rawInitial, rawDesired *NetworkFirewallPolicyRule) (*NetworkFirewallPolicyRule, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeNetworkFirewallPolicyRuleDesiredState(rawDesired, rawInitial *NetworkFirewallPolicyRule, opts ...dcl.ApplyOption) (*NetworkFirewallPolicyRule, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Match = canonicalizeNetworkFirewallPolicyRuleMatch(rawDesired.Match, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &NetworkFirewallPolicyRule{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.RuleName, rawInitial.RuleName) {
		canonicalDesired.RuleName = rawInitial.RuleName
	} else {
		canonicalDesired.RuleName = rawDesired.RuleName
	}
	if dcl.IsZeroValue(rawDesired.Priority) || (dcl.IsEmptyValueIndirect(rawDesired.Priority) && dcl.IsEmptyValueIndirect(rawInitial.Priority)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Priority = rawInitial.Priority
	} else {
		canonicalDesired.Priority = rawDesired.Priority
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	canonicalDesired.Match = canonicalizeNetworkFirewallPolicyRuleMatch(rawDesired.Match, rawInitial.Match, opts...)
	if dcl.StringCanonicalize(rawDesired.Action, rawInitial.Action) {
		canonicalDesired.Action = rawInitial.Action
	} else {
		canonicalDesired.Action = rawDesired.Action
	}
	if dcl.IsZeroValue(rawDesired.Direction) || (dcl.IsEmptyValueIndirect(rawDesired.Direction) && dcl.IsEmptyValueIndirect(rawInitial.Direction)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Direction = rawInitial.Direction
	} else {
		canonicalDesired.Direction = rawDesired.Direction
	}
	if dcl.BoolCanonicalize(rawDesired.EnableLogging, rawInitial.EnableLogging) {
		canonicalDesired.EnableLogging = rawInitial.EnableLogging
	} else {
		canonicalDesired.EnableLogging = rawDesired.EnableLogging
	}
	if dcl.StringArrayCanonicalize(rawDesired.TargetServiceAccounts, rawInitial.TargetServiceAccounts) {
		canonicalDesired.TargetServiceAccounts = rawInitial.TargetServiceAccounts
	} else {
		canonicalDesired.TargetServiceAccounts = rawDesired.TargetServiceAccounts
	}
	canonicalDesired.TargetSecureTags = canonicalizeNetworkFirewallPolicyRuleTargetSecureTagsSlice(rawDesired.TargetSecureTags, rawInitial.TargetSecureTags, opts...)
	if dcl.BoolCanonicalize(rawDesired.Disabled, rawInitial.Disabled) {
		canonicalDesired.Disabled = rawInitial.Disabled
	} else {
		canonicalDesired.Disabled = rawDesired.Disabled
	}
	if dcl.IsZeroValue(rawDesired.FirewallPolicy) || (dcl.IsEmptyValueIndirect(rawDesired.FirewallPolicy) && dcl.IsEmptyValueIndirect(rawInitial.FirewallPolicy)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.FirewallPolicy = rawInitial.FirewallPolicy
	} else {
		canonicalDesired.FirewallPolicy = rawDesired.FirewallPolicy
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeNetworkFirewallPolicyRuleNewState(c *Client, rawNew, rawDesired *NetworkFirewallPolicyRule) (*NetworkFirewallPolicyRule, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RuleName) && dcl.IsEmptyValueIndirect(rawDesired.RuleName) {
		rawNew.RuleName = rawDesired.RuleName
	} else {
		if dcl.StringCanonicalize(rawDesired.RuleName, rawNew.RuleName) {
			rawNew.RuleName = rawDesired.RuleName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Priority) && dcl.IsEmptyValueIndirect(rawDesired.Priority) {
		rawNew.Priority = rawDesired.Priority
	} else {
	}

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.Match) && dcl.IsEmptyValueIndirect(rawDesired.Match) {
		rawNew.Match = rawDesired.Match
	} else {
		rawNew.Match = canonicalizeNewNetworkFirewallPolicyRuleMatch(c, rawDesired.Match, rawNew.Match)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Action) && dcl.IsEmptyValueIndirect(rawDesired.Action) {
		rawNew.Action = rawDesired.Action
	} else {
		if dcl.StringCanonicalize(rawDesired.Action, rawNew.Action) {
			rawNew.Action = rawDesired.Action
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Direction) && dcl.IsEmptyValueIndirect(rawDesired.Direction) {
		rawNew.Direction = rawDesired.Direction
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.EnableLogging) && dcl.IsEmptyValueIndirect(rawDesired.EnableLogging) {
		rawNew.EnableLogging = rawDesired.EnableLogging
	} else {
		if dcl.BoolCanonicalize(rawDesired.EnableLogging, rawNew.EnableLogging) {
			rawNew.EnableLogging = rawDesired.EnableLogging
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RuleTupleCount) && dcl.IsEmptyValueIndirect(rawDesired.RuleTupleCount) {
		rawNew.RuleTupleCount = rawDesired.RuleTupleCount
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.TargetServiceAccounts) && dcl.IsEmptyValueIndirect(rawDesired.TargetServiceAccounts) {
		rawNew.TargetServiceAccounts = rawDesired.TargetServiceAccounts
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.TargetServiceAccounts, rawNew.TargetServiceAccounts) {
			rawNew.TargetServiceAccounts = rawDesired.TargetServiceAccounts
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.TargetSecureTags) && dcl.IsEmptyValueIndirect(rawDesired.TargetSecureTags) {
		rawNew.TargetSecureTags = rawDesired.TargetSecureTags
	} else {
		rawNew.TargetSecureTags = canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTagsSlice(c, rawDesired.TargetSecureTags, rawNew.TargetSecureTags)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Disabled) && dcl.IsEmptyValueIndirect(rawDesired.Disabled) {
		rawNew.Disabled = rawDesired.Disabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.Disabled, rawNew.Disabled) {
			rawNew.Disabled = rawDesired.Disabled
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Kind) && dcl.IsEmptyValueIndirect(rawDesired.Kind) {
		rawNew.Kind = rawDesired.Kind
	} else {
		if dcl.StringCanonicalize(rawDesired.Kind, rawNew.Kind) {
			rawNew.Kind = rawDesired.Kind
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.FirewallPolicy) && dcl.IsEmptyValueIndirect(rawDesired.FirewallPolicy) {
		rawNew.FirewallPolicy = rawDesired.FirewallPolicy
	} else {
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeNetworkFirewallPolicyRuleMatch(des, initial *NetworkFirewallPolicyRuleMatch, opts ...dcl.ApplyOption) *NetworkFirewallPolicyRuleMatch {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NetworkFirewallPolicyRuleMatch{}

	if dcl.StringArrayCanonicalize(des.SrcIPRanges, initial.SrcIPRanges) {
		cDes.SrcIPRanges = initial.SrcIPRanges
	} else {
		cDes.SrcIPRanges = des.SrcIPRanges
	}
	if dcl.StringArrayCanonicalize(des.DestIPRanges, initial.DestIPRanges) {
		cDes.DestIPRanges = initial.DestIPRanges
	} else {
		cDes.DestIPRanges = des.DestIPRanges
	}
	cDes.Layer4Configs = canonicalizeNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(des.Layer4Configs, initial.Layer4Configs, opts...)
	cDes.SrcSecureTags = canonicalizeNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(des.SrcSecureTags, initial.SrcSecureTags, opts...)
	if dcl.StringArrayCanonicalize(des.SrcRegionCodes, initial.SrcRegionCodes) {
		cDes.SrcRegionCodes = initial.SrcRegionCodes
	} else {
		cDes.SrcRegionCodes = des.SrcRegionCodes
	}
	if dcl.StringArrayCanonicalize(des.DestRegionCodes, initial.DestRegionCodes) {
		cDes.DestRegionCodes = initial.DestRegionCodes
	} else {
		cDes.DestRegionCodes = des.DestRegionCodes
	}
	if dcl.StringArrayCanonicalize(des.SrcThreatIntelligences, initial.SrcThreatIntelligences) {
		cDes.SrcThreatIntelligences = initial.SrcThreatIntelligences
	} else {
		cDes.SrcThreatIntelligences = des.SrcThreatIntelligences
	}
	if dcl.StringArrayCanonicalize(des.DestThreatIntelligences, initial.DestThreatIntelligences) {
		cDes.DestThreatIntelligences = initial.DestThreatIntelligences
	} else {
		cDes.DestThreatIntelligences = des.DestThreatIntelligences
	}
	if dcl.StringArrayCanonicalize(des.SrcFqdns, initial.SrcFqdns) {
		cDes.SrcFqdns = initial.SrcFqdns
	} else {
		cDes.SrcFqdns = des.SrcFqdns
	}
	if dcl.StringArrayCanonicalize(des.DestFqdns, initial.DestFqdns) {
		cDes.DestFqdns = initial.DestFqdns
	} else {
		cDes.DestFqdns = des.DestFqdns
	}
	if dcl.StringArrayCanonicalize(des.SrcAddressGroups, initial.SrcAddressGroups) {
		cDes.SrcAddressGroups = initial.SrcAddressGroups
	} else {
		cDes.SrcAddressGroups = des.SrcAddressGroups
	}
	if dcl.StringArrayCanonicalize(des.DestAddressGroups, initial.DestAddressGroups) {
		cDes.DestAddressGroups = initial.DestAddressGroups
	} else {
		cDes.DestAddressGroups = des.DestAddressGroups
	}

	return cDes
}

func canonicalizeNetworkFirewallPolicyRuleMatchSlice(des, initial []NetworkFirewallPolicyRuleMatch, opts ...dcl.ApplyOption) []NetworkFirewallPolicyRuleMatch {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NetworkFirewallPolicyRuleMatch, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNetworkFirewallPolicyRuleMatch(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NetworkFirewallPolicyRuleMatch, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNetworkFirewallPolicyRuleMatch(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNetworkFirewallPolicyRuleMatch(c *Client, des, nw *NetworkFirewallPolicyRuleMatch) *NetworkFirewallPolicyRuleMatch {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NetworkFirewallPolicyRuleMatch while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.SrcIPRanges, nw.SrcIPRanges) {
		nw.SrcIPRanges = des.SrcIPRanges
	}
	if dcl.StringArrayCanonicalize(des.DestIPRanges, nw.DestIPRanges) {
		nw.DestIPRanges = des.DestIPRanges
	}
	nw.Layer4Configs = canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c, des.Layer4Configs, nw.Layer4Configs)
	nw.SrcSecureTags = canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c, des.SrcSecureTags, nw.SrcSecureTags)
	if dcl.StringArrayCanonicalize(des.SrcRegionCodes, nw.SrcRegionCodes) {
		nw.SrcRegionCodes = des.SrcRegionCodes
	}
	if dcl.StringArrayCanonicalize(des.DestRegionCodes, nw.DestRegionCodes) {
		nw.DestRegionCodes = des.DestRegionCodes
	}
	if dcl.StringArrayCanonicalize(des.SrcThreatIntelligences, nw.SrcThreatIntelligences) {
		nw.SrcThreatIntelligences = des.SrcThreatIntelligences
	}
	if dcl.StringArrayCanonicalize(des.DestThreatIntelligences, nw.DestThreatIntelligences) {
		nw.DestThreatIntelligences = des.DestThreatIntelligences
	}
	if dcl.StringArrayCanonicalize(des.SrcFqdns, nw.SrcFqdns) {
		nw.SrcFqdns = des.SrcFqdns
	}
	if dcl.StringArrayCanonicalize(des.DestFqdns, nw.DestFqdns) {
		nw.DestFqdns = des.DestFqdns
	}
	if dcl.StringArrayCanonicalize(des.SrcAddressGroups, nw.SrcAddressGroups) {
		nw.SrcAddressGroups = des.SrcAddressGroups
	}
	if dcl.StringArrayCanonicalize(des.DestAddressGroups, nw.DestAddressGroups) {
		nw.DestAddressGroups = des.DestAddressGroups
	}

	return nw
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchSet(c *Client, des, nw []NetworkFirewallPolicyRuleMatch) []NetworkFirewallPolicyRuleMatch {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NetworkFirewallPolicyRuleMatch
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNetworkFirewallPolicyRuleMatchNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatch(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchSlice(c *Client, des, nw []NetworkFirewallPolicyRuleMatch) []NetworkFirewallPolicyRuleMatch {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NetworkFirewallPolicyRuleMatch
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatch(c, &d, &n))
	}

	return items
}

func canonicalizeNetworkFirewallPolicyRuleMatchLayer4Configs(des, initial *NetworkFirewallPolicyRuleMatchLayer4Configs, opts ...dcl.ApplyOption) *NetworkFirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NetworkFirewallPolicyRuleMatchLayer4Configs{}

	if dcl.StringCanonicalize(des.IPProtocol, initial.IPProtocol) || dcl.IsZeroValue(des.IPProtocol) {
		cDes.IPProtocol = initial.IPProtocol
	} else {
		cDes.IPProtocol = des.IPProtocol
	}
	if dcl.StringArrayCanonicalize(des.Ports, initial.Ports) {
		cDes.Ports = initial.Ports
	} else {
		cDes.Ports = des.Ports
	}

	return cDes
}

func canonicalizeNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(des, initial []NetworkFirewallPolicyRuleMatchLayer4Configs, opts ...dcl.ApplyOption) []NetworkFirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NetworkFirewallPolicyRuleMatchLayer4Configs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNetworkFirewallPolicyRuleMatchLayer4Configs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NetworkFirewallPolicyRuleMatchLayer4Configs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNetworkFirewallPolicyRuleMatchLayer4Configs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4Configs(c *Client, des, nw *NetworkFirewallPolicyRuleMatchLayer4Configs) *NetworkFirewallPolicyRuleMatchLayer4Configs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NetworkFirewallPolicyRuleMatchLayer4Configs while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IPProtocol, nw.IPProtocol) {
		nw.IPProtocol = des.IPProtocol
	}
	if dcl.StringArrayCanonicalize(des.Ports, nw.Ports) {
		nw.Ports = des.Ports
	}

	return nw
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4ConfigsSet(c *Client, des, nw []NetworkFirewallPolicyRuleMatchLayer4Configs) []NetworkFirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NetworkFirewallPolicyRuleMatchLayer4Configs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNetworkFirewallPolicyRuleMatchLayer4ConfigsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4Configs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, des, nw []NetworkFirewallPolicyRuleMatchLayer4Configs) []NetworkFirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NetworkFirewallPolicyRuleMatchLayer4Configs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatchLayer4Configs(c, &d, &n))
	}

	return items
}

func canonicalizeNetworkFirewallPolicyRuleMatchSrcSecureTags(des, initial *NetworkFirewallPolicyRuleMatchSrcSecureTags, opts ...dcl.ApplyOption) *NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NetworkFirewallPolicyRuleMatchSrcSecureTags{}

	if dcl.IsZeroValue(des.Name) || (dcl.IsEmptyValueIndirect(des.Name) && dcl.IsEmptyValueIndirect(initial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(des, initial []NetworkFirewallPolicyRuleMatchSrcSecureTags, opts ...dcl.ApplyOption) []NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NetworkFirewallPolicyRuleMatchSrcSecureTags, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNetworkFirewallPolicyRuleMatchSrcSecureTags(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NetworkFirewallPolicyRuleMatchSrcSecureTags, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNetworkFirewallPolicyRuleMatchSrcSecureTags(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTags(c *Client, des, nw *NetworkFirewallPolicyRuleMatchSrcSecureTags) *NetworkFirewallPolicyRuleMatchSrcSecureTags {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NetworkFirewallPolicyRuleMatchSrcSecureTags while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTagsSet(c *Client, des, nw []NetworkFirewallPolicyRuleMatchSrcSecureTags) []NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NetworkFirewallPolicyRuleMatchSrcSecureTags
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNetworkFirewallPolicyRuleMatchSrcSecureTagsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTags(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c *Client, des, nw []NetworkFirewallPolicyRuleMatchSrcSecureTags) []NetworkFirewallPolicyRuleMatchSrcSecureTags {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NetworkFirewallPolicyRuleMatchSrcSecureTags
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleMatchSrcSecureTags(c, &d, &n))
	}

	return items
}

func canonicalizeNetworkFirewallPolicyRuleTargetSecureTags(des, initial *NetworkFirewallPolicyRuleTargetSecureTags, opts ...dcl.ApplyOption) *NetworkFirewallPolicyRuleTargetSecureTags {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NetworkFirewallPolicyRuleTargetSecureTags{}

	if dcl.IsZeroValue(des.Name) || (dcl.IsEmptyValueIndirect(des.Name) && dcl.IsEmptyValueIndirect(initial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeNetworkFirewallPolicyRuleTargetSecureTagsSlice(des, initial []NetworkFirewallPolicyRuleTargetSecureTags, opts ...dcl.ApplyOption) []NetworkFirewallPolicyRuleTargetSecureTags {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NetworkFirewallPolicyRuleTargetSecureTags, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNetworkFirewallPolicyRuleTargetSecureTags(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NetworkFirewallPolicyRuleTargetSecureTags, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNetworkFirewallPolicyRuleTargetSecureTags(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTags(c *Client, des, nw *NetworkFirewallPolicyRuleTargetSecureTags) *NetworkFirewallPolicyRuleTargetSecureTags {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NetworkFirewallPolicyRuleTargetSecureTags while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTagsSet(c *Client, des, nw []NetworkFirewallPolicyRuleTargetSecureTags) []NetworkFirewallPolicyRuleTargetSecureTags {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NetworkFirewallPolicyRuleTargetSecureTags
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNetworkFirewallPolicyRuleTargetSecureTagsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTags(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTagsSlice(c *Client, des, nw []NetworkFirewallPolicyRuleTargetSecureTags) []NetworkFirewallPolicyRuleTargetSecureTags {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NetworkFirewallPolicyRuleTargetSecureTags
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNetworkFirewallPolicyRuleTargetSecureTags(c, &d, &n))
	}

	return items
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffNetworkFirewallPolicyRule(c *Client, desired, actual *NetworkFirewallPolicyRule, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RuleName, actual.RuleName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("RuleName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Priority, actual.Priority, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Priority")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Match, actual.Match, dcl.DiffInfo{ObjectFunction: compareNetworkFirewallPolicyRuleMatchNewStyle, EmptyObject: EmptyNetworkFirewallPolicyRuleMatch, OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Match")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Action, actual.Action, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Action")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Direction, actual.Direction, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Direction")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableLogging, actual.EnableLogging, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("EnableLogging")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RuleTupleCount, actual.RuleTupleCount, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RuleTupleCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetServiceAccounts, actual.TargetServiceAccounts, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("TargetServiceAccounts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetSecureTags, actual.TargetSecureTags, dcl.DiffInfo{ObjectFunction: compareNetworkFirewallPolicyRuleTargetSecureTagsNewStyle, EmptyObject: EmptyNetworkFirewallPolicyRuleTargetSecureTags, OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("TargetSecureTags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Disabled, actual.Disabled, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Disabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Kind, actual.Kind, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Kind")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FirewallPolicy, actual.FirewallPolicy, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FirewallPolicy")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if len(newDiffs) > 0 {
		c.Config.Logger.Infof("Diff function found diffs: %v", newDiffs)
	}
	return newDiffs, nil
}
func compareNetworkFirewallPolicyRuleMatchNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NetworkFirewallPolicyRuleMatch)
	if !ok {
		desiredNotPointer, ok := d.(NetworkFirewallPolicyRuleMatch)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatch or *NetworkFirewallPolicyRuleMatch", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NetworkFirewallPolicyRuleMatch)
	if !ok {
		actualNotPointer, ok := a.(NetworkFirewallPolicyRuleMatch)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatch", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SrcIPRanges, actual.SrcIPRanges, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestIPRanges, actual.DestIPRanges, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("DestIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Layer4Configs, actual.Layer4Configs, dcl.DiffInfo{ObjectFunction: compareNetworkFirewallPolicyRuleMatchLayer4ConfigsNewStyle, EmptyObject: EmptyNetworkFirewallPolicyRuleMatchLayer4Configs, OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Layer4Configs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SrcSecureTags, actual.SrcSecureTags, dcl.DiffInfo{ObjectFunction: compareNetworkFirewallPolicyRuleMatchSrcSecureTagsNewStyle, EmptyObject: EmptyNetworkFirewallPolicyRuleMatchSrcSecureTags, OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcSecureTags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SrcRegionCodes, actual.SrcRegionCodes, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcRegionCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestRegionCodes, actual.DestRegionCodes, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("DestRegionCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SrcThreatIntelligences, actual.SrcThreatIntelligences, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcThreatIntelligences")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestThreatIntelligences, actual.DestThreatIntelligences, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("DestThreatIntelligences")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SrcFqdns, actual.SrcFqdns, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcFqdns")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestFqdns, actual.DestFqdns, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("DestFqdns")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SrcAddressGroups, actual.SrcAddressGroups, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("SrcAddressGroups")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestAddressGroups, actual.DestAddressGroups, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("DestAddressGroups")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNetworkFirewallPolicyRuleMatchLayer4ConfigsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NetworkFirewallPolicyRuleMatchLayer4Configs)
	if !ok {
		desiredNotPointer, ok := d.(NetworkFirewallPolicyRuleMatchLayer4Configs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatchLayer4Configs or *NetworkFirewallPolicyRuleMatchLayer4Configs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NetworkFirewallPolicyRuleMatchLayer4Configs)
	if !ok {
		actualNotPointer, ok := a.(NetworkFirewallPolicyRuleMatchLayer4Configs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatchLayer4Configs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IPProtocol, actual.IPProtocol, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("IpProtocol")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Ports, actual.Ports, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Ports")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNetworkFirewallPolicyRuleMatchSrcSecureTagsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NetworkFirewallPolicyRuleMatchSrcSecureTags)
	if !ok {
		desiredNotPointer, ok := d.(NetworkFirewallPolicyRuleMatchSrcSecureTags)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatchSrcSecureTags or *NetworkFirewallPolicyRuleMatchSrcSecureTags", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NetworkFirewallPolicyRuleMatchSrcSecureTags)
	if !ok {
		actualNotPointer, ok := a.(NetworkFirewallPolicyRuleMatchSrcSecureTags)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleMatchSrcSecureTags", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNetworkFirewallPolicyRuleTargetSecureTagsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NetworkFirewallPolicyRuleTargetSecureTags)
	if !ok {
		desiredNotPointer, ok := d.(NetworkFirewallPolicyRuleTargetSecureTags)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleTargetSecureTags or *NetworkFirewallPolicyRuleTargetSecureTags", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NetworkFirewallPolicyRuleTargetSecureTags)
	if !ok {
		actualNotPointer, ok := a.(NetworkFirewallPolicyRuleTargetSecureTags)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkFirewallPolicyRuleTargetSecureTags", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *NetworkFirewallPolicyRule) urlNormalized() *NetworkFirewallPolicyRule {
	normalized := dcl.Copy(*r).(NetworkFirewallPolicyRule)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.RuleName = dcl.SelfLinkToName(r.RuleName)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Action = dcl.SelfLinkToName(r.Action)
	normalized.Kind = dcl.SelfLinkToName(r.Kind)
	normalized.FirewallPolicy = dcl.SelfLinkToName(r.FirewallPolicy)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *NetworkFirewallPolicyRule) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "PatchRule" {
		fields := map[string]interface{}{
			"project":        dcl.ValueOrEmptyString(nr.Project),
			"location":       dcl.ValueOrEmptyString(nr.Location),
			"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
			"priority":       dcl.ValueOrEmptyString(nr.Priority),
		}
		if dcl.IsRegion(nr.Location) {
			return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/patchRule?priority={{priority}}", nr.basePath(), userBasePath, fields), nil
		}

		return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/patchRule?priority={{priority}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the NetworkFirewallPolicyRule resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *NetworkFirewallPolicyRule) marshal(c *Client) ([]byte, error) {
	m, err := expandNetworkFirewallPolicyRule(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling NetworkFirewallPolicyRule: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalNetworkFirewallPolicyRule decodes JSON responses into the NetworkFirewallPolicyRule resource schema.
func unmarshalNetworkFirewallPolicyRule(b []byte, c *Client, res *NetworkFirewallPolicyRule) (*NetworkFirewallPolicyRule, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNetworkFirewallPolicyRule(m, c, res)
}

func unmarshalMapNetworkFirewallPolicyRule(m map[string]interface{}, c *Client, res *NetworkFirewallPolicyRule) (*NetworkFirewallPolicyRule, error) {

	flattened := flattenNetworkFirewallPolicyRule(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNetworkFirewallPolicyRule expands NetworkFirewallPolicyRule into a JSON request object.
func expandNetworkFirewallPolicyRule(c *Client, f *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.RuleName; dcl.ValueShouldBeSent(v) {
		m["ruleName"] = v
	}
	if v := f.Priority; dcl.ValueShouldBeSent(v) {
		m["priority"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleMatch(c, f.Match, res); err != nil {
		return nil, fmt.Errorf("error expanding Match into match: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["match"] = v
	}
	if v := f.Action; dcl.ValueShouldBeSent(v) {
		m["action"] = v
	}
	if v := f.Direction; dcl.ValueShouldBeSent(v) {
		m["direction"] = v
	}
	if v := f.EnableLogging; dcl.ValueShouldBeSent(v) {
		m["enableLogging"] = v
	}
	if v := f.TargetServiceAccounts; v != nil {
		m["targetServiceAccounts"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleTargetSecureTagsSlice(c, f.TargetSecureTags, res); err != nil {
		return nil, fmt.Errorf("error expanding TargetSecureTags into targetSecureTags: %w", err)
	} else if v != nil {
		m["targetSecureTags"] = v
	}
	if v := f.Disabled; dcl.ValueShouldBeSent(v) {
		m["disabled"] = v
	}
	if v := f.FirewallPolicy; dcl.ValueShouldBeSent(v) {
		m["firewallPolicy"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyRule flattens NetworkFirewallPolicyRule from a JSON request object into the
// NetworkFirewallPolicyRule type.
func flattenNetworkFirewallPolicyRule(c *Client, i interface{}, res *NetworkFirewallPolicyRule) *NetworkFirewallPolicyRule {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &NetworkFirewallPolicyRule{}
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.RuleName = dcl.FlattenString(m["ruleName"])
	resultRes.Priority = dcl.FlattenInteger(m["priority"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Match = flattenNetworkFirewallPolicyRuleMatch(c, m["match"], res)
	resultRes.Action = dcl.FlattenString(m["action"])
	resultRes.Direction = flattenNetworkFirewallPolicyRuleDirectionEnum(m["direction"])
	resultRes.EnableLogging = dcl.FlattenBool(m["enableLogging"])
	resultRes.RuleTupleCount = dcl.FlattenInteger(m["ruleTupleCount"])
	resultRes.TargetServiceAccounts = dcl.FlattenStringSlice(m["targetServiceAccounts"])
	resultRes.TargetSecureTags = flattenNetworkFirewallPolicyRuleTargetSecureTagsSlice(c, m["targetSecureTags"], res)
	resultRes.Disabled = dcl.FlattenBool(m["disabled"])
	resultRes.Kind = dcl.FlattenString(m["kind"])
	resultRes.FirewallPolicy = dcl.FlattenString(m["firewallPolicy"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandNetworkFirewallPolicyRuleMatchMap expands the contents of NetworkFirewallPolicyRuleMatch into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchMap(c *Client, f map[string]NetworkFirewallPolicyRuleMatch, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatch(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNetworkFirewallPolicyRuleMatchSlice expands the contents of NetworkFirewallPolicyRuleMatch into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchSlice(c *Client, f []NetworkFirewallPolicyRuleMatch, res *NetworkFirewallPolicyRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatch(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNetworkFirewallPolicyRuleMatchMap flattens the contents of NetworkFirewallPolicyRuleMatch from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleMatch {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleMatch{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleMatch{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleMatch)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleMatch(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNetworkFirewallPolicyRuleMatchSlice flattens the contents of NetworkFirewallPolicyRuleMatch from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleMatch {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleMatch{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleMatch{}
	}

	items := make([]NetworkFirewallPolicyRuleMatch, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleMatch(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNetworkFirewallPolicyRuleMatch expands an instance of NetworkFirewallPolicyRuleMatch into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatch(c *Client, f *NetworkFirewallPolicyRuleMatch, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SrcIPRanges; v != nil {
		m["srcIpRanges"] = v
	}
	if v := f.DestIPRanges; v != nil {
		m["destIpRanges"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c, f.Layer4Configs, res); err != nil {
		return nil, fmt.Errorf("error expanding Layer4Configs into layer4Configs: %w", err)
	} else if v != nil {
		m["layer4Configs"] = v
	}
	if v, err := expandNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c, f.SrcSecureTags, res); err != nil {
		return nil, fmt.Errorf("error expanding SrcSecureTags into srcSecureTags: %w", err)
	} else if v != nil {
		m["srcSecureTags"] = v
	}
	if v := f.SrcRegionCodes; v != nil {
		m["srcRegionCodes"] = v
	}
	if v := f.DestRegionCodes; v != nil {
		m["destRegionCodes"] = v
	}
	if v := f.SrcThreatIntelligences; v != nil {
		m["srcThreatIntelligences"] = v
	}
	if v := f.DestThreatIntelligences; v != nil {
		m["destThreatIntelligences"] = v
	}
	if v := f.SrcFqdns; v != nil {
		m["srcFqdns"] = v
	}
	if v := f.DestFqdns; v != nil {
		m["destFqdns"] = v
	}
	if v := f.SrcAddressGroups; v != nil {
		m["srcAddressGroups"] = v
	}
	if v := f.DestAddressGroups; v != nil {
		m["destAddressGroups"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyRuleMatch flattens an instance of NetworkFirewallPolicyRuleMatch from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatch(c *Client, i interface{}, res *NetworkFirewallPolicyRule) *NetworkFirewallPolicyRuleMatch {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NetworkFirewallPolicyRuleMatch{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNetworkFirewallPolicyRuleMatch
	}
	r.SrcIPRanges = dcl.FlattenStringSlice(m["srcIpRanges"])
	r.DestIPRanges = dcl.FlattenStringSlice(m["destIpRanges"])
	r.Layer4Configs = flattenNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c, m["layer4Configs"], res)
	r.SrcSecureTags = flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c, m["srcSecureTags"], res)
	r.SrcRegionCodes = dcl.FlattenStringSlice(m["srcRegionCodes"])
	r.DestRegionCodes = dcl.FlattenStringSlice(m["destRegionCodes"])
	r.SrcThreatIntelligences = dcl.FlattenStringSlice(m["srcThreatIntelligences"])
	r.DestThreatIntelligences = dcl.FlattenStringSlice(m["destThreatIntelligences"])
	r.SrcFqdns = dcl.FlattenStringSlice(m["srcFqdns"])
	r.DestFqdns = dcl.FlattenStringSlice(m["destFqdns"])
	r.SrcAddressGroups = dcl.FlattenStringSlice(m["srcAddressGroups"])
	r.DestAddressGroups = dcl.FlattenStringSlice(m["destAddressGroups"])

	return r
}

// expandNetworkFirewallPolicyRuleMatchLayer4ConfigsMap expands the contents of NetworkFirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchLayer4ConfigsMap(c *Client, f map[string]NetworkFirewallPolicyRuleMatchLayer4Configs, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatchLayer4Configs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice expands the contents of NetworkFirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, f []NetworkFirewallPolicyRuleMatchLayer4Configs, res *NetworkFirewallPolicyRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatchLayer4Configs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNetworkFirewallPolicyRuleMatchLayer4ConfigsMap flattens the contents of NetworkFirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchLayer4ConfigsMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleMatchLayer4Configs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleMatchLayer4Configs{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleMatchLayer4Configs{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleMatchLayer4Configs)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleMatchLayer4Configs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice flattens the contents of NetworkFirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleMatchLayer4Configs {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleMatchLayer4Configs{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleMatchLayer4Configs{}
	}

	items := make([]NetworkFirewallPolicyRuleMatchLayer4Configs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleMatchLayer4Configs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNetworkFirewallPolicyRuleMatchLayer4Configs expands an instance of NetworkFirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchLayer4Configs(c *Client, f *NetworkFirewallPolicyRuleMatchLayer4Configs, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IPProtocol; !dcl.IsEmptyValueIndirect(v) {
		m["ipProtocol"] = v
	}
	if v := f.Ports; v != nil {
		m["ports"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyRuleMatchLayer4Configs flattens an instance of NetworkFirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchLayer4Configs(c *Client, i interface{}, res *NetworkFirewallPolicyRule) *NetworkFirewallPolicyRuleMatchLayer4Configs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NetworkFirewallPolicyRuleMatchLayer4Configs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNetworkFirewallPolicyRuleMatchLayer4Configs
	}
	r.IPProtocol = dcl.FlattenString(m["ipProtocol"])
	r.Ports = dcl.FlattenStringSlice(m["ports"])

	return r
}

// expandNetworkFirewallPolicyRuleMatchSrcSecureTagsMap expands the contents of NetworkFirewallPolicyRuleMatchSrcSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchSrcSecureTagsMap(c *Client, f map[string]NetworkFirewallPolicyRuleMatchSrcSecureTags, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatchSrcSecureTags(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice expands the contents of NetworkFirewallPolicyRuleMatchSrcSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c *Client, f []NetworkFirewallPolicyRuleMatchSrcSecureTags, res *NetworkFirewallPolicyRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNetworkFirewallPolicyRuleMatchSrcSecureTags(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsMap flattens the contents of NetworkFirewallPolicyRuleMatchSrcSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleMatchSrcSecureTags {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleMatchSrcSecureTags{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleMatchSrcSecureTags{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleMatchSrcSecureTags)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleMatchSrcSecureTags(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice flattens the contents of NetworkFirewallPolicyRuleMatchSrcSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleMatchSrcSecureTags {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleMatchSrcSecureTags{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleMatchSrcSecureTags{}
	}

	items := make([]NetworkFirewallPolicyRuleMatchSrcSecureTags, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleMatchSrcSecureTags(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNetworkFirewallPolicyRuleMatchSrcSecureTags expands an instance of NetworkFirewallPolicyRuleMatchSrcSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleMatchSrcSecureTags(c *Client, f *NetworkFirewallPolicyRuleMatchSrcSecureTags, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTags flattens an instance of NetworkFirewallPolicyRuleMatchSrcSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTags(c *Client, i interface{}, res *NetworkFirewallPolicyRule) *NetworkFirewallPolicyRuleMatchSrcSecureTags {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NetworkFirewallPolicyRuleMatchSrcSecureTags{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNetworkFirewallPolicyRuleMatchSrcSecureTags
	}
	r.Name = dcl.FlattenString(m["name"])
	r.State = flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum(m["state"])

	return r
}

// expandNetworkFirewallPolicyRuleTargetSecureTagsMap expands the contents of NetworkFirewallPolicyRuleTargetSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleTargetSecureTagsMap(c *Client, f map[string]NetworkFirewallPolicyRuleTargetSecureTags, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNetworkFirewallPolicyRuleTargetSecureTags(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNetworkFirewallPolicyRuleTargetSecureTagsSlice expands the contents of NetworkFirewallPolicyRuleTargetSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleTargetSecureTagsSlice(c *Client, f []NetworkFirewallPolicyRuleTargetSecureTags, res *NetworkFirewallPolicyRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNetworkFirewallPolicyRuleTargetSecureTags(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNetworkFirewallPolicyRuleTargetSecureTagsMap flattens the contents of NetworkFirewallPolicyRuleTargetSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleTargetSecureTagsMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleTargetSecureTags {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleTargetSecureTags{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleTargetSecureTags{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleTargetSecureTags)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleTargetSecureTags(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNetworkFirewallPolicyRuleTargetSecureTagsSlice flattens the contents of NetworkFirewallPolicyRuleTargetSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleTargetSecureTagsSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleTargetSecureTags {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleTargetSecureTags{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleTargetSecureTags{}
	}

	items := make([]NetworkFirewallPolicyRuleTargetSecureTags, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleTargetSecureTags(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNetworkFirewallPolicyRuleTargetSecureTags expands an instance of NetworkFirewallPolicyRuleTargetSecureTags into a JSON
// request object.
func expandNetworkFirewallPolicyRuleTargetSecureTags(c *Client, f *NetworkFirewallPolicyRuleTargetSecureTags, res *NetworkFirewallPolicyRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyRuleTargetSecureTags flattens an instance of NetworkFirewallPolicyRuleTargetSecureTags from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleTargetSecureTags(c *Client, i interface{}, res *NetworkFirewallPolicyRule) *NetworkFirewallPolicyRuleTargetSecureTags {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NetworkFirewallPolicyRuleTargetSecureTags{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNetworkFirewallPolicyRuleTargetSecureTags
	}
	r.Name = dcl.FlattenString(m["name"])
	r.State = flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnum(m["state"])

	return r
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumMap flattens the contents of NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum(item.(interface{}))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumSlice flattens the contents of NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum{}
	}

	items := make([]NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum(item.(interface{})))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum asserts that an interface is a string, and returns a
// pointer to a *NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum with the same value as that string.
func flattenNetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum(i interface{}) *NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnumRef(s)
}

// flattenNetworkFirewallPolicyRuleDirectionEnumMap flattens the contents of NetworkFirewallPolicyRuleDirectionEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleDirectionEnumMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleDirectionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleDirectionEnum{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleDirectionEnum{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleDirectionEnum)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleDirectionEnum(item.(interface{}))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleDirectionEnumSlice flattens the contents of NetworkFirewallPolicyRuleDirectionEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleDirectionEnumSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleDirectionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleDirectionEnum{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleDirectionEnum{}
	}

	items := make([]NetworkFirewallPolicyRuleDirectionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleDirectionEnum(item.(interface{})))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleDirectionEnum asserts that an interface is a string, and returns a
// pointer to a *NetworkFirewallPolicyRuleDirectionEnum with the same value as that string.
func flattenNetworkFirewallPolicyRuleDirectionEnum(i interface{}) *NetworkFirewallPolicyRuleDirectionEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return NetworkFirewallPolicyRuleDirectionEnumRef(s)
}

// flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnumMap flattens the contents of NetworkFirewallPolicyRuleTargetSecureTagsStateEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnumMap(c *Client, i interface{}, res *NetworkFirewallPolicyRule) map[string]NetworkFirewallPolicyRuleTargetSecureTagsStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkFirewallPolicyRuleTargetSecureTagsStateEnum{}
	}

	if len(a) == 0 {
		return map[string]NetworkFirewallPolicyRuleTargetSecureTagsStateEnum{}
	}

	items := make(map[string]NetworkFirewallPolicyRuleTargetSecureTagsStateEnum)
	for k, item := range a {
		items[k] = *flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnum(item.(interface{}))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnumSlice flattens the contents of NetworkFirewallPolicyRuleTargetSecureTagsStateEnum from a JSON
// response object.
func flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnumSlice(c *Client, i interface{}, res *NetworkFirewallPolicyRule) []NetworkFirewallPolicyRuleTargetSecureTagsStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkFirewallPolicyRuleTargetSecureTagsStateEnum{}
	}

	if len(a) == 0 {
		return []NetworkFirewallPolicyRuleTargetSecureTagsStateEnum{}
	}

	items := make([]NetworkFirewallPolicyRuleTargetSecureTagsStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnum(item.(interface{})))
	}

	return items
}

// flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnum asserts that an interface is a string, and returns a
// pointer to a *NetworkFirewallPolicyRuleTargetSecureTagsStateEnum with the same value as that string.
func flattenNetworkFirewallPolicyRuleTargetSecureTagsStateEnum(i interface{}) *NetworkFirewallPolicyRuleTargetSecureTagsStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return NetworkFirewallPolicyRuleTargetSecureTagsStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *NetworkFirewallPolicyRule) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNetworkFirewallPolicyRule(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Project == nil && ncr.Project == nil {
			c.Config.Logger.Info("Both Project fields null - considering equal.")
		} else if nr.Project == nil || ncr.Project == nil {
			c.Config.Logger.Info("Only one Project field is null - considering unequal.")
			return false
		} else if *nr.Project != *ncr.Project {
			return false
		}
		if nr.Location == nil && ncr.Location == nil {
			c.Config.Logger.Info("Both Location fields null - considering equal.")
		} else if nr.Location == nil || ncr.Location == nil {
			c.Config.Logger.Info("Only one Location field is null - considering unequal.")
			return false
		} else if *nr.Location != *ncr.Location {
			return false
		}
		if nr.FirewallPolicy == nil && ncr.FirewallPolicy == nil {
			c.Config.Logger.Info("Both FirewallPolicy fields null - considering equal.")
		} else if nr.FirewallPolicy == nil || ncr.FirewallPolicy == nil {
			c.Config.Logger.Info("Only one FirewallPolicy field is null - considering unequal.")
			return false
		} else if *nr.FirewallPolicy != *ncr.FirewallPolicy {
			return false
		}
		if nr.Priority == nil && ncr.Priority == nil {
			c.Config.Logger.Info("Both Priority fields null - considering equal.")
		} else if nr.Priority == nil || ncr.Priority == nil {
			c.Config.Logger.Info("Only one Priority field is null - considering unequal.")
			return false
		} else if *nr.Priority != *ncr.Priority {
			return false
		}
		return true
	}
}

type networkFirewallPolicyRuleDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         networkFirewallPolicyRuleApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToNetworkFirewallPolicyRuleDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]networkFirewallPolicyRuleDiff, error) {
	opNamesToFieldDiffs := make(map[string][]*dcl.FieldDiff)
	// Map each operation name to the field diffs associated with it.
	for _, fd := range fds {
		for _, ro := range fd.ResultingOperation {
			if fieldDiffs, ok := opNamesToFieldDiffs[ro]; ok {
				fieldDiffs = append(fieldDiffs, fd)
				opNamesToFieldDiffs[ro] = fieldDiffs
			} else {
				config.Logger.Infof("%s required due to diff: %v", ro, fd)
				opNamesToFieldDiffs[ro] = []*dcl.FieldDiff{fd}
			}
		}
	}
	var diffs []networkFirewallPolicyRuleDiff
	// For each operation name, create a networkFirewallPolicyRuleDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := networkFirewallPolicyRuleDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToNetworkFirewallPolicyRuleApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToNetworkFirewallPolicyRuleApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (networkFirewallPolicyRuleApiOperation, error) {
	switch opName {

	case "updateNetworkFirewallPolicyRulePatchRuleOperation":
		return &updateNetworkFirewallPolicyRulePatchRuleOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractNetworkFirewallPolicyRuleFields(r *NetworkFirewallPolicyRule) error {
	vMatch := r.Match
	if vMatch == nil {
		// note: explicitly not the empty object.
		vMatch = &NetworkFirewallPolicyRuleMatch{}
	}
	if err := extractNetworkFirewallPolicyRuleMatchFields(r, vMatch); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMatch) {
		r.Match = vMatch
	}
	return nil
}
func extractNetworkFirewallPolicyRuleMatchFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatch) error {
	return nil
}
func extractNetworkFirewallPolicyRuleMatchLayer4ConfigsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatchLayer4Configs) error {
	return nil
}
func extractNetworkFirewallPolicyRuleMatchSrcSecureTagsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatchSrcSecureTags) error {
	return nil
}
func extractNetworkFirewallPolicyRuleTargetSecureTagsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleTargetSecureTags) error {
	return nil
}

func postReadExtractNetworkFirewallPolicyRuleFields(r *NetworkFirewallPolicyRule) error {
	vMatch := r.Match
	if vMatch == nil {
		// note: explicitly not the empty object.
		vMatch = &NetworkFirewallPolicyRuleMatch{}
	}
	if err := postReadExtractNetworkFirewallPolicyRuleMatchFields(r, vMatch); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMatch) {
		r.Match = vMatch
	}
	return nil
}
func postReadExtractNetworkFirewallPolicyRuleMatchFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatch) error {
	return nil
}
func postReadExtractNetworkFirewallPolicyRuleMatchLayer4ConfigsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatchLayer4Configs) error {
	return nil
}
func postReadExtractNetworkFirewallPolicyRuleMatchSrcSecureTagsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleMatchSrcSecureTags) error {
	return nil
}
func postReadExtractNetworkFirewallPolicyRuleTargetSecureTagsFields(r *NetworkFirewallPolicyRule, o *NetworkFirewallPolicyRuleTargetSecureTags) error {
	return nil
}
