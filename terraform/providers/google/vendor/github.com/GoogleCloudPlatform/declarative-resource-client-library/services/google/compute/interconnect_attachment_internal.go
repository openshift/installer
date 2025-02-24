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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *InterconnectAttachment) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "region"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.PrivateInterconnectInfo) {
		if err := r.PrivateInterconnectInfo.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PartnerMetadata) {
		if err := r.PartnerMetadata.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *InterconnectAttachmentPrivateInterconnectInfo) validate() error {
	return nil
}
func (r *InterconnectAttachmentPartnerMetadata) validate() error {
	return nil
}
func (r *InterconnectAttachment) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *InterconnectAttachment) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *InterconnectAttachment) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/interconnectAttachments", nr.basePath(), userBasePath, params), nil

}

func (r *InterconnectAttachment) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/interconnectAttachments", nr.basePath(), userBasePath, params), nil

}

func (r *InterconnectAttachment) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}", nr.basePath(), userBasePath, params), nil
}

// interconnectAttachmentApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type interconnectAttachmentApiOperation interface {
	do(context.Context, *InterconnectAttachment, *Client) error
}

// newUpdateInterconnectAttachmentPatchRequest creates a request for an
// InterconnectAttachment resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInterconnectAttachmentPatchRequest(ctx context.Context, f *InterconnectAttachment, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		req["name"] = v
	}
	if v := f.Interconnect; !dcl.IsEmptyValueIndirect(v) {
		req["interconnect"] = v
	}
	if v := f.Router; !dcl.IsEmptyValueIndirect(v) {
		req["router"] = v
	}
	if v := f.Mtu; !dcl.IsEmptyValueIndirect(v) {
		req["mtu"] = v
	}
	if v := f.Type; !dcl.IsEmptyValueIndirect(v) {
		req["type"] = v
	}
	if v := f.PairingKey; !dcl.IsEmptyValueIndirect(v) {
		req["pairingKey"] = v
	}
	if v := f.AdminEnabled; !dcl.IsEmptyValueIndirect(v) {
		req["adminEnabled"] = v
	}
	if v := f.VlanTag8021q; !dcl.IsEmptyValueIndirect(v) {
		req["vlanTag8021q"] = v
	}
	if v := f.EdgeAvailabilityDomain; !dcl.IsEmptyValueIndirect(v) {
		req["edgeAvailabilityDomain"] = v
	}
	if v := f.CandidateSubnets; v != nil {
		req["candidateSubnets"] = v
	}
	if v := f.Bandwidth; !dcl.IsEmptyValueIndirect(v) {
		req["bandwidth"] = v
	}
	if v, err := expandInterconnectAttachmentPartnerMetadata(c, f.PartnerMetadata, res); err != nil {
		return nil, fmt.Errorf("error expanding PartnerMetadata into partnerMetadata: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["partnerMetadata"] = v
	}
	if v := f.PartnerAsn; !dcl.IsEmptyValueIndirect(v) {
		req["partnerAsn"] = v
	}
	if v := f.Encryption; !dcl.IsEmptyValueIndirect(v) {
		req["encryption"] = v
	}
	if v := f.IpsecInternalAddresses; v != nil {
		req["ipsecInternalAddresses"] = v
	}
	if v := f.DataplaneVersion; !dcl.IsEmptyValueIndirect(v) {
		req["dataplaneVersion"] = v
	}
	return req, nil
}

