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

func (r *Route) validate() error {

	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"NextHopVpnTunnel", "NextHopIP", "NextHopInstance", "NextHopGateway", "NextHopIlb"}, r.NextHopVpnTunnel, r.NextHopIP, r.NextHopInstance, r.NextHopGateway, r.NextHopIlb); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "network"); err != nil {
		return err
	}
	if err := dcl.Required(r, "destRange"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *RouteWarning) validate() error {
	return nil
}
func (r *Route) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *Route) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/routes/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Route) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/routes", nr.basePath(), userBasePath, params), nil

}

func (r *Route) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/routes", nr.basePath(), userBasePath, params), nil

}

func (r *Route) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/routes/{{name}}", nr.basePath(), userBasePath, params), nil
}

// routeApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type routeApiOperation interface {
	do(context.Context, *Route, *Client) error
}

func (c *Client) listRouteRaw(ctx context.Context, r *Route, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != RouteMaxPage {
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

type listRouteOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listRoute(ctx context.Context, r *Route, pageToken string, pageSize int32) ([]*Route, string, error) {
	b, err := c.listRouteRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listRouteOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Route
	for _, v := range m.Items {
		res, err := unmarshalMapRoute(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllRoute(ctx context.Context, f func(*Route) bool, resources []*Route) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteRoute(ctx, res)
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

type deleteRouteOperation struct{}

func (op *deleteRouteOperation) do(ctx context.Context, r *Route, c *Client) error {
	r, err := c.GetRoute(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Route not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetRoute checking for existence. error: %v", err)
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
		_, err := c.GetRoute(ctx, r)
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
type createRouteOperation struct {
	response map[string]interface{}
}

func (op *createRouteOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createRouteOperation) do(ctx context.Context, r *Route, c *Client) error {
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

	if _, err := c.GetRoute(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getRouteRaw(ctx context.Context, r *Route) ([]byte, error) {
	if dcl.IsZeroValue(r.Priority) {
		r.Priority = dcl.Int64(1000)
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

func (c *Client) routeDiffsForRawDesired(ctx context.Context, rawDesired *Route, opts ...dcl.ApplyOption) (initial, desired *Route, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Route
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Route); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Route, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetRoute(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Route resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Route resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Route resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeRouteDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Route: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Route: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractRouteFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeRouteInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Route: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeRouteDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Route: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffRoute(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeRouteInitialState(rawInitial, rawDesired *Route) (*Route, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.NextHopVpnTunnel) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NextHopIP, rawInitial.NextHopInstance, rawInitial.NextHopGateway, rawInitial.NextHopIlb) {
			rawInitial.NextHopVpnTunnel = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.NextHopIP) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NextHopVpnTunnel, rawInitial.NextHopInstance, rawInitial.NextHopGateway, rawInitial.NextHopIlb) {
			rawInitial.NextHopIP = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.NextHopInstance) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NextHopVpnTunnel, rawInitial.NextHopIP, rawInitial.NextHopGateway, rawInitial.NextHopIlb) {
			rawInitial.NextHopInstance = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.NextHopGateway) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NextHopVpnTunnel, rawInitial.NextHopIP, rawInitial.NextHopInstance, rawInitial.NextHopIlb) {
			rawInitial.NextHopGateway = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.NextHopIlb) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NextHopVpnTunnel, rawInitial.NextHopIP, rawInitial.NextHopInstance, rawInitial.NextHopGateway) {
			rawInitial.NextHopIlb = dcl.String("")
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

func canonicalizeRouteDesiredState(rawDesired, rawInitial *Route, opts ...dcl.ApplyOption) (*Route, error) {

	if dcl.IsZeroValue(rawDesired.Priority) {
		rawDesired.Priority = dcl.Int64(1000)
	}

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Route{}
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
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Network, rawInitial.Network) {
		canonicalDesired.Network = rawInitial.Network
	} else {
		canonicalDesired.Network = rawDesired.Network
	}
	if dcl.StringArrayCanonicalize(rawDesired.Tag, rawInitial.Tag) {
		canonicalDesired.Tag = rawInitial.Tag
	} else {
		canonicalDesired.Tag = rawDesired.Tag
	}
	if dcl.StringCanonicalize(rawDesired.DestRange, rawInitial.DestRange) {
		canonicalDesired.DestRange = rawInitial.DestRange
	} else {
		canonicalDesired.DestRange = rawDesired.DestRange
	}
	if dcl.IsZeroValue(rawDesired.Priority) || (dcl.IsEmptyValueIndirect(rawDesired.Priority) && dcl.IsEmptyValueIndirect(rawInitial.Priority)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Priority = rawInitial.Priority
	} else {
		canonicalDesired.Priority = rawDesired.Priority
	}
	if dcl.StringCanonicalize(rawDesired.NextHopInstance, rawInitial.NextHopInstance) {
		canonicalDesired.NextHopInstance = rawInitial.NextHopInstance
	} else {
		canonicalDesired.NextHopInstance = rawDesired.NextHopInstance
	}
	if dcl.StringCanonicalize(rawDesired.NextHopIP, rawInitial.NextHopIP) {
		canonicalDesired.NextHopIP = rawInitial.NextHopIP
	} else {
		canonicalDesired.NextHopIP = rawDesired.NextHopIP
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.NextHopGateway, rawInitial.NextHopGateway) {
		canonicalDesired.NextHopGateway = rawInitial.NextHopGateway
	} else {
		canonicalDesired.NextHopGateway = rawDesired.NextHopGateway
	}
	if dcl.StringCanonicalize(rawDesired.NextHopIlb, rawInitial.NextHopIlb) {
		canonicalDesired.NextHopIlb = rawInitial.NextHopIlb
	} else {
		canonicalDesired.NextHopIlb = rawDesired.NextHopIlb
	}
	if dcl.StringCanonicalize(rawDesired.NextHopVpnTunnel, rawInitial.NextHopVpnTunnel) {
		canonicalDesired.NextHopVpnTunnel = rawInitial.NextHopVpnTunnel
	} else {
		canonicalDesired.NextHopVpnTunnel = rawDesired.NextHopVpnTunnel
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}

	if canonicalDesired.NextHopVpnTunnel != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NextHopIP, rawDesired.NextHopInstance, rawDesired.NextHopGateway, rawDesired.NextHopIlb) {
			canonicalDesired.NextHopVpnTunnel = dcl.String("")
		}
	}

	if canonicalDesired.NextHopIP != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NextHopVpnTunnel, rawDesired.NextHopInstance, rawDesired.NextHopGateway, rawDesired.NextHopIlb) {
			canonicalDesired.NextHopIP = dcl.String("")
		}
	}

	if canonicalDesired.NextHopInstance != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NextHopVpnTunnel, rawDesired.NextHopIP, rawDesired.NextHopGateway, rawDesired.NextHopIlb) {
			canonicalDesired.NextHopInstance = dcl.String("")
		}
	}

	if canonicalDesired.NextHopGateway != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NextHopVpnTunnel, rawDesired.NextHopIP, rawDesired.NextHopInstance, rawDesired.NextHopIlb) {
			canonicalDesired.NextHopGateway = dcl.String("")
		}
	}

	if canonicalDesired.NextHopIlb != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NextHopVpnTunnel, rawDesired.NextHopIP, rawDesired.NextHopInstance, rawDesired.NextHopGateway) {
			canonicalDesired.NextHopIlb = dcl.String("")
		}
	}

	return canonicalDesired, nil
}

