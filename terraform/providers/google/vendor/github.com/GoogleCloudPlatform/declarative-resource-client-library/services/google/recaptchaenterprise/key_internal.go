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
package recaptchaenterprise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Key) validate() error {

	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"WebSettings", "AndroidSettings", "IosSettings"}, r.WebSettings, r.AndroidSettings, r.IosSettings); err != nil {
		return err
	}
	if err := dcl.Required(r, "displayName"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.WebSettings) {
		if err := r.WebSettings.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AndroidSettings) {
		if err := r.AndroidSettings.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.IosSettings) {
		if err := r.IosSettings.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.TestingOptions) {
		if err := r.TestingOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *KeyWebSettings) validate() error {
	if err := dcl.Required(r, "integrationType"); err != nil {
		return err
	}
	return nil
}
func (r *KeyAndroidSettings) validate() error {
	return nil
}
func (r *KeyIosSettings) validate() error {
	return nil
}
func (r *KeyTestingOptions) validate() error {
	return nil
}
func (r *Key) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://recaptchaenterprise.googleapis.com/v1/", params)
}

func (r *Key) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/keys/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Key) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/keys", nr.basePath(), userBasePath, params), nil

}

func (r *Key) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/keys", nr.basePath(), userBasePath, params), nil

}

func (r *Key) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/keys/{{name}}", nr.basePath(), userBasePath, params), nil
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
	if v, err := expandKeyWebSettings(c, f.WebSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding WebSettings into webSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["webSettings"] = v
	}
	if v, err := expandKeyAndroidSettings(c, f.AndroidSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding AndroidSettings into androidSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["androidSettings"] = v
	}
	if v, err := expandKeyIosSettings(c, f.IosSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding IosSettings into iosSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["iosSettings"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
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
	mask := dcl.UpdateMask(op.FieldDiffs)
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
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete Key: %w", err)
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
	if r.Name != nil {
		// Allowing creation to continue with Name set could result in a Key with the wrong Name.
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

	if _, err := c.GetKey(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getKeyRaw(ctx context.Context, r *Key) ([]byte, error) {

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

	if fetchState.Name == nil {
		// We cannot perform a get because of lack of information. We have to assume
		// that this is being created for the first time.
		desired, err := canonicalizeKeyDesiredState(rawDesired, nil)
		return nil, desired, nil, err
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

	if !dcl.IsZeroValue(rawInitial.WebSettings) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.AndroidSettings, rawInitial.IosSettings) {
			rawInitial.WebSettings = EmptyKeyWebSettings
		}
	}

	if !dcl.IsZeroValue(rawInitial.AndroidSettings) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.WebSettings, rawInitial.IosSettings) {
			rawInitial.AndroidSettings = EmptyKeyAndroidSettings
		}
	}

	if !dcl.IsZeroValue(rawInitial.IosSettings) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.WebSettings, rawInitial.AndroidSettings) {
			rawInitial.IosSettings = EmptyKeyIosSettings
		}
	}

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
		rawDesired.WebSettings = canonicalizeKeyWebSettings(rawDesired.WebSettings, nil, opts...)
		rawDesired.AndroidSettings = canonicalizeKeyAndroidSettings(rawDesired.AndroidSettings, nil, opts...)
		rawDesired.IosSettings = canonicalizeKeyIosSettings(rawDesired.IosSettings, nil, opts...)
		rawDesired.TestingOptions = canonicalizeKeyTestingOptions(rawDesired.TestingOptions, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Key{}
	if dcl.IsZeroValue(rawDesired.Name) || (dcl.IsEmptyValueIndirect(rawDesired.Name) && dcl.IsEmptyValueIndirect(rawInitial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	canonicalDesired.WebSettings = canonicalizeKeyWebSettings(rawDesired.WebSettings, rawInitial.WebSettings, opts...)
	canonicalDesired.AndroidSettings = canonicalizeKeyAndroidSettings(rawDesired.AndroidSettings, rawInitial.AndroidSettings, opts...)
	canonicalDesired.IosSettings = canonicalizeKeyIosSettings(rawDesired.IosSettings, rawInitial.IosSettings, opts...)
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	canonicalDesired.TestingOptions = canonicalizeKeyTestingOptions(rawDesired.TestingOptions, rawInitial.TestingOptions, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}

	if canonicalDesired.WebSettings != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.AndroidSettings, rawDesired.IosSettings) {
			canonicalDesired.WebSettings = EmptyKeyWebSettings
		}
	}

	if canonicalDesired.AndroidSettings != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.WebSettings, rawDesired.IosSettings) {
			canonicalDesired.AndroidSettings = EmptyKeyAndroidSettings
		}
	}

	if canonicalDesired.IosSettings != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.WebSettings, rawDesired.AndroidSettings) {
			canonicalDesired.IosSettings = EmptyKeyIosSettings
		}
	}

	return canonicalDesired, nil
}

