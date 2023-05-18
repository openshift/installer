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
package firebaserules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Ruleset) validate() error {

	if err := dcl.Required(r, "source"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Source) {
		if err := r.Source.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Metadata) {
		if err := r.Metadata.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *RulesetSource) validate() error {
	if err := dcl.Required(r, "files"); err != nil {
		return err
	}
	return nil
}
func (r *RulesetSourceFiles) validate() error {
	if err := dcl.Required(r, "content"); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *RulesetMetadata) validate() error {
	return nil
}
func (r *Ruleset) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://firebaserules.googleapis.com/v1/", params)
}

func (r *Ruleset) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/rulesets/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Ruleset) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/rulesets", nr.basePath(), userBasePath, params), nil

}

func (r *Ruleset) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/rulesets", nr.basePath(), userBasePath, params), nil

}

func (r *Ruleset) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/rulesets/{{name}}", nr.basePath(), userBasePath, params), nil
}

// rulesetApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type rulesetApiOperation interface {
	do(context.Context, *Ruleset, *Client) error
}

func (c *Client) listRulesetRaw(ctx context.Context, r *Ruleset, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != RulesetMaxPage {
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

type listRulesetOperation struct {
	Rulesets []map[string]interface{} `json:"rulesets"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listRuleset(ctx context.Context, r *Ruleset, pageToken string, pageSize int32) ([]*Ruleset, string, error) {
	b, err := c.listRulesetRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listRulesetOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Ruleset
	for _, v := range m.Rulesets {
		res, err := unmarshalMapRuleset(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllRuleset(ctx context.Context, f func(*Ruleset) bool, resources []*Ruleset) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteRuleset(ctx, res)
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

type deleteRulesetOperation struct{}

func (op *deleteRulesetOperation) do(ctx context.Context, r *Ruleset, c *Client) error {
	r, err := c.GetRuleset(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.InfoWithContextf(ctx, "Ruleset not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetRuleset checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete Ruleset: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createRulesetOperation struct {
	response map[string]interface{}
}

func (op *createRulesetOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createRulesetOperation) do(ctx context.Context, r *Ruleset, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	if r.Name != nil {
		// Allowing creation to continue with Name set could result in a Ruleset with the wrong Name.
		return fmt.Errorf("server-generated parameter Name was specified by user as %v, should be unspecified", dcl.ValueOrEmptyString(r.Name))
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

	// Include Name in URL substitution for initial GET request.
	m := op.response
	r.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))

	if _, err := c.GetRuleset(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getRulesetRaw(ctx context.Context, r *Ruleset) ([]byte, error) {

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

func (c *Client) rulesetDiffsForRawDesired(ctx context.Context, rawDesired *Ruleset, opts ...dcl.ApplyOption) (initial, desired *Ruleset, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Ruleset
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Ruleset); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Ruleset, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	if fetchState.Name == nil {
		// We cannot perform a get because of lack of information. We have to assume
		// that this is being created for the first time.
		desired, err := canonicalizeRulesetDesiredState(rawDesired, nil)
		return nil, desired, nil, err
	}
	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetRuleset(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Ruleset resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Ruleset resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Ruleset resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeRulesetDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Ruleset: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Ruleset: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractRulesetFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeRulesetInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Ruleset: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeRulesetDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Ruleset: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffRuleset(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeRulesetInitialState(rawInitial, rawDesired *Ruleset) (*Ruleset, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeRulesetDesiredState(rawDesired, rawInitial *Ruleset, opts ...dcl.ApplyOption) (*Ruleset, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Source = canonicalizeRulesetSource(rawDesired.Source, nil, opts...)
		rawDesired.Metadata = canonicalizeRulesetMetadata(rawDesired.Metadata, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Ruleset{}
	if dcl.IsZeroValue(rawDesired.Name) || (dcl.IsEmptyValueIndirect(rawDesired.Name) && dcl.IsEmptyValueIndirect(rawInitial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.Source = canonicalizeRulesetSource(rawDesired.Source, rawInitial.Source, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeRulesetNewState(c *Client, rawNew, rawDesired *Ruleset) (*Ruleset, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Source) && dcl.IsEmptyValueIndirect(rawDesired.Source) {
		rawNew.Source = rawDesired.Source
	} else {
		rawNew.Source = canonicalizeNewRulesetSource(c, rawDesired.Source, rawNew.Source)
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Metadata) && dcl.IsEmptyValueIndirect(rawDesired.Metadata) {
		rawNew.Metadata = rawDesired.Metadata
	} else {
		rawNew.Metadata = canonicalizeNewRulesetMetadata(c, rawDesired.Metadata, rawNew.Metadata)
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeRulesetSource(des, initial *RulesetSource, opts ...dcl.ApplyOption) *RulesetSource {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &RulesetSource{}

	cDes.Files = canonicalizeRulesetSourceFilesSlice(des.Files, initial.Files, opts...)
	if dcl.IsZeroValue(des.Language) || (dcl.IsEmptyValueIndirect(des.Language) && dcl.IsEmptyValueIndirect(initial.Language)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Language = initial.Language
	} else {
		cDes.Language = des.Language
	}

	return cDes
}

func canonicalizeRulesetSourceSlice(des, initial []RulesetSource, opts ...dcl.ApplyOption) []RulesetSource {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]RulesetSource, 0, len(des))
		for _, d := range des {
			cd := canonicalizeRulesetSource(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]RulesetSource, 0, len(des))
	for i, d := range des {
		cd := canonicalizeRulesetSource(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewRulesetSource(c *Client, des, nw *RulesetSource) *RulesetSource {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for RulesetSource while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Files = canonicalizeNewRulesetSourceFilesSlice(c, des.Files, nw.Files)

	return nw
}

func canonicalizeNewRulesetSourceSet(c *Client, des, nw []RulesetSource) []RulesetSource {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []RulesetSource
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareRulesetSourceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewRulesetSource(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewRulesetSourceSlice(c *Client, des, nw []RulesetSource) []RulesetSource {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []RulesetSource
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewRulesetSource(c, &d, &n))
	}

	return items
}

func canonicalizeRulesetSourceFiles(des, initial *RulesetSourceFiles, opts ...dcl.ApplyOption) *RulesetSourceFiles {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &RulesetSourceFiles{}

	if dcl.StringCanonicalize(des.Content, initial.Content) || dcl.IsZeroValue(des.Content) {
		cDes.Content = initial.Content
	} else {
		cDes.Content = des.Content
	}
	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Fingerprint, initial.Fingerprint) || dcl.IsZeroValue(des.Fingerprint) {
		cDes.Fingerprint = initial.Fingerprint
	} else {
		cDes.Fingerprint = des.Fingerprint
	}

	return cDes
}

func canonicalizeRulesetSourceFilesSlice(des, initial []RulesetSourceFiles, opts ...dcl.ApplyOption) []RulesetSourceFiles {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]RulesetSourceFiles, 0, len(des))
		for _, d := range des {
			cd := canonicalizeRulesetSourceFiles(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]RulesetSourceFiles, 0, len(des))
	for i, d := range des {
		cd := canonicalizeRulesetSourceFiles(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewRulesetSourceFiles(c *Client, des, nw *RulesetSourceFiles) *RulesetSourceFiles {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for RulesetSourceFiles while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Content, nw.Content) {
		nw.Content = des.Content
	}
	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Fingerprint, nw.Fingerprint) {
		nw.Fingerprint = des.Fingerprint
	}

	return nw
}

func canonicalizeNewRulesetSourceFilesSet(c *Client, des, nw []RulesetSourceFiles) []RulesetSourceFiles {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []RulesetSourceFiles
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareRulesetSourceFilesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewRulesetSourceFiles(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewRulesetSourceFilesSlice(c *Client, des, nw []RulesetSourceFiles) []RulesetSourceFiles {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []RulesetSourceFiles
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewRulesetSourceFiles(c, &d, &n))
	}

	return items
}

func canonicalizeRulesetMetadata(des, initial *RulesetMetadata, opts ...dcl.ApplyOption) *RulesetMetadata {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &RulesetMetadata{}

	if dcl.StringArrayCanonicalize(des.Services, initial.Services) {
		cDes.Services = initial.Services
	} else {
		cDes.Services = des.Services
	}

	return cDes
}

func canonicalizeRulesetMetadataSlice(des, initial []RulesetMetadata, opts ...dcl.ApplyOption) []RulesetMetadata {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]RulesetMetadata, 0, len(des))
		for _, d := range des {
			cd := canonicalizeRulesetMetadata(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]RulesetMetadata, 0, len(des))
	for i, d := range des {
		cd := canonicalizeRulesetMetadata(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewRulesetMetadata(c *Client, des, nw *RulesetMetadata) *RulesetMetadata {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for RulesetMetadata while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Services, nw.Services) {
		nw.Services = des.Services
	}

	return nw
}

func canonicalizeNewRulesetMetadataSet(c *Client, des, nw []RulesetMetadata) []RulesetMetadata {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []RulesetMetadata
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareRulesetMetadataNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewRulesetMetadata(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewRulesetMetadataSlice(c *Client, des, nw []RulesetMetadata) []RulesetMetadata {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []RulesetMetadata
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewRulesetMetadata(c, &d, &n))
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
func diffRuleset(c *Client, desired, actual *Ruleset, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.DiffInfo{ObjectFunction: compareRulesetSourceNewStyle, EmptyObject: EmptyRulesetSource, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareRulesetMetadataNewStyle, EmptyObject: EmptyRulesetMetadata, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
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
func compareRulesetSourceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*RulesetSource)
	if !ok {
		desiredNotPointer, ok := d.(RulesetSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetSource or *RulesetSource", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*RulesetSource)
	if !ok {
		actualNotPointer, ok := a.(RulesetSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetSource", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Files, actual.Files, dcl.DiffInfo{ObjectFunction: compareRulesetSourceFilesNewStyle, EmptyObject: EmptyRulesetSourceFiles, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Files")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Language, actual.Language, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Language")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareRulesetSourceFilesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*RulesetSourceFiles)
	if !ok {
		desiredNotPointer, ok := d.(RulesetSourceFiles)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetSourceFiles or *RulesetSourceFiles", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*RulesetSourceFiles)
	if !ok {
		actualNotPointer, ok := a.(RulesetSourceFiles)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetSourceFiles", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Content, actual.Content, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Content")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fingerprint, actual.Fingerprint, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareRulesetMetadataNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*RulesetMetadata)
	if !ok {
		desiredNotPointer, ok := d.(RulesetMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetMetadata or *RulesetMetadata", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*RulesetMetadata)
	if !ok {
		actualNotPointer, ok := a.(RulesetMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RulesetMetadata", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Services, actual.Services, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Services")); len(ds) != 0 || err != nil {
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
func (r *Ruleset) urlNormalized() *Ruleset {
	normalized := dcl.Copy(*r).(Ruleset)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Ruleset) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Ruleset resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Ruleset) marshal(c *Client) ([]byte, error) {
	m, err := expandRuleset(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Ruleset: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalRuleset decodes JSON responses into the Ruleset resource schema.
func unmarshalRuleset(b []byte, c *Client, res *Ruleset) (*Ruleset, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapRuleset(m, c, res)
}

func unmarshalMapRuleset(m map[string]interface{}, c *Client, res *Ruleset) (*Ruleset, error) {

	flattened := flattenRuleset(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandRuleset expands Ruleset into a JSON request object.
func expandRuleset(c *Client, f *Ruleset) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/rulesets/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := expandRulesetSource(c, f.Source, res); err != nil {
		return nil, fmt.Errorf("error expanding Source into source: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenRuleset flattens Ruleset from a JSON request object into the
// Ruleset type.
func flattenRuleset(c *Client, i interface{}, res *Ruleset) *Ruleset {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Ruleset{}
	resultRes.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))
	resultRes.Source = flattenRulesetSource(c, m["source"], res)
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.Metadata = flattenRulesetMetadata(c, m["metadata"], res)
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandRulesetSourceMap expands the contents of RulesetSource into a JSON
// request object.
func expandRulesetSourceMap(c *Client, f map[string]RulesetSource, res *Ruleset) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandRulesetSource(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandRulesetSourceSlice expands the contents of RulesetSource into a JSON
// request object.
func expandRulesetSourceSlice(c *Client, f []RulesetSource, res *Ruleset) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandRulesetSource(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenRulesetSourceMap flattens the contents of RulesetSource from a JSON
// response object.
func flattenRulesetSourceMap(c *Client, i interface{}, res *Ruleset) map[string]RulesetSource {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RulesetSource{}
	}

	if len(a) == 0 {
		return map[string]RulesetSource{}
	}

	items := make(map[string]RulesetSource)
	for k, item := range a {
		items[k] = *flattenRulesetSource(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenRulesetSourceSlice flattens the contents of RulesetSource from a JSON
// response object.
func flattenRulesetSourceSlice(c *Client, i interface{}, res *Ruleset) []RulesetSource {
	a, ok := i.([]interface{})
	if !ok {
		return []RulesetSource{}
	}

	if len(a) == 0 {
		return []RulesetSource{}
	}

	items := make([]RulesetSource, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRulesetSource(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandRulesetSource expands an instance of RulesetSource into a JSON
// request object.
func expandRulesetSource(c *Client, f *RulesetSource, res *Ruleset) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandRulesetSourceFilesSlice(c, f.Files, res); err != nil {
		return nil, fmt.Errorf("error expanding Files into files: %w", err)
	} else if v != nil {
		m["files"] = v
	}
	if v := f.Language; !dcl.IsEmptyValueIndirect(v) {
		m["language"] = v
	}

	return m, nil
}

// flattenRulesetSource flattens an instance of RulesetSource from a JSON
// response object.
func flattenRulesetSource(c *Client, i interface{}, res *Ruleset) *RulesetSource {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &RulesetSource{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyRulesetSource
	}
	r.Files = flattenRulesetSourceFilesSlice(c, m["files"], res)
	r.Language = flattenRulesetSourceLanguageEnum(m["language"])

	return r
}

// expandRulesetSourceFilesMap expands the contents of RulesetSourceFiles into a JSON
// request object.
func expandRulesetSourceFilesMap(c *Client, f map[string]RulesetSourceFiles, res *Ruleset) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandRulesetSourceFiles(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandRulesetSourceFilesSlice expands the contents of RulesetSourceFiles into a JSON
// request object.
func expandRulesetSourceFilesSlice(c *Client, f []RulesetSourceFiles, res *Ruleset) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandRulesetSourceFiles(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenRulesetSourceFilesMap flattens the contents of RulesetSourceFiles from a JSON
// response object.
func flattenRulesetSourceFilesMap(c *Client, i interface{}, res *Ruleset) map[string]RulesetSourceFiles {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RulesetSourceFiles{}
	}

	if len(a) == 0 {
		return map[string]RulesetSourceFiles{}
	}

	items := make(map[string]RulesetSourceFiles)
	for k, item := range a {
		items[k] = *flattenRulesetSourceFiles(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenRulesetSourceFilesSlice flattens the contents of RulesetSourceFiles from a JSON
// response object.
func flattenRulesetSourceFilesSlice(c *Client, i interface{}, res *Ruleset) []RulesetSourceFiles {
	a, ok := i.([]interface{})
	if !ok {
		return []RulesetSourceFiles{}
	}

	if len(a) == 0 {
		return []RulesetSourceFiles{}
	}

	items := make([]RulesetSourceFiles, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRulesetSourceFiles(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandRulesetSourceFiles expands an instance of RulesetSourceFiles into a JSON
// request object.
func expandRulesetSourceFiles(c *Client, f *RulesetSourceFiles, res *Ruleset) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Content; !dcl.IsEmptyValueIndirect(v) {
		m["content"] = v
	}
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Fingerprint; !dcl.IsEmptyValueIndirect(v) {
		m["fingerprint"] = v
	}

	return m, nil
}

// flattenRulesetSourceFiles flattens an instance of RulesetSourceFiles from a JSON
// response object.
func flattenRulesetSourceFiles(c *Client, i interface{}, res *Ruleset) *RulesetSourceFiles {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &RulesetSourceFiles{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyRulesetSourceFiles
	}
	r.Content = dcl.FlattenString(m["content"])
	r.Name = dcl.FlattenString(m["name"])
	r.Fingerprint = dcl.FlattenString(m["fingerprint"])

	return r
}

// expandRulesetMetadataMap expands the contents of RulesetMetadata into a JSON
// request object.
func expandRulesetMetadataMap(c *Client, f map[string]RulesetMetadata, res *Ruleset) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandRulesetMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandRulesetMetadataSlice expands the contents of RulesetMetadata into a JSON
// request object.
func expandRulesetMetadataSlice(c *Client, f []RulesetMetadata, res *Ruleset) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandRulesetMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenRulesetMetadataMap flattens the contents of RulesetMetadata from a JSON
// response object.
func flattenRulesetMetadataMap(c *Client, i interface{}, res *Ruleset) map[string]RulesetMetadata {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RulesetMetadata{}
	}

	if len(a) == 0 {
		return map[string]RulesetMetadata{}
	}

	items := make(map[string]RulesetMetadata)
	for k, item := range a {
		items[k] = *flattenRulesetMetadata(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenRulesetMetadataSlice flattens the contents of RulesetMetadata from a JSON
// response object.
func flattenRulesetMetadataSlice(c *Client, i interface{}, res *Ruleset) []RulesetMetadata {
	a, ok := i.([]interface{})
	if !ok {
		return []RulesetMetadata{}
	}

	if len(a) == 0 {
		return []RulesetMetadata{}
	}

	items := make([]RulesetMetadata, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRulesetMetadata(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandRulesetMetadata expands an instance of RulesetMetadata into a JSON
// request object.
func expandRulesetMetadata(c *Client, f *RulesetMetadata, res *Ruleset) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Services; v != nil {
		m["services"] = v
	}

	return m, nil
}

// flattenRulesetMetadata flattens an instance of RulesetMetadata from a JSON
// response object.
func flattenRulesetMetadata(c *Client, i interface{}, res *Ruleset) *RulesetMetadata {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &RulesetMetadata{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyRulesetMetadata
	}
	r.Services = dcl.FlattenStringSlice(m["services"])

	return r
}

// flattenRulesetSourceLanguageEnumMap flattens the contents of RulesetSourceLanguageEnum from a JSON
// response object.
func flattenRulesetSourceLanguageEnumMap(c *Client, i interface{}, res *Ruleset) map[string]RulesetSourceLanguageEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RulesetSourceLanguageEnum{}
	}

	if len(a) == 0 {
		return map[string]RulesetSourceLanguageEnum{}
	}

	items := make(map[string]RulesetSourceLanguageEnum)
	for k, item := range a {
		items[k] = *flattenRulesetSourceLanguageEnum(item.(interface{}))
	}

	return items
}

// flattenRulesetSourceLanguageEnumSlice flattens the contents of RulesetSourceLanguageEnum from a JSON
// response object.
func flattenRulesetSourceLanguageEnumSlice(c *Client, i interface{}, res *Ruleset) []RulesetSourceLanguageEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []RulesetSourceLanguageEnum{}
	}

	if len(a) == 0 {
		return []RulesetSourceLanguageEnum{}
	}

	items := make([]RulesetSourceLanguageEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRulesetSourceLanguageEnum(item.(interface{})))
	}

	return items
}

// flattenRulesetSourceLanguageEnum asserts that an interface is a string, and returns a
// pointer to a *RulesetSourceLanguageEnum with the same value as that string.
func flattenRulesetSourceLanguageEnum(i interface{}) *RulesetSourceLanguageEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return RulesetSourceLanguageEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Ruleset) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalRuleset(b, c, r)
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

type rulesetDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         rulesetApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToRulesetDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]rulesetDiff, error) {
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
	var diffs []rulesetDiff
	// For each operation name, create a rulesetDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := rulesetDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToRulesetApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToRulesetApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (rulesetApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractRulesetFields(r *Ruleset) error {
	vSource := r.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &RulesetSource{}
	}
	if err := extractRulesetSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		r.Source = vSource
	}
	vMetadata := r.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &RulesetMetadata{}
	}
	if err := extractRulesetMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		r.Metadata = vMetadata
	}
	return nil
}
func extractRulesetSourceFields(r *Ruleset, o *RulesetSource) error {
	return nil
}
func extractRulesetSourceFilesFields(r *Ruleset, o *RulesetSourceFiles) error {
	return nil
}
func extractRulesetMetadataFields(r *Ruleset, o *RulesetMetadata) error {
	return nil
}

func postReadExtractRulesetFields(r *Ruleset) error {
	vSource := r.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &RulesetSource{}
	}
	if err := postReadExtractRulesetSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		r.Source = vSource
	}
	vMetadata := r.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &RulesetMetadata{}
	}
	if err := postReadExtractRulesetMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		r.Metadata = vMetadata
	}
	return nil
}
func postReadExtractRulesetSourceFields(r *Ruleset, o *RulesetSource) error {
	return nil
}
func postReadExtractRulesetSourceFilesFields(r *Ruleset, o *RulesetSourceFiles) error {
	return nil
}
func postReadExtractRulesetMetadataFields(r *Ruleset, o *RulesetMetadata) error {
	return nil
}
