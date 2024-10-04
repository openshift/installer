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

type ForwardingRule struct {
	Labels                        map[string]string                             `json:"labels"`
	AllPorts                      *bool                                         `json:"allPorts"`
	AllowGlobalAccess             *bool                                         `json:"allowGlobalAccess"`
	LabelFingerprint              *string                                       `json:"labelFingerprint"`
	BackendService                *string                                       `json:"backendService"`
	CreationTimestamp             *string                                       `json:"creationTimestamp"`
	Description                   *string                                       `json:"description"`
	IPAddress                     *string                                       `json:"ipAddress"`
	IPProtocol                    *ForwardingRuleIPProtocolEnum                 `json:"ipProtocol"`
	IPVersion                     *ForwardingRuleIPVersionEnum                  `json:"ipVersion"`
	IsMirroringCollector          *bool                                         `json:"isMirroringCollector"`
	LoadBalancingScheme           *ForwardingRuleLoadBalancingSchemeEnum        `json:"loadBalancingScheme"`
	MetadataFilter                []ForwardingRuleMetadataFilter                `json:"metadataFilter"`
	Name                          *string                                       `json:"name"`
	Network                       *string                                       `json:"network"`
	NetworkTier                   *ForwardingRuleNetworkTierEnum                `json:"networkTier"`
	PortRange                     *string                                       `json:"portRange"`
	Ports                         []string                                      `json:"ports"`
	Region                        *string                                       `json:"region"`
	SelfLink                      *string                                       `json:"selfLink"`
	ServiceLabel                  *string                                       `json:"serviceLabel"`
	ServiceName                   *string                                       `json:"serviceName"`
	Subnetwork                    *string                                       `json:"subnetwork"`
	Target                        *string                                       `json:"target"`
	Project                       *string                                       `json:"project"`
	Location                      *string                                       `json:"location"`
	ServiceDirectoryRegistrations []ForwardingRuleServiceDirectoryRegistrations `json:"serviceDirectoryRegistrations"`
	PscConnectionId               *string                                       `json:"pscConnectionId"`
	PscConnectionStatus           *ForwardingRulePscConnectionStatusEnum        `json:"pscConnectionStatus"`
	SourceIPRanges                []string                                      `json:"sourceIPRanges"`
	BaseForwardingRule            *string                                       `json:"baseForwardingRule"`
}

func (r *ForwardingRule) String() string {
	return dcl.SprintResource(r)
}

// The enum ForwardingRuleIPProtocolEnum.
type ForwardingRuleIPProtocolEnum string

// ForwardingRuleIPProtocolEnumRef returns a *ForwardingRuleIPProtocolEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRuleIPProtocolEnumRef(s string) *ForwardingRuleIPProtocolEnum {
	v := ForwardingRuleIPProtocolEnum(s)
	return &v
}

func (v ForwardingRuleIPProtocolEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TCP", "UDP", "ESP", "AH", "SCTP", "ICMP", "L3_DEFAULT"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRuleIPProtocolEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ForwardingRuleIPVersionEnum.
type ForwardingRuleIPVersionEnum string

// ForwardingRuleIPVersionEnumRef returns a *ForwardingRuleIPVersionEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRuleIPVersionEnumRef(s string) *ForwardingRuleIPVersionEnum {
	v := ForwardingRuleIPVersionEnum(s)
	return &v
}

func (v ForwardingRuleIPVersionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNSPECIFIED_VERSION", "IPV4", "IPV6"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRuleIPVersionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ForwardingRuleLoadBalancingSchemeEnum.
type ForwardingRuleLoadBalancingSchemeEnum string

// ForwardingRuleLoadBalancingSchemeEnumRef returns a *ForwardingRuleLoadBalancingSchemeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRuleLoadBalancingSchemeEnumRef(s string) *ForwardingRuleLoadBalancingSchemeEnum {
	v := ForwardingRuleLoadBalancingSchemeEnum(s)
	return &v
}

func (v ForwardingRuleLoadBalancingSchemeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"INVALID", "INTERNAL", "INTERNAL_MANAGED", "INTERNAL_SELF_MANAGED", "EXTERNAL", "EXTERNAL_MANAGED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRuleLoadBalancingSchemeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ForwardingRuleMetadataFilterFilterMatchCriteriaEnum.
type ForwardingRuleMetadataFilterFilterMatchCriteriaEnum string

// ForwardingRuleMetadataFilterFilterMatchCriteriaEnumRef returns a *ForwardingRuleMetadataFilterFilterMatchCriteriaEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRuleMetadataFilterFilterMatchCriteriaEnumRef(s string) *ForwardingRuleMetadataFilterFilterMatchCriteriaEnum {
	v := ForwardingRuleMetadataFilterFilterMatchCriteriaEnum(s)
	return &v
}

func (v ForwardingRuleMetadataFilterFilterMatchCriteriaEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"NOT_SET", "MATCH_ALL", "MATCH_ANY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRuleMetadataFilterFilterMatchCriteriaEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ForwardingRuleNetworkTierEnum.
type ForwardingRuleNetworkTierEnum string

// ForwardingRuleNetworkTierEnumRef returns a *ForwardingRuleNetworkTierEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRuleNetworkTierEnumRef(s string) *ForwardingRuleNetworkTierEnum {
	v := ForwardingRuleNetworkTierEnum(s)
	return &v
}

func (v ForwardingRuleNetworkTierEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PREMIUM", "STANDARD"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRuleNetworkTierEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ForwardingRulePscConnectionStatusEnum.
type ForwardingRulePscConnectionStatusEnum string

// ForwardingRulePscConnectionStatusEnumRef returns a *ForwardingRulePscConnectionStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func ForwardingRulePscConnectionStatusEnumRef(s string) *ForwardingRulePscConnectionStatusEnum {
	v := ForwardingRulePscConnectionStatusEnum(s)
	return &v
}

func (v ForwardingRulePscConnectionStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATUS_UNSPECIFIED", "PENDING", "ACCEPTED", "REJECTED", "CLOSED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ForwardingRulePscConnectionStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ForwardingRuleMetadataFilter struct {
	empty               bool                                                 `json:"-"`
	FilterMatchCriteria *ForwardingRuleMetadataFilterFilterMatchCriteriaEnum `json:"filterMatchCriteria"`
	FilterLabel         []ForwardingRuleMetadataFilterFilterLabel            `json:"filterLabel"`
}

type jsonForwardingRuleMetadataFilter ForwardingRuleMetadataFilter

func (r *ForwardingRuleMetadataFilter) UnmarshalJSON(data []byte) error {
	var res jsonForwardingRuleMetadataFilter
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyForwardingRuleMetadataFilter
	} else {

		r.FilterMatchCriteria = res.FilterMatchCriteria

		r.FilterLabel = res.FilterLabel

	}
	return nil
}

// This object is used to assert a desired state where this ForwardingRuleMetadataFilter is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyForwardingRuleMetadataFilter *ForwardingRuleMetadataFilter = &ForwardingRuleMetadataFilter{empty: true}

func (r *ForwardingRuleMetadataFilter) Empty() bool {
	return r.empty
}

func (r *ForwardingRuleMetadataFilter) String() string {
	return dcl.SprintResource(r)
}

func (r *ForwardingRuleMetadataFilter) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ForwardingRuleMetadataFilterFilterLabel struct {
	empty bool    `json:"-"`
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type jsonForwardingRuleMetadataFilterFilterLabel ForwardingRuleMetadataFilterFilterLabel

func (r *ForwardingRuleMetadataFilterFilterLabel) UnmarshalJSON(data []byte) error {
	var res jsonForwardingRuleMetadataFilterFilterLabel
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyForwardingRuleMetadataFilterFilterLabel
	} else {

		r.Name = res.Name

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this ForwardingRuleMetadataFilterFilterLabel is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyForwardingRuleMetadataFilterFilterLabel *ForwardingRuleMetadataFilterFilterLabel = &ForwardingRuleMetadataFilterFilterLabel{empty: true}

func (r *ForwardingRuleMetadataFilterFilterLabel) Empty() bool {
	return r.empty
}

func (r *ForwardingRuleMetadataFilterFilterLabel) String() string {
	return dcl.SprintResource(r)
}

func (r *ForwardingRuleMetadataFilterFilterLabel) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ForwardingRuleServiceDirectoryRegistrations struct {
	empty     bool    `json:"-"`
	Namespace *string `json:"namespace"`
	Service   *string `json:"service"`
}

type jsonForwardingRuleServiceDirectoryRegistrations ForwardingRuleServiceDirectoryRegistrations

func (r *ForwardingRuleServiceDirectoryRegistrations) UnmarshalJSON(data []byte) error {
	var res jsonForwardingRuleServiceDirectoryRegistrations
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyForwardingRuleServiceDirectoryRegistrations
	} else {

		r.Namespace = res.Namespace

		r.Service = res.Service

	}
	return nil
}

// This object is used to assert a desired state where this ForwardingRuleServiceDirectoryRegistrations is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyForwardingRuleServiceDirectoryRegistrations *ForwardingRuleServiceDirectoryRegistrations = &ForwardingRuleServiceDirectoryRegistrations{empty: true}

func (r *ForwardingRuleServiceDirectoryRegistrations) Empty() bool {
	return r.empty
}

func (r *ForwardingRuleServiceDirectoryRegistrations) String() string {
	return dcl.SprintResource(r)
}

func (r *ForwardingRuleServiceDirectoryRegistrations) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *ForwardingRule) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "ForwardingRule",
		Version: "compute",
	}
}

