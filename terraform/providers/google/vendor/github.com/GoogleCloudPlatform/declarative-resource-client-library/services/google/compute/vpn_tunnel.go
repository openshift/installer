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
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type VpnTunnel struct {
	Id                           *int64               `json:"id"`
	Name                         *string              `json:"name"`
	Description                  *string              `json:"description"`
	Location                     *string              `json:"location"`
	TargetVpnGateway             *string              `json:"targetVpnGateway"`
	VpnGateway                   *string              `json:"vpnGateway"`
	VpnGatewayInterface          *int64               `json:"vpnGatewayInterface"`
	PeerExternalGateway          *string              `json:"peerExternalGateway"`
	PeerExternalGatewayInterface *int64               `json:"peerExternalGatewayInterface"`
	PeerGcpGateway               *string              `json:"peerGcpGateway"`
	Router                       *string              `json:"router"`
	PeerIP                       *string              `json:"peerIP"`
	SharedSecret                 *string              `json:"sharedSecret"`
	SharedSecretHash             *string              `json:"sharedSecretHash"`
	Status                       *VpnTunnelStatusEnum `json:"status"`
	SelfLink                     *string              `json:"selfLink"`
	IkeVersion                   *int64               `json:"ikeVersion"`
	DetailedStatus               *string              `json:"detailedStatus"`
	LocalTrafficSelector         []string             `json:"localTrafficSelector"`
	RemoteTrafficSelector        []string             `json:"remoteTrafficSelector"`
	Project                      *string              `json:"project"`
}

func (r *VpnTunnel) String() string {
	return dcl.SprintResource(r)
}

// The enum VpnTunnelStatusEnum.
type VpnTunnelStatusEnum string

// VpnTunnelStatusEnumRef returns a *VpnTunnelStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func VpnTunnelStatusEnumRef(s string) *VpnTunnelStatusEnum {
	v := VpnTunnelStatusEnum(s)
	return &v
}

func (v VpnTunnelStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PROVISIONING", "WAITING_FOR_FULL_CONFIG", "FIRST_HANDSHAKE", "ESTABLISHED", "NO_INCOMING_PACKETS", "AUTHORIZATION_ERROR", "NEGOTIATION_FAILURE", "DEPROVISIONING", "FAILED", "REJECTED", "ALLOCATING_RESOURCES", "STOPPED", "PEER_IDENTITY_MISMATCH", "TS_NARROWING_NOT_ALLOWED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "VpnTunnelStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *VpnTunnel) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "VpnTunnel",
		Version: "compute",
	}
}

func (r *VpnTunnel) ID() (string, error) {
	if err := extractVpnTunnelFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"id":                              dcl.ValueOrEmptyString(nr.Id),
		"name":                            dcl.ValueOrEmptyString(nr.Name),
		"description":                     dcl.ValueOrEmptyString(nr.Description),
		"location":                        dcl.ValueOrEmptyString(nr.Location),
		"target_vpn_gateway":              dcl.ValueOrEmptyString(nr.TargetVpnGateway),
		"vpn_gateway":                     dcl.ValueOrEmptyString(nr.VpnGateway),
		"vpn_gateway_interface":           dcl.ValueOrEmptyString(nr.VpnGatewayInterface),
		"peer_external_gateway":           dcl.ValueOrEmptyString(nr.PeerExternalGateway),
		"peer_external_gateway_interface": dcl.ValueOrEmptyString(nr.PeerExternalGatewayInterface),
		"peer_gcp_gateway":                dcl.ValueOrEmptyString(nr.PeerGcpGateway),
		"router":                          dcl.ValueOrEmptyString(nr.Router),
		"peer_ip":                         dcl.ValueOrEmptyString(nr.PeerIP),
		"shared_secret":                   dcl.ValueOrEmptyString(nr.SharedSecret),
		"shared_secret_hash":              dcl.ValueOrEmptyString(nr.SharedSecretHash),
		"status":                          dcl.ValueOrEmptyString(nr.Status),
		"self_link":                       dcl.ValueOrEmptyString(nr.SelfLink),
		"ike_version":                     dcl.ValueOrEmptyString(nr.IkeVersion),
		"detailed_status":                 dcl.ValueOrEmptyString(nr.DetailedStatus),
		"local_traffic_selector":          dcl.ValueOrEmptyString(nr.LocalTrafficSelector),
		"remote_traffic_selector":         dcl.ValueOrEmptyString(nr.RemoteTrafficSelector),
		"project":                         dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/regions/{{location}}/vpnTunnels/{{name}}", params), nil
}

const VpnTunnelMaxPage = -1

type VpnTunnelList struct {
	Items []*VpnTunnel

	nextToken string

	pageSize int32

	resource *VpnTunnel
}

func (l *VpnTunnelList) HasNext() bool {
	return l.nextToken != ""
}

func (l *VpnTunnelList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listVpnTunnel(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListVpnTunnel(ctx context.Context, project, location string) (*VpnTunnelList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListVpnTunnelWithMaxResults(ctx, project, location, VpnTunnelMaxPage)

}

func (c *Client) ListVpnTunnelWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*VpnTunnelList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &VpnTunnel{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listVpnTunnel(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &VpnTunnelList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetVpnTunnel(ctx context.Context, r *VpnTunnel) (*VpnTunnel, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractVpnTunnelFields(r)

	b, err := c.getVpnTunnelRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalVpnTunnel(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name
	if dcl.IsZeroValue(result.IkeVersion) {
		result.IkeVersion = dcl.Int64(2)
	}

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeVpnTunnelNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractVpnTunnelFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteVpnTunnel(ctx context.Context, r *VpnTunnel) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("VpnTunnel resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting VpnTunnel...")
	deleteOp := deleteVpnTunnelOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllVpnTunnel deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllVpnTunnel(ctx context.Context, project, location string, filter func(*VpnTunnel) bool) error {
	listObj, err := c.ListVpnTunnel(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllVpnTunnel(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllVpnTunnel(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyVpnTunnel(ctx context.Context, rawDesired *VpnTunnel, opts ...dcl.ApplyOption) (*VpnTunnel, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *VpnTunnel
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyVpnTunnelHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyVpnTunnelHelper(c *Client, ctx context.Context, rawDesired *VpnTunnel, opts ...dcl.ApplyOption) (*VpnTunnel, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyVpnTunnel...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractVpnTunnelFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.vpnTunnelDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToVpnTunnelDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []vpnTunnelApiOperation
	if create {
		ops = append(ops, &createVpnTunnelOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyVpnTunnelDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyVpnTunnelDiff(c *Client, ctx context.Context, desired *VpnTunnel, rawDesired *VpnTunnel, ops []vpnTunnelApiOperation, opts ...dcl.ApplyOption) (*VpnTunnel, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetVpnTunnel(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createVpnTunnelOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapVpnTunnel(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeVpnTunnelNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeVpnTunnelNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeVpnTunnelDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractVpnTunnelFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractVpnTunnelFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffVpnTunnel(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}
