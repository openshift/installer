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
package networkconnectivity

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

func (r *Hub) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *HubRoutingVpcs) validate() error {
	return nil
}
func (r *Hub) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://networkconnectivity.googleapis.com/v1/", params)
}

func (r *Hub) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/hubs/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Hub) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/locations/global/hubs", nr.basePath(), userBasePath, params), nil

}

func (r *Hub) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/hubs?hubId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Hub) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/global/hubs/{{name}}", nr.basePath(), userBasePath, params), nil
}

// hubApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type hubApiOperation interface {
	do(context.Context, *Hub, *Client) error
}

// newUpdateHubUpdateHubRequest creates a request for an
// Hub resource's UpdateHub update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateHubUpdateHubRequest(ctx context.Context, f *Hub, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	return req, nil
}

// marshalUpdateHubUpdateHubRequest converts the update into
// the final JSON request body.
func marshalUpdateHubUpdateHubRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateHubUpdateHubOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateHubUpdateHubOperation) do(ctx context.Context, r *Hub, c *Client) error {
	_, err := c.GetHub(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateHub")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateHubUpdateHubRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateHubUpdateHubRequest(c, req)
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

func (c *Client) listHubRaw(ctx context.Context, r *Hub, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != HubMaxPage {
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

type listHubOperation struct {
	Hubs  []map[string]interface{} `json:"hubs"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listHub(ctx context.Context, r *Hub, pageToken string, pageSize int32) ([]*Hub, string, error) {
	b, err := c.listHubRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listHubOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Hub
	for _, v := range m.Hubs {
		res, err := unmarshalMapHub(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllHub(ctx context.Context, f func(*Hub) bool, resources []*Hub) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteHub(ctx, res)
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

type deleteHubOperation struct{}

func (op *deleteHubOperation) do(ctx context.Context, r *Hub, c *Client) error {
	r, err := c.GetHub(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Hub not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetHub checking for existence. error: %v", err)
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
		_, err := c.GetHub(ctx, r)
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
type createHubOperation struct {
	response map[string]interface{}
}

func (op *createHubOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createHubOperation) do(ctx context.Context, r *Hub, c *Client) error {
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

	if _, err := c.GetHub(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getHubRaw(ctx context.Context, r *Hub) ([]byte, error) {

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

func (c *Client) hubDiffsForRawDesired(ctx context.Context, rawDesired *Hub, opts ...dcl.ApplyOption) (initial, desired *Hub, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Hub
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Hub); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Hub, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetHub(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Hub resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Hub resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Hub resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeHubDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Hub: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Hub: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractHubFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeHubInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Hub: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeHubDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Hub: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffHub(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeHubInitialState(rawInitial, rawDesired *Hub) (*Hub, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeHubDesiredState(rawDesired, rawInitial *Hub, opts ...dcl.ApplyOption) (*Hub, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Hub{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
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
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeHubNewState(c *Client, rawNew, rawDesired *Hub) (*Hub, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
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

	if dcl.IsEmptyValueIndirect(rawNew.UniqueId) && dcl.IsEmptyValueIndirect(rawDesired.UniqueId) {
		rawNew.UniqueId = rawDesired.UniqueId
	} else {
		if dcl.StringCanonicalize(rawDesired.UniqueId, rawNew.UniqueId) {
			rawNew.UniqueId = rawDesired.UniqueId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	rawNew.Project = rawDesired.Project

	if dcl.IsEmptyValueIndirect(rawNew.RoutingVpcs) && dcl.IsEmptyValueIndirect(rawDesired.RoutingVpcs) {
		rawNew.RoutingVpcs = rawDesired.RoutingVpcs
	} else {
		rawNew.RoutingVpcs = canonicalizeNewHubRoutingVpcsSlice(c, rawDesired.RoutingVpcs, rawNew.RoutingVpcs)
	}

	return rawNew, nil
}

func canonicalizeHubRoutingVpcs(des, initial *HubRoutingVpcs, opts ...dcl.ApplyOption) *HubRoutingVpcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &HubRoutingVpcs{}

	if dcl.IsZeroValue(des.Uri) || (dcl.IsEmptyValueIndirect(des.Uri) && dcl.IsEmptyValueIndirect(initial.Uri)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Uri = initial.Uri
	} else {
		cDes.Uri = des.Uri
	}

	return cDes
}

func canonicalizeHubRoutingVpcsSlice(des, initial []HubRoutingVpcs, opts ...dcl.ApplyOption) []HubRoutingVpcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]HubRoutingVpcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeHubRoutingVpcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]HubRoutingVpcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeHubRoutingVpcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewHubRoutingVpcs(c *Client, des, nw *HubRoutingVpcs) *HubRoutingVpcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for HubRoutingVpcs while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewHubRoutingVpcsSet(c *Client, des, nw []HubRoutingVpcs) []HubRoutingVpcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []HubRoutingVpcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareHubRoutingVpcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewHubRoutingVpcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewHubRoutingVpcsSlice(c *Client, des, nw []HubRoutingVpcs) []HubRoutingVpcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []HubRoutingVpcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewHubRoutingVpcs(c, &d, &n))
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
func diffHub(c *Client, desired, actual *Hub, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateHubUpdateHubOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateHubUpdateHubOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UniqueId, actual.UniqueId, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UniqueId")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RoutingVpcs, actual.RoutingVpcs, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareHubRoutingVpcsNewStyle, EmptyObject: EmptyHubRoutingVpcs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RoutingVpcs")); len(ds) != 0 || err != nil {
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
func compareHubRoutingVpcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*HubRoutingVpcs)
	if !ok {
		desiredNotPointer, ok := d.(HubRoutingVpcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a HubRoutingVpcs or *HubRoutingVpcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*HubRoutingVpcs)
	if !ok {
		actualNotPointer, ok := a.(HubRoutingVpcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a HubRoutingVpcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateHubUpdateHubOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
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
func (r *Hub) urlNormalized() *Hub {
	normalized := dcl.Copy(*r).(Hub)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.UniqueId = dcl.SelfLinkToName(r.UniqueId)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Hub) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateHub" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/global/hubs/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Hub resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Hub) marshal(c *Client) ([]byte, error) {
	m, err := expandHub(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Hub: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalHub decodes JSON responses into the Hub resource schema.
func unmarshalHub(b []byte, c *Client, res *Hub) (*Hub, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapHub(m, c, res)
}

func unmarshalMapHub(m map[string]interface{}, c *Client, res *Hub) (*Hub, error) {

	flattened := flattenHub(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandHub expands Hub into a JSON request object.
func expandHub(c *Client, f *Hub) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/global/hubs/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenHub flattens Hub from a JSON request object into the
// Hub type.
func flattenHub(c *Client, i interface{}, res *Hub) *Hub {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Hub{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.UniqueId = dcl.FlattenString(m["uniqueId"])
	resultRes.State = flattenHubStateEnum(m["state"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.RoutingVpcs = flattenHubRoutingVpcsSlice(c, m["routingVpcs"], res)

	return resultRes
}

// expandHubRoutingVpcsMap expands the contents of HubRoutingVpcs into a JSON
// request object.
func expandHubRoutingVpcsMap(c *Client, f map[string]HubRoutingVpcs, res *Hub) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandHubRoutingVpcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandHubRoutingVpcsSlice expands the contents of HubRoutingVpcs into a JSON
// request object.
func expandHubRoutingVpcsSlice(c *Client, f []HubRoutingVpcs, res *Hub) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandHubRoutingVpcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenHubRoutingVpcsMap flattens the contents of HubRoutingVpcs from a JSON
// response object.
func flattenHubRoutingVpcsMap(c *Client, i interface{}, res *Hub) map[string]HubRoutingVpcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]HubRoutingVpcs{}
	}

	if len(a) == 0 {
		return map[string]HubRoutingVpcs{}
	}

	items := make(map[string]HubRoutingVpcs)
	for k, item := range a {
		items[k] = *flattenHubRoutingVpcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenHubRoutingVpcsSlice flattens the contents of HubRoutingVpcs from a JSON
// response object.
func flattenHubRoutingVpcsSlice(c *Client, i interface{}, res *Hub) []HubRoutingVpcs {
	a, ok := i.([]interface{})
	if !ok {
		return []HubRoutingVpcs{}
	}

	if len(a) == 0 {
		return []HubRoutingVpcs{}
	}

	items := make([]HubRoutingVpcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenHubRoutingVpcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandHubRoutingVpcs expands an instance of HubRoutingVpcs into a JSON
// request object.
func expandHubRoutingVpcs(c *Client, f *HubRoutingVpcs, res *Hub) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uri; !dcl.IsEmptyValueIndirect(v) {
		m["uri"] = v
	}

	return m, nil
}

// flattenHubRoutingVpcs flattens an instance of HubRoutingVpcs from a JSON
// response object.
func flattenHubRoutingVpcs(c *Client, i interface{}, res *Hub) *HubRoutingVpcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &HubRoutingVpcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyHubRoutingVpcs
	}
	r.Uri = dcl.FlattenString(m["uri"])

	return r
}

// flattenHubStateEnumMap flattens the contents of HubStateEnum from a JSON
// response object.
func flattenHubStateEnumMap(c *Client, i interface{}, res *Hub) map[string]HubStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]HubStateEnum{}
	}

	if len(a) == 0 {
		return map[string]HubStateEnum{}
	}

	items := make(map[string]HubStateEnum)
	for k, item := range a {
		items[k] = *flattenHubStateEnum(item.(interface{}))
	}

	return items
}

// flattenHubStateEnumSlice flattens the contents of HubStateEnum from a JSON
// response object.
func flattenHubStateEnumSlice(c *Client, i interface{}, res *Hub) []HubStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []HubStateEnum{}
	}

	if len(a) == 0 {
		return []HubStateEnum{}
	}

	items := make([]HubStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenHubStateEnum(item.(interface{})))
	}

	return items
}

// flattenHubStateEnum asserts that an interface is a string, and returns a
// pointer to a *HubStateEnum with the same value as that string.
func flattenHubStateEnum(i interface{}) *HubStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return HubStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Hub) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalHub(b, c, r)
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

type hubDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         hubApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToHubDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]hubDiff, error) {
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
	var diffs []hubDiff
	// For each operation name, create a hubDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := hubDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToHubApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToHubApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (hubApiOperation, error) {
	switch opName {

	case "updateHubUpdateHubOperation":
		return &updateHubUpdateHubOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractHubFields(r *Hub) error {
	return nil
}
func extractHubRoutingVpcsFields(r *Hub, o *HubRoutingVpcs) error {
	return nil
}

func postReadExtractHubFields(r *Hub) error {
	return nil
}
func postReadExtractHubRoutingVpcsFields(r *Hub, o *HubRoutingVpcs) error {
	return nil
}
