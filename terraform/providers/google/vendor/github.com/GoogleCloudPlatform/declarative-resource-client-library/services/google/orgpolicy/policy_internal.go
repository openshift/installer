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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Policy) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Parent, "Parent"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Spec) {
		if err := r.Spec.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.DryRunSpec) {
		if err := r.DryRunSpec.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PolicySpec) validate() error {
	return nil
}
func (r *PolicySpecRules) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Values", "AllowAll", "DenyAll", "Enforce"}, r.Values, r.AllowAll, r.DenyAll, r.Enforce); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Values) {
		if err := r.Values.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Condition) {
		if err := r.Condition.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PolicySpecRulesValues) validate() error {
	return nil
}
func (r *PolicySpecRulesCondition) validate() error {
	return nil
}
func (r *PolicyDryRunSpec) validate() error {
	return nil
}
func (r *PolicyDryRunSpecRules) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Values", "AllowAll", "DenyAll", "Enforce"}, r.Values, r.AllowAll, r.DenyAll, r.Enforce); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Values) {
		if err := r.Values.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Condition) {
		if err := r.Condition.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PolicyDryRunSpecRulesValues) validate() error {
	return nil
}
func (r *PolicyDryRunSpecRulesCondition) validate() error {
	return nil
}
func (r *Policy) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://orgpolicy.googleapis.com/v2/", params)
}

func (r *Policy) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
		"name":   dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("{{parent}}/policies/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Policy) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("{{parent}}/policies", nr.basePath(), userBasePath, params), nil

}

func (r *Policy) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("{{parent}}/policies", nr.basePath(), userBasePath, params), nil

}

func (r *Policy) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
		"name":   dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("{{parent}}/policies/{{name}}", nr.basePath(), userBasePath, params), nil
}

// policyApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type policyApiOperation interface {
	do(context.Context, *Policy, *Client) error
}

