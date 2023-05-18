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
package dataplex

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

func (r *Zone) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "type"); err != nil {
		return err
	}
	if err := dcl.Required(r, "discoverySpec"); err != nil {
		return err
	}
	if err := dcl.Required(r, "resourceSpec"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Lake, "Lake"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.DiscoverySpec) {
		if err := r.DiscoverySpec.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ResourceSpec) {
		if err := r.ResourceSpec.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AssetStatus) {
		if err := r.AssetStatus.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ZoneDiscoverySpec) validate() error {
	if err := dcl.Required(r, "enabled"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.CsvOptions) {
		if err := r.CsvOptions.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.JsonOptions) {
		if err := r.JsonOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ZoneDiscoverySpecCsvOptions) validate() error {
	return nil
}
func (r *ZoneDiscoverySpecJsonOptions) validate() error {
	return nil
}
func (r *ZoneResourceSpec) validate() error {
	if err := dcl.Required(r, "locationType"); err != nil {
		return err
	}
	return nil
}
func (r *ZoneAssetStatus) validate() error {
	return nil
}
func (r *Zone) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://dataplex.googleapis.com/v1/", params)
}

func (r *Zone) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"lake":     dcl.ValueOrEmptyString(nr.Lake),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Zone) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"lake":     dcl.ValueOrEmptyString(nr.Lake),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones", nr.basePath(), userBasePath, params), nil

}

func (r *Zone) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"lake":     dcl.ValueOrEmptyString(nr.Lake),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones?zoneId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Zone) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"lake":     dcl.ValueOrEmptyString(nr.Lake),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Zone) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *Zone) SetPolicyVerb() string {
	return ""
}

func (r *Zone) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *Zone) IAMPolicyVersion() int {
	return 3
}

// zoneApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type zoneApiOperation interface {
	do(context.Context, *Zone, *Client) error
}

