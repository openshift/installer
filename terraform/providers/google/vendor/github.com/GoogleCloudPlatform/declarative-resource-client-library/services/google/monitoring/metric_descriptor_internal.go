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
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *MetricDescriptor) validate() error {

	if err := dcl.Required(r, "type"); err != nil {
		return err
	}
	if err := dcl.Required(r, "metricKind"); err != nil {
		return err
	}
	if err := dcl.Required(r, "valueType"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Metadata) {
		if err := r.Metadata.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *MetricDescriptorLabels) validate() error {
	return nil
}
func (r *MetricDescriptorMetadata) validate() error {
	return nil
}
func (r *MetricDescriptor) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://monitoring.googleapis.com/v3/", params)
}

func (r *MetricDescriptor) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"type":    dcl.ValueOrEmptyString(nr.Type),
	}
	return dcl.URL("projects/{{project}}/metricDescriptors/{{type}}", nr.basePath(), userBasePath, params), nil
}

func (r *MetricDescriptor) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/metricDescriptors", nr.basePath(), userBasePath, params), nil

}

func (r *MetricDescriptor) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/metricDescriptors", nr.basePath(), userBasePath, params), nil

}

func (r *MetricDescriptor) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"type":    dcl.ValueOrEmptyString(nr.Type),
	}
	return dcl.URL("projects/{{project}}/metricDescriptors/{{type}}", nr.basePath(), userBasePath, params), nil
}

// metricDescriptorApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type metricDescriptorApiOperation interface {
	do(context.Context, *MetricDescriptor, *Client) error
}

