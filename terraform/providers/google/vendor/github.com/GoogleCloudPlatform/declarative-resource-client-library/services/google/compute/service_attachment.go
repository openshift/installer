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

type ServiceAttachment struct {
	Id                     *int64                                     `json:"id"`
	Name                   *string                                    `json:"name"`
	Description            *string                                    `json:"description"`
	SelfLink               *string                                    `json:"selfLink"`
	Region                 *string                                    `json:"region"`
	TargetService          *string                                    `json:"targetService"`
	ConnectionPreference   *ServiceAttachmentConnectionPreferenceEnum `json:"connectionPreference"`
	ConnectedEndpoints     []ServiceAttachmentConnectedEndpoints      `json:"connectedEndpoints"`
	NatSubnets             []string                                   `json:"natSubnets"`
	EnableProxyProtocol    *bool                                      `json:"enableProxyProtocol"`
	ConsumerRejectLists    []string                                   `json:"consumerRejectLists"`
	ConsumerAcceptLists    []ServiceAttachmentConsumerAcceptLists     `json:"consumerAcceptLists"`
	PscServiceAttachmentId *ServiceAttachmentPscServiceAttachmentId   `json:"pscServiceAttachmentId"`
	Fingerprint            *string                                    `json:"fingerprint"`
	Project                *string                                    `json:"project"`
	Location               *string                                    `json:"location"`
}

func (r *ServiceAttachment) String() string {
	return dcl.SprintResource(r)
}

// The enum ServiceAttachmentConnectionPreferenceEnum.
type ServiceAttachmentConnectionPreferenceEnum string

// ServiceAttachmentConnectionPreferenceEnumRef returns a *ServiceAttachmentConnectionPreferenceEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceAttachmentConnectionPreferenceEnumRef(s string) *ServiceAttachmentConnectionPreferenceEnum {
	v := ServiceAttachmentConnectionPreferenceEnum(s)
	return &v
}

func (v ServiceAttachmentConnectionPreferenceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"CONNECTION_PREFERENCE_UNSPECIFIED", "ACCEPT_AUTOMATIC", "ACCEPT_MANUAL"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceAttachmentConnectionPreferenceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ServiceAttachmentConnectedEndpointsStatusEnum.
type ServiceAttachmentConnectedEndpointsStatusEnum string

// ServiceAttachmentConnectedEndpointsStatusEnumRef returns a *ServiceAttachmentConnectedEndpointsStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func ServiceAttachmentConnectedEndpointsStatusEnumRef(s string) *ServiceAttachmentConnectedEndpointsStatusEnum {
	v := ServiceAttachmentConnectedEndpointsStatusEnum(s)
	return &v
}

func (v ServiceAttachmentConnectedEndpointsStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PENDING", "RUNNING", "DONE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ServiceAttachmentConnectedEndpointsStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ServiceAttachmentConnectedEndpoints struct {
	empty           bool                                           `json:"-"`
	Status          *ServiceAttachmentConnectedEndpointsStatusEnum `json:"status"`
	PscConnectionId *int64                                         `json:"pscConnectionId"`
	Endpoint        *string                                        `json:"endpoint"`
}

type jsonServiceAttachmentConnectedEndpoints ServiceAttachmentConnectedEndpoints

func (r *ServiceAttachmentConnectedEndpoints) UnmarshalJSON(data []byte) error {
	var res jsonServiceAttachmentConnectedEndpoints
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceAttachmentConnectedEndpoints
	} else {

		r.Status = res.Status

		r.PscConnectionId = res.PscConnectionId

		r.Endpoint = res.Endpoint

	}
	return nil
}

// This object is used to assert a desired state where this ServiceAttachmentConnectedEndpoints is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceAttachmentConnectedEndpoints *ServiceAttachmentConnectedEndpoints = &ServiceAttachmentConnectedEndpoints{empty: true}

func (r *ServiceAttachmentConnectedEndpoints) Empty() bool {
	return r.empty
}

func (r *ServiceAttachmentConnectedEndpoints) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceAttachmentConnectedEndpoints) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceAttachmentConsumerAcceptLists struct {
	empty           bool    `json:"-"`
	ProjectIdOrNum  *string `json:"projectIdOrNum"`
	ConnectionLimit *int64  `json:"connectionLimit"`
}

type jsonServiceAttachmentConsumerAcceptLists ServiceAttachmentConsumerAcceptLists

func (r *ServiceAttachmentConsumerAcceptLists) UnmarshalJSON(data []byte) error {
	var res jsonServiceAttachmentConsumerAcceptLists
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceAttachmentConsumerAcceptLists
	} else {

		r.ProjectIdOrNum = res.ProjectIdOrNum

		r.ConnectionLimit = res.ConnectionLimit

	}
	return nil
}

// This object is used to assert a desired state where this ServiceAttachmentConsumerAcceptLists is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceAttachmentConsumerAcceptLists *ServiceAttachmentConsumerAcceptLists = &ServiceAttachmentConsumerAcceptLists{empty: true}

func (r *ServiceAttachmentConsumerAcceptLists) Empty() bool {
	return r.empty
}

func (r *ServiceAttachmentConsumerAcceptLists) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceAttachmentConsumerAcceptLists) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ServiceAttachmentPscServiceAttachmentId struct {
	empty bool   `json:"-"`
	High  *int64 `json:"high"`
	Low   *int64 `json:"low"`
}

