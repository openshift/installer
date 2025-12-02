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
package compute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *ServiceAttachment) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "targetService"); err != nil {
		return err
	}
	if err := dcl.Required(r, "connectionPreference"); err != nil {
		return err
	}
	if err := dcl.Required(r, "natSubnets"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.PscServiceAttachmentId) {
		if err := r.PscServiceAttachmentId.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceAttachmentConnectedEndpoints) validate() error {
	return nil
}
func (r *ServiceAttachmentConsumerAcceptLists) validate() error {
	if err := dcl.Required(r, "projectIdOrNum"); err != nil {
		return err
	}
	return nil
}
func (r *ServiceAttachmentPscServiceAttachmentId) validate() error {
	return nil
}
func (r *ServiceAttachment) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *ServiceAttachment) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/serviceAttachments/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *ServiceAttachment) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/serviceAttachments", nr.basePath(), userBasePath, params), nil

}

func (r *ServiceAttachment) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/serviceAttachments", nr.basePath(), userBasePath, params), nil

}

func (r *ServiceAttachment) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/serviceAttachments/{{name}}", nr.basePath(), userBasePath, params), nil
}

// serviceAttachmentApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type serviceAttachmentApiOperation interface {
	do(context.Context, *ServiceAttachment, *Client) error
}

// newUpdateServiceAttachmentPatchRequest creates a request for an
// ServiceAttachment resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateServiceAttachmentPatchRequest(ctx context.Context, f *ServiceAttachment, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.ConnectionPreference; !dcl.IsEmptyValueIndirect(v) {
		req["connectionPreference"] = v
	}
	if v := f.NatSubnets; v != nil {
		req["natSubnets"] = v
	}
	if v, err := dcl.SelfLinkToNameArrayExpander(f.ConsumerRejectLists); err != nil {
		return nil, fmt.Errorf("error expanding ConsumerRejectLists into consumerRejectLists: %w", err)
	} else if v != nil {
		req["consumerRejectLists"] = v
	}
	if v, err := expandServiceAttachmentConsumerAcceptListsSlice(c, f.ConsumerAcceptLists, res); err != nil {
		return nil, fmt.Errorf("error expanding ConsumerAcceptLists into consumerAcceptLists: %w", err)
	} else if v != nil {
		req["consumerAcceptLists"] = v
	}
	b, err := c.getServiceAttachmentRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"fingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["fingerprint"] = rawFingerprint.(string)
	}
	req["name"] = fmt.Sprintf("%s", *f.Name)

	return req, nil
}