func canonicalizeKeyNewState(c *Client, rawNew, rawDesired *Key) (*Key, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.WebSettings) && dcl.IsEmptyValueIndirect(rawDesired.WebSettings) {
		rawNew.WebSettings = rawDesired.WebSettings
	} else {
		rawNew.WebSettings = canonicalizeNewKeyWebSettings(c, rawDesired.WebSettings, rawNew.WebSettings)
	}

	if dcl.IsEmptyValueIndirect(rawNew.AndroidSettings) && dcl.IsEmptyValueIndirect(rawDesired.AndroidSettings) {
		rawNew.AndroidSettings = rawDesired.AndroidSettings
	} else {
		rawNew.AndroidSettings = canonicalizeNewKeyAndroidSettings(c, rawDesired.AndroidSettings, rawNew.AndroidSettings)
	}

	if dcl.IsEmptyValueIndirect(rawNew.IosSettings) && dcl.IsEmptyValueIndirect(rawDesired.IosSettings) {
		rawNew.IosSettings = rawDesired.IosSettings
	} else {
		rawNew.IosSettings = canonicalizeNewKeyIosSettings(c, rawDesired.IosSettings, rawNew.IosSettings)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.TestingOptions) && dcl.IsEmptyValueIndirect(rawDesired.TestingOptions) {
		rawNew.TestingOptions = rawDesired.TestingOptions
	} else {
		rawNew.TestingOptions = canonicalizeNewKeyTestingOptions(c, rawDesired.TestingOptions, rawNew.TestingOptions)
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeKeyWebSettings(des, initial *KeyWebSettings, opts ...dcl.ApplyOption) *KeyWebSettings {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyWebSettings{}

	if dcl.BoolCanonicalize(des.AllowAllDomains, initial.AllowAllDomains) || dcl.IsZeroValue(des.AllowAllDomains) {
		cDes.AllowAllDomains = initial.AllowAllDomains
	} else {
		cDes.AllowAllDomains = des.AllowAllDomains
	}
	if dcl.StringArrayCanonicalize(des.AllowedDomains, initial.AllowedDomains) {
		cDes.AllowedDomains = initial.AllowedDomains
	} else {
		cDes.AllowedDomains = des.AllowedDomains
	}
	if dcl.BoolCanonicalize(des.AllowAmpTraffic, initial.AllowAmpTraffic) || dcl.IsZeroValue(des.AllowAmpTraffic) {
		cDes.AllowAmpTraffic = initial.AllowAmpTraffic
	} else {
		cDes.AllowAmpTraffic = des.AllowAmpTraffic
	}
	if dcl.IsZeroValue(des.IntegrationType) || (dcl.IsEmptyValueIndirect(des.IntegrationType) && dcl.IsEmptyValueIndirect(initial.IntegrationType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.IntegrationType = initial.IntegrationType
	} else {
		cDes.IntegrationType = des.IntegrationType
	}
	if dcl.IsZeroValue(des.ChallengeSecurityPreference) || (dcl.IsEmptyValueIndirect(des.ChallengeSecurityPreference) && dcl.IsEmptyValueIndirect(initial.ChallengeSecurityPreference)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ChallengeSecurityPreference = initial.ChallengeSecurityPreference
	} else {
		cDes.ChallengeSecurityPreference = des.ChallengeSecurityPreference
	}

	return cDes
}

func canonicalizeKeyWebSettingsSlice(des, initial []KeyWebSettings, opts ...dcl.ApplyOption) []KeyWebSettings {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyWebSettings, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyWebSettings(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyWebSettings, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyWebSettings(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyWebSettings(c *Client, des, nw *KeyWebSettings) *KeyWebSettings {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyWebSettings while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AllowAllDomains, nw.AllowAllDomains) {
		nw.AllowAllDomains = des.AllowAllDomains
	}
	if dcl.StringArrayCanonicalize(des.AllowedDomains, nw.AllowedDomains) {
		nw.AllowedDomains = des.AllowedDomains
	}
	if dcl.BoolCanonicalize(des.AllowAmpTraffic, nw.AllowAmpTraffic) {
		nw.AllowAmpTraffic = des.AllowAmpTraffic
	}

	return nw
}

func canonicalizeNewKeyWebSettingsSet(c *Client, des, nw []KeyWebSettings) []KeyWebSettings {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyWebSettings
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyWebSettingsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyWebSettings(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyWebSettingsSlice(c *Client, des, nw []KeyWebSettings) []KeyWebSettings {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyWebSettings
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyWebSettings(c, &d, &n))
	}

	return items
}

func canonicalizeKeyAndroidSettings(des, initial *KeyAndroidSettings, opts ...dcl.ApplyOption) *KeyAndroidSettings {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyAndroidSettings{}

	if dcl.BoolCanonicalize(des.AllowAllPackageNames, initial.AllowAllPackageNames) || dcl.IsZeroValue(des.AllowAllPackageNames) {
		cDes.AllowAllPackageNames = initial.AllowAllPackageNames
	} else {
		cDes.AllowAllPackageNames = des.AllowAllPackageNames
	}
	if dcl.StringArrayCanonicalize(des.AllowedPackageNames, initial.AllowedPackageNames) {
		cDes.AllowedPackageNames = initial.AllowedPackageNames
	} else {
		cDes.AllowedPackageNames = des.AllowedPackageNames
	}

	return cDes
}

func canonicalizeKeyAndroidSettingsSlice(des, initial []KeyAndroidSettings, opts ...dcl.ApplyOption) []KeyAndroidSettings {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyAndroidSettings, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyAndroidSettings(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyAndroidSettings, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyAndroidSettings(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyAndroidSettings(c *Client, des, nw *KeyAndroidSettings) *KeyAndroidSettings {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyAndroidSettings while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AllowAllPackageNames, nw.AllowAllPackageNames) {
		nw.AllowAllPackageNames = des.AllowAllPackageNames
	}
	if dcl.StringArrayCanonicalize(des.AllowedPackageNames, nw.AllowedPackageNames) {
		nw.AllowedPackageNames = des.AllowedPackageNames
	}

	return nw
}

func canonicalizeNewKeyAndroidSettingsSet(c *Client, des, nw []KeyAndroidSettings) []KeyAndroidSettings {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyAndroidSettings
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyAndroidSettingsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyAndroidSettings(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyAndroidSettingsSlice(c *Client, des, nw []KeyAndroidSettings) []KeyAndroidSettings {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyAndroidSettings
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyAndroidSettings(c, &d, &n))
	}

	return items
}

func canonicalizeKeyIosSettings(des, initial *KeyIosSettings, opts ...dcl.ApplyOption) *KeyIosSettings {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyIosSettings{}

	if dcl.BoolCanonicalize(des.AllowAllBundleIds, initial.AllowAllBundleIds) || dcl.IsZeroValue(des.AllowAllBundleIds) {
		cDes.AllowAllBundleIds = initial.AllowAllBundleIds
	} else {
		cDes.AllowAllBundleIds = des.AllowAllBundleIds
	}
	if dcl.StringArrayCanonicalize(des.AllowedBundleIds, initial.AllowedBundleIds) {
		cDes.AllowedBundleIds = initial.AllowedBundleIds
	} else {
		cDes.AllowedBundleIds = des.AllowedBundleIds
	}

	return cDes
}

func canonicalizeKeyIosSettingsSlice(des, initial []KeyIosSettings, opts ...dcl.ApplyOption) []KeyIosSettings {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyIosSettings, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyIosSettings(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyIosSettings, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyIosSettings(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyIosSettings(c *Client, des, nw *KeyIosSettings) *KeyIosSettings {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyIosSettings while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AllowAllBundleIds, nw.AllowAllBundleIds) {
		nw.AllowAllBundleIds = des.AllowAllBundleIds
	}
	if dcl.StringArrayCanonicalize(des.AllowedBundleIds, nw.AllowedBundleIds) {
		nw.AllowedBundleIds = des.AllowedBundleIds
	}

	return nw
}

func canonicalizeNewKeyIosSettingsSet(c *Client, des, nw []KeyIosSettings) []KeyIosSettings {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyIosSettings
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyIosSettingsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyIosSettings(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyIosSettingsSlice(c *Client, des, nw []KeyIosSettings) []KeyIosSettings {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyIosSettings
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyIosSettings(c, &d, &n))
	}

	return items
}

func canonicalizeKeyTestingOptions(des, initial *KeyTestingOptions, opts ...dcl.ApplyOption) *KeyTestingOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &KeyTestingOptions{}

	if dcl.IsZeroValue(des.TestingScore) || (dcl.IsEmptyValueIndirect(des.TestingScore) && dcl.IsEmptyValueIndirect(initial.TestingScore)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.TestingScore = initial.TestingScore
	} else {
		cDes.TestingScore = des.TestingScore
	}
	if dcl.IsZeroValue(des.TestingChallenge) || (dcl.IsEmptyValueIndirect(des.TestingChallenge) && dcl.IsEmptyValueIndirect(initial.TestingChallenge)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.TestingChallenge = initial.TestingChallenge
	} else {
		cDes.TestingChallenge = des.TestingChallenge
	}

	return cDes
}

func canonicalizeKeyTestingOptionsSlice(des, initial []KeyTestingOptions, opts ...dcl.ApplyOption) []KeyTestingOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]KeyTestingOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeKeyTestingOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]KeyTestingOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeKeyTestingOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewKeyTestingOptions(c *Client, des, nw *KeyTestingOptions) *KeyTestingOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for KeyTestingOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewKeyTestingOptionsSet(c *Client, des, nw []KeyTestingOptions) []KeyTestingOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []KeyTestingOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareKeyTestingOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewKeyTestingOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewKeyTestingOptionsSlice(c *Client, des, nw []KeyTestingOptions) []KeyTestingOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []KeyTestingOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewKeyTestingOptions(c, &d, &n))
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
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.WebSettings, actual.WebSettings, dcl.DiffInfo{ObjectFunction: compareKeyWebSettingsNewStyle, EmptyObject: EmptyKeyWebSettings, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WebSettings")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AndroidSettings, actual.AndroidSettings, dcl.DiffInfo{ObjectFunction: compareKeyAndroidSettingsNewStyle, EmptyObject: EmptyKeyAndroidSettings, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AndroidSettings")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IosSettings, actual.IosSettings, dcl.DiffInfo{ObjectFunction: compareKeyIosSettingsNewStyle, EmptyObject: EmptyKeyIosSettings, OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("IosSettings")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.TestingOptions, actual.TestingOptions, dcl.DiffInfo{ObjectFunction: compareKeyTestingOptionsNewStyle, EmptyObject: EmptyKeyTestingOptions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TestingOptions")); len(ds) != 0 || err != nil {
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
func compareKeyWebSettingsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyWebSettings)
	if !ok {
		desiredNotPointer, ok := d.(KeyWebSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyWebSettings or *KeyWebSettings", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyWebSettings)
	if !ok {
		actualNotPointer, ok := a.(KeyWebSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyWebSettings", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowAllDomains, actual.AllowAllDomains, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowAllDomains")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedDomains, actual.AllowedDomains, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedDomains")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowAmpTraffic, actual.AllowAmpTraffic, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowAmpTraffic")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IntegrationType, actual.IntegrationType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IntegrationType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ChallengeSecurityPreference, actual.ChallengeSecurityPreference, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("ChallengeSecurityPreference")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyAndroidSettingsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyAndroidSettings)
	if !ok {
		desiredNotPointer, ok := d.(KeyAndroidSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyAndroidSettings or *KeyAndroidSettings", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyAndroidSettings)
	if !ok {
		actualNotPointer, ok := a.(KeyAndroidSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyAndroidSettings", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowAllPackageNames, actual.AllowAllPackageNames, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowAllPackageNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedPackageNames, actual.AllowedPackageNames, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedPackageNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyIosSettingsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyIosSettings)
	if !ok {
		desiredNotPointer, ok := d.(KeyIosSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyIosSettings or *KeyIosSettings", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyIosSettings)
	if !ok {
		actualNotPointer, ok := a.(KeyIosSettings)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyIosSettings", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowAllBundleIds, actual.AllowAllBundleIds, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowAllBundleIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedBundleIds, actual.AllowedBundleIds, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateKeyUpdateKeyOperation")}, fn.AddNest("AllowedBundleIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareKeyTestingOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*KeyTestingOptions)
	if !ok {
		desiredNotPointer, ok := d.(KeyTestingOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyTestingOptions or *KeyTestingOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*KeyTestingOptions)
	if !ok {
		actualNotPointer, ok := a.(KeyTestingOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a KeyTestingOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.TestingScore, actual.TestingScore, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TestingScore")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TestingChallenge, actual.TestingChallenge, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TestingChallenge")); len(ds) != 0 || err != nil {
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
		return dcl.URL("projects/{{project}}/keys/{{name}}", nr.basePath(), userBasePath, fields), nil

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
	if v, err := dcl.DeriveField("projects/%s/keys/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v, err := expandKeyWebSettings(c, f.WebSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding WebSettings into webSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["webSettings"] = v
	}
	if v, err := expandKeyAndroidSettings(c, f.AndroidSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding AndroidSettings into androidSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["androidSettings"] = v
	}
	if v, err := expandKeyIosSettings(c, f.IosSettings, res); err != nil {
		return nil, fmt.Errorf("error expanding IosSettings into iosSettings: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["iosSettings"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := expandKeyTestingOptions(c, f.TestingOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding TestingOptions into testingOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["testingOptions"] = v
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
	resultRes.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.WebSettings = flattenKeyWebSettings(c, m["webSettings"], res)
	resultRes.AndroidSettings = flattenKeyAndroidSettings(c, m["androidSettings"], res)
	resultRes.IosSettings = flattenKeyIosSettings(c, m["iosSettings"], res)
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.TestingOptions = flattenKeyTestingOptions(c, m["testingOptions"], res)
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandKeyWebSettingsMap expands the contents of KeyWebSettings into a JSON
// request object.
func expandKeyWebSettingsMap(c *Client, f map[string]KeyWebSettings, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyWebSettings(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyWebSettingsSlice expands the contents of KeyWebSettings into a JSON
// request object.
func expandKeyWebSettingsSlice(c *Client, f []KeyWebSettings, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyWebSettings(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyWebSettingsMap flattens the contents of KeyWebSettings from a JSON
// response object.
func flattenKeyWebSettingsMap(c *Client, i interface{}, res *Key) map[string]KeyWebSettings {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyWebSettings{}
	}

	if len(a) == 0 {
		return map[string]KeyWebSettings{}
	}

	items := make(map[string]KeyWebSettings)
	for k, item := range a {
		items[k] = *flattenKeyWebSettings(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyWebSettingsSlice flattens the contents of KeyWebSettings from a JSON
// response object.
func flattenKeyWebSettingsSlice(c *Client, i interface{}, res *Key) []KeyWebSettings {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyWebSettings{}
	}

	if len(a) == 0 {
		return []KeyWebSettings{}
	}

	items := make([]KeyWebSettings, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyWebSettings(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyWebSettings expands an instance of KeyWebSettings into a JSON
// request object.
func expandKeyWebSettings(c *Client, f *KeyWebSettings, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowAllDomains; !dcl.IsEmptyValueIndirect(v) {
		m["allowAllDomains"] = v
	}
	if v := f.AllowedDomains; v != nil {
		m["allowedDomains"] = v
	}
	if v := f.AllowAmpTraffic; !dcl.IsEmptyValueIndirect(v) {
		m["allowAmpTraffic"] = v
	}
	if v := f.IntegrationType; !dcl.IsEmptyValueIndirect(v) {
		m["integrationType"] = v
	}
	if v := f.ChallengeSecurityPreference; !dcl.IsEmptyValueIndirect(v) {
		m["challengeSecurityPreference"] = v
	}

	return m, nil
}

// flattenKeyWebSettings flattens an instance of KeyWebSettings from a JSON
// response object.
func flattenKeyWebSettings(c *Client, i interface{}, res *Key) *KeyWebSettings {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyWebSettings{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyWebSettings
	}
	r.AllowAllDomains = dcl.FlattenBool(m["allowAllDomains"])
	r.AllowedDomains = dcl.FlattenStringSlice(m["allowedDomains"])
	r.AllowAmpTraffic = dcl.FlattenBool(m["allowAmpTraffic"])
	r.IntegrationType = flattenKeyWebSettingsIntegrationTypeEnum(m["integrationType"])
	r.ChallengeSecurityPreference = flattenKeyWebSettingsChallengeSecurityPreferenceEnum(m["challengeSecurityPreference"])

	return r
}

// expandKeyAndroidSettingsMap expands the contents of KeyAndroidSettings into a JSON
// request object.
func expandKeyAndroidSettingsMap(c *Client, f map[string]KeyAndroidSettings, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyAndroidSettings(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyAndroidSettingsSlice expands the contents of KeyAndroidSettings into a JSON
// request object.
func expandKeyAndroidSettingsSlice(c *Client, f []KeyAndroidSettings, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyAndroidSettings(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyAndroidSettingsMap flattens the contents of KeyAndroidSettings from a JSON
// response object.
func flattenKeyAndroidSettingsMap(c *Client, i interface{}, res *Key) map[string]KeyAndroidSettings {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyAndroidSettings{}
	}

	if len(a) == 0 {
		return map[string]KeyAndroidSettings{}
	}

	items := make(map[string]KeyAndroidSettings)
	for k, item := range a {
		items[k] = *flattenKeyAndroidSettings(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyAndroidSettingsSlice flattens the contents of KeyAndroidSettings from a JSON
// response object.
func flattenKeyAndroidSettingsSlice(c *Client, i interface{}, res *Key) []KeyAndroidSettings {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyAndroidSettings{}
	}

	if len(a) == 0 {
		return []KeyAndroidSettings{}
	}

	items := make([]KeyAndroidSettings, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyAndroidSettings(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyAndroidSettings expands an instance of KeyAndroidSettings into a JSON
// request object.
func expandKeyAndroidSettings(c *Client, f *KeyAndroidSettings, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowAllPackageNames; !dcl.IsEmptyValueIndirect(v) {
		m["allowAllPackageNames"] = v
	}
	if v := f.AllowedPackageNames; v != nil {
		m["allowedPackageNames"] = v
	}

	return m, nil
}

// flattenKeyAndroidSettings flattens an instance of KeyAndroidSettings from a JSON
// response object.
func flattenKeyAndroidSettings(c *Client, i interface{}, res *Key) *KeyAndroidSettings {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyAndroidSettings{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyAndroidSettings
	}
	r.AllowAllPackageNames = dcl.FlattenBool(m["allowAllPackageNames"])
	r.AllowedPackageNames = dcl.FlattenStringSlice(m["allowedPackageNames"])

	return r
}

// expandKeyIosSettingsMap expands the contents of KeyIosSettings into a JSON
// request object.
func expandKeyIosSettingsMap(c *Client, f map[string]KeyIosSettings, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyIosSettings(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyIosSettingsSlice expands the contents of KeyIosSettings into a JSON
// request object.
func expandKeyIosSettingsSlice(c *Client, f []KeyIosSettings, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyIosSettings(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyIosSettingsMap flattens the contents of KeyIosSettings from a JSON
// response object.
func flattenKeyIosSettingsMap(c *Client, i interface{}, res *Key) map[string]KeyIosSettings {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyIosSettings{}
	}

	if len(a) == 0 {
		return map[string]KeyIosSettings{}
	}

	items := make(map[string]KeyIosSettings)
	for k, item := range a {
		items[k] = *flattenKeyIosSettings(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyIosSettingsSlice flattens the contents of KeyIosSettings from a JSON
// response object.
func flattenKeyIosSettingsSlice(c *Client, i interface{}, res *Key) []KeyIosSettings {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyIosSettings{}
	}

	if len(a) == 0 {
		return []KeyIosSettings{}
	}

	items := make([]KeyIosSettings, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyIosSettings(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyIosSettings expands an instance of KeyIosSettings into a JSON
// request object.
func expandKeyIosSettings(c *Client, f *KeyIosSettings, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowAllBundleIds; !dcl.IsEmptyValueIndirect(v) {
		m["allowAllBundleIds"] = v
	}
	if v := f.AllowedBundleIds; v != nil {
		m["allowedBundleIds"] = v
	}

	return m, nil
}

// flattenKeyIosSettings flattens an instance of KeyIosSettings from a JSON
// response object.
func flattenKeyIosSettings(c *Client, i interface{}, res *Key) *KeyIosSettings {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyIosSettings{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyIosSettings
	}
	r.AllowAllBundleIds = dcl.FlattenBool(m["allowAllBundleIds"])
	r.AllowedBundleIds = dcl.FlattenStringSlice(m["allowedBundleIds"])

	return r
}

// expandKeyTestingOptionsMap expands the contents of KeyTestingOptions into a JSON
// request object.
func expandKeyTestingOptionsMap(c *Client, f map[string]KeyTestingOptions, res *Key) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandKeyTestingOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandKeyTestingOptionsSlice expands the contents of KeyTestingOptions into a JSON
// request object.
func expandKeyTestingOptionsSlice(c *Client, f []KeyTestingOptions, res *Key) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandKeyTestingOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenKeyTestingOptionsMap flattens the contents of KeyTestingOptions from a JSON
// response object.
func flattenKeyTestingOptionsMap(c *Client, i interface{}, res *Key) map[string]KeyTestingOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyTestingOptions{}
	}

	if len(a) == 0 {
		return map[string]KeyTestingOptions{}
	}

	items := make(map[string]KeyTestingOptions)
	for k, item := range a {
		items[k] = *flattenKeyTestingOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenKeyTestingOptionsSlice flattens the contents of KeyTestingOptions from a JSON
// response object.
func flattenKeyTestingOptionsSlice(c *Client, i interface{}, res *Key) []KeyTestingOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyTestingOptions{}
	}

	if len(a) == 0 {
		return []KeyTestingOptions{}
	}

	items := make([]KeyTestingOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyTestingOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandKeyTestingOptions expands an instance of KeyTestingOptions into a JSON
// request object.
func expandKeyTestingOptions(c *Client, f *KeyTestingOptions, res *Key) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.TestingScore; !dcl.IsEmptyValueIndirect(v) {
		m["testingScore"] = v
	}
	if v := f.TestingChallenge; !dcl.IsEmptyValueIndirect(v) {
		m["testingChallenge"] = v
	}

	return m, nil
}

// flattenKeyTestingOptions flattens an instance of KeyTestingOptions from a JSON
// response object.
func flattenKeyTestingOptions(c *Client, i interface{}, res *Key) *KeyTestingOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &KeyTestingOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyKeyTestingOptions
	}
	r.TestingScore = dcl.FlattenDouble(m["testingScore"])
	r.TestingChallenge = flattenKeyTestingOptionsTestingChallengeEnum(m["testingChallenge"])

	return r
}

// flattenKeyWebSettingsIntegrationTypeEnumMap flattens the contents of KeyWebSettingsIntegrationTypeEnum from a JSON
// response object.
func flattenKeyWebSettingsIntegrationTypeEnumMap(c *Client, i interface{}, res *Key) map[string]KeyWebSettingsIntegrationTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyWebSettingsIntegrationTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]KeyWebSettingsIntegrationTypeEnum{}
	}

	items := make(map[string]KeyWebSettingsIntegrationTypeEnum)
	for k, item := range a {
		items[k] = *flattenKeyWebSettingsIntegrationTypeEnum(item.(interface{}))
	}

	return items
}

// flattenKeyWebSettingsIntegrationTypeEnumSlice flattens the contents of KeyWebSettingsIntegrationTypeEnum from a JSON
// response object.
func flattenKeyWebSettingsIntegrationTypeEnumSlice(c *Client, i interface{}, res *Key) []KeyWebSettingsIntegrationTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyWebSettingsIntegrationTypeEnum{}
	}

	if len(a) == 0 {
		return []KeyWebSettingsIntegrationTypeEnum{}
	}

	items := make([]KeyWebSettingsIntegrationTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyWebSettingsIntegrationTypeEnum(item.(interface{})))
	}

	return items
}

// flattenKeyWebSettingsIntegrationTypeEnum asserts that an interface is a string, and returns a
// pointer to a *KeyWebSettingsIntegrationTypeEnum with the same value as that string.
func flattenKeyWebSettingsIntegrationTypeEnum(i interface{}) *KeyWebSettingsIntegrationTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return KeyWebSettingsIntegrationTypeEnumRef(s)
}

// flattenKeyWebSettingsChallengeSecurityPreferenceEnumMap flattens the contents of KeyWebSettingsChallengeSecurityPreferenceEnum from a JSON
// response object.
func flattenKeyWebSettingsChallengeSecurityPreferenceEnumMap(c *Client, i interface{}, res *Key) map[string]KeyWebSettingsChallengeSecurityPreferenceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyWebSettingsChallengeSecurityPreferenceEnum{}
	}

	if len(a) == 0 {
		return map[string]KeyWebSettingsChallengeSecurityPreferenceEnum{}
	}

	items := make(map[string]KeyWebSettingsChallengeSecurityPreferenceEnum)
	for k, item := range a {
		items[k] = *flattenKeyWebSettingsChallengeSecurityPreferenceEnum(item.(interface{}))
	}

	return items
}

// flattenKeyWebSettingsChallengeSecurityPreferenceEnumSlice flattens the contents of KeyWebSettingsChallengeSecurityPreferenceEnum from a JSON
// response object.
func flattenKeyWebSettingsChallengeSecurityPreferenceEnumSlice(c *Client, i interface{}, res *Key) []KeyWebSettingsChallengeSecurityPreferenceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyWebSettingsChallengeSecurityPreferenceEnum{}
	}

	if len(a) == 0 {
		return []KeyWebSettingsChallengeSecurityPreferenceEnum{}
	}

	items := make([]KeyWebSettingsChallengeSecurityPreferenceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyWebSettingsChallengeSecurityPreferenceEnum(item.(interface{})))
	}

	return items
}

// flattenKeyWebSettingsChallengeSecurityPreferenceEnum asserts that an interface is a string, and returns a
// pointer to a *KeyWebSettingsChallengeSecurityPreferenceEnum with the same value as that string.
func flattenKeyWebSettingsChallengeSecurityPreferenceEnum(i interface{}) *KeyWebSettingsChallengeSecurityPreferenceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return KeyWebSettingsChallengeSecurityPreferenceEnumRef(s)
}

// flattenKeyTestingOptionsTestingChallengeEnumMap flattens the contents of KeyTestingOptionsTestingChallengeEnum from a JSON
// response object.
func flattenKeyTestingOptionsTestingChallengeEnumMap(c *Client, i interface{}, res *Key) map[string]KeyTestingOptionsTestingChallengeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]KeyTestingOptionsTestingChallengeEnum{}
	}

	if len(a) == 0 {
		return map[string]KeyTestingOptionsTestingChallengeEnum{}
	}

	items := make(map[string]KeyTestingOptionsTestingChallengeEnum)
	for k, item := range a {
		items[k] = *flattenKeyTestingOptionsTestingChallengeEnum(item.(interface{}))
	}

	return items
}

// flattenKeyTestingOptionsTestingChallengeEnumSlice flattens the contents of KeyTestingOptionsTestingChallengeEnum from a JSON
// response object.
func flattenKeyTestingOptionsTestingChallengeEnumSlice(c *Client, i interface{}, res *Key) []KeyTestingOptionsTestingChallengeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []KeyTestingOptionsTestingChallengeEnum{}
	}

	if len(a) == 0 {
		return []KeyTestingOptionsTestingChallengeEnum{}
	}

	items := make([]KeyTestingOptionsTestingChallengeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenKeyTestingOptionsTestingChallengeEnum(item.(interface{})))
	}

	return items
}

// flattenKeyTestingOptionsTestingChallengeEnum asserts that an interface is a string, and returns a
// pointer to a *KeyTestingOptionsTestingChallengeEnum with the same value as that string.
func flattenKeyTestingOptionsTestingChallengeEnum(i interface{}) *KeyTestingOptionsTestingChallengeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return KeyTestingOptionsTestingChallengeEnumRef(s)
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
	vWebSettings := r.WebSettings
	if vWebSettings == nil {
		// note: explicitly not the empty object.
		vWebSettings = &KeyWebSettings{}
	}
	if err := extractKeyWebSettingsFields(r, vWebSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWebSettings) {
		r.WebSettings = vWebSettings
	}
	vAndroidSettings := r.AndroidSettings
	if vAndroidSettings == nil {
		// note: explicitly not the empty object.
		vAndroidSettings = &KeyAndroidSettings{}
	}
	if err := extractKeyAndroidSettingsFields(r, vAndroidSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAndroidSettings) {
		r.AndroidSettings = vAndroidSettings
	}
	vIosSettings := r.IosSettings
	if vIosSettings == nil {
		// note: explicitly not the empty object.
		vIosSettings = &KeyIosSettings{}
	}
	if err := extractKeyIosSettingsFields(r, vIosSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIosSettings) {
		r.IosSettings = vIosSettings
	}
	vTestingOptions := r.TestingOptions
	if vTestingOptions == nil {
		// note: explicitly not the empty object.
		vTestingOptions = &KeyTestingOptions{}
	}
	if err := extractKeyTestingOptionsFields(r, vTestingOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTestingOptions) {
		r.TestingOptions = vTestingOptions
	}
	return nil
}
func extractKeyWebSettingsFields(r *Key, o *KeyWebSettings) error {
	return nil
}
func extractKeyAndroidSettingsFields(r *Key, o *KeyAndroidSettings) error {
	return nil
}
func extractKeyIosSettingsFields(r *Key, o *KeyIosSettings) error {
	return nil
}
func extractKeyTestingOptionsFields(r *Key, o *KeyTestingOptions) error {
	return nil
}

func postReadExtractKeyFields(r *Key) error {
	vWebSettings := r.WebSettings
	if vWebSettings == nil {
		// note: explicitly not the empty object.
		vWebSettings = &KeyWebSettings{}
	}
	if err := postReadExtractKeyWebSettingsFields(r, vWebSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWebSettings) {
		r.WebSettings = vWebSettings
	}
	vAndroidSettings := r.AndroidSettings
	if vAndroidSettings == nil {
		// note: explicitly not the empty object.
		vAndroidSettings = &KeyAndroidSettings{}
	}
	if err := postReadExtractKeyAndroidSettingsFields(r, vAndroidSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAndroidSettings) {
		r.AndroidSettings = vAndroidSettings
	}
	vIosSettings := r.IosSettings
	if vIosSettings == nil {
		// note: explicitly not the empty object.
		vIosSettings = &KeyIosSettings{}
	}
	if err := postReadExtractKeyIosSettingsFields(r, vIosSettings); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIosSettings) {
		r.IosSettings = vIosSettings
	}
	vTestingOptions := r.TestingOptions
	if vTestingOptions == nil {
		// note: explicitly not the empty object.
		vTestingOptions = &KeyTestingOptions{}
	}
	if err := postReadExtractKeyTestingOptionsFields(r, vTestingOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTestingOptions) {
		r.TestingOptions = vTestingOptions
	}
	return nil
}
func postReadExtractKeyWebSettingsFields(r *Key, o *KeyWebSettings) error {
	return nil
}
func postReadExtractKeyAndroidSettingsFields(r *Key, o *KeyAndroidSettings) error {
	return nil
}
func postReadExtractKeyIosSettingsFields(r *Key, o *KeyIosSettings) error {
	return nil
}
func postReadExtractKeyTestingOptionsFields(r *Key, o *KeyTestingOptions) error {
	return nil
}
