// Copyright 2021 Google LLC. All Rights Reserved.
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

type Interconnect struct {
	Description             *string                            `json:"description"`
	SelfLink                *string                            `json:"selfLink"`
	Id                      *int64                             `json:"id"`
	Name                    *string                            `json:"name"`
	Location                *string                            `json:"location"`
	LinkType                *InterconnectLinkTypeEnum          `json:"linkType"`
	RequestedLinkCount      *int64                             `json:"requestedLinkCount"`
	InterconnectType        *InterconnectInterconnectTypeEnum  `json:"interconnectType"`
	AdminEnabled            *bool                              `json:"adminEnabled"`
	NocContactEmail         *string                            `json:"nocContactEmail"`
	CustomerName            *string                            `json:"customerName"`
	OperationalStatus       *InterconnectOperationalStatusEnum `json:"operationalStatus"`
	ProvisionedLinkCount    *int64                             `json:"provisionedLinkCount"`
	InterconnectAttachments []string                           `json:"interconnectAttachments"`
	PeerIPAddress           *string                            `json:"peerIPAddress"`
	GoogleIPAddress         *string                            `json:"googleIPAddress"`
	GoogleReferenceId       *string                            `json:"googleReferenceId"`
	ExpectedOutages         []InterconnectExpectedOutages      `json:"expectedOutages"`
	CircuitInfos            []InterconnectCircuitInfos         `json:"circuitInfos"`
	State                   *InterconnectStateEnum             `json:"state"`
	Project                 *string                            `json:"project"`
}

func (r *Interconnect) String() string {
	return dcl.SprintResource(r)
}

// The enum InterconnectLinkTypeEnum.
type InterconnectLinkTypeEnum string

// InterconnectLinkTypeEnumRef returns a *InterconnectLinkTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectLinkTypeEnumRef(s string) *InterconnectLinkTypeEnum {
	v := InterconnectLinkTypeEnum(s)
	return &v
}

func (v InterconnectLinkTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LINK_TYPE_ETHERNET_10G_LR", "LINK_TYPE_ETHERNET_100G_LR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectLinkTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectInterconnectTypeEnum.
type InterconnectInterconnectTypeEnum string

// InterconnectInterconnectTypeEnumRef returns a *InterconnectInterconnectTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectInterconnectTypeEnumRef(s string) *InterconnectInterconnectTypeEnum {
	v := InterconnectInterconnectTypeEnum(s)
	return &v
}

func (v InterconnectInterconnectTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"IT_PRIVATE", "PARTNER", "DEDICATED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectInterconnectTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectOperationalStatusEnum.
type InterconnectOperationalStatusEnum string

// InterconnectOperationalStatusEnumRef returns a *InterconnectOperationalStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectOperationalStatusEnumRef(s string) *InterconnectOperationalStatusEnum {
	v := InterconnectOperationalStatusEnum(s)
	return &v
}

func (v InterconnectOperationalStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"OS_ACTIVE", "OS_UNPROVISIONED", "OS_UNDER_MAINTENANCE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectOperationalStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectExpectedOutagesSourceEnum.
type InterconnectExpectedOutagesSourceEnum string

// InterconnectExpectedOutagesSourceEnumRef returns a *InterconnectExpectedOutagesSourceEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectExpectedOutagesSourceEnumRef(s string) *InterconnectExpectedOutagesSourceEnum {
	v := InterconnectExpectedOutagesSourceEnum(s)
	return &v
}

func (v InterconnectExpectedOutagesSourceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"GOOGLE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectExpectedOutagesSourceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectExpectedOutagesStateEnum.
type InterconnectExpectedOutagesStateEnum string

// InterconnectExpectedOutagesStateEnumRef returns a *InterconnectExpectedOutagesStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectExpectedOutagesStateEnumRef(s string) *InterconnectExpectedOutagesStateEnum {
	v := InterconnectExpectedOutagesStateEnum(s)
	return &v
}

func (v InterconnectExpectedOutagesStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ACTIVE", "CANCELLED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectExpectedOutagesStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectExpectedOutagesIssueTypeEnum.
type InterconnectExpectedOutagesIssueTypeEnum string

// InterconnectExpectedOutagesIssueTypeEnumRef returns a *InterconnectExpectedOutagesIssueTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectExpectedOutagesIssueTypeEnumRef(s string) *InterconnectExpectedOutagesIssueTypeEnum {
	v := InterconnectExpectedOutagesIssueTypeEnum(s)
	return &v
}

func (v InterconnectExpectedOutagesIssueTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"OUTAGE", "PARTIAL_OUTAGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectExpectedOutagesIssueTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectStateEnum.
type InterconnectStateEnum string

// InterconnectStateEnumRef returns a *InterconnectStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectStateEnumRef(s string) *InterconnectStateEnum {
	v := InterconnectStateEnum(s)
	return &v
}

func (v InterconnectStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"DEPRECATED", "OBSOLETE", "DELETED", "ACTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type InterconnectExpectedOutages struct {
	empty            bool                                      `json:"-"`
	Name             *string                                   `json:"name"`
	Description      *string                                   `json:"description"`
	Source           *InterconnectExpectedOutagesSourceEnum    `json:"source"`
	State            *InterconnectExpectedOutagesStateEnum     `json:"state"`
	IssueType        *InterconnectExpectedOutagesIssueTypeEnum `json:"issueType"`
	AffectedCircuits []string                                  `json:"affectedCircuits"`
	StartTime        *int64                                    `json:"startTime"`
	EndTime          *int64                                    `json:"endTime"`
}

type jsonInterconnectExpectedOutages InterconnectExpectedOutages

func (r *InterconnectExpectedOutages) UnmarshalJSON(data []byte) error {
	var res jsonInterconnectExpectedOutages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInterconnectExpectedOutages
	} else {

		r.Name = res.Name

		r.Description = res.Description

		r.Source = res.Source

		r.State = res.State

		r.IssueType = res.IssueType

		r.AffectedCircuits = res.AffectedCircuits

		r.StartTime = res.StartTime

		r.EndTime = res.EndTime

	}
	return nil
}

// This object is used to assert a desired state where this InterconnectExpectedOutages is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyInterconnectExpectedOutages *InterconnectExpectedOutages = &InterconnectExpectedOutages{empty: true}

func (r *InterconnectExpectedOutages) Empty() bool {
	return r.empty
}

func (r *InterconnectExpectedOutages) String() string {
	return dcl.SprintResource(r)
}

func (r *InterconnectExpectedOutages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InterconnectCircuitInfos struct {
	empty            bool    `json:"-"`
	GoogleCircuitId  *string `json:"googleCircuitId"`
	GoogleDemarcId   *string `json:"googleDemarcId"`
	CustomerDemarcId *string `json:"customerDemarcId"`
}

type jsonInterconnectCircuitInfos InterconnectCircuitInfos

func (r *InterconnectCircuitInfos) UnmarshalJSON(data []byte) error {
	var res jsonInterconnectCircuitInfos
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInterconnectCircuitInfos
	} else {

		r.GoogleCircuitId = res.GoogleCircuitId

		r.GoogleDemarcId = res.GoogleDemarcId

		r.CustomerDemarcId = res.CustomerDemarcId

	}
	return nil
}

// This object is used to assert a desired state where this InterconnectCircuitInfos is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyInterconnectCircuitInfos *InterconnectCircuitInfos = &InterconnectCircuitInfos{empty: true}

func (r *InterconnectCircuitInfos) Empty() bool {
	return r.empty
}

func (r *InterconnectCircuitInfos) String() string {
	return dcl.SprintResource(r)
}

func (r *InterconnectCircuitInfos) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Interconnect) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "Interconnect",
		Version: "compute",
	}
}

