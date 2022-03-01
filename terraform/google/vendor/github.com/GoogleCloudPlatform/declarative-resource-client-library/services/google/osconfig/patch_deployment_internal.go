// Copyright 2021 Google LLC. All Rights Reserved.
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
package osconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *PatchDeployment) validate() error {

	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"OneTimeSchedule", "RecurringSchedule"}, r.OneTimeSchedule, r.RecurringSchedule); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "instanceFilter"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.InstanceFilter) {
		if err := r.InstanceFilter.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PatchConfig) {
		if err := r.PatchConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.OneTimeSchedule) {
		if err := r.OneTimeSchedule.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.RecurringSchedule) {
		if err := r.RecurringSchedule.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Rollout) {
		if err := r.Rollout.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentInstanceFilter) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"All", "GroupLabels"}, r.All, r.GroupLabels); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"All", "Zones"}, r.All, r.Zones); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"All", "Instances"}, r.All, r.Instances); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"All", "InstanceNamePrefixes"}, r.All, r.InstanceNamePrefixes); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentInstanceFilterGroupLabels) validate() error {
	return nil
}
func (r *PatchDeploymentPatchConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Apt) {
		if err := r.Apt.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Yum) {
		if err := r.Yum.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Goo) {
		if err := r.Goo.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Zypper) {
		if err := r.Zypper.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WindowsUpdate) {
		if err := r.WindowsUpdate.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PreStep) {
		if err := r.PreStep.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PostStep) {
		if err := r.PostStep.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigApt) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Excludes", "ExclusivePackages"}, r.Excludes, r.ExclusivePackages); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigYum) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Excludes", "ExclusivePackages"}, r.Excludes, r.ExclusivePackages); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigGoo) validate() error {
	return nil
}
func (r *PatchDeploymentPatchConfigZypper) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"WithOptional", "ExclusivePatches"}, r.WithOptional, r.ExclusivePatches); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"WithUpdate", "ExclusivePatches"}, r.WithUpdate, r.ExclusivePatches); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Categories", "ExclusivePatches"}, r.Categories, r.ExclusivePatches); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Severities", "ExclusivePatches"}, r.Severities, r.ExclusivePatches); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Excludes", "ExclusivePatches"}, r.Excludes, r.ExclusivePatches); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigWindowsUpdate) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Classifications", "ExclusivePatches"}, r.Classifications, r.ExclusivePatches); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Excludes", "ExclusivePatches"}, r.Excludes, r.ExclusivePatches); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPreStep) validate() error {
	if !dcl.IsEmptyValueIndirect(r.LinuxExecStepConfig) {
		if err := r.LinuxExecStepConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WindowsExecStepConfig) {
		if err := r.WindowsExecStepConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"LocalPath", "GcsObject"}, r.LocalPath, r.GcsObject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GcsObject) {
		if err := r.GcsObject.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	if err := dcl.Required(r, "generationNumber"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) validate() error {
	if err := dcl.Required(r, "interpreter"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"LocalPath", "GcsObject"}, r.LocalPath, r.GcsObject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GcsObject) {
		if err := r.GcsObject.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	if err := dcl.Required(r, "generationNumber"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPostStep) validate() error {
	if !dcl.IsEmptyValueIndirect(r.LinuxExecStepConfig) {
		if err := r.LinuxExecStepConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WindowsExecStepConfig) {
		if err := r.WindowsExecStepConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"LocalPath", "GcsObject"}, r.LocalPath, r.GcsObject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GcsObject) {
		if err := r.GcsObject.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	if err := dcl.Required(r, "generationNumber"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) validate() error {
	if err := dcl.Required(r, "interpreter"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"LocalPath", "GcsObject"}, r.LocalPath, r.GcsObject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GcsObject) {
		if err := r.GcsObject.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	if err := dcl.Required(r, "generationNumber"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentOneTimeSchedule) validate() error {
	if err := dcl.Required(r, "executeTime"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentRecurringSchedule) validate() error {
	if err := dcl.Required(r, "timeZone"); err != nil {
		return err
	}
	if err := dcl.Required(r, "timeOfDay"); err != nil {
		return err
	}
	if err := dcl.Required(r, "frequency"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Weekly", "Monthly"}, r.Weekly, r.Monthly); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Weekly", "Monthly"}, r.Weekly, r.Monthly); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.TimeZone) {
		if err := r.TimeZone.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.TimeOfDay) {
		if err := r.TimeOfDay.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Weekly) {
		if err := r.Weekly.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Monthly) {
		if err := r.Monthly.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentRecurringScheduleTimeZone) validate() error {
	return nil
}
func (r *PatchDeploymentRecurringScheduleTimeOfDay) validate() error {
	return nil
}
func (r *PatchDeploymentRecurringScheduleWeekly) validate() error {
	if err := dcl.Required(r, "dayOfWeek"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentRecurringScheduleMonthly) validate() error {
	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"WeekDayOfMonth", "MonthDay"}, r.WeekDayOfMonth, r.MonthDay); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.WeekDayOfMonth) {
		if err := r.WeekDayOfMonth.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) validate() error {
	if err := dcl.Required(r, "weekOrdinal"); err != nil {
		return err
	}
	if err := dcl.Required(r, "dayOfWeek"); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeploymentRollout) validate() error {
	if err := dcl.Required(r, "mode"); err != nil {
		return err
	}
	if err := dcl.Required(r, "disruptionBudget"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.DisruptionBudget) {
		if err := r.DisruptionBudget.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PatchDeploymentRolloutDisruptionBudget) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Fixed", "Percent"}, r.Fixed, r.Percent); err != nil {
		return err
	}
	return nil
}
func (r *PatchDeployment) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://osconfig.googleapis.com/v1", params)
}

func (r *PatchDeployment) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/patchDeployments/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *PatchDeployment) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/patchDeployments", nr.basePath(), userBasePath, params), nil

}

func (r *PatchDeployment) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/patchDeployments?patchDeploymentId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *PatchDeployment) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/patchDeployments/{{name}}", nr.basePath(), userBasePath, params), nil
}

// patchDeploymentApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type patchDeploymentApiOperation interface {
	do(context.Context, *PatchDeployment, *Client) error
}

