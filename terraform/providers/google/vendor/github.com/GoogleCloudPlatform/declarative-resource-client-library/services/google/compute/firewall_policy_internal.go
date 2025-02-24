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

func (r *FirewallPolicy) validate() error {

	if err := dcl.Required(r, "shortName"); err != nil {
		return err
	}
	if err := dcl.Required(r, "parent"); err != nil {
		return err
	}
	return nil
}
func (r *FirewallPolicy) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *FirewallPolicy) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *FirewallPolicy) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("locations/global/firewallPolicies?parentId={{parent}}", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicy) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("locations/global/firewallPolicies?parentId={{parent}}", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicy) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

// firewallPolicyApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type firewallPolicyApiOperation interface {
	do(context.Context, *FirewallPolicy, *Client) error
}

// newUpdateFirewallPolicyPatchRequest creates a request for an
// FirewallPolicy resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateFirewallPolicyPatchRequest(ctx context.Context, f *FirewallPolicy, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	b, err := c.getFirewallPolicyRaw(ctx, f)
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

// marshalUpdateFirewallPolicyPatchRequest converts the update into
// the final JSON request body.
func marshalUpdateFirewallPolicyPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateFirewallPolicyPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateFirewallPolicyPatchOperation) do(ctx context.Context, r *FirewallPolicy, c *Client) error {
	_, err := c.GetFirewallPolicy(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdateFirewallPolicyPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateFirewallPolicyPatchRequest(c, req)
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

func (c *Client) listFirewallPolicyRaw(ctx context.Context, r *FirewallPolicy, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != FirewallPolicyMaxPage {
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

type listFirewallPolicyOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listFirewallPolicy(ctx context.Context, r *FirewallPolicy, pageToken string, pageSize int32) ([]*FirewallPolicy, string, error) {
	b, err := c.listFirewallPolicyRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listFirewallPolicyOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*FirewallPolicy
	for _, v := range m.Items {
		res, err := unmarshalMapFirewallPolicy(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Parent = r.Parent
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllFirewallPolicy(ctx context.Context, f func(*FirewallPolicy) bool, resources []*FirewallPolicy) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteFirewallPolicy(ctx, res)
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

type deleteFirewallPolicyOperation struct{}

func (op *deleteFirewallPolicyOperation) do(ctx context.Context, r *FirewallPolicy, c *Client) error {
	r, err := c.GetFirewallPolicy(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "FirewallPolicy not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetFirewallPolicy checking for existence. error: %v", err)
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
		_, err := c.GetFirewallPolicy(ctx, r)
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
type createFirewallPolicyOperation struct {
	response map[string]interface{}
}

func (op *createFirewallPolicyOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (c *Client) getFirewallPolicyRaw(ctx context.Context, r *FirewallPolicy) ([]byte, error) {

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

func (c *Client) firewallPolicyDiffsForRawDesired(ctx context.Context, rawDesired *FirewallPolicy, opts ...dcl.ApplyOption) (initial, desired *FirewallPolicy, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *FirewallPolicy
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*FirewallPolicy); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected FirewallPolicy, got %T", sh)
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
		desired, err := canonicalizeFirewallPolicyDesiredState(rawDesired, nil)
		return nil, desired, nil, err
	}
	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetFirewallPolicy(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a FirewallPolicy resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve FirewallPolicy resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that FirewallPolicy resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeFirewallPolicyDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for FirewallPolicy: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for FirewallPolicy: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractFirewallPolicyFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeFirewallPolicyInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for FirewallPolicy: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeFirewallPolicyDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for FirewallPolicy: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffFirewallPolicy(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeFirewallPolicyInitialState(rawInitial, rawDesired *FirewallPolicy) (*FirewallPolicy, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeFirewallPolicyDesiredState(rawDesired, rawInitial *FirewallPolicy, opts ...dcl.ApplyOption) (*FirewallPolicy, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &FirewallPolicy{}
	if dcl.IsZeroValue(rawDesired.Name) || (dcl.IsEmptyValueIndirect(rawDesired.Name) && dcl.IsEmptyValueIndirect(rawInitial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.ShortName, rawInitial.ShortName) {
		canonicalDesired.ShortName = rawInitial.ShortName
	} else {
		canonicalDesired.ShortName = rawDesired.ShortName
	}
	if dcl.StringCanonicalize(rawDesired.Parent, rawInitial.Parent) {
		canonicalDesired.Parent = rawInitial.Parent
	} else {
		canonicalDesired.Parent = rawDesired.Parent
	}
	return canonicalDesired, nil
}

func canonicalizeFirewallPolicyNewState(c *Client, rawNew, rawDesired *FirewallPolicy) (*FirewallPolicy, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
		if dcl.StringCanonicalize(rawDesired.Id, rawNew.Id) {
			rawNew.Id = rawDesired.Id
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreationTimestamp) && dcl.IsEmptyValueIndirect(rawDesired.CreationTimestamp) {
		rawNew.CreationTimestamp = rawDesired.CreationTimestamp
	} else {
		if dcl.StringCanonicalize(rawDesired.CreationTimestamp, rawNew.CreationTimestamp) {
			rawNew.CreationTimestamp = rawDesired.CreationTimestamp
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

	if dcl.IsEmptyValueIndirect(rawNew.ShortName) && dcl.IsEmptyValueIndirect(rawDesired.ShortName) {
		rawNew.ShortName = rawDesired.ShortName
	} else {
		if dcl.StringCanonicalize(rawDesired.ShortName, rawNew.ShortName) {
			rawNew.ShortName = rawDesired.ShortName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Parent) && dcl.IsEmptyValueIndirect(rawDesired.Parent) {
		rawNew.Parent = rawDesired.Parent
	} else {
		if dcl.StringCanonicalize(rawDesired.Parent, rawNew.Parent) {
			rawNew.Parent = rawDesired.Parent
		}
	}

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffFirewallPolicy(c *Client, desired, actual *FirewallPolicy, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateFirewallPolicyPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.ShortName, actual.ShortName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ShortName")); len(ds) != 0 || err != nil {
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

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *FirewallPolicy) urlNormalized() *FirewallPolicy {
	normalized := dcl.Copy(*r).(FirewallPolicy)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Id = dcl.SelfLinkToName(r.Id)
	normalized.CreationTimestamp = dcl.SelfLinkToName(r.CreationTimestamp)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Fingerprint = dcl.SelfLinkToName(r.Fingerprint)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.SelfLinkWithId = dcl.SelfLinkToName(r.SelfLinkWithId)
	normalized.ShortName = dcl.SelfLinkToName(r.ShortName)
	normalized.Parent = r.Parent
	return &normalized
}

func (r *FirewallPolicy) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"name": dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("locations/global/firewallPolicies/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the FirewallPolicy resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *FirewallPolicy) marshal(c *Client) ([]byte, error) {
	m, err := expandFirewallPolicy(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling FirewallPolicy: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalFirewallPolicy decodes JSON responses into the FirewallPolicy resource schema.
func unmarshalFirewallPolicy(b []byte, c *Client, res *FirewallPolicy) (*FirewallPolicy, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapFirewallPolicy(m, c, res)
}

func unmarshalMapFirewallPolicy(m map[string]interface{}, c *Client, res *FirewallPolicy) (*FirewallPolicy, error) {

	flattened := flattenFirewallPolicy(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandFirewallPolicy expands FirewallPolicy into a JSON request object.
func expandFirewallPolicy(c *Client, f *FirewallPolicy) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.ShortName; dcl.ValueShouldBeSent(v) {
		m["shortName"] = v
	}
	if v := f.Parent; dcl.ValueShouldBeSent(v) {
		m["parent"] = v
	}

	return m, nil
}

// flattenFirewallPolicy flattens FirewallPolicy from a JSON request object into the
// FirewallPolicy type.
func flattenFirewallPolicy(c *Client, i interface{}, res *FirewallPolicy) *FirewallPolicy {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &FirewallPolicy{}
	resultRes.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))
	resultRes.Id = dcl.FlattenString(m["id"])
	resultRes.CreationTimestamp = dcl.FlattenString(m["creationTimestamp"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Fingerprint = dcl.FlattenString(m["fingerprint"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.SelfLinkWithId = dcl.FlattenString(m["selfLinkWithId"])
	resultRes.RuleTupleCount = dcl.FlattenInteger(m["ruleTupleCount"])
	resultRes.ShortName = dcl.FlattenString(m["shortName"])
	resultRes.Parent = dcl.FlattenString(m["parent"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *FirewallPolicy) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalFirewallPolicy(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

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

type firewallPolicyDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         firewallPolicyApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToFirewallPolicyDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]firewallPolicyDiff, error) {
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
	var diffs []firewallPolicyDiff
	// For each operation name, create a firewallPolicyDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := firewallPolicyDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToFirewallPolicyApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToFirewallPolicyApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (firewallPolicyApiOperation, error) {
	switch opName {

	case "updateFirewallPolicyPatchOperation":
		return &updateFirewallPolicyPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractFirewallPolicyFields(r *FirewallPolicy) error {
	return nil
}

func postReadExtractFirewallPolicyFields(r *FirewallPolicy) error {
	return nil
}