func (r *Interconnect) ID() (string, error) {
	if err := extractInterconnectFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"description":             dcl.ValueOrEmptyString(nr.Description),
		"selfLink":                dcl.ValueOrEmptyString(nr.SelfLink),
		"id":                      dcl.ValueOrEmptyString(nr.Id),
		"name":                    dcl.ValueOrEmptyString(nr.Name),
		"location":                dcl.ValueOrEmptyString(nr.Location),
		"linkType":                dcl.ValueOrEmptyString(nr.LinkType),
		"requestedLinkCount":      dcl.ValueOrEmptyString(nr.RequestedLinkCount),
		"interconnectType":        dcl.ValueOrEmptyString(nr.InterconnectType),
		"adminEnabled":            dcl.ValueOrEmptyString(nr.AdminEnabled),
		"nocContactEmail":         dcl.ValueOrEmptyString(nr.NocContactEmail),
		"customerName":            dcl.ValueOrEmptyString(nr.CustomerName),
		"operationalStatus":       dcl.ValueOrEmptyString(nr.OperationalStatus),
		"provisionedLinkCount":    dcl.ValueOrEmptyString(nr.ProvisionedLinkCount),
		"interconnectAttachments": dcl.ValueOrEmptyString(nr.InterconnectAttachments),
		"peerIPAddress":           dcl.ValueOrEmptyString(nr.PeerIPAddress),
		"googleIPAddress":         dcl.ValueOrEmptyString(nr.GoogleIPAddress),
		"googleReferenceId":       dcl.ValueOrEmptyString(nr.GoogleReferenceId),
		"expectedOutages":         dcl.ValueOrEmptyString(nr.ExpectedOutages),
		"circuitInfos":            dcl.ValueOrEmptyString(nr.CircuitInfos),
		"state":                   dcl.ValueOrEmptyString(nr.State),
		"project":                 dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/global/interconnects/{{name}}", params), nil
}

const InterconnectMaxPage = -1

type InterconnectList struct {
	Items []*Interconnect

	nextToken string

	pageSize int32

	resource *Interconnect
}

func (l *InterconnectList) HasNext() bool {
	return l.nextToken != ""
}

func (l *InterconnectList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listInterconnect(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListInterconnect(ctx context.Context, project string) (*InterconnectList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListInterconnectWithMaxResults(ctx, project, InterconnectMaxPage)

}

func (c *Client) ListInterconnectWithMaxResults(ctx context.Context, project string, pageSize int32) (*InterconnectList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Interconnect{
		Project: &project,
	}
	items, token, err := c.listInterconnect(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &InterconnectList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetInterconnect(ctx context.Context, r *Interconnect) (*Interconnect, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractInterconnectFields(r)

	b, err := c.getInterconnectRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalInterconnect(b, c)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeInterconnectNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractInterconnectFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteInterconnect(ctx context.Context, r *Interconnect) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Interconnect resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Interconnect...")
	deleteOp := deleteInterconnectOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllInterconnect deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllInterconnect(ctx context.Context, project string, filter func(*Interconnect) bool) error {
	listObj, err := c.ListInterconnect(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllInterconnect(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllInterconnect(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyInterconnect(ctx context.Context, rawDesired *Interconnect, opts ...dcl.ApplyOption) (*Interconnect, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Interconnect
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyInterconnectHelper(c, ctx, rawDesired, opts...)
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

func applyInterconnectHelper(c *Client, ctx context.Context, rawDesired *Interconnect, opts ...dcl.ApplyOption) (*Interconnect, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyInterconnect...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractInterconnectFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.interconnectDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToInterconnectDiffs(c.Config, fieldDiffs, opts)
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
	var ops []interconnectApiOperation
	if create {
		ops = append(ops, &createInterconnectOperation{})
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
	return applyInterconnectDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyInterconnectDiff(c *Client, ctx context.Context, desired *Interconnect, rawDesired *Interconnect, ops []interconnectApiOperation, opts ...dcl.ApplyOption) (*Interconnect, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetInterconnect(ctx, desired.urlNormalized())
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createInterconnectOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapInterconnect(r, c)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeInterconnectNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeInterconnectNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeInterconnectDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractInterconnectFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractInterconnectFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffInterconnect(c, newDesired, newState)
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