// marshalUpdateServiceAttachmentPatchRequest converts the update into
// the final JSON request body.
func marshalUpdateServiceAttachmentPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateServiceAttachmentPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateServiceAttachmentPatchOperation) do(ctx context.Context, r *ServiceAttachment, c *Client) error {
	_, err := c.GetServiceAttachment(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdateServiceAttachmentPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateServiceAttachmentPatchRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listServiceAttachmentRaw(ctx context.Context, r *ServiceAttachment, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ServiceAttachmentMaxPage {
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

type listServiceAttachmentOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listServiceAttachment(ctx context.Context, r *ServiceAttachment, pageToken string, pageSize int32) ([]*ServiceAttachment, string, error) {
	b, err := c.listServiceAttachmentRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listServiceAttachmentOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*ServiceAttachment
	for _, v := range m.Items {
		res, err := unmarshalMapServiceAttachment(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllServiceAttachment(ctx context.Context, f func(*ServiceAttachment) bool, resources []*ServiceAttachment) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteServiceAttachment(ctx, res)
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

type deleteServiceAttachmentOperation struct{}

func (op *deleteServiceAttachmentOperation) do(ctx context.Context, r *ServiceAttachment, c *Client) error {
	r, err := c.GetServiceAttachment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "ServiceAttachment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetServiceAttachment checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		return err
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetServiceAttachment(ctx, r)
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
type createServiceAttachmentOperation struct {
	response map[string]interface{}
}

func (op *createServiceAttachmentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createServiceAttachmentOperation) do(ctx context.Context, r *ServiceAttachment, c *Client) error {
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
	// wait for object to be created.
	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetServiceAttachment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getServiceAttachmentRaw(ctx context.Context, r *ServiceAttachment) ([]byte, error) {

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

func (c *Client) serviceAttachmentDiffsForRawDesired(ctx context.Context, rawDesired *ServiceAttachment, opts ...dcl.ApplyOption) (initial, desired *ServiceAttachment, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *ServiceAttachment
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*ServiceAttachment); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected ServiceAttachment, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetServiceAttachment(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a ServiceAttachment resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve ServiceAttachment resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that ServiceAttachment resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeServiceAttachmentDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for ServiceAttachment: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for ServiceAttachment: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractServiceAttachmentFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeServiceAttachmentInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for ServiceAttachment: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeServiceAttachmentDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for ServiceAttachment: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffServiceAttachment(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeServiceAttachmentInitialState(rawInitial, rawDesired *ServiceAttachment) (*ServiceAttachment, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeServiceAttachmentDesiredState(rawDesired, rawInitial *ServiceAttachment, opts ...dcl.ApplyOption) (*ServiceAttachment, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.PscServiceAttachmentId = canonicalizeServiceAttachmentPscServiceAttachmentId(rawDesired.PscServiceAttachmentId, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &ServiceAttachment{}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.TargetService) || (dcl.IsEmptyValueIndirect(rawDesired.TargetService) && dcl.IsEmptyValueIndirect(rawInitial.TargetService)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.TargetService = rawInitial.TargetService
	} else {
		canonicalDesired.TargetService = rawDesired.TargetService
	}
	if dcl.IsZeroValue(rawDesired.ConnectionPreference) || (dcl.IsEmptyValueIndirect(rawDesired.ConnectionPreference) && dcl.IsEmptyValueIndirect(rawInitial.ConnectionPreference)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.ConnectionPreference = rawInitial.ConnectionPreference
	} else {
		canonicalDesired.ConnectionPreference = rawDesired.ConnectionPreference
	}
	if dcl.StringArrayCanonicalize(rawDesired.NatSubnets, rawInitial.NatSubnets) {
		canonicalDesired.NatSubnets = rawInitial.NatSubnets
	} else {
		canonicalDesired.NatSubnets = rawDesired.NatSubnets
	}
	if dcl.BoolCanonicalize(rawDesired.EnableProxyProtocol, rawInitial.EnableProxyProtocol) {
		canonicalDesired.EnableProxyProtocol = rawInitial.EnableProxyProtocol
	} else {
		canonicalDesired.EnableProxyProtocol = rawDesired.EnableProxyProtocol
	}
	if dcl.StringArrayCanonicalize(rawDesired.ConsumerRejectLists, rawInitial.ConsumerRejectLists) {
		canonicalDesired.ConsumerRejectLists = rawInitial.ConsumerRejectLists
	} else {
		canonicalDesired.ConsumerRejectLists = rawDesired.ConsumerRejectLists
	}
	canonicalDesired.ConsumerAcceptLists = canonicalizeServiceAttachmentConsumerAcceptListsSlice(rawDesired.ConsumerAcceptLists, rawInitial.ConsumerAcceptLists, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	return canonicalDesired, nil
}

func canonicalizeServiceAttachmentNewState(c *Client, rawNew, rawDesired *ServiceAttachment) (*ServiceAttachment, error) {

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

	if dcl.IsEmptyValueIndirect(rawNew.Region) && dcl.IsEmptyValueIndirect(rawDesired.Region) {
		rawNew.Region = rawDesired.Region
	} else {
		if dcl.StringCanonicalize(rawDesired.Region, rawNew.Region) {
			rawNew.Region = rawDesired.Region
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.TargetService) && dcl.IsEmptyValueIndirect(rawDesired.TargetService) {
		rawNew.TargetService = rawDesired.TargetService
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ConnectionPreference) && dcl.IsEmptyValueIndirect(rawDesired.ConnectionPreference) {
		rawNew.ConnectionPreference = rawDesired.ConnectionPreference
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ConnectedEndpoints) && dcl.IsEmptyValueIndirect(rawDesired.ConnectedEndpoints) {
		rawNew.ConnectedEndpoints = rawDesired.ConnectedEndpoints
	} else {
		rawNew.ConnectedEndpoints = canonicalizeNewServiceAttachmentConnectedEndpointsSlice(c, rawDesired.ConnectedEndpoints, rawNew.ConnectedEndpoints)
	}

	if dcl.IsEmptyValueIndirect(rawNew.NatSubnets) && dcl.IsEmptyValueIndirect(rawDesired.NatSubnets) {
		rawNew.NatSubnets = rawDesired.NatSubnets
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.NatSubnets, rawNew.NatSubnets) {
			rawNew.NatSubnets = rawDesired.NatSubnets
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.EnableProxyProtocol) && dcl.IsEmptyValueIndirect(rawDesired.EnableProxyProtocol) {
		rawNew.EnableProxyProtocol = rawDesired.EnableProxyProtocol
	} else {
		if dcl.BoolCanonicalize(rawDesired.EnableProxyProtocol, rawNew.EnableProxyProtocol) {
			rawNew.EnableProxyProtocol = rawDesired.EnableProxyProtocol
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ConsumerRejectLists) && dcl.IsEmptyValueIndirect(rawDesired.ConsumerRejectLists) {
		rawNew.ConsumerRejectLists = rawDesired.ConsumerRejectLists
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.ConsumerRejectLists, rawNew.ConsumerRejectLists) {
			rawNew.ConsumerRejectLists = rawDesired.ConsumerRejectLists
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ConsumerAcceptLists) && dcl.IsEmptyValueIndirect(rawDesired.ConsumerAcceptLists) {
		rawNew.ConsumerAcceptLists = rawDesired.ConsumerAcceptLists
	} else {
		rawNew.ConsumerAcceptLists = canonicalizeNewServiceAttachmentConsumerAcceptListsSlice(c, rawDesired.ConsumerAcceptLists, rawNew.ConsumerAcceptLists)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PscServiceAttachmentId) && dcl.IsEmptyValueIndirect(rawDesired.PscServiceAttachmentId) {
		rawNew.PscServiceAttachmentId = rawDesired.PscServiceAttachmentId
	} else {
		rawNew.PscServiceAttachmentId = canonicalizeNewServiceAttachmentPscServiceAttachmentId(c, rawDesired.PscServiceAttachmentId, rawNew.PscServiceAttachmentId)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Fingerprint) && dcl.IsEmptyValueIndirect(rawDesired.Fingerprint) {
		rawNew.Fingerprint = rawDesired.Fingerprint
	} else {
		if dcl.StringCanonicalize(rawDesired.Fingerprint, rawNew.Fingerprint) {
			rawNew.Fingerprint = rawDesired.Fingerprint
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeServiceAttachmentConnectedEndpoints(des, initial *ServiceAttachmentConnectedEndpoints, opts ...dcl.ApplyOption) *ServiceAttachmentConnectedEndpoints {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceAttachmentConnectedEndpoints{}

	if dcl.IsZeroValue(des.Status) || (dcl.IsEmptyValueIndirect(des.Status) && dcl.IsEmptyValueIndirect(initial.Status)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Status = initial.Status
	} else {
		cDes.Status = des.Status
	}
	if dcl.IsZeroValue(des.PscConnectionId) || (dcl.IsEmptyValueIndirect(des.PscConnectionId) && dcl.IsEmptyValueIndirect(initial.PscConnectionId)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PscConnectionId = initial.PscConnectionId
	} else {
		cDes.PscConnectionId = des.PscConnectionId
	}
	if dcl.StringCanonicalize(des.Endpoint, initial.Endpoint) || dcl.IsZeroValue(des.Endpoint) {
		cDes.Endpoint = initial.Endpoint
	} else {
		cDes.Endpoint = des.Endpoint
	}

	return cDes
}

func canonicalizeServiceAttachmentConnectedEndpointsSlice(des, initial []ServiceAttachmentConnectedEndpoints, opts ...dcl.ApplyOption) []ServiceAttachmentConnectedEndpoints {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceAttachmentConnectedEndpoints, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceAttachmentConnectedEndpoints(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceAttachmentConnectedEndpoints, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceAttachmentConnectedEndpoints(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceAttachmentConnectedEndpoints(c *Client, des, nw *ServiceAttachmentConnectedEndpoints) *ServiceAttachmentConnectedEndpoints {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceAttachmentConnectedEndpoints while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Endpoint, nw.Endpoint) {
		nw.Endpoint = des.Endpoint
	}

	return nw
}

func canonicalizeNewServiceAttachmentConnectedEndpointsSet(c *Client, des, nw []ServiceAttachmentConnectedEndpoints) []ServiceAttachmentConnectedEndpoints {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceAttachmentConnectedEndpoints
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceAttachmentConnectedEndpointsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceAttachmentConnectedEndpoints(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceAttachmentConnectedEndpointsSlice(c *Client, des, nw []ServiceAttachmentConnectedEndpoints) []ServiceAttachmentConnectedEndpoints {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceAttachmentConnectedEndpoints
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceAttachmentConnectedEndpoints(c, &d, &n))
	}

	return items
}

func canonicalizeServiceAttachmentConsumerAcceptLists(des, initial *ServiceAttachmentConsumerAcceptLists, opts ...dcl.ApplyOption) *ServiceAttachmentConsumerAcceptLists {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceAttachmentConsumerAcceptLists{}

	if dcl.IsZeroValue(des.ProjectIdOrNum) || (dcl.IsEmptyValueIndirect(des.ProjectIdOrNum) && dcl.IsEmptyValueIndirect(initial.ProjectIdOrNum)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ProjectIdOrNum = initial.ProjectIdOrNum
	} else {
		cDes.ProjectIdOrNum = des.ProjectIdOrNum
	}
	if dcl.IsZeroValue(des.ConnectionLimit) || (dcl.IsEmptyValueIndirect(des.ConnectionLimit) && dcl.IsEmptyValueIndirect(initial.ConnectionLimit)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ConnectionLimit = initial.ConnectionLimit
	} else {
		cDes.ConnectionLimit = des.ConnectionLimit
	}

	return cDes
}

func canonicalizeServiceAttachmentConsumerAcceptListsSlice(des, initial []ServiceAttachmentConsumerAcceptLists, opts ...dcl.ApplyOption) []ServiceAttachmentConsumerAcceptLists {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceAttachmentConsumerAcceptLists, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceAttachmentConsumerAcceptLists(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceAttachmentConsumerAcceptLists, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceAttachmentConsumerAcceptLists(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceAttachmentConsumerAcceptLists(c *Client, des, nw *ServiceAttachmentConsumerAcceptLists) *ServiceAttachmentConsumerAcceptLists {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceAttachmentConsumerAcceptLists while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceAttachmentConsumerAcceptListsSet(c *Client, des, nw []ServiceAttachmentConsumerAcceptLists) []ServiceAttachmentConsumerAcceptLists {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceAttachmentConsumerAcceptLists
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceAttachmentConsumerAcceptListsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceAttachmentConsumerAcceptLists(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceAttachmentConsumerAcceptListsSlice(c *Client, des, nw []ServiceAttachmentConsumerAcceptLists) []ServiceAttachmentConsumerAcceptLists {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceAttachmentConsumerAcceptLists
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceAttachmentConsumerAcceptLists(c, &d, &n))
	}

	return items
}

func canonicalizeServiceAttachmentPscServiceAttachmentId(des, initial *ServiceAttachmentPscServiceAttachmentId, opts ...dcl.ApplyOption) *ServiceAttachmentPscServiceAttachmentId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceAttachmentPscServiceAttachmentId{}

	return cDes
}

func canonicalizeServiceAttachmentPscServiceAttachmentIdSlice(des, initial []ServiceAttachmentPscServiceAttachmentId, opts ...dcl.ApplyOption) []ServiceAttachmentPscServiceAttachmentId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceAttachmentPscServiceAttachmentId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceAttachmentPscServiceAttachmentId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceAttachmentPscServiceAttachmentId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceAttachmentPscServiceAttachmentId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceAttachmentPscServiceAttachmentId(c *Client, des, nw *ServiceAttachmentPscServiceAttachmentId) *ServiceAttachmentPscServiceAttachmentId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceAttachmentPscServiceAttachmentId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceAttachmentPscServiceAttachmentIdSet(c *Client, des, nw []ServiceAttachmentPscServiceAttachmentId) []ServiceAttachmentPscServiceAttachmentId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceAttachmentPscServiceAttachmentId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceAttachmentPscServiceAttachmentIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceAttachmentPscServiceAttachmentId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceAttachmentPscServiceAttachmentIdSlice(c *Client, des, nw []ServiceAttachmentPscServiceAttachmentId) []ServiceAttachmentPscServiceAttachmentId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceAttachmentPscServiceAttachmentId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceAttachmentPscServiceAttachmentId(c, &d, &n))
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
func diffServiceAttachment(c *Client, desired, actual *ServiceAttachment, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetService, actual.TargetService, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TargetService")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConnectionPreference, actual.ConnectionPreference, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("ConnectionPreference")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConnectedEndpoints, actual.ConnectedEndpoints, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareServiceAttachmentConnectedEndpointsNewStyle, EmptyObject: EmptyServiceAttachmentConnectedEndpoints, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConnectedEndpoints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NatSubnets, actual.NatSubnets, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("NatSubnets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableProxyProtocol, actual.EnableProxyProtocol, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableProxyProtocol")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConsumerRejectLists, actual.ConsumerRejectLists, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("ConsumerRejectLists")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConsumerAcceptLists, actual.ConsumerAcceptLists, dcl.DiffInfo{ObjectFunction: compareServiceAttachmentConsumerAcceptListsNewStyle, EmptyObject: EmptyServiceAttachmentConsumerAcceptLists, OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("ConsumerAcceptLists")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PscServiceAttachmentId, actual.PscServiceAttachmentId, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareServiceAttachmentPscServiceAttachmentIdNewStyle, EmptyObject: EmptyServiceAttachmentPscServiceAttachmentId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PscServiceAttachmentId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fingerprint, actual.Fingerprint, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fingerprint")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
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
func compareServiceAttachmentConnectedEndpointsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceAttachmentConnectedEndpoints)
	if !ok {
		desiredNotPointer, ok := d.(ServiceAttachmentConnectedEndpoints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentConnectedEndpoints or *ServiceAttachmentConnectedEndpoints", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceAttachmentConnectedEndpoints)
	if !ok {
		actualNotPointer, ok := a.(ServiceAttachmentConnectedEndpoints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentConnectedEndpoints", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PscConnectionId, actual.PscConnectionId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("PscConnectionId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Endpoint, actual.Endpoint, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("Endpoint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceAttachmentConsumerAcceptListsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceAttachmentConsumerAcceptLists)
	if !ok {
		desiredNotPointer, ok := d.(ServiceAttachmentConsumerAcceptLists)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentConsumerAcceptLists or *ServiceAttachmentConsumerAcceptLists", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceAttachmentConsumerAcceptLists)
	if !ok {
		actualNotPointer, ok := a.(ServiceAttachmentConsumerAcceptLists)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentConsumerAcceptLists", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ProjectIdOrNum, actual.ProjectIdOrNum, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("ProjectIdOrNum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConnectionLimit, actual.ConnectionLimit, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceAttachmentPatchOperation")}, fn.AddNest("ConnectionLimit")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceAttachmentPscServiceAttachmentIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceAttachmentPscServiceAttachmentId)
	if !ok {
		desiredNotPointer, ok := d.(ServiceAttachmentPscServiceAttachmentId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentPscServiceAttachmentId or *ServiceAttachmentPscServiceAttachmentId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceAttachmentPscServiceAttachmentId)
	if !ok {
		actualNotPointer, ok := a.(ServiceAttachmentPscServiceAttachmentId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceAttachmentPscServiceAttachmentId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.High, actual.High, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("High")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Low, actual.Low, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Low")); len(ds) != 0 || err != nil {
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
func (r *ServiceAttachment) urlNormalized() *ServiceAttachment {
	normalized := dcl.Copy(*r).(ServiceAttachment)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.TargetService = dcl.SelfLinkToName(r.TargetService)
	normalized.Fingerprint = dcl.SelfLinkToName(r.Fingerprint)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *ServiceAttachment) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{location}}/serviceAttachments/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the ServiceAttachment resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *ServiceAttachment) marshal(c *Client) ([]byte, error) {
	m, err := expandServiceAttachment(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ServiceAttachment: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalServiceAttachment decodes JSON responses into the ServiceAttachment resource schema.
func unmarshalServiceAttachment(b []byte, c *Client, res *ServiceAttachment) (*ServiceAttachment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapServiceAttachment(m, c, res)
}

func unmarshalMapServiceAttachment(m map[string]interface{}, c *Client, res *ServiceAttachment) (*ServiceAttachment, error) {

	flattened := flattenServiceAttachment(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandServiceAttachment expands ServiceAttachment into a JSON request object.
func expandServiceAttachment(c *Client, f *ServiceAttachment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.TargetService; dcl.ValueShouldBeSent(v) {
		m["targetService"] = v
	}
	if v := f.ConnectionPreference; dcl.ValueShouldBeSent(v) {
		m["connectionPreference"] = v
	}
	if v := f.NatSubnets; v != nil {
		m["natSubnets"] = v
	}
	if v := f.EnableProxyProtocol; dcl.ValueShouldBeSent(v) {
		m["enableProxyProtocol"] = v
	}
	if v, err := dcl.SelfLinkToNameArrayExpander(f.ConsumerRejectLists); err != nil {
		return nil, fmt.Errorf("error expanding ConsumerRejectLists into consumerRejectLists: %w", err)
	} else if v != nil {
		m["consumerRejectLists"] = v
	}
	if v, err := expandServiceAttachmentConsumerAcceptListsSlice(c, f.ConsumerAcceptLists, res); err != nil {
		return nil, fmt.Errorf("error expanding ConsumerAcceptLists into consumerAcceptLists: %w", err)
	} else if v != nil {
		m["consumerAcceptLists"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenServiceAttachment flattens ServiceAttachment from a JSON request object into the
// ServiceAttachment type.
func flattenServiceAttachment(c *Client, i interface{}, res *ServiceAttachment) *ServiceAttachment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &ServiceAttachment{}
	resultRes.Id = dcl.FlattenInteger(m["id"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.TargetService = dcl.FlattenString(m["targetService"])
	resultRes.ConnectionPreference = flattenServiceAttachmentConnectionPreferenceEnum(m["connectionPreference"])
	resultRes.ConnectedEndpoints = flattenServiceAttachmentConnectedEndpointsSlice(c, m["connectedEndpoints"], res)
	resultRes.NatSubnets = dcl.FlattenStringSlice(m["natSubnets"])
	resultRes.EnableProxyProtocol = dcl.FlattenBool(m["enableProxyProtocol"])
	resultRes.ConsumerRejectLists = dcl.FlattenStringSlice(m["consumerRejectLists"])
	resultRes.ConsumerAcceptLists = flattenServiceAttachmentConsumerAcceptListsSlice(c, m["consumerAcceptLists"], res)
	resultRes.PscServiceAttachmentId = flattenServiceAttachmentPscServiceAttachmentId(c, m["pscServiceAttachmentId"], res)
	resultRes.Fingerprint = dcl.FlattenString(m["fingerprint"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandServiceAttachmentConnectedEndpointsMap expands the contents of ServiceAttachmentConnectedEndpoints into a JSON
// request object.
func expandServiceAttachmentConnectedEndpointsMap(c *Client, f map[string]ServiceAttachmentConnectedEndpoints, res *ServiceAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceAttachmentConnectedEndpoints(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceAttachmentConnectedEndpointsSlice expands the contents of ServiceAttachmentConnectedEndpoints into a JSON
// request object.
func expandServiceAttachmentConnectedEndpointsSlice(c *Client, f []ServiceAttachmentConnectedEndpoints, res *ServiceAttachment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceAttachmentConnectedEndpoints(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceAttachmentConnectedEndpointsMap flattens the contents of ServiceAttachmentConnectedEndpoints from a JSON
// response object.
func flattenServiceAttachmentConnectedEndpointsMap(c *Client, i interface{}, res *ServiceAttachment) map[string]ServiceAttachmentConnectedEndpoints {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceAttachmentConnectedEndpoints{}
	}

	if len(a) == 0 {
		return map[string]ServiceAttachmentConnectedEndpoints{}
	}

	items := make(map[string]ServiceAttachmentConnectedEndpoints)
	for k, item := range a {
		items[k] = *flattenServiceAttachmentConnectedEndpoints(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceAttachmentConnectedEndpointsSlice flattens the contents of ServiceAttachmentConnectedEndpoints from a JSON
// response object.
func flattenServiceAttachmentConnectedEndpointsSlice(c *Client, i interface{}, res *ServiceAttachment) []ServiceAttachmentConnectedEndpoints {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceAttachmentConnectedEndpoints{}
	}

	if len(a) == 0 {
		return []ServiceAttachmentConnectedEndpoints{}
	}

	items := make([]ServiceAttachmentConnectedEndpoints, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceAttachmentConnectedEndpoints(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceAttachmentConnectedEndpoints expands an instance of ServiceAttachmentConnectedEndpoints into a JSON
// request object.
func expandServiceAttachmentConnectedEndpoints(c *Client, f *ServiceAttachmentConnectedEndpoints, res *ServiceAttachment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Status; !dcl.IsEmptyValueIndirect(v) {
		m["status"] = v
	}
	if v := f.PscConnectionId; !dcl.IsEmptyValueIndirect(v) {
		m["pscConnectionId"] = v
	}
	if v := f.Endpoint; !dcl.IsEmptyValueIndirect(v) {
		m["endpoint"] = v
	}

	return m, nil
}

// flattenServiceAttachmentConnectedEndpoints flattens an instance of ServiceAttachmentConnectedEndpoints from a JSON
// response object.
func flattenServiceAttachmentConnectedEndpoints(c *Client, i interface{}, res *ServiceAttachment) *ServiceAttachmentConnectedEndpoints {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceAttachmentConnectedEndpoints{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceAttachmentConnectedEndpoints
	}
	r.Status = flattenServiceAttachmentConnectedEndpointsStatusEnum(m["status"])
	r.PscConnectionId = dcl.FlattenInteger(m["pscConnectionId"])
	r.Endpoint = dcl.FlattenString(m["endpoint"])

	return r
}

// expandServiceAttachmentConsumerAcceptListsMap expands the contents of ServiceAttachmentConsumerAcceptLists into a JSON
// request object.
func expandServiceAttachmentConsumerAcceptListsMap(c *Client, f map[string]ServiceAttachmentConsumerAcceptLists, res *ServiceAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceAttachmentConsumerAcceptLists(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceAttachmentConsumerAcceptListsSlice expands the contents of ServiceAttachmentConsumerAcceptLists into a JSON
// request object.
func expandServiceAttachmentConsumerAcceptListsSlice(c *Client, f []ServiceAttachmentConsumerAcceptLists, res *ServiceAttachment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceAttachmentConsumerAcceptLists(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceAttachmentConsumerAcceptListsMap flattens the contents of ServiceAttachmentConsumerAcceptLists from a JSON
// response object.
func flattenServiceAttachmentConsumerAcceptListsMap(c *Client, i interface{}, res *ServiceAttachment) map[string]ServiceAttachmentConsumerAcceptLists {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceAttachmentConsumerAcceptLists{}
	}

	if len(a) == 0 {
		return map[string]ServiceAttachmentConsumerAcceptLists{}
	}

	items := make(map[string]ServiceAttachmentConsumerAcceptLists)
	for k, item := range a {
		items[k] = *flattenServiceAttachmentConsumerAcceptLists(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceAttachmentConsumerAcceptListsSlice flattens the contents of ServiceAttachmentConsumerAcceptLists from a JSON
// response object.
func flattenServiceAttachmentConsumerAcceptListsSlice(c *Client, i interface{}, res *ServiceAttachment) []ServiceAttachmentConsumerAcceptLists {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceAttachmentConsumerAcceptLists{}
	}

	if len(a) == 0 {
		return []ServiceAttachmentConsumerAcceptLists{}
	}

	items := make([]ServiceAttachmentConsumerAcceptLists, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceAttachmentConsumerAcceptLists(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceAttachmentConsumerAcceptLists expands an instance of ServiceAttachmentConsumerAcceptLists into a JSON
// request object.
func expandServiceAttachmentConsumerAcceptLists(c *Client, f *ServiceAttachmentConsumerAcceptLists, res *ServiceAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := dcl.SelfLinkToNameExpander(f.ProjectIdOrNum); err != nil {
		return nil, fmt.Errorf("error expanding ProjectIdOrNum into projectIdOrNum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["projectIdOrNum"] = v
	}
	if v := f.ConnectionLimit; !dcl.IsEmptyValueIndirect(v) {
		m["connectionLimit"] = v
	}

	return m, nil
}

// flattenServiceAttachmentConsumerAcceptLists flattens an instance of ServiceAttachmentConsumerAcceptLists from a JSON
// response object.
func flattenServiceAttachmentConsumerAcceptLists(c *Client, i interface{}, res *ServiceAttachment) *ServiceAttachmentConsumerAcceptLists {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceAttachmentConsumerAcceptLists{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceAttachmentConsumerAcceptLists
	}
	r.ProjectIdOrNum = dcl.FlattenString(m["projectIdOrNum"])
	r.ConnectionLimit = dcl.FlattenInteger(m["connectionLimit"])

	return r
}

// expandServiceAttachmentPscServiceAttachmentIdMap expands the contents of ServiceAttachmentPscServiceAttachmentId into a JSON
// request object.
func expandServiceAttachmentPscServiceAttachmentIdMap(c *Client, f map[string]ServiceAttachmentPscServiceAttachmentId, res *ServiceAttachment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceAttachmentPscServiceAttachmentId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceAttachmentPscServiceAttachmentIdSlice expands the contents of ServiceAttachmentPscServiceAttachmentId into a JSON
// request object.
func expandServiceAttachmentPscServiceAttachmentIdSlice(c *Client, f []ServiceAttachmentPscServiceAttachmentId, res *ServiceAttachment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceAttachmentPscServiceAttachmentId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceAttachmentPscServiceAttachmentIdMap flattens the contents of ServiceAttachmentPscServiceAttachmentId from a JSON
// response object.
func flattenServiceAttachmentPscServiceAttachmentIdMap(c *Client, i interface{}, res *ServiceAttachment) map[string]ServiceAttachmentPscServiceAttachmentId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceAttachmentPscServiceAttachmentId{}
	}

	if len(a) == 0 {
		return map[string]ServiceAttachmentPscServiceAttachmentId{}
	}

	items := make(map[string]ServiceAttachmentPscServiceAttachmentId)
	for k, item := range a {
		items[k] = *flattenServiceAttachmentPscServiceAttachmentId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceAttachmentPscServiceAttachmentIdSlice flattens the contents of ServiceAttachmentPscServiceAttachmentId from a JSON
// response object.
func flattenServiceAttachmentPscServiceAttachmentIdSlice(c *Client, i interface{}, res *ServiceAttachment) []ServiceAttachmentPscServiceAttachmentId {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceAttachmentPscServiceAttachmentId{}
	}

	if len(a) == 0 {
		return []ServiceAttachmentPscServiceAttachmentId{}
	}

	items := make([]ServiceAttachmentPscServiceAttachmentId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceAttachmentPscServiceAttachmentId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceAttachmentPscServiceAttachmentId expands an instance of ServiceAttachmentPscServiceAttachmentId into a JSON
// request object.
func expandServiceAttachmentPscServiceAttachmentId(c *Client, f *ServiceAttachmentPscServiceAttachmentId, res *ServiceAttachment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceAttachmentPscServiceAttachmentId flattens an instance of ServiceAttachmentPscServiceAttachmentId from a JSON
// response object.
func flattenServiceAttachmentPscServiceAttachmentId(c *Client, i interface{}, res *ServiceAttachment) *ServiceAttachmentPscServiceAttachmentId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceAttachmentPscServiceAttachmentId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceAttachmentPscServiceAttachmentId
	}
	r.High = dcl.FlattenInteger(m["high"])
	r.Low = dcl.FlattenInteger(m["low"])

	return r
}

// flattenServiceAttachmentConnectionPreferenceEnumMap flattens the contents of ServiceAttachmentConnectionPreferenceEnum from a JSON
// response object.
func flattenServiceAttachmentConnectionPreferenceEnumMap(c *Client, i interface{}, res *ServiceAttachment) map[string]ServiceAttachmentConnectionPreferenceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceAttachmentConnectionPreferenceEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceAttachmentConnectionPreferenceEnum{}
	}

	items := make(map[string]ServiceAttachmentConnectionPreferenceEnum)
	for k, item := range a {
		items[k] = *flattenServiceAttachmentConnectionPreferenceEnum(item.(interface{}))
	}

	return items
}

// flattenServiceAttachmentConnectionPreferenceEnumSlice flattens the contents of ServiceAttachmentConnectionPreferenceEnum from a JSON
// response object.
func flattenServiceAttachmentConnectionPreferenceEnumSlice(c *Client, i interface{}, res *ServiceAttachment) []ServiceAttachmentConnectionPreferenceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceAttachmentConnectionPreferenceEnum{}
	}

	if len(a) == 0 {
		return []ServiceAttachmentConnectionPreferenceEnum{}
	}

	items := make([]ServiceAttachmentConnectionPreferenceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceAttachmentConnectionPreferenceEnum(item.(interface{})))
	}

	return items
}

// flattenServiceAttachmentConnectionPreferenceEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceAttachmentConnectionPreferenceEnum with the same value as that string.
func flattenServiceAttachmentConnectionPreferenceEnum(i interface{}) *ServiceAttachmentConnectionPreferenceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceAttachmentConnectionPreferenceEnumRef(s)
}

// flattenServiceAttachmentConnectedEndpointsStatusEnumMap flattens the contents of ServiceAttachmentConnectedEndpointsStatusEnum from a JSON
// response object.
func flattenServiceAttachmentConnectedEndpointsStatusEnumMap(c *Client, i interface{}, res *ServiceAttachment) map[string]ServiceAttachmentConnectedEndpointsStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceAttachmentConnectedEndpointsStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceAttachmentConnectedEndpointsStatusEnum{}
	}

	items := make(map[string]ServiceAttachmentConnectedEndpointsStatusEnum)
	for k, item := range a {
		items[k] = *flattenServiceAttachmentConnectedEndpointsStatusEnum(item.(interface{}))
	}

	return items
}

// flattenServiceAttachmentConnectedEndpointsStatusEnumSlice flattens the contents of ServiceAttachmentConnectedEndpointsStatusEnum from a JSON
// response object.
func flattenServiceAttachmentConnectedEndpointsStatusEnumSlice(c *Client, i interface{}, res *ServiceAttachment) []ServiceAttachmentConnectedEndpointsStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceAttachmentConnectedEndpointsStatusEnum{}
	}

	if len(a) == 0 {
		return []ServiceAttachmentConnectedEndpointsStatusEnum{}
	}

	items := make([]ServiceAttachmentConnectedEndpointsStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceAttachmentConnectedEndpointsStatusEnum(item.(interface{})))
	}

	return items
}

// flattenServiceAttachmentConnectedEndpointsStatusEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceAttachmentConnectedEndpointsStatusEnum with the same value as that string.
func flattenServiceAttachmentConnectedEndpointsStatusEnum(i interface{}) *ServiceAttachmentConnectedEndpointsStatusEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceAttachmentConnectedEndpointsStatusEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *ServiceAttachment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalServiceAttachment(b, c, r)
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
		if nr.Location == nil && ncr.Location == nil {
			c.Config.Logger.Info("Both Location fields null - considering equal.")
		} else if nr.Location == nil || ncr.Location == nil {
			c.Config.Logger.Info("Only one Location field is null - considering unequal.")
			return false
		} else if *nr.Location != *ncr.Location {
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

type serviceAttachmentDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         serviceAttachmentApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToServiceAttachmentDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]serviceAttachmentDiff, error) {
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
	var diffs []serviceAttachmentDiff
	// For each operation name, create a serviceAttachmentDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := serviceAttachmentDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToServiceAttachmentApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToServiceAttachmentApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (serviceAttachmentApiOperation, error) {
	switch opName {

	case "updateServiceAttachmentPatchOperation":
		return &updateServiceAttachmentPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractServiceAttachmentFields(r *ServiceAttachment) error {
	vPscServiceAttachmentId := r.PscServiceAttachmentId
	if vPscServiceAttachmentId == nil {
		// note: explicitly not the empty object.
		vPscServiceAttachmentId = &ServiceAttachmentPscServiceAttachmentId{}
	}
	if err := extractServiceAttachmentPscServiceAttachmentIdFields(r, vPscServiceAttachmentId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPscServiceAttachmentId) {
		r.PscServiceAttachmentId = vPscServiceAttachmentId
	}
	return nil
}
func extractServiceAttachmentConnectedEndpointsFields(r *ServiceAttachment, o *ServiceAttachmentConnectedEndpoints) error {
	return nil
}
func extractServiceAttachmentConsumerAcceptListsFields(r *ServiceAttachment, o *ServiceAttachmentConsumerAcceptLists) error {
	return nil
}
func extractServiceAttachmentPscServiceAttachmentIdFields(r *ServiceAttachment, o *ServiceAttachmentPscServiceAttachmentId) error {
	return nil
}

func postReadExtractServiceAttachmentFields(r *ServiceAttachment) error {
	vPscServiceAttachmentId := r.PscServiceAttachmentId
	if vPscServiceAttachmentId == nil {
		// note: explicitly not the empty object.
		vPscServiceAttachmentId = &ServiceAttachmentPscServiceAttachmentId{}
	}
	if err := postReadExtractServiceAttachmentPscServiceAttachmentIdFields(r, vPscServiceAttachmentId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPscServiceAttachmentId) {
		r.PscServiceAttachmentId = vPscServiceAttachmentId
	}
	return nil
}
func postReadExtractServiceAttachmentConnectedEndpointsFields(r *ServiceAttachment, o *ServiceAttachmentConnectedEndpoints) error {
	return nil
}
func postReadExtractServiceAttachmentConsumerAcceptListsFields(r *ServiceAttachment, o *ServiceAttachmentConsumerAcceptLists) error {
	return nil
}
func postReadExtractServiceAttachmentPscServiceAttachmentIdFields(r *ServiceAttachment, o *ServiceAttachmentPscServiceAttachmentId) error {
	return nil
}