func (r *ForwardingRule) ID() (string, error) {
	if err := extractForwardingRuleFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"labels":                          dcl.ValueOrEmptyString(nr.Labels),
		"all_ports":                       dcl.ValueOrEmptyString(nr.AllPorts),
		"allow_global_access":             dcl.ValueOrEmptyString(nr.AllowGlobalAccess),
		"label_fingerprint":               dcl.ValueOrEmptyString(nr.LabelFingerprint),
		"backend_service":                 dcl.ValueOrEmptyString(nr.BackendService),
		"creation_timestamp":              dcl.ValueOrEmptyString(nr.CreationTimestamp),
		"description":                     dcl.ValueOrEmptyString(nr.Description),
		"ip_address":                      dcl.ValueOrEmptyString(nr.IPAddress),
		"ip_protocol":                     dcl.ValueOrEmptyString(nr.IPProtocol),
		"ip_version":                      dcl.ValueOrEmptyString(nr.IPVersion),
		"is_mirroring_collector":          dcl.ValueOrEmptyString(nr.IsMirroringCollector),
		"load_balancing_scheme":           dcl.ValueOrEmptyString(nr.LoadBalancingScheme),
		"metadata_filter":                 dcl.ValueOrEmptyString(nr.MetadataFilter),
		"name":                            dcl.ValueOrEmptyString(nr.Name),
		"network":                         dcl.ValueOrEmptyString(nr.Network),
		"network_tier":                    dcl.ValueOrEmptyString(nr.NetworkTier),
		"port_range":                      dcl.ValueOrEmptyString(nr.PortRange),
		"ports":                           dcl.ValueOrEmptyString(nr.Ports),
		"region":                          dcl.ValueOrEmptyString(nr.Region),
		"self_link":                       dcl.ValueOrEmptyString(nr.SelfLink),
		"service_label":                   dcl.ValueOrEmptyString(nr.ServiceLabel),
		"service_name":                    dcl.ValueOrEmptyString(nr.ServiceName),
		"subnetwork":                      dcl.ValueOrEmptyString(nr.Subnetwork),
		"target":                          dcl.ValueOrEmptyString(nr.Target),
		"project":                         dcl.ValueOrEmptyString(nr.Project),
		"location":                        dcl.ValueOrEmptyString(nr.Location),
		"service_directory_registrations": dcl.ValueOrEmptyString(nr.ServiceDirectoryRegistrations),
		"psc_connection_id":               dcl.ValueOrEmptyString(nr.PscConnectionId),
		"psc_connection_status":           dcl.ValueOrEmptyString(nr.PscConnectionStatus),
		"source_ip_ranges":                dcl.ValueOrEmptyString(nr.SourceIPRanges),
		"base_forwarding_rule":            dcl.ValueOrEmptyString(nr.BaseForwardingRule),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.Nprintf("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}", params), nil
	}

	return dcl.Nprintf("projects/{{project}}/global/forwardingRules/{{name}}", params), nil
}

const ForwardingRuleMaxPage = -1

type ForwardingRuleList struct {
	Items []*ForwardingRule

	nextToken string

	pageSize int32

	resource *ForwardingRule
}

func (l *ForwardingRuleList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ForwardingRuleList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listForwardingRule(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListForwardingRule(ctx context.Context, project, location string) (*ForwardingRuleList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListForwardingRuleWithMaxResults(ctx, project, location, ForwardingRuleMaxPage)

}

func (c *Client) ListForwardingRuleWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ForwardingRuleList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &ForwardingRule{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listForwardingRule(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ForwardingRuleList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetForwardingRule(ctx context.Context, r *ForwardingRule) (*ForwardingRule, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractForwardingRuleFields(r)

	b, err := c.getForwardingRuleRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalForwardingRule(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeForwardingRuleNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractForwardingRuleFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteForwardingRule(ctx context.Context, r *ForwardingRule) error {
	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("ForwardingRule resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting ForwardingRule...")
	deleteOp := deleteForwardingRuleOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllForwardingRule deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllForwardingRule(ctx context.Context, project, location string, filter func(*ForwardingRule) bool) error {
	listObj, err := c.ListForwardingRule(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllForwardingRule(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllForwardingRule(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyForwardingRule(ctx context.Context, rawDesired *ForwardingRule, opts ...dcl.ApplyOption) (*ForwardingRule, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	c = NewClient(c.Config.Clone(dcl.WithCodeRetryability(map[int]dcl.Retryability{
		412: dcl.Retryability{
			Retryable: false,
			Pattern:   "",
			Timeout:   0,
		},
	})))
	var resultNewState *ForwardingRule
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyForwardingRuleHelper(c, ctx, rawDesired, opts...)
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

func applyForwardingRuleHelper(c *Client, ctx context.Context, rawDesired *ForwardingRule, opts ...dcl.ApplyOption) (*ForwardingRule, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyForwardingRule...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractForwardingRuleFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.forwardingRuleDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToForwardingRuleDiffs(c.Config, fieldDiffs, opts)
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
	var ops []forwardingRuleApiOperation
	if create {
		ops = append(ops, &createForwardingRuleOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	ops, err = forwardingRuleSetLabelsPostCreate(ops)
	if err != nil {
		return nil, err
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
	return applyForwardingRuleDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyForwardingRuleDiff(c *Client, ctx context.Context, desired *ForwardingRule, rawDesired *ForwardingRule, ops []forwardingRuleApiOperation, opts ...dcl.ApplyOption) (*ForwardingRule, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetForwardingRule(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createForwardingRuleOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapForwardingRule(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeForwardingRuleNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeForwardingRuleNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeForwardingRuleDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractForwardingRuleFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractForwardingRuleFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffForwardingRule(c, newDesired, newState)
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
