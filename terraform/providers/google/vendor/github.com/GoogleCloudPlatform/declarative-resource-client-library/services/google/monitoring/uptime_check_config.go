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
package monitoring

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type UptimeCheckConfig struct {
	Name              *string                             `json:"name"`
	DisplayName       *string                             `json:"displayName"`
	MonitoredResource *UptimeCheckConfigMonitoredResource `json:"monitoredResource"`
	ResourceGroup     *UptimeCheckConfigResourceGroup     `json:"resourceGroup"`
	HttpCheck         *UptimeCheckConfigHttpCheck         `json:"httpCheck"`
	TcpCheck          *UptimeCheckConfigTcpCheck          `json:"tcpCheck"`
	Period            *string                             `json:"period"`
	Timeout           *string                             `json:"timeout"`
	ContentMatchers   []UptimeCheckConfigContentMatchers  `json:"contentMatchers"`
	SelectedRegions   []string                            `json:"selectedRegions"`
	Project           *string                             `json:"project"`
}

func (r *UptimeCheckConfig) String() string {
	return dcl.SprintResource(r)
}

// The enum UptimeCheckConfigResourceGroupResourceTypeEnum.
type UptimeCheckConfigResourceGroupResourceTypeEnum string

// UptimeCheckConfigResourceGroupResourceTypeEnumRef returns a *UptimeCheckConfigResourceGroupResourceTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func UptimeCheckConfigResourceGroupResourceTypeEnumRef(s string) *UptimeCheckConfigResourceGroupResourceTypeEnum {
	v := UptimeCheckConfigResourceGroupResourceTypeEnum(s)
	return &v
}

func (v UptimeCheckConfigResourceGroupResourceTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"RESOURCE_TYPE_UNSPECIFIED", "INSTANCE", "AWS_ELB_LOAD_BALANCER"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "UptimeCheckConfigResourceGroupResourceTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum UptimeCheckConfigHttpCheckRequestMethodEnum.
type UptimeCheckConfigHttpCheckRequestMethodEnum string

// UptimeCheckConfigHttpCheckRequestMethodEnumRef returns a *UptimeCheckConfigHttpCheckRequestMethodEnum with the value of string s
// If the empty string is provided, nil is returned.
func UptimeCheckConfigHttpCheckRequestMethodEnumRef(s string) *UptimeCheckConfigHttpCheckRequestMethodEnum {
	v := UptimeCheckConfigHttpCheckRequestMethodEnum(s)
	return &v
}

func (v UptimeCheckConfigHttpCheckRequestMethodEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"METHOD_UNSPECIFIED", "GET", "POST"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "UptimeCheckConfigHttpCheckRequestMethodEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum UptimeCheckConfigHttpCheckContentTypeEnum.
type UptimeCheckConfigHttpCheckContentTypeEnum string

// UptimeCheckConfigHttpCheckContentTypeEnumRef returns a *UptimeCheckConfigHttpCheckContentTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func UptimeCheckConfigHttpCheckContentTypeEnumRef(s string) *UptimeCheckConfigHttpCheckContentTypeEnum {
	v := UptimeCheckConfigHttpCheckContentTypeEnum(s)
	return &v
}

func (v UptimeCheckConfigHttpCheckContentTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TYPE_UNSPECIFIED", "URL_ENCODED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "UptimeCheckConfigHttpCheckContentTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum UptimeCheckConfigContentMatchersMatcherEnum.
type UptimeCheckConfigContentMatchersMatcherEnum string

// UptimeCheckConfigContentMatchersMatcherEnumRef returns a *UptimeCheckConfigContentMatchersMatcherEnum with the value of string s
// If the empty string is provided, nil is returned.
func UptimeCheckConfigContentMatchersMatcherEnumRef(s string) *UptimeCheckConfigContentMatchersMatcherEnum {
	v := UptimeCheckConfigContentMatchersMatcherEnum(s)
	return &v
}

func (v UptimeCheckConfigContentMatchersMatcherEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"CONTENT_MATCHER_OPTION_UNSPECIFIED", "CONTAINS_STRING", "NOT_CONTAINS_STRING", "MATCHES_REGEX", "NOT_MATCHES_REGEX"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "UptimeCheckConfigContentMatchersMatcherEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type UptimeCheckConfigMonitoredResource struct {
	empty        bool              `json:"-"`
	Type         *string           `json:"type"`
	FilterLabels map[string]string `json:"filterLabels"`
}

type jsonUptimeCheckConfigMonitoredResource UptimeCheckConfigMonitoredResource

func (r *UptimeCheckConfigMonitoredResource) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigMonitoredResource
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigMonitoredResource
	} else {

		r.Type = res.Type

		r.FilterLabels = res.FilterLabels

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigMonitoredResource is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigMonitoredResource *UptimeCheckConfigMonitoredResource = &UptimeCheckConfigMonitoredResource{empty: true}

func (r *UptimeCheckConfigMonitoredResource) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigMonitoredResource) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigMonitoredResource) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type UptimeCheckConfigResourceGroup struct {
	empty        bool                                            `json:"-"`
	GroupId      *string                                         `json:"groupId"`
	ResourceType *UptimeCheckConfigResourceGroupResourceTypeEnum `json:"resourceType"`
}

