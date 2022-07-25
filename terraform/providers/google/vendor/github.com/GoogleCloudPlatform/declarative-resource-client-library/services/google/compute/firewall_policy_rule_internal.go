// Copyright 2021 Google LLC. All Rights Reserved.
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
)

func (r *FirewallPolicyRule) validate() error {

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
	if !dcl.IsEmptyValueIndirect(r.Match) {
		if err := r.Match.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *FirewallPolicyRuleMatch) validate() error {
	if err := dcl.Required(r, "layer4Configs"); err != nil {
		return err
	}
	return nil
}
func (r *FirewallPolicyRuleMatchLayer4Configs) validate() error {
	if err := dcl.Required(r, "ipProtocol"); err != nil {
		return err
	}
	return nil
}
func (r *FirewallPolicyRule) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *FirewallPolicyRule) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"priority":       dcl.ValueOrEmptyString(nr.Priority),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/getRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
}

func (r *FirewallPolicyRule) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicyRule) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/addRule", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicyRule) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"priority":       dcl.ValueOrEmptyString(nr.Priority),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/removeRule?priority={{priority}}", nr.basePath(), userBasePath, params), nil
}

// firewallPolicyRuleApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type firewallPolicyRuleApiOperation interface {
	do(context.Context, *FirewallPolicyRule, *Client) error
}

// newUpdateFirewallPolicyRulePatchRuleRequest creates a request for an
// FirewallPolicyRule resource's PatchRule update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateFirewallPolicyRulePatchRuleRequest(ctx context.Context, f *FirewallPolicyRule, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandFirewallPolicyRuleMatch(c, f.Match); err != nil {
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
	if v := f.TargetResources; v != nil {
		req["targetResources"] = v
	}
	if v := f.EnableLogging; !dcl.IsEmptyValueIndirect(v) {
		req["enableLogging"] = v
	}
	if v := f.RuleTupleCount; !dcl.IsEmptyValueIndirect(v) {
		req["ruleTupleCount"] = v
	}
	if v := f.TargetServiceAccounts; v != nil {
		req["targetServiceAccounts"] = v
	}
	if v := f.Disabled; !dcl.IsEmptyValueIndirect(v) {
		req["disabled"] = v
	}
	return req, nil
}

