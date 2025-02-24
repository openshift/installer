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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *NetworkFirewallPolicy) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicy) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *NetworkFirewallPolicy) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *NetworkFirewallPolicy) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicy) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicy) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

// networkFirewallPolicyApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type networkFirewallPolicyApiOperation interface {
	do(context.Context, *NetworkFirewallPolicy, *Client) error
}

// newUpdateNetworkFirewallPolicyPatchRequest creates a request for an
// NetworkFirewallPolicy resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateNetworkFirewallPolicyPatchRequest(ctx context.Context, f *NetworkFirewallPolicy, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	b, err := c.getNetworkFirewallPolicyRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"fingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["fingerprint"] = rawFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateNetworkFirewallPolicyPatchRequest converts the update into
// the final JSON request body.
func marshalUpdateNetworkFirewallPolicyPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateNetworkFirewallPolicyPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateNetworkFirewallPolicyPatchOperation) do(ctx context.Context, r *NetworkFirewallPolicy, c *Client) error {
	_, err := c.GetNetworkFirewallPolicy(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdateNetworkFirewallPolicyPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateNetworkFirewallPolicyPatchRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

func (c *Client) listNetworkFirewallPolicyRaw(ctx context.Context, r *NetworkFirewallPolicy, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != NetworkFirewallPolicyMaxPage {
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

type listNetworkFirewallPolicyOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listNetworkFirewallPolicy(ctx context.Context, r *NetworkFirewallPolicy, pageToken string, pageSize int32) ([]*NetworkFirewallPolicy, string, error) {
	b, err := c.listNetworkFirewallPolicyRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listNetworkFirewallPolicyOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*NetworkFirewallPolicy
	for _, v := range m.Items {
		res, err := unmarshalMapNetworkFirewallPolicy(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllNetworkFirewallPolicy(ctx context.Context, f func(*NetworkFirewallPolicy) bool, resources []*NetworkFirewallPolicy) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteNetworkFirewallPolicy(ctx, res)
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

type deleteNetworkFirewallPolicyOperation struct{}

func (op *deleteNetworkFirewallPolicyOperation) do(ctx context.Context, r *NetworkFirewallPolicy, c *Client) error {
	r, err := c.GetNetworkFirewallPolicy(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "NetworkFirewallPolicy not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetNetworkFirewallPolicy checking for existence. error: %v", err)
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
		_, err := c.GetNetworkFirewallPolicy(ctx, r)
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
type createNetworkFirewallPolicyOperation struct {
	response map[string]interface{}
}

func (op *createNetworkFirewallPolicyOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createNetworkFirewallPolicyOperation) do(ctx context.Context, r *NetworkFirewallPolicy, c *Client) error {
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

	if _, err := c.GetNetworkFirewallPolicy(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getNetworkFirewallPolicyRaw(ctx context.Context, r *NetworkFirewallPolicy) ([]byte, error) {

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

func (c *Client) networkFirewallPolicyDiffsForRawDesired(ctx context.Context, rawDesired *NetworkFirewallPolicy, opts ...dcl.ApplyOption) (initial, desired *NetworkFirewallPolicy, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *NetworkFirewallPolicy
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*NetworkFirewallPolicy); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected NetworkFirewallPolicy, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetNetworkFirewallPolicy(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a NetworkFirewallPolicy resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve NetworkFirewallPolicy resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that NetworkFirewallPolicy resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeNetworkFirewallPolicyDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for NetworkFirewallPolicy: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for NetworkFirewallPolicy: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractNetworkFirewallPolicyFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeNetworkFirewallPolicyInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for NetworkFirewallPolicy: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeNetworkFirewallPolicyDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for NetworkFirewallPolicy: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffNetworkFirewallPolicy(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeNetworkFirewallPolicyInitialState(rawInitial, rawDesired *NetworkFirewallPolicy) (*NetworkFirewallPolicy, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeNetworkFirewallPolicyDesiredState(rawDesired, rawInitial *NetworkFirewallPolicy, opts ...dcl.ApplyOption) (*NetworkFirewallPolicy, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &NetworkFirewallPolicy{}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
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

func canonicalizeNetworkFirewallPolicyNewState(c *Client, rawNew, rawDesired *NetworkFirewallPolicy) (*NetworkFirewallPolicy, error) {

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.CreationTimestamp) && dcl.IsEmptyValueIndirect(rawDesired.CreationTimestamp) {
		rawNew.CreationTimestamp = rawDesired.CreationTimestamp
	} else {
		if dcl.StringCanonicalize(rawDesired.CreationTimestamp, rawNew.CreationTimestamp) {
			rawNew.CreationTimestamp = rawDesired.CreationTimestamp
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
		if dcl.StringCanonicalize(rawDesired.Id, rawNew.Id) {
			rawNew.Id = rawDesired.Id
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Fingerprint) && dcl.IsEmptyValueIndirect(rawDesired.Fingerprint) {
		rawNew.Fingerprint = rawDesired.Fingerprint
	} else {
		if dcl.StringCanonicalize(rawDesired.Fingerprint, rawNew.Fingerprint) {
			rawNew.Fingerprint = rawDesired.Fingerprint
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLinkWithId) && dcl.IsEmptyValueIndirect(rawDesired.SelfLinkWithId) {
		rawNew.SelfLinkWithId = rawDesired.SelfLinkWithId
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLinkWithId, rawNew.SelfLinkWithId) {
			rawNew.SelfLinkWithId = rawDesired.SelfLinkWithId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RuleTupleCount) && dcl.IsEmptyValueIndirect(rawDesired.RuleTupleCount) {
		rawNew.RuleTupleCount = rawDesired.RuleTupleCount
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Region) && dcl.IsEmptyValueIndirect(rawDesired.Region) {
		rawNew.Region = rawDesired.Region
	} else {
		if dcl.StringCanonicalize(rawDesired.Region, rawNew.Region) {
			rawNew.Region = rawDesired.Region
		}
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffNetworkFirewallPolicy(c *Client, desired, actual *NetworkFirewallPolicy, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreationTimestamp, actual.CreationTimestamp, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreationTimestamp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkFirewallPolicyPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fingerprint, actual.Fingerprint, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SelfLinkWithId, actual.SelfLinkWithId, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLinkWithId")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
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

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *NetworkFirewallPolicy) urlNormalized() *NetworkFirewallPolicy {
	normalized := dcl.Copy(*r).(NetworkFirewallPolicy)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.CreationTimestamp = dcl.SelfLinkToName(r.CreationTimestamp)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Id = dcl.SelfLinkToName(r.Id)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Fingerprint = dcl.SelfLinkToName(r.Fingerprint)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.SelfLinkWithId = dcl.SelfLinkToName(r.SelfLinkWithId)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *NetworkFirewallPolicy) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		if dcl.IsRegion(nr.Location) {
			return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{name}}", nr.basePath(), userBasePath, fields), nil
		}

		return dcl.URL("projects/{{project}}/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the NetworkFirewallPolicy resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *NetworkFirewallPolicy) marshal(c *Client) ([]byte, error) {
	m, err := expandNetworkFirewallPolicy(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling NetworkFirewallPolicy: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalNetworkFirewallPolicy decodes JSON responses into the NetworkFirewallPolicy resource schema.
func unmarshalNetworkFirewallPolicy(b []byte, c *Client, res *NetworkFirewallPolicy) (*NetworkFirewallPolicy, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNetworkFirewallPolicy(m, c, res)
}

func unmarshalMapNetworkFirewallPolicy(m map[string]interface{}, c *Client, res *NetworkFirewallPolicy) (*NetworkFirewallPolicy, error) {

	flattened := flattenNetworkFirewallPolicy(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNetworkFirewallPolicy expands NetworkFirewallPolicy into a JSON request object.
func expandNetworkFirewallPolicy(c *Client, f *NetworkFirewallPolicy) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
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

// flattenNetworkFirewallPolicy flattens NetworkFirewallPolicy from a JSON request object into the
// NetworkFirewallPolicy type.
func flattenNetworkFirewallPolicy(c *Client, i interface{}, res *NetworkFirewallPolicy) *NetworkFirewallPolicy {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &NetworkFirewallPolicy{}
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.CreationTimestamp = dcl.FlattenString(m["creationTimestamp"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Id = dcl.FlattenString(m["id"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Fingerprint = dcl.FlattenString(m["fingerprint"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.SelfLinkWithId = dcl.FlattenString(m["selfLinkWithId"])
	resultRes.RuleTupleCount = dcl.FlattenInteger(m["ruleTupleCount"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *NetworkFirewallPolicy) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNetworkFirewallPolicy(b, c, r)
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

type networkFirewallPolicyDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         networkFirewallPolicyApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToNetworkFirewallPolicyDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]networkFirewallPolicyDiff, error) {
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
	var diffs []networkFirewallPolicyDiff
	// For each operation name, create a networkFirewallPolicyDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := networkFirewallPolicyDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToNetworkFirewallPolicyApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToNetworkFirewallPolicyApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (networkFirewallPolicyApiOperation, error) {
	switch opName {

	case "updateNetworkFirewallPolicyPatchOperation":
		return &updateNetworkFirewallPolicyPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractNetworkFirewallPolicyFields(r *NetworkFirewallPolicy) error {
	return nil
}

func postReadExtractNetworkFirewallPolicyFields(r *NetworkFirewallPolicy) error {
	return nil
}
