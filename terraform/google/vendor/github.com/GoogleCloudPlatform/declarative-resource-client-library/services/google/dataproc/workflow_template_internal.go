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
package dataproc

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

func (r *WorkflowTemplate) validate() error {

	if err := dcl.RequiredParameter(r.Name, "Name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "placement"); err != nil {
		return err
	}
	if err := dcl.Required(r, "jobs"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Placement) {
		if err := r.Placement.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacement) validate() error {
	if !dcl.IsEmptyValueIndirect(r.ManagedCluster) {
		if err := r.ManagedCluster.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ClusterSelector) {
		if err := r.ClusterSelector.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedCluster) validate() error {
	if err := dcl.Required(r, "clusterName"); err != nil {
		return err
	}
	if err := dcl.Required(r, "config"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Config) {
		if err := r.Config.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementClusterSelector) validate() error {
	if err := dcl.Required(r, "clusterLabels"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateJobs) validate() error {
	if err := dcl.Required(r, "stepId"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.HadoopJob) {
		if err := r.HadoopJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SparkJob) {
		if err := r.SparkJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PysparkJob) {
		if err := r.PysparkJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.HiveJob) {
		if err := r.HiveJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PigJob) {
		if err := r.PigJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SparkRJob) {
		if err := r.SparkRJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SparkSqlJob) {
		if err := r.SparkSqlJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PrestoJob) {
		if err := r.PrestoJob.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Scheduling) {
		if err := r.Scheduling.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsHadoopJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsHadoopJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsSparkJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsSparkJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsPysparkJob) validate() error {
	if err := dcl.Required(r, "mainPythonFileUri"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsPysparkJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsHiveJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.QueryList) {
		if err := r.QueryList.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsHiveJobQueryList) validate() error {
	if err := dcl.Required(r, "queries"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateJobsPigJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.QueryList) {
		if err := r.QueryList.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsPigJobQueryList) validate() error {
	if err := dcl.Required(r, "queries"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateJobsPigJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsSparkRJob) validate() error {
	if err := dcl.Required(r, "mainRFileUri"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsSparkRJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsSparkSqlJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.QueryList) {
		if err := r.QueryList.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsSparkSqlJobQueryList) validate() error {
	if err := dcl.Required(r, "queries"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateJobsSparkSqlJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsPrestoJob) validate() error {
	if !dcl.IsEmptyValueIndirect(r.QueryList) {
		if err := r.QueryList.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LoggingConfig) {
		if err := r.LoggingConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateJobsPrestoJobQueryList) validate() error {
	if err := dcl.Required(r, "queries"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateJobsPrestoJobLoggingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplateJobsScheduling) validate() error {
	return nil
}
func (r *WorkflowTemplateParameters) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "fields"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Validation) {
		if err := r.Validation.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateParametersValidation) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Regex) {
		if err := r.Regex.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Values) {
		if err := r.Values.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplateParametersValidationRegex) validate() error {
	if err := dcl.Required(r, "regexes"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplateParametersValidationValues) validate() error {
	if err := dcl.Required(r, "values"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplate) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://dataproc.googleapis.com/v1/", params)
}

func (r *WorkflowTemplate) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *WorkflowTemplate) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workflowTemplates", nr.basePath(), userBasePath, params), nil

}

func (r *WorkflowTemplate) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workflowTemplates", nr.basePath(), userBasePath, params), nil

}

func (r *WorkflowTemplate) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}", nr.basePath(), userBasePath, params), nil
}

// workflowTemplateApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type workflowTemplateApiOperation interface {
	do(context.Context, *WorkflowTemplate, *Client) error
}

func (c *Client) listWorkflowTemplateRaw(ctx context.Context, r *WorkflowTemplate, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != WorkflowTemplateMaxPage {
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

type listWorkflowTemplateOperation struct {
	Templates []map[string]interface{} `json:"templates"`
	Token     string                   `json:"nextPageToken"`
}

func (c *Client) listWorkflowTemplate(ctx context.Context, r *WorkflowTemplate, pageToken string, pageSize int32) ([]*WorkflowTemplate, string, error) {
	b, err := c.listWorkflowTemplateRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listWorkflowTemplateOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*WorkflowTemplate
	for _, v := range m.Templates {
		res, err := unmarshalMapWorkflowTemplate(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllWorkflowTemplate(ctx context.Context, f func(*WorkflowTemplate) bool, resources []*WorkflowTemplate) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteWorkflowTemplate(ctx, res)
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

type deleteWorkflowTemplateOperation struct{}

func (op *deleteWorkflowTemplateOperation) do(ctx context.Context, r *WorkflowTemplate, c *Client) error {
	r, err := c.GetWorkflowTemplate(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "WorkflowTemplate not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetWorkflowTemplate checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete WorkflowTemplate: %w", err)
	}

	// we saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// this is the reason we are adding retry to handle that case.
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		_, err = c.GetWorkflowTemplate(ctx, r)
		if !dcl.IsNotFound(err) {
			if i == maxRetry {
				return dcl.NotDeletedError{ExistingResource: r}
			}
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createWorkflowTemplateOperation struct {
	response map[string]interface{}
}

func (op *createWorkflowTemplateOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createWorkflowTemplateOperation) do(ctx context.Context, r *WorkflowTemplate, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(req, &m); err != nil {
		return err
	}
	normalized := r.urlNormalized()
	m["id"] = fmt.Sprintf("%s", *normalized.Name)

	req, err = json.Marshal(m)
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

	if _, err := c.GetWorkflowTemplate(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getWorkflowTemplateRaw(ctx context.Context, r *WorkflowTemplate) ([]byte, error) {

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

func (c *Client) workflowTemplateDiffsForRawDesired(ctx context.Context, rawDesired *WorkflowTemplate, opts ...dcl.ApplyOption) (initial, desired *WorkflowTemplate, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *WorkflowTemplate
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*WorkflowTemplate); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected WorkflowTemplate, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetWorkflowTemplate(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a WorkflowTemplate resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve WorkflowTemplate resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that WorkflowTemplate resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeWorkflowTemplateDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for WorkflowTemplate: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for WorkflowTemplate: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeWorkflowTemplateInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for WorkflowTemplate: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeWorkflowTemplateDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for WorkflowTemplate: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffWorkflowTemplate(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeWorkflowTemplateInitialState(rawInitial, rawDesired *WorkflowTemplate) (*WorkflowTemplate, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeWorkflowTemplateDesiredState(rawDesired, rawInitial *WorkflowTemplate, opts ...dcl.ApplyOption) (*WorkflowTemplate, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Placement = canonicalizeWorkflowTemplatePlacement(rawDesired.Placement, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &WorkflowTemplate{}
	if dcl.NameToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Labels) {
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	canonicalDesired.Placement = canonicalizeWorkflowTemplatePlacement(rawDesired.Placement, rawInitial.Placement, opts...)
	canonicalDesired.Jobs = canonicalizeWorkflowTemplateJobsSlice(rawDesired.Jobs, rawInitial.Jobs, opts...)
	canonicalDesired.Parameters = canonicalizeWorkflowTemplateParametersSlice(rawDesired.Parameters, rawInitial.Parameters, opts...)
	if dcl.StringCanonicalize(rawDesired.DagTimeout, rawInitial.DagTimeout) {
		canonicalDesired.DagTimeout = rawInitial.DagTimeout
	} else {
		canonicalDesired.DagTimeout = rawDesired.DagTimeout
	}
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

func canonicalizeWorkflowTemplateNewState(c *Client, rawNew, rawDesired *WorkflowTemplate) (*WorkflowTemplate, error) {

	rawNew.Name = rawDesired.Name

	if dcl.IsNotReturnedByServer(rawNew.Version) && dcl.IsNotReturnedByServer(rawDesired.Version) {
		rawNew.Version = rawDesired.Version
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.CreateTime) && dcl.IsNotReturnedByServer(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.UpdateTime) && dcl.IsNotReturnedByServer(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Labels) && dcl.IsNotReturnedByServer(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Placement) && dcl.IsNotReturnedByServer(rawDesired.Placement) {
		rawNew.Placement = rawDesired.Placement
	} else {
		rawNew.Placement = canonicalizeNewWorkflowTemplatePlacement(c, rawDesired.Placement, rawNew.Placement)
	}

	if dcl.IsNotReturnedByServer(rawNew.Jobs) && dcl.IsNotReturnedByServer(rawDesired.Jobs) {
		rawNew.Jobs = rawDesired.Jobs
	} else {
		rawNew.Jobs = canonicalizeNewWorkflowTemplateJobsSlice(c, rawDesired.Jobs, rawNew.Jobs)
	}

	if dcl.IsNotReturnedByServer(rawNew.Parameters) && dcl.IsNotReturnedByServer(rawDesired.Parameters) {
		rawNew.Parameters = rawDesired.Parameters
	} else {
		rawNew.Parameters = canonicalizeNewWorkflowTemplateParametersSlice(c, rawDesired.Parameters, rawNew.Parameters)
	}

	if dcl.IsNotReturnedByServer(rawNew.DagTimeout) && dcl.IsNotReturnedByServer(rawDesired.DagTimeout) {
		rawNew.DagTimeout = rawDesired.DagTimeout
	} else {
		if dcl.StringCanonicalize(rawDesired.DagTimeout, rawNew.DagTimeout) {
			rawNew.DagTimeout = rawDesired.DagTimeout
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeWorkflowTemplatePlacement(des, initial *WorkflowTemplatePlacement, opts ...dcl.ApplyOption) *WorkflowTemplatePlacement {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacement{}

	cDes.ManagedCluster = canonicalizeWorkflowTemplatePlacementManagedCluster(des.ManagedCluster, initial.ManagedCluster, opts...)
	cDes.ClusterSelector = canonicalizeWorkflowTemplatePlacementClusterSelector(des.ClusterSelector, initial.ClusterSelector, opts...)

	return cDes
}

func canonicalizeWorkflowTemplatePlacementSlice(des, initial []WorkflowTemplatePlacement, opts ...dcl.ApplyOption) []WorkflowTemplatePlacement {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacement, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacement(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacement, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacement(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacement(c *Client, des, nw *WorkflowTemplatePlacement) *WorkflowTemplatePlacement {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacement while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ManagedCluster = canonicalizeNewWorkflowTemplatePlacementManagedCluster(c, des.ManagedCluster, nw.ManagedCluster)
	nw.ClusterSelector = canonicalizeNewWorkflowTemplatePlacementClusterSelector(c, des.ClusterSelector, nw.ClusterSelector)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementSet(c *Client, des, nw []WorkflowTemplatePlacement) []WorkflowTemplatePlacement {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplatePlacement
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplatePlacementSlice(c *Client, des, nw []WorkflowTemplatePlacement) []WorkflowTemplatePlacement {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacement
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacement(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedCluster(des, initial *WorkflowTemplatePlacementManagedCluster, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedCluster {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedCluster{}

	if dcl.StringCanonicalize(des.ClusterName, initial.ClusterName) || dcl.IsZeroValue(des.ClusterName) {
		cDes.ClusterName = initial.ClusterName
	} else {
		cDes.ClusterName = des.ClusterName
	}
	cDes.Config = canonicalizeClusterClusterConfig(des.Config, initial.Config, opts...)
	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterSlice(des, initial []WorkflowTemplatePlacementManagedCluster, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedCluster {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedCluster, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedCluster(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedCluster, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedCluster(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedCluster(c *Client, des, nw *WorkflowTemplatePlacementManagedCluster) *WorkflowTemplatePlacementManagedCluster {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedCluster while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ClusterName, nw.ClusterName) {
		nw.ClusterName = des.ClusterName
	}
	nw.Config = canonicalizeNewClusterClusterConfig(c, des.Config, nw.Config)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterSet(c *Client, des, nw []WorkflowTemplatePlacementManagedCluster) []WorkflowTemplatePlacementManagedCluster {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplatePlacementManagedCluster
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplatePlacementManagedClusterSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedCluster) []WorkflowTemplatePlacementManagedCluster {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedCluster
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedCluster(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementClusterSelector(des, initial *WorkflowTemplatePlacementClusterSelector, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementClusterSelector {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementClusterSelector{}

	if dcl.StringCanonicalize(des.Zone, initial.Zone) || dcl.IsZeroValue(des.Zone) {
		cDes.Zone = initial.Zone
	} else {
		cDes.Zone = des.Zone
	}
	if dcl.IsZeroValue(des.ClusterLabels) {
		cDes.ClusterLabels = initial.ClusterLabels
	} else {
		cDes.ClusterLabels = des.ClusterLabels
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementClusterSelectorSlice(des, initial []WorkflowTemplatePlacementClusterSelector, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementClusterSelector {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementClusterSelector, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementClusterSelector(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementClusterSelector, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementClusterSelector(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementClusterSelector(c *Client, des, nw *WorkflowTemplatePlacementClusterSelector) *WorkflowTemplatePlacementClusterSelector {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementClusterSelector while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Zone, nw.Zone) {
		nw.Zone = des.Zone
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementClusterSelectorSet(c *Client, des, nw []WorkflowTemplatePlacementClusterSelector) []WorkflowTemplatePlacementClusterSelector {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplatePlacementClusterSelector
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementClusterSelectorNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplatePlacementClusterSelectorSlice(c *Client, des, nw []WorkflowTemplatePlacementClusterSelector) []WorkflowTemplatePlacementClusterSelector {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementClusterSelector
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementClusterSelector(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobs(des, initial *WorkflowTemplateJobs, opts ...dcl.ApplyOption) *WorkflowTemplateJobs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobs{}

	if dcl.StringCanonicalize(des.StepId, initial.StepId) || dcl.IsZeroValue(des.StepId) {
		cDes.StepId = initial.StepId
	} else {
		cDes.StepId = des.StepId
	}
	cDes.HadoopJob = canonicalizeWorkflowTemplateJobsHadoopJob(des.HadoopJob, initial.HadoopJob, opts...)
	cDes.SparkJob = canonicalizeWorkflowTemplateJobsSparkJob(des.SparkJob, initial.SparkJob, opts...)
	cDes.PysparkJob = canonicalizeWorkflowTemplateJobsPysparkJob(des.PysparkJob, initial.PysparkJob, opts...)
	cDes.HiveJob = canonicalizeWorkflowTemplateJobsHiveJob(des.HiveJob, initial.HiveJob, opts...)
	cDes.PigJob = canonicalizeWorkflowTemplateJobsPigJob(des.PigJob, initial.PigJob, opts...)
	cDes.SparkRJob = canonicalizeWorkflowTemplateJobsSparkRJob(des.SparkRJob, initial.SparkRJob, opts...)
	cDes.SparkSqlJob = canonicalizeWorkflowTemplateJobsSparkSqlJob(des.SparkSqlJob, initial.SparkSqlJob, opts...)
	cDes.PrestoJob = canonicalizeWorkflowTemplateJobsPrestoJob(des.PrestoJob, initial.PrestoJob, opts...)
	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}
	cDes.Scheduling = canonicalizeWorkflowTemplateJobsScheduling(des.Scheduling, initial.Scheduling, opts...)
	if dcl.StringArrayCanonicalize(des.PrerequisiteStepIds, initial.PrerequisiteStepIds) || dcl.IsZeroValue(des.PrerequisiteStepIds) {
		cDes.PrerequisiteStepIds = initial.PrerequisiteStepIds
	} else {
		cDes.PrerequisiteStepIds = des.PrerequisiteStepIds
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSlice(des, initial []WorkflowTemplateJobs, opts ...dcl.ApplyOption) []WorkflowTemplateJobs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobs(c *Client, des, nw *WorkflowTemplateJobs) *WorkflowTemplateJobs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobs while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.StepId, nw.StepId) {
		nw.StepId = des.StepId
	}
	nw.HadoopJob = canonicalizeNewWorkflowTemplateJobsHadoopJob(c, des.HadoopJob, nw.HadoopJob)
	nw.SparkJob = canonicalizeNewWorkflowTemplateJobsSparkJob(c, des.SparkJob, nw.SparkJob)
	nw.PysparkJob = canonicalizeNewWorkflowTemplateJobsPysparkJob(c, des.PysparkJob, nw.PysparkJob)
	nw.HiveJob = canonicalizeNewWorkflowTemplateJobsHiveJob(c, des.HiveJob, nw.HiveJob)
	nw.PigJob = canonicalizeNewWorkflowTemplateJobsPigJob(c, des.PigJob, nw.PigJob)
	nw.SparkRJob = canonicalizeNewWorkflowTemplateJobsSparkRJob(c, des.SparkRJob, nw.SparkRJob)
	nw.SparkSqlJob = canonicalizeNewWorkflowTemplateJobsSparkSqlJob(c, des.SparkSqlJob, nw.SparkSqlJob)
	nw.PrestoJob = canonicalizeNewWorkflowTemplateJobsPrestoJob(c, des.PrestoJob, nw.PrestoJob)
	nw.Scheduling = canonicalizeNewWorkflowTemplateJobsScheduling(c, des.Scheduling, nw.Scheduling)
	if dcl.StringArrayCanonicalize(des.PrerequisiteStepIds, nw.PrerequisiteStepIds) {
		nw.PrerequisiteStepIds = des.PrerequisiteStepIds
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSet(c *Client, des, nw []WorkflowTemplateJobs) []WorkflowTemplateJobs {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobs
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSlice(c *Client, des, nw []WorkflowTemplateJobs) []WorkflowTemplateJobs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobs(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsHadoopJob(des, initial *WorkflowTemplateJobsHadoopJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsHadoopJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsHadoopJob{}

	if dcl.StringCanonicalize(des.MainJarFileUri, initial.MainJarFileUri) || dcl.IsZeroValue(des.MainJarFileUri) {
		cDes.MainJarFileUri = initial.MainJarFileUri
	} else {
		cDes.MainJarFileUri = des.MainJarFileUri
	}
	if dcl.StringCanonicalize(des.MainClass, initial.MainClass) || dcl.IsZeroValue(des.MainClass) {
		cDes.MainClass = initial.MainClass
	} else {
		cDes.MainClass = des.MainClass
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) || dcl.IsZeroValue(des.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) || dcl.IsZeroValue(des.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) || dcl.IsZeroValue(des.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsHadoopJobSlice(des, initial []WorkflowTemplateJobsHadoopJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHadoopJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsHadoopJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsHadoopJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsHadoopJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsHadoopJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsHadoopJob(c *Client, des, nw *WorkflowTemplateJobsHadoopJob) *WorkflowTemplateJobsHadoopJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsHadoopJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MainJarFileUri, nw.MainJarFileUri) {
		nw.MainJarFileUri = des.MainJarFileUri
	}
	if dcl.StringCanonicalize(des.MainClass, nw.MainClass) {
		nw.MainClass = des.MainClass
	}
	if dcl.StringArrayCanonicalize(des.Args, nw.Args) {
		nw.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, nw.FileUris) {
		nw.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, nw.ArchiveUris) {
		nw.ArchiveUris = des.ArchiveUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsHadoopJobSet(c *Client, des, nw []WorkflowTemplateJobsHadoopJob) []WorkflowTemplateJobsHadoopJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsHadoopJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHadoopJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsHadoopJobSlice(c *Client, des, nw []WorkflowTemplateJobsHadoopJob) []WorkflowTemplateJobsHadoopJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsHadoopJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsHadoopJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfig(des, initial *WorkflowTemplateJobsHadoopJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsHadoopJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsHadoopJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsHadoopJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsHadoopJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsHadoopJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsHadoopJobLoggingConfig) *WorkflowTemplateJobsHadoopJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsHadoopJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsHadoopJobLoggingConfig) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsHadoopJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHadoopJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsHadoopJobLoggingConfig) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsHadoopJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkJob(des, initial *WorkflowTemplateJobsSparkJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkJob{}

	if dcl.StringCanonicalize(des.MainJarFileUri, initial.MainJarFileUri) || dcl.IsZeroValue(des.MainJarFileUri) {
		cDes.MainJarFileUri = initial.MainJarFileUri
	} else {
		cDes.MainJarFileUri = des.MainJarFileUri
	}
	if dcl.StringCanonicalize(des.MainClass, initial.MainClass) || dcl.IsZeroValue(des.MainClass) {
		cDes.MainClass = initial.MainClass
	} else {
		cDes.MainClass = des.MainClass
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) || dcl.IsZeroValue(des.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) || dcl.IsZeroValue(des.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) || dcl.IsZeroValue(des.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkJobSlice(des, initial []WorkflowTemplateJobsSparkJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkJob(c *Client, des, nw *WorkflowTemplateJobsSparkJob) *WorkflowTemplateJobsSparkJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MainJarFileUri, nw.MainJarFileUri) {
		nw.MainJarFileUri = des.MainJarFileUri
	}
	if dcl.StringCanonicalize(des.MainClass, nw.MainClass) {
		nw.MainClass = des.MainClass
	}
	if dcl.StringArrayCanonicalize(des.Args, nw.Args) {
		nw.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, nw.FileUris) {
		nw.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, nw.ArchiveUris) {
		nw.ArchiveUris = des.ArchiveUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkJobSet(c *Client, des, nw []WorkflowTemplateJobsSparkJob) []WorkflowTemplateJobsSparkJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkJobSlice(c *Client, des, nw []WorkflowTemplateJobsSparkJob) []WorkflowTemplateJobsSparkJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkJobLoggingConfig(des, initial *WorkflowTemplateJobsSparkJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsSparkJobLoggingConfig) *WorkflowTemplateJobsSparkJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsSparkJobLoggingConfig) []WorkflowTemplateJobsSparkJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsSparkJobLoggingConfig) []WorkflowTemplateJobsSparkJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPysparkJob(des, initial *WorkflowTemplateJobsPysparkJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPysparkJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPysparkJob{}

	if dcl.StringCanonicalize(des.MainPythonFileUri, initial.MainPythonFileUri) || dcl.IsZeroValue(des.MainPythonFileUri) {
		cDes.MainPythonFileUri = initial.MainPythonFileUri
	} else {
		cDes.MainPythonFileUri = des.MainPythonFileUri
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) || dcl.IsZeroValue(des.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.PythonFileUris, initial.PythonFileUris) || dcl.IsZeroValue(des.PythonFileUris) {
		cDes.PythonFileUris = initial.PythonFileUris
	} else {
		cDes.PythonFileUris = des.PythonFileUris
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) || dcl.IsZeroValue(des.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) || dcl.IsZeroValue(des.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPysparkJobSlice(des, initial []WorkflowTemplateJobsPysparkJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPysparkJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPysparkJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPysparkJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPysparkJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPysparkJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPysparkJob(c *Client, des, nw *WorkflowTemplateJobsPysparkJob) *WorkflowTemplateJobsPysparkJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPysparkJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MainPythonFileUri, nw.MainPythonFileUri) {
		nw.MainPythonFileUri = des.MainPythonFileUri
	}
	if dcl.StringArrayCanonicalize(des.Args, nw.Args) {
		nw.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.PythonFileUris, nw.PythonFileUris) {
		nw.PythonFileUris = des.PythonFileUris
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, nw.FileUris) {
		nw.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, nw.ArchiveUris) {
		nw.ArchiveUris = des.ArchiveUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPysparkJobSet(c *Client, des, nw []WorkflowTemplateJobsPysparkJob) []WorkflowTemplateJobsPysparkJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPysparkJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPysparkJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPysparkJobSlice(c *Client, des, nw []WorkflowTemplateJobsPysparkJob) []WorkflowTemplateJobsPysparkJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPysparkJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPysparkJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfig(des, initial *WorkflowTemplateJobsPysparkJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPysparkJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPysparkJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPysparkJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPysparkJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPysparkJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsPysparkJobLoggingConfig) *WorkflowTemplateJobsPysparkJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPysparkJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsPysparkJobLoggingConfig) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPysparkJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPysparkJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsPysparkJobLoggingConfig) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPysparkJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsHiveJob(des, initial *WorkflowTemplateJobsHiveJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsHiveJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsHiveJob{}

	if dcl.StringCanonicalize(des.QueryFileUri, initial.QueryFileUri) || dcl.IsZeroValue(des.QueryFileUri) {
		cDes.QueryFileUri = initial.QueryFileUri
	} else {
		cDes.QueryFileUri = des.QueryFileUri
	}
	cDes.QueryList = canonicalizeWorkflowTemplateJobsHiveJobQueryList(des.QueryList, initial.QueryList, opts...)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, initial.ContinueOnFailure) || dcl.IsZeroValue(des.ContinueOnFailure) {
		cDes.ContinueOnFailure = initial.ContinueOnFailure
	} else {
		cDes.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.IsZeroValue(des.ScriptVariables) {
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHiveJobSlice(des, initial []WorkflowTemplateJobsHiveJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHiveJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsHiveJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsHiveJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsHiveJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsHiveJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsHiveJob(c *Client, des, nw *WorkflowTemplateJobsHiveJob) *WorkflowTemplateJobsHiveJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsHiveJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.QueryFileUri, nw.QueryFileUri) {
		nw.QueryFileUri = des.QueryFileUri
	}
	nw.QueryList = canonicalizeNewWorkflowTemplateJobsHiveJobQueryList(c, des.QueryList, nw.QueryList)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, nw.ContinueOnFailure) {
		nw.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsHiveJobSet(c *Client, des, nw []WorkflowTemplateJobsHiveJob) []WorkflowTemplateJobsHiveJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsHiveJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHiveJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsHiveJobSlice(c *Client, des, nw []WorkflowTemplateJobsHiveJob) []WorkflowTemplateJobsHiveJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsHiveJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsHiveJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsHiveJobQueryList(des, initial *WorkflowTemplateJobsHiveJobQueryList, opts ...dcl.ApplyOption) *WorkflowTemplateJobsHiveJobQueryList {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsHiveJobQueryList{}

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) || dcl.IsZeroValue(des.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHiveJobQueryListSlice(des, initial []WorkflowTemplateJobsHiveJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHiveJobQueryList {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsHiveJobQueryList, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsHiveJobQueryList(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsHiveJobQueryList, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsHiveJobQueryList(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsHiveJobQueryList(c *Client, des, nw *WorkflowTemplateJobsHiveJobQueryList) *WorkflowTemplateJobsHiveJobQueryList {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsHiveJobQueryList while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Queries, nw.Queries) {
		nw.Queries = des.Queries
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsHiveJobQueryListSet(c *Client, des, nw []WorkflowTemplateJobsHiveJobQueryList) []WorkflowTemplateJobsHiveJobQueryList {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsHiveJobQueryList
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHiveJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsHiveJobQueryListSlice(c *Client, des, nw []WorkflowTemplateJobsHiveJobQueryList) []WorkflowTemplateJobsHiveJobQueryList {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsHiveJobQueryList
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsHiveJobQueryList(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPigJob(des, initial *WorkflowTemplateJobsPigJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPigJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPigJob{}

	if dcl.StringCanonicalize(des.QueryFileUri, initial.QueryFileUri) || dcl.IsZeroValue(des.QueryFileUri) {
		cDes.QueryFileUri = initial.QueryFileUri
	} else {
		cDes.QueryFileUri = des.QueryFileUri
	}
	cDes.QueryList = canonicalizeWorkflowTemplateJobsPigJobQueryList(des.QueryList, initial.QueryList, opts...)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, initial.ContinueOnFailure) || dcl.IsZeroValue(des.ContinueOnFailure) {
		cDes.ContinueOnFailure = initial.ContinueOnFailure
	} else {
		cDes.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.IsZeroValue(des.ScriptVariables) {
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPigJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobSlice(des, initial []WorkflowTemplateJobsPigJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPigJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPigJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPigJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPigJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPigJob(c *Client, des, nw *WorkflowTemplateJobsPigJob) *WorkflowTemplateJobsPigJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPigJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.QueryFileUri, nw.QueryFileUri) {
		nw.QueryFileUri = des.QueryFileUri
	}
	nw.QueryList = canonicalizeNewWorkflowTemplateJobsPigJobQueryList(c, des.QueryList, nw.QueryList)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, nw.ContinueOnFailure) {
		nw.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPigJobSet(c *Client, des, nw []WorkflowTemplateJobsPigJob) []WorkflowTemplateJobsPigJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPigJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPigJobSlice(c *Client, des, nw []WorkflowTemplateJobsPigJob) []WorkflowTemplateJobsPigJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPigJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPigJobQueryList(des, initial *WorkflowTemplateJobsPigJobQueryList, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPigJobQueryList {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPigJobQueryList{}

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) || dcl.IsZeroValue(des.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobQueryListSlice(des, initial []WorkflowTemplateJobsPigJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJobQueryList {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPigJobQueryList, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPigJobQueryList(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPigJobQueryList, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPigJobQueryList(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPigJobQueryList(c *Client, des, nw *WorkflowTemplateJobsPigJobQueryList) *WorkflowTemplateJobsPigJobQueryList {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPigJobQueryList while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Queries, nw.Queries) {
		nw.Queries = des.Queries
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPigJobQueryListSet(c *Client, des, nw []WorkflowTemplateJobsPigJobQueryList) []WorkflowTemplateJobsPigJobQueryList {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPigJobQueryList
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPigJobQueryListSlice(c *Client, des, nw []WorkflowTemplateJobsPigJobQueryList) []WorkflowTemplateJobsPigJobQueryList {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPigJobQueryList
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJobQueryList(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPigJobLoggingConfig(des, initial *WorkflowTemplateJobsPigJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPigJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPigJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPigJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPigJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPigJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPigJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPigJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsPigJobLoggingConfig) *WorkflowTemplateJobsPigJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPigJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsPigJobLoggingConfig) []WorkflowTemplateJobsPigJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPigJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsPigJobLoggingConfig) []WorkflowTemplateJobsPigJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPigJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkRJob(des, initial *WorkflowTemplateJobsSparkRJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkRJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkRJob{}

	if dcl.StringCanonicalize(des.MainRFileUri, initial.MainRFileUri) || dcl.IsZeroValue(des.MainRFileUri) {
		cDes.MainRFileUri = initial.MainRFileUri
	} else {
		cDes.MainRFileUri = des.MainRFileUri
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) || dcl.IsZeroValue(des.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) || dcl.IsZeroValue(des.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) || dcl.IsZeroValue(des.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkRJobSlice(des, initial []WorkflowTemplateJobsSparkRJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkRJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkRJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkRJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkRJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkRJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkRJob(c *Client, des, nw *WorkflowTemplateJobsSparkRJob) *WorkflowTemplateJobsSparkRJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkRJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MainRFileUri, nw.MainRFileUri) {
		nw.MainRFileUri = des.MainRFileUri
	}
	if dcl.StringArrayCanonicalize(des.Args, nw.Args) {
		nw.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.FileUris, nw.FileUris) {
		nw.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, nw.ArchiveUris) {
		nw.ArchiveUris = des.ArchiveUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkRJobSet(c *Client, des, nw []WorkflowTemplateJobsSparkRJob) []WorkflowTemplateJobsSparkRJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkRJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkRJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkRJobSlice(c *Client, des, nw []WorkflowTemplateJobsSparkRJob) []WorkflowTemplateJobsSparkRJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkRJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkRJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfig(des, initial *WorkflowTemplateJobsSparkRJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkRJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkRJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkRJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkRJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkRJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsSparkRJobLoggingConfig) *WorkflowTemplateJobsSparkRJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkRJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsSparkRJobLoggingConfig) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkRJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkRJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsSparkRJobLoggingConfig) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkRJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkSqlJob(des, initial *WorkflowTemplateJobsSparkSqlJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkSqlJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkSqlJob{}

	if dcl.StringCanonicalize(des.QueryFileUri, initial.QueryFileUri) || dcl.IsZeroValue(des.QueryFileUri) {
		cDes.QueryFileUri = initial.QueryFileUri
	} else {
		cDes.QueryFileUri = des.QueryFileUri
	}
	cDes.QueryList = canonicalizeWorkflowTemplateJobsSparkSqlJobQueryList(des.QueryList, initial.QueryList, opts...)
	if dcl.IsZeroValue(des.ScriptVariables) {
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) || dcl.IsZeroValue(des.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobSlice(des, initial []WorkflowTemplateJobsSparkSqlJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkSqlJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkSqlJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkSqlJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkSqlJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJob(c *Client, des, nw *WorkflowTemplateJobsSparkSqlJob) *WorkflowTemplateJobsSparkSqlJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkSqlJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.QueryFileUri, nw.QueryFileUri) {
		nw.QueryFileUri = des.QueryFileUri
	}
	nw.QueryList = canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryList(c, des.QueryList, nw.QueryList)
	if dcl.StringArrayCanonicalize(des.JarFileUris, nw.JarFileUris) {
		nw.JarFileUris = des.JarFileUris
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobSet(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJob) []WorkflowTemplateJobsSparkSqlJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkSqlJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobSlice(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJob) []WorkflowTemplateJobsSparkSqlJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkSqlJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobQueryList(des, initial *WorkflowTemplateJobsSparkSqlJobQueryList, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkSqlJobQueryList {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkSqlJobQueryList{}

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) || dcl.IsZeroValue(des.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobQueryListSlice(des, initial []WorkflowTemplateJobsSparkSqlJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJobQueryList {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkSqlJobQueryList, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkSqlJobQueryList(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobQueryList, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkSqlJobQueryList(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryList(c *Client, des, nw *WorkflowTemplateJobsSparkSqlJobQueryList) *WorkflowTemplateJobsSparkSqlJobQueryList {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkSqlJobQueryList while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Queries, nw.Queries) {
		nw.Queries = des.Queries
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryListSet(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJobQueryList) []WorkflowTemplateJobsSparkSqlJobQueryList {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkSqlJobQueryList
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryListSlice(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJobQueryList) []WorkflowTemplateJobsSparkSqlJobQueryList {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkSqlJobQueryList
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryList(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfig(des, initial *WorkflowTemplateJobsSparkSqlJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsSparkSqlJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkSqlJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsSparkSqlJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsSparkSqlJobLoggingConfig) *WorkflowTemplateJobsSparkSqlJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsSparkSqlJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJobLoggingConfig) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsSparkSqlJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsSparkSqlJobLoggingConfig) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsSparkSqlJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPrestoJob(des, initial *WorkflowTemplateJobsPrestoJob, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPrestoJob {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPrestoJob{}

	if dcl.StringCanonicalize(des.QueryFileUri, initial.QueryFileUri) || dcl.IsZeroValue(des.QueryFileUri) {
		cDes.QueryFileUri = initial.QueryFileUri
	} else {
		cDes.QueryFileUri = des.QueryFileUri
	}
	cDes.QueryList = canonicalizeWorkflowTemplateJobsPrestoJobQueryList(des.QueryList, initial.QueryList, opts...)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, initial.ContinueOnFailure) || dcl.IsZeroValue(des.ContinueOnFailure) {
		cDes.ContinueOnFailure = initial.ContinueOnFailure
	} else {
		cDes.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.StringCanonicalize(des.OutputFormat, initial.OutputFormat) || dcl.IsZeroValue(des.OutputFormat) {
		cDes.OutputFormat = initial.OutputFormat
	} else {
		cDes.OutputFormat = des.OutputFormat
	}
	if dcl.StringArrayCanonicalize(des.ClientTags, initial.ClientTags) || dcl.IsZeroValue(des.ClientTags) {
		cDes.ClientTags = initial.ClientTags
	} else {
		cDes.ClientTags = des.ClientTags
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobSlice(des, initial []WorkflowTemplateJobsPrestoJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJob {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPrestoJob, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPrestoJob(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPrestoJob, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPrestoJob(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPrestoJob(c *Client, des, nw *WorkflowTemplateJobsPrestoJob) *WorkflowTemplateJobsPrestoJob {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPrestoJob while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.QueryFileUri, nw.QueryFileUri) {
		nw.QueryFileUri = des.QueryFileUri
	}
	nw.QueryList = canonicalizeNewWorkflowTemplateJobsPrestoJobQueryList(c, des.QueryList, nw.QueryList)
	if dcl.BoolCanonicalize(des.ContinueOnFailure, nw.ContinueOnFailure) {
		nw.ContinueOnFailure = des.ContinueOnFailure
	}
	if dcl.StringCanonicalize(des.OutputFormat, nw.OutputFormat) {
		nw.OutputFormat = des.OutputFormat
	}
	if dcl.StringArrayCanonicalize(des.ClientTags, nw.ClientTags) {
		nw.ClientTags = des.ClientTags
	}
	nw.LoggingConfig = canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfig(c, des.LoggingConfig, nw.LoggingConfig)

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPrestoJobSet(c *Client, des, nw []WorkflowTemplateJobsPrestoJob) []WorkflowTemplateJobsPrestoJob {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPrestoJob
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPrestoJobSlice(c *Client, des, nw []WorkflowTemplateJobsPrestoJob) []WorkflowTemplateJobsPrestoJob {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPrestoJob
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJob(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPrestoJobQueryList(des, initial *WorkflowTemplateJobsPrestoJobQueryList, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPrestoJobQueryList {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPrestoJobQueryList{}

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) || dcl.IsZeroValue(des.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobQueryListSlice(des, initial []WorkflowTemplateJobsPrestoJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJobQueryList {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPrestoJobQueryList, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPrestoJobQueryList(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPrestoJobQueryList, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPrestoJobQueryList(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPrestoJobQueryList(c *Client, des, nw *WorkflowTemplateJobsPrestoJobQueryList) *WorkflowTemplateJobsPrestoJobQueryList {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPrestoJobQueryList while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Queries, nw.Queries) {
		nw.Queries = des.Queries
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPrestoJobQueryListSet(c *Client, des, nw []WorkflowTemplateJobsPrestoJobQueryList) []WorkflowTemplateJobsPrestoJobQueryList {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPrestoJobQueryList
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPrestoJobQueryListSlice(c *Client, des, nw []WorkflowTemplateJobsPrestoJobQueryList) []WorkflowTemplateJobsPrestoJobQueryList {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPrestoJobQueryList
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJobQueryList(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfig(des, initial *WorkflowTemplateJobsPrestoJobLoggingConfig, opts ...dcl.ApplyOption) *WorkflowTemplateJobsPrestoJobLoggingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsPrestoJobLoggingConfig{}

	if dcl.IsZeroValue(des.DriverLogLevels) {
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPrestoJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsPrestoJobLoggingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsPrestoJobLoggingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfig(c *Client, des, nw *WorkflowTemplateJobsPrestoJobLoggingConfig) *WorkflowTemplateJobsPrestoJobLoggingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsPrestoJobLoggingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfigSet(c *Client, des, nw []WorkflowTemplateJobsPrestoJobLoggingConfig) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsPrestoJobLoggingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfigSlice(c *Client, des, nw []WorkflowTemplateJobsPrestoJobLoggingConfig) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsPrestoJobLoggingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateJobsScheduling(des, initial *WorkflowTemplateJobsScheduling, opts ...dcl.ApplyOption) *WorkflowTemplateJobsScheduling {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateJobsScheduling{}

	if dcl.IsZeroValue(des.MaxFailuresPerHour) {
		cDes.MaxFailuresPerHour = initial.MaxFailuresPerHour
	} else {
		cDes.MaxFailuresPerHour = des.MaxFailuresPerHour
	}
	if dcl.IsZeroValue(des.MaxFailuresTotal) {
		cDes.MaxFailuresTotal = initial.MaxFailuresTotal
	} else {
		cDes.MaxFailuresTotal = des.MaxFailuresTotal
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSchedulingSlice(des, initial []WorkflowTemplateJobsScheduling, opts ...dcl.ApplyOption) []WorkflowTemplateJobsScheduling {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateJobsScheduling, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateJobsScheduling(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateJobsScheduling, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateJobsScheduling(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateJobsScheduling(c *Client, des, nw *WorkflowTemplateJobsScheduling) *WorkflowTemplateJobsScheduling {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateJobsScheduling while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplateJobsSchedulingSet(c *Client, des, nw []WorkflowTemplateJobsScheduling) []WorkflowTemplateJobsScheduling {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateJobsScheduling
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSchedulingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateJobsSchedulingSlice(c *Client, des, nw []WorkflowTemplateJobsScheduling) []WorkflowTemplateJobsScheduling {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateJobsScheduling
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateJobsScheduling(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateParameters(des, initial *WorkflowTemplateParameters, opts ...dcl.ApplyOption) *WorkflowTemplateParameters {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateParameters{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.StringArrayCanonicalize(des.Fields, initial.Fields) || dcl.IsZeroValue(des.Fields) {
		cDes.Fields = initial.Fields
	} else {
		cDes.Fields = des.Fields
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	cDes.Validation = canonicalizeWorkflowTemplateParametersValidation(des.Validation, initial.Validation, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateParametersSlice(des, initial []WorkflowTemplateParameters, opts ...dcl.ApplyOption) []WorkflowTemplateParameters {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateParameters, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateParameters(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateParameters, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateParameters(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateParameters(c *Client, des, nw *WorkflowTemplateParameters) *WorkflowTemplateParameters {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateParameters while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringArrayCanonicalize(des.Fields, nw.Fields) {
		nw.Fields = des.Fields
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	nw.Validation = canonicalizeNewWorkflowTemplateParametersValidation(c, des.Validation, nw.Validation)

	return nw
}

func canonicalizeNewWorkflowTemplateParametersSet(c *Client, des, nw []WorkflowTemplateParameters) []WorkflowTemplateParameters {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateParameters
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateParametersSlice(c *Client, des, nw []WorkflowTemplateParameters) []WorkflowTemplateParameters {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateParameters
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateParameters(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateParametersValidation(des, initial *WorkflowTemplateParametersValidation, opts ...dcl.ApplyOption) *WorkflowTemplateParametersValidation {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateParametersValidation{}

	cDes.Regex = canonicalizeWorkflowTemplateParametersValidationRegex(des.Regex, initial.Regex, opts...)
	cDes.Values = canonicalizeWorkflowTemplateParametersValidationValues(des.Values, initial.Values, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateParametersValidationSlice(des, initial []WorkflowTemplateParametersValidation, opts ...dcl.ApplyOption) []WorkflowTemplateParametersValidation {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateParametersValidation, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateParametersValidation(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateParametersValidation, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateParametersValidation(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateParametersValidation(c *Client, des, nw *WorkflowTemplateParametersValidation) *WorkflowTemplateParametersValidation {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateParametersValidation while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Regex = canonicalizeNewWorkflowTemplateParametersValidationRegex(c, des.Regex, nw.Regex)
	nw.Values = canonicalizeNewWorkflowTemplateParametersValidationValues(c, des.Values, nw.Values)

	return nw
}

func canonicalizeNewWorkflowTemplateParametersValidationSet(c *Client, des, nw []WorkflowTemplateParametersValidation) []WorkflowTemplateParametersValidation {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateParametersValidation
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateParametersValidationSlice(c *Client, des, nw []WorkflowTemplateParametersValidation) []WorkflowTemplateParametersValidation {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateParametersValidation
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateParametersValidation(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateParametersValidationRegex(des, initial *WorkflowTemplateParametersValidationRegex, opts ...dcl.ApplyOption) *WorkflowTemplateParametersValidationRegex {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateParametersValidationRegex{}

	if dcl.StringArrayCanonicalize(des.Regexes, initial.Regexes) || dcl.IsZeroValue(des.Regexes) {
		cDes.Regexes = initial.Regexes
	} else {
		cDes.Regexes = des.Regexes
	}

	return cDes
}

func canonicalizeWorkflowTemplateParametersValidationRegexSlice(des, initial []WorkflowTemplateParametersValidationRegex, opts ...dcl.ApplyOption) []WorkflowTemplateParametersValidationRegex {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateParametersValidationRegex, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateParametersValidationRegex(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateParametersValidationRegex, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateParametersValidationRegex(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateParametersValidationRegex(c *Client, des, nw *WorkflowTemplateParametersValidationRegex) *WorkflowTemplateParametersValidationRegex {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateParametersValidationRegex while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Regexes, nw.Regexes) {
		nw.Regexes = des.Regexes
	}

	return nw
}

func canonicalizeNewWorkflowTemplateParametersValidationRegexSet(c *Client, des, nw []WorkflowTemplateParametersValidationRegex) []WorkflowTemplateParametersValidationRegex {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateParametersValidationRegex
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationRegexNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateParametersValidationRegexSlice(c *Client, des, nw []WorkflowTemplateParametersValidationRegex) []WorkflowTemplateParametersValidationRegex {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateParametersValidationRegex
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateParametersValidationRegex(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplateParametersValidationValues(des, initial *WorkflowTemplateParametersValidationValues, opts ...dcl.ApplyOption) *WorkflowTemplateParametersValidationValues {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplateParametersValidationValues{}

	if dcl.StringArrayCanonicalize(des.Values, initial.Values) || dcl.IsZeroValue(des.Values) {
		cDes.Values = initial.Values
	} else {
		cDes.Values = des.Values
	}

	return cDes
}

func canonicalizeWorkflowTemplateParametersValidationValuesSlice(des, initial []WorkflowTemplateParametersValidationValues, opts ...dcl.ApplyOption) []WorkflowTemplateParametersValidationValues {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplateParametersValidationValues, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplateParametersValidationValues(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplateParametersValidationValues, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplateParametersValidationValues(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplateParametersValidationValues(c *Client, des, nw *WorkflowTemplateParametersValidationValues) *WorkflowTemplateParametersValidationValues {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplateParametersValidationValues while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Values, nw.Values) {
		nw.Values = des.Values
	}

	return nw
}

func canonicalizeNewWorkflowTemplateParametersValidationValuesSet(c *Client, des, nw []WorkflowTemplateParametersValidationValues) []WorkflowTemplateParametersValidationValues {
	if des == nil {
		return nw
	}
	var reorderedNew []WorkflowTemplateParametersValidationValues
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewWorkflowTemplateParametersValidationValuesSlice(c *Client, des, nw []WorkflowTemplateParametersValidationValues) []WorkflowTemplateParametersValidationValues {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplateParametersValidationValues
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplateParametersValidationValues(c, &d, &n))
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
func diffWorkflowTemplate(c *Client, desired, actual *WorkflowTemplate, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Placement, actual.Placement, dcl.Info{ObjectFunction: compareWorkflowTemplatePlacementNewStyle, EmptyObject: EmptyWorkflowTemplatePlacement, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Placement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Jobs, actual.Jobs, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsNewStyle, EmptyObject: EmptyWorkflowTemplateJobs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Jobs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Parameters, actual.Parameters, dcl.Info{ObjectFunction: compareWorkflowTemplateParametersNewStyle, EmptyObject: EmptyWorkflowTemplateParameters, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Parameters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DagTimeout, actual.DagTimeout, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DagTimeout")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	return newDiffs, nil
}
func compareWorkflowTemplatePlacementNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacement)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacement)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacement or *WorkflowTemplatePlacement", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacement)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacement)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacement", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ManagedCluster, actual.ManagedCluster, dcl.Info{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedCluster, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedCluster")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterSelector, actual.ClusterSelector, dcl.Info{ObjectFunction: compareWorkflowTemplatePlacementClusterSelectorNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementClusterSelector, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterSelector")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedCluster)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedCluster)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedCluster or *WorkflowTemplatePlacementManagedCluster", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedCluster)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedCluster)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedCluster", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ClusterName, actual.ClusterName, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.Info{ObjectFunction: compareClusterClusterConfigNewStyle, EmptyObject: EmptyClusterClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementClusterSelectorNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementClusterSelector)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementClusterSelector)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementClusterSelector or *WorkflowTemplatePlacementClusterSelector", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementClusterSelector)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementClusterSelector)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementClusterSelector", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Zone")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterLabels, actual.ClusterLabels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobs)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobs or *WorkflowTemplateJobs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobs)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.StepId, actual.StepId, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StepId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HadoopJob, actual.HadoopJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsHadoopJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHadoopJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HadoopJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkJob, actual.SparkJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PysparkJob, actual.PysparkJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPysparkJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPysparkJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PysparkJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HiveJob, actual.HiveJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsHiveJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHiveJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HiveJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PigJob, actual.PigJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPigJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PigJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkRJob, actual.SparkRJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkRJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkRJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkRJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkSqlJob, actual.SparkSqlJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkSqlJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrestoJob, actual.PrestoJob, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPrestoJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrestoJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Scheduling, actual.Scheduling, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSchedulingNewStyle, EmptyObject: EmptyWorkflowTemplateJobsScheduling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Scheduling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrerequisiteStepIds, actual.PrerequisiteStepIds, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrerequisiteStepIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsHadoopJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsHadoopJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsHadoopJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHadoopJob or *WorkflowTemplateJobsHadoopJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsHadoopJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsHadoopJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHadoopJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MainJarFileUri, actual.MainJarFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainJarFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MainClass, actual.MainClass, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainClass")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsHadoopJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHadoopJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsHadoopJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsHadoopJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsHadoopJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHadoopJobLoggingConfig or *WorkflowTemplateJobsHadoopJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsHadoopJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsHadoopJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHadoopJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkJob or *WorkflowTemplateJobsSparkJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MainJarFileUri, actual.MainJarFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainJarFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MainClass, actual.MainClass, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainClass")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkJobLoggingConfig or *WorkflowTemplateJobsSparkJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPysparkJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPysparkJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPysparkJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPysparkJob or *WorkflowTemplateJobsPysparkJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPysparkJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPysparkJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPysparkJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MainPythonFileUri, actual.MainPythonFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainPythonFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PythonFileUris, actual.PythonFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PythonFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPysparkJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPysparkJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPysparkJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPysparkJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPysparkJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPysparkJobLoggingConfig or *WorkflowTemplateJobsPysparkJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPysparkJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPysparkJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPysparkJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsHiveJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsHiveJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsHiveJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHiveJob or *WorkflowTemplateJobsHiveJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsHiveJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsHiveJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHiveJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsHiveJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHiveJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsHiveJobQueryListNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsHiveJobQueryList)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsHiveJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHiveJobQueryList or *WorkflowTemplateJobsHiveJobQueryList", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsHiveJobQueryList)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsHiveJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsHiveJobQueryList", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPigJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPigJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPigJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJob or *WorkflowTemplateJobsPigJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPigJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPigJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPigJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPigJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPigJobQueryListNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPigJobQueryList)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPigJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJobQueryList or *WorkflowTemplateJobsPigJobQueryList", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPigJobQueryList)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPigJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJobQueryList", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPigJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPigJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPigJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJobLoggingConfig or *WorkflowTemplateJobsPigJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPigJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPigJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPigJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkRJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkRJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkRJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkRJob or *WorkflowTemplateJobsSparkRJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkRJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkRJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkRJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MainRFileUri, actual.MainRFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainRFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkRJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkRJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkRJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkRJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkRJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkRJobLoggingConfig or *WorkflowTemplateJobsSparkRJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkRJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkRJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkRJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkSqlJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkSqlJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkSqlJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJob or *WorkflowTemplateJobsSparkSqlJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkSqlJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkSqlJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkSqlJobQueryListNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkSqlJobQueryList)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkSqlJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJobQueryList or *WorkflowTemplateJobsSparkSqlJobQueryList", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkSqlJobQueryList)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkSqlJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJobQueryList", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSparkSqlJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsSparkSqlJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsSparkSqlJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJobLoggingConfig or *WorkflowTemplateJobsSparkSqlJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsSparkSqlJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsSparkSqlJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsSparkSqlJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPrestoJobNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPrestoJob)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPrestoJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJob or *WorkflowTemplateJobsPrestoJob", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPrestoJob)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPrestoJob)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJob", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPrestoJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OutputFormat, actual.OutputFormat, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OutputFormat")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientTags, actual.ClientTags, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClientTags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.Info{ObjectFunction: compareWorkflowTemplateJobsPrestoJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPrestoJobQueryListNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPrestoJobQueryList)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPrestoJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJobQueryList or *WorkflowTemplateJobsPrestoJobQueryList", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPrestoJobQueryList)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPrestoJobQueryList)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJobQueryList", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsPrestoJobLoggingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsPrestoJobLoggingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsPrestoJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJobLoggingConfig or *WorkflowTemplateJobsPrestoJobLoggingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsPrestoJobLoggingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsPrestoJobLoggingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsPrestoJobLoggingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateJobsSchedulingNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateJobsScheduling)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateJobsScheduling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsScheduling or *WorkflowTemplateJobsScheduling", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateJobsScheduling)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateJobsScheduling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateJobsScheduling", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MaxFailuresPerHour, actual.MaxFailuresPerHour, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxFailuresPerHour")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxFailuresTotal, actual.MaxFailuresTotal, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxFailuresTotal")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateParametersNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateParameters)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateParameters)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParameters or *WorkflowTemplateParameters", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateParameters)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateParameters)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParameters", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fields, actual.Fields, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fields")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Validation, actual.Validation, dcl.Info{ObjectFunction: compareWorkflowTemplateParametersValidationNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidation, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Validation")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateParametersValidationNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateParametersValidation)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateParametersValidation)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidation or *WorkflowTemplateParametersValidation", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateParametersValidation)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateParametersValidation)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidation", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Regex, actual.Regex, dcl.Info{ObjectFunction: compareWorkflowTemplateParametersValidationRegexNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidationRegex, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Regex")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.Info{ObjectFunction: compareWorkflowTemplateParametersValidationValuesNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidationValues, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateParametersValidationRegexNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateParametersValidationRegex)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateParametersValidationRegex)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidationRegex or *WorkflowTemplateParametersValidationRegex", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateParametersValidationRegex)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateParametersValidationRegex)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidationRegex", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Regexes, actual.Regexes, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Regexes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplateParametersValidationValuesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplateParametersValidationValues)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplateParametersValidationValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidationValues or *WorkflowTemplateParametersValidationValues", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplateParametersValidationValues)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplateParametersValidationValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplateParametersValidationValues", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
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
func (r *WorkflowTemplate) urlNormalized() *WorkflowTemplate {
	normalized := dcl.Copy(*r).(WorkflowTemplate)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DagTimeout = dcl.SelfLinkToName(r.DagTimeout)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *WorkflowTemplate) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the WorkflowTemplate resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *WorkflowTemplate) marshal(c *Client) ([]byte, error) {
	m, err := expandWorkflowTemplate(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling WorkflowTemplate: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalWorkflowTemplate decodes JSON responses into the WorkflowTemplate resource schema.
func unmarshalWorkflowTemplate(b []byte, c *Client) (*WorkflowTemplate, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapWorkflowTemplate(m, c)
}

func unmarshalMapWorkflowTemplate(m map[string]interface{}, c *Client) (*WorkflowTemplate, error) {

	flattened := flattenWorkflowTemplate(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandWorkflowTemplate expands WorkflowTemplate into a JSON request object.
func expandWorkflowTemplate(c *Client, f *WorkflowTemplate) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if v != nil {
		m["name"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := expandWorkflowTemplatePlacement(c, f.Placement); err != nil {
		return nil, fmt.Errorf("error expanding Placement into placement: %w", err)
	} else if v != nil {
		m["placement"] = v
	}
	if v, err := expandWorkflowTemplateJobsSlice(c, f.Jobs); err != nil {
		return nil, fmt.Errorf("error expanding Jobs into jobs: %w", err)
	} else {
		m["jobs"] = v
	}
	if v, err := expandWorkflowTemplateParametersSlice(c, f.Parameters); err != nil {
		return nil, fmt.Errorf("error expanding Parameters into parameters: %w", err)
	} else {
		m["parameters"] = v
	}
	if v := f.DagTimeout; dcl.ValueShouldBeSent(v) {
		m["dagTimeout"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if v != nil {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if v != nil {
		m["location"] = v
	}

	return m, nil
}

// flattenWorkflowTemplate flattens WorkflowTemplate from a JSON request object into the
// WorkflowTemplate type.
func flattenWorkflowTemplate(c *Client, i interface{}) *WorkflowTemplate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &WorkflowTemplate{}
	res.Name = dcl.FlattenString(m["name"])
	res.Version = dcl.FlattenInteger(m["version"])
	res.CreateTime = dcl.FlattenString(m["createTime"])
	res.UpdateTime = dcl.FlattenString(m["updateTime"])
	res.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	res.Placement = flattenWorkflowTemplatePlacement(c, m["placement"])
	res.Jobs = flattenWorkflowTemplateJobsSlice(c, m["jobs"])
	res.Parameters = flattenWorkflowTemplateParametersSlice(c, m["parameters"])
	res.DagTimeout = dcl.FlattenString(m["dagTimeout"])
	res.Project = dcl.FlattenString(m["project"])
	res.Location = dcl.FlattenString(m["location"])

	return res
}

// expandWorkflowTemplatePlacementMap expands the contents of WorkflowTemplatePlacement into a JSON
// request object.
func expandWorkflowTemplatePlacementMap(c *Client, f map[string]WorkflowTemplatePlacement) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacement(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementSlice expands the contents of WorkflowTemplatePlacement into a JSON
// request object.
func expandWorkflowTemplatePlacementSlice(c *Client, f []WorkflowTemplatePlacement) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacement(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementMap flattens the contents of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacementMap(c *Client, i interface{}) map[string]WorkflowTemplatePlacement {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacement{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacement{}
	}

	items := make(map[string]WorkflowTemplatePlacement)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacement(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementSlice flattens the contents of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacementSlice(c *Client, i interface{}) []WorkflowTemplatePlacement {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacement{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacement{}
	}

	items := make([]WorkflowTemplatePlacement, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacement(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplatePlacement expands an instance of WorkflowTemplatePlacement into a JSON
// request object.
func expandWorkflowTemplatePlacement(c *Client, f *WorkflowTemplatePlacement) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkflowTemplatePlacementManagedCluster(c, f.ManagedCluster); err != nil {
		return nil, fmt.Errorf("error expanding ManagedCluster into managedCluster: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["managedCluster"] = v
	}
	if v, err := expandWorkflowTemplatePlacementClusterSelector(c, f.ClusterSelector); err != nil {
		return nil, fmt.Errorf("error expanding ClusterSelector into clusterSelector: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["clusterSelector"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacement flattens an instance of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacement(c *Client, i interface{}) *WorkflowTemplatePlacement {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacement{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacement
	}
	r.ManagedCluster = flattenWorkflowTemplatePlacementManagedCluster(c, m["managedCluster"])
	r.ClusterSelector = flattenWorkflowTemplatePlacementClusterSelector(c, m["clusterSelector"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterMap expands the contents of WorkflowTemplatePlacementManagedCluster into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterMap(c *Client, f map[string]WorkflowTemplatePlacementManagedCluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedCluster(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterSlice expands the contents of WorkflowTemplatePlacementManagedCluster into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterSlice(c *Client, f []WorkflowTemplatePlacementManagedCluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedCluster(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterMap flattens the contents of WorkflowTemplatePlacementManagedCluster from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterMap(c *Client, i interface{}) map[string]WorkflowTemplatePlacementManagedCluster {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedCluster{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedCluster{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedCluster)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedCluster(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterSlice flattens the contents of WorkflowTemplatePlacementManagedCluster from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterSlice(c *Client, i interface{}) []WorkflowTemplatePlacementManagedCluster {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedCluster{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedCluster{}
	}

	items := make([]WorkflowTemplatePlacementManagedCluster, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedCluster(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedCluster expands an instance of WorkflowTemplatePlacementManagedCluster into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedCluster(c *Client, f *WorkflowTemplatePlacementManagedCluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ClusterName; !dcl.IsEmptyValueIndirect(v) {
		m["clusterName"] = v
	}
	if v, err := expandClusterClusterConfig(c, f.Config); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedCluster flattens an instance of WorkflowTemplatePlacementManagedCluster from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedCluster(c *Client, i interface{}) *WorkflowTemplatePlacementManagedCluster {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedCluster{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedCluster
	}
	r.ClusterName = dcl.FlattenString(m["clusterName"])
	r.Config = flattenClusterClusterConfig(c, m["config"])
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])

	return r
}

// expandWorkflowTemplatePlacementClusterSelectorMap expands the contents of WorkflowTemplatePlacementClusterSelector into a JSON
// request object.
func expandWorkflowTemplatePlacementClusterSelectorMap(c *Client, f map[string]WorkflowTemplatePlacementClusterSelector) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementClusterSelector(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementClusterSelectorSlice expands the contents of WorkflowTemplatePlacementClusterSelector into a JSON
// request object.
func expandWorkflowTemplatePlacementClusterSelectorSlice(c *Client, f []WorkflowTemplatePlacementClusterSelector) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementClusterSelector(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementClusterSelectorMap flattens the contents of WorkflowTemplatePlacementClusterSelector from a JSON
// response object.
func flattenWorkflowTemplatePlacementClusterSelectorMap(c *Client, i interface{}) map[string]WorkflowTemplatePlacementClusterSelector {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementClusterSelector{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementClusterSelector{}
	}

	items := make(map[string]WorkflowTemplatePlacementClusterSelector)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementClusterSelector(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementClusterSelectorSlice flattens the contents of WorkflowTemplatePlacementClusterSelector from a JSON
// response object.
func flattenWorkflowTemplatePlacementClusterSelectorSlice(c *Client, i interface{}) []WorkflowTemplatePlacementClusterSelector {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementClusterSelector{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementClusterSelector{}
	}

	items := make([]WorkflowTemplatePlacementClusterSelector, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementClusterSelector(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplatePlacementClusterSelector expands an instance of WorkflowTemplatePlacementClusterSelector into a JSON
// request object.
func expandWorkflowTemplatePlacementClusterSelector(c *Client, f *WorkflowTemplatePlacementClusterSelector) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Zone; !dcl.IsEmptyValueIndirect(v) {
		m["zone"] = v
	}
	if v := f.ClusterLabels; !dcl.IsEmptyValueIndirect(v) {
		m["clusterLabels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementClusterSelector flattens an instance of WorkflowTemplatePlacementClusterSelector from a JSON
// response object.
func flattenWorkflowTemplatePlacementClusterSelector(c *Client, i interface{}) *WorkflowTemplatePlacementClusterSelector {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementClusterSelector{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementClusterSelector
	}
	r.Zone = dcl.FlattenString(m["zone"])
	r.ClusterLabels = dcl.FlattenKeyValuePairs(m["clusterLabels"])

	return r
}

// expandWorkflowTemplateJobsMap expands the contents of WorkflowTemplateJobs into a JSON
// request object.
func expandWorkflowTemplateJobsMap(c *Client, f map[string]WorkflowTemplateJobs) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobs(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSlice expands the contents of WorkflowTemplateJobs into a JSON
// request object.
func expandWorkflowTemplateJobsSlice(c *Client, f []WorkflowTemplateJobs) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobs(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsMap flattens the contents of WorkflowTemplateJobs from a JSON
// response object.
func flattenWorkflowTemplateJobsMap(c *Client, i interface{}) map[string]WorkflowTemplateJobs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobs{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobs{}
	}

	items := make(map[string]WorkflowTemplateJobs)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobs(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSlice flattens the contents of WorkflowTemplateJobs from a JSON
// response object.
func flattenWorkflowTemplateJobsSlice(c *Client, i interface{}) []WorkflowTemplateJobs {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobs{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobs{}
	}

	items := make([]WorkflowTemplateJobs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobs(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobs expands an instance of WorkflowTemplateJobs into a JSON
// request object.
func expandWorkflowTemplateJobs(c *Client, f *WorkflowTemplateJobs) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.StepId; !dcl.IsEmptyValueIndirect(v) {
		m["stepId"] = v
	}
	if v, err := expandWorkflowTemplateJobsHadoopJob(c, f.HadoopJob); err != nil {
		return nil, fmt.Errorf("error expanding HadoopJob into hadoopJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["hadoopJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkJob(c, f.SparkJob); err != nil {
		return nil, fmt.Errorf("error expanding SparkJob into sparkJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPysparkJob(c, f.PysparkJob); err != nil {
		return nil, fmt.Errorf("error expanding PysparkJob into pysparkJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pysparkJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsHiveJob(c, f.HiveJob); err != nil {
		return nil, fmt.Errorf("error expanding HiveJob into hiveJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["hiveJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPigJob(c, f.PigJob); err != nil {
		return nil, fmt.Errorf("error expanding PigJob into pigJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pigJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkRJob(c, f.SparkRJob); err != nil {
		return nil, fmt.Errorf("error expanding SparkRJob into sparkRJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkRJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkSqlJob(c, f.SparkSqlJob); err != nil {
		return nil, fmt.Errorf("error expanding SparkSqlJob into sparkSqlJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkSqlJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPrestoJob(c, f.PrestoJob); err != nil {
		return nil, fmt.Errorf("error expanding PrestoJob into prestoJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["prestoJob"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}
	if v, err := expandWorkflowTemplateJobsScheduling(c, f.Scheduling); err != nil {
		return nil, fmt.Errorf("error expanding Scheduling into scheduling: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["scheduling"] = v
	}
	if v := f.PrerequisiteStepIds; v != nil {
		m["prerequisiteStepIds"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobs flattens an instance of WorkflowTemplateJobs from a JSON
// response object.
func flattenWorkflowTemplateJobs(c *Client, i interface{}) *WorkflowTemplateJobs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobs
	}
	r.StepId = dcl.FlattenString(m["stepId"])
	r.HadoopJob = flattenWorkflowTemplateJobsHadoopJob(c, m["hadoopJob"])
	r.SparkJob = flattenWorkflowTemplateJobsSparkJob(c, m["sparkJob"])
	r.PysparkJob = flattenWorkflowTemplateJobsPysparkJob(c, m["pysparkJob"])
	r.HiveJob = flattenWorkflowTemplateJobsHiveJob(c, m["hiveJob"])
	r.PigJob = flattenWorkflowTemplateJobsPigJob(c, m["pigJob"])
	r.SparkRJob = flattenWorkflowTemplateJobsSparkRJob(c, m["sparkRJob"])
	r.SparkSqlJob = flattenWorkflowTemplateJobsSparkSqlJob(c, m["sparkSqlJob"])
	r.PrestoJob = flattenWorkflowTemplateJobsPrestoJob(c, m["prestoJob"])
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	r.Scheduling = flattenWorkflowTemplateJobsScheduling(c, m["scheduling"])
	r.PrerequisiteStepIds = dcl.FlattenStringSlice(m["prerequisiteStepIds"])

	return r
}

// expandWorkflowTemplateJobsHadoopJobMap expands the contents of WorkflowTemplateJobsHadoopJob into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobMap(c *Client, f map[string]WorkflowTemplateJobsHadoopJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsHadoopJobSlice expands the contents of WorkflowTemplateJobsHadoopJob into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobSlice(c *Client, f []WorkflowTemplateJobsHadoopJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHadoopJobMap flattens the contents of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsHadoopJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHadoopJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHadoopJob{}
	}

	items := make(map[string]WorkflowTemplateJobsHadoopJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHadoopJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsHadoopJobSlice flattens the contents of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsHadoopJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHadoopJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHadoopJob{}
	}

	items := make([]WorkflowTemplateJobsHadoopJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHadoopJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsHadoopJob expands an instance of WorkflowTemplateJobsHadoopJob into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJob(c *Client, f *WorkflowTemplateJobsHadoopJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MainJarFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["mainJarFileUri"] = v
	}
	if v := f.MainClass; !dcl.IsEmptyValueIndirect(v) {
		m["mainClass"] = v
	}
	if v := f.Args; v != nil {
		m["args"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}
	if v := f.FileUris; v != nil {
		m["fileUris"] = v
	}
	if v := f.ArchiveUris; v != nil {
		m["archiveUris"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsHadoopJob flattens an instance of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJob(c *Client, i interface{}) *WorkflowTemplateJobsHadoopJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsHadoopJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsHadoopJob
	}
	r.MainJarFileUri = dcl.FlattenString(m["mainJarFileUri"])
	r.MainClass = dcl.FlattenString(m["mainClass"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.FileUris = dcl.FlattenStringSlice(m["fileUris"])
	r.ArchiveUris = dcl.FlattenStringSlice(m["archiveUris"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsHadoopJobLoggingConfigMap expands the contents of WorkflowTemplateJobsHadoopJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsHadoopJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsHadoopJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsHadoopJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsHadoopJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHadoopJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsHadoopJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsHadoopJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsHadoopJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsHadoopJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsHadoopJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsHadoopJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsHadoopJobLoggingConfig expands an instance of WorkflowTemplateJobsHadoopJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobLoggingConfig(c *Client, f *WorkflowTemplateJobsHadoopJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsHadoopJobLoggingConfig flattens an instance of WorkflowTemplateJobsHadoopJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsHadoopJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsHadoopJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsHadoopJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsSparkJobMap expands the contents of WorkflowTemplateJobsSparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkJobSlice expands the contents of WorkflowTemplateJobsSparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobSlice(c *Client, f []WorkflowTemplateJobsSparkJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkJobMap flattens the contents of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkJobSlice flattens the contents of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkJob{}
	}

	items := make([]WorkflowTemplateJobsSparkJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkJob expands an instance of WorkflowTemplateJobsSparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJob(c *Client, f *WorkflowTemplateJobsSparkJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MainJarFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["mainJarFileUri"] = v
	}
	if v := f.MainClass; !dcl.IsEmptyValueIndirect(v) {
		m["mainClass"] = v
	}
	if v := f.Args; v != nil {
		m["args"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}
	if v := f.FileUris; v != nil {
		m["fileUris"] = v
	}
	if v := f.ArchiveUris; v != nil {
		m["archiveUris"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkJob flattens an instance of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJob(c *Client, i interface{}) *WorkflowTemplateJobsSparkJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkJob
	}
	r.MainJarFileUri = dcl.FlattenString(m["mainJarFileUri"])
	r.MainClass = dcl.FlattenString(m["mainClass"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.FileUris = dcl.FlattenStringSlice(m["fileUris"])
	r.ArchiveUris = dcl.FlattenStringSlice(m["archiveUris"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsSparkJobLoggingConfigMap expands the contents of WorkflowTemplateJobsSparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsSparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkJobLoggingConfig flattens an instance of WorkflowTemplateJobsSparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsSparkJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsPysparkJobMap expands the contents of WorkflowTemplateJobsPysparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobMap(c *Client, f map[string]WorkflowTemplateJobsPysparkJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPysparkJobSlice expands the contents of WorkflowTemplateJobsPysparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobSlice(c *Client, f []WorkflowTemplateJobsPysparkJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPysparkJobMap flattens the contents of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPysparkJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPysparkJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPysparkJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPysparkJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPysparkJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPysparkJobSlice flattens the contents of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsPysparkJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPysparkJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPysparkJob{}
	}

	items := make([]WorkflowTemplateJobsPysparkJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPysparkJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPysparkJob expands an instance of WorkflowTemplateJobsPysparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJob(c *Client, f *WorkflowTemplateJobsPysparkJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MainPythonFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["mainPythonFileUri"] = v
	}
	if v := f.Args; v != nil {
		m["args"] = v
	}
	if v := f.PythonFileUris; v != nil {
		m["pythonFileUris"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}
	if v := f.FileUris; v != nil {
		m["fileUris"] = v
	}
	if v := f.ArchiveUris; v != nil {
		m["archiveUris"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPysparkJob flattens an instance of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJob(c *Client, i interface{}) *WorkflowTemplateJobsPysparkJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPysparkJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPysparkJob
	}
	r.MainPythonFileUri = dcl.FlattenString(m["mainPythonFileUri"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.PythonFileUris = dcl.FlattenStringSlice(m["pythonFileUris"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.FileUris = dcl.FlattenStringSlice(m["fileUris"])
	r.ArchiveUris = dcl.FlattenStringSlice(m["archiveUris"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsPysparkJobLoggingConfigMap expands the contents of WorkflowTemplateJobsPysparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPysparkJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPysparkJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsPysparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPysparkJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPysparkJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPysparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPysparkJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPysparkJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPysparkJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPysparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPysparkJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPysparkJobLoggingConfig expands an instance of WorkflowTemplateJobsPysparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPysparkJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPysparkJobLoggingConfig flattens an instance of WorkflowTemplateJobsPysparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsPysparkJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPysparkJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPysparkJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsHiveJobMap expands the contents of WorkflowTemplateJobsHiveJob into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobMap(c *Client, f map[string]WorkflowTemplateJobsHiveJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsHiveJobSlice expands the contents of WorkflowTemplateJobsHiveJob into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobSlice(c *Client, f []WorkflowTemplateJobsHiveJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHiveJobMap flattens the contents of WorkflowTemplateJobsHiveJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsHiveJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHiveJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHiveJob{}
	}

	items := make(map[string]WorkflowTemplateJobsHiveJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHiveJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsHiveJobSlice flattens the contents of WorkflowTemplateJobsHiveJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsHiveJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHiveJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHiveJob{}
	}

	items := make([]WorkflowTemplateJobsHiveJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHiveJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsHiveJob expands an instance of WorkflowTemplateJobsHiveJob into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJob(c *Client, f *WorkflowTemplateJobsHiveJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsHiveJobQueryList(c, f.QueryList); err != nil {
		return nil, fmt.Errorf("error expanding QueryList into queryList: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["queryList"] = v
	}
	if v := f.ContinueOnFailure; !dcl.IsEmptyValueIndirect(v) {
		m["continueOnFailure"] = v
	}
	if v := f.ScriptVariables; !dcl.IsEmptyValueIndirect(v) {
		m["scriptVariables"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsHiveJob flattens an instance of WorkflowTemplateJobsHiveJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJob(c *Client, i interface{}) *WorkflowTemplateJobsHiveJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsHiveJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsHiveJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsHiveJobQueryList(c, m["queryList"])
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])

	return r
}

// expandWorkflowTemplateJobsHiveJobQueryListMap expands the contents of WorkflowTemplateJobsHiveJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsHiveJobQueryList) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsHiveJobQueryListSlice expands the contents of WorkflowTemplateJobsHiveJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobQueryListSlice(c *Client, f []WorkflowTemplateJobsHiveJobQueryList) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHiveJobQueryListMap flattens the contents of WorkflowTemplateJobsHiveJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobQueryListMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsHiveJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHiveJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHiveJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsHiveJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHiveJobQueryList(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsHiveJobQueryListSlice flattens the contents of WorkflowTemplateJobsHiveJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobQueryListSlice(c *Client, i interface{}) []WorkflowTemplateJobsHiveJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHiveJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHiveJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsHiveJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHiveJobQueryList(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsHiveJobQueryList expands an instance of WorkflowTemplateJobsHiveJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobQueryList(c *Client, f *WorkflowTemplateJobsHiveJobQueryList) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Queries; v != nil {
		m["queries"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsHiveJobQueryList flattens an instance of WorkflowTemplateJobsHiveJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobQueryList(c *Client, i interface{}) *WorkflowTemplateJobsHiveJobQueryList {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsHiveJobQueryList{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsHiveJobQueryList
	}
	r.Queries = dcl.FlattenStringSlice(m["queries"])

	return r
}

// expandWorkflowTemplateJobsPigJobMap expands the contents of WorkflowTemplateJobsPigJob into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobMap(c *Client, f map[string]WorkflowTemplateJobsPigJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPigJobSlice expands the contents of WorkflowTemplateJobsPigJob into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobSlice(c *Client, f []WorkflowTemplateJobsPigJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobMap flattens the contents of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPigJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobSlice flattens the contents of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsPigJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJob{}
	}

	items := make([]WorkflowTemplateJobsPigJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPigJob expands an instance of WorkflowTemplateJobsPigJob into a JSON
// request object.
func expandWorkflowTemplateJobsPigJob(c *Client, f *WorkflowTemplateJobsPigJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsPigJobQueryList(c, f.QueryList); err != nil {
		return nil, fmt.Errorf("error expanding QueryList into queryList: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["queryList"] = v
	}
	if v := f.ContinueOnFailure; !dcl.IsEmptyValueIndirect(v) {
		m["continueOnFailure"] = v
	}
	if v := f.ScriptVariables; !dcl.IsEmptyValueIndirect(v) {
		m["scriptVariables"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}
	if v, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPigJob flattens an instance of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJob(c *Client, i interface{}) *WorkflowTemplateJobsPigJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPigJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPigJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsPigJobQueryList(c, m["queryList"])
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.LoggingConfig = flattenWorkflowTemplateJobsPigJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsPigJobQueryListMap expands the contents of WorkflowTemplateJobsPigJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsPigJobQueryList) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPigJobQueryListSlice expands the contents of WorkflowTemplateJobsPigJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobQueryListSlice(c *Client, f []WorkflowTemplateJobsPigJobQueryList) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobQueryListMap flattens the contents of WorkflowTemplateJobsPigJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobQueryListMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPigJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJobQueryList(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobQueryListSlice flattens the contents of WorkflowTemplateJobsPigJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobQueryListSlice(c *Client, i interface{}) []WorkflowTemplateJobsPigJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsPigJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJobQueryList(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPigJobQueryList expands an instance of WorkflowTemplateJobsPigJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobQueryList(c *Client, f *WorkflowTemplateJobsPigJobQueryList) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Queries; v != nil {
		m["queries"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPigJobQueryList flattens an instance of WorkflowTemplateJobsPigJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobQueryList(c *Client, i interface{}) *WorkflowTemplateJobsPigJobQueryList {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPigJobQueryList{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPigJobQueryList
	}
	r.Queries = dcl.FlattenStringSlice(m["queries"])

	return r
}

// expandWorkflowTemplateJobsPigJobLoggingConfigMap expands the contents of WorkflowTemplateJobsPigJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPigJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPigJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsPigJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPigJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPigJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPigJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPigJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsPigJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPigJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPigJobLoggingConfig expands an instance of WorkflowTemplateJobsPigJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPigJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPigJobLoggingConfig flattens an instance of WorkflowTemplateJobsPigJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsPigJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPigJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPigJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsSparkRJobMap expands the contents of WorkflowTemplateJobsSparkRJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkRJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkRJobSlice expands the contents of WorkflowTemplateJobsSparkRJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobSlice(c *Client, f []WorkflowTemplateJobsSparkRJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkRJobMap flattens the contents of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkRJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkRJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkRJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkRJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkRJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkRJobSlice flattens the contents of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkRJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkRJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkRJob{}
	}

	items := make([]WorkflowTemplateJobsSparkRJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkRJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkRJob expands an instance of WorkflowTemplateJobsSparkRJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJob(c *Client, f *WorkflowTemplateJobsSparkRJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MainRFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["mainRFileUri"] = v
	}
	if v := f.Args; v != nil {
		m["args"] = v
	}
	if v := f.FileUris; v != nil {
		m["fileUris"] = v
	}
	if v := f.ArchiveUris; v != nil {
		m["archiveUris"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkRJob flattens an instance of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJob(c *Client, i interface{}) *WorkflowTemplateJobsSparkRJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkRJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkRJob
	}
	r.MainRFileUri = dcl.FlattenString(m["mainRFileUri"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.FileUris = dcl.FlattenStringSlice(m["fileUris"])
	r.ArchiveUris = dcl.FlattenStringSlice(m["archiveUris"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsSparkRJobLoggingConfigMap expands the contents of WorkflowTemplateJobsSparkRJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkRJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkRJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsSparkRJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkRJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkRJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkRJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkRJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkRJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkRJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkRJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkRJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkRJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkRJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkRJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkRJobLoggingConfig flattens an instance of WorkflowTemplateJobsSparkRJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsSparkRJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkRJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkRJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsSparkSqlJobMap expands the contents of WorkflowTemplateJobsSparkSqlJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkSqlJobSlice expands the contents of WorkflowTemplateJobsSparkSqlJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobMap flattens the contents of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkSqlJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobSlice flattens the contents of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkSqlJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJob{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJob expands an instance of WorkflowTemplateJobsSparkSqlJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJob(c *Client, f *WorkflowTemplateJobsSparkSqlJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, f.QueryList); err != nil {
		return nil, fmt.Errorf("error expanding QueryList into queryList: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["queryList"] = v
	}
	if v := f.ScriptVariables; !dcl.IsEmptyValueIndirect(v) {
		m["scriptVariables"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v := f.JarFileUris; v != nil {
		m["jarFileUris"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkSqlJob flattens an instance of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJob(c *Client, i interface{}) *WorkflowTemplateJobsSparkSqlJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkSqlJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkSqlJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, m["queryList"])
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsSparkSqlJobQueryListMap expands the contents of WorkflowTemplateJobsSparkSqlJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJobQueryList) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkSqlJobQueryListSlice expands the contents of WorkflowTemplateJobsSparkSqlJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobQueryListSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJobQueryList) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobQueryListMap flattens the contents of WorkflowTemplateJobsSparkSqlJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobQueryListMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkSqlJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobQueryListSlice flattens the contents of WorkflowTemplateJobsSparkSqlJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobQueryListSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkSqlJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJobQueryList expands an instance of WorkflowTemplateJobsSparkSqlJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobQueryList(c *Client, f *WorkflowTemplateJobsSparkSqlJobQueryList) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Queries; v != nil {
		m["queries"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobQueryList flattens an instance of WorkflowTemplateJobsSparkSqlJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobQueryList(c *Client, i interface{}) *WorkflowTemplateJobsSparkSqlJobQueryList {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkSqlJobQueryList{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkSqlJobQueryList
	}
	r.Queries = dcl.FlattenStringSlice(m["queries"])

	return r
}

// expandWorkflowTemplateJobsSparkSqlJobLoggingConfigMap expands the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkSqlJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkSqlJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig flattens an instance of WorkflowTemplateJobsSparkSqlJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkSqlJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsPrestoJobMap expands the contents of WorkflowTemplateJobsPrestoJob into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJob) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJob(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPrestoJobSlice expands the contents of WorkflowTemplateJobsPrestoJob into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobSlice(c *Client, f []WorkflowTemplateJobsPrestoJob) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJob(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobMap flattens the contents of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPrestoJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJob(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobSlice flattens the contents of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobSlice(c *Client, i interface{}) []WorkflowTemplateJobsPrestoJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJob{}
	}

	items := make([]WorkflowTemplateJobsPrestoJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJob(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJob expands an instance of WorkflowTemplateJobsPrestoJob into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJob(c *Client, f *WorkflowTemplateJobsPrestoJob) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, f.QueryList); err != nil {
		return nil, fmt.Errorf("error expanding QueryList into queryList: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["queryList"] = v
	}
	if v := f.ContinueOnFailure; !dcl.IsEmptyValueIndirect(v) {
		m["continueOnFailure"] = v
	}
	if v := f.OutputFormat; !dcl.IsEmptyValueIndirect(v) {
		m["outputFormat"] = v
	}
	if v := f.ClientTags; v != nil {
		m["clientTags"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, f.LoggingConfig); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPrestoJob flattens an instance of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJob(c *Client, i interface{}) *WorkflowTemplateJobsPrestoJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPrestoJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPrestoJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsPrestoJobQueryList(c, m["queryList"])
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.OutputFormat = dcl.FlattenString(m["outputFormat"])
	r.ClientTags = dcl.FlattenStringSlice(m["clientTags"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, m["loggingConfig"])

	return r
}

// expandWorkflowTemplateJobsPrestoJobQueryListMap expands the contents of WorkflowTemplateJobsPrestoJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJobQueryList) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPrestoJobQueryListSlice expands the contents of WorkflowTemplateJobsPrestoJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobQueryListSlice(c *Client, f []WorkflowTemplateJobsPrestoJobQueryList) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobQueryListMap flattens the contents of WorkflowTemplateJobsPrestoJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobQueryListMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPrestoJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJobQueryList(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobQueryListSlice flattens the contents of WorkflowTemplateJobsPrestoJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobQueryListSlice(c *Client, i interface{}) []WorkflowTemplateJobsPrestoJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsPrestoJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJobQueryList(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJobQueryList expands an instance of WorkflowTemplateJobsPrestoJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobQueryList(c *Client, f *WorkflowTemplateJobsPrestoJobQueryList) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Queries; v != nil {
		m["queries"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPrestoJobQueryList flattens an instance of WorkflowTemplateJobsPrestoJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobQueryList(c *Client, i interface{}) *WorkflowTemplateJobsPrestoJobQueryList {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPrestoJobQueryList{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPrestoJobQueryList
	}
	r.Queries = dcl.FlattenStringSlice(m["queries"])

	return r
}

// expandWorkflowTemplateJobsPrestoJobLoggingConfigMap expands the contents of WorkflowTemplateJobsPrestoJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJobLoggingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsPrestoJobLoggingConfigSlice expands the contents of WorkflowTemplateJobsPrestoJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPrestoJobLoggingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPrestoJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobLoggingConfigMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsPrestoJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPrestoJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobLoggingConfigSlice(c *Client, i interface{}) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPrestoJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJobLoggingConfig expands an instance of WorkflowTemplateJobsPrestoJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPrestoJobLoggingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DriverLogLevels; !dcl.IsEmptyValueIndirect(v) {
		m["driverLogLevels"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPrestoJobLoggingConfig flattens an instance of WorkflowTemplateJobsPrestoJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c *Client, i interface{}) *WorkflowTemplateJobsPrestoJobLoggingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPrestoJobLoggingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPrestoJobLoggingConfig
	}
	r.DriverLogLevels = dcl.FlattenKeyValuePairs(m["driverLogLevels"])

	return r
}

// expandWorkflowTemplateJobsSchedulingMap expands the contents of WorkflowTemplateJobsScheduling into a JSON
// request object.
func expandWorkflowTemplateJobsSchedulingMap(c *Client, f map[string]WorkflowTemplateJobsScheduling) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsScheduling(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateJobsSchedulingSlice expands the contents of WorkflowTemplateJobsScheduling into a JSON
// request object.
func expandWorkflowTemplateJobsSchedulingSlice(c *Client, f []WorkflowTemplateJobsScheduling) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsScheduling(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSchedulingMap flattens the contents of WorkflowTemplateJobsScheduling from a JSON
// response object.
func flattenWorkflowTemplateJobsSchedulingMap(c *Client, i interface{}) map[string]WorkflowTemplateJobsScheduling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsScheduling{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsScheduling{}
	}

	items := make(map[string]WorkflowTemplateJobsScheduling)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsScheduling(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateJobsSchedulingSlice flattens the contents of WorkflowTemplateJobsScheduling from a JSON
// response object.
func flattenWorkflowTemplateJobsSchedulingSlice(c *Client, i interface{}) []WorkflowTemplateJobsScheduling {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsScheduling{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsScheduling{}
	}

	items := make([]WorkflowTemplateJobsScheduling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsScheduling(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateJobsScheduling expands an instance of WorkflowTemplateJobsScheduling into a JSON
// request object.
func expandWorkflowTemplateJobsScheduling(c *Client, f *WorkflowTemplateJobsScheduling) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MaxFailuresPerHour; !dcl.IsEmptyValueIndirect(v) {
		m["maxFailuresPerHour"] = v
	}
	if v := f.MaxFailuresTotal; !dcl.IsEmptyValueIndirect(v) {
		m["maxFailuresTotal"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsScheduling flattens an instance of WorkflowTemplateJobsScheduling from a JSON
// response object.
func flattenWorkflowTemplateJobsScheduling(c *Client, i interface{}) *WorkflowTemplateJobsScheduling {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsScheduling{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsScheduling
	}
	r.MaxFailuresPerHour = dcl.FlattenInteger(m["maxFailuresPerHour"])
	r.MaxFailuresTotal = dcl.FlattenInteger(m["maxFailuresTotal"])

	return r
}

// expandWorkflowTemplateParametersMap expands the contents of WorkflowTemplateParameters into a JSON
// request object.
func expandWorkflowTemplateParametersMap(c *Client, f map[string]WorkflowTemplateParameters) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParameters(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateParametersSlice expands the contents of WorkflowTemplateParameters into a JSON
// request object.
func expandWorkflowTemplateParametersSlice(c *Client, f []WorkflowTemplateParameters) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParameters(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersMap flattens the contents of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParametersMap(c *Client, i interface{}) map[string]WorkflowTemplateParameters {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParameters{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParameters{}
	}

	items := make(map[string]WorkflowTemplateParameters)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParameters(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateParametersSlice flattens the contents of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParametersSlice(c *Client, i interface{}) []WorkflowTemplateParameters {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParameters{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParameters{}
	}

	items := make([]WorkflowTemplateParameters, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParameters(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateParameters expands an instance of WorkflowTemplateParameters into a JSON
// request object.
func expandWorkflowTemplateParameters(c *Client, f *WorkflowTemplateParameters) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Fields; v != nil {
		m["fields"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v, err := expandWorkflowTemplateParametersValidation(c, f.Validation); err != nil {
		return nil, fmt.Errorf("error expanding Validation into validation: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["validation"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParameters flattens an instance of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParameters(c *Client, i interface{}) *WorkflowTemplateParameters {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateParameters{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateParameters
	}
	r.Name = dcl.FlattenString(m["name"])
	r.Fields = dcl.FlattenStringSlice(m["fields"])
	r.Description = dcl.FlattenString(m["description"])
	r.Validation = flattenWorkflowTemplateParametersValidation(c, m["validation"])

	return r
}

// expandWorkflowTemplateParametersValidationMap expands the contents of WorkflowTemplateParametersValidation into a JSON
// request object.
func expandWorkflowTemplateParametersValidationMap(c *Client, f map[string]WorkflowTemplateParametersValidation) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidation(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateParametersValidationSlice expands the contents of WorkflowTemplateParametersValidation into a JSON
// request object.
func expandWorkflowTemplateParametersValidationSlice(c *Client, f []WorkflowTemplateParametersValidation) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidation(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationMap flattens the contents of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationMap(c *Client, i interface{}) map[string]WorkflowTemplateParametersValidation {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidation{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidation{}
	}

	items := make(map[string]WorkflowTemplateParametersValidation)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidation(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateParametersValidationSlice flattens the contents of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationSlice(c *Client, i interface{}) []WorkflowTemplateParametersValidation {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidation{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidation{}
	}

	items := make([]WorkflowTemplateParametersValidation, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidation(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateParametersValidation expands an instance of WorkflowTemplateParametersValidation into a JSON
// request object.
func expandWorkflowTemplateParametersValidation(c *Client, f *WorkflowTemplateParametersValidation) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkflowTemplateParametersValidationRegex(c, f.Regex); err != nil {
		return nil, fmt.Errorf("error expanding Regex into regex: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["regex"] = v
	}
	if v, err := expandWorkflowTemplateParametersValidationValues(c, f.Values); err != nil {
		return nil, fmt.Errorf("error expanding Values into values: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["values"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParametersValidation flattens an instance of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidation(c *Client, i interface{}) *WorkflowTemplateParametersValidation {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateParametersValidation{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateParametersValidation
	}
	r.Regex = flattenWorkflowTemplateParametersValidationRegex(c, m["regex"])
	r.Values = flattenWorkflowTemplateParametersValidationValues(c, m["values"])

	return r
}

// expandWorkflowTemplateParametersValidationRegexMap expands the contents of WorkflowTemplateParametersValidationRegex into a JSON
// request object.
func expandWorkflowTemplateParametersValidationRegexMap(c *Client, f map[string]WorkflowTemplateParametersValidationRegex) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidationRegex(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateParametersValidationRegexSlice expands the contents of WorkflowTemplateParametersValidationRegex into a JSON
// request object.
func expandWorkflowTemplateParametersValidationRegexSlice(c *Client, f []WorkflowTemplateParametersValidationRegex) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidationRegex(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationRegexMap flattens the contents of WorkflowTemplateParametersValidationRegex from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationRegexMap(c *Client, i interface{}) map[string]WorkflowTemplateParametersValidationRegex {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidationRegex{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidationRegex{}
	}

	items := make(map[string]WorkflowTemplateParametersValidationRegex)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidationRegex(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateParametersValidationRegexSlice flattens the contents of WorkflowTemplateParametersValidationRegex from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationRegexSlice(c *Client, i interface{}) []WorkflowTemplateParametersValidationRegex {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidationRegex{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidationRegex{}
	}

	items := make([]WorkflowTemplateParametersValidationRegex, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidationRegex(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateParametersValidationRegex expands an instance of WorkflowTemplateParametersValidationRegex into a JSON
// request object.
func expandWorkflowTemplateParametersValidationRegex(c *Client, f *WorkflowTemplateParametersValidationRegex) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Regexes; v != nil {
		m["regexes"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParametersValidationRegex flattens an instance of WorkflowTemplateParametersValidationRegex from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationRegex(c *Client, i interface{}) *WorkflowTemplateParametersValidationRegex {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateParametersValidationRegex{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateParametersValidationRegex
	}
	r.Regexes = dcl.FlattenStringSlice(m["regexes"])

	return r
}

// expandWorkflowTemplateParametersValidationValuesMap expands the contents of WorkflowTemplateParametersValidationValues into a JSON
// request object.
func expandWorkflowTemplateParametersValidationValuesMap(c *Client, f map[string]WorkflowTemplateParametersValidationValues) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidationValues(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplateParametersValidationValuesSlice expands the contents of WorkflowTemplateParametersValidationValues into a JSON
// request object.
func expandWorkflowTemplateParametersValidationValuesSlice(c *Client, f []WorkflowTemplateParametersValidationValues) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidationValues(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationValuesMap flattens the contents of WorkflowTemplateParametersValidationValues from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationValuesMap(c *Client, i interface{}) map[string]WorkflowTemplateParametersValidationValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidationValues{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidationValues{}
	}

	items := make(map[string]WorkflowTemplateParametersValidationValues)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidationValues(c, item.(map[string]interface{}))
	}

	return items
}

// flattenWorkflowTemplateParametersValidationValuesSlice flattens the contents of WorkflowTemplateParametersValidationValues from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationValuesSlice(c *Client, i interface{}) []WorkflowTemplateParametersValidationValues {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidationValues{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidationValues{}
	}

	items := make([]WorkflowTemplateParametersValidationValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidationValues(c, item.(map[string]interface{})))
	}

	return items
}

// expandWorkflowTemplateParametersValidationValues expands an instance of WorkflowTemplateParametersValidationValues into a JSON
// request object.
func expandWorkflowTemplateParametersValidationValues(c *Client, f *WorkflowTemplateParametersValidationValues) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Values; v != nil {
		m["values"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParametersValidationValues flattens an instance of WorkflowTemplateParametersValidationValues from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationValues(c *Client, i interface{}) *WorkflowTemplateParametersValidationValues {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateParametersValidationValues{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateParametersValidationValues
	}
	r.Values = dcl.FlattenStringSlice(m["values"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *WorkflowTemplate) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalWorkflowTemplate(b, c)
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

type workflowTemplateDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         workflowTemplateApiOperation
}

func convertFieldDiffsToWorkflowTemplateDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]workflowTemplateDiff, error) {
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
	var diffs []workflowTemplateDiff
	// For each operation name, create a workflowTemplateDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := workflowTemplateDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToWorkflowTemplateApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToWorkflowTemplateApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (workflowTemplateApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractWorkflowTemplateFields(r *WorkflowTemplate) error {
	vPlacement := r.Placement
	if vPlacement == nil {
		// note: explicitly not the empty object.
		vPlacement = &WorkflowTemplatePlacement{}
	}
	if err := extractWorkflowTemplatePlacementFields(r, vPlacement); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPlacement) {
		r.Placement = vPlacement
	}
	return nil
}
func extractWorkflowTemplatePlacementFields(r *WorkflowTemplate, o *WorkflowTemplatePlacement) error {
	vManagedCluster := o.ManagedCluster
	if vManagedCluster == nil {
		// note: explicitly not the empty object.
		vManagedCluster = &WorkflowTemplatePlacementManagedCluster{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterFields(r, vManagedCluster); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vManagedCluster) {
		o.ManagedCluster = vManagedCluster
	}
	vClusterSelector := o.ClusterSelector
	if vClusterSelector == nil {
		// note: explicitly not the empty object.
		vClusterSelector = &WorkflowTemplatePlacementClusterSelector{}
	}
	if err := extractWorkflowTemplatePlacementClusterSelectorFields(r, vClusterSelector); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vClusterSelector) {
		o.ClusterSelector = vClusterSelector
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedCluster) error {
	// *ClusterClusterConfig is a reused type - that's not compatible with function extractors.

	return nil
}
func extractWorkflowTemplatePlacementClusterSelectorFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementClusterSelector) error {
	return nil
}
func extractWorkflowTemplateJobsFields(r *WorkflowTemplate, o *WorkflowTemplateJobs) error {
	vHadoopJob := o.HadoopJob
	if vHadoopJob == nil {
		// note: explicitly not the empty object.
		vHadoopJob = &WorkflowTemplateJobsHadoopJob{}
	}
	if err := extractWorkflowTemplateJobsHadoopJobFields(r, vHadoopJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vHadoopJob) {
		o.HadoopJob = vHadoopJob
	}
	vSparkJob := o.SparkJob
	if vSparkJob == nil {
		// note: explicitly not the empty object.
		vSparkJob = &WorkflowTemplateJobsSparkJob{}
	}
	if err := extractWorkflowTemplateJobsSparkJobFields(r, vSparkJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkJob) {
		o.SparkJob = vSparkJob
	}
	vPysparkJob := o.PysparkJob
	if vPysparkJob == nil {
		// note: explicitly not the empty object.
		vPysparkJob = &WorkflowTemplateJobsPysparkJob{}
	}
	if err := extractWorkflowTemplateJobsPysparkJobFields(r, vPysparkJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPysparkJob) {
		o.PysparkJob = vPysparkJob
	}
	vHiveJob := o.HiveJob
	if vHiveJob == nil {
		// note: explicitly not the empty object.
		vHiveJob = &WorkflowTemplateJobsHiveJob{}
	}
	if err := extractWorkflowTemplateJobsHiveJobFields(r, vHiveJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vHiveJob) {
		o.HiveJob = vHiveJob
	}
	vPigJob := o.PigJob
	if vPigJob == nil {
		// note: explicitly not the empty object.
		vPigJob = &WorkflowTemplateJobsPigJob{}
	}
	if err := extractWorkflowTemplateJobsPigJobFields(r, vPigJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPigJob) {
		o.PigJob = vPigJob
	}
	vSparkRJob := o.SparkRJob
	if vSparkRJob == nil {
		// note: explicitly not the empty object.
		vSparkRJob = &WorkflowTemplateJobsSparkRJob{}
	}
	if err := extractWorkflowTemplateJobsSparkRJobFields(r, vSparkRJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkRJob) {
		o.SparkRJob = vSparkRJob
	}
	vSparkSqlJob := o.SparkSqlJob
	if vSparkSqlJob == nil {
		// note: explicitly not the empty object.
		vSparkSqlJob = &WorkflowTemplateJobsSparkSqlJob{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobFields(r, vSparkSqlJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkSqlJob) {
		o.SparkSqlJob = vSparkSqlJob
	}
	vPrestoJob := o.PrestoJob
	if vPrestoJob == nil {
		// note: explicitly not the empty object.
		vPrestoJob = &WorkflowTemplateJobsPrestoJob{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobFields(r, vPrestoJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPrestoJob) {
		o.PrestoJob = vPrestoJob
	}
	vScheduling := o.Scheduling
	if vScheduling == nil {
		// note: explicitly not the empty object.
		vScheduling = &WorkflowTemplateJobsScheduling{}
	}
	if err := extractWorkflowTemplateJobsSchedulingFields(r, vScheduling); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vScheduling) {
		o.Scheduling = vScheduling
	}
	return nil
}
func extractWorkflowTemplateJobsHadoopJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHadoopJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsHadoopJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsHadoopJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHadoopJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsSparkJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsSparkJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsPysparkJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPysparkJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPysparkJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsPysparkJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPysparkJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsHiveJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHiveJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsHiveJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsHiveJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	return nil
}
func extractWorkflowTemplateJobsHiveJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHiveJobQueryList) error {
	return nil
}
func extractWorkflowTemplateJobsPigJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsPigJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsPigJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPigJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPigJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsPigJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJobQueryList) error {
	return nil
}
func extractWorkflowTemplateJobsPigJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsSparkRJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkRJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkRJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsSparkRJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkRJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsSparkSqlJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsSparkSqlJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsSparkSqlJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJobQueryList) error {
	return nil
}
func extractWorkflowTemplateJobsSparkSqlJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsPrestoJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsPrestoJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func extractWorkflowTemplateJobsPrestoJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJobQueryList) error {
	return nil
}
func extractWorkflowTemplateJobsPrestoJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJobLoggingConfig) error {
	return nil
}
func extractWorkflowTemplateJobsSchedulingFields(r *WorkflowTemplate, o *WorkflowTemplateJobsScheduling) error {
	return nil
}
func extractWorkflowTemplateParametersFields(r *WorkflowTemplate, o *WorkflowTemplateParameters) error {
	vValidation := o.Validation
	if vValidation == nil {
		// note: explicitly not the empty object.
		vValidation = &WorkflowTemplateParametersValidation{}
	}
	if err := extractWorkflowTemplateParametersValidationFields(r, vValidation); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vValidation) {
		o.Validation = vValidation
	}
	return nil
}
func extractWorkflowTemplateParametersValidationFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidation) error {
	vRegex := o.Regex
	if vRegex == nil {
		// note: explicitly not the empty object.
		vRegex = &WorkflowTemplateParametersValidationRegex{}
	}
	if err := extractWorkflowTemplateParametersValidationRegexFields(r, vRegex); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRegex) {
		o.Regex = vRegex
	}
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &WorkflowTemplateParametersValidationValues{}
	}
	if err := extractWorkflowTemplateParametersValidationValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vValues) {
		o.Values = vValues
	}
	return nil
}
func extractWorkflowTemplateParametersValidationRegexFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidationRegex) error {
	return nil
}
func extractWorkflowTemplateParametersValidationValuesFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidationValues) error {
	return nil
}

func postReadExtractWorkflowTemplateFields(r *WorkflowTemplate) error {
	vPlacement := r.Placement
	if vPlacement == nil {
		// note: explicitly not the empty object.
		vPlacement = &WorkflowTemplatePlacement{}
	}
	if err := postReadExtractWorkflowTemplatePlacementFields(r, vPlacement); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPlacement) {
		r.Placement = vPlacement
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementFields(r *WorkflowTemplate, o *WorkflowTemplatePlacement) error {
	vManagedCluster := o.ManagedCluster
	if vManagedCluster == nil {
		// note: explicitly not the empty object.
		vManagedCluster = &WorkflowTemplatePlacementManagedCluster{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterFields(r, vManagedCluster); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vManagedCluster) {
		o.ManagedCluster = vManagedCluster
	}
	vClusterSelector := o.ClusterSelector
	if vClusterSelector == nil {
		// note: explicitly not the empty object.
		vClusterSelector = &WorkflowTemplatePlacementClusterSelector{}
	}
	if err := extractWorkflowTemplatePlacementClusterSelectorFields(r, vClusterSelector); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vClusterSelector) {
		o.ClusterSelector = vClusterSelector
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedCluster) error {
	// *ClusterClusterConfig is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractWorkflowTemplatePlacementClusterSelectorFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementClusterSelector) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsFields(r *WorkflowTemplate, o *WorkflowTemplateJobs) error {
	vHadoopJob := o.HadoopJob
	if vHadoopJob == nil {
		// note: explicitly not the empty object.
		vHadoopJob = &WorkflowTemplateJobsHadoopJob{}
	}
	if err := extractWorkflowTemplateJobsHadoopJobFields(r, vHadoopJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vHadoopJob) {
		o.HadoopJob = vHadoopJob
	}
	vSparkJob := o.SparkJob
	if vSparkJob == nil {
		// note: explicitly not the empty object.
		vSparkJob = &WorkflowTemplateJobsSparkJob{}
	}
	if err := extractWorkflowTemplateJobsSparkJobFields(r, vSparkJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkJob) {
		o.SparkJob = vSparkJob
	}
	vPysparkJob := o.PysparkJob
	if vPysparkJob == nil {
		// note: explicitly not the empty object.
		vPysparkJob = &WorkflowTemplateJobsPysparkJob{}
	}
	if err := extractWorkflowTemplateJobsPysparkJobFields(r, vPysparkJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPysparkJob) {
		o.PysparkJob = vPysparkJob
	}
	vHiveJob := o.HiveJob
	if vHiveJob == nil {
		// note: explicitly not the empty object.
		vHiveJob = &WorkflowTemplateJobsHiveJob{}
	}
	if err := extractWorkflowTemplateJobsHiveJobFields(r, vHiveJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vHiveJob) {
		o.HiveJob = vHiveJob
	}
	vPigJob := o.PigJob
	if vPigJob == nil {
		// note: explicitly not the empty object.
		vPigJob = &WorkflowTemplateJobsPigJob{}
	}
	if err := extractWorkflowTemplateJobsPigJobFields(r, vPigJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPigJob) {
		o.PigJob = vPigJob
	}
	vSparkRJob := o.SparkRJob
	if vSparkRJob == nil {
		// note: explicitly not the empty object.
		vSparkRJob = &WorkflowTemplateJobsSparkRJob{}
	}
	if err := extractWorkflowTemplateJobsSparkRJobFields(r, vSparkRJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkRJob) {
		o.SparkRJob = vSparkRJob
	}
	vSparkSqlJob := o.SparkSqlJob
	if vSparkSqlJob == nil {
		// note: explicitly not the empty object.
		vSparkSqlJob = &WorkflowTemplateJobsSparkSqlJob{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobFields(r, vSparkSqlJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSparkSqlJob) {
		o.SparkSqlJob = vSparkSqlJob
	}
	vPrestoJob := o.PrestoJob
	if vPrestoJob == nil {
		// note: explicitly not the empty object.
		vPrestoJob = &WorkflowTemplateJobsPrestoJob{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobFields(r, vPrestoJob); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPrestoJob) {
		o.PrestoJob = vPrestoJob
	}
	vScheduling := o.Scheduling
	if vScheduling == nil {
		// note: explicitly not the empty object.
		vScheduling = &WorkflowTemplateJobsScheduling{}
	}
	if err := extractWorkflowTemplateJobsSchedulingFields(r, vScheduling); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vScheduling) {
		o.Scheduling = vScheduling
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsHadoopJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHadoopJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsHadoopJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsHadoopJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHadoopJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsPysparkJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPysparkJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPysparkJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsPysparkJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPysparkJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsHiveJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHiveJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsHiveJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsHiveJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsHiveJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsHiveJobQueryList) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsPigJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsPigJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsPigJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPigJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPigJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsPigJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJobQueryList) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsPigJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPigJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkRJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkRJob) error {
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkRJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkRJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkRJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkSqlJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsSparkSqlJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsSparkSqlJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkSqlJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJobQueryList) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsSparkSqlJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsSparkSqlJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsPrestoJobFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJob) error {
	vQueryList := o.QueryList
	if vQueryList == nil {
		// note: explicitly not the empty object.
		vQueryList = &WorkflowTemplateJobsPrestoJobQueryList{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobQueryListFields(r, vQueryList); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vQueryList) {
		o.QueryList = vQueryList
	}
	vLoggingConfig := o.LoggingConfig
	if vLoggingConfig == nil {
		// note: explicitly not the empty object.
		vLoggingConfig = &WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}
	if err := extractWorkflowTemplateJobsPrestoJobLoggingConfigFields(r, vLoggingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLoggingConfig) {
		o.LoggingConfig = vLoggingConfig
	}
	return nil
}
func postReadExtractWorkflowTemplateJobsPrestoJobQueryListFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJobQueryList) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsPrestoJobLoggingConfigFields(r *WorkflowTemplate, o *WorkflowTemplateJobsPrestoJobLoggingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplateJobsSchedulingFields(r *WorkflowTemplate, o *WorkflowTemplateJobsScheduling) error {
	return nil
}
func postReadExtractWorkflowTemplateParametersFields(r *WorkflowTemplate, o *WorkflowTemplateParameters) error {
	vValidation := o.Validation
	if vValidation == nil {
		// note: explicitly not the empty object.
		vValidation = &WorkflowTemplateParametersValidation{}
	}
	if err := extractWorkflowTemplateParametersValidationFields(r, vValidation); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vValidation) {
		o.Validation = vValidation
	}
	return nil
}
func postReadExtractWorkflowTemplateParametersValidationFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidation) error {
	vRegex := o.Regex
	if vRegex == nil {
		// note: explicitly not the empty object.
		vRegex = &WorkflowTemplateParametersValidationRegex{}
	}
	if err := extractWorkflowTemplateParametersValidationRegexFields(r, vRegex); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRegex) {
		o.Regex = vRegex
	}
	vValues := o.Values
	if vValues == nil {
		// note: explicitly not the empty object.
		vValues = &WorkflowTemplateParametersValidationValues{}
	}
	if err := extractWorkflowTemplateParametersValidationValuesFields(r, vValues); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vValues) {
		o.Values = vValues
	}
	return nil
}
func postReadExtractWorkflowTemplateParametersValidationRegexFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidationRegex) error {
	return nil
}
func postReadExtractWorkflowTemplateParametersValidationValuesFields(r *WorkflowTemplate, o *WorkflowTemplateParametersValidationValues) error {
	return nil
}
