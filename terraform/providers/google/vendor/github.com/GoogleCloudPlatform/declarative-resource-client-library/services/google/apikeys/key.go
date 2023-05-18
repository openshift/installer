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
package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Key struct {
	Name         *string          `json:"name"`
	DisplayName  *string          `json:"displayName"`
	KeyString    *string          `json:"keyString"`
	Uid          *string          `json:"uid"`
	Restrictions *KeyRestrictions `json:"restrictions"`
	Project      *string          `json:"project"`
}

func (r *Key) String() string {
	return dcl.SprintResource(r)
}

type KeyRestrictions struct {
	empty                  bool                                   `json:"-"`
	BrowserKeyRestrictions *KeyRestrictionsBrowserKeyRestrictions `json:"browserKeyRestrictions"`
	ServerKeyRestrictions  *KeyRestrictionsServerKeyRestrictions  `json:"serverKeyRestrictions"`
	AndroidKeyRestrictions *KeyRestrictionsAndroidKeyRestrictions `json:"androidKeyRestrictions"`
	IosKeyRestrictions     *KeyRestrictionsIosKeyRestrictions     `json:"iosKeyRestrictions"`
	ApiTargets             []KeyRestrictionsApiTargets            `json:"apiTargets"`
}

type jsonKeyRestrictions KeyRestrictions

func (r *KeyRestrictions) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictions
	} else {

		r.BrowserKeyRestrictions = res.BrowserKeyRestrictions

		r.ServerKeyRestrictions = res.ServerKeyRestrictions

		r.AndroidKeyRestrictions = res.AndroidKeyRestrictions

		r.IosKeyRestrictions = res.IosKeyRestrictions

		r.ApiTargets = res.ApiTargets

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictions *KeyRestrictions = &KeyRestrictions{empty: true}

func (r *KeyRestrictions) Empty() bool {
	return r.empty
}

func (r *KeyRestrictions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsBrowserKeyRestrictions struct {
	empty            bool     `json:"-"`
	AllowedReferrers []string `json:"allowedReferrers"`
}

type jsonKeyRestrictionsBrowserKeyRestrictions KeyRestrictionsBrowserKeyRestrictions

func (r *KeyRestrictionsBrowserKeyRestrictions) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsBrowserKeyRestrictions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsBrowserKeyRestrictions
	} else {

		r.AllowedReferrers = res.AllowedReferrers

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsBrowserKeyRestrictions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsBrowserKeyRestrictions *KeyRestrictionsBrowserKeyRestrictions = &KeyRestrictionsBrowserKeyRestrictions{empty: true}

func (r *KeyRestrictionsBrowserKeyRestrictions) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsBrowserKeyRestrictions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsBrowserKeyRestrictions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsServerKeyRestrictions struct {
	empty      bool     `json:"-"`
	AllowedIps []string `json:"allowedIps"`
}

type jsonKeyRestrictionsServerKeyRestrictions KeyRestrictionsServerKeyRestrictions

func (r *KeyRestrictionsServerKeyRestrictions) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsServerKeyRestrictions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsServerKeyRestrictions
	} else {

		r.AllowedIps = res.AllowedIps

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsServerKeyRestrictions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsServerKeyRestrictions *KeyRestrictionsServerKeyRestrictions = &KeyRestrictionsServerKeyRestrictions{empty: true}

func (r *KeyRestrictionsServerKeyRestrictions) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsServerKeyRestrictions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsServerKeyRestrictions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsAndroidKeyRestrictions struct {
	empty               bool                                                       `json:"-"`
	AllowedApplications []KeyRestrictionsAndroidKeyRestrictionsAllowedApplications `json:"allowedApplications"`
}

type jsonKeyRestrictionsAndroidKeyRestrictions KeyRestrictionsAndroidKeyRestrictions

func (r *KeyRestrictionsAndroidKeyRestrictions) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsAndroidKeyRestrictions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsAndroidKeyRestrictions
	} else {

		r.AllowedApplications = res.AllowedApplications

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsAndroidKeyRestrictions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsAndroidKeyRestrictions *KeyRestrictionsAndroidKeyRestrictions = &KeyRestrictionsAndroidKeyRestrictions{empty: true}

func (r *KeyRestrictionsAndroidKeyRestrictions) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsAndroidKeyRestrictions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsAndroidKeyRestrictions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsAndroidKeyRestrictionsAllowedApplications struct {
	empty           bool    `json:"-"`
	Sha1Fingerprint *string `json:"sha1Fingerprint"`
	PackageName     *string `json:"packageName"`
}

type jsonKeyRestrictionsAndroidKeyRestrictionsAllowedApplications KeyRestrictionsAndroidKeyRestrictionsAllowedApplications

func (r *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsAndroidKeyRestrictionsAllowedApplications
	} else {

		r.Sha1Fingerprint = res.Sha1Fingerprint

		r.PackageName = res.PackageName

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsAndroidKeyRestrictionsAllowedApplications is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsAndroidKeyRestrictionsAllowedApplications *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications = &KeyRestrictionsAndroidKeyRestrictionsAllowedApplications{empty: true}

func (r *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsAndroidKeyRestrictionsAllowedApplications) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsIosKeyRestrictions struct {
	empty            bool     `json:"-"`
	AllowedBundleIds []string `json:"allowedBundleIds"`
}

type jsonKeyRestrictionsIosKeyRestrictions KeyRestrictionsIosKeyRestrictions

func (r *KeyRestrictionsIosKeyRestrictions) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsIosKeyRestrictions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsIosKeyRestrictions
	} else {

		r.AllowedBundleIds = res.AllowedBundleIds

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsIosKeyRestrictions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsIosKeyRestrictions *KeyRestrictionsIosKeyRestrictions = &KeyRestrictionsIosKeyRestrictions{empty: true}

func (r *KeyRestrictionsIosKeyRestrictions) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsIosKeyRestrictions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsIosKeyRestrictions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyRestrictionsApiTargets struct {
	empty   bool     `json:"-"`
	Service *string  `json:"service"`
	Methods []string `json:"methods"`
}

type jsonKeyRestrictionsApiTargets KeyRestrictionsApiTargets

func (r *KeyRestrictionsApiTargets) UnmarshalJSON(data []byte) error {
	var res jsonKeyRestrictionsApiTargets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyRestrictionsApiTargets
	} else {

		r.Service = res.Service

		r.Methods = res.Methods

	}
	return nil
}

// This object is used to assert a desired state where this KeyRestrictionsApiTargets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyRestrictionsApiTargets *KeyRestrictionsApiTargets = &KeyRestrictionsApiTargets{empty: true}

func (r *KeyRestrictionsApiTargets) Empty() bool {
	return r.empty
}

func (r *KeyRestrictionsApiTargets) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyRestrictionsApiTargets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Key) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "apikeys",
		Type:    "Key",
		Version: "apikeys",
	}
}