// newUpdatePatchDeploymentUpdatePatchDeploymentRequest creates a request for an
// PatchDeployment resource's UpdatePatchDeployment update type by filling in the update
// fields based on the intended state of the resource.
func newUpdatePatchDeploymentUpdatePatchDeploymentRequest(ctx context.Context, f *PatchDeployment, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandPatchDeploymentInstanceFilter(c, f.InstanceFilter); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["instanceFilter"] = v
	}
	if v, err := expandPatchDeploymentPatchConfig(c, f.PatchConfig); err != nil {
		return nil, fmt.Errorf("error expanding PatchConfig into patchConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["patchConfig"] = v
	}
	if v := f.Duration; !dcl.IsEmptyValueIndirect(v) {
		req["duration"] = v
	}
	if v, err := expandPatchDeploymentOneTimeSchedule(c, f.OneTimeSchedule); err != nil {
		return nil, fmt.Errorf("error expanding OneTimeSchedule into oneTimeSchedule: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["oneTimeSchedule"] = v
	}
	if v, err := expandPatchDeploymentRecurringSchedule(c, f.RecurringSchedule); err != nil {
		return nil, fmt.Errorf("error expanding RecurringSchedule into recurringSchedule: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["recurringSchedule"] = v
	}
	if v, err := expandPatchDeploymentRollout(c, f.Rollout); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["rollout"] = v
	}
	return req, nil
}

// marshalUpdatePatchDeploymentUpdatePatchDeploymentRequest converts the update into
// the final JSON request body.
func marshalUpdatePatchDeploymentUpdatePatchDeploymentRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updatePatchDeploymentUpdatePatchDeploymentOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updatePatchDeploymentUpdatePatchDeploymentOperation) do(ctx context.Context, r *PatchDeployment, c *Client) error {
	_, err := c.GetPatchDeployment(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdatePatchDeployment")
	if err != nil {
		return err
	}
	mask := dcl.TopLevelUpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdatePatchDeploymentUpdatePatchDeploymentRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdatePatchDeploymentUpdatePatchDeploymentRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listPatchDeploymentRaw(ctx context.Context, r *PatchDeployment, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != PatchDeploymentMaxPage {
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

type listPatchDeploymentOperation struct {
	PatchDeployments []map[string]interface{} `json:"patchDeployments"`
	Token            string                   `json:"nextPageToken"`
}

func (c *Client) listPatchDeployment(ctx context.Context, r *PatchDeployment, pageToken string, pageSize int32) ([]*PatchDeployment, string, error) {
	b, err := c.listPatchDeploymentRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listPatchDeploymentOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*PatchDeployment
	for _, v := range m.PatchDeployments {
		res, err := unmarshalMapPatchDeployment(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllPatchDeployment(ctx context.Context, f func(*PatchDeployment) bool, resources []*PatchDeployment) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeletePatchDeployment(ctx, res)
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

type deletePatchDeploymentOperation struct{}

func (op *deletePatchDeploymentOperation) do(ctx context.Context, r *PatchDeployment, c *Client) error {
	r, err := c.GetPatchDeployment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "PatchDeployment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetPatchDeployment checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete PatchDeployment: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createPatchDeploymentOperation struct {
	response map[string]interface{}
}

func (op *createPatchDeploymentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createPatchDeploymentOperation) do(ctx context.Context, r *PatchDeployment, c *Client) error {
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

	if _, err := c.GetPatchDeployment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getPatchDeploymentRaw(ctx context.Context, r *PatchDeployment) ([]byte, error) {

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

func (c *Client) patchDeploymentDiffsForRawDesired(ctx context.Context, rawDesired *PatchDeployment, opts ...dcl.ApplyOption) (initial, desired *PatchDeployment, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *PatchDeployment
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*PatchDeployment); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected PatchDeployment, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetPatchDeployment(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a PatchDeployment resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve PatchDeployment resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that PatchDeployment resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizePatchDeploymentDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for PatchDeployment: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for PatchDeployment: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizePatchDeploymentInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for PatchDeployment: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizePatchDeploymentDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for PatchDeployment: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffPatchDeployment(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizePatchDeploymentInitialState(rawInitial, rawDesired *PatchDeployment) (*PatchDeployment, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.OneTimeSchedule) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.RecurringSchedule) {
			rawInitial.OneTimeSchedule = EmptyPatchDeploymentOneTimeSchedule
		}
	}

	if !dcl.IsZeroValue(rawInitial.RecurringSchedule) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.OneTimeSchedule) {
			rawInitial.RecurringSchedule = EmptyPatchDeploymentRecurringSchedule
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

func canonicalizePatchDeploymentDesiredState(rawDesired, rawInitial *PatchDeployment, opts ...dcl.ApplyOption) (*PatchDeployment, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.InstanceFilter = canonicalizePatchDeploymentInstanceFilter(rawDesired.InstanceFilter, nil, opts...)
		rawDesired.PatchConfig = canonicalizePatchDeploymentPatchConfig(rawDesired.PatchConfig, nil, opts...)
		rawDesired.OneTimeSchedule = canonicalizePatchDeploymentOneTimeSchedule(rawDesired.OneTimeSchedule, nil, opts...)
		rawDesired.RecurringSchedule = canonicalizePatchDeploymentRecurringSchedule(rawDesired.RecurringSchedule, nil, opts...)
		rawDesired.Rollout = canonicalizePatchDeploymentRollout(rawDesired.Rollout, nil, opts...)

		return rawDesired, nil
	}

	if rawDesired.OneTimeSchedule != nil || rawInitial.OneTimeSchedule != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.RecurringSchedule) {
			rawDesired.OneTimeSchedule = nil
			rawInitial.OneTimeSchedule = nil
		}
	}

	if rawDesired.RecurringSchedule != nil || rawInitial.RecurringSchedule != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.OneTimeSchedule) {
			rawDesired.RecurringSchedule = nil
			rawInitial.RecurringSchedule = nil
		}
	}

	canonicalDesired := &PatchDeployment{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	canonicalDesired.InstanceFilter = canonicalizePatchDeploymentInstanceFilter(rawDesired.InstanceFilter, rawInitial.InstanceFilter, opts...)
	canonicalDesired.PatchConfig = canonicalizePatchDeploymentPatchConfig(rawDesired.PatchConfig, rawInitial.PatchConfig, opts...)
	if dcl.StringCanonicalize(rawDesired.Duration, rawInitial.Duration) {
		canonicalDesired.Duration = rawInitial.Duration
	} else {
		canonicalDesired.Duration = rawDesired.Duration
	}
	canonicalDesired.OneTimeSchedule = canonicalizePatchDeploymentOneTimeSchedule(rawDesired.OneTimeSchedule, rawInitial.OneTimeSchedule, opts...)
	canonicalDesired.RecurringSchedule = canonicalizePatchDeploymentRecurringSchedule(rawDesired.RecurringSchedule, rawInitial.RecurringSchedule, opts...)
	canonicalDesired.Rollout = canonicalizePatchDeploymentRollout(rawDesired.Rollout, rawInitial.Rollout, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}

	return canonicalDesired, nil
}

func canonicalizePatchDeploymentNewState(c *Client, rawNew, rawDesired *PatchDeployment) (*PatchDeployment, error) {

	if dcl.IsNotReturnedByServer(rawNew.Name) && dcl.IsNotReturnedByServer(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Description) && dcl.IsNotReturnedByServer(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.InstanceFilter) && dcl.IsNotReturnedByServer(rawDesired.InstanceFilter) {
		rawNew.InstanceFilter = rawDesired.InstanceFilter
	} else {
		rawNew.InstanceFilter = canonicalizeNewPatchDeploymentInstanceFilter(c, rawDesired.InstanceFilter, rawNew.InstanceFilter)
	}

	if dcl.IsNotReturnedByServer(rawNew.PatchConfig) && dcl.IsNotReturnedByServer(rawDesired.PatchConfig) {
		rawNew.PatchConfig = rawDesired.PatchConfig
	} else {
		rawNew.PatchConfig = canonicalizeNewPatchDeploymentPatchConfig(c, rawDesired.PatchConfig, rawNew.PatchConfig)
	}

	if dcl.IsNotReturnedByServer(rawNew.Duration) && dcl.IsNotReturnedByServer(rawDesired.Duration) {
		rawNew.Duration = rawDesired.Duration
	} else {
		if dcl.StringCanonicalize(rawDesired.Duration, rawNew.Duration) {
			rawNew.Duration = rawDesired.Duration
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.OneTimeSchedule) && dcl.IsNotReturnedByServer(rawDesired.OneTimeSchedule) {
		rawNew.OneTimeSchedule = rawDesired.OneTimeSchedule
	} else {
		rawNew.OneTimeSchedule = canonicalizeNewPatchDeploymentOneTimeSchedule(c, rawDesired.OneTimeSchedule, rawNew.OneTimeSchedule)
	}

	if dcl.IsNotReturnedByServer(rawNew.RecurringSchedule) && dcl.IsNotReturnedByServer(rawDesired.RecurringSchedule) {
		rawNew.RecurringSchedule = rawDesired.RecurringSchedule
	} else {
		rawNew.RecurringSchedule = canonicalizeNewPatchDeploymentRecurringSchedule(c, rawDesired.RecurringSchedule, rawNew.RecurringSchedule)
	}

	if dcl.IsNotReturnedByServer(rawNew.CreateTime) && dcl.IsNotReturnedByServer(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.UpdateTime) && dcl.IsNotReturnedByServer(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.LastExecuteTime) && dcl.IsNotReturnedByServer(rawDesired.LastExecuteTime) {
		rawNew.LastExecuteTime = rawDesired.LastExecuteTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Rollout) && dcl.IsNotReturnedByServer(rawDesired.Rollout) {
		rawNew.Rollout = rawDesired.Rollout
	} else {
		rawNew.Rollout = canonicalizeNewPatchDeploymentRollout(c, rawDesired.Rollout, rawNew.Rollout)
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizePatchDeploymentInstanceFilter(des, initial *PatchDeploymentInstanceFilter, opts ...dcl.ApplyOption) *PatchDeploymentInstanceFilter {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.All != nil || (initial != nil && initial.All != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GroupLabels) {
			des.All = nil
			if initial != nil {
				initial.All = nil
			}
		}
	}

	if des.GroupLabels != nil || (initial != nil && initial.GroupLabels != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.All) {
			des.GroupLabels = nil
			if initial != nil {
				initial.GroupLabels = nil
			}
		}
	}

	if des.All != nil || (initial != nil && initial.All != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Zones) {
			des.All = nil
			if initial != nil {
				initial.All = nil
			}
		}
	}

	if des.Zones != nil || (initial != nil && initial.Zones != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.All) {
			des.Zones = nil
			if initial != nil {
				initial.Zones = nil
			}
		}
	}

	if des.All != nil || (initial != nil && initial.All != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Instances) {
			des.All = nil
			if initial != nil {
				initial.All = nil
			}
		}
	}

	if des.Instances != nil || (initial != nil && initial.Instances != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.All) {
			des.Instances = nil
			if initial != nil {
				initial.Instances = nil
			}
		}
	}

	if des.All != nil || (initial != nil && initial.All != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.InstanceNamePrefixes) {
			des.All = nil
			if initial != nil {
				initial.All = nil
			}
		}
	}

	if des.InstanceNamePrefixes != nil || (initial != nil && initial.InstanceNamePrefixes != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.All) {
			des.InstanceNamePrefixes = nil
			if initial != nil {
				initial.InstanceNamePrefixes = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentInstanceFilter{}

	if dcl.BoolCanonicalize(des.All, initial.All) || dcl.IsZeroValue(des.All) {
		cDes.All = initial.All
	} else {
		cDes.All = des.All
	}
	cDes.GroupLabels = canonicalizePatchDeploymentInstanceFilterGroupLabelsSlice(des.GroupLabels, initial.GroupLabels, opts...)
	if dcl.StringArrayCanonicalize(des.Zones, initial.Zones) || dcl.IsZeroValue(des.Zones) {
		cDes.Zones = initial.Zones
	} else {
		cDes.Zones = des.Zones
	}
	if dcl.StringArrayCanonicalize(des.Instances, initial.Instances) || dcl.IsZeroValue(des.Instances) {
		cDes.Instances = initial.Instances
	} else {
		cDes.Instances = des.Instances
	}
	if dcl.StringArrayCanonicalize(des.InstanceNamePrefixes, initial.InstanceNamePrefixes) || dcl.IsZeroValue(des.InstanceNamePrefixes) {
		cDes.InstanceNamePrefixes = initial.InstanceNamePrefixes
	} else {
		cDes.InstanceNamePrefixes = des.InstanceNamePrefixes
	}

	return cDes
}

func canonicalizePatchDeploymentInstanceFilterSlice(des, initial []PatchDeploymentInstanceFilter, opts ...dcl.ApplyOption) []PatchDeploymentInstanceFilter {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentInstanceFilter, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentInstanceFilter(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentInstanceFilter, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentInstanceFilter(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentInstanceFilter(c *Client, des, nw *PatchDeploymentInstanceFilter) *PatchDeploymentInstanceFilter {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentInstanceFilter while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.All, nw.All) {
		nw.All = des.All
	}
	nw.GroupLabels = canonicalizeNewPatchDeploymentInstanceFilterGroupLabelsSlice(c, des.GroupLabels, nw.GroupLabels)
	if dcl.StringArrayCanonicalize(des.Zones, nw.Zones) {
		nw.Zones = des.Zones
	}
	if dcl.StringArrayCanonicalize(des.Instances, nw.Instances) {
		nw.Instances = des.Instances
	}
	if dcl.StringArrayCanonicalize(des.InstanceNamePrefixes, nw.InstanceNamePrefixes) {
		nw.InstanceNamePrefixes = des.InstanceNamePrefixes
	}

	return nw
}

func canonicalizeNewPatchDeploymentInstanceFilterSet(c *Client, des, nw []PatchDeploymentInstanceFilter) []PatchDeploymentInstanceFilter {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentInstanceFilter
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentInstanceFilterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentInstanceFilterSlice(c *Client, des, nw []PatchDeploymentInstanceFilter) []PatchDeploymentInstanceFilter {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentInstanceFilter
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentInstanceFilter(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentInstanceFilterGroupLabels(des, initial *PatchDeploymentInstanceFilterGroupLabels, opts ...dcl.ApplyOption) *PatchDeploymentInstanceFilterGroupLabels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentInstanceFilterGroupLabels{}

	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}

	return cDes
}

func canonicalizePatchDeploymentInstanceFilterGroupLabelsSlice(des, initial []PatchDeploymentInstanceFilterGroupLabels, opts ...dcl.ApplyOption) []PatchDeploymentInstanceFilterGroupLabels {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentInstanceFilterGroupLabels, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentInstanceFilterGroupLabels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentInstanceFilterGroupLabels, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentInstanceFilterGroupLabels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentInstanceFilterGroupLabels(c *Client, des, nw *PatchDeploymentInstanceFilterGroupLabels) *PatchDeploymentInstanceFilterGroupLabels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentInstanceFilterGroupLabels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentInstanceFilterGroupLabelsSet(c *Client, des, nw []PatchDeploymentInstanceFilterGroupLabels) []PatchDeploymentInstanceFilterGroupLabels {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentInstanceFilterGroupLabels
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentInstanceFilterGroupLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentInstanceFilterGroupLabelsSlice(c *Client, des, nw []PatchDeploymentInstanceFilterGroupLabels) []PatchDeploymentInstanceFilterGroupLabels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentInstanceFilterGroupLabels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentInstanceFilterGroupLabels(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfig(des, initial *PatchDeploymentPatchConfig, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfig{}

	if dcl.IsZeroValue(des.RebootConfig) {
		cDes.RebootConfig = initial.RebootConfig
	} else {
		cDes.RebootConfig = des.RebootConfig
	}
	cDes.Apt = canonicalizePatchDeploymentPatchConfigApt(des.Apt, initial.Apt, opts...)
	cDes.Yum = canonicalizePatchDeploymentPatchConfigYum(des.Yum, initial.Yum, opts...)
	cDes.Goo = canonicalizePatchDeploymentPatchConfigGoo(des.Goo, initial.Goo, opts...)
	cDes.Zypper = canonicalizePatchDeploymentPatchConfigZypper(des.Zypper, initial.Zypper, opts...)
	cDes.WindowsUpdate = canonicalizePatchDeploymentPatchConfigWindowsUpdate(des.WindowsUpdate, initial.WindowsUpdate, opts...)
	cDes.PreStep = canonicalizePatchDeploymentPatchConfigPreStep(des.PreStep, initial.PreStep, opts...)
	cDes.PostStep = canonicalizePatchDeploymentPatchConfigPostStep(des.PostStep, initial.PostStep, opts...)

	return cDes
}

func canonicalizePatchDeploymentPatchConfigSlice(des, initial []PatchDeploymentPatchConfig, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfig(c *Client, des, nw *PatchDeploymentPatchConfig) *PatchDeploymentPatchConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Apt = canonicalizeNewPatchDeploymentPatchConfigApt(c, des.Apt, nw.Apt)
	nw.Yum = canonicalizeNewPatchDeploymentPatchConfigYum(c, des.Yum, nw.Yum)
	nw.Goo = canonicalizeNewPatchDeploymentPatchConfigGoo(c, des.Goo, nw.Goo)
	nw.Zypper = canonicalizeNewPatchDeploymentPatchConfigZypper(c, des.Zypper, nw.Zypper)
	nw.WindowsUpdate = canonicalizeNewPatchDeploymentPatchConfigWindowsUpdate(c, des.WindowsUpdate, nw.WindowsUpdate)
	nw.PreStep = canonicalizeNewPatchDeploymentPatchConfigPreStep(c, des.PreStep, nw.PreStep)
	nw.PostStep = canonicalizeNewPatchDeploymentPatchConfigPostStep(c, des.PostStep, nw.PostStep)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigSet(c *Client, des, nw []PatchDeploymentPatchConfig) []PatchDeploymentPatchConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigSlice(c *Client, des, nw []PatchDeploymentPatchConfig) []PatchDeploymentPatchConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfig(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigApt(des, initial *PatchDeploymentPatchConfigApt, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigApt {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Excludes != nil || (initial != nil && initial.Excludes != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePackages) {
			des.Excludes = nil
			if initial != nil {
				initial.Excludes = nil
			}
		}
	}

	if des.ExclusivePackages != nil || (initial != nil && initial.ExclusivePackages != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Excludes) {
			des.ExclusivePackages = nil
			if initial != nil {
				initial.ExclusivePackages = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigApt{}

	if dcl.IsZeroValue(des.Type) {
		cDes.Type = initial.Type
	} else {
		cDes.Type = des.Type
	}
	if dcl.StringArrayCanonicalize(des.Excludes, initial.Excludes) || dcl.IsZeroValue(des.Excludes) {
		cDes.Excludes = initial.Excludes
	} else {
		cDes.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePackages, initial.ExclusivePackages) || dcl.IsZeroValue(des.ExclusivePackages) {
		cDes.ExclusivePackages = initial.ExclusivePackages
	} else {
		cDes.ExclusivePackages = des.ExclusivePackages
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigAptSlice(des, initial []PatchDeploymentPatchConfigApt, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigApt {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigApt, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigApt(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigApt, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigApt(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigApt(c *Client, des, nw *PatchDeploymentPatchConfigApt) *PatchDeploymentPatchConfigApt {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigApt while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Excludes, nw.Excludes) {
		nw.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePackages, nw.ExclusivePackages) {
		nw.ExclusivePackages = des.ExclusivePackages
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigAptSet(c *Client, des, nw []PatchDeploymentPatchConfigApt) []PatchDeploymentPatchConfigApt {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigApt
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigAptNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigAptSlice(c *Client, des, nw []PatchDeploymentPatchConfigApt) []PatchDeploymentPatchConfigApt {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigApt
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigApt(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigYum(des, initial *PatchDeploymentPatchConfigYum, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigYum {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Excludes != nil || (initial != nil && initial.Excludes != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePackages) {
			des.Excludes = nil
			if initial != nil {
				initial.Excludes = nil
			}
		}
	}

	if des.ExclusivePackages != nil || (initial != nil && initial.ExclusivePackages != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Excludes) {
			des.ExclusivePackages = nil
			if initial != nil {
				initial.ExclusivePackages = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigYum{}

	if dcl.BoolCanonicalize(des.Security, initial.Security) || dcl.IsZeroValue(des.Security) {
		cDes.Security = initial.Security
	} else {
		cDes.Security = des.Security
	}
	if dcl.BoolCanonicalize(des.Minimal, initial.Minimal) || dcl.IsZeroValue(des.Minimal) {
		cDes.Minimal = initial.Minimal
	} else {
		cDes.Minimal = des.Minimal
	}
	if dcl.StringArrayCanonicalize(des.Excludes, initial.Excludes) || dcl.IsZeroValue(des.Excludes) {
		cDes.Excludes = initial.Excludes
	} else {
		cDes.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePackages, initial.ExclusivePackages) || dcl.IsZeroValue(des.ExclusivePackages) {
		cDes.ExclusivePackages = initial.ExclusivePackages
	} else {
		cDes.ExclusivePackages = des.ExclusivePackages
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigYumSlice(des, initial []PatchDeploymentPatchConfigYum, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigYum {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigYum, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigYum(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigYum, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigYum(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigYum(c *Client, des, nw *PatchDeploymentPatchConfigYum) *PatchDeploymentPatchConfigYum {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigYum while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.Security, nw.Security) {
		nw.Security = des.Security
	}
	if dcl.BoolCanonicalize(des.Minimal, nw.Minimal) {
		nw.Minimal = des.Minimal
	}
	if dcl.StringArrayCanonicalize(des.Excludes, nw.Excludes) {
		nw.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePackages, nw.ExclusivePackages) {
		nw.ExclusivePackages = des.ExclusivePackages
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigYumSet(c *Client, des, nw []PatchDeploymentPatchConfigYum) []PatchDeploymentPatchConfigYum {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigYum
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigYumNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigYumSlice(c *Client, des, nw []PatchDeploymentPatchConfigYum) []PatchDeploymentPatchConfigYum {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigYum
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigYum(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigGoo(des, initial *PatchDeploymentPatchConfigGoo, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigGoo {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}
	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigGoo{}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigGooSlice(des, initial []PatchDeploymentPatchConfigGoo, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigGoo {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigGoo, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigGoo(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigGoo, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigGoo(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigGoo(c *Client, des, nw *PatchDeploymentPatchConfigGoo) *PatchDeploymentPatchConfigGoo {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigGoo while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigGooSet(c *Client, des, nw []PatchDeploymentPatchConfigGoo) []PatchDeploymentPatchConfigGoo {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigGoo
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigGooNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigGooSlice(c *Client, des, nw []PatchDeploymentPatchConfigGoo) []PatchDeploymentPatchConfigGoo {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigGoo
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigGoo(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigZypper(des, initial *PatchDeploymentPatchConfigZypper, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigZypper {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.WithOptional != nil || (initial != nil && initial.WithOptional != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.WithOptional = nil
			if initial != nil {
				initial.WithOptional = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.WithOptional) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if des.WithUpdate != nil || (initial != nil && initial.WithUpdate != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.WithUpdate = nil
			if initial != nil {
				initial.WithUpdate = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.WithUpdate) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if des.Categories != nil || (initial != nil && initial.Categories != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.Categories = nil
			if initial != nil {
				initial.Categories = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Categories) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if des.Severities != nil || (initial != nil && initial.Severities != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.Severities = nil
			if initial != nil {
				initial.Severities = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Severities) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if des.Excludes != nil || (initial != nil && initial.Excludes != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.Excludes = nil
			if initial != nil {
				initial.Excludes = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Excludes) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigZypper{}

	if dcl.BoolCanonicalize(des.WithOptional, initial.WithOptional) || dcl.IsZeroValue(des.WithOptional) {
		cDes.WithOptional = initial.WithOptional
	} else {
		cDes.WithOptional = des.WithOptional
	}
	if dcl.BoolCanonicalize(des.WithUpdate, initial.WithUpdate) || dcl.IsZeroValue(des.WithUpdate) {
		cDes.WithUpdate = initial.WithUpdate
	} else {
		cDes.WithUpdate = des.WithUpdate
	}
	if dcl.StringArrayCanonicalize(des.Categories, initial.Categories) || dcl.IsZeroValue(des.Categories) {
		cDes.Categories = initial.Categories
	} else {
		cDes.Categories = des.Categories
	}
	if dcl.StringArrayCanonicalize(des.Severities, initial.Severities) || dcl.IsZeroValue(des.Severities) {
		cDes.Severities = initial.Severities
	} else {
		cDes.Severities = des.Severities
	}
	if dcl.StringArrayCanonicalize(des.Excludes, initial.Excludes) || dcl.IsZeroValue(des.Excludes) {
		cDes.Excludes = initial.Excludes
	} else {
		cDes.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePatches, initial.ExclusivePatches) || dcl.IsZeroValue(des.ExclusivePatches) {
		cDes.ExclusivePatches = initial.ExclusivePatches
	} else {
		cDes.ExclusivePatches = des.ExclusivePatches
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigZypperSlice(des, initial []PatchDeploymentPatchConfigZypper, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigZypper {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigZypper, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigZypper(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigZypper, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigZypper(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigZypper(c *Client, des, nw *PatchDeploymentPatchConfigZypper) *PatchDeploymentPatchConfigZypper {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigZypper while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.WithOptional, nw.WithOptional) {
		nw.WithOptional = des.WithOptional
	}
	if dcl.BoolCanonicalize(des.WithUpdate, nw.WithUpdate) {
		nw.WithUpdate = des.WithUpdate
	}
	if dcl.StringArrayCanonicalize(des.Categories, nw.Categories) {
		nw.Categories = des.Categories
	}
	if dcl.StringArrayCanonicalize(des.Severities, nw.Severities) {
		nw.Severities = des.Severities
	}
	if dcl.StringArrayCanonicalize(des.Excludes, nw.Excludes) {
		nw.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePatches, nw.ExclusivePatches) {
		nw.ExclusivePatches = des.ExclusivePatches
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigZypperSet(c *Client, des, nw []PatchDeploymentPatchConfigZypper) []PatchDeploymentPatchConfigZypper {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigZypper
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigZypperNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigZypperSlice(c *Client, des, nw []PatchDeploymentPatchConfigZypper) []PatchDeploymentPatchConfigZypper {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigZypper
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigZypper(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigWindowsUpdate(des, initial *PatchDeploymentPatchConfigWindowsUpdate, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigWindowsUpdate {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Classifications != nil || (initial != nil && initial.Classifications != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.Classifications = nil
			if initial != nil {
				initial.Classifications = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Classifications) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if des.Excludes != nil || (initial != nil && initial.Excludes != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExclusivePatches) {
			des.Excludes = nil
			if initial != nil {
				initial.Excludes = nil
			}
		}
	}

	if des.ExclusivePatches != nil || (initial != nil && initial.ExclusivePatches != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Excludes) {
			des.ExclusivePatches = nil
			if initial != nil {
				initial.ExclusivePatches = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigWindowsUpdate{}

	if dcl.IsZeroValue(des.Classifications) {
		cDes.Classifications = initial.Classifications
	} else {
		cDes.Classifications = des.Classifications
	}
	if dcl.StringArrayCanonicalize(des.Excludes, initial.Excludes) || dcl.IsZeroValue(des.Excludes) {
		cDes.Excludes = initial.Excludes
	} else {
		cDes.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePatches, initial.ExclusivePatches) || dcl.IsZeroValue(des.ExclusivePatches) {
		cDes.ExclusivePatches = initial.ExclusivePatches
	} else {
		cDes.ExclusivePatches = des.ExclusivePatches
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigWindowsUpdateSlice(des, initial []PatchDeploymentPatchConfigWindowsUpdate, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigWindowsUpdate {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigWindowsUpdate, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigWindowsUpdate(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigWindowsUpdate, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigWindowsUpdate(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigWindowsUpdate(c *Client, des, nw *PatchDeploymentPatchConfigWindowsUpdate) *PatchDeploymentPatchConfigWindowsUpdate {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigWindowsUpdate while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Excludes, nw.Excludes) {
		nw.Excludes = des.Excludes
	}
	if dcl.StringArrayCanonicalize(des.ExclusivePatches, nw.ExclusivePatches) {
		nw.ExclusivePatches = des.ExclusivePatches
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigWindowsUpdateSet(c *Client, des, nw []PatchDeploymentPatchConfigWindowsUpdate) []PatchDeploymentPatchConfigWindowsUpdate {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigWindowsUpdate
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigWindowsUpdateNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigWindowsUpdateSlice(c *Client, des, nw []PatchDeploymentPatchConfigWindowsUpdate) []PatchDeploymentPatchConfigWindowsUpdate {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigWindowsUpdate
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigWindowsUpdate(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPreStep(des, initial *PatchDeploymentPatchConfigPreStep, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPreStep {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPreStep{}

	cDes.LinuxExecStepConfig = canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfig(des.LinuxExecStepConfig, initial.LinuxExecStepConfig, opts...)
	cDes.WindowsExecStepConfig = canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfig(des.WindowsExecStepConfig, initial.WindowsExecStepConfig, opts...)

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPreStepSlice(des, initial []PatchDeploymentPatchConfigPreStep, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPreStep {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPreStep, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPreStep(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPreStep, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPreStep(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPreStep(c *Client, des, nw *PatchDeploymentPatchConfigPreStep) *PatchDeploymentPatchConfigPreStep {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPreStep while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.LinuxExecStepConfig = canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, des.LinuxExecStepConfig, nw.LinuxExecStepConfig)
	nw.WindowsExecStepConfig = canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, des.WindowsExecStepConfig, nw.WindowsExecStepConfig)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepSet(c *Client, des, nw []PatchDeploymentPatchConfigPreStep) []PatchDeploymentPatchConfigPreStep {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPreStep
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPreStepNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepSlice(c *Client, des, nw []PatchDeploymentPatchConfigPreStep) []PatchDeploymentPatchConfigPreStep {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPreStep
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPreStep(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfig(des, initial *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.LocalPath != nil || (initial != nil && initial.LocalPath != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GcsObject) {
			des.LocalPath = nil
			if initial != nil {
				initial.LocalPath = nil
			}
		}
	}

	if des.GcsObject != nil || (initial != nil && initial.GcsObject != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LocalPath) {
			des.GcsObject = nil
			if initial != nil {
				initial.GcsObject = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}

	if dcl.StringCanonicalize(des.LocalPath, initial.LocalPath) || dcl.IsZeroValue(des.LocalPath) {
		cDes.LocalPath = initial.LocalPath
	} else {
		cDes.LocalPath = des.LocalPath
	}
	cDes.GcsObject = canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(des.GcsObject, initial.GcsObject, opts...)
	if dcl.IsZeroValue(des.AllowedSuccessCodes) {
		cDes.AllowedSuccessCodes = initial.AllowedSuccessCodes
	} else {
		cDes.AllowedSuccessCodes = des.AllowedSuccessCodes
	}
	if dcl.IsZeroValue(des.Interpreter) {
		cDes.Interpreter = initial.Interpreter
	} else {
		cDes.Interpreter = des.Interpreter
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice(des, initial []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c *Client, des, nw *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPreStepLinuxExecStepConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	nw.GcsObject = canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, des.GcsObject, nw.GcsObject)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSet(c *Client, des, nw []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice(c *Client, des, nw []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(des, initial *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}

	if dcl.StringCanonicalize(des.Bucket, initial.Bucket) || dcl.IsZeroValue(des.Bucket) {
		cDes.Bucket = initial.Bucket
	} else {
		cDes.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, initial.Object) || dcl.IsZeroValue(des.Object) {
		cDes.Object = initial.Object
	} else {
		cDes.Object = des.Object
	}
	if dcl.IsZeroValue(des.GenerationNumber) {
		cDes.GenerationNumber = initial.GenerationNumber
	} else {
		cDes.GenerationNumber = des.GenerationNumber
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice(des, initial []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c *Client, des, nw *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Bucket, nw.Bucket) {
		nw.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, nw.Object) {
		nw.Object = des.Object
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSet(c *Client, des, nw []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice(c *Client, des, nw []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfig(des, initial *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.LocalPath != nil || (initial != nil && initial.LocalPath != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GcsObject) {
			des.LocalPath = nil
			if initial != nil {
				initial.LocalPath = nil
			}
		}
	}

	if des.GcsObject != nil || (initial != nil && initial.GcsObject != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LocalPath) {
			des.GcsObject = nil
			if initial != nil {
				initial.GcsObject = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}

	if dcl.StringCanonicalize(des.LocalPath, initial.LocalPath) || dcl.IsZeroValue(des.LocalPath) {
		cDes.LocalPath = initial.LocalPath
	} else {
		cDes.LocalPath = des.LocalPath
	}
	cDes.GcsObject = canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(des.GcsObject, initial.GcsObject, opts...)
	if dcl.IsZeroValue(des.AllowedSuccessCodes) {
		cDes.AllowedSuccessCodes = initial.AllowedSuccessCodes
	} else {
		cDes.AllowedSuccessCodes = des.AllowedSuccessCodes
	}
	if dcl.IsZeroValue(des.Interpreter) {
		cDes.Interpreter = initial.Interpreter
	} else {
		cDes.Interpreter = des.Interpreter
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice(des, initial []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c *Client, des, nw *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPreStepWindowsExecStepConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	nw.GcsObject = canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, des.GcsObject, nw.GcsObject)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSet(c *Client, des, nw []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice(c *Client, des, nw []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(des, initial *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}

	if dcl.StringCanonicalize(des.Bucket, initial.Bucket) || dcl.IsZeroValue(des.Bucket) {
		cDes.Bucket = initial.Bucket
	} else {
		cDes.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, initial.Object) || dcl.IsZeroValue(des.Object) {
		cDes.Object = initial.Object
	} else {
		cDes.Object = des.Object
	}
	if dcl.IsZeroValue(des.GenerationNumber) {
		cDes.GenerationNumber = initial.GenerationNumber
	} else {
		cDes.GenerationNumber = des.GenerationNumber
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice(des, initial []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c *Client, des, nw *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Bucket, nw.Bucket) {
		nw.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, nw.Object) {
		nw.Object = des.Object
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSet(c *Client, des, nw []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice(c *Client, des, nw []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPostStep(des, initial *PatchDeploymentPatchConfigPostStep, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPostStep {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPostStep{}

	cDes.LinuxExecStepConfig = canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfig(des.LinuxExecStepConfig, initial.LinuxExecStepConfig, opts...)
	cDes.WindowsExecStepConfig = canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfig(des.WindowsExecStepConfig, initial.WindowsExecStepConfig, opts...)

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPostStepSlice(des, initial []PatchDeploymentPatchConfigPostStep, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPostStep {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPostStep, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPostStep(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPostStep, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPostStep(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPostStep(c *Client, des, nw *PatchDeploymentPatchConfigPostStep) *PatchDeploymentPatchConfigPostStep {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPostStep while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.LinuxExecStepConfig = canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, des.LinuxExecStepConfig, nw.LinuxExecStepConfig)
	nw.WindowsExecStepConfig = canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, des.WindowsExecStepConfig, nw.WindowsExecStepConfig)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepSet(c *Client, des, nw []PatchDeploymentPatchConfigPostStep) []PatchDeploymentPatchConfigPostStep {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPostStep
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPostStepNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepSlice(c *Client, des, nw []PatchDeploymentPatchConfigPostStep) []PatchDeploymentPatchConfigPostStep {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPostStep
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPostStep(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfig(des, initial *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.LocalPath != nil || (initial != nil && initial.LocalPath != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GcsObject) {
			des.LocalPath = nil
			if initial != nil {
				initial.LocalPath = nil
			}
		}
	}

	if des.GcsObject != nil || (initial != nil && initial.GcsObject != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LocalPath) {
			des.GcsObject = nil
			if initial != nil {
				initial.GcsObject = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}

	if dcl.StringCanonicalize(des.LocalPath, initial.LocalPath) || dcl.IsZeroValue(des.LocalPath) {
		cDes.LocalPath = initial.LocalPath
	} else {
		cDes.LocalPath = des.LocalPath
	}
	cDes.GcsObject = canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(des.GcsObject, initial.GcsObject, opts...)
	if dcl.IsZeroValue(des.AllowedSuccessCodes) {
		cDes.AllowedSuccessCodes = initial.AllowedSuccessCodes
	} else {
		cDes.AllowedSuccessCodes = des.AllowedSuccessCodes
	}
	if dcl.IsZeroValue(des.Interpreter) {
		cDes.Interpreter = initial.Interpreter
	} else {
		cDes.Interpreter = des.Interpreter
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice(des, initial []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c *Client, des, nw *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPostStepLinuxExecStepConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	nw.GcsObject = canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, des.GcsObject, nw.GcsObject)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSet(c *Client, des, nw []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice(c *Client, des, nw []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(des, initial *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}

	if dcl.StringCanonicalize(des.Bucket, initial.Bucket) || dcl.IsZeroValue(des.Bucket) {
		cDes.Bucket = initial.Bucket
	} else {
		cDes.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, initial.Object) || dcl.IsZeroValue(des.Object) {
		cDes.Object = initial.Object
	} else {
		cDes.Object = des.Object
	}
	if dcl.IsZeroValue(des.GenerationNumber) {
		cDes.GenerationNumber = initial.GenerationNumber
	} else {
		cDes.GenerationNumber = des.GenerationNumber
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice(des, initial []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c *Client, des, nw *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Bucket, nw.Bucket) {
		nw.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, nw.Object) {
		nw.Object = des.Object
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSet(c *Client, des, nw []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice(c *Client, des, nw []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfig(des, initial *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.LocalPath != nil || (initial != nil && initial.LocalPath != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.GcsObject) {
			des.LocalPath = nil
			if initial != nil {
				initial.LocalPath = nil
			}
		}
	}

	if des.GcsObject != nil || (initial != nil && initial.GcsObject != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LocalPath) {
			des.GcsObject = nil
			if initial != nil {
				initial.GcsObject = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}

	if dcl.StringCanonicalize(des.LocalPath, initial.LocalPath) || dcl.IsZeroValue(des.LocalPath) {
		cDes.LocalPath = initial.LocalPath
	} else {
		cDes.LocalPath = des.LocalPath
	}
	cDes.GcsObject = canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(des.GcsObject, initial.GcsObject, opts...)
	if dcl.IsZeroValue(des.AllowedSuccessCodes) {
		cDes.AllowedSuccessCodes = initial.AllowedSuccessCodes
	} else {
		cDes.AllowedSuccessCodes = des.AllowedSuccessCodes
	}
	if dcl.IsZeroValue(des.Interpreter) {
		cDes.Interpreter = initial.Interpreter
	} else {
		cDes.Interpreter = des.Interpreter
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice(des, initial []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c *Client, des, nw *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPostStepWindowsExecStepConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	nw.GcsObject = canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, des.GcsObject, nw.GcsObject)

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSet(c *Client, des, nw []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice(c *Client, des, nw []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(des, initial *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, opts ...dcl.ApplyOption) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}

	if dcl.StringCanonicalize(des.Bucket, initial.Bucket) || dcl.IsZeroValue(des.Bucket) {
		cDes.Bucket = initial.Bucket
	} else {
		cDes.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, initial.Object) || dcl.IsZeroValue(des.Object) {
		cDes.Object = initial.Object
	} else {
		cDes.Object = des.Object
	}
	if dcl.IsZeroValue(des.GenerationNumber) {
		cDes.GenerationNumber = initial.GenerationNumber
	} else {
		cDes.GenerationNumber = des.GenerationNumber
	}

	return cDes
}

func canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice(des, initial []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, opts ...dcl.ApplyOption) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c *Client, des, nw *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Bucket, nw.Bucket) {
		nw.Bucket = des.Bucket
	}
	if dcl.StringCanonicalize(des.Object, nw.Object) {
		nw.Object = des.Object
	}

	return nw
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSet(c *Client, des, nw []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice(c *Client, des, nw []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentOneTimeSchedule(des, initial *PatchDeploymentOneTimeSchedule, opts ...dcl.ApplyOption) *PatchDeploymentOneTimeSchedule {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentOneTimeSchedule{}

	if dcl.IsZeroValue(des.ExecuteTime) {
		cDes.ExecuteTime = initial.ExecuteTime
	} else {
		cDes.ExecuteTime = des.ExecuteTime
	}

	return cDes
}

func canonicalizePatchDeploymentOneTimeScheduleSlice(des, initial []PatchDeploymentOneTimeSchedule, opts ...dcl.ApplyOption) []PatchDeploymentOneTimeSchedule {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentOneTimeSchedule, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentOneTimeSchedule(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentOneTimeSchedule, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentOneTimeSchedule(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentOneTimeSchedule(c *Client, des, nw *PatchDeploymentOneTimeSchedule) *PatchDeploymentOneTimeSchedule {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentOneTimeSchedule while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentOneTimeScheduleSet(c *Client, des, nw []PatchDeploymentOneTimeSchedule) []PatchDeploymentOneTimeSchedule {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentOneTimeSchedule
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentOneTimeScheduleNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentOneTimeScheduleSlice(c *Client, des, nw []PatchDeploymentOneTimeSchedule) []PatchDeploymentOneTimeSchedule {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentOneTimeSchedule
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentOneTimeSchedule(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringSchedule(des, initial *PatchDeploymentRecurringSchedule, opts ...dcl.ApplyOption) *PatchDeploymentRecurringSchedule {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Weekly != nil || (initial != nil && initial.Weekly != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Monthly) {
			des.Weekly = nil
			if initial != nil {
				initial.Weekly = nil
			}
		}
	}

	if des.Monthly != nil || (initial != nil && initial.Monthly != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Weekly) {
			des.Monthly = nil
			if initial != nil {
				initial.Monthly = nil
			}
		}
	}

	if des.Weekly != nil || (initial != nil && initial.Weekly != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Monthly) {
			des.Weekly = nil
			if initial != nil {
				initial.Weekly = nil
			}
		}
	}

	if des.Monthly != nil || (initial != nil && initial.Monthly != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Weekly) {
			des.Monthly = nil
			if initial != nil {
				initial.Monthly = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringSchedule{}

	cDes.TimeZone = canonicalizePatchDeploymentRecurringScheduleTimeZone(des.TimeZone, initial.TimeZone, opts...)
	if dcl.IsZeroValue(des.StartTime) {
		cDes.StartTime = initial.StartTime
	} else {
		cDes.StartTime = des.StartTime
	}
	if dcl.IsZeroValue(des.EndTime) {
		cDes.EndTime = initial.EndTime
	} else {
		cDes.EndTime = des.EndTime
	}
	cDes.TimeOfDay = canonicalizePatchDeploymentRecurringScheduleTimeOfDay(des.TimeOfDay, initial.TimeOfDay, opts...)
	if dcl.IsZeroValue(des.Frequency) {
		cDes.Frequency = initial.Frequency
	} else {
		cDes.Frequency = des.Frequency
	}
	cDes.Weekly = canonicalizePatchDeploymentRecurringScheduleWeekly(des.Weekly, initial.Weekly, opts...)
	cDes.Monthly = canonicalizePatchDeploymentRecurringScheduleMonthly(des.Monthly, initial.Monthly, opts...)

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleSlice(des, initial []PatchDeploymentRecurringSchedule, opts ...dcl.ApplyOption) []PatchDeploymentRecurringSchedule {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringSchedule, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringSchedule(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringSchedule, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringSchedule(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringSchedule(c *Client, des, nw *PatchDeploymentRecurringSchedule) *PatchDeploymentRecurringSchedule {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringSchedule while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.TimeZone = canonicalizeNewPatchDeploymentRecurringScheduleTimeZone(c, des.TimeZone, nw.TimeZone)
	nw.TimeOfDay = canonicalizeNewPatchDeploymentRecurringScheduleTimeOfDay(c, des.TimeOfDay, nw.TimeOfDay)
	nw.Weekly = canonicalizeNewPatchDeploymentRecurringScheduleWeekly(c, des.Weekly, nw.Weekly)
	nw.Monthly = canonicalizeNewPatchDeploymentRecurringScheduleMonthly(c, des.Monthly, nw.Monthly)

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleSet(c *Client, des, nw []PatchDeploymentRecurringSchedule) []PatchDeploymentRecurringSchedule {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringSchedule
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleSlice(c *Client, des, nw []PatchDeploymentRecurringSchedule) []PatchDeploymentRecurringSchedule {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringSchedule
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringSchedule(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringScheduleTimeZone(des, initial *PatchDeploymentRecurringScheduleTimeZone, opts ...dcl.ApplyOption) *PatchDeploymentRecurringScheduleTimeZone {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringScheduleTimeZone{}

	if dcl.StringCanonicalize(des.Id, initial.Id) || dcl.IsZeroValue(des.Id) {
		cDes.Id = initial.Id
	} else {
		cDes.Id = des.Id
	}
	if dcl.StringCanonicalize(des.Version, initial.Version) || dcl.IsZeroValue(des.Version) {
		cDes.Version = initial.Version
	} else {
		cDes.Version = des.Version
	}

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleTimeZoneSlice(des, initial []PatchDeploymentRecurringScheduleTimeZone, opts ...dcl.ApplyOption) []PatchDeploymentRecurringScheduleTimeZone {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringScheduleTimeZone, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringScheduleTimeZone(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringScheduleTimeZone, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringScheduleTimeZone(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeZone(c *Client, des, nw *PatchDeploymentRecurringScheduleTimeZone) *PatchDeploymentRecurringScheduleTimeZone {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringScheduleTimeZone while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Id, nw.Id) {
		nw.Id = des.Id
	}
	if dcl.StringCanonicalize(des.Version, nw.Version) {
		nw.Version = des.Version
	}

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeZoneSet(c *Client, des, nw []PatchDeploymentRecurringScheduleTimeZone) []PatchDeploymentRecurringScheduleTimeZone {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringScheduleTimeZone
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleTimeZoneNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeZoneSlice(c *Client, des, nw []PatchDeploymentRecurringScheduleTimeZone) []PatchDeploymentRecurringScheduleTimeZone {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringScheduleTimeZone
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringScheduleTimeZone(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringScheduleTimeOfDay(des, initial *PatchDeploymentRecurringScheduleTimeOfDay, opts ...dcl.ApplyOption) *PatchDeploymentRecurringScheduleTimeOfDay {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringScheduleTimeOfDay{}

	if dcl.IsZeroValue(des.Hours) {
		cDes.Hours = initial.Hours
	} else {
		cDes.Hours = des.Hours
	}
	if dcl.IsZeroValue(des.Minutes) {
		cDes.Minutes = initial.Minutes
	} else {
		cDes.Minutes = des.Minutes
	}
	if dcl.IsZeroValue(des.Seconds) {
		cDes.Seconds = initial.Seconds
	} else {
		cDes.Seconds = des.Seconds
	}
	if dcl.IsZeroValue(des.Nanos) {
		cDes.Nanos = initial.Nanos
	} else {
		cDes.Nanos = des.Nanos
	}

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleTimeOfDaySlice(des, initial []PatchDeploymentRecurringScheduleTimeOfDay, opts ...dcl.ApplyOption) []PatchDeploymentRecurringScheduleTimeOfDay {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringScheduleTimeOfDay, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringScheduleTimeOfDay(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringScheduleTimeOfDay, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringScheduleTimeOfDay(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeOfDay(c *Client, des, nw *PatchDeploymentRecurringScheduleTimeOfDay) *PatchDeploymentRecurringScheduleTimeOfDay {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringScheduleTimeOfDay while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeOfDaySet(c *Client, des, nw []PatchDeploymentRecurringScheduleTimeOfDay) []PatchDeploymentRecurringScheduleTimeOfDay {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringScheduleTimeOfDay
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleTimeOfDayNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleTimeOfDaySlice(c *Client, des, nw []PatchDeploymentRecurringScheduleTimeOfDay) []PatchDeploymentRecurringScheduleTimeOfDay {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringScheduleTimeOfDay
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringScheduleTimeOfDay(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringScheduleWeekly(des, initial *PatchDeploymentRecurringScheduleWeekly, opts ...dcl.ApplyOption) *PatchDeploymentRecurringScheduleWeekly {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringScheduleWeekly{}

	if dcl.IsZeroValue(des.DayOfWeek) {
		cDes.DayOfWeek = initial.DayOfWeek
	} else {
		cDes.DayOfWeek = des.DayOfWeek
	}

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleWeeklySlice(des, initial []PatchDeploymentRecurringScheduleWeekly, opts ...dcl.ApplyOption) []PatchDeploymentRecurringScheduleWeekly {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringScheduleWeekly, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringScheduleWeekly(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringScheduleWeekly, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringScheduleWeekly(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringScheduleWeekly(c *Client, des, nw *PatchDeploymentRecurringScheduleWeekly) *PatchDeploymentRecurringScheduleWeekly {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringScheduleWeekly while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleWeeklySet(c *Client, des, nw []PatchDeploymentRecurringScheduleWeekly) []PatchDeploymentRecurringScheduleWeekly {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringScheduleWeekly
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleWeeklyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleWeeklySlice(c *Client, des, nw []PatchDeploymentRecurringScheduleWeekly) []PatchDeploymentRecurringScheduleWeekly {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringScheduleWeekly
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringScheduleWeekly(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringScheduleMonthly(des, initial *PatchDeploymentRecurringScheduleMonthly, opts ...dcl.ApplyOption) *PatchDeploymentRecurringScheduleMonthly {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.WeekDayOfMonth != nil || (initial != nil && initial.WeekDayOfMonth != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.MonthDay) {
			des.WeekDayOfMonth = nil
			if initial != nil {
				initial.WeekDayOfMonth = nil
			}
		}
	}

	if des.MonthDay != nil || (initial != nil && initial.MonthDay != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.WeekDayOfMonth) {
			des.MonthDay = nil
			if initial != nil {
				initial.MonthDay = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringScheduleMonthly{}

	cDes.WeekDayOfMonth = canonicalizePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(des.WeekDayOfMonth, initial.WeekDayOfMonth, opts...)
	if dcl.IsZeroValue(des.MonthDay) {
		cDes.MonthDay = initial.MonthDay
	} else {
		cDes.MonthDay = des.MonthDay
	}

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleMonthlySlice(des, initial []PatchDeploymentRecurringScheduleMonthly, opts ...dcl.ApplyOption) []PatchDeploymentRecurringScheduleMonthly {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringScheduleMonthly, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringScheduleMonthly(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringScheduleMonthly, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringScheduleMonthly(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthly(c *Client, des, nw *PatchDeploymentRecurringScheduleMonthly) *PatchDeploymentRecurringScheduleMonthly {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringScheduleMonthly while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.WeekDayOfMonth = canonicalizeNewPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, des.WeekDayOfMonth, nw.WeekDayOfMonth)

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthlySet(c *Client, des, nw []PatchDeploymentRecurringScheduleMonthly) []PatchDeploymentRecurringScheduleMonthly {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringScheduleMonthly
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleMonthlyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthlySlice(c *Client, des, nw []PatchDeploymentRecurringScheduleMonthly) []PatchDeploymentRecurringScheduleMonthly {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringScheduleMonthly
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringScheduleMonthly(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(des, initial *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, opts ...dcl.ApplyOption) *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}

	if dcl.IsZeroValue(des.WeekOrdinal) {
		cDes.WeekOrdinal = initial.WeekOrdinal
	} else {
		cDes.WeekOrdinal = des.WeekOrdinal
	}
	if dcl.IsZeroValue(des.DayOfWeek) {
		cDes.DayOfWeek = initial.DayOfWeek
	} else {
		cDes.DayOfWeek = des.DayOfWeek
	}

	return cDes
}

func canonicalizePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice(des, initial []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, opts ...dcl.ApplyOption) []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c *Client, des, nw *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSet(c *Client, des, nw []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice(c *Client, des, nw []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRollout(des, initial *PatchDeploymentRollout, opts ...dcl.ApplyOption) *PatchDeploymentRollout {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRollout{}

	if dcl.IsZeroValue(des.Mode) {
		cDes.Mode = initial.Mode
	} else {
		cDes.Mode = des.Mode
	}
	cDes.DisruptionBudget = canonicalizePatchDeploymentRolloutDisruptionBudget(des.DisruptionBudget, initial.DisruptionBudget, opts...)

	return cDes
}

func canonicalizePatchDeploymentRolloutSlice(des, initial []PatchDeploymentRollout, opts ...dcl.ApplyOption) []PatchDeploymentRollout {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRollout, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRollout(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRollout, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRollout(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRollout(c *Client, des, nw *PatchDeploymentRollout) *PatchDeploymentRollout {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRollout while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.DisruptionBudget = canonicalizeNewPatchDeploymentRolloutDisruptionBudget(c, des.DisruptionBudget, nw.DisruptionBudget)

	return nw
}

func canonicalizeNewPatchDeploymentRolloutSet(c *Client, des, nw []PatchDeploymentRollout) []PatchDeploymentRollout {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRollout
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRolloutNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRolloutSlice(c *Client, des, nw []PatchDeploymentRollout) []PatchDeploymentRollout {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRollout
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRollout(c, &d, &n))
	}

	return items
}

func canonicalizePatchDeploymentRolloutDisruptionBudget(des, initial *PatchDeploymentRolloutDisruptionBudget, opts ...dcl.ApplyOption) *PatchDeploymentRolloutDisruptionBudget {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Fixed != nil || (initial != nil && initial.Fixed != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Percent) {
			des.Fixed = nil
			if initial != nil {
				initial.Fixed = nil
			}
		}
	}

	if des.Percent != nil || (initial != nil && initial.Percent != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Fixed) {
			des.Percent = nil
			if initial != nil {
				initial.Percent = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &PatchDeploymentRolloutDisruptionBudget{}

	if dcl.IsZeroValue(des.Fixed) {
		cDes.Fixed = initial.Fixed
	} else {
		cDes.Fixed = des.Fixed
	}
	if dcl.IsZeroValue(des.Percent) {
		cDes.Percent = initial.Percent
	} else {
		cDes.Percent = des.Percent
	}

	return cDes
}

func canonicalizePatchDeploymentRolloutDisruptionBudgetSlice(des, initial []PatchDeploymentRolloutDisruptionBudget, opts ...dcl.ApplyOption) []PatchDeploymentRolloutDisruptionBudget {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PatchDeploymentRolloutDisruptionBudget, 0, len(des))
		for _, d := range des {
			cd := canonicalizePatchDeploymentRolloutDisruptionBudget(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PatchDeploymentRolloutDisruptionBudget, 0, len(des))
	for i, d := range des {
		cd := canonicalizePatchDeploymentRolloutDisruptionBudget(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPatchDeploymentRolloutDisruptionBudget(c *Client, des, nw *PatchDeploymentRolloutDisruptionBudget) *PatchDeploymentRolloutDisruptionBudget {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for PatchDeploymentRolloutDisruptionBudget while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewPatchDeploymentRolloutDisruptionBudgetSet(c *Client, des, nw []PatchDeploymentRolloutDisruptionBudget) []PatchDeploymentRolloutDisruptionBudget {
	if des == nil {
		return nw
	}
	var reorderedNew []PatchDeploymentRolloutDisruptionBudget
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := comparePatchDeploymentRolloutDisruptionBudgetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewPatchDeploymentRolloutDisruptionBudgetSlice(c *Client, des, nw []PatchDeploymentRolloutDisruptionBudget) []PatchDeploymentRolloutDisruptionBudget {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PatchDeploymentRolloutDisruptionBudget
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPatchDeploymentRolloutDisruptionBudget(c, &d, &n))
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
func diffPatchDeployment(c *Client, desired, actual *PatchDeployment, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceFilter, actual.InstanceFilter, dcl.Info{ObjectFunction: comparePatchDeploymentInstanceFilterNewStyle, EmptyObject: EmptyPatchDeploymentInstanceFilter, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("InstanceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PatchConfig, actual.PatchConfig, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfig, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("PatchConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Duration, actual.Duration, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Duration")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OneTimeSchedule, actual.OneTimeSchedule, dcl.Info{ObjectFunction: comparePatchDeploymentOneTimeScheduleNewStyle, EmptyObject: EmptyPatchDeploymentOneTimeSchedule, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("OneTimeSchedule")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RecurringSchedule, actual.RecurringSchedule, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleNewStyle, EmptyObject: EmptyPatchDeploymentRecurringSchedule, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("RecurringSchedule")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LastExecuteTime, actual.LastExecuteTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LastExecuteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rollout, actual.Rollout, dcl.Info{ObjectFunction: comparePatchDeploymentRolloutNewStyle, EmptyObject: EmptyPatchDeploymentRollout, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Rollout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	return newDiffs, nil
}
func comparePatchDeploymentInstanceFilterNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentInstanceFilter)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentInstanceFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentInstanceFilter or *PatchDeploymentInstanceFilter", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentInstanceFilter)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentInstanceFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentInstanceFilter", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.All, actual.All, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("All")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GroupLabels, actual.GroupLabels, dcl.Info{ObjectFunction: comparePatchDeploymentInstanceFilterGroupLabelsNewStyle, EmptyObject: EmptyPatchDeploymentInstanceFilterGroupLabels, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GroupLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zones, actual.Zones, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Zones")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Instances, actual.Instances, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Instances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNamePrefixes, actual.InstanceNamePrefixes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("InstanceNamePrefixes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentInstanceFilterGroupLabelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentInstanceFilterGroupLabels)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentInstanceFilterGroupLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentInstanceFilterGroupLabels or *PatchDeploymentInstanceFilterGroupLabels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentInstanceFilterGroupLabels)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentInstanceFilterGroupLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentInstanceFilterGroupLabels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfig)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfig or *PatchDeploymentPatchConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfig)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RebootConfig, actual.RebootConfig, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("RebootConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Apt, actual.Apt, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigAptNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigApt, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Apt")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Yum, actual.Yum, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigYumNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigYum, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Yum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Goo, actual.Goo, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigGooNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigGoo, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Goo")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zypper, actual.Zypper, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigZypperNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigZypper, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Zypper")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WindowsUpdate, actual.WindowsUpdate, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigWindowsUpdateNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigWindowsUpdate, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WindowsUpdate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PreStep, actual.PreStep, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPreStepNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPreStep, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("PreStep")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PostStep, actual.PostStep, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPostStepNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPostStep, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("PostStep")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigAptNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigApt)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigApt or *PatchDeploymentPatchConfigApt", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigApt)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigApt", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Excludes, actual.Excludes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Excludes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusivePackages, actual.ExclusivePackages, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("ExclusivePackages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigYumNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigYum)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigYum or *PatchDeploymentPatchConfigYum", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigYum)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigYum", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Security, actual.Security, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Security")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Minimal, actual.Minimal, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Minimal")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Excludes, actual.Excludes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Excludes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusivePackages, actual.ExclusivePackages, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("ExclusivePackages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigGooNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	return diffs, nil
}

func comparePatchDeploymentPatchConfigZypperNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigZypper)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigZypper or *PatchDeploymentPatchConfigZypper", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigZypper)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigZypper", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.WithOptional, actual.WithOptional, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WithOptional")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WithUpdate, actual.WithUpdate, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WithUpdate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Categories, actual.Categories, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Categories")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Severities, actual.Severities, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Severities")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Excludes, actual.Excludes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Excludes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusivePatches, actual.ExclusivePatches, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("ExclusivePatches")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigWindowsUpdateNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigWindowsUpdate)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigWindowsUpdate)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigWindowsUpdate or *PatchDeploymentPatchConfigWindowsUpdate", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigWindowsUpdate)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigWindowsUpdate)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigWindowsUpdate", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Classifications, actual.Classifications, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Classifications")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Excludes, actual.Excludes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Excludes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusivePatches, actual.ExclusivePatches, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("ExclusivePatches")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPreStepNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPreStep)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPreStep)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStep or *PatchDeploymentPatchConfigPreStep", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPreStep)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPreStep)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStep", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LinuxExecStepConfig, actual.LinuxExecStepConfig, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfig, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LinuxExecStepConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WindowsExecStepConfig, actual.WindowsExecStepConfig, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfig, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WindowsExecStepConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPreStepLinuxExecStepConfig)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPreStepLinuxExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepLinuxExecStepConfig or *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPreStepLinuxExecStepConfig)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPreStepLinuxExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepLinuxExecStepConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GcsObject, actual.GcsObject, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GcsObject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedSuccessCodes, actual.AllowedSuccessCodes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("AllowedSuccessCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject or *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GenerationNumber, actual.GenerationNumber, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GenerationNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPreStepWindowsExecStepConfig)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPreStepWindowsExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepWindowsExecStepConfig or *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPreStepWindowsExecStepConfig)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPreStepWindowsExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepWindowsExecStepConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GcsObject, actual.GcsObject, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GcsObject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedSuccessCodes, actual.AllowedSuccessCodes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("AllowedSuccessCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject or *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GenerationNumber, actual.GenerationNumber, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GenerationNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPostStepNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPostStep)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPostStep)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStep or *PatchDeploymentPatchConfigPostStep", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPostStep)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPostStep)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStep", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LinuxExecStepConfig, actual.LinuxExecStepConfig, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfig, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LinuxExecStepConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WindowsExecStepConfig, actual.WindowsExecStepConfig, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfig, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WindowsExecStepConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPostStepLinuxExecStepConfig)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPostStepLinuxExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepLinuxExecStepConfig or *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPostStepLinuxExecStepConfig)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPostStepLinuxExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepLinuxExecStepConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GcsObject, actual.GcsObject, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GcsObject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedSuccessCodes, actual.AllowedSuccessCodes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("AllowedSuccessCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject or *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GenerationNumber, actual.GenerationNumber, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GenerationNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPostStepWindowsExecStepConfig)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPostStepWindowsExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepWindowsExecStepConfig or *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPostStepWindowsExecStepConfig)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPostStepWindowsExecStepConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepWindowsExecStepConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GcsObject, actual.GcsObject, dcl.Info{ObjectFunction: comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectNewStyle, EmptyObject: EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GcsObject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedSuccessCodes, actual.AllowedSuccessCodes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("AllowedSuccessCodes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject or *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GenerationNumber, actual.GenerationNumber, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("GenerationNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentOneTimeScheduleNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentOneTimeSchedule)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentOneTimeSchedule)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentOneTimeSchedule or *PatchDeploymentOneTimeSchedule", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentOneTimeSchedule)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentOneTimeSchedule)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentOneTimeSchedule", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ExecuteTime, actual.ExecuteTime, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("ExecuteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringSchedule)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringSchedule)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringSchedule or *PatchDeploymentRecurringSchedule", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringSchedule)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringSchedule)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringSchedule", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.TimeZone, actual.TimeZone, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleTimeZoneNewStyle, EmptyObject: EmptyPatchDeploymentRecurringScheduleTimeZone, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("TimeZone")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StartTime, actual.StartTime, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("StartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EndTime, actual.EndTime, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("EndTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TimeOfDay, actual.TimeOfDay, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleTimeOfDayNewStyle, EmptyObject: EmptyPatchDeploymentRecurringScheduleTimeOfDay, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("TimeOfDay")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Frequency, actual.Frequency, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Frequency")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Weekly, actual.Weekly, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleWeeklyNewStyle, EmptyObject: EmptyPatchDeploymentRecurringScheduleWeekly, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Weekly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Monthly, actual.Monthly, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleMonthlyNewStyle, EmptyObject: EmptyPatchDeploymentRecurringScheduleMonthly, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Monthly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LastExecuteTime, actual.LastExecuteTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LastExecuteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NextExecuteTime, actual.NextExecuteTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NextExecuteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleTimeZoneNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringScheduleTimeZone)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringScheduleTimeZone)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleTimeZone or *PatchDeploymentRecurringScheduleTimeZone", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringScheduleTimeZone)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringScheduleTimeZone)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleTimeZone", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleTimeOfDayNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringScheduleTimeOfDay)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringScheduleTimeOfDay)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleTimeOfDay or *PatchDeploymentRecurringScheduleTimeOfDay", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringScheduleTimeOfDay)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringScheduleTimeOfDay)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleTimeOfDay", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Hours, actual.Hours, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Hours")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Minutes, actual.Minutes, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Minutes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Seconds, actual.Seconds, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Seconds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Nanos, actual.Nanos, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Nanos")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleWeeklyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringScheduleWeekly)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringScheduleWeekly)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleWeekly or *PatchDeploymentRecurringScheduleWeekly", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringScheduleWeekly)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringScheduleWeekly)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleWeekly", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DayOfWeek, actual.DayOfWeek, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("DayOfWeek")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleMonthlyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringScheduleMonthly)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringScheduleMonthly)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleMonthly or *PatchDeploymentRecurringScheduleMonthly", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringScheduleMonthly)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringScheduleMonthly)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleMonthly", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.WeekDayOfMonth, actual.WeekDayOfMonth, dcl.Info{ObjectFunction: comparePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthNewStyle, EmptyObject: EmptyPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WeekDayOfMonth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MonthDay, actual.MonthDay, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("MonthDay")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth or *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.WeekOrdinal, actual.WeekOrdinal, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("WeekOrdinal")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DayOfWeek, actual.DayOfWeek, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("DayOfWeek")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRolloutNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRollout)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRollout)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRollout or *PatchDeploymentRollout", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRollout)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRollout)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRollout", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Mode, actual.Mode, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Mode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisruptionBudget, actual.DisruptionBudget, dcl.Info{ObjectFunction: comparePatchDeploymentRolloutDisruptionBudgetNewStyle, EmptyObject: EmptyPatchDeploymentRolloutDisruptionBudget, OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("DisruptionBudget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePatchDeploymentRolloutDisruptionBudgetNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PatchDeploymentRolloutDisruptionBudget)
	if !ok {
		desiredNotPointer, ok := d.(PatchDeploymentRolloutDisruptionBudget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRolloutDisruptionBudget or *PatchDeploymentRolloutDisruptionBudget", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PatchDeploymentRolloutDisruptionBudget)
	if !ok {
		actualNotPointer, ok := a.(PatchDeploymentRolloutDisruptionBudget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PatchDeploymentRolloutDisruptionBudget", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Fixed, actual.Fixed, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Fixed")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Percent, actual.Percent, dcl.Info{OperationSelector: dcl.TriggersOperation("updatePatchDeploymentUpdatePatchDeploymentOperation")}, fn.AddNest("Percent")); len(ds) != 0 || err != nil {
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
func (r *PatchDeployment) urlNormalized() *PatchDeployment {
	normalized := dcl.Copy(*r).(PatchDeployment)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Duration = dcl.SelfLinkToName(r.Duration)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *PatchDeployment) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdatePatchDeployment" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/patchDeployments/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the PatchDeployment resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *PatchDeployment) marshal(c *Client) ([]byte, error) {
	m, err := expandPatchDeployment(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling PatchDeployment: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalPatchDeployment decodes JSON responses into the PatchDeployment resource schema.
func unmarshalPatchDeployment(b []byte, c *Client) (*PatchDeployment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapPatchDeployment(m, c)
}

func unmarshalMapPatchDeployment(m map[string]interface{}, c *Client) (*PatchDeployment, error) {

	flattened := flattenPatchDeployment(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandPatchDeployment expands PatchDeployment into a JSON request object.
func expandPatchDeployment(c *Client, f *PatchDeployment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s/patchDeployments/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if v != nil {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandPatchDeploymentInstanceFilter(c, f.InstanceFilter); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if v != nil {
		m["instanceFilter"] = v
	}
	if v, err := expandPatchDeploymentPatchConfig(c, f.PatchConfig); err != nil {
		return nil, fmt.Errorf("error expanding PatchConfig into patchConfig: %w", err)
	} else if v != nil {
		m["patchConfig"] = v
	}
	if v := f.Duration; dcl.ValueShouldBeSent(v) {
		m["duration"] = v
	}
	if v, err := expandPatchDeploymentOneTimeSchedule(c, f.OneTimeSchedule); err != nil {
		return nil, fmt.Errorf("error expanding OneTimeSchedule into oneTimeSchedule: %w", err)
	} else if v != nil {
		m["oneTimeSchedule"] = v
	}
	if v, err := expandPatchDeploymentRecurringSchedule(c, f.RecurringSchedule); err != nil {
		return nil, fmt.Errorf("error expanding RecurringSchedule into recurringSchedule: %w", err)
	} else if v != nil {
		m["recurringSchedule"] = v
	}
	if v, err := expandPatchDeploymentRollout(c, f.Rollout); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if v != nil {
		m["rollout"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if v != nil {
		m["project"] = v
	}

	return m, nil
}

// flattenPatchDeployment flattens PatchDeployment from a JSON request object into the
// PatchDeployment type.
func flattenPatchDeployment(c *Client, i interface{}) *PatchDeployment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &PatchDeployment{}
	res.Name = dcl.FlattenString(m["name"])
	res.Description = dcl.FlattenString(m["description"])
	res.InstanceFilter = flattenPatchDeploymentInstanceFilter(c, m["instanceFilter"])
	res.PatchConfig = flattenPatchDeploymentPatchConfig(c, m["patchConfig"])
	res.Duration = dcl.FlattenString(m["duration"])
	res.OneTimeSchedule = flattenPatchDeploymentOneTimeSchedule(c, m["oneTimeSchedule"])
	res.RecurringSchedule = flattenPatchDeploymentRecurringSchedule(c, m["recurringSchedule"])
	res.CreateTime = dcl.FlattenString(m["createTime"])
	res.UpdateTime = dcl.FlattenString(m["updateTime"])
	res.LastExecuteTime = dcl.FlattenString(m["lastExecuteTime"])
	res.Rollout = flattenPatchDeploymentRollout(c, m["rollout"])
	res.Project = dcl.FlattenString(m["project"])

	return res
}

// expandPatchDeploymentInstanceFilterMap expands the contents of PatchDeploymentInstanceFilter into a JSON
// request object.
func expandPatchDeploymentInstanceFilterMap(c *Client, f map[string]PatchDeploymentInstanceFilter) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentInstanceFilter(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentInstanceFilterSlice expands the contents of PatchDeploymentInstanceFilter into a JSON
// request object.
func expandPatchDeploymentInstanceFilterSlice(c *Client, f []PatchDeploymentInstanceFilter) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentInstanceFilter(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentInstanceFilterMap flattens the contents of PatchDeploymentInstanceFilter from a JSON
// response object.
func flattenPatchDeploymentInstanceFilterMap(c *Client, i interface{}) map[string]PatchDeploymentInstanceFilter {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentInstanceFilter{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentInstanceFilter{}
	}

	items := make(map[string]PatchDeploymentInstanceFilter)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentInstanceFilter(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentInstanceFilterSlice flattens the contents of PatchDeploymentInstanceFilter from a JSON
// response object.
func flattenPatchDeploymentInstanceFilterSlice(c *Client, i interface{}) []PatchDeploymentInstanceFilter {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentInstanceFilter{}
	}

	if len(a) == 0 {
		return []PatchDeploymentInstanceFilter{}
	}

	items := make([]PatchDeploymentInstanceFilter, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentInstanceFilter(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentInstanceFilter expands an instance of PatchDeploymentInstanceFilter into a JSON
// request object.
func expandPatchDeploymentInstanceFilter(c *Client, f *PatchDeploymentInstanceFilter) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.All; !dcl.IsEmptyValueIndirect(v) {
		m["all"] = v
	}
	if v, err := expandPatchDeploymentInstanceFilterGroupLabelsSlice(c, f.GroupLabels); err != nil {
		return nil, fmt.Errorf("error expanding GroupLabels into groupLabels: %w", err)
	} else if v != nil {
		m["groupLabels"] = v
	}
	if v := f.Zones; v != nil {
		m["zones"] = v
	}
	if v := f.Instances; v != nil {
		m["instances"] = v
	}
	if v := f.InstanceNamePrefixes; v != nil {
		m["instanceNamePrefixes"] = v
	}

	return m, nil
}

// flattenPatchDeploymentInstanceFilter flattens an instance of PatchDeploymentInstanceFilter from a JSON
// response object.
func flattenPatchDeploymentInstanceFilter(c *Client, i interface{}) *PatchDeploymentInstanceFilter {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentInstanceFilter{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentInstanceFilter
	}
	r.All = dcl.FlattenBool(m["all"])
	r.GroupLabels = flattenPatchDeploymentInstanceFilterGroupLabelsSlice(c, m["groupLabels"])
	r.Zones = dcl.FlattenStringSlice(m["zones"])
	r.Instances = dcl.FlattenStringSlice(m["instances"])
	r.InstanceNamePrefixes = dcl.FlattenStringSlice(m["instanceNamePrefixes"])

	return r
}

// expandPatchDeploymentInstanceFilterGroupLabelsMap expands the contents of PatchDeploymentInstanceFilterGroupLabels into a JSON
// request object.
func expandPatchDeploymentInstanceFilterGroupLabelsMap(c *Client, f map[string]PatchDeploymentInstanceFilterGroupLabels) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentInstanceFilterGroupLabels(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentInstanceFilterGroupLabelsSlice expands the contents of PatchDeploymentInstanceFilterGroupLabels into a JSON
// request object.
func expandPatchDeploymentInstanceFilterGroupLabelsSlice(c *Client, f []PatchDeploymentInstanceFilterGroupLabels) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentInstanceFilterGroupLabels(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentInstanceFilterGroupLabelsMap flattens the contents of PatchDeploymentInstanceFilterGroupLabels from a JSON
// response object.
func flattenPatchDeploymentInstanceFilterGroupLabelsMap(c *Client, i interface{}) map[string]PatchDeploymentInstanceFilterGroupLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentInstanceFilterGroupLabels{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentInstanceFilterGroupLabels{}
	}

	items := make(map[string]PatchDeploymentInstanceFilterGroupLabels)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentInstanceFilterGroupLabels(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentInstanceFilterGroupLabelsSlice flattens the contents of PatchDeploymentInstanceFilterGroupLabels from a JSON
// response object.
func flattenPatchDeploymentInstanceFilterGroupLabelsSlice(c *Client, i interface{}) []PatchDeploymentInstanceFilterGroupLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentInstanceFilterGroupLabels{}
	}

	if len(a) == 0 {
		return []PatchDeploymentInstanceFilterGroupLabels{}
	}

	items := make([]PatchDeploymentInstanceFilterGroupLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentInstanceFilterGroupLabels(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentInstanceFilterGroupLabels expands an instance of PatchDeploymentInstanceFilterGroupLabels into a JSON
// request object.
func expandPatchDeploymentInstanceFilterGroupLabels(c *Client, f *PatchDeploymentInstanceFilterGroupLabels) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}

	return m, nil
}

// flattenPatchDeploymentInstanceFilterGroupLabels flattens an instance of PatchDeploymentInstanceFilterGroupLabels from a JSON
// response object.
func flattenPatchDeploymentInstanceFilterGroupLabels(c *Client, i interface{}) *PatchDeploymentInstanceFilterGroupLabels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentInstanceFilterGroupLabels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentInstanceFilterGroupLabels
	}
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])

	return r
}

// expandPatchDeploymentPatchConfigMap expands the contents of PatchDeploymentPatchConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigMap(c *Client, f map[string]PatchDeploymentPatchConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigSlice expands the contents of PatchDeploymentPatchConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigSlice(c *Client, f []PatchDeploymentPatchConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigMap flattens the contents of PatchDeploymentPatchConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfig{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfig{}
	}

	items := make(map[string]PatchDeploymentPatchConfig)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigSlice flattens the contents of PatchDeploymentPatchConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigSlice(c *Client, i interface{}) []PatchDeploymentPatchConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfig{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfig{}
	}

	items := make([]PatchDeploymentPatchConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfig expands an instance of PatchDeploymentPatchConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfig(c *Client, f *PatchDeploymentPatchConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RebootConfig; !dcl.IsEmptyValueIndirect(v) {
		m["rebootConfig"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigApt(c, f.Apt); err != nil {
		return nil, fmt.Errorf("error expanding Apt into apt: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["apt"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigYum(c, f.Yum); err != nil {
		return nil, fmt.Errorf("error expanding Yum into yum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yum"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigGoo(c, f.Goo); err != nil {
		return nil, fmt.Errorf("error expanding Goo into goo: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goo"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigZypper(c, f.Zypper); err != nil {
		return nil, fmt.Errorf("error expanding Zypper into zypper: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["zypper"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigWindowsUpdate(c, f.WindowsUpdate); err != nil {
		return nil, fmt.Errorf("error expanding WindowsUpdate into windowsUpdate: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["windowsUpdate"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPreStep(c, f.PreStep); err != nil {
		return nil, fmt.Errorf("error expanding PreStep into preStep: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["preStep"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPostStep(c, f.PostStep); err != nil {
		return nil, fmt.Errorf("error expanding PostStep into postStep: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["postStep"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfig flattens an instance of PatchDeploymentPatchConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfig(c *Client, i interface{}) *PatchDeploymentPatchConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfig
	}
	r.RebootConfig = flattenPatchDeploymentPatchConfigRebootConfigEnum(m["rebootConfig"])
	r.Apt = flattenPatchDeploymentPatchConfigApt(c, m["apt"])
	r.Yum = flattenPatchDeploymentPatchConfigYum(c, m["yum"])
	r.Goo = flattenPatchDeploymentPatchConfigGoo(c, m["goo"])
	r.Zypper = flattenPatchDeploymentPatchConfigZypper(c, m["zypper"])
	r.WindowsUpdate = flattenPatchDeploymentPatchConfigWindowsUpdate(c, m["windowsUpdate"])
	r.PreStep = flattenPatchDeploymentPatchConfigPreStep(c, m["preStep"])
	r.PostStep = flattenPatchDeploymentPatchConfigPostStep(c, m["postStep"])

	return r
}

// expandPatchDeploymentPatchConfigAptMap expands the contents of PatchDeploymentPatchConfigApt into a JSON
// request object.
func expandPatchDeploymentPatchConfigAptMap(c *Client, f map[string]PatchDeploymentPatchConfigApt) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigApt(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigAptSlice expands the contents of PatchDeploymentPatchConfigApt into a JSON
// request object.
func expandPatchDeploymentPatchConfigAptSlice(c *Client, f []PatchDeploymentPatchConfigApt) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigApt(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigAptMap flattens the contents of PatchDeploymentPatchConfigApt from a JSON
// response object.
func flattenPatchDeploymentPatchConfigAptMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigApt {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigApt{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigApt{}
	}

	items := make(map[string]PatchDeploymentPatchConfigApt)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigApt(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigAptSlice flattens the contents of PatchDeploymentPatchConfigApt from a JSON
// response object.
func flattenPatchDeploymentPatchConfigAptSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigApt {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigApt{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigApt{}
	}

	items := make([]PatchDeploymentPatchConfigApt, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigApt(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigApt expands an instance of PatchDeploymentPatchConfigApt into a JSON
// request object.
func expandPatchDeploymentPatchConfigApt(c *Client, f *PatchDeploymentPatchConfigApt) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Type; !dcl.IsEmptyValueIndirect(v) {
		m["type"] = v
	}
	if v := f.Excludes; v != nil {
		m["excludes"] = v
	}
	if v := f.ExclusivePackages; v != nil {
		m["exclusivePackages"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigApt flattens an instance of PatchDeploymentPatchConfigApt from a JSON
// response object.
func flattenPatchDeploymentPatchConfigApt(c *Client, i interface{}) *PatchDeploymentPatchConfigApt {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigApt{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigApt
	}
	r.Type = flattenPatchDeploymentPatchConfigAptTypeEnum(m["type"])
	r.Excludes = dcl.FlattenStringSlice(m["excludes"])
	r.ExclusivePackages = dcl.FlattenStringSlice(m["exclusivePackages"])

	return r
}

// expandPatchDeploymentPatchConfigYumMap expands the contents of PatchDeploymentPatchConfigYum into a JSON
// request object.
func expandPatchDeploymentPatchConfigYumMap(c *Client, f map[string]PatchDeploymentPatchConfigYum) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigYum(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigYumSlice expands the contents of PatchDeploymentPatchConfigYum into a JSON
// request object.
func expandPatchDeploymentPatchConfigYumSlice(c *Client, f []PatchDeploymentPatchConfigYum) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigYum(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigYumMap flattens the contents of PatchDeploymentPatchConfigYum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigYumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigYum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigYum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigYum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigYum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigYum(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigYumSlice flattens the contents of PatchDeploymentPatchConfigYum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigYumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigYum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigYum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigYum{}
	}

	items := make([]PatchDeploymentPatchConfigYum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigYum(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigYum expands an instance of PatchDeploymentPatchConfigYum into a JSON
// request object.
func expandPatchDeploymentPatchConfigYum(c *Client, f *PatchDeploymentPatchConfigYum) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Security; !dcl.IsEmptyValueIndirect(v) {
		m["security"] = v
	}
	if v := f.Minimal; !dcl.IsEmptyValueIndirect(v) {
		m["minimal"] = v
	}
	if v := f.Excludes; v != nil {
		m["excludes"] = v
	}
	if v := f.ExclusivePackages; v != nil {
		m["exclusivePackages"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigYum flattens an instance of PatchDeploymentPatchConfigYum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigYum(c *Client, i interface{}) *PatchDeploymentPatchConfigYum {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigYum{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigYum
	}
	r.Security = dcl.FlattenBool(m["security"])
	r.Minimal = dcl.FlattenBool(m["minimal"])
	r.Excludes = dcl.FlattenStringSlice(m["excludes"])
	r.ExclusivePackages = dcl.FlattenStringSlice(m["exclusivePackages"])

	return r
}

// expandPatchDeploymentPatchConfigGooMap expands the contents of PatchDeploymentPatchConfigGoo into a JSON
// request object.
func expandPatchDeploymentPatchConfigGooMap(c *Client, f map[string]PatchDeploymentPatchConfigGoo) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigGoo(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigGooSlice expands the contents of PatchDeploymentPatchConfigGoo into a JSON
// request object.
func expandPatchDeploymentPatchConfigGooSlice(c *Client, f []PatchDeploymentPatchConfigGoo) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigGoo(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigGooMap flattens the contents of PatchDeploymentPatchConfigGoo from a JSON
// response object.
func flattenPatchDeploymentPatchConfigGooMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigGoo {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigGoo{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigGoo{}
	}

	items := make(map[string]PatchDeploymentPatchConfigGoo)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigGoo(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigGooSlice flattens the contents of PatchDeploymentPatchConfigGoo from a JSON
// response object.
func flattenPatchDeploymentPatchConfigGooSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigGoo {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigGoo{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigGoo{}
	}

	items := make([]PatchDeploymentPatchConfigGoo, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigGoo(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigGoo expands an instance of PatchDeploymentPatchConfigGoo into a JSON
// request object.
func expandPatchDeploymentPatchConfigGoo(c *Client, f *PatchDeploymentPatchConfigGoo) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenPatchDeploymentPatchConfigGoo flattens an instance of PatchDeploymentPatchConfigGoo from a JSON
// response object.
func flattenPatchDeploymentPatchConfigGoo(c *Client, i interface{}) *PatchDeploymentPatchConfigGoo {
	_, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigGoo{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigGoo
	}

	return r
}

// expandPatchDeploymentPatchConfigZypperMap expands the contents of PatchDeploymentPatchConfigZypper into a JSON
// request object.
func expandPatchDeploymentPatchConfigZypperMap(c *Client, f map[string]PatchDeploymentPatchConfigZypper) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigZypper(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigZypperSlice expands the contents of PatchDeploymentPatchConfigZypper into a JSON
// request object.
func expandPatchDeploymentPatchConfigZypperSlice(c *Client, f []PatchDeploymentPatchConfigZypper) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigZypper(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigZypperMap flattens the contents of PatchDeploymentPatchConfigZypper from a JSON
// response object.
func flattenPatchDeploymentPatchConfigZypperMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigZypper {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigZypper{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigZypper{}
	}

	items := make(map[string]PatchDeploymentPatchConfigZypper)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigZypper(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigZypperSlice flattens the contents of PatchDeploymentPatchConfigZypper from a JSON
// response object.
func flattenPatchDeploymentPatchConfigZypperSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigZypper {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigZypper{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigZypper{}
	}

	items := make([]PatchDeploymentPatchConfigZypper, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigZypper(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigZypper expands an instance of PatchDeploymentPatchConfigZypper into a JSON
// request object.
func expandPatchDeploymentPatchConfigZypper(c *Client, f *PatchDeploymentPatchConfigZypper) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.WithOptional; !dcl.IsEmptyValueIndirect(v) {
		m["withOptional"] = v
	}
	if v := f.WithUpdate; !dcl.IsEmptyValueIndirect(v) {
		m["withUpdate"] = v
	}
	if v := f.Categories; v != nil {
		m["categories"] = v
	}
	if v := f.Severities; v != nil {
		m["severities"] = v
	}
	if v := f.Excludes; v != nil {
		m["excludes"] = v
	}
	if v := f.ExclusivePatches; v != nil {
		m["exclusivePatches"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigZypper flattens an instance of PatchDeploymentPatchConfigZypper from a JSON
// response object.
func flattenPatchDeploymentPatchConfigZypper(c *Client, i interface{}) *PatchDeploymentPatchConfigZypper {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigZypper{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigZypper
	}
	r.WithOptional = dcl.FlattenBool(m["withOptional"])
	r.WithUpdate = dcl.FlattenBool(m["withUpdate"])
	r.Categories = dcl.FlattenStringSlice(m["categories"])
	r.Severities = dcl.FlattenStringSlice(m["severities"])
	r.Excludes = dcl.FlattenStringSlice(m["excludes"])
	r.ExclusivePatches = dcl.FlattenStringSlice(m["exclusivePatches"])

	return r
}

// expandPatchDeploymentPatchConfigWindowsUpdateMap expands the contents of PatchDeploymentPatchConfigWindowsUpdate into a JSON
// request object.
func expandPatchDeploymentPatchConfigWindowsUpdateMap(c *Client, f map[string]PatchDeploymentPatchConfigWindowsUpdate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigWindowsUpdate(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigWindowsUpdateSlice expands the contents of PatchDeploymentPatchConfigWindowsUpdate into a JSON
// request object.
func expandPatchDeploymentPatchConfigWindowsUpdateSlice(c *Client, f []PatchDeploymentPatchConfigWindowsUpdate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigWindowsUpdate(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigWindowsUpdateMap flattens the contents of PatchDeploymentPatchConfigWindowsUpdate from a JSON
// response object.
func flattenPatchDeploymentPatchConfigWindowsUpdateMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigWindowsUpdate {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigWindowsUpdate{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigWindowsUpdate{}
	}

	items := make(map[string]PatchDeploymentPatchConfigWindowsUpdate)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigWindowsUpdate(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigWindowsUpdateSlice flattens the contents of PatchDeploymentPatchConfigWindowsUpdate from a JSON
// response object.
func flattenPatchDeploymentPatchConfigWindowsUpdateSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigWindowsUpdate {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigWindowsUpdate{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigWindowsUpdate{}
	}

	items := make([]PatchDeploymentPatchConfigWindowsUpdate, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigWindowsUpdate(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigWindowsUpdate expands an instance of PatchDeploymentPatchConfigWindowsUpdate into a JSON
// request object.
func expandPatchDeploymentPatchConfigWindowsUpdate(c *Client, f *PatchDeploymentPatchConfigWindowsUpdate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Classifications; v != nil {
		m["classifications"] = v
	}
	if v := f.Excludes; v != nil {
		m["excludes"] = v
	}
	if v := f.ExclusivePatches; v != nil {
		m["exclusivePatches"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigWindowsUpdate flattens an instance of PatchDeploymentPatchConfigWindowsUpdate from a JSON
// response object.
func flattenPatchDeploymentPatchConfigWindowsUpdate(c *Client, i interface{}) *PatchDeploymentPatchConfigWindowsUpdate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigWindowsUpdate{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigWindowsUpdate
	}
	r.Classifications = flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnumSlice(c, m["classifications"])
	r.Excludes = dcl.FlattenStringSlice(m["excludes"])
	r.ExclusivePatches = dcl.FlattenStringSlice(m["exclusivePatches"])

	return r
}

// expandPatchDeploymentPatchConfigPreStepMap expands the contents of PatchDeploymentPatchConfigPreStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepMap(c *Client, f map[string]PatchDeploymentPatchConfigPreStep) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStep(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPreStepSlice expands the contents of PatchDeploymentPatchConfigPreStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepSlice(c *Client, f []PatchDeploymentPatchConfigPreStep) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStep(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPreStepMap flattens the contents of PatchDeploymentPatchConfigPreStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStep {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStep{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStep{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStep)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStep(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepSlice flattens the contents of PatchDeploymentPatchConfigPreStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStep {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStep{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStep{}
	}

	items := make([]PatchDeploymentPatchConfigPreStep, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStep(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPreStep expands an instance of PatchDeploymentPatchConfigPreStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStep(c *Client, f *PatchDeploymentPatchConfigPreStep) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, f.LinuxExecStepConfig); err != nil {
		return nil, fmt.Errorf("error expanding LinuxExecStepConfig into linuxExecStepConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linuxExecStepConfig"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, f.WindowsExecStepConfig); err != nil {
		return nil, fmt.Errorf("error expanding WindowsExecStepConfig into windowsExecStepConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["windowsExecStepConfig"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPreStep flattens an instance of PatchDeploymentPatchConfigPreStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStep(c *Client, i interface{}) *PatchDeploymentPatchConfigPreStep {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPreStep{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPreStep
	}
	r.LinuxExecStepConfig = flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, m["linuxExecStepConfig"])
	r.WindowsExecStepConfig = flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, m["windowsExecStepConfig"])

	return r
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigMap expands the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigMap(c *Client, f map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice expands the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice(c *Client, f []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigMap flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfig expands an instance of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c *Client, f *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocalPath; !dcl.IsEmptyValueIndirect(v) {
		m["localPath"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, f.GcsObject); err != nil {
		return nil, fmt.Errorf("error expanding GcsObject into gcsObject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gcsObject"] = v
	}
	if v := f.AllowedSuccessCodes; v != nil {
		m["allowedSuccessCodes"] = v
	}
	if v := f.Interpreter; !dcl.IsEmptyValueIndirect(v) {
		m["interpreter"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfig flattens an instance of PatchDeploymentPatchConfigPreStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfig(c *Client, i interface{}) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfig
	}
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.GcsObject = flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, m["gcsObject"])
	r.AllowedSuccessCodes = dcl.FlattenIntSlice(m["allowedSuccessCodes"])
	r.Interpreter = flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum(m["interpreter"])

	return r
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectMap expands the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectMap(c *Client, f map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice expands the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice(c *Client, f []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectMap flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject expands an instance of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c *Client, f *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Bucket; !dcl.IsEmptyValueIndirect(v) {
		m["bucket"] = v
	}
	if v := f.Object; !dcl.IsEmptyValueIndirect(v) {
		m["object"] = v
	}
	if v := f.GenerationNumber; !dcl.IsEmptyValueIndirect(v) {
		m["generationNumber"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject flattens an instance of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject(c *Client, i interface{}) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.GenerationNumber = dcl.FlattenInteger(m["generationNumber"])

	return r
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigMap expands the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigMap(c *Client, f map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice expands the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice(c *Client, f []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigMap flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfig expands an instance of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c *Client, f *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocalPath; !dcl.IsEmptyValueIndirect(v) {
		m["localPath"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, f.GcsObject); err != nil {
		return nil, fmt.Errorf("error expanding GcsObject into gcsObject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gcsObject"] = v
	}
	if v := f.AllowedSuccessCodes; v != nil {
		m["allowedSuccessCodes"] = v
	}
	if v := f.Interpreter; !dcl.IsEmptyValueIndirect(v) {
		m["interpreter"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfig flattens an instance of PatchDeploymentPatchConfigPreStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfig(c *Client, i interface{}) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfig
	}
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.GcsObject = flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, m["gcsObject"])
	r.AllowedSuccessCodes = dcl.FlattenIntSlice(m["allowedSuccessCodes"])
	r.Interpreter = flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum(m["interpreter"])

	return r
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectMap expands the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectMap(c *Client, f map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice expands the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice(c *Client, f []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectMap flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject expands an instance of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c *Client, f *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Bucket; !dcl.IsEmptyValueIndirect(v) {
		m["bucket"] = v
	}
	if v := f.Object; !dcl.IsEmptyValueIndirect(v) {
		m["object"] = v
	}
	if v := f.GenerationNumber; !dcl.IsEmptyValueIndirect(v) {
		m["generationNumber"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject flattens an instance of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject(c *Client, i interface{}) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.GenerationNumber = dcl.FlattenInteger(m["generationNumber"])

	return r
}

// expandPatchDeploymentPatchConfigPostStepMap expands the contents of PatchDeploymentPatchConfigPostStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepMap(c *Client, f map[string]PatchDeploymentPatchConfigPostStep) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStep(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPostStepSlice expands the contents of PatchDeploymentPatchConfigPostStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepSlice(c *Client, f []PatchDeploymentPatchConfigPostStep) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStep(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPostStepMap flattens the contents of PatchDeploymentPatchConfigPostStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStep {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStep{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStep{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStep)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStep(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepSlice flattens the contents of PatchDeploymentPatchConfigPostStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStep {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStep{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStep{}
	}

	items := make([]PatchDeploymentPatchConfigPostStep, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStep(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPostStep expands an instance of PatchDeploymentPatchConfigPostStep into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStep(c *Client, f *PatchDeploymentPatchConfigPostStep) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, f.LinuxExecStepConfig); err != nil {
		return nil, fmt.Errorf("error expanding LinuxExecStepConfig into linuxExecStepConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linuxExecStepConfig"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, f.WindowsExecStepConfig); err != nil {
		return nil, fmt.Errorf("error expanding WindowsExecStepConfig into windowsExecStepConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["windowsExecStepConfig"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPostStep flattens an instance of PatchDeploymentPatchConfigPostStep from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStep(c *Client, i interface{}) *PatchDeploymentPatchConfigPostStep {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPostStep{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPostStep
	}
	r.LinuxExecStepConfig = flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, m["linuxExecStepConfig"])
	r.WindowsExecStepConfig = flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, m["windowsExecStepConfig"])

	return r
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigMap expands the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigMap(c *Client, f map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice expands the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice(c *Client, f []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigMap flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfig expands an instance of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c *Client, f *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocalPath; !dcl.IsEmptyValueIndirect(v) {
		m["localPath"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, f.GcsObject); err != nil {
		return nil, fmt.Errorf("error expanding GcsObject into gcsObject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gcsObject"] = v
	}
	if v := f.AllowedSuccessCodes; v != nil {
		m["allowedSuccessCodes"] = v
	}
	if v := f.Interpreter; !dcl.IsEmptyValueIndirect(v) {
		m["interpreter"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfig flattens an instance of PatchDeploymentPatchConfigPostStepLinuxExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfig(c *Client, i interface{}) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfig
	}
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.GcsObject = flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, m["gcsObject"])
	r.AllowedSuccessCodes = dcl.FlattenIntSlice(m["allowedSuccessCodes"])
	r.Interpreter = flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum(m["interpreter"])

	return r
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectMap expands the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectMap(c *Client, f map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice expands the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice(c *Client, f []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectMap flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject expands an instance of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c *Client, f *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Bucket; !dcl.IsEmptyValueIndirect(v) {
		m["bucket"] = v
	}
	if v := f.Object; !dcl.IsEmptyValueIndirect(v) {
		m["object"] = v
	}
	if v := f.GenerationNumber; !dcl.IsEmptyValueIndirect(v) {
		m["generationNumber"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject flattens an instance of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject(c *Client, i interface{}) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.GenerationNumber = dcl.FlattenInteger(m["generationNumber"])

	return r
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigMap expands the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigMap(c *Client, f map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice expands the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice(c *Client, f []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigMap flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfig expands an instance of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c *Client, f *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocalPath; !dcl.IsEmptyValueIndirect(v) {
		m["localPath"] = v
	}
	if v, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, f.GcsObject); err != nil {
		return nil, fmt.Errorf("error expanding GcsObject into gcsObject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gcsObject"] = v
	}
	if v := f.AllowedSuccessCodes; v != nil {
		m["allowedSuccessCodes"] = v
	}
	if v := f.Interpreter; !dcl.IsEmptyValueIndirect(v) {
		m["interpreter"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfig flattens an instance of PatchDeploymentPatchConfigPostStepWindowsExecStepConfig from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfig(c *Client, i interface{}) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfig
	}
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.GcsObject = flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, m["gcsObject"])
	r.AllowedSuccessCodes = dcl.FlattenIntSlice(m["allowedSuccessCodes"])
	r.Interpreter = flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum(m["interpreter"])

	return r
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectMap expands the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectMap(c *Client, f map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice expands the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice(c *Client, f []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectMap flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject expands an instance of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject into a JSON
// request object.
func expandPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c *Client, f *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Bucket; !dcl.IsEmptyValueIndirect(v) {
		m["bucket"] = v
	}
	if v := f.Object; !dcl.IsEmptyValueIndirect(v) {
		m["object"] = v
	}
	if v := f.GenerationNumber; !dcl.IsEmptyValueIndirect(v) {
		m["generationNumber"] = v
	}

	return m, nil
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject flattens an instance of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject(c *Client, i interface{}) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.GenerationNumber = dcl.FlattenInteger(m["generationNumber"])

	return r
}

// expandPatchDeploymentOneTimeScheduleMap expands the contents of PatchDeploymentOneTimeSchedule into a JSON
// request object.
func expandPatchDeploymentOneTimeScheduleMap(c *Client, f map[string]PatchDeploymentOneTimeSchedule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentOneTimeSchedule(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentOneTimeScheduleSlice expands the contents of PatchDeploymentOneTimeSchedule into a JSON
// request object.
func expandPatchDeploymentOneTimeScheduleSlice(c *Client, f []PatchDeploymentOneTimeSchedule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentOneTimeSchedule(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentOneTimeScheduleMap flattens the contents of PatchDeploymentOneTimeSchedule from a JSON
// response object.
func flattenPatchDeploymentOneTimeScheduleMap(c *Client, i interface{}) map[string]PatchDeploymentOneTimeSchedule {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentOneTimeSchedule{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentOneTimeSchedule{}
	}

	items := make(map[string]PatchDeploymentOneTimeSchedule)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentOneTimeSchedule(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentOneTimeScheduleSlice flattens the contents of PatchDeploymentOneTimeSchedule from a JSON
// response object.
func flattenPatchDeploymentOneTimeScheduleSlice(c *Client, i interface{}) []PatchDeploymentOneTimeSchedule {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentOneTimeSchedule{}
	}

	if len(a) == 0 {
		return []PatchDeploymentOneTimeSchedule{}
	}

	items := make([]PatchDeploymentOneTimeSchedule, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentOneTimeSchedule(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentOneTimeSchedule expands an instance of PatchDeploymentOneTimeSchedule into a JSON
// request object.
func expandPatchDeploymentOneTimeSchedule(c *Client, f *PatchDeploymentOneTimeSchedule) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ExecuteTime; !dcl.IsEmptyValueIndirect(v) {
		m["executeTime"] = v
	}

	return m, nil
}

// flattenPatchDeploymentOneTimeSchedule flattens an instance of PatchDeploymentOneTimeSchedule from a JSON
// response object.
func flattenPatchDeploymentOneTimeSchedule(c *Client, i interface{}) *PatchDeploymentOneTimeSchedule {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentOneTimeSchedule{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentOneTimeSchedule
	}
	r.ExecuteTime = dcl.FlattenString(m["executeTime"])

	return r
}

// expandPatchDeploymentRecurringScheduleMap expands the contents of PatchDeploymentRecurringSchedule into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMap(c *Client, f map[string]PatchDeploymentRecurringSchedule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringSchedule(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleSlice expands the contents of PatchDeploymentRecurringSchedule into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleSlice(c *Client, f []PatchDeploymentRecurringSchedule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringSchedule(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleMap flattens the contents of PatchDeploymentRecurringSchedule from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringSchedule {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringSchedule{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringSchedule{}
	}

	items := make(map[string]PatchDeploymentRecurringSchedule)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringSchedule(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleSlice flattens the contents of PatchDeploymentRecurringSchedule from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleSlice(c *Client, i interface{}) []PatchDeploymentRecurringSchedule {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringSchedule{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringSchedule{}
	}

	items := make([]PatchDeploymentRecurringSchedule, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringSchedule(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringSchedule expands an instance of PatchDeploymentRecurringSchedule into a JSON
// request object.
func expandPatchDeploymentRecurringSchedule(c *Client, f *PatchDeploymentRecurringSchedule) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPatchDeploymentRecurringScheduleTimeZone(c, f.TimeZone); err != nil {
		return nil, fmt.Errorf("error expanding TimeZone into timeZone: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["timeZone"] = v
	}
	if v := f.StartTime; !dcl.IsEmptyValueIndirect(v) {
		m["startTime"] = v
	}
	if v := f.EndTime; !dcl.IsEmptyValueIndirect(v) {
		m["endTime"] = v
	}
	if v, err := expandPatchDeploymentRecurringScheduleTimeOfDay(c, f.TimeOfDay); err != nil {
		return nil, fmt.Errorf("error expanding TimeOfDay into timeOfDay: %w", err)
	} else if v != nil {
		m["timeOfDay"] = v
	}
	if v := f.Frequency; !dcl.IsEmptyValueIndirect(v) {
		m["frequency"] = v
	}
	if v, err := expandPatchDeploymentRecurringScheduleWeekly(c, f.Weekly); err != nil {
		return nil, fmt.Errorf("error expanding Weekly into weekly: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["weekly"] = v
	}
	if v, err := expandPatchDeploymentRecurringScheduleMonthly(c, f.Monthly); err != nil {
		return nil, fmt.Errorf("error expanding Monthly into monthly: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["monthly"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringSchedule flattens an instance of PatchDeploymentRecurringSchedule from a JSON
// response object.
func flattenPatchDeploymentRecurringSchedule(c *Client, i interface{}) *PatchDeploymentRecurringSchedule {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringSchedule{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringSchedule
	}
	r.TimeZone = flattenPatchDeploymentRecurringScheduleTimeZone(c, m["timeZone"])
	r.StartTime = dcl.FlattenString(m["startTime"])
	r.EndTime = dcl.FlattenString(m["endTime"])
	r.TimeOfDay = flattenPatchDeploymentRecurringScheduleTimeOfDay(c, m["timeOfDay"])
	r.Frequency = flattenPatchDeploymentRecurringScheduleFrequencyEnum(m["frequency"])
	r.Weekly = flattenPatchDeploymentRecurringScheduleWeekly(c, m["weekly"])
	r.Monthly = flattenPatchDeploymentRecurringScheduleMonthly(c, m["monthly"])
	r.LastExecuteTime = dcl.FlattenString(m["lastExecuteTime"])
	r.NextExecuteTime = dcl.FlattenString(m["nextExecuteTime"])

	return r
}

// expandPatchDeploymentRecurringScheduleTimeZoneMap expands the contents of PatchDeploymentRecurringScheduleTimeZone into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeZoneMap(c *Client, f map[string]PatchDeploymentRecurringScheduleTimeZone) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleTimeZone(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleTimeZoneSlice expands the contents of PatchDeploymentRecurringScheduleTimeZone into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeZoneSlice(c *Client, f []PatchDeploymentRecurringScheduleTimeZone) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleTimeZone(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleTimeZoneMap flattens the contents of PatchDeploymentRecurringScheduleTimeZone from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeZoneMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleTimeZone {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleTimeZone{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleTimeZone{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleTimeZone)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleTimeZone(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleTimeZoneSlice flattens the contents of PatchDeploymentRecurringScheduleTimeZone from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeZoneSlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleTimeZone {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleTimeZone{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleTimeZone{}
	}

	items := make([]PatchDeploymentRecurringScheduleTimeZone, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleTimeZone(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringScheduleTimeZone expands an instance of PatchDeploymentRecurringScheduleTimeZone into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeZone(c *Client, f *PatchDeploymentRecurringScheduleTimeZone) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v := f.Version; !dcl.IsEmptyValueIndirect(v) {
		m["version"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringScheduleTimeZone flattens an instance of PatchDeploymentRecurringScheduleTimeZone from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeZone(c *Client, i interface{}) *PatchDeploymentRecurringScheduleTimeZone {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringScheduleTimeZone{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringScheduleTimeZone
	}
	r.Id = dcl.FlattenString(m["id"])
	r.Version = dcl.FlattenString(m["version"])

	return r
}

// expandPatchDeploymentRecurringScheduleTimeOfDayMap expands the contents of PatchDeploymentRecurringScheduleTimeOfDay into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeOfDayMap(c *Client, f map[string]PatchDeploymentRecurringScheduleTimeOfDay) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleTimeOfDay(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleTimeOfDaySlice expands the contents of PatchDeploymentRecurringScheduleTimeOfDay into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeOfDaySlice(c *Client, f []PatchDeploymentRecurringScheduleTimeOfDay) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleTimeOfDay(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleTimeOfDayMap flattens the contents of PatchDeploymentRecurringScheduleTimeOfDay from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeOfDayMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleTimeOfDay {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleTimeOfDay{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleTimeOfDay{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleTimeOfDay)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleTimeOfDay(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleTimeOfDaySlice flattens the contents of PatchDeploymentRecurringScheduleTimeOfDay from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeOfDaySlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleTimeOfDay {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleTimeOfDay{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleTimeOfDay{}
	}

	items := make([]PatchDeploymentRecurringScheduleTimeOfDay, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleTimeOfDay(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringScheduleTimeOfDay expands an instance of PatchDeploymentRecurringScheduleTimeOfDay into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleTimeOfDay(c *Client, f *PatchDeploymentRecurringScheduleTimeOfDay) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Hours; !dcl.IsEmptyValueIndirect(v) {
		m["hours"] = v
	}
	if v := f.Minutes; !dcl.IsEmptyValueIndirect(v) {
		m["minutes"] = v
	}
	if v := f.Seconds; !dcl.IsEmptyValueIndirect(v) {
		m["seconds"] = v
	}
	if v := f.Nanos; !dcl.IsEmptyValueIndirect(v) {
		m["nanos"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringScheduleTimeOfDay flattens an instance of PatchDeploymentRecurringScheduleTimeOfDay from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleTimeOfDay(c *Client, i interface{}) *PatchDeploymentRecurringScheduleTimeOfDay {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringScheduleTimeOfDay{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringScheduleTimeOfDay
	}
	r.Hours = dcl.FlattenInteger(m["hours"])
	r.Minutes = dcl.FlattenInteger(m["minutes"])
	r.Seconds = dcl.FlattenInteger(m["seconds"])
	r.Nanos = dcl.FlattenInteger(m["nanos"])

	return r
}

// expandPatchDeploymentRecurringScheduleWeeklyMap expands the contents of PatchDeploymentRecurringScheduleWeekly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleWeeklyMap(c *Client, f map[string]PatchDeploymentRecurringScheduleWeekly) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleWeekly(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleWeeklySlice expands the contents of PatchDeploymentRecurringScheduleWeekly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleWeeklySlice(c *Client, f []PatchDeploymentRecurringScheduleWeekly) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleWeekly(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleWeeklyMap flattens the contents of PatchDeploymentRecurringScheduleWeekly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleWeeklyMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleWeekly {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleWeekly{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleWeekly{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleWeekly)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleWeekly(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleWeeklySlice flattens the contents of PatchDeploymentRecurringScheduleWeekly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleWeeklySlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleWeekly {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleWeekly{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleWeekly{}
	}

	items := make([]PatchDeploymentRecurringScheduleWeekly, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleWeekly(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringScheduleWeekly expands an instance of PatchDeploymentRecurringScheduleWeekly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleWeekly(c *Client, f *PatchDeploymentRecurringScheduleWeekly) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DayOfWeek; !dcl.IsEmptyValueIndirect(v) {
		m["dayOfWeek"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringScheduleWeekly flattens an instance of PatchDeploymentRecurringScheduleWeekly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleWeekly(c *Client, i interface{}) *PatchDeploymentRecurringScheduleWeekly {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringScheduleWeekly{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringScheduleWeekly
	}
	r.DayOfWeek = flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum(m["dayOfWeek"])

	return r
}

// expandPatchDeploymentRecurringScheduleMonthlyMap expands the contents of PatchDeploymentRecurringScheduleMonthly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthlyMap(c *Client, f map[string]PatchDeploymentRecurringScheduleMonthly) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleMonthly(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleMonthlySlice expands the contents of PatchDeploymentRecurringScheduleMonthly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthlySlice(c *Client, f []PatchDeploymentRecurringScheduleMonthly) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleMonthly(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleMonthlyMap flattens the contents of PatchDeploymentRecurringScheduleMonthly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleMonthly {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleMonthly{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleMonthly{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleMonthly)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleMonthly(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleMonthlySlice flattens the contents of PatchDeploymentRecurringScheduleMonthly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlySlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleMonthly {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleMonthly{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleMonthly{}
	}

	items := make([]PatchDeploymentRecurringScheduleMonthly, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleMonthly(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringScheduleMonthly expands an instance of PatchDeploymentRecurringScheduleMonthly into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthly(c *Client, f *PatchDeploymentRecurringScheduleMonthly) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, f.WeekDayOfMonth); err != nil {
		return nil, fmt.Errorf("error expanding WeekDayOfMonth into weekDayOfMonth: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["weekDayOfMonth"] = v
	}
	if v := f.MonthDay; !dcl.IsEmptyValueIndirect(v) {
		m["monthDay"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringScheduleMonthly flattens an instance of PatchDeploymentRecurringScheduleMonthly from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthly(c *Client, i interface{}) *PatchDeploymentRecurringScheduleMonthly {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringScheduleMonthly{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringScheduleMonthly
	}
	r.WeekDayOfMonth = flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, m["weekDayOfMonth"])
	r.MonthDay = dcl.FlattenInteger(m["monthDay"])

	return r
}

// expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthMap expands the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthMap(c *Client, f map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice expands the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice(c *Client, f []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthMap flattens the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice flattens the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthSlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}

	items := make([]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth expands an instance of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth into a JSON
// request object.
func expandPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c *Client, f *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.WeekOrdinal; !dcl.IsEmptyValueIndirect(v) {
		m["weekOrdinal"] = v
	}
	if v := f.DayOfWeek; !dcl.IsEmptyValueIndirect(v) {
		m["dayOfWeek"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth flattens an instance of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth(c *Client, i interface{}) *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth
	}
	r.WeekOrdinal = dcl.FlattenInteger(m["weekOrdinal"])
	r.DayOfWeek = flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum(m["dayOfWeek"])

	return r
}

// expandPatchDeploymentRolloutMap expands the contents of PatchDeploymentRollout into a JSON
// request object.
func expandPatchDeploymentRolloutMap(c *Client, f map[string]PatchDeploymentRollout) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRollout(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRolloutSlice expands the contents of PatchDeploymentRollout into a JSON
// request object.
func expandPatchDeploymentRolloutSlice(c *Client, f []PatchDeploymentRollout) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRollout(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRolloutMap flattens the contents of PatchDeploymentRollout from a JSON
// response object.
func flattenPatchDeploymentRolloutMap(c *Client, i interface{}) map[string]PatchDeploymentRollout {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRollout{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRollout{}
	}

	items := make(map[string]PatchDeploymentRollout)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRollout(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRolloutSlice flattens the contents of PatchDeploymentRollout from a JSON
// response object.
func flattenPatchDeploymentRolloutSlice(c *Client, i interface{}) []PatchDeploymentRollout {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRollout{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRollout{}
	}

	items := make([]PatchDeploymentRollout, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRollout(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRollout expands an instance of PatchDeploymentRollout into a JSON
// request object.
func expandPatchDeploymentRollout(c *Client, f *PatchDeploymentRollout) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Mode; !dcl.IsEmptyValueIndirect(v) {
		m["mode"] = v
	}
	if v, err := expandPatchDeploymentRolloutDisruptionBudget(c, f.DisruptionBudget); err != nil {
		return nil, fmt.Errorf("error expanding DisruptionBudget into disruptionBudget: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["disruptionBudget"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRollout flattens an instance of PatchDeploymentRollout from a JSON
// response object.
func flattenPatchDeploymentRollout(c *Client, i interface{}) *PatchDeploymentRollout {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRollout{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRollout
	}
	r.Mode = flattenPatchDeploymentRolloutModeEnum(m["mode"])
	r.DisruptionBudget = flattenPatchDeploymentRolloutDisruptionBudget(c, m["disruptionBudget"])

	return r
}

// expandPatchDeploymentRolloutDisruptionBudgetMap expands the contents of PatchDeploymentRolloutDisruptionBudget into a JSON
// request object.
func expandPatchDeploymentRolloutDisruptionBudgetMap(c *Client, f map[string]PatchDeploymentRolloutDisruptionBudget) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPatchDeploymentRolloutDisruptionBudget(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPatchDeploymentRolloutDisruptionBudgetSlice expands the contents of PatchDeploymentRolloutDisruptionBudget into a JSON
// request object.
func expandPatchDeploymentRolloutDisruptionBudgetSlice(c *Client, f []PatchDeploymentRolloutDisruptionBudget) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPatchDeploymentRolloutDisruptionBudget(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPatchDeploymentRolloutDisruptionBudgetMap flattens the contents of PatchDeploymentRolloutDisruptionBudget from a JSON
// response object.
func flattenPatchDeploymentRolloutDisruptionBudgetMap(c *Client, i interface{}) map[string]PatchDeploymentRolloutDisruptionBudget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRolloutDisruptionBudget{}
	}

	items := make(map[string]PatchDeploymentRolloutDisruptionBudget)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRolloutDisruptionBudget(c, item.(map[string]interface{}))
	}

	return items
}

// flattenPatchDeploymentRolloutDisruptionBudgetSlice flattens the contents of PatchDeploymentRolloutDisruptionBudget from a JSON
// response object.
func flattenPatchDeploymentRolloutDisruptionBudgetSlice(c *Client, i interface{}) []PatchDeploymentRolloutDisruptionBudget {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRolloutDisruptionBudget{}
	}

	items := make([]PatchDeploymentRolloutDisruptionBudget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRolloutDisruptionBudget(c, item.(map[string]interface{})))
	}

	return items
}

// expandPatchDeploymentRolloutDisruptionBudget expands an instance of PatchDeploymentRolloutDisruptionBudget into a JSON
// request object.
func expandPatchDeploymentRolloutDisruptionBudget(c *Client, f *PatchDeploymentRolloutDisruptionBudget) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Fixed; !dcl.IsEmptyValueIndirect(v) {
		m["fixed"] = v
	}
	if v := f.Percent; !dcl.IsEmptyValueIndirect(v) {
		m["percent"] = v
	}

	return m, nil
}

// flattenPatchDeploymentRolloutDisruptionBudget flattens an instance of PatchDeploymentRolloutDisruptionBudget from a JSON
// response object.
func flattenPatchDeploymentRolloutDisruptionBudget(c *Client, i interface{}) *PatchDeploymentRolloutDisruptionBudget {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PatchDeploymentRolloutDisruptionBudget{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPatchDeploymentRolloutDisruptionBudget
	}
	r.Fixed = dcl.FlattenInteger(m["fixed"])
	r.Percent = dcl.FlattenInteger(m["percent"])

	return r
}

// flattenPatchDeploymentPatchConfigRebootConfigEnumMap flattens the contents of PatchDeploymentPatchConfigRebootConfigEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigRebootConfigEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigRebootConfigEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigRebootConfigEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigRebootConfigEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigRebootConfigEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigRebootConfigEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigRebootConfigEnumSlice flattens the contents of PatchDeploymentPatchConfigRebootConfigEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigRebootConfigEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigRebootConfigEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigRebootConfigEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigRebootConfigEnum{}
	}

	items := make([]PatchDeploymentPatchConfigRebootConfigEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigRebootConfigEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigRebootConfigEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigRebootConfigEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigRebootConfigEnum(i interface{}) *PatchDeploymentPatchConfigRebootConfigEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigRebootConfigEnumRef("")
	}

	return PatchDeploymentPatchConfigRebootConfigEnumRef(s)
}

// flattenPatchDeploymentPatchConfigAptTypeEnumMap flattens the contents of PatchDeploymentPatchConfigAptTypeEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigAptTypeEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigAptTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigAptTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigAptTypeEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigAptTypeEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigAptTypeEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigAptTypeEnumSlice flattens the contents of PatchDeploymentPatchConfigAptTypeEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigAptTypeEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigAptTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigAptTypeEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigAptTypeEnum{}
	}

	items := make([]PatchDeploymentPatchConfigAptTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigAptTypeEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigAptTypeEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigAptTypeEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigAptTypeEnum(i interface{}) *PatchDeploymentPatchConfigAptTypeEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigAptTypeEnumRef("")
	}

	return PatchDeploymentPatchConfigAptTypeEnumRef(s)
}

// flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnumMap flattens the contents of PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnumSlice flattens the contents of PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum{}
	}

	items := make([]PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigWindowsUpdateClassificationsEnum(i interface{}) *PatchDeploymentPatchConfigWindowsUpdateClassificationsEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigWindowsUpdateClassificationsEnumRef("")
	}

	return PatchDeploymentPatchConfigWindowsUpdateClassificationsEnumRef(s)
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumMap flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumSlice flattens the contents of PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum(i interface{}) *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumRef("")
	}

	return PatchDeploymentPatchConfigPreStepLinuxExecStepConfigInterpreterEnumRef(s)
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumMap flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumSlice flattens the contents of PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum{}
	}

	items := make([]PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum(i interface{}) *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumRef("")
	}

	return PatchDeploymentPatchConfigPreStepWindowsExecStepConfigInterpreterEnumRef(s)
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumMap flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumSlice flattens the contents of PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum(i interface{}) *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumRef("")
	}

	return PatchDeploymentPatchConfigPostStepLinuxExecStepConfigInterpreterEnumRef(s)
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumMap flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumMap(c *Client, i interface{}) map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum{}
	}

	items := make(map[string]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumSlice flattens the contents of PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum from a JSON
// response object.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumSlice(c *Client, i interface{}) []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum{}
	}

	items := make([]PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum with the same value as that string.
func flattenPatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum(i interface{}) *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumRef("")
	}

	return PatchDeploymentPatchConfigPostStepWindowsExecStepConfigInterpreterEnumRef(s)
}

// flattenPatchDeploymentRecurringScheduleFrequencyEnumMap flattens the contents of PatchDeploymentRecurringScheduleFrequencyEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleFrequencyEnumMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleFrequencyEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleFrequencyEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleFrequencyEnum{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleFrequencyEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleFrequencyEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleFrequencyEnumSlice flattens the contents of PatchDeploymentRecurringScheduleFrequencyEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleFrequencyEnumSlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleFrequencyEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleFrequencyEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleFrequencyEnum{}
	}

	items := make([]PatchDeploymentRecurringScheduleFrequencyEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleFrequencyEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleFrequencyEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentRecurringScheduleFrequencyEnum with the same value as that string.
func flattenPatchDeploymentRecurringScheduleFrequencyEnum(i interface{}) *PatchDeploymentRecurringScheduleFrequencyEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentRecurringScheduleFrequencyEnumRef("")
	}

	return PatchDeploymentRecurringScheduleFrequencyEnumRef(s)
}

// flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumMap flattens the contents of PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumSlice flattens the contents of PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumSlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum{}
	}

	items := make([]PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum with the same value as that string.
func flattenPatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum(i interface{}) *PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumRef("")
	}

	return PatchDeploymentRecurringScheduleWeeklyDayOfWeekEnumRef(s)
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumMap flattens the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumMap(c *Client, i interface{}) map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum{}
	}

	items := make(map[string]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumSlice flattens the contents of PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum from a JSON
// response object.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumSlice(c *Client, i interface{}) []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum{}
	}

	items := make([]PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum with the same value as that string.
func flattenPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum(i interface{}) *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumRef("")
	}

	return PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthDayOfWeekEnumRef(s)
}

// flattenPatchDeploymentRolloutModeEnumMap flattens the contents of PatchDeploymentRolloutModeEnum from a JSON
// response object.
func flattenPatchDeploymentRolloutModeEnumMap(c *Client, i interface{}) map[string]PatchDeploymentRolloutModeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PatchDeploymentRolloutModeEnum{}
	}

	if len(a) == 0 {
		return map[string]PatchDeploymentRolloutModeEnum{}
	}

	items := make(map[string]PatchDeploymentRolloutModeEnum)
	for k, item := range a {
		items[k] = *flattenPatchDeploymentRolloutModeEnum(item.(interface{}))
	}

	return items
}

// flattenPatchDeploymentRolloutModeEnumSlice flattens the contents of PatchDeploymentRolloutModeEnum from a JSON
// response object.
func flattenPatchDeploymentRolloutModeEnumSlice(c *Client, i interface{}) []PatchDeploymentRolloutModeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PatchDeploymentRolloutModeEnum{}
	}

	if len(a) == 0 {
		return []PatchDeploymentRolloutModeEnum{}
	}

	items := make([]PatchDeploymentRolloutModeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPatchDeploymentRolloutModeEnum(item.(interface{})))
	}

	return items
}

// flattenPatchDeploymentRolloutModeEnum asserts that an interface is a string, and returns a
// pointer to a *PatchDeploymentRolloutModeEnum with the same value as that string.
func flattenPatchDeploymentRolloutModeEnum(i interface{}) *PatchDeploymentRolloutModeEnum {
	s, ok := i.(string)
	if !ok {
		return PatchDeploymentRolloutModeEnumRef("")
	}

	return PatchDeploymentRolloutModeEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *PatchDeployment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalPatchDeployment(b, c)
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

type patchDeploymentDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         patchDeploymentApiOperation
}

func convertFieldDiffsToPatchDeploymentDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]patchDeploymentDiff, error) {
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
	var diffs []patchDeploymentDiff
	// For each operation name, create a patchDeploymentDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := patchDeploymentDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToPatchDeploymentApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToPatchDeploymentApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (patchDeploymentApiOperation, error) {
	switch opName {

	case "updatePatchDeploymentUpdatePatchDeploymentOperation":
		return &updatePatchDeploymentUpdatePatchDeploymentOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractPatchDeploymentFields(r *PatchDeployment) error {
	vInstanceFilter := r.InstanceFilter
	if vInstanceFilter == nil {
		// note: explicitly not the empty object.
		vInstanceFilter = &PatchDeploymentInstanceFilter{}
	}
	if err := extractPatchDeploymentInstanceFilterFields(r, vInstanceFilter); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vInstanceFilter) {
		r.InstanceFilter = vInstanceFilter
	}
	vPatchConfig := r.PatchConfig
	if vPatchConfig == nil {
		// note: explicitly not the empty object.
		vPatchConfig = &PatchDeploymentPatchConfig{}
	}
	if err := extractPatchDeploymentPatchConfigFields(r, vPatchConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPatchConfig) {
		r.PatchConfig = vPatchConfig
	}
	vOneTimeSchedule := r.OneTimeSchedule
	if vOneTimeSchedule == nil {
		// note: explicitly not the empty object.
		vOneTimeSchedule = &PatchDeploymentOneTimeSchedule{}
	}
	if err := extractPatchDeploymentOneTimeScheduleFields(r, vOneTimeSchedule); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vOneTimeSchedule) {
		r.OneTimeSchedule = vOneTimeSchedule
	}
	vRecurringSchedule := r.RecurringSchedule
	if vRecurringSchedule == nil {
		// note: explicitly not the empty object.
		vRecurringSchedule = &PatchDeploymentRecurringSchedule{}
	}
	if err := extractPatchDeploymentRecurringScheduleFields(r, vRecurringSchedule); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRecurringSchedule) {
		r.RecurringSchedule = vRecurringSchedule
	}
	vRollout := r.Rollout
	if vRollout == nil {
		// note: explicitly not the empty object.
		vRollout = &PatchDeploymentRollout{}
	}
	if err := extractPatchDeploymentRolloutFields(r, vRollout); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRollout) {
		r.Rollout = vRollout
	}
	return nil
}
func extractPatchDeploymentInstanceFilterFields(r *PatchDeployment, o *PatchDeploymentInstanceFilter) error {
	return nil
}
func extractPatchDeploymentInstanceFilterGroupLabelsFields(r *PatchDeployment, o *PatchDeploymentInstanceFilterGroupLabels) error {
	return nil
}
func extractPatchDeploymentPatchConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfig) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &PatchDeploymentPatchConfigApt{}
	}
	if err := extractPatchDeploymentPatchConfigAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &PatchDeploymentPatchConfigYum{}
	}
	if err := extractPatchDeploymentPatchConfigYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vGoo := o.Goo
	if vGoo == nil {
		// note: explicitly not the empty object.
		vGoo = &PatchDeploymentPatchConfigGoo{}
	}
	if err := extractPatchDeploymentPatchConfigGooFields(r, vGoo); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGoo) {
		o.Goo = vGoo
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &PatchDeploymentPatchConfigZypper{}
	}
	if err := extractPatchDeploymentPatchConfigZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vWindowsUpdate := o.WindowsUpdate
	if vWindowsUpdate == nil {
		// note: explicitly not the empty object.
		vWindowsUpdate = &PatchDeploymentPatchConfigWindowsUpdate{}
	}
	if err := extractPatchDeploymentPatchConfigWindowsUpdateFields(r, vWindowsUpdate); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsUpdate) {
		o.WindowsUpdate = vWindowsUpdate
	}
	vPreStep := o.PreStep
	if vPreStep == nil {
		// note: explicitly not the empty object.
		vPreStep = &PatchDeploymentPatchConfigPreStep{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepFields(r, vPreStep); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPreStep) {
		o.PreStep = vPreStep
	}
	vPostStep := o.PostStep
	if vPostStep == nil {
		// note: explicitly not the empty object.
		vPostStep = &PatchDeploymentPatchConfigPostStep{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepFields(r, vPostStep); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPostStep) {
		o.PostStep = vPostStep
	}
	return nil
}
func extractPatchDeploymentPatchConfigAptFields(r *PatchDeployment, o *PatchDeploymentPatchConfigApt) error {
	return nil
}
func extractPatchDeploymentPatchConfigYumFields(r *PatchDeployment, o *PatchDeploymentPatchConfigYum) error {
	return nil
}
func extractPatchDeploymentPatchConfigGooFields(r *PatchDeployment, o *PatchDeploymentPatchConfigGoo) error {
	return nil
}
func extractPatchDeploymentPatchConfigZypperFields(r *PatchDeployment, o *PatchDeploymentPatchConfigZypper) error {
	return nil
}
func extractPatchDeploymentPatchConfigWindowsUpdateFields(r *PatchDeployment, o *PatchDeploymentPatchConfigWindowsUpdate) error {
	return nil
}
func extractPatchDeploymentPatchConfigPreStepFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStep) error {
	vLinuxExecStepConfig := o.LinuxExecStepConfig
	if vLinuxExecStepConfig == nil {
		// note: explicitly not the empty object.
		vLinuxExecStepConfig = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigFields(r, vLinuxExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLinuxExecStepConfig) {
		o.LinuxExecStepConfig = vLinuxExecStepConfig
	}
	vWindowsExecStepConfig := o.WindowsExecStepConfig
	if vWindowsExecStepConfig == nil {
		// note: explicitly not the empty object.
		vWindowsExecStepConfig = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigFields(r, vWindowsExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsExecStepConfig) {
		o.WindowsExecStepConfig = vWindowsExecStepConfig
	}
	return nil
}
func extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) error {
	return nil
}
func extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) error {
	return nil
}
func extractPatchDeploymentPatchConfigPostStepFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStep) error {
	vLinuxExecStepConfig := o.LinuxExecStepConfig
	if vLinuxExecStepConfig == nil {
		// note: explicitly not the empty object.
		vLinuxExecStepConfig = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigFields(r, vLinuxExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLinuxExecStepConfig) {
		o.LinuxExecStepConfig = vLinuxExecStepConfig
	}
	vWindowsExecStepConfig := o.WindowsExecStepConfig
	if vWindowsExecStepConfig == nil {
		// note: explicitly not the empty object.
		vWindowsExecStepConfig = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigFields(r, vWindowsExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsExecStepConfig) {
		o.WindowsExecStepConfig = vWindowsExecStepConfig
	}
	return nil
}
func extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) error {
	return nil
}
func extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) error {
	return nil
}
func extractPatchDeploymentOneTimeScheduleFields(r *PatchDeployment, o *PatchDeploymentOneTimeSchedule) error {
	return nil
}
func extractPatchDeploymentRecurringScheduleFields(r *PatchDeployment, o *PatchDeploymentRecurringSchedule) error {
	vTimeZone := o.TimeZone
	if vTimeZone == nil {
		// note: explicitly not the empty object.
		vTimeZone = &PatchDeploymentRecurringScheduleTimeZone{}
	}
	if err := extractPatchDeploymentRecurringScheduleTimeZoneFields(r, vTimeZone); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vTimeZone) {
		o.TimeZone = vTimeZone
	}
	vTimeOfDay := o.TimeOfDay
	if vTimeOfDay == nil {
		// note: explicitly not the empty object.
		vTimeOfDay = &PatchDeploymentRecurringScheduleTimeOfDay{}
	}
	if err := extractPatchDeploymentRecurringScheduleTimeOfDayFields(r, vTimeOfDay); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vTimeOfDay) {
		o.TimeOfDay = vTimeOfDay
	}
	vWeekly := o.Weekly
	if vWeekly == nil {
		// note: explicitly not the empty object.
		vWeekly = &PatchDeploymentRecurringScheduleWeekly{}
	}
	if err := extractPatchDeploymentRecurringScheduleWeeklyFields(r, vWeekly); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWeekly) {
		o.Weekly = vWeekly
	}
	vMonthly := o.Monthly
	if vMonthly == nil {
		// note: explicitly not the empty object.
		vMonthly = &PatchDeploymentRecurringScheduleMonthly{}
	}
	if err := extractPatchDeploymentRecurringScheduleMonthlyFields(r, vMonthly); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMonthly) {
		o.Monthly = vMonthly
	}
	return nil
}
func extractPatchDeploymentRecurringScheduleTimeZoneFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleTimeZone) error {
	return nil
}
func extractPatchDeploymentRecurringScheduleTimeOfDayFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleTimeOfDay) error {
	return nil
}
func extractPatchDeploymentRecurringScheduleWeeklyFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleWeekly) error {
	return nil
}
func extractPatchDeploymentRecurringScheduleMonthlyFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleMonthly) error {
	vWeekDayOfMonth := o.WeekDayOfMonth
	if vWeekDayOfMonth == nil {
		// note: explicitly not the empty object.
		vWeekDayOfMonth = &PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}
	if err := extractPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthFields(r, vWeekDayOfMonth); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWeekDayOfMonth) {
		o.WeekDayOfMonth = vWeekDayOfMonth
	}
	return nil
}
func extractPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) error {
	return nil
}
func extractPatchDeploymentRolloutFields(r *PatchDeployment, o *PatchDeploymentRollout) error {
	vDisruptionBudget := o.DisruptionBudget
	if vDisruptionBudget == nil {
		// note: explicitly not the empty object.
		vDisruptionBudget = &PatchDeploymentRolloutDisruptionBudget{}
	}
	if err := extractPatchDeploymentRolloutDisruptionBudgetFields(r, vDisruptionBudget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDisruptionBudget) {
		o.DisruptionBudget = vDisruptionBudget
	}
	return nil
}
func extractPatchDeploymentRolloutDisruptionBudgetFields(r *PatchDeployment, o *PatchDeploymentRolloutDisruptionBudget) error {
	return nil
}

func postReadExtractPatchDeploymentFields(r *PatchDeployment) error {
	vInstanceFilter := r.InstanceFilter
	if vInstanceFilter == nil {
		// note: explicitly not the empty object.
		vInstanceFilter = &PatchDeploymentInstanceFilter{}
	}
	if err := postReadExtractPatchDeploymentInstanceFilterFields(r, vInstanceFilter); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vInstanceFilter) {
		r.InstanceFilter = vInstanceFilter
	}
	vPatchConfig := r.PatchConfig
	if vPatchConfig == nil {
		// note: explicitly not the empty object.
		vPatchConfig = &PatchDeploymentPatchConfig{}
	}
	if err := postReadExtractPatchDeploymentPatchConfigFields(r, vPatchConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPatchConfig) {
		r.PatchConfig = vPatchConfig
	}
	vOneTimeSchedule := r.OneTimeSchedule
	if vOneTimeSchedule == nil {
		// note: explicitly not the empty object.
		vOneTimeSchedule = &PatchDeploymentOneTimeSchedule{}
	}
	if err := postReadExtractPatchDeploymentOneTimeScheduleFields(r, vOneTimeSchedule); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vOneTimeSchedule) {
		r.OneTimeSchedule = vOneTimeSchedule
	}
	vRecurringSchedule := r.RecurringSchedule
	if vRecurringSchedule == nil {
		// note: explicitly not the empty object.
		vRecurringSchedule = &PatchDeploymentRecurringSchedule{}
	}
	if err := postReadExtractPatchDeploymentRecurringScheduleFields(r, vRecurringSchedule); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRecurringSchedule) {
		r.RecurringSchedule = vRecurringSchedule
	}
	vRollout := r.Rollout
	if vRollout == nil {
		// note: explicitly not the empty object.
		vRollout = &PatchDeploymentRollout{}
	}
	if err := postReadExtractPatchDeploymentRolloutFields(r, vRollout); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRollout) {
		r.Rollout = vRollout
	}
	return nil
}
func postReadExtractPatchDeploymentInstanceFilterFields(r *PatchDeployment, o *PatchDeploymentInstanceFilter) error {
	return nil
}
func postReadExtractPatchDeploymentInstanceFilterGroupLabelsFields(r *PatchDeployment, o *PatchDeploymentInstanceFilterGroupLabels) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfig) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &PatchDeploymentPatchConfigApt{}
	}
	if err := extractPatchDeploymentPatchConfigAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &PatchDeploymentPatchConfigYum{}
	}
	if err := extractPatchDeploymentPatchConfigYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vGoo := o.Goo
	if vGoo == nil {
		// note: explicitly not the empty object.
		vGoo = &PatchDeploymentPatchConfigGoo{}
	}
	if err := extractPatchDeploymentPatchConfigGooFields(r, vGoo); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGoo) {
		o.Goo = vGoo
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &PatchDeploymentPatchConfigZypper{}
	}
	if err := extractPatchDeploymentPatchConfigZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vWindowsUpdate := o.WindowsUpdate
	if vWindowsUpdate == nil {
		// note: explicitly not the empty object.
		vWindowsUpdate = &PatchDeploymentPatchConfigWindowsUpdate{}
	}
	if err := extractPatchDeploymentPatchConfigWindowsUpdateFields(r, vWindowsUpdate); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsUpdate) {
		o.WindowsUpdate = vWindowsUpdate
	}
	vPreStep := o.PreStep
	if vPreStep == nil {
		// note: explicitly not the empty object.
		vPreStep = &PatchDeploymentPatchConfigPreStep{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepFields(r, vPreStep); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPreStep) {
		o.PreStep = vPreStep
	}
	vPostStep := o.PostStep
	if vPostStep == nil {
		// note: explicitly not the empty object.
		vPostStep = &PatchDeploymentPatchConfigPostStep{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepFields(r, vPostStep); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPostStep) {
		o.PostStep = vPostStep
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigAptFields(r *PatchDeployment, o *PatchDeploymentPatchConfigApt) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigYumFields(r *PatchDeployment, o *PatchDeploymentPatchConfigYum) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigGooFields(r *PatchDeployment, o *PatchDeploymentPatchConfigGoo) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigZypperFields(r *PatchDeployment, o *PatchDeploymentPatchConfigZypper) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigWindowsUpdateFields(r *PatchDeployment, o *PatchDeploymentPatchConfigWindowsUpdate) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPreStepFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStep) error {
	vLinuxExecStepConfig := o.LinuxExecStepConfig
	if vLinuxExecStepConfig == nil {
		// note: explicitly not the empty object.
		vLinuxExecStepConfig = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigFields(r, vLinuxExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLinuxExecStepConfig) {
		o.LinuxExecStepConfig = vLinuxExecStepConfig
	}
	vWindowsExecStepConfig := o.WindowsExecStepConfig
	if vWindowsExecStepConfig == nil {
		// note: explicitly not the empty object.
		vWindowsExecStepConfig = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigFields(r, vWindowsExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsExecStepConfig) {
		o.WindowsExecStepConfig = vWindowsExecStepConfig
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepLinuxExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepLinuxExecStepConfigGcsObject) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepWindowsExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPreStepWindowsExecStepConfigGcsObject) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPostStepFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStep) error {
	vLinuxExecStepConfig := o.LinuxExecStepConfig
	if vLinuxExecStepConfig == nil {
		// note: explicitly not the empty object.
		vLinuxExecStepConfig = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigFields(r, vLinuxExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLinuxExecStepConfig) {
		o.LinuxExecStepConfig = vLinuxExecStepConfig
	}
	vWindowsExecStepConfig := o.WindowsExecStepConfig
	if vWindowsExecStepConfig == nil {
		// note: explicitly not the empty object.
		vWindowsExecStepConfig = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfig{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigFields(r, vWindowsExecStepConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWindowsExecStepConfig) {
		o.WindowsExecStepConfig = vWindowsExecStepConfig
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepLinuxExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepLinuxExecStepConfigGcsObject) error {
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepWindowsExecStepConfig) error {
	vGcsObject := o.GcsObject
	if vGcsObject == nil {
		// note: explicitly not the empty object.
		vGcsObject = &PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject{}
	}
	if err := extractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectFields(r, vGcsObject); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcsObject) {
		o.GcsObject = vGcsObject
	}
	return nil
}
func postReadExtractPatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObjectFields(r *PatchDeployment, o *PatchDeploymentPatchConfigPostStepWindowsExecStepConfigGcsObject) error {
	return nil
}
func postReadExtractPatchDeploymentOneTimeScheduleFields(r *PatchDeployment, o *PatchDeploymentOneTimeSchedule) error {
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleFields(r *PatchDeployment, o *PatchDeploymentRecurringSchedule) error {
	vTimeZone := o.TimeZone
	if vTimeZone == nil {
		// note: explicitly not the empty object.
		vTimeZone = &PatchDeploymentRecurringScheduleTimeZone{}
	}
	if err := extractPatchDeploymentRecurringScheduleTimeZoneFields(r, vTimeZone); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vTimeZone) {
		o.TimeZone = vTimeZone
	}
	vTimeOfDay := o.TimeOfDay
	if vTimeOfDay == nil {
		// note: explicitly not the empty object.
		vTimeOfDay = &PatchDeploymentRecurringScheduleTimeOfDay{}
	}
	if err := extractPatchDeploymentRecurringScheduleTimeOfDayFields(r, vTimeOfDay); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vTimeOfDay) {
		o.TimeOfDay = vTimeOfDay
	}
	vWeekly := o.Weekly
	if vWeekly == nil {
		// note: explicitly not the empty object.
		vWeekly = &PatchDeploymentRecurringScheduleWeekly{}
	}
	if err := extractPatchDeploymentRecurringScheduleWeeklyFields(r, vWeekly); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWeekly) {
		o.Weekly = vWeekly
	}
	vMonthly := o.Monthly
	if vMonthly == nil {
		// note: explicitly not the empty object.
		vMonthly = &PatchDeploymentRecurringScheduleMonthly{}
	}
	if err := extractPatchDeploymentRecurringScheduleMonthlyFields(r, vMonthly); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMonthly) {
		o.Monthly = vMonthly
	}
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleTimeZoneFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleTimeZone) error {
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleTimeOfDayFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleTimeOfDay) error {
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleWeeklyFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleWeekly) error {
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleMonthlyFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleMonthly) error {
	vWeekDayOfMonth := o.WeekDayOfMonth
	if vWeekDayOfMonth == nil {
		// note: explicitly not the empty object.
		vWeekDayOfMonth = &PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth{}
	}
	if err := extractPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthFields(r, vWeekDayOfMonth); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vWeekDayOfMonth) {
		o.WeekDayOfMonth = vWeekDayOfMonth
	}
	return nil
}
func postReadExtractPatchDeploymentRecurringScheduleMonthlyWeekDayOfMonthFields(r *PatchDeployment, o *PatchDeploymentRecurringScheduleMonthlyWeekDayOfMonth) error {
	return nil
}
func postReadExtractPatchDeploymentRolloutFields(r *PatchDeployment, o *PatchDeploymentRollout) error {
	vDisruptionBudget := o.DisruptionBudget
	if vDisruptionBudget == nil {
		// note: explicitly not the empty object.
		vDisruptionBudget = &PatchDeploymentRolloutDisruptionBudget{}
	}
	if err := extractPatchDeploymentRolloutDisruptionBudgetFields(r, vDisruptionBudget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDisruptionBudget) {
		o.DisruptionBudget = vDisruptionBudget
	}
	return nil
}
func postReadExtractPatchDeploymentRolloutDisruptionBudgetFields(r *PatchDeployment, o *PatchDeploymentRolloutDisruptionBudget) error {
	return nil
}
