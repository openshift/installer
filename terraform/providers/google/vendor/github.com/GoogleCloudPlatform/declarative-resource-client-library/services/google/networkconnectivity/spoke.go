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
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Spoke struct {
	Name                           *string                              `json:"name"`
	CreateTime                     *string                              `json:"createTime"`
	UpdateTime                     *string                              `json:"updateTime"`
	Labels                         map[string]string                    `json:"labels"`
	Description                    *string                              `json:"description"`
	Hub                            *string                              `json:"hub"`
	LinkedVpnTunnels               *SpokeLinkedVpnTunnels               `json:"linkedVpnTunnels"`
	LinkedInterconnectAttachments  *SpokeLinkedInterconnectAttachments  `json:"linkedInterconnectAttachments"`
	LinkedRouterApplianceInstances *SpokeLinkedRouterApplianceInstances `json:"linkedRouterApplianceInstances"`
	LinkedVPCNetwork               *SpokeLinkedVPCNetwork               `json:"linkedVPCNetwork"`
	UniqueId                       *string                              `json:"uniqueId"`
	State                          *SpokeStateEnum                      `json:"state"`
	Project                        *string                              `json:"project"`
	Location                       *string                              `json:"location"`
}

func (r *Spoke) String() string {
	return dcl.SprintResource(r)
}

// The enum SpokeStateEnum.
type SpokeStateEnum string

// SpokeStateEnumRef returns a *SpokeStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func SpokeStateEnumRef(s string) *SpokeStateEnum {
	v := SpokeStateEnum(s)
	return &v
}

func (v SpokeStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "CREATING", "ACTIVE", "DELETING"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "SpokeStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type SpokeLinkedVpnTunnels struct {
	empty                  bool     `json:"-"`
	Uris                   []string `json:"uris"`
	SiteToSiteDataTransfer *bool    `json:"siteToSiteDataTransfer"`
}

type jsonSpokeLinkedVpnTunnels SpokeLinkedVpnTunnels

func (r *SpokeLinkedVpnTunnels) UnmarshalJSON(data []byte) error {
	var res jsonSpokeLinkedVpnTunnels
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptySpokeLinkedVpnTunnels
	} else {

		r.Uris = res.Uris

		r.SiteToSiteDataTransfer = res.SiteToSiteDataTransfer

	}
	return nil
}

// This object is used to assert a desired state where this SpokeLinkedVpnTunnels is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptySpokeLinkedVpnTunnels *SpokeLinkedVpnTunnels = &SpokeLinkedVpnTunnels{empty: true}

func (r *SpokeLinkedVpnTunnels) Empty() bool {
	return r.empty
}

func (r *SpokeLinkedVpnTunnels) String() string {
	return dcl.SprintResource(r)
}

func (r *SpokeLinkedVpnTunnels) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type SpokeLinkedInterconnectAttachments struct {
	empty                  bool     `json:"-"`
	Uris                   []string `json:"uris"`
	SiteToSiteDataTransfer *bool    `json:"siteToSiteDataTransfer"`
}

type jsonSpokeLinkedInterconnectAttachments SpokeLinkedInterconnectAttachments

func (r *SpokeLinkedInterconnectAttachments) UnmarshalJSON(data []byte) error {
	var res jsonSpokeLinkedInterconnectAttachments
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptySpokeLinkedInterconnectAttachments
	} else {

		r.Uris = res.Uris

		r.SiteToSiteDataTransfer = res.SiteToSiteDataTransfer

	}
	return nil
}

// This object is used to assert a desired state where this SpokeLinkedInterconnectAttachments is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptySpokeLinkedInterconnectAttachments *SpokeLinkedInterconnectAttachments = &SpokeLinkedInterconnectAttachments{empty: true}

func (r *SpokeLinkedInterconnectAttachments) Empty() bool {
	return r.empty
}

func (r *SpokeLinkedInterconnectAttachments) String() string {
	return dcl.SprintResource(r)
}

func (r *SpokeLinkedInterconnectAttachments) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type SpokeLinkedRouterApplianceInstances struct {
	empty                  bool                                           `json:"-"`
	Instances              []SpokeLinkedRouterApplianceInstancesInstances `json:"instances"`
	SiteToSiteDataTransfer *bool                                          `json:"siteToSiteDataTransfer"`
}

type jsonSpokeLinkedRouterApplianceInstances SpokeLinkedRouterApplianceInstances

func (r *SpokeLinkedRouterApplianceInstances) UnmarshalJSON(data []byte) error {
	var res jsonSpokeLinkedRouterApplianceInstances
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptySpokeLinkedRouterApplianceInstances
	} else {

		r.Instances = res.Instances

		r.SiteToSiteDataTransfer = res.SiteToSiteDataTransfer

	}
	return nil
}

// This object is used to assert a desired state where this SpokeLinkedRouterApplianceInstances is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptySpokeLinkedRouterApplianceInstances *SpokeLinkedRouterApplianceInstances = &SpokeLinkedRouterApplianceInstances{empty: true}

func (r *SpokeLinkedRouterApplianceInstances) Empty() bool {
	return r.empty
}

func (r *SpokeLinkedRouterApplianceInstances) String() string {
	return dcl.SprintResource(r)
}

func (r *SpokeLinkedRouterApplianceInstances) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type SpokeLinkedRouterApplianceInstancesInstances struct {
	empty          bool    `json:"-"`
	VirtualMachine *string `json:"virtualMachine"`
	IPAddress      *string `json:"ipAddress"`
}