// marshalUpdateFirewallPolicyRulePatchRuleRequest converts the update into
// the final JSON request body.
func marshalUpdateFirewallPolicyRulePatchRuleRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateFirewallPolicyRulePatchRuleOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (c *Client) listFirewallPolicyRuleRaw(ctx context.Context, r *FirewallPolicyRule, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != FirewallPolicyRuleMaxPage {
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

type listFirewallPolicyRuleOperation struct {
	Rules []map[string]interface{} `json:"rules"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listFirewallPolicyRule(ctx context.Context, r *FirewallPolicyRule, pageToken string, pageSize int32) ([]*FirewallPolicyRule, string, error) {
	b, err := c.listFirewallPolicyRuleRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listFirewallPolicyRuleOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*FirewallPolicyRule
	for _, v := range m.Rules {
		res, err := unmarshalMapFirewallPolicyRule(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.FirewallPolicy = r.FirewallPolicy
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllFirewallPolicyRule(ctx context.Context, f func(*FirewallPolicyRule) bool, resources []*FirewallPolicyRule) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteFirewallPolicyRule(ctx, res)
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

type deleteFirewallPolicyRuleOperation struct{}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createFirewallPolicyRuleOperation struct {
	response map[string]interface{}
}

func (op *createFirewallPolicyRuleOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (c *Client) getFirewallPolicyRuleRaw(ctx context.Context, r *FirewallPolicyRule) ([]byte, error) {

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

func (c *Client) firewallPolicyRuleDiffsForRawDesired(ctx context.Context, rawDesired *FirewallPolicyRule, opts ...dcl.ApplyOption) (initial, desired *FirewallPolicyRule, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *FirewallPolicyRule
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*FirewallPolicyRule); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected FirewallPolicyRule, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetFirewallPolicyRule(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a FirewallPolicyRule resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve FirewallPolicyRule resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that FirewallPolicyRule resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeFirewallPolicyRuleDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for FirewallPolicyRule: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for FirewallPolicyRule: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeFirewallPolicyRuleInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for FirewallPolicyRule: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeFirewallPolicyRuleDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for FirewallPolicyRule: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffFirewallPolicyRule(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeFirewallPolicyRuleInitialState(rawInitial, rawDesired *FirewallPolicyRule) (*FirewallPolicyRule, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeFirewallPolicyRuleDesiredState(rawDesired, rawInitial *FirewallPolicyRule, opts ...dcl.ApplyOption) (*FirewallPolicyRule, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Match = canonicalizeFirewallPolicyRuleMatch(rawDesired.Match, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &FirewallPolicyRule{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.Priority) {
		canonicalDesired.Priority = rawInitial.Priority
	} else {
		canonicalDesired.Priority = rawDesired.Priority
	}
	canonicalDesired.Match = canonicalizeFirewallPolicyRuleMatch(rawDesired.Match, rawInitial.Match, opts...)
	if dcl.StringCanonicalize(rawDesired.Action, rawInitial.Action) {
		canonicalDesired.Action = rawInitial.Action
	} else {
		canonicalDesired.Action = rawDesired.Action
	}
	if dcl.IsZeroValue(rawDesired.Direction) {
		canonicalDesired.Direction = rawInitial.Direction
	} else {
		canonicalDesired.Direction = rawDesired.Direction
	}
	if dcl.StringArrayCanonicalize(rawDesired.TargetResources, rawInitial.TargetResources) {
		canonicalDesired.TargetResources = rawInitial.TargetResources
	} else {
		canonicalDesired.TargetResources = rawDesired.TargetResources
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
	if dcl.BoolCanonicalize(rawDesired.Disabled, rawInitial.Disabled) {
		canonicalDesired.Disabled = rawInitial.Disabled
	} else {
		canonicalDesired.Disabled = rawDesired.Disabled
	}
	if dcl.NameToSelfLink(rawDesired.FirewallPolicy, rawInitial.FirewallPolicy) {
		canonicalDesired.FirewallPolicy = rawInitial.FirewallPolicy
	} else {
		canonicalDesired.FirewallPolicy = rawDesired.FirewallPolicy
	}

	return canonicalDesired, nil
}

func canonicalizeFirewallPolicyRuleNewState(c *Client, rawNew, rawDesired *FirewallPolicyRule) (*FirewallPolicyRule, error) {

	if dcl.IsNotReturnedByServer(rawNew.Description) && dcl.IsNotReturnedByServer(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Priority) && dcl.IsNotReturnedByServer(rawDesired.Priority) {
		rawNew.Priority = rawDesired.Priority
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Match) && dcl.IsNotReturnedByServer(rawDesired.Match) {
		rawNew.Match = rawDesired.Match
	} else {
		rawNew.Match = canonicalizeNewFirewallPolicyRuleMatch(c, rawDesired.Match, rawNew.Match)
	}

	if dcl.IsNotReturnedByServer(rawNew.Action) && dcl.IsNotReturnedByServer(rawDesired.Action) {
		rawNew.Action = rawDesired.Action
	} else {
		if dcl.StringCanonicalize(rawDesired.Action, rawNew.Action) {
			rawNew.Action = rawDesired.Action
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Direction) && dcl.IsNotReturnedByServer(rawDesired.Direction) {
		rawNew.Direction = rawDesired.Direction
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.TargetResources) && dcl.IsNotReturnedByServer(rawDesired.TargetResources) {
		rawNew.TargetResources = rawDesired.TargetResources
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.TargetResources, rawNew.TargetResources) {
			rawNew.TargetResources = rawDesired.TargetResources
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.EnableLogging) && dcl.IsNotReturnedByServer(rawDesired.EnableLogging) {
		rawNew.EnableLogging = rawDesired.EnableLogging
	} else {
		if dcl.BoolCanonicalize(rawDesired.EnableLogging, rawNew.EnableLogging) {
			rawNew.EnableLogging = rawDesired.EnableLogging
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.RuleTupleCount) && dcl.IsNotReturnedByServer(rawDesired.RuleTupleCount) {
		rawNew.RuleTupleCount = rawDesired.RuleTupleCount
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.TargetServiceAccounts) && dcl.IsNotReturnedByServer(rawDesired.TargetServiceAccounts) {
		rawNew.TargetServiceAccounts = rawDesired.TargetServiceAccounts
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.TargetServiceAccounts, rawNew.TargetServiceAccounts) {
			rawNew.TargetServiceAccounts = rawDesired.TargetServiceAccounts
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Disabled) && dcl.IsNotReturnedByServer(rawDesired.Disabled) {
		rawNew.Disabled = rawDesired.Disabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.Disabled, rawNew.Disabled) {
			rawNew.Disabled = rawDesired.Disabled
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Kind) && dcl.IsNotReturnedByServer(rawDesired.Kind) {
		rawNew.Kind = rawDesired.Kind
	} else {
		if dcl.StringCanonicalize(rawDesired.Kind, rawNew.Kind) {
			rawNew.Kind = rawDesired.Kind
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.FirewallPolicy) && dcl.IsNotReturnedByServer(rawDesired.FirewallPolicy) {
		rawNew.FirewallPolicy = rawDesired.FirewallPolicy
	} else {
		if dcl.NameToSelfLink(rawDesired.FirewallPolicy, rawNew.FirewallPolicy) {
			rawNew.FirewallPolicy = rawDesired.FirewallPolicy
		}
	}

	return rawNew, nil
}

func canonicalizeFirewallPolicyRuleMatch(des, initial *FirewallPolicyRuleMatch, opts ...dcl.ApplyOption) *FirewallPolicyRuleMatch {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &FirewallPolicyRuleMatch{}

	if dcl.StringArrayCanonicalize(des.SrcIPRanges, initial.SrcIPRanges) || dcl.IsZeroValue(des.SrcIPRanges) {
		cDes.SrcIPRanges = initial.SrcIPRanges
	} else {
		cDes.SrcIPRanges = des.SrcIPRanges
	}
	if dcl.StringArrayCanonicalize(des.DestIPRanges, initial.DestIPRanges) || dcl.IsZeroValue(des.DestIPRanges) {
		cDes.DestIPRanges = initial.DestIPRanges
	} else {
		cDes.DestIPRanges = des.DestIPRanges
	}
	cDes.Layer4Configs = canonicalizeFirewallPolicyRuleMatchLayer4ConfigsSlice(des.Layer4Configs, initial.Layer4Configs, opts...)

	return cDes
}

func canonicalizeFirewallPolicyRuleMatchSlice(des, initial []FirewallPolicyRuleMatch, opts ...dcl.ApplyOption) []FirewallPolicyRuleMatch {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]FirewallPolicyRuleMatch, 0, len(des))
		for _, d := range des {
			cd := canonicalizeFirewallPolicyRuleMatch(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]FirewallPolicyRuleMatch, 0, len(des))
	for i, d := range des {
		cd := canonicalizeFirewallPolicyRuleMatch(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewFirewallPolicyRuleMatch(c *Client, des, nw *FirewallPolicyRuleMatch) *FirewallPolicyRuleMatch {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for FirewallPolicyRuleMatch while comparing non-nil desired to nil actual.  Returning desired object.")
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
	nw.Layer4Configs = canonicalizeNewFirewallPolicyRuleMatchLayer4ConfigsSlice(c, des.Layer4Configs, nw.Layer4Configs)

	return nw
}

func canonicalizeNewFirewallPolicyRuleMatchSet(c *Client, des, nw []FirewallPolicyRuleMatch) []FirewallPolicyRuleMatch {
	if des == nil {
		return nw
	}
	var reorderedNew []FirewallPolicyRuleMatch
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareFirewallPolicyRuleMatchNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewFirewallPolicyRuleMatchSlice(c *Client, des, nw []FirewallPolicyRuleMatch) []FirewallPolicyRuleMatch {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []FirewallPolicyRuleMatch
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewFirewallPolicyRuleMatch(c, &d, &n))
	}

	return items
}

func canonicalizeFirewallPolicyRuleMatchLayer4Configs(des, initial *FirewallPolicyRuleMatchLayer4Configs, opts ...dcl.ApplyOption) *FirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &FirewallPolicyRuleMatchLayer4Configs{}

	if dcl.StringCanonicalize(des.IPProtocol, initial.IPProtocol) || dcl.IsZeroValue(des.IPProtocol) {
		cDes.IPProtocol = initial.IPProtocol
	} else {
		cDes.IPProtocol = des.IPProtocol
	}
	if dcl.StringArrayCanonicalize(des.Ports, initial.Ports) || dcl.IsZeroValue(des.Ports) {
		cDes.Ports = initial.Ports
	} else {
		cDes.Ports = des.Ports
	}

	return cDes
}

func canonicalizeFirewallPolicyRuleMatchLayer4ConfigsSlice(des, initial []FirewallPolicyRuleMatchLayer4Configs, opts ...dcl.ApplyOption) []FirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]FirewallPolicyRuleMatchLayer4Configs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeFirewallPolicyRuleMatchLayer4Configs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]FirewallPolicyRuleMatchLayer4Configs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeFirewallPolicyRuleMatchLayer4Configs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewFirewallPolicyRuleMatchLayer4Configs(c *Client, des, nw *FirewallPolicyRuleMatchLayer4Configs) *FirewallPolicyRuleMatchLayer4Configs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for FirewallPolicyRuleMatchLayer4Configs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewFirewallPolicyRuleMatchLayer4ConfigsSet(c *Client, des, nw []FirewallPolicyRuleMatchLayer4Configs) []FirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return nw
	}
	var reorderedNew []FirewallPolicyRuleMatchLayer4Configs
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareFirewallPolicyRuleMatchLayer4ConfigsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, des, nw []FirewallPolicyRuleMatchLayer4Configs) []FirewallPolicyRuleMatchLayer4Configs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []FirewallPolicyRuleMatchLayer4Configs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewFirewallPolicyRuleMatchLayer4Configs(c, &d, &n))
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
func diffFirewallPolicyRule(c *Client, desired, actual *FirewallPolicyRule, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Priority, actual.Priority, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Priority")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Match, actual.Match, dcl.Info{ObjectFunction: compareFirewallPolicyRuleMatchNewStyle, EmptyObject: EmptyFirewallPolicyRuleMatch, OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Match")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Action, actual.Action, dcl.Info{OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Action")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Direction, actual.Direction, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Direction")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetResources, actual.TargetResources, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("TargetResources")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableLogging, actual.EnableLogging, dcl.Info{OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("EnableLogging")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RuleTupleCount, actual.RuleTupleCount, dcl.Info{OutputOnly: true, OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("RuleTupleCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetServiceAccounts, actual.TargetServiceAccounts, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("TargetServiceAccounts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Disabled, actual.Disabled, dcl.Info{OperationSelector: dcl.TriggersOperation("updateFirewallPolicyRulePatchRuleOperation")}, fn.AddNest("Disabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Kind, actual.Kind, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Kind")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FirewallPolicy, actual.FirewallPolicy, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FirewallPolicy")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	return newDiffs, nil
}
func compareFirewallPolicyRuleMatchNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*FirewallPolicyRuleMatch)
	if !ok {
		desiredNotPointer, ok := d.(FirewallPolicyRuleMatch)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a FirewallPolicyRuleMatch or *FirewallPolicyRuleMatch", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*FirewallPolicyRuleMatch)
	if !ok {
		actualNotPointer, ok := a.(FirewallPolicyRuleMatch)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a FirewallPolicyRuleMatch", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SrcIPRanges, actual.SrcIPRanges, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SrcIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestIPRanges, actual.DestIPRanges, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DestIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Layer4Configs, actual.Layer4Configs, dcl.Info{ObjectFunction: compareFirewallPolicyRuleMatchLayer4ConfigsNewStyle, EmptyObject: EmptyFirewallPolicyRuleMatchLayer4Configs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Layer4Configs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareFirewallPolicyRuleMatchLayer4ConfigsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*FirewallPolicyRuleMatchLayer4Configs)
	if !ok {
		desiredNotPointer, ok := d.(FirewallPolicyRuleMatchLayer4Configs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a FirewallPolicyRuleMatchLayer4Configs or *FirewallPolicyRuleMatchLayer4Configs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*FirewallPolicyRuleMatchLayer4Configs)
	if !ok {
		actualNotPointer, ok := a.(FirewallPolicyRuleMatchLayer4Configs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a FirewallPolicyRuleMatchLayer4Configs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IPProtocol, actual.IPProtocol, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpProtocol")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Ports, actual.Ports, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Ports")); len(ds) != 0 || err != nil {
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
func (r *FirewallPolicyRule) urlNormalized() *FirewallPolicyRule {
	normalized := dcl.Copy(*r).(FirewallPolicyRule)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Action = dcl.SelfLinkToName(r.Action)
	normalized.Kind = dcl.SelfLinkToName(r.Kind)
	normalized.FirewallPolicy = dcl.SelfLinkToName(r.FirewallPolicy)
	return &normalized
}

func (r *FirewallPolicyRule) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "PatchRule" {
		fields := map[string]interface{}{
			"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
			"priority":       dcl.ValueOrEmptyString(nr.Priority),
		}
		return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/patchRule?priority={{priority}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the FirewallPolicyRule resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *FirewallPolicyRule) marshal(c *Client) ([]byte, error) {
	m, err := expandFirewallPolicyRule(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling FirewallPolicyRule: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalFirewallPolicyRule decodes JSON responses into the FirewallPolicyRule resource schema.
func unmarshalFirewallPolicyRule(b []byte, c *Client) (*FirewallPolicyRule, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapFirewallPolicyRule(m, c)
}

func unmarshalMapFirewallPolicyRule(m map[string]interface{}, c *Client) (*FirewallPolicyRule, error) {

	flattened := flattenFirewallPolicyRule(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandFirewallPolicyRule expands FirewallPolicyRule into a JSON request object.
func expandFirewallPolicyRule(c *Client, f *FirewallPolicyRule) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Priority; dcl.ValueShouldBeSent(v) {
		m["priority"] = v
	}
	if v, err := expandFirewallPolicyRuleMatch(c, f.Match); err != nil {
		return nil, fmt.Errorf("error expanding Match into match: %w", err)
	} else if v != nil {
		m["match"] = v
	}
	if v := f.Action; dcl.ValueShouldBeSent(v) {
		m["action"] = v
	}
	if v := f.Direction; dcl.ValueShouldBeSent(v) {
		m["direction"] = v
	}
	m["targetResources"] = f.TargetResources
	if v := f.EnableLogging; dcl.ValueShouldBeSent(v) {
		m["enableLogging"] = v
	}
	m["targetServiceAccounts"] = f.TargetServiceAccounts
	if v := f.Disabled; dcl.ValueShouldBeSent(v) {
		m["disabled"] = v
	}
	if v := f.FirewallPolicy; dcl.ValueShouldBeSent(v) {
		m["firewallPolicy"] = v
	}

	return m, nil
}

// flattenFirewallPolicyRule flattens FirewallPolicyRule from a JSON request object into the
// FirewallPolicyRule type.
func flattenFirewallPolicyRule(c *Client, i interface{}) *FirewallPolicyRule {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &FirewallPolicyRule{}
	res.Description = dcl.FlattenString(m["description"])
	res.Priority = dcl.FlattenInteger(m["priority"])
	res.Match = flattenFirewallPolicyRuleMatch(c, m["match"])
	res.Action = dcl.FlattenString(m["action"])
	res.Direction = flattenFirewallPolicyRuleDirectionEnum(m["direction"])
	res.TargetResources = dcl.FlattenStringSlice(m["targetResources"])
	res.EnableLogging = dcl.FlattenBool(m["enableLogging"])
	res.RuleTupleCount = dcl.FlattenInteger(m["ruleTupleCount"])
	res.TargetServiceAccounts = dcl.FlattenStringSlice(m["targetServiceAccounts"])
	res.Disabled = dcl.FlattenBool(m["disabled"])
	res.Kind = dcl.FlattenString(m["kind"])
	res.FirewallPolicy = dcl.FlattenString(m["firewallPolicy"])

	return res
}

// expandFirewallPolicyRuleMatchMap expands the contents of FirewallPolicyRuleMatch into a JSON
// request object.
func expandFirewallPolicyRuleMatchMap(c *Client, f map[string]FirewallPolicyRuleMatch) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandFirewallPolicyRuleMatch(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandFirewallPolicyRuleMatchSlice expands the contents of FirewallPolicyRuleMatch into a JSON
// request object.
func expandFirewallPolicyRuleMatchSlice(c *Client, f []FirewallPolicyRuleMatch) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandFirewallPolicyRuleMatch(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenFirewallPolicyRuleMatchMap flattens the contents of FirewallPolicyRuleMatch from a JSON
// response object.
func flattenFirewallPolicyRuleMatchMap(c *Client, i interface{}) map[string]FirewallPolicyRuleMatch {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]FirewallPolicyRuleMatch{}
	}

	if len(a) == 0 {
		return map[string]FirewallPolicyRuleMatch{}
	}

	items := make(map[string]FirewallPolicyRuleMatch)
	for k, item := range a {
		items[k] = *flattenFirewallPolicyRuleMatch(c, item.(map[string]interface{}))
	}

	return items
}

// flattenFirewallPolicyRuleMatchSlice flattens the contents of FirewallPolicyRuleMatch from a JSON
// response object.
func flattenFirewallPolicyRuleMatchSlice(c *Client, i interface{}) []FirewallPolicyRuleMatch {
	a, ok := i.([]interface{})
	if !ok {
		return []FirewallPolicyRuleMatch{}
	}

	if len(a) == 0 {
		return []FirewallPolicyRuleMatch{}
	}

	items := make([]FirewallPolicyRuleMatch, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenFirewallPolicyRuleMatch(c, item.(map[string]interface{})))
	}

	return items
}

// expandFirewallPolicyRuleMatch expands an instance of FirewallPolicyRuleMatch into a JSON
// request object.
func expandFirewallPolicyRuleMatch(c *Client, f *FirewallPolicyRuleMatch) (map[string]interface{}, error) {
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
	if v, err := expandFirewallPolicyRuleMatchLayer4ConfigsSlice(c, f.Layer4Configs); err != nil {
		return nil, fmt.Errorf("error expanding Layer4Configs into layer4Configs: %w", err)
	} else if v != nil {
		m["layer4Configs"] = v
	}

	return m, nil
}

// flattenFirewallPolicyRuleMatch flattens an instance of FirewallPolicyRuleMatch from a JSON
// response object.
func flattenFirewallPolicyRuleMatch(c *Client, i interface{}) *FirewallPolicyRuleMatch {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &FirewallPolicyRuleMatch{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyFirewallPolicyRuleMatch
	}
	r.SrcIPRanges = dcl.FlattenStringSlice(m["srcIpRanges"])
	r.DestIPRanges = dcl.FlattenStringSlice(m["destIpRanges"])
	r.Layer4Configs = flattenFirewallPolicyRuleMatchLayer4ConfigsSlice(c, m["layer4Configs"])

	return r
}

// expandFirewallPolicyRuleMatchLayer4ConfigsMap expands the contents of FirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandFirewallPolicyRuleMatchLayer4ConfigsMap(c *Client, f map[string]FirewallPolicyRuleMatchLayer4Configs) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandFirewallPolicyRuleMatchLayer4Configs(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandFirewallPolicyRuleMatchLayer4ConfigsSlice expands the contents of FirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, f []FirewallPolicyRuleMatchLayer4Configs) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandFirewallPolicyRuleMatchLayer4Configs(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenFirewallPolicyRuleMatchLayer4ConfigsMap flattens the contents of FirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenFirewallPolicyRuleMatchLayer4ConfigsMap(c *Client, i interface{}) map[string]FirewallPolicyRuleMatchLayer4Configs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]FirewallPolicyRuleMatchLayer4Configs{}
	}

	if len(a) == 0 {
		return map[string]FirewallPolicyRuleMatchLayer4Configs{}
	}

	items := make(map[string]FirewallPolicyRuleMatchLayer4Configs)
	for k, item := range a {
		items[k] = *flattenFirewallPolicyRuleMatchLayer4Configs(c, item.(map[string]interface{}))
	}

	return items
}

// flattenFirewallPolicyRuleMatchLayer4ConfigsSlice flattens the contents of FirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenFirewallPolicyRuleMatchLayer4ConfigsSlice(c *Client, i interface{}) []FirewallPolicyRuleMatchLayer4Configs {
	a, ok := i.([]interface{})
	if !ok {
		return []FirewallPolicyRuleMatchLayer4Configs{}
	}

	if len(a) == 0 {
		return []FirewallPolicyRuleMatchLayer4Configs{}
	}

	items := make([]FirewallPolicyRuleMatchLayer4Configs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenFirewallPolicyRuleMatchLayer4Configs(c, item.(map[string]interface{})))
	}

	return items
}

// expandFirewallPolicyRuleMatchLayer4Configs expands an instance of FirewallPolicyRuleMatchLayer4Configs into a JSON
// request object.
func expandFirewallPolicyRuleMatchLayer4Configs(c *Client, f *FirewallPolicyRuleMatchLayer4Configs) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
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

// flattenFirewallPolicyRuleMatchLayer4Configs flattens an instance of FirewallPolicyRuleMatchLayer4Configs from a JSON
// response object.
func flattenFirewallPolicyRuleMatchLayer4Configs(c *Client, i interface{}) *FirewallPolicyRuleMatchLayer4Configs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &FirewallPolicyRuleMatchLayer4Configs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyFirewallPolicyRuleMatchLayer4Configs
	}
	r.IPProtocol = dcl.FlattenString(m["ipProtocol"])
	r.Ports = dcl.FlattenStringSlice(m["ports"])

	return r
}

// flattenFirewallPolicyRuleDirectionEnumMap flattens the contents of FirewallPolicyRuleDirectionEnum from a JSON
// response object.
func flattenFirewallPolicyRuleDirectionEnumMap(c *Client, i interface{}) map[string]FirewallPolicyRuleDirectionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]FirewallPolicyRuleDirectionEnum{}
	}

	if len(a) == 0 {
		return map[string]FirewallPolicyRuleDirectionEnum{}
	}

	items := make(map[string]FirewallPolicyRuleDirectionEnum)
	for k, item := range a {
		items[k] = *flattenFirewallPolicyRuleDirectionEnum(item.(interface{}))
	}

	return items
}

// flattenFirewallPolicyRuleDirectionEnumSlice flattens the contents of FirewallPolicyRuleDirectionEnum from a JSON
// response object.
func flattenFirewallPolicyRuleDirectionEnumSlice(c *Client, i interface{}) []FirewallPolicyRuleDirectionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []FirewallPolicyRuleDirectionEnum{}
	}

	if len(a) == 0 {
		return []FirewallPolicyRuleDirectionEnum{}
	}

	items := make([]FirewallPolicyRuleDirectionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenFirewallPolicyRuleDirectionEnum(item.(interface{})))
	}

	return items
}

// flattenFirewallPolicyRuleDirectionEnum asserts that an interface is a string, and returns a
// pointer to a *FirewallPolicyRuleDirectionEnum with the same value as that string.
func flattenFirewallPolicyRuleDirectionEnum(i interface{}) *FirewallPolicyRuleDirectionEnum {
	s, ok := i.(string)
	if !ok {
		return FirewallPolicyRuleDirectionEnumRef("")
	}

	return FirewallPolicyRuleDirectionEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *FirewallPolicyRule) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalFirewallPolicyRule(b, c)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

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

type firewallPolicyRuleDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         firewallPolicyRuleApiOperation
}

func convertFieldDiffsToFirewallPolicyRuleDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]firewallPolicyRuleDiff, error) {
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
	var diffs []firewallPolicyRuleDiff
	// For each operation name, create a firewallPolicyRuleDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := firewallPolicyRuleDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToFirewallPolicyRuleApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToFirewallPolicyRuleApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (firewallPolicyRuleApiOperation, error) {
	switch opName {

	case "updateFirewallPolicyRulePatchRuleOperation":
		return &updateFirewallPolicyRulePatchRuleOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractFirewallPolicyRuleFields(r *FirewallPolicyRule) error {
	vMatch := r.Match
	if vMatch == nil {
		// note: explicitly not the empty object.
		vMatch = &FirewallPolicyRuleMatch{}
	}
	if err := extractFirewallPolicyRuleMatchFields(r, vMatch); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMatch) {
		r.Match = vMatch
	}
	return nil
}
func extractFirewallPolicyRuleMatchFields(r *FirewallPolicyRule, o *FirewallPolicyRuleMatch) error {
	return nil
}
func extractFirewallPolicyRuleMatchLayer4ConfigsFields(r *FirewallPolicyRule, o *FirewallPolicyRuleMatchLayer4Configs) error {
	return nil
}

func postReadExtractFirewallPolicyRuleFields(r *FirewallPolicyRule) error {
	vMatch := r.Match
	if vMatch == nil {
		// note: explicitly not the empty object.
		vMatch = &FirewallPolicyRuleMatch{}
	}
	if err := postReadExtractFirewallPolicyRuleMatchFields(r, vMatch); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMatch) {
		r.Match = vMatch
	}
	return nil
}
func postReadExtractFirewallPolicyRuleMatchFields(r *FirewallPolicyRule, o *FirewallPolicyRuleMatch) error {
	return nil
}
func postReadExtractFirewallPolicyRuleMatchLayer4ConfigsFields(r *FirewallPolicyRule, o *FirewallPolicyRuleMatchLayer4Configs) error {
	return nil
}
