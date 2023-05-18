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
package dataplex

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

func (r *Lake) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Metastore) {
		if err := r.Metastore.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AssetStatus) {
		if err := r.AssetStatus.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MetastoreStatus) {
		if err := r.MetastoreStatus.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *LakeMetastore) validate() error {
	return nil
}
func (r *LakeAssetStatus) validate() error {
	return nil
}
func (r *LakeMetastoreStatus) validate() error {
	return nil
}
func (r *Lake) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://dataplex.googleapis.com/v1/", params)
}

func (r *Lake) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Lake) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes", nr.basePath(), userBasePath, params), nil

}

func (r *Lake) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes?lakeId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Lake) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Lake) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *Lake) SetPolicyVerb() string {
	return ""
}

func (r *Lake) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *Lake) IAMPolicyVersion() int {
	return 3
}

// lakeApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type lakeApiOperation interface {
	do(context.Context, *Lake, *Client) error
}

// newUpdateLakeUpdateLakeRequest creates a request for an
// Lake resource's UpdateLake update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateLakeUpdateLakeRequest(ctx context.Context, f *Lake, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.DeriveField("projects/%s/locations/%s/lakes/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["name"] = v
	}
	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandLakeMetastore(c, f.Metastore, res); err != nil {
		return nil, fmt.Errorf("error expanding Metastore into metastore: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["metastore"] = v
	}
	if v, err := expandLakeAssetStatus(c, f.AssetStatus, res); err != nil {
		return nil, fmt.Errorf("error expanding AssetStatus into assetStatus: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["assetStatus"] = v
	}
	if v, err := expandLakeMetastoreStatus(c, f.MetastoreStatus, res); err != nil {
		return nil, fmt.Errorf("error expanding MetastoreStatus into metastoreStatus: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["metastoreStatus"] = v
	}
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/lakes/%s", *f.Project, *f.Location, *f.Name)

	return req, nil
}

// marshalUpdateLakeUpdateLakeRequest converts the update into
// the final JSON request body.
func marshalUpdateLakeUpdateLakeRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateLakeUpdateLakeOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateLakeUpdateLakeOperation) do(ctx context.Context, r *Lake, c *Client) error {
	_, err := c.GetLake(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateLake")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateLakeUpdateLakeRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateLakeUpdateLakeRequest(c, req)
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

func (c *Client) listLakeRaw(ctx context.Context, r *Lake, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != LakeMaxPage {
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

type listLakeOperation struct {
	Lakes []map[string]interface{} `json:"lakes"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listLake(ctx context.Context, r *Lake, pageToken string, pageSize int32) ([]*Lake, string, error) {
	b, err := c.listLakeRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listLakeOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Lake
	for _, v := range m.Lakes {
		res, err := unmarshalMapLake(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllLake(ctx context.Context, f func(*Lake) bool, resources []*Lake) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteLake(ctx, res)
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

type deleteLakeOperation struct{}

func (op *deleteLakeOperation) do(ctx context.Context, r *Lake, c *Client) error {
	r, err := c.GetLake(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Lake not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetLake checking for existence. error: %v", err)
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

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetLake(ctx, r)
		if dcl.IsNotFound(err) {
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
type createLakeOperation struct {
	response map[string]interface{}
}

func (op *createLakeOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createLakeOperation) do(ctx context.Context, r *Lake, c *Client) error {
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

	if _, err := c.GetLake(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getLakeRaw(ctx context.Context, r *Lake) ([]byte, error) {

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

func (c *Client) lakeDiffsForRawDesired(ctx context.Context, rawDesired *Lake, opts ...dcl.ApplyOption) (initial, desired *Lake, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Lake
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Lake); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Lake, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetLake(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Lake resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Lake resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Lake resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeLakeDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Lake: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Lake: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractLakeFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeLakeInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Lake: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeLakeDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Lake: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffLake(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeLakeInitialState(rawInitial, rawDesired *Lake) (*Lake, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeLakeDesiredState(rawDesired, rawInitial *Lake, opts ...dcl.ApplyOption) (*Lake, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Metastore = canonicalizeLakeMetastore(rawDesired.Metastore, nil, opts...)
		rawDesired.AssetStatus = canonicalizeLakeAssetStatus(rawDesired.AssetStatus, nil, opts...)
		rawDesired.MetastoreStatus = canonicalizeLakeMetastoreStatus(rawDesired.MetastoreStatus, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Lake{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	canonicalDesired.Metastore = canonicalizeLakeMetastore(rawDesired.Metastore, rawInitial.Metastore, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	return canonicalDesired, nil
}

func canonicalizeLakeNewState(c *Client, rawNew, rawDesired *Lake) (*Lake, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceAccount) && dcl.IsEmptyValueIndirect(rawDesired.ServiceAccount) {
		rawNew.ServiceAccount = rawDesired.ServiceAccount
	} else {
		if dcl.StringCanonicalize(rawDesired.ServiceAccount, rawNew.ServiceAccount) {
			rawNew.ServiceAccount = rawDesired.ServiceAccount
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Metastore) && dcl.IsEmptyValueIndirect(rawDesired.Metastore) {
		rawNew.Metastore = rawDesired.Metastore
	} else {
		rawNew.Metastore = canonicalizeNewLakeMetastore(c, rawDesired.Metastore, rawNew.Metastore)
	}

	if dcl.IsEmptyValueIndirect(rawNew.AssetStatus) && dcl.IsEmptyValueIndirect(rawDesired.AssetStatus) {
		rawNew.AssetStatus = rawDesired.AssetStatus
	} else {
		rawNew.AssetStatus = canonicalizeNewLakeAssetStatus(c, rawDesired.AssetStatus, rawNew.AssetStatus)
	}

	if dcl.IsEmptyValueIndirect(rawNew.MetastoreStatus) && dcl.IsEmptyValueIndirect(rawDesired.MetastoreStatus) {
		rawNew.MetastoreStatus = rawDesired.MetastoreStatus
	} else {
		rawNew.MetastoreStatus = canonicalizeNewLakeMetastoreStatus(c, rawDesired.MetastoreStatus, rawNew.MetastoreStatus)
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeLakeMetastore(des, initial *LakeMetastore, opts ...dcl.ApplyOption) *LakeMetastore {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LakeMetastore{}

	if dcl.StringCanonicalize(des.Service, initial.Service) || dcl.IsZeroValue(des.Service) {
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}

	return cDes
}

func canonicalizeLakeMetastoreSlice(des, initial []LakeMetastore, opts ...dcl.ApplyOption) []LakeMetastore {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LakeMetastore, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLakeMetastore(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LakeMetastore, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLakeMetastore(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLakeMetastore(c *Client, des, nw *LakeMetastore) *LakeMetastore {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LakeMetastore while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Service, nw.Service) {
		nw.Service = des.Service
	}

	return nw
}

func canonicalizeNewLakeMetastoreSet(c *Client, des, nw []LakeMetastore) []LakeMetastore {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LakeMetastore
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLakeMetastoreNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLakeMetastore(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLakeMetastoreSlice(c *Client, des, nw []LakeMetastore) []LakeMetastore {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LakeMetastore
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLakeMetastore(c, &d, &n))
	}

	return items
}

func canonicalizeLakeAssetStatus(des, initial *LakeAssetStatus, opts ...dcl.ApplyOption) *LakeAssetStatus {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LakeAssetStatus{}

	if dcl.IsZeroValue(des.UpdateTime) || (dcl.IsEmptyValueIndirect(des.UpdateTime) && dcl.IsEmptyValueIndirect(initial.UpdateTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UpdateTime = initial.UpdateTime
	} else {
		cDes.UpdateTime = des.UpdateTime
	}
	if dcl.IsZeroValue(des.ActiveAssets) || (dcl.IsEmptyValueIndirect(des.ActiveAssets) && dcl.IsEmptyValueIndirect(initial.ActiveAssets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ActiveAssets = initial.ActiveAssets
	} else {
		cDes.ActiveAssets = des.ActiveAssets
	}
	if dcl.IsZeroValue(des.SecurityPolicyApplyingAssets) || (dcl.IsEmptyValueIndirect(des.SecurityPolicyApplyingAssets) && dcl.IsEmptyValueIndirect(initial.SecurityPolicyApplyingAssets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SecurityPolicyApplyingAssets = initial.SecurityPolicyApplyingAssets
	} else {
		cDes.SecurityPolicyApplyingAssets = des.SecurityPolicyApplyingAssets
	}

	return cDes
}

func canonicalizeLakeAssetStatusSlice(des, initial []LakeAssetStatus, opts ...dcl.ApplyOption) []LakeAssetStatus {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LakeAssetStatus, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLakeAssetStatus(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LakeAssetStatus, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLakeAssetStatus(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLakeAssetStatus(c *Client, des, nw *LakeAssetStatus) *LakeAssetStatus {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LakeAssetStatus while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewLakeAssetStatusSet(c *Client, des, nw []LakeAssetStatus) []LakeAssetStatus {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LakeAssetStatus
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLakeAssetStatusNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLakeAssetStatus(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLakeAssetStatusSlice(c *Client, des, nw []LakeAssetStatus) []LakeAssetStatus {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LakeAssetStatus
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLakeAssetStatus(c, &d, &n))
	}

	return items
}

func canonicalizeLakeMetastoreStatus(des, initial *LakeMetastoreStatus, opts ...dcl.ApplyOption) *LakeMetastoreStatus {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LakeMetastoreStatus{}

	if dcl.IsZeroValue(des.State) || (dcl.IsEmptyValueIndirect(des.State) && dcl.IsEmptyValueIndirect(initial.State)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.State = initial.State
	} else {
		cDes.State = des.State
	}
	if dcl.StringCanonicalize(des.Message, initial.Message) || dcl.IsZeroValue(des.Message) {
		cDes.Message = initial.Message
	} else {
		cDes.Message = des.Message
	}
	if dcl.IsZeroValue(des.UpdateTime) || (dcl.IsEmptyValueIndirect(des.UpdateTime) && dcl.IsEmptyValueIndirect(initial.UpdateTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UpdateTime = initial.UpdateTime
	} else {
		cDes.UpdateTime = des.UpdateTime
	}
	if dcl.StringCanonicalize(des.Endpoint, initial.Endpoint) || dcl.IsZeroValue(des.Endpoint) {
		cDes.Endpoint = initial.Endpoint
	} else {
		cDes.Endpoint = des.Endpoint
	}

	return cDes
}

func canonicalizeLakeMetastoreStatusSlice(des, initial []LakeMetastoreStatus, opts ...dcl.ApplyOption) []LakeMetastoreStatus {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LakeMetastoreStatus, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLakeMetastoreStatus(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LakeMetastoreStatus, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLakeMetastoreStatus(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLakeMetastoreStatus(c *Client, des, nw *LakeMetastoreStatus) *LakeMetastoreStatus {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LakeMetastoreStatus while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Message, nw.Message) {
		nw.Message = des.Message
	}
	if dcl.StringCanonicalize(des.Endpoint, nw.Endpoint) {
		nw.Endpoint = des.Endpoint
	}

	return nw
}

func canonicalizeNewLakeMetastoreStatusSet(c *Client, des, nw []LakeMetastoreStatus) []LakeMetastoreStatus {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LakeMetastoreStatus
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLakeMetastoreStatusNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLakeMetastoreStatus(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLakeMetastoreStatusSlice(c *Client, des, nw []LakeMetastoreStatus) []LakeMetastoreStatus {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LakeMetastoreStatus
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLakeMetastoreStatus(c, &d, &n))
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
func diffLake(c *Client, desired, actual *Lake, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccount, actual.ServiceAccount, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metastore, actual.Metastore, dcl.DiffInfo{ObjectFunction: compareLakeMetastoreNewStyle, EmptyObject: EmptyLakeMetastore, OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Metastore")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AssetStatus, actual.AssetStatus, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareLakeAssetStatusNewStyle, EmptyObject: EmptyLakeAssetStatus, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AssetStatus")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetastoreStatus, actual.MetastoreStatus, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareLakeMetastoreStatusNewStyle, EmptyObject: EmptyLakeMetastoreStatus, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetastoreStatus")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
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
func compareLakeMetastoreNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LakeMetastore)
	if !ok {
		desiredNotPointer, ok := d.(LakeMetastore)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeMetastore or *LakeMetastore", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LakeMetastore)
	if !ok {
		actualNotPointer, ok := a.(LakeMetastore)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeMetastore", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLakeAssetStatusNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LakeAssetStatus)
	if !ok {
		desiredNotPointer, ok := d.(LakeAssetStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeAssetStatus or *LakeAssetStatus", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LakeAssetStatus)
	if !ok {
		actualNotPointer, ok := a.(LakeAssetStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeAssetStatus", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ActiveAssets, actual.ActiveAssets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("ActiveAssets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityPolicyApplyingAssets, actual.SecurityPolicyApplyingAssets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("SecurityPolicyApplyingAssets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLakeMetastoreStatusNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LakeMetastoreStatus)
	if !ok {
		desiredNotPointer, ok := d.(LakeMetastoreStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeMetastoreStatus or *LakeMetastoreStatus", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LakeMetastoreStatus)
	if !ok {
		actualNotPointer, ok := a.(LakeMetastoreStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LakeMetastoreStatus", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Message, actual.Message, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Message")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Endpoint, actual.Endpoint, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLakeUpdateLakeOperation")}, fn.AddNest("Endpoint")); len(ds) != 0 || err != nil {
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
func (r *Lake) urlNormalized() *Lake {
	normalized := dcl.Copy(*r).(Lake)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.ServiceAccount = dcl.SelfLinkToName(r.ServiceAccount)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *Lake) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateLake" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Lake resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Lake) marshal(c *Client) ([]byte, error) {
	m, err := expandLake(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Lake: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalLake decodes JSON responses into the Lake resource schema.
func unmarshalLake(b []byte, c *Client, res *Lake) (*Lake, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapLake(m, c, res)
}

func unmarshalMapLake(m map[string]interface{}, c *Client, res *Lake) (*Lake, error) {

	flattened := flattenLake(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandLake expands Lake into a JSON request object.
func expandLake(c *Client, f *Lake) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/lakes/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandLakeMetastore(c, f.Metastore, res); err != nil {
		return nil, fmt.Errorf("error expanding Metastore into metastore: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metastore"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenLake flattens Lake from a JSON request object into the
// Lake type.
func flattenLake(c *Client, i interface{}, res *Lake) *Lake {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Lake{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.State = flattenLakeStateEnum(m["state"])
	resultRes.ServiceAccount = dcl.FlattenString(m["serviceAccount"])
	resultRes.Metastore = flattenLakeMetastore(c, m["metastore"], res)
	resultRes.AssetStatus = flattenLakeAssetStatus(c, m["assetStatus"], res)
	resultRes.MetastoreStatus = flattenLakeMetastoreStatus(c, m["metastoreStatus"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandLakeMetastoreMap expands the contents of LakeMetastore into a JSON
// request object.
func expandLakeMetastoreMap(c *Client, f map[string]LakeMetastore, res *Lake) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLakeMetastore(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLakeMetastoreSlice expands the contents of LakeMetastore into a JSON
// request object.
func expandLakeMetastoreSlice(c *Client, f []LakeMetastore, res *Lake) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLakeMetastore(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLakeMetastoreMap flattens the contents of LakeMetastore from a JSON
// response object.
func flattenLakeMetastoreMap(c *Client, i interface{}, res *Lake) map[string]LakeMetastore {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LakeMetastore{}
	}

	if len(a) == 0 {
		return map[string]LakeMetastore{}
	}

	items := make(map[string]LakeMetastore)
	for k, item := range a {
		items[k] = *flattenLakeMetastore(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLakeMetastoreSlice flattens the contents of LakeMetastore from a JSON
// response object.
func flattenLakeMetastoreSlice(c *Client, i interface{}, res *Lake) []LakeMetastore {
	a, ok := i.([]interface{})
	if !ok {
		return []LakeMetastore{}
	}

	if len(a) == 0 {
		return []LakeMetastore{}
	}

	items := make([]LakeMetastore, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLakeMetastore(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLakeMetastore expands an instance of LakeMetastore into a JSON
// request object.
func expandLakeMetastore(c *Client, f *LakeMetastore, res *Lake) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}

	return m, nil
}

// flattenLakeMetastore flattens an instance of LakeMetastore from a JSON
// response object.
func flattenLakeMetastore(c *Client, i interface{}, res *Lake) *LakeMetastore {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LakeMetastore{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLakeMetastore
	}
	r.Service = dcl.FlattenString(m["service"])

	return r
}

// expandLakeAssetStatusMap expands the contents of LakeAssetStatus into a JSON
// request object.
func expandLakeAssetStatusMap(c *Client, f map[string]LakeAssetStatus, res *Lake) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLakeAssetStatus(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLakeAssetStatusSlice expands the contents of LakeAssetStatus into a JSON
// request object.
func expandLakeAssetStatusSlice(c *Client, f []LakeAssetStatus, res *Lake) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLakeAssetStatus(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLakeAssetStatusMap flattens the contents of LakeAssetStatus from a JSON
// response object.
func flattenLakeAssetStatusMap(c *Client, i interface{}, res *Lake) map[string]LakeAssetStatus {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LakeAssetStatus{}
	}

	if len(a) == 0 {
		return map[string]LakeAssetStatus{}
	}

	items := make(map[string]LakeAssetStatus)
	for k, item := range a {
		items[k] = *flattenLakeAssetStatus(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLakeAssetStatusSlice flattens the contents of LakeAssetStatus from a JSON
// response object.
func flattenLakeAssetStatusSlice(c *Client, i interface{}, res *Lake) []LakeAssetStatus {
	a, ok := i.([]interface{})
	if !ok {
		return []LakeAssetStatus{}
	}

	if len(a) == 0 {
		return []LakeAssetStatus{}
	}

	items := make([]LakeAssetStatus, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLakeAssetStatus(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLakeAssetStatus expands an instance of LakeAssetStatus into a JSON
// request object.
func expandLakeAssetStatus(c *Client, f *LakeAssetStatus, res *Lake) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.UpdateTime; !dcl.IsEmptyValueIndirect(v) {
		m["updateTime"] = v
	}
	if v := f.ActiveAssets; !dcl.IsEmptyValueIndirect(v) {
		m["activeAssets"] = v
	}
	if v := f.SecurityPolicyApplyingAssets; !dcl.IsEmptyValueIndirect(v) {
		m["securityPolicyApplyingAssets"] = v
	}

	return m, nil
}

// flattenLakeAssetStatus flattens an instance of LakeAssetStatus from a JSON
// response object.
func flattenLakeAssetStatus(c *Client, i interface{}, res *Lake) *LakeAssetStatus {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LakeAssetStatus{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLakeAssetStatus
	}
	r.UpdateTime = dcl.FlattenString(m["updateTime"])
	r.ActiveAssets = dcl.FlattenInteger(m["activeAssets"])
	r.SecurityPolicyApplyingAssets = dcl.FlattenInteger(m["securityPolicyApplyingAssets"])

	return r
}

// expandLakeMetastoreStatusMap expands the contents of LakeMetastoreStatus into a JSON
// request object.
func expandLakeMetastoreStatusMap(c *Client, f map[string]LakeMetastoreStatus, res *Lake) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLakeMetastoreStatus(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLakeMetastoreStatusSlice expands the contents of LakeMetastoreStatus into a JSON
// request object.
func expandLakeMetastoreStatusSlice(c *Client, f []LakeMetastoreStatus, res *Lake) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLakeMetastoreStatus(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLakeMetastoreStatusMap flattens the contents of LakeMetastoreStatus from a JSON
// response object.
func flattenLakeMetastoreStatusMap(c *Client, i interface{}, res *Lake) map[string]LakeMetastoreStatus {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LakeMetastoreStatus{}
	}

	if len(a) == 0 {
		return map[string]LakeMetastoreStatus{}
	}

	items := make(map[string]LakeMetastoreStatus)
	for k, item := range a {
		items[k] = *flattenLakeMetastoreStatus(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLakeMetastoreStatusSlice flattens the contents of LakeMetastoreStatus from a JSON
// response object.
func flattenLakeMetastoreStatusSlice(c *Client, i interface{}, res *Lake) []LakeMetastoreStatus {
	a, ok := i.([]interface{})
	if !ok {
		return []LakeMetastoreStatus{}
	}

	if len(a) == 0 {
		return []LakeMetastoreStatus{}
	}

	items := make([]LakeMetastoreStatus, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLakeMetastoreStatus(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLakeMetastoreStatus expands an instance of LakeMetastoreStatus into a JSON
// request object.
func expandLakeMetastoreStatus(c *Client, f *LakeMetastoreStatus, res *Lake) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.State; !dcl.IsEmptyValueIndirect(v) {
		m["state"] = v
	}
	if v := f.Message; !dcl.IsEmptyValueIndirect(v) {
		m["message"] = v
	}
	if v := f.UpdateTime; !dcl.IsEmptyValueIndirect(v) {
		m["updateTime"] = v
	}
	if v := f.Endpoint; !dcl.IsEmptyValueIndirect(v) {
		m["endpoint"] = v
	}

	return m, nil
}

// flattenLakeMetastoreStatus flattens an instance of LakeMetastoreStatus from a JSON
// response object.
func flattenLakeMetastoreStatus(c *Client, i interface{}, res *Lake) *LakeMetastoreStatus {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LakeMetastoreStatus{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLakeMetastoreStatus
	}
	r.State = flattenLakeMetastoreStatusStateEnum(m["state"])
	r.Message = dcl.FlattenString(m["message"])
	r.UpdateTime = dcl.FlattenString(m["updateTime"])
	r.Endpoint = dcl.FlattenString(m["endpoint"])

	return r
}

// flattenLakeStateEnumMap flattens the contents of LakeStateEnum from a JSON
// response object.
func flattenLakeStateEnumMap(c *Client, i interface{}, res *Lake) map[string]LakeStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LakeStateEnum{}
	}

	if len(a) == 0 {
		return map[string]LakeStateEnum{}
	}

	items := make(map[string]LakeStateEnum)
	for k, item := range a {
		items[k] = *flattenLakeStateEnum(item.(interface{}))
	}

	return items
}

// flattenLakeStateEnumSlice flattens the contents of LakeStateEnum from a JSON
// response object.
func flattenLakeStateEnumSlice(c *Client, i interface{}, res *Lake) []LakeStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LakeStateEnum{}
	}

	if len(a) == 0 {
		return []LakeStateEnum{}
	}

	items := make([]LakeStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLakeStateEnum(item.(interface{})))
	}

	return items
}

// flattenLakeStateEnum asserts that an interface is a string, and returns a
// pointer to a *LakeStateEnum with the same value as that string.
func flattenLakeStateEnum(i interface{}) *LakeStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LakeStateEnumRef(s)
}

// flattenLakeMetastoreStatusStateEnumMap flattens the contents of LakeMetastoreStatusStateEnum from a JSON
// response object.
func flattenLakeMetastoreStatusStateEnumMap(c *Client, i interface{}, res *Lake) map[string]LakeMetastoreStatusStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LakeMetastoreStatusStateEnum{}
	}

	if len(a) == 0 {
		return map[string]LakeMetastoreStatusStateEnum{}
	}

	items := make(map[string]LakeMetastoreStatusStateEnum)
	for k, item := range a {
		items[k] = *flattenLakeMetastoreStatusStateEnum(item.(interface{}))
	}

	return items
}

// flattenLakeMetastoreStatusStateEnumSlice flattens the contents of LakeMetastoreStatusStateEnum from a JSON
// response object.
func flattenLakeMetastoreStatusStateEnumSlice(c *Client, i interface{}, res *Lake) []LakeMetastoreStatusStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LakeMetastoreStatusStateEnum{}
	}

	if len(a) == 0 {
		return []LakeMetastoreStatusStateEnum{}
	}

	items := make([]LakeMetastoreStatusStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLakeMetastoreStatusStateEnum(item.(interface{})))
	}

	return items
}

// flattenLakeMetastoreStatusStateEnum asserts that an interface is a string, and returns a
// pointer to a *LakeMetastoreStatusStateEnum with the same value as that string.
func flattenLakeMetastoreStatusStateEnum(i interface{}) *LakeMetastoreStatusStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LakeMetastoreStatusStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Lake) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalLake(b, c, r)
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

type lakeDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         lakeApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToLakeDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]lakeDiff, error) {
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
	var diffs []lakeDiff
	// For each operation name, create a lakeDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := lakeDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToLakeApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToLakeApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (lakeApiOperation, error) {
	switch opName {

	case "updateLakeUpdateLakeOperation":
		return &updateLakeUpdateLakeOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractLakeFields(r *Lake) error {
	vMetastore := r.Metastore
	if vMetastore == nil {
		// note: explicitly not the empty object.
		vMetastore = &LakeMetastore{}
	}
	if err := extractLakeMetastoreFields(r, vMetastore); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastore) {
		r.Metastore = vMetastore
	}
	vAssetStatus := r.AssetStatus
	if vAssetStatus == nil {
		// note: explicitly not the empty object.
		vAssetStatus = &LakeAssetStatus{}
	}
	if err := extractLakeAssetStatusFields(r, vAssetStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAssetStatus) {
		r.AssetStatus = vAssetStatus
	}
	vMetastoreStatus := r.MetastoreStatus
	if vMetastoreStatus == nil {
		// note: explicitly not the empty object.
		vMetastoreStatus = &LakeMetastoreStatus{}
	}
	if err := extractLakeMetastoreStatusFields(r, vMetastoreStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreStatus) {
		r.MetastoreStatus = vMetastoreStatus
	}
	return nil
}
func extractLakeMetastoreFields(r *Lake, o *LakeMetastore) error {
	return nil
}
func extractLakeAssetStatusFields(r *Lake, o *LakeAssetStatus) error {
	return nil
}
func extractLakeMetastoreStatusFields(r *Lake, o *LakeMetastoreStatus) error {
	return nil
}

func postReadExtractLakeFields(r *Lake) error {
	vMetastore := r.Metastore
	if vMetastore == nil {
		// note: explicitly not the empty object.
		vMetastore = &LakeMetastore{}
	}
	if err := postReadExtractLakeMetastoreFields(r, vMetastore); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastore) {
		r.Metastore = vMetastore
	}
	vAssetStatus := r.AssetStatus
	if vAssetStatus == nil {
		// note: explicitly not the empty object.
		vAssetStatus = &LakeAssetStatus{}
	}
	if err := postReadExtractLakeAssetStatusFields(r, vAssetStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAssetStatus) {
		r.AssetStatus = vAssetStatus
	}
	vMetastoreStatus := r.MetastoreStatus
	if vMetastoreStatus == nil {
		// note: explicitly not the empty object.
		vMetastoreStatus = &LakeMetastoreStatus{}
	}
	if err := postReadExtractLakeMetastoreStatusFields(r, vMetastoreStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreStatus) {
		r.MetastoreStatus = vMetastoreStatus
	}
	return nil
}
func postReadExtractLakeMetastoreFields(r *Lake, o *LakeMetastore) error {
	return nil
}
func postReadExtractLakeAssetStatusFields(r *Lake, o *LakeAssetStatus) error {
	return nil
}
func postReadExtractLakeMetastoreStatusFields(r *Lake, o *LakeMetastoreStatus) error {
	return nil
}