func canonicalizeRouteNewState(c *Client, rawNew, rawDesired *Route) (*Route, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Network) && dcl.IsEmptyValueIndirect(rawDesired.Network) {
		rawNew.Network = rawDesired.Network
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Network, rawNew.Network) {
			rawNew.Network = rawDesired.Network
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Tag) && dcl.IsEmptyValueIndirect(rawDesired.Tag) {
		rawNew.Tag = rawDesired.Tag
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.Tag, rawNew.Tag) {
			rawNew.Tag = rawDesired.Tag
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DestRange) && dcl.IsEmptyValueIndirect(rawDesired.DestRange) {
		rawNew.DestRange = rawDesired.DestRange
	} else {
		if dcl.StringCanonicalize(rawDesired.DestRange, rawNew.DestRange) {
			rawNew.DestRange = rawDesired.DestRange
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Priority) && dcl.IsEmptyValueIndirect(rawDesired.Priority) {
		rawNew.Priority = rawDesired.Priority
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopInstance) && dcl.IsEmptyValueIndirect(rawDesired.NextHopInstance) {
		rawNew.NextHopInstance = rawDesired.NextHopInstance
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopInstance, rawNew.NextHopInstance) {
			rawNew.NextHopInstance = rawDesired.NextHopInstance
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopIP) && dcl.IsEmptyValueIndirect(rawDesired.NextHopIP) {
		rawNew.NextHopIP = rawDesired.NextHopIP
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopIP, rawNew.NextHopIP) {
			rawNew.NextHopIP = rawDesired.NextHopIP
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopNetwork) && dcl.IsEmptyValueIndirect(rawDesired.NextHopNetwork) {
		rawNew.NextHopNetwork = rawDesired.NextHopNetwork
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopNetwork, rawNew.NextHopNetwork) {
			rawNew.NextHopNetwork = rawDesired.NextHopNetwork
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopGateway) && dcl.IsEmptyValueIndirect(rawDesired.NextHopGateway) {
		rawNew.NextHopGateway = rawDesired.NextHopGateway
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.NextHopGateway, rawNew.NextHopGateway) {
			rawNew.NextHopGateway = rawDesired.NextHopGateway
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopPeering) && dcl.IsEmptyValueIndirect(rawDesired.NextHopPeering) {
		rawNew.NextHopPeering = rawDesired.NextHopPeering
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopPeering, rawNew.NextHopPeering) {
			rawNew.NextHopPeering = rawDesired.NextHopPeering
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopIlb) && dcl.IsEmptyValueIndirect(rawDesired.NextHopIlb) {
		rawNew.NextHopIlb = rawDesired.NextHopIlb
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopIlb, rawNew.NextHopIlb) {
			rawNew.NextHopIlb = rawDesired.NextHopIlb
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Warning) && dcl.IsEmptyValueIndirect(rawDesired.Warning) {
		rawNew.Warning = rawDesired.Warning
	} else {
		rawNew.Warning = canonicalizeNewRouteWarningSlice(c, rawDesired.Warning, rawNew.Warning)
	}

	if dcl.IsEmptyValueIndirect(rawNew.NextHopVpnTunnel) && dcl.IsEmptyValueIndirect(rawDesired.NextHopVpnTunnel) {
		rawNew.NextHopVpnTunnel = rawDesired.NextHopVpnTunnel
	} else {
		if dcl.StringCanonicalize(rawDesired.NextHopVpnTunnel, rawNew.NextHopVpnTunnel) {
			rawNew.NextHopVpnTunnel = rawDesired.NextHopVpnTunnel
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeRouteWarning(des, initial *RouteWarning, opts ...dcl.ApplyOption) *RouteWarning {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &RouteWarning{}

	return cDes
}

func canonicalizeRouteWarningSlice(des, initial []RouteWarning, opts ...dcl.ApplyOption) []RouteWarning {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]RouteWarning, 0, len(des))
		for _, d := range des {
			cd := canonicalizeRouteWarning(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]RouteWarning, 0, len(des))
	for i, d := range des {
		cd := canonicalizeRouteWarning(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewRouteWarning(c *Client, des, nw *RouteWarning) *RouteWarning {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for RouteWarning while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Message, nw.Message) {
		nw.Message = des.Message
	}

	return nw
}

func canonicalizeNewRouteWarningSet(c *Client, des, nw []RouteWarning) []RouteWarning {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []RouteWarning
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareRouteWarningNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewRouteWarning(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewRouteWarningSlice(c *Client, des, nw []RouteWarning) []RouteWarning {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []RouteWarning
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewRouteWarning(c, &d, &n))
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
func diffRoute(c *Client, desired, actual *Route, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Network")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tag, actual.Tag, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DestRange, actual.DestRange, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DestRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Priority, actual.Priority, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Priority")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopInstance, actual.NextHopInstance, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopInstance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopIP, actual.NextHopIP, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopIp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopNetwork, actual.NextHopNetwork, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopNetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopGateway, actual.NextHopGateway, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopGateway")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopPeering, actual.NextHopPeering, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopPeering")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopIlb, actual.NextHopIlb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopIlb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Warning, actual.Warning, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareRouteWarningNewStyle, EmptyObject: EmptyRouteWarning, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Warnings")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextHopVpnTunnel, actual.NextHopVpnTunnel, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextHopVpnTunnel")); len(ds) != 0 || err != nil {
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
func compareRouteWarningNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*RouteWarning)
	if !ok {
		desiredNotPointer, ok := d.(RouteWarning)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RouteWarning or *RouteWarning", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*RouteWarning)
	if !ok {
		actualNotPointer, ok := a.(RouteWarning)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a RouteWarning", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Code, actual.Code, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Code")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Message, actual.Message, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Message")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Data, actual.Data, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Data")); len(ds) != 0 || err != nil {
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
func (r *Route) urlNormalized() *Route {
	normalized := dcl.Copy(*r).(Route)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Network = dcl.SelfLinkToName(r.Network)
	normalized.DestRange = dcl.SelfLinkToName(r.DestRange)
	normalized.NextHopInstance = dcl.SelfLinkToName(r.NextHopInstance)
	normalized.NextHopIP = dcl.SelfLinkToName(r.NextHopIP)
	normalized.NextHopNetwork = dcl.SelfLinkToName(r.NextHopNetwork)
	normalized.NextHopGateway = dcl.SelfLinkToName(r.NextHopGateway)
	normalized.NextHopPeering = dcl.SelfLinkToName(r.NextHopPeering)
	normalized.NextHopIlb = dcl.SelfLinkToName(r.NextHopIlb)
	normalized.NextHopVpnTunnel = dcl.SelfLinkToName(r.NextHopVpnTunnel)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Route) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Route resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Route) marshal(c *Client) ([]byte, error) {
	m, err := expandRoute(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Route: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalRoute decodes JSON responses into the Route resource schema.
func unmarshalRoute(b []byte, c *Client, res *Route) (*Route, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapRoute(m, c, res)
}

func unmarshalMapRoute(m map[string]interface{}, c *Client, res *Route) (*Route, error) {
	if v, err := dcl.MapFromListOfKeyValues(m, []string{"warnings", "data", "items"}, "key", "value"); err != nil {
		return nil, err
	} else {
		dcl.PutMapEntry(
			m,
			[]string{"warnings", "data"},
			v,
		)
	}

	flattened := flattenRoute(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandRoute expands Route into a JSON request object.
func expandRoute(c *Client, f *Route) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := dcl.DeriveField("global/networks/%s", f.Network, dcl.SelfLinkToName(f.Network)); err != nil {
		return nil, fmt.Errorf("error expanding Network into network: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["network"] = v
	}
	if v := f.Tag; v != nil {
		m["tags"] = v
	}
	if v := f.DestRange; dcl.ValueShouldBeSent(v) {
		m["destRange"] = v
	}
	if v := f.Priority; dcl.ValueShouldBeSent(v) {
		m["priority"] = v
	}
	if v := f.NextHopInstance; dcl.ValueShouldBeSent(v) {
		m["nextHopInstance"] = v
	}
	if v := f.NextHopIP; dcl.ValueShouldBeSent(v) {
		m["nextHopIp"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/global/gateways/%s", f.NextHopGateway, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.NextHopGateway)); err != nil {
		return nil, fmt.Errorf("error expanding NextHopGateway into nextHopGateway: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["nextHopGateway"] = v
	}
	if v := f.NextHopIlb; dcl.ValueShouldBeSent(v) {
		m["nextHopIlb"] = v
	}
	if v := f.NextHopVpnTunnel; dcl.ValueShouldBeSent(v) {
		m["nextHopVpnTunnel"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenRoute flattens Route from a JSON request object into the
// Route type.
func flattenRoute(c *Client, i interface{}, res *Route) *Route {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Route{}
	resultRes.Id = dcl.FlattenInteger(m["id"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Network = dcl.FlattenString(m["network"])
	resultRes.Tag = dcl.FlattenStringSlice(m["tags"])
	resultRes.DestRange = dcl.FlattenString(m["destRange"])
	resultRes.Priority = dcl.FlattenInteger(m["priority"])
	if _, ok := m["priority"]; !ok {
		c.Config.Logger.Info("Using default value for priority")
		resultRes.Priority = dcl.Int64(1000)
	}
	resultRes.NextHopInstance = dcl.FlattenString(m["nextHopInstance"])
	resultRes.NextHopIP = dcl.FlattenString(m["nextHopIp"])
	resultRes.NextHopNetwork = dcl.FlattenString(m["nextHopNetwork"])
	resultRes.NextHopGateway = dcl.FlattenString(m["nextHopGateway"])
	resultRes.NextHopPeering = dcl.FlattenString(m["nextHopPeering"])
	resultRes.NextHopIlb = dcl.FlattenString(m["nextHopIlb"])
	resultRes.Warning = flattenRouteWarningSlice(c, m["warnings"], res)
	resultRes.NextHopVpnTunnel = dcl.FlattenString(m["nextHopVpnTunnel"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandRouteWarningMap expands the contents of RouteWarning into a JSON
// request object.
func expandRouteWarningMap(c *Client, f map[string]RouteWarning, res *Route) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandRouteWarning(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandRouteWarningSlice expands the contents of RouteWarning into a JSON
// request object.
func expandRouteWarningSlice(c *Client, f []RouteWarning, res *Route) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandRouteWarning(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenRouteWarningMap flattens the contents of RouteWarning from a JSON
// response object.
func flattenRouteWarningMap(c *Client, i interface{}, res *Route) map[string]RouteWarning {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RouteWarning{}
	}

	if len(a) == 0 {
		return map[string]RouteWarning{}
	}

	items := make(map[string]RouteWarning)
	for k, item := range a {
		items[k] = *flattenRouteWarning(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenRouteWarningSlice flattens the contents of RouteWarning from a JSON
// response object.
func flattenRouteWarningSlice(c *Client, i interface{}, res *Route) []RouteWarning {
	a, ok := i.([]interface{})
	if !ok {
		return []RouteWarning{}
	}

	if len(a) == 0 {
		return []RouteWarning{}
	}

	items := make([]RouteWarning, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRouteWarning(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandRouteWarning expands an instance of RouteWarning into a JSON
// request object.
func expandRouteWarning(c *Client, f *RouteWarning, res *Route) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenRouteWarning flattens an instance of RouteWarning from a JSON
// response object.
func flattenRouteWarning(c *Client, i interface{}, res *Route) *RouteWarning {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &RouteWarning{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyRouteWarning
	}
	r.Code = flattenRouteWarningCodeEnum(m["code"])
	r.Message = dcl.FlattenString(m["message"])
	r.Data = dcl.FlattenKeyValuePairs(m["data"])

	return r
}

// flattenRouteWarningCodeEnumMap flattens the contents of RouteWarningCodeEnum from a JSON
// response object.
func flattenRouteWarningCodeEnumMap(c *Client, i interface{}, res *Route) map[string]RouteWarningCodeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]RouteWarningCodeEnum{}
	}

	if len(a) == 0 {
		return map[string]RouteWarningCodeEnum{}
	}

	items := make(map[string]RouteWarningCodeEnum)
	for k, item := range a {
		items[k] = *flattenRouteWarningCodeEnum(item.(interface{}))
	}

	return items
}

// flattenRouteWarningCodeEnumSlice flattens the contents of RouteWarningCodeEnum from a JSON
// response object.
func flattenRouteWarningCodeEnumSlice(c *Client, i interface{}, res *Route) []RouteWarningCodeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []RouteWarningCodeEnum{}
	}

	if len(a) == 0 {
		return []RouteWarningCodeEnum{}
	}

	items := make([]RouteWarningCodeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenRouteWarningCodeEnum(item.(interface{})))
	}

	return items
}

// flattenRouteWarningCodeEnum asserts that an interface is a string, and returns a
// pointer to a *RouteWarningCodeEnum with the same value as that string.
func flattenRouteWarningCodeEnum(i interface{}) *RouteWarningCodeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return RouteWarningCodeEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Route) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalRoute(b, c, r)
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

type routeDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         routeApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToRouteDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]routeDiff, error) {
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
	var diffs []routeDiff
	// For each operation name, create a routeDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := routeDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToRouteApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToRouteApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (routeApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractRouteFields(r *Route) error {
	return nil
}
func extractRouteWarningFields(r *Route, o *RouteWarning) error {
	return nil
}

func postReadExtractRouteFields(r *Route) error {
	return nil
}
func postReadExtractRouteWarningFields(r *Route, o *RouteWarning) error {
	return nil
}
