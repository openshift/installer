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
package eventarc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Channel struct {
	Name               *string           `json:"name"`
	Uid                *string           `json:"uid"`
	CreateTime         *string           `json:"createTime"`
	UpdateTime         *string           `json:"updateTime"`
	ThirdPartyProvider *string           `json:"thirdPartyProvider"`
	PubsubTopic        *string           `json:"pubsubTopic"`
	State              *ChannelStateEnum `json:"state"`
	ActivationToken    *string           `json:"activationToken"`
	CryptoKeyName      *string           `json:"cryptoKeyName"`
	Project            *string           `json:"project"`
	Location           *string           `json:"location"`
}

func (r *Channel) String() string {
	return dcl.SprintResource(r)
}

// The enum ChannelStateEnum.
type ChannelStateEnum string

// ChannelStateEnumRef returns a *ChannelStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ChannelStateEnumRef(s string) *ChannelStateEnum {
	v := ChannelStateEnum(s)
	return &v
}

func (v ChannelStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "PENDING", "ACTIVE", "INACTIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ChannelStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Channel) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "eventarc",
		Type:    "Channel",
		Version: "eventarc",
	}
}

func (r *Channel) ID() (string, error) {
	if err := extractChannelFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                 dcl.ValueOrEmptyString(nr.Name),
		"uid":                  dcl.ValueOrEmptyString(nr.Uid),
		"create_time":          dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":          dcl.ValueOrEmptyString(nr.UpdateTime),
		"third_party_provider": dcl.ValueOrEmptyString(nr.ThirdPartyProvider),
		"pubsub_topic":         dcl.ValueOrEmptyString(nr.PubsubTopic),
		"state":                dcl.ValueOrEmptyString(nr.State),
		"activation_token":     dcl.ValueOrEmptyString(nr.ActivationToken),
		"crypto_key_name":      dcl.ValueOrEmptyString(nr.CryptoKeyName),
		"project":              dcl.ValueOrEmptyString(nr.Project),
		"location":             dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/channels/{{name}}", params), nil
}

const ChannelMaxPage = -1

type ChannelList struct {
	Items []*Channel

	nextToken string

	pageSize int32

	resource *Channel
}

func (l *ChannelList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ChannelList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listChannel(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListChannel(ctx context.Context, project, location string) (*ChannelList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "",
			Timeout:   10000000000,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListChannelWithMaxResults(ctx, project, location, ChannelMaxPage)

}

func (c *Client) ListChannelWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ChannelList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Channel{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listChannel(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ChannelList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetChannel(ctx context.Context, r *Channel) (*Channel, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "",
			Timeout:   10000000000,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractChannelFields(r)

	b, err := c.getChannelRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalChannel(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeChannelNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractChannelFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteChannel(ctx context.Context, r *Channel) error {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "",
			Timeout:   10000000000,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Channel resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Channel...")
	deleteOp := deleteChannelOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllChannel deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllChannel(ctx context.Context, project, location string, filter func(*Channel) bool) error {
	listObj, err := c.ListChannel(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllChannel(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllChannel(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyChannel(ctx context.Context, rawDesired *Channel, opts ...dcl.ApplyOption) (*Channel, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		403: dcl.Retryability{
			Retryable: true,
			Pattern:   "",
			Timeout:   10000000000,
		},
	})))
	var resultNewState *Channel
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyChannelHelper(c, ctx, rawDesired, opts...)
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

func applyChannelHelper(c *Client, ctx context.Context, rawDesired *Channel, opts ...dcl.ApplyOption) (*Channel, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyChannel...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractChannelFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.channelDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToChannelDiffs(c.Config, fieldDiffs, opts)
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
	var ops []channelApiOperation
	if create {
		ops = append(ops, &createChannelOperation{})
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
	return applyChannelDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyChannelDiff(c *Client, ctx context.Context, desired *Channel, rawDesired *Channel, ops []channelApiOperation, opts ...dcl.ApplyOption) (*Channel, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetChannel(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createChannelOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapChannel(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeChannelNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeChannelNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeChannelDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractChannelFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractChannelFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffChannel(c, newDesired, newState)
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
