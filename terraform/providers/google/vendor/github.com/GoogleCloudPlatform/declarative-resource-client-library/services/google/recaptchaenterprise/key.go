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
package recaptchaenterprise

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Key struct {
	Name            *string             `json:"name"`
	DisplayName     *string             `json:"displayName"`
	WebSettings     *KeyWebSettings     `json:"webSettings"`
	AndroidSettings *KeyAndroidSettings `json:"androidSettings"`
	IosSettings     *KeyIosSettings     `json:"iosSettings"`
	Labels          map[string]string   `json:"labels"`
	CreateTime      *string             `json:"createTime"`
	TestingOptions  *KeyTestingOptions  `json:"testingOptions"`
	Project         *string             `json:"project"`
}

func (r *Key) String() string {
	return dcl.SprintResource(r)
}

// The enum KeyWebSettingsIntegrationTypeEnum.
type KeyWebSettingsIntegrationTypeEnum string

// KeyWebSettingsIntegrationTypeEnumRef returns a *KeyWebSettingsIntegrationTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func KeyWebSettingsIntegrationTypeEnumRef(s string) *KeyWebSettingsIntegrationTypeEnum {
	v := KeyWebSettingsIntegrationTypeEnum(s)
	return &v
}

func (v KeyWebSettingsIntegrationTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCORE", "CHECKBOX", "INVISIBLE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "KeyWebSettingsIntegrationTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum KeyWebSettingsChallengeSecurityPreferenceEnum.
type KeyWebSettingsChallengeSecurityPreferenceEnum string

// KeyWebSettingsChallengeSecurityPreferenceEnumRef returns a *KeyWebSettingsChallengeSecurityPreferenceEnum with the value of string s
// If the empty string is provided, nil is returned.
func KeyWebSettingsChallengeSecurityPreferenceEnumRef(s string) *KeyWebSettingsChallengeSecurityPreferenceEnum {
	v := KeyWebSettingsChallengeSecurityPreferenceEnum(s)
	return &v
}

func (v KeyWebSettingsChallengeSecurityPreferenceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"CHALLENGE_SECURITY_PREFERENCE_UNSPECIFIED", "USABILITY", "BALANCE", "SECURITY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "KeyWebSettingsChallengeSecurityPreferenceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum KeyTestingOptionsTestingChallengeEnum.
type KeyTestingOptionsTestingChallengeEnum string

// KeyTestingOptionsTestingChallengeEnumRef returns a *KeyTestingOptionsTestingChallengeEnum with the value of string s
// If the empty string is provided, nil is returned.
func KeyTestingOptionsTestingChallengeEnumRef(s string) *KeyTestingOptionsTestingChallengeEnum {
	v := KeyTestingOptionsTestingChallengeEnum(s)
	return &v
}

func (v KeyTestingOptionsTestingChallengeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TESTING_CHALLENGE_UNSPECIFIED", "NOCAPTCHA", "UNSOLVABLE_CHALLENGE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "KeyTestingOptionsTestingChallengeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type KeyWebSettings struct {
	empty                       bool                                           `json:"-"`
	AllowAllDomains             *bool                                          `json:"allowAllDomains"`
	AllowedDomains              []string                                       `json:"allowedDomains"`
	AllowAmpTraffic             *bool                                          `json:"allowAmpTraffic"`
	IntegrationType             *KeyWebSettingsIntegrationTypeEnum             `json:"integrationType"`
	ChallengeSecurityPreference *KeyWebSettingsChallengeSecurityPreferenceEnum `json:"challengeSecurityPreference"`
}

type jsonKeyWebSettings KeyWebSettings

func (r *KeyWebSettings) UnmarshalJSON(data []byte) error {
	var res jsonKeyWebSettings
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyWebSettings
	} else {

		r.AllowAllDomains = res.AllowAllDomains

		r.AllowedDomains = res.AllowedDomains

		r.AllowAmpTraffic = res.AllowAmpTraffic

		r.IntegrationType = res.IntegrationType

		r.ChallengeSecurityPreference = res.ChallengeSecurityPreference

	}
	return nil
}

// This object is used to assert a desired state where this KeyWebSettings is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyWebSettings *KeyWebSettings = &KeyWebSettings{empty: true}

func (r *KeyWebSettings) Empty() bool {
	return r.empty
}

func (r *KeyWebSettings) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyWebSettings) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyAndroidSettings struct {
	empty                bool     `json:"-"`
	AllowAllPackageNames *bool    `json:"allowAllPackageNames"`
	AllowedPackageNames  []string `json:"allowedPackageNames"`
}

