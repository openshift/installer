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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Service) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Custom) {
		if err := r.Custom.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Telemetry) {
		if err := r.Telemetry.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceCustom) validate() error {
	return nil
}
func (r *ServiceTelemetry) validate() error {
	return nil
}
func (r *Service) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://monitoring.googleapis.com/v3/", params)
}

func (r *Service) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Service) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/services", nr.basePath(), userBasePath, params), nil

}

func (r *Service) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services?serviceId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Service) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services/{{name}}", nr.basePath(), userBasePath, params), nil
}

// serviceApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type serviceApiOperation interface {
	do(context.Context, *Service, *Client) error
}

// newUpdateServiceUpdateServiceRequest creates a request for an
// Service resource's UpdateService update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateServiceUpdateServiceRequest(ctx context.Context, f *Service, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v, err := expandServiceCustom(c, f.Custom, res); err != nil {
		return nil, fmt.Errorf("error expanding Custom into custom: %w", err)
	} else if v != nil {
		req["custom"] = v
	}
	if v, err := expandServiceTelemetry(c, f.Telemetry, res); err != nil {
		return nil, fmt.Errorf("error expanding Telemetry into telemetry: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["telemetry"] = v
	}
	if v := f.UserLabels; !dcl.IsEmptyValueIndirect(v) {
		req["userLabels"] = v
	}
	return req, nil
}

