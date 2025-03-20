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

func (r *VpnTunnel) validate() error {

	if err := dcl.Required(r, "sharedSecret"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *VpnTunnel) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *VpnTunnel) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/vpnTunnels/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *VpnTunnel) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/vpnTunnels", nr.basePath(), userBasePath, params), nil

}

func (r *VpnTunnel) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/vpnTunnels", nr.basePath(), userBasePath, params), nil

}

func (r *VpnTunnel) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/vpnTunnels/{{name}}", nr.basePath(), userBasePath, params), nil
}

// vpnTunnelApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type vpnTunnelApiOperation interface {
	do(context.Context, *VpnTunnel, *Client) error
}

func (c *Client) listVpnTunnelRaw(ctx context.Context, r *VpnTunnel, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != VpnTunnelMaxPage {
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

type listVpnTunnelOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listVpnTunnel(ctx context.Context, r *VpnTunnel, pageToken string, pageSize int32) ([]*VpnTunnel, string, error) {
	b, err := c.listVpnTunnelRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listVpnTunnelOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*VpnTunnel
	for _, v := range m.Items {
		res, err := unmarshalMapVpnTunnel(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllVpnTunnel(ctx context.Context, f func(*VpnTunnel) bool, resources []*VpnTunnel) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteVpnTunnel(ctx, res)
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

type deleteVpnTunnelOperation struct{}

func (op *deleteVpnTunnelOperation) do(ctx context.Context, r *VpnTunnel, c *Client) error {
	r, err := c.GetVpnTunnel(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "VpnTunnel not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetVpnTunnel checking for existence. error: %v", err)
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
		_, err := c.GetVpnTunnel(ctx, r)
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
type createVpnTunnelOperation struct {
	response map[string]interface{}
}

func (op *createVpnTunnelOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createVpnTunnelOperation) do(ctx context.Context, r *VpnTunnel, c *Client) error {
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

	if _, err := c.GetVpnTunnel(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getVpnTunnelRaw(ctx context.Context, r *VpnTunnel) ([]byte, error) {
	if dcl.IsZeroValue(r.IkeVersion) {
		r.IkeVersion = dcl.Int64(2)
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

func (c *Client) vpnTunnelDiffsForRawDesired(ctx context.Context, rawDesired *VpnTunnel, opts ...dcl.ApplyOption) (initial, desired *VpnTunnel, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *VpnTunnel
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*VpnTunnel); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected VpnTunnel, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetVpnTunnel(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a VpnTunnel resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve VpnTunnel resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that VpnTunnel resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeVpnTunnelDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for VpnTunnel: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for VpnTunnel: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractVpnTunnelFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeVpnTunnelInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for VpnTunnel: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeVpnTunnelDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for VpnTunnel: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffVpnTunnel(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeVpnTunnelInitialState(rawInitial, rawDesired *VpnTunnel) (*VpnTunnel, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeVpnTunnelDesiredState(rawDesired, rawInitial *VpnTunnel, opts ...dcl.ApplyOption) (*VpnTunnel, error) {

	if dcl.IsZeroValue(rawDesired.IkeVersion) {
		rawDesired.IkeVersion = dcl.Int64(2)
	}

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &VpnTunnel{}
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
	if dcl.StringCanonicalize(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	if dcl.IsZeroValue(rawDesired.TargetVpnGateway) || (dcl.IsEmptyValueIndirect(rawDesired.TargetVpnGateway) && dcl.IsEmptyValueIndirect(rawInitial.TargetVpnGateway)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.TargetVpnGateway = rawInitial.TargetVpnGateway
	} else {
		canonicalDesired.TargetVpnGateway = rawDesired.TargetVpnGateway
	}
	if dcl.IsZeroValue(rawDesired.VpnGateway) || (dcl.IsEmptyValueIndirect(rawDesired.VpnGateway) && dcl.IsEmptyValueIndirect(rawInitial.VpnGateway)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.VpnGateway = rawInitial.VpnGateway
	} else {
		canonicalDesired.VpnGateway = rawDesired.VpnGateway
	}
	if dcl.IsZeroValue(rawDesired.VpnGatewayInterface) || (dcl.IsEmptyValueIndirect(rawDesired.VpnGatewayInterface) && dcl.IsEmptyValueIndirect(rawInitial.VpnGatewayInterface)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.VpnGatewayInterface = rawInitial.VpnGatewayInterface
	} else {
		canonicalDesired.VpnGatewayInterface = rawDesired.VpnGatewayInterface
	}
	if dcl.StringCanonicalize(rawDesired.PeerExternalGateway, rawInitial.PeerExternalGateway) {
		canonicalDesired.PeerExternalGateway = rawInitial.PeerExternalGateway
	} else {
		canonicalDesired.PeerExternalGateway = rawDesired.PeerExternalGateway
	}
	if dcl.IsZeroValue(rawDesired.PeerExternalGatewayInterface) || (dcl.IsEmptyValueIndirect(rawDesired.PeerExternalGatewayInterface) && dcl.IsEmptyValueIndirect(rawInitial.PeerExternalGatewayInterface)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.PeerExternalGatewayInterface = rawInitial.PeerExternalGatewayInterface
	} else {
		canonicalDesired.PeerExternalGatewayInterface = rawDesired.PeerExternalGatewayInterface
	}
	if dcl.StringCanonicalize(rawDesired.PeerGcpGateway, rawInitial.PeerGcpGateway) {
		canonicalDesired.PeerGcpGateway = rawInitial.PeerGcpGateway
	} else {
		canonicalDesired.PeerGcpGateway = rawDesired.PeerGcpGateway
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Router, rawInitial.Router) {
		canonicalDesired.Router = rawInitial.Router
	} else {
		canonicalDesired.Router = rawDesired.Router
	}
	if dcl.StringCanonicalize(rawDesired.PeerIP, rawInitial.PeerIP) {
		canonicalDesired.PeerIP = rawInitial.PeerIP
	} else {
		canonicalDesired.PeerIP = rawDesired.PeerIP
	}
	if dcl.StringCanonicalize(rawDesired.SharedSecret, rawInitial.SharedSecret) {
		canonicalDesired.SharedSecret = rawInitial.SharedSecret
	} else {
		canonicalDesired.SharedSecret = rawDesired.SharedSecret
	}
	if dcl.IsZeroValue(rawDesired.IkeVersion) || (dcl.IsEmptyValueIndirect(rawDesired.IkeVersion) && dcl.IsEmptyValueIndirect(rawInitial.IkeVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.IkeVersion = rawInitial.IkeVersion
	} else {
		canonicalDesired.IkeVersion = rawDesired.IkeVersion
	}
	if dcl.StringArrayCanonicalize(rawDesired.LocalTrafficSelector, rawInitial.LocalTrafficSelector) {
		canonicalDesired.LocalTrafficSelector = rawInitial.LocalTrafficSelector
	} else {
		canonicalDesired.LocalTrafficSelector = rawDesired.LocalTrafficSelector
	}
	if dcl.StringArrayCanonicalize(rawDesired.RemoteTrafficSelector, rawInitial.RemoteTrafficSelector) {
		canonicalDesired.RemoteTrafficSelector = rawInitial.RemoteTrafficSelector
	} else {
		canonicalDesired.RemoteTrafficSelector = rawDesired.RemoteTrafficSelector
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeVpnTunnelNewState(c *Client, rawNew, rawDesired *VpnTunnel) (*VpnTunnel, error) {

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

	if dcl.IsEmptyValueIndirect(rawNew.Location) && dcl.IsEmptyValueIndirect(rawDesired.Location) {
		rawNew.Location = rawDesired.Location
	} else {
		if dcl.StringCanonicalize(rawDesired.Location, rawNew.Location) {
			rawNew.Location = rawDesired.Location
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.TargetVpnGateway) && dcl.IsEmptyValueIndirect(rawDesired.TargetVpnGateway) {
		rawNew.TargetVpnGateway = rawDesired.TargetVpnGateway
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.VpnGateway) && dcl.IsEmptyValueIndirect(rawDesired.VpnGateway) {
		rawNew.VpnGateway = rawDesired.VpnGateway
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.VpnGatewayInterface) && dcl.IsEmptyValueIndirect(rawDesired.VpnGatewayInterface) {
		rawNew.VpnGatewayInterface = rawDesired.VpnGatewayInterface
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PeerExternalGateway) && dcl.IsEmptyValueIndirect(rawDesired.PeerExternalGateway) {
		rawNew.PeerExternalGateway = rawDesired.PeerExternalGateway
	} else {
		if dcl.StringCanonicalize(rawDesired.PeerExternalGateway, rawNew.PeerExternalGateway) {
			rawNew.PeerExternalGateway = rawDesired.PeerExternalGateway
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.PeerExternalGatewayInterface) && dcl.IsEmptyValueIndirect(rawDesired.PeerExternalGatewayInterface) {
		rawNew.PeerExternalGatewayInterface = rawDesired.PeerExternalGatewayInterface
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PeerGcpGateway) && dcl.IsEmptyValueIndirect(rawDesired.PeerGcpGateway) {
		rawNew.PeerGcpGateway = rawDesired.PeerGcpGateway
	} else {
		if dcl.StringCanonicalize(rawDesired.PeerGcpGateway, rawNew.PeerGcpGateway) {
			rawNew.PeerGcpGateway = rawDesired.PeerGcpGateway
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Router) && dcl.IsEmptyValueIndirect(rawDesired.Router) {
		rawNew.Router = rawDesired.Router
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Router, rawNew.Router) {
			rawNew.Router = rawDesired.Router
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.PeerIP) && dcl.IsEmptyValueIndirect(rawDesired.PeerIP) {
		rawNew.PeerIP = rawDesired.PeerIP
	} else {
		if dcl.StringCanonicalize(rawDesired.PeerIP, rawNew.PeerIP) {
			rawNew.PeerIP = rawDesired.PeerIP
		}
	}

	rawNew.SharedSecret = rawDesired.SharedSecret

	if dcl.IsEmptyValueIndirect(rawNew.SharedSecretHash) && dcl.IsEmptyValueIndirect(rawDesired.SharedSecretHash) {
		rawNew.SharedSecretHash = rawDesired.SharedSecretHash
	} else {
		if dcl.StringCanonicalize(rawDesired.SharedSecretHash, rawNew.SharedSecretHash) {
			rawNew.SharedSecretHash = rawDesired.SharedSecretHash
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Status) && dcl.IsEmptyValueIndirect(rawDesired.Status) {
		rawNew.Status = rawDesired.Status
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.IkeVersion) && dcl.IsEmptyValueIndirect(rawDesired.IkeVersion) {
		rawNew.IkeVersion = rawDesired.IkeVersion
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DetailedStatus) && dcl.IsEmptyValueIndirect(rawDesired.DetailedStatus) {
		rawNew.DetailedStatus = rawDesired.DetailedStatus
	} else {
		if dcl.StringCanonicalize(rawDesired.DetailedStatus, rawNew.DetailedStatus) {
			rawNew.DetailedStatus = rawDesired.DetailedStatus
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.LocalTrafficSelector) && dcl.IsEmptyValueIndirect(rawDesired.LocalTrafficSelector) {
		rawNew.LocalTrafficSelector = rawDesired.LocalTrafficSelector
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.LocalTrafficSelector, rawNew.LocalTrafficSelector) {
			rawNew.LocalTrafficSelector = rawDesired.LocalTrafficSelector
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RemoteTrafficSelector) && dcl.IsEmptyValueIndirect(rawDesired.RemoteTrafficSelector) {
		rawNew.RemoteTrafficSelector = rawDesired.RemoteTrafficSelector
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.RemoteTrafficSelector, rawNew.RemoteTrafficSelector) {
			rawNew.RemoteTrafficSelector = rawDesired.RemoteTrafficSelector
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
func diffVpnTunnel(c *Client, desired, actual *VpnTunnel, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetVpnGateway, actual.TargetVpnGateway, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TargetVpnGateway")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VpnGateway, actual.VpnGateway, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VpnGateway")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VpnGatewayInterface, actual.VpnGatewayInterface, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VpnGatewayInterface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeerExternalGateway, actual.PeerExternalGateway, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeerExternalGateway")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeerExternalGatewayInterface, actual.PeerExternalGatewayInterface, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeerExternalGatewayInterface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeerGcpGateway, actual.PeerGcpGateway, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeerGcpGateway")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Router, actual.Router, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Router")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeerIP, actual.PeerIP, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeerIp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SharedSecret, actual.SharedSecret, dcl.DiffInfo{Ignore: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SharedSecret")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SharedSecretHash, actual.SharedSecretHash, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SharedSecretHash")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.IkeVersion, actual.IkeVersion, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IkeVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DetailedStatus, actual.DetailedStatus, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DetailedStatus")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalTrafficSelector, actual.LocalTrafficSelector, dcl.DiffInfo{ServerDefault: true, Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalTrafficSelector")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RemoteTrafficSelector, actual.RemoteTrafficSelector, dcl.DiffInfo{ServerDefault: true, Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RemoteTrafficSelector")); len(ds) != 0 || err != nil {
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
func (r *VpnTunnel) urlNormalized() *VpnTunnel {
	normalized := dcl.Copy(*r).(VpnTunnel)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.TargetVpnGateway = dcl.SelfLinkToName(r.TargetVpnGateway)
	normalized.VpnGateway = dcl.SelfLinkToName(r.VpnGateway)
	normalized.PeerExternalGateway = dcl.SelfLinkToName(r.PeerExternalGateway)
	normalized.PeerGcpGateway = dcl.SelfLinkToName(r.PeerGcpGateway)
	normalized.Router = dcl.SelfLinkToName(r.Router)
	normalized.PeerIP = dcl.SelfLinkToName(r.PeerIP)
	normalized.SharedSecret = dcl.SelfLinkToName(r.SharedSecret)
	normalized.SharedSecretHash = dcl.SelfLinkToName(r.SharedSecretHash)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.DetailedStatus = dcl.SelfLinkToName(r.DetailedStatus)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *VpnTunnel) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the VpnTunnel resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *VpnTunnel) marshal(c *Client) ([]byte, error) {
	m, err := expandVpnTunnel(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling VpnTunnel: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalVpnTunnel decodes JSON responses into the VpnTunnel resource schema.
func unmarshalVpnTunnel(b []byte, c *Client, res *VpnTunnel) (*VpnTunnel, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapVpnTunnel(m, c, res)
}

func unmarshalMapVpnTunnel(m map[string]interface{}, c *Client, res *VpnTunnel) (*VpnTunnel, error) {

	flattened := flattenVpnTunnel(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandVpnTunnel expands VpnTunnel into a JSON request object.
func expandVpnTunnel(c *Client, f *VpnTunnel) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Location; dcl.ValueShouldBeSent(v) {
		m["region"] = v
	}
	if v := f.TargetVpnGateway; dcl.ValueShouldBeSent(v) {
		m["targetVpnGateway"] = v
	}
	if v := f.VpnGateway; dcl.ValueShouldBeSent(v) {
		m["vpnGateway"] = v
	}
	if v := f.VpnGatewayInterface; v != nil {
		m["vpnGatewayInterface"] = v
	}
	if v := f.PeerExternalGateway; dcl.ValueShouldBeSent(v) {
		m["peerExternalGateway"] = v
	}
	if v := f.PeerExternalGatewayInterface; dcl.ValueShouldBeSent(v) {
		m["peerExternalGatewayInterface"] = v
	}
	if v := f.PeerGcpGateway; dcl.ValueShouldBeSent(v) {
		m["peerGcpGateway"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/regions/%s/routers/%s", f.Router, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Router)); err != nil {
		return nil, fmt.Errorf("error expanding Router into router: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["router"] = v
	}
	if v := f.PeerIP; dcl.ValueShouldBeSent(v) {
		m["peerIp"] = v
	}
	if v := f.SharedSecret; dcl.ValueShouldBeSent(v) {
		m["sharedSecret"] = v
	}
	if v := f.IkeVersion; dcl.ValueShouldBeSent(v) {
		m["ikeVersion"] = v
	}
	if v := f.LocalTrafficSelector; v != nil {
		m["localTrafficSelector"] = v
	}
	if v := f.RemoteTrafficSelector; v != nil {
		m["remoteTrafficSelector"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenVpnTunnel flattens VpnTunnel from a JSON request object into the
// VpnTunnel type.
func flattenVpnTunnel(c *Client, i interface{}, res *VpnTunnel) *VpnTunnel {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &VpnTunnel{}
	resultRes.Id = dcl.FlattenInteger(m["id"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Location = dcl.FlattenString(m["region"])
	resultRes.TargetVpnGateway = dcl.FlattenString(m["targetVpnGateway"])
	resultRes.VpnGateway = dcl.FlattenString(m["vpnGateway"])
	resultRes.VpnGatewayInterface = dcl.FlattenInteger(m["vpnGatewayInterface"])
	resultRes.PeerExternalGateway = dcl.FlattenString(m["peerExternalGateway"])
	resultRes.PeerExternalGatewayInterface = dcl.FlattenInteger(m["peerExternalGatewayInterface"])
	resultRes.PeerGcpGateway = dcl.FlattenString(m["peerGcpGateway"])
	resultRes.Router = dcl.FlattenString(m["router"])
	resultRes.PeerIP = dcl.FlattenString(m["peerIp"])
	resultRes.SharedSecret = dcl.FlattenSecretValue(m["sharedSecret"])
	resultRes.SharedSecretHash = dcl.FlattenString(m["sharedSecretHash"])
	resultRes.Status = flattenVpnTunnelStatusEnum(m["status"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.IkeVersion = dcl.FlattenInteger(m["ikeVersion"])
	if _, ok := m["ikeVersion"]; !ok {
		c.Config.Logger.Info("Using default value for ikeVersion")
		resultRes.IkeVersion = dcl.Int64(2)
	}
	resultRes.DetailedStatus = dcl.FlattenString(m["detailedStatus"])
	resultRes.LocalTrafficSelector = dcl.FlattenStringSlice(m["localTrafficSelector"])
	resultRes.RemoteTrafficSelector = dcl.FlattenStringSlice(m["remoteTrafficSelector"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// flattenVpnTunnelStatusEnumMap flattens the contents of VpnTunnelStatusEnum from a JSON
// response object.
func flattenVpnTunnelStatusEnumMap(c *Client, i interface{}, res *VpnTunnel) map[string]VpnTunnelStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]VpnTunnelStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]VpnTunnelStatusEnum{}
	}

	items := make(map[string]VpnTunnelStatusEnum)
	for k, item := range a {
		items[k] = *flattenVpnTunnelStatusEnum(item.(interface{}))
	}

	return items
}

// flattenVpnTunnelStatusEnumSlice flattens the contents of VpnTunnelStatusEnum from a JSON
// response object.
func flattenVpnTunnelStatusEnumSlice(c *Client, i interface{}, res *VpnTunnel) []VpnTunnelStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []VpnTunnelStatusEnum{}
	}

	if len(a) == 0 {
		return []VpnTunnelStatusEnum{}
	}

	items := make([]VpnTunnelStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenVpnTunnelStatusEnum(item.(interface{})))
	}

	return items
}

// flattenVpnTunnelStatusEnum asserts that an interface is a string, and returns a
// pointer to a *VpnTunnelStatusEnum with the same value as that string.
func flattenVpnTunnelStatusEnum(i interface{}) *VpnTunnelStatusEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return VpnTunnelStatusEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *VpnTunnel) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalVpnTunnel(b, c, r)
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

type vpnTunnelDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         vpnTunnelApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToVpnTunnelDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]vpnTunnelDiff, error) {
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
	var diffs []vpnTunnelDiff
	// For each operation name, create a vpnTunnelDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := vpnTunnelDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToVpnTunnelApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToVpnTunnelApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (vpnTunnelApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractVpnTunnelFields(r *VpnTunnel) error {
	return nil
}

func postReadExtractVpnTunnelFields(r *VpnTunnel) error {
	return nil
}