func (r *Key) ID() (string, error) {
	if err := extractKeyFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":         dcl.ValueOrEmptyString(nr.Name),
		"display_name": dcl.ValueOrEmptyString(nr.DisplayName),
		"key_string":   dcl.ValueOrEmptyString(nr.KeyString),
		"uid":          dcl.ValueOrEmptyString(nr.Uid),
		"restrictions": dcl.ValueOrEmptyString(nr.Restrictions),
		"project":      dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/locations/global/keys/{{name}}", params), nil
}

const KeyMaxPage = -1

type KeyList struct {
	Items []*Key

	nextToken string

	pageSize int32

	resource *Key
}

func (l *KeyList) HasNext() bool {
	return l.nextToken != ""
}

func (l *KeyList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listKey(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListKey(ctx context.Context, project string) (*KeyList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListKeyWithMaxResults(ctx, project, KeyMaxPage)

}

func (c *Client) ListKeyWithMaxResults(ctx context.Context, project string, pageSize int32) (*KeyList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Key{
		Project: &project,
	}
	items, token, err := c.listKey(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &KeyList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) DeleteKey(ctx context.Context, r *Key) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Key resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Key...")
	deleteOp := deleteKeyOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllKey deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllKey(ctx context.Context, project string, filter func(*Key) bool) error {
	listObj, err := c.ListKey(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllKey(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllKey(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyKey(ctx context.Context, rawDesired *Key, opts ...dcl.ApplyOption) (*Key, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Key
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyKeyHelper(c, ctx, rawDesired, opts...)
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

func applyKeyHelper(c *Client, ctx context.Context, rawDesired *Key, opts ...dcl.ApplyOption) (*Key, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyKey...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractKeyFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.keyDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToKeyDiffs(c.Config, fieldDiffs, opts)
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
	var ops []keyApiOperation
	if create {
		ops = append(ops, &createKeyOperation{})
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
	return applyKeyDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyKeyDiff(c *Client, ctx context.Context, desired *Key, rawDesired *Key, ops []keyApiOperation, opts ...dcl.ApplyOption) (*Key, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetKey(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createKeyOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapKey(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeKeyNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeKeyNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeKeyDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractKeyFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractKeyFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffKey(c, newDesired, newState)
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