// marshalUpdateServiceUpdateServiceRequest converts the update into
// the final JSON request body.
func marshalUpdateServiceUpdateServiceRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateServiceUpdateServiceOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateServiceUpdateServiceOperation) do(ctx context.Context, r *Service, c *Client) error {
	_, err := c.GetService(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateService")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateServiceUpdateServiceRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateServiceUpdateServiceRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listServiceRaw(ctx context.Context, r *Service, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ServiceMaxPage {
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

type listServiceOperation struct {
	Services []map[string]interface{} `json:"services"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listService(ctx context.Context, r *Service, pageToken string, pageSize int32) ([]*Service, string, error) {
	b, err := c.listServiceRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listServiceOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Service
	for _, v := range m.Services {
		res, err := unmarshalMapService(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllService(ctx context.Context, f func(*Service) bool, resources []*Service) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteService(ctx, res)
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

type deleteServiceOperation struct{}

func (op *deleteServiceOperation) do(ctx context.Context, r *Service, c *Client) error {
	r, err := c.GetService(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Service not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetService checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete Service: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createServiceOperation struct {
	response map[string]interface{}
}

func (op *createServiceOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createServiceOperation) do(ctx context.Context, r *Service, c *Client) error {
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

	if _, err := c.GetService(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getServiceRaw(ctx context.Context, r *Service) ([]byte, error) {
	if dcl.IsZeroValue(r.Custom) {
		r.Custom = &ServiceCustom{}
	}

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

func (c *Client) serviceDiffsForRawDesired(ctx context.Context, rawDesired *Service, opts ...dcl.ApplyOption) (initial, desired *Service, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Service
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Service); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Service, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetService(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Service resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Service resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Service resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeServiceDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Service: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Service: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractServiceFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeServiceInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Service: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeServiceDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Service: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffService(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeServiceInitialState(rawInitial, rawDesired *Service) (*Service, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeServiceDesiredState(rawDesired, rawInitial *Service, opts ...dcl.ApplyOption) (*Service, error) {

	if dcl.IsZeroValue(rawDesired.Custom) {
		rawDesired.Custom = &ServiceCustom{}
	}

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Custom = canonicalizeServiceCustom(rawDesired.Custom, nil, opts...)
		rawDesired.Telemetry = canonicalizeServiceTelemetry(rawDesired.Telemetry, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Service{}
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
	canonicalDesired.Custom = canonicalizeServiceCustom(rawDesired.Custom, rawInitial.Custom, opts...)
	canonicalDesired.Telemetry = canonicalizeServiceTelemetry(rawDesired.Telemetry, rawInitial.Telemetry, opts...)
	if dcl.IsZeroValue(rawDesired.UserLabels) || (dcl.IsEmptyValueIndirect(rawDesired.UserLabels) && dcl.IsEmptyValueIndirect(rawInitial.UserLabels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.UserLabels = rawInitial.UserLabels
	} else {
		canonicalDesired.UserLabels = rawDesired.UserLabels
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeServiceNewState(c *Client, rawNew, rawDesired *Service) (*Service, error) {

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

	if dcl.IsEmptyValueIndirect(rawNew.Custom) && dcl.IsEmptyValueIndirect(rawDesired.Custom) {
		rawNew.Custom = rawDesired.Custom
	} else {
		rawNew.Custom = canonicalizeNewServiceCustom(c, rawDesired.Custom, rawNew.Custom)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Telemetry) && dcl.IsEmptyValueIndirect(rawDesired.Telemetry) {
		rawNew.Telemetry = rawDesired.Telemetry
	} else {
		rawNew.Telemetry = canonicalizeNewServiceTelemetry(c, rawDesired.Telemetry, rawNew.Telemetry)
	}

	if dcl.IsEmptyValueIndirect(rawNew.UserLabels) && dcl.IsEmptyValueIndirect(rawDesired.UserLabels) {
		rawNew.UserLabels = rawDesired.UserLabels
	} else {
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeServiceCustom(des, initial *ServiceCustom, opts ...dcl.ApplyOption) *ServiceCustom {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &ServiceCustom{}

	return cDes
}

func canonicalizeServiceCustomSlice(des, initial []ServiceCustom, opts ...dcl.ApplyOption) []ServiceCustom {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceCustom, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceCustom(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceCustom, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceCustom(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceCustom(c *Client, des, nw *ServiceCustom) *ServiceCustom {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceCustom while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceCustomSet(c *Client, des, nw []ServiceCustom) []ServiceCustom {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceCustom
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceCustomNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceCustom(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceCustomSlice(c *Client, des, nw []ServiceCustom) []ServiceCustom {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceCustom
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceCustom(c, &d, &n))
	}

	return items
}

func canonicalizeServiceTelemetry(des, initial *ServiceTelemetry, opts ...dcl.ApplyOption) *ServiceTelemetry {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceTelemetry{}

	if dcl.StringCanonicalize(des.ResourceName, initial.ResourceName) || dcl.IsZeroValue(des.ResourceName) {
		cDes.ResourceName = initial.ResourceName
	} else {
		cDes.ResourceName = des.ResourceName
	}

	return cDes
}

func canonicalizeServiceTelemetrySlice(des, initial []ServiceTelemetry, opts ...dcl.ApplyOption) []ServiceTelemetry {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceTelemetry, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceTelemetry(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceTelemetry, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceTelemetry(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceTelemetry(c *Client, des, nw *ServiceTelemetry) *ServiceTelemetry {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceTelemetry while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ResourceName, nw.ResourceName) {
		nw.ResourceName = des.ResourceName
	}

	return nw
}

func canonicalizeNewServiceTelemetrySet(c *Client, des, nw []ServiceTelemetry) []ServiceTelemetry {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceTelemetry
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceTelemetryNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceTelemetry(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceTelemetrySlice(c *Client, des, nw []ServiceTelemetry) []ServiceTelemetry {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceTelemetry
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceTelemetry(c, &d, &n))
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
func diffService(c *Client, desired, actual *Service, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceUpdateServiceOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Custom, actual.Custom, dcl.DiffInfo{ObjectFunction: compareServiceCustomNewStyle, EmptyObject: EmptyServiceCustom, OperationSelector: dcl.TriggersOperation("updateServiceUpdateServiceOperation")}, fn.AddNest("Custom")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Telemetry, actual.Telemetry, dcl.DiffInfo{ObjectFunction: compareServiceTelemetryNewStyle, EmptyObject: EmptyServiceTelemetry, OperationSelector: dcl.TriggersOperation("updateServiceUpdateServiceOperation")}, fn.AddNest("Telemetry")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UserLabels, actual.UserLabels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceUpdateServiceOperation")}, fn.AddNest("UserLabels")); len(ds) != 0 || err != nil {
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
func compareServiceCustomNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func compareServiceTelemetryNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceTelemetry)
	if !ok {
		desiredNotPointer, ok := d.(ServiceTelemetry)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceTelemetry or *ServiceTelemetry", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceTelemetry)
	if !ok {
		actualNotPointer, ok := a.(ServiceTelemetry)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceTelemetry", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ResourceName, actual.ResourceName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceUpdateServiceOperation")}, fn.AddNest("ResourceName")); len(ds) != 0 || err != nil {
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
func (r *Service) urlNormalized() *Service {
	normalized := dcl.Copy(*r).(Service)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Service) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateService" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/services/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Service resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Service) marshal(c *Client) ([]byte, error) {
	m, err := expandService(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Service: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalService decodes JSON responses into the Service resource schema.
func unmarshalService(b []byte, c *Client, res *Service) (*Service, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapService(m, c, res)
}

func unmarshalMapService(m map[string]interface{}, c *Client, res *Service) (*Service, error) {

	flattened := flattenService(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandService expands Service into a JSON request object.
func expandService(c *Client, f *Service) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/services/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v, err := expandServiceCustom(c, f.Custom, res); err != nil {
		return nil, fmt.Errorf("error expanding Custom into custom: %w", err)
	} else if v != nil {
		m["custom"] = v
	}
	if v, err := expandServiceTelemetry(c, f.Telemetry, res); err != nil {
		return nil, fmt.Errorf("error expanding Telemetry into telemetry: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["telemetry"] = v
	}
	if v := f.UserLabels; dcl.ValueShouldBeSent(v) {
		m["userLabels"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenService flattens Service from a JSON request object into the
// Service type.
func flattenService(c *Client, i interface{}, res *Service) *Service {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Service{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.Custom = flattenServiceCustom(c, m["custom"], res)
	if _, ok := m["custom"]; !ok {
		c.Config.Logger.Info("Using default value for custom")
		resultRes.Custom = &ServiceCustom{}
	}
	resultRes.Telemetry = flattenServiceTelemetry(c, m["telemetry"], res)
	resultRes.UserLabels = dcl.FlattenKeyValuePairs(m["userLabels"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandServiceCustomMap expands the contents of ServiceCustom into a JSON
// request object.
func expandServiceCustomMap(c *Client, f map[string]ServiceCustom, res *Service) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceCustom(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceCustomSlice expands the contents of ServiceCustom into a JSON
// request object.
func expandServiceCustomSlice(c *Client, f []ServiceCustom, res *Service) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceCustom(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceCustomMap flattens the contents of ServiceCustom from a JSON
// response object.
func flattenServiceCustomMap(c *Client, i interface{}, res *Service) map[string]ServiceCustom {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceCustom{}
	}

	if len(a) == 0 {
		return map[string]ServiceCustom{}
	}

	items := make(map[string]ServiceCustom)
	for k, item := range a {
		items[k] = *flattenServiceCustom(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceCustomSlice flattens the contents of ServiceCustom from a JSON
// response object.
func flattenServiceCustomSlice(c *Client, i interface{}, res *Service) []ServiceCustom {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceCustom{}
	}

	if len(a) == 0 {
		return []ServiceCustom{}
	}

	items := make([]ServiceCustom, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceCustom(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceCustom expands an instance of ServiceCustom into a JSON
// request object.
func expandServiceCustom(c *Client, f *ServiceCustom, res *Service) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceCustom flattens an instance of ServiceCustom from a JSON
// response object.
func flattenServiceCustom(c *Client, i interface{}, res *Service) *ServiceCustom {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceCustom{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceCustom
	}

	return r
}

// expandServiceTelemetryMap expands the contents of ServiceTelemetry into a JSON
// request object.
func expandServiceTelemetryMap(c *Client, f map[string]ServiceTelemetry, res *Service) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceTelemetry(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceTelemetrySlice expands the contents of ServiceTelemetry into a JSON
// request object.
func expandServiceTelemetrySlice(c *Client, f []ServiceTelemetry, res *Service) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceTelemetry(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceTelemetryMap flattens the contents of ServiceTelemetry from a JSON
// response object.
func flattenServiceTelemetryMap(c *Client, i interface{}, res *Service) map[string]ServiceTelemetry {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceTelemetry{}
	}

	if len(a) == 0 {
		return map[string]ServiceTelemetry{}
	}

	items := make(map[string]ServiceTelemetry)
	for k, item := range a {
		items[k] = *flattenServiceTelemetry(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceTelemetrySlice flattens the contents of ServiceTelemetry from a JSON
// response object.
func flattenServiceTelemetrySlice(c *Client, i interface{}, res *Service) []ServiceTelemetry {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceTelemetry{}
	}

	if len(a) == 0 {
		return []ServiceTelemetry{}
	}

	items := make([]ServiceTelemetry, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceTelemetry(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceTelemetry expands an instance of ServiceTelemetry into a JSON
// request object.
func expandServiceTelemetry(c *Client, f *ServiceTelemetry, res *Service) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ResourceName; !dcl.IsEmptyValueIndirect(v) {
		m["resourceName"] = v
	}

	return m, nil
}

// flattenServiceTelemetry flattens an instance of ServiceTelemetry from a JSON
// response object.
func flattenServiceTelemetry(c *Client, i interface{}, res *Service) *ServiceTelemetry {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceTelemetry{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceTelemetry
	}
	r.ResourceName = dcl.FlattenString(m["resourceName"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Service) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalService(b, c, r)
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

type serviceDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         serviceApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToServiceDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]serviceDiff, error) {
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
	var diffs []serviceDiff
	// For each operation name, create a serviceDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := serviceDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToServiceApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToServiceApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (serviceApiOperation, error) {
	switch opName {

	case "updateServiceUpdateServiceOperation":
		return &updateServiceUpdateServiceOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractServiceFields(r *Service) error {
	vCustom := r.Custom
	if vCustom == nil {
		// note: explicitly not the empty object.
		vCustom = &ServiceCustom{}
	}
	if err := extractServiceCustomFields(r, vCustom); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCustom) {
		r.Custom = vCustom
	}
	vTelemetry := r.Telemetry
	if vTelemetry == nil {
		// note: explicitly not the empty object.
		vTelemetry = &ServiceTelemetry{}
	}
	if err := extractServiceTelemetryFields(r, vTelemetry); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTelemetry) {
		r.Telemetry = vTelemetry
	}
	return nil
}
func extractServiceCustomFields(r *Service, o *ServiceCustom) error {
	return nil
}
func extractServiceTelemetryFields(r *Service, o *ServiceTelemetry) error {
	return nil
}

func postReadExtractServiceFields(r *Service) error {
	vCustom := r.Custom
	if vCustom == nil {
		// note: explicitly not the empty object.
		vCustom = &ServiceCustom{}
	}
	if err := postReadExtractServiceCustomFields(r, vCustom); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCustom) {
		r.Custom = vCustom
	}
	vTelemetry := r.Telemetry
	if vTelemetry == nil {
		// note: explicitly not the empty object.
		vTelemetry = &ServiceTelemetry{}
	}
	if err := postReadExtractServiceTelemetryFields(r, vTelemetry); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTelemetry) {
		r.Telemetry = vTelemetry
	}
	return nil
}
func postReadExtractServiceCustomFields(r *Service, o *ServiceCustom) error {
	return nil
}
func postReadExtractServiceTelemetryFields(r *Service, o *ServiceTelemetry) error {
	return nil
}