type jsonUptimeCheckConfigResourceGroup UptimeCheckConfigResourceGroup

func (r *UptimeCheckConfigResourceGroup) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigResourceGroup
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigResourceGroup
	} else {

		r.GroupId = res.GroupId

		r.ResourceType = res.ResourceType

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigResourceGroup is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigResourceGroup *UptimeCheckConfigResourceGroup = &UptimeCheckConfigResourceGroup{empty: true}

func (r *UptimeCheckConfigResourceGroup) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigResourceGroup) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigResourceGroup) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type UptimeCheckConfigHttpCheck struct {
	empty         bool                                         `json:"-"`
	RequestMethod *UptimeCheckConfigHttpCheckRequestMethodEnum `json:"requestMethod"`
	UseSsl        *bool                                        `json:"useSsl"`
	Path          *string                                      `json:"path"`
	Port          *int64                                       `json:"port"`
	AuthInfo      *UptimeCheckConfigHttpCheckAuthInfo          `json:"authInfo"`
	MaskHeaders   *bool                                        `json:"maskHeaders"`
	Headers       map[string]string                            `json:"headers"`
	ContentType   *UptimeCheckConfigHttpCheckContentTypeEnum   `json:"contentType"`
	ValidateSsl   *bool                                        `json:"validateSsl"`
	Body          *string                                      `json:"body"`
}

type jsonUptimeCheckConfigHttpCheck UptimeCheckConfigHttpCheck

