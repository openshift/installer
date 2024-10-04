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

type InterconnectAttachment struct {
	Description             *string                                           `json:"description"`
	SelfLink                *string                                           `json:"selfLink"`
	Id                      *int64                                            `json:"id"`
	Name                    *string                                           `json:"name"`
	Interconnect            *string                                           `json:"interconnect"`
	Router                  *string                                           `json:"router"`
	Region                  *string                                           `json:"region"`
	Mtu                     *int64                                            `json:"mtu"`
	PrivateInterconnectInfo *InterconnectAttachmentPrivateInterconnectInfo    `json:"privateInterconnectInfo"`
	OperationalStatus       *InterconnectAttachmentOperationalStatusEnum      `json:"operationalStatus"`
	CloudRouterIPAddress    *string                                           `json:"cloudRouterIPAddress"`
	CustomerRouterIPAddress *string                                           `json:"customerRouterIPAddress"`
	Type                    *InterconnectAttachmentTypeEnum                   `json:"type"`
	PairingKey              *string                                           `json:"pairingKey"`
	AdminEnabled            *bool                                             `json:"adminEnabled"`
	VlanTag8021q            *int64                                            `json:"vlanTag8021q"`
	EdgeAvailabilityDomain  *InterconnectAttachmentEdgeAvailabilityDomainEnum `json:"edgeAvailabilityDomain"`
	CandidateSubnets        []string                                          `json:"candidateSubnets"`
	Bandwidth               *InterconnectAttachmentBandwidthEnum              `json:"bandwidth"`
	PartnerMetadata         *InterconnectAttachmentPartnerMetadata            `json:"partnerMetadata"`
	State                   *InterconnectAttachmentStateEnum                  `json:"state"`
	PartnerAsn              *int64                                            `json:"partnerAsn"`
	Encryption              *InterconnectAttachmentEncryptionEnum             `json:"encryption"`
	IpsecInternalAddresses  []string                                          `json:"ipsecInternalAddresses"`
	DataplaneVersion        *int64                                            `json:"dataplaneVersion"`
	SatisfiesPzs            *bool                                             `json:"satisfiesPzs"`
	Project                 *string                                           `json:"project"`
}

func (r *InterconnectAttachment) String() string {
	return dcl.SprintResource(r)
}

// The enum InterconnectAttachmentOperationalStatusEnum.
type InterconnectAttachmentOperationalStatusEnum string

// InterconnectAttachmentOperationalStatusEnumRef returns a *InterconnectAttachmentOperationalStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentOperationalStatusEnumRef(s string) *InterconnectAttachmentOperationalStatusEnum {
	v := InterconnectAttachmentOperationalStatusEnum(s)
	return &v
}

func (v InterconnectAttachmentOperationalStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"OS_ACTIVE", "OS_UNPROVISIONED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectAttachmentOperationalStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectAttachmentTypeEnum.
type InterconnectAttachmentTypeEnum string

// InterconnectAttachmentTypeEnumRef returns a *InterconnectAttachmentTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentTypeEnumRef(s string) *InterconnectAttachmentTypeEnum {
	v := InterconnectAttachmentTypeEnum(s)
	return &v
}

func (v InterconnectAttachmentTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PATH", "OTHER", "PARAMETER"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectAttachmentTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectAttachmentEdgeAvailabilityDomainEnum.
type InterconnectAttachmentEdgeAvailabilityDomainEnum string

// InterconnectAttachmentEdgeAvailabilityDomainEnumRef returns a *InterconnectAttachmentEdgeAvailabilityDomainEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentEdgeAvailabilityDomainEnumRef(s string) *InterconnectAttachmentEdgeAvailabilityDomainEnum {
	v := InterconnectAttachmentEdgeAvailabilityDomainEnum(s)
	return &v
}

func (v InterconnectAttachmentEdgeAvailabilityDomainEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"AVAILABILITY_DOMAIN_ANY", "AVAILABILITY_DOMAIN_1", "AVAILABILITY_DOMAIN_2"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectAttachmentEdgeAvailabilityDomainEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectAttachmentBandwidthEnum.
type InterconnectAttachmentBandwidthEnum string

// InterconnectAttachmentBandwidthEnumRef returns a *InterconnectAttachmentBandwidthEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentBandwidthEnumRef(s string) *InterconnectAttachmentBandwidthEnum {
	v := InterconnectAttachmentBandwidthEnum(s)
	return &v
}

func (v InterconnectAttachmentBandwidthEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"BPS_50M", "BPS_100M", "BPS_200M", "BPS_300M", "BPS_400M", "BPS_500M", "BPS_1G", "BPS_2G", "BPS_5G", "BPS_10G", "BPS_20G", "BPS_50G"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectAttachmentBandwidthEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectAttachmentStateEnum.
type InterconnectAttachmentStateEnum string

// InterconnectAttachmentStateEnumRef returns a *InterconnectAttachmentStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentStateEnumRef(s string) *InterconnectAttachmentStateEnum {
	v := InterconnectAttachmentStateEnum(s)
	return &v
}

func (v InterconnectAttachmentStateEnum) Validate() error {
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
		Enum:  "InterconnectAttachmentStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InterconnectAttachmentEncryptionEnum.
type InterconnectAttachmentEncryptionEnum string

// InterconnectAttachmentEncryptionEnumRef returns a *InterconnectAttachmentEncryptionEnum with the value of string s
// If the empty string is provided, nil is returned.
func InterconnectAttachmentEncryptionEnumRef(s string) *InterconnectAttachmentEncryptionEnum {
	v := InterconnectAttachmentEncryptionEnum(s)
	return &v
}

func (v InterconnectAttachmentEncryptionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"NONE", "IPSEC"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InterconnectAttachmentEncryptionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type InterconnectAttachmentPrivateInterconnectInfo struct {
	empty    bool   `json:"-"`
	Tag8021q *int64 `json:"tag8021q"`
}

type jsonInterconnectAttachmentPrivateInterconnectInfo InterconnectAttachmentPrivateInterconnectInfo

func (r *InterconnectAttachmentPrivateInterconnectInfo) UnmarshalJSON(data []byte) error {
	var res jsonInterconnectAttachmentPrivateInterconnectInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInterconnectAttachmentPrivateInterconnectInfo
	} else {

		r.Tag8021q = res.Tag8021q

	}
	return nil
}

// This object is used to assert a desired state where this InterconnectAttachmentPrivateInterconnectInfo is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInterconnectAttachmentPrivateInterconnectInfo *InterconnectAttachmentPrivateInterconnectInfo = &InterconnectAttachmentPrivateInterconnectInfo{empty: true}

func (r *InterconnectAttachmentPrivateInterconnectInfo) Empty() bool {
	return r.empty
}

func (r *InterconnectAttachmentPrivateInterconnectInfo) String() string {
	return dcl.SprintResource(r)
}

func (r *InterconnectAttachmentPrivateInterconnectInfo) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InterconnectAttachmentPartnerMetadata struct {
	empty            bool    `json:"-"`
	PartnerName      *string `json:"partnerName"`
	InterconnectName *string `json:"interconnectName"`
	PortalUrl        *string `json:"portalUrl"`
}

type jsonInterconnectAttachmentPartnerMetadata InterconnectAttachmentPartnerMetadata

func (r *InterconnectAttachmentPartnerMetadata) UnmarshalJSON(data []byte) error {
	var res jsonInterconnectAttachmentPartnerMetadata
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInterconnectAttachmentPartnerMetadata
	} else {

		r.PartnerName = res.PartnerName

		r.InterconnectName = res.InterconnectName

		r.PortalUrl = res.PortalUrl

	}
	return nil
}

// This object is used to assert a desired state where this InterconnectAttachmentPartnerMetadata is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInterconnectAttachmentPartnerMetadata *InterconnectAttachmentPartnerMetadata = &InterconnectAttachmentPartnerMetadata{empty: true}

func (r *InterconnectAttachmentPartnerMetadata) Empty() bool {
	return r.empty
}

func (r *InterconnectAttachmentPartnerMetadata) String() string {
	return dcl.SprintResource(r)
}

func (r *InterconnectAttachmentPartnerMetadata) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *InterconnectAttachment) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "InterconnectAttachment",
		Version: "compute",
	}
}

