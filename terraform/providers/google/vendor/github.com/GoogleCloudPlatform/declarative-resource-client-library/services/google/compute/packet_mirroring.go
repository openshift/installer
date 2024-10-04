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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type PacketMirroring struct {
	Id                *int64                            `json:"id"`
	SelfLink          *string                           `json:"selfLink"`
	Name              *string                           `json:"name"`
	Description       *string                           `json:"description"`
	Region            *string                           `json:"region"`
	Network           *PacketMirroringNetwork           `json:"network"`
	Priority          *int64                            `json:"priority"`
	CollectorIlb      *PacketMirroringCollectorIlb      `json:"collectorIlb"`
	MirroredResources *PacketMirroringMirroredResources `json:"mirroredResources"`
	Filter            *PacketMirroringFilter            `json:"filter"`
	Enable            *PacketMirroringEnableEnum        `json:"enable"`
	Project           *string                           `json:"project"`
	Location          *string                           `json:"location"`
}

func (r *PacketMirroring) String() string {
	return dcl.SprintResource(r)
}

// The enum PacketMirroringFilterDirectionEnum.
type PacketMirroringFilterDirectionEnum string

// PacketMirroringFilterDirectionEnumRef returns a *PacketMirroringFilterDirectionEnum with the value of string s
// If the empty string is provided, nil is returned.
func PacketMirroringFilterDirectionEnumRef(s string) *PacketMirroringFilterDirectionEnum {
	v := PacketMirroringFilterDirectionEnum(s)
	return &v
}

func (v PacketMirroringFilterDirectionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"INGRESS", "EGRESS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PacketMirroringFilterDirectionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum PacketMirroringEnableEnum.
type PacketMirroringEnableEnum string

// PacketMirroringEnableEnumRef returns a *PacketMirroringEnableEnum with the value of string s
// If the empty string is provided, nil is returned.
func PacketMirroringEnableEnumRef(s string) *PacketMirroringEnableEnum {
	v := PacketMirroringEnableEnum(s)
	return &v
}

func (v PacketMirroringEnableEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TRUE", "FALSE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "PacketMirroringEnableEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type PacketMirroringNetwork struct {
	empty        bool    `json:"-"`
	Url          *string `json:"url"`
	CanonicalUrl *string `json:"canonicalUrl"`
}

type jsonPacketMirroringNetwork PacketMirroringNetwork

func (r *PacketMirroringNetwork) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringNetwork
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringNetwork
	} else {

		r.Url = res.Url

		r.CanonicalUrl = res.CanonicalUrl

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringNetwork is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringNetwork *PacketMirroringNetwork = &PacketMirroringNetwork{empty: true}

func (r *PacketMirroringNetwork) Empty() bool {
	return r.empty
}

func (r *PacketMirroringNetwork) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringNetwork) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PacketMirroringCollectorIlb struct {
	empty        bool    `json:"-"`
	Url          *string `json:"url"`
	CanonicalUrl *string `json:"canonicalUrl"`
}

type jsonPacketMirroringCollectorIlb PacketMirroringCollectorIlb

func (r *PacketMirroringCollectorIlb) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringCollectorIlb
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringCollectorIlb
	} else {

		r.Url = res.Url

		r.CanonicalUrl = res.CanonicalUrl

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringCollectorIlb is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringCollectorIlb *PacketMirroringCollectorIlb = &PacketMirroringCollectorIlb{empty: true}

func (r *PacketMirroringCollectorIlb) Empty() bool {
	return r.empty
}

func (r *PacketMirroringCollectorIlb) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringCollectorIlb) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PacketMirroringMirroredResources struct {
	empty       bool                                          `json:"-"`
	Subnetworks []PacketMirroringMirroredResourcesSubnetworks `json:"subnetworks"`
	Instances   []PacketMirroringMirroredResourcesInstances   `json:"instances"`
	Tags        []string                                      `json:"tags"`
}

type jsonPacketMirroringMirroredResources PacketMirroringMirroredResources

func (r *PacketMirroringMirroredResources) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringMirroredResources
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringMirroredResources
	} else {

		r.Subnetworks = res.Subnetworks

		r.Instances = res.Instances

		r.Tags = res.Tags

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringMirroredResources is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringMirroredResources *PacketMirroringMirroredResources = &PacketMirroringMirroredResources{empty: true}

func (r *PacketMirroringMirroredResources) Empty() bool {
	return r.empty
}

func (r *PacketMirroringMirroredResources) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringMirroredResources) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PacketMirroringMirroredResourcesSubnetworks struct {
	empty        bool    `json:"-"`
	Url          *string `json:"url"`
	CanonicalUrl *string `json:"canonicalUrl"`
}

type jsonPacketMirroringMirroredResourcesSubnetworks PacketMirroringMirroredResourcesSubnetworks

func (r *PacketMirroringMirroredResourcesSubnetworks) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringMirroredResourcesSubnetworks
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringMirroredResourcesSubnetworks
	} else {

		r.Url = res.Url

		r.CanonicalUrl = res.CanonicalUrl

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringMirroredResourcesSubnetworks is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringMirroredResourcesSubnetworks *PacketMirroringMirroredResourcesSubnetworks = &PacketMirroringMirroredResourcesSubnetworks{empty: true}

func (r *PacketMirroringMirroredResourcesSubnetworks) Empty() bool {
	return r.empty
}

func (r *PacketMirroringMirroredResourcesSubnetworks) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringMirroredResourcesSubnetworks) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PacketMirroringMirroredResourcesInstances struct {
	empty        bool    `json:"-"`
	Url          *string `json:"url"`
	CanonicalUrl *string `json:"canonicalUrl"`
}

