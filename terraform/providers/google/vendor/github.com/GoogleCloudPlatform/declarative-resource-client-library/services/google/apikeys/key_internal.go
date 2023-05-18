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
package apikeys

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

func (r *Key) validate() error {

	if err := dcl.RequiredParameter(r.Name, "Name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Restrictions) {
		if err := r.Restrictions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *KeyRestrictions) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"BrowserKeyRestrictions", "ServerKeyRestrictions", "AndroidKeyRestrictions", "IosKeyRestrictions"}, r.BrowserKeyRestrictions, r.ServerKeyRestrictions, r.AndroidKeyRestrictions, r.IosKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.BrowserKeyRestrictions) {
		if err := r.BrowserKeyRestrictions.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ServerKeyRestrictions) {
		if err := r.ServerKeyRestrictions.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AndroidKeyRestrictions) {
		if err := r.AndroidKeyRestrictions.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.IosKeyRestrictions) {
		if err := r.IosKeyRestrictions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *KeyRestrictionsBrowserKeyRestrictions) validate() error {
	if err := dcl.Required(r, "allowedReferrers"); err != nil {
		return err
	}
	return nil
}
func (r *KeyRestrictionsServerKeyRestrictions) validate() error {
	if err := dcl.Required(r, "allowedIps"); err != nil {
		return err
	}
	return nil
}
func (r *KeyRestrictionsAndroidKeyRestrictions) validate() error {
	if err := dcl.Required(r, "allowedApplications"); err != nil {
		return err
	}
	return nil
}
func (r *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) validate() error {
	if err := dcl.Required(r, "sha1Fingerprint"); err != nil {
		return err
	}
	if err := dcl.Required(r, "packageName"); err != nil {
		return err
	}
	return nil
}
func (r *KeyRestrictionsIosKeyRestrictions) validate() error {
	if err := dcl.Required(r, "allowedBundleIds"); err != nil {
		return err
	}
	return nil
}
func (r *KeyRestrictionsApiTargets) validate() error {
	if err := dcl.Required(r, "service"); err != nil {
		return err
	}
	return nil
}
func (r *Key) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://apikeys.googleapis.com/v2/", params)
}

func (r *Key) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/keys/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Key) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/locations/global/keys", nr.basePath(), userBasePath, params), nil

}

func (r *Key) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/keys?keyId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Key) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/keys/{{name}}", nr.basePath(), userBasePath, params), nil
}

// keyApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type keyApiOperation interface {
	do(context.Context, *Key, *Client) error
}

// newUpdateKeyUpdateKeyRequest creates a request for an
// Key resource's UpdateKey update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateKeyUpdateKeyRequest(ctx context.Context, f *Key, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v, err := expandKeyRestrictions(c, f.Restrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding Restrictions into restrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["restrictions"] = v
	}
	return req, nil
}