type jsonSpokeLinkedRouterApplianceInstancesInstances SpokeLinkedRouterApplianceInstancesInstances

func (r *SpokeLinkedRouterApplianceInstancesInstances) UnmarshalJSON(data []byte) error {
	var res jsonSpokeLinkedRouterApplianceInstancesInstances
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptySpokeLinkedRouterApplianceInstancesInstances
	} else {

		r.VirtualMachine = res.VirtualMachine

		r.IPAddress = res.IPAddress

	}
	return nil
}

// This object is used to assert a desired state where this SpokeLinkedRouterApplianceInstancesInstances is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptySpokeLinkedRouterApplianceInstancesInstances *SpokeLinkedRouterApplianceInstancesInstances = &SpokeLinkedRouterApplianceInstancesInstances{empty: true}

func (r *SpokeLinkedRouterApplianceInstancesInstances) Empty() bool {
	return r.empty
}

func (r *SpokeLinkedRouterApplianceInstancesInstances) String() string {
	return dcl.SprintResource(r)
}

func (r *SpokeLinkedRouterApplianceInstancesInstances) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type SpokeLinkedVPCNetwork struct {
	empty               bool     `json:"-"`
	Uri                 *string  `json:"uri"`
	ExcludeExportRanges []string `json:"excludeExportRanges"`
}

type jsonSpokeLinkedVPCNetwork SpokeLinkedVPCNetwork

func (r *SpokeLinkedVPCNetwork) UnmarshalJSON(data []byte) error {
	var res jsonSpokeLinkedVPCNetwork
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptySpokeLinkedVPCNetwork
	} else {

		r.Uri = res.Uri

		r.ExcludeExportRanges = res.ExcludeExportRanges

	}
	return nil
}

// This object is used to assert a desired state where this SpokeLinkedVPCNetwork is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptySpokeLinkedVPCNetwork *SpokeLinkedVPCNetwork = &SpokeLinkedVPCNetwork{empty: true}

func (r *SpokeLinkedVPCNetwork) Empty() bool {
	return r.empty
}

func (r *SpokeLinkedVPCNetwork) String() string {
	return dcl.SprintResource(r)
}

func (r *SpokeLinkedVPCNetwork) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Spoke) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "network_connectivity",
		Type:    "Spoke",
		Version: "networkconnectivity",
	}
}

func (r *Spoke) ID() (string, error) {
	if err := extractSpokeFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                              dcl.ValueOrEmptyString(nr.Name),
		"create_time":                       dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":                       dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":                            dcl.ValueOrEmptyString(nr.Labels),
		"description":                       dcl.ValueOrEmptyString(nr.Description),
		"hub":                               dcl.ValueOrEmptyString(nr.Hub),
		"linked_vpn_tunnels":                dcl.ValueOrEmptyString(nr.LinkedVpnTunnels),
		"linked_interconnect_attachments":   dcl.ValueOrEmptyString(nr.LinkedInterconnectAttachments),
		"linked_router_appliance_instances": dcl.ValueOrEmptyString(nr.LinkedRouterApplianceInstances),
		"linked_vpc_network":                dcl.ValueOrEmptyString(nr.LinkedVPCNetwork),
		"unique_id":                         dcl.ValueOrEmptyString(nr.UniqueId),
		"state":                             dcl.ValueOrEmptyString(nr.State),
		"project":                           dcl.ValueOrEmptyString(nr.Project),
		"location":                          dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/spokes/{{name}}", params), nil
}

const SpokeMaxPage = -1

type SpokeList struct {
	Items []*Spoke

	nextToken string

	pageSize int32

	resource *Spoke
}

func (l *SpokeList) HasNext() bool {
	return l.nextToken != ""
}

func (l *SpokeList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listSpoke(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListSpoke(ctx context.Context, project, location string) (*SpokeList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListSpokeWithMaxResults(ctx, project, location, SpokeMaxPage)

}

func (c *Client) ListSpokeWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*SpokeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Spoke{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listSpoke(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &SpokeList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetSpoke(ctx context.Context, r *Spoke) (*Spoke, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractSpokeFields(r)

	b, err := c.getSpokeRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalSpoke(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeSpokeNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractSpokeFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteSpoke(ctx context.Context, r *Spoke) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Spoke resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Spoke...")
	deleteOp := deleteSpokeOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllSpoke deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllSpoke(ctx context.Context, project, location string, filter func(*Spoke) bool) error {
	listObj, err := c.ListSpoke(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllSpoke(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllSpoke(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplySpoke(ctx context.Context, rawDesired *Spoke, opts ...dcl.ApplyOption) (*Spoke, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Spoke
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applySpokeHelper(c, ctx, rawDesired, opts...)
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

func applySpokeHelper(c *Client, ctx context.Context, rawDesired *Spoke, opts ...dcl.ApplyOption) (*Spoke, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplySpoke...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractSpokeFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.spokeDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToSpokeDiffs(c.Config, fieldDiffs, opts)
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
	var ops []spokeApiOperation
	if create {
		ops = append(ops, &createSpokeOperation{})
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
	return applySpokeDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applySpokeDiff(c *Client, ctx context.Context, desired *Spoke, rawDesired *Spoke, ops []spokeApiOperation, opts ...dcl.ApplyOption) (*Spoke, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetSpoke(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createSpokeOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapSpoke(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeSpokeNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeSpokeNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeSpokeDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractSpokeFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractSpokeFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffSpoke(c, newDesired, newState)
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