type jsonPacketMirroringMirroredResourcesInstances PacketMirroringMirroredResourcesInstances

func (r *PacketMirroringMirroredResourcesInstances) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringMirroredResourcesInstances
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringMirroredResourcesInstances
	} else {

		r.Url = res.Url

		r.CanonicalUrl = res.CanonicalUrl

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringMirroredResourcesInstances is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringMirroredResourcesInstances *PacketMirroringMirroredResourcesInstances = &PacketMirroringMirroredResourcesInstances{empty: true}

func (r *PacketMirroringMirroredResourcesInstances) Empty() bool {
	return r.empty
}

func (r *PacketMirroringMirroredResourcesInstances) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringMirroredResourcesInstances) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type PacketMirroringFilter struct {
	empty       bool                                `json:"-"`
	CidrRanges  []string                            `json:"cidrRanges"`
	IPProtocols []string                            `json:"ipProtocols"`
	Direction   *PacketMirroringFilterDirectionEnum `json:"direction"`
}

type jsonPacketMirroringFilter PacketMirroringFilter

func (r *PacketMirroringFilter) UnmarshalJSON(data []byte) error {
	var res jsonPacketMirroringFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyPacketMirroringFilter
	} else {

		r.CidrRanges = res.CidrRanges

		r.IPProtocols = res.IPProtocols

		r.Direction = res.Direction

	}
	return nil
}

// This object is used to assert a desired state where this PacketMirroringFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyPacketMirroringFilter *PacketMirroringFilter = &PacketMirroringFilter{empty: true}

func (r *PacketMirroringFilter) Empty() bool {
	return r.empty
}

func (r *PacketMirroringFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *PacketMirroringFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *PacketMirroring) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "PacketMirroring",
		Version: "compute",
	}
}

func (r *PacketMirroring) ID() (string, error) {
	if err := extractPacketMirroringFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"id":                 dcl.ValueOrEmptyString(nr.Id),
		"self_link":          dcl.ValueOrEmptyString(nr.SelfLink),
		"name":               dcl.ValueOrEmptyString(nr.Name),
		"description":        dcl.ValueOrEmptyString(nr.Description),
		"region":             dcl.ValueOrEmptyString(nr.Region),
		"network":            dcl.ValueOrEmptyString(nr.Network),
		"priority":           dcl.ValueOrEmptyString(nr.Priority),
		"collector_ilb":      dcl.ValueOrEmptyString(nr.CollectorIlb),
		"mirrored_resources": dcl.ValueOrEmptyString(nr.MirroredResources),
		"filter":             dcl.ValueOrEmptyString(nr.Filter),
		"enable":             dcl.ValueOrEmptyString(nr.Enable),
		"project":            dcl.ValueOrEmptyString(nr.Project),
		"location":           dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/regions/{{location}}/packetMirrorings/{{name}}", params), nil
}

const PacketMirroringMaxPage = -1

type PacketMirroringList struct {
	Items []*PacketMirroring

	nextToken string

	pageSize int32

	resource *PacketMirroring
}

func (l *PacketMirroringList) HasNext() bool {
	return l.nextToken != ""
}

func (l *PacketMirroringList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listPacketMirroring(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListPacketMirroring(ctx context.Context, project, location string) (*PacketMirroringList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListPacketMirroringWithMaxResults(ctx, project, location, PacketMirroringMaxPage)

}

func (c *Client) ListPacketMirroringWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*PacketMirroringList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &PacketMirroring{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listPacketMirroring(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &PacketMirroringList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetPacketMirroring(ctx context.Context, r *PacketMirroring) (*PacketMirroring, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractPacketMirroringFields(r)

	b, err := c.getPacketMirroringRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalPacketMirroring(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizePacketMirroringNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractPacketMirroringFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeletePacketMirroring(ctx context.Context, r *PacketMirroring) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("PacketMirroring resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting PacketMirroring...")
	deleteOp := deletePacketMirroringOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllPacketMirroring deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllPacketMirroring(ctx context.Context, project, location string, filter func(*PacketMirroring) bool) error {
	listObj, err := c.ListPacketMirroring(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllPacketMirroring(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllPacketMirroring(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyPacketMirroring(ctx context.Context, rawDesired *PacketMirroring, opts ...dcl.ApplyOption) (*PacketMirroring, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *PacketMirroring
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyPacketMirroringHelper(c, ctx, rawDesired, opts...)
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

func applyPacketMirroringHelper(c *Client, ctx context.Context, rawDesired *PacketMirroring, opts ...dcl.ApplyOption) (*PacketMirroring, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyPacketMirroring...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractPacketMirroringFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.packetMirroringDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToPacketMirroringDiffs(c.Config, fieldDiffs, opts)
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
	var ops []packetMirroringApiOperation
	if create {
		ops = append(ops, &createPacketMirroringOperation{})
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
	return applyPacketMirroringDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyPacketMirroringDiff(c *Client, ctx context.Context, desired *PacketMirroring, rawDesired *PacketMirroring, ops []packetMirroringApiOperation, opts ...dcl.ApplyOption) (*PacketMirroring, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetPacketMirroring(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createPacketMirroringOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapPacketMirroring(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizePacketMirroringNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizePacketMirroringNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizePacketMirroringDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractPacketMirroringFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractPacketMirroringFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffPacketMirroring(c, newDesired, newState)
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