type jsonServiceAttachmentPscServiceAttachmentId ServiceAttachmentPscServiceAttachmentId

func (r *ServiceAttachmentPscServiceAttachmentId) UnmarshalJSON(data []byte) error {
	var res jsonServiceAttachmentPscServiceAttachmentId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyServiceAttachmentPscServiceAttachmentId
	} else {

		r.High = res.High

		r.Low = res.Low

	}
	return nil
}

// This object is used to assert a desired state where this ServiceAttachmentPscServiceAttachmentId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyServiceAttachmentPscServiceAttachmentId *ServiceAttachmentPscServiceAttachmentId = &ServiceAttachmentPscServiceAttachmentId{empty: true}

func (r *ServiceAttachmentPscServiceAttachmentId) Empty() bool {
	return r.empty
}

func (r *ServiceAttachmentPscServiceAttachmentId) String() string {
	return dcl.SprintResource(r)
}

func (r *ServiceAttachmentPscServiceAttachmentId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *ServiceAttachment) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "ServiceAttachment",
		Version: "compute",
	}
}

func (r *ServiceAttachment) ID() (string, error) {
	if err := extractServiceAttachmentFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"id":                        dcl.ValueOrEmptyString(nr.Id),
		"name":                      dcl.ValueOrEmptyString(nr.Name),
		"description":               dcl.ValueOrEmptyString(nr.Description),
		"self_link":                 dcl.ValueOrEmptyString(nr.SelfLink),
		"region":                    dcl.ValueOrEmptyString(nr.Region),
		"target_service":            dcl.ValueOrEmptyString(nr.TargetService),
		"connection_preference":     dcl.ValueOrEmptyString(nr.ConnectionPreference),
		"connected_endpoints":       dcl.ValueOrEmptyString(nr.ConnectedEndpoints),
		"nat_subnets":               dcl.ValueOrEmptyString(nr.NatSubnets),
		"enable_proxy_protocol":     dcl.ValueOrEmptyString(nr.EnableProxyProtocol),
		"consumer_reject_lists":     dcl.ValueOrEmptyString(nr.ConsumerRejectLists),
		"consumer_accept_lists":     dcl.ValueOrEmptyString(nr.ConsumerAcceptLists),
		"psc_service_attachment_id": dcl.ValueOrEmptyString(nr.PscServiceAttachmentId),
		"fingerprint":               dcl.ValueOrEmptyString(nr.Fingerprint),
		"project":                   dcl.ValueOrEmptyString(nr.Project),
		"location":                  dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/regions/{{location}}/serviceAttachments/{{name}}", params), nil
}

const ServiceAttachmentMaxPage = -1

type ServiceAttachmentList struct {
	Items []*ServiceAttachment

	nextToken string

	pageSize int32

	resource *ServiceAttachment
}

func (l *ServiceAttachmentList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ServiceAttachmentList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listServiceAttachment(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListServiceAttachment(ctx context.Context, project, location string) (*ServiceAttachmentList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListServiceAttachmentWithMaxResults(ctx, project, location, ServiceAttachmentMaxPage)

}

func (c *Client) ListServiceAttachmentWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ServiceAttachmentList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &ServiceAttachment{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listServiceAttachment(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ServiceAttachmentList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetServiceAttachment(ctx context.Context, r *ServiceAttachment) (*ServiceAttachment, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractServiceAttachmentFields(r)

	b, err := c.getServiceAttachmentRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalServiceAttachment(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeServiceAttachmentNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractServiceAttachmentFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteServiceAttachment(ctx context.Context, r *ServiceAttachment) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("ServiceAttachment resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting ServiceAttachment...")
	deleteOp := deleteServiceAttachmentOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllServiceAttachment deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllServiceAttachment(ctx context.Context, project, location string, filter func(*ServiceAttachment) bool) error {
	listObj, err := c.ListServiceAttachment(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllServiceAttachment(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllServiceAttachment(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyServiceAttachment(ctx context.Context, rawDesired *ServiceAttachment, opts ...dcl.ApplyOption) (*ServiceAttachment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *ServiceAttachment
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyServiceAttachmentHelper(c, ctx, rawDesired, opts...)
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

func applyServiceAttachmentHelper(c *Client, ctx context.Context, rawDesired *ServiceAttachment, opts ...dcl.ApplyOption) (*ServiceAttachment, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyServiceAttachment...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractServiceAttachmentFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.serviceAttachmentDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToServiceAttachmentDiffs(c.Config, fieldDiffs, opts)
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
	var ops []serviceAttachmentApiOperation
	if create {
		ops = append(ops, &createServiceAttachmentOperation{})
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
	return applyServiceAttachmentDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyServiceAttachmentDiff(c *Client, ctx context.Context, desired *ServiceAttachment, rawDesired *ServiceAttachment, ops []serviceAttachmentApiOperation, opts ...dcl.ApplyOption) (*ServiceAttachment, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetServiceAttachment(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createServiceAttachmentOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapServiceAttachment(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeServiceAttachmentNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeServiceAttachmentNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeServiceAttachmentDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractServiceAttachmentFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractServiceAttachmentFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffServiceAttachment(c, newDesired, newState)
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