func (c *Client) listMetricDescriptorRaw(ctx context.Context, r *MetricDescriptor, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != MetricDescriptorMaxPage {
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

type listMetricDescriptorOperation struct {
	MetricDescriptors []map[string]interface{} `json:"metricDescriptors"`
	Token             string                   `json:"nextPageToken"`
}

func (c *Client) listMetricDescriptor(ctx context.Context, r *MetricDescriptor, pageToken string, pageSize int32) ([]*MetricDescriptor, string, error) {
	b, err := c.listMetricDescriptorRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listMetricDescriptorOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*MetricDescriptor
	for _, v := range m.MetricDescriptors {
		res, err := unmarshalMapMetricDescriptor(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllMetricDescriptor(ctx context.Context, f func(*MetricDescriptor) bool, resources []*MetricDescriptor) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteMetricDescriptor(ctx, res)
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

type deleteMetricDescriptorOperation struct{}

func (op *deleteMetricDescriptorOperation) do(ctx context.Context, r *MetricDescriptor, c *Client) error {
	r, err := c.GetMetricDescriptor(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "MetricDescriptor not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetMetricDescriptor checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete MetricDescriptor: %w", err)
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetMetricDescriptor(ctx, r)
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
type createMetricDescriptorOperation struct {
	response map[string]interface{}
}

func (op *createMetricDescriptorOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createMetricDescriptorOperation) do(ctx context.Context, r *MetricDescriptor, c *Client) error {
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

	// Poll for the MetricDescriptor resource to be created. MetricDescriptor resources are eventually consistent but do not support operations
	// so we must repeatedly poll to check for their creation.
	requiredSuccesses := 10
	start := time.Now()
	err = dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		u, err := r.getURL(c.Config.BasePath)
		if err != nil {
			return nil, err
		}
		getResp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, nil)
		if err != nil {
			// If the error is a transient server error (e.g., 500) or not found (i.e., the resource has not yet been created),
			// continue retrying until the transient error is resolved, the resource is created, or we time out.
			if dcl.IsRetryableRequestError(c.Config, err, true, start) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		getResp.Response.Body.Close()
		requiredSuccesses--
		if requiredSuccesses > 0 {
			return &dcl.RetryDetails{}, dcl.OperationNotDone{}
		}
		return getResp, nil
	}, c.Config.RetryProvider)

	if _, err := c.GetMetricDescriptor(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getMetricDescriptorRaw(ctx context.Context, r *MetricDescriptor) ([]byte, error) {

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

func (c *Client) metricDescriptorDiffsForRawDesired(ctx context.Context, rawDesired *MetricDescriptor, opts ...dcl.ApplyOption) (initial, desired *MetricDescriptor, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *MetricDescriptor
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*MetricDescriptor); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected MetricDescriptor, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetMetricDescriptor(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a MetricDescriptor resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve MetricDescriptor resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that MetricDescriptor resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeMetricDescriptorDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for MetricDescriptor: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for MetricDescriptor: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractMetricDescriptorFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeMetricDescriptorInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for MetricDescriptor: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeMetricDescriptorDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for MetricDescriptor: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffMetricDescriptor(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeMetricDescriptorInitialState(rawInitial, rawDesired *MetricDescriptor) (*MetricDescriptor, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeMetricDescriptorDesiredState(rawDesired, rawInitial *MetricDescriptor, opts ...dcl.ApplyOption) (*MetricDescriptor, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Metadata = canonicalizeMetricDescriptorMetadata(rawDesired.Metadata, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &MetricDescriptor{}
	if dcl.StringCanonicalize(rawDesired.Type, rawInitial.Type) {
		canonicalDesired.Type = rawInitial.Type
	} else {
		canonicalDesired.Type = rawDesired.Type
	}
	canonicalDesired.Labels = canonicalizeMetricDescriptorLabelsSlice(rawDesired.Labels, rawInitial.Labels, opts...)
	if dcl.IsZeroValue(rawDesired.MetricKind) || (dcl.IsEmptyValueIndirect(rawDesired.MetricKind) && dcl.IsEmptyValueIndirect(rawInitial.MetricKind)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.MetricKind = rawInitial.MetricKind
	} else {
		canonicalDesired.MetricKind = rawDesired.MetricKind
	}
	if canonicalizeMetricDescriptorValueType(rawDesired.ValueType, rawInitial.ValueType) {
		canonicalDesired.ValueType = rawInitial.ValueType
	} else {
		canonicalDesired.ValueType = rawDesired.ValueType
	}
	if dcl.StringCanonicalize(rawDesired.Unit, rawInitial.Unit) {
		canonicalDesired.Unit = rawInitial.Unit
	} else {
		canonicalDesired.Unit = rawDesired.Unit
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	canonicalDesired.Metadata = canonicalizeMetricDescriptorMetadata(rawDesired.Metadata, rawInitial.Metadata, opts...)
	if dcl.IsZeroValue(rawDesired.LaunchStage) || (dcl.IsEmptyValueIndirect(rawDesired.LaunchStage) && dcl.IsEmptyValueIndirect(rawInitial.LaunchStage)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.LaunchStage = rawInitial.LaunchStage
	} else {
		canonicalDesired.LaunchStage = rawDesired.LaunchStage
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeMetricDescriptorNewState(c *Client, rawNew, rawDesired *MetricDescriptor) (*MetricDescriptor, error) {

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Type) && dcl.IsEmptyValueIndirect(rawDesired.Type) {
		rawNew.Type = rawDesired.Type
	} else {
		if dcl.StringCanonicalize(rawDesired.Type, rawNew.Type) {
			rawNew.Type = rawDesired.Type
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
		rawNew.Labels = canonicalizeNewMetricDescriptorLabelsSet(c, rawDesired.Labels, rawNew.Labels)
	}

	if dcl.IsEmptyValueIndirect(rawNew.MetricKind) && dcl.IsEmptyValueIndirect(rawDesired.MetricKind) {
		rawNew.MetricKind = rawDesired.MetricKind
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ValueType) && dcl.IsEmptyValueIndirect(rawDesired.ValueType) {
		rawNew.ValueType = rawDesired.ValueType
	} else {
		if canonicalizeMetricDescriptorValueType(rawDesired.ValueType, rawNew.ValueType) {
			rawNew.ValueType = rawDesired.ValueType
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Unit) && dcl.IsEmptyValueIndirect(rawDesired.Unit) {
		rawNew.Unit = rawDesired.Unit
	} else {
		if dcl.StringCanonicalize(rawDesired.Unit, rawNew.Unit) {
			rawNew.Unit = rawDesired.Unit
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	rawNew.Metadata = rawDesired.Metadata

	rawNew.LaunchStage = rawDesired.LaunchStage

	if dcl.IsEmptyValueIndirect(rawNew.MonitoredResourceTypes) && dcl.IsEmptyValueIndirect(rawDesired.MonitoredResourceTypes) {
		rawNew.MonitoredResourceTypes = rawDesired.MonitoredResourceTypes
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.MonitoredResourceTypes, rawNew.MonitoredResourceTypes) {
			rawNew.MonitoredResourceTypes = rawDesired.MonitoredResourceTypes
		}
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeMetricDescriptorLabels(des, initial *MetricDescriptorLabels, opts ...dcl.ApplyOption) *MetricDescriptorLabels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &MetricDescriptorLabels{}

	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if canonicalizeMetricDescriptorLabelsValueType(des.ValueType, initial.ValueType) || dcl.IsZeroValue(des.ValueType) {
		cDes.ValueType = initial.ValueType
	} else {
		cDes.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}

	return cDes
}

func canonicalizeMetricDescriptorLabelsSlice(des, initial []MetricDescriptorLabels, opts ...dcl.ApplyOption) []MetricDescriptorLabels {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]MetricDescriptorLabels, 0, len(des))
		for _, d := range des {
			cd := canonicalizeMetricDescriptorLabels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]MetricDescriptorLabels, 0, len(des))
	for i, d := range des {
		cd := canonicalizeMetricDescriptorLabels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewMetricDescriptorLabels(c *Client, des, nw *MetricDescriptorLabels) *MetricDescriptorLabels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for MetricDescriptorLabels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}
	if canonicalizeMetricDescriptorLabelsValueType(des.ValueType, nw.ValueType) {
		nw.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}

	return nw
}

func canonicalizeNewMetricDescriptorLabelsSet(c *Client, des, nw []MetricDescriptorLabels) []MetricDescriptorLabels {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []MetricDescriptorLabels
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareMetricDescriptorLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewMetricDescriptorLabels(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewMetricDescriptorLabelsSlice(c *Client, des, nw []MetricDescriptorLabels) []MetricDescriptorLabels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []MetricDescriptorLabels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewMetricDescriptorLabels(c, &d, &n))
	}

	return items
}

func canonicalizeMetricDescriptorMetadata(des, initial *MetricDescriptorMetadata, opts ...dcl.ApplyOption) *MetricDescriptorMetadata {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &MetricDescriptorMetadata{}

	if dcl.IsZeroValue(des.LaunchStage) || (dcl.IsEmptyValueIndirect(des.LaunchStage) && dcl.IsEmptyValueIndirect(initial.LaunchStage)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.LaunchStage = initial.LaunchStage
	} else {
		cDes.LaunchStage = des.LaunchStage
	}
	if dcl.StringCanonicalize(des.SamplePeriod, initial.SamplePeriod) || dcl.IsZeroValue(des.SamplePeriod) {
		cDes.SamplePeriod = initial.SamplePeriod
	} else {
		cDes.SamplePeriod = des.SamplePeriod
	}
	if dcl.StringCanonicalize(des.IngestDelay, initial.IngestDelay) || dcl.IsZeroValue(des.IngestDelay) {
		cDes.IngestDelay = initial.IngestDelay
	} else {
		cDes.IngestDelay = des.IngestDelay
	}

	return cDes
}

func canonicalizeMetricDescriptorMetadataSlice(des, initial []MetricDescriptorMetadata, opts ...dcl.ApplyOption) []MetricDescriptorMetadata {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]MetricDescriptorMetadata, 0, len(des))
		for _, d := range des {
			cd := canonicalizeMetricDescriptorMetadata(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]MetricDescriptorMetadata, 0, len(des))
	for i, d := range des {
		cd := canonicalizeMetricDescriptorMetadata(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewMetricDescriptorMetadata(c *Client, des, nw *MetricDescriptorMetadata) *MetricDescriptorMetadata {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for MetricDescriptorMetadata while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.SamplePeriod, nw.SamplePeriod) {
		nw.SamplePeriod = des.SamplePeriod
	}
	if dcl.StringCanonicalize(des.IngestDelay, nw.IngestDelay) {
		nw.IngestDelay = des.IngestDelay
	}

	return nw
}

func canonicalizeNewMetricDescriptorMetadataSet(c *Client, des, nw []MetricDescriptorMetadata) []MetricDescriptorMetadata {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []MetricDescriptorMetadata
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareMetricDescriptorMetadataNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewMetricDescriptorMetadata(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewMetricDescriptorMetadataSlice(c *Client, des, nw []MetricDescriptorMetadata) []MetricDescriptorMetadata {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []MetricDescriptorMetadata
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewMetricDescriptorMetadata(c, &d, &n))
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
func diffMetricDescriptor(c *Client, desired, actual *MetricDescriptor, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{Type: "Set", ObjectFunction: compareMetricDescriptorLabelsNewStyle, EmptyObject: EmptyMetricDescriptorLabels, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricKind, actual.MetricKind, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricKind")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ValueType, actual.ValueType, dcl.DiffInfo{Type: "EnumType", CustomDiff: canonicalizeMetricDescriptorValueType, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ValueType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Unit, actual.Unit, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Unit")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{Ignore: true, ObjectFunction: compareMetricDescriptorMetadataNewStyle, EmptyObject: EmptyMetricDescriptorMetadata, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LaunchStage, actual.LaunchStage, dcl.DiffInfo{Ignore: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LaunchStage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MonitoredResourceTypes, actual.MonitoredResourceTypes, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MonitoredResourceTypes")); len(ds) != 0 || err != nil {
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
func compareMetricDescriptorLabelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*MetricDescriptorLabels)
	if !ok {
		desiredNotPointer, ok := d.(MetricDescriptorLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricDescriptorLabels or *MetricDescriptorLabels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*MetricDescriptorLabels)
	if !ok {
		actualNotPointer, ok := a.(MetricDescriptorLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricDescriptorLabels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ValueType, actual.ValueType, dcl.DiffInfo{Type: "EnumType", CustomDiff: canonicalizeMetricDescriptorLabelsValueType, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ValueType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareMetricDescriptorMetadataNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*MetricDescriptorMetadata)
	if !ok {
		desiredNotPointer, ok := d.(MetricDescriptorMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricDescriptorMetadata or *MetricDescriptorMetadata", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*MetricDescriptorMetadata)
	if !ok {
		actualNotPointer, ok := a.(MetricDescriptorMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricDescriptorMetadata", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LaunchStage, actual.LaunchStage, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LaunchStage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SamplePeriod, actual.SamplePeriod, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SamplePeriod")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IngestDelay, actual.IngestDelay, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IngestDelay")); len(ds) != 0 || err != nil {
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
func (r *MetricDescriptor) urlNormalized() *MetricDescriptor {
	normalized := dcl.Copy(*r).(MetricDescriptor)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Type = r.Type
	normalized.Unit = dcl.SelfLinkToName(r.Unit)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *MetricDescriptor) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the MetricDescriptor resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *MetricDescriptor) marshal(c *Client) ([]byte, error) {
	m, err := expandMetricDescriptor(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling MetricDescriptor: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalMetricDescriptor decodes JSON responses into the MetricDescriptor resource schema.
func unmarshalMetricDescriptor(b []byte, c *Client, res *MetricDescriptor) (*MetricDescriptor, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapMetricDescriptor(m, c, res)
}

func unmarshalMapMetricDescriptor(m map[string]interface{}, c *Client, res *MetricDescriptor) (*MetricDescriptor, error) {

	flattened := flattenMetricDescriptor(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandMetricDescriptor expands MetricDescriptor into a JSON request object.
func expandMetricDescriptor(c *Client, f *MetricDescriptor) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Type; dcl.ValueShouldBeSent(v) {
		m["type"] = v
	}
	if v, err := expandMetricDescriptorLabelsSlice(c, f.Labels, res); err != nil {
		return nil, fmt.Errorf("error expanding Labels into labels: %w", err)
	} else if v != nil {
		m["labels"] = v
	}
	if v := f.MetricKind; dcl.ValueShouldBeSent(v) {
		m["metricKind"] = v
	}
	if v := f.ValueType; dcl.ValueShouldBeSent(v) {
		m["valueType"] = v
	}
	if v := f.Unit; dcl.ValueShouldBeSent(v) {
		m["unit"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v, err := expandMetricDescriptorMetadata(c, f.Metadata, res); err != nil {
		return nil, fmt.Errorf("error expanding Metadata into metadata: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}
	if v := f.LaunchStage; dcl.ValueShouldBeSent(v) {
		m["launchStage"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenMetricDescriptor flattens MetricDescriptor from a JSON request object into the
// MetricDescriptor type.
func flattenMetricDescriptor(c *Client, i interface{}, res *MetricDescriptor) *MetricDescriptor {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &MetricDescriptor{}
	resultRes.SelfLink = dcl.FlattenString(m["name"])
	resultRes.Type = dcl.FlattenString(m["type"])
	resultRes.Labels = flattenMetricDescriptorLabelsSlice(c, m["labels"], res)
	resultRes.MetricKind = flattenMetricDescriptorMetricKindEnum(m["metricKind"])
	resultRes.ValueType = flattenMetricDescriptorValueTypeEnum(m["valueType"])
	resultRes.Unit = dcl.FlattenString(m["unit"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.Metadata = flattenMetricDescriptorMetadata(c, m["metadata"], res)
	resultRes.LaunchStage = flattenMetricDescriptorLaunchStageEnum(m["launchStage"])
	resultRes.MonitoredResourceTypes = dcl.FlattenStringSlice(m["monitoredResourceTypes"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandMetricDescriptorLabelsMap expands the contents of MetricDescriptorLabels into a JSON
// request object.
func expandMetricDescriptorLabelsMap(c *Client, f map[string]MetricDescriptorLabels, res *MetricDescriptor) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandMetricDescriptorLabels(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandMetricDescriptorLabelsSlice expands the contents of MetricDescriptorLabels into a JSON
// request object.
func expandMetricDescriptorLabelsSlice(c *Client, f []MetricDescriptorLabels, res *MetricDescriptor) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandMetricDescriptorLabels(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenMetricDescriptorLabelsMap flattens the contents of MetricDescriptorLabels from a JSON
// response object.
func flattenMetricDescriptorLabelsMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorLabels{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorLabels{}
	}

	items := make(map[string]MetricDescriptorLabels)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorLabels(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenMetricDescriptorLabelsSlice flattens the contents of MetricDescriptorLabels from a JSON
// response object.
func flattenMetricDescriptorLabelsSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorLabels{}
	}

	if len(a) == 0 {
		return []MetricDescriptorLabels{}
	}

	items := make([]MetricDescriptorLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorLabels(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandMetricDescriptorLabels expands an instance of MetricDescriptorLabels into a JSON
// request object.
func expandMetricDescriptorLabels(c *Client, f *MetricDescriptorLabels, res *MetricDescriptor) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.ValueType; !dcl.IsEmptyValueIndirect(v) {
		m["valueType"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}

	return m, nil
}

// flattenMetricDescriptorLabels flattens an instance of MetricDescriptorLabels from a JSON
// response object.
func flattenMetricDescriptorLabels(c *Client, i interface{}, res *MetricDescriptor) *MetricDescriptorLabels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &MetricDescriptorLabels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyMetricDescriptorLabels
	}
	r.Key = dcl.FlattenString(m["key"])
	r.ValueType = flattenMetricDescriptorLabelsValueTypeEnum(m["valueType"])
	r.Description = dcl.FlattenString(m["description"])

	return r
}

// expandMetricDescriptorMetadataMap expands the contents of MetricDescriptorMetadata into a JSON
// request object.
func expandMetricDescriptorMetadataMap(c *Client, f map[string]MetricDescriptorMetadata, res *MetricDescriptor) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandMetricDescriptorMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandMetricDescriptorMetadataSlice expands the contents of MetricDescriptorMetadata into a JSON
// request object.
func expandMetricDescriptorMetadataSlice(c *Client, f []MetricDescriptorMetadata, res *MetricDescriptor) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandMetricDescriptorMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenMetricDescriptorMetadataMap flattens the contents of MetricDescriptorMetadata from a JSON
// response object.
func flattenMetricDescriptorMetadataMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorMetadata {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorMetadata{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorMetadata{}
	}

	items := make(map[string]MetricDescriptorMetadata)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorMetadata(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenMetricDescriptorMetadataSlice flattens the contents of MetricDescriptorMetadata from a JSON
// response object.
func flattenMetricDescriptorMetadataSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorMetadata {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorMetadata{}
	}

	if len(a) == 0 {
		return []MetricDescriptorMetadata{}
	}

	items := make([]MetricDescriptorMetadata, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorMetadata(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandMetricDescriptorMetadata expands an instance of MetricDescriptorMetadata into a JSON
// request object.
func expandMetricDescriptorMetadata(c *Client, f *MetricDescriptorMetadata, res *MetricDescriptor) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LaunchStage; !dcl.IsEmptyValueIndirect(v) {
		m["launchStage"] = v
	}
	if v := f.SamplePeriod; !dcl.IsEmptyValueIndirect(v) {
		m["samplePeriod"] = v
	}
	if v := f.IngestDelay; !dcl.IsEmptyValueIndirect(v) {
		m["ingestDelay"] = v
	}

	return m, nil
}

// flattenMetricDescriptorMetadata flattens an instance of MetricDescriptorMetadata from a JSON
// response object.
func flattenMetricDescriptorMetadata(c *Client, i interface{}, res *MetricDescriptor) *MetricDescriptorMetadata {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &MetricDescriptorMetadata{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyMetricDescriptorMetadata
	}
	r.LaunchStage = flattenMetricDescriptorMetadataLaunchStageEnum(m["launchStage"])
	r.SamplePeriod = dcl.FlattenString(m["samplePeriod"])
	r.IngestDelay = dcl.FlattenString(m["ingestDelay"])

	return r
}

// flattenMetricDescriptorLabelsValueTypeEnumMap flattens the contents of MetricDescriptorLabelsValueTypeEnum from a JSON
// response object.
func flattenMetricDescriptorLabelsValueTypeEnumMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorLabelsValueTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorLabelsValueTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorLabelsValueTypeEnum{}
	}

	items := make(map[string]MetricDescriptorLabelsValueTypeEnum)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorLabelsValueTypeEnum(item.(interface{}))
	}

	return items
}

// flattenMetricDescriptorLabelsValueTypeEnumSlice flattens the contents of MetricDescriptorLabelsValueTypeEnum from a JSON
// response object.
func flattenMetricDescriptorLabelsValueTypeEnumSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorLabelsValueTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorLabelsValueTypeEnum{}
	}

	if len(a) == 0 {
		return []MetricDescriptorLabelsValueTypeEnum{}
	}

	items := make([]MetricDescriptorLabelsValueTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorLabelsValueTypeEnum(item.(interface{})))
	}

	return items
}

// flattenMetricDescriptorLabelsValueTypeEnum asserts that an interface is a string, and returns a
// pointer to a *MetricDescriptorLabelsValueTypeEnum with the same value as that string.
func flattenMetricDescriptorLabelsValueTypeEnum(i interface{}) *MetricDescriptorLabelsValueTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return MetricDescriptorLabelsValueTypeEnumRef(s)
}

// flattenMetricDescriptorMetricKindEnumMap flattens the contents of MetricDescriptorMetricKindEnum from a JSON
// response object.
func flattenMetricDescriptorMetricKindEnumMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorMetricKindEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorMetricKindEnum{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorMetricKindEnum{}
	}

	items := make(map[string]MetricDescriptorMetricKindEnum)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorMetricKindEnum(item.(interface{}))
	}

	return items
}

// flattenMetricDescriptorMetricKindEnumSlice flattens the contents of MetricDescriptorMetricKindEnum from a JSON
// response object.
func flattenMetricDescriptorMetricKindEnumSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorMetricKindEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorMetricKindEnum{}
	}

	if len(a) == 0 {
		return []MetricDescriptorMetricKindEnum{}
	}

	items := make([]MetricDescriptorMetricKindEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorMetricKindEnum(item.(interface{})))
	}

	return items
}

// flattenMetricDescriptorMetricKindEnum asserts that an interface is a string, and returns a
// pointer to a *MetricDescriptorMetricKindEnum with the same value as that string.
func flattenMetricDescriptorMetricKindEnum(i interface{}) *MetricDescriptorMetricKindEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return MetricDescriptorMetricKindEnumRef(s)
}

// flattenMetricDescriptorValueTypeEnumMap flattens the contents of MetricDescriptorValueTypeEnum from a JSON
// response object.
func flattenMetricDescriptorValueTypeEnumMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorValueTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorValueTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorValueTypeEnum{}
	}

	items := make(map[string]MetricDescriptorValueTypeEnum)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorValueTypeEnum(item.(interface{}))
	}

	return items
}

// flattenMetricDescriptorValueTypeEnumSlice flattens the contents of MetricDescriptorValueTypeEnum from a JSON
// response object.
func flattenMetricDescriptorValueTypeEnumSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorValueTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorValueTypeEnum{}
	}

	if len(a) == 0 {
		return []MetricDescriptorValueTypeEnum{}
	}

	items := make([]MetricDescriptorValueTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorValueTypeEnum(item.(interface{})))
	}

	return items
}

// flattenMetricDescriptorValueTypeEnum asserts that an interface is a string, and returns a
// pointer to a *MetricDescriptorValueTypeEnum with the same value as that string.
func flattenMetricDescriptorValueTypeEnum(i interface{}) *MetricDescriptorValueTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return MetricDescriptorValueTypeEnumRef(s)
}

// flattenMetricDescriptorMetadataLaunchStageEnumMap flattens the contents of MetricDescriptorMetadataLaunchStageEnum from a JSON
// response object.
func flattenMetricDescriptorMetadataLaunchStageEnumMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorMetadataLaunchStageEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorMetadataLaunchStageEnum{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorMetadataLaunchStageEnum{}
	}

	items := make(map[string]MetricDescriptorMetadataLaunchStageEnum)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorMetadataLaunchStageEnum(item.(interface{}))
	}

	return items
}

// flattenMetricDescriptorMetadataLaunchStageEnumSlice flattens the contents of MetricDescriptorMetadataLaunchStageEnum from a JSON
// response object.
func flattenMetricDescriptorMetadataLaunchStageEnumSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorMetadataLaunchStageEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorMetadataLaunchStageEnum{}
	}

	if len(a) == 0 {
		return []MetricDescriptorMetadataLaunchStageEnum{}
	}

	items := make([]MetricDescriptorMetadataLaunchStageEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorMetadataLaunchStageEnum(item.(interface{})))
	}

	return items
}

// flattenMetricDescriptorMetadataLaunchStageEnum asserts that an interface is a string, and returns a
// pointer to a *MetricDescriptorMetadataLaunchStageEnum with the same value as that string.
func flattenMetricDescriptorMetadataLaunchStageEnum(i interface{}) *MetricDescriptorMetadataLaunchStageEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return MetricDescriptorMetadataLaunchStageEnumRef(s)
}

// flattenMetricDescriptorLaunchStageEnumMap flattens the contents of MetricDescriptorLaunchStageEnum from a JSON
// response object.
func flattenMetricDescriptorLaunchStageEnumMap(c *Client, i interface{}, res *MetricDescriptor) map[string]MetricDescriptorLaunchStageEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricDescriptorLaunchStageEnum{}
	}

	if len(a) == 0 {
		return map[string]MetricDescriptorLaunchStageEnum{}
	}

	items := make(map[string]MetricDescriptorLaunchStageEnum)
	for k, item := range a {
		items[k] = *flattenMetricDescriptorLaunchStageEnum(item.(interface{}))
	}

	return items
}

// flattenMetricDescriptorLaunchStageEnumSlice flattens the contents of MetricDescriptorLaunchStageEnum from a JSON
// response object.
func flattenMetricDescriptorLaunchStageEnumSlice(c *Client, i interface{}, res *MetricDescriptor) []MetricDescriptorLaunchStageEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricDescriptorLaunchStageEnum{}
	}

	if len(a) == 0 {
		return []MetricDescriptorLaunchStageEnum{}
	}

	items := make([]MetricDescriptorLaunchStageEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricDescriptorLaunchStageEnum(item.(interface{})))
	}

	return items
}

// flattenMetricDescriptorLaunchStageEnum asserts that an interface is a string, and returns a
// pointer to a *MetricDescriptorLaunchStageEnum with the same value as that string.
func flattenMetricDescriptorLaunchStageEnum(i interface{}) *MetricDescriptorLaunchStageEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return MetricDescriptorLaunchStageEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *MetricDescriptor) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalMetricDescriptor(b, c, r)
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
		if nr.Type == nil && ncr.Type == nil {
			c.Config.Logger.Info("Both Type fields null - considering equal.")
		} else if nr.Type == nil || ncr.Type == nil {
			c.Config.Logger.Info("Only one Type field is null - considering unequal.")
			return false
		} else if *nr.Type != *ncr.Type {
			return false
		}
		return true
	}
}

type metricDescriptorDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         metricDescriptorApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToMetricDescriptorDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]metricDescriptorDiff, error) {
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
	var diffs []metricDescriptorDiff
	// For each operation name, create a metricDescriptorDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := metricDescriptorDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToMetricDescriptorApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToMetricDescriptorApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (metricDescriptorApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractMetricDescriptorFields(r *MetricDescriptor) error {
	vMetadata := r.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &MetricDescriptorMetadata{}
	}
	if err := extractMetricDescriptorMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		r.Metadata = vMetadata
	}
	return nil
}
func extractMetricDescriptorLabelsFields(r *MetricDescriptor, o *MetricDescriptorLabels) error {
	return nil
}
func extractMetricDescriptorMetadataFields(r *MetricDescriptor, o *MetricDescriptorMetadata) error {
	return nil
}

func postReadExtractMetricDescriptorFields(r *MetricDescriptor) error {
	vMetadata := r.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &MetricDescriptorMetadata{}
	}
	if err := postReadExtractMetricDescriptorMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		r.Metadata = vMetadata
	}
	return nil
}
func postReadExtractMetricDescriptorLabelsFields(r *MetricDescriptor, o *MetricDescriptorLabels) error {
	return nil
}
func postReadExtractMetricDescriptorMetadataFields(r *MetricDescriptor, o *MetricDescriptorMetadata) error {
	return nil
}