// newUpdatePolicyUpdatePolicyRequest creates a request for an
// Policy resource's UpdatePolicy update type by filling in the update
// fields based on the intended state of the resource.
func newUpdatePolicyUpdatePolicyRequest(ctx context.Context, f *Policy, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandPolicySpec(c, f.Spec, res); err != nil {
		return nil, fmt.Errorf("error expanding Spec into spec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["spec"] = v
	}
	if v, err := expandPolicyDryRunSpec(c, f.DryRunSpec, res); err != nil {
		return nil, fmt.Errorf("error expanding DryRunSpec into dryRunSpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["dryRunSpec"] = v
	}
	b, err := c.getPolicyRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawEtag, err := dcl.GetMapEntry(
		m,
		[]string{"etag"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["etag"] = rawEtag.(string)
	}
	return req, nil
}

// marshalUpdatePolicyUpdatePolicyRequest converts the update into
// the final JSON request body.
func marshalUpdatePolicyUpdatePolicyRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updatePolicyUpdatePolicyOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (c *Client) listPolicyRaw(ctx context.Context, r *Policy, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != PolicyMaxPage {
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

type listPolicyOperation struct {
	Policies []map[string]interface{} `json:"policies"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listPolicy(ctx context.Context, r *Policy, pageToken string, pageSize int32) ([]*Policy, string, error) {
	b, err := c.listPolicyRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listPolicyOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Policy
	for _, v := range m.Policies {
		res, err := unmarshalMapPolicy(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Parent = r.Parent
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllPolicy(ctx context.Context, f func(*Policy) bool, resources []*Policy) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeletePolicy(ctx, res)
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

type deletePolicyOperation struct{}

func (op *deletePolicyOperation) do(ctx context.Context, r *Policy, c *Client) error {
	r, err := c.GetPolicy(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Policy not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetPolicy checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete Policy: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createPolicyOperation struct {
	response map[string]interface{}
}

func (op *createPolicyOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createPolicyOperation) do(ctx context.Context, r *Policy, c *Client) error {
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

	o, err := dcl.ResponseBodyAsJSON(resp)
	if err != nil {
		return fmt.Errorf("error decoding response body into JSON: %w", err)
	}
	op.response = o

	if _, err := c.GetPolicy(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getPolicyRaw(ctx context.Context, r *Policy) ([]byte, error) {

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

func (c *Client) policyDiffsForRawDesired(ctx context.Context, rawDesired *Policy, opts ...dcl.ApplyOption) (initial, desired *Policy, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Policy
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Policy); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Policy, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetPolicy(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Policy resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Policy resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Policy resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizePolicyDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Policy: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Policy: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractPolicyFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizePolicyInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Policy: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizePolicyDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Policy: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffPolicy(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizePolicyInitialState(rawInitial, rawDesired *Policy) (*Policy, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizePolicyDesiredState(rawDesired, rawInitial *Policy, opts ...dcl.ApplyOption) (*Policy, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Spec = canonicalizePolicySpec(rawDesired.Spec, nil, opts...)
		rawDesired.DryRunSpec = canonicalizePolicyDryRunSpec(rawDesired.DryRunSpec, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Policy{}
	if canonicalizePolicyName(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.Spec = canonicalizePolicySpec(rawDesired.Spec, rawInitial.Spec, opts...)
	canonicalDesired.DryRunSpec = canonicalizePolicyDryRunSpec(rawDesired.DryRunSpec, rawInitial.DryRunSpec, opts...)
	if dcl.NameToSelfLink(rawDesired.Parent, rawInitial.Parent) {
		canonicalDesired.Parent = rawInitial.Parent
	} else {
		canonicalDesired.Parent = rawDesired.Parent
	}
	return canonicalDesired, nil
}

func canonicalizePolicyNewState(c *Client, rawNew, rawDesired *Policy) (*Policy, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if canonicalizePolicyName(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Spec) && dcl.IsEmptyValueIndirect(rawDesired.Spec) {
		rawNew.Spec = rawDesired.Spec
	} else {
		rawNew.Spec = canonicalizeNewPolicySpec(c, rawDesired.Spec, rawNew.Spec)
	}

	if dcl.IsEmptyValueIndirect(rawNew.DryRunSpec) && dcl.IsEmptyValueIndirect(rawDesired.DryRunSpec) {
		rawNew.DryRunSpec = rawDesired.DryRunSpec
	} else {
		rawNew.DryRunSpec = canonicalizeNewPolicyDryRunSpec(c, rawDesired.DryRunSpec, rawNew.DryRunSpec)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	rawNew.Parent = rawDesired.Parent

	return rawNew, nil
}

func canonicalizePolicySpec(des, initial *PolicySpec, opts ...dcl.ApplyOption) *PolicySpec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicySpec{}

	cDes.Rules = canonicalizePolicySpecRulesSlice(des.Rules, initial.Rules, opts...)
	if dcl.BoolCanonicalize(des.InheritFromParent, initial.InheritFromParent) || dcl.IsZeroValue(des.InheritFromParent) {
		cDes.InheritFromParent = initial.InheritFromParent
	} else {
		cDes.InheritFromParent = des.InheritFromParent
	}
	if dcl.BoolCanonicalize(des.Reset, initial.Reset) || dcl.IsZeroValue(des.Reset) {
		cDes.Reset = initial.Reset
	} else {
		cDes.Reset = des.Reset
	}

	return cDes
}

func canonicalizePolicySpecSlice(des, initial []PolicySpec, opts ...dcl.ApplyOption) []PolicySpec {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicySpec, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicySpec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicySpec, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicySpec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicySpec(c *Client, des, nw *PolicySpec) *PolicySpec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicySpec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Etag, nw.Etag) {
		nw.Etag = des.Etag
	}
	nw.Rules = canonicalizeNewPolicySpecRulesSet(c, des.Rules, nw.Rules)
	if dcl.BoolCanonicalize(des.InheritFromParent, nw.InheritFromParent) {
		nw.InheritFromParent = des.InheritFromParent
	}
	if dcl.BoolCanonicalize(des.Reset, nw.Reset) {
		nw.Reset = des.Reset
	}

	return nw
}

func canonicalizeNewPolicySpecSet(c *Client, des, nw []PolicySpec) []PolicySpec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicySpec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicySpecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicySpec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicySpecSlice(c *Client, des, nw []PolicySpec) []PolicySpec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicySpec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicySpec(c, &d, &n))
	}

	return items
}

func canonicalizePolicySpecRules(des, initial *PolicySpecRules, opts ...dcl.ApplyOption) *PolicySpecRules {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Values != nil || (initial != nil && initial.Values != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.AllowAll, des.DenyAll, des.Enforce) {
			des.Values = nil
			if initial != nil {
				initial.Values = nil
			}
		}
	}

	if des.AllowAll != nil || (initial != nil && initial.AllowAll != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.DenyAll, des.Enforce) {
			des.AllowAll = nil
			if initial != nil {
				initial.AllowAll = nil
			}
		}
	}

	if des.DenyAll != nil || (initial != nil && initial.DenyAll != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.AllowAll, des.Enforce) {
			des.DenyAll = nil
			if initial != nil {
				initial.DenyAll = nil
			}
		}
	}

	if des.Enforce != nil || (initial != nil && initial.Enforce != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.AllowAll, des.DenyAll) {
			des.Enforce = nil
			if initial != nil {
				initial.Enforce = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PolicySpecRules{}

	cDes.Values = canonicalizePolicySpecRulesValues(des.Values, initial.Values, opts...)
	if dcl.BoolCanonicalize(des.AllowAll, initial.AllowAll) || dcl.IsZeroValue(des.AllowAll) {
		cDes.AllowAll = initial.AllowAll
	} else {
		cDes.AllowAll = des.AllowAll
	}
	if dcl.BoolCanonicalize(des.DenyAll, initial.DenyAll) || dcl.IsZeroValue(des.DenyAll) {
		cDes.DenyAll = initial.DenyAll
	} else {
		cDes.DenyAll = des.DenyAll
	}
	if dcl.BoolCanonicalize(des.Enforce, initial.Enforce) || dcl.IsZeroValue(des.Enforce) {
		cDes.Enforce = initial.Enforce
	} else {
		cDes.Enforce = des.Enforce
	}
	cDes.Condition = canonicalizePolicySpecRulesCondition(des.Condition, initial.Condition, opts...)

	return cDes
}

func canonicalizePolicySpecRulesSlice(des, initial []PolicySpecRules, opts ...dcl.ApplyOption) []PolicySpecRules {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicySpecRules, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicySpecRules(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicySpecRules, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicySpecRules(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicySpecRules(c *Client, des, nw *PolicySpecRules) *PolicySpecRules {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicySpecRules while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Values = canonicalizeNewPolicySpecRulesValues(c, des.Values, nw.Values)
	if dcl.BoolCanonicalize(des.AllowAll, nw.AllowAll) {
		nw.AllowAll = des.AllowAll
	}
	if dcl.BoolCanonicalize(des.DenyAll, nw.DenyAll) {
		nw.DenyAll = des.DenyAll
	}
	if dcl.BoolCanonicalize(des.Enforce, nw.Enforce) {
		nw.Enforce = des.Enforce
	}
	nw.Condition = canonicalizeNewPolicySpecRulesCondition(c, des.Condition, nw.Condition)

	return nw
}

func canonicalizeNewPolicySpecRulesSet(c *Client, des, nw []PolicySpecRules) []PolicySpecRules {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicySpecRules
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicySpecRulesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicySpecRules(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicySpecRulesSlice(c *Client, des, nw []PolicySpecRules) []PolicySpecRules {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicySpecRules
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicySpecRules(c, &d, &n))
	}

	return items
}

func canonicalizePolicySpecRulesValues(des, initial *PolicySpecRulesValues, opts ...dcl.ApplyOption) *PolicySpecRulesValues {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicySpecRulesValues{}

	if dcl.StringArrayCanonicalize(des.AllowedValues, initial.AllowedValues) {
		cDes.AllowedValues = initial.AllowedValues
	} else {
		cDes.AllowedValues = des.AllowedValues
	}
	if dcl.StringArrayCanonicalize(des.DeniedValues, initial.DeniedValues) {
		cDes.DeniedValues = initial.DeniedValues
	} else {
		cDes.DeniedValues = des.DeniedValues
	}

	return cDes
}

func canonicalizePolicySpecRulesValuesSlice(des, initial []PolicySpecRulesValues, opts ...dcl.ApplyOption) []PolicySpecRulesValues {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicySpecRulesValues, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicySpecRulesValues(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicySpecRulesValues, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicySpecRulesValues(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicySpecRulesValues(c *Client, des, nw *PolicySpecRulesValues) *PolicySpecRulesValues {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicySpecRulesValues while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.AllowedValues, nw.AllowedValues) {
		nw.AllowedValues = des.AllowedValues
	}
	if dcl.StringArrayCanonicalize(des.DeniedValues, nw.DeniedValues) {
		nw.DeniedValues = des.DeniedValues
	}

	return nw
}

func canonicalizeNewPolicySpecRulesValuesSet(c *Client, des, nw []PolicySpecRulesValues) []PolicySpecRulesValues {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicySpecRulesValues
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicySpecRulesValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicySpecRulesValues(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicySpecRulesValuesSlice(c *Client, des, nw []PolicySpecRulesValues) []PolicySpecRulesValues {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicySpecRulesValues
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicySpecRulesValues(c, &d, &n))
	}

	return items
}

func canonicalizePolicySpecRulesCondition(des, initial *PolicySpecRulesCondition, opts ...dcl.ApplyOption) *PolicySpecRulesCondition {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicySpecRulesCondition{}

	if canonicalizePolicyRulesConditionExpression(des.Expression, initial.Expression) || dcl.IsZeroValue(des.Expression) {
		cDes.Expression = initial.Expression
	} else {
		cDes.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, initial.Title) || dcl.IsZeroValue(des.Title) {
		cDes.Title = initial.Title
	} else {
		cDes.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, initial.Location) || dcl.IsZeroValue(des.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}

	return cDes
}

func canonicalizePolicySpecRulesConditionSlice(des, initial []PolicySpecRulesCondition, opts ...dcl.ApplyOption) []PolicySpecRulesCondition {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicySpecRulesCondition, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicySpecRulesCondition(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicySpecRulesCondition, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicySpecRulesCondition(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicySpecRulesCondition(c *Client, des, nw *PolicySpecRulesCondition) *PolicySpecRulesCondition {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicySpecRulesCondition while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if canonicalizePolicyRulesConditionExpression(des.Expression, nw.Expression) {
		nw.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, nw.Title) {
		nw.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}

	return nw
}

func canonicalizeNewPolicySpecRulesConditionSet(c *Client, des, nw []PolicySpecRulesCondition) []PolicySpecRulesCondition {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicySpecRulesCondition
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicySpecRulesConditionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicySpecRulesCondition(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicySpecRulesConditionSlice(c *Client, des, nw []PolicySpecRulesCondition) []PolicySpecRulesCondition {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicySpecRulesCondition
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicySpecRulesCondition(c, &d, &n))
	}

	return items
}

func canonicalizePolicyDryRunSpec(des, initial *PolicyDryRunSpec, opts ...dcl.ApplyOption) *PolicyDryRunSpec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicyDryRunSpec{}

	cDes.Rules = canonicalizePolicyDryRunSpecRulesSlice(des.Rules, initial.Rules, opts...)
	if dcl.BoolCanonicalize(des.InheritFromParent, initial.InheritFromParent) || dcl.IsZeroValue(des.InheritFromParent) {
		cDes.InheritFromParent = initial.InheritFromParent
	} else {
		cDes.InheritFromParent = des.InheritFromParent
	}
	if dcl.BoolCanonicalize(des.Reset, initial.Reset) || dcl.IsZeroValue(des.Reset) {
		cDes.Reset = initial.Reset
	} else {
		cDes.Reset = des.Reset
	}

	return cDes
}

func canonicalizePolicyDryRunSpecSlice(des, initial []PolicyDryRunSpec, opts ...dcl.ApplyOption) []PolicyDryRunSpec {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicyDryRunSpec, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicyDryRunSpec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicyDryRunSpec, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicyDryRunSpec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicyDryRunSpec(c *Client, des, nw *PolicyDryRunSpec) *PolicyDryRunSpec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicyDryRunSpec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Etag, nw.Etag) {
		nw.Etag = des.Etag
	}
	nw.Rules = canonicalizeNewPolicyDryRunSpecRulesSlice(c, des.Rules, nw.Rules)
	if dcl.BoolCanonicalize(des.InheritFromParent, nw.InheritFromParent) {
		nw.InheritFromParent = des.InheritFromParent
	}
	if dcl.BoolCanonicalize(des.Reset, nw.Reset) {
		nw.Reset = des.Reset
	}

	return nw
}

func canonicalizeNewPolicyDryRunSpecSet(c *Client, des, nw []PolicyDryRunSpec) []PolicyDryRunSpec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicyDryRunSpec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicyDryRunSpecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicyDryRunSpec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicyDryRunSpecSlice(c *Client, des, nw []PolicyDryRunSpec) []PolicyDryRunSpec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicyDryRunSpec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicyDryRunSpec(c, &d, &n))
	}

	return items
}

func canonicalizePolicyDryRunSpecRules(des, initial *PolicyDryRunSpecRules, opts ...dcl.ApplyOption) *PolicyDryRunSpecRules {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Values != nil || (initial != nil && initial.Values != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.AllowAll, des.DenyAll, des.Enforce) {
			des.Values = nil
			if initial != nil {
				initial.Values = nil
			}
		}
	}

	if des.AllowAll != nil || (initial != nil && initial.AllowAll != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.DenyAll, des.Enforce) {
			des.AllowAll = nil
			if initial != nil {
				initial.AllowAll = nil
			}
		}
	}

	if des.DenyAll != nil || (initial != nil && initial.DenyAll != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.AllowAll, des.Enforce) {
			des.DenyAll = nil
			if initial != nil {
				initial.DenyAll = nil
			}
		}
	}

	if des.Enforce != nil || (initial != nil && initial.Enforce != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Values, des.AllowAll, des.DenyAll) {
			des.Enforce = nil
			if initial != nil {
				initial.Enforce = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PolicyDryRunSpecRules{}

	cDes.Values = canonicalizePolicyDryRunSpecRulesValues(des.Values, initial.Values, opts...)
	if dcl.BoolCanonicalize(des.AllowAll, initial.AllowAll) || dcl.IsZeroValue(des.AllowAll) {
		cDes.AllowAll = initial.AllowAll
	} else {
		cDes.AllowAll = des.AllowAll
	}
	if dcl.BoolCanonicalize(des.DenyAll, initial.DenyAll) || dcl.IsZeroValue(des.DenyAll) {
		cDes.DenyAll = initial.DenyAll
	} else {
		cDes.DenyAll = des.DenyAll
	}
	if dcl.BoolCanonicalize(des.Enforce, initial.Enforce) || dcl.IsZeroValue(des.Enforce) {
		cDes.Enforce = initial.Enforce
	} else {
		cDes.Enforce = des.Enforce
	}
	cDes.Condition = canonicalizePolicyDryRunSpecRulesCondition(des.Condition, initial.Condition, opts...)

	return cDes
}

func canonicalizePolicyDryRunSpecRulesSlice(des, initial []PolicyDryRunSpecRules, opts ...dcl.ApplyOption) []PolicyDryRunSpecRules {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicyDryRunSpecRules, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicyDryRunSpecRules(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicyDryRunSpecRules, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicyDryRunSpecRules(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicyDryRunSpecRules(c *Client, des, nw *PolicyDryRunSpecRules) *PolicyDryRunSpecRules {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicyDryRunSpecRules while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Values = canonicalizeNewPolicyDryRunSpecRulesValues(c, des.Values, nw.Values)
	if dcl.BoolCanonicalize(des.AllowAll, nw.AllowAll) {
		nw.AllowAll = des.AllowAll
	}
	if dcl.BoolCanonicalize(des.DenyAll, nw.DenyAll) {
		nw.DenyAll = des.DenyAll
	}
	if dcl.BoolCanonicalize(des.Enforce, nw.Enforce) {
		nw.Enforce = des.Enforce
	}
	nw.Condition = canonicalizeNewPolicyDryRunSpecRulesCondition(c, des.Condition, nw.Condition)

	return nw
}

func canonicalizeNewPolicyDryRunSpecRulesSet(c *Client, des, nw []PolicyDryRunSpecRules) []PolicyDryRunSpecRules {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicyDryRunSpecRules
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicyDryRunSpecRulesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicyDryRunSpecRules(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicyDryRunSpecRulesSlice(c *Client, des, nw []PolicyDryRunSpecRules) []PolicyDryRunSpecRules {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicyDryRunSpecRules
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicyDryRunSpecRules(c, &d, &n))
	}

	return items
}

func canonicalizePolicyDryRunSpecRulesValues(des, initial *PolicyDryRunSpecRulesValues, opts ...dcl.ApplyOption) *PolicyDryRunSpecRulesValues {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicyDryRunSpecRulesValues{}

	if dcl.StringArrayCanonicalize(des.AllowedValues, initial.AllowedValues) {
		cDes.AllowedValues = initial.AllowedValues
	} else {
		cDes.AllowedValues = des.AllowedValues
	}
	if dcl.StringArrayCanonicalize(des.DeniedValues, initial.DeniedValues) {
		cDes.DeniedValues = initial.DeniedValues
	} else {
		cDes.DeniedValues = des.DeniedValues
	}

	return cDes
}

func canonicalizePolicyDryRunSpecRulesValuesSlice(des, initial []PolicyDryRunSpecRulesValues, opts ...dcl.ApplyOption) []PolicyDryRunSpecRulesValues {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicyDryRunSpecRulesValues, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicyDryRunSpecRulesValues(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicyDryRunSpecRulesValues, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicyDryRunSpecRulesValues(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicyDryRunSpecRulesValues(c *Client, des, nw *PolicyDryRunSpecRulesValues) *PolicyDryRunSpecRulesValues {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicyDryRunSpecRulesValues while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.AllowedValues, nw.AllowedValues) {
		nw.AllowedValues = des.AllowedValues
	}
	if dcl.StringArrayCanonicalize(des.DeniedValues, nw.DeniedValues) {
		nw.DeniedValues = des.DeniedValues
	}

	return nw
}

func canonicalizeNewPolicyDryRunSpecRulesValuesSet(c *Client, des, nw []PolicyDryRunSpecRulesValues) []PolicyDryRunSpecRulesValues {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicyDryRunSpecRulesValues
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicyDryRunSpecRulesValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicyDryRunSpecRulesValues(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicyDryRunSpecRulesValuesSlice(c *Client, des, nw []PolicyDryRunSpecRulesValues) []PolicyDryRunSpecRulesValues {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicyDryRunSpecRulesValues
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicyDryRunSpecRulesValues(c, &d, &n))
	}

	return items
}

func canonicalizePolicyDryRunSpecRulesCondition(des, initial *PolicyDryRunSpecRulesCondition, opts ...dcl.ApplyOption) *PolicyDryRunSpecRulesCondition {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PolicyDryRunSpecRulesCondition{}

	if dcl.StringCanonicalize(des.Expression, initial.Expression) || dcl.IsZeroValue(des.Expression) {
		cDes.Expression = initial.Expression
	} else {
		cDes.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, initial.Title) || dcl.IsZeroValue(des.Title) {
		cDes.Title = initial.Title
	} else {
		cDes.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, initial.Location) || dcl.IsZeroValue(des.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}

	return cDes
}

func canonicalizePolicyDryRunSpecRulesConditionSlice(des, initial []PolicyDryRunSpecRulesCondition, opts ...dcl.ApplyOption) []PolicyDryRunSpecRulesCondition {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PolicyDryRunSpecRulesCondition, 0, len(des))
		for _, d := range des {
			cd := canonicalizePolicyDryRunSpecRulesCondition(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PolicyDryRunSpecRulesCondition, 0, len(des))
	for i, d := range des {
		cd := canonicalizePolicyDryRunSpecRulesCondition(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPolicyDryRunSpecRulesCondition(c *Client, des, nw *PolicyDryRunSpecRulesCondition) *PolicyDryRunSpecRulesCondition {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PolicyDryRunSpecRulesCondition while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Expression, nw.Expression) {
		nw.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, nw.Title) {
		nw.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}

	return nw
}

func canonicalizeNewPolicyDryRunSpecRulesConditionSet(c *Client, des, nw []PolicyDryRunSpecRulesCondition) []PolicyDryRunSpecRulesCondition {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PolicyDryRunSpecRulesCondition
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePolicyDryRunSpecRulesConditionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPolicyDryRunSpecRulesCondition(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPolicyDryRunSpecRulesConditionSlice(c *Client, des, nw []PolicyDryRunSpecRulesCondition) []PolicyDryRunSpecRulesCondition {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PolicyDryRunSpecRulesCondition
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPolicyDryRunSpecRulesCondition(c, &d, &n))
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
func diffPolicy(c *Client, desired, actual *Policy, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{CustomDiff: canonicalizePolicyName, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Spec, actual.Spec, dcl.DiffInfo{ObjectFunction: comparePolicySpecNewStyle, EmptyObject: EmptyPolicySpec, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Spec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DryRunSpec, actual.DryRunSpec, dcl.DiffInfo{ObjectFunction: comparePolicyDryRunSpecNewStyle, EmptyObject: EmptyPolicyDryRunSpec, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("DryRunSpec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Parent, actual.Parent, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Parent")); len(ds) != 0 || err != nil {
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
func comparePolicySpecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicySpec)
	if !ok {
		desiredNotPointer, ok := d.(PolicySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpec or *PolicySpec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicySpec)
	if !ok {
		actualNotPointer, ok := a.(PolicySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rules, actual.Rules, dcl.DiffInfo{Type: "Set", ObjectFunction: comparePolicySpecRulesNewStyle, EmptyObject: EmptyPolicySpecRules, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Rules")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InheritFromParent, actual.InheritFromParent, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("InheritFromParent")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reset, actual.Reset, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Reset")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicySpecRulesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicySpecRules)
	if !ok {
		desiredNotPointer, ok := d.(PolicySpecRules)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRules or *PolicySpecRules", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicySpecRules)
	if !ok {
		actualNotPointer, ok := a.(PolicySpecRules)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRules", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{ObjectFunction: comparePolicySpecRulesValuesNewStyle, EmptyObject: EmptyPolicySpecRulesValues, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowAll, actual.AllowAll, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("AllowAll")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DenyAll, actual.DenyAll, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("DenyAll")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Enforce, actual.Enforce, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Enforce")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Condition, actual.Condition, dcl.DiffInfo{ObjectFunction: comparePolicySpecRulesConditionNewStyle, EmptyObject: EmptyPolicySpecRulesCondition, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Condition")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicySpecRulesValuesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicySpecRulesValues)
	if !ok {
		desiredNotPointer, ok := d.(PolicySpecRulesValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRulesValues or *PolicySpecRulesValues", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicySpecRulesValues)
	if !ok {
		actualNotPointer, ok := a.(PolicySpecRulesValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRulesValues", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedValues, actual.AllowedValues, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("AllowedValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DeniedValues, actual.DeniedValues, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("DeniedValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicySpecRulesConditionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicySpecRulesCondition)
	if !ok {
		desiredNotPointer, ok := d.(PolicySpecRulesCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRulesCondition or *PolicySpecRulesCondition", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicySpecRulesCondition)
	if !ok {
		actualNotPointer, ok := a.(PolicySpecRulesCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicySpecRulesCondition", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Expression, actual.Expression, dcl.DiffInfo{CustomDiff: canonicalizePolicyRulesConditionExpression, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Expression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Title, actual.Title, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Title")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicyDryRunSpecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicyDryRunSpec)
	if !ok {
		desiredNotPointer, ok := d.(PolicyDryRunSpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpec or *PolicyDryRunSpec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicyDryRunSpec)
	if !ok {
		actualNotPointer, ok := a.(PolicyDryRunSpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rules, actual.Rules, dcl.DiffInfo{ObjectFunction: comparePolicyDryRunSpecRulesNewStyle, EmptyObject: EmptyPolicyDryRunSpecRules, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Rules")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InheritFromParent, actual.InheritFromParent, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("InheritFromParent")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reset, actual.Reset, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Reset")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicyDryRunSpecRulesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicyDryRunSpecRules)
	if !ok {
		desiredNotPointer, ok := d.(PolicyDryRunSpecRules)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRules or *PolicyDryRunSpecRules", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicyDryRunSpecRules)
	if !ok {
		actualNotPointer, ok := a.(PolicyDryRunSpecRules)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRules", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{ObjectFunction: comparePolicyDryRunSpecRulesValuesNewStyle, EmptyObject: EmptyPolicyDryRunSpecRulesValues, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowAll, actual.AllowAll, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("AllowAll")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DenyAll, actual.DenyAll, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("DenyAll")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Enforce, actual.Enforce, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Enforce")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Condition, actual.Condition, dcl.DiffInfo{ObjectFunction: comparePolicyDryRunSpecRulesConditionNewStyle, EmptyObject: EmptyPolicyDryRunSpecRulesCondition, OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Condition")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicyDryRunSpecRulesValuesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicyDryRunSpecRulesValues)
	if !ok {
		desiredNotPointer, ok := d.(PolicyDryRunSpecRulesValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRulesValues or *PolicyDryRunSpecRulesValues", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicyDryRunSpecRulesValues)
	if !ok {
		actualNotPointer, ok := a.(PolicyDryRunSpecRulesValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRulesValues", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedValues, actual.AllowedValues, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("AllowedValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DeniedValues, actual.DeniedValues, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("DeniedValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePolicyDryRunSpecRulesConditionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PolicyDryRunSpecRulesCondition)
	if !ok {
		desiredNotPointer, ok := d.(PolicyDryRunSpecRulesCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRulesCondition or *PolicyDryRunSpecRulesCondition", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PolicyDryRunSpecRulesCondition)
	if !ok {
		actualNotPointer, ok := a.(PolicyDryRunSpecRulesCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PolicyDryRunSpecRulesCondition", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Expression, actual.Expression, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Expression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Title, actual.Title, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Title")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePolicyUpdatePolicyOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
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
func (r *Policy) urlNormalized() *Policy {
	normalized := dcl.Copy(*r).(Policy)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Parent = r.Parent
	return &normalized
}

func (r *Policy) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdatePolicy" {
		fields := map[string]interface{}{
			"parent": dcl.ValueOrEmptyString(nr.Parent),
			"name":   dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("{{parent}}/policies/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Policy resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Policy) marshal(c *Client) ([]byte, error) {
	m, err := expandPolicy(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Policy: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalPolicy decodes JSON responses into the Policy resource schema.
func unmarshalPolicy(b []byte, c *Client, res *Policy) (*Policy, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapPolicy(m, c, res)
}

func unmarshalMapPolicy(m map[string]interface{}, c *Client, res *Policy) (*Policy, error) {

	flattened := flattenPolicy(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandPolicy expands Policy into a JSON request object.
func expandPolicy(c *Client, f *Policy) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := expandPolicyName(c, f.Name, res); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := expandPolicySpec(c, f.Spec, res); err != nil {
		return nil, fmt.Errorf("error expanding Spec into spec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["spec"] = v
	}
	if v, err := expandPolicyDryRunSpec(c, f.DryRunSpec, res); err != nil {
		return nil, fmt.Errorf("error expanding DryRunSpec into dryRunSpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["dryRunSpec"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Parent into parent: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["parent"] = v
	}

	return m, nil
}

// flattenPolicy flattens Policy from a JSON request object into the
// Policy type.
func flattenPolicy(c *Client, i interface{}, res *Policy) *Policy {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Policy{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Spec = flattenPolicySpec(c, m["spec"], res)
	resultRes.DryRunSpec = flattenPolicyDryRunSpec(c, m["dryRunSpec"], res)
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Parent = dcl.FlattenString(m["parent"])

	return resultRes
}

// expandPolicySpecMap expands the contents of PolicySpec into a JSON
// request object.
func expandPolicySpecMap(c *Client, f map[string]PolicySpec, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicySpec(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicySpecSlice expands the contents of PolicySpec into a JSON
// request object.
func expandPolicySpecSlice(c *Client, f []PolicySpec, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicySpec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicySpecMap flattens the contents of PolicySpec from a JSON
// response object.
func flattenPolicySpecMap(c *Client, i interface{}, res *Policy) map[string]PolicySpec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicySpec{}
	}

	if len(a) == 0 {
		return map[string]PolicySpec{}
	}

	items := make(map[string]PolicySpec)
	for k, item := range a {
		items[k] = *flattenPolicySpec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicySpecSlice flattens the contents of PolicySpec from a JSON
// response object.
func flattenPolicySpecSlice(c *Client, i interface{}, res *Policy) []PolicySpec {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicySpec{}
	}

	if len(a) == 0 {
		return []PolicySpec{}
	}

	items := make([]PolicySpec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicySpec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicySpec expands an instance of PolicySpec into a JSON
// request object.
func expandPolicySpec(c *Client, f *PolicySpec, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPolicySpecRulesSlice(c, f.Rules, res); err != nil {
		return nil, fmt.Errorf("error expanding Rules into rules: %w", err)
	} else if v != nil {
		m["rules"] = v
	}
	if v := f.InheritFromParent; !dcl.IsEmptyValueIndirect(v) {
		m["inheritFromParent"] = v
	}
	if v := f.Reset; !dcl.IsEmptyValueIndirect(v) {
		m["reset"] = v
	}

	return m, nil
}

// flattenPolicySpec flattens an instance of PolicySpec from a JSON
// response object.
func flattenPolicySpec(c *Client, i interface{}, res *Policy) *PolicySpec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicySpec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicySpec
	}
	r.Etag = dcl.FlattenString(m["etag"])
	r.UpdateTime = dcl.FlattenString(m["updateTime"])
	r.Rules = flattenPolicySpecRulesSlice(c, m["rules"], res)
	r.InheritFromParent = dcl.FlattenBool(m["inheritFromParent"])
	r.Reset = dcl.FlattenBool(m["reset"])

	return r
}

// expandPolicySpecRulesMap expands the contents of PolicySpecRules into a JSON
// request object.
func expandPolicySpecRulesMap(c *Client, f map[string]PolicySpecRules, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicySpecRules(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicySpecRulesSlice expands the contents of PolicySpecRules into a JSON
// request object.
func expandPolicySpecRulesSlice(c *Client, f []PolicySpecRules, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicySpecRules(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicySpecRulesMap flattens the contents of PolicySpecRules from a JSON
// response object.
func flattenPolicySpecRulesMap(c *Client, i interface{}, res *Policy) map[string]PolicySpecRules {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicySpecRules{}
	}

	if len(a) == 0 {
		return map[string]PolicySpecRules{}
	}

	items := make(map[string]PolicySpecRules)
	for k, item := range a {
		items[k] = *flattenPolicySpecRules(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicySpecRulesSlice flattens the contents of PolicySpecRules from a JSON
// response object.
func flattenPolicySpecRulesSlice(c *Client, i interface{}, res *Policy) []PolicySpecRules {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicySpecRules{}
	}

	if len(a) == 0 {
		return []PolicySpecRules{}
	}

	items := make([]PolicySpecRules, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicySpecRules(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicySpecRules expands an instance of PolicySpecRules into a JSON
// request object.
func expandPolicySpecRules(c *Client, f *PolicySpecRules, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPolicySpecRulesValues(c, f.Values, res); err != nil {
		return nil, fmt.Errorf("error expanding Values into values: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["values"] = v
	}
	if v := f.AllowAll; !dcl.IsEmptyValueIndirect(v) {
		m["allowAll"] = v
	}
	if v := f.DenyAll; !dcl.IsEmptyValueIndirect(v) {
		m["denyAll"] = v
	}
	if v := f.Enforce; !dcl.IsEmptyValueIndirect(v) {
		m["enforce"] = v
	}
	if v, err := expandPolicySpecRulesCondition(c, f.Condition, res); err != nil {
		return nil, fmt.Errorf("error expanding Condition into condition: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["condition"] = v
	}

	return m, nil
}

// flattenPolicySpecRules flattens an instance of PolicySpecRules from a JSON
// response object.
func flattenPolicySpecRules(c *Client, i interface{}, res *Policy) *PolicySpecRules {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicySpecRules{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicySpecRules
	}
	r.Values = flattenPolicySpecRulesValues(c, m["values"], res)
	r.AllowAll = dcl.FlattenBool(m["allowAll"])
	r.DenyAll = dcl.FlattenBool(m["denyAll"])
	r.Enforce = dcl.FlattenBool(m["enforce"])
	r.Condition = flattenPolicySpecRulesCondition(c, m["condition"], res)

	return r
}

// expandPolicySpecRulesValuesMap expands the contents of PolicySpecRulesValues into a JSON
// request object.
func expandPolicySpecRulesValuesMap(c *Client, f map[string]PolicySpecRulesValues, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicySpecRulesValues(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicySpecRulesValuesSlice expands the contents of PolicySpecRulesValues into a JSON
// request object.
func expandPolicySpecRulesValuesSlice(c *Client, f []PolicySpecRulesValues, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicySpecRulesValues(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicySpecRulesValuesMap flattens the contents of PolicySpecRulesValues from a JSON
// response object.
func flattenPolicySpecRulesValuesMap(c *Client, i interface{}, res *Policy) map[string]PolicySpecRulesValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicySpecRulesValues{}
	}

	if len(a) == 0 {
		return map[string]PolicySpecRulesValues{}
	}

	items := make(map[string]PolicySpecRulesValues)
	for k, item := range a {
		items[k] = *flattenPolicySpecRulesValues(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicySpecRulesValuesSlice flattens the contents of PolicySpecRulesValues from a JSON
// response object.
func flattenPolicySpecRulesValuesSlice(c *Client, i interface{}, res *Policy) []PolicySpecRulesValues {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicySpecRulesValues{}
	}

	if len(a) == 0 {
		return []PolicySpecRulesValues{}
	}

	items := make([]PolicySpecRulesValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicySpecRulesValues(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicySpecRulesValues expands an instance of PolicySpecRulesValues into a JSON
// request object.
func expandPolicySpecRulesValues(c *Client, f *PolicySpecRulesValues, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowedValues; v != nil {
		m["allowedValues"] = v
	}
	if v := f.DeniedValues; v != nil {
		m["deniedValues"] = v
	}

	return m, nil
}

// flattenPolicySpecRulesValues flattens an instance of PolicySpecRulesValues from a JSON
// response object.
func flattenPolicySpecRulesValues(c *Client, i interface{}, res *Policy) *PolicySpecRulesValues {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicySpecRulesValues{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicySpecRulesValues
	}
	r.AllowedValues = dcl.FlattenStringSlice(m["allowedValues"])
	r.DeniedValues = dcl.FlattenStringSlice(m["deniedValues"])

	return r
}

// expandPolicySpecRulesConditionMap expands the contents of PolicySpecRulesCondition into a JSON
// request object.
func expandPolicySpecRulesConditionMap(c *Client, f map[string]PolicySpecRulesCondition, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicySpecRulesCondition(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicySpecRulesConditionSlice expands the contents of PolicySpecRulesCondition into a JSON
// request object.
func expandPolicySpecRulesConditionSlice(c *Client, f []PolicySpecRulesCondition, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicySpecRulesCondition(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicySpecRulesConditionMap flattens the contents of PolicySpecRulesCondition from a JSON
// response object.
func flattenPolicySpecRulesConditionMap(c *Client, i interface{}, res *Policy) map[string]PolicySpecRulesCondition {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicySpecRulesCondition{}
	}

	if len(a) == 0 {
		return map[string]PolicySpecRulesCondition{}
	}

	items := make(map[string]PolicySpecRulesCondition)
	for k, item := range a {
		items[k] = *flattenPolicySpecRulesCondition(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicySpecRulesConditionSlice flattens the contents of PolicySpecRulesCondition from a JSON
// response object.
func flattenPolicySpecRulesConditionSlice(c *Client, i interface{}, res *Policy) []PolicySpecRulesCondition {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicySpecRulesCondition{}
	}

	if len(a) == 0 {
		return []PolicySpecRulesCondition{}
	}

	items := make([]PolicySpecRulesCondition, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicySpecRulesCondition(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicySpecRulesCondition expands an instance of PolicySpecRulesCondition into a JSON
// request object.
func expandPolicySpecRulesCondition(c *Client, f *PolicySpecRulesCondition, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Expression; !dcl.IsEmptyValueIndirect(v) {
		m["expression"] = v
	}
	if v := f.Title; !dcl.IsEmptyValueIndirect(v) {
		m["title"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenPolicySpecRulesCondition flattens an instance of PolicySpecRulesCondition from a JSON
// response object.
func flattenPolicySpecRulesCondition(c *Client, i interface{}, res *Policy) *PolicySpecRulesCondition {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicySpecRulesCondition{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicySpecRulesCondition
	}
	r.Expression = dcl.FlattenString(m["expression"])
	r.Title = dcl.FlattenString(m["title"])
	r.Description = dcl.FlattenString(m["description"])
	r.Location = dcl.FlattenString(m["location"])

	return r
}

// expandPolicyDryRunSpecMap expands the contents of PolicyDryRunSpec into a JSON
// request object.
func expandPolicyDryRunSpecMap(c *Client, f map[string]PolicyDryRunSpec, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicyDryRunSpec(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicyDryRunSpecSlice expands the contents of PolicyDryRunSpec into a JSON
// request object.
func expandPolicyDryRunSpecSlice(c *Client, f []PolicyDryRunSpec, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicyDryRunSpec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicyDryRunSpecMap flattens the contents of PolicyDryRunSpec from a JSON
// response object.
func flattenPolicyDryRunSpecMap(c *Client, i interface{}, res *Policy) map[string]PolicyDryRunSpec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicyDryRunSpec{}
	}

	if len(a) == 0 {
		return map[string]PolicyDryRunSpec{}
	}

	items := make(map[string]PolicyDryRunSpec)
	for k, item := range a {
		items[k] = *flattenPolicyDryRunSpec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicyDryRunSpecSlice flattens the contents of PolicyDryRunSpec from a JSON
// response object.
func flattenPolicyDryRunSpecSlice(c *Client, i interface{}, res *Policy) []PolicyDryRunSpec {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicyDryRunSpec{}
	}

	if len(a) == 0 {
		return []PolicyDryRunSpec{}
	}

	items := make([]PolicyDryRunSpec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicyDryRunSpec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicyDryRunSpec expands an instance of PolicyDryRunSpec into a JSON
// request object.
func expandPolicyDryRunSpec(c *Client, f *PolicyDryRunSpec, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPolicyDryRunSpecRulesSlice(c, f.Rules, res); err != nil {
		return nil, fmt.Errorf("error expanding Rules into rules: %w", err)
	} else if v != nil {
		m["rules"] = v
	}
	if v := f.InheritFromParent; !dcl.IsEmptyValueIndirect(v) {
		m["inheritFromParent"] = v
	}
	if v := f.Reset; !dcl.IsEmptyValueIndirect(v) {
		m["reset"] = v
	}

	return m, nil
}

// flattenPolicyDryRunSpec flattens an instance of PolicyDryRunSpec from a JSON
// response object.
func flattenPolicyDryRunSpec(c *Client, i interface{}, res *Policy) *PolicyDryRunSpec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicyDryRunSpec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicyDryRunSpec
	}
	r.Etag = dcl.FlattenString(m["etag"])
	r.UpdateTime = dcl.FlattenString(m["updateTime"])
	r.Rules = flattenPolicyDryRunSpecRulesSlice(c, m["rules"], res)
	r.InheritFromParent = dcl.FlattenBool(m["inheritFromParent"])
	r.Reset = dcl.FlattenBool(m["reset"])

	return r
}

// expandPolicyDryRunSpecRulesMap expands the contents of PolicyDryRunSpecRules into a JSON
// request object.
func expandPolicyDryRunSpecRulesMap(c *Client, f map[string]PolicyDryRunSpecRules, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicyDryRunSpecRules(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicyDryRunSpecRulesSlice expands the contents of PolicyDryRunSpecRules into a JSON
// request object.
func expandPolicyDryRunSpecRulesSlice(c *Client, f []PolicyDryRunSpecRules, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicyDryRunSpecRules(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicyDryRunSpecRulesMap flattens the contents of PolicyDryRunSpecRules from a JSON
// response object.
func flattenPolicyDryRunSpecRulesMap(c *Client, i interface{}, res *Policy) map[string]PolicyDryRunSpecRules {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicyDryRunSpecRules{}
	}

	if len(a) == 0 {
		return map[string]PolicyDryRunSpecRules{}
	}

	items := make(map[string]PolicyDryRunSpecRules)
	for k, item := range a {
		items[k] = *flattenPolicyDryRunSpecRules(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicyDryRunSpecRulesSlice flattens the contents of PolicyDryRunSpecRules from a JSON
// response object.
func flattenPolicyDryRunSpecRulesSlice(c *Client, i interface{}, res *Policy) []PolicyDryRunSpecRules {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicyDryRunSpecRules{}
	}

	if len(a) == 0 {
		return []PolicyDryRunSpecRules{}
	}

	items := make([]PolicyDryRunSpecRules, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicyDryRunSpecRules(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicyDryRunSpecRules expands an instance of PolicyDryRunSpecRules into a JSON
// request object.
func expandPolicyDryRunSpecRules(c *Client, f *PolicyDryRunSpecRules, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPolicyDryRunSpecRulesValues(c, f.Values, res); err != nil {
		return nil, fmt.Errorf("error expanding Values into values: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["values"] = v
	}
	if v := f.AllowAll; !dcl.IsEmptyValueIndirect(v) {
		m["allowAll"] = v
	}
	if v := f.DenyAll; !dcl.IsEmptyValueIndirect(v) {
		m["denyAll"] = v
	}
	if v := f.Enforce; !dcl.IsEmptyValueIndirect(v) {
		m["enforce"] = v
	}
	if v, err := expandPolicyDryRunSpecRulesCondition(c, f.Condition, res); err != nil {
		return nil, fmt.Errorf("error expanding Condition into condition: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["condition"] = v
	}

	return m, nil
}

// flattenPolicyDryRunSpecRules flattens an instance of PolicyDryRunSpecRules from a JSON
// response object.
func flattenPolicyDryRunSpecRules(c *Client, i interface{}, res *Policy) *PolicyDryRunSpecRules {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicyDryRunSpecRules{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicyDryRunSpecRules
	}
	r.Values = flattenPolicyDryRunSpecRulesValues(c, m["values"], res)
	r.AllowAll = dcl.FlattenBool(m["allowAll"])
	r.DenyAll = dcl.FlattenBool(m["denyAll"])
	r.Enforce = dcl.FlattenBool(m["enforce"])
	r.Condition = flattenPolicyDryRunSpecRulesCondition(c, m["condition"], res)

	return r
}

// expandPolicyDryRunSpecRulesValuesMap expands the contents of PolicyDryRunSpecRulesValues into a JSON
// request object.
func expandPolicyDryRunSpecRulesValuesMap(c *Client, f map[string]PolicyDryRunSpecRulesValues, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicyDryRunSpecRulesValues(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicyDryRunSpecRulesValuesSlice expands the contents of PolicyDryRunSpecRulesValues into a JSON
// request object.
func expandPolicyDryRunSpecRulesValuesSlice(c *Client, f []PolicyDryRunSpecRulesValues, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicyDryRunSpecRulesValues(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicyDryRunSpecRulesValuesMap flattens the contents of PolicyDryRunSpecRulesValues from a JSON
// response object.
func flattenPolicyDryRunSpecRulesValuesMap(c *Client, i interface{}, res *Policy) map[string]PolicyDryRunSpecRulesValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicyDryRunSpecRulesValues{}
	}

	if len(a) == 0 {
		return map[string]PolicyDryRunSpecRulesValues{}
	}

	items := make(map[string]PolicyDryRunSpecRulesValues)
	for k, item := range a {
		items[k] = *flattenPolicyDryRunSpecRulesValues(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicyDryRunSpecRulesValuesSlice flattens the contents of PolicyDryRunSpecRulesValues from a JSON
// response object.
func flattenPolicyDryRunSpecRulesValuesSlice(c *Client, i interface{}, res *Policy) []PolicyDryRunSpecRulesValues {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicyDryRunSpecRulesValues{}
	}

	if len(a) == 0 {
		return []PolicyDryRunSpecRulesValues{}
	}

	items := make([]PolicyDryRunSpecRulesValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicyDryRunSpecRulesValues(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicyDryRunSpecRulesValues expands an instance of PolicyDryRunSpecRulesValues into a JSON
// request object.
func expandPolicyDryRunSpecRulesValues(c *Client, f *PolicyDryRunSpecRulesValues, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowedValues; v != nil {
		m["allowedValues"] = v
	}
	if v := f.DeniedValues; v != nil {
		m["deniedValues"] = v
	}

	return m, nil
}

// flattenPolicyDryRunSpecRulesValues flattens an instance of PolicyDryRunSpecRulesValues from a JSON
// response object.
func flattenPolicyDryRunSpecRulesValues(c *Client, i interface{}, res *Policy) *PolicyDryRunSpecRulesValues {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicyDryRunSpecRulesValues{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicyDryRunSpecRulesValues
	}
	r.AllowedValues = dcl.FlattenStringSlice(m["allowedValues"])
	r.DeniedValues = dcl.FlattenStringSlice(m["deniedValues"])

	return r
}

// expandPolicyDryRunSpecRulesConditionMap expands the contents of PolicyDryRunSpecRulesCondition into a JSON
// request object.
func expandPolicyDryRunSpecRulesConditionMap(c *Client, f map[string]PolicyDryRunSpecRulesCondition, res *Policy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPolicyDryRunSpecRulesCondition(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPolicyDryRunSpecRulesConditionSlice expands the contents of PolicyDryRunSpecRulesCondition into a JSON
// request object.
func expandPolicyDryRunSpecRulesConditionSlice(c *Client, f []PolicyDryRunSpecRulesCondition, res *Policy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPolicyDryRunSpecRulesCondition(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPolicyDryRunSpecRulesConditionMap flattens the contents of PolicyDryRunSpecRulesCondition from a JSON
// response object.
func flattenPolicyDryRunSpecRulesConditionMap(c *Client, i interface{}, res *Policy) map[string]PolicyDryRunSpecRulesCondition {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PolicyDryRunSpecRulesCondition{}
	}

	if len(a) == 0 {
		return map[string]PolicyDryRunSpecRulesCondition{}
	}

	items := make(map[string]PolicyDryRunSpecRulesCondition)
	for k, item := range a {
		items[k] = *flattenPolicyDryRunSpecRulesCondition(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPolicyDryRunSpecRulesConditionSlice flattens the contents of PolicyDryRunSpecRulesCondition from a JSON
// response object.
func flattenPolicyDryRunSpecRulesConditionSlice(c *Client, i interface{}, res *Policy) []PolicyDryRunSpecRulesCondition {
	a, ok := i.([]interface{})
	if !ok {
		return []PolicyDryRunSpecRulesCondition{}
	}

	if len(a) == 0 {
		return []PolicyDryRunSpecRulesCondition{}
	}

	items := make([]PolicyDryRunSpecRulesCondition, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPolicyDryRunSpecRulesCondition(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPolicyDryRunSpecRulesCondition expands an instance of PolicyDryRunSpecRulesCondition into a JSON
// request object.
func expandPolicyDryRunSpecRulesCondition(c *Client, f *PolicyDryRunSpecRulesCondition, res *Policy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Expression; !dcl.IsEmptyValueIndirect(v) {
		m["expression"] = v
	}
	if v := f.Title; !dcl.IsEmptyValueIndirect(v) {
		m["title"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenPolicyDryRunSpecRulesCondition flattens an instance of PolicyDryRunSpecRulesCondition from a JSON
// response object.
func flattenPolicyDryRunSpecRulesCondition(c *Client, i interface{}, res *Policy) *PolicyDryRunSpecRulesCondition {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PolicyDryRunSpecRulesCondition{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPolicyDryRunSpecRulesCondition
	}
	r.Expression = dcl.FlattenString(m["expression"])
	r.Title = dcl.FlattenString(m["title"])
	r.Description = dcl.FlattenString(m["description"])
	r.Location = dcl.FlattenString(m["location"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Policy) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalPolicy(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Parent == nil && ncr.Parent == nil {
			c.Config.Logger.Info("Both Parent fields null - considering equal.")
		} else if nr.Parent == nil || ncr.Parent == nil {
			c.Config.Logger.Info("Only one Parent field is null - considering unequal.")
			return false
		} else if *nr.Parent != *ncr.Parent {
			return false
		}
		if nr.Name == nil && ncr.Name == nil {
			c.Config.Logger.Info("Both Name fields null - considering equal.")
		} else if nr.Name == nil || ncr.Name == nil {
			c.Config.Logger.Info("Only one Name field is null - considering unequal.")
			return false
		} else if *nr.Name != *ncr.Name {
			return false
		}
		return true
	}
}

type policyDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         policyApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToPolicyDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]policyDiff, error) {
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
	var diffs []policyDiff
	// For each operation name, create a policyDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := policyDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToPolicyApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToPolicyApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (policyApiOperation, error) {
	switch opName {

	case "updatePolicyUpdatePolicyOperation":
		return &updatePolicyUpdatePolicyOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractPolicyFields(r *Policy) error {
	vSpec := r.Spec
	if vSpec == nil {
		// note: explicitly not the empty object.
		vSpec = &PolicySpec{}
	}
	if err := extractPolicySpecFields(r, vSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSpec) {
		r.Spec = vSpec
	}
	vDryRunSpec := r.DryRunSpec
	if vDryRunSpec == nil {
		// note: explicitly not the empty object.
		vDryRunSpec = &PolicyDryRunSpec{}
	}
	if err := extractPolicyDryRunSpecFields(r, vDryRunSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDryRunSpec) {
		r.DryRunSpec = vDryRunSpec
	}
	return nil
}
func extractPolicySpecFields(r *Policy, o *PolicySpec) error {
	return nil
}
func extractPolicySpecRulesFields(r *Policy, o *PolicySpecRules) error {
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &PolicySpecRulesValues{}
	}
	if err := extractPolicySpecRulesValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValues) {
		o.Values = vValues
	}
	vCondition := o.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &PolicySpecRulesCondition{}
	}
	if err := extractPolicySpecRulesConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		o.Condition = vCondition
	}
	return nil
}
func extractPolicySpecRulesValuesFields(r *Policy, o *PolicySpecRulesValues) error {
	return nil
}
func extractPolicySpecRulesConditionFields(r *Policy, o *PolicySpecRulesCondition) error {
	return nil
}
func extractPolicyDryRunSpecFields(r *Policy, o *PolicyDryRunSpec) error {
	return nil
}
func extractPolicyDryRunSpecRulesFields(r *Policy, o *PolicyDryRunSpecRules) error {
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &PolicyDryRunSpecRulesValues{}
	}
	if err := extractPolicyDryRunSpecRulesValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValues) {
		o.Values = vValues
	}
	vCondition := o.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &PolicyDryRunSpecRulesCondition{}
	}
	if err := extractPolicyDryRunSpecRulesConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		o.Condition = vCondition
	}
	return nil
}
func extractPolicyDryRunSpecRulesValuesFields(r *Policy, o *PolicyDryRunSpecRulesValues) error {
	return nil
}
func extractPolicyDryRunSpecRulesConditionFields(r *Policy, o *PolicyDryRunSpecRulesCondition) error {
	return nil
}

func postReadExtractPolicyFields(r *Policy) error {
	vSpec := r.Spec
	if vSpec == nil {
		// note: explicitly not the empty object.
		vSpec = &PolicySpec{}
	}
	if err := postReadExtractPolicySpecFields(r, vSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSpec) {
		r.Spec = vSpec
	}
	vDryRunSpec := r.DryRunSpec
	if vDryRunSpec == nil {
		// note: explicitly not the empty object.
		vDryRunSpec = &PolicyDryRunSpec{}
	}
	if err := postReadExtractPolicyDryRunSpecFields(r, vDryRunSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDryRunSpec) {
		r.DryRunSpec = vDryRunSpec
	}
	return nil
}
func postReadExtractPolicySpecFields(r *Policy, o *PolicySpec) error {
	return nil
}
func postReadExtractPolicySpecRulesFields(r *Policy, o *PolicySpecRules) error {
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &PolicySpecRulesValues{}
	}
	if err := extractPolicySpecRulesValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValues) {
		o.Values = vValues
	}
	vCondition := o.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &PolicySpecRulesCondition{}
	}
	if err := extractPolicySpecRulesConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		o.Condition = vCondition
	}
	return nil
}
func postReadExtractPolicySpecRulesValuesFields(r *Policy, o *PolicySpecRulesValues) error {
	return nil
}
func postReadExtractPolicySpecRulesConditionFields(r *Policy, o *PolicySpecRulesCondition) error {
	return nil
}
func postReadExtractPolicyDryRunSpecFields(r *Policy, o *PolicyDryRunSpec) error {
	return nil
}
func postReadExtractPolicyDryRunSpecRulesFields(r *Policy, o *PolicyDryRunSpecRules) error {
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &PolicyDryRunSpecRulesValues{}
	}
	if err := extractPolicyDryRunSpecRulesValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValues) {
		o.Values = vValues
	}
	vCondition := o.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &PolicyDryRunSpecRulesCondition{}
	}
	if err := extractPolicyDryRunSpecRulesConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		o.Condition = vCondition
	}
	return nil
}
func postReadExtractPolicyDryRunSpecRulesValuesFields(r *Policy, o *PolicyDryRunSpecRulesValues) error {
	return nil
}
func postReadExtractPolicyDryRunSpecRulesConditionFields(r *Policy, o *PolicyDryRunSpecRulesCondition) error {
	return nil
}