// marshalUpdateInterconnectAttachmentPatchRequest converts the update into
// the final JSON request body.
func marshalUpdateInterconnectAttachmentPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateInterconnectAttachmentPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInterconnectAttachmentPatchOperation) do(ctx context.Context, r *InterconnectAttachment, c *Client) error {
	_, err := c.GetInterconnectAttachment(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdateInterconnectAttachmentPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInterconnectAttachmentPatchRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listInterconnectAttachmentRaw(ctx context.Context, r *InterconnectAttachment, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != InterconnectAttachmentMaxPage {
		m["pageSize"] = fmt.Sprintf("%v", pageSize)
	}

	u, err = dcl.AddQueryParams(u, m)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	return ioutil.ReadAll(resp.Response.Body)
}

type listInterconnectAttachmentOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listInterconnectAttachment(ctx context.Context, r *InterconnectAttachment, pageToken string, pageSize int32) ([]*InterconnectAttachment, string, error) {
	b, err := c.listInterconnectAttachmentRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listInterconnectAttachmentOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*InterconnectAttachment
	for _, v := range m.Items {
		res, err := unmarshalMapInterconnectAttachment(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Region = r.Region
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllInterconnectAttachment(ctx context.Context, f func(*InterconnectAttachment) bool, resources []*InterconnectAttachment) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteInterconnectAttachment(ctx, res)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", strings.Join(errors, "\n"))
	} else {
		return nil
	}
}

type deleteInterconnectAttachmentOperation struct{}

func (op *deleteInterconnectAttachmentOperation) do(ctx context.Context, r *InterconnectAttachment, c *Client) error {
	r, err := c.GetInterconnectAttachment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "InterconnectAttachment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetInterconnectAttachment checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete InterconnectAttachment: %w", err)
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetInterconnectAttachment(ctx, r)
		if dcl.IsNotFound(err) {
			return nil, nil
		}
		if retriesRemaining > 0 {
			retriesRemaining--
			return &dcl.RetryDetails{}, dcl.OperationNotDone{}
		}
		return nil, dcl.NotDeletedError{ExistingResource: r}
	}, c.Config.RetryProvider)
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createInterconnectAttachmentOperation struct {
	response map[string]interface{}
}

func (op *createInterconnectAttachmentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createInterconnectAttachmentOperation) do(ctx context.Context, r *InterconnectAttachment, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(req), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	o, err := dcl.ResponseBodyAsJSON(resp)
	if err != nil {
		return fmt.Errorf("error decoding response body into JSON: %w", err)
	}
	op.response = o

	if _, err := c.GetInterconnectAttachment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getInterconnectAttachmentRaw(ctx context.Context, r *InterconnectAttachment) ([]byte, error) {

	u, err := r.getURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	b, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) interconnectAttachmentDiffsForRawDesired(ctx context.Context, rawDesired *InterconnectAttachment, opts ...dcl.ApplyOption) (initial, desired *InterconnectAttachment, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *InterconnectAttachment
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*InterconnectAttachment); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected InterconnectAttachment, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetInterconnectAttachment(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a InterconnectAttachment resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve InterconnectAttachment resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that InterconnectAttachment resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeInterconnectAttachmentDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for InterconnectAttachment: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for InterconnectAttachment: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractInterconnectAttachmentFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeInterconnectAttachmentInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for InterconnectAttachment: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeInterconnectAttachmentDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for InterconnectAttachment: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffInterconnectAttachment(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeInterconnectAttachmentInitialState(rawInitial, rawDesired *InterconnectAttachment) (*InterconnectAttachment, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeInterconnectAttachmentDesiredState(rawDesired, rawInitial *InterconnectAttachment, opts ...dcl.ApplyOption) (*InterconnectAttachment, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.PrivateInterconnectInfo = canonicalizeInterconnectAttachmentPrivateInterconnectInfo(rawDesired.PrivateInterconnectInfo, nil, opts...)
		rawDesired.PartnerMetadata = canonicalizeInterconnectAttachmentPartnerMetadata(rawDesired.PartnerMetadata, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &InterconnectAttachment{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Interconnect, rawInitial.Interconnect) {
		canonicalDesired.Interconnect = rawInitial.Interconnect
	} else {
		canonicalDesired.Interconnect = rawDesired.Interconnect
	}
	if dcl.StringCanonicalize(rawDesired.Router, rawInitial.Router) {
		canonicalDesired.Router = rawInitial.Router
	} else {
		canonicalDesired.Router = rawDesired.Router
	}
	if dcl.StringCanonicalize(rawDesired.Region, rawInitial.Region) {
		canonicalDesired.Region = rawInitial.Region
	} else {
		canonicalDesired.Region = rawDesired.Region
	}
	if dcl.IsZeroValue(rawDesired.Mtu) || (dcl.IsEmptyValueIndirect(rawDesired.Mtu) && dcl.IsEmptyValueIndirect(rawInitial.Mtu)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Mtu = rawInitial.Mtu
	} else {
		canonicalDesired.Mtu = rawDesired.Mtu
	}
	if dcl.IsZeroValue(rawDesired.Type) || (dcl.IsEmptyValueIndirect(rawDesired.Type) && dcl.IsEmptyValueIndirect(rawInitial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Type = rawInitial.Type
	} else {
		canonicalDesired.Type = rawDesired.Type
	}
	if dcl.StringCanonicalize(rawDesired.PairingKey, rawInitial.PairingKey) {
		canonicalDesired.PairingKey = rawInitial.PairingKey
	} else {
		canonicalDesired.PairingKey = rawDesired.PairingKey
	}
	if dcl.BoolCanonicalize(rawDesired.AdminEnabled, rawInitial.AdminEnabled) {
		canonicalDesired.AdminEnabled = rawInitial.AdminEnabled
	} else {
		canonicalDesired.AdminEnabled = rawDesired.AdminEnabled
	}
	if dcl.IsZeroValue(rawDesired.VlanTag8021q) || (dcl.IsEmptyValueIndirect(rawDesired.VlanTag8021q) && dcl.IsEmptyValueIndirect(rawInitial.VlanTag8021q)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.VlanTag8021q = rawInitial.VlanTag8021q
	} else {
		canonicalDesired.VlanTag8021q = rawDesired.VlanTag8021q
	}
	if dcl.IsZeroValue(rawDesired.EdgeAvailabilityDomain) || (dcl.IsEmptyValueIndirect(rawDesired.EdgeAvailabilityDomain) && dcl.IsEmptyValueIndirect(rawInitial.EdgeAvailabilityDomain)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.EdgeAvailabilityDomain = rawInitial.EdgeAvailabilityDomain
	} else {
		canonicalDesired.EdgeAvailabilityDomain = rawDesired.EdgeAvailabilityDomain
	}
	if dcl.StringArrayCanonicalize(rawDesired.CandidateSubnets, rawInitial.CandidateSubnets) {
		canonicalDesired.CandidateSubnets = rawInitial.CandidateSubnets
	} else {
		canonicalDesired.CandidateSubnets = rawDesired.CandidateSubnets
	}
	if dcl.IsZeroValue(rawDesired.Bandwidth) || (dcl.IsEmptyValueIndirect(rawDesired.Bandwidth) && dcl.IsEmptyValueIndirect(rawInitial.Bandwidth)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Bandwidth = rawInitial.Bandwidth
	} else {
		canonicalDesired.Bandwidth = rawDesired.Bandwidth
	}
	canonicalDesired.PartnerMetadata = canonicalizeInterconnectAttachmentPartnerMetadata(rawDesired.PartnerMetadata, rawInitial.PartnerMetadata, opts...)
	if dcl.IsZeroValue(rawDesired.PartnerAsn) || (dcl.IsEmptyValueIndirect(rawDesired.PartnerAsn) && dcl.IsEmptyValueIndirect(rawInitial.PartnerAsn)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.PartnerAsn = rawInitial.PartnerAsn
	} else {
		canonicalDesired.PartnerAsn = rawDesired.PartnerAsn
	}
	if dcl.IsZeroValue(rawDesired.Encryption) || (dcl.IsEmptyValueIndirect(rawDesired.Encryption) && dcl.IsEmptyValueIndirect(rawInitial.Encryption)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Encryption = rawInitial.Encryption
	} else {
		canonicalDesired.Encryption = rawDesired.Encryption
	}
	if dcl.StringArrayCanonicalize(rawDesired.IpsecInternalAddresses, rawInitial.IpsecInternalAddresses) {
		canonicalDesired.IpsecInternalAddresses = rawInitial.IpsecInternalAddresses
	} else {
		canonicalDesired.IpsecInternalAddresses = rawDesired.IpsecInternalAddresses
	}
	if dcl.IsZeroValue(rawDesired.DataplaneVersion) || (dcl.IsEmptyValueIndirect(rawDesired.DataplaneVersion) && dcl.IsEmptyValueIndirect(rawInitial.DataplaneVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.DataplaneVersion = rawInitial.DataplaneVersion
	} else {
		canonicalDesired.DataplaneVersion = rawDesired.DataplaneVersion
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeInterconnectAttachmentNewState(c *Client, rawNew, rawDesired *InterconnectAttachment) (*InterconnectAttachment, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Interconnect) && dcl.IsEmptyValueIndirect(rawDesired.Interconnect) {
		rawNew.Interconnect = rawDesired.Interconnect
	} else {
		if dcl.StringCanonicalize(rawDesired.Interconnect, rawNew.Interconnect) {
			rawNew.Interconnect = rawDesired.Interconnect
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Router) && dcl.IsEmptyValueIndirect(rawDesired.Router) {
		rawNew.Router = rawDesired.Router
	} else {
		if dcl.StringCanonicalize(rawDesired.Router, rawNew.Router) {
			rawNew.Router = rawDesired.Router
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Region) && dcl.IsEmptyValueIndirect(rawDesired.Region) {
		rawNew.Region = rawDesired.Region
	} else {
		if dcl.StringCanonicalize(rawDesired.Region, rawNew.Region) {
			rawNew.Region = rawDesired.Region
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Mtu) && dcl.IsEmptyValueIndirect(rawDesired.Mtu) {
		rawNew.Mtu = rawDesired.Mtu
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PrivateInterconnectInfo) && dcl.IsEmptyValueIndirect(rawDesired.PrivateInterconnectInfo) {
		rawNew.PrivateInterconnectInfo = rawDesired.PrivateInterconnectInfo
	} else {
		rawNew.PrivateInterconnectInfo = canonicalizeNewInterconnectAttachmentPrivateInterconnectInfo(c, rawDesired.PrivateInterconnectInfo, rawNew.PrivateInterconnectInfo)
	}

	if dcl.IsEmptyValueIndirect(rawNew.OperationalStatus) && dcl.IsEmptyValueIndirect(rawDesired.OperationalStatus) {
		rawNew.OperationalStatus = rawDesired.OperationalStatus
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CloudRouterIPAddress) && dcl.IsEmptyValueIndirect(rawDesired.CloudRouterIPAddress) {
		rawNew.CloudRouterIPAddress = rawDesired.CloudRouterIPAddress
	} else {
		if dcl.StringCanonicalize(rawDesired.CloudRouterIPAddress, rawNew.CloudRouterIPAddress) {
			rawNew.CloudRouterIPAddress = rawDesired.CloudRouterIPAddress
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CustomerRouterIPAddress) && dcl.IsEmptyValueIndirect(rawDesired.CustomerRouterIPAddress) {
		rawNew.CustomerRouterIPAddress = rawDesired.CustomerRouterIPAddress
	} else {
		if dcl.StringCanonicalize(rawDesired.CustomerRouterIPAddress, rawNew.CustomerRouterIPAddress) {
			rawNew.CustomerRouterIPAddress = rawDesired.CustomerRouterIPAddress
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Type) && dcl.IsEmptyValueIndirect(rawDesired.Type) {
		rawNew.Type = rawDesired.Type
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PairingKey) && dcl.IsEmptyValueIndirect(rawDesired.PairingKey) {
		rawNew.PairingKey = rawDesired.PairingKey
	} else {
		if dcl.StringCanonicalize(rawDesired.PairingKey, rawNew.PairingKey) {
			rawNew.PairingKey = rawDesired.PairingKey
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.AdminEnabled) && dcl.IsEmptyValueIndirect(rawDesired.AdminEnabled) {
		rawNew.AdminEnabled = rawDesired.AdminEnabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.AdminEnabled, rawNew.AdminEnabled) {
			rawNew.AdminEnabled = rawDesired.AdminEnabled
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.VlanTag8021q) && dcl.IsEmptyValueIndirect(rawDesired.VlanTag8021q) {
		rawNew.VlanTag8021q = rawDesired.VlanTag8021q
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.EdgeAvailabilityDomain) && dcl.IsEmptyValueIndirect(rawDesired.EdgeAvailabilityDomain) {
		rawNew.EdgeAvailabilityDomain = rawDesired.EdgeAvailabilityDomain
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CandidateSubnets) && dcl.IsEmptyValueIndirect(rawDesired.CandidateSubnets) {
		rawNew.CandidateSubnets = rawDesired.CandidateSubnets
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.CandidateSubnets, rawNew.CandidateSubnets) {
			rawNew.CandidateSubnets = rawDesired.CandidateSubnets
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Bandwidth) && dcl.IsEmptyValueIndirect(rawDesired.Bandwidth) {
		rawNew.Bandwidth = rawDesired.Bandwidth
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PartnerMetadata) && dcl.IsEmptyValueIndirect(rawDesired.PartnerMetadata) {
		rawNew.PartnerMetadata = rawDesired.PartnerMetadata
	} else {
		rawNew.PartnerMetadata = canonicalizeNewInterconnectAttachmentPartnerMetadata(c, rawDesired.PartnerMetadata, rawNew.PartnerMetadata)
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PartnerAsn) && dcl.IsEmptyValueIndirect(rawDesired.PartnerAsn) {
		rawNew.PartnerAsn = rawDesired.PartnerAsn
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Encryption) && dcl.IsEmptyValueIndirect(rawDesired.Encryption) {
		rawNew.Encryption = rawDesired.Encryption
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.IpsecInternalAddresses) && dcl.IsEmptyValueIndirect(rawDesired.IpsecInternalAddresses) {
		rawNew.IpsecInternalAddresses = rawDesired.IpsecInternalAddresses
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.IpsecInternalAddresses, rawNew.IpsecInternalAddresses) {
			rawNew.IpsecInternalAddresses = rawDesired.IpsecInternalAddresses
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DataplaneVersion) && dcl.IsEmptyValueIndirect(rawDesired.DataplaneVersion) {
		rawNew.DataplaneVersion = rawDesired.DataplaneVersion
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SatisfiesPzs) && dcl.IsEmptyValueIndirect(rawDesired.SatisfiesPzs) {
		rawNew.SatisfiesPzs = rawDesired.SatisfiesPzs
	} else {
		if dcl.BoolCanonicalize(rawDesired.SatisfiesPzs, rawNew.SatisfiesPzs) {
			rawNew.SatisfiesPzs = rawDesired.SatisfiesPzs
		}
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeInterconnectAttachmentPrivateInterconnectInfo(des, initial *InterconnectAttachmentPrivateInterconnectInfo, opts ...dcl.ApplyOption) *InterconnectAttachmentPrivateInterconnectInfo {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InterconnectAttachmentPrivateInterconnectInfo{}

	return cDes
}

func canonicalizeInterconnectAttachmentPrivateInterconnectInfoSlice(des, initial []InterconnectAttachmentPrivateInterconnectInfo, opts ...dcl.ApplyOption) []InterconnectAttachmentPrivateInterconnectInfo {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InterconnectAttachmentPrivateInterconnectInfo, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInterconnectAttachmentPrivateInterconnectInfo(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InterconnectAttachmentPrivateInterconnectInfo, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInterconnectAttachmentPrivateInterconnectInfo(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInterconnectAttachmentPrivateInterconnectInfo(c *Client, des, nw *InterconnectAttachmentPrivateInterconnectInfo) *InterconnectAttachmentPrivateInterconnectInfo {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InterconnectAttachmentPrivateInterconnectInfo while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewInterconnectAttachmentPrivateInterconnectInfoSet(c *Client, des, nw []InterconnectAttachmentPrivateInterconnectInfo) []InterconnectAttachmentPrivateInterconnectInfo {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InterconnectAttachmentPrivateInterconnectInfo
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInterconnectAttachmentPrivateInterconnectInfoNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInterconnectAttachmentPrivateInterconnectInfo(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInterconnectAttachmentPrivateInterconnectInfoSlice(c *Client, des, nw []InterconnectAttachmentPrivateInterconnectInfo) []InterconnectAttachmentPrivateInterconnectInfo {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InterconnectAttachmentPrivateInterconnectInfo
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInterconnectAttachmentPrivateInterconnectInfo(c, &d, &n))
	}

	return items
}

func canonicalizeInterconnectAttachmentPartnerMetadata(des, initial *InterconnectAttachmentPartnerMetadata, opts ...dcl.ApplyOption) *InterconnectAttachmentPartnerMetadata {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InterconnectAttachmentPartnerMetadata{}

	if dcl.StringCanonicalize(des.PartnerName, initial.PartnerName) || dcl.IsZeroValue(des.PartnerName) {
		cDes.PartnerName = initial.PartnerName
	} else {
		cDes.PartnerName = des.PartnerName
	}
	if dcl.StringCanonicalize(des.InterconnectName, initial.InterconnectName) || dcl.IsZeroValue(des.InterconnectName) {
		cDes.InterconnectName = initial.InterconnectName
	} else {
		cDes.InterconnectName = des.InterconnectName
	}
	if dcl.StringCanonicalize(des.PortalUrl, initial.PortalUrl) || dcl.IsZeroValue(des.PortalUrl) {
		cDes.PortalUrl = initial.PortalUrl
	} else {
		cDes.PortalUrl = des.PortalUrl
	}

	return cDes
}

func canonicalizeInterconnectAttachmentPartnerMetadataSlice(des, initial []InterconnectAttachmentPartnerMetadata, opts ...dcl.ApplyOption) []InterconnectAttachmentPartnerMetadata {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InterconnectAttachmentPartnerMetadata, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInterconnectAttachmentPartnerMetadata(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InterconnectAttachmentPartnerMetadata, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInterconnectAttachmentPartnerMetadata(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInterconnectAttachmentPartnerMetadata(c *Client, des, nw *InterconnectAttachmentPartnerMetadata) *InterconnectAttachmentPartnerMetadata {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InterconnectAttachmentPartnerMetadata while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.PartnerName, nw.PartnerName) {
		nw.PartnerName = des.PartnerName
	}
	if dcl.StringCanonicalize(des.InterconnectName, nw.InterconnectName) {
		nw.InterconnectName = des.InterconnectName
	}
	if dcl.StringCanonicalize(des.PortalUrl, nw.PortalUrl) {
		nw.PortalUrl = des.PortalUrl
	}

	return nw
}

func canonicalizeNewInterconnectAttachmentPartnerMetadataSet(c *Client, des, nw []InterconnectAttachmentPartnerMetadata) []InterconnectAttachmentPartnerMetadata {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InterconnectAttachmentPartnerMetadata
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInterconnectAttachmentPartnerMetadataNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInterconnectAttachmentPartnerMetadata(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInterconnectAttachmentPartnerMetadataSlice(c *Client, des, nw []InterconnectAttachmentPartnerMetadata) []InterconnectAttachmentPartnerMetadata {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InterconnectAttachmentPartnerMetadata
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInterconnectAttachmentPartnerMetadata(c, &d, &n))
	}

	return items
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffInterconnectAttachment(c *Client, desired, actual *InterconnectAttachment, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interconnect, actual.Interconnect, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Interconnect")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Router, actual.Router, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Router")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Mtu, actual.Mtu, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Mtu")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateInterconnectInfo, actual.PrivateInterconnectInfo, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareInterconnectAttachmentPrivateInterconnectInfoNewStyle, EmptyObject: EmptyInterconnectAttachmentPrivateInterconnectInfo, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrivateInterconnectInfo")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationalStatus, actual.OperationalStatus, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OperationalStatus")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CloudRouterIPAddress, actual.CloudRouterIPAddress, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CloudRouterIpAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CustomerRouterIPAddress, actual.CustomerRouterIPAddress, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CustomerRouterIpAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PairingKey, actual.PairingKey, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("PairingKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdminEnabled, actual.AdminEnabled, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("AdminEnabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VlanTag8021q, actual.VlanTag8021q, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("VlanTag8021q")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EdgeAvailabilityDomain, actual.EdgeAvailabilityDomain, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("EdgeAvailabilityDomain")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CandidateSubnets, actual.CandidateSubnets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("CandidateSubnets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Bandwidth, actual.Bandwidth, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Bandwidth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PartnerMetadata, actual.PartnerMetadata, dcl.DiffInfo{ObjectFunction: compareInterconnectAttachmentPartnerMetadataNewStyle, EmptyObject: EmptyInterconnectAttachmentPartnerMetadata, OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("PartnerMetadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PartnerAsn, actual.PartnerAsn, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("PartnerAsn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Encryption, actual.Encryption, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("Encryption")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IpsecInternalAddresses, actual.IpsecInternalAddresses, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("IpsecInternalAddresses")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataplaneVersion, actual.DataplaneVersion, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("DataplaneVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SatisfiesPzs, actual.SatisfiesPzs, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SatisfiesPzs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if len(newDiffs) > 0 {
		c.Config.Logger.Infof("Diff function found diffs: %v", newDiffs)
	}
	return newDiffs, nil
}
func compareInterconnectAttachmentPrivateInterconnectInfoNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InterconnectAttachmentPrivateInterconnectInfo)
	if !ok {
		desiredNotPointer, ok := d.(InterconnectAttachmentPrivateInterconnectInfo)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectAttachmentPrivateInterconnectInfo or *InterconnectAttachmentPrivateInterconnectInfo", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InterconnectAttachmentPrivateInterconnectInfo)
	if !ok {
		actualNotPointer, ok := a.(InterconnectAttachmentPrivateInterconnectInfo)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectAttachmentPrivateInterconnectInfo", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Tag8021q, actual.Tag8021q, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tag8021q")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInterconnectAttachmentPartnerMetadataNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InterconnectAttachmentPartnerMetadata)
	if !ok {
		desiredNotPointer, ok := d.(InterconnectAttachmentPartnerMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectAttachmentPartnerMetadata or *InterconnectAttachmentPartnerMetadata", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InterconnectAttachmentPartnerMetadata)
	if !ok {
		actualNotPointer, ok := a.(InterconnectAttachmentPartnerMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectAttachmentPartnerMetadata", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PartnerName, actual.PartnerName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("PartnerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InterconnectName, actual.InterconnectName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("InterconnectName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PortalUrl, actual.PortalUrl, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInterconnectAttachmentPatchOperation")}, fn.AddNest("PortalUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *InterconnectAttachment) urlNormalized() *InterconnectAttachment {
	normalized := dcl.Copy(*r).(InterconnectAttachment)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Interconnect = dcl.SelfLinkToName(r.Interconnect)
	normalized.Router = dcl.SelfLinkToName(r.Router)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.CloudRouterIPAddress = dcl.SelfLinkToName(r.CloudRouterIPAddress)
	normalized.CustomerRouterIPAddress = dcl.SelfLinkToName(r.CustomerRouterIPAddress)
	normalized.PairingKey = dcl.SelfLinkToName(r.PairingKey)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *InterconnectAttachment) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"region":  dcl.ValueOrEmptyString(nr.Region),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the InterconnectAttachment resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *InterconnectAttachment) marshal(c *Client) ([]byte, error) {
	m, err := expandInterconnectAttachment(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling InterconnectAttachment: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalInterconnectAttachment decodes JSON responses into the InterconnectAttachment resource schema.
func unmarshalInterconnectAttachment(b []byte, c *Client, res *InterconnectAttachment) (*InterconnectAttachment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapInterconnectAttachment(m, c, res)
}

func unmarshalMapInterconnectAttachment(m map[string]interface{}, c *Client, res *InterconnectAttachment) (*InterconnectAttachment, error) {

	flattened := flattenInterconnectAttachment(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandInterconnectAttachment expands InterconnectAttachment into a JSON request object.
func expandInterconnectAttachment(c *Client, f *InterconnectAttachment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Interconnect; dcl.ValueShouldBeSent(v) {
		m["interconnect"] = v
	}
	if v := f.Router; dcl.ValueShouldBeSent(v) {
		m["router"] = v
	}
	if v := f.Region; dcl.ValueShouldBeSent(v) {
		m["region"] = v
	}
	if v := f.Mtu; dcl.ValueShouldBeSent(v) {
		m["mtu"] = v
	}
	if v := f.Type; dcl.ValueShouldBeSent(v) {
		m["type"] = v
	}
	if v := f.PairingKey; dcl.ValueShouldBeSent(v) {
		m["pairingKey"] = v
	}
	if v := f.AdminEnabled; dcl.ValueShouldBeSent(v) {
		m["adminEnabled"] = v
	}
	if v := f.VlanTag8021q; dcl.ValueShouldBeSent(v) {
		m["vlanTag8021q"] = v
	}
	if v := f.EdgeAvailabilityDomain; dcl.ValueShouldBeSent(v) {
		m["edgeAvailabilityDomain"] = v
	}
	if v := f.CandidateSubnets; v != nil {
		m["candidateSubnets"] = v
	}
	if v := f.Bandwidth; dcl.ValueShouldBeSent(v) {
		m["bandwidth"] = v
	}
	if v, err := expandInterconnectAttachmentPartnerMetadata(c, f.PartnerMetadata, res); err != nil {
		return nil, fmt.Errorf("error expanding PartnerMetadata into partnerMetadata: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["partnerMetadata"] = v
	}
	if v := f.PartnerAsn; dcl.ValueShouldBeSent(v) {
		m["partnerAsn"] = v
	}
	if v := f.Encryption; dcl.ValueShouldBeSent(v) {
		m["encryption"] = v
	}
	if v := f.IpsecInternalAddresses; v != nil {
		m["ipsecInternalAddresses"] = v
	}
	if v := f.DataplaneVersion; dcl.ValueShouldBeSent(v) {
		m["dataplaneVersion"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenInterconnectAttachment flattens InterconnectAttachment from a JSON request object into the
// InterconnectAttachment type.
func flattenInterconnectAttachment(c *Client, i interface{}, res *InterconnectAttachment) *InterconnectAttachment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &InterconnectAttachment{}
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.Id = dcl.FlattenInteger(m["id"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Interconnect = dcl.FlattenString(m["interconnect"])
	resultRes.Router = dcl.FlattenString(m["router"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.Mtu = dcl.FlattenInteger(m["mtu"])
	resultRes.PrivateInterconnectInfo = flattenInterconnectAttachmentPrivateInterconnectInfo(c, m["privateInterconnectInfo"], res)
	resultRes.OperationalStatus = flattenInterconnectAttachmentOperationalStatusEnum(m["operationalStatus"])
	resultRes.CloudRouterIPAddress = dcl.FlattenString(m["cloudRouterIpAddress"])
	resultRes.CustomerRouterIPAddress = dcl.FlattenString(m["customerRouterIpAddress"])
	resultRes.Type = flattenInterconnectAttachmentTypeEnum(m["type"])
	resultRes.PairingKey = dcl.FlattenString(m["pairingKey"])
	resultRes.AdminEnabled = dcl.FlattenBool(m["adminEnabled"])
	resultRes.VlanTag8021q = dcl.FlattenInteger(m["vlanTag8021q"])
	resultRes.EdgeAvailabilityDomain = flattenInterconnectAttachmentEdgeAvailabilityDomainEnum(m["edgeAvailabilityDomain"])
	resultRes.CandidateSubnets = dcl.FlattenStringSlice(m["candidateSubnets"])
	resultRes.Bandwidth = flattenInterconnectAttachmentBandwidthEnum(m["bandwidth"])
	resultRes.PartnerMetadata = flattenInterconnectAttachmentPartnerMetadata(c, m["partnerMetadata"], res)
	resultRes.State = flattenInterconnectAttachmentStateEnum(m["state"])
	resultRes.PartnerAsn = dcl.FlattenInteger(m["partnerAsn"])
	resultRes.Encryption = flattenInterconnectAttachmentEncryptionEnum(m["encryption"])
	resultRes.IpsecInternalAddresses = dcl.FlattenStringSlice(m["ipsecInternalAddresses"])
	resultRes.DataplaneVersion = dcl.FlattenInteger(m["dataplaneVersion"])
	resultRes.SatisfiesPzs = dcl.FlattenBool(m["satisfiesPzs"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandInterconnectAttachmentPrivateInterconnectInfoMap expands the contents of InterconnectAttachmentPrivateInterconnectInfo into a JSON
// request object.
func expandInterconnectAttachmentPrivateInterconnectInfoMap(c *Client, f map[string]InterconnectAttachmentPrivateInterconnectInfo, res *InterconnectAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInterconnectAttachmentPrivateInterconnectInfo(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInterconnectAttachmentPrivateInterconnectInfoSlice expands the contents of InterconnectAttachmentPrivateInterconnectInfo into a JSON
// request object.
func expandInterconnectAttachmentPrivateInterconnectInfoSlice(c *Client, f []InterconnectAttachmentPrivateInterconnectInfo, res *InterconnectAttachment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInterconnectAttachmentPrivateInterconnectInfo(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInterconnectAttachmentPrivateInterconnectInfoMap flattens the contents of InterconnectAttachmentPrivateInterconnectInfo from a JSON
// response object.
func flattenInterconnectAttachmentPrivateInterconnectInfoMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentPrivateInterconnectInfo {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentPrivateInterconnectInfo{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentPrivateInterconnectInfo{}
	}

	items := make(map[string]InterconnectAttachmentPrivateInterconnectInfo)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentPrivateInterconnectInfo(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInterconnectAttachmentPrivateInterconnectInfoSlice flattens the contents of InterconnectAttachmentPrivateInterconnectInfo from a JSON
// response object.
func flattenInterconnectAttachmentPrivateInterconnectInfoSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentPrivateInterconnectInfo {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentPrivateInterconnectInfo{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentPrivateInterconnectInfo{}
	}

	items := make([]InterconnectAttachmentPrivateInterconnectInfo, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentPrivateInterconnectInfo(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInterconnectAttachmentPrivateInterconnectInfo expands an instance of InterconnectAttachmentPrivateInterconnectInfo into a JSON
// request object.
func expandInterconnectAttachmentPrivateInterconnectInfo(c *Client, f *InterconnectAttachmentPrivateInterconnectInfo, res *InterconnectAttachment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenInterconnectAttachmentPrivateInterconnectInfo flattens an instance of InterconnectAttachmentPrivateInterconnectInfo from a JSON
// response object.
func flattenInterconnectAttachmentPrivateInterconnectInfo(c *Client, i interface{}, res *InterconnectAttachment) *InterconnectAttachmentPrivateInterconnectInfo {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InterconnectAttachmentPrivateInterconnectInfo{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInterconnectAttachmentPrivateInterconnectInfo
	}
	r.Tag8021q = dcl.FlattenInteger(m["tag8021q"])

	return r
}

// expandInterconnectAttachmentPartnerMetadataMap expands the contents of InterconnectAttachmentPartnerMetadata into a JSON
// request object.
func expandInterconnectAttachmentPartnerMetadataMap(c *Client, f map[string]InterconnectAttachmentPartnerMetadata, res *InterconnectAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInterconnectAttachmentPartnerMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInterconnectAttachmentPartnerMetadataSlice expands the contents of InterconnectAttachmentPartnerMetadata into a JSON
// request object.
func expandInterconnectAttachmentPartnerMetadataSlice(c *Client, f []InterconnectAttachmentPartnerMetadata, res *InterconnectAttachment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInterconnectAttachmentPartnerMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInterconnectAttachmentPartnerMetadataMap flattens the contents of InterconnectAttachmentPartnerMetadata from a JSON
// response object.
func flattenInterconnectAttachmentPartnerMetadataMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentPartnerMetadata {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentPartnerMetadata{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentPartnerMetadata{}
	}

	items := make(map[string]InterconnectAttachmentPartnerMetadata)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentPartnerMetadata(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInterconnectAttachmentPartnerMetadataSlice flattens the contents of InterconnectAttachmentPartnerMetadata from a JSON
// response object.
func flattenInterconnectAttachmentPartnerMetadataSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentPartnerMetadata {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentPartnerMetadata{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentPartnerMetadata{}
	}

	items := make([]InterconnectAttachmentPartnerMetadata, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentPartnerMetadata(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInterconnectAttachmentPartnerMetadata expands an instance of InterconnectAttachmentPartnerMetadata into a JSON
// request object.
func expandInterconnectAttachmentPartnerMetadata(c *Client, f *InterconnectAttachmentPartnerMetadata, res *InterconnectAttachment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.PartnerName; !dcl.IsEmptyValueIndirect(v) {
		m["partnerName"] = v
	}
	if v := f.InterconnectName; !dcl.IsEmptyValueIndirect(v) {
		m["interconnectName"] = v
	}
	if v := f.PortalUrl; !dcl.IsEmptyValueIndirect(v) {
		m["portalUrl"] = v
	}

	return m, nil
}

// flattenInterconnectAttachmentPartnerMetadata flattens an instance of InterconnectAttachmentPartnerMetadata from a JSON
// response object.
func flattenInterconnectAttachmentPartnerMetadata(c *Client, i interface{}, res *InterconnectAttachment) *InterconnectAttachmentPartnerMetadata {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InterconnectAttachmentPartnerMetadata{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInterconnectAttachmentPartnerMetadata
	}
	r.PartnerName = dcl.FlattenString(m["partnerName"])
	r.InterconnectName = dcl.FlattenString(m["interconnectName"])
	r.PortalUrl = dcl.FlattenString(m["portalUrl"])

	return r
}

// flattenInterconnectAttachmentOperationalStatusEnumMap flattens the contents of InterconnectAttachmentOperationalStatusEnum from a JSON
// response object.
func flattenInterconnectAttachmentOperationalStatusEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentOperationalStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentOperationalStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentOperationalStatusEnum{}
	}

	items := make(map[string]InterconnectAttachmentOperationalStatusEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentOperationalStatusEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentOperationalStatusEnumSlice flattens the contents of InterconnectAttachmentOperationalStatusEnum from a JSON
// response object.
func flattenInterconnectAttachmentOperationalStatusEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentOperationalStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentOperationalStatusEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentOperationalStatusEnum{}
	}

	items := make([]InterconnectAttachmentOperationalStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentOperationalStatusEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentOperationalStatusEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentOperationalStatusEnum with the same value as that string.
func flattenInterconnectAttachmentOperationalStatusEnum(i interface{}) *InterconnectAttachmentOperationalStatusEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentOperationalStatusEnumRef(s)
}

// flattenInterconnectAttachmentTypeEnumMap flattens the contents of InterconnectAttachmentTypeEnum from a JSON
// response object.
func flattenInterconnectAttachmentTypeEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentTypeEnum{}
	}

	items := make(map[string]InterconnectAttachmentTypeEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentTypeEnumSlice flattens the contents of InterconnectAttachmentTypeEnum from a JSON
// response object.
func flattenInterconnectAttachmentTypeEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentTypeEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentTypeEnum{}
	}

	items := make([]InterconnectAttachmentTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentTypeEnum with the same value as that string.
func flattenInterconnectAttachmentTypeEnum(i interface{}) *InterconnectAttachmentTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentTypeEnumRef(s)
}

// flattenInterconnectAttachmentEdgeAvailabilityDomainEnumMap flattens the contents of InterconnectAttachmentEdgeAvailabilityDomainEnum from a JSON
// response object.
func flattenInterconnectAttachmentEdgeAvailabilityDomainEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentEdgeAvailabilityDomainEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentEdgeAvailabilityDomainEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentEdgeAvailabilityDomainEnum{}
	}

	items := make(map[string]InterconnectAttachmentEdgeAvailabilityDomainEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentEdgeAvailabilityDomainEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentEdgeAvailabilityDomainEnumSlice flattens the contents of InterconnectAttachmentEdgeAvailabilityDomainEnum from a JSON
// response object.
func flattenInterconnectAttachmentEdgeAvailabilityDomainEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentEdgeAvailabilityDomainEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentEdgeAvailabilityDomainEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentEdgeAvailabilityDomainEnum{}
	}

	items := make([]InterconnectAttachmentEdgeAvailabilityDomainEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentEdgeAvailabilityDomainEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentEdgeAvailabilityDomainEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentEdgeAvailabilityDomainEnum with the same value as that string.
func flattenInterconnectAttachmentEdgeAvailabilityDomainEnum(i interface{}) *InterconnectAttachmentEdgeAvailabilityDomainEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentEdgeAvailabilityDomainEnumRef(s)
}

// flattenInterconnectAttachmentBandwidthEnumMap flattens the contents of InterconnectAttachmentBandwidthEnum from a JSON
// response object.
func flattenInterconnectAttachmentBandwidthEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentBandwidthEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentBandwidthEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentBandwidthEnum{}
	}

	items := make(map[string]InterconnectAttachmentBandwidthEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentBandwidthEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentBandwidthEnumSlice flattens the contents of InterconnectAttachmentBandwidthEnum from a JSON
// response object.
func flattenInterconnectAttachmentBandwidthEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentBandwidthEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentBandwidthEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentBandwidthEnum{}
	}

	items := make([]InterconnectAttachmentBandwidthEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentBandwidthEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentBandwidthEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentBandwidthEnum with the same value as that string.
func flattenInterconnectAttachmentBandwidthEnum(i interface{}) *InterconnectAttachmentBandwidthEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentBandwidthEnumRef(s)
}

// flattenInterconnectAttachmentStateEnumMap flattens the contents of InterconnectAttachmentStateEnum from a JSON
// response object.
func flattenInterconnectAttachmentStateEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentStateEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentStateEnum{}
	}

	items := make(map[string]InterconnectAttachmentStateEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentStateEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentStateEnumSlice flattens the contents of InterconnectAttachmentStateEnum from a JSON
// response object.
func flattenInterconnectAttachmentStateEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentStateEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentStateEnum{}
	}

	items := make([]InterconnectAttachmentStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentStateEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentStateEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentStateEnum with the same value as that string.
func flattenInterconnectAttachmentStateEnum(i interface{}) *InterconnectAttachmentStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentStateEnumRef(s)
}

// flattenInterconnectAttachmentEncryptionEnumMap flattens the contents of InterconnectAttachmentEncryptionEnum from a JSON
// response object.
func flattenInterconnectAttachmentEncryptionEnumMap(c *Client, i interface{}, res *InterconnectAttachment) map[string]InterconnectAttachmentEncryptionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectAttachmentEncryptionEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectAttachmentEncryptionEnum{}
	}

	items := make(map[string]InterconnectAttachmentEncryptionEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectAttachmentEncryptionEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectAttachmentEncryptionEnumSlice flattens the contents of InterconnectAttachmentEncryptionEnum from a JSON
// response object.
func flattenInterconnectAttachmentEncryptionEnumSlice(c *Client, i interface{}, res *InterconnectAttachment) []InterconnectAttachmentEncryptionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectAttachmentEncryptionEnum{}
	}

	if len(a) == 0 {
		return []InterconnectAttachmentEncryptionEnum{}
	}

	items := make([]InterconnectAttachmentEncryptionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectAttachmentEncryptionEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectAttachmentEncryptionEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectAttachmentEncryptionEnum with the same value as that string.
func flattenInterconnectAttachmentEncryptionEnum(i interface{}) *InterconnectAttachmentEncryptionEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InterconnectAttachmentEncryptionEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *InterconnectAttachment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalInterconnectAttachment(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Project == nil && ncr.Project == nil {
			c.Config.Logger.Info("Both Project fields null - considering equal.")
		} else if nr.Project == nil || ncr.Project == nil {
			c.Config.Logger.Info("Only one Project field is null - considering unequal.")
			return false
		} else if *nr.Project != *ncr.Project {
			return false
		}
		if nr.Region == nil && ncr.Region == nil {
			c.Config.Logger.Info("Both Region fields null - considering equal.")
		} else if nr.Region == nil || ncr.Region == nil {
			c.Config.Logger.Info("Only one Region field is null - considering unequal.")
			return false
		} else if *nr.Region != *ncr.Region {
			return false
		}
		if nr.Name == nil && ncr.Name == nil {
			c.Config.Logger.Info("Both Name fields null - considering equal.")
		} else if nr.Name == nil || ncr.Name == nil {
			c.Config.Logger.Info("Only one Name field is null - considering unequal.")
			return false
		} else if *nr.Name != *ncr.Name {
			return false
		}
		return true
	}
}

type interconnectAttachmentDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         interconnectAttachmentApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToInterconnectAttachmentDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]interconnectAttachmentDiff, error) {
	opNamesToFieldDiffs := make(map[string][]*dcl.FieldDiff)
	// Map each operation name to the field diffs associated with it.
	for _, fd := range fds {
		for _, ro := range fd.ResultingOperation {
			if fieldDiffs, ok := opNamesToFieldDiffs[ro]; ok {
				fieldDiffs = append(fieldDiffs, fd)
				opNamesToFieldDiffs[ro] = fieldDiffs
			} else {
				config.Logger.Infof("%s required due to diff: %v", ro, fd)
				opNamesToFieldDiffs[ro] = []*dcl.FieldDiff{fd}
			}
		}
	}
	var diffs []interconnectAttachmentDiff
	// For each operation name, create a interconnectAttachmentDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := interconnectAttachmentDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToInterconnectAttachmentApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToInterconnectAttachmentApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (interconnectAttachmentApiOperation, error) {
	switch opName {

	case "updateInterconnectAttachmentPatchOperation":
		return &updateInterconnectAttachmentPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractInterconnectAttachmentFields(r *InterconnectAttachment) error {
	vPrivateInterconnectInfo := r.PrivateInterconnectInfo
	if vPrivateInterconnectInfo == nil {
		// note: explicitly not the empty object.
		vPrivateInterconnectInfo = &InterconnectAttachmentPrivateInterconnectInfo{}
	}
	if err := extractInterconnectAttachmentPrivateInterconnectInfoFields(r, vPrivateInterconnectInfo); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPrivateInterconnectInfo) {
		r.PrivateInterconnectInfo = vPrivateInterconnectInfo
	}
	vPartnerMetadata := r.PartnerMetadata
	if vPartnerMetadata == nil {
		// note: explicitly not the empty object.
		vPartnerMetadata = &InterconnectAttachmentPartnerMetadata{}
	}
	if err := extractInterconnectAttachmentPartnerMetadataFields(r, vPartnerMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPartnerMetadata) {
		r.PartnerMetadata = vPartnerMetadata
	}
	return nil
}
func extractInterconnectAttachmentPrivateInterconnectInfoFields(r *InterconnectAttachment, o *InterconnectAttachmentPrivateInterconnectInfo) error {
	return nil
}
func extractInterconnectAttachmentPartnerMetadataFields(r *InterconnectAttachment, o *InterconnectAttachmentPartnerMetadata) error {
	return nil
}

func postReadExtractInterconnectAttachmentFields(r *InterconnectAttachment) error {
	vPrivateInterconnectInfo := r.PrivateInterconnectInfo
	if vPrivateInterconnectInfo == nil {
		// note: explicitly not the empty object.
		vPrivateInterconnectInfo = &InterconnectAttachmentPrivateInterconnectInfo{}
	}
	if err := postReadExtractInterconnectAttachmentPrivateInterconnectInfoFields(r, vPrivateInterconnectInfo); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPrivateInterconnectInfo) {
		r.PrivateInterconnectInfo = vPrivateInterconnectInfo
	}
	vPartnerMetadata := r.PartnerMetadata
	if vPartnerMetadata == nil {
		// note: explicitly not the empty object.
		vPartnerMetadata = &InterconnectAttachmentPartnerMetadata{}
	}
	if err := postReadExtractInterconnectAttachmentPartnerMetadataFields(r, vPartnerMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPartnerMetadata) {
		r.PartnerMetadata = vPartnerMetadata
	}
	return nil
}
func postReadExtractInterconnectAttachmentPrivateInterconnectInfoFields(r *InterconnectAttachment, o *InterconnectAttachmentPrivateInterconnectInfo) error {
	return nil
}
func postReadExtractInterconnectAttachmentPartnerMetadataFields(r *InterconnectAttachment, o *InterconnectAttachmentPartnerMetadata) error {
	return nil
}
