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

func (r *Network) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.RoutingConfig) {
		if err := r.RoutingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *NetworkRoutingConfig) validate() error {
	return nil
}
func (r *Network) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *Network) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/networks/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Network) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/networks", nr.basePath(), userBasePath, params), nil

}

func (r *Network) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/networks", nr.basePath(), userBasePath, params), nil

}

func (r *Network) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/networks/{{name}}", nr.basePath(), userBasePath, params), nil
}

// networkApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type networkApiOperation interface {
	do(context.Context, *Network, *Client) error
}

// newUpdateNetworkUpdateRequest creates a request for an
// Network resource's update update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateNetworkUpdateRequest(ctx context.Context, f *Network, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandNetworkRoutingConfig(c, f.RoutingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding RoutingConfig into routingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["routingConfig"] = v
	}
	if v := f.Mtu; !dcl.IsEmptyValueIndirect(v) {
		req["mtu"] = v
	}
	return req, nil
}

// marshalUpdateNetworkUpdateRequest converts the update into
// the final JSON request body.
func marshalUpdateNetworkUpdateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateNetworkUpdateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (c *Client) listNetworkRaw(ctx context.Context, r *Network, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != NetworkMaxPage {
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

type listNetworkOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listNetwork(ctx context.Context, r *Network, pageToken string, pageSize int32) ([]*Network, string, error) {
	b, err := c.listNetworkRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listNetworkOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Network
	for _, v := range m.Items {
		res, err := unmarshalMapNetwork(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllNetwork(ctx context.Context, f func(*Network) bool, resources []*Network) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteNetwork(ctx, res)
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

type deleteNetworkOperation struct{}

func (op *deleteNetworkOperation) do(ctx context.Context, r *Network, c *Client) error {
	r, err := c.GetNetwork(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Network not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetNetwork checking for existence. error: %v", err)
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
		_, err := c.GetNetwork(ctx, r)
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
type createNetworkOperation struct {
	response map[string]interface{}
}

func (op *createNetworkOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createNetworkOperation) do(ctx context.Context, r *Network, c *Client) error {
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

	if _, err := c.GetNetwork(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getNetworkRaw(ctx context.Context, r *Network) ([]byte, error) {
	if dcl.IsZeroValue(r.AutoCreateSubnetworks) {
		r.AutoCreateSubnetworks = dcl.Bool(true)
	}

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

func (c *Client) networkDiffsForRawDesired(ctx context.Context, rawDesired *Network, opts ...dcl.ApplyOption) (initial, desired *Network, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Network
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Network); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Network, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetNetwork(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Network resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Network resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Network resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeNetworkDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Network: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Network: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractNetworkFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeNetworkInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Network: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeNetworkDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Network: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffNetwork(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeNetworkInitialState(rawInitial, rawDesired *Network) (*Network, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeNetworkDesiredState(rawDesired, rawInitial *Network, opts ...dcl.ApplyOption) (*Network, error) {

	if dcl.IsZeroValue(rawDesired.AutoCreateSubnetworks) {
		rawDesired.AutoCreateSubnetworks = dcl.Bool(true)
	}

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.RoutingConfig = canonicalizeNetworkRoutingConfig(rawDesired.RoutingConfig, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Network{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.BoolCanonicalize(rawDesired.AutoCreateSubnetworks, rawInitial.AutoCreateSubnetworks) {
		canonicalDesired.AutoCreateSubnetworks = rawInitial.AutoCreateSubnetworks
	} else {
		canonicalDesired.AutoCreateSubnetworks = rawDesired.AutoCreateSubnetworks
	}
	canonicalDesired.RoutingConfig = canonicalizeNetworkRoutingConfig(rawDesired.RoutingConfig, rawInitial.RoutingConfig, opts...)
	if dcl.IsZeroValue(rawDesired.Mtu) || (dcl.IsEmptyValueIndirect(rawDesired.Mtu) && dcl.IsEmptyValueIndirect(rawInitial.Mtu)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Mtu = rawInitial.Mtu
	} else {
		canonicalDesired.Mtu = rawDesired.Mtu
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeNetworkNewState(c *Client, rawNew, rawDesired *Network) (*Network, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.GatewayIPv4) && dcl.IsEmptyValueIndirect(rawDesired.GatewayIPv4) {
		rawNew.GatewayIPv4 = rawDesired.GatewayIPv4
	} else {
		if dcl.StringCanonicalize(rawDesired.GatewayIPv4, rawNew.GatewayIPv4) {
			rawNew.GatewayIPv4 = rawDesired.GatewayIPv4
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.AutoCreateSubnetworks) && dcl.IsEmptyValueIndirect(rawDesired.AutoCreateSubnetworks) {
		rawNew.AutoCreateSubnetworks = rawDesired.AutoCreateSubnetworks
	} else {
		if dcl.BoolCanonicalize(rawDesired.AutoCreateSubnetworks, rawNew.AutoCreateSubnetworks) {
			rawNew.AutoCreateSubnetworks = rawDesired.AutoCreateSubnetworks
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RoutingConfig) && dcl.IsEmptyValueIndirect(rawDesired.RoutingConfig) {
		rawNew.RoutingConfig = rawDesired.RoutingConfig
	} else {
		rawNew.RoutingConfig = canonicalizeNewNetworkRoutingConfig(c, rawDesired.RoutingConfig, rawNew.RoutingConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Mtu) && dcl.IsEmptyValueIndirect(rawDesired.Mtu) {
		rawNew.Mtu = rawDesired.Mtu
	} else {
	}

	rawNew.Project = rawDesired.Project

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

	return rawNew, nil
}

func canonicalizeNetworkRoutingConfig(des, initial *NetworkRoutingConfig, opts ...dcl.ApplyOption) *NetworkRoutingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NetworkRoutingConfig{}

	if dcl.IsZeroValue(des.RoutingMode) || (dcl.IsEmptyValueIndirect(des.RoutingMode) && dcl.IsEmptyValueIndirect(initial.RoutingMode)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.RoutingMode = initial.RoutingMode
	} else {
		cDes.RoutingMode = des.RoutingMode
	}

	return cDes
}

func canonicalizeNetworkRoutingConfigSlice(des, initial []NetworkRoutingConfig, opts ...dcl.ApplyOption) []NetworkRoutingConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NetworkRoutingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNetworkRoutingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NetworkRoutingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNetworkRoutingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNetworkRoutingConfig(c *Client, des, nw *NetworkRoutingConfig) *NetworkRoutingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NetworkRoutingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNetworkRoutingConfigSet(c *Client, des, nw []NetworkRoutingConfig) []NetworkRoutingConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NetworkRoutingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNetworkRoutingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNetworkRoutingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNetworkRoutingConfigSlice(c *Client, des, nw []NetworkRoutingConfig) []NetworkRoutingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NetworkRoutingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNetworkRoutingConfig(c, &d, &n))
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
func diffNetwork(c *Client, desired, actual *Network, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GatewayIPv4, actual.GatewayIPv4, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GatewayIPv4")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.AutoCreateSubnetworks, actual.AutoCreateSubnetworks, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoCreateSubnetworks")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RoutingConfig, actual.RoutingConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareNetworkRoutingConfigNewStyle, EmptyObject: EmptyNetworkRoutingConfig, OperationSelector: dcl.TriggersOperation("updateNetworkUpdateOperation")}, fn.AddNest("RoutingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Mtu, actual.Mtu, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNetworkUpdateOperation")}, fn.AddNest("Mtu")); len(ds) != 0 || err != nil {
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

	if len(newDiffs) > 0 {
		c.Config.Logger.Infof("Diff function found diffs: %v", newDiffs)
	}
	return newDiffs, nil
}
func compareNetworkRoutingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NetworkRoutingConfig)
	if !ok {
		desiredNotPointer, ok := d.(NetworkRoutingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkRoutingConfig or *NetworkRoutingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NetworkRoutingConfig)
	if !ok {
		actualNotPointer, ok := a.(NetworkRoutingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NetworkRoutingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RoutingMode, actual.RoutingMode, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateNetworkUpdateOperation")}, fn.AddNest("RoutingMode")); len(ds) != 0 || err != nil {
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
func (r *Network) urlNormalized() *Network {
	normalized := dcl.Copy(*r).(Network)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.GatewayIPv4 = dcl.SelfLinkToName(r.GatewayIPv4)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.SelfLinkWithId = dcl.SelfLinkToName(r.SelfLinkWithId)
	return &normalized
}

func (r *Network) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "update" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/global/networks/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Network resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Network) marshal(c *Client) ([]byte, error) {
	m, err := expandNetwork(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Network: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalNetwork decodes JSON responses into the Network resource schema.
func unmarshalNetwork(b []byte, c *Client, res *Network) (*Network, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNetwork(m, c, res)
}

func unmarshalMapNetwork(m map[string]interface{}, c *Client, res *Network) (*Network, error) {

	flattened := flattenNetwork(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNetwork expands Network into a JSON request object.
func expandNetwork(c *Client, f *Network) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.AutoCreateSubnetworks; v != nil {
		m["autoCreateSubnetworks"] = v
	}
	if v, err := expandNetworkRoutingConfig(c, f.RoutingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding RoutingConfig into routingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["routingConfig"] = v
	}
	if v := f.Mtu; dcl.ValueShouldBeSent(v) {
		m["mtu"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenNetwork flattens Network from a JSON request object into the
// Network type.
func flattenNetwork(c *Client, i interface{}, res *Network) *Network {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Network{}
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.GatewayIPv4 = dcl.FlattenString(m["gatewayIPv4"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.AutoCreateSubnetworks = dcl.FlattenBool(m["autoCreateSubnetworks"])
	if _, ok := m["autoCreateSubnetworks"]; !ok {
		c.Config.Logger.Info("Using default value for autoCreateSubnetworks")
		resultRes.AutoCreateSubnetworks = dcl.Bool(true)
	}
	resultRes.RoutingConfig = flattenNetworkRoutingConfig(c, m["routingConfig"], res)
	resultRes.Mtu = dcl.FlattenInteger(m["mtu"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.SelfLinkWithId = flattenNetworkSelfLinkWithID(c, m["selfLinkWithId"], res, m)

	return resultRes
}

// expandNetworkRoutingConfigMap expands the contents of NetworkRoutingConfig into a JSON
// request object.
func expandNetworkRoutingConfigMap(c *Client, f map[string]NetworkRoutingConfig, res *Network) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNetworkRoutingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNetworkRoutingConfigSlice expands the contents of NetworkRoutingConfig into a JSON
// request object.
func expandNetworkRoutingConfigSlice(c *Client, f []NetworkRoutingConfig, res *Network) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNetworkRoutingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNetworkRoutingConfigMap flattens the contents of NetworkRoutingConfig from a JSON
// response object.
func flattenNetworkRoutingConfigMap(c *Client, i interface{}, res *Network) map[string]NetworkRoutingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkRoutingConfig{}
	}

	if len(a) == 0 {
		return map[string]NetworkRoutingConfig{}
	}

	items := make(map[string]NetworkRoutingConfig)
	for k, item := range a {
		items[k] = *flattenNetworkRoutingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNetworkRoutingConfigSlice flattens the contents of NetworkRoutingConfig from a JSON
// response object.
func flattenNetworkRoutingConfigSlice(c *Client, i interface{}, res *Network) []NetworkRoutingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkRoutingConfig{}
	}

	if len(a) == 0 {
		return []NetworkRoutingConfig{}
	}

	items := make([]NetworkRoutingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkRoutingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNetworkRoutingConfig expands an instance of NetworkRoutingConfig into a JSON
// request object.
func expandNetworkRoutingConfig(c *Client, f *NetworkRoutingConfig, res *Network) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RoutingMode; !dcl.IsEmptyValueIndirect(v) {
		m["routingMode"] = v
	}

	return m, nil
}

// flattenNetworkRoutingConfig flattens an instance of NetworkRoutingConfig from a JSON
// response object.
func flattenNetworkRoutingConfig(c *Client, i interface{}, res *Network) *NetworkRoutingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NetworkRoutingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNetworkRoutingConfig
	}
	r.RoutingMode = flattenNetworkRoutingConfigRoutingModeEnum(m["routingMode"])

	return r
}

// flattenNetworkRoutingConfigRoutingModeEnumMap flattens the contents of NetworkRoutingConfigRoutingModeEnum from a JSON
// response object.
func flattenNetworkRoutingConfigRoutingModeEnumMap(c *Client, i interface{}, res *Network) map[string]NetworkRoutingConfigRoutingModeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NetworkRoutingConfigRoutingModeEnum{}
	}

	if len(a) == 0 {
		return map[string]NetworkRoutingConfigRoutingModeEnum{}
	}

	items := make(map[string]NetworkRoutingConfigRoutingModeEnum)
	for k, item := range a {
		items[k] = *flattenNetworkRoutingConfigRoutingModeEnum(item.(interface{}))
	}

	return items
}

// flattenNetworkRoutingConfigRoutingModeEnumSlice flattens the contents of NetworkRoutingConfigRoutingModeEnum from a JSON
// response object.
func flattenNetworkRoutingConfigRoutingModeEnumSlice(c *Client, i interface{}, res *Network) []NetworkRoutingConfigRoutingModeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NetworkRoutingConfigRoutingModeEnum{}
	}

	if len(a) == 0 {
		return []NetworkRoutingConfigRoutingModeEnum{}
	}

	items := make([]NetworkRoutingConfigRoutingModeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNetworkRoutingConfigRoutingModeEnum(item.(interface{})))
	}

	return items
}

// flattenNetworkRoutingConfigRoutingModeEnum asserts that an interface is a string, and returns a
// pointer to a *NetworkRoutingConfigRoutingModeEnum with the same value as that string.
func flattenNetworkRoutingConfigRoutingModeEnum(i interface{}) *NetworkRoutingConfigRoutingModeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return NetworkRoutingConfigRoutingModeEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Network) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNetwork(b, c, r)
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

type networkDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         networkApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToNetworkDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]networkDiff, error) {
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
	var diffs []networkDiff
	// For each operation name, create a networkDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := networkDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToNetworkApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToNetworkApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (networkApiOperation, error) {
	switch opName {

	case "updateNetworkUpdateOperation":
		return &updateNetworkUpdateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractNetworkFields(r *Network) error {
	vRoutingConfig := r.RoutingConfig
	if vRoutingConfig == nil {
		// note: explicitly not the empty object.
		vRoutingConfig = &NetworkRoutingConfig{}
	}
	if err := extractNetworkRoutingConfigFields(r, vRoutingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRoutingConfig) {
		r.RoutingConfig = vRoutingConfig
	}
	return nil
}
func extractNetworkRoutingConfigFields(r *Network, o *NetworkRoutingConfig) error {
	return nil
}

func postReadExtractNetworkFields(r *Network) error {
	vRoutingConfig := r.RoutingConfig
	if vRoutingConfig == nil {
		// note: explicitly not the empty object.
		vRoutingConfig = &NetworkRoutingConfig{}
	}
	if err := postReadExtractNetworkRoutingConfigFields(r, vRoutingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRoutingConfig) {
		r.RoutingConfig = vRoutingConfig
	}
	return nil
}
func postReadExtractNetworkRoutingConfigFields(r *Network, o *NetworkRoutingConfig) error {
	return nil
}
