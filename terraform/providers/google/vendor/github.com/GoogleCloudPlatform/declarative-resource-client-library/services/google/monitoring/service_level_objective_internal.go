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

func (r *ServiceLevelObjective) validate() error {

	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"RollingPeriod", "CalendarPeriod"}, r.RollingPeriod, r.CalendarPeriod); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "goal"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Service, "Service"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.ServiceLevelIndicator) {
		if err := r.ServiceLevelIndicator.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicator) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"BasicSli", "RequestBased", "WindowsBased"}, r.BasicSli, r.RequestBased, r.WindowsBased); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.BasicSli) {
		if err := r.BasicSli.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.RequestBased) {
		if err := r.RequestBased.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WindowsBased) {
		if err := r.WindowsBased.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Availability", "Latency", "OperationAvailability", "OperationLatency"}, r.Availability, r.Latency, r.OperationAvailability, r.OperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Availability) {
		if err := r.Availability.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Latency) {
		if err := r.Latency.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.OperationAvailability) {
		if err := r.OperationAvailability.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.OperationLatency) {
		if err := r.OperationLatency.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"GoodTotalRatio", "DistributionCut"}, r.GoodTotalRatio, r.DistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GoodTotalRatio) {
		if err := r.GoodTotalRatio.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.DistributionCut) {
		if err := r.DistributionCut.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Range) {
		if err := r.Range.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"GoodBadMetricFilter", "GoodTotalRatioThreshold", "MetricMeanInRange", "MetricSumInRange"}, r.GoodBadMetricFilter, r.GoodTotalRatioThreshold, r.MetricMeanInRange, r.MetricSumInRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GoodTotalRatioThreshold) {
		if err := r.GoodTotalRatioThreshold.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MetricMeanInRange) {
		if err := r.MetricMeanInRange.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MetricSumInRange) {
		if err := r.MetricSumInRange.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Performance", "BasicSliPerformance"}, r.Performance, r.BasicSliPerformance); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Performance) {
		if err := r.Performance.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.BasicSliPerformance) {
		if err := r.BasicSliPerformance.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"GoodTotalRatio", "DistributionCut"}, r.GoodTotalRatio, r.DistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GoodTotalRatio) {
		if err := r.GoodTotalRatio.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.DistributionCut) {
		if err := r.DistributionCut.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Range) {
		if err := r.Range.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Availability", "Latency", "OperationAvailability", "OperationLatency"}, r.Availability, r.Latency, r.OperationAvailability, r.OperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Availability) {
		if err := r.Availability.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Latency) {
		if err := r.Latency.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.OperationAvailability) {
		if err := r.OperationAvailability.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.OperationLatency) {
		if err := r.OperationLatency.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Range) {
		if err := r.Range.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) validate() error {
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Range) {
		if err := r.Range.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) validate() error {
	return nil
}
func (r *ServiceLevelObjective) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://monitoring.googleapis.com/v3/", params)
}

func (r *ServiceLevelObjective) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"service": dcl.ValueOrEmptyString(nr.Service),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services/{{service}}/serviceLevelObjectives/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *ServiceLevelObjective) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"service": dcl.ValueOrEmptyString(nr.Service),
	}
	return dcl.URL("projects/{{project}}/services/{{service}}/serviceLevelObjectives", nr.basePath(), userBasePath, params), nil

}

func (r *ServiceLevelObjective) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"service": dcl.ValueOrEmptyString(nr.Service),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services/{{service}}/serviceLevelObjectives?serviceLevelObjectiveId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *ServiceLevelObjective) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"service": dcl.ValueOrEmptyString(nr.Service),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/services/{{service}}/serviceLevelObjectives/{{name}}", nr.basePath(), userBasePath, params), nil
}

// serviceLevelObjectiveApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type serviceLevelObjectiveApiOperation interface {
	do(context.Context, *ServiceLevelObjective, *Client) error
}

// newUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest creates a request for an
// ServiceLevelObjective resource's UpdateServiceLevelObjective update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest(ctx context.Context, f *ServiceLevelObjective, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicator(c, f.ServiceLevelIndicator, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceLevelIndicator into serviceLevelIndicator: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["serviceLevelIndicator"] = v
	}
	if v := f.Goal; !dcl.IsEmptyValueIndirect(v) {
		req["goal"] = v
	}
	if v := f.RollingPeriod; !dcl.IsEmptyValueIndirect(v) {
		req["rollingPeriod"] = v
	}
	if v := f.CalendarPeriod; !dcl.IsEmptyValueIndirect(v) {
		req["calendarPeriod"] = v
	}
	if v := f.UserLabels; !dcl.IsEmptyValueIndirect(v) {
		req["userLabels"] = v
	}
	return req, nil
}

// marshalUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest converts the update into
// the final JSON request body.
func marshalUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation) do(ctx context.Context, r *ServiceLevelObjective, c *Client) error {
	_, err := c.GetServiceLevelObjective(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateServiceLevelObjective")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateServiceLevelObjectiveUpdateServiceLevelObjectiveRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listServiceLevelObjectiveRaw(ctx context.Context, r *ServiceLevelObjective, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ServiceLevelObjectiveMaxPage {
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

type listServiceLevelObjectiveOperation struct {
	ServiceLevelObjectives []map[string]interface{} `json:"serviceLevelObjectives"`
	Token                  string                   `json:"nextPageToken"`
}

func (c *Client) listServiceLevelObjective(ctx context.Context, r *ServiceLevelObjective, pageToken string, pageSize int32) ([]*ServiceLevelObjective, string, error) {
	b, err := c.listServiceLevelObjectiveRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listServiceLevelObjectiveOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*ServiceLevelObjective
	for _, v := range m.ServiceLevelObjectives {
		res, err := unmarshalMapServiceLevelObjective(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Service = r.Service
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllServiceLevelObjective(ctx context.Context, f func(*ServiceLevelObjective) bool, resources []*ServiceLevelObjective) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteServiceLevelObjective(ctx, res)
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

type deleteServiceLevelObjectiveOperation struct{}

func (op *deleteServiceLevelObjectiveOperation) do(ctx context.Context, r *ServiceLevelObjective, c *Client) error {
	r, err := c.GetServiceLevelObjective(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "ServiceLevelObjective not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetServiceLevelObjective checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete ServiceLevelObjective: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createServiceLevelObjectiveOperation struct {
	response map[string]interface{}
}

func (op *createServiceLevelObjectiveOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createServiceLevelObjectiveOperation) do(ctx context.Context, r *ServiceLevelObjective, c *Client) error {
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

	if _, err := c.GetServiceLevelObjective(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getServiceLevelObjectiveRaw(ctx context.Context, r *ServiceLevelObjective) ([]byte, error) {

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

func (c *Client) serviceLevelObjectiveDiffsForRawDesired(ctx context.Context, rawDesired *ServiceLevelObjective, opts ...dcl.ApplyOption) (initial, desired *ServiceLevelObjective, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *ServiceLevelObjective
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*ServiceLevelObjective); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected ServiceLevelObjective, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetServiceLevelObjective(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a ServiceLevelObjective resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve ServiceLevelObjective resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that ServiceLevelObjective resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeServiceLevelObjectiveDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for ServiceLevelObjective: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for ServiceLevelObjective: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractServiceLevelObjectiveFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeServiceLevelObjectiveInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for ServiceLevelObjective: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeServiceLevelObjectiveDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for ServiceLevelObjective: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffServiceLevelObjective(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeServiceLevelObjectiveInitialState(rawInitial, rawDesired *ServiceLevelObjective) (*ServiceLevelObjective, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.RollingPeriod) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.CalendarPeriod) {
			rawInitial.RollingPeriod = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.CalendarPeriod) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.RollingPeriod) {
			rawInitial.CalendarPeriod = ServiceLevelObjectiveCalendarPeriodEnumRef("")
		}
	}

	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeServiceLevelObjectiveDesiredState(rawDesired, rawInitial *ServiceLevelObjective, opts ...dcl.ApplyOption) (*ServiceLevelObjective, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.ServiceLevelIndicator = canonicalizeServiceLevelObjectiveServiceLevelIndicator(rawDesired.ServiceLevelIndicator, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &ServiceLevelObjective{}
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
	canonicalDesired.ServiceLevelIndicator = canonicalizeServiceLevelObjectiveServiceLevelIndicator(rawDesired.ServiceLevelIndicator, rawInitial.ServiceLevelIndicator, opts...)
	if dcl.IsZeroValue(rawDesired.Goal) || (dcl.IsEmptyValueIndirect(rawDesired.Goal) && dcl.IsEmptyValueIndirect(rawInitial.Goal)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Goal = rawInitial.Goal
	} else {
		canonicalDesired.Goal = rawDesired.Goal
	}
	if dcl.StringCanonicalize(rawDesired.RollingPeriod, rawInitial.RollingPeriod) {
		canonicalDesired.RollingPeriod = rawInitial.RollingPeriod
	} else {
		canonicalDesired.RollingPeriod = rawDesired.RollingPeriod
	}
	if dcl.IsZeroValue(rawDesired.CalendarPeriod) || (dcl.IsEmptyValueIndirect(rawDesired.CalendarPeriod) && dcl.IsEmptyValueIndirect(rawInitial.CalendarPeriod)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.CalendarPeriod = rawInitial.CalendarPeriod
	} else {
		canonicalDesired.CalendarPeriod = rawDesired.CalendarPeriod
	}
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
	if dcl.NameToSelfLink(rawDesired.Service, rawInitial.Service) {
		canonicalDesired.Service = rawInitial.Service
	} else {
		canonicalDesired.Service = rawDesired.Service
	}

	if canonicalDesired.RollingPeriod != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.CalendarPeriod) {
			canonicalDesired.RollingPeriod = dcl.String("")
		}
	}

	if canonicalDesired.CalendarPeriod != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.RollingPeriod) {
			canonicalDesired.CalendarPeriod = ServiceLevelObjectiveCalendarPeriodEnumRef("")
		}
	}

	return canonicalDesired, nil
}

func canonicalizeServiceLevelObjectiveNewState(c *Client, rawNew, rawDesired *ServiceLevelObjective) (*ServiceLevelObjective, error) {

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

	if dcl.IsEmptyValueIndirect(rawNew.ServiceLevelIndicator) && dcl.IsEmptyValueIndirect(rawDesired.ServiceLevelIndicator) {
		rawNew.ServiceLevelIndicator = rawDesired.ServiceLevelIndicator
	} else {
		rawNew.ServiceLevelIndicator = canonicalizeNewServiceLevelObjectiveServiceLevelIndicator(c, rawDesired.ServiceLevelIndicator, rawNew.ServiceLevelIndicator)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Goal) && dcl.IsEmptyValueIndirect(rawDesired.Goal) {
		rawNew.Goal = rawDesired.Goal
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.RollingPeriod) && dcl.IsEmptyValueIndirect(rawDesired.RollingPeriod) {
		rawNew.RollingPeriod = rawDesired.RollingPeriod
	} else {
		if dcl.StringCanonicalize(rawDesired.RollingPeriod, rawNew.RollingPeriod) {
			rawNew.RollingPeriod = rawDesired.RollingPeriod
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CalendarPeriod) && dcl.IsEmptyValueIndirect(rawDesired.CalendarPeriod) {
		rawNew.CalendarPeriod = rawDesired.CalendarPeriod
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DeleteTime) && dcl.IsEmptyValueIndirect(rawDesired.DeleteTime) {
		rawNew.DeleteTime = rawDesired.DeleteTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceManagementOwned) && dcl.IsEmptyValueIndirect(rawDesired.ServiceManagementOwned) {
		rawNew.ServiceManagementOwned = rawDesired.ServiceManagementOwned
	} else {
		if dcl.BoolCanonicalize(rawDesired.ServiceManagementOwned, rawNew.ServiceManagementOwned) {
			rawNew.ServiceManagementOwned = rawDesired.ServiceManagementOwned
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.UserLabels) && dcl.IsEmptyValueIndirect(rawDesired.UserLabels) {
		rawNew.UserLabels = rawDesired.UserLabels
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Service = rawDesired.Service

	return rawNew, nil
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicator(des, initial *ServiceLevelObjectiveServiceLevelIndicator, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicator {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.BasicSli != nil || (initial != nil && initial.BasicSli != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.RequestBased, des.WindowsBased) {
			des.BasicSli = nil
			if initial != nil {
				initial.BasicSli = nil
			}
		}
	}

	if des.RequestBased != nil || (initial != nil && initial.RequestBased != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BasicSli, des.WindowsBased) {
			des.RequestBased = nil
			if initial != nil {
				initial.RequestBased = nil
			}
		}
	}

	if des.WindowsBased != nil || (initial != nil && initial.WindowsBased != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BasicSli, des.RequestBased) {
			des.WindowsBased = nil
			if initial != nil {
				initial.WindowsBased = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicator{}

	cDes.BasicSli = canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSli(des.BasicSli, initial.BasicSli, opts...)
	cDes.RequestBased = canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBased(des.RequestBased, initial.RequestBased, opts...)
	cDes.WindowsBased = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBased(des.WindowsBased, initial.WindowsBased, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicator, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicator {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicator, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicator(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicator, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicator(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicator(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicator) *ServiceLevelObjectiveServiceLevelIndicator {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicator while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BasicSli = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, des.BasicSli, nw.BasicSli)
	nw.RequestBased = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, des.RequestBased, nw.RequestBased)
	nw.WindowsBased = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, des.WindowsBased, nw.WindowsBased)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicator) []ServiceLevelObjectiveServiceLevelIndicator {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicator
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicator(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicator) []ServiceLevelObjectiveServiceLevelIndicator {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicator
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicator(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSli(des, initial *ServiceLevelObjectiveServiceLevelIndicatorBasicSli, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Availability != nil || (initial != nil && initial.Availability != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Latency, des.OperationAvailability, des.OperationLatency) {
			des.Availability = nil
			if initial != nil {
				initial.Availability = nil
			}
		}
	}

	if des.Latency != nil || (initial != nil && initial.Latency != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.OperationAvailability, des.OperationLatency) {
			des.Latency = nil
			if initial != nil {
				initial.Latency = nil
			}
		}
	}

	if des.OperationAvailability != nil || (initial != nil && initial.OperationAvailability != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.Latency, des.OperationLatency) {
			des.OperationAvailability = nil
			if initial != nil {
				initial.OperationAvailability = nil
			}
		}
	}

	if des.OperationLatency != nil || (initial != nil && initial.OperationLatency != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.Latency, des.OperationAvailability) {
			des.OperationLatency = nil
			if initial != nil {
				initial.OperationLatency = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}

	if dcl.StringArrayCanonicalize(des.Method, initial.Method) {
		cDes.Method = initial.Method
	} else {
		cDes.Method = des.Method
	}
	if dcl.StringArrayCanonicalize(des.Location, initial.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}
	if dcl.StringArrayCanonicalize(des.Version, initial.Version) {
		cDes.Version = initial.Version
	} else {
		cDes.Version = des.Version
	}
	cDes.Availability = canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(des.Availability, initial.Availability, opts...)
	cDes.Latency = canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(des.Latency, initial.Latency, opts...)
	cDes.OperationAvailability = canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(des.OperationAvailability, initial.OperationAvailability, opts...)
	cDes.OperationLatency = canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(des.OperationLatency, initial.OperationLatency, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorBasicSli, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSli, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSli(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSli, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSli(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSli(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) *ServiceLevelObjectiveServiceLevelIndicatorBasicSli {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorBasicSli while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Method, nw.Method) {
		nw.Method = des.Method
	}
	if dcl.StringArrayCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}
	if dcl.StringArrayCanonicalize(des.Version, nw.Version) {
		nw.Version = des.Version
	}
	nw.Availability = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, des.Availability, nw.Availability)
	nw.Latency = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, des.Latency, nw.Latency)
	nw.OperationAvailability = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, des.OperationAvailability, nw.OperationAvailability)
	nw.OperationLatency = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, des.OperationLatency, nw.OperationLatency)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSli) []ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSli
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorBasicSliNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSli) []ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSli
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(des, initial *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(des, initial *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}

	if dcl.StringCanonicalize(des.Threshold, initial.Threshold) || dcl.IsZeroValue(des.Threshold) {
		cDes.Threshold = initial.Threshold
	} else {
		cDes.Threshold = des.Threshold
	}
	if dcl.IsZeroValue(des.Experience) || (dcl.IsEmptyValueIndirect(des.Experience) && dcl.IsEmptyValueIndirect(initial.Experience)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Experience = initial.Experience
	} else {
		cDes.Experience = des.Experience
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Threshold, nw.Threshold) {
		nw.Threshold = des.Threshold
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(des, initial *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(des, initial *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}

	if dcl.StringCanonicalize(des.Threshold, initial.Threshold) || dcl.IsZeroValue(des.Threshold) {
		cDes.Threshold = initial.Threshold
	} else {
		cDes.Threshold = des.Threshold
	}
	if dcl.IsZeroValue(des.Experience) || (dcl.IsEmptyValueIndirect(des.Experience) && dcl.IsEmptyValueIndirect(initial.Experience)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Experience = initial.Experience
	} else {
		cDes.Experience = des.Experience
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Threshold, nw.Threshold) {
		nw.Threshold = des.Threshold
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBased(des, initial *ServiceLevelObjectiveServiceLevelIndicatorRequestBased, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.GoodTotalRatio != nil || (initial != nil && initial.GoodTotalRatio != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.DistributionCut) {
			des.GoodTotalRatio = nil
			if initial != nil {
				initial.GoodTotalRatio = nil
			}
		}
	}

	if des.DistributionCut != nil || (initial != nil && initial.DistributionCut != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodTotalRatio) {
			des.DistributionCut = nil
			if initial != nil {
				initial.DistributionCut = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}

	cDes.GoodTotalRatio = canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(des.GoodTotalRatio, initial.GoodTotalRatio, opts...)
	cDes.DistributionCut = canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(des.DistributionCut, initial.DistributionCut, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorRequestBased, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBased, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBased(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBased, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBased(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBased(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) *ServiceLevelObjectiveServiceLevelIndicatorRequestBased {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorRequestBased while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.GoodTotalRatio = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, des.GoodTotalRatio, nw.GoodTotalRatio)
	nw.DistributionCut = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, des.DistributionCut, nw.DistributionCut)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBased) []ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBased
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBased) []ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBased
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(des, initial *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}

	if dcl.StringCanonicalize(des.GoodServiceFilter, initial.GoodServiceFilter) || dcl.IsZeroValue(des.GoodServiceFilter) {
		cDes.GoodServiceFilter = initial.GoodServiceFilter
	} else {
		cDes.GoodServiceFilter = des.GoodServiceFilter
	}
	if dcl.StringCanonicalize(des.BadServiceFilter, initial.BadServiceFilter) || dcl.IsZeroValue(des.BadServiceFilter) {
		cDes.BadServiceFilter = initial.BadServiceFilter
	} else {
		cDes.BadServiceFilter = des.BadServiceFilter
	}
	if dcl.StringCanonicalize(des.TotalServiceFilter, initial.TotalServiceFilter) || dcl.IsZeroValue(des.TotalServiceFilter) {
		cDes.TotalServiceFilter = initial.TotalServiceFilter
	} else {
		cDes.TotalServiceFilter = des.TotalServiceFilter
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.GoodServiceFilter, nw.GoodServiceFilter) {
		nw.GoodServiceFilter = des.GoodServiceFilter
	}
	if dcl.StringCanonicalize(des.BadServiceFilter, nw.BadServiceFilter) {
		nw.BadServiceFilter = des.BadServiceFilter
	}
	if dcl.StringCanonicalize(des.TotalServiceFilter, nw.TotalServiceFilter) {
		nw.TotalServiceFilter = des.TotalServiceFilter
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(des, initial *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}

	if dcl.StringCanonicalize(des.DistributionFilter, initial.DistributionFilter) || dcl.IsZeroValue(des.DistributionFilter) {
		cDes.DistributionFilter = initial.DistributionFilter
	} else {
		cDes.DistributionFilter = des.DistributionFilter
	}
	cDes.Range = canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(des.Range, initial.Range, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.DistributionFilter, nw.DistributionFilter) {
		nw.DistributionFilter = des.DistributionFilter
	}
	nw.Range = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, des.Range, nw.Range)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}

	if dcl.IsZeroValue(des.Min) || (dcl.IsEmptyValueIndirect(des.Min) && dcl.IsEmptyValueIndirect(initial.Min)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Min = initial.Min
	} else {
		cDes.Min = des.Min
	}
	if dcl.IsZeroValue(des.Max) || (dcl.IsEmptyValueIndirect(des.Max) && dcl.IsEmptyValueIndirect(initial.Max)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Max = initial.Max
	} else {
		cDes.Max = des.Max
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBased(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.GoodBadMetricFilter != nil || (initial != nil && initial.GoodBadMetricFilter != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodTotalRatioThreshold, des.MetricMeanInRange, des.MetricSumInRange) {
			des.GoodBadMetricFilter = nil
			if initial != nil {
				initial.GoodBadMetricFilter = nil
			}
		}
	}

	if des.GoodTotalRatioThreshold != nil || (initial != nil && initial.GoodTotalRatioThreshold != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodBadMetricFilter, des.MetricMeanInRange, des.MetricSumInRange) {
			des.GoodTotalRatioThreshold = nil
			if initial != nil {
				initial.GoodTotalRatioThreshold = nil
			}
		}
	}

	if des.MetricMeanInRange != nil || (initial != nil && initial.MetricMeanInRange != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodBadMetricFilter, des.GoodTotalRatioThreshold, des.MetricSumInRange) {
			des.MetricMeanInRange = nil
			if initial != nil {
				initial.MetricMeanInRange = nil
			}
		}
	}

	if des.MetricSumInRange != nil || (initial != nil && initial.MetricSumInRange != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodBadMetricFilter, des.GoodTotalRatioThreshold, des.MetricMeanInRange) {
			des.MetricSumInRange = nil
			if initial != nil {
				initial.MetricSumInRange = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}

	if dcl.StringCanonicalize(des.GoodBadMetricFilter, initial.GoodBadMetricFilter) || dcl.IsZeroValue(des.GoodBadMetricFilter) {
		cDes.GoodBadMetricFilter = initial.GoodBadMetricFilter
	} else {
		cDes.GoodBadMetricFilter = des.GoodBadMetricFilter
	}
	cDes.GoodTotalRatioThreshold = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(des.GoodTotalRatioThreshold, initial.GoodTotalRatioThreshold, opts...)
	cDes.MetricMeanInRange = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(des.MetricMeanInRange, initial.MetricMeanInRange, opts...)
	cDes.MetricSumInRange = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(des.MetricSumInRange, initial.MetricSumInRange, opts...)
	if dcl.StringCanonicalize(des.WindowPeriod, initial.WindowPeriod) || dcl.IsZeroValue(des.WindowPeriod) {
		cDes.WindowPeriod = initial.WindowPeriod
	} else {
		cDes.WindowPeriod = des.WindowPeriod
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBased(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBased(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBased while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.GoodBadMetricFilter, nw.GoodBadMetricFilter) {
		nw.GoodBadMetricFilter = des.GoodBadMetricFilter
	}
	nw.GoodTotalRatioThreshold = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, des.GoodTotalRatioThreshold, nw.GoodTotalRatioThreshold)
	nw.MetricMeanInRange = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, des.MetricMeanInRange, nw.MetricMeanInRange)
	nw.MetricSumInRange = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, des.MetricSumInRange, nw.MetricSumInRange)
	if dcl.StringCanonicalize(des.WindowPeriod, nw.WindowPeriod) {
		nw.WindowPeriod = des.WindowPeriod
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Performance != nil || (initial != nil && initial.Performance != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.BasicSliPerformance) {
			des.Performance = nil
			if initial != nil {
				initial.Performance = nil
			}
		}
	}

	if des.BasicSliPerformance != nil || (initial != nil && initial.BasicSliPerformance != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Performance) {
			des.BasicSliPerformance = nil
			if initial != nil {
				initial.BasicSliPerformance = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}

	cDes.Performance = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(des.Performance, initial.Performance, opts...)
	cDes.BasicSliPerformance = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(des.BasicSliPerformance, initial.BasicSliPerformance, opts...)
	if dcl.IsZeroValue(des.Threshold) || (dcl.IsEmptyValueIndirect(des.Threshold) && dcl.IsEmptyValueIndirect(initial.Threshold)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Threshold = initial.Threshold
	} else {
		cDes.Threshold = des.Threshold
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Performance = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, des.Performance, nw.Performance)
	nw.BasicSliPerformance = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, des.BasicSliPerformance, nw.BasicSliPerformance)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.GoodTotalRatio != nil || (initial != nil && initial.GoodTotalRatio != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.DistributionCut) {
			des.GoodTotalRatio = nil
			if initial != nil {
				initial.GoodTotalRatio = nil
			}
		}
	}

	if des.DistributionCut != nil || (initial != nil && initial.DistributionCut != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GoodTotalRatio) {
			des.DistributionCut = nil
			if initial != nil {
				initial.DistributionCut = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}

	cDes.GoodTotalRatio = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(des.GoodTotalRatio, initial.GoodTotalRatio, opts...)
	cDes.DistributionCut = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(des.DistributionCut, initial.DistributionCut, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.GoodTotalRatio = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, des.GoodTotalRatio, nw.GoodTotalRatio)
	nw.DistributionCut = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, des.DistributionCut, nw.DistributionCut)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}

	if dcl.StringCanonicalize(des.GoodServiceFilter, initial.GoodServiceFilter) || dcl.IsZeroValue(des.GoodServiceFilter) {
		cDes.GoodServiceFilter = initial.GoodServiceFilter
	} else {
		cDes.GoodServiceFilter = des.GoodServiceFilter
	}
	if dcl.StringCanonicalize(des.BadServiceFilter, initial.BadServiceFilter) || dcl.IsZeroValue(des.BadServiceFilter) {
		cDes.BadServiceFilter = initial.BadServiceFilter
	} else {
		cDes.BadServiceFilter = des.BadServiceFilter
	}
	if dcl.StringCanonicalize(des.TotalServiceFilter, initial.TotalServiceFilter) || dcl.IsZeroValue(des.TotalServiceFilter) {
		cDes.TotalServiceFilter = initial.TotalServiceFilter
	} else {
		cDes.TotalServiceFilter = des.TotalServiceFilter
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.GoodServiceFilter, nw.GoodServiceFilter) {
		nw.GoodServiceFilter = des.GoodServiceFilter
	}
	if dcl.StringCanonicalize(des.BadServiceFilter, nw.BadServiceFilter) {
		nw.BadServiceFilter = des.BadServiceFilter
	}
	if dcl.StringCanonicalize(des.TotalServiceFilter, nw.TotalServiceFilter) {
		nw.TotalServiceFilter = des.TotalServiceFilter
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}

	if dcl.StringCanonicalize(des.DistributionFilter, initial.DistributionFilter) || dcl.IsZeroValue(des.DistributionFilter) {
		cDes.DistributionFilter = initial.DistributionFilter
	} else {
		cDes.DistributionFilter = des.DistributionFilter
	}
	cDes.Range = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(des.Range, initial.Range, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.DistributionFilter, nw.DistributionFilter) {
		nw.DistributionFilter = des.DistributionFilter
	}
	nw.Range = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, des.Range, nw.Range)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}

	if dcl.IsZeroValue(des.Min) || (dcl.IsEmptyValueIndirect(des.Min) && dcl.IsEmptyValueIndirect(initial.Min)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Min = initial.Min
	} else {
		cDes.Min = des.Min
	}
	if dcl.IsZeroValue(des.Max) || (dcl.IsEmptyValueIndirect(des.Max) && dcl.IsEmptyValueIndirect(initial.Max)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Max = initial.Max
	} else {
		cDes.Max = des.Max
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Availability != nil || (initial != nil && initial.Availability != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Latency, des.OperationAvailability, des.OperationLatency) {
			des.Availability = nil
			if initial != nil {
				initial.Availability = nil
			}
		}
	}

	if des.Latency != nil || (initial != nil && initial.Latency != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.OperationAvailability, des.OperationLatency) {
			des.Latency = nil
			if initial != nil {
				initial.Latency = nil
			}
		}
	}

	if des.OperationAvailability != nil || (initial != nil && initial.OperationAvailability != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.Latency, des.OperationLatency) {
			des.OperationAvailability = nil
			if initial != nil {
				initial.OperationAvailability = nil
			}
		}
	}

	if des.OperationLatency != nil || (initial != nil && initial.OperationLatency != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Availability, des.Latency, des.OperationAvailability) {
			des.OperationLatency = nil
			if initial != nil {
				initial.OperationLatency = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}

	if dcl.StringArrayCanonicalize(des.Method, initial.Method) {
		cDes.Method = initial.Method
	} else {
		cDes.Method = des.Method
	}
	if dcl.StringArrayCanonicalize(des.Location, initial.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}
	if dcl.StringArrayCanonicalize(des.Version, initial.Version) {
		cDes.Version = initial.Version
	} else {
		cDes.Version = des.Version
	}
	cDes.Availability = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(des.Availability, initial.Availability, opts...)
	cDes.Latency = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(des.Latency, initial.Latency, opts...)
	cDes.OperationAvailability = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(des.OperationAvailability, initial.OperationAvailability, opts...)
	cDes.OperationLatency = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(des.OperationLatency, initial.OperationLatency, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Method, nw.Method) {
		nw.Method = des.Method
	}
	if dcl.StringArrayCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}
	if dcl.StringArrayCanonicalize(des.Version, nw.Version) {
		nw.Version = des.Version
	}
	nw.Availability = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, des.Availability, nw.Availability)
	nw.Latency = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, des.Latency, nw.Latency)
	nw.OperationAvailability = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, des.OperationAvailability, nw.OperationAvailability)
	nw.OperationLatency = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, des.OperationLatency, nw.OperationLatency)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}

	if dcl.StringCanonicalize(des.Threshold, initial.Threshold) || dcl.IsZeroValue(des.Threshold) {
		cDes.Threshold = initial.Threshold
	} else {
		cDes.Threshold = des.Threshold
	}
	if dcl.IsZeroValue(des.Experience) || (dcl.IsEmptyValueIndirect(des.Experience) && dcl.IsEmptyValueIndirect(initial.Experience)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Experience = initial.Experience
	} else {
		cDes.Experience = des.Experience
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Threshold, nw.Threshold) {
		nw.Threshold = des.Threshold
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}

	if dcl.StringCanonicalize(des.Threshold, initial.Threshold) || dcl.IsZeroValue(des.Threshold) {
		cDes.Threshold = initial.Threshold
	} else {
		cDes.Threshold = des.Threshold
	}
	if dcl.IsZeroValue(des.Experience) || (dcl.IsEmptyValueIndirect(des.Experience) && dcl.IsEmptyValueIndirect(initial.Experience)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Experience = initial.Experience
	} else {
		cDes.Experience = des.Experience
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Threshold, nw.Threshold) {
		nw.Threshold = des.Threshold
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}

	if dcl.StringCanonicalize(des.TimeSeries, initial.TimeSeries) || dcl.IsZeroValue(des.TimeSeries) {
		cDes.TimeSeries = initial.TimeSeries
	} else {
		cDes.TimeSeries = des.TimeSeries
	}
	cDes.Range = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(des.Range, initial.Range, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.TimeSeries, nw.TimeSeries) {
		nw.TimeSeries = des.TimeSeries
	}
	nw.Range = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, des.Range, nw.Range)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}

	if dcl.IsZeroValue(des.Min) || (dcl.IsEmptyValueIndirect(des.Min) && dcl.IsEmptyValueIndirect(initial.Min)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Min = initial.Min
	} else {
		cDes.Min = des.Min
	}
	if dcl.IsZeroValue(des.Max) || (dcl.IsEmptyValueIndirect(des.Max) && dcl.IsEmptyValueIndirect(initial.Max)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Max = initial.Max
	} else {
		cDes.Max = des.Max
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}

	if dcl.StringCanonicalize(des.TimeSeries, initial.TimeSeries) || dcl.IsZeroValue(des.TimeSeries) {
		cDes.TimeSeries = initial.TimeSeries
	} else {
		cDes.TimeSeries = des.TimeSeries
	}
	cDes.Range = canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(des.Range, initial.Range, opts...)

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.TimeSeries, nw.TimeSeries) {
		nw.TimeSeries = des.TimeSeries
	}
	nw.Range = canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, des.Range, nw.Range)

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, &d, &n))
	}

	return items
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(des, initial *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, opts ...dcl.ApplyOption) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}

	if dcl.IsZeroValue(des.Min) || (dcl.IsEmptyValueIndirect(des.Min) && dcl.IsEmptyValueIndirect(initial.Min)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Min = initial.Min
	} else {
		cDes.Min = des.Min
	}
	if dcl.IsZeroValue(des.Max) || (dcl.IsEmptyValueIndirect(des.Max) && dcl.IsEmptyValueIndirect(initial.Max)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Max = initial.Max
	} else {
		cDes.Max = des.Max
	}

	return cDes
}

func canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice(des, initial []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, opts ...dcl.ApplyOption) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, 0, len(des))
		for _, d := range des {
			cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, 0, len(des))
	for i, d := range des {
		cd := canonicalizeServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c *Client, des, nw *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSet(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice(c *Client, des, nw []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, &d, &n))
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
func diffServiceLevelObjective(c *Client, desired, actual *ServiceLevelObjective, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceLevelIndicator, actual.ServiceLevelIndicator, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicator, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("ServiceLevelIndicator")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Goal, actual.Goal, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Goal")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RollingPeriod, actual.RollingPeriod, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("RollingPeriod")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CalendarPeriod, actual.CalendarPeriod, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("CalendarPeriod")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DeleteTime, actual.DeleteTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DeleteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceManagementOwned, actual.ServiceManagementOwned, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceManagementOwned")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UserLabels, actual.UserLabels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("UserLabels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
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
func compareServiceLevelObjectiveServiceLevelIndicatorNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicator)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicator)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicator or *ServiceLevelObjectiveServiceLevelIndicator", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicator)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicator)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicator", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BasicSli, actual.BasicSli, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorBasicSliNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSli, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("BasicSli")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RequestBased, actual.RequestBased, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBased, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("RequestBased")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WindowsBased, actual.WindowsBased, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBased, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("WindowsBased")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorBasicSliNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSli)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorBasicSli)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSli or *ServiceLevelObjectiveServiceLevelIndicatorBasicSli", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSli)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorBasicSli)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSli", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Method, actual.Method, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Method")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Availability, actual.Availability, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Availability")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Latency, actual.Latency, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Latency")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationAvailability, actual.OperationAvailability, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("OperationAvailability")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationLatency, actual.OperationLatency, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("OperationLatency")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency or *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Threshold, actual.Threshold, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Threshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Experience, actual.Experience, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Experience")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency or *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Threshold, actual.Threshold, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Threshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Experience, actual.Experience, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Experience")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBased)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorRequestBased)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBased or *ServiceLevelObjectiveServiceLevelIndicatorRequestBased", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBased)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorRequestBased)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBased", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoodTotalRatio, actual.GoodTotalRatio, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodTotalRatio")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DistributionCut, actual.DistributionCut, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("DistributionCut")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio or *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoodServiceFilter, actual.GoodServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BadServiceFilter, actual.BadServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("BadServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TotalServiceFilter, actual.TotalServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("TotalServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut or *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DistributionFilter, actual.DistributionFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("DistributionFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Range, actual.Range, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Range")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange or *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Min, actual.Min, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Min")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Max, actual.Max, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Max")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBased)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBased)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBased or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBased)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBased)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBased", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoodBadMetricFilter, actual.GoodBadMetricFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodBadMetricFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GoodTotalRatioThreshold, actual.GoodTotalRatioThreshold, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodTotalRatioThreshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricMeanInRange, actual.MetricMeanInRange, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("MetricMeanInRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricSumInRange, actual.MetricSumInRange, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("MetricSumInRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WindowPeriod, actual.WindowPeriod, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("WindowPeriod")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Performance, actual.Performance, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Performance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BasicSliPerformance, actual.BasicSliPerformance, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("BasicSliPerformance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Threshold, actual.Threshold, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Threshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoodTotalRatio, actual.GoodTotalRatio, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodTotalRatio")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DistributionCut, actual.DistributionCut, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("DistributionCut")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoodServiceFilter, actual.GoodServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("GoodServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BadServiceFilter, actual.BadServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("BadServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TotalServiceFilter, actual.TotalServiceFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("TotalServiceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DistributionFilter, actual.DistributionFilter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("DistributionFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Range, actual.Range, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Range")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Min, actual.Min, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Min")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Max, actual.Max, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Max")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Method, actual.Method, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Method")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Availability, actual.Availability, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Availability")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Latency, actual.Latency, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Latency")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationAvailability, actual.OperationAvailability, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("OperationAvailability")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationLatency, actual.OperationLatency, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("OperationLatency")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Threshold, actual.Threshold, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Threshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Experience, actual.Experience, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Experience")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Threshold, actual.Threshold, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Threshold")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Experience, actual.Experience, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Experience")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.TimeSeries, actual.TimeSeries, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("TimeSeries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Range, actual.Range, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Range")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Min, actual.Min, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Min")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Max, actual.Max, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Max")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.TimeSeries, actual.TimeSeries, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("TimeSeries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Range, actual.Range, dcl.DiffInfo{ObjectFunction: compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeNewStyle, EmptyObject: EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Range")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange)
	if !ok {
		desiredNotPointer, ok := d.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange or *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange)
	if !ok {
		actualNotPointer, ok := a.(ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Min, actual.Min, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Min")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Max, actual.Max, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation")}, fn.AddNest("Max")); len(ds) != 0 || err != nil {
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
func (r *ServiceLevelObjective) urlNormalized() *ServiceLevelObjective {
	normalized := dcl.Copy(*r).(ServiceLevelObjective)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.RollingPeriod = dcl.SelfLinkToName(r.RollingPeriod)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Service = dcl.SelfLinkToName(r.Service)
	return &normalized
}

func (r *ServiceLevelObjective) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateServiceLevelObjective" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"service": dcl.ValueOrEmptyString(nr.Service),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/services/{{service}}/serviceLevelObjectives/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the ServiceLevelObjective resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *ServiceLevelObjective) marshal(c *Client) ([]byte, error) {
	m, err := expandServiceLevelObjective(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ServiceLevelObjective: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalServiceLevelObjective decodes JSON responses into the ServiceLevelObjective resource schema.
func unmarshalServiceLevelObjective(b []byte, c *Client, res *ServiceLevelObjective) (*ServiceLevelObjective, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapServiceLevelObjective(m, c, res)
}

func unmarshalMapServiceLevelObjective(m map[string]interface{}, c *Client, res *ServiceLevelObjective) (*ServiceLevelObjective, error) {

	flattened := flattenServiceLevelObjective(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandServiceLevelObjective expands ServiceLevelObjective into a JSON request object.
func expandServiceLevelObjective(c *Client, f *ServiceLevelObjective) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/services/%s/serviceLevelObjectives/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Service), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicator(c, f.ServiceLevelIndicator, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceLevelIndicator into serviceLevelIndicator: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["serviceLevelIndicator"] = v
	}
	if v := f.Goal; dcl.ValueShouldBeSent(v) {
		m["goal"] = v
	}
	if v := f.RollingPeriod; dcl.ValueShouldBeSent(v) {
		m["rollingPeriod"] = v
	}
	if v := f.CalendarPeriod; dcl.ValueShouldBeSent(v) {
		m["calendarPeriod"] = v
	}
	if v := f.UserLabels; dcl.ValueShouldBeSent(v) {
		m["userLabels"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Service into service: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}

	return m, nil
}

// flattenServiceLevelObjective flattens ServiceLevelObjective from a JSON request object into the
// ServiceLevelObjective type.
func flattenServiceLevelObjective(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjective {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &ServiceLevelObjective{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.ServiceLevelIndicator = flattenServiceLevelObjectiveServiceLevelIndicator(c, m["serviceLevelIndicator"], res)
	resultRes.Goal = dcl.FlattenDouble(m["goal"])
	resultRes.RollingPeriod = dcl.FlattenString(m["rollingPeriod"])
	resultRes.CalendarPeriod = flattenServiceLevelObjectiveCalendarPeriodEnum(m["calendarPeriod"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.DeleteTime = dcl.FlattenString(m["deleteTime"])
	resultRes.ServiceManagementOwned = dcl.FlattenBool(m["serviceManagementOwned"])
	resultRes.UserLabels = dcl.FlattenKeyValuePairs(m["userLabels"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Service = dcl.FlattenString(m["service"])

	return resultRes
}

// expandServiceLevelObjectiveServiceLevelIndicatorMap expands the contents of ServiceLevelObjectiveServiceLevelIndicator into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicator, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicator(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicator into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicator, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicator(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicator from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicator {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicator{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicator{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicator)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicator(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicator from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicator {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicator{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicator{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicator, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicator(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicator expands an instance of ServiceLevelObjectiveServiceLevelIndicator into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicator(c *Client, f *ServiceLevelObjectiveServiceLevelIndicator, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, f.BasicSli, res); err != nil {
		return nil, fmt.Errorf("error expanding BasicSli into basicSli: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["basicSli"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, f.RequestBased, res); err != nil {
		return nil, fmt.Errorf("error expanding RequestBased into requestBased: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["requestBased"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, f.WindowsBased, res); err != nil {
		return nil, fmt.Errorf("error expanding WindowsBased into windowsBased: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["windowsBased"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicator flattens an instance of ServiceLevelObjectiveServiceLevelIndicator from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicator(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicator {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicator{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicator
	}
	r.BasicSli = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, m["basicSli"], res)
	r.RequestBased = flattenServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, m["requestBased"], res)
	r.WindowsBased = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, m["windowsBased"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSli into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSli, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSli into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorBasicSli, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSli from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSli)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSli from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSli, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSli(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSli expands an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSli into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSli(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorBasicSli, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Method; v != nil {
		m["method"] = v
	}
	if v := f.Location; v != nil {
		m["location"] = v
	}
	if v := f.Version; v != nil {
		m["version"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, f.Availability, res); err != nil {
		return nil, fmt.Errorf("error expanding Availability into availability: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["availability"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, f.Latency, res); err != nil {
		return nil, fmt.Errorf("error expanding Latency into latency: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["latency"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, f.OperationAvailability, res); err != nil {
		return nil, fmt.Errorf("error expanding OperationAvailability into operationAvailability: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["operationAvailability"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, f.OperationLatency, res); err != nil {
		return nil, fmt.Errorf("error expanding OperationLatency into operationLatency: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["operationLatency"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSli flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSli from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSli(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorBasicSli {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSli
	}
	r.Method = dcl.FlattenStringSlice(m["method"])
	r.Location = dcl.FlattenStringSlice(m["location"])
	r.Version = dcl.FlattenStringSlice(m["version"])
	r.Availability = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, m["availability"], res)
	r.Latency = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, m["latency"], res)
	r.OperationAvailability = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, m["operationAvailability"], res)
	r.OperationLatency = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, m["operationLatency"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilitySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability expands an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability
	}

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency expands an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Threshold; !dcl.IsEmptyValueIndirect(v) {
		m["threshold"] = v
	}
	if v := f.Experience; !dcl.IsEmptyValueIndirect(v) {
		m["experience"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency
	}
	r.Threshold = dcl.FlattenString(m["threshold"])
	r.Experience = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum(m["experience"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilitySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability expands an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability
	}

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency expands an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Threshold; !dcl.IsEmptyValueIndirect(v) {
		m["threshold"] = v
	}
	if v := f.Experience; !dcl.IsEmptyValueIndirect(v) {
		m["experience"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency
	}
	r.Threshold = dcl.FlattenString(m["threshold"])
	r.Experience = flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum(m["experience"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBased, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorRequestBased, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBased)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBased, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBased(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBased expands an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBased(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorRequestBased, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, f.GoodTotalRatio, res); err != nil {
		return nil, fmt.Errorf("error expanding GoodTotalRatio into goodTotalRatio: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goodTotalRatio"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, f.DistributionCut, res); err != nil {
		return nil, fmt.Errorf("error expanding DistributionCut into distributionCut: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["distributionCut"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBased flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBased(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorRequestBased {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBased
	}
	r.GoodTotalRatio = flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, m["goodTotalRatio"], res)
	r.DistributionCut = flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, m["distributionCut"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio expands an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GoodServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["goodServiceFilter"] = v
	}
	if v := f.BadServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["badServiceFilter"] = v
	}
	if v := f.TotalServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["totalServiceFilter"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio
	}
	r.GoodServiceFilter = dcl.FlattenString(m["goodServiceFilter"])
	r.BadServiceFilter = dcl.FlattenString(m["badServiceFilter"])
	r.TotalServiceFilter = dcl.FlattenString(m["totalServiceFilter"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut expands an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DistributionFilter; !dcl.IsEmptyValueIndirect(v) {
		m["distributionFilter"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, f.Range, res); err != nil {
		return nil, fmt.Errorf("error expanding Range into range: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["range"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut
	}
	r.DistributionFilter = dcl.FlattenString(m["distributionFilter"])
	r.Range = flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, m["range"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Min; !dcl.IsEmptyValueIndirect(v) {
		m["min"] = v
	}
	if v := f.Max; !dcl.IsEmptyValueIndirect(v) {
		m["max"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange
	}
	r.Min = dcl.FlattenDouble(m["min"])
	r.Max = dcl.FlattenDouble(m["max"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBased expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GoodBadMetricFilter; !dcl.IsEmptyValueIndirect(v) {
		m["goodBadMetricFilter"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, f.GoodTotalRatioThreshold, res); err != nil {
		return nil, fmt.Errorf("error expanding GoodTotalRatioThreshold into goodTotalRatioThreshold: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goodTotalRatioThreshold"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, f.MetricMeanInRange, res); err != nil {
		return nil, fmt.Errorf("error expanding MetricMeanInRange into metricMeanInRange: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metricMeanInRange"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, f.MetricSumInRange, res); err != nil {
		return nil, fmt.Errorf("error expanding MetricSumInRange into metricSumInRange: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metricSumInRange"] = v
	}
	if v := f.WindowPeriod; !dcl.IsEmptyValueIndirect(v) {
		m["windowPeriod"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBased flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBased from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBased(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBased
	}
	r.GoodBadMetricFilter = dcl.FlattenString(m["goodBadMetricFilter"])
	r.GoodTotalRatioThreshold = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, m["goodTotalRatioThreshold"], res)
	r.MetricMeanInRange = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, m["metricMeanInRange"], res)
	r.MetricSumInRange = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, m["metricSumInRange"], res)
	r.WindowPeriod = dcl.FlattenString(m["windowPeriod"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, f.Performance, res); err != nil {
		return nil, fmt.Errorf("error expanding Performance into performance: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["performance"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, f.BasicSliPerformance, res); err != nil {
		return nil, fmt.Errorf("error expanding BasicSliPerformance into basicSliPerformance: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["basicSliPerformance"] = v
	}
	if v := f.Threshold; !dcl.IsEmptyValueIndirect(v) {
		m["threshold"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold
	}
	r.Performance = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, m["performance"], res)
	r.BasicSliPerformance = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, m["basicSliPerformance"], res)
	r.Threshold = dcl.FlattenDouble(m["threshold"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, f.GoodTotalRatio, res); err != nil {
		return nil, fmt.Errorf("error expanding GoodTotalRatio into goodTotalRatio: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goodTotalRatio"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, f.DistributionCut, res); err != nil {
		return nil, fmt.Errorf("error expanding DistributionCut into distributionCut: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["distributionCut"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance
	}
	r.GoodTotalRatio = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, m["goodTotalRatio"], res)
	r.DistributionCut = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, m["distributionCut"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GoodServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["goodServiceFilter"] = v
	}
	if v := f.BadServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["badServiceFilter"] = v
	}
	if v := f.TotalServiceFilter; !dcl.IsEmptyValueIndirect(v) {
		m["totalServiceFilter"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio
	}
	r.GoodServiceFilter = dcl.FlattenString(m["goodServiceFilter"])
	r.BadServiceFilter = dcl.FlattenString(m["badServiceFilter"])
	r.TotalServiceFilter = dcl.FlattenString(m["totalServiceFilter"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DistributionFilter; !dcl.IsEmptyValueIndirect(v) {
		m["distributionFilter"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, f.Range, res); err != nil {
		return nil, fmt.Errorf("error expanding Range into range: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["range"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut
	}
	r.DistributionFilter = dcl.FlattenString(m["distributionFilter"])
	r.Range = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, m["range"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Min; !dcl.IsEmptyValueIndirect(v) {
		m["min"] = v
	}
	if v := f.Max; !dcl.IsEmptyValueIndirect(v) {
		m["max"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange
	}
	r.Min = dcl.FlattenDouble(m["min"])
	r.Max = dcl.FlattenDouble(m["max"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Method; v != nil {
		m["method"] = v
	}
	if v := f.Location; v != nil {
		m["location"] = v
	}
	if v := f.Version; v != nil {
		m["version"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, f.Availability, res); err != nil {
		return nil, fmt.Errorf("error expanding Availability into availability: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["availability"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, f.Latency, res); err != nil {
		return nil, fmt.Errorf("error expanding Latency into latency: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["latency"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, f.OperationAvailability, res); err != nil {
		return nil, fmt.Errorf("error expanding OperationAvailability into operationAvailability: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["operationAvailability"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, f.OperationLatency, res); err != nil {
		return nil, fmt.Errorf("error expanding OperationLatency into operationLatency: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["operationLatency"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance
	}
	r.Method = dcl.FlattenStringSlice(m["method"])
	r.Location = dcl.FlattenStringSlice(m["location"])
	r.Version = dcl.FlattenStringSlice(m["version"])
	r.Availability = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, m["availability"], res)
	r.Latency = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, m["latency"], res)
	r.OperationAvailability = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, m["operationAvailability"], res)
	r.OperationLatency = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, m["operationLatency"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilitySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability
	}

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Threshold; !dcl.IsEmptyValueIndirect(v) {
		m["threshold"] = v
	}
	if v := f.Experience; !dcl.IsEmptyValueIndirect(v) {
		m["experience"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency
	}
	r.Threshold = dcl.FlattenString(m["threshold"])
	r.Experience = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum(m["experience"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilitySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability
	}

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencySlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Threshold; !dcl.IsEmptyValueIndirect(v) {
		m["threshold"] = v
	}
	if v := f.Experience; !dcl.IsEmptyValueIndirect(v) {
		m["experience"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency
	}
	r.Threshold = dcl.FlattenString(m["threshold"])
	r.Experience = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum(m["experience"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.TimeSeries; !dcl.IsEmptyValueIndirect(v) {
		m["timeSeries"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, f.Range, res); err != nil {
		return nil, fmt.Errorf("error expanding Range into range: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["range"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange
	}
	r.TimeSeries = dcl.FlattenString(m["timeSeries"])
	r.Range = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, m["range"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Min; !dcl.IsEmptyValueIndirect(v) {
		m["min"] = v
	}
	if v := f.Max; !dcl.IsEmptyValueIndirect(v) {
		m["max"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange
	}
	r.Min = dcl.FlattenDouble(m["min"])
	r.Max = dcl.FlattenDouble(m["max"])

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.TimeSeries; !dcl.IsEmptyValueIndirect(v) {
		m["timeSeries"] = v
	}
	if v, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, f.Range, res); err != nil {
		return nil, fmt.Errorf("error expanding Range into range: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["range"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange
	}
	r.TimeSeries = dcl.FlattenString(m["timeSeries"])
	r.Range = flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, m["range"], res)

	return r
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeMap expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeMap(c *Client, f map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice expands the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice(c *Client, f []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, res *ServiceLevelObjective) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange expands an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange into a JSON
// request object.
func expandServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c *Client, f *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange, res *ServiceLevelObjective) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Min; !dcl.IsEmptyValueIndirect(v) {
		m["min"] = v
	}
	if v := f.Max; !dcl.IsEmptyValueIndirect(v) {
		m["max"] = v
	}

	return m, nil
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange flattens an instance of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange(c *Client, i interface{}, res *ServiceLevelObjective) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange
	}
	r.Min = dcl.FlattenDouble(m["min"])
	r.Max = dcl.FlattenDouble(m["max"])

	return r
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum(item.(interface{}))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum(item.(interface{})))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum with the same value as that string.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum(i interface{}) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnumRef(s)
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum(item.(interface{}))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum(item.(interface{})))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum with the same value as that string.
func flattenServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum(i interface{}) *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnumRef(s)
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum(item.(interface{}))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum(item.(interface{})))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum with the same value as that string.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum(i interface{}) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnumRef(s)
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumMap flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum{}
	}

	items := make(map[string]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum(item.(interface{}))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumSlice flattens the contents of ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum from a JSON
// response object.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum{}
	}

	items := make([]ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum(item.(interface{})))
	}

	return items
}

// flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum with the same value as that string.
func flattenServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum(i interface{}) *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnumRef(s)
}

// flattenServiceLevelObjectiveCalendarPeriodEnumMap flattens the contents of ServiceLevelObjectiveCalendarPeriodEnum from a JSON
// response object.
func flattenServiceLevelObjectiveCalendarPeriodEnumMap(c *Client, i interface{}, res *ServiceLevelObjective) map[string]ServiceLevelObjectiveCalendarPeriodEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ServiceLevelObjectiveCalendarPeriodEnum{}
	}

	if len(a) == 0 {
		return map[string]ServiceLevelObjectiveCalendarPeriodEnum{}
	}

	items := make(map[string]ServiceLevelObjectiveCalendarPeriodEnum)
	for k, item := range a {
		items[k] = *flattenServiceLevelObjectiveCalendarPeriodEnum(item.(interface{}))
	}

	return items
}

// flattenServiceLevelObjectiveCalendarPeriodEnumSlice flattens the contents of ServiceLevelObjectiveCalendarPeriodEnum from a JSON
// response object.
func flattenServiceLevelObjectiveCalendarPeriodEnumSlice(c *Client, i interface{}, res *ServiceLevelObjective) []ServiceLevelObjectiveCalendarPeriodEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ServiceLevelObjectiveCalendarPeriodEnum{}
	}

	if len(a) == 0 {
		return []ServiceLevelObjectiveCalendarPeriodEnum{}
	}

	items := make([]ServiceLevelObjectiveCalendarPeriodEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenServiceLevelObjectiveCalendarPeriodEnum(item.(interface{})))
	}

	return items
}

// flattenServiceLevelObjectiveCalendarPeriodEnum asserts that an interface is a string, and returns a
// pointer to a *ServiceLevelObjectiveCalendarPeriodEnum with the same value as that string.
func flattenServiceLevelObjectiveCalendarPeriodEnum(i interface{}) *ServiceLevelObjectiveCalendarPeriodEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ServiceLevelObjectiveCalendarPeriodEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *ServiceLevelObjective) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalServiceLevelObjective(b, c, r)
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
		if nr.Service == nil && ncr.Service == nil {
			c.Config.Logger.Info("Both Service fields null - considering equal.")
		} else if nr.Service == nil || ncr.Service == nil {
			c.Config.Logger.Info("Only one Service field is null - considering unequal.")
			return false
		} else if *nr.Service != *ncr.Service {
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

type serviceLevelObjectiveDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         serviceLevelObjectiveApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToServiceLevelObjectiveDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]serviceLevelObjectiveDiff, error) {
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
	var diffs []serviceLevelObjectiveDiff
	// For each operation name, create a serviceLevelObjectiveDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := serviceLevelObjectiveDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToServiceLevelObjectiveApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToServiceLevelObjectiveApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (serviceLevelObjectiveApiOperation, error) {
	switch opName {

	case "updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation":
		return &updateServiceLevelObjectiveUpdateServiceLevelObjectiveOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractServiceLevelObjectiveFields(r *ServiceLevelObjective) error {
	vServiceLevelIndicator := r.ServiceLevelIndicator
	if vServiceLevelIndicator == nil {
		// note: explicitly not the empty object.
		vServiceLevelIndicator = &ServiceLevelObjectiveServiceLevelIndicator{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorFields(r, vServiceLevelIndicator); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceLevelIndicator) {
		r.ServiceLevelIndicator = vServiceLevelIndicator
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicator) error {
	vBasicSli := o.BasicSli
	if vBasicSli == nil {
		// note: explicitly not the empty object.
		vBasicSli = &ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliFields(r, vBasicSli); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicSli) {
		o.BasicSli = vBasicSli
	}
	vRequestBased := o.RequestBased
	if vRequestBased == nil {
		// note: explicitly not the empty object.
		vRequestBased = &ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedFields(r, vRequestBased); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRequestBased) {
		o.RequestBased = vRequestBased
	}
	vWindowsBased := o.WindowsBased
	if vWindowsBased == nil {
		// note: explicitly not the empty object.
		vWindowsBased = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedFields(r, vWindowsBased); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWindowsBased) {
		o.WindowsBased = vWindowsBased
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorBasicSliFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) error {
	vAvailability := o.Availability
	if vAvailability == nil {
		// note: explicitly not the empty object.
		vAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityFields(r, vAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAvailability) {
		o.Availability = vAvailability
	}
	vLatency := o.Latency
	if vLatency == nil {
		// note: explicitly not the empty object.
		vLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyFields(r, vLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLatency) {
		o.Latency = vLatency
	}
	vOperationAvailability := o.OperationAvailability
	if vOperationAvailability == nil {
		// note: explicitly not the empty object.
		vOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityFields(r, vOperationAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationAvailability) {
		o.OperationAvailability = vOperationAvailability
	}
	vOperationLatency := o.OperationLatency
	if vOperationLatency == nil {
		// note: explicitly not the empty object.
		vOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyFields(r, vOperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationLatency) {
		o.OperationLatency = vOperationLatency
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) error {
	vGoodTotalRatio := o.GoodTotalRatio
	if vGoodTotalRatio == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioFields(r, vGoodTotalRatio); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatio) {
		o.GoodTotalRatio = vGoodTotalRatio
	}
	vDistributionCut := o.DistributionCut
	if vDistributionCut == nil {
		// note: explicitly not the empty object.
		vDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutFields(r, vDistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDistributionCut) {
		o.DistributionCut = vDistributionCut
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) error {
	vGoodTotalRatioThreshold := o.GoodTotalRatioThreshold
	if vGoodTotalRatioThreshold == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatioThreshold = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdFields(r, vGoodTotalRatioThreshold); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatioThreshold) {
		o.GoodTotalRatioThreshold = vGoodTotalRatioThreshold
	}
	vMetricMeanInRange := o.MetricMeanInRange
	if vMetricMeanInRange == nil {
		// note: explicitly not the empty object.
		vMetricMeanInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeFields(r, vMetricMeanInRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricMeanInRange) {
		o.MetricMeanInRange = vMetricMeanInRange
	}
	vMetricSumInRange := o.MetricSumInRange
	if vMetricSumInRange == nil {
		// note: explicitly not the empty object.
		vMetricSumInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeFields(r, vMetricSumInRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricSumInRange) {
		o.MetricSumInRange = vMetricSumInRange
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) error {
	vPerformance := o.Performance
	if vPerformance == nil {
		// note: explicitly not the empty object.
		vPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceFields(r, vPerformance); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPerformance) {
		o.Performance = vPerformance
	}
	vBasicSliPerformance := o.BasicSliPerformance
	if vBasicSliPerformance == nil {
		// note: explicitly not the empty object.
		vBasicSliPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceFields(r, vBasicSliPerformance); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicSliPerformance) {
		o.BasicSliPerformance = vBasicSliPerformance
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) error {
	vGoodTotalRatio := o.GoodTotalRatio
	if vGoodTotalRatio == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioFields(r, vGoodTotalRatio); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatio) {
		o.GoodTotalRatio = vGoodTotalRatio
	}
	vDistributionCut := o.DistributionCut
	if vDistributionCut == nil {
		// note: explicitly not the empty object.
		vDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutFields(r, vDistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDistributionCut) {
		o.DistributionCut = vDistributionCut
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) error {
	vAvailability := o.Availability
	if vAvailability == nil {
		// note: explicitly not the empty object.
		vAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityFields(r, vAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAvailability) {
		o.Availability = vAvailability
	}
	vLatency := o.Latency
	if vLatency == nil {
		// note: explicitly not the empty object.
		vLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyFields(r, vLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLatency) {
		o.Latency = vLatency
	}
	vOperationAvailability := o.OperationAvailability
	if vOperationAvailability == nil {
		// note: explicitly not the empty object.
		vOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityFields(r, vOperationAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationAvailability) {
		o.OperationAvailability = vOperationAvailability
	}
	vOperationLatency := o.OperationLatency
	if vOperationLatency == nil {
		// note: explicitly not the empty object.
		vOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyFields(r, vOperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationLatency) {
		o.OperationLatency = vOperationLatency
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) error {
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) error {
	return nil
}

func postReadExtractServiceLevelObjectiveFields(r *ServiceLevelObjective) error {
	vServiceLevelIndicator := r.ServiceLevelIndicator
	if vServiceLevelIndicator == nil {
		// note: explicitly not the empty object.
		vServiceLevelIndicator = &ServiceLevelObjectiveServiceLevelIndicator{}
	}
	if err := postReadExtractServiceLevelObjectiveServiceLevelIndicatorFields(r, vServiceLevelIndicator); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceLevelIndicator) {
		r.ServiceLevelIndicator = vServiceLevelIndicator
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicator) error {
	vBasicSli := o.BasicSli
	if vBasicSli == nil {
		// note: explicitly not the empty object.
		vBasicSli = &ServiceLevelObjectiveServiceLevelIndicatorBasicSli{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliFields(r, vBasicSli); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicSli) {
		o.BasicSli = vBasicSli
	}
	vRequestBased := o.RequestBased
	if vRequestBased == nil {
		// note: explicitly not the empty object.
		vRequestBased = &ServiceLevelObjectiveServiceLevelIndicatorRequestBased{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedFields(r, vRequestBased); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRequestBased) {
		o.RequestBased = vRequestBased
	}
	vWindowsBased := o.WindowsBased
	if vWindowsBased == nil {
		// note: explicitly not the empty object.
		vWindowsBased = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBased{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedFields(r, vWindowsBased); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWindowsBased) {
		o.WindowsBased = vWindowsBased
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorBasicSliFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSli) error {
	vAvailability := o.Availability
	if vAvailability == nil {
		// note: explicitly not the empty object.
		vAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityFields(r, vAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAvailability) {
		o.Availability = vAvailability
	}
	vLatency := o.Latency
	if vLatency == nil {
		// note: explicitly not the empty object.
		vLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyFields(r, vLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLatency) {
		o.Latency = vLatency
	}
	vOperationAvailability := o.OperationAvailability
	if vOperationAvailability == nil {
		// note: explicitly not the empty object.
		vOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityFields(r, vOperationAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationAvailability) {
		o.OperationAvailability = vOperationAvailability
	}
	vOperationLatency := o.OperationLatency
	if vOperationLatency == nil {
		// note: explicitly not the empty object.
		vOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyFields(r, vOperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationLatency) {
		o.OperationLatency = vOperationLatency
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorRequestBasedFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBased) error {
	vGoodTotalRatio := o.GoodTotalRatio
	if vGoodTotalRatio == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioFields(r, vGoodTotalRatio); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatio) {
		o.GoodTotalRatio = vGoodTotalRatio
	}
	vDistributionCut := o.DistributionCut
	if vDistributionCut == nil {
		// note: explicitly not the empty object.
		vDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutFields(r, vDistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDistributionCut) {
		o.DistributionCut = vDistributionCut
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatioFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBased) error {
	vGoodTotalRatioThreshold := o.GoodTotalRatioThreshold
	if vGoodTotalRatioThreshold == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatioThreshold = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdFields(r, vGoodTotalRatioThreshold); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatioThreshold) {
		o.GoodTotalRatioThreshold = vGoodTotalRatioThreshold
	}
	vMetricMeanInRange := o.MetricMeanInRange
	if vMetricMeanInRange == nil {
		// note: explicitly not the empty object.
		vMetricMeanInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeFields(r, vMetricMeanInRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricMeanInRange) {
		o.MetricMeanInRange = vMetricMeanInRange
	}
	vMetricSumInRange := o.MetricSumInRange
	if vMetricSumInRange == nil {
		// note: explicitly not the empty object.
		vMetricSumInRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeFields(r, vMetricSumInRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricSumInRange) {
		o.MetricSumInRange = vMetricSumInRange
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold) error {
	vPerformance := o.Performance
	if vPerformance == nil {
		// note: explicitly not the empty object.
		vPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceFields(r, vPerformance); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPerformance) {
		o.Performance = vPerformance
	}
	vBasicSliPerformance := o.BasicSliPerformance
	if vBasicSliPerformance == nil {
		// note: explicitly not the empty object.
		vBasicSliPerformance = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceFields(r, vBasicSliPerformance); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicSliPerformance) {
		o.BasicSliPerformance = vBasicSliPerformance
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance) error {
	vGoodTotalRatio := o.GoodTotalRatio
	if vGoodTotalRatio == nil {
		// note: explicitly not the empty object.
		vGoodTotalRatio = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioFields(r, vGoodTotalRatio); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGoodTotalRatio) {
		o.GoodTotalRatio = vGoodTotalRatio
	}
	vDistributionCut := o.DistributionCut
	if vDistributionCut == nil {
		// note: explicitly not the empty object.
		vDistributionCut = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutFields(r, vDistributionCut); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDistributionCut) {
		o.DistributionCut = vDistributionCut
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatioFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance) error {
	vAvailability := o.Availability
	if vAvailability == nil {
		// note: explicitly not the empty object.
		vAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityFields(r, vAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAvailability) {
		o.Availability = vAvailability
	}
	vLatency := o.Latency
	if vLatency == nil {
		// note: explicitly not the empty object.
		vLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyFields(r, vLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLatency) {
		o.Latency = vLatency
	}
	vOperationAvailability := o.OperationAvailability
	if vOperationAvailability == nil {
		// note: explicitly not the empty object.
		vOperationAvailability = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityFields(r, vOperationAvailability); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationAvailability) {
		o.OperationAvailability = vOperationAvailability
	}
	vOperationLatency := o.OperationLatency
	if vOperationLatency == nil {
		// note: explicitly not the empty object.
		vOperationLatency = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyFields(r, vOperationLatency); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vOperationLatency) {
		o.OperationLatency = vOperationLatency
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailabilityFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange) error {
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange) error {
	vRange := o.Range
	if vRange == nil {
		// note: explicitly not the empty object.
		vRange = &ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange{}
	}
	if err := extractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeFields(r, vRange); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRange) {
		o.Range = vRange
	}
	return nil
}
func postReadExtractServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRangeFields(r *ServiceLevelObjective, o *ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange) error {
	return nil
}