func (r *UptimeCheckConfigHttpCheck) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigHttpCheck
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigHttpCheck
	} else {

		r.RequestMethod = res.RequestMethod

		r.UseSsl = res.UseSsl

		r.Path = res.Path

		r.Port = res.Port

		r.AuthInfo = res.AuthInfo

		r.MaskHeaders = res.MaskHeaders

		r.Headers = res.Headers

		r.ContentType = res.ContentType

		r.ValidateSsl = res.ValidateSsl

		r.Body = res.Body

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigHttpCheck is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigHttpCheck *UptimeCheckConfigHttpCheck = &UptimeCheckConfigHttpCheck{empty: true}

func (r *UptimeCheckConfigHttpCheck) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigHttpCheck) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigHttpCheck) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type UptimeCheckConfigHttpCheckAuthInfo struct {
	empty    bool    `json:"-"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type jsonUptimeCheckConfigHttpCheckAuthInfo UptimeCheckConfigHttpCheckAuthInfo

func (r *UptimeCheckConfigHttpCheckAuthInfo) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigHttpCheckAuthInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigHttpCheckAuthInfo
	} else {

		r.Username = res.Username

		r.Password = res.Password

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigHttpCheckAuthInfo is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigHttpCheckAuthInfo *UptimeCheckConfigHttpCheckAuthInfo = &UptimeCheckConfigHttpCheckAuthInfo{empty: true}

func (r *UptimeCheckConfigHttpCheckAuthInfo) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigHttpCheckAuthInfo) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigHttpCheckAuthInfo) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type UptimeCheckConfigTcpCheck struct {
	empty bool   `json:"-"`
	Port  *int64 `json:"port"`
}

type jsonUptimeCheckConfigTcpCheck UptimeCheckConfigTcpCheck

func (r *UptimeCheckConfigTcpCheck) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigTcpCheck
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigTcpCheck
	} else {

		r.Port = res.Port

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigTcpCheck is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigTcpCheck *UptimeCheckConfigTcpCheck = &UptimeCheckConfigTcpCheck{empty: true}

func (r *UptimeCheckConfigTcpCheck) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigTcpCheck) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigTcpCheck) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type UptimeCheckConfigContentMatchers struct {
	empty   bool                                         `json:"-"`
	Content *string                                      `json:"content"`
	Matcher *UptimeCheckConfigContentMatchersMatcherEnum `json:"matcher"`
}

type jsonUptimeCheckConfigContentMatchers UptimeCheckConfigContentMatchers

func (r *UptimeCheckConfigContentMatchers) UnmarshalJSON(data []byte) error {
	var res jsonUptimeCheckConfigContentMatchers
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyUptimeCheckConfigContentMatchers
	} else {

		r.Content = res.Content

		r.Matcher = res.Matcher

	}
	return nil
}

// This object is used to assert a desired state where this UptimeCheckConfigContentMatchers is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyUptimeCheckConfigContentMatchers *UptimeCheckConfigContentMatchers = &UptimeCheckConfigContentMatchers{empty: true}

func (r *UptimeCheckConfigContentMatchers) Empty() bool {
	return r.empty
}

func (r *UptimeCheckConfigContentMatchers) String() string {
	return dcl.SprintResource(r)
}

func (r *UptimeCheckConfigContentMatchers) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *UptimeCheckConfig) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "monitoring",
		Type:    "UptimeCheckConfig",
		Version: "monitoring",
	}
}

func (r *UptimeCheckConfig) ID() (string, error) {
	if err := extractUptimeCheckConfigFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":               dcl.ValueOrEmptyString(nr.Name),
		"display_name":       dcl.ValueOrEmptyString(nr.DisplayName),
		"monitored_resource": dcl.ValueOrEmptyString(nr.MonitoredResource),
		"resource_group":     dcl.ValueOrEmptyString(nr.ResourceGroup),
		"http_check":         dcl.ValueOrEmptyString(nr.HttpCheck),
		"tcp_check":          dcl.ValueOrEmptyString(nr.TcpCheck),
		"period":             dcl.ValueOrEmptyString(nr.Period),
		"timeout":            dcl.ValueOrEmptyString(nr.Timeout),
		"content_matchers":   dcl.ValueOrEmptyString(nr.ContentMatchers),
		"selected_regions":   dcl.ValueOrEmptyString(nr.SelectedRegions),
		"project":            dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/uptimeCheckConfigs/{{name}}", params), nil
}

const UptimeCheckConfigMaxPage = -1

type UptimeCheckConfigList struct {
	Items []*UptimeCheckConfig

	nextToken string

	pageSize int32

	resource *UptimeCheckConfig
}

func (l *UptimeCheckConfigList) HasNext() bool {
	return l.nextToken != ""
}

func (l *UptimeCheckConfigList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listUptimeCheckConfig(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListUptimeCheckConfig(ctx context.Context, project string) (*UptimeCheckConfigList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListUptimeCheckConfigWithMaxResults(ctx, project, UptimeCheckConfigMaxPage)

}

func (c *Client) ListUptimeCheckConfigWithMaxResults(ctx context.Context, project string, pageSize int32) (*UptimeCheckConfigList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &UptimeCheckConfig{
		Project: &project,
	}
	items, token, err := c.listUptimeCheckConfig(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &UptimeCheckConfigList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetUptimeCheckConfig(ctx context.Context, r *UptimeCheckConfig) (*UptimeCheckConfig, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractUptimeCheckConfigFields(r)

	b, err := c.getUptimeCheckConfigRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalUptimeCheckConfig(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name
	if dcl.IsZeroValue(result.Period) {
		result.Period = dcl.String("60s")
	}

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeUptimeCheckConfigNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractUptimeCheckConfigFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteUptimeCheckConfig(ctx context.Context, r *UptimeCheckConfig) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("UptimeCheckConfig resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting UptimeCheckConfig...")
	deleteOp := deleteUptimeCheckConfigOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllUptimeCheckConfig deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllUptimeCheckConfig(ctx context.Context, project string, filter func(*UptimeCheckConfig) bool) error {
	listObj, err := c.ListUptimeCheckConfig(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllUptimeCheckConfig(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllUptimeCheckConfig(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyUptimeCheckConfig(ctx context.Context, rawDesired *UptimeCheckConfig, opts ...dcl.ApplyOption) (*UptimeCheckConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *UptimeCheckConfig
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyUptimeCheckConfigHelper(c, ctx, rawDesired, opts...)
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

func applyUptimeCheckConfigHelper(c *Client, ctx context.Context, rawDesired *UptimeCheckConfig, opts ...dcl.ApplyOption) (*UptimeCheckConfig, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyUptimeCheckConfig...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractUptimeCheckConfigFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.uptimeCheckConfigDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToUptimeCheckConfigDiffs(c.Config, fieldDiffs, opts)
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
	var ops []uptimeCheckConfigApiOperation
	if create {
		ops = append(ops, &createUptimeCheckConfigOperation{})
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
	return applyUptimeCheckConfigDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyUptimeCheckConfigDiff(c *Client, ctx context.Context, desired *UptimeCheckConfig, rawDesired *UptimeCheckConfig, ops []uptimeCheckConfigApiOperation, opts ...dcl.ApplyOption) (*UptimeCheckConfig, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetUptimeCheckConfig(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createUptimeCheckConfigOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapUptimeCheckConfig(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeUptimeCheckConfigNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeUptimeCheckConfigNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeUptimeCheckConfigDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractUptimeCheckConfigFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractUptimeCheckConfigFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffUptimeCheckConfig(c, newDesired, newState)
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