func (r *InterconnectAttachment) ID() (string, error) {
	if err := extractInterconnectAttachmentFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"description":                dcl.ValueOrEmptyString(nr.Description),
		"self_link":                  dcl.ValueOrEmptyString(nr.SelfLink),
		"id":                         dcl.ValueOrEmptyString(nr.Id),
		"name":                       dcl.ValueOrEmptyString(nr.Name),
		"interconnect":               dcl.ValueOrEmptyString(nr.Interconnect),
		"router":                     dcl.ValueOrEmptyString(nr.Router),
		"region":                     dcl.ValueOrEmptyString(nr.Region),
		"mtu":                        dcl.ValueOrEmptyString(nr.Mtu),
		"private_interconnect_info":  dcl.ValueOrEmptyString(nr.PrivateInterconnectInfo),
		"operational_status":         dcl.ValueOrEmptyString(nr.OperationalStatus),
		"cloud_router_ip_address":    dcl.ValueOrEmptyString(nr.CloudRouterIPAddress),
		"customer_router_ip_address": dcl.ValueOrEmptyString(nr.CustomerRouterIPAddress),
		"type":                       dcl.ValueOrEmptyString(nr.Type),
		"pairing_key":                dcl.ValueOrEmptyString(nr.PairingKey),
		"admin_enabled":              dcl.ValueOrEmptyString(nr.AdminEnabled),
		"vlan_tag8021q":              dcl.ValueOrEmptyString(nr.VlanTag8021q),
		"edge_availability_domain":   dcl.ValueOrEmptyString(nr.EdgeAvailabilityDomain),
		"candidate_subnets":          dcl.ValueOrEmptyString(nr.CandidateSubnets),
		"bandwidth":                  dcl.ValueOrEmptyString(nr.Bandwidth),
		"partner_metadata":           dcl.ValueOrEmptyString(nr.PartnerMetadata),
		"state":                      dcl.ValueOrEmptyString(nr.State),
		"partner_asn":                dcl.ValueOrEmptyString(nr.PartnerAsn),
		"encryption":                 dcl.ValueOrEmptyString(nr.Encryption),
		"ipsec_internal_addresses":   dcl.ValueOrEmptyString(nr.IpsecInternalAddresses),
		"dataplane_version":          dcl.ValueOrEmptyString(nr.DataplaneVersion),
		"satisfies_pzs":              dcl.ValueOrEmptyString(nr.SatisfiesPzs),
		"project":                    dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}", params), nil
}

const InterconnectAttachmentMaxPage = -1

type InterconnectAttachmentList struct {
	Items []*InterconnectAttachment

	nextToken string

	pageSize int32

	resource *InterconnectAttachment
}

func (l *InterconnectAttachmentList) HasNext() bool {
	return l.nextToken != ""
}

func (l *InterconnectAttachmentList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listInterconnectAttachment(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListInterconnectAttachment(ctx context.Context, project, region string) (*InterconnectAttachmentList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListInterconnectAttachmentWithMaxResults(ctx, project, region, InterconnectAttachmentMaxPage)

}

func (c *Client) ListInterconnectAttachmentWithMaxResults(ctx context.Context, project, region string, pageSize int32) (*InterconnectAttachmentList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &InterconnectAttachment{
		Project: &project,
		Region:  &region,
	}
	items, token, err := c.listInterconnectAttachment(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &InterconnectAttachmentList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetInterconnectAttachment(ctx context.Context, r *InterconnectAttachment) (*InterconnectAttachment, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractInterconnectAttachmentFields(r)

	b, err := c.getInterconnectAttachmentRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalInterconnectAttachment(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Region = r.Region
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeInterconnectAttachmentNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractInterconnectAttachmentFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteInterconnectAttachment(ctx context.Context, r *InterconnectAttachment) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("InterconnectAttachment resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting InterconnectAttachment...")
	deleteOp := deleteInterconnectAttachmentOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllInterconnectAttachment deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllInterconnectAttachment(ctx context.Context, project, region string, filter func(*InterconnectAttachment) bool) error {
	listObj, err := c.ListInterconnectAttachment(ctx, project, region)
	if err != nil {
		return err
	}

	err = c.deleteAllInterconnectAttachment(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllInterconnectAttachment(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyInterconnectAttachment(ctx context.Context, rawDesired *InterconnectAttachment, opts ...dcl.ApplyOption) (*InterconnectAttachment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *InterconnectAttachment
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyInterconnectAttachmentHelper(c, ctx, rawDesired, opts...)
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

func applyInterconnectAttachmentHelper(c *Client, ctx context.Context, rawDesired *InterconnectAttachment, opts ...dcl.ApplyOption) (*InterconnectAttachment, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyInterconnectAttachment...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractInterconnectAttachmentFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.interconnectAttachmentDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToInterconnectAttachmentDiffs(c.Config, fieldDiffs, opts)
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
	var ops []interconnectAttachmentApiOperation
	if create {
		ops = append(ops, &createInterconnectAttachmentOperation{})
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
	return applyInterconnectAttachmentDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyInterconnectAttachmentDiff(c *Client, ctx context.Context, desired *InterconnectAttachment, rawDesired *InterconnectAttachment, ops []interconnectAttachmentApiOperation, opts ...dcl.ApplyOption) (*InterconnectAttachment, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetInterconnectAttachment(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createInterconnectAttachmentOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapInterconnectAttachment(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeInterconnectAttachmentNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeInterconnectAttachmentNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeInterconnectAttachmentDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractInterconnectAttachmentFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractInterconnectAttachmentFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffInterconnectAttachment(c, newDesired, newState)
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