// newUpdateZoneUpdateZoneRequest creates a request for an
// Zone resource's UpdateZone update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateZoneUpdateZoneRequest(ctx context.Context, f *Zone, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.DeriveField("projects/%s/locations/%s/lakes/%s/zones/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Lake), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["name"] = v
	}
	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandZoneDiscoverySpec(c, f.DiscoverySpec, res); err != nil {
		return nil, fmt.Errorf("error expanding DiscoverySpec into discoverySpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["discoverySpec"] = v
	}
	if v, err := expandZoneAssetStatus(c, f.AssetStatus, res); err != nil {
		return nil, fmt.Errorf("error expanding AssetStatus into assetStatus: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["assetStatus"] = v
	}
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/lakes/%s/zones/%s", *f.Project, *f.Location, *f.Lake, *f.Name)

	return req, nil
}

// marshalUpdateZoneUpdateZoneRequest converts the update into
// the final JSON request body.
func marshalUpdateZoneUpdateZoneRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateZoneUpdateZoneOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateZoneUpdateZoneOperation) do(ctx context.Context, r *Zone, c *Client) error {
	_, err := c.GetZone(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateZone")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateZoneUpdateZoneRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateZoneUpdateZoneRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listZoneRaw(ctx context.Context, r *Zone, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ZoneMaxPage {
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

type listZoneOperation struct {
	Zones []map[string]interface{} `json:"zones"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listZone(ctx context.Context, r *Zone, pageToken string, pageSize int32) ([]*Zone, string, error) {
	b, err := c.listZoneRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listZoneOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Zone
	for _, v := range m.Zones {
		res, err := unmarshalMapZone(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.Lake = r.Lake
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllZone(ctx context.Context, f func(*Zone) bool, resources []*Zone) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteZone(ctx, res)
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

type deleteZoneOperation struct{}

func (op *deleteZoneOperation) do(ctx context.Context, r *Zone, c *Client) error {
	r, err := c.GetZone(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Zone not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetZone checking for existence. error: %v", err)
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
	var o operations.StandardGCPOperation
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
		_, err := c.GetZone(ctx, r)
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
type createZoneOperation struct {
	response map[string]interface{}
}

func (op *createZoneOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createZoneOperation) do(ctx context.Context, r *Zone, c *Client) error {
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
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetZone(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getZoneRaw(ctx context.Context, r *Zone) ([]byte, error) {

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

func (c *Client) zoneDiffsForRawDesired(ctx context.Context, rawDesired *Zone, opts ...dcl.ApplyOption) (initial, desired *Zone, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Zone
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Zone); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Zone, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetZone(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Zone resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Zone resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Zone resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeZoneDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Zone: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Zone: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractZoneFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeZoneInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Zone: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeZoneDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Zone: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffZone(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeZoneInitialState(rawInitial, rawDesired *Zone) (*Zone, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeZoneDesiredState(rawDesired, rawInitial *Zone, opts ...dcl.ApplyOption) (*Zone, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.DiscoverySpec = canonicalizeZoneDiscoverySpec(rawDesired.DiscoverySpec, nil, opts...)
		rawDesired.ResourceSpec = canonicalizeZoneResourceSpec(rawDesired.ResourceSpec, nil, opts...)
		rawDesired.AssetStatus = canonicalizeZoneAssetStatus(rawDesired.AssetStatus, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Zone{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.Type) || (dcl.IsEmptyValueIndirect(rawDesired.Type) && dcl.IsEmptyValueIndirect(rawInitial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Type = rawInitial.Type
	} else {
		canonicalDesired.Type = rawDesired.Type
	}
	canonicalDesired.DiscoverySpec = canonicalizeZoneDiscoverySpec(rawDesired.DiscoverySpec, rawInitial.DiscoverySpec, opts...)
	canonicalDesired.ResourceSpec = canonicalizeZoneResourceSpec(rawDesired.ResourceSpec, rawInitial.ResourceSpec, opts...)
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
	if dcl.NameToSelfLink(rawDesired.Lake, rawInitial.Lake) {
		canonicalDesired.Lake = rawInitial.Lake
	} else {
		canonicalDesired.Lake = rawDesired.Lake
	}
	return canonicalDesired, nil
}

func canonicalizeZoneNewState(c *Client, rawNew, rawDesired *Zone) (*Zone, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Type) && dcl.IsEmptyValueIndirect(rawDesired.Type) {
		rawNew.Type = rawDesired.Type
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DiscoverySpec) && dcl.IsEmptyValueIndirect(rawDesired.DiscoverySpec) {
		rawNew.DiscoverySpec = rawDesired.DiscoverySpec
	} else {
		rawNew.DiscoverySpec = canonicalizeNewZoneDiscoverySpec(c, rawDesired.DiscoverySpec, rawNew.DiscoverySpec)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ResourceSpec) && dcl.IsEmptyValueIndirect(rawDesired.ResourceSpec) {
		rawNew.ResourceSpec = rawDesired.ResourceSpec
	} else {
		rawNew.ResourceSpec = canonicalizeNewZoneResourceSpec(c, rawDesired.ResourceSpec, rawNew.ResourceSpec)
	}

	if dcl.IsEmptyValueIndirect(rawNew.AssetStatus) && dcl.IsEmptyValueIndirect(rawDesired.AssetStatus) {
		rawNew.AssetStatus = rawDesired.AssetStatus
	} else {
		rawNew.AssetStatus = canonicalizeNewZoneAssetStatus(c, rawDesired.AssetStatus, rawNew.AssetStatus)
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	rawNew.Lake = rawDesired.Lake

	return rawNew, nil
}

func canonicalizeZoneDiscoverySpec(des, initial *ZoneDiscoverySpec, opts ...dcl.ApplyOption) *ZoneDiscoverySpec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ZoneDiscoverySpec{}

	if dcl.BoolCanonicalize(des.Enabled, initial.Enabled) || dcl.IsZeroValue(des.Enabled) {
		cDes.Enabled = initial.Enabled
	} else {
		cDes.Enabled = des.Enabled
	}
	if dcl.StringArrayCanonicalize(des.IncludePatterns, initial.IncludePatterns) {
		cDes.IncludePatterns = initial.IncludePatterns
	} else {
		cDes.IncludePatterns = des.IncludePatterns
	}
	if dcl.StringArrayCanonicalize(des.ExcludePatterns, initial.ExcludePatterns) {
		cDes.ExcludePatterns = initial.ExcludePatterns
	} else {
		cDes.ExcludePatterns = des.ExcludePatterns
	}
	cDes.CsvOptions = canonicalizeZoneDiscoverySpecCsvOptions(des.CsvOptions, initial.CsvOptions, opts...)
	cDes.JsonOptions = canonicalizeZoneDiscoverySpecJsonOptions(des.JsonOptions, initial.JsonOptions, opts...)
	if dcl.StringCanonicalize(des.Schedule, initial.Schedule) || dcl.IsZeroValue(des.Schedule) {
		cDes.Schedule = initial.Schedule
	} else {
		cDes.Schedule = des.Schedule
	}

	return cDes
}

func canonicalizeZoneDiscoverySpecSlice(des, initial []ZoneDiscoverySpec, opts ...dcl.ApplyOption) []ZoneDiscoverySpec {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ZoneDiscoverySpec, 0, len(des))
		for _, d := range des {
			cd := canonicalizeZoneDiscoverySpec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ZoneDiscoverySpec, 0, len(des))
	for i, d := range des {
		cd := canonicalizeZoneDiscoverySpec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewZoneDiscoverySpec(c *Client, des, nw *ZoneDiscoverySpec) *ZoneDiscoverySpec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ZoneDiscoverySpec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.Enabled, nw.Enabled) {
		nw.Enabled = des.Enabled
	}
	if dcl.StringArrayCanonicalize(des.IncludePatterns, nw.IncludePatterns) {
		nw.IncludePatterns = des.IncludePatterns
	}
	if dcl.StringArrayCanonicalize(des.ExcludePatterns, nw.ExcludePatterns) {
		nw.ExcludePatterns = des.ExcludePatterns
	}
	nw.CsvOptions = canonicalizeNewZoneDiscoverySpecCsvOptions(c, des.CsvOptions, nw.CsvOptions)
	nw.JsonOptions = canonicalizeNewZoneDiscoverySpecJsonOptions(c, des.JsonOptions, nw.JsonOptions)
	if dcl.StringCanonicalize(des.Schedule, nw.Schedule) {
		nw.Schedule = des.Schedule
	}

	return nw
}

func canonicalizeNewZoneDiscoverySpecSet(c *Client, des, nw []ZoneDiscoverySpec) []ZoneDiscoverySpec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ZoneDiscoverySpec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareZoneDiscoverySpecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewZoneDiscoverySpec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewZoneDiscoverySpecSlice(c *Client, des, nw []ZoneDiscoverySpec) []ZoneDiscoverySpec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ZoneDiscoverySpec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewZoneDiscoverySpec(c, &d, &n))
	}

	return items
}

func canonicalizeZoneDiscoverySpecCsvOptions(des, initial *ZoneDiscoverySpecCsvOptions, opts ...dcl.ApplyOption) *ZoneDiscoverySpecCsvOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ZoneDiscoverySpecCsvOptions{}

	if dcl.IsZeroValue(des.HeaderRows) || (dcl.IsEmptyValueIndirect(des.HeaderRows) && dcl.IsEmptyValueIndirect(initial.HeaderRows)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.HeaderRows = initial.HeaderRows
	} else {
		cDes.HeaderRows = des.HeaderRows
	}
	if dcl.StringCanonicalize(des.Delimiter, initial.Delimiter) || dcl.IsZeroValue(des.Delimiter) {
		cDes.Delimiter = initial.Delimiter
	} else {
		cDes.Delimiter = des.Delimiter
	}
	if dcl.StringCanonicalize(des.Encoding, initial.Encoding) || dcl.IsZeroValue(des.Encoding) {
		cDes.Encoding = initial.Encoding
	} else {
		cDes.Encoding = des.Encoding
	}
	if dcl.BoolCanonicalize(des.DisableTypeInference, initial.DisableTypeInference) || dcl.IsZeroValue(des.DisableTypeInference) {
		cDes.DisableTypeInference = initial.DisableTypeInference
	} else {
		cDes.DisableTypeInference = des.DisableTypeInference
	}

	return cDes
}

func canonicalizeZoneDiscoverySpecCsvOptionsSlice(des, initial []ZoneDiscoverySpecCsvOptions, opts ...dcl.ApplyOption) []ZoneDiscoverySpecCsvOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ZoneDiscoverySpecCsvOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeZoneDiscoverySpecCsvOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ZoneDiscoverySpecCsvOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeZoneDiscoverySpecCsvOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewZoneDiscoverySpecCsvOptions(c *Client, des, nw *ZoneDiscoverySpecCsvOptions) *ZoneDiscoverySpecCsvOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ZoneDiscoverySpecCsvOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Delimiter, nw.Delimiter) {
		nw.Delimiter = des.Delimiter
	}
	if dcl.StringCanonicalize(des.Encoding, nw.Encoding) {
		nw.Encoding = des.Encoding
	}
	if dcl.BoolCanonicalize(des.DisableTypeInference, nw.DisableTypeInference) {
		nw.DisableTypeInference = des.DisableTypeInference
	}

	return nw
}

func canonicalizeNewZoneDiscoverySpecCsvOptionsSet(c *Client, des, nw []ZoneDiscoverySpecCsvOptions) []ZoneDiscoverySpecCsvOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ZoneDiscoverySpecCsvOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareZoneDiscoverySpecCsvOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewZoneDiscoverySpecCsvOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewZoneDiscoverySpecCsvOptionsSlice(c *Client, des, nw []ZoneDiscoverySpecCsvOptions) []ZoneDiscoverySpecCsvOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ZoneDiscoverySpecCsvOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewZoneDiscoverySpecCsvOptions(c, &d, &n))
	}

	return items
}

func canonicalizeZoneDiscoverySpecJsonOptions(des, initial *ZoneDiscoverySpecJsonOptions, opts ...dcl.ApplyOption) *ZoneDiscoverySpecJsonOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ZoneDiscoverySpecJsonOptions{}

	if dcl.StringCanonicalize(des.Encoding, initial.Encoding) || dcl.IsZeroValue(des.Encoding) {
		cDes.Encoding = initial.Encoding
	} else {
		cDes.Encoding = des.Encoding
	}
	if dcl.BoolCanonicalize(des.DisableTypeInference, initial.DisableTypeInference) || dcl.IsZeroValue(des.DisableTypeInference) {
		cDes.DisableTypeInference = initial.DisableTypeInference
	} else {
		cDes.DisableTypeInference = des.DisableTypeInference
	}

	return cDes
}

func canonicalizeZoneDiscoverySpecJsonOptionsSlice(des, initial []ZoneDiscoverySpecJsonOptions, opts ...dcl.ApplyOption) []ZoneDiscoverySpecJsonOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ZoneDiscoverySpecJsonOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeZoneDiscoverySpecJsonOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ZoneDiscoverySpecJsonOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeZoneDiscoverySpecJsonOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewZoneDiscoverySpecJsonOptions(c *Client, des, nw *ZoneDiscoverySpecJsonOptions) *ZoneDiscoverySpecJsonOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ZoneDiscoverySpecJsonOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Encoding, nw.Encoding) {
		nw.Encoding = des.Encoding
	}
	if dcl.BoolCanonicalize(des.DisableTypeInference, nw.DisableTypeInference) {
		nw.DisableTypeInference = des.DisableTypeInference
	}

	return nw
}

func canonicalizeNewZoneDiscoverySpecJsonOptionsSet(c *Client, des, nw []ZoneDiscoverySpecJsonOptions) []ZoneDiscoverySpecJsonOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ZoneDiscoverySpecJsonOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareZoneDiscoverySpecJsonOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewZoneDiscoverySpecJsonOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewZoneDiscoverySpecJsonOptionsSlice(c *Client, des, nw []ZoneDiscoverySpecJsonOptions) []ZoneDiscoverySpecJsonOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ZoneDiscoverySpecJsonOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewZoneDiscoverySpecJsonOptions(c, &d, &n))
	}

	return items
}

func canonicalizeZoneResourceSpec(des, initial *ZoneResourceSpec, opts ...dcl.ApplyOption) *ZoneResourceSpec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ZoneResourceSpec{}

	if dcl.IsZeroValue(des.LocationType) || (dcl.IsEmptyValueIndirect(des.LocationType) && dcl.IsEmptyValueIndirect(initial.LocationType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.LocationType = initial.LocationType
	} else {
		cDes.LocationType = des.LocationType
	}

	return cDes
}

func canonicalizeZoneResourceSpecSlice(des, initial []ZoneResourceSpec, opts ...dcl.ApplyOption) []ZoneResourceSpec {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ZoneResourceSpec, 0, len(des))
		for _, d := range des {
			cd := canonicalizeZoneResourceSpec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ZoneResourceSpec, 0, len(des))
	for i, d := range des {
		cd := canonicalizeZoneResourceSpec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewZoneResourceSpec(c *Client, des, nw *ZoneResourceSpec) *ZoneResourceSpec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ZoneResourceSpec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewZoneResourceSpecSet(c *Client, des, nw []ZoneResourceSpec) []ZoneResourceSpec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ZoneResourceSpec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareZoneResourceSpecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewZoneResourceSpec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewZoneResourceSpecSlice(c *Client, des, nw []ZoneResourceSpec) []ZoneResourceSpec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ZoneResourceSpec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewZoneResourceSpec(c, &d, &n))
	}

	return items
}

func canonicalizeZoneAssetStatus(des, initial *ZoneAssetStatus, opts ...dcl.ApplyOption) *ZoneAssetStatus {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ZoneAssetStatus{}

	if dcl.IsZeroValue(des.UpdateTime) || (dcl.IsEmptyValueIndirect(des.UpdateTime) && dcl.IsEmptyValueIndirect(initial.UpdateTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UpdateTime = initial.UpdateTime
	} else {
		cDes.UpdateTime = des.UpdateTime
	}
	if dcl.IsZeroValue(des.ActiveAssets) || (dcl.IsEmptyValueIndirect(des.ActiveAssets) && dcl.IsEmptyValueIndirect(initial.ActiveAssets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ActiveAssets = initial.ActiveAssets
	} else {
		cDes.ActiveAssets = des.ActiveAssets
	}
	if dcl.IsZeroValue(des.SecurityPolicyApplyingAssets) || (dcl.IsEmptyValueIndirect(des.SecurityPolicyApplyingAssets) && dcl.IsEmptyValueIndirect(initial.SecurityPolicyApplyingAssets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SecurityPolicyApplyingAssets = initial.SecurityPolicyApplyingAssets
	} else {
		cDes.SecurityPolicyApplyingAssets = des.SecurityPolicyApplyingAssets
	}

	return cDes
}

func canonicalizeZoneAssetStatusSlice(des, initial []ZoneAssetStatus, opts ...dcl.ApplyOption) []ZoneAssetStatus {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ZoneAssetStatus, 0, len(des))
		for _, d := range des {
			cd := canonicalizeZoneAssetStatus(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ZoneAssetStatus, 0, len(des))
	for i, d := range des {
		cd := canonicalizeZoneAssetStatus(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewZoneAssetStatus(c *Client, des, nw *ZoneAssetStatus) *ZoneAssetStatus {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ZoneAssetStatus while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewZoneAssetStatusSet(c *Client, des, nw []ZoneAssetStatus) []ZoneAssetStatus {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ZoneAssetStatus
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareZoneAssetStatusNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewZoneAssetStatus(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewZoneAssetStatusSlice(c *Client, des, nw []ZoneAssetStatus) []ZoneAssetStatus {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ZoneAssetStatus
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewZoneAssetStatus(c, &d, &n))
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
func diffZone(c *Client, desired, actual *Zone, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiscoverySpec, actual.DiscoverySpec, dcl.DiffInfo{ObjectFunction: compareZoneDiscoverySpecNewStyle, EmptyObject: EmptyZoneDiscoverySpec, OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("DiscoverySpec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ResourceSpec, actual.ResourceSpec, dcl.DiffInfo{ObjectFunction: compareZoneResourceSpecNewStyle, EmptyObject: EmptyZoneResourceSpec, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ResourceSpec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AssetStatus, actual.AssetStatus, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareZoneAssetStatusNewStyle, EmptyObject: EmptyZoneAssetStatus, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AssetStatus")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Lake, actual.Lake, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Lake")); len(ds) != 0 || err != nil {
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
func compareZoneDiscoverySpecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ZoneDiscoverySpec)
	if !ok {
		desiredNotPointer, ok := d.(ZoneDiscoverySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpec or *ZoneDiscoverySpec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ZoneDiscoverySpec)
	if !ok {
		actualNotPointer, ok := a.(ZoneDiscoverySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Enabled, actual.Enabled, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Enabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IncludePatterns, actual.IncludePatterns, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("IncludePatterns")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExcludePatterns, actual.ExcludePatterns, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("ExcludePatterns")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CsvOptions, actual.CsvOptions, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareZoneDiscoverySpecCsvOptionsNewStyle, EmptyObject: EmptyZoneDiscoverySpecCsvOptions, OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("CsvOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JsonOptions, actual.JsonOptions, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareZoneDiscoverySpecJsonOptionsNewStyle, EmptyObject: EmptyZoneDiscoverySpecJsonOptions, OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("JsonOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Schedule, actual.Schedule, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Schedule")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareZoneDiscoverySpecCsvOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ZoneDiscoverySpecCsvOptions)
	if !ok {
		desiredNotPointer, ok := d.(ZoneDiscoverySpecCsvOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpecCsvOptions or *ZoneDiscoverySpecCsvOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ZoneDiscoverySpecCsvOptions)
	if !ok {
		actualNotPointer, ok := a.(ZoneDiscoverySpecCsvOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpecCsvOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HeaderRows, actual.HeaderRows, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("HeaderRows")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Delimiter, actual.Delimiter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Delimiter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Encoding, actual.Encoding, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Encoding")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisableTypeInference, actual.DisableTypeInference, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("DisableTypeInference")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareZoneDiscoverySpecJsonOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ZoneDiscoverySpecJsonOptions)
	if !ok {
		desiredNotPointer, ok := d.(ZoneDiscoverySpecJsonOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpecJsonOptions or *ZoneDiscoverySpecJsonOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ZoneDiscoverySpecJsonOptions)
	if !ok {
		actualNotPointer, ok := a.(ZoneDiscoverySpecJsonOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneDiscoverySpecJsonOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Encoding, actual.Encoding, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("Encoding")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisableTypeInference, actual.DisableTypeInference, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("DisableTypeInference")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareZoneResourceSpecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ZoneResourceSpec)
	if !ok {
		desiredNotPointer, ok := d.(ZoneResourceSpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneResourceSpec or *ZoneResourceSpec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ZoneResourceSpec)
	if !ok {
		actualNotPointer, ok := a.(ZoneResourceSpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneResourceSpec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocationType, actual.LocationType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocationType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareZoneAssetStatusNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ZoneAssetStatus)
	if !ok {
		desiredNotPointer, ok := d.(ZoneAssetStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneAssetStatus or *ZoneAssetStatus", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ZoneAssetStatus)
	if !ok {
		actualNotPointer, ok := a.(ZoneAssetStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ZoneAssetStatus", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ActiveAssets, actual.ActiveAssets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("ActiveAssets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityPolicyApplyingAssets, actual.SecurityPolicyApplyingAssets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateZoneUpdateZoneOperation")}, fn.AddNest("SecurityPolicyApplyingAssets")); len(ds) != 0 || err != nil {
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
func (r *Zone) urlNormalized() *Zone {
	normalized := dcl.Copy(*r).(Zone)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Lake = dcl.SelfLinkToName(r.Lake)
	return &normalized
}

func (r *Zone) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateZone" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"lake":     dcl.ValueOrEmptyString(nr.Lake),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Zone resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Zone) marshal(c *Client) ([]byte, error) {
	m, err := expandZone(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Zone: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalZone decodes JSON responses into the Zone resource schema.
func unmarshalZone(b []byte, c *Client, res *Zone) (*Zone, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapZone(m, c, res)
}

func unmarshalMapZone(m map[string]interface{}, c *Client, res *Zone) (*Zone, error) {

	flattened := flattenZone(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandZone expands Zone into a JSON request object.
func expandZone(c *Client, f *Zone) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/lakes/%s/zones/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Lake), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Type; dcl.ValueShouldBeSent(v) {
		m["type"] = v
	}
	if v, err := expandZoneDiscoverySpec(c, f.DiscoverySpec, res); err != nil {
		return nil, fmt.Errorf("error expanding DiscoverySpec into discoverySpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["discoverySpec"] = v
	}
	if v, err := expandZoneResourceSpec(c, f.ResourceSpec, res); err != nil {
		return nil, fmt.Errorf("error expanding ResourceSpec into resourceSpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["resourceSpec"] = v
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
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Lake into lake: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["lake"] = v
	}

	return m, nil
}

// flattenZone flattens Zone from a JSON request object into the
// Zone type.
func flattenZone(c *Client, i interface{}, res *Zone) *Zone {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Zone{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.State = flattenZoneStateEnum(m["state"])
	resultRes.Type = flattenZoneTypeEnum(m["type"])
	resultRes.DiscoverySpec = flattenZoneDiscoverySpec(c, m["discoverySpec"], res)
	resultRes.ResourceSpec = flattenZoneResourceSpec(c, m["resourceSpec"], res)
	resultRes.AssetStatus = flattenZoneAssetStatus(c, m["assetStatus"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Lake = dcl.FlattenString(m["lake"])

	return resultRes
}

// expandZoneDiscoverySpecMap expands the contents of ZoneDiscoverySpec into a JSON
// request object.
func expandZoneDiscoverySpecMap(c *Client, f map[string]ZoneDiscoverySpec, res *Zone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandZoneDiscoverySpec(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandZoneDiscoverySpecSlice expands the contents of ZoneDiscoverySpec into a JSON
// request object.
func expandZoneDiscoverySpecSlice(c *Client, f []ZoneDiscoverySpec, res *Zone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandZoneDiscoverySpec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenZoneDiscoverySpecMap flattens the contents of ZoneDiscoverySpec from a JSON
// response object.
func flattenZoneDiscoverySpecMap(c *Client, i interface{}, res *Zone) map[string]ZoneDiscoverySpec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneDiscoverySpec{}
	}

	if len(a) == 0 {
		return map[string]ZoneDiscoverySpec{}
	}

	items := make(map[string]ZoneDiscoverySpec)
	for k, item := range a {
		items[k] = *flattenZoneDiscoverySpec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenZoneDiscoverySpecSlice flattens the contents of ZoneDiscoverySpec from a JSON
// response object.
func flattenZoneDiscoverySpecSlice(c *Client, i interface{}, res *Zone) []ZoneDiscoverySpec {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneDiscoverySpec{}
	}

	if len(a) == 0 {
		return []ZoneDiscoverySpec{}
	}

	items := make([]ZoneDiscoverySpec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneDiscoverySpec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandZoneDiscoverySpec expands an instance of ZoneDiscoverySpec into a JSON
// request object.
func expandZoneDiscoverySpec(c *Client, f *ZoneDiscoverySpec, res *Zone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Enabled; !dcl.IsEmptyValueIndirect(v) {
		m["enabled"] = v
	}
	if v := f.IncludePatterns; v != nil {
		m["includePatterns"] = v
	}
	if v := f.ExcludePatterns; v != nil {
		m["excludePatterns"] = v
	}
	if v, err := expandZoneDiscoverySpecCsvOptions(c, f.CsvOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CsvOptions into csvOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["csvOptions"] = v
	}
	if v, err := expandZoneDiscoverySpecJsonOptions(c, f.JsonOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding JsonOptions into jsonOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["jsonOptions"] = v
	}
	if v := f.Schedule; !dcl.IsEmptyValueIndirect(v) {
		m["schedule"] = v
	}

	return m, nil
}

// flattenZoneDiscoverySpec flattens an instance of ZoneDiscoverySpec from a JSON
// response object.
func flattenZoneDiscoverySpec(c *Client, i interface{}, res *Zone) *ZoneDiscoverySpec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ZoneDiscoverySpec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyZoneDiscoverySpec
	}
	r.Enabled = flattenZoneDiscoverySpecEnable(c, m["enabled"], res)
	r.IncludePatterns = dcl.FlattenStringSlice(m["includePatterns"])
	r.ExcludePatterns = dcl.FlattenStringSlice(m["excludePatterns"])
	r.CsvOptions = flattenZoneDiscoverySpecCsvOptions(c, m["csvOptions"], res)
	r.JsonOptions = flattenZoneDiscoverySpecJsonOptions(c, m["jsonOptions"], res)
	r.Schedule = dcl.FlattenString(m["schedule"])

	return r
}

// expandZoneDiscoverySpecCsvOptionsMap expands the contents of ZoneDiscoverySpecCsvOptions into a JSON
// request object.
func expandZoneDiscoverySpecCsvOptionsMap(c *Client, f map[string]ZoneDiscoverySpecCsvOptions, res *Zone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandZoneDiscoverySpecCsvOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandZoneDiscoverySpecCsvOptionsSlice expands the contents of ZoneDiscoverySpecCsvOptions into a JSON
// request object.
func expandZoneDiscoverySpecCsvOptionsSlice(c *Client, f []ZoneDiscoverySpecCsvOptions, res *Zone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandZoneDiscoverySpecCsvOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenZoneDiscoverySpecCsvOptionsMap flattens the contents of ZoneDiscoverySpecCsvOptions from a JSON
// response object.
func flattenZoneDiscoverySpecCsvOptionsMap(c *Client, i interface{}, res *Zone) map[string]ZoneDiscoverySpecCsvOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneDiscoverySpecCsvOptions{}
	}

	if len(a) == 0 {
		return map[string]ZoneDiscoverySpecCsvOptions{}
	}

	items := make(map[string]ZoneDiscoverySpecCsvOptions)
	for k, item := range a {
		items[k] = *flattenZoneDiscoverySpecCsvOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenZoneDiscoverySpecCsvOptionsSlice flattens the contents of ZoneDiscoverySpecCsvOptions from a JSON
// response object.
func flattenZoneDiscoverySpecCsvOptionsSlice(c *Client, i interface{}, res *Zone) []ZoneDiscoverySpecCsvOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneDiscoverySpecCsvOptions{}
	}

	if len(a) == 0 {
		return []ZoneDiscoverySpecCsvOptions{}
	}

	items := make([]ZoneDiscoverySpecCsvOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneDiscoverySpecCsvOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandZoneDiscoverySpecCsvOptions expands an instance of ZoneDiscoverySpecCsvOptions into a JSON
// request object.
func expandZoneDiscoverySpecCsvOptions(c *Client, f *ZoneDiscoverySpecCsvOptions, res *Zone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.HeaderRows; !dcl.IsEmptyValueIndirect(v) {
		m["headerRows"] = v
	}
	if v := f.Delimiter; !dcl.IsEmptyValueIndirect(v) {
		m["delimiter"] = v
	}
	if v := f.Encoding; !dcl.IsEmptyValueIndirect(v) {
		m["encoding"] = v
	}
	if v := f.DisableTypeInference; !dcl.IsEmptyValueIndirect(v) {
		m["disableTypeInference"] = v
	}

	return m, nil
}

// flattenZoneDiscoverySpecCsvOptions flattens an instance of ZoneDiscoverySpecCsvOptions from a JSON
// response object.
func flattenZoneDiscoverySpecCsvOptions(c *Client, i interface{}, res *Zone) *ZoneDiscoverySpecCsvOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ZoneDiscoverySpecCsvOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyZoneDiscoverySpecCsvOptions
	}
	r.HeaderRows = dcl.FlattenInteger(m["headerRows"])
	r.Delimiter = dcl.FlattenString(m["delimiter"])
	r.Encoding = dcl.FlattenString(m["encoding"])
	r.DisableTypeInference = dcl.FlattenBool(m["disableTypeInference"])

	return r
}

// expandZoneDiscoverySpecJsonOptionsMap expands the contents of ZoneDiscoverySpecJsonOptions into a JSON
// request object.
func expandZoneDiscoverySpecJsonOptionsMap(c *Client, f map[string]ZoneDiscoverySpecJsonOptions, res *Zone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandZoneDiscoverySpecJsonOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandZoneDiscoverySpecJsonOptionsSlice expands the contents of ZoneDiscoverySpecJsonOptions into a JSON
// request object.
func expandZoneDiscoverySpecJsonOptionsSlice(c *Client, f []ZoneDiscoverySpecJsonOptions, res *Zone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandZoneDiscoverySpecJsonOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenZoneDiscoverySpecJsonOptionsMap flattens the contents of ZoneDiscoverySpecJsonOptions from a JSON
// response object.
func flattenZoneDiscoverySpecJsonOptionsMap(c *Client, i interface{}, res *Zone) map[string]ZoneDiscoverySpecJsonOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneDiscoverySpecJsonOptions{}
	}

	if len(a) == 0 {
		return map[string]ZoneDiscoverySpecJsonOptions{}
	}

	items := make(map[string]ZoneDiscoverySpecJsonOptions)
	for k, item := range a {
		items[k] = *flattenZoneDiscoverySpecJsonOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenZoneDiscoverySpecJsonOptionsSlice flattens the contents of ZoneDiscoverySpecJsonOptions from a JSON
// response object.
func flattenZoneDiscoverySpecJsonOptionsSlice(c *Client, i interface{}, res *Zone) []ZoneDiscoverySpecJsonOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneDiscoverySpecJsonOptions{}
	}

	if len(a) == 0 {
		return []ZoneDiscoverySpecJsonOptions{}
	}

	items := make([]ZoneDiscoverySpecJsonOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneDiscoverySpecJsonOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandZoneDiscoverySpecJsonOptions expands an instance of ZoneDiscoverySpecJsonOptions into a JSON
// request object.
func expandZoneDiscoverySpecJsonOptions(c *Client, f *ZoneDiscoverySpecJsonOptions, res *Zone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Encoding; !dcl.IsEmptyValueIndirect(v) {
		m["encoding"] = v
	}
	if v := f.DisableTypeInference; !dcl.IsEmptyValueIndirect(v) {
		m["disableTypeInference"] = v
	}

	return m, nil
}

// flattenZoneDiscoverySpecJsonOptions flattens an instance of ZoneDiscoverySpecJsonOptions from a JSON
// response object.
func flattenZoneDiscoverySpecJsonOptions(c *Client, i interface{}, res *Zone) *ZoneDiscoverySpecJsonOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ZoneDiscoverySpecJsonOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyZoneDiscoverySpecJsonOptions
	}
	r.Encoding = dcl.FlattenString(m["encoding"])
	r.DisableTypeInference = dcl.FlattenBool(m["disableTypeInference"])

	return r
}

// expandZoneResourceSpecMap expands the contents of ZoneResourceSpec into a JSON
// request object.
func expandZoneResourceSpecMap(c *Client, f map[string]ZoneResourceSpec, res *Zone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandZoneResourceSpec(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandZoneResourceSpecSlice expands the contents of ZoneResourceSpec into a JSON
// request object.
func expandZoneResourceSpecSlice(c *Client, f []ZoneResourceSpec, res *Zone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandZoneResourceSpec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenZoneResourceSpecMap flattens the contents of ZoneResourceSpec from a JSON
// response object.
func flattenZoneResourceSpecMap(c *Client, i interface{}, res *Zone) map[string]ZoneResourceSpec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneResourceSpec{}
	}

	if len(a) == 0 {
		return map[string]ZoneResourceSpec{}
	}

	items := make(map[string]ZoneResourceSpec)
	for k, item := range a {
		items[k] = *flattenZoneResourceSpec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenZoneResourceSpecSlice flattens the contents of ZoneResourceSpec from a JSON
// response object.
func flattenZoneResourceSpecSlice(c *Client, i interface{}, res *Zone) []ZoneResourceSpec {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneResourceSpec{}
	}

	if len(a) == 0 {
		return []ZoneResourceSpec{}
	}

	items := make([]ZoneResourceSpec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneResourceSpec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandZoneResourceSpec expands an instance of ZoneResourceSpec into a JSON
// request object.
func expandZoneResourceSpec(c *Client, f *ZoneResourceSpec, res *Zone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocationType; !dcl.IsEmptyValueIndirect(v) {
		m["locationType"] = v
	}

	return m, nil
}

// flattenZoneResourceSpec flattens an instance of ZoneResourceSpec from a JSON
// response object.
func flattenZoneResourceSpec(c *Client, i interface{}, res *Zone) *ZoneResourceSpec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ZoneResourceSpec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyZoneResourceSpec
	}
	r.LocationType = flattenZoneResourceSpecLocationTypeEnum(m["locationType"])

	return r
}

// expandZoneAssetStatusMap expands the contents of ZoneAssetStatus into a JSON
// request object.
func expandZoneAssetStatusMap(c *Client, f map[string]ZoneAssetStatus, res *Zone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandZoneAssetStatus(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandZoneAssetStatusSlice expands the contents of ZoneAssetStatus into a JSON
// request object.
func expandZoneAssetStatusSlice(c *Client, f []ZoneAssetStatus, res *Zone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandZoneAssetStatus(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenZoneAssetStatusMap flattens the contents of ZoneAssetStatus from a JSON
// response object.
func flattenZoneAssetStatusMap(c *Client, i interface{}, res *Zone) map[string]ZoneAssetStatus {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneAssetStatus{}
	}

	if len(a) == 0 {
		return map[string]ZoneAssetStatus{}
	}

	items := make(map[string]ZoneAssetStatus)
	for k, item := range a {
		items[k] = *flattenZoneAssetStatus(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenZoneAssetStatusSlice flattens the contents of ZoneAssetStatus from a JSON
// response object.
func flattenZoneAssetStatusSlice(c *Client, i interface{}, res *Zone) []ZoneAssetStatus {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneAssetStatus{}
	}

	if len(a) == 0 {
		return []ZoneAssetStatus{}
	}

	items := make([]ZoneAssetStatus, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneAssetStatus(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandZoneAssetStatus expands an instance of ZoneAssetStatus into a JSON
// request object.
func expandZoneAssetStatus(c *Client, f *ZoneAssetStatus, res *Zone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.UpdateTime; !dcl.IsEmptyValueIndirect(v) {
		m["updateTime"] = v
	}
	if v := f.ActiveAssets; !dcl.IsEmptyValueIndirect(v) {
		m["activeAssets"] = v
	}
	if v := f.SecurityPolicyApplyingAssets; !dcl.IsEmptyValueIndirect(v) {
		m["securityPolicyApplyingAssets"] = v
	}

	return m, nil
}

// flattenZoneAssetStatus flattens an instance of ZoneAssetStatus from a JSON
// response object.
func flattenZoneAssetStatus(c *Client, i interface{}, res *Zone) *ZoneAssetStatus {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ZoneAssetStatus{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyZoneAssetStatus
	}
	r.UpdateTime = dcl.FlattenString(m["updateTime"])
	r.ActiveAssets = dcl.FlattenInteger(m["activeAssets"])
	r.SecurityPolicyApplyingAssets = dcl.FlattenInteger(m["securityPolicyApplyingAssets"])

	return r
}

// flattenZoneStateEnumMap flattens the contents of ZoneStateEnum from a JSON
// response object.
func flattenZoneStateEnumMap(c *Client, i interface{}, res *Zone) map[string]ZoneStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneStateEnum{}
	}

	if len(a) == 0 {
		return map[string]ZoneStateEnum{}
	}

	items := make(map[string]ZoneStateEnum)
	for k, item := range a {
		items[k] = *flattenZoneStateEnum(item.(interface{}))
	}

	return items
}

// flattenZoneStateEnumSlice flattens the contents of ZoneStateEnum from a JSON
// response object.
func flattenZoneStateEnumSlice(c *Client, i interface{}, res *Zone) []ZoneStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneStateEnum{}
	}

	if len(a) == 0 {
		return []ZoneStateEnum{}
	}

	items := make([]ZoneStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneStateEnum(item.(interface{})))
	}

	return items
}

// flattenZoneStateEnum asserts that an interface is a string, and returns a
// pointer to a *ZoneStateEnum with the same value as that string.
func flattenZoneStateEnum(i interface{}) *ZoneStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ZoneStateEnumRef(s)
}

// flattenZoneTypeEnumMap flattens the contents of ZoneTypeEnum from a JSON
// response object.
func flattenZoneTypeEnumMap(c *Client, i interface{}, res *Zone) map[string]ZoneTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ZoneTypeEnum{}
	}

	items := make(map[string]ZoneTypeEnum)
	for k, item := range a {
		items[k] = *flattenZoneTypeEnum(item.(interface{}))
	}

	return items
}

// flattenZoneTypeEnumSlice flattens the contents of ZoneTypeEnum from a JSON
// response object.
func flattenZoneTypeEnumSlice(c *Client, i interface{}, res *Zone) []ZoneTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneTypeEnum{}
	}

	if len(a) == 0 {
		return []ZoneTypeEnum{}
	}

	items := make([]ZoneTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneTypeEnum(item.(interface{})))
	}

	return items
}

// flattenZoneTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ZoneTypeEnum with the same value as that string.
func flattenZoneTypeEnum(i interface{}) *ZoneTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ZoneTypeEnumRef(s)
}

// flattenZoneResourceSpecLocationTypeEnumMap flattens the contents of ZoneResourceSpecLocationTypeEnum from a JSON
// response object.
func flattenZoneResourceSpecLocationTypeEnumMap(c *Client, i interface{}, res *Zone) map[string]ZoneResourceSpecLocationTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ZoneResourceSpecLocationTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ZoneResourceSpecLocationTypeEnum{}
	}

	items := make(map[string]ZoneResourceSpecLocationTypeEnum)
	for k, item := range a {
		items[k] = *flattenZoneResourceSpecLocationTypeEnum(item.(interface{}))
	}

	return items
}

// flattenZoneResourceSpecLocationTypeEnumSlice flattens the contents of ZoneResourceSpecLocationTypeEnum from a JSON
// response object.
func flattenZoneResourceSpecLocationTypeEnumSlice(c *Client, i interface{}, res *Zone) []ZoneResourceSpecLocationTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ZoneResourceSpecLocationTypeEnum{}
	}

	if len(a) == 0 {
		return []ZoneResourceSpecLocationTypeEnum{}
	}

	items := make([]ZoneResourceSpecLocationTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenZoneResourceSpecLocationTypeEnum(item.(interface{})))
	}

	return items
}

// flattenZoneResourceSpecLocationTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ZoneResourceSpecLocationTypeEnum with the same value as that string.
func flattenZoneResourceSpecLocationTypeEnum(i interface{}) *ZoneResourceSpecLocationTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ZoneResourceSpecLocationTypeEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Zone) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalZone(b, c, r)
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
		if nr.Lake == nil && ncr.Lake == nil {
			c.Config.Logger.Info("Both Lake fields null - considering equal.")
		} else if nr.Lake == nil || ncr.Lake == nil {
			c.Config.Logger.Info("Only one Lake field is null - considering unequal.")
			return false
		} else if *nr.Lake != *ncr.Lake {
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

type zoneDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         zoneApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToZoneDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]zoneDiff, error) {
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
	var diffs []zoneDiff
	// For each operation name, create a zoneDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := zoneDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToZoneApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToZoneApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (zoneApiOperation, error) {
	switch opName {

	case "updateZoneUpdateZoneOperation":
		return &updateZoneUpdateZoneOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractZoneFields(r *Zone) error {
	vDiscoverySpec := r.DiscoverySpec
	if vDiscoverySpec == nil {
		// note: explicitly not the empty object.
		vDiscoverySpec = &ZoneDiscoverySpec{}
	}
	if err := extractZoneDiscoverySpecFields(r, vDiscoverySpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiscoverySpec) {
		r.DiscoverySpec = vDiscoverySpec
	}
	vResourceSpec := r.ResourceSpec
	if vResourceSpec == nil {
		// note: explicitly not the empty object.
		vResourceSpec = &ZoneResourceSpec{}
	}
	if err := extractZoneResourceSpecFields(r, vResourceSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vResourceSpec) {
		r.ResourceSpec = vResourceSpec
	}
	vAssetStatus := r.AssetStatus
	if vAssetStatus == nil {
		// note: explicitly not the empty object.
		vAssetStatus = &ZoneAssetStatus{}
	}
	if err := extractZoneAssetStatusFields(r, vAssetStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAssetStatus) {
		r.AssetStatus = vAssetStatus
	}
	return nil
}
func extractZoneDiscoverySpecFields(r *Zone, o *ZoneDiscoverySpec) error {
	vCsvOptions := o.CsvOptions
	if vCsvOptions == nil {
		// note: explicitly not the empty object.
		vCsvOptions = &ZoneDiscoverySpecCsvOptions{}
	}
	if err := extractZoneDiscoverySpecCsvOptionsFields(r, vCsvOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCsvOptions) {
		o.CsvOptions = vCsvOptions
	}
	vJsonOptions := o.JsonOptions
	if vJsonOptions == nil {
		// note: explicitly not the empty object.
		vJsonOptions = &ZoneDiscoverySpecJsonOptions{}
	}
	if err := extractZoneDiscoverySpecJsonOptionsFields(r, vJsonOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vJsonOptions) {
		o.JsonOptions = vJsonOptions
	}
	return nil
}
func extractZoneDiscoverySpecCsvOptionsFields(r *Zone, o *ZoneDiscoverySpecCsvOptions) error {
	return nil
}
func extractZoneDiscoverySpecJsonOptionsFields(r *Zone, o *ZoneDiscoverySpecJsonOptions) error {
	return nil
}
func extractZoneResourceSpecFields(r *Zone, o *ZoneResourceSpec) error {
	return nil
}
func extractZoneAssetStatusFields(r *Zone, o *ZoneAssetStatus) error {
	return nil
}

func postReadExtractZoneFields(r *Zone) error {
	vDiscoverySpec := r.DiscoverySpec
	if vDiscoverySpec == nil {
		// note: explicitly not the empty object.
		vDiscoverySpec = &ZoneDiscoverySpec{}
	}
	if err := postReadExtractZoneDiscoverySpecFields(r, vDiscoverySpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiscoverySpec) {
		r.DiscoverySpec = vDiscoverySpec
	}
	vResourceSpec := r.ResourceSpec
	if vResourceSpec == nil {
		// note: explicitly not the empty object.
		vResourceSpec = &ZoneResourceSpec{}
	}
	if err := postReadExtractZoneResourceSpecFields(r, vResourceSpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vResourceSpec) {
		r.ResourceSpec = vResourceSpec
	}
	vAssetStatus := r.AssetStatus
	if vAssetStatus == nil {
		// note: explicitly not the empty object.
		vAssetStatus = &ZoneAssetStatus{}
	}
	if err := postReadExtractZoneAssetStatusFields(r, vAssetStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAssetStatus) {
		r.AssetStatus = vAssetStatus
	}
	return nil
}
func postReadExtractZoneDiscoverySpecFields(r *Zone, o *ZoneDiscoverySpec) error {
	vCsvOptions := o.CsvOptions
	if vCsvOptions == nil {
		// note: explicitly not the empty object.
		vCsvOptions = &ZoneDiscoverySpecCsvOptions{}
	}
	if err := extractZoneDiscoverySpecCsvOptionsFields(r, vCsvOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCsvOptions) {
		o.CsvOptions = vCsvOptions
	}
	vJsonOptions := o.JsonOptions
	if vJsonOptions == nil {
		// note: explicitly not the empty object.
		vJsonOptions = &ZoneDiscoverySpecJsonOptions{}
	}
	if err := extractZoneDiscoverySpecJsonOptionsFields(r, vJsonOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vJsonOptions) {
		o.JsonOptions = vJsonOptions
	}
	return nil
}
func postReadExtractZoneDiscoverySpecCsvOptionsFields(r *Zone, o *ZoneDiscoverySpecCsvOptions) error {
	return nil
}
func postReadExtractZoneDiscoverySpecJsonOptionsFields(r *Zone, o *ZoneDiscoverySpecJsonOptions) error {
	return nil
}
func postReadExtractZoneResourceSpecFields(r *Zone, o *ZoneResourceSpec) error {
	return nil
}
func postReadExtractZoneAssetStatusFields(r *Zone, o *ZoneAssetStatus) error {
	return nil
}