type jsonKeyAndroidSettings KeyAndroidSettings

func (r *KeyAndroidSettings) UnmarshalJSON(data []byte) error {
	var res jsonKeyAndroidSettings
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyAndroidSettings
	} else {

		r.AllowAllPackageNames = res.AllowAllPackageNames

		r.AllowedPackageNames = res.AllowedPackageNames

	}
	return nil
}

// This object is used to assert a desired state where this KeyAndroidSettings is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyAndroidSettings *KeyAndroidSettings = &KeyAndroidSettings{empty: true}

func (r *KeyAndroidSettings) Empty() bool {
	return r.empty
}

func (r *KeyAndroidSettings) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyAndroidSettings) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyIosSettings struct {
	empty             bool     `json:"-"`
	AllowAllBundleIds *bool    `json:"allowAllBundleIds"`
	AllowedBundleIds  []string `json:"allowedBundleIds"`
}

type jsonKeyIosSettings KeyIosSettings

func (r *KeyIosSettings) UnmarshalJSON(data []byte) error {
	var res jsonKeyIosSettings
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyIosSettings
	} else {

		r.AllowAllBundleIds = res.AllowAllBundleIds

		r.AllowedBundleIds = res.AllowedBundleIds

	}
	return nil
}

// This object is used to assert a desired state where this KeyIosSettings is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyIosSettings *KeyIosSettings = &KeyIosSettings{empty: true}

func (r *KeyIosSettings) Empty() bool {
	return r.empty
}

func (r *KeyIosSettings) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyIosSettings) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type KeyTestingOptions struct {
	empty            bool                                   `json:"-"`
	TestingScore     *float64                               `json:"testingScore"`
	TestingChallenge *KeyTestingOptionsTestingChallengeEnum `json:"testingChallenge"`
}

type jsonKeyTestingOptions KeyTestingOptions

func (r *KeyTestingOptions) UnmarshalJSON(data []byte) error {
	var res jsonKeyTestingOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyKeyTestingOptions
	} else {

		r.TestingScore = res.TestingScore

		r.TestingChallenge = res.TestingChallenge

	}
	return nil
}

// This object is used to assert a desired state where this KeyTestingOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyKeyTestingOptions *KeyTestingOptions = &KeyTestingOptions{empty: true}

func (r *KeyTestingOptions) Empty() bool {
	return r.empty
}

func (r *KeyTestingOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *KeyTestingOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Key) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "recaptcha_enterprise",
		Type:    "Key",
		Version: "recaptchaenterprise",
	}
}

func (r *Key) ID() (string, error) {
	if err := extractKeyFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":             dcl.ValueOrEmptyString(nr.Name),
		"display_name":     dcl.ValueOrEmptyString(nr.DisplayName),
		"web_settings":     dcl.ValueOrEmptyString(nr.WebSettings),
		"android_settings": dcl.ValueOrEmptyString(nr.AndroidSettings),
		"ios_settings":     dcl.ValueOrEmptyString(nr.IosSettings),
		"labels":           dcl.ValueOrEmptyString(nr.Labels),
		"create_time":      dcl.ValueOrEmptyString(nr.CreateTime),
		"testing_options":  dcl.ValueOrEmptyString(nr.TestingOptions),
		"project":          dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/keys/{{name}}", params), nil
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

func (c *Client) GetKey(ctx context.Context, r *Key) (*Key, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractKeyFields(r)

	b, err := c.getKeyRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalKey(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeKeyNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractKeyFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
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