// marshalUpdateKeyUpdateKeyRequest converts the update into
// the final JSON request body.
func marshalUpdateKeyUpdateKeyRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateKeyUpdateKeyOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateKeyUpdateKeyOperation) do(ctx context.Context, r *Key, c *Client) error {
	_, err := c.GetKey(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateKey")
	if err != nil {
		return err
	}
	mask := dcl.TopLevelUpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateKeyUpdateKeyRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateKeyUpdateKeyRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listKeyRaw(ctx context.Context, r *Key, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != KeyMaxPage {
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

type listKeyOperation struct {
	Keys  []map[string]interface{} `json:"keys"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listKey(ctx context.Context, r *Key, pageToken string, pageSize int32) ([]*Key, string, error) {
	b, err := c.listKeyRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listKeyOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Key
	for _, v := range m.Keys {
		res, err := unmarshalMapKey(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllKey(ctx context.Context, f func(*Key) bool, resources []*Key) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteKey(ctx, res)
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

type deleteKeyOperation struct{}

func (op *deleteKeyOperation) do(ctx context.Context, r *Key, c *Client) error {
	r, err := c.GetKey(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Key not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetKey checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		return err
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createKeyOperation struct {
	response map[string]interface{}
}

func (op *createKeyOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createKeyOperation) do(ctx context.Context, r *Key, c *Client) error {
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
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetKey(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) keyDiffsForRawDesired(ctx context.Context, rawDesired *Key, opts ...dcl.ApplyOption) (initial, desired *Key, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Key
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Key); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Key, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetKey(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Key resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Key resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Key resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeKeyDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Key: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Key: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractKeyFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeKeyInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Key: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeKeyDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Key: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffKey(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeKeyInitialState(rawInitial, rawDesired *Key) (*Key, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeKeyDesiredState(rawDesired, rawInitial *Key, opts ...dcl.ApplyOption) (*Key, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Restrictions = canonicalizeKeyRestrictions(rawDesired.Restrictions, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Key{}
	if dcl.NameToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	canonicalDesired.Restrictions = canonicalizeKeyRestrictions(rawDesired.Restrictions, rawInitial.Restrictions, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeKeyNewState(c *Client, rawNew, rawDesired *Key) (*Key, error) {

	rawNew.Name = rawDesired.Name

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.KeyString) && dcl.IsEmptyValueIndirect(rawDesired.KeyString) {
		rawNew.KeyString = rawDesired.KeyString
	} else {
		if dcl.StringCanonicalize(rawDesired.KeyString, rawNew.KeyString) {
			rawNew.KeyString = rawDesired.KeyString
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Restrictions) && dcl.IsEmptyValueIndirect(rawDesired.Restrictions) {
		rawNew.Restrictions = rawDesired.Restrictions
	} else {
		rawNew.Restrictions = canonicalizeNewKeyRestrictions(c, rawDesired.Restrictions, rawNew.Restrictions)
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeKeyRestrictions(des, initial *KeyRestrictions, opts ...dcl.ApplyOption) *KeyRestrictions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.BrowserKeyRestrictions != nil || (initial != nil && initial.BrowserKeyRestrictions != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ServerKeyRestrictions, des.AndroidKeyRestrictions, des.IosKeyRestrictions) {
			des.BrowserKeyRestrictions = nil
			if initial != nil {
				initial.BrowserKeyRestrictions = nil
			}
		}
	}

	if des.ServerKeyRestrictions != nil || (initial != nil && initial.ServerKeyRestrictions != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BrowserKeyRestrictions, des.AndroidKeyRestrictions, des.IosKeyRestrictions) {
			des.ServerKeyRestrictions = nil
			if initial != nil {
				initial.ServerKeyRestrictions = nil
			}
		}
	}

	if des.AndroidKeyRestrictions != nil || (initial != nil && initial.AndroidKeyRestrictions != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BrowserKeyRestrictions, des.ServerKeyRestrictions, des.IosKeyRestrictions) {
			des.AndroidKeyRestrictions = nil
			if initial != nil {
				initial.AndroidKeyRestrictions = nil
			}
		}
	}

	if des.IosKeyRestrictions != nil || (initial != nil && initial.IosKeyRestrictions != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BrowserKeyRestrictions, des.ServerKeyRestrictions, des.AndroidKeyRestrictions) {
			des.IosKeyRestrictions = nil
			if initial != nil {
				initial.IosKeyRestrictions = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictions{}

	cDes.BrowserKeyRestrictions = canonicalizeKeyRestrictionsBrowserKeyRestrictions(des.BrowserKeyRestrictions, initial.BrowserKeyRestrictions, opts...)
	cDes.ServerKeyRestrictions = canonicalizeKeyRestrictionsServerKeyRestrictions(des.ServerKeyRestrictions, initial.ServerKeyRestrictions, opts...)
	cDes.AndroidKeyRestrictions = canonicalizeKeyRestrictionsAndroidKeyRestrictions(des.AndroidKeyRestrictions, initial.AndroidKeyRestrictions, opts...)
	cDes.IosKeyRestrictions = canonicalizeKeyRestrictionsIosKeyRestrictions(des.IosKeyRestrictions, initial.IosKeyRestrictions, opts...)
	cDes.ApiTargets = canonicalizeKeyRestrictionsApiTargetsSlice(des.ApiTargets, initial.ApiTargets, opts...)

	return cDes
}

func canonicalizeKeyRestrictionsSlice(des, initial []KeyRestrictions, opts ...dcl.ApplyOption) []KeyRestrictions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictions(c *Client, des, nw *KeyRestrictions) *KeyRestrictions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BrowserKeyRestrictions = canonicalizeNewKeyRestrictionsBrowserKeyRestrictions(c, des.BrowserKeyRestrictions, nw.BrowserKeyRestrictions)
	nw.ServerKeyRestrictions = canonicalizeNewKeyRestrictionsServerKeyRestrictions(c, des.ServerKeyRestrictions, nw.ServerKeyRestrictions)
	nw.AndroidKeyRestrictions = canonicalizeNewKeyRestrictionsAndroidKeyRestrictions(c, des.AndroidKeyRestrictions, nw.AndroidKeyRestrictions)
	nw.IosKeyRestrictions = canonicalizeNewKeyRestrictionsIosKeyRestrictions(c, des.IosKeyRestrictions, nw.IosKeyRestrictions)
	nw.ApiTargets = canonicalizeNewKeyRestrictionsApiTargetsSlice(c, des.ApiTargets, nw.ApiTargets)

	return nw
}

func canonicalizeNewKeyRestrictionsSet(c *Client, des, nw []KeyRestrictions) []KeyRestrictions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsSlice(c *Client, des, nw []KeyRestrictions) []KeyRestrictions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictions(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsBrowserKeyRestrictions(des, initial *KeyRestrictionsBrowserKeyRestrictions, opts ...dcl.ApplyOption) *KeyRestrictionsBrowserKeyRestrictions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsBrowserKeyRestrictions{}

	if dcl.StringArrayCanonicalize(des.AllowedReferrers, initial.AllowedReferrers) {
		cDes.AllowedReferrers = initial.AllowedReferrers
	} else {
		cDes.AllowedReferrers = des.AllowedReferrers
	}

	return cDes
}

func canonicalizeKeyRestrictionsBrowserKeyRestrictionsSlice(des, initial []KeyRestrictionsBrowserKeyRestrictions, opts ...dcl.ApplyOption) []KeyRestrictionsBrowserKeyRestrictions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsBrowserKeyRestrictions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsBrowserKeyRestrictions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsBrowserKeyRestrictions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsBrowserKeyRestrictions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsBrowserKeyRestrictions(c *Client, des, nw *KeyRestrictionsBrowserKeyRestrictions) *KeyRestrictionsBrowserKeyRestrictions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsBrowserKeyRestrictions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.AllowedReferrers, nw.AllowedReferrers) {
		nw.AllowedReferrers = des.AllowedReferrers
	}

	return nw
}

func canonicalizeNewKeyRestrictionsBrowserKeyRestrictionsSet(c *Client, des, nw []KeyRestrictionsBrowserKeyRestrictions) []KeyRestrictionsBrowserKeyRestrictions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsBrowserKeyRestrictions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsBrowserKeyRestrictionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsBrowserKeyRestrictions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsBrowserKeyRestrictionsSlice(c *Client, des, nw []KeyRestrictionsBrowserKeyRestrictions) []KeyRestrictionsBrowserKeyRestrictions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsBrowserKeyRestrictions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsBrowserKeyRestrictions(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsServerKeyRestrictions(des, initial *KeyRestrictionsServerKeyRestrictions, opts ...dcl.ApplyOption) *KeyRestrictionsServerKeyRestrictions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsServerKeyRestrictions{}

	if dcl.StringArrayCanonicalize(des.AllowedIps, initial.AllowedIps) {
		cDes.AllowedIps = initial.AllowedIps
	} else {
		cDes.AllowedIps = des.AllowedIps
	}

	return cDes
}

func canonicalizeKeyRestrictionsServerKeyRestrictionsSlice(des, initial []KeyRestrictionsServerKeyRestrictions, opts ...dcl.ApplyOption) []KeyRestrictionsServerKeyRestrictions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsServerKeyRestrictions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsServerKeyRestrictions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsServerKeyRestrictions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsServerKeyRestrictions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsServerKeyRestrictions(c *Client, des, nw *KeyRestrictionsServerKeyRestrictions) *KeyRestrictionsServerKeyRestrictions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsServerKeyRestrictions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.AllowedIps, nw.AllowedIps) {
		nw.AllowedIps = des.AllowedIps
	}

	return nw
}

func canonicalizeNewKeyRestrictionsServerKeyRestrictionsSet(c *Client, des, nw []KeyRestrictionsServerKeyRestrictions) []KeyRestrictionsServerKeyRestrictions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsServerKeyRestrictions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsServerKeyRestrictionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsServerKeyRestrictions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsServerKeyRestrictionsSlice(c *Client, des, nw []KeyRestrictionsServerKeyRestrictions) []KeyRestrictionsServerKeyRestrictions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsServerKeyRestrictions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsServerKeyRestrictions(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsAndroidKeyRestrictions(des, initial *KeyRestrictionsAndroidKeyRestrictions, opts ...dcl.ApplyOption) *KeyRestrictionsAndroidKeyRestrictions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsAndroidKeyRestrictions{}

	cDes.AllowedApplications = canonicalizeKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(des.AllowedApplications, initial.AllowedApplications, opts...)

	return cDes
}

func canonicalizeKeyRestrictionsAndroidKeyRestrictionsSlice(des, initial []KeyRestrictionsAndroidKeyRestrictions, opts ...dcl.ApplyOption) []KeyRestrictionsAndroidKeyRestrictions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsAndroidKeyRestrictions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsAndroidKeyRestrictions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsAndroidKeyRestrictions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsAndroidKeyRestrictions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictions(c *Client, des, nw *KeyRestrictionsAndroidKeyRestrictions) *KeyRestrictionsAndroidKeyRestrictions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsAndroidKeyRestrictions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AllowedApplications = canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c, des.AllowedApplications, nw.AllowedApplications)

	return nw
}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsSet(c *Client, des, nw []KeyRestrictionsAndroidKeyRestrictions) []KeyRestrictionsAndroidKeyRestrictions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsAndroidKeyRestrictions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsAndroidKeyRestrictionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsAndroidKeyRestrictions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsSlice(c *Client, des, nw []KeyRestrictionsAndroidKeyRestrictions) []KeyRestrictionsAndroidKeyRestrictions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsAndroidKeyRestrictions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsAndroidKeyRestrictions(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(des, initial *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, opts ...dcl.ApplyOption) *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}

	if dcl.StringCanonicalize(des.Sha1Fingerprint, initial.Sha1Fingerprint) || dcl.IsZeroValue(des.Sha1Fingerprint) {
		cDes.Sha1Fingerprint = initial.Sha1Fingerprint
	} else {
		cDes.Sha1Fingerprint = des.Sha1Fingerprint
	}
	if dcl.StringCanonicalize(des.PackageName, initial.PackageName) || dcl.IsZeroValue(des.PackageName) {
		cDes.PackageName = initial.PackageName
	} else {
		cDes.PackageName = des.PackageName
	}

	return cDes
}

func canonicalizeKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(des, initial []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, opts ...dcl.ApplyOption) []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c *Client, des, nw *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsAndroidKeyRestrictionsAllowedApplications while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Sha1Fingerprint, nw.Sha1Fingerprint) {
		nw.Sha1Fingerprint = des.Sha1Fingerprint
	}
	if dcl.StringCanonicalize(des.PackageName, nw.PackageName) {
		nw.PackageName = des.PackageName
	}

	return nw
}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSet(c *Client, des, nw []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c *Client, des, nw []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsIosKeyRestrictions(des, initial *KeyRestrictionsIosKeyRestrictions, opts ...dcl.ApplyOption) *KeyRestrictionsIosKeyRestrictions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsIosKeyRestrictions{}

	if dcl.StringArrayCanonicalize(des.AllowedBundleIds, initial.AllowedBundleIds) {
		cDes.AllowedBundleIds = initial.AllowedBundleIds
	} else {
		cDes.AllowedBundleIds = des.AllowedBundleIds
	}

	return cDes
}

func canonicalizeKeyRestrictionsIosKeyRestrictionsSlice(des, initial []KeyRestrictionsIosKeyRestrictions, opts ...dcl.ApplyOption) []KeyRestrictionsIosKeyRestrictions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsIosKeyRestrictions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsIosKeyRestrictions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsIosKeyRestrictions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsIosKeyRestrictions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsIosKeyRestrictions(c *Client, des, nw *KeyRestrictionsIosKeyRestrictions) *KeyRestrictionsIosKeyRestrictions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsIosKeyRestrictions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.AllowedBundleIds, nw.AllowedBundleIds) {
		nw.AllowedBundleIds = des.AllowedBundleIds
	}

	return nw
}

func canonicalizeNewKeyRestrictionsIosKeyRestrictionsSet(c *Client, des, nw []KeyRestrictionsIosKeyRestrictions) []KeyRestrictionsIosKeyRestrictions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsIosKeyRestrictions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsIosKeyRestrictionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsIosKeyRestrictions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsIosKeyRestrictionsSlice(c *Client, des, nw []KeyRestrictionsIosKeyRestrictions) []KeyRestrictionsIosKeyRestrictions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsIosKeyRestrictions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsIosKeyRestrictions(c, &d, &n))
	}

	return items
}

func canonicalizeKeyRestrictionsApiTargets(des, initial *KeyRestrictionsApiTargets, opts ...dcl.ApplyOption) *KeyRestrictionsApiTargets {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyRestrictionsApiTargets{}

	if dcl.StringCanonicalize(des.Service, initial.Service) || dcl.IsZeroValue(des.Service) {
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}
	if dcl.StringArrayCanonicalize(des.Methods, initial.Methods) {
		cDes.Methods = initial.Methods
	} else {
		cDes.Methods = des.Methods
	}

	return cDes
}

func canonicalizeKeyRestrictionsApiTargetsSlice(des, initial []KeyRestrictionsApiTargets, opts ...dcl.ApplyOption) []KeyRestrictionsApiTargets {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyRestrictionsApiTargets, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyRestrictionsApiTargets(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyRestrictionsApiTargets, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyRestrictionsApiTargets(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyRestrictionsApiTargets(c *Client, des, nw *KeyRestrictionsApiTargets) *KeyRestrictionsApiTargets {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyRestrictionsApiTargets while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Service, nw.Service) {
		nw.Service = des.Service
	}
	if dcl.StringArrayCanonicalize(des.Methods, nw.Methods) {
		nw.Methods = des.Methods
	}

	return nw
}

func canonicalizeNewKeyRestrictionsApiTargetsSet(c *Client, des, nw []KeyRestrictionsApiTargets) []KeyRestrictionsApiTargets {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyRestrictionsApiTargets
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyRestrictionsApiTargetsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyRestrictionsApiTargets(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyRestrictionsApiTargetsSlice(c *Client, des, nw []KeyRestrictionsApiTargets) []KeyRestrictionsApiTargets {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyRestrictionsApiTargets
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyRestrictionsApiTargets(c, &d, &n))
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
func diffKey(c *Client, desired, actual *Key, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyString, actual.KeyString, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyString")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Restrictions, actual.Restrictions, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsNewStyle, EmptyObject: EmptyKeyRestrictions, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("Restrictions")); len(ds) != 0 || err != nil {
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
func compareKeyRestrictionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictions)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictions or *KeyRestrictions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictions)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BrowserKeyRestrictions, actual.BrowserKeyRestrictions, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsBrowserKeyRestrictionsNewStyle, EmptyObject: EmptyKeyRestrictionsBrowserKeyRestrictions, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("BrowserKeyRestrictions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServerKeyRestrictions, actual.ServerKeyRestrictions, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsServerKeyRestrictionsNewStyle, EmptyObject: EmptyKeyRestrictionsServerKeyRestrictions, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("ServerKeyRestrictions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AndroidKeyRestrictions, actual.AndroidKeyRestrictions, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsAndroidKeyRestrictionsNewStyle, EmptyObject: EmptyKeyRestrictionsAndroidKeyRestrictions, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AndroidKeyRestrictions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IosKeyRestrictions, actual.IosKeyRestrictions, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsIosKeyRestrictionsNewStyle, EmptyObject: EmptyKeyRestrictionsIosKeyRestrictions, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("IosKeyRestrictions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ApiTargets, actual.ApiTargets, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsApiTargetsNewStyle, EmptyObject: EmptyKeyRestrictionsApiTargets, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("ApiTargets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsBrowserKeyRestrictionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsBrowserKeyRestrictions)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsBrowserKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsBrowserKeyRestrictions or *KeyRestrictionsBrowserKeyRestrictions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsBrowserKeyRestrictions)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsBrowserKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsBrowserKeyRestrictions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedReferrers, actual.AllowedReferrers, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedReferrers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsServerKeyRestrictionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsServerKeyRestrictions)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsServerKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsServerKeyRestrictions or *KeyRestrictionsServerKeyRestrictions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsServerKeyRestrictions)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsServerKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsServerKeyRestrictions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedIps, actual.AllowedIps, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedIps")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsAndroidKeyRestrictionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsAndroidKeyRestrictions)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsAndroidKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsAndroidKeyRestrictions or *KeyRestrictionsAndroidKeyRestrictions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsAndroidKeyRestrictions)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsAndroidKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsAndroidKeyRestrictions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedApplications, actual.AllowedApplications, dcl.DiffInfo{ObjectFunction: compareKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsNewStyle, EmptyObject: EmptyKeyRestrictionsAndroidKeyRestrictionsAllowedApplications, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedApplications")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsAndroidKeyRestrictionsAllowedApplications)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsAndroidKeyRestrictionsAllowedApplications)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsAndroidKeyRestrictionsAllowedApplications or *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsAndroidKeyRestrictionsAllowedApplications)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsAndroidKeyRestrictionsAllowedApplications)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsAndroidKeyRestrictionsAllowedApplications", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Sha1Fingerprint, actual.Sha1Fingerprint, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("Sha1Fingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PackageName, actual.PackageName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("PackageName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsIosKeyRestrictionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsIosKeyRestrictions)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsIosKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsIosKeyRestrictions or *KeyRestrictionsIosKeyRestrictions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsIosKeyRestrictions)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsIosKeyRestrictions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsIosKeyRestrictions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedBundleIds, actual.AllowedBundleIds, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedBundleIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyRestrictionsApiTargetsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyRestrictionsApiTargets)
	if !ok {
		desiredNotPointer, ok := d.(KeyRestrictionsApiTargets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsApiTargets or *KeyRestrictionsApiTargets", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyRestrictionsApiTargets)
	if !ok {
		actualNotPointer, ok := a.(KeyRestrictionsApiTargets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyRestrictionsApiTargets", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Methods, actual.Methods, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("Methods")); len(ds) != 0 || err != nil {
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
func (r *Key) urlNormalized() *Key {
	normalized := dcl.Copy(*r).(Key)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.KeyString = dcl.SelfLinkToName(r.KeyString)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Key) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateKey" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/global/keys/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Key resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Key) marshal(c *Client) ([]byte, error) {
	m, err := expandKey(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Key: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalKey decodes JSON responses into the Key resource schema.
func unmarshalKey(b []byte, c *Client, res *Key) (*Key, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapKey(m, c, res)
}

func unmarshalMapKey(m map[string]interface{}, c *Client, res *Key) (*Key, error) {

	flattened := flattenKey(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandKey expands Key into a JSON request object.
func expandKey(c *Client, f *Key) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v, err := expandKeyRestrictions(c, f.Restrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding Restrictions into restrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["restrictions"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenKey flattens Key from a JSON request object into the
// Key type.
func flattenKey(c *Client, i interface{}, res *Key) *Key {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Key{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.KeyString = dcl.FlattenString(m["keyString"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Restrictions = flattenKeyRestrictions(c, m["restrictions"], res)
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandKeyRestrictionsMap expands the contents of KeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsMap(c *Client, f map[string]KeyRestrictions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsSlice expands the contents of KeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsSlice(c *Client, f []KeyRestrictions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsMap flattens the contents of KeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictions{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictions{}
	}

	items := make(map[string]KeyRestrictions)
	for k, item := range a {
		items[k] = *flattenKeyRestrictions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsSlice flattens the contents of KeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsSlice(c *Client, i interface{}, res *Key) []KeyRestrictions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictions{}
	}

	if len(a) == 0 {
		return []KeyRestrictions{}
	}

	items := make([]KeyRestrictions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictions expands an instance of KeyRestrictions into a JSON
// request object.
func expandKeyRestrictions(c *Client, f *KeyRestrictions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandKeyRestrictionsBrowserKeyRestrictions(c, f.BrowserKeyRestrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding BrowserKeyRestrictions into browserKeyRestrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["browserKeyRestrictions"] = v
	}
	if v, err := expandKeyRestrictionsServerKeyRestrictions(c, f.ServerKeyRestrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding ServerKeyRestrictions into serverKeyRestrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["serverKeyRestrictions"] = v
	}
	if v, err := expandKeyRestrictionsAndroidKeyRestrictions(c, f.AndroidKeyRestrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding AndroidKeyRestrictions into androidKeyRestrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["androidKeyRestrictions"] = v
	}
	if v, err := expandKeyRestrictionsIosKeyRestrictions(c, f.IosKeyRestrictions, res); err != nil {
		return nil, fmt.Errorf("error expanding IosKeyRestrictions into iosKeyRestrictions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["iosKeyRestrictions"] = v
	}
	if v, err := expandKeyRestrictionsApiTargetsSlice(c, f.ApiTargets, res); err != nil {
		return nil, fmt.Errorf("error expanding ApiTargets into apiTargets: %w", err)
	} else if v != nil {
		m["apiTargets"] = v
	}

	return m, nil
}

// flattenKeyRestrictions flattens an instance of KeyRestrictions from a JSON
// response object.
func flattenKeyRestrictions(c *Client, i interface{}, res *Key) *KeyRestrictions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictions
	}
	r.BrowserKeyRestrictions = flattenKeyRestrictionsBrowserKeyRestrictions(c, m["browserKeyRestrictions"], res)
	r.ServerKeyRestrictions = flattenKeyRestrictionsServerKeyRestrictions(c, m["serverKeyRestrictions"], res)
	r.AndroidKeyRestrictions = flattenKeyRestrictionsAndroidKeyRestrictions(c, m["androidKeyRestrictions"], res)
	r.IosKeyRestrictions = flattenKeyRestrictionsIosKeyRestrictions(c, m["iosKeyRestrictions"], res)
	r.ApiTargets = flattenKeyRestrictionsApiTargetsSlice(c, m["apiTargets"], res)

	return r
}

// expandKeyRestrictionsBrowserKeyRestrictionsMap expands the contents of KeyRestrictionsBrowserKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsBrowserKeyRestrictionsMap(c *Client, f map[string]KeyRestrictionsBrowserKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsBrowserKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsBrowserKeyRestrictionsSlice expands the contents of KeyRestrictionsBrowserKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsBrowserKeyRestrictionsSlice(c *Client, f []KeyRestrictionsBrowserKeyRestrictions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsBrowserKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsBrowserKeyRestrictionsMap flattens the contents of KeyRestrictionsBrowserKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsBrowserKeyRestrictionsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsBrowserKeyRestrictions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsBrowserKeyRestrictions{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsBrowserKeyRestrictions{}
	}

	items := make(map[string]KeyRestrictionsBrowserKeyRestrictions)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsBrowserKeyRestrictions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsBrowserKeyRestrictionsSlice flattens the contents of KeyRestrictionsBrowserKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsBrowserKeyRestrictionsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsBrowserKeyRestrictions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsBrowserKeyRestrictions{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsBrowserKeyRestrictions{}
	}

	items := make([]KeyRestrictionsBrowserKeyRestrictions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsBrowserKeyRestrictions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsBrowserKeyRestrictions expands an instance of KeyRestrictionsBrowserKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsBrowserKeyRestrictions(c *Client, f *KeyRestrictionsBrowserKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowedReferrers; v != nil {
		m["allowedReferrers"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsBrowserKeyRestrictions flattens an instance of KeyRestrictionsBrowserKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsBrowserKeyRestrictions(c *Client, i interface{}, res *Key) *KeyRestrictionsBrowserKeyRestrictions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsBrowserKeyRestrictions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsBrowserKeyRestrictions
	}
	r.AllowedReferrers = dcl.FlattenStringSlice(m["allowedReferrers"])

	return r
}

// expandKeyRestrictionsServerKeyRestrictionsMap expands the contents of KeyRestrictionsServerKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsServerKeyRestrictionsMap(c *Client, f map[string]KeyRestrictionsServerKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsServerKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsServerKeyRestrictionsSlice expands the contents of KeyRestrictionsServerKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsServerKeyRestrictionsSlice(c *Client, f []KeyRestrictionsServerKeyRestrictions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsServerKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsServerKeyRestrictionsMap flattens the contents of KeyRestrictionsServerKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsServerKeyRestrictionsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsServerKeyRestrictions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsServerKeyRestrictions{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsServerKeyRestrictions{}
	}

	items := make(map[string]KeyRestrictionsServerKeyRestrictions)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsServerKeyRestrictions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsServerKeyRestrictionsSlice flattens the contents of KeyRestrictionsServerKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsServerKeyRestrictionsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsServerKeyRestrictions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsServerKeyRestrictions{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsServerKeyRestrictions{}
	}

	items := make([]KeyRestrictionsServerKeyRestrictions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsServerKeyRestrictions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsServerKeyRestrictions expands an instance of KeyRestrictionsServerKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsServerKeyRestrictions(c *Client, f *KeyRestrictionsServerKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowedIps; v != nil {
		m["allowedIps"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsServerKeyRestrictions flattens an instance of KeyRestrictionsServerKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsServerKeyRestrictions(c *Client, i interface{}, res *Key) *KeyRestrictionsServerKeyRestrictions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsServerKeyRestrictions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsServerKeyRestrictions
	}
	r.AllowedIps = dcl.FlattenStringSlice(m["allowedIps"])

	return r
}

// expandKeyRestrictionsAndroidKeyRestrictionsMap expands the contents of KeyRestrictionsAndroidKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictionsMap(c *Client, f map[string]KeyRestrictionsAndroidKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsAndroidKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsAndroidKeyRestrictionsSlice expands the contents of KeyRestrictionsAndroidKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictionsSlice(c *Client, f []KeyRestrictionsAndroidKeyRestrictions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsAndroidKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsAndroidKeyRestrictionsMap flattens the contents of KeyRestrictionsAndroidKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictionsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsAndroidKeyRestrictions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsAndroidKeyRestrictions{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsAndroidKeyRestrictions{}
	}

	items := make(map[string]KeyRestrictionsAndroidKeyRestrictions)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsAndroidKeyRestrictions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsAndroidKeyRestrictionsSlice flattens the contents of KeyRestrictionsAndroidKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictionsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsAndroidKeyRestrictions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsAndroidKeyRestrictions{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsAndroidKeyRestrictions{}
	}

	items := make([]KeyRestrictionsAndroidKeyRestrictions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsAndroidKeyRestrictions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsAndroidKeyRestrictions expands an instance of KeyRestrictionsAndroidKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictions(c *Client, f *KeyRestrictionsAndroidKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c, f.AllowedApplications, res); err != nil {
		return nil, fmt.Errorf("error expanding AllowedApplications into allowedApplications: %w", err)
	} else if v != nil {
		m["allowedApplications"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsAndroidKeyRestrictions flattens an instance of KeyRestrictionsAndroidKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictions(c *Client, i interface{}, res *Key) *KeyRestrictionsAndroidKeyRestrictions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsAndroidKeyRestrictions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsAndroidKeyRestrictions
	}
	r.AllowedApplications = flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c, m["allowedApplications"], res)

	return r
}

// expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsMap expands the contents of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsMap(c *Client, f map[string]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice expands the contents of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c *Client, f []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsMap flattens the contents of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}
	}

	items := make(map[string]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice flattens the contents of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}
	}

	items := make([]KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplications expands an instance of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications into a JSON
// request object.
func expandKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c *Client, f *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Sha1Fingerprint; !dcl.IsEmptyValueIndirect(v) {
		m["sha1Fingerprint"] = v
	}
	if v := f.PackageName; !dcl.IsEmptyValueIndirect(v) {
		m["packageName"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplications flattens an instance of KeyRestrictionsAndroidKeyRestrictionsAllowedApplications from a JSON
// response object.
func flattenKeyRestrictionsAndroidKeyRestrictionsAllowedApplications(c *Client, i interface{}, res *Key) *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	}
	r.Sha1Fingerprint = dcl.FlattenString(m["sha1Fingerprint"])
	r.PackageName = dcl.FlattenString(m["packageName"])

	return r
}

// expandKeyRestrictionsIosKeyRestrictionsMap expands the contents of KeyRestrictionsIosKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsIosKeyRestrictionsMap(c *Client, f map[string]KeyRestrictionsIosKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsIosKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsIosKeyRestrictionsSlice expands the contents of KeyRestrictionsIosKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsIosKeyRestrictionsSlice(c *Client, f []KeyRestrictionsIosKeyRestrictions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsIosKeyRestrictions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsIosKeyRestrictionsMap flattens the contents of KeyRestrictionsIosKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsIosKeyRestrictionsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsIosKeyRestrictions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsIosKeyRestrictions{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsIosKeyRestrictions{}
	}

	items := make(map[string]KeyRestrictionsIosKeyRestrictions)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsIosKeyRestrictions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsIosKeyRestrictionsSlice flattens the contents of KeyRestrictionsIosKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsIosKeyRestrictionsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsIosKeyRestrictions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsIosKeyRestrictions{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsIosKeyRestrictions{}
	}

	items := make([]KeyRestrictionsIosKeyRestrictions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsIosKeyRestrictions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsIosKeyRestrictions expands an instance of KeyRestrictionsIosKeyRestrictions into a JSON
// request object.
func expandKeyRestrictionsIosKeyRestrictions(c *Client, f *KeyRestrictionsIosKeyRestrictions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowedBundleIds; v != nil {
		m["allowedBundleIds"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsIosKeyRestrictions flattens an instance of KeyRestrictionsIosKeyRestrictions from a JSON
// response object.
func flattenKeyRestrictionsIosKeyRestrictions(c *Client, i interface{}, res *Key) *KeyRestrictionsIosKeyRestrictions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsIosKeyRestrictions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsIosKeyRestrictions
	}
	r.AllowedBundleIds = dcl.FlattenStringSlice(m["allowedBundleIds"])

	return r
}

// expandKeyRestrictionsApiTargetsMap expands the contents of KeyRestrictionsApiTargets into a JSON
// request object.
func expandKeyRestrictionsApiTargetsMap(c *Client, f map[string]KeyRestrictionsApiTargets, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyRestrictionsApiTargets(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyRestrictionsApiTargetsSlice expands the contents of KeyRestrictionsApiTargets into a JSON
// request object.
func expandKeyRestrictionsApiTargetsSlice(c *Client, f []KeyRestrictionsApiTargets, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyRestrictionsApiTargets(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyRestrictionsApiTargetsMap flattens the contents of KeyRestrictionsApiTargets from a JSON
// response object.
func flattenKeyRestrictionsApiTargetsMap(c *Client, i interface{}, res *Key) map[string]KeyRestrictionsApiTargets {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyRestrictionsApiTargets{}
	}

	if len(a) == 0 {
		return map[string]KeyRestrictionsApiTargets{}
	}

	items := make(map[string]KeyRestrictionsApiTargets)
	for k, item := range a {
		items[k] = *flattenKeyRestrictionsApiTargets(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyRestrictionsApiTargetsSlice flattens the contents of KeyRestrictionsApiTargets from a JSON
// response object.
func flattenKeyRestrictionsApiTargetsSlice(c *Client, i interface{}, res *Key) []KeyRestrictionsApiTargets {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyRestrictionsApiTargets{}
	}

	if len(a) == 0 {
		return []KeyRestrictionsApiTargets{}
	}

	items := make([]KeyRestrictionsApiTargets, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyRestrictionsApiTargets(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyRestrictionsApiTargets expands an instance of KeyRestrictionsApiTargets into a JSON
// request object.
func expandKeyRestrictionsApiTargets(c *Client, f *KeyRestrictionsApiTargets, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}
	if v := f.Methods; v != nil {
		m["methods"] = v
	}

	return m, nil
}

// flattenKeyRestrictionsApiTargets flattens an instance of KeyRestrictionsApiTargets from a JSON
// response object.
func flattenKeyRestrictionsApiTargets(c *Client, i interface{}, res *Key) *KeyRestrictionsApiTargets {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyRestrictionsApiTargets{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyRestrictionsApiTargets
	}
	r.Service = dcl.FlattenString(m["service"])
	r.Methods = dcl.FlattenStringSlice(m["methods"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Key) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalKey(b, c, r)
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

type keyDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         keyApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToKeyDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]keyDiff, error) {
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
	var diffs []keyDiff
	// For each operation name, create a keyDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := keyDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToKeyApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToKeyApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (keyApiOperation, error) {
	switch opName {

	case "updateKeyUpdateKeyOperation":
		return &updateKeyUpdateKeyOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractKeyFields(r *Key) error {
	vRestrictions := r.Restrictions
	if vRestrictions == nil {
		// note: explicitly not the empty object.
		vRestrictions = &KeyRestrictions{}
	}
	if err := extractKeyRestrictionsFields(r, vRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRestrictions) {
		r.Restrictions = vRestrictions
	}
	return nil
}
func extractKeyRestrictionsFields(r *Key, o *KeyRestrictions) error {
	vBrowserKeyRestrictions := o.BrowserKeyRestrictions
	if vBrowserKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vBrowserKeyRestrictions = &KeyRestrictionsBrowserKeyRestrictions{}
	}
	if err := extractKeyRestrictionsBrowserKeyRestrictionsFields(r, vBrowserKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBrowserKeyRestrictions) {
		o.BrowserKeyRestrictions = vBrowserKeyRestrictions
	}
	vServerKeyRestrictions := o.ServerKeyRestrictions
	if vServerKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vServerKeyRestrictions = &KeyRestrictionsServerKeyRestrictions{}
	}
	if err := extractKeyRestrictionsServerKeyRestrictionsFields(r, vServerKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServerKeyRestrictions) {
		o.ServerKeyRestrictions = vServerKeyRestrictions
	}
	vAndroidKeyRestrictions := o.AndroidKeyRestrictions
	if vAndroidKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vAndroidKeyRestrictions = &KeyRestrictionsAndroidKeyRestrictions{}
	}
	if err := extractKeyRestrictionsAndroidKeyRestrictionsFields(r, vAndroidKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAndroidKeyRestrictions) {
		o.AndroidKeyRestrictions = vAndroidKeyRestrictions
	}
	vIosKeyRestrictions := o.IosKeyRestrictions
	if vIosKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vIosKeyRestrictions = &KeyRestrictionsIosKeyRestrictions{}
	}
	if err := extractKeyRestrictionsIosKeyRestrictionsFields(r, vIosKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIosKeyRestrictions) {
		o.IosKeyRestrictions = vIosKeyRestrictions
	}
	return nil
}
func extractKeyRestrictionsBrowserKeyRestrictionsFields(r *Key, o *KeyRestrictionsBrowserKeyRestrictions) error {
	return nil
}
func extractKeyRestrictionsServerKeyRestrictionsFields(r *Key, o *KeyRestrictionsServerKeyRestrictions) error {
	return nil
}
func extractKeyRestrictionsAndroidKeyRestrictionsFields(r *Key, o *KeyRestrictionsAndroidKeyRestrictions) error {
	return nil
}
func extractKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsFields(r *Key, o *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) error {
	return nil
}
func extractKeyRestrictionsIosKeyRestrictionsFields(r *Key, o *KeyRestrictionsIosKeyRestrictions) error {
	return nil
}
func extractKeyRestrictionsApiTargetsFields(r *Key, o *KeyRestrictionsApiTargets) error {
	return nil
}

func postReadExtractKeyFields(r *Key) error {
	vRestrictions := r.Restrictions
	if vRestrictions == nil {
		// note: explicitly not the empty object.
		vRestrictions = &KeyRestrictions{}
	}
	if err := postReadExtractKeyRestrictionsFields(r, vRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRestrictions) {
		r.Restrictions = vRestrictions
	}
	return nil
}
func postReadExtractKeyRestrictionsFields(r *Key, o *KeyRestrictions) error {
	vBrowserKeyRestrictions := o.BrowserKeyRestrictions
	if vBrowserKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vBrowserKeyRestrictions = &KeyRestrictionsBrowserKeyRestrictions{}
	}
	if err := extractKeyRestrictionsBrowserKeyRestrictionsFields(r, vBrowserKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBrowserKeyRestrictions) {
		o.BrowserKeyRestrictions = vBrowserKeyRestrictions
	}
	vServerKeyRestrictions := o.ServerKeyRestrictions
	if vServerKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vServerKeyRestrictions = &KeyRestrictionsServerKeyRestrictions{}
	}
	if err := extractKeyRestrictionsServerKeyRestrictionsFields(r, vServerKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServerKeyRestrictions) {
		o.ServerKeyRestrictions = vServerKeyRestrictions
	}
	vAndroidKeyRestrictions := o.AndroidKeyRestrictions
	if vAndroidKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vAndroidKeyRestrictions = &KeyRestrictionsAndroidKeyRestrictions{}
	}
	if err := extractKeyRestrictionsAndroidKeyRestrictionsFields(r, vAndroidKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAndroidKeyRestrictions) {
		o.AndroidKeyRestrictions = vAndroidKeyRestrictions
	}
	vIosKeyRestrictions := o.IosKeyRestrictions
	if vIosKeyRestrictions == nil {
		// note: explicitly not the empty object.
		vIosKeyRestrictions = &KeyRestrictionsIosKeyRestrictions{}
	}
	if err := extractKeyRestrictionsIosKeyRestrictionsFields(r, vIosKeyRestrictions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIosKeyRestrictions) {
		o.IosKeyRestrictions = vIosKeyRestrictions
	}
	return nil
}
func postReadExtractKeyRestrictionsBrowserKeyRestrictionsFields(r *Key, o *KeyRestrictionsBrowserKeyRestrictions) error {
	return nil
}
func postReadExtractKeyRestrictionsServerKeyRestrictionsFields(r *Key, o *KeyRestrictionsServerKeyRestrictions) error {
	return nil
}
func postReadExtractKeyRestrictionsAndroidKeyRestrictionsFields(r *Key, o *KeyRestrictionsAndroidKeyRestrictions) error {
	return nil
}
func postReadExtractKeyRestrictionsAndroidKeyRestrictionsAllowedApplicationsFields(r *Key, o *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) error {
	return nil
}
func postReadExtractKeyRestrictionsIosKeyRestrictionsFields(r *Key, o *KeyRestrictionsIosKeyRestrictions) error {
	return nil
}
func postReadExtractKeyRestrictionsApiTargetsFields(r *Key, o *KeyRestrictionsApiTargets) error {
	return nil
}
